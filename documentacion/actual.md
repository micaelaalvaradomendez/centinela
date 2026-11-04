### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (04/11/2026):** según la última sincronización de submódulos (`backend` `4011b08` y `frontend` `1e5ec49`), en este archivo quedan únicamente las tareas **pendientes de implementación** de la Etapa 2, Etapa 3 y etapa adaptacion-proxmox. Lo ya verificado con las suites está en [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md) / [`terminado3.md`](terminado3.md), y las correcciones de las tareas implementadas cuyos tests fallan están documentadas en [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md) / [`futuro3.md`](futuro3.md).
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
> - `FIX-127` (resuelto el 03/11/2026 en `83a19b5`/`14db393`/`30a79a5`/`109f266`: `useOptionalEventsContext` y `instanceRef ?? vmid`; seguimiento de mocks internos en `FIX-141` de [`futuro2.md`](futuro2.md)).
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
> - `FIX-138` (resuelto y con tests en verde el 03/11/2026 en `3fe055b`: ampliación de `CaracteresEspeciales` NIST SP 800-63B en `password.go`).
> - `FIX-140` (resuelto y con tests en verde el 03/11/2026 en `4e4d2c4`: `emailProvider` en `/api/version` y cabecera `X-Centinela-Email-Provider: mock`).
> - `FIX-143` (resuelto y con tests en verde el 04/11/2026 en `4011b08`: vocabulario cerrado de `auditoria.accion` y `auditoria.resultado`, enums publicados en Swagger y `contrato-etapa1.md`, test `TestContratoVocabularioAuditoria` y migración histórica `20261104_vocabulario_auditoria.sql`; `TestFIX143_VocabularioAuditoriaContratoCerrado` 4/4 en verde).
>
> **Pasados a `terminado3.md` (Etapa Adaptación Proxmox PVE 8.4):**
> - `ADP-01`, `ADP-02`, `ADP-03`, `ADP-04`, `ADP-05` y `ADP-08` (resueltos y con tests en verde el 01/11/2026).
> - `ADP-06` (`SIM-PVE8-02`, resuelto y con tests en verde el 03/11/2026 en `aee232a`: catálogo de plantillas e ISOs Debian 12 Bookworm en `storage.go`).
>
> **Seguimiento de Fixes en `futuro2.md` y `futuro3.md`:**
> - `FIX-141` (en [`futuro2.md`](futuro2.md): sincronizar tests unitarios internos de `frontend/centinela/src/test/` con los cambios recientes de mocks y UI).
> - `FIX-144` (Frontend, en [`futuro3.md`](futuro3.md) y en este tablero: Alineación 1 a 1 de `auditActionLabels.ts`, `Auditoria.tsx` y `Usuarios.tsx` con el enum cerrado de `auditoria.accion` y `auditoria.resultado` del Backend y retiro de lectura doble).

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


### `FIX-141` - Actualización de mocks de `@/context/EventsContext` (`useOptionalEventsContext`) y sincronización de expectativas UI en tests unitarios internos del submódulo `frontend` (`FIX-127` / `99c7a18` / `b93654b`) (Testing Frontend)

- **Área:** Frontend (Tests internos del submódulo)
- **Asignado:** Cristian / Luz
- **Estimación:** 0.5 h
- **Depende de:** `FIX-127`, `FRN-29`.
- **Nota de diagnóstico:** La implementación de producto de `FIX-127` en `EventsContext.tsx`, `SnapshotsTabContent.tsx` e `InstanceDetailPage.tsx` **está bien implementada** y pasa las pruebas de aceptación (`test/front/etapa3-pendientes.test.tsx` 3/3 ✅ y `test/front/instance-detail.test.tsx` 3/3 ✅). El fallo ocurre únicamente en archivos de test internos del submódulo `frontend/centinela/src/test/` que quedaron desincronizados tras cambios recientes:
- **Tests internos del submódulo que fallan actualmente:**
  1. `frontend/centinela/src/test/snapshots.test.tsx` (4/4 fallan con `Error: [vitest] No "useOptionalEventsContext" export is defined on the "@/context/EventsContext" mock`)
  2. `frontend/centinela/src/test/instance-transient-states.test.tsx` (1/3 falla con `Error: [vitest] No "useOptionalEventsContext" export is defined on the "@/context/EventsContext" mock`)
  3. `frontend/centinela/src/test/instances-table.test.tsx` (3 tests fallan tras el commit `99c7a18` que eliminó las columnas `Checkbox` e `ID` de `Instances.tsx` sin actualizar ese test interno).
  4. `frontend/centinela/src/test/dashboard-node-status.test.tsx` (3 tests fallan tras renombrar los botones `"Ver todas"` y `"Ver todo"` a `"Ver todas las instancias"` y `"Ver toda la actividad"` en el Dashboard).
  5. `frontend/centinela/src/test/accessibility.test.tsx` (1 test falla porque `button.tsx` reemplazó `focus-visible:ring-4` por `focus-visible:ring-0 focus-visible:bg-[#93c5fdbd]`).
