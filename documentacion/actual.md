### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (19/10/2026):** según ClickUp, en este archivo quedan las tareas **pendientes de resolución** de la Etapa 1 y Etapa 2 y los `FIX` abiertos. Lo ya verificado con las suites está en [`terminado.md`](terminado.md) / [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md), y las correcciones de las implementadas con problemas, en [`futuro.md`](futuro.md) / [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md).
>
> **Pasados a `terminado-1.md` el 19/10/2026 (Etapa 1):**
> - `FIX-81`, `FRN-21` / `FIX-95`, `FRN-17B`, `FIX-86`, `FIX-97` (resueltos y con tests en verde).
> - `FIX-79`, `FIX-76`, `FIX-85`, `FRN-19C`, `FIX-83`, `FIX-99` (implementados con problemas; correcciones abiertas en `FIX-101`, `FIX-102`, `FIX-103` y `FIX-104` de `futuro-1.md`).
>
> **Pasados a `terminado2.md` el 19/10/2026 (Etapa 2):**
> - `BAC-39`, `BAC-40`, `BAC-41`, `BAC-50`, `BAC-51`, `BAC-52`, `BAC-53` (resueltos y con tests en verde).

---

### ETAPA 1: Tareas y Fixes Pendientes

#### `BAC-24C` - Limpieza de permisos al eliminar una instancia (nueva)
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.0 h
- **Depende de:** `BAC-24B` y `BAC-25B`.
- **Problema:** `/cluster/nextid` devuelve el menor VMID libre. Si se borra la 110 y después se crea otra instancia, recibe el 110 y **hereda las filas de `permisos_instancia`** de la borrada: un operador tendría acceso a una máquina que no se le asignó.
- **Entregable:**
  - Cuando una tarea de borrado termina con `COMPLETED`, eliminar en una transacción las filas de `permisos_instancia` de ese VMID.
  - Registrar la limpieza en `auditoria` con `accion: INSTANCE_PERMISSIONS_PURGED` y la lista de usuarios afectados.
  - Publicar el `TASK_FINISHED` **antes** de borrar los permisos, para que los operadores asignados lo reciban.
  - Si la tarea falla, no se toca nada.
- **Criterio de éxito:** después de borrar la 110, `GET /api/admin/users/:id/permissions` del operador ya no la incluye. Al crear una instancia nueva con el VMID 110, el operador no la ve.


### `FIX-91` - Eliminar el componente muerto `InstanceTester` (`FRN-15` / `FRN-16`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FRN-15`, `FRN-16` (en `terminado-1.md`).
- **Problema técnico detallado:**
  `frontend/centinela/src/components/features/instances/components/InstanceTester.tsx` se titula **"Probador de Instancias (Mock Local)"**, define `INITIAL_INSTANCES` con datos falsos (ids 1..3, "servidor-web-produccion", "base-datos-staging", etc.) y **simula** las transiciones con `setTimeout`, además de decir "Recargá la página para reiniciar el mock". **No se importa desde ningún archivo**: es el harness de desarrollo usado para `FRN-15`/`FRN-16`, y hoy es código muerto con datos mock que puede confundirse con código de producto.
- **Qué debe hacer el equipo de frontend:** eliminar el componente (o moverlo a `src/test` si se quiere conservar como utilidad de desarrollo).
- **Criterio de éxito:** no queda código de producto con datos mock ni componentes sin importar.


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

---

### ETAPA 2: Tareas Pendientes

#### `INF-10` - Verificación y disponibilidad de plantillas LXC e ISOs en el servidor Proxmox real

- **Área:** Infraestructura
- **Asignados:** Nico y Lucas
- **Estimación:** 1.5 h
- **Depende de:** `INF-07C` (permisos de token en Proxmox).
- **Problema:** Para que la validación en el servidor real (`INT-06`) funcione, los almacenamientos de Proxmox deben contener las imágenes base necesarias.
- **Entregable:**
  1. Descargar o verificar la presencia de una plantilla LXC liviana oficial (ej. `alpine-3.19-default` o `debian-12-standard`) en el storage `local:vztmpl`.
  2. Verificar una imagen ISO de prueba en `local:iso`.
  3. Comprobar que el API Token puede listar ambos contenidos mediante consulta HTTPS a `/storage/{storage}/content`.
- **Criterio de éxito:** El backend desplegado en el servidor lista las plantillas e ISOs reales a través de `GET /api/node/templates` sin errores de permisos.


