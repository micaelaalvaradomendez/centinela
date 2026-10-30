# Requerimientos de El Centinela

> **Alcance del documento:** Especificación técnica y funcional completa de la plataforma **El Centinela**, alineada con la arquitectura hexagonal implementada en la **Fase Base** y la **Etapa 1**, las definiciones de las **Etapas 2 y 3**, las Decisiones de Diseño Cerradas (`D1`–`D6`) y las remediaciones de las auditorías de seguridad e infraestructura (`FIX-105` a `FIX-119`).
>
> **Criterio de clasificación de las remediaciones de seguridad (`FIX-105` a `FIX-119`):**
> - **En los Requerimientos Funcionales (`RF`):** se especifican aquellos controles que **definen o alteran un flujo de negocio, una regla de validación de dominio o la interacción del usuario** (por ejemplo: doble canal con contraseña + 2FA para cambiar el email en `RF-01`, límites de reintentos y bloqueo de fuerza bruta en `RF-01`/`RF-13`, segregación de métricas del host según el rol en `RF-02`, purga transaccional de permisos huérfanos al borrar una VM en `RF-04`/`RF-09`, validación de existencia real y protección de VMIDs al asignar permisos en `RF-09`).
> - **En los Requerimientos No Funcionales (`RNF`):** se especifican las **propiedades transversales de arquitectura, criptografía, ocultamiento de topología y saneamiento** que protegen a toda la plataforma sin alterar la lógica funcional del usuario (por ejemplo: enmascaramiento de `UPID` tras `tareaId`, ofuscación reversible de `VMID` contra IDOR, alias lógicos para nodos y storages, sanitización de trazas Perl/Linux de Proxmox y prohibición de credenciales por defecto en producción).

---

## Requerimientos funcionales

### RF-01. Autenticación, sesiones seguras, gestión de cuenta propia y control de acceso

El sistema debe gestionar el inicio de sesión mediante correo electrónico (normalizado a minúsculas) y contraseña, sin permitir el auto-registro público: las cuentas son aprovisionadas exclusivamente por un administrador (`RF-09`). Después de validar las credenciales primarias, debe exigir obligatoriamente autenticación de doble factor (2FA) mediante contraseñas temporales basadas en tiempo (TOTP, RFC 6238, 6 dígitos, ventana de 30 segundos).

El flujo de autenticación y gestión de credenciales debe cumplir las siguientes reglas de negocio:
1. **Fase de pre-autenticación (`POST /api/auth/login`):** Al validar correo, estado activo (`activo = true`, `eliminado_en IS NULL`) y contraseña (`bcrypt`), el sistema emite únicamente un token temporal (`jwtTemporal` con claim `type: "pre-2fa"`, `jti` en Redis y TTL de 5 minutos) e informa si el usuario ya tiene vinculado su 2FA (`totpVinculado`) y si debe cambiar su contraseña temporal (`cambioContrasenaRequerido`).
2. **Fase de segundo factor (`GET /api/auth/2fa/qr` y `POST /api/auth/2fa/verify`):**
   - Si el usuario **no** tiene vinculado el 2FA (`totpVinculado == false`), el sistema genera el secreto TOTP, lo cifra con **AES-256-GCM** y presenta el código QR junto con el secreto manual.
   - Si el usuario **ya** tiene vinculado el 2FA (`totpVinculado == true`), el sistema solicita el código TOTP de 6 dígitos y **bloquea cualquier intento de volver a consultar o regenerar el código QR** (`409 TWO_FACTOR_ALREADY_ENABLED`). La revinculación por pérdida de dispositivo solo puede ejecutarla un administrador (`RF-09`).
   - Al verificar el código TOTP, el sistema debe aplicar protección **anti-replay** registrando la ventana temporal consumida (`ultimo_totp_periodo`) para impedir que un mismo código de 6 dígitos se reutilice dentro de los mismos 30 segundos.
3. **Protección contra fuerza bruta en Login y 2FA (`FIX-112`):**
   - En `POST /api/auth/login`: permitir un máximo de **5 intentos fallidos** por par `(IP, Email)` dentro de una ventana de **15 minutos**; al sexto intento, bloquear con **`429 Too Many Requests`** y cabecera `Retry-After: 900`.
   - En `POST /api/auth/2fa/verify`: permitir un máximo de **3 intentos fallidos** por `jtiTemporal`. Al tercer fallo consecutivo, revocar y eliminar de inmediato la sesión pre-auth en Redis (`auth:pre2fa:<jti>`), respondiendo **`401 TOKEN_REVOKED`** y obligando al usuario a reiniciar el login desde el primer paso.
4. **Política de longitud y complejidad de contraseñas conforme a NIST SP 800-63B (`FIX-115A` / `FIX-115B`):**
   - Toda contraseña (temporal, nueva o restablecida) debe tener entre **8 y 128 caracteres** (eliminando el límite restrictivo previo de 12 caracteres para permitir el uso de *passphrases* y gestores de contraseñas como Bitwarden, 1Password o KeePass), e incluir al menos 1 letra mayúscula, 1 número y 1 carácter especial (`!@#$%^&*-_=+`).
