# Etapa 1: Core Operativo de Proxmox y MVP

> **Estado:** Planificada. Revisión del 01/10/2026 contra el código (backend `44a2339`, frontend `3d1e84a`).
> **Dependencia previa:** Cierre de etapa base (`actual.md`, `futuro.md`, `terminado.md` con hito `LOGIN-04` aprobado e infraestructura base operativa). Algunas tareas de `actual.md` frenan tareas puntuales de esta etapa; se indican en cada una y en la §3.
> **Requerimientos cubiertos:** `RF-02` (Dashboard Host), `RF-03` (Inventario completo), `RF-04` (Ciclo de Vida con UPID), `RF-11` parcial (notificación en tiempo real de fin de tarea) y `RNF-04` parcial (reanudación de tareas).

---

> [!IMPORTANT]
> **Revisión del 01/10/2026: qué cambió en este documento**
> 1. **Se revisó el código.** Parte de lo que pedía la etapa ya existe (§2). Esas tareas pasan a ser de **ajuste o verificación** y se reestimaron: `BAC-23B`, `BAC-24A`, `BAC-24B`, `BAC-25A`, `BAC-25B` y `BAC-26`.
> 2. **Contrato de eventos:** vale el de `BAC-21` (terminado), definido en `backend/docs/contrato-eventos.md` y `frontend/centinela/src/types/notifications.ts`:
>    - los campos van en español: `tipo`, `severidad`, `recursoTipo`, `recursoId`, `mensaje`, `fechaHora` y `detalles`;
>    - la correlación con la acción se hace por `detalles.tareaId`, que también devuelve el `202`;
>    - el UPID crudo no viaja al frontend.
>
>    Se corrigió el ejemplo en inglés (`type`, `payload.upid`) que tenía `BAC-26`.
> 3. **`BAC-21B` está en `terminado.md` como implementada con problema, y lo que le falta es `FIX-39` (`futuro.md`).** `BAC-21B` cumple `start` y `stop` con `FULL_ACCESS` y su registro en `tareas_asincronas`. `FIX-39` agrega:
>    - en `GET /api/instances`, los campos `ip`, `cpuUsage`, `ramUsage`, `maxRam` y `activeTask` en `null`, y `nivelAcceso` con su valor real;
>    - `shutdown` y `reboot` con `FULL_ACCESS`;
>    - la auditoría de la orden despachada (`PENDING`) con su `upid`.
>
>    Los campos en `null` los completan `BAC-23B` (IP), `BAC-22B` (CPU y RAM) y `BAC-25C` (`activeTask`), así que `FIX-39` no depende de la Etapa 1. `instancesSummary` va en `GET /api/node/status` (`BAC-22B`), no en cada instancia.
>    - **La ruta no es obligatoria:** `/status/:action` o una ruta por acción (`/start`, `/stop`, `/shutdown`, `/reboot`). Lo que se exige es `FULL_ACCESS` y el registro. La forma que se elija la documenta `BAC-29`.
>    - **El endpoint `DELETE` no existe** y queda entero en `BAC-24B`: solo `ADMIN`, `409` si la instancia está encendida, VMIDs protegidos, seguimiento y auditoría.
> 4. **Contrato primero (`BAC-29`, nueva).** El backend publica los payloads de la etapa antes de implementarlos. El frontend maqueta e integra contra ese contrato con datos de prueba, y cierra cada tarea cuando la tarea de backend correspondiente está terminada.
> 5. **Sentido de las dependencias.** Ninguna tarea de backend depende de infraestructura: todo se desarrolla contra el simulador y Redis local. Infraestructura depende de lo que entrega el backend (endpoints, variables de entorno). El frontend depende del contrato (`BAC-29`) para empezar y de la implementación del backend para cerrar.
> 6. **Tareas nuevas** (lo que falta en el código y no cubría ninguna tarea):
>    - `BAC-29`: contrato de la etapa.
>    - `BAC-24C`: borrar los permisos de una instancia eliminada. Como `/cluster/nextid` reutiliza el menor VMID libre, una VM nueva heredaría los permisos de la eliminada.
>    - `FRN-20C`: se separó de `FRN-19C`.
> 7. **Renumeración del frontend:** `FRN-13A/B` → `FRN-19A/B` y `FRN-14A/B` → `FRN-20A/B`, porque `FRN-13` y `FRN-14` ya existen y están terminadas. `FRN-19C` se divide en `FRN-19C` (Dashboard) y `FRN-20C` (Inventario).
> 8. **Dependencias de la fase base que faltaban:**
>    - `FIX-37`, porque sin ella `canOperateInstance` da `false` para todo OPERATOR;
>    - `SEC-03`, para `PermissionGate`;
>    - `FIX-38`;
>    - `FRN-17C`, que es la conexión del mismo hook;
>    - `BAC-18B`, porque cambia el esquema de `auditoria`;
>    - `FIX-39`, por los campos, `shutdown`/`reboot` y la auditoría de lo despachado.
>
>    `FIX-31` ya está completa: el Nginx del borde está versionado en este repositorio, en `docker/nginx-edge.conf`, y lo usa la suite de pruebas.
> 9. **Se mantienen del análisis del 26/09/2026:**
>    - Las tablas reales son `auditoria`, `tareas_asincronas` y `permisos_instancia`.
>    - Cada código de error nuevo se agrega al inventario de `FIX-08` y a Swagger (`RNF-05`) en la tarea que lo introduce.
>    - El enlace de la notificación al detalle de la instancia (`RF-11`) queda pospuesto, porque en esta etapa no existe esa vista.

## 1. Criterio de Planificación y Decisiones de Arquitectura

