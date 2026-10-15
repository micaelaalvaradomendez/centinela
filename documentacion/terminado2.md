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