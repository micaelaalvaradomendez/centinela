### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (24/10/2026):** según ClickUp, en este archivo quedan las tareas **pendientes de implementación** de la Etapa 1 y Etapa 2. Lo ya verificado con las suites está en [`terminado.md`](terminado.md) / [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md), y las correcciones de las implementadas con problemas, en [`futuro.md`](futuro.md) / [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md).
>
> **Pasados a `terminado-1.md` (Etapa 1):**
> - `FIX-104B` / `FIX-105` y `FIX-102` (resueltos y con tests en verde el 21/10/2026).
> - `FIX-100` (resuelto y con tests en verde el 22/10/2026: guarda en `formatCount` y `formatShare`).
> - `FIX-91` (pendiente de implementación: eliminar el código muerto `InstanceTester.tsx`; sigue en `futuro-1.md`).
> - `FIX-101`, `FIX-103` y `FIX-104` (cerrados el 22/10/2026: expectativas desactualizadas de `test/front`, no defectos del submódulo; pruebas en verde).
> - `BAC-24C` / `FIX-113` (resuelto y con tests en verde el 24/10/2026: purga atómica de permisos al eliminar instancia; pruebas `TestFIX113` y pasos de integración en verde).
>
> **Pasados a `terminado2.md` (Etapa 2):**
> - `BAC-36`, `FIX-107`, `FIX-109`, `FIX-110`, `FIX-111A`, `FIX-112`, `FIX-114`, `FIX-115A`, `FIX-116`, `FIX-117A`, `FIX-118`, `FIX-119` (resueltos y con tests en verde el 21/10/2026).
> - `FIX-111B` y `FIX-115B` (resueltos y con tests en verde el 22/10/2026: modal de reautenticación y claves hasta 128 chars).
> - `FRN-24` y `FRN-25` (cerradas el 22/10/2026: implementación funcional y pruebas alineadas al contrato real).
> - `INF-10` (pendiente de verificación en Proxmox físico; seguimiento en `FIX-123` dentro de `futuro2.md`).
> - `FIX-117B` (resuelto y con tests en verde el 24/10/2026: adaptación de Dashboard para operador; `fix-117b-dashboard-operador.test.tsx` 4/4 en verde).
> - `FIX-108` (implementada el 24/10/2026 con fallas de validación de entrada; seguimiento en `FIX-124` dentro de `futuro2.md`).
> - `FRN-22`, `FRN-23`, `FRN-26`, `FRN-27` y `FRN-28` (resueltos y con tests en verde el 26/10/2026: frontend `41cad44`; telemetría en tiempo real WS con 12 tests, curvas SVG de recursos, seguimiento de cuota y edición reactiva de vCPU/RAM).

---

### ETAPA 1: Tareas Pendientes de Implementación

*Actualmente no hay tareas pendientes de implementación en la Etapa 1. Todas las tareas fueron completadas y verificadas en [`terminado-1.md`](terminado-1.md).*

---

### ETAPA 2: Tareas Pendientes de Implementación

*Las tareas base de la Etapa 2 fueron completadas y verificadas en [`terminado2.md`](terminado2.md). Restan las tareas puente de integración hacia la Etapa 3:*

#### `FRN-36` (`BRG-10-FRN`) - Modularización de `InstanceDetailPage` para Pestaña de Snapshots y Estados Transitorios
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `FRN-21`, `FRN-28`.
- **Problema:** La página de detalle actual tiene pestañas fijas y no propaga el estado de bloqueo a la cabecera cuando una máquina está en transición asíncrona.
- **Entregable:**
  1. Reestructuración de `InstanceDetailPage.tsx` con soporte para pestañas modulares (dejando listo el slot para "Snapshots").
  2. Soporte en el badge de estado para etiquetas transitorias: `SNAPSHOTTING` (azul), `ROLLING_BACK` (amarillo/naranja).
  3. Deshabilitación reactiva de la botonera de acciones rápidas de la cabecera cuando `instance.activeTask != null`.
- **Criterio de éxito:** Si la máquina reporta `activeTask`, los botones de la cabecera se bloquean con spinner y la pestaña de Snapshots cuenta con su contenedor maquetado.


