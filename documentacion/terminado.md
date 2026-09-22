

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

> [!WARNING]
> **Estado: Parcial.** La tabla ya consume la API real (`GET /api/admin/users`) y renderiza datos, pero las rutas `/users` en `applicationRoutes.tsx` carecen de `loader: loadAdminSession`, permitiendo el acceso a usuarios `OPERATOR`. Defecto registrado en [`futuro.md`](futuro.md) como `FIX-12`.

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 4 h
- **Ventana propuesta:** 14/09/2026, 09:00-13:00
- **Depende de:** `LOGIN-03` y contrato inicial de `BAC-06`.
- **Entregable:** vista `/admin/users` con nombre, correo, rol, estado de 2FA y acciones. Solo debe ser visible para administradores.
- **Criterio de éxito:** el administrador ve el listado real y un operador no puede acceder a la vista.

### `FRN-06` - Modal de creación y desactivación de usuarios

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 3 h
- **Ventana propuesta:** 14/09/2026, 14:00-17:00
- **Depende de:** `FRN-05` y `BAC-06`.
- **Entregable:** modal con usuario, correo y rol; visualización controlada de la contraseña temporal; confirmación para desactivar usuarios y toasts de resultado.
- **Criterio de éxito:** el administrador puede completar altas y bajas desde la interfaz y los errores se muestran de forma comprensible.

### FRN-06B — Modal de Edición de Usuario y Cambio de Rol:
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

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (posterior a `BAC-05`).
- **Depende de:** `BAC-03` y `BAC-05`.
- **Entregable:** `POST /api/auth/logout` y un mecanismo de invalidación de sesión (tabla `sessions` o lista de revocación con TTL) que puedan reutilizar `BAC-13` y `BAC-15` para revocar sesiones activas al resetear 2FA o contraseña.
- **Criterio de éxito:** un JWT revocado deja de autorizar peticiones aunque no haya expirado por tiempo; `BAC-13` y `BAC-15` consumen este mecanismo en lugar de simular la revocación.

### `BAC-18` - Base transversal de auditoría (append-only)

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

> [!WARNING]
> **Estado: Implementado con fallo en pruebas de integración.** Se implementó el puerto `InstanceRepository` (`internal/core/ports/instance_port.go`), el repositorio en PostgreSQL (`internal/adapters/secondary/postgres/instance_repository.go`) y la factoría `RequireInstanceAccess` (`internal/adapters/primary/http/middleware/instance_guard.go`), cableado en `cmd/api/main.go`. En las pruebas de integración (`test/back/resource_access_acceptance_test.go`), la llamada `GET /instances/9999` devolvió `404 Not Found` en lugar de `403 Forbidden` debido a que el grupo de rutas `/instances` aún no está montado activamente en el router de Gin a la espera de los handlers de Proxmox (`BAC-14`). Defecto registrado en [`futuro.md`](futuro.md) como `FIX-16`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 3 h
- **Depende de:** `BAC-07`.
- **Entregable:** guard que valide `(user_id, instance_id)` en `user_instances` antes de consultar o enviar órdenes a Proxmox VE.
- **Criterio de éxito:** un operador sin permiso recibe `403` y la API no realiza ninguna llamada a Proxmox VE.

### `FRN-07` - Selector de asignación de instancias

> [!WARNING]
> **Estado: Implementado con fallo en pruebas unitarias.** En `UserDetail.tsx` se maquetó la solapa de roles y permisos con el componente `AccessList`, pero se encuentra desconectada de la API: no realiza la consulta a `GET /api/instances` ni despacha la actualización a `PUT /api/admin/users/:id/permissions`. Los tests en `test/front/admin-users.test.tsx` fallan al esperar llamadas de red a dichos endpoints. Defecto registrado en [`futuro.md`](futuro.md) como `FIX-14`.

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

