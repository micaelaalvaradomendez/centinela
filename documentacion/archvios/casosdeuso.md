# Especificación de Casos de Uso — El Centinela

> **Alcance del documento:** Este documento especifica los casos de uso reales de la plataforma **El Centinela**, integrando la arquitectura hexagonal implementada en la **Fase Base** y la **Etapa 1**, junto con los flujos, contratos y salvaguardas definidos para las **Etapas 2 y 3** y las tareas de endurecimiento de seguridad (`SEC-*`, `FIX-*`).
>
> **Nota sobre RF-12 (Gestión de Organizaciones):** Por decisión de alcance arquitectónico del MVP, el sistema opera bajo un modelo **mono-organización**. El requerimiento `RF-12` queda fuera del alcance operativo y el aislamiento de recursos se aplica directamente mediante Control de Acceso Basado en Roles (`RBAC`) y niveles de acceso por instancia (`FULL_ACCESS` / `READ_ONLY`).

---

## 1. Catálogo de Actores

### 1.1. Actores Principales (Humanos)

| Actor | Rol / Estado en el Sistema | Responsabilidades y Alcance de Permisos |
| :--- | :--- | :--- |
| **Usuario No Autenticado / Pre-autenticado** | Sin sesión o con `jwtTemporal` (`type: "pre-2fa"`, TTL 5 min) | Puede iniciar el flujo de autenticación (`POST /api/auth/login`), vincular por primera vez su dispositivo TOTP (`GET /api/auth/2fa/qr`), validar su código de 6 dígitos (`POST /api/auth/2fa/verify`) o solicitar la recuperación de contraseña olvidada mediante OTP por correo (`RF-13`). |
| **Usuario con Cambio de Contraseña Pendiente** | Autenticado tras 2FA con flag `cambio_contrasena = true` | Estado transitorio posterior al primer ingreso con contraseña temporal o tras un reseteo administrativo. El middleware `RequireAuth` bloquea el acceso a toda la plataforma (`403 PASSWORD_CHANGE_REQUIRED`) excepto a `PUT /api/account/password` y `POST /api/auth/logout`. |
| **Operador de Solo Lectura** | `OPERATOR` con permiso `READ_ONLY` sobre un `vmid` en `permisos_instancia` | Puede visualizar el dashboard, consultar en el inventario únicamente las instancias que tiene asignadas, inspeccionar su detalle, ver métricas en tiempo real e histórico (`RF-05`), listar snapshots (`RF-06`) y recibir notificaciones SSE (`TASK_FINISHED`, alertas) de esas instancias. Tiene bloqueada toda acción de energía, edición de hardware o mutación de snapshots (`403 INSTANCE_ACCESS_DENIED`). |
| **Operador con Control Operativo** | `OPERATOR` con permiso `FULL_ACCESS` sobre un `vmid` en `permisos_instancia` | Posee todos los permisos de `READ_ONLY` sobre el `vmid` asignado y además puede ejecutar acciones de energía (`start`, `stop`, `shutdown`, `reboot`), modificar recursos de hardware (`vCPU` y `RAM`) y gestionar snapshots (`crear`, `revertir`, `eliminar`), siempre que la instancia no pertenezca al rango protegido (`100..105`). No puede eliminar instancias, aprovisionar nuevas máquinas ni acceder a funciones administrativas. |
| **Administrador** | `ADMIN` | Posee permiso implícito `FULL_ACCESS` sobre todas las instancias del nodo sin requerir filas en `permisos_instancia` (respetando la protección de las instancias del sistema `100..105`). Es el único rol autorizado para aprovisionar nuevas instancias (`RF-07`), eliminar instancias detenidas (`DELETE /api/instances/:vmid`), administrar usuarios, roles y permisos (`RF-09`), consultar/exportar el registro inmutable de auditoría (`RF-08`) y recibir alertas globales del nodo físico (`RF-11`). |

### 1.2. Actores Secundarios (Sistemas Externos y Procesos Internos)

| Actor | Tipo | Rol en las Transacciones |
| :--- | :--- | :--- |
| **API de Proxmox VE (`ProxmoxPort`)** | Sistema externo (Hipervisor) | Interfaz REST/JSON consumida vía HTTPS con verificación de certificado TLS (`PROXMOX_CA_CERT_PATH`) y autenticación de mínimo privilegio mediante API Token (`PVEAPIToken=user@realm!tokenid=secret`). Ejecuta operaciones sobre el hipervisor y retorna identificadores de tareas asíncronas (`UPID`). |
| **Servicio de Correo (`EmailService`)** | Sistema externo (SMTP) | Despacha correos transaccionales: entrega de contraseñas temporales en altas o reseteos administrativos (`RF-01`, `RF-09`), códigos OTP de 6 dígitos para recuperación de contraseña (`RF-13`) y alertas críticas de infraestructura (`RF-11`). |
| **Worker de Seguimiento de Tareas (`SeguimientoTareas`)** | Proceso interno en segundo plano (Worker Pool Go) | Pool acotado de goroutines (`UPID_WORKERS=8`) que monitorea mediante polling el estado de los `UPID` en Proxmox VE hasta por un máximo de 3 minutos (`D4`), actualiza la tabla `tareas_asincronas`, registra el resultado definitivo en `auditoria` y publica el evento `TASK_FINISHED` en Redis Pub/Sub. |
| **Worker Evaluador de Umbrales (`EvaluadorUmbrales`)** | Proceso interno en segundo plano (Ticker Go) | Sondea cada 30 segundos la telemetría del nodo físico y de las instancias, evalúa umbrales de saturación con histéresis en Redis (2 ciclos consecutivos) y dispara alertas proactivas (`ALERT_TRIGGERED` / `ALERT_RESOLVED`) por SSE y correo electrónico. |

---

## 2. Salvaguardas Transversales y Reglas de Negocio Globales

Todas las transacciones de los casos de uso están gobernadas por las siguientes barreras de seguridad e integridad implementadas en el sistema:

1. **Protección Inquebrantable de la Infraestructura Propia (`SEC-05` / `RejectProtectedInstance`):**
   Las máquinas virtuales y contenedores que alojan la propia plataforma El Centinela y su entorno base (`PROXMOX_PROTECTED_VMIDS = 100, 101, 102, 103, 104, 105`) están blindadas a nivel de middleware. Cualquier intento de ejecutar acciones de energía, eliminar la instancia, alterar su hardware, revertir snapshots o asignarla a un operador es rechazado con **`403 Forbidden` (`INSTANCE_PROTECTED`)**, incluso si quien lo solicita tiene rol `ADMIN`.
2. **Orden Estricto de Evaluación de Guardas (Prevención de Fuga de Información):**
   En todo endpoint que opera sobre un `:vmid`, el backend evalúa primero la autenticación (`RequireAuth`), luego el permiso RBAC del usuario sobre el recurso (`RequireInstanceAccess`) y recién después la protección de instancia (`RejectProtectedInstance`) y el estado de ocupación (`RequireInstanceIdle`). Así, un `OPERATOR` sin permisos sobre la VM `100` recibe `403 INSTANCE_ACCESS_DENIED` y nunca puede inferir qué `VMID` corresponden a infraestructura crítica.
3. **Ocultamiento de Topología, Ofuscación de Identificadores y Saneamiento de Trazas (`FIX-105` a `FIX-110`):**
   - **`UPID` oculto tras `tareaId` (`FIX-105` / `FIX-106`):** El identificador crudo `UPID` de Proxmox VE (que expone el nombre interno del nodo, PIDs del kernel Linux y el usuario del token PVE) jamás se expone en las respuestas HTTP, en la auditoría visible ni en los eventos SSE hacia el frontend; se almacena en `tareas_asincronas.upid_proxmox` y solo se expone el UUIDv4 opaco `tareaId`.
   - **Ofuscación reversible de `VMID` (`FIX-107` / `FIX-108`):** Para prevenir la enumeración secuencial de máquinas y ataques IDOR (`100`, `101`, `110`), la capa HTTP cifra/ofusca bidireccionalmente los identificadores de instancia mediante `INSTANCE_ID_SALT` (produciendo tokens públicos opacos como `"inst_k8Pq2v"`).
   - **Aliases de Nodos y Storages (`FIX-109`) y Saneamiento de Errores (`FIX-110`):** Los nombres reales de red del host (`"proxmox"`, `"pve-01"`) y de los pools de disco (`"local-lvm"`, `"zfs-pool"`) se traducen a aliases lógicos seguros, y toda traza de error de Proxmox (rutas `/usr/share/perl5/...`) se registra solo en el log interno devolviendo códigos canonizados al cliente.
4. **Exclusión Mutua de Operaciones por Instancia (`BAC-51` / `FIX-33` / `RequireInstanceIdle`):**
   Ninguna instancia puede recibir dos órdenes mutables o asíncronas en simultáneo. Si existe un registro en `tareas_asincronas` con estado `RUNNING` para ese `vmid`, cualquier nueva solicitud de energía, cambio de hardware, snapshot o borrado se rechaza con **`409 Conflict` (`INSTANCE_BUSY`)**.
5. **Trazabilidad en Dos Tiempos sobre Auditoría Inmutable (`RF-08`):**
   Toda operación asíncrona sobre Proxmox genera dos registros en la tabla *append-only* `auditoria`: un primer registro al aceptar la orden (`resultado = PENDING`) y un segundo registro cuando el worker `SeguimientoTareas` verifica el término del `UPID` (`resultado = EXITO` o `FALLA`). La tabla `auditoria` posee un trigger en PostgreSQL (`trg_auditoria_inmutable`) que bloquea cualquier `UPDATE`, `DELETE` o `TRUNCATE`.

---

## 3. Especificación Detallada de Casos de Uso

---

### CU-01: Autenticarse, Vincular/Verificar 2FA y Gestionar Sesión

- **Requerimientos asociados:** `RF-01`, `RNF-01`.
- **Actor principal:** Usuario No Autenticado / Pre-autenticado / Usuario Autenticado (`ADMIN` u `OPERATOR`).
- **Actores secundarios:** PostgreSQL (`usuarios`, `sesiones_activas`, `auditoria`), Redis (`auth:pre2fa:*`, `sesion:*`), Canal SSE.
- **Endpoints involucrados:**
  - `POST /api/auth/login`
  - `GET /api/auth/2fa/qr`
  - `POST /api/auth/2fa/verify`
  - `PUT /api/account/password`
  - `POST /api/auth/refresh`
  - `POST /api/auth/logout`
  - `GET /api/account/profile` y `PUT /api/account/profile`

#### Descripción
El usuario ingresa sus credenciales primarias (correo electrónico y contraseña). Tras su validación, el sistema exige la verificación de segundo factor mediante un código TOTP de 6 dígitos (guiando al usuario en la vinculación por código QR si es su primer ingreso) y, en caso de poseer una contraseña temporal, fuerza su actualización antes de habilitar el acceso operativo a la plataforma.

