# Pruebas de aceptación del frontend

Estas pruebas contrastan el frontend (`frontend/centinela`, React + TS) con los criterios de éxito de `documentacion/actual.md` y funcionan como prueba de regresión de `documentacion/terminado.md`, `documentacion/terminado2.md` y `documentacion/terminado3.md`. Usan Vitest, React Testing Library y jsdom, sin agregar dependencias ni archivos al submódulo.

## Principios

- **Una prueba que pasa verifica el criterio de éxito.** Se exige la llamada HTTP real (método, ruta, cabeceras, payload) y su efecto en la interfaz. No se usan patrones del tipo `if (llamada) … else expect(boton).toBeVisible()`.
- **Los mocks de `fetch` reproducen el contrato real del backend**, tomado de `backend/internal/adapters/primary/http`. Ejemplos: `/auth/password/forgot`, `{ permisos: [{ vmid, nivelAcceso }] }` en permisos, `GET /instances` con `{ id, name, type: vm|lxc }` y la ausencia de `contrasenaTemp` por BAC-16. Si el frontend espera otro contrato, la prueba falla: es un problema de integración.
- **La sesión se siembra donde la guarda la app.** `centinela_access` y `centinela_pending_login` van en sessionStorage; `centinela_user` va en localStorage.

## Cobertura

