# Informe de estado actual — login, 2FA e integración

Fecha: 15/09/2026.
Fuentes: `documentacion/actual.md`, `test/back` (suite Go, ejecución real contra PostgreSQL) y `test/front` (suite Vitest sobre el código real del frontend).

Commits auditados (coinciden con `origin/main` de cada submódulo):

| Componente | Commit |
|---|---|
| `frontend` | `9001530` |
| `backend` | `4919a88` |

## 1. Resumen ejecutivo

- El **backend cumple** todas las tareas `BAC-01`–`BAC-04` y `LOGIN-01`: la suite `test/back` lo valida con Docker Compose real, PostgreSQL real y HTTP real. `LOGIN-03` está aprobada solo del lado backend.
- El **frontend cumple** `FRN-01`, `FRN-02` y `FRN-03` (formulario, validaciones de cliente, navegación protegida).
- `FRN-04`, `LOGIN-02` y `LOGIN-03` **fallan** porque el frontend y el backend no hablan el mismo contrato HTTP. No es un problema de lógica interna de ninguno de los dos lados: es un desacople entre lo que envía/espera cada submódulo.

## 2. Estado por tarea

| ID | Estado | Evidencia |
|---|---|---|
| `BAC-01` | ✅ Completa (validado por `test/back`) | El compose aislado de `test/back` levanta PostgreSQL, GORM automigra el esquema (tablas `organizaciones`, `usuarios`, `sesiones_activas`) y crea el admin seed con roles `ADMIN`/`OPERATOR`. |
| `BAC-02` | ✅ Completa | Hash bcrypt verificado; login distingue credenciales válidas e inválidas en pruebas reales. |
| `BAC-03` | ✅ Completa | `POST /api/auth/login` consulta PostgreSQL y emite JWT temporal HS256 con `id`, `role`, `org`. |
| `BAC-04` | ✅ Completa | `400` para payload incompleto y `401` para credenciales inválidas, con JSON unificado. |
| `FRN-01` | ✅ Completa | Formulario de login (usuario/correo, contraseña, envío, estado de carga) validado por `test/front`. |
| `FRN-02` | ✅ Completa | Validaciones de cliente bloquean campos vacíos y correo inválido antes de enviar. |
| `FRN-03` | ✅ Completa | Rutas protegidas y navegación Dashboard/Instancias funcionan. |
| `FRN-04` | ❌ En proceso — **bloqueada por contrato** | El frontend envía `password`/`recordarSesion` y espera `challengeToken`; el backend exige `contrasena` y devuelve `jwtTemporal`, `totpVinculado`, `cambioContrasenaRequerido`. |
| `LOGIN-01` | ✅ Completa | El backend genera TOTP, cifra el secreto (AES-256-GCM), lo persiste y valida códigos correctamente. |
| `LOGIN-02` | ❌ En proceso — **bloqueada por contrato** | La pantalla de TOTP existe, pero llama `POST /auth/2fa/setup` (el backend expone `GET /api/auth/2fa/qr`), envía `challengeToken` en el body en vez de `Authorization: Bearer`, y usa `code` donde el backend espera `codigo`. Además el input acepta letras, aunque el criterio pide 6 dígitos. |
| `LOGIN-03` | ⚠️ Parcial — **backend aprobado, integración no** | El circuito completo (credenciales → TOTP → Dashboard) funciona de punta a punta en `test/back` (solo HTTP), pero no puede completarse desde el navegador por los desacoples de `FRN-04` y `LOGIN-02`. |

## 3. Errores a corregir

### 3.1 Contrato de login (`FRN-04`)

| Campo | Frontend envía/espera | Backend expone |
|---|---|---|
| Contraseña | `password` | `contrasena` |
| Token previo a 2FA | `challengeToken` | `jwtTemporal` |
| Indicador de enrolamiento | `requiresTwoFactorSetup` | `totpVinculado` |
| Respuesta esperada | sesión y usuario completos | datos mínimos del paso pre-2FA |

**Corrección pendiente:** alinear `frontend/centinela/src/components/features/auth/services/authService.ts` y sus tipos (`authentication.ts`) con los nombres de campo reales del backend, o exponer un adaptador único que traduzca el contrato en un solo punto.

### 3.2 Contrato de 2FA (`LOGIN-02`)

| Aspecto | Frontend | Backend |
|---|---|---|
| Preparación | `POST /auth/2fa/setup` | `GET /api/auth/2fa/qr` |
| Transporte del token temporal | `challengeToken` en JSON | `Authorization: Bearer` |
| Campo del código | `code` | `codigo` |
| Respuesta de verificación | sesión con usuario y banderas | `accessToken`, `refreshToken`, `expiresIn` |
| Validación de input | acepta letras | debe exigir 6 dígitos numéricos |

**Corrección pendiente:** ajustar `2fa.service.ts` y `LoginContinuation.tsx` al método/URL/headers reales, y agregar validación numérica estricta en el input OTP del frontend.

## 4. Validaciones ejecutadas

| Suite | Resultado |
|---|---|
| `test/back` (`go vet`, `go test -short`, `go test -v -count=1`, contra Docker Compose real) | ✅ Todo aprobado; Docker quedó limpio al finalizar |
| `test/front` (Vitest, 4 archivos / 14 pruebas) | 8 aprobadas, 6 fallidas (`authentication-contract.test.ts`, `two-factor-form.test.tsx`) por los contratos descritos arriba |

## 5. Próximos pasos sugeridos

1. Unificar el contrato de login (`FRN-04`) entre frontend y backend.
2. Unificar el contrato de 2FA (`LOGIN-02`), incluida la validación numérica del input.
3. Re-ejecutar `test/front` y `test/back` para confirmar `LOGIN-03` de punta a punta.
