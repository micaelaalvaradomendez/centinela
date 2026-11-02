### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (30/10/2026):** según ClickUp, en este archivo quedan las tareas **pendientes de implementación** de la Etapa 2, Etapa 3 y etapa adaptacion-proxmox. Lo ya verificado con las suites está en [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md) / [`terminado3.md`](terminado3.md), y las correcciones de las implementadas con problemas, en [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md) / [`futuro3.md`](futuro3.md).
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
> - `FRN-22`, `FRN-23`, `FRN-26`, `FRN-27` y `FRN-28` (resueltos y con tests en verde el 26/10/2026: frontend `41cad44`; telemetría en tiempo real WS con 12 tests, curvas SVG de recursos, seguimiento de cuota y edición reactiva de vCPU/RAM).
> - `FIX-91` (resuelto y con tests en verde el 30/10/2026: eliminación de código muerto `InstanceTester.tsx`; `fix-91-dead-code.test.tsx` 2/2 en verde).
> - `FIX-108` y `FIX-124` (resueltos y con tests en verde el 30/10/2026: identificadores opacos `InstanceRef` y sanitización estricta anti-traversal; `fix-108-ids-ofuscados.test.tsx` 11/11 y `ids-ofuscados.test.tsx` 14/14 en verde).
> - `FRN-36` (resuelto y con tests en verde el 30/10/2026: modularización de `InstanceDetailPage`, tabs transitorios y bloqueo de cabecera; `frn-36-29-snapshots-ui.test.tsx` en verde).
>
> **Pasados a `terminado3.md` (Etapa 3):**
> - `BAC-42`, `BAC-43`, `BAC-46`, `BAC-47` (resueltos en cierres previos).
> - `BAC-44` (resuelto y con tests en verde el 30/10/2026: interceptor global, filtros compuestos y orden `fecha_hora DESC`).
> - `BAC-45` (resuelto y con tests en verde el 30/10/2026: streaming HTTP con cursor `.Rows()`, sanitización CSV y cabecera `auditoria_centinela_YYYYMMDD`).
> - `BAC-48` (resuelto y con tests en verde el 30/10/2026: cookies `HttpOnly`, `SameSite=Strict`, `Secure` y Rate Limiting 429).
> - `BAC-54` (resuelto y con tests en verde el 30/10/2026: lectura dinámica de `PROXMOX_NODE`, `PROXMOX_URL` y `PROXMOX_PROTECTED_VMIDS`).
> - `FRN-29` (resuelto y con tests en verde el 30/10/2026: pestaña Snapshots, validación de `snapname` y feedback reactivo; `frn-36-29-snapshots-ui.test.tsx` y `snapshots.test.tsx` en verde).
> - `FRN-30` (resuelto y con tests en verde el 30/10/2026: modal destructivo de rollback con confirmación antierror tipear nombre idéntico; `frn-30-snapshots-rollback-ui.test.tsx` 2/2 en verde).
> - `FRN-34` (resuelto y con tests en verde el 30/10/2026: adaptación de cliente a cookies seguras de sesión y auditoría WCAG AA; `frn-34-cookies-a11y.test.tsx` 3/3 en verde).
>
> **Pasados a `terminado3.md` (Etapa Adaptación Proxmox PVE 8.4):**
> - `ADP-01`, `ADP-02`, `ADP-03`, `ADP-04`, `ADP-05` y `ADP-08` (resueltos y con tests en verde el 01/11/2026: backend `ec5922c` y frontend `ce38ea9`).
>
> **Seguimiento de Fixes en `futuro2.md` y `futuro3.md`:**
> - `FIX-126` (seguimiento en [`futuro2.md`](futuro2.md): ocultar `<StorageMeter>` en Dashboard para `OPERATOR` o `data.storage == null`).
> - `FIX-127` (seguimiento en [`futuro2.md`](futuro2.md): resiliencia de `SnapshotsTabContent` ante providers desacoplados y alias `:vmid`/`:instanceRef`).
> - `FIX-128` (seguimiento en [`futuro3.md`](futuro3.md): acciones `SNAPSHOT_DELETE`, `ALERT_TRIGGERED` y exportación JSON en `Auditoria.tsx`).
> - `FIX-129` (seguimiento en [`futuro3.md`](futuro3.md): cabeceras `HSTS includeSubDomains`, `X-Frame-Options`, `X-Content-Type-Options` y `CSP` en `docker/nginx-edge.conf`).
> - `FIX-130` (seguimiento en [`futuro3.md`](futuro3.md): documentación de compatibilidad de storages con snapshots y placeholders seguros en `documentacion/api-proxmox.md`).
> - `FIX-131` (seguimiento en [`futuro3.md`](futuro3.md): métricas y gráfico de Tráfico de Red `netin`/`netout` en detalle de instancia para `CU-05`).
> - `FIX-132` (seguimiento en [`futuro3.md`](futuro3.md): discriminación de errores 400 en recuperación de contraseña para no quemar OTP ante `SAME_PASSWORD` para `CU-12`).

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


