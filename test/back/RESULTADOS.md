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

## Pull del 20/09/2026 (`caf71c2` -> `52ebce1`) — regresión crítica detectada

Se trajeron 3 commits nuevos (`3052a02`, `adcf261`, `52ebce1`) que dicen resolver los dos hallazgos de la revisión anterior. Al re-ejecutar `go test -v -count=1 ./...` la suite **empeoró**: pasó de 3 fallas a 8 fallas, incluyendo tareas que antes aprobaban (`BAC-03`, `LOGIN-01`, `BAC-09`, `LOGIN-03`).

### Causa raíz

El commit `3052a02` agregó `uniqueIndex:idx_usuario_vmid` a `PermisoInstancia.UsuarioID` y `PermisoInstancia.VmidProxmox` (correcto, esa combinación es la que se necesitaba), pero **el mismo tag `uniqueIndex:idx_usuario_vmid` se copió también al campo `UsuarioID` de `SesionActiva`, `Auditoria` y `Notificacion`** en `internal/core/domain/models.go`, cada uno como único campo del grupo. Se verificó contra la base real (`\d sesiones_activas`, `\d permisos_instancia`, `\d auditoria`, `\d notificaciones`):

- `sesiones_activas` terminó con un índice **`UNIQUE, btree (usuario_id)`** de una sola columna. Efecto: un usuario no puede tener más de una fila en `sesiones_activas` en toda la vida de la cuenta. Se reprodujo con `curl` puro (sin la suite): el primer login funciona, pero **el propio flujo de verificación TOTP del primer login ya falla** (`TOTP_FAILED: ... duplicate key value violates unique constraint "idx_usuario_vmid"`) porque el login crea una sesión temporal y la verificación intenta crear la sesión de refresh, y la segunda inserción choca contra la primera.
- `permisos_instancia` **no tiene ningún índice único** (solo la PK). Como los nombres de índice son únicos por esquema en PostgreSQL, y `idx_usuario_vmid` ya había sido tomado por `sesiones_activas` durante el `AutoMigrate`, la creación del índice en `permisos_instancia` fue silenciosamente omitida. **El bug original de `BAC-05` (VMIDs duplicados) sigue sin solución real**, pese a que el commit dice arreglarlo.
- `auditoria` y `notificaciones` tampoco recibieron el índice (mismo motivo de colisión de nombre), así que no hay evidencia de que estas tablas estén rotas hoy, pero el patrón es una señal de que el mismo copy-paste puede repetirse.

La mitigación en `ReemplazarPermisos` (deduplicar `vmids` en memoria antes de insertar) sí funciona para pedidos individuales, pero no reemplaza la restricción de base de datos: dos peticiones separadas todavía pueden crear filas duplicadas.

### Impacto

**Ningún usuario puede completar el login + verificación TOTP en una base nueva.** Esto bloquea `BAC-03`, `LOGIN-01`, `BAC-10`/`BAC-11`, `BAC-09`, `LOGIN-03` y cualquier prueba de `BAC-05`/`BAC-06`/`BAC-06B`/`BAC-07` que dependa de un `accessToken` real, no por un problema de las pruebas sino por esta regresión.

El fix de ruta de `BAC-06`/`BAC-06B` (mover el grupo `admin` a `Group("/admin", ...)`) sí se aplicó correctamente: se confirmó manualmente que `GET /api/users` ahora devuelve `404` y `GET /api/admin/users` es la ruta montada. No se pudo validar el comportamiento completo (200 para ADMIN, 403 para OPERATOR) porque no se puede obtener un `accessToken` válido mientras la regresión de sesiones siga presente.

### Recomendación antes de seguir

1. En `internal/core/domain/models.go`, renombrar el índice de `PermisoInstancia` a algo específico, por ejemplo `uniqueIndex:idx_permiso_instancia_usuario_vmid`, y **revertir el tag `uniqueIndex:idx_usuario_vmid` en `SesionActiva`, `Auditoria` y `Notificacion`** a `index` simple (como estaba antes de `3052a02`).
2. Re-ejecutar `go test -v -count=1 ./...` recién después de ese fix; hoy la suite no es representativa de ningún estado nuevo real, solo confirma la regresión.
3. No implementar todavía las pruebas reservadas de `BAC-07` (`GET /permissions`), `BAC-08` o `BAC-14`: el backend no agregó ningún endpoint de instancias, guard de recurso ni inventario de Proxmox en este pull; siguen sin contrato real contra el cual escribir una prueba.
