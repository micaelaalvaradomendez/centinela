# Etapa 1: correcciones pendientes (FIX)

Correcciones de tareas de la Etapa 1 que ya están implementadas pero no cumplen del todo su criterio de éxito. Las tareas originales están en [`terminado-1.md`](terminado-1.md).

---

## 🛠️ Fixes detectados en la verificación del 03/10/2026

---

## 🛠️ Fixes detectados en la verificación del 05/10/2026

### `FIX-46` - Completar el contrato de la Etapa 1 y alinear `activeTask` (`BAC-29`) (Backend)

- **Área:** Backend
- **Asignado:** Nico (autor de `5787179`, `21b1332` y `4ef36e3`)
- **Estimación:** 1 h
- **Depende de:** `BAC-29` (en `terminado-1.md`).
- **Problema y evidencia (verificado el 05/10/2026, backend `e1f2df4`):** `docs/contrato-etapa1.md` existe, pero le faltan partes del entregable de `BAC-29`. Prueba que falla: `test/back/etapa1_acceptance_test.go`, caso *"BAC-29 el contrato de la Etapa 1 está publicado…"*.
  1. **Entregable 3:** no documenta `DELETE /api/instances/:vmid` ni su `202 { upid, tareaId }`.
  2. **Entregable 4:** los `detalles` de `TASK_FINISHED` no incluyen `exitstatus` ni `motivo` (`PROXMOX_ERROR` | `TIMEOUT`, D2).
  3. **Entregable 5:** no dice que `GET /api/node/status` lo puede consultar cualquier usuario autenticado (D1).
  4. **Entregable 6:** falta el código `INSTANCE_INVALID_STATE`, que ya usa `BAC-24A`.
  5. **`activeTask`:** el entregable 2 pide `{ tareaId, action, status } | null`. El contrato y el backend (`ports.InstanciaListadaDTO.ActiveTask`, *"tareaId en curso o null"*) lo definen como el id de la tarea. `FRN-16B` necesita la acción para mostrar qué operación sigue en curso después de `F5`.
- **Entregable:**
  1. Agregar al contrato los puntos 1 a 4.
  2. Llevar `activeTask` a `{ tareaId, action, status } | null` en el contrato y en `GET /api/instances`. Si el equipo decide quedarse con el id, hay que registrar la decisión en `etapa1.md` y en `FRN-16B`.
  3. Actualizar la versión `.docx` del contrato.
- **Criterio de éxito:**
  - Pasa el caso `BAC-29…` de `etapa1_acceptance_test.go`.
  - `GET /api/instances` informa `activeTask` con la forma documentada.

### `FIX-47` - Deduplicar eventos por `id` y limpiar restos del canal anterior (`FRN-17A`) (Frontend)

- **Área:** Frontend
- **Asignado:** Nico (autor de `8739022`)
- **Estimación:** 1 h
- **Depende de:** `FRN-17A` (en `terminado-1.md`).
- **Problema y evidencia (verificado el 05/10/2026, frontend `d47999d`):**
  1. **No deduplica:** `parseCentinelaEventsMessage` (`services/eventsClient.ts`) arma el mensaje con `tipo`, `severidad`, `recursoId` y `detalles` y **descarta el `id`**, y el cliente entrega cada mensaje en `ultimoMensaje` sin revisar si ya llegó. Un evento repetido, por ejemplo después de una reconexión, se procesa dos veces. Prueba que falla: `events-client.test.tsx`, caso *"un evento repetido (mismo id) se procesa una sola vez"*.
  2. **Restos del canal anterior:** sigue existiendo `hooks/useWebSocket.js`, vacío, y `types/notifications.ts:4` sigue diciendo *"El servidor WebSocket todavía no existe"*. Prueba que falla: *"se eliminó el hook vacío useWebSocket.js y notifications.ts ya no dice que el servidor no existe"*.
- **Entregable:**
  1. Conservar el `id` del `RealtimeEvent` en el mensaje y descartar los que ya se recibieron (por ejemplo, con un conjunto de ids vistos que se limpie al detener el cliente).
  2. Eliminar `hooks/useWebSocket.js` y actualizar el comentario de `types/notifications.ts`.
