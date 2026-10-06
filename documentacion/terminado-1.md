# Etapa 1: tareas terminadas

Tareas de la Etapa 1 (ver [`etapa1.md`](etapa1.md)) cuyo criterio de éxito ya se verificó. Las correcciones pendientes están en [`futuro-1.md`](futuro-1.md).

## Integración del PR #3 (06/10/2026)

> [!IMPORTANT]
> Se conservan las verificaciones históricas de este archivo. El PR #3 actualiza backend a `eec77ff` y frontend a `5a86dce`, ambos descendientes de los commits locales anteriores. Los resultados siguientes fueron declarados por el PR, no ejecutados nuevamente en esta integración; no se trasladan ni eliminan los bloques activos hasta revalidarlos.

| Tarea | Avance incorporado / declarado por el PR | Seguimiento conservando IDs locales |
|---|---|---|
| `BAC-25A` | Commit `3988546`: pool acotado `UPID_WORKERS`; el PR declara aceptación en verde | Revalidar el caso existente de `etapa1_acceptance_test.go` |
| `BAC-22` | Commit `13f9c35`: telemetría de nodo con Redis y lectura stale | Hipótesis de recuperación en `FIX-62` de [`futuro-1.md`](futuro-1.md) |
| `BAC-23A` | Commit `d1dec4f`: inventario consolidado y resolución concurrente de IP | Cobertura de rutas en `FIX-63` de [`futuro-1.md`](futuro-1.md) |
| `BAC-29` | Commit `0850ec7`: contrato ampliado, `activeTask`, DELETE y eventos | Mantener `FIX-46` de [`actual.md`](actual.md) hasta revalidar todo su entregable |
| `FRN-20A` / `FIX-43` | El PR declara 5/5 pruebas de tabla en verde | Confirmar con `instances-table.test.tsx` antes de cerrar `FIX-43` |
| `FRN-17A` | Contexto y cliente de eventos incorporados | Deduplicación en `FIX-47`; conexión y suscripciones en `FIX-64` |

También se incorporan commits de recuperación de tareas y auditoría (`f17295c`, `b476636`); su presencia no cierra automáticamente `BAC-25C` ni `FIX-53`. Los diagnósticos anteriores se conservan como evidencia de sus revisiones originales.

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
> La corrección es `FIX-35`, en [`futuro-1.md`](futuro-1.md).
>
> **Actualización (01/10/2026):** `FIX-35` está completo (ver la verificación del 01/10/2026 en este archivo). El simulador ya responde `lxc/{vmid}/interfaces` igual que el Proxmox real. La IP de las VMs por guest agent no se pudo comparar, porque el servidor real no tiene VMs.

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

---

# Verificación del 01/10/2026

#### `FIX-35` - Formato de las IP de contenedores en el simulador de Proxmox (`BAC-28`) (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 01/10/2026, backend `44a2339`, commit "mas arreglos del simulador").**
> - `test/back/simulador_proxmox_acceptance_test.go` compila y levanta el simulador del submódulo y pasa los 2 casos `FIX-35…`:
>   - en `lxc/101/interfaces`, todos los `ip-addresses[].prefix` son **string**;
>   - `lxc/201/interfaces`, con el contenedor apagado, responde **`200 {"data": null}`**, igual que el Proxmox real.
> - Entregable 1: `agent/network-get-interfaces` de qemu se revisó contra la especificación QAPI de QEMU (`GuestIpAddress.prefix` es entero) y se conserva numérico. `red.go` y `docs/simulador-proxmox.md` documentan la diferencia entre los dos endpoints.
> - Entregable 3: `cmd/proxmox-simulador/simulador_test.go` suma `TestRed_InterfacesLXC_PrefixComoString`, `TestRed_InterfacesLXC_ApagadoDevuelveDataNull` y `TestRed_AgenteQemu_PrefixComoNumero`; `go test ./cmd/proxmox-simulador` pasa.
>
> Con esto se cierra la advertencia de `BAC-28`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 0,5 h
- **Depende de:** `BAC-28`.
- **Problema y evidencia:** se comparó `GET /nodes/proxmox/lxc/{vmid}/interfaces` del simulador con el del Proxmox real 9.2.2 (lecturas del 30/09/2026):
  1. **Tipo de dato:** `ip-addresses[].prefix` es **string** en el real (`"prefix":"24"`) y **número** en el simulador (`"prefix":24`). `BAC-23A` va a leer la IP de ahí, y un decodificador escrito contra el simulador falla contra el real, o al revés.
  2. **Contenedor apagado:** el real responde **`200 {"data":null}`** (verificado con la 104, apagada); el simulador responde `500 CT 201 not running`. Con el simulador, `BAC-23A` trataría como error algo que en el real es "sin IP".
