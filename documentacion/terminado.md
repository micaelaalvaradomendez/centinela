

### `BAC-01` - Script Docker Compose de base de datos local

- **Área:** Backend / Infraestructura
- **Asignados:** Tayra y Lucas
- **Estimación:** 2 h
- **Ventana propuesta:** 10/09/2026, 14:00-16:00
- **Depende de:** ninguna.
- **Entregable:** PostgreSQL local mediante Docker Compose, tablas mínimas de usuarios y roles `ADMIN`/`OPERATOR`, y usuario de prueba.
- **Criterio de éxito:** el equipo puede levantar la base con `docker compose up` y conectarse usando las credenciales documentadas.

### `BAC-02` - Modelo y hashing de contraseñas

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 10/09/2026, 14:00-16:00
- **Depende de:** ninguna; puede ejecutarse en paralelo con `BAC-01`.
- **Entregable:** modelo de usuario y verificación de contraseñas con bcrypt o Argon2.
- **Criterio de éxito:** la contraseña nunca se persiste en texto plano y una credencial válida se distingue de una inválida mediante pruebas automatizadas o manuales documentadas.

### `BAC-03` - Endpoint `POST /api/auth/login`

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 3 h
- **Ventana propuesta:** 11/09/2026, 09:00-12:00
- **Depende de:** `BAC-01` y `BAC-02`.
- **Entregable:** endpoint que recibe usuario o correo y contraseña, valida contra PostgreSQL y devuelve un JWT firmado con `id` y `role`.
- **Criterio de éxito:** el usuario de prueba puede iniciar sesión y el JWT se valida correctamente con la clave configurada.

### `BAC-04` - Manejo de errores HTTP

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Ventana propuesta:** 11/09/2026, 12:00-13:30
- **Depende de:** `BAC-03`.
- **Entregable:** respuestas JSON unificadas con `400` para payload incompleto y `401` para credenciales inválidas.
- **Criterio de éxito:** los casos válidos e inválidos devuelven código HTTP y estructura JSON documentados.

### `LOGIN-01` - Flujo 2FA/TOTP en memoria

- **Área:** Backend
- **Asignado:** Tayra y Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 11/09/2026, 13:30-15:30
- **Depende de:** `BAC-03` y `BAC-04`.
- **Entregable:** generación de secreto TOTP temporal, validación de código de seis dígitos y estado de verificación en memoria. No reemplaza todavía la persistencia segura definitiva del RF-01.
- **Criterio de éxito:** el login solicita el TOTP después de validar la contraseña y rechaza códigos inválidos o vencidos.


### `FRN-01` - Maquetado del formulario de login

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 3 h
- **Ventana propuesta:** 10/09/2026, 14:00-17:00
- **Depende de:** ninguna.
- **Entregable:** campos de usuario o correo, contraseña, botón de envío y estado visual de carga.
- **Criterio de éxito:** la vista se renderiza y permite recorrer los estados inicial, cargando, éxito y error.

### `FRN-02` - Validación en cliente del login

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 2 h
- **Ventana propuesta:** 10/09/2026, 17:00-19:00
- **Depende de:** `FRN-01`.
- **Entregable:** validación de formato y mensajes de error antes de enviar la petición.
- **Criterio de éxito:** no se envían formularios incompletos y cada error se muestra junto al campo correspondiente.

### `FRN-03` - Navbar y rutas base

- **Área:** Frontend
- **Asignados:** Cristian y Luz
- **Estimación:** 3 h
- **Ventana propuesta:** 10/09/2026, 14:00-17:00
- **Depende de:** ninguna.
- **Entregable:** navegación y vistas iniciales de Dashboard e Instancias, aunque sean estáticas.
- **Criterio de éxito:** el usuario puede navegar entre las pantallas base sin romper la aplicación.


### `FRN-04` - Conexión del login con localhost

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 3 h
- **Ventana propuesta:** 11/09/2026, 15:30-18:30
- **Depende de:** `BAC-03`, `BAC-04` y `FRN-01`.
- **Entregable:** petición HTTP al backend local y almacenamiento del JWT en memoria.
- **Criterio de éxito:** el usuario de prueba puede iniciar sesión desde el frontend y el cliente conserva el JWT durante la sesión.

### `LOGIN-02` - Pantalla y validación del código TOTP

- **Área:** Frontend
- **Asignados:** Belinda y Luz
- **Estimación:** 2 h
- **Ventana propuesta:** 11/09/2026, 15:30-17:30
- **Depende de:** `LOGIN-01` y `FRN-01`.
- **Entregable:** pantalla de ingreso del código TOTP de seis dígitos y mensajes de error.
- **Criterio de éxito:** el frontend permite ingresar el código, muestra el estado de verificación y habilita la navegación solo cuando la validación es correcta.

### `LOGIN-03` - Integración y prueba del prototipo

- **Área:** Frontend / Backend
- **Asignados:** Cristian, Tayra y Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 11/09/2026, 18:30-20:30
- **Depende de:** `BAC-04`, `FRN-04` y `LOGIN-02`.
- **Entregable:** recorrido completo usuario de prueba -> contraseña -> TOTP en memoria -> pantalla base.
- **Criterio de éxito:** el flujo válido termina en Dashboard, el flujo inválido muestra un error y no permite acceder a las vistas protegidas.

### `BAC-09` - Endpoint de roles del sistema (GET /api/roles)

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1 h
- **Depende de:** `BAC-05`.
- **Entregable:** endpoint público o protegido que retorne el listado de roles con su identificador y descripción básica.
- **Criterio de éxito:** devuelve un JSON con los roles disponibles (`ADMIN`, `OPERATOR`) con código HTTP 200.

### `INF-03` - PostgreSQL persistente en el servidor

- **Área:** Infraestructura
- **Asignados:** Lucas y Nico
- **Estimación:** 3 h
- **Ventana propuesta:** 16/09/2026, 09:00-12:00
- **Depende de:** `BAC-05`.
- **Entregable:** PostgreSQL desplegado en el LXC o Docker de prueba, con volumen persistente y credenciales seguras.
- **Criterio de éxito:** la base es accesible por la red interna o VPN y los datos sobreviven al reinicio del contenedor.


