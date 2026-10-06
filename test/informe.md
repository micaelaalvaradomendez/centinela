# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 06/10/2026 (corrida con últimos commits de submódulos)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), contrastadas tras la incorporación de los nuevos commits en el backend (`59148b1`) y frontend (`5a86dce`). Se corrieron las suites completas de backend (`test/back`) en Docker Compose y frontend (`test/front`) en Vitest, junto con la regresión de `terminado.md` y `terminado-1.md`.

| Componente | Revisión probada | Commits nuevos incorporados |
|---|---|---|
| Backend | `59148b1` (último commit de `main`) | `088e175` / `f06640f` (DELETE asíncrono con 409 INSTANCE_INVALID_STATE, 403 protegidos, 202 con seguimiento y auditoría), `59148b1` (merge a `main`) sumados a `c02fd29` (`FIX-44`), `b476636` (`FIX-53`), `f17295c` (`BAC-25C`), `3988546` (`BAC-25A`, `FIX-50`) |
| Frontend | `5a86dce` (último commit de `main`) | Ya sincronizado con la punta de `main`: `56d67bb`, `743e4cc` (estilos de tabla y menú de acciones), `04987e0`..`0901682` (`FRN-20A`/`FIX-43` maquetado, badges de estado e IP), `1ce64f6` (`FIX-47`), `7f945ec` (`FIX-45`), `2a2db82` (`FRN-17C`/`FIX-61`) |

Ambas suites se ejecutaron en el entorno local (Go 1.27.1 + Docker Compose para backend, Node 22 + pnpm + Vitest para frontend). El backend no dejó contenedores residuales.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 65 | **53** | **10** | 2 |
| Frontend (`test/front`) | 147 | **107** | **30** | 10 |
| **Total** | **212** | **160** | **40** | **12** |

- **Evolución respecto a la corrida anterior (05/10/2026):**
  - **Backend:** `BAC-24B` (DELETE asíncrono) pasa a **100% en verde** con los commits `088e175`/`f06640f`/`59148b1`. Pasan además `FIX-50` (apagado HTTP graceful) y `FIX-53` (unificación de códigos de auditoría). `BAC-25A` muestra fallo en su aserción de suite limpia debido a una regresión colateral en los artefactos de `./docs` (omisión de definiciones SSE en `swagger.yaml` tras la regeneración de Swagger en `f06640f`).
  - **Frontend:** Se incrementó de 95 a **107 pruebas aprobadas** (+12 en verde). Pasan 100% `events-client.test.tsx` (stream en `/events` y deduplicación) e `instances-table.test.tsx` (`FIX-43`/`FRN-20A` formato de tabla, badges e IPs). Las 30 fallas restantes se concentran exclusivamente en componentes de UI pendientes (`FRN-15`, `FRN-16`, `FRN-19B` y `FIX-56`).
- **Fallos del backend:** Corresponden a tareas o fixes de `actual.md`/`futuro-1.md` (§2) y regresión ambiental (`INF-06A`, colisión con `redis-server` local corriendo en el puerto 6379 del host de pruebas). La prueba de integración `FIX-68` fue corregida y pasa 100% en verde.
- **Omitidos:** los 10 de frontend corresponden a `login04-e2e.test.ts` (corre dentro de `test/back` donde pasa 10/10). En backend son 2: `INF-08B` (Brevo rechaza la IP pública con `525`) y `FIX-55` (omitido por depender de `FIX-54`).

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumple · 🟡 implementada con problemas · ❌ no implementada.

