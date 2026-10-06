# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 05/10/2026
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), incluidas las nuevas: `FIX-44`, `FIX-43`, `BAC-24A`, `BAC-24B`, `FRN-15`, `FRN-19B`, `FRN-16` y `BAC-25C`. Se revisó que cada una tenga pruebas que evalúen su criterio de éxito; se corrigieron las que estaban mal y se corrió la regresión de `terminado.md` y `terminado-1.md`.

| Componente | Revisión probada | Commits nuevos |
|---|---|---|
| Backend | `e1f2df4` (último commit de `main`) | Contrato de la Etapa 1 (`5787179` a `8ba0c59`), `FIX-39` (`42e94ca`), `FIX-37` (`4e204f1`), `BAC-24A` (`8591e90`) |
| Frontend | `d47999d` (último commit de `main`) | `FIX-36` (`fc27ba2`), `FRN-19A` (`6572011`), `FRN-17C` (`2a2db82`), `FRN-17A` (`8739022`, `0407fbd`) |

Cada suite se corrió 2 veces con el mismo resultado. El backend no dejó contenedores residuales.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 58 | **49** | **8** | 1 |
| Frontend (`test/front`) | 139 | **95** | **34** | 10 |
| **Total** | **197** | **144** | **42** | **11** |

- **Todas las tareas de `actual.md` tienen pruebas.** Las de las tareas nuevas ya existían, pero 7 tenían errores y se reescribieron (§3).
- **Todos los fallos son del producto.** No queda ningún fallo de las pruebas.
- **Omitidos:** los 10 del frontend son `LOGIN-04`, que corre desde `test/back` y pasa 10/10. El del backend es `INF-08B`: Brevo rechaza la IP pública.

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumple · 🟡 implementada con problemas · ❌ no implementada.

| Tarea | Estado | Detalle |
|---|---|---|
| `FRN-17C` cliente con ticket | 🟡 | Cumplen el ticket con `Bearer`, la reconexión con ticket nuevo, el retroceso exponencial y el `401`. **Pero abre el stream en `/api/events/stream`, y el backend lo sirve en `/api/events`**: contra el backend real no llegaría ningún evento |
| `FIX-44` índice y purga | 🟡 | La purga horaria cumple; falta el índice parcial con `jti_access`. Su entregable 2 venía de una prueba mal hecha (§3) |
| `BAC-29` contrato | 🟡 | `contrato-etapa1.md` existe, pero le faltan `DELETE`, `exitstatus`, `PROXMOX_ERROR`, `INSTANCE_INVALID_STATE` y D1. Además define `activeTask` como el id de la tarea, y la tarea pide `{ tareaId, action, status }` |
| `BAC-22` telemetría | ❌ | `GET /api/node/status` responde 404 |
| `BAC-23A` adaptador de IP | ❌ | — |
| `BAC-25A` pool de UPID | ❌ | — |
| `FRN-17A` consumo de eventos | 🟡 | Una sola conexión con `EventsProvider` y descarte de mensajes inválidos: cumplen. **No deduplica** (el cliente descarta el `id`), y sigue `hooks/useWebSocket.js` |
| `FIX-43` tabla de instancias | ❌ | Sin cambios |
| `BAC-24A` estado previo y protegidos | 🟡 | Todo el comportamiento cumple. Falta documentar `INSTANCE_INVALID_STATE` en Swagger (entregable 4) |
| `BAC-24B` DELETE | 🟡 | Existe: solo `ADMIN`, protegidos con `403`. Pero responde `409 INSTANCE_NOT_STOPPED` en lugar de `INSTANCE_INVALID_STATE` (D2), y `204` en lugar de `202 { upid, tareaId }`: no hay seguimiento ni `TASK_FINISHED` |
| `FRN-15` modales | ❌ | La tabla no tiene botonera ni modales |
| `FRN-19B` semáforo | ❌ | El Dashboard no consulta `/node/status` |
| `FRN-16` operación en progreso | ❌ | Depende de `FRN-15` |
| `BAC-25C` reanudación | 🟡 | `activeTask` ya se informa (parte 2). Al reiniciar no se retoma el sondeo |

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

