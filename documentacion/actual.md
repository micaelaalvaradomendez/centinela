
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 30/09/2026** (detalle en [test/informe.md](../test/informe.md)).
> - `BAC-17A`, `BAC-17B`, `BAC-21C` e `INF-06A` están completas, y `INF-08B` está implementada con problema (su corrección es `FIX-34`). Las cinco pasaron a [`terminado.md`](terminado.md).
> - De la Etapa 1, `FIX-32`, `FIX-33` y `BAC-28` pasaron a [`terminado-1.md`](terminado-1.md). `BAC-28` tiene problemas y su corrección es `FIX-35`, en [`futuro-1.md`](futuro-1.md).
>
> **Quedan en este archivo solo tareas sin implementar**, todas con pruebas que hoy fallan:
>
> | Tarea | Área | Qué falta |
> |---|---|---|
> | `SEC-03` | Frontend | `usePermissions()` / `PermissionGate`, y ocultar "Auditoría" al OPERATOR |
> | `FIX-27`, `FIX-28`, `FIX-29`, `FIX-30` | Frontend | Mensaje del 502, logout con Bearer, mayúscula en el validador, `canOperateInstance` |
> | `FRN-17C` | Frontend | `useEvents` con ticket efímero (ya desbloqueada por `BAC-21C`) |
> | `BAC-16B` | Backend | `SmtpEmailService` elegido por `EMAIL_PROVIDER` |
> | `BAC-18B` | Backend | Índice parcial en `sesiones_activas` (`jti_access`) y particionado de `auditoria` |
> | `BAC-21B` | Backend | Campos nuevos en `GET /instances`, `/status/:action` y `DELETE` |
> | `FIX-31` | Infraestructura | Configuración de Nginx con TLS versionada |
> | `INF-06B`, `INF-08A` | Infraestructura | Redis y SMTP en el servidor (sin prueba automatizada) |

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

---
# ETAPA 1

> Las tareas de la Etapa 1 ya verificadas están en [`terminado-1.md`](terminado-1.md), y sus correcciones en [`futuro-1.md`](futuro-1.md).

