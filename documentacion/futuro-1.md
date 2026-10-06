# Etapa 1: correcciones pendientes (FIX)

Correcciones de tareas de la Etapa 1 que ya están implementadas pero no cumplen del todo su criterio de éxito. Las tareas originales están en [`terminado-1.md`](terminado-1.md).

---

## 🛠️ Fixes detectados en la verificación del 03/10/2026

### `FIX-43` - Corregir maquetado, visualización de IP, badges de estado y columnas en tabla de instancias (`FRN-20A` / RF-03) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz / Cristian (PR #74)
- **Estimación:** 1.0 h
- **Depende de:** `FRN-20A` (en `terminado-1.md`).
- **Problema y evidencia:**
  El PR #74 implementó la integración viva de `Instances.tsx` con `useInstances.ts` y `instanceService.ts`, consumiendo `GET /instances` y aplicando permisos. Sin embargo, la suite de pruebas de aceptación (`test/front/instances-table.test.tsx`) falla 5/5 por discrepancias de maquetado e interfaz:
  1. **Nombre y VMID:** `Instances.tsx:113` renderiza `<p>{instance.name} ({instance.id})</p>`. La prueba y el diseño requieren que el nombre y el VMID se presenten como elementos claramente identificables en la fila (o celdas diferenciadas), permitiendo consultar `within(row).getByText('101')` y `screen.findByText('servidor-web')` de forma unívoca.
  2. **Columna de IP:** `Instances.tsx:125` tiene hardcodeado un guión fijo (`<td className="px-4 py-4">—</td>`). `instanceService.ts` debe leer `instance.ip` de la respuesta, y `Instances.tsx` debe renderizar la IP (ej. `192.168.1.50`) o el texto `"No detectada"` cuando sea `null`, junto a un botón interactivo para copiar la dirección al portapapeles (`navigator.clipboard.writeText`).
  3. **Badges de estado:** `Instances.tsx:119` renderiza `{instance.status}` como texto simple sin estilos distintivos. El criterio de aceptación exige badges diferenciados con estilos visuales estándar: verde para `running` / `en ejecución` y gris para `stopped` / `detenida`.
  4. **Botonera de acciones:** cada fila debe contar con su botonera de acciones presente en la tabla, independientemente de si la fila tiene IP o si las acciones operativas están condicionadas.
- **Entregable:**
  1. En `instanceService.ts`, agregar el campo opcional `ip?: string | null` en `InventoryInstance` y leerlo en el mapeo de `fetchInstanceInventory`.
  2. En `Instances.tsx`, ajustar el renderizado de la columna de nombre para que el texto del nombre (`instance.name`) y el VMID (`instance.id`) estén en elementos o nodos de texto separados.
  3. Renderizar la columna IP mostrando `instance.ip` si existe (con botón de copia con icono y `aria-label="Copiar IP"`) o `"No detectada"` si es `null`.
  4. Envolver el estado en un badge con clases de Tailwind que apliquen fondo y texto verde para `running` (ej. `bg-green-100 text-green-700` o variante shadcn correspondiente) y gris para `stopped`.
  5. Asegurar que la columna tipo exponga claramente `VM` o `LXC`.
- **Criterio de éxito:**
  - Los 5 casos de `test/front/instances-table.test.tsx` pasan 100% en verde.

> [!NOTE]
> **Actualización (06/10/2026):** En el frontend `5a86dce` (PR #83 y #84 de Luz / Cristian, commits `04987e0`, `189ab69`, `98deff1`, `f69c468`, `d9fcdbb`), se implementaron las correcciones de formato de nombre/VMID, IP con botón de copia, badges de estado y estilos de tabla. La suite `test/front/instances-table.test.tsx` pasa **5/5 (100% en verde)**.

---

## 🛠️ Fixes detectados en la verificación del 06/10/2026

### `FIX-46` - Alinear especificación de permisos de `GET /api/node/status` en contrato (`BAC-29`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro (autor de `contrato-etapa1.md`)
- **Estimación:** 0,5 h
- **Depende de:** `BAC-29` (en `terminado-1.md`).
- **Problema y evidencia:**
  En `backend/docs/contrato-etapa1.md`, la descripción de `GET /api/node/status` y los permisos de acceso (`ADMIN` u `OPERATOR`) se encuentran redactados en secciones separadas. La prueba de aceptación `test/back/etapa1_acceptance_test.go:51` busca validar mediante expresión regular que la sección del endpoint declare explícitamente en su propio bloque o tabla que cualquier usuario autenticado (`ADMIN` u `OPERATOR`) tiene acceso permitido (D1).
- **Entregable:**
  1. En `backend/docs/contrato-etapa1.md`, incluir en la sección de `GET /api/node/status` la indicación explícita de autenticación y roles permitidos: "Permisos: cualquier usuario autenticado (`ADMIN` u `OPERATOR`)".
- **Criterio de éxito:**
  - El caso `TestEtapa1/BAC-29_contrato_etapa1_incluye_node_status_y_codigos_de_error` pasa 100% en verde.

### `FIX-47` - Corrección de ciclo de recuperación y respuesta tras caída de Proxmox en telemetría (`BAC-22`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra (autora de `13f9c35`)
- **Estimación:** 1.0 h
- **Depende de:** `BAC-22` (en `terminado-1.md`).
- **Problema y evidencia:**
  En `TestEtapa1/BAC-22_telemetria_de_nodo_con_cache_redis_y_stale_reading` (`etapa1_acceptance_test.go:89`), tras simular la caída de Proxmox (donde el servicio responde con datos `stale` o error `502`), al recuperar el simulador/cliente y vencer el TTL, una nueva consulta a `GET /api/node/status` responde con `502 PROXMOX_UNAVAILABLE` en lugar de refrescar exitosamente los datos desde Proxmox y responder `200 OK`.
- **Entregable:**
  1. En `telemetria_service.go` / `node_handler.go`, asegurar que al recuperarse la conectividad con Proxmox, las lecturas posteriores actualicen la caché de Redis y retornen la respuesta viva con código `200 OK` y `stale: false`.
- **Criterio de éxito:**
  - El caso `TestEtapa1/BAC-22_telemetria_de_nodo_con_cache_redis_y_stale_reading` pasa 100% en verde.

### `FIX-48` - Literales de rutas de red en pruebas unitarias del adaptador de inventario (`BAC-23A`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro (autor de `d1dec4f`)
- **Estimación:** 0,5 h
- **Depende de:** `BAC-23A` (en `terminado-1.md`).
- **Problema y evidencia:**
  En `test/back/etapa1_acceptance_test.go:156`, la prueba verifica que las pruebas unitarias del backend invoquen o contemplen explícitamente los endpoints de Proxmox para obtención de interfaces: `qemu/{vmid}/agent/network-get-interfaces` y `lxc/{vmid}/interfaces`. En la suite unitaria actual, los mocks utilizan métodos abstractos o interfaces mockeadas donde no figuran las cadenas literales requeridas.
- **Entregable:**
  1. En los tests unitarios de `inventario_service_test.go` o del cliente Proxmox, incorporar casos o asserts donde se verifiquen explícitamente las rutas `network-get-interfaces` y `/interfaces`.
- **Criterio de éxito:**
  - El caso `TestEtapa1/BAC-23A_adaptador_de_inventario_resuelve_IP_de_VM_y_LXC_en_paralelo` pasa 100% en verde.

### `FIX-49` - Suscripción selectiva, deduplicación de eventos y conexión SSE compartida (`FRN-17A`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian (autor de PR #80/#81/#82)
- **Estimación:** 1.5 h
- **Depende de:** `FRN-17A` (en `terminado-1.md`).
- **Problema y evidencia:**
  La implementación de `useEvents.ts` y `EventsContext.tsx` no cumple con tres requisitos esenciales evaluados en `test/front/events-client.test.tsx`:
  1. No deduplica eventos por `id`: si el backend retransmite o reenvía un evento con un `id` ya procesado, no se descarta.
  2. No comparte una única conexión: múltiples consumidores de `useEvents` generan múltiples instancias de `EventSource` hacia `/api/events` en vez de suscribirse a un canal único administrado centralmente por `EventsProvider`.
  3. No provee suscripción granular: los componentes necesitan poder suscribirse a eventos específicos (ej. por `tipo` como `TASK_FINISHED` o por `recursoId` / VMID).
- **Entregable:**
  1. En `EventsContext.tsx`, mantener el estado de IDs de eventos procesados en un `Set` para descartar mensajes repetidos.
  2. Centralizar la conexión SSE en el `EventsProvider`, de modo que todos los llamados a `useEvents()` consuman del mismo stream.
  3. Exponer métodos de suscripción con limpieza al desmontar (ej. `subscribe(filtro, callback)`).
- **Criterio de éxito:**
  - Los casos de deduplicación, conexión única y suscripción por instancia en `test/front/events-client.test.tsx` pasan en verde.