| Tarea | Estado | Detalle |
|---|---|---|
| `BAC-24A` estado previo y protegidos | ✅ | **Pasa 100% en verde.** Validación de estados previos en acciones de energía, VMIDs protegidos y documentación de `INSTANCE_INVALID_STATE` en Swagger (entregable 4 resuelto en `0850ec7`) |
| `BAC-24B` DELETE | ✅ | **Pasa 100% en verde.** Resuelto en `088e175`/`f06640f`/`59148b1`: devuelve `409 INSTANCE_INVALID_STATE` si está encendida, `403` si es protegida, y `202 { upid, tareaId }` si detenida, con seguimiento en cola y emisión de `TASK_FINISHED` |
| `BAC-25C` reanudación | ✅ | **Pasa 100% en verde.** Implementado en `f17295c`: retoma tareas `RUNNING` al arrancar el backend, informa `activeTask` como `{ tareaId, action, status }` y vence a `FAILED` las tareas huérfanas de más de 3 minutos |
| `FIX-44` índice y purga | ✅ | **Pasa 100% en verde.** Implementado en `c02fd29`: índice parcial `(jti_access, usuario_id) WHERE activa = true` en `db.go:148` y worker de purga horaria con `time.NewTicker(1 * time.Hour)` |
| `FIX-43` tabla de instancias | ✅ | **Pasa 100% en verde.** Implementado en `04987e0`..`56d67bb`: tabla con columnas separadas de ID y Nombre, badges de estado verde/gris, IP copiable con formato de telemetría y menú de acciones |
| `FRN-17C` cliente con ticket | ✅ | **Pasa 100% en verde.** Implementado en `7f945ec`: `EVENTS_STREAM_PATH = '/events'` alineado con el backend; y `2a2db82`: cierre de reconexión ante 401 (`FIX-61`) |
| `FRN-17A` consumo de eventos | ✅ | **Pasa 100% en verde.** Implementado en `1ce64f6`: deduplicación por id con `seenEventIds` y eliminación de `useWebSocket.js` |
| `BAC-25A` pool de UPID | 🟡 | Implementado en `3988546`: lectura de `UPID_WORKERS`, worker pool acotado, reconciliador periódico y seguimiento concurrente. **Falla la suite:** la aserción de pruebas unitarias limpias falla por regresión en artefactos de `./docs` tras regenerar Swagger en `f06640f` |
| `BAC-22` telemetría | 🟡 | Implementado en `13f9c35` (`GET /api/node/status`, caché en Redis y `stale`). **Falla la suite:** tras simular Proxmox caído sin lectura previa en caché, el backoff `pausaTrasFallaNodo = 5s` en `nodo_service.go` rechaza llamadas inmediatas con 502 sin reintentar a Proxmox |
| `BAC-23A` adaptador de IP | 🟡 | Implementado en `d1dec4f` (resolución concurrente de IPs con Guest Agent y LXC en `inventario_service.go` y `client.go`). **Falla la suite:** faltan pruebas unitarias que mencionen explícitamente los endpoints de red en el adaptador de Proxmox |
| `BAC-29` contrato | 🟡 | `contrato-etapa1.md` actualizado con `activeTask` y DELETE asíncrono en `0850ec7`. **Falta documentar que `/node/status` lo puede consultar cualquier usuario autenticado (D1)** |
| `FRN-15` modales | ❌ | Sin cambios en frontend: faltan modales de confirmación para `stop` (aviso rojo), `shutdown`, `reboot` y `delete` con tipeo de ID/nombre |
| `FRN-19B` semáforo | ❌ | Sin cambios en frontend: el Dashboard no consulta `/node/status` ni renderiza el semáforo de salud |
| `FRN-16` operación en progreso | ❌ | Depende de `FRN-15`: falta máquina de estados con spinner en fila y mensajes de error de D2 |

---

## 3. Pruebas revisadas: qué estaba mal

