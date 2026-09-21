# Informe de estado actual — login, 2FA e integración
# Informe de estado actual — login, 2FA, administración e integración

Fecha: 18/09/2026.
Fuentes: [documentacion/actual.md](documentacion/actual.md), [test/back](test/back) (suite Go, ejecución real contra PostgreSQL) y [test/front](test/front) (suite Vitest sobre el código real del frontend).
Fecha: 20/09/2026.  
Fuentes: [documentacion/actual.md](documentacion/actual.md), [test/back](test/back) (suite Go, ejecución real con Docker Compose y PostgreSQL 16) y [test/front](test/front) (suite Vitest sobre el código real del frontend).

Commits auditados (coinciden con `origin/main` de cada submódulo):

| Componente | Commit |
|---|---|
| Frontend | `7775b62` |
| Backend | `caf71c2` |
| Frontend | `68bd23f` |
| Backend | `52ebce1` |

---


## 1. Resumen ejecutivo

- **Actualización de submódulos:** ambos submódulos están al día con sus ramas `main` (`backend@caf71c2`, `frontend@7775b62`); el último pull no modificó ningún archivo de `BAC-05` a `BAC-14` ni de `FRN-05` a `FRN-08`.
- **Frontend administrativo (`FRN-05`, `FRN-06`, `FRN-06B`): sin implementar.** Existen pantallas visuales estáticas ([Users.jsx](frontend/centinela/src/pages/Users.jsx), [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx)), pero sin guard de rol, sin llamadas a la API, sin selector de rol ni envío de formulario. 7 pruebas reales nuevas en [test/front/admin-users.test.tsx](test/front/admin-users.test.tsx) lo confirman: **7 de 7 fallan**. Detalle en 3.1.
- **Backend administrativo (`BAC-05`, `BAC-06`, `BAC-06B`, `BAC-09`): lógica funcional correcta, con un bug real de duplicados y un desvío de ruta.** El CRUD de usuarios, la contraseña temporal y el RBAC funcionan correctamente en las cinco operaciones (crear, listar, ver detalle, editar, dar de baja). Dos hallazgos distintos:
  1. **Bug funcional (no de nomenclatura):** usar un `id` propio como PK en `permisos_instancia` en vez de una compuesta es válido por sí solo. El problema real es que no existe **ninguna** restricción de unicidad sobre `(usuario_id, vmid_proxmox)`. Se demostró en vivo: enviar `{"vmids": [201, 201]}` persiste **2 filas duplicadas** en vez de 1. Afecta a `BAC-05` y, por compartir repositorio, también a `BAC-07`.
  2. **Desvío de ruta (`BAC-06`/`BAC-06B`):** en [backend/cmd/api/main.go](backend/cmd/api/main.go) la variable del grupo de rutas se llama `admin`, pero se monta con `Group("/")` en vez de `Group("/admin")`, por lo que `GET /api/admin/users` no existe (solo `/api/users`). El resto de la API (`/api/auth/*`, `/api/roles`) sí coincide con `actual.md`, lo que sugiere un detalle de implementación pendiente, no un cambio de diseño deliberado.
- **Actualización de submódulos:**
  - **Backend (`52ebce1`):** Se avanzaron 3 commits respecto al informe previo (`3052a02`, `adcf261`, `52ebce1`). El commit `adcf261` corrigió el montaje del router de administración bajo `/api/admin/users`. Sin embargo, el commit `3052a02` introdujo una **regresión crítica** (`FIX-01`): reutilizó el tag `uniqueIndex:idx_usuario_vmid` en `SesionActiva.UsuarioID`, provocando que PostgreSQL cree un índice único sobre la sesión. En consecuencia, el segundo login de un usuario o el salto inmediato a verificar TOTP falla con `401 AUTH_FAILED` (`SQLSTATE 23505 duplicate key`). Además, como el nombre del índice colisionó, la restricción de unicidad sobre `permisos_instancia` (`BAC-05`) **nunca se creó en la base**.
  - **Frontend (`68bd23f`):** Se avanzaron 5 commits (`acc5ecf`, `12b0779`, `6b08566`, `c7a8800`, `68bd23f`). Se agregaron componentes y maquetas ricas en TypeScript ([Users.tsx](frontend/centinela/src/pages/Users.tsx), [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) con selector de rol, [UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx), [Auditoria.tsx](frontend/centinela/src/pages/Auditoria.tsx)), y se resolvió `FIX-07` en [apiClient.ts](frontend/centinela/src/services/apiClient.ts) soportando `errorCode`. No obstante, el commit `6b08566` introdujo una **regresión crítica de producto** (`FIX-08`): movió `/dashboard`, `/instances`, `/users`, `/users/new`, etc., fuera de `ProtectedLayout` al layout público (`MainLayoutAuth`) como "rutas temporales de diseño", **desprotegiendo el acceso al sistema sin autenticación ni 2FA** y rompiendo `FRN-03`.
