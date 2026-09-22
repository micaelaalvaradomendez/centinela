
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:


### `BAC-14` - Lectura mínima del inventario de Proxmox

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 3 h
- **Ventana propuesta:** 17/09/2026, 11:00-14:00
- **Depende de:** `BAC-08` y de las credenciales de lectura de Proxmox VE.
- **Entregable:** `GET /api/instances` consumiendo `/cluster/resources` o `/nodes/{node}/resources`, con una respuesta normalizada mínima que incluya ID, nombre, tipo, nodo y estado. Un administrador recibe todo el inventario y un operador solo las instancias asignadas.
- **Criterio de éxito:** `FRN-07` puede cargar IDs reales de VMs y LXC; una cuenta no puede descubrir instancias fuera de su alcance y los errores de Proxmox se traducen a una respuesta HTTP controlada.



### `FRN-10` - Cambio obligatorio de contraseña temporal

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2 h
- **Ventana propuesta:** 16/09/2026, 11:00-13:00
- **Depende de:** `BAC-12` y `FRN-04`.
- **Entregable:** vista de nueva contraseña y confirmación que detecte `must_change_password`, bloquee la navegación general y llame a `POST /api/auth/change-password`.
- **Criterio de éxito:** una cuenta con contraseña temporal solo puede cerrar sesión o cambiarla; después del cambio continúa al enrolamiento o validación 2FA que corresponda.

### `FRN-11` - Acciones administrativas de recuperación

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2 h
- **Ventana propuesta:** 18/09/2026, 12:00-14:00
- **Depende de:** `FRN-05`, `BAC-13` y `BAC-15`.
- **Entregable:** acciones separadas para restablecer contraseña y 2FA desde el panel de usuarios, ambas con confirmación explícita, estado de carga y notificación del resultado.
- **Criterio de éxito:** un administrador puede iniciar cada recuperación sin confundir sus efectos; la tabla refleja que el 2FA quedó desvinculado y nunca muestra secretos ni hashes.

### `FRN-12` - Vistas de recuperación de contraseña (RF-13)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 3 h
- **Ventana propuesta:** A definir (posterior a `BAC-20`).
- **Depende de:** `BAC-19`, `BAC-20` y el maquetado existente de `RecoverPassword.tsx`.
- **Entregable:** conectar `RecoverPassword.tsx` al flujo real: paso de ingreso de correo, paso de ingreso del código de seis dígitos y paso de nueva contraseña, con manejo de errores del backend en cada paso.
- **Criterio de éxito:** una cuenta puede recuperar el acceso sin intervención de un administrador, y los errores de código inválido o vencido se muestran junto al campo correspondiente.

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

**Hito:** *Control de Acceso Basado en Recursos*
* El administrador abre un usuario en `FRN-07`, el frontend lista las instancias de Proxmox (`BAC-14`), selecciona un subconjunto y guarda (`BAC-07`).
* Al iniciar sesión como ese Operador, `GET /api/instances` solo devuelve las instancias que tiene permitidas.
* Si intenta forzar una petición sobre una instancia no asignada, el backend rechaza con `403` (`BAC-08`) y el frontend muestra el error en un toast sin cerrar la sesión (`FRN-08`).
---


### `BAC-17` - Cierre de sesión y revocación atómica de sesiones (Backend)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** A definir.
- **Depende de:** `BAC-03` y `BAC-04`.
- **Estado de implementación actual en submódulo:** En `backend/cmd/api/main.go`, el endpoint `POST /api/auth/logout` está montado sin middleware de autenticación, por lo que el JTI del token de acceso (`ContextKeyJTI`) nunca se inyecta en el contexto y llega vacío a `authService.CerrarSesion`. En consecuencia, solo se desactiva el registro del `refreshToken` en `sesiones_activas`, mientras que el `accessToken` sigue figurando como activo en base de datos hasta que expire.
- **Entregable:**
  1. Permitir que `POST /api/auth/logout` reciba el encabezado `Authorization: Bearer <accessToken>` (o extraerlo del contexto) para obtener el JTI del token de acceso y revocar atómicamente ambas sesiones (`jtiRefresh` y `jtiAccess`) en la tabla `sesiones_activas`.
  2. Asegurar que tras el logout, cualquier llamada subsiguiente con el `accessToken` reciba inmediatamente `401 Unauthorized` (`TOKEN_REVOKED`) en el middleware `RequireAuth`.
  3. Registrar formalmente la acción en la tabla de auditoría con la acción `LOGOUT`.
