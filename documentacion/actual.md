### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (24/10/2026):** según ClickUp, en este archivo quedan las tareas **pendientes de implementación** de la Etapa 1 y Etapa 2. Lo ya verificado con las suites está en [`terminado.md`](terminado.md) / [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md), y las correcciones de las implementadas con problemas, en [`futuro.md`](futuro.md) / [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md).
>
> **Pasados a `terminado-1.md` (Etapa 1):**
> - `FIX-104B` / `FIX-105` y `FIX-102` (resueltos y con tests en verde el 21/10/2026).
> - `FIX-100` (resuelto y con tests en verde el 22/10/2026: guarda en `formatCount` y `formatShare`).
> - `FIX-91` (pendiente de implementación: eliminar el código muerto `InstanceTester.tsx`; sigue en `futuro-1.md`).
> - `FIX-101`, `FIX-103` y `FIX-104` (cerrados el 22/10/2026: expectativas desactualizadas de `test/front`, no defectos del submódulo; pruebas en verde).
> - `BAC-24C` / `FIX-113` (resuelto y con tests en verde el 24/10/2026: purga atómica de permisos al eliminar instancia; pruebas `TestFIX113` y pasos de integración en verde).
>
> **Pasados a `terminado2.md` (Etapa 2):**
> - `BAC-36`, `FIX-107`, `FIX-109`, `FIX-110`, `FIX-111A`, `FIX-112`, `FIX-114`, `FIX-115A`, `FIX-116`, `FIX-117A`, `FIX-118`, `FIX-119` (resueltos y con tests en verde el 21/10/2026).
> - `FIX-111B` y `FIX-115B` (resueltos y con tests en verde el 22/10/2026: modal de reautenticación y claves hasta 128 chars).
> - `FRN-24` y `FRN-25` (cerradas el 22/10/2026: implementación funcional y pruebas alineadas al contrato real).
> - `INF-10` (pendiente de verificación en Proxmox físico; seguimiento en `FIX-123` dentro de `futuro2.md`).
> - `FIX-117B` (resuelto y con tests en verde el 24/10/2026: adaptación de Dashboard para operador; `fix-117b-dashboard-operador.test.tsx` 4/4 en verde).
> - `FIX-108` (implementada el 24/10/2026 con fallas de validación de entrada; seguimiento en `FIX-124` dentro de `futuro2.md`).
> - `FRN-22`, `FRN-23`, `FRN-26`, `FRN-27` y `FRN-28` (resueltos y con tests en verde el 26/10/2026: frontend `41cad44`; telemetría en tiempo real WS con 12 tests, curvas SVG de recursos, seguimiento de cuota y edición reactiva de vCPU/RAM).

---

### ETAPA 1: Tareas Pendientes de Implementación

*Actualmente no hay tareas pendientes de implementación en la Etapa 1. Todas las tareas fueron completadas y verificadas en [`terminado-1.md`](terminado-1.md).*

---

### ETAPA 2: Tareas Pendientes de Implementación

*Las tareas base de la Etapa 2 fueron completadas y verificadas en [`terminado2.md`](terminado2.md). Restan las tareas puente de integración hacia la Etapa 3:*

#### `FRN-36` (`BRG-10-FRN`) - Modularización de `InstanceDetailPage` para Pestaña de Snapshots y Estados Transitorios
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `FRN-21`, `FRN-28`.
- **Problema:** La página de detalle actual tiene pestañas fijas y no propaga el estado de bloqueo a la cabecera cuando una máquina está en transición asíncrona.
- **Entregable:**
  1. Reestructuración de `InstanceDetailPage.tsx` con soporte para pestañas modulares (dejando listo el slot para "Snapshots").
  2. Soporte en el badge de estado para etiquetas transitorias: `SNAPSHOTTING` (azul), `ROLLING_BACK` (amarillo/naranja).
  3. Deshabilitación reactiva de la botonera de acciones rápidas de la cabecera cuando `instance.activeTask != null`.
- **Criterio de éxito:** Si la máquina reporta `activeTask`, los botones de la cabecera se bloquean con spinner y la pestaña de Snapshots cuenta con su contenedor maquetado.


#### `INF-13` (`BRG-11-INF`) - Verificación y Validación de Storage Compatible con Snapshots en Proxmox Real
- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.5 h
- **Depende de:** `INF-10`.
- **Problema:** Si el Proxmox físico (`100.81.49.19`) utiliza discos `raw` sobre LVM estándar, la API de snapshots fallará en la validación final (`INT-09`).
- **Entregable:**
  1. Verificación del almacenamiento configurado en Proxmox (`local-lvm` con thin-pool, `zfs` o `qcow2`).
  2. Si el storage por defecto no soporta snapshots, configurar un pool LVM-thin o storage de directorio para las instancias de prueba del Centinela.
  3. Documentar en `documentacion/api-proxmox.md` los storages aptos para snapshots.
