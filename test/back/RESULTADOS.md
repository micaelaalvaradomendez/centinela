# Resultado base de las pruebas backend

Fecha de ejecución: 16/09/2026.

## Resumen

| Validación | Resultado |
|---|---|
| `go vet ./...` | Correcta |
| `go test -short ./...` | Correcta |
| `go test -v -count=1 ./...` | Correcta; pruebas implementadas aprobadas y BAC-08/BAC-14 omitidas explícitamente |
| Limpieza de contenedores | Correcta; Docker Compose cerró y eliminó volúmenes limpiamente |

## Estado derivado

| Tarea | Resultado | Conclusión |
|---|---|---|
| `BAC-01` | Aprobada | La base aislada inicia, GORM crea el esquema, existe el administrador de prueba y las tablas están correctamente estructuradas. |
| `BAC-02` | Aprobada | La contraseña se persiste como bcrypt y el login distingue credenciales válidas e inválidas. |
| `BAC-03` | Aprobada | El login consulta PostgreSQL y emite un JWT temporal HS256 válido. |
| `BAC-04` | Aprobada | Los casos incompleto e inválido devuelven `400`/`401` con JSON unificado. |
| `LOGIN-01` | Aprobada | El backend genera QR, cifra el secreto y valida códigos TOTP inválidos, vencidos y vigentes. |
| `LOGIN-03` | Aprobada | El flujo backend completo emite tokens finales y protege las rutas. |
| `BAC-07` | Parcial aprobada | Se asignan y reemplazan VMIDs mediante `PUT /api/users/:id/instances`, se leen en el detalle y un `OPERATOR` recibe `403`; falta el contrato documentado `/permissions`. |
| `BAC-08` | Pendiente | Test reservado; aún no existe guard por recurso ni operación de instancia observable. |
| `BAC-14` | Pendiente | Test reservado; aún no existe `GET /api/instances` ni integración Proxmox verificable. |

## Interpretación

El commit `4e8e9c4` del backend modificó el contrato de `POST /api/auth/login` (`LoginRequest`) para exigir `password` en lugar de `contrasena`. La suite fue actualizada y ahora valida correctamente todo el recorrido implementado con Docker Compose y PostgreSQL reales. El hito de recursos tiene cobertura parcial ejecutable y dos pruebas omitidas hasta que se implementen el guard y el inventario.