- **Criterio de éxito:** tras ejecutar `POST /api/auth/logout`, tanto el refresh token como el access token quedan con `activa: false` en la base de datos; cualquier petición con el access token revocado es rechazada con código 401; la auditoría registra el evento de cierre de sesión.


### `FRN-13` - Flujo integral de logout y limpieza de sesión en cliente (Frontend)

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 2 h
- **Ventana propuesta:** A definir.
- **Depende de:** `BAC-17`.
- **Estado de implementación actual en submódulo:** En `frontend/centinela`, el servicio `logoutSession()` (`src/components/features/auth/services/authService.ts`) y el hook `useLogout()` (`src/components/features/auth/hooks/useAuth.ts`) ejecutan `POST /api/auth/logout` enviando únicamente `{ refreshToken }` y limpian el almacenamiento local. El botón "Cerrar sesión" se encuentra en el popover del usuario en `Header.tsx`. Sin embargo, falta garantizar el envío del token de acceso en la cabecera, gestionar el caso en que el refresh token no exista o haya expirado, y reaccionar ante eventos globales de revocación.
- **Entregable:**
  1. Asegurar que `logoutSession()` incluya la cabecera `Authorization: Bearer <accessToken>` junto al body `{ refreshToken }` para permitir al backend la revocación de ambos tokens.
  2. Garantizar que `useLogout()` limpie exhaustivamente `sessionStorage` (`centinela_access`, `centinela_refresh`, `centinela_user`, `centinela_pending_login`) sin importar si la llamada a la API tiene éxito o falla por red/expiración.
  3. Implementar un listener global para el evento `centinela:api-unauthorized` (o respuestas 401 `TOKEN_REVOKED`) que invoque la limpieza de sesión y redirección inmediata a `/login` para expulsar al usuario si su sesión fue revocada remotamente o cerrada en otra pestaña.
- **Criterio de éxito:** el usuario presiona "Cerrar sesión" en el menú de usuario, el backend revoca ambos tokens en la base de datos, el cliente elimina todas las credenciales de `sessionStorage` y redirige a `/login` con `replace: true`, impidiendo volver atrás mediante el historial del navegador; si el backend revoca la sesión, el frontend detecta el 401 y expulsa al usuario al login.


### `SEC-01` - Refresh token en cookie HttpOnly en backend

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** A definir.
- **Depende de:** `BAC-17`.
- **Estado actual:** la revocación de sesiones y la renovación de tokens están implementadas, pero el backend todavía recibe el refresh token en JSON y devuelve el token en el body. No existe emisión ni lectura de cookie HttpOnly.
- **Entregable:**
	1. Emitir el `refreshToken` mediante `Set-Cookie` al completar `POST /api/auth/2fa/verify`.
	2. Leerlo desde la cookie en `POST /api/auth/refresh` y `POST /api/auth/logout`.
	3. Dejar de devolver el refresh token en las respuestas JSON públicas.
	4. Configurar `HttpOnly`, `Secure`, `SameSite` y `Path` de forma segura y configurable para desarrollo y producción.
	5. Eliminar la cookie al cerrar sesión o cuando la sesión sea inválida.
- **Criterio de éxito:** el backend renueva y revoca sesiones usando exclusivamente la cookie HttpOnly; el refresh token no aparece en cuerpos JSON ni en logs; una cookie inválida o vencida produce un error controlado.

