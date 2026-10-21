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

---

## Verificación del 21/10/2026

Verificadas contra backend `a25d1b3` y frontend `b0a7105`.


#### `BAC-36` - Endpoint REST de métricas históricas de instancia (`timeframe=2h|24h`)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commits `ad6565f` y `16a7ba2`. Handler `InstanceMetricsHandler` registrado en `/api/instances/:instanceRef/metrics` protegido con `RequireAuth` y `RequireInstanceAccess`, con soporte de rangos temporales `2h` y `24h` y caché Redis con TTL de 30 segundos.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-34`.
- **Problema:** Al abrir la vista de una máquina, los gráficos necesitan renderizar de inmediato la tendencia histórica sin aguardar la llegada de los primeros paquetes del WebSocket.
- **Entregable:**
  1. Endpoint `GET /api/instances/:instanceRef/metrics?timeframe=2h|24h` protegido.
  2. Mapeo a intervalos RRD de Proxmox (`hour` y `day`).
  3. Caché Redis de 30 s para evitar saturar el hipervisor ante recargas de página.
- **Criterio de éxito:** Responde en menos de 100 ms desde caché con el arreglo histórico de métricas formateado para el consumo de gráficos.


### `FIX-107` - Cifrado y ofuscación reversible de VMIDs de Proxmox contra IDOR y enumeración de hipervisor (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commit `c51c77d`. Servicio de ofuscación simétrica en `internal/infrastructure/crypto` con clave secreta y middleware `ResolveInstanceRef` en rutas HTTP `/api/instances/:instanceRef`. La suite `TestFIX107_OfuscacionDeVMIDs` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 2,5 h
- **Depende de:** `FIX-106`.
- **Problema:** Los identificadores se exponían como enteros secuenciales (`100`, `101`), permitiendo ataques de enumeración (IDOR) y deducción de topología.
- **Entregable:**
  1. Utilidad `OfuscarVMID(vmid int)` y `DesofuscarVMID(hash string)` con clave secreta.
  2. Middleware de resolución en handlers de instancias `:instanceRef`.
  3. Respuestas 404 homogéneas ante tokens inválidos o inexistentes.
- **Criterio de éxito:** `/api/instances` expone tokens opacos y las acciones operativas operan transparentemente sin exponer números de VMID en red.


### `FIX-109` - Enmascaramiento de nombres de host, nodos y pools de almacenamiento (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commit `1d84e7b`. Mapeo lógico de nombres físicos de nodos (`proxmox` $\to$ alias lógico) y pools de storage para usuarios que no poseen rol `ADMIN`.

- **Área:** Backend
- **Asignado:** Lisandro / Tayra
- **Estimación:** 2,0 h
- **Depende de:** `FIX-107`.
- **Problema:** Se exponían nombres reales de red del servidor físico y volúmenes de almacenamiento locales (`local-lvm`, `ceph`), revelando la topología interna.
- **Entregable:**
  1. Diccionario de aliases lógicos para hipervisores y storages.
  2. Resolución transparente hacia identificadores físicos en llamadas a Proxmox.
  3. Restricción estricta de metadatos técnicos para rol `ADMIN`.
- **Criterio de éxito:** el frontend nunca recibe hostnames reales ni rutas físicas de disco en los JSONs de respuesta.


### `FIX-110` - Saneamiento estricto de mensajes de error y canales en tiempo real (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commit `c9deb62`. `mapearErrorProxmox` y canales en tiempo real filtran cualquier traza de rutas `/usr/share/perl5/...` o identificadores internos de Proxmox. `TestFIX110_SaneamientoDeErrores` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1,5 h
- **Depende de:** `FIX-105`, `FIX-106`.
- **Problema:** Respuestas de error 5xx exponían fragmentos crudos retornados por la API de Proxmox con rutas absolutas de código Perl y módulos de Linux.
- **Entregable:**
  1. Registro de trazas crudas exclusivamente en logs internos del servidor.
  2. Emisión de mensajes de error canónicos y estructurados (`PROXMOX_ERROR`, `PROXMOX_UNAVAILABLE`).
  3. Sanitización de eventos en `/api/events` y métricas WebSocket.
- **Criterio de éxito:** ningún error 4xx/5xx ni evento en tiempo real contiene trazas de archivos del sistema operativo anfitrión.


### `FIX-111A` - Blindaje contra secuestro de cuenta en cambio de email (`PUT /api/account/profile`) (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commits `dc50c0a`, `526d0f2` y `10e2c84`. Exige obligatoriamente `contrasenaActual` y `codigo2fa`, genera un token temporal en Redis y despacha confirmación por correo antes de aplicar la modificación en PostgreSQL, registrando en `auditoria`. `TestFIX111A_BlindajeCambioDeEmail` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Importancia:** 🔴 **CRÍTICA** (CWE-284, OWASP A01:2021).
- **Problema:** Un usuario autenticado podía mutar directamente su email sin verificar contraseña ni 2FA y sin confirmación por correo, facilitando el secuestro definitivo de cuenta.
- **Entregable:**
  1. Validación obligatoria de contraseña actual y TOTP 2FA.
  2. Generación de token criptográfico de confirmación en Redis (TTL 15 min).
  3. Modificación del correo únicamente tras la validación del código/enlace en el nuevo buzón.
  4. Registro obligatorio e inmutable en `auditoria`.
