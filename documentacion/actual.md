
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Verificación del 23/09/2026** ([test/informe.md](../test/informe.md)): las tareas completas (`BAC-14`, `FIX-16`, `BAC-17`, `FRN-13`, `SEC-02`, `FIX-14`) y las implementadas con problemas (`FRN-10`, `FRN-14`) pasaron a [`terminado.md`](terminado.md), y los problemas se registraron como `FIX-21` a `FIX-25`.
>
> **Verificación del 25/09/2026:** `SEC-01`, `FIX-17` y `FIX-25` pasaron todas sus pruebas y se movieron a [`terminado.md`](terminado.md). `SEC-04` y `FIX-24` están implementadas con un problema cada una: también pasaron a `terminado.md`, y sus correcciones son `FIX-26` y `FIX-27` en [`futuro.md`](futuro.md).

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (25/09/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación y aviso de correo, edición del perfil (`FIX-25`, terminada), baja con confirmación y 403 al `OPERATOR`. Solo queda el mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta (`FIX-27`).

---

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

