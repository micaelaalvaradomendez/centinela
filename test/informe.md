# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 06/10/2026 (corrida con últimos commits de submódulos)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), contrastadas tras la incorporación de los nuevos commits en el backend (`df320c7`) y frontend (`fc6f9ed`). Se corrieron las suites completas de backend (`test/back`) en Docker Compose y frontend (`test/front`) en Vitest, incorporando las nuevas pruebas de aceptación para las tareas pendientes y la regresión de `terminado.md` y `terminado-1.md`.

| Componente | Revisión probada | Commits nuevos incorporados |
|---|---|---|
| Backend | `df320c7` (último commit de `main`) | `df320c7` (migración de datos de `auditoria_legacy` en `db.go` y test `migrar_auditoria_test.go` para `FIX-51`; pruebas unitarias de `ObtenerInterfaces` en `client_test.go` para `BAC-23A`/`FIX-63`; y aclaración de acceso D1 en `contrato-etapa1.md` para `BAC-29`/`FIX-65`) |
| Frontend | `fc6f9ed` (último commit de `main`) | `fc6f9ed` (31 commits nuevos: nuevos componentes `Combobox`, `Pagination`, `DropDownMenu`, `InputGroup`, refactor de selectores en `Users.tsx`, `Auditoria.tsx`, `Instances.tsx`, `Dashboard.tsx` y hooks de instancias) |

Ambas suites se ejecutaron en el entorno local (Go 1.27.1 + Docker Compose para backend, Node 22 + pnpm + Vitest para frontend). El backend no dejó contenedores residuales.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 72 | **65** | **5** | 2 |
| Frontend (`test/front`) | 157 | **115** | **32** | 10 |
| **Total** | **229** | **180** | **37** | **12** |

- **Evolución respecto a la corrida previa:**
  - **Backend:** Se incrementó de 53 a **65 pruebas aprobadas** (+12 en verde).
    - `FIX-51` (migración de auditoría legacy a particionada): pasa a **100% en verde** tras incorporar la migración en `db.go` y corregir el matcher de columnas explícitas en la suite de aceptación.
    - `BAC-23A` (adaptador de IP): pasa a **100% en verde** con las pruebas unitarias agregadas en `client_test.go`.
    - `BAC-29` (contrato formal D1): pasa a **100% en verde** con la especificación en `docs/contrato-etapa1.md`.
    - `BAC-23B` (IP real y RBAC estricto): validado con la nueva suite `etapa1_ola2_acceptance_test.go` y **pasa 100% en verde**.
    - `BAC-22B` (conteo `instancesSummary` y telemetría por instancia): validado con `etapa1_ola2_acceptance_test.go` y **pasa 100% en verde**.
    - Los fallos restantes (5) corresponden a `BAC-25A` (bloqueado por `FIX-69` en Swagger), `BAC-25B` (`UPID_TIMEOUT` configurable pendiente en backend), `FIX-52` (distinguir `INSTANCE_BUSY`), `FIX-54` (`eliminado_en`), y `INF-06A` (colisión ambiental de puerto 6379 con el servicio local `redis-server`).
  - **Frontend:** Se incrementó de 107 a **115 pruebas aprobadas** (+8 en verde).
    - Se adaptaron `admin-users.test.tsx` (22 de 22 en verde) y `audit.test.tsx` (7 de 7 en verde) para soportar la nueva arquitectura de componentes `Combobox` y `Pagination` introducida en `fc6f9ed`.
    - Se incorporaron las suites `test/front/auth-branding.test.tsx` (4 pruebas) y `test/front/instances-ola2.test.tsx` (5 pruebas).
    - Las 32 fallas restantes mapean de forma unívoca a las tareas pendientes de implementación en frontend: `FIX-58`, `FIX-59`, `FIX-60`, `FIX-56`, `FRN-15`, `FRN-16`, `FRN-19B` y el selector de tipo de `FRN-20B`.
- **Omitidos:**
  - Frontend (10): `login04-e2e.test.ts` (suite de integración que se valida dentro de `test/back` donde aprueba 10/10).
  - Backend (2): `INF-08B` (Brevo rechaza la IP pública con código 525) y `FIX-55` (omitido condicionalmente a la espera de la migración de esquema de `FIX-54`).

