# Contraste del estado actual con frontend y backend

## Resultado ejecutivo

Al 15/09/2026, los submódulos auditados corresponden a los últimos commits de sus ramas `main`:

| Componente | Commit auditado | Relación con `origin/main` |
|---|---|---|
| Frontend | `17c5a8a` | Sin diferencias |
| Backend | `8d7bce5` | Sin diferencias |

De las once tareas descritas en `actual.md`, cinco están implementadas en código pero todavía requieren una validación funcional de aceptación, tres están completas en su alcance visual o de navegación, dos están en proceso por incompatibilidades entre frontend y backend y la integración punta a punta no está completada.

## Criterio de clasificación

- **Completa en alcance:** el entregable solicitado existe y la validación disponible lo respalda.
- **Implementada, pendiente de validación:** el código existe y compila, pero no hay pruebas automatizadas ni una comprobación funcional suficiente para confirmar el criterio de éxito.
- **En proceso:** existe una parte sustancial, pero falta corregir comportamiento o integración.
- **No completada:** el criterio principal no puede recorrerse con el estado actual.

## Estado por tarea

| ID | Estado | Evidencia y contraste |
|---|---|---|
| `BAC-01` | **En proceso** | Existe PostgreSQL en Docker Compose, migración GORM y creación de un administrador inicial. Sin embargo, el entorno integrado no inicia porque Compose configura `DATABASE_URL` y el backend lee `DB_DSN`; dentro del contenedor intenta conectarse a `localhost:5433`. No se confirmó el criterio completo de levantar la base y conectarse con las credenciales documentadas. |
| `BAC-02` | **Implementada, pendiente de validación** | El backend usa bcrypt con costo 12 para generar y verificar hashes. No existen tests automatizados que demuestren credenciales válidas e inválidas. |
| `BAC-03` | **Implementada, pendiente de validación** | Está registrado `POST /api/auth/login`; consulta PostgreSQL, verifica contraseña y emite un JWT temporal con identificador, rol y organización. El endpoint no pudo probarse dentro del entorno integrado debido al bloqueo de `BAC-01`. |
| `BAC-04` | **Implementada, pendiente de validación** | El handler devuelve `400` para un cuerpo inválido y `401` para errores de autenticación mediante una estructura JSON común. No hay tests HTTP automatizados para confirmar todos los casos. |
| `FRN-01` | **Completa en alcance** | El formulario de login contiene correo, contraseña, opción de sesión, botón y estado de envío. El frontend compila correctamente. |
| `FRN-02` | **Completa en alcance** | Existen validaciones previas al envío y errores asociados a los campos. El formulario evita el envío cuando encuentra errores. |
| `FRN-03` | **Completa en alcance** | Dashboard e Instancias están registradas en el router y accesibles mediante la navegación protegida. Las vistas son estáticas, como permite el entregable. |
| `FRN-04` | **En proceso** | El frontend realiza una petición real y conserva estado temporal, pero su contrato no coincide con el backend: envía `password`, espera `challengeToken` y una estructura de usuario; el backend exige `contrasena` y devuelve `jwtTemporal`, `totpVinculado` y `cambioContrasenaRequerido`. |
| `LOGIN-01` | **Implementada, pendiente de validación** | El backend supera el prototipo en memoria: genera TOTP, cifra el secreto con AES-256-GCM, lo persiste en PostgreSQL, valida códigos y emite tokens definitivos. Falta una prueba funcional completa y no existen tests automatizados. |
| `LOGIN-02` | **En proceso** | Existe una pantalla real de QR/código, estados de carga y mensajes de error. No puede completar la validación porque llama `POST /auth/2fa/setup` con `challengeToken`, mientras el backend expone `GET /api/auth/2fa/qr` y requiere `jwtTemporal` como Bearer; también envía `code` donde el backend exige `codigo`. |
| `LOGIN-03` | **No completada** | El recorrido usuario de prueba -> contraseña -> TOTP -> Dashboard no funciona punta a punta por el fallo de configuración Docker y por los contratos incompatibles de login y 2FA. |

## Evidencia principal

### Backend

