# Informe de Estado de Tareas Verificado por Pruebas

**Fecha de ejecución:** 24/10/2026 (tras actualizar los submódulos a los últimos commits de `origin/main`)  
**Alcance:** Evaluación integral de las **18 tareas activas** en [documentacion/actual.md](../documentacion/actual.md) (1 de Etapa 1, 9 de Etapa 2 y 8 de Etapa 3), además de la verificación de regresión de las **20 tareas** referenciadas en el encabezado de `actual.md`. Todas las pruebas fueron ejecutadas contra las suites completas en `test/front` (Vitest + jsdom) y `test/back` (Go).

| Submódulo | Revisión probada (`origin/main`) | Últimos commits incorporados |
|---|---|---|
| **Backend** | `c08af54` (`fix: purgar permisos al destruir instancias`) | `c08af54` (**`BAC-24C` / `FIX-113`**: purga transaccional de permisos al destruir instancias, auditoría `INSTANCE_PERMISSIONS_PURGED` y notificación SSE `TASK_FINISHED` a `UsuariosAfectados`), `1801e33` (**`FIX-118`**: eliminación de credenciales por defecto en sembrado de administrador), `a25d1b3` (**`FIX-116`**: validación contra inyección de directivas en creación de instancias), `94e95a2` (**`FIX-117A`**: segregación de telemetría del host por rol para operadores), `5c72bbb` (**`FIX-119`**: validación de existencia de VMIDs y bloqueo de instancias protegidas en asignación de permisos). |
| **Frontend** | `d4e99b2` (`Merge branch 'fix/tests-verdes': suite del frente en verde`) | `d4e99b2` / `019d2de` / `b61d117` / `e78736c` (ajustes de pruebas internas del submódulo frontend), `59e5661` / `02cf6d8` (**`FRN-27` / `FRN-28` / `FRN-26`**: modal `EditInstanceConfigModal`, hook `useEditInstanceConfig`, servicio `updateInstanceConfiguration` y pestañas de `InstanceDetailPage`), `5c993ee` (**`FIX-117B`**: adaptación del Dashboard para operadores ante segregación de métricas globales). |

---

## 1. Resumen Ejecutivo de las 18 Tareas Activas de `documentacion/actual.md`

### 1.1 Implementadas y pasan los tests al 100% (4 tareas)
1. **`BAC-24C` / `FIX-113`** — Limpieza de permisos al eliminar una instancia (`Backend`):
   - **Evidencia:** Commit `c08af54` en `internal/core/services/seguimiento_tareas.go`. Captura `UsuariosAfectados` antes de eliminar filas de `permisos_instancia`, registra `INSTANCE_PERMISSIONS_PURGED` en auditoría y despacha `TASK_FINISHED` mediante `PublicarParaUsuarios`.
   - **Prueba:** `TestFIX113_PurgaDePermisosAlEliminarInstancia` en `test/back/etapa2_seguridad_acceptance_test.go` — **PASS**.
2. **`FIX-117B`** — Adaptación del Dashboard de operador ante segregación de métricas globales (`Frontend`):
   - **Evidencia:** Commit `5c993ee` en `Home.tsx` y `useNodeStatus.ts`. Oculta el bloque de almacenamiento global cuando `storage` se omite para el rol `OPERATOR` y renderiza el resumen acotado de instancias sin errores de consola.
   - **Prueba:** `test/front/fix-117b-dashboard-operador.test.tsx` (4/4 pruebas) — **PASS**.
3. **`FRN-27`** — Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug (`Frontend`):
   - **Evidencia:** Commit `02cf6d8` (`EditInstanceConfigModal.tsx`, `InstanceResourcesPanel.tsx`, `useEditInstanceConfig.ts`). Precarga `cores` y `ramMb`, calcula $\Delta\text{RAM}$ reactivamente contra `allocatable.freeRamMb`, bloquea excesos de cuota, oculta/deshabilita edición para `READ_ONLY` y muestra la advertencia de hot-plug cuando la instancia está encendida.
   - **Prueba:** `test/front/frn-27-28-edit-instance-config.test.tsx` (bloque `FRN-27`, 3/3 pruebas) — **PASS**.