- **Estado de las suites de prueba:**
  - **Frontend (`test/front`):** Se crearon y activaron todas las pruebas pendientes (eliminando todos los `it.todo`). Se incorporó `api-client.test.ts` (3/3 aprobadas). Resultado actual: **19 aprobadas, 12 fallidas, 0 pendientes (31 pruebas totales)**. Las 12 fallas documentan la rotura del guard de rutas (`navigation.test.tsx`, 2 fallas) y la falta de conexión real con la API en administración y recursos (`admin-users.test.tsx`, 10 fallas).
  - **Backend (`test/back`):** Se activaron las pruebas del hito de recursos en [resource_access_acceptance_test.go](test/back/resource_access_acceptance_test.go) (eliminando los `t.Skip`). Resultado actual: **3 aprobadas, 10 fallidas**. Las pruebas que aprueban son `BAC-01` (esquema inicial y seed), `BAC-02` (bcrypt) y `BAC-04` (errores HTTP 400/401). Las 10 fallas reflejan la regresión crítica de login (`BAC-03`), la falta de unicidad en `BAC-05`, la cascada provocada por falta de token (`LOGIN-01`, `BAC-09`, `BAC-10/11`, `LOGIN-03`) y la ausencia de los endpoints de recursos (`BAC-07`, `BAC-08`, `BAC-14`).

  3 pruebas reales nuevas en [test/back/backend_acceptance_test.go](test/back/backend_acceptance_test.go) confirman ambos hallazgos. Detalle en 3.2.
- **Hito "Gestión Administrativa de Usuarios": no cumplido de punta a punta.** Ver 3.3.
- **Hito "Control de Acceso Basado en Recursos": sin empezar (`BAC-08`, `BAC-14`, `FRN-07`, `FRN-08`).** `BAC-07` tiene una implementación parcial (`PUT /api/users/:id/instances`) que funciona pero con el mismo desvío de prefijo de ruta.
- **Nuevas tareas 2FA (`BAC-10`, `BAC-11`, `FRN-09`):** hay implementación real, no solo preparación. El backend genera QR, cifra el secreto y permite login posterior; el frontend `LoginContinuation` muestra QR, clave manual y OTP. Las pruebas nuevas detectan un fallo de BAC-10: un usuario ya vinculado puede regenerar su secreto llamando nuevamente a `GET /api/auth/2fa/qr` sin pasar por reset/relink. FRN-09 tiene una prueba enfocada aprobada.
- **Historial de correcciones de este informe:** el 18/09/2026 se detectó que los criterios `FRN-05`/`FRN-06`/`FRN-06B` estaban declarados como `it.todo` (Vitest nunca los ejecuta) y que los tests backend de `BAC-05`/`BAC-06`/`BAC-06B` no verificaban dos requisitos textuales de `actual.md`. Ambos huecos se corrigieron con pruebas reales que hoy fallan, y ese es el estado reflejado en este informe.
- **Nota sobre la base de datos:** el repositorio no tiene ningún archivo `.sql` con `CREATE TABLE` (`backend/scripts/init.sql` solo define la función `uuid_generate_v7()`). El esquema real lo genera `gorm.AutoMigrate` en tiempo de ejecución a partir de las structs de [backend/internal/core/domain/models.go](backend/internal/core/domain/models.go), y se verifica consultando en vivo el PostgreSQL 16 que levanta `test/back` con Docker Compose (no es una base simulada).
---

