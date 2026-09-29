
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 29/09/2026** (`FIX-22` y `FIX-23` se completaron y pasaron a [`terminado.md`](terminado.md)) (detalle en [test/informe.md](../test/informe.md)). En este archivo quedan **solo tareas sin implementar**; todas tienen pruebas que hoy fallan porque el código todavía no existe:
>
> | Tarea | Área | Qué falta |
> |---|---|---|
> | `SEC-03` | Frontend | `usePermissions()` / `PermissionGate` y ocultar "Auditoría" al OPERATOR en el menú |
> | `INF-08` | Infraestructura | Documentar las variables `EMAIL_PROVIDER` y `SMTP_*` |
> | `BAC-16B` | Backend | Adaptador `SmtpEmailService` elegido por `EMAIL_PROVIDER` |
> | `BAC-17B` | Backend | 1 sesión de usuario = 1 registro en `sesiones_activas` |
> | `BAC-18B` | Backend | Índice parcial en `sesiones_activas` y particionado de `auditoria` |
>
> **Dónde está el resto:**
> - Las tareas implementadas, completas o con problema, están en [`terminado.md`](terminado.md). Por ejemplo, `FIX-24` (con problema) y `FIX-25` (completa).
> - Las correcciones de las que tienen problema están en [`futuro.md`](futuro.md): `FIX-27` a `FIX-31` (`FIX-26` se descartó).

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (28/09/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. Solo queda el mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta, cuya corrección es `FIX-27` en [`futuro.md`](futuro.md), no en este archivo.

---

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


---


### `FIX-27` - Mensaje específico cuando falla el envío del correo en el alta (`FIX-24`) (Frontend)

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

### `FIX-29` - Validar la mayúscula en `validatePasswordComplexity` (`FIX-20`) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** `FIX-20`.
- **Problema y evidencia:** en `frontend/centinela/src/components/features/auth/utils/validateAuthenticationFields.ts`, `validatePasswordComplexity` tiene tres condiciones. La segunda prueba `/[0-9]/` pero agrega el mensaje *"Debe contener al menos una letra mayúscula."*, y no hay ninguna condición que pruebe mayúsculas. El dígito se valida una sola vez y la mayúscula nunca, así que el cliente acepta claves como `nueva1234!` o `sinmayus1!`, que el backend rechaza con `400 PASSWORD_CHANGE_FAILED` (`crypto.ValidarComplejidadContrasena`). Además, el mensaje de largo dice *"Debe tener 8 y 12 caracteres."*.
  - Pruebas que fallan:
    - `test/front/password-change.test.tsx`: *"FIX-29 no llama a la API si la contraseña nueva no cumple la complejidad del backend (sin mayúscula)"*.
    - `test/front/recover-password.test.tsx`: *"FIX-29 el paso 3 no envía una contraseña que no cumple la complejidad del backend"*.
- **Entregable:**
  1. Separar las condiciones: `/[A-Z]/` con el mensaje de mayúscula y `/[0-9]/` con el mensaje *"Debe contener al menos un número."*.
  2. Corregir el mensaje de largo: *"Debe tener entre 8 y 12 caracteres."*.
- **Criterio de éxito:** los casos `FIX-29…` pasan, y siguen en verde los de dígito, símbolo y largo, y los casos de `FRN-12` y `FIX-21`.


### `FIX-30` - Helper `canOperateInstance` en `usePermissions()` (`FRN-18`) (Frontend)

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

### `BAC-21B` (`BRG-01`) - Alineación de esquema (`auditoria`/`tareas_asincronas`), rutas de energía (`FULL_ACCESS`), regla de `DELETE` y extensión de `GET /api/instances`

- **Área:** Backend
- **Asignados:** Tayra y Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** Previo al inicio de `BAC-23B`, `BAC-24A/B` y `BAC-27`.
- **Depende de:** `BAC-14`, `FIX-16`, `SEC-04`, `FIX-23`.
- **Problema y evidencia (`cierre-fase-base.md`):**
  1. `etapa1.md` menciona `audit_logs` y `user_instances`, cuando las tablas reales en `domain/models.go` son `auditoria`, `tareas_asincronas` y `permisos_instancia`.
  2. `GET /api/instances` ya existe (`BAC-14`) y lo consume el selector `FRN-07`. Si se reemplaza su formato en vez de extenderlo, se rompe la pantalla de permisos.
  3. `FIX-16` creó `/api/instances/:vmid/start` y `/stop`, mientras que `BAC-24A` pide `/api/instances/:vmid/status/:action`. Además, falta exigir `FULL_ACCESS` en energía y restringir `DELETE /api/instances/:vmid` solo a `ADMIN`.
- **Entregable:**
  1. Persistir la auditoría de ciclo de vida en la tabla real `auditoria` (guardando `upid`, `action` y `resource_type` en `detalles` JSONB) y el estado de ejecución en `tareas_asincronas`, consultando permisos contra `permisos_instancia`.
  2. Extender `GET /api/instances` de forma 100% retrocompatible: mantener `{ id, name, type, node, status }` y agregar `{ ip, cpuUsage, ramUsage, maxRam, nivelAcceso, activeTask, instancesSummary }`.
  3. Montar `POST /api/instances/:vmid/status/:action` manteniendo alias en `/start` y `/stop`, protegidos con `RequireInstanceAccess(repo, "vmid", "FULL_ACCESS")`. Proteger `DELETE /api/instances/:vmid` (`BAC-24B`) con `RequireRole("ADMIN")` y validación de estado `stopped` (409 Conflict si está encendida).
- **Criterio de éxito:** `GET /api/instances` responde con los campos nuevos sin romper `FIX-14`; un operador `READ_ONLY` recibe 403 en acciones de energía; un `OPERATOR` recibe 403 en `DELETE`; todo se registra en `auditoria` y `tareas_asincronas`.

### `BAC-21C` (`BRG-02-BAC`) - Autenticación de `/api/events` por ticket efímero en Redis, bus Pub/Sub y revocación en vivo

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** Junto a `BAC-25B` y `BAC-26`.
- **Depende de:** `INF-06`, `BAC-17B`.
- **Problema y evidencia (`cierre-fase-base.md`):** `EventSource` (SSE) y WebSockets no permiten enviar el header `Authorization: Bearer`, y la cookie `HttpOnly` de `SEC-01` solo viaja a `/api/auth`. Además, si se revoca la sesión o un permiso de instancia, las conexiones abiertas deben cerrarse o actualizar su filtro, y el bus de eventos debe usar Redis Pub/Sub para funcionar con más de una réplica.
- **Entregable:**
  1. Crear el endpoint `POST /api/events/ticket` (bajo `RequireAuth`) que guarde en Redis un ticket de un solo uso con TTL de 30 segundos (`SET ws_ticket:<uuid> <usuario_id> EX 30`) y devuelva `{ ticket }`.
  2. En `GET /api/events?ticket=<uuid>`, validar y consumir atómicamente (`GETDEL`) el ticket en Redis antes de abrir el canal SSE/WebSocket.
  3. Publicar los eventos `TASK_FINISHED` de `BAC-25B` mediante **Redis Pub/Sub** (`centinela:events`) y cortar inmediatamente la conexión activa del usuario si recibe un evento de revocación de sesión (`BAC-17`) o recargar su filtro si cambian sus permisos (`BAC-07`).
- **Criterio de éxito:** `/api/events` rechaza con 401 tickets inválidos o reutilizados; al hacer logout o desactivar al usuario, el backend corta el stream inmediatamente.


