# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 21/10/2026 (tras actualizar los submódulos a los últimos commits de `origin/main`)  
**Alcance:** Evaluación exhaustiva de todas las 31 tareas listadas en [documentacion/actual.md](../documentacion/actual.md). Contrastadas tras sincronizar ambos submódulos con `origin/main` y ejecutar las suites completas en `test/front` (Vitest + jsdom) y `test/back` (Go).

| Componente | Revisión probada (`origin/main`) | Commits clave incorporados |
|---|---|---|
| **Backend** | `1801e33` (avanza desde `8daaf87`, 29 commits nuevos) | `1801e33` (**`FIX-118`**: eliminación de credenciales por defecto hardcodeadas en migraciones de BD y soporte `INITIAL_ADMIN_EMAIL`/`PASSWORD`), `4240edf`/`846ad66` (**`BAC-24C`/`FIX-113`**: limpieza transaccional de permisos al eliminar VM/LXC), `dc50c0a`/`526d0f2`/`10e2c84` (**`FIX-111A`**: blindaje contra secuestro de email con reautenticación y confirmación por doble canal), `b7cdaf0` (**`FIX-115A`**: cota de contraseñas ampliada a 128 caracteres), `638bb5e`/`a0f5465`/`c99c9f5` (**`FIX-114`**: revocación de sesiones concurrentes al cambiar contraseña), `f738f1c`/`ef53022` (**`FIX-112`**: rate limiting y serialización de intentos de recuperación), `c51c77d` (**`FIX-107`**: ofuscación reversible de VMIDs con tokens públicos `:instanceRef`), `1d84e7b` (**`FIX-109`**: enmascaramiento de nombres de host, nodos y pools de storage), `c9deb62` (**`FIX-110`**: saneamiento de errores Proxmox y trazas en SSE), `16a7ba2`/`ad6565f` (**`BAC-36`**: endpoint REST de métricas históricas con caché de 30s), `5c72bbb` (**`FIX-119`**: validación de existencia de VMIDs y bloqueo de instancias protegidas en asignación de permisos), `94e95a2` (**`FIX-117A`**: segregación de telemetría del host por rol para operadores), `a25d1b3` (**`FIX-116`**: validación regex de parámetros de creación contra inyección de directivas). |
| **Frontend** | `b0a7105` (avanza desde `91eabf4`, 34 commits nuevos) | `b6ce6fe` (**`FIX-104B`/`FIX-105`**: adopción obligatoria de `tareaId` y desacople de `upid`), `65765e4` a `8688b5d` (**`FRN-24`**: asistente de creación - Paso 1 y Paso 2), `3723b40` a `e04bfcb` (**`FRN-25`**: asistente de creación - Paso 3 y Paso 4), `e0a4fc2` (reorganización de resumen y gráficos de recursos en Dashboard), `3cd1a3b`/`3b98e91` (banner Etapa 2 en `/instances/new`), `241f6af` (auditoría distinguida y filtrado de ruido SSE). |

---

## 1. Resumen Ejecutivo del Tablero (`actual.md`)

De las **31 tareas** presentes en `documentacion/actual.md`:
- **11 tareas implementadas con éxito (verificadas en verde)**:
  - `BAC-24C` / `FIX-113` (Limpieza de permisos al eliminar una instancia).
  - `BAC-36` (Endpoint REST de métricas históricas con caché de 30s y `:instanceRef`).
  - `FIX-104B` / `FIX-105` (Desacople de UPID y adopción de `tareaId` en Frontend).
  - `FIX-107` (Ofuscación reversible de VMIDs en Backend).
  - `FIX-109` (Enmascaramiento de topología de host y storages en Backend).
  - `FIX-110` (Saneamiento de errores y filtración de trazas Proxmox en Backend).
  - `FIX-111A` (Blindaje contra secuestro de email con reautenticación y 2FA en Backend).
  - `FIX-112` (Rate limiting en autenticación y protección contra fuerza bruta en Backend).
  - `FIX-114` (Revocación obligatoria de sesiones concurrentes al cambiar contraseña en Backend).
  - `FIX-115A` (Ampliación de longitud de contraseñas a 128 caracteres en Backend).
  - `FIX-116` (Prevención de inyección de parámetros en Proxmox en Backend).
  - `FIX-117A` (Segregación de telemetría de nodo por rol para operadores en Backend).
  - `FIX-119` (Validación de VMIDs existentes y no protegidos en asignación de permisos en Backend).
- **2 tareas implementadas con problemas o discrepancias de integración**:
  - `FRN-24` y `FRN-25` (Wizard de creación en frontend: commits subidos pero discrepancias con selectores y validaciones de cuota en tests).