- **Criterio de éxito:** cambios de email no confirmados o sin credenciales son rechazados con 400/401 y auditados.


### `FIX-112` - Middleware de Rate Limiting y protección contra fuerza bruta en autenticación (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commits `f738f1c` y `ef53022`. Implementa rate limiting distribuido mediante Redis en `/login`, `/2fa/verify` (máximo 3 intentos por token pre-auth con revocación inmediata) y `/password/forgot`, además de serialización transaccional con `FOR UPDATE` en `ConfirmarRecuperacionContrasena`. `TestFIX112_MitigacionFuerzaBrutaRateLimit` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Lisandro / Nico
- **Estimación:** 2,5 h
- **Importancia:** 🔴 **CRÍTICA** (CWE-307, OWASP A07:2021).
- **Problema:** Rutas de autenticación sin cotas de intentos, permitiendo ataques de diccionario, fuerza bruta sobre tokens temporales y saturación SMTP.
- **Entregable:**
  1. Middleware `RateLimiter` con algoritmo sliding window en Redis.
  2. Revocación de sesión temporal de Redis al tercer intento erróneo de 2FA.
  3. Cota estricta de solicitudes en reseteo de contraseña.
  4. Bloqueo pesimista de fila (`FOR UPDATE`) en base de datos.
- **Criterio de éxito:** rechazo HTTP 429 ante excesos de intentos y destrucción de sesiones pre-auth vulneradas.


### `FIX-114` - Revocación obligatoria de sesiones activas en cambio voluntario de contraseña (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commits `638bb5e`, `a0f5465` y `c99c9f5`. `CambiarContrasena` en `user_service.go` invalida todas las sesiones activas concurrentes en Redis y PostgreSQL y notifica por bus de eventos para cerrar WebSockets y SSE. `TestFIX114_RevocacionSesionesCambioPassword` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,0 h
- **Importancia:** 🟠 **ALTA** (CWE-613, OWASP A07:2021).
- **Problema:** Al cambiar voluntariamente la contraseña en `PUT /api/account/password`, las sesiones activas en otros dispositivos permanecían vigentes, impidiendo expulsar a un atacante.
- **Entregable:**
  1. Invocación de revocación de sesiones concurrentes preservando la sesión solicitante (o forzando re-login).
  2. Emisión de `MensajeSesionesRevocadas` para cortar canales en tiempo real.
  3. Registro en `auditoria`.
- **Criterio de éxito:** tokens y sesiones anteriores responden `401 TOKEN_REVOKED`.


### `FIX-115A` - Ampliación de cota máxima de contraseñas a 128 caracteres en Backend y DTOs (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commit `b7cdaf0`. `ValidarComplejidadContrasena` en `password.go` y los tags de validación Gin en DTOs se ampliaron a `max=128`. `TestFIX115A_LongitudPassword128Chars` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 0,5 h
- **Importancia:** 🟠 **MEDIA / ALTA** (CWE-521, NIST SP 800-63B).
- **Problema:** Restringir contraseñas a un máximo de 12 caracteres impedía el uso de frases seguras (*passphrases*) o contraseñas de gestores de contraseñas.
- **Entregable:**
  1. Cota superior de 128 caracteres en `ValidarComplejidadContrasena`.
  2. Ajuste de tags de validación `max=128` en DTOs de `auth_handler.go`, `user_port.go`, etc.
- **Criterio de éxito:** el sistema acepta contraseñas de alta longitud sin truncarlas ni rechazarlas.


### `FIX-116` - Sanitización estricta y prevención de inyección de parámetros en Proxmox (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commit `a25d1b3`. Validación defensiva de formato en interfaces de red (`Bridge`), hostnames RFC 1123 y rechazo de caracteres delimitadores de directivas (`,`) antes de contactar a Proxmox. `TestFIX116_InyeccionDeParametrosProxmox` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Lisandro / Tayra
- **Estimación:** 2,0 h
- **Importancia:** 🟠 **MEDIA** (CWE-88, CWE-20).
- **Problema:** La API de Proxmox utiliza comas (`,`) como separador de directivas de hardware; caracteres no saneados en nombres de bridge o parámetros de red permitían inyectar opciones de configuración no autorizadas.
- **Entregable:**
  1. Validación regex estricta para nombres de bridge (`^[a-zA-Z0-9_.-]+$`).
  2. Validación de hostnames RFC 1123.
  3. Rechazo de `,` en `client.go`.
- **Criterio de éxito:** directivas manipuladas son rechazadas con `400 INVALID_REQUEST` antes de enviar tráfico a Proxmox.