4. **`FRN-28`** — Integración de edición en la página de Detalle de Instancia y refresco reactivo (`Frontend`):
   - **Evidencia:** Commit `02cf6d8` (`instanceDetailService.ts`, `useInstanceDetail.ts`). Envía `PUT /api/instances/:instanceRef/config`, distingue `requiresReboot: false` (toast verde y actualización reactiva de `vCPU` y `Memoria RAM` en el detalle y tarjeta `Hardware`) frente a `requiresReboot: true` (toast amarillo y banner `role="status"` de reinicio pendiente), y maneja errores `409 RESOURCE_QUOTA_EXCEEDED` y `409 INSTANCE_BUSY`.
   - **Prueba:** `test/front/frn-27-28-edit-instance-config.test.tsx` (bloque `FRN-28`, 3/3 pruebas) — **PASS**.

---

### 1.2 Implementadas (total o parcialmente) pero NO pasan los tests (8 tareas)
1. **`FIX-108`** — Adaptación del Frontend para consumo de identificadores de instancia ofuscados (`Frontend`):
   - **Estado de implementación:** Parcialmente implementada. `InventoryInstance.id` acepta `string | number` y las vistas operan con `inst_...` (5/11 pruebas pasan en `fix-108-ids-ofuscados.test.tsx`).
   - **Por qué falla el test (6/11 fallos en `test/front/fix-108-ids-ofuscados.test.tsx`):** `validateInstanceRef` en `frontend/centinela/src/components/features/instances/services/instanceService.ts:12-16` únicamente comprueba `typeof ref === 'string' && ref.trim()` sin validar el formato `/^[a-zA-Z0-9_-]+$/`, permitiendo pasar identificadores maliciosos con path traversal (`../admin/users`), espacios (`inst 101`), query strings (`inst_101?admin=true`), barras (`inst/101`), punto y coma (`inst_101;rm`) o saltos de línea.
2. **`FRN-23`** — Componentes de gráficos interactivos de rendimiento y selector de rango temporal (`Frontend`):
   - **Estado de implementación:** Parcialmente implementada en `InstanceResourceUsage.tsx` (renderiza selector `2h`/`24h` y consulta `GET /instances/:instanceRef/metrics?timeframe=2h`; pasa 1/9 pruebas en `frn-23-metrics-charts.test.tsx`).
   - **Por qué falla el test (8/9 fallos en `test/front/frn-23-metrics-charts.test.tsx`):**
     - Solo existen 2 gráficos (`CPU` y `Memoria`); faltan los gráficos de **Tráfico de Red (KB/s In/Out)** y **Disco I/O (KB/s Read/Write)** exigidos en `actual.md:118-119`.
     - Falta la opción **"En vivo"** en el selector de rango temporal conectada al WebSocket (`FRN-22`).
     - No se muestra el indicador **"Instancia detenida"** cuando `cpu` y `mem` son `null`.
3. **`FRN-26`** — Envío de creación, manejo de errores de cuota (`409`) y seguimiento en la UI (`Frontend`):
   - **Estado de implementación:** Parcialmente implementada en `createNewInstances.tsx` y `useInstances.ts` (despacha `POST /api/instances`, deshabilita el botón mostrando `"Creando instancia..."`, navega a `/instances` tras `202 Accepted` y procesa `TASK_FINISHED`; pasa 1/4 pruebas en `frn-26-create-submit-quota-sse.test.tsx`).
   - **Por qué falla el test (3/4 fallos en `test/front/frn-26-create-submit-quota-sse.test.tsx`):**
     - Ante `409 Conflict (RESOURCE_QUOTA_EXCEEDED)`, `createNewInstances.tsx` solo emite un toast flotante global y **no renderiza un banner de alerta en línea (`role="alert"`) dentro del wizard** como exige `actual.md:139`.
     - En `/instances`, cuando una instancia recién creada tiene `activeTask: { action: 'CREATE', status: 'RUNNING' }`, `Instances.tsx` muestra el texto genérico `"Ejecutando..."` en el botón de acciones pero **no muestra el estado transitorio `"Creando..."`** en la fila (`actual.md:140`).
