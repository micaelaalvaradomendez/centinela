# Informe de Estado de Tareas y Verificación por Pruebas

**Fecha de evaluación:** 22/09/2026  
**Fuentes analizadas:** [documentacion/actual.md](documentacion/actual.md), [documentacion/terminado.md](documentacion/terminado.md) y suites en [test](test)  
**Modalidad de ejecución:** **Auditoría Integral en Crudo (0 Skips, 0 Todos)**  

**Submódulos actualizados (Estado en `main`):**
- **Backend:** `7cadbce` (`Merge branch 'feature/user-resource-authorization'`)
- **Frontend:** `d088df2` (`Merge pull request #47 from luzpacello/rama-Beli`)

---

## 1. Resumen Ejecutivo de Métricas

Todas las aserciones corren activamente contra los submódulos reales en ejecución (PostgreSQL 16 en contenedor Docker para Backend; Vitest y React Testing Library para Frontend).

| Suite | Total Pruebas | Aprobadas | Fallidas | Omitidas (`skip` / `todo`) | % Aprobación |
|---|---|---|---|---|---|
| **Backend** ([test/back](test/back)) | 27 | 23 | 4 | **0** | **85.2%** |
| **Frontend** ([test/front](test/front)) | 39 | 35 | 4 | **0** | **89.7%** |
| **Total General** | **66** | **58** | **8** | **0** | **87.9%** |

---

## 2. Estado Real de las Tareas

### 2.1 Tareas Completas y Verificadas al 100% (51 Pruebas Pasando)

| ID | Área | Descripción y Alcance | Verificación / Test |
|---|---|---|---|
| `BAC-01` | Backend | Contenedor PostgreSQL 16 con Docker Compose, migración y seed admin inicial. | `TestTareasBackendEIntegracion/BAC-01` |
| `BAC-02` | Backend | Hash seguro con bcrypt, sin contraseñas planas y validación de credenciales. | `TestTareasBackendEIntegracion/BAC-02` |
| `BAC-03` | Backend | Endpoint `POST /api/auth/login` emite JWT firmado con `id` y `role`. | `TestTareasBackendEIntegracion/BAC-03` |
| `BAC-04` | Backend | Respuestas de error estandarizadas: `400` (`INVALID_REQUEST`) y `401` (`AUTH_FAILED`). | `TestTareasBackendEIntegracion/BAC-04` |
| `LOGIN-01` | Backend | Flujo TOTP: generación de secreto temporal, validación de 6 dígitos y anti-replay. | `TestTareasBackendEIntegracion/LOGIN-01` |
| `LOGIN-02` | Frontend | Pantalla y validación en cliente del código TOTP de 6 dígitos. | `two-factor-form.test.tsx` |
| `LOGIN-03` | Fullstack | Recorrido punta a punta de autenticación (credenciales -> TOTP -> Dashboard). | `LOGIN-03` y `login-form.test.tsx` |
| `BAC-05` | Backend | Esquema relacional; restricción única en `permisos_instancia(usuario_id, vmid_proxmox)`. | `TestTareasBackendEIntegracion/BAC-05` |
| `BAC-06` | Backend | CRUD de usuarios para administradores bajo `/api/admin/users`. | `TestTareasBackendEIntegracion/BAC-06` |
| `BAC-06B` | Backend | Modificación y cambio de rol mediante `PUT /api/admin/users/:id`. | `TestTareasBackendEIntegracion/BAC-06B` |
| `BAC-07` | Backend | Asignación atómica de permisos en `GET` y `PUT /api/admin/users/:id/permissions`. | `TestHitoControlDeAccesoBasadoEnRecursos/BAC-07` |
| `BAC-09` | Backend | Catálogo de roles del sistema vía `GET /api/roles`. | `TestTareasBackendEIntegracion/BAC-09` |
| `BAC-10` | Backend | Enrolamiento 2FA (`GET /api/auth/2fa/qr`) con rechazo 409 si ya está vinculado. | `TestTareasBackendEIntegracion/BAC-10` |
| `BAC-11` | Backend | Cifrado AES-256 de secreto TOTP y login posterior con secreto persistido. | `TestTareasBackendEIntegracion/BAC-11` |
| `BAC-12` | Backend | Cambio forzado de clave temporal en `PUT /api/account/password`. | `TestTareasBackendEIntegracion/BAC-12` |
| `BAC-13` | Backend | Reset administrativo de 2FA en `POST /api/admin/users/:id/2fa/reset` (403 para operador). | `TestHitoRecuperacionDeContrasenasYNotificaciones/BAC-13` |
| `BAC-15` | Backend | Reset administrativo de contraseña en `POST /api/admin/users/:id/password/reset`. | `TestHitoRecuperacionDeContrasenasYNotificaciones/BAC-15` |
| `BAC-16` | Backend | Mailer Service seguro simulado vía puerto `EmailService` (`MockEmailService`). | `TestHitoRecuperacionDeContrasenasYNotificaciones/BAC-16` |
| `BAC-17` | Backend | Revocación de sesión/JWT en `POST /api/auth/logout` registrada en `sesiones_activas`. | `TestHitoRecuperacionDeContrasenasYNotificaciones/BAC-17` |
| `BAC-18` | Backend | Auditoría append-only: registro transversal de acciones y `GET /api/admin/audit`. | `TestHitoRecuperacionDeContrasenasYNotificaciones/BAC-18` |
| `BAC-19` | Backend | Solicitud de recuperación de contraseña vía `POST /api/auth/password/forgot`. | `TestHitoRecuperacionDeContrasenasYNotificaciones/BAC-19` |
| `BAC-20` | Backend | Confirmación de recuperación de contraseña vía `POST /api/auth/password/reset`. | `TestHitoRecuperacionDeContrasenasYNotificaciones/BAC-20` |
| `BAC-21` | Fullstack | Contrato base de notificaciones `RealtimeEvent` sincronizado en Go y TypeScript. | `TestHitoRecuperacionDeContrasenasYNotificaciones/BAC-21` |
| `FRN-01` | Frontend | Formulario de login, campos y estados de carga. | `login-form.test.tsx` |
| `FRN-02` | Frontend | Validación en cliente del formulario de inicio de sesión. | `login-form.test.tsx` |
| `FRN-03` | Frontend | Rutas y navegación protegidas bajo `ProtectedLayout` con guard de sesión. | `navigation.test.tsx` |
| `FRN-04` | Frontend | Cliente de autenticación con guardado de token en `sessionStorage`. | `api-client.test.ts` |
| `FRN-05` | Frontend | Panel de usuarios con consumo de API real y guard administrativo. | `admin-users.test.tsx` |
| `FRN-06` | Frontend | Modal de alta de usuarios con clave temporal y acción de baja (`DELETE`). | `admin-users.test.tsx` |
| `FRN-06B` | Frontend | Edición de usuarios y actualización de rol mediante `PUT /api/admin/users/:id`. | `admin-users.test.tsx` |
| `FRN-08` | Frontend | Manejo de error 403 Forbidden sin invalidar sesión de forma abrupta. | `admin-users.test.tsx` |
| `FRN-09` | Frontend | Enrolamiento 2FA con QR, clave manual y validación en primer acceso. | `two-factor-enrollment.test.tsx` |

