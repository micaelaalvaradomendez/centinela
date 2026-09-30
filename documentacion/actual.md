
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 29/09/2026** (`FIX-22` y `FIX-23` se completaron y pasaron a [`terminado.md`](terminado.md)) (detalle en [test/informe.md](../test/informe.md)). En este archivo quedan **solo tareas sin implementar**; todas tienen pruebas que hoy fallan porque el código todavía no existe:
>
> | Tarea | Área | Qué falta |
> |---|---|---|
> | `SEC-03` | Frontend | `usePermissions()` / `PermissionGate` y ocultar "Auditoría" al OPERATOR en el menú |
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


### `BAC-16B` - Adaptador Go `SmtpEmailService` para envío real de credenciales y códigos OTP (`RF-09` / `RF-13`)

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


### `BAC-17B` - Corrección de lógica de inserción en `sesiones_activas` (1 sesión = 1 registro) y almacenamiento en Redis con TTL

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
- **Problema y evidencia (análisis de cierre de la fase base):**
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
- **Depende de:** `INF-06A` (Redis local en el repo del backend), `BAC-17A` y `BAC-17B`.
- **Problema y evidencia (análisis de cierre de la fase base):** `EventSource` (SSE) y WebSockets no permiten enviar el header `Authorization: Bearer`, y la cookie `HttpOnly` de `SEC-01` solo viaja a `/api/auth`. Además, si se revoca la sesión o un permiso de instancia, las conexiones abiertas deben cerrarse o actualizar su filtro, y el bus de eventos debe usar Redis Pub/Sub para funcionar con más de una réplica.
- **Entregable:**
  1. Crear el endpoint `POST /api/events/ticket` (bajo `RequireAuth`) que guarde en Redis un ticket de un solo uso con TTL de 30 segundos (`SET ws_ticket:<uuid> <usuario_id> EX 30`) y devuelva `{ ticket }`.
  2. En `GET /api/events?ticket=<uuid>`, validar y consumir atómicamente (`GETDEL`) el ticket en Redis antes de abrir el canal SSE/WebSocket.
  3. Publicar los eventos `TASK_FINISHED` de `BAC-25B` mediante **Redis Pub/Sub** (`centinela:events`) y cortar inmediatamente la conexión activa del usuario si recibe un evento de revocación de sesión (`BAC-17`) o recargar su filtro si cambian sus permisos (`BAC-07`).
- **Criterio de éxito:** `/api/events` rechaza con 401 tickets inválidos o reutilizados; al hacer logout o desactivar al usuario, el backend corta el stream inmediatamente.


---
### `FRN-17C` (`BRG-02-FRN`) - Cliente de eventos con solicitud previa de ticket efímero y reconexión segura

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1,5 h
- **Ventana propuesta:** Junto a `FRN-17A`.
- **Depende de:** `BAC-21C` (`BRG-02-BAC`).
- **Problema y contexto:** El hook `useEvents` del frontend no puede pasar el JWT por header en `EventSource`/WebSocket ni exponer el access token largo en la URL. Debe solicitar primero el ticket efímero al backend.
- **Entregable:**
  1. En `useEvents` (`FRN-17A`), antes de abrir la conexión hacia `/api/events`, invocar `POST /api/events/ticket` con el interceptor autenticado (`Bearer`) y conectar a `/api/events?ticket=<uuid>`.
  2. Ante una desconexión de red, solicitar un nuevo ticket efímero aplicando retroceso exponencial; si `/api/events/ticket` responde `401`, disparar el cierre de sesión local y redirigir a `/login`.
- **Criterio de éxito:** El frontend se conecta a `/api/events` usando tickets de un solo uso sin exponer el JWT en la URL y se reconecta pidiendo un ticket fresco.


### `FIX-31` - Versionar y verificar el TLS de Nginx en el borde (`INF-05`) (Infraestructura)

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1 h
- **Ventana propuesta:** A definir.
- **Depende de:** `INF-04` e `INF-05`.
- **Problema y evidencia:** el criterio de `INF-05` exige que *"el tráfico hacia el sistema se sirva únicamente sobre HTTPS"*. La configuración de Nginx del servidor (CT 103) no está versionada en ningún repositorio. El deploy (`frontend/.github/workflows/deploy-front-test.yml`) solo ejecuta `nginx -t` y `reload` sobre lo que ya existe en el contenedor, y el único `nginx.conf` versionado (`docker/nginx.conf`, el del escenario de integración) escucha solo en `:80`. Por eso no se puede verificar el TLS. Prueba omitida: `test/back/cierre_fase_base_acceptance_test.go`, caso *"FIX-31 INF-05 TLS en el borde con Nginx"*.
- **Entregable:**
  1. Versionar la configuración de Nginx del borde, por ejemplo en `frontend/centinela/deploy/nginx.conf` o en `docker/`, con:
     - `listen 443 ssl` y las rutas de `ssl_certificate` y `ssl_certificate_key` (sin incluir los certificados);
     - `listen 80` que responda `return 301 https://$host$request_uri`;
     - `proxy_pass` de `/api/` al backend.
  2. Hacer que el deploy use ese archivo versionado.
