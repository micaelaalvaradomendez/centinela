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
> A continuación se detallan las tareas pendientes que requieren intervención directa por parte de los equipos de desarrollo en los submódulos (**`backend`** y **`frontend`**), con su causa raíz, archivos afectados y entregables precisos.

---

## 🛠️ Fixes detectados en la verificación del 09/10/2026



## 🛠️ Fixes detectados en la verificación del 13/10/2026 (Etapa 1)


---


## 🛠️ Stubs y placeholders detectados en los submódulos (13/10/2026)

> Revisión de `backend` (`6afc017`) y `frontend` (`7fcc0bd`) buscando stubs/placeholders en código ya declarado terminado. Cada hallazgo cita el archivo y la tarea de `terminado-1.md` que afecta.


---

## 🛠️ Fixes detectados en la verificación del 19/10/2026 y 21/10/2026 (Etapa 1)

Verificados contra frontend `b0a7105` y backend `a25d1b3`.

---

### `FIX-91` - Eliminar el componente muerto `InstanceTester` (`FRN-15` / `FRN-16`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FRN-15`, `FRN-16` (en `terminado-1.md`).
- **Problema técnico detallado:**
  `frontend/centinela/src/components/features/instances/components/InstanceTester.tsx` se titula **"Probador de Instancias (Mock Local)"**, define `INITIAL_INSTANCES` con datos falsos (ids 1..3, "servidor-web-produccion", "base-datos-staging", etc.) y simula las transiciones con `setTimeout`. No se importa desde ningún archivo de producción y persiste como código muerto con datos mock.
- **Qué debe hacer el equipo de frontend:** eliminar el archivo `InstanceTester.tsx`.
- **Criterio de éxito:** no queda código de producto con datos mock ni componentes sin importar.


> [!NOTE]
> **Verificación del 22/10/2026 — retirados de este listado:** `FIX-101`, `FIX-103` y `FIX-104` no eran
> defectos del submódulo `frontend`, sino expectativas desactualizadas de las pruebas de `test/front`:
> - `FIX-101`: `Instances.tsx` ya filtra en memoria de forma reactiva; el test buscaba una pestaña
>   "En ejecución" inexistente en lugar del filtro real (`Filtros` → `Encendidas` / `Detenidas`).
> - `FIX-103`: `useNodeStatus.ts` ya valida `Array.isArray(serie)`; las pruebas fallaban por simular
>   `GET /node/templates` con un objeto envoltorio cuando el contrato real responde un array.
> - `FIX-104`: el asistente `/instances/new` dejó de ser un cascarón estático y hoy es el wizard
>   funcional de 5 pasos; ya no corresponde exigir el banner "Próximamente". `fixes-pending.test.tsx`
>   ahora valida el wizard real.
>
> Pruebas ajustadas y en verde: `instances-ola2.test.tsx` (6/6) y `fix-103-node-history.test.tsx` (7/7).
> `FIX-91` permanece porque sí requiere una acción de implementación (eliminar `InstanceTester.tsx`).


