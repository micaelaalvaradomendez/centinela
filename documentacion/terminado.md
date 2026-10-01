

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
> **Estado: Implementado con problema (verificado el 28/09/2026, frontend `8d7ecab`).** El retroceso del 23/09 se resolvió con `FIX-24`. `test/front/admin-users.test.tsx`, bloque FRN-06 (6/7):
> - **Funciona:** selector de rol, `POST /api/admin/users` con el contrato real, toast de alta con la clave enviada por correo (sin mostrar contraseña), formulario sin campos de contraseña y "Eliminar usuario" con confirmación y `DELETE`.
> - **Falla** el mensaje ante `502 EMAIL_DELIVERY_FAILED`: se muestra uno genérico. Corrección registrada en [`futuro.md`](futuro.md) como `FIX-27`.
>
> **Nota:** el entregable original pedía una "visualización controlada de la contraseña temporal". Desde `BAC-16` la clave solo se envía por correo, y así lo verifica la prueba.

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

> [!NOTE]
> **Estado: Completada según el contrato vigente (verificado el 28/09/2026).** El selector de rol, el `POST` de alta y la edición con `PUT /api/admin/users/{id}` funcionan (`admin-users.test.tsx`, bloques FRN-06 y FRN-06B). El criterio original pedía *"ver su contraseña temporal"*, pero desde `BAC-16` la clave solo se envía por correo y no se muestra.

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

> [!NOTE]
> **Actualización (01/10/2026):** el adaptador SMTP real quedó hecho con `BAC-16B` (backend `860b3c9`, ver la verificación del 01/10/2026 en este archivo).

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

> [!NOTE]
> **Estado: Completada (verificado el 29/09/2026, backend `9554efa`).** El problema detectado el 23/09 (la base permitía `UPDATE` y `DELETE` sobre `auditoria`) se resolvió con `FIX-23`: trigger `trg_auditoria_inmutable` y `REVOKE`. `test/back/password_recovery_acceptance_test.go`, caso `BAC-18 … append-only`, pasa: registro, consulta, filtros, CSV, 403 al operador, y `UPDATE`, `DELETE` y `TRUNCATE` rechazados por el motor.

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

> [!NOTE]
> **Actualización (01/10/2026):** la regresión se corrigió con `FIX-28` (commit `cfb88f7`, ver la verificación del 01/10/2026 en este archivo). El logout revoca la sesión en el servidor y la prueba LOGIN-04 pasa 10/10.

> [!WARNING]
> **Regresión detectada (28/09/2026, frontend `8d7ecab`).** El commit `deb59cb` quitó el `Authorization: Bearer` del logout (`skipAuthorization: true`). El backend responde `401 MISSING_TOKEN` y la sesión queda activa en el servidor, aunque el cliente limpia su almacenamiento y redirige. La corrección está registrada en [`futuro.md`](futuro.md) como `FIX-28`.

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
> **Estado: Implementado con problema (verificado el 28/09/2026, frontend `8d7ecab`).** `test/front/password-change.test.tsx` (9/10):
> - **Funciona:** se deriva ante `403 PASSWORD_CHANGE_REQUIRED`, se bloquea la navegación, se envía `PUT /api/account/password`, se muestran los errores del backend, se valida el largo, el dígito y el símbolo, y la opción "Cerrar sesión" (`FIX-21`, terminada) está disponible.
> - **Falla:** la validación de mayúscula del cliente (`validatePasswordComplexity` no la verifica). Corrección registrada en [`futuro.md`](futuro.md) como `FIX-29`.

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2 h
- **Depende de:** `BAC-12` y `FRN-04`.
- **Entregable:** vista de nueva contraseña y confirmación que detecte `must_change_password`, bloquee la navegación general y llame al endpoint de cambio. El endpoint real es `PUT /api/account/password`; la versión anterior de esta tarea citaba `POST /api/auth/change-password`, que no existe.
- **Criterio de éxito:** una cuenta con contraseña temporal solo puede cerrar sesión o cambiarla; después del cambio continúa al enrolamiento o validación 2FA que corresponda.

### `FRN-14` - Suite de pruebas unitarias y de integración para la vista de Auditoría (`Auditoria.tsx`)

> [!WARNING]
> **Regresión detectada (01/10/2026, frontend `3d1e84a`).** El commit `56b88f5` (*"eliminar columnas que no serán utilizadas"*) quitó de `Auditoria.tsx` el filtro **"Acción"** y el parámetro `accion`, que exige RF-08 (`requerimientos.md`: *"filtros por fecha, usuario, tipo de acción y resultado"*). `audit.test.tsx` pasa 6/7; falla *"actualiza los query parameters de filtro…"*. La corrección es **`FIX-36`**, en [`futuro.md`](futuro.md).

> [!NOTE]
> **Estado: Completada (verificado el 29/09/2026, frontend `749194e`).** `test/front/audit.test.tsx` pasa 7/7: carga paginada con Bearer, registros reales, filtros de acción, resultado y fechas, paginación, exportación CSV con Bearer y redirección del OPERATOR a `/dashboard` (el guard se corrigió con `FIX-22`).

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