- **Criterio de éxito:** Ejecución manual de `qm snapshot <vmid> test` por CLI en el Proxmox real completa exitosamente sin errores de formato de disco.

---
# ETAPA 3
---


#### `BAC-42` - Adaptador de Proxmox para Snapshots, endpoints de Creación y Listado
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-50` (`BRG-10-BAC`: simulador con soporte de snapshots).
- **Problema:** Los usuarios autorizados necesitan registrar puntos de restauración para congelar el estado de sus máquinas antes de realizar operaciones de riesgo.
- **Entregable:**
  1. Métodos en el cliente Proxmox: `CrearSnapshot(ctx, node, tipo, vmid, snapname, descripcion)` y `ListarSnapshots(ctx, node, tipo, vmid)`.
  2. Endpoint `GET /api/instances/:vmid/snapshots` con guard de acceso `RequireInstanceAccess`.
  3. Endpoint `POST /api/instances/:vmid/snapshots` validando nombre alfanumérico (sin espacios ni caracteres especiales), captura de UPID e inserción en `tareas_asincronas` con respuesta `HTTP 202 Accepted { upid, tareaId, snapname }`.
  4. Pruebas unitarias en `snapshot_service_test.go` contra el simulador.
- **Criterio de éxito:** Un usuario con permisos crea y lista snapshots recibiendo la jerarquía temporal ordenada sin errores de serialización.

#### `BAC-43` - Reversión (Rollback) y Eliminación de Snapshots con bloqueo de concurrencia
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-42`, `BAC-51` (`BRG-11-BAC`: guarda de concurrencia `409 INSTANCE_BUSY`).
- **Problema:** La reversión o borrado de snapshots altera profundamente el almacenamiento del hipervisor. Si se ejecutan en simultáneo con órdenes de energía o creación, se corre riesgo de corrupción de datos.
- **Entregable:**
  1. Endpoint `POST /api/instances/:vmid/snapshots/:snapname/rollback` (requiere `ADMIN` o `FULL_ACCESS`).
  2. Endpoint `DELETE /api/instances/:vmid/snapshots/:snapname` (requiere `ADMIN` o `FULL_ACCESS`).
  3. Validación de estado: si la instancia tiene una tarea activa (`activeTask != nil`), rechazar inmediatamente con `HTTP 409 Conflict` (`INSTANCE_BUSY`).
  4. Al despachar la reversión, marcar la máquina en estado transitorio `ROLLING_BACK` y emitir `TASK_FINISHED` al concluir el UPID.
- **Criterio de éxito:** Intentos de revertir snapshots mientras otra tarea está en curso responden 409. Al completar el rollback, la máquina regresa al estado anterior y se liberan los bloqueos.


#### `BAC-44` - Interceptor global y consultas complejas sobre auditoría particionada
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-41`, `BAC-53` (`BRG-13-BAC`: particionado preventivo e índices compuestos de auditoría).
- **Problema:** La auditoría debe consolidar eventos de todas las etapas con filtros combinados eficientes (por usuario, rango de fechas, tipo de recurso, acción y severidad) aprovechando las particiones de PostgreSQL.
- **Entregable:**
  1. Extensión del interceptor de auditoría para registrar: `SNAPSHOT_CREATE`, `SNAPSHOT_ROLLBACK`, `SNAPSHOT_DELETE`, `ALERT_TRIGGERED` y `CONFIG_UPDATED`.
  2. Endpoint `GET /api/admin/audit` con parámetros de consulta: `?desde=...&hasta=...&usuarioId=...&accion=...&recursoTipo=...&recursoId=...&page=1&limit=25`.
  3. Optimización de consulta SQL utilizando exclusión de particiones (`partition pruning`) por `fecha_hora`.
- **Criterio de éxito:** Consultas sobre millones de registros históricos responden en menos de 150 ms sin realizar escaneos completos de tablas.

#### `BAC-45` - Streaming de exportación estructurada de auditoría a CSV y JSON
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-44`.
- **Problema:** Exportar reportes de auditoría extensos cargando todos los registros en memoria del servidor puede provocar caídas por `Out of Memory` (OOM).
- **Entregable:**
  1. Endpoint `GET /api/admin/audit/export` accesible solo para usuarios con rol `ADMIN`.
  2. Parámetro `format=csv|json`, respetando exactamente los mismos filtros aplicados en la vista.
  3. Implementación de streaming HTTP: lectura en lotes con cursores de base de datos (`pgx` / `gorm.Rows`) y escritura inmediata en el stream de respuesta.
  4. Cabeceras HTTP apropiadas: `Content-Disposition: attachment; filename="auditoria_centinela_YYYYMMDD.csv"`.
