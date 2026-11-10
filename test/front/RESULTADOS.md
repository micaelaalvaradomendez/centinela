# Resultados de las pruebas frontend

**Fecha:** 10/11/2026. **Frontend:** submódulo en `a268063` (últimos commits de `origin/main`: PR #198 `a268063` / `10beb6f` resolviendo **`FRN-33`** con Centro de Notificaciones y banners contextuales de alerta en Dashboard; PR #197 `da3f422` / `087be9c` resolviendo **`ADP-07` / `FRN-37`** con selector y normalizador dinámico de plantillas Debian 12 en wizard de creación; PR #196 `bed77ff` y PR #193 `020b499` resolviendo **`FIX-141`** con el 100 % de los tests unitarios internos en verde; PR #192 `c110791` resolviendo **`FIX-144`** con el vocabulario cerrado `EXITO`/`FALLA`/`PENDIENTE`). **Comando:** `npm --prefix test/front test` (`vitest run`).

| Métrica | Valor |
|---|---:|
| Archivos de prueba (`test/front`) | 49 (**47 en verde**, **0 con fallos**, **2 omitidos** E2E) |
| Pruebas individuales (`test/front`) | 316 |
| Aprueban | **300 (100 % de las pruebas ejecutadas)** |
| Fallan | **0** |
| Omitidas | 16 (`login04-e2e.test.ts` e `int02-e2e.test.ts`: se ejecutan desde `test/back` contra el stack Docker real) |

## Tareas de Frontend verificadas y estado en `terminado3.md` / `terminado2.md`

| Tarea | Archivo | Estado | Detalle |
|---|---|---:|---|
| `FRN-33` (Centro de notificaciones y banners `CRITICAL`) | `frn-33-alerts-center.test.tsx` | ✅ **Pasa (2/2)** | Resuelto en PR #198 (`10beb6f` / `a268063`): `NotificationBell.tsx`, `DashboardAlertBanner.tsx` y `AlertsContext.tsx`. Banner contextual `role="alert"` y notificaciones en vivo vía SSE. |
| `FRN-37` / `ADP-07` (Plantillas Debian 12 Bookworm en Wizard) | `frn-37-adp-07-wizard-debian12.test.tsx` | ✅ **Pasa (2/2)** | Resuelto en PR #197 (`087be9c` / `da3f422`): `parseTemplateOption` y `formatTemplateDisplay` en `createInstanceOptions.ts` e `InstanceCreationSummary.tsx`, reconociendo Debian 12 Bookworm y arquitectura `amd64`. |
| `FIX-141` (Suite unitaria interna `frontend/centinela`) | `frontend/centinela/src/test/*.test.ts(x)` (46 archivos) | ✅ **Pasa (298/298)** | Resuelto en PR #193, PR #196, PR #197 y PR #198: mocks de `useOptionalEventsContext`, `alerts-context` y expectativas de paneles/tablas sincronizadas al 100 %. |
| `FIX-144` (Alineación con enum cerrado de auditoría `RF-08`) | `etapa3-pendientes.test.tsx` (bloque `FIX-144`), `frn-31-audit-etapa3.test.tsx` | ✅ **Pasa (2/2 y 4/4)** | Resuelto en PR #192 (`c110791`): `AUDIT_ACTION_LABELS` coincide 31/31 con `ports/audit_port.go`, `AUDIT_ACTION_OPTIONS` expone las 29 acciones visibles sin duplicados, y `Auditoria.tsx` + `deletedUsersService.ts` usan el enum cerrado `EXITO`/`FALLA`/`PENDIENTE`. |
| `FIX-142` (`FRN-22` + `FRN-23`) | `frn-22-use-instance-metrics.test.tsx`, `frn-23-metrics-charts.test.tsx` | ✅ **Pasa (9/9 y 9/9)** | Restaurado en PR #182 (`26af24f`): `getWsTicket()`, `metricsWebSocket.ts`, `useInstanceMetrics.ts`, pestaña `"Métricas"`, modo `"En vivo"`, aviso `"Instancia detenida"` y gráfico de Red. |
| `FIX-131` (Gráfico de Tráfico de Red `CU-05`) | `etapa3-pendientes.test.tsx` (bloque `FIX-131`) | ✅ **Pasa (3/3)** | Implementado en PR #182 (`fa92405` / `26af24f`): `InstanceMetricPoint` declara `netInBps`/`netOutBps`, `instanceDetailService.ts` los normaliza e `InstanceResourceUsage.tsx` renderiza las 4 tarjetas (`CPU`, `Memoria`, `Disco`, `Tráfico de Red`). |
| `FIX-127` (`SnapshotsTabContent` + `:instanceRef`/`:vmid`) | `etapa3-pendientes.test.tsx` (bloque `FIX-127`) | ✅ **Pasa (3/3)** | Implementado en `83a19b5`/`14db393`/`30a79a5`/`109f266`: `useOptionalEventsContext` y `instanceRef ?? vmid`. |
| `FIX-132` (Errores 400 de contraseña no queman OTP) | `etapa3-pendientes.test.tsx` (bloque `FIX-132`), `recover-password.test.tsx` | ✅ **Pasa (4/4 y 7/7)** | Implementado en `3d06236`: `SAME_PASSWORD`, `WEAK_PASSWORD` y `PASSWORD_RECENTLY_USED` no incrementan `failedAttempts`. |
| `FRN-32` (Menú desplegable de exportación CSV/JSON) | `etapa3-pendientes.test.tsx` (bloque `FRN-32`) | ✅ **Pasa (4/4)** | Implementado en `404c784`: `AuditExportButton` unificado con estado de carga, Toast y transferencia de filtros. |
| `FIX-128` (`FRN-31` Auditoría Etapa 3) | `frn-31-audit-etapa3.test.tsx` | ✅ **Pasa (4/4)** | Implementado en `b38ab17`, `408043e` y `c110791`: acciones de Etapa 3, filtro por `recursoTipo`, badges `EXITO`/`PENDIENTE`/`FALLA`, visor JSON y exportación CSV/JSON. |

## Regresión (Etapas 1, 2 y 3)

Pasan **47 de los 47 archivos ejecutados en local** (300 pruebas en verde), incluyendo:
- **Etapa 1:** `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client`, `navigation`, `admin-users`, `admin-users-deletion`, `admin-recovery`, `recover-password`, `password-change`, `session-security`, `audit`, `auth-branding`, `events-client`, `dashboard-metrics`, `dashboard-node`, `instances-table`, `instances-modals`, `instances-operation`, `instances-ola2`, `instances-metrics`, `instances-events`, `dashboard-summary`, `instance-detail`, `fixes-pending`, `adversarial-security`.
- **Etapa 2:** `frn-22-use-instance-metrics` (9/9), `frn-23-metrics-charts` (9/9), `frn-24-25-create-wizard` (11/11), `frn-26-create-submit-quota-sse` (4/4), `frn-27-28-edit-instance-config` (6/6), `fix-91-dead-code` (2/2), `fix-103-node-history` (7/7), `fix-104b-105-tareaid` (12/12), `fix-108-ids-ofuscados` (11/11), `fix-111b-cambio-email-reauth` (4/4), `fix-115b-contrasenas-largas` (11/11), `fix-117b-dashboard-operador` (4/4).
- **Etapa 3:** `frn-33-alerts-center` (2/2), `frn-37-adp-07-wizard-debian12` (2/2), `frn-36-29-snapshots-ui` (7/7), `frn-30-snapshots-rollback-ui` (2/2), `frn-31-audit-etapa3` (4/4), `frn-34-cookies-a11y` (3/3) y todos los bloques (`FIX-127`, `FIX-131`, `FIX-132`, `FRN-32` y `FIX-144`) en `etapa3-pendientes.test.tsx` (**16/16 en verde**).

## Suite interna del submódulo (`frontend/centinela/src/test/`)

Al ejecutar `npm --prefix frontend/centinela test -- --run` dentro del submódulo:
- **Archivos de prueba:** **46 / 46 en verde (100 %)**.
- **Pruebas individuales:** **298 / 298 en verde (100 %)**.
