# Resultado base de las pruebas backend

Fecha de ejecución: 18/09/2026 (revisión en detalle de los tests contra `documentacion/actual.md`).

## Resumen

| Validación | Resultado |
|---|---|
| `go vet ./...` | Correcta |
| `go test -short ./...` | Correcta |
| `go test -v -count=1 ./...` | **3 fallas reales**: un bug funcional de duplicados en `BAC-05`/`BAC-07` y un desvío de ruta en `BAC-06`/`BAC-06B`; el resto de las pruebas implementadas pasa y `BAC-08`/`BAC-14` quedan omitidas explícitamente |
| Limpieza de contenedores | Correcta; Docker Compose cerró y eliminó volúmenes limpiamente |

## Estado derivado

| Tarea | Resultado | Conclusión |
|---|---|---|
| `BAC-01` | Aprobada | La base aislada inicia, GORM crea el esquema, existe el administrador de prueba y las tablas están correctamente estructuradas. |
| `BAC-02` | Aprobada | La contraseña se persiste como bcrypt y el login distingue credenciales válidas e inválidas. |
| `BAC-03` | Aprobada | El login consulta PostgreSQL y emite un JWT temporal HS256 válido. |
| `BAC-04` | Aprobada | Los casos incompleto e inválido devuelven `400`/`401` con JSON unificado. |
| `LOGIN-01` | Aprobada | El backend genera QR, cifra el secreto y valida códigos TOTP inválidos, vencidos y vigentes. |
| `BAC-10` | **Fallida** | El QR y secreto se generan, pero una cuenta ya vinculada puede pedir otro QR por `GET /api/auth/2fa/qr` sin pasar por reset/relink. |
| `BAC-11` | Parcial | La persistencia cifrada y el login posterior se ejercitan en la prueba nueva; la verificación queda bloqueada por el fallo previo de BAC-10. |
| `LOGIN-03` | Aprobada | El flujo backend completo emite tokens finales y protege las rutas. |
| `BAC-05` | **Fallida — bug funcional, no de nomenclatura** | El alta de usuarios y sus campos de seguridad funcionan. Usar un `id` propio como PK en `permisos_instancia` (en vez de una compuesta) es válido; el problema real es que **no hay ninguna restricción de unicidad** sobre `(usuario_id, vmid_proxmox)` (ni PK compuesta ni `UNIQUE`). Se demostró en vivo: al enviar `{"vmids": [201, 201]}` se persistieron **2 filas** para el mismo par usuario/instancia en lugar de 1. |
| `BAC-06` | Parcial — **prueba de contrato fallida** | Alta, listado filtrado y desactivación funcionan correctamente bajo `/api/users`. `GET /api/admin/users` devuelve `404`: en `cmd/api/main.go` la variable del grupo se llama `admin` pero se monta con `Group("/")` en vez de `Group("/admin")`, a diferencia de `/api/auth/*` y `/api/roles`, que sí calzan con `actual.md`. |
| `BAC-06B` | Aprobada (bajo `/api/users`) | `PUT /api/users/:id` actualiza nombre, correo, rol y estado (`activo`); un `OPERATOR` recibe `403` en detalle, edición y baja. Mismo desvío de prefijo que `BAC-06`, sin el bug de duplicados de `BAC-05`. |
| `BAC-09` | Aprobada | `GET /api/roles` coincide exactamente con la ruta documentada en `actual.md`. |
| `BAC-07` | Parcial — **hereda el bug de `BAC-05`** | Se asignan y reemplazan VMIDs mediante `PUT /api/users/:id/instances` y se leen en el detalle; un `OPERATOR` recibe `403`. Usa el mismo repositorio sin restricción de unicidad, y falta el contrato documentado `/permissions`. |
| `BAC-08` | Pendiente | Test reservado; aún no existe guard por recurso ni operación de instancia observable. |
| `BAC-14` | Pendiente | Test reservado; aún no existe `GET /api/instances` ni integración Proxmox verificable. |

## Interpretación

El commit `4e8e9c4` del backend modificó el contrato de `POST /api/auth/login` (`LoginRequest`) para exigir `password` en lugar de `contrasena`. La suite fue actualizada y ahora valida correctamente todo el recorrido implementado con Docker Compose y PostgreSQL reales. El hito de recursos tiene cobertura parcial ejecutable y dos pruebas omitidas hasta que se implementen el guard y el inventario.

Pull del 18/09/2026 (`4536776` -> `caf71c2`): normaliza el email a minúsculas en login/alta/edición y elimina el endpoint duplicado `DELETE /account/sessions/current` (el logout real, `POST /api/auth/logout`, no cambió). Ningún archivo de `BAC-05` a `BAC-14` fue modificado en ese pull.

Revisión del 18/09/2026: al preguntarnos si usar un `id` propio en vez de la clave compuesta literal era válido, se reemplazó el chequeo de "tiene que llamarse así" por uno funcional: ¿existe alguna restricción de unicidad sobre `(usuario_id, vmid_proxmox)`, sea PK o `UNIQUE`? No existe ninguna, y se demostró con un test que envía VMIDs duplicados y verifica cuántas filas quedaron persistidas: **2, no 1**. Esto confirma que el hallazgo no es una cuestión de nomenclatura sino un bug funcional real. El desvío de ruta (`/api/admin/users` vs `/api/users`) sigue siendo un tema aparte, de contrato, no de lógica.
