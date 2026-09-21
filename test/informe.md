# Informe de Estado de Tareas y Verificación por Pruebas

**Fecha de evaluación:** 21/09/2026  
**Fuentes analizadas:** [documentacion/actual.md](documentacion/actual.md), [documentacion/terminado.md](documentacion/terminado.md)  
**Entorno y suites ejecutadas:**
- Backend: [test/back/backend_acceptance_test.go](test/back/backend_acceptance_test.go) y [test/back/resource_access_acceptance_test.go](test/back/resource_access_acceptance_test.go) (Go, PostgreSQL 16, Docker Compose)
- Frontend: [test/front](test/front) (Vitest, React Testing Library)

**Submódulos auditados (Fuente de verdad en código):**
- Backend: `019017b` (`refactor(auth): eliminar endpoint legacy relink y estandarizar reset administrativo`)
- Frontend: `68bd23f` (`Merge pull request #36 from luzpacello/rama-Beli`)

---

## 1. Marco Metodológico y Resumen Ejecutivo

### Criterio de Análisis y Fuente de Verdad
- **La fuente de verdad es el código real alojado en los submódulos:** Solo lo que está implementado en `backend/` y `frontend/` determina el comportamiento del sistema.
- **Rol de `documentacion/actual.md`:** Define la especificación y los criterios de éxito esperados. Las suites de pruebas de aceptación (`test/back` y `test/front`) se ejecutan contra los submódulos para comprobar empíricamente:
  1. **Tareas No Iniciadas:** Tareas de `actual.md` cuyo desarrollo aún no comenzó en los submódulos (ej. `BAC-08`, `BAC-14`, `FRN-07`, `FRN-08`). Responden `404 page not found` o carecen de componente. **No son bugs ni fixes**, sino trabajo pendiente de desarrollo.
  2. **Tareas en Desarrollo / Parciales:** Tareas con código ya presente en los submódulos pero que presentan desvíos frente a los criterios de `actual.md` (ej. `BAC-07` que implementó `/instances` en vez de `/permissions`, o `FRN-05`/`FRN-06` con maquetas estáticas desconectadas de la API).
  3. **Defectos y Regresiones en Código Existente (Fixes):** Fallas introducidas en código previamente desarrollado (ej. `FIX-01` que rompió el login repetido y desató fallos en cascada en backend; `FIX-08` que desprotegió rutas en frontend).
  4. **Tareas Completadas y Verificadas:** Tareas cuyo código está en el submódulo y satisface todos los criterios validados por las pruebas.

### Métricas de Ejecución de Pruebas
| Suite | Total Pruebas | Aprobadas | Fallidas | Pendientes (`todo`) | Estado General |
|---|---|---|---|---|---|
| **Backend** (`test/back`) | 13 | 3 | 10 | 0 | ❌ Bloqueo crítico por `FIX-01` + endpoints de recursos |
| **Frontend** (`test/front`) | 31 | 19 | 12 | 0 | ❌ Bloqueo por `FIX-08` + maquetas desconectadas |
| **Total General** | 44 | 22 | 22 | 0 | **50% Aprobación global** |

### Diagnóstico Técnico Principal
1. **Regresión Crítica en Backend (`FIX-01`):** En [backend/internal/core/domain/models.go](backend/internal/core/domain/models.go), el tag `uniqueIndex:idx_usuario_vmid` está compartido erróneamente en cuatro tablas (`SesionActiva`, `PermisoInstancia`, `Auditoria`, `Notificacion`). Esto provoca que PostgreSQL rechace logins repetidos con error `duplicate key value violates unique constraint "idx_usuario_vmid"` (SQLSTATE 23505), desatando un fallo en cascada que bloquea 5 pruebas de autenticación y RBAC (`LOGIN-01`, `LOGIN-03`, `BAC-09`, `BAC-10/11`, `BAC-05/06`).
2. **Regresión Crítica en Frontend (`FIX-08`):** En [frontend/centinela/src/routes/applicationRoutes.tsx](frontend/centinela/src/routes/applicationRoutes.tsx), las rutas `/dashboard`, `/instances`, `/users`, etc., fueron ubicadas bajo el layout público `MainLayoutAuth` como maquetas de diseño, desprotegiendo el acceso sin login ni 2FA y rompiendo el guard de `FRN-03`.
3. **Frontend Administrativo Sin Integración (`FIX-05`, `FIX-06`, `FIX-06B`):** Las páginas [frontend/centinela/src/pages/Users.tsx](frontend/centinela/src/pages/Users.tsx), [frontend/centinela/src/pages/CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) y [frontend/centinela/src/pages/UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx) son maquetas estáticas que utilizan arrays locales fijos, no llaman a los endpoints `GET /api/admin/users`, `POST /api/admin/users` ni `PUT /api/admin/users/:id`, no tienen guard de rol (`OPERATOR` accede libremente) y carecen de submit funcional.
4. **Hito de Control de Acceso por Recursos Sin Implementar:** Ni backend ni frontend tienen implementados `BAC-08` (middleware por recurso), `BAC-14` (inventario Proxmox normalizado), `FRN-07` (selector de instancias) ni `FRN-08` (interceptor de 403). Las llamadas a `/api/instances` devuelven `404 page not found`.

