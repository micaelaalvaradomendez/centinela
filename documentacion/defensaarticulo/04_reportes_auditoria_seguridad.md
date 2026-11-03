# 4. Reportes de Auditoría de Seguridad: Vulnerabilidades de Código, Modelo de Acceso e Infraestructura

**Contexto en el Artículo:** Sección V.C (*"Cuadro XII: Tipos de auditoría automatizada sobre el código"*) y Sección IV (*"Superación del judgment gap: liderazgo, estrategia y personas"*).

---

## 4.1. Diagnóstico Estructural: La Brecha entre el Diseño Visual y el Modelo de Datos

En el desarrollo de software acelerado por IA, existe un fenómeno recurrente: la generación de **interfaces visuales que aparentan controles sofisticados pero que carecen de respaldo en la arquitectura del backend**. En *El Centinela*, la auditoría de seguridad estática y dinámica realizada el **19/10/2026** (documentada en [`analisis-seguridad-vulnerabilidades-permisos.md`](file:///home/micaela/centinela/documentacion/analisis/analisis-seguridad-vulnerabilidades-permisos.md)) desnudó una desconexión crítica entre frontend y backend:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        FRONTEND (UI / VISTAS)                          │
│ Promete matriz granular: "Gestionar snapshots", "Ver red",            │
│ "Gestionar plantillas", "Gestionar configuración global"               │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │  DESCONEXIÓN ESTRUCTURAL
┌───────────────────────────────────▼────────────────────────────────────┐
│                        BACKEND (MODELO REAL)                           │
│ - Rol estrictamente binario: ADMIN (bypass total) u OPERATOR           │
│ - Instancias: PermisoInstancia (UsuarioID, Vmid, FULL_ACCESS|READ_ONLY)│
│ - Huérfanos: Si se borra la VM, los permisos NO se eliminaban          │
│ - Multi-tenancy: OrganizacionID residual ocultaba eventos de sistema   │
└────────────────────────────────────────────────────────────────────────┘
```

### La Falsa Granularidad en la Interfaz:
En las pantallas de administración de usuarios ([`CrearUsuarios.tsx`](file:///home/micaela/centinela/frontend/centinela/src/pages/CrearUsuarios.tsx#L161-L172)), la interfaz presentaba checkboxes para habilitar permisos específicos (red, snapshots, configuración global). Sin embargo, en el backend Go **no existía ninguna tabla de capacidades (`capabilities`), scopes de JWT ni middleware granular**. El backend solo evaluaba:
```go
// backend/internal/adapters/primary/http/middleware/instance_guard.go
if rol == "ADMIN" {
    c.Next()
    return // Bypass absoluto de todas las comprobaciones
}
```
Si el usuario era `OPERATOR`, únicamente se validaba si existía un registro en la tabla `permisos_instancia` con el par `(usuario_id, vmid)`. Todo el panel visual de permisos granulares era puramente decorativo. La detección de este "espejismo" requirió una auditoría de juicio humano que contrastara los esquemas SQL relacionales contra el árbol de componentes de React.

---

## 4.2. Matriz de las 9 Brechas Críticas de Seguridad en Código Detectadas (CWE)

La auditoría técnica formalizó **9 vulnerabilidades críticas** clasificadas bajo el estándar internacional de Common Weakness Enumeration (CWE), estructurando el plan de remediación asignado en los ítems `FIX-111` a `FIX-119` (y ampliado en `FIX-138`):

| Ref. Fix | Vulnerabilidad Detectada | Taxonomía CWE | Vector de Ataque y Escenario de Explotación | Remediación Técnica Implementada |
|---|---|---|---|---|
| **`FIX-111A` / `FIX-111B`** | **Secuestro de Cuenta por Cambio No Verificado de Correo** | **CWE-287** / **CWE-640** | El endpoint `PUT /api/account/email` aceptaba un nuevo email sin solicitar la contraseña actual ni el token 2FA. Un atacante con acceso momentáneo a una terminal desbloqueada cambiaba el correo por el suyo y ejecutaba el flujo de olvido de clave, adueñándose de cuentas administradoras. | Se implementó el modal con reautenticación obligatoria exigiendo contraseña actual y código TOTP de 6 dígitos antes de procesar el cambio (commits `dc50c0a`, `57512c9`). |
| **`FIX-112` / `BAC-48`** | **Fuerza Bruta y Condición de Carrera en Recuperación OTP** | **CWE-307** / **CWE-799** | Ausencia de rate limiting en `/auth/login` y `/auth/password/forgot`. El endpoint `/password/reset` leía `intentos_fallidos` sin bloqueo de fila (`SELECT ... FOR UPDATE`), permitiendo ataques de fuerza bruta paralela sobre el código numérico de 6 dígitos. | Rate limiter distribuido en Redis (máximo 5 intentos/min con respuesta HTTP 429) y serialización pesimista en PostgreSQL (commits `f738f1c`, `ef53022`, `651e8de`). |
| **`BAC-24C` / `FIX-113`** | **Permisos Huérfanos por Reciclaje de VMID en Proxmox** | **CWE-284** | Al destruir una máquina (`DELETE /api/instances/:vmid`), el backend eliminaba la VM en el hipervisor pero **no borraba las filas en `permisos_instancia`**. Cuando Proxmox reciclaba el VMID mediante `/cluster/nextid`, los operadores anteriores heredaban acceso automático sobre la nueva máquina ajena. | Purga atómica transaccional de tuplas en `permisos_instancia` dentro del worker de eliminación, emitiendo el evento de auditoría `INSTANCE_PERMISSIONS_PURGED` (commit `c08af54`). |
| **`FIX-114`** | **Persistencia de Sesiones tras Cambio Voluntario de Contraseña** | **CWE-613** | Al cambiar la clave (`PUT /api/account/password`), se actualizaba el hash bcrypt en la base de datos pero **no se invalidaban las sesiones activas en Redis ni en `sesiones_activas`**. Un intruso con una cookie de refresh o JWT seguía operando indefinidamente. | Invalidación transaccional de todas las sesiones concurrentes en Redis y base de datos, forzando la desconexión inmediata de canales SSE y WebSockets (commits `638bb5e`, `a0f5465`). |
| **`FIX-115A/B` + `FIX-138`** | **Restricción Anticuada de Contraseñas (12 chars y símbolos limitados)** | **CWE-521** | Validadores en backend (`password.go`) y frontend imponían una longitud máxima de 12 caracteres (`max=12`) y aceptaban solo 12 símbolos (`!@#$%^&*-_=+`), violando la guía **NIST SP 800-63B** y rechazando claves de Bitwarden/1Password con `.`, `,`, `?`, `/`, `()`. | Ampliación del límite superior a 128 caracteres (`b7cdaf0`, `5926c44`) y extensión del juego de caracteres especiales ASCII y passphrases con espacios en `FIX-138` (commit `3fe055b`). |
| **`FIX-116`** | **Inyección de Parámetros y Directivas hacia Proxmox VE** | **CWE-88** | La construcción de la cadena de red interpolaba directamente el puente: `fmt.Sprintf("name=eth0,bridge=%s,...", input.Bridge)`. Un usuario malicioso podía enviar `vmbr0,rate=50,tag=100`, inyectando asignaciones de VLANs y cotas de tráfico indebidas en el hipervisor. | Sanitización estricta mediante expresión regular `^[a-zA-Z0-9]+$` rechazando comas, saltos de línea y caracteres especiales (commit `a25d1b3`). |
| **`FIX-117A` / `FIX-117B`** | **Fuga de Información Global del Host a Operadores** | **CWE-200** | `/api/node/status` devolvía a cualquier usuario autenticado la totalidad del almacenamiento en disco del nodo físico, permitiendo a operadores no privilegiados deducir la actividad global del clúster. | Segregación de DTOs en backend (omitiendo el bloque `storage` para el rol `OPERATOR`, commit `94e95a2`) y adaptación condicional del Dashboard en frontend (`5c993ee`). |
| **`FIX-118`** | **Credenciales por Defecto Hardcodeadas en Migraciones** | **CWE-798** | Los scripts de inicialización SQL y semillas de GORM insertaban un usuario administrador con credenciales fijas (`admin@elcentinela.com` / `Centinela2026!`), creando una puerta trasera conocida en entornos productivos. | Eliminación del sembrado automático de contraseñas fijas; configuración inicial segura en primer arranque (commit `1801e33`). |
| **`FIX-119` + `ADP-04`** | **Asignación de Permisos sobre VMIDs Inexistentes o Protegidos** | **CWE-285** | La API de permisos permitía asignar a operadores acceso sobre VMIDs inexistentes o sobre las máquinas internas de infraestructura (`100..105`), donde reside el propio sistema El Centinela. | Validación previa contra el inventario real y bloqueo estricto configurable dinámicamente vía `PROXMOX_PROTECTED_VMIDS` en `.env` (commits `5c72bbb`, `c09a867`). |

---

## 4.3. Prevención de IDOR y Air-Gap de Metadatos (`FIX-105` a `FIX-110`, `FIX-124`)

Una vulnerabilidad de diseño de alto impacto detectada en la auditoría ([`analisis-cifrado-ids-vulnerabilidades-infraestructura.md`](file:///home/micaela/centinela/documentacion/analisis/analisis-cifrado-ids-vulnerabilidades-infraestructura.md)) fue la exposición directa de los identificadores internos de Proxmox en las respuestas JSON, rutas REST y URLs del navegador:
1.  **Ocultamiento de UPID (`FIX-105` / `FIX-106`):** El formato crudo `UPID:<node>:<pid>:<pstart>:<starttime>:<type>:<vmid>:<user>:` revelaba el nombre del nodo, el PID del proceso en el kernel de Linux y el usuario de sistema (`root@pam`). Se reemplazó hacia el cliente por un identificador opaco `tareaId` (UUIDv4), manteniendo el mapeo real únicamente en la tabla interna `tareas_asincronas`.
2.  **Sustitución de VMIDs por Tokens Opacos (`FIX-107` / `FIX-108` / `FIX-124`):**
    *   Rutas como `GET /api/instances/101` exponían enteros secuenciales predecibles susceptibles a enumeración (CWE-639 / IDOR).
    *   Se implementó una capa de ofuscación bidireccional basada en **Tokens Opacos Criptográficos (`InstanceRef`)** con formato `inst_[a-zA-Z0-9_-]{12,64}`, validados estrictamente contra *Path Traversal* (`../`) en frontend y backend (`fix-108-ids-ofuscados.test.tsx` 11/11 PASS e `ids-ofuscados.test.tsx` 14/14 PASS).
3.  **Saneamiento de Trazas y Exportación Segura (`FIX-109` / `FIX-110` / `BAC-45`):** Enmascaramiento de nombres de nodos y storages para usuarios no administradores, limpieza de trazas internas de Proxmox en errores HTTP/SSE, y sanitización contra **Inyección de Fórmulas CSV (CWE-1236)** prefijando caracteres peligrosos (`=`, `+`, `-`, `@`) en la exportación de auditoría en streaming.

---

## 4.4. Auditoría de Seguridad de Red e Infraestructura en el Servidor Final (`pve-is3-g1`)

El **29 y 30 de octubre de 2026**, el alcance de la auditoría liderada por la PM se extendió desde el código fuente hacia la infraestructura real del **Servidor Final del Grupo 1 (`pve-is3-g1`, Proxmox VE 8.4.21)**, documentándose en [`comparativa-servidor-prueba-vs-final.md`](file:///home/micaela/centinela/documentacion/servidorfinal/comparativa-servidor-prueba-vs-final.md), [`analisis-estado-servidor.md`](file:///home/micaela/centinela/documentacion/servidorfinal/analisis-estado-servidor.md) y [`SMTP-y-Brevo.md`](file:///home/micaela/centinela/documentacion/SMTP-y-Brevo.md) cuatro hallazgos críticos de seguridad e infraestructura:

1.  **Violación de Aislamiento de Capas en Red Virtual (`centinela-redis` en `vmbr1`):**
    *   El nodo físico se diseñó con tres puentes virtuales aislados: `vmbr1` (aplicación/lógica), `vmbr2` (DMZ/clientes) y `vmbr3` (capa de datos, donde reside PostgreSQL `10.10.3.60`).
    *   Sin embargo, el contenedor `104 (centinela-redis)` fue desplegado en `vmbr1` (`10.10.1.15`), compartiendo segmento de red con el contenedor del frontend (`103: 10.10.1.20`) y el proxy (`100: 10.10.1.10`). Dado que Redis almacena tokens de sesión, tickets efímeros de WebSocket/SSE y bloqueos distribuidos, la auditoría ordenó su migración estricta a la red de datos `vmbr3` (`10.10.3.15/24`).
2.  **Egreso Acotado de Firewall para el Relay SMTP (`centinela-smtp-brevo`):**
    *   Los correos de alta de usuarios y recuperación OTP quedaban suspendidos en cola (`Connection timed out` hacia `smtp-relay.brevo.com:587`) porque el firewall del host bloqueaba la salida de las redes internas.
    *   En lugar de abrir la salida a internet a todos los contenedores (lo que violaría el aislamiento del clúster), se diseñó e implementó una regla quirúrgica de `iptables` (`POSTROUTING MASQUERADE` y `FORWARD ACCEPT` antes del `DROP` general) restringida exclusivamente a la IP `/32` del contenedor relay y a los puertos TCP `587,465`, verificando negociación **STARTTLS con TLS 1.3** y entrega en 2.3 segundos (PR #5, commit `4d8b3f8`).
3.  **Aislamiento de Acceso con Jump Host Bastión (`centi-mica`, LXC 106):**
    *   Ante el cierre de los puertos 80/443 en la interfaz Tailscale del host Proxmox (`100.116.49.111`), se aprovisionó el contenedor bastión `centi-mica` (`100.114.44.18`) con `tailscale serve` apuntando exclusivamente al proxy reverso interno (`10.10.1.10:80`), evitando exponer los contenedores internos de base de datos o API directamente a la VPN.
4.  **Endurecimiento Perimetral de Cabeceras HTTP (`INF-11` / `FIX-129`):**
    *   La suite `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS` verificó la redirección 80 $\rightarrow$ 443 y protocolos `TLSv1.2 TLSv1.3`, e identificó la necesidad de incorporar `Strict-Transport-Security` con `includeSubDomains`, `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff` y `Content-Security-Policy` con soporte explícito para `ws:` y `wss:` en `docker/nginx-edge.conf` (`FIX-129`).

---

## 4.5. Aporte al Artículo: La Auditoría Algorítmica y Sociotécnica Asistida por el Gestor

La literatura actual de Project Management destaca que los modelos de lenguaje pueden generar código rápidamente, pero tienden a reproducir malas prácticas de seguridad o a introducir vulnerabilidades lógicas sutiles entre capas que pasan inadvertidas en revisiones superficiales.

El caso de *El Centinela* ilustra cómo la auditoría liderada por el AI Manager:
1.  **Transforma hallazgos forenses en especificaciones ejecutables:** Cada brecha de seguridad fue formalizada con su respectivo código CWE, su escenario de ataque y su prueba automatizada en `etapa2_seguridad_acceptance_test.go` y `etapa3_fixes_pendientes_test.go`.
2.  **Garantiza la inmutabilidad de la auditoría:** Se protegió la tabla `auditoria` de PostgreSQL mediante el trigger `trg_auditoria_inmutable` y revocación de permisos `UPDATE`/`DELETE`, asegurando además en `FIX-137` que ningún filtro residual de `organizacion_id` oculte los eventos críticos generados por el sistema.
3.  **Une la seguridad del código con la seguridad de la infraestructura:** Demuestra que el rol estratégico del gestor abarca tanto los contratos de software como la topología real de red donde el sistema se ejecuta.
