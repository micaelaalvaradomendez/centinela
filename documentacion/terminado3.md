

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

#### `BAC-47` - Despacho multicanal de alertas (SSE en tiempo real y correo estructurado)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-46`, `BAC-52` (`BRG-12-BAC`: broker SSE generalizado para alertas globales).
- **Problema:** Una vez detectada la saturación, todos los administradores conectados deben ser prevenidos en sus pantallas y por canal de contingencia (correo).
- **Entregable:**
  1. Emisión de eventos `ALERT_TRIGGERED` y `ALERT_RESOLVED` a través del canal `/api/events` hacia los clientes SSE.
  2. Despacho estructurado de correo electrónico a través del adaptador `EmailService`:
     - Asunto: `[ALERTA CRÍTICA] Saturación de Recursos en {Recurso}`.
     - Cuerpo con métricas exactas, fecha y recomendación operativa.
  3. Registro automático de la alerta en la tabla inmutable `auditoria`.
- **Criterio de éxito:** Al dispararse la saturación, los clientes SSE reciben el evento en menos de 1 s y el log del servidor registra el correo despachado por `EmailService`.

---

## Verificación del 30/10/2026

Verificadas contra backend `908522e` y frontend `594a4e0`.

#### `BAC-44` - Interceptor global y consultas complejas sobre auditoría particionada

> [!NOTE]
> **Estado verificado (30/10/2026, backend `908522e`): implementada y con tests en verde.** Pasan las pruebas unitarias y de integración (`TestAuditHandlerExcluirAccion`, `TestAuditHandlerFiltrosRF08`, `TestAuditoria_CursorTraduceFiltrosYFilas`, etc.). Soporta constantes de acción Etapa 3 (`SNAPSHOT_CREATE`, `SNAPSHOT_ROLLBACK`, `SNAPSHOT_DELETE`, `ALERT_TRIGGERED`, `CONFIG_UPDATED`), filtros combinados (`desde`, `hasta`, `usuarioId`, `recursoTipo`, `recursoId`, `page`, `limit`), rechazo de rangos invertidos con 400 y orden `fecha_hora DESC`.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-41`, `BAC-53`.
- **Problema:** La auditoría debe consolidar eventos de todas las etapas con filtros combinados eficientes aprovechando las particiones de PostgreSQL.
- **Entregable:**
  1. Extensión del interceptor de auditoría para registrar eventos de snapshots, alertas y configuración.
  2. Endpoint `GET /api/admin/audit` con parámetros combinados y paginación.
  3. Optimización de consulta SQL utilizando exclusión de particiones (`partition pruning`) por `fecha_hora`.
- **Criterio de éxito:** Consultas sobre registros históricos responden sin escaneos completos de tablas y validan el contrato.

#### `BAC-45` - Streaming de exportación estructurada de auditoría a CSV y JSON

> [!NOTE]
> **Estado verificado (30/10/2026, backend `908522e`): implementada y con tests en verde.** Pasan `TestAuditExportStreamingFormatosYVacio`, `TestAuditExportStreamingPermisosYErrores`, `TestAuditExportStreamingCancelaYCierraCursor`, `TestAuditExportStreamingFalloParcialAborta`. Transmisión de registros mediante cursor SQL `.Rows()` sin saturar memoria RAM, cabecera `Content-Disposition` con `auditoria_centinela_YYYYMMDD.csv|json` y sanitización contra inyección de fórmulas CSV (`=`, `+`, `-`, `@`).

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-44`.
- **Problema:** Exportar reportes extensos cargando todos los registros en memoria puede provocar caídas por Out of Memory (OOM).
- **Entregable:**
  1. Endpoint `GET /api/admin/audit/export` accesible solo para `ADMIN`.
  2. Parámetros `format=csv|json`, respetando filtros aplicados.
  3. Streaming HTTP por cursor de base de datos sin Find() global.
  4. Cabecera `Content-Disposition: attachment; filename="auditoria_centinela_YYYYMMDD.csv|json"`.
- **Criterio de éxito:** Exportaciones masivas se descargan en flujo continuo sin picos de RAM en el backend.

#### `BAC-48` - Hardening de autenticación: Cookies `HttpOnly`, Rate Limiting y Sanitización

