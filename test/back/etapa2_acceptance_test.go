package back_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

// TestEtapa2BAC34_AdaptadorMetricasRRD verifica los entregables de BAC-34 (Etapa 2):
// Adaptador de Proxmox RRDdata en el cliente y servicio de telemetría de instancia.
// Entregables:
//  1. Método en el cliente Proxmox ObtenerMetricasRRD(ctx, nodo, tipo, vmid, timeframe).
//  2. DTO normalizado PuntoMetricaDTO: { timestamp, cpuPercent, ramUsedBytes, ramTotalBytes, netInBps, netOutBps, diskReadBps, diskWriteBps }.
//  3. Pruebas unitarias con mocks en internal/adapters/secondary/proxmox/client_test.go.
// Criterio de éxito: consulta exitosa de métricas para VM y LXC devolviendo series temporales con valores numéricos coherentes sin fallos de conversión.
func TestEtapa2BAC34_AdaptadorMetricasRRD(t *testing.T) {
	// 1. Verificación del DTO normalizado en los puertos o modelos del dominio:
	portsCode := string(readFile(t, sourcePath("backend", "internal/core/ports/proxmox_port.go")))
	clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))
	clientTestCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client_test.go")))

	camposRequeridos := []string{
		"timestamp",
		"cpuPercent",
		"ramUsedBytes",
		"ramTotalBytes",
		"netInBps",
		"netOutBps",
		"diskReadBps",
		"diskWriteBps",
	}

	hasPuntoMetricaDTO := strings.Contains(portsCode, "PuntoMetricaDTO") || strings.Contains(clientCode, "PuntoMetricaDTO")
	if !hasPuntoMetricaDTO {
		t.Errorf("BAC-34: no se encuentra declarado el DTO PuntoMetricaDTO en internal/core/ports/proxmox_port.go ni en proxmox/client.go")
	} else {
		for _, campo := range camposRequeridos {
			if !regexp.MustCompile(`(?i)`+campo).MatchString(portsCode) && !regexp.MustCompile(`(?i)`+campo).MatchString(clientCode) {
				t.Errorf("BAC-34: PuntoMetricaDTO no incluye el campo requerido %q", campo)
			}
		}
	}

	// 2. Método en el cliente Proxmox: ObtenerMetricasRRD
	if !regexp.MustCompile(`(?i)ObtenerMetricasRRD\s*\(`).MatchString(portsCode) && !regexp.MustCompile(`(?i)ObtenerMetricasRRD\s*\(`).MatchString(clientCode) {
		t.Errorf("BAC-34: el cliente Proxmox o su puerto no define el método ObtenerMetricasRRD(ctx, nodo, tipo, vmid, timeframe)")
	}

	// 3. Pruebas unitarias con mocks en client_test.go:
	if !strings.Contains(clientTestCode, "ObtenerMetricasRRD") && !strings.Contains(clientTestCode, "TestClient_ObtenerMetricasRRD") {
		t.Errorf("BAC-34: internal/adapters/secondary/proxmox/client_test.go carece de pruebas unitarias para ObtenerMetricasRRD")
	}

	// 4. Verificación de ruta RRD en el cliente: debe consultar el endpoint /rrddata de Proxmox VE
	if !regexp.MustCompile(`(?i)/rrddata`).MatchString(clientCode) {
		t.Errorf("BAC-34: client.go no consulta el endpoint /rrddata de Proxmox VE")
	}
}

