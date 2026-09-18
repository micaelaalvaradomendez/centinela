
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:






---


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

### BAC-06B — Edición de Usuario y Cambio de Rol (PUT /api/admin/users/{id}):
- Asignado: Lisandro | Estimación: 2h
- Depende de: BAC-06
- Entregable: Endpoint PUT /api/admin/users/{id} para actualizar nombre, correo, estado (isActive) y rol (role: ADMIN u OPERATOR).  
- Criterio de éxito: Un administrador puede cambiarle el rol a un usuario o desactivarlo; los cambios se reflejan de inmediato en la base de datos.  

### `BAC-07` - Asignación de permisos por recurso

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 3 h
- **Ventana propuesta:** 15/09/2026, 09:00-12:00
- **Depende de:** `BAC-05` y `BAC-06`.
- **Entregable:** `PUT /api/admin/users/{id}/permissions` y `GET /api/admin/users/{id}/permissions`, con actualización atómica de `user_instances`.
- **Criterio de éxito:** la relación usuario-instancia se persiste, reemplaza el conjunto anterior de forma atómica y solo puede gestionarla un administrador autorizado.

### `BAC-08` - Middleware de autorización por recurso

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** 15/09/2026, 13:00-16:00
- **Depende de:** `BAC-07`.
- **Entregable:** guard que valide `(user_id, instance_id)` en `user_instances` antes de consultar o enviar órdenes a Proxmox VE.
- **Criterio de éxito:** un operador sin permiso recibe `403` y la API no realiza ninguna llamada a Proxmox VE.

### `BAC-14` - Lectura mínima del inventario de Proxmox

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 3 h
- **Ventana propuesta:** 17/09/2026, 11:00-14:00
- **Depende de:** `BAC-08` y de las credenciales de lectura de Proxmox VE.
- **Entregable:** `GET /api/instances` consumiendo `/cluster/resources` o `/nodes/{node}/resources`, con una respuesta normalizada mínima que incluya ID, nombre, tipo, nodo y estado. Un administrador recibe todo el inventario y un operador solo las instancias asignadas.
- **Criterio de éxito:** `FRN-07` puede cargar IDs reales de VMs y LXC; una cuenta no puede descubrir instancias fuera de su alcance y los errores de Proxmox se traducen a una respuesta HTTP controlada.

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

### `FRN-07` - Selector de asignación de instancias

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 4 h
- **Ventana propuesta:** 15/09/2026, 09:00-13:00
- **Depende de:** `FRN-05`, `BAC-07` y `BAC-14`.
- **Entregable:** selector con instancias reales obtenidas de `GET /api/instances`, instancias asignadas y botón que envíe el array de IDs a `PUT /api/admin/users/{id}/permissions`.
- **Criterio de éxito:** el administrador puede guardar una matriz de permisos y verla nuevamente al abrir el usuario.

### `FRN-08` - Interceptor HTTP y manejo de `403 Forbidden`

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2 h
- **Ventana propuesta:** 15/09/2026, 14:00-16:00
- **Depende de:** `BAC-08` y `FRN-04`.
- **Entregable:** interceptor Axios o Fetch para enviar `Authorization: Bearer <JWT>` y mostrar un mensaje amigable ante un `403` sin cerrar la sesión.
- **Criterio de éxito:** las peticiones incluyen el JWT y un operador sin permiso recibe una respuesta visual clara.


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




