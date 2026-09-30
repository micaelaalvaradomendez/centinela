# Etapa 1: tareas terminadas

Tareas de la Etapa 1 (ver [`etapa1.md`](etapa1.md)) cuyo criterio de éxito ya se verificó. Las correcciones pendientes están en [`futuro-1.md`](futuro-1.md).

---

# Verificación del 30/09/2026

Revisiones probadas: backend `43a0b06` (simulador `cmd/proxmox-simulador`) y el Proxmox real `100.81.49.19`, sobre el que solo se hicieron lecturas.

#### `FIX-32` - Fidelidad del simulador de Proxmox con la API real (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 30/09/2026, backend `43a0b06`), commits `c997398`, `2c56bab` y `43a0b06`.** Se comparó el simulador con el Proxmox real `100.81.49.19`, con las mismas credenciales y solo lecturas sobre el real:
> - `unprivileged` y los numéricos de `config` ahora son números.
> - `status/current` incluye `ha: {managed: 0}`.
> - `/cluster/nextid` devuelve el próximo VMID libre.
> - `/nodes/{node}/tasks` devuelve `{total, data[]}`, con los mismos campos que el real (`upid`, `type`, `id`, `user`, `tokenid`, `status`, `starttime`, `endtime`, `node`, `pid`, `pstart`). Igual que el real, lista solo las tareas terminadas.
> - El 401 responde `Authentication failed!` sin cuerpo.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Depende de:** ninguna.
- **Problema y evidencia:** comparando endpoint por endpoint, el simulador responde distinto al Proxmox real 9.2.2 en estos puntos:
  1. **Tipo de dato:** en `GET /nodes/{node}/lxc/{vmid}/config`, `unprivileged` llega como número en el real (`1`) y como texto en el simulador (`"1"`). Un código que decodifique a entero funciona con uno y falla con el otro. Conviene revisar también los demás campos numéricos de `config`, como `cores`, `memory` y `swap`.
  2. **Campo faltante:** `GET /nodes/{node}/lxc/{vmid}/status/current` devuelve `"ha": {"managed": 0}` en el real, y el simulador no incluye `ha`.
  3. **Endpoints que el real tiene y el simulador responde con `501`:**
     - `GET /cluster/nextid` (real: `{"data":"106"}`): lo usa el alta de instancias (`RF-07`, `api-proxmox.md` §2).
     - `GET /nodes/{node}/tasks` (real: lista con `total` y `data[]` de `upid`, `status`, `starttime`, `endtime`, `user`): sirve para resincronizar tareas (`BAC-25C`, `RNF-04`).
  4. **Respuesta 401:** el real responde `401 Authentication failed!` **sin cuerpo**; el simulador devuelve un JSON `{"data":null,"message":"invalid token value!"}`. El backend no se ve afectado porque decide por el código HTTP, pero conviene imitar el real.
- **Entregable:**
  1. Devolver `unprivileged` y los demás numéricos de `config` como números.
  2. Agregar `ha: {managed: 0}` a `status/current` (con `managed: 1` y `state` en las instancias con HA, como la 100).
  3. Implementar `GET /cluster/nextid`, que devuelva el menor VMID libre desde 100 como string, y `GET /nodes/{node}/tasks` con las tareas creadas en la sesión del simulador.
  4. Responder el 401 sin cuerpo, con el status text `Authentication failed!`.
  5. Agregar a `simulador_test.go` casos para cada punto.
- **Criterio de éxito:** la comparación contra el Proxmox real no muestra diferencias de estructura ni de tipos en los endpoints de solo lectura.