#### `INF-13` (`BRG-11-INF`) - Verificación y Validación de Storage Compatible con Snapshots en Proxmox Real
- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.5 h
- **Depende de:** `INF-10`.
- **Problema:** Si el Proxmox físico (`100.81.49.19`) utiliza discos `raw` sobre LVM estándar, la API de snapshots fallará en la validación final (`INT-09`).
- **Entregable:**
  1. Verificación del almacenamiento configurado en Proxmox (`local-lvm` con thin-pool, `zfs` o `qcow2`).
  2. Si el storage por defecto no soporta snapshots, configurar un pool LVM-thin o storage de directorio para las instancias de prueba del Centinela.
  3. Documentar en `documentacion/api-proxmox.md` los storages aptos para snapshots.
- **Criterio de éxito:** Ejecución manual de `qm snapshot <vmid> test` por CLI en el Proxmox real completa exitosamente sin errores de formato de disco.

### `FIX-91` - Eliminar el componente muerto `InstanceTester` (`FRN-15` / `FRN-16`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FRN-15`, `FRN-16` (en `terminado-1.md`).
- **Problema técnico detallado:**
  `frontend/centinela/src/components/features/instances/components/InstanceTester.tsx` se titula **"Probador de Instancias (Mock Local)"**, define `INITIAL_INSTANCES` con datos falsos (ids 1..3, "servidor-web-produccion", "base-datos-staging", etc.) y simula las transiciones con `setTimeout`. No se importa desde ningún archivo de producción y persiste como código muerto con datos mock.
- **Qué debe hacer el equipo de frontend:** eliminar el archivo `InstanceTester.tsx`.
- **Criterio de éxito:** no queda código de producto con datos mock ni componentes sin importar.


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


### `FIX-124` - Sanitización estricta de identificadores opacos de instancia en Frontend (FIX-108) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1.0 h
- **Depende de:** `FIX-108` (en `terminado2.md`).
- **Problema técnico detallado:**
  En la implementación de `FIX-108` (frontend `d4e99b2`), la función `validateInstanceRef` en `src/components/features/instances/services/instanceService.ts:12-16` solo ejecuta un recorte de espacios en blanco (`ref.trim()`) cuando recibe un `string`, omitiendo validar contra la expresión regular de identificadores seguros `/^[a-zA-Z0-9_-]+$/`. Como consecuencia, permite entradas hostiles como secuencias de path traversal (`../admin/users`), inyección de parámetros (`?force=1`), separadores de ruta (`inst/1`), punto y coma (`inst_1;rm`), espacios internos (`inst abc`) o saltos de línea (`inst_1\nX`). Adicionalmente, en `src/pages/InstanceDetailPage.tsx:18` persiste una conversión residual `Number(instanceReference)` que infringe la restricción de manipular identificadores opacos exclusivamente como cadenas de texto.
- **Qué debe hacer el equipo de frontend:**
  1. En `src/components/features/instances/services/instanceService.ts`:
     - Restringir la validación de referencias de instancia en `validateInstanceRef` contra `/^[a-zA-Z0-9_-]+$/`:
       ```typescript
       const OPAQUE_ID_REGEX = /^[a-zA-Z0-9_-]+$/;
       function validateInstanceRef(ref: unknown): InstanceRef {
         if (typeof ref === 'string') {
           const trimmed = ref.trim();
           if (trimmed && OPAQUE_ID_REGEX.test(trimmed)) return trimmed;
         }
         if (typeof ref === 'number' && Number.isSafeInteger(ref) && ref > 0) return ref;
         throw new ApiRequestError('La instancia no tiene una referencia válida.');
       }
       ```
  2. En `src/pages/InstanceDetailPage.tsx`:
     - Eliminar la conversión residual `Number(instanceReference)` y validar la referencia directamente sin forzar conversiones numéricas con `Number(...)` o `parseInt(...)` sobre parámetros de ruta.
  3. Ejecutar y certificar la suite de pruebas `test/front/fix-108-ids-ofuscados.test.tsx`.
- **Criterio de éxito:**
  - Todas las entradas maliciosas o malformadas son rechazadas localmente en `instanceService.ts` antes de despachar solicitudes de red.
  - La suite `test/front/fix-108-ids-ofuscados.test.tsx` pasa 11 de 11 pruebas en verde (100% ✅).





