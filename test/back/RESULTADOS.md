# Resultados de las pruebas backend

**Fecha:** 02/10/2026. **Backend:** submódulo en `44a2339` (último commit de `main`, sin commits nuevos). **Comando:** `go test -v -count=1 ./...`. 3 corridas completas con el mismo resultado (~2 min) y sin contenedores residuales.

| Métrica | Valor |
|---|---:|
| Casos | 55 |
| Aprueban | **43** |
| Fallan | **11** |
| Omitidos | 1 (DELETE de BAC-21B: el endpoint lo crea BAC-24B) |

## Por caso

| Caso | Tarea | Resultado |
|---|---|---|
| BAC-01 a BAC-12, BAC-09, BAC-05/06/06B, prefijo `/api/admin/users`, LOGIN-03 | regresión | ✅ 12/12 |
| LOGIN-04 circuito desde Go y LOGIN-04 front ↔ back (10/10 pasos) | regresión | ✅ 2/2 |
| INF-05 CORS, FIX-31, BAC-17A, INF-06A, INF-08B, BAC-16B, BAC-17B | regresión | ✅ 7/7 |
| **BAC-18B** índice parcial y particiones; purga horaria | BAC-18B | ❌ 2 |
| BAC-19, BAC-20, BAC-16, BAC-21, BAC-18, FIX-17, BAC-13, BAC-15 | regresión | ✅ 8/8 |
| BAC-07, FIX-16/BAC-08, BAC-14, SEC-04 (2) | regresión | ✅ 5/5 |
| **FIX-37** nivel de acceso en `GET /account/profile` | FIX-37 | ❌ |
| BAC-21B energía con FULL_ACCESS + `tareas_asincronas` | BAC-21B | ✅ |
| **FIX-39** campos de `GET /instances`, `shutdown`/`reboot` y auditoría con upid | FIX-39 | ❌ 3 |
| BAC-21B DELETE solo ADMIN | BAC-21B | ⏭️ omitido |
| BAC-21C (2) | regresión | ✅ 2/2 |
| **FIX-40** `504 PROXMOX_TIMEOUT` y `502 PROXMOX_UNAVAILABLE` (nuevo) | FIX-40 | ❌ (el 504 llega con `PROXMOX_UNAVAILABLE`; el 502 ya cumple) |
| INF-07B `/api/events` a través del borde | INF-07B | ✅ (`TASK_FINISHED` en ~1 s) |
| **BAC-29** contrato de la etapa (nuevo) | BAC-29 | ❌ (no existe `docs/contrato-etapa1.md`) |
| **BAC-22** `GET /api/node/status` (nuevo) | BAC-22 | ❌ (404) |
| **BAC-23A** adaptador de IP (nuevo) | BAC-23A | ❌ (el backend no consulta los endpoints de IP) |
| **BAC-25A** pool de UPID (nuevo) | BAC-25A | ❌ (falta `UPID_WORKERS`; las 20 tareas sí terminan `COMPLETED`) |
| BAC-17, SEC-01, integración SEC-01/SEC-02, FIX-28, FIX-08 | regresión | ✅ 5/5 |
| FIX-35 simulador (2) | regresión | ✅ 2/2 |

## Fallos y causa

| Caso | Causa |
|---|---|
| BAC-18B (2) | Sin índice parcial, sin particiones trimestrales ni índice `(fecha_hora, accion, resultado)`, y sin purga horaria |
| FIX-37 | `UsuarioDetalleDTO` solo tiene `instanciasPermitidas` |
| FIX-39 (3) | Faltan los campos nuevos en `GET /instances`; `shutdown` y `reboot` no existen (404); las acciones de energía no se auditan |
| FIX-40 | `mapearErrorProxmox` responde `504 PROXMOX_UNAVAILABLE` ante un timeout |
| BAC-29 | No existe `backend/docs/contrato-etapa1.md` |
| BAC-22 | `GET /api/node/status` responde 404 |
| BAC-23A | Ningún archivo del backend consulta `agent/network-get-interfaces` ni `lxc/{vmid}/interfaces`, y no hay pruebas unitarias de IP |
| BAC-25A | El backend no lee `UPID_WORKERS` (hoy, una goroutine por tarea) |

