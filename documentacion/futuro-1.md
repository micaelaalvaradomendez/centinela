# Etapa 1: correcciones pendientes (FIX)

Correcciones de tareas de la Etapa 1 que ya están implementadas pero no cumplen del todo su criterio de éxito. Las tareas originales están en [`terminado-1.md`](terminado-1.md).

---

## 🛠️ Fixes detectados en la verificación del 03/10/2026

---

## 🛠️ Fixes detectados en la verificación del 05/10/2026

--- 

proximas tareas a cargar 

---


#### `BAC-23B` - IP real en `GET /api/instances` y verificación del filtrado RBAC (`RF-03`)
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `BAC-23A` y `FIX-39` (que agrega los campos en `null` y `nivelAcceso`).
- **Contexto:** el endpoint y el filtrado por rol ya existen (`BAC-14`). Esta tarea solo conecta el adaptador nuevo.
- **Entregable:**
  1. `ListarInstancias` usa el adaptador de `BAC-23A` y completa `ip`.
  2. La IP se resuelve **solo para las instancias que el usuario puede ver**, después del filtro, para no consultar el guest agent de máquinas ajenas.
  3. Mantiene el contrato retrocompatible para `FRN-07`/`FIX-14`.
- **Criterio de éxito:**
  - Un OPERATOR solo ve sus instancias, con su IP o `null`.
  - Sin token, `401`; sin permisos, lista vacía.
  - Con 20 instancias en el simulador, el listado responde en menos de 2 s.
  - Las pruebas `BAC-14`, `BAC-21B` y `FIX-39` siguen en verde.

#### `BAC-22B` (`BRG-05-BAC`) - Conteo de instancias por estado en `GET /api/node/status` (`RF-02`) y métricas por instancia (`RF-03`)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-22`, `BAC-23A` y `FIX-39`.
- **Entregable:**
  1. En `GET /api/node/status`, incluir `instancesSummary: { vms: { running, stopped, paused, total }, lxc: { running, stopped, paused, total } }`, cacheado junto con la telemetría.
  2. En `GET /api/instances`, completar `cpuUsage` (0-100), `ramUsage` y `maxRam` (bytes) con los datos de `BAC-23A`.
- **Criterio de éxito:** `GET /api/node/status` devuelve el desglose por estado coincidiendo con el simulador, y `GET /api/instances` devuelve CPU y RAM de cada instancia (`null` si está apagada).


#### `BAC-25B` - Timeout configurable, reintentos y `exitstatus` en el evento
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25A`.
- **Entregable:**
  1. Cortar el seguimiento a los **3 minutos** (D4): `UPID_TIMEOUT` con valor por defecto `3m`, en lugar de los 10 min fijos de `limiteSeguimiento`. Al vencer, la tarea queda `FAILED` con `motivo: TIMEOUT`.
  2. Reintentos con retroceso cuando falla la consulta del estado, en lugar del reintento fijo cada 1 s.
  3. Completar los `detalles` del `TASK_FINISHED` según `BAC-29`:
     - agregar `exitstatus`;
     - agregar `motivo` cuando la tarea falla (D2): `PROXMOX_ERROR` si Proxmox la terminó con un `exitstatus` distinto de `OK`, y `TIMEOUT` si venció el plazo. Cada caso con su propio `mensaje`;
     - agregar los textos de `shutdown`, `reboot` y `delete` en `nombresAccion`.
  4. Publicar el resultado por el bus existente (`eventosService.Publicar`) y exponerlo para `BAC-27` y `BAC-24C`, por ejemplo con un callback o un suscriptor interno.
- **Criterio de éxito:**
  - Sin `UPID_TIMEOUT` configurado, una tarea que no termina se marca `FAILED` con `motivo: TIMEOUT` a los 3 minutos.
  - Todo `TASK_FINISHED` trae `tareaId`, `accion`, `estado` y `exitstatus`. Los fallidos traen además `motivo`, y es distinto en cada uno de los dos casos.
  - `BAC-27` y `BAC-24C` reciben todos los resultados.


#### `FRN-20B` (ex `FRN-14B`) - Filtros reactivos por tipo, estado y buscador dinámico
- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.0 h
- **Depende de:** `FRN-20A`.
- **Entregable:** barra de control sobre la tabla de inventario:
  - búsqueda en tiempo real por ID o nombre;
  - filtro por tipo (Todas, VM, LXC);
  - filtro por estado (Todas, En ejecución, Detenidas);
  - contador de elementos mostrados sobre el total.
- **Criterio de éxito:** el filtrado es instantáneo, en memoria del cliente, sin peticiones al backend.


#### `FRN-16B` (`BRG-04-FRN`) - Resincronización del estado "Operación en progreso" tras recarga (`F5`) o reconexión

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25C`, `FRN-16` y `FRN-17C`.
- **Entregable:**
  1. Al cargar la tabla, poner en `transitioning` (spinner y botones bloqueados) las filas cuyo ítem de `GET /api/instances` trae `activeTask !== null`, guardando su `tareaId`.
  2. Al reconectarse `/api/events` (`FRN-17C`), volver a pedir `GET /api/instances` sin mostrar carga, para tomar los estados finales de las tareas que terminaron durante la desconexión.
- **Criterio de éxito:** si se recarga con `F5` en medio de una acción, la fila sigue con el spinner y los controles bloqueados hasta que llega `TASK_FINISHED` o termina la tarea.

