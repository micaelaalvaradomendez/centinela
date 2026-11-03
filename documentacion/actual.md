### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (03/11/2026):** según la última sincronización de submódulos (`backend` `fb8a41c` y `frontend` `877eddd`), en este archivo quedan únicamente las tareas **pendientes de implementación** de la Etapa 2, Etapa 3 y etapa adaptacion-proxmox. Lo ya verificado con las suites está en [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md) / [`terminado3.md`](terminado3.md), y las correcciones de las tareas implementadas cuyos tests fallan están documentadas en [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md) / [`futuro3.md`](futuro3.md).
>
> **Pasados a `terminado-1.md` / `terminado2.md` (Etapa 1 y 2):**
> - `FIX-104B` / `FIX-105` y `FIX-102` (resueltos y con tests en verde el 21/10/2026).
> - `FIX-100` (resuelto y con tests en verde el 22/10/2026: guarda en `formatCount` y `formatShare`).
> - `FIX-101`, `FIX-103` y `FIX-104` (cerrados el 22/10/2026: expectativas desactualizadas de `test/front`, no defectos del submódulo; pruebas en verde).
> - `BAC-24C` / `FIX-113` (resuelto y con tests en verde el 24/10/2026: purga atómica de permisos al eliminar instancia; pruebas `TestFIX113` y pasos de integración en verde).
> - `BAC-36`, `FIX-107`, `FIX-109`, `FIX-110`, `FIX-111A`, `FIX-112`, `FIX-114`, `FIX-115A`, `FIX-116`, `FIX-117A`, `FIX-118`, `FIX-119` (resueltos y con tests en verde el 21/10/2026).
> - `FIX-111B` y `FIX-115B` (resueltos y con tests en verde el 22/10/2026: modal de reautenticación y claves hasta 128 chars).
> - `FRN-24` y `FRN-25` (cerradas el 22/10/2026: implementación funcional y pruebas alineadas al contrato real).
> - `FIX-117B` (resuelto y con tests en verde el 24/10/2026: adaptación de Dashboard para operador; `fix-117b-dashboard-operador.test.tsx` 4/4 en verde).
> - `FRN-26`, `FRN-27` y `FRN-28` (resueltos y con tests en verde; `FRN-22` y `FRN-23` tienen seguimiento en `FIX-142` de [`futuro2.md`](futuro2.md) tras la reversión en `4f743e4`).
> - `FIX-91` (resuelto y con tests en verde el 30/10/2026: eliminación de código muerto `InstanceTester.tsx`; `fix-91-dead-code.test.tsx` 2/2 en verde).
> - `FIX-108` y `FIX-124` (resueltos y con tests en verde el 30/10/2026: identificadores opacos `InstanceRef` y sanitización estricta anti-traversal; `fix-108-ids-ofuscados.test.tsx` 11/11 y `ids-ofuscados.test.tsx` 14/14 en verde).
> - `FRN-36` (resuelto y con tests en verde el 30/10/2026: modularización de `InstanceDetailPage`, tabs transitorios y bloqueo de cabecera; `frn-36-29-snapshots-ui.test.tsx` en verde).
> - `FIX-127` (resuelto el 03/11/2026 en `83a19b5`/`14db393`/`30a79a5`/`109f266`: `useOptionalEventsContext` y `instanceRef ?? vmid`; seguimiento de mocks internos y ruta `:vmid` en `FIX-141` de [`futuro2.md`](futuro2.md)).
> - `FIX-128` (Backend VLAN 0, resuelto y con tests en verde el 03/11/2026 en `b2d5fec`: normalización de `*input.VLAN == 0` a `nil`).
>
> **Pasados a `terminado3.md` (Etapa 3):**
> - `BAC-42`, `BAC-43`, `BAC-46`, `BAC-47` (resueltos en cierres previos).
> - `BAC-44` (resuelto y con tests en verde el 30/10/2026: interceptor global, filtros compuestos y orden `fecha_hora DESC`).
> - `BAC-45` (resuelto y con tests en verde el 30/10/2026: streaming HTTP con cursor `.Rows()`, sanitización CSV y cabecera `auditoria_centinela_YYYYMMDD`).
> - `BAC-48` (resuelto y con tests en verde el 30/10/2026: cookies `HttpOnly`, `SameSite=Strict`, `Secure` y Rate Limiting 429).
> - `BAC-54` (resuelto y con tests en verde el 30/10/2026: lectura dinámica de `PROXMOX_NODE`, `PROXMOX_URL` y `PROXMOX_PROTECTED_VMIDS`).
> - `FRN-29` (resuelto y con tests en verde el 30/10/2026: pestaña Snapshots, validación de `snapname` y feedback reactivo; `frn-36-29-snapshots-ui.test.tsx` en verde).
> - `FRN-30` (resuelto y con tests en verde el 30/10/2026: modal destructivo de rollback con confirmación antierror tipear nombre idéntico; `frn-30-snapshots-rollback-ui.test.tsx` 2/2 en verde).
> - `FRN-34` (resuelto y con tests en verde el 30/10/2026: adaptación de cliente a cookies seguras de sesión y auditoría WCAG AA; `frn-34-cookies-a11y.test.tsx` 3/3 en verde).
> - `FRN-32` (resuelto y con tests en verde el 03/11/2026 en `404c784`: menú desplegable `AuditExportButton` CSV/JSON con estado de carga, Toast y transferencia de filtros).
> - `BAC-55` (`BRG-15-SIM`, resuelto y con tests en verde el 03/11/2026 en `61b973d`/`ec5922c`: switch `PROXMOX_SIM_TARGET_VERSION`).
> - `FIX-128` (`FRN-31`, resuelto y con tests en verde el 03/11/2026 en `b38ab17`: `SNAPSHOT_DELETE`, `ALERT_TRIGGERED`, `'Revertir snapshot'` y exportación en `Auditoria.tsx`).
> - `FIX-132` (resuelto y con tests en verde el 03/11/2026 en `3d06236`: discriminación de `SAME_PASSWORD` / `WEAK_PASSWORD` en `RecoverPassword.tsx`).
> - `FIX-133` (resuelto y con tests en verde el 03/11/2026 en `fb8a41c`: CPU Overcommit permitido en escalado vertical cuando $\Delta\text{RAM} \le 0$).
> - `FIX-138` (resuelto y con tests en verde el 03/11/2026 en `3fe055b`: ampliación de `CaracteresEspeciales` NIST SP 800-63B en `password.go`).
> - `FIX-140` (resuelto y con tests en verde el 03/11/2026 en `4e4d2c4`: `emailProvider` en `/api/version` y cabecera `X-Centinela-Email-Provider: mock`).
>
> **Pasados a `terminado3.md` (Etapa Adaptación Proxmox PVE 8.4):**
> - `ADP-01`, `ADP-02`, `ADP-03`, `ADP-04`, `ADP-05` y `ADP-08` (resueltos y con tests en verde el 01/11/2026).
> - `ADP-06` (`SIM-PVE8-02`, resuelto y con tests en verde el 03/11/2026 en `aee232a`: catálogo de plantillas e ISOs Debian 12 Bookworm en `storage.go`).
>
> **Seguimiento de Fixes en `futuro2.md`:**
> - `FIX-141` (en [`futuro2.md`](futuro2.md): actualizar mocks de `@/context/EventsContext` con `useOptionalEventsContext` en 2 tests internos de `frontend/centinela/src/test/`).
> - `FIX-142` (en [`futuro2.md`](futuro2.md): restaurar `FRN-22` revertido en `4f743e4` —`useInstanceMetrics`, `getWsTicket`, `metricsWebSocket`— y completar `FRN-23` con pestaña `"Métricas"`, modo `"En vivo"` y serie de Red).

