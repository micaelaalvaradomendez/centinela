# Informe de Estado de Tareas y Verificación por Pruebas

**Fecha de evaluación:** 21/09/2026  
**Fuentes analizadas:** [documentacion/actual.md](../documentacion/actual.md), [documentacion/terminado.md](../documentacion/terminado.md), [documentacion/futuro.md](../documentacion/futuro.md)  
**Modalidad de ejecución:** **Auditoría Integral en Crudo (0 Skips, 0 Todos)**  
**Entorno y suites ejecutadas:**
- Backend: [test/back/backend_acceptance_test.go](back/backend_acceptance_test.go), [test/back/resource_access_acceptance_test.go](back/resource_access_acceptance_test.go) y [test/back/password_recovery_acceptance_test.go](back/password_recovery_acceptance_test.go) (Go 1.27, PostgreSQL 16, Docker Compose)
- Frontend: [test/front](front) (Vitest, React Testing Library)

**Submódulos auditados (Fuente de verdad en código):**
- Backend: `3dadc48` (`Resolucion de conflictos y correccion de build tras el merge` - `origin/main`)
- Frontend: `1832e85` (`Merge pull request #43 from luzpacello/Nico` - `origin/main`)

---

## 1. Marco Metodológico y Resumen Ejecutivo

### Criterio de Verificación Exhaustiva (0 Skips, 0 Todos)
A solicitud de la dirección técnica, **se eliminaron todos los `t.Skip` y los `it.todo` de la suite de pruebas**. Todas las aserciones de contrato corren activamente contra los submódulos reales:
1. **Tareas Terminadas y Verificadas:** Se valida que no tengan regresiones y que satisfagan al 100% sus requerimientos.
2. **Tareas en Curso / Desvíos de Contrato:** Se audita la discrepancia exacta de rutas o componentes (`BAC-07` con `FIX-02`, `FRN-05` con `FIX-12`).
3. **Tareas Pendientes:** Se expone explícitamente el fallo por ausencia de endpoint (`404 Not Found`) o de componente de UI para dar seguimiento visual riguroso al avance pendiente de desarrollo.

### Métricas de Ejecución de Pruebas
| Suite | Total Pruebas Ejecutadas | Aprobadas | Fallidas | Omitidas (`skip` / `todo`) | % Aprobación |
|---|---|---|---|---|---|
| **Backend** (`test/back`) | 24 | 21 | 3 | **0** | **87.5%** |
| **Frontend** (`test/front`) | 31 | 28 | 3 | **0** | **90.3%** |
| **Total General** | **55** | **49** | **6** | **0** | **89.1%** |

*(Nota: La tasa de aprobación global aumentó del 82.4% al **89.1%** tras la integración de la rama de gestión de credenciales en el backend y la alineación de contratos en los arneses de prueba).*

### Diagnóstico Técnico Consolidado
1. **49 Pruebas en Verde (89.1% del sistema verificado empíricamente):**
   - **Auth & 2FA Base:** `BAC-01` a `BAC-04`, `LOGIN-01` a `LOGIN-03`, `BAC-10/11` (anti-replay y cifrado), `FRN-01`, `FRN-02`, `FRN-04`, `FRN-09`.
   - **Gestión de Usuarios & RBAC:** `BAC-05` (constraint único `usuario_id, vmid_proxmox`), `BAC-06` y `BAC-06B` (CRUD en `/api/admin/users`), `BAC-09` (roles), `FRN-03`, `FRN-06` y `FRN-06B` (alta/baja y edición en UI), `FRN-08` (403 con sesión preservada).
   - **Gestión de Seguridad Administrativa y Credenciales:** `BAC-12` (cambio obligatorio de clave validado en PostgreSQL), `BAC-13` (Reset 2FA), `BAC-15` (Reset contraseña), `BAC-16` (Entrega segura mediante puerto `EmailService` con `MockEmailService`; SMTP real descartado) y `BAC-17` (Logout y revocación de sesión en `sesiones_activas`).
   - **Recuperación de Contraseña por Usuario (RF-13):** `BAC-19` (`POST /api/auth/password/forgot`) y `BAC-20` (`POST /api/auth/password/reset`) operativos y verificados tras la integración en `main`.
   - **Auditoría Transversal (`BAC-18`):** `GET /api/admin/audit` y exportación CSV `GET /api/admin/audit/export` operan con éxito y rechazan al `OPERATOR` con 403.
   - **Contrato de Notificaciones y Eventos (`BAC-21`):** Struct Go `RealtimeEvent` e interfaz TypeScript sincronizados y validados.