#### Precondiciones
1. El usuario debe haber sido dado de alta previamente por un `ADMIN` (no existe auto-registro público en el sistema).
2. La cuenta debe encontrarse activa (`activo = true`) y no eliminada (`eliminado_en IS NULL`).

#### Postcondiciones
- **De éxito:** Se registra una sesión única en `sesiones_activas` (`sid`) y en Redis; el cliente recibe un `accessToken` JWT de corta duración en memoria y una cookie `centinela_refresh` (`HttpOnly`, `Secure`, `SameSite=Strict`); queda registrado el ingreso en `auditoria` (`LOGIN`, `EXITO`).
- **De fallo:** No se emite sesión definitiva y se registra el intento fallido en `auditoria` (`LOGIN`, `FALLA`).

#### Flujo Principal (Transacciones)
1. El usuario accede a la pantalla de inicio de sesión (`/login`) e ingresa su correo electrónico y contraseña.
2. El frontend envía `POST /api/auth/login` con `{ email, password }`.
3. El backend normaliza el correo a minúsculas, verifica en PostgreSQL que el usuario exista, esté activo (`activo = true`, `eliminado_en IS NULL`) y valida la contraseña contra el `hash_contrasena` (bcrypt costo 12).
4. El backend genera un token pre-autenticación `jwtTemporal` (con claim `type: "pre-2fa"`, `jti` único y expiración de 5 minutos), registra el estado pre-2FA en Redis (`auth:pre2fa:<jti>` con contador de intentos en `0`) y responde `200 OK` con `{ jwtTemporal, totpVinculado, cambioContrasenaRequerido }`.
5. **CONDICIONAL — Estado de vinculación 2FA (`totpVinculado`):**
   - **5A. Sin vinculación previa (`totpVinculado == false`):**
     1. El frontend invoca `GET /api/auth/2fa/qr` enviando `Authorization: Bearer <jwtTemporal>`.
     2. El middleware `RequirePreAuth` valida el token temporal. El backend verifica que `totp_vinculado` sea `false`, genera un secreto TOTP aleatorio compatible con RFC 6238 (6 dígitos, ventana de 30 s), lo cifra con **AES-256-GCM** usando la clave maestra `TOTP_ENCRYPTION_KEY`, lo guarda en `usuarios.secreto_totp` y retorna `{ qrBase64, secretoManual }`.
     3. El frontend renderiza el código QR y la clave manual para que el usuario lo escanee en su aplicación autenticadora (Google Authenticator, Authy, FreeOTP) y habilita el campo para ingresar el primer código de 6 dígitos.
   - **5B. Con vinculación previa (`totpVinculado == true`):**
     1. El frontend renderiza directamente la vista de verificación solicitando el código TOTP de 6 dígitos.
6. El usuario ingresa el código TOTP de 6 dígitos y el frontend envía `POST /api/auth/2fa/verify` con `Authorization: Bearer <jwtTemporal>` y `{ codigo }`.
7. El backend recupera y descifra en memoria el `secreto_totp` con AES-256-GCM, valida el código TOTP ingresado y verifica la **salvaguarda anti-replay** (`periodo_actual > usuario.ultimo_totp_periodo`).
8. El backend actualiza `totp_vinculado = true` y `ultimo_totp_periodo = periodo_actual`, elimina el `jti` temporal de Redis, crea el registro de sesión en `sesiones_activas` (`sid` UUID, `jti_access`, `jti_refresh`, IP y User-Agent) y su clave espejo en Redis (`sesion:<sid>`).
9. El backend registra `LOGIN` (`EXITO`) en `auditoria`, setea la cookie `centinela_refresh` (`HttpOnly`, `Secure`, `SameSite=Strict`, `Path=/api/auth`, `Max-Age=7d`) y responde `200 OK` con `{ accessToken, usuario, cambioContrasenaRequerido }`.
10. **CONDICIONAL — Contraseña temporal (`cambioContrasenaRequerido == true`):**
    - **10A. Si `cambioContrasenaRequerido == true`:**
      1. El frontend redirige obligatoriamente a la vista de cambio de contraseña. Mientras este flag siga en `true`, el middleware `RequireAuth` rechaza cualquier petición a otro endpoint del sistema con **`403 Forbidden` (`PASSWORD_CHANGE_REQUIRED`)**.
      2. El usuario ingresa su `contrasenaActual` y la `contrasenaNueva`, enviando `PUT /api/account/password`.
      3. El backend verifica que la contraseña actual sea correcta, que la nueva sea diferente de la actual y que cumpla la política de complejidad (al menos 1 mayúscula, 1 número, 1 carácter especial `!@#$%^&*-_=+` y longitud válida).
      4. El backend actualiza `hash_contrasena`, establece `cambio_contrasena = false`, revoca cualquier otra sesión activa del usuario y registra `CAMBIAR_CONTRASENA` (`EXITO`) en `auditoria`.
    - **10B. Si `cambioContrasenaRequerido == false`:**
      1. El frontend almacena el `accessToken` en el estado de sesión y redirige al usuario al Dashboard principal (`/dashboard`).

#### Subflujos de Sesión y Cuenta Propia
- **Renovación silenciosa de token (`POST /api/auth/refresh`):** Cuando el `accessToken` expira, el interceptor HTTP del frontend invoca `/api/auth/refresh` enviando automáticamente la cookie `centinela_refresh`. El backend valida que la sesión (`sid`) siga activa en PostgreSQL/Redis y que el `jti_refresh` coincida, rota ambos identificadores (`jti_access` y `jti_refresh`) *in-place* sobre la misma fila `sid` (evitando multiplicar sesiones al recargar pestañas) y retorna un nuevo `accessToken` junto con la cookie renovada.
- **Cierre de sesión (`POST /api/auth/logout`):** El usuario presiona «Cerrar sesión». El backend marca la sesión `sid` como inactiva en PostgreSQL, elimina `sesion:<sid>` de Redis, limpia la cookie `centinela_refresh` (`Max-Age=-1`), registra `LOGOUT` en `auditoria` y emite por el canal interno el evento SSE `event: cierre` con `{"motivo": "LOGOUT"}` para cortar de inmediato cualquier conexión de eventos abierta en esa sesión.
- **Actualización de perfil propio (`GET/PUT /api/account/profile`):** El usuario consulta sus datos y permisos asignados, o actualiza su `nombreCompleto`. Si solicita modificar su `emailUsuario` (`FIX-111`), debe suministrar su contraseña actual y un código TOTP válido, enviándose un enlace/token de confirmación al nuevo correo y un aviso de seguridad al correo anterior.

#### Flujos de Excepción y Salvaguardas
- **Credenciales inválidas, usuario inactivo o eliminado (`401 AUTH_FAILED`):** El backend responde con un único mensaje genérico («Credenciales inválidas») sin revelar si falló el correo, la contraseña o el estado de la cuenta, y audita el fallo.
- **Intento de regenerar QR con 2FA ya vinculado (`409 TWO_FACTOR_ALREADY_ENABLED`):** Si `totp_vinculado == true` y se invoca `GET /api/auth/2fa/qr`, el backend rechaza la petición para impedir que un atacante que conozca la contraseña robe el secreto TOTP. Solo un `ADMIN` puede resetear el 2FA (`CU-09`).
- **Código TOTP inválido o reutilizado (`401 TOTP_FAILED`):** Si el código de 6 dígitos no coincide o ya fue utilizado en la misma ventana de 30 segundos (`ultimo_totp_periodo`), se rechaza e incrementa el contador de fallos en Redis. Al alcanzar **3 intentos fallidos**, el backend destruye el `jwtTemporal` en Redis y obliga al usuario a volver a iniciar sesión con usuario y contraseña.
- **Exceso de intentos por fuerza bruta (`429 TOO_MANY_REQUESTS` — `FIX-118`):** El rate limiter en Redis bloquea temporalmente la IP/cuenta tras superar el umbral de intentos fallidos en `/auth/login` o `/auth/2fa/verify`.

---

### CU-02: Visualizar Dashboard General y Salud del Nodo Físico

- **Requerimientos asociados:** `RF-02`, `RNF-02`, `RNF-04`.
- **Actor principal:** Usuario Autenticado (`ADMIN` u `OPERATOR`).
- **Actores secundarios:** API de Proxmox VE, Redis (caché de estado del nodo).
- **Endpoints involucrados:** `GET /api/node/status`.

#### Descripción
El usuario consulta el estado de salud global del hipervisor (uso de CPU, memoria RAM, almacenamiento físico, tiempo de actividad, capacidad disponible para aprovisionamiento y resumen de instancias por estado), evaluado mediante un semáforo de criticidad.

#### Precondiciones
1. El usuario posee una sesión activa válida y `cambio_contrasena == false`.

#### Postcondiciones
- **De éxito:** El frontend renderiza los medidores de recursos, el semáforo de salud del nodo y los contadores de máquinas virtuales y contenedores.

#### Flujo Principal (Transacciones)
1. El usuario accede a la vista principal (`/dashboard`) y el frontend solicita `GET /api/node/status`.
2. El backend verifica en Redis si existe una instantánea fresca del estado del nodo (TTL de caché de 5 a 10 segundos para evitar saturar la API de Proxmox ante múltiples usuarios concurrentes).
3. Si la caché expiró, el backend consulta a la API de Proxmox VE los endpoints `/nodes/{node}/status` y `/cluster/resources`.
4. Proxmox VE retorna las métricas crudas de CPU, memoria física, disco del nodo, `uptime` y el listado de recursos.
5. El backend normaliza las métricas, calcula la capacidad asignable (`allocatable`, descontando el **15% de reserva de seguridad de RAM** para el sistema operativo del host) y determina el estado del **Semáforo de Salud (`D3`)**:
   - **`Saludable` (Verde):** Uso de CPU, RAM y Almacenamiento por debajo del `70%`.
   - **`Advertencia` (Amarillo):** Al menos un recurso se encuentra en `>= 70%` (o en nivel de saturación `>= 85% / 90%`).
6. **CONDICIONAL — Filtrado de superficie según el rol (`FIX-117`):**
   - Si el usuario es **`ADMIN`**: recibe la telemetría física completa del host, la capacidad de almacenamiento y el conteo global de instancias del servidor (`instancesSummary`).
   - Si el usuario es **`OPERATOR`**: el resumen de instancias se acota estrictamente a las instancias sobre las que posee permisos en `permisos_instancia`, ocultando datos sensibles de topología física del host.
7. El backend almacena el resultado en caché y responde `200 OK` al frontend.
8. El frontend actualiza las tarjetas de resumen, las barras de progreso y el indicador visual del semáforo, manteniendo la vista actualizada mediante polling periódico sin recargar la página.