---

### ETAPA 1: Tareas Pendientes de Implementación

*Actualmente no hay tareas pendientes de implementación en la Etapa 1. Todas las tareas fueron completadas y verificadas en [`terminado-1.md`](terminado-1.md).*

---

### ETAPA 2: Tareas Pendientes de Implementación

*Las tareas base de la Etapa 2 y puente hacia Etapa 3 fueron completadas y verificadas en [`terminado2.md`](terminado2.md). Resta la verificación sobre infraestructura física:*

### `FIX-123` - Descarga y verificación de disponibilidad de plantillas LXC e ISOs en Proxmox físico (`INF-10`) (Infraestructura)

- **Área:** Infraestructura
- **Asignado:** Nico / Lucas
- **Estimación:** 1.5 h
- **Depende de:** `INF-10` (en `terminado2.md`).
- **Problema técnico detallado:**
  El servidor Proxmox VE real requiere la presencia verificada de imágenes base en los almacenamientos `local:vztmpl` (ej. Alpine/Debian) y `local:iso` con permisos de lectura en el API Token para que el aprovisionamiento de máquinas virtuales y contenedores opere sin errores de almacenamiento inexistente.
- **Qué debe hacer el equipo de infraestructura:**
  1. Descargar o verificar la plantilla LXC liviana oficial en `local:vztmpl`.
  2. Verificar una imagen ISO en `local:iso`.
  3. Comprobar que `GET /storage/{storage}/content` con el API Token liste ambos recursos.
- **Criterio de éxito:** `GET /api/node/templates` responde con el catálogo real sin errores de permisos.


### `FIX-141` - Actualización de mocks de `@/context/EventsContext` (`useOptionalEventsContext`) en tests unitarios internos del submódulo `frontend` (`FIX-127`) (Testing Frontend)

- **Área:** Frontend (Tests internos del submódulo)
- **Asignado:** Cristian / Luz
- **Estimación:** 0.25 h
- **Depende de:** `FIX-127`, `FRN-29`.
- **Nota de diagnóstico:** La implementación de producto de `FIX-127` en `EventsContext.tsx`, `SnapshotsTabContent.tsx` e `InstanceDetailPage.tsx` **está bien implementada** y pasa las pruebas de aceptación (`test/front/etapa3-pendientes.test.tsx` 3/3 ✅ y `test/front/instance-detail.test.tsx` 3/3 ✅). El fallo ocurre únicamente en 2 archivos de test internos del submódulo `frontend/centinela/src/test/`.
- **Tests internos del submódulo que fallan actualmente:**
  - `frontend/centinela/src/test/snapshots.test.tsx` (4/4 fallan con `Error: [vitest] No "useOptionalEventsContext" export is defined on the "@/context/EventsContext" mock`)
  - `frontend/centinela/src/test/instance-transient-states.test.tsx` (1/5 falla con `Error: [vitest] No "useOptionalEventsContext" export is defined on the "@/context/EventsContext" mock`)
- **Problema técnico detallado (Por qué fallan los tests del submódulo):**
  En el commit `14db393` (`FIX-127`), `frontend/centinela/src/components/features/instances/components/detail/SnapshotsTabContent.tsx:69` reemplazó `useEventsContext()` por `useOptionalEventsContext()`. Sin embargo, los mocks manuales `vi.mock('@/context/EventsContext', ...)` definidos en `frontend/centinela/src/test/snapshots.test.tsx:15-17` y `frontend/centinela/src/test/instance-transient-states.test.tsx:15-18` quedaron con la exportación vieja (`useEventsContext`) sin exportar `useOptionalEventsContext`.
