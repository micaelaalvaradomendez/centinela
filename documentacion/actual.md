### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (10/11/2026):** tras la sincronización de submódulos (`backend` `d0f3f57` y `frontend` `a268063`) y la ejecución completa de las suites de prueba, **todas las tareas de desarrollo de software (Backend y Frontend) están 100 % implementadas y en verde**. En este archivo restan únicamente las tareas de **Infraestructura sobre el host Proxmox VE físico** (`FIX-123`, `INF-14`, `ADP-09`, `ADP-10`).
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
> - `FRN-22`, `FRN-23`, `FRN-26`, `FRN-27`, `FRN-28` y `FIX-142` (resueltos y con tests en verde el 04/11/2026 en `26af24f`: restauración de `useInstanceMetrics`/`getWsTicket`, pestaña `"Métricas"`, modo `"En vivo"`, estado `"Instancia detenida"` y gráfico de Red; `frn-22-use-instance-metrics.test.tsx` 9/9 y `frn-23-metrics-charts.test.tsx` 9/9 en verde).
> - `FIX-91` (resuelto y con tests en verde el 30/10/2026: eliminación de código muerto `InstanceTester.tsx`; `fix-91-dead-code.test.tsx` 2/2 en verde).
> - `FIX-108` y `FIX-124` (resueltos y con tests en verde el 30/10/2026: identificadores opacos `InstanceRef` y sanitización estricta anti-traversal; `fix-108-ids-ofuscados.test.tsx` 11/11 y `ids-ofuscados.test.tsx` 14/14 en verde).
> - `FRN-36` (resuelto y con tests en verde el 30/10/2026: modularización de `InstanceDetailPage`, tabs transitorios y bloqueo de cabecera; `frn-36-29-snapshots-ui.test.tsx` en verde).
> - `FIX-127` (resuelto el 03/11/2026 en `83a19b5`/`14db393`/`30a79a5`/`109f266`: `useOptionalEventsContext` y `instanceRef ?? vmid`).
> - `FIX-141` (resuelto y con tests en verde el 05/11/2026 en `020b499` y `bed77ff`: 43/43 archivos y 277/277 tests unitarios internos en `frontend/centinela` en verde).
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
> - `FIX-131` (resuelto y con tests en verde el 04/11/2026 en `fa92405`/`26af24f`: `netInBps`/`netOutBps` y tarjeta `"Tráfico de Red"` en `InstanceResourceUsage.tsx`; 3/3 en `test/front/etapa3-pendientes.test.tsx`).
> - `FIX-132` (resuelto y con tests en verde el 03/11/2026 en `3d06236`: discriminación de `SAME_PASSWORD` / `WEAK_PASSWORD` en `RecoverPassword.tsx`).
> - `FIX-133` (resuelto y con tests en verde el 03/11/2026 en `fb8a41c`: CPU Overcommit permitido en escalado vertical cuando $\Delta\text{RAM} \le 0$).
> - `FIX-134` (resuelto y con tests en verde el 04/11/2026 en `42093c3`: soporte operativo para instancias en estado `paused` en `stop`, `resume` y `EliminarInstancia`; `TestFIX134_InstanciasPausadas` 4/4 en verde).
> - `FIX-135` (resuelto y con tests en verde el 04/11/2026 en `5f8bb5e`: reconciliación con tareas de Proxmox en `explicadaPorCentinela`, auditoría `SYSTEM_INSTANCE_STOPPED_EXTERNALLY` y `ALERT_IGNORE_EXTERNAL_SHUTDOWN`; `TestFIX135_AlertasParadasExternas` 3/3 en verde).
> - `FIX-136` (resuelto y con tests en verde el 04/11/2026 en `08126ae`: parametrización de umbrales `ConfigAlertas` y `ALERT_THRESHOLD_*` / `ALERT_CYCLES_COUNT`; `TestFIX136_UmbralesDeAlertaConfigurables` 4/4 en verde).
> - `FIX-137` (resuelto y con tests en verde el 05/11/2026 en `5b3ac80`: eliminación de stub `PausarInstancia`, inclusión de eventos de sistema `organizacion_id IS NULL` en `audit_repository.go` y limpieza de `contrato_etapa1.go`; `TestFIX137_CodigoInerteYAuditoriaDeSistema` 3/3 en verde).
> - `FIX-138` (resuelto y con tests en verde el 03/11/2026 en `3fe055b`: ampliación de `CaracteresEspeciales` NIST SP 800-63B en `password.go`).
> - `FIX-139` (resuelto y con tests en verde el 05/11/2026 en `d92a568`: `RecuperarYReconciliarAlArranque` en `seguimiento_tareas.go` e invocación previa en `cmd/api/main.go`; `TestFIX139_ReconciliacionDeTareasAlArranque` 3/3 en verde).
> - `FIX-140` (resuelto y con tests en verde el 03/11/2026 en `4e4d2c4`: `emailProvider` en `/api/version` y cabecera `X-Centinela-Email-Provider: mock`).
> - `FIX-143` (resuelto y con tests en verde el 04/11/2026 en `4011b08`: vocabulario cerrado de `auditoria.accion` y `auditoria.resultado`, enums publicados en Swagger y `contrato-etapa1.md`, test `TestContratoVocabularioAuditoria` y migración histórica `20261104_vocabulario_auditoria.sql`; `TestFIX143_VocabularioAuditoriaContratoCerrado` 4/4 en verde).
> - `FIX-144` (resuelto y con tests en verde el 05/11/2026 en `c110791`: alineación 1 a 1 de `auditActionLabels.ts`, `Auditoria.tsx` y `deletedUsersService.ts` con el enum cerrado `EXITO`/`FALLA`/`PENDIENTE` y retiro de lectura doble; `test/front/etapa3-pendientes.test.tsx` 2/2 en verde).
> - `FRN-33` (resuelto y con tests en verde el 10/11/2026 en `10beb6f` y `a268063`: Centro de Notificaciones con campana, contador de no leídas y banner contextual `role="alert"` en Dashboard; `frn-33-alerts-center.test.tsx` 2/2 en verde).
> - `FIX-129` (`INF-11`, resuelto y con tests en verde el 10/11/2026: cabeceras HSTS `includeSubDomains`, `X-Frame-Options`, `X-Content-Type-Options` y CSP en `docker/nginx-edge.conf`; `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS` 3/3 en verde).
>
> **Pasados a `terminado3.md` (Etapa Adaptación Proxmox PVE 8.4):**
> - `ADP-01`, `ADP-02`, `ADP-03`, `ADP-04`, `ADP-05` y `ADP-08` (resueltos y con tests en verde el 01/11/2026).
> - `ADP-06` (`SIM-PVE8-02`, resuelto y con tests en verde el 03/11/2026 en `aee232a`: catálogo de plantillas e ISOs Debian 12 Bookworm en `storage.go`).
> - `ADP-07` / `FRN-37` (resuelto y con tests en verde el 10/11/2026 en `087be9c` y `da3f422`: selector y extracción de metadatos de plantillas Debian 12 en wizard de creación; `frn-37-adp-07-wizard-debian12.test.tsx` 2/2 en verde).

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