2. **6 Pruebas en Rojo (Seguimiento exacto de los defectos de submódulo pendientes):**
   - **3 en Backend (`TestHitoControlDeAccesoBasadoEnRecursos`):**
     - `BAC-07` (`FIX-02`): En backend la ruta se montó como `PUT /api/admin/users/:id/instances` en vez de `/permissions`, y falta el `GET /permissions`.
     - `BAC-08`: Middleware de autorización por recurso (`(user_id, instance_id)`) aún no implementado; devuelve 404 en `/instances/9999`.
     - `BAC-14`: Endpoint `GET /api/instances` aún no implementado; devuelve 404.
   - **3 en Frontend (`admin-users.test.tsx`):**
     - `FRN-05` (`FIX-12`): En [applicationRoutes.tsx](../frontend/centinela/src/routes/applicationRoutes.tsx) falta asociar `loader: loadAdminSession` a `/users` para expulsar al `OPERATOR` hacia `/dashboard`.
     - `FRN-07` (2 pruebas, `FIX-14`): En [UserDetail.tsx](../frontend/centinela/src/pages/UserDetail.tsx) falta conectar la consulta `GET /api/instances` y el botón de guardado hacia `PUT /api/admin/users/:id/permissions`.

---

## 2. Análisis y Clasificación Integral de Fixes

A continuación se presenta el análisis detallado de cada fix registrado en `actual.md`, `terminado.md` y `futuro.md`, discriminando si el origen fue un **problema del test** (arnés desactualizado, falso positivo, aserción errónea) o un **problema de submódulo** (defecto de producto, lógica faltante o desvío de especificación).

