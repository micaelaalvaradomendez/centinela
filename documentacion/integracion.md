# Tareas de Integración y Validación End-to-End (Todas las Etapas)

Este documento consolida y centraliza **todas las tareas de integración, suites automatizadas y pruebas de aceptación (Smoke Tests / UAT)** del proyecto **El Centinela**, cubriendo desde la Etapa 1 hasta la etapa de Adaptación a Proxmox VE 8.4.21 (e incluyendo las tareas de validación archivadas de la Etapa 4).

---

## 🗺️ Mapa General de Integración

| Hito | Tarea | Área | Responsables | Estimación | Ámbito | Dependencia Principal |
|---|---|---|---|---|---|---|
| **Etapa 1** | `INT-01` | Backend / Testing | Tayra y Cristian | 2.5 h | Local (CI/Docker) | `BAC-23B` a `BAC-27` |
| **Etapa 1** | `INT-02` | Integración / E2E | Todo el equipo | 2.0 h | Local (Stack MVP) | Todas las tareas Etapa 1 |
| **Etapa 1** | `INT-03` | Integración / Infra | Nico y Backend | 1.5 h | Servidor / Edge TLS | `INT-02`, `INF-07`, `INF-07B` |
| **Etapa 2** | `INT-04` | Backend / Testing | Tayra y Lisandro | 2.5 h | Local (CI/Docker) | `BAC-36`, `BAC-39`, `BAC-41` |
| **Etapa 2** | `INT-05` | Integración / E2E | Todo el equipo | 2.0 h | Local (Stack Completo) | `INT-04`, `FRN-23`, `FRN-26`, `FRN-28` |
| **Etapa 2** | `INT-06` | Integración / Infra | Nico y Backend | 2.0 h | Proxmox Real | `INT-05`, `INF-10` |
| **Etapa 3** | `INT-07` | Backend / Testing | Tayra y Lisandro | 2.5 h | Local (CI/Docker) | `BAC-43`, `BAC-45`, `BAC-47`, `BAC-48` |
| **Etapa 3** | `INT-08` | Integración / E2E | Todo el equipo | 2.0 h | Local (UAT Pre-Defensa)| `INT-07`, `FRN-30`, `FRN-32`–`FRN-34` |
| **Etapa 3** | `INT-09` | Integración / Infra | Nico, Lucas y Back | 2.0 h | Proxmox Real Prod | `INT-08`, `INF-11`, `INF-13` |
| **Puente PVE8** | `INT-10` (`BRG-13-INT`) | Testing / QA | Tayra | 1.5 h | Test Harness Dual | `INT-07` |
| **Adaptación PVE8** | `ADP-11` (`INT-PVE8-01`) | Testing / QA | Tayra / Lisandro | 2.0 h | Suites Duales 8/9 | `ADP-03`, `ADP-05` |
| **Adaptación PVE8** | `ADP-12` (`INT-PVE8-02`) | Testing / Backend | Lisandro | 2.5 h | Simulador PVE 8.4 | `ADP-01` a `ADP-11` |
| **Adaptación PVE8** | `ADP-13` (`INT-PVE8-03`) | Integración Global | Equipo Completo | 2.5 h | Servidor PVE 8.4 Real | `ADP-10`, `ADP-12` |
| *Etapa 4 (PaaS)* | *`INT-10` (PaaS)* | Integración | Tayra / Lisandro | 2.5 h | Local (Ingress/Build) | `BAC-63`, `INF-15` |
| *Etapa 4 (PaaS)* | *`INT-11` (PaaS)* | Integración E2E | Todo el equipo | 2.0 h | Dominio Público | `INT-10`, `FRN-42` |

---

## 📌 ETAPA 1: Visibilidad y Control Operativo (Hito MVP)

### `INT-01` - Pruebas de integración automatizadas de ciclo de vida y worker UPID

