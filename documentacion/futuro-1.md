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

### `FIX-62` - Heurística de circuit breaker y recuperación de telemetría en `nodo_service` (`BAC-22`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra / Lisandro
- **Estado:** Pendiente en el submódulo `backend`. La suite de tests de este repositorio ya adecuó la tolerancia de sincronización de recuperación, pero el servicio en el backend presenta una anomalía de diseño cuando no existe caché previa.
- **Estimación:** 1 h
- **Depende de:** `BAC-22` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En `backend/internal/core/services/nodo_service.go:24`, se definió `pausaTrasFallaNodo = 5 * time.Second` y en la línea 64:
  ```go
  if err := s.fallaReciente(); err != nil {
      return s.respaldo(ctx, err)
  }
  ```
  `s.respaldo` busca la clave `claveNodoUltimo` (`"node:status:last_known"`) en Redis. Si dicha clave no existe (por ejemplo, arranque en frío o caída de Proxmox antes de registrar lecturas exitosas), `s.respaldo` retorna `nil, false, errProxmox` (HTTP 502 inmediato). Esto genera que durante 5 segundos continuos cualquier petición sea rechazada a ciegas sin siquiera intentar consultar Proxmox, incluso si Proxmox ya se recuperó de inmediato. Un circuit breaker solo debe evitar llamadas si tiene datos *stale* que servir como fallback, o debe implementar un estado *half-open* que permita sondear la recuperación.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/core/services/nodo_service.go`, en la función `ObtenerEstado`, condicionar la llamada a `fallaReciente()`: si no existe lectura previa en `claveNodoUltimo`, **no bloquear las peticiones entrantes durante 5 segundos a ciegas**; permitir el intento hacia Proxmox o reducir la penalización sin caché.
  2. Implementar un mecanismo de sondeo o reintento inmediato tras la expiración del cooldown que limpie `ultimaFalla` y `errFalla` al obtener una respuesta exitosa (`200 OK`).
  3. Asegurar que `TestNodoService` en `nodo_service_test.go` cubra el escenario de recuperación tras fallo sin caché previa y con caché previa.
- **Criterio de éxito:**
  - El backend se recupera inmediatamente tras el restablecimiento del hipervisor sin quedar en un blackout de 5 segundos cuando no hay caché previa disponible.
  - Pasan las pruebas unitarias de `nodo_service_test.go` y la prueba de integración `TestEtapa1/BAC-22`.

### `FIX-63` - Cobertura de pruebas unitarias para `ObtenerInterfaces` en el cliente Proxmox (`BAC-23A`) (Backend)

- **Área:** Backend / pruebas
- **Asignado:** Lisandro
- **Estado:** Pendiente en el submódulo `backend`.
- **Estimación:** 1 h
- **Depende de:** `BAC-23A` (en `terminado-1.md`).
- **Problema técnico detallado:**
  El método `ObtenerInterfaces(ctx context.Context, nodo, tipo string, vmid int)` en `backend/internal/adapters/secondary/proxmox/client.go:474` implementa el consumo HTTP de:
  - QEMU VMs: `/api2/json/nodes/{node}/qemu/{vmid}/agent/network-get-interfaces`
  - LXC Containers: `/api2/json/nodes/{node}/lxc/{vmid}/interfaces`
  Sin embargo, en `backend/internal/adapters/secondary/proxmox/client_test.go` **no se escribió ninguna prueba unitaria para `ObtenerInterfaces`**. La suite de aceptación de `BAC-23A` (`etapa1_acceptance_test.go:171`) verifica que existan pruebas unitarias para ambas rutas y falla con:
  `no hay pruebas unitarias del adaptador de IP (ningún _test.go del backend menciona network-get-interfaces ni /interfaces)`.
  Aunque `inventario_service_test.go` prueba la lógica interna con stubs en memoria, el adaptador HTTP contra Proxmox quedó sin cobertura.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/adapters/secondary/proxmox/client_test.go`, agregar casos de prueba con `httptest.Server`:
     - `TestClient_ObtenerInterfaces_Qemu`: verifica que una VM consulte la ruta `/api2/json/nodes/{node}/qemu/{vmid}/agent/network-get-interfaces` y parsee correctamente las interfaces y direcciones IPv4/IPv6 devueltas por el Guest Agent (con `prefix` entero).
     - `TestClient_ObtenerInterfaces_LXC`: verifica que un contenedor consulte la ruta `/api2/json/nodes/{node}/lxc/{vmid}/interfaces` y parsee las interfaces y direcciones (soportando `prefix` string o número, y devolviendo vacío sin error si `data: null` por estar detenido).
     - `TestClient_ObtenerInterfaces_ErrorYTimeout`: verifica el manejo de errores (ej. VM sin guest agent respondiendo error o 500) devolviendo `ErrGuestAgentNoDisponible` o `nil` controlado según el contrato.