#### Flujos de Excepción y Salvaguardas
- **Hipervisor Proxmox VE inaccesible o timeout (`502 PROXMOX_UNAVAILABLE` / `504 PROXMOX_TIMEOUT`):** Si Proxmox no responde dentro del timeout configurado, el backend retorna el código de error estandarizado (o el último estado conocido marcado como degradado) y el frontend cambia inmediatamente el semáforo global a estado **`Inaccesible` (Rojo)** mostrando un banner de degradación de servicio sin bloquear la navegación por el resto de la plataforma.

---

### CU-03: Consultar, Filtrar e Inventariar Instancias (VMs y LXC)

- **Requerimientos asociados:** `RF-03`, `RNF-01`, `RNF-02`.
- **Actor principal:** Usuario Autenticado (`ADMIN` u `OPERATOR`).
- **Actores secundarios:** API de Proxmox VE, PostgreSQL (`permisos_instancia`, `tareas_asincronas`).
- **Endpoints involucrados:** `GET /api/instances`, `GET /api/instances/:vmid`.

#### Descripción
El sistema presenta el inventario unificado de máquinas virtuales (`QEMU`) y contenedores (`LXC`) filtrado estrictamente según la matriz de permisos del usuario autenticado, enriquecido en vivo con su dirección IPv4 privada, su nivel de acceso (`FULL_ACCESS` o `READ_ONLY`) y el estado de cualquier tarea asíncrona en curso (`activeTask`).

#### Precondiciones
1. El usuario posee una sesión activa válida.

#### Postcondiciones
- **De éxito:** El usuario visualiza únicamente las instancias autorizadas, con controles de acción habilitados, deshabilitados o bloqueados por tarea en curso según corresponda a su `nivelAcceso` y al estado de cada máquina.

#### Flujo Principal (Transacciones)
1. El usuario navega al panel de inventario o al dashboard y el frontend solicita `GET /api/instances`.
2. El backend consulta en vivo a la API de Proxmox VE (`GET /cluster/resources?type=vm`) obteniendo el listado crudo de recursos `qemu` y `lxc`.
3. **Filtrado RBAC en el Dominio:**
   - Si el usuario tiene rol **`ADMIN`**: el backend conserva todas las instancias del nodo y les asigna `nivelAcceso: "FULL_ACCESS"`.
   - Si el usuario tiene rol **`OPERATOR`**: el backend consulta en PostgreSQL la tabla `permisos_instancia` para ese `usuario_id`, descarta toda instancia cuyo `vmid` no figure asignado e inyecta en cada instancia autorizada su nivel específico (`"FULL_ACCESS"` o `"READ_ONLY"`).
4. **Resolución Concurrente de IP en Vivo (`BAC-23B`):**
   - Únicamente sobre la lista ya filtrada por RBAC y solo para aquellas instancias cuyo estado sea `"running"`, el backend lanza consultas concurrentes acotadas (worker pool con timeout corto) hacia Proxmox VE: `GET /nodes/{node}/qemu/{vmid}/agent/network-get-interfaces` para VMs y `GET /nodes/{node}/lxc/{vmid}/interfaces` para contenedores LXC, extrayendo la primera dirección IPv4 privada no-loopback (`ip`). Si una instancia está detenida o no tiene QEMU Guest Agent activo, asigna `ip: null` sin hacer fallar el listado.
5. **Cruce con Tareas Asíncronas Activas:**
   - El backend consulta en PostgreSQL (`BuscarTareasActivasPorVmids`) si existen registros en `tareas_asincronas` con estado `RUNNING` para los `vmid` listados y adjunta el objeto `activeTask` (`tareaId`, `tipoAccion`, `estado`) en cada instancia correspondiente.
6. El backend unifica las propiedades (`vmid`, `name`, `type: "vm" | "lxc"`, `status`, `cpuUsage`, `maxCpu`, `ramUsage`, `maxRam`, `diskUsage`, `maxDisk`, `uptime`, `ip`, `nivelAcceso`, `activeTask`) y responde `200 OK`.
7. El frontend renderiza la tabla de inventario:
   - Si `nivelAcceso == "READ_ONLY"`, oculta o deshabilita con tooltip explicativo los botones de energía y modificación.
   - Si `activeTask != null` (incluso si el usuario acaba de recargar la página con `F5`), muestra el indicador de progreso en la fila y mantiene bloqueados los botones de esa instancia.
8. El usuario utiliza el buscador dinámico por nombre/ID o los filtros por tipo (`Todos`, `VM`, `LXC`) y estado (`Running`, `Stopped`, `Paused`), y el frontend filtra la vista de forma reactiva.

#### Flujos de Excepción y Salvaguardas
- **Consulta de detalle sin permiso (`403 INSTANCE_ACCESS_DENIED`):** Si un `OPERATOR` intenta acceder directamente por URL a `GET /api/instances/:vmid` de una instancia no asignada, el middleware `RequireInstanceAccess(READ_ONLY)` bloquea la petición y audita el acceso denegado.
- **Instancia inexistente en el hipervisor (`404 INSTANCE_NOT_FOUND`):** Si el `:vmid` solicitado no existe en Proxmox VE, el backend retorna `404`.

---

### CU-04: Ejecutar Acciones de Control de Energía y Ciclo de Vida (Manejo de UPID)

- **Requerimientos asociados:** `RF-04`, `RF-09`, `RF-11`, `RNF-03`, `RNF-04`.
- **Actor principal:** `ADMIN` u `OPERATOR` con permiso `FULL_ACCESS` sobre la instancia (para `delete`: exclusivamente `ADMIN`).
- **Actores secundarios:** API de Proxmox VE, Worker `SeguimientoTareas`, PostgreSQL (`tareas_asincronas`, `auditoria`, `permisos_instancia`), Redis Pub/Sub y Canal SSE.
- **Endpoints involucrados:**
  - `POST /api/instances/:vmid/status/:action` (acciones: `start`, `stop`, `shutdown`, `reboot`)
  - `DELETE /api/instances/:vmid` (exclusivo `ADMIN`)

#### Descripción
El usuario autorizado ejecuta una orden asíncrona sobre el ciclo de vida de una máquina virtual o contenedor: encendido (`start`), apagado ordenado por sistema operativo (`shutdown`), corte forzado de energía (`stop`), reinicio (`reboot`) o eliminación definitiva de la instancia (`delete`). El sistema valida permisos, protección de infraestructura, ausencia de tareas concurrentes y coherencia del estado actual antes de invocar a Proxmox VE, monitorea el `UPID` en segundo plano y notifica el desenlace en tiempo real vía Server-Sent Events (SSE).

#### Precondiciones
1. El usuario posee rol `ADMIN` o `OPERATOR` con `nivelAcceso = FULL_ACCESS` sobre el `:vmid` (o rol `ADMIN` si la acción es eliminar).
2. El `:vmid` **no** pertenece a la lista de instancias protegidas del sistema (`100..105`).
3. La instancia no tiene ninguna otra tarea asíncrona con estado `RUNNING` en `tareas_asincronas`.

#### Postcondiciones
- **De éxito:** La instancia transiciona al nuevo estado en Proxmox VE (o es eliminada junto con todos sus permisos asociados en BD si fue un `DELETE`), la tarea en `tareas_asincronas` queda en `COMPLETED`, se registran dos entradas en `auditoria` (`PENDING` al iniciar y `EXITO` al finalizar) y se emite el evento SSE `TASK_FINISHED` desbloqueando la interfaz.
- **De fallo:** Si falla una validación previa, no se invoca a Proxmox VE. Si la tarea falla durante su ejecución en Proxmox o supera los 3 minutos, `tareas_asincronas` pasa a `FAILED`/`TIMEOUT`, se registra `FALLA` en `auditoria` y se notifica el error por SSE.

#### Flujo Principal — Acciones de Energía (`start`, `shutdown`, `stop`, `reboot`)
1. El usuario presiona un botón de acción de energía sobre una instancia en la tabla o vista de detalle.
2. El frontend despliega un modal de confirmación explícita (advirtiendo riesgo de pérdida de datos en memoria si la acción elegida es `stop` —apagado forzado— en lugar de `shutdown`).
3. El usuario confirma la operación y el frontend envía `POST /api/instances/:vmid/status/:action`.
4. El backend ejecuta en orden la **cadena de guardas de seguridad**:
   1. `RequireInstanceAccess(FULL_ACCESS)`: verifica que sea `ADMIN` o tenga `FULL_ACCESS` sobre `:vmid`.
   2. `RejectProtectedInstance`: verifica que `:vmid` no esté en `100..105`.
   3. `RequireInstanceIdle`: verifica en PostgreSQL que `:vmid` no tenga una tarea en estado `RUNNING`.
   4. `validarEstadoParaAccion`: consulta a Proxmox VE el estado actual de la instancia y comprueba la compatibilidad:
      - `start` requiere estado actual `stopped`.
      - `stop`, `shutdown` y `reboot` requieren estado actual `running`.
5. Superadas las validaciones, el backend envía la orden a Proxmox VE (`POST /nodes/{node}/{type}/{vmid}/status/{action}`).
6. Proxmox VE encola la operación en el hipervisor y retorna inmediatamente un identificador `UPID`.
7. El backend registra la tarea en la tabla `tareas_asincronas` con estado `RUNNING` (guardando el `UPID` solo internamente), inserta el primer registro en `auditoria` con resultado `PENDING`, encola el seguimiento en el worker pool `SeguimientoTareas` y responde al frontend **`202 Accepted`** con `{ tareaId, estado: "RUNNING" }` (sin exponer el `UPID` crudo, según `FIX-105`).
8. Al recibir `202 Accepted`, el frontend bloquea todos los controles de acción de esa fila/instancia y muestra el indicador visual «Operación en progreso...».
9. En segundo plano, el worker `SeguimientoTareas` consulta periódicamente a Proxmox VE (`GET /nodes/{node}/tasks/{upid}/status`) hasta que la tarea pasa a `stopped` (con un límite máximo de **3 minutos** según `D4`).
10. Cuando Proxmox VE reporta la finalización con `exitstatus == "OK"`, el worker actualiza `tareas_asincronas` a `COMPLETED`, inserta el segundo registro en `auditoria` con resultado `EXITO` y publica el evento `TASK_FINISHED` en Redis Pub/Sub.
11. El servicio SSE transmite el evento `TASK_FINISHED` a las sesiones activas de los usuarios con permiso sobre ese `vmid`.
12. El frontend recibe el evento SSE, desbloquea los controles de la instancia, refresca su estado en vivo y despliega una notificación flotante (toast) confirmando el éxito de la operación.

