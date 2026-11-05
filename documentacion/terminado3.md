

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

---

## Verificación del 01/11/2026 (Etapa Adaptación Proxmox PVE 8.4)

Verificadas contra backend `ec5922c` y frontend `ce38ea9`.

#### `ADP-01` (`BAC-PVE8-01`) - Ajuste de Parámetros de Aprovisionamiento QEMU y LXC en `client.go`

> [!NOTE]
> **Estado verificado (01/11/2026, backend `ec5922c`): implementada y con tests en verde.** Commit `7505aa5`. Inyección de `ostype=l26` y `agent=1` para QEMU, y `unprivileged=1` con `features=nesting=1` para LXC en los parámetros urlencoded de creación de instancias. Pasan las pruebas unitarias en `client_test.go` y la suite `TestClient_AprovisionarInstancia`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** Cierre de Etapa 3 (`INT-07`).
- **Problema:** `client.go` no incluía `ostype`, `agent=1`, `unprivileged=1` ni `nesting=1` al llamar a `POST /nodes/{node}/qemu` y `/lxc`, provocando que en PVE 8.4 las VMs no activaran el socket de red del agente y los contenedores no arrancaran con systemd.
- **Entregables:**
  1. En `backend/internal/adapters/secondary/proxmox/client.go:parametrosCreacion`:
     - Se agregaron `params.Set("ostype", "l26")` y `params.Set("agent", "1")` para QEMU.
     - Se agregaron `params.Set("unprivileged", "1")` y `params.Set("features", "nesting=1")` para LXC.
  2. Pruebas unitarias en `client_test.go` verificando que los parámetros urlencoded incluyan estas directivas.
- **Criterio de éxito:** `go test -v -run TestClient_AprovisionarInstancia ./internal/adapters/secondary/proxmox` pasa al 100% validando la inclusión de `agent=1`, `ostype=l26`, `unprivileged=1` y `nesting=1`.

#### `ADP-02` (`BAC-PVE8-02`) - Parámetros de Purga Profunda en Eliminación de Instancias

> [!NOTE]
> **Estado verificado (01/11/2026, backend `ec5922c`): implementada y con tests en verde.** Commit `084ad9a`. Incorporación de query string `?purge=1&destroy-unreferenced-disks=1` en la llamada HTTP DELETE de `EliminarInstancia`. Pasan las pruebas unitarias de purga en `client_test.go`.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `ADP-01`.
- **Problema:** En PVE 8.4, eliminar una instancia mediante `DELETE /nodes/{node}/{tipo}/{vmid}` sin `purge=1` falla si la instancia tiene metadatos en HA o deja discos huérfanos sin desasociar.
- **Entregables:**
  1. En `client.go:EliminarInstancia`:
     - Query string `?purge=1&destroy-unreferenced-disks=1` anexada a la llamada HTTP DELETE.
  2. Pruebas unitarias en `client_test.go` asegurando que la URL de borrado contenga los parámetros de purga.
- **Criterio de éxito:** La llamada de borrado emite la query esperada y la prueba unitaria pasa en verde.

#### `ADP-03` (`BAC-PVE8-03`) - Flexibilización y Tolerancia en Parser de Interfaces y Configuración

> [!NOTE]
> **Estado verificado (01/11/2026, backend `ec5922c`): implementada y con tests en verde.** Commit `166c108`. Tolerancia en `ObtenerConfig` para digests con memoria como string o entero, y en `ObtenerInterfaces` para respuestas donde `ip-addresses` sea null, ausente o use CIDRs tradicionales en `inet`. Pasan las suites `TestClient_ObtenerInterfaces_LXC` y `TestObtenerConfig_QemuYLXC`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `ADP-01`.
- **Problema:** Si PVE 8.4 entrega `memory` como entero en `qemu config` o `prefix` como entero en `lxc interfaces`, no debe ocurrir ningún desajuste en el cómputo de digests de concurrencia ni en la extracción de IPs.
- **Entregables:**
  1. `ObtenerConfig` mantiene consistencia en el digest tanto si `memory` llega como string o como entero en el JSON.
  2. `ObtenerInterfaces` gestiona respuestas donde `ip-addresses` sea `null`, esté ausente o use CIDRs tradicionales en `inet`.
  3. Tests de regresión en `client_test.go` con payloads de PVE 8.4 reales.
- **Criterio de éxito:** Suites `TestClient_ObtenerInterfaces_LXC` y `TestObtenerConfig_QemuYLXC` pasan sin errores contra payloads mixtos (string e integer).

#### `ADP-04` (`BAC-PVE8-04`) - Parametrización Dinámica de VMIDs Protegidos por Entorno

> [!NOTE]
> **Estado verificado (01/11/2026, backend `ec5922c`): implementada y con tests en verde.** Commit `c09a867`. `PROXMOX_PROTECTED_VMIDS` admite listas vacías o personalizadas vía variable de entorno sin romper el arranque de la API ni bloquear instancias no protegidas.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** Ninguna.
- **Problema:** El rango `100-105` está hardcodeado por defecto para la infraestructura de Centinela. En el servidor final con PVE 8.4.21, Centinela podría usar otros IDs y no debe bloquear VMs preexistentes del hipervisor.
- **Entregables:**
  1. En `backend/internal/adapters/primary/http/middleware/protected_instance.go`:
     - Verificación de que `PROXMOX_PROTECTED_VMIDS` permita listas vacías o personalizadas mediante `.env` sin romper el arranque de la API.
  2. Actualización de `SETUP.md` y `.env.example` documentando la configuración de VMIDs protegidos para el servidor productivo.
