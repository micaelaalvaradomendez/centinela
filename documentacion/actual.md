
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