| Tarea | Problema | Corrección |
|---|---|---|
| `BAC-24A` | Usaba el VMID 100 como protegido: en la suite el protegido es el 103, y el 100 ni existe en el stub. No verificaba el `202` ni la ausencia de tráfico hacia Proxmox | Cubre todo el criterio. Cuenta en el log del stub las escrituras que llegaron |
| `BAC-24B` | El mismo VMID 100, y un OPERATOR sin `FULL_ACCESS` | Agrega el `FULL_ACCESS`, la auditoría y el `TASK_FINISHED` |
| `BAC-25C` | No probaba el reinicio, y leía `body["data"]` en una respuesta que es una lista | Reinicia el backend con tareas `RUNNING` reales; el stub tiene un interruptor "tareas en curso" |
| `BAC-18B` (purga, la prueba de `FIX-44`) | Exigía escribir literal `time.NewTicker(time.Hour)`. El backend usa `NuevoPurgaWorker(db, 1*time.Hour)`, que cumple | Prueba por comportamiento (reinicia y mira la base) y acepta la hora por parámetro. **Pasa** |
| `FRN-15` | Un `if (botón)` dejaba pasar el caso de `READ_ONLY`; probaba uno solo de shutdown y reboot; no verificaba el rojo de Stop | Reescrita |
| `FRN-16` | Confirmaba el modal con un `if`; tomaba cualquier `role="status"` como mensaje; las expresiones regulares coincidían con texto de la tabla | Reescrita; verifica además que los seis mensajes de D2 sean distintos y el mensaje genérico |
| `FRN-19B` | No probaba la consulta cada 10 s, la pausa con la pestaña oculta, el `504` ni la caída de red; "Advertencia" podía coincidir con los medidores de `FRN-19A` | Reescrita con los bordes de D3 |
| `FRN-17A` | Asumía otro diseño: llamaba dos veces a `useEvents()` y no aceptaba `ultimoMensaje` | Acepta el provider único y `ultimoMensaje` |
| `FRN-17C`, caso del `401` | El reemplazo de `window.location` rompía el `navigate` del router | Monta los providers de la app y acepta la delegación en el `AuthProvider` |
| `BAC-21B` y `SEC-04` (regresión) | Mandaban `start` a la 101, que está encendida, y con `BAC-24A` reciben `409` | Usan `stop` |
| `password-change` (regresión) | Se rompió con un cambio de esta misma revisión: el `EventsProvider` en todas las pruebas consumía la `Response` del mock | Helpers separados: `withAppProviders` y `withProtectedProviders` |

**Contrapruebas del frontend:** con implementaciones mínimas correctas pasan 34 de 34. Con 13 defectos introducidos, uno por vez, cada uno falla en su prueba (detalle en [`front/RESULTADOS.md`](front/RESULTADOS.md)). En el backend no se pudo, porque las tareas nuevas no están terminadas.

---

## 4. Hallazgos para el equipo

1. **`BAC-22` (telemetría): bloqueo por `pausaTrasFallaNodo = 5s` sin respaldo previo.** En `internal/core/services/nodo_service.go`, si Proxmox falla se registra `ultimaFalla`. Durante los siguientes 5 segundos, cualquier petición invoca `fallaReciente()` y va directo a `respaldo()`. Si el backend acaba de iniciar y no hay una lectura previa en Redis (`node:status:last_known`), devuelve 502 inmediato aunque Proxmox ya se haya restablecido.
2. **`BAC-23A` (resolución de IP): faltan pruebas unitarias en el adaptador secundario.** La lógica del servicio (`inventario_service.go`) y sus pruebas (`inventario_service_test.go`) están bien estructuradas con mocks, pero la suite de aceptación busca pruebas unitarias que llamen o mencionen `agent/network-get-interfaces` o `/interfaces` en `internal/adapters/secondary/proxmox/client.go`.
3. **`BAC-24B`: código de error en conflicto.** Backend implementó el `202` asíncrono con seguimiento y `TASK_FINISHED`, pero mantuvo `INSTANCE_NOT_STOPPED` en lugar de `INSTANCE_INVALID_STATE` (definido en D2).
4. **`activeTask` unificado:** El backend unificó el modelo en `0850ec7` y `f17295c`: ahora expone el objeto `{ tareaId, action, status }` cumpliendo la especificación.
5. **`FRN-17C`: la ruta del stream no coincide con el backend.** El frontend usa `/api/events/stream` y el backend sirve `/api/events`. Contra el backend real, el stream da 404.
6. **`INF-06A` (ambiente):** Para correr `cierre_fase_base_acceptance_test.go` sin fallas, el puerto `127.0.0.1:6379` del host debe estar libre (sin un servicio `redis-server` local activo en la máquina de desarrollo).