- **Criterio de éxito:** Si `PROXMOX_PROTECTED_VMIDS=""` o contiene IDs arbitrarios (ej. `"200,201"`), la API protege exactamente esos IDs y permite gestionar los IDs 100-105 sin retornar 403.

#### `ADP-05` (`SIM-PVE8-01`) - Soporte de Modo Dual de Versión en Simulador (`PROXMOX_SIM_TARGET_VERSION`)

> [!NOTE]
> **Estado verificado (01/11/2026, backend `ec5922c`): implementada y con tests en verde.** Commit `ec5922c`. Configuración `TargetVersion` en simulador con default `8.4.21`. Respuestas condicionadas en `/nodes/{node}/status` (firma `pve-manager/8.4-21` y kernel `Linux 6.8.8-2-pve`), validación de `ostype` y `agent` en creación QEMU, y requerimiento de `purge=1` en borrado. Pasan las pruebas de contrato en `simulador_test.go`.

- **Área:** Backend / Testing
- **Asignado:** Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `ADP-01`, `ADP-03`.
- **Problema:** El simulador reportaba rígidamente `pve-manager/9.2.2` y kernel 7.0.2 en `/nodes/{node}/status`, y no validaba los parámetros de `agent` ni `ostype` al crear VMs.
- **Entregables:**
  1. En `backend/cmd/proxmox-simulador/simulador.go` y `config.go`:
     - Variable de configuración `TargetVersion` (default `"8.4.21"`).
     - `/nodes/{node}/status` responde `pveversion: "pve-manager/8.4-21/12345678"` y `kversion: "Linux 6.8.8-2-pve"`.
     - `POST /nodes/{node}/qemu` requiere `ostype` y `agent` para habilitar el socket del guest agent.
     - `DELETE` exige `purge=1` si la instancia tiene el flag HA activado.
- **Criterio de éxito:** `PROXMOX_SIM_TARGET_VERSION=8.4.21 go run ./cmd/proxmox-simulador` expone las firmas y validaciones de PVE 8.4.21.

#### `ADP-08` (`FRN-PVE8-02`) - Verificación de Renderizado de Estado de Storage y Nodos en Dashboard