// TestEtapa2BAC35_CanalWebSocketMetricasConTicket verifica los entregables de BAC-35 (Etapa 2):
// Canal WebSocket de métricas en tiempo real con ticket efímero.
// Entregables:
//  1. Endpoint de emisión de ticket efímero POST /api/auth/ws-ticket (TTL 30 s en Redis).
//  2. Handler WebSocket /api/ws/instances/:vmid/metrics?ticket=... que valide el ticket y permisos.
//  3. Rechazo 401 Unauthorized sin ticket o con ticket expirado antes del handshake.
func TestEtapa2BAC35_CanalWebSocketMetricasConTicket(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	// 1. Endpoint POST /api/auth/ws-ticket:
	// Sin autenticación debe responder 401 Unauthorized
	unauthStatus, _ := requestValue(t, http.MethodPost, "/auth/ws-ticket", "", nil)
	if unauthStatus != http.StatusUnauthorized {
		t.Errorf("BAC-35: POST /api/auth/ws-ticket sin autenticar debe responder 401 Unauthorized, recibido %d", unauthStatus)
	}

	// Con autenticación válida debe emitir un ticket efímero
	statusTicket, bodyTicket := requestJSON(t, http.MethodPost, "/auth/ws-ticket", adminToken, nil)
	if statusTicket != http.StatusOK {
		t.Errorf("BAC-35: POST /api/auth/ws-ticket con ADMIN: esperado 200 OK, recibido %d: %#v", statusTicket, bodyTicket)
	} else {
		ticket, ok := bodyTicket["ticket"].(string)
		if !ok || strings.TrimSpace(ticket) == "" {
			t.Errorf("BAC-35: POST /api/auth/ws-ticket no devolvió el campo 'ticket': %#v", bodyTicket)
		}
	}

	ref101 := instanceRefByVMID(t, adminToken, 101)

	// 2. Conexión WebSocket /api/ws/instances/:instanceRef/metrics sin ticket o con ticket inválido:
	// Debe ser rechazada con 401 antes del handshake
	statusSinTicket, _ := requestValue(t, http.MethodGet, "/ws/instances/"+ref101+"/metrics", "", nil)
	if statusSinTicket != http.StatusUnauthorized {
		t.Errorf("BAC-35: GET /api/ws/instances/%s/metrics sin ticket debe responder 401 Unauthorized, recibido %d", ref101, statusSinTicket)
	}

	statusTicketInvalido, _ := requestValue(t, http.MethodGet, "/ws/instances/"+ref101+"/metrics?ticket=invalido-12345", "", nil)
	if statusTicketInvalido != http.StatusUnauthorized {
		t.Errorf("BAC-35: GET /api/ws/instances/%s/metrics con ticket inválido debe responder 401 Unauthorized, recibido %d", ref101, statusTicketInvalido)
	}

	// 3. Verificación estática del handler WebSocket en el código fuente:
	routerCode := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))
	instanceHandlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/instance_handler.go")))
	authHandlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/auth_handler.go")))

	hasWSRoute := strings.Contains(routerCode, "/ws/instances/:vmid/metrics") || strings.Contains(routerCode, "/ws/instances/:instanceRef/metrics") || strings.Contains(instanceHandlerCode, "/ws/instances/:vmid/metrics")
	hasTicketEndpoint := strings.Contains(routerCode, "/auth/ws-ticket") || strings.Contains(authHandlerCode, "EmitirTicketWS") || strings.Contains(authHandlerCode, "ws-ticket")

	if !hasWSRoute {
		t.Errorf("BAC-35: no está declarada la ruta /api/ws/instances/:vmid/metrics en los routers ni handlers")
	}
	if !hasTicketEndpoint {
		t.Errorf("BAC-35: no está declarado el endpoint /api/auth/ws-ticket en los routers ni handlers")
	}
}