1. **Regla de Granularidad ($\le$ 2.5 h):** Ninguna tarea supera las 2.5 horas de estimación. Cada ítem tiene una responsabilidad técnica única (separando capa de datos, lógica de negocio y capa visual).
2. **Premisa de Nodo Fijo:** Esta etapa opera asumiendo un **único nodo Proxmox VE fijo** configurado por variables de entorno (`PROXMOX_NODE`). El soporte multi-nodo (`RNF-07`) queda fuera del MVP.
3. **Exclusiones Explícitas:**
   - `RF-12` (organizaciones / multi-tenant) permanece descartado. Ninguna tabla o payload debe incluir `organization_id`.
   - La consola web remota (SSH, noVNC, xterm.js) queda excluida por motivos de seguridad operacional.
   - Métricas históricas (`RF-05`), aprovisionamiento (`RF-07`), snapshots (`RF-06`) y alertas por correo de `RF-11` corresponden a las Etapas 2 y 3.
4. **Caché Compartida (Redis):** `GET /api/node/status` consulta `/nodes/{node}/status` en Proxmox VE y guarda el resultado en Redis con TTL de 5 a 10 s. Además, guarda el **último estado conocido** para responder si Proxmox no está disponible (`RF-02`, Infraestructura).
5. **Reutilización de Cimientos Base:**
   - **Auditoría (`BAC-18`):** toda acción de ciclo de vida y borrado se registra en la tabla append-only **`auditoria`** (`FIX-23`), con `upid`, `action` y `resource_type` en `detalles` (JSONB). El estado de ejecución va en **`tareas_asincronas`**. No se crean tablas paralelas.
   - **Contrato de Eventos (`BAC-21`):** el fin de tarea se notifica con `TASK_FINISHED` respetando el schema de `backend/docs/contrato-eventos.md`.
   - **Canal en tiempo real (`BAC-21C`):** `/api/events` (SSE con ticket efímero, bus Redis Pub/Sub y filtro por permiso) ya existe. Esta etapa lo usa; no crea un canal nuevo.
6. **Poller de UPID Desacoplado:** El backend responde `HTTP 202 Accepted` de inmediato con `{ upid, tareaId }` y delega el seguimiento a un worker pool concurrente en segundo plano.
7. **Diseño Antierror en Frontend:** Patrón "Operación en progreso": al confirmar una acción, los controles de esa instancia se bloquean con spinner hasta recibir `TASK_FINISHED` con el mismo `tareaId`.
8. **VMIDs protegidos:** las acciones destructivas (`stop`, `shutdown`, `reboot`, `DELETE`) sobre los VMIDs de `PROXMOX_PROTECTED_VMIDS` responden `403 INSTANCE_PROTECTED` (`RejectProtectedInstance`, ya existe para `stop`).

### Decisiones a confirmar por el equipo antes de cargar

| # | Decisión | Propuesta | Afecta |
|---|---|---|---|
| D1 | ¿Quién puede ver `GET /api/node/status`? | Cualquier usuario autenticado (`RF-02` no lo restringe). El `instancesSummary` cuenta todas las instancias del nodo. | `BAC-22`, `BAC-22B`, `FRN-19B` |
| D2 | Código para una acción incompatible con el estado actual (`start` sobre una instancia encendida, `DELETE` sobre una encendida) | `409 INSTANCE_INVALID_STATE` (distinto de `409 INSTANCE_BUSY`, que es "hay otra tarea en curso") | `BAC-24A`, `BAC-24B`, `FRN-16` |
| D3 | Umbrales del semáforo de salud | `Saludable` < 80 % en CPU, RAM y disco; `Advertencia` ≥ 80 % en alguno; `Inaccesible` si el backend responde `502`/`504` o no responde | `FRN-19B` |
| D4 | Timeout del seguimiento de UPID | `UPID_TIMEOUT` configurable, por defecto 3 min (hoy está fijo en 10 min) | `BAC-25B` |

---

## 2. Estado del código al 01/10/2026 (punto de partida)

| Pieza | Estado | Dónde | Tarea que lo completa |
|---|---|---|---|
| Simulador de Proxmox: estado del nodo, inventario, acciones `start`/`stop`/`shutdown`/`reboot`, estado de tareas, `DELETE`, IP de VM y LXC | ✅ Completo | `backend/cmd/proxmox-simulador` | `BAC-28`, `FIX-32`, `FIX-35` (terminadas) |
| Redis local y adaptador (`KeyValueStore`, Pub/Sub) | ✅ | `internal/adapters/secondary/redis` | `INF-06A`, `BAC-17A` (terminadas) |
| `GET /api/instances` con `{ id, name, type, node, status }` y filtro RBAC | ✅ | `instance_handler.go` | `BAC-14` (terminada). Los campos nuevos los agrega `FIX-39` |
| `POST /instances/:vmid/start` y `/stop` con `FULL_ACCESS`, `202 { upid, tareaId }`, `409 INSTANCE_BUSY` y `403 INSTANCE_PROTECTED` (solo `stop`) | ✅ | `main.go`, `instance_handler.go` | `FIX-16`, `SEC-04`, `FIX-33`, `BAC-21B` (con problema) |
| `shutdown` y `reboot` | ❌ | `ports/proxmox_port.go`, `proxmox/client.go`, `main.go` | `FIX-39` |
| Endpoint `DELETE /api/instances/:vmid` | ❌ No existe | | `BAC-24B` |
| Seguimiento de UPID: una goroutine por tarea, consulta cada 1 s, corta a los 10 min, guarda en `tareas_asincronas` y publica `TASK_FINISHED` | 🟡 Funciona, pero sin pool acotado ni timeout configurable, y solo con los textos de `start`/`stop` | `services/seguimiento_tareas.go` | `BAC-25A`, `BAC-25B` |
| `/api/events` (SSE): ticket efímero, Redis Pub/Sub, filtro por permiso y corte en vivo | ✅ | `services/eventos_service.go` | `BAC-21C` (terminada) |
| Auditoría de las acciones de energía | ❌ No se registra nada | `instance_handler.go` | `FIX-39` (orden despachada) y `BAC-27` (resultado) |
| `GET /api/node/status` | ❌ No existe | | `BAC-22` |
| IP, CPU y RAM por instancia | ❌ | | `BAC-23A`, `BAC-23B`, `BAC-22B` |
| Reanudar tareas al reiniciar el backend | ❌ | | `BAC-25C` |
| Borrar permisos de una instancia eliminada | ❌ | | `BAC-24C` (nueva) |
| Dashboard e Inventario en el frontend | 🟡 Maquetas estáticas, sin llamadas a la API | `pages/Dashboard.tsx`, `pages/Instances.tsx` | `FRN-19A/B/C`, `FRN-20A/B/C` |
| Carpetas `features/dashboard` y `features/intances`, `hooks/useWebSocket.js` | Archivos vacíos | | `FRN-19A`, `FRN-20A` y `FRN-17A` los reemplazan o eliminan |
| `usePermissions` con `canAccessInstance` y `canOperateInstance` | 🟡 Sin datos reales para el OPERATOR | `hooks/usePermissions.ts` | `FIX-37`, `SEC-03`, `FIX-38` (fase base) |
| Hook `useEvents` | ❌ | | `FRN-17C` (fase base) y `FRN-17A` |
| Contrato `RealtimeEvent` en TypeScript | ✅ | `types/notifications.ts` | `BAC-21` (terminada) |