---

### 2.2 Tareas en Desarrollo / Parciales

| ID | Área | Estado Actual en Código | Faltante para Completar |
|---|---|---|---|
| `BAC-08` | Backend | **Implementado en capa de middleware/repo, no conectado a rutas.** En [backend/internal/adapters/primary/http/middleware/instance_guard.go](backend/internal/adapters/primary/http/middleware/instance_guard.go) se agregó `RequireInstanceAccess` y su puerto de persistencia en PostgreSQL (`eb84a58`). | Montar el middleware sobre los endpoints de Proxmox en [backend/cmd/api/main.go](backend/cmd/api/main.go) una vez que existan los handlers de instancias (`BAC-14`). Actualmente `GET /instances/:vmid` responde `404 Not Found`. |
| `FRN-07` | Frontend | **Maquetado visual presente, desconectado de API.** En [frontend/centinela/src/pages/UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx) existe la pestaña de asignación de instancias. | 1. Ejecutar petición `GET /api/instances` para listar las VMs disponibles.<br>2. Enviar selección mediante `PUT /api/admin/users/:id/permissions` al guardar. |
| `FRN-10` | Frontend | **Lógica iniciada en ramas de autenticación.** Manejo del flag `must_change_password` en flujo de login. | Consolidar la vista obligatoria de cambio de contraseña conectada a `PUT /api/account/password` e incorporar suite de aceptación en [test/front](test/front). |

---

### 2.3 Tareas No Empezadas / Pendientes

| ID | Área | Descripción | Bloqueante / Dependencia |
|---|---|---|---|
| `BAC-14` | Backend | Lectura y normalización del inventario Proxmox VE (`GET /api/instances`). | Requiere adaptador secundario Proxmox VE con credenciales de lectura. Bloquea a `BAC-08` y `FRN-07`. |
| `FRN-11` | Frontend | Botones de reset administrativo de 2FA y de contraseña desde la tabla de usuarios. | Depende de conectar UI a `POST /api/admin/users/:id/2fa/reset` y `POST /api/admin/users/:id/password/reset` (ambos ya operativos en backend). |
| `FRN-12` | Frontend | Conexión de la vista `RecoverPassword.tsx` a los endpoints `POST /api/auth/password/forgot` y `POST /api/auth/password/reset`. | Vistas maquetadas pero aún sin despacho HTTP real. |

---

## 3. Detalle de los Tests Fallidos y Auditoría de Tareas/Fixes

