
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 28/09/2026** (detalle en [test/informe.md](../test/informe.md)). En este archivo quedan **solo tareas sin implementar**; todas tienen pruebas que hoy fallan porque el código todavía no existe:
>
> | Tarea | Área | Qué falta |
> |---|---|---|
> | `SEC-03` | Frontend | `usePermissions()` / `PermissionGate` y ocultar "Auditoría" al OPERATOR en el menú |
> | `FIX-22` | Frontend | Poner `/auditoria` bajo el guard `loadAdminSession` |
> | `FIX-23` | Backend | Trigger que haga append-only la tabla `auditoria` |
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


### `FIX-22` - Guard administrativo en la ruta `/auditoria` (`FRN-14` / `BAC-18`) (Frontend)

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