// TestEtapa2INF10_DisponibilidadPlantillasEIsos verifica los entregables de INF-10 (Etapa 2):
// Disponibilidad y verificación de plantillas LXC (vztmpl) e ISOs (iso) a través de la API.
// Criterio de éxito: el backend lista plantillas e ISOs reales a través de GET /api/node/templates sin errores de permisos.
func TestEtapa2INF10_DisponibilidadPlantillasEIsos(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	status, body := requestValue(t, http.MethodGet, "/node/templates", adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("INF-10: GET /api/node/templates con ADMIN: esperado 200 OK, recibido %d: %#v", status, body)
	}

	// Normalizamos la lista de plantillas devuelta
	var templates []map[string]any
	switch val := body.(type) {
	case []any:
		for _, item := range val {
			if m, ok := item.(map[string]any); ok {
				templates = append(templates, m)
			}
		}
	case map[string]any:
		if items, ok := val["items"].([]any); ok {
			for _, item := range items {
				if m, ok := item.(map[string]any); ok {
					templates = append(templates, m)
				}
			}
		} else if data, ok := val["data"].([]any); ok {
			for _, item := range data {
				if m, ok := item.(map[string]any); ok {
					templates = append(templates, m)
				}
			}
		}
	}

	if len(templates) == 0 {
		t.Fatalf("INF-10: GET /api/node/templates devolvió una lista vacía de plantillas: %#v", body)
	}

	// Criterio: debe existir al menos una plantilla LXC (vztmpl o LXC_TEMPLATE) y una ISO (iso o ISO_IMAGE)
	tieneLXC := false
	tieneISO := false

	for _, tmpl := range templates {
		tipo, _ := tmpl["type"].(string)
		formato, _ := tmpl["format"].(string)
		id, _ := tmpl["id"].(string)
		nombre, _ := tmpl["name"].(string)

		cadena := strings.ToLower(fmt.Sprintf("%s %s %s %s", tipo, formato, id, nombre))
		if strings.Contains(cadena, "vztmpl") || strings.Contains(cadena, "lxc") || strings.Contains(cadena, "tar.zst") || strings.Contains(cadena, "tar.xz") {
			tieneLXC = true
		}
		if strings.Contains(cadena, "iso") {
			tieneISO = true
		}
	}

	if !tieneLXC {
		t.Errorf("INF-10: no se encontró ninguna plantilla LXC (vztmpl / debian / alpine) en /api/node/templates: %#v", templates)
	}
	if !tieneISO {
		t.Errorf("INF-10: no se encontró ninguna imagen ISO en /api/node/templates: %#v", templates)
	}
}

// TestEtapa2BAC36_MetricasHistoricasEndpoint verifica los entregables de BAC-36 (Etapa 2):
// Endpoint REST de métricas históricas de instancia (timeframe=2h|24h).
// Entregables:
//  1. Endpoint GET /api/instances/:vmid/metrics?timeframe=2h|24h protegido con RequireAuth y RequireInstanceAccess.
//  2. Parámetro timeframe configurable (por defecto 2h, soportando 24h), mapeado a los intervalos de Proxmox (hour y day).
//  3. Caché corta en Redis (TTL 30 s) para evitar consultas repetitivas al hipervisor.
// Criterio de éxito: responde con el arreglo histórico de métricas formateado para el consumo de gráficos.
func TestEtapa2BAC36_MetricasHistoricasEndpoint(t *testing.T) {
	mainCode := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))
	instanceHandlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/instance_handler.go")))

	hasMetricsRoute := strings.Contains(mainCode, "/instances/:vmid/metrics") || strings.Contains(mainCode, "/instances/:id/metrics") ||
		strings.Contains(mainCode, "/instances/:instanceRef/metrics") || strings.Contains(mainCode, "/metrics") ||
		strings.Contains(instanceHandlerCode, "/instances/:vmid/metrics") || strings.Contains(instanceHandlerCode, "ObtenerMetricas")
	if !hasMetricsRoute {
		t.Errorf("BAC-36: la ruta GET /api/instances/:vmid/metrics no está declarada en cmd/api/main.go ni en instance_handler.go")
	}

	if testing.Short() || setupErr != nil {
		return
	}

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	ref101 := instanceRefByVMID(t, adminToken, 101)

	unauthStatus, _ := requestValue(t, http.MethodGet, "/instances/"+ref101+"/metrics", "", nil)
	if unauthStatus != http.StatusUnauthorized {
		t.Errorf("BAC-36: GET /api/instances/%s/metrics sin autenticar debe responder 401 Unauthorized, recibido %d", ref101, unauthStatus)
	}

	status2h, body2h := requestValue(t, http.MethodGet, "/instances/"+ref101+"/metrics?timeframe=2h", adminToken, nil)
	if status2h != http.StatusOK {
		t.Errorf("BAC-36: GET /api/instances/%s/metrics?timeframe=2h con ADMIN: esperado 200 OK, recibido %d: %#v", ref101, status2h, body2h)
	}

	status24h, body24h := requestValue(t, http.MethodGet, "/instances/"+ref101+"/metrics?timeframe=24h", adminToken, nil)
	if status24h != http.StatusOK {
		t.Errorf("BAC-36: GET /api/instances/%s/metrics?timeframe=24h con ADMIN: esperado 200 OK, recibido %d: %#v", ref101, status24h, body24h)
	}

	statusInv, _ := requestValue(t, http.MethodGet, "/instances/"+ref101+"/metrics?timeframe=invalido", adminToken, nil)
	if statusInv != http.StatusBadRequest {
		t.Errorf("BAC-36: GET /api/instances/%s/metrics con timeframe inválido debe responder 400 Bad Request, recibido %d", ref101, statusInv)
	}
}

