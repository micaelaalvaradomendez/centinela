
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Verificación del 23/09/2026** ([test/informe.md](../test/informe.md)): las tareas completas (`BAC-14`, `FIX-16`, `BAC-17`, `FRN-13`, `SEC-02`, `FIX-14`) y las implementadas con problemas (`FRN-10`, `FRN-14`) pasaron a [`terminado.md`](terminado.md). Los problemas se registraron como `FIX-21` a `FIX-25` en [`futuro.md`](futuro.md). En este archivo quedan solo las tareas **no implementadas**.

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (23/09/2026):** el backend cumple el hito completo. En el frontend, la tabla real y el guard para `OPERATOR` funcionan, pero el hito sigue abierto por dos problemas:
> - Tras el alta no se muestra ninguna confirmación (`FIX-24`).
> - No se puede editar el perfil ni desactivar al usuario desde la interfaz (`FIX-24`, `FIX-25`).

---

### `SEC-01` - Refresh token en cookie HttpOnly en backend

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** A definir. **Prioridad alta.**
- **Depende de:** `BAC-17`.
- **Estado actual:** la revocación de sesiones y la renovación de tokens están implementadas, pero el backend todavía recibe el refresh token en JSON y lo devuelve en el body. No existe emisión ni lectura de cookie HttpOnly.
- **Estado verificado (23/09/2026):** no implementada. `POST /auth/2fa/verify` no emite `Set-Cookie`, y `/auth/refresh` y `/auth/logout` exigen `refreshToken` en el body (`auth_handler.go:69-71`, `:190-192`).
  - **Bloquea la integración:** el frontend (`SEC-02`, ya terminada) envía `{}` en refresh y logout y el backend responde `400 INVALID_REQUEST`.
  - Con ambos `main` desplegados, el logout de la interfaz no revoca la sesión en el servidor y la renovación silenciosa falla.
  - Pruebas: `test/back/session_security_acceptance_test.go`, casos `SEC-01…` y `SEC-01 SEC-02 integracion…`.
- **Entregable:**
	1. Emitir el `refreshToken` mediante `Set-Cookie` al completar `POST /api/auth/2fa/verify`.
	2. Leerlo desde la cookie en `POST /api/auth/refresh` y `POST /api/auth/logout`.
	3. Dejar de devolver el refresh token en las respuestas JSON públicas.
	4. Configurar `HttpOnly`, `Secure`, `SameSite` y `Path` de forma segura y configurable para desarrollo y producción.
	5. Eliminar la cookie al cerrar sesión o cuando la sesión sea inválida.
- **Criterio de éxito:** el backend renueva y revoca sesiones usando exclusivamente la cookie HttpOnly; el refresh token no aparece en cuerpos JSON ni en logs; una cookie inválida o vencida produce un error controlado.

### `SEC-04` - Esquema extensible de niveles de acceso a recursos en Backend y Modelo de Datos

- **Área:** Backend
- **Asignado:** Lisandro / Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 2).
- **Depende de:** `BAC-05`, `BAC-07` y `BAC-08`.
- **Estado verificado (23/09/2026):** no implementada.
  - `domain.PermisoInstancia` no tiene `NivelAcceso` (`models.go:79-83`).
  - `VerificarAcceso` y `RequireInstanceAccess` no reciben el nivel requerido.
  - El frontend ya ofrece "Solo lectura" en `rolesAndPermissions.tsx`, pero el backend lo ignora.
  - Prueba: `test/back/resource_access_acceptance_test.go`, caso `SEC-04…`.
- **Problema:** La tabla `permisos_instancia` solo almacena `(usuario_id, vmid_proxmox)` como una relación binaria (asignado o no asignado). Sin embargo, en el documento de diseño `Diseño de endpoints para front.md` y en las siguientes etapas (Etapa 1 Ciclo de vida y Etapa 3 Snapshots) se proyectan niveles de acceso sobre las instancias (por ejemplo `FULL_ACCESS` para operar energía/snapshots vs `READ_ONLY` para monitoreo de métricas sin emitir órdenes destructivas). Si no se sienta esta base en el esquema y en el guard ahora, agregar niveles de acceso en la Etapa 1 requerirá migrar tablas en producción y alterar contratos.
- **Entregable:**
  1. Agregar en `domain.PermisoInstancia` la columna `nivel_acceso` (VARCHAR(30) default `'FULL_ACCESS'`) con restricción CHECK o enum para los valores `FULL_ACCESS` y `READ_ONLY`.
  2. Actualizar el puerto `InstanceRepository.VerificarAcceso` para aceptar opcionalmente el nivel de acceso requerido (`requiredLevel: string`).
  3. Extender el middleware `RequireInstanceAccess(repo, paramName, requiredLevel)` para permitir guards como `RequireInstanceAccess(repo, "vmid", "FULL_ACCESS")` en endpoints mutantes (`POST /instances/:vmid/start`, `POST /instances/:vmid/stop`) y tolerar `READ_ONLY` en consultas (`GET /instances/:vmid`).
  4. Mantener retrocompatibilidad total: si el payload de `PUT /permissions` solo envía `vmids: [101]`, asignar `FULL_ACCESS` por defecto.