#### Subflujo Alternativo — Eliminación Definitiva de Instancia (`DELETE /api/instances/:vmid`)
1. Un usuario con rol **`ADMIN`** solicita eliminar una instancia detenida.
2. El frontend abre un **modal de confirmación destructiva reforzada (`FRN-16B`)** que exige al administrador **escribir manualmente el nombre exacto de la instancia** para habilitar el botón «Eliminar definitivamente».
3. El administrador escribe el nombre, confirma y el frontend invoca `DELETE /api/instances/:vmid`.
4. El backend valida `RequireRole("ADMIN")`, `RejectProtectedInstance` (`100..105`), `RequireInstanceIdle` y consulta a Proxmox VE para comprobar que la instancia esté en estado **`stopped`**.
5. El backend invoca `DELETE /nodes/{node}/{type}/{vmid}` en Proxmox VE, monitorea el `UPID` y, al confirmarse el borrado, ejecuta la **purga automática de permisos huérfanos (`BAC-24C` / `FIX-113`)**: elimina todas las filas de `permisos_instancia` asociadas a ese `vmid` para evitar que un futuro reciclaje del ID en Proxmox otorgue acceso accidental a operadores anteriores, registrando `ELIMINAR_VM` e `INSTANCE_PERMISSIONS_PURGED` en `auditoria`.

#### Flujos de Excepción y Salvaguardas
- **Permiso insuficiente (`403 INSTANCE_ACCESS_DENIED`):** Si un `OPERATOR` con `READ_ONLY` o sin asignación intenta ejecutar una acción de energía, o si cualquier `OPERATOR` intenta un `DELETE`.
- **Instancia protegida del sistema (`403 INSTANCE_PROTECTED`):** Si el `:vmid` pertenece a `100..105`. Se audita el intento bloqueado.
- **Operación concurrente en curso (`409 INSTANCE_BUSY`):** Si ya hay un `UPID` activo en `tareas_asincronas` para ese `:vmid`.
- **Transición de estado inválida (`409 INSTANCE_INVALID_STATE` / `409 INSTANCE_NOT_STOPPED`):** Si se intenta encender una instancia que ya está `running`, apagar/reiniciar una que ya está `stopped`, o eliminar (`DELETE`) una instancia que no está `stopped`.
- **Acción no soportada (`400 INVALID_ACTION`):** Si se envía una acción fuera del conjunto `{start, stop, shutdown, reboot}` (por ejemplo, `pause`).
- **Timeout de seguimiento de UPID (`D4` — 3 minutos):** Si la tarea en Proxmox supera los 3 minutos sin concluir, el worker la marca como `TIMEOUT` en `tareas_asincronas`, registra `FALLA` en `auditoria` y emite `TASK_FINISHED` con `status: "FAILED"` para destrabar la interfaz del usuario.
- **Diferenciación entre Caída y Timeout de Proxmox (`502 PROXMOX_UNAVAILABLE` vs `504 PROXMOX_TIMEOUT`):** Si Proxmox rechaza la conexión (`502`), el frontend informa que la orden **no** llegó a ejecutarse y es seguro reintentar; si vence el timeout de red (`504`), advierte que el estado es ambiguo y refresca el inventario antes de habilitar reintentos.
- **Recuperación tras reinicio del Backend (`BAC-25C`):** Al iniciar el servidor Go, el sistema busca en `tareas_asincronas` todas las filas que hayan quedado en `RUNNING` antes del reinicio y reanuda automáticamente su monitoreo en el worker pool.

---

### CU-05: Consultar Telemetría en Tiempo Real e Histórico de Métricas por Instancia

- **Requerimientos asociados:** `RF-05`, `RNF-02`, `RNF-06`.
- **Actor principal:** `ADMIN` u `OPERATOR` con permiso `READ_ONLY` o `FULL_ACCESS` sobre la instancia.
- **Actores secundarios:** API de Proxmox VE, Redis (tickets WebSocket y caché RRDdata).
- **Endpoints involucrados:**
  - `POST /api/auth/ws-ticket`
  - `GET /api/ws/instances/:vmid/metrics?ticket=<uuid>` (WebSocket)
  - `GET /api/instances/:vmid/metrics?timeframe=2h|24h` (REST histórico)

#### Descripción
El usuario visualiza en el panel de detalle de una instancia los gráficos dinámicos de consumo de CPU (%), memoria RAM (usada vs asignada), tráfico de red (`netin` / `netout`) y lectura/escritura de disco (`diskread` / `diskwrite`), combinando un flujo en vivo cada 2 segundos vía WebSocket y series históricas recientes (`2h`, `24h`) obtenidas desde RRDdata de Proxmox VE.

#### Precondiciones
1. El usuario posee una sesión activa y al menos permiso `READ_ONLY` sobre el `:vmid`.

#### Postcondiciones
- **De éxito:** Se establecen y actualizan los gráficos en tiempo real e histórico; al abandonar la vista, la conexión WebSocket y el ticker del backend se cierran limpiamente liberando recursos.

#### Flujo Principal (Transacciones)
1. El usuario ingresa a la vista de detalle de una instancia (`/instances/:vmid`) y abre la sección de Telemetría y Rendimiento.
2. **Carga de Serie Histórica Inicial:**
   1. El frontend solicita `GET /api/instances/:vmid/metrics?timeframe=2h`.
   2. El backend valida `RequireInstanceAccess(READ_ONLY)` y verifica si la serie RRD para ese `:vmid` y `timeframe` está en la caché de Redis (TTL 30 segundos).
   3. Si no está en caché, consulta `GET /nodes/{node}/{type}/{vmid}/rrddata` en Proxmox VE, normaliza los timestamps a formato Unix UTC y reemplaza valores nulos por `0`, guardando el resultado en Redis y devolviéndolo al frontend.
   4. El frontend renderiza los gráficos de área/líneas iniciales con escalado dinámico de unidades (`MB`/`GB`, `KB/s`/`MB/s`).
3. **Suscripción al Flujo en Vivo por WebSocket con Ticket Efímero (`SEC-09`):**
   1. Para no exponer el JWT en la URL del WebSocket, el frontend solicita primero `POST /api/auth/ws-ticket` con su header `Authorization: Bearer <accessToken>`.
   2. El backend valida el permiso sobre la instancia, genera un `ticket` UUID criptográfico en Redis (`ws:ticket:<uuid>`, TTL 30 segundos, un solo uso) y lo devuelve al frontend.
   3. El frontend abre la conexión WebSocket hacia `GET /api/ws/instances/:vmid/metrics?ticket=<uuid>`.
   4. El backend consume y elimina atómicamente el ticket en Redis (`GETDEL`), actualiza la conexión HTTP a WebSocket e inicia un ticker cada **2 segundos** que consulta el estado actual de la instancia en Proxmox VE y emite frames JSON con `{ timestamp, status, cpuPercent, ramUsedBytes, ramMaxBytes, netInBytes, netOutBytes, diskReadBytes, diskWriteBytes }`.
4. El frontend agrega cada nuevo punto recibido a una ventana deslizante en memoria (descartando los puntos más antiguos para no degradar la RAM del navegador) y actualiza los gráficos suavemente.
5. Si el usuario cambia el selector de rango histórico (por ejemplo, de `2h` a `24h`), el frontend invoca `GET /api/instances/:vmid/metrics?timeframe=24h` y actualiza la gráfica de tendencias.
6. Cuando el usuario cierra la vista o navega a otra pantalla, el frontend cierra el WebSocket y el backend detiene inmediatamente la goroutine del ticker asociada.

#### Flujos de Excepción y Salvaguardas
- **Instancia detenida (`status == "stopped"`):** Si la máquina está apagada o se apaga mientras el usuario mira los gráficos, el backend **no** cierra abruptamente el socket con error; emite frames con `status: "stopped"` y métricas en `0`, y el frontend muestra un overlay informativo «Instancia detenida — sin telemetría activa».
- **Ticket WebSocket expirado, inválido o reutilizado (`401 UNAUTHORIZED`):** Si el ticket supera los 30 segundos de TTL o intenta usarse por segunda vez, el handshake WebSocket es rechazado antes de hacer el upgrade.

---

### CU-06: Administrar Snapshots (Puntos de Restauración)

- **Requerimientos asociados:** `RF-06`, `RNF-03`, `RNF-04`.
- **Actor principal:**
  - Para listar: `ADMIN` u `OPERATOR` con `READ_ONLY` o `FULL_ACCESS`.
  - Para crear, revertir (rollback) o eliminar: `ADMIN` u `OPERATOR` con `FULL_ACCESS`.
- **Actores secundarios:** API de Proxmox VE, Worker `SeguimientoTareas`, PostgreSQL (`tareas_asincronas`, `auditoria`), Redis Pub/Sub y Canal SSE.
- **Endpoints involucrados:**
  - `GET /api/instances/:vmid/snapshots`
  - `POST /api/instances/:vmid/snapshots`
  - `POST /api/instances/:vmid/snapshots/:snapname/rollback`
  - `DELETE /api/instances/:vmid/snapshots/:snapname`

#### Descripción
El usuario gestiona los puntos de restauración (*snapshots*) de una instancia: consulta el historial existente, crea un nuevo snapshot antes de un cambio crítico, revierte la instancia a un estado anterior (*rollback*) bajo confirmación destructiva reforzada o elimina snapshots obsoletos. Todas las operaciones de mutación se ejecutan de manera asíncrona monitoreando su `UPID`.

#### Precondiciones
1. El usuario posee permiso sobre el `:vmid` (`READ_ONLY` para listar; `FULL_ACCESS` o `ADMIN` para crear, revertir o eliminar).
2. Para crear, revertir o eliminar: el `:vmid` no es una instancia protegida (`100..105`) y no tiene otra tarea activa en curso (`RequireInstanceIdle`).
3. El almacenamiento subyacente de la instancia en Proxmox VE soporta snapshots (ej. LVM-Thin, ZFS, Ceph, qcow2).

#### Postcondiciones
- **De éxito:** El snapshot es creado, restaurado o eliminado en Proxmox VE; se registra el inicio (`PENDING`) y fin (`EXITO`) en `auditoria`; y se notifica a la interfaz mediante el evento SSE `TASK_FINISHED`.

#### Flujo Principal (Transacciones)
1. **Listado de Snapshots:**
   1. El usuario abre la pestaña «Snapshots» dentro del detalle de la instancia (`/instances/:vmid`).
   2. El frontend solicita `GET /api/instances/:vmid/snapshots`.
   3. El backend valida `RequireInstanceAccess(READ_ONLY)`, consulta `GET /nodes/{node}/{type}/{vmid}/snapshot` en Proxmox VE, filtra/identifica el nodo simbólico `current` y retorna la lista ordenada cronológicamente con `{ name, description, snaptime, vmstate, parent }`.
   4. El frontend muestra la tabla/árbol de puntos de restauración (ocultando los botones de crear, restaurar y borrar si el operador solo tiene `READ_ONLY`).
