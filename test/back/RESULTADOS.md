# Resultados de las pruebas backend

**Fecha:** 01/10/2026 (2ª verificación). **Backend:** submódulo en `44a2339` (último commit de `main`). **Comando:** `go test -v -count=1 ./...`. 2 corridas completas con el mismo resultado (~2,5 min) y sin contenedores residuales.

| Métrica | Valor |
|---|---:|
| Casos | 50 |
| Aprueban | **43** |
| Fallan | **6** |
| Omitidos | 1 (DELETE de BAC-21B: el endpoint todavía no existe) |

## Por caso

| Caso | Tarea | Resultado |
|---|---|---|
| BAC-01 a BAC-12, BAC-09, BAC-05/06/06B, prefijo `/api/admin/users`, LOGIN-03 | regresión | ✅ 12/12 |
| LOGIN-04 circuito desde Go y LOGIN-04 front ↔ back (10/10 pasos) | regresión | ✅ 2/2 |
| INF-05 CORS, INF-06A, BAC-16B, BAC-17A, BAC-17B | regresión | ✅ 5/5 |
| FIX-31 borde TLS (`docker/nginx-edge.conf` + servicio `edge`) | FIX-31 | ✅ |
| INF-08B credenciales reales contra Brevo | INF-08B | ✅ |
| **BAC-18B índice parcial y particiones** | BAC-18B | ❌ |
| **BAC-18B purga horaria** (nuevo) | BAC-18B | ❌ |
| BAC-19, BAC-20, BAC-16, BAC-21, BAC-18, FIX-17, BAC-13, BAC-15 | regresión | ✅ 8/8 |
| BAC-07, FIX-16/BAC-08, BAC-14, SEC-04 (2) | regresión | ✅ 5/5 |
| **FIX-37 nivel de acceso en `GET /account/profile`** | FIX-37 | ❌ |
| BAC-21B energía con FULL_ACCESS + `tareas_asincronas` | BAC-21B | ✅ |
| **BAC-21B campos nuevos en `GET /instances`** | BAC-21B / FIX-39 | ❌ |
| **BAC-21B auditoría de energía con upid** | BAC-21B / FIX-39 | ❌ |
| **FIX-39 shutdown y reboot con FULL_ACCESS** (nuevo) | FIX-39 | ❌ (404) |
| BAC-21B DELETE solo ADMIN | BAC-21B | ⏭️ omitido (lo crea BAC-24B) |
| BAC-21C (2 casos) | regresión | ✅ 2/2 |
| **INF-07B `/api/events` a través del borde Nginx** (nuevo, Etapa 1) | INF-07B | ✅ (`TASK_FINISHED` en ~1 s) |
| BAC-17, SEC-01, integración SEC-01/SEC-02, FIX-28, FIX-08 | regresión | ✅ 5/5 |
| FIX-35 simulador: `prefix` string y LXC apagado `200 {"data": null}` (2, nuevos) | FIX-35 | ✅ 2/2 |

## Fallos y causa

| Caso | Causa |
|---|---|
| BAC-18B (2) | Sin índice parcial, sin particiones trimestrales ni índice `(fecha_hora, accion, resultado)`, y sin purga horaria. `(usuario_id, fecha_hora)` ya existe |
| FIX-37 | `UsuarioDetalleDTO` solo tiene `instanciasPermitidas` |
| BAC-21B / FIX-39 | Faltan los campos nuevos en `GET /instances`; las acciones de energía no se escriben en `auditoria`; `shutdown` y `reboot` no existen (404) |

## Cambios en la suite

- **Nuevo `simulador_proxmox_acceptance_test.go` (FIX-35):** compila `backend/cmd/proxmox-simulador` en un directorio temporal y lo levanta en un puerto libre.
- **BAC-18B:** caso nuevo de purga horaria, verificado en el código fuente. El de índices y particiones se reescribió contra el criterio, sin nombres exactos, y se validó con un esquema correcto en un Postgres temporal.
- **FIX-39:** caso nuevo para `shutdown` y `reboot`; el `GET` exige el valor real de `nivelAcceso` (el resto puede ir en `null`).
- **BAC-21B:** reescrito contra el criterio de éxito. Acepta `/status/:action` o `/start` y `/stop`, la auditoría puede tener cualquier forma con el `upid`, y el `DELETE` se omite si no existe.
- **FIX-31:** el servicio `edge` (nginx con `docker/nginx-edge.conf` y un certificado de `edge-tls/`) se prueba en funcionamiento por HTTPS y HTTP.
- **INF-07B (nuevo, Etapa 1):** verifica `docker/nginx-edge.conf` para el stream: `proxy_http_version 1.1`, `proxy_buffering off` y `proxy_read_timeout` de 5 min o más. También verifica el circuito real por el servicio `edge` en HTTPS: ticket, `GET /api/events`, `start` de la 102 y `TASK_FINISHED` con el mismo `tareaId` en menos de 8 s. La verificación de la configuración se probó contra variantes incorrectas, y las rechaza.
- **Stub de Proxmox:** responde `GET /nodes/{node}/tasks/{upid}/status` con la tarea terminada (`stopped`, `exitstatus: OK`). Antes no existía, así que el seguimiento de UPID nunca publicaba `TASK_FINISHED` en la suite. No cambió el resultado de ningún otro caso.
- **INF-08B:** reintenta ante errores de red o DNS, y ante `525 Unauthorized IP` (depende de la red) omite el caso con el motivo.

## Lo que la suite no cubre

- El tiempo de consulta (< 20 ms) de BAC-18B.
- INF-07B en el servidor real (CT 103): el caso valida la referencia versionada; el servidor se valida en `INT-03`.
- Que el stream siga abierto más allá del `proxy_read_timeout`: se verifica el valor configurado (el backend manda `: ping` cada 25 s), no una espera de minutos.
- `instancesSummary`: va en `GET /api/node/status` (`BAC-22B`), fuera de BAC-21B.
