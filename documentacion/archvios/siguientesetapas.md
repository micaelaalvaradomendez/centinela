Una vez que terminen **`futuro.md`** (con la prueba integral `LOGIN-04` aprobada y el despliegue de base de datos/proxy de `INF-03`/`INF-04` funcionando), el sistema tiene los cimientos completamente resueltos: autenticación segura, 2FA persistido, usuarios, roles, bloqueo por recursos, auditoría append-only (`BAC-18`) y el contrato base del canal de eventos (`BAC-21`).

**RF-12 (organizaciones / multi-tenant) queda descartado del MVP por decisión del equipo.** Ninguna tarea de las tres etapas siguientes debe reservar `organization_id` en ninguna tabla; si se retoma en el futuro, va a requerir migración y backfill de todo lo ya construido.

**La integración con un servidor SMTP real queda formalmente descartada del alcance del proyecto.** La entrega de correos queda resuelta de forma definitiva mediante el simulador por consola `MockEmailService` (desacoplado a través del puerto `EmailService`), cumpliendo con la política de mínimo privilegio (Zero-Trust), sanitización de logs y control estructurado de errores sin dependencia de infraestructura externa.

A partir de ahí, entran de lleno a **las funcionalidades operativas del hipervisor**. Para no improvisar y tener todo listo en ClickUp antes de que terminen la semana que viene, el trabajo que sigue debe estructurarse en **3 etapas consecutivas** hasta la entrega final:

---

### 🗺️ Visión General de las Etapas Posteriores

```mermaid
flowchart LR
    FUTURO[Cierre futuro.md: LOGIN-04] --> ETAPA1[Etapa 1: Core Operativo y MVP]
    ETAPA1 --> ETAPA2[Etapa 2: Telemetría y Aprovisionamiento]
    ETAPA2 --> ETAPA3[Etapa 3: Hardening, Auditoría y Deploy Final]

```

---

### 🚀 ETAPA 1: Core Operativo de Proxmox y MVP (El corazón del sistema)

> **Objetivo:** Cumplir el hito del primer entregable funcional. Que un usuario pueda entrar, ver la salud del nodo físico, listar sus máquinas reales y operarlas (prender/apagar) sin romper nada. Esta etapa asume **un único nodo Proxmox fijo**; el soporte multi-nodo (`RNF-07`) no está incluido y, si se necesita para el MVP, requiere modelar el nodo en la base de datos, un selector en el frontend y desacoplar la lógica de servicios del nodo fijo antes de empezar esta etapa.
> 
> 
> **Requerimientos cubiertos:** `RF-02` (Dashboard Host), `RF-03` (Inventario completo), `RF-04` (Ciclo de Vida con UPID) y `RF-11` de forma parcial (notificación de fin de tarea; el motor de alertas por saturación se completa recién en la Etapa 3).
> 
> 

#### 1. Backend (`Tayra` y `Lisandro`):

* **Adaptador de Telemetría del Nodo (`RF-02`):**
* Endpoint `GET /api/node/status` consumiendo `/nodes/{node}/status` de la API de Proxmox.


* Transformación y normalización de métricas (cálculo de % de CPU, conversión de bytes a GB en RAM y disco, uptime). Caché corta en Redis (no en memoria del proceso, para que sobreviva a múltiples réplicas del backend) para no saturar al hipervisor si hay múltiples consultas; el despliegue de Redis se agrega como tarea de infraestructura de esta etapa (ver punto 3).




* **Inventario Unificado Definitivo (`RF-03`):**
* Evolucionar `BAC-14` hacia el endpoint formal `GET /api/instances` que unifique VMs (QEMU) y Contenedores (LXC).


* Debe incluir: ID, nombre, tipo, estado (`running`/`stopped`), IP asignada (vía QEMU Guest Agent) y consumo actual.


* Aplicación estricta del filtrado RBAC: el Operador solo recibe su subconjunto de `user_instances`.




* **Motor de Tareas Asíncronas y Energía (`RF-04` / Crítico):**
* Endpoints de control: `POST /api/instances/{id}/status/{action}` (`start`, `shutdown`, `stop`, `reboot`).


* Captura obligatoria del string **UPID** devuelto por Proxmox.


* *Worker / Task Poller* en segundo plano que sondee periódicamente `/nodes/{node}/tasks/{upid}/status` hasta que termine o falle.