#### `BAC-36` - Endpoint REST de métricas históricas de instancia (`timeframe=2h|24h`)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-34`.
- **Problema:** Al abrir la vista de una máquina, los gráficos necesitan renderizar de inmediato la tendencia histórica sin aguardar la llegada de los primeros paquetes del WebSocket.
- **Entregable:**
  1. Endpoint `GET /api/instances/:vmid/metrics?timeframe=2h|24h` protegido con `RequireAuth` y `RequireInstanceAccess`.
  2. Parámetro `timeframe` configurable (por defecto `2h`, soportando `24h`), mapeado a los intervalos de Proxmox (`hour` y `day`).
  3. Caché corta en Redis (TTL 30 s) para evitar consultas repetitivas al hipervisor ante recargas de página.
- **Criterio de éxito:** Responde en menos de 100 ms desde caché con el arreglo histórico de métricas formateado para el consumo de gráficos.


---



### `FIX-101` - Corrección de filtros reactivos en memoria (búsqueda de texto y alternancia en ejecución / detenidas) (`FIX-79` / `FRN-20B`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 1 h
- **Depende de:** `FIX-79`, `FRN-20B` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En `frontend/centinela/src/pages/Instances.tsx`, el commit `623bc28` expuso los botones visibles para el selector de tipo ("Todas", "VM", "LXC") resolviendo el acceso directo. Sin embargo, en `test/front/instances-ola2.test.tsx` persisten 2 fallos de filtrado reactivo en memoria:
  1. La búsqueda en memoria (`search` / `q`) sobre el nombre o ID de instancia no actualiza el listado visible al ingresar texto en el input de búsqueda.
  2. Al alternar entre los estados "En ejecución" (`running`) y "Detenidas" (`stopped`), la tabla no filtra adecuadamente el conjunto de filas visibles conservando los contadores sincronizados.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/pages/Instances.tsx`, revisar la pipeline de filtrado en memoria (`useMemo` sobre `instances` aplicando secuencialmente filtro de tipo, filtro de estado y filtro de texto del buscador).
  2. Asegurar que el cambio en el `input` de búsqueda actualice el estado local de forma reactiva y filtre las filas sin depender de consultas al backend.
  3. Ejecutar `npx vitest run instances-ola2.test.tsx` y certificar que la suite pase al 100% (6/6 en verde).
- **Criterio de éxito:**
  - `test/front/instances-ola2.test.tsx` pasa 6 de 6 pruebas en verde.


### `FIX-102` - Corrección de atributo de accesibilidad `aria-label` en `InstanceAction` para permisos de operador (`FIX-76` / `FIX-85` / `FRN-15`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Luz
- **Estimación:** 0,5 h
- **Depende de:** `FIX-76`, `FIX-85` (en `terminado-1.md`), `FRN-15`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/components/features/instances/components/InstanceAction.tsx` (líneas 222-266), las opciones de acción de energía dentro del menú desplegable (`DropdownMenuItem`) contienen un error tipográfico en el atributo de accesibilidad: están escritas como `arial-label="..."` (con 'l' intermedia) en lugar de `aria-label="..."`.
  Por este motivo, las herramientas de accesibilidad y las pruebas de `test/front/instances-modals.test.tsx` que buscan elementos por nombre accesible (`screen.getByRole('menuitem', { name: ... })`) no pueden localizar las acciones correspondientes a operadores cuando se evalúa la visibilidad de acciones de energía (`un OPERATOR ve acciones de energía solo donde tiene FULL_ACCESS y nunca ve Delete`).
- **Qué debe hacer el equipo de frontend:**
  1. En `InstanceAction.tsx:222-266`, corregir todos los atributos `arial-label` cambiándolos a `aria-label`.
  2. Verificar que los textos coincidan con las expresiones canónicas esperadas (`Iniciar`, `Reiniciar`, `Apagar`, `Detener`).
  3. Ejecutar `npx vitest run instances-modals.test.tsx` y certificar que pase 6 de 6 pruebas en verde.
- **Criterio de éxito:**
  - La suite `instances-modals.test.tsx` pasa al 100% (6/6 en verde).


### `FIX-103` - Manejo defensivo de arrays de telemetría de nodo en `useNodeStatus.ts` (`TypeError: serie.map is not a function`) (`FRN-19C` / `FIX-83`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 0,5 h
- **Depende de:** `FRN-19C`, `FIX-83` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En el commit `54efc3d` (`frontend/centinela/src/components/features/dashboard/hooks/useNodeStatus.ts:146`), se introdujo la transformación de series históricas de telemetría llamando directamente a `serie.map(...)`:
  ```ts
  const mapped = serie.map(...)
  ```
  Cuando la API o los mocks de prueba devuelven un objeto, una serie ausente o una respuesta que no es un array nativo, se produce una excepción no capturada `TypeError: serie.map is not a function`. Esto provoca el desmontaje completo del componente `<Dashboard>`, causando la falla en cascada de `test/front/dashboard-summary.test.tsx` (2 fallos) y `test/front/dashboard-node.test.tsx` (12 fallos).
