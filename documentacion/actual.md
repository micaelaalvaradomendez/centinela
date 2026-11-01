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

#### `ADP-01` (`BAC-PVE8-01`) - Ajuste de Parámetros de Aprovisionamiento QEMU y LXC en `client.go`
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** Cierre de Etapa 3 (`INT-07`).
- **Problema:** `client.go` no incluye `ostype`, `agent=1`, `unprivileged=1` ni `nesting=1` al llamar a `POST /nodes/{node}/qemu` y `/lxc`, provocando que en PVE 8.4 las VMs no activen el socket de red del agente y los contenedores no arranquen con systemd.
- **Entregables:**
  1. En `backend/internal/adapters/secondary/proxmox/client.go:parametrosCreacion`:
     - Agregar `params.Set("ostype", "l26")` y `params.Set("agent", "1")` para QEMU.
     - Agregar `params.Set("unprivileged", "1")` y `params.Set("features", "nesting=1")` para LXC.
  2. Actualizar pruebas unitarias en `client_test.go` verificando que los parámetros urlencoded incluyan estas directivas.
- **Criterio de éxito:** `go test -v -run TestClient_AprovisionarInstancia ./internal/adapters/secondary/proxmox` pasa al 100% validando la inclusión de `agent=1`, `ostype=l26`, `unprivileged=1` y `nesting=1`.

#### `ADP-02` (`BAC-PVE8-02`) - Parámetros de Purga Profunda en Eliminación de Instancias
- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `ADP-01`.
- **Problema:** En PVE 8.4, eliminar una instancia mediante `DELETE /nodes/{node}/{tipo}/{vmid}` sin `purge=1` falla si la instancia tiene metadatos en HA o deja discos huérfanos sin desasociar.
- **Entregables:**
  1. En `client.go:EliminarInstancia`:
     - Agregar query string `?purge=1&destroy-unreferenced-disks=1` a la llamada HTTP DELETE.
  2. Agregar prueba unitaria en `client_test.go` asegurando que la URL de borrado contenga los parámetros de purga.
- **Criterio de éxito:** La llamada de borrado emite la query esperada y la prueba unitaria pasa en verde.

#### `ADP-03` (`BAC-PVE8-03`) - Flexibilización y Tolerancia en Parser de Interfaces y Configuración
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `ADP-01`.
- **Problema:** Si PVE 8.4 entrega `memory` como entero en `qemu config` o `prefix` como entero en `lxc interfaces`, no debe ocurrir ningún desajuste en el cómputo de digests de concurrencia ni en la extracción de IPs.
- **Entregables:**
  1. Verificar que `ObtenerConfig` mantenga consistencia en el digest tanto si `memory` llega como string o como entero en el JSON.
  2. Asegurar que `ObtenerInterfaces` gestione respuestas donde `ip-addresses` sea `null`, esté ausente o use CIDRs tradicionales en `inet`.
  3. Tests de regresión en `client_test.go` con payloads de PVE 8.4 reales.
- **Criterio de éxito:** Suites `TestClient_ObtenerInterfaces_LXC` y `TestObtenerConfig_QemuYLXC` pasan sin errores contra payloads mixtos (string e integer).

#### `ADP-04` (`BAC-PVE8-04`) - Parametrización Dinámica de VMIDs Protegidos por Entorno
- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** Ninguna.
- **Problema:** El rango `100-105` está hardcodeado por defecto para la infraestructura de Centinela. En el servidor final con PVE 8.4.21, Centinela podría usar otros IDs y no debe bloquear VMs preexistentes del hipervisor.
- **Entregables:**
  1. En `backend/internal/adapters/primary/http/middleware/protected_instance.go`:
     - Verificar que `PROXMOX_PROTECTED_VMIDS` permita listas vacías o personalizadas mediante `.env` sin romper el arranque de la API.
  2. Actualizar `SETUP.md` y `.env.example` documentando la configuración de VMIDs protegidos para el servidor productivo.
- **Criterio de éxito:** Si `PROXMOX_PROTECTED_VMIDS=""` o contiene IDs arbitrarios (ej. `"200,201"`), la API protege exactamente esos IDs y permite gestionar los IDs 100-105 sin retornar 403.


#### `ADP-05` (`SIM-PVE8-01`) - Soporte de Modo Dual de Versión en Simulador (`PROXMOX_SIM_TARGET_VERSION`)
- **Área:** Backend / Testing
- **Asignado:** Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `ADP-01`, `ADP-03`.
- **Problema:** El simulador reporta rígidamente `pve-manager/9.2.2` y kernel 7.0.2 en `/nodes/{node}/status`, y no valida los parámetros de `agent` ni `ostype` al crear VMs.
- **Entregables:**
  1. En `backend/cmd/proxmox-simulador/simulador.go` y `config.go`:
     - Incorporar variable de configuración `TargetVersion` (default `"8.4.21"`).
     - Cuando `TargetVersion == "8.4.21"`:
       - `/nodes/{node}/status` responde `pveversion: "pve-manager/8.4-21/12345678"` y `kversion: "Linux 6.8.8-2-pve"`.
       - `POST /nodes/{node}/qemu` requiere `ostype` y `agent` para habilitar el socket del guest agent.
       - `DELETE` exige `purge=1` si la instancia tiene el flag HA activado.
- **Criterio de éxito:** `PROXMOX_SIM_TARGET_VERSION=8.4.21 go run ./cmd/proxmox-simulador` expone las firmas y validaciones de PVE 8.4.21.

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

#### `ADP-08` (`FRN-PVE8-02`) - Verificación de Renderizado de Estado de Storage y Nodos en Dashboard
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1.5 h
- **Depende de:** `ADP-07`.
- **Problema:** En PVE 8.4, los storages reportan métricas sin el content type `import` y con ligeras variaciones de redondeo en bytes disponibles.
- **Entregables:**
  1. En `Dashboard.tsx` y componentes de medición (`StorageMeter.tsx`, `InstanceResourcesPanel.tsx`):
     - Validar que los cálculos de porcentaje de uso no produzcan NaN ni desbordamientos ante los tamaños de bloque devueltos por PVE 8.4.
- **Criterio de éxito:** Tarjetas de almacenamiento y nodo renderizan correctamente en el Dashboard con datos simulados de PVE 8.4.

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
