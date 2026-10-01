# Resultados de las pruebas frontend

**Fecha:** 01/10/2026 (2ª verificación, sin commits nuevos del frontend). **Frontend:** submódulo en `3d1e84a` (último commit de `main`). **Comando:** `pnpm test`. 2 corridas con el mismo resultado.

| Métrica | Valor |
|---|---:|
| Pruebas | 100 |
| Aprueban | 76 |
| Fallan | **14** |
| Omitidas | 10 (`login04-e2e.test.ts`: se ejecuta desde `test/back`, donde pasa 10/10) |

## Por archivo

| Archivo | Pruebas | Fallan | Tareas |
|---|---:|---:|---|
| `admin-users.test.tsx` | 20 | 1 | FRN-05, FRN-06, FRN-06B, FIX-25, FIX-14, FRN-18, FIX-27 ✅; **FIX-38** ❌ (nuevo) |
| `admin-recovery.test.tsx` | 4 | 0 | FRN-11 / FIX-19 ✅ |
| `recover-password.test.tsx` | 7 | 1 | FRN-12 / FIX-18 ✅; **FIX-29** ❌ |
| `password-change.test.tsx` | 12 | 3 | FRN-10, FIX-21 ✅; **FIX-29** ❌ (mayúscula, y 2 nuevos: mensaje de dígito y de largo) |
| `navigation.test.tsx` | 8 | 4 | FRN-03, FIX-30 ✅; **SEC-03** ❌ (3); **FIX-38** ❌ (nuevo) |
| `events-client.test.tsx` | 4 | 4 | **FRN-17C** ❌ (incluye el caso nuevo de retroceso exponencial) |
| `audit.test.tsx` | 7 | 1 | FRN-14 / FIX-22 ✅ 6/7; **FIX-36** ❌ |
| `session-security.test.ts` | 12 | 0 | SEC-02, FIX-08, FIX-28 ✅ |
| `login04-e2e.test.ts` | 10 | — | Omitida en esta suite; ver `../back/RESULTADOS.md` |
| `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client` | 16 | 0 | FRN-01, FRN-02, FRN-04, LOGIN-02, LOGIN-03, FRN-09, FIX-07 ✅ |

## Fallos y causa

| Tarea | Causa en el código |
|---|---|
| FIX-29 | `validatePasswordComplexity` prueba `/[0-9]/` donde debería probar mayúsculas; el mensaje de dígito no existe y el de largo dice "8 y 12" |
| SEC-03 | `usePermissions` no tiene `isOperator` ni `hasRole`; no existe `PermissionGate`, y `Sidebar.tsx` muestra "Auditoría" a todos los roles |
| FRN-17C | No existe `useEvents` |
| FIX-36 | `Auditoria.tsx` no tiene el filtro "Acción" (`56b88f5`) |
| FIX-38 | `informationOfUser.tsx` ofrece `READ_ONLY` en el selector de rol; `canAccessInstance` acepta el rol `READ_ONLY` |

## Cambios en la suite

- **`password-change.test.tsx`:** 2 casos FIX-29 que verifican los mensajes de dígito y de largo.
- **`events-client.test.tsx`:** caso de retroceso exponencial (FRN-17C).
- **`admin-users.test.tsx` y `navigation.test.tsx`:** casos de FIX-38.
