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

### `FIX-52` - Distinguir `INSTANCE_BUSY` de `INSTANCE_INVALID_STATE` y limpiar código muerto (`BAC-24A`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro (autor de `8591e90`)
- **Estimación:** 1,5 h
- **Depende de:** `BAC-24A` (en `terminado-1.md`).
- **Problema y evidencia (detectado en la revisión manual de Lucas del 05/10/2026, verificada en su PC):**
  1. **`INSTANCE_BUSY` quedó casi inalcanzable.** `validarEstadoParaAccion` (`instance_handler.go:111`) decide únicamente con el `estado` que devuelve `ObtenerInstancia`, comparado contra `estadosRequeridosPorAccion`. No consulta si ya existe una tarea `RUNNING` para ese `vmid` en `tareas_asincronas`. Un `start` seguido de inmediato por un `stop` sobre la misma instancia da `409 INSTANCE_INVALID_STATE` ("estado incompatible"), cuando la causa real es que la instancia está ocupada con la tarea del `start` anterior: Proxmox todavía no reflejó el cambio de estado. El código para distinguir ambos casos ya existe (`tareas_asincronas`, campo `activeTask` de `GET /api/instances`), pero `validarEstadoParaAccion` no lo usa.
  2. **Código muerto.** `ProxmoxPort.ReiniciarInstancia` (`client.go:291`, `ports/proxmox_port.go:99`) nunca se invoca: `CambiarEstado` reemplazó su uso por `Reboot(node, vmid, vmType)` (`instance_handler.go:485`), que sí resuelve nodo y tipo. `ReiniciarInstancia` quedó huérfano en el puerto y en el cliente.
  3. **Comentario desactualizado.** El comentario de `EliminarInstancia` (`client.go:323`) dice *"si está encendida Proxmox responde 500 con un mensaje de bloqueo y se retorna `ErrInstanciaOcupada`"*, pero el handler (`instance_handler.go`, validación de estado previa al `DELETE`) ya rechaza una instancia no `stopped` con `409` antes de invocar a `proxmox.EliminarInstancia`: esa rama de `ErrInstanciaOcupada` es inalcanzable con el flujo actual y el comentario induce a error a quien lea el cliente.
- **Entregable:**
  1. Antes de aplicar `estadosRequeridosPorAccion`, consultar si el `vmid` tiene una tarea `RUNNING` en `tareas_asincronas`; si la tiene, responder `409 INSTANCE_BUSY` sin tráfico de escritura a Proxmox, sin evaluar la matriz de estados.
  2. Eliminar `ReiniciarInstancia` de `ProxmoxPort` y de `client.go` (o, si se prefiere conservarlo, usarlo de forma consistente en lugar de `Reboot` y actualizar el puerto para que reciba `node`/`vmType`).
  3. Corregir el comentario de `EliminarInstancia` para reflejar que el `409` por instancia encendida lo decide el handler, no Proxmox.
- **Criterio de éxito:**
  - Un `start` inmediatamente seguido de un `stop` (u otra combinación con una tarea `RUNNING` en curso) responde `409 INSTANCE_BUSY`, no `409 INSTANCE_INVALID_STATE`.
  - `start` sobre una instancia realmente `running` (sin tarea en curso) sigue respondiendo `409 INSTANCE_INVALID_STATE`.
  - No quedan referencias a `ReiniciarInstancia` sin uso, y el comentario de `EliminarInstancia` describe el flujo real.
  - Siguen en verde todos los casos de `BAC-24A…` y `BAC-24B…` de `etapa1_acceptance_test.go`.