* Servicio de notificación (WebSockets o SSE) que emita el evento de fin de tarea hacia el frontend, usando el esquema de evento genérico ya definido en la fase base (`BAC-21`), porque `RF-11` lo reutiliza en la Etapa 3 para alertas de saturación.


* `DELETE /api/instances/{id}` para eliminar una instancia (VM o LXC), tratada como acción destructiva con confirmación explícita en el frontend y seguimiento asíncrono por UPID igual que el resto de las acciones de ciclo de vida.





#### 2. Frontend (`Belinda`, `Luz` y `Cristian`):

* **Dashboard de Salud Global (`RF-02`):**
* Tarjetas de métricas (gauges / barras de progreso) de CPU, RAM, Almacenamiento y Uptime general del servidor host.


* Semáforo o badge de estado global del hipervisor (Normal / Advertencia / Inaccesible).




* **Vista de Inventario Filtrable (`RF-03`):**
* Tabla interactiva de instancias con badges de estado con código de colores (Verde = *Running*, Rojo/Gris = *Stopped*).


* Filtros dinámicos por nombre, tipo (VM vs LXC) y estado.




* **Controles de Energía Antierror (`RF-04`):**
* Botones de acción rápida con **modales de confirmación explícita** para evitar apagados accidentales.


* Patrón de interfaz *"Operación en progreso"*: al pulsar una acción, el botón entra en estado de carga (spinner), se bloquean los demás botones de esa máquina para evitar órdenes duplicadas y se escucha la resolución del UPID para mostrar un Toast de éxito/fallo.





#### 3. Infraestructura (`Nico` y `Lucas`):

* Asignar los permisos específicos `VM.PowerMgmt` y `Sys.Audit` al API Token de Proxmox con separación de privilegios activa.


* Validar la resolución de nombres y red interna entre el Backend y la API de Proxmox sobre `vmbr1`.


* Desplegar Redis (u otro almacén compartido equivalente) para la caché de estado del nodo, de modo que no dependa de la memoria de un único proceso backend.



> 🎯 **Prueba de Integración Etapa 1 (Hito MVP):**
> Un operador entra al sistema, ve solo las 2 máquinas que tiene asignadas en el inventario. Le da "Start" a una VM: la UI se bloquea con spinner, el Back recibe el UPID de Proxmox, espera la confirmación en segundo plano y la máquina pasa a estado verde con un Toast de éxito en el navegador.
> 
> 

---

### 📈 ETAPA 2: Telemetría Avanzada, Aprovisionamiento y Edición

> **Objetivo:** Dotar a la plataforma de monitoreo en tiempo real y capacidad de gestión de recursos con validación de cuotas.
> 
> 
> **Requerimientos cubiertos:** `RF-05` (Métricas por instancia), `RF-07` (Wizard de creación) y `RF-10` (Edición de recursos).
> 
> 

#### 1. Backend:

* **Streaming de Métricas (`RF-05`):**
* Canal WebSocket `/ws/metrics/{id}` que consuma periódicamente la serie temporal / RRDdata de Proxmox (CPU, RAM, Disco, Red) y la transmita sin sobrecargar la API REST.


* **Decisión de retención:** definir si el histórico corto (2h/24h) se sirve directamente desde el RRDdata que ya mantiene Proxmox o si el sistema necesita almacenamiento propio (por ejemplo, TimescaleDB) para consultas de tendencias más flexibles. Esta decisión condiciona el diseño del endpoint y debe tomarse antes de programarlo.




* **Aprovisionamiento Asistido con Control de Cuotas (`RF-07`):**
* Endpoint `POST /api/instances` para crear VMs o LXCs.


* **Algoritmo de validación de cuotas:** antes de enviar la orden a Proxmox, comprueba la RAM y disco disponibles en el host. Si los recursos no alcanzan, rechaza inmediatamente con `HTTP 409 Conflict` evitando sobreaprovisionar el servidor.


* Asocia automáticamente la nueva instancia creada a la matriz de permisos del creador (`user_instances`).




* **Edición Dinámica de Recursos (`RF-10`):**
* Endpoint `PUT /api/instances/{id}/config` para editar vCPU y RAM.


* Cálculo diferencial ($\Delta\text{RAM}$): valida que el host soporte el incremento solicitado antes de aplicar el cambio.