5. **Cambio obligatorio o voluntario de contraseña y revocación de sesiones (`PUT /api/account/password` — `FIX-114`):**
   - Mientras `cambio_contrasena == true`, el middleware `RequireAuth` debe bloquear el acceso a cualquier recurso del sistema (`403 PASSWORD_CHANGE_REQUIRED`) excepto a `PUT /api/account/password` y `POST /api/auth/logout`.
   - Al ejecutarse un cambio de contraseña (obligatorio o voluntario), el backend debe verificar la contraseña actual, prohibir que la nueva sea idéntica a la anterior y **revocar obligatoriamente todas las demás sesiones activas del usuario** en `sesiones_activas` y Redis (emitiendo `MensajeSesionesRevocadas` por SSE/WebSocket para expulsar inmediatamente a cualquier intruso con tokens previos), auditando `sesiones_revocadas: true`.
6. **Blindaje contra secuestro de cuenta (*Account Takeover*) en cambio de correo propio (`PUT /api/account/profile` — `FIX-111A` / `FIX-111B`):**
   - Si el usuario solicita modificar su `emailUsuario`, el sistema debe exigir obligatoriamente su **`contrasenaActual`** y un **código `codigo2fa` TOTP válido**.
   - El correo **no debe mutar inmediatamente** en la tabla `usuarios`: el backend debe generar un token criptográfico de 32 bytes en Redis (`auth:email_change:<token>`, TTL 15 minutos), enviar el enlace/código de verificación a la nueva dirección, enviar una alerta de seguridad al correo anterior y registrar en `auditoria` tanto la solicitud (`UPDATE_PROFILE_EMAIL_REQUESTED`) como la confirmación (`UPDATE_PROFILE_EMAIL_CONFIRMED`).

**Frontend**

- Maquetar las vistas de inicio de sesión, vinculación inicial de 2FA mediante QR/clave manual, verificación de código TOTP de 6 dígitos y cambio obligatorio de contraseña temporal.
- Admitir contraseñas de **8 a 128 caracteres** (`maxLength={128}`) en todos los formularios y validadores (`validatePasswordComplexity.ts` — `FIX-115B`).
- En la edición de perfil propio (`FIX-111B`), interceptar cualquier cambio en el campo de correo electrónico y abrir un modal de seguridad que exija ingresar `contrasenaActual` y `codigo2fa` (TOTP), informando luego que se envió un correo de confirmación al nuevo destino.
- Mantener el `accessToken` JWT en memoria y utilizar la cookie `centinela_refresh` (`HttpOnly`, `Secure`, `SameSite=Strict`, `Path=/api/auth`) para la renovación silenciosa (`POST /api/auth/refresh`).

**Backend**

- Implementar `/api/auth/login`, `/api/auth/2fa/qr`, `/api/auth/2fa/verify`, `/api/auth/refresh`, `/api/auth/logout`, `/api/account/profile` y `/api/account/password`.
- Aplicar el middleware `RateLimiter` en Redis (*Sliding Window* / *Token Bucket*) sobre `/api/auth/login` y `/api/auth/2fa/verify` (`FIX-112`).
- Rotar *in-place* los identificadores `jti_access` y `jti_refresh` sobre la misma fila `sid` en `sesiones_activas` durante `/api/auth/refresh`.
- Revocar las sesiones activas en `CambiarContrasena` (`FIX-114`) y exigir re-autenticación + doble canal + auditoría en `ActualizarPerfil` (`FIX-111A`).

**Infraestructura**

- Persistir usuarios y sesiones activas en PostgreSQL y mantener el estado efímero (`pre2fa`, `sesion:<sid>`, tokens de cambio de email y contadores de rate limit) en Redis.
- Prohibir semillas con credenciales fijas (`Admin123!`) en producción (`FIX-118`, detallado en `RNF-01`).

---

### RF-02. Dashboard de estado del nodo y segregación por rol

Al ingresar al sistema, debe presentarse una vista general con el estado de salud y el consumo de recursos del nodo físico de Proxmox VE, evaluado mediante el **Semáforo de Salud (`D3`)**:
- **`Saludable` (Verde):** uso de CPU, RAM y almacenamiento inferior al `70%`.
- **`Advertencia` (Amarillo):** al menos un recurso se ubica en `>= 70%` (o en zona crítica `>= 85% / 90%`).
- **`Inaccesible` (Rojo):** pérdida de conectividad o timeout (`502` / `504`) al consultar la API de Proxmox VE.

**Regla de segregación de información según el rol (`FIX-117A` / `FIX-117B`):**
- **Para el rol `ADMIN`:** `GET /api/node/status` y `GET /api/node/history` entregan la telemetría global completa del servidor anfitrión: CPU total y núcleos, memoria RAM usada/total y capacidad asignable (`allocatable`, descontando el **15% de reserva de seguridad** para el host), almacenamiento físico de los discos del hipervisor, `uptime` y `instancesSummary` global de todo el clúster (`running`, `stopped`, `paused`, `total`).
- **Para el rol `OPERATOR`:** para evitar la fuga de información de infraestructura física y de cargas de trabajo de otras áreas (CWE-200), el backend debe **omitir el almacenamiento físico global del nodo y el `instancesSummary` global del clúster**, entregando únicamente el conteo y resumen acotado a las instancias que el operador tiene asignadas en `permisos_instancia`.

**Frontend**

- Diseñar tarjetas de resumen, medidores porcentuales y el indicador visual del semáforo de salud (`Saludable`, `Advertencia`, `Inaccesible`).
- Adaptar el Dashboard según el rol (`FIX-117B`) y manejar de forma defensiva valores `null`/`undefined` o respuestas no-array en los contadores y series (`FIX-100`, `FIX-103`) para que la vista de un `OPERATOR` renderice limpiamente su resumen acotado sin arrojar `TypeError`.

