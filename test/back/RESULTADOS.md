# Resultado base de las pruebas backend

Fecha de ejecución: 15/09/2026.

## Resumen

| Validación | Resultado |
|---|---|
| `go vet ./...` | Correcta |
| `go test -short ./...` | Correcta |
| `go test -v -count=1 ./...` | Correcta |
| Limpieza de contenedores | Correcta; no quedaron recursos activos |

## Estado derivado

| Tarea | Resultado | Conclusión |
|---|---|---|
| `BAC-01` | Aprobada | La base aislada inicia, GORM crea el esquema, existe el administrador de prueba y la API expone los roles `ADMIN` y `OPERATOR`. |
| `BAC-02` | Aprobada | La contraseña se persiste como bcrypt y el login distingue credenciales válidas e inválidas. |
| `BAC-03` | Aprobada | El login consulta PostgreSQL y emite un JWT temporal HS256 válido con identidad, rol y organización. |
| `BAC-04` | Aprobada | Los casos incompleto e inválido devuelven `400`/`401` y una estructura JSON unificada. |
| `LOGIN-01` | Aprobada | El backend genera QR y secreto, persiste el secreto cifrado, rechaza códigos inválidos o vencidos y acepta un TOTP vigente. |
| `LOGIN-03` | Aprobada parcialmente | La parte backend completa credenciales -> TOTP -> access token y protege rutas. El recorrido Front + Back continúa bloqueado por los contratos detectados en `test/front`. |

## Interpretación

Estas pruebas confirman el comportamiento del commit backend auditado mediante HTTP y PostgreSQL reales. No sustituyen las pruebas frontend ni demuestran por sí solas que el usuario pueda completar el flujo desde el navegador.