### `INF-04` - Red interna y reverse proxy con Nginx

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 3 h
- **Ventana propuesta:** 16/09/2026, 13:00-16:00
- **Depende de:** `INF-03` y del backend desplegable.
- **Entregable:** Nginx enruta `/api/*` al backend y `/` al frontend; la subred virtual `vmbr1` comunica backend, base de datos y API de Proxmox VE.
- **Criterio de éxito:** el dominio o IP local resuelve, el frontend alcanza el backend y el backend alcanza PostgreSQL y Proxmox VE sin exponer la base públicamente.

### BAC-06B — Edición de Usuario y Cambio de Rol (PUT /api/admin/users/{id}):
- Asignado: Lisandro | Estimación: 2h
- Depende de: BAC-06
- Entregable: Endpoint PUT /api/admin/users/{id} para actualizar nombre, correo, estado (isActive) y rol (role: ADMIN u OPERATOR).  
- Criterio de éxito: Un administrador puede cambiarle el rol a un usuario o desactivarlo; los cambios se reflejan de inmediato en la base de datos.  

---

# REVISAR

### `BAC-05` - Migración y esquema relacional de usuarios y RBAC

- **Área:** Backend
- **Asignados:** Tayra y Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** 14/09/2026, 09:00-12:00
- **Depende de:** `BAC-01` y `LOGIN-03`.
- **Entregable:** migración SQL u ORM con `users` y `user_instances`. `users` debe incluir `id` UUID, `username`, `email`, `password_hash`, `role`, `totp_secret` protegido, `is_2fa_enabled` y `must_change_password`. `user_instances` debe tener `user_id`, `instance_id` y clave primaria compuesta.
- **Criterio de éxito:** la migración se ejecuta desde cero, crea las relaciones y carga un administrador inicial.

### `BAC-06` - CRUD de usuarios para administradores

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 4 h
- **Ventana propuesta:** 14/09/2026, 13:00-17:00
- **Depende de:** `BAC-05`.
- **Entregable:** `GET /api/admin/users`, `POST /api/admin/users` y `DELETE /api/admin/users/{id}`. La creación debe generar una contraseña temporal y establecer `must_change_password: true`.
- **Criterio de éxito:** un administrador puede listar, crear y desactivar usuarios; un operador recibe `403 Forbidden`.




### `FRN-05` - Panel de gestión de usuarios

> [!NOTE]
> **Estado: Completada (verificado el 23/09/2026).** El guard administrativo se resolvió con `FIX-12`: las rutas `/users`, `/users/new` y `/users/:userId` están bajo `loadAdminSession`. `test/front/admin-users.test.tsx`, bloque FRN-05 (3/3), verifica la redirección del OPERATOR, `GET /api/admin/users` con Bearer y la tabla con datos reales.

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 4 h
- **Ventana propuesta:** 14/09/2026, 09:00-13:00
- **Depende de:** `LOGIN-03` y contrato inicial de `BAC-06`.
- **Entregable:** vista `/admin/users` con nombre, correo, rol, estado de 2FA y acciones. Solo debe ser visible para administradores.
- **Criterio de éxito:** el administrador ve el listado real y un operador no puede acceder a la vista.

### `FRN-06` - Modal de creación y desactivación de usuarios

> [!WARNING]
> **Estado: Implementado con retroceso detectado (23/09/2026).** El alta envía `POST /api/admin/users` correctamente, pero hay dos problemas:
> - Desde `BAC-16` el backend ya no devuelve `contrasenaTemp`, y `CrearUsuarios.tsx` depende de ese campo para mostrar el resultado. Un alta exitosa **no muestra ninguna confirmación**.
> - "Eliminar usuario" (`detailsUserPage.tsx:110`) no tiene acción y no envía `DELETE`.
>
> Prueba: `test/front/admin-users.test.tsx`, bloque FRN-06. Defecto registrado en [`futuro.md`](futuro.md) como `FIX-24`.

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 3 h
- **Ventana propuesta:** 14/09/2026, 14:00-17:00
- **Depende de:** `FRN-05` y `BAC-06`.
- **Entregable:** modal con usuario, correo y rol; visualización controlada de la contraseña temporal; confirmación para desactivar usuarios y toasts de resultado.
- **Criterio de éxito:** el administrador puede completar altas y bajas desde la interfaz y los errores se muestran de forma comprensible.

### FRN-06B — Modal de Edición de Usuario y Cambio de Rol:

> [!NOTE]
> **Estado: Completada (verificado el 25/09/2026).** El retroceso detectado el 23/09 (la edición del perfil no se persistía) se resolvió con `FIX-25`. Prueba: `test/front/admin-users.test.tsx`, bloque FRN-06B.
- Asignado: Luz | Estimación: 2h
- Depende de: FRN-05, FRN-06 y BAC-06B
- Entregable: Formulario precargado con los datos del usuario seleccionado para modificar su información básica y cambiar su rol mediante un desplegable.  
- Criterio de éxito: Al confirmar la edición, se llama a PUT /api/admin/users/{id}, se actualiza la tabla y se muestra una alerta visual (toast) de éxito.  


### `FIX-07` - El frontend lee `body.code` en vez de `errorCode`

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1 h
- **Ventana propuesta:** A definir.
- **Depende de:** ninguna.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 6, hallazgo H1):** el backend responde los errores con el campo `errorCode` (por ejemplo `AUTH_FAILED`), pero `api.js` intenta leer `body.code`, que no existe en la respuesta real.
- **Entregable:** corregir `api.js` (y cualquier consumidor) para leer `errorCode` del cuerpo de la respuesta.
- **Criterio de éxito:** los mensajes de error específicos del backend (por ejemplo credenciales inválidas o TOTP vencido) se muestran correctamente en la interfaz en vez de un error genérico.

### `FRN-09` - Vinculación 2FA mediante QR

- **Área:** Frontend
- **Asignadas:** Belinda y Luz
- **Estimación:** 3 h
- **Ventana propuesta:** 16/09/2026, 14:00-17:00
- **Depende de:** `BAC-10`, `BAC-11` y `LOGIN-02`.
- **Entregable:** vista o modal obligatorio para cuentas sin 2FA, con QR, clave de vinculación manual, ingreso del código de seis dígitos y estados de carga y error.
- **Criterio de éxito:** la cuenta no puede acceder a las rutas protegidas hasta confirmar un código válido; al finalizar, continúa el login sin mostrar nuevamente el secreto.