### `FIX-127` - Resiliencia de `SnapshotsTabContent` ante Providers desacoplados y alias de ruta `:vmid` / `:instanceRef` (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Luz
- **Estimación:** 1.0 h
- **Depende de:** `FRN-36`, `FRN-29`.
- **Problema técnico detallado:**
  1. En `frontend/centinela/src/components/features/instances/components/detail/SnapshotsTabContent.tsx:69`, se ejecuta `const { ultimoMensaje } = useEventsContext();` sin comprobar si el componente está envuelto en un `EventsProvider`. Cuando se monta en pruebas unitarias aisladas (`src/test/instance-transient-states.test.tsx`), lanza `Error: useEventsContext debe utilizarse dentro de un EventsProvider` y falla el renderizado del estado vacío.
  2. En `frontend/centinela/src/pages/InstanceDetailPage.tsx:12`, la lectura de parámetros utiliza estrictamente `const { instanceRef } = useParams<{ instanceRef: string }>()`. Si componentes, tests legados o enlaces invocan la ruta con el parámetro `:vmid`, `instanceRef` queda como `undefined` y rechaza la navegación como "Identificador inválido".
- **Qué debe hacer el equipo de frontend:**
  1. En `EventsContext.tsx` o `SnapshotsTabContent.tsx`, permitir que `useEventsContext()` devuelva `null` o un valor seguro por defecto cuando no esté presente el `EventsProvider` (o usar un hook defensivo `useOptionalEventsContext()`).
  2. En `InstanceDetailPage.tsx`, leer ambos parámetros con fallback: `const { instanceRef, vmid } = useParams(); const ref = instanceRef ?? vmid;`.
  3. En `applicationRoutes.tsx`, asegurar que ambas variantes de ruta (`/instances/:instanceRef` y `/instances/:vmid`) dirijan al mismo componente.
- **Criterio de éxito:**
  La prueba unitaria `src/test/instance-transient-states.test.tsx` y `src/test/instance-detail.test.tsx` pasan al 100% sin requerir providers globales en tests unitarios.

---

### `FIX-128` - Admisión de `vlan: 0` o su omisión como tráfico no etiquetado (*untagged*) en la creación de instancias (Backend)