### `SEC-02` - Cliente frontend compatible con refresh token HttpOnly

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2 h
- **Ventana propuesta:** A definir, posterior a `SEC-01`.
- **Depende de:** `SEC-01` y `BAC-17`.
- **Estado actual:** `credentials: 'include'` ya está configurado en el cliente HTTP, pero el frontend todavía guarda el `refreshToken` en `sessionStorage`, lo lee para renovar la sesión y lo envía en el body de refresh/logout.
- **Entregable:**
	1. Eliminar el almacenamiento, lectura y tipado de `refreshToken` en el frontend.
	2. Mantener `credentials: 'include'` en login/2FA, refresh y logout.
	3. Adaptar la validación de respuestas para aceptar respuestas que ya no incluyan `refreshToken`.
	4. Adaptar el interceptor para renovar la sesión sin construir un body con el token.
	5. Mantener el `accessToken` en `sessionStorage` mientras esta tarea se limite al refresh token.
- **Criterio de éxito:** JavaScript no puede leer el refresh token desde storage, memoria de aplicación ni respuestas HTTP; la renovación y el logout funcionan mediante cookies y la sesión local conserva únicamente el access token.


### `FIX-14` - Integración frontend del selector de instancias (`FRN-07`)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 4 h
- **Ventana propuesta:** A definir.
- **Depende de:** `BAC-07` y `BAC-14`.
- **Problema y evidencia en pruebas (`test/front/admin-users.test.tsx`):**
  1. *Falla de consulta:* `FRN-07 - selector de asignación de instancias > consulta GET /api/instances con Authorization Bearer para listar instancias` falló con `AssertionError: expected undefined to be defined` en la comprobación de llamadas a `fetchMock`. `UserDetail.tsx` nunca realiza la consulta HTTP a `/instances` para obtener el listado dinámico de máquinas virtuales y contenedores LXC.
  2. *Falla de guardado:* `FRN-07 - selector de asignación de instancias > permite seleccionar VMIDs y enviarlos al endpoint de permisos` falló con `AssertionError: expected undefined to be defined` tras hacer click en el botón de guardado. `handleSaveChanges()` solo despacha la actualización de datos generales a `PUT /api/admin/users/:id`, omitiendo el despacho de los VMIDs seleccionados hacia el endpoint de permisos.
  3. *Diagnóstico en código:* en [frontend/centinela/src/pages/UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx), el subcomponente `AccessList()` renderiza elementos estáticos fijos (`['Ubuntu Server (101)', 'Desarrollo (104)', ...]`) sin estado reactivo, sin consumir `apiClient.get('/instances')` y sin vincular los checkboxes a una mutación real.
- **Entregable:**
  1. Implementar en `UserDetail.tsx` (o hook dedicado) la carga asíncrona de instancias consumiendo `GET /api/instances` con token de autorización Bearer.
  2. Mapear los VMIDs actualmente asignados al usuario (obtenidos desde `GET /api/admin/users/:id/permissions` o `instanciasPermitidas` del DTO) en el estado de selección de la pestaña "Roles y permisos".
  3. Al presionar "Guardar cambios", despachar la petición `PUT /api/admin/users/:id/permissions` con el payload `{ vmids: number[] }` de forma atómica junto a la edición del perfil de usuario, mostrando estados de carga, notificación visual (toast) y control de errores.
- **Criterio de éxito:** la suite `test/front/admin-users.test.tsx` pasa al 100% (2 pruebas de `FRN-07` en verde), verificando la llamada con Bearer a `GET /api/instances` y el despacho de `PUT /api/admin/users/:id/permissions` con el array de VMIDs seleccionados.