## 2. Estado por tarea

| ID | Estado | Evidencia |
|---|---|---|
| `BAC-01` | ✅ Completa | El compose aislado de [test/back](test/back) levanta PostgreSQL, GORM automigra el esquema (tablas `organizaciones`, `usuarios`, `sesiones_activas`) y crea el admin seed con roles `ADMIN`/`OPERATOR`. |
| `BAC-02` | ✅ Completa | Hash bcrypt verificado y login válido e inválido comprobados contra PostgreSQL. |
| `BAC-03` | ✅ Completa | `POST /api/auth/login` emite un JWT temporal válido usando el campo `password`. |
| `BAC-04` | ✅ Completa | Payload incompleto e inválido devuelven `400` y `401` con el contrato JSON esperado. |
| `FRN-01` | ✅ Completa | Formulario de login (usuario/correo, contraseña, envío, estado de carga) validado por [test/front](test/front). |
| `FRN-02` | ✅ Completa | Validaciones de cliente bloquean campos vacíos y correo inválido antes de enviar. |
| `FRN-03` | ✅ Completa | Protección de rutas y navegación Dashboard -> Instancias comprobadas con los encabezados actuales. |
| `FRN-04` | ✅ Completa | El frontend envía `password` y acepta la respuesta pre-2FA emitida por el backend. |
| `LOGIN-01` | ✅ Backend implementado | El backend expone `GET /api/auth/2fa/qr` y `POST /api/auth/2fa/verify`, cifra el secreto con AES-256-GCM y valida códigos TOTP. |
| `LOGIN-02` | ✅ Completa | El contrato de QR/verificación y la validación numérica del OTP pasan. |
| `LOGIN-03` | ✅ En integración activa | El flujo de autenticación conecta login con redirección a 2FA ([test/front/authentication-contract.test.ts](test/front/authentication-contract.test.ts) validado). |
| `BAC-10` | ❌ No completa | `GET /api/auth/2fa/qr` genera QR y secreto, pero no rechaza la solicitud cuando `totp_vinculado=true`; se demostró con una prueba de login posterior. |
| `BAC-11` | ⚠️ Parcial | El secreto se cifra y el login posterior consulta el secreto persistido; la prueba también prepara replay, pero el escenario queda bloqueado por la regeneración indebida del QR. |
| `FRN-09` | ✅ Parcialmente validada | La UI real muestra QR, clave manual formateada y formulario OTP en `LoginContinuation`; falta validar el recorrido completo contra BAC-10/BAC-11 verde. |
| `BAC-05` | ❌ No completa — **bug funcional, no de nomenclatura** | El alta de usuarios y sus campos de seguridad funcionan. Usar un `id` propio como PK en `permisos_instancia` en vez de una compuesta es válido; lo que falta es cualquier restricción de unicidad sobre `(usuario_id, vmid_proxmox)`. Demostrado en vivo: `{"vmids": [201, 201]}` persiste **2 filas** en vez de 1. |
| `BAC-06` | ⚠️ Parcial — **contrato de ruta incumplido** | El CRUD, la contraseña temporal y la restricción de `OPERATOR` funcionan bajo `/api/users`. `GET /api/admin/users` devuelve `404`: en `cmd/api/main.go` la variable se llama `admin` pero el `Group()` usa `"/"` en lugar de `"/admin"`. |
| `BAC-06B` | ⚠️ Parcial — **mismo desvío** | `PUT /api/users/:id` actualiza nombre, correo, rol y estado (`activo`) correctamente, y un `OPERATOR` recibe `403` al consultar, editar o dar de baja; el mismo detalle de montaje de rutas afecta a `/api/admin/users/:id`. |
| `BAC-09` | ✅ Completa | `GET /api/roles` coincide exactamente con la ruta documentada (sin prefijo `/admin`, tal como pide `actual.md`). |
| `FRN-05` | ❌ No completa | 3 pruebas reales, **3 fallan**: la ruta `/users` no bloquea a `OPERATOR` (sin guard de rol), la vista no llama a `GET /api/users` y la tabla ignora cualquier dato real (siempre muestra "Sin usuarios"). Ver 3.1. |
| `FRN-06` | ❌ No completa | 3 pruebas reales, **3 fallan**: no hay selector de rol, el botón "Crear usuario" no dispara ningún `POST` (no hay `<form>` ni `onSubmit`) y no se muestra la contraseña temporal. Ver 3.1. |
| `FRN-06B` | ❌ No completa | 1 prueba real, **falla**: no existe ninguna acción de edición por fila porque la tabla nunca renderiza usuarios. Ver 3.1. |
| `BAC-07` | ⚠️ Parcial | La implementación actual permite reemplazar permisos mediante `PUT /api/users/:id/instances` y leerlos en `GET /api/users/:id`; falta el contrato documentado `GET/PUT /api/admin/users/:id/permissions`. |
| `BAC-08` | ⏳ Preparada | La prueba está reservada, pero aún no existe el guard por recurso ni una operación de instancia que permita comprobar el `403` sin llamada a Proxmox. |
| `BAC-14` | ⏳ Preparada | La prueba está reservada, pero aún no existe `GET /api/instances` ni el contrato normalizado del adaptador Proxmox. |
| `FRN-07` | ⏳ Preparada | Hay criterios frontend `todo` para listado, selección y persistencia de VMIDs; todavía no existe la vista. |
| `FRN-08` | ⏳ Preparada | Hay criterios frontend `todo` para Bearer y `403` sin cerrar sesión; falta la integración de recursos. |
| `BAC-01` | ✅ Completa | El compose aislado levanta PostgreSQL 16, GORM automigra y crea admin seed con roles `ADMIN`/`OPERATOR`. |
| `BAC-02` | ✅ Completa | Hash bcrypt verificado; login distingue credenciales válidas e inválidas. |
| `BAC-03` | ❌ Rota por regresión (`FIX-01`) | La lógica existe, pero `POST /api/auth/login` falla en el segundo intento con `401 duplicate key value violates unique constraint "idx_usuario_vmid"` en `sesiones_activas`. |
| `BAC-04` | ✅ Completa | Payload incompleto e inválido devuelven `400` y `401` con contrato JSON unificado (`errorCode`). |
| `FRN-01` | ✅ Completa | Formulario de login (campos, envío, estado de carga) validado por `test/front`. |
| `FRN-02` | ✅ Completa | Validaciones de cliente bloquean campos vacíos y formato de correo inválido antes de enviar. |
| `FRN-03` | ❌ Rota por regresión (`FIX-08`) | El commit `6b08566` sacó las rutas del `ProtectedLayout`; cualquiera puede acceder a `/dashboard`, `/instances` y `/users` sin sesión. |
| `FRN-04` | ✅ Completa | El frontend envía payload compatible (`email`, `password`) y procesa la respuesta pre-2FA. |
| `LOGIN-01` | ⚠️ Bloqueada | El backend expone `/api/auth/2fa/qr` y `/api/auth/2fa/verify`, pero la verificación falla por la colisión de sesión en base. |
| `LOGIN-02` | ✅ Completa | Contrato de QR/verificación y validación numérica de 6 dígitos probados y aprobados. |
| `LOGIN-03` | ⚠️ Bloqueada | La integración entre login y 2FA está implementada en el front, pero bloqueada de validar punta a punta por la regresión de sesión en el back. |
| `BAC-10` | ❌ No completa (`FIX-03`) | `GET /api/auth/2fa/qr` no rechaza solicitudes cuando `totp_vinculado=true`; permite regenerar el secreto sin pasar por reset. |
| `BAC-11` | ⚠️ Parcial | Cifrado AES-GCM implementado; la verificación completa de replay queda bloqueada tras la falla de BAC-10/BAC-03. |
| `FRN-09` | ✅ Completa | `LoginContinuation` renderiza QR, clave manual formateada y formulario OTP de 6 dígitos. |
| `BAC-05` | ❌ No completa | Falta restricción de unicidad sobre `(usuario_id, vmid_proxmox)`. La colisión de nombre de índice impidió la creación del constraint en Postgres. |
| `BAC-06` | ⚠️ Parcial | El router ya está montado bajo `/api/admin/users`, pero su validación punta a punta está bloqueada por la rotura de login. |
| `BAC-06B` | ⚠️ Parcial | `PUT /api/admin/users/:id` implementado en backend; bloqueado de validar por login. |
| `BAC-09` | ⚠️ Parcial | `GET /api/roles` expuesto en backend; requiere token válido para responder 200. |
| `FRN-05` | ❌ No completa (`FIX-05`) | Vista maquetada en `Users.tsx`, pero usa datos mockeados fijos (`usersList`), no llama a `GET /api/admin/users` y no tiene guard de rol. |
| `FRN-06` | ❌ No completa (`FIX-06`) | `CrearUsuarios.tsx` tiene selector visual de rol, pero los botones no disparan `POST /api/admin/users`, no muestran la contraseña temporal y no hay baja funcional. |
| `FRN-06B` | ❌ No completa (`FIX-06`) | `UserDetail.tsx` maquetado, pero el botón "Guardar cambios" no envía `PUT /api/admin/users/:id` a la API. |
| `BAC-07` | ❌ No completa (`FIX-02`) | Expone `PUT /api/admin/users/:id/instances` en vez de `/permissions`; no existe `GET /permissions` aislado y persiste duplicados. |
| `BAC-08` | ❌ No implementada | No existe middleware de autorización por recurso ni endpoints protegidos de instancias (responde 404). |
| `BAC-14` | ❌ No implementada | No existe `GET /api/instances` ni normalización de Proxmox en backend (responde 404). |
| `FRN-07` | ❌ No implementada | No existe selector interactivo conectado a `GET /api/instances` ni envío del array de permisos. |
| `FRN-08` | ❌ No implementada | No existe interceptor global que maneje respuestas `403` en recursos preservando la sesión. |
| `FIX-07` | ✅ Completa | `apiClient.ts` extrae correctamente el campo `errorCode` emitido por el backend (3 pruebas aprobadas). |
| `FIX-08` | ❌ Pendiente | Falta restaurar los guards de sesión en `applicationRoutes.tsx` para volver a proteger las rutas del sistema. |