- **Área:** Backend
- **Asignado:** Lisandro / Tayra
- **Estimación:** 0.5 h
- **Depende de:** `BAC-37`, `RF-07`
- **Problema técnico detallado:**
  En `backend/internal/adapters/primary/http/instance_handler.go:427-429`, la función `validarCreacionNormalizada` aplica la siguiente regla de validación sobre el identificador de red VLAN:
  ```go
  if input.VLAN != nil && (*input.VLAN < 1 || *input.VLAN > 4094) {
      return ports.CrearInstanciaProxmoxInput{}, fmt.Errorf("vlan debe estar entre 1 y 4094")
  }
  ```
  Al enviar peticiones desde el frontend o clientes API donde no se configura una VLAN explícita (tráfico nativo *untagged* por defecto), los serializadores JSON suelen enviar `"vlan": 0` (o `null`). Debido a que el DTO `CrearInstanciaRedRequest.VLAN` se bindea como un puntero a entero (`*int`), un valor numérico `0` hace que `input.VLAN != nil` sea verdadero y que `*input.VLAN < 1` falle la condición. Esto provoca una respuesta inmediata `400 Bad Request: vlan debe estar entre 1 y 4094`, bloqueando la creación de cualquier instancia (VM o LXC).

  En redes basadas en el estándar IEEE 802.1Q y en Proxmox VE, el valor `0` representa la ausencia de etiqueta VLAN (*tráfico nativo / sin etiquetar*). En el hipervisor Proxmox VE, los parámetros de red únicamente llevan la directiva `,tag=<id>` cuando la VLAN es efectiva ($\ge 1$).
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/adapters/primary/http/instance_handler.go`, dentro de `validarCreacionNormalizada`, admitir explícitamente `0` transformándolo en `nil` (o saneando el puntero antes de validar):
     ```go
     if input.VLAN != nil {
         if *input.VLAN == 0 {
             input.VLAN = nil
         } else if *input.VLAN < 1 || *input.VLAN > 4094 {
             return ports.CrearInstanciaProxmoxInput{}, fmt.Errorf("vlan debe estar entre 1 y 4094")
         }
     }
     ```
  2. Verificar que en `internal/adapters/secondary/proxmox/client.go` (`AprovisionarInstancia`), la directiva de red `net0` no incluya `,tag=0` cuando el valor sea `nil` o `0` (el cliente actual ya valida `if input.VLAN != nil { net0 += ",tag=" + strconv.Itoa(*input.VLAN) }`).
  3. Añadir pruebas unitarias en `instance_handler_test.go` verificando que un payload de creación con `"vlan": 0`, `"vlan": null` o con el campo omitido sea aceptado exitosamente con respuesta `202 Accepted` y despachado sin error.
- **Criterio de éxito:**
  - `POST /api/instances` acepta peticiones con `"vlan": 0`, interpretándolas como interfaz de red nativa sin etiquetar (*untagged*).
  - No se rechaza con `400 INVALID_REQUEST` ante la ausencia o valor neutro (`0`) de VLAN.
  - La suite de pruebas de handlers y creación de instancias pasa al 100% en verde.




---
# ETAPA 3: Tareas Pendientes de Implementación
---

*Las tareas `BAC-44`, `BAC-45`, `BAC-48`, `BAC-54`, `FRN-29`, `FRN-30` y `FRN-34` fueron completadas y verificadas en [`terminado3.md`](terminado3.md). Los ajustes de `INF-11` (`FIX-129`) y `FRN-31` (`FIX-128`) se encuentran en [`futuro3.md`](futuro3.md). Restan las siguientes tareas pendientes de la Etapa 3:*

#### `FRN-32` - Componente y flujo de exportación de reportes de auditoría
- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 1.5 h
- **Depende de:** `FRN-31`, `BAC-45`.
- **Problema:** Los operadores de cumplimiento requieren descargar los registros en archivos locales para archivo legal o análisis externo.
- **Entregable:**
  1. Botón "Exportar registros" con menú desplegable para elegir "Exportar como CSV" o "Exportar como JSON".
  2. Integración con `auditService.exportLogs(filtros, format)` transfiriendo los filtros actualmente activos en pantalla.
  3. Estado de carga en el botón durante la generación del reporte y disparo automático de la descarga nativa en el navegador.
  4. Notificación Toast al completar la descarga.
- **Criterio de éxito:** Al hacer clic en exportar, se descarga el archivo con los filtros aplicados sin recargar la página ni bloquear la UI.


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


#### `BAC-55` (`BRG-15-SIM`) - Switch de Perfil de Versión en `proxmox-simulador` (`PROXMOX_SIM_TARGET_VERSION`)
- **Área:** Backend / Simulador
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-54`.
- **Problema:** El simulador responde únicamente con firmas de PVE 9.2.2. Se requiere que admita perfiles de versión mediante variable de entorno para poder contrastar pruebas tanto en versión 9 como en versión 8.4.21.
- **Entregable:**
  1. En `backend/cmd/proxmox-simulador/simulador.go`:
     - Incorporar lectura de `PROXMOX_SIM_TARGET_VERSION` (valores: `"9.2.2"`, `"8.4.21"`; default `"9.2.2"` para no romper Etapa 3).
  2. Implementación de respuestas condicionadas:
     - En `nodes/{node}/status`: exponer firma `pve-manager/8.4-21` y kernel `6.8` cuando la variable esté en modo 8.4.21.
     - Hook de preparación para que las tareas de adaptación (`ADP-05`, `ADP-06`) conecten sus validaciones específicas.