- **Criterio de éxito:** La migración crea el campo sin romper registros previos; el middleware `RequireInstanceAccess` verifica tanto la pertenencia de la instancia como el nivel de permiso; si un usuario tiene permiso `READ_ONLY` sobre la VM 101, puede consultar su estado pero recibe 403 al intentar ejecutar una acción de apagado/encendido.

### `FRN-11` - Acciones administrativas de recuperación

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2 h
- **Ventana propuesta:** 18/09/2026, 12:00-14:00
- **Depende de:** `FRN-05`, `BAC-13` y `BAC-15`.
- **Estado verificado (23/09/2026):** no implementada.
  - El botón "Restablecer contraseña" de `Users.tsx:354` no tiene `onClick`.
  - No existe la acción de reset de 2FA ni el diálogo de confirmación.
  - Prueba: `test/front/admin-recovery.test.tsx` (1/4). El detalle de implementación ya está en `FIX-19` de [`futuro.md`](futuro.md).
- **Entregable:** acciones separadas para restablecer contraseña y 2FA desde el panel de usuarios, ambas con confirmación explícita, estado de carga y notificación del resultado.
- **Criterio de éxito:** un administrador puede iniciar cada recuperación sin confundir sus efectos; la tabla refleja que el 2FA quedó desvinculado y nunca muestra secretos ni hashes.

### `FRN-12` - Vistas de recuperación de contraseña (RF-13)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 3 h
- **Ventana propuesta:** A definir (posterior a `BAC-20`).
- **Depende de:** `BAC-19`, `BAC-20` y el maquetado existente de `RecoverPassword.tsx`.
- **Estado verificado (23/09/2026):** no implementada. `RecoverPassword.tsx` solo cambia de paso y no llama a `/auth/password/forgot` ni a `/auth/password/reset`. Prueba: `test/front/recover-password.test.tsx` (1/4). El detalle de implementación ya está en `FIX-18` de [`futuro.md`](futuro.md).
- **Entregable:** conectar `RecoverPassword.tsx` al flujo real: paso de ingreso de correo, paso de ingreso del código de seis dígitos y paso de nueva contraseña, con manejo de errores del backend en cada paso.
- **Criterio de éxito:** una cuenta puede recuperar el acceso sin intervención de un administrador, y los errores de código inválido o vencido se muestran junto al campo correspondiente.

### `SEC-03` - Contexto y sistema reactivo de permisos en Frontend (UI/UX RBAC + Resources)

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 2,5 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 2).
- **Depende de:** `FRN-04`, `BAC-07` y `BAC-09`.
- **Estado verificado (23/09/2026):** no implementada.
  - `src/context/AuthContext.js` está vacío.
  - No existen `usePermissions` ni `PermissionGate`.
  - El menú solo filtra con un `isAdmin` ad hoc en `Sidebar.tsx` y no tiene un acceso a Auditoría para el administrador.
  - Prueba: `test/front/navigation.test.tsx` (1/4).
- **Problema:** En el frontend actual, el control de roles y permisos está disperso y acoplado únicamente a los loaders de rutas (`loadAdminSession`). No existe un mecanismo reactivo a nivel de componentes (`usePermissions` / `PermissionGate`) para ocultar o deshabilitar condicionalmente acciones según el rol (`ADMIN` vs `OPERATOR`) o según las instancias asignadas al operador. Esto genera que en vistas como `UserDetail.tsx` o `Header.tsx` se muestren controles estáticos o se dependa de que el backend rechace con 403, en lugar de ofrecer una experiencia fluida y consistente.
- **Entregable:**
  1. Crear un hook y contexto `usePermissions()` / `useAuthUser()` en `frontend/centinela/src/context/` que exponga helpers como `isAdmin`, `isOperator`, `canAccessInstance(vmid: number)` y `hasRole(role: string)`.
  2. Crear un componente wrapper `<PermissionGate requiredRole="ADMIN" fallback={...}>` para condicionar la renderización de botones y secciones administrativas (ej: botón "Eliminar usuario", accesos a auditoría, etc.).
  3. Integrar la reactividad en el menú de navegación y en el header para que las opciones no permitidas a un operador no aparezcan en la interfaz.
- **Criterio de éxito:** Si un operador inicia sesión, la interfaz no muestra accesos directos ni botones exclusivos de administrador; si intenta interactuar con un componente restringido, el helper `canAccessInstance` evalúa en memoria las instancias permitidas cargadas en la sesión; las pruebas unitarias de componentes verifican el render condicional.



---


### `FIX-17` - Alinear y formalizar endpoints y contratos de contraseñas (Backend)