### `FIX-16` - Montar middleware de autorización por recurso en rutas de instancias (`BAC-08`)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (junto a `BAC-14`).
- **Depende de:** `BAC-08` y `BAC-14`.
- **Problema y evidencia en pruebas (`test/back/resource_access_acceptance_test.go`):** la prueba `BAC-08_rechaza_instancia_no_asignada_con_403_sin_llamar_a_Proxmox` falló con código HTTP `404 Not Found` en lugar de `403 Forbidden` al consultar `GET /instances/9999` con token de un usuario con rol `OPERATOR`.
- **Diagnóstico y causa raíz:** el middleware `RequireInstanceAccess` (`internal/adapters/primary/http/middleware/instance_guard.go`), el puerto `InstanceRepository` (`internal/core/ports/instance_port.go`) y el repositorio PostgreSQL `postgres.NewInstanceRepository(db)` (`internal/adapters/secondary/postgres/instance_repository.go`) se encuentran completamente programados e instanciados en `main.go`. No obstante, el grupo de rutas `/instances` permanece comentado en `backend/cmd/api/main.go` a la espera de que se incorporen los handlers de Proxmox VE (`BAC-14`). Al no estar montadas las rutas en Gin, la solicitud no alcanza el guard y responde `404`.
- **Entregable:**
  1. Habilitar y montar el grupo de rutas `/instances` en `cmd/api/main.go` bajo `middleware.RequireAuth()`.
  2. Aplicar el middleware `RequireInstanceAccess(instanceRepo, "vmid")` en las operaciones por recurso (`GET /instances/:vmid`, `POST /instances/:vmid/start`, etc.).
  3. Asegurar que las solicitudes de operadores sobre VMIDs no autorizados sean interceptadas tempranamente con código `403 Forbidden` y `{ "errorCode": "INSTANCE_ACCESS_DENIED" }`, abortando la cadena antes de invocar cualquier llamada hacia Proxmox VE.
- **Criterio de éxito:** un operador que intente acceder a un VMID no asignado recibe `403 Forbidden` con `INSTANCE_ACCESS_DENIED`; la aserción de `BAC-08_rechaza_instancia_no_asignada_con_403_sin_llamar_a_Proxmox` en `resource_access_acceptance_test.go` pasa al 100%.

### `SEC-03` - Contexto y sistema reactivo de permisos en Frontend (UI/UX RBAC + Resources)

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 2,5 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 2).
- **Depende de:** `FRN-04`, `BAC-07` y `BAC-09`.
- **Problema:** En el frontend actual, el control de roles y permisos está disperso y acoplado únicamente a los loaders de rutas (`loadAdminSession`). No existe un mecanismo reactivo a nivel de componentes (`usePermissions` / `PermissionGate`) para ocultar o deshabilitar condicionalmente acciones según el rol (`ADMIN` vs `OPERATOR`) o según las instancias asignadas al operador. Esto genera que en vistas como `UserDetail.tsx` o `Header.tsx` se muestren controles estáticos o se dependa de que el backend rechace con 403, en lugar de ofrecer una experiencia fluida y consistente.
- **Entregable:**
  1. Crear un hook y contexto `usePermissions()` / `useAuthUser()` en `frontend/centinela/src/context/` que exponga helpers como `isAdmin`, `isOperator`, `canAccessInstance(vmid: number)` y `hasRole(role: string)`.
  2. Crear un componente wrapper `<PermissionGate requiredRole="ADMIN" fallback={...}>` para condicionar la renderización de botones y secciones administrativas (ej: botón "Eliminar usuario", accesos a auditoría, etc.).
  3. Integrar la reactividad en el menú de navegación y en el header para que las opciones no permitidas a un operador no aparezcan en la interfaz.
- **Criterio de éxito:** Si un operador inicia sesión, la interfaz no muestra accesos directos ni botones exclusivos de administrador; si intenta interactuar con un componente restringido, el helper `canAccessInstance` evalúa en memoria las instancias permitidas cargadas en la sesión; las pruebas unitarias de componentes verifican el render condicional.

