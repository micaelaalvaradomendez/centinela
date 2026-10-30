# 4. Reportes de Auditoría de Seguridad: Vulnerabilidades, Modelo de Acceso y Mitigaciones

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
│ - Multi-tenancy: OrganizacionID no existía en Permisos ni en Auditoría │
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

## 4.2. Matriz de las 9 Brechas Críticas de Seguridad Detectadas (CWE)

La auditoría técnica formalizó **9 vulnerabilidades críticas** clasificadas bajo el estándar internacional de Common Weakness Enumeration (CWE), estructurando el plan de remediación inmediata asignado en los ítems `FIX-111` a `FIX-119`:

| Ref. Fix | Vulnerabilidad Detectada | Taxonomía CWE | Vector de Ataque y Escenario de Explotación | Remediació Técnican Implementada |
|---|---|---|---|---|
| **`FIX-111A` / `FIX-111B`** | **Secuestro de Cuenta por Cambio No Verificado de Correo** | **CWE-287** / **CWE-640** | El endpoint `PUT /api/account/email` aceptaba un nuevo email sin solicitar la contraseña actual ni el token 2FA. Un atacante con acceso momentáneo a una terminal desbloqueada cambiaba el correo por el suyo y ejecutaba el flujo de olvido de clave, adueñándose de cuentas administradoras. | Se implementó el modal con reautenticación obligatoria exigiendo contraseña actual y código TOTP de 6 dígitos antes de procesar el cambio. |
| **`FIX-112`** | **Fuerza Bruta y Condición de Carrera en Recuperación OTP** | **CWE-307** / **CWE-799** | Ausencia de rate limiting en `/auth/login` y `/auth/password/forgot`. El endpoint `/password/reset` leía `intentos_fallidos` sin bloqueo de fila (`SELECT ... FOR UPDATE`), permitiendo ataques de fuerza bruta paralela sobre el código numérico de 6 dígitos dentro de la ventana de 5 minutos. | Rate limiter distribuido en Redis con ventana deslizante (máximo 5 intentos por IP/usuario) y bloqueo pesimista en PostgreSQL. |
| **`BAC-24C` / `FIX-113`** | **Permisos Huérfanos por Reciclaje de VMID en Proxmox** | **CWE-284** | Al destruir una máquina (`DELETE /api/instances/:vmid`), el backend eliminaba la VM en el hipervisor pero **no borraba las filas en `permisos_instancia`**. Cuando Proxmox reciclaba el VMID mediante `/cluster/nextid`, los operadores anteriores heredaban acceso automático sobre la nueva máquina ajena. | Purga atómica transaccional de tuplas en `permisos_instancia` dentro del worker de eliminación, emitiendo el evento de auditoría `INSTANCE_PERMISSIONS_PURGED`. |
| **`FIX-114`** | **Persistencia de Sesiones tras Cambio Voluntario de Contraseña** | **CWE-613** | Al cambiar la clave (`PUT /api/account/password`), se actualizaba el hash bcrypt en la base de datos pero **no se invalidaban las sesiones activas en Redis ni en `sesiones_activas`**. Un intruso con una cookie de refresh o JWT seguía operando indefinidamente. | Invalidación transaccional de todas las sesiones concurrentes en Redis y base de datos, forzando la desconexión inmediata de sockets SSE y WebSockets. |
| **`FIX-115A` / `FIX-115B`** | **Restricción Anticuada de Contraseñas a 12 Caracteres** | **CWE-521** | Validadores en backend (`password.go`) y frontend imponían una longitud máxima de 12 caracteres (`max=12`), violando la guía NIST SP 800-63B e impidiendo el uso de contraseñas de alta entropía (*passphrases*) o gestores de claves. | Ampliación del límite superior a 128 caracteres tanto en esquemas Zod de React como en validadores de Go. |
| **`FIX-116`** | **Inyección de Parámetros y Directivas hacia Proxmox VE** | **CWE-88** | La construcción de la cadena de red interpolaba directamente el puente: `fmt.Sprintf("name=eth0,bridge=%s,...", input.Bridge)`. Un usuario malicioso podía enviar `vmbr0,rate=50,tag=100`, inyectando asignaciones de VLANs y cotas de tráfico indebidas en el hipervisor. | Sanitización estricta mediante expresión regular `^[a-zA-Z0-9]+$` rechazando comas, saltos de línea y caracteres especiales. |
| **`FIX-117A` / `FIX-117B`** | **Fuga de Información Global del Host a Operadores** | **CWE-200** | `/api/node/status` devolvía a cualquier usuario autenticado la totalidad del almacenamiento en disco del nodo físico, permitiendo a operadores no privilegiados deducir la actividad global del clúster. | Segregación de DTOs en backend (omitiendo el bloque `storage` para el rol `OPERATOR`) y adaptación condicional del Dashboard en frontend. |
| **`FIX-118`** | **Credenciales por Defecto Hardcodeadas en Migraciones** | **CWE-798** | Los scripts de inicialización SQL y semillas de GORM insertaban un usuario administrador con credenciales fijas (`admin@elcentinela.com` / `Centinela2026!`), creando una puerta trasera conocida en entornos productivos. | Eliminación del sembrado automático de contraseñas fijas; implementación de un flujo de configuración inicial en primer arranque. |
| **`FIX-119`** | **Asignación de Permisos sobre VMIDs Inexistentes o Protegidos** | **CWE-285** | La API de permisos permitía asignar a operadores acceso sobre VMIDs inexistentes o sobre las máquinas internas de infraestructura (`100..105`), donde reside el propio sistema El Centinela. | Validación previa contra el inventario real y bloqueo estricto de asignación sobre la lista de VMIDs protegidos configurados en el sistema. |