---

## 3. Dependencias con la fase base (`actual.md` / `terminado.md`)

| Tarea de la fase base | Estado al 01/10 | Frena en la Etapa 1 | Tipo |
|---|---|---|---|
| `FIX-39` (completa `BAC-21B`) | ❌ (`BAC-21B` está en `terminado.md` con problema) | Cierre de `BAC-22B`, `BAC-23B`, `BAC-24A`, `BAC-24B`, `BAC-25C`, `BAC-27`, `FRN-16` y `FRN-20A` | Bloqueante |
| `FIX-37` | ❌ | Cierre de `FRN-15`, `FRN-16` e `INT-02`, porque sin ella ningún OPERATOR puede operar | Bloqueante |
| `SEC-03` | 🟡 | Cierre de `FRN-15` y `FRN-16` (`PermissionGate` oculta Delete al OPERATOR) | Bloqueante |
| `FIX-38` | ❌ | Cierre de `FRN-15` (`canAccessInstance` no debe aceptar `READ_ONLY` como rol) | Bloqueante |
| `FRN-17C` | ❌ | `FRN-17A` y `FRN-16B`: es la conexión del mismo hook. Conviene hacerlas juntas | Bloqueante |
| `BAC-18B` | ❌ | `BAC-27`: el particionamiento cambia el esquema de `auditoria`, así que tiene que estar mergeado antes de sumar los registros de la etapa | Bloqueante blanda |
| `INF-06B`, `INF-08A` | En `terminado.md` sin verificar | `INT-03` | Bloqueante |
| `FIX-29`, `FIX-36` | ❌ | Ninguna (solo el cierre de la fase base) | — |

---

## 4. Mapa de Dependencias de la Etapa 1

Las líneas continuas son dependencias bloqueantes. Las punteadas significan que la tarea **puede empezar** contra el contrato con datos de prueba, pero **cierra** con la otra. Los nodos grises son tareas de la fase base.

```mermaid
flowchart TD
    subgraph Base [Fase base pendiente]
        FIX39[FIX-39: Campos, shutdown/reboot, auditoría]:::base
        FIX37[FIX-37: permisos en /account/profile]:::base
        SEC03[SEC-03: PermissionGate]:::base
        FIX38[FIX-38: READ_ONLY no es rol]:::base
        FRN17C[FRN-17C: conexión useEvents + ticket]:::base
        BAC18B[BAC-18B: particiones auditoria]:::base
        INF06B[INF-06B: Redis servidor]:::base
        INF08A[INF-08A: SMTP servidor]:::base
    end

    subgraph Infra [Infraestructura]
        INF07[INF-07: Token, red y VMIDs protegidos]
        INF07B[INF-07B: Nginx para /api/events]
    end

    subgraph Back [Backend]
        BAC29[BAC-29: Contrato de la etapa]
        BAC22[BAC-22: Telemetría + Redis]
        BAC22B[BAC-22B: Resumen por estado + CPU/RAM]
        BAC23A[BAC-23A: Adaptador inventario + IP]
        BAC23B[BAC-23B: IP en GET /api/instances]
        BAC24A[BAC-24A: Validación de estado y protegidos]
        BAC24B[BAC-24B: DELETE seguido y protegido]
        BAC24C[BAC-24C: Limpieza de permisos al eliminar]
        BAC25A[BAC-25A: Pool acotado de UPID]
        BAC25B[BAC-25B: Timeout, reintentos, exitstatus]
        BAC25C[BAC-25C: Reanudación + activeTask]
        BAC26[BAC-26: TASK_FINISHED de todas las acciones]
        BAC27[BAC-27: Auditoría del resultado]
    end

    subgraph Front [Frontend]
        FRN19A[FRN-19A: Gauges CPU/RAM/Disco]
        FRN19B[FRN-19B: Semáforo + integración nodo]
        FRN19C[FRN-19C: Tarjetas VMs/LXC]
        FRN20A[FRN-20A: Tabla de inventario]
        FRN20B[FRN-20B: Filtros y buscador]
        FRN20C[FRN-20C: CPU/RAM + estado en vivo]
        FRN15[FRN-15: Modales antierror]
        FRN16[FRN-16: Operación en progreso]
        FRN16B[FRN-16B: Resincronización tras F5]
        FRN17A[FRN-17A: Consumo de eventos]
        FRN17B[FRN-17B: Desbloqueo + Toasts]
    end

    subgraph Hito [Verificación]
        INT01[INT-01: Integración backend]
        INT02[INT-02: Smoke test local]
        INT03[INT-03: Validación en servidor]
    end

    BAC22 --> BAC22B
    BAC23A --> BAC22B
    BAC23A --> BAC23B
    FIX39 --> BAC23B
    FIX39 --> BAC24A
    FIX39 --> BAC24B
    BAC25A --> BAC25B
    BAC25B --> BAC26
    BAC25B --> BAC27
    FIX39 --> BAC27
    BAC18B -.-> BAC27
    BAC24B --> BAC24C
    BAC25B --> BAC24C
    BAC25A --> BAC25C
    FIX39 --> BAC25C

    BAC29 -.-> FRN19B
    BAC29 -.-> FRN20A
    BAC29 -.-> FRN16
    FRN19A --> FRN19B
    BAC22 --> FRN19B
    FRN19B --> FRN19C
    BAC22B --> FRN19C
    FRN17A --> FRN19C
    FIX39 --> FRN20A
    FRN20A --> FRN20B
    FRN20A --> FRN20C
    BAC22B --> FRN20C
    FRN17A --> FRN20C
    FRN20A --> FRN15
    SEC03 --> FRN15
    FIX37 --> FRN15
    FIX38 --> FRN15
    FRN15 --> FRN16
    BAC24A --> FRN16
    BAC24B --> FRN16
    FRN17C --> FRN17A
    FRN16 --> FRN17B
    FRN17A --> FRN17B
    BAC26 --> FRN17B
    BAC25C --> FRN16B
    FRN16 --> FRN16B
    FRN17C --> FRN16B


    BAC23B --> INT01
    BAC24C --> INT01
    BAC25C --> INT01
    BAC26 --> INT01
    BAC27 --> INT01
    INT01 --> INT02
    FRN17B --> INT02
    FRN16B --> INT02
    FRN19C --> INT02
    FRN20B --> INT02
    FRN20C --> INT02
    INT02 --> INT03
    INF07 --> INT03
    INF07B --> INT03
    INF06B --> INT03
    INF08A --> INT03

    classDef base fill:#e5e7eb,stroke:#6b7280,color:#111827
```