4. **`FRN-36` (`BRG-10-FRN`)** — Modularización de `InstanceDetailPage` para Pestaña de Snapshots y Estados Transitorios (`Frontend`):
   - **Estado de implementación:** Parcialmente implementada en `InstanceDetailContent.tsx` (la pestaña `"Snapshots"` ya existe y es seleccionable; pasa 1/3 pruebas del bloque `FRN-36` en `frn-36-29-snapshots-ui.test.tsx`).
   - **Por qué falla el test (2/3 fallos del bloque `FRN-36` en `test/front/frn-36-29-snapshots-ui.test.tsx`):**
     - `InstanceStatusBadge` (`InstanceBadges.tsx:21-40`) solo contempla `start` (`"Encendiendo…"`) y `stop` (`"Apagando…"`), ignorando `SNAPSHOTTING` (`"Creando snapshot..."`) y `ROLLING_BACK` (`"Revirtiendo..."`).
     - En `InstanceDetailContent.tsx:65`, `<InstanceAction variant="header" disabled={isSavingResources} ... />` no recibe `isBusy`, y `InstanceAction.tsx:206-219` no incluye `isTransitioning` en la propiedad `disabled` de los botones `Encender`, `Apagar` y `Reiniciar` cuando `variant === 'header'`.
5. **`BAC-44`** — Interceptor global y consultas complejas sobre auditoría particionada (`Backend`):
   - **Estado de implementación:** Parcialmente implementada en `audit_handler.go` y `audit_repository.go` (soporta `desde`, `hasta`, `accion`, `usuarioId`, `resultado`, validación `desde > hasta` y orden `fecha_hora DESC`).
   - **Por qué falla el test (`TestEtapa3BAC44_InterceptorGlobalYConsultasAuditoria` en `test/back/etapa3_acceptance_test.go`):**
     - `internal/core/ports/audit_port.go` aún no declara las constantes de acción de Etapa 3: `SNAPSHOT_CREATE`, `SNAPSHOT_ROLLBACK`, `SNAPSHOT_DELETE` ni `CONFIG_UPDATED` (actualmente usa `CONFIG_INSTANCE`).
     - `audit_handler.go` no lee los query params `recursoTipo`, `recursoId`, `page` ni `limit` especificados en `actual.md:239` (solo lee `pagina` y `tamano`).
6. **`BAC-45`** — Streaming de exportación estructurada de auditoría a CSV y JSON (`Backend`):
   - **Estado de implementación:** Parcialmente implementada en `audit_handler.go`, `audit_service.go` y `audit_repository.go` (existe `GET /api/admin/audit/export` para CSV y JSON con sanitización de fórmulas `=`, `+`, `-`, `@`).
   - **Por qué falla el test (`TestEtapa3BAC45_StreamingExportacionAuditoriaCSVyJSON` en `test/back/etapa3_acceptance_test.go`):**
     - `audit_repository.go` (`ListarParaExportacion`) carga hasta 100.000 registros en memoria con `.Limit(limiteExportacion).Find(& filas)` en lugar de transmitir mediante cursor SQL (`.Rows()`, `ScanRows` o `c.Stream`) como exige `actual.md:252`.
     - `audit_handler.go` genera `filename="auditoria.csv"` / `"auditoria.json"` en lugar del formato `filename="auditoria_centinela_YYYYMMDD.csv|json"` exigido en `actual.md:253`, y lee `formato` pero no el alias `format` de `actual.md:251`.
7. **`FRN-31`** — Vista administrativa de Auditoría con filtros reactivos completos (`Frontend`):
   - **Estado de implementación:** Parcialmente implementada en `Auditoria.tsx` (renderiza tabla paginada, badges semánticos verde/amarillo/rojo para `SUCCESS`/`WARNING`/`FAILED` y exportación CSV; pasa 1/4 pruebas en `frn-31-audit-etapa3.test.tsx`).
   - **Por qué falla el test (3/4 fallos en `test/front/frn-31-audit-etapa3.test.tsx`):**
     - El filtro de fecha es un único botón disparador de calendario sin controles con `aria-label="Desde"` y `aria-label="Hasta"` (`actual.md:309`).
     - Faltan el selector de filtro por `recursoTipo` y las opciones de acción de Etapa 3 (`SNAPSHOT_CREATE`, `SNAPSHOT_ROLLBACK`, `SNAPSHOT_DELETE`, `ALERT_TRIGGERED`, `CONFIG_UPDATED`).
     - Al pulsar `"Ver detalles"` no se abre un modal (`role="dialog"`) con el payload JSON formateado (`actual.md:313`).
     - Falta el botón/opción para exportar en formato **JSON** (`actual.md:251`).