#### `FIX-33` - Conflicto de bloqueo de Proxmox devuelto como `502 PROXMOX_UNAVAILABLE` (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 30/09/2026, backend `43a0b06`).** Con el backend apuntando al simulador, un `start` seguido de inmediato por un `stop` responde **`409 INSTANCE_BUSY`** ("La instancia se encuentra ejecutando otra tarea. Aguarde a que finalice."), y `502` queda para cuando Proxmox no responde. Además, `start` ahora devuelve `tareaId` junto al `upid`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1 h
- **Depende de:** `BAC-21B` / `BAC-24A` (rutas de energía).
- **Problema y evidencia:** con el backend apuntando al simulador, `POST /api/instances/110/start` seguido de inmediato por `POST /api/instances/110/stop` hace que Proxmox rechace la segunda orden con `500 can't lock file '/var/lock/qemu-server/lock-110.conf' - got timeout`, que es el comportamiento real documentado en las capturas del equipo. El backend la traduce a **`502 PROXMOX_UNAVAILABLE`, "Error al consultar la infraestructura subyacente"**. Pero Proxmox está disponible: la instancia está ocupada con otra tarea. El frontend no puede distinguir "servidor caído" de "esperá a que termine la operación en curso", y con ese mensaje FRN-16 y FRN-17B mostrarían un error equivocado.
- **Entregable:**
  1. En `proxmox/client.go`, reconocer las respuestas `500` cuyo mensaje contiene `can't lock file` o `is locked` y devolver un error de dominio, por ejemplo `ErrInstanciaOcupada`.
  2. En `instance_handler.go`, mapearlo a **`409 INSTANCE_BUSY`** con un mensaje claro.
  3. Agregar el código al inventario de `errorCode` (`FIX-08`) y a Swagger.
- **Criterio de éxito:** una segunda acción sobre una instancia con una tarea en curso responde `409 INSTANCE_BUSY`, y `502` queda reservado para cuando Proxmox no responde.

#### `BAC-28` - Completar el simulador de Proxmox para la Etapa 1

> [!WARNING]
> **Estado: Implementado con problema (verificado el 30/09/2026, backend `43a0b06`).** Se verificó el simulador en local:
> - **Funciona:**
>   - `DELETE /nodes/{node}/{tipo}/{vmid}` devuelve un UPID `qmdestroy`/`vzdestroy` y saca la instancia del inventario al terminar la tarea.
>   - Si la instancia está encendida, responde `500 CT … is running - destroy failed`.
>   - `agent/network-get-interfaces` responde `QEMU guest agent is not running` o `VM … is not running`.
>   - `lxc/{vmid}/interfaces` devuelve las interfaces con el formato de Proxmox.
> - **Diferencias con el Proxmox real** (se comparó `GET /nodes/proxmox/lxc/101/interfaces` y la 104 apagada):
>   - `ip-addresses[].prefix` es **string** en el real (`"8"`, `"24"`) y **número** en el simulador.
>   - Un contenedor apagado responde **`200 {"data":null}`** en el real, y el simulador responde `500 CT … not running`.
>
> La corrección es `FIX-35`, en [`futuro-1.md`](futuro-1.md). La IP de las VMs por guest agent no se pudo comparar, porque el servidor real no tiene VMs.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** ninguna (el simulador ya existe desde `de407a5`).
- **Problema y evidencia:** el simulador imita inventario, estado del nodo, acciones con UPID, estado de tareas, métricas, snapshots, creación y configuración. Le faltan dos endpoints que necesita la Etapa 1, así que sin ellos esas tareas no se pueden desarrollar en local:
  1. `DELETE /nodes/{node}/{tipo}/{vmid}` (lo usa `BAC-24B`). Hoy responde `501`.
  2. La IP de cada instancia: `GET /nodes/{node}/qemu/{vmid}/agent/network-get-interfaces` para VMs y `GET /nodes/{node}/lxc/{vmid}/interfaces` para contenedores (lo usa `BAC-23A`).
- **Entregable:**
  1. `DELETE` que devuelva un UPID de tipo `qmdestroy`/`vzdestroy`, saque la instancia del inventario al terminar la tarea y responda error si está encendida, igual que Proxmox.
  2. Los dos endpoints de IP, con el formato real de Proxmox (capturarlo del servidor real con el token de solo lectura) y el caso "guest agent no está corriendo" para VMs sin agente.
  3. Pruebas en `cmd/proxmox-simulador/simulador_test.go`.
- **Criterio de éxito:** con el simulador, `BAC-23A` obtiene IPs y `BAC-24B` elimina una instancia detenida, sin tocar el servidor.
