
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 02/10/2026, 2ª verificación** (backend `0167b96`, frontend `6fd2c7c`; detalle en [test/informe.md](../test/informe.md)).
> - **Pasaron a [`terminado.md`](terminado.md):**
>   - `FIX-29` y `FIX-40`, completas (1ª verificación);
>   - `SEC-03`, implementada con problema (1ª verificación). Sus dos correcciones, `FIX-41` y `FIX-42`, quedaron resueltas por las PR #72 y #73 del frontend;
>   - `FIX-38`, completa (PR #73, frontend `6fd2c7c`).
> - El backend no tuvo commits nuevos. Las demás tareas siguen sin implementar: se revisó el código y no hay avance parcial. Todas tienen pruebas automatizadas, y hoy fallan.
>
> | Tarea | Área | Estado | Pruebas |
> |---|---|---|---|
> | `FRN-17C` | Frontend | No implementada: no existe `useEvents` | `events-client.test.tsx` (4) |
> | `FIX-36` | Frontend | No implementado: `Auditoria.tsx` no tiene el filtro "Acción" | `audit.test.tsx` (1) |
> | `BAC-18B` | Backend | No implementada: faltan el índice parcial, las particiones y la purga horaria. El índice `(usuario_id, fecha_hora)` ya existe (`idx_auditoria_usuario_fecha`) | `cierre_fase_base…` (2) |
> | `FIX-37` | Backend | No implementado: falta `permisos` en `GET /account/profile` | `resource_access…` (1) |
> | `FIX-39` | Backend | No implementado: campos de `GET /instances`, `shutdown`/`reboot` (404) y auditoría de energía | `puente_etapa1…` (3) |
> | `INF-07B` | Infraestructura | Referencia local verificada (`docker/nginx-edge.conf`); falta aplicarla en el CT 103 | `puente_etapa1…` (1, pasa) |
> | `BAC-29`, `BAC-22`, `BAC-23A`, `BAC-25A` | Backend | No implementadas (`BAC-25A`: sigue una goroutine por tarea y `limiteSeguimiento` fijo en 10 min) | `etapa1_acceptance…` (4) |
> | `FRN-19A`, `FRN-20A`, `FRN-17A` | Frontend | No implementadas: no existen `features/dashboard` ni `features/intances`, e `Instances.tsx` sigue sin consultar la API | `dashboard-metrics` (3), `instances-table` (5), `events-client` (4) |

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

---

### `BAC-18B` - Índice parcial y purga en `sesiones_activas`, y particionamiento trimestral en `auditoria` (PostgreSQL)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `BAC-17B`, `FIX-23`.
- **Estado verificado (01/10/2026, backend `44a2339`):** no implementada. Las pruebas verifican el criterio y no los nombres sugeridos:
  - un índice parcial cualquiera que cubra `jti_access` con `WHERE activa…`;
  - `auditoria` particionada por `RANGE (fecha_hora)`, con tramos trimestrales para el trimestre actual y el siguiente, y una partición `DEFAULT`;
  - un índice compuesto que empiece por `fecha_hora` e incluya `accion` y `resultado`, y otro `(usuario_id, fecha_hora…)`;
  - una purga horaria de sesiones inactivas o vencidas.

  Ya existe el índice `(usuario_id, fecha_hora)` (`idx_auditoria_usuario_fecha`, del modelo GORM); falta todo lo demás.
- **Problema y contexto:**
  1. Las filas históricas o con `activa = false` en `sesiones_activas` penalizan las lecturas en PostgreSQL si no existe un índice parcial ni una purga de registros vencidos.
  2. La tabla `auditoria` es *append-only* (`FIX-23`, no admite `DELETE` bajo ningún concepto) y en la Etapa 1 registrará múltiples eventos por cada operación de Proxmox VE (`PENDING` y `SUCCESS`/`FAILED`). Sin particionamiento por fechas e índices compuestos, las consultas de `GET /api/admin/audit` y la exportación CSV se degradarán progresivamente.
- **Entregable:**
  1. Crear en PostgreSQL (`init.sql` / migración) el índice parcial para `sesiones_activas`:
     ```sql
     CREATE INDEX IF NOT EXISTS idx_sesiones_activas_vigentes
       ON sesiones_activas (jti_token, usuario_id)
       WHERE activa = true;
     ```
  2. Agregar una rutina de limpieza en el backend (ticker cada 1 hora) que ejecute `DELETE FROM sesiones_activas WHERE activa = false OR fecha_expiracion < NOW();`.
  3. Configurar en PostgreSQL el **particionamiento declarativo trimestral por rango de fechas** sobre la tabla `auditoria` (`PARTITION BY RANGE (fecha_hora)`), creando las particiones trimestrales (`auditoria_2026_q3`, `auditoria_2026_q4`, `auditoria_2027_q1` y `auditoria_default`) junto con los índices compuestos `(fecha_hora DESC, accion, resultado)` y `(usuario_id, fecha_hora DESC)`.
- **Criterio de éxito:** Las sesiones muertas se purgan automáticamente de PostgreSQL; la tabla `auditoria` opera sobre particiones trimestrales manteniendo tiempos de consulta constantes (< 20 ms) e inmutabilidad append-only.

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


### `FIX-36` - Regresión: restaurar el filtro "Acción" en Auditoría (`FRN-14` / RF-08) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian (autor de `56b88f5`)
- **Estimación:** 0,5 h
- **Depende de:** `FRN-14` y `BAC-18`.
- **Problema y evidencia:** el commit `56b88f5` (*"eliminar columnas que no serán utilizadas"*) quitó de `frontend/centinela/src/pages/Auditoria.tsx` dos cosas: las columnas *Nodo* e *IP origen*, y también el estado `accion`, el campo `aria-label="Acción"` y el `params.set('accion', …)` de la consulta. Las columnas eran opcionales; el filtro no:
  - RF-08 exige *"filtros por fecha, usuario, tipo de acción y resultado"* (`requerimientos.md:193`).
  - El backend sigue aceptando `accion` en `GET /api/admin/audit` (BAC-18 pasa).
  - Prueba que falla: `test/front/audit.test.tsx`, caso *"actualiza los query parameters de filtro al cambiar los selectores reactivos"* (`Unable to find a label with the text of: /acción/i`). Antes pasaba 7/7; ahora 6/7.
- **Entregable:**
  1. Restaurar en `Auditoria.tsx` el filtro "Acción" (`aria-label="Acción"`), con el estado `accion`, su `updateFilter`, la dependencia del `useEffect` y `params.set('accion', accion)`. La exportación CSV tiene que usar los mismos filtros.
  2. Si el equipo decide que el filtro no va, primero hay que modificar RF-08 y avisar para ajustar la prueba. No se quita en silencio.
- **Criterio de éxito:** `audit.test.tsx` vuelve a pasar 7/7 y al elegir una acción la consulta incluye `accion=<valor>`.

### `FIX-37` - Informar el nivel de acceso por instancia en `GET /account/profile` (`FIX-30` / `FRN-18` / `SEC-04`) (Backend)

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1 h
- **Depende de:** `SEC-04` y `FIX-30`.
- **Problema y evidencia:**
  - `FIX-30` (frontend) arma `canOperateInstance(vmid)` con `perfil.permisos: [{ vmid, nivelAcceso }]` de `GET /account/profile`, como proponía su entregable 2. El backend no envía ese campo: `ports.UsuarioDetalleDTO` (`internal/core/ports/user_port.go:75`) solo tiene `instanciasPermitidas: []int`.
  - Con el backend real, `mapPerfilToUserSession` guarda `permisos: []`, y **`canOperateInstance` da `false` para todo OPERATOR, aunque tenga `FULL_ACCESS`**. Cuando la UI de la Etapa 1 use el helper para habilitar encender, apagar o reiniciar, ningún operador va a poder operar.
  - Prueba que falla: `test/back/resource_access_acceptance_test.go`, caso *"FIX-37 FRN-18 GET /account/profile expone el nivel de acceso por instancia…"*. Las claves que se reciben no incluyen `permisos`.
- **Entregable:**
  1. Agregar a la respuesta de `GET /account/profile` el campo `permisos: [{ vmid, nivelAcceso: "FULL_ACCESS" | "READ_ONLY" }]`, con el mismo formato de `GET /api/admin/users/:id/permissions`, leído de `permisos_instancia`.
  2. Mantener `instanciasPermitidas` para no romper a quienes ya lo usan.
- **Criterio de éxito:**
  - El caso `FIX-37…` pasa: con la 101 en `FULL_ACCESS` y la 102 en `READ_ONLY`, el perfil del OPERATOR informa esos niveles.
  - En el frontend, `canOperateInstance(101)` da `true` y `canOperateInstance(102)` da `false` con datos reales.
  - BAC-14, SEC-04 y LOGIN-04 siguen en verde.

### `FIX-39` - Completar `BAC-21B`: campos de `GET /api/instances`, auditoría de energía, y `shutdown` y `reboot` (Backend)

- **Área:** Backend
- **Asignados:** Tayra y Lisandro
- **Estimación:** 1,5 h
- **Depende de:** `BAC-21B` (en `terminado.md`).
- **Problema y evidencia (verificado el 01/10/2026, backend `44a2339`):** `BAC-21B` cumple la parte de energía con `FULL_ACCESS` y `tareas_asincronas`, pero le faltan tres puntos de su criterio de éxito y de la aclaración del 01/10/2026:
  1. **`GET /api/instances`** sigue con los 5 campos de `BAC-14` (`ports.InstanciaListadaDTO`). Faltan `ip`, `cpuUsage`, `ramUsage`, `maxRam`, `nivelAcceso` y `activeTask`. Prueba que falla: `puente_etapa1…`, caso *"GET /api/instances agrega campos sin romper el contrato de BAC-14"*.
  2. **Auditoría de energía:** ninguna acción de energía se escribe en `auditoria`. `instance_handler.go` responde con `upid` y `tareaId` sin auditar, y `eventos_service.go` solo audita la apertura y el cierre del stream. El criterio dice *"todo se registra en `auditoria` y `tareas_asincronas`"*. Prueba que falla: *"las acciones de energía quedan en auditoria con el upid de la tarea"*.
  3. **`shutdown` y `reboot`:** solo existen `start` y `stop`.
- **Aclaración (01/10/2026, revisión de `etapa1.md`; antes estaba en `BAC-21B`):**
  - **Campos nuevos:** `ip`, `cpuUsage`, `ramUsage`, `maxRam` y `activeTask` se agregan **en `null`**; solo `nivelAcceso` lleva su valor real. Los completan tareas de la Etapa 1: `BAC-23B` (IP), `BAC-22B` (CPU y RAM) y `BAC-25C` (`activeTask`). Así esta corrección no depende de ellas.
  - **`instancesSummary`** no va en cada instancia sino en `GET /api/node/status` (`BAC-22B`); queda fuera de este FIX.
  - **Energía:** se suman `shutdown` y `reboot`, que se agregan a `ProxmoxPort`. La ruta puede ser `/status/:action` o una por acción, como `/start` y `/stop`: lo que se exige es `FULL_ACCESS` y el registro.
  - **Auditoría:** el registro de esta tarea es la orden despachada (`PENDING`). El resultado lo registra `BAC-27`.
  - **Pendiente para la Etapa 1:** la validación de estado previo y los VMIDs protegidos en las acciones nuevas (`BAC-24A`/`BAC-24B`), y la regla de `DELETE` (`BAC-24B`).
- **Entregable:**
  1. Extender `ports.InstanciaListadaDTO` con `ip`, `cpuUsage`, `ramUsage`, `maxRam` y `activeTask` (en `null` por ahora) y con `nivelAcceso`, tomado de `permisos_instancia` (`FULL_ACCESS` o `READ_ONLY`; para el ADMIN, `FULL_ACCESS`). Los 5 campos actuales no cambian.
  2. Auditar cada acción de energía despachada en `auditoria`: `accion`, `instancia_id`, `resultado` (`PENDING`) y el `upid` en `detalles`. Los nombres de las claves de `detalles` son libres.
  3. Agregar `shutdown` y `reboot` a `ProxmoxPort` y exponerlos con `RequireInstanceAccess(..., FULL_ACCESS)`, con el mismo seguimiento en `tareas_asincronas`.
- **Criterio de éxito:** en `test/back/puente_etapa1_acceptance_test.go`:
  - pasan *"GET /api/instances agrega campos…"* (con `nivelAcceso` real), *"las acciones de energía quedan en auditoria con el upid…"* y *"shutdown y reboot exigen FULL_ACCESS…"*;
  - siguen en verde *"las acciones de energía exigen FULL_ACCESS…"*, `BAC-14`, `FIX-14`/`FRN-07` y `SEC-04`.

---
# ETAPA 1
> Las tareas de la Etapa 1 ya verificadas están en [`terminado-1.md`](terminado-1.md), y sus correcciones en [`futuro-1.md`](futuro-1.md).
---

> Tareas de la Ola 1 de [`etapa1.md`](etapa1.md), en curso. Ninguna está implementada todavía (verificado el 02/10/2026).

#### `INF-07B` (`BRG-03`) - Configuración de Nginx para SSE en el servidor (`RNF-06`)

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.0 h
- **Depende de:** `FIX-31` y `BAC-21C`, las dos terminadas. No depende de `BAC-26`: el stream ya emite `TASK_FINISHED` para `start` y `stop`.
- **Problema:** sin configuración específica, Nginx corta las conexiones de `/api/events` a los 60 s o retiene en buffer los mensajes SSE.
- **Referencia:** `docker/nginx-edge.conf` en el repositorio integrador (`FIX-31`). Su `location /api/` ya tiene `proxy_http_version 1.1`, `proxy_buffering off` y `proxy_read_timeout 1h`, y la suite de pruebas lo usa delante del backend. La configuración del servidor (CT 103) no está en los submódulos.
- **Entregable:**
  1. Llevar al Nginx del CT 103 la configuración de `docker/nginx-edge.conf`, con los certificados reales en `/etc/nginx/tls/`. Es el pendiente para infraestructura de `FIX-31`.
  2. Confirmar que `/api/events` queda sin buffer ni corte por inactividad: `proxy_buffering off`, `proxy_cache off` y `proxy_read_timeout` ≥ 1 h. Si en el servidor hace falta algo más, como `gzip off` o `proxy_set_header Connection ''`, agregarlo también en `docker/nginx-edge.conf`, para que la referencia y el servidor no se separen.
- **Verificación local (01/10/2026):** `test/back/puente_etapa1_acceptance_test.go`, caso *"INF-07B BRG-03 /api/events a través del borde Nginx…"*, **pasa** con la referencia. Verifica:
  - en `docker/nginx-edge.conf`: `proxy_http_version 1.1`, `proxy_buffering off` y `proxy_read_timeout` ≥ 5 min;
  - por el servicio `edge` en HTTPS: ticket, stream y `start`, con el `TASK_FINISHED` llegando en ~1 s.

  Queda pendiente aplicarlo en el CT 103, que se valida en `INT-03`.
- **Criterio de éxito:** a través de Nginx (HTTPS), un cliente conectado a `/api/events?ticket=…` recibe el `TASK_FINISHED` de un `start` sobre una instancia de prueba sin demora visible, y la conexión sigue abierta después de 5 minutos sin tráfico.


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

#### `FRN-20A` (ex `FRN-14A`) - Tabla interactiva de inventario con badges de estado e IP (`RF-03`)
- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.5 h
- **Depende de:** empieza ya contra `GET /api/instances` (`BAC-14`: `id`, `name`, `type`, `node` y `status`) y `BAC-29`. Cierra con `FIX-39` (campos nuevos).
- **Entregable:** tabla en `pages/Instances.tsx` con su servicio y su hook en `features/intances`. Reemplaza la maqueta y los archivos vacíos. Columnas:
  - ID, Nombre, Tipo (`VM` / `LXC`);
  - Estado (`Running` en verde, `Stopped` en gris);
  - IP con botón para copiar al portapapeles;
  - botonera de acciones (solo la estructura; los botones los conecta `FRN-15`).
- **Criterio de éxito:** VMs y LXC se ven en la misma tabla. Si la IP es `null`, se muestra `No detectada`. La tabla tolera que los campos nuevos vengan en `null`.


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