1. **`FRN-17C`: la ruta del stream no coincide con el backend.** El frontend usa `/api/events/stream` y el backend sirve `/api/events`. Contra el backend real, el stream da 404.
2. **`FIX-44`, entregable 2:** no hace falta. El worker ya corre cada una hora, y el pedido venía de una prueba que exigía una forma de código. Conviene sacarlo de la tarea.
3. **`BAC-24B`:**
   - D2 eligió `INSTANCE_INVALID_STATE` para "no está detenida", y el backend usa `INSTANCE_NOT_STOPPED`;
   - devuelve `204` en lugar de `202 { upid, tareaId }`, así que la orden de borrado no se sigue.
4. **`activeTask`:** el contrato publicado (`BAC-29`) lo define como el id de la tarea; `BAC-25C` y el entregable 2 de `BAC-29` piden `{ tareaId, action, status }`. Hay que decidir uno de los dos (la prueba acepta ambos por ahora).
5. **Dependencia que falta en `BAC-25C`:** su criterio pide actualizar `auditoria` al terminar, pero esa entrada de resolución la define `BAC-27`, que no está en `actual.md`.

---

## 4b. Migración de tareas (05/10/2026)

- **A `terminado.md`**, implementada con problema: `FRN-17C`. Su FIX es `FIX-45` (ruta del stream), en `futuro.md`.
- **A `terminado-1.md`**, implementadas con problema. Sus FIX están en `futuro-1.md`:
  - `BAC-29`, con `FIX-46` (completar el contrato y `activeTask`);
  - `FRN-17A`, con `FIX-47` (deduplicar y limpiar);
  - `BAC-24A`, con `FIX-48` (Swagger);
  - `BAC-24B`, con `FIX-49` (`INSTANCE_INVALID_STATE` y `202` con seguimiento).
- **Siguen en `actual.md`**, sin implementar: `FIX-44`, `BAC-22`, `BAC-23A`, `BAC-25A`, `FIX-43`, `FRN-15`, `FRN-19B`, `FRN-16` y `BAC-25C`.
- **Sin FIX, porque son temas de planificación:**
  - el entregable 2 de `FIX-44` sobra;
  - a `BAC-25C` le falta la dependencia con `BAC-27`.

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

## 5b. Sincronización con últimos commits de submódulos y estado actual (06/10/2026)

Se actualizaron los submódulos contra sus ramas remotas principales (`origin/main`):

| Componente | Commit anterior | Commit actual | Novedades incorporadas |
|---|---|---|---|
| **Backend** | `eec77ff` | **`c02fd29`** (punta de `origin/main`) | **`FIX-44`** (`c02fd29`): índice parcial en `sesiones_activas` sobre `(jti_access, usuario_id) WHERE activa = true` y worker de purga con `time.NewTicker(1 * time.Hour)`. Sumado a los commits previos de `main`: `b476636` (`FIX-53` auditoría unificada), `f17295c` (`BAC-25C` recuperación en arranque), `0850ec7` (`BAC-29`/`FIX-46`/`FIX-48`/`FIX-49`), `d1dec4f` (`BAC-23A` IPs), `3988546` (`BAC-25A` pool de UPID y `FIX-50` shutdown), `13f9c35` (`BAC-22` telemetría). |
| **Frontend** | `5a86dce` | **`5a86dce`** (punta de `origin/main`) | Ya sincronizado con la punta de `main`: `56d67bb` y `743e4cc` (estilos de tabla y menú desplegable de acciones), `04987e0`..`0901682` (`FRN-20A`/`FIX-43` maquetado, badges de estado e IP), `1ce64f6` (`FIX-47` deduplicación y limpieza), `7f945ec` (`FIX-45` stream en `/events`), `2a2db82` (`FRN-17C`/`FIX-61` ticket efímero y cierre de reconexión ante 401). |

### Diagnóstico de ejecución de pruebas en el entorno local