- **Criterio de éxito:** Exportaciones de más de 50.000 filas se descargan fluidamente en el navegador sin picos de consumo de RAM en el contenedor backend.


#### `INF-11` - Configuración de producción en Nginx: TLS 1.3, HSTS, CSP y WSS
- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 2.0 h
- **Depende de:** `INF-04`, `INF-09`.
- **Problema:** El entorno de producción final requiere estándares estrictos de cifrado, mitigación de ataques web y estabilidad total en canales bidireccionales.
- **Entregable:**
  1. Configuración de Nginx en el servidor de producción con protocolos modernos (TLS 1.2 y TLS 1.3 exclusivamente).
  2. Inyección de encabezados de seguridad:
     - `Strict-Transport-Security: max-age=31536000; includeSubDomains`
     - `X-Content-Type-Options: nosniff`
     - `X-Frame-Options: DENY`
     - `Content-Security-Policy` ajustada para WebSockets y fuentes locales.
  3. Verificación de proxying sin caídas de sesión ni timeout en WebSockets `/api/ws/*` y eventos SSE `/api/events`.
- **Criterio de éxito:** Análisis SSL Labs o verificación con curl reporta calificación A/A+ con todos los encabezados de seguridad presentes.

#### `BAC-46` - Background worker evaluador de umbrales con histéresis y deduplicación
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `BAC-22` (telemetría de nodo) y `BAC-34` (telemetría de instancias).
- **Problema:** La saturación prolongada de CPU, RAM o almacenamiento en el host o en las máquinas críticas puede degradar los servicios. Se requiere detección proactiva antes de que ocurra una caída.
- **Entregable:**
  1. Worker periódico de evaluación configurable (`ALERT_POLL_INTERVAL`, por defecto 30 s).
  2. Verificación de umbrales críticos:
     - Host: RAM $\ge 85\%$, CPU $\ge 90\%$, Almacenamiento $\ge 90\%$.
     - Instancias: Memoria asignada saturada o estado inaccesible.
  3. Algoritmo de histéresis en Redis: marcar alerta activa en clave con TTL; solo disparar evento si la condición persiste durante 2 ciclos consecutivos y resolverla cuando caiga por debajo del 80%.
- **Criterio de éxito:** No se emiten alertas espurias por picos momentáneos de 1 segundo. Al sostenerse la saturación, se genera el estado de alarma de forma determinística.

#### `FRN-29` - Pestaña "Snapshots" en Detalle de Instancia y formulario de nuevo snapshot
- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.0 h
- **Depende de:** `BAC-42`, `FRN-36` (`BRG-10-FRN`: detalle modularizado y estados transitorios).
- **Problema:** La vista de detalle no cuenta con una interfaz visual para inspeccionar los puntos de restauración existentes ni disparar un nuevo guardado.
- **Entregable:**
  1. Pestaña "Snapshots" dentro de `InstanceDetailPage.tsx`.
  2. Línea de tiempo o tabla interactiva de snapshots con nombre, descripción, fecha de creación y badge de estado.
  3. Botón "Crear Snapshot" que despliegue un modal sencillo solicitando nombre (validado contra `/^[a-zA-Z0-9_\-]+$/`) y descripción opcional.
  4. Feedback inmediato con spinner de carga durante la creación.
- **Criterio de éxito:** El usuario visualiza la lista de snapshots y puede crear uno nuevo viendo la lista actualizada reactivamente.


#### `FRN-31` - Vista administrativa de Auditoría con filtros reactivos completos
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.5 h
- **Depende de:** `BAC-44`, `FIX-74`.
- **Problema:** Los administradores necesitan auditar la actividad institucional de forma ágil, buscando incidentes o filtrando anomalías operativas.
- **Entregable:**
  1. Rediseño funcional de `Auditoria.tsx`:
     - Rango de fechas con controles accesibles (`aria-label="Desde"` y `aria-label="Hasta"`).
     - Desplegables con autocompletado para tipo de acción, usuario y tipo de recurso.
     - Indicador visual de severidad (Badge verde = `SUCCESS`, rojo = `FAILED`, amarillo = `WARNING`).
  2. Tabla paginada con controles de avance/retroceso y selector de elementos por página.
  3. Modal para inspeccionar el payload JSON completo (`detalles`) de un evento seleccionado.
- **Criterio de éxito:** Los cambios en los filtros actualizan reactivamente la tabla vía query parameters cumpliendo con todas las aserciones de accesibilidad.