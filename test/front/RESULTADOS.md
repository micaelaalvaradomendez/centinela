# Resultados de las pruebas frontend

**Fecha:** 05/11/2026. **Frontend:** submódulo en `bed77ff` (últimos commits de `origin/main`: PR #196 `bed77ff` y PR #193 `020b499` resolviendo **`FIX-141`** con el 100 % de los tests unitarios internos en verde; PR #192 `c110791` resolviendo **`FIX-144`** al alinear `auditActionLabels.ts`, `Auditoria.tsx`, `auditFilters.ts` y `deletedUsersService.ts` con el vocabulario cerrado `EXITO`/`FALLA`/`PENDIENTE`; `5370f31` red estática por defecto en el asistente de creación y manejo de `403 INSTANCE_NOT_OWNED`). **Comando:** `npm --prefix test/front test` (`vitest run`).

| Métrica | Valor |
|---|---:|
| Archivos de prueba (`test/front`) | 49 (**45 en verde**, **2 con fallos**, **2 omitidos** E2E) |
| Pruebas individuales (`test/front`) | 316 |
| Aprueban | **298** |
| Fallan | **2** (correspondientes a las únicas 2 tareas de Frontend pendientes en `documentacion/actual.md`: `FRN-33` y `ADP-07`/`FRN-37`) |
| Omitidas | 16 (`login04-e2e.test.ts` e `int02-e2e.test.ts`: se ejecutan desde `test/back` contra el stack Docker real) |

## Tareas de Frontend verificadas y estado en `actual.md`

| Tarea | Archivo | Estado | Detalle |
|---|---|---:|---|
| `FIX-141` (Suite unitaria interna `frontend/centinela`) | `frontend/centinela/src/test/*.test.ts(x)` (43 archivos) | ✅ **Pasa (277/277)** | Resuelto en PR #193 (`020b499`) y PR #196 (`bed77ff`): mocks de `useOptionalEventsContext` y expectativas de paneles/tablas sincronizadas al 100 %. |
| `FIX-144` (Alineación con enum cerrado de auditoría `RF-08`) | `etapa3-pendientes.test.tsx` (bloque `FIX-144`), `frn-31-audit-etapa3.test.tsx` | ✅ **Pasa (2/2 y 4/4)** | Resuelto en PR #192 (`c110791`): `AUDIT_ACTION_LABELS` coincide 31/31 con `ports/audit_port.go`, `AUDIT_ACTION_OPTIONS` expone las 29 acciones visibles sin duplicados ni códigos fantasma, y `Auditoria.tsx` + `deletedUsersService.ts` usan el enum cerrado `EXITO`/`FALLA`/`PENDIENTE`. |
| `FIX-142` (`FRN-22` + `FRN-23`) | `frn-22-use-instance-metrics.test.tsx`, `frn-23-metrics-charts.test.tsx` | ✅ **Pasa (9/9 y 9/9)** | Restaurado en PR #182 (`26af24f`): `getWsTicket()`, `metricsWebSocket.ts`, `useInstanceMetrics.ts`, pestaña `"Métricas"`, modo `"En vivo"`, aviso `"Instancia detenida"` y gráfico de Red. |
| `FIX-131` (Gráfico de Tráfico de Red `CU-05`) | `etapa3-pendientes.test.tsx` (bloque `FIX-131`) | ✅ **Pasa (3/3)** | Implementado en PR #182 (`fa92405` / `26af24f`): `InstanceMetricPoint` declara `netInBps`/`netOutBps`, `instanceDetailService.ts` los normaliza e `InstanceResourceUsage.tsx` renderiza las 4 tarjetas (`CPU`, `Memoria`, `Disco`, `Tráfico de Red`). |
| `FIX-127` (`SnapshotsTabContent` + `:instanceRef`/`:vmid`) | `etapa3-pendientes.test.tsx` (bloque `FIX-127`) | ✅ **Pasa (3/3)** | Implementado en `83a19b5`/`14db393`/`30a79a5`/`109f266`: `useOptionalEventsContext` y `instanceRef ?? vmid`. |
| `FIX-132` (Errores 400 de contraseña no queman OTP) | `etapa3-pendientes.test.tsx` (bloque `FIX-132`), `recover-password.test.tsx` | ✅ **Pasa (4/4 y 7/7)** | Implementado en `3d06236`: `SAME_PASSWORD`, `WEAK_PASSWORD` y `PASSWORD_RECENTLY_USED` no incrementan `failedAttempts`. |
| `FRN-32` (Menú desplegable de exportación CSV/JSON) | `etapa3-pendientes.test.tsx` (bloque `FRN-32`) | ✅ **Pasa (4/4)** | Implementado en `404c784`: `AuditExportButton` unificado con estado de carga, Toast y transferencia de filtros. |
| `FIX-128` (`FRN-31` Auditoría Etapa 3) | `frn-31-audit-etapa3.test.tsx` | ✅ **Pasa (4/4)** | Implementado en `b38ab17`, `408043e` y `c110791`: acciones de Etapa 3, filtro por `recursoTipo`, badges `EXITO`/`PENDIENTE`/`FALLA`, visor JSON y exportación CSV/JSON. |
| `FRN-33` (Centro de notificaciones y banners `CRITICAL`) | `frn-33-alerts-center.test.tsx` | ❌ **Falla (1 de 2)** | **Pendiente en `actual.md`:** falta que el banner rojo contextual (`role="alert"`) en la cabecera del Dashboard se renderice ante eventos `ALERT_TRIGGERED` de saturación crítica sin depender de que el mensaje contenga `'crítica'` con tilde. |
| `FRN-37` / `ADP-07` (Plantillas Debian 12 Bookworm en Wizard) | `frn-37-adp-07-wizard-debian12.test.tsx` | ❌ **Falla (1 de 3)** | **Pendiente en `actual.md`:** `fetchCreationTemplates` devuelve `'—'` en `operatingSystem` para `debian-12-standard_12.7-1_amd64.tar.zst` en lugar de extraer dinámicamente `"Debian"` y versión `"12"`. |

## Regresión (Etapas 1, 2 y 3)

Pasan **45 de los 47 archivos ejecutados en local** (298 pruebas en verde), incluyendo:
- **Etapa 1:** `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client`, `navigation`, `admin-users`, `admin-users-deletion`, `admin-recovery`, `recover-password`, `password-change`, `session-security`, `audit`, `auth-branding`, `events-client`, `dashboard-metrics`, `dashboard-node`, `instances-table`, `instances-modals`, `instances-operation`, `instances-ola2`, `instances-metrics`, `instances-events`, `dashboard-summary`, `instance-detail`, `fixes-pending`, `adversarial-security`.
- **Etapa 2:** `frn-22-use-instance-metrics` (9/9), `frn-23-metrics-charts` (9/9), `frn-24-25-create-wizard` (11/11), `frn-26-create-submit-quota-sse` (4/4), `frn-27-28-edit-instance-config` (6/6), `fix-91-dead-code` (2/2), `fix-103-node-history` (7/7), `fix-104b-105-tareaid` (12/12), `fix-108-ids-ofuscados` (11/11), `fix-111b-cambio-email-reauth` (4/4), `fix-115b-contrasenas-largas` (11/11), `fix-117b-dashboard-operador` (4/4).
- **Etapa 3:** `frn-36-29-snapshots-ui` (7/7), `frn-30-snapshots-rollback-ui` (2/2), `frn-31-audit-etapa3` (4/4), `frn-34-cookies-a11y` (3/3) y todos los bloques (`FIX-127`, `FIX-131`, `FIX-132`, `FRN-32` y `FIX-144`) en `etapa3-pendientes.test.tsx` (**16/16 en verde**).

## Suite interna del submódulo (`frontend/centinela/src/test/`) — `FIX-141` Cerrado

Al ejecutar `npm --prefix frontend/centinela test -- --run` dentro del submódulo:
- **Archivos de prueba:** **43 / 43 en verde (100 %)**.
- **Pruebas individuales:** **277 / 277 en verde (100 %)**.
