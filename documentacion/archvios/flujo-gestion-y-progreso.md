# Flujo de Gestión, Progreso y Verificación Continua sobre los Submódulos

**Proyecto:** El Centinela (Orquestación y gestión de entornos virtualizados sobre Proxmox VE)  
**Período documentado:** 10/09/2026 – 19/10/2026 (73 commits en el repositorio orquestador)  
**Submódulos gestionados:** [`backend/`](../backend) (`centinela-back`) y [`frontend/`](../frontend) (`centinela-front`)

---

## 1. Propósito y Filosofía de Gestión

Este repositorio nació para resolver un desafío concreto de liderazgo técnico: **cómo coordinar, auditar y unificar el trabajo simultáneo de tres áreas (Backend, Frontend e Infraestructura) sin depender de reportes subjetivos ni frenar el desarrollo cuando el servidor físico o la otra capa aún no están listos**.

Para lograrlo, se diseñó una metodología apoyada en cuatro pilares:
1. **Planificación Atómica orientada a ClickUp:** Ninguna tarea entra al tablero como un bloque difuso o de larga duración. Todas las épicas se desglosaron en **micro-tareas de 0,5 h a 2,5 h** (máximo 4 h en el arranque inicial), con responsable asignado, estimación horaria, dependencias explícitas, archivos exactos a modificar y un criterio de éxito automatizable.
2. **Desacople *Local-First* mediante Contratos y Simuladores:** Antes de arrancar cada etapa, se congela el contrato JSON/HTTP/SSE/WebSocket entre capas y se construye un simulador de Proxmox VE en Docker/Go. De este modo, Frontend y Backend avanzan en paralelo sin depender del servidor Proxmox físico.
3. **Verificación Empírica Independiente (`/test`):** Las tareas marcadas como completadas por los desarrolladores en ClickUp **no se dan por cerradas por simple declaración**. Deben aprobar las suites de caja negra, integración E2E y seguridad adversarial alojadas en [`test/back`](../test/back) y [`test/front`](../test/front), que evalúan los submódulos desde afuera sin contaminar su código.
4. **Trazabilidad Tripartita (`actual` $\rightarrow$ `terminado` / `futuro` + `FIX`):** Cada corte de verificación registra el commit exacto (`SHA`) evaluado de cada submódulo, promueve a `terminado*.md` lo que pasa al 100% en verde, y ante cualquier desfasaje de contrato, caso de borde incumplido, placeholder o vulnerabilidad de seguridad, genera una tarea atómica **`FIX-XX`** con diagnóstico de causa raíz y línea exacta de código.

---

## 2. Metodología de Planificación y Carga en ClickUp

### 2.1. Estructura de Etapas y Bloques Puente (*Bridge Tasks*)
El proyecto se organizó en **4 fases evolutivas** conectadas mediante **Bloques Puente (`BRG-*`)**. La función de cada Bloque Puente es preparar la infraestructura, los simuladores y los contratos de la etapa siguiente mientras el equipo cierra los últimos detalles de la etapa actual:

| Fase / Etapa | Documentos de Gestión | Alcance Funcional Principal | Bloque Puente Asociado |
| :--- | :--- | :--- | :--- |
| **Fase Base** | [`terminado.md`](./terminado.md), [`futuro.md`](./futuro.md), [`resumenfasebasefaseuno.md`](./archvios/resumenfasebasefaseuno.md) | Arquitectura hexagonal en Go y React, PostgreSQL (`init.sql`), Login + 2FA TOTP (AES-256-GCM), sesión única en Redis (`1 sesión = 1 fila`), usuarios y permisos (`READ_ONLY` / `FULL_ACCESS`), recuperación de clave con OTP por SMTP (Mailpit STARTTLS), auditoría *append-only* con trigger inmutable, Nginx Edge TLS y cierre E2E **`LOGIN-04`**. | **Puente a Etapa 1 (`BAC-21B`, `BAC-21C`, `INF-07B`, `FRN-17C`, `BAC-28`):** Canal SSE `/api/events` con tickets efímeros en Redis, bus Pub/Sub, simulador Proxmox base y extensión de DTOs de instancias. |
| **Etapa 1** *(Visibilidad y Control Operativo)* | [`etapa1.md`](./etapa1.md), [`terminado-1.md`](./terminado-1.md), [`futuro-1.md`](./futuro-1.md) | Salud del nodo (`/api/node/status` con caché Redis y semáforo `D3`), inventario unificado QEMU/LXC con resolución concurrente de IPs por Guest Agent (`BAC-23B`), control de energía (`start`, `stop`, `shutdown`, `reboot`, `DELETE` con `BAC-24B`), worker pool acotado `UPID_WORKERS` con reanudación tras reinicio (`BAC-25C`), máquina de estados antierror en UI con catálogo **D2** e hitos **`INT-01`**, **`INT-02`** e **`INT-03`**. | **Puente a Etapa 2 (Ola 6: `BAC-30` a `BAC-33`, `FRN-21`, `INF-09`, `INF-07C`):** Simulador con storages/templates/RRDdata, endpoints `/node/storages` y `/node/templates`, worker UPID con soporte `create`, cálculo de `allocatable` del host, vista `/instances/:vmid` y upgrade WebSocket en Nginx. |
| **Etapa 2** *(Telemetría, Wizard y Escalado)* | [`etapa2.md`](./etapa2.md), [`terminado2.md`](./terminado2.md), [`futuro2.md`](./futuro2.md) | Streaming de métricas en vivo cada 2 s por WebSocket con ticket efímero (`BAC-35`, `FRN-22`) e histórico RRDdata `2h`/`24h` (`BAC-34`, `BAC-36`, `FRN-23`), motor de cuotas con 15% de reserva de RAM para el host y lock distribuido en Redis `lock:instance:provisioning` (`BAC-37`), Wizard de creación de 4 pasos (`BAC-38`, `BAC-39`, `FRN-24`..`26`) y edición de vCPU/RAM con $\Delta\text{RAM}$ y detección de *hot-plug* (`BAC-40`, `FRN-27`, `FRN-28`). | **Puente a Etapa 3 (Bloque 8: `BAC-50` a `BAC-53`, `FRN-36`, `INF-13`):** Simulador de snapshots/rollback (`BAC-50`), guarda de exclusión mutua `RequireInstanceIdle` (`BAC-51` $\rightarrow$ `409 INSTANCE_BUSY`), broker SSE para alertas `HOST` (`BAC-52`) y particionado preventivo de auditoría (`BAC-53`). |
| **Etapa 3** *(Snapshots, Exportación, Alertas y Hardening)* | [`etapa3.md`](./etapa3.md), [`futuro2.md`](./futuro2.md), [`actual.md`](./actual.md) | CRUD asíncrono y rollback de snapshots con confirmación destructiva escribiendo el nombre (`BAC-42`, `BAC-43`, `FRN-29`, `FRN-30`), exportación de auditoría CSV/JSON en streaming por cursor (`BAC-44`, `BAC-45`, `FRN-31`, `FRN-32`), evaluador de umbrales con **histéresis de 2 ciclos en Redis** (`BAC-46`, `BAC-47`, `FRN-33`) y blindaje integral de seguridad e identificadores (**`FIX-105` a `FIX-119`**). | Hitos finales de certificación E2E y producción (**`INT-07`**, **`INT-08`**, **`INT-09`**). |