> [!NOTE]
> **Estado: Completada (verificado el 25/09/2026 y confirmado el 28/09/2026, backend `eb0c9af` / `d36bc50`).**
> - La columna `nivel_acceso` tiene default `FULL_ACCESS` y un CHECK con `FULL_ACCESS` y `READ_ONLY`.
> - El contrato `{ permisos: [{ vmid, nivelAcceso }] }` funciona en `PUT` y `GET /api/admin/users/:id/permissions`; si se omite `nivelAcceso` en un ítem, se asume `FULL_ACCESS`.
> - `VerificarAcceso` y `RequireInstanceAccess` reciben el nivel requerido: con `READ_ONLY` se puede hacer `GET /instances/:vmid`, pero `start` y `stop` responden 403; con `FULL_ACCESS` se puede hacer `start`.
>
> **Decisión (28/09/2026):** el contrato oficial es el formato nuevo `{ permisos }`, que ya usa el frontend (`userInstanceService.ts`). **No se mantiene la retrocompatibilidad con `{ vmids }`:** el entregable 4 original se reemplazó y `FIX-26` quedó descartado. El payload anterior se rechaza con `400 INVALID_REQUEST`.
>
> Pruebas: `test/back/resource_access_acceptance_test.go`, casos `SEC-04 niveles de acceso…` y `SEC-04 el payload anterior { vmids } se rechaza…`.

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
  4. ~~Mantener retrocompatibilidad total: si el payload de `PUT /permissions` solo envía `vmids: [101]`, asignar `FULL_ACCESS` por defecto.~~ **Reemplazado (28/09/2026):** `PUT /permissions` recibe `{ permisos: [{ vmid, nivelAcceso? }] }`. Si se omite `nivelAcceso` en un ítem, se asume `FULL_ACCESS`, y el formato anterior `{ vmids }` se rechaza con `400 INVALID_REQUEST`.
- **Criterio de éxito:** La migración crea el campo sin romper registros previos; el middleware `RequireInstanceAccess` verifica tanto la pertenencia de la instancia como el nivel de permiso; si un usuario tiene permiso `READ_ONLY` sobre la VM 101, puede consultar su estado pero recibe 403 al intentar ejecutar una acción de apagado/encendido.

### `FIX-24` - Confirmación del alta sin contraseña temporal y baja de usuario (`FRN-06` / `BAC-16`) (Frontend)

> [!NOTE]
> **Actualización (01/10/2026):** el entregable 2 se corrigió con `FIX-27` (commit `3ea7fa4`). FRN-06 pasa 6/6.

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

---

# Verificación del 28/09/2026: tareas movidas desde `actual.md`

Evidencia completa en [test/informe.md](../test/informe.md). Revisiones probadas: backend `d36bc50` y frontend `8d7ecab`, el último commit de `main` en ambos submódulos.

### `FRN-11` - Acciones administrativas de recuperación

> [!NOTE]
> **Estado: Completada (verificado el 28/09/2026, frontend `8d7ecab` / backend `d36bc50`).** `test/front/admin-recovery.test.tsx` (4/4, commit `fe158ab`). Se verifica:
> - "Restablecer contraseña" y "Restablecer 2FA" son acciones separadas en el menú de `Users.tsx`, cada una con su modal de confirmación (`ConfirmUserAction`).
> - No se llama a la API antes de confirmar.
> - Se envía `POST /api/admin/users/:id/password/reset` o `…/2fa/reset` con Bearer y se muestra el toast de resultado.
> - La tabla pasa a "Desactivado".
> - Nunca se muestran secretos ni hashes.

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2 h
- **Ventana propuesta:** 18/09/2026, 12:00-14:00
- **Depende de:** `FRN-05`, `BAC-13` y `BAC-15`.

- **Entregable:** acciones separadas para restablecer contraseña y 2FA desde el panel de usuarios, ambas con confirmación explícita, estado de carga y notificación del resultado.
- **Criterio de éxito:** un administrador puede iniciar cada recuperación sin confundir sus efectos; la tabla refleja que el 2FA quedó desvinculado y nunca muestra secretos ni hashes.

### `FRN-12` - Vistas de recuperación de contraseña (RF-13)

> [!NOTE]
> **Estado: Completada (verificado el 28/09/2026, frontend `8d7ecab` / backend `d36bc50`).** `test/front/recover-password.test.tsx` (commit `fe6d784`). Se verifica:
> - El paso 1 envía `POST /api/auth/password/forgot` con `{ email }` y no avanza con un correo inválido.
> - El paso 2 exige el código de 6 dígitos.
> - El paso 3 envía `POST /api/auth/password/reset` con `{ email, codigo, nuevaContrasena }`.
> - Muestra `RESET_FAILED` como "código incorrecto o expirado" y redirige a `/login`.
>
> El único caso que falla del archivo (complejidad en el paso 3) corresponde a `FIX-20`.

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 3 h
- **Ventana propuesta:** A definir (posterior a `BAC-20`).
- **Depende de:** `BAC-19`, `BAC-20` y el maquetado existente de `RecoverPassword.tsx`.

- **Entregable:** conectar `RecoverPassword.tsx` al flujo real: paso de ingreso de correo, paso de ingreso del código de seis dígitos y paso de nueva contraseña, con manejo de errores del backend en cada paso.
- **Criterio de éxito:** una cuenta puede recuperar el acceso sin intervención de un administrador, y los errores de código inválido o vencido se muestran junto al campo correspondiente.

### `FIX-18` - Conexión de `RecoverPassword.tsx` a la API y alineación de tests (`FRN-12` / RF-13) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 28/09/2026, frontend `8d7ecab` / backend `d36bc50`).** Ídem `FRN-12`: `test/front/recover-password.test.tsx` (6/6 de sus casos, commit `fe6d784`), incluidas la redirección a `/login` y la validación del código en el paso 2.

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

