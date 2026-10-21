### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Estado del tablero (21/10/2026):** según ClickUp, en este archivo quedan las tareas **pendientes de implementación** de la Etapa 1 y Etapa 2. Lo ya verificado con las suites está en [`terminado.md`](terminado.md) / [`terminado-1.md`](terminado-1.md) / [`terminado2.md`](terminado2.md), y las correcciones de las implementadas con problemas, en [`futuro.md`](futuro.md) / [`futuro-1.md`](futuro-1.md) / [`futuro2.md`](futuro2.md).
>
> **Pasados a `terminado-1.md` el 21/10/2026 (Etapa 1):**
> - `FIX-104B` / `FIX-105` y `FIX-102` (resueltos y con tests en verde).
> - `FIX-91`, `FIX-101`, `FIX-103`, `FIX-104` (implementados con problemas / correcciones registradas en `futuro-1.md`).
>
> **Pasados a `terminado2.md` el 21/10/2026 (Etapa 2):**
> - `BAC-36`, `FIX-107`, `FIX-109`, `FIX-110`, `FIX-111A`, `FIX-112`, `FIX-114`, `FIX-115A`, `FIX-116`, `FIX-117A`, `FIX-118`, `FIX-119` (resueltos y con tests en verde).
> - `FRN-24`, `FRN-25`, `INF-10` (implementados con problemas / correcciones registradas en `futuro2.md` como `FIX-120`, `FIX-121` y `FIX-123`).

---

### ETAPA 1: Tareas Pendientes de Implementación

#### `BAC-24C` / `FIX-113` - Limpieza de permisos al eliminar una instancia
- **Área:** Backend
- **Asignada:** Tayra / Lisandro
- **Estimación:** 1.0 h
- **Depende de:** `BAC-24B` y `BAC-25B`.
- **Importancia:** 🔴 **ALTA** (CWE-284: Improper Access Control, CWE-639: Authorization Bypass Through User-Controlled Key).
- **Problema:** `/cluster/nextid` devuelve el menor VMID libre. Si se borra la 110 y después se crea otra instancia, recibe el 110 y **hereda las filas de `permisos_instancia`** de la borrada: un operador tendría acceso a una máquina que no se le asignó.
- **Entregable:**
  - Cuando una tarea de borrado termina con `COMPLETED`, eliminar en una transacción las filas de `permisos_instancia` de ese VMID.
  - Registrar la limpieza en `auditoria` con `accion: INSTANCE_PERMISSIONS_PURGED` y la lista de usuarios afectados.
  - Publicar el `TASK_FINISHED` **antes** de borrar los permisos, para que los operadores asignados lo reciban.
  - Si la tarea falla, no se toca nada.
- **Criterio de éxito:** después de borrar la 110, `GET /api/admin/users/:id/permissions` del operador ya no la incluye. Al crear una instancia nueva con el VMID 110, el operador no la ve.


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

### ETAPA 2: Tareas Pendientes de Implementación

