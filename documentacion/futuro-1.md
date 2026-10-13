# Etapa 1: correcciones pendientes (FIX)

Correcciones de tareas de la Etapa 1 que ya están implementadas pero no cumplen del todo su criterio de éxito. Las tareas originales están en [`terminado-1.md`](terminado-1.md).

---

## 🛠️ Fixes detectados en la verificación del 03/10/2026

---

## 🛠️ Fixes detectados en la verificación del 05/10/2026

--- 

## Hallazgos incorporados del PR #3 (06/10/2026)

> [!IMPORTANT]
> `FIX-46` a `FIX-49`, `FIX-43` y `FIX-64` fueron promovidas a [`terminado-1.md`](terminado-1.md).
> Los problemas de sincronización de la suite de pruebas en este repositorio (`FIX-68` y ajuste de sincronización de `FIX-62` y `FIX-65`) fueron corregidos localmente y se validan con sus pruebas en verde.
> A continuación se detallan las tareas pendientes que requieren intervención directa por parte de los equipos de desarrollo en los submódulos (**`backend`**), con su causa raíz, archivos afectados y entregables precisos.

---

## 🛠️ Fixes detectados en la verificación del 09/10/2026


### `FIX-73` - Sincronización y manejo determinista de tareas consecutivas en acciones de energía (`BAC-21B` / `FIX-39` / `FIX-33`) (Backend / pruebas)

- **Área:** Backend / pruebas
- **Asignado:** Lisandro / Tayra
- **Estimación:** 1,5 h
- **Depende de:** `BAC-21B` (en `terminado-1.md`), `FIX-33`, `BAC-24A`.
- **Problema técnico detallado:**
  En `test/back/puente_etapa1_acceptance_test.go`, las pruebas `FIX-39 BAC-21B shutdown y reboot exigen FULL_ACCESS` y `BAC-21B las acciones de energía quedan en auditoria` fallan con:
  ```
  409 INSTANCE_BUSY: "La instancia se encuentra ejecutando otra tarea. Aguarde a que finalice."
  ```
  Al emitirse una orden de energía (ej. `stop` sobre la 101), el backend registra la tarea en estado `RUNNING`. Cuando de inmediato se despacha otra acción de ciclo de vida (`shutdown`, `reboot` u otro `stop`), el control de concurrencia de `FIX-33` rechaza la segunda orden con `409 INSTANCE_BUSY`. Aunque el rechazo con 409 es la conducta esperada ante instancias ocupadas, la prueba o el cliente que encadena acciones debe sincronizar la finalización de la tarea precedente antes de disparar la siguiente, y la función de sondeo `energyPath` no debe emitir peticiones mutantes con token administrativo que inicien tareas no esperadas.
- **Qué debe hacer el equipo:**
  1. En las pruebas de aceptación de ciclo de vida, esperar a que la tarea previa pase a `COMPLETED` o `FAILED` antes de despachar una nueva orden sobre el mismo VMID.
  2. Asegurar que `energyPath` inspeccione las rutas sin mutar el estado de la instancia en Proxmox.
  3. Validar que la transición de estados y la auditoría registren cada tarea con su UPID correspondiente.
- **Criterio de éxito:**
  - La suite `TestPuenteEtapa1` pasa 100% en verde sin fallos espurios de `INSTANCE_BUSY`.
  - Cada orden de energía se procesa y audita ordenadamente.




## 🛠️ Fixes detectados en la verificación del 13/10/2026 (Etapa 1)

### `FIX-83` - Refrescar el conteo de instancias ante `TASK_FINISHED` (`FRN-19C`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 1 h
- **Depende de:** `FRN-19C` (en `terminado-1.md`), `FRN-19B`, `FRN-17A`.
- **Problema técnico detallado:**
  Las tarjetas de conteo del Dashboard (`FRN-19C`) se renderizan desde `instancesSummary` (1 caso de `dashboard-summary.test.tsx` en verde), pero **no se refrescan al recibir `TASK_FINISHED`**: `useNodeStatus` sólo hace polling cada 10 s y no observa el canal de eventos, así que el conteo puede tardar hasta 10 s en cambiar (el criterio pide menos de 2 s).
- **Qué debe hacer el equipo de frontend:**
  1. Suscribir el Dashboard al canal (`useEventsContext` / `useEventSubscription`) y, ante cualquier `TASK_FINISHED` (o `INSTANCE_STATE_CHANGED`), volver a consultar `GET /api/node/status` sin loader.
  2. Ejecutar `npx vitest run dashboard-summary.test.tsx` hasta 2/2.