> [!NOTE]
> **Estado: Completada (verificado el 28/09/2026, frontend `8d7ecab` / backend `d36bc50`).** Ídem `FRN-11`: `test/front/admin-recovery.test.tsx` (4/4, commit `fe158ab`).

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

> [!NOTE]
> **Estado: Completada (verificado el 28/09/2026, frontend `8d7ecab` / backend `d36bc50`).** `test/front/password-change.test.tsx` (commit `b87ea1b`): desde `/change-password`, "Cerrar sesión" usa `useLogout()`, llama a `/auth/logout`, borra `centinela_access` y navega a `/login`. **Nota:** la revocación en el servidor está afectada por la regresión de `FRN-13` (logout sin Bearer, ver `FIX-28` en `futuro.md`).

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



### `FIX-20` - Sincronización de políticas de complejidad y distinción UX entre Cambio y Restablecimiento (`FRN-10` / `FRN-12`) (Frontend / UX)

> [!WARNING]
> **Estado: Implementado con problema (verificado el 28/09/2026, frontend `8d7ecab`, commit `fe6d784`).**
>
> **Funciona:**
> - Existe `validatePasswordComplexity` en `components/features/auth/utils/validateAuthenticationFields.ts`, integrado en `ChangePassword.tsx` y en el paso 3 de `RecoverPassword.tsx`.
> - Frena claves fuera de 8-12 caracteres, sin dígito o sin símbolo, y muestra textos de ayuda.
>
> **Falla:** no verifica la mayúscula. La segunda condición prueba `/[0-9]/` con el mensaje "al menos una letra mayúscula", así que una clave como `nueva1234!` pasa el cliente y el backend la rechaza con 400. Además, el mensaje de largo dice "Debe tener 8 y 12 caracteres".
>
> Pruebas que fallan: `password-change.test.tsx` (caso "sin mayúscula") y `recover-password.test.tsx` (paso 3). Corrección registrada en [`futuro.md`](futuro.md) como `FIX-29`.

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

> [!WARNING]
> **Estado: Implementado con problema (verificado el 28/09/2026, frontend `8d7ecab`).**
>
> **Entregable 1 cumple:** `detailsUserPage.tsx`, `rolesAndPermissions.tsx` y `userInstanceService.ts` leen el `nivelAcceso` desde `GET /permissions` (READ_ONLY se muestra como "Solo lectura") y envían `{ vmid, nivelAcceso }` en el `PUT`. Pruebas: los 2 casos `FRN-18…` de `admin-users.test.tsx` pasan.
>
> **Entregable 2 no existe:** no hay `canOperateInstance`, porque el hook `usePermissions()` de `SEC-03` no está implementado. Prueba que falla: `navigation.test.tsx`, caso `FIX-30 canOperateInstance…`. Corrección registrada en [`futuro.md`](futuro.md) como `FIX-30`.
>
> **Actualización (01/10/2026):** `FIX-30` ya agregó `canOperateInstance` en `src/hooks/usePermissions.ts` (ver la verificación del 01/10/2026). Pero el backend no informa el nivel por instancia en `GET /account/profile`, así que con datos reales el helper da `false` para todo OPERATOR. La corrección es **`FIX-37`**, en [`futuro.md`](futuro.md).

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (inmediatamente posterior a `SEC-03` y `SEC-04`).
- **Depende de:** `SEC-03`, `SEC-04` y `FIX-14`.
- **Problema y evidencia (análisis de cierre de la fase base):** `SEC-04` implementa `nivel_acceso` (`FULL_ACCESS` y `READ_ONLY`) en el backend, pero ninguna tarea de Frontend tenía asignado enviar ese nivel en `PUT /api/admin/users/:id/permissions`, leerlo al abrir la ficha del usuario ni exponerlo en `usePermissions()` (`SEC-03`) para distinguir quién puede solo ver una máquina de quién puede apagarla o reiniciarla.
- **Entregable:**
  1. En `detailsUserPage.tsx` y `rolesAndPermissions.tsx`, leer el `nivelAcceso` de cada instancia desde `GET /api/admin/users/:id/permissions` y enviar `{ vmid, nivelAcceso: 'FULL_ACCESS' | 'READ_ONLY' }` en `PUT /api/admin/users/:id/permissions`.
  2. En el contexto `SEC-03` (`usePermissions()`), agregar el helper `canOperateInstance(vmid: number): boolean` (devuelve `true` solo si es `ADMIN` o si tiene `FULL_ACCESS` sobre ese `vmid`), diferenciándolo de `canAccessInstance(vmid)` (que devuelve `true` tanto para `READ_ONLY` como `FULL_ACCESS`).
- **Criterio de éxito:** El administrador puede asignar y guardar el nivel "Solo lectura" o "Control total" por instancia desde la UI; `canOperateInstance` retorna `false` para instancias en modo `READ_ONLY`.

### `INF-05` - CORS y TLS en el borde

> [!WARNING]
> **Estado: Implementado con problema (verificado el 28/09/2026, backend `d36bc50`, commit `c624a57`).**
>
> **CORS cumple:** `middleware.CORS()` aplica la lista blanca de `ALLOWED_ORIGINS`. El origen permitido recibe `204` con `Access-Control-Allow-Origin` exacto y `Allow-Credentials: true`; uno no autorizado recibe `403` sin cabeceras CORS. Prueba: `cierre_fase_base_acceptance_test.go`, caso `INF-05 CORS…`.
>
> **TLS no verificable:** la configuración de Nginx del servidor (CT 103) no está versionada en ningún repositorio, así que no se puede comprobar que el tráfico se sirva solo por HTTPS. Prueba omitida: `FIX-31 INF-05 TLS…`. Corrección registrada en [`futuro.md`](futuro.md) como `FIX-31`.

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (junto con `INF-04`).
- **Depende de:** `INF-04`.
- **Entregable:** whitelist de orígenes permitidos en CORS y certificados TLS configurados en Nginx para todo el tráfico hacia el frontend y la API.
- **Criterio de éxito:** una petición desde un origen no autorizado es rechazada por CORS y el tráfico hacia el sistema se sirve únicamente sobre HTTPS.