### `FIX-53` - Unificar códigos de acción y resultado en la auditoría de instancias (`BAC-24A`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro (autor de `8591e90`)
- **Estimación:** 1 h
- **Depende de:** `BAC-24A` (en `terminado-1.md`) y `FIX-39` (patrón de auditoría `PENDING`).
- **Problema y evidencia (detectado en la revisión manual de Lucas del 05/10/2026, verificada en su PC):** una misma acción real termina con distinto `accion` en `auditoria` según la ruta y el desenlace:
  - `IniciarInstancia` (alias `/start`): éxito audita `accion: "START"` (`instance_handler.go:344`); falla audita `accion: ports.AccionIniciarVM` = `"INICIAR_VM"` (línea 338).
  - `DetenerInstancia` (alias `/stop`): éxito `"STOP"` (línea 388); falla `ports.AccionDetenerVM` = `"DETENER_VM"` (línea 382).
  - `CambiarEstado` (`/status/:action`): tanto éxito como falla usan el literal `accionAudit` (`"START"`/`"STOP"`/`"SHUTDOWN"`/`"REBOOT"`, líneas 490 y 496), que nunca coincide con las constantes `ports.AccionIniciarVM`/`AccionDetenerVM`/`AccionReiniciarVM`.
  - `EliminarInstancia` sí es consistente: usa `ports.AccionEliminarVM` tanto en éxito como en falla (líneas 548 y 555).

  Filtrar la auditoría por una sola acción no trae todos los registros de la misma operación real. Además, `"PENDING"` (líneas 344, 388, 496) es un literal suelto: no es uno de los resultados definidos en `ports.ResultadoExito`/`ResultadoFalla` (`audit_port.go:105`), y no hay una constante `ports.ResultadoPendiente` equivalente.

  Por separado, `IniciarInstancia` y `DetenerInstancia` descartan la instancia que ya obtuvieron en `validarEstadoParaAccion` (`if _, ok := h.validarEstadoParaAccion(...)`) y auditan `instanciaNombre: ""` (líneas 338 y 344, 382 y 388) y `"resource_type": "vm_or_lxc"` (un valor fijo que no es ni `"qemu"/"lxc"` como usa `CambiarEstado` -`instancia.Tipo`- ni `"VM"/"LXC"` como usan los eventos `TASK_FINISHED`). `CambiarEstado` sí captura `instancia, ok := h.validarEstadoParaAccion(...)` y audita `instancia.Nombre` e `instancia.Tipo` (líneas 490 y 496).
- **Entregable:**
  1. Definir `ports.ResultadoPendiente = "PENDIENTE"` (o el nombre que el equipo prefiera) en `audit_port.go`, junto a `ResultadoExito`/`ResultadoFalla`, y usarlo en los cuatro lugares que hoy escriben el literal `"PENDING"`.
  2. Unificar el `accion` auditado para cada operación real, usando siempre las constantes de `ports` (`AccionIniciarVM`, `AccionDetenerVM`, `AccionReiniciarVM`) tanto para el éxito/`PENDIENTE` como para la falla, en `IniciarInstancia`, `DetenerInstancia` y `CambiarEstado`. Si no existe una constante para `shutdown` (apagado ordenado), agregarla (por ejemplo `AccionApagarVM = "APAGAR_VM"`) en vez de usar el literal `"SHUTDOWN"`.
  3. En `IniciarInstancia` y `DetenerInstancia`, dejar de descartar la instancia de `validarEstadoParaAccion` y auditar `instancia.Nombre` e `instancia.Tipo` como ya hace `CambiarEstado`, en lugar de `""` y `"vm_or_lxc"`.
- **Criterio de éxito:**
  - El mismo `vmid` y la misma acción real auditan siempre el mismo código de `accion`, sin importar si se invocó por `/start`/`/stop` o por `/status/:action`, y sin importar el desenlace (falla o `PENDIENTE`).
  - `instanciaNombre` y `resource_type` en los detalles de auditoría son correctos también al usar los alias `/start` y `/stop`.
  - Una consulta de auditoría filtrando por una sola acción (por ejemplo `AccionIniciarVM`) trae tanto los intentos fallidos como las órdenes despachadas con éxito.
  - Siguen en verde todos los casos de `BAC-24A…` y `BAC-24B…` de `etapa1_acceptance_test.go`.

