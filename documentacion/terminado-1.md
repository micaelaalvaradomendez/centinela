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

#### `BAC-23B` - IP real en `GET /api/instances` y verificación del filtrado RBAC (`RF-03`)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, backend `df320c7`).**
> - `internal/core/services/instance_service.go` (`ListarInstancias`) y adaptadores resuelven `ip` de manera concurrente únicamente para las instancias visibles por el usuario (`Permitidas`), protegiendo la infraestructura contra accesos no autorizados y llamadas innecesarias al guest agent.
> - Preserva retrocompatibilidad para `FRN-07`/`FIX-14` y control RBAC: usuarios `OPERATOR` solo visualizan sus instancias asignadas y usuarios `ADMIN` acceden al inventario global del nodo.
> - Pasan las pruebas de aceptación en `test/back/etapa1_ola2_acceptance_test.go` (`BAC-23B_IP_real_en_GET_api_instances_y_filtrado_RBAC_estricto`).

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `BAC-23A` y `FIX-39`.
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

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, backend `df320c7`).**
> - En `GET /api/node/status` se incluye `instancesSummary` con desglose segregado de VMs y LXC (`running`, `stopped`, `paused`, `total`), integrado con la telemetría del nodo y cache en Redis.
> - En `GET /api/instances` se completan métricas de recursos en tiempo real (`cpuUsage`, `ramUsage`, `maxRam`), devolviendo valores nulos o 0 cuando la instancia está detenida.
> - Pasan las pruebas de aceptación en `test/back/etapa1_ola2_acceptance_test.go` (`BAC-22B_instancesSummary_en_GET_api_node_status_y_metricas_por_instancia`).

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-22`, `BAC-23A` y `FIX-39`.
- **Entregable:**
  1. En `GET /api/node/status`, incluir `instancesSummary: { vms: { running, stopped, paused, total }, lxc: { running, stopped, paused, total } }`, cacheado junto con la telemetría.
  2. En `GET /api/instances`, completar `cpuUsage` (0-100), `ramUsage` y `maxRam` (bytes) con los datos de `BAC-23A`.
- **Criterio de éxito:** `GET /api/node/status` devuelve el desglose por estado coincidiendo con el simulador, y `GET /api/instances` devuelve CPU y RAM de cada instancia (`null` si está apagada).

#### `FRN-16B` (`BRG-04-FRN`) - Resincronización del estado "Operación en progreso" tras recarga (`F5`) o reconexión

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, frontend `fc6f9ed`).**
> - Al cargar la tabla de instancias (`/instances`), las filas con `activeTask !== null` se inicializan inmediatamente en estado `transitioning` con spinner y botones de acción bloqueados/deshabilitados, conservando el `tareaId`.
> - Al reconectarse el canal SSE `/api/events`, se revalida el listado en segundo plano sin parpadeos de carga para sincronizar estados terminados durante la desconexión.
> - Pasan las pruebas en `test/front/instances-ola2.test.tsx` (`FRN-16B - Resincronización del estado de operación en progreso tras recarga (activeTask)`).

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25C`, `FRN-16` y `FRN-17C`.
- **Entregable:**
  1. Al cargar la tabla, poner en `transitioning` (spinner y botones bloqueados) las filas cuyo ítem de `GET /api/instances` trae `activeTask !== null`, guardando su `tareaId`.
  2. Al reconectarse `/api/events` (`FRN-17C`), volver a pedir `GET /api/instances` sin mostrar carga, para tomar los estados finales de las tareas que terminaron durante la desconexión.
- **Criterio de éxito:** si se recarga con `F5` en medio de una acción, la fila sigue con el spinner y los controles bloqueados hasta que llega `TASK_FINISHED` o termina la tarea.