**Backend**

- Exponer `GET /api/node/status` y `GET /api/node/history` con caché corta en Redis (5 a 10 segundos).
- Enmascarar el nombre técnico del nodo físico mediante un alias lógico seguro (`FIX-109`, ej. `"nodo-principal"`) y segregar el payload de respuesta inspeccionando el rol en el contexto (`ContextKeyRol` — `FIX-117A`).

**Infraestructura**

- Garantizar que el token de API de Proxmox VE posea el privilegio `Sys.Audit` sobre el nodo.

---

### RF-03. Inventario de instancias con identificadores ofuscados

El sistema debe presentar un listado unificado de todas las máquinas virtuales (`QEMU` normalizado a `"vm"`) y contenedores (`LXC` normalizado a `"lxc"`) gestionados en Proxmox VE a los que el usuario tenga acceso (`GET /api/instances` y `GET /api/instances/:instanceRef`).

Cada elemento del inventario debe incluir: identificador público ofuscado (`id` / `instanceRef`, ej. `"inst_k8Pq2v"` según `FIX-107` / `FIX-108` para evitar exponer el `VMID` secuencial entero), nombre (`name`), tipo (`vm` o `lxc`), estado (`running`, `stopped`, `paused`), dirección IPv4 privada en vivo (`ip`), uso y tope de CPU (`cpuUsage`, `maxCpu`), uso y tope de RAM (`ramUsage`, `maxRam`), uso y tope de disco (`diskUsage`, `maxDisk`), tiempo de actividad (`uptime`), el nivel de permiso efectivo del usuario (`nivelAcceso`: `"FULL_ACCESS"` o `"READ_ONLY"`) y, si existe una operación asíncrona en ejecución sobre la máquina, el objeto `activeTask` (`tareaId`, `tipoAccion`, `estado`).

**Frontend**

- Soportar identificadores de instancia como cadenas alfanuméricas opacas (`id: string | number` — `FIX-108`) tanto en la tabla como en las rutas (`/instances/:id`) y servicios, sin forzar conversiones numéricas destructivas.
- Implementar filtrado reactivo en memoria (`FIX-101`) combinando búsqueda por texto (nombre o identificador), selector de tipo (`Todas`, `VM`, `LXC`) y estado (`En ejecución`, `Detenidas`, `Pausadas`) manteniendo los contadores sincronizados.
- Habilitar o deshabilitar acciones según `nivelAcceso` (`FULL_ACCESS` vs `READ_ONLY`) con atributos de accesibilidad `aria-label` canónicos (`FIX-102`), y bloquear la fila cuando `activeTask != null`.

**Backend**

- Desarrollar `GET /api/instances` y `GET /api/instances/:instanceRef`, traduciendo internamente el token público ofuscado al `vmid` numérico de Proxmox VE (`OfuscarVMID` / `DesofuscarVMID` — `FIX-107`).
- Aplicar el filtrado RBAC estricto en la capa de dominio (`ADMIN` ve todas con `FULL_ACCESS`; `OPERATOR` ve únicamente las asignadas en `permisos_instancia`).
- Resolver la dirección IPv4 privada en vivo de forma concurrente y acotada (`BAC-23B`) **únicamente después de filtrar por RBAC** y solo para instancias en estado `"running"`.
- Cruzar con `tareas_asincronas` para adjuntar `activeTask` (exponiendo `tareaId` y nunca el `UPID` crudo — `FIX-105`).

**Infraestructura**

- Configurar la variable secreta `INSTANCE_ID_SALT` para la ofuscación reversible de identificadores de instancia (`FIX-107`).

---

### RF-04. Control del ciclo de vida de instancias, seguimiento por `tareaId` y purga de permisos

El sistema debe permitir ejecutar acciones de energía sobre una instancia: iniciar (`start`), apagar ordenadamente mediante señal ACPI (`shutdown`), forzar el corte inmediato de energía (`stop`) y reiniciar (`reboot`), además de permitir la eliminación definitiva de una instancia detenida (`delete`, exclusivo para `ADMIN`).

Toda solicitud sobre una instancia debe atravesar una **cadena de guardas en orden estricto** antes de contactar a Proxmox VE:
1. **Autorización RBAC (`RequireInstanceAccess`):** Para acciones de energía se exige rol `ADMIN` o `OPERATOR` con `FULL_ACCESS` sobre la instancia (`403 INSTANCE_ACCESS_DENIED` si es `READ_ONLY` o sin acceso). Para `DELETE /api/instances/:id` se exige exclusivamente rol `ADMIN`.
2. **Protección de infraestructura crítica (`SEC-05` / `RejectProtectedInstance`):** Si la instancia pertenece a la lista de recursos del propio sistema (`PROXMOX_PROTECTED_VMIDS = 100, 101, 102, 103, 104, 105`), la petición se rechaza incondicionalmente con **`403 INSTANCE_PROTECTED`** para todos los roles (incluido `ADMIN`).
3. **Bloqueo de concurrencia por instancia (`BAC-51` / `FIX-33` / `RequireInstanceIdle`):** Si la instancia ya tiene una tarea activa (`RUNNING`) en `tareas_asincronas` o un lock en Redis, se rechaza con **`409 INSTANCE_BUSY`**.
4. **Validación previa de estado (`FIX-34`):** `start` exige `stopped`; `stop`, `shutdown` y `reboot` exigen `running`; `delete` exige `stopped`. Si no coincide, se rechaza con **`409 INSTANCE_INVALID_STATE`** (o `409 INSTANCE_NOT_STOPPED`).