---

## 5. Orden de carga en ClickUp

| Ola | Se carga y se empieza | Condición para empezar |
|---|---|---|
| **1** | `INF-07`, `INF-07B`, `BAC-29`, `BAC-22`, `BAC-23A`, `BAC-25A`, `FRN-19A`, `FRN-20A` (contra `BAC-14`/`BAC-29`), `FRN-17A` (junto con `FRN-17C`) | Ninguna pendiente. `FRN-20A` cierra cuando termina `BAC-21B` |
| **2** | `BAC-22B`, `BAC-23B`, `BAC-24A`, `BAC-24B`, `BAC-25B`, `FRN-19B`, `FRN-20B`, `FRN-15` (los modales se maquetan antes) | Tareas de la Ola 1 y `FIX-39`. `FRN-15` cierra con `SEC-03`, `FIX-37` y `FIX-38` |
| **3** | `BAC-24C`, `BAC-25C`, `BAC-26`, `BAC-27`, `FRN-16`, `FRN-19C`, `FRN-20C` | Ola 2. `BAC-27` necesita además `BAC-18B` mergeada |
| **4** | `FRN-16B`, `FRN-17B`, `INT-01` | Ola 3 |
| **5** | `INT-02`, y después `INT-03` | Todo lo anterior. `INT-03` necesita además `INF-06B` e `INF-08A` |

**Camino crítico:** `FIX-39` → `BAC-24A/B` → `FRN-16` → `FRN-17B` → `INT-02` → `INT-03`. En paralelo, `SEC-03`, `FIX-37` y `FIX-38` tienen que estar listas antes de que cierre `FRN-15`.

---

## 6. Desglose Detallado de Tareas (ClickUp)

### Bloque 1: Infraestructura y Entorno

> [!IMPORTANT]
> **Segmentación local / servidor.** **El backend y el frontend desarrollan y prueban en local, sin esperar al servidor**, e infraestructura configura el servidor en tareas aparte. Las dos partes se juntan en la validación final (`INT-03`). Infraestructura solo depende de cosas que el backend ya entregó.
>
> | Lo que necesita el desarrollo | En local (no bloquea) | En el servidor (infraestructura) |
> |---|---|---|
> | Proxmox VE | Simulador `cmd/proxmox-simulador` (completo: `BAC-28`, `FIX-32`, `FIX-35`) | `INF-07` |
> | Redis | `INF-06A` y `BAC-17A` (terminadas) | `INF-06B` (fase base) |
> | SMTP | `INF-08B` (terminada) | `INF-08A` (fase base) |
> | Nginx para SSE | No hace falta: en local el frontend habla directo con el backend | `INF-07B` |

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

---

### Bloque 2: Backend - Contrato, Telemetría e Inventario

#### `BAC-29` - Contrato HTTP y de eventos de la Etapa 1 (nueva)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.0 h
- **Depende de:** ninguna. Se hace primero, para que el frontend no espere a la implementación.
- **Entregable:** `backend/docs/contrato-etapa1.md` y las anotaciones Swagger (sin implementación) de:
  1. `GET /api/node/status`: `{ cpu: { usagePercent, cores }, ram: { usedGb, totalGb, usagePercent }, storage: { usedGb, totalGb, usagePercent }, uptimeSeconds, instancesSummary: { vms: { running, stopped, paused, total }, lxc: { … } }, stale, fetchedAt }`.
  2. `GET /api/instances`: lo de `BAC-14`, más `ip | null`, `cpuUsage | null` (0-100), `ramUsage | null`, `maxRam | null` (bytes), `nivelAcceso` y `activeTask: { tareaId, action, status } | null`.
  3. Las rutas de energía tal como queden con `FIX-39`. Pueden ser `/status/:action` o una por acción, como las actuales `/start` y `/stop`. Para `start`, `shutdown`, `stop` y `reboot`, y para el `DELETE /api/instances/:vmid` de `BAC-24B`: `202 { upid, tareaId }`.
  4. Los `detalles` de `TASK_FINISHED`: `{ tareaId, accion, estado: COMPLETED|FAILED, exitstatus, error? }`.
  5. La tabla de códigos de error de la etapa con su estado HTTP:
     - `INSTANCE_ACCESS_DENIED`, `INSTANCE_PROTECTED`, `INSTANCE_BUSY`, `INSTANCE_INVALID_STATE` (D2) e `INSTANCE_NOT_FOUND`;
     - `PROXMOX_UNAVAILABLE` (`502`/`504`) e `INVALID_ACTION`.

     Se agregan también al inventario de `FIX-08`.