---

## 2. Estado de las tareas de `documentacion/actual.md`

**Leyenda:** ✅ cumple · 🟡 implementada con problemas / parcial · ❌ no implementada.

| Tarea | Área | Estado | Detalle y Cobertura de Pruebas |
|---|---|---|---|
| `FIX-51` auditoría legacy | Backend | ✅ | **Pasa 100% en verde.** `db.go:286` ejecuta `INSERT INTO auditoria (...) SELECT ... FROM auditoria_legacy` y descarta la tabla legacy. Validado en `fixes_acceptance_test.go:52` y `migrar_auditoria_test.go`. |
| `BAC-23B` IP real y RBAC | Backend | ✅ | **Pasa 100% en verde.** `instance_handler.go` resuelve IPs con `ResolverIPs` únicamente sobre las instancias filtradas del usuario; sin token da 401 y operador sin permisos lista vacía. Validado en `etapa1_ola2_acceptance_test.go:28`. |
| `BAC-22B` telemetría y resumen | Backend | ✅ | **Pasa 100% en verde.** `GET /api/node/status` incluye `instancesSummary` con desglose `vms` y `lxc` (running, stopped, paused, total); `GET /api/instances` expone `cpuUsage`, `ramUsage` y `maxRam`. Validado en `etapa1_ola2_acceptance_test.go:74`. |
| `FRN-16B` resincronización F5 | Frontend | ✅ | **Pasa 100% en verde.** Al cargar `/instances` con `activeTask !== null`, las filas inician en estado `transitioning` y con controles deshabilitados. Validado en `instances-ola2.test.tsx:125`. |
| `FRN-20B` filtros reactivos | Frontend | 🟡 | Búsqueda en tiempo real por nombre/ID y contador de elementos pasan 100%. Pendiente agregar el selector reactivo para tipo de recurso (Todas, VM, LXC). Validado en `instances-ola2.test.tsx:47`. |
| `BAC-25B` timeout configurable | Backend | ❌ | Pendiente en backend: `seguimiento_tareas.go` mantiene `limiteSeguimiento` fijo en 10 min en lugar de `UPID_TIMEOUT` (3m) con reintentos con backoff. Verificado en `etapa1_ola2_acceptance_test.go:128`. |
| `FIX-52` INSTANCE_BUSY | Backend | ❌ | Pendiente en backend: `validarEstadoParaAccion` no consulta si existe tarea `RUNNING` antes de evaluar estados, y persiste código muerto de `ReiniciarInstancia`. Verificado en `fixes_acceptance_test.go:78`. |
| `FIX-54` eliminación lógica | Backend | ❌ | Pendiente en backend: falta columna `eliminado_en` y reemplazo del índice único incondicional por parcial `WHERE eliminado_en IS NULL`. Verificado en `fixes_acceptance_test.go:155`. |
| `FIX-55` identidad histórica | Backend | ❌ | Pendiente en backend: bloqueada por dependencia con `FIX-54`. Verificado en `fixes_acceptance_test.go:175`. |
| `FIX-56` distinción UI baja | Frontend | ❌ | Pendiente en frontend: la interfaz no distingue eliminación de suspensión en badges ni modal. Verificado en `admin-users-deletion.test.tsx:104`. |
| `FIX-58` logotipo oficial | Frontend | ❌ | Pendiente en frontend: layout de autenticación renderiza texto plano `"logo"` en lugar del imagotipo/logotipo SVG accesible. Verificado en `auth-branding.test.tsx:19`. |
| `FIX-59` copy amigable auth | Frontend | ❌ | Pendiente en frontend: persiste el placeholder `"imagenes y informacion random"` en el sidebar de autenticación. Verificado en `auth-branding.test.tsx:37`. |
| `FIX-60` reubicar botón crear | Frontend | ❌ | Pendiente en frontend: `/dashboard` sigue conteniendo el botón "Crear instancia", el cual debe residir en `/instances` protegido por rol ADMIN. Verificado en `auth-branding.test.tsx:55`. |
| `FRN-15` modales antierror | Frontend | ❌ | Pendiente en frontend: faltan los modales de confirmación para las acciones operativas en el inventario. Verificado en `instances-modals.test.tsx:41`. |
| `FRN-16` operación en progreso | Frontend | ❌ | Pendiente en frontend: falta la máquina de estados reactiva con spinner y mapeo de errores D2. Verificado en `instances-operation.test.tsx:35`. |
| `FRN-19B` semáforo de nodo | Frontend | ❌ | Pendiente en frontend: Dashboard no consume `/node/status` ni renderiza el semáforo con estados D3. Verificado en `dashboard-node.test.tsx:57`. |