#### 2. Frontend:

* **Gráficos en Tiempo Real (`RF-05`):**
* Integración de componentes gráficos (Chart.js / Recharts) conectados al canal de WebSockets para mostrar histórico reciente y picos de consumo.




* **Wizard Guiado de 4 Pasos (`RF-07`):**
* Paso 1: Tipo (VM / LXC) $\rightarrow$ Paso 2: Recursos (Sliders de CPU/RAM con barras de cuota disponible en tiempo real) $\rightarrow$ Paso 3: Template/ISO y Red $\rightarrow$ Paso 4: Resumen y confirmación.


* **Validación de red en el Paso 3:** definir si la IP y la puerta de enlace se toman de un pool DHCP administrado por Proxmox o si el backend valida y reserva la IP antes de crear la instancia. Sin esta definición, el paso 3 no puede implementarse de forma consistente con `RF-07`.




* **Modal de Edición de Recursos (`RF-10`):**
* Ajuste de parámetros con advertencias visuales si el cambio exige reinicio de la máquina.





#### 3. Infraestructura:

* Mantener plantillas LXC (ej. Alpine/Debian) e imágenes ISO base optimizadas en los storages de Proxmox.


* Otorgar permisos `VM.Allocate` y `Datastore.AllocateSpace` al token de Proxmox.



---

### 🛡️ ETAPA 3: Auditoría Inmutable, Snapshots, Hardening y Cierre Final

> **Objetivo:** Completar los requerimientos de trazabilidad y respaldo, cerrar las notificaciones por saturación de recursos, auditar la seguridad y dejar el entorno listo para la defensa final.
> 
> 
> **Requerimientos cubiertos:** `RF-06` (Snapshots), `RF-08` (Audit Log inalterable), `RF-11` de forma completa (alertas por saturación y correo) y `RNF-01` a `RNF-06`.
> 
> 

#### 1. Backend:

* **Registro Inmutable de Auditoría (`RF-08`):**
* La tabla `audit_logs` y el interceptor de acciones administrativas (usuarios, roles, permisos, resets) ya existen desde la fase base (`BAC-18`, ver `futuro.md`), para no necesitar backfill. Acá solo falta ampliar el interceptor para cubrir las acciones operativas de Proxmox (energía, creación, borrado, UPID, resultado OK/Fallo).


* Endpoint paginado y filtrable `GET /api/admin/audit-logs`.


* Endpoint de exportación que respete los filtros aplicados.


* **Motor de alertas por saturación (`RF-11`):** evaluador periódico de la salud/consumo del nodo (CPU, RAM, disco) que dispare eventos de alta severidad sobre el mismo canal definido en `BAC-21`, con umbrales configurables y emisión a través del canal en tiempo real y el simulador de correo `MockEmailService`.




* **Puntos de Restauración / Snapshots (`RF-06`):**
* Endpoints para crear (`POST`), listar (`GET`) y revertir (`POST /rollback`) snapshots asociados al ciclo asíncrono de Proxmox.





#### 2. Frontend:

* **Vista de Auditoría (`RF-08`):**
* Tabla exclusiva para Administradores con filtros por usuario, fecha, máquina y tipo de evento.




* **Gestor de Snapshots (`RF-06`):**
* Pestaña dentro del detalle de cada instancia con alertas explícitas en rojo para acciones destructivas (revertir snapshot).





#### 3. Infraestructura y Equipo Completo (Cierre):

* **Configuración Final de Producción:** Nginx reverse proxy con terminación SSL/TLS y soporte completo de encabezados `Upgrade` para WebSockets sin caídas de sesión.


* **Smoke Test y UAT:** Simulación integral punta a punta de todos los casos de uso para la defensa académica.


