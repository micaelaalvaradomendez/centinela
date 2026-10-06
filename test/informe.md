# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 05/10/2026 (segunda corrida)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), contrastadas tras la incorporación de los nuevos commits en el backend (`13f9c35` a `f17295c`). Se corrieron las suites completas de backend (`test/back`) y frontend (`test/front`) y la regresión de `terminado.md` y `terminado-1.md`.

| Componente | Revisión probada | Commits nuevos |
|---|---|---|
| Backend | `f17295c` (último commit de `main`) | `13f9c35` (`BAC-22`), `3988546` (`BAC-25A`), `d1dec4f` (`BAC-23A`), `0850ec7` (`BAC-29`, `BAC-24B`, `FIX-46`, `FIX-48`, `FIX-49`), `f17295c` (`BAC-25C`) |
| Frontend | `d47999d` (último commit de `main`) | Sin nuevos commits respecto a la revisión previa (`d47999d`) |

Cada suite se corrió 2 veces con el mismo resultado. El backend no dejó contenedores residuales.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 59 | **52** | **6** | 1 |
| Frontend (`test/front`) | 139 | **95** | **34** | 10 |
| **Total** | **198** | **147** | **40** | **11** |

- **Evolución respecto a la corrida anterior:** se sumaron **3 tareas en verde** en backend: `BAC-25A` (pool de UPID), `BAC-25C` (reanudación al reiniciar) y `BAC-24A` (documentación Swagger completada).
- **Fallos del backend:** 5 corresponden a tareas de `actual.md` (§2) y 1 es de regresión ambiental (`INF-06A`, colisión con `redis-server` corriendo en el puerto 6379 del host de pruebas).
- **Frontend sin cambios:** al no haber nuevos commits en `main`, mantiene 95 aprobadas, 34 fallas y 10 omitidas.
- **Omitidos:** los 10 del frontend son `LOGIN-04`, que corre desde `test/back` y pasa 10/10. El del backend es `INF-08B`: Brevo rechaza la IP pública con `525`.

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumple · 🟡 implementada con problemas · ❌ no implementada.

| Tarea | Estado | Detalle |
|---|---|---|
| `BAC-24A` estado previo y protegidos | ✅ | **Pasa 100% en verde.** Validación de estados previos en acciones de energía, VMIDs protegidos y documentación de `INSTANCE_INVALID_STATE` en Swagger (entregable 4 resuelto en `0850ec7`) |
| `BAC-25A` pool de UPID | ✅ | **Pasa 100% en verde.** Implementado en `3988546`: lectura de `UPID_WORKERS`, worker pool acotado, reconciliador periódico y seguimiento concurrente de 20 órdenes completadas sin pérdida |
| `BAC-25C` reanudación | ✅ | **Pasa 100% en verde.** Implementado en `f17295c`: retoma tareas `RUNNING` al arrancar el backend, informa `activeTask` como `{ tareaId, action, status }` y vence a `FAILED` las tareas huérfanas de más de 3 minutos |
| `BAC-24B` DELETE | 🟡 | Avanzó en `0850ec7`: ahora devuelve `202 { upid, tareaId }` con seguimiento y `TASK_FINISHED`, auditoría y control de rol/protegidos. **Pero responde `409 INSTANCE_NOT_STOPPED` en lugar de `INSTANCE_INVALID_STATE` (D2)** |
| `BAC-22` telemetría | 🟡 | Implementado en `13f9c35` (`GET /api/node/status`, caché en Redis y `stale`). **Falla la suite:** tras simular Proxmox caído sin lectura previa en caché, el backoff `pausaTrasFallaNodo = 5s` en `nodo_service.go` rechaza llamadas inmediatas con 502 sin reintentar a Proxmox |
| `BAC-23A` adaptador de IP | 🟡 | Implementado en `d1dec4f` (resolución concurrente de IPs con Guest Agent y LXC en `inventario_service.go` y `client.go`). **Falla la suite:** faltan pruebas unitarias que mencionen explícitamente los endpoints de red en el adaptador de Proxmox |
| `BAC-29` contrato | 🟡 | `contrato-etapa1.md` actualizado con `activeTask` y DELETE asíncrono en `0850ec7`. **Falta documentar que `/node/status` lo puede consultar cualquier usuario autenticado (D1)** |
| `FIX-44` índice y purga | 🟡 | La purga horaria cumple; falta el índice parcial con `jti_access` |
| `FRN-17C` cliente con ticket | 🟡 | Abre el stream en `/api/events/stream` en lugar de `/api/events` (servido por backend) |
| `FRN-17A` consumo de eventos | 🟡 | Conexión única y descarte de inválidos cumplen; no deduplica (el cliente descarta el `id`) y persiste `hooks/useWebSocket.js` |
| `FIX-43` tabla de instancias | ❌ | Sin cambios en frontend |
| `FRN-15` modales | ❌ | Sin cambios en frontend (la tabla no tiene botonera ni modales) |
| `FRN-19B` semáforo | ❌ | Sin cambios en frontend (el Dashboard no consulta `/node/status`) |
| `FRN-16` operación en progreso | ❌ | Depende de `FRN-15` |

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

## 6. Cómo reproducir

```bash
for s in backend frontend; do
  git -C $s fetch origin --prune && git -C $s checkout -B main origin/main --force && git -C $s reset --hard origin/main
done
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)
(cd test/front && pnpm vitest run --no-file-parallelism)
```

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