### `FIX-01` - Colisión de índice único rompe el login repetido (CRÍTICA)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (máxima prioridad, antes que cualquier otra tarea de `BAC-05` en adelante).
- **Depende de:** ninguna; bloquea la verificación de `BAC-06`, `BAC-06B`, `BAC-07`, `BAC-08`, `BAC-09`, `BAC-10`, `BAC-11`, `BAC-14` y `LOGIN-03`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 3):** `models.go` reutiliza el mismo nombre `uniqueIndex:idx_usuario_vmid` en cuatro tablas (`SesionActiva.UsuarioID`, `PermisoInstancia.UsuarioID+VmidProxmox`, `Auditoria.UsuarioID`, `Notificacion.UsuarioID`). El índice terminó aplicado sobre `sesiones_activas.usuario_id`, así que la segunda sesión del mismo usuario viola la unicidad y el login devuelve `401 AUTH_FAILED` con `duplicate key value violates unique constraint "idx_usuario_vmid"`. Al mismo tiempo, el unique real que debía existir sobre `permisos_instancia(usuario_id, vmid_proxmox)` (parte del criterio de éxito de `BAC-05`) nunca se creó.
- **Entregable:** renombrar el índice compuesto de `PermisoInstancia` a algo único (por ejemplo `idx_permisos_usuario_vmid`), quitar el `uniqueIndex` de `SesionActiva.UsuarioID` (una sesión no es única por usuario) y agregar el constraint real faltante sobre `permisos_instancia(usuario_id, vmid_proxmox)`.
- **Criterio de éxito:** un mismo usuario puede loguearse más de una vez sin recibir `401`; `test/back` deja de reportar la colisión de `idx_usuario_vmid` en `backend_acceptance_test.go`; el chequeo de esquema de `BAC-05` confirma el unique real sobre `permisos_instancia`.


### `FIX-06` - Modal de alta y edición de usuarios sin funcionalidad real (`FRN-06` / `FRN-06B`)

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 3 h
- **Ventana propuesta:** A definir (posterior a `FIX-01` y `FIX-05`).
- **Depende de:** `FIX-01` y `FIX-05`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 5):** el modal no tiene selector de rol, no dispara el `POST` al confirmar el alta, no muestra la contraseña temporal devuelta por el backend y no existe la acción de editar (`PUT /api/admin/users/{id}`).
- **Entregable:** completar el modal de alta con selector de rol y `POST` real mostrando la contraseña temporal, y agregar la acción de edición que llame a `PUT /api/admin/users/{id}` y actualice la tabla.
- **Criterio de éxito:** un administrador puede crear un usuario y ver su contraseña temporal, editar su rol y ver el cambio reflejado sin recargar manualmente la tabla.



### `FIX-08` - Restaurar protección de rutas y guards de sesión (`FRN-03` / `FRN-05`)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 1 h
- **Ventana propuesta:** A definir (alta prioridad).
- **Depende de:** ninguna.
- **Problema (evidencia en `test/front/RESULTADOS.md`, sección pull 20/09):** el commit `6b08566` en `centinela/src/routes/applicationRoutes.tsx` movió `/dashboard`, `/instances`, `/users`, `/users/new`, `/users/:userId` y `/auditoria` al layout público (`MainLayoutAuth`) como "rutas temporales de diseño". Esto anuló el guard `loadProtectedSession`, permitiendo que cualquier usuario acceda a estas vistas sin iniciar sesión ni pasar el 2FA, rompiendo el criterio de éxito de `FRN-03` y `FRN-05` en producto y haciendo fallar `navigation.test.tsx`.
- **Entregable:** reubicar las rutas protegidas (`/dashboard`, `/instances`, `/users`, `/users/new`, `/users/:userId`, `/auditoria`) dentro del grupo con `ProtectedLayout` y `loader: loadProtectedSession`. Si se requieren vistas de diseño sin backend, utilizar mocks dentro del contexto autenticado en lugar de remover la protección de rutas.
- **Criterio de éxito:** un usuario no autenticado que intenta navegar a `/dashboard` o `/users` es redirigido a `/login`; las pruebas de navegación de `test/front` (`navigation.test.tsx`) validan la protección de rutas correctamente.

### `FIX-03` - El guard de re-enrolamiento 2FA no rechaza reemplazar un secreto ya vinculado (`BAC-10`)

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (posterior a `FIX-01`).
- **Depende de:** `FIX-01`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 6):** el criterio de éxito documentado de `BAC-10` exige que "un usuario ya vinculado no pueda reemplazar su secreto sin pasar por el flujo de restablecimiento", pero el server actual sí permite generar un secreto nuevo sobre una cuenta con `is_2fa_enabled: true`.
- **Entregable:** `GET /api/auth/2fa/qr` debe rechazar la generación de un nuevo secreto si la cuenta ya tiene 2FA habilitado, devolviendo un error que indique que debe usarse el flujo de restablecimiento administrativo (`BAC-13`).
- **Criterio de éxito:** una cuenta con `is_2fa_enabled: true` recibe un error controlado al pedir un nuevo QR; solo puede repetir el enrolamiento después de un reset administrativo.

### `FRN-08` - Interceptor HTTP y manejo de `403 Forbidden`

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2 h
- **Ventana propuesta:** 15/09/2026, 14:00-16:00
- **Depende de:** `BAC-08` y `FRN-04`.
- **Entregable:** interceptor Axios o Fetch para enviar `Authorization: Bearer <JWT>` y mostrar un mensaje amigable ante un `403` sin cerrar la sesión.
- **Criterio de éxito:** las peticiones incluyen el JWT y un operador sin permiso recibe una respuesta visual clara.


### `BAC-10` - Enrolamiento 2FA y generación de QR

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** 16/09/2026, 09:00-11:00
- **Depende de:** `BAC-05` y `LOGIN-03`.
- **Entregable:** `GET /api/auth/2fa/qr`, accesible mediante una sesión restringida, que genere un secreto TOTP único, una URI `otpauth://` y un código QR. También debe devolver el secreto para vinculación manual, sin persistirlo como habilitado antes de la confirmación.
- **Criterio de éxito:** un usuario sin 2FA obtiene un QR y una clave manual compatibles con una aplicación autenticadora; un usuario ya vinculado no puede reemplazar su secreto sin pasar por el flujo de restablecimiento.