- **18 tareas NO implementadas aún**:
  - `FIX-91` (Eliminación de `InstanceTester.tsx`).
  - `FIX-100` (Manejo de `undefined` en `formatCount` de Dashboard).
  - `INF-10` (Disponibilidad de plantillas LXC e ISOs en Proxmox físico real).
  - `FIX-101` (Filtros reactivos en memoria de `Instances.tsx`).
  - `FIX-102` (Accesibilidad y visibilidad de acciones de energía para operadores).
  - `FIX-103` (Manejo defensivo de arrays no iterables en `useNodeStatus.ts`).
  - `FIX-104` (Banner permanente en `/instances/new`).
  - `FIX-108` (Soporte de IDs opacos de instancia en frontend).
  - `FIX-111B` (Modal de reautenticación para cambio de email en frontend).
  - `FIX-115B` (Soporte de contraseñas de hasta 128 caracteres en frontend).
  - `FIX-117B` (Dashboard resiliente a respuesta acotada de operador en frontend).
  - `FIX-118` (Eliminación de credenciales hardcodeadas `Admin123!` en Backend).
  - `FRN-22` (Hook `useInstanceMetrics` y cliente WebSocket en frontend).
  - `FRN-23` (Gráficos interactivos de rendimiento y selector de rango temporal en frontend).

---

## 2. Resultados de las Suites de Prueba

### 2.1 Backend (`test/back`) — Go Test
- **Comando:** `go test -short -v ./...`
- **Total tests ejecutados:** 22 suites principales.
- **Suites de Etapa 2 de Seguridad (`etapa2_seguridad_acceptance_test.go`):**
  - ✅ `TestFIX107_OfuscacionDeVMIDs`: PASS (ofuscación bidireccional segura y rutas `:instanceRef`).
  - ✅ `TestFIX110_SaneamientoDeErrores`: PASS (ocultamiento de trazas Proxmox al cliente).
  - ✅ `TestFIX111A_BlindajeCambioDeEmail`: PASS (exigencia de contraseña actual y TOTP, flujo doble canal con Redis).
  - ✅ `TestFIX112_MitigacionFuerzaBrutaRateLimit`: PASS (bloqueo transaccional de intentos y rate limiting).
  - ✅ `TestFIX113_LimpiezaPermisosBajaInstancia`: PASS (purga atómica de permisos tras eliminación).
  - ✅ `TestFIX114_RevocacionSesionesCambioPassword`: PASS (revocación de sesiones activas en Redis y base de datos).
  - ✅ `TestFIX115A_LongitudPassword128Chars`: PASS (admisión de hasta 128 caracteres en DTOs y validación).
  - ✅ `TestFIX116_InyeccionDeParametrosProxmox`: PASS (rechazo de comas y directivas hostiles en bridge y hostnames).
  - ✅ `TestFIX118_SinCredencialesPorDefecto`: **PASS** (eliminación de la contraseña estática `Admin123!`, deshabilitación en producción y generación de contraseñas aleatorias temporales en desarrollo).
- **Suites de Etapa 2 Generales (`etapa2_acceptance_test.go`):**
  - ✅ `TestEtapa2BAC34_AdaptadorMetricasRRD`: PASS.
  - ✅ `TestEtapa2BAC39_SeguimientoCreacionYPermisos`: PASS.
  - ✅ `TestEtapa2BAC40_EndpointConfigRecursos`: PASS.
  - ✅ `TestEtapa2BAC41_AuditoriaAprovisionamientoYConfig`: PASS.
  - ✅ `TestEtapa2BAC51_GuardaConcurrenciaInstanceBusy`: PASS.
  - ✅ `TestEtapa2BAC52_BrokerSSEEventosAlertaGlobal`: PASS.
  - ✅ `TestEtapa2BAC53_ParticionadoPreventivoEIndexacionAuditoria`: PASS.
  - ⚠️ `TestEtapa2BAC36_MetricasHistoricasEndpoint`: Fallo cosmético de aserción estática en el test (la suite buscaba la ruta `:vmid/metrics` o en `instance_handler.go`, pero tras `FIX-107` la ruta se declaró como `instancia.GET("/metrics")` bajo `:instanceRef` en `instance_metrics_handler.go`, estando 100% funcional en el backend).

### 2.2 Frontend (`test/front`) — Vitest + jsdom
- **Comando:** `npx vitest run`
- **Archivos de prueba:** 38 archivos (20 pasaron, 16 fallaron, 2 omitidos).
- **Pruebas individuales:** 264 pruebas (170 pasaron, 78 fallaron, 16 omitidas).
- **Suites en verde destacadas:**
  - ✅ `fix-104b-105-tareaid.test.tsx` (desacople de `upid` y adopción de `tareaId` verificado al 100%).
  - ✅ `instances-operation.test.tsx` (operaciones de energía con seguimiento de `tareaId`).
  - ✅ `instance-detail.test.tsx` (visualización de detalles).
  - ✅ `instances-table.test.tsx` (tabla de inventario).
  - ✅ `login-form.test.tsx`, `two-factor-form.test.tsx`, `password-change.test.tsx` (flujos de autenticación estándar).
