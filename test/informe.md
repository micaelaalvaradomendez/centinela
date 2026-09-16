# Informe de estado actual — login, 2FA e integración

Fecha: 16/09/2026.
Fuentes: [documentacion/actual.md](documentacion/actual.md), [test/back](test/back) (suite Go, ejecución real contra PostgreSQL) y [test/front](test/front) (suite Vitest sobre el código real del frontend).

Commits auditados (coinciden con `origin/main` de cada submódulo):

| Componente | Commit |
|---|---|
| Frontend | `dc605af` |
| Backend | `4536776` |

## 1. Resumen ejecutivo

- **Actualización de submódulos:** Se sincronizaron ambos submódulos con los últimos cambios de sus respectivas ramas `main`.
  - En `backend` (`4919a88` -> `4536776`): se modificó el DTO de login en [backend/internal/adapters/primary/http/auth_handler.go](backend/internal/adapters/primary/http/auth_handler.go) para requerir `password` en lugar de `contrasena`, se agregó el endpoint de logout (`POST /api/auth/logout`) y el endpoint de verificación de build (`GET /api/version`).
  - En `frontend` (`86b9a02` -> `dc605af`): se adaptaron los servicios de autenticación ([frontend/centinela/src/components/features/auth/services/authService.ts](frontend/centinela/src/components/features/auth/services/authService.ts) y [frontend/centinela/src/components/features/2fa/services/2fa.service.ts](frontend/centinela/src/components/features/2fa/services/2fa.service.ts)) para enviar `password`, usar cabecera `Authorization: Bearer` en 2FA y consumir `GET /api/auth/2fa/qr` y `POST /api/auth/2fa/verify`. Además se renovaron las vistas de Dashboard, Instancias y NotFound.
- **Resultado de pruebas frontend ([test/front](test/front)):** 13 aprobadas y 1 fallida. Se actualizaron las expectativas de contrato (`password`, método `qr`) y navegación. La única falla restante confirma que el input OTP todavía acepta letras.
- **Resultado de pruebas backend ([test/back](test/back)):** 9 subtests aprobados. Se actualizó la suite para enviar `password` y el recorrido completo con PostgreSQL, TOTP y JWT ahora pasa.
- **Cobertura preparada para Gestión Administrativa:** la suite backend ahora ejecuta BAC-05, BAC-06 y BAC-06B sobre `GET/POST /api/users`, `PUT/DELETE /api/users/:id`, validando contraseña temporal, `cambio_contrasena`, cambio de rol, desactivación y rechazo de `OPERATOR`. La suite frontend incorpora 12 criterios `todo` para FRN-05, FRN-06 y FRN-06B, pendientes de convertirse en pruebas ejecutables cuando exista `/admin/users`.
- **Cobertura preparada para el hito Control de Acceso Basado en Recursos:** BAC-07 ya valida asignación, reemplazo atómico observable y lectura de VMIDs permitidos en el detalle del usuario. BAC-08 y BAC-14 quedan como pruebas omitidas hasta que existan el guard por `(user_id, instance_id)`, `GET /api/instances` y el adaptador Proxmox. FRN-07 y FRN-08 tienen criterios `todo` para selector de instancias, Bearer y manejo visual de `403`.

## 2. Estado por tarea

