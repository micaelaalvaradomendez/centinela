
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:


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


### `BAC-10` - Enrolamiento 2FA y generación de QR

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** 16/09/2026, 09:00-11:00
- **Depende de:** `BAC-05` y `LOGIN-03`.
- **Entregable:** `GET /api/auth/2fa/setup`, accesible mediante una sesión restringida, que genere un secreto TOTP único, una URI `otpauth://` y un código QR. También debe devolver el secreto para vinculación manual, sin persistirlo como habilitado antes de la confirmación.
- **Criterio de éxito:** un usuario sin 2FA obtiene un QR y una clave manual compatibles con una aplicación autenticadora; un usuario ya vinculado no puede reemplazar su secreto sin pasar por el flujo de restablecimiento.

### `BAC-11` - Validación y persistencia del secreto TOTP

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** 16/09/2026, 11:00-14:00
- **Depende de:** `BAC-10` y `BAC-05`.
- **Entregable:** `POST /api/auth/2fa/enable` para validar el primer código de seis dígitos, cifrar el secreto con AES-256 usando una clave externa a la base de datos y establecer `is_2fa_enabled: true`. El flujo de login también debe validar los códigos posteriores contra el secreto persistido antes de emitir el JWT de acceso completo.
- **Criterio de éxito:** el secreto nunca se almacena ni se expone en texto plano después del enrolamiento; un código válido habilita el 2FA y permite completar el login, mientras que códigos inválidos, vencidos o reutilizados son rechazados.

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


### `FRN-09` - Vinculación 2FA mediante QR

- **Área:** Frontend
- **Asignadas:** Belinda y Luz
- **Estimación:** 3 h
- **Ventana propuesta:** 16/09/2026, 14:00-17:00
- **Depende de:** `BAC-10`, `BAC-11` y `LOGIN-02`.
- **Entregable:** vista o modal obligatorio para cuentas sin 2FA, con QR, clave de vinculación manual, ingreso del código de seis dígitos y estados de carga y error.
- **Criterio de éxito:** la cuenta no puede acceder a las rutas protegidas hasta confirmar un código válido; al finalizar, continúa el login sin mostrar nuevamente el secreto.

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

## Fixes — hallazgos del informe de readiness (19/09/2026)

Estas tareas nacen de `documentacion/ANALISIS-READINESS.md`, que corrió las suites de aceptación (`test/back`, `test/front`) contra los commits realmente deployados en PRUEBAS (`52ebce1` backend, `9e8b45c` frontend). No son trabajo nuevo de alcance: son defectos en tareas que ya figuran como en curso en este mismo documento o como cerradas en `terminado.md`, y quedaron sin corregir porque nadie corrió la suite antes de desplegar. Cada tarea referencia la sección del informe donde está la evidencia.

### `FIX-01` - Colisión de índice único rompe el login repetido (CRÍTICA)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (máxima prioridad, antes que cualquier otra tarea de `BAC-05` en adelante).
- **Depende de:** ninguna; bloquea la verificación de `BAC-06`, `BAC-06B`, `BAC-07`, `BAC-08`, `BAC-09`, `BAC-10`, `BAC-11`, `BAC-14` y `LOGIN-03`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 3):** `models.go` reutiliza el mismo nombre `uniqueIndex:idx_usuario_vmid` en cuatro tablas (`SesionActiva.UsuarioID`, `PermisoInstancia.UsuarioID+VmidProxmox`, `Auditoria.UsuarioID`, `Notificacion.UsuarioID`). El índice terminó aplicado sobre `sesiones_activas.usuario_id`, así que la segunda sesión del mismo usuario viola la unicidad y el login devuelve `401 AUTH_FAILED` con `duplicate key value violates unique constraint "idx_usuario_vmid"`. Al mismo tiempo, el unique real que debía existir sobre `permisos_instancia(usuario_id, vmid_proxmox)` (parte del criterio de éxito de `BAC-05`) nunca se creó.
- **Entregable:** renombrar el índice compuesto de `PermisoInstancia` a algo único (por ejemplo `idx_permisos_usuario_vmid`), quitar el `uniqueIndex` de `SesionActiva.UsuarioID` (una sesión no es única por usuario) y agregar el constraint real faltante sobre `permisos_instancia(usuario_id, vmid_proxmox)`.
- **Criterio de éxito:** un mismo usuario puede loguearse más de una vez sin recibir `401`; `test/back` deja de reportar la colisión de `idx_usuario_vmid` en `backend_acceptance_test.go`; el chequeo de esquema de `BAC-05` confirma el unique real sobre `permisos_instancia`.

