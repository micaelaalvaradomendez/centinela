# Resultados de las pruebas frontend

**Fecha:** 23/09/2026. **Frontend probado:** `3192cf4` (`origin/main`), exportado con `git archive` y ejecutado con `CENTINELA_FRONTEND_DIR`. **Comando:** `pnpm test` (2 corridas, mismo resultado, ~21 s).

| Métrica | `origin/main` (3192cf4) | Submódulo local (5b91e80 + 3 parches sin commitear) |
|---|---:|---:|
| Archivos | 12 (6 con fallos) | 12 (6 con fallos) |
| Pruebas | 72 | 72 |
| Aprueban | **55** | 50 |
| Fallan | **17** | 22 |

La referencia es la columna de `origin/main`. El submódulo local está 13 commits atrás: le falta SEC-02, por eso fallan 5 pruebas de `session-security` y una de FRN-10. Además, un parche local mueve `/auditoria` bajo el guard de admin y **oculta** el fallo de FRN-14.

## Por archivo (`origin/main`)

| Archivo | Pruebas | Fallan | Tareas |
|---|---:|---:|---|
| `admin-users.test.tsx` | 13 | 3 | FRN-05 ✅, FRN-06 ❌ (2), FRN-06B ❌ (1), FIX-14/FRN-07 ✅, FRN-08 ✅ |
| `admin-recovery.test.tsx` | 4 | 3 | FRN-11 ❌ |
| `recover-password.test.tsx` | 4 | 3 | FRN-12 ❌ |
| `navigation.test.tsx` | 6 | 3 | FRN-03 ✅, SEC-03 ❌ |
| `password-change.test.tsx` | 10 | 4 | FRN-10 🟡 (FIX-21, FIX-20) |
| `audit.test.tsx` | 7 | 1 | FRN-14 🟡 |
| `session-security.test.ts` | 12 | 0 | FRN-13 ✅, SEC-02 ✅, FIX-08 ✅ |
| `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client` | 16 | 0 | FRN-01, FRN-02, FRN-04, LOGIN-02, LOGIN-03, FRN-09, FIX-07 ✅ |

## Fallos y causa

Los 17 fallos son **del producto**. En la primera corrida también fallaron 4 pruebas por problemas del propio test, que ya se corrigieron:
- FRN-13 sembraba `centinela_user` y `centinela_refresh` en sessionStorage, que la app ya no usa.
- FRN-08 buscaba un texto que el toast renderiza dos veces.

| Tarea | Prueba | Error | Causa en el código |
|---|---|---|---|
| FRN-11 | restablecer contraseña pide confirmación | `Unable to find role="dialog"` | `Users.tsx:354`: el botón "Restablecer contraseña" no tiene `onClick` |
| FRN-11 | al confirmar envía POST …/password/reset | `Unable to find role="dialog"` | ídem; no hay diálogo ni llamada |
| FRN-11 | restablecer 2FA es una acción separada | `Unable to find role="button" … restablecer 2fa` | No existe la acción de reset 2FA en el menú de acciones |
| FRN-12 | paso 1 envía POST /auth/password/forgot | `expected [] to have a length of 1` | `RecoverPassword.tsx` solo hace `setCurrentStep`; no llama a la API |
| FRN-12 | paso final envía POST /auth/password/reset | `expected [] to have a length of 1` | ídem (el submit del paso 3 vuelve al paso 1) |
| FRN-12 | código inválido o vencido se informa | `expected [] to have a length of 1` | ídem; no se muestra `RESET_FAILED` |
| SEC-03 | usePermissions expone helpers | `ningún módulo de src/context exporta usePermissions y PermissionGate` | `src/context/AuthContext.js` está vacío |
| SEC-03 | PermissionGate contenido/fallback | ídem | ídem |
| SEC-03 | un ADMIN ve accesos, incluida Auditoría | `Unable to find … link "Auditoría"` | `Sidebar.tsx` no tiene ningún enlace a `/auditoria` |
| FRN-14 | un operador en /auditoria vuelve a /dashboard | `expected '/auditoria' to be '/dashboard'` | `applicationRoutes.tsx`: `/auditoria` está bajo `loadProtectedSession` y no bajo `loadAdminSession` |
| FRN-10 | la pantalla ofrece cerrar sesión | `Unable to find role="button" … cerrar sesión` | `ChangePassword.tsx` no tiene opción de logout (FIX-21) |
| FRN-10 | 3 casos: no cumple la complejidad del backend | `expected [ [ '/api/account/password', … ] ] to have a length of +0` | `ChangePassword.tsx:27` solo valida el largo; faltan las reglas de `crypto/password.go:89` (FIX-20) |
| FRN-06 | tras el alta informa el resultado (BAC-16) | `Unable to find … usuario creado\|correo…` | `CrearUsuarios.tsx:56` solo muestra algo si llega `contrasenaTemp`, pero el backend ya no la envía (BAC-16). Un alta exitosa **no muestra ninguna confirmación** |
| FRN-06 | "Eliminar usuario" solicita DELETE | `la acción no envió DELETE a la API` | `detailsUserPage.tsx:110`: botón sin `onClick` |
| FRN-06B | editar y guardar envía PUT /admin/users/:id | `guardar no envió PUT /admin/users/u2` | "Guardar cambios" solo persiste permisos y se deshabilita si no cambió ninguno. No existe `PUT /admin/users/:id` en el frontend |

## Nota de integración

FRN-13 y SEC-02 pasan del lado del cliente, pero **el backend todavía rechaza sus peticiones**: `POST /auth/logout` y `POST /auth/refresh` con body `{}` responden 400. Esto se verifica en `test/back`, en el caso `SEC-01 SEC-02 integracion…`, y se resuelve con SEC-01 en el backend.
