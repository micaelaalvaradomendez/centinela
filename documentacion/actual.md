### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Integración y verificación local (09/10/2026):** tras ejecutar las suites de pruebas automatizadas sobre backend `4f68e44` y frontend `ee4644c`:
> - Se verificaron y promovieron a [`terminado.md`](terminado.md): `FIX-58`, `FIX-59` (con observación en `FIX-80`), `FIX-60`, `FIX-61` y `FIX-70`.
> - Se verificaron y promovieron a [`terminado-1.md`](terminado-1.md): `FRN-15`, `FRN-19B` (con observación en `FIX-78`), `FRN-16` (con observación en `FIX-77`), `FIX-52`, `BAC-25B`, `FRN-20B` (con observación en `FIX-79`), `FIX-48` y `FIX-62`.
> - Las observaciones detectadas en tareas implementadas se canalizan en [`futuro.md`](futuro.md) (`FIX-80`) y [`futuro-1.md`](futuro-1.md) (`FIX-77`, `FIX-78`, `FIX-79`).
> - En este archivo permanecen únicamente las tareas **pendientes de implementación**:
>
> | Tarea | Área | Estado | Pruebas / Alcance |
> |---|---|---|---|
> | `FIX-54` | Backend | No implementado: columna `eliminado_en` y separación de baja lógica | `fixes_acceptance_test.go` / `admin-users-deletion.test.tsx` |
> | `FIX-55` | Backend | No implementado: alineación de baja y autenticación con identidad histórica | `fixes_acceptance_test.go` |
> | `FIX-56` | Frontend | No implementado: UI distingue suspensión de eliminación | `admin-users-deletion.test.tsx` |

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

> **Cobertura incorporada del PR #3:** las pruebas nuevas comprueban el comportamiento actual: conflicto al reutilizar correo o username de un usuario dado de baja, y alta con datos distintos. No implementan ni validan la reutilización del correo requerida por `FIX-54` a `FIX-57`; esa regresión debe cambiar cuando se implemente el nuevo contrato.

