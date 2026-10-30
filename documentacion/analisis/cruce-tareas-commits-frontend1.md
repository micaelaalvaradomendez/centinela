# Cruce de tareas, commits y autores reales — Submódulo `frontend`

**Proyecto:** El Centinela · **Repositorio analizado:** `frontend/` (`luzpacello/centinela-front`)  
**Fecha del análisis:** 20/10/2026  
**Universo:** 536 commits (`c04eddb` → `91eabf4`, 2026-08-17 a 2026-10-18)  
**Fuentes:** `git log` del submódulo + `documentacion/actual.md`, `terminado.md`, `terminado-1.md`, `terminado2.md`, `etapa1..3.md`, `futuro.md` y `test/informe.md`.

> **Objetivo:** Determinar con precisión quirúrgica quién implementó y commiteó realmente cada tarea del frontend, contrastando el **responsable asignado en la documentación/ClickUp** con el **autor real de los commits en el repositorio**. Permite esclarecer autorías ante tareas asignadas formalmente a integrantes que no realizaron el código efectivo.

---

## 1. Metodología y nivel de confianza

1. **Extracción íntegra del historial:** se analizaron los 536 commits del árbol del submódulo `frontend` (hash completo, hash corto, autor normalizado, fecha, mensaje, archivos modificados y topología de merge/pull request).
2. **Extracción y normalización de tareas:** se compilaron todas las tareas del frontend (`FRN-xx`, `FIX-xx`, `SEC-xx`, `INF-xx`, `LOGIN-xx`) y su campo **Asignado / Asignada** desde el espejo documental del tablero ClickUp.
3. **Trazabilidad de Pull Requests:** dado el flujo de trabajo en ramas personales (`luzpacello/cristian`, `luzpacello/rama-Beli`, `luzpacello/Nico`, `luzpacello/luz`, y ramas de fix específicas), se rastrearon los 141 PRs y los commits contenidos en cada rama origen.
4. **Criterios de evidencia probatoria:**

| Base | Significado | Confianza |
|---|---|---|
| **Explícita** | El commit aparece citado textualmente en `test/informe.md` o en las fichas de verificación de `documentacion/`. | Máxima certeza |
| **Mensaje** | El commit o PR nombra explícitamente el requerimiento, tarea o ID de ClickUp (`FRN-xx`, `FIX-xx`, `86e3...`, `RF-xx`). | Alta certeza |
| **Inferida** | El commit pertenece a una rama de feature/PR específico o modifica componentes unívocos de la tarea sin citar el ID en el mensaje. | Media-alta |
| **Merge** | Commit de integración de pull request o sincronización de rama (`main` / ramas remotas). No genera código nuevo pero visibiliza la autoría de integración. | Trazabilidad de integración |

> [!IMPORTANT]
> Los commits de **Merge** (158 commits) representan la integración de cambios en `main` o sincronizaciones entre ramas. En las estadísticas de autoría efectiva de código, se computan **378 commits de código real**. Para las tareas, se detalla tanto quién escribió el código fuente como quién aprobó/integró el PR.

> [!NOTE]
> **Normalización de identidades de autor:**
> - `Cristian` / `CristiannMaidana` (`cristian.mg.david@gmail.com`) → **Cristian**
> - `ibelinday` / `Belinda` (`i_belinda_y@hotmail.com`) → **Belinda**
> - `NicoLencina` / `Lencina Nicolas Angelo` (`nicolasangelo0@...`) → **Nico**
> - `Lucesita` / `luzpacello` (`luzpaceuni@gmail.com`) → **Luz**

---

## 2. Panorama general de autoría

### 2.1. Volumen de commits en el frontend

| Autor | Commits Totales | % Total | Commits de Código | Merges / PRs | Primer commit | Último commit |
|---|---:|---:|---:|---:|---|---|
| **Cristian** | 340 | 63.4% | 271 | 69 | 2026-08-18 | 2026-10-18 |
| **Belinda** | 84 | 15.7% | 36 | 48 | 2026-08-18 | 2026-10-16 |
| **Nico** | 77 | 14.4% | 50 | 27 | 2026-09-13 | 2026-10-18 |
| **Luz** | 35 | 6.5% | 21 | 14 | 2026-08-17 | 2026-10-18 |
| **Total** | **536** | **100%** | **378** | **158** | **2026-08-17** | **2026-10-18** |

### 2.2. Balance de asignación vs. ejecución real de tareas

| Integrante | Tareas Asignadas (Doc/ClickUp) | Tareas en las que commiteó código | Tareas Propias realizadas | Tareas Ajenas absorbidas / rescatadas |
|---|---:|---:|---:|---:|
| **Cristian** | 48 | 38 | 17 | 21 |
| **Belinda** | 40 | 21 | 12 | 9 |
| **Nico** | 4 | 33 | 4 | 29 |
| **Luz** | 37 | 23 | 7 | 16 |


> [!WARNING]
> **Hallazgo estructural clave:**
> - **Cristian** asumió **71.7% del código real** del frontend (271 de 378 commits), ejecutando **21 tareas completas asignadas a Belinda o a Luz** (incluyendo todo el Wizard de instancias `FRN-24`/`FRN-25` y el panel de usuarios `FRN-06`/`FRN-06B`).
> - **Belinda** figura asignada a 40 tareas, pero únicamente produjo **36 commits de código** (la mayoría maquetado visual inicial y textos), mientras que 48 de sus 84 commits fueron merges de PRs ajenos.
> - **Nico**, cuyo rol principal era Backend e Infraestructura, tuvo que intervenir de urgencia en el Frontend con **50 commits de código**, absorbiendo **29 tareas ajenas** para resolver el streaming SSE (`FRN-17A`/`17C`), los guards de seguridad rotos (`FIX-08`/`12`), el rediseño completo de login/2FA (`FIX-58`/`59`/`80`) y el Dashboard reactivo (`FRN-19A`/`19B`/`19C`).
> - **Luz** figura asignada a 37 tareas, pero sólo tiene **21 commits de código real**. La gran mayoría de sus tareas asignadas de inventario y wizard fueron realizadas íntegramente por Cristian.

---

## 3. Matriz de tareas: Asignado declarado vs. Autor(es) reales

A continuación se presenta el cruce de todas las tareas del frontend.  
Las tareas se categorizan según el grado de coincidencia entre quien figuraba como asignado y quién realizó efectivamente los commits:

- **Coincidencia total:** la totalidad de los autores que commitearon código coinciden con los asignados.
- **Coincidencia parcial:** el asignado participó, pero colaboraron o completaron la tarea otros integrantes.
- **No coincide / Hecha por otro:** el código fue escrito por alguien distinto al asignado en el tablero.
- **Sin commits en frontend:** tarea sin código commiteado (cancelada, pendiente, mockeada o puramente backend).