2. **Creación de un Nuevo Snapshot:**
   1. El usuario con `FULL_ACCESS` o `ADMIN` hace clic en «Crear Snapshot» y completa el formulario indicando `snapname` (nombre identificador), `description` y opcionalmente si desea incluir el estado de la memoria RAM (`vmstate`, aplicable a VMs en ejecución).
   2. El backend valida `RequireInstanceAccess(FULL_ACCESS)`, `RejectProtectedInstance`, `RequireInstanceIdle` y comprueba mediante expresión regular que `snapname` contenga únicamente caracteres alfanuméricos, guiones o guiones bajos (comenzando por una letra y con longitud máxima de 40 caracteres) y sanitiza la descripción (`SEC-12` / `FIX-116`).
   3. El backend envía `POST /nodes/{node}/{type}/{vmid}/snapshot` a Proxmox VE, recibe el `UPID`, registra la tarea en `tareas_asincronas` (`RUNNING`) y en `auditoria` (`CREATE_SNAPSHOT`, `PENDING`), encola el `UPID` en `SeguimientoTareas` y responde `202 Accepted` con `{ tareaId }`.
   4. Al finalizar el `UPID`, el worker audita `EXITO` y emite `TASK_FINISHED` por SSE; el frontend desbloquea la vista y recarga la lista de snapshots.
3. **Reversión (Rollback) a un Snapshot Existente:**
   1. El usuario selecciona «Revertir» sobre un snapshot de la tabla.
   2. El frontend abre un **modal de advertencia destructiva (`FRN-30`)** explicando que todo cambio posterior a la fecha del snapshot se perderá irreversiblemente, y exige al usuario **escribir manualmente el nombre exacto del snapshot (`snapname`)** para habilitar el botón de confirmación.
   3. El usuario escribe el nombre, confirma y el frontend invoca `POST /api/instances/:vmid/snapshots/:snapname/rollback`.
   4. El backend valida `RequireInstanceAccess(FULL_ACCESS)`, `RejectProtectedInstance` y `RequireInstanceIdle`, envía `POST /nodes/{node}/{type}/{vmid}/snapshot/{snapname}/rollback` a Proxmox VE y captura el `UPID`.
   5. El backend registra la tarea como `ROLLING_BACK` / `RUNNING`, audita `ROLLBACK_SNAPSHOT` (`PENDING`) y retorna `202 Accepted`.
   6. El frontend coloca la instancia en estado visual «Revirtiendo...» bloqueando todos sus controles hasta recibir el evento SSE `TASK_FINISHED`.
4. **Eliminación de un Snapshot:**
   1. El usuario solicita eliminar un snapshot antiguo y confirma en el modal.
   2. El frontend envía `DELETE /api/instances/:vmid/snapshots/:snapname`.
   3. El backend valida permisos e inactividad, dispara `DELETE /nodes/{node}/{type}/{vmid}/snapshot/{snapname}` en Proxmox VE, monitorea el `UPID`, audita `DELETE_SNAPSHOT` y notifica por SSE al concluir.

#### Flujos de Excepción y Salvaguardas
- **Nombre de snapshot inválido o duplicado (`400 INVALID_SNAPSHOT_NAME` / `409 SNAPSHOT_ALREADY_EXISTS`):** Si `snapname` contiene espacios, acentos o caracteres especiales no admitidos por Proxmox VE, o si ya existe un snapshot con ese identificador en la instancia.
- **Almacenamiento sin soporte de snapshots (`422 SNAPSHOT_NOT_SUPPORTED`):** Si el disco de la instancia reside en un storage que no admite snapshots nativos, el backend traduce el error técnico de Proxmox a un mensaje claro para el usuario.
- **Instancia ocupada (`409 INSTANCE_BUSY`):** Si se intenta crear, revertir o borrar un snapshot mientras otra tarea asíncrona está ejecutándose sobre la misma instancia.

---

### CU-07: Aprovisionar Nueva Instancia mediante Asistente Guiado (Wizard)

- **Requerimientos asociados:** `RF-07`, `RNF-01`, `RNF-03`, `RNF-04`.
- **Actor principal:** Usuario Administrador (`ADMIN` exclusivamente).
- **Actores secundarios:** API de Proxmox VE, Redis (Lock distribuido `lock:instance:provisioning`), Worker `SeguimientoTareas`, PostgreSQL (`tareas_asincronas`, `permisos_instancia`, `auditoria`), Canal SSE.
- **Endpoints involucrados:**
  - `GET /api/node/storages`
  - `GET /api/node/templates`
  - `POST /api/instances`

#### Descripción
El administrador despliega una nueva máquina virtual (`QEMU`) o contenedor (`LXC`) a través de un asistente paso a paso que abstrae los parámetros complejos del hipervisor, muestra en tiempo real la cuota de recursos disponible en el servidor y previene condiciones de carrera mediante un bloqueo distribuido y una validación estricta de capacidad física antes de enviar la orden a Proxmox VE.

#### Precondiciones
1. El usuario posee rol **`ADMIN`** (`RequireRole("ADMIN")`).
2. Existen plantillas LXC o imágenes/plantillas QEMU disponibles en los almacenamientos activos del nodo Proxmox VE.

#### Postcondiciones
- **De éxito:** Se aprovisiona la nueva instancia en Proxmox VE con el próximo `VMID` libre del clúster, se asocia automáticamente permiso `FULL_ACCESS` en `permisos_instancia`, se registra la operación en `auditoria` (`CREATE_INSTANCE`, `EXITO`) y la nueva instancia aparece en el inventario al recibir `TASK_FINISHED`.
- **De fallo:** Si la cuota de RAM o disco es insuficiente, la petición se rechaza con `409 RESOURCE_QUOTA_EXCEEDED` sin crear recursos parciales en Proxmox VE.

#### Flujo Principal (Transacciones)
1. El administrador abre el «Asistente de Nueva Instancia» (Wizard).
2. El frontend consulta en paralelo `GET /api/node/storages`, `GET /api/node/templates` y `GET /api/node/status` para conocer los almacenamientos válidos para discos (`images` / `rootdir`), el catálogo de plantillas/ISOs y la cuota libre asignable del nodo (`allocatable`).
3. El administrador completa los **4 pasos guiados del Wizard**:
   - **Paso 1 (Tipo y Base):** Selecciona Contenedor (`LXC`) o Máquina Virtual (`QEMU`), define el `hostname`/`name` y elige la plantilla o imagen base de la lista descubierta.
   - **Paso 2 (Recursos de Cómputo y Almacenamiento):** Selecciona el `storage` destino y ajusta los deslizadores de `vCPUs`, `RAM (MB)` y `Disco (GB)` viendo en vivo cuánto porcentaje de la cuota libre del nodo consumirá la nueva instancia.
   - **Paso 3 (Red):** Define la configuración de red sobre el bridge autorizado (DHCP o dirección IPv4 estática con máscara CIDR y puerta de enlace).
   - **Paso 4 (Resumen y Confirmación):** Revisa la ficha completa de la máquina a crear y confirma el aprovisionamiento.
4. El frontend bloquea el botón de envío y ejecuta `POST /api/instances` con el payload del asistente.
5. El backend ejecuta las **validaciones de seguridad y control de concurrencia**:
   1. **Sanitización Anti-Inyección de Parámetros (`SEC-10` / `FIX-116`):** Verifica que `hostname`, `storage`, `template` y los parámetros de red cumplan expresiones regulares estrictas (prohibiendo comas `,`, espacios o caracteres de control que pudieran inyectar directivas adicionales en los strings de configuración de Proxmox VE).
   2. **Adquisición de Lock Distribuido en Redis (`BAC-37`):** Adquiere el candado `lock:instance:provisioning` (`SETNX` con TTL de 30 segundos) para serializar las solicitudes de creación y evitar que dos administradores consuman al mismo tiempo la misma RAM libre o colisionen al pedir el próximo `VMID`.
   3. **Validación de Cuota Real del Host con Reserva del 15% (`BAC-36`):** Consulta el estado real del nodo y de los storages en Proxmox VE y verifica que:
      $$\text{RAM}_{\text{solicitada}} \le (\text{RAM}_{\text{total\_host}} \times 0.85) - \text{RAM}_{\text{asignada\_instancias}}$$
      y que el `storage` elegido posea espacio libre superior al disco solicitado.
6. Confirmada la disponibilidad de recursos, el backend solicita a Proxmox VE el próximo identificador libre (`GET /cluster/nextid`), construye los parámetros nativos e invoca la creación (`POST /nodes/{node}/lxc` o `POST /nodes/{node}/qemu`).
7. Proxmox VE inicia el aprovisionamiento y retorna el `UPID`.
8. El backend persiste la tarea en `tareas_asincronas` (`RUNNING`), registra `CREATE_INSTANCE` (`PENDING`) en `auditoria`, **libera el lock distribuido en Redis**, delega el monitoreo del `UPID` al worker `SeguimientoTareas` y responde al frontend **`202 Accepted`** con `{ tareaId, vmid }`.
9. Cuando el worker `SeguimientoTareas` verifica que el `UPID` finalizó con `exitstatus == "OK"`, registra automáticamente el permiso `FULL_ACCESS` para ese nuevo `vmid` en `permisos_instancia`, graba `CREATE_INSTANCE` (`EXITO`) en `auditoria` y publica `TASK_FINISHED` por SSE.
10. El frontend recibe la notificación en tiempo real, informa que la instancia está lista y redirige al inventario o al detalle de la nueva máquina.

#### Flujos de Excepción y Salvaguardas
- **Cuota de recursos insuficiente en el nodo (`409 RESOURCE_QUOTA_EXCEEDED`):** Si la RAM o el disco solicitados superan la capacidad asignable del host (respetando el 15% de reserva de RAM para el hipervisor), el backend libera el lock de Redis, rechaza la solicitud con `409` detallando el recurso agotado y el frontend resalta en rojo el paso de recursos en el Wizard sin perder los datos cargados.
- **Conflicto por aprovisionamiento simultáneo (`409 PROVISIONING_LOCK_BUSY`):** Si otro aprovisionamiento está reteniendo el lock de Redis y no se libera dentro de la ventana de espera.
- **Intento de acceso por un Operador (`403 FORBIDDEN`):** Cualquier intento de un `OPERATOR` de invocar `POST /api/instances` es bloqueado por `RequireRole("ADMIN")`.

---

### CU-08: Consultar y Exportar Registros Inmutables de Auditoría

- **Requerimientos asociados:** `RF-08`, `RNF-01`, `RNF-02`.
- **Actor principal:** Usuario Administrador (`ADMIN` exclusivamente).
- **Actores secundarios:** PostgreSQL (tabla particionada e inmutable `auditoria`).
- **Endpoints involucrados:**
  - `GET /api/admin/audit`
  - `GET /api/admin/audit/export?formato=csv|json`