---

## 3. Detalle de fallas detectadas en la ejecución

### 3.1 Suite Frontend (14 pasadas / 7 fallidas / 7 pendientes)
### 3.1 Suite Frontend (19 pasadas / 12 fallidas / 0 pendientes)

Los 7 fallos son pruebas reales nuevas (no `it.todo`) que documentan por qué `FRN-05`, `FRN-06` y `FRN-06B` de [documentacion/actual.md](documentacion/actual.md) **no están completas**, aunque existan pantallas visuales:
Las 12 fallas corresponden a pruebas reales sobre el código fuente:

1. **FRN-05 — sin guard de rol:** un usuario `OPERATOR` puede renderizar el mismo panel `Gestión de usuarios` que un `ADMIN`; [sessionGuard.ts](frontend/centinela/src/components/features/auth/routes/sessionGuard.ts) solo exige sesión activa, no rol.
2. **FRN-05 — sin llamada a la API:** [Users.jsx](frontend/centinela/src/pages/Users.jsx) no importa `fetch` ni los helpers de [api.js](frontend/centinela/src/services/api.js); nunca consulta `GET /api/users`.
3. **FRN-05 — datos fijos:** la tabla siempre muestra "Sin usuarios" y contadores en `0`, sin importar la respuesta simulada del backend.
4. **FRN-06 — sin selector de rol:** el comentario `{/* Rol de usuario */}` en [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) no tiene ningún campo debajo.
5. **FRN-06 — sin envío:** los botones "Cancelar"/"Crear usuario" son `type="submit"` sin un `<form>` que los contenga ni `onClick`/`onSubmit`; no se dispara ningún `POST /api/users`.
6. **FRN-06 — sin contraseña temporal:** no hay estado ni lógica para mostrar `contrasenaTemp` tras la creación, porque no hay creación real.
7. **FRN-06B — sin acción de edición:** no existe ningún botón "Editar" porque la tabla nunca renderiza filas de usuarios.
1. **`navigation.test.tsx` (2 fallos — `FRN-03` / `FIX-08`):**
   - Un usuario sin sesión que navega a `/dashboard` no es redirigido a `/login` porque la ruta se declaró en `MainLayoutAuth` sin guard.
   - El enlace de navegación a `/instances` no cumple con el flujo protegido requerido.
