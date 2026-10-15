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

---


### `FIX-97` - Normalización de telemetría de inventario (escala porcentual 0-100, desacople de stopped y glifo em-dash) (`FRN-20C`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Luz
- **Estimación:** 0,5 h
- **Depende de:** `FRN-20C` (en `terminado-1.md`).
- **Problema técnico detallado:**
  A raíz de la unificación en el Backend de todas las métricas de CPU a porcentaje directo (escala `0 a 100`) y el envío de `null` garantizado en `cpuUsage` y `ramUsage` cuando la instancia está apagada (`stopped`), se identificaron las siguientes discrepancias en el Frontend:
  1. **Bug de multiplicación por 100 en `formatCpuUsage`:** En `frontend/centinela/src/components/features/instances/utils/formatInstanceTelemetry.ts`, la función utiliza `new Intl.NumberFormat('es-AR', { style: 'percent' })`. En JavaScript, `Intl.NumberFormat` con estilo porcentual multiplica por 100 la entrada esperando una fracción `0..1`. Al recibir del backend un valor en escala porcentual (ej. `25` por 25%), genera salidas distorsionadas como `"2.500 %"` (o `"93 %"` ante `0.93%`).
  2. **Intercepción fija en `Instances.tsx`:** En `frontend/centinela/src/pages/Instances.tsx:144-145`, el código evalúa `instance.status === 'stopped' ? '-' : formatCpuUsage(...)`. Esto hardcodea un guion simple `'-'` y puentea los formateadores. Al enviar el backend `null` cuando está apagada, la tabla debe delegar directamente en los formateadores sin bifurcar por `status`.
  3. **Glifo de ausencia en nulos:** Tanto `formatCpuUsage` como `formatRamUsage` retornan `'-'` (`U+002D`) en vez del em-dash `'—'` (`U+2014`) exigido por la especificación para campos nulos o ausentes.
- **Qué debe hacer el equipo de frontend:**
  1. En `formatInstanceTelemetry.ts`, adaptar `formatCpuUsage` para que reciba el valor en escala 0-100 (ej. pasando `usage / 100` a `Intl.NumberFormat` o formateando con sufijo ` %`) y retorne `'—'` (`U+2014`) ante valores `null`, `undefined` o negativos.
  2. En `formatInstanceTelemetry.ts`, asegurar que `formatRamUsage` retorne `'—'` ante valores `null`, `undefined` o negativos.
  3. En `Instances.tsx`, remover el ternario `instance.status === 'stopped' ? '-' : ...` en las filas de telemetría y actualizar la clase de centrado `cpuUsage === '—'` / `ramUsage === '—'`, delegando directamente en las funciones utilitarias.
  4. En `test/front/instances-metrics.test.tsx`, actualizar el mock de prueba `cpuUsage: 0.25` a `cpuUsage: 25` para reflejar el contrato unificado en porcentaje.
  5. Ejecutar `npx vitest run instances-metrics.test.tsx` y certificar 3/3 en verde.
- **Criterio de éxito:** `instances-metrics.test.tsx` pasa 3 de 3 en verde y los consumos de CPU se visualizan en escala real (ej. `25 %`) sin multiplicaciones indebidas ni guiones simples ante máquinas detenidas.



### `FIX-99` - Advertencia visible "Próximamente" en vista directa `/instances/new` (`FIX-89`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FIX-89` (en `terminado-1.md`).
- **Problema técnico detallado:**
  El commit `4e666a0` interceptó con éxito el botón "Crear instancia" en `Instances.tsx` para emitir un toast de advertencia. Sin embargo, si un usuario accede directamente por URL a `/instances/new`, la página `frontend/centinela/src/pages/createNewInstances.tsx` renderiza el asistente estático de 6 pasos sin ningún aviso o banner que informe que la funcionalidad se encuentra en desarrollo para la Etapa 2, provocando el fallo en `test/front/fixes-pending.test.tsx:118`:
  ```
  AssertionError: el asistente de /instances/new no advierte de forma explícita que está en desarrollo / próximamente: expected null not to be null
  ```
- **Qué debe hacer el equipo de frontend:**
  1. En `createNewInstances.tsx`, incorporar una alerta o badge destacado en la cabecera del formulario que indique explícitamente: `"Aprovisionamiento en desarrollo — Próximamente en Etapa 2"`.
  2. Ejecutar `npx vitest run fixes-pending.test.tsx -t "FIX-89"` y comprobar que pase en verde.
- **Criterio de éxito:** `fixes-pending.test.tsx` pasa el caso `FIX-89` en verde.

---

### `FIX-100` - Manejo seguro de valores `undefined` en `formatCount` del Dashboard (regresión commit `57c739f` / `FRN-19B`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz / Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FRN-19B` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En el commit `57c739f` (`frontend/centinela/src/pages/Dashboard.tsx:58`), se introdujo la función utilitaria de formateo de contadores:
  ```ts
  const formatCount = (value: number | null) => value === null ? '—' : value.toLocaleString('es');
  ```
  Y en la línea 133 (tarjeta de Plantillas para Admin):
  ```tsx
  <p className="text-metrica">{formatCount(templates.loading ? null : templates.count)}</p>
  ```
  Cuando `templates.count` es `undefined` (estado inicial de carga, respuesta pendiente o sin inicializar en mocks de pruebas), el argumento `value` recibe `undefined`. Como la guarda evalúa estrictamente `value === null`, la condición es falsa y se ejecuta `undefined.toLocaleString('es')`, lanzando una excepción no controlada:
  ```
  TypeError: Cannot read properties of undefined (reading 'toLocaleString')
  ```
  Esto provoca el desmontaje completo de la vista `<Dashboard>` (`<body><div /></body>`), rompiendo en cascada 12 pruebas en `test/front/dashboard-node.test.tsx` y 2 pruebas en `test/front/dashboard-summary.test.tsx`.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/pages/Dashboard.tsx:58`, adaptar `formatCount` y `formatShare` para emplear una guarda defensiva ante valores no numéricos (`undefined`, `null`, `NaN`):
     ```ts
     const formatCount = (value: number | null | undefined) => (typeof value !== 'number' || Number.isNaN(value)) ? '—' : value.toLocaleString('es');
     const formatShare = (value: number | null | undefined) => (typeof value !== 'number' || !totalInstances) ? '—' : `${Math.round((value / totalInstances) * 100)}%`;
     ```
  2. Ejecutar `npx vitest run dashboard-node.test.tsx` en `test/front` y certificar 12 de 12 pruebas en verde.
- **Criterio de éxito:**
  El Dashboard tolera contadores numéricos indefinidos o nulos sin lanzar `TypeError` en consola, y la suite `test/front/dashboard-node.test.tsx` pasa al 100% (12/12 en verde).