### `FIX-117A` - Restricción de acceso y segregación de telemetría global del host para operadores (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commit `94e95a2`. En `GET /api/node/status` y `GET /api/node/history`, se detecta el rol del usuario autenticado: para rol `OPERATOR` se omite el almacenamiento total de los discos del host y el inventario global de máquinas ajenas. `TestFIX117A_SegregacionDeTelemetriaGlobal` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1,5 h
- **Importancia:** 🟡 **MEDIA** (CWE-200).
- **Problema:** Operadores que solo gestionaban una máquina de prueba podían consultar la carga global de CPU/RAM de todo el hipervisor físico y el conteo de máquinas de todos los demás departamentos.
- **Entregable:**
  1. Segregación de payload según `ContextKeyRol`.
  2. Ocultamiento de storage global y máquinas no asignadas para operadores.
- **Criterio de éxito:** operadores no reciben la topología global ni capacidad física de los discos del anfitrión.


### `FIX-119` - Validación de existencia de VMIDs y bloqueo de instancias protegidas en asignación de permisos (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `a25d1b3`): implementada y con test en verde.** Commit `5c72bbb`. `AsignarPermisos` en `user_service.go` consulta el inventario de Proxmox, bloquea la asignación de identificadores inexistentes y rechaza con `403 PROTECTED_INSTANCE_CANNOT_BE_ASSIGNED` ante recursos de `PROXMOX_PROTECTED_VMIDS`. `TestFIX119_ValidacionDeVMIDsEnAsignacionDePermisos` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Lisandro / Tayra
- **Estimación:** 1,5 h
- **Importancia:** 🟡 **MEDIA** (CWE-20).
- **Problema:** Se podían persistir permisos sobre VMIDs que no existían o sobre máquinas protegidas de la infraestructura interna de Centinela.
- **Entregable:**
  1. Inyección de catálogo de instancias en `AsignarPermisos`.
  2. Validación de existencia en Proxmox (`400 INVALID_REQUEST`).
  3. Rechazo de VMIDs protegidos (`403`).
- **Criterio de éxito:** la tabla de permisos solo contiene máquinas existentes y operables.


#### `FRN-24` - Wizard de creación: Paso 1 (Tipo de instancia) y Paso 2 (Recursos y Cuotas en Vivo)

> [!WARNING]
> **Estado verificado (21/10/2026, frontend `b0a7105`): implementada con problemas.** Commits `65765e4` a `8688b5d`. Se desarrollaron los componentes del asistente, pero la suite `test/front/frn-24-25-create-wizard.test.tsx` falla en las comprobaciones de medidores reactivos de cuotas y bloqueo de avance al exceder la RAM libre del nodo. Corrección abierta en **`FIX-120`** ([futuro2.md](futuro2.md)).

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.5 h
- **Depende de:** `BAC-33`, `FIX-60`.
- **Entregable:** selector de tarjetas QEMU/LXC, sliders de vCPU/RAM/Disco y medidor reactivo de cuota del host.
- **Criterio de éxito:** el usuario no puede avanzar al Paso 3 si los recursos exceden la memoria o disco libres.


#### `FRN-25` - Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen)

> [!WARNING]
> **Estado verificado (21/10/2026, frontend `b0a7105`): implementada con problemas.** Commits `3723b40` a `e04bfcb`. Se incorporaron los selectores de storage y resumen, pero fallan las aserciones de filtrado de imágenes según tipo (`vztmpl` vs `iso`) y la validación de red IP estática. Corrección abierta en **`FIX-121`** ([futuro2.md](futuro2.md)).

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.0 h
- **Depende de:** `FRN-24`, `BAC-31`.
- **Entregable:** selector de storages reales, filtrado de templates/ISOs según tipo, configuración de red validada y ficha de resumen.
- **Criterio de éxito:** validación estricta de todos los campos antes de habilitar el botón "Crear Instancia".


### `FIX-118` - Eliminación de credenciales por defecto hardcodeadas en migraciones de base de datos (Backend)

> [!NOTE]
> **Estado verificado (21/10/2026, backend `1801e33`): implementada y con test en verde.** Commit `1801e33`. `credencialesIniciales` en `db.go` suprime la contraseña estática `Admin123!`, desactiva el sembrado automático si `APP_ENV=production`, genera contraseñas temporales aleatorias con `crypto.GenerarContrasenaTemp()` en desarrollo y permite parametrizar mediante `INITIAL_ADMIN_EMAIL` / `INITIAL_ADMIN_PASSWORD`. La suite `TestFIX118_SinCredencialesPorDefecto` pasa al 100% en verde.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,0 h
- **Problema:** inicialización desatendida con credenciales públicas hardcodeadas en repositorios.
- **Criterio de éxito:** deshabilitar en producción y parametrizar vía variables de entorno.


#### `INF-10` - Verificación y disponibilidad de plantillas LXC e ISOs en el servidor Proxmox real

> [!WARNING]
> **Estado verificado (21/10/2026): pendiente de verificación en servidor Proxmox físico real.** Tarea de infraestructura pendiente en el servidor Proxmox real con storages `local:vztmpl` e `iso`. Corrección y seguimiento abierto en **`FIX-123`** ([futuro2.md](futuro2.md)).

- **Área:** Infraestructura
- **Asignados:** Nico y Lucas
- **Estimación:** 1.5 h
- **Criterio de éxito:** el backend lista plantillas e ISOs reales sin errores de permisos.