### 2.2. Anatomía de las Micro-Tareas y Control de Tiempos en ClickUp
Cada tarjeta cargada en ClickUp desde los documentos de planificación sigue una plantilla estricta pensada para eliminar ambigüedades:
- **Identificador taxonómico:** `BAC-01`..`53` (Backend), `FRN-01`..`36` (Frontend), `INF-01`..`13` (Infraestructura), `SEC-01`..`04` (Seguridad), `LOGIN-01`..`04` / `INT-01`..`09` (Hitos de integración) y `FIX-01`..`119` (Correcciones).
- **Asignación nominal por especialidad:**
  - **Backend (`centinela-back`):** `Tayra` y `Lisandro` (con soporte puntual de `Lucas` y `Nico`).
  - **Frontend (`centinela-front`):** `Belinda`, `Luz` y `Cristian` (con soporte de `Nico` en cliente de eventos SSE).
  - **Infraestructura y Servidor Proxmox VE:** `Nico` y `Lucas`.
  - **Dirección de Arquitectura, Desglose en ClickUp, Auditorías de Seguridad y Suites `/test`:** `Micaela`.
- **Estimación horaria acotada (`0,5 h` a `2,5 h`):** Dimensionada para completarse y verificarse dentro de la misma jornada o sprint corto, permitiendo medir desvíos reales de tiempo en ClickUp.
- **Separación entre "Dependencia para empezar" y "Dependencia para cerrar":** Por ejemplo, en Etapa 1 y Etapa 2 los desarrolladores de Frontend nunca tuvieron que esperar a que Backend terminara un endpoint para empezar a programar su vista: arrancaban contra el documento de contrato cerrado (`BAC-29`) y los mocks tipados, y solo dependían del endpoint real para el cierre de integración.
- **Criterio de éxito ejecutable:** Toda tarea especifica qué test de `test/back` o `test/front` certifica su cumplimiento.

---

## 3. El Motor de Verificación Independiente (`/test`)

A partir del **15/09/2026** (commit `6a9b2ad`), se introdujo en la raíz del repositorio el directorio [`test/`](../test) con una premisa innegociable: **las pruebas de aceptación no viven dentro de los submódulos de los desarrolladores, sino en el repositorio orquestador**.

### 3.1. Principios del Arnés de Pruebas
1. **Cero Falsos Positivos:** Ningún test aprueba con `t.Logf` ni condiciones laxas. Cada prueba verifica el código de estado HTTP, el `errorCode` del catálogo estandarizado, la mutación real en PostgreSQL/Redis y el cambio visible en el DOM.
2. **Circuito Real de Credenciales (*Zero-Shortcut*):** Las pruebas de backend no inyectan JWTs falsos ni saltean el 2FA. Extraen las contraseñas temporales y códigos OTP desde los logs del contenedor o desde la API de **Mailpit (`:18025` con STARTTLS obligatorio)** y completan el flujo real `POST /auth/login` $\rightarrow$ `GET /auth/2fa/qr` $\rightarrow$ `POST /auth/2fa/verify`.
3. **Aislamiento Total del Servidor Físico (`dns_search: invalid`):** En `test/back/compose.yaml`, todos los contenedores bloquean la resolución DNS de sufijos locales para impedir que cualquier llamada caiga por accidente en el servidor Proxmox real de la red Tailscale (`100.81.49.19`).
4. **Validación de la Propia Suite mediante *Mutation Testing*:** Como consta en `test/front/RESULTADOS.md` (05/10/2026), la suite de frontend fue sometida a contrapruebas inyectando 13 defectos deliberados uno por uno sobre una implementación de referencia, comprobando que cada defecto hiciera fallar exactamente el test diseñado para detectarlo.