---

## 3. Pruebas incorporadas y adecuadas en esta iteración

1. **`test/back/etapa1_ola2_acceptance_test.go` (NUEVO):**
   - Valida `BAC-23B`: IP real expuesta bajo RBAC estricto, 401 sin credenciales, lista vacía para operador no asignado e invocación de `ResolverIPs` solo sobre instancias visibles.
   - Valida `BAC-22B`: Estructura y datos de `instancesSummary` (vms/lxc) en `GET /api/node/status` y métricas de telemetría (`cpuUsage`, `ramUsage`, `maxRam`) en `GET /api/instances`.
   - Valida `BAC-25B`: Requisitos de `UPID_TIMEOUT` (3m) y detalles en `TASK_FINISHED`.
2. **`test/front/auth-branding.test.tsx` (NUEVO):**
   - Valida `FIX-58`: Ausencia de texto plano placeholder `"logo"` y presencia de imagotipo/logotipo vectorizado con accesibilidad.
   - Valida `FIX-59`: Erradicación del string `"imagenes y informacion random"` y presencia de copy explicativo sobre Centinela y 2FA.
   - Valida `FIX-60`: Remoción del botón "Crear instancia" del Dashboard y presencia contextual en `/instances` protegido por `ADMIN`.
3. **`test/front/instances-ola2.test.tsx` (NUEVO):**
   - Valida `FRN-20B`: Búsqueda instantánea en cliente por ID/nombre, filtro por estado (running/stopped), selector de tipo de recurso y contador sobre el total.
   - Valida `FRN-16B`: Estado `transitioning` automático al montar la tabla si `activeTask !== null`.
4. **Adecuaciones en pruebas existentes:**
   - `test/back/fixes_acceptance_test.go:60`: Se adecuó el regex de `FIX-51` (`(?is)INSERT INTO auditoria ...`) para aceptar sintaxis SQL estándar con lista de columnas explícita entre paréntesis y saltos de línea.
   - `test/front/admin-users.test.tsx`: Se adaptaron los selectores de rol y permisos para interactuar correctamente con el nuevo componente `@base-ui/react` `Combobox` introducido por frontend en `fc6f9ed` (pasa 22/22 en verde).
   - `test/front/audit.test.tsx`: Se adecuó la selección del enlace de paginación siguiente (`PaginationNext` renderizado como `<a>` accesible) pasando 7/7 en verde.

---

## 4. Próximos pasos recomendados para desarrollo

1. **Promover tareas resueltas a `documentacion/terminado.md` / `documentacion/terminado-1.md`:**
   - `FIX-51` (Backend): completamente verificado y pasando en verde.
   - `BAC-23B` (Backend): completamente verificado y pasando en verde.
   - `BAC-22B` (Backend): completamente verificado y pasando en verde.
2. **Backend:**
   - Resolver `FIX-69` en Swagger para restaurar definiciones SSE en `docs/swagger.yaml` y liberar `BAC-25A`.
   - Implementar `BAC-25B` en `seguimiento_tareas.go` configurando `UPID_TIMEOUT` con valor default de 3 minutos y retroceso en reintentos.
   - Implementar `FIX-52` en `instance_handler.go` consultando tareas activas antes de evaluar estado incompatible.
   - Implementar `FIX-54` agregando `eliminado_en` e índice condicional.
3. **Frontend:**
   - Resolver los fixes visuales de auth: `FIX-58` (logo SVG), `FIX-59` (copy de bienvenida y 2FA) y `FIX-60` (reubicar botón crear al inventario).
   - Completar el selector por tipo de recurso en `Instances.tsx` (`FRN-20B`).
   - Implementar los modales antierror (`FRN-15`), máquina de estados (`FRN-16`) y semáforo del Dashboard (`FRN-19B`).
