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

### `FIX-72` - Sincronizar firma de `NewSeguimientoTareas` en pruebas unitarias de servicios (`BAC-25A`) (Backend / pruebas)

- **Área:** Backend / pruebas
- **Asignado:** Lisandro / Lucas
- **Estimación:** 1 h
- **Depende de:** `BAC-25A` (en `terminado-1.md`).
- **Problema técnico detallado:**
  Al ejecutar las pruebas unitarias del backend invocadas por la prueba de aceptación `TestEtapa1/BAC-25A` (`etapa1_acceptance_test.go:185`), la suite `go test ./internal/core/services/...` falla en compilación:
  ```
  internal/core/services/seguimiento_pool_test.go:89:60: too many arguments in call to services.NewSeguimientoTareas
  	have (*proxmoxPool, *tareasEnMemoria, *publicadorFalso, *auditoriaGrabadora, services.ConfigSeguimiento)
  	want (ports.ProxmoxPort, ports.TareaRepository, services.ConfigSeguimiento)
  internal/core/services/seguimiento_tareas_test.go:131:54: too many arguments in call to services.NewSeguimientoTareas
  	have (*proxmoxTareas, *tareasEnMemoria, *publicadorFalso, *auditoriaGrabadora, services.ConfigSeguimiento)
  	want (ports.ProxmoxPort, ports.TareaRepository, services.ConfigSeguimiento)
  ```
  La firma del constructor `services.NewSeguimientoTareas` fue simplificada para recibir únicamente `(ports.ProxmoxPort, ports.TareaRepository, services.ConfigSeguimiento)`, pero los archivos de pruebas unitarias de servicios (`seguimiento_pool_test.go` y `seguimiento_tareas_test.go`) no fueron actualizados y siguen pasando 5 argumentos con stubs residuales.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/core/services/seguimiento_pool_test.go`, actualizar la llamada a `services.NewSeguimientoTareas` para pasar únicamente los 3 parámetros requeridos.
  2. En `backend/internal/core/services/seguimiento_tareas_test.go`, actualizar la llamada a `services.NewSeguimientoTareas` a la firma de 3 parámetros.
  3. Ejecutar `go test -v ./internal/core/services/...` en el repositorio backend y verificar que compilen y pasen todas las pruebas unitarias.
- **Criterio de éxito:**
  - `go test ./internal/core/services/...` compila limpiamente y pasa en verde.
  - La prueba de aceptación `TestEtapa1/BAC-25A` en `test/back/etapa1_acceptance_test.go` pasa 100% en verde.


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


### `FIX-76` - Desacoplar guarda de `nivelAcceso` para rol `ADMIN` en modales de energía (`FRN-15`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Luz
- **Estimación:** 0,5 h
- **Depende de:** `FRN-15`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/components/features/instances/components/InstanceAction.tsx:87`:
  ```tsx
  const canOperate = canOperateInstance(instance.id) && instance.nivelAcceso !== 'READ_ONLY';
  ```
  Un usuario con rol `ADMIN` dispone de privilegios absolutos para operar cualquier instancia sin estar restringido por las listas de acceso granular de operadores (`nivelAcceso`). Sin embargo, la condición evalúa `instance.nivelAcceso !== 'READ_ONLY'` sin contemplar `isAdmin()`. Cuando un fixture o instancia mockeada contiene `nivelAcceso: 'READ_ONLY'` (asignado para un perfil de operador), un administrador queda bloqueado y la acción `start` ("Encender") no se renderiza en la fila ni abre su modal de confirmación, provocando el fallo en `instances-modals.test.tsx`:
  ```
  AssertionError: la fila no ofrece la acción "start" (se busca un botón o ítem de menú con nombre /^(iniciar|encender|start)\b/i)
  ```
- **Qué debe hacer el equipo de frontend:**
  1. En `InstanceAction.tsx`, actualizar la asignación de `canOperate` para que un administrador siempre esté habilitado a operar:
     ```tsx
     const canOperate = isAdmin() || (canOperateInstance(instance.id) && instance.nivelAcceso !== 'READ_ONLY');
     ```
  2. Aplicar la misma corrección en la validación interna de `confirmPowerAction` (`InstanceAction.tsx:114`).
  3. Ejecutar `npx vitest run instances-modals.test.tsx` y certificar que la prueba `Start pide confirmación` pase en verde.
- **Criterio de éxito:**
  - Un usuario `ADMIN` puede disparar la acción de inicio/encendido y ver su modal de confirmación sin ser bloqueado por `nivelAcceso`.
  - La suite `instances-modals.test.tsx` pasa al 100% (6/6 en verde).


### `FIX-77` - Mapeo de errores HTTP de ciclo de vida (`BAC-29`) y resiliencia de la máquina de estados (`FRN-16`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian / Belinda
- **Estimación:** 1,5 h
- **Depende de:** `FRN-16`, `BAC-29`.
- **Problema técnico detallado:**
  Al despachar una orden de ciclo de vida sobre una instancia (`requestInstancePowerAction`), ante un rechazo del servidor la interfaz no mapea los códigos HTTP y claves de error acordados en el contrato D2 (`BAC-29`) a sus mensajes descriptivos específicos, y `ConfirmUserAction` suprime o desvía los errores 403. Al ejecutar `npx vitest run instances-operation.test.tsx`, fallan las pruebas por no encontrar los mensajes esperados:
  - `409 INSTANCE_INVALID_STATE`: debe mostrar *"La instancia ya está encendida"* o *"El estado actual no corresponde a la acción"*.
  - `409 INSTANCE_BUSY`: debe mostrar *"La instancia se encuentra ejecutando otra tarea. Aguarde a que finalice."*.
  - `403 INSTANCE_PROTECTED`: debe indicar que la instancia es de infraestructura protegida.
  - `403 INSTANCE_ACCESS_DENIED`: debe advertir que el operador no posee permisos sobre esa máquina.
  - `502 PROXMOX_UNAVAILABLE`: debe informar que el hipervisor Proxmox no está disponible.
  - `504 PROXMOX_TIMEOUT`: debe avisar que la respuesta de Proxmox agotó el tiempo de espera.
  Además, al ocurrir un error, todos los controles de la fila deben desbloquearse inmediatamente saliendo del estado de spinner.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/components/features/instances/components/InstanceAction.tsx` (o servicio asociado), implementar una función de traducción que reciba el `ApiRequestError` (evaluando `status` y `errorCode`) y genere el mensaje exacto exigido por el contrato D2.
  2. Asegurar que ante cualquier fallo se notifique el mensaje al usuario (vía toast o texto en modal), se resetee `pendingAction` y la fila recupere su interactividad normal.
  3. Verificar que `npx vitest run instances-operation.test.tsx` pase en verde para todos los códigos de error documentados.
- **Criterio de éxito:**
  - Cada uno de los seis códigos de error despliega su mensaje diferenciado y la fila se desbloquea tras el fallo.
  - La suite `instances-operation.test.tsx` pasa 9 de 9 pruebas en verde.


### `FIX-78` - Desduplicación de indicadores de salud del nodo y accesibilidad en Dashboard (`FRN-19B`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Luz
- **Estimación:** 1 h
- **Depende de:** `FRN-19B`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/pages/Dashboard.tsx`, el estado de salud del nodo (`nodeHealth.label`: "Saludable", "Advertencia", "Inaccesible") se renderiza simultáneamente en dos componentes:
  1. En la cabecera de la tarjeta de recursos (`Dashboard.tsx:136`): `<Badge variant="outline">{nodeHealth.label}</Badge>`.
  2. En el pie de página de estado del sistema (`Dashboard.tsx:211`): `<p className="text-caption">{nodeHealth.label}</p>`.
  Esta duplicación genera colisiones al evaluar la interfaz con lectores de pantalla o Testing Library (`findByText(/saludable/i)` y `findByText(/inaccesible/i)` fallan con `Found multiple elements with the text of: ...`). Asimismo, cuando el nodo está cargando o en error 502/504, los skeletons y las alertas de accesibilidad deben identificarse unívocamente sin interferir con otras tarjetas.
- **Qué debe hacer el equipo de frontend:**
  1. En `Dashboard.tsx`, diferenciar semánticamente los textos o asignar atributos de accesibilidad inequívocos (`aria-label="Estado general del nodo: Saludable"` o dejar el semáforo principal en la tarjeta de recursos y usar una descripción extendida o etiqueta diferenciada en el pie).
  2. Asegurar que las respuestas 502 y 504 reflejen inequívocamente la etiqueta `Inaccesible` en el semáforo de recursos.
  3. Ejecutar `npx vitest run dashboard-node.test.tsx` y comprobar que todos los casos pasen en verde.
- **Criterio de éxito:**
  - `dashboard-node.test.tsx` pasa al 100% sin colisiones de texto duplicado.
  - La experiencia visual y de accesibilidad es nítida tanto para usuarios videntes como para lectores de pantalla.


### `FIX-79` - Exposición directa de controles para selector reactivo por tipo de recurso (`FRN-20B`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz
- **Estimación:** 1 h
- **Depende de:** `FRN-20B`.
- **Problema técnico detallado:**
  La lógica de filtrado reactivo en memoria por tipo (`VM`, `LXC`, `all`) y buscador ya fue implementada en `frontend/centinela/src/pages/Instances.tsx`. No obstante, el selector fue colocado dentro de un menú desplegable contextual (`DropdownMenu`), requiriendo clics previos para ser descubierto. La suite de pruebas `test/front/instances-ola2.test.tsx:120` exige que los controles por tipo de recurso estén expuestos de forma visible y directa en la barra de control (como pestañas `tab` o botones accesibles `button` con etiquetas legibles como "Todas", "VM", "LXC"):
  ```
  AssertionError: debe existir un control para filtrar reactivamente por máquinas virtuales (VM): expected null not to be null
  ```
- **Qué debe hacer el equipo de frontend:**
  1. En `Instances.tsx`, complementar o reemplazar el dropdown por botones o pestañas accesibles y visibles en la barra de inventario que permitan alternar directamente entre "Todas", "VM (Máquinas virtuales)" y "LXC (Contenedores)".
  2. Mantener la reactividad instantánea en memoria sin peticiones al backend y la actualización sincronizada de los contadores.
  3. Ejecutar `npx vitest run instances-ola2.test.tsx` y certificar que la suite pase al 100%.
- **Criterio de éxito:**
  - El selector por tipo de recurso está directamente visible y accesible sin pasos ocultos.
  - La suite `instances-ola2.test.tsx` pasa 6 de 6 pruebas en verde.



