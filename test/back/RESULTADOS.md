# Resultados de las pruebas backend

**Fecha:** 25/09/2026. **Backend:** submódulo actualizado a `eb0c9af` (último commit de `main`). **Comando:** `go test -v -count=1 ./...`. Se hicieron 2 corridas con el mismo resultado (~80 s) y no quedaron contenedores residuales.

| Métrica | Valor |
|---|---:|
| Casos | 30 |
| Aprueban | 28 |
| Fallan | **2** |

## Por caso

| Caso | Tarea | Resultado |
|---|---|---|
| BAC-01, BAC-02, BAC-03, BAC-04, LOGIN-01, BAC-10/11, BAC-12, BAC-09, BAC-05/06/06B, prefijo `/api/admin/users`, LOGIN-03 | regresión | ✅ 12/12 |
| LOGIN-04 circuito completo | regresión | ✅ |
| BAC-19, BAC-20, BAC-16, BAC-21, BAC-13, BAC-15 | regresión | ✅ 6/6 |
| **BAC-18 auditoría append-only** | BAC-18 / **FIX-23** | ❌ |
| FIX-17 contratos canónicos de contraseñas | FIX-17 | ✅ |
| BAC-07 permisos (contrato `{ permisos }`) | regresión | ✅ |
| FIX-16 / BAC-08 guard por recurso | regresión | ✅ |
| BAC-14 inventario normalizado y filtrado | regresión | ✅ |
| SEC-04 niveles FULL_ACCESS / READ_ONLY | SEC-04 | ✅ |
| **FIX-26 SEC-04 retrocompatibilidad `{ vmids }`** | SEC-04, entregable 4 → **FIX-26** | ❌ |
| BAC-17 logout atómico (con cookie) | regresión | ✅ |
| SEC-01 refresh solo en cookie HttpOnly | SEC-01 | ✅ |
| SEC-01 / SEC-02 integración con el frontend | SEC-01 | ✅ |
| FIX-08 estructura de error | regresión | ✅ (4 subcasos) |

## Fallos

### FIX-23 / BAC-18: la auditoría no es append-only
```
password_recovery_acceptance_test.go:227: la base permitió UPDATE sobre auditoria con el usuario de la aplicación; la tabla no es append-only
password_recovery_acceptance_test.go:230: la base permitió DELETE sobre auditoria con el usuario de la aplicación; la tabla no es append-only
password_recovery_acceptance_test.go:235: la base permitió TRUNCATE sobre auditoria con el usuario de la aplicación (FIX-23)
```
No hay ningún trigger, `REVOKE` ni regla en `backend/`. La solución está detallada en FIX-23.

### FIX-26 (SEC-04, entregable 4): no hay retrocompatibilidad con `{ vmids }`
```
resource_access_acceptance_test.go:231: PUT permissions con el payload anterior { vmids: [103] }: esperado 204, recibido 400:
  {"errorCode":"INVALID_REQUEST", "message":"Se requiere el campo permisos (array de {vmid, nivelAcceso})."}
```
`asignarPermisosRequest` (`internal/adapters/primary/http/user_handler.go`) solo acepta `permisos` y, además, lo marca con `binding:"required"`.

## Lo que la suite no cubre
- `502 PROXMOX_UNAVAILABLE`, porque el stub de Proxmox siempre responde.
- El envío real de correo, que está fuera de alcance según BAC-16.
