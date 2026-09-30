# Resultados de las pruebas frontend

**Fecha:** 30/09/2026. **Frontend:** submódulo en `749194e` (último commit de `main`). **Comando:** `pnpm test`. 2 corridas con el mismo resultado.

| Métrica | Valor |
|---|---:|
| Pruebas | 95 |
| Aprueban | **74** |
| Fallan | **11** |
| Omitidas | 10 (`login04-e2e.test.ts`: corre desde `test/back` contra el backend real) |

## Por archivo

| Archivo | Pruebas | Fallan | Tareas |
|---|---:|---:|---|
| `admin-users.test.tsx` | 19 | 1 | FRN-05, FRN-06, FRN-06B, FIX-25, FIX-14, FRN-18 ✅; **FIX-27** ❌ |
| `admin-recovery.test.tsx` | 4 | 0 | FRN-11 / FIX-19 ✅ |
| `recover-password.test.tsx` | 7 | 1 | FRN-12 / FIX-18 ✅; **FIX-29** (paso 3) ❌ |
| `password-change.test.tsx` | 10 | 1 | FRN-10, FIX-21 ✅; **FIX-29** (sin mayúscula) ❌ |
| `navigation.test.tsx` | 7 | 4 | FRN-03 ✅; **SEC-03** ❌ (3); **FIX-30** ❌ |
| `events-client.test.tsx` **(nuevo)** | 3 | 3 | **FRN-17C** ❌ |
| `audit.test.tsx` | 7 | 0 | FRN-14 / FIX-22 ✅ |
| `session-security.test.ts` | 12 | 1 | SEC-02, FIX-08 ✅; **FIX-28** ❌ |
| `login04-e2e.test.ts` | 10 | — | Omitida en esta suite; ver `../back/RESULTADOS.md` |
| `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client` | 16 | 0 | FRN-01, FRN-02, FRN-04, LOGIN-02, LOGIN-03, FRN-09, FIX-07 ✅ |

## Fallos y causa

| Tarea | Causa en el código |
|---|---|
| FIX-28 | `authService.logoutSession()` usa `skipAuthorization: true` (`deb59cb`) |
| FIX-29 | `validatePasswordComplexity` prueba `/[0-9]/` donde debería probar mayúsculas |
| SEC-03 / FIX-30 | `context/AuthContext.js` no exporta `usePermissions` ni `PermissionGate`; `Sidebar.tsx` muestra "Auditoría" a todos los roles |
| FRN-17C | No existe ningún `useEvents` que pida `POST /api/events/ticket` antes de conectarse |
| FIX-27 | `useCreateUser.ts` reemplaza el mensaje del 502 por uno genérico |