## Cambios en la suite (02/10/2026)

- **Nuevo `etapa1_acceptance_test.go`** (Ola 1 de la Etapa 1). Reemplaza a `etapa1_core_acceptance_test.go`, que buscaba texto en archivos fijos: un `"lxc"` en `client.go`, la palabra `"workers"`, una expresión regular sobre `instance_handler.go`. Además pedía `instancesSummary` en `BAC-22`, cuando eso es de `BAC-22B`.
  - **BAC-29:** `docs/contrato-etapa1.md` con los nombres de su entregable, el `motivo` de D2, D1, Swagger con `/node/status`, y los códigos nuevos en `docs/estandar_http.md` (inventario de `FIX-08`).
  - **BAC-22**, por comportamiento:
    - sin estado previo y con Proxmox caído, `502` o `504`;
    - `401` sin token y `200` con un OPERATOR real (D1);
    - valores normalizados contra el stub (CPU 25 % y 8 núcleos, RAM 8 de 16 GB, disco 25 de 100 GB, uptime);
    - una clave en Redis con TTL de 10 s o menos;
    - con la caché vigente y Proxmox caído, `200` con `stale: false` en menos de 50 ms;
    - vencido el TTL, `200` con `stale: true`.
  - **BAC-23A** (función interna, sin endpoint todavía): el código consulta los dos endpoints de IP, hay pruebas unitarias, y `go test` de esos paquetes pasa.
  - **BAC-25A:**
    - el backend lee `UPID_WORKERS`;
    - las pruebas unitarias del seguimiento pasan;
    - de 20 órdenes seguidas, las 20 quedan `COMPLETED` en menos de 3 s.
- **FIX-40 (en `puente_etapa1…`), por comportamiento.** Nuevos servicios `proxmox-lento` (acepta la conexión y no responde) y `backend-lento`: el backend tiene que responder `504 PROXMOX_TIMEOUT`. La contraprueba usa el stub caído y espera `502 PROXMOX_UNAVAILABLE`.
- **Stub de Proxmox:**
  - `GET /nodes/{node}/status` con valores conocidos;
  - **un UPID distinto por tarea** (`$request_id`): antes era fijo por instancia, y desde el segundo `start` sobre la misma instancia el backend no podía registrar la tarea (`upid` único en `tareas_asincronas`);
  - una caída simulada con el archivo `/tmp/proxmox-caido` (responde `503`).
- **BAC-21B (auditoría):** busca el `upid` que devolvió la respuesta, en lugar de un UPID fijo.

### Hallazgo de seguridad en la suite (corregido)

Al detener el contenedor del stub, el nombre `proxmox` deja de resolver dentro de Docker. El backend probaba entonces con el dominio de búsqueda del host (`tail6bb3f3.ts.net`, Tailscale), y **`proxmox.tail6bb3f3.ts.net` es el Proxmox real (`100.81.49.19`)**. En una corrida intermedia, el backend de pruebas se conectó a ese servidor con el token falso del stub: **todas las respuestas fueron `401` y no se ejecutó ninguna acción**.

Se corrigió de dos formas:
- la suite ya no detiene el stub: simula la caída con un `503`;
- los tres backends de prueba tienen `dns_search: invalid`, así que un nombre que no resuelve dentro de Docker ya no puede salir por Tailscale.

## Lo que la suite no cubre

- BAC-18B: el tiempo de consulta (< 20 ms).
- BAC-23A: el timeout de 2 s por consulta de IP y la concurrencia acotada. Se verifican con el simulador cuando `BAC-23B` exponga la IP en `GET /api/instances`.
- BAC-25A: la cota de N consultas en paralelo no se observa desde la API. Queda a cargo de las pruebas unitarias del backend, que la suite ejecuta.
- FIX-40: el `502` con el token rechazado. Necesitaría otro backend con un token distinto; lo cubren las pruebas unitarias de `instance_handler_test.go`.
- BAC-22: no hubo contraprueba contra una implementación, porque el endpoint no existe. Los pasos que siguen al 404 se verán cuando se implemente.