8. **`INF-11`** — Configuración de producción en Nginx: TLS 1.3, HSTS, CSP y WSS (`Infraestructura`):
   - **Estado de implementación:** Parcialmente implementada en `docker/nginx-edge.conf` (ya redirige HTTP 80 a HTTPS 443 con `301`, habilita `TLSv1.2 TLSv1.3`, reenvía `Upgrade` para WebSockets en `/api/ws/` y deshabilita `proxy_buffering` en `/api/events`; pasa 2/3 subtests en `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS`).
   - **Por qué falla el test (1/3 subtests en `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS`):** Faltan las cabeceras de endurecimiento `Strict-Transport-Security` (con `includeSubDomains`), `X-Frame-Options`, `X-Content-Type-Options` y `Content-Security-Policy` (con `default-src 'self'` y `wss:`) en `docker/nginx-edge.conf`.

---

### 1.3 No implementadas (6 tareas)
1. **`FRN-22`** — Hook `useInstanceMetrics` y cliente WebSocket para telemetría en tiempo real (`Frontend`):
   - **Falta:** No existen `getWsTicket()` (`POST /api/auth/ws-ticket`) ni el hook `useInstanceMetrics.ts`. Fallan 9/9 pruebas en `test/front/frn-22-use-instance-metrics.test.tsx`.
2. **`INF-13` (`BRG-11-INF`)** — Verificación y Validación de Storage Compatible con Snapshots en Proxmox Real (`Infraestructura` / `Documentación`):
   - **Falta:** `documentacion/api-proxmox.md` no documenta aún los tipos de storage compatibles con snapshots (`lvmthin`, `zfspool`/`zfs`, `qcow2`) frente a los incompatibles (`raw`), y mantiene un UUID en texto plano en `PROXMOX_TOKEN_SECRET=` (línea 150). Falla `TestEtapa2INF13_StorageCompatibleConSnapshots` en `test/back/etapa3_acceptance_test.go`.
3. **`BAC-42`** — Adaptador de Proxmox para Snapshots, endpoints de Creación y Listado (`Backend`):
   - **Falta:** No existen `ListarSnapshots`, `CrearSnapshot`, `RevertirSnapshot` ni `EliminarSnapshot` en `proxmox_port.go` / `client.go`, ni las rutas `GET /api/instances/:instanceRef/snapshots` y `POST /api/instances/:instanceRef/snapshots`. Falla `TestEtapa3BAC42_AdaptadorYEndpointsSnapshotsCreacionYListado` en `test/back/etapa3_acceptance_test.go`.
4. **`BAC-43`** — Reversión (Rollback) y Eliminación de Snapshots con bloqueo de concurrencia (`Backend`):
   - **Falta:** No existen las rutas `POST /api/instances/:instanceRef/snapshots/:snapname/rollback` ni `DELETE /api/instances/:instanceRef/snapshots/:snapname`. Falla `TestEtapa3BAC43_RollbackYEliminacionSnapshotsConcurrencia` en `test/back/etapa3_acceptance_test.go`.
5. **`BAC-46`** — Background worker evaluador de umbrales con histéresis y deduplicación (`Backend`):
   - **Falta:** No existe el worker evaluador de umbrales (`ALERT_CPU_THRESHOLD`, `ALERT_RAM_THRESHOLD`, `ALERT_STORAGE_THRESHOLD`) con histéresis en Redis ni sus pruebas unitarias. Falla `TestEtapa3BAC46_WorkerEvaluadorUmbralesHisteresis` en `test/back/etapa3_acceptance_test.go`.
6. **`FRN-29`** — Pestaña "Snapshots" en Detalle de Instancia y formulario de nuevo snapshot (`Frontend`):
   - **Falta:** La pestaña `"Snapshots"` en `InstanceDetailContent.tsx:123-128` solo muestra el texto estático `"La consulta de snapshots todavía no está disponible."` sin tabla de snapshots, ni modal `"Crear Snapshot"` / `"Tomar Snapshot"`, ni acciones de Rollback/Eliminación. Fallan las 4 pruebas de `FRN-29` en `test/front/frn-36-29-snapshots-ui.test.tsx`.

---

## 2. Estado de las 20 Tareas Referenciadas en el Encabezado de `actual.md`

