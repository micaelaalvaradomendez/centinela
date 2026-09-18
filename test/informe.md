# Informe de estado actual — login, 2FA e integración

Fecha: 18/09/2026.
Fuentes: [documentacion/actual.md](documentacion/actual.md), [test/back](test/back) (suite Go, ejecución real contra PostgreSQL) y [test/front](test/front) (suite Vitest sobre el código real del frontend).

Commits auditados (coinciden con `origin/main` de cada submódulo):

| Componente | Commit |
|---|---|
| Frontend | `7775b62` |
| Backend | `caf71c2` |



## 1. Resumen ejecutivo

- **Actualización de submódulos:** ambos submódulos están al día con sus ramas `main` (`backend@caf71c2`, `frontend@7775b62`); el último pull no modificó ningún archivo de `BAC-05` a `BAC-14` ni de `FRN-05` a `FRN-08`.
- **Frontend administrativo (`FRN-05`, `FRN-06`, `FRN-06B`): sin implementar.** Existen pantallas visuales estáticas ([Users.jsx](frontend/centinela/src/pages/Users.jsx), [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx)), pero sin guard de rol, sin llamadas a la API, sin selector de rol ni envío de formulario. 7 pruebas reales nuevas en [test/front/admin-users.test.tsx](test/front/admin-users.test.tsx) lo confirman: **7 de 7 fallan**. Detalle en 3.1.
- **Backend administrativo (`BAC-05`, `BAC-06`, `BAC-06B`, `BAC-09`): lógica funcional correcta, con un bug real de duplicados y un desvío de ruta.** El CRUD de usuarios, la contraseña temporal y el RBAC funcionan correctamente en las cinco operaciones (crear, listar, ver detalle, editar, dar de baja). Dos hallazgos distintos:
  1. **Bug funcional (no de nomenclatura):** usar un `id` propio como PK en `permisos_instancia` en vez de una compuesta es válido por sí solo. El problema real es que no existe **ninguna** restricción de unicidad sobre `(usuario_id, vmid_proxmox)`. Se demostró en vivo: enviar `{"vmids": [201, 201]}` persiste **2 filas duplicadas** en vez de 1. Afecta a `BAC-05` y, por compartir repositorio, también a `BAC-07`.
  2. **Desvío de ruta (`BAC-06`/`BAC-06B`):** en [backend/cmd/api/main.go](backend/cmd/api/main.go) la variable del grupo de rutas se llama `admin`, pero se monta con `Group("/")` en vez de `Group("/admin")`, por lo que `GET /api/admin/users` no existe (solo `/api/users`). El resto de la API (`/api/auth/*`, `/api/roles`) sí coincide con `actual.md`, lo que sugiere un detalle de implementación pendiente, no un cambio de diseño deliberado.

  3 pruebas reales nuevas en [test/back/backend_acceptance_test.go](test/back/backend_acceptance_test.go) confirman ambos hallazgos. Detalle en 3.2.
- **Hito "Gestión Administrativa de Usuarios": no cumplido de punta a punta.** Ver 3.3.
- **Hito "Control de Acceso Basado en Recursos": sin empezar (`BAC-08`, `BAC-14`, `FRN-07`, `FRN-08`).** `BAC-07` tiene una implementación parcial (`PUT /api/users/:id/instances`) que funciona pero con el mismo desvío de prefijo de ruta.
- **Nuevas tareas 2FA (`BAC-10`, `BAC-11`, `FRN-09`):** hay implementación real, no solo preparación. El backend genera QR, cifra el secreto y permite login posterior; el frontend `LoginContinuation` muestra QR, clave manual y OTP. Las pruebas nuevas detectan un fallo de BAC-10: un usuario ya vinculado puede regenerar su secreto llamando nuevamente a `GET /api/auth/2fa/qr` sin pasar por reset/relink. FRN-09 tiene una prueba enfocada aprobada.
- **Historial de correcciones de este informe:** el 18/09/2026 se detectó que los criterios `FRN-05`/`FRN-06`/`FRN-06B` estaban declarados como `it.todo` (Vitest nunca los ejecuta) y que los tests backend de `BAC-05`/`BAC-06`/`BAC-06B` no verificaban dos requisitos textuales de `actual.md`. Ambos huecos se corrigieron con pruebas reales que hoy fallan, y ese es el estado reflejado en este informe.
- **Nota sobre la base de datos:** el repositorio no tiene ningún archivo `.sql` con `CREATE TABLE` (`backend/scripts/init.sql` solo define la función `uuid_generate_v7()`). El esquema real lo genera `gorm.AutoMigrate` en tiempo de ejecución a partir de las structs de [backend/internal/core/domain/models.go](backend/internal/core/domain/models.go), y se verifica consultando en vivo el PostgreSQL 16 que levanta `test/back` con Docker Compose (no es una base simulada).

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