- **Criterio de éxito:**
  - `go test -v ./internal/adapters/secondary/proxmox/...` pasa al 100% ejecutando los nuevos tests.
  - La suite de aceptación `TestEtapa1/BAC-23A` encuentra las pruebas en `client_test.go`, las ejecuta y pasa 100% en verde.

### `FIX-65` - Explicitar en `docs/contrato-etapa1.md` el acceso de cualquier usuario autenticado a `/node/status` (D1) (`BAC-29` / `FIX-46`) (Backend)

- **Área:** Backend / documentación
- **Asignado:** Nico
- **Estado:** Pendiente en el submódulo `backend`. La suite de tests de este repo ya adecuó el matcher multilínea, pero el documento de contrato debe explicitar la especificación formal del requisito D1.
- **Estimación:** 0,5 h
- **Depende de:** `FIX-46` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En `backend/docs/contrato-etapa1.md`, la sección `GET /api/node/status` describe el formato JSON de la respuesta y la telemetría, pero no define con precisión formal en el encabezado del endpoint la matriz de control de acceso requerida por el entregable 5 de `BAC-29` (D1). El texto actual menciona de forma genérica en párrafos separados que ambas operaciones usan `Authorization: Bearer <accessToken>`, pero no explicita que `GET /api/node/status` es accesible tanto para el rol `ADMIN` como para el rol `OPERATOR` sin restricciones de permisos de instancia.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/docs/contrato-etapa1.md`, en la sección correspondiente a `GET /api/node/status`, agregar explícitamente:
     > **Control de acceso (D1):** Endpoint accesible para cualquier usuario autenticado (`Authorization: Bearer <accessToken>`), habilitado tanto para el rol `ADMIN` como para el rol `OPERATOR`. No requiere permisos específicos sobre instancias.
  2. Mantener la alineación entre la versión Markdown y el archivo `.docx` correspondiente si aplica.
- **Criterio de éxito:**
  - `backend/docs/contrato-etapa1.md` especifica formalmente el requisito D1.
  - La prueba de contrato `TestEtapa1/BAC-29` valida el texto contractual y pasa 100% en verde.

### `FIX-69` - Restaurar definiciones SSE, esquemas de eventos y contratos en `swagger.yaml` / docs (`BAC-25A` / `BAC-29`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro / Lucas
- **Estado:** Pendiente en el submódulo `backend`. Correr `go test ./docs/...` en el backend falla con errores de esquemas faltantes.
- **Estimación:** 1,5 h
- **Depende de:** `BAC-25A` y `BAC-29` (en `terminado-1.md`).
- **Problema técnico detallado:**
  Al regenerar la documentación de Swagger con `swag init` en el commit `f06640f` del backend, se sobreescribieron `docs/swagger.yaml`, `docs/swagger.json` y `docs/docs.go`. Como los handlers no contenían anotaciones declarativas completas para los modelos de streaming SSE ni para el ticket efímero, se perdieron del archivo YAML/JSON las siguientes definiciones esenciales de la Etapa 1:
  - `http.SSEEventPayload`
  - `http.SSEDetallesEvento`
  - `http.TaskSuccess`
  - `http.TaskFailed`
  - `http.SSETicketResponse`
  - `http.SSECierrePayload`
  - `http.SSETiposPayload`
  Además, `contrato_etapa1_test.go` falla porque la respuesta 200 de `/node/status` quedó documentada con el tipo interno `#/definitions/internal_adapters_primary_http.EstadoNodeResponse` en lugar de la definición de contrato esperada `#/definitions/http.EstadoNodeResponse`.
- **Qué debe hacer el equipo de backend:**
  1. En los handlers de HTTP (`events_handler.go` y `node_handler.go`), incorporar o corregir las anotaciones de Swagger (`@Success`, `@Failure`, `@Produce`, `@Router` y tipos `@Model` o referencias a los DTOs de eventos).
  2. Asegurar que las estructuras de DTO de eventos SSE (`RealtimeEvent`, `TaskFinishedDetails`, `TicketResponse`, etc.) estén expuestas y referenciadas en las anotaciones de swag para que se generen en `docs/swagger.yaml` y `docs/swagger.json`.
  3. Ejecutar la regeneración (`swag init` con los flags de formato adecuados, o configurar alias de paquetes) para que los nombres de los esquemas en `#/definitions` no contengan prefijos de rutas internas (`internal_adapters_primary_http`).
  4. Ejecutar localmente `go test -v ./docs/...` en el backend hasta que todas las aserciones de artefactos (`TestContratoEtapa1ArtefactosJSONYAML`, `TestContratoEtapa1ArtefactosJSONGo` y `TestContratoEtapa1`) pasen en verde.
- **Criterio de éxito:**
  - `go test -v ./docs/...` en el repositorio backend pasa 100% en verde sin fallos de definiciones faltantes.
  - `docs/swagger.yaml` y `docs/swagger.json` contienen los esquemas completos de SSE y telemetría de nodo.