- **Criterio de éxito:** Ejecutar `PROXMOX_SIM_TARGET_VERSION=8.4.21 go run ./cmd/proxmox-simulador` emite la cabecera y el JSON de status del nodo correspondientes a la versión 8.4.21.
     - En `nodes/{node}/status`: exponer firma `pve-manager/8.4-21` y kernel `6.8` cuando la variable esté en modo 8.4.21.
     - Hook de preparación para que las tareas de adaptación (`ADP-05`, `ADP-06`) conecten sus validaciones específicas.
- **Criterio de éxito:** Ejecutar `PROXMOX_SIM_TARGET_VERSION=8.4.21 go run ./cmd/proxmox-simulador` emite la cabecera y el JSON de status del nodo correspondientes a la versión 8.4.21.

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

*Las tareas `ADP-01`, `ADP-02`, `ADP-03`, `ADP-04`, `ADP-05` y `ADP-08` fueron completadas y verificadas en [`terminado3.md`](terminado3.md). Restan las siguientes tareas pendientes de adaptación:*

#### `ADP-06` (`SIM-PVE8-02`) - Actualización del Catálogo Inicial de Plantillas e ISOs a Debian 12 Bookworm
- **Área:** Backend / Simulador
- **Asignado:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `ADP-05`.
- **Problema:** El inventario de plantillas de `storage.go` expone `debian-13-standard_13.6-1_amd64.tar.zst` y `debian-13.6.0-amd64-netinst.iso`, inexistentes en repositorios oficiales de PVE 8.4.
- **Entregables:**
  1. En `backend/cmd/proxmox-simulador/storage.go`:
     - Reemplazar ítems del catálogo de `local:vztmpl` por `debian-12-standard_12.7-1_amd64.tar.zst` y `ubuntu-22.04-standard_22.04-1_amd64.tar.zst`.
     - Reemplazar ISO en `local:iso` por `debian-12.7.0-amd64-netinst.iso`.
     - Remover el content type `import` en el almacenamiento `local` para coincidir con PVE 8.4.
  2. Ajustar pruebas unitarias del simulador en `simulador_test.go`.
- **Criterio de éxito:** `GET /nodes/{node}/storage/local/content` devuelve las plantillas de Debian 12 y `simulador_test.go` pasa al 100%.


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


### `FIX-128` - Opciones completas de acciones de auditoría Etapa 3 y selectores en Auditoria.tsx (`FRN-31`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-31`, `BAC-44`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/utils/auditActionLabels.ts`, el listado de acciones en `AUDIT_ACTION_LABELS` y `AUDIT_ACTION_ORDER` omite las acciones `SNAPSHOT_DELETE` ("Snapshot eliminado") y `ALERT_TRIGGERED` ("Alerta disparada"). Asimismo, `SNAPSHOT_ROLLBACK` está rotulado como `"Snapshot restaurado"`, lo que no coincide con las expresiones esperadas de reversión/rollback (`/revertir.*snapshot|rollback|snapshot_rollback/i`). Por último, `Auditoria.tsx` renderiza un único botón estático de `"Exportar CSV"` en lugar de ofrecer la opción de descarga tanto en CSV como en JSON integrada con `BAC-45`.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/utils/auditActionLabels.ts`:
     - Incorporar en `AUDIT_ACTION_LABELS` y `AUDIT_ACTION_ORDER`:
       - `SNAPSHOT_DELETE: 'Snapshot eliminado'`
       - `ALERT_TRIGGERED: 'Alerta disparada'`
     - Asegurar que `SNAPSHOT_ROLLBACK` sea etiquetado como `"Revertir snapshot"` o `"Snapshot revertido"` para satisfacer la búsqueda por vocabulario canónico.
  2. En `frontend/centinela/src/pages/Auditoria.tsx`:
     - Incorporar en el filtro de tipo de recurso la opción `SNAPSHOT`.
     - Añadir botón o menú de exportación que permita seleccionar entre `"Exportar como CSV"` y `"Exportar como JSON"`, disparando la descarga de `auditoria_centinela_YYYYMMDD.json` con `format=json`.