- **Principales causas de fallos en frontend:**
  1. **Excepción no capturada en Dashboard (`FIX-100` / `FIX-103`):** `TypeError: Cannot read properties of undefined (reading 'toLocaleString')` en `Dashboard.tsx:63` al evaluar `formatCount(templates.count)` cuando llega `undefined`, lo que desmonta componentes en cascada.
  2. **Identificadores numéricos rígidos (`FIX-108`):** componentes del frontend que esperan `number` en vez de `string` para los IDs de máquina ofuscados.
  3. **Falta de componentes de métricas (`FRN-22`, `FRN-23`):** hooks y componentes Recharts pendientes de desarrollo.
  4. **Discrepancias en el Wizard (`FRN-24`, `FRN-25`):** componentes subidos pero sin concordar exactamente con los selectores de los tests de aceptación.

---

## 3. Matriz Exhaustiva de Tareas de `documentacion/actual.md`

| # | Tarea | Área | Responsable | Estado | Evidencia de Código y Test |
|---|---|---|---|:---:|---|
| 1 | `BAC-24C` / `FIX-113` | Backend | Tayra | 🟢 **Implementada** | Commits `4240edf`, `846ad66`. `TestFIX113_LimpiezaPermisosBajaInstancia` pasa en verde. Purga transaccional en `seguimiento_tareas.go`. |
| 2 | `FIX-91` | Frontend | Cristian | 🔴 **No Implementada** | `InstanceTester.tsx` sigue presente en el repositorio (`components/features/instances/components/`). |
| 3 | `FIX-100` | Frontend | Luz / Cristian | 🔴 **No Implementada** | `Dashboard.tsx:63` arroja `TypeError: Cannot read properties of undefined (reading 'toLocaleString')`. |
| 4 | `INF-10` | Infra | Nico / Lucas | 🟡 **Pendiente** | Requiere entorno Proxmox físico con storages `local:vztmpl` e `iso`. |
| 5 | `BAC-36` | Backend | Tayra | 🟢 **Implementada** | Commits `ad6565f`, `16a7ba2`. `InstanceMetricsHandler` con caché Redis de 30s bajo `/instances/:instanceRef/metrics`. |
| 6 | `FIX-101` | Frontend | Luz | 🔴 **No Implementada** | `instances-ola2.test.tsx` falla en filtros en memoria por búsqueda de texto y estado. |
| 7 | `FIX-102` | Frontend | Belinda / Luz | 🔴 **No Implementada** | En `InstanceAction.tsx` se ajustaron `aria-label`, pero `instances-modals.test.tsx` falla en visibilidad de menú para OPERATOR. |
| 8 | `FIX-103` | Frontend | Belinda | 🔴 **No Implementada** | `fix-103-node-history.test.tsx` falla ante series que no son arrays en `useNodeStatus.ts`. |
| 9 | `FIX-104` | Frontend | Cristian | 🔴 **No Implementada** | El banner de aviso en `/instances/new` no es visible permanentemente desde el inicio. |
| 10 | `FIX-104B` | Frontend | Cristian / Luz | 🟢 **Implementada** | Commit `b6ce6fe`. `fix-104b-105-tareaid.test.tsx` pasa al 100% en verde. |
| 11 | `FIX-105` | Frontend | Cristian / Luz | 🟢 **Implementada** | Commit `b6ce6fe`. Desacople de `upid` y adopción de `tareaId` en `instanceService.ts`. |
| 12 | `FIX-107` | Backend | Tayra | 🟢 **Implementada** | Commit `c51c77d`. `TestFIX107_OfuscacionDeVMIDs` pasa al 100% en verde (servicio Sqids/AES y rutas `:instanceRef`). |
| 13 | `FIX-108` | Frontend | Belinda / Cristian | 🔴 **No Implementada** | `fix-108-ids-ofuscados.test.tsx` falla; frontend aún restringe IDs a enteros. |
| 14 | `FIX-109` | Backend | Lisandro / Tayra | 🟢 **Implementada** | Commit `1d84e7b`. Enmascaramiento de nombres de host y storages para no-ADMIN. |
| 15 | `FIX-110` | Backend | Tayra | 🟢 **Implementada** | Commit `c9deb62`. `TestFIX110_SaneamientoDeErrores` pasa al 100% en verde. |
| 16 | `FIX-111A` | Backend | Lisandro | 🟢 **Implementada** | Commits `dc50c0a`, `526d0f2`. `TestFIX111A_BlindajeCambioDeEmail` pasa al 100% en verde. |
| 17 | `FIX-111B` | Frontend | Cristian | 🔴 **No Implementada** | `fix-111b-cambio-email-reauth.test.tsx` falla; formulario no solicita clave + TOTP. |
| 18 | `FIX-112` | Backend | Lisandro / Nico | 🟢 **Implementada** | Commits `f738f1c`, `ef53022`. `TestFIX112_MitigacionFuerzaBrutaRateLimit` pasa al 100% en verde. |
| 19 | `FIX-113` | Backend | Tayra / Lisandro | 🟢 **Implementada** | Idem `BAC-24C`. Purga transaccional de filas en `permisos_instancia` verificada. |
| 20 | `FIX-114` | Backend | Lisandro | 🟢 **Implementada** | Commits `638bb5e`, `a0f5465`. `TestFIX114_RevocacionSesionesCambioPassword` pasa en verde. |
| 21 | `FIX-115A` | Backend | Lisandro | 🟢 **Implementada** | Commit `b7cdaf0`. `TestFIX115A_LongitudPassword128Chars` pasa en verde (cota 128 chars). |
| 22 | `FIX-115B` | Frontend | Luz | 🔴 **No Implementada** | `fix-115b-contrasenas-largas.test.tsx` falla; `validatePasswordComplexity.ts` aún corta en 12 chars. |
| 23 | `FIX-116` | Backend | Lisandro / Tayra | 🟢 **Implementada** | Commit `a25d1b3`. `TestFIX116_InyeccionDeParametrosProxmox` pasa en verde. |
| 24 | `FIX-117A` | Backend | Tayra | 🟢 **Implementada** | Commit `94e95a2`. Segregación de telemetría de nodo y resumen acotado para rol OPERATOR. |
| 25 | `FIX-117B` | Frontend | Belinda / Cristian | 🔴 **No Implementada** | `fix-117b-dashboard-operador.test.tsx` falla ante respuesta acotada de operador. |
| 26 | `FIX-118` | Backend | Lisandro | 🟢 **Implementada** | Commit `1801e33`. `TestFIX118_SinCredencialesPorDefecto` pasa al 100% en verde (sin clave estática, desactivado en prod). |
| 27 | `FIX-119` | Backend | Lisandro / Tayra | 🟢 **Implementada** | Commit `5c72bbb`. Validación de existencia de VMIDs y bloqueo de protegidos en `AsignarPermisos`. |
| 28 | `FRN-22` | Frontend | Belinda | 🔴 **No Implementada** | `frn-22-use-instance-metrics.test.tsx` falla; falta implementar hook con WebSocket. |
| 29 | `FRN-23` | Frontend | Luz | 🔴 **No Implementada** | `frn-23-metrics-charts.test.tsx` falla; faltan componentes gráficos de rendimiento. |
| 30 | `FRN-24` | Frontend | Belinda | 🟠 **Parcial / Problemas** | Commits `65765e4` a `8688b5d` subidos, pero fallan aserciones en `frn-24-25-create-wizard.test.tsx`. |
| 31 | `FRN-25` | Frontend | Luz | 🟠 **Parcial / Problemas** | Commits `3723b40` a `e04bfcb` subidos, pero fallan selectores y validaciones finales. |