---

## 4.3. Prevención de IDOR: Sustitución de VMIDs por Tokens Opacos (`FIX-105` a `FIX-110`, `FIX-124`)

Una vulnerabilidad de diseño de alto impacto detectada en la auditoría fue la exposición directa de los identificadores internos de Proxmox en las rutas de la API pública y en la URL del navegador:
*   **Problema (Insecure Direct Object References - IDOR):** Rutas como `GET /api/instances/101` o `/instances/102/metrics` exponían enteros secuenciales predecibles. Esto permitía a cualquier usuario enumerar mediante scripting todas las máquinas virtuales del centro de cómputos.
*   **Solución Arquitectónica:** Se implementó una capa de ofuscación bidireccional basada en **Tokens Opacos Criptográficos (`InstanceRef`)**:
    *   Formato público: `inst_[a-zA-Z0-9_-]{12,64}` (tokens codificados en base64url con firma o mapeo interno).
    *   El backend traduce transparentemente el token opaco al `VMID` interno únicamente dentro de los adaptadores de infraestructura hacia Proxmox.
    *   La suite de pruebas `test/front/fix-108-ids-ofuscados.test.tsx` y `fixes_acceptance_test.go` verifican que ningún endpoint público acepte enteros planos ni exponga trazas de depuración que revelen la numeración del hipervisor.

---

## 4.4. Aporte al Artículo: La Auditoría Algorítmica Asistida por el Gestor

La literatura actual de Project Management destaca que los modelos de lenguaje pueden generar código rápidamente, pero tienden a reproducir malas prácticas de seguridad o a introducir vulnerabilidades lógicas sutiles que pasan inadvertidas para los desarrolladores menos experimentados.

El caso de *El Centinela* ilustra cómo la auditoría de seguridad liderada por el AI Manager:
1.  **Transforma hallazgos forenses en especificaciones ejecutables:** Cada brecha de seguridad fue formalizada con su respectivo código CWE, su escenario de ataque y su prueba automatizada adversarial (`adversarial_integrity_test.go`).
2.  **Garantiza la inmutabilidad de la auditoría:** Se protegió la tabla `auditoria` de PostgreSQL mediante triggers que prohíben operaciones `UPDATE` y `DELETE`, asegurando que ni siquiera un administrador comprometido pueda alterar las trazas de ejecución.
3.  **Evita el sesgo de complacencia:** Demuestra que la métrica de éxito de un proyecto no es la cantidad de código generado, sino la robustez y la seguridad demostrable de la arquitectura desplegada.