### 3.2. Topología de `test/back` (Go + 9 Contenedores en Docker Compose)
Al ejecutar `go test -v -count=1 ./...` dentro de `test/back`, `main_test.go` levanta automáticamente un entorno efímero de **9 contenedores**:
- **`db` (`postgres:16-alpine`):** Ejecuta `backend/scripts/init.sql` en limpio para validar migraciones, triggers de inmutabilidad (`trg_auditoria_inmutable`), particionado trimestral y claves foráneas.
- **`redis` (`redis:7-alpine`):** Almacena sesiones (`sesion:<sid>`), tickets efímeros de un solo uso para SSE/WebSocket (`GETDEL`), caché de nodo/RRD, locks distribuidos (`lock:instance:provisioning`) y el bus Pub/Sub.
- **`proxmox` (Stub Nginx con interruptores en vivo):** Simula la API de Proxmox VE y permite activar caídas o fallos en caliente creando archivos dentro del contenedor (`/tmp/proxmox-caido` para forzar `502`, `/tmp/tareas-en-curso`, `/tmp/tareas-fallidas`).
- **`proxmox-sim` (Simulador dinámico en Go):** Compila `backend/cmd/proxmox-simulador` para probar aprovisionamiento de VMs/LXC, consumo de cuotas, `/cluster/nextid`, series RRDdata y ciclo completo de snapshots y rollback.
- **`proxmox-lento` + `backend-lento` (`:18082`):** Un contenedor `busybox` que acepta conexiones TCP sin responder jamás, usado exclusivamente para certificar que el backend corta a tiempo y devuelve `504 PROXMOX_TIMEOUT` (`FIX-40`).
- **`mailpit` + `backend-smtp` (`:18081`):** Certifica el envío real de correos por SMTP con verificación estricta de certificados TLS (`MP_SMTP_REQUIRE_STARTTLS=1`).
- **`backend` (`:18080`) y `edge` (`nginx:1.29-alpine` en `:18089` HTTP / `:18453` HTTPS):** Certifican toda la API REST, el streaming SSE sin buffering (`proxy_buffering off`), el upgrade de WebSockets y las cabeceras TLS/HSTS/CORS de borde.

### 3.3. Suites E2E Cruzadas y Suites Adversariales
Además de las pruebas por etapa (`backend_acceptance_test.go`, `etapa1_acceptance_test.go`, `etapa1_ola2_acceptance_test.go`, `etapa1_ola3_acceptance_test.go`, `etapa2_acceptance_test.go`, `fixes_acceptance_test.go`), el repositorio ejecuta dos categorías especiales de validación:
1. **Pruebas E2E reales Front $\leftrightarrow$ Back (`login04_integracion_front_back_test.go` e `int02_integracion_front_back_test.go`):** Desde el test de Go se levanta todo el stack Docker del backend y se invoca Vitest (`login04-e2e.test.ts` e `int02-e2e.test.ts`) haciendo que los **servicios reales de TypeScript del submódulo `frontend`** disparen peticiones HTTP y abran conexiones SSE reales contra el **binario Go del submódulo `backend`**.
2. **Suites de Integridad y Seguridad Adversarial (`adversarial_integrity_test.go` y `adversarial-security.test.ts`):** Incorporadas el 13/10/2026 (`12e7b46`), someten a ambos submódulos a ataques de escalada horizontal/vertical de privilegios RBAC, intentos de operar o apagar la VM protegida `103`, ráfagas concurrentes para forzar condiciones de carrera (`409 INSTANCE_BUSY`), adulteración de firmas JWT, reutilización de tickets efímeros SSE/WS, inyección de payloads XSS en mensajes de error, Path Traversal y Prototype Pollution.

---

## 4. Línea de Tiempo de Progreso y Cortes de Verificación

A lo largo de las 6 semanas documentadas en el historial de Git, el proyecto atravesó los siguientes hitos de planificación, auditoría y remediación:

### 4.1. Arranque y Construcción del Arnés de Control (10/09/2026 – 18/09/2026)
- **10/09 (`93d5061`):** Se estructura el repositorio con los requerimientos iniciales y el esquema de archivos `actual.md`, `futuro.md` y `terminado.md`. Se definen las primeras tareas de Fase Base (`BAC-01..04`, `FRN-01..04`, `LOGIN-01..03`).
- **14/09 (`83f7ed0`, `07761c2`):** Se enlazan los submódulos `frontend` y `backend` en `.gitmodules` y se construye el entorno Docker integrado (`compose.yaml`, `docker/`, `scripts/centinela.sh` y workflow de CI en `.github/workflows/integration-tests.yml`).
- **15/09 (`6a9b2ad`, `e0051b8`):** Se inauguran las suites independientes `test/back` y `test/front` junto con `test/informe.md`.
- **16/09 – 18/09 (`d544496`):** Primer contraste empírico contra tareas declaradas completas en ClickUp. Las suites detectan discrepancias de contrato entre lo que enviaba el frontend y lo que esperaba el backend en login, 2FA y permisos por instancia; nacen los primeros **`FIX-01` a `FIX-16`**.