- **Solución técnica paso a paso:**
  1. En `frontend/centinela/src/test/snapshots.test.tsx:15-17` y `frontend/centinela/src/test/instance-transient-states.test.tsx:15-18`, agregar `useOptionalEventsContext` al objeto devuelto por `vi.mock('@/context/EventsContext')`:
     ```ts
     vi.mock('@/context/EventsContext', () => ({
         useEventsContext: () => ({ ultimoMensaje: null }),
         useOptionalEventsContext: () => ({ ultimoMensaje: null }),
         useEventSubscription: () => {},
     }));
     ```
  2. Sincronizar las aserciones de `instances-table.test.tsx`, `dashboard-node-status.test.tsx` y `accessibility.test.tsx` con los cambios visuales introducidos en `99c7a18` y `b93654b`.
- **Criterio de éxito:**
  - `npm --prefix frontend/centinela test -- --run` pasa los 43 archivos de test al 100% en verde.

---
# ETAPA 3: Tareas Pendientes de Implementación
---

*Las tareas `BAC-44`, `BAC-45`, `BAC-48`, `BAC-54`, `BAC-55`, `FRN-29`, `FRN-30`, `FRN-32`, `FRN-34`, `FIX-128` (`FRN-31`), `FIX-131`, `FIX-132`, `FIX-133`, `FIX-134`, `FIX-135`, `FIX-136`, `FIX-138` (Backend) y `FIX-140` fueron completadas y verificadas en [`terminado3.md`](terminado3.md). Restan las siguientes tareas pendientes de la Etapa 3:*

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

---
---

### `FIX-144` - Alineación de filtros, contadores y diccionario de acciones de auditoría con el enum cerrado del Backend y retiro de lectura doble (`RF-08` / `FRN-31` / `FRN-12`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1.0 h
- **Depende de:** `FIX-143`, `FRN-31`, `FRN-12`.
- **Objetivo:**
  Sincronizar 1 a 1 el diccionario de acciones (`auditActionLabels.ts`), los contadores/filtros de resultado de la vista de Auditoría (`Auditoria.tsx`) y la consulta del historial de bajas de usuarios (`Usuarios.tsx`) con el único enum cerrado de `auditoria.accion` y `auditoria.resultado` publicado por el Backend en `FIX-143`, retirando los códigos fantasma y la lectura doble temporal.

