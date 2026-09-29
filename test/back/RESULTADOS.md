# Resultados de las pruebas backend

**Fecha:** 29/09/2026. **Backend:** submódulo en `9554efa` (último commit de `main`). **Comando:** `go test -v -count=1 ./...`. Se hicieron 2 corridas con el mismo resultado (~2 min) y no quedaron contenedores residuales.

| Métrica | Valor |
|---|---:|
| Casos | 38 |
| Aprueban | 32 |
| Fallan | **5** |
| Omitidos | 1 |

## Por caso

| Caso | Tarea | Resultado |
|---|---|---|
| BAC-01, BAC-02, BAC-03, BAC-04, LOGIN-01, BAC-10/11, BAC-12, BAC-09, BAC-05/06/06B, prefijo `/api/admin/users`, LOGIN-03 | regresión | ✅ 12/12 |
| LOGIN-04 circuito completo | regresión | ✅ |
| INF-05 CORS con lista blanca | INF-05 | ✅ |
| INF-05 TLS en Nginx | INF-05 | ⏭️ omitido: la configuración de Nginx del servidor no está versionada |
| INF-06 Redis del `docker-compose.yml` del backend | INF-06 | ✅ |
| **INF-08 variables SMTP** | INF-08 | ❌ |
| **BAC-16B envío por SMTP** | BAC-16B | ❌ |
| **BAC-17B una sesión = un registro** | BAC-17B | ❌ |
| **BAC-18B índice parcial y particiones** | BAC-18B | ❌ |
| BAC-19, BAC-20, BAC-16, BAC-21, FIX-17, BAC-13, BAC-15 | regresión | ✅ 7/7 |
| BAC-18 auditoría append-only (trigger + `REVOKE`, commits `ad52485` y `9554efa`) | BAC-18 / FIX-23 | ✅ |
| BAC-07, FIX-16/BAC-08, BAC-14, SEC-04 | regresión | ✅ 4/4 |
| SEC-04: el payload anterior `{ vmids }` se rechaza con 400 (contrato `{ permisos }`) | SEC-04 (FIX-26 descartado) | ✅ |
| BAC-17, SEC-01, integración SEC-01/SEC-02 y FIX-08 | regresión | ✅ 4/4 (FIX-08 con 4 subcasos) |
| **FIX-28 FRN-13 integración: logout tal como lo envía el frontend** | FRN-13 (regresión del frontend) | ❌ |

## Fallos

| Caso | Mensaje | Causa |
|---|---|---|
| INF-08 | `EMAIL_PROVIDER / SMTP_HOST / SMTP_PORT / SMTP_USER / SMTP_PASS / SMTP_FROM no está documentada…` | No figuran en `backend/.env.example` ni en `backend/docker-compose.yml` |
| BAC-16B | `con EMAIL_PROVIDER=smtp el alta no envió ningún correo SMTP…` | `cmd/api/main.go:85` siempre instancia `email.NewMockEmailService()` |
| BAC-17B | `se crearon 13 filas y hay 13 activas` / `tras el logout quedaron 11 registros activos` | Login, 2FA y cada refresh insertan filas nuevas, y el logout solo desactiva el último access y el último refresh |
| BAC-18B | `falta el índice parcial…`, `relkind actual "r"`, `falta la partición…`, `falta el índice compuesto…` | No hay migración |
| FRN-13 integración | `esperado 204, recibido 401: MISSING_TOKEN` / `el access token … recibido 200 (la sesión sigue activa)` | El frontend envía el logout sin Bearer (`skipAuthorization: true`, commit `deb59cb`) |

## Lo que la suite no cubre
- `502 PROXMOX_UNAVAILABLE`, porque el stub siempre responde.
- **INF-05:** el TLS de Nginx en el servidor.
- **INF-06:** la red `vmbr1`.
- **INF-08:** la conectividad hacia el SMTP real.
- **BAC-18B:** la purga horaria y el tiempo de consulta.