- **Criterio de éxito:** la configuración versionada sirve solo HTTPS y redirige HTTP → HTTPS. A partir de ahí, el caso `FIX-31…` deja de omitirse y verifica el archivo (443 ssl, redirección 301 y ningún `server` que sirva contenido por `:80`).


### `INF-06A` - Redis local para desarrollo en el repositorio del backend

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

### `INF-06B` - Redis en el servidor (entornos de pruebas y estable)

- **Área:** Infraestructura
- **Asignados:** Nico y Lucas
- **Estimación:** 1 h
- **Ventana propuesta:** A definir. **No bloquea al backend**, que trabaja con `INF-06A`. Tiene que estar antes de desplegar `BAC-17B` y `BAC-21C` en el servidor.
- **Depende de:** `INF-03` e `INF-04`.
- **Entregable:**
  1. Desplegar Redis en el servidor, en el LXC del backend o en uno propio de `vmbr1`, con contraseña y sin exponerlo fuera de la red interna.
  2. Cargar `REDIS_ADDR` y `REDIS_PASSWORD` en el `.env` del backend de pruebas y del estable.
- **Criterio de éxito:** desde el LXC del backend, `redis-cli -h <host> -a <pass> PING` responde `PONG` y sin contraseña responde `NOAUTH`. Desde fuera de la red interna no se alcanza.

### `INF-08B` - Credenciales SMTP en `.env.example` y `.env` del repositorio del backend

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

### `INF-08A` - Configuración del SMTP en el servidor

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir. **No bloquea al backend**, que trabaja con `INF-08B`.
- **Depende de:** `INF-03`, `INF-04` e `INF-08B`.
- **Entregable:**
  1. Cargar `EMAIL_PROVIDER` y `SMTP_*` en el `.env` del backend del servidor (pruebas y estable).
  2. Habilitar y verificar la salida de red desde el LXC del backend hacia el host SMTP, por el puerto 587 (STARTTLS) o 465 (TLS). Por ejemplo, con `openssl s_client -starttls smtp -connect <host>:587`.
- **Criterio de éxito:** desde el LXC del backend se establece la conexión TLS con el servidor SMTP, y el backend desplegado envía el correo de alta de usuario.

### `BAC-17A` - Adaptador base de Redis en el backend (conexión, configuración y puerto)

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

---
# ETAPA 1 

### Bloque 8: Fixes detectados al comparar el simulador con el Proxmox real (30/09/2026)

*(Verificación: `cmd/proxmox-simulador` del backend (`e6e7dc5`) contra el Proxmox real `100.81.49.19`, con las mismas credenciales de `api-proxmox.md`, comparando campos y tipos de cada respuesta. Sobre el Proxmox real se hicieron solo lecturas (GET). Coinciden en estructura: `cluster/resources`, `nodes/{node}/status`, `nodes/{node}/lxc`, `rrddata`, `snapshot` y los errores de VMID inexistente, nodo inexistente, ruta inexistente y token inválido. No se pudo comparar `nodes/{node}/qemu` porque el servidor real no tiene VMs.)*

#### `FIX-32` - Fidelidad del simulador de Proxmox con la API real (Backend)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Depende de:** ninguna.
- **Problema y evidencia:** comparando endpoint por endpoint, el simulador responde distinto al Proxmox real 9.2.2 en estos puntos:
  1. **Tipo de dato:** en `GET /nodes/{node}/lxc/{vmid}/config`, `unprivileged` llega como número en el real (`1`) y como texto en el simulador (`"1"`). Un código que decodifique a entero funciona con uno y falla con el otro. Conviene revisar también los demás campos numéricos de `config`, como `cores`, `memory` y `swap`.
  2. **Campo faltante:** `GET /nodes/{node}/lxc/{vmid}/status/current` devuelve `"ha": {"managed": 0}` en el real, y el simulador no incluye `ha`.
  3. **Endpoints que el real tiene y el simulador responde con `501`:**
     - `GET /cluster/nextid` (real: `{"data":"106"}`): lo usa el alta de instancias (`RF-07`, `api-proxmox.md` §2).
     - `GET /nodes/{node}/tasks` (real: lista con `total` y `data[]` de `upid`, `status`, `starttime`, `endtime`, `user`): sirve para resincronizar tareas (`BAC-25C`, `RNF-04`).
  4. **Respuesta 401:** el real responde `401 Authentication failed!` **sin cuerpo**; el simulador devuelve un JSON `{"data":null,"message":"invalid token value!"}`. El backend no se ve afectado porque decide por el código HTTP, pero conviene imitar el real.