| ID | Estado | Evidencia |
|---|---|---|
| `BAC-01` | ✅ Completa | El compose aislado de [test/back](test/back) levanta PostgreSQL, GORM automigra el esquema (tablas `organizaciones`, `usuarios`, `sesiones_activas`) y crea el admin seed con roles `ADMIN`/`OPERATOR`. |
| `BAC-02` | ✅ Completa | Hash bcrypt verificado y login válido e inválido comprobados contra PostgreSQL. |
| `BAC-03` | ✅ Completa | `POST /api/auth/login` emite un JWT temporal válido usando el campo `password`. |
| `BAC-04` | ✅ Completa | Payload incompleto e inválido devuelven `400` y `401` con el contrato JSON esperado. |
| `FRN-01` | ✅ Completa | Formulario de login (usuario/correo, contraseña, envío, estado de carga) validado por [test/front](test/front). |
| `FRN-02` | ✅ Completa | Validaciones de cliente bloquean campos vacíos y correo inválido antes de enviar. |
| `FRN-03` | ✅ Completa | Protección de rutas y navegación Dashboard -> Instancias comprobadas con los encabezados actuales. |
| `FRN-04` | ✅ Completa | El frontend envía `password` y acepta la respuesta pre-2FA emitida por el backend. |
| `LOGIN-01` | ✅ Backend implementado | El backend expone `GET /api/auth/2fa/qr` y `POST /api/auth/2fa/verify`, cifra el secreto con AES-256-GCM y valida códigos TOTP. |
| `LOGIN-02` | ⚠️ Parcial | El contrato de QR/verificación pasa usando `Authorization: Bearer` y `codigo`; el input OTP todavía permite caracteres no numéricos. |
| `LOGIN-03` | ✅ En integración activa | El flujo de autenticación conecta login con redirección a 2FA ([test/front/authentication-contract.test.ts](test/front/authentication-contract.test.ts) validado). |
| `BAC-07` | ⚠️ Parcial | La implementación actual permite reemplazar permisos mediante `PUT /api/users/:id/instances` y leerlos en `GET /api/users/:id`; falta el contrato documentado `GET/PUT /api/admin/users/:id/permissions`. |
| `BAC-08` | ⏳ Preparada | La prueba está reservada, pero aún no existe el guard por recurso ni una operación de instancia que permita comprobar el `403` sin llamada a Proxmox. |
| `BAC-14` | ⏳ Preparada | La prueba está reservada, pero aún no existe `GET /api/instances` ni el contrato normalizado del adaptador Proxmox. |
| `FRN-07` | ⏳ Preparada | Hay criterios frontend `todo` para listado, selección y persistencia de VMIDs; todavía no existe la vista. |
| `FRN-08` | ⏳ Preparada | Hay criterios frontend `todo` para Bearer y `403` sin cerrar sesión; falta la integración de recursos. |

## 3. Detalle de fallas detectadas en la ejecución

### 3.1 Suite Frontend (13 pasadas / 1 fallida)

1. **`two-factor-form.test.tsx` — Validación estricta de dígitos OTP:**
   - Causa: El componente [frontend/centinela/src/pages/TwoFactor.tsx](frontend/centinela/src/pages/TwoFactor.tsx) permite escribir caracteres alfabéticos (`ABCDEF`) sin sanitizar o restringir el input a solo dígitos.

### 3.2 Suite Backend (9 subtests aprobados)

La suite se actualizó para enviar `{"email": "...", "password": "..."}`. Con ese contrato, todas las validaciones de backend pasan, incluyendo login, QR, TOTP, roles y protección de rutas.

## 4. Validaciones ejecutadas

| Suite | Comando | Resultado |
|---|---|---|
| `test/front` | `corepack pnpm --dir test/front test` | 13 aprobadas, 1 fallida (validación numérica OTP) |
| `test/back` | `go test -v -count=1 ./...` (con Docker Compose) | Todas las pruebas aprobadas |

## 5. Próximos pasos recomendados

1. Agregar sanitización numérica estricta en el input OTP de [frontend/centinela/src/pages/TwoFactor.tsx](frontend/centinela/src/pages/TwoFactor.tsx).
2. Re-ejecutar `test/front` para confirmar el 100% de aprobación.
3. Implementar `/admin/users` y su servicio frontend; convertir los `it.todo` de [test/front/admin-users.test.tsx](test/front/admin-users.test.tsx) en pruebas ejecutables.
4. Decidir si el contrato definitivo conserva `/api/users` o cambia a `/api/admin/users`; actualizar tests y frontend en el mismo cambio.
5. Definir el contrato final del hito de recursos: rutas de permisos, forma del inventario normalizado, operación protegida sobre instancias y mecanismo de mock/spy para garantizar que un `403` no llegue a Proxmox.