| Fix ID | Tarea / Área | Clasificación de Origen | Diagnóstico y Causa Raíz | Estado Actual |
|---|---|---|---|---|
| `FIX-01` | `BAC-05` / Back | **Problema de Submódulo** | Colisión de nombre de índice `uniqueIndex:idx_usuario_vmid` en `models.go` que bloqueaba múltiples sesiones en `sesiones_activas`. | **Resuelto en Backend**. Verificado en `TestTareasBackendEIntegracion`. |
| `FIX-02` | `BAC-07` / Back | **Problema de Submódulo** | El backend montó `PUT /api/admin/users/:id/instances` en vez de `/permissions` y omitió `GET /permissions`. El test exige el contrato oficial. | **Pendiente en Backend** (`futuro.md`). Falla con 404. |
| `FIX-03` | `BAC-10` / Back | **Problema de Submódulo** | `GET /api/auth/2fa/qr` no rechazaba re-enrolar una cuenta con `is_2fa_enabled == true`. | **Resuelto en Backend**. Verificado en `TestTareasBackendEIntegracion`. |
| `FIX-05` | `FRN-05` / Front | **Problema de Submódulo** | La vista `/admin/users` mostraba datos estáticos en el JSX sin consultar `GET /api/admin/users`. | **Resuelto en Frontend**. La tabla consume y renderiza la API real. |
| `FIX-06` | `FRN-06` / Front | **Problema de Submódulo** | Modal de alta y edición de usuario no disparaba peticiones `POST` ni `PUT`. | **Resuelto en Frontend**. Verificado en `admin-users.test.tsx`. |
| `FIX-07` | `FRN-02` / Front | **Problema de Submódulo** | El cliente HTTP intentaba leer `body.code` en lugar del campo estándar `errorCode`. | **Resuelto en Frontend**. Verificado en `terminado.md`. |
| `FIX-08 (terminado)` | `FRN-03` / Front | **Problema de Submódulo** | Rutas protegidas (`/dashboard`, `/users`) se habían movido transitoriamente a layout público sin guards. | **Resuelto en Frontend**. Reubicadas bajo `ProtectedLayout`. |
| `FIX-08 (futuro)` | Front / Back | **Problema de Submódulos (Contrato)** | Desalineación entre valores de `errorCode` del backend y códigos de traducción visual en el frontend. | **Pendiente en Futuro** (`futuro.md`). |
| `FIX-09` | Front / Back | **Problema de Submódulo (Alcance)** | Llamada `/signup` en frontend sin endpoint (MVP es cerrado por invitación de admin). | **Pendiente en Futuro** (`futuro.md`). |
| `FIX-10` | `FRN-03` / Front | **Problema de los Test** | `navigation.test.tsx` escribía el token en `localStorage`, pero la aplicación leía de `sessionStorage`. El producto funcionaba bien; el test fallaba falsamente. | **Arreglado en los Test**. Se corrigió el test a `sessionStorage` y pasa 100%. |
| `FIX-11` | QA / Suites | **Tarea de Verificación** | Re-ejecución y validación integral continua del arnés de pruebas tras `FIX-01`. | **En ejecución continua**. |
| `FIX-12` | `FRN-05` / Front | **Problema de Submódulo** | `applicationRoutes.tsx` carece de `loader: loadAdminSession` en `/users`, permitiendo el acceso a `OPERATOR`. El test evalúa correctamente el PRD. | **Pendiente en Frontend** (`futuro.md`). |
| `FIX-13` | `BAC-19` / `BAC-20` | **Mixto (Submódulo + Test)** | 1) Backend tenía el código en `feat/gestion-credenciales` sin mergear (submódulo). 2) El test usaba `/password-recovery/*` en vez de `/password/forgot` y `/password/reset` estipulados en `Diseño de endpoints para front.md` (test). | **Arreglado en ambos**. Submódulo mergeado en `3dadc48` y test ajustado a rutas reales. Pasa 100%. |
| `FIX-14` | `FRN-07` / Front | **Problema de Submódulo** | `UserDetail.tsx` no consume `GET /api/instances` ni despacha `PUT /permissions` al guardar. Los tests fallan porque falta la lógica en UI. | **Pendiente en Frontend** (`futuro.md`). |
| `FIX-15` | Seguridad | **Mejora Arquitectónica** | Evaluación y diseño de migración a cookies `HttpOnly` para `refreshToken` contra XSS. | **Registrado en Futuro** (`futuro.md` y `siguientesetapas.md`). |

### Arreglos Realizados en el Arnés de Pruebas de Aceptación
1. **`FIX-10` en `navigation.test.tsx`:** Ajustado para sembrar tokens en `sessionStorage` conforme a la arquitectura de `tokenStorage.ts`.
2. **`FIX-13` en `password_recovery_acceptance_test.go`:** Ajustadas las rutas de prueba a `/api/auth/password/forgot` y `/api/auth/password/reset` según el contrato documentado.
3. **`signedAccessToken` y Sesiones en `backend_acceptance_test.go`:** Adaptado el generador de tokens de prueba para registrar sesiones activas en PostgreSQL (`sesiones_activas`), respetando la validación estricta de `jti` implementada en el backend.
4. **`BAC-12` en `backend_acceptance_test.go`:** Se adaptó la prueba para contemplar el cambio obligatorio de contraseña (`PUT /api/account/password`) liberando la sesión antes de invocar endpoints restringidos.
5. **`BAC-21` en `password_recovery_acceptance_test.go`:** Se corrigió la aserción que intentaba un HTTP `GET /events` inexistente; se validó el contrato estructural `RealtimeEvent` compartido entre Go y TypeScript.

---

## 3. Matriz de Estado Consolidada por Tarea

