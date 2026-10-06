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

#### `FRN-19A` - Componente visual de medidores de recursos del host (Frontend)

> [!NOTE]
> **Estado: Completada (verificado el 03/10/2026, frontend `907efe5`, PR #78).**
> - Se implementó `ResourceMeter.tsx` en `features/dashboard` con barras de progreso para CPU, RAM y almacenamiento, con umbrales visuales al 70%.
> - Integrado en `Dashboard.tsx` y validado con suite unitaria `resource-meters.test.tsx` y suite de aceptación `dashboard-metrics.test.tsx` (3/3 en verde).

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** ninguna.

---

# Verificación del 06/10/2026

Revisiones probadas: backend `eec77ff` (commits `3988546`, `5787179`, `13f9c35`, `d1dec4f`), frontend `5a86dce` (PR #80, #81, #82, #83, #84).

#### `BAC-25A` - Worker pool acotado para el seguimiento de UPID (Backend)

> [!NOTE]
> **Estado: Completada (verificado el 06/10/2026, backend `eec77ff`, commit `3988546`).**
> - Se implementó el pool acotado de workers (`UPID_WORKERS`, por defecto 8) con canal bufferizado en `internal/services/seguimiento_tareas.go`.
> - La función `Seguir` no bloquea: si el buffer está lleno, la tarea queda persistida en `tareas_asincronas` con estado `RUNNING` para posterior rescate sin pérdida.
> - Cierre ordenado respetando el contexto de aplicación (`context.Context`).
> - **Pruebas:** `TestEtapa1/BAC-25A_el_seguimiento_de_UPID_usa_un_pool_acotado_y_no_pierde_tareas` en `test/back/etapa1_acceptance_test.go` pasa **100% en verde** (53.16s).

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** ninguna pendiente. Parte de `services/seguimiento_tareas.go` (`BAC-21C`).
- **Contexto:** cada tarea lanzaba su propia goroutine sin límite.
- **Entregable:**
  1. Reemplazar la goroutine por tarea con un pool de N workers (`UPID_WORKERS`, por defecto 8) que lee de un canal con buffer.
  2. `Seguir` sigue siendo no bloqueante: si el canal está lleno, la tarea queda en `tareas_asincronas` como `RUNNING` y no se pierde.
  3. Cerrar ordenadamente con el contexto del proceso.
  4. Pruebas unitarias con un Proxmox falso.
- **Criterio de éxito original:** con 50 tareas simultáneas nunca hay más de N consultas a Proxmox en paralelo; el fin de cada tarea se detecta en menos de 2 s; pruebas de `seguimiento_tareas_test.go` en verde.

#### `BAC-29` - Contrato HTTP y de eventos de la Etapa 1 (Backend)

> [!WARNING]
> **Estado: Implementada con problemas (verificado el 06/10/2026, backend `eec77ff`, commits `5787179`, `21b1332`, `4ef36e3`, `0850ec7`).**
> - **Completado:**
>   - Se redactó y publicó `backend/docs/contrato-etapa1.md` detallando las especificaciones de `GET /api/node/status`, `GET /api/instances`, rutas de energía (`shutdown`, `reboot`), estructura de `TASK_FINISHED`, matriz de errores de la etapa y permisos.
>   - Se incorporaron las anotaciones Swagger para los nuevos endpoints de la etapa.
> - **Problemas detectados (se corrigen en `FIX-46` en `futuro-1.md`):**
>   - La prueba de aceptación `etapa1_acceptance_test.go:51` busca validar mediante expresión regular que la sección de `GET /node/status` en `contrato-etapa1.md` documente explícitamente en su bloque que cualquier usuario autenticado (`ADMIN` u `OPERATOR`) puede consultarlo. Actualmente la redacción ubica la especificación de permisos en una sección separada del encabezado, haciendo fallar la verificación sintáctica.
> - **Pruebas:** `TestEtapa1/BAC-29_contrato_etapa1_incluye_node_status_y_codigos_de_error` falla la aserción de documentación de permisos por línea.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.0 h
- **Depende de:** ninguna.
- **Entregable:** `backend/docs/contrato-etapa1.md` y anotaciones Swagger de `GET /api/node/status`, `GET /api/instances`, rutas de energía, detalles de `TASK_FINISHED`, matriz de códigos de error y autorizaciones.
- **Criterio de éxito original:** el frontend puede maquetar las vistas de la etapa con el documento y Swagger muestra los endpoints con ejemplos.

#### `BAC-22` - Adaptador de telemetría del nodo con caché en Redis (`RF-02`) (Backend)

> [!WARNING]
> **Estado: Implementada con problemas (verificado el 06/10/2026, backend `eec77ff`, commit `13f9c35`).**
> - **Completado:**
>   - Se creó el endpoint `GET /api/node/status` (`node_handler.go`, `telemetria_service.go`) protegido con autenticación (`RequireAuth`).
>   - Integra `ObtenerEstadoNodo` en `ProxmoxPort` y `client.go`, normalizando métricas de CPU (porcentaje y cores), RAM y disco (conversión a GB), y uptime.
>   - Implementa almacenamiento en caché Redis con TTL corto y soporte para devolver el último estado conocido (`stale: true`) ante indisponibilidad.
> - **Problemas detectados (se corrigen en `FIX-47` en `futuro-1.md`):**
>   - En la prueba `etapa1_acceptance_test.go:89`, tras simular una caída temporal de Proxmox y restablecer el stub mock, una llamada posterior a `GET /api/node/status` retorna `502 PROXMOX_UNAVAILABLE` en vez de refrescar la caché y responder `200 OK` con datos actualizados.
> - **Pruebas:** `TestEtapa1/BAC-22_telemetria_de_nodo_con_cache_redis_y_stale_reading` falla en la verificación de recuperación tras fallo.

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.5 h
- **Depende de:** `INF-06A` y `BAC-17A`.
- **Entregable:** endpoint `GET /api/node/status` con normalización de métricas, caché en Redis con TTL de 5 a 10s, clave de fallback sin TTL (`stale: true`), y manejo de errores.
- **Criterio de éxito original:** responde en <50ms con caché; refresca con 200 al vencer TTL; ante caída responde stale con 200; 401 sin token y 200 para OPERATOR/ADMIN.

#### `BAC-23A` - Adaptador y normalización de inventario Proxmox (QEMU / LXC / IP) (Backend)

> [!WARNING]
> **Estado: Implementada con problemas (verificado el 06/10/2026, backend `eec77ff`, commit `d1dec4f`).**
> - **Completado:**
>   - Implementado en `internal/services/inventario_service.go` y `internal/adapters/secondary/proxmox/client.go` unificando recursos de VMs y LXC en una estructura común con CPU, memoria y resolución de IPs.
>   - Soporte para consultar interfaces de red de QEMU y LXC, seleccionando la primera IPv4 no loopback o primera IPv6 global.
> - **Problemas detectados (se corrigen en `FIX-48` en `futuro-1.md`):**
>   - La suite de aceptación estática/dinámica (`etapa1_acceptance_test.go:156`) inspecciona la presencia explícita de los literales de ruta `network-get-interfaces` o `/interfaces` en la configuración de peticiones de clientes o stubs de test unitarios. Como la implementación utiliza constantes o métodos de abstracción interna en el cliente HTTP, la aserción de prueba estática falla al no hallar los strings literales esperados en los archivos de prueba unitaria.
> - **Pruebas:** `TestEtapa1/BAC-23A_adaptador_de_inventario_resuelve_IP_de_VM_y_LXC_en_paralelo` falla en la verificación de endpoints de interfaces en tests unitarios.

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `BAC-14`, `BAC-28` y `FIX-35`.
- **Entregable:** servicio en Go que unifica VMs y LXC, incluye telemetría cruda y resuelve IP mediante guest agent o interfaces LXC concurrentemente con timeout.
- **Criterio de éxito original:** función interna con pruebas unitarias que devuelve lista consolidada; instancias sin IP devuelven `null` sin error; consultas paralelas con timeout acotado.

#### `FRN-17A` - Consumo de eventos en tiempo real y distribución por instancia (Frontend)

> [!WARNING]
> **Estado: Implementada con problemas (verificado el 06/10/2026, frontend `5a86dce`, PR #80, #81, #82, commits `28ae029`, `74ea3d5`, `e8810c9`).**
> - **Completado:**
>   - Se implementó `EventsContext.tsx` y el hook `useEvents.ts`, proveyendo contexto para el consumo de eventos SSE (`RealtimeEvent`) y gestión de estado de la conexión (`estado`, `ultimoMensaje`, `error`).
> - **Problemas detectados (se corrigen en `FIX-49` en `futuro-1.md`):**
>   - `useEvents` no expone una interfaz de suscripción (`subscribe`, `events`) que permita a componentes suscribirse por tipo de evento o por `recursoId` / instancia.
>   - No implementa deduplicación de eventos por campo `id`.
>   - Múltiples componentes que invocan `useEvents` abren conexiones independientes a `/api/events` en lugar de compartir una única conexión viva gestionada por el provider.
> - **Pruebas:** en `test/front/events-client.test.tsx`, fallan 3 casos: *"deduplica eventos por id para no reprocesar"*, *"mantiene una sola conexión viva para toda la app"* y *"permite a los componentes suscribirse a eventos de una instancia"*.

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1.5 h
- **Depende de:** `FRN-17C`.
- **Entregable:** cliente de eventos sobre la conexión de `FRN-17C`, parseo de `RealtimeEvent`, deduplicación por `id`, suscripción por tipo y recurso, y conexión compartida única por pestaña.
- **Criterio de éxito original:** Dashboard y tabla reciben `TASK_FINISHED` en tiempo real con una conexión por pestaña; eventos duplicados se procesan una sola vez; pruebas de `events-client.test.tsx` en verde.