- **Área:** Backend
- **Asignado:** Lisandro / Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 3).
- **Depende de:** `BAC-12`, `BAC-15`, `BAC-19`, `BAC-20`.
- **Problema y evidencia:**
  1. *Rutas ambiguas y divergencia con la documentación:* `actual.md` y versiones previas de planificación referenciaban `POST /api/auth/change-password` para el cambio obligatorio, mientras el backend implementó `PUT /api/account/password`.
  2. *Inconsistencia en reset administrativo:* En `terminado.md` (`BAC-15`) se documentó `POST /api/admin/users/{id}/reset-password`, pero en `cmd/api/main.go` se montó `POST /api/admin/users/:id/password/reset`.
  3. *Confusión entre reset público y reset administrativo:* Ambos comparten el sufijo `/password/reset` pero con semánticas, autorizaciones y payloads totalmente diferentes.
- **Entregable:**
  1. Formalizar como canónico `PUT /api/account/password` para el cambio autenticado de contraseña (payload: `{ contrasenaActual, contrasenaNueva }`), manteniendo bajo `RequireAuth` el bloqueo 403 `PASSWORD_CHANGE_REQUIRED` a otras rutas.
  2. Mantener `POST /api/auth/password/forgot` (payload: `{ email }`) y `POST /api/auth/password/reset` (payload: `{ email, codigo, nuevaContrasena }`) para la recuperación pública (RF-13), garantizando que tras el reset exitoso la cuenta quede con `cambio_contrasena = false` y todas las sesiones previas revocadas en `sesiones_activas`.
  3. Confirmar como canónico `POST /api/admin/users/:id/password/reset` (sin body, requiere rol `ADMIN`), el cual genera una contraseña temporal, activa `cambio_contrasena = true`, revoca sesiones y despacha la clave por `EmailService` (`MockEmailService`).
  4. Actualizar Swagger/OpenAPI y eliminar alias no utilizados o rutas contradictorias.
  5. Asegurar consistencia de códigos de error estructurados: `PASSWORD_CHANGE_REQUIRED` (403), `PASSWORD_CHANGE_FAILED` (400), `RESET_FAILED` (400) y `INVALID_REQUEST` (400).
- **Criterio de éxito:** Swagger expone los contratos precisos y unificados; `PUT /api/account/password` valida la contraseña actual y libera la cuenta; `POST /api/auth/password/reset` valida el OTP de 6 dígitos sin requerir contraseña actual; las llamadas a cada endpoint responden con los códigos de estado y payloads esperados.


### `FIX-18` - Conexión de `RecoverPassword.tsx` a la API y alineación de tests (`FRN-12` / RF-13) (Frontend)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 3 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 3).
- **Depende de:** `BAC-19`, `BAC-20` y `FIX-17`.
- **Problema y evidencia en código y pruebas (`test/front/recover-password.test.tsx`):**
  1. *Formulario desconectado de la red:* `frontend/centinela/src/pages/RecoverPassword.tsx` es actualmente un cascarón visual que solo maneja estado local mediante `handleNext()`, sin invocar `apiClient` ni interactuar con el backend en ninguno de sus tres pasos.
  2. *Desalineación de rutas en tests:* `test/front/recover-password.test.tsx` espera que el componente llame a URLs ficticias (`/recover/request`, `/auth/recovery`, `/recover/confirm`, `/recovery/confirm`), las cuales nunca existieron en el backend.
  3. *Flujo de finalización ausente:* Al finalizar el tercer paso, el formulario no consume `POST /api/auth/password/reset` ni redirige al usuario al login.
- **Entregable:**
  1. Conectar el Paso 1 (Solicitud) a `POST /api/auth/password/forgot` enviando `{ email }`, mostrando spinner/estado de carga y avanzando al Paso 2 al recibir HTTP 200.
  2. En el Paso 2 (Código), vincular `InputOTP` al estado del código de 6 dígitos y validar que no esté vacío antes de avanzar al Paso 3.
  3. Conectar el Paso 3 (Nueva contraseña) a `POST /api/auth/password/reset` enviando `{ email, codigo, nuevaContrasena }` (sin solicitar contraseña actual).
  4. Procesar la respuesta HTTP 200 mostrando un toast de éxito («Contraseña restablecida con éxito») y redirigiendo a `/login` con `{ replace: true }`.
  5. Procesar los errores del backend (`RESET_FAILED`, código expirado o superación de 3 intentos) y mostrarlos junto al campo correspondiente en la interfaz.
  6. Actualizar la suite `test/front/recover-password.test.tsx` para mockear y validar los llamados reales a `/auth/password/forgot` y `/auth/password/reset`.
- **Criterio de éxito:** El usuario puede recuperar su cuenta de punta a punta consumiendo la API; la suite `test/front/recover-password.test.tsx` pasa al 100% verificando los endpoints canónicos del backend.

### `FIX-19` - Acciones de Restablecimiento Administrativo de Contraseña y 2FA en `Users.tsx` (`FRN-11`) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz / Cristian
- **Estimación:** 2,5 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 3).
- **Depende de:** `BAC-13`, `BAC-15`, `FRN-05` y `test/front/admin-recovery.test.tsx`.
- **Problema y evidencia en código:**
  1. *Botón sin acción:* En `frontend/centinela/src/pages/Users.tsx` (línea 354), el botón «Restablecer contraseña» carece de manejador `onClick`.
  2. *Acción de 2FA faltante:* En el menú de acciones no existe la opción para «Restablecer 2FA» / «Desvincular 2FA», requerida por `FRN-11` y verificada en `test/front/admin-recovery.test.tsx`.
  3. *Falta de reactividad y feedback:* No se ofrece modal de confirmación ni notificación al administrador sobre el envío de la clave temporal, y la tabla no actualiza el estado de 2FA tras un reset.