- Base, automigración y seed: [`backend/internal/adapters/secondary/postgres/db.go`](../backend/internal/adapters/secondary/postgres/db.go).
- Hash de contraseñas: [`backend/internal/infrastructure/crypto/password.go`](../backend/internal/infrastructure/crypto/password.go).
- Rutas HTTP activas: [`backend/cmd/api/main.go`](../backend/cmd/api/main.go).
- Login y respuestas HTTP: [`backend/internal/adapters/primary/http/auth_handler.go`](../backend/internal/adapters/primary/http/auth_handler.go).
- Login, TOTP y emisión de tokens: [`backend/internal/core/services/auth_service.go`](../backend/internal/core/services/auth_service.go).
- Cifrado y validación TOTP: [`backend/internal/infrastructure/crypto/totp.go`](../backend/internal/infrastructure/crypto/totp.go).

### Frontend

- Formulario de login: [`frontend/centinela/src/components/features/auth/components/LoginForm.tsx`](../frontend/centinela/src/components/features/auth/components/LoginForm.tsx).
- Validaciones: [`frontend/centinela/src/components/features/auth/utils/validateAuthenticationFields.ts`](../frontend/centinela/src/components/features/auth/utils/validateAuthenticationFields.ts).
- Contrato y petición de login: [`frontend/centinela/src/components/features/auth/services/authService.ts`](../frontend/centinela/src/components/features/auth/services/authService.ts).
- Contrato esperado por el cliente: [`frontend/centinela/src/components/features/auth/types/authentication.ts`](../frontend/centinela/src/components/features/auth/types/authentication.ts).
- Servicio 2FA: [`frontend/centinela/src/components/features/2fa/services/2fa.service.ts`](../frontend/centinela/src/components/features/2fa/services/2fa.service.ts).
- Continuación de login y QR: [`frontend/centinela/src/components/features/auth/components/LoginContinuation.tsx`](../frontend/centinela/src/components/features/auth/components/LoginContinuation.tsx).
- Rutas de aplicación: [`frontend/centinela/src/routes/applicationRoutes.tsx`](../frontend/centinela/src/routes/applicationRoutes.tsx).

## Bloqueos que impiden cerrar la fase

### 1. Contrato de login incompatible

| Aspecto | Frontend | Backend |
|---|---|---|
| Contraseña enviada | `password` | `contrasena` |
| Token previo a 2FA | `challengeToken` | `jwtTemporal` |
| Indicador de enrolamiento | `requiresTwoFactorSetup` | `totpVinculado` |
| Respuesta esperada | sesión y usuario completos | datos mínimos del paso pre-2FA |

### 2. Contrato de 2FA incompatible

| Aspecto | Frontend | Backend |
|---|---|---|
| Preparación inicial | `POST /auth/2fa/setup` | `GET /api/auth/2fa/qr` |
| Transporte del token temporal | `challengeToken` en JSON | Bearer token en `Authorization` |
| Campo del código | `code` | `codigo` |
| Respuesta de verificación | sesión con usuario y banderas | `accessToken`, `refreshToken`, `expiresIn` |

### 3. Integración Docker bloqueada

El archivo [`compose.yaml`](../compose.yaml) entrega la conexión como `DATABASE_URL`, pero el backend consulta `DB_DSN`. Como resultado, el backend usa su valor local predeterminado y busca PostgreSQL en `localhost:5433` desde el contenedor.

También debe revisarse el proxy: Nginx elimina el prefijo `/api` al reenviar, mientras el backend registra sus rutas con ese prefijo.

## Validaciones ejecutadas

| Validación | Resultado |
|---|---|
| `go test ./...` | Correcta; valida compilación, pero el backend no contiene archivos de test |
| `corepack pnpm --dir frontend/centinela build` | Correcta |
| `corepack pnpm --dir frontend/centinela lint` | Fallida por dos errores existentes en `use2FA.ts` y `button.tsx` |
| `docker compose up --build --detach --wait` | Fallida: el backend no logra conectarse a PostgreSQL |
| Prueba HTTP punta a punta | No ejecutable porque el stack no alcanzó estado saludable |

## Conclusión operativa

Pueden considerarse entregadas en su alcance actual `FRN-01`, `FRN-02` y `FRN-03`. `BAC-02`, `BAC-03`, `BAC-04` y `LOGIN-01` tienen implementación sustancial, pero todavía necesitan pruebas de aceptación antes de marcarlas como terminadas. `BAC-01`, `FRN-04` y `LOGIN-02` requieren correcciones. `LOGIN-03` permanece bloqueada hasta alinear los contratos y conseguir que el entorno integrado arranque.