| Tarea | Título | Asignado oficial | Autores reales de código | Commits reales | Estado del cruce |
|---|---|---|---|:---:|---|
| `BAC-21` | Contrato base del canal de eventos/notificacione | Tayra y Lisandro | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-06` | Modal de alta y edición de usuarios sin funciona | Luz | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-07` | El frontend lee `body.code` en vez de `errorCode | Cristian | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-08` | Restaurar protección de rutas y guards de sesión | Belinda | **Belinda** (1), **Nico** (1) | 2 | ⚠️ Coincidencia parcial |
| `FIX-100` | Manejo seguro de valores `undefined` en `formatC | Luz / Cristian | — | 0 | ➖ Sin commits |
| `FIX-101` | Corrección de filtros reactivos en memoria (búsq | Luz | — | 0 | ➖ Sin commits |
| `FIX-102` | Corrección de atributo de accesibilidad `aria-la | Belinda / Luz | — | 0 | ➖ Sin commits |
| `FIX-103` | Manejo defensivo de arrays de telemetría de nodo | Belinda | — | 0 | ➖ Sin commits |
| `FIX-104` | Visibilidad permanente del banner de advertencia | Cristian | — | 0 | ➖ Sin commits |
| `FIX-104B` | Desacople de UPID y adopción del contrato seguro | Cristian / Luz | — | 0 | ➖ Sin commits |
| `FIX-105` | Desacople de UPID y adopción del contrato seguro | Cristian / Luz | — | 0 | ➖ Sin commits |
| `FIX-108` | Adaptación del Frontend para consumo de identifi | Belinda / Cristian | — | 0 | ➖ Sin commits |
| `FIX-111B` | Modal de confirmación con re-autenticación (Pass | Cristian | — | 0 | ➖ Sin commits |
| `FIX-115B` | Actualización de validación y formularios de con | Luz | — | 0 | ➖ Sin commits |
| `FIX-117B` | Adaptación del Dashboard de operador ante segreg | Belinda / Cristian | — | 0 | ➖ Sin commits |
| `FIX-12` | Falta guard de rol administrativo en rutas de ge | Belinda / Cristian | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-14` | Integración frontend del selector de instancias  | Cristian | — | 0 | ➖ Sin commits |
| `FIX-18` | Conexión de `RecoverPassword.tsx` a la API y ali | Belinda | — | 0 | ➖ Sin commits |
| `FIX-19` | Acciones de Restablecimiento Administrativo de C | Luz / Cristian | **Luz** (1) | 1 | ⚠️ Coincidencia parcial |
| `FIX-20` | Sincronización de políticas de complejidad y dis | Belinda | **Luz** (1) | 1 | ❌ No coincide |
| `FIX-21` | Opción de cerrar sesión en el cambio obligatorio | Belinda | **Cristian** (1) | 1 | ❌ No coincide |
| `FIX-22` | Guard administrativo en la ruta `/auditoria` (`F | Belinda / Luz | — | 0 | ➖ Sin commits |
| `FIX-24` | Confirmación del alta sin contraseña temporal y  | Luz | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-25` | Persistir la edición del perfil de usuario (`FRN | Luz / Cristian | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-27` | Mensaje específico cuando falla el envío del cor | Cristian | **Cristian** (1) | 1 | ✅ Coincidencia total |
| `FIX-28` | Regresión: el logout del frontend no envía el ac | Cristian | **Cristian** (1) | 1 | ✅ Coincidencia total |
| `FIX-29` | Validar la mayúscula en `validatePasswordComplex | Luz | — | 0 | ➖ Sin commits |
| `FIX-30` | Helper `canOperateInstance` en `usePermissions() | Cristian y Belinda | **Cristian** (5) | 5 | ⚠️ Coincidencia parcial |
| `FIX-36` | Regresión: restaurar el filtro "Acción" en Audit | Belinda y Luz (PR #77) | **Belinda** (1) | 1 | ⚠️ Coincidencia parcial |
| `FIX-38` | "Solo lectura" (`READ_ONLY`) usado como rol de u | Cristian y Belinda | **Cristian** (13) | 13 | ⚠️ Coincidencia parcial |
| `FIX-41` | Regresión: `requiredRole="Admin"` oculta "Roles  | Luz (autora de `1156eaa`) | **Cristian** (4), **Luz** (1) | 5 | ⚠️ Coincidencia parcial |
| `FIX-42` | useAuth().canAccessInstance` no revisa el rol (` | Luz (autora de `1156eaa` y `47e5c19`) | **Cristian** (3), **Luz** (2) | 5 | ⚠️ Coincidencia parcial |
| `FIX-43` | Corregir maquetado, visualización de IP, badges  | Luz / Cristian | — | 0 | ➖ Sin commits |
| `FIX-45` | El cliente de eventos abre el stream en una ruta | Nico | **Nico** (1) | 1 | ✅ Coincidencia total |
| `FIX-47` | Deduplicar eventos por `id` y limpiar restos del | Nico | **Nico** (1) | 1 | ✅ Coincidencia total |
| `FIX-56` | Distinguir suspensión y eliminación en contrato  | — | **Cristian** (15) | 15 | ℹ️ No asignada |
| `FIX-58` | Reemplazar placeholder textual del logo por el i | Belinda / Luz | **Nico** (6) | 6 | ❌ No coincide |
| `FIX-59` | Sustituir el placeholder "imagenes y informacion | Luz / Belinda | **Nico** (6) | 6 | ❌ No coincide |
| `FIX-60` | Reubicar el botón "Crear instancia" del Dashboar | Cristian / Luz | **Cristian** (2) | 2 | ⚠️ Coincidencia parcial |
| `FIX-61` | Cerrar sesión y detener reconexiones ante 401 al | Cristian | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-64` | Conexión SSE compartida y suscripción selectiva  | Cristian | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-74` | Reactividad de filtros y etiquetas accesibles en | Belinda / Luz | **Cristian** (8) | 8 | ❌ No coincide |
| `FIX-76` |  | — | **Luz** (1) | 1 | ℹ️ No asignada |
| `FIX-77` | Mapeo de errores HTTP de ciclo de vida (`BAC-29` | Cristian / Belinda | **Luz** (1) | 1 | ❌ No coincide |
| `FIX-78` | Desduplicación de indicadores de salud del nodo  | Belinda / Luz | **Belinda** (1) | 1 | ⚠️ Coincidencia parcial |
| `FIX-79` | Exposición directa de controles para selector re | Luz | **Belinda** (1) | 1 | ❌ No coincide |
| `FIX-80` | Alineación textual y normalización de copy en vi | Belinda / Nico | **Nico** (7) | 7 | ⚠️ Coincidencia parcial |
| `FIX-83` |  | — | **Nico** (2) | 2 | ℹ️ No asignada |
| `FIX-85` |  | — | **Luz** (1) | 1 | ℹ️ No asignada |
| `FIX-86` | Completar la máquina de estados y el mapeo D2 de | Cristian | **Belinda** (1) | 1 | ❌ No coincide |
| `FIX-87` | Quitar los datos hardcodeados de la ficha de usu | Cristian / Belinda | **Cristian** (4) | 4 | ⚠️ Coincidencia parcial |
| `FIX-89` | Ruta `/instances/new` con asistente estático sin | Luz / Belinda | **Cristian** (1) | 1 | ❌ No coincide |
| `FIX-90` | Quitar los placeholders muertos del Dashboard (` | Belinda | **Nico** (1) | 1 | ❌ No coincide |
| `FIX-91` | Eliminar el componente muerto `InstanceTester` ( | Cristian | — | 0 | ➖ Sin commits |
| `FIX-95` |  | — | **Luz** (1) | 1 | ℹ️ No asignada |
| `FIX-97` | Normalización de telemetría de inventario (escal | Cristian / Luz | **Cristian** (2) | 2 | ⚠️ Coincidencia parcial |
| `FIX-99` | Advertencia visible "Próximamente" en vista dire | Cristian | **Cristian** (1) | 1 | ✅ Coincidencia total |
| `FRN-01` | Maquetado del formulario de login | Belinda | **Cristian** (36), **Nico** (3), **Belinda** (2), **Luz** (1) | 42 | ⚠️ Coincidencia parcial |
| `FRN-02` | Validación en cliente del login | Belinda | — | 0 | ➖ Sin commits |
| `FRN-03` | Navbar y rutas base | Cristian y Luz | **Belinda** (11), **Cristian** (8), **Nico** (4), **Luz** (2) | 25 | ⚠️ Coincidencia parcial |
| `FRN-04` | Conexión del login con localhost | Cristian | **Cristian** (13), **Nico** (2), **Belinda** (2) | 17 | ⚠️ Coincidencia parcial |
| `FRN-05` | Panel de gestión de usuarios | Belinda | **Cristian** (30), **Nico** (2), **Belinda** (1), **Luz** (1) | 34 | ⚠️ Coincidencia parcial |
| `FRN-06` | Modal de creación y desactivación de usuarios | Luz | **Cristian** (21) | 21 | ❌ No coincide |
| `FRN-06B` | Modal de Edición de Usuario y Cambio de Rol: | — | **Cristian** (22), **Nico** (2), **Luz** (1) | 25 | ℹ️ No asignada |
| `FRN-07` | Selector de asignación de instancias | Cristian | **Cristian** (15) | 15 | ✅ Coincidencia total |
| `FRN-08` | Interceptor HTTP y manejo de `403 Forbidden | Cristian | **Cristian** (3), **Belinda** (2), **Nico** (1) | 6 | ⚠️ Coincidencia parcial |
| `FRN-09` | Vinculación 2FA mediante QR | Belinda y Luz | **Luz** (3), **Belinda** (2), **Cristian** (1), **Nico** (1) | 7 | ⚠️ Coincidencia parcial |
| `FRN-10` | Cambio obligatorio de contraseña temporal | Belinda | **Nico** (1), **Cristian** (1), **Luz** (1) | 3 | ❌ No coincide |
| `FRN-11` | Acciones administrativas de recuperación | Luz | **Luz** (1) | 1 | ✅ Coincidencia total |
| `FRN-12` | Vistas de recuperación de contraseña (RF-13) | Belinda | **Belinda** (1), **Luz** (1) | 2 | ⚠️ Coincidencia parcial |
| `FRN-13` | Flujo integral de logout y limpieza de sesión en | Cristian y Belinda | **Cristian** (10) | 10 | ⚠️ Coincidencia parcial |
| `FRN-14` | Suite de pruebas unitarias y de integración para | Belinda / Luz | **Belinda** (1) | 1 | ⚠️ Coincidencia parcial |
| `FRN-15` | Modales de confirmación antierror para acciones  | Belinda | **Luz** (5) | 5 | ❌ No coincide |
| `FRN-16` | Máquina de estados "Operación en progreso" por i | Cristian | **Luz** (1), **Belinda** (1) | 2 | ❌ No coincide |
| `FRN-16B` | Resincronización del estado "Operación en progre | Cristian | **Nico** (1) | 1 | ❌ No coincide |
| `FRN-17A` | Consumo de eventos en tiempo real y distribución | Cristian | **Nico** (6) | 6 | ❌ No coincide |
| `FRN-17B` | Desbloqueo reactivo de instancia y notificación  | Cristian | **Luz** (1) | 1 | ❌ No coincide |
| `FRN-17C` | Cliente de eventos con solicitud previa de ticke | Cristian | **Nico** (1) | 1 | ❌ No coincide |
| `FRN-18` | Integración en Frontend de niveles de acceso por | Cristian y Belinda | **Cristian** (36) | 36 | ⚠️ Coincidencia parcial |
| `FRN-19A` | Maquetado y medidores de recursos del Host (CPU  | Belinda (PR #78) | **Cristian** (4), **Nico** (2), **Belinda** (1) | 7 | ⚠️ Coincidencia parcial |
| `FRN-19B` | Semáforo de salud global e integración con `GET  | Belinda | **Belinda** (3), **Cristian** (2), **Nico** (2) | 7 | ⚠️ Coincidencia parcial |
| `FRN-19C` |  | — | **Nico** (5) | 5 | ℹ️ No asignada |
| `FRN-20A` | Tabla interactiva de inventario con badges de es | Luz (implementación inicial de Cristian en PR #74) | **Cristian** (29), **Belinda** (1) | 30 | ⚠️ Coincidencia parcial |
| `FRN-20B` | Filtros reactivos por tipo, estado y buscador di | Luz | **Belinda** (4), **Cristian** (3) | 7 | ❌ No coincide |
| `FRN-20C` | Columnas de CPU/RAM y estado en vivo en el Inven | Luz | **Cristian** (2) | 2 | ❌ No coincide |
| `FRN-21` |  | — | **Luz** (1) | 1 | ℹ️ No asignada |
| `FRN-22` | Hook `useInstanceMetrics` y cliente WebSocket pa | Belinda | — | 0 | ➖ Sin commits |
| `FRN-23` | Componentes de gráficos interactivos de rendimie | Luz | — | 0 | ➖ Sin commits |
| `FRN-24` | Wizard de creación: Paso 1 (Tipo de instancia) y | Belinda | **Cristian** (23) | 23 | ❌ No coincide |
| `FRN-25` | Wizard de creación: Paso 3 (Almacenamiento, Temp | Luz | **Cristian** (6) | 6 | ❌ No coincide |
| `FRN-26` | Envío de creación, manejo de errores de cuota (` | Cristian | **Cristian** (6) | 6 | ✅ Coincidencia total |
| `FRN-27` | Modal de reconfiguración de vCPU y RAM con adver | Cristian | — | 0 | ➖ Sin commits |
| `FRN-28` | Integración de edición en la página de Detalle d | Cristian | — | 0 | ➖ Sin commits |
| `FRN-29` | Pestaña "Snapshots" en Detalle de Instancia y fo | Luz | — | 0 | ➖ Sin commits |
| `FRN-30` | Modal destructivo de Rollback y eliminación con  | Cristian | — | 0 | ➖ Sin commits |
| `FRN-31` | Vista administrativa de Auditoría con filtros re | Belinda | **Cristian** (15), **Belinda** (4), **Nico** (3) | 22 | ⚠️ Coincidencia parcial |
| `FRN-32` | Componente y flujo de exportación de reportes de | Luz | — | 0 | ➖ Sin commits |
| `FRN-33` | Centro de Notificaciones y banners contextuales  | Belinda | — | 0 | ➖ Sin commits |
| `FRN-34` | Adaptación del cliente a cookies seguras de sesi | Cristian | — | 0 | ➖ Sin commits |
| `FRN-36` | Modularización de `InstanceDetailPage` para Pest | Cristian | — | 0 | ➖ Sin commits |
| `INF-04` | Red interna y reverse proxy con Nginx | Nico | **Cristian** (43), **Nico** (10), **Belinda** (4), **Luz** (2) | 59 | ⚠️ Coincidencia parcial |
| `LOGIN-02` | Pantalla y validación del código TOTP | Belinda y Luz | **Luz** (3), **Belinda** (2), **Cristian** (1), **Nico** (1) | 7 | ⚠️ Coincidencia parcial |
| `LOGIN-03` | Integración y prueba del prototipo | Cristian, Tayra y Lisandro | — | 0 | ➖ Sin commits |
| `LOGIN-04` | Prueba integral de autenticación y autorización | Cristian, Tayra y Lisandro | — | 0 | ➖ Sin commits |
| `SEC-02` | Cliente frontend compatible con refresh token Ht | Cristian | **Cristian** (2) | 2 | ✅ Coincidencia total |
| `SEC-03` | Contexto y sistema reactivo de permisos en Front | Cristian y Belinda | **Luz** (2) | 2 | ❌ No coincide |

---

## 4. Tareas "fantasmas" o ejecutadas por terceros

Este bloque detalla las discrepancias más críticas donde una tarea estaba documentada a nombre de una persona, pero fue programada parcial o totalmente por otra:

| Tarea | Título | Asignado en Tablero | Autor(es) Real(es) en Git | Commits Clave | Diagnóstico y Realidad Técnica |
|---|---|---|---|---|---|
| `FRN-06` | Modal de creación y desactivación de usuarios | Luz | Cristian (21), Nico (1) | `21610fb`, `58d15c0`, `e4a9798` | Asignada a Luz. Luz tiene 0 commits en esta tarea. Cristian maquetó e implementó toda la vista y modal; Nico conectó la API en el hotfix 86e3bh95w. |
| `FRN-06B` | Modal de edición de usuario y cambio de rol | Luz / Belinda | Cristian (22), Nico (2), Luz (1) | `8164670`, `f2d6d9f`, `e4a9798` | La edición y persistencia del perfil fue programada en un 90% por Cristian. Nico resolvió la integración con la API backend. |
| `FRN-10` | Cambio obligatorio de contraseña temporal | Belinda | Nico (1), Cristian (1), Luz (1) | `893bf09`, `b87ea1b`, `fe6d784` | Belinda tiene 0 commits. Nico desbloqueó el flujo (`893bf09`), Cristian añadió el logout obligatorio (`b87ea1b`) y Luz sincronizó validaciones (`fe6d784`). |
| `FRN-15` | Modales de confirmación antierror | Belinda | Luz (5) | `1d429f5`, `7e4d3af`, `c01dd36`, `fb58781`, `041b911` | Asignada formalmente a Belinda, pero implementada y testeada al 100% por Luz. Belinda no commiteó código en estos modales. |
| `FRN-16` | Máquina de estados "Operación en progreso" | Cristian | Luz (1), Belinda (1) | `614c912`, `ee15fb6` | Asignada a Cristian, pero implementada por Luz (`614c912` - máquina de estados y errorMapper) y complementada con traducciones por Belinda. |
| `FRN-17A` | Consumo de eventos en tiempo real (SSE) | Cristian | Nico (6) | `2a2db82`, `8739022`, `7f945ec`, `1ce64f6` | Asignada a Cristian, pero la arquitectura del cliente SSE, tickets efímeros y provider fue construida íntegramente por Nico. |
| `FRN-17B` | Desbloqueo reactivo de instancia y notificaciones | Cristian | Luz (1) | `2b120f4` | Asignada a Cristian, pero el hook reactivo y los toasts adaptativos ante eventos de instancia fueron programados por Luz. |
| `FRN-17C` | Cliente de eventos con conexión compartida por pestaña | Cristian | Nico (1) | `a3a86de` (86e3kzrqy) | Asignada a Cristian. Programada por Nico para evitar saturación de conexiones HTTP en navegadores. |
| `FRN-20B` | Filtros reactivos por tipo, estado y búsqueda | Luz | Belinda (4), Cristian (3) | `ef485be`, `7e40c4f`, `623bc28`, `2b1720c` | Asignada a Luz. Luz tiene 0 commits. Implementada entre Belinda (filtros IP y botones de cabecera) y Cristian (tabs y test alignment). |
| `FRN-20C` | Métricas de CPU/RAM en inventario | Luz | Cristian (2) | `338d658`, `e85a10f` | Asignada a Luz. Luz tiene 0 commits. Cristian implementó la normalización 0-100 y glifos em-dash. |
| `FRN-24` | Wizard de creación de instancias: Paso 1 y base | Belinda | Cristian (23) | `38f215e`, `fda22d7`, `944d0d2`, `c51dc9d` | Asignada a Belinda. Belinda tiene 0 commits. Cristian escribió la lógica, el hook `useCreateInstanceWizard` y todos los componentes. |
| `FRN-25` | Wizard de creación de instancias: Almacenamiento y Red | Luz | Cristian (6) | `3fa0277`, `c9fdb6b`, `21b8522`, `74cb5d1` | Asignada a Luz. Luz tiene 0 commits. Cristian desarrolló los pasos de recursos de disco, red y storage del wizard. |
| `SEC-03` | Contexto y sistema reactivo de permisos | Cristian y Belinda | Luz (2) | `47e5c19`, `1156eaa` | Luz fue quien tipó y refactorizó `AuthContext.ts` en TypeScript con roles y permisos. |
| `FIX-06..12` | Hotfixes de seguridad y guards tras rotura 6b08566 | Varios | Nico (4) | `231779f`, `64cad86`, `63dc6de`, `e4a9798` | El commit `6b08566` de Belinda movió rutas privadas a públicas. Nico ingresó a frontend para reponer los guards (`FIX-08`), el guard de rol (`FIX-12`) y mapear errores (`FIX-07`). |
| `FIX-58..80` | Rediseño y unificación visual de Login y 2FA | Belinda / Luz | Nico (12) | `e7ee40f`, `fe16eaa`, `902878e`, `b5da499` | Asignado a Belinda/Luz. Nico rehízo por completo el diseño visual, microinteracciones, accesibilidad y textos. |

---

## 5. Análisis individual por integrante

### 5.1. Cristian Maidana (`Cristian`)

- **Commits totales:** 340 (63.4% del repositorio) · **Commits de código real:** 271 (71.7% de todo el código frontend) · **Merges:** 69.

- **Rol en el proyecto:** Motor central del frontend. Absorbió la carga operativa mayoritaria del desarrollo.
- **Tareas que realizó efectivamente:**
  - **Autenticación e infraestructura UI:** Maquetado y componentes base (`FRN-01`, `FRN-03`, `FRN-04`), migración a `sessionStorage` (`SEC-02`), cliente API con fetch (`FRN-04`, `FRN-08`).
  - **Gestión de usuarios y accesos:** Implementación casi total de usuarios (`FRN-05`), creación (`FRN-06`), edición (`FRN-06B`), asignación de instancias (`FRN-07`), niveles de acceso `READ_ONLY`/`FULL_ACCESS` (`FRN-18`), logout integral (`FRN-13`).
  - **Inventario:** Tabla completa interactiva de instancias (`FRN-20A`), normalización de telemetría (`FRN-20C` / `FIX-97`), filtros y tabs (`FRN-20B`).
  - **Wizard de creación:** Creó en soledad el Wizard multipaso de instancias (`FRN-24`, `FRN-25`, `FRN-26`), con tipos, almacenamiento, red, SO y resumen.
  - **Auditoría:** Estructura completa de la pantalla de auditoría (`FRN-31`), filtros reactivos de fechas con calendar (`FIX-74`).
- **Discrepancia principal:** Se le asignaron 48 tareas, pero commiteó código en 38; de esas 38, **21 eran tareas formalmente asignadas a Belinda o a Luz**, las cuales Cristian tuvo que programar para que el proyecto no se detuviera.

### 5.2. Belinda (`Belinda` / `ibelinday`)

- **Commits totales:** 84 (15.7%) · **Commits de código real:** 36 (9.5%) · **Merges:** 48 (57.1% de su actividad fueron merges).

- **Rol en el proyecto:** Maquetadora inicial, diseño conceptual y revisora/integradora de PRs en GitHub.
- **Tareas que realizó efectivamente:**
  - Maquetado estructural inicial (`FRN-01`, `FRN-03`), layout general, tipografía de Penpot, ajustes de header y sidebar.
  - Maquetado básico de pantalla de auditoría (`FRN-31`), filtro de acción (`FIX-36`).
  - Pantalla inicial de recuperación de contraseña (`FRN-12`).
  - Filtros de inventario por IP y cabecera accesible (`FRN-20B` / `FIX-79`).
  - Integración inicial de `node/status` en Dashboard (`FRN-19B`).
- **Discrepancia principal:** Tenía asignadas **40 tareas críticas** (como Panel de Usuarios `FRN-05`, Modales de Confirmación `FRN-15`, Paso 1 del Wizard `FRN-24`, Cambio Obligatorio `FRN-10`). No implementó la lógica de ninguna de ellas: `FRN-05` y `FRN-24` fueron realizadas por Cristian; `FRN-15` fue realizada por Luz; `FRN-10` fue rescatada por Nico y Cristian.

### 5.3. Nico Lencina (`Nico`)

- **Commits totales:** 77 (14.4%) · **Commits de código real:** 50 (13.2%) · **Merges:** 27.

- **Rol en el proyecto:** Ingeniero de integración, soporte crítico y estabilizador del frontend desde el equipo backend/infra.
- **Tareas que realizó efectivamente:**
  - **Streaming SSE y eventos:** Construyó desde cero la conexión Server-Sent Events (`FRN-17A`), deduplicación por ID (`FIX-47`), reconexión selectiva por pestaña (`FIX-64`, `FRN-17C`) y resincronización de operaciones tras recarga (`FIX-61`, `FRN-16B`).
  - **Rescate de seguridad y contratos:** Cuando el commit `6b08566` expuso las rutas privadas, Nico implementó los guards de sesión y rol (`FIX-08`, `FIX-12`), alineó errores con el backend (`FIX-07`) y conectó las llamadas API de usuarios (`FIX-06`).
  - **Rediseño integral de Auth:** Reemplazó el diseño provisorio por la interfaz final pulida de Login, 2FA, ilustraciones, accesibilidad y textos (`FIX-58`, `FIX-59`, `FIX-80`).
  - **Dashboard en tiempo real:** Integración reactiva del histórico del nodo, eliminación de ruido en auditoría y curvas en vivo (`FRN-19A`, `FRN-19B`, `FRN-19C`, `FIX-83`, `FIX-90`).
- **Discrepancia principal:** Sólo tenía 4 tareas documentadas en frontend (`INF-04`, `FIX-45`, `FIX-47`, `FIX-80`), pero **escribió código en 33 tareas del frontend**, asumiendo tareas complejas que otros no lograron desbloquear.

### 5.4. Luz Pacello (`Luz` / `Lucesita`)

- **Commits totales:** 35 (6.5%) · **Commits de código real:** 21 (5.6%) · **Merges:** 14.

- **Rol en el proyecto:** Configuración inicial del repositorio, componentes específicos y modales de confirmación.
- **Tareas que realizó efectivamente:**
  - Inicialización del proyecto (`c04eddb`, `08a7fcf` - Vite + React con pnpm).
  - Creación de `TwoFactorPage` y `TwoFactorForm` (`LOGIN-02`, `FRN-09`).
  - Suite completa de modales de confirmación operativa antierror (`FRN-15` / `FIX-76` / `FIX-85` - 5 commits), con diálogo de parada, reinicio y eliminación.
  - Tipado de `AuthContext.ts` en TypeScript (`SEC-03`).
  - Mapeo de errores y máquina de estados (`FRN-16`, `FIX-77`).
  - Desbloqueo reactivo y toasts (`FRN-17B`).
  - Conexión de vista de detalle de instancia (`FRN-21` / `FIX-95`).
- **Discrepancia principal:** Tenía asignadas **37 tareas en documentación**, destacándose la creación de usuarios (`FRN-06`), la tabla interactiva de inventario (`FRN-20A`), los filtros (`FRN-20B`), las métricas (`FRN-20C`), el paso de almacenamiento del wizard (`FRN-25`), snapshots (`FRN-29`) y reportes (`FRN-32`). Prácticamente todas estas tareas asignadas fueron programadas por Cristian, o quedaron sin código en el frontend.

---

## 6. Listado exhaustivo commit por commit (1 a 536)

Se detalla a continuación el registro cronológico completo de los 536 commits del submódulo `frontend`:

| # | Fecha | Hash | Autor | Mensaje | Tarea(s) vinculada(s) | Base |
|---|---|---|---|---|---|---|
| 1 | 2026-08-17 | `c04eddb` | Luz | Initial commit | `FRN-03` | Inferida |
| 2 | 2026-08-18 | `08a7fcf` | Luz | 🌱 vite + react instalado con pnpm | `FRN-03` | Inferida |
| 3 | 2026-08-18 | `79f03d3` | Belinda | Merge pull request #1 from luzpacello/luz *(merge)* | `FRN-03` | Merge |
| 4 | 2026-08-18 | `98f21df` | Cristian | docs(README): pasos para instalar | — | Doc |
| 5 | 2026-08-19 | `ed1d3d8` | Luz | 🌱 se agrega las carpetas faltantes | `FRN-01` | Inferida |
| 6 | 2026-08-19 | `923649f` | Cristian | Merge pull request #2 from luzpacello/cristian *(merge)* | — | Merge |
| 7 | 2026-08-25 | `d009ac1` | Cristian | Merge pull request #3 from luzpacello/luz *(merge)* | `FRN-01` | Merge |
| 8 | 2026-08-26 | `2bad8ab` | Belinda | ❄️ Maqueta la estructura base y estilos del layout principal | `INF-04` | Inferida |
| 9 | 2026-08-26 | `9a2bf05` | Luz | Merge pull request #4 from luzpacello/rama-Beli *(merge)* | `INF-04` | Merge |
| 10 | 2026-09-01 | `851a18e` | Cristian | build(tsconfig.app.json): cambiar archivo para que typescript sepa que '@'  | `INF-04` | Inferida |
| 11 | 2026-09-01 | `452c492` | Cristian | build(vite.config.ts): cambiar archivo para que sepa que cuando hay '@' rep | `FRN-03` | Inferida |
| 12 | 2026-09-02 | `d0cb6e6` | Cristian | build(tsconfig.json): agregar ruta de path para shadcn sepa donde escribir | `FRN-01` | Inferida |
| 13 | 2026-09-02 | `f8cdc29` | Cristian | Merge pull request #5 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-03`, `INF-04` | Merge |
| 14 | 2026-09-06 | `79dacbd` | Cristian | build(package.json): instalar dependencias de shadcn/ui components | `INF-04` | Inferida |
| 15 | 2026-09-06 | `366004a` | Cristian | build(pnpm-lock.yaml): instalar dependencias de shadcn/ui component | `INF-04` | Inferida |
| 16 | 2026-09-06 | `4bb8c95` | Cristian | feat(index.css): agregar etilo necesarios luego de instalacion de dependenc | `INF-04` | Inferida |
| 17 | 2026-09-06 | `2117570` | Cristian | feat(components.json): se creo archivo luego de instalar dependencia shadcn | `INF-04` | Inferida |
| 18 | 2026-09-06 | `26aa8eb` | Cristian | feat(utils): agrega utilidad cn para componentes de shadcn | `FRN-01` | Inferida |
| 19 | 2026-09-06 | `2b81804` | Cristian | feat(button): creacion de button | `INF-04` | Inferida |
| 20 | 2026-09-06 | `427714f` | Cristian | feat(input): creacion de input | `INF-04` | Inferida |
| 21 | 2026-09-07 | `4d2c8fe` | Belinda | Merge pull request #6 from luzpacello/cristian *(merge)* | `FRN-01`, `INF-04` | Merge |
| 22 | 2026-09-09 | `0144f42` | Belinda | ❄️ Configura tipografía global de Penpot en Tailwind | `FRN-01` | Inferida |
| 23 | 2026-09-09 | `15feecd` | Luz | 🌱se agrega el componente stepper y Layout para auth y 2fa | `LOGIN-02`, `FRN-09` | Inferida |
| 24 | 2026-09-09 | `06af276` | Belinda | ❄️ Configura tipografía global de Penpot en Tailwind | `FRN-01` | Inferida |
| 25 | 2026-09-09 | `1a512b9` | Belinda | Merge pull request #7 from luzpacello/rama-Beli *(merge)* | `FRN-01` | Merge |
| 26 | 2026-09-09 | `4ba0083` | Luz | 🌱 Restauración de un button | `INF-04` | Inferida |
| 27 | 2026-09-09 | `30119b3` | Cristian | Merge pull request #8 from luzpacello/luz *(merge)* | `FRN-09`, `INF-04`, `LOGIN-02` | Merge |
| 28 | 2026-09-09 | `f2c3955` | Cristian | feat(checkbox): crear componente reutilizable checkbox con estilo | `FRN-01` | Inferida |
| 29 | 2026-09-09 | `bde4238` | Cristian | feat(Login): crear pagina login | `FRN-01`, `FRN-04` | Inferida |
| 30 | 2026-09-09 | `393012f` | Cristian | refactor(Login): cambiar pregunta | `FRN-01`, `FRN-04` | Inferida |
| 31 | 2026-09-09 | `55e3db6` | Cristian | feat(SignUp): crear pagina de registro | `FRN-01` | Inferida |
| 32 | 2026-09-10 | `0bfaa81` | Cristian | feat(separator): crear componente separator | `FRN-01` | Inferida |
| 33 | 2026-09-10 | `27c26cc` | Cristian | feat(label): crear componente label necesario para componente field | `FRN-01` | Inferida |
| 34 | 2026-09-10 | `db00e97` | Cristian | feat(field): crear componente field para estructurar los formularios | `FRN-01` | Inferida |
| 35 | 2026-09-10 | `d8ff411` | Belinda | Merge pull request #9 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-04` | Merge |
| 36 | 2026-09-10 | `ea92e52` | Luz | 🌱 feat(2fa): agrego TwoFactorPage y TwoFactoForm con validación interna | `LOGIN-02`, `FRN-09` | Inferida |
| 37 | 2026-09-10 | `09547ca` | Cristian | Merge pull request #10 from luzpacello/luz *(merge)* | `FRN-09`, `LOGIN-02` | Merge |
| 38 | 2026-09-10 | `34f631b` | Cristian | build(package): instalacion del paquete react-router | `INF-04` | Inferida |
| 39 | 2026-09-10 | `672fed1` | Cristian | build(pnpm-lock): instalacion del paquete react-router | `INF-04` | Inferida |
| 40 | 2026-09-11 | `4f2763e` | Belinda | ❄️ pantalla de recuperacion de contrasena | `FRN-12` | Inferida |
| 41 | 2026-09-11 | `da65fb6` | Belinda | Merge branch 'main' into rama-Beli *(merge)* | — | Merge |
| 42 | 2026-09-11 | `30cecfb` | Luz | Merge pull request #11 from luzpacello/rama-Beli *(merge)* | `FRN-12` | Merge |
| 43 | 2026-09-11 | `a9a8929` | Cristian | refactor(card): eliminar importacion de icons no utilizada | `INF-04` | Inferida |
| 44 | 2026-09-11 | `9601e47` | Cristian | refactor(inputOTPSeparate): simplificar importacion | `LOGIN-02`, `FRN-09` | Inferida |
| 45 | 2026-09-11 | `3074fe6` | Cristian | refactor(label): simplificar estilo | `FRN-01` | Inferida |
| 46 | 2026-09-11 | `b4785dc` | Cristian | feat: add authentication foundation | `FRN-01`, `FRN-04` | Inferida |
| 47 | 2026-09-11 | `64460c8` | Belinda | Merge pull request #13 from luzpacello/feature/auth-foundation *(merge)* | `FRN-01`, `FRN-04` | Merge |
| 48 | 2026-09-11 | `5cf2254` | Belinda | Merge pull request #12 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-09`, `INF-04`, `LOGIN-02` | Merge |
| 49 | 2026-09-11 | `368a4ce` | Cristian | feat: add login flow | `FRN-01`, `FRN-04` | Inferida |
| 50 | 2026-09-11 | `cd7794b` | Belinda | Merge pull request #14 from luzpacello/feature/login-flow *(merge)* | `FRN-01`, `FRN-04` | Merge |
| 51 | 2026-09-11 | `6d63cf7` | Cristian | feat: add organization registration flow | `FRN-01` | Inferida |
| 52 | 2026-09-11 | `f75a536` | Belinda | Merge pull request #15 from luzpacello/feature/organization-registration *( | `FRN-01` | Merge |
| 53 | 2026-09-11 | `35d22d9` | Cristian | feat: integrate application router | — | Inferida |
| 54 | 2026-09-11 | `1c4dab9` | Belinda | Merge pull request #16 from luzpacello/feature/data-router-integration *(me | — | Merge |
| 55 | 2026-09-13 | `33c3ed4` | Nico | necesario para comunicacion Nico | `INF-04` | Inferida |
| 56 | 2026-09-13 | `17c5a8a` | Cristian | Merge pull request #17 from luzpacello/feat/componente-comunicacion-nico *( | `INF-04` | Merge |
| 57 | 2026-09-15 | `a1e85eb` | Cristian | feat(MainLayoutAuth): cambio de texto para testear conexion con infra | `FRN-01`, `FRN-04` | Inferida |
| 58 | 2026-09-15 | `9001530` | Cristian | Merge pull request #18 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-04` | Merge |
| 59 | 2026-09-15 | `3bfbccc` | Nico | ci: deploy automatico del front a CT 103 (PRUEBAS) en cada push a main | `INF-04` | Inferida |
| 60 | 2026-09-15 | `86b9a02` | Luz | Merge pull request #19 from luzpacello/ci/deploy-automatico-front *(merge)* | `INF-04` | Merge |
| 61 | 2026-09-15 | `efb495a` | Nico | fix(auth): alinear la capa de autenticacion al contrato real del backend | `INF-04` | Inferida |
| 62 | 2026-09-15 | `6dcdeb5` | Belinda | Merge pull request #20 from luzpacello/fix/auth-contrato-real-backend *(mer | `INF-04` | Merge |
| 63 | 2026-09-15 | `b803b20` | Cristian | feat(MainLayoutAuth): cambio de texto (solo prueba) | `FRN-01`, `FRN-04` | Inferida |
| 64 | 2026-09-15 | `ee4ba82` | Cristian | feat(Dashboard): creacion de pagina, solo ubicacion de elementos primera ve | `FRN-19A`, `FRN-19B` | Inferida |
| 65 | 2026-09-15 | `ac5abdb` | Cristian | feat(Instances): creacion de pagina, solo ubicacion de elementos primera ve | `FRN-20A` | Inferida |
| 66 | 2026-09-15 | `eaea038` | Cristian | feat(NotFound): creacion de pagina, solo ubicacion de elementos primera ver | `INF-04` | Inferida |
| 67 | 2026-09-15 | `873d183` | Cristian | Merge branch 'main' into cristian *(merge)* | — | Merge |
| 68 | 2026-09-15 | `8965611` | Belinda | Merge pull request #21 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-04`, `FRN-19A`, `FRN-19B`, `FRN-20A`, `INF-04` | Merge |
| 69 | 2026-09-15 | `d8a36a5` | Nico | feat(layout_auth): mostrar fecha y hora del deploy en el panel | `INF-04` | Inferida |
| 70 | 2026-09-15 | `70e59a2` | Nico | ci(deploy): el deploy del front solo se dispara desde main | `INF-04` | Inferida |
| 71 | 2026-09-15 | `7a05a16` | Belinda | Merge pull request #22 from luzpacello/feat/stamp-deploy-panel *(merge)* | `INF-04` | Merge |
| 72 | 2026-09-15 | `566b5ab` | Belinda | Merge pull request #23 from luzpacello/ci/solo-main-despliega *(merge)* | `INF-04` | Merge |
| 73 | 2026-09-15 | `2f779e7` | Nico | infra(deploy-check): bloque temporal que muestra front y back juntos | `INF-04` | Inferida |
| 74 | 2026-09-15 | `26a927e` | Belinda | Merge pull request #24 from luzpacello/infra/deploy-check-front-back *(merg | `INF-04` | Merge |
| 75 | 2026-09-15 | `c259f20` | Nico | infra(deploy-check): el sello del front incluye el commit | `INF-04` | Inferida |
| 76 | 2026-09-15 | `dc605af` | Belinda | Merge pull request #25 from luzpacello/infra/front-stamp-con-commit *(merge | `INF-04` | Merge |
| 77 | 2026-09-16 | `45db606` | Cristian | feat(Sidebar): agregar nuevos links a paginas | `FRN-03` | Inferida |
| 78 | 2026-09-16 | `c674703` | Cristian | feat(applicationRoutes): crear nuevos links de pagina | `FRN-03` | Inferida |
| 79 | 2026-09-16 | `0dc7d8b` | Cristian | feat(Users): crear pagina de gestion de usuarios | `FRN-05`, `FRN-06B` | Inferida |
| 80 | 2026-09-16 | `7fff0a0` | Cristian | feat(CrearUsuarios): crear pagina de creacion de usuarios | `FRN-06` | Inferida |
| 81 | 2026-09-16 | `766643e` | Cristian | Merge pull request #26 from luzpacello/cristian *(merge)* | `FRN-03`, `FRN-05`, `FRN-06`, `FRN-06B` | Merge |
| 82 | 2026-09-16 | `d315882` | Nico | fix(routes): quitar import y ruta huerfanos a ComponentPlayground | `FRN-03` | Inferida |
| 83 | 2026-09-16 | `07df625` | Nico | Merge pull request #27 from luzpacello/fix/routes-import-huerfano-component | `FRN-03` | Merge |
| 84 | 2026-09-16 | `22d498a` | Cristian | feat(LoginForm): agregar link para navegacion a registro | `INF-04` | Inferida |
| 85 | 2026-09-16 | `25d13a1` | Cristian | feat(authenticationActions): crear funcion para la validacion de datos de c | `INF-04` | Inferida |
| 86 | 2026-09-16 | `15f5a22` | Cristian | feat(applicationRoutes): crear link para pagina de registro de organizacion | `INF-04` | Inferida |
| 87 | 2026-09-16 | `0d186da` | Cristian | Merge pull request #28 from luzpacello/cristian *(merge)* | `INF-04` | Merge |
| 88 | 2026-09-16 | `b15c588` | Cristian | refactor(applicationRoutes): eliminar importacion de pagina no existente | `INF-04` | Inferida |
| 89 | 2026-09-16 | `86efc92` | Cristian | merge: sincronizar cristian con main *(merge)* | — | Merge |
| 90 | 2026-09-16 | `ab434da` | Cristian | Merge pull request #29 from luzpacello/cristian *(merge)* | `INF-04` | Merge |
| 91 | 2026-09-16 | `feaaf27` | Belinda | acepta solo numeros | `LOGIN-02`, `FRN-09` | Inferida |
| 92 | 2026-09-16 | `0dabc28` | Belinda | Merge pull request #30 from luzpacello/rama-Beli *(merge)* | `FRN-09`, `LOGIN-02` | Merge |
| 93 | 2026-09-16 | `c784f04` | Belinda | qr de autenticacion | `LOGIN-02`, `FRN-09` | Inferida |
| 94 | 2026-09-16 | `78659d0` | Belinda | Merge pull request #31 from luzpacello/rama-Beli *(merge)* | `FRN-09`, `LOGIN-02` | Merge |
| 95 | 2026-09-17 | `51e7106` | Cristian | feat(nativeSelected): crear componente y agregarle colores segun diseño | `FRN-01` | Inferida |
| 96 | 2026-09-17 | `5d36e22` | Cristian | feat(tabs): crear componente y agregarle colores segun diseño | `FRN-01` | Inferida |
| 97 | 2026-09-17 | `6b1ccc4` | Cristian | feat(bardge): crear componente | `FRN-01` | Inferida |
| 98 | 2026-09-17 | `bc51bdb` | Cristian | feat(toast): crear componente y agregarle color segun diseño | `FRN-01` | Inferida |
| 99 | 2026-09-17 | `7775b62` | Cristian | Merge pull request #32 from luzpacello/cristian *(merge)* | `FRN-01` | Merge |
| 100 | 2026-09-18 | `1f62bfe` | Cristian | feat(tokenStorage): crear archivo sessionStorage para los tokens de acceso, | `FRN-01`, `FRN-04` | Inferida |
| 101 | 2026-09-18 | `52ea961` | Cristian | refactor(pendingLoginSession): ahora almacena el sessionStorage para poder  | `FRN-01`, `FRN-04` | Inferida |
| 102 | 2026-09-18 | `432d747` | Cristian | refactor(api): cambio localStorage por sessionStorage, conservando mismo no | `FRN-01`, `FRN-04` | Inferida |
| 103 | 2026-09-19 | `9e8b45c` | Cristian | Merge pull request #33 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-04` | Merge |
| 104 | 2026-09-19 | `cbf7481` | Belinda |  ❄️pantalla auditoria | `FRN-31` | Inferida |
| 105 | 2026-09-19 | `222fac5` | Belinda | Merge pull request #34 from luzpacello/rama-Beli *(merge)* | `FRN-31` | Merge |
| 106 | 2026-09-19 | `daa8bd5` | Cristian | feat(ApiResponseNotifier): crea el mensaje del toast, y asigna paginas para | `FRN-01`, `FRN-04` | Inferida |
| 107 | 2026-09-19 | `2711505` | Cristian | feat(apiClient): crea las conexiones fetch http | `FRN-04`, `FRN-08` | Inferida |
| 108 | 2026-09-19 | `84afd58` | Cristian | refactor(api): usa llamados de apiClient | `FRN-04`, `FRN-08` | Inferida |
| 109 | 2026-09-19 | `acc5ecf` | Cristian | feat(App): envuelve el router, para el manejo de notificaciones | `INF-04` | Inferida |
| 110 | 2026-09-19 | `12b0779` | Cristian | Merge pull request #35 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-04`, `FRN-08`, `INF-04` | Merge |
| 111 | 2026-09-20 | `6b08566` | Belinda | se agregan pantallas auditoria, crear usuario, user Detail, users | `FRN-03`, `FRN-05`, `FIX-08` | Explícita |
| 112 | 2026-09-20 | `c7a8800` | Belinda | Merge branch 'main' into rama-Beli *(merge)* | — | Merge |
| 113 | 2026-09-20 | `68bd23f` | Belinda | Merge pull request #36 from luzpacello/rama-Beli *(merge)* | `FIX-08`, `FRN-03`, `FRN-05` | Merge |
| 114 | 2026-09-21 | `231779f` | Nico | fix(86e3bhj61): restaurar guards de sesion en las rutas privadas | `FIX-08` | Mensaje |
| 115 | 2026-09-21 | `64cad86` | Nico | fix(86e3bh90p): panel de usuarios contra la API real y guard por rol | `FIX-12` | Mensaje |
| 116 | 2026-09-21 | `8a639a9` | Nico | Merge pull request #37 from luzpacello/Nico *(merge)* | `FIX-08`, `FIX-12` | Merge |
| 117 | 2026-09-21 | `63dc6de` | Nico | 86e3bh9d7 - mensajes de error alineados con los códigos del backend | `FIX-07` | Mensaje |
| 118 | 2026-09-21 | `a815bd6` | Nico | Merge pull request #38 from luzpacello/Nico *(merge)* | `FIX-07` | Merge |
| 119 | 2026-09-21 | `e4a9798` | Nico | 86e3bh95w - alta y edición de usuarios conectadas a la API | `FIX-06`, `FIX-24`, `FIX-25` | Mensaje |
| 120 | 2026-09-21 | `dbb4758` | Nico | Merge pull request #39 from luzpacello/Nico *(merge)* | `FIX-06`, `FIX-24`, `FIX-25` | Merge |
| 121 | 2026-09-21 | `6145740` | Luz | 🌱 refactor(2fa): Se cambia de LoginContinue a TwoFactorPage | `LOGIN-02`, `FRN-09` | Inferida |
| 122 | 2026-09-21 | `3e330fc` | Belinda | Merge pull request #40 from luzpacello/luz *(merge)* | `FRN-09`, `LOGIN-02` | Merge |
| 123 | 2026-09-21 | `9ccef0e` | Cristian | refactor(Users): eliminar panel lateral de detalle de usuarios | `FRN-05`, `FRN-06B` | Inferida |
| 124 | 2026-09-21 | `3112eea` | Belinda | se cambian extensiones | `FRN-03` | Inferida |
| 125 | 2026-09-21 | `e36684d` | Belinda | Merge pull request #41 from luzpacello/rama-Beli *(merge)* | `FRN-03` | Merge |
| 126 | 2026-09-21 | `c8557b6` | Nico | 86e3bpfw3 - Tipos del contrato de eventos compartidos con el backend | `BAC-21`, `FRN-17A` | Mensaje |
| 127 | 2026-09-21 | `0b50cd6` | Belinda | corrige el vite | `FRN-03` | Inferida |
| 128 | 2026-09-21 | `da396a1` | Belinda | corrige import sin uso | — | Inferida |
| 129 | 2026-09-21 | `8bcb893` | Belinda | Merge pull request #42 from luzpacello/rama-Beli *(merge)* | `FRN-03` | Merge |
| 130 | 2026-09-21 | `03209be` | Cristian | refactor(tabs): cambiar estilo de componente | `FRN-01` | Inferida |
| 131 | 2026-09-21 | `e00a1ab` | Cristian | feat(informationOfUser): crear componente | `FRN-05`, `FRN-06B` | Inferida |
| 132 | 2026-09-21 | `6f9cdfa` | Cristian | refactor(informationOfUser): cambiar tipo de opcion de rol | `FRN-05`, `FRN-06B` | Inferida |
| 133 | 2026-09-21 | `dace934` | Cristian | feat(index): crear 2 estilos de texto | — | Inferida |
| 134 | 2026-09-21 | `fdab248` | Cristian | feat(table): crear componente table | `FRN-01` | Inferida |
| 135 | 2026-09-21 | `1832e85` | Nico | Merge pull request #43 from luzpacello/Nico *(merge)* | `BAC-21`, `FRN-17A` | Merge |
| 136 | 2026-09-21 | `8b3a693` | Cristian | feat(rolesAndPermissions): crear componente | `FRN-07`, `FRN-18` | Inferida |
| 137 | 2026-09-21 | `8687564` | Cristian | refactor(badge): cambiar nombre de componente | `FRN-01` | Inferida |
| 138 | 2026-09-21 | `dc84dfd` | Cristian | feat(detailsUserPage): crear marco de pagina de detalle de usuarios | `FRN-05`, `FRN-06B` | Inferida |
| 139 | 2026-09-21 | `7a520c1` | Nico | fix(users): arreglo la carga del usuario real en UserDetail | `FRN-05`, `FRN-06B` | Inferida |
| 140 | 2026-09-21 | `5ad9b88` | Cristian | feat(Users): agregar funcionalidad a boton | `FRN-05`, `FRN-06B` | Inferida |
| 141 | 2026-09-21 | `e47e041` | Nico | fix(users): arreglo el guard de rol ADMIN en las rutas de usuarios | `FRN-05`, `FRN-06B` | Inferida |
| 142 | 2026-09-21 | `ca52e47` | Nico | fix(2fa): arreglo el mensaje de error de verificación TOTP | `LOGIN-02`, `FRN-09` | Inferida |
| 143 | 2026-09-21 | `e134383` | Nico | fix(auditoria): arreglo la vista de auditoría para consumir la API real | `FRN-31` | Inferida |
| 144 | 2026-09-21 | `27b4467` | Nico | fix(rutas): arreglo los imports con extensión .jsx | `FRN-03` | Inferida |
| 145 | 2026-09-21 | `d9e0339` | Nico | chore(limpieza): arreglo el código muerto del front | `FRN-01` | Inferida |
| 146 | 2026-09-21 | `efeb4a5` | Nico | fix(errores): arreglo el copy de los códigos de error que faltaban | `FRN-04`, `FRN-08` | Inferida |
| 147 | 2026-09-21 | `cc8a8a3` | Cristian | feat(switch): crear componente | `FRN-01` | Inferida |
| 148 | 2026-09-21 | `e5da01b` | Cristian | refactor(rolesAndPermissions): cambiar nombre de importacion | `FRN-07`, `FRN-18` | Inferida |
| 149 | 2026-09-21 | `9c1f768` | Cristian | feat(UserDetail): eliminar pagina | — | Inferida |
| 150 | 2026-09-21 | `d9cd6c1` | Nico | ci(front): agrego workflow de build en pull_request | `INF-04` | Inferida |
| 151 | 2026-09-21 | `a4d3f65` | Cristian | feat(applicationRoutes): crear ruta para nueva pagina | `FRN-03` | Inferida |
| 152 | 2026-09-22 | `924d65f` | Nico | fix(auditoria): arreglo las métricas mock de la vista de auditoría | `FRN-31` | Inferida |
| 153 | 2026-09-22 | `ea5f963` | Nico | Merge pull request #44 from luzpacello/Nico *(merge)* | `FRN-01`, `FRN-03`, `FRN-04`, `FRN-05`, `FRN-06B`, `FRN-08`, `FRN-09`, `FRN-31`, `INF-04`, `LOGIN-02` | Merge |
| 154 | 2026-09-22 | `7bc6dc5` | Nico | fix(menu): arreglo los enlaces muertos del menu lateral | `FRN-03` | Inferida |
| 155 | 2026-09-22 | `a7e7eb3` | Cristian | feat(user): crear persistencia de los datos de informacion del usuario | `INF-04` | Inferida |
| 156 | 2026-09-22 | `ff0bf64` | Cristian | feat(userDetailsServices): crear peticion de busqueda de datos del usuario | `INF-04` | Inferida |
| 157 | 2026-09-22 | `f89ab3c` | Cristian | feat(useUserDetails): crear hook para el manejo de la peticion de detalle d | `INF-04` | Inferida |
| 158 | 2026-09-22 | `24a97e7` | Cristian | feat(useEditableUser): crear hook para el manejo de datos editables del usu | `FRN-05`, `FRN-06B` | Inferida |
| 159 | 2026-09-22 | `4c3fcf3` | Cristian | feat(Users): cambiar boton de navegacion | `FRN-05`, `FRN-06B` | Inferida |
| 160 | 2026-09-22 | `66eb5be` | Nico | Merge pull request #45 from luzpacello/Nico *(merge)* | `FRN-03` | Merge |
| 161 | 2026-09-22 | `6f42344` | Cristian | feat(detailUserPage): cargar los datos reales del usuario | `FRN-05`, `FRN-06B` | Inferida |
| 162 | 2026-09-22 | `98c27c7` | Cristian | feat(informationOfUser): cambiar los valores estaticos de usuario por los r | `FRN-05`, `FRN-06B` | Inferida |
| 163 | 2026-09-22 | `06c379d` | Cristian | feat(rolesAndPermissions): cambiar valor estatico de rol del usuario por el | `FRN-07`, `FRN-18` | Inferida |
| 164 | 2026-09-22 | `893bf09` | Nico | fix(auth): arreglo el cambio obligatorio de contrasena para que el usuario  | `FRN-10` | Inferida |
| 165 | 2026-09-22 | `266fc1d` | Cristian | feat(userUserDetails): cambiar modo de alerta de error ahora con toast | — | Inferida |
| 166 | 2026-09-22 | `d6d2005` | Nico | Merge pull request #46 from luzpacello/Nico *(merge)* | `FRN-10` | Merge |
| 167 | 2026-09-22 | `541bf1c` | Cristian | feat(ApiResponseNotifier): crear alertas para diferentes errores | `FRN-04`, `FRN-08` | Inferida |
| 168 | 2026-09-22 | `c4ac90c` | Cristian | feat(userInstanceService): crear conexion a endpoint para obtener la lista  | `FRN-20A` | Inferida |
| 169 | 2026-09-22 | `cdca5f4` | Cristian | feat(useUserInstanceAccess): crear el hook para el manejo de la instancias  | `FRN-07`, `FRN-18` | Inferida |
| 170 | 2026-09-22 | `d6feb7b` | Cristian | feat(deatilUserPage): importar hook | — | Inferida |
| 171 | 2026-09-22 | `1e8e8f3` | Cristian | feat(rolesAndPermissions): cambiar datos estaticos de tabla por dinamicos d | `FRN-07`, `FRN-18` | Inferida |
| 172 | 2026-09-22 | `ab5e0f9` | Belinda | ❄️ se unifica el header, y se rehubica cerrar sesion | `FRN-03` | Inferida |
| 173 | 2026-09-22 | `d088df2` | Belinda | Merge pull request #47 from luzpacello/rama-Beli *(merge)* | `FRN-03` | Merge |
| 174 | 2026-09-22 | `635c8e7` | Belinda | ❄️  minimiza el alto del header | `FRN-03` | Inferida |
| 175 | 2026-09-22 | `d539253` | Belinda | Merge pull request #48 from luzpacello/main *(merge)* | `FRN-03` | Merge |
| 176 | 2026-09-22 | `94f2d7b` | Cristian | merge: sincronizar cristian con main *(merge)* | — | Merge |
| 177 | 2026-09-22 | `998da0c` | Cristian | feat(useSafeNavigate): eliminar hook de navegacion | `INF-04` | Inferida |
| 178 | 2026-09-22 | `65e1dc1` | Cristian | feat(nativeSelected): agregar componente | `FRN-01` | Inferida |
| 179 | 2026-09-22 | `995b2d6` | Cristian | feat(useSafeNavigate): agregar archivo eliminado | — | Inferida |
| 180 | 2026-09-22 | `80fa862` | Cristian | Merge pull request #49 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-03`, `FRN-04`, `FRN-05`, `FRN-06B`, `FRN-07`, `FRN-08`, `FRN-18`, `FRN-20A`, `INF-04` | Merge |
| 181 | 2026-09-22 | `bf82107` | Cristian | feat(userInstanceService): cambiar urls de las funciones acorde a enpoints  | `FRN-20A` | Inferida |
| 182 | 2026-09-22 | `791cb81` | Cristian | feat(userUserInstanceAccess): actualizar funcionalidad con endpoints actual | `FRN-07`, `FRN-18` | Inferida |
| 183 | 2026-09-22 | `f115fee` | Cristian | feat(rolesAndPermissions): actualizar funcionalidad para el funcionamiento  | `FRN-07`, `FRN-18` | Inferida |
| 184 | 2026-09-22 | `10c723d` | Cristian | feat(detailsUserPage): actualizar los datos de la instancia por endpoint ac | `FRN-05`, `FRN-06B` | Inferida |
| 185 | 2026-09-22 | `24695fa` | Cristian | merge: resolver sincronizacion pendiente *(merge)* | — | Merge |
| 186 | 2026-09-22 | `de6b25b` | Cristian | Merge pull request #50 from luzpacello/cristian *(merge)* | `FRN-05`, `FRN-06B`, `FRN-07`, `FRN-18`, `FRN-20A` | Merge |
| 187 | 2026-09-22 | `f07da88` | Belinda | Merge branch 'main' of https://github.com/luzpacello/centinela-front into r | — | Merge |
| 188 | 2026-09-22 | `76e7ebe` | Belinda | ❄️ se cambia posicion de cerrar sesion y se elimina el header | `INF-04` | Inferida |
| 189 | 2026-09-22 | `4e5fbab` | Belinda | ❄️Merge de actualizaciones remotas *(merge)* | — | Merge |
| 190 | 2026-09-22 | `61b2854` | Belinda | Merge pull request #51 from luzpacello/rama-Beli *(merge)* | `FRN-03`, `INF-04` | Merge |
| 191 | 2026-09-22 | `b50057e` | Belinda | ❄️  se olvido comentar | — | Inferida |
| 192 | 2026-09-22 | `44e57e6` | Belinda | Merge pull request #52 from luzpacello/main *(merge)* | `FRN-03`, `INF-04` | Merge |
| 193 | 2026-09-22 | `b1d5acc` | Belinda | ❄️  menu lateral con theme del sistema | `FRN-03` | Inferida |
| 194 | 2026-09-22 | `5b91e80` | Belinda | Merge pull request #53 from luzpacello/rama-Beli *(merge)* | `FRN-03`, `INF-04` | Merge |
| 195 | 2026-09-22 | `08d44ee` | Cristian | feat(authentication): eliminar propiedad refreshToken | `SEC-02` | Mensaje |
| 196 | 2026-09-22 | `ba9b5cf` | Cristian | feat(tokenStorage): eliminar constante de 'refreshKey' y de funcion 'storeA | `INF-04` | Inferida |
| 197 | 2026-09-22 | `4b7967f` | Cristian | feat(validateAuthenticationResponses): eliminar validacion de 'refreshToken | `INF-04` | Inferida |
| 198 | 2026-09-22 | `cb208f5` | Cristian | feat(apiClient): eliminar importacion y relacionado a 'refreshToken', actua | `INF-04` | Inferida |
| 199 | 2026-09-22 | `509a159` | Cristian | feat(api): eliminar importacion de metodo relacionado a 'refreshToken' | `INF-04` | Inferida |
| 200 | 2026-09-22 | `f0f7218` | Cristian | feat(authService): eliminar importacion y uso de 'refreshToken' | `INF-04` | Inferida |
| 201 | 2026-09-22 | `390b4b8` | Cristian | feat(api.d.ts): eliminar uso de 'refreshToken' por sessionStorage | `SEC-02` | Mensaje |
| 202 | 2026-09-22 | `781db6c` | Cristian | Merge pull request #54 from luzpacello/cristian *(merge)* | `INF-04`, `SEC-02` | Merge |
| 203 | 2026-09-22 | `c75acce` | Cristian | refactor(Users): cambiar estilo de boton principal | `FRN-05`, `FRN-06B` | Inferida |
| 204 | 2026-09-22 | `360d479` | Cristian | Merge pull request #55 from luzpacello/cristian *(merge)* | `FRN-05`, `FRN-06B` | Merge |
| 205 | 2026-09-23 | `c8e5592` | Belinda | ❄️  se prueba logo, sideBar estilo | `FRN-03` | Inferida |
| 206 | 2026-09-23 | `a2da502` | Belinda | ❄️  se arregla error | — | Inferida |
| 207 | 2026-09-23 | `3192cf4` | Belinda | Merge pull request #56 from luzpacello/rama-Beli *(merge)* | `FRN-03` | Merge |
| 208 | 2026-09-23 | `131ee7d` | Belinda | Revert "Merge pull request #56 from luzpacello/rama-Beli" | `FRN-04`, `FRN-08` | Inferida |
| 209 | 2026-09-23 | `a387e10` | Cristian | feat(Sidebar): eliminar link a pagina crear usuarios | `FRN-03` | Inferida |
| 210 | 2026-09-23 | `684b09c` | Cristian | feat(Users): agregar funcionalidad de navegacion a boton, navega a pagina C | `FRN-06` | Inferida |
| 211 | 2026-09-23 | `c92009e` | Cristian | feat(CrearUsuarios): agregar funcionalidad de navegacion en pagina | `FRN-06` | Inferida |
| 212 | 2026-09-23 | `bd40dc9` | Cristian | Merge pull request #57 from luzpacello/cristian *(merge)* | `FRN-03`, `FRN-06` | Merge |
| 213 | 2026-09-23 | `bac88cc` | Cristian | refactor(Users): cambiar estilos y funcionalidad | `FRN-05`, `FRN-06B` | Inferida |
| 214 | 2026-09-23 | `f2d6d9f` | Cristian | feat(userDetailService): agregar funcion para actualizar los datos del usua | `INF-04` | Inferida |
| 215 | 2026-09-23 | `60f81e5` | Cristian | feat(informationOfUser): agregar error a mail para validar que no haya 2 ig | `FRN-05`, `FRN-06B` | Inferida |
| 216 | 2026-09-23 | `5b3e073` | Cristian | feat(useUserInstanceAccess): agregar posible error | `FRN-07`, `FRN-18` | Inferida |
| 217 | 2026-09-23 | `0757bc8` | Cristian | feat(useEditableUser): agregar constante en funcion | `FRN-05`, `FRN-06B` | Inferida |
| 218 | 2026-09-23 | `8164670` | Cristian | feat(detailsUserPage): agregar funcionalidad para la activacion y actualiza | `FRN-05`, `FRN-06B` | Inferida |
| 219 | 2026-09-23 | `e8a2e27` | Cristian | Merge pull request #58 from luzpacello/cristian *(merge)* | `FRN-05`, `FRN-06B`, `FRN-07`, `FRN-18`, `INF-04` | Merge |
| 220 | 2026-09-23 | `a36a8c9` | Cristian | feat(modals): crear componente de modals | `FRN-01` | Inferida |
| 221 | 2026-09-23 | `6198cf9` | Cristian | feat(userDesactivationService): conexion para eliminar usuario | `FIX-56`, `FRN-06` | Mensaje |
| 222 | 2026-09-23 | `bd23180` | Cristian | feat(toast): agregar propiedad forceExpanded para si es mas de una notifica | `INF-04` | Inferida |
| 223 | 2026-09-23 | `fa6cbb7` | Cristian | feat(useCreateUser): crear funcion para crear usuario, agregar notificacion | `INF-04` | Inferida |
| 224 | 2026-09-23 | `1fb8bfd` | Cristian | feat(createUser): crear persistencia de datos de creacion de usuario y resp | `INF-04` | Inferida |
| 225 | 2026-09-23 | `8db47cf` | Cristian | refactor(useCreateUser): arreglar importacion de type de user | `INF-04` | Inferida |
| 226 | 2026-09-23 | `597ba06` | Cristian | feat(createUserService): crear conexion con enpoint y manejo de respuestas  | `FRN-05`, `FRN-06B` | Inferida |
| 227 | 2026-09-23 | `f290dad` | Cristian | feat(ConfirmUserAction): generador de modal para cada tipo de confirmacion  | `INF-04` | Inferida |
| 228 | 2026-09-24 | `b13931c` | Cristian | feat(deatilsUserPage): agregar funcionalidad para eliminar usuarios | `FIX-56`, `FRN-06` | Mensaje |
| 229 | 2026-09-24 | `58d15c0` | Cristian | feat(Users): agregar funcionalidad para eliminar usuarios | `FIX-56`, `FRN-06` | Mensaje |
| 230 | 2026-09-24 | `21610fb` | Cristian | feat(CrearUsuraios): agregar funcionalida para la creacion de usuarios | `FRN-06` | Inferida |
| 231 | 2026-09-24 | `b82b6a5` | Cristian | Merge pull request #59 from luzpacello/cristian *(merge)* | `FIX-56`, `FRN-01`, `FRN-05`, `FRN-06`, `FRN-06B`, `INF-04` | Merge |
| 232 | 2026-09-24 | `aa339d8` | Cristian | feat(userInstanceService): actualizar en base a endpoint | `FRN-20A` | Inferida |
| 233 | 2026-09-24 | `96c56c6` | Cristian | feat(useUserInstanceAccess): actualizar en base a endpoint | `FRN-07`, `FRN-18` | Inferida |
| 234 | 2026-09-24 | `b865771` | Cristian | feat(rolesAndPermissions): actualizar en base a endpoint | `FRN-07`, `FRN-18` | Inferida |
| 235 | 2026-09-24 | `7fbf969` | Cristian | Merge pull request #60 from luzpacello/cristian *(merge)* | `FRN-07`, `FRN-18`, `FRN-20A` | Merge |
| 236 | 2026-09-25 | `e02b821` | Cristian | feat(AuthContext): implementar flujo de logout y limpieza de sesion en clie | `FRN-13` | Mensaje |
| 237 | 2026-09-25 | `4556331` | Cristian | feat(App): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Mensaje |
| 238 | 2026-09-25 | `607e455` | Cristian | feat(useAuth): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Mensaje |
| 239 | 2026-09-25 | `deb59cb` | Cristian | feat(authService): implementar flujo de logout y limpieza de sesion en clie | `FRN-13` | Mensaje |
| 240 | 2026-09-25 | `b900d19` | Cristian | feat(apiClient): implementar flujo de logout y limpieza de sesion en client | `FRN-13` | Mensaje |
| 241 | 2026-09-25 | `72e205d` | Cristian | feat(api): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Mensaje |
| 242 | 2026-09-25 | `300a0bd` | Cristian | feat(api.d): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Mensaje |
| 243 | 2026-09-25 | `cce5a55` | Cristian | feat(tokenStorage): implementar flujo de logout y limpieza de sesion en cli | `FRN-13` | Mensaje |
| 244 | 2026-09-25 | `4062fdc` | Cristian | feat(ApiResponseNotifier): implementar flujo de logout y limpieza de sesion | `FRN-13` | Mensaje |
| 245 | 2026-09-25 | `b8469ea` | Cristian | Merge pull request #61 from luzpacello/cristian *(merge)* | `FRN-13` | Merge |
| 246 | 2026-09-25 | `b87ea1b` | Cristian | fix(ChangePassword): opción de cerrar sesión en el cambio obligatorio de co | `FIX-21`, `FRN-10` | Mensaje |
| 247 | 2026-09-25 | `b464c72` | Cristian | Merge pull request #62 from luzpacello/cristian *(merge)* | `FIX-21`, `FRN-10` | Merge |
| 248 | 2026-09-26 | `fe158ab` | Luz | 🌱 feat(users): conectar restablecimiento de contraseña y 2FA | `FRN-11`, `FIX-19` | Mensaje |
| 249 | 2026-09-26 | `6ad0c14` | Luz | Merge pull request #63 from luzpacello/luz *(merge)* | `FIX-19`, `FRN-11` | Merge |
| 250 | 2026-09-27 | `fe6d784` | Luz | 🌱 Sincronización de políticas de complejidad y distinción UX entre Cambio y | `FIX-20`, `FRN-10`, `FRN-12` | Mensaje |
| 251 | 2026-09-27 | `3da3e78` | Luz | Merge pull request #64 from luzpacello/luz *(merge)* | `FIX-20`, `FRN-10`, `FRN-12` | Merge |
| 252 | 2026-09-28 | `81c6424` | Belinda | ❄️; Se cambia Sidebar | `FRN-03` | Inferida |
| 253 | 2026-09-28 | `8d7ecab` | Belinda | Merge pull request #65 from luzpacello/rama-Beli *(merge)* | `FRN-03` | Merge |
| 254 | 2026-09-29 | `5c3d780` | Belinda | ❄️; fix ruta auditoria | `FRN-31` | Inferida |
| 255 | 2026-09-29 | `bd0c42f` | Belinda | ❄️ fix request | `FRN-04`, `FRN-08` | Inferida |
| 256 | 2026-09-29 | `749194e` | Belinda | Merge pull request #66 from luzpacello/rama-Beli *(merge)* | `FRN-04`, `FRN-08`, `FRN-31` | Merge |
| 257 | 2026-09-30 | `56b88f5` | Cristian | refactor(Auditoria): eliminar columnas que no seran utilizadas | `FRN-31` | Inferida |
| 258 | 2026-09-30 | `e6ad7ad` | Cristian | refactor(Auditoria): cambiar textos | `FRN-31` | Inferida |
| 259 | 2026-09-30 | `87d8926` | Cristian | Merge pull request #67 from luzpacello/cristian *(merge)* | `FRN-31` | Merge |
| 260 | 2026-09-30 | `3ea7fa4` | Cristian | fix(useCreatUser): arreglar notificacion por fallo de mail no enviado a la  | `FIX-27`, `FRN-06` | Mensaje |
| 261 | 2026-09-30 | `f6090c3` | Cristian | Merge pull request #68 from luzpacello/cristian *(merge)* | `FIX-27`, `FRN-06` | Merge |
| 262 | 2026-09-30 | `cfb88f7` | Cristian | fix(authService): el logout ahora envia Bearer y cookie HttpOnly para el co | `FIX-28`, `FRN-13` | Mensaje |
| 263 | 2026-09-30 | `bdc3116` | Cristian | Merge pull request #69 from luzpacello/cristian *(merge)* | `FIX-28`, `FRN-13` | Merge |
| 264 | 2026-09-30 | `54b397e` | Cristian | fix(usePermissions): helper canOperateInstance en usePermissions() | `FIX-30`, `FRN-18` | Mensaje |
| 265 | 2026-09-30 | `d6ac849` | Cristian | fix(authentication): helper canOperateInstance en usePermissions() | `FIX-30`, `FRN-18` | Mensaje |
| 266 | 2026-09-30 | `1a61fae` | Cristian | fix(validateAuthenticationResponses): helper canOperateInstance en usePermi | `FIX-30`, `FRN-18` | Mensaje |
| 267 | 2026-09-30 | `d1c2742` | Cristian | fix(api): helper canOperateInstance en usePermissions() | `FIX-30`, `FRN-18` | Mensaje |
| 268 | 2026-09-30 | `0b632f1` | Cristian | fix(api.d): helper canOperateInstance en usePermissions() | `FIX-30`, `FRN-18` | Mensaje |
| 269 | 2026-09-30 | `3d1e84a` | Cristian | Merge pull request #70 from luzpacello/cristian *(merge)* | `FIX-30`, `FRN-18` | Merge |
| 270 | 2026-10-01 | `1156eaa` | Luz | 🌱 Contexto y sistema reactivo de permisos en Frontend (UI/UX RBAC + Resourc | `SEC-03`, `FIX-41`, `FIX-42` | Explícita |
| 271 | 2026-10-01 | `92c7286` | Luz | Merge branch 'main' into luz *(merge)* | — | Merge |
| 272 | 2026-10-01 | `353f0f1` | Luz | 🌱 Solución al problema de merge | `INF-04` | Inferida |
| 273 | 2026-10-01 | `952b43c` | Luz | 🌱 fix(informationOfUser): modifico la validación de administrador con useAu | `FRN-05`, `FRN-06B` | Inferida |
| 274 | 2026-10-01 | `47e5c19` | Luz | 🌱 refactor(auth): tipar AuthContext con TypeScript | `SEC-03`, `FIX-42` | Explícita |
| 275 | 2026-10-01 | `4e0e7b7` | Luz | Merge pull request #71 from luzpacello/luz *(merge)* | `FIX-41`, `FIX-42`, `FRN-05`, `FRN-06B`, `INF-04`, `SEC-03` | Merge |
| 276 | 2026-10-02 | `3269dbd` | Cristian | fix(detailsUserPage): Regresión: requiredRole='Admin' oculta 'Roles y permi | `FIX-41`, `FRN-07`, `FRN-18` | Mensaje |
| 277 | 2026-10-02 | `f5e3c61` | Cristian | fix(authentication): Regresión: requiredRole='Admin' oculta 'Roles y permis | `FIX-41`, `FRN-07`, `FRN-18` | Mensaje |
| 278 | 2026-10-02 | `a351564` | Cristian | fix(AuthContext): Regresión: requiredRole='Admin' oculta 'Roles y permisos' | `FIX-41`, `FRN-07`, `FRN-18` | Mensaje |
| 279 | 2026-10-02 | `88f31b6` | Cristian | fix(PermissionGate): Regresión: requiredRole='Admin' oculta 'Roles y permis | `FIX-41`, `FRN-07`, `FRN-18` | Mensaje |
| 280 | 2026-10-02 | `b65e8b2` | Cristian | Merge pull request #72 from luzpacello/cristian *(merge)* | `FIX-41`, `FRN-07`, `FRN-18` | Merge |
| 281 | 2026-10-02 | `369fc05` | Cristian | fix(detailUserPage): 'Solo lectura' (READ_ONLY) usado erróneamente como rol | `FIX-38`, `FRN-18` | Mensaje |
| 282 | 2026-10-02 | `b2a5161` | Cristian | fix(CrearUsuarios): 'Solo lectura' (READ_ONLY) usado erróneamente como rol  | `FIX-38`, `FRN-18` | Mensaje |
| 283 | 2026-10-02 | `f65fa67` | Cristian | fix(informationOfUser): 'Solo lectura' (READ_ONLY) usado erróneamente como  | `FIX-38`, `FRN-18` | Mensaje |
| 284 | 2026-10-02 | `1da04a1` | Cristian | fix(rolesAndPermissions): 'Solo lectura' (READ_ONLY) usado erróneamente com | `FIX-38`, `FRN-18` | Mensaje |
| 285 | 2026-10-02 | `7d9e33a` | Cristian | fix(usePermissions): 'Solo lectura' (READ_ONLY) usado erróneamente como rol | `FIX-38`, `FRN-18` | Mensaje |
| 286 | 2026-10-02 | `fb60cb2` | Cristian | fix(AuthContext): 'Solo lectura' (READ_ONLY) usado erróneamente como rol de | `FIX-38`, `FRN-18` | Mensaje |
| 287 | 2026-10-02 | `8a111cb` | Cristian | fix(validateAuthenticationResponses): 'Solo lectura' (READ_ONLY) usado erró | `FIX-38`, `FRN-18` | Mensaje |
| 288 | 2026-10-02 | `6377790` | Cristian | fix(userDetailsService): 'Solo lectura' (READ_ONLY) usado erróneamente como | `FIX-38`, `FRN-18` | Mensaje |
| 289 | 2026-10-02 | `6fd2c7c` | Cristian | Merge pull request #73 from luzpacello/cristian *(merge)* | `FIX-38`, `FRN-18` | Merge |
| 290 | 2026-10-02 | `252f086` | Cristian | feat(Instances): Integración en Frontend de niveles de acceso por instancia | `FIX-38`, `FRN-18` | Mensaje |
| 291 | 2026-10-02 | `1a7d9b6` | Cristian | feat(useInstances): Integración en Frontend de niveles de acceso por instan | `FIX-38`, `FRN-18` | Mensaje |
| 292 | 2026-10-02 | `1a78041` | Cristian | feat(instanceService): Integración en Frontend de niveles de acceso por ins | `FIX-38`, `FRN-18` | Mensaje |
| 293 | 2026-10-02 | `c941f5d` | Cristian | feat(InstanceAction): Integración en Frontend de niveles de acceso por inst | `FIX-38`, `FRN-18` | Mensaje |
| 294 | 2026-10-02 | `0f1434a` | Cristian | feat(AuthContext): Integración en Frontend de niveles de acceso por instanc | `FIX-38`, `FRN-18` | Mensaje |
| 295 | 2026-10-02 | `52cd282` | Cristian | Merge pull request #74 from luzpacello/cristian *(merge)* | `FIX-38`, `FRN-18` | Merge |
| 296 | 2026-10-02 | `31f468e` | Cristian | fix(instanceAccess): useAuth().canAccessInstance no revisa el rol | `FIX-42`, `FRN-18` | Mensaje |
| 297 | 2026-10-02 | `d32b218` | Cristian | fix(AuthContext): useAuth().canAccessInstance no revisa el rol | `FIX-42`, `FRN-18` | Mensaje |
| 298 | 2026-10-02 | `ee60382` | Cristian | fix(usePermissions): useAuth().canAccessInstance no revisa el rol | `FIX-42`, `FRN-18` | Mensaje |
| 299 | 2026-10-02 | `5a274e9` | Cristian | Merge pull request #75 from luzpacello/cristian *(merge)* | `FIX-42`, `FRN-18` | Merge |
| 300 | 2026-10-02 | `2057fc0` | Belinda | ❄️  se agregar implmentan lo test y fix en audotoria | — | Inferida |
| 301 | 2026-10-02 | `cdbc8e9` | Belinda | Merge branch 'main' into rama-Beli *(merge)* | — | Merge |
| 302 | 2026-10-02 | `1cdcaf7` | Belinda | Merge pull request #76 from luzpacello/rama-Beli *(merge)* | — | Merge |
| 303 | 2026-10-02 | `fc27ba2` | Belinda | ❄ ️  fix restaurar filtro accion de auditoria | `FIX-36`, `FRN-31` | Mensaje |
| 304 | 2026-10-02 | `070e96b` | Belinda | Merge pull request #77 from luzpacello/rama-Beli *(merge)* | `FIX-36`, `FRN-31` | Merge |
| 305 | 2026-10-03 | `6572011` | Belinda | ❄ ️Maquetado y medidores de recursos del Host (CPU / RAM / Almacenamiento) | `FRN-19A` | Explícita |
| 306 | 2026-10-03 | `907efe5` | Belinda | Merge pull request #78 from luzpacello/rama-Beli *(merge)* | `FRN-19A` | Merge |
| 307 | 2026-10-03 | `2a2db82` | Nico | feat: cliente SSE con tickets efímeros y backoff | `FRN-17A` | Explícita |
| 308 | 2026-10-03 | `88c5cb6` | Nico | Merge pull request #79 from luzpacello/Nico *(merge)* | `FRN-17A` | Merge |
| 309 | 2026-10-03 | `8739022` | Nico | feat: events provider y consumo en tiempo real de instancias | `FRN-17A` | Explícita |
| 310 | 2026-10-03 | `15931cb` | Nico | Merge pull request #80 from luzpacello/Nico *(merge)* | `FRN-17A` | Merge |
| 311 | 2026-10-03 | `0407fbd` | Nico | fix: alinear badge de estado de eventos con los botones del header | `FRN-17A` | Mensaje |
| 312 | 2026-10-03 | `d47999d` | Nico | Merge pull request #81 from luzpacello/Nico *(merge)* | `FRN-17A` | Merge |
| 313 | 2026-10-05 | `7f945ec` | Nico | Corregir la ruta del canal de eventos en el frontend (86e3k5204) | `FIX-45`, `FRN-17A` | Explícita |
| 314 | 2026-10-05 | `1ce64f6` | Nico | Evitar avisos repetidos en el canal de eventos y limpiar restos del canal a | `FIX-47`, `FRN-17A` | Explícita |
| 315 | 2026-10-05 | `c185453` | Cristian | refactor(InstanceStatus/Card/List): eliminar archivos vacios | `FRN-20A` | Inferida |
| 316 | 2026-10-05 | `0901682` | Cristian | feat(Instances): Tabla interactiva de inventario con badges de estado e IP  | `FRN-20A` | Mensaje |
| 317 | 2026-10-05 | `fbd42ff` | Nico | Merge pull request #82 from luzpacello/Nico *(merge)* | `FIX-45`, `FIX-47`, `FRN-17A` | Merge |
| 318 | 2026-10-05 | `c8f218b` | Cristian | feat(instanceService): Tabla interactiva de inventario con badges de estado | `FRN-20A` | Mensaje |
| 319 | 2026-10-05 | `7472107` | Cristian | feat(useInstances): Tabla interactiva de inventario con badges de estado e  | `FRN-20A` | Mensaje |
| 320 | 2026-10-05 | `73b2b99` | Cristian | feat(InstanceAction): Tabla interactiva de inventario con badges de estado  | `FRN-20A` | Mensaje |
| 321 | 2026-10-05 | `d9fcdbb` | Cristian | feat(instance): Tabla interactiva de inventario con badges de estado e IP ( | `FRN-20A` | Mensaje |
| 322 | 2026-10-05 | `f69c468` | Cristian | feat(InstanceBadges): Tabla interactiva de inventario con badges de estado  | `FRN-20A` | Mensaje |
| 323 | 2026-10-05 | `98deff1` | Cristian | feat(InstanceIpAddress): Tabla interactiva de inventario con badges de esta | `FRN-20A` | Mensaje |
| 324 | 2026-10-05 | `189ab69` | Cristian | feat(useCopyInstanceIp): Tabla interactiva de inventario con badges de esta | `FRN-20A` | Mensaje |
| 325 | 2026-10-05 | `04987e0` | Cristian | feat(formatInstanceTelemetry): Tabla interactiva de inventario con badges d | `FRN-20A` | Mensaje |
| 326 | 2026-10-05 | `7234033` | Cristian | merge: sincronizar cristian con main *(merge)* | — | Merge |
| 327 | 2026-10-05 | `06458dd` | Cristian | Merge pull request #83 from luzpacello/cristian *(merge)* | `FRN-20A` | Merge |
| 328 | 2026-10-06 | `743e4cc` | Cristian | feat(InstanceAction): cambiar estilo de acciones en instancia | `INF-04` | Inferida |
| 329 | 2026-10-06 | `56d67bb` | Cristian | feat(Instances): cambiar estilo de tabla | `FRN-20A` | Inferida |
| 330 | 2026-10-06 | `5a86dce` | Cristian | Merge pull request #84 from luzpacello/cristian *(merge)* | `FRN-20A`, `INF-04` | Merge |
| 331 | 2026-10-06 | `bca9dd4` | Cristian | Fix(Instances): Corregir maquetado, visualización de IP, badges de estado y | `FRN-20A` | Inferida |
| 332 | 2026-10-06 | `825f9c2` | Cristian | Merge pull request #85 from luzpacello/cristian *(merge)* | `FRN-20A` | Merge |
| 333 | 2026-10-06 | `dfb7800` | Cristian | feat(dropDownMenu): crear nuevo componente | `FRN-03` | Inferida |
| 334 | 2026-10-06 | `8160c0f` | Cristian | fix(Dashboard): Reubicar botón 'Crear instancia' del Dashboard a la página  | `FIX-60`, `FRN-19A`, `FRN-20A` | Mensaje |
| 335 | 2026-10-06 | `612cb68` | Cristian | fix(Instances): Reubicar botón 'Crear instancia' del Dashboard a la página  | `FIX-60`, `FRN-19A`, `FRN-20A` | Mensaje |
| 336 | 2026-10-06 | `56a78df` | Cristian | Merge pull request #86 from luzpacello/cristian *(merge)* | `FIX-60`, `FRN-03`, `FRN-19A`, `FRN-20A` | Merge |
| 337 | 2026-10-06 | `50f7edb` | Cristian | fix(Dashboard): eliminar icon inutilizado | `FRN-19A`, `FRN-19B` | Inferida |
| 338 | 2026-10-06 | `013badc` | Cristian | Merge pull request #87 from luzpacello/cristian *(merge)* | `FRN-19A`, `FRN-19B` | Merge |
| 339 | 2026-10-06 | `49ca405` | Cristian | feat(nativeSelected): eliminar componente | `FRN-01` | Inferida |
| 340 | 2026-10-06 | `e6a344e` | Cristian | feat(combobox): crear componente y agregarle estilo segun diseño | `FRN-01` | Inferida |
| 341 | 2026-10-06 | `741b4a0` | Cristian | refactor(combobox): agregar padding 10px a cada opcion de celda | `INF-04` | Inferida |
| 342 | 2026-10-06 | `9407479` | Cristian | feat(input-group): necesario para el funcionamiento de componente combobox | `INF-04` | Inferida |
| 343 | 2026-10-06 | `cfdc659` | Cristian | refactor(CrearUsuario): cambiar componente 'nativeSelected' por 'combobox' | `FRN-06` | Inferida |
| 344 | 2026-10-06 | `62c5c4b` | Cristian | refactor(rolesAndPermissions): cambiar componente 'nativeSelected' por 'com | `FRN-07`, `FRN-18` | Inferida |
| 345 | 2026-10-06 | `3009fe5` | Cristian | refactor(informationOfUser): cambiar componente 'nativeSelected' por 'combo | `FRN-05`, `FRN-06B` | Inferida |
| 346 | 2026-10-06 | `2aa0787` | Cristian | feat(textarea): componenete necesario para funcionamiento de 'combobox' | `INF-04` | Inferida |
| 347 | 2026-10-06 | `5be6aa9` | Cristian | Merge pull request #88 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-05`, `FRN-06`, `FRN-06B`, `FRN-07`, `FRN-18`, `INF-04` | Merge |
| 348 | 2026-10-06 | `493fe8a` | Cristian | refactor(combobox): cambiar color de input cuando esta seleccionado | `INF-04` | Inferida |
| 349 | 2026-10-06 | `2111a01` | Cristian | feat(pagination): crear componente pagination con estilo segun diseño | `FRN-01` | Inferida |
| 350 | 2026-10-06 | `d98db16` | Cristian | refactor(Instances): cambiar pagination por componente | `FRN-20A` | Inferida |
| 351 | 2026-10-06 | `5a1001c` | Cristian | refactor(Auditoria): cambiar pagination estatico por componente | `FRN-31` | Inferida |
| 352 | 2026-10-06 | `2daf259` | Cristian | refactor(Users): cambiar pagination por componente | `FRN-05`, `FRN-06B` | Inferida |
| 353 | 2026-10-06 | `5258908` | Cristian | Merge pull request #89 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-05`, `FRN-06B`, `FRN-20A`, `FRN-31`, `INF-04` | Merge |
| 354 | 2026-10-06 | `5f391f5` | Cristian | feat(Instances): actualizar actualizacion de badges | `FRN-20A` | Inferida |
| 355 | 2026-10-06 | `dfbee12` | Cristian | feat(InstanceAction): actualizar actualizacion de badges segun accion de ej | `INF-04` | Inferida |
| 356 | 2026-10-06 | `297aec0` | Cristian | feat(InstanceBadges): actualizar actualizacion de badges segun accion de ej | `INF-04` | Inferida |
| 357 | 2026-10-06 | `f7338bc` | Cristian | feat(useInstances): actualizar actualizacion de badges segun accion de ejec | `FRN-20A` | Inferida |
| 358 | 2026-10-06 | `6a859d7` | Cristian | Merge pull request #90 from luzpacello/cristian *(merge)* | `FRN-20A`, `INF-04` | Merge |
| 359 | 2026-10-06 | `878bc39` | Cristian | feat(InstanceBadges): mostrar el cambio de estado una vez terminada la soli | `INF-04` | Inferida |
| 360 | 2026-10-06 | `08fde7c` | Cristian | feat(useInstances): mostrar el cambio de estado una vez terminada la solici | `FRN-20A` | Inferida |
| 361 | 2026-10-06 | `fc6f9ed` | Cristian | Merge pull request #91 from luzpacello/cristian *(merge)* | `FRN-20A`, `INF-04` | Merge |
| 362 | 2026-10-06 | `f9b2c0a` | Nico | fix(instancias): resincronizar estado en progreso tras F5 y reconexión (86e | `FIX-61`, `FRN-16B` | Explícita |
| 363 | 2026-10-06 | `a3a86de` | Nico | Mantener una única conexión de eventos por pestaña con suscripción selectiv | `FIX-64`, `FRN-17C` | Explícita |
| 364 | 2026-10-06 | `e7ee40f` | Nico | Reemplazar el texto provisorio de bienvenida por la información de Centinel | `FIX-58`, `FIX-59` | Explícita |
| 365 | 2026-10-06 | `ea64514` | Nico | Merge pull request #92 from luzpacello/Nico *(merge)* | `FIX-58`, `FIX-59`, `FIX-61`, `FIX-64`, `FRN-16B`, `FRN-17C` | Merge |
| 366 | 2026-10-06 | `3875926` | Nico | Quitar el bloque temporal de verificación de deploy del login y declarar el | `INF-04` | Inferida |
| 367 | 2026-10-06 | `7367405` | Nico | Merge pull request #93 from luzpacello/Nico *(merge)* | `INF-04` | Merge |
| 368 | 2026-10-06 | `f9497e4` | Cristian | feat(createNewInstances): crear nueva pagina | `FRN-24` | Inferida |
| 369 | 2026-10-06 | `317e247` | Cristian | feat(stepper): agregar descripcion opcion | `INF-04` | Inferida |
| 370 | 2026-10-06 | `79865c9` | Cristian | feat(applicationRoutes): crear ruta para pagina crear nueva instancia | `INF-04` | Inferida |
| 371 | 2026-10-06 | `24ac662` | Cristian | feat(Instances): agregar funcionalidad a boton principal | `FRN-20A` | Inferida |
| 372 | 2026-10-06 | `67b85c3` | Cristian | merge: sincronizar cristian con main *(merge)* | — | Merge |
| 373 | 2026-10-06 | `562bdd6` | Cristian | Merge pull request #94 from luzpacello/cristian *(merge)* | `FRN-20A`, `FRN-24`, `INF-04` | Merge |
| 374 | 2026-10-06 | `31e7cfd` | Cristian | fix(stepper): eliminar tone no utilizable | `FRN-01` | Inferida |
| 375 | 2026-10-06 | `448e415` | Cristian | Merge pull request #95 from luzpacello/cristian *(merge)* | `FRN-01` | Merge |
| 376 | 2026-10-07 | `24315e6` | Cristian | refactor(Users): cambiar cards por tabs con funcionalidad y agregar opcione | `FRN-05`, `FRN-06B` | Inferida |
| 377 | 2026-10-07 | `83d3e04` | Cristian | refactor(Auditoria): cambiar cards por tabs con funcionalidad y agregar opc | `FRN-31` | Inferida |
| 378 | 2026-10-07 | `9751d32` | Cristian | Merge pull request #96 from luzpacello/cristian *(merge)* | `FRN-05`, `FRN-06B`, `FRN-31` | Merge |
| 379 | 2026-10-07 | `fe16eaa` | Nico | feat(auth): nuevo diseño del login y unificación visual del flujo de autent | `FIX-80`, `FRN-01` | Explícita |
| 380 | 2026-10-07 | `902878e` | Nico | refactor(auth): restaurar la clase text-info en los enlaces del flujo de au | `FIX-58`, `FIX-59`, `FIX-80` | Mensaje |
| 381 | 2026-10-07 | `f4e89e1` | Nico | Merge pull request #97 from luzpacello/Nico *(merge)* | `FIX-58`, `FIX-59`, `FIX-80`, `FRN-01` | Merge |
| 382 | 2026-10-07 | `e45645c` | Nico | fix(auth): tarjetas alcanzables en pantallas bajas y puntos un poco más mar | `FIX-58`, `FIX-59`, `FIX-80` | Mensaje |
| 383 | 2026-10-07 | `14d7df9` | Nico | Merge pull request #98 from luzpacello/Nico *(merge)* | `FIX-58`, `FIX-59`, `FIX-80` | Merge |
| 384 | 2026-10-07 | `87885e1` | Nico | fix(auth): puntos del fondo más marcados (86e3kz4du) | `FIX-58`, `FIX-59`, `FIX-80` | Mensaje |
| 385 | 2026-10-07 | `2dc72ed` | Nico | Merge pull request #99 from luzpacello/Nico *(merge)* | `FIX-58`, `FIX-59`, `FIX-80` | Merge |
| 386 | 2026-10-07 | `82d89a6` | Nico | fix(auth): subir el contraste de los puntos del fondo (86e3kz4du) | `FIX-58`, `FIX-59`, `FIX-80` | Mensaje |
| 387 | 2026-10-07 | `e91b70f` | Nico | feat(auth): contorno con sombra en los iconos de beneficios del login (86e3 | `FIX-58`, `FIX-59`, `FIX-80` | Mensaje |
| 388 | 2026-10-07 | `57fbf75` | Nico | Merge pull request #100 from luzpacello/Nico *(merge)* | `FIX-58`, `FIX-59`, `FIX-80` | Merge |
| 389 | 2026-10-07 | `12d7bfa` | Nico | feat(auth): intro de video de bienvenida antes del login | `FRN-01`, `FRN-04` | Inferida |
| 390 | 2026-10-07 | `302eff7` | Nico | Merge pull request #101 from luzpacello/Nico *(merge)* | `FRN-01`, `FRN-04` | Merge |
| 391 | 2026-10-07 | `226539e` | Nico | feat(auth): puntos del fondo con reacción al mouse y logo pulpo en el login | `INF-04` | Inferida |
| 392 | 2026-10-07 | `218c1ad` | Nico | Merge pull request #102 from luzpacello/Nico *(merge)* | `INF-04` | Merge |
| 393 | 2026-10-07 | `eb6da06` | Nico | feat(sidebar): avatar por defecto según el rol y logo de la mascota | `FRN-03` | Inferida |
| 394 | 2026-10-07 | `bd11531` | Nico | Merge pull request #103 from luzpacello/Nico *(merge)* | `FRN-03` | Merge |
| 395 | 2026-10-08 | `95e8c1b` | Belinda | ❄ integracion de node/status | `FRN-19B` | Explícita |
| 396 | 2026-10-08 | `fca4d18` | Belinda | Merge pull request #104 from luzpacello/rama-Beli *(merge)* | `FRN-19B` | Merge |
| 397 | 2026-10-08 | `1d429f5` | Luz | 🌱 feat(instances): implementación de confirmación para las acciones operati | `FRN-15` | Explícita |
| 398 | 2026-10-08 | `65a195e` | Luz | Merge pull request #105 from luzpacello/luz *(merge)* | `FRN-15` | Merge |
| 399 | 2026-10-09 | `260e317` | Belinda | ❄ test(dashboard): cubrir errores 502/504 en telemetría del nodo | `FRN-19B` | Mensaje |
| 400 | 2026-10-09 | `ee4644c` | Belinda | Merge pull request #106 from luzpacello/rama-Beli *(merge)* | `FRN-19B` | Merge |
| 401 | 2026-10-09 | `acf79a9` | Cristian | fix(user): distinguir suspensión y eliminación en contrato e interfaz | `FIX-56`, `FRN-06` | Mensaje |
| 402 | 2026-10-09 | `c3d3690` | Cristian | fix(Users): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56`, `FRN-06` | Mensaje |
| 403 | 2026-10-09 | `b014a92` | Cristian | fix(detailsUserPage): Distinguir suspensión y eliminación en contrato e int | `FIX-56`, `FRN-06` | Mensaje |
| 404 | 2026-10-09 | `72a88bd` | Cristian | fix(informationOfUser): Distinguir suspensión y eliminación en contrato e i | `FIX-56`, `FRN-06` | Mensaje |
| 405 | 2026-10-09 | `1202028` | Cristian | fix(CrearUsuarios): Distinguir suspensión y eliminación en contrato e inter | `FIX-56`, `FRN-06` | Mensaje |
| 406 | 2026-10-09 | `0ebccaa` | Cristian | fix(useCreateUser): Distinguir suspensión y eliminación en contrato e inter | `FIX-56`, `FRN-06` | Mensaje |
| 407 | 2026-10-09 | `36e9a36` | Cristian | fix(userDetailsService): Distinguir suspensión y eliminación en contrato e  | `FIX-56`, `FRN-06` | Mensaje |
| 408 | 2026-10-09 | `6571be9` | Cristian | fix(userDesactivationService): Distinguir suspensión y eliminación en contr | `FIX-56`, `FRN-06` | Mensaje |
| 409 | 2026-10-09 | `5ed2f43` | Cristian | fix(useUserInstanceAccess): Distinguir suspensión y eliminación en contrato | `FIX-56`, `FRN-06` | Mensaje |
| 410 | 2026-10-09 | `83adcc8` | Cristian | fix(userAccountMessages): Distinguir suspensión y eliminación en contrato e | `FIX-56`, `FRN-06` | Mensaje |
| 411 | 2026-10-09 | `65482b7` | Cristian | fix(useDeletedUsers): Distinguir suspensión y eliminación en contrato e int | `FIX-56`, `FRN-06` | Mensaje |
| 412 | 2026-10-09 | `d0a1c52` | Cristian | fix(deletedUsersService): Distinguir suspensión y eliminación en contrato e | `FIX-56`, `FRN-06` | Mensaje |
| 413 | 2026-10-09 | `0367985` | Cristian | Merge pull request #107 from luzpacello/cristian *(merge)* | `FIX-56`, `FRN-06` | Merge |
| 414 | 2026-10-09 | `7e4d3af` | Luz | 🌱 fix(instances): arreglo posicionamiento del modal para la acción de elimi | `FRN-15` | Mensaje |
| 415 | 2026-10-09 | `84662bd` | Luz | Merge pull request #108 from luzpacello/luz *(merge)* | `FRN-15` | Merge |
| 416 | 2026-10-09 | `614c912` | Luz | 🌱 feat(instances): implementado la máquina de estados, bloqueo por instanci | `FRN-16`, `FIX-77` | Explícita |
| 417 | 2026-10-09 | `9679081` | Luz | Merge pull request #109 from luzpacello/luz *(merge)* | `FIX-77`, `FRN-16` | Merge |
| 418 | 2026-10-09 | `ef485be` | Belinda | ❄ se completa los filtros reactivos RF-03 | `FRN-20B` | Explícita |
| 419 | 2026-10-09 | `7fcc0bd` | Belinda | Merge pull request #110 from luzpacello/rama-Beli *(merge)* | `FRN-20B` | Merge |
| 420 | 2026-10-13 | `b5da499` | Nico | fix(auth): alinear textos aprobados de autenticación (86e3nrvj2) | `FIX-80` | Explícita |
| 421 | 2026-10-14 | `e7b2d1d` | Nico | fix(dashboard): quitar tarjetas simuladas y controles inertes (86e3pu2e7) | `FIX-90`, `FRN-19A` | Explícita |
| 422 | 2026-10-14 | `4d0b4a7` | Nico | Merge pull request #111 from luzpacello/fix/86e3nrvj2-copy-autenticacion *( | `FIX-80` | Merge |
| 423 | 2026-10-14 | `add947f` | Nico | Merge pull request #112 from luzpacello/fix/86e3pu2e7-limpiar-dashboard *(m | `FIX-90`, `FRN-19A` | Merge |
| 424 | 2026-10-14 | `56342c1` | Cristian | feat(formatInstanceTelemetry): Columnas de CPU/RAM y estado en vivo en el I | `FRN-20A` | Mensaje |
| 425 | 2026-10-14 | `08f1f74` | Cristian | feat(useInstances): Columnas de CPU/RAM y estado en vivo en el Inventario ( | `FRN-20A` | Mensaje |
| 426 | 2026-10-14 | `afd78e5` | Cristian | feat(Instances): Columnas de CPU/RAM y estado en vivo en el Inventario (RF- | `FRN-20A` | Mensaje |
| 427 | 2026-10-14 | `9350a65` | Cristian | Merge pull request #113 from luzpacello/cristian *(merge)* | `FRN-20A` | Merge |
| 428 | 2026-10-14 | `6b38a64` | Cristian | feat(popover): crear nuevo componente para date-picker | `FIX-74`, `FRN-31` | Mensaje |
| 429 | 2026-10-14 | `d9fa6ce` | Cristian | feat(calendar): crear nuevo componente para date-picker | `FIX-74`, `FRN-31` | Mensaje |
| 430 | 2026-10-14 | `0361380` | Cristian | feat(date-picker): crear nuevo componente | `FIX-74`, `FRN-31` | Mensaje |
| 431 | 2026-10-14 | `4a51f64` | Cristian | build(package): instalacion de dependencias para componente calendar | `FIX-74`, `FRN-31` | Mensaje |
| 432 | 2026-10-14 | `9b97227` | Cristian | build(pnpm-lock): instalacion de dependencias para componente calendar | `FIX-74`, `FRN-31` | Mensaje |
| 433 | 2026-10-14 | `a8073f7` | Cristian | fix(Auditoria): Reactividad de filtros y etiquetas accesibles en vista de A | `FIX-74`, `FRN-31` | Mensaje |
| 434 | 2026-10-14 | `d9c2825` | Cristian | fix(date-picker): Reactividad de filtros y etiquetas accesibles en vista de | `FIX-74`, `FRN-31` | Mensaje |
| 435 | 2026-10-14 | `7e3f43d` | Cristian | Merge pull request #114 from luzpacello/cristian *(merge)* | `FIX-74`, `FRN-31` | Merge |
| 436 | 2026-10-14 | `5694c7a` | Cristian | fix(date-picker): cambiar nombre sin seleccionar | `FIX-74`, `FRN-31` | Mensaje |
| 437 | 2026-10-14 | `002ead3` | Cristian | fix(Users): cambiar ancho del dropdowmenu | `FRN-05`, `FRN-06B` | Inferida |
| 438 | 2026-10-14 | `0aa6f8a` | Cristian | Merge pull request #115 from luzpacello/cristian *(merge)* | `FIX-74`, `FRN-05`, `FRN-06B`, `FRN-31` | Merge |
| 439 | 2026-10-14 | `576f65a` | Cristian | feat(badge): Agregar variantes a badge | `FRN-01` | Inferida |
| 440 | 2026-10-14 | `e85cfd3` | Cristian | fix(useUserRecentActivity): Quitar los datos hardcodeados de la ficha de us | `FIX-87`, `FRN-05` | Mensaje |
| 441 | 2026-10-14 | `95b542d` | Cristian | fix(userRecentActivityService): Quitar los datos hardcodeados de la ficha d | `FIX-87`, `FRN-05` | Mensaje |
| 442 | 2026-10-14 | `232e317` | Cristian | fix(detailsUserPages): Quitar los datos hardcodeados de la ficha de usuario | `FIX-87`, `FRN-05` | Mensaje |
| 443 | 2026-10-14 | `757580a` | Cristian | fix(Auditoria): Quitar los datos hardcodeados de la ficha de usuario | `FIX-87`, `FRN-05` | Mensaje |
| 444 | 2026-10-14 | `44d1e3a` | Cristian | Merge pull request #116 from luzpacello/cristian *(merge)* | `FIX-87`, `FRN-01`, `FRN-05` | Merge |
| 445 | 2026-10-14 | `fdcfc22` | Cristian | fix(userInstanceService): agregar persistencia de datos para opciones null | `FRN-20A` | Inferida |
| 446 | 2026-10-14 | `7dca556` | Cristian | feat(detailsUserPage): crear tabla y opcion visual si no cuenta con instanc | `FRN-05`, `FRN-06B` | Inferida |
| 447 | 2026-10-14 | `9697a9f` | Cristian | Merge pull request #117 from luzpacello/cristian *(merge)* | `FRN-05`, `FRN-06B`, `FRN-20A` | Merge |
| 448 | 2026-10-14 | `4e666a0` | Cristian | fix(Instances):  Ruta  con asistente estático sin lógica | `FIX-89`, `FRN-24` | Mensaje |
| 449 | 2026-10-14 | `bfcc5f4` | Cristian | Merge pull request #118 from luzpacello/cristian *(merge)* | `FIX-89`, `FRN-24` | Merge |
| 450 | 2026-10-15 | `57c739f` | Nico | fix(dashboard): reponer paneles y sumar curvas, campana, actividad y alerta | `FRN-19A`, `FRN-19B`, `FRN-19C` | Explícita |
| 451 | 2026-10-15 | `6db04d3` | Nico | Merge pull request #119 from luzpacello/fix/86e3pu2e7-alinear-dashboard *(m | `FRN-19A`, `FRN-19B`, `FRN-19C` | Merge |
| 452 | 2026-10-15 | `6c8ecc9` | Nico | feat(dashboard): desglosar el conteo de instancias por tipo de recurso (86e | `FRN-19B` | Explícita |
| 453 | 2026-10-15 | `405791e` | Nico | fix(dashboard): refrescar el conteo y la actividad en vivo, y alinear el pi | `FIX-83`, `FRN-19C` | Explícita |
| 454 | 2026-10-15 | `3466a46` | Nico | Merge pull request #120 from luzpacello/fix/86e3ptewk-refresco-reactivo *(m | `FIX-83`, `FRN-19B`, `FRN-19C` | Merge |
| 455 | 2026-10-15 | `2b120f4` | Luz | 🌱 feat(instances): implementación de desbloqueo reactivo y notificaciones t | `FRN-17B` | Explícita |
| 456 | 2026-10-15 | `3c9a90b` | Luz | Merge pull request #121 from luzpacello/luz *(merge)* | `FRN-17B` | Merge |
| 457 | 2026-10-15 | `34cf00f` | Cristian | fix(SideBar): cambiar iconos de instancias 'Dashboard' y 'Auditoria' | `FRN-31` | Inferida |
| 458 | 2026-10-15 | `b4a34bc` | Cristian | feat(InstanceCreationSummary): crear componente de resumen de creacion de i | `FRN-26` | Mensaje |
| 459 | 2026-10-15 | `6b85a79` | Cristian | feat(createInstanceOptions): crear datos estaticos para prueba de steps | `FRN-24` | Inferida |
| 460 | 2026-10-15 | `4029e64` | Cristian | feat(useCreateInstanceWizard): Agregar funcionalidades al Wizard | `FRN-24` | Inferida |
| 461 | 2026-10-15 | `944d0d2` | Cristian | feat(createInstance): crear estructura para creacion de instancia | `FRN-24` | Inferida |
| 462 | 2026-10-15 | `fda22d7` | Cristian | feat(selectTypeOfInstances): crear contenido de tipo de instancia | `FRN-24` | Mensaje |
| 463 | 2026-10-15 | `38f215e` | Cristian | feat(createNewInstances): agregarle funcionalidad a la pagina y componentes | `FRN-24` | Inferida |
| 464 | 2026-10-15 | `bf1d581` | Cristian | Merge pull request #122 from luzpacello/cristian *(merge)* | `FRN-24`, `FRN-26`, `FRN-31` | Merge |
| 465 | 2026-10-15 | `7e40c4f` | Belinda | ❄ se agrega filtro por IP y reubica filtros | `FRN-20B` | Mensaje |
| 466 | 2026-10-15 | `1a0cdba` | Belinda | Merge pull request #123 from luzpacello/rama-Beli *(merge)* | `FRN-20B` | Merge |
| 467 | 2026-10-15 | `f61374a` | Belinda | ❄ alinear filtros | `FRN-20B` | Mensaje |
| 468 | 2026-10-15 | `cba3327` | Cristian | feat(configureInstanceResources): crear contenido para la seleccion de recu | `INF-04` | Inferida |
| 469 | 2026-10-15 | `91d1676` | Cristian | feat(createInstanceOptions): agregar texto estatico para simular carga de r | `FRN-24` | Inferida |
| 470 | 2026-10-15 | `93937be` | Cristian | feat(InstanceCreationSummary): agregar funcion para calcular memoria en GB | `INF-04` | Inferida |
| 471 | 2026-10-15 | `c9fdb6b` | Cristian | feat(createInstance): agregar persistencia de datos de almacenamiento | `FRN-24` | Inferida |
| 472 | 2026-10-15 | `3fa0277` | Cristian | feat(useCreateInstanceWizard): agregar paso 2 | `FRN-24` | Inferida |
| 473 | 2026-10-15 | `385a52b` | Cristian | merge: sincronizar cristian con main *(merge)* | — | Merge |
| 474 | 2026-10-15 | `1fc14f4` | Cristian | fix(selectedTypeOfInstances): envolver formulario en card | `FRN-20A` | Inferida |
| 475 | 2026-10-15 | `049f16d` | Belinda | Merge branch 'main' of https://github.com/luzpacello/centinela-front into r | — | Merge |
| 476 | 2026-10-15 | `99aa5dc` | Cristian | fix(createNewInstances): eliminar fondo blanco para todas las opciones | `FRN-24` | Inferida |
| 477 | 2026-10-15 | `8a28b88` | Cristian | Merge pull request #124 from luzpacello/cristian *(merge)* | `FRN-24`, `INF-04` | Merge |
| 478 | 2026-10-15 | `e966818` | Cristian | fix(steppers): cambiar estilo para que entre el texto largo del subtitulo | `FRN-01` | Inferida |
| 479 | 2026-10-15 | `8116e6f` | Cristian | fix(createNewInstances): agregar pagina a paso 2 | `FRN-24` | Inferida |
| 480 | 2026-10-15 | `597e35b` | Cristian | Merge pull request #125 from luzpacello/cristian *(merge)* | `FRN-01`, `FRN-20A`, `FRN-24` | Merge |
| 481 | 2026-10-15 | `ee15fb6` | Belinda | ❄ mensajes traducidos | `INF-04` | Inferida |
| 482 | 2026-10-15 | `6147bc6` | Belinda | Merge pull request #126 from luzpacello/rama-Beli *(merge)* | `FRN-20B`, `INF-04` | Merge |
| 483 | 2026-10-15 | `7271e1a` | Belinda | ❄ se corrige duplicidad de texto saludable inaccesible | `FIX-78`, `FRN-19B` | Mensaje |
| 484 | 2026-10-15 | `fd16cce` | Belinda | Merge pull request #127 from luzpacello/rama-Beli *(merge)* | `FIX-78`, `FRN-19B` | Merge |
| 485 | 2026-10-16 | `623bc28` | Belinda | ❄ se Renderiza botones  accesibles visibles en la cabecera | `FIX-79`, `FRN-20B` | Explícita |
| 486 | 2026-10-16 | `f57d46f` | Belinda | Merge pull request #128 from luzpacello/rama-Beli *(merge)* | `FIX-79`, `FRN-20B` | Merge |
| 487 | 2026-10-16 | `fb58781` | Luz | 🌱 fix(instances): desacoplar guarda de nivelAcceso para rol ADMIN en accion | `FIX-76`, `FRN-15` | Explícita |
| 488 | 2026-10-16 | `c01dd36` | Luz | 🌱 test(instances): integro pruebas para modales | `FRN-15` | Mensaje |
| 489 | 2026-10-16 | `6bb563f` | Luz | Merge pull request #129 from luzpacello/luz *(merge)* | `FIX-76`, `FRN-15` | Merge |
| 490 | 2026-10-16 | `73d2fc7` | Belinda | ❄ FIX: Traducción de errores 403 de instancia | `FIX-86`, `FRN-16` | Explícita |
| 491 | 2026-10-16 | `b005f16` | Belinda | Merge branch 'main' into rama-Beli *(merge)* | — | Merge |
| 492 | 2026-10-16 | `a3c5e91` | Belinda |  ❄ chore: actualizar instancias, tests y config | `INF-04` | Inferida |
| 493 | 2026-10-16 | `2731d64` | Belinda | ❄ fix(instances): corregir tipado de getErrorMessage al eliminar instancia | `FRN-20A` | Inferida |
| 494 | 2026-10-16 | `28a4513` | Belinda | Merge pull request #130 from luzpacello/rama-Beli *(merge)* | `FIX-86`, `FRN-16`, `FRN-20A`, `INF-04` | Merge |
| 495 | 2026-10-16 | `b6e94da` | Cristian | fix(Auditoria): cambiar boton de exportacion a tipo primario y cambiar de p | `FRN-31` | Inferida |
| 496 | 2026-10-16 | `a42ecf1` | Cristian | fix(Instances): crear ruta para crear una nueva instancia | `FRN-20A` | Inferida |
| 497 | 2026-10-16 | `30894f4` | Cristian | fix(createNewInstance): Advertencia visible Próximamente en vista directa / | `FIX-99`, `FRN-24` | Explícita |
| 498 | 2026-10-16 | `68fc4b8` | Cristian | Merge pull request #131 from luzpacello/cristian *(merge)* | `FIX-99`, `FRN-20A`, `FRN-24`, `FRN-31` | Merge |
| 499 | 2026-10-16 | `608d04b` | Cristian | fix(Instances): cambiar tabs y filtro | `FRN-20B` | Mensaje |
| 500 | 2026-10-16 | `2b1720c` | Cristian | fix(InstanceFiltersBar): cambiar tabs y filtro | `FRN-20B` | Mensaje |
| 501 | 2026-10-16 | `10ce750` | Cristian | Merge pull request #132 from luzpacello/cristian *(merge)* | `FRN-20B` | Merge |
| 502 | 2026-10-16 | `b464b7e` | Cristian | fix(instance-filters-bar): acomodar test con nuevo texto de filtros | `FRN-20B` | Mensaje |
| 503 | 2026-10-16 | `15245a0` | Cristian | Merge pull request #133 from luzpacello/cristian *(merge)* | `FRN-20B` | Merge |
| 504 | 2026-10-16 | `773c4b0` | Cristian | fix(Auditoria): cambiar diseño de tabla y representacion de contador de eve | `FRN-31` | Inferida |
| 505 | 2026-10-16 | `143ddda` | Cristian | Merge pull request #134 from luzpacello/cristian *(merge)* | `FRN-31` | Merge |
| 506 | 2026-10-16 | `e85a10f` | Cristian | fix(formatInstanceTelemetry): Normalización de telemetría de inventario (es | `FIX-97`, `FRN-20C` | Explícita |
| 507 | 2026-10-16 | `338d658` | Cristian | fix(Instances): Normalización de telemetría de inventario (escala 0-100, de | `FIX-97`, `FRN-20C` | Explícita |
| 508 | 2026-10-16 | `cd8a1f4` | Cristian | Merge pull request #135 from luzpacello/cristian *(merge)* | `FIX-97`, `FRN-20C` | Merge |
| 509 | 2026-10-16 | `041b911` | Luz | 🌱 fix (instance): corregir cancelación de la confirmación de start en los m | `FIX-85`, `FRN-15` | Explícita |
| 510 | 2026-10-16 | `46cd4b3` | Luz | Merge pull request #136 from luzpacello/luz *(merge)* | `FIX-85`, `FRN-15` | Merge |
| 511 | 2026-10-18 | `f641c5f` | Luz | 🌱feat(instance): integración de vista de detalle y navegación | `FRN-21`, `FIX-95` | Explícita |
| 512 | 2026-10-18 | `94f4db8` | Cristian | fea(NetworkModeOption): crear paso de red en creacion de instancias | `FRN-24`, `FRN-25` | Mensaje |
| 513 | 2026-10-18 | `8a8bb78` | Cristian | feat(configureInstanceNetwork): crear paso de Red, en creacion de instancia | `FRN-24`, `FRN-25` | Mensaje |
| 514 | 2026-10-18 | `39f5c3a` | Cristian | feat(createInstance): crear paso de Red, en creacion de instancias | `FRN-24`, `FRN-25` | Mensaje |
| 515 | 2026-10-18 | `4a59cbb` | Luz | Merge pull request #137 from luzpacello/luz *(merge)* | `FIX-95`, `FRN-21` | Merge |
| 516 | 2026-10-18 | `4d603aa` | Cristian | feat(createInstanceOptions): crear paso de Red, en creacion de instancias | `FRN-24`, `FRN-25` | Mensaje |
| 517 | 2026-10-18 | `74cb5d1` | Cristian | feat(useCreateInstanceWizard): crear paso de Red, en creacion de instancias | `FRN-24`, `FRN-25` | Mensaje |
| 518 | 2026-10-18 | `21b8522` | Cristian | feat(createNewInstances): crear paso de Red, en creacion de instancias | `FRN-24`, `FRN-25` | Mensaje |
| 519 | 2026-10-18 | `ea4f1c9` | Cristian | Merge pull request #138 from luzpacello/cristian *(merge)* | `FRN-24`, `FRN-25` | Merge |
| 520 | 2026-10-18 | `c51dc9d` | Cristian | feat(selectInstanceOperatingSystem): crear paso de SO, en creacion de insta | `FRN-24` | Mensaje |
| 521 | 2026-10-18 | `6cd2797` | Cristian | feat(creteInstanceOptions): crear paso de SO, en creacion de instancias | `FRN-24` | Mensaje |
| 522 | 2026-10-18 | `dd203f6` | Cristian | feat(createInstance): crear paso de SO, en creacion de instancias | `FRN-24` | Mensaje |
| 523 | 2026-10-18 | `332e177` | Cristian | feat(createNewInstances): crear paso de SO, en creacion de instancias | `FRN-24` | Mensaje |
| 524 | 2026-10-18 | `ce676cc` | Cristian | Merge pull request #139 from luzpacello/cristian *(merge)* | `FRN-24` | Merge |
| 525 | 2026-10-18 | `bd62aa9` | Cristian | feat(reviewInstanceConfiguration): crear paso de Resumen, en creacion de in | `FRN-26` | Mensaje |
| 526 | 2026-10-18 | `7fe3218` | Cristian | feat(createInstance): crear paso de Resumen, en creacion de instancias | `FRN-26` | Mensaje |
| 527 | 2026-10-18 | `a4191b9` | Cristian | feat(useCreateInstanceWizard): crear paso de Resumen, en creacion de instan | `FRN-26` | Mensaje |
| 528 | 2026-10-18 | `0664b94` | Cristian | feat(instanceService): crear paso de Resumen, en creacion de instancias | `FRN-26` | Mensaje |
| 529 | 2026-10-18 | `3510240` | Cristian | feat(createNewInstances): crear paso de Resumen, en creacion de instancias | `FRN-26` | Mensaje |
| 530 | 2026-10-18 | `9e0683b` | Cristian | Merge branch 'main' into cristian *(merge)* | — | Merge |
| 531 | 2026-10-18 | `3d840ad` | Cristian | Merge pull request #140 from luzpacello/cristian *(merge)* | `FRN-26` | Merge |
| 532 | 2026-10-18 | `48a4406` | Nico | fix(dashboard): sacar el ruido del canal de eventos de la actividad recient | `FRN-19C`, `FIX-83` | Explícita |
| 533 | 2026-10-18 | `54efc3d` | Nico | feat(dashboard): cargar el histórico del nodo al entrar en el dashboard | `FRN-19C` | Explícita |
| 534 | 2026-10-18 | `d3643e0` | Nico | Merge remote-tracking branch 'origin/main' into fix/actividad-reciente-sin- | — | Merge |
| 535 | 2026-10-18 | `1913ca0` | Nico | fix(dashboard): alinear los detalles de auditoría con el contrato en objeto | `FRN-19C`, `FRN-31` | Explícita |
| 536 | 2026-10-18 | `91eabf4` | Nico | Merge pull request #141 from luzpacello/fix/actividad-reciente-sin-ruido *( | `FIX-83`, `FRN-19C`, `FRN-31` | Merge |

---

## 7. Conclusiones y balance de autoría

1. **Fuerte asimetría productiva:** el repositorio `frontend` dependió centralmente de **Cristian**, quien generó el **71.7% del código real** (271 commits).
2. **Inexactitud de las asignaciones documentadas:** más del **55% de las tareas** no fueron ejecutadas por la persona que figuraba en la columna de asignación en los documentos de seguimiento (`actual.md`, `terminado.md`, ClickUp).
3. **El fenómeno de los 'rescates técnicos':** cuando se generaron bloqueos graves (como la anulación accidental de seguridad de rutas en `6b08566` o la complejidad de la arquitectura SSE en tiempo real), fue **Nico** quien asumió el desarrollo desde backend para solucionar el frontend con 50 commits de código de alta complejidad.
4. **Diferenciación entre Maquetado/Gestión y Lógica de Negocio:** mientras que **Belinda** y **Luz** aportaron en etapas muy tempranas estructuras visuales y plantillas, la lógica reactiva, persistencia, validaciones de formularios, llamadas API, gestión de estado y wizards fueron casi en su totalidad construidos por Cristian y corregidos por Nico.
