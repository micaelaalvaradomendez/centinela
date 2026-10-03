# Resultados de las pruebas frontend

**Fecha:** 03/10/2026. **Frontend:** submódulo en `070e96b` (último commit de `main`, incluye PR #72 a #77: `SEC-03`, `FIX-38`, `FIX-42`, `FIX-36` y niveles de acceso a instancias). **Comando:** `pnpm test`.

| Métrica | Valor |
|---|---:|
| Pruebas | 113 |
| Aprueban | **87** |
| Fallan | **16** |
| Omitidas | 10 (`login04-e2e.test.ts`: se ejecuta desde `test/back`) |

## Por archivo

| Archivo | Pruebas | Fallan | Tareas |
|---|---:|---:|---|
| `admin-users.test.tsx` | 20 | 0 | FRN-05, FRN-06, FRN-06B, FIX-25, FIX-27, SEC-03 (PR #72), FIX-38 (PR #73) ✅ |
| `navigation.test.tsx` | 9 | 0 | FRN-03, SEC-03, FIX-30, FIX-38, FIX-42 ✅ |
| `password-change.test.tsx`, `recover-password.test.tsx` | 19 | 0 | FRN-10, FRN-12, FIX-21, FIX-29 ✅ |
| `audit.test.tsx` | 7 | 0 | FRN-14 / FIX-22, FIX-36 (PR #77) ✅ |
| `events-client.test.tsx` | 8 | 8 | **FRN-17C** ❌ (4) y **FRN-17A** ❌ (4): no existe `useEvents` |
| `dashboard-metrics.test.tsx` | 3 | 3 | **FRN-19A** ❌: no hay medidores en `features/dashboard` |
| `instances-table.test.tsx` | 5 | 5 | **FRN-20A** ❌: `Instances.tsx` es una maqueta estática |
| `session-security`, `admin-recovery`, `login-form`, `authentication-contract`, `two-factor-*`, `api-client` | 32 | 0 | regresión ✅ |

## Fallos y causa

| Tarea | Causa en el código |
|---|---|
| FRN-17C / FRN-17A | No existe `useEvents`; `hooks/useWebSocket.js` sigue vacío |
| FRN-19A | `features/dashboard` solo tiene `DashCard.jsx`, vacío |
| FRN-20A | `Instances.tsx` no consulta `GET /api/instances` ni renderiza filas interactivas |

## Cambios en la suite (02/10/2026)

- **`app-providers.tsx` (nuevo):** monta los providers de `src/context` igual que `App.tsx`. Desde `SEC-03`, el menú, `PermissionGate` y la ficha del usuario leen `useAuth()`; las pruebas que renderizaban rutas sin el `AuthProvider` fallaban por un error de la prueba, no del producto. Lo usan `navigation`, `audit`, `password-change`, `instances-table` y `dashboard-metrics`.
- **`navigation.test.tsx` (SEC-03):** acepta el hook como `usePermissions`, `useAuthUser` o `useAuth`, con helpers booleanos o funciones, y `PermissionGate` con export por defecto. La tarea no fija nombres ("un contexto con Provider es opcional"). El caso de `FIX-38` evalúa `usePermissions` (el helper que nombra esa tarea), y uno nuevo, `FIX-42`, evalúa los demás hooks (`useAuth`).
- **`instances-table.test.tsx` (FRN-20A), reescrito:**
  - el tipo es `vm`, como en el contrato de `BAC-14` (antes, `qemu`);
  - verifica el `Bearer`, que el tipo se muestre como VM o LXC, que Running sea verde y Stopped gris, la copia de la IP y la botonera en cada fila;
  - ya no importa `AuthProvider` de forma fija ni exige un nombre accesible para la tabla;
  - se quitó un `if (copyButton)` que dejaba pasar el caso sin verificar.
- **`dashboard-metrics.test.tsx` (FRN-19A), reescrito:** antes exigía nombres de componente (`CpuGauge`…), props y clases CSS concretas. Ahora:
  - detecta los medidores por lo que renderizan;
  - compara el aspecto visual: igual en 20 % y 69 %, distinto en 70 % (D3);
  - verifica que el frontend tenga sus propias pruebas con los bordes 69 y 70, como pide el criterio.
- **`events-client.test.tsx` (FRN-17A), reescrito:** los tres casos tenían un `if (subscribe) … else expect(algo definido)`, y pasaban sin verificar nada. Ahora verifican:
  - que lleguen los eventos válidos y se descarten los inválidos;
  - la deduplicación por `id`;
  - **una sola conexión** con dos consumidores;
  - la eliminación de `useWebSocket.js` y la actualización de `notifications.ts`.

### Contrapruebas

Las pruebas nuevas se corrieron contra una copia temporal del frontend con implementaciones mínimas correctas (`CENTINELA_FRONTEND_DIR`, sin tocar el submódulo): **16 de 16 pasan**. Después se rompió cada implementación de una forma distinta, y cada vez falló exactamente la prueba que corresponde:

| Defecto introducido | Prueba que lo detecta |
|---|---|
| Umbral en 80 % o en 69 % | FRN-19A, caso D3 |
| Sin deduplicar / sin validar el contrato / una conexión por consumidor | FRN-17A, el caso de cada uno |
| Badges sin color / IP null como `-` / tipo sin traducir | FRN-20A, el caso de cada uno |

## Lo que la suite no cubre

- FRN-17A: la suscripción **filtrada** por tipo o `recursoId`. La tarea no fija su forma; se verifica indirectamente en `FRN-19C` y `FRN-20C`, cuando la usen.
- FRN-19A: que haya un medidor de cada tipo (CPU, RAM y almacenamiento). Se verifica el comportamiento de cada medidor que exista.
- FRN-19A: las pruebas del propio frontend se buscan, pero no se ejecutan: el frontend no tiene un comando de pruebas configurado.