- **Criterio de éxito:**
  La suite `test/front/frn-31-audit-etapa3.test.tsx` pasa 4 de 4 pruebas en verde (100% ✅).

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

### `FIX-132` - Discriminación de errores 400 en recuperación de contraseña para no quemar intentos de OTP ante contraseñas inválidas o duplicadas (`CU-12`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz / Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-34`.
- **Problema técnico detallado:**
  En el flujo de excepciones de **CU-12** (`casosdeuso.md:612`), se establece taxativamente:
  > *"Contraseña débil o igual a la anterior (`400 WEAK_PASSWORD` / `400 SAME_PASSWORD`): El backend rechaza el cambio sin quemar el código OTP para que el usuario pueda corregir la contraseña elegida."*
  
  Sin embargo, en `frontend/centinela/src/pages/RecoverPassword.tsx:130-135`:
  ```tsx
  if (error.errorCode === 'RESET_FAILED' || error.status === 400) {
      const nextAttempts = failedAttempts + 1;
      ...
      setCodeError('El código ingresado es incorrecto o ha expirado.');
  ```
  La condición `error.status === 400` captura de forma indiscriminada cualquier error 400 devuelto por `POST /api/auth/password/reset` (incluyendo `400 SAME_PASSWORD`, `400 WEAK_PASSWORD` o errores de validación de contraseña). Esto provoca que la interfaz acuse erróneamente al código OTP de ser incorrecto, incremente `failedAttempts` y bloquee el proceso (`setResetBlocked(true)`) al 3.er intento, quemando en el cliente un código OTP que es completamente válido en el servidor.
- **Qué debe hacer el equipo de frontend:**
  En `RecoverPassword.tsx`:
  1. Capturar específicamente los códigos de error de contraseña previa al conteo de intentos:
     ```tsx
     if (error.errorCode === 'SAME_PASSWORD' || error.errorCode === 'WEAK_PASSWORD' || error.errorCode === 'PASSWORD_RECENTLY_USED') {
         setPasswordError(error.message || 'La nueva contraseña no puede ser igual a la anterior ni violar las políticas de seguridad.');
         return;
     }
     ```
  2. Incrementar `failedAttempts` y aplicar el mensaje `'El código ingresado es incorrecto o ha expirado.'` **únicamente** cuando `error.errorCode === 'RESET_FAILED'`.
- **Criterio de éxito:**
  Si el usuario ingresa un código OTP válido pero una contraseña igual a la actual o no conforme a la política del backend, la interfaz muestra el error sobre el campo de contraseña, no incrementa los intentos fallidos del OTP y permite reintentar con otra contraseña sin bloquear el proceso.


