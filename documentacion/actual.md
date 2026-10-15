### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (15/10/2026):** según ClickUp, en este archivo quedan las tareas **pendientes de resolución** de la Etapa 1 y los `FIX` abiertos. Lo ya verificado con las suites está en [`terminado.md`](terminado.md) / [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md), y las correcciones de las implementadas con problemas, en [`futuro.md`](futuro.md) / [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md).
>
> **Pasados a `terminado-1.md` el 15/10/2026:**
> - `BAC-26`, `BAC-32`, `BAC-33`, `FIX-84`, `FIX-74`, `FIX-87` (resueltos y con tests en verde).
> - `FRN-20C`, `FIX-89` (implementados; correcciones abiertas en `FIX-97` y `FIX-99` de `futuro-1.md`).
>
> **Pasados a `terminado2.md` el 15/10/2026 (Etapa 2):**
> - `BAC-34`, `BAC-35`, `BAC-37`, `BAC-38` (resueltos y con tests en verde).

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

> **Cobertura incorporada del PR #3:** las pruebas nuevas comprueban el comportamiento actual: conflicto al reutilizar correo o username de un usuario dado de baja, y alta con datos distintos. No implementan ni validan la reutilización del correo requerida por `FIX-54` a `FIX-57`; esa regresión debe cambiar cuando se implemente el nuevo contrato.