### `FIX-126` - Renderizado condicional de StorageMeter en Dashboard para OPERATOR (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 0.5 h
- **Depende de:** `FIX-117A`, `FIX-117B`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/pages/Dashboard.tsx:140`, el componente `<StorageMeter>` se renderiza de forma incondicional aun cuando `data?.storage` es `null` (comportamiento canónico de `FIX-117A` para usuarios con rol `OPERATOR`), mostrando un medidor en 0% en lugar de ocultar la tarjeta de capacidad de almacenamiento del host según lo exigido en `FIX-117B`.
- **Qué debe hacer el equipo de frontend:**
  En `Dashboard.tsx`, condicionar el renderizado del bloque de almacenamiento a la presencia efectiva de `data?.storage` o al rol `user.rol === 'ADMIN'`, impidiendo que un operador visualice métricas de disco vacías.
- **Criterio de éxito:**
  El Dashboard no renderiza `<StorageMeter>` para operadores ni ante `data.storage === null`.


### `FIX-108` - Adaptación del Frontend para consumo de identificadores de instancia ofuscados (Frontend)

> [!WARNING]
> **Estado verificado (24/10/2026, frontend `d4e99b2`): implementada con fallas de validación de entrada.** El frontend adoptó el identificador de instancia como `string | number` (`id: string | number`), y los servicios y rutas despachan hacia `/instances/:instanceRef`. Sin embargo, `validateInstanceRef` en `instanceService.ts` solo realiza un recorte de espacios (`ref.trim()`) y no valida contra `/^[a-zA-Z0-9_-]+$/`, permitiendo path traversal, caracteres hostiles y espacios. Corrección y seguimiento abierto en **`FIX-124`** ([futuro2.md](futuro2.md)). La suite `test/front/fix-108-ids-ofuscados.test.tsx` pasa 5 de 11 pruebas (6 fallidas).

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1,5 h
- **Depende de:** `FIX-107`.
- **Problema:** El frontend modelaba el identificador de instancia estrictamente como un número entero (`vmid: number` e `id: number` en `InventoryInstance`), validándolo con `Number.isSafeInteger(vmid) && vmid > 0`. Al migrar el backend hacia identificadores opacos/cifrados no secuenciales (ej. `"inst_k8Pq2v"`), el frontend rechazaba las respuestas con errores de validación de esquema en `instanceService.ts`.
- **Entregables:**
  1. En `instance.ts`: actualizar `InventoryInstance` para que `id` admita `string | number` y definir el tipo `InstanceRef`.
  2. En `instanceService.ts`: flexibilizar la validación para admitir identificadores de instancia alfanuméricos y enteros legados.
  3. En rutas y componentes: asegurar navegación `/instances/:instanceRef` sin conversiones destructivas a entero.
- **Criterio de éxito:** La tabla de instancias renderiza identificadores opacos y las operaciones de ciclo de vida despachan la referencia en la URL sin conversiones numéricas destructivas. Los defectos de sanitización estricta quedan derivados a `FIX-124`.

---
# ETAPA 3
---


#### `BAC-44` - Interceptor global y consultas complejas sobre auditoría particionada
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-41`, `BAC-53` (`BRG-13-BAC`: particionado preventivo e índices compuestos de auditoría).
- **Problema:** La auditoría debe consolidar eventos de todas las etapas con filtros combinados eficientes (por usuario, rango de fechas, tipo de recurso, acción y severidad) aprovechando las particiones de PostgreSQL.
- **Entregable:**
  1. Extensión del interceptor de auditoría para registrar: `SNAPSHOT_CREATE`, `SNAPSHOT_ROLLBACK`, `SNAPSHOT_DELETE`, `ALERT_TRIGGERED` y `CONFIG_UPDATED`.
  2. Endpoint `GET /api/admin/audit` con parámetros de consulta: `?desde=...&hasta=...&usuarioId=...&accion=...&recursoTipo=...&recursoId=...&page=1&limit=25`.
  3. Optimización de consulta SQL utilizando exclusión de particiones (`partition pruning`) por `fecha_hora`.
- **Criterio de éxito:** Consultas sobre millones de registros históricos responden en menos de 150 ms sin realizar escaneos completos de tablas.

