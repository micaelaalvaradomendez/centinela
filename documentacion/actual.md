### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (13/10/2026):** según ClickUp, en este archivo quedan las tareas **pendientes de resolución** de la Etapa 1 y los `FIX` abiertos. Lo ya verificado con las suites está en [`terminado.md`](terminado.md) / [`terminado-1.md`](terminado-1.md), y las correcciones de las implementadas con problemas, en [`futuro.md`](futuro.md) / [`futuro-1.md`](futuro-1.md).
>
> **Pasados a `terminado-1.md` el 13/10/2026 (ya resueltos, sin trabajo pendiente):**
> - `FIX-73` — sincronización de tareas consecutivas en `test/back` (ya aplicada en este repo; `TestPuenteEtapa1` en verde).
> - `FIX-78` — desduplicación del semáforo del nodo (`dashboard-node.test.tsx` 12/12).

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

> **Cobertura incorporada del PR #3:** las pruebas nuevas comprueban el comportamiento actual: conflicto al reutilizar correo o username de un usuario dado de baja, y alta con datos distintos. No implementan ni validan la reutilización del correo requerida por `FIX-54` a `FIX-57`; esa regresión debe cambiar cuando se implemente el nuevo contrato.

> **Nota de alcance (06/10/2026, revisión estática):** “recorrido completo” se refiere a los casos anteriores, no a crear una cuenta nueva con el correo de una eliminada y preservar la identidad histórica. Esta brecha de la fase base queda pendiente en `FIX-54` a `FIX-57` de [`futuro.md`](futuro.md#fix-54---separar-eliminación-lógica-y-suspensión-migrar-unicidad-del-correo-backend). No se ejecutaron suites ni se verificó el esquema desplegado.

---


#### `BAC-26` - Verificación de `TASK_FINISHED` por `/api/events` para todas las acciones (`RF-11` parcial)

> [!NOTE]
> **Estado verificado (13/10/2026, backend `6afc017`): implementada y con test en verde.** `TestEtapa1Ola3BAC26` pasa: TASK_FINISHED por `/api/events` para las acciones de energía con el payload de BAC-29 (`tareaId`, `accion`, `estado`, `exitstatus`) y filtrado por permiso.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.0 h
- **Depende de:** `BAC-25B` y `BAC-21C` (terminada: canal SSE, Pub/Sub y filtro por permiso).
- **Contexto:** el canal ya existe; esta tarea **no crea uno nuevo**. Verifica y documenta la emisión para las cinco acciones.
- **Entregable:**
  1. Pruebas en `eventos_service_test.go` (o de integración) para `start`, `shutdown`, `stop`, `reboot` y `delete`, con el payload de `BAC-29`. Ejemplo:
     ```json
     {
       "id": "0192…",
       "tipo": "TASK_FINISHED",
       "severidad": "INFO",
       "recursoTipo": "VM",
       "recursoId": "101",
       "mensaje": "La tarea de encendido finalizó correctamente",
       "fechaHora": "2026-10-01T13:30:00Z",
       "detalles": { "tareaId": "0192…", "accion": "start", "estado": "COMPLETED", "exitstatus": "OK" }
     }
     ```
  2. Actualizar `docs/contrato-eventos.md` y `docs/eventos-tiempo-real.md`.
- **Criterio de éxito:** los clientes conectados reciben el evento al terminar cada tarea. Un OPERATOR sin permiso sobre la instancia no lo recibe; uno con `READ_ONLY` sí.




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





#### `BAC-32` (`BRG-08-BAC`) - Extensión del worker UPID para tareas de creación y autoasignación de permisos (`RF-04`, `RF-07`)

> [!WARNING]
> **Estado verificado (14/10/2026, backend `a363c9a`): implementada parcialmente, con problema pendiente.** Falta acción `create` en el worker UPID y la autoasignación en `permisos_instancia` (`TestEtapa1Ola3BAC32`). Corrección en **`FIX-93`** ([futuro-1.md](futuro-1.md)).

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-25B`, `BAC-24C`.
- **Problema:** En la Etapa 1, el worker de seguimiento de UPID solo maneja acciones de energía y borrado. Cuando se crea una instancia (`qmcreate`/`vzcreate`), la máquina aún no existe en el inventario durante la tarea. Al finalizar la creación con éxito:
  1. El creador debe quedar automáticamente registrado con `FULL_ACCESS` en `permisos_instancia` para ese nuevo `vmid` (`RF-07`).
  2. El evento `TASK_FINISHED` debe emitirse con `accion: "create"` y `recursoId: vmid`.
- **Entregable:**
  - Soporte de acciones `create` en `seguimiento_tareas.go`.
  - Integración con el repositorio de permisos: al registrarse `COMPLETED` para una tarea `create`, insertar atómicamente `(usuario_id, vmid, FULL_ACCESS)` en `permisos_instancia`.
- **Criterio de éxito:** Al concluir una tarea de creación, el usuario creador ve inmediatamente la instancia asignada con permisos operativos y recibe el evento `TASK_FINISHED` correspondiente.

#### `BAC-33` (`BRG-09-BAC`) - Cálculo y exposición de cuotas asignables del Host (`RF-02`, `RF-07`)

> [!WARNING]
> **Estado verificado (14/10/2026, backend `a363c9a`): implementada parcialmente, con problema pendiente.** `GET /api/node/status` no devuelve el objeto `allocatable` (`TestEtapa1Ola3BAC33`). Corrección en **`FIX-94`** ([futuro-1.md](futuro-1.md)).

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-22`, `BAC-29`.
- **Problema:** `GET /api/node/status` informa el consumo físico actual del nodo, pero el asistente de creación necesita saber exactamente cuánta memoria RAM y disco pueden asignarse con seguridad sin comprometer la estabilidad del host físico (margen de seguridad / headroom).
- **Entregable:**
  - Exponer en `GET /api/node/status` (o endpoint auxiliar `GET /api/node/quota`) el objeto `allocatable`:
    `{ maxCpuCores: number, freeRamMb: number, freeStorageGb: number, headroomPercent: number }`, reservando un 15% de RAM física para el sistema operativo base de Proxmox.
- **Criterio de éxito:** El backend devuelve las cotas máximas seguras para nuevas asignaciones, sirviendo de base al validador de cuotas de la Etapa 2.

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



---


### `FIX-74` - Reactividad de filtros y etiquetas accesibles en vista de Auditoría (`FRN-14` / `FIX-36`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Luz
- **Estimación:** 1,5 h
- **Depende de:** `FRN-14`, `FIX-36` (en `terminado.md`).
- **Problema técnico detallado:**
  Al ejecutar las pruebas unitarias y de integración de la vista de Auditoría (`frontend/centinela/src/test/audit.test.tsx` y `test/front/audit.test.tsx`), fallan los casos:
  - `AuditoriaPage > Reactividad de Filtros: debe invocar al endpoint con los query params al cambiar los filtros`
  - `FRN-14 - Suite de pruebas para vista de Auditoría > actualiza los query parameters de filtro al cambiar los selectores reactivos`
  El error lanzado por Testing Library es:
  ```
  TestingLibraryElementError: Unable to find a label with the text of: Desde
  ```
  En `Auditoria.tsx` (y sus controles subordinados), los inputs de filtrado temporal ("Desde", "Hasta") y el selector de "Acción" no exponen un `aria-label` o `<label>` asociado que coincida con el contrato de accesibilidad requerido por las pruebas. Además, al modificar interactivamente estos selectores, la vista no dispara la actualización reactiva ni traslada los query parameters (`?desde=...&hasta=...&accion=...`) a la petición `GET /api/admin/audit`.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/pages/Auditoria.tsx`, asignar `aria-label="Desde"` y `aria-label="Hasta"` a los campos de rango de fechas, y `aria-label="Acción"` al selector de tipo de acción.
  2. Implementar la sincronización reactiva del estado para que ante cualquier cambio en estos filtros se invoque a `auditService.getAuditLogs` con los parámetros URL correspondientes.
  3. Ejecutar `npm test src/test/audit.test.tsx` en `frontend/centinela` y certificar que todos los casos pasen en verde.
- **Criterio de éxito:**
  - `npm test src/test/audit.test.tsx` en el frontend pasa al 100% sin errores de elementos no encontrados.
  - La vista de auditoría filtra reactivamente en función de los criterios seleccionados y cumple con accesibilidad.




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



#### `FRN-20C` (`BRG-05-FRN`, parte Inventario) - Columnas de CPU/RAM y estado en vivo en el Inventario (`RF-03`) (nueva, separada de `FRN-19C`)

> [!NOTE]
> **Estado verificado (13/10/2026, frontend `7fcc0bd`): implementada y con test en verde.** `instances-metrics.test.tsx` (3/3) pasa: columnas de CPU (%) y RAM (usado / total) con `—` en `null`; la re-consulta ante `TASK_FINISHED` la aporta `useInstances` (`FRN-17A/B`).


- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 1.5 h
- **Depende de:** `BAC-22B`, `FRN-20A` y `FRN-17A`.
- **Entregable:**
  1. Columnas de CPU (`%`) y RAM (`GB usados / GB totales`), con `-` si vienen en `null`.
  2. Ante cualquier `TASK_FINISHED` recibido, aunque la acción la haya iniciado otro usuario, volver a consultar la fila o la lista para actualizar estado y métricas. Si la acción fue `delete` y terminó `COMPLETED`, quitar la fila.
- **Criterio de éxito:** si otro usuario apaga una VM, la fila cambia a `Stopped` sin `F5`.



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


### `FIX-87` - Quitar los datos hardcodeados de la ficha de usuario (`FRN-05` / `FRN-06B` / `FRN-18` / `FIX-25`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Belinda
- **Estimación:** 1 h
- **Depende de:** `FRN-05`, `FRN-06B`, `FRN-18`, `FIX-25` (en `terminado.md`).
- **Problema técnico detallado:**
  En `frontend/centinela/src/pages/detailsUserPage.tsx` hay dos arrays **hardcodeados** que se renderizan como si fueran datos reales del usuario:
  ```ts
  const additionalSecurityInformation = [
    { label: 'Último cambio de contraseña', value: '02/05/2024, 09:10', ... },
    { label: 'Intentos de inicio fallidos', value: '0', ... },
    { label: 'Bloqueado', value: 'No', ... },
  ];
  const recentActivity = [
    { title: 'Inició sesión correctamente', time: 'Hoy, 17:41', address: '192.168.1.77', ... },
    { title: 'Inició instancia Windows 11 (102)', ... },
    { title: 'Tomó snapshot de Ubuntu Server (101)', ... },
  ];
  ```
  Son **placeholders de diseño**: muestran una fecha de cambio de contraseña inventada y una "actividad reciente" ficticia (instancias 101/102 y una snapshot que el sistema todavía no toma). Quedaron en una vista declarada terminada.
- **Qué debe hacer el equipo de frontend:**
  1. Alimentar "Actividad reciente" desde datos reales (por ejemplo `GET /api/admin/audit` filtrado por el usuario) o **quitarla** hasta que exista el endpoint.
  2. Quitar/blanquear `additionalSecurityInformation` hasta que el perfil exponga esos campos (último cambio de contraseña, intentos fallidos, bloqueo).
- **Criterio de éxito:** la ficha de usuario no muestra valores inventados ni textos de placeholder.

### `FIX-89` - Ruta `/instances/new` con asistente estático sin lógica (`FIX-60`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz / Belinda
- **Estimación:** 0,5 h
- **Depende de:** `FIX-60` (en `terminado.md`).
- **Problema técnico detallado:**
  El botón "Crear instancia" del inventario (`FIX-60`) navega a `/instances/new`, pero `frontend/centinela/src/pages/createNewInstances.tsx` es un **placeholder**: declara `steps = ['General','Sistema','Disco','Red','SO','Resumen']` y un `Stepper`, sin ninguna llamada al backend ni servicio (no importa `apiClient` ni `*Service`). Un ADMIN que llega desde un botón real termina en una pantalla sin funcionalidad.
- **Qué debe hacer el equipo de frontend:**
  1. Ocultar/deshabilitar la ruta y el botón hasta que exista el asistente real (Etapa 2), o
  2. Etiquetar la vista como "Próximamente" de forma explícita.
- **Criterio de éxito:** ningún botón lleva a una pantalla sin funcionalidad sin avisar.



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

### `FIX-84` - Completar storages y contenido en el simulador de Proxmox (`BAC-30`) (Backend / Testing)

- **Área:** Backend / Testing
- **Asignado:** Lisandro
- **Estimación:** 1 h
- **Depende de:** `BAC-30` (en `terminado-1.md`).
- **Problema técnico detallado:**
  El simulador ya responde `POST /nodes/{node}/qemu|lxc` (UPID `qmcreate`/`vzcreate`), `PUT /nodes/{node}/{tipo}/{vmid}/config` y `GET /nodes/{node}/{tipo}/{vmid}/rrddata`, pero **no** implementa `GET /nodes/{node}/storage` ni `GET /nodes/{node}/storage/{storage}/content`, que el asistente de creación (`RF-07`) necesita. Hoy responden `501 Not Implemented`.
- **Qué debe hacer el equipo de backend:**
  1. Agregar al simulador `GET /nodes/{node}/storage` (storages `local` y `local-lvm` con tipo, total y libre) y `GET /nodes/{node}/storage/{storage}/content` (templates LXC e ISOs), con el mismo formato que Proxmox VE 9.x.
  2. Dejar `go test ./cmd/proxmox-simulador/...` en verde y la suite `test/back` (BAC-30) sin bloques rojos.
- **Criterio de éxito:** el desarrollo de la Etapa 2 puede consultar storages y contenido contra el simulador sin `501 Not Implemented`.

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


---


### `FIX-81` - Corregir el orden de migraciones: `migrarUnicidadEmail` antes de crear `auditoria` (`FIX-54`) (Backend)

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


#### `BAC-34` - Adaptador de Proxmox RRDdata en el cliente y servicio de telemetría de instancia

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-30` (simulador con soporte RRD).
- **Problema:** Proxmox almacena las métricas históricas de rendimiento en formato RRD consumibles vía `/nodes/{node}/{qemu|lxc}/{vmid}/rrddata`. El cliente backend carece del adaptador para consultar este endpoint y normalizar las métricas heterogéneas entre VMs y LXCs.
- **Entregable:**
  1. Método en el cliente Proxmox `ObtenerMetricasRRD(ctx, nodo, tipo, vmid, timeframe)`.
  2. DTO normalizado `PuntoMetricaDTO`: `{ timestamp: int64, cpuPercent: float64, ramUsedBytes: int64, ramTotalBytes: int64, netInBps: int64, netOutBps: int64, diskReadBps: int64, diskWriteBps: int64 }`.
  3. Pruebas unitarias con mocks en `internal/adapters/secondary/proxmox/client_test.go`.
- **Criterio de éxito:** Consulta exitosa de métricas para VM y LXC devolviendo series temporales con valores numéricos coherentes sin fallos de conversión.


#### `BAC-35` - Canal WebSocket de métricas en tiempo real con ticket efímero

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.5 h
- **Depende de:** `BAC-34`, `BAC-21C` (patrón de tickets en Redis).
- **Problema:** Enviar peticiones REST cada pocos segundos para monitorear métricas degrada el rendimiento. Se requiere un canal bidireccional eficiente por WebSocket `/api/ws/instances/:vmid/metrics` que no exponga credenciales JWT en los query parameters.
- **Entregable:**
  1. Endpoint de emisión de ticket efímero `POST /api/auth/ws-ticket` (TTL 30 s en Redis).
  2. Handler WebSocket `/api/ws/instances/:vmid/metrics?ticket=...` que valide el ticket, verifique acceso mediante `RequireInstanceAccess` y establezca la conexión.
  3. Poller desacoplado por conexión que sondee el estado actual de la instancia cada 2 segundos y envíe el paquete JSON al cliente.
  4. Cierre ordenado y limpieza de goroutines al desconectar el cliente.
- **Criterio de éxito:** Un cliente autenticado con ticket recibe métricas periódicas por WebSocket. Un ticket inválido o expirado responde `401 Unauthorized` antes del handshake.

#### `BAC-37` - Servicio de validación de cuotas del host y bloqueo distribuido

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-33` (allocatable del host).
- **Problema:** Permitir la creación de máquinas sin validar los límites físicos del hipervisor puede provocar el colapso del host (`Out Of Memory` o disco lleno). Además, peticiones concurrentes simultáneas pueden sobreasignar la misma memoria libre si no existe un bloqueo transaccional.
- **Entregable:**
  1. Servicio de dominio `CuotasService.ValidarDisponibilidad(cores, ramMb, discoGb, storageId)`.
  2. Verificación de seguridad: $\text{RAM\_Solicitada} \le \text{RAM\_Libre} - \text{Headroom}(15\%)$ y $\text{Disco\_Solicitado} \le \text{Disco\_Libre}(storageId)$.
  3. Bloqueo distribuido en Redis `lock:instance:provisioning` mediante `SETNX` (TTL 30 s) para garantizar atomicidad durante la evaluación y despacho.
  4. Mapeo a error `409 Conflict` con código estructurado `RESOURCE_QUOTA_EXCEEDED` ante saturación.
- **Criterio de éxito:** Peticiones que superan los recursos seguros son rechazadas inmediatamente sin tocar Proxmox. Peticiones simultáneas compitiendo por el remanente resuelven una con éxito y otra con 409 sin condición de carrera.


#### `BAC-38` - Endpoint `POST /api/instances` para creación de VMs y Contenedores LXC

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `BAC-30`, `BAC-37`.
- **Problema:** No existe el endpoint de aprovisionamiento de recursos. Se requiere un endpoint unificado que oculte las diferencias nativas de Proxmox entre la creación de máquinas virtuales (`/nodes/{node}/qemu`) y contenedores (`/nodes/{node}/lxc`).
- **Entregable:**
  1. Endpoint `POST /api/instances` con rol obligatorio `ADMIN`.
  2. Payload validado:
     ```json
     {
       "tipo": "vm", // "vm" o "lxc"
       "nombre": "srv-prod-01",
       "cores": 2,
       "ramMb": 2048,
       "discoGb": 20,
       "storageId": "local-lvm",
       "templateOso": "debian-12-standard.tar.zst",
       "red": { "modo": "dhcp" } // o "static" con ip y gateway
     }
     ```
  3. Asignación automática del próximo identificador libre mediante `/cluster/nextid`.
  4. Envío de la orden a Proxmox, captura del UPID de creación (`qmcreate`/`vzcreate`), inserción en `tareas_asincronas` con estado `RUNNING` y respuesta `HTTP 202 Accepted { upid, tareaId, vmid }`.
- **Criterio de éxito:** Una solicitud válida responde `202 Accepted` en menos de 1 s, Proxmox comienza el aprovisionamiento y el UPID queda registrado para seguimiento.


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