- **Entregable:**
  1. Devolver `prefix` como string en `lxc/{vmid}/interfaces`. Revisar también `agent/network-get-interfaces` contra la documentación de Proxmox, porque no se pudo comparar: el servidor no tiene VMs.
  2. Para un LXC apagado, responder `200` con `{"data": null}`.
  3. Agregar ambos casos a `cmd/proxmox-simulador/simulador_test.go`.
- **Criterio de éxito:** la comparación contra el Proxmox real de `lxc/{vmid}/interfaces` no muestra diferencias de tipos, y un contenedor apagado responde igual en los dos.

---

# Verificación del 03/10/2026

#### `INF-07B` (`BRG-03`) - Configuración de Nginx para SSE en el servidor (`RNF-06`)

> [!NOTE]
> **Estado: Referencia local completada (verificado el 03/10/2026, backend `4e204f1`).**
> - La configuración en `docker/nginx-edge.conf` establece `proxy_http_version 1.1`, `proxy_buffering off` y `proxy_read_timeout 1h`.
> - **Pruebas:** en `test/back/puente_etapa1_acceptance_test.go`, el caso *"INF-07B BRG-03 /api/events a través del borde Nginx entrega TASK_FINISHED sin buffer"* pasa en verde (entrega en ~1.0 s).
> - **Pendiente operacional:** aplicar esta configuración en el contenedor Nginx del CT 103 en Proxmox con certificados reales, lo cual se validará en `INT-03`.

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.0 h
- **Depende de:** `FIX-31` y `BAC-21C`.
- **Entregable:**
  1. Llevar al Nginx del CT 103 la configuración de `docker/nginx-edge.conf` con certificados TLS.
  2. Confirmar que `/api/events` no retiene buffer y tolera conexiones abiertas sin tráfico.
- **Criterio de éxito:** un cliente conectado a `/api/events` recibe `TASK_FINISHED` sin retardo y la conexión no se cierra por inactividad.

#### `FRN-20A` (ex `FRN-14A`) - Tabla interactiva de inventario con badges de estado e IP (`RF-03`) (Frontend)

