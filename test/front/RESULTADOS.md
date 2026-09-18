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