---

# Verificación del 29/09/2026: tareas movidas desde `actual.md`

Evidencia completa en [test/informe.md](../test/informe.md). Revisiones probadas: backend `9554efa` y frontend `749194e`, el último commit de `main` en ambos submódulos.

### `FIX-22` - Guard administrativo en la ruta `/auditoria` (`FRN-14` / `BAC-18`) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 29/09/2026, frontend `749194e`, commit `5c3d780`).** `/auditoria` se movió al grupo `loadAdminSession` de `applicationRoutes.tsx`. `test/front/audit.test.tsx`, caso *"si un operador intenta entrar a /auditoria, el guard administrativo lo redirige a /dashboard"*, pasa, y los otros 6 casos de `FRN-14` siguen en verde. **Pendiente relacionado:** el enlace "Auditoría" del menú lateral sigue visible para el OPERATOR; eso lo cubre `SEC-03` en `actual.md`.

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

> [!NOTE]
> **Estado: Completada (verificado el 29/09/2026, backend `9554efa`, commits `ad52485` y `9554efa`).** `postgres/db.go` y `scripts/init.sql` crean la función `audit_inmutabilidad()` y el trigger `trg_auditoria_inmutable` (`BEFORE UPDATE OR DELETE OR TRUNCATE … FOR EACH STATEMENT`). Además, se revocan `UPDATE`, `DELETE` y `TRUNCATE` al usuario de la aplicación, y la FK de `auditoria.usuario_id` pasa a `ON DELETE RESTRICT`. `test/back/password_recovery_acceptance_test.go`, caso `BAC-18 … append-only`, pasa: los tres comandos fallan en PostgreSQL, y el registro, la consulta, los filtros y el CSV siguen funcionando.

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

---

# Verificación del 30/09/2026: tareas movidas desde `actual.md`

Evidencia completa en [test/informe.md](../test/informe.md). Revisiones probadas: backend `43a0b06` y frontend `749194e`.

### `BAC-17A` - Adaptador base de Redis en el backend (conexión, configuración y puerto)

> [!NOTE]
> **Estado: Completada (verificado el 30/09/2026, backend `43a0b06`), commit `0aac034`.** `test/back/cierre_fase_base_acceptance_test.go`, caso `BAC-17A…`:
> - `github.com/redis/go-redis/v9` en `go.mod`, adaptador `internal/adapters/secondary/redis/` y puerto `KeyValueStore` con `GetDel`, `Publish` y `Subscribe`.
> - Con `REDIS_ADDR` configurado, el backend abre la conexión al arrancar. Si Redis no responde, arranca en **modo degradado** (almacén en memoria), según `docs/redis.md`.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** Ya. Es la base de `BAC-17B` y `BAC-21C`.
- **Depende de:** `INF-06A` (Redis local). No depende del servidor (`INF-06B`).
- **Problema y evidencia:** `BAC-17B` (sesiones en Redis con TTL) y `BAC-21C` (tickets de `/api/events` y bus Pub/Sub) dan por hecho que el backend ya habla con Redis. Hoy no es así: `go.mod` no incluye ningún cliente de Redis, no hay adaptador en `internal/adapters/secondary` y `cmd/api/main.go` no lee `REDIS_ADDR` ni `REDIS_PASSWORD`. Si cada tarea arma su propia conexión, se duplica el trabajo y se pisan entre sí.
- **Entregable:**
  1. Agregar `github.com/redis/go-redis/v9` y crear `internal/adapters/secondary/redis/` con un constructor que lea `REDIS_ADDR`, `REDIS_PASSWORD` y `REDIS_DB` (opcional, default `0`).
  2. Definir un puerto en `internal/core/ports`, por ejemplo `KeyValueStore`, con lo mínimo que usan las tareas siguientes:
     - `Set(key, value, ttl)`, `Get`, `GetDel` (atómico) y `Del`, para las sesiones de BAC-17B y los tickets de BAC-21C;
     - `Publish(canal, mensaje)` y `Subscribe(canal)`, para el bus de BAC-21C.
  3. En `main.go`, conectarse al arrancar con un `PING` y un log claro. Hay que definir y documentar qué pasa si Redis no está: que el backend no arranque, o que funcione en modo degradado sin tickets ni bus.
  4. Pruebas unitarias del adaptador con `miniredis` o contra el Redis del compose.
- **Criterio de éxito:** con `docker compose up redis`, el backend en local arranca y confirma la conexión con Redis. `BAC-17B` y `BAC-21C` usan este puerto en lugar de crear sus propias conexiones.

### `BAC-17B` - Corrección de lógica de inserción en `sesiones_activas` (1 sesión = 1 registro) y almacenamiento en Redis con TTL