// TestEtapa2BAC39_SeguimientoCreacionYPermisos verifica los entregables de BAC-39 (Etapa 2):
// Seguimiento de creación, autoasignación de permisos y eventos (RF-07, RF-11).
// Entregables:
//  1. Integración en seguimiento_tareas.go: al detectar fin exitoso (exitstatus: "OK") de un UPID de creación,
//     insertar atómicamente (creador_id, vmid, FULL_ACCESS) en permisos_instancia.
//  2. Emisión por el canal SSE /api/events de TASK_FINISHED con accion: "create", estado: "COMPLETED".
//  3. Si la tarea de creación falla en Proxmox, liberar bloqueos y registrar el evento con estado: "FAILED".
func TestEtapa2BAC39_SeguimientoCreacionYPermisos(t *testing.T) {
	seguimientoCode := string(readFile(t, sourcePath("backend", "internal/core/services/seguimiento_tareas.go")))
	tareaRepoCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/postgres/tarea_repository.go")))

	hasCreacionHandler := strings.Contains(seguimientoCode, "create") || strings.Contains(seguimientoCode, "qmcreate") || strings.Contains(seguimientoCode, "vzcreate") ||
		strings.Contains(tareaRepoCode, "debeAsignarPermisoCreacion") || strings.Contains(tareaRepoCode, "esCreacion")
	if !hasCreacionHandler {
		t.Errorf("BAC-39: no se contempla el manejo de tareas de creación (create / qmcreate / vzcreate)")
	}

	hasAutoPermiso := strings.Contains(tareaRepoCode, "permisos_instancia") || strings.Contains(tareaRepoCode, "PermisoInstancia") ||
		strings.Contains(tareaRepoCode, "NivelAccesoFullAccess") || strings.Contains(tareaRepoCode, "FULL_ACCESS") ||
		strings.Contains(seguimientoCode, "permisos_instancia") || strings.Contains(seguimientoCode, "FULL_ACCESS")
	if !hasAutoPermiso {
		t.Errorf("BAC-39: no se implementa la autoasignación atómica de FULL_ACCESS en permisos_instancia al creador")
	}

	hasTaskFinishedEvent := strings.Contains(seguimientoCode, "TASK_FINISHED")
	if !hasTaskFinishedEvent {
		t.Errorf("BAC-39: seguimiento_tareas.go no emite el evento TASK_FINISHED")
	}
}