### `FIX-146` - Restauración de columna unívoca (`ID` / `InstanceRef`) y selección en la tabla de inventario (`FRN-20A` / `FIX-43` / `FIX-108`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-20A`, `FIX-43`, `FIX-108`.
- **Problema técnico detallado:**
  En los commits `99c7a18` y `0bf4f33`, se eliminaron silenciosamente las columnas `Checkbox` e `ID` de `Instances.tsx`, modificando los tests unitarios internos (`instances-table.test.tsx`) para ocultar la regresión. Esto incumple el criterio de aceptación de `FIX-43` (*"Muestra columnas separadas para ID y Nombre, permitiendo identificación unívoca"*) y `FIX-108`, impidiendo a los operadores distinguir instancias homónimas o correlacionar identificadores con registros de auditoría y logs de Proxmox.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/pages/Instances.tsx` (y componentes de tabla asociados):
     - Restaurar la columna `ID` exhibiendo el identificador público opaco (`inst_...`) de cada instancia con botón de copia rápida.
     - Restaurar la columna de selección `Checkbox` para operaciones en lote según la especificación de `FRN-20A`.
  2. En `frontend/centinela/src/test/instances-table.test.tsx`:
     - Reincorporar las aserciones que validan la presencia visible de la columna `ID` y las referencias públicas `inst_...`.
- **Criterio de éxito:**
  La tabla de inventario en `/instances` vuelve a exhibir la columna `ID` unívoca y pasa `npm --prefix frontend/centinela test -- --run` al 100%.

---
# ETAPA 3: Tareas Pendientes de Implementación
---

*Todas las tareas de desarrollo de software (Backend y Frontend) de la Etapa 3 fueron completadas y verificadas en [`terminado3.md`](terminado3.md). Resta la siguiente tarea de auditoría sobre el host físico:*

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


### `FIX-145` - Eliminación de la restricción de pertenencia rígida en borrado (`403 INSTANCE_NOT_OWNED`) y desacoplamiento mediante configuración de instancias protegidas (`RF-07` / `BAC-24B`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas / Nico
- **Estimación:** 1.5 h
- **Depende de:** `RF-07`, `BAC-24B`, `PR #30` (`df4cada`).
- **Problema técnico detallado:**
  En el commit `df4cada` (PR #30), se introdujo en `instance_handler.go:EliminarInstancia` una consulta bloqueante `FueCreadaPorCentinela` contra la tabla `tareas_asincronas`, rechazando con `403 INSTANCE_NOT_OWNED` cualquier orden `DELETE` sobre instancias cuya tarea de creación no figure en la base de datos local.
  Esto genera un problema grave de desacople y operación en el servidor real:
  1. Impide a los administradores gestionar o destruir instancias preexistentes en Proxmox, instancias migradas o recreadas tras limpiezas de base de datos.
  2. En el servidor físico los nombres, identificadores y rangos de VMIDs cambian constantemente respecto a los entornos locales de desarrollo.
  3. La protección contra la destrucción accidental de las instancias propias del proyecto Centinela (el propio nodo Docker, DNS, base de datos del host) debe definirse explícitamente en la variable de entorno `PROXMOX_PROTECTED_VMIDS`, que ya responde de forma canónica con `403 INSTANCE_PROTECTED`.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/adapters/primary/http/instance_handler.go`:
     - Retirar de `EliminarInstancia` la guarda `FueCreadaPorCentinela` y la emisión del código `403 INSTANCE_NOT_OWNED`.
     - Mantener la verificación de estado operativo (`stopped` o `paused`) para responder `409 INSTANCE_INVALID_STATE` si la instancia está encendida.
     - Confiar la protección de infraestructura exclusivamente a la lista declarativa `PROXMOX_PROTECTED_VMIDS` en `RejectProtectedInstance(vmidsProtegidos)`.
  2. En `backend/docs/contrato-etapa1.md`:
     - Retirar la cláusula `INSTANCE_NOT_OWNED` del contrato y de la documentación de errores de `DELETE /instances/{instanceRef}`.
  3. En `backend/internal/core/ports/tarea_port.go` y `tarea_repository.go`:
     - Deprecar el método `FueCreadaPorCentinela` si no tiene otros consumidores funcionales.
- **Criterio de éxito:**
  Un usuario con rol `ADMIN` puede solicitar `DELETE` sobre cualquier instancia detenida del hipervisor que no forme parte de `PROXMOX_PROTECTED_VMIDS`, recibiendo `202 Accepted` y despachando el borrado en Proxmox VE.


---

# etapa adaptacion-proxmox

---

*Las tareas `ADP-01` a `ADP-08` fueron completadas y verificadas en [`terminado3.md`](terminado3.md). Restan las tareas sobre el hipervisor físico:*

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