| Tarea | Área | Estado en Suites (`test/back` y `test/front`) | Detalle |
|---|---|:---:|---|
| `FIX-104B` / `FIX-105` | Frontend | 🟢 **PASS** | `fix-104b-105-tareaid.test.tsx` (12/12 pruebas en verde). |
| `FIX-100` | Frontend | 🟢 **PASS** | `dashboard-node.test.tsx` (12/12 pruebas en verde). |
| `FIX-103` | Frontend | 🟢 **PASS** | `fix-103-node-history.test.tsx` (7/7 pruebas en verde). |
| `FIX-104` / `FRN-24` / `FRN-25` | Frontend | 🟢 **PASS** | `frn-24-25-create-wizard.test.tsx` (11/11 pruebas en verde). |
| `FIX-111B` | Frontend | 🟢 **PASS** | `fix-111b-cambio-email-reauth.test.tsx` (4/4 pruebas en verde). |
| `FIX-115B` | Frontend | 🟢 **PASS** | `fix-115b-contrasenas-largas.test.tsx` (11/11 pruebas en verde). |
| `BAC-36` | Backend | 🟢 **PASS** | `TestEtapa2BAC36_MetricasHistoricasEndpoint` en verde. |
| `FIX-107` | Backend | 🟢 **PASS** | `TestFIX107_OfuscacionDeVMIDs` en verde. |
| `FIX-109` | Backend | 🟢 **PASS** | Implementado en `node_handler.go` y `instance_handler.go`. |
| `FIX-110` | Backend | 🟢 **PASS** | `TestFIX110_SaneamientoDeErrores` en verde. |
| `FIX-111A` | Backend | 🟢 **PASS** | Implementado en `auth_service.go` / `user_handler.go`. |
| `FIX-114` | Backend | 🟢 **PASS** | Implementado en `auth_service.go` (`RevocarTodasLasSesiones`). |
| `FIX-117A` | Backend | 🟢 **PASS** | Commit `94e95a2` en `node_handler.go`. |
| `FIX-118` | Backend | 🟢 **PASS** | `TestFIX118_SinCredencialesPorDefecto` en verde (commit `1801e33`). |
| `FIX-119` | Backend | 🟢 **PASS** | Commit `5c72bbb` en `user_service.go`. |
| `FIX-112` | Backend | 🟠 **Fallo Parcial** | `TestFIX112_RateLimitingAutenticacion`: falta bloqueo pesimista `SELECT ... FOR UPDATE` al leer `IntentosRecuperacion` en `ConfirmarRecuperacionContrasena`. |
| `FIX-115A` | Backend | 🟠 **Fallo Parcial** | `TestFIX115A_ContrasenasDeHasta128Caracteres`: `auth_handler.go` y `user_port.go` aún conservan un tag `max=12` en un DTO secundario y `password.go` no verifica explícitamente la cota superior de 128 en `ValidarComplejidad`. |
| `FIX-116` | Backend | 🟠 **Fallo Parcial** | `TestFIX116_InyeccionDeParametrosProxmox`: la validación contra comas se agregó en `instance_service.go`, pero el test exige defensa en profundidad también dentro de `internal/adapters/secondary/proxmox/client.go` antes de concatenar `net0`. |
| `FIX-101` | Frontend | 🟠 **Fallo Parcial** | `instances-ola2.test.tsx` (5/6 pasan; falla 1 test porque el buscador en memoria filtra por nombre/IP pero no por `String(instance.id)`). |
| `FIX-102` | Frontend | 🟠 **Fallo Parcial** | `instances-modals.test.tsx` (5/6 pasan; falla 1 test de visibilidad de acciones para `OPERATOR`). |
| `FIX-91` | Frontend | 🔴 **No Implementada** | `fixes-pending.test.tsx` falla únicamente en `FIX-91`: el archivo muerto `frontend/centinela/src/components/features/instances/components/InstanceTester.tsx` aún no fue eliminado. |

---

## 3. Resultados Globales de Ejecución de Suites

