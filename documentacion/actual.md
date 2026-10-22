### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (22/10/2026):** según ClickUp, en este archivo quedan las tareas **pendientes de implementación** de la Etapa 1 y Etapa 2. Lo ya verificado con las suites está en [`terminado.md`](terminado.md) / [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md), y las correcciones de las implementadas con problemas, en [`futuro.md`](futuro.md) / [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md).
>
> **Pasados a `terminado-1.md` (Etapa 1):**
> - `FIX-104B` / `FIX-105` y `FIX-102` (resueltos y con tests en verde el 21/10/2026).
> - `FIX-100` (resuelto y con tests en verde el 22/10/2026: guarda en `formatCount` y `formatShare`).
> - `FIX-91` (pendiente de implementación: eliminar el código muerto `InstanceTester.tsx`; sigue en `futuro-1.md`).
> - `FIX-101`, `FIX-103` y `FIX-104` (cerrados el 22/10/2026: expectativas desactualizadas de `test/front`, no defectos del submódulo; pruebas en verde).
>
> **Pasados a `terminado2.md` (Etapa 2):**
> - `BAC-36`, `FIX-107`, `FIX-109`, `FIX-110`, `FIX-111A`, `FIX-112`, `FIX-114`, `FIX-115A`, `FIX-116`, `FIX-117A`, `FIX-118`, `FIX-119` (resueltos y con tests en verde el 21/10/2026).
> - `FIX-111B` y `FIX-115B` (resueltos y con tests en verde el 22/10/2026: modal de reautenticación y claves hasta 128 chars).
> - `FRN-24` y `FRN-25` (cerradas el 22/10/2026: implementación funcional y pruebas alineadas al contrato real).
> - `INF-10` (pendiente de verificación en Proxmox físico; seguimiento en `FIX-123` dentro de `futuro2.md`).

---

### ETAPA 1: Tareas Pendientes de Implementación

#### `BAC-24C` / `FIX-113` - Limpieza de permisos al eliminar una instancia
- **Área:** Backend
- **Asignada:** Tayra / Lisandro
- **Estimación:** 1.0 h
- **Depende de:** `BAC-24B` y `BAC-25B`.
- **Importancia:** 🔴 **ALTA** (CWE-284: Improper Access Control, CWE-639: Authorization Bypass Through User-Controlled Key).
- **Problema:** `/cluster/nextid` devuelve el menor VMID libre. Si se borra la 110 y después se crea otra instancia, recibe el 110 y **hereda las filas de `permisos_instancia`** de la borrada: un operador tendría acceso a una máquina que no se le asignó.
- **Entregable:**
  - Cuando una tarea de borrado termina con `COMPLETED`, eliminar en una transacción las filas de `permisos_instancia` de ese VMID.
  - Registrar la limpieza en `auditoria` con `accion: INSTANCE_PERMISSIONS_PURGED` y la lista de usuarios afectados.
  - Publicar el `TASK_FINISHED` **antes** de borrar los permisos, para que los operadores asignados lo reciban.
  - Si la tarea falla, no se toca nada.
- **Criterio de éxito:** después de borrar la 110, `GET /api/admin/users/:id/permissions` del operador ya no la incluye. Al crear una instancia nueva con el VMID 110, el operador no la ve.

---

### ETAPA 2: Tareas Pendientes de Implementación