* **Integración con proveedor SMTP real (Continuación de BAC-16 / Si alcanza el tiempo):**
  - Conectar el backend a un servidor SMTP real (SendGrid, Mailgun o servidor propio) mediante un nuevo adaptador `SmtpEmailService` que reemplace al actual `MockEmailService`.
  - Inyección de variables de entorno reales (`SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, etc.).
  - *Garantía de Arquitectura Hexagonal:* Gracias al puerto `EmailService`, ni la capa de dominio (`AuthService`, `UserService`) ni los handlers HTTP sufrirán ningún cambio; solo bastará con crear el adaptador SMTP y conectarlo en el `main.go`. Si el tiempo no alcanza antes de la entrega final, se deja como se implementó en la fase base (con el simulador de consola `MockEmailService`).


* **Evaluación de almacenamiento seguro de tokens (Mejora post FIX-10):**
  - Se implementó el almacenamiento de los JWT mediante `sessionStorage`.
  - Como mejora de seguridad para evaluar en esta etapa, analizar el uso de cookies `HttpOnly` y `Secure`, principalmente para el `refreshToken`, ya que evita que JavaScript pueda acceder directamente a este token y reduce su exposición ante ataques XSS. Esto requeriría adaptar el manejo del endpoint de refresh en el backend para emitir y leer la cookie.

---

### 📋 Cuadro Resumen de Carga para ClickUp

| Etapa | Foco Funcional | Backlog Backend | Backlog Frontend | Infraestructura |
| --- | --- | --- | --- | --- |
| **Etapa 1** *(Inmediata tras futuro.md)* | **Core MVP Proxmox** (`RF-02`, `RF-03`, `RF-04`, `RF-11` parcial)

 | Poller del host, unificación DTOs instancias, endpoints de energía y borrado, worker de sondeo UPID sobre el canal de eventos de `BAC-21`.

 | Dashboard del nodo, tabla de instancias con filtros, botones antierror y toasts.

 | Permisos `VM.PowerMgmt` y `Sys.Audit` en token Proxmox, despliegue de Redis para caché de estado.

 |
| **Etapa 2** | **Telemetría y Cuotas** (`RF-05`, `RF-07`, `RF-10`)

 | WebSockets de métricas, validación de cuotas libres en host y creación de instancias.

 | Gráficos en vivo (Chart.js), Wizard de 4 pasos con barras de cuota y modal de edición.

 | Storage con ISOs y plantillas LXC livianas.

 |
| **Etapa 3** | **Seguridad y Cierre** (`RF-06`, `RF-08`, `RF-11`, Deploy)

 | Ampliación del interceptor de `audit_logs` a acciones operativas, endpoint de consulta/exportación, motor de alertas por saturación con SMTP, endpoints de snapshots asíncronos.

 | Visor de auditoría con filtros, modales destructivos de snapshot.

 | Nginx con TLS/WSS enrutado y ensayo de demo en vivo.

 |


---

## 🔍 Evaluación de Cobertura Integral y Tareas Faltantes a Incorporar

> **Auditoría contra requerimientos.md y centinela.md:** Esta sección evalúa la cobertura de la visión del sistema y agrega formalmente las tareas y precisiones necesarias para cerrar el 100% de los requerimientos funcionales y no funcionales sin modificar el texto histórico anterior.

### 1. Matriz de Cobertura de Requerimientos del Sistema

| Requerimiento | Descripción | Estado en el Proyecto | Implementación / Planificación |
|---|---|:---:|---|
| **RF-01** | Autenticación y Control de Acceso (JWT, RBAC `ADMIN`/`OPERATOR`, 2FA TOTP con QR y secreto manual) | ✅ **100% Cubierto** | Terminado en fase base (`BAC-01` a `BAC-04`, `BAC-08` a `BAC-10`, `LOGIN-01` a `LOGIN-04`). Revinculación/reset de 2FA asistido por Admin (`BAC-13`, `FIX-19`). Hardening con cookies `HttpOnly` en `etapa3.md` (`BAC-48`, `FRN-34`). |
| **RF-02** | Dashboard de Estado del Nodo Host (CPU %, núcleos, RAM, Disco, Uptime, conteo VMs/LXC, semáforo, caché Redis) | ✅ **100% Cubierto** | En `terminado-1.md` (`BAC-22`, `FRN-19A`), `actual.md` (`FRN-19C`, `FIX-78`, `FIX-83`) y `etapa1.md`. |
| **RF-03** | Inventario Unificado de Instancias (VMs/LXC, IP QEMU Agent, CPU/RAM, filtros dinámicos, RBAC estricto) | ✅ **100% Cubierto** | En `terminado-1.md` (`BAC-23A`, `FRN-20A`), `actual.md` (`FRN-20C`, `FIX-43`, `FIX-63`) y `etapa1.md`. |
| **RF-04** | Ciclo de Vida de Instancias (Start, Shutdown, Stop, Reboot, Delete, modales antierror, worker UPID, bloqueo mutuo) | ✅ **100% Cubierto** | En `terminado-1.md` (`BAC-25A`, `BAC-29`), `actual.md` (`BAC-24C`, `FRN-15`, `FRN-16`, `FRN-17B`) y `etapa1.md`. |
| **RF-05** | Monitor de Métricas por Instancia (WebSocket en vivo, RRDdata, CPU/RAM/Red/Disco, histórico 2h/24h) | ✅ **100% Cubierto** | En `actual.md` (`BAC-34`, `BAC-35`) y `etapa2.md` (`BAC-36`, `FRN-22`, `FRN-23`). |
| **RF-06** | Gestión de Snapshots (Crear, listar, eliminar, rollback destructivo con tipeo de seguridad, lock `409 INSTANCE_BUSY`) | ✅ **100% Cubierto** | En `etapa3.md` (`BAC-42`, `BAC-43`, `FRN-29`, `FRN-30`). |
| **RF-07** | Asistente de Creación de Instancias (Wizard 4 pasos VM/LXC, validación transaccional de cuotas, error 409, autoasignación) | ✅ **100% Cubierto** | En `actual.md` (`BAC-33`, `BAC-37`, `BAC-38`, `INF-10`) y `etapa2.md` (`BAC-39`, `FRN-24` a `FRN-26`). |
| **RF-08** | Registro de Auditoría Inmutable (Append-only en `auditoria`, particionado por fecha, filtros combinados, exportación streaming CSV/JSON) | ✅ **100% Cubierto** | En `terminado.md` (`BAC-18`), `etapa2.md` (`BAC-41`) y `etapa3.md` (`BAC-44`, `BAC-45`, `FRN-31`, `FRN-32`). |
| **RF-09** | Gestión de Usuarios y Permisos (CRUD usuarios, roles, permisos por instancia `FULL_ACCESS`/`READ_ONLY`, clave temporal y reset) | ✅ **100% Cubierto** | En `terminado.md` (`BAC-06`, `BAC-07`, `BAC-10`, `BAC-15`, `FRN-05` a `FRN-07`, `SEC-04`). |
| **RF-10** | Edición Dinámica de Recursos (vCPU y RAM, cálculo diferencial $\Delta\text{RAM}$, hot-plug vs reinicio) | ✅ **100% Cubierto** | En `etapa2.md` (`BAC-40`, `FRN-27`, `FRN-28`). |
| **RF-11** | Notificaciones y Alertas (Toasts reactivos, `TASK_FINISHED`, alertas por saturación con histéresis, correo y enlace a `/instances/:vmid`) | ✅ **100% Cubierto** | En `terminado.md`, `actual.md` (`BAC-26`, `FRN-21`) y `etapa3.md` (`BAC-46`, `BAC-47`, `FRN-33`). |
| **RF-12** | Gestión de Organizaciones (Multi-tenant) | ⚠️ **Descartado del MVP** | Descarte formal acordado por el equipo (`futuro.md`). El sistema opera como **mono-organización**. |
| **RF-13** | Recuperación de Contraseña por Usuario (OTP de 6 dígitos al correo, validación con expiración e invalidación) | ✅ **100% Cubierto** | En `terminado.md` (`BAC-16`, `BAC-17`, `FRN-12`, `FRN-13`). |
| **RNF-01** | Seguridad (Mínimo privilegio API Proxmox, bcrypt, cifrado AES-256 de secretos TOTP, cookies HttpOnly, Rate Limiting, TLS) | ✅ **100% Cubierto** | En `terminado.md` (`BAC-02`, `BAC-10`) y `etapa3.md` (`BAC-48`, `INF-11`). |
| **RNF-02** | Rendimiento (Respuestas $\le$ 5s, Go concurrente, caché Redis de telemetría y tickets, streaming HTTP en exportación) | ✅ **100% Cubierto** | En `terminado-1.md` (`BAC-22`), `etapa2.md` y `etapa3.md` (`BAC-45`). |
| **RNF-03** | Usabilidad (UI moderna, abstracción de complejidad de Proxmox, catálogo D2 de errores amigables, modales antierror) | ✅ **100% Cubierto** | En `terminado.md` (`FIX-08`), `actual.md` (`FIX-77`, `FIX-86`), `etapa2.md` y `etapa3.md`. |
| **RNF-04** | Confiabilidad (Tolerancia a caídas, UPID worker asíncrono, timeout 504 vs caído 502, auto-restart en Compose) | ✅ **100% Cubierto** | En `terminado.md` (`FIX-40`), `terminado-1.md` (`BAC-25A`), `etapa1.md` y `etapa3.md`. |
| **RNF-05** | Mantenibilidad (Arquitectura hexagonal, código testeable, **Documentación interactiva Swagger UI / OpenAPI**) | 🔄 **Incorporado a Etapa 3** | Arquitectura cubierta; se formaliza la tarea de publicación interactiva de Swagger UI (`BAC-49`). |
| **RNF-06** | Compatibilidad (Navegadores modernos, CORS restringido, proxy Nginx con upgrades WSS y SSE) | ✅ **100% Cubierto** | En `terminado.md`, `etapa1.md` (`FIX-31`) y `etapa3.md` (`INF-11`). |
| **RNF-07** | Escalabilidad (Soporte multi-nodo Proxmox) | ⚠️ **Postergado post-MVP** | Alcance del MVP limitado a un **nodo físico único** (`PROXMOX_NODE`). No bloquea la entrega. |
| **Exclusión** | Consola Remota Web (noVNC / xterm.js / SSH en navegador) | 🚫 **Exclusión de Seguridad** | Excluida formalmente en `requerimientos.md` y `centinela.md` para evitar vectores de ejecución remota. |

---

### 2. Precisiones y Rectificaciones sobre el Plan Original

1. **Alineación del Servidor SMTP Real (`INF-08` / `BAC-16B` / `BAC-47`):**
   - La nota preliminar de este documento indicaba el descarte de SMTP. Dicha decisión fue rectificada en `futuro.md` y `etapa3.md`:
   - El puerto hexagonal `EmailService` utiliza `MockEmailService` únicamente para pruebas unitarias/CI.
   - En desarrollo y producción final se activa el adaptador real `SmtpEmailService` para credenciales temporales (`BAC-06`), códigos OTP de recuperación de contraseña (`RF-13`) y despacho de alertas críticas de saturación (`RF-11` / `BAC-47`).
2. **Modelo Oficial de Base de Datos (Sin nombres provisorios):**
   - La tabla de auditoría no se denomina `audit_logs` sino **`auditoria`**, estructurada como *append-only* y particionada trimestralmente por rango de fechas (`PARTITION BY RANGE (fecha_hora)`).
   - La tabla de permisos de instancia no se denomina `user_instances` sino **`permisos_instancia`**, con niveles de acceso `FULL_ACCESS` y `READ_ONLY`.
   - Las tareas asíncronas se persisten formalmente en **`tareas_asincronas`**.
3. **Flujo de Revinculación de 2FA (`RF-01`):**
   - La solicitud de regeneración o revinculación de 2FA por extravío del dispositivo autenticador se resuelve mediante **restablecimiento administrativo** (`POST /api/admin/users/:id/2fa/reset` / `BAC-13` / `FIX-19`) para evitar vulnerabilidades de suplantación en endpoints públicos no autenticados.

---

### 3. Tareas Faltantes Agregadas para Cubrir el 100% de los Requerimientos

Para cubrir los requerimientos formales de `requerimientos.md` que no contaban con una tarea técnica individual en las etapas 1, 2 y 3, se agregan formalmente al backlog de la Etapa 3:

#### `BAC-49` - Publicación y Automatización de Swagger UI / OpenAPI Interactivo (`RNF-05`)
- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `BAC-45`, `BAC-48`.
- **Requerimiento:** `RNF-05` Backend (*"Configurar la generación automática de documentación interactiva de la API"*).
- **Entregable:**
  1. Integración de interfaz interactiva Swagger UI embebida bajo `/api/docs` o `/swagger/*`.
  2. Regeneración automatizada de las anotaciones OpenAPI para la totalidad de endpoints del sistema (Autenticación, Nodo, Instancias, Métricas, Creación, Snapshots, Auditoría y Eventos) con sus contratos de error D2.
  3. Capacidad de ejecutar pruebas en vivo contra el backend inyectando el token JWT o cookies de sesión.
- **Criterio de éxito:** Cualquier evaluador puede ingresar a `/api/docs`, consultar todos los contratos vigentes y probarlos interactivamente.

#### `INF-12` - Automatización de Respaldos Periódicos de Base de Datos y Script de Restauración (`RF-08` / `RNF-04`)
- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.5 h
- **Depende de:** `INF-03`, `INF-11`.
- **Requerimiento:** `RF-08` Infraestructura (*"Garantizar la persistencia segura y los respaldos periódicos de la base de datos... Aislar la base de datos para evitar modificaciones externas directas"*).
- **Entregable:**
  1. Contenedor o cron job de respaldo programado (`centinela-db-backup`) ejecutando `pg_dump` con compresión sobre un volumen montado.
  2. Política de retención con rotación (respaldos diarios retenidos por 7 días y respaldo semanal).
  3. Script de prueba de restauración documentado `scripts/restore-db.sh`.
- **Criterio de éxito:** El script de backup genera los volcados sin interrumpir la operación y `restore-db.sh` recupera íntegramente el esquema y datos (incluidas las particiones de `auditoria`) en una base limpia.

#### `FRN-35` - Enlaces Contextuales Universales desde Notificaciones hacia el Detalle de Instancia (`RF-11`)
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-21`, `FRN-33`.
- **Requerimiento:** `RF-11` Frontend (*"Cada notificación relacionada con una VM o LXC debe incluir un enlace hacia el detalle de la instancia y sus acciones disponibles"*).
- **Entregable:**
  1. Adaptación de las notificaciones flotantes (Toasts) y de los elementos del Centro de Notificaciones para incluir un botón o enlace directo `"Ver Instancia"` cuando la notificación tenga un `recursoId` asociado.
  2. Navegación directa hacia `/instances/:vmid`, activando la pestaña pertinente (ej. pestaña "Snapshots" si fue un evento de rollback, o "Métricas" si fue una alerta de saturación).
- **Criterio de éxito:** Al hacer clic en la notificación de una instancia, el usuario accede directamente a la vista de detalle correspondiente sin pasos intermedios.

---

### 4. Bloque Puente Consolidado (Etapa 2 $\rightarrow$ Etapa 3)

Para garantizar un paso sin fricciones técnicas entre el cierre de la Etapa 2 y el inicio de la Etapa 3 (replicando la estructura del Bloque 7 de la Etapa 1), se incorporaron formalmente en [`etapa2.md`](etapa2.md) y [`etapa3.md`](etapa3.md) las siguientes tareas puente:

| Tarea | Área | Responsable | Est. | Problema que resuelve y propósito como puente |
|---|---|---|---|---|
| `BAC-50` (`BRG-10-BAC`) | Backend / Test | Lisandro | 2.0 h | **Simulador Proxmox:** Soporte mock de endpoints de snapshots y UPIDs en memoria para permitir el desarrollo local de `BAC-42` y `BAC-43` sin depender del hipervisor físico. |
| `BAC-51` (`BRG-11-BAC`) | Backend | Tayra | 2.0 h | **Guarda Concurrente:** Middleware de exclusión mutua por instancia (`409 INSTANCE_BUSY`) para bloquear órdenes de energía y config mientras corre un snapshot o rollback. |
| `FRN-36` (`BRG-10-FRN`) | Frontend | Cristian | 2.0 h | **Layout Modular:** Reestructuración de `InstanceDetailPage.tsx` para admitir la pestaña de Snapshots (`FRN-29`) y badges de estados transitorios (`SNAPSHOTTING`/`ROLLING_BACK`). |
| `BAC-52` (`BRG-12-BAC`) | Backend | Tayra | 1.5 h | **Broker SSE Global:** Generalización de `/api/events` para transmitir eventos de alerta a nivel `HOST` (`ALERT_TRIGGERED`) a todos los administradores sin requerir un VMID obligatorio. |
| `INF-13` (`BRG-11-INF`) | Infraestructura | Nico | 1.5 h | **Storage Proxmox Real:** Verificación y habilitación de pools compatibles con snapshots nativos (`qcow2`, `lvm-thin`, `zfs`) en el Proxmox físico antes de la validación final. |
| `BAC-53` (`BRG-13-BAC`) | Backend / DB | Lisandro | 1.5 h | **Particionado Preventivo:** Creación preventiva de particiones trimestrales e índices compuestos en `auditoria` para habilitar el partition pruning y streaming de `BAC-45`. |