### 3.1 Suite de Backend (`test/back`)
- **Compilación (`go test -c ./...`):** Exitosa (0 errores).
- **Ejecución (`go test -short -v ./...`):**
  - **En verde (`PASS`):** `TestINF09WebSocketsNginxEdge`, `TestINF07CPermisosTokenProxmox`, `TestEtapa1Ola3BAC32`, `TestEtapa2BAC34_AdaptadorMetricasRRD`, `TestEtapa2BAC36_MetricasHistoricasEndpoint`, `TestEtapa2BAC39_SeguimientoCreacionYPermisos`, `TestEtapa2BAC40_EndpointConfigRecursos`, `TestEtapa2BAC41_AuditoriaAprovisionamientoYConfig`, `TestEtapa2BAC51_GuardaConcurrenciaInstanceBusy`, `TestEtapa2BAC52_BrokerSSEEventosAlertaGlobal`, `TestEtapa2BAC53_ParticionadoPreventivoEIndexacionAuditoria`, `TestFIX107_OfuscacionDeVMIDs`, `TestFIX110_SaneamientoDeErrores`, `TestFIX113_PurgaDePermisosAlEliminarInstancia`, `TestFIX118_SinCredencialesPorDefecto`, `TestFIX81_OrdenMigracionesBaseLimpia`, `TestEtapa1SimuladorProxmox`, `TestEtapa2BAC50_SimuladorSnapshotsYRollback`.
  - **En rojo (`FAIL` — detectando trabajo pendiente real):** `TestFIX112_RateLimitingAutenticacion`, `TestFIX115A_ContrasenasDeHasta128Caracteres`, `TestFIX116_InyeccionDeParametrosProxmox`, `TestEtapa2INF13_StorageCompatibleConSnapshots`, `TestEtapa3BAC42_AdaptadorYEndpointsSnapshotsCreacionYListado`, `TestEtapa3BAC43_RollbackYEliminacionSnapshotsConcurrencia`, `TestEtapa3BAC44_InterceptorGlobalYConsultasAuditoria`, `TestEtapa3BAC45_StreamingExportacionAuditoriaCSVyJSON`, `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS`, `TestEtapa3BAC46_WorkerEvaluadorUmbralesHisteresis`.

### 3.2 Suite de Frontend (`test/front`)
- **Comando:** `pnpm --dir test/front exec vitest run`
- **Archivos de prueba:** **42 archivos** (**27 en verde**, **13 con fallos**, **2 omitidos** por requerir entorno Docker E2E activo).
- **Pruebas individuales:** **287 pruebas** (**225 en verde**, **46 fallidas**, **16 omitidas**).

---

## 4. Matriz Completa de las 18 Tareas Activas en `documentacion/actual.md`