- **Entregable:**
  1. Conectar el botón «Restablecer contraseña» a una función con modal de confirmación explícito que invoque `POST /api/admin/users/:id/password/reset` con token Bearer, cerrando el dropdown y emitiendo un toast que informe el envío por correo.
  2. Agregar en el dropdown de acciones la opción «Restablecer 2FA» (con icono de seguridad y advertencia de acción crítica) con modal de confirmación, consumiendo `POST /api/admin/users/:id/2fa/reset`.
  3. Al completar con éxito el restablecimiento de 2FA, mutar reactivamente el estado local (`totpVinculado: false`) para que el badge de la tabla pase inmediatamente de «Activado» a «Desactivado» sin necesidad de recargar la página.
  4. Garantizar que bajo ninguna circunstancia se muestren secretos ni hashes en la interfaz (política Zero-Trust).
- **Criterio de éxito:** Un administrador puede resetear la contraseña y el 2FA de cualquier operador desde el panel; la tabla actualiza el estado de 2FA reactivamente; la suite `test/front/admin-recovery.test.tsx` pasa al 100%.



### `FIX-21` - Opción de cerrar sesión en el cambio obligatorio de contraseña (`FRN-10`) (Frontend)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** `FRN-10` y `FRN-13`.
- **Problema y evidencia:** el criterio de `FRN-10` dice que "una cuenta con contraseña temporal solo puede cerrar sesión o cambiarla". `frontend/centinela/src/pages/ChangePassword.tsx` solo ofrece el formulario de cambio; no hay forma de salir sin cambiar la clave. La prueba `test/front/password-change.test.tsx`, caso *"la pantalla de cambio ofrece cerrar sesión como única alternativa"*, falla con `Unable to find role="button" and name /cerrar sesión/i`.
- **Entregable:**
  1. Agregar en `ChangePassword.tsx` un botón secundario "Cerrar sesión" que use el hook existente `useLogout()` (`components/features/auth/hooks/useAuth.ts`). Ese hook ya llama a `POST /api/auth/logout`, limpia la sesión y navega a `/login` con `replace`.
  2. Deshabilitar el botón mientras se envía el formulario o el logout (`isSubmitting` / `isLoggingOut`).
- **Criterio de éxito:** desde `/change-password` el usuario puede cerrar sesión: se llama a `/auth/logout`, se borra `centinela_access` y se navega a `/login`. `password-change.test.tsx` pasa 7/7.

### `FIX-22` - Guard administrativo en la ruta `/auditoria` (`FRN-14` / `BAC-18`) (Frontend)

- **Área:** Frontend
- **Asignada:** Belinda / Luz
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir (alta prioridad: es un control de acceso).
- **Depende de:** `FRN-14` y `FIX-12`.
- **Problema y evidencia:** en `frontend/centinela/src/routes/applicationRoutes.tsx` la ruta `{ path: '/auditoria', Component: AuditoriaPage }` está dentro del grupo `loadProtectedSession` y **no** dentro del grupo `loadAdminSession`. Un usuario `OPERATOR` puede abrir la vista de auditoría; el backend le responde 403 a los datos, pero la pantalla administrativa se muestra igual. La prueba `test/front/audit.test.tsx`, caso *"si un operador intenta entrar a /auditoria, el guard administrativo lo redirige a /dashboard"*, falla con `expected '/auditoria' to be '/dashboard'`.
- **Entregable:** mover `{ path: '/auditoria', Component: AuditoriaPage }` al arreglo `children` del grupo con `loader: loadAdminSession`, junto a `/users`, `/users/new` y `/users/:userId`.
- **Criterio de éxito:** un `OPERATOR` que navega a `/auditoria` es redirigido a `/dashboard`, y un `ADMIN` accede normalmente. `audit.test.tsx` pasa 7/7.
- **Nota:** el acceso visible a Auditoría en el menú lateral, solo para administradores, forma parte de `SEC-03`, que sigue en `actual.md`.

### `FIX-23` - Hacer append-only la tabla `auditoria` a nivel de base de datos (`BAC-18`) (Backend)

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** `BAC-18`.
- **Problema y evidencia:** el criterio de `BAC-18` exige que "un intento de `UPDATE` o `DELETE` sobre la auditoría con las credenciales de la aplicación falle a nivel de base de datos, no solo por convención de código". Hoy no existe ningún trigger, `REVOKE` ni regla: "append-only" solo aparece en comentarios (`internal/core/ports/audit_port.go:106`, `internal/core/domain/models.go:86`). La prueba `test/back/password_recovery_acceptance_test.go`, caso *"BAC-18 auditoria … append-only en la base"*, falla con:
  ```
  la base permitió UPDATE sobre auditoria con el usuario de la aplicación; la tabla no es append-only
  la base permitió DELETE sobre auditoria con el usuario de la aplicación; la tabla no es append-only
  ```