### `FIX-69` - Restaurar definiciones SSE, esquemas de eventos y contratos en `swagger.yaml` / docs (`BAC-25A` / `BAC-29`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro / Lucas
- **Estado:** Pendiente en el submódulo `backend`. Correr `go test ./docs/...` en el backend falla con errores de esquemas faltantes.
- **Estimación:** 1,5 h
- **Depende de:** `BAC-25A` y `BAC-29` (en `terminado-1.md`).
- **Problema técnico detallado:**
  Al regenerar la documentación de Swagger con `swag init` en el commit `f06640f` del backend, se sobreescribieron `docs/swagger.yaml`, `docs/swagger.json` y `docs/docs.go`. Como los handlers no contenían anotaciones declarativas completas para los modelos de streaming SSE ni para el ticket efímero, se perdieron del archivo YAML/JSON las siguientes definiciones esenciales de la Etapa 1:
  - `http.SSEEventPayload`
  - `http.SSEDetallesEvento`
  - `http.TaskSuccess`
  - `http.TaskFailed`
  - `http.SSETicketResponse`
  - `http.SSECierrePayload`
  - `http.SSETiposPayload`
  Además, `contrato_etapa1_test.go` falla porque la respuesta 200 de `/node/status` quedó documentada con el tipo interno `#/definitions/internal_adapters_primary_http.EstadoNodeResponse` en lugar de la definición de contrato esperada `#/definitions/http.EstadoNodeResponse`.
- **Qué debe hacer el equipo de backend:**
  1. En los handlers de HTTP (`events_handler.go` y `node_handler.go`), incorporar o corregir las anotaciones de Swagger (`@Success`, `@Failure`, `@Produce`, `@Router` y tipos `@Model` o referencias a los DTOs de eventos).
  2. Asegurar que las estructuras de DTO de eventos SSE (`RealtimeEvent`, `TaskFinishedDetails`, `TicketResponse`, etc.) estén expuestas y referenciadas en las anotaciones de swag para que se generen en `docs/swagger.yaml` y `docs/swagger.json`.
  3. Ejecutar la regeneración (`swag init` con los flags de formato adecuados, o configurar alias de paquetes) para que los nombres de los esquemas en `#/definitions` no contengan prefijos de rutas internas (`internal_adapters_primary_http`).
  4. Ejecutar localmente `go test -v ./docs/...` en el backend hasta que todas las aserciones de artefactos (`TestContratoEtapa1ArtefactosJSONYAML`, `TestContratoEtapa1ArtefactosJSONGo` y `TestContratoEtapa1`) pasen en verde.
- **Criterio de éxito:**
  - `go test -v ./docs/...` en el repositorio backend pasa 100% en verde sin fallos de definiciones faltantes.
  - `docs/swagger.yaml` y `docs/swagger.json` contienen los esquemas completos de SSE y telemetría de nodo.


### `FIX-65` - Explicitar en `docs/contrato-etapa1.md` el acceso de cualquier usuario autenticado a `/node/status` (D1) (`BAC-29` / `FIX-46`) (Backend)