**Enmascaramiento de `UPID` y seguimiento por `tareaId` (`FIX-105` / `FIX-106`):**
Al despachar la orden a Proxmox VE, el backend captura el `UPID`, lo guarda en la columna privada `tareas_asincronas.upid_proxmox` y responde inmediatamente **`202 Accepted`** entregando **únicamente el UUID opaco `{ "tareaId": "<uuid>" }`** (jamás el string `UPID:...`). El worker pool `SeguimientoTareas` monitorea el `UPID` internamente hasta por **3 minutos (`D4`)** y notifica el resultado final por SSE (`TASK_FINISHED` correlacionado por `tareaId`).

**Purga transaccional de permisos huérfanos al eliminar una instancia (`BAC-24C` / `FIX-113`):**
Dado que Proxmox VE recicla los números de `VMID` libres al invocar `/cluster/nextid`, cuando una tarea de eliminación (`DELETE`) finaliza con éxito (`COMPLETED`), el backend debe **eliminar atómicamente en una transacción todas las filas de `permisos_instancia` asociadas a ese `vmid`** y registrar en `auditoria` la acción `INSTANCE_PERMISSIONS_PURGED` con la lista de usuarios afectados (evitando que una futura máquina creada con el mismo ID herede accesos de operadores anteriores). Si el borrado en Proxmox falla, los permisos no se alteran.

**Frontend**

- Consumir exclusivamente `response.tareaId` en `requestInstancePowerAction`, `deleteInstance` y `createInstance`, sin exigir la presencia del campo `upid` (`FIX-105`).
- Solicitar confirmación explícita mediante modal antes de cada acción de energía y exigir en la eliminación (`FRN-16B`) que el administrador **escriba manualmente el nombre exacto de la instancia**.
- Bloquear los controles de la instancia mientras `tareaId` esté en progreso y desbloquearlos al recibir `TASK_FINISHED` por SSE.

**Backend**

- Implementar `POST /api/instances/:id/status/:action` y `DELETE /api/instances/:id` respondiendo `202 Accepted` con `{ tareaId }`.
- Ejecutar en `seguimiento_tareas.go` la purga transaccional de `permisos_instancia` y el registro `INSTANCE_PERMISSIONS_PURGED` al completarse una eliminación (`FIX-113`), emitiendo `TASK_FINISHED` para que los operadores conectados desmonten la fila eliminada.
- Reanudar al arranque del backend el seguimiento de cualquier tarea que haya quedado en `RUNNING` (`BAC-25C`).

**Infraestructura**

- Otorgar al token de API de Proxmox VE el privilegio `VM.PowerMgmt` y `VM.Allocate`, excluyendo las ACL sobre los VMIDs protegidos `100..105`.

---

### RF-05. Monitor de métricas por instancia

El sistema debe mostrar gráficos dinámicos con la telemetría específica de cada instancia (porcentaje de CPU en escala directa $0-100\%$, memoria RAM utilizada y asignada, tráfico de red `netin`/`netout` y lectura/escritura de disco `diskread`/`diskwrite`) combinando:
1. **Flujo en tiempo real (WebSocket cada 2 segundos):** autenticado mediante ticket efímero de un solo uso (`POST /api/auth/ws-ticket`, TTL 30s en Redis).
2. **Histórico reciente (REST sobre RRDdata):** consulta de series temporales (`GET /api/instances/:id/metrics?timeframe=2h|24h`, con caché de 30s en Redis).

**Frontend**

- Solicitar el ticket efímero (`getWsTicket()`) y conectar el hook `useInstanceMetrics` (`FRN-22`) almacenando los puntos en una ventana deslizante en memoria, con reconexión por *exponential backoff* y cierre limpio al desmontar la vista.
- Manejar de forma **tolerante a nulos (`null-safe`)** los paquetes donde la instancia esté apagada (`status: "stopped"`, `cpu: null`, `mem: null`), pintando el estado «Instancia detenida» sin invocar métodos numéricos sobre `null` (`FRN-22`, `FRN-23`).

**Backend**

- Exponer `POST /api/auth/ws-ticket`, `GET /api/ws/instances/:id/metrics?ticket=<uuid>` (consumo atómico `GETDEL`) y `GET /api/instances/:id/metrics?timeframe=2h|24h` (`BAC-34`, `BAC-35`, `BAC-36`).
- Sanitizar todos los mensajes enviados por el WebSocket para que nunca filtren errores crudos de Perl/Linux ni identificadores `UPID` (`FIX-110`).

**Infraestructura**

- Configurar Nginx para soportar el *upgrade* de conexiones `WebSocket` (`Upgrade` y `Connection: upgrade`).

---

### RF-06. Gestión de snapshots

El sistema debe permitir administrar puntos de restauración (*snapshots*) de una instancia sobre `/api/instances/:id/snapshots`:
- **Listar snapshots (`GET`):** permitido para usuarios con acceso `READ_ONLY` o `FULL_ACCESS` (y `ADMIN`).
- **Crear (`POST`), Revertir/Rollback (`POST .../:snapname/rollback`) y Eliminar (`DELETE .../:snapname`):** permitidos únicamente para `ADMIN` u `OPERATOR` con `FULL_ACCESS`, respetando `RejectProtectedInstance` (`100..105`) y `RequireInstanceIdle` (`409 INSTANCE_BUSY`).

**Frontend**