- **Entregable:**
  1. Crear, después del `AutoMigrate` en `postgres.InitDB()` o en `scripts/init.sql`, una función y un trigger idempotentes:
     ```sql
     CREATE OR REPLACE FUNCTION auditoria_append_only() RETURNS trigger AS $$
     BEGIN
       RAISE EXCEPTION 'la tabla auditoria es append-only (% no permitido)', TG_OP;
     END; $$ LANGUAGE plpgsql;

     DROP TRIGGER IF EXISTS trg_auditoria_append_only ON auditoria;
     CREATE TRIGGER trg_auditoria_append_only
       BEFORE UPDATE OR DELETE ON auditoria
       FOR EACH ROW EXECUTE FUNCTION auditoria_append_only();
     ```
     Agregar también `BEFORE TRUNCATE ... FOR EACH STATEMENT`.
  2. El trigger funciona aunque la aplicación se conecte con el usuario dueño de la tabla, que es el caso del entorno de pruebas y del compose actual. Como defensa adicional en producción, conectar la aplicación con un rol sin privilegios `UPDATE`/`DELETE` sobre `auditoria` (`REVOKE UPDATE, DELETE, TRUNCATE ON auditoria FROM <rol_app>`).
  3. Revisar que ninguna ruta del backend haga `UPDATE`/`DELETE` sobre `auditoria`. Por ejemplo, el `ON DELETE SET NULL` de `usuario_id` al eliminar un usuario: si hace falta, resolverlo con eliminación lógica de usuarios, que es lo que ya hace `DELETE /admin/users/:id`.
- **Criterio de éxito:** `UPDATE`, `DELETE` y `TRUNCATE` sobre `auditoria` fallan en PostgreSQL con las credenciales de la aplicación; el registro y la consulta de auditoría siguen funcionando. El caso `BAC-18` de `password_recovery_acceptance_test.go` pasa.

### `FIX-24` - Confirmación del alta sin contraseña temporal y baja de usuario (`FRN-06` / `BAC-16`) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2 h
- **Ventana propuesta:** A definir.
- **Depende de:** `FRN-06`, `BAC-06` y `BAC-16`.
- **Problema y evidencia:**
  1. *Alta sin respuesta visible:* desde `BAC-16`, `POST /api/admin/users` responde `201 { id, rol, activo }` **sin `contrasenaTemp`**, porque la clave se envía por correo. `frontend/centinela/src/pages/CrearUsuarios.tsx:56` hace `setTempPassword(created.contrasenaTemp ?? null)`, y solo muestra algo si llega ese campo. Resultado: tras un alta exitosa el administrador no ve ninguna confirmación. Además, la vista conserva campos de "Contraseña temporal" y "Confirmar contraseña temporal" que ya no tienen función. Prueba: `test/front/admin-users.test.tsx`, caso *"tras el alta informa el resultado y que la clave temporal se envió por correo (BAC-16)…"*.
  2. *Baja sin acción:* el botón "Eliminar usuario" de `frontend/centinela/src/pages/detailsUserPage.tsx:110` no tiene `onClick` y nunca envía `DELETE /api/admin/users/:id`. Prueba: `admin-users.test.tsx`, caso *"el botón 'Eliminar usuario' solicita DELETE /api/admin/users/:id"*.
- **Entregable:**
  1. En `CrearUsuarios.tsx`, al recibir 201:
     - mostrar un toast o mensaje con `role="status"`, por ejemplo *"Usuario creado. La contraseña temporal se envió a <correo>."*;
     - limpiar el formulario o volver a `/users`;
     - eliminar la caja "Contraseña temporal generada" y los campos de contraseña del formulario.
  2. Manejar `502 EMAIL_DELIVERY_FAILED` con un mensaje claro: *"No se pudo enviar el correo; el usuario no fue creado"*.
  3. En `detailsUserPage.tsx`, conectar "Eliminar usuario" a un diálogo de confirmación explícita que envíe `DELETE /api/admin/users/:id` y espere 204. Después, mostrar un toast y volver a `/users`, o reflejar `activo: false`, que es una baja lógica.
- **Criterio de éxito:** el administrador ve la confirmación de cada alta sin que se muestre ninguna contraseña, y puede dar de baja a un usuario desde el detalle. Los dos casos de FRN-06 de `admin-users.test.tsx` pasan.