// TestEtapa2BAC40_EndpointConfigRecursos verifica los entregables de BAC-40 (Etapa 2):
// Endpoint PUT /api/instances/:vmid/config con cálculo diferencial (ΔRAM) y detección de hot-plug.
// Entregables:
//  1. Endpoint PUT /api/instances/:vmid/config accesible para usuarios con rol ADMIN o FULL_ACCESS.
//  2. Payload: { cores?: number, ramMb?: number }.
//  3. Cálculo diferencial: si ΔRAM > 0, validar contra la cuota libre con CuotasService.
//  4. Invocación de PUT /nodes/{node}/{tipo}/{vmid}/config en Proxmox.
//  5. Detección de hot-plug: responde { requiresReboot: bool, message?: string }. Rechazo 409 RESOURCE_QUOTA_EXCEEDED ante exceso.
func TestEtapa2BAC40_EndpointConfigRecursos(t *testing.T) {
	mainCode := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))
	configHandlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/instance_config_handler.go")))

	hasConfigRoute := (strings.Contains(mainCode, "/instances") && strings.Contains(mainCode, "/:vmid/config")) ||
		strings.Contains(mainCode, "/instances/:vmid/config") ||
		strings.Contains(configHandlerCode, "ActualizarConfig")
	if !hasConfigRoute {
		t.Errorf("BAC-40: la ruta PUT /api/instances/:vmid/config no está declarada en main.go ni en instance_config_handler.go")
	}

	if testing.Short() || setupErr != nil {
		return
	}

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	ref101 := instanceRefByVMID(t, adminToken, 101)

	unauthStatus, _ := requestValue(t, http.MethodPut, "/instances/"+ref101+"/config", "", map[string]any{"cores": 2})
	if unauthStatus != http.StatusUnauthorized {
		t.Errorf("BAC-40: PUT /api/instances/%s/config sin autenticar debe responder 401 Unauthorized, recibido %d", ref101, unauthStatus)
	}

	status, body := requestJSON(t, http.MethodPut, "/instances/"+ref101+"/config", adminToken, map[string]any{"cores": 2, "ramMb": 2048})
	if status != http.StatusOK && status != http.StatusAccepted {
		t.Errorf("BAC-40: PUT /api/instances/%s/config con ADMIN: esperado 200 OK o 202 Accepted, recibido %d: %#v", ref101, status, body)
	} else {
		if _, ok := body["requiresReboot"]; !ok {
			t.Errorf("BAC-40: respuesta no incluye el campo 'requiresReboot': %#v", body)
		}
	}
}

// TestEtapa2BAC41_AuditoriaAprovisionamientoYConfig verifica los entregables de BAC-41 (Etapa 2):
// Auditoría inmutable de aprovisionamiento y cambios de recursos (RF-08).
// Entregables:
//  1. Registro de CREATE_INSTANCE con estados PENDING y final SUCCESS o FAILED.
//  2. Registro de CONFIG_INSTANCE con valores anteriores y nuevos en detalles JSONB.
//  3. Disponibilidad en GET /api/admin/audit.
func TestEtapa2BAC41_AuditoriaAprovisionamientoYConfig(t *testing.T) {
	auditPortCode := string(readFile(t, sourcePath("backend", "internal/core/ports/audit_port.go")))

	hasCreateAction := strings.Contains(auditPortCode, "CREATE_INSTANCE") || strings.Contains(auditPortCode, "CREAR_VM")
	if !hasCreateAction {
		t.Errorf("BAC-41: audit_port.go no define la acción CREATE_INSTANCE o CREAR_VM")
	}

	hasConfigAction := strings.Contains(auditPortCode, "CONFIG_INSTANCE") || strings.Contains(auditPortCode, "CONFIGURAR_INSTANCIA")
	if !hasConfigAction {
		t.Errorf("BAC-41: audit_port.go no define la acción CONFIG_INSTANCE")
	}
}