- **Criterio de éxito:**
  - Pasan los 4 casos `FRN-17A…` de `events-client.test.tsx`.
  - Siguen en verde los de `FRN-17C`.

### `FIX-48` - Documentar `INSTANCE_INVALID_STATE` en Swagger (`BAC-24A`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas (autor de `8591e90`)
- **Estimación:** 0,5 h
- **Depende de:** `BAC-24A` (en `terminado-1.md`).
- **Problema y evidencia (verificado el 05/10/2026, backend `e1f2df4`):** el entregable 4 de `BAC-24A` pide documentar los códigos en Swagger y en `FIX-08`. `INSTANCE_INVALID_STATE` está en `docs/estandar_http.md` (inventario de `FIX-08`), pero **no aparece en `docs/swagger.json`**: las respuestas `409` de `/instances/{vmid}/start`, `/stop` y `/status/{action}` no lo nombran. Prueba que falla: `etapa1_acceptance_test.go`, caso *"BAC-24A energía…"*: `INSTANCE_INVALID_STATE no está documentado en docs/estandar_http.md y docs/swagger.json`.
- **Entregable:** agregar la anotación `@Failure 409 … "INSTANCE_INVALID_STATE …"` a los handlers de energía y regenerar `swagger.json`, `swagger.yaml` y `docs.go`.
- **Criterio de éxito:** pasa el caso `BAC-24A…`, y su comportamiento sigue en verde.

### `FIX-49` - `DELETE /api/instances/:vmid`: código de estado (D2) y `202` con seguimiento (`BAC-24B`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra (autora de `b8d2631`)
- **Estimación:** 1,5 h
- **Depende de:** `BAC-24B` (en `terminado-1.md`) y `FIX-39` (auditoría de la orden despachada).
- **Problema y evidencia (verificado el 05/10/2026, backend `e1f2df4`):** prueba que falla: `test/back/etapa1_acceptance_test.go`, caso *"BAC-24B DELETE solo ADMIN…"*.
  1. **Código de error:** con la instancia encendida, `EliminarInstancia` responde `409 INSTANCE_NOT_STOPPED` (`instance_handler.go:542`). D2 y el entregable 2 de `BAC-24B` piden `409 INSTANCE_INVALID_STATE`, el mismo código de `BAC-24A`, para que el frontend (`FRN-16`) muestre un único mensaje para "la acción no corresponde al estado".
  2. **Respuesta y seguimiento:** con la instancia detenida responde `204` sin cuerpo (`instance_handler.go:558`). `proxmox.EliminarInstancia` descarta el UPID de `qmdestroy`/`vzdestroy`, la tarea no se registra en `tareas_asincronas` (`Seguir`) y no se emite `TASK_FINISHED`. El frontend no puede mostrar "operación en progreso" ni enterarse del final.
  3. **Auditoría:** se registra un resultado final (`ResultadoExito` o `ResultadoFalla`) en el momento de la petición, sin `upid` y sin la orden despachada `PENDING` que pide el entregable 3 (como en `FIX-39`).
  4. **Tipo de recurso:** el entregable 4 pide guardarlo antes de borrar. Hay que verificarlo cuando exista el seguimiento.
- **Entregable:**
  1. Responder `409 INSTANCE_INVALID_STATE` cuando la instancia no está `stopped`.
  2. Que `ProxmoxPort.EliminarInstancia` devuelva el UPID; registrar la tarea con `Seguir`, guardando antes el tipo (`VM`/`LXC`), y responder `202 { upid, tareaId }`.
  3. Auditar la orden despachada (`PENDING`) con su `upid`, como las acciones de energía.
  4. Actualizar Swagger y `docs/estandar_http.md`, y quitar `INSTANCE_NOT_STOPPED` si deja de usarse.
- **Criterio de éxito:**
  - Pasa el caso `BAC-24B…`: `409 INSTANCE_INVALID_STATE`, `202` con `upid` y `tareaId`, la orden en `auditoria` con su `upid`, y `TASK_FINISHED` con `recursoTipo` `LXC` para la 102.
  - Siguen en verde el `401`, el `403` al OPERATOR y el `403 INSTANCE_PROTECTED`.

