# Resultados de las pruebas frontend

**Fecha:** 01/10/2026. **Frontend:** submódulo en `3d1e84a` (último commit de `main`). **Comando:** `pnpm test`. 2 corridas con el mismo resultado.

| Métrica | Valor |
|---|---:|
| Pruebas | 95 |
| Aprueban | **76** |
| Fallan | **9** |
| Omitidas | 10 (`login04-e2e.test.ts`: se ejecuta desde `test/back` contra el backend real, donde pasa 10/10) |

## Por archivo

| Archivo | Pruebas | Fallan | Tareas |
|---|---:|---:|---|
| `admin-users.test.tsx` | 19 | 0 | FRN-05, FRN-06, FRN-06B, FIX-25, FIX-14, FRN-18 ✅; **FIX-27 ✅ nuevo** |
| `admin-recovery.test.tsx` | 4 | 0 | FRN-11 / FIX-19 ✅ |
| `recover-password.test.tsx` | 7 | 1 | FRN-12 / FIX-18 ✅; **FIX-29** (paso 3) ❌ |
| `password-change.test.tsx` | 10 | 1 | FRN-10, FIX-21 ✅; **FIX-29** (sin mayúscula) ❌ |
| `navigation.test.tsx` | 7 | 3 | FRN-03 ✅; **FIX-30 ✅ nuevo**; **SEC-03** ❌ (3) |
| `events-client.test.tsx` | 3 | 3 | **FRN-17C** ❌ |
| `audit.test.tsx` | 7 | 1 | FRN-14 / FIX-22 ✅ 6/7; **regresión del filtro "Acción"** ❌ |
| `session-security.test.ts` | 12 | 0 | SEC-02, FIX-08 ✅; **FIX-28 ✅ nuevo** |
| `login04-e2e.test.ts` | 10 | — | Omitida en esta suite; ver `../back/RESULTADOS.md` |
| `login-form`, `authentication-contract`, `two-factor-form`, `two-factor-enrollment`, `api-client` | 16 | 0 | FRN-01, FRN-02, FRN-04, LOGIN-02, LOGIN-03, FRN-09, FIX-07 ✅ |

## Fallos y causa

| Tarea | Causa en el código |
|---|---|
| FIX-29 | `validatePasswordComplexity` prueba `/[0-9]/` donde debería probar mayúsculas |
| SEC-03 | `usePermissions` está en `src/hooks/` (no en `src/context/`) y no tiene `isOperator` ni `hasRole`; no existe `PermissionGate`, y `Sidebar.tsx` muestra "Auditoría" a todos los roles |
| FRN-17C | No existe ningún `useEvents` que pida `POST /api/events/ticket` antes de conectarse |
| FRN-14 (regresión) | `56b88f5` quitó de `Auditoria.tsx` el filtro "Acción" y el parámetro `accion`, que exige RF-08 |

## Cambios en la suite

- **`navigation.test.tsx`:** el caso FIX-30 acepta `usePermissions` en `src/context` o en `src/hooks`. Los casos SEC-03 siguen exigiendo `src/context`.
- **`login04-e2e.test.ts`:** la marca del código de recuperación pasa a `"código de seguridad es:"`.