- **Entregable:**
  1. Devolver `unprivileged` y los demás numéricos de `config` como números.
  2. Agregar `ha: {managed: 0}` a `status/current` (con `managed: 1` y `state` en las instancias con HA, como la 100).
  3. Implementar `GET /cluster/nextid`, que devuelva el menor VMID libre desde 100 como string, y `GET /nodes/{node}/tasks` con las tareas creadas en la sesión del simulador.
  4. Responder el 401 sin cuerpo, con el status text `Authentication failed!`.
  5. Agregar a `simulador_test.go` casos para cada punto.
- **Criterio de éxito:** la comparación contra el Proxmox real no muestra diferencias de estructura ni de tipos en los endpoints de solo lectura.

#### `FIX-33` - Conflicto de bloqueo de Proxmox devuelto como `502 PROXMOX_UNAVAILABLE` (Backend)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1 h
- **Depende de:** `BAC-21B` / `BAC-24A` (rutas de energía).
- **Problema y evidencia:** con el backend apuntando al simulador, `POST /api/instances/110/start` seguido de inmediato por `POST /api/instances/110/stop` hace que Proxmox rechace la segunda orden con `500 can't lock file '/var/lock/qemu-server/lock-110.conf' - got timeout`, que es el comportamiento real documentado en las capturas del equipo. El backend la traduce a **`502 PROXMOX_UNAVAILABLE`, "Error al consultar la infraestructura subyacente"**. Pero Proxmox está disponible: la instancia está ocupada con otra tarea. El frontend no puede distinguir "servidor caído" de "esperá a que termine la operación en curso", y con ese mensaje FRN-16 y FRN-17B mostrarían un error equivocado.
- **Entregable:**
  1. En `proxmox/client.go`, reconocer las respuestas `500` cuyo mensaje contiene `can't lock file` o `is locked` y devolver un error de dominio, por ejemplo `ErrInstanciaOcupada`.
  2. En `instance_handler.go`, mapearlo a **`409 INSTANCE_BUSY`** con un mensaje claro.
  3. Agregar el código al inventario de `errorCode` (`FIX-08`) y a Swagger.
- **Criterio de éxito:** una segunda acción sobre una instancia con una tarea en curso responde `409 INSTANCE_BUSY`, y `502` queda reservado para cuando Proxmox no responde.


#### `BAC-28` - Completar el simulador de Proxmox para la Etapa 1

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** ninguna (el simulador ya existe desde `de407a5`).
- **Problema y evidencia:** el simulador imita inventario, estado del nodo, acciones con UPID, estado de tareas, métricas, snapshots, creación y configuración. Le faltan dos endpoints que necesita la Etapa 1, así que sin ellos esas tareas no se pueden desarrollar en local:
  1. `DELETE /nodes/{node}/{tipo}/{vmid}` (lo usa `BAC-24B`). Hoy responde `501`.
  2. La IP de cada instancia: `GET /nodes/{node}/qemu/{vmid}/agent/network-get-interfaces` para VMs y `GET /nodes/{node}/lxc/{vmid}/interfaces` para contenedores (lo usa `BAC-23A`).
- **Entregable:**
  1. `DELETE` que devuelva un UPID de tipo `qmdestroy`/`vzdestroy`, saque la instancia del inventario al terminar la tarea y responda error si está encendida, igual que Proxmox.
  2. Los dos endpoints de IP, con el formato real de Proxmox (capturarlo del servidor real con el token de solo lectura) y el caso "guest agent no está corriendo" para VMs sin agente.
  3. Pruebas en `cmd/proxmox-simulador/simulador_test.go`.
- **Criterio de éxito:** con el simulador, `BAC-23A` obtiene IPs y `BAC-24B` elimina una instancia detenida, sin tocar el servidor.