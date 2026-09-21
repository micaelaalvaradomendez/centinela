# Resultado base de las pruebas frontend

Fecha de ejecución: 18/09/2026 (commit `7775b62`).

## Resumen

| Métrica | Resultado |
|---|---:|
| Archivos de prueba | 5 |
| Pruebas totales | 28 |
| Aprobadas | 14 |
| Fallidas | 7 |
| Pendientes (`todo`) | 7 |

## Estado derivado

| Tarea | Resultado | Conclusión |
|---|---|---|
| `FRN-01` | Pruebas aprobadas | El formulario y su estado de carga cumplen el alcance comprobado. |
| `FRN-02` | Pruebas aprobadas | Los campos vacíos y el correo inválido se bloquean antes del envío. |
| `FRN-03` | Pruebas aprobadas | La protección de rutas y navegación Dashboard -> Instancias pasan con los encabezados actuales. |
| `FRN-04` | Pruebas aprobadas | El payload usa `password` y la respuesta pre-2FA del backend es aceptada. |
| `LOGIN-02` | Pruebas aprobadas | El contrato QR/verificación y la validación numérica del OTP pasan. |
| `LOGIN-03` | Pruebas aprobadas | La redirección al enrolamiento 2FA tras un login válido funciona correctamente. |
| `FRN-09` | Prueba enfocada aprobada | `LoginContinuation` muestra el QR, la clave manual formateada y el formulario de seis dígitos usando el JWT temporal. |
| `FRN-05` | **3 fallidas** | Sin guard de rol para `OPERATOR`, sin llamada a `GET /api/users`, tabla con datos fijos. |
| `FRN-06` | **3 fallidas** | Sin selector de rol, sin `onSubmit`/`POST /api/users`, sin mostrar la contraseña temporal. |
| `FRN-06B` | **1 fallida** | No existe ninguna acción de edición por usuario. |
| `FRN-07` | Pendiente (`todo`) | No existe ningún componente de selector de instancias contra el cual escribir una prueba real. |
| `FRN-08` | Pendiente (`todo`) | No existe ningún interceptor ni superficie de error `403` contra la cual escribir una prueba real. |

## Fallos detectados

Los criterios de `FRN-05`, `FRN-06` y `FRN-06B` dejaron de ser `it.todo` (que Vitest nunca ejecuta) y ahora son 7 pruebas reales en [test/front/admin-users.test.tsx](test/front/admin-users.test.tsx); **las 7 fallan** contra el código actual:

1. Un `OPERATOR` puede ver el panel `/users` igual que un `ADMIN` (sin guard de rol en [sessionGuard.ts](frontend/centinela/src/components/features/auth/routes/sessionGuard.ts)).
2. [Users.jsx](frontend/centinela/src/pages/Users.jsx) nunca llama a `GET /api/users`.
3. La tabla de usuarios ignora cualquier dato real y siempre muestra "Sin usuarios".
4. [CrearUsuarios.tsx](frontend/centinela/src/pages/CrearUsuarios.tsx) no tiene selector de rol.
5. El botón "Crear usuario" no dispara ningún `POST /api/users` (no hay `<form>` ni `onSubmit`).
6. No se muestra la contraseña temporal (`contrasenaTemp`) tras crear un usuario.
7. No existe ninguna acción "Editar" por fila.

## Alcance

Esta suite valida las tareas con responsabilidad frontend presentes en [documentacion/actual.md](documentacion/actual.md). Los criterios administrativos y del hito de recursos están preparados como `todo` porque todavía no existe la vista `/admin/users`, el selector de instancias ni la integración frontend de `GET /api/instances`.

Pull del 18/09/2026 (`78659d0` -> `7775b62`): se agregaron cuatro componentes UI aislados (`toast`, `bardge`, `tabs`, `nativeSelected`) que ningún archivo de página importa todavía. No cambian el estado de `FRN-05` a `FRN-08`; quedan disponibles para cuando se conecten los formularios y tablas administrativas.

## Pull del 20/09/2026 (`7775b62` -> `68bd23f`)

Se trajeron los commits que reescriben `Users.jsx` -> `Users.tsx`, amplían `CrearUsuarios.tsx`, agregan `UserDetail.tsx` y `Auditoria.tsx`, migran el almacenamiento de token de `localStorage` a `sessionStorage` (`tokenStorage.ts`) y suman `apiClient.ts`/`ApiResponseNotifier.tsx`. Al re-ejecutar `pnpm test` el resultado pasó de 14 aprobadas/7 fallidas/7 `todo` (28 pruebas) a **15 aprobadas/7 fallidas/7 `todo` (29 pruebas)**.