> [!NOTE]
> **Estado verificado (01/11/2026, frontend `ce38ea9` / `877eddd`): implementada y con tests en verde.** Commit `8d66b75` (PR #177). Manejo robusto en `Dashboard.tsx`, `StorageMeter.tsx` e `InstanceResourcesPanel.tsx` ante variaciones de redondeo y ausencia de content type `import`. Pasa `test/front/adp-08-dashboard-storage-pve8.test.tsx` (2/2 en verde).

- **Área:** Frontend
- **Asignado:** Cristian / Belinda
- **Estimación:** 1.5 h
- **Depende de:** `ADP-07`.
- **Problema:** En PVE 8.4, los storages reportan métricas sin el content type `import` y con ligeras variaciones de redondeo en bytes disponibles.
- **Entregables:**
  1. En `Dashboard.tsx` y componentes de medición (`StorageMeter.tsx`, `InstanceResourcesPanel.tsx`):
     - Cálculos de porcentaje de uso tolerantes sin NaN ni desbordamientos ante los tamaños de bloque devueltos por PVE 8.4.
- **Criterio de éxito:** Tarjetas de almacenamiento y nodo renderizan correctamente en el Dashboard con datos simulados de PVE 8.4.

#### `FRN-32` - Componente y flujo de exportación de reportes de auditoría

> [!NOTE]
> **Estado verificado (03/11/2026, frontend `877eddd`): implementada y con tests en verde.** Commit `404c784` (PR #180). `AuditExportButton.tsx` y `auditService.ts` implementan el menú desplegable CSV/JSON, estado de carga, notificación Toast y transferencia de filtros activos. Pasan 4/4 en `test/front/etapa3-pendientes.test.tsx` y `src/test/audit-service.test.ts`.

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 1.5 h
- **Depende de:** `FRN-31`, `BAC-45`.
- **Problema:** Los operadores de cumplimiento requieren descargar los registros en archivos locales para archivo legal o análisis externo.
- **Entregable:**
  1. Botón "Exportar registros" con menú desplegable para elegir "Exportar como CSV" o "Exportar como JSON".
  2. Integración con `exportAuditLogs(filtros, format)` transfiriendo los filtros actualmente activos en pantalla.
  3. Estado de carga en el botón durante la generación del reporte y disparo automático de la descarga nativa en el navegador.
  4. Notificación Toast al completar la descarga.
- **Criterio de éxito:** Al hacer clic en exportar, se descarga el archivo con los filtros aplicados sin recargar la página ni bloquear la UI.

#### `BAC-55` (`BRG-15-SIM`) - Switch de Perfil de Versión en `proxmox-simulador` (`PROXMOX_SIM_TARGET_VERSION`)

> [!NOTE]
> **Estado verificado (03/11/2026, backend `fb8a41c`): implementada y con tests en verde.** Commits `61b973d` y `ec5922c`. El simulador lee `PROXMOX_SIM_TARGET_VERSION` y conmuta firmas de versión y kernel entre PVE `8.4.21` y `9.2.2`. Pasan `TestADP05_SimuladorModoDualVersion` y `cmd/proxmox-simulador`.

- **Área:** Backend / Simulador
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-54`.
- **Problema:** Se requería que el simulador admitiera perfiles de versión mediante variable de entorno para contrastar pruebas tanto en versión 9 como en versión 8.4.21.
- **Entregable:**
  1. Lectura de `PROXMOX_SIM_TARGET_VERSION` en `backend/cmd/proxmox-simulador`.
  2. Respuestas condicionadas en `nodes/{node}/status` y validaciones de contrato por versión.
- **Criterio de éxito:** `PROXMOX_SIM_TARGET_VERSION=8.4.21` emite la cabecera y el JSON de status del nodo correspondientes a la versión 8.4.21.

#### `ADP-06` (`SIM-PVE8-02`) - Actualización del Catálogo Inicial de Plantillas e ISOs a Debian 12 Bookworm

> [!NOTE]
> **Estado verificado (03/11/2026, backend `fb8a41c`): implementada y con tests en verde.** Commit `aee232a`. `backend/cmd/proxmox-simulador/storage.go` expone `debian-12-standard_12.7-1_amd64.tar.zst`, `ubuntu-22.04-standard_22.04-1_amd64.tar.zst` y `debian-12.7.0-amd64-netinst.iso`, sin el content type `import` en modo PVE 8.4. Pasan `TestADP06_CatalogoDebian12Simulador` y `simulador_test.go`.

- **Área:** Backend / Simulador
- **Asignado:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `ADP-05`.
- **Problema:** El inventario de plantillas de `storage.go` exponía imágenes Debian 13 inexistentes en repositorios oficiales de PVE 8.4.
- **Entregables:**
  1. Reemplazo en `backend/cmd/proxmox-simulador/storage.go` por plantillas e ISOs de Debian 12 Bookworm y Ubuntu 22.04, removiendo `import` en modo PVE 8.4.
  2. Ajuste de pruebas unitarias en `simulador_test.go`.
- **Criterio de éxito:** `GET /nodes/{node}/storage/local/content` devuelve las plantillas de Debian 12 y `simulador_test.go` pasa al 100%.

#### `FIX-128` (`FRN-31`) - Opciones completas de acciones de auditoría Etapa 3 y selectores en `Auditoria.tsx` (Frontend)

> [!NOTE]
> **Estado verificado (03/11/2026, frontend `877eddd`): implementada y con tests en verde.** Commit `b38ab17` (PR #181). `auditActionLabels.ts` incluye `SNAPSHOT_DELETE`, `ALERT_TRIGGERED` y `SNAPSHOT_ROLLBACK: 'Revertir snapshot'`, y `Auditoria.tsx` integra `AuditExportButton`. Pasan 4/4 pruebas en `test/front/frn-31-audit-etapa3.test.tsx`.

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-31`, `BAC-44`.
- **Problema técnico detallado:** `auditActionLabels.ts` omitía `SNAPSHOT_DELETE` y `ALERT_TRIGGERED`, rotulaba `SNAPSHOT_ROLLBACK` sin vocabulario de reversión y `Auditoria.tsx` no ofrecía exportación JSON.
- **Entregable:**
  1. Incorporación de `SNAPSHOT_DELETE`, `ALERT_TRIGGERED` y `'Revertir snapshot'` en `auditActionLabels.ts`.
  2. Integración de `AuditExportButton` (CSV y JSON) en `Auditoria.tsx`.
- **Criterio de éxito:** La suite `test/front/frn-31-audit-etapa3.test.tsx` pasa 4 de 4 pruebas en verde (100% ✅).

#### `FIX-132` - Discriminación de errores 400 en recuperación de contraseña para no quemar intentos de OTP ante contraseñas inválidas o duplicadas (`CU-12`) (Frontend)

> [!NOTE]
> **Estado verificado (03/11/2026, frontend `877eddd`): implementada y con tests en verde.** Commit `3d06236` (PR #179). `RecoverPassword.tsx` discrimina `SAME_PASSWORD`, `WEAK_PASSWORD` y `PASSWORD_RECENTLY_USED` mostrando el error en el campo de contraseña sin incrementar `failedAttempts`. Pasan 4/4 en `test/front/etapa3-pendientes.test.tsx` y 7/7 en `test/front/recover-password.test.tsx`.

- **Área:** Frontend
- **Asignado:** Luz / Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-34`.
- **Problema técnico detallado:** `RecoverPassword.tsx` capturaba cualquier `error.status === 400` como fallo de código OTP, incrementando `failedAttempts` y bloqueando el formulario al 3.er intento ante `SAME_PASSWORD` o `WEAK_PASSWORD`.
- **Entregable:** Discriminación explícita de códigos de error de contraseña (`SAME_PASSWORD`, `WEAK_PASSWORD`, `PASSWORD_RECENTLY_USED`) respecto de `RESET_FAILED`.
- **Criterio de éxito:** Un error de política o repetición de contraseña no incrementa los intentos fallidos del OTP ni bloquea el flujo.

#### `FIX-133` - Corrección de validación diferencial de cuota de vCPU en escalado vertical para permitir CPU Overcommit (`RF-10` / `BAC-40` / `CU-10`) (Backend)

> [!NOTE]
> **Estado verificado (03/11/2026, backend `fb8a41c`): implementada y con tests en verde.** Commit `fb8a41c`. `AmpliarInstancia` en `cuotas_service.go` retorna inmediatamente si `deltaMemoria <= 0` y valida únicamente la memoria libre del nodo sin bloquear por la suma lineal de vCPUs. Pasan `TestFIX133_CuotaVCPUSinOvercommit` y las pruebas de `internal/core/services`.

- **Área:** Backend
- **Asignado:** Lucas / Micaela
- **Estimación:** 1.5 h
- **Depende de:** `BAC-40`, `RF-10`, `CU-10`.
- **Problema técnico detallado:** `AmpliarInstancia` rechazaba cualquier incremento de vCPU cuando la suma de vCPUs asignadas en el nodo igualaba o superaba los núcleos físicos (`coresDisponibles <= 0`), impidiendo el *CPU Overcommit* estándar de Proxmox VE e incumpliendo la regla de `RF-10` ($\Delta\text{RAM} \le 0$ debe autorizarse directamente).
- **Entregable:** Validación diferencial exclusiva sobre $\Delta\text{RAM} > 0$ en `cuotas_service.go` y `configuracion_instancias.go`.
- **Criterio de éxito:** Aumentar vCPUs con $\Delta\text{RAM} \le 0$ se autoriza sin errores `RECURSOS_INSUFICIENTES`.

#### `FIX-138` - Ampliación de caracteres especiales admitidos en validación de contraseñas conforme a NIST SP 800-63B (`RF-01` / `FIX-115A/B`) (Backend / Criptografía)

> [!NOTE]
> **Estado verificado (03/11/2026, backend `fb8a41c`): implementada y con tests en verde.** Commit `3fe055b`. `password.go` exporta `CaracteresEspeciales = "!@#$%^&*-_=+.~,;:'\"()[]{}<>\\/|?"` y admite espacios en frases de contraseña de más de 12 caracteres. Pasan `TestFIX138_CaracteresEspecialesNIST` y `internal/infrastructure/crypto`.

- **Área:** Backend / Criptografía
- **Asignado:** Lucas / Micaela
- **Estimación:** 1.0 h
- **Depende de:** `RF-01`, `FIX-115A`, `FIX-115B`.
- **Problema técnico detallado:** `ValidarComplejidadContrasena` solo admitía `!@#$%^&*-_=+`, rechazando contraseñas generadas con `.`, `,`, `?`, `~`, `()`, `[]`, `/` o espacios en passphrases largas.
- **Entregable:** Ampliación de `CaracteresEspeciales` y soporte de espacio en passphrases $> 12$ caracteres en `internal/infrastructure/crypto/password.go`.
- **Criterio de éxito:** El backend acepta contraseñas con símbolos ASCII estándar y frases de contraseña largas (`TestFIX138_CaracteresEspecialesNIST` en verde).

#### `FIX-140` - Trazabilidad y visibilidad de proveedor de email en modo Mock para prevenir falsos éxitos de entrega (`RF-09` / `RF-11` / `RF-13`) (Backend)

> [!NOTE]
> **Estado verificado (03/11/2026, backend `fb8a41c`): implementada y con tests en verde.** Commit `4e4d2c4`. `GET /api/version` expone `emailProvider` y los handlers que despachan correo en modo `mock` inyectan la cabecera `X-Centinela-Email-Provider: mock`. Pasa `TestFIX140_ProveedorDeEmailVisible`.

- **Área:** Backend
- **Asignado:** Lucas / Micaela
- **Estimación:** 1.0 h
- **Depende de:** `RF-09`, `RF-11`, `RF-13`.
- **Problema técnico detallado:** En entornos sin SMTP configurado, el backend operaba en modo `mock` retornando `200 OK` sin indicar en la respuesta HTTP ni en `/api/version` que los correos no salían a la red.
- **Entregable:**
  1. Cabecera `X-Centinela-Email-Provider: mock` en respuestas de `auth_handler.go`, `user_handler.go` y `account_handler.go` cuando el proveedor es mock.
  2. Campo `"emailProvider"` en `GET /api/version`.
- **Criterio de éxito:** `TestFIX140_ProveedorDeEmailVisible` pasa al 100% en verde.

#### `FIX-131` - Incorporación de métricas y gráfico de Tráfico de Red (`netin`/`netout`) en detalle de instancia (`CU-05`) (Frontend)

> [!NOTE]
> **Estado verificado (04/11/2026, frontend `1e5ec49`): implementada y con tests en verde.** Commits `fa92405` y `26af24f` (PR #182). `InstanceMetricPoint` declara `netInBps` y `netOutBps`, `instanceDetailService.ts` los normaliza y `InstanceResourceUsage.tsx` renderiza las 4 tarjetas y curvas (`CPU`, `Memoria`, `Disco`, `Tráfico de Red`). Pasan 3/3 subtests de `FIX-131` en `test/front/etapa3-pendientes.test.tsx` y 9/9 en `test/front/frn-23-metrics-charts.test.tsx`.

- **Área:** Frontend
- **Asignado:** Belinda / Cristian / Luz
- **Estimación:** 1.5 h
- **Depende de:** `FRN-22`, `FRN-23`.
- **Entregable:** Propiedades `netInBps`/`netOutBps` en `instanceDetail.ts`, normalización en `instanceDetailService.ts` y 4.ª columna `"Tráfico de Red"` en `InstanceResourceUsage.tsx`.
- **Criterio de éxito:** La grilla de uso de recursos renderiza las 4 métricas (`CPU`, `Memoria`, `Disco` y `Tráfico de Red`) con sus gráficos de evolución temporal.

#### `FIX-134` - Soporte operativo para instancias en estado `paused` y prevención de deadlocks en controles de energía y borrado (`RF-03` / `RF-04` / `CU-04`) (Backend)

> [!NOTE]
> **Estado verificado (04/11/2026, backend `08126ae`): implementada y con tests en verde.** Commit `42093c3` (`fix: admitir instancias pausadas en controles de energía y borrado`). Se implementaron `esEstadoCompatible` admitiendo `paused` en `stop` y `resume`, el método `Resume` en `ProxmoxPort` / `client.go` / simulador y la admisión de `paused` al eliminar instancia. Pasan 4/4 subtests de `TestFIX134_InstanciasPausadas` en `test/back/etapa3_fixes_pendientes_test.go` y la suite interna del backend.

- **Área:** Backend
- **Asignado:** Lucas / Nico / Tayra
- **Estimación:** 2.0 h
- **Depende de:** `RF-03`, `RF-04`, `CU-04`.
- **Entregable:** Soporte de acción `resume` (`/status/resume`), admisión de estado `paused` en `stop` y en `EliminarInstancia`, e implementación de `Resume` en `ProxmoxPort` y cliente Proxmox.
- **Criterio de éxito:** `TestFIX134_InstanciasPausadas` pasa 4/4 en verde.

#### `FIX-135` - Reconciliación con tareas de Proxmox en Evaluador de Alertas para prevenir falsos positivos de `instance_down` ante paradas directas en el hipervisor (`RF-11` / `CU-11`) (Backend)

> [!NOTE]
> **Estado verificado (04/11/2026, backend `08126ae`): implementada y con tests en verde.** Commit `5f8bb5e` (`fix: conciliar paradas externas de Proxmox para evitar falsas alertas`). `alert_evaluator_worker.go` verifica tareas recientes del nodo Proxmox (`qmstop`, `qmshutdown`, `vzstop`, `vzdump`), registra `SYSTEM_INSTANCE_STOPPED_EXTERNALLY` y parametrizó `ALERT_IGNORE_EXTERNAL_SHUTDOWN`. Pasan 3/3 subtests de `TestFIX135_AlertasParadasExternas` en `test/back/etapa3_fixes_pendientes_test.go` y `alert_evaluator_test.go`.

- **Área:** Backend
- **Asignado:** Lucas / Belinda / Tayra
- **Estimación:** 2.0 h
- **Depende de:** `RF-11`, `CU-11`.
- **Entregable:** Consulta secundaria a tareas de Proxmox en `explicadaPorCentinela`, auditoría `SYSTEM_INSTANCE_STOPPED_EXTERNALLY` y variable `ALERT_IGNORE_EXTERNAL_SHUTDOWN`.
- **Criterio de éxito:** `TestFIX135_AlertasParadasExternas` pasa 3/3 en verde.

#### `FIX-136` - Parametrización por variables de entorno de umbrales hardcodeados y persistencia en el Evaluador de Alertas (`RF-11` / `RNF-07`) (Backend / Infraestructura)

> [!NOTE]
> **Estado verificado (04/11/2026, backend `08126ae`): implementada y con tests en verde.** Commit `08126ae` (`fix: parametrizar umbrales y persistencia de alertas por entorno`). Se incorporó `ConfigAlertas` en `alert_evaluator_worker.go`, lectura de `ALERT_THRESHOLD_HOST_RAM`, `ALERT_THRESHOLD_HOST_CPU`, `ALERT_THRESHOLD_HOST_DISK`, `ALERT_THRESHOLD_VM_RAM`, `ALERT_THRESHOLD_VM_CPU`, `ALERT_THRESHOLD_RECOVERY` y `ALERT_CYCLES_COUNT` en `cmd/api/main.go` y documentación en `.env.example`. Pasan 4/4 subtests de `TestFIX136_UmbralesDeAlertaConfigurables` en `test/back/etapa3_fixes_pendientes_test.go`.

- **Área:** Backend / Infraestructura
- **Asignado:** Nico / Lucas / Tayra
- **Estimación:** 1.5 h
- **Depende de:** `RF-11`, `RNF-07`.
- **Entregable:** Struct `ConfigAlertas` inyectable en `NewEvaluadorAlertas`, variables `ALERT_THRESHOLD_*` y `ALERT_CYCLES_COUNT` en `main.go` y `.env.example`.
- **Criterio de éxito:** `TestFIX136_UmbralesDeAlertaConfigurables` pasa 4/4 en verde.

#### `FIX-143` - Vocabulario de auditoría con contrato cerrado, eliminación de códigos duplicados, normalización histórica y publicación de enums en Swagger (`RF-08` / `BAC-44` / `BAC-45`) (Backend / Base de Datos)

> [!NOTE]
> **Estado verificado (04/11/2026, backend `4011b08`): implementada y con tests en verde.** Commit `4011b08` (`fix(auditoria): vocabulario cerrado de accion y resultado, enums en Swagger y normalización histórica`). Se unificaron `auditoria.resultado` (`EXITO`, `FALLA`, `PENDIENTE`) y `auditoria.accion` mediante `AccionesAuditoriaValidas()` y `ResultadosAuditoriaValidos()` en `ports/audit_port.go`, se eliminaron las constantes muertas (`CREAR_VM`, `CREAR_LXC`, `SUSPENDER_VM`, `ACTUALIZAR_PERFIL`, `ResultadoSuccess`, `ResultadoFailed`), se agregaron validación de contrato en repositorio (`ErrAuditoriaFueraDeContrato`), enums en `docs/swagger.json`, tabla en `docs/contrato-etapa1.md`, `TestContratoVocabularioAuditoria` en `docs/contrato_etapa1_test.go` y la migración `migrations/20261104_vocabulario_auditoria.sql`. Pasan 4/4 subtests de `TestFIX143_VocabularioAuditoriaContratoCerrado` en `test/back/etapa3_fixes_pendientes_test.go` y la suite interna del backend.

- **Área:** Backend / Base de Datos
- **Asignado:** Lisandro / Lucas / Nico
- **Estimación:** 2.0 h
- **Depende de:** `RF-08`, `BAC-44`, `BAC-45`.
- **Entregable:** Vocabulario cerrado en `ports/audit_port.go`, validación estricta en `audit_handler.go` y `audit_repository.go`, enums en Swagger y `docs/contrato-etapa1.md`, test `TestContratoVocabularioAuditoria` y migración histórica `20261104_vocabulario_auditoria.sql`.
- **Criterio de éxito:** `TestFIX143_VocabularioAuditoriaContratoCerrado` (4/4) y `TestContratoVocabularioAuditoria` pasan al 100 % en verde.



### `FIX-143` - Vocabulario de auditoría con contrato cerrado, eliminación de códigos duplicados, normalización histórica y publicación de enums en Swagger (`RF-08` / `BAC-44` / `BAC-45`) (Backend / Base de Datos)

- **Área:** Backend / Base de Datos
- **Asignado:** Lucas / Nico
- **Estimación:** 2.0 h
- **Depende de:** `RF-08`, `BAC-44`, `BAC-45`.
- **Objetivo:**
  Dejar en un **único conjunto cerrado, tipado y documentado** en el Backend los códigos que se escriben y filtran en `auditoria.accion` y `auditoria.resultado`, publicando sus enums en Swagger (`docs/swagger.json`) y en `docs/contrato-etapa1.md`, validándolos con un test de contrato y normalizando las filas históricas de la base de datos.

- **Contexto y Problemática Técnica Detallada (Backend y Base de Datos):**
  Actualmente los campos `accion` y `resultado` se exponen tal cual en `GET /api/admin/audit` y `GET /api/admin/audit/export`, sin `enum` en Swagger ni tabla de contrato en `docs/contrato-etapa1.md`, y sin un test de contrato que fije el conjunto cerrado. Esto permitió que distintas capas del backend escribieran **dos familias incompatibles de códigos**:

  1. **Esquizofrenia en `auditoria.resultado` (dos familias activas en paralelo):**
     - En [`backend/internal/core/ports/audit_port.go`](file:///home/micaela/centinela/backend/internal/core/ports/audit_port.go#L73) y líneas [`130-146`](file:///home/micaela/centinela/backend/internal/core/ports/audit_port.go#L130-L146), el propio puerto declara dos vocabularios simultáneos:
       ```go
       Resultado string // Ciclo de vida: PENDIENTE/SUCCESS/FAILED; otras operaciones: EXITO/FALLA.
       ...
       ResultadoExito     = "EXITO"
       ResultadoFalla     = "FALLA"
       ResultadoPendiente = "PENDIENTE"
       ResultadoSuccess   = "SUCCESS"
       ResultadoFailed    = "FAILED"
       ```
     - **Servicios de Autenticación y Usuarios** (`auth_service.go:121, 146, 162, 232, 315, 360, 405, 444, 480, 520, 560` y `user_service.go:186, 297, 369, 467`) persisten `EXITO` y `FALLA`.
     - **Servicios y Handlers de Instancias, Snapshots, Tareas y Alertas** (`audit_service.go:285-330`, `instance_handler.go`, `configuracion_instancias.go`, `snapshot_service.go`, `seguimiento_tareas.go`, `alert_evaluator_worker.go`) persisten `PENDIENTE`, `SUCCESS` y `FAILED`.
     - **Evidencia en base de datos de pruebas:** `EXITO` (8.041 filas), `FALLA` (87 filas), `PENDIENTE` (45 filas), `SUCCESS` (39 filas), `PENDING` (33 filas), `FAILED` (8 filas).
     - En [`audit_handler.go:409-413`](file:///home/micaela/centinela/backend/internal/adapters/primary/http/audit_handler.go#L409-L413), `validarConsultaAuditoria` acepta la mezcla de 5 literales (`"SUCCESS", "FAILED", "PENDIENTE", "EXITO", "FALLA"`) y rechaza `"PENDING"` con `400 Bad Request`.
     - En [`audit_repository.go:267-269`](file:///home/micaela/centinela/backend/internal/adapters/secondary/postgres/audit_repository.go#L267-L269), `aplicarFiltros` filtra por igualdad exacta (`WHERE a.resultado = ?`) sin traducción ni agrupación.

  2. **Acciones duplicadas o muertas en `auditoria.accion`:**
     - Antecedente de regresión cruzada: el commit `86e3k5fqp` (05/10) había unificado a `START`/`STOP` + `EXITO`/`FALLA`, y el commit `d76eb35` de `RF-08` (10/10) lo revirtió parcialmente a `INICIAR_VM` + `SUCCESS`/`FAILED`, dejando en la base `START` (31) vs `INICIAR_VM` (47) y `STOP` (27) vs `DETENER_VM`.
     - En [`ports/audit_port.go:98-128`](file:///home/micaela/centinela/backend/internal/core/ports/audit_port.go#L98-L128) coexisten constantes en español (`CREAR_USUARIO`, `ACTUALIZAR_USUARIO`, `ELIMINAR_USUARIO`, `ASIGNAR_PERMISOS`, `INICIAR_VM`, `DETENER_VM`, `APAGAR_VM`, `REINICIAR_VM`, `ELIMINAR_VM`), constantes en inglés (`CREATE_INSTANCE`, `CONFIG_UPDATED`, `PASSWORD_CHANGED`, `SNAPSHOT_CREATE`, `SNAPSHOT_ROLLBACK`, `SNAPSHOT_DELETE`, `ALERT_TRIGGERED`, `ALERT_RESOLVED`, `INSTANCE_PERMISSIONS_PURGED`, `UPDATE_PROFILE_EMAIL_REQUESTED`, `UPDATE_PROFILE_EMAIL_CONFIRMED`) y **7 constantes declaradas en `audit_port.go` sin ningún escritor activo en servicios/handlers** (`AccionEventosConexion = "EVENTOS_CONEXION"`, `AccionEventosCierre = "EVENTOS_CIERRE"`, `AccionActualizarPerfil = "ACTUALIZAR_PERFIL"`, `AccionCrearVM = "CREAR_VM"`, `AccionCrearLXC = "CREAR_LXC"`, `AccionSuspenderVM = "SUSPENDER_VM"`, `AccionResumeVM = "REANUDAR_VM"`).

- **Especificación Técnica (Qué debe hacer el equipo de Backend):**
  1. **Contrato cerrado en `internal/core/ports/audit_port.go` y servicios:**
     - Definir un **único conjunto cerrado para `resultado`** (se recomienda mirar el volumen persistido y unificar en `EXITO`, `FALLA`, `PENDIENTE` —8.173 filas vs 80 filas—, o el conjunto único que el equipo acuerde, eliminando la familia paralela).
     - Definir un **único conjunto cerrado para `accion`** en `ports/audit_port.go`, eliminando las constantes muertas o duplicadas sin escritor activo (`CREAR_VM`, `CREAR_LXC`, `SUSPENDER_VM`, `REANUDAR_VM`, `ACTUALIZAR_PERFIL`, `EVENTOS_CONEXION`, `EVENTOS_CIERRE`) y exponiendo funciones/slices canónicos (`AccionesAuditoriaValidas() []string` y `ResultadosAuditoriaValidos() []string`).
     - Reemplazar todos los usos de la familia descartada en `audit_service.go`, `instance_handler.go`, `configuracion_instancias.go`, `snapshot_service.go`, `seguimiento_tareas.go`, `alert_evaluator_worker.go`, `auth_service.go` y `user_service.go` para que **ningún escritor pueda emitir un valor fuera del conjunto cerrado** ni existan literales sueltos fuera de las constantes de `ports`.
  2. **Validación HTTP y publicación de Enums en Swagger y Contrato:**
     - En [`audit_handler.go:409-413`](file:///home/micaela/centinela/backend/internal/adapters/primary/http/audit_handler.go#L409-L413) (`validarConsultaAuditoria`), validar `f.Resultado` (y `f.Accion` cuando se envíe) estrictamente contra el conjunto cerrado de `ports`, devolviendo `400 Bad Request` (`INVALID_PARAMETERS`) para cualquier valor fuera del enum.
     - En `audit_handler.go` / `responses.go`, declarar las anotaciones Swagger `enums(...)` tanto en los `@Param accion` y `@Param resultado` de `GET /admin/audit` y `GET /admin/audit/export` como en los campos `accion` y `resultado` de `ItemAuditoriaResponse` / `ports.AuditoriaDTO`, y regenerar `backend/docs/swagger.json`.
     - Documentar la tabla oficial del vocabulario cerrado de `auditoria.accion` y `auditoria.resultado` en [`backend/docs/contrato-etapa1.md`](file:///home/micaela/centinela/backend/docs/contrato-etapa1.md).
     - En [`backend/docs/contrato_etapa1_test.go`](file:///home/micaela/centinela/backend/docs/contrato_etapa1_test.go), agregar una aserción con `verificarEnum` sobre `swagger.json` y `contrato-etapa1.md` que falle si se agrega, renombra o quita un código en `ports` sin actualizar Swagger y el documento de contrato.
  3. **Normalización de filas históricas en Base de Datos:**
     - Dado que la tabla `auditoria` es particionada y *append-only* (protegida por el trigger `trg_auditoria_inmutable`), crear una migración SQL controlada en `backend/migrations/` (para ejecutar con backup previo y ventana) que dentro de una transacción atómica deshabilite temporalmente el trigger de inmutabilidad, normalice las filas históricas al conjunto cerrado elegido (`SUCCESS -> EXITO`, `FAILED -> FALLA`, `PENDING -> PENDIENTE`, `START -> INICIAR_VM`, `STOP -> DETENER_VM`, etc.) y rehabilite el trigger antes del `COMMIT`.
  4. **Fuera de alcance:**
     - `detalles.accion` del evento SSE (`/api/events`), que pertenece al contrato de seguimiento de tareas en tiempo real.

- **Criterios de Aceptación (DoD — Backend):**
  1. Toda fila nueva en `auditoria` usa únicamente códigos del conjunto cerrado de `ports`; ningún escritor emite valores de otra familia ni literales sueltos.
  2. El enum de `accion` y `resultado` está publicado en `backend/docs/swagger.json` y documentado como tabla en `backend/docs/contrato-etapa1.md`.
  3. El test de contrato `backend/docs/contrato_etapa1_test.go` y la prueba `TestFIX143_VocabularioAuditoriaContratoCerrado` en [`test/back/etapa3_fixes_pendientes_test.go`](file:///home/micaela/centinela/test/back/etapa3_fixes_pendientes_test.go#L345-L461) pasan en verde (y fallan si se agrega, cambia o quita un código sin documentarlo o si coexisten dos familias de `resultado`).
  4. El filtro `resultado` en `GET /api/admin/audit` y `GET /api/admin/audit/export` acepta exactamente el conjunto cerrado y devuelve `400 Bad Request` para cualquier otro valor.
  5. Aviso al equipo de frontend cuando el conjunto quede fijo y desplegado para ejecutar `FIX-144`.

#### `FIX-137` - Saneamiento de código inerte: eliminación de stub documental `PausarInstancia` y simplificación de residuos multi-tenancy (`RF-12` / `RNF-05`) (Backend)

> [!NOTE]
> **Estado verificado (05/11/2026, backend `d0f3f57`): implementada y con tests en verde.** Commit `5b3ac80` (`fix: elimina endpoints fantasma e incluye auditorías globales`). Se eliminó el stub `PausarInstancia` y la ruta `/pause` de Swagger, se actualizó `audit_repository.go` para incluir eventos globales del sistema (`a.organizacion_id IS NULL`) y se limpiaron los tipos DTO huérfanos en `contrato_etapa1.go`. Pasa `TestFIX137_CodigoInerteYAuditoriaDeSistema` (3/3 subtests ✅) en `test/back/etapa3_fixes_pendientes_test.go`.

- **Área:** Backend
- **Asignado:** Lucas / Micaela
- **Estimación:** 1.5 h
- **Depende de:** `RF-12`, `RNF-05`.
- **Entregable:**
  1. Eliminación del stub `PausarInstancia` y sus anotaciones Swagger en `instance_handler.go` y `docs/swagger.json`.
  2. Inclusión de eventos del sistema (`a.organizacion_id IS NULL`) en `audit_repository.go`.
  3. Eliminación de tipos DTO huérfanos en `contrato_etapa1.go`.
- **Criterio de éxito:** `TestFIX137_CodigoInerteYAuditoriaDeSistema` (3/3) en verde.

#### `FIX-139` - Reconciliación inmediata de tareas al inicio del backend para mitigar falsos positivos de `409 INSTANCE_BUSY` (`RF-04` / `RNF-04`) (Backend)

> [!NOTE]
> **Estado verificado (05/11/2026, backend `d0f3f57`): implementada y con tests en verde.** Commit `d92a568` (`fix: reconcilia tareas antes de abrir la API`). Se incorporó `RecuperarYReconciliarAlArranque` en `seguimiento_tareas.go`, se invoca en `cmd/api/main.go` antes de iniciar el servidor HTTP y se agregaron pruebas unitarias del arranque con tareas `RUNNING`. Pasa `TestFIX139_ReconciliacionDeTareasAlArranque` (3/3 subtests ✅) en `test/back/etapa3_fixes_pendientes_test.go`.

- **Área:** Backend
- **Asignado:** Lucas / Nico
- **Estimación:** 1.5 h
- **Depende de:** `RF-04`, `RNF-04`.
- **Entregable:**
  1. Método `RecuperarYReconciliarAlArranque(ctx context.Context)` en `seguimiento_tareas.go`.
  2. Invocación en `cmd/api/main.go` antes de `srv.ListenAndServe()`.
  3. Pruebas unitarias en `seguimiento_tareas_test.go`.
- **Criterio de éxito:** `TestFIX139_ReconciliacionDeTareasAlArranque` (3/3) en verde.

#### `FIX-144` - Alineación de filtros, contadores y diccionario de acciones de auditoría con el enum cerrado del Backend y retiro de lectura doble (`RF-08` / `FRN-31` / `FRN-12`) (Frontend)

> [!NOTE]
> **Estado verificado (05/11/2026, frontend `bed77ff`): implementada y con tests en verde.** Commits `a2a2761`, `798a93f`, `7ca0558`, `af78395`, `efd59a5` y `c110791` (PR #192). Se sincronizó `AUDIT_ACTION_LABELS` (31/31 códigos del backend) y `AUDIT_ACTION_OPTIONS` en `auditActionLabels.ts`, y se unificaron `Auditoria.tsx`, `auditFilters.ts`, `auditService.ts` y `deletedUsersService.ts` con el vocabulario cerrado `EXITO`, `FALLA`, `PENDIENTE`, retirando la lectura doble temporal. Pasan el bloque `FIX-144` en `test/front/etapa3-pendientes.test.tsx` (2/2 ✅), `test/front/frn-31-audit-etapa3.test.tsx` (4/4 ✅) y `frontend/centinela/src/test/audit-action-labels.test.ts` (3/3 ✅).

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1.0 h
- **Depende de:** `FIX-143`, `FRN-31`, `FRN-12`.
- **Entregable:**
  1. Alineación 1 a 1 de `auditActionLabels.ts` con el enum cerrado de `auditoria.accion` del backend.
  2. Unificación de `Auditoria.tsx` y `deletedUsersService.ts` sobre el enum cerrado `EXITO`, `FALLA`, `PENDIENTE`.
- **Criterio de éxito:** Suite `FIX-144` en `test/front/etapa3-pendientes.test.tsx` (2/2) en verde.