> [!NOTE]
> **Estado: Completada (verificado el 30/09/2026, backend `43a0b06`), commit `9e0f43c`.** `test/back/cierre_fase_base_acceptance_test.go`, caso `BAC-17B…`:
> - Login + 2FA + 10 refresh generan **1 sola fila** en `sesiones_activas`, con `id` igual al claim `sid` y las columnas `jti_access` y `jti_refresh`.
> - Después del logout no quedan filas activas.
> - Cada refresh rota el `jti_access` y el access token anterior deja de valer (`session_security…`, caso de integración SEC-01/SEC-02).
> - La prueba integral LOGIN-04 sigue en verde con el esquema nuevo.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2,5 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `INF-06A` (Redis local en el repo del backend), `BAC-17A` y `BAC-17`.
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

### `BAC-21C` (`BRG-02-BAC`) - Autenticación de `/api/events` por ticket efímero en Redis, bus Pub/Sub y revocación en vivo

> [!NOTE]
> **Estado: Completada (verificado el 30/09/2026, backend `43a0b06`), commit `45ecad7`.** `test/back/puente_etapa1_acceptance_test.go`, casos `BAC-21C…` (2/2):
> - `POST /api/events/ticket` exige Bearer y guarda `ws_ticket:<uuid>` en Redis con TTL de 30 s o menos.
> - `GET /api/events?ticket=` rechaza con 401 un ticket inválido o reutilizado.
> - El stream SSE se corta al hacer logout.
>
> Con esto queda desbloqueada `FRN-17C` en el frontend.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** Junto a `BAC-25B` y `BAC-26`.
- **Depende de:** `INF-06A` (Redis local en el repo del backend), `BAC-17A` y `BAC-17B`.
- **Problema y evidencia (análisis de cierre de la fase base):** `EventSource` (SSE) y WebSockets no permiten enviar el header `Authorization: Bearer`, y la cookie `HttpOnly` de `SEC-01` solo viaja a `/api/auth`. Además, si se revoca la sesión o un permiso de instancia, las conexiones abiertas deben cerrarse o actualizar su filtro, y el bus de eventos debe usar Redis Pub/Sub para funcionar con más de una réplica.
- **Entregable:**
  1. Crear el endpoint `POST /api/events/ticket` (bajo `RequireAuth`) que guarde en Redis un ticket de un solo uso con TTL de 30 segundos (`SET ws_ticket:<uuid> <usuario_id> EX 30`) y devuelva `{ ticket }`.
  2. En `GET /api/events?ticket=<uuid>`, validar y consumir atómicamente (`GETDEL`) el ticket en Redis antes de abrir el canal SSE/WebSocket.
  3. Publicar los eventos `TASK_FINISHED` de `BAC-25B` mediante **Redis Pub/Sub** (`centinela:events`) y cortar inmediatamente la conexión activa del usuario si recibe un evento de revocación de sesión (`BAC-17`) o recargar su filtro si cambian sus permisos (`BAC-07`).
- **Criterio de éxito:** `/api/events` rechaza con 401 tickets inválidos o reutilizados; al hacer logout o desactivar al usuario, el backend corta el stream inmediatamente.

### `INF-06A` - Redis local para desarrollo en el repositorio del backend

> [!NOTE]
> **Estado: Completada (verificado el 30/09/2026, backend `43a0b06`), commit `546e1d5`.** `test/back/cierre_fase_base_acceptance_test.go`, caso `INF-06A…`:
> - Redis publicado en `127.0.0.1:6379:6379` y `REDIS_ADDR=localhost:6379` en `.env.example`.
> - La misma contraseña en el compose y en `.env.example` (`centinela_redis_pass`).
> - `REDIS_DB` agregado.

- **Área:** Backend (en el repositorio del backend)
- **Asignada:** Tayra
- **Estimación:** 0,5 h
- **Ventana propuesta:** Ya.
- **Depende de:** ninguna.
- **Problema y evidencia:** el servicio `redis` ya está en `backend/docker-compose.yml`, con contraseña, `maxmemory 256mb` y `volatile-lru`, verificado por `test/back`. Pero el backend en local (`go run ./cmd/api`) **no lo puede usar**, por tres motivos:
  - el servicio no publica ningún puerto;
  - `REDIS_ADDR=redis:6379` solo resuelve dentro de Docker;
  - la contraseña por defecto del compose (`centinela_redis_pass`) no coincide con la de `.env.example` (`centinela_redis_password`).
- **Entregable:**
  1. Publicar Redis **solo en loopback** en `backend/docker-compose.yml` (`127.0.0.1:6379:6379`).
  2. En `backend/.env.example`, poner `REDIS_ADDR=localhost:6379` para desarrollo, con `redis:6379` comentado para cuando el backend corre en un contenedor, y unificar la contraseña con la del compose.
- **Criterio de éxito:** con `docker compose up redis`, desde la máquina de desarrollo `redis-cli -h 127.0.0.1 -a <pass> PING` responde `PONG`, y sin contraseña responde `NOAUTH`. La prueba `INF-06A…` de `test/back` pasa.

### `INF-08B` - Credenciales SMTP en `.env.example` y `.env` del repositorio del backend

> [!WARNING]
> **Estado: Implementado con problema (verificado el 30/09/2026, backend `43a0b06`), commit `233e804`.** Las 6 variables (`EMAIL_PROVIDER=smtp`, `SMTP_HOST=smtp-relay.brevo.com`, `SMTP_PORT=587`, `SMTP_USER`, `SMTP_PASS` y `SMTP_FROM`) ya están en `backend/.env.example`. Pero `SMTP_USER` (`tu_correo@ejemplo.com`), `SMTP_PASS` (`tu_clave_secreta_aqui`) y `SMTP_FROM` (`no-reply@tudominio.com`) son **valores de ejemplo**, así que no se puede autenticar contra el relay. Prueba: `cierre_fase_base_acceptance_test.go`, caso `INF-08B…`. La corrección es `FIX-34`, en [`futuro.md`](futuro.md).