| ID | Área | Estado en Documentación | Estado en Pruebas | Resultado Empírico |
|---|---|---|---|---|
| `BAC-01` | Backend | `terminado.md` | ✅ PASS | PostgreSQL 16 con Docker Compose, migración y seed admin validados. |
| `BAC-02` | Backend | `terminado.md` | ✅ PASS | Hashing bcrypt persistido; validación y rechazo comprobados. |
| `BAC-03` | Backend | `terminado.md` | ✅ PASS | `POST /api/auth/login` emite JWT con rol e identidad tras validar en base. |
| `BAC-04` | Backend | `terminado.md` | ✅ PASS | Respuestas de error unificadas en 400 (`INVALID_REQUEST`) y 401 (`AUTH_FAILED`). |
| `LOGIN-01` | Backend | `terminado.md` | ✅ PASS | Flujo TOTP temporal emitido y verificado. |
| `LOGIN-02` | Frontend | `terminado.md` | ✅ PASS | Vista TOTP valida 6 dígitos y maneja estados visuales. |
| `LOGIN-03` | Fullstack | `terminado.md` | ✅ PASS | Recorrido e2e completo de login verificado de punta a punta. |
| `BAC-05` | Backend | `terminado.md` | ✅ PASS | Restricción única en `permisos_instancia(usuario_id, vmid_proxmox)` confirmada. |
| `BAC-06` | Backend | `terminado.md` | ✅ PASS | CRUD usuarios operativo bajo `/api/admin/users`. |
| `BAC-06B` | Backend | `terminado.md` | ✅ PASS | Edición y cambio de rol operativos bajo `PUT /api/admin/users/:id`. |
| `BAC-07` | Backend | `actual.md` | ❌ FAIL | Expone `/instances` en vez de `/permissions` y falta `GET /permissions` (`FIX-02`). |
| `BAC-08` | Backend | `actual.md` | ❌ FAIL | `GET /instances/9999` devuelve 404 (middleware por recurso no implementado). |
| `BAC-09` | Backend | `terminado.md` | ✅ PASS | `GET /api/roles` devuelve roles `ADMIN` y `OPERATOR`. |
| `BAC-10` | Backend | `terminado.md` | ✅ PASS | `GET /api/auth/2fa/qr` con rechazo 409 si ya está vinculado. |
| `BAC-11` | Backend | `terminado.md` | ✅ PASS | Validación de código, cifrado AES-256 y anti-replay de 30s verificados. |
| `BAC-12` | Backend | `actual.md` | ✅ PASS | `PUT /api/account/password` actualiza hash y desbloquea flag `cambio_contrasena`. |
| `BAC-13` | Backend | `actual.md` | ✅ PASS | `POST /api/admin/users/:id/2fa/reset` resetea secreto y rechaza operador con 403. |
| `BAC-14` | Backend | `actual.md` | ❌ FAIL | `GET /api/instances` responde 404 (inventario Proxmox VE no implementado). |
| `BAC-15` | Backend | `actual.md` | ✅ PASS | `POST /api/admin/users/:id/password/reset` genera clave temporal y rechaza operador. |
| `BAC-16` | Backend | `terminado.md` | ✅ PASS | Mailer Service simulado (`MockEmailService`) vía puerto `EmailService`; SMTP real descartado. |
| `BAC-17` | Backend | `actual.md` | ✅ PASS | `POST /api/auth/logout` operativo con revocación de `jti` en PostgreSQL. |
| `BAC-18` | Backend | `actual.md` | ✅ PASS | Base append-only de auditoría: consulta `GET /api/admin/audit` y exportación CSV activas. |
| `BAC-19` | Backend | `actual.md` | ✅ PASS | `POST /api/auth/password/forgot` emite respuesta genérica y genera código. |
| `BAC-20` | Backend | `actual.md` | ✅ PASS | `POST /api/auth/password/reset` valida código de 6 dígitos y actualiza contraseña. |
| `BAC-21` | Backend | `actual.md` | ✅ PASS | Contrato de datos `RealtimeEvent` sincronizado entre Go y TypeScript. |
| `FRN-01` | Frontend | `terminado.md` | ✅ PASS | Formulario de login visual y estados de carga. |
| `FRN-02` | Frontend | `terminado.md` | ✅ PASS | Validación en cliente de campos y formato. |
| `FRN-03` | Frontend | `terminado.md` | ✅ PASS | Rutas privadas protegidas bajo `ProtectedLayout`. |
| `FRN-04` | Frontend | `terminado.md` | ✅ PASS | Conexión con backend y almacenamiento de JWT en `sessionStorage`. |
| `FRN-05` | Frontend | `terminado.md` (revisar) | ❌ FAIL | Tabla consume API real, pero falta guard `loadAdminSession` (`FIX-12`). |
| `FRN-06` | Frontend | `terminado.md` | ✅ PASS | Modal alta de usuario con rol, clave temporal y baja (DELETE). |
| `FRN-06B` | Frontend | `terminado.md` | ✅ PASS | Detalle y edición de usuario con `PUT /api/admin/users/:id`. |
| `FRN-07` | Frontend | `actual.md` | ❌ FAIL (2) | Falta llamado a `GET /api/instances` y envío de permisos seleccionados en UI (`FIX-14`). |
| `FRN-08` | Frontend | `terminado.md` | ✅ PASS | Interceptor HTTP no destruye sesión ante 403 Forbidden. |
| `FRN-09` | Frontend | `terminado.md` | ✅ PASS | Pantalla de vinculación 2FA con QR, clave manual y campo OTP. |
| `FRN-10` | Frontend | `actual.md` | ⏳ Pendiente | Pantalla obligatoria de cambio de contraseña temporal. |
| `FRN-11` | Frontend | `actual.md` | ⏳ Pendiente | Botones de reset administrativo de 2FA y clave en UI. |
| `FRN-12` | Frontend | `actual.md` | ⏳ Pendiente | Conexión de `RecoverPassword.tsx` a la API. |
| `INF-03` | Infra | `terminado.md` | ℹ️ Docker | Contenedor persistente PostgreSQL funcional en pruebas locales. |
| `INF-04` | Infra | `terminado.md` | ⏳ En despliegue | Configuración de Nginx y vmbr1 fuera de suite local. |