#### Descripción
El administrador inspecciona el historial cronológico e inalterable de todas las operaciones de seguridad e infraestructura ocurridas en la plataforma, aplicando filtros combinados, inspeccionando el payload JSON de cada evento y exportando los resultados en formato CSV o JSON para cumplimiento y análisis forense.

#### Precondiciones
1. El usuario posee rol **`ADMIN`** (`RequireRole("ADMIN")`).

#### Postcondiciones
- **De éxito:** Se presentan los registros paginados o se descarga el archivo exportado, quedando registrada la propia exportación en la tabla `auditoria` (`EXPORTAR_AUDITORIA`).

#### Flujo Principal (Transacciones)
1. El administrador navega a la vista de «Auditoría» (`/audit`).
2. El frontend solicita `GET /api/admin/audit?pagina=1&tamano=20&ordenPor=fechaHora&ordenDir=DESC`.
3. El backend valida `RequireRole("ADMIN")` y ejecuta la consulta indexada sobre las particiones correspondientes de la tabla `auditoria` en PostgreSQL, retornando los eventos junto con los metadatos de paginación (`totalRegistros`, `totalPaginas`).
4. El frontend renderiza la tabla mostrando fecha y hora exacta, usuario ejecutor (nombre, email y rol), acción (`LOGIN`, `START`, `STOP`, `SHUTDOWN`, `REBOOT`, `ELIMINAR_VM`, `CREATE_INSTANCE`, `UPDATE_INSTANCE_CONFIG`, `CREATE_SNAPSHOT`, `ROLLBACK_SNAPSHOT`, `CREAR_USUARIO`, `ACTUALIZAR_PERMISOS`, etc.), instancia afectada (`instanciaId`), código de error (si hubo) y badge de resultado (`EXITO`, `FALLA`, `PENDING`).
5. El administrador aplica filtros por usuario, tipo de acción, ID de instancia, resultado o rango de fechas (`desde` / `hasta`) y el frontend refresca la consulta paginada.
6. Al hacer clic sobre una fila, el frontend abre un panel lateral (Drawer) que formatea el campo JSONB `detalles` (mostrando valores anteriores vs nuevos en cambios de configuración o el motivo técnico de un fallo en Proxmox).
7. **Exportación de Registros (`BAC-45` / `FRN-32`):**
   1. El administrador presiona «Exportar CSV» o «Exportar JSON».
   2. El frontend invoca `GET /api/admin/audit/export?formato=csv|json` conservando exactamente los mismos filtros activos en pantalla.
   3. El backend ejecuta la consulta en **modo streaming por cursor** (escribiendo fila a fila en el `http.ResponseWriter` con encabezado `Content-Disposition: attachment` para no agotar la memoria RAM del servidor ante miles de registros) y registra el evento `EXPORTAR_AUDITORIA` en la base de datos.

#### Salvaguardas de Inmutabilidad
- **Protección a Nivel Motor PostgreSQL (`INF-15` / `SEC-11`):** No existe en el backend ningún endpoint ni método de repositorio para editar o borrar filas de `auditoria`. Además, a nivel de base de datos, el rol de la aplicación tiene revocados los permisos `UPDATE`, `DELETE` y `TRUNCATE` sobre `auditoria`, y el trigger `trg_auditoria_inmutable` aborta con excepción cualquier intento de mutación directa.

---

### CU-09: Gestionar Usuarios, Roles, Ciclo de Vida de Cuentas y Permisos por Instancia

- **Requerimientos asociados:** `RF-01`, `RF-09`, `RNF-01`.
- **Actor principal:** Usuario Administrador (`ADMIN` exclusivamente).
- **Actores secundarios:** Servicio de Correo (`EmailService` SMTP), API de Proxmox VE, PostgreSQL (`usuarios`, `permisos_instancia`, `sesiones_activas`, `auditoria`), Redis y Canal SSE.
- **Endpoints involucrados:**
  - `GET /api/admin/users`, `GET /api/admin/users/:id`, `GET /api/roles`
  - `POST /api/admin/users`
  - `PUT /api/admin/users/:id`
  - `DELETE /api/admin/users/:id`
  - `GET /api/admin/users/:id/permissions` y `PUT /api/admin/users/:id/permissions`
  - `POST /api/admin/users/:id/password/reset`
  - `POST /api/admin/users/:id/2fa/reset`

#### Descripción
El administrador gestiona el ciclo de vida completo de las cuentas del personal: da de alta nuevos usuarios (con generación y envío de contraseña temporal por correo), edita sus datos o rol (`ADMIN` / `OPERATOR`), asigna o revoca permisos granulares por instancia (`FULL_ACCESS` vs `READ_ONLY`), ejecuta reseteos administrativos de contraseña o 2FA en caso de pérdida de dispositivo, suspende cuentas temporalmente o aplica la baja lógica definitiva.

#### Precondiciones
1. El usuario ejecutor posee rol **`ADMIN`**.

#### Postcondiciones
- **De éxito:** Los cambios de usuario o matriz de permisos persisten transaccionalmente en PostgreSQL; ante cualquier suspensión, baja o reseteo de credenciales/2FA, todas las sesiones activas y conexiones SSE del usuario afectado son revocadas en el acto; todas las acciones quedan registradas en `auditoria`.

#### Flujo Principal y Subflujos Administrativos

##### 9.1. Alta de Nuevo Usuario con Contraseña Temporal (`POST /api/admin/users`)
1. El administrador accede a la sección «Usuarios» (`/users`), presiona «Nuevo Usuario» y completa `nombreCompleto`, `nombreUsuario`, `emailUsuario` y `rol` (`ADMIN` o `OPERATOR`).
2. El backend valida que el `emailUsuario` y el `nombreUsuario` no estén en uso por ningún usuario no eliminado (`WHERE eliminado_en IS NULL`).
3. El backend genera aleatoriamente una **contraseña temporal segura** que cumple la política de complejidad, calcula su hash bcrypt (costo 12) y prepara la creación del usuario con `cambio_contrasena = true`, `totp_vinculado = false` y `activo = true`.
4. El backend envía la contraseña temporal e instrucciones de acceso al correo del nuevo usuario mediante `EmailService`.
5. **Salvaguarda Transaccional de Entrega de Correo (`D5`):** Si el servidor SMTP falla al enviar el correo, el backend revierte la operación y responde **`502 Bad Gateway` (`EMAIL_DELIVERY_FAILED`)**, evitando dejar cuentas creadas cuya contraseña temporal nadie conoce.
6. Si el correo se envía con éxito, confirma el alta en PostgreSQL, registra `CREAR_USUARIO` (`EXITO`) en `auditoria` y retorna `201 Created`.

##### 9.2. Asignación Granular de Permisos por Instancia (`PUT /api/admin/users/:id/permissions`)
1. El administrador selecciona un usuario con rol `OPERATOR` y abre el modal «Gestionar Permisos de Instancias».
2. El frontend consulta el inventario disponible y los permisos actuales (`GET /api/admin/users/:id/permissions`).
3. El administrador selecciona a qué instancias (`vmid`) tendrá acceso el operador y define para cada una el nivel de acceso:
   - **`READ_ONLY`:** Solo visualización de estado, detalle, métricas y listado de snapshots.
   - **`FULL_ACCESS`:** Visualización más control de energía (`start`, `stop`, `shutdown`, `reboot`), edición de recursos (`vCPU`/`RAM`) y gestión de snapshots.
4. El frontend envía `PUT /api/admin/users/:id/permissions` con `{"permisos": [{"vmid": 106, "nivelAcceso": "FULL_ACCESS"}, ...]}`.
5. **Salvaguarda de Instancias Protegidas (`FIX-119`):** El backend verifica que ningún `vmid` de la lista pertenezca a `PROXMOX_PROTECTED_VMIDS` (`100..105`); si incluye alguno, rechaza toda la petición con **`403 Forbidden` (`INSTANCE_PROTECTED`)**.
6. El backend reemplaza de forma atómica (dentro de una transacción SQL) las filas de `permisos_instancia` para ese usuario, registra `ACTUALIZAR_PERMISOS` en `auditoria` y responde `200 OK`.

##### 9.3. Reseteo Administrativo de 2FA o Contraseña
- **Reseteo de 2FA (`POST /api/admin/users/:id/2fa/reset`):** Cuando un usuario pierde su dispositivo autenticador, el administrador ejecuta el reseteo de 2FA. El backend limpia `secreto_totp = ''`, marca `totp_vinculado = false`, **revoca todas las sesiones activas del usuario** en PostgreSQL y Redis (cortando su SSE con motivo `SESSIONS_REVOKED`) y audita `RESETEAR_2FA`. En su próximo ingreso (`CU-01`), el sistema le presentará un nuevo código QR.
- **Reseteo de Contraseña (`POST /api/admin/users/:id/password/reset`):** El backend genera una nueva contraseña temporal, la envía por correo (`502 EMAIL_DELIVERY_FAILED` si falla SMTP), actualiza `hash_contrasena`, fuerza `cambio_contrasena = true`, **revoca todas las sesiones activas del usuario** y audita `RESETEAR_CONTRASENA`.

##### 9.4. Suspensión Reversible vs. Eliminación Lógica Definitiva (`FIX-54`, `FIX-55`, `FIX-56`)
- **Suspensión temporal (`PUT /api/admin/users/:id` con `activo: false`):** Deshabilita el acceso del usuario de manera **reversible**. El backend marca `activo = false`, revoca inmediatamente todas sus sesiones activas y cierra su canal SSE emitiendo `event: cierre` (`{"motivo": "USER_INACTIVE"}`). El usuario puede ser reactivado en el futuro (`activo: true`).
- **Eliminación definitiva (`DELETE /api/admin/users/:id`):** Ejecuta una **baja lógica irreversible** (`eliminado_en = NOW()`, `activo = false`).
  - **Salvaguarda de Integridad de Auditoría:** Nunca se ejecuta un `DELETE FROM usuarios` físico en PostgreSQL porque violaría la clave foránea de los registros históricos en la tabla inmutable `auditoria`. En su lugar, el índice único parcial (`WHERE eliminado_en IS NULL`) libera inmediatamente el `email` y el `nombre_usuario` por si se necesita volver a registrar a la persona desde cero, a la vez que limpia sus `permisos_instancia` y revoca todas sus sesiones.

