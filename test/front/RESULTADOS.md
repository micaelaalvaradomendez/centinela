# Resultado base de las pruebas frontend

Fecha de ejecución: 16/09/2026.

## Resumen

| Métrica | Resultado |
|---|---:|
| Archivos de prueba | 5 |
| Pruebas totales | 32 |
| Aprobadas | 13 |
| Fallidas | 1 |
| Pendientes (`todo`) | 18 |

## Estado derivado

| Tarea | Resultado | Conclusión |
|---|---|---|
| `FRN-01` | Pruebas aprobadas | El formulario y su estado de carga cumplen el alcance comprobado. |
| `FRN-02` | Pruebas aprobadas | Los campos vacíos y el correo inválido se bloquean antes del envío. |
| `FRN-03` | Pruebas aprobadas | La protección de rutas y navegación Dashboard -> Instancias pasan con los encabezados actuales. |
| `FRN-04` | Pruebas aprobadas | El payload usa `password` y la respuesta pre-2FA del backend es aceptada. |
| `LOGIN-02` | 3 aprobadas, 1 fallida | El contrato QR/verificación pasa. El input OTP aún acepta letras. |
| `LOGIN-03` | Pruebas aprobadas | La redirección al enrolamiento 2FA tras un login válido funciona correctamente. |
| `FRN-05` | Pendiente | Criterios `todo` para listado administrativo y control visual por rol. |
| `FRN-06` | Pendiente | Criterios `todo` para alta, contraseña temporal, desactivación y errores. |
| `FRN-06B` | Pendiente | Criterios `todo` para edición y cambio de rol. |
| `FRN-07` | Pendiente | Criterios `todo` para selector de instancias y persistencia de VMIDs. |
| `FRN-08` | Pendiente | Criterios `todo` para Bearer y manejo de `403` sin cerrar sesión. |

## Fallos detectados

1. [test/front/two-factor-form.test.tsx](test/front/two-factor-form.test.tsx): el campo TOTP de [frontend/centinela/src/pages/TwoFactor.tsx](frontend/centinela/src/pages/TwoFactor.tsx) permite escribir caracteres no numéricos.

## Alcance

Esta suite valida las tareas con responsabilidad frontend presentes en [documentacion/actual.md](documentacion/actual.md). Los criterios administrativos y del hito de recursos están preparados como `todo` porque todavía no existe la vista `/admin/users`, el selector de instancias ni la integración frontend de `GET /api/instances`.