2. **`admin-users.test.tsx` (10 fallos — `FRN-05`, `FRN-06`, `FRN-06B`, `FRN-07`, `FRN-08`):**
   - **`FRN-05` (3 fallos):** `/users` no bloquea a `OPERATOR` (sin guard de rol), la página no llama a `GET /api/users` con Bearer, y la tabla ignora datos reales mostrando mocks en memoria.
   - **`FRN-06` (3 fallos):** El formulario de alta no dispara `POST /api/users`, no muestra la contraseña temporal tras la creación, y no existe llamada `DELETE /api/users/:id` al presionar eliminar.
   - **`FRN-06B` (1 fallo):** Al guardar cambios en `UserDetail.tsx`, no se dispara ninguna petición `PUT /api/users/:id`.
   - **`FRN-07` (2 fallos):** `UserDetail.tsx` no consulta `GET /api/instances` ni envía el array de VMIDs seleccionados a la ruta de permisos.
   - **`FRN-08` (1 fallo):** No hay manejo controlado de `403` al interactuar con recursos protegidos.

Los 7 criterios restantes de `FRN-07`/`FRN-08` siguen en `it.todo`: a diferencia de los anteriores, no existe ningún componente, ruta ni servicio contra el cual afirmar o refutar el criterio (no hay selector de instancias ni interceptor de errores en el código actual).
### 3.2 Suite Backend (3 pasadas / 10 fallidas)