### `FIX-133` - Corrección de validación diferencial de cuota de vCPU en escalado vertical para permitir CPU Overcommit (`RF-10` / `BAC-40` / `CU-10`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas / Micaela
- **Estimación:** 1.5 h
- **Depende de:** `BAC-40`, `RF-10`, `CU-10`.
- **Problema técnico detallado:**
  En `documentacion/archvios/requerimientos.md` (RF-10, líneas 266-268), la especificación de negocio `BAC-40` establece textualmente:
  > *"El backend debe aplicar la validación diferencial de cuota (`BAC-40`):*
  > *- Si $\Delta\text{RAM} \le 0$ (reducción o cambio solo de vCPU), autorizar directamente.*
  > *- Si $\Delta\text{RAM} > 0$, verificar que la memoria solicitada quepa en el margen seguro del nodo..."*
  
  Sin embargo, en la implementación de `internal/core/services/configuracion_instancias.go:63-67`:
  ```go
  if deltaCores > 0 || deltaMemoria > 0 {
      err = s.cuotas.AmpliarInstancia(ctx, inst.Nodo, max(0, deltaCores), max(int64(0), deltaMemoria), aplicar)
  } else {
      err = aplicar()
  }
  ```
  Y en `internal/core/services/cuotas_service.go:207, 244-246`:
  ```go
  coresDisponibles = capacidad.Nodo.CPUs - coresAsignados - coresReservados
  ...
  if deltaCores > coresDisponibles {
      return fmt.Errorf("%w: vCPU disponibles insuficientes en el nodo %s", ports.ErrRecursosInsuficientes, nodo)
  }
  ```
  `coresAsignados` acumula la suma de `MaxCpu` de todas las instancias del inventario. En plataformas de virtualización como Proxmox VE (KVM y LXC), la sobreasignación de CPU (*CPU Overcommit*) es una práctica arquitectónica estándar (por ejemplo, asignar 24 vCPUs acumuladas en un nodo físico de 8 o 16 núcleos). Al validar de forma rígida `deltaCores > coresDisponibles`, `coresDisponibles` resulta $\le 0$ en la mayoría de los hipervisores en producción.
  
  Esto provoca un **falso positivo de "Recursos Insuficientes"** que rechaza arbitrariamente cualquier aumento de vCPU con un error `400/422 RECURSOS_INSUFICIENTES`, bloqueando la operación de reconfiguración incluso cuando el usuario no modifica la memoria RAM ($\Delta\text{RAM} \le 0$).
- **Qué debe hacer el equipo de backend:**
  1. En `internal/core/services/configuracion_instancias.go`:
     - Modificar la condición de delegación a `CuotasService` para que el control estricto de margen de capacidad se ejecute **únicamente cuando $\Delta\text{RAM} > 0$**, tal como exige `RF-10`:
     ```go
     if deltaMemoria > 0 {
         err = s.cuotas.AmpliarInstancia(ctx, inst.Nodo, max(0, deltaCores), deltaMemoria, aplicar)
     } else {
         err = aplicar()
     }
     ```
  2. En `internal/core/services/cuotas_service.go`:
     - En `AmpliarInstancia`, verificar únicamente `deltaMemoria > memoriaDisponible`. Para `deltaCores`, validar únicamente límites de sanidad de la instancia (ej. que la VM no supere el máximo de núcleos físicos admisibles por configuración o un límite máximo configurable de overcommit, pero nunca rechazar por la suma lineal de vCPUs del nodo si no se excede la RAM).
  3. Actualizar la suite de pruebas unitarias en `configuracion_instancias_test.go` para verificar que un cambio con `deltaCores > 0` y `deltaMemoria <= 0` aplica exitosamente la llamada a `proxmox.ActualizarConfig` sin ser bloqueado por cuotas de vCPU.
- **Criterio de éxito:**
  Al invocar `PUT /api/instances/:instanceRef/config` solicitando aumentar vCPUs sin alterar la RAM en un nodo donde la suma de vCPUs de las VMs supera los cores físicos del host, la API responde `200 OK` con `applied: true` y no rechaza con `RECURSOS_INSUFICIENTES`.

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

### `FIX-138` - Ampliación de caracteres especiales admitidos en validación de contraseñas conforme a NIST SP 800-63B (`RF-01` / `FIX-115A/B`) (Backend / Criptografía)

- **Área:** Backend / Criptografía
- **Asignado:** Lucas / Micaela
- **Estimación:** 1.0 h
- **Depende de:** `RF-01`, `FIX-115A`, `FIX-115B`.
- **Problema técnico detallado:**
  En `internal/infrastructure/crypto/password.go:131-152`, la función `ValidarComplejidadContrasena` define la lista de caracteres especiales como:
  ```go
  especiales := "!@#$%^&*-_=+"
  ```
  Las recomendaciones de seguridad de **NIST SP 800-63B** (Digital Identity Guidelines), expresamente adoptadas en `requerimientos.md` (`FIX-115A`), establecen que se debe admitir cualquier carácter ASCII imprimible (códigos 32 a 126) y espacios para permitir el uso de frases de contraseña (*passphrases*) y gestores de contraseñas comerciales.
  
  Actualmente, si un usuario ingresa una contraseña de alta entropía generada por Bitwarden, KeePass o 1Password que contiene símbolos habituales como puntos (`.`), comas (`,`), dos puntos (`:`), punto y coma (`;`), tildes (`~`), paréntesis (`()`), corchetes (`[]`), barras (`/`) o signos de interrogación (`?`), la función rechaza la contraseña con el error:
  `"la contraseña debe contener al menos un carácter especial (!@#$%^&*-_=+)"`.
  
  Esto representa un **falso rechazo de complejidad** que frustra a los usuarios e impide utilizar contraseñas seguras estándar.