- **Criterio de éxito:** el frontend puede maquetar `FRN-19B`, `FRN-20A`, `FRN-16` y `FRN-17B` solo con este documento, y Swagger muestra los endpoints con ejemplos. Si una tarea posterior cambia el contrato, actualiza este archivo y avisa al frontend.

#### `BAC-22` - Adaptador de telemetría del nodo con caché en Redis (`RF-02`)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.5 h
- **Depende de:** `INF-06A` y `BAC-17A` (terminadas). Se desarrolla contra el simulador, que ya responde `/nodes/{node}/status`.
- **Entregable:** endpoint `GET /api/node/status` (`RequireAuth`, según D1) que consume `/nodes/{node}/status` de Proxmox:
  - Agrega `ObtenerEstadoNodo` a `ProxmoxPort` y al cliente.
  - Normaliza el porcentaje y los núcleos de CPU, convierte de bytes a GB la RAM y el almacenamiento, y devuelve el uptime en segundos, con el formato de `BAC-29`.
  - Guarda el resultado en Redis con TTL de 5 a 10 s, y además el último estado conocido sin TTL.
  - Si Proxmox no responde y hay un último estado conocido, responde `200` con `stale: true`. Si no hay ninguno, responde `502`/`504 PROXMOX_UNAVAILABLE`.
- **Criterio de éxito:**
  - Responde en menos de 50 ms con la caché vigente.
  - Al vencer el TTL consulta Proxmox, actualiza Redis y responde `200`.
  - Con el simulador detenido, responde el último estado con `stale: true`.

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

---

### Bloque 3: Backend - Ciclo de Vida Asíncrono, UPID y Notificaciones

> [!NOTE]
> **Punto de partida en el backend:**
> - **`BAC-21B`** (en `terminado.md`, con problema): `start` y `stop` exigen `FULL_ACCESS` y quedan en `tareas_asincronas`.
> - **`FIX-39`** (en `futuro.md`) agrega:
>   - `shutdown` y `reboot` con `FULL_ACCESS`;
>   - la auditoría de la orden despachada (`PENDING`, con el `upid`);
>   - los campos nuevos de `GET /api/instances`.
>
> `BAC-24A` completa la energía y **`BAC-24B` crea el `DELETE` completo**, que no hace ninguna de las dos.

#### `BAC-24A` - Validación de estado previo y VMIDs protegidos en acciones de energía (`RF-04`)
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

#### `BAC-24C` - Limpieza de permisos al eliminar una instancia (nueva)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.0 h
- **Depende de:** `BAC-24B` y `BAC-25B`.
- **Problema:** `/cluster/nextid` devuelve el menor VMID libre. Si se borra la 110 y después se crea otra instancia, recibe el 110 y **hereda las filas de `permisos_instancia`** de la borrada: un operador tendría acceso a una máquina que no se le asignó.
- **Entregable:**
  - Cuando una tarea de borrado termina con `COMPLETED`, eliminar en una transacción las filas de `permisos_instancia` de ese VMID.
  - Registrar la limpieza en `auditoria` con `accion: INSTANCE_PERMISSIONS_PURGED` y la lista de usuarios afectados.
  - Publicar el `TASK_FINISHED` **antes** de borrar los permisos, para que los operadores asignados lo reciban.
  - Si la tarea falla, no se toca nada.
- **Criterio de éxito:** después de borrar la 110, `GET /api/admin/users/:id/permissions` del operador ya no la incluye. Al crear una instancia nueva con el VMID 110, el operador no la ve.

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

#### `BAC-25B` - Timeout configurable, reintentos y `exitstatus` en el evento
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25A`.
- **Entregable:**
  1. Timeout configurable `UPID_TIMEOUT` (D4, por defecto 3 min) en lugar de los 10 min fijos. Al vencer: `FAILED` con un error controlado.
  2. Reintentos con retroceso cuando falla la consulta del estado, en lugar del reintento fijo cada 1 s.
  3. Completar los `detalles` del `TASK_FINISHED` según `BAC-29`:
     - agregar `exitstatus`;
     - agregar los textos de `shutdown`, `reboot` y `delete` en `nombresAccion`.
  4. Publicar el resultado por el bus existente (`eventosService.Publicar`) y exponerlo para `BAC-27` y `BAC-24C`, por ejemplo con un callback o un suscriptor interno.
- **Criterio de éxito:**
  - Una tarea que no termina se marca `FAILED` al vencer `UPID_TIMEOUT`.
  - Todo `TASK_FINISHED` trae `tareaId`, `accion`, `estado` y `exitstatus`.
  - `BAC-27` y `BAC-24C` reciben todos los resultados.

#### `BAC-25C` (`BRG-04-BAC`) - Reanudación de UPIDs en curso al arrancar el Backend (`RNF-04`)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25A` y `FIX-39` (campo `activeTask`).
- **Entregable:**
  1. Al iniciar el backend, leer de `tareas_asincronas` las filas con `estado = 'RUNNING'` y volver a encolarlas en el pool de `BAC-25A`.
  2. Completar `activeTask: { tareaId, action, status } | null` en cada instancia de `GET /api/instances`, cruzando con las tareas `RUNNING`.
- **Criterio de éxito:** si el backend se reinicia durante una tarea de Proxmox, al levantar retoma el sondeo, actualiza `tareas_asincronas` y `auditoria`, y emite `TASK_FINISHED`. Mientras tanto, `GET /api/instances` muestra la tarea en `activeTask`.