### 3.2 Suite Backend (4 fallas/escenarios relevantes nuevos: BAC-05 x2, BAC-06/06B y BAC-10/11)
1. **`BAC-05` — Falta de restricción de unicidad en base de datos:**
   - La consulta al catálogo de PostgreSQL confirma que no existe índice único sobre `(usuario_id, vmid_proxmox)`.
2. **`BAC-03` — Causa raíz de la regresión de sesiones (`FIX-01`):**
   - `SesionActiva.UsuarioID` tiene `uniqueIndex:idx_usuario_vmid`, impidiendo logins sucesivos (`duplicate key value violates unique constraint "idx_usuario_vmid"`).
3. **Cascada por falta de sesión válida (5 fallos):**
   - `LOGIN-01`, `BAC-10/11`, `BAC-09`, `BAC-05/06/06B` y `LOGIN-03` fallan con `401 MISSING_TOKEN` o `401 AUTH_FAILED` al intentar autenticarse con el token de sesión.
4. **Hito de recursos (`resource_access_acceptance_test.go` — 3 fallos):**
   - `BAC-07`: `PUT /admin/users/:id/permissions` no responde en la ruta documentada.
   - `BAC-08`: `GET /instances/9999` devuelve `404 page not found` (falta endpoint de recursos y middleware).
   - `BAC-14`: `GET /instances` devuelve `404 page not found` (falta inventario normalizado de Proxmox).

La suite se actualizó para enviar `{"email": "...", "password": "..."}`. Con ese contrato, login, QR, TOTP, roles y protección de rutas pasan. La revisión en detalle contra `documentacion/actual.md` encontró dos problemas de naturaleza distinta:
---

1. **BAC-05 — bug funcional confirmado con evidencia, no un tema de nomenclatura:** usar un `id` propio como PK en `permisos_instancia` (definida en [backend/internal/core/domain/models.go](backend/internal/core/domain/models.go)) en vez de una clave primaria compuesta `(usuario_id, vmid_proxmox)` **es válido por sí solo**; muchos diseños prefieren una PK subrogada. El problema real es que **no existe ninguna restricción de unicidad** sobre esas dos columnas juntas (ni PK compuesta ni `UNIQUE`): `usuario_id` y `vmid_proxmox` solo tienen índices simples (`gorm:"index"`). Esto se comprobó de dos formas:
   - **Esquema:** una consulta a `pg_index`/`pg_attribute` contra el PostgreSQL real de la suite confirma que no hay ningún índice único sobre ese par de columnas.
   - **Comportamiento real:** se envió `PUT /users/:id/instances` con `{"vmids": [201, 201]}` y se contó cuántas filas quedó en la tabla para ese usuario e instancia. Resultado: **2 filas**, no 1. La lógica de "reemplazar permisos" (borrar e insertar) evita duplicados solo entre llamadas sucesivas, pero no dentro de una misma llamada con VMIDs repetidos.
   - Este bug también afecta a `BAC-07`, porque usa el mismo repositorio (`ReemplazarPermisos`).