---

## 2. Matriz de Estado Consolidada por Tarea

| ID | Área | Estado Real | Verificación por Test | Detalle Técnico |
|---|---|---|---|---|
| `BAC-01` | Backend | ✅ Completa | PASS (`test/back`) | Docker Compose levanta PostgreSQL 16, GORM migra esquema base y crea admin seed. |
| `BAC-02` | Backend | ✅ Completa | PASS (`test/back`) | Hashing seguro con bcrypt; valida credenciales válidas y rechaza inválidas. |
| `BAC-03` | Backend | ❌ Rota por regresión | FAIL (`test/back`) | Lógica de login existe, pero falla en segundo login por colisión de índice `idx_usuario_vmid` (`FIX-01`). |
| `BAC-04` | Backend | ✅ Completa | PASS (`test/back`) | Manejo unificado de errores HTTP 400 (payload incompleto) y 401 (credenciales inválidas) con `errorCode`. |
| `LOGIN-01` | Backend | ⚠️ Bloqueada | FAIL (`test/back`) | Endpoints `/api/auth/2fa/qr` y `/api/auth/2fa/verify` implementados; bloqueados por falta de sesión válida (`FIX-01`). |
| `LOGIN-02` | Frontend | ✅ Completa | PASS (`test/front`) | Formulario TOTP valida código numérico de 6 dígitos y maneja errores de validación. |
| `LOGIN-03` | Fullstack | ⚠️ Bloqueada | PASS front / FAIL back | Frontend redirecciona correctamente a 2FA tras login, pero el recorrido e-2-e falla por login backend. |
| `BAC-05` | Backend | ❌ Incompleta | FAIL (`test/back`) | Alta de usuario funciona, pero falta restricción de unicidad en base para `(usuario_id, vmid_proxmox)`. |
| `BAC-06` | Backend | ⚠️ Parcial | FAIL (`test/back`) | CRUD y baja lógica implementados; falla por login bloqueante y discrepancia de prefijo de ruta `/admin`. |
| `BAC-06B` | Backend | ⚠️ Parcial | FAIL (`test/back`) | Edición de rol y usuario (`PUT /api/admin/users/:id`) implementada; bloqueada por login en la suite. |
| `BAC-07` | Backend | ❌ Incompleta | FAIL (`test/back`) | Endpoint `/api/admin/users/:id/permissions` responde 404 (expone `/instances`), sin `GET` aislado y con duplicados. |
| `BAC-08` | Backend | ❌ No iniciada | FAIL (`test/back`) | Middleware de recursos no existe; rutas `/api/instances/:id` devuelven `404 page not found`. |
| `BAC-09` | Backend | ⚠️ Bloqueada | FAIL (`test/back`) | Endpoint `GET /api/roles` existe y responde roles `ADMIN` y `OPERATOR`; bloqueado por falta de token. |
| `BAC-10` | Backend | ⚠️ Parcial | FAIL (`test/back`) | Genera QR y secreto cifrado. En `019017b` se agregó rechazo si ya está vinculado, pero suite está bloqueada por `FIX-01`. |
| `BAC-11` | Backend | ⚠️ Parcial | FAIL (`test/back`) | Cifrado AES-GCM y anti-replay de 30s implementados; verificación e-2-e bloqueada por `FIX-01`. |
| `BAC-12` | Backend | ⏳ No iniciada | Sin tests | Endpoint `POST /api/auth/change-password` con flag `must_change_password`. |
| `BAC-13` | Backend | ⏳ No iniciada | Sin tests | Restablecimiento administrativo de 2FA (`POST /api/admin/users/:id/reset-2fa`). |
| `BAC-14` | Backend | ❌ No iniciada | FAIL (`test/back`) | Endpoint `GET /api/instances` no existe (404); falta adaptador y normalización de inventario Proxmox. |
| `BAC-15` | Backend | ⏳ No iniciada | Sin tests | Reset administrativo de contraseña (`POST /api/admin/users/:id/reset-password`). |
| `BAC-16` | Backend | ⏳ No iniciada | Sin tests | Entrega de credenciales temporales vía SMTP. |
| `BAC-17` | Backend | ⏳ No iniciada | Sin tests | Revocación de sesiones y JWT (`POST /api/auth/logout`). |
| `BAC-18` | Backend | ⏳ No iniciada | Sin tests | Auditoría append-only a nivel motor de base de datos. |
| `BAC-19` | Backend | ⏳ No iniciada | Sin tests | Solicitud de código de recuperación por correo. |
| `BAC-20` | Backend | ⏳ No iniciada | Sin tests | Confirmación de recuperación de contraseña con código de 6 dígitos. |
| `BAC-21` | Backend | ⏳ No iniciada | Sin tests | Esquema genérico de eventos en tiempo real (WebSocket/SSE). |
| `FRN-01` | Frontend | ✅ Completa | PASS (`test/front`) | Formulario de login con campos de usuario, contraseña, botón y estados de carga. |
| `FRN-02` | Frontend | ✅ Completa | PASS (`test/front`) | Validación en cliente de campos requeridos y formato de correo antes del submit. |
| `FRN-03` | Frontend | ❌ Rota por regresión | FAIL (`test/front`) | Rutas desprotegidas por `FIX-08`; usuario sin sesión puede navegar a `/dashboard` e `/instances`. |
| `FRN-04` | Frontend | ✅ Completa | PASS (`test/front`) | Petición HTTP al backend enviando `email` y `password`, gestionando respuesta previa a 2FA. |
| `FRN-05` | Frontend | ❌ Incompleta | FAIL (`test/front`) | Panel `/users` usa lista hardcodeada, no consume `GET /api/admin/users` y no restringe rol `OPERATOR`. |
| `FRN-06` | Frontend | ❌ Incompleta | FAIL (`test/front`) | Formulario de alta no envía `POST /api/admin/users`, no muestra contraseña temporal ni ejecuta baja. |
| `FRN-06B` | Frontend | ❌ Incompleta | FAIL (`test/front`) | [UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx) no envía `PUT /api/admin/users/:id` al presionar "Guardar cambios". |
| `FRN-07` | Frontend | ❌ No iniciada | FAIL (`test/front`) | No existe selector de instancias reales ni envío de permisos por usuario. |
| `FRN-08` | Frontend | ❌ No iniciada | FAIL (`test/front`) | No hay interceptor para manejar error `403` en recursos preservando la sesión activa. |
| `FRN-09` | Frontend | ✅ Completa | PASS (`test/front`) | Pantalla `LoginContinuation` renderiza QR, muestra clave manual formateada y campo OTP. |
| `FRN-10` | Frontend | ⏳ No iniciada | Sin tests | Vista de cambio obligatorio de contraseña temporal. |
| `FRN-11` | Frontend | ⏳ No iniciada | Sin tests | Acciones de restablecimiento administrativo de 2FA y contraseña en UI. |
| `FRN-12` | Frontend | ⏳ No iniciada | Sin tests | Conexión del flujo de recuperación de contraseña en `RecoverPassword.tsx`. |
| `INF-03` | Infra | ℹ️ Verificada local | Docker Compose | PostgreSQL con volumen y persistencia validado en contenedor de testing. |
| `INF-04` | Infra | ⏳ En despliegue | Fuera de suite | Configuración de reverse proxy Nginx y subred `vmbr1`. |