- **Solución técnica paso a paso:**
  En `frontend/centinela/src/test/snapshots.test.tsx:15-17` y `frontend/centinela/src/test/instance-transient-states.test.tsx:15-18`, agregar `useOptionalEventsContext` al objeto devuelto por `vi.mock('@/context/EventsContext')`:
  ```ts
  vi.mock('@/context/EventsContext', () => ({
      useEventsContext: () => ({ ultimoMensaje: null }),
      useOptionalEventsContext: () => ({ ultimoMensaje: null }),
      useEventSubscription: () => {},
  }));
  ```
- **Criterio de éxito:**
  - `frontend/centinela/src/test/snapshots.test.tsx` (4/4) y `frontend/centinela/src/test/instance-transient-states.test.tsx` (5/5) pasan al 100% en verde.

---

### `FIX-142` - Restauración de `FRN-22` (`useInstanceMetrics`, `getWsTicket`, `metricsWebSocket`) revertida en `4f743e4` e integración de pestaña `"Métricas"`, modo `"En vivo"` y gráfico de Red en `FRN-23` (Frontend)

- **Área:** Frontend
- **Asignado:** Luz / Cristian / Belinda
- **Estimación:** 2.5 h
- **Depende de:** `FRN-21`, `FRN-22`, `FRN-23`, `BAC-35`, `BAC-36`.
- **Tests que fallan actualmente:**
  - `test/front/frn-22-use-instance-metrics.test.tsx` (5/5 fallan: `no hay un servicio exportado getWsTicket() en el frontend`, `no existe el hook useInstanceMetrics en el frontend (FRN-22)`)
  - `test/front/frn-23-metrics-charts.test.tsx` (7/7 fallan: `la vista de detalle no tiene la pestaña "Métricas"`)
- **Problema técnico detallado (Por qué fallan los tests):**
  1. **Reversión de `FRN-22` en `origin/main` (`4f743e4`):** La tarea `FRN-22` fue implementada originalmente en el commit `60833b3` (`🌱 feat(instances):se implementa telemetría en tiempo real de instancias.`) y pasada a `terminado2.md`. Sin embargo, en el commit `4f743e4` (`Revert "🌱 feat(instances):se implementa telemetría en tiempo real de instancias."`) se eliminaron por completo del árbol `src/` los 5 archivos de la tarea:
     - `src/components/features/instances/metrics/metrics.types.ts`
     - `src/components/features/instances/metrics/metricsService.ts` (`getWsTicket()`)
     - `src/components/features/instances/metrics/metricsWebSocket.ts`
     - `src/components/features/instances/metrics/useInstanceMetrics.ts` (`useInstanceMetrics(vmid)`)
     - `src/test/instance-metrics.test.tsx`
     Como consecuencia, `useInstanceMetrics` y `getWsTicket` no existen actualmente en `origin/main` (`877eddd`).
  2. **Incumplimiento de entregables de `FRN-23` en `InstanceDetailContent.tsx` e `InstanceResourceUsage.tsx`:**
     - En `frontend/centinela/src/components/features/instances/components/detail/InstanceDetailContent.tsx:116-120`, las pestañas se nombraron `"Resumen"`, `"Recursos"`, `"Red"` y `"Snapshots"`, omitiendo la pestaña **`"Métricas"`** exigida por `FRN-21` y `FRN-23` (`test/front/frn-23-metrics-charts.test.tsx:118`).
     - En `frontend/centinela/src/components/features/instances/components/detail/InstanceResourceUsage.tsx:19-22`, el selector de rango (`timeframeOptions`) solo ofrece `'2h'` (`"Últimas 2 horas"`) y `'24h'` (`"Últimas 24 horas"`), omitiendo la opción **`"En vivo"`** conectada al WebSocket de `useInstanceMetrics` (`FRN-23` entregable 2).
     - En `InstanceResourceUsage.tsx:28-57`, solo se grafican 3 métricas (`CPU`, `Memoria`, `Disco`), omitiendo el gráfico de **Tráfico de Red (`netInBps` / `netOutBps`)** (`FRN-23` entregable 1 y `FIX-131`).
     - Cuando la instancia está detenida (`detail.status === 'stopped'` o muestras con `cpuPercent: null` y `ramUsedBytes: null`), `InstanceResourceUsage.tsx` no muestra el estado vacío explícito **`"Instancia detenida"`** (`FRN-23` entregable 5).
- **Solución técnica paso a paso:**
  1. Restaurar desde el commit `60833b3` (o `01c1567`) los módulos `metrics.types.ts`, `metricsService.ts` (exportando `getWsTicket()` que invoca `POST /auth/ws-ticket`), `metricsWebSocket.ts`, `useInstanceMetrics.ts` y `src/test/instance-metrics.test.tsx`, adaptándolos para aceptar `InstanceRef` (`string | number`) sin convertir IDs opacos a `Number`.
  2. En `InstanceDetailContent.tsx:116-120`, exponer la pestaña `"Métricas"` (o incluir `"Métricas"` en el nombre accesible / pestaña dedicada que renderice `InstanceResourceUsage`).
  3. En `InstanceResourceUsage.tsx` y `useInstanceDetail.ts`:
     - Añadir la opción `{ value: 'live', label: 'En vivo' }` en `timeframeOptions`, consumiendo `useInstanceMetrics(detail.reference)` cuando `timeframe === 'live'` y `GET /api/instances/:ref/metrics?timeframe=2h|24h` en los modos históricos.
     - Incorporar la 4.ª tarjeta/gráfico de **Tráfico de red** (`netInBps` / `netOutBps` formateados en `KB/s` o `MB/s`).
     - Mostrar el aviso `"Instancia detenida"` cuando `detail.status === 'stopped'` o todos los puntos de CPU/RAM sean `null`, sin producir `NaN` ni invocar `.toFixed()` sobre `null`.