- **Área:** Backend / documentación
- **Asignado:** Nico
- **Estado:** Pendiente en el submódulo `backend`. La suite de tests de este repo ya adecuó el matcher multilínea, pero el documento de contrato debe explicitar la especificación formal del requisito D1.
- **Estimación:** 0,5 h
- **Depende de:** `FIX-46` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En `backend/docs/contrato-etapa1.md`, la sección `GET /api/node/status` describe el formato JSON de la respuesta y la telemetría, pero no define con precisión formal en el encabezado del endpoint la matriz de control de acceso requerida por el entregable 5 de `BAC-29` (D1). El texto actual menciona de forma genérica en párrafos separados que ambas operaciones usan `Authorization: Bearer <accessToken>`, pero no explicita que `GET /api/node/status` es accesible tanto para el rol `ADMIN` como para el rol `OPERATOR` sin restricciones de permisos de instancia.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/docs/contrato-etapa1.md`, en la sección correspondiente a `GET /api/node/status`, agregar explícitamente:
     > **Control de acceso (D1):** Endpoint accesible para cualquier usuario autenticado (`Authorization: Bearer <accessToken>`), habilitado tanto para el rol `ADMIN` como para el rol `OPERATOR`. No requiere permisos específicos sobre instancias.
  2. Mantener la alineación entre la versión Markdown y el archivo `.docx` correspondiente si aplica.
- **Criterio de éxito:**
  - `backend/docs/contrato-etapa1.md` especifica formalmente el requisito D1.
  - La prueba de contrato `TestEtapa1/BAC-29` valida el texto contractual y pasa 100% en verde.



### `FIX-63` - Cobertura de pruebas unitarias para `ObtenerInterfaces` en el cliente Proxmox (`BAC-23A`) (Backend)

- **Área:** Backend / pruebas
- **Asignado:** Lisandro
- **Estado:** Pendiente en el submódulo `backend`.
- **Estimación:** 1 h
- **Depende de:** `BAC-23A` (en `terminado-1.md`).
- **Problema técnico detallado:**
  El método `ObtenerInterfaces(ctx context.Context, nodo, tipo string, vmid int)` en `backend/internal/adapters/secondary/proxmox/client.go:474` implementa el consumo HTTP de:
  - QEMU VMs: `/api2/json/nodes/{node}/qemu/{vmid}/agent/network-get-interfaces`
  - LXC Containers: `/api2/json/nodes/{node}/lxc/{vmid}/interfaces`
  Sin embargo, en `backend/internal/adapters/secondary/proxmox/client_test.go` **no se escribió ninguna prueba unitaria para `ObtenerInterfaces`**. La suite de aceptación de `BAC-23A` (`etapa1_acceptance_test.go:171`) verifica que existan pruebas unitarias para ambas rutas y falla con:
  `no hay pruebas unitarias del adaptador de IP (ningún _test.go del backend menciona network-get-interfaces ni /interfaces)`.
  Aunque `inventario_service_test.go` prueba la lógica interna con stubs en memoria, el adaptador HTTP contra Proxmox quedó sin cobertura.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/adapters/secondary/proxmox/client_test.go`, agregar casos de prueba con `httptest.Server`:
     - `TestClient_ObtenerInterfaces_Qemu`: verifica que una VM consulte la ruta `/api2/json/nodes/{node}/qemu/{vmid}/agent/network-get-interfaces` y parsee correctamente las interfaces y direcciones IPv4/IPv6 devueltas por el Guest Agent (con `prefix` entero).
     - `TestClient_ObtenerInterfaces_LXC`: verifica que un contenedor consulte la ruta `/api2/json/nodes/{node}/lxc/{vmid}/interfaces` y parsee las interfaces y direcciones (soportando `prefix` string o número, y devolviendo vacío sin error si `data: null` por estar detenido).
     - `TestClient_ObtenerInterfaces_ErrorYTimeout`: verifica el manejo de errores (ej. VM sin guest agent respondiendo error o 500) devolviendo `ErrGuestAgentNoDisponible` o `nil` controlado según el contrato.
- **Criterio de éxito:**
  - `go test -v ./internal/adapters/secondary/proxmox/...` pasa al 100% ejecutando los nuevos tests.
  - La suite de aceptación `TestEtapa1/BAC-23A` encuentra las pruebas en `client_test.go`, las ejecuta y pasa 100% en verde.


---

# Verificación del 09/10/2026: tareas movidas desde `actual.md`

Revisiones probadas: backend `4f68e44` y frontend `ee4644c`.

#### `FRN-15` - Modales de confirmación antierror para acciones operativas

> [!NOTE]
> **Estado: Completada (verificado el 09/10/2026, frontend `ee4644c`).**
> - Se implementaron modales de confirmación diferenciados por criticidad: `Start` (confirmación estándar), `Shutdown`/`Reboot` (aviso de apagado/reinicio de OS huésped), `Stop` (advertencia en rojo sobre posible pérdida de datos) y `Delete` (modal destructivo exigiendo tipear el nombre o ID de la instancia).
> - La eliminación destructiva (`Delete`) queda restringida exclusivamente al rol `ADMIN`. Los operadores con `READ_ONLY` no ven controles mutantes; usuarios con rol `ADMIN` siempre están habilitados a operar.
> - Pasan 6 de 6 pruebas en `test/front/instances-modals.test.tsx`.

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.5 h
- **Depende de:** `FRN-20A`, `SEC-03` (`PermissionGate`), `FIX-37` y `FIX-38`.
- **Criterio de éxito:** Ninguna acción operativa se dispara sin pasar por el modal; cancelar no emite peticiones; `OPERATOR` con `READ_ONLY` no ve energía y nunca ve `Delete`; `ADMIN` puede operar y eliminar.