| # | Tarea | Etapa | Área | Responsable | Clasificación | Archivo de Prueba y Resultado |
|---|---|---|---|:---:|:---:|---|
| 1 | `BAC-24C` / `FIX-113` | Etapa 1 | Backend | Tayra / Lisandro | 🟢 **Implementada y pasa tests** | `test/back/etapa2_seguridad_acceptance_test.go` (`TestFIX113_PurgaDePermisosAlEliminarInstancia`) — **PASS** |
| 2 | `FIX-108` | Etapa 2 | Frontend | Belinda / Cristian | 🟠 **Implementada parcial (no pasa tests)** | `test/front/fix-108-ids-ofuscados.test.tsx` — **5/11 PASS, 6/11 FAIL** (falta sanitización regex en `validateInstanceRef`) |
| 3 | `FIX-117B` | Etapa 2 | Frontend | Belinda / Cristian | 🟢 **Implementada y pasa tests** | `test/front/fix-117b-dashboard-operador.test.tsx` — **4/4 PASS** |
| 4 | `FRN-22` | Etapa 2 | Frontend | Belinda | 🔴 **No implementada** | `test/front/frn-22-use-instance-metrics.test.tsx` — **0/9 PASS** (falta `useInstanceMetrics.ts` y `getWsTicket`) |
| 5 | `FRN-23` | Etapa 2 | Frontend | Luz | 🟠 **Implementada parcial (no pasa tests)** | `test/front/frn-23-metrics-charts.test.tsx` — **1/9 PASS, 8/9 FAIL** (faltan gráficos Red/Disco, modo "En vivo" y cartel "Instancia detenida") |
| 6 | `FRN-26` | Etapa 2 | Frontend | Cristian | 🟠 **Implementada parcial (no pasa tests)** | `test/front/frn-26-create-submit-quota-sse.test.tsx` — **1/4 PASS, 3/4 FAIL** (falta alerta `role="alert"` interna en 409 y etiqueta `"Creando..."` en tabla) |
| 7 | `FRN-27` | Etapa 2 | Frontend | Cristian | 🟢 **Implementada y pasa tests** | `test/front/frn-27-28-edit-instance-config.test.tsx` — **3/3 PASS** |
| 8 | `FRN-28` | Etapa 2 | Frontend | Cristian | 🟢 **Implementada y pasa tests** | `test/front/frn-27-28-edit-instance-config.test.tsx` — **3/3 PASS** |
| 9 | `FRN-36` (`BRG-10-FRN`) | Etapa 2 | Frontend | Cristian | 🟠 **Implementada parcial (no pasa tests)** | `test/front/frn-36-29-snapshots-ui.test.tsx` — **1/3 PASS, 2/3 FAIL** (pestaña creada; faltan badges `SNAPSHOTTING`/`ROLLING_BACK` y bloqueo de botones en cabecera) |
| 10 | `INF-13` (`BRG-11-INF`) | Etapa 2 | Infra / Docs | Nico | 🔴 **No implementada** | `test/back/etapa3_acceptance_test.go` (`TestEtapa2INF13_StorageCompatibleConSnapshots`) — **FAIL** (falta documentar storages compatibles en `api-proxmox.md` y limpiar UUID en `PROXMOX_TOKEN_SECRET`) |
| 11 | `BAC-42` | Etapa 3 | Backend | Lisandro | 🔴 **No implementada** | `test/back/etapa3_acceptance_test.go` (`TestEtapa3BAC42_AdaptadorYEndpointsSnapshotsCreacionYListado`) — **FAIL** (faltan métodos Proxmox y endpoints GET/POST de snapshots) |
| 12 | `BAC-43` | Etapa 3 | Backend | Tayra | 🔴 **No implementada** | `test/back/etapa3_acceptance_test.go` (`TestEtapa3BAC43_RollbackYEliminacionSnapshotsConcurrencia`) — **FAIL** (faltan endpoints de rollback y delete de snapshots) |
| 13 | `BAC-44` | Etapa 3 | Backend | Tayra | 🟠 **Implementada parcial (no pasa tests)** | `test/back/etapa3_acceptance_test.go` (`TestEtapa3BAC44_InterceptorGlobalYConsultasAuditoria`) — **FAIL** (faltan constantes Etapa 3 en `audit_port.go` y filtros `recursoTipo`, `recursoId`, `page`, `limit` en `audit_handler.go`) |
| 14 | `BAC-45` | Etapa 3 | Backend | Lisandro | 🟠 **Implementada parcial (no pasa tests)** | `test/back/etapa3_acceptance_test.go` (`TestEtapa3BAC45_StreamingExportacionAuditoriaCSVyJSON`) — **FAIL** (usa `.Find(&filas)` en memoria en vez de cursor `.Rows()` y nombra `auditoria.csv` en vez de `auditoria_centinela_YYYYMMDD.csv`) |
| 15 | `INF-11` | Etapa 3 | Infra | Nico | 🟠 **Implementada parcial (no pasa tests)** | `test/back/etapa3_acceptance_test.go` (`TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS`) — **2/3 subtests PASS, 1/3 FAIL** (faltan cabeceras `Strict-Transport-Security`, `Content-Security-Policy`, `X-Frame-Options` y `X-Content-Type-Options` en `docker/nginx-edge.conf`) |
| 16 | `BAC-46` | Etapa 3 | Backend | Lisandro | 🔴 **No implementada** | `test/back/etapa3_acceptance_test.go` (`TestEtapa3BAC46_WorkerEvaluadorUmbralesHisteresis`) — **FAIL** (falta worker de umbrales con histéresis en Redis y pruebas unitarias) |
| 17 | `FRN-29` | Etapa 3 | Frontend | Luz | 🔴 **No implementada** | `test/front/frn-36-29-snapshots-ui.test.tsx` — **0/4 PASS** (la pestaña Snapshots muestra solo un placeholder sin tabla, modal de creación ni acciones) |
| 18 | `FRN-31` | Etapa 3 | Frontend | Belinda | 🟠 **Implementada parcial (no pasa tests)** | `test/front/frn-31-audit-etapa3.test.tsx` — **1/4 PASS, 3/4 FAIL** (badges de severidad listos; faltan `aria-label="Desde"/"Hasta"`, filtro `recursoTipo`, acciones Etapa 3, modal JSON y botón Exportar JSON) |
