# Informe de Estado de Tareas Verificado por Pruebas

**Fecha de ejecución:** 10/11/2026 (tras actualizar los submódulos a los últimos commits de `origin/main` — Backend `d0f3f57`, Frontend `a268063` — y ejecutar todas las suites de pruebas)  
**Alcance:** Evaluación integral de las tareas de [documentacion/actual.md](../documentacion/actual.md), verificación de las tareas completadas y movidas a [documentacion/terminado3.md](../documentacion/terminado3.md), y estado de regresión en las suites `test/front`, `test/back`, `frontend/centinela` y `backend`.

| Submódulo | Revisión probada (`origin/main`) | Últimos commits incorporados |
|---|---|---|
| **Backend** | `d0f3f57` (`Merge pull request #30`) | `d0f3f57` / `df4cada` (no crear en VMID protegidos y exigir pertenencia para borrar), `d92a568` (**`FIX-139`**: reconcilia tareas antes de abrir la API), `5b3ac80` (**`FIX-137`**: elimina endpoints fantasma e incluye auditorías globales), `4011b08` (**`FIX-143`**: vocabulario cerrado de `auditoria.accion` y `auditoria.resultado`). |
| **Frontend** | `a268063` (`Merge pull request #198`) | `a268063` / `10beb6f` (**`FRN-33`**: Centro de Notificaciones con campana y banner contextual `role="alert"` en Dashboard), `da3f422` / `087be9c` (**`ADP-07` / `FRN-37`**: adaptación del wizard y resumen para Debian 12 Bookworm), `f1b77b6` (métricas de tráfico de red en CU-05), `bed77ff` (**`FIX-141`**: suite interna en verde), `c110791` (**`FIX-144`**: enum cerrado de auditoría). |

---

## 1. Tareas Recién Verificadas al 100% y Movidas a `terminado3.md`

En esta corrida se verificaron con pruebas automatizadas y se cerraron las siguientes **3 tareas** que figuraban pendientes en `actual.md`:

1. **`FRN-33`** — Centro de Notificaciones y banners contextuales de alerta por saturación (`Frontend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commits `10beb6f` y `a268063` (PR #198). Incorpora `DashboardAlertBanner.tsx`, `NotificationBell.tsx`, `AlertsContext.tsx` y su integración en `ProtectedLayout.tsx` y `Dashboard.tsx`.
   - **Prueba:** `test/front/frn-33-alerts-center.test.tsx` (2/2 subtests) y `frontend/centinela/src/test/notifications-components.test.tsx` (4/4 subtests) — **100 % PASS**.
2. **`ADP-07` (`FRN-37`)** — Desacoplamiento de constantes y selección de plantillas Debian 12 Bookworm en el Asistente de Creación (`Frontend`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Commits `087be9c` y `da3f422` (PR #197). Incorpora `parseTemplateOption` y `formatTemplateDisplay` en `createInstanceOptions.ts` e `InstanceCreationSummary.tsx`, clasificando dinámicamente plantillas Debian (11, 12, 13), Ubuntu y Alpine con extracción de versión y arquitectura (`amd64`).
   - **Prueba:** `test/front/frn-37-adp-07-wizard-debian12.test.tsx` (2/2 subtests) y `frontend/centinela/src/test/create-instance-wizard.test.tsx` (11/11 subtests) — **100 % PASS**.
3. **`FIX-129` (`INF-11`)** — Endurecimiento de cabeceras de seguridad en Nginx de producción (`Infraestructura`, movida a [terminado3.md](../documentacion/terminado3.md)):
   - **Evidencia:** Configuración en `docker/nginx-edge.conf` con `Strict-Transport-Security "max-age=31536000; includeSubDomains"`, `X-Frame-Options "DENY"`, `X-Content-Type-Options "nosniff"` y `Content-Security-Policy "default-src 'self' ... connect-src 'self' ws: wss: ..."`.
   - **Prueba:** `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS` (3/3 subtests) en `test/back/etapa3_acceptance_test.go` — **100 % PASS**.

---

## 2. Estado de las Tareas Pendientes en `documentacion/actual.md`

*(Todas las tareas de desarrollo de software —tanto de Frontend como de Backend— están 100 % completadas y con pruebas en verde. Únicamente restan las tareas que requieren conexión con el servidor Proxmox VE físico de producción).*

| # | Tarea | Etapa | Área | Estado en Pruebas | Diagnóstico y Requerimiento |
|---|---|---|---|:---:|---|
| 1 | **`FIX-123`** | Etapa 2 | Infra | ⚪ **Omitida sin host físico** (`SKIP` en `TestFIX123_PlantillasEISOsEnProxmoxReal`) | Requiere `PROXMOX_REAL_URL`, `PROXMOX_REAL_TOKEN_ID` y `PROXMOX_REAL_TOKEN_SECRET` contra el servidor Proxmox VE físico. |
| 2 | **`INF-14`** | Etapa 3 | Infra / Docs | 🔴 **Pendiente en host físico** | Requiere ejecutar `scripts/pve8-audit-probe.sh` en el hipervisor PVE 8.4.21 y archivar capturas en `backend/API proxmox respuestas/pve8-capturas/`. |
| 3 | **`ADP-09`** | PVE 8.4 | Infra / Docs | 🔴 **Pendiente en host físico** | Requiere configurar usuario `centi-api@pve`, rol `CentinelaAdmin` y API token en el nodo físico `10.x.x.x`. |
| 4 | **`ADP-10`** | PVE 8.4 | Infra / Docs | 🔴 **Pendiente en host físico** | Requiere aprovisionar storage LVM-Thin (`local-lvm`) y descargar plantillas Debian 12 en `local:vztmpl` en el host físico. |

---

## 3. Resultados Globales de Ejecución de las Suites

| Suite | Comando | Archivos / Paquetes | Pruebas Individuales | Estado |
|---|---|---|---|---|
| **Aceptación Frontend (`test/front`)** | `npm --prefix test/front test` | **49 archivos** (**47 PASS**, **0 FAIL**, **2 SKIPPED** E2E) | **316 pruebas** (**300 PASS**, **0 FAIL**, **16 SKIPPED**) | 🟢 **100 % de pruebas ejecutadas en verde** (`FRN-33` y `ADP-07`/`FRN-37` cerradas al 100 %). |
| **Interna Frontend (`frontend/centinela`)** | `npm --prefix frontend/centinela test -- --run` | **46 archivos** (**46 PASS**, **0 FAIL**) | **298 pruebas** (**298 PASS**, **0 FAIL**) | 🟢 **100 % en verde** (incluyendo `notifications-components`, `alerts-context`, `create-instance-wizard`). |
| **Interna Backend (`backend`)** | `go test -C backend ./...` | **14 paquetes** (**14 PASS**, **0 FAIL**) | **100 % PASS** | 🟢 **100 % en verde** (todos los servicios, repositorios, HTTP y migraciones). |
| **Aceptación e Integración Backend (`test/back`)** | `go test -C test/back -v ./...` | **24 archivos de test** | Todas las pruebas de integración, seguridad, control de acceso y regresión ejecutadas | 🟢 **100 % de integración y backend en verde** (`TestTareasBackendEIntegracion`, `TestHitoSeguridadDeSesionesYCookies`, `TestAdversarialSystemIntegrity`, `TestHitoControlDeAccesoBasadoEnRecursos`, `TestFixesPendientes`, `TestEtapa1`, `TestEtapa1Ola2`, `TestEtapa1Ola3BAC26`, `TestPuenteEtapa1`, `TestEtapa3INF11`). |