- Ocultar los botones de creación, reversión y borrado si el usuario solo posee permiso `READ_ONLY`.
- Para la reversión (*rollback*), exigir en un **modal destructivo reforzado (`FRN-30`)** que el usuario **escriba manualmente el nombre exacto del snapshot** antes de habilitar el botón, mostrando el estado `ROLLING_BACK` y siguiendo la tarea mediante `tareaId` (`FIX-105`).

**Backend**

- Validar con expresión regular estricta que `snapname` sea alfanumérico/guiones (máx. 40 caracteres) y sanitizar la descripción (`SEC-12` / `FIX-116`).
- Gestionar la creación, reversión y borrado como tareas asíncronas devolviendo `202 Accepted` con `{ tareaId }`, auditando en dos tiempos (`PENDING` $\rightarrow$ `EXITO`/`FALLA`) y traduciendo la falta de soporte del storage a `422 SNAPSHOT_NOT_SUPPORTED`.

**Infraestructura**

- Otorgar al token de API de Proxmox VE los privilegios `VM.Snapshot` y `VM.Snapshot.Rollback`.

---

### RF-07. Asistente de creación de instancias (Wizard) y sanitización de hardware

El sistema debe proporcionar un asistente paso a paso (**exclusivo para el rol `ADMIN`**) para aprovisionar nuevas máquinas virtuales (`QEMU`) o contenedores (`LXC`) en 4 pasos (Tipo y Nombre, Recursos con medidor de cuota en vivo, Almacenamiento/Plantilla/Red y Resumen).

Antes de crear la instancia en Proxmox VE, el backend debe aplicar tres controles obligatorios:
1. **Sanitización estricta contra *Parameter Injection* (`SEC-10` / `FIX-116`):** La API de Proxmox interpreta parámetros como `net0` mediante listas clave-valor separadas por comas (`name=eth0,bridge=vmbr0,firewall=1`). Para impedir que un input como `vmbr0,rate=100,tag=20` inyecte directivas de hardware o VLANs no autorizadas (CWE-88), el backend debe validar:
   - `Bridge`: expresión regular estricta de interfaz de red `^[a-zA-Z0-9_.-]+$` (prohibiendo comas, espacios o caracteres de control).
   - `Nombre` (`hostname`): convención RFC 1123 `^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`.
   - `StorageID` y `TemplateID`: libres de comas y caracteres de control, resolviendo los aliases lógicos de almacenamiento (`FIX-109`) hacia el identificador real de Proxmox.
2. **Bloqueo distribuido en Redis (`BAC-37`):** adquirir el lock `lock:instance:provisioning` (`SETNX`, TTL 30s) para serializar peticiones concurrentes.
3. **Validación de cuota con reserva del 15% de RAM para el host (`BAC-37`):** verificar que $\text{RAM}_{\text{solicitada}} \le (\text{RAM}_{\text{total\_host}} \times 0.85) - \text{RAM}_{\text{asignada\_instancias}}$ y que exista espacio libre en el disco destino; en caso contrario, rechazar con **`409 RESOURCE_QUOTA_EXCEEDED`**.

**Frontend**

- Diseñar el wizard de 4 pasos (`FRN-24`, `FRN-25`) con barras de capacidad reactivas que bloqueen el avance si la configuración supera el 100% del remanente asignable del host.
- Consumir la respuesta `202 Accepted` validando `{ vmid, tareaId }` sin requerir el campo privado `upid` (`FIX-105`).

**Backend**

- Implementar `GET /api/node/storages` (enmascarando nombres técnicos de pools de almacenamiento con etiquetas/aliases seguros según `FIX-109`), `GET /api/node/templates` y `POST /api/instances` (`RequireRole("ADMIN")`).
- Al finalizar con éxito el `UPID` de creación en `seguimiento_tareas.go` (`BAC-39`), registrar `CREATE_INSTANCE` (`EXITO`) en `auditoria` y emitir `TASK_FINISHED` por SSE.

**Infraestructura**

- Otorgar al token de API de Proxmox VE los privilegios `VM.Allocate` y `Datastore.AllocateSpace`, manteniendo plantillas LXC e ISOs disponibles (`INF-10`).

---

### RF-08. Registro de auditoría inmutable y exportación

El sistema debe registrar de forma **estrictamente inmutable (*append-only*)** todas las operaciones de autenticación, cambios de perfil/correo (`FIX-111A`), administración de usuarios/permisos, purga de permisos por borrado de VM (`INSTANCE_PERMISSIONS_PURGED` — `FIX-113`), cambios de configuración, acciones de ciclo de vida, snapshots y alertas.

Por seguridad (`FIX-105`), los registros visibles y exportables de auditoría **no deben exponer la cadena cruda `UPID` de Proxmox VE**, correlacionando las operaciones asíncronas mediante el identificador `tareaId`.

**Frontend**

- Diseñar la vista de auditoría (`/audit`) exclusiva para `ADMIN`, con filtros por usuario, acción, instancia, resultado y rango de fechas, visor lateral (Drawer) para el campo JSONB `detalles` (`FRN-31`) y botones de descarga «Exportar CSV» y «Exportar JSON» (`FRN-32`).

**Backend**

- Exponer `GET /api/admin/audit` y `GET /api/admin/audit/export?formato=csv|json` en **modo streaming por cursor** (`BAC-45`), auditando cada descarga con la acción `EXPORTAR_AUDITORIA`.

**Infraestructura**

