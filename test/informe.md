# Informe de Estado de Tareas Verificado por Pruebas

**Fecha de ejecución:** 04/11/2026 (tras actualizar los submódulos a los últimos commits de `origin/main` y ejecutar todas las suites de pruebas)  
**Alcance:** Evaluación integral de las tareas de [documentacion/actual.md](../documentacion/actual.md), verificación de las tareas recién completadas y movidas a [documentacion/terminado2.md](../documentacion/terminado2.md) y [documentacion/terminado3.md](../documentacion/terminado3.md), y estado de regresión en las suites `test/front`, `test/back`, `frontend/centinela` y `backend`.

| Submódulo | Revisión probada (`origin/main`) | Últimos commits incorporados |
|---|---|---|
| **Backend** | `4011b08` (`fix(auditoria): vocabulario cerrado de accion y resultado, enums en Swagger y normalización histórica`) | `4011b08` (**`FIX-143`**: vocabulario cerrado de `auditoria.accion` y `auditoria.resultado`, enums en Swagger y `contrato-etapa1.md`, `TestContratoVocabularioAuditoria` y migración `20261104_vocabulario_auditoria.sql`), `08126ae` (**`FIX-136`**: umbrales e histéresis de alertas configurables por entorno), `5f8bb5e` (**`FIX-135`**: detección y alerta de paradas externas de instancias), `42093c3` (**`FIX-134`**: soporte para instancias `PAUSED`). |
| **Frontend** | `1e5ec49` (`Merge pull request #189`) | `1e5ec49` / `edd581c` (mejoras en `Auditoria.tsx`), `26af24f` (**`FIX-142` + `FIX-131`**: normalización de `AlertEventPayload` en `eventsClient.ts`, banner crítico `role="alert"` en `Home.tsx`, tarjeta de almacenamiento para `OPERATOR` y selector `"En vivo"` en `InstanceResourceUsage.tsx`), `99c7a18` (`Instances.tsx`). |

---

## 1. Tareas Verificadas al 100% y Movidas a `terminado2.md` / `terminado3.md`

En esta corrida se verificaron con pruebas automatizadas y se cerraron las siguientes **6 tareas** que estaban en `actual.md`:

1. **`FIX-143`** — Vocabulario de auditoría con contrato cerrado, eliminación de códigos duplicados, normalización histórica y publicación de enums en Swagger (`Backend / Base de Datos`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commit `4011b08` en `ports/audit_port.go`, `audit_handler.go`, `audit_repository.go`, `docs/swagger.json`, `docs/contrato-etapa1.md`, `docs/contrato_etapa1_test.go` y `migrations/20261104_vocabulario_auditoria.sql`. Unifica `auditoria.resultado` en `EXITO`, `FALLA`, `PENDIENTE`, elimina constantes muertas/duplicadas, publica los `enum` en Swagger y documenta la tabla oficial en `contrato-etapa1.md`.
   - **Prueba:** `TestFIX143_VocabularioAuditoriaContratoCerrado` (4/4 subtests) en `test/back/etapa3_fixes_pendientes_test.go` y `TestContratoVocabularioAuditoria` en `backend/docs/contrato_etapa1_test.go` — **PASS**.
2. **`FIX-142`** — Tarjeta de almacenamiento en Dashboard para rol `OPERATOR` (`Frontend`, movida a [terminado2.md](../documentacion/terminado2.md)):
   - **Evidencia:** Commit `26af24f` en `Home.tsx`, `NodeCapacitySection.tsx` y `node.types.ts`. El bloque de almacenamiento se renderiza tanto para `ADMIN` como para `OPERATOR` consumiendo `nodeStatus.storage`.
   - **Prueba:** `test/front/fix-117b-dashboard-operador.test.tsx` (4/4 pruebas) — **PASS**.
3. **`FIX-131`** — Selector `"En vivo"` (WebSocket) y gráfico de Tráfico de Red en detalle de instancia (`Frontend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commit `26af24f` en `InstanceResourceUsage.tsx`. Incorpora la opción `"En vivo"` en el selector de período, invoca `useInstanceMetrics(instanceId, { enabled: timeframe === 'live' })` y muestra el badge de estado de conexión (`En vivo` / `Reconectando...`).
   - **Prueba:** `test/front/frn-23-metrics-charts.test.tsx` (9/9 pruebas) y `test/front/frn-22-use-instance-metrics.test.tsx` (9/9 pruebas) — **PASS**.
4. **`FIX-134`** — Soporte completo para instancias en estado `PAUSED` (`Backend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commit `42093c3` en `instance.go`, `instance_service.go`, `proxmox_port.go` y `client.go`. Mapea `"paused"` a `EstadoInstanciaPausada = "PAUSED"`, rechaza `start`/`stop`/`shutdown`/`reboot` y `DELETE` sobre instancias pausadas con `409 Conflict (INSTANCE_INVALID_STATE)` y agrega soporte para `resume`.
   - **Prueba:** `TestFIX134_InstanciasPausadas` (4/4 subtests) en `test/back/etapa3_fixes_pendientes_test.go` — **PASS**.
5. **`FIX-135`** — Detección y alerta de paradas externas de instancias (`Backend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commit `5f8bb5e` en `instance_service.go`, `seguimiento_tareas.go` y `alert_evaluator_service.go`. Registra intención de parada en Redis (`instance:stop_intent:{vmid}`), detecta transiciones `RUNNING -> STOPPED` no iniciadas desde Centinela, audita `SYSTEM_INSTANCE_STOPPED_EXTERNALLY` y despacha alerta crítica por SSE y correo.
   - **Prueba:** `TestFIX135_AlertasParadasExternas` (3/3 subtests) en `test/back/etapa3_fixes_pendientes_test.go` — **PASS**.
6. **`FIX-136`** — Parametrización por entorno de umbrales e histéresis del evaluador de alertas (`Backend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commit `08126ae` en `cmd/api/main.go` y `alert_evaluator_worker.go`. Lee los umbrales y ciclos de persistencia desde variables de entorno con validación de rango e inyección en `NewEvaluadorAlertas`.
   - **Prueba:** `TestFIX136_UmbralesDeAlertaConfigurables` (4/4 subtests) en `test/back/etapa3_fixes_pendientes_test.go` — **PASS**.

---

## 2. Estado de las 11 Tareas Pendientes en `documentacion/actual.md`

| # | Tarea | Etapa | Área | Estado en Pruebas | Diagnóstico y Trabajo Pendiente |
|---|---|---|---|:---:|---|
| 1 | **`FIX-141`** | Etapa 2 | Frontend | 🔴 **Pendiente** (`5 archivos / 12 tests FAIL` en `frontend/centinela`) | Desalineación en la suite interna del submódulo `frontend/centinela`: mocks incompletos en `Auditoria.test.tsx`, `useInstanceDetail.test.ts`, `InstanceSummarySection.test.tsx`, assertions obsoletas en `Home.test.tsx` (`OPERATOR`) y validación de `validateInstanceRef` en `instanceService.test.ts`. |
| 2 | **`ADP-07` (`FRN-37`)** | PVE 8.4 | Frontend | 🔴 **Pendiente** (`2/3 PASS, 1/3 FAIL` en `frn-37-adp-07-wizard-debian12.test.tsx`) | `extractTemplateMetadata()` en `instances.service.ts:276-306` solo reconoce Ubuntu (`24.04`/`22.04`) y devuelve `'—'` como sistema operativo para `local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst` en lugar de `operatingSystem: 'Debian'` y `version: '12'`. |
| 3 | **`FRN-33`** | Etapa 3 | Frontend | 🟠 **Parcial** (`1/2 PASS, 1/2 FAIL` en `frn-33-alerts-center.test.tsx`) | `eventsClient.ts` ya normaliza `AlertEventPayload` (`metric`, `value`, `threshold`, `severity`, `message`) y `Home.tsx` renderiza el banner, pero falta que el contenedor de alerta crítica se renderice en la cabecera ante cualquier `ALERT_TRIGGERED` sin depender de que el texto libre contenga la subcadena `'crítica'` con tilde (`severity === 'CRITICAL'`). |
| 4 | **`FIX-137`** | Etapa 3 | Backend | 🔴 **Pendiente** (`0/3 PASS` en `TestFIX137_CodigoInerteYAuditoriaDeSistema`) | 1) El stub `PausarInstancia` (`501 Not Implemented`) sigue expuesto en `instance_handler.go` y Swagger; 2) `audit_repository.go` oculta eventos del sistema con `usuario_id IS NULL AND organizacion_id IS NULL` cuando el admin tiene `organizacion_id`; 3) `contrato_etapa1.go` conserva structs DTO huérfanos. |
| 5 | **`FIX-139` (`FIX-123`)** | Etapa 3 | Backend | 🔴 **Pendiente** (`0/3 PASS` en `TestFIX139_ReconciliacionDeTareasAlArranque`) | No existe `RecuperarYReconciliarAlArranque` en `seguimiento_tareas.go` ni se invoca en `cmd/api/main.go`. Además, las tareas `CREATE` y `DELETE` retomadas tras un reinicio pierden su callback de autoasignación/purga de permisos. |
| 6 | **`FIX-144`** | Etapa 3 | Frontend | 🔴 **Pendiente** (`0/2 PASS` en `etapa3-pendientes.test.tsx` bloque `FIX-144`) | Ya cerrado `FIX-143` en el backend (`EXITO`/`FALLA`/`PENDIENTE`), falta sincronizar el frontend: 1) `AUDIT_ACTION_OPTIONS` en `auditActionLabels.ts` aún incluye códigos inexistentes (`START_VM`, `STOP_VM`, `CAMBIAR_CONTRASENA`, `EXPORTAR_AUDITORIA`, `SUSPENDER_VM`, `ACTUALIZAR_PERFIL`) y omite acciones reales (`INSTANCE_PERMISSIONS_PURGED`, `UPDATE_PROFILE_EMAIL_REQUESTED`, `UPDATE_PROFILE_EMAIL_CONFIRMED`, `SYSTEM_INSTANCE_STOPPED_EXTERNALLY`, `EVENTOS_CONEXION`, `EVENTOS_CIERRE`); 2) `Auditoria.tsx` sigue filtrando con `SUCCESS`/`FAILED`/`PENDING` en lugar de `EXITO`/`FALLA`/`PENDIENTE`. |
| 7 | **`FIX-129` (`INF-11`)** | Etapa 3 | Infra | 🟠 **Parcial** (`2/3 PASS, 1/3 FAIL` en `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS`) | `docker/nginx-edge.conf` ya tiene redirección `301` HTTP->HTTPS, `TLSv1.2 TLSv1.3`, WSS y SSE sin buffering, pero le faltan las 4 cabeceras de seguridad: `Strict-Transport-Security` (`includeSubDomains`), `Content-Security-Policy`, `X-Frame-Options` y `X-Content-Type-Options`. |
| 8 | **`FIX-123` (original)** | Etapa 3 | Docs / Back | ➡️ **Reemplazada por `FIX-139`** | El usuario indicó referenciarla bajo el identificador `FIX-139` en `actual.md`. |
| 9 | **`INF-14`** | Etapa 3 | Infra / Docs | 🔴 **Pendiente** (`FAIL` en `TestEtapa3INF14_AuditoriaYSondaPVE8Real`) | Falta el script `scripts/pve8-audit-probe.sh` y la sección de auditoría empírica PVE 8.4 en `documentacion/api-proxmox.md`. |
| 10 | **`ADP-09`** | PVE 8.4 | Infra | 🔴 **Pendiente** (`FAIL` en `TestADP09_ADP10_InfraestructuraPVE8`) | Falta aprovisionar la plantilla `debian-12-standard` e ISO Debian 12 en el hipervisor y documentar los volúmenes verificados en `documentacion/api-proxmox.md`. |
| 11 | **`ADP-10`** | PVE 8.4 | Infra / QA | 🔴 **Pendiente** (`FAIL` en `TestADP09_ADP10_InfraestructuraPVE8`) | Falta documentar en `documentacion/api-proxmox.md` la matriz de validación E2E contra el servidor Proxmox VE 8.4 real. |

---

## 3. Resultados Globales de Ejecución de las 4 Suites

| Suite | Comando | Archivos / Paquetes | Pruebas Individuales | Estado |
|---|---|---|---|---|
| **Aceptación Frontend (`test/front`)** | `npm --prefix test/front test` | **49 archivos** (**44 PASS**, **3 FAIL**, **2 SKIPPED** E2E) | **316 pruebas** (**296 PASS**, **4 FAIL**, **16 SKIPPED**) | 🟢 **98,7 % de pruebas ejecutadas en verde** (los únicos 4 fallos son las 3 tareas pendientes de front en `actual.md`: `FRN-33`, `ADP-07`/`FRN-37` y `FIX-144`). |
| **Interna Frontend (`frontend/centinela`)** | `npm --prefix frontend/centinela test -- --run` | **43 archivos** (**38 PASS**, **5 FAIL**) | **272 pruebas** (**260 PASS**, **12 FAIL**) | 🟠 Los 12 fallos internos están documentados y rastreados en **`FIX-141`**. |
| **Interna Backend (`backend`)** | `go test -C backend ./...` | **14 paquetes** (**14 PASS**, **0 FAIL**) | **100 % PASS** | 🟢 **100 % en verde** (incluyendo `TestContratoVocabularioAuditoria` y `migrations/embed_test.go`). |
| **Aceptación Backend (`test/back`)** | `go test -C test/back -v ./...` | **24 archivos de test** | Todas las pruebas de Etapas 1, 2, 3 y PVE 8.4 ejecutadas | 🟢 Pasan todas las tareas cerradas (`FIX-143`, `FIX-134`, `FIX-135`, `FIX-136`, `FIX-128`, `FIX-133`, `FIX-138`, `FIX-140`, `BAC-42..55`, `ADP-01..06`); fallan únicamente las pruebas de las tareas aún abiertas en `actual.md` (`FIX-137`, `FIX-139`, `FIX-129`/`INF-11`, `INF-14`, `ADP-09`/`ADP-10`). |