### `FIX-25` - Persistir la edición del perfil de usuario (`FRN-06B` / `FIX-14`) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz / Cristian
- **Estimación:** 2 h
- **Ventana propuesta:** A definir.
- **Depende de:** `FRN-06B`, `BAC-06B` y `FIX-14`.
- **Problema y evidencia:** en `frontend/centinela/src/pages/detailsUserPage.tsx`, el formulario de "Información general" (`informationOfUser.tsx` + `useEditableUser.ts`) modifica solo el estado local. "Guardar cambios" (`detailsUserPage.tsx:114`) llama únicamente a `instanceAccess.saveAssignments()` y está deshabilitado si no cambió ningún permiso de instancia. En todo el frontend **no existe ninguna llamada a `PUT /api/admin/users/:id`**: editar nombre, correo, rol o estado no se guarda nunca. Esto incumple `FRN-06B` y el entregable 3 de `FIX-14` ("guardar de forma atómica junto a la edición del perfil"). Prueba: `test/front/admin-users.test.tsx`, caso *"editar el nombre y guardar envía PUT /api/admin/users/:id con los datos modificados"*, que falla con `guardar no envió PUT /admin/users/u2`.
- **Entregable:**
  1. Agregar a `components/features/users/services/userDetailsService.ts` una función `updateUserDetails(userId, values)`. Debe enviar `PUT /api/admin/users/:id` solo con los campos modificados de `{ nombreCompleto, emailUsuario, rol, activo }` y validar la respuesta con el `isUserDetailsResponse` existente.
  2. Habilitar "Guardar cambios" cuando haya cambios en el perfil **o** en los permisos. Al guardar, enviar primero el `PUT` del perfil y luego el de permisos, con un solo estado de carga, un toast de resultado y la vista refrescada con la respuesta.
  3. Mostrar los errores del backend: `USER_CONFLICT` (409) para correo duplicado y 400 para datos inválidos.
- **Criterio de éxito:** un administrador cambia el nombre, el correo, el rol o el estado de un usuario, guarda, y el cambio se refleja en el detalle y en la tabla sin recargar. El caso FRN-06B de `admin-users.test.tsx` pasa, y los 3 casos de FIX-14 siguen en verde.

- ---

