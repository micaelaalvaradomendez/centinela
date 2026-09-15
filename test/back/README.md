# Pruebas de aceptación del backend

Esta suite contrasta las tareas backend y la parte backend de integración descritas en `documentacion/actual.md`. Usa el stack real: Go, Gin, PostgreSQL 16 y Docker Compose.

## Cobertura

| Tarea | Comportamiento validado |
|---|---|
| `BAC-01` | PostgreSQL, esquema generado, rol y usuario administrador seed |
| `BAC-02` | Hash bcrypt, ausencia de contraseña plana y rechazo de contraseña inválida |
| `BAC-03` | Login contra PostgreSQL, firma HS256 y claims del JWT temporal |
| `BAC-04` | Errores JSON `400` y `401` con código y mensaje |
| `LOGIN-01` | QR, secreto cifrado, códigos inválidos/vencidos y vinculación con TOTP válido |
| `LOGIN-03` | Emisión de tokens finales y protección de rutas en la parte backend |

## Requisitos

- Docker con `docker compose`.
- Go compatible con el módulo de pruebas.
- Puertos locales `15433` y `18080` disponibles.

## Ejecutar

Desde la raíz del repositorio:

```bash
cd test/back
go test -v -count=1 ./...
```

La suite construye los contenedores, usa una base aislada y elimina contenedores y volumen al terminar.

Para validar solamente que los tests compilan, sin Docker:

```bash
cd test/back
go test -short ./...
```

## Alcance

El éxito de `LOGIN-03` confirma el recorrido del backend desde credenciales hasta un access token protegido por TOTP. El recorrido completo en navegador también requiere que pase la suite `test/front`, porque actualmente los contratos HTTP de ambos componentes no coinciden.
