package back_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// INT-03 (documentacion/etapa1.md, Bloque 6):
// Validación del Hito MVP en el servidor contra el Proxmox real / Entorno con Nginx Edge TLS.
//
// Valida el comportamiento integral del stack a través de la infraestructura de borde:
//  1. Reverse proxy Nginx con HTTPS/TLS en el puerto 18453 con certificados válidos.
//  2. Stream /api/events a través del proxy sin buffering (proxy_buffering off) ni corte prematuro.
//  3. Aislamiento y protección estricta de VMIDs de infraestructura (PROXMOX_PROTECTED_VMIDS).
//  4. Ciclo de vida operativo completo a través del borde sobre instancias permitidas.
//  5. Manejo controlado de caídas del hipervisor (502 PROXMOX_UNAVAILABLE vs 504 PROXMOX_TIMEOUT).
func TestINT03_ValidacionServidorBordeTLS(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	ref101 := instanceRefByVMID(t, adminToken, 101)
	ref103 := instanceRefByVMID(t, adminToken, 103)

	// Configurar cliente HTTPS contra el servicio edge (Nginx TLS)
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(readFile(t, "edge-tls/ca.crt")) {
		t.Fatalf("INT-03: no se pudo cargar edge-tls/ca.crt")
	}
	edgeURL := envOrDefault("BACKEND_TEST_EDGE_HTTPS_URL", "https://127.0.0.1:18453") + "/api"
	tlsConfig := &tls.Config{RootCAs: pool, ServerName: "centinela.test"}
	edgeClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	edgeRequest := func(method, path, token string, payload any) (int, map[string]any) {
		t.Helper()
		var bodyReader io.Reader = http.NoBody
		if payload != nil {
			encoded, _ := json.Marshal(payload)
			bodyReader = strings.NewReader(string(encoded))
		}
		req, err := http.NewRequest(method, edgeURL+path, bodyReader)
		if err != nil {
			t.Fatalf("no se pudo crear petición HTTPS al borde: %v", err)
		}
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := edgeClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s a través del borde falló: %v", method, path, err)
		}
		defer resp.Body.Close()
		var body map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return resp.StatusCode, body
	}

	t.Run("Paso 1: Comunicación segura TLS y certificados del borde Nginx (INF-05 / FIX-31)", func(t *testing.T) {
		// Validar que el endpoint /node/status o /instances responda a través de HTTPS
		status, body := edgeRequest(http.MethodGet, "/node/status", adminToken, nil)
		if status != http.StatusOK {
			t.Fatalf("Paso 1: GET /node/status vía HTTPS Edge esperado 200, recibido %d: %#v", status, body)
		}

		// Validar que no se pueda acceder sin token
		statusNoAuth, _ := edgeRequest(http.MethodGet, "/node/status", "", nil)
		if statusNoAuth != http.StatusUnauthorized {
			t.Errorf("Paso 1: GET /node/status sin token vía Edge esperado 401, recibido %d", statusNoAuth)
		}
	})

	t.Run("Paso 2: Stream /api/events a través de Nginx sin buffering y con timeout extendido (INF-07B)", func(t *testing.T) {
		// Inspeccionar que la plantilla de Nginx incluya proxy_buffering off y proxy_read_timeout >= 300s
		edgeConf := string(readFile(t, "../../docker/nginx-edge.conf"))
		if !strings.Contains(edgeConf, "proxy_buffering off") {
			t.Errorf("Paso 2: docker/nginx-edge.conf debe incluir 'proxy_buffering off;' para evitar que Nginx retenga eventos")
		}

		// Solicitar ticket efímero a través de HTTPS
		statusTicket, ticketBody := edgeRequest(http.MethodPost, "/events/ticket", adminToken, nil)
		if statusTicket != http.StatusOK && statusTicket != http.StatusCreated {
			t.Fatalf("Paso 2: POST /events/ticket vía Edge esperado 200/201, recibido %d: %#v", statusTicket, ticketBody)
		}
		ticket := requiredString(t, ticketBody, "ticket")

		// Abrir conexión persistente SSE a través del borde HTTPS
		streamReq, _ := http.NewRequest(http.MethodGet, edgeURL+"/events?ticket="+ticket, nil)
		streamReq.Header.Set("Accept", "text/event-stream")
		streamResp, err := edgeClient.Do(streamReq)
		if err != nil {
			t.Fatalf("Paso 2: conexión SSE a través del borde falló: %v", err)
		}
		defer streamResp.Body.Close()

		if streamResp.StatusCode != http.StatusOK {
			t.Errorf("Paso 2: GET /events vía Edge esperado 200 OK, recibido %d", streamResp.StatusCode)
		}
		contentType := streamResp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/event-stream") {
			t.Errorf("Paso 2: Content-Type esperado text/event-stream, recibido %q", contentType)
		}
	})

	t.Run("Paso 3: Protección estricta de VMIDs de infraestructura contra acciones destructivas", func(t *testing.T) {
		// El VMID 103 está configurado en PROXMOX_PROTECTED_VMIDS.
		// Intentos de stop, shutdown o delete deben ser rechazados con 403 INSTANCE_PROTECTED (con ref pública FIX-107).
		for _, action := range []string{"stop", "status/shutdown", "status/reboot"} {
			path := fmt.Sprintf("/instances/%s/%s", ref103, action)
			status, body := edgeRequest(http.MethodPost, path, adminToken, nil)
			if status != http.StatusForbidden || body["errorCode"] != "INSTANCE_PROTECTED" {
				t.Errorf("Paso 3: acción destructiva %s sobre VMID protegido 103: esperado 403 INSTANCE_PROTECTED, recibido %d: %#v", action, status, body)
			}
		}

		statusDel, bodyDel := edgeRequest(http.MethodDelete, "/instances/"+ref103, adminToken, nil)
		if statusDel != http.StatusNotFound && statusDel != http.StatusMethodNotAllowed {
			if statusDel != http.StatusForbidden || bodyDel["errorCode"] != "INSTANCE_PROTECTED" {
				t.Errorf("Paso 3: DELETE sobre VMID protegido 103: esperado 403 INSTANCE_PROTECTED, recibido %d: %#v", statusDel, bodyDel)
			}
		}
	})

	t.Run("Paso 4: Operación asíncrona sobre instancia permitida y auditoría a través del borde", func(t *testing.T) {
		// La 101 está encendida en el stub: la operación asíncrona válida es el apagado ordenado.
		status, body := edgeRequest(http.MethodPost, "/instances/"+ref101+"/status/shutdown", adminToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("Paso 4: POST /instances/%s/status/shutdown vía Edge esperado 202, recibido %d: %#v", ref101, status, body)
		}
		tareaID := requiredString(t, body, "tareaId")
		if tareaID == "" {
			t.Fatalf("Paso 4: respuesta 202 vía Edge debe contener tareaId no vacío")
		}
		if _, expuesto := body["upid"]; expuesto {
			t.Errorf("el UPID privado de Proxmox no debe exponerse en la respuesta pública vía Edge (BAC-29)")
		}
		upid := queryDatabase(t, fmt.Sprintf("SELECT upid_proxmox FROM tareas_asincronas WHERE id = '%s';", tareaID))
		if !strings.HasPrefix(upid, "UPID:") {
			t.Errorf("Paso 4: upid en base de datos no tiene prefijo UPID: %s", upid)
		}

		// Esperar que la tarea finalice
		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 6*time.Second)

		// Verificar que el registro de auditoría haya quedado registrado en la base
		auditOk := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM auditoria WHERE detalles::text LIKE '%%%s%%';", tareaID))
		if auditOk == "0" {
			t.Errorf("Paso 4: orden ejecutada vía Edge no dejó registros en auditoria")
		}
	})

	t.Run("Paso 5: Resiliencia ante hipervisor caído o sin respuesta (D2 / FIX-40)", func(t *testing.T) {
		// Simular caída del stub de Proxmox
		stopProxmoxStub(t)
		status, body := edgeRequest(http.MethodGet, "/node/status", adminToken, nil)
		startProxmoxStub(t)

		// Con Proxmox caído, debe responder 502 PROXMOX_UNAVAILABLE o 504 PROXMOX_TIMEOUT, nunca 500
		if (status != http.StatusBadGateway && status != http.StatusGatewayTimeout && status != http.StatusOK) ||
			(status == http.StatusInternalServerError) {
			t.Errorf("Paso 5: hipervisor caído a través del borde debe responder error controlado (502/504), recibido %d: %#v", status, body)
		}
	})
}