- **Área:** Backend / Testing
- **Asignados:** Tayra y Cristian
- **Estimación:** 2.5 h
- **Depende de:** `BAC-23B`, `BAC-24A`, `BAC-24B`, `BAC-24C`, `BAC-25A`, `BAC-25B`, `BAC-25C`, `BAC-26` y `BAC-27`.
- **Validación en el repositorio:**
  - **Archivo de prueba:** [`test/back/int01_acceptance_test.go`](file:///home/micaela/centinela/test/back/int01_acceptance_test.go) (`TestINT01_CicloDeVidaYWorkerUPID`).
  - **Comando:** `go test -v -run TestINT01_CicloDeVidaYWorkerUPID .` en `test/back`.
- **Entregables y cobertura automatizada:**
  1. **Despacho asíncrono y rechazos canónicos:** cada acción operativa devuelve HTTP `202 Accepted` con `upid` y `tareaId`. Órdenes inválidas devuelven `409 INSTANCE_INVALID_STATE`, VMIDs protegidos devuelven `403 INSTANCE_PROTECTED`, falta de permisos devuelve `403 INSTANCE_ACCESS_DENIED`, y colisiones concurrentes devuelven `409 INSTANCE_BUSY`.
  2. **Ciclo del worker UPID:** sondeo del worker hasta `COMPLETED` (`exitstatus: OK`) y vencimiento por `UPID_TIMEOUT` (3 minutos) o fallos `PROXMOX_ERROR`.
  3. **Eventos en tiempo real:** emisión de `TASK_FINISHED` por `/api/events` con el payload de `BAC-29` y filtrado estricto por permiso de instancia (un operador sin permiso no lo recibe; uno con `READ_ONLY` o `FULL_ACCESS` sí).
  4. **Auditoría dual:** dos registros correlacionados en `auditoria` por cada acción (`PENDING` al despachar la orden y `SUCCESS`/`FAILED` al completarse).
  5. **Resiliencia ante caídas:** reanudación de UPIDs en curso al reiniciar el proceso backend (`RNF-04`).
  6. **Limpieza de permisos:** al completarse un `DELETE`, eliminación atómica de las tuplas del VMID en `permisos_instancia` (`BAC-24C`).
- **Criterio de éxito:** La suite pasa al 100% en verde sin pasos manuales ni mocks desconectados.

---

### `INT-02` - Validación integral del Hito MVP Operativo en local (Smoke Test End-to-End)

- **Área:** Integración / Todo el equipo
- **Asignados:** Equipo completo
- **Estimación:** 2.0 h
- **Depende de:**
  - Todas las tareas de desarrollo de la Etapa 1.
  - Tareas base de seguridad: `FIX-37`, `SEC-03`, `FIX-38` y `FRN-17C`.
- **Validación en el repositorio:**
  - **Suite backend:** [`test/back/int02_acceptance_test.go`](file:///home/micaela/centinela/test/back/int02_acceptance_test.go) (`TestINT02_SmokeTestEndToEndMVP`).
  - **Suite frontend E2E:** [`test/front/int02-e2e.test.ts`](file:///home/micaela/centinela/test/front/int02-e2e.test.ts).
  - **Orquestación integrada:** [`test/back/int02_integracion_front_back_test.go`](file:///home/micaela/centinela/test/back/int02_integracion_front_back_test.go) (`TestINT02IntegracionFrontBack`).
  - **Comandos:**
    - `go test -v -run TestINT02_SmokeTestEndToEndMVP .` en `test/back`
    - `npx vitest run int02-e2e` en `test/front`
- **Entregable y cobertura automatizada (7 pasos del MVP):**
  1. **Login y 2FA:** un `OPERATOR` con `FULL_ACCESS` asignado sobre la VM 101 completa login con contraseña y 2FA TOTP contra el backend real.
  2. **Dashboard:** consulta `GET /api/node/status` validando telemetría del host (CPU, RAM, Storage) y resumen de instancias (`instancesSummary`) sin recarga de página.
  3. **Inventario:** consulta `GET /api/instances` validando filtrado RBAC (solo visualiza la instancia 101 asignada con `nivelAcceso: FULL_ACCESS`).
  4. **Ciclo de vida:** envía orden de apagado ordenado (`POST /instances/101/stop` o `/status/shutdown`), recibe `202 Accepted` con `{ upid, tareaId }` y valida estado transitorio.
  5. **Notificación:** solicita ticket efímero (`POST /events/ticket`), se conecta al stream SSE y recibe `TASK_FINISHED` liberando el bloqueo de la fila.
  6. **Auditoría:** el administrador consulta `GET /api/admin/audit` y constata los registros auditados de la acción de ciclo de vida.
  7. **Borrado y revocación:** el administrador elimina una instancia (`DELETE /instances/:vmid`) y el operador pierde inmediatamente el acceso en su inventario.
- **Criterio de éxito:** Los 7 pasos pasan al 100% en verde ejecutando el código real de servicios del frontend y el stack real del backend.

---

### `INT-03` - Validación del Hito MVP en el servidor contra el Proxmox real / Entorno Edge TLS

- **Área:** Integración / Infraestructura
- **Asignados:** Nico y un integrante de backend
- **Estimación:** 1.5 h
- **Depende de:** `INT-02`, `INF-07`, `INF-07B`, `INF-06B` e `INF-08A`.
- **Validación en el repositorio:**
  - **Archivo de prueba:** [`test/back/int03_acceptance_test.go`](file:///home/micaela/centinela/test/back/int03_acceptance_test.go) (`TestINT03_ValidacionServidorBordeTLS`).
  - **Comando:** `go test -v -run TestINT03_ValidacionServidorBordeTLS .` en `test/back`.
- **Entregables y cobertura automatizada:**
  1. **Terminación TLS y certificados:** verificación de canal HTTPS seguro en el puerto `18453` a través de Nginx Edge con validación de CA (`edge-tls/ca.crt`).
  2. **Streaming SSE sin cortes:** verificación de que `/api/events` no sufra buffering en Nginx (`proxy_buffering off`) y fluya en tiempo real.
  3. **Protección estricta de infraestructura:** intento de mutación o borrado sobre instancias protegidas (`100` a `105`) rechazado con `403 INSTANCE_PROTECTED`.
  4. **Ciclo de vida sobre instancia permitida:** ejecución de acciones operativas sobre instancias dinámicas ($\text{VMID} \ge 106$) a través del proxy.
  5. **Tolerancia a fallos de borde:** diferenciación exacta entre indisponibilidad del hipervisor (`502 PROXMOX_UNAVAILABLE`) y agotamiento de tiempo de respuesta (`504 PROXMOX_TIMEOUT`).
- **Criterio de éxito:** La suite `TestINT03_ValidacionServidorBordeTLS` pasa en verde validando que la capa de red y reverse proxy no degrade las garantías operativas ni de seguridad del sistema.

---

## 📌 ETAPA 2: Aprovisionamiento, Telemetría WS y Edición de Recursos

### `INT-04` - Suite automatizada de pruebas de integración de la Etapa 2

- **Área:** Backend / Testing
- **Asignados:** Tayra y Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `BAC-36`, `BAC-39`, `BAC-41`.
- **Problema:** Se requiere una suite de aceptación automatizada que valide en CI/local todo el flujo de cuotas, creación, telemetría y edición contra el simulador.
- **Entregables:**
  - Archivo de prueba [`test/back/etapa2_acceptance_test.go`](file:///home/micaela/centinela/test/back/etapa2_acceptance_test.go):
    1. Consumo de métricas históricas RRD y conexión WebSocket con ticket efímero.
    2. Creación exitosa de instancia con autoasignación de permisos en `permisos_instancia`.
    3. Rechazo estricto por cuota insuficiente (`409 RESOURCE_QUOTA_EXCEEDED`) y prueba de concurrencia con candado Redis.
    4. Modificación de configuración con cálculo diferencial y verificación de auditoría.
- **Criterio de éxito:** `go test ./test/back/etapa2_acceptance_test.go` pasa 100% en verde sin dependencias externas.

---

### `INT-05` - Smoke Test End-to-End local de la Etapa 2

- **Área:** Integración / Todo el equipo
- **Asignados:** Equipo completo
- **Estimación:** 2.0 h
- **Depende de:** `INT-04`, `FRN-23`, `FRN-26`, `FRN-28`.
- **Problema:** Verificación manual e interactiva integral de experiencia de usuario y coherencia funcional antes de pasar al servidor físico.
- **Entregables (Recorrido interactivo en entorno local):**
  1. Admin accede a `/instances`, pulsa "Crear instancia", completa el Wizard de 5 pasos para un contenedor LXC y confirma.
  2. Ve la fila en estado "Creando...", recibe notificación de tarea finalizada y la máquina aparece lista.
  3. Abre `/instances/:instanceRef`, verifica los gráficos de métricas actualizándose en vivo por WebSocket y prueba el selector de 2h/24h.
  4. Abre el modal de edición, amplía la memoria RAM, comprueba el cálculo diferencial y confirma.
  5. Verifica en `/admin/audit` que la creación y la edición quedaron registradas inmutablemente.
- **Criterio de éxito:** Flujo completo ejecutado sin errores en la consola del navegador ni inconsistencias en la base de datos.

---

### `INT-06` - Validación en servidor de la Etapa 2 contra Proxmox real

- **Área:** Integración / Infraestructura
- **Asignados:** Nico y un integrante de backend
- **Estimación:** 2.0 h
- **Depende de:** `INT-05`, `INF-10`.
- **Problema:** Validación definitiva sobre el hipervisor real y a través del proxy inverso de producción antes de cerrar la etapa.
- **Entregables:**
  - Ejecutar los pasos del Smoke Test (`INT-05`) en el entorno de pruebas contra el Proxmox real:
    1. Crear un contenedor LXC real con la plantilla de prueba en un VMID libre ($\ge 106$).
    2. Comprobar que Nginx enruta el WebSocket de métricas sin cortes.
    3. Modificar la memoria del contenedor.
    4. Destruir el contenedor de prueba con `DELETE /api/instances/:vmid` (`BAC-24B`/`BAC-24C`) para dejar el entorno limpio.
- **Criterio de éxito:** El aprovisionamiento, telemetría y edición funcionan idénticamente en el servidor físico, registrando auditoría y respetando los VMIDs protegidos.

---

## 📌 ETAPA 3: Snapshots, Exportación, Alertas y Hardening

### `INT-07` - Suite automatizada de pruebas de integración de la Etapa 3

- **Área:** Backend / Testing
- **Asignados:** Tayra y Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `BAC-43`, `BAC-45`, `BAC-47`, `BAC-48`.
- **Problema:** Se requiere una suite de aceptación automatizada integral en CI/local que verifique todo el circuito de snapshots, auditoría extendida, alertas y seguridad.
- **Entregables:**
  - Archivo de prueba [`test/back/etapa3_acceptance_test.go`](file:///home/micaela/centinela/test/back/etapa3_acceptance_test.go):
    1. Creación, listado, rollback y borrado de snapshots con validación de lock `409 INSTANCE_BUSY`.
    2. Filtrado y streaming de exportación de auditoría en CSV y JSON.
    3. Simulación de saturación de nodo y verificación de emisión de alertas en tiempo real.
    4. Comprobación de cookies `HttpOnly` y rechazo por rate limiting (429).
- **Criterio de éxito:** `go test ./test/back/etapa3_acceptance_test.go` pasa 100% en verde de forma determinística contra el simulador.

---

### `INT-08` - Smoke Test End-to-End local del sistema completo (UAT pre-defensa)

- **Área:** Integración / Todo el equipo
- **Asignados:** Equipo completo
- **Estimación:** 2.0 h
- **Depende de:** `INT-07`, `FRN-30`, `FRN-32`, `FRN-33`, `FRN-34`.
- **Problema:** Ensayo general de la plataforma completa previo a la validación final en hardware físico.
- **Entregables (Recorrido funcional interactivo en entorno local):**
  1. Login de operador con 2FA y cookie segura `centinela_refresh`.
  2. Monitoreo de telemetría y visualización de gráficos de recursos en vivo.
  3. Creación de snapshot sobre una máquina y posterior restauración con confirmación antierror tipear nombre idéntico.
  4. Generación deliberada de carga para disparar banner de alerta en tiempo real.
  5. Consulta de auditoría administrativa y descarga del reporte CSV y JSON.
- **Criterio de éxito:** Los 5 flujos completan con éxito sin errores en la consola de navegación ni anomalías de base de datos.

---

### `INT-09` - Validación final en servidor de producción con Proxmox real

- **Área:** Integración / Infraestructura
- **Asignados:** Nico, Lucas y un integrante de backend
- **Estimación:** 2.0 h
- **Depende de:** `INT-08`, `INF-11`, `INF-13` (`BRG-11-INF`: storage compatible con snapshots en Proxmox real).
- **Problema:** Certificación definitiva sobre el hipervisor real a través de Tailscale y proxy inverso con certificados antes del cierre formal del proyecto.
- **Entregables:**
  1. Ejecución de snapshots y rollback sobre una instancia real en Proxmox VE.
  2. Comprobación del stream de alertas y métricas a través de HTTPS/TLS sin interrupciones.
  3. Auditoría exhaustiva del entorno, dejando únicamente las instancias autorizadas y protegidas (`100` a `105`).
  4. Acta técnica de conformidad de la entrega del MVP final.
- **Criterio de éxito:** La plataforma supera todas las validaciones sobre hardware físico cumpliendo el 100% de los requerimientos funcionales (`RF-01` a `RF-11`) y no funcionales.

---

### `INT-10` (`BRG-13-INT`) - Baseline Dual-Test Harness para Suites de Aceptación (Tolerancia PVE 8 y PVE 9)

- **Área:** Testing / QA
- **Asignado:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `INT-07`.
- **Problema:** Pruebas como `TestEtapa1SimuladorProxmox` en `test/back` fallarían si el simulador o un mock responde con tipos enteros en `prefix` o `memory`.
- **Entregables:**
  1. En `test/back/`: crear helper `assertNetworkPrefix(t, val)` que valide prefijos CIDR de red representados tanto como string `"24"` como número `24`.
  2. Actualizar las aserciones de integración existentes para que sean agnósticas de la versión subyacente.
- **Criterio de éxito:** Las pruebas de aceptación corren exitosamente sin fallar ante payloads con prefijo numérico o textual.

---

## 📌 ETAPA ADAPTACIÓN A PROXMOX VE 8.4.21

### `ADP-11` (`INT-PVE8-01`) - Ajuste de Aserciones de Tests de Aceptación para PVE 8 y PVE 9

- **Área:** Testing / QA
- **Asignado:** Tayra / Lisandro
- **Estimación:** 2.0 h
- **Depende de:** `ADP-03`, `ADP-05`.
- **Problema:** En `test/back/simulador_proxmox_acceptance_test.go:42`, la aserción `address["prefix"].(string)` falla si se corre contra PVE 8.4 donde el prefijo sea entero o se normalice numéricamente.
- **Entregables:**
  1. En `test/back/simulador_proxmox_acceptance_test.go`:
     - Flexibilizar la verificación de `prefix`: aceptar `string` ("24") o `float64/int` (24).
     - Asegurar que la prueba valide la IP resultante sin exigir representación rígida del prefijo en el JSON crudo.
  2. En `test/back/etapa3_acceptance_test.go`:
     - Garantizar que las verificaciones de storage admitan `lvmthin`, `zfspool` o `qcow2` según las definiciones de PVE 8.4.
- **Criterio de éxito:** `go test -v ./test/back -run TestEtapa1SimuladorProxmox` pasa en verde tanto con respuestas string como enteras.

---

### `ADP-12` (`INT-PVE8-02`) - Smoke Test Automatizado contra Simulador en Modo PVE 8.4.21

- **Área:** Testing / Backend
- **Asignado:** Lisandro
- **Estimación:** 2.5 h
- **Depende de:** `ADP-01` a `ADP-11`.
- **Problema:** Se requiere certificar que todo el ciclo operativo (inventario, encendido, apagado, aprovisionamiento, edición de recursos y snapshots) funciona íntegramente con el simulador configurado en modo PVE 8.4.21.
- **Entregables:**
  - Nueva suite de prueba [`test/back/pve8_adaptation_acceptance_test.go`](file:///home/micaela/centinela/test/back/pve8_adaptation_acceptance_test.go):
    1. Levanta el simulador con `PROXMOX_SIM_TARGET_VERSION=8.4.21`.
    2. Ejecuta:
       - Creación de VM QEMU con verificación de `agent=1` y `ostype=l26`.
       - Creación de LXC con verificación de `unprivileged=1` y `nesting=1`.
       - Consulta de interfaces y selección de IP.
       - Creación y reversión de snapshot.
       - Eliminación de instancia con parámetros `purge=1`.
- **Criterio de éxito:** La suite `pve8_adaptation_acceptance_test.go` pasa 100% en verde.

---

### `ADP-13` (`INT-PVE8-03`) - UAT y Validación Final sobre Servidor Proxmox VE 8.4.21 Físico

- **Área:** Integración / Toda el área
- **Asignado:** Lisandro, Tayra, Nico, Belinda, Cristian
- **Estimación:** 2.5 h
- **Depende de:** `ADP-10`, `ADP-12`.
- **Problema:** Validación definitiva sobre hardware anfitrión real en la red corporativa para certificar la puesta en producción.
- **Entregables:**
  1. Configuración de variables de entorno productivas en el host de despliegue (`PROXMOX_URL=https://<HOST-PVE8>:8006/api2/json`, token real y `PROXMOX_NODE`).
  2. Ejecución de prueba de ciclo de vida completo:
     - Login administrativo con 2FA.
     - Visualización del dashboard de salud del nodo PVE 8.4.21.
     - Aprovisionamiento de 1 contenedor LXC Debian 12.
     - Aprovisionamiento de 1 VM QEMU Debian 12.
     - Toma de snapshot y rollback.
     - Eliminación limpia de las instancias de prueba.
- **Criterio de éxito:** Informe UAT firmado sin errores 500/502 en los logs de backend ni advertencias de infraestructura en Proxmox.



