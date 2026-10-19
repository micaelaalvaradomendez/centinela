# Etapa 2: tareas terminadas

Aca se listan todas las tareas terminadas e implementadas en los submódulos correspondientes a la Etapa 2.

---

## Verificación del 15/10/2026

Verificadas contra backend `dbf7075`.


#### `BAC-34` - Adaptador de Proxmox RRDdata en el cliente y servicio de telemetría de instancia

> [!NOTE]
> **Estado verificado (15/10/2026, backend `dbf7075`): implementada y con test en verde.** Pasan `TestClient_ObtenerMetricasRRD` (casos VM y LXC para timeframes `hour` y `day`) y la suite de integración de Etapa 2 `TestEtapa2BAC34_AdaptadorMetricasRRD`. Devuelve la estructura `PuntoMetricaDTO` normalizada con valores numéricos coherentes.

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

> [!NOTE]
> **Estado verificado (15/10/2026, backend `dbf7075`): implementada y con test en verde.** Pasan `TestAuth_EmitirWsTicket` y `TestMetricsWS_HandshakeYStreaming` (5 subtests: emisión de ticket en Redis con TTL 30 s, handshake exitoso, streaming periódico, rechazo 401 con ticket inválido/expirado y cierre limpio).

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

> [!NOTE]
> **Estado verificado (15/10/2026, backend `dbf7075`): implementada y con test en verde.** Pasan las pruebas de aprovisionamiento (`TestAprovisionamiento_ValidarCuotas_Suficientes`, `TestAprovisionamiento_ValidarCuotas_InsuficientesExcedeRAM`, `TestAprovisionamiento_BloqueoDistribuidoRedis`, `TestCrearInstancia_RecursosInsuficientesEs409`): bloqueo distribuido Redis `lock:instance:provisioning` y respuesta 409 `RESOURCE_QUOTA_EXCEEDED` ante saturación de cuota.

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

> [!NOTE]
> **Estado verificado (15/10/2026, backend `dbf7075`): implementada y con test en verde.** Pasa la suite `TestCrearInstancia_*` (6/6 subtests: creación de VM y LXC, respuesta `202 Accepted` con `{upid, tareaId, vmid}`, registro de tarea en `RUNNING` y restricción RBAC 403 Forbidden para operadores sin rol ADMIN).

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

---

## Verificación del 19/10/2026

Verificadas contra backend `8daaf87`.


#### `BAC-39` - Seguimiento de creación, autoasignación de permisos y eventos (`RF-07`, `RF-11`)

> [!NOTE]
> **Estado verificado (19/10/2026, backend `8daaf87`): implementada y con test en verde.** Pasa la suite de integración `TestEtapa2BAC39_SeguimientoCreacionYPermisos`. Al detectar fin exitoso (`exitstatus: "OK"`), inserta atómicamente `(creador_id, vmid, FULL_ACCESS)` en `permisos_instancia` y emite `TASK_FINISHED` por SSE.

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

> [!NOTE]
> **Estado verificado (19/10/2026, backend `8daaf87`): implementada y con test en verde.** Pasan `TestEtapa2BAC40_EndpointConfigRecursos` (validación de cuota libre con $\Delta\text{RAM}$, rechazo 409 `RESOURCE_QUOTA_EXCEEDED` ante excedente, detección de hot-plug y flag `requiresReboot: true` ante instancias `running` que no soportan hot-plug).

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

> [!NOTE]
> **Estado verificado (19/10/2026, backend `8daaf87`): implementada y con test en verde.** Pasa `TestEtapa2BAC41_AuditoriaAprovisionamientoYConfig`. Registra `CREATE_INSTANCE` con estado inicial `PENDING` y actualización final correlacionada por `tareaId`/`upid`, y `CONFIG_INSTANCE` con parámetros anteriores y nuevos en `detalles` (JSONB).

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

> [!NOTE]
> **Estado verificado (19/10/2026, backend `8daaf87`): implementada y con test en verde.** Pasan los 4 subtests de `TestEtapa2BAC50_SimuladorSnapshotsYRollback` (`CrearSnapshot`, `ListarSnapshots`, `RollbackSnapshot`, `EliminarSnapshot`) sobre el simulador `cmd/proxmox-simulador`.

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

> [!NOTE]
> **Estado verificado (19/10/2026, backend `8daaf87`): implementada y con test en verde.** Pasa `TestEtapa2BAC51_GuardaConcurrenciaInstanceBusy`. El middleware `RequireInstanceIdle` rechaza determinísticamente con HTTP 409 y código `INSTANCE_BUSY` ante solicitudes de ciclo de vida o reconfiguración si la instancia posee una tarea activa en ejecución.

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

> [!NOTE]
> **Estado verificado (19/10/2026, backend `8daaf87`): implementada y con test en verde.** Pasa `TestEtapa2BAC52_BrokerSSEEventosAlertaGlobal`. El modelo `Evento` admite `recursoTipo: "HOST"` sin `recursoId`, difundiendo alertas exclusivamente a usuarios con rol `ADMIN` y descartándolas silenciosamente para operadores.

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

> [!NOTE]
> **Estado verificado (19/10/2026, backend `8daaf87`): implementada y con test en verde.** Pasa `TestEtapa2BAC53_ParticionadoPreventivoEIndexacionAuditoria`. Las particiones `auditoria_2026_q3`, `auditoria_2026_q4` y `auditoria_2027_q1` y el índice `idx_auditoria_query` existen y garantizan *partition pruning* eficiente.

- **Área:** Backend / Base de Datos
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `BAC-18`, `BAC-41`.
- **Problema:** En PostgreSQL, si no existen las particiones para los trimestres de trabajo (Q3/Q4 2026), las inserciones de auditoría fallarán. Además, el streaming de `BAC-45` requiere índices compuestos para evitar escaneos secuenciales.
- **Entregable:**
  1. Script de migración SQL DDL que crea preventivamente las particiones trimestrales: `auditoria_2026_q3`, `auditoria_2026_q4`, `auditoria_2027_q1`.
  2. Índices compuestos para partition pruning: `CREATE INDEX idx_auditoria_query ON auditoria (fecha_hora DESC, usuario_id, accion)`.
- **Criterio de éxito:** `EXPLAIN ANALYZE` sobre una consulta por rango de fechas confirma que PostgreSQL solo lee la partición correspondiente (*partition pruning*) sin recorrer la tabla completa.