- **Qué debe hacer el equipo de frontend:**
  1. En `useNodeStatus.ts:146`, validar que `Array.isArray(serie)` antes de invocar `.map()`, retornando un array vacío `[]` como valor por defecto si la serie no es iterable.
  2. Ejecutar `npx vitest run dashboard-summary.test.tsx` y certificar que pase en verde.
- **Criterio de éxito:**
  - `useNodeStatus` no arroja excepciones ante respuestas no-array y `dashboard-summary.test.tsx` pasa al 100% en verde.


### `FIX-104` - Visibilidad permanente del banner de advertencia "Próximamente" en todos los pasos de `/instances/new` (`FIX-99` / `FIX-89`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FIX-99`, `FIX-89` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En el commit `30894f4` (`frontend/centinela/src/pages/createNewInstances.tsx:55`), se agregó el banner de advertencia con el texto requerido: `"Aprovisionamiento en desarrollo — Próximamente en Etapa 2"`. Sin embargo, se condicionó con:
  ```tsx
  {currentStep >= 3 && (
      <div role="alert" className={styles.developmentNotice}>
          ...
      </div>
  )}
  ```
  Al ingresar directamente a `/instances/new`, el asistente arranca en el paso inicial (`currentStep === 0` o `1`), por lo que el banner no se renderiza. La suite `test/front/fixes-pending.test.tsx:118` evalúa la vista en su montaje inicial y falla con:
  ```
  AssertionError: el asistente de /instances/new no advierte de forma explícita que está en desarrollo / próximamente: expected null not to be null
  ```
- **Qué debe hacer el equipo de frontend:**
  1. En `createNewInstances.tsx`, remover la condición `{currentStep >= 3 && ...}` para que el banner de alerta `role="alert"` se muestre de forma permanente en la cabecera del asistente durante todos los pasos de la creación hasta que la funcionalidad esté plenamente operativa en la Etapa 2.
  2. Ejecutar `npx vitest run fixes-pending.test.tsx -t "FIX-89"` y certificar que pase en verde.
- **Criterio de éxito:**
  - `fixes-pending.test.tsx` pasa el caso `FIX-89` en verde.


