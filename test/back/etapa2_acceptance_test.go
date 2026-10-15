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

	// 2. Conexión WebSocket /api/ws/instances/:vmid/metrics sin ticket o con ticket inválido:
	// Debe ser rechazada con 401 antes del handshake
	statusSinTicket, _ := requestValue(t, http.MethodGet, "/ws/instances/101/metrics", "", nil)
	if statusSinTicket != http.StatusUnauthorized {
		t.Errorf("BAC-35: GET /api/ws/instances/101/metrics sin ticket debe responder 401 Unauthorized, recibido %d", statusSinTicket)
	}

	statusTicketInvalido, _ := requestValue(t, http.MethodGet, "/ws/instances/101/metrics?ticket=invalido-12345", "", nil)
	if statusTicketInvalido != http.StatusUnauthorized {
		t.Errorf("BAC-35: GET /api/ws/instances/101/metrics con ticket inválido debe responder 401 Unauthorized, recibido %d", statusTicketInvalido)
	}

	// 3. Verificación estática del handler WebSocket en el código fuente:
	routerCode := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))
	instanceHandlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/instance_handler.go")))
	authHandlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/auth_handler.go")))

	hasWSRoute := strings.Contains(routerCode, "/ws/instances/:vmid/metrics") || strings.Contains(instanceHandlerCode, "/ws/instances/:vmid/metrics")
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