### `INF-08` - Configuración de servidor/cuenta SMTP y variables de entorno para correo saliente

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `INF-03`, `INF-04`.
- **Problema y contexto:** Para habilitar el envío real de correos (`RF-09` y `RF-13`) desde el backend, la infraestructura debe proveer las credenciales del relay/servidor SMTP y habilitar la salida de red en los puertos correspondientes (`587` STARTTLS / `465` TLS).
- **Entregable:**
  1. Configurar la cuenta de servicio o relay SMTP e inyectar en el entorno del servidor y en `.env.example` / `docker-compose.yml` las variables: `EMAIL_PROVIDER=smtp`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` y `SMTP_FROM`.
  2. Verificar conectividad de red saliente desde el contenedor del backend hacia el host SMTP.
- **Criterio de éxito:** Variables documentadas y conectividad validada desde el contenedor del backend hacia el puerto SMTP sin bloqueos de firewall.

### `INF-06` - Despliegue de servicio Redis en Docker Compose y red `vmbr1`

- **Área:** Infraestructura
- **Asignados:** Nico y Lucas
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (Cierre de Fase Base, previo a `BAC-17B`).
- **Depende de:** `INF-03` e `INF-04`.
- **Problema y contexto:** Se adelanta el despliegue de Redis al cierre de la fase base para alojar las sesiones activas con expiración automática (`TTL`), los tickets efímeros de WebSockets/SSE y la caché de Proxmox.
- **Entregable:**
  1. Agregar el contenedor `redis:7-alpine` en `docker-compose.yml` y en la red interna `vmbr1`, configurado con contraseña (`REDIS_PASSWORD`), límite de memoria (`maxmemory 256mb`, política `volatile-lru`) y variables `REDIS_ADDR` / `REDIS_PASSWORD` expuestas al contenedor del backend.
- **Criterio de éxito:** Servicio Redis operativo y accesible desde la red interna `vmbr1`, respondiendo `PONG` a `redis-cli PING` con autenticación.


### `INF-05` - CORS y TLS en el borde

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (junto con `INF-04`).
- **Depende de:** `INF-04`.
- **Entregable:** whitelist de orígenes permitidos en CORS y certificados TLS configurados en Nginx para todo el tráfico hacia el frontend y la API.
- **Criterio de éxito:** una petición desde un origen no autorizado es rechazada por CORS y el tráfico hacia el sistema se sirve únicamente sobre HTTPS.


### `BAC-16B` - Adaptador Go `SmtpEmailService` para envío real de credenciales y códigos OTP (`RF-09` / `RF-13`)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `BAC-16`, `BAC-19`, `INF-08`.
- **Problema y contexto:** `BAC-16` dejó creado el puerto hexagonal `ports.EmailService`, pero solo implementó `MockEmailService` por consola. El backend necesita el adaptador SMTP real para enviar las contraseñas temporales y los códigos de recuperación de 6 dígitos.
- **Entregable:**
  1. Implementar `SmtpEmailService` en `backend/internal/adapters/secondary/email/smtp_service.go` cumpliendo la interfaz `ports.EmailService` (`EnviarCredencialesTemporales` y `EnviarCodigoRecuperacion`) con soporte `STARTTLS`/`TLS`.
  2. En `cmd/api/main.go`, instanciar `SmtpEmailService` cuando `EMAIL_PROVIDER=smtp` y mantener `MockEmailService` cuando `EMAIL_PROVIDER=mock` (usado por los tests automatizados).
  3. Estructurar los cuerpos de correo para alta de cuenta (`RF-09`), reset administrativo (`BAC-15`) y código OTP de 6 dígitos (`RF-13`).
- **Criterio de éxito:** Con `EMAIL_PROVIDER=smtp`, el backend envía correos reales; ante un fallo de entrega aborta la operación y devuelve `502 EMAIL_DELIVERY_FAILED`; la suite de tests sigue pasando en modo `mock`.


### `BAC-17B` - Corrección de lógica de inserción en `sesiones_activas` (1 sesión = 1 registro) y almacenamiento en Redis con TTL

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2,5 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `INF-06`, `BAC-17`.
- **Problema y diagnóstico en código (`backend/internal/core/services/auth_service.go`):**
  Actualmente el código de `auth_service.go` tiene un diseño ilógico que multiplica las filas en `sesiones_activas` por cada usuario:
  1. En `Login()` (líneas 76-84) inserta la **Fila 1** para el `jwtTemporal` (pre-2FA) con `activa = true`.
  2. En `VerificarTotp()` (líneas 249-313), en vez de transformar esa sesión o eliminar la temporal, deja la **Fila 1** con `activa = true`, inserta la **Fila 2** para el `refreshToken` y encima inserta la **Fila 3** para el `accessToken`. **Un solo inicio de sesión genera 3 filas en la base de datos.**
  3. En `RefrescarToken()` (líneas 388-399), cada vez que el frontend renueva el `accessToken`, el backend hace **otro `INSERT`** (`Fila 4, Fila 5, Fila 6...`) dejando todas las filas de `accessToken` anteriores con `activa = true`.
  4. En `CerrarSesion()` (líneas 433-438), solo pasa a `activa = false` el último `access` y `refresh`, dejando huérfanas la fila pre-2FA y todas las filas de renovaciones intermedias.
- **Entregable (Solución arquitectónica):**
  1. **Regla de oro (1 Sesión de Usuario = 1 único registro activo):**
     - **Paso Pre-2FA (`Login`):** Guardar el `jtiTemporal` exclusivamente en Redis (`SET auth:pre2fa:<jti> <usuario_id> EX 300`) con expiración automática de 5 minutos (o si se usa PostgreSQL, que sea la única fila creada que luego se actualiza en el paso 2FA). **Nunca dejar filas pre-2FA sueltas.**
     - **Paso Post-2FA (`VerificarTotp`):** Consumir y eliminar (`DEL`) el `jtiTemporal` pre-2FA. Crear **1 única sesión** para el navegador del usuario (en Redis `SET auth:session:<session_id> ... EX <ttl>` y **1 sola fila** en `sesiones_activas` que represente la sesión activa, no 2 filas separadas). Actualizar en ese mismo acto `fecha_ultimo_acceso = NOW()` en la tabla `usuarios`.
     - **Paso Renovación (`RefrescarToken`):** **Prohibido hacer `INSERT` en `RefrescarToken`**. Al renovar el token, hacer `UPDATE` sobre el **mismo registro existente** de esa sesión (actualizando el `jti_token` vigente y su `fecha_expiracion` en la fila única y en Redis). Así, aunque un usuario renueve su token 100 veces en el día, sigue ocupando **1 sola fila**.
     - **Paso Cierre (`CerrarSesion` / Expiración):** Eliminar la clave de Redis (`DEL`) y eliminar (`DELETE`) o desactivar esa única fila en `sesiones_activas`.
- **Criterio de éxito:** Un usuario que inicia sesión, verifica 2FA y refresca su token 10 veces genera **exactamente 1 sesión activa** (no 12 filas); al cerrar sesión o vencer el TTL, no quedan filas residuales activas y `fecha_ultimo_acceso` se persiste correctamente en `usuarios`.

### `BAC-18B` - Índice parcial y purga en `sesiones_activas`, y particionamiento trimestral en `auditoria` (PostgreSQL)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `BAC-17B`, `FIX-23`.
- **Problema y contexto:**
  1. Las filas históricas o con `activa = false` en `sesiones_activas` penalizan las lecturas en PostgreSQL si no existe un índice parcial ni una purga de registros vencidos.
  2. La tabla `auditoria` es *append-only* (`FIX-23`, no admite `DELETE` bajo ningún concepto) y en la Etapa 1 registrará múltiples eventos por cada operación de Proxmox VE (`PENDING` y `SUCCESS`/`FAILED`). Sin particionamiento por fechas e índices compuestos, las consultas de `GET /api/admin/audit` y la exportación CSV se degradarán progresivamente.
- **Entregable:**
  1. Crear en PostgreSQL (`init.sql` / migración) el índice parcial para `sesiones_activas`:
     ```sql
     CREATE INDEX IF NOT EXISTS idx_sesiones_activas_vigentes
       ON sesiones_activas (jti_token, usuario_id)
       WHERE activa = true;
     ```
  2. Agregar una rutina de limpieza en el backend (ticker cada 1 hora) que ejecute `DELETE FROM sesiones_activas WHERE activa = false OR fecha_expiracion < NOW();`.
  3. Configurar en PostgreSQL el **particionamiento declarativo trimestral por rango de fechas** sobre la tabla `auditoria` (`PARTITION BY RANGE (fecha_hora)`), creando las particiones trimestrales (`auditoria_2026_q3`, `auditoria_2026_q4`, `auditoria_2027_q1` y `auditoria_default`) junto con los índices compuestos `(fecha_hora DESC, accion, resultado)` y `(usuario_id, fecha_hora DESC)`.
- **Criterio de éxito:** Las sesiones muertas se purgan automáticamente de PostgreSQL; la tabla `auditoria` opera sobre particiones trimestrales manteniendo tiempos de consulta constantes (< 20 ms) e inmutabilidad append-only.


### `FIX-20` - Sincronización de políticas de complejidad y distinción UX entre Cambio y Restablecimiento (`FRN-10` / `FRN-12`) (Frontend / UX)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 3).
- **Depende de:** `FRN-10`, `FRN-12` y `BAC-02`.
- **Problema:**
  1. *Disparidad en reglas de validación:* `ChangePassword.tsx` (`validateChangePasswordFields`, líneas 22-34) valida únicamente que la longitud esté entre 8 y 12 caracteres. Omite las reglas que exige el backend en `crypto.ValidarComplejidadContrasena` (`backend/internal/infrastructure/crypto/password.go:89`): al menos una mayúscula, un número y un carácter especial `!@#$%^&*-_=+`. Si el usuario ingresa una clave que no cumple estas reglas, el cliente la envía y el backend responde 400 `PASSWORD_CHANGE_FAILED`.
     - **Evidencia (23/09/2026):** en `test/front/password-change.test.tsx`, los casos *"no llama a la API si la contraseña nueva no cumple la complejidad del backend (sin mayúscula / sin dígito / sin carácter especial)"* fallan, porque el cliente envía `PUT /api/account/password` con `nueva1234!`, `NuevaClave!` y `Nueva12345`.
  2. *Ambigüedad visual:* Falta claridad en los textos de ayuda de los formularios para distinguir que en el **Cambio** se requiere la clave temporal previa, mientras que en el **Restablecimiento** solo se define una nueva contraseña.