Se intentó correr la suite completa de pruebas de backend y frontend:
1. **Backend (`go test`):** `backend/go.mod` declara `go 1.27.1`. El entorno local cuenta con `go1.26.0`; al intentar compilar, Go toolchain requiere descargar automáticamente Go 1.27.1, lo cual falla debido a que la salida a Internet pasa por un proxy corporativo autenticado (`407 Proxy Authentication Required`).
2. **Frontend (`vitest`):** ni `frontend/centinela` ni `test/front` cuentan con `node_modules` instalados en disco. Intentar instalarlos con `npm`/`pnpm` falla por la misma restricción de proxy de red (`E407`).
3. **Entorno Docker:** el usuario no posee permisos sobre el socket del demonio de Docker (`permission denied while trying to connect to unix:///var/run/docker.sock`).

Debido a estas restricciones del entorno, la auditoría del estado actual de las tareas se realizó mediante **análisis estático y trazabilidad exhaustiva de código, commits y suites unitarias existentes**.

### Matriz consolidada de estado de tareas y fixes (Backend `c02fd29` / Frontend `5a86dce`)

| Tarea / Fix | Área | Estado Real | Evidencia en el código |
|---|---|:---:|---|
| **`FIX-44`** (índice parcial y purga) | Backend | ✅ Cumple | Commit `c02fd29` en `postgres/db.go:148` indexa `(jti_access, usuario_id) WHERE activa = true`, y `purga_worker.go:39` usa `time.NewTicker(1 * time.Hour)`. |
| **`FIX-45`** (ruta del stream SSE) | Frontend | ✅ Cumple | Commit `7f945ec` en `services/eventsClient.ts:19` define `EVENTS_STREAM_PATH = '/events'`. |
| **`FIX-46`** / **`BAC-29`** (contrato y `activeTask`) | Backend | ✅ Cumple | Commit `0850ec7`: `contrato-etapa1.md`, Swagger y `ports.InstanciaListadaDTO.ActiveTask` alineados como objeto `{ tareaId, action, status }`, con `DELETE` y `TASK_FINISHED` completos. |
| **`FIX-47`** / **`FRN-17A`** (deduplicación y limpieza) | Frontend | ✅ Cumple | Commit `1ce64f6`: `eventsClient.ts` deduplica por `id` con `seenEventIds`, se eliminó `useWebSocket.js` y se actualizó `types/notifications.ts`. |
| **`FIX-48`** (Swagger `INSTANCE_INVALID_STATE`) | Backend | ✅ Cumple | Commit `0850ec7`: anotado `@Failure 409` en handlers de ciclo de vida y regenerado `docs/swagger.json`. |
| **`FIX-49`** / **`BAC-24B`** (`DELETE` asíncrono) | Backend | ✅ Cumple | Commit `0850ec7`: `DELETE` responde `409 INSTANCE_INVALID_STATE` si está encendida, y `202 { upid, tareaId }` si está detenida, encolando el UPID con `Seguir` y emitiendo `TASK_FINISHED`. |
| **`FIX-50`** (apagado HTTP ordenado) | Backend | ✅ Cumple | Commit `3988546`: `main.go:308-330` usa `http.Server`, escucha `<-ctx.Done()` ante `SIGINT`/`SIGTERM` y ejecuta `srv.Shutdown(ctxApagado)`. |
| **`FIX-53`** (unificar auditoría de instancias) | Backend | ✅ Cumple | Commit `b476636`: unifica vocabulario canónico `START`/`STOP`/`SHUTDOWN`/`REBOOT`/`DELETE`, define `ResultadoPendiente = "PENDING"` y preserva `instanciaNombre` y `resource_type` en alias `/start` y `/stop`. |
| **`FIX-61`** (detener reconexión ante 401 en ticket) | Frontend | ✅ Cumple | `eventsClient.ts:300` ejecuta `terminate('unauthorized')` ante error 401 al solicitar ticket. |
| **`BAC-25A`** (pool acotado de UPID) | Backend | ✅ Cumple | Commit `3988546`: `seguimiento_tareas.go` lee `UPID_WORKERS` (default 8) con reconciliador en base de datos. |
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
pnpm --dir test/front vitest run --no-file-parallelism
```

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
