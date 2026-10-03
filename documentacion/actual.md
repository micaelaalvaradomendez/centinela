### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 03/10/2026** (backend `4e204f1`, frontend `070e96b`; detalle en [test/informe.md](../test/informe.md)).
> - **Pasaron a [`terminado.md`](terminado.md):**
>   - `FIX-29` y `FIX-40`, completas;
>   - `SEC-03`, `FIX-41`, `FIX-42` y `FIX-38`, completas;
>   - `FIX-36`, completa (PR #77, frontend `070e96b`);
>   - `FIX-37` y `FIX-39`, completas (backend `4e204f1`);
>   - `BAC-18B`, implementada con problemas en backend `4e204f1` (particionamiento trimestral OK, purga periódica OK; falta indexar `jti_access`, derivado a `FIX-44` en `futuro.md`).
> - **Pasaron a [`terminado-1.md`](terminado-1.md):**
>   - `INF-07B`, referencia local verificada (pendiente despliegue CT 103 en `INT-03`);
>   - `FRN-20A`, implementada con problemas en frontend `070e96b` (integración a la API y modales implementados; maquetado, IP y badges derivados a `FIX-43` en `futuro-1.md`).
> - **Tareas que permanecen en este archivo (sin implementar en los submódulos):**
>
> | Tarea | Área | Estado | Pruebas |
> |---|---|---|---|
> | `FRN-17C` | Frontend | No implementada: no existe `useEvents` ni pedido de ticket en front | `events-client.test.tsx` (4) |
> | `BAC-29` | Backend | No implementada: falta `contrato-etapa1.md` y Swagger de Etapa 1 | `etapa1_acceptance…` (1) |
> | `BAC-22` | Backend | No implementada: falta `GET /api/node/status` y telemetría en Redis | `etapa1_acceptance…` (1) |
> | `BAC-23A` | Backend | No implementada: falta resolución de IP por guest agent / LXC | `etapa1_acceptance…` (1) |
> | `BAC-25A` | Backend | No implementada: falta pool acotado de workers (`UPID_WORKERS`) | `etapa1_acceptance…` (1) |
> | `FRN-19A` | Frontend | No implementada: medidores de host en `features/dashboard` vacíos | `dashboard-metrics` (3) |
> | `FRN-17A` | Frontend | No implementada: consumo y distribución de eventos en front | `events-client` (4) |

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

---

### `FRN-17C` (`BRG-02-FRN`) - Cliente de eventos con solicitud previa de ticket efímero y reconexión segura

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1,5 h
- **Ventana propuesta:** Junto a `FRN-17A`.
- **Depende de:** `BAC-21C` (`BRG-02-BAC`).
- **Problema y contexto:** El hook `useEvents` del frontend no puede pasar el JWT por header en `EventSource`/WebSocket ni exponer el access token largo en la URL. Debe solicitar primero el ticket efímero al backend.
- **Entregable:**
  1. En `useEvents` (`FRN-17A`), antes de abrir la conexión hacia `/api/events`, invocar `POST /api/events/ticket` con el interceptor autenticado (`Bearer`) y conectar a `/api/events?ticket=<uuid>`.
  2. Ante una desconexión de red, solicitar un nuevo ticket efímero aplicando retroceso exponencial; si `/api/events/ticket` responde `401`, disparar el cierre de sesión local y redirigir a `/login`.
- **Criterio de éxito:** El frontend se conecta a `/api/events` usando tickets de un solo uso sin exponer el JWT en la URL y se reconecta pidiendo un ticket fresco.

---
# ETAPA 1
> Las tareas de la Etapa 1 ya verificadas están en [`terminado-1.md`](terminado-1.md), y sus correcciones en [`futuro-1.md`](futuro-1.md).
---

> Tareas de la Ola 1 de [`etapa1.md`](etapa1.md), en curso. Ninguna está implementada todavía (verificado el 03/10/2026).

#### `BAC-29` - Contrato HTTP y de eventos de la Etapa 1 (nueva)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.0 h
- **Depende de:** ninguna. Se hace primero, para que el frontend no espere a la implementación.
- **Entregable:** `backend/docs/contrato-etapa1.md` y las anotaciones Swagger (sin implementación) de:
  1. `GET /api/node/status`: `{ cpu: { usagePercent, cores }, ram: { usedGb, totalGb, usagePercent }, storage: { usedGb, totalGb, usagePercent }, uptimeSeconds, instancesSummary: { vms: { running, stopped, paused, total }, lxc: { … } }, stale, fetchedAt }`.
  2. `GET /api/instances`: lo de `BAC-14`, más `ip | null`, `cpuUsage | null` (0-100), `ramUsage | null`, `maxRam | null` (bytes), `nivelAcceso` y `activeTask: { tareaId, action, status } | null`.
  3. Las rutas de energía tal como queden con `FIX-39`. Pueden ser `/status/:action` o una por acción, como las actuales `/start` y `/stop`. Para `start`, `shutdown`, `stop` y `reboot`, y para el `DELETE /api/instances/:vmid` de `BAC-24B`: `202 { upid, tareaId }`.
  4. Los `detalles` de `TASK_FINISHED`: `{ tareaId, accion, estado: COMPLETED|FAILED, exitstatus, motivo?: PROXMOX_ERROR|TIMEOUT, error? }`. `motivo` va solo cuando `estado` es `FAILED` (D2).
  5. Quién puede consultar cada endpoint: `GET /api/node/status`, cualquier usuario autenticado (D1).
  6. La tabla de códigos de error de la etapa con su estado HTTP:
     - `INSTANCE_ACCESS_DENIED`, `INSTANCE_PROTECTED`, `INSTANCE_BUSY`, `INSTANCE_INVALID_STATE` (D2) e `INSTANCE_NOT_FOUND`;
     - `PROXMOX_UNAVAILABLE` (`502`), `PROXMOX_TIMEOUT` (`504`, `FIX-40`) e `INVALID_ACTION`.

     Cada código lleva una línea con su significado, para que el frontend arme un mensaje distinto para cada uno (D2).

     Se agregan también al inventario de `FIX-08`.
- **Criterio de éxito:** el frontend puede maquetar `FRN-19B`, `FRN-20A`, `FRN-16` y `FRN-17B` solo con este documento, y Swagger muestra los endpoints con ejemplos. Si una tarea posterior cambia el contrato, actualiza este archivo y avisa al frontend.

#### `BAC-22` - Adaptador de telemetría del nodo con caché en Redis (`RF-02`)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.5 h
- **Depende de:** `INF-06A` y `BAC-17A` (terminadas). Se desarrolla contra el simulador, que ya responde `/nodes/{node}/status`.
- **Entregable:** endpoint `GET /api/node/status` que consume `/nodes/{node}/status` de Proxmox:
  - Agrega `ObtenerEstadoNodo` a `ProxmoxPort` y al cliente.
  - Lo protege con `RequireAuth`: lo puede consultar cualquier usuario autenticado, `ADMIN` u `OPERATOR` (D1).
  - Normaliza el porcentaje y los núcleos de CPU, convierte de bytes a GB la RAM y el almacenamiento, y devuelve el uptime en segundos, con el formato de `BAC-29`.
  - Guarda el resultado en Redis con TTL de 5 a 10 s, y además el último estado conocido sin TTL.
  - Si Proxmox no responde y hay un último estado conocido, responde `200` con `stale: true`. Si no hay ninguno, responde `502 PROXMOX_UNAVAILABLE` o `504 PROXMOX_TIMEOUT` (D2, `FIX-40`).
- **Criterio de éxito:**
  - Responde en menos de 50 ms con la caché vigente.
  - Al vencer el TTL consulta Proxmox, actualiza Redis y responde `200`.
  - Con el simulador detenido, responde el último estado con `stale: true`.
  - Sin token responde `401`; con un token de `OPERATOR` responde `200` (D1).

#### `BAC-23A` - Adaptador y normalización de inventario Proxmox (QEMU / LXC / IP)
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `BAC-14`, `BAC-28` y `FIX-35` (terminadas). Se desarrolla contra el simulador. Toca los mismos archivos que `FIX-39` (`proxmox/client.go`, `InstanciaListadaDTO`): hay que coordinarse con quien la haga.
- **Entregable:** servicio en Go que:
  - unifica VMs y LXC en una estructura común, desde `/cluster/resources` (como hoy) o desde `/nodes/{node}/qemu` + `/nodes/{node}/lxc`;
  - incluye `cpu`, `mem` y `maxmem` crudos para `BAC-22B`;
  - resuelve la IP con `qemu/{vmid}/agent/network-get-interfaces` (`prefix` numérico) y `lxc/{vmid}/interfaces` (`prefix` string; `data: null` si el contenedor está apagado).
  - Para la IP: primera IPv4 que no sea loopback; si no hay, la primera IPv6 global.
- **Criterio de éxito:**
  - Función interna con pruebas unitarias que devuelve la lista consolidada.
  - Las instancias apagadas o sin guest agent devuelven `ip: null` sin error.
  - Cada consulta de IP tiene un timeout propio (≤ 2 s) y se hacen en paralelo con concurrencia acotada.
 
#### `BAC-25A` - Worker pool acotado para el seguimiento de UPID
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** ninguna pendiente. Parte de `services/seguimiento_tareas.go` (`BAC-21C`). El pool no depende de qué acción originó el UPID.
- **Contexto:** hoy cada tarea lanza su propia goroutine, sin límite.
- **Entregable:**
  1. Reemplazar la goroutine por tarea con un pool de N workers (`UPID_WORKERS`, por defecto 8) que lee de un canal con buffer.
  2. `Seguir` sigue siendo no bloqueante: si el canal está lleno, la tarea queda en `tareas_asincronas` como `RUNNING` y no se pierde, porque `BAC-25C` la retoma.
  3. Cerrar ordenadamente con el contexto del proceso.
  4. Pruebas unitarias con un Proxmox falso.
- **Criterio de éxito:**
  - Con 50 tareas simultáneas nunca hay más de N consultas a Proxmox en paralelo.
  - El fin de cada tarea se detecta en menos de 2 s.
  - Las pruebas de `seguimiento_tareas_test.go` siguen en verde.

#### `FRN-19A` (ex `FRN-13A`) - Maquetado y medidores de recursos del Host (CPU / RAM / Almacenamiento)
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** ninguna (datos de prueba). Puede partir de la maqueta de `pages/Dashboard.tsx`.
- **Entregable:** componentes reutilizables en `features/dashboard`, con Tailwind y shadcn:
  - barras o gauges de CPU (% y núcleos);
  - uso de RAM (GB usados sobre el total);
  - uso de almacenamiento (GB o TB usados sobre el total).
- **Criterio de éxito:** los componentes son responsive y renderizan valores de 0 % a 100 %, con cambio de color según la saturación: normal por debajo del 70 % y advertencia desde el 70 % (D3). Tienen pruebas de componente, incluidos los valores de borde 69 % y 70 %.

#### `FRN-17A` - Consumo de eventos en tiempo real y distribución por instancia
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1.5 h
- **Depende de:** `FRN-17C` (fase base: conexión con ticket y reconexión del mismo hook `useEvents`). Conviene hacerlas juntas. En el backend, `BAC-21C` ya está terminada.
- **Entregable:**
  - Sobre la conexión de `FRN-17C`:
    - parsear cada mensaje como `RealtimeEvent` (`types/notifications.ts`) y descartar los que no cumplen el contrato;
    - deduplicar por `id`;
    - exponer una suscripción por tipo y por `recursoId`, por ejemplo con un provider montado una sola vez en el layout protegido.
  - Eliminar el archivo vacío `hooks/useWebSocket.js` y actualizar el comentario desactualizado de `notifications.ts`.
- **Criterio de éxito:**
  - El Dashboard y la tabla reciben los `TASK_FINISHED` del backend en tiempo real con una sola conexión abierta por pestaña.
  - Un evento repetido se procesa una sola vez.
  - Las pruebas de `events-client.test.tsx` pasan.
