# Resultados de las pruebas backend

**Fecha:** 30/09/2026. **Backend:** submódulo en `e6e7dc5` (último commit de `main`). **Comando:** `go test -v -count=1 ./...`. 2 corridas con el mismo resultado (~2,5 min) y sin contenedores residuales.

| Métrica | Valor |
|---|---:|
| Casos | 43 |
| Aprueban | 30 |
| Fallan | **13** |
| Omitidos | 0 |

## Por caso

| Caso | Tarea | Resultado |
|---|---|---|
| BAC-01 a BAC-12, BAC-09, BAC-05/06/06B, prefijo `/api/admin/users`, LOGIN-03 | regresión | ✅ 12/12 |
| LOGIN-04 circuito desde Go | regresión | ✅ |
| **LOGIN-04 front ↔ back** (`TestLOGIN04IntegracionFrontBack`) | LOGIN-04 | ❌ 9/10 pasos (falla el paso 7 por FIX-28) |
| INF-05 CORS | regresión | ✅ |
| **FIX-31 TLS de Nginx versionado** | FIX-31 | ❌ |
| **INF-06A Redis local** | INF-06A | ❌ |
| **INF-08B credenciales SMTP** | INF-08B | ❌ |
| **BAC-16B envío por SMTP** | BAC-16B | ❌ |
| **BAC-17A adaptador de Redis** | BAC-17A | ❌ |
| **BAC-17B una sesión = un registro** | BAC-17B | ❌ |
| **BAC-18B índice parcial y particiones** | BAC-18B | ❌ |
| BAC-19, BAC-20, BAC-16, BAC-21, BAC-18, FIX-17, BAC-13, BAC-15 | regresión | ✅ 8/8 |
| BAC-07, FIX-16/BAC-08 (incluido `INSTANCE_PROTECTED`), BAC-14, SEC-04 (2) | regresión | ✅ 5/5 |
| **BAC-21B** (3 casos) | BAC-21B | ❌ 0/3 |
| **BAC-21C** (2 casos) | BAC-21C | ❌ 0/2 |
| BAC-17, SEC-01, integración SEC-01/SEC-02, FIX-08 | regresión | ✅ 4/4 |
| **FIX-28 integración del logout del frontend** | FIX-28 | ❌ |

## Fallos y causa

| Caso | Causa |
|---|---|
| FIX-31 | No hay ninguna configuración de Nginx versionada con `listen 443 ssl` |
| INF-06A | Falta `127.0.0.1:6379:6379` en el compose; `REDIS_ADDR=redis:6379`; las contraseñas del compose y de `.env.example` no coinciden |
| INF-08B | Las 6 variables SMTP no están en `backend/.env.example` |
| BAC-16B | `cmd/api/main.go` siempre instancia `email.NewMockEmailService()` |
| BAC-17A | Falta `go-redis`, el adaptador y los puertos; el backend no se conecta a Redis |
| BAC-17B | 13 filas activas por sesión; quedan 11 activas después del logout |
| BAC-18B | Sin índice parcial, sin particiones y sin índices compuestos |
| BAC-21B | Faltan los campos nuevos en `GET /instances`; `/status/:action` → 404; `DELETE /instances/:vmid` → 405 |
| BAC-21C | `/events/ticket` y `/events` → 404 |
| FIX-28 (y LOGIN-04 paso 7) | El frontend envía el logout sin Bearer (`skipAuthorization: true`, `deb59cb`): `401 MISSING_TOKEN` y la sesión sigue activa |

## Lo que la suite no cubre
- `INF-06B` e `INF-08A`: configuración del servidor.
- Conectividad hacia un SMTP y un Proxmox reales. La comparación entre el simulador y el Proxmox real se hizo a mano; ver `../informe.md` §4.
- La purga horaria y el tiempo de consulta de BAC-18B.