#### `FRN-19B` (ex `FRN-13B`) - Semáforo de salud global e integración con `GET /api/node/status` (`RF-02`)

> [!NOTE]
> **Estado: Completada (verificado el 09/10/2026, frontend `ee4644c`).**
> - `Dashboard.tsx` conecta los medidores a `GET /api/node/status`: porcentaje de CPU, memoria RAM, disco/almacenamiento, uptime formateado en días/horas/minutos y resumen global de instancias.
> - Semáforo de salud visual implementado: `Saludable` (< 70%), `Advertencia` (>= 70%) e `Inaccesible` ante respuestas 502/504 o caída de conexión.
> - Muestra advertencia visual de datos desactualizados si la respuesta contiene `stale: true`.
> - Pasan 12 de 12 pruebas en `test/front/dashboard-node.test.tsx`.
> - Observación menor sobre desduplicación de texto en lectores de pantalla canalizada en `FIX-78` de `futuro-1.md`.

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** `FRN-19A`, `BAC-29`, `BAC-22`.
- **Criterio de éxito:** Los datos reales del nodo se visualizan en pantalla y se actualizan sin `F5`. Ante caída de backend la interfaz muestra `Inaccesible` sin romperse.


#### `FRN-16` - Máquina de estados "Operación en progreso" por instancia