| Archivo | Tareas |
|---|---|
| `login-form.test.tsx` | FRN-01, FRN-02 |
| `authentication-contract.test.ts` | FRN-04, LOGIN-02, LOGIN-03 |
| `two-factor-form.test.tsx` | LOGIN-02 |
| `two-factor-enrollment.test.tsx` | FRN-09 |
| `api-client.test.ts` | FIX-07 |
| `navigation.test.tsx` | FRN-03, **SEC-03** |
| `admin-users.test.tsx` | FRN-05, FRN-06, FRN-06B, **FIX-14/FRN-07**, FRN-08 |
| `admin-users-deletion.test.tsx` | **FIX-56** (distinguir suspensión de eliminación) |
| `admin-recovery.test.tsx` | **FRN-11** |
| `recover-password.test.tsx` | **FRN-12** |
| `password-change.test.tsx` | **FRN-10** |
| `session-security.test.ts` | **FRN-13**, **SEC-02**, FIX-08 |
| `audit.test.tsx` | **FRN-14** |
| `auth-branding.test.tsx` | **FIX-58**, **FIX-59** |
| `events-client.test.ts` | **FRN-17A**, **FRN-17C/FIX-61** |
| `dashboard-metrics.test.tsx` | **FRN-19A** (medidores CPU/RAM/almacenamiento, bordes 69% y 70%) |
| `dashboard-node.test.tsx` | **FRN-19B** (semáforo de salud global, 502/504, FIX-60) |
| `instances-table.test.tsx` | **FRN-20A** (tabla de inventario, normalización y RBAC) |
| `instances-modals.test.tsx` | **FRN-15**, **FIX-76**, **FIX-85**, **FIX-102** (modales de confirmación antierror, desacople nivelAcceso ADMIN y cierre de Start) |
| `instances-operation.test.tsx` | **FRN-16**, **FIX-86** (máquina de estados operación en progreso, errores D2 y mapeo 403) |
| `instances-ola2.test.tsx` | **FRN-20B**, **FIX-79**, **FIX-101** (filtros reactivos, selector visible Todas/VM/LXC, buscador dinámico por nombre/IP/ID) |
| `instances-metrics.test.tsx` | **FRN-20C**, **FIX-97** (columnas de CPU en porcentaje 0-100, RAM y em-dash '—' en nulos) |
| `instances-events.test.tsx` | **FRN-17B** (desbloqueo reactivo y toast adaptativo ante TASK_FINISHED) |
| `dashboard-summary.test.tsx` | **FRN-19C**, **FIX-83** (tarjetas de resumen VMs/LXC y refresco por evento) |
| `instance-detail.test.tsx` | **FRN-21**, **FIX-95** (vista Detalle de Instancia, registro en AppRoutes y enlace desde inventario) |
| `fixes-pending.test.tsx` | **FIX-87**, **FIX-89 / FIX-99**, **FIX-90**, **FIX-91**, **FIX-100** |
| `frn-22-use-instance-metrics.test.tsx` | **FRN-22** (hook `useInstanceMetrics`, ticket efímero WebSocket y reconexión con backoff exponencial) |
| `frn-23-metrics-charts.test.tsx` | **FRN-23**, **FIX-131** (gráficos CPU, Memoria, Red y Disco I/O, selector `2h`/`24h`/`7d`/`En vivo` y cartel `"Instancia detenida"`) |
| `frn-24-25-create-wizard.test.tsx` | **FRN-24**, **FRN-25**, **FIX-104** (wizard de creación de instancias, catálogo de plantillas y validación reactiva de cuotas) |
| `frn-26-create-submit-quota-sse.test.tsx` | **FRN-26** (envío de `POST /api/instances`, alerta `409 RESOURCE_QUOTA_EXCEEDED`, estado `"Creando..."` y evento SSE) |
| `frn-27-28-edit-instance-config.test.tsx` | **FRN-27**, **FRN-28** (modal de edición de vCPU y RAM, advertencia hot-plug/reboot y refresco en Detalle de Instancia) |
| `frn-36-29-snapshots-ui.test.tsx` | **FRN-36**, **FRN-29** (pestaña Snapshots, estados transitorios `SNAPSHOTTING`/`ROLLING_BACK`, creación, rollback y eliminación) |
| `frn-31-audit-etapa3.test.tsx` | **FRN-31** (filtros de auditoría de Etapa 3, badges semánticos, modal JSON de detalles y exportación CSV/JSON) |
| `frn-33-alerts-center.test.tsx` | **FRN-33** (centro de notificaciones de alertas críticas en cabecera, banner rojo y limpieza en `ALERT_RESOLVED`) |
| `frn-37-adp-07-wizard-debian12.test.tsx` | **FRN-37 / ADP-07** (soporte de plantillas Debian 12 en el Wizard de creación de instancias) |
| `etapa3-pendientes.test.tsx` | **FRN-32**, **FIX-127**, **FIX-128**, **FIX-132**, **FIX-144** (persistencia de sesión en memoria, estados transitorios, VLAN 0, y alineación de vocabulario de auditoría) |
| `fix-103-node-history.test.tsx` | **FIX-103** (historial de telemetría del nodo) |
| `fix-104b-105-tareaid.test.tsx` | **FIX-104B**, **FIX-105** (correlación de `tareaId` en eventos SSE) |
| `fix-108-ids-ofuscados.test.tsx` | **FIX-108** (consumo y validación estricta `/^[a-zA-Z0-9_-]+$/` de identificadores ofuscados `inst_...`) |
| `fix-111b-cambio-email-reauth.test.tsx` | **FIX-111B** (reautenticación al cambiar email del propio usuario) |
| `fix-115b-contrasenas-largas.test.tsx` | **FIX-115B** (soporte de contraseñas de hasta 128 caracteres en formularios) |
| `fix-117b-dashboard-operador.test.tsx` | **FIX-117B**, **FIX-142** (Dashboard para rol `OPERATOR` con tarjeta de almacenamiento visible y resumen acotado) |
| `adversarial-security.test.ts` | **Seguridad adversarial** (prevención de XSS en errores, inyección de rutas en API client, resiliencia a JSON corrupto, prototype pollution) |
| `login04-e2e.test.ts` | **LOGIN-04** (circuito E2E completo contra backend real en Docker) |
| `int02-e2e.test.ts` | **INT-02** (Smoke test integral de 7 pasos con código real de frontend contra el backend real) |

## Ejecutar

```bash
pnpm --dir test/front install
pnpm --dir test/front test                 # contra el submódulo tal como está checkouteado
```

Antes de correr la suite se trae el último commit del submódulo; ante un conflicto prevalece el remoto:

```bash
git -C frontend fetch origin --prune
git -C frontend checkout -B main origin/main --force
git -C frontend reset --hard origin/main
```

`vitest.config.mjs` también acepta `CENTINELA_FRONTEND_DIR` para apuntar a otra copia de `centinela/`, pero el procedimiento estándar es correr la suite sobre el submódulo actualizado.

Si cambian las dependencias del frontend, primero hay que correr `pnpm install` en `frontend/centinela`.

## Entorno

`setup.ts` agrega los polyfills que jsdom no trae: `document.elementFromPoint` y `window.matchMedia` (este último lo usa `hooks/use-mobile.ts` del sidebar). Después de cada prueba limpia sessionStorage y localStorage.

Resultado vigente: [`RESULTADOS.md`](RESULTADOS.md). Informe consolidado: [`../informe.md`](../informe.md).
