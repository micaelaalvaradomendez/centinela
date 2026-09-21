

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