- **Criterio de éxito:** tras encender o apagar una instancia, el conteo cambia en menos de 2 s sin `F5`.

### `FIX-84` - Completar storages y contenido en el simulador de Proxmox (`BAC-30`) (Backend / Testing)

- **Área:** Backend / Testing
- **Asignado:** Lisandro
- **Estimación:** 1 h
- **Depende de:** `BAC-30` (en `terminado-1.md`).
- **Problema técnico detallado:**
  El simulador ya responde `POST /nodes/{node}/qemu|lxc` (UPID `qmcreate`/`vzcreate`), `PUT /nodes/{node}/{tipo}/{vmid}/config` y `GET /nodes/{node}/{tipo}/{vmid}/rrddata`, pero **no** implementa `GET /nodes/{node}/storage` ni `GET /nodes/{node}/storage/{storage}/content`, que el asistente de creación (`RF-07`) necesita. Hoy responden `501 Not Implemented`.
- **Qué debe hacer el equipo de backend:**
  1. Agregar al simulador `GET /nodes/{node}/storage` (storages `local` y `local-lvm` con tipo, total y libre) y `GET /nodes/{node}/storage/{storage}/content` (templates LXC e ISOs), con el mismo formato que Proxmox VE 9.x.
  2. Dejar `go test ./cmd/proxmox-simulador/...` en verde y la suite `test/back` (BAC-30) sin bloques rojos.
- **Criterio de éxito:** el desarrollo de la Etapa 2 puede consultar storages y contenido contra el simulador sin `501 Not Implemented`.

### `FIX-85` - Cerrar la confirmación de `Start` en los modales de energía (`FIX-76` / `FRN-15`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 1 h
- **Depende de:** `FIX-76` (en `terminado-1.md`), `FRN-15`.
- **Problema técnico detallado:**
  `test/front/instances-modals.test.tsx` quedó con **1 fallo**: `Start pide confirmación y cancelar no envía ninguna petición`. El resto de los casos (shutdown/reboot/stop/delete y el filtro por rol) pasa.
- **Qué debe hacer el equipo de frontend:**
  1. Revisar `InstanceAction` / `ConfirmUserAction` para que `Start` abra el modal de confirmación y que cancelarlo **no** dispare ninguna petición (`POST /instances/:vmid/start` o `/status/start`).
  2. Ejecutar `npx vitest run instances-modals.test.tsx` hasta 6/6.
- **Criterio de éxito:** `instances-modals.test.tsx` pasa 6 de 6.

### `FIX-86` - Completar la máquina de estados y el mapeo D2 de operaciones (`FIX-77` / `FRN-16`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2 h
- **Depende de:** `FIX-77` (en `terminado-1.md`), `FRN-16`, `BAC-29`.
- **Problema técnico detallado:**
  `test/front/instances-operation.test.tsx` queda con **4 fallos** (de 9) tras corregir en este repo la redacción esperada del 502 (`PROXMOX_UNAVAILABLE`), que ya no cuenta:
  1. `al confirmar envía la orden por la ruta del contrato y, mientras está en curso, bloquea toda la fila con un spinner`.
  2. `D2: 403 'INSTANCE_PROTECTED' muestra su mensaje y desbloquea la fila`.
  3. `D2: 403 'INSTANCE_ACCESS_DENIED' muestra su mensaje y desbloquea la fila`.
  4. `los seis mensajes de D2 son distintos entre sí`.
  Los 403 no exponen el mensaje propio del código: el flujo de `403` (aviso global de `apiClient`/`ApiResponseNotifier`) muestra el texto crudo del backend en lugar de la traducción de `errorMapper` (`INSTANCE_PROTECTED` → "infraestructura crítica"; `INSTANCE_ACCESS_DENIED` → "no posee permisos operativos"), mientras que los `409` y `PROXMOX_TIMEOUT` sí se traducen.
- **Qué debe hacer el equipo de frontend:**
  1. Mantener la fila bloqueada (spinner y botones deshabilitados) tras el `202` hasta el `TASK_FINISHED`, y desbloquearla ante el error.
  2. Garantizar un mensaje propio y distinto por cada código de D2 (y un mensaje genérico para el resto).
  3. Ejecutar `npx vitest run instances-operation.test.tsx` hasta 9/9.
- **Criterio de éxito:** `instances-operation.test.tsx` pasa 9 de 9.

---