#### `BAC-45` - Streaming de exportación estructurada de auditoría a CSV y JSON
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-44`.
- **Problema:** Exportar reportes de auditoría extensos cargando todos los registros en memoria del servidor puede provocar caídas por `Out of Memory` (OOM).
- **Entregable:**
  1. Endpoint `GET /api/admin/audit/export` accesible solo para usuarios con rol `ADMIN`.
  2. Parámetro `format=csv|json`, respetando exactamente los mismos filtros aplicados en la vista.
  3. Implementación de streaming HTTP: lectura en lotes con cursores de base de datos (`pgx` / `gorm.Rows`) y escritura inmediata en el stream de respuesta.
  4. Cabeceras HTTP apropiadas: `Content-Disposition: attachment; filename="auditoria_centinela_YYYYMMDD.csv"`.
- **Criterio de éxito:** Exportaciones de más de 50.000 filas se descargan fluidamente en el navegador sin picos de consumo de RAM en el contenedor backend.


#### `INF-11` - Configuración de producción en Nginx: TLS 1.3, HSTS, CSP y WSS
- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 2.0 h
- **Depende de:** `INF-04`, `INF-09`.
- **Problema:** El entorno de producción final requiere estándares estrictos de cifrado, mitigación de ataques web y estabilidad total en canales bidireccionales.
- **Entregable:**
  1. Configuración de Nginx en el servidor de producción con protocolos modernos (TLS 1.2 y TLS 1.3 exclusivamente).
  2. Inyección de encabezados de seguridad:
     - `Strict-Transport-Security: max-age=31536000; includeSubDomains`
     - `X-Content-Type-Options: nosniff`
     - `X-Frame-Options: DENY`
     - `Content-Security-Policy` ajustada para WebSockets y fuentes locales.
  3. Verificación de proxying sin caídas de sesión ni timeout en WebSockets `/api/ws/*` y eventos SSE `/api/events`.
- **Criterio de éxito:** Análisis SSL Labs o verificación con curl reporta calificación A/A+ con todos los encabezados de seguridad presentes.



#### `FRN-29` - Pestaña "Snapshots" en Detalle de Instancia y formulario de nuevo snapshot
- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.0 h
- **Depende de:** `BAC-42`, `FRN-36` (`BRG-10-FRN`: detalle modularizado y estados transitorios).
- **Problema:** La vista de detalle no cuenta con una interfaz visual para inspeccionar los puntos de restauración existentes ni disparar un nuevo guardado.
- **Entregable:**
  1. Pestaña "Snapshots" dentro de `InstanceDetailPage.tsx`.
  2. Línea de tiempo o tabla interactiva de snapshots con nombre, descripción, fecha de creación y badge de estado.
  3. Botón "Crear Snapshot" que despliegue un modal sencillo solicitando nombre (validado contra `/^[a-zA-Z0-9_\-]+$/`) y descripción opcional.
  4. Feedback inmediato con spinner de carga durante la creación.
- **Criterio de éxito:** El usuario visualiza la lista de snapshots y puede crear uno nuevo viendo la lista actualizada reactivamente.


#### `FRN-31` - Vista administrativa de Auditoría con filtros reactivos completos
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.5 h
- **Depende de:** `BAC-44`, `FIX-74`.
- **Problema:** Los administradores necesitan auditar la actividad institucional de forma ágil, buscando incidentes o filtrando anomalías operativas.
- **Entregable:**
  1. Rediseño funcional de `Auditoria.tsx`:
     - Rango de fechas con controles accesibles (`aria-label="Desde"` y `aria-label="Hasta"`).
     - Desplegables con autocompletado para tipo de acción, usuario y tipo de recurso.
     - Indicador visual de severidad (Badge verde = `SUCCESS`, rojo = `FAILED`, amarillo = `WARNING`).
  2. Tabla paginada con controles de avance/retroceso y selector de elementos por página.
  3. Modal para inspeccionar el payload JSON completo (`detalles`) de un evento seleccionado.
- **Criterio de éxito:** Los cambios en los filtros actualizan reactivamente la tabla vía query parameters cumpliendo con todas las aserciones de accesibilidad.


#### `FRN-30` - Modal destructivo de Rollback y eliminación con confirmación antierror
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `FRN-29`, `BAC-43`.
- **Problema:** Revertir un snapshot descarta todos los cambios generados desde su captura. Un clic involuntario causaría pérdida irreparable de datos en el cliente.
- **Entregable:**
  1. Modal antierror `RollbackSnapshotModal.tsx` con advertencia en rojo: *"Esta acción restaurará la máquina al punto seleccionado. Todos los datos posteriores se perderán."*
  2. Input de seguridad que exija escribir textualmente el nombre del snapshot antes de habilitar el botón "Revertir".
  3. Modal de confirmación para eliminación de snapshots.
  4. Bloqueo de la botonera de la instancia con indicador *"Restaurando estado..."* mientras la tarea asíncrona se resuelve.
- **Criterio de éxito:** El botón de rollback solo es presionable tras tipear el nombre idéntico. Al confirmar, la interfaz entra en modo de espera y muestra un Toast al finalizar.


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


#### `BAC-48` - Hardening de autenticación: Cookies `HttpOnly`, Rate Limiting y Sanitización
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2.0 h
- **Depende de:** `BAC-03`, `BAC-17B`.
- **Problema:** Almacenar tokens de larga duración en `sessionStorage` deja una ventana de exposición ante vulnerabilidades XSS. Además, los endpoints de login y recuperación requieren protección estricta contra fuerza bruta.
- **Entregable:**
  1. Modificación de `POST /api/auth/login` y `POST /api/auth/refresh`: emisión del `refreshToken` en cookie HTTP con directivas `HttpOnly`, `Secure`, `SameSite=Strict`, `Path=/api/auth`.
  2. Middleware de Rate Limiting por IP para endpoints públicos (`/auth/login`, `/auth/password/forgot`) limitando a 5 intentos por minuto con respuesta `429 Too Many Requests`.
  3. Sanitización estricta de encabezados y payloads JSON contra inyecciones de caracteres de control.
- **Criterio de éxito:** Pruebas automatizadas confirman que JavaScript no puede leer la cookie de refresh y que peticiones abusivas consecutivas reciben 429.

#### `FRN-34` - Adaptación del cliente a cookies seguras de sesión y auditoría de accesibilidad
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.0 h
- **Depende de:** `BAC-48`.
- **Problema:** El cliente de API debe adaptarse para enviar credenciales automáticas en llamadas de refresh y la interfaz debe cumplir con estándares de accesibilidad para la entrega final.
- **Entregable:**
  1. Ajuste de `apiClient.ts` para enviar `credentials: 'include'` en peticiones al backend.
  2. Eliminación de referencias al almacenamiento manual del `refreshToken` en `sessionStorage`.
  3. Auditoría y corrección de accesibilidad (WCAG 2.1 AA) en modales, botones y selectores, garantizando navegación completa por teclado y etiquetas `aria-*`.
- **Criterio de éxito:** El flujo de rotación de tokens funciona automáticamente mediante cookies y las pruebas de accesibilidad pasan sin advertencias.


#### `BAC-54` (`BRG-14-BAC`) - Extracción y Parametrización en `.env` para Servidores Heterogéneos
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1.5 h
- **Depende de:** `INT-07`.
- **Problema:** En el código actual persisten suposiciones sobre el nombre de nodo (`proxmox` o `pve`) y el rango de IDs protegidos (`100-105`). En el servidor nuevo con PVE 8.4.21, el nodo y los VMIDs protegidos varían y no deben bloquear cargas del anfitrión.
- **Entregable:**
  1. Parametrización estricta en `backend/cmd/api/main.go` y `middleware/protected_instance.go`: lectura prioritaria de `PROXMOX_NODE`, `PROXMOX_URL` y `PROXMOX_PROTECTED_VMIDS`.
  2. Fallback seguro: si `PROXMOX_PROTECTED_VMIDS` viene vacío o no se define, no bloquear IDs arbitrarios ni fallar el arranque.
  3. Soporte para la bandera opcional `PROXMOX_TARGET_PVE_VERSION` (por defecto `"auto"`).
- **Criterio de éxito:** Modificar el `.env` con un nodo alternativo y VMIDs protegidos personalizados altera el comportamiento del backend sin requerir recompilación.

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