#### `BAC-26` - Verificación de `TASK_FINISHED` por `/api/events` para todas las acciones (`RF-11` parcial)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.0 h
- **Depende de:** `BAC-25B` y `BAC-21C` (terminada: canal SSE, Pub/Sub y filtro por permiso).
- **Contexto:** el canal ya existe; esta tarea **no crea uno nuevo**. Verifica y documenta la emisión para las cinco acciones.
- **Entregable:**
  1. Pruebas en `eventos_service_test.go` (o de integración) para `start`, `shutdown`, `stop`, `reboot` y `delete`, con el payload de `BAC-29`. Ejemplo:
     ```json
     {
       "id": "0192…",
       "tipo": "TASK_FINISHED",
       "severidad": "INFO",
       "recursoTipo": "VM",
       "recursoId": "101",
       "mensaje": "La tarea de encendido finalizó correctamente",
       "fechaHora": "2026-10-01T13:30:00Z",
       "detalles": { "tareaId": "0192…", "accion": "start", "estado": "COMPLETED", "exitstatus": "OK" }
     }
     ```
  2. Actualizar `docs/contrato-eventos.md` y `docs/eventos-tiempo-real.md`.
- **Criterio de éxito:** los clientes conectados reciben el evento al terminar cada tarea. Un OPERATOR sin permiso sobre la instancia no lo recibe; uno con `READ_ONLY` sí.

#### `BAC-27` - Auditoría del resultado de las acciones de ciclo de vida en `auditoria` (`RF-08`)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:**
  - `BAC-18` y `FIX-23` (terminadas);
  - `FIX-39` (registra la orden despachada; la del `DELETE`, `BAC-24B`);
  - `BAC-25B`;
  - `BAC-18B` (particionamiento de `auditoria`, que tiene que estar mergeado antes).
- **Entregable:**
  - Para cada tarea que termina, registrar en `auditoria` una segunda entrada con el mismo `upid` y `tareaId` que la de `FIX-39`, más `status: "SUCCESS"` o `"FAILED"` y `exitstatus`.
  - La entrada de la orden despachada (`PENDING`) la hacen `FIX-39` y `BAC-24B`. Acá solo se verifica que tenga `action`, `user_id`, `resource_id` y `upid`.
- **Criterio de éxito:**
  - Cada acción deja exactamente dos registros correlacionados por `upid`, sin migraciones de esquema nuevas.
  - Las tareas vencidas por timeout quedan como `FAILED`.
  - `GET /api/admin/audit` permite filtrar por `accion`.

---

### Bloque 4: Frontend - Vistas del Host e Inventario

> [!NOTE]
> Cada tarea de frontend indica **con qué puede empezar** (contrato `BAC-29` y datos de prueba) y **con qué cierra** (la implementación de backend). Las pantallas actuales (`Dashboard.tsx`, `Instances.tsx`) son maquetas estáticas que se reemplazan. Las carpetas vacías `features/dashboard` y `features/intances` se usan o se eliminan.

#### `FRN-19A` (ex `FRN-13A`) - Maquetado y medidores de recursos del Host (CPU / RAM / Almacenamiento)
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** ninguna (datos de prueba). Puede partir de la maqueta de `pages/Dashboard.tsx`.
- **Entregable:** componentes reutilizables en `features/dashboard`, con Tailwind y shadcn:
  - barras o gauges de CPU (% y núcleos);
  - uso de RAM (GB usados sobre el total);
  - uso de almacenamiento (GB o TB usados sobre el total).
- **Criterio de éxito:** los componentes son responsive y renderizan valores de 0 % a 100 %, con cambio de color según la saturación (D3). Tienen pruebas de componente.

#### `FRN-19B` (ex `FRN-13B`) - Semáforo de salud global e integración con `GET /api/node/status` (`RF-02`)
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** `FRN-19A`. Empieza con `BAC-29` (datos de prueba) y cierra con `BAC-22`.
- **Entregable:** el Dashboard conecta los medidores a `GET /api/node/status` mediante un servicio en `features/dashboard/services`:
  - uptime legible (días, horas, minutos);
  - semáforo de salud según D3 (`Saludable`, `Advertencia`, `Inaccesible`), y aviso de "datos desactualizados" si `stale: true`;
  - consulta automática cada 10 s, que se pausa con la pestaña oculta;
  - skeletons mientras carga y reintento visual ante una desconexión.
- **Criterio de éxito:** los datos reales del nodo se ven en pantalla y se actualizan sin `F5`. Si el backend cae, la UI muestra `Inaccesible` sin romperse.

#### `FRN-19C` (`BRG-05-FRN`, parte Dashboard) - Tarjetas de conteo de VMs/LXC con auto-actualización (`RF-02`)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 1.0 h
- **Depende de:** `BAC-22B`, `FRN-19B` y `FRN-17A`.
- **Entregable:** tarjetas de resumen de VMs y LXC (`En ejecución`, `Detenidas`, `Total`) con los datos de `instancesSummary`. Usan la misma consulta cada 10 s de `FRN-19B` y vuelven a consultar al recibir cualquier `TASK_FINISHED`.
- **Criterio de éxito:** después de encender una VM, el conteo cambia sin `F5` en menos de 2 s desde el evento.

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

#### `FRN-20C` (`BRG-05-FRN`, parte Inventario) - Columnas de CPU/RAM y estado en vivo en el Inventario (`RF-03`) (nueva, separada de `FRN-19C`)

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 1.5 h
- **Depende de:** `BAC-22B`, `FRN-20A` y `FRN-17A`.
- **Entregable:**
  1. Columnas de CPU (`%`) y RAM (`GB usados / GB totales`), con `-` si vienen en `null`.
  2. Ante cualquier `TASK_FINISHED` recibido, aunque la acción la haya iniciado otro usuario, volver a consultar la fila o la lista para actualizar estado y métricas. Si la acción fue `delete` y terminó `COMPLETED`, quitar la fila.
- **Criterio de éxito:** si otro usuario apaga una VM, la fila cambia a `Stopped` sin `F5`.

---

### Bloque 5: Frontend - Controles de Energía y Gestión de Estados

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