> **Nota de alcance (06/10/2026, revisión estática):** “recorrido completo” se refiere a los casos anteriores, no a crear una cuenta nueva con el correo de una eliminada y preservar la identidad histórica. Esta brecha de la fase base queda pendiente en `FIX-54` a `FIX-57` de [`futuro.md`](futuro.md#fix-54---separar-eliminación-lógica-y-suspensión-migrar-unicidad-del-correo-backend). No se ejecutaron suites ni se verificó el esquema desplegado.

---

### `FIX-54` - Separar eliminación lógica y suspensión; migrar unicidad del correo (Backend)

- **Área:** Backend
- **Estado:** Pendiente; no implementado ni probado.
- **Estimación:** 2 h
- **Depende de:** `BAC-05`, `BAC-06` y la clasificación controlada de datos históricos descrita arriba; coordinar disponibilidad de evidencia con `FIX-51`, sin mezclar su solución.
- **Problema:** la baja conserva correctamente la identidad, pero la unicidad incondicional reserva su correo; `activo` no distingue eliminación de suspensión.
- **Entregables:**
  1. Agregar `eliminado_en timestamptz` nullable y definir la invariante de eliminado siempre inactivo, conservando identidad y relaciones. Aplicar el plan seguro de clasificación histórica; no inferir eliminación solo desde `activo=false`.
  2. Quitar `uniqueIndex` incondicional de `EmailUsuario` del modelo e inspeccionar el nombre real del índice/constraint existente antes de reemplazarlo explícitamente por `UNIQUE (lower(btrim(email_usuario))) WHERE eliminado_en IS NULL`. Evitar que `AutoMigrate` recree la unicidad anterior; mantener la del username.
  3. Prevalidar duplicados normalizados entre no eliminados sin fusionar identidades ni modificar auditoría; definir resolución controlada y transacción/orden seguro de migración. Documentar base nueva, base existente, reinicio e idempotencia.
  4. Corregir el comentario `SET NULL` para reflejar `RESTRICT`, sin cambiar FK ni aplicar scopes de borrado a JOIN de auditoría o actividad histórica.
- **Criterios de aceptación pendientes:** múltiples eliminados pueden conservar el mismo correo; solo una fila no eliminada puede reservarlo, incluso con mayúsculas/espacios; suspendidos lo reservan. Migración repetible y reinicio no recrean el índice anterior ni pierden filas/relaciones/auditoría.

### `FIX-55` - Alinear baja, validación y autenticación con identidad histórica (Backend)

- **Área:** Backend
- **Estado:** Pendiente; no implementado ni probado.
- **Estimación:** 2,5 h
- **Depende de:** `FIX-54`, `BAC-06`, `BAC-06B`, `BAC-13`, `BAC-15` y los mecanismos existentes de autenticación/revocación.
- **Problema:** cambiar solo la unicidad dejaría lookup ambiguo, reactivación de eliminados y flujos con estado temporal anterior a la baja.
- **Entregables:**
  1. `DELETE` fija la marca y desactiva conservando identidad; `PUT activo=false` solo suspende y `PUT activo=true` no restaura eliminados. Bloquear edición, cambios de permisos y resets administrativos sobre eliminados; conservar consultas históricas explícitas por UUID.
  2. `ExisteEmailEnOrg` y validaciones de crear/editar/perfil excluyen eliminados, incluyen suspendidos y comparten normalización con DB. El lookup de autenticación por email excluye eliminados, mantiene el rechazo de suspendidos y no usa `First` ambiguo entre generaciones.
  3. Incorporar guardas de eliminado por UUID en login, 2FA, refresh y recuperación y revisar los caminos que emiten/usan acceso. Invalidar OTP de recuperación y expiración, preauth, tickets, sesiones PostgreSQL/Redis y streams con los mecanismos existentes. Definir consistencia/rollback o estrategia fail-closed: un error de revocación no debe producir éxito engañoso ni dejar acceso habilitado.
  4. Crear con UUID y clave nuevos, cambio obligatorio y 2FA nuevo, sin permisos/sesiones heredados; mantener relaciones históricas en el UUID anterior. Conservar la unicidad DB como última barrera ante carreras y mapear su violación de correo a `409` estructurado, con código acordado en contrato, sin SQL crudo y sin convertir cualquier fallo DB en conflicto.
  5. Revisar el envío SMTP previo al INSERT al probar altas concurrentes y documentar su riesgo; no prometer una transacción distribuida email/DB ni incorporar un refactor SMTP ajeno como requisito obligatorio.
- **Criterios de aceptación pendientes:** nueva cuenta accesible sin seleccionar la histórica; eliminado no puede reactivarse, editarse ni acceder mediante credenciales, tokens, preauth u OTP previos. Suspensión sigue siendo reversible y reserva correo. Conflictos y fallos de revocación tienen respuestas controladas y coherentes con el estado persistido.

### `FIX-56` - Distinguir suspensión y eliminación en contrato e interfaz (Frontend)

- **Área:** Frontend
- **Estado:** Pendiente; no implementado ni probado.
- **Estimación:** 1,5 h
- **Depende de:** DTO/contrato de `FIX-54`/`FIX-55`, `FRN-05`, `FRN-06`, `FRN-06B` y `FIX-24`.
- **Problema:** la interfaz llama eliminación a una desactivación y no dispone de estado separado; no es la causa del bloqueo de unicidad backend.
- **Entregables:**
  1. Diferenciar “Inactivo” de “Eliminado” según DTO; deshabilitar reactivar, editar y resetear eliminados. Ocultarlos del listado operativo por defecto manteniendo acceso histórico explícito y auditoría.
  2. Modal de eliminación: explicar conservación de historial y liberación del correo; suspensión: explicar reversibilidad y reserva del correo. Alinear acciones, mensajes y navegación con el contrato.
  3. Mostrar el conflicto de correo de cuentas no eliminadas con el código acordado por backend; mantener alta `201`, baja `204` y restricción `403` salvo cambio contractual explícito documentado.
- **Criterios de aceptación pendientes:** mensajes y acciones distinguen ambos estados; un eliminado no ofrece acciones operativas; historial accesible y errores de correo claros, sin atribuir la restricción DB al frontend.


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


#### `BAC-26` - Verificación de `TASK_FINISHED` por `/api/events` para todas las acciones (`RF-11` parcial)
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

#### `BAC-27` - Auditoría del resultado de las acciones de ciclo de vida en `auditoria` (`RF-08`)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:**
  - `BAC-18` y `FIX-23` (terminadas);
  - `FIX-39` (registra la orden despachada; la del `DELETE`, `BAC-24B`);
  - `BAC-25B`;
  - `BAC-18B` (particionamiento de `auditoria`, que tiene que estar mergeado antes).
- **Entregable:**
  - Para cada tarea que termina, registrar en `auditoria` una segunda entrada con el mismo `upid` y `tareaId` que la de `FIX-39`, más `status: "SUCCESS"` o `"FAILED"` y `exitstatus`.
  - La entrada de la orden despachada (`PENDING`) la hacen `FIX-39` y `BAC-24B`. Acá solo se verifica que tenga `action`, `user_id`, `resource_id` y `upid`.
- **Criterio de éxito:**
  - Cada acción deja exactamente dos registros correlacionados por `upid`, sin migraciones de esquema nuevas.
  - Las tareas vencidas por timeout quedan como `FAILED`.
  - `GET /api/admin/audit` permite filtrar por `accion`.



#### `FRN-19C` (`BRG-05-FRN`, parte Dashboard) - Tarjetas de conteo de VMs/LXC con auto-actualización (`RF-02`)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 1.0 h
- **Depende de:** `BAC-22B`, `FRN-19B` y `FRN-17A`.
- **Entregable:** tarjetas de resumen de VMs y LXC (`En ejecución`, `Detenidas`, `Total`) con los datos de `instancesSummary`. Usan la misma consulta cada 10 s de `FRN-19B` y vuelven a consultar al recibir cualquier `TASK_FINISHED`.
- **Criterio de éxito:** después de encender una VM, el conteo cambia sin `F5` en menos de 2 s desde el evento.



#### `FRN-20C` (`BRG-05-FRN`, parte Inventario) - Columnas de CPU/RAM y estado en vivo en el Inventario (`RF-03`) (nueva, separada de `FRN-19C`)

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 1.5 h
- **Depende de:** `BAC-22B`, `FRN-20A` y `FRN-17A`.
- **Entregable:**
  1. Columnas de CPU (`%`) y RAM (`GB usados / GB totales`), con `-` si vienen en `null`.
  2. Ante cualquier `TASK_FINISHED` recibido, aunque la acción la haya iniciado otro usuario, volver a consultar la fila o la lista para actualizar estado y métricas. Si la acción fue `delete` y terminó `COMPLETED`, quitar la fila.
- **Criterio de éxito:** si otro usuario apaga una VM, la fila cambia a `Stopped` sin `F5`.



#### `FRN-17B` - Desbloqueo reactivo de instancia y notificación Toast adaptativa
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



#### `BAC-30` (`BRG-06-BAC`) - Extensión del simulador de Proxmox para Aprovisionamiento, Storages y RRDdata (`RF-05`, `RF-07`)

- **Área:** Backend / Testing
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-28` y `FIX-35` (en `terminado-1.md`).
- **Problema:** El simulador actual solo responde a acciones de ciclo de vida de energía, inventario, borrado e interfaces de red. Para desarrollar y probar la Etapa 2 en local (siguiendo la premisa de no depender del servidor real), el simulador debe soportar:
  1. `POST /nodes/{node}/qemu` y `POST /nodes/{node}/lxc`: retornar un UPID válido de tipo `qmcreate` / `vzcreate`.
  2. `PUT /nodes/{node}/{tipo}/{vmid}/config`: actualizar núcleos y memoria simulados.
  3. `GET /nodes/{node}/storage` y `GET /nodes/{node}/storage/{storage}/content`: listar storages (`local`, `local-lvm`) y contenidos (templates LXC e imágenes ISO).
  4. `GET /nodes/{node}/{tipo}/{vmid}/rrddata`: retornar array de métricas simuladas (CPU, RAM, red, disco).
- **Entregable:**
  - Nuevos handlers y rutas en `backend/cmd/proxmox-simulador`.
  - Casos de prueba en `backend/cmd/proxmox-simulador/simulador_test.go` verificando respuestas y tipos de datos idénticos a Proxmox VE 9.x.
- **Criterio de éxito:** Las suites de prueba y el desarrollo de la Etapa 2 pueden crear, editar y consultar métricas contra el simulador local sin errores `501 Not Implemented`.

#### `BAC-31` (`BRG-07-BAC`) - Endpoints de inventario de Almacenamiento, Plantillas e ISOs (`RF-07`)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-14`, `BAC-30`.
- **Problema:** El asistente de creación (Paso 3 de `RF-07`) requiere seleccionar en qué storage se creará el disco y qué plantilla LXC o imagen ISO se usará. Actualmente el backend no expone ningún endpoint de storages ni de contenidos del hipervisor.
- **Entregable:**
  1. Método en el cliente Proxmox para consumir `/nodes/{node}/storage` y `/nodes/{node}/storage/{storage}/content`.
  2. Endpoints protegidos con rol `ADMIN`:
     - `GET /api/node/storages`: lista storages habilitados con tipo, espacio total y libre.
     - `GET /api/node/templates`: lista imágenes ISO y templates LXC disponibles normalizados.
  3. Documentación en Swagger e inventario de errores.
- **Criterio de éxito:** Un administrador consulta `GET /api/node/storages` y `GET /api/node/templates` y recibe el listado estructurado de discos e imágenes disponibles. Un operador recibe `403 Forbidden`.

#### `BAC-32` (`BRG-08-BAC`) - Extensión del worker UPID para tareas de creación y autoasignación de permisos (`RF-04`, `RF-07`)

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

#### `INF-09` (`BRG-06-INF`) - Configuración de Nginx para WebSockets y Streaming bidireccional (`RNF-06`)

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.5 h
- **Depende de:** `INF-07B`.
- **Problema:** `INF-07B` configuró Nginx para SSE unidereccional. La telemetría en tiempo real de la Etapa 2 (`RF-05`) utiliza WebSockets (`/api/ws/*` o `/ws/*`), que requieren encabezados `Upgrade` y `Connection: "upgrade"`, además de timeouts extendidos para no interrumpir sesiones activas de monitoreo.
- **Entregable:**
  - Actualizar `docker/nginx-edge.conf` y la plantilla de despliegue del CT 103 con soporte de WebSockets:
    `proxy_http_version 1.1; proxy_set_header Upgrade $http_upgrade; proxy_set_header Connection "upgrade"; proxy_read_timeout 3600s;`.
- **Criterio de éxito:** Una conexión WebSocket establecida a través de Nginx completa el handshake HTTP 101 y se mantiene abierta transmitiendo datos sin cortes por buffer o timeout.

#### `INF-07C` (`BRG-07-INF`) - Ampliación de privilegios del Token de Proxmox para Aprovisionamiento (`RF-07`, `RF-10`)

- **Área:** Infraestructura
- **Asignado:** Nico / Lucas
- **Estimación:** 1.0 h
- **Depende de:** `INF-07`.
- **Problema:** El API Token actual (`centi-api@pve!backend-token`) solo posee `VM.PowerMgmt`, `Sys.Audit` y lectura. Si intenta crear instancias o reconfigurar cores/RAM en la Etapa 2, Proxmox rechazará las peticiones con `403 Permission check failed`.
- **Entregable:**
  - Incorporar los permisos `VM.Allocate`, `Datastore.AllocateSpace`, `VM.Config.CPU`, `VM.Config.Memory`, `VM.Config.Disk`, `VM.Config.Network` al rol del token en Proxmox VE.
  - Verificar permisos con `pveum user permissions` o llamada de prueba con `curl`.
  - Actualizar la documentación en `documentacion/api-proxmox.md`.
- **Criterio de éxito:** El token permite consultar storages, crear contenedores de prueba y modificar configuraciones sin errores de autorización.

---


### `FIX-71` - Robustecer reconexión y reintentos ante arranque en frío de Redis (`INF-06` / `BAC-17A`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro / Lucas
- **Estimación:** 1 h
- **Depende de:** `BAC-17A` e `INF-06` (en `terminado.md`).
- **Problema técnico detallado:**
  Al ejecutar las pruebas unitarias del paquete `backend/internal/adapters/secondary/redis` (`go test ./internal/adapters/secondary/redis/...`), la prueba `TestConectar/se_recupera_si_Redis_aparece_durante_los_reintentos` falla sistemáticamente:
  ```
  kv_store_test.go:79: Debe conectar cuando Redis aparece en un reintento: redis no responde en 127.0.0.1:46429: dial tcp 127.0.0.1:46429: connect: connection refused
  ```
  En `backend/internal/adapters/secondary/redis/kv_store.go`, la función `Conectar` utiliza `store := Nuevo(cfg)` y ejecuta un bucle de 3 intentos de `Ping` espaciados con `esperaArranque = 500ms`. Sin embargo, el cliente `goredis.NewClient` inicializa internamente un pool de conexiones con `DialerRetries: 2` y `MaxRetries: 3`. Cuando el socket no está disponible de inmediato (arranque en frío o Redis levantando de forma asíncrona unos cientos de milisegundos más tarde), el pool retiene errores de conexión o agota sus reintentos internos sin refrescar el estado del socket entre ciclos de sondeo, o el contexto de la prueba simula el reinicio de forma que el cliente no logra negociar la reconexión limpia en el tiempo estipulado.
- **Qué debe hacer el equipo de backend:**
  1. En `kv_store.go`, asegurar que la estrategia de arranque en frío y reintentos de `Conectar` renueve el intento de conexión o limpie las conexiones muertas del pool cuando `Ping` retorne `connection refused`.
  2. Ajustar la sincronización de reintento en `kv_store_test.go` para garantizar que la reactivación del servidor de prueba en segundo plano (`srv.StartAddr` / `srv.Restart`) sea tolerada y verificada de manera determinística por la función `Conectar`.
  3. Ejecutar localmente `go test -v ./internal/adapters/secondary/redis/...` en el backend y certificar que `TestConectar` pase 100% en verde.
- **Criterio de éxito:**
  - `go test ./internal/adapters/secondary/redis/...` pasa en verde sin fallos de conexión.
  - El backend tolera arrancar cuando Redis tarda hasta 2-3 segundos en estar listo para aceptar conexiones TCP.


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


### `FIX-80` - Alineación textual y normalización de copy en vistas de autenticación (`FIX-59`) (Frontend)

- **Área:** Frontend / UI
- **Asignado:** Belinda / Nico
- **Estimación:** 0,5 h
- **Depende de:** `FIX-59`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/components/layout_auth/MainLayoutAuth.tsx`, el placeholder informal `"imagenes y informacion random"` fue eliminado con éxito y se incorporaron explicaciones sobre la infraestructura Proxmox VE y el valor de seguridad de 2FA. Sin embargo, la suite de pruebas `test/front/auth-branding.test.tsx` verifica la existencia de encabezados y títulos estructurados específicos acordados en las especificaciones de diseño:
  - *"Tu infraestructura virtual, simplificada y bajo control"* (en el headline principal).
  - *"Protección de infraestructura con doble factor"* (en la tarjeta de seguridad 2FA).
  Adicionalmente, en la suite de pruebas de testing library, el uso de expresiones regulares amplias sobre selectores singulares provoca fallos por coincidencia múltiple de texto si no se estructuran jerárquicamente.
- **Qué debe hacer el equipo de frontend:**
  1. En `MainLayoutAuth.tsx`, alinear el titular `HEADLINE` y el título del beneficio de 2FA en el array de `features` con las frases literales aprobadas en la especificación.
  2. Mantener la estética responsiva, los efectos de animación y la legibilidad en pantallas reducidas.
  3. Ejecutar `npx vitest run auth-branding.test.tsx` y comprobar que la suite pase al 100% en verde.
- **Criterio de éxito:**
  - El layout de autenticación contiene los textos institucionales precisos aprobados para Centinela y 2FA.
  - La suite `auth-branding.test.tsx` pasa todos sus casos (incluido `FIX-59`).


### `FIX-72` - Sincronizar firma de `NewSeguimientoTareas` en pruebas unitarias de servicios (`BAC-25A`) (Backend / pruebas)

- **Área:** Backend / pruebas
- **Asignado:** Lisandro / Lucas
- **Estimación:** 1 h
- **Depende de:** `BAC-25A` (en `terminado-1.md`).
- **Problema técnico detallado:**
  Al ejecutar las pruebas unitarias del backend invocadas por la prueba de aceptación `TestEtapa1/BAC-25A` (`etapa1_acceptance_test.go:185`), la suite `go test ./internal/core/services/...` falla en compilación:
  ```
  internal/core/services/seguimiento_pool_test.go:89:60: too many arguments in call to services.NewSeguimientoTareas
  	have (*proxmoxPool, *tareasEnMemoria, *publicadorFalso, *auditoriaGrabadora, services.ConfigSeguimiento)
  	want (ports.ProxmoxPort, ports.TareaRepository, services.ConfigSeguimiento)
  internal/core/services/seguimiento_tareas_test.go:131:54: too many arguments in call to services.NewSeguimientoTareas
  	have (*proxmoxTareas, *tareasEnMemoria, *publicadorFalso, *auditoriaGrabadora, services.ConfigSeguimiento)
  	want (ports.ProxmoxPort, ports.TareaRepository, services.ConfigSeguimiento)
  ```
  La firma del constructor `services.NewSeguimientoTareas` fue simplificada para recibir únicamente `(ports.ProxmoxPort, ports.TareaRepository, services.ConfigSeguimiento)`, pero los archivos de pruebas unitarias de servicios (`seguimiento_pool_test.go` y `seguimiento_tareas_test.go`) no fueron actualizados y siguen pasando 5 argumentos con stubs residuales.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/core/services/seguimiento_pool_test.go`, actualizar la llamada a `services.NewSeguimientoTareas` para pasar únicamente los 3 parámetros requeridos.
  2. En `backend/internal/core/services/seguimiento_tareas_test.go`, actualizar la llamada a `services.NewSeguimientoTareas` a la firma de 3 parámetros.
  3. Ejecutar `go test -v ./internal/core/services/...` en el repositorio backend y verificar que compilen y pasen todas las pruebas unitarias.
- **Criterio de éxito:**
  - `go test ./internal/core/services/...` compila limpiamente y pasa en verde.
  - La prueba de aceptación `TestEtapa1/BAC-25A` en `test/back/etapa1_acceptance_test.go` pasa 100% en verde.


### `FIX-76` - Desacoplar guarda de `nivelAcceso` para rol `ADMIN` en modales de energía (`FRN-15`) (Frontend)

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


### `FIX-77` - Mapeo de errores HTTP de ciclo de vida (`BAC-29`) y resiliencia de la máquina de estados (`FRN-16`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Belinda
- **Estimación:** 1,5 h
- **Depende de:** `FRN-16`, `BAC-29`.
- **Problema técnico detallado:**
  Al despachar una orden de ciclo de vida sobre una instancia (`requestInstancePowerAction`), ante un rechazo del servidor la interfaz no mapea los códigos HTTP y claves de error acordados en el contrato D2 (`BAC-29`) a sus mensajes descriptivos específicos, y `ConfirmUserAction` suprime o desvía los errores 403. Al ejecutar `npx vitest run instances-operation.test.tsx`, fallan las pruebas por no encontrar los mensajes esperados:
  - `409 INSTANCE_INVALID_STATE`: debe mostrar *"La instancia ya está encendida"* o *"El estado actual no corresponde a la acción"*.
  - `409 INSTANCE_BUSY`: debe mostrar *"La instancia se encuentra ejecutando otra tarea. Aguarde a que finalice."*.
  - `403 INSTANCE_PROTECTED`: debe indicar que la instancia es de infraestructura protegida.
  - `403 INSTANCE_ACCESS_DENIED`: debe advertir que el operador no posee permisos sobre esa máquina.
  - `502 PROXMOX_UNAVAILABLE`: debe informar que el hipervisor Proxmox no está disponible.
  - `504 PROXMOX_TIMEOUT`: debe avisar que la respuesta de Proxmox agotó el tiempo de espera.
  Además, al ocurrir un error, todos los controles de la fila deben desbloquearse inmediatamente saliendo del estado de spinner.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/components/features/instances/components/InstanceAction.tsx` (o servicio asociado), implementar una función de traducción que reciba el `ApiRequestError` (evaluando `status` y `errorCode`) y genere el mensaje exacto exigido por el contrato D2.
  2. Asegurar que ante cualquier fallo se notifique el mensaje al usuario (vía toast o texto en modal), se resetee `pendingAction` y la fila recupere su interactividad normal.
  3. Verificar que `npx vitest run instances-operation.test.tsx` pase en verde para todos los códigos de error documentados.
- **Criterio de éxito:**
  - Cada uno de los seis códigos de error despliega su mensaje diferenciado y la fila se desbloquea tras el fallo.
  - La suite `instances-operation.test.tsx` pasa 9 de 9 pruebas en verde.


### `FIX-78` - Desduplicación de indicadores de salud del nodo y accesibilidad en Dashboard (`FRN-19B`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Luz
- **Estimación:** 1 h
- **Depende de:** `FRN-19B`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/pages/Dashboard.tsx`, el estado de salud del nodo (`nodeHealth.label`: "Saludable", "Advertencia", "Inaccesible") se renderiza simultáneamente en dos componentes:
  1. En la cabecera de la tarjeta de recursos (`Dashboard.tsx:136`): `<Badge variant="outline">{nodeHealth.label}</Badge>`.
  2. En el pie de página de estado del sistema (`Dashboard.tsx:211`): `<p className="text-caption">{nodeHealth.label}</p>`.
  Esta duplicación genera colisiones al evaluar la interfaz con lectores de pantalla o Testing Library (`findByText(/saludable/i)` y `findByText(/inaccesible/i)` fallan con `Found multiple elements with the text of: ...`). Asimismo, cuando el nodo está cargando o en error 502/504, los skeletons y las alertas de accesibilidad deben identificarse unívocamente sin interferir con otras tarjetas.
- **Qué debe hacer el equipo de frontend:**
  1. En `Dashboard.tsx`, diferenciar semánticamente los textos o asignar atributos de accesibilidad inequívocos (`aria-label="Estado general del nodo: Saludable"` o dejar el semáforo principal en la tarjeta de recursos y usar una descripción extendida o etiqueta diferenciada en el pie).
  2. Asegurar que las respuestas 502 y 504 reflejen inequívocamente la etiqueta `Inaccesible` en el semáforo de recursos.
  3. Ejecutar `npx vitest run dashboard-node.test.tsx` y comprobar que todos los casos pasen en verde.
- **Criterio de éxito:**
  - `dashboard-node.test.tsx` pasa al 100% sin colisiones de texto duplicado.
  - La experiencia visual y de accesibilidad es nítida tanto para usuarios videntes como para lectores de pantalla.


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