- **Qué debe hacer el equipo de backend:**
  1. En `internal/infrastructure/crypto/password.go`:
     - Ampliar el conjunto `especiales` para admitir todos los caracteres de puntuación y símbolos estándar de ASCII:
       ```go
       especiales := `!@#$%^&*-_=+.~,;:'"()[]{}<>\/|?`
       ```
     - Permitir también el espacio (`' '`) cuando forma parte de una passphrase de más de 12 caracteres.
  2. En `frontend/centinela/src/utils/validatePasswordComplexity.ts`:
     - Alinear la expresión regular del validador de frontend con el mismo juego de caracteres especiales para garantizar paridad contractual entre cliente y servidor.
  3. Ejecutar los tests unitarios de `password_test.go` verificando que contraseñas generadas con símbolos extendidos son validadas exitosamente.
- **Criterio de éxito:**
  Contraseñas que contengan caracteres como `.`, `,`, `?`, `~`, `()`, o `/` son aceptadas exitosamente en el registro, reseteo (`POST /api/auth/password/reset`) y cambio de contraseña (`PUT /api/account/password`).


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

### `FIX-140` - Trazabilidad y visibilidad de proveedor de email en modo Mock para prevenir falsos éxitos de entrega (`RF-09` / `RF-11` / `RF-13`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas / Micaela
- **Estimación:** 1.0 h
- **Depende de:** `RF-09`, `RF-11`, `RF-13`.
- **Problema técnico detallado:**
  En `cmd/api/main.go:480-494`:
  ```go
  case "":
      if production {
          return nil, errors.New("EMAIL_PROVIDER no configurado...")
      }
      log.Println("[WARN] EMAIL_PROVIDER ausente. Usando MOCK por defecto")
      return email.NewMockEmailService(), nil
  ```
  Cuando un desarrollador o tester levanta el backend en un entorno de pruebas o pre-producción sin haber configurado las credenciales SMTP en `.env`, el sistema entra silenciosamente en modo `MOCK`.
  
  Cuando se prueba la recuperación de contraseña (`POST /api/auth/password/forgot`), la creación de nuevos usuarios (`POST /api/admin/users`) o el cambio de correo electrónico (`PUT /api/account/profile`), el backend retorna exitosamente `200 OK` informando que el correo fue despachado. Sin embargo, el correo nunca sale a la red y queda sepultado en los logs de la consola del contenedor. Esto genera **falsos éxitos de entrega** que confunden al equipo de frontend y a los usuarios de prueba, quienes creen que el servicio de correo o la red fallaron al no recibir el mensaje.
- **Qué debe hacer el equipo de backend:**
  1. En `internal/adapters/primary/http/auth_handler.go`, `user_handler.go` y `account_handler.go`:
     - Cuando el proveedor activo sea `mock`, inyectar una cabecera HTTP de diagnóstico en las respuestas:
       `X-Centinela-Email-Provider: mock`
     - En los logs del servidor, imprimir un banner destacado con el enlace o código generado para facilitar las pruebas manuales sin revisar líneas dispersas.
  2. En `cmd/api/main.go`:
     - En el endpoint `GET /api/version`, incluir el campo `"emailProvider": "smtp" | "mock"` para que el frontend o los tests de integración puedan conocer de antemano el modo de transporte de correo configurado en el backend.
- **Criterio de éxito:**
  El endpoint `GET /api/version` informa el proveedor de email activo, y cualquier respuesta que despache correos en modo de desarrollo incluye la cabecera `X-Centinela-Email-Provider: mock` para alertar inmediatamente al operador.