### 4.2. Consolidación de Fase Base y Desacople *Local-First* (20/09/2026 – 01/10/2026)
- **20/09 – 28/09 (`833f161`, `07f60e9`, `55732c8`):** Se documenta `api-proxmox.md` definiendo el blindaje de los VMIDs de infraestructura (`100..105`). En los cortes del **23/09**, **25/09** y **28/09** se verifican y promueven `SEC-01`, `SEC-02`, `SEC-04`, `FRN-05`, `FRN-11`, `FRN-12`, `INF-05` y `FIX-17` a `FIX-25`.
- **29/09 – 30/09 (`5df8274` — *"rompo dependencia de back con estado de tareas de infra en servidor"*):** Hito clave de gestión. Para evitar que el atraso en tareas del servidor físico frenara al equipo de Go y React, se introduce el simulador de Proxmox, Redis local y Mailpit en Docker. Se verifican `FIX-22`, `FIX-23` (auditoría inmutable en PostgreSQL), `BAC-17A/B` (1 sesión = 1 fila tras múltiples refresh) y `BAC-21C` (SSE con tickets efímeros en Redis). Se publica la planificación completa de `etapa1.md`.
- **01/10 (`97b8484`..`826530c`):** **Cierre formal de la Fase Base.** El hito **`LOGIN-04` pasa 10/10 pruebas E2E** integrando el frontend real contra el backend real. Se cierran las **Decisiones de Diseño `D1` a `D6`** en `etapa1.md` (contratos de telemetría, matriz de errores HTTP `D2`, umbrales del semáforo `D3`, timeout UPID de 3 min `D4`, rollback ante fallo SMTP `D5` y reglas de OTP `D6`).

### 4.3. Ejecución de Etapa 1, Rediseño de Baja Lógica y Puente a Etapa 2 (02/10/2026 – 11/10/2026)
- **02/10 – 03/10 (`cb0388e`, `b8e66c0`):** Se verifica `FIX-40` (`504 PROXMOX_TIMEOUT` mediante el contenedor `proxmox-lento`), `SEC-03` (`PermissionGate`), `INF-07B` (SSE sin buffering en Nginx Edge) y `FRN-19A` (medidores del nodo).
- **04/10 – 06/10 (`1e73ce1`, `cb90b1b`):** Al probar la eliminación de usuarios con historial en la tabla inmutable `auditoria`, se detecta la colisión con la restricción `ON DELETE RESTRICT` y con la reutilización de correos electrónicos. Se diseñan y cargan **`FIX-54` a `FIX-57`**, separando la **suspensión reversible** (`activo = false`) de la **baja lógica irreversible** (`eliminado_en = NOW()` con índice único parcial `WHERE eliminado_en IS NULL`). Se promueven `BAC-22`, `BAC-22B`, `BAC-23A/B`, `BAC-24A/B`, `BAC-25A/C`, `FRN-16B`, `FRN-17C` y se publica la planificación de **`etapa2.md`** junto con las tareas puente (`BRG-*`).
- **08/10 – 11/10 (`21567ab`, `bbe3313`):** Corte de verificación del **09/10** promoviendo `FRN-15`, `FRN-16`, `FRN-19B`, `FRN-20B`, `BAC-25B` y `FIX-58` a `FIX-70`. Se incorporan las suites de aceptación para Ola 2 y Ola 3 de Etapa 1 y los hitos integrales `INT-01`, `INT-02` e `INT-03`.

### 4.4. Auditoría de Stubs, Avance de Etapa 2, Diseño de Etapa 3 y Auditorías de Seguridad (13/10/2026 – 19/10/2026)
- **13/10 (`c78de57`, `12e7b46`):**
  - Se realiza un barrido estático de código sobre submódulos buscando **stubs, mocks olvidados y placeholders** en tareas que habían sido declaradas terminadas: se detectan datos hardcodeados en la vista de detalle de usuario (`FIX-87`), formulario maquetado sin conexión en `/instances/new` (`FIX-89`), tarjetas muertas en Dashboard (`FIX-90`) y el componente muerto `InstanceTester.tsx` con datos falsos (`FIX-91`).
  - Se incorporan las suites de seguridad adversarial (`adversarial_integrity_test.go` y `adversarial-security.test.ts`).
  - Se promueven a `terminado-1.md` `FIX-54..56`, `BAC-27`, `BAC-30`, `BAC-31`, `INF-07C` e `INF-09`.
- **14/10 – 16/10 (`6477aef`, `12e7b46`, `43adebe`, `2e11d46`):**
  - Se desglosa íntegramente la **`etapa3.md`**.
  - En el corte del **15/10** (contra backend `dbf7075` y frontend `6db04d3`), se inaugura `terminado2.md` promoviendo las primeras tareas de Etapa 2 verificadas en verde: `BAC-34` (adaptador RRDdata), `BAC-35` (WebSocket de métricas con ticket efímero), `BAC-37` (validación de cuotas con 15% de reserva y lock Redis `lock:instance:provisioning`) y `BAC-38` (`POST /api/instances`).
  - Se redacta [`informe-evolucion-fasebase-etapa1-etapa2.md`](./archvios/informe-evolucion-fasebase-etapa1-etapa2.md) analizando la asimetría de velocidad entre el equipo de Backend (que ya avanzaba sobre Etapa 2) y el de Frontend (que aún cerraba pendientes de Etapa 1).
- **19/10 (`adcb642`, `ea3828b` — Corte actual contra backend `8daaf87` y frontend `91eabf4`):**
  - **Promociones en Backend (Etapa 1 y Etapa 2):** Pasan al 100% en verde y se archivan en `terminado-1.md` y `terminado2.md` las tareas `FIX-81` (orden de migraciones en base limpia), `BAC-39` (seguimiento de creación y autoasignación de permisos), `BAC-40` (`PUT /api/instances/:vmid/config` con $\Delta\text{RAM}$ y `requiresReboot`), `BAC-41` (auditoría de creación y config), `BAC-50` (simulador de snapshots y rollback), `BAC-51` (guarda `RequireInstanceIdle` $\rightarrow$ `409 INSTANCE_BUSY`), `BAC-52` (broker SSE para alertas `HOST` solo a `ADMIN`) y `BAC-53` (particionado preventivo `2026_q3`..`2027_q1` e índice compuesto).
  - **Promociones y Regresiones Detectadas en Frontend:** Se promueven en verde `FRN-21`/`FIX-95` (detalle de instancia `/instances/:vmid`), `FRN-17B` (reacción a `TASK_FINISHED`), `FIX-86` y `FIX-97`. Se detecta que el commit `54efc3d` introdujo una regresión en `useNodeStatus.ts:146` (`TypeError: serie.map is not a function`, documentada en `FIX-103`) y que el ocultamiento del `UPID` en el backend (`bd22c3d`) rompió la validación rígida de `response.upid` en `instanceService.ts` (`FIX-104B` / `FIX-105`).
  - **Auditorías Profundas de Ciberseguridad e Infraestructura:** Se ejecutan dos auditorías exhaustivas de arquitectura y código fuente ([`analisis-cifrado-ids-vulnerabilidades-infraestructura.md`](./archvios/analisis-cifrado-ids-vulnerabilidades-infraestructura.md) y [`analisis-seguridad-vulnerabilidades-permisos.md`](./archvios/analisis-seguridad-vulnerabilidades-permisos.md)), de las cuales se desprenden y planifican las tareas atómicas **`FIX-105` a `FIX-119`** en [`actual.md`](./actual.md) y [`futuro2.md`](./futuro2.md) junto con las tareas de Frontend de Etapa 2 (`FRN-22` a `FRN-25`).

---

## 5. Catálogo de Desfasajes y Vulnerabilidades Detectadas por el Control Centralizado

La tabla siguiente resume los casos más representativos donde la auditoría cruzada de este repositorio evitó fallas de integración o brechas graves de seguridad:

| Categoría | Caso Detectado | Impacto si no se auditaba | Tarea / Fix Generado |
| :--- | :--- | :--- | :--- |
| **Ruptura de Contrato Front $\leftrightarrow$ Back** | El backend (`bd22c3d`) ocultó el `UPID` de Proxmox devolviendo `{ tareaId }`, pero `instanceService.ts` en el frontend lanzaba `ApiRequestError` si `response.upid` no era un `string`. | Todas las órdenes de encender, apagar, reiniciar, eliminar y crear instancias fallaban en la interfaz web pese a ejecutarse en el backend. | **`FIX-104B` / `FIX-105` y `FIX-106`** |
| **Fuga de Metadatos del Hipervisor (CWE-200 / CWE-639)** | Exposición de `VMID` secuenciales (`100`, `101`...), nombres internos de nodos (`pve-01`), pools de storage (`local-lvm`) y rutas Perl en errores (`/usr/share/perl5/PVE/...`). | Permite enumerar máquinas (IDOR), mapear la topología interna del datacenter y dirigir ataques al sistema de archivos del host. | **`FIX-107` a `FIX-110`** (Hashids con `INSTANCE_ID_SALT`, aliases de nodo/storage y saneamiento de errores) |
| **Secuestro de Cuenta (*Account Takeover* — CWE-306)** | `PUT /api/account/profile` permitía cambiar el `emailUsuario` sin pedir contraseña actual, sin TOTP, sin confirmar el nuevo correo y sin registrar en `auditoria`. | Un atacante con acceso momentáneo a una sesión cambiaba el email por el suyo, pedía recuperación de clave por OTP y se adueñaba de la cuenta (incluso de un `ADMIN`). | **`FIX-111A` (Back) y `FIX-111B` (Front)** |
| **Fuerza Bruta y Condición de Carrera en OTP (CWE-307)** | Ausencia de *Rate Limiting* en `/auth/login`, `/auth/2fa/verify` y `/auth/password/forgot`, y lectura sin lock de `intentos_fallidos` en `/password/reset`. | Ataques de diccionario, fuerza bruta paralela sobre los 6 dígitos del TOTP en la ventana de 5 min y *Mail Bombing* al servidor SMTP. | **`FIX-112`** (Rate Limiter en Redis + `SELECT ... FOR UPDATE` en PostgreSQL) |
| **Permisos Huérfanos por Reciclaje de VMID (CWE-284)** | Al eliminar una instancia (`DELETE /api/instances/:vmid`), no se borraban sus filas en `permisos_instancia`. Al crear otra máquina, `/cluster/nextid` reutilizaba el `VMID`. | Los operadores que tenían permiso sobre la máquina vieja heredaban acceso automático no autorizado sobre la máquina nueva. | **`BAC-24C` / `FIX-113`** (Purga transaccional + auditoría `INSTANCE_PERMISSIONS_PURGED`) |
| **Persistencia de Sesiones tras Cambio de Clave (CWE-613)** | `PUT /api/account/password` actualizaba el hash bcrypt pero no revocaba las sesiones activas en `sesiones_activas` ni en Redis. | Si un usuario cambiaba su clave por sospecha de intrusión, el atacante conectado en otro navegador seguía operando con su token y refresh cookie. | **`FIX-114`** (Revocación de sesiones + corte SSE/WS) |
| **Restricción Anticuada de Contraseñas (CWE-521)** | Validadores en Go (`password.go`) y React rechazaban contraseñas de más de 12 caracteres (`max=12`). | Impedía usar gestores de contraseñas y *passphrases*, incumpliendo NIST SP 800-63B. | **`FIX-115A` y `FIX-115B`** (Ampliación a `max=128`) |
| **Inyección de Parámetros hacia Proxmox (CWE-88)** | `client.go` interpolaba `input.Bridge` directamente en `name=eth0,bridge=%s,firewall=%s` sin filtrar comas (`,`). | Enviar `vmbr0,rate=50,tag=100` inyectaba directivas de VLAN y límite de ancho de banda no autorizadas en Proxmox VE. | **`FIX-116`** (Sanitización estricta por expresión regular) |
| **Exposición de Telemetría Global y Asignación de VMs Protegidas** | `/api/node/status` entregaba a los `OPERATOR` el conteo global del clúster y discos del host, y `/permissions` permitía asignarles los VMIDs protegidos `100..105`. | Fuga de actividad global del servidor y monitoreo indebido de las máquinas internas donde corre El Centinela. | **`FIX-117A/B` y `FIX-119`** |
| **Colisión entre Baja de Usuario e Inmutabilidad de Auditoría** | Eliminar físicamente un usuario chocaba con la FK de la tabla inmutable `auditoria`, pero solo desactivarlo impedía reutilizar su email. | Imposibilidad de dar de baja definitiva o reciclar correos sin violar la integridad forense del registro de auditoría. | **`FIX-54`, `FIX-55` y `FIX-56`** (Separación de suspensión vs baja lógica con `eliminado_en` e índice único parcial) |

---

## 6. Guía Rápida del Flujo Operativo en cada Corte de Verificación

Cada vez que se realiza un corte de control sobre el avance del equipo en ClickUp, se ejecuta el siguiente protocolo estandarizado en este repositorio:

1. **Sincronización de punteros de submódulos:**
   ```bash
   git -C backend fetch origin --prune && git -C backend reset --hard origin/main
   git -C frontend fetch origin --prune && git -C frontend reset --hard origin/main
   ```
2. **Inspección de commits nuevos y auditoría estática:**
   Se revisan los commits ingresados desde el corte anterior (`git -C backend log <sha_anterior>..HEAD` y `git -C frontend log <sha_anterior>..HEAD`), verificando que no introduzcan stubs, credenciales fijas ni rupturas de contrato.
3. **Ejecución de las suites independientes de aceptación y seguridad:**
   ```bash
   cd test/back && go test -v -count=1 ./...
   pnpm --dir test/front test
   ```
4. **Actualización del dictamen técnico y de los tableros:**
   - Se vuelcan los resultados, métricas y causas raíz en [`test/informe.md`](../test/informe.md).
   - Las tareas con 100% de sus pruebas en verde se mueven de [`actual.md`](./actual.md) a [`terminado.md`](./terminado.md) / [`terminado-1.md`](./terminado-1.md) / [`terminado2.md`](./terminado2.md) dejando constancia de la fecha y el `SHA` certificado.
   - Si una tarea presenta un defecto parcial o regresión, se documenta el diagnóstico exacto (archivo, línea y comportamiento esperado) en un nuevo `FIX-XX` dentro de [`actual.md`](./actual.md) y [`futuro-1.md`](./futuro-1.md) / [`futuro2.md`](./futuro2.md) para su asignación inmediata en ClickUp.
5. **Commit en el repositorio orquestador:**
   Se versionan conjuntamente los punteros actualizados de `backend` y `frontend`, las nuevas pruebas de regresión en `test/` y la actualización de `documentacion/`, manteniendo la trazabilidad histórica completa del proyecto.