## 3. Detalle de fallas detectadas en la ejecución

### 3.1 Suite Frontend (14 pasadas / 7 fallidas / 7 pendientes)

Los 7 fallos son pruebas reales nuevas (no `it.todo`) que documentan por qué `FRN-05`, `FRN-06` y `FRN-06B` de [documentacion/actual.md](documentacion/actual.md) **no están completas**, aunque existan pantallas visuales:

1. **FRN-05 — sin guard de rol:** un usuario `OPERATOR` puede renderizar el mismo panel `Gestión de usuarios` que un `ADMIN`; [sessionGuard.ts](frontend/centinela/src/components/features/auth/routes/sessionGuard.ts) solo exige sesión activa, no rol.
2. **FRN-05 — sin llamada a la API:** [Users.jsx](frontend/centinela/src/pages/Users.jsx) no importa `fetch` ni los helpers de [api.js](frontend/centinela/src/services/api.js); nunca consulta `GET /api/users`.
3. **FRN-05 — datos fijos:** la tabla siempre muestra "Sin usuarios" y contadores en `0`, sin importar la respuesta simulada del backend.
4. **FRN-06 — sin selector de rol:** el comentario `{/* Rol de usuario */}` en [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) no tiene ningún campo debajo.
5. **FRN-06 — sin envío:** los botones "Cancelar"/"Crear usuario" son `type="submit"` sin un `<form>` que los contenga ni `onClick`/`onSubmit`; no se dispara ningún `POST /api/users`.
6. **FRN-06 — sin contraseña temporal:** no hay estado ni lógica para mostrar `contrasenaTemp` tras la creación, porque no hay creación real.
7. **FRN-06B — sin acción de edición:** no existe ningún botón "Editar" porque la tabla nunca renderiza filas de usuarios.

Los 7 criterios restantes de `FRN-07`/`FRN-08` siguen en `it.todo`: a diferencia de los anteriores, no existe ningún componente, ruta ni servicio contra el cual afirmar o refutar el criterio (no hay selector de instancias ni interceptor de errores en el código actual).

### 3.2 Suite Backend (4 fallas/escenarios relevantes nuevos: BAC-05 x2, BAC-06/06B y BAC-10/11)

La suite se actualizó para enviar `{"email": "...", "password": "..."}`. Con ese contrato, login, QR, TOTP, roles y protección de rutas pasan. La revisión en detalle contra `documentacion/actual.md` encontró dos problemas de naturaleza distinta:

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

Las 10 fallas totales (7 frontend + 3 backend) son intencionales: documentan con evidencia ejecutable qué falta para cerrar `FRN-05`, `FRN-06`, `FRN-06B`, `BAC-05`, `BAC-06` y `BAC-06B`, corrigiendo el informe previo que las daba por completas o "preparadas" sin una prueba real que lo respaldara.

## 5. Próximos pasos recomendados

1. Corregir BAC-10: rechazar `GET /api/auth/2fa/qr` para usuarios con `totp_vinculado=true` y permitir regeneración solo mediante reset/relink autorizado; luego completar la verificación de replay de BAC-11.
2. Implementar el guard de rol para `/users` y `/users/new` (solo `ADMIN`), conectar [Users.jsx](frontend/centinela/src/pages/Users.jsx) a `GET /api/users` y agregar el selector de rol y el `onSubmit` real en [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) para que las 7 pruebas de FRN-05/FRN-06/FRN-06B pasen.
3. Decidir si el contrato definitivo conserva `/api/users` o corrige `main.go` para montar el grupo `admin` bajo `/api/admin` (afecta a `BAC-06`, `BAC-06B` y `BAC-07`, que hoy comparten el mismo grupo de rutas); actualizar backend, frontend y tests en el mismo cambio para no dejar rutas duplicadas.
4. Agregar una restricción de unicidad sobre `(usuario_id, vmid_proxmox)` en la struct `PermisoInstancia` de [backend/internal/core/domain/models.go](backend/internal/core/domain/models.go) — como PK compuesta o como `gorm:"uniqueIndex:..."` sobre ambas columnas; cualquiera de las dos formas resuelve el bug de duplicados demostrado en 3.2.
5. Definir el contrato final del hito de recursos: rutas de permisos, forma del inventario normalizado, operación protegida sobre instancias y mecanismo de mock/spy para garantizar que un `403` no llegue a Proxmox.