### Mejoras confirmadas

- La prueba de selector de rol en `CrearUsuarios.tsx` (parte de `FRN-06`) y la prueba de acción "Editar" por fila (`FRN-06B`) **ya no fallan**: `CrearUsuarios.tsx` ahora tiene selector de rol y `UserDetail.tsx` cubre la edición.

### Regresión CONFIRMADA en `FRN-03` (navbar y rutas base) — bug real de producto, no del test

Las dos pruebas de `navigation.test.tsx`, que en el pull anterior estaban **aprobadas**, ahora **fallan**. La primera hipótesis (más abajo, tachada por la verificación) era que el test estaba desalineado por el cambio de `localStorage` a `sessionStorage`. Se verificó esa hipótesis modificando temporalmente el test para sembrar la sesión en `sessionStorage` (la clave real que usa `tokenStorage.ts`) y volviendo a correrlo: **las dos pruebas siguieron fallando**, incluida "redirige a login cuando no existe una sesión" corriendo con el storage completamente vacío. Eso descarta el storage como causa y expone la causa real:

**`centinela/src/routes/applicationRoutes.tsx`, commit `6b08566` (Belinda, 20/09/2026 06:01), sacó `/dashboard`, `/instances`, `/users` y `/users/new` del grupo protegido por `loadProtectedSession` y las movió al grupo público (`MainLayoutAuth`), y sumó ahí mismo las rutas nuevas `/auditoria` y `/users/:userId`.** El comentario del propio commit dice: *"Rutas temporales de diseño: se pueden visualizar sin sesión ni backend."* El grupo protegido (`ProtectedLayout` + `loadProtectedSession`) quedó con una sola ruta índice que redirige a `/dashboard`, pero `/dashboard` ya no está adentro del grupo protegido — o sea que el guard no se ejecuta nunca para esas rutas.

**Efecto real, no solo de prueba:** cualquiera puede navegar directamente a `/dashboard`, `/instances`, `/users`, `/users/new`, `/users/:userId` o `/auditoria` sin haber iniciado sesión ni pasado el 2FA. Esto revierte el criterio de éxito de `FRN-03` (aislamiento por rol/sesión, `RF-01`) y de `FRN-05` (panel de administración solo para `ADMIN`) directamente en el producto, más allá de lo que cualquier suite de test pueda decir.

**Corrección sobre el informe previo:** [documentacion/ANALISIS-READINESS.md](../../documentacion/ANALISIS-READINESS.md) (19/09/2026) había atribuido esta falla únicamente al desalineamiento de storage del test (*"con login real, el dashboard se renderiza"*, concluyendo que el criterio "probablemente sigue cumpliéndose"). Esa conclusión es incorrecta para el commit `6b08566` (posterior a ese análisis, del 20/09): el dashboard se renderiza tanto con sesión como sin ella, porque la ruta ya no tiene guard. El desalineamiento de storage en `setup.ts`/`navigation.test.tsx` también existe y conviene corregirlo, pero no es la causa de esta falla puntual.

### Recomendación antes de seguir

1. **Prioridad alta, bug de producto:** revertir el commit `6b08566` en la parte de `applicationRoutes.tsx`, o mover `/dashboard`, `/instances`, `/users`, `/users/new`, `/users/:userId` y `/auditoria` de vuelta al grupo con `loader: loadProtectedSession`. Si se necesitan para diseño sin backend, usar datos mockeados detrás del guard, no sacarlas de la protección.
2. Corregir `test/front/setup.ts` para también hacer `window.sessionStorage.clear()` en `afterEach`, y actualizar `navigation.test.tsx` para sembrar la sesión en `sessionStorage` (vía las mismas claves que usa `tokenStorage.ts`). Esto es necesario mantenimiento del test, independiente del bug de arriba.
3. No conviene todavía escribir pruebas nuevas para `FRN-07` (selector de instancias) o `FRN-08` (interceptor 403): ningún componente para eso apareció en este pull. La vista `Auditoria.tsx` tampoco tiene contraparte de backend (el endpoint de auditoría es una tarea reservada para la Etapa 3, ver `documentacion/siguientesetapas.md`), así que escribir un test de contrato contra ella sería prematuro; como mucho cabría un smoke test de que la vista renderiza, pero no valida ningún `RF` todavía implementado según `actual.md`.
