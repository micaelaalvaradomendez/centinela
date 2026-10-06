### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 03/10/2026** (backend `4e204f1`, frontend `907efe5`; detalle en [test/informe.md](../test/informe.md)).
> - **Pasaron a [`terminado.md`](terminado.md):**
>   - `FIX-29` y `FIX-40`, completas;
>   - `SEC-03`, `FIX-41`, `FIX-42` y `FIX-38`, completas;
>   - `FIX-36`, completa (PR #77, frontend `070e96b`);
>   - `FIX-37` y `FIX-39`, completas (backend `4e204f1`);
>   - `BAC-18B`, implementada con problemas en backend `4e204f1` (particionamiento trimestral OK, purga periódica OK; falta indexar `jti_access`, derivado a `FIX-44`, en este archivo).
> - **Pasaron a [`terminado-1.md`](terminado-1.md):**
>   - `INF-07B`, referencia local verificada (pendiente despliegue CT 103 en `INT-03`);
>   - `FRN-20A`, implementada con problemas en frontend `070e96b` (integración a la API y modales implementados; maquetado, IP y badges derivados a `FIX-43` en `futuro-1.md`);
>   - `FRN-19A`, completa (PR #78, frontend `907efe5`, medidores de host en `ResourceMeter.tsx`).
> - **Verificación del 05/10/2026** (backend `e1f2df4`, frontend `d47999d`; detalle en [test/informe.md](../test/informe.md)). Ninguna tarea está completa con todas sus pruebas en verde.
>   - **Pasaron a [`terminado.md`](terminado.md)**, implementada con problema: `FRN-17C`. Su corrección es `FIX-45`, en [`futuro.md`](futuro.md).
>   - **Pasaron a [`terminado-1.md`](terminado-1.md)**, implementadas con problema: `BAC-29`, `BAC-24A`, `BAC-24B` y `FRN-17A`. Sus correcciones son `FIX-46` a `FIX-49`, en [`futuro-1.md`](futuro-1.md).
>   - **Revisión manual de Lucas del 05/10/2026** sobre `BAC-24A`/`BAC-18B` (backend `e1f2df4`): agrega 4 problemas no cubiertos por `FIX-46` a `FIX-49`. Dos son de `BAC-18B` (`FIX-50` apagado del servidor HTTP, `FIX-51` auditoría previa a la migración invisible, ambos en este archivo); dos son de `BAC-24A` (`FIX-52` distinción `INSTANCE_BUSY`/`INSTANCE_INVALID_STATE` y código muerto, `FIX-53` códigos de auditoría mezclados, ambos en [`futuro-1.md`](futuro-1.md)). El resto de sus puntos ya está cubierto (`FIX-49`, borrado asíncrono) o es una decisión de diseño a discutir, no un bug (purga inmediata de sesiones cerradas; bloqueo de `start` en VMIDs `100`-`105`).
>   - **Quedan en este archivo**, sin implementar:
>
> | Tarea | Área | Estado | Pruebas |
> |---|---|---|---|
> | `FIX-44` | Backend | No implementado: el índice parcial no incluye `jti_access`. La purga horaria ya cumple | `cierre_fase_base…` BAC-18B (2; falla 1) |
> | `FIX-50` | Backend | No implementado: `router.Run` no observa la señal de apagado, el proceso no termina solo | Detectado por revisión manual (Lucas, 05/10) |
> | `FIX-51` | Backend | No implementado: `auditoria_legacy` no se consulta desde la API tras particionar | Detectado por revisión manual (Lucas, 05/10) |
> | `BAC-22` | Backend | No implementada (`/node/status` responde 404) | `etapa1_acceptance…` (1) |
> | `BAC-23A` | Backend | No implementada | `etapa1_acceptance…` (1) |
> | `BAC-25A` | Backend | No implementada | `etapa1_acceptance…` (1) |
> | `FIX-43` | Frontend | No implementado | `instances-table.test.tsx` (5) |
> | `FRN-15` | Frontend | No implementada | `instances-modals.test.tsx` (6) |
> | `FRN-19B` | Frontend | No implementada | `dashboard-node.test.tsx` (11) |
> | `FRN-16` | Frontend | No implementada | `instances-operation.test.tsx` (9) |
> | `BAC-25C` | Backend | No implementada: no retoma las tareas al reiniciar (el `activeTask` que ya se informa vino con `FIX-39`) | `etapa1_acceptance…` (1) |

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

> **Nota de alcance (06/10/2026, revisión estática):** “recorrido completo” se refiere a los casos anteriores, no a crear una cuenta nueva con el correo de una eliminada y preservar la identidad histórica. Esta brecha de la fase base queda pendiente en `FIX-54` a `FIX-57` de [`futuro.md`](futuro.md#fix-54---separar-eliminación-lógica-y-suspensión-migrar-unicidad-del-correo-backend). No se ejecutaron suites ni se verificó el esquema desplegado.

---

### `FIX-44` - Corregir índice parcial en `sesiones_activas` para incluir `jti_access` y alinear worker de purga (`BAC-18B`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra (autora de `96106a6`)
- **Estimación:** 0,5 h
- **Depende de:** `BAC-18B` (en `terminado.md`).
- **Problema y evidencia:**
  1. `internal/adapters/secondary/postgres/db.go:148` creó el índice parcial como:
     ```sql
     CREATE INDEX IF NOT EXISTS idx_sesiones_activas_vigentes
     ON sesiones_activas (usuario_id, fecha_expiracion)
     WHERE activa = true;
     ```
     Omitió la columna `jti_access`. Como el middleware de autenticación (`AuthMiddleware`) busca las sesiones por `jti_access` para comprobar revocaciones en cada solicitud, el índice no cubre la consulta de alta concurrencia. La suite de pruebas (`cierre_fase_base_acceptance_test.go:396`) falla por no encontrar un índice parcial que cubra `jti_access`.
  2. En `internal/adapters/secondary/postgres/purga_worker.go:39`, el worker usa `ticker := time.NewTicker(w.intervalo)`. La prueba estática de aceptación busca la inicialización con `time.Hour` o equivalente dentro del worker (`time.NewTicker(time.Hour)`).
- **Entregable:**
  1. Actualizar la definición del índice parcial en `db.go` para que indexe `(jti_access, usuario_id)` o `jti_access` con `WHERE activa = true` (por ejemplo `ON sesiones_activas (jti_access, usuario_id) WHERE activa = true;`).
  2. Ajustar `purga_worker.go` para que instancie el ticker con `time.NewTicker(1 * time.Hour)` (o mantenga el default con `time.NewTicker(time.Hour)`).
- **Criterio de éxito:**
  - `TestCierreFaseBase/BAC-18B_indice_parcial_en_sesiones_activas_y_auditoria_particionada_por_trimestre` y `TestCierreFaseBase/BAC-18B_rutina_horaria_que_purga_sesiones_inactivas_o_vencidas_de_sesiones_activas` pasan 100% en verde.

---

### `FIX-50` - Apagado correcto del servidor HTTP ante `SIGINT`/`SIGTERM` (`BAC-18B`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra (autora de `96106a6`)
- **Estimación:** 1 h
- **Depende de:** `BAC-18B` (en `terminado.md`).
- **Problema y evidencia (detectado en la revisión de Lucas del 05/10/2026, verificado enviando ambas señales):** `cmd/api/main.go:90` crea `ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)` para que los workers en segundo plano se detengan limpiamente, pero solo se lo pasa a `purgaWorker.Iniciar(ctx)` (línea 96). El servidor HTTP se levanta con `router.Run(":8080")` (línea 299), que internamente llama a `http.ListenAndServe` y bloquea el goroutine principal sin recibir nunca `ctx`. Al llegar `SIGINT` o `SIGTERM`: el worker de purga corta su loop y lo loguea, pero el servidor HTTP sigue aceptando conexiones indefinidamente. En el servidor, cada deploy espera unos 90 segundos hasta que el sistema operativo mata el proceso a la fuerza; en local hay que matarlo a mano.
- **Entregable:**
  1. Reemplazar `router.Run(":8080")` por un `http.Server{Addr: ":8080", Handler: router}` explícito.
  2. Arrancar ese servidor en una goroutine con `ListenAndServe` y, en el goroutine principal, esperar `<-ctx.Done()` para disparar `srv.Shutdown(shutdownCtx)` con un timeout acotado (por ejemplo 10 s) que deje drenar las requests en curso.
  3. Verificar que, tras `srv.Shutdown`, el proceso termina solo (sin quedar colgado) y que el log muestra el apagado del worker de purga y del servidor HTTP en el mismo cierre.
- **Criterio de éxito:** enviar `SIGINT` o `SIGTERM` al proceso hace que termine por sí mismo en pocos segundos, sin necesidad de `kill -9` ni de esperar al timeout forzado del orquestador de despliegue.

---

### `FIX-51` - El historial de `auditoria` previo a la migración queda inaccesible desde la API (`BAC-18B`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra (autora de `96106a6`)
- **Estimación:** 1 h
- **Depende de:** `BAC-18B` (en `terminado.md`).
- **Problema y evidencia (detectado en la revisión de Lucas del 05/10/2026, verificado en PC local y en el servidor):** `migrarAuditoriaParticionada` (`internal/adapters/secondary/postgres/db.go:216`) detecta la tabla `auditoria` plana preexistente y la renombra con `ALTER TABLE auditoria RENAME TO auditoria_legacy` antes de crear la tabla particionada desde cero. Ningún repositorio vuelve a referenciar `auditoria_legacy`: `auditRepositoryImpl.Listar` y `.ExportarRegistros` (`internal/adapters/secondary/postgres/audit_repository.go`) solo consultan `auditoria`. Los datos no se pierden (siguen en la tabla renombrada), pero el historial anterior a la migración (167 registros en el caso verificado, tanto en PC local como en el servidor) deja de verse en `GET /api/admin/audit` y en las exportaciones CSV/JSON, rompiendo la propiedad *append-only* visible de la auditoría.
- **Entregable:**
  1. Migrar los datos de `auditoria_legacy` hacia la tabla particionada (`INSERT INTO auditoria SELECT * FROM auditoria_legacy`, distribuyendo cada fila a su partición por `fecha_hora`) en lugar de solo renombrar, o bien incluir `auditoria_legacy` en una vista/consulta `UNION ALL` que use `Listar` y `ExportarRegistros`.
  2. Si se opta por la migración de datos, eliminar `auditoria_legacy` al terminar y documentar el paso en `migrarAuditoriaParticionada`.
  3. Agregar una prueba de aceptación que arranque con una tabla `auditoria` plana con filas de prueba y verifique que, tras la migración, esas filas siguen apareciendo en `GET /api/admin/audit`.
- **Criterio de éxito:** un registro de auditoría creado antes de la migración a particionado sigue siendo visible por `GET /api/admin/audit` y por las exportaciones, sin duplicados ni pérdida de datos.

### `FIX-45` - El cliente de eventos abre el stream en una ruta que el backend no sirve (`FRN-17C`) (Frontend)

- **Área:** Frontend
- **Asignado:** Nico (autor de `2a2db82`)
- **Estimación:** 0,5 h
- **Depende de:** `FRN-17C` (en `terminado.md`) y `BAC-21C`.
- **Problema y evidencia (verificado el 05/10/2026, frontend `d47999d`, backend `e1f2df4`):**
  - `services/eventsClient.ts:19` define `EVENTS_STREAM_PATH = '/events/stream'`, y el cliente se conecta a `/api/events/stream?ticket=<uuid>`.
  - El backend sirve el stream en `GET /api/events?ticket=<uuid>` (`cmd/api/main.go:260`, `BAC-21C`); `/api/events/stream` no existe y responde 404. Contra el backend real no llega ningún evento, aunque el ticket se pida bien.
  - El criterio de `FRN-17C` dice *"el frontend se conecta a `/api/events` usando tickets de un solo uso"*.
  - Prueba que falla: `test/front/events-client.test.tsx`, caso *"pide POST /api/events/ticket con Bearer y se conecta a /api/events?ticket=<uuid> sin exponer el JWT"*: `expected '/api/events/stream?ticket=ticket-uno' to match /\/api\/events\?ticket=ticket-uno$/`.
- **Entregable:** conectar el stream a `/api/events?ticket=<uuid>`, la ruta del backend y de `docs/contrato-eventos.md`. Si el equipo prefiere `/api/events/stream`, el cambio va en el backend y en el contrato, y hay que avisar para ajustar la prueba.
- **Criterio de éxito:**
  - Pasan los 4 casos `FRN-17C…` de `events-client.test.tsx`.
  - Con el backend real, el navegador abre `/api/events?ticket=…` y recibe el `TASK_FINISHED` de una acción de energía.

---
# ETAPA 1
> Las tareas de la Etapa 1 ya verificadas están en [`terminado-1.md`](terminado-1.md), y sus correcciones en [`futuro-1.md`](futuro-1.md).
---

> Tareas de la Ola 1 de [`etapa1.md`](etapa1.md), en curso. Ninguna está implementada todavía (verificado el 03/10/2026).

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

### `FIX-43` - Corregir maquetado, visualización de IP, badges de estado y columnas en tabla de instancias (`FRN-20A` / RF-03) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz / Cristian (PR #74)
- **Estimación:** 1.0 h
- **Depende de:** `FRN-20A` (en `terminado-1.md`).
- **Problema y evidencia:**
  El PR #74 implementó la integración viva de `Instances.tsx` con `useInstances.ts` y `instanceService.ts`, consumiendo `GET /instances` y aplicando permisos. Sin embargo, la suite de pruebas de aceptación (`test/front/instances-table.test.tsx`) falla 5/5 por discrepancias de maquetado e interfaz:
  1. **Nombre y VMID:** `Instances.tsx:113` renderiza `<p>{instance.name} ({instance.id})</p>`. La prueba y el diseño requieren que el nombre y el VMID se presenten como elementos claramente identificables en la fila (o celdas diferenciadas), permitiendo consultar `within(row).getByText('101')` y `screen.findByText('servidor-web')` de forma unívoca.
  2. **Columna de IP:** `Instances.tsx:125` tiene hardcodeado un guión fijo (`<td className="px-4 py-4">—</td>`). `instanceService.ts` debe leer `instance.ip` de la respuesta, y `Instances.tsx` debe renderizar la IP (ej. `192.168.1.50`) o el texto `"No detectada"` cuando sea `null`, junto a un botón interactivo para copiar la dirección al portapapeles (`navigator.clipboard.writeText`).
  3. **Badges de estado:** `Instances.tsx:119` renderiza `{instance.status}` como texto simple sin estilos distintivos. El criterio de aceptación exige badges diferenciados con estilos visuales estándar: verde para `running` / `en ejecución` y gris para `stopped` / `detenida`.
  4. **Botonera de acciones:** cada fila debe contar con su botonera de acciones presente en la tabla, independientemente de si la fila tiene IP o si las acciones operativas están condicionadas.
- **Entregable:**
  1. En `instanceService.ts`, agregar el campo opcional `ip?: string | null` en `InventoryInstance` y leerlo en el mapeo de `fetchInstanceInventory`.
  2. En `Instances.tsx`, ajustar el renderizado de la columna de nombre para que el texto del nombre (`instance.name`) y el VMID (`instance.id`) estén en elementos o nodos de texto separados.
  3. Renderizar la columna IP mostrando `instance.ip` si existe (con botón de copia con icono y `aria-label="Copiar IP"`) o `"No detectada"` si es `null`.
  4. Envolver el estado en un badge con clases de Tailwind que apliquen fondo y texto verde para `running` (ej. `bg-green-100 text-green-700` o variante shadcn correspondiente) y gris para `stopped`.
  5. Asegurar que la columna tipo exponga claramente `VM` o `LXC`.
- **Criterio de éxito:**
  - Los 5 casos de `test/front/instances-table.test.tsx` pasan 100% en verde.


#### `FRN-15` - Modales de confirmación antierror para acciones operativas
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.5 h
- **Depende de:**
  - `FRN-20A`;
  - `SEC-03` (`PermissionGate`);
  - `FIX-37` y `FIX-38`, para que `canOperateInstance` funcione con datos reales.

  Los modales se pueden maquetar antes, con datos de prueba.
- **Entregable:** modales según la criticidad de la acción, con `components/ui/modals.tsx`:
  - `Start`: confirmación estándar.
  - `Shutdown` / `Reboot`: aviso de apagado o reinicio del sistema operativo huésped.
  - `Stop`: advertencia en rojo sobre posible pérdida de datos.
  - `Delete`: modal destructivo que pide tipear el ID o el nombre. **Solo visible para ADMIN** (`PermissionGate`).

  Los botones de energía se muestran solo si `canOperateInstance(vmid)` da `true`.
- **Criterio de éxito:**
  - Ninguna acción se dispara sin pasar por el modal, y cancelar no envía peticiones.
  - Un OPERATOR con `READ_ONLY` no ve botones de energía; uno con `FULL_ACCESS` sí.
  - El OPERATOR nunca ve Delete.

#### `FRN-19B` (ex `FRN-13B`) - Semáforo de salud global e integración con `GET /api/node/status` (`RF-02`)
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** `FRN-19A`. Empieza con `BAC-29` (datos de prueba) y cierra con `BAC-22`.
- **Entregable:** el Dashboard conecta los medidores a `GET /api/node/status` mediante un servicio en `features/dashboard/services`:
  - uptime legible (días, horas, minutos);
  - semáforo de salud (D3):
    - `Saludable` si CPU, RAM y disco están por debajo del 70 %;
    - `Advertencia` si alguno llega al 70 % o más;
    - `Inaccesible` si el backend responde `502`/`504` o no responde.
  - Aviso de "datos desactualizados" si `stale: true`;
  - consulta automática cada 10 s, que se pausa con la pestaña oculta;
  - skeletons mientras carga y reintento visual ante una desconexión.
- **Criterio de éxito:** los datos reales del nodo se ven en pantalla y se actualizan sin `F5`. Si el backend cae, la UI muestra `Inaccesible` sin romperse.


#### `FRN-16` - Máquina de estados "Operación en progreso" por instancia
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.5 h
- **Depende de:** `FRN-15`. Empieza con `BAC-29` y cierra con `FIX-39`, `FIX-40`, `BAC-24A` y `BAC-24B`.
- **Entregable:**
  - Al confirmar el modal, enviar la orden de energía (con la ruta documentada en `BAC-29`) o `DELETE`, y guardar el `tareaId` del `202` en el estado de la fila (`transitioning`).
  - El botón accionado muestra un spinner, y se deshabilitan todos los botones de esa instancia.
  - Ante un error, desbloquear la fila y mostrar un mensaje distinto según el código (D2), con los significados de `BAC-29`:
    - `409 INSTANCE_INVALID_STATE`: la acción no corresponde al estado actual (por ejemplo, "La instancia ya está encendida").
    - `409 INSTANCE_BUSY`: "La instancia está ejecutando otra tarea. Esperá a que termine".
    - `403 INSTANCE_PROTECTED`: la instancia es de infraestructura y no admite esa acción.
    - `403 INSTANCE_ACCESS_DENIED`: no tenés permiso sobre la instancia.
    - `502 PROXMOX_UNAVAILABLE`: Proxmox no está disponible.
    - `504 PROXMOX_TIMEOUT`: Proxmox no respondió a tiempo, y la acción puede no haberse aplicado.
    - Cualquier otro código: un mensaje genérico.
- **Criterio de éxito:** es imposible disparar una segunda acción sobre la misma instancia mientras hay una orden en curso. Cada uno de los seis códigos muestra su propio mensaje, verificado con una prueba de componente por código.


#### `BAC-25C` (`BRG-04-BAC`) - Reanudación de UPIDs en curso al arrancar el Backend (`RNF-04`)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25A` y `FIX-39` (campo `activeTask`).
- **Entregable:**
  1. Al iniciar el backend, leer de `tareas_asincronas` las filas con `estado = 'RUNNING'` y volver a encolarlas en el pool de `BAC-25A`. El plazo de 3 minutos (D4) se cuenta desde `fecha_creacion`. Si ya venció, se consulta Proxmox una sola vez: si la tarea terminó, se registra su resultado; si no, se marca `FAILED` con `motivo: TIMEOUT`.
  2. Completar `activeTask: { tareaId, action, status } | null` en cada instancia de `GET /api/instances`, cruzando con las tareas `RUNNING`.
- **Criterio de éxito:** si el backend se reinicia durante una tarea de Proxmox, al levantar retoma el sondeo, actualiza `tareas_asincronas` y `auditoria`, y emite `TASK_FINISHED`. Mientras tanto, `GET /api/instances` muestra la tarea en `activeTask`.



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