- **Criterio de éxito:**
  - `test/front/frn-22-use-instance-metrics.test.tsx` (5/5) pasa al 100% en verde.
  - `test/front/frn-23-metrics-charts.test.tsx` (7/7) pasa al 100% en verde.

---
# ETAPA 3: Tareas Pendientes de Implementación
---

*Las tareas `BAC-44`, `BAC-45`, `BAC-48`, `BAC-54`, `BAC-55`, `FRN-29`, `FRN-30`, `FRN-32`, `FRN-34`, `FIX-128` (`FRN-31`), `FIX-132`, `FIX-133`, `FIX-138` (Backend) y `FIX-140` fueron completadas y verificadas en [`terminado3.md`](terminado3.md). Restan las siguientes tareas pendientes de la Etapa 3:*

#### `FRN-33` - Centro de Notificaciones y banners contextuales de alerta por saturación
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** `BAC-47`.
- **Problema:** Los administradores en el frontend necesitan identificar de un vistazo las alertas activas del sistema sin tener que revisar manualmente cada gráfica.
- **Entregable:**
  1. Campana / Centro de Notificaciones en la barra superior con contador de alertas no leídas.
  2. Banner contextual persistente de color rojo en la cabecera del Dashboard ante alertas de severidad `CRITICAL`.
  3. Badge de advertencia en la fila del inventario si una instancia específica tiene alertas de saturación activas.
  4. Sonido o feedback visual sutil al recibir un evento `ALERT_TRIGGERED` en vivo.
- **Criterio de éxito:** Al emitirse una alerta de saturación, el banner aparece automáticamente en la pantalla de los administradores sin requerir `F5`.

#### `INF-14` (`BRG-12-INF`) - Auditoría Preliminar y Recolección de Telemetría Cruda en Servidor PVE 8.4.21 Físico
- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 2.0 h
- **Depende de:** `INT-09`.
- **Problema:** Se necesita el inventario real y las firmas JSON exactas que devuelve el nuevo servidor de producción PVE 8.4.21 para alimentar las fixtures de desarrollo y detectar posibles extensiones de storage no documentadas.
- **Entregable:**
  1. Script de recolección segura `scripts/pve8-audit-probe.sh` que consulte con token de prueba los endpoints:
     - `GET /api2/json/version`
     - `GET /api2/json/nodes/{node}/status`
     - `GET /api2/json/cluster/resources`
     - `GET /api2/json/nodes/{node}/storage`
     - `GET /api2/json/nodes/{node}/storage/{storage}/content`
  2. Archivar las capturas sanitizadas (sin tokens ni secretos) en `backend/API proxmox respuestas/pve8-capturas/`.
- **Criterio de éxito:** Carpeta de capturas creada con las respuestas reales del host físico 8.4.21 documentadas.

#### `FRN-37` (`BRG-11-FRN`) - Desacoplamiento de Constantes de Plantillas de SO en Asistente de Creación
- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 1.5 h
- **Depende de:** `FRN-30`.
- **Problema:** El wizard de creación (`useCreateInstanceWizard.ts` y `createInstanceOptions.ts`) contiene referencias duras a nombres de plantillas Debian 13.
- **Entregable:**
  1. Abstracción del parser de plantillas en `frontend/centinela/src/components/features/instances/data/createInstanceOptions.ts`:
     - Normalizador que clasifique cualquier plantilla devuelta por el backend según su familia (`Debian`, `Ubuntu`, `Alpine`, `Otro`) extrayendo la versión mayor dinámicamente (`debian-12`, `debian-13`, etc.).
  2. Pruebas unitarias de frontend verificando que el selector renderice correctamente plantillas Debian 12 Bookworm.
- **Criterio de éxito:** El wizard muestra la tarjeta y badge "Debian 12" al recibir `debian-12-standard...` sin depender de valores estáticos Debian 13.

---

# etapa adaptacion-proxmox

---

*Las tareas `ADP-01`, `ADP-02`, `ADP-03`, `ADP-04`, `ADP-05`, `ADP-06` y `ADP-08` fueron completadas y verificadas en [`terminado3.md`](terminado3.md). Restan las siguientes tareas pendientes de adaptación:*

#### `ADP-07` (`FRN-PVE8-01`) - Adaptación del Wizard de Creación para Nombres de Plantillas Debian 12
- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** `ADP-06`.
- **Problema:** El asistente de aprovisionamiento (`useCreateInstanceWizard.ts` y componentes asociados) asume plantillas de prueba Debian 13 en sus placeholders y validaciones de interfaz.
- **Entregables:**
  1. En `frontend/centinela/src/components/features/instances/hooks/useCreateInstanceWizard.ts`:
     - Asegurar que la selección de plantillas maneje dinámicamente cualquier versión de Debian (11, 12, 13) o Ubuntu sin parseos rígidos de versión.
  2. En `InstanceCreationSummary.tsx`:
     - Formatear adecuadamente los identificadores de volid de Debian 12 (`local:vztmpl/debian-12-standard...`).