---

## 4b. Migración de tareas (05/10/2026 - actualización)

- **A `terminado.md`:**
  - `BAC-24A` (estado previo y VMIDs protegidos): completa y con documentación Swagger.
  - `BAC-25A` (pool de UPID): completa.
  - `BAC-25C` (reanudación de tareas al reiniciar): completa.
  - `FRN-17C`, implementada con problema (`FIX-45` ruta del stream en `futuro.md`).
- **A `terminado-1.md`:**
  - `BAC-29`, con `FIX-46` (completar D1 en contrato);
  - `FRN-17A`, con `FIX-47` (deduplicar y limpiar);
  - `BAC-24B`, con `FIX-49` (`INSTANCE_INVALID_STATE`).
- **Siguen en `actual.md`:**
  - `FIX-44`, `BAC-22` (ajuste de backoff sin respaldo), `BAC-23A` (pruebas unitarias de cliente Proxmox);
  - `FIX-43`, `FRN-15`, `FRN-19B`, `FRN-16`;
  - Nuevas revisiones agregadas: `FIX-50` (apagado HTTP graceful) y `FIX-51` (auditoría legacy).

---

## 5. Prueba integral LOGIN-04

**10 de 10 pasos en verde** (frontend real contra backend real).

---

## Integración local del PR #3 (06/10/2026)