---

## 3. Estado de Defectos y Regresiones en Código Existente (Fixes)

> **Aclaración conceptual clave:** Los "Fixes" corresponden únicamente a errores, regresiones o desvíos sobre código ya implementado en los submódulos. Las tareas como `BAC-08`, `BAC-14`, `FRN-07` y `FRN-08` **no son fixes**: son tareas de feature especificadas en [documentacion/actual.md](documentacion/actual.md) cuyo desarrollo en los submódulos aún no ha comenzado (por eso devuelven 404).

| Fix ID | Prioridad | Estado | Afecta a | Causa / Situación Actual en Submódulos |
|---|---|---|---|---|
| `FIX-01` | **CRÍTICA** | ❌ **Abierta** | Backend (`BAC-03`, `BAC-05`, `LOGIN-01`, `LOGIN-03`) | Tag `uniqueIndex:idx_usuario_vmid` duplicado en [backend/internal/core/domain/models.go](backend/internal/core/domain/models.go). Bloquea logins repetidos con error 401/23505 y provoca fallos en cascada. |
| `FIX-02` | Media | ❌ **Abierta** | Backend (`BAC-07`) | En el backend se implementó `PUT /api/admin/users/:id/instances` en vez de `PUT/GET /permissions` y no previene duplicados en un mismo payload. Registrado en [documentacion/futuro.md](documentacion/futuro.md). |
| `FIX-03` | Alta | ⚠️ **En verificación** | Backend (`BAC-10`) | Commits `924641f` y `019017b` en backend agregaron rechazo con 409 `TOTP_ALREADY_LINKED`; validación bloqueada por `FIX-01`. |
| `FIX-05` | Alta | ❌ **Abierta** | Frontend (`FRN-05`) | [frontend/centinela/src/pages/Users.tsx](frontend/centinela/src/pages/Users.tsx) usa mock `usersList` y no tiene guard de rol `ADMIN`. |
| `FIX-06` | Alta | ❌ **Abierta** | Frontend (`FRN-06`, `FRN-06B`) | [frontend/centinela/src/pages/CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) y [frontend/centinela/src/pages/UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx) no tienen handlers de submit hacia la API. |
| `FIX-07` | Media | ✅ **Resuelta** | Frontend (`apiClient.ts`) | El cliente extrae correctamente `errorCode` de los errores del backend (3 pruebas aprobadas en [test/front/api-client.test.ts](test/front/api-client.test.ts)). |
| `FIX-08` | **CRÍTICA** | ❌ **Abierta** | Frontend (`FRN-03`, `FRN-05`) | Rutas protegidas movidas a `MainLayoutAuth` en commit `6b08566`, permitiendo acceso público anónimo. |