#### Flujos de Excepción y Salvaguardas
- **Prohibición de Auto-Sabotaje (`400 SELF_DELETE_NOT_ALLOWED`):** Un administrador no puede eliminarse a sí mismo, ni suspender su propia cuenta (`activo: false`), ni degradarse su propio rol de `ADMIN` a `OPERATOR`.
- **Protección contra Escalada Horizontal entre Administradores (`403 FORBIDDEN` — `FIX-112`):** Un `ADMIN` estándar no puede eliminar, desactivar ni resetear las credenciales/2FA de otro par `ADMIN` (reservado a la cuenta `SUPERADMIN` inicial).
- **Conflicto de Unicidad (`409 USER_CONFLICT`):** Si al crear o editar un usuario se ingresa un correo o nombre de usuario que ya pertenece a otra cuenta activa/suspendida no eliminada.

---

### CU-10: Modificar Recursos de Hardware de una Instancia (Escalado Vertical)

- **Requerimientos asociados:** `RF-10`, `RNF-03`, `RNF-04`.
- **Actor principal:** `ADMIN` u `OPERATOR` con permiso `FULL_ACCESS` sobre la instancia.
- **Actores secundarios:** API de Proxmox VE, PostgreSQL (`auditoria`, `tareas_asincronas`).
- **Endpoints involucrados:** `PUT /api/instances/:vmid/config`.

#### Descripción
El usuario con control operativo ajusta la cantidad de núcleos virtuales (`vCPUs`) y/o memoria `RAM (MB)` asignados a una instancia existente. El sistema calcula el diferencial de recursos ($\Delta\text{RAM}$), verifica que cualquier incremento respete la cuota libre del servidor físico y aplica el cambio en Proxmox VE informando si requiere reinicio para hacerse efectivo.

#### Precondiciones
1. El usuario posee rol `ADMIN` o permiso `FULL_ACCESS` sobre el `:vmid`.
2. El `:vmid` no es una instancia protegida (`100..105`) y no tiene tareas asíncronas en curso (`RequireInstanceIdle`).

#### Postcondiciones
- **De éxito:** La configuración de CPU/RAM se actualiza en Proxmox VE, se registra el cambio en `auditoria` (`UPDATE_INSTANCE_CONFIG` detallando valores previos y nuevos) y la interfaz informa si los cambios se aplicaron en caliente o si quedan pendientes de reinicio (`requiresReboot`).

#### Flujo Principal (Transacciones)
1. El usuario abre el modal «Modificar Recursos» dentro de la vista de detalle de la instancia (`/instances/:vmid`).
2. El frontend precarga los controles numéricos con la asignación actual de la instancia (`maxCpu` y `maxRam`) y muestra la capacidad libre del nodo.
3. El usuario modifica los valores de `cores` (vCPUs) y/o `memoryMB` (RAM) y confirma la operación.
4. El frontend envía `PUT /api/instances/:vmid/config` con `{ cores, memoryMB }`.
5. El backend ejecuta `RequireInstanceAccess(FULL_ACCESS)`, `RejectProtectedInstance` y `RequireInstanceIdle`, y valida los límites formales (`cores >= 1` sin superar los núcleos físicos del host; `memoryMB >= 128`).
6. El backend consulta en Proxmox VE la configuración actual de la instancia y calcula el **diferencial de memoria**:
   $$\Delta\text{RAM} = \text{memoryMB}_{\text{nuevo}} - \text{memoryMB}_{\text{actual}}$$
7. **CONDICIONAL — Evaluación Inteligente de Cuota (`BAC-39`):**
   - **7A. Incremento de recursos ($\Delta\text{RAM} > 0$):** El backend consulta la memoria libre asignable del nodo físico (respetando el 15% de reserva de seguridad del host). Si $\Delta\text{RAM}$ supera la memoria disponible, aborta con **`409 Conflict` (`RESOURCE_QUOTA_EXCEEDED`)**.
   - **7B. Reducción o mantenimiento de RAM ($\Delta\text{RAM} \le 0$):** El backend autoriza la operación directamente sin bloquearla aunque el nodo físico se encuentre saturado, dado que reducir memoria ayuda a descomprimir el servidor.
8. El backend envía la orden de reconfiguración a Proxmox VE (`PUT /nodes/{node}/{type}/{vmid}/config`).
9. El backend determina si el cambio requiere reinicio (`requiresReboot: true` cuando una VM `QEMU` está `running` sin soporte hot-plug completo o al reducir RAM en caliente; `requiresReboot: false` si la instancia está `stopped` o es un contenedor `LXC` que aplica límites de cgroups en vivo).
10. El backend registra `UPDATE_INSTANCE_CONFIG` (`EXITO`) en `auditoria` guardando el valor anterior y el nuevo, y responde `200 OK` con `{ cores, memoryMB, requiresReboot }`.
11. El frontend actualiza los indicadores de la instancia y, si `requiresReboot == true`, muestra un aviso destacado ofreciendo programar o ejecutar un reinicio (`CU-04`) para aplicar los nuevos recursos.

#### Flujos de Excepción y Salvaguardas
- **Cuota insuficiente para el incremento (`409 RESOURCE_QUOTA_EXCEEDED`):** Cuando $\Delta\text{RAM} > 0$ excede la RAM física disponible descontando la reserva del 15% del host.
- **Instancia protegida (`403 INSTANCE_PROTECTED`) o sin permiso (`403 INSTANCE_ACCESS_DENIED`).**

---

### CU-11: Recibir Notificaciones en Tiempo Real y Alertas Proactivas de Infraestructura

- **Requerimientos asociados:** `RF-11`, `RNF-02`, `RNF-04`.
- **Actor principal:** Usuario Autenticado (`ADMIN` u `OPERATOR`).
- **Actores secundarios:** Worker `EvaluadorUmbrales`, Worker `SeguimientoTareas`, Redis (estado de histéresis `alert:state:*`, tickets SSE y Pub/Sub), Servicio de Correo (`EmailService`), API de Proxmox VE.
- **Endpoints involucrados:**
  - `POST /api/events/ticket`
  - `GET /api/events?ticket=<uuid>` (Stream SSE)

#### Descripción
El sistema mantiene un canal unidireccional persistente Server-Sent Events (SSE) hacia el navegador de cada usuario autenticado y monitorea proactivamente en segundo plano la salud del nodo y de las instancias. Ante la finalización de una tarea asíncrona (`TASK_FINISHED`) o la superación sostenida de un umbral de recursos (`ALERT_TRIGGERED` / `ALERT_RESOLVED`), el sistema filtra los destinatarios según sus permisos RBAC, emite notificaciones emergentes en vivo con redirección directa (*deep-linking*) hacia la instancia afectada y despacha correos electrónicos de alta severidad a los administradores.

#### Precondiciones
1. Para el canal en vivo: el usuario posee una sesión activa válida.
2. Para el motor de alertas: el worker `EvaluadorUmbrales` se encuentra en ejecución continua en el backend.

#### Postcondiciones
- **De éxito:** Las alertas se registran en auditoría/historial, se entregan sin recargar la página únicamente a los usuarios autorizados para ver ese recurso, y se envían por correo electrónico a los administradores cuando la severidad es crítica.

#### Flujo Principal (Transacciones)

##### 11.1. Establecimiento Seguro del Canal SSE con Ticket de Un Solo Uso
1. Al autenticarse e ingresar al panel web, el frontend solicita `POST /api/events/ticket` enviando su `Authorization: Bearer <accessToken>`.
2. El backend verifica la sesión, genera un `ticket` UUID efímero en Redis (`sse:ticket:<uuid>`, TTL 30 segundos, asociado al `usuarioId`, `rol` y `sid` de la sesión) y responde `200 OK` con `{ ticket }`.
3. El frontend abre una conexión `EventSource` hacia `GET /api/events?ticket=<uuid>`.
4. El backend consume y borra atómicamente el ticket en Redis (`GETDEL`, garantizando que no pueda reutilizarse), establece los encabezados `Content-Type: text/event-stream`, emite el comentario inicial `: conectado` y mantiene la conexión viva enviando un latido `: ping` cada **25 segundos**.

##### 11.2. Detección Proactiva con Histéresis y Despacho Multicanal (`BAC-46`, `BAC-47`)
1. Cada **30 segundos**, el worker `EvaluadorUmbrales` consulta en segundo plano la telemetría del nodo físico y de las instancias en ejecución en Proxmox VE.
2. El worker compara el consumo actual contra los umbrales de saturación configurados (por ejemplo: RAM $\ge 85\%$, CPU $\ge 90\%$, Disco $\ge 90\%$).
3. **Salvaguarda Anti-Flapping (Histéresis en Redis):**
   - Si un recurso supera el umbral en un único muestreo aislado, el worker incrementa el contador en Redis (`alert:state:<recurso>:<metrica>`) pero **no** dispara la alerta todavía para evitar falsos positivos por picos momentáneos de CPU/red.
   - Si el recurso supera el umbral durante **2 ciclos consecutivos**, el estado pasa a `ACTIVE`, el backend registra la alerta en `auditoria`, publica el evento `ALERT_TRIGGERED` en Redis Pub/Sub (con `severity: "WARNING" | "CRITICAL"`, `vmid`, `instanceName`, `metric`, `value`) y, si es crítica, envía un correo electrónico de alerta a los usuarios `ADMIN` vía `EmailService`.
   - Mientras la alerta permanezca `ACTIVE`, el worker no vuelve a inundar de correos ni notificaciones repetidas.
   - Cuando la métrica desciende por debajo del **umbral de recuperación (ej. $< 80\%$)**, el worker limpia el estado en Redis y publica el evento `ALERT_RESOLVED`.
4. **Filtrado RBAC en el Hub SSE:**
   - Las alertas del nodo físico se transmiten a los clientes SSE con rol `ADMIN`.
   - Las alertas y eventos `TASK_FINISHED` de una instancia específica (`vmid`) se transmiten a todos los `ADMIN` y únicamente a aquellos `OPERATOR` que posean permiso `READ_ONLY` o `FULL_ACCESS` sobre ese `vmid`.
5. **Interacción y Redirección en el Frontend (`FRN-33`, `FRN-35`):**
   - El frontend recibe el evento SSE, despliega un toast flotante codificado por color de severidad e incrementa el contador de la campana del Centro de Notificaciones en la barra superior.
   - El usuario hace clic en la notificación o en el botón «Ver instancia».
   - El frontend redirige automáticamente a la vista de detalle de la instancia afectada (`/instances/:vmid`), donde el usuario visualiza las métricas en vivo (`CU-05`) y los botones de control para tomar acciones correctivas como reiniciar (`CU-04`) o ampliar recursos (`CU-10`).

#### Flujos de Excepción y Salvaguardas
- **Revocación Reactiva de Sesión en Vivo (`event: cierre`):** Si mientras el usuario navega un `ADMIN` suspende/elimina su cuenta, le resetea la contraseña/2FA, o el usuario cierra sesión en otra pestaña, el backend envía por el stream SSE el evento de control `event: cierre` con `{"motivo": "USER_INACTIVE" | "SESSIONS_REVOKED" | "LOGOUT"}` y cierra la conexión HTTP; el frontend intercepta este evento, limpia inmediatamente el estado en memoria y redirige a `/login`.

