# Informe de Auditoría de Seguridad: Vulnerabilidades, Modelo de Acceso y Brechas Críticas

**Fecha:** 19 de Octubre de 2026  
**Área:** Arquitectura de Software, Ciberseguridad y Control de Acceso  
**Repositorios Auditados:**  
*   `backend` (`centinela-back`, commit `8daaf87`)  
*   `frontend` (`centinela-front`, commit `91eabf4`)  
**Documentos de Referencia:** `terminado.md`, `terminado-1.md`, `terminado2.md`, `actual.md`, `futuro2.md`.

---

## 1. Introducción y Alcance

Este documento consolida el análisis técnico de seguridad estática, revisión de código fuente y evaluación de arquitectura realizado sobre el ecosistema **El Centinela**.

El objetivo de la auditoría consistió en:
1.  Examinar la correspondencia entre los requerimientos funcionales declarados en las etapas base, 1 y 2 frente a la implementación real en código.
2.  Evaluar la solidez del modelo de control de acceso basado en roles (RBAC) y control de acceso basado en recursos (instancias virtuales QEMU y contenedores LXC en Proxmox VE).
3.  Identificar brechas de seguridad críticas (vulnerabilidades lógicas, fallas de autenticación, escalamiento de privilegios, fugas de información e inyecciones de configuración).
4.  Fundamentar detalladamente el plan de remediación plasmado en los ítems `FIX-111` a `FIX-119` del archivo [`documentacion/futuro2.md`](../futuro2.md).

---

## 2. Diagnóstico Estructural: Por qué los Roles, Permisos y Acceso a Instancias están "Modelados pero a la vez No"