> **Nota de alcance (06/10/2026, revisión estática):** “recorrido completo” se refiere a los casos anteriores, no a crear una cuenta nueva con el correo de una eliminada y preservar la identidad histórica. Esta brecha de la fase base queda pendiente en `FIX-54` a `FIX-57` de [`futuro.md`](futuro.md#fix-54---separar-eliminación-lógica-y-suspensión-migrar-unicidad-del-correo-backend). No se ejecutaron suites ni se verificó el esquema desplegado.

---

#### `BAC-24C` - Limpieza de permisos al eliminar una instancia (nueva)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.0 h
- **Depende de:** `BAC-24B` y `BAC-25B`.
- **Problema:** `/cluster/nextid` devuelve el menor VMID libre. Si se borra la 110 y después se crea otra instancia, recibe el 110 y **hereda las filas de `permisos_instancia`** de la borrada: un operador tendría acceso a una máquina que no se le asignó.
- **Entregable:**
  - Cuando una tarea de borrado termina con `COMPLETED`, eliminar en una transacción las filas de `permisos_instancia` de ese VMID.
  - Registrar la limpieza en `auditoria` con `accion: INSTANCE_PERMISSIONS_PURGED` y la lista de usuarios afectados.
  - Publicar el `TASK_FINISHED` **antes** de borrar los permisos, para que los operadores asignados lo reciban.
  - Si la tarea falla, no se toca nada.
- **Criterio de éxito:** después de borrar la 110, `GET /api/admin/users/:id/permissions` del operador ya no la incluye. Al crear una instancia nueva con el VMID 110, el operador no la ve.



#### `FIX-81` - Corregir el orden de migraciones: `migrarUnicidadEmail` antes de crear `auditoria` (`FIX-54`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro / Tayra
- **Estimación:** 0,5 h
- **Depende de:** `FIX-54` (en `terminado.md`).
- **Problema técnico detallado:**
  En `backend/internal/adapters/secondary/postgres/db.go`, después de `AutoMigrate` se llama `migrarUnicidadEmail(db)` **antes** de `migrarAuditoriaParticionada(db)`. `migrarUnicidadEmail` → `clasificarEliminados` ejecuta `SELECT … FROM auditoria` (`migrar_usuarios.go`), pero la tabla `auditoria` (particionada) la crea `migrarAuditoriaParticionada`. En una **base nueva** (deploy limpio o el stack de `test/back`, que parte de una base vacía) el arranque aborta:
  ```
  ❌ Error fatal al iniciar la base de datos: error al clasificar usuarios eliminados:
     ERROR: relation "auditoria" does not exist (SQLSTATE 42P01)
  ```
  En el servidor de pruebas no se ve porque `auditoria` ya existía; en un entorno limpio el backend sale con código 1 y **toda** la suite de aceptación del backend (incluyendo la Etapa 1) queda bloqueada.
- **Qué debe hacer el equipo de backend:**
  1. En `db.go`, ejecutar `migrarUnicidadEmail(db)` **después** de `migrarAuditoriaParticionada(db)` y `asegurarParticionesFuturas(db)`.
  2. Alternativa equivalente: que `clasificarEliminados` no consulte `auditoria` si la tabla todavía no existe (guarda con `existeTabla`).
  3. Verificar con una base **vacía** (`docker compose down --volumes` + `up`) que el backend arranca y que la suite `test/back` corre completa sin parches.
- **Criterio de éxito:**
  - El backend arranca contra una base nueva sin errores y sin perder la clasificación de usuarios eliminados.
  - `test/back` deja de requerir el reordenamiento manual.

---

#### `FRN-21` (`BRG-06-FRN`) - Maquetado de la vista Detalle de Instancia (`/instances/:vmid`) y enlace desde notificaciones (`RF-03`, `RF-11`)

> [!WARNING]
> **Estado verificado (14/10/2026, frontend `add947f`): implementada parcialmente, con problema pendiente.** La vista existe pero falta registrar `/instances/:vmid` en `AppRoutes.tsx` y enlazar filas en `Instances.tsx` (`instance-detail.test.tsx`). Corrección en **`FIX-95`** ([futuro-1.md](futuro-1.md)).

- **Área:** Frontend
- **Asignado:** Cristian / Belinda
- **Estimación:** 2.5 h
- **Depende de:** `FRN-20A`, `FRN-17B`.
- **Problema:** En `etapa1.md`, la nota 45 pospuso el enlace de los toasts al detalle de la instancia porque la vista no existía. Sin esta pantalla, la Etapa 2 no tiene contenedor donde montar los gráficos de métricas (`RF-05`) ni la edición de recursos (`RF-10`).
- **Entregable:**
  1. Ruta `/instances/:vmid` bajo `ProtectedLayout` con guard de sesión y permiso de acceso.
  2. Componente de vista `InstanceDetailPage.tsx` con:
     - Cabecera con nombre, ID, tipo, IP, badge de estado y botonera de acciones (`InstanceAction.tsx`).
     - Navegación por pestañas: "General" (información básica y métricas estáticas), "Métricas" (contenedor placeholder para gráficos de Etapa 2) y "Configuración" (contenedor placeholder para edición de Etapa 2).
     - Botón de retorno / breadcrumb hacia `/instances`.
  3. Vincular la tabla de inventario (`Instances.tsx`) y los toasts de notificación (`FRN-17B`) con enlace directo a `/instances/:vmid`.
- **Criterio de éxito:** Al hacer clic en una fila o en un toast de tarea completada, la aplicación navega a `/instances/:vmid` mostrando los datos reales de la máquina y sus pestañas base listas.



### `FIX-79` - Exposición directa de controles para selector reactivo por tipo de recurso (`FRN-20B`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 1 h
- **Depende de:** `FRN-20B`.
- **Problema técnico detallado:**
  La lógica de filtrado reactivo en memoria por tipo (`VM`, `LXC`, `all`) y buscador ya fue implementada en `frontend/centinela/src/pages/Instances.tsx`. No obstante, el selector fue colocado dentro de un menú desplegable contextual (`DropdownMenu`), requiriendo clics previos para ser descubierto. La suite de pruebas `test/front/instances-ola2.test.tsx:120` exige que los controles por tipo de recurso estén expuestos de forma visible y directa en la barra de control (como pestañas `tab` o botones accesibles `button` con etiquetas legibles como "Todas", "VM", "LXC"):
  ```
  AssertionError: debe existir un control para filtrar reactivamente por máquinas virtuales (VM): expected null not to be null
  ```
- **Qué debe hacer el equipo de frontend:**
  1. En `Instances.tsx`, complementar o reemplazar el dropdown por botones o pestañas accesibles y visibles en la barra de inventario que permitan alternar directamente entre "Todas", "VM (Máquinas virtuales)" y "LXC (Contenedores)".
  2. Mantener la reactividad instantánea en memoria sin peticiones al backend y la actualización sincronizada de los contadores.
  3. Ejecutar `npx vitest run instances-ola2.test.tsx` y certificar que la suite pase al 100%.
- **Criterio de éxito:**
  - El selector por tipo de recurso está directamente visible y accesible sin pasos ocultos.
  - La suite `instances-ola2.test.tsx` pasa 6 de 6 pruebas en verde.


### `FIX-76` - Desacoplar guarda de `nivelAcceso` para rol `ADMIN` en modales de energía (`FRN-15`) (Frontend)

> [!WARNING]
> **Estado verificado (13/10/2026, frontend `7fcc0bd`): implementada parcialmente, con problema pendiente.** `instances-modals.test.tsx` quedó con **1 fallo** (`Start pide confirmación y cancelar no envía ninguna petición`). Corrección en **`FIX-85`** ([futuro-1.md](futuro-1.md)).


- **Área:** Frontend
- **Asignado:** Belinda / Luz
- **Estimación:** 0,5 h
- **Depende de:** `FRN-15`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/components/features/instances/components/InstanceAction.tsx:87`:
  ```tsx
  const canOperate = canOperateInstance(instance.id) && instance.nivelAcceso !== 'READ_ONLY';
  ```
  Un usuario con rol `ADMIN` dispone de privilegios absolutos para operar cualquier instancia sin estar restringido por las listas de acceso granular de operadores (`nivelAcceso`). Sin embargo, la condición evalúa `instance.nivelAcceso !== 'READ_ONLY'` sin contemplar `isAdmin()`. Cuando un fixture o instancia mockeada contiene `nivelAcceso: 'READ_ONLY'` (asignado para un perfil de operador), un administrador queda bloqueado y la acción `start` ("Encender") no se renderiza en la fila ni abre su modal de confirmación, provocando el fallo en `instances-modals.test.tsx`:
  ```
  AssertionError: la fila no ofrece la acción "start" (se busca un botón o ítem de menú con nombre /^(iniciar|encender|start)\b/i)
  ```
- **Qué debe hacer el equipo de frontend:**
  1. En `InstanceAction.tsx`, actualizar la asignación de `canOperate` para que un administrador siempre esté habilitado a operar:
     ```tsx
     const canOperate = isAdmin() || (canOperateInstance(instance.id) && instance.nivelAcceso !== 'READ_ONLY');
     ```
  2. Aplicar la misma corrección en la validación interna de `confirmPowerAction` (`InstanceAction.tsx:114`).
  3. Ejecutar `npx vitest run instances-modals.test.tsx` y certificar que la prueba `Start pide confirmación` pase en verde.
- **Criterio de éxito:**
  - Un usuario `ADMIN` puede disparar la acción de inicio/encendido y ver su modal de confirmación sin ser bloqueado por `nivelAcceso`.
  - La suite `instances-modals.test.tsx` pasa al 100% (6/6 en verde).


#### `FRN-19C` (`BRG-05-FRN`, parte Dashboard) - Tarjetas de conteo de VMs/LXC con auto-actualización (`RF-02`)

> [!WARNING]
> **Estado verificado (13/10/2026, frontend `7fcc0bd`): implementada parcialmente, con problema pendiente.** Las tarjetas de conteo se renderizan desde `instancesSummary` (1 caso nuevo en verde), pero **no se refrescan al recibir `TASK_FINISHED`** (solo por el polling de 10 s) → 1 caso en rojo. Corrección en **`FIX-83`** ([futuro-1.md](futuro-1.md)).


- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 1.0 h
- **Depende de:** `BAC-22B`, `FRN-19B` y `FRN-17A`.
- **Entregable:** tarjetas de resumen de VMs y LXC (`En ejecución`, `Detenidas`, `Total`) con los datos de `instancesSummary`. Usan la misma consulta cada 10 s de `FRN-19B` y vuelven a consultar al recibir cualquier `TASK_FINISHED`.
- **Criterio de éxito:** después de encender una VM, el conteo cambia sin `F5` en menos de 2 s desde el evento.





#### `FRN-17B` - Desbloqueo reactivo de instancia y notificación Toast adaptativa

> [!NOTE]
> **Estado verificado (13/10/2026, frontend `7fcc0bd`): implementada y con test en verde.** `instances-events.test.tsx` (2/2) pasa: ante `TASK_FINISHED` se libera la fila y se muestra el toast (éxito en COMPLETED; error con el motivo en FAILED).

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `FRN-16`, `FRN-17A` y `BAC-26` (`exitstatus` en `detalles`).
- **Entregable:**
  - Al recibir un `TASK_FINISHED` con el mismo `detalles.tareaId` que una fila en transición, quitar el bloqueo y el spinner.
  - Actualizar el badge de estado (`Running` / `Stopped`).
  - Mostrar un toast (`components/ui/toast.tsx`):
    - verde si `estado: COMPLETED`;
    - rojo si `FAILED`, con un mensaje distinto según `detalles.motivo` (D2):
      - `PROXMOX_ERROR`: "Proxmox no pudo completar la acción", con `detalles.error`;
      - `TIMEOUT`: "La acción no terminó en 3 minutos; revisá el estado de la instancia".
  - Si la acción fue `delete`, quitar la fila.
- **Criterio de éxito:** la interfaz actualiza el estado y libera los botones sin `F5`.

---




### `FIX-83` - Refrescar el conteo de instancias ante `TASK_FINISHED` (`FRN-19C`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 1 h
- **Depende de:** `FRN-19C` (en `terminado-1.md`), `FRN-19B`, `FRN-17A`.
- **Problema técnico detallado:**
  Las tarjetas de conteo del Dashboard (`FRN-19C`) se renderizan desde `instancesSummary` (1 caso de `dashboard-summary.test.tsx` en verde), pero **no se refrescan al recibir `TASK_FINISHED`**: `useNodeStatus` sólo hace polling cada 10 s y no observa el canal de eventos, así que el conteo puede tardar hasta 10 s en cambiar (el criterio pide menos de 2 s).
- **Qué debe hacer el equipo de frontend:**
  1. Suscribir el Dashboard al canal (`useEventsContext` / `useEventSubscription`) y, ante cualquier `TASK_FINISHED` (o `INSTANCE_STATE_CHANGED`), volver a consultar `GET /api/node/status` sin loader.
  2. Ejecutar `npx vitest run dashboard-summary.test.tsx` hasta 2/2.
- **Criterio de éxito:** tras encender o apagar una instancia, el conteo cambia en menos de 2 s sin `F5`.


### `FIX-85` - Cerrar la confirmación de `Start` en los modales de energía (`FIX-76` / `FRN-15`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 1 h
- **Depende de:** `FIX-76` (en `terminado-1.md`), `FRN-15`.
- **Problema técnico detallado:**
  `test/front/instances-modals.test.tsx` quedó con **1 fallo**: `Start pide confirmación y cancelar no envía ninguna petición`. El resto de los casos (shutdown/reboot/stop/delete y el filtro por rol) pasa.
- **Qué debe hacer el equipo de frontend:**
  1. Revisar `InstanceAction` / `ConfirmUserAction` para que `Start` abra el modal de confirmación y que cancelarlo **no** dispare ninguna petición (`POST /instances/:vmid/start` o `/status/start`).
  2. Ejecutar `npx vitest run instances-modals.test.tsx` hasta 6/6.
- **Criterio de éxito:** `instances-modals.test.tsx` pasa 6 de 6.


### `FIX-95` - Declarar ruta `/instances/:vmid` y enlazar inventario con la vista de detalle (`FRN-21`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-21`.
- **Problema técnico detallado:**
  La vista `InstanceDetailPage.tsx` fue maquetada con éxito bajo `src/pages/InstanceDetailPage.tsx` e incluye la estructura de cabecera y pestañas ("General", "Métricas", "Configuración"). Sin embargo, la integración global evaluada por `test/front/instance-detail.test.tsx` falla en 2 aserciones:
  1. La ruta protegida `/instances/:vmid` no está registrada en el enrutador principal (`frontend/centinela/src/routes/AppRoutes.tsx`).
  2. La tabla de inventario (`frontend/centinela/src/pages/Instances.tsx`) no contiene enlaces ni eventos de navegación para acceder al detalle de la máquina al hacer clic sobre la fila o el identificador/nombre de la instancia.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/routes/AppRoutes.tsx`, importar `InstanceDetailPage` y declarar la ruta `<Route path="/instances/:vmid" element={<InstanceDetailPage />} />` dentro de `ProtectedLayout`.
  2. En `frontend/centinela/src/pages/Instances.tsx`, hacer que el nombre o ID de cada fila sea un `<Link to={`/instances/${instance.id}`}>` o incorporar navegación programática con `useNavigate`.
  3. Ejecutar `npx vitest run instance-detail.test.tsx` en `test/front` y certificar que pase 3 de 3 en verde.
- **Criterio de éxito:**
  - La ruta `/instances/:vmid` es accesible para usuarios autenticados.
  - Al interactuar con una fila del inventario, la aplicación navega fluidamente hacia la vista de detalle.
  - `test/front/instance-detail.test.tsx` pasa al 100% (3/3 en verde).

### `FIX-86` - Completar la máquina de estados y el mapeo D2 de operaciones (`FIX-77` / `FRN-16`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2 h
- **Depende de:** `FIX-77` (en `terminado-1.md`), `FRN-16`, `BAC-29`.
- **Problema técnico detallado:**
  `test/front/instances-operation.test.tsx` queda con **3 fallos** (de 9) y **todos son del mapeo de errores `403`**.
  En `components/common/ConfirmUserAction.tsx` los `401/403` **no** se muestran en el modal (se delegan al aviso global): `if (!(error instanceof ApiRequestError && [401,403].includes(error.status))) setErrorMessage(...)`. Y el aviso global (`ApiResponseNotifier.handleForbidden`) muestra el **mensaje crudo** del backend con el título genérico `Acción rechazada`, porque el mapa `forbiddenTitles` sólo conoce `INSUFFICIENT_PERMISSIONS`, `NO_ROLE`, `INVALID_ROLE`, `2FA_REQUIRED` y `PASSWORD_CHANGE_REQUIRED`. `utils/errorMapper.ts` **sí** traduce `INSTANCE_PROTECTED` ("infraestructura crítica") e `INSTANCE_ACCESS_DENIED` ("permisos"), pero ese camino no lo usa.
- **Qué debe hacer el equipo de frontend:**
  1. Mostrar el mensaje diferenciado de estos códigos: usar `getErrorMessage(error)` en `ApiResponseNotifier.handleForbidden` (o agregar los códigos de instancia a `forbiddenTitles`).
  2. Alternativa equivalente: dejar que `ConfirmUserAction` muestre `getErrorMessage(error)` también para los `403` de instancia, sin delegarlos ciegamente al aviso global.
  3. Ejecutar `npx vitest run instances-operation.test.tsx` hasta 9/9.
- **Criterio de éxito:** `instances-operation.test.tsx` pasa 9 de 9.

---
ETAPA 2
---

#### `INF-10` - Verificación y disponibilidad de plantillas LXC e ISOs en el servidor Proxmox real

- **Área:** Infraestructura
- **Asignados:** Nico y Lucas
- **Estimación:** 1.5 h
- **Depende de:** `INF-07C` (permisos de token en Proxmox).
- **Problema:** Para que la validación en el servidor real (`INT-06`) funcione, los almacenamientos de Proxmox deben contener las imágenes base necesarias.
- **Entregable:**
  1. Descargar o verificar la presencia de una plantilla LXC liviana oficial (ej. `alpine-3.19-default` o `debian-12-standard`) en el storage `local:vztmpl`.
  2. Verificar una imagen ISO de prueba en `local:iso`.
  3. Comprobar que el API Token puede listar ambos contenidos mediante consulta HTTPS a `/storage/{storage}/content`.
- **Criterio de éxito:** El backend desplegado en el servidor lista las plantillas e ISOs reales a través de `GET /api/node/templates` sin errores de permisos.

#### `BAC-36` - Endpoint REST de métricas históricas de instancia (`timeframe=2h|24h`)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-34`.
- **Problema:** Al abrir la vista de una máquina, los gráficos necesitan renderizar de inmediato la tendencia histórica sin aguardar la llegada de los primeros paquetes del WebSocket.
- **Entregable:**
  1. Endpoint `GET /api/instances/:vmid/metrics?timeframe=2h|24h` protegido con `RequireAuth` y `RequireInstanceAccess`.
  2. Parámetro `timeframe` configurable (por defecto `2h`, soportando `24h`), mapeado a los intervalos de Proxmox (`hour` y `day`).
  3. Caché corta en Redis (TTL 30 s) para evitar consultas repetitivas al hipervisor ante recargas de página.
- **Criterio de éxito:** Responde en menos de 100 ms desde caché con el arreglo histórico de métricas formateado para el consumo de gráficos.

#### `BAC-39` - Seguimiento de creación, autoasignación de permisos y eventos (`RF-07`, `RF-11`)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-38`, `BAC-32` (worker con soporte de creación).
- **Problema:** Una máquina recién creada no pertenece a ningún operador por defecto. El creador debe obtener acceso inmediato sin requerir que otro administrador configure manualmente sus permisos en `/admin/users`.
- **Entregable:**
  1. Integración en `seguimiento_tareas.go`: al detectar fin exitoso (`exitstatus: "OK"`) de un UPID de creación, insertar atómicamente `(creador_id, vmid, FULL_ACCESS)` en `permisos_instancia`.
  2. Emisión por el canal SSE `/api/events` de `TASK_FINISHED`:
     `{ tipo: "TASK_FINISHED", recursoTipo: "VM", recursoId: "108", detalles: { tareaId, accion: "create", estado: "COMPLETED" } }`.
  3. Si la tarea de creación falla en Proxmox, liberar bloqueos y registrar el evento con `estado: "FAILED"`.
- **Criterio de éxito:** Al finalizar la creación, el inventario del usuario incluye la nueva máquina y la interfaz recibe la notificación en tiempo real.


#### `BAC-40` - Endpoint `PUT /api/instances/:vmid/config` con cálculo diferencial ($\Delta\text{RAM}$) y detección de hot-plug

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-37`.
- **Problema:** La modificación de recursos debe validar que los incrementos solicitados no saturen el nodo y debe advertir con precisión si el cambio exige reiniciar el sistema operativo huésped.
- **Entregable:**
  1. Endpoint `PUT /api/instances/:vmid/config` accesible para usuarios con rol `ADMIN` o con `FULL_ACCESS` sobre la instancia.
  2. Payload: `{ cores?: number, ramMb?: number }`.
  3. Cálculo diferencial: si $\Delta\text{RAM} > 0$, validar contra la cuota libre con `CuotasService`.
  4. Invocación de `PUT /nodes/{node}/{tipo}/{vmid}/config` en Proxmox.
  5. Inspección del resultado: si la instancia está `running` y no aplica en caliente, responder `HTTP 200 OK` con `{ requiresReboot: true, message: "Recursos modificados; los cambios se aplicarán tras reiniciar la instancia." }`. Si aplica en caliente, `{ requiresReboot: false }`.
- **Criterio de éxito:** Incrementos excesivos son rechazados con `409 RESOURCE_QUOTA_EXCEEDED`. Cambios válidos se persisten y notifican el estado de hot-plug.

#### `BAC-41` - Auditoría inmutable de aprovisionamiento y cambios de recursos (`RF-08`)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-38`, `BAC-40`.
- **Problema:** Toda alteración en la infraestructura debe quedar registrada inalterablemente en la tabla `auditoria` para trazabilidad institucional de los operadores.
- **Entregable:**
  1. Registro de `CREATE_INSTANCE`: entrada inicial con `resultado: "PENDING"` al recibir la orden, y entrada final con `resultado: "SUCCESS"` o `"FAILED"` vinculada al `upid` y `tareaId`.
  2. Registro de `CONFIG_INSTANCE`: entrada auditando los valores anteriores y los nuevos parámetros en `detalles` (JSONB).
  3. Verificación de inmutabilidad contra la base particionada.
- **Criterio de éxito:** `GET /api/admin/audit` refleja las acciones de creación y edición correlacionadas por identificador de tarea sin omisiones de datos.


#### `BAC-50` (`BRG-10-BAC`) - Extensión del simulador de Proxmox para Snapshots y Rollback
- **Área:** Backend / Testing
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-30` (simulador base).
- **Problema:** Lisandro y Tayra no pueden desarrollar `BAC-42` ni `BAC-43` de forma local si el simulador no responde los endpoints de snapshot.
- **Entregable:**
  1. Handlers simulados en `backend/cmd/proxmox-simulador`:
     - `POST /api2/json/nodes/:node/qemu/:vmid/snapshot` y `POST .../lxc/...` (persiste snapshot en memoria y devuelve UPID).
     - `GET /api2/json/nodes/:node/qemu/:vmid/snapshot` (lista ordenada con `name`, `snaptime`, `description`).
     - `POST .../rollback` y `DELETE .../snapshot/:snapname` (devuelven UPID y simulan la demora asíncrona).
- **Criterio de éxito:** Pruebas con `curl` al simulador permiten crear, listar y revertir un snapshot simulado con seguimiento de UPID.

#### `BAC-51` (`BRG-11-BAC`) - Guarda centralizada de concurrencia operativa y exclusión mutua (`409 INSTANCE_BUSY`)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-29` (`activeTask`), `BAC-37` (lock de cuotas).
- **Problema:** Si una instancia entra en snapshot o rollback, los endpoints de ciclo de vida (`start`, `stop`, `reboot` de Etapa 1) y edición de hardware (`PUT /config` de Etapa 2) permitirían mutaciones concurrentes que corrompen el estado.
- **Entregable:**
  1. Middleware/guard de dominio `RequireInstanceIdle`: valida si la instancia tiene una tarea activa en `tareas_asincronas` o un candado en Redis (`lock:instance:{vmid}:operation`).
  2. Si está ocupada, rechaza inmediatamente con `HTTP 409 Conflict` (`INSTANCE_BUSY`).
  3. Integración de la guarda en: `POST /api/instances/:id/status/:action` y `PUT /api/instances/:id/config`.
- **Criterio de éxito:** Una solicitud de encendido o reconfiguración sobre una máquina con tarea activa devuelve `409` estructurado de forma determinística.


#### `BAC-52` (`BRG-12-BAC`) - Generalización del Broker SSE para Eventos Globales de Alerta (`ALERT_TRIGGERED`)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-21C`, `BAC-26`.
- **Problema:** El canal `/api/events` asume que todo evento es `TASK_FINISHED` sobre una VM específica; las alertas del host de `BAC-47` no tienen VMID y deben transmitirse a todos los administradores.
- **Entregable:**
  1. Extensión del modelo `Evento` en Go para admitir `recursoTipo: "HOST" | "VM" | "LXC"`, siendo `recursoId` opcional.
  2. Regla de enrutamiento en Redis Pub/Sub: los eventos con `recursoTipo == "HOST"` se difunden a todos los clientes con rol `ADMIN`.
  3. Actualización de `docs/contrato-eventos.md`.
- **Criterio de éxito:** Un evento de prueba `ALERT_TRIGGERED` a nivel `HOST` es recibido en el navegador por un administrador y descartado silenciosamente para un operador.


#### `BAC-53` (`BRG-13-BAC`) - Particionado Preventivo e Indexación Compuesta de la Tabla `auditoria`
- **Área:** Backend / Base de Datos
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `BAC-18`, `BAC-41`.
- **Problema:** En PostgreSQL, si no existen las particiones para los trimestres de trabajo (Q3/Q4 2026), las inserciones de auditoría fallarán. Además, el streaming de `BAC-45` requiere índices compuestos para evitar escaneos secuenciales.
- **Entregable:**
  1. Script de migración SQL DDL que crea preventivamente las particiones trimestrales: `auditoria_2026_q3`, `auditoria_2026_q4`, `auditoria_2027_q1`.
  2. Índices compuestos para partition pruning: `CREATE INDEX idx_auditoria_query ON auditoria (fecha_hora DESC, usuario_id, accion)`.
- **Criterio de éxito:** `EXPLAIN ANALYZE` sobre una consulta por rango de fechas confirma que PostgreSQL solo lee la partición correspondiente (*partition pruning*) sin recorrer la tabla completa.