---

### CU-12: Recuperar Contraseña Olvidada mediante Código Temporal (OTP)

- **Requerimientos asociados:** `RF-13`, `RNF-01`.
- **Actor principal:** Usuario No Autenticado.
- **Actores secundarios:** Servicio de Correo (`EmailService` SMTP), PostgreSQL (`usuarios`, `codigos_recuperacion`, `sesiones_activas`, `auditoria`), Redis y Canal SSE.
- **Endpoints involucrados:**
  - `POST /api/auth/password/forgot`
  - `POST /api/auth/password/reset`

#### Descripción
El usuario que olvidó su contraseña solicita su restablecimiento ingresando su correo electrónico, recibe un código numérico temporal de un solo uso (OTP de 6 dígitos) con validez de 15 minutos y define una nueva contraseña segura, revocándose automáticamente cualquier sesión activa previa.

#### Precondiciones
1. El usuario conoce el correo electrónico asociado a su cuenta y tiene acceso a su bandeja de entrada.

#### Postcondiciones
- **De éxito:** Se actualiza `hash_contrasena` en PostgreSQL, se marca el código OTP como utilizado, se revocan todas las sesiones activas de la cuenta y se registra el evento en `auditoria`.

#### Flujo Principal (Transacciones)
1. En la pantalla de login (`/login`), el usuario selecciona «¿Olvidaste tu contraseña?» (`/forgot-password`), ingresa su dirección de correo electrónico y confirma.
2. El frontend envía `POST /api/auth/password/forgot` con `{ email }`.
3. **Salvaguarda contra Enumeración de Usuarios (`D6`):** El backend normaliza el correo y busca si pertenece a un usuario activo y no eliminado (`activo = true`, `eliminado_en IS NULL`). Exista o no la cuenta, el endpoint responde **siempre `200 OK`** con el mismo mensaje genérico («Si el correo electrónico se encuentra registrado, recibirás un código de verificación en los próximos minutos»).
4. Si la cuenta existe y está activa:
   1. El backend invalida en `codigos_recuperacion` cualquier código anterior pendiente para ese usuario.
   2. Genera criptográficamente un código numérico aleatorio de **6 dígitos**, almacena su hash en `codigos_recuperacion` con fecha de expiración a **15 minutos** (`expira_en = NOW() + 15m`) y contador `intentos_fallidos = 0`.
   3. Envía el código de 6 dígitos al correo del usuario mediante `EmailService` y registra la solicitud en `auditoria`.
5. El frontend avanza a la vista de verificación (`/reset-password`), donde el usuario ingresa el código de 6 dígitos recibido en su correo y su `nuevaContrasena`.
6. El frontend envía `POST /api/auth/password/reset` con `{ email, codigo, nuevaContrasena }`.
7. El backend verifica que exista un código vigente para ese usuario (`usado_en IS NULL`, `expira_en > NOW()`, `intentos_fallidos < 3`) y compara el código ingresado.
8. El backend valida que `nuevaContrasena` cumpla la política de complejidad (mínimo 8 caracteres, al menos 1 mayúscula, 1 número y 1 carácter especial) y que **no sea idéntica a la contraseña actual** del usuario.
9. En una transacción atómica, el backend actualiza `hash_contrasena` (bcrypt costo 12), establece `cambio_contrasena = false`, marca el código OTP como consumido (`usado_en = NOW()`), **revoca todas las sesiones activas del usuario** en `sesiones_activas` y Redis (emitiendo `SESSIONS_REVOKED` por SSE) y registra el restablecimiento en `auditoria` (`EXITO`).
10. El frontend muestra la confirmación de éxito y redirige al usuario a `/login` para que inicie sesión con su nueva contraseña y su segundo factor TOTP (`CU-01`).

#### Flujos de Excepción y Salvaguardas
- **Quema de Código por Fuerza Bruta (Máximo 3 Intentos — `D6`):** Cada vez que el usuario envía un código incorrecto en `POST /api/auth/password/reset`, el backend incrementa `intentos_fallidos` en `codigos_recuperacion` y responde `400 Bad Request` (`RESET_FAILED`). Al alcanzar el **tercer intento fallido**, el código queda invalidado definitivamente aunque no hayan transcurrido los 15 minutos, obligando a solicitar uno nuevo.
- **Código expirado o ya utilizado (`400 RESET_FAILED`):** Si transcurrieron más de 15 minutos desde su emisión o si ya fue consumido o reemplazado por uno más reciente.
- **Contraseña débil o igual a la anterior (`400 WEAK_PASSWORD` / `400 SAME_PASSWORD`):** El backend rechaza el cambio sin quemar el código OTP para que el usuario pueda corregir la contraseña elegida.
- **Rate Limiting por IP/Correo (`429 TOO_MANY_REQUESTS` — `FIX-118`):** Limita la cantidad de solicitudes a `/password/forgot` para evitar abuso del servidor SMTP o bombardeo de correos.

---

## 4. Matriz de Trazabilidad: Casos de Uso vs. Requerimientos y Salvaguardas Clave

| Caso de Uso | Requerimientos Cubiertos | Endpoints Principales | Salvaguardas Críticas Aplicadas |
| :--- | :--- | :--- | :--- |
| **CU-01** | `RF-01`, `RNF-01` | `/api/auth/login`, `/api/auth/2fa/*`, `/api/auth/refresh`, `/api/auth/logout`, `/api/account/*` | `jwtTemporal` (5m), bloqueo de QR si ya vinculado (`409 TWO_FACTOR_ALREADY_ENABLED`), cifrado TOTP AES-256-GCM, anti-replay `ultimo_totp_periodo`, bloqueo `403 PASSWORD_CHANGE_REQUIRED`, cookie `HttpOnly` + `SameSite=Strict`, corte reactivo SSE en logout. |
| **CU-02** | `RF-02`, `RNF-02`, `RNF-04` | `/api/node/status` | Caché Redis (5-10s), semáforo `Saludable` (<70%) / `Advertencia` (≥70%) / `Inaccesible` (502/504), cálculo de reserva del 15% de RAM del host, acotamiento de superficie para `OPERATOR` (`FIX-117`). |
| **CU-03** | `RF-03`, `RNF-01`, `RNF-02` | `/api/instances`, `/api/instances/:vmid` | Filtrado RBAC en dominio (`ADMIN` vs `OPERATOR`), inyección de `nivelAcceso` (`FULL_ACCESS` / `READ_ONLY`), resolución concurrente de IP post-RBAC (`BAC-23B`), persistencia de bloqueo UI vía `activeTask`. |
| **CU-04** | `RF-04`, `RF-09`, `RF-11`, `RNF-04` | `/api/instances/:vmid/status/:action`, `DELETE /api/instances/:vmid` | Orden estricto de guardas: `RequireInstanceAccess` -> `RejectProtectedInstance` (`100..105`) -> `RequireInstanceIdle` (`409 INSTANCE_BUSY`) -> `validarEstadoParaAccion` (`409 INSTANCE_INVALID_STATE`). Ocultamiento de `UPID` (`FIX-105`), timeout de 3 min (`D4`), doble auditoría (`PENDING` -> `EXITO`/`FALLA`), purga de permisos huérfanos al eliminar (`BAC-24C` / `FIX-113`). |
| **CU-05** | `RF-05`, `RNF-02`, `RNF-06` | `/api/auth/ws-ticket`, `/api/ws/instances/:vmid/metrics`, `/api/instances/:vmid/metrics` | Autenticación WebSocket mediante ticket efímero de un solo uso en Redis (`GETDEL`, TTL 30s), manejo *null-safe* cuando la instancia está `stopped`, caché RRDdata de 30s. |
| **CU-06** | `RF-06`, `RNF-03`, `RNF-04` | `/api/instances/:vmid/snapshots`, `/api/instances/:vmid/snapshots/:snapname/rollback` | Separación `READ_ONLY` (listar) vs `FULL_ACCESS` (crear/revertir/borrar), regex alfanumérico en `snapname` (`SEC-12`), confirmación destructiva escribiendo el nombre del snapshot (`FRN-30`), bloqueo `ROLLING_BACK`. |
| **CU-07** | `RF-07`, `RNF-01`, `RNF-04` | `/api/node/storages`, `/api/node/templates`, `POST /api/instances` | Exclusivo `ADMIN`, sanitización anti-inyección en strings de Proxmox (`FIX-116`), lock distribuido en Redis (`lock:instance:provisioning`), validación de cuota con reserva del 15% de RAM del host (`409 RESOURCE_QUOTA_EXCEEDED`), auto-asignación de permisos al concluir `UPID`. |
| **CU-08** | `RF-08`, `RNF-01`, `RNF-02` | `/api/admin/audit`, `/api/admin/audit/export` | Exclusivo `ADMIN`, inmutabilidad por trigger PostgreSQL (`trg_auditoria_inmutable`) + `REVOKE UPDATE/DELETE/TRUNCATE`, particionado trimestral, exportación CSV/JSON en streaming por cursor. |
| **CU-09** | `RF-01`, `RF-09`, `RNF-01` | `/api/admin/users/*` | Exclusivo `ADMIN`, rollback si falla SMTP (`502 EMAIL_DELIVERY_FAILED` — `D5`), suspensión reversible (`activo=false`) vs baja lógica irreversible (`eliminado_en` con índice único parcial para no romper FK de `auditoria`), prohibición de asignar VMIDs `100..105` (`FIX-119`), prohibición de auto-baja y escalada entre admins (`FIX-112`), revocación inmediata de sesiones/SSE. |
| **CU-10** | `RF-10`, `RNF-03`, `RNF-04` | `PUT /api/instances/:vmid/config` | Requiere `FULL_ACCESS` o `ADMIN`, validación de cuota únicamente sobre el delta positivo ($\Delta\text{RAM} > 0$), permiso directo para reducciones ($\Delta\text{RAM} \le 0$), indicador `requiresReboot`. |
| **CU-11** | `RF-11`, `RNF-02`, `RNF-04` | `/api/events/ticket`, `/api/events` | Ticket SSE efímero de un solo uso (`GETDEL`), filtrado RBAC de eventos por `vmid`, histéresis de 2 ciclos consecutivos en Redis antes de emitir `ALERT_TRIGGERED` y umbral de recuperación para `ALERT_RESOLVED`, despacho multicanal (UI + SMTP) y *deep-linking* a la instancia. |
| **CU-12** | `RF-13`, `RNF-01` | `/api/auth/password/forgot`, `/api/auth/password/reset` | Respuesta genérica `200 OK` contra enumeración de correos (`D6`), OTP de 6 dígitos con TTL de 15 min, quema automática del código al 3.er intento fallido, revocación global de sesiones activas tras el cambio. |
