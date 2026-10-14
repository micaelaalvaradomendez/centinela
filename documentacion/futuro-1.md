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


### `FIX-91` - Eliminar el componente muerto `InstanceTester` (`FRN-15` / `FRN-16`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FRN-15`, `FRN-16` (en `terminado-1.md`).
- **Problema técnico detallado:**
  `frontend/centinela/src/components/features/instances/components/InstanceTester.tsx` se titula **"Probador de Instancias (Mock Local)"**, define `INITIAL_INSTANCES` con datos falsos (ids 1..3, "servidor-web-produccion", "base-datos-staging", etc.) y **simula** las transiciones con `setTimeout`, además de decir "Recargá la página para reiniciar el mock". **No se importa desde ningún archivo**: es el harness de desarrollo usado para `FRN-15`/`FRN-16`, y hoy es código muerto con datos mock que puede confundirse con código de producto.
- **Qué debe hacer el equipo de frontend:** eliminar el componente (o moverlo a `src/test` si se quiere conservar como utilidad de desarrollo).
- **Criterio de éxito:** no queda código de producto con datos mock ni componentes sin importar.