- Garantizar la inmutabilidad en PostgreSQL mediante `REVOKE UPDATE, DELETE, TRUNCATE ON auditoria` y el trigger `trg_auditoria_inmutable` (`INF-15` / `SEC-11`), junto con particionado trimestral preventivo (`auditoria_2026_q3`, `auditoria_2026_q4`, `auditoria_2027_q1`) e índice compuesto `idx_auditoria_query` (`BAC-53`).

---

### RF-09. Gestión de usuarios, ciclo de vida de cuentas y validación estricta de permisos por instancia

Los administradores (`ADMIN`) deben poder gestionar las cuentas de usuario del sistema y asignar a cada operador (`OPERATOR`) las instancias autorizadas con nivel `READ_ONLY` o `FULL_ACCESS`.

El sistema debe aplicar las siguientes reglas y salvaguardas en la gestión de cuentas y permisos:
1. **Alta con contraseña temporal y transaccionalidad SMTP (`D5`):** `POST /api/admin/users` genera una contraseña temporal aleatoria, la envía por correo y marca `cambio_contrasena = true`. Si falla el servidor SMTP, revierte la creación y devuelve `502 EMAIL_DELIVERY_FAILED`.
2. **Suspensión reversible vs. Baja lógica irreversible (`FIX-54..56`):**
   - `PUT /api/admin/users/:id` con `activo: false` suspende temporalmente al usuario y revoca sus sesiones activas (`USER_INACTIVE`).
   - `DELETE /api/admin/users/:id` aplica baja lógica definitiva (`eliminado_en = NOW()`, `activo = false`), revoca sesiones, limpia sus permisos y libera su correo/username mediante un índice único parcial (`WHERE eliminado_en IS NULL`) sin borrar físicamente la fila para no romper la clave foránea histórica de `auditoria`.
3. **Validación de existencia real y bloqueo de instancias protegidas al asignar permisos (`PUT /api/admin/users/:id/permissions` — `FIX-119`):**
   - Antes de persistir en `permisos_instancia`, el backend debe consultar el catálogo actual de instancias en Proxmox VE y verificar que **cada `vmid` exista realmente**; si alguno no existe, rechazar con **`400 INVALID_REQUEST`**.
   - Debe verificar que **ningún `vmid` pertenezca a las instancias protegidas del sistema (`PROXMOX_PROTECTED_VMIDS = 100..105`)**. Si se intenta asignar una máquina protegida a un operador, rechazar con **`403 PROTECTED_INSTANCE_CANNOT_BE_ASSIGNED` (`INSTANCE_PROTECTED`)** para impedir que un operador descubra o monitoree por WebSocket la infraestructura interna de El Centinela.
4. **Auto-protección y reseteo administrativo:** un `ADMIN` no puede eliminarse, desactivarse ni quitarse el rol a sí mismo (`400 SELF_DELETE_NOT_ALLOWED`). Los endpoints `POST /api/admin/users/:id/password/reset` y `/2fa/reset` deben revocar en el acto todas las sesiones activas del usuario afectado.

**Frontend**

- Construir el panel de administración de usuarios (`/users`) y el selector de permisos por instancia (`READ_ONLY` / `FULL_ACCESS`), mostrando únicamente las instancias operables del inventario y excluyendo las protegidas.

**Backend**

- Implementar `/api/admin/users`, `/api/admin/users/:id/permissions`, `/password/reset` y `/2fa/reset` con las validaciones de `FIX-119` y la purga automática al eliminar instancias (`FIX-113`).

**Infraestructura**

- Mantener la matriz de roles y permisos íntegramente en PostgreSQL (`usuarios`, `permisos_instancia`), aislada de los usuarios del sistema operativo Linux y de Proxmox VE.

---

### RF-10. Edición de recursos de instancias (Escalado vertical)

El sistema debe permitir a los usuarios con rol `ADMIN` o permiso `FULL_ACCESS` sobre una instancia modificar la cantidad de núcleos virtuales (`cores`) y memoria `RAM (MB)` asignados (`PUT /api/instances/:id/config`).

El backend debe aplicar la **validación diferencial de cuota (`BAC-40`)**:
- Si $\Delta\text{RAM} = \text{memoryMB}_{\text{nuevo}} - \text{memoryMB}_{\text{actual}} > 0$, validar que el incremento no supere la memoria libre asignable del nodo físico (respetando el 15% de reserva de RAM para el host); si la excede, rechazar con **`409 RESOURCE_QUOTA_EXCEEDED`**.
- Si $\Delta\text{RAM} \le 0$ (reducción o cambio solo de vCPU), autorizar directamente.
- Inspeccionar si la instancia está `running` y si soporta *hot-plug*, retornando `{ requiresReboot: true | false }` y registrando en `auditoria` (`CONFIG_INSTANCE` / `UPDATE_INSTANCE_CONFIG`) los valores anteriores y nuevos (`BAC-41`).

---

### RF-11. Notificaciones en tiempo real y alertas proactivas