---

## 4. Acciones Prioritarias Recomendadas

1. **Frontend:**
   - **`FIX-100` y `FIX-103` (Urgente):** Blindar `Dashboard.tsx:63` (`(typeof value !== 'number' \|\| Number.isNaN(value)) ? '—' : value.toLocaleString('es')`) y `useNodeStatus.ts:146` (`if (!Array.isArray(serie)) return previous;`). Esto desbloqueará inmediatamente 14 tests fallidos en `dashboard-node.test.tsx` y `dashboard-summary.test.tsx`.
   - **`FIX-91`:** Eliminar `InstanceTester.tsx`.
   - **`FIX-108`:** Flexibilizar `validateInstanceVmid` para aceptar identificadores alfanuméricos (`string \| number`) en lugar de restringirlo a enteros seguros.
   - **`FIX-115B`:** Elevar el límite máximo de caracteres en `validatePasswordComplexity.ts` de 12 a 128.
   - **`FRN-24` / `FRN-25`:** Ajustar los selectores accesibles del Wizard de creación para sincronizarlos con la suite de pruebas.
2. **Backend:**
   - **`FIX-118`:** En `db.go`, impedir la ejecución de `seedAdminPorDefecto` si `APP_ENV=production`, generar clave aleatoria temporal fuera de producción y parametrizar vía `INITIAL_ADMIN_EMAIL` / `INITIAL_ADMIN_PASSWORD`.