#### `FRN-16` - Máquina de estados "Operación en progreso" por instancia
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.5 h
- **Depende de:** `FRN-15`. Empieza con `BAC-29` y cierra con `FIX-39`, `BAC-24A` y `BAC-24B`.
- **Entregable:**
  - Al confirmar el modal, enviar la orden de energía (con la ruta documentada en `BAC-29`) o `DELETE`, y guardar el `tareaId` del `202` en el estado de la fila (`transitioning`).
  - El botón accionado muestra un spinner, y se deshabilitan todos los botones de esa instancia.
  - Ante un error, desbloquear la fila y mostrar un mensaje según el código:
    - `409 INSTANCE_BUSY`: "La instancia está ejecutando otra tarea".
    - `409 INSTANCE_INVALID_STATE`.
    - `403 INSTANCE_PROTECTED`.
    - `403 INSTANCE_ACCESS_DENIED`.
    - `502`/`504 PROXMOX_UNAVAILABLE`.
- **Criterio de éxito:** es imposible disparar una segunda acción sobre la misma instancia mientras hay una orden en curso, y cada error muestra su mensaje específico.

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

#### `FRN-17B` - Desbloqueo reactivo de instancia y notificación Toast adaptativa
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `FRN-16`, `FRN-17A` y `BAC-26` (`exitstatus` en `detalles`).
- **Entregable:**
  - Al recibir un `TASK_FINISHED` con el mismo `detalles.tareaId` que una fila en transición, quitar el bloqueo y el spinner.
  - Actualizar el badge de estado (`Running` / `Stopped`).
  - Mostrar un toast (`components/ui/toast.tsx`):
    - verde si `estado: COMPLETED`;
    - rojo con `detalles.error` si `FAILED`.
  - Si la acción fue `delete`, quitar la fila.
- **Criterio de éxito:** la interfaz actualiza el estado y libera los botones sin `F5`.