### `SEC-04` - Esquema extensible de niveles de acceso a recursos en Backend y Modelo de Datos

- **Área:** Backend
- **Asignado:** Lisandro / Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 2).
- **Depende de:** `BAC-05`, `BAC-07` y `BAC-08`.
- **Problema:** La tabla `permisos_instancia` solo almacena `(usuario_id, vmid_proxmox)` como una relación binaria (asignado o no asignado). Sin embargo, en el documento de diseño `Diseño de endpoints para front.md` y en las siguientes etapas (Etapa 1 Ciclo de vida y Etapa 3 Snapshots) se proyectan niveles de acceso sobre las instancias (por ejemplo `FULL_ACCESS` para operar energía/snapshots vs `READ_ONLY` para monitoreo de métricas sin emitir órdenes destructivas). Si no se sienta esta base en el esquema y en el guard ahora, agregar niveles de acceso en la Etapa 1 requerirá migrar tablas en producción y alterar contratos.
- **Entregable:**
  1. Agregar en `domain.PermisoInstancia` la columna `nivel_acceso` (VARCHAR(30) default `'FULL_ACCESS'`) con restricción CHECK o enum para los valores `FULL_ACCESS` y `READ_ONLY`.
  2. Actualizar el puerto `InstanceRepository.VerificarAcceso` para aceptar opcionalmente el nivel de acceso requerido (`requiredLevel: string`).
  3. Extender el middleware `RequireInstanceAccess(repo, paramName, requiredLevel)` para permitir guards como `RequireInstanceAccess(repo, "vmid", "FULL_ACCESS")` en endpoints mutantes (`POST /instances/:vmid/start`, `POST /instances/:vmid/stop`) y tolerar `READ_ONLY` en consultas (`GET /instances/:vmid`).
  4. Mantener retrocompatibilidad total: si el payload de `PUT /permissions` solo envía `vmids: [101]`, asignar `FULL_ACCESS` por defecto.
- **Criterio de éxito:** La migración crea el campo sin romper registros previos; el middleware `RequireInstanceAccess` verifica tanto la pertenencia de la instancia como el nivel de permiso; si un usuario tiene permiso `READ_ONLY` sobre la VM 101, puede consultar su estado pero recibe 403 al intentar ejecutar una acción de apagado/encendido.

### `FRN-14` - Suite de pruebas unitarias y de integración para la vista de Auditoría (`Auditoria.tsx`)

- **Área:** Frontend
- **Asignada:** Belinda / Luz
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Fase Base / Cierre de Auditoría).
- **Depende de:** `BAC-18` y `FRN-03`.
- **Problema:** El backend ya tiene verificada la auditoría append-only y la exportación CSV (`BAC-18` en `password_recovery_acceptance_test.go`), y el frontend tiene maquetado [frontend/centinela/src/pages/Auditoria.tsx](frontend/centinela/src/pages/Auditoria.tsx), pero **no existe ninguna prueba automatizada en `test/front`** que valide que la tabla cargue correctamente los registros de `/admin/audit`, que los filtros (por acción, resultado y fechas) envíen los query params correctos y que el botón de exportación descargue el archivo CSV.
- **Entregable:**
  1. Crear la suite `test/front/audit.test.tsx`.
  2. Testear que la vista realiza la petición `GET /api/admin/audit?pagina=1&tamano=10` con cabecera `Authorization: Bearer`.
  3. Testear que la interacción con los filtros reactivos y la paginación actualiza los query parameters.
  4. Testear que el botón "Exportar CSV" dispara la llamada a `/api/admin/audit/export?formato=csv`.
  5. Testear que si un operador entra a la ruta `/auditoria`, el guard lo redirige a `/dashboard`.
- **Criterio de éxito:** `npm test` en `test/front` ejecuta y aprueba los casos de prueba de `audit.test.tsx`, garantizando que la auditoría administrativa esté verificada de punta a punta.