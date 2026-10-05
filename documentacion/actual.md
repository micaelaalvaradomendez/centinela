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