### 3.1 Backend (4 Fallos en `test/back`)
1. **`BAC-08_rechaza_instancia_no_asignada_con_403_sin_llamar_a_Proxmox` (`FIX-16` / `BAC-08`):**
   - *Resultado recibido:* HTTP `404 Not Found` en `GET /instances/9999`.
   - *Causa:* El middleware `RequireInstanceAccess` ya fue programado en el submódulo backend, pero la ruta `/instances/:vmid` todavía no está expuesta en el enrutador HTTP activo de [backend/cmd/api/main.go](backend/cmd/api/main.go). Registrado como `FIX-16` en [documentacion/futuro.md](documentacion/futuro.md).
2. **`BAC-14_filtra_GET_/api/instances_por_permisos_y_normaliza_respuesta` (`BAC-14`):**
   - *Resultado recibido:* HTTP `404 Not Found` en `GET /instances`.
   - *Causa:* El endpoint de inventario Proxmox VE no ha sido implementado en el backend.
3. **`BAC-17_revocacion_atomica_de_accessToken_y_refreshToken_en_logout` (`BAC-17`):**
   - *Resultado recibido:* `accessToken` sigue autorizado (HTTP 200) tras invocar `POST /api/auth/logout`.
   - *Causa:* En [backend/cmd/api/main.go](backend/cmd/api/main.go) el endpoint `/auth/logout` no cuenta con `middleware.RequireAuth()`, por lo que el JTI del `accessToken` nunca ingresa al contexto ni se desactiva en `sesiones_activas`.
4. **`SEC-01_emision_de_refreshToken_en_cookie_HttpOnly_y_soporte_de_rotacion` (`SEC-01`):**
   - *Resultado recibido:* `POST /api/auth/2fa/verify` no emite `Set-Cookie` con flag `HttpOnly`, `POST /api/auth/refresh` exige body y no lee cookie, y `logout` no invalida la cookie.
   - *Causa:* El backend opera exclusivamente mediante transporte de tokens por payload JSON y no tiene implementado el soporte de cookies `HttpOnly`.

### 3.2 Frontend (4 Fallos en `test/front`)
5. **`FRN-07 consulta GET /api/instances con Authorization Bearer para listar instancias` (`FIX-14` / `FRN-07`):**
   - *Resultado recibido:* `AssertionError: expected undefined to be defined`.
   - *Causa:* [frontend/centinela/src/pages/UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx) no realiza la llamada fetch a `/instances` al abrir el detalle de permisos. Registrado como `FIX-14` en [documentacion/futuro.md](documentacion/futuro.md).
6. **`FRN-07 permite seleccionar VMIDs y enviarlos al endpoint de permisos` (`FIX-14` / `FRN-07`):**
   - *Resultado recibido:* `AssertionError: expected undefined to be defined`.
   - *Causa:* El botón de guardar en [frontend/centinela/src/pages/UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx) no despacha el `PUT` a `/permissions`. Registrado como `FIX-14` en [documentacion/futuro.md](documentacion/futuro.md).
7. **`FRN-13 logoutSession limpia exhaustivamente sessionStorage ante fallos de red o errores HTTP` (`FRN-13`):**
   - *Resultado recibido:* `AssertionError: expected '{"id":"u1",...}' to be null`.
   - *Causa:* `logoutSession()` en [frontend/centinela/src/components/features/auth/services/authService.ts](frontend/centinela/src/components/features/auth/services/authService.ts) solo ejecuta `clearAuthTokens()`, pero no purga las claves `centinela_user` ni `centinela_pending_login` de `sessionStorage`.
8. **`SEC-02 renovación de sesión es compatible con respuestas del backend que no incluyen refreshToken en el cuerpo JSON` (`SEC-02`):**
   - *Resultado recibido:* `ApiRequestError: Expirado`.
   - *Causa:* En [frontend/centinela/src/services/apiClient.ts](frontend/centinela/src/services/apiClient.ts), `tryToRenewSession` exige obligatoriamente que `refreshToken` esté presente en el cuerpo JSON de `/auth/refresh`, abortando si este viaja en la cookie HttpOnly.

---

## 4. Estado de los Submódulos tras la Actualización

- **Backend (`7cadbce`):**
  - Se incorporó la rama `feature/user-resource-authorization`.
  - Se unificó el contrato hacia `GET` y `PUT /api/admin/users/:id/permissions` con deduplicación y persistencia atómica verificada.
  - Se incorporó la lógica del guard de autorización por recurso (`instance_guard.go`), listo para acoplarse al inventario de Proxmox.
- **Frontend (`d088df2`):**
  - Mantiene aprobadas el 100% de las pruebas de autenticación, 2FA, navegación protegida y gestión básica de usuarios.
  - Pendiente la integración de llamadas a `/instances` y `/permissions` en `UserDetail.tsx`.