El sistema debe mantener un canal de eventos en tiempo real mediante **Server-Sent Events (SSE)** (`POST /api/events/ticket` $\rightarrow$ `GET /api/events?ticket=<uuid>`) y un motor de evaluación de umbrales (`EvaluadorUmbrales` cada 30s) con **histéresis de 2 ciclos consecutivos en Redis (`BAC-46`)** para evitar falsos positivos:
- **Eventos de instancia (`recursoTipo: "VM" | "LXC"`):** `TASK_FINISHED`, `ALERT_TRIGGERED` y `ALERT_RESOLVED`, distribuidos a todos los `ADMIN` y únicamente a los `OPERATOR` con permiso `READ_ONLY` o `FULL_ACCESS` sobre esa instancia.
- **Eventos globales del nodo (`recursoTipo: "HOST"` — `BAC-52`):** difundidos exclusivamente a usuarios con rol `ADMIN` y descartados silenciosamente para los `OPERATOR` (en coherencia con `FIX-117A`).
- **Eventos de control de sesión (`event: cierre`):** `LOGOUT`, `SESSIONS_REVOKED` (incluyendo cambio de contraseña según `FIX-114`) y `USER_INACTIVE`.

Por seguridad (`FIX-105`, `FIX-110`), ningún evento transmitido por SSE debe contener el campo `upid` crudo ni trazas de error internas del sistema de archivos de Proxmox VE.

---

### RF-12. Gestión de organizaciones (Descartado del alcance del MVP)

> **Decisión arquitectónica de alcance:** El requerimiento original de soporte multi-organización (`RF-12`) fue **descartado del alcance del MVP**. La plataforma opera bajo un modelo **mono-organización** gobernado por los roles `ADMIN` y `OPERATOR` y la tabla `permisos_instancia`.

---

### RF-13. Recuperación de contraseña mediante código OTP con bloqueo transaccional y Rate Limiting

El sistema debe permitir restablecer una contraseña olvidada mediante un código numérico temporal de un solo uso (OTP de **6 dígitos**, TTL 15 minutos) enviado por correo electrónico (`POST /api/auth/password/forgot` y `POST /api/auth/password/reset`), cumpliendo las reglas `D6` y el endurecimiento de seguridad **`FIX-112` y `FIX-115`**:
- **Prevención de enumeración y Anti-Flooding SMTP (`FIX-112`):** `POST /api/auth/password/forgot` responde siempre `200 OK` con mensaje genérico, pero aplica un **Rate Limit estricto en Redis de máximo 3 solicitudes por dirección de correo y por IP cada 60 minutos** (`429 Too Many Requests`) para impedir *Mail Bombing* o agotamiento de cuotas del relay SMTP.
- **Bloqueo pesimista contra condiciones de carrera en intentos fallidos (`FIX-112`):** En `POST /api/auth/password/reset`, la lectura, validación e incremento del contador de intentos fallidos (máximo **3 intentos** antes de quemar definitivamente el código OTP) debe ejecutarse dentro de una **transacción en PostgreSQL con bloqueo pesimista de fila (`SELECT ... FOR UPDATE`)**, impidiendo que peticiones paralelas simultáneas burlen el límite de 3 intentos.
- **Longitud NIST SP 800-63B y revocación global:** la nueva contraseña admite entre **8 y 128 caracteres** (`FIX-115`), no puede coincidir con la actual y, al confirmarse, revoca todas las sesiones activas del usuario en PostgreSQL y Redis (`SESSIONS_REVOKED`).

---

## Requerimientos no funcionales

### RNF-01. Seguridad, blindaje de infraestructura y aislamiento de metadatos (`FIX-105` a `FIX-119`)

Además del principio de mínimo privilegio en el token de Proxmox VE (`PVEAPIToken`), el hashing `bcrypt` (costo 12) y el cifrado **AES-256-GCM** para secretos TOTP, el sistema debe implementar una capa de **Air-Gap de Metadatos de Infraestructura** para que, aun si el cliente web fuera comprometido, un atacante no obtenga información útil para atacar el hipervisor subyacente:

1. **Enmascaramiento absoluto del `UPID` tras `tareaId` (`FIX-105` / `FIX-106` — CWE-200):**
   El identificador `UPID` nativo de Proxmox (`UPID:<node>:<pid_hex>:<pstart_hex>:<starttime_hex>:<type>:<vmid>:<user>:`) revela el hostname interno del servidor, el PID del proceso en el kernel Linux, timestamps de arranque y el usuario del token PVE. El backend debe almacenarlo exclusivamente en la columna interna `tareas_asincronas.upid_proxmox` y exponer hacia el frontend y la auditoría visible únicamente un **UUIDv4 opaco (`tareaId`)**.
2. **Ofuscación / Cifrado reversible de `VMID` contra IDOR y enumeración secuencial (`FIX-107` / `FIX-108` — CWE-639):**
   Los identificadores enteros secuenciales de Proxmox (`100`, `101`, `110`) permiten enumerar cuántas máquinas existen e intentar ataques IDOR iterativos. La capa HTTP del backend debe cifrar/ofuscar bidireccionalmente los `VMID` mediante un algoritmo con sal criptográfica (`INSTANCE_ID_SALT`, ej. *Hashids/Sqids* o cifrado simétrico produciendo tokens opacos como `"inst_k8Pq2v"`), y no debe distinguir en sus respuestas de error entre un token malformado y un recurso inexistente (`404`) para evitar oráculos de validación.
3. **Enmascaramiento de topología de Nodos y Pools de Almacenamiento (`FIX-109` — CWE-200):**
   El backend no debe exponer en los JSON públicos los nombres reales de red de los hipervisores (`"proxmox"`, `"pve-01"`) ni los identificadores técnicos de los drivers/pools de disco (`"local-lvm"`, `"zfs-pool"`). Debe traducirlos mediante un mapa de aliases lógicos seguros (`"nodo-principal"`, `"Almacenamiento Local Rápido"`) y resolverlos internamente antes de invocar a Proxmox VE.