### `FIX-02` - Contrato de permisos por recurso desalineado con la documentación (`BAC-07`)

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1 h
- **Ventana propuesta:** A definir (posterior a `FIX-01`, para poder verificar con login funcionando).
- **Depende de:** `FIX-01`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 6):** `BAC-07` documenta `PUT/GET /api/admin/users/{id}/permissions`, pero el server expone `/api/users/:id/instances`. `FRN-07` (que todavía no está implementado) va a apuntar a la ruta documentada si nadie corrige el desvío.
- **Entregable:** alinear la ruta real con `/api/admin/users/{id}/permissions` (o, si el equipo decide mantener `/api/users/:id/instances`, actualizar la descripción de `BAC-07` en este documento para que documentación y server coincidan).
- **Criterio de éxito:** la ruta que consume `FRN-07` es la misma que describe `BAC-07` en `actual.md`; no quedan dos nombres distintos para el mismo endpoint entre doc y código.

### `FIX-03` - El guard de re-enrolamiento 2FA no rechaza reemplazar un secreto ya vinculado (`BAC-10`)

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (posterior a `FIX-01`).
- **Depende de:** `FIX-01`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 6):** el criterio de éxito documentado de `BAC-10` exige que "un usuario ya vinculado no pueda reemplazar su secreto sin pasar por el flujo de restablecimiento", pero el server actual sí permite generar un secreto nuevo sobre una cuenta con `is_2fa_enabled: true`.
- **Entregable:** `GET /api/auth/2fa/setup` (o el endpoint equivalente real) debe rechazar la generación de un nuevo secreto si la cuenta ya tiene 2FA habilitado, devolviendo un error que indique que debe usarse el flujo de restablecimiento administrativo (`BAC-13`).
- **Criterio de éxito:** una cuenta con `is_2fa_enabled: true` recibe un error controlado al pedir un nuevo QR; solo puede repetir el enrolamiento después de un reset administrativo.

### `FIX-04` - Contrato de endpoints 2FA desalineado entre documentación y servidor

- **Área:** Backend
- **Asignados:** Tayra y Lisandro
- **Estimación:** 1 h
- **Ventana propuesta:** A definir.
- **Depende de:** ninguna.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 6, nota "DOCS"):** `BAC-10`/`BAC-11` documentan `GET /api/auth/2fa/setup` y `POST /api/auth/2fa/enable`; el server real usa nombres distintos (`qr`/`verify`/`relink`). Nadie actualizó `actual.md` cuando cambió el contrato implementado.
- **Entregable:** decidir cuál nomenclatura es la definitiva (documentación o server) y aplicar el cambio en el lado que quedó desactualizado, para que `actual.md` describa exactamente las rutas que expone el backend.
- **Criterio de éxito:** las rutas descritas en `BAC-10`/`BAC-11` en este documento son las mismas que responde el server; `FRN-09` se integra sin adivinar nombres de endpoint.



### `FIX-05` - Panel de gestión de usuarios no consume la API real (`FRN-05`)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 3 h
- **Ventana propuesta:** A definir (posterior a `FIX-01`, para poder probar contra el backend real).
- **Depende de:** `FIX-01`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 5):** la vista `/admin/users` no llama a `GET /api/admin/users`; muestra datos fijos y no tiene guard de rol verificado contra la API.
- **Entregable:** conectar la tabla a `GET /api/admin/users`, mostrar nombre, correo, rol y estado de 2FA reales, y validar que un operador no pueda acceder a la vista.
- **Criterio de éxito:** el administrador ve el listado real devuelto por el backend; un usuario con rol `OPERATOR` no puede renderizar la vista.

### `FIX-06` - Modal de alta y edición de usuarios sin funcionalidad real (`FRN-06` / `FRN-06B`)

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 3 h
- **Ventana propuesta:** A definir (posterior a `FIX-01` y `FIX-05`).
- **Depende de:** `FIX-01` y `FIX-05`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 5):** el modal no tiene selector de rol, no dispara el `POST` al confirmar el alta, no muestra la contraseña temporal devuelta por el backend y no existe la acción de editar (`PUT /api/admin/users/{id}`).
- **Entregable:** completar el modal de alta con selector de rol y `POST` real mostrando la contraseña temporal, y agregar la acción de edición que llame a `PUT /api/admin/users/{id}` y actualice la tabla.
- **Criterio de éxito:** un administrador puede crear un usuario y ver su contraseña temporal, editar su rol y ver el cambio reflejado sin recargar manualmente la tabla.

