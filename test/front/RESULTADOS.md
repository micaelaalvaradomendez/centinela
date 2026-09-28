# Resultados de las pruebas frontend

**Fecha:** 25/09/2026. **Frontend:** submódulo actualizado a `7fbf969` (último commit de `main`). **Comando:** `pnpm test`. Se hicieron 2 corridas con el mismo resultado (~25 s).

| Métrica | Valor |
|---|---:|
| Pruebas | 80 |
| Aprueban | **62** |
| Fallan | **18** |

## Por archivo

| Archivo | Pruebas | Fallan | Tareas |
|---|---:|---:|---|
| `admin-users.test.tsx` | 18 | 1 | FRN-05 ✅, FRN-06/FIX-24 🟡 (6/7, falla **FIX-27**), FRN-06B/**FIX-25** ✅, FIX-14 ✅ (4/4, incluido READ_ONLY), FRN-08 ✅ |
| `admin-recovery.test.tsx` | 4 | 3 | **FRN-11 / FIX-19** ❌ |
| `recover-password.test.tsx` | 7 | 6 | **FRN-12 / FIX-18** ❌, **FIX-20** (paso 3) ❌ |
| `navigation.test.tsx` | 6 | 3 | FRN-03 ✅, **SEC-03** ❌ (1/4) |
| `password-change.test.tsx` | 10 | 4 | FRN-10 ✅ (6), **FIX-21** ❌, **FIX-20** ❌ (3) |
| `audit.test.tsx` | 7 | 1 | FRN-14 ✅ (6), **FIX-22** ❌ |
| `session-security.test.ts` | 12 | 0 | FRN-13 ✅, SEC-02 ✅, FIX-08 ✅ |
| `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client` | 16 | 0 | FRN-01, FRN-02, FRN-04, LOGIN-02, LOGIN-03, FRN-09, FIX-07 ✅ |

## Fallos y causa

| Tarea | Prueba | Error | Causa en el código |
|---|---|---|---|
| FRN-11 / FIX-19 | restablecer contraseña pide confirmación | `Unable to find role="dialog"` | `Users.tsx:347`: el botón no tiene `onClick` |
| FRN-11 / FIX-19 | al confirmar envía POST …/password/reset | `Unable to find role="dialog"` | ídem |
| FRN-11 / FIX-19 | restablecer 2FA es una acción separada | `Unable to find role="button" … restablecer 2fa` | La acción no existe |
| FRN-12 / FIX-18 | paso 1, paso final y código inválido (3 casos) | `expected [] to have a length of 1` | `RecoverPassword.tsx` no llama a la API |
| FIX-18 | al terminar redirige a /login | `expected '/recover-password' to be '/login'` | No hay redirección |
| FIX-18 | el paso 2 no avanza sin código | `expect(element).not.toBeInTheDocument()` | `handleNext` avanza sin validar el código |
| FIX-20 | paso 3 no envía una clave no conforme | `Unable to find a label … nueva contraseña` | El submit del paso 3 vuelve al paso 1 sin validar ni enviar |
| FIX-20 | 3 casos en `ChangePassword` (sin mayúscula, sin dígito, sin especial) | `expected [[ '/api/account/password', … ]] to have a length of +0` | `ChangePassword.tsx` solo valida el largo |
| FIX-21 | la pantalla de cambio ofrece cerrar sesión | `Unable to find role="button" … cerrar sesión` | No hay botón |
| FIX-22 | un operador en /auditoria vuelve a /dashboard | `expected '/auditoria' to be '/dashboard'` | `applicationRoutes.tsx:50`: la ruta no está bajo `loadAdminSession` |
| SEC-03 | usePermissions / PermissionGate (2 casos) | `ningún módulo de src/context exporta usePermissions y PermissionGate` | `src/context/AuthContext.js` está vacío |
| SEC-03 | el ADMIN ve Usuarios y Auditoría | `Unable to find … link "Auditoría"` | El sidebar no tiene un enlace a `/auditoria` |
| FIX-27 | 502 `EMAIL_DELIVERY_FAILED` informa que no se creó | `Unable to find … no se pudo enviar el correo` | `useCreateUser.ts` reemplaza el mensaje de `createUserService.ts` por uno genérico |
