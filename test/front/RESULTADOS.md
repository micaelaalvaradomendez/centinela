# Resultados de las pruebas frontend

**Fecha:** 04/11/2026. **Frontend:** submódulo en `1e5ec49` (último commit de `origin/main`; nuevos: PR #182 `26af24f` restaurando `useInstanceMetrics`, `getWsTicket`, pestaña `"Métricas"`, modo `"En vivo"` y gráfico de Tráfico de Red `FIX-142` / `FIX-131`; PR #189 `edd581c` catálogo y filtros desplegables en `Auditoria.tsx`; `99c7a18` rediseño de tabla en `Instances.tsx`). **Comando:** `npm --prefix test/front test` (`vitest run`).

| Métrica | Valor |
|---|---:|
| Archivos de prueba | 49 (**44 en verde**, **3 con fallos**, **2 omitidos** E2E) |
| Pruebas individuales | 316 |
| Aprueban | **296** |
| Fallan | **4** (todas correspondientes a las 3 tareas de Frontend pendientes en `documentacion/actual.md`) |
| Omitidas | 16 (`login04-e2e.test.ts` e `int02-e2e.test.ts`: se ejecutan desde `test/back` contra el stack Docker real) |

## Tareas de `actual.md` (Frontend)

| Tarea | Archivo | Estado | Detalle |
|---|---|---:|---|
| `FIX-142` (`FRN-22` + `FRN-23`) | `frn-22-use-instance-metrics.test.tsx`, `frn-23-metrics-charts.test.tsx` | ✅ **Pasa (9/9 y 9/9)** | Restaurado en PR #182 (`26af24f`): `getWsTicket()`, `metricsWebSocket.ts`, `useInstanceMetrics.ts`, pestaña `"Métricas"`, modo `"En vivo"`, aviso `"Instancia detenida"` y gráfico de Red. |
| `FIX-131` (Gráfico de Tráfico de Red `CU-05`) | `etapa3-pendientes.test.tsx` (bloque `FIX-131`) | ✅ **Pasa (3/3)** | Implementado en PR #182 (`fa92405` / `26af24f`): `InstanceMetricPoint` declara `netInBps`/`netOutBps`, `instanceDetailService.ts` los normaliza e `InstanceResourceUsage.tsx` renderiza las 4 tarjetas (`CPU`, `Memoria`, `Disco`, `Tráfico de Red`). |
| `FIX-127` (`SnapshotsTabContent` + `:instanceRef`/`:vmid`) | `etapa3-pendientes.test.tsx` (bloque `FIX-127`) | ✅ **Pasa (3/3)** | Implementado en `83a19b5`/`14db393`/`30a79a5`/`109f266`: `useOptionalEventsContext` y `instanceRef ?? vmid`. *(Queda `FIX-141` para sincronizar tests internos del submódulo `frontend/centinela/src/test/`).* |
| `FIX-132` (Errores 400 de contraseña no queman OTP) | `etapa3-pendientes.test.tsx` (bloque `FIX-132`), `recover-password.test.tsx` | ✅ **Pasa (4/4 y 7/7)** | Implementado en `3d06236`: `SAME_PASSWORD`, `WEAK_PASSWORD` y `PASSWORD_RECENTLY_USED` no incrementan `failedAttempts`. |
| `FRN-32` (Menú desplegable de exportación CSV/JSON) | `etapa3-pendientes.test.tsx` (bloque `FRN-32`) | ✅ **Pasa (4/4)** | Implementado en `404c784`: `AuditExportButton` unificado con estado de carga, Toast y transferencia de filtros. |
| `FIX-128` (`FRN-31` Auditoría Etapa 3) | `frn-31-audit-etapa3.test.tsx` | ✅ **Pasa (4/4)** | Implementado en `b38ab17` y `408043e`: acciones de Etapa 3 (`SNAPSHOT_CREATE`, `SNAPSHOT_ROLLBACK`, `SNAPSHOT_DELETE`, `ALERT_TRIGGERED`, `CONFIG_UPDATED`), filtro por `recursoTipo`, visor de detalle JSON y exportación CSV/JSON. |
| `FRN-33` (Centro de notificaciones y banners `CRITICAL`) | `frn-33-alerts-center.test.tsx` | ❌ **Falla (1 de 2)** | **Pendiente de implementación:** falta el banner rojo contextual (`role="alert"`) en la cabecera del Dashboard ante eventos `ALERT_TRIGGERED` de saturación crítica. |
| `FRN-37` / `ADP-07` (Plantillas Debian 12 Bookworm en Wizard) | `frn-37-adp-07-wizard-debian12.test.tsx` | ❌ **Falla (1 de 3)** | **Pendiente de implementación:** `fetchCreationTemplates` devuelve `'—'` en `operatingSystem` para `debian-12-standard_12.7-1_amd64.tar.zst` en lugar de extraer dinámicamente `"Debian"` y versión `"12"`. |
| `FIX-144` (Alineación con enum cerrado de auditoría `RF-08`) | `etapa3-pendientes.test.tsx` (bloque `FIX-144`) | ❌ **Falla (2 de 2)** | **Pendiente de implementación (depende de `FIX-143`):** `AUDIT_ACTION_OPTIONS` contiene 4 códigos inexistentes en `ports/audit_port.go` (`START_VM`, `STOP_VM`, `CAMBIAR_CONTRASENA`, `EXPORTAR_AUDITORIA`), y `Auditoria.tsx` (`SUCCESS`/`FAILED`) vs `Usuarios.tsx` (`EXITO`) consultan familias distintas de `resultado`. |

## Regresión (Etapas 1, 2 y 3)

Pasan **44 de los 47 archivos ejecutados en local** (296 pruebas en verde), incluyendo:
- **Etapa 1:** `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client`, `navigation`, `admin-users`, `admin-users-deletion`, `admin-recovery`, `recover-password`, `password-change`, `session-security`, `audit`, `auth-branding`, `events-client`, `dashboard-metrics`, `dashboard-node`, `instances-table`, `instances-modals`, `instances-operation`, `instances-ola2`, `instances-metrics`, `instances-events`, `dashboard-summary`, `instance-detail`, `fixes-pending`, `adversarial-security`.
- **Etapa 2:** `frn-22-use-instance-metrics` (9/9), `frn-23-metrics-charts` (9/9), `frn-24-25-create-wizard` (11/11), `frn-26-create-submit-quota-sse` (4/4), `frn-27-28-edit-instance-config` (6/6), `fix-91-dead-code` (2/2), `fix-103-node-history` (7/7), `fix-104b-105-tareaid` (12/12), `fix-108-ids-ofuscados` (11/11), `fix-111b-cambio-email-reauth` (4/4), `fix-115b-contrasenas-largas` (11/11), `fix-117b-dashboard-operador` (4/4).
- **Etapa 3:** `frn-36-29-snapshots-ui` (7/7), `frn-30-snapshots-rollback-ui` (2/2), `frn-31-audit-etapa3` (4/4), `frn-34-cookies-a11y` (3/3) y los bloques `FIX-127`, `FIX-131`, `FIX-132` y `FRN-32` en `etapa3-pendientes.test.tsx` (14/16).

## Suites internas del submódulo (`frontend/centinela/src/test/`) — Seguimiento en `FIX-141`

Al ejecutar `npm --prefix frontend/centinela test -- --run` dentro del submódulo:
- **Archivos de prueba:** 43 (**38 en verde**, **5 con fallos**).
- **Pruebas individuales:** 272 (**260 en verde**, **12 fallidas**):
  1. `src/test/snapshots.test.tsx` (4 fallos) y `src/test/instance-transient-states.test.tsx` (1 fallo): falta exportar `useOptionalEventsContext` en el `vi.mock('@/context/EventsContext')` manual.
  2. `src/test/instances-table.test.tsx` (3 fallos): el commit `99c7a18` eliminó las columnas `Checkbox` e `ID` de `Instances.tsx` sin actualizar este test interno.
  3. `src/test/dashboard-node-status.test.tsx` (3 fallos): los botones `"Ver todas"` y `"Ver todo"` pasaron a llamarse `"Ver todas las instancias"` y `"Ver toda la actividad"` en el Dashboard.
  4. `src/test/accessibility.test.tsx` (1 fallo): `button.tsx` reemplazó `focus-visible:ring-4` por `focus-visible:ring-0 focus-visible:bg-[#93c5fdbd]`.