2. **BAC-06 / BAC-06B — el grupo de rutas `admin` no se montó bajo `/admin`:** `actual.md` pide `GET/POST/DELETE /api/admin/users` y `PUT /api/admin/users/{id}`. En [backend/cmd/api/main.go](backend/cmd/api/main.go), la línea `admin := api.Group("/", middleware.RequireAuth(), middleware.RequireRole("ADMIN"))` nombra la variable `admin` (reconociendo el concepto), pero le pasa `"/"` como prefijo en vez de `"/admin"`; por eso el CRUD termina expuesto en `/api/users` y `GET /api/admin/users` responde `404 page not found`. El resto de la API sí respeta los prefijos documentados (`/api/auth/*`, `/api/roles`), lo que sugiere un detalle pendiente de terminar, no un cambio de diseño.

Los 3 hallazgos quedan como pruebas reales que fallan (no como notas en markdown). El resto del comportamiento de `BAC-06`/`BAC-06B` (contraseña temporal, `cambio_contrasena`, edición de nombre/correo/rol/estado, rechazo `403` a `OPERATOR` en las cinco operaciones CRUD) sí está implementado y pasa.

En contraste, el frontend (`FRN-05`, `FRN-06`, `FRN-06B`) no tiene ninguna de estas funcionalidades implementadas: no es solo un desvío de contrato, es la ausencia total de la lógica (ver 3.1).

### 3.4 Nuevas tareas 2FA

1. **BAC-10:** la ruta implementada es `GET /api/auth/2fa/qr` (no `GET /api/auth/2fa/setup`), pero la lógica de generación inicial sí existe. El fallo funcional no es el nombre de la ruta: después de completar la vinculación, la misma sesión de login puede solicitar otro QR y reemplazar el secreto cifrado. Debe rechazarse salvo reset/relink autorizado.
2. **BAC-11:** el secreto se cifra con AES-GCM usando `TOTP_ENCRYPTION_KEY`, se persiste cifrado y el login posterior lo usa para validar TOTP. La prueba nueva deja preparado el caso de código reutilizado en una sesión nueva; debe completarse una vez corregido BAC-10 para determinar si existe replay.
3. **FRN-09:** `LoginContinuation` ya consume `GET /api/auth/2fa/qr`, renderiza `qrBase64`, muestra `secretoManual` con formato visual y habilita el formulario de seis dígitos. [test/front/two-factor-enrollment.test.tsx](test/front/two-factor-enrollment.test.tsx) pasa de forma enfocada.

### 3.3 Estado del hito "Gestión Administrativa de Usuarios"

El hito exige, textualmente: administrador entra a `/admin/users`, ve la tabla real de `GET /api/admin/users`, crea/edita/desactiva usuarios, y un `OPERATOR` recibe `403` en "estos endpoints o la vista". Con la evidencia de 3.1 y 3.2, el hito **no está cumplido de punta a punta**:

- Backend: el CRUD funciona, pero el grupo de rutas quedó montado en `/api/users` en vez de `/api/admin/users` por el detalle de código señalado en 3.2 (falla real).
- Frontend: no existe `/admin/users`; existe `/users`, sin guard de rol, sin conexión a la API y sin alta/edición funcionales (falla real, ver 3.1).
- La condición "un `OPERATOR` recibe 403 en la vista" solo se cumple del lado backend; del lado frontend un `OPERATOR` puede navegar `/users` sin restricción.

## 4. Validaciones ejecutadas