> [!NOTE]
> **Estado: Implementada con observaciones (verificado el 09/10/2026, frontend `ee4644c`).**
> - Al confirmar una acción operativa, se despacha la orden y la fila entra en estado `transitioning`, mostrando un spinner en el botón accionado y bloqueando todos los controles de esa instancia para evitar órdenes concurrentes.
> - **Observaciones detectadas:** El mapeo diferenciado de los 6 códigos de error acordados en el contrato D2 (`BAC-29`) (`INSTANCE_INVALID_STATE`, `INSTANCE_BUSY`, `INSTANCE_PROTECTED`, `INSTANCE_ACCESS_DENIED`, `PROXMOX_UNAVAILABLE`, `PROXMOX_TIMEOUT`) y el desbloqueo interactivo inmediato de la fila tras un fallo se canalizan en [`FIX-77`](futuro-1.md#fix-77---mapeo-de-errores-http-de-ciclo-de-vida-bac-29-y-resiliencia-de-la-máquina-de-estados-frn-16-frontend) de `futuro-1.md`.

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.5 h
- **Depende de:** `FRN-15`, `BAC-29`, `FIX-39`, `FIX-40`, `BAC-24A` y `BAC-24B`.
- **Criterio de éxito:** Es imposible disparar una segunda acción sobre la misma instancia mientras hay una orden en curso. Los errores devuelven feedback visual claro.


#### `FIX-52` - Distinguir `INSTANCE_BUSY` de `INSTANCE_INVALID_STATE` y limpiar código muerto (`BAC-24A`) (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 09/10/2026, backend `4f68e44`).**
> - Commit `222760d`: `validarEstadoParaAccion` en `instance_handler.go` consulta `tareas_asincronas` antes de evaluar la matriz de estados; si existe una tarea `RUNNING` para el VMID responde inmediatamente `409 INSTANCE_BUSY`.
> - Se eliminó el método huérfano `ReiniciarInstancia` en `ProxmoxPort` y `client.go`.
> - Se corrigió el comentario en `EliminarInstancia` reflejando que la validación de estado detenido la realiza el handler.
> - Pasan las pruebas de regresión en `test/back/fixes_acceptance_test.go`.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Depende de:** `BAC-24A`.
- **Criterio de éxito:** Una acción consecutiva sobre una instancia con tarea en curso responde `409 INSTANCE_BUSY`; sobre una ya detenida/incompatible responde `409 INSTANCE_INVALID_STATE`.


#### `BAC-25B` - Timeout configurable, reintentos y `exitstatus` en el evento

> [!NOTE]
> **Estado: Completada (verificado el 09/10/2026, backend `4f68e44`).**
> - Se incorporó la variable de entorno `UPID_TIMEOUT` con valor por defecto de 3 minutos (configurable a 5 s en pruebas); al expirar, la tarea pasa a `FAILED` con `motivo: TIMEOUT`.
> - Se implementaron reintentos con retroceso exponencial ante fallos transitorios en la consulta de estado de tareas en Proxmox.
> - El evento `TASK_FINISHED` incluye `exitstatus` y `motivo` (`PROXMOX_ERROR` o `TIMEOUT`) con su mensaje descriptivo correspondiente.
> - Pasan todas las pruebas de la Ola 2 en `test/back/etapa1_ola2_acceptance_test.go`.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25A`.
- **Criterio de éxito:** Seguimiento acotado por `UPID_TIMEOUT`, reintentos resilientes y payload completo de `TASK_FINISHED` según contrato `BAC-29`.


#### `FRN-20B` (ex `FRN-14B`) - Filtros reactivos por tipo, estado y buscador dinámico

> [!NOTE]
> **Estado: Implementada con observaciones (verificado el 09/10/2026, frontend `ee4644c`).**
> - En `Instances.tsx`, la barra superior incluye buscador en tiempo real por ID y nombre, selector por estado (Todas, En ejecución, Detenidas) y contador de instancias reactivo que opera instantáneamente en memoria.
> - **Observación detectada:** El selector reactivo por tipo de recurso (VM / LXC) fue implementado dentro de un menú desplegable contextual (`DropdownMenu`) en lugar de botones o pestañas directamente visibles y accesibles en la barra superior; canalizado en [`FIX-79`](futuro-1.md#fix-79---exposición-directa-de-controles-para-selector-reactivo-por-tipo-de-recurso-frn-20b-frontend) de `futuro-1.md`.

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.0 h
- **Depende de:** `FRN-20A`.
- **Criterio de éxito:** El filtrado es instantáneo en memoria del cliente sin peticiones al backend.


#### `FIX-48` - Documentar `INSTANCE_INVALID_STATE` en Swagger (`BAC-24A`) (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 09/10/2026, backend `4f68e44`).**
> - Las especificaciones OpenAPI/Swagger (`backend/docs/swagger.json` y `swagger.yaml`) formalizan el código de respuesta HTTP `409` con clave `INSTANCE_INVALID_STATE` para todas las acciones de ciclo de vida (`start`, `stop`, `shutdown`, `reboot`, `delete`).

- **Área:** Backend
- **Asignado:** Lucas
- **Estimación:** 0,5 h
- **Depende de:** `BAC-24A`.
- **Criterio de éxito:** Swagger documenta formalmente `INSTANCE_INVALID_STATE`.


#### `FIX-62` - Heurística de circuit breaker y recuperación de telemetría en `nodo_service` (`BAC-22`) (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 09/10/2026, backend `4f68e44`).**
> - Commit `eca422d`: En `backend/internal/core/services/nodo_service.go`, se condicionó la penalización `fallaReciente()`: si no existe lectura previa en caché (`claveNodoUltimo`), no bloquea peticiones durante 5 segundos a ciegas.
> - Se limpió el estado de fallo permitiendo la recuperación instantánea de telemetría tan pronto como Proxmox vuelve a responder.
> - Pasan las pruebas unitarias en `nodo_service_test.go` y la prueba de integración `TestEtapa1/BAC-22`.

- **Área:** Backend
- **Asignada:** Tayra / Lisandro
- **Estimación:** 1 h
- **Depende de:** `BAC-22`.
- **Criterio de éxito:** El backend se recupera inmediatamente tras el restablecimiento del hipervisor sin quedar en un blackout de 5 segundos cuando no hay caché previa disponible.