- **Área:** Backend (en el repositorio del backend)
- **Asignado:** Lisandro
- **Estimación:** 0,5 h
- **Ventana propuesta:** Ya. El `.env` ya está hecho en local, pero falta el `push`.
- **Depende de:** ninguna.
- **Problema y evidencia:** `BAC-16B` necesita credenciales SMTP reales para desarrollar y probar el envío desde local. Hoy `backend/.env.example` no tiene ninguna de las variables.
- **Entregable:**
  1. Agregar a `backend/.env.example` las variables sin comentar: `EMAIL_PROVIDER=smtp`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` y `SMTP_FROM`.
  2. Tener las mismas variables en el `.env` local de cada desarrollador del backend (el `.env` no se versiona).
  3. Hacer el `push` a `main`.
- **Criterio de éxito:** la prueba `INF-08B…` de `test/back` pasa. La prueba lee `backend/.env.example`, verifica que estén las 6 variables con valores reales (no de ejemplo) y **se autentica contra el servidor SMTP**, sin enviar correos.
- **Advertencia de seguridad:** `.env.example` queda versionado. Si contiene la contraseña real, cualquiera con acceso al repositorio puede enviar correo como El Centinela. Conviene una cuenta de uso exclusivo con permisos mínimos.

---

# Verificación del 01/10/2026: tareas movidas desde `actual.md`

Evidencia completa en [test/informe.md](../test/informe.md). Revisiones probadas: backend `860b3c9` y frontend `3d1e84a`.

### `BAC-16B` - Adaptador Go `SmtpEmailService` para envío real de credenciales y códigos OTP (`RF-09` / `RF-13`)

> [!NOTE]
> **Estado: Completada (verificado el 01/10/2026, backend `860b3c9`).** `test/back/cierre_fase_base_acceptance_test.go`, caso `BAC-16B…`, con Mailpit en **STARTTLS obligatorio**:
> - Con `EMAIL_PROVIDER=smtp`, `cmd/api/main.go` instancia `SMTPEmailService` (`go-mail`; STARTTLS en el 587 y TLS directo en el 465). Sin `EMAIL_PROVIDER=smtp` sigue usando `MockEmailService`, y la suite pasa en ese modo.
> - Correos multiparte (texto y HTML) para alta y reset administrativo (`EnviarCredencialesTemporales`) y para el código OTP (`EnviarCodigoRecuperacion`). La clave temporal que llega permite iniciar sesión, y el código de 6 dígitos también llega.
> - Con el servidor SMTP caído, el alta responde `502 EMAIL_DELIVERY_FAILED` y el usuario no se guarda.
>
> **Notas:**
> - El puerto cambió de nombre: `EnviarContrasenaTemporal` pasó a `EnviarCredencialesTemporales(destinatario, nombre, clave)`, y se agregó `SendMail`.
> - El envío con las credenciales reales del relay depende de `FIX-34` (`INF-08B`).
> - Se verificó contra Mailpit, no contra Brevo.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `BAC-16`, `BAC-19` e `INF-08B` (credenciales SMTP en el repo del backend, para desarrollar en local). No depende de la configuración del servidor (`INF-08A`).
- **Problema y contexto:** `BAC-16` dejó creado el puerto hexagonal `ports.EmailService`, pero solo implementó `MockEmailService` por consola. El backend necesita el adaptador SMTP real para enviar las contraseñas temporales y los códigos de recuperación de 6 dígitos.
- **Entregable:**
  1. Implementar `SmtpEmailService` en `backend/internal/adapters/secondary/email/smtp_service.go` cumpliendo la interfaz `ports.EmailService` (`EnviarContrasenaTemporal` y `EnviarCodigoRecuperacion`, los nombres reales de `internal/core/ports/email_port.go`) con soporte `STARTTLS`/`TLS`.
  2. En `cmd/api/main.go`, instanciar `SmtpEmailService` cuando `EMAIL_PROVIDER=smtp` y mantener `MockEmailService` cuando `EMAIL_PROVIDER=mock` (usado por los tests automatizados).
  3. Estructurar los cuerpos de correo para alta de cuenta (`RF-09`), reset administrativo (`BAC-15`) y código OTP de 6 dígitos (`RF-13`).
- **Criterio de éxito:** Con `EMAIL_PROVIDER=smtp` y las credenciales de `INF-08B`, el backend envía correos reales desde local; ante un fallo de entrega aborta la operación y devuelve `502 EMAIL_DELIVERY_FAILED`; la suite de tests sigue pasando en modo `mock`.

### `FIX-27` - Mensaje específico cuando falla el envío del correo en el alta (`FIX-24`) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 01/10/2026, frontend `3d1e84a`), commit `3ea7fa4`.** `test/front/admin-users.test.tsx`, caso `FIX-27…`, pasa: ante `502 EMAIL_DELIVERY_FAILED` se muestra el mensaje específico, el formulario no navega y conserva los datos. Los demás casos de FRN-06 siguen en verde (19/19).

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** `FIX-24` y `BAC-16`.
- **Problema y evidencia:** el entregable 2 de `FIX-24` pide *"manejar `502 EMAIL_DELIVERY_FAILED` con un mensaje claro: 'No se pudo enviar el correo; el usuario no fue creado'"*. El mensaje se arma bien en `components/features/createuser/services/createUserService.ts`, que relanza `ApiRequestError('No se pudo enviar el correo de activación; el usuario no fue dado de alta.', 502, 'EMAIL_DELIVERY_FAILED')`. Pero el `catch` de `components/features/createuser/hooks/useCreateUser.ts` lo ignora y siempre muestra el texto genérico *"No se pudo completar la creación del usuario."*. El administrador no se entera de que el problema fue el correo. La prueba `test/front/admin-users.test.tsx`, caso *"FIX-27 si el correo no se pudo enviar (502 EMAIL_DELIVERY_FAILED)…"*, falla con `Unable to find an element with the text: /no se pudo enviar el correo|servidor de correo no está disponible/i`.
- **Entregable:**
  1. En el `catch` de `useCreateUser.submitNewUser`, si `error instanceof ApiRequestError && error.errorCode === 'EMAIL_DELIVERY_FAILED'`, usar `error.message` como descripción del toast de error. Mantener el mensaje genérico para el resto de los errores.
  2. Conservar el comportamiento actual ante el error: no navegar a `/users` y dejar cargados los datos del formulario para reintentar.
- **Criterio de éxito:** ante un 502 del correo el administrador ve el mensaje específico, sigue en el formulario de alta y los datos se conservan. Los casos `FIX-27…` de `admin-users.test.tsx` pasan y los demás casos de FRN-06 siguen en verde.

### `FIX-28` - Regresión: el logout del frontend no envía el access token y la sesión no se revoca (`FRN-13` / `BAC-17`) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 01/10/2026, frontend `3d1e84a`), commit `cfb88f7`.** Se quitó `skipAuthorization: true` de `logoutSession()`: el logout viaja con `Authorization: Bearer` y la cookie `centinela_refresh`. Pasan los tres controles:
> - `session-security.test.ts`, caso `FIX-28…` (12/12 en el archivo).
> - `session_security_acceptance_test.go`, caso de integración `FIX-28…`: `204`, y después el access token viejo recibe `401`.
> - Prueba integral LOGIN-04 frontend ↔ backend, **paso 7**: la prueba completa pasa 10/10 por primera vez.

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir. **Prioridad alta**: la sesión queda abierta en el servidor después de "Cerrar sesión".
- **Depende de:** `FRN-13`, `BAC-17` y `SEC-01`.
- **Problema y evidencia:**
  1. En el commit `deb59cb` (*"Se limpia implementacion vieja de endpoint donde se utiliza body y header"*), `authService.logoutSession()` pasó a llamar a `POST /api/auth/logout` con `skipAuthorization: true`, es decir, **sin `Authorization: Bearer <accessToken>`**.
  2. El backend no cambió ese contrato: `cmd/api/main.go:157` monta `/auth/logout` con `middleware.RequireAuth`, porque `BAC-17` necesita el JTI del access token para revocar atómicamente el access y el refresh. Sin el Bearer responde `401 MISSING_TOKEN` antes de llegar al handler.
  3. El usuario no ve el error, porque `useLogout()` limpia la sesión local y redirige a `/login` igual, pero **el access token sigue siendo válido hasta que expira**. Esto incumple el entregable 1 de `FRN-13` y el criterio de `BAC-17`.
  - Pruebas que lo muestran:
    - `test/front/session-security.test.ts`, caso *"FIX-28 logoutSession envía POST /auth/logout con Authorization: Bearer…"*: `expected undefined to be 'Bearer access-123'`.
    - `test/back/session_security_acceptance_test.go`, caso *"FIX-28 FRN-13 integracion el logout tal como lo envia el frontend…"*: `esperado 204, recibido 401: MISSING_TOKEN`, y después del logout `/account/profile` sigue respondiendo 200.
- **Entregable:**
  1. En `components/features/auth/services/authService.ts`, quitar `skipAuthorization: true` de `logoutSession()` para que el cliente adjunte el `Bearer` del access token, manteniendo `credentials: 'include'` para la cookie `centinela_refresh`.
  2. Mantener la limpieza local actual de `useLogout()`: `catch` + `finally` con `clearAuthTokens()` y `navigate('/login', { replace: true })`.
- **Criterio de éxito:** los dos casos `FIX-28…` pasan. "Cerrar sesión" revoca en el servidor el access y el refresh (el access token revocado responde `401 TOKEN_REVOKED`), y los demás casos de `session-security.test.ts` siguen en verde.
- **Alternativa descartada:** aceptar en el backend un logout solo con la cookie. Sin el access token, el backend no puede revocar el JTI del access, que es lo que exige `BAC-17`.

### `FIX-30` - Helper `canOperateInstance` en `usePermissions()` (`FRN-18`) (Frontend)

> [!WARNING]
> **Estado: Implementado con problema (verificado el 01/10/2026, frontend `3d1e84a`), commits `54b397e` a `0b632f1`.**
>
> **Frontend cumple:** `src/hooks/usePermissions.ts` expone `canOperateInstance(vmid)`. Da `true` para `ADMIN` o `FULL_ACCESS` y `false` para `READ_ONLY` o un VMID sin asignar. Lo arma con `permisos: [{ vmid, nivelAcceso }]` de la sesión, que se carga desde `GET /account/profile` (`mapPerfilToUserSession`). Prueba: `navigation.test.tsx`, caso `FIX-30…`, que siembra ese formato.
>
> **Falta la fuente de datos real:** el backend no devuelve `permisos` en `GET /account/profile`; `UsuarioDetalleDTO` solo tiene `instanciasPermitidas`. Con el backend real, la sesión queda con `permisos: []` y **`canOperateInstance` da `false` para todo OPERATOR, aunque tenga `FULL_ACCESS`**. Prueba: `resource_access_acceptance_test.go`, caso `FIX-37…`. La corrección es **`FIX-37`** (backend), en [`futuro.md`](futuro.md).
>
> **Nota:** el hook quedó en `src/hooks/`. Es válido: el criterio de `SEC-03` no fija carpeta, y la prueba se ajustó el 01/10/2026.

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 0,5 h (una vez que exista `SEC-03`)
- **Ventana propuesta:** A definir, junto con `SEC-03` o inmediatamente después.
- **Depende de:** `SEC-03` y `FRN-18`.
- **Problema y evidencia:** el entregable 2 de `FRN-18` pide agregar a `usePermissions()` el helper `canOperateInstance(vmid)`. Debe devolver `true` solo para `ADMIN` o para instancias con `FULL_ACCESS`, a diferencia de `canAccessInstance(vmid)`, que acepta `READ_ONLY` y `FULL_ACCESS`. No existe, porque `usePermissions()` (`SEC-03`) no está implementado: `context/AuthContext.js` no exporta ningún hook. Prueba que falla: `test/front/navigation.test.tsx`, caso *"FIX-30 canOperateInstance distingue FULL_ACCESS de READ_ONLY (FRN-18)"*.
- **Entregable:**
  1. Exponer `canOperateInstance(vmid: number): boolean` en `usePermissions()`.
  2. Tomar el nivel por instancia de la sesión del usuario: por ejemplo `permisos: [{ vmid, nivelAcceso }]` en `centinela_user`, cargado con el mismo contrato de `GET /api/admin/users/:id/permissions`, o con un campo equivalente en `GET /account/profile`. La prueba siembra ese formato. Si se elige otra fuente, hay que avisar para alinear la prueba.
- **Criterio de éxito:** `canOperateInstance(101)` es `true` con `FULL_ACCESS`, es `false` con `READ_ONLY` y para un VMID no asignado, y siempre es `true` para `ADMIN`. El caso `FIX-30…` pasa.

### `LOGIN-04` - Prueba integral de autenticación y autorización

> [!NOTE]
> **Estado: Completada (verificado el 01/10/2026, backend `860b3c9` y frontend `3d1e84a`).** Las dos pruebas automatizadas pasan, dos corridas completas con el mismo resultado:
> - `test/back/login04_acceptance_test.go`: el circuito desde Go, pasos 1 a 7.
> - `test/back/login04_integracion_front_back_test.go`, que corre `test/front/login04-e2e.test.ts`: el **código real del frontend contra el backend real**, **10/10 pasos**. Era la primera vez que pasaba completa: el paso 7 se destrabó con `FIX-28`.
>
> **Cobertura del entregable:**
> - Alta con la clave temporal entregada solo por correo; cambio obligatorio de esa clave.
> - Enrolamiento y login con TOTP.
> - Acceso según rol; filtro por instancias (el OPERATOR solo ve la 101).
> - Renovación silenciosa con la cookie HttpOnly.
> - Logout con revocación en el servidor.
> - Recuperación de contraseña por el propio usuario, con el código enviado por correo.
> - Resets administrativos de 2FA y de contraseña.
> - Acciones administrativas registradas en la auditoría y visibles para el ADMIN.
>
> **Cobertura del criterio de éxito (rechazos con códigos controlados):**
>
> | Intento | Resultado |
> |---|---|
> | Omitir pasos | El JWT pre-2FA en una ruta protegida recibe `403`. Con el cambio de clave pendiente, la respuesta es `403 PASSWORD_CHANGE_REQUIRED`. Un TOTP inválido recibe `401` |
> | Usar credenciales anteriores | La clave temporal después del cambio, la contraseña anterior después de la recuperación o del reset y un código de recuperación inválido reciben `401` o `400` |
> | Acceder con otro rol | El OPERATOR en rutas de ADMIN recibe `403`, y la sesión no se cierra |
> | Consultar una instancia no asignada | `403 INSTANCE_ACCESS_DENIED` |
> | Reutilizar un token revocado | El access token después del logout recibe `401 TOKEN_REVOKED`; la sesión previa a los resets queda revocada. La revocación del refresh la cubre además `BAC-17` (`session_security…`) |
>
> **Nota:** el envío de correo se verificó con el mock y, en `BAC-16B`, con Mailpit en STARTTLS. Con el relay real falta `FIX-34`, que no es parte del criterio de esta tarea.

- **Área:** Frontend / Backend
- **Asignados:** Cristian, Tayra y Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** 18/09/2026, 14:00-17:00
- **Depende de:** `FRN-07`, `FRN-09`, `FRN-10`, `FRN-11`, `FRN-12`, `BAC-14`, `BAC-16`, `BAC-17`, `BAC-18`, `BAC-20` y `BAC-21`.
- **Entregable:** pruebas documentadas o automatizadas de creación de usuario, entrega y cambio de clave temporal, enrolamiento y login con TOTP, acceso según rol, filtro por instancias, recuperación administrativa de contraseña y 2FA, recuperación de contraseña por el propio usuario, logout/revocación de sesión y verificación de que las acciones administrativas quedan auditadas.
- **Criterio de éxito:** todos los recorridos válidos terminan con el acceso esperado y los intentos de omitir pasos, usar credenciales anteriores, acceder con otro rol, consultar una instancia no asignada o reutilizar un token revocado son rechazados con códigos HTTP controlados.
