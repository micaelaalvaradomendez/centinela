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


### `FIX-101` - Corrección de filtros reactivos en memoria (búsqueda de texto y alternancia en ejecución / detenidas) (`FIX-79` / `FRN-20B`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 1 h
- **Depende de:** `FIX-79`, `FRN-20B` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En `frontend/centinela/src/pages/Instances.tsx`, la pipeline de filtrado en memoria (`useMemo` sobre `instances`) no actualiza reactivamente el listado visible al ingresar texto en el input de búsqueda ni conserva los contadores sincronizados al alternar entre estados "En ejecución" (`running`) y "Detenidas" (`stopped`).
- **Qué debe hacer el equipo de frontend:**
  1. En `Instances.tsx`, asegurar que el cambio en el `input` de búsqueda actualice el estado local de forma reactiva y filtre las filas en memoria sin llamadas extra al backend.
  2. Ajustar el filtrado por estado (`running` / `stopped`).
  3. Certificar que `test/front/instances-ola2.test.tsx` pase 6 de 6 pruebas en verde.
- **Criterio de éxito:** `test/front/instances-ola2.test.tsx` pasa al 100% en verde.


### `FIX-103` - Manejo defensivo de arrays de telemetría de nodo en `useNodeStatus.ts` (`TypeError: serie.map is not a function`) (`FRN-19C` / `FIX-83`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 0,5 h
- **Depende de:** `FRN-19C`, `FIX-83` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En `frontend/centinela/src/components/features/dashboard/hooks/useNodeStatus.ts:146`, se invoca directamente `serie.map(...)`. Cuando la API o los mocks devuelven un objeto envoltorio (`{ items }` o `{ data }`) o un valor no iterable, se produce la excepción `TypeError: serie.map is not a function`, desmontando el componente `<Dashboard>`.
- **Qué debe hacer el equipo de frontend:**
  1. En `useNodeStatus.ts:146`, validar `if (!Array.isArray(serie)) return previous || [];` antes de invocar `.map()`.
  2. Certificar que `dashboard-summary.test.tsx` y `fix-103-node-history.test.tsx` pasen al 100% en verde.
- **Criterio de éxito:** `useNodeStatus` no arroja excepciones ante respuestas no-array.


### `FIX-104` - Visibilidad permanente del banner de advertencia "Próximamente" en todos los pasos de `/instances/new` (`FIX-99` / `FIX-89`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FIX-99`, `FIX-89` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En `frontend/centinela/src/pages/createNewInstances.tsx`, el banner con texto `"Aprovisionamiento en desarrollo — Próximamente en Etapa 2"` se condicionó con `{currentStep >= 3 && ...}`. Al ingresar a `/instances/new` en el paso inicial (`currentStep === 0` o `1`), el banner no se renderiza y `test/front/fixes-pending.test.tsx:118` falla.
- **Qué debe hacer el equipo de frontend:**
  1. En `createNewInstances.tsx`, remover la condición `{currentStep >= 3 && ...}` para que el banner `role="alert"` sea visible permanentemente en la cabecera en todos los pasos.
  2. Certificar que `fixes-pending.test.tsx` pase el caso `FIX-89` en verde.
- **Criterio de éxito:** `fixes-pending.test.tsx` pasa el caso `FIX-89` en verde.