### `BAC-11` - Validación y persistencia del secreto TOTP

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** 16/09/2026, 11:00-14:00
- **Depende de:** `BAC-10` y `BAC-05`.
- **Entregable:** `POST /api/auth/2fa/enable` para validar el primer código de seis dígitos, cifrar el secreto con AES-256 usando una clave externa a la base de datos y establecer `is_2fa_enabled: true`. El flujo de login también debe validar los códigos posteriores contra el secreto persistido antes de emitir el JWT de acceso completo.
- **Criterio de éxito:** el secreto nunca se almacena ni se expone en texto plano después del enrolamiento; un código válido habilita el 2FA y permite completar el login, mientras que códigos inválidos, vencidos o reutilizados son rechazados.

### `BAC-13` - Restablecimiento administrativo de 2FA

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** 17/09/2026, 09:00-11:00
- **Depende de:** `BAC-06`, `BAC-08` y `BAC-11`.
- **Entregable:** `POST /api/admin/users/{id}/2fa/reset`, restringido a administradores, que invalide el secreto TOTP, establezca `is_2fa_enabled: false` y revoque las sesiones activas del usuario afectado.
- **Criterio de éxito:** un operador recibe `403 Forbidden`; tras el restablecimiento, los códigos del secreto anterior dejan de funcionar y el usuario debe repetir `BAC-10`, `BAC-11` y `FRN-09` en su siguiente acceso. La suite de pruebas `password_recovery_acceptance_test.go` confirma que responde HTTP 200/204 para administradores y 403 para operadores.

### `BAC-15` - Restablecimiento administrativo de contraseña

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 17/09/2026, 14:00-16:00
- **Depende de:** `BAC-06` y `BAC-12`.
- **Entregable:** `POST /api/admin/users/{id}/reset-password`, restringido a administradores, que genere una contraseña temporal segura, actualice su hash, establezca `must_change_password: true` y revoque las sesiones activas del usuario.
- **Criterio de éxito:** un operador recibe `403 Forbidden`; la clave anterior deja de funcionar y el usuario debe cambiar la nueva contraseña temporal en el siguiente acceso.


### `BAC-16` - Entrega segura de credenciales temporales (Mailer Service RF-13)

- **Área:** Backend / Infraestructura
- **Asignados:** Lisandro y Nico
- **Estimación:** 3 h
- **Estado:** Completada la arquitectura y los criterios de seguridad mediante Mailer Service simulado (`MockEmailService`).
- **Depende de:** `BAC-06`, `BAC-15`.
- **Implementación (Cumple con DoD y Seguridad):**
  - **Desacoplamiento (Arquitectura Hexagonal):** Se creó el puerto `EmailService` y se inyectó en los servicios correspondientes (`AuthService` y `UserService`).
  - **Mínimo Privilegio (Cero Confianza):** Se modificaron los endpoints `POST /api/admin/users` y `POST /api/admin/users/{id}/password/reset`. Las contraseñas temporales ya no se retornan en los payloads JSON de respuesta pública.
  - **Sanitización de Logs:** Ningún log del backend expone la contraseña en texto plano, a excepción del entorno controlado del simulador de correo (`MockEmailService`) que actúa como bandeja de entrada de consola.
  - **Manejo de Fallos Estructurado:** Se implementó el control de errores. Si el envío simulado falla, el sistema aborta la transacción antes de persistir inconsistencias y devuelve un `502 Bad Gateway` con el `errorCode: "EMAIL_DELIVERY_FAILED"`, sin exponer trazas internas.
- **Decisión de alcance:** La integración con servidor SMTP real queda formalmente descartada del alcance. La solución definitiva del proyecto es `MockEmailService` por consola mediante el puerto `EmailService`.
- **Criterio de éxito:** el usuario recibe su credencial de forma segura por el puerto de correo sin exposición en JSON ni persistencia en texto plano; las suites de aceptación confirman hash seguro y control estructurado.

### `BAC-17` - Logout y revocación de sesión/JWT

> [!NOTE]
> La revocación atómica de access y refresh token se completó en una segunda iteración: `BAC-17` - *Cierre de sesión y revocación atómica*, al final de este archivo.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (posterior a `BAC-05`).
- **Depende de:** `BAC-03` y `BAC-05`.
- **Entregable:** `POST /api/auth/logout` y un mecanismo de invalidación de sesión (tabla `sessions` o lista de revocación con TTL) que puedan reutilizar `BAC-13` y `BAC-15` para revocar sesiones activas al resetear 2FA o contraseña.
- **Criterio de éxito:** un JWT revocado deja de autorizar peticiones aunque no haya expirado por tiempo; `BAC-13` y `BAC-15` consumen este mecanismo en lugar de simular la revocación.

### `BAC-18` - Base transversal de auditoría (append-only)