> [!WARNING]
> **Estado: Implementada con problemas (verificado el 03/10/2026, frontend `070e96b`, PR #74, commit `52cd282`).**
> - **Completado:**
>   - Reemplaza la maqueta estática de `Instances.tsx` con integración reactiva a `useInstances.ts` y `instanceService.ts`, consultando dinámicamente `GET /instances` con el token Bearer del usuario autenticado.
>   - Filtra instancias por permisos con `canAccessInstance(instance.id)` y monta `InstanceAction.tsx` con modales de confirmación para `start` y `stop`.
> - **Problemas detectados (se corrigen en `FIX-43` en `futuro-1.md`):**
>   1. La celda de nombre renderiza `{instance.name} ({instance.id})` en un único bloque de texto, lo que rompe las consultas de prueba y la separación semántica entre nombre e identificador numérico.
>   2. La columna de IP renderiza un guión fijo `—` en lugar de mostrar la IP de la instancia o `"No detectada"` cuando viene en `null`, y no cuenta con botón para copiar al portapapeles.
>   3. El estado de la instancia se muestra como texto simple `{instance.status}` sin badge ni clases visuales diferenciadas (verde para Running, gris para Stopped).
>   4. La estructura de la fila no provee la botonera de acciones completa cuando no hay acciones inmediatas o falta la IP.
> - **Pruebas:** `test/front/instances-table.test.tsx` falla 5/5 casos debido a estas discrepancias de maquetado e interfaz.

- **Área:** Frontend
- **Asignada:** Luz (implementación inicial de Cristian en PR #74)
- **Estimación:** 2.5 h
- **Depende de:** `BAC-14`, `BAC-29`, `FIX-39`.
- **Criterio de éxito original:** VMs y LXC en la misma tabla; IP null muestra "No detectada"; tolera campos nuevos en null; badges de estado diferenciados.

#### `INF-07` - Token de Proxmox, conectividad y VMIDs protegidos en el servidor

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.0 h
- **Depende de:** ninguna. Usa variables que el backend ya lee: `PROXMOX_*` (`BAC-14`) y `PROXMOX_PROTECTED_VMIDS` (`FIX-33`/`SEC-04`).
- **No bloquea al desarrollo:** solo hace falta para `INT-03`.
- **Estado verificado (29/09/2026):** el token `centi-api@pve!backend-token` responde 200 en `/version`, `/nodes/proxmox/status`, `/cluster/resources` y `/cluster/nextid`, y tiene `Sys.Audit`, `VM.Audit`, `VM.PowerMgmt`, `VM.Allocate` y `VM.GuestAgent.Audit`. Alcanza para telemetría, inventario, energía, borrado y lectura de IP.
- **Entregable:**
  1. **Rotar el secreto del token**, que hoy está en texto plano en `documentacion/api-proxmox.md`. Documentar en ese archivo el token y sus privilegios, sin el secreto.
  2. Verificar con `curl`, desde el LXC del backend (pruebas y estable), la conectividad HTTPS hacia `/nodes/{node}/status` y `/cluster/resources`.
  3. Cargar en el `.env` del servidor `PROXMOX_URL`, `PROXMOX_NODE`, `PROXMOX_TOKEN_ID`, `PROXMOX_TOKEN_SECRET` (el nuevo) y `PROXMOX_PROTECTED_VMIDS=100,101,102,103,104,105`.
- **Criterio de éxito:**
  - Desde el LXC del backend las dos consultas responden `200 OK` con el token nuevo, y el anterior responde `401`.
  - El backend desplegado lista las instancias reales.
  - Un `stop` sobre un VMID protegido responde `403 INSTANCE_PROTECTED`.

#### `FRN-19A` (ex `FRN-13A`) - Maquetado y medidores de recursos del Host (CPU / RAM / Almacenamiento) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 03/10/2026, frontend `907efe5`, PR #78, commit `6572011`).**
> - Se implementó el componente `ResourceMeter.tsx` en `features/dashboard/components`, soportando barras de progreso y medidores para CPU (% y núcleos), RAM (usada y total) y almacenamiento (formateo dinámico GB/TB).
> - Aplica estilos reactivos acordes a la saturación: visualización estándar por debajo del 70% y estado de advertencia (`warning`) a partir del 70% (D3).
> - Se integró en `Dashboard.tsx` y se incorporó la suite de pruebas unitarias `centinela/src/test/resource-meters.test.tsx` evaluando umbrales de borde al 69% y 70%.
> - **Pruebas:** `test/front/dashboard-metrics.test.tsx` pasa **3/3 en verde ✅**.

- **Área:** Frontend
- **Asignada:** Belinda (PR #78)
- **Estimación:** 2.0 h
- **Depende de:** ninguna.
- **Entregable:** componentes reutilizables en `features/dashboard` con cambio de color según saturación y pruebas unitarias de borde.
- **Criterio de éxito:** medidores responsive de 0% a 100%, advertencia desde el 70% y tests de componente aprobados.

---

# Verificación del 05/10/2026

Revisiones probadas: backend `e1f2df4` y frontend `d47999d`. Evidencia completa en [test/informe.md](../test/informe.md).

#### `BAC-29` - Contrato HTTP y de eventos de la Etapa 1 (nueva)

> [!WARNING]
> **Estado: Implementada con problema (verificado el 05/10/2026, backend `e1f2df4`, commits `5787179`, `21b1332`, `4ef36e3` y `8ba0c59`).** Existe `docs/contrato-etapa1.md`, con su versión `.docx`, y Swagger documenta `/node/status`. `test/back/etapa1_acceptance_test.go`, caso `BAC-29…`, falla porque al contrato le faltan partes del entregable:
> - entregable 3: el `DELETE /api/instances/:vmid` y su `202 { upid, tareaId }`;
> - entregable 4: `exitstatus` y `motivo: PROXMOX_ERROR` en los `detalles` de `TASK_FINISHED`;
> - entregable 5: que `GET /api/node/status` lo puede consultar cualquier usuario autenticado (D1);
> - entregable 6: el código `INSTANCE_INVALID_STATE`.
>
> **Diferencia con la tarea:** el entregable 2 pide `activeTask: { tareaId, action, status } | null`, y el contrato lo documenta como el id de la tarea en curso; el backend lo implementa así (`ports.InstanciaListadaDTO.ActiveTask`).
>
> La corrección es **`FIX-46`**, en [`futuro-1.md`](futuro-1.md).

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

#### `BAC-24A` - Validación de estado previo y VMIDs protegidos en acciones de energía (`RF-04`)

> [!WARNING]
> **Estado: Implementada con problema (verificado el 05/10/2026, backend `e1f2df4`, commit `8591e90`).** `test/back/etapa1_acceptance_test.go`, caso `BAC-24A…`. **Todo el comportamiento cumple**, en la ruta genérica `/status/:action` y en los alias `/start` y `/stop`:
> - `409 INSTANCE_INVALID_STATE` para `start` sobre una instancia encendida, y para `stop`, `shutdown` y `reboot` sobre una detenida;
> - `403 INSTANCE_PROTECTED` para `shutdown`, `reboot` y `stop` sobre un VMID protegido;
> - `400 INVALID_ACTION` para una acción desconocida;
> - `403` para un usuario sin la instancia asignada;
> - `202 { upid, tareaId }` en las órdenes válidas;
> - ninguna escritura en Proxmox al rechazar (verificado en el log del stub).
>
> **Falla el entregable 4:** `INSTANCE_INVALID_STATE` está en `docs/estandar_http.md` (inventario de `FIX-08`), pero no en Swagger (`docs/swagger.json`). Siguen en verde `FIX-16`, `SEC-04`, `FIX-39`, `FIX-40` y `LOGIN-04`.
>
> La corrección es **`FIX-48`**, en [`futuro-1.md`](futuro-1.md).
>
> **Problemas adicionales detectados en la revisión manual de Lucas del 05/10/2026 (no los cubre `FIX-48`):**
> - **La validación de estado previo confunde "estado incompatible" con "instancia ocupada".** `validarEstadoParaAccion` (`instance_handler.go:111`) solo compara el `estado` que informa Proxmox contra `estadosRequeridosPorAccion`, sin mirar si ya hay una tarea `RUNNING` para ese `vmid` en `tareas_asincronas`. Un `start` seguido de inmediato por un `stop` sobre la misma instancia responde `409 INSTANCE_INVALID_STATE` en vez de `409 INSTANCE_BUSY`, porque el estado todavía no reflejó el cambio disparado por el `start`. Corrección: `FIX-52`.
> - **Código muerto y comentario desactualizado.** `ProxmoxPort.ReiniciarInstancia` (`client.go:291`) nunca se invoca: `CambiarEstado` solo usa `Reboot(node, vmid, vmType)` (`instance_handler.go:485`). El comentario de `EliminarInstancia` (`client.go:323`) dice que borrar una instancia encendida devuelve `ErrInstanciaOcupada` de Proxmox, pero el handler ya la rechaza antes con `409` sin llegar a llamarlo: la rama de Proxmox es inalcanzable con el código actual. Corrección: `FIX-52`.
> - **Códigos de auditoría mezclados entre rutas.** Un `start` exitoso vía `/start` audita `accion: "START"` (`instance_handler.go:344`), uno fallido vía la misma ruta audita `ports.AccionIniciarVM` = `"INICIAR_VM"` (línea 338), y el mismo `start` disparado vía `/status/start` audita `accionAudit = "START"` en ambos casos (líneas 490 y 496). Filtrar la auditoría por una sola acción no trae todos los registros de la misma operación real. Además, `"PENDING"` (líneas 344, 388, 496) no es uno de los resultados definidos en `ports.ResultadoExito`/`ResultadoFalla`. Las rutas `/start` y `/stop` tampoco registran `instanciaNombre` (pasan `""` en lugar de `instancia.Nombre`, que si obtienen pero descartan) ni el tipo real de recurso (`"resource_type": "vm_or_lxc"`, un valor fijo que no es ni `"qemu"/"lxc"` ni `"VM"/"LXC"`), a diferencia de `CambiarEstado`, que sí usa `instancia.Nombre` e `instancia.Tipo`. Corrección: `FIX-53`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `FIX-39` (`shutdown` y `reboot`).
- **Entregable:**
  1. Antes de enviar la orden, validar el estado actual de la instancia:
     - `start` solo si está `stopped`;
     - `shutdown`, `stop` y `reboot` solo si está `running`.

     Si no corresponde, responder `409 INSTANCE_INVALID_STATE` (D2) sin enviar tráfico de escritura a Proxmox.
  2. Aplicar `RejectProtectedInstance` a `shutdown` y `reboot`, en todas las rutas que se usen para esas acciones.
  3. Si se usa una ruta genérica como `/status/:action`, una acción desconocida responde `400 INVALID_ACTION`.
  4. Documentar los códigos en Swagger y `FIX-08`.
- **Criterio de éxito:**
  - Toda orden válida devuelve `202 { upid, tareaId }`.
  - `start` sobre una instancia encendida devuelve `409 INSTANCE_INVALID_STATE`.
  - `reboot` sobre un VMID protegido devuelve `403 INSTANCE_PROTECTED`.
  - Un usuario sin la instancia asignada recibe `403` sin tráfico hacia Proxmox.

#### `BAC-24B` - Endpoint de eliminación `DELETE /api/instances/:vmid`

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, backend `59148b1`, commits `088e175` y `f06640f`).** El endpoint destructivo `DELETE /api/instances/:vmid` cumple todos sus entregables:
> - sin token, `401`;
> - un OPERATOR recibe `403`;
> - un VMID protegido recibe `403 INSTANCE_PROTECTED`;
> - con la instancia encendida responde `409 INSTANCE_INVALID_STATE` (D2);
> - con la instancia detenida responde `202 { upid, tareaId }`, encolando el UPID en `Seguir`, despachando la orden auditada con `PENDING` y emitiendo `TASK_FINISHED` con `recursoTipo` correcto.
> Pasa 100% en verde en `etapa1_acceptance_test.go` y `puente_etapa1_acceptance_test.go`. Resuelto por `FIX-49`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `FIX-39`, porque reutiliza su auditoría de la orden despachada y toca los mismos archivos. Se desarrolla contra el simulador, que ya implementa el `DELETE` (`BAC-28`).
- **Contexto:** el endpoint no existe. La regla de `DELETE` estaba en `BAC-21B`, pero quedó fuera de su alcance y de `FIX-39`. La prueba `puente_etapa1…` del `DELETE` se omite hasta que exista.
- **Entregable:**
  1. Agregar `EliminarInstancia` a `ProxmoxPort` y al cliente: `DELETE /nodes/{node}/{qemu|lxc}/{vmid}`.
  2. Montar `DELETE /api/instances/:vmid` con:
     - `RequireRole("ADMIN")`;
     - `RejectProtectedInstance`;
     - la validación de estado: si la instancia no está `stopped`, `409 INSTANCE_INVALID_STATE` (D2) sin enviar la orden a Proxmox.
  3. Registrar el UPID de `qmdestroy`/`vzdestroy` en el seguimiento (`Seguir`) y en la auditoría de la orden despachada, y responder `202 { upid, tareaId }`.
  4. Guardar el tipo de recurso (`VM`/`LXC`) **antes** de borrar: hoy `seguir` lo averigua con `ObtenerInstancia` al terminar, y después del borrado esa consulta falla.
  5. Documentar el endpoint y sus códigos en Swagger y `FIX-08`.
- **Criterio de éxito:**
  - Un `OPERATOR`, aunque tenga `FULL_ACCESS`, recibe `403`.
  - Borrar una instancia encendida responde `409 INSTANCE_INVALID_STATE`.
  - Sobre una detenida responde `202`, el simulador la saca del inventario y llega `TASK_FINISHED` con `recursoTipo` correcto.
  - Sobre un VMID protegido responde `403 INSTANCE_PROTECTED`.

#### `FRN-17A` - Consumo de eventos en tiempo real y distribución por instancia

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, frontend `5a86dce`, commit `1ce64f6` `FIX-47`).** `test/front/events-client.test.tsx`, bloque `FRN-17A`:
> - una sola conexión por pestaña: `EventsProvider` en `ProtectedLayout`, y los consumidores con `useEventsContext()`;
> - descarta los mensajes que no cumplen `RealtimeEvent`;
> - deduplica por `id` en `parseCentinelaEventsMessage` manteniendo `seenEventIds`;
> - se eliminó `hooks/useWebSocket.js` y se actualizó `types/notifications.ts`.
> Pasan 4 de 4 pruebas de `FRN-17A` en `events-client.test.tsx`. Resuelto por `FIX-47`.

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

---

# Verificación del 06/10/2026: tareas movidas desde `actual.md`

Evidencia completa en [test/informe.md](../test/informe.md). Revisiones probadas: backend `59148b1` y frontend `5a86dce`.

#### `BAC-22` - Adaptador de telemetría del nodo con caché en Redis (`RF-02`)

> [!WARNING]
> **Estado: Implementada con problema (verificado el 06/10/2026, backend `59148b1`, commit `13f9c35`).**
> - **Cumple:** implementa `GET /api/node/status` en `nodo_service.go`, normaliza CPU/RAM/disco y uptime, cachea en Redis con TTL y sirve fallback `stale` si Proxmox cae teniendo lectura previa.
> - **Problema:** si el backend arranca y Proxmox falla antes de poblar la caché (`node:status:last_known`), `nodo_service.go:24` aplica `pausaTrasFallaNodo = 5s`. Durante ese intervalo, cualquier consulta devuelve `502 PROXMOX_UNAVAILABLE` sin reintentar a Proxmox, haciendo fallar el caso de recuperación en `etapa1_acceptance_test.go:91`.
>
> La corrección es **`FIX-62`**, en [`futuro-1.md`](futuro-1.md).

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.5 h
- **Depende de:** `INF-06A` y `BAC-17A`.
- **Criterio de éxito original:** Responde en < 50 ms en caché; refresca al vencer TTL; ante caída de Proxmox sirve último estado con `stale: true`; accesible para cualquier usuario autenticado (`ADMIN` u `OPERATOR`).

#### `BAC-23A` - Adaptador y normalización de inventario Proxmox (QEMU / LXC / IP)

> [!WARNING]
> **Estado: Implementada con problema (verificado el 06/10/2026, backend `59148b1`, commit `d1dec4f`).**
> - **Cumple:** unifica inventario de VMs y LXCs, resuelve IPs de forma concurrente consultando QEMU Guest Agent e interfaces de LXC, y devuelve `ip: null` si la máquina está apagada.
> - **Problema:** la suite de aceptación (`etapa1_acceptance_test.go:165`) busca pruebas unitarias que mencionen explícitamente `network-get-interfaces` o `/interfaces` en archivos `_test.go` del cliente/adaptador de Proxmox.
>
> La corrección es **`FIX-63`** (y profundización en `FIX-67`), en [`futuro-1.md`](futuro-1.md).

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `BAC-14`, `BAC-28` y `FIX-35`.
- **Criterio de éxito original:** Devuelve inventario consolidado con IP resuelta concurrentemente y timeout acotado por petición.

#### `BAC-25A` - Worker pool acotado para el seguimiento de UPID

> [!WARNING]
> **Estado: Implementada con problema (verificado el 06/10/2026, backend `59148b1`, commit `3988546`).**
> - **Cumple:** pool acotado con `UPID_WORKERS` (default 8), encolado no bloqueante con reconciliador en base de datos (`seguimiento_tareas.go`), y procesamiento concurrente sin pérdida de tareas.
> - **Problema:** la aserción de suite limpia en `etapa1_acceptance_test.go:178` ejecuta `go test ./docs ...`, la cual falla debido a que la regeneración de Swagger en `f06640f` omitió las definiciones de modelos SSE en `swagger.yaml`.
>
> La corrección es **`FIX-69`** (y verificación estricta de concurrencia en `FIX-67`), en [`futuro-1.md`](futuro-1.md).

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-21C`.
- **Criterio de éxito original:** Con 50 tareas simultáneas no supera $N$ workers en vuelo; detección en menos de 2 s.

#### `BAC-25C` (`BRG-04-BAC`) - Reanudación de UPIDs en curso al arrancar el Backend (`RNF-04`)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, backend `59148b1`, commit `f17295c`).**
> - Reanuda tareas `RUNNING` ≤ 3 min al iniciar el backend y marca `FAILED (TIMEOUT)` las vencidas.
> - Informa `activeTask: { tareaId, action, status }` en `GET /api/instances` cruzando con las tareas en curso.
> - Pasa 100% en verde en `etapa1_acceptance_test.go`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25A` y `FIX-39`.
- **Criterio de éxito:** Tareas en curso se retoman tras reinicio y se informa `activeTask` en inventario.

#### `FIX-43` - Corregir maquetado, visualización de IP, badges de estado y columnas en tabla de instancias (`FRN-20A` / RF-03) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, frontend `5a86dce`, commits `0901682`..`56d67bb`, `743e4cc`).**
> - Muestra columnas separadas para ID y Nombre, permitiendo identificación unívoca.
> - Columna IP con valor real o "No detectada" y botón copiar interactivo.
> - Badges de estado con estilos verde/gris según `running`/`stopped`.
> - Menú desplegable y botonera de acciones operativa.
> - Pasan 5 de 5 pruebas en `test/front/instances-table.test.tsx`.

- **Área:** Frontend
- **Asignada:** Luz / Cristian
- **Estimación:** 1.0 h
- **Depende de:** `FRN-20A`.
- **Criterio de éxito:** Pasan las 5 pruebas de `instances-table.test.tsx`.

#### `FIX-46` - Completar el contrato de la Etapa 1 y alinear `activeTask` (`BAC-29`) (Backend)

> [!WARNING]
> **Estado: Implementada con problema (verificado el 06/10/2026, backend `59148b1`, commit `0850ec7`).**
> - **Cumple:** contrato actualizado con `activeTask` estructurado, DELETE asíncrono y eventos.
> - **Problema:** falta documentar explícitamente en `docs/contrato-etapa1.md` que `/node/status` lo puede consultar cualquier usuario autenticado (D1).
>
> Pendiente D1 registrado en [`futuro-1.md`](futuro-1.md).

- **Área:** Backend
- **Asignado:** Nico
- **Estimación:** 1 h
- **Depende de:** `BAC-29`.
- **Criterio de éxito:** Contrato alineado y paso de `etapa1_acceptance_test.go`.

#### `FIX-47` - Deduplicar eventos por `id` y limpiar restos del canal anterior (`FRN-17A`) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, frontend `5a86dce`, commit `1ce64f6`).**
> - Deduplica eventos por `id` en `eventsClient.ts` con `seenEventIds`.
> - Se eliminó el archivo residual `hooks/useWebSocket.js` y se actualizó `types/notifications.ts`.
> - Pasan los casos de `FRN-17A` en `events-client.test.tsx`.

- **Área:** Frontend
- **Asignado:** Nico
- **Estimación:** 1 h
- **Depende de:** `FRN-17A`.
- **Criterio de éxito:** Eventos repetidos se procesan una sola vez y no quedan restos del canal anterior.

#### `FIX-48` - Documentar `INSTANCE_INVALID_STATE` en Swagger (`BAC-24A`) (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, backend `59148b1`, commit `0850ec7`).**
> - Se documentó la respuesta `@Failure 409` con `INSTANCE_INVALID_STATE` en los handlers de acciones de energía y se regeneró `swagger.json`.
> - Entregable 4 verificado y en verde en `etapa1_acceptance_test.go`.

- **Área:** Backend
- **Asignado:** Lucas
- **Estimación:** 0,5 h
- **Depende de:** `BAC-24A`.
- **Criterio de éxito:** Swagger documenta `INSTANCE_INVALID_STATE` en acciones de ciclo de vida.

#### `FIX-49` - `DELETE /api/instances/:vmid`: código de estado (D2) y `202` con seguimiento (`BAC-24B`) (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, backend `59148b1`, commits `088e175` y `f06640f`).**
> - Responde `409 INSTANCE_INVALID_STATE` si la instancia no está detenida.
> - Responde `202 { upid, tareaId }` si está detenida, encolando el UPID en `Seguir`, despachando auditoría `PENDING` y emitiendo `TASK_FINISHED` con `resource_type`.
> - Pasa 100% en verde en `etapa1_acceptance_test.go` y `puente_etapa1_acceptance_test.go`.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1,5 h
- **Depende de:** `BAC-24B`.
- **Criterio de éxito:** DELETE responde 409 si encendida, 202 con seguimiento si detenida y emite TASK_FINISHED.

#### `FIX-53` - Unificar códigos de acción y resultado en la auditoría de instancias (`BAC-24A`) (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, backend `59148b1`, commit `b476636`).**
> - Unifica vocabulario canónico `START`/`STOP`/`SHUTDOWN`/`REBOOT`/`DELETE`.
> - Define `ResultadoPendiente = "PENDING"` y preserva `instanciaNombre` y `resource_type` en alias `/start` y `/stop`.
> - Pasa la prueba `FIX-53_unificacion_de_codigos_de_accion_y_resultado_en_la_auditoria_de_instancias` en `fixes_acceptance_test.go`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1 h
- **Depende de:** `BAC-24A`.
- **Criterio de éxito:** Códigos de auditoría unificados entre rutas y alias.

#### `FIX-64` - Conexión SSE compartida y suscripción selectiva (`FRN-17A`) (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, frontend `5a86dce`).**
> - Implementa conexión compartida única vía `EventsProvider` en `ProtectedLayout`.
> - Suscripción y distribución de eventos por contexto (`useEventsContext`).
> - Pasan las pruebas correspondientes en `test/front/events-client.test.tsx`.

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1,5 h
- **Depende de:** `FRN-17A`, `FRN-17C`.
- **Criterio de éxito:** Una sola conexión SSE compartida por pestaña con suscripciones selectivas.
