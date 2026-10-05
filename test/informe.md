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

## 6. Cómo reproducir

```bash
for s in backend frontend; do
  git -C $s fetch origin --prune && git -C $s checkout -B main origin/main --force && git -C $s reset --hard origin/main
done
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)
pnpm --dir test/front vitest run --no-file-parallelism
```

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