### `FIX-108` - Adaptación del Frontend para consumo de identificadores de instancia ofuscados (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1,5 h
- **Depende de:** `FIX-107`.
- **Problema técnico detallado:**
  Actualmente, el frontend modela el identificador de instancia estrictamente como un número entero (`vmid: number` y `id: number` en `InventoryInstance`), validándolo con `Number.isSafeInteger(vmid) && vmid > 0` en `validateInstanceVmid`.
  Al migrar el backend hacia identificadores opacos/cifrados no secuenciales (ej. `"inst_k8Pq2v"`), el frontend rechaza las respuestas con errores de validación de esquema en `instanceService.ts`.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/components/features/instances/types/instance.ts`:
     - Actualizar el tipo `InventoryInstance`:
       ```typescript
       export interface InventoryInstance {
         id: string | number; // Soporta el identificador opaco público
         name: string;
         type: 'VM' | 'LXC';
         ...
       }
       ```
  2. En `instanceService.ts`:
     - Flexibilizar `validateInstanceVmid` o renombrarlo a `validateInstanceId(id: string | number): string` para admitir tanto cadenas alfanuméricas ofuscadas como enteros legados, sanitizando que no contengan caracteres especiales ni inyecciones.
     - Actualizar las URLs de invocación (`/instances/${validId}/${action}`) para enviar el string del identificador público.
  3. En las rutas de la aplicación (`src/routes` o componentes con parámetros de URL `/instances/:id`):
     - Asegurar que los componentes de detalle de instancia y modales de acción reciban y operen con el identificador en formato string sin forzar conversiones numéricas destructivas (`parseInt` o `Number(...)`).
- **Criterio de éxito:**
  - La tabla de instancias renderiza correctamente el listado proveniente del backend con IDs ofuscados.
  - Los botones de acción de ciclo de vida (iniciar, apagar, ver métricas) operan de manera transparente utilizando el identificador opaco en las URLs del navegador.


### `FIX-117B` - Adaptación del Dashboard de operador ante segregación de métricas globales (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1,0 h
- **Depende de:** `FIX-117A`.
- **Importancia:** 🟡 **MEDIA** (Resiliencia UI y prevención de errores por campos omitidos).
- **Problema técnico detallado:**
  Al implementarse `FIX-117A`, el backend deja de enviar métricas globales del host (`instancesSummary`, almacenamiento global) cuando quien consulta es un usuario con rol `OPERATOR`. Los componentes del Dashboard principal que asumen la presencia incondicional de estos objetos arrojan excepciones en tiempo de ejecución o renderizan valores vacíos/NaN.
- **Qué se debe hacer:**
  1. En los componentes de Dashboard y resumen de métricas del frontend:
     - Condicionar la visualización de tarjetas de métricas globales de hardware y almacenamiento según el rol del usuario autenticado (`user.rol === 'ADMIN'`).
     - Para usuarios con rol `OPERATOR`, renderizar un resumen enfocado exclusivamente a las instancias asignadas y accesibles, sin intentar acceder a las propiedades omitidas por el backend.
- **Criterio de éxito:**
  - El Dashboard de un operador renderiza limpiamente sin advertencias ni errores de consola ante la respuesta acotada de `/node/status`.


#### `FRN-22` - Hook `useInstanceMetrics` y cliente WebSocket para telemetría en tiempo real

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** `BAC-35`, `FRN-21` (página de detalle).
- **Problema:** El cliente frontend debe gestionar la conexión WebSocket de métricas de forma resiliente, solicitando el ticket de autenticación previo y manejando caídas de red sin fugar memoria ni sockets huérfanos.
- **Entregable:**
  1. Servicio de solicitud de ticket efímero `getWsTicket()`.
  2. Hook personalizado `useInstanceMetrics(vmid)` que gestione conexión, estado (`conectando`, `conectado`, `reconectando`, `error`) y almacenamiento en búfer de los últimos $N$ puntos de telemetría.
  3. Reintento con backoff exponencial ante desconexión accidental.
  4. Cierre explícito del socket al desmontar el componente o cambiar de pestaña.
  5. **Contrato de telemetría unificado (escala y nulos):**
     - La CPU transmitida en vivo por WebSocket (`MetricasEnVivoDTO.cpu`) y en el histórico RRDdata (`PuntoMetricaDTO.cpuPercent`) viaja **en escala porcentual directa ($0$ a $100$)**, no como fracción $0..1$. El hook debe almacenar la telemetría en dicha escala sin multiplicar por 100.
     - Cuando la instancia esté apagada (`stopped`), el backend envía `cpu: null` y `mem: null`. El hook debe manejar valores nulos de forma segura (sin invocar `.toFixed()` directamente sobre `null`) e indicar en el estado del búfer la inactividad temporal.
- **Criterio de éxito:** El hook entrega series de datos actualizadas cada 2 s sin fugas de memoria en pruebas de componente con Vitest y maneja instancias apagadas (`null`) sin excepciones de tipo en consola.


#### `FRN-23` - Componentes de gráficos interactivos de rendimiento y selector de rango temporal

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.5 h
- **Depende de:** `FRN-22`, `BAC-36`.
- **Problema:** La pestaña "Métricas" de `InstanceDetailPage.tsx` carece de componentes visuales para interpretar el consumo de recursos de la máquina.
- **Entregable:**
  1. Incorporación de componentes de gráficos interactivos (utilizando Recharts o Chart.js):
     - Gráfico de CPU (%) en escala directa $0-100\%$ (compatible con `PuntoMetricaDTO.cpuPercent` y WebSocket).
     - Gráfico de Memoria RAM (Usada vs Total).
     - Gráfico de Tráfico de Red (KB/s In/Out).
     - Gráfico de I/O de Disco (KB/s Read/Write).
  2. Selector de rango temporal: "En vivo" (WebSocket), "Últimas 2 horas" y "Últimas 24 horas" (REST `BAC-36`).
  3. Formateo legible de unidades (bytes a MB/GB, bps a KB/s).
  4. Tooltips explicativos y estados de carga (skeletons) mientras se obtienen los primeros datos.
  5. **Tolerancia a valores nulos:** cuando la máquina esté detenida (`cpu: null`, `mem: null`), los gráficos deben pintar discontinuidades temporales o mostrar el estado "Instancia detenida" sin romper el ciclo de vida del componente con errores `NaN` o `TypeError`.
- **Criterio de éxito:** Al abrir la pestaña "Métricas", se visualiza la serie histórica y los gráficos continúan actualizándose dinámicamente con los eventos en vivo, tolerando paradas de la máquina sin fallos visuales.


---


#### `FRN-26` - Envío de creación, manejo de errores de cuota (`409`) y seguimiento en la UI

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `FRN-25`, `BAC-38`.
- **Problema:** Durante el envío, la interfaz debe bloquearse para evitar dobles clics, manejar amigablemente los rechazos por cuota insuficiente y reflejar el progreso de la creación en la tabla de inventario.
- **Entregable:**
  1. Despacho HTTP `POST /api/instances` con estado de carga en el botón de confirmación.
  2. Manejo de error `409 RESOURCE_QUOTA_EXCEEDED`: banner en rojo dentro del modal indicando el motivo exacto sin perder los datos cargados en los pasos anteriores.
  3. Ante `202 Accepted`: cierre del modal, inserción de la nueva fila en `/instances` en estado *"Creando..."* con spinner activo.
  4. Escucha del evento `TASK_FINISHED` de creación vía `useEventsContext()` para actualizar la fila a *"Stopped"*, emitir Toast de éxito y habilitar sus controles.
- **Criterio de éxito:** Flujo completo de alta probado sin recargar la pantalla (`F5`), pasando de creación a inventario activo de forma fluida.

#### `FRN-27` - Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `BAC-40`, `FRN-21`.
- **Problema:** Los operadores autorizados necesitan ajustar la capacidad de sus máquinas con controles seguros que indiquen claramente si la acción requiere reinicio.
- **Entregable:**
  1. Modal `EditInstanceConfigModal.tsx` accesible desde la barra de acciones o pestaña "Configuración" de `/instances/:vmid`.
  2. Controles precargados con los valores actuales de vCPU y RAM.
  3. Visualización del cálculo diferencial ($\Delta\text{RAM}$) y comparación contra el espacio disponible en el host.
  4. Mensaje informativo: *"Si la instancia está encendida y no admite hot-plug, los cambios se harán efectivos tras reiniciar el sistema operativo huésped."*
- **Criterio de éxito:** El modal bloquea valores que excedan la memoria libre del host y calcula correctamente el delta antes de confirmar.

#### `FRN-28` - Integración de edición en la página de Detalle de Instancia y refresco reactivo

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1.5 h
- **Depende de:** `FRN-27`.
- **Problema:** Tras confirmar una edición, la interfaz debe reflejar los nuevos límites sin requerir una recarga manual del navegador.
- **Entregable:**
  1. Conexión del modal con `PUT /api/instances/:vmid/config`.
  2. Toasts informativos:
     - Verde: *"Recursos actualizados exitosamente en caliente."*
     - Amarillo: *"Configuración guardada. La máquina requiere reinicio para aplicar los nuevos recursos."*
  3. Actualización inmediata de los valores de CPU/RAM en la cabecera del detalle y en la tabla de `/instances`.
- **Criterio de éxito:** La modificación de recursos se refleja instantáneamente en la interfaz y emite el feedback correspondiente al estado de hot-plug.


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