- **Contexto y Problemática Técnica Detallada (Frontend):**
  Debido a la coexistencia de dos familias de códigos en el backend antes de `FIX-143`, el frontend presenta hoy tres defectos concretos de integración:

  1. **Contradicción y pérdida de registros al filtrar por `resultado` entre `Auditoria.tsx` y `Usuarios.tsx`:**
     - **En [`Auditoria.tsx:220-237`](file:///home/micaela/centinela/frontend/centinela/src/pages/Auditoria.tsx#L220-L237) (`loadResultCounts`) y líneas [`425-438`](file:///home/micaela/centinela/frontend/centinela/src/pages/Auditoria.tsx#L425-L438):** Las pestañas rápidas `"Exitosos"` / `"Fallidos"`, los contadores superiores y el `<select aria-label="Resultado">` envían `resultado=SUCCESS` y `resultado=FAILED`. Como PostgreSQL filtra por igualdad exacta (`WHERE a.resultado = ?`), al filtrar por `"Exitosos"` **no se cuentan ni listan los 8.041 registros de `LOGIN`, `LOGOUT`, `2FA`, `CREAR_USUARIO`, `ELIMINAR_USUARIO` y `ASIGNAR_PERMISOS`** (guardados como `EXITO`), y al filtrar por `"Fallidos"` **no aparece ningún intento fallido de login o 2FA** (guardados como `FALLA`).
     - **En [`Usuarios.tsx:244-246`](file:///home/micaela/centinela/frontend/centinela/src/pages/Usuarios.tsx#L244-L246):** El historial de usuarios eliminados consulta `getAuditLogs({ accion: 'ELIMINAR_USUARIO', resultado: 'EXITO' })`, usando la familia opuesta a `Auditoria.tsx`.

  2. **Códigos fantasma, duplicados y faltantes en `AUDIT_ACTION_ORDER` ([`auditActionLabels.ts:55-86`](file:///home/micaela/centinela/frontend/centinela/src/utils/auditActionLabels.ts#L55-L86)):**
     - El selector de filtro de acciones en el frontend ofrece **4 códigos que no existen en `ports/audit_port.go` del backend**:
       - `CAMBIAR_CONTRASENA`: el backend emite `PASSWORD_CHANGED` (`audit_port.go:106`), por lo que filtrar por `"Contraseña actualizada"` en la UI envía `accion=CAMBIAR_CONTRASENA` y devuelve **0 resultados**.
       - `START_VM` y `STOP_VM`: conviven en `AUDIT_ACTION_ORDER` junto con `INICIAR_VM` y `DETENER_VM`, duplicando en el desplegable las opciones `"Instancia iniciada"` e `"Instancia detenida"` con códigos que el backend no escribe.
       - `EXPORTAR_AUDITORIA`, `SUSPENDER_VM`, `REANUDAR_VM` y `ACTUALIZAR_PERFIL`: se ofrecen en el filtro aunque ningún servicio del backend los persiste.
     - A la inversa, `AUDIT_ACTION_ORDER` **omite 3 acciones reales emitidas por el backend**: `UPDATE_PROFILE_EMAIL_REQUESTED`, `UPDATE_PROFILE_EMAIL_CONFIRMED` e `INSTANCE_PERMISSIONS_PURGED` (que sí están en `AUDIT_ACTION_LABELS` en las líneas 18-20 pero fueron omitidas de `AUDIT_ACTION_ORDER`).

- **Especificación Técnica (Qué debe hacer el equipo de Frontend):**
  1. **En [`frontend/centinela/src/utils/auditActionLabels.ts`](file:///home/micaela/centinela/frontend/centinela/src/utils/auditActionLabels.ts):**
     - Alinear `AUDIT_ACTION_ORDER` y `AUDIT_ACTION_LABELS` **1 a 1** con el enum cerrado de `auditoria.accion` publicado por el backend en `FIX-143`:
       - Reemplazar `CAMBIAR_CONTRASENA` por el código canónico del backend (`PASSWORD_CHANGED`, o el código unificado en `FIX-143`).
       - Eliminar de `AUDIT_ACTION_ORDER` los códigos duplicados/inexistentes (`START_VM`, `STOP_VM`, `EXPORTAR_AUDITORIA`, `SUSPENDER_VM`, `REANUDAR_VM`, `ACTUALIZAR_PERFIL`).
       - Incorporar en `AUDIT_ACTION_ORDER` las acciones activas faltantes (`INSTANCE_PERMISSIONS_PURGED`, `UPDATE_PROFILE_EMAIL_REQUESTED`, `UPDATE_PROFILE_EMAIL_CONFIRMED`).
       - Una vez ejecutada la normalización histórica de `FIX-143`, retirar los alias temporales (`START`, `STOP`, `SHUTDOWN`, `REBOOT`, `CREAR_INSTANCIA`, `CREAR_VM`, `CREAR_LXC`, `CONFIG_INSTANCE`) señalados en el comentario de la línea 26.
  2. **En [`frontend/centinela/src/pages/Auditoria.tsx`](file:///home/micaela/centinela/frontend/centinela/src/pages/Auditoria.tsx) y [`frontend/centinela/src/pages/Usuarios.tsx`](file:///home/micaela/centinela/frontend/centinela/src/pages/Usuarios.tsx):**
     - Actualizar `ResultFilter`, `loadResultCounts` (`Auditoria.tsx:220-237`), las pestañas de filtro rápido (`Auditoria.tsx:425-428`), las opciones del `<select aria-label="Resultado">` (`Auditoria.tsx:436-438`), el badge de estado (`resultLabel`) y la consulta de historial de bajas en `Usuarios.tsx:244-246` para que utilicen **exclusivamente el único enum cerrado de `auditoria.resultado`** definido en `FIX-143` (sin mezclar `SUCCESS` en una vista y `EXITO` en otra).

- **Criterios de Aceptación (DoD — Frontend):**
  1. `AUDIT_ACTION_OPTIONS` en `auditActionLabels.ts` coincide 1 a 1 con el enum cerrado de `auditoria.accion` del backend, sin opciones repetidas ni códigos inexistentes.
  2. `Auditoria.tsx` (contadores, pestañas y selector de resultado) y `Usuarios.tsx` (historial de usuarios eliminados) consultan exactamente los códigos del enum cerrado de `auditoria.resultado`.
  3. La suite `FIX-144` en [`test/front/etapa3-pendientes.test.tsx`](file:///home/micaela/centinela/test/front/etapa3-pendientes.test.tsx#L261-L293) pasa íntegramente en verde.