---

## 4. Detalle de los 6 Fallos de Integración Pendientes

### 4.1 Defectos en Código Entregado (Fixes Prioritarios de Submódulo)
1. **Frontend — `FRN-05` (`FIX-12`):**
   - *Falla:* `admin-users.test.tsx > un usuario con rol OPERATOR no debe poder ver el panel de administración`.
   - *Causa:* En [applicationRoutes.tsx](../frontend/centinela/src/routes/applicationRoutes.tsx), la ruta `/users` solo tiene `loadProtectedSession`. Falta asignarle `loader: loadAdminSession` (ya escrito en [sessionGuard.ts](../frontend/centinela/src/components/features/auth/routes/sessionGuard.ts)).
2. **Backend — `BAC-07` (`FIX-02`):**
   - *Falla:* `resource_access_acceptance_test.go > PUT /api/admin/users/:id/permissions devuelve 404`.
   - *Causa:* En [main.go](../backend/cmd/api/main.go), la ruta fue nombrada `/instances` en lugar de `/permissions`, y no se expuso el método `GET`.

### 4.2 Tareas Pendientes de Implementación (Hito de Recursos)
3. **Backend — `BAC-08` (Middleware de Recursos):**
   - *Falla:* `resource_access_acceptance_test.go > GET /instances/9999 devuelve 404 en vez de 403`.
   - *Causa:* Middleware de validación de permisos por VMID aún no implementado.
4. **Backend — `BAC-14` (Inventario Proxmox VE):**
   - *Falla:* `resource_access_acceptance_test.go > GET /api/instances devuelve 404`.
   - *Causa:* Adaptador de Proxmox e inventario normalizado aún no implementados.
5. **Frontend — `FRN-07` (Listar Instancias en Detalle, `FIX-14`):**
   - *Falla:* `admin-users.test.tsx > consulta GET /api/instances para listar instancias`.
   - *Causa:* [UserDetail.tsx](../frontend/centinela/src/pages/UserDetail.tsx) no dispara el fetch a `/instances`.
6. **Frontend — `FRN-07` (Guardar Permisos de Instancias, `FIX-14`):**
   - *Falla:* `admin-users.test.tsx > permite seleccionar VMIDs y enviarlos al endpoint de permisos`.
   - *Causa:* [UserDetail.tsx](../frontend/centinela/src/pages/UserDetail.tsx) no dispara el `PUT` a `/permissions`.