### `FIX-104B` - Desacople de UPID y adopción del contrato seguro `tareaId` en Frontend (`instanceService.ts`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Luz
- **Estimación:** 1,0 h
- **Depende de:** `bd22c3d` (Backend: ocultamiento de UPID tras `tareaId`), `FRN-15`, `FRN-16` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En el backend (commit `bd22c3d` y contrato `BAC-29`), la respuesta HTTP `202 Accepted` de las acciones operativas (`start`, `stop`, `shutdown`, `reboot`, `delete`) dejó de exponer el campo privado `upid` por motivos de seguridad (evitando fugas de topología y nombres de nodos Proxmox), devolviendo únicamente el identificador opaco `{ "tareaId": "<uuid>" }` (`AccionAceptadaResponse`).
  Sin embargo, en el frontend:
  1. En `frontend/centinela/src/components/features/instances/services/instanceService.ts` (líneas 88-90 en `requestInstancePowerAction`):
     ```typescript
     if (!response || typeof response.upid !== 'string' || !response.upid.trim()) {
       throw new ApiRequestError('La orden fue aceptada, pero no se pudo interpretar su seguimiento...');
     }
     ```
  2. En la función `deleteInstance` (líneas 101-103):
     ```typescript
     if (!response || typeof response.upid !== 'string' || !response.upid.trim()) {
       throw new ApiRequestError('La eliminación fue aceptada, peno no se pudo interpretar su seguimiento.');
     }
     ```
  3. Los tipos en `types/instance.ts` declaran obligatoria la propiedad `upid: string`.
  Esto provoca que ante cualquier acción operativa de ciclo de vida (iniciar, apagar, reiniciar o eliminar), aunque el backend procese y acepte la orden con 202, el frontend arroje inmediatamente una excepción `ApiRequestError` al usuario por no encontrar `response.upid`.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/components/features/instances/types/instance.ts`:
     - Modificar `InstancePowerActionResponse` para que espere `{ tareaId: string; upid?: string }` (haciendo `upid` opcional o reemplazándolo por `tareaId`).
  2. En `instanceService.ts`:
     - En `requestInstancePowerAction`: reemplazar la validación de `response.upid` por la validación de `response.tareaId`:
       ```typescript
       if (!response || typeof response.tareaId !== 'string' || !response.tareaId.trim()) {
         throw new ApiRequestError('La orden fue aceptada, pero no se pudo interpretar su identificador de seguimiento.');
       }
       ```
     - En `deleteInstance`: validar la presencia de `response.tareaId` en lugar de `response.upid`.
  3. En modales y componentes (`ConfirmUserAction`, `InstanceAction`), usar `tareaId` para el seguimiento del evento `TASK_FINISHED`.
- **Criterio de éxito:**
  - El usuario puede encender, apagar, reiniciar y borrar instancias desde la interfaz sin recibir errores de seguimiento.
  - La suite `test/front/instances-operation.test.tsx` y `test/front/int02-e2e.test.ts` pasan en verde sin errores por ausencia de `upid`.



### `FIX-105` - Desacople de UPID y adopción del contrato seguro `tareaId` en Frontend (`instanceService.ts`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Luz
- **Estimación:** 1,0 h
- **Depende de:** `bd22c3d` (Backend: ocultamiento de UPID tras `tareaId`).
- **Problema técnico detallado:**
  En el backend (commit `bd22c3d`), el cuerpo HTTP `202 Accepted` de las acciones operativas y creación de instancias dejó de exponer el campo `upid` por motivos de seguridad, devolviendo únicamente el identificador opaco `{ "tareaId": "<uuid>" }` (`AccionAceptadaResponse`).
  Sin embargo, en el frontend:
  1. En `frontend/centinela/src/components/features/instances/services/instanceService.ts` (líneas 88-90):
     ```typescript
     if (!response || typeof response.upid !== 'string' || !response.upid.trim()) {
       throw new ApiRequestError('La orden fue aceptada, pero no se pudo interpretar su seguimiento...');
     }
     ```
  2. En la función `deleteInstance` (líneas 101-103):
     ```typescript
     if (!response || typeof response.upid !== 'string' || !response.upid.trim()) {
       throw new ApiRequestError('La eliminación fue aceptada, peno no se pudo interpretar su seguimiento.');
     }
     ```
  3. En `createInstance` (líneas 150-154):
     ```typescript
     if (typeof creation.upid !== 'string') {
       throw new ApiRequestError('La creación fue aceptada, pero no se pudo interpretar su seguimiento.');
     }
     ```
  4. Los tipos en `types/instance.ts` y `types/createInstance.ts` exigen obligatoriamente la propiedad `upid: string`.
  
  Esto provoca que **todas las acciones del frontend (iniciar, detener, apagar, reiniciar, borrar y crear) fallen de inmediato en la interfaz** arrojando excepciones al usuario, a pesar de que el backend procesó y aceptó la orden.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/components/features/instances/types/instance.ts`:
     - Modificar `InstancePowerActionResponse` para que espere `{ tareaId: string; upid?: string }` (haciendo `upid` opcional o eliminándolo del contrato público).
  2. En `frontend/centinela/src/components/features/instances/types/createInstance.ts`:
     - Actualizar `CreateInstanceResponse` para que no exija `upid` obligatorio: `{ vmid: number; tareaId: string }`.
  3. En `instanceService.ts`:
     - En `requestInstancePowerAction`: reemplazar la validación de `response.upid` por la validación estricta de `response.tareaId`:
       ```typescript
       if (!response || typeof response.tareaId !== 'string' || !response.tareaId.trim()) {
         throw new ApiRequestError('La orden fue aceptada, pero no se pudo interpretar su identificador de seguimiento.');
       }
       ```
     - En `deleteInstance`: validar la presencia de `response.tareaId` en lugar de `response.upid`.
     - En `createInstance`: validar `creation.vmid` y `creation.tareaId`, eliminando la exigencia de `creation.upid`.
  4. En los modales y hooks que rastrean tareas activas (`InstanceActionModal`, `useInstances`), correlacionar la tarea en curso mediante `tareaId` (coincidente con `activeTask.tareaId` de la telemetría del inventario y los eventos Server-Sent Events).
- **Criterio de éxito:**
  - El usuario puede encender, apagar, reiniciar y borrar instancias desde la interfaz gráfica sin recibir alertas de error.
  - La suite de tests unitarios de frontend pasa en verde sin aserciones rotas por ausencia de `upid`.