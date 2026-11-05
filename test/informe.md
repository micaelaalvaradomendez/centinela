# Informe de Estado de Tareas Verificado por Pruebas

**Fecha de ejecución:** 05/11/2026 (tras actualizar los submódulos a los últimos commits de `origin/main` y ejecutar todas las suites de pruebas)  
**Alcance:** Evaluación integral de las tareas de [documentacion/actual.md](../documentacion/actual.md), verificación de las tareas recién completadas y movidas a [documentacion/terminado2.md](../documentacion/terminado2.md) y [documentacion/terminado3.md](../documentacion/terminado3.md), y estado de regresión en las suites `test/front`, `test/back`, `frontend/centinela` y `backend`.

| Submódulo | Revisión probada (`origin/main`) | Últimos commits incorporados |
|---|---|---|
| **Backend** | `d0f3f57` (`Merge pull request #30`) | `d0f3f57` / `df4cada` (no crear en VMID protegidos y exigir pertenencia para borrar), `d92a568` (**`FIX-139`**: reconcilia tareas antes de abrir la API), `5b3ac80` (**`FIX-137`**: elimina endpoints fantasma e incluye auditorías globales), `4011b08` (**`FIX-143`**: vocabulario cerrado de `auditoria.accion` y `auditoria.resultado`). |
| **Frontend** | `bed77ff` (`Merge pull request #196`) | `bed77ff` / `0bf4f33` y `020b499` / `1b7e7d4` / `bc44649` (**`FIX-141`**: alinea pruebas de paneles, tablas y mocks internos en `frontend/centinela`), `c110791` / `efd59a5`..`a2a2761` (**`FIX-144`**: alinea auditoría con el vocabulario cerrado `EXITO`/`FALLA`/`PENDIENTE` y retira lectura doble), `5370f31` (red estática por defecto en el asistente de creación y manejo de `403 INSTANCE_NOT_OWNED`). |

---

## 1. Tareas Verificadas al 100% y Movidas a `terminado2.md` / `terminado3.md`

En esta corrida se verificaron con pruebas automatizadas y se cerraron las siguientes **4 tareas** que estaban pendientes en `actual.md`:

1. **`FIX-137`** — Saneamiento de código inerte: eliminación de stub documental `PausarInstancia`, inclusión de auditorías globales del sistema (`organizacion_id IS NULL`) y limpieza de tipos DTO huérfanos (`Backend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commit `5b3ac80` en `instance_handler.go`, `audit_repository.go`, `contrato_etapa1.go` y `docs/swagger.json`.
   - **Prueba:** `TestFIX137_CodigoInerteYAuditoriaDeSistema` (3/3 subtests) en `test/back/etapa3_fixes_pendientes_test.go` — **PASS**.
2. **`FIX-139`** — Reconciliación inmediata de tareas al inicio del backend para mitigar falsos positivos de `409 INSTANCE_BUSY` (`Backend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commit `d92a568` en `seguimiento_tareas.go`, `seguimiento_tareas_test.go` y `cmd/api/main.go`. Expone `RecuperarYReconciliarAlArranque`, la invoca antes de abrir la API HTTP e incluye tests unitarios.
   - **Prueba:** `TestFIX139_ReconciliacionDeTareasAlArranque` (3/3 subtests) en `test/back/etapa3_fixes_pendientes_test.go` — **PASS**.
3. **`FIX-141`** — Actualización de mocks de `@/context/EventsContext` (`useOptionalEventsContext`) y sincronización de expectativas UI en tests unitarios internos del submódulo `frontend` (`Testing Frontend`, movida a [terminado2.md](../documentacion/terminado2.md)):
   - **Evidencia:** Commits `bc44649`, `1b7e7d4`, `020b499` (PR #193) y `0bf4f33`, `bed77ff` (PR #196).
   - **Prueba:** `npm --prefix frontend/centinela test -- --run` (**43/43 archivos y 277/277 pruebas individuales**) — **100 % PASS**.
4. **`FIX-144`** — Alineación de filtros, contadores y diccionario de acciones de auditoría con el enum cerrado del Backend y retiro de lectura doble (`Frontend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commits `a2a2761`, `798a93f`, `7ca0558`, `af78395`, `efd59a5` y `c110791` (PR #192) en `auditActionLabels.ts`, `Auditoria.tsx`, `auditFilters.ts`, `auditService.ts` y `deletedUsersService.ts`.
   - **Prueba:** Bloque `FIX-144` en `test/front/etapa3-pendientes.test.tsx` (2/2 pruebas), `test/front/frn-31-audit-etapa3.test.tsx` (4/4 pruebas) y `frontend/centinela/src/test/audit-action-labels.test.ts` (3/3 pruebas) — **PASS**.

---

## 2. Estado de las Tareas Pendientes en `documentacion/actual.md`

*(Todas las tareas de **Backend** en `actual.md` ya están implementadas y en verde; solo restan 2 tareas de Frontend y las tareas de Infraestructura/Documentación de Proxmox físico y Nginx).*

| # | Tarea | Etapa | Área | Estado en Pruebas | Diagnóstico y Trabajo Pendiente |
|---|---|---|---|:---:|---|
| 1 | **`ADP-07` (`FRN-37`)** | PVE 8.4 | Frontend | 🔴 **Pendiente** (`2/3 PASS, 1/3 FAIL` en `frn-37-adp-07-wizard-debian12.test.tsx`) | `extractTemplateMetadata()` en `instances.service.ts` solo reconoce Ubuntu (`24.04`/`22.04`) y devuelve `'—'` como sistema operativo para `local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst` en lugar de `operatingSystem: 'Debian'` y `version: '12'`. |
| 2 | **`FRN-33`** | Etapa 3 | Frontend | 🟠 **Parcial** (`1/2 PASS, 1/2 FAIL` en `frn-33-alerts-center.test.tsx`) | `eventsClient.ts` ya normaliza `AlertEventPayload` y `Home.tsx` renderiza el banner, pero falta que el contenedor de alerta crítica (`role="alert"`) se renderice en la cabecera ante cualquier `ALERT_TRIGGERED` sin depender de que el texto libre contenga la subcadena `'crítica'` con tilde (`severity === 'CRITICAL'`). |
| 3 | **`FIX-129` (`INF-11`)** | Etapa 3 | Infra | 🟠 **Parcial** (`2/3 PASS, 1/3 FAIL` en `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS`) | `docker/nginx-edge.conf` ya tiene redirección `301` HTTP->HTTPS, `TLSv1.2 TLSv1.3`, WSS y SSE sin buffering, pero le faltan las 4 cabeceras de seguridad: `Strict-Transport-Security` (`includeSubDomains`), `Content-Security-Policy`, `X-Frame-Options` y `X-Content-Type-Options`. |
| 4 | **`FIX-123`** | Etapa 2 | Infra | ⚪ **Omitida sin host físico** (`SKIP` en `TestFIX123_PlantillasEISOsEnProxmoxReal`) | Requiere `PROXMOX_REAL_URL`, `PROXMOX_REAL_TOKEN_ID` y `PROXMOX_REAL_TOKEN_SECRET` contra el servidor Proxmox VE físico. |
| 5 | **`INF-14`** | Etapa 3 | Infra / Docs | 🔴 **Pendiente** (`FAIL` en `TestEtapa3INF14_AuditoriaYSondaPVE8Real`) | Falta el script `scripts/pve8-audit-probe.sh` y la sección de auditoría empírica PVE 8.4 en `documentacion/api-proxmox.md`. |
| 6 | **`ADP-09` / `ADP-10`** | PVE 8.4 | Infra / Docs | 🔴 **Pendiente** (`FAIL` en `TestADP09_ADP10_InfraestructuraPVE8`) | Falta el script `scripts/pve8-audit-probe.sh` y documentar en `documentacion/api-proxmox.md` los volúmenes y permisos verificados en PVE 8.4.21. |

---

## 3. Resultados Globales de Ejecución de las 4 Suites

| Suite | Comando | Archivos / Paquetes | Pruebas Individuales | Estado |
|---|---|---|---|---|
| **Aceptación Frontend (`test/front`)** | `npm --prefix test/front test` | **49 archivos** (**45 PASS**, **2 FAIL**, **2 SKIPPED** E2E) | **316 pruebas** (**298 PASS**, **2 FAIL**, **16 SKIPPED**) | 🟢 **99,3 % de pruebas ejecutadas en verde** (los únicos 2 fallos son las 2 tareas pendientes de front en `actual.md`: `FRN-33` y `ADP-07`/`FRN-37`). |
| **Interna Frontend (`frontend/centinela`)** | `npm --prefix frontend/centinela test -- --run` | **43 archivos** (**43 PASS**, **0 FAIL**) | **277 pruebas** (**277 PASS**, **0 FAIL**) | 🟢 **100 % en verde** (**`FIX-141`** cerrado). |
| **Interna Backend (`backend`)** | `go test -C backend ./...` | **14 paquetes** (**14 PASS**, **0 FAIL**) | **100 % PASS** | 🟢 **100 % en verde** (incluyendo `TestContratoVocabularioAuditoria`, `seguimiento_tareas_test.go` y `migrations/embed_test.go`). |
| **Aceptación Backend (`test/back`)** | `go test -C test/back -v ./...` | **24 archivos de test** | Todas las pruebas de Etapas 1, 2, 3 y PVE 8.4 ejecutadas | 🟢 **100 % de las tareas de Backend en verde** (`FIX-128`, `FIX-133`..`FIX-140`, `FIX-143`, `BAC-42`..`BAC-55`, `ADP-01`..`ADP-06`); en pruebas sin Docker solo fallan las 3 tareas de Infra/Docs abiertas en `actual.md` (`FIX-129`/`INF-11`, `INF-14`, `ADP-09`/`ADP-10`). |
