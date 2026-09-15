# Resultado base de las pruebas frontend

Fecha de ejecución: 15/09/2026.

## Resumen

| Métrica | Resultado |
|---|---:|
| Archivos de prueba | 4 |
| Pruebas totales | 14 |
| Aprobadas | 8 |
| Fallidas | 6 |

## Estado derivado

| Tarea | Resultado | Conclusión |
|---|---|---|
| `FRN-01` | Pruebas aprobadas | El formulario y su estado de carga cumplen el alcance comprobado. |
| `FRN-02` | Pruebas aprobadas | Los campos vacíos y el correo inválido se bloquean antes del envío. |
| `FRN-03` | Pruebas aprobadas | La sesión protege las rutas y permite navegar entre Dashboard e Instancias. |
| `FRN-04` | 2 pruebas fallidas | El payload y la respuesta esperada por el frontend no coinciden con la API Go. |
| `LOGIN-02` | 3 pruebas aprobadas y 3 fallidas | La UI muestra y envía seis caracteres, pero acepta letras y usa un contrato 2FA incompatible. |
| `LOGIN-03` | 1 prueba fallida | Una respuesta válida del backend no permite continuar del login al enrolamiento 2FA. |

## Fallos detectados

1. El login envía `password` y `recordarSesion`; el backend espera `contrasena`.
2. El frontend rechaza la respuesta `{ jwtTemporal, totpVinculado, cambioContrasenaRequerido }` del backend.
3. El enrolamiento llama `POST /auth/2fa/setup`; el backend expone `GET /api/auth/2fa/qr`.
4. El JWT temporal no se envía como Bearer token durante el flujo 2FA.
5. La verificación envía `challengeToken` y `code`; el backend espera `codigo` y obtiene la sesión desde el JWT temporal.
6. El campo TOTP permite letras aunque la tarea exige seis dígitos.

## Alcance

Esta suite valida las tareas con responsabilidad frontend presentes en `documentacion/actual.md`. No confirma `BAC-01` a `BAC-04` ni la implementación interna de `LOGIN-01`; esas tareas requieren una suite backend y pruebas de integración con PostgreSQL.