Durante la inspección de ambos submódulos se constató una fractura entre la conceptualización de la seguridad en el diseño original y su materialización técnica en el backend y frontend:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        FRONTEND (UI / VISTAS)                          │
│ Promete matriz granular: "Gestionar snapshots", "Ver red",            │
│ "Gestionar plantillas", "Gestionar configuración global"               │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │  DESCONEXIÓN CONCEPTUAL
┌───────────────────────────────────▼────────────────────────────────────┐
│                        BACKEND (MODELO REAL)                           │
│ - Rol estrictamente binario: ADMIN (bypass total) u OPERATOR           │
│ - Instancias: PermisoInstancia (UsuarioID, Vmid, FULL_ACCESS|READ_ONLY)│
│ - Huérfanos: Si se borra la VM, los permisos NO se eliminan           │
│ - Multi-tenancy: OrganizacionID no existe en Permisos ni en Auditoría  │
└────────────────────────────────────────────────────────────────────────┘
```

### A. Falsa granularidad en la Interfaz de Usuario
En [`frontend/centinela/src/pages/CrearUsuarios.tsx`](../../frontend/centinela/src/pages/CrearUsuarios.tsx#L161-L172) y en [`rolesAndPermissions.tsx`](../../frontend/centinela/src/components/features/users/components/rolesAndPermissions.tsx#L122-L141), la interfaz gráfica exhibe un panel informativo con una supuesta matriz de permisos granulares:
*   *“Gestionar snapshots y backups propios”*
*   *“Ver plantillas y almacenamiento”*
*   *“Gestionar configuración global”*
*   *“Ver red”*

**Evidencia de la brecha:**  
En el backend, **no existe ningún claim, scope, tabla ni middleware de permisos por acción**. No hay tabla de capacidades (`capabilities` o `roles_permissions`). El modelo es **estrictamente binario**: o el usuario tiene `claims.Rol == "ADMIN"` o tiene `claims.Rol == "OPERATOR"`. La interfaz muestra una sofisticación de control de acceso que es puramente decorativa.

### B. El Bypass absoluto del Rol ADMIN y la Inconsistencia de BAC-39
En [`backend/internal/adapters/primary/http/middleware/instance_guard.go`](../../backend/internal/adapters/primary/http/middleware/instance_guard.go#L46-L49):
```go
if rol == "ADMIN" {
    c.Next()
    return
}
```
El middleware de autorización por recurso (`RequireInstanceAccess`) omite de forma total e inmediata cualquier verificación en la base de datos si el usuario posee rol `ADMIN`. Los administradores tienen acceso implícito e irrestricto a todos los VMIDs.

**La contradicción en BAC-39:**  
En [`backend/internal/core/services/seguimiento_tareas.go`](../../backend/internal/core/services/seguimiento_tareas.go), al finalizar con éxito la creación de una instancia (`BAC-39`), el sistema inserta una fila en `permisos_instancia` vinculando al creador con `FULL_ACCESS`. Dado que solo un `ADMIN` puede invocar `POST /api/instances`, esa fila almacenada en la base de datos es **completamente inútil y redundante**, pues el middleware de acceso jamás la consulta. Esto evidencia un desacople entre el equipo que diseñó las historias de usuario y quienes implementaron los guards de seguridad.

### C. Desconexión en el Flujo de Creación de Usuarios
Al dar de alta un usuario en el formulario [`CrearUsuarios.tsx`](../../frontend/centinela/src/pages/CrearUsuarios.tsx) (`POST /api/admin/users`), el payload únicamente admite `{ nombreCompleto, nombreUsuario, emailUsuario, rol }`. No existe mecanismo en el alta para asignarle máquinas al nuevo usuario.
Un operador recién creado nace obligatoriamente con **cero instancias asignadas**. Para que pueda operar, el administrador debe salir, ingresar a la vista de detalle de usuario ([`DetailsUserPage.tsx`](../../frontend/centinela/src/pages/detailsUserPage.tsx)) y configurar manualmente las máquinas en una segunda petición desacoplada.

### D. Fuga de Permisos Huérfanos por Reciclaje de VMID (`BAC-24C`)
En Proxmox VE, el clúster asigna identificadores numéricos a través de `/cluster/nextid` retornando el menor número entero libre.
Cuando se ejecuta `DELETE /api/instances/:vmid`, la máquina se borra físicamente en el hipervisor, pero **la tabla `permisos_instancia` no es limpiada**. Si con posterioridad se crea una máquina que vuelve a tomar ese número (por ejemplo, el VMID `110`), cualquier operador que tuvo asignada la máquina vieja **obtiene acceso inmediato y descontrolado sobre la nueva**, lo que constituye una falla crítica de autorización (IDOR).

### E. El Espejismo de la Organización (Multi-Tenancy Fantasma)
1.  Existe una tabla `organizaciones` y una columna `Usuario.OrganizacionID`.
2.  Sin embargo, la tabla [`PermisoInstancia`](../../backend/internal/core/domain/models.go#L91-L96) **no contiene `OrganizacionID`**.
3.  Proxmox VE es un único clúster compartido con un espacio plano y global de VMIDs.
4.  La tabla `auditoria` no tiene `organizacion_id`: para filtrar por organización realiza un `LEFT JOIN usuarios u ON u.id = a.usuario_id WHERE u.organizacion_id = ?`. Como consecuencia técnica, **todas las operaciones del sistema, tareas en segundo plano sin actor (`usuario_id = NULL`) o acciones de usuarios eliminados son invisibilizadas en el registro de auditoría**.
5.  El endpoint `POST /organizations` que el frontend invoca en [`organizationService.ts`](../../frontend/centinela/src/components/features/auth/services/organizationService.ts#L9) **no existe en el backend** (devuelve 404), dejando el flujo de registro como código muerto.

---

## 3. Análisis Detallado de las 9 Brechas de Seguridad Detectadas

---

### Brecha 1: Secuestro de Cuenta (Account Takeover) vía Cambio de Email sin Verificación
*   **Fix Asociado:** `FIX-111`
*   **Clasificación:** 🔴 **CRÍTICA** | CWE-284 (Improper Access Control) | CWE-306 (Missing Authentication for Sensitive Function)
*   **Archivos Afectados:**
    *   [`backend/internal/adapters/primary/http/account_handler.go`](../../backend/internal/adapters/primary/http/account_handler.go#L66-L84) (Líneas 66-84)
    *   [`backend/internal/core/services/user_service.go`](../../backend/internal/core/services/user_service.go#L404-L440) (Líneas 404-440)

#### Cómo se Detectó:
Al auditar los métodos expuestos en `AccountHandler`, se analizó la función `ActualizarPerfil`:
```go
// account_handler.go:78
userID := extraerUserID(c)
usuario, err := h.service.ActualizarPerfil(c.Request.Context(), userID, input)
```
Rastreando la implementación en `user_service.go:416-430`:
```go
if input.EmailUsuario != "" && input.EmailUsuario != usuario.EmailUsuario {
    // Solo valida que el email no esté duplicado en la org
    if existe, err := s.userRepo.ExisteEmailEnOrg(...); ... { ... }
    cambios["email_usuario"] = input.EmailUsuario
}
if len(cambios) > 0 {
    s.userRepo.ActualizarUsuario(ctx, usuarioID, cambios)
}
```
Se constató que:
1.  No se exige la contraseña actual (`input.ContrasenaActual`).
2.  No se exige un código TOTP de 2FA.
3.  No se envía confirmación previa al nuevo email antes de persistirlo.
4.  **No existe llamada a `s.auditSvc.Registrar`**, haciendo la modificación completamente silenciosa.

#### Vector de Ataque y Escenario de Explotación:
1.  El atacante aprovecha una sesión desatendida en una oficina o extrae el Access Token mediante una vulnerabilidad XSS o malware en el cliente.
2.  Emite un request:
    ```http
    PUT /api/account/profile HTTP/1.1
    Authorization: Bearer <token_robado>
    Content-Type: application/json

    {
      "emailUsuario": "atacante@evil.com"
    }
    ```
3.  El backend actualiza la dirección inmediatamente sin pedir ninguna credencial adicional.
4.  El atacante acude a `POST /api/auth/password/forgot` enviando `"email": "atacante@evil.com"`.
5.  Recibe el código OTP de 6 dígitos en su propia casilla, invoca `POST /api/auth/password/reset` y cambia la contraseña de la cuenta.
6.  El usuario original queda despojado de su cuenta sin haber recibido ninguna advertencia ni existir registro en la auditoría del sistema.

#### Importancia:
Permite la toma de control total, definitiva e inadvertida de cualquier cuenta activa en la plataforma, incluyendo administradores.

---

### Brecha 2: Ausencia Total de Rate Limiting en Autenticación y Fuerza Bruta
*   **Fix Asociado:** `FIX-112`
*   **Clasificación:** 🔴 **CRÍTICA** | CWE-307 (Improper Restriction of Excessive Authentication Attempts) | OWASP A07:2021
*   **Archivos Afectados:**
    *   [`backend/cmd/api/main.go`](../../backend/cmd/api/main.go#L214-L245) (Líneas 214-245)
    *   [`backend/internal/core/services/auth_service.go`](../../backend/internal/core/services/auth_service.go#L239-L253) (Líneas 239-253)

#### Cómo se Detectó:
Se revisó la cadena de middlewares en `main.go`. Se encontraron `Recovery()`, `RequestLogger()`, `CORS()` y `SecurityHeaders()`, pero **ningún middleware de limitación de tasa (Rate Limiting)**.
Al auditar `VerificarTotp` en `auth_service.go`:
```go
valido, err := crypto.ValidarCodigo(usuario.SecretoTotpCifrado, encKey, codigo)
if !valido {
    s.auditSvc.Registrar(...)
    return nil, fmt.Errorf("código TOTP incorrecto")
}
```
Si el código falla, simplemente devuelve error. No hay contador de reintentos sobre el token temporal `jtiTemporal`.
En `SolicitarRecuperacionContrasena` (`auth_service.go:594`), el endpoint devuelve 200 siempre, pero no acota la cantidad de llamadas que un cliente puede hacer hacia el servicio SMTP.

#### Vector de Ataque y Escenario de Explotación:
1.  **Fuerza bruta sobre contraseñas:** Un atacante con una lista de diccionarios envía cientos de peticiones por segundo a `/api/auth/login` hasta vulnerar claves débiles.
2.  **Fuerza bruta sobre 2FA:** Tras superar el primer paso de login, el token temporal `pre-auth` tiene una vigencia de 5 minutos (300 segundos). La ventana de validación de TOTP admite $\pm 1$ período (3 códigos válidos de 1.000.000 de combinaciones). Mediante un script que ejecute hilos paralelos, el atacante puede probar decenas de miles de códigos en los 5 minutos sin ser bloqueado.
3.  **Saturación de Correo:** Envío masivo de peticiones a `/api/auth/password/forgot` provocando el agotamiento de cuotas del relay SMTP y denegación de servicio.

#### Importancia:
Anula la eficacia del doble factor de autenticación (2FA) y expone la plataforma a ataques automatizados masivos de fuerza bruta y saturación de infraestructura.

---

### Brecha 3: Persistencia de Permisos Huérfanos por Reciclaje de VMID (`BAC-24C`)
*   **Fix Asociado:** `FIX-113`
*   **Clasificación:** 🔴 **ALTA** | CWE-284 (Improper Access Control) | CWE-639 (Authorization Bypass Through User-Controlled Key)
*   **Archivos Afectados:**
    *   [`backend/internal/adapters/primary/http/instance_handler.go`](../../backend/internal/adapters/primary/http/instance_handler.go#L870-L935) (`EliminarInstancia`)
    *   [`backend/internal/core/services/seguimiento_tareas.go`](../../backend/internal/core/services/seguimiento_tareas.go)

#### Cómo se Detectó:
Se analizó el ciclo de vida destructivo en `instance_handler.go:EliminarInstancia`. El handler invoca:
```go
tareaID, err := h.seguimiento.Seguir(ctx, upid, vmid, usuarioID, ports.AccionProxmoxEliminarVM, ...)
```
En `seguimiento_tareas.go`, cuando la tarea asíncrona de Proxmox termina con éxito (`exitstatus: OK`), el código notifica a la auditoría y publica `TASK_FINISHED`. **En ningún punto se ejecuta un DELETE sobre la tabla `permisos_instancia`**.
Esto fue ratificado por el informe de pruebas en [`test/informe.md`](../../test/informe.md#L48-L50), donde la prueba de aceptación `TestINT01` paso 6 (`BAC-24C`) figura como fallida/pendiente.

#### Vector de Ataque y Escenario de Explotación:
1.  Un operador tenía permiso asignado sobre la máquina `110`.
2.  El administrador elimina la máquina `110` mediante la interfaz.
3.  Más adelante, el administrador da de alta un nuevo servidor de base de datos sensible. Proxmox le asigna el ID `110` porque es el menor VMID libre disponible.
4.  El operador abre su panel y descubre que tiene control total sobre el nuevo servidor de base de datos sin que nadie se lo haya otorgado.

#### Importancia:
Falla de control de acceso basada en objetos (IDOR) que vulnera la confidencialidad e integridad de nuevas instancias aprovisionadas.

---

### Brecha 4: Cambio Voluntario de Contraseña no Revoca Sesiones Activas
*   **Fix Asociado:** `FIX-114`
*   **Clasificación:** 🟠 **ALTA** | CWE-613 (Insufficient Session Expiration) | OWASP A07:2021
*   **Archivos Afectados:**
    *   [`backend/internal/core/services/user_service.go`](../../backend/internal/core/services/user_service.go#L443-L483) (`CambiarContrasena`)

#### Cómo se Detectó:
Se contrastó el método `ConfirmarRecuperacionContrasena` (que sí invoca `revocarSesionesDeUsuario`) con el método de cambio voluntario en `user_service.go:443`:
```go
// user_service.go:469
if err := s.userRepo.ActualizarUsuario(ctx, usuarioID, map[string]any{
    "contrasena_hash":   hash,
    "cambio_contrasena": false,
}); err != nil {
    return err
}
// Registra auditoría y retorna nil sin tocar sesiones_activas ni Redis
```
El método actualiza el hash de la clave pero **no invalida ninguna sesión activa en la tabla `sesiones_activas` ni en Redis**.

#### Vector de Ataque y Escenario de Explotación:
1.  Un usuario advierte actividad extraña en su cuenta y procede de inmediato a cambiar su contraseña desde su equipo personal.
2.  El atacante, que posee una sesión iniciada en otro dispositivo, no se ve afectado: su token JWT de acceso (válido por hasta 8 horas) y su cookie de refresh token (válida por 30 días) continúan operando con normalidad.

#### Importancia:
Invalida el propósito defensivo primordial del cambio de contraseñas ante sospechas de intrusión, permitiendo la persistencia no autorizada del atacante.

---

### Brecha 5: Límite Máximo Defectuoso de 12 Caracteres en Contraseñas
*   **Fix Asociado:** `FIX-115`
*   **Clasificación:** 🟠 **MEDIA / ALTA** | CWE-521 (Weak Password Requirements) | Violación de NIST SP 800-63B
*   **Archivos Afectados:**
    *   [`backend/internal/infrastructure/crypto/password.go`](../../backend/internal/infrastructure/crypto/password.go#L93-L95) (Líneas 93-95)
    *   [`backend/internal/adapters/primary/http/auth_handler.go`](../../backend/internal/adapters/primary/http/auth_handler.go#L280) (Línea 280)
    *   [`backend/internal/core/ports/user_port.go`](../../backend/internal/core/ports/user_port.go#L142) (Línea 142)

#### Cómo se Detectó:
En `password.go:93`:
```go
if len(contrasena) > 12 {
    return fmt.Errorf("la contraseña no puede superar los 12 caracteres")
}
```
Esta misma regla está declarada en las etiquetas de validación de los DTOs de Gin:
```go
binding:"required,min=8,max=12"
```

#### Vector de Ataque e Impacto:
1.  **Imposibilidad de usar Gestores de Contraseñas:** Gestores estándar de la industria generan por defecto cadenas de 16 a 64 caracteres; todas son rechazadas.
2.  **Imposibilidad de usar Passphrases:** No se pueden utilizar frases compuestas de varias palabras.
3.  **Violación de Estándares Internacionales:** Las directrices de NIST SP 800-63B prohíben expresamente truncar o fijar cotas máximas cortas en contraseñas (exigen admitir contraseñas de al menos 64 caracteres).
4.  **Debilitamiento ante Cracking Offline:** Reduce drásticamente el espacio de búsqueda computacional en caso de filtración de hashes bcrypt.

#### Importancia:
Fomenta contraseñas cortas y predecibles en un sistema que administra infraestructura de hipervisores críticos.

---

### Brecha 6: Riesgo de Inyección de Parámetros en Configuración de Proxmox (`net0`)
*   **Fix Asociado:** `FIX-116`
*   **Clasificación:** 🟠 **MEDIA** | CWE-88 (Improper Neutralization of Argument Delimiters in a Command) | CWE-20
*   **Archivos Afectados:**
    *   [`backend/internal/adapters/secondary/proxmox/client.go`](../../backend/internal/adapters/secondary/proxmox/client.go#L240-L249) (Líneas 240-249)
    *   [`backend/internal/adapters/primary/http/instance_handler.go`](../../backend/internal/adapters/primary/http/instance_handler.go#L350-L389) (`validarCreacionNormalizada`)

#### Cómo se Detectó:
En `client.go:240`:
```go
net0 := fmt.Sprintf("name=eth0,bridge=%s,firewall=%s", input.Bridge, firewall)
```
La API REST de Proxmox toma el parámetro `net0` y lo interpreta como una lista de directivas clave-valor separadas por comas.
Al auditar `validarCreacionNormalizada` en `instance_handler.go:371`:
```go
if input.Bridge == "" || input.TemplateID == "" {
    return ports.CrearInstanciaProxmoxInput{}, fmt.Errorf("bridge y template son obligatorios")
}
```
El backend verifica que `input.Bridge` no esté vacío, pero **no valida sus caracteres ni sanitiza comas**.

#### Vector de Ataque y Escenario de Explotación:
Si un usuario con permisos de creación introduce en el campo bridge:
`vmbr0,rate=50,tag=100`
La cadena interpolada resultante es:
`name=eth0,bridge=vmbr0,rate=50,tag=100,firewall=1`
Proxmox deserializa esas directivas inyectadas alterando la configuración de red física del contenedor (modificando VLANs, ancho de banda o interfaces) sin pasar por el esquema de validación formal de Centinela.

#### Importancia:
Permite la configuración indebida de parámetros de red de bajo nivel en el hipervisor anfitrión.

---

### Brecha 7: Exposición Indebida de Métricas y Topología Global del Host a Operadores
*   **Fix Asociado:** `FIX-117`
*   **Clasificación:** 🟡 **MEDIA** | CWE-200 (Exposure of Sensitive Information to an Unauthorized Actor)
*   **Archivos Afectados:**
    *   [`backend/cmd/api/main.go`](../../backend/cmd/api/main.go#L292-L295) (Líneas 292-295)
    *   [`backend/internal/adapters/primary/http/node_handler.go`](../../backend/internal/adapters/primary/http/node_handler.go#L47-L79) (Líneas 47-79)

#### Cómo se Detectó:
En `main.go`:
```go
api.GET("/node/status", middleware.RequireAuth(authService), nodeHandler.ObtenerEstado)
api.GET("/node/history", middleware.RequireAuth(authService), nodeHandler.ObtenerHistorial)
```
Ambos endpoints están protegidos únicamente por `RequireAuth()`, sin ningún filtro de rol `RequireRole("ADMIN")` ni segregación de contenido.
En `node_handler.go`, la respuesta contiene `instancesSummary`, que expone el recuento total de todas las máquinas virtuales y contenedores del hipervisor físico, junto con el desglose exacto de cuántas están encendidas, apagadas y pausadas, además del almacenamiento total y libre de los discos locales.

#### Vector de Ataque e Impacto:
Un operador con acceso únicamente a un contenedor de pruebas puede monitorear de forma continua la carga de CPU/RAM de todo el servidor físico, identificar cuándo se encienden o apagan máquinas de otras áreas y deducir períodos de alta actividad, facilitando tareas de reconocimiento interno para planificar ataques de denegación de servicio.

#### Importancia:
Fuga de información de infraestructura física hacia usuarios no privilegiados.

---

### Brecha 8: Credenciales por Defecto Hardcodeadas en Migraciones de Base de Datos
*   **Fix Asociado:** `FIX-118`
*   **Clasificación:** 🟡 **MEDIA / BAJA** | CWE-798 (Use of Hard-coded Credentials)
*   **Archivos Afectados:**
    *   [`backend/internal/adapters/secondary/postgres/db.go`](../../backend/internal/adapters/secondary/postgres/db.go#L434-L480) (`seedAdminPorDefecto`)

#### Cómo se Detectó:
En `db.go:455-478`:
```go
hash, _ := crypto.HashContrasena("Admin123!")
admin := domain.Usuario{
    EmailUsuario:   "admin@elcentinela.com",
    ContrasenaHash: hash,
    Rol:            "ADMIN",
    ...
}
```
Si la tabla de usuarios está vacía, el servidor inserta este usuario administrador y vuelca las credenciales en texto plano a los logs:
```
log.Println("   👉 Email: admin@elcentinela.com")
log.Println("   👉 Clave: Admin123!")
```

#### Vector de Ataque e Impacto:
Si una base de datos se despliega en producción o staging vacía sin ejecutar el seed interactivo, la cuenta queda activa con credenciales universalmente legibles en el código del repositorio.

#### Importancia:
Riesgo de compromiso total de la infraestructura en despliegues desatendidos o automatizados.

---

### Brecha 9: Asignación a Ciegas de VMIDs Inexistentes o Protegidos en Permisos
*   **Fix Asociado:** `FIX-119`
*   **Clasificación:** 🟡 **MEDIA** | CWE-20 (Improper Input Validation)
*   **Archivos Afectados:**
    *   [`backend/internal/adapters/primary/http/user_handler.go`](../../backend/internal/adapters/primary/http/user_handler.go#L285-L310) (`AsignarPermisos`)
    *   [`backend/internal/core/services/user_service.go`](../../backend/internal/core/services/user_service.go#L268-L286)

#### Cómo se Detectó:
Al inspeccionar `AsignarPermisos` en `user_service.go`:
```go
func (s *userServiceImpl) AsignarPermisos(ctx context.Context, usuarioID, orgID uuid.UUID, actorID uuid.UUID, permisos []ports.PermisoInstanciaInput) error {
    ...
    if err := s.userRepo.ReemplazarPermisos(ctx, usuarioID, permisos); err != nil {
        return err
    }
}
```
El servicio persiste cualquier número entero provisto en `permisos`.
No verifica:
1.  Si el VMID existe realmente en Proxmox VE.
2.  Si el VMID corresponde a una máquina protegida de infraestructura del clúster (`PROXMOX_PROTECTED_VMIDS`, ej. 100, 101, 102).

#### Vector de Ataque e Impacto:
Permite que un operador visualice en su inventario máquinas de la infraestructura crítica de Centinela (servidor de base de datos, proxies Nginx, etc.), consulte su estado de telemetría y establezca conexiones WebSocket hacia sus métricas internas.

#### Importancia:
Permite la exposición de recursos críticos protegidos y ensucia la base de datos con identificadores inválidos.

---

## 4. Matriz Resumen de Prioridades y Plan de Remediación

| ID Fix | Vulnerabilidad Clave | Ubicación en Código | Severidad | Esfuerzo Estimado |
| :--- | :--- | :--- | :---: | :---: |
| **`FIX-111`** | Secuestro de cuenta por cambio de email sin clave ni 2FA | `account_handler.go` / `user_service.go` | 🔴 **CRÍTICA** | 2,5 h |
| **`FIX-112`** | Ausencia total de Rate Limiting en autenticación | `main.go` / `auth_service.go` | 🔴 **CRÍTICA** | 2,5 h |
| **`FIX-113`** | Permisos huérfanos por reciclaje de VMID (`BAC-24C`) | `instance_handler.go` / `seguimiento_tareas.go` | 🔴 **ALTA** | 1,5 h |
| **`FIX-114`** | Cambio voluntario de contraseña no revoca sesiones | `user_service.go` | 🟠 **ALTA** | 1,0 h |
| **`FIX-115`** | Cota máxima defectuosa de 12 caracteres en contraseñas | `password.go` / DTOs | 🟠 **MEDIA / ALTA** | 1,0 h |
| **`FIX-116`** | Inyección de directivas de hardware en Proxmox (`net0`) | `client.go` / `instance_handler.go` | 🟠 **MEDIA** | 2,0 h |
| **`FIX-117`** | Exposición de telemetría global del host a operadores | `node_handler.go` / `main.go` | 🟡 **MEDIA** | 2,0 h |
| **`FIX-118`** | Credenciales por defecto fijas en migraciones (`db.go`) | `db.go` | 🟡 **MEDIA / BAJA** | 1,0 h |
| **`FIX-119`** | Asignación de VMIDs protegidos o inexistentes en permisos | `user_service.go` | 🟡 **MEDIA** | 1,5 h |

---

## 5. Conclusión y Recomendación

El sistema cuenta con bases criptográficas adecuadas en aspectos específicos (hashing con bcrypt cost 12, cifrado AES-256-GCM para secretos TOTP y tablas append-only en PostgreSQL). Sin embargo, **la ausencia de controles de flujo esenciales (rate limiting, re-autenticación en cambios críticos y sincronización del ciclo de vida de infraestructura con permisos) genera brechas que comprometen la seguridad general de la plataforma**.

La ejecución de las tareas `FIX-111` a `FIX-119` detalladas en este informe y especificadas en [`documentacion/futuro2.md`](../futuro2.md) resolverá de forma definitiva estas vulnerabilidades, transformando a El Centinela en un sistema blindado y preparado para entornos de producción.
