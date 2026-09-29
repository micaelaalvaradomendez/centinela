# Resultados de las pruebas frontend

**Fecha:** 29/09/2026. **Frontend:** submódulo en `749194e` (último commit de `main`). **Comando:** `pnpm test`. Se hicieron 2 corridas con el mismo resultado (~25 s).

| Métrica | Valor |
|---|---:|
| Pruebas | 82 |
| Aprueban | **74** |
| Fallan | **8** |

## Por archivo

| Archivo | Pruebas | Fallan | Tareas |
|---|---:|---:|---|
| `admin-users.test.tsx` | 19 | 1 | FRN-05 ✅, FRN-06/FIX-24 (falla **FIX-27**), FRN-06B/FIX-25 ✅, FIX-14 ✅, **FRN-18** entregable 1 ✅ (2), FRN-08 ✅ |
| `admin-recovery.test.tsx` | 4 | 0 | **FRN-11 / FIX-19** ✅ |
| `recover-password.test.tsx` | 7 | 1 | **FRN-12 / FIX-18** ✅ (6); **FIX-20** paso 3 ❌ |
| `navigation.test.tsx` | 7 | 4 | FRN-03 ✅ (2), **SEC-03** ❌ (1/4), **FIX-30** `canOperateInstance` ❌ |
| `password-change.test.tsx` | 10 | 1 | FRN-10 ✅, **FIX-21** ✅, **FIX-20** 🟡 (2/3: falla "sin mayúscula") |
| `audit.test.tsx` | 7 | 0 | FRN-14 ✅, FIX-22 ✅ (commit `5c3d780`) |
| `session-security.test.ts` | 12 | 1 | **FRN-13** ❌ (logout sin Bearer), SEC-02 ✅, FIX-08 ✅ |
| `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client` | 16 | 0 | FRN-01, FRN-02, FRN-04, LOGIN-02, LOGIN-03, FRN-09, FIX-07 ✅ |

## Fallos y causa

| Tarea | Prueba | Error | Causa en el código |
|---|---|---|---|
| **FRN-13** (regresión) | logoutSession envía … Authorization: Bearer | `expected undefined to be 'Bearer access-123'` | `authService.logoutSession()` usa `skipAuthorization: true` (commit `deb59cb`) |
| FIX-20 | cambio de clave "sin mayúscula" y recuperación, paso 3 | `expected [[ '/api/account/password' … ]] to have a length of +0` / la llamada a `/auth/password/reset` se envía | `validatePasswordComplexity` prueba `/[0-9]/` donde debería probar mayúsculas |
| SEC-03 | usePermissions, PermissionGate y `canOperateInstance` (FRN-18) | `ningún módulo de src/context exporta usePermissions y PermissionGate` | `AuthContext.js` solo maneja `TOKEN_REVOKED` |
| SEC-03 | un OPERATOR no ve accesos administrativos | `expect(element).not.toBeInTheDocument()` (enlace "Auditoría") | `Sidebar.tsx` (commit `81c6424`) muestra "Auditoría" a todos los roles |
| FIX-27 | 502 `EMAIL_DELIVERY_FAILED` informa que no se creó | `Unable to find … no se pudo enviar el correo` | `useCreateUser.ts` muestra un mensaje genérico |
