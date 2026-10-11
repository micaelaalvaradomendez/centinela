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