> [!WARNING]
> **Estado: Implementado con retroceso detectado (23/09/2026).** El registro, la consulta, los filtros, la exportación CSV y el 403 al operador funcionan. Pero la base de datos **permite `UPDATE` y `DELETE` sobre `auditoria`** con las credenciales de la aplicación: no hay trigger, `REVOKE` ni regla, así que no se cumple el criterio "falla a nivel de base de datos". Prueba: `test/back/password_recovery_acceptance_test.go`, caso `BAC-18…`. Defecto registrado en [`futuro.md`](futuro.md) como `FIX-23`.

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 4 h
- **Ventana propuesta:** A definir (junto con `BAC-05`).
- **Depende de:** `BAC-05`.
- **Entregable:** tabla `audit_logs` con columnas genéricas y reutilizables por cualquier etapa futura: `user_id`, `timestamp`, `accion`, `resource_type`, `resource_id`, `upid` (nullable, para cuando la acción dispare una tarea de Proxmox), `resultado` y `detalle` (JSON). La tabla debe crearse **append-only**: el rol de aplicación no debe tener permisos `UPDATE`/`DELETE` sobre ella (a nivel de motor de base de datos, no solo por convención de código). El servicio de auditoría se implementa como un middleware/interceptor central de la capa de servicios, no como llamadas sueltas repetidas en cada handler, para que la Etapa 1 (energía de instancias), la Etapa 2 (aprovisionamiento) y la Etapa 3 (snapshots) lo reutilicen sin tocar el esquema. En la fase base debe registrar, sin exponer secretos, la creación y eliminación de usuarios, los cambios de rol y de permisos por instancia, y los resets de contraseña y 2FA (`BAC-06`, `BAC-06B`, `BAC-07`, `BAC-13`, `BAC-15`).
- **Criterio de éxito:** cada acción administrativa de la fase base queda registrada desde que ocurre; un intento de `UPDATE` o `DELETE` sobre `audit_logs` con las credenciales de la aplicación falla a nivel de base de datos; una acción nueva agregada en una etapa posterior (por ejemplo, `start` de una VM) se audita sin migrar la tabla. El endpoint de consulta con filtros y exportación queda fuera de esta tarea; corresponde a `RF-08` en la Etapa 3.

### `BAC-12` - Cambio obligatorio de contraseña temporal

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Depende de:** `BAC-02`, `BAC-05` y `BAC-06`.
- **Entregable:** `PUT /api/account/password`, accesible con una sesión restringida, que compruebe la contraseña temporal, valide la nueva clave, actualice su hash e indique `cambio_contrasena: false`. Mientras el indicador sea verdadero, el resto de endpoints protegidos permanece bloqueado con HTTP 403 `PASSWORD_CHANGE_REQUIRED`.
- **Criterio de éxito:** la contraseña temporal deja de ser válida después del cambio, la nueva contraseña nunca se guarda en texto plano y el usuario no obtiene acceso completo antes de finalizar el proceso. Verificado 100% en `TestTareasBackendEIntegracion`.

### `BAC-19` - Solicitud de recuperación de contraseña (RF-13)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Depende de:** `BAC-05` y `BAC-16`.
- **Entregable:** `POST /api/auth/password/forgot`, que valide el formato del correo, genere un código temporal de seis dígitos con expiración de 15 minutos, invalide códigos previos de la misma cuenta y lo envíe mediante el puerto `EmailService` (`MockEmailService`).
- **Criterio de éxito:** solicitar un nuevo código invalida el anterior; el endpoint responde de forma genérica exista o no la cuenta para no filtrar información de usuarios registrados. Verificado 100% en `TestHitoRecuperacionDeContrasenasYNotificaciones`.

### `BAC-20` - Confirmación de recuperación de contraseña (RF-13)

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 2 h
- **Depende de:** `BAC-19` y `BAC-17`.
- **Entregable:** `POST /api/auth/password/reset`, que valide el código de seis dígitos y su expiración, actualice el hash de la nueva contraseña y revoque las sesiones activas de la cuenta en `sesiones_activas`.
- **Criterio de éxito:** un código vencido o ya usado se rechaza con error controlado; un código válido cambia la contraseña y cierra las sesiones anteriores. Verificado 100% en `TestHitoRecuperacionDeContrasenasYNotificaciones`.

### `BAC-21` - Contrato base del canal de eventos/notificaciones (RF-11)

- **Área:** Backend / Frontend
- **Asignados:** Tayra y Lisandro
- **Estimación:** 2 h
- **Depende de:** `BAC-05`.
- **Entregable:** definición del esquema genérico de evento (`type`, `severity`, `resource_type`, `resource_id`, `message`, `timestamp`, `payload`) compartido entre backend y frontend. En el backend quedó el struct `RealtimeEvent` en `internal/core/ports/event_port.go` con constantes y validador `NewRealtimeEvent`; en el frontend quedó la interfaz en `centinela/src/types/notifications.ts`.
- **Criterio de éxito:** existe un tipo/interfaz compartida y sincronizada entre Go y TypeScript. Verificado 100% en `TestHitoRecuperacionDeContrasenasYNotificaciones`.

### `BAC-07` - Asignación de permisos por recurso

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 3 h
- **Depende de:** `BAC-05` y `BAC-06`.
- **Entregable:** `PUT /api/admin/users/{id}/permissions` y `GET /api/admin/users/{id}/permissions`, con actualización atómica de `user_instances`.
- **Criterio de éxito:** la relación usuario-instancia se persiste, reemplaza el conjunto anterior de forma atómica y solo puede gestionarla un administrador autorizado.

### `BAC-08` - Middleware de autorización por recurso

> [!NOTE]
> **Estado: Completada (verificado el 23/09/2026).** El defecto de rutas no montadas se resolvió con `FIX-16`, ver más abajo. `test/back/resource_access_acceptance_test.go` verifica:
> - `403 INSTANCE_ACCESS_DENIED` en `GET`, `start` y `stop` sobre instancias no asignadas;
> - que Proxmox no recibe la orden;
> - que las instancias asignadas y el ADMIN pasan el guard.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 3 h
- **Depende de:** `BAC-07`.
- **Entregable:** guard que valide `(user_id, instance_id)` en `user_instances` antes de consultar o enviar órdenes a Proxmox VE.
- **Criterio de éxito:** un operador sin permiso recibe `403` y la API no realiza ninguna llamada a Proxmox VE.

### `FRN-07` - Selector de asignación de instancias

> [!NOTE]
> **Estado: Completada (verificado el 23/09/2026).** El selector se conectó a la API con `FIX-14`, ver más abajo. `detailsUserPage.tsx`, `rolesAndPermissions.tsx` y `useUserInstanceAccess.ts`:
> - consultan `GET /api/instances` con Bearer;
> - marcan las instancias asignadas;
> - envían `PUT /api/admin/users/:id/permissions` con `{ vmids }`.
>
> El guardado del perfil junto con los permisos se completó con `FIX-25` (25/09/2026).

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 4 h
- **Depende de:** `FRN-05`, `BAC-07` y `BAC-14`.
- **Entregable:** selector con instancias reales obtenidas de `GET /api/instances`, instancias asignadas y botón que envíe el array de IDs a `PUT /api/admin/users/{id}/permissions`.
- **Criterio de éxito:** el administrador puede guardar una matriz de permisos y verla nuevamente al abrir el usuario.