---

## 4. Desglose Detallado de Pruebas Fallidas

### 4.1 Suite Backend (10 Fallas en 13 Pruebas)
1. **`BAC-05` (Unicidad de permisos):** No existe restricción única sobre `(usuario_id, vmid_proxmox)`. Se permiten pares duplicados en base de datos.
2. **`BAC-03` (Login repetido / colisión `FIX-01`):** `POST /api/auth/login` falla con HTTP 401 y código `duplicate key value violates unique constraint "idx_usuario_vmid"`.
3. **`LOGIN-01` (Generación/Validación TOTP):** Falla en cascada con 401 `MISSING_TOKEN` porque el login previo no emitió credenciales.
4. **`BAC-10` / `BAC-11` (TOTP vinculado y persistencia):** Falla al realizar el login posterior por la misma colisión de sesión.
5. **`BAC-09` (Listado de roles):** Responde 401 `MISSING_TOKEN` al requerir token de administrador.
6. **`BAC-05` / `BAC-06` / `BAC-06B` (Gestión de usuarios y RBAC):** Falla por falta de sesión válida en login previo.
7. **`BAC-06` / `BAC-06B` (Prefijo de ruta):** Discrepancia con el prefijo `/api/admin/users` documentado en pruebas aisladas.
8. **`LOGIN-03` (Recorrido de autenticación completo):** Falla por token JWT vacío a raíz del login bloqueado.
9. **`BAC-07` (Permisos por usuario):** Corregido el token del test con el `org_id` real de PostgreSQL; la creación del usuario pasa correctamente, pero `PUT /api/admin/users/:id/permissions` responde `404 page not found` porque el backend implementó la ruta `/instances` (`FIX-02`).
10. **`BAC-08` y `BAC-14` (Recursos e Inventario Proxmox):** Endpoints `/api/instances` y `/api/instances/:id` responden `404 page not found` porque **no se han empezado a desarrollar en el backend** (tareas pendientes de [documentacion/actual.md](documentacion/actual.md), no son un fix).