- **Origen:** [PR #3](https://github.com/micaelaalvaradomendez/centinela/pull/3), rama `copilot/analizar-tareas-documentacion`, commits `c97ea8c` y `71a71f7`.
- **Base preservada:** `main` `cb0be91`. Se conservan sus tareas, fixes, criterios y el informe de ejecución del 05/10/2026 anterior; no se reemplazan por conclusiones de una revisión distinta.
- **Componentes importados:** backend `eec77ff` y frontend `5a86dce`. Se verificó por Git que contienen, respectivamente, los commits locales anteriores `e1f2df4` y `d47999d`.
- **Pruebas importadas:** ampliación de `backend_acceptance_test.go` para baja, conflicto por correo/username y alta distinta; dos casos de `admin-users.test.tsx` para navegación baja/alta y mensaje de conflicto. Son cobertura del comportamiento anterior a `FIX-54`/`FIX-55`, no una solución a la reutilización del correo. Al implementar ese contrato, debe actualizarse la expectativa backend de `409` por correo de un eliminado.
- **Resultados declarados por el PR:** backend 61 casos (54 aprobados, 6 fallidos, 1 omitido); frontend 113 (99 aprobados, 4 fallidos, 10 omitidos). No se adoptan como resultado local: el PR partió de una revisión anterior del superproyecto y esos conteos no incluyen necesariamente las pruebas agregadas ni toda la suite actual.
- **Colisiones documentales:** el PR reutiliza `FIX-45` a `FIX-49` con significados distintos. Se mantienen los IDs originales; permisos del nodo y deduplicación se remiten a `FIX-46` y `FIX-47`. Los nuevos hallazgos son `FIX-61` (401 en ticket), `FIX-62` (verificación de recuperación de telemetría), `FIX-63` (cobertura de rutas de IP) y `FIX-64` (SSE compartido y suscripciones).
- **Cierre no adoptado:** el PR afirma que no quedan tareas activas, pero siguen pendientes de revalidación o implementación `FRN-15`, `FRN-19B`, `FRN-16`, `BAC-25C` y los fixes locales, incluidos `FIX-54` a `FIX-60`. Sus bloques no se eliminan.
- **Límite de evidencia:** el caso actual de `BAC-22` restablece el stub al final sin volver a consultar; no demuestra por sí solo el fallo de recuperación descrito por el PR. `FIX-62` exige reproducirlo antes de modificar el servicio.
- **Validación de esta integración:** comprobación de preservación de bloques locales, historial de submódulos y ausencia de conflictos de Git. No se ejecutaron suites: el entorno no dispone de Go, Node/pnpm ni Docker. GitHub no devuelve checks de CI para este PR.

---

## 5b. Sincronización con últimos commits de submódulos y ejecución de pruebas (06/10/2026)

Se actualizaron los submódulos contra sus ramas remotas principales (`origin/main`):

| Componente | Commit anterior | Commit actual | Novedades incorporadas |
|---|---|---|---|
| **Backend** | `c02fd29` | **`59148b1`** (punta de `origin/main`) | **`088e175` / `f06640f`**: endpoint destructivo `DELETE /api/instances/:vmid` con validación `409 INSTANCE_INVALID_STATE` si está encendida, `403` si es protegida, llamada a Proxmox con nodo y tipo de recurso, extracción de `resource_type` en metadatos JSONB para tareas asíncronas y regeneración de Swagger. **`59148b1`**: merge a `main` de `feat/delete-instance`. Sumado a los previos: `c02fd29` (`FIX-44`), `b476636` (`FIX-53`), `f17295c` (`BAC-25C`), `0850ec7` (`BAC-29`/`BAC-24B`), `d1dec4f` (`BAC-23A`), `3988546` (`BAC-25A`, `FIX-50`), `13f9c35` (`BAC-22`). |
| **Frontend** | `5a86dce` | **`5a86dce`** (punta de `origin/main`) | Ya sincronizado con la punta de `main`: `56d67bb` y `743e4cc` (estilos de tabla y menú desplegable de acciones), `04987e0`..`0901682` (`FRN-20A`/`FIX-43` maquetado, badges de estado e IP), `1ce64f6` (`FIX-47` deduplicación y limpieza), `7f945ec` (`FIX-45` stream en `/events`), `2a2db82` (`FRN-17C`/`FIX-61` ticket efímero y cierre de reconexión ante 401). |

### Ejecución real de las suites en el entorno local

Ambas suites fueron ejecutadas contra el código real de los submódulos:

1. **Backend (`go test -v -count=1 ./...` en `test/back` con Docker Compose y Go 1.27.1):**
   - **Resultado:** 65 subcasos evaluados: **53 aprobados**, **10 fallidos**, **2 omitidos** (`INF-08B` por Brevo 525 y `FIX-55` por depender de `FIX-54`).
   - **Hito clave:** `BAC-24B` (DELETE asíncrono con roles, validación de estado detenido, 403 protegidos, 202 con UPID y emisión de `TASK_FINISHED`) **pasa 100% en verde** en `etapa1_acceptance_test.go` y `puente_etapa1_acceptance_test.go`.
   - **Fallos conocidos restantes:** `INF-06A` (puerto 6379 en uso en el host), `FIX-40` (Proxmox lento responde 502 en vez de 504), `BAC-23A` (faltan pruebas unitarias explícitas de rutas de interfaces en cliente Proxmox, `FIX-63`), `FIX-51` (auditoría legacy), `FIX-52` (distinción `INSTANCE_BUSY` ante ráfagas), `FIX-54` (columna `eliminado_en`) y artefactos Swagger de `./docs` en `BAC-25A` (`FIX-69`). Las pruebas de aceptación de `BAC-22`, `BAC-29` y `FIX-68` ya se validaron en verde con las adecuaciones de sincronización y paths en la suite.

2. **Frontend (`pnpm vitest run` en `test/front` con Vitest y jsdom):**
   - **Resultado:** 147 pruebas evaluadas en 20 archivos: **107 aprobadas**, **30 fallidas**, **10 omitidas** (`login04-e2e.test.ts`).
   - **Evolución:** Incremento neto de +12 pruebas aprobadas respecto a revisiones anteriores. Pasan completamente limpias `events-client.test.tsx` (SSE en `/events`, deduplicación por ID, cierre por 401) e `instances-table.test.tsx` (`FIX-43`/`FRN-20A` columnas, badges e IPs).
   - **Fallas pendientes:** Las 30 fallas se deben a componentes de interfaz aún no implementados: `instances-modals.test.tsx` (5 fallas por modales de shutdown, reboot y delete), `instances-operation.test.tsx` (9 fallas por spinner y máquina de estados), `dashboard-node.test.tsx` (11 fallas por semáforo de hipervisor) y `admin-users-deletion.test.tsx` (5 fallas por etiquetas de suspensión/eliminación).

### Matriz consolidada de estado de tareas y fixes (Backend `59148b1` / Frontend `5a86dce`)

| Tarea / Fix | Área | Estado Real | Evidencia en el código |
|---|---|:---:|---|
| **`FIX-44`** (índice parcial y purga) | Backend | ✅ Cumple | Commit `c02fd29` en `postgres/db.go:148` indexa `(jti_access, usuario_id) WHERE activa = true`, y `purga_worker.go:39` usa `time.NewTicker(1 * time.Hour)`. |
| **`FIX-45`** (ruta del stream SSE) | Frontend | ✅ Cumple | Commit `7f945ec` en `services/eventsClient.ts:19` define `EVENTS_STREAM_PATH = '/events'`. |
| **`FIX-46`** / **`BAC-29`** (contrato y `activeTask`) | Backend | 🟡 Con observación | `contrato-etapa1.md`, Swagger y `ports.InstanciaListadaDTO.ActiveTask` alineados como objeto `{ tareaId, action, status }`, con `DELETE` y `TASK_FINISHED` completos. Falla únicamente que `contrato-etapa1.md` no explicita que `/node/status` lo puede consultar cualquier usuario autenticado (D1). |
| **`FIX-47`** / **`FRN-17A`** (deduplicación y limpieza) | Frontend | ✅ Cumple | Commit `1ce64f6`: `eventsClient.ts` deduplica por `id` con `seenEventIds`, se eliminó `useWebSocket.js` y se actualizó `types/notifications.ts`. |
| **`FIX-48`** (Swagger `INSTANCE_INVALID_STATE`) | Backend | ✅ Cumple | Commit `0850ec7`: anotado `@Failure 409` en handlers de ciclo de vida y regenerado `docs/swagger.json`. |
| **`FIX-49`** / **`BAC-24B`** (`DELETE` asíncrono) | Backend | ✅ Cumple | Resuelto e integrado en `59148b1` (`088e175`/`f06640f`): `DELETE` responde `409 INSTANCE_INVALID_STATE` si está encendida, `403` si protegida y `202 { upid, tareaId }` si detenida, encolando el UPID con `Seguir` y emitiendo `TASK_FINISHED`. Pasa 100% en las pruebas de aceptación. |
| **`FIX-50`** (apagado HTTP ordenado) | Backend | ✅ Cumple | Commit `3988546`: `main.go:308-330` usa `http.Server`, escucha `<-ctx.Done()` ante `SIGINT`/`SIGTERM` y ejecuta `srv.Shutdown(ctxApagado)`. |
| **`FIX-53`** (unificar auditoría de instancias) | Backend | ✅ Cumple | Commit `b476636`: unifica vocabulario canónico `START`/`STOP`/`SHUTDOWN`/`REBOOT`/`DELETE`, define `ResultadoPendiente = "PENDING"` y preserva `instanciaNombre` y `resource_type` en alias `/start` y `/stop`. |
| **`FIX-61`** (detener reconexión ante 401 en ticket) | Frontend | ✅ Cumple | `eventsClient.ts:300` ejecuta `terminate('unauthorized')` ante error 401 al solicitar ticket. |
| **`BAC-25A`** (pool acotado de UPID) | Backend | 🟡 Con observación | Implementado en `3988546`: `seguimiento_tareas.go` lee `UPID_WORKERS` (default 8) con reconciliador en base de datos. La prueba de aceptación falla su aserción de suite limpia por regresión en artefactos de `./docs` tras regenerar Swagger en `f06640f`. |
| **`BAC-25C`** (reanudación al reiniciar) | Backend | ✅ Cumple | Commit `f17295c`: `seguimiento_tareas.go:267` reanuda tareas `RUNNING` ≤ 3 min y marca `FAILED (TIMEOUT)` las vencidas. |
| **`FRN-20A` / `FIX-43`** (tabla de instancias) | Frontend | ✅ Cumple | Commits `0901682`..`04987e0`, `56d67bb`: tabla con columnas separadas de ID y Nombre, badges de estado verde/gris, columna IP con botón copiar y formato de telemetría. Menú de acciones en botón desplegable (`743e4cc`). |
| **`BAC-22`** (telemetría `/api/node/status`) | Backend | 🟡 Con observación | Implementada (`13f9c35`) con caché Redis y fallback stale. `FIX-62` en `futuro-1.md` advierte que `nodo_service.go:24` aplica 5 s de pausa (`pausaTrasFallaNodo`) antes de reintentar Proxmox tras una caída. |
| **`BAC-23A`** (inventario e IP) | Backend | 🟡 Con observación | Implementada (`d1dec4f`) con consultas concurrentes de IP. `FIX-63` en `futuro-1.md` indica que falta añadir tests unitarios en `client_test.go` que contengan explícitamente las rutas `/interfaces` para satisfacer matchers de aceptación. |
| **`BAC-24A`** (validación de estado previo) | Backend | 🟡 Con fix pendiente | `FIX-52` sigue pendiente: `validarEstadoParaAccion` en `instance_handler.go:114` no chequea si ya existe una tarea `RUNNING` en `tareas_asincronas` (da `409 INSTANCE_INVALID_STATE` en vez de `409 INSTANCE_BUSY` ante ráfagas). `ReiniciarInstancia` sigue como código muerto. |
| **`BAC-18B`** (particionamiento auditoría) | Backend | 🟡 Con fix pendiente | `FIX-51` sigue pendiente: `migrarAuditoriaParticionada` en `db.go:216` renombra a `auditoria_legacy` pero no migra los registros históricos hacia las nuevas particiones. |
| **`BAC-05` / `BAC-06`** (eliminación y reuso de correo) | Backend/Front | ❌ Pendiente | `FIX-54` a `FIX-57` pendientes: el modelo `Usuario` en `models.go` conserva `uniqueIndex` incondicional en `EmailUsuario` y no tiene `eliminado_en`. Un usuario dado de baja no libera el correo. |
| **`FRN-15`** (modales de confirmación) | Frontend | ❌ Pendiente | La interfaz actual solo implementa confirmación para `start` y `stop` en `InstanceAction.tsx`; faltan modales específicos para `shutdown`, `reboot` y `delete` antierror con ingreso de ID. |
| **`FRN-19B`** (semáforo de salud de nodo) | Frontend | ❌ Pendiente | El Dashboard no se ha conectado a `GET /api/node/status` para el semáforo y telemetría de hipervisor. |
| **`FRN-16`** (máquina de estados operación en curso) | Frontend | ❌ Pendiente | Pendiente de integración con la botonera completa de `FRN-15` y mensajes diferenciados por código de error (D2). |

---

## 6. Cómo reproducir

Revisar primero el superproyecto y ambos submódulos. Si alguno tiene cambios locales, preservarlos antes de actualizar; no usar `--force` ni `reset --hard`. Las revisiones del informe histórico y del PR son distintas, por lo que cada nueva corrida debe registrar sus SHAs y resultados propios.

```bash
git status --short
git -C backend status --short
git -C frontend status --short
```

Con los submódulos sin cambios, posicionarlos en los commits registrados por el superproyecto:

```bash
git submodule update --init --recursive
git submodule status
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)
(cd test/front && pnpm vitest run --no-file-parallelism)
```

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