- **Entregable:**
  1. Implementar un validador unificado de contraseñas (`validatePasswordComplexity`) en el archivo existente `frontend/centinela/src/utils/validators.js`, que replique exactamente los criterios del backend:
     - Longitud: 8 a 12 caracteres.
     - Al menos una letra mayúscula.
     - Al menos un dígito numérico.
     - Al menos un símbolo permitido (`[!@#$%^&*-_=+]`).
  2. Integrar el validador unificado en `ChangePassword.tsx` y `RecoverPassword.tsx` (paso 3), mostrando el error junto al campo y textos de ayuda explicativos debajo. `CrearUsuarios.tsx` no lo necesita: desde `BAC-16` la clave temporal la genera el backend y se envía por correo, y los campos de contraseña de esa vista se eliminan en `FIX-24`.
- **Criterio de éxito:** Las validaciones de cliente previenen el envío de contraseñas no conformes; el usuario recibe retroalimentación inmediata; se eliminan los errores 400 por rechazo de complejidad en el backend. Los 3 casos de complejidad de `test/front/password-change.test.tsx` pasan.
- **Nota:** la corrección de la referencia de `FRN-10` a `PUT /api/account/password` en la documentación ya se hizo el 23/09/2026, en `terminado.md`, y no forma parte de este FIX.


### `FRN-18` (`SEC-04B`) - Integración en Frontend de niveles de acceso por instancia (`READ_ONLY` / `FULL_ACCESS`)

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (inmediatamente posterior a `SEC-03` y `SEC-04`).
- **Depende de:** `SEC-03`, `SEC-04` y `FIX-14`.
- **Problema y evidencia (`cierre-fase-base.md`):** `SEC-04` implementa `nivel_acceso` (`FULL_ACCESS` y `READ_ONLY`) en el backend, pero ninguna tarea de Frontend tenía asignado enviar ese nivel en `PUT /api/admin/users/:id/permissions`, leerlo al abrir la ficha del usuario ni exponerlo en `usePermissions()` (`SEC-03`) para distinguir quién puede solo ver una máquina de quién puede apagarla o reiniciarla.
- **Entregable:**
  1. En `detailsUserPage.tsx` y `rolesAndPermissions.tsx`, leer el `nivelAcceso` de cada instancia desde `GET /api/admin/users/:id/permissions` y enviar `{ vmid, nivelAcceso: 'FULL_ACCESS' | 'READ_ONLY' }` en `PUT /api/admin/users/:id/permissions`.
  2. En el contexto `SEC-03` (`usePermissions()`), agregar el helper `canOperateInstance(vmid: number): boolean` (devuelve `true` solo si es `ADMIN` o si tiene `FULL_ACCESS` sobre ese `vmid`), diferenciándolo de `canAccessInstance(vmid)` (que devuelve `true` tanto para `READ_ONLY` como `FULL_ACCESS`).
- **Criterio de éxito:** El administrador puede asignar y guardar el nivel "Solo lectura" o "Control total" por instancia desde la UI; `canOperateInstance` retorna `false` para instancias en modo `READ_ONLY`.