> [!NOTE]
> **Estado verificado (30/10/2026, backend `908522e`): implementada y con tests en verde.** Pasan al 100% las pruebas de seguridad `TestAuthHardening_CookieSegura`, `TestAuthHardening_CookieSinSecureEnLocal`, `TestAuthHardening_LogoutBorraLaCookie`, `TestAuthHardening_RafagasPorIP` y `TestAuthHardening_Sanitizacion`. Emisión de `refreshToken` en cookie `centinela_refresh` con directivas `HttpOnly`, `Secure`, `SameSite=Strict`, `Path=/api/auth` y rate limiting de 5 intentos por minuto con HTTP 429.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-03`, `BAC-17B`.
- **Problema:** Evitar exposición de tokens de larga duración ante vectores XSS y proteger endpoints de autenticación contra fuerza bruta.
- **Entregable:**
  1. Emisión del `refreshToken` exclusivamente en cookie HTTP segura (`HttpOnly`, `Secure`, `SameSite=Strict`, `Path=/api/auth`).
  2. Middleware de Rate Limiting por IP para endpoints públicos con respuesta `429 Too Many Requests`.
  3. Sanitización estricta de encabezados y payloads JSON contra inyección de caracteres de control.
- **Criterio de éxito:** JavaScript no puede leer la cookie de refresh y peticiones abusivas consecutivas reciben 429.

#### `BAC-54` (`BRG-14-BAC`) - Extracción y Parametrización en `.env` para Servidores Heterogéneos

> [!NOTE]
> **Estado verificado (30/10/2026, backend `908522e`): implementada y con tests en verde.** Pasa la suite de aceptación `TestEtapa3BAC54_ParametrizacionEntornoServidoresHeterogeneos` y `TestRejectProtectedInstance_ListaDesdeEntorno`. Lectura dinámica de `PROXMOX_NODE`, `PROXMOX_URL` y `PROXMOX_PROTECTED_VMIDS` desde variables de entorno con fallback seguro ante listas vacías.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `INT-07`.
- **Problema:** En servidores con configuraciones heterogéneas (PVE 8 / 9), el nodo y los VMIDs protegidos varían y no deben requerir recompilación del backend.
- **Entregable:**
  1. Parametrización en `main.go` y `middleware/protected_instance.go`: lectura prioritaria de `PROXMOX_NODE`, `PROXMOX_URL` y `PROXMOX_PROTECTED_VMIDS`.
  2. Fallback seguro ante listas vacías sin fallar el arranque.
- **Criterio de éxito:** Modificar el `.env` altera el nodo y VMIDs protegidos sin recompilar el código.

#### `FRN-29` - Pestaña "Snapshots" en Detalle de Instancia y formulario de nuevo snapshot

> [!NOTE]
> **Estado verificado (30/10/2026, frontend `594a4e0`): implementada y con tests en verde.** Pasa la suite `test/front/frn-36-29-snapshots-ui.test.tsx` (7 de 7 pruebas en verde) y `frontend/centinela/src/test/snapshots.test.tsx` (3 de 3 pruebas en verde). Tabla interactiva de snapshots con fecha, descripción y estado de RAM, botón "Tomar Snapshot" con modal validando `snapname` alfanumérico y `description` <= 250, despacho asíncrono con HTTP 202 y recarga automática al recibir `TASK_FINISHED`.

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.0 h
- **Depende de:** `BAC-42`, `FRN-36`.
- **Problema:** La vista de detalle carecía de una interfaz para inspeccionar puntos de restauración y disparar nuevos snapshots.
- **Entregable:**
  1. Pestaña "Snapshots" dentro de `InstanceDetailPage.tsx`.
  2. Tabla interactiva de snapshots con metadatos y estado.
  3. Modal "Tomar Snapshot" con validación de nombre (`/^[a-zA-Z0-9_\-]+$/`).
  4. Manejo del ciclo asíncrono y feedback visual con spinner.
- **Criterio de éxito:** El usuario visualiza la lista de snapshots y puede crear uno nuevo viendo la lista actualizada reactivamente.

#### `FRN-30` - Modal destructivo de Rollback y eliminación con confirmación antierror

> [!NOTE]
> **Estado verificado (30/10/2026, frontend `594a4e0`): implementada y con tests en verde.** Pasa la suite `test/front/frn-30-snapshots-rollback-ui.test.tsx` (2 de 2 pruebas en verde, 100% ✅). Modal destructivo `RollbackSnapshotModal.tsx` con advertencia de pérdida de datos, input de seguridad que exige tipear el nombre idéntico del snapshot antes de habilitar "Revertir", y modal de confirmación para eliminación.

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `FRN-29`, `BAC-43`.
- **Problema:** Revertir un snapshot descarta datos posteriores; un clic accidental causaría pérdida irreparable.
- **Entregable:**
  1. Modal antierror `RollbackSnapshotModal.tsx` con advertencia en rojo.
  2. Input de confirmación exigiendo tipear el nombre textual del snapshot.
  3. Modal de confirmación para eliminación de snapshots.
  4. Bloqueo de la botonera mientras la tarea asíncrona se resuelve.
- **Criterio de éxito:** El botón de rollback solo es presionable tras tipear el nombre idéntico y muestra feedback reactivo.

#### `FRN-34` - Adaptación del cliente a cookies seguras de sesión y auditoría de accesibilidad

> [!NOTE]
> **Estado verificado (30/10/2026, frontend `594a4e0`): implementada y con tests en verde.** Pasa la suite `test/front/frn-34-cookies-a11y.test.tsx` (3 de 3 pruebas en verde, 100% ✅). `apiClient.ts` envía `credentials: 'include'` en peticiones al backend, el cliente no almacena `refreshToken` en `sessionStorage` ni `localStorage`, y los modales y componentes cumplen con directivas de accesibilidad WCAG (`alertdialog`, `aria-labelledby`).

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `BAC-48`.
- **Problema:** El cliente de API debe enviar credenciales automáticas en llamadas de refresh y la interfaz debe cumplir con estándares de accesibilidad WCAG.
- **Entregable:**
  1. Ajuste de `apiClient.ts` para enviar `credentials: 'include'`.
  2. Eliminación de almacenamiento manual de `refreshToken` en almacenamiento local del navegador.
  3. Modales accesibles con navegación por teclado y roles ARIA.
- **Criterio de éxito:** La rotación de tokens opera mediante cookies HttpOnly y las pruebas de accesibilidad pasan limpiamente.