### `FIX-07` - El frontend lee `body.code` en vez de `errorCode`

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1 h
- **Ventana propuesta:** A definir.
- **Depende de:** ninguna.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 6, hallazgo H1):** el backend responde los errores con el campo `errorCode` (por ejemplo `AUTH_FAILED`), pero `api.js` intenta leer `body.code`, que no existe en la respuesta real.
- **Entregable:** corregir `api.js` (y cualquier consumidor) para leer `errorCode` del cuerpo de la respuesta.
- **Criterio de éxito:** los mensajes de error específicos del backend (por ejemplo credenciales inválidas o TOTP vencido) se muestran correctamente en la interfaz en vez de un error genérico.

### `FIX-08` - Restaurar protección de rutas y guards de sesión (`FRN-03` / `FRN-05`)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 1 h
- **Ventana propuesta:** A definir (alta prioridad).
- **Depende de:** ninguna.
- **Problema (evidencia en `test/front/RESULTADOS.md`, sección pull 20/09):** el commit `6b08566` en `centinela/src/routes/applicationRoutes.tsx` movió `/dashboard`, `/instances`, `/users`, `/users/new`, `/users/:userId` y `/auditoria` al layout público (`MainLayoutAuth`) como "rutas temporales de diseño". Esto anuló el guard `loadProtectedSession`, permitiendo que cualquier usuario acceda a estas vistas sin iniciar sesión ni pasar el 2FA, rompiendo el criterio de éxito de `FRN-03` y `FRN-05` en producto y haciendo fallar `navigation.test.tsx`.
- **Entregable:** reubicar las rutas protegidas (`/dashboard`, `/instances`, `/users`, `/users/new`, `/users/:userId`, `/auditoria`) dentro del grupo con `ProtectedLayout` y `loader: loadProtectedSession`. Si se requieren vistas de diseño sin backend, utilizar mocks dentro del contexto autenticado en lugar de remover la protección de rutas.
- **Criterio de éxito:** un usuario no autenticado que intenta navegar a `/dashboard` o `/users` es redirigido a `/login`; las pruebas de navegación de `test/front` (`navigation.test.tsx`) validan la protección de rutas correctamente.

--- sin test ---

### `BAC-17` - Logout y revocación de sesión/JWT

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (posterior a `BAC-05`).
- **Depende de:** `BAC-03` y `BAC-05`.
- **Entregable:** `POST /api/auth/logout` y un mecanismo de invalidación de sesión (tabla `sessions` o lista de revocación con TTL) que puedan reutilizar `BAC-13` y `BAC-15` para revocar sesiones activas al resetear 2FA o contraseña.
- **Criterio de éxito:** un JWT revocado deja de autorizar peticiones aunque no haya expirado por tiempo; `BAC-13` y `BAC-15` consumen este mecanismo en lugar de simular la revocación.

### `BAC-12` - Cambio obligatorio de contraseña

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 16/09/2026, 09:00-11:00
- **Depende de:** `BAC-02`, `BAC-05` y `BAC-06`.
- **Entregable:** `POST /api/auth/change-password`, accesible con una sesión restringida, que compruebe la contraseña temporal, valide la nueva clave, actualice su hash e indique `must_change_password: false`. Mientras el indicador sea verdadero, el resto de endpoints protegidos debe permanecer bloqueado.
- **Criterio de éxito:** la contraseña temporal deja de ser válida después del cambio, la nueva contraseña nunca se guarda en texto plano y el usuario no obtiene acceso completo antes de finalizar el proceso.



### `BAC-15` - Restablecimiento administrativo de contraseña

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 17/09/2026, 14:00-16:00
- **Depende de:** `BAC-06` y `BAC-12`.
- **Entregable:** `POST /api/admin/users/{id}/reset-password`, restringido a administradores, que genere una contraseña temporal segura, actualice su hash, establezca `must_change_password: true` y revoque las sesiones activas del usuario.
- **Criterio de éxito:** un operador recibe `403 Forbidden`; la clave anterior deja de funcionar y el usuario debe cambiar la nueva contraseña temporal en el siguiente acceso.

### `BAC-13` - Restablecimiento administrativo de 2FA

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** 17/09/2026, 09:00-11:00
- **Depende de:** `BAC-06`, `BAC-08` y `BAC-11`.
- **Entregable:** `POST /api/admin/users/{id}/reset-2fa`, restringido a administradores, que invalide el secreto TOTP, establezca `is_2fa_enabled: false` y revoque las sesiones activas del usuario afectado.
- **Criterio de éxito:** un operador recibe `403 Forbidden`; tras el restablecimiento, los códigos del secreto anterior dejan de funcionar y el usuario debe repetir `BAC-10`, `BAC-11` y `FRN-09` en su siguiente acceso.