4. **Saneamiento estricto de trazas de error y canales en tiempo real (`FIX-110` — CWE-209):**
   Cuando la API de Proxmox VE devuelve errores internos (que suelen incluir rutas del sistema de archivos como `/usr/share/perl5/PVE/API2/Qemu.pm line 450`, nombres de módulos Perl o sockets UNIX), el backend debe registrar el error crudo únicamente en el log interno del servidor y devolver al cliente HTTP, SSE o WebSocket un mensaje canonizado y libre de rutas o metadatos del sistema operativo (`PROXMOX_ERROR`, `PROXMOX_UNAVAILABLE`, `PROXMOX_TIMEOUT`).
5. **Prohibición de credenciales por defecto en despliegues productivos (`FIX-118` — CWE-798):**
   Las migraciones e inicialización de base de datos (`db.go`) **nunca deben insertar usuarios con contraseñas hardcodeadas** (como `admin@elcentinela.com` / `Admin123!`) ni imprimir credenciales en texto plano en los logs cuando `APP_ENV=production`. En entornos de desarrollo se debe generar una clave aleatoria con `crypto.GenerarContrasenaTemp()`, y en producción exigir variables de entorno secretas (`INITIAL_ADMIN_EMAIL`, `INITIAL_ADMIN_PASSWORD`) o el uso del comando interactivo `cmd/seed`.
6. **Blindaje de instancias críticas del sistema (`SEC-05` y `FIX-119`):**
   Las instancias en `PROXMOX_PROTECTED_VMIDS` (`100..105`) están protegidas tanto frente a operaciones de ciclo de vida/configuración (`RejectProtectedInstance` para todos los roles) como frente a su asignación en `permisos_instancia` (`FIX-119`).

### RNF-02. Rendimiento y concurrencia

El sistema debe responder en menos de 2 segundos para consultas en caché/base de datos y en un máximo de 5 segundos para lecturas en vivo hacia Proxmox VE, soportando múltiples operadores concurrentes sin degradación.

- **Concurrencia acotada en Go:** worker pool fijo `UPID_WORKERS=8` para el seguimiento de tareas y goroutines acotadas con timeout corto para la resolución de IPs (`BAC-23B`).
- **Particionado e indexación de base de datos (`BAC-53`):** particionado trimestral por rango en `auditoria` con *partition pruning* mediante `idx_auditoria_query` y exportación en streaming por cursor con memoria constante $O(1)$ (`BAC-45`).

### RNF-03. Usabilidad y prevención de errores humanos

- **Confirmaciones graduadas según el riesgo:** modales estándar para encendido/apagado limpio, advertencia destacada para corte de energía (`stop`) y **escritura manual obligatoria del nombre exacto** para eliminar una instancia (`FRN-16B`) o revertir un snapshot (`FRN-30`).
- **Compatibilidad con gestores de contraseñas (`FIX-115`):** soporte de contraseñas e inputs de hasta 128 caracteres conforme a NIST SP 800-63B.

### RNF-04. Confiabilidad, resiliencia y exclusión mutua

- **Exclusión mutua operativa (`BAC-51` / `RequireInstanceIdle`):** impide ejecutar acciones solapadas sobre una misma instancia (`409 INSTANCE_BUSY`).
- **Prevención de condiciones de carrera (`BAC-37` y `FIX-112`):** lock distribuido en Redis (`lock:instance:provisioning`) para aprovisionamiento y bloqueo pesimista `SELECT ... FOR UPDATE` en PostgreSQL para validación de códigos de recuperación.
- **Recuperación tras reinicios (`BAC-25C`) y timeout de tareas (`D4`):** reanudación automática de tareas `RUNNING` al iniciar el backend y corte a los 3 minutos si un `UPID` queda colgado.

### RNF-05. Mantenibilidad y arquitectura

- **Arquitectura Hexagonal (Puertos y Adaptadores):** desacoplamiento estricto entre el núcleo de dominio (`internal/core`) y los adaptadores primarios (HTTP, SSE, WebSocket) y secundarios (PostgreSQL, Redis, Proxmox VE, SMTP, Crypto).

### RNF-06. Compatibilidad y seguridad de transporte

- Compatibilidad con Chrome, Firefox, Brave, Edge y Safari recientes.
- CORS restringido a `CORS_ALLOWED_ORIGIN`, cookies `HttpOnly` + `SameSite=Strict` y verificación de certificado TLS contra Proxmox VE (`PROXMOX_CA_CERT_PATH`).

### RNF-07. Escalabilidad

- Almacenamiento de sesiones, tickets efímeros (`GETDEL`), caché, rate limiters, estado de histéresis de alertas y locks distribuidos en **Redis**, permitiendo escalar el backend o gestionar múltiples nodos sin acoplar estado en memoria local.

---

## Funcionalidades excluidas

### 1. Consola remota interactiva web (SSH / VNC / xterm.js / noVNC)
Excluida para mitigar el riesgo de ejecución de comandos arbitrarios en el sistema operativo huésped.

### 2. Despliegue de aplicaciones internas (CI/CD dentro de las VMs/LXC)
El alcance se limita al aprovisionamiento, hardware virtual, monitoreo y ciclo de vida de la infraestructura sobre Proxmox VE.

### 3. Multi-tenancy de múltiples organizaciones (`RF-12`)
Excluido del MVP en favor de un modelo mono-organización con aislamiento por roles (`ADMIN` / `OPERATOR`) y permisos por instancia (`FULL_ACCESS` / `READ_ONLY`).