### 4.2 Suite Frontend (12 Fallas en 31 Pruebas)
1. **[test/front/navigation.test.tsx](test/front/navigation.test.tsx) (2 fallas - `FRN-03` / `FIX-08`):**
   - Un usuario anónimo accede a `/dashboard` sin ser redirigido a `/login`.
   - Navegación a `/instances` no cumple con la protección requerida.
2. **[test/front/admin-users.test.tsx](test/front/admin-users.test.tsx) (10 fallas - `FRN-05`, `FRN-06`, `FRN-06B`, `FRN-07`, `FRN-08`):**
   - **`FRN-05` (3 fallas):** `/users` permite acceso a rol `OPERATOR`; no se invoca `GET /api/admin/users` con Bearer token; la tabla ignora los datos provistos y muestra mocks.
   - **`FRN-06` (3 fallas):** El botón de crear no dispara `POST /api/admin/users`; no se muestra la modal con la contraseña temporal; no se ejecuta `DELETE /api/admin/users/:id`.
   - **`FRN-06B` (1 falla):** Al pulsar "Guardar cambios" en [frontend/centinela/src/pages/UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx), no se ejecuta `PUT /api/admin/users/:id`.
   - **`FRN-07` (2 fallas):** No se consultan instancias en `GET /api/instances` ni se envían VMIDs a la API de permisos.
   - **`FRN-08` (1 falla):** No existe manejo controlado de respuestas `403 Forbidden` preservando la sesión.

---

## 5. Próximos Pasos Técnicos Priorizados

1. **Resolver `FIX-01` en Backend ([backend/internal/core/domain/models.go](backend/internal/core/domain/models.go)):**
   - Eliminar `uniqueIndex:idx_usuario_vmid` de `SesionActiva.UsuarioID`, `Auditoria.UsuarioID` y `Notificacion.UsuarioID`.
   - Nombrar el índice compuesto de `PermisoInstancia` explícitamente como `uniqueIndex:idx_permiso_usuario_vmid` en `UsuarioID` y `VmidProxmox`.
   - *Impacto:* Desbloquea de inmediato 6 pruebas de backend (`BAC-03`, `LOGIN-01`, `LOGIN-03`, `BAC-09`, `BAC-10/11`, `BAC-05/06`).
2. **Resolver `FIX-08` en Frontend ([frontend/centinela/src/routes/applicationRoutes.tsx](frontend/centinela/src/routes/applicationRoutes.tsx)):**
   - Reubicar `/dashboard`, `/instances`, `/users`, `/users/new` y `/users/:userId` dentro de `ProtectedLayout` con `loader: loadProtectedSession`.
   - *Impacto:* Corrige las 2 fallas de [test/front/navigation.test.tsx](test/front/navigation.test.tsx) y restablece la seguridad del sistema.
3. **Conectar el Módulo Administrativo en Frontend (`FIX-05`, `FIX-06`, `FIX-06B`):**
   - Implementar guard de rol `ADMIN` en la ruta `/users`.
   - Reemplazar array `usersList` en [frontend/centinela/src/pages/Users.tsx](frontend/centinela/src/pages/Users.tsx) por llamada a `GET /api/admin/users`.
   - Conectar formularios en [frontend/centinela/src/pages/CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) y [frontend/centinela/src/pages/UserDetail.tsx](frontend/centinela/src/pages/UserDetail.tsx) a los endpoints `POST` y `PUT /api/admin/users`.
4. **Implementar Endpoints de Recursos en Backend (`BAC-08`, `BAC-14`, `BAC-07`):**
   - Implementar `GET /api/instances` con adaptador Proxmox normalizado.
   - Implementar middleware de validación `(user_id, instance_id)` que devuelva `403` sin consultar a Proxmox.

