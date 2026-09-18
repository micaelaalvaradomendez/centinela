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
| `BAC-05` | Esquema de usuarios y estado inicial de cuentas nuevas. Usar un `id` propio en `permisos_instancia` en vez de clave primaria compuesta es válido; el problema real es que no hay **ninguna restricción de unicidad** sobre `(usuario_id, vmid_proxmox)`. Se demuestra con un test funcional: enviar `{"vmids": [201, 201]}` persiste **2 filas duplicadas** en vez de 1. |
| `BAC-06` | Alta, listado filtrado, contraseña temporal y desactivación administrativa, validados bajo `/api/users`. Un chequeo real contra `GET /api/admin/users` (el prefijo que exige `actual.md`) devuelve `404`: en [backend/cmd/api/main.go](../../backend/cmd/api/main.go) la variable se llama `admin` pero su `Group()` se monta en `"/"`, no en `"/admin"`. |
| `BAC-06B` | Edición de nombre, correo, rol y estado (`activo`) mediante `PUT /api/users/:id`, más rechazo `403` para `OPERATOR` en detalle, edición y baja. |
| `BAC-09` | Listado de roles `ADMIN` y `OPERATOR` mediante `GET /api/roles` (ruta idéntica a la documentada; esta zona de la API sí coincide con `actual.md`). |
| `BAC-07` | Asignación y reemplazo de VMIDs mediante `PUT /api/users/:id/instances`, lectura en detalle y rechazo para `OPERATOR`. Comparte el bug de duplicados de `BAC-05` porque usa el mismo repositorio. |
| `BAC-08` | Contrato reservado para guard de autorización por recurso y comprobación de que no se llama a Proxmox ante un `403` |
| `BAC-14` | Contrato reservado para inventario normalizado y filtrado de instancias desde Proxmox |

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

`BAC-05` y `BAC-06`/`BAC-06B` tienen pruebas que **fallan intencionalmente** por dos motivos distintos:

- `BAC-05` no es un problema de nomenclatura: es un bug funcional reproducible (VMIDs duplicados generan filas duplicadas) porque no existe ninguna restricción de unicidad, ni PK compuesta ni `UNIQUE`, sobre `(usuario_id, vmid_proxmox)`. Usar un `id` propio en vez de esa PK compuesta es válido por sí solo; lo que falta es la restricción.
- `BAC-06`/`BAC-06B` sí es una cuestión de ruta: falta montar el grupo `admin` bajo `/admin` en `cmd/api/main.go`.

El resto del comportamiento (CRUD, contraseña temporal, RBAC) sí está implementado y validado bajo `/api/users`.

La gestión administrativa se prueba contra las rutas implementadas actualmente: `GET/POST /api/users`, `PUT/DELETE /api/users/:id` y `GET /api/roles`. La documentación de tareas menciona `/api/admin/users`, pero ese prefijo no corresponde al router actual; si se cambia durante el desarrollo, habrá que actualizar las constantes de ruta de esta suite y del frontend.

El hito **Control de Acceso Basado en Recursos** queda preparado en `resource_access_acceptance_test.go`: BAC-07 se ejecuta sobre la implementación actual y BAC-08/BAC-14 aparecen como pruebas omitidas hasta que existan el guard de recurso, `GET /api/instances` y un adaptador Proxmox observable por la suite.