- **Criterio de éxito:** El wizard permite seleccionar y resumir plantillas de Debian 12 sin advertencias ni errores en consola.

#### `ADP-09` (`INF-PVE8-01`) - Verificación y Configuración de Privilegios del Token en PVE 8.4.21
- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 2.0 h
- **Depende de:** Ninguna.
- **Problema:** El servidor productivo debe contar con el usuario `centi-api@pve`, el rol `CentinelaAdmin` y el API Token con los permisos exactos requeridos por las Etapas 1, 2 y 3.
- **Entregables:**
  1. Conexión SSH o consola en el servidor Proxmox VE 8.4.21 (`10.x.x.x`):
     - Creación del usuario: `pveum user add centi-api@pve`.
     - Creación del rol `CentinelaAdmin` con:  
       `VM.Allocate`, `Datastore.AllocateSpace`, `VM.Config.CPU`, `VM.Config.Memory`, `VM.Config.Disk`, `VM.Config.Network`, `VM.PowerMgmt`, `VM.Audit`, `VM.Snapshot`, `VM.Snapshot.Rollback`, `Datastore.Audit`.
     - Generación del token `pveum user token add centi-api@pve backend-token --privsep 0`.
  2. Actualizar `documentacion/api-proxmox.md` con las instrucciones específicas para PVE 8.4.21.
- **Criterio de éxito:** Ejecución de `curl -H "Authorization: PVEAPIToken=..." https://<PVE8-HOST>:8006/api2/json/cluster/resources` responde 200 con el inventario del host.

#### `ADP-10` (`INF-PVE8-02`) - Aprovisionamiento de Storage LVM-Thin y Descarga de Plantillas Debian 12
- **Área:** Infraestructura
- **Asignado:** Nico / Lucas
- **Estimación:** 2.0 h
- **Depende de:** `ADP-09`.
- **Problema:** Para que las funciones de aprovisionamiento rápido y Snapshots (`RF-06`, `INF-13`) operen sin fallas en PVE 8.4.21, el storage debe soportar snapshots en caliente y contar con las plantillas descargadas en disco.
- **Entregables:**
  1. Verificar que `local-lvm` sea de tipo `lvmthin` (o que exista un pool `zfspool`).
  2. Descargar la plantilla oficial de Debian 12 en `local:vztmpl`:  
     `pveam update && pveam download local debian-12-standard_12.7-1_amd64.tar.zst`.
  3. Colocar una ISO de instalación de Debian 12 en `local:iso`.
- **Criterio de éxito:** `pveam list local` muestra la plantilla de Debian 12 y `GET /nodes/{node}/storage/local/content` la lista con su volid.

---

### `FIX-129` - Endurecimiento de cabeceras de seguridad en Nginx de producción (`INF-11`) (Infraestructura)

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.0 h
- **Depende de:** `INF-11`.
- **Problema técnico detallado:**
  En el archivo de configuración `docker/nginx-edge.conf:36`:
  1. La cabecera `Strict-Transport-Security` está configurada como `"max-age=31536000"` pero omite la directiva obligatoria `includeSubDomains`.
  2. Faltan las directivas de seguridad para mitigar clickjacking y sniffing de MIME:
     - `add_header X-Frame-Options "DENY" always;`
     - `add_header X-Content-Type-Options "nosniff" always;`
  3. Falta la directiva `Content-Security-Policy` ajustada para scripts, fuentes locales y canales bidireccionales WebSocket (`connect-src 'self' wss: ws:;`).
- **Qué debe hacer el equipo de infraestructura:**
  En `docker/nginx-edge.conf`, agregar y ajustar las cabeceras requeridas en el bloque `server` HTTPS (puerto 443):
  ```nginx
  add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
  add_header X-Frame-Options "DENY" always;
  add_header X-Content-Type-Options "nosniff" always;
  add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; connect-src 'self' ws: wss:; img-src 'self' data:;" always;
  ```
- **Criterio de éxito:**
  La prueba `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS` en `test/back/etapa3_acceptance_test.go` pasa todos sus subtests en verde.

---



---

### `FIX-131` - Incorporación de métricas y gráfico de Tráfico de Red (`netin`/`netout`) en detalle de instancia (`CU-05`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1.5 h
- **Depende de:** `FRN-22`, `FRN-23`.
- **Problema técnico detallado:**
  En el flujo de **CU-05** (`casosdeuso.md:275`), se especifica que:
  > *"El frontend renderiza gráficos de líneas (CPU, RAM y Tráfico de Red) mostrando la evolución temporal."*
  
  El backend en `domain.PuntoMetricaDTO` ya expone `netInBps` y `netOutBps` provenientes de Proxmox RRDdata. Sin embargo, en el frontend (`frontend/centinela/src/components/features/instances/types/instanceDetail.ts`), la interfaz `InstanceMetricPoint` omite las propiedades de red, y `InstanceResourceUsage.tsx` sólo instancia columnas y curvas para CPU, Memoria y Disco, dejando sin cobertura la gráfica de tráfico de red exigida por el caso de uso.
- **Qué debe hacer el equipo de frontend:**
  1. En `instanceDetail.ts`:
     - Agregar a `InstanceMetricPoint`: `netInBps: number | null; netOutBps: number | null;`.
  2. En `instanceDetailService.ts`:
     - Normalizar los campos `netInBps` y `netOutBps` en el mapeador de puntos métricos.
  3. En `InstanceResourceUsage.tsx`:
     - Incorporar la métrica y columna de **Tráfico de Red** utilizando el icono `Network` o `ArrowUpDown`, formateando las tasas con una función auxiliar `formatNetworkRate(bytesPerSec)` (`KB/s`, `MB/s`) e instanciando un `InstanceResourceChart` para la serie de tráfico de red combinado (`(netInBps ?? 0) + (netOutBps ?? 0)` o líneas segregadas de entrada/salida).
- **Criterio de éxito:**
  Al ingresar a la vista de detalle de una instancia con métricas históricas (2h o 24h), la grilla de uso de recursos renderiza las 4 métricas (CPU, Memoria, Disco y Tráfico de Red) con sus respectivos gráficos de evolución temporal.

---

### `FIX-134` - Soporte operativo para instancias en estado `paused` y prevención de deadlocks en controles de energía y borrado (`RF-03` / `RF-04` / `CU-04`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas / Nico
- **Estimación:** 2.0 h
- **Depende de:** `RF-03`, `RF-04`, `CU-04`.
- **Problema técnico detallado:**
  En `documentacion/archvios/centinela.md` (línea 67) y `casosdeuso.md` (CU-03/CU-04), se define explícitamente que las instancias pueden presentar los estados `Running`, `Stopped` o `Paused`. Un hipervisor Proxmox puede pausar una VM automáticamente ante saturación de almacenamiento (I/O freeze), operaciones de backup/snapshot en suspensión, o por intervención manual de un sysadmin en la GUI nativa de Proxmox.
  
  Sin embargo, en `internal/adapters/primary/http/instance_handler.go:209-253`:
  ```go
  var estadosRequeridosPorAccion = map[string]string{
      "start":    "stopped",
      "stop":     "running",
      "shutdown": "running",
      "reboot":   "running",
  }
  ```
  Y en `EliminarInstancia` (`instance_handler.go:1016`):
  ```go
  if instancia.Estado != "stopped" {
      SendError(c, http.StatusConflict, "INSTANCE_INVALID_STATE", "La instancia debe estar detenida...")
  }
  ```
  Si una instancia se encuentra en estado `"paused"`:
  1. `POST /instances/:ref/start` falla con `409 INSTANCE_INVALID_STATE` (exige `stopped`).
  2. `POST /instances/:ref/stop` (apagado forzado) falla con `409 INSTANCE_INVALID_STATE` (exige `running`).
  3. `POST /instances/:ref/status/shutdown` y `reboot` fallan con `409 INSTANCE_INVALID_STATE`.
  4. `DELETE /instances/:ref` falla con `409 INSTANCE_INVALID_STATE`.
  5. La acción `resume` o `unpause` no existe en `CambiarEstado`.
  
  Esto crea un **deadlock operativo completo (falso positivo de incompatibilidad de estado)**: una vez que una máquina se pausa en el hipervisor, El Centinela no permite iniciarla, ni apagarla forzosamente para recuperarla, ni eliminarla.
- **Qué debe hacer el equipo de backend:**
  1. En `instance_handler.go`:
     - Modificar la función `validarEstadoParaAccion` y la matriz `estadosRequeridosPorAccion` para que la acción **`stop` (apagado forzado)** admita tanto `running` como `paused`:
       ```go
       func esEstadoCompatible(accion, estadoActual string) bool {
           switch accion {
           case "start":
               return estadoActual == "stopped"
           case "stop":
               return estadoActual == "running" || estadoActual == "paused"
           case "shutdown", "reboot":
               return estadoActual == "running"
           case "resume":
               return estadoActual == "paused"
           default:
               return false
           }
       }
       ```
     - En `CambiarEstado`, incorporar la acción `"resume"` en el selector de acciones admitidas, despachando la llamada `resume` a Proxmox VE (`POST /nodes/{node}/{type}/{vmid}/status/resume`).
  2. En `internal/adapters/secondary/proxmox/client.go`:
     - Implementar el método `Resume(ctx context.Context, node string, vmid int, vmType string) (string, error)` que envíe la petición al endpoint `/status/resume` de Proxmox y retorne el `UPID` de la tarea.
  3. Actualizar `ports.ProxmoxPort` y las pruebas de `instance_handler_test.go` para cubrir la recuperación de VMs pausadas mediante `stop` forzado y mediante `resume`.
- **Criterio de éxito:**
  Una máquina virtual en estado `paused` puede ser reanudada exitosamente mediante `POST /api/instances/:ref/status/resume` o apagada forzosamente mediante `POST /api/instances/:ref/stop`, retornando `202 Accepted` con su respectivo `tareaId`, sin ser bloqueada con `409 INSTANCE_INVALID_STATE`.

---

### `FIX-135` - Reconciliación con tareas de Proxmox en Evaluador de Alertas para prevenir falsos positivos de `instance_down` ante paradas directas en el hipervisor (`RF-11` / `CU-11`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas / Belinda
- **Estimación:** 2.0 h
- **Depende de:** `RF-11`, `CU-11`.
- **Problema técnico detallado:**
  En `internal/core/services/alert_evaluator_worker.go:249-251`:
  ```go
  if antes, ok := anterior[vmid]; ok && antes.Estado == "running" && !e.explicadaPorCentinela(ctx, vmid) {
      e.disparar(ctx, activa, m, TipoAlertaCaida, fmt.Sprintf("La instancia %s se detuvo sin una orden de El Centinela.", ahora.Nombre))
  }
  ```
  La función `explicadaPorCentinela()` consulta exclusivamente la tabla local `tareas_asincronas` de PostgreSQL dentro de una ventana de 10 minutos (`ventanaTareasCentinela`).
  
  En infraestructuras reales, administradores de sistemas y operadores de nivel superior realizan frecuentemente tareas de mantenimiento directamente en Proxmox VE (desde la consola web oficial, comandos CLI como `qm stop`, `pct stop`, o scripts automáticos de backup en frío del clúster). Debido a que estas tareas legítimas no se originan en la base de datos de El Centinela, el worker asume erróneamente que la máquina se cayó por un fallo catastrófico, disparando una **alerta CRITICAL (`instance_down`)** y despachando correos masivos a todos los administradores (`📧 [ALERTAS] correo de la alerta ... enviado`).
  
  Esto representa un **falso positivo severo de alarma de caída** que deteriora la confianza operativa en el sistema de alertas.
- **Qué debe hacer el equipo de backend:**
  1. En `internal/core/services/alert_evaluator_worker.go`:
     - En `explicadaPorCentinela()`, si la base de datos local no registra una tarea reciente, agregar una verificación secundaria contra la API de tareas de Proxmox VE (`GET /nodes/{node}/tasks?vmid={vmid}`).
     - Si Proxmox reporta una tarea reciente de tipo `qmstop`, `qmshutdown`, `vzstop` o `vzdump` finalizada con estado `OK` en los últimos 10 minutos:
       - No clasificar el evento como caída no autorizada (`instance_down`).
       - En su lugar, registrar un evento de auditoría informativa `SYSTEM_INSTANCE_STOPPED_EXTERNALLY` (Severidad `INFO`), sin disparar la alerta crítica ni enviar correos SMTP de emergencia.
  2. Agregar una variable de entorno opcional `ALERT_IGNORE_EXTERNAL_SHUTDOWN=true|false` (por defecto `false` para entornos donde Centinela es el único gestor, `true` para entornos híbridos).
  3. Crear tests unitarios en `alert_evaluator_test.go` simulando la parada de una VM con tarea externa en Proxmox, comprobando que se inhibe la alerta `instance_down`.
- **Criterio de éxito:**
  Cuando una VM se apaga mediante una orden directa en el nodo Proxmox, el evaluador de alertas no emite un evento `ALERT_TRIGGERED` con tipo `instance_down` ni despacha correos críticos a los administradores.

---

### `FIX-136` - Parametrización por variables de entorno de umbrales hardcodeados y persistencia en el Evaluador de Alertas (`RF-11` / `RNF-07`) (Backend / Infraestructura)

- **Área:** Backend / Infraestructura
- **Asignado:** Nico / Lucas
- **Estimación:** 1.5 h
- **Depende de:** `RF-11`, `RNF-07`.
- **Problema técnico detallado:**
  En `internal/core/services/alert_evaluator_worker.go:24-35`, todos los umbrales de saturación y parámetros de histéresis están codificados como constantes fijas en el código Go:
  ```go
  const (
      umbralHostRAM       = 85.0
      umbralHostCPU       = 90.0
      umbralHostDisco     = 90.0
      umbralInstanciaMem  = 95.0
      umbralInstanciaCPU  = 95.0
      umbralRecuperacion  = 80.0
      ciclosParaAlertar   = 2
      ttlAlertaActiva     = 24 * time.Hour
      ventanaTareasCentinela = 10 * time.Minute
  )
  ```
  En servidores de cómputo intensivo o bases de datos, operar al 88% de RAM o 92% de CPU es un comportamiento esperado y regular. Al no poder configurar estos valores desde `.env`, el sistema genera **falsos positivos continuos de saturación** (`host_saturation`, `resource_saturation`), saturando el registro de auditoría y la casilla de correo de los administradores. Asimismo, en hosts de pruebas con recursos limitados no es posible ajustar umbrales más bajos para verificar las alertas.
- **Qué debe hacer el equipo de backend:**
  1. En `internal/core/services/alert_evaluator_worker.go`:
     - Crear un struct de configuración inyectable:
       ```go
       type ConfigAlertas struct {
           UmbralHostRAM       float64
           UmbralHostCPU       float64
           UmbralHostDisco     float64
           UmbralInstanciaMem  float64
           UmbralInstanciaCPU  float64
           UmbralRecuperacion  float64
           CiclosParaAlertar   int
           VentanaTareas       time.Duration
           Intervalo           time.Duration
       }
       ```
     - En `NewEvaluadorAlertas`, recibir `ConfigAlertas` en lugar de leer constantes directas.
  2. En `cmd/api/main.go`:
     - Leer las variables de entorno con fallbacks a los defaults recomendados:
       - `ALERT_THRESHOLD_HOST_RAM` (default: 85.0)
       - `ALERT_THRESHOLD_HOST_CPU` (default: 90.0)
       - `ALERT_THRESHOLD_HOST_DISK` (default: 90.0)
       - `ALERT_THRESHOLD_VM_RAM` (default: 95.0)
       - `ALERT_THRESHOLD_VM_CPU` (default: 95.0)
       - `ALERT_THRESHOLD_RECOVERY` (default: 80.0)
       - `ALERT_CYCLES_COUNT` (default: 2)
  3. En `backend/.env.example` y `SETUP.md`:
     - Documentar los nuevos parámetros opcionales del Evaluador de Alertas.
- **Criterio de éxito:**
  Al configurar `ALERT_THRESHOLD_HOST_RAM=92` en el archivo `.env`, el host operando al 88% de memoria RAM no dispara ninguna alerta crítica `host_saturation`.

---

### `FIX-137` - Saneamiento de código inerte: eliminación de stub documental `PausarInstancia` y simplificación de residuos multi-tenancy (`RF-12` / `RNF-05`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas / Micaela
- **Estimación:** 1.5 h
- **Depende de:** `RF-12`, `RNF-05`.
- **Problema técnico detallado:**
  1. **Stub muerto `PausarInstancia`:** En `internal/adapters/primary/http/instance_handler.go:872-896`, la función `PausarInstancia` contiene anotaciones Swagger para una ruta `/instances/{instanceRef}/pause` que no está registrada en el motor de Gin en `main.go`. Esto genera confusión en la documentación pública de Swagger (expone un endpoint planificado que en la práctica da 404 Not Found genérico de Gin o no rutea) e incrementa el código muerto.
  2. **Residuos de Multi-Tenancy (`RF-12` descartado):** Aunque el soporte multi-organización fue formalmente descartado del MVP en `requerimientos.md:280` (`RF-12`), el backend arrastra una sobrecarga arquitectónica en `UserRepository`, `AuditRepository` y `jwt.go` obligando a poblar y validar `organizacion_id`. En particular, en `audit_repository.go:109`:
     ```sql
     Where("(a.organizacion_id = ? OR (a.organizacion_id IS NULL AND u.organizacion_id = ?))", orgID, orgID)
     ```
     Si un evento de auditoría del sistema no tiene `usuario_id` ni `organizacion_id`, queda invisible para los administradores en el listado de auditoría, generando un **falso negativo de registro de auditoría**.
- **Qué debe hacer el equipo de backend:**
  1. En `instance_handler.go`:
     - Si la acción Pause no forma parte del alcance operativo inmediato, remover la función stub `PausarInstancia` y sus anotaciones Swagger, actualizando la documentación de Swagger (`swag init`) para reflejar solo los endpoints operativos.
  2. En `internal/adapters/secondary/postgres/audit_repository.go`:
     - Asegurar que la consulta de auditoría capture los eventos generados por el sistema (`usuario_id IS NULL`) cuando el administrador consulta los registros, sin requerir una coincidencia estricta de `organizacion_id` que los oculte:
       ```sql
       Where("a.organizacion_id = ? OR a.organizacion_id IS NULL OR (u.organizacion_id = ?)", orgID, orgID)
       ```
  3. Eliminar DTOs y tipos huérfanos sin uso en `internal/adapters/primary/http/contrato_etapa1.go` que generen duplicación con `ports` y `responses.go`.
- **Criterio de éxito:**
  Los eventos de auditoría disparados automáticamente por el sistema (`ALERT_TRIGGERED`, `ALERT_RESOLVED`, limpiezas automáticas) son visibles para el rol `ADMIN` al invocar `GET /api/admin/audit` sin ser filtrados por omisión de `organizacion_id`.

---

### `FIX-139` - Reconciliación inmediata de tareas al inicio del backend para mitigar falsos positivos de `409 INSTANCE_BUSY` (`RF-04` / `RNF-04`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas / Nico
- **Estimación:** 1.5 h
- **Depende de:** `RF-04`, `RNF-04`.
- **Problema técnico detallado:**
  En `internal/core/services/exclusion_instancias.go:68-75`, el middleware `RequireInstanceIdle` consulta las tareas en curso en la base de datos:
  ```go
  activas, err := s.tareas.BuscarTareasActivasPorVmids(ctx, []int{vmid})
  if _, enCurso := activas[vmid]; enCurso {
      liberar()
      return nil, ports.ErrInstanciaOcupada
  }
  ```
  Si el contenedor Docker o servicio del backend se reinicia (por despliegue, actualización o reinicio del host) mientras una tarea asíncrona de Proxmox estaba ejecutándose (o si Proxmox ya la finalizó pero el evento de finalización no se guardó antes del reinicio), la fila en `tareas_asincronas` permanece con `estado = 'RUNNING'`.
  
  Cuando el backend vuelve a encender, el reconciliador periódico de `SeguimientoTareas` (`seguimiento_tareas.go:480`) corre en un ticker de fondo con intervalos de 30 segundos y espera hasta `UPID_TIMEOUT` (3 minutos). Durante todo ese tiempo, cualquier usuario que intente encender, apagar o reconfigurar la VM recibe un **falso positivo de `409 INSTANCE_BUSY`**, creyendo que la instancia está bloqueada cuando en realidad la tarea ya concluyó en Proxmox VE.
- **Qué debe hacer el equipo de backend:**
  1. En `internal/core/services/seguimiento_tareas.go`:
     - Exponer un método de inicialización síncrono o de alta prioridad: `RecuperarYReconciliarAlArranque(ctx context.Context)`.
     - Este método debe ejecutarse al inicio de `Iniciar(ctx)` antes de aceptar tráfico HTTP, consultando a Proxmox el estado real de todas las tareas en `RUNNING` persistidas en la BD.
     - Si la tarea ya culminó en Proxmox (`status: "stopped"`), marcarla de inmediato como `COMPLETED` o `FAILED` en PostgreSQL, liberando la instancia en el segundo cero del arranque.
  2. En `cmd/api/main.go`:
     - Invocar la reconciliación inicial de tareas antes de abrir el puerto HTTP (`srv.ListenAndServe()`).
  3. Probar el escenario mediante un test en `seguimiento_tareas_test.go` que simule un reinicio del proceso con tareas `RUNNING` pendientes cuya ejecución ya finalizó en Proxmox.
- **Criterio de éxito:**
  Al reiniciar el backend con tareas huérfanas en la base de datos, las tareas finalizadas en Proxmox son saneadas inmediatamente al iniciar y las instancias quedan operables sin demoras de 3 minutos ni respuestas `409 INSTANCE_BUSY`.