| Suite | Comando | Resultado |
|---|---|---|
| `test/front` | `corepack pnpm --dir test/front test` | 14 aprobadas, **7 fallidas**, 7 `todo` |
| `test/back` | `go test -v -count=1 ./...` (con Docker Compose) | **3 fallidas** (`BAC-05` x2, `BAC-06`/`BAC-06B`), resto de pruebas implementadas aprobadas, `BAC-08`/`BAC-14` en `SKIP` explícito |
| `test/front` | `pnpm --dir test/front test` | **19 aprobadas**, **12 fallidas**, **0 pendientes** (31 pruebas) |
| `test/back` | `/snap/bin/go test -v -count=1 ./...` (con Docker Compose y Postgres 16) | **3 aprobadas**, **10 fallidas** (13 pruebas) |

Las 10 fallas totales (7 frontend + 3 backend) son intencionales: documentan con evidencia ejecutable qué falta para cerrar `FRN-05`, `FRN-06`, `FRN-06B`, `BAC-05`, `BAC-06` y `BAC-06B`, corrigiendo el informe previo que las daba por completas o "preparadas" sin una prueba real que lo respaldara.
---

## 5. Próximos pasos recomendados

1. Corregir BAC-10: rechazar `GET /api/auth/2fa/qr` para usuarios con `totp_vinculado=true` y permitir regeneración solo mediante reset/relink autorizado; luego completar la verificación de replay de BAC-11.
2. Implementar el guard de rol para `/users` y `/users/new` (solo `ADMIN`), conectar [Users.jsx](frontend/centinela/src/pages/Users.jsx) a `GET /api/users` y agregar el selector de rol y el `onSubmit` real en [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) para que las 7 pruebas de FRN-05/FRN-06/FRN-06B pasen.
3. Decidir si el contrato definitivo conserva `/api/users` o corrige `main.go` para montar el grupo `admin` bajo `/api/admin` (afecta a `BAC-06`, `BAC-06B` y `BAC-07`, que hoy comparten el mismo grupo de rutas); actualizar backend, frontend y tests en el mismo cambio para no dejar rutas duplicadas.
4. Agregar una restricción de unicidad sobre `(usuario_id, vmid_proxmox)` en la struct `PermisoInstancia` de [backend/internal/core/domain/models.go](backend/internal/core/domain/models.go) — como PK compuesta o como `gorm:"uniqueIndex:..."` sobre ambas columnas; cualquiera de las dos formas resuelve el bug de duplicados demostrado en 3.2.
5. Definir el contrato final del hito de recursos: rutas de permisos, forma del inventario normalizado, operación protegida sobre instancias y mecanismo de mock/spy para garantizar que un `403` no llegue a Proxmox.
1. **Corregir `FIX-01` (Backend - Prioridad Crítica):**
   - En [models.go](backend/internal/core/domain/models.go), quitar el tag `uniqueIndex:idx_usuario_vmid` de `SesionActiva.UsuarioID`, `Auditoria` y `Notificacion`.
   - Asignar un nombre específico a la restricción única en `PermisoInstancia`: `gorm:"uniqueIndex:idx_permisos_usuario_vmid"`.
   - Esto destraba inmediatamente el login repetido, la verificación TOTP y toda la cascada de tests backend.
2. **Corregir `FIX-08` (Frontend - Prioridad Crítica):**
   - En [applicationRoutes.tsx](frontend/centinela/src/routes/applicationRoutes.tsx), mover `/dashboard`, `/instances`, `/users`, `/users/new`, `/users/:userId` y `/auditoria` de vuelta al grupo protegido con `loader: loadProtectedSession` y `ProtectedLayout`.
3. **Conectar el panel de administración a la API (`FIX-05`, `FIX-06`, `FIX-06B`):**
   - Reemplazar el array estático `usersList` en [Users.tsx](frontend/centinela/src/pages/Users.tsx) por una llamada a `GET /api/admin/users`.
   - Conectar el formulario de [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) para disparar `POST /api/admin/users` y mostrar el modal con la contraseña temporal.
   - Conectar el formulario de [UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx) para llamar a `PUT /api/admin/users/:id` y `DELETE /api/admin/users/:id`.
4. **Implementar endpoints del Hito de Recursos (Backend):**
   - Definir contrato final para permisos (`/api/admin/users/:id/permissions`).
   - Implementar `GET /api/instances` normalizando inventario de Proxmox (`BAC-14`) y el middleware de validación por recurso (`BAC-08`).
