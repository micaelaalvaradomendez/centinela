# Cruce de Tareas vs Commits — Consolidado Submódulos Backend y Frontend (Parte D: 03/11/2026 al 10/11/2026)

**Fecha de generación:** 10/11/2026  
**Submódulos analizados:** `backend/` y `frontend/`  
**Rango Backend (Parte D):** `42093c3` (2026-11-03, commit #274) → `d0f3f57` (2026-11-05, commit #281, HEAD actual)  
**Rango Frontend (Parte D):** `fa74412` (2026-11-03, commit #709) → `a268063` (2026-11-07, commit #785, HEAD actual, PR #198)  
**Cortes previos (Parte C):** Backend `fb8a41c` (commit #273, 02/11/2026) | Frontend `877eddd` (commit #708, PR #181, 02/11/2026)  
**Total de commits analizados en el tramo Parte D:** 85 commits (8 en Backend + 77 en Frontend)  
**Total acumulado en los repositorios:** 281 commits en `backend/` (`cd6a83e` → `d0f3f57`) y 785 commits en `frontend/` (`c04eddb` → `a268063`)

---

## 1. Metodología y Continuidad Histórica

Este documento continúa y completa la trazabilidad metodológica de las auditorías previas archivadas en `documentacion/analisis/`:
1. **Parte A (20/10/2026):**
   - Backend (`cruce-tareas-commits-backend.md`): Commits `#1 cd6a83e` al `#218 8daaf87` (218 commits).
   - Frontend (`cruce-tareas-commits-frontend.md`): Commits `#1 c04eddb` al `#536 91eabf4` (536 commits).
   - Foco: Requerimientos funcionales base (`BAC-01..45`, `FRN-01..29`) y primeros 82 fixes (`FIX-01..82`).
2. **Parte B (30/10/2026):**
   - Backend (`cruce-tareas-commits-backend-parte-b.md`): Commits `#219 43b011e` al `#262 908522e` (44 commits).
   - Frontend (`cruce-tareas-commits-frontend-parte-b.md`): Commits `#537 3b98e91` al `#667 594a4e0` (131 commits).
   - Foco: Ocultamiento de metadatos de hipervisor, autenticación robusta y fixes intermedios (`FIX-83..125`).
3. **Parte C (03/11/2026):**
   - Backend (`cruce-tareas-commits-backend-parte-c.md`): Commits `#263 61b973d` al `#273 fb8a41c` (11 commits).
   - Frontend (`cruce-tareas-commits-frontend-parte-c.md`): Commits `#668 443e97a` al `#708 877eddd` (41 commits).
   - Foco: Adaptación inicial al hipervisor real Proxmox VE 8.4 (`ADP-01..08`) y primeros fixes de hardening (`FIX-126..140`).
4. **Parte D (este informe — 10/11/2026):**
   - **Backend:** 8 commits posteriores (`#274 42093c3` al `#281 d0f3f57`), cubriendo el cierre de fixes de backend (`FIX-134`, `FIX-135`, `FIX-136`, `FIX-137`, `FIX-139`, `FIX-143`) y la política de instancias protegidas/borrado (`PR #30` / origen de `FIX-145`).
   - **Frontend:** 77 commits posteriores (`#709 fa74412` al `#785 a268063`, PRs #182 al #198), cubriendo la restauración de telemetría y gráficos de red (`FIX-131`, `FRN-22..28`), alineación del enum de auditoría (`FIX-144`), estabilización de tests unitarios (`FIX-141`), wizard Debian 12 (`ADP-07`), Centro de Notificaciones y banners reactivos (`FRN-33`), y la detección de la regresión de inventario (`FIX-146`).

---

## 2. Panorama General del Tramo Parte D (03/11/2026 al 10/11/2026)

### 2.1 Distribución en Submódulo Backend (`#274` al `#281`)

| Integrante | Usuario / Email en Git | Commits Totales | Merges | Commits de Código Real | % Código del Tramo |
|---|---|---:|---:|---:|---:|
| **Tayra** | `tayraag` / `Tayra` (`tayraag07@gmail.com`) | 5 | 0 | 5 | 71.4% |
| **Lisandro** | `lisandroureta` (`lisandroureta@gmail.com`) | 1 | 0 | 1 | 14.3% |
| **Nicolás** | `NicoLencina` (`nicolasangelo0@hotmail.com`) | 2 | 1 | 1 | 14.3% |
| **Lucas** | `Lucas Vela` (`lucas.tito1@hotmail.com`) | 0 | 0 | 0 | 0.0% |
| **Total Backend Tramo D** | — | **8** | **1** | **7** | **100.0%** |

### 2.2 Distribución en Submódulo Frontend (`#709` al `#785`)

| Integrante | Usuario / Email en Git | Commits Totales | Merges | Commits de Código Real | % Código del Tramo |
|---|---|---:|---:|---:|---:|
| **Cristian** | `Cristian` / `CristiannMaidana` (`cristian.mg.david@gmail.com`) | 69 | 12 | 57 | 91.9% |
| **Nicolás** | `NicoLencina` / `Lencina Nicolas Angelo` (`nicolasangelo0@hotmail.com`) | 3 | 1 | 2 | 3.2% |
| **Luz** | `luzpacello` / `Lucesita` (`luzpaceuni@gmail.com`) | 3 | 1 | 2 | 3.2% |
| **Belinda** | `Belinda` / `ibelinday` (`i_belinda_y@hotmail.com`) | 2 | 1 | 1 | 1.6% |
| **Total Frontend Tramo D** | — | **77** | **15** | **62** | **100.0%** |

---

## 3. Submódulo Backend: Detalle Cronológico y Cruce de Tareas

### 3.1 Lista Cronológica de Commits (#274 al #281)

| # | Fecha | Hash | Autor | Mensaje del Commit | Archivos Principales Afectados | Tarea Asociada |
|---:|---|---|---|---|---|---|
| 274 | 2026-11-03 | `42093c3` | Tayra | `fix: admitir instancias pausadas en controles de energía y borrado` | `instance_handler.go`, `client.go`, `proxmox_port.go`, `contrato-etapa1.md`, `swagger.json` | **FIX-134** |
| 275 | 2026-11-03 | `5f8bb5e` | Tayra | `fix: conciliar paradas externas de Proxmox para evitar falsas alertas` | `alert_evaluator_worker.go`, `client.go`, `audit_port.go`, `contrato-eventos.md` | **FIX-135** |
| 276 | 2026-11-03 | `08126ae` | Tayra | `fix: parametrizar umbrales y persistencia de alertas por entorno` | `alert_evaluator_worker.go`, `cmd/api/main.go`, `.env.example`, `SETUP.md` | **FIX-136** |
| 277 | 2026-11-04 | `4011b08` | Lisandro | `fix(auditoria): vocabulario cerrado de accion y resultado, enums en Swagger y normalización histórica` | `ports/audit_port.go`, `audit_handler.go`, `20261104_vocabulario_auditoria.sql`, `contrato-etapa1.md` (41 archivos) | **FIX-143** |
| 278 | 2026-11-04 | `5b3ac80` | Tayra | `fix: elimina endpoints fantasma e incluye auditorías globales` | `instance_handler.go`, `audit_repository.go`, `swagger.json`, `contrato_etapa1.go` | **FIX-137** |
| 279 | 2026-11-04 | `d92a568` | Tayra | `fix: reconcilia tareas antes de abrir la API` | `seguimiento_tareas.go`, `seguimiento_tareas_test.go`, `cmd/api/main.go` | **FIX-139** |
| 280 | 2026-11-05 | `df4cada` | Nicolás | `fix(instancias): no crear en VMID protegidos y exigir pertenencia para borrar` | `instance_handler.go`, `tarea_repository.go`, `client.go`, `contrato-etapa1.md` | **ADP-04 Ext. / Origen FIX-145** |
| 281 | 2026-11-05 | `d0f3f57` | Nicolás | `Merge pull request #30 from tayraag/fix/instancias-protegidas-y-borrado` | *(Merge PR #30)* | Merge **PR #30** |

### 3.2 Matriz de Cruce: Asignación Documental vs. Autoría Real en Backend

| Código Tarea | Descripción Técnica | Responsable Asignado (`actual.md` / `terminado3.md`) | Autor Real en Git | Commits Asociados | Coincidencia de Autoría | Estado de Verificación |
|---|---|---|---|---|---|---|
| **FIX-134** | Soporte para instancias en estado `paused` en controles de energía (`resume`, `stop`) y borrado | **Lucas / Nico / Tayra** | **Tayra** | `42093c3` (#274) | ⚠️ **Parcial** *(Asumida 100% por Tayra)* | ✅ **Pasa** (`TestFIX134_InstanciasPausadas` 4/4) |
| **FIX-135** | Reconciliación con tareas de Proxmox en Evaluador de Alertas para prevenir falsos positivos de caída ante apagados externos | **Lucas / Belinda / Tayra** | **Tayra** | `5f8bb5e` (#275) | ⚠️ **Parcial** *(Asumida 100% por Tayra)* | ✅ **Pasa** (`TestFIX135_AlertasParadasExternas` 3/3) |
| **FIX-136** | Parametrización de umbrales `ConfigAlertas` por variables de entorno y persistencia en Redis | **Nico / Lucas / Tayra** | **Tayra** | `08126ae` (#276) | ⚠️ **Parcial** *(Asumida 100% por Tayra)* | ✅ **Pasa** (`TestFIX136_UmbralesDeAlertaConfigurables` 4/4) |
| **FIX-143** | Vocabulario cerrado de auditoría (`EXITO`/`FALLA`/`PENDIENTE`), eliminación de constantes muertas, enums en Swagger y migración SQL | **Lucas / Nico** | **Lisandro** | `4011b08` (#277) | ❌ **No coincide** *(Absorbida por Lisandro)* | ✅ **Pasa** (`TestFIX143_VocabularioAuditoriaContratoCerrado` 4/4) |
| **FIX-137** | Saneamiento de código inerte: remoción de stub `/pause` y consulta de auditoría global de sistema (`organizacion_id IS NULL`) | **Lucas / Micaela** | **Tayra** | `5b3ac80` (#278) | ❌ **No coincide** *(Absorbida por Tayra)* | ✅ **Pasa** (`TestFIX137_CodigoInerteYAuditoriaDeSistema` 3/3) |
| **FIX-139** | Reconciliación de tareas asíncronas en `seguimiento_tareas` al arranque de la API | **Lucas / Nico** | **Tayra** | `d92a568` (#279) | ❌ **No coincide** *(Absorbida por Tayra)* | ✅ **Pasa** (`TestFIX139_ReconciliacionDeTareasAlArranque` 3/3) |
| **ADP-04 Ext. / PR #30** | Salteo de VMIDs protegidos y validación de pertenencia en borrado (`FueCreadaPorCentinela`) | **Tayra / Lisandro** *(en ADP-04)* | **Nicolás** | `df4cada` (#280), `d0f3f57` (#281) | ⚠️ **Iniciativa propia** *(Nicolás)* | ⚠️ **Origen de FIX-145**: La guarda 403 bloquea instancias reales preexistentes. |

---

## 4. Submódulo Frontend: Detalle Cronológico y Cruce de Tareas

### 4.1 Resumen de Pull Requests y Commits del Tramo Parte D (#709 al #785)

El tramo comprende 77 commits articulados en 15 Pull Requests (#182 al #198):

- **PR #182 (`fa74412` → `fa92405`, merge `26af24f` — 03/11/2026):**
  - **Autor:** Cristian (13 commits de código + merge).
  - **Alcance:** Restauración integral de la infraestructura de telemetría revertida en `4f743e4` (`useInstanceMetrics`, `getWsTicket`, `metricsWebSocket`), pestaña *"Métricas"*, modo *"En vivo"* y gráfico de Tráfico de Red (`netInBps`/`netOutBps`).
  - **Tareas:** **`FRN-22`, `FRN-23`, `FRN-26`, `FRN-27`, `FRN-28`, `FIX-131`, `FIX-142`**.
- **PR #183 (`74fce00` → `9668bfd`, merge `6e3758f` — 03/11/2026):**
  - **Autor:** Cristian (4 commits de código + merge).
  - **Alcance:** Pulido de estilos y alineación visual en vistas de Auditoría, menú contextual desplegable y botón de exportación.
  - **Tareas:** **`FRN-31`, `FRN-32`**.
- **PR #184 (`99c7a18`, merge `1d9d229` — 03/11/2026):**
  - **Autor:** Cristian (1 commit de código + merge).
  - **Alcance:** Reestructuración de la tabla de inventario en `Instances.tsx`.
  - **Impacto colateral:** *Eliminó silenciosamente las columnas `Checkbox` e `ID`*, lo que junto al commit `0bf4f33` constituyó la regresión documentada en **`FIX-146`**.
- **PR #185 (`d7a1958`, `71809f7`, merge `9787187` — 03/11/2026):**
  - **Autor:** Cristian (2 commits de código + merge).
  - **Alcance:** Incorporación de botón de retorno ("Volver") y navegación contextual en `InstanceDetailPage` e `InstanceDetailContent`.
- **PR #186 (`1b5285d` → `8d34b67`, merge `8b1d641` — 03/11/2026):**
  - **Autor:** Cristian (7 commits de código + merge).
  - **Alcance:** Formulario de Creación de Usuarios: sección de contraseña, selector combobox de roles, componente `RolesOptionCard`, corrección de tarjetas de instancias recientes y celdas de tabla de auditoría.
  - **Tareas:** **`FRN-04`, `FRN-05`**.
- **PR #187 (`4847927` → `6f4f4b2`, merge `e88078d` — 03/11/2026):**
  - **Autor:** Cristian (8 commits de código + merge).
  - **Alcance:** Refactorización visual de autenticación: eliminación de animaciones pesadas de fondo (`DotsHover`, `MainLayoutAuth`), normalización de botones en `LoginForm` y `RecoverPassword`, título del Dashboard y alineación de permisos.
- **PR #188 (`169daa7`, merge `b3569a1` — 03/11/2026):**
  - **Autor:** Cristian (1 commit de código + merge).
  - **Alcance:** Limpieza de opciones redundantes en barra lateral (`SideBar`) y enlace directo a configuración.
- **PR #189 (`408043e` → `f02f3e7`, merge `edd581c` — 03/11/2026):**
  - **Autor:** Cristian (8 commits de código + merge).
  - **Alcance:** Homogeneización de filtros y catálogo en pantalla de Auditoría (`useAuditCatalog`, `auditFilters`, `auditDataBoundary`, `auditLogService`) y estabilización de pruebas.
- **PR #190 (`eef5c38` → `8f5ec14`, merge `dccc522` — 03/11/2026):**
  - **Autor:** Cristian (3 commits de código + merge).
  - **Alcance:** Navegación cruzada bidireccional: navegación desde el detalle de un usuario (`detailsUserPage`) hacia la Auditoría prefiltrada por su identificador.
- **PR #191 (`e9fe857`, merge `1e5ec49` — 03/11/2026):**
  - **Autor:** Cristian (1 commit de código + merge).
  - **Alcance:** Estado activo y feedback de carga en botones de acción en `detailsUserPage`.
- **PR #192 (`a2a2761` → `efd59a5`, merge `c110791` — 04/11/2026):**
  - **Autor:** Cristian (7 commits de código + merge).
  - **Alcance:** Sincronización 1 a 1 de `auditActionLabels.ts` (31 códigos canónicos) con el enum cerrado del Backend (`EXITO`/`FALLA`/`PENDIENTE`), retiro de la lectura doble temporal y adaptación de `deletedUsersService.ts`.
  - **Tareas:** **`FIX-144` / `RF-08`**.
- **PR #193 (`bc44649`, `1b7e7d4`, merge `020b499` — 04/11/2026):**
  - **Autor:** Cristian (2 commits de código + merge).
  - **Alcance:** Actualización de mocks de `@/context/EventsContext` (`useOptionalEventsContext`) en suites internas `snapshots.test` e `instance-transient-states.test` para lograr el 100% (277/277 tests unitarios en verde).
  - **Tareas:** **`FIX-141` / `FIX-127`**.
- **PR #196 (`5370f31`, `0bf4f33`, merge `bed77ff` — 05/11/2026):**
  - **Autor:** Nicolás (2 commits de código + merge).
  - **Alcance:** `5370f31`: Ajuste en wizard de creación para solicitar IP estática por defecto y mapear el error `403 INSTANCE_NOT_OWNED`. `0bf4f33`: Adaptación de las aserciones de tests para ocultar la ausencia de la columna `ID` en la tabla de inventario.
  - **Tareas:** Soporte de red estática y sincronización con PR #30 de backend.
- **PR #197 (`087be9c`, `f1b77b6`, merge `da3f422` — 06/11/2026):**
  - **Autor:** Luz (2 commits de código + merge).
  - **Alcance:** `087be9c`: Adaptación del wizard de creación para clasificar y extraer metadatos de plantillas Debian 12 Bookworm (`amd64`). `f1b77b6`: Completar métricas de tráfico de red en detalle de instancia (`CU-05`).
  - **Tareas:** **`ADP-07` / `FRN-37`** y **`FIX-131` (Cierre)**.
- **PR #198 (`10beb6f`, merge `a268063` — 07/11/2026 al 10/11/2026):**
  - **Autor:** Belinda (1 commit de código + merge).
  - **Alcance:** Implementación del Centro de Notificaciones con campana reactiva (`NotificationBell.tsx`), contador de alertas no leídas y banners contextuales de alerta por saturación (`role="alert"`, `DashboardAlertBanner.tsx`) sincronizados por SSE en `AlertsContext.tsx`.
  - **Tareas:** **`FRN-33` (`RF-11`, `RNF-03`)**.

### 4.2 Matriz de Cruce: Asignación Documental vs. Autoría Real en Frontend

| Código Tarea | Descripción Técnica | Responsable Asignado (`actual.md` / `etapa3.md`) | Autor Real en Git | PR / Commits Asociados | Coincidencia de Autoría | Estado de Verificación |
|---|---|---|---|---|---|---|
| **FRN-22..28 / FIX-131 / FIX-142** | Restauración de WebSocket de métricas, pestaña "Métricas", modo "En vivo" y tarjeta de Tráfico de Red | **Belinda / Cristian / Luz** | **Cristian** | PR #182 (`fa74412`..`fa92405`, `26af24f`) | ⚠️ **Parcial** *(Cristian ejecutó la restauración)* | ✅ **Pasa** (`frn-22`, `frn-23`, `etapa3-pendientes`) |
| **FIX-144** | Alineación de filtros, contadores y diccionario de acciones con enum cerrado del Backend (`EXITO`/`FALLA`) | **Belinda / Cristian** | **Cristian** | PR #192 (`a2a2761`..`efd59a5`, `c110791`) | ✅ **Coincide** *(Cristian)* | ✅ **Pasa** (2/2 en `etapa3-pendientes.test.tsx`) |
| **FIX-141** | Mocks de `useOptionalEventsContext` y estabilización de vitest unitario (277/277 en verde) | **Equipo Frontend** | **Cristian** | PR #193 (`bc44649`, `1b7e7d4`, `020b499`) | ✅ **Coincide** *(Cristian)* | ✅ **Pasa** (43/43 suites internas) |
| **ADP-07 / FRN-37** | Soporte y extracción de metadata de plantillas Debian 12 Bookworm en wizard de creación | **Belinda / Luz** *(en `actual.md` Belinda)* | **Luz** | PR #197 (`087be9c`, `da3f422`) | ❌ **No coincide con titular** *(Ejecutado por Luz)* | ✅ **Pasa** (`frn-37-adp-07-wizard-debian12.test.tsx`) |
| **FIX-131 (Cierre)** | Afianzamiento de métricas de red en vista de telemetría de instancia (`CU-05`) | **Belinda / Cristian / Luz** | **Luz** | PR #197 (`f1b77b6`, `da3f422`) | ✅ **Coincide** *(Luz)* | ✅ **Pasa** (3/3 en suite externa) |
| **FRN-33** | Centro de Notificaciones (campana, contador) y banners contextuales de alerta crítica por saturación | **Belinda** | **Belinda** | PR #198 (`10beb6f`, `a268063`) | ✅ **Coincide 100%** *(Belinda)* | ✅ **Pasa** (`frn-33-alerts-center.test.tsx` 2/2) |
| **FIX-146 (Origen)** | Remoción de columna `ID` y checkbox en tabla de instancias (`Instances.tsx`) | *Defecto introducido* | **Cristian** (`99c7a18`) + **Nicolás** (`0bf4f33`) | PR #184 y PR #196 | ⚠️ **Regresión introducida** | ❌ **Pendiente de remediación en hardware real** |

---

## 5. Cuadro Consolidado Acumulado Histórico de los Repositorios

Alcanzado el hito de cierre de desarrollo de software al 10/11/2026, el historial completo de ambos submódulos refleja la siguiente distribución cuantitativa:

### 5.1 Submódulo Backend (`cd6a83e` → `d0f3f57`, 281 commits totales)

| Integrante | Commits Totales | % Total Repositorio | Commits de Código Real | Merges | % Código Real Efectivo | Rol Técnico Observado |
|---|---:|---:|---:|---:|---:|---|
| **Tayra** | **109** | 38.8% | **103** | 6 | 43.3% | Implementación central de servicios, Proxmox client y workers de alerta |
| **Lisandro** | **95** | 33.8% | **91** | 4 | 38.2% | Arquitectura base, simulador dual PVE 8/9, contratos y auditoría SQL |
| **Nicolás** | **54** | 19.2% | **33** | 21 | 13.9% | Gestión de PRs, integración y guardas de red/instancias |
| **Lucas** | **23** | 8.2% | **11** | 12 | 4.6% | Commits iniciales (inactivo en Backend desde el 14/10/2026) |
| **Total Backend** | **281** | **100.0%** | **238** | **43** | **100.0%** | — |

### 5.2 Submódulo Frontend (`c04eddb` → `a268063`, 785 commits totales)

| Integrante | Commits Totales | % Total Repositorio | Commits de Código Real | Merges | % Código Real Efectivo | Rol Técnico Observado |
|---|---:|---:|---:|---:|---:|---|
| **Cristian** | **509** | 64.8% | **414** | 95 | 72.5% | Motor principal de UI, vistas, refactorización, vitest y auditoría |
| **Nicolás** | **118** | 15.0% | **78** | 40 | 13.7% | Estabilización de suites, coordinación de PRs y adaptaciones de red |
| **Belinda** | **109** | 13.9% | **48** | 61 | 8.4% | Banners de saturación, exportación, modales de rollback y alertas |
| **Luz** | **49** | 6.2% | **31** | 18 | 5.4% | Wizard Debian 12, métricas de red y estructura fundacional inicial |
| **Total Frontend** | **785** | **100.0%** | **571** | **214** | **100.0%** | — |

### 5.3 Totales Globales Combinados del Proyecto (1.066 commits)

| Integrante | Commits Totales (Ambos Repos) | % del Proyecto | Commits de Código Real | Merges Totales | % Código Real Combinado |
|---|---:|---:|---:|---:|---:|
| **Cristian** | **509** | 47.7% | **414** | 95 | 51.2% |
| **Nicolás** | **172** | 16.1% | **111** | 61 | 13.7% |
| **Tayra** | **109** | 10.2% | **103** | 6 | 12.7% |
| **Belinda** | **109** | 10.2% | **48** | 61 | 5.9% |
| **Lisandro** | **95** | 8.9% | **91** | 4 | 11.2% |
| **Luz** | **49** | 4.6% | **31** | 18 | 3.8% |
| **Lucas** | **23** | 2.2% | **11** | 12 | 1.4% |
| **Consolidado General** | **1.066** | **100.0%** | **809** | **257** | **100.0%** |

---

## 6. Hallazgos Sociotécnicos Clave del Tramo Parte D

1. **Absorción de la Deuda Técnica de Lucas Vela por Tayra y Lisandro en Backend:**
   - En la planificación documental (`actual.md`, `futuro3.md`), **Lucas Vela** figuraba como asignado titular o cotitular de 6 tareas críticas de Etapa 3 (`FIX-133`, `FIX-134`, `FIX-135`, `FIX-136`, `FIX-137`, `FIX-139`, `FIX-143`).
   - La evidencia en Git demuestra que **Lucas no realizó ningún commit** en el repositorio `backend/` durante este tramo (su último commit data del 14/10/2026).
   - El 100% de la carga técnica fue absorbida en caliente por **Tayra** (resolviendo `FIX-134`, `FIX-135`, `FIX-136`, `FIX-137` y `FIX-139` entre el 03/11 y el 04/11) y por **Lisandro** (resolviendo `FIX-143` en el commit masivo `4011b08`).

2. **Hegemonía Operativa de Cristian en Frontend:**
   - En el submódulo `frontend/`, **Cristian ejecutó 57 de los 62 commits de código real (91.9%)** y fusionó 12 de los 15 PRs del tramo.
   - Destrabó la regresión de telemetría en PR #182 (13 commits), sincronizó el diccionario de auditoría en PR #192 y estabilizó los tests unitarios en PR #193.

3. **Reactivación y Cierre Especializado de Belinda y Luz:**
   - **Belinda** completó con éxito la implementación del Centro de Notificaciones y banners reactivos de saturación (`FRN-33`, PR #198), cerrando la última épica pendiente de UI de Etapa 3.
   - **Luz** intervino decisivamente en el PR #197, implementando la compatibilidad de nombres de plantillas Debian 12 Bookworm en el wizard (`ADP-07` / `FRN-37`), tarea que originalmente figuraba asignada a Belinda.

4. **Intervención de Nicolás y Génesis de los Defectos `FIX-145` y `FIX-146`:**
   - Nicolás intervino de forma cruzada en ambos submódulos (PR #30 en `backend/` y PR #196 en `frontend/`).
   - Su commit `df4cada` en Backend introdujo la verificación `FueCreadaPorCentinela` con respuesta `403 INSTANCE_NOT_OWNED` al borrar instancias. Aunque buscaba evitar la eliminación de máquinas del sistema, generó un bloqueo operativo sobre máquinas preexistentes en Proxmox físico, obligando a abrir el reporte **`FIX-145`**.
   - En Frontend, el commit `99c7a18` (Cristian) y el ajuste de tests `0bf4f33` (Nicolás) retiraron silenciosamente la columna `ID` unívoca de `Instances.tsx`, originando la tarea pendiente **`FIX-146`** para restaurar la identificación de recursos requerida en la Etapa 2.

5. **Estado Técnico Final del Software:**
   - Con la integración de los PRs analizados, **el 100% de las tareas de desarrollo de software (Backend y Frontend) se encuentra implementado y verificado en verde**.
   - Las únicas tareas pendientes corresponden a infraestructura y validación sobre el hardware Proxmox VE 8.4.21 físico (`FIX-123`, `INF-14`, `ADP-09`, `ADP-10`) y los dos ajustes de desacople identificados (`FIX-145` y `FIX-146`).