### `FIX-12` - Falta guard de rol administrativo en rutas de gestión de usuarios (`FRN-05`)

- **Área:** Frontend
- **Asignada:** Belinda / Cristian
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir (alta prioridad).
- **Depende de:** `FIX-08`.
- **Problema:** las rutas administrativas `/users`, `/users/new` y `/users/:userId` deben verificar el rol del usuario autenticado, no solamente la existencia de un JWT. Sin este guard, un usuario `OPERATOR` puede acceder al panel de administración, incumpliendo `RF-09` y el criterio de éxito de `FRN-05`.
- **Entregable:** enlazar `loadAdminSession` a las rutas administrativas o agruparlas bajo un layout con dicho guard.
- **Criterio de éxito:** un usuario `OPERATOR` que navega a `/users` es redirigido a `/dashboard`; un usuario `ADMIN` puede acceder normalmente al panel.

### `FIX-13` - Desalineación de endpoints y PR pendiente en recuperación de contraseñas (`BAC-19` / `BAC-20`)

- **Área:** Backend
- **Asignados:** Lisandro y Tayra
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** la rama o implementación que contenga el flujo de recuperación.
- **Problema:** el informe histórico registró endpoints de recuperación ausentes o con rutas diferentes a las definidas en el contrato, lo que provocaba respuestas `404` y requería integrar la implementación correcta en `main`.
- **Entregable:** confirmar en `main` las rutas oficiales de recuperación, alinear handlers, documentación Swagger y pruebas, y eliminar alias o rutas obsoletas que generen ambigüedad.
- **Criterio de éxito:** las rutas documentadas para solicitar y confirmar la recuperación responden según contrato en `main`, y las pruebas de `BAC-19` y `BAC-20` pasan sin depender de una rama remota no integrada.

### `FIX-08` - Auditoría del contrato de códigos de error

- **Área:** Frontend / Backend
- **Asignados:** Cristian y Lisandro
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** ninguna.
- **Problema:** el backend y el frontend deben compartir una lista única de valores `errorCode`; el informe de pruebas registró una posible desalineación entre los códigos emitidos y los mensajes manejados por el cliente.
- **Entregable:** inventario de todos los `errorCode` emitidos por el backend, su estado HTTP y su manejo explícito en frontend. Corregir las diferencias verificadas y documentar los códigos que deliberadamente usan mensaje genérico.
- **Criterio de éxito:** cada código del contrato tiene una respuesta y un mensaje de frontend verificables, sin depender de nombres históricos como `body.code`.


---

# Verificación del 23/09/2026: tareas movidas desde `actual.md`

Evidencia completa en [test/informe.md](../test/informe.md). Revisiones probadas: backend `15032da` y frontend `3192cf4` (`origin/main` en ambos casos).

**Hito: Control de Acceso Basado en Recursos (completado).** El administrador abre un usuario en `FRN-07`, lista las instancias reales (`BAC-14`), asigna un subconjunto (`BAC-07`). El operador solo ve esas instancias en `GET /api/instances`, recibe `403 INSTANCE_ACCESS_DENIED` sobre las demás (`BAC-08`/`FIX-16`), y el frontend lo muestra en un toast sin cerrar la sesión (`FRN-08`). Verificado de punta a punta en `test/back/login04_acceptance_test.go`, pasos 6.2 a 6.4, en `resource_access_acceptance_test.go` y en `test/front/admin-users.test.tsx`.

### `BAC-14` - Lectura mínima del inventario de Proxmox

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 3 h
- **Depende de:** `BAC-08` y de las credenciales de lectura de Proxmox VE.
- **Entregable:** `GET /api/instances` consumiendo `/cluster/resources` o `/nodes/{node}/resources`, con una respuesta normalizada mínima que incluya ID, nombre, tipo, nodo y estado. Un administrador recibe todo el inventario y un operador solo las instancias asignadas.
- **Criterio de éxito:** `FRN-07` puede cargar IDs reales de VMs y LXC; una cuenta no puede descubrir instancias fuera de su alcance y los errores de Proxmox se traducen a una respuesta HTTP controlada.
- **Verificación:** ✅ `resource_access_acceptance_test.go`, caso `BAC-14…`, contra un stub de Proxmox VE. Verifica:
  - 401 sin token.
  - El ADMIN recibe exactamente las VMs y LXC, sin nodos ni storages, con `{ id, name, type: vm|lxc, node, status }`.
  - El operador ve solo lo asignado y existente.
  - `404 INSTANCE_NOT_FOUND` para un VMID inexistente.
  - Implementación: `instance_handler.go` y `proxmox/client.go`.

### `FIX-16` - Montar middleware de autorización por recurso en rutas de instancias (`BAC-08`)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Depende de:** `BAC-08` y `BAC-14`.
- **Entregable:** habilitar el grupo `/instances` bajo `RequireAuth()` y aplicar `RequireInstanceAccess(instanceRepo, "vmid")` en `GET /instances/:vmid`, `POST /instances/:vmid/start` y `POST /instances/:vmid/stop`, cortando con `403 INSTANCE_ACCESS_DENIED` antes de llamar a Proxmox.
- **Criterio de éxito:** un operador que intente acceder a un VMID no asignado recibe `403 Forbidden` con `INSTANCE_ACCESS_DENIED`.
- **Verificación:** ✅ `main.go:211-216`. En `resource_access_acceptance_test.go`, caso `FIX-16…`, el guard corta GET, start y stop, y el stub de Proxmox no registra la orden.

### `BAC-17` - Cierre de sesión y revocación atómica de sesiones (Backend)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Depende de:** `BAC-03` y `BAC-04`.
- **Entregable:** `POST /api/auth/logout` bajo `RequireAuth`, que obtiene el JTI del access token y revoca atómicamente `jtiRefresh` y `jtiAccess` en `sesiones_activas`. Tras el logout, el access token recibe `401 TOKEN_REVOKED`, y la acción se audita como `LOGOUT`.
- **Criterio de éxito:** tras el logout, ambos tokens quedan con `activa: false`; cualquier petición con el access token revocado es rechazada con 401; la auditoría registra el evento.
- **Verificación:** ✅ `session_security_acceptance_test.go`, caso `BAC-17…`, con tokens reales del circuito login → 2FA. Verifica:
  - 401 sin Bearer.
  - 204 con Bearer.
  - `TOKEN_REVOKED` posterior al logout.
  - El refresh posterior falla.
  - Ambas sesiones inactivas en la base.
  - Un registro `LOGOUT` nuevo.