#### `FRN-16B` (`BRG-04-FRN`) - Resincronización del estado "Operación en progreso" tras recarga (`F5`) o reconexión

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25C`, `FRN-16` y `FRN-17C`.
- **Entregable:**
  1. Al cargar la tabla, poner en `transitioning` (spinner y botones bloqueados) las filas cuyo ítem de `GET /api/instances` trae `activeTask !== null`, guardando su `tareaId`.
  2. Al reconectarse `/api/events` (`FRN-17C`), volver a pedir `GET /api/instances` sin mostrar carga, para tomar los estados finales de las tareas que terminaron durante la desconexión.
- **Criterio de éxito:** si se recarga con `F5` en medio de una acción, la fila sigue con el spinner y los controles bloqueados hasta que llega `TASK_FINISHED` o termina la tarea.

---

### Bloque 6: Integración y Verificación del Hito MVP

#### `INT-01` - Pruebas de integración automatizadas de ciclo de vida y worker UPID
- **Área:** Backend / Testing
- **Asignados:** Tayra y Cristian
- **Estimación:** 2.5 h
- **Depende de:** `BAC-23B`, `BAC-24A`, `BAC-24B`, `BAC-24C`, `BAC-25A`, `BAC-25B`, `BAC-25C`, `BAC-26` y `BAC-27`.
- **Entregable:** suite de integración automática, contra el simulador y Redis local, que cubra:
  1. Envío de cada acción y retorno de `202` con `upid` y `tareaId`. Rechazos `409`/`403`.
  2. Sondeo del worker hasta `COMPLETED` y vencimiento por `UPID_TIMEOUT`.
  3. `TASK_FINISHED` por `/api/events` con el payload de `BAC-29`, y su filtro por permiso.
  4. Los dos registros en `auditoria`.
  5. Reanudación después de reiniciar el backend.
  6. Limpieza de permisos después de un `DELETE`.
- **Criterio de éxito:** la suite pasa al 100 % en local o en CI, sin pasos manuales.

#### `INT-02` - Validación integral del Hito MVP Operativo en local (Smoke Test End-to-End)
- **Área:** Integración / Todo el equipo
- **Asignados:** Equipo completo
- **Estimación:** 2.0 h
- **Depende de:**
  - todas las tareas de desarrollo de la etapa;
  - de la fase base: `FIX-37`, `SEC-03`, `FIX-38` y `FRN-17C`.

  **No depende de infraestructura:** corre en local, con backend, frontend, Redis local y el simulador.
- **Entregable:** prueba de aceptación manual de extremo a extremo, en local:
  1. **Login:** un OPERATOR con `FULL_ACCESS` sobre una instancia inicia sesión con 2FA.
  2. **Dashboard:** ve los recursos del host y los conteos actualizándose sin `F5`.
  3. **Inventario:** solo ve sus instancias, con IP, CPU y RAM. Los filtros funcionan.
  4. **Ciclo de vida:** enciende una VM apagada, confirma el modal y ve el spinner y el bloqueo. Recarga con `F5` y el bloqueo sigue.
  5. **Notificación:** el backend emite el evento, la fila pasa a verde, salta el toast de éxito y el conteo del Dashboard cambia.
  6. **Auditoría:** el ADMIN ve los dos registros de la acción en `auditoria`.
  7. **Borrado:** el ADMIN elimina una instancia detenida y el operador pierde el permiso sobre ese VMID.
- **Criterio de éxito:** el flujo completo funciona sin errores de consola ni inconsistencias de interfaz.

#### `INT-03` - Validación del Hito MVP en el servidor contra el Proxmox real

- **Área:** Integración / Infraestructura
- **Asignados:** Nico y un integrante de backend
- **Estimación:** 1.5 h
- **Depende de:**
  - `INT-02`;
  - `INF-07` e `INF-07B`;
  - `INF-06B` e `INF-08A` (fase base);
  - backend y frontend desplegados.
- **Entregable:** repetir en el entorno de pruebas los pasos 1 a 7 de `INT-02`, contra el Proxmox real y a través de Nginx (HTTPS y `/api/events`). Las acciones de energía y borrado se hacen **solo sobre una instancia de prueba creada para la validación** (VMID ≥ 106, obtenido de `/cluster/nextid`), **nunca sobre los VMIDs protegidos 100 a 105**.
- **Criterio de éxito:** el flujo funciona en el servidor igual que en local:
  - el stream `/api/events` no se corta ni queda en buffer;
  - las IPs se leen del guest agent real;
  - la instancia de prueba se enciende, se apaga y se elimina, con su auditoría registrada;
  - un `stop` sobre la 100 responde `403 INSTANCE_PROTECTED`.

---

## 7. Resumen de Distribución y Carga de Trabajo

**Ya terminadas (no se cargan):** `BAC-28`, `FIX-32`, `FIX-33` y `FIX-35`, en `terminado-1.md`.

| Ola | Tarea | Área | Responsable(s) | Estimación | Empieza con | Cierra con (bloqueantes) |
|---|---|---|---|---|---|---|
| 1 | `INF-07` | Infraestructura | Nico | 1.0 h | — | — |
| 1 | `INF-07B` (BRG-03) | Infraestructura | Nico | 1.0 h | — | — (`FIX-31` terminada) |
| 1 | `BAC-29` (nueva) | Backend | Lisandro | 1.0 h | — | — |
| 1 | `BAC-22` | Backend | Tayra | 2.5 h | — | — |
| 1 | `BAC-23A` | Backend | Lisandro | 2.5 h | — | — |
| 1 | `BAC-25A` | Backend | Lisandro | 2.0 h | — | — |
| 1 | `FRN-19A` | Frontend | Belinda | 2.0 h | — | — |
| 1 | `FRN-20A` | Frontend | Luz | 2.5 h | `BAC-14`, `BAC-29` | `FIX-39` |
| 1 | `FRN-17A` | Frontend | Cristian | 1.5 h | — | `FRN-17C` |
| 2 | `BAC-22B` (BRG-05) | Backend | Tayra | 1.5 h | — | `BAC-22`, `BAC-23A`, `FIX-39` |
| 2 | `BAC-23B` | Backend | Lisandro | 1.5 h | — | `BAC-23A`, `FIX-39` |
| 2 | `BAC-24A` | Backend | Lisandro | 1.5 h | — | `FIX-39` |
| 2 | `BAC-24B` | Backend | Lisandro | 2.0 h | — | `FIX-39` |
| 2 | `BAC-25B` | Backend | Tayra | 1.5 h | — | `BAC-25A` |
| 2 | `FRN-19B` | Frontend | Belinda | 2.0 h | `BAC-29` | `FRN-19A`, `BAC-22` |
| 2 | `FRN-20B` | Frontend | Luz | 2.0 h | — | `FRN-20A` |
| 2 | `FRN-15` | Frontend | Belinda | 2.5 h | `FRN-20A` (modales con datos de prueba) | `SEC-03`, `FIX-37`, `FIX-38` |
| 3 | `BAC-24C` (nueva) | Backend | Tayra | 1.0 h | — | `BAC-24B`, `BAC-25B` |
| 3 | `BAC-25C` (BRG-04) | Backend | Lisandro | 1.5 h | — | `BAC-25A`, `FIX-39` |
| 3 | `BAC-26` | Backend | Tayra | 1.0 h | — | `BAC-25B` |
| 3 | `BAC-27` | Backend | Tayra | 2.0 h | — | `FIX-39`, `BAC-25B`, `BAC-18B` |
| 3 | `FRN-16` | Frontend | Cristian | 2.5 h | `BAC-29` | `FRN-15`, `BAC-24A`, `BAC-24B` |
| 3 | `FRN-19C` (BRG-05) | Frontend | Belinda | 1.0 h | — | `BAC-22B`, `FRN-19B`, `FRN-17A` |
| 3 | `FRN-20C` (BRG-05, nueva) | Frontend | Luz | 1.5 h | — | `BAC-22B`, `FRN-20A`, `FRN-17A` |
| 4 | `FRN-16B` (BRG-04) | Frontend | Cristian | 1.5 h | — | `BAC-25C`, `FRN-16`, `FRN-17C` |
| 4 | `FRN-17B` | Frontend | Cristian | 2.0 h | — | `FRN-16`, `FRN-17A`, `BAC-26` |
| 4 | `INT-01` | Testing | Tayra, Cristian | 2.5 h | — | Bloques 2 y 3 |
| 5 | `INT-02` | Integración | Equipo completo | 2.0 h | — | Toda la etapa, más `FIX-37`, `SEC-03`, `FIX-38` y `FRN-17C` |
| 5 | `INT-03` | Integración / Infra | Nico + backend | 1.5 h | — | `INT-02`, `INF-07`, `INF-07B`, `INF-06B`, `INF-08A` |
| | **Total** | | | **50.5 h** | | **29 tareas, promedio 1.74 h** |

**Carga por persona en la Etapa 1:**

| Persona | Horas | Tareas | Además tiene en la fase base |
|---|---:|---|---|
| Lisandro | 12.0 | `BAC-29`, `BAC-23A`, `BAC-23B`, `BAC-24A`, `BAC-24B`, `BAC-25A`, `BAC-25C` | `FIX-39` (compartida) |
| Tayra | 9.5 + `INT-01` | `BAC-22`, `BAC-22B`, `BAC-25B`, `BAC-24C`, `BAC-26`, `BAC-27` | `FIX-39` (compartida), `BAC-18B` |
| Belinda | 7.5 | `FRN-19A`, `FRN-19B`, `FRN-19C`, `FRN-15` | `SEC-03`, `FIX-38` |
| Luz | 6.0 | `FRN-20A`, `FRN-20B`, `FRN-20C` | `FIX-29` |
| Cristian | 7.5 + `INT-01` | `FRN-16`, `FRN-16B`, `FRN-17A`, `FRN-17B` | `FRN-17C`, `FIX-36`, `SEC-03`, `FIX-38` |
| Nico | 2.0 + `INT-03` | `INF-07`, `INF-07B` | `INF-06B`, `INF-08A` |

> `BAC-23B` pasó de Tayra a Lisandro para equilibrar la carga, porque Tayra además tiene `BAC-18B` y `FIX-39` en la fase base. `FIX-37` (backend, "A definir") conviene asignarla a Lisandro o a Tayra según quién termine primero `FIX-39`, porque toca `permisos_instancia`.