// TestEtapa2BAC51_GuardaConcurrenciaInstanceBusy verifica los entregables de BAC-51 (Etapa 2):
// Guarda centralizada de concurrencia operativa y exclusión mutua (409 INSTANCE_BUSY).
// Entregables:
//  1. Middleware/guard RequireInstanceIdle: valida si la instancia tiene una tarea activa o lock en Redis.
//  2. Rechazo 409 Conflict (INSTANCE_BUSY) si la instancia está ocupada.
//  3. Integración en POST /api/instances/:id/status/:action y PUT /api/instances/:id/config.
func TestEtapa2BAC51_GuardaConcurrenciaInstanceBusy(t *testing.T) {
	mainCode := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))
	instanceHandlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/instance_handler.go")))

	hasGuard := strings.Contains(mainCode, "RequireInstanceIdle") || strings.Contains(instanceHandlerCode, "RequireInstanceIdle") ||
		strings.Contains(instanceHandlerCode, "INSTANCE_BUSY") || strings.Contains(instanceHandlerCode, "activeTask")
	if !hasGuard {
		t.Errorf("BAC-51: no se encuentra la guarda RequireInstanceIdle ni la validación centralizada de INSTANCE_BUSY en main.go o instance_handler.go")
	}
}

// TestEtapa2BAC52_BrokerSSEEventosAlertaGlobal verifica los entregables de BAC-52 (Etapa 2):
// Generalización del Broker SSE para Eventos Globales de Alerta (ALERT_TRIGGERED).
// Entregables:
//  1. Extensión del modelo Evento en Go para admitir recursoTipo: "HOST" | "VM" | "LXC", siendo recursoId opcional.
//  2. Regla de enrutamiento en Redis Pub/Sub: los eventos con recursoTipo == "HOST" se difunden a todos los clientes con rol ADMIN.
//  3. Actualización de docs/contrato-eventos.md.
func TestEtapa2BAC52_BrokerSSEEventosAlertaGlobal(t *testing.T) {
	eventPortCode := string(readFile(t, sourcePath("backend", "internal/core/ports/event_port.go")))

	hasAlertTriggered := strings.Contains(eventPortCode, "ALERT_TRIGGERED")
	if !hasAlertTriggered {
		t.Errorf("BAC-52: event_port.go no incluye el tipo de evento ALERT_TRIGGERED")
	}

	hasHostResource := strings.Contains(eventPortCode, `"HOST"`) || strings.Contains(eventPortCode, "RecursoHost")
	if !hasHostResource {
		t.Errorf("BAC-52: event_port.go no incluye el tipo de recurso HOST para alertas globales")
	}

	docCode := string(readFile(t, sourcePath("backend", "docs/contrato-eventos.md")))
	if !strings.Contains(docCode, "ALERT_TRIGGERED") && !strings.Contains(docCode, "HOST") {
		t.Errorf("BAC-52: docs/contrato-eventos.md no documenta los eventos de alerta ALERT_TRIGGERED a nivel HOST")
	}
}

// TestEtapa2BAC53_ParticionadoPreventivoEIndexacionAuditoria verifica los entregables de BAC-53 (Etapa 2):
// Particionado Preventivo e Indexación Compuesta de la Tabla auditoria.
// Entregables:
//  1. Script de migración SQL DDL que crea preventivamente las particiones trimestrales: auditoria_2026_q3, auditoria_2026_q4, auditoria_2027_q1.
//  2. Índices compuestos para partition pruning: CREATE INDEX idx_auditoria_query ON auditoria (fecha_hora DESC, usuario_id, accion).
func TestEtapa2BAC53_ParticionadoPreventivoEIndexacionAuditoria(t *testing.T) {
	dbCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/postgres/db.go")))

	particionesEsperadas := []string{"auditoria_2026_q3", "auditoria_2026_q4", "auditoria_2027_q1"}
	for _, p := range particionesEsperadas {
		if !strings.Contains(dbCode, p) {
			t.Errorf("BAC-53: db.go no declara preventivamente la partición trimestral %q", p)
		}
	}

	hasQueryIndex := strings.Contains(dbCode, "idx_auditoria_query") ||
		(strings.Contains(dbCode, "fecha_hora") && strings.Contains(dbCode, "usuario_id") && strings.Contains(dbCode, "accion"))
	if !hasQueryIndex {
		t.Errorf("BAC-53: db.go no define el índice compuesto idx_auditoria_query ON auditoria (fecha_hora DESC, usuario_id, accion)")
	}
}