- **Integración:** ✅ resuelta con `SEC-01` (25/09/2026). El logout desde la interfaz (Bearer + cookie `centinela_refresh`) revoca ambas sesiones.

### `FRN-13` - Flujo integral de logout y limpieza de sesión en cliente (Frontend)

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 2 h
- **Depende de:** `BAC-17`.
- **Entregable:**
  1. `logoutSession()` envía `Authorization: Bearer <accessToken>`.
  2. La sesión local se limpia siempre, aunque la llamada falle.
  3. Hay un listener global de `centinela:api-unauthorized` / `TOKEN_REVOKED` que limpia la sesión y redirige a `/login`.
- **Criterio de éxito:** "Cerrar sesión" revoca en el backend, limpia las credenciales y redirige a `/login` con `replace: true`; un 401 remoto expulsa al usuario.
- **Verificación:** ✅ `test/front/session-security.test.ts`, bloque FRN-13 (6/6). Verifica:
  - Bearer y `credentials: include` en el logout.
  - Limpieza ante caída de red, 500 y 401.
  - `historyAction: REPLACE`.
  - Listener de `TOKEN_REVOKED` (`ApiResponseNotifier.tsx`).
- **Integración:** ✅ resuelta con `SEC-01` (25/09/2026).

### `SEC-02` - Cliente frontend compatible con refresh token HttpOnly

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2 h
- **Depende de:** `SEC-01` y `BAC-17`.
- **Entregable:**
  1. Eliminar el almacenamiento, la lectura y el tipado de `refreshToken`.
  2. Mantener `credentials: 'include'`.
  3. Aceptar respuestas sin `refreshToken`.
  4. Renovar la sesión sin enviar el token en el body.
  5. Conservar solo el `accessToken` en `sessionStorage`.
- **Criterio de éxito:** JavaScript no puede leer el refresh token desde el almacenamiento, la memoria de la aplicación ni las respuestas HTTP; la renovación y el logout funcionan mediante cookies.
- **Verificación:** ✅ `session-security.test.ts`, bloque SEC-02 (4/4); commits `08d44ee` a `390b4b8` en `origin/main`. Verifica:
  - `tokenStorage` sin `getRefreshToken`.
  - No se persiste el refresh token que llegue en el body.
  - El interceptor renueva con `POST /auth/refresh` sin el token y con `credentials: include`.
  - El logout no envía el token.
- **Integración:** ✅ resuelta con `SEC-01` (25/09/2026). El caso `SEC-01 SEC-02 integracion…` de `test/back` pasa.

### `FIX-14` - Integración frontend del selector de instancias (`FRN-07`)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 4 h
- **Depende de:** `BAC-07` y `BAC-14`.
- **Entregable:**
  1. Carga de `GET /api/instances` con Bearer.
  2. Marcado de los VMIDs asignados (`GET /api/admin/users/:id/permissions`).
  3. Envío de `PUT /api/admin/users/:id/permissions` con `{ vmids: number[] }`, con estados de carga y toast.
- **Criterio de éxito:** las pruebas de `FRN-07` en `test/front/admin-users.test.tsx` pasan al 100%.
- **Verificación:** ✅ `admin-users.test.tsx`, bloque FIX-14 / FRN-07 (4/4, con el contrato `{ permisos: [{ vmid, nivelAcceso }] }`). La parte de guardar "junto a la edición del perfil" se completó con `FIX-25`.

### `FRN-10` - Cambio obligatorio de contraseña temporal

> [!WARNING]
> **Estado: Implementado con problema (23/09/2026).** La vista `/change-password` funciona: se deriva ante `403 PASSWORD_CHANGE_REQUIRED`, se bloquea la navegación, se envía `PUT /api/account/password`, se muestran los errores del backend y se valida la longitud. Pero **no ofrece la opción de cerrar sesión**, y el criterio exige que la cuenta "solo pueda cerrar sesión o cambiarla". Prueba: `test/front/password-change.test.tsx` (6/7). Defecto registrado en [`futuro.md`](futuro.md) como `FIX-21`. La validación de complejidad ya figura en `FIX-20`.

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2 h
- **Depende de:** `BAC-12` y `FRN-04`.
- **Entregable:** vista de nueva contraseña y confirmación que detecte `must_change_password`, bloquee la navegación general y llame al endpoint de cambio. El endpoint real es `PUT /api/account/password`; la versión anterior de esta tarea citaba `POST /api/auth/change-password`, que no existe.
- **Criterio de éxito:** una cuenta con contraseña temporal solo puede cerrar sesión o cambiarla; después del cambio continúa al enrolamiento o validación 2FA que corresponda.

### `FRN-14` - Suite de pruebas unitarias y de integración para la vista de Auditoría (`Auditoria.tsx`)

> [!WARNING]
> **Estado: Implementado con problema (23/09/2026).** La suite `test/front/audit.test.tsx` existe y verifica:
> - la carga paginada con Bearer;
> - los registros reales;
> - los filtros de acción, resultado y fechas;
> - la paginación;
> - la exportación CSV con Bearer.
>
> **Falla el caso 5:** en el código commiteado, `/auditoria` está bajo `loadProtectedSession` y no bajo `loadAdminSession` (`applicationRoutes.tsx`), así que **un OPERATOR puede entrar a la vista**. Resultado: 6/7. Defecto registrado en [`futuro.md`](futuro.md) como `FIX-22`.

- **Área:** Frontend
- **Asignada:** Belinda / Luz
- **Estimación:** 2 h
- **Depende de:** `BAC-18` y `FRN-03`.
- **Entregable:** suite `test/front/audit.test.tsx` que verifique:
  1. `GET /api/admin/audit?pagina=1&tamano=10` con Bearer.
  2. Los filtros y la paginación.
  3. La exportación `/api/admin/audit/export?formato=csv`.
  4. La redirección de un operador desde `/auditoria` a `/dashboard`.
- **Criterio de éxito:** `npm test` en `test/front` aprueba los casos de `audit.test.tsx`.

---

# Verificación del 25/09/2026: tareas movidas desde `actual.md`

Evidencia completa en [test/informe.md](../test/informe.md). Revisiones probadas: backend `eb0c9af` y frontend `7fbf969`, el último commit de `main` en ambos submódulos.

### `SEC-01` - Refresh token en cookie HttpOnly en backend

> [!NOTE]
> **Estado: Completada (verificado el 25/09/2026, backend `eb0c9af`, commits `c719c9b` y `47bb8f2`).** `test/back/session_security_acceptance_test.go`: los casos `SEC-01…` e integración `SEC-01 SEC-02…` pasan. Se verifica:
> - `POST /auth/2fa/verify` emite `centinela_refresh` con `HttpOnly`, `Secure` por defecto (configurable con `COOKIE_SECURE`), `SameSite=Strict` y `Path=/api/auth`.
> - El refresh token no aparece en el JSON (`TokenResult.RefreshToken` con `json:"-"`) ni en los logs.
> - `/auth/refresh` renueva solo con la cookie (y la rota) y responde 401 ante una cookie inválida.
> - `/auth/logout` revoca y borra la cookie.
> - Las peticiones reales del frontend (body `{}` con `credentials: 'include'`) funcionan. Esto cierra la integración pendiente de `BAC-17`, `FRN-13` y `SEC-02`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** A definir.
- **Depende de:** `BAC-17`.
- **Entregable:**
	1. Emitir el `refreshToken` mediante `Set-Cookie` al completar `POST /api/auth/2fa/verify`.
	2. Leerlo desde la cookie en `POST /api/auth/refresh` y `POST /api/auth/logout`.
	3. Dejar de devolver el refresh token en las respuestas JSON públicas.
	4. Configurar `HttpOnly`, `Secure`, `SameSite` y `Path` de forma segura y configurable para desarrollo y producción.
	5. Eliminar la cookie al cerrar sesión o cuando la sesión sea inválida.
- **Criterio de éxito:** el backend renueva y revoca sesiones usando exclusivamente la cookie HttpOnly; el refresh token no aparece en cuerpos JSON ni en logs; una cookie inválida o vencida produce un error controlado.

### `FIX-17` - Alinear y formalizar endpoints y contratos de contraseñas (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 25/09/2026, backend `eb0c9af`, commit `5407599`).** `test/back/password_recovery_acceptance_test.go`, caso `FIX-17 contratos canonicos…`, verifica:
> - `docs/swagger.json` documenta `PUT /account/password`, `POST /auth/password/forgot`, `POST /auth/password/reset` y `POST /admin/users/{id}/password/reset`, sin ninguna ruta histórica.
> - `/auth/change-password` y `/admin/users/:id/reset-password` responden 404.
> - Los códigos de error son `PASSWORD_CHANGE_FAILED` (400), `INVALID_REQUEST` (400) y `RESET_FAILED` (400).
> - La recuperación pública deja `cambio_contrasena=false` y revoca todas las sesiones.

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

### `FIX-25` - Persistir la edición del perfil de usuario (`FRN-06B` / `FIX-14`) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 25/09/2026, frontend `7fbf969`, commits `f2d6d9f`, `8164670` y `60f81e5`).** `test/front/admin-users.test.tsx`: los bloques FRN-06B y FIX-25 pasan (3/3), y los 4 casos de FIX-14 siguen en verde. Se verifica:
> - `userDetailsService.updateUserDetails` envía `PUT /api/admin/users/:id` solo con los campos modificados.
> - "Guardar cambios" se habilita con cambios de perfil o de permisos.
> - El `409 USER_CONFLICT` se muestra junto al campo de correo.

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

### `SEC-04` - Esquema extensible de niveles de acceso a recursos en Backend y Modelo de Datos

> [!WARNING]
> **Estado: Implementado con problema (verificado el 25/09/2026, backend `eb0c9af`).**
>
> **Funciona:**
> - La columna `nivel_acceso` tiene default `FULL_ACCESS` y un CHECK con `FULL_ACCESS` y `READ_ONLY`.
> - El contrato `{ permisos: [{ vmid, nivelAcceso }] }` funciona en `PUT` y `GET /api/admin/users/:id/permissions`.
> - `VerificarAcceso` y `RequireInstanceAccess` reciben el nivel requerido: con `READ_ONLY` se puede hacer `GET /instances/:vmid`, pero `start` y `stop` responden 403; con `FULL_ACCESS` se puede hacer `start`.
>
> **Falla el entregable 4 (retrocompatibilidad):** `PUT /permissions` con el payload anterior `{ "vmids": [103] }` responde `400 INVALID_REQUEST`, porque `asignarPermisosRequest` solo acepta `permisos`.
>
> Pruebas: `test/back/resource_access_acceptance_test.go`. El caso `SEC-04 niveles de acceso…` pasa; el caso `FIX-26 SEC-04 retrocompatibilidad…` falla. Defecto registrado en [`futuro.md`](futuro.md) como `FIX-26`.

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

### `FIX-24` - Confirmación del alta sin contraseña temporal y baja de usuario (`FRN-06` / `BAC-16`) (Frontend)

> [!WARNING]
> **Estado: Implementado con problema (verificado el 25/09/2026, frontend `7fbf969`).**
>
> **Funciona:**
> - El alta pide confirmar el correo y confirmar en un modal (`ConfirmUserAction`).
> - Muestra el toast *"Usuario creado exitosamente… enviada por correo"*.
> - El formulario ya no tiene campos de contraseña.
> - "Eliminar usuario" abre un modal y envía `DELETE /api/admin/users/:id`.
>
> **Falla el entregable 2:** ante `502 EMAIL_DELIVERY_FAILED`, `createUserService.ts` arma el mensaje correcto, pero `useCreateUser.ts` lo descarta y muestra uno genérico.
>
> Prueba: `test/front/admin-users.test.tsx`, bloque FRN-06 (5/6); falla el caso `FIX-27…`. Defecto registrado en [`futuro.md`](futuro.md) como `FIX-27`.

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