### `FIX-108` - Adaptación del Frontend para consumo de identificadores de instancia ofuscados (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1,5 h
- **Depende de:** `FIX-107`.
- **Problema técnico detallado:**
  Actualmente, el frontend modela el identificador de instancia estrictamente como un número entero (`vmid: number` y `id: number` en `InventoryInstance`), validándolo con `Number.isSafeInteger(vmid) && vmid > 0` en `validateInstanceVmid`.
  Al migrar el backend hacia identificadores opacos/cifrados no secuenciales (ej. `"inst_k8Pq2v"`), el frontend rechaza las respuestas con errores de validación de esquema en `instanceService.ts`.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/components/features/instances/types/instance.ts`:
     - Actualizar el tipo `InventoryInstance`:
       ```typescript
       export interface InventoryInstance {
         id: string | number; // Soporta el identificador opaco público
         name: string;
         type: 'VM' | 'LXC';
         ...
       }
       ```
  2. En `instanceService.ts`:
     - Flexibilizar `validateInstanceVmid` o renombrarlo a `validateInstanceId(id: string | number): string` para admitir tanto cadenas alfanuméricas ofuscadas como enteros legados, sanitizando que no contengan caracteres especiales ni inyecciones.
     - Actualizar las URLs de invocación (`/instances/${validId}/${action}`) para enviar el string del identificador público.
  3. En las rutas de la aplicación (`src/routes` o componentes con parámetros de URL `/instances/:id`):
     - Asegurar que los componentes de detalle de instancia y modales de acción reciban y operen con el identificador en formato string sin forzar conversiones numéricas destructivas (`parseInt` o `Number(...)`).
- **Criterio de éxito:**
  - La tabla de instancias renderiza correctamente el listado proveniente del backend con IDs ofuscados.
  - Los botones de acción de ciclo de vida (iniciar, apagar, ver métricas) operan de manera transparente utilizando el identificador opaco en las URLs del navegador.


### `FIX-111B` - Modal de confirmación con re-autenticación (Password + 2FA) en cambio de email (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1,0 h
- **Depende de:** `FIX-111A`.
- **Importancia:** 🔴 **CRÍTICA** (Mitigación UI contra secuestro de cuenta).
- **Problema técnico detallado:**
  Actualmente los componentes de visualización y edición del perfil de usuario permiten editar el input de email y enviar directamente el payload `PUT /api/account/profile` sin requerir credenciales de confirmación. Con la implementación del backend seguro en `FIX-111A`, cualquier intento de actualización de email sin los parámetros de re-autenticación es rechazado con error 400/401.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/components/features/users/components/informationOfUser.tsx` y `DetailsUserPage.tsx`:
     - Al detectar una edición en el campo de correo electrónico, interceptar el guardado y desplegar un modal de confirmación de seguridad.
     - El modal debe solicitar la contraseña actual del usuario y el código de seis dígitos de su aplicación autenticadora (TOTP 2FA).
     - Despachar la solicitud `PUT /api/account/profile` adjuntando `contrasenaActual` y `codigo2fa`.
     - Presentar al usuario un banner o notificación de feedback indicando que se ha enviado un correo de confirmación a la nueva dirección para activar el cambio.
- **Criterio de éxito:**
  - La interfaz gráfica no permite enviar el cambio de correo sin solicitar contraseña actual y código 2FA.
  - El usuario recibe retroalimentación visual clara sobre el flujo de confirmación por doble canal.


### `FIX-115B` - Actualización de validación y formularios de contraseñas de alta longitud en Frontend (Frontend)

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 0,5 h
- **Depende de:** `FIX-115A`.
- **Importancia:** 🟠 **MEDIA / ALTA** (Experiencia de usuario y seguridad NIST SP 800-63B).
- **Problema técnico detallado:**
  En el frontend, la función de utilidad `validatePasswordComplexity.ts` y las reglas de validación en formularios de cambio y reseteo restringen la entrada a un máximo de 12 caracteres, impidiendo al usuario pegar o escribir contraseñas robustas generadas por administradores de contraseñas.
- **Qué se debe hacer:**
  1. En `frontend/centinela/src/components/features/auth/utils/validatePasswordComplexity.ts`:
     - Actualizar la regla de longitud máxima para admitir hasta 128 caracteres (`password.length <= 128`).
  2. En los formularios de cambio y reseteo de contraseña (`ChangePassword.tsx`, `RecoverPassword.tsx`):
     - Modificar los inputs para permitir hasta 128 caracteres (`maxLength={128}`).
     - Actualizar los textos de ayuda y mensajes de validación eliminando la mención restrictiva de 12 caracteres.
- **Criterio de éxito:**
  - El frontend valida correctamente contraseñas largas de hasta 128 caracteres sin bloquear el botón de envío ni mostrar mensajes de error espurios.
  - Los formularios permiten pegar credenciales generadas por administradores de contraseñas estándares.


### `FIX-117B` - Adaptación del Dashboard de operador ante segregación de métricas globales (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1,0 h
- **Depende de:** `FIX-117A`.
- **Importancia:** 🟡 **MEDIA** (Resiliencia UI y prevención de errores por campos omitidos).
- **Problema técnico detallado:**
  Al implementarse `FIX-117A`, el backend deja de enviar métricas globales del host (`instancesSummary`, almacenamiento global) cuando quien consulta es un usuario con rol `OPERATOR`. Los componentes del Dashboard principal que asumen la presencia incondicional de estos objetos arrojan excepciones en tiempo de ejecución o renderizan valores vacíos/NaN.
- **Qué se debe hacer:**
  1. En los componentes de Dashboard y resumen de métricas del frontend:
     - Condicionar la visualización de tarjetas de métricas globales de hardware y almacenamiento según el rol del usuario autenticado (`user.rol === 'ADMIN'`).
     - Para usuarios con rol `OPERATOR`, renderizar un resumen enfocado exclusivamente a las instancias asignadas y accesibles, sin intentar acceder a las propiedades omitidas por el backend.
- **Criterio de éxito:**
  - El Dashboard de un operador renderiza limpiamente sin advertencias ni errores de consola ante la respuesta acotada de `/node/status`.


#### `FRN-22` - Hook `useInstanceMetrics` y cliente WebSocket para telemetría en tiempo real

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** `BAC-35`, `FRN-21` (página de detalle).
- **Problema:** El cliente frontend debe gestionar la conexión WebSocket de métricas de forma resiliente, solicitando el ticket de autenticación previo y manejando caídas de red sin fugar memoria ni sockets huérfanos.
- **Entregable:**
  1. Servicio de solicitud de ticket efímero `getWsTicket()`.
  2. Hook personalizado `useInstanceMetrics(vmid)` que gestione conexión, estado (`conectando`, `conectado`, `reconectando`, `error`) y almacenamiento en búfer de los últimos $N$ puntos de telemetría.
  3. Reintento con backoff exponencial ante desconexión accidental.
  4. Cierre explícito del socket al desmontar el componente o cambiar de pestaña.
  5. **Contrato de telemetría unificado (escala y nulos):**
     - La CPU transmitida en vivo por WebSocket (`MetricasEnVivoDTO.cpu`) y en el histórico RRDdata (`PuntoMetricaDTO.cpuPercent`) viaja **en escala porcentual directa ($0$ a $100$)**, no como fracción $0..1$. El hook debe almacenar la telemetría en dicha escala sin multiplicar por 100.
     - Cuando la instancia esté apagada (`stopped`), el backend envía `cpu: null` y `mem: null`. El hook debe manejar valores nulos de forma segura (sin invocar `.toFixed()` directamente sobre `null`) e indicar en el estado del búfer la inactividad temporal.
- **Criterio de éxito:** El hook entrega series de datos actualizadas cada 2 s sin fugas de memoria en pruebas de componente con Vitest y maneja instancias apagadas (`null`) sin excepciones de tipo en consola.


#### `FRN-23` - Componentes de gráficos interactivos de rendimiento y selector de rango temporal

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.5 h
- **Depende de:** `FRN-22`, `BAC-36`.
- **Problema:** La pestaña "Métricas" de `InstanceDetailPage.tsx` carece de componentes visuales para interpretar el consumo de recursos de la máquina.
- **Entregable:**
  1. Incorporación de componentes de gráficos interactivos (utilizando Recharts o Chart.js):
     - Gráfico de CPU (%) en escala directa $0-100\%$ (compatible con `PuntoMetricaDTO.cpuPercent` y WebSocket).
     - Gráfico de Memoria RAM (Usada vs Total).
     - Gráfico de Tráfico de Red (KB/s In/Out).
     - Gráfico de I/O de Disco (KB/s Read/Write).
  2. Selector de rango temporal: "En vivo" (WebSocket), "Últimas 2 horas" y "Últimas 24 horas" (REST `BAC-36`).
  3. Formateo legible de unidades (bytes a MB/GB, bps a KB/s).
  4. Tooltips explicativos y estados de carga (skeletons) mientras se obtienen los primeros datos.
  5. **Tolerancia a valores nulos:** cuando la máquina esté detenida (`cpu: null`, `mem: null`), los gráficos deben pintar discontinuidades temporales o mostrar el estado "Instancia detenida" sin romper el ciclo de vida del componente con errores `NaN` o `TypeError`.
- **Criterio de éxito:** Al abrir la pestaña "Métricas", se visualiza la serie histórica y los gráficos continúan actualizándose dinámicamente con los eventos en vivo, tolerando paradas de la máquina sin fallos visuales.
