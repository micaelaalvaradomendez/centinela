# Cruce de Tareas vs Commits — Submódulo Frontend (Parte C: 29/10/2026 al 03/11/2026)

**Fecha de generación:** 03/11/2026  
**Submódulo:** `frontend/`  
**Rango analizado (Parte C):** `443e97a` (2026-10-29, commit #668) → `877eddd` (2026-11-02, commit #708, HEAD actual, PR #181)  
**Corte previo (Parte B):** `594a4e0` (2026-10-30, commit #667, Merge PR #168 — `FRN-30`)  
**Total de commits en el tramo Parte C:** 41 commits (12 commits de merge y 29 commits de código reales)  
**Total acumulado en el repositorio:** 708 commits (`c04eddb` → `877eddd`)

---

## 1. Metodología y Continuidad con Partes A y B

Este documento continúa la trazabilidad histórica del submódulo `frontend/`:
1. **Parte A (`cruce-tareas-commits-frontend.md` / `cruce-tareas-commits-frontend1.md` — 20/10/2026):** Cubrió los primeros **536 commits** (`#1 c04eddb` al `#536 91eabf4`, del 17/08/2026 al 18/10/2026), abarcando `FRN-01` a `FRN-29` y `FIX-01` a `FIX-82`.
2. **Parte B (`cruce-tareas-commits-frontend-parte-b.md` — 30/10/2026):** Cubrió **131 commits** (`#537 3b98e91` al `#667 594a4e0`, del 19/10/2026 al 30/10/2026, hasta el PR #168), abarcando los cierres e iteraciones de `FRN-30`, `FRN-31`, `FRN-33`, `FRN-34`, `FRN-35`, `FRN-36` y los fixes `FIX-88` a `FIX-124`.
3. **Parte C (este informe — 03/11/2026):** Analiza los **41 commits posteriores** (`#668 443e97a` al `#708 877eddd`, PRs #169 al #181, del 29/10/2026 al 02/11/2026) y los cruza contra las tareas de adaptación a PVE 8.4 (`FRN-37`/`ADP-07` y `ADP-08`), funcionalidades pendientes de Etapa 3 (`FRN-32`, `FRN-33`, `FRN-34`, `FRN-36`) y correcciones (`FIX-126`, `FIX-127`, `FIX-128`, `FIX-131`, `FIX-132`) documentadas en `documentacion/actual.md`, `documentacion/tareas-ajustes-etapa3.md` y `documentacion/servidorfinal/01_tareas_adaptacion_pve8.md`.

---

## 2. Panorama General de Commits por Integrante

### 2.1 Tramo Parte C (`#668` al `#708` — 29/10/2026 al 02/11/2026)

| Integrante | Usuario / Email en Git | Commits Totales | Merges | Commits de Código Real | % del Tramo (Código) |
|---|---|---:|---:|---:|---:|
| **Nicolás** | `Nico` / `Nicolas Butterfield` (`nicolasbutterfield02@gmail.com`) | 23 | 6 | 17 | 58.6% |
| **Cristian** | `Cristian` / `Cristianusable` (`cristian.e.pacheco@gmail.com`) | 12 | 3 | 9 | 31.0% |
| **Belinda** | `Belinda` / `BelindaMeriles` (`belindameriles@gmail.com`) | 6 | 3 | 3 | 10.3% |
| **Luz** | `luzpacello` (`luzpacello@gmail.com`) | 0 | 0 | 0 | 0.0% |
| **Total Parte C** | — | **41** | **12** | **29** | **100.0%** |

### 2.2 Acumulado Histórico Completo (`#1` al `#708` — 17/08/2026 al 02/11/2026)

| Integrante | Parte A (#1–#536) | Parte B (#537–#667) | Parte C (#668–#708) | Total Commits | Merges Totales | Código Real Total | % Código Real |
|---|---:|---:|---:|---:|---:|---:|---:|
| **Luz** | 252 | 9 | 0 | **261** | 43 | **218** | 40.3% |
| **Cristian** | 131 | 52 | 12 | **195** | 50 | **145** | 26.8% |
| **Belinda** | 146 | 41 | 6 | **193** | 67 | **126** | 23.3% |
| **Nicolás** | 7 | 29 | 23 | **59** | 7 | **52** | 9.6% |
| **Total Histórico** | **536** | **131** | **41** | **708** | **167** | **541** | **100.0%** |

---

## 3. Lista Cronológica de Commits del Tramo Parte C (#668 al #708)

| # | Fecha | Hash | Autor | Mensaje del Commit | Archivos Principales Modificados | PR / Tarea Asociada |
|---:|---|---|---|---|---|---|
| 668 | 2026-10-29 | `443e97a` | Nicolás | `fix(dashboard): ajustar medidor de almacenamiento 1dfdzjeajw9` | `StorageHealthCard.tsx`, `StorageHealthCard.test.tsx` | PR #169 → **FIX-117B / FIX-126** |
| 669 | 2026-10-30 | `b9a67a7` | Nicolás | `fix(accesibilidad): foco accesible en modales y contraste del ingreso (1dfdzjeajxy)` | `useModalFocus.ts`, `Login.tsx`, `SessionExpiryWarningModal.tsx`, `frn-34-session-renewal-security.test.tsx` | PR #170 → **FRN-34** |
| 670 | 2026-10-30 | `851a455` | Nicolás | `Merge pull request #169 from luzpacello/1dfdzjeajw9` | *(Merge PR #169)* | Merge **FIX-117B / FIX-126** |
| 671 | 2026-10-30 | `455aaee` | Nicolás | `Merge pull request #170 from luzpacello/1dfdzjeajxy` | *(Merge PR #170)* | Merge **FRN-34** |
| 672 | 2026-10-30 | `53125e1` | Cristian | `feat(instances): desacoplar plantillas de SO del asistente de creacion` | `CreateInstanceModal.tsx`, `instanceConstants.ts`, `proxmoxService.ts` | PR #171 → **FRN-37 / ADP-07** |
| 673 | 2026-10-30 | `5b28f63` | Cristian | `refactor(instances): ajustar selector de plantillas y pruebas del asistente` | `CreateInstanceModal.tsx`, `CreateInstanceModal.test.tsx` | PR #171 → **FRN-37 / ADP-07** |
| 674 | 2026-10-30 | `3b03bf4` | Cristian | `fix(instances): restaurar ayudas visuales en seleccion de plantillas` | `CreateInstanceModal.tsx` | PR #171 → **FRN-37 / ADP-07** |
| 675 | 2026-10-30 | `cdc4221` | Cristian | `test(instances): estabilizar pruebas del asistente de creacion` | `CreateInstanceModal.test.tsx` | PR #171 → **FRN-37 / ADP-07** |
| 676 | 2026-10-30 | `28f591c` | Cristian | `Merge pull request #171 from luzpacello/cristian` | *(Merge PR #171)* | Merge **FRN-37 / ADP-07** |
| 677 | 2026-10-30 | `0df9639` | Nicolás | `test(frontend): estabilizar suite interna de vitest` | `AuditLogsView.test.tsx`, `InstancesView.test.tsx`, `UsersView.test.tsx` | PR #172 → Estabilización Vitest |
| 678 | 2026-10-30 | `52eea52` | Nicolás | `fix(permisos): aceptar instanceRef de texto en la vista de permisos` | `InstancePermissionsView.tsx`, `InstancePermissionsView.test.tsx` | PR #172 → **FIX-108** |
| 679 | 2026-10-30 | `64df620` | Nicolás | `fix(instancias): cerrar acción de energía cuando llega TASK_FINISHED` | `useInstancePowerActions.ts`, `InstancesView.tsx` | PR #172 → **FRN-16 / FRN-36** |
| 680 | 2026-10-30 | `a2c6a93` | Nicolás | `fix(auditoria): unificar opciones de acción con el vocabulario de eventos` | `AuditLogsView.tsx`, `auditConstants.ts` | PR #172 → **FRN-31** |
| 681 | 2026-10-30 | `1b0cc16` | Nicolás | `fix(usuarios): alinear textos de baja y refresco con las pruebas de la vista` | `UsersView.tsx`, `DeleteUserModal.tsx` | PR #172 → **FIX-56** |
| 682 | 2026-10-30 | `29bd04c` | Nicolás | `fix(instancias): mantener visible el estado transitorio antes de cerrar con TASK_FINISHED` | `useInstancePowerActions.ts`, `InstanceStatusBadge.tsx` | PR #172 → **FRN-16 / FRN-36** |
| 683 | 2026-10-30 | `5cf6037` | Nicolás | `fix(instancias): esperar confirmación del listado antes de limpiar el estado en curso` | `useInstancePowerActions.ts`, `InstancesView.tsx` | PR #172 → **FRN-16 / FRN-26** |
| 684 | 2026-10-30 | `0561ccc` | Nicolás | `fix(instancias): cerrar estado en curso al recibir TASK_FINISHED sin esperar otro refresco` | `useInstancePowerActions.ts` | PR #172 → **FRN-16 / FRN-36** |
| 685 | 2026-10-30 | `cc0b943` | Nicolás | `Merge pull request #172 from luzpacello/1dfdzjeakg2` | *(Merge PR #172)* | Merge Estabilización / **FRN-16 / FRN-36** |
| 686 | 2026-10-30 | `d966ceb` | Nicolás | `fix(instancias): mostrar etiquetas en español en la columna de acciones` | `InstanceRow.tsx`, `InstancesTable.tsx` | PR #173 → **FRN-35 / FRN-36** |
| 687 | 2026-10-30 | `af9db56` | Nicolás | `fix(instancias): refrescar listado al cerrar tarea antes de limpiar el estado en curso` | `useInstancePowerActions.ts`, `InstancesView.tsx` | PR #173 → **FRN-16 / FRN-26** |
| 688 | 2026-10-30 | `6df975c` | Nicolás | `fix(instancias): esperar el cambio real de estado antes de limpiar la acción en curso` | `useInstancePowerActions.ts` | PR #173 → **FRN-16 / FRN-36** |
| 689 | 2026-10-30 | `4eb5fab` | Nicolás | `fix(instancias): sondear el listado tras TASK_FINISHED hasta ver el nuevo estado` | `useInstancePowerActions.ts`, `InstancesView.tsx` | PR #173 → **FRN-16 / FRN-26** |
| 690 | 2026-10-30 | `4299d09` | Nicolás | `fix(instancias): sostener el sondeo de estado hasta que Proxmox refleje el cambio` | `useInstancePowerActions.ts` | PR #173 → **FRN-16 / FRN-26** |
| 691 | 2026-10-30 | `6a8a2bf` | Nicolás | `Merge pull request #173 from luzpacello/1dfdzjeakl3` | *(Merge PR #173)* | Merge **FRN-16 / FRN-26 / FRN-36** |
| 692 | 2026-10-30 | `2e6cc83` | Nicolás | `fix(instances): mostrar el estado de apagado y reinicio en el badge` | `InstanceStatusBadge.tsx`, `InstancesTable.tsx` | PR #174 → **FRN-36** |
| 693 | 2026-10-30 | `667eef9` | Nicolás | `Merge pull request #174 from luzpacello/1dfdzjeakp8` | *(Merge PR #174)* | Merge **FRN-36** |
| 694 | 2026-10-30 | `101955a` | Nicolás | `fix(instances): sostener el estado en curso hasta que la pantalla se actualice` | `useInstancePowerActions.ts`, `InstancesView.tsx` | PR #176 → **FRN-16 / FRN-36** |
| 695 | 2026-10-30 | `4c16d05` | Nicolás | `Merge pull request #176 from luzpacello/1dfdzjeakqs` | *(Merge PR #176)* | Merge **FRN-16 / FRN-36** |
| 696 | 2026-11-01 | `8d66b75` | Belinda | `Verificación de Renderizado de Estado de Storage y Nodos en Dashboard` | `Dashboard.tsx`, `StorageHealthCard.tsx`, `ClusterHealthBanner.tsx` | PR #177 → **ADP-08** |
| 697 | 2026-11-01 | `ce38ea9` | Belinda | `Merge pull request #177 from luzpacello/Belinda` | *(Merge PR #177)* | Merge **ADP-08** |
| 698 | 2026-11-02 | `83a19b5` | Cristian | `fix(snapshots): soportar montaje aislado y alias de ruta en pestaña de snapshots` | `SnapshotsTabContent.tsx`, `InstanceDetailView.tsx` | PR #178 → **FIX-127** |
| 699 | 2026-11-02 | `14db393` | Cristian | `test(snapshots): cubrir montaje aislado de SnapshotsTabContent con ruta :instanceRef` | `SnapshotsTabContent.test.tsx` | PR #178 → **FIX-127** |
| 700 | 2026-11-02 | `30a79a5` | Cristian | `fix(snapshots): agregar aria-label explícito en acciones de snapshot` | `SnapshotsTabContent.tsx` | PR #178 → **FIX-127** |
| 701 | 2026-11-02 | `109f266` | Cristian | `fix(snapshots): evitar colisión de texto con selector getByText(/Restaurar/i)` | `SnapshotsTabContent.tsx`, `SnapshotsTabContent.test.tsx` | PR #178 → **FIX-127** |
| 702 | 2026-11-02 | `2d088e6` | Cristian | `Merge pull request #178 from luzpacello/cristian` | *(Merge PR #178)* | Merge **FIX-127** |
| 703 | 2026-11-02 | `3d06236` | Cristian | `fix(auth): no descontar intentos de OTP ante errores 400 de formato de contraseña` | `ForgotPasswordModal.tsx`, `ForgotPasswordModal.test.tsx` | PR #179 → **FIX-132** |
| 704 | 2026-11-02 | `2260b5b` | Cristian | `Merge pull request #179 from luzpacello/cristian` | *(Merge PR #179)* | Merge **FIX-132** |
| 705 | 2026-11-02 | `404c784` | Belinda | `Componente y flujo de exportación de reportes de auditoría` | `AuditExportModal.tsx`, `AuditLogsView.tsx`, `auditService.ts` | PR #180 → **FRN-32** |
| 706 | 2026-11-02 | `b3169d0` | Belinda | `Merge pull request #180 from luzpacello/Belinda` | *(Merge PR #180)* | Merge **FRN-32** |
| 707 | 2026-11-02 | `b38ab17` | Belinda | `FIX: FIX-128` | `CreateInstanceModal.tsx`, `EditInstanceNetworkModal.tsx`, `instanceValidation.ts` | PR #181 → **FIX-128** (Front) |
| 708 | 2026-11-02 | `877eddd` | Belinda | `Merge pull request #181 from luzpacello/Belinda` | *(Merge PR #181)* | Merge **FIX-128** (Front) |

---

## 4. Cruce Detallado por Tarea (Adaptación PVE 8.4 y Etapa 3)

### 4.1 Tareas de Adaptación al Servidor Real Proxmox VE 8.4 (`FRN-37`/`ADP-07` y `ADP-08`)

| Código | Descripción de la Tarea | Responsable Asignado (`01_tareas_adaptacion_pve8.md` / `actual.md`) | Autor(es) en Git | Commits / PR Asociados | Coincidencia de Autoría | Estado en Suite Externa (`test/front`) |
|---|---|---|---|---|---|---|
| **FRN-37 / ADP-07** | Desacoplamiento de constantes hardcodeadas de plantillas de SO (`Debian 13` / `Ubuntu 24.04`) en `CreateInstanceModal.tsx` para consumir dinámicamente el catálogo del nodo PVE 8.4 (`Debian 12 Bookworm`) | **Belinda** | **Cristian** | `53125e1`, `5b28f63`, `3b03bf4`, `cdc4221` (PR #171 `28f591c`) | ❌ **No coincide** *(Cristian)* | ✅ **Pasa** (`adp-pve8-frontend.test.tsx`) |
| **ADP-08** | Verificación de renderizado de estado de Storage (`local`, `local-lvm`) y salud del nodo (`pve`) en `Dashboard.tsx` | **Cristian / Belinda** | **Belinda** | `8d66b75` (PR #177 `ce38ea9`) | ✅ **Coincide** *(Belinda)* | ✅ **Pasa** (`adp-pve8-frontend.test.tsx`) |

### 4.2 Tareas Funcionales y Correcciones de Etapa 3 (`FRN-32` a `FRN-36` y `FIX-126` a `FIX-132`)

| Código | Descripción de la Tarea | Responsable Asignado (`actual.md`) | Autor(es) en Git (Parte C) | Commits / PR Asociados | Coincidencia de Autoría | Estado en Suite Externa (`test/front`) |
|---|---|---|---|---|---|---|
| **FIX-126 / FIX-117B** | Ajuste del medidor de almacenamiento e indicador de sobrecarga en `StorageHealthCard.tsx` | **Belinda / Cristian** | **Nicolás** | `443e97a` (PR #169 `851a455`) | ❌ **No coincide** *(Nicolás)* | ✅ **Pasa** (`fix-117b-storage-overload-indicator.test.tsx`) |
| **FRN-34** | Foco accesible en modales (`useModalFocus.ts`), contraste WCAG en login y modal de expiración de sesión | **Cristian** | **Nicolás** *(cierre en Parte C sobre base de Cristian en Parte B)* | `b9a67a7` (PR #170 `455aaee`) | ⚠️ **Compartida** *(Cristian + Nicolás)* | ✅ **Pasa** (`frn-34-session-renewal-security.test.tsx`) |
| **FRN-16 / FRN-26 / FRN-36 / FIX-108** | Estabilización de estados transitorios (`starting`, `stopping`, `rebooting`), cierre y sondeo con `TASK_FINISHED`, badges en español y soporte de `instanceRef` de texto en permisos | **Luz / Cristian / Belinda** | **Nicolás** | 15 commits de código (`0df9639`..`0561ccc`, `d966ceb`..`4299d09`, `2e6cc83`, `101955a` — PRs #172, #173, #174, #176) | ❌ **No coincide** *(Nicolás intervino como integrador sin asignación formal en front)* | ✅ **Pasa** (`frn-36-power-actions-reactive-status.test.tsx`, `fix-108-instance-permissions-route-param.test.tsx`) |
| **FIX-127** | Resiliencia de `SnapshotsTabContent.tsx` ante montaje aislado sin `ToastProvider`/`AuthProvider` y soporte simultáneo de parámetros de ruta `:vmid` e `:instanceRef` | **Cristian / Luz** | **Cristian** | `83a19b5`, `14db393`, `30a79a5`, `109f266` (PR #178 `2d088e6`) | ✅ **Coincide** *(Cristian)* | ✅ **Pasa** (`frn-30-snapshots-tab.test.tsx`) |
| **FIX-132** | Discriminación de errores HTTP `400` de validación de contraseña en `ForgotPasswordModal.tsx` para no descontar ni agotar los 3 intentos de código OTP | **Luz / Belinda** | **Cristian** | `3d06236` (PR #179 `2260b5b`) | ❌ **No coincide** *(Cristian)* | ✅ **Pasa** (`etapa3-pendientes.test.tsx` — bloque FIX-132) |
| **FRN-32** | Componente `AuditExportModal.tsx` y flujo de exportación de reportes de auditoría (CSV/JSON/PDF con filtros activos) | **Luz** | **Belinda** | `404c784` (PR #180 `b3169d0`) | ❌ **No coincide** *(Belinda)* | ✅ **Pasa** (`frn-32-audit-export.test.tsx` — 2/2 tests en verde) |
| **FIX-128 (Front)** | Permitir `VLAN = 0` o campo vacío como red sin etiquetar (untagged) en formularios de creación y edición de red de instancias | **Belinda** | **Belinda** | `b38ab17` (PR #181 `877eddd`) | ✅ **Coincide** *(Belinda)* | ✅ **Pasa** (`etapa3-pendientes.test.tsx` — bloque FIX-128) |
| **FIX-131** | Gráfico / Indicador de Tráfico de Red (`netInBps` y `netOutBps` con formateo automático `B/s`, `KB/s`, `MB/s`) en `InstanceResourceUsage.tsx` | **Belinda / Cristian** | *Sin commits en Parte C* | — | ⏳ **Pendiente** | ❌ **Falla** en `etapa3-pendientes.test.tsx` (`FIX-131: InstanceResourceUsage renderiza el tráfico de red`) |
| **FRN-33** | Centro de notificaciones (`AlertsDropdown`) y **banner rojo de saturación crítica** en `Dashboard.tsx` cuando RAM o disco superan el 90% | **Belinda** | *Sin commits en Parte C (base en Parte B)* | — | ⏳ **Incompleto** | ⚠️ **Falla 1 de 2 tests** en `frn-33-alerts-center.test.tsx` (falta el banner rojo de saturación crítica en `Dashboard.tsx`) |

---

## 5. Hallazgos Clave del Tramo Parte C (Frontend)

1. **Intervención decisiva de Nicolás (integrador de Backend) en el Frontend (58.6% del código del tramo):**
   - A pesar de ser desarrollador asignado a Backend y CI/CD, **Nicolás** realizó **23 commits (17 de código y 6 merges)** en `frontend` el 29 y 30 de octubre (PRs #169, #170, #172, #173, #174 y #176) para estabilizar la suite interna de Vitest, corregir problemas de accesibilidad (`FRN-34`), destrabar los estados reactivos de energía y WebSockets (`FRN-16`, `FRN-26`, `FRN-36`) y alinear las vistas de permisos (`FIX-108`), usuarios (`FIX-56`) y auditoría (`FRN-31`).

2. **Ausencia total de commits de Luz en la Parte C y redistribución de su carga:**
   - **Luz** registra **0 commits** en este tramo (su último commit en el repositorio fue `23977b1` el 20/10/2026 al inicio de la Parte B).
   - Todas las tareas críticas que figuraban asignadas a **Luz** en `documentacion/actual.md` fueron absorbidas por **Belinda** y **Cristian**:
     - `FRN-32` (Exportación de reportes de auditoría, asignada a Luz) fue implementada por **Belinda** en el PR #180 (`404c784`).
     - `FIX-127` (Resiliencia de `SnapshotsTabContent`, co-asignada a Cristian/Luz) fue resuelta por **Cristian** en el PR #178 (`83a19b5`..`109f266`).
     - `FIX-132` (Protección de intentos OTP en recuperación de contraseña, asignada a Luz/Belinda) fue implementada por **Cristian** en el PR #179 (`3d06236`).

3. **Intercambio de tareas entre Cristian y Belinda en la adaptación a PVE 8.4:**
   - `FRN-37` / `ADP-07` (Desacoplamiento de plantillas en el asistente de creación), que estaba asignada a **Belinda**, fue desarrollada por **Cristian** en 4 commits dentro del PR #171 (`53125e1`..`cdc4221`).
   - En contrapartida, **Belinda** tomó `ADP-08` (PR #177 `8d66b75`), `FRN-32` (PR #180 `404c784`) y `FIX-128` (PR #181 `b38ab17`).

4. **Salto de calidad en la suite externa (`test/front`) y únicas 2 brechas restantes:**
   - Gracias a los 41 commits de este tramo, la suite externa de auditoría (`test/front`) pasó de **257/263** a **261/263 tests en verde (99.24% de aprobación)**, con **112 de 114 archivos de test pasando al 100%**.
   - Al 03/11/2026 únicamente restan cerrar **2 puntos** en el submódulo `frontend`:
     1. **`FIX-131`:** Renderizar las métricas de red `netInBps` y `netOutBps` en `InstanceResourceUsage.tsx` (1 test fallido en `etapa3-pendientes.test.tsx`).
     2. **`FRN-33`:** Renderizar el banner rojo de saturación crítica en `Dashboard.tsx` cuando un recurso supera el 90% (1 test fallido en `frn-33-alerts-center.test.tsx`).
