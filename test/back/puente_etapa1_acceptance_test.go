package back_test

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Tareas puente entre la fase base y la Etapa 1 (documentacion/actual.md):
//   - BAC-21B: GET /api/instances extendido de forma retrocompatible; acciones de energía con FULL_ACCESS
//     (por /status/:action o por /start y /stop) registradas en tareas_asincronas y auditoria; DELETE
//     solo ADMIN con 409 si está encendida (cuando exista, BAC-24B).
//   - BAC-21C: tickets efímeros en Redis para /api/events y corte del stream al revocar la sesión.
//   - FIX-40 (decisión D2): 504 PROXMOX_TIMEOUT distinto de 502 PROXMOX_UNAVAILABLE (backend-lento en compose.yaml).
//   - INF-07B (BRG-03, etapa1.md): /api/events a través del Nginx del borde (docker/nginx-edge.conf)
//     sin buffer ni corte por inactividad, con un TASK_FINISHED real de punta a punta.
//
// Stub de Proxmox (proxmox-stub/nginx.conf): 101 qemu running, 102 lxc stopped, 103 qemu running
// (103 es VMID protegido por PROXMOX_PROTECTED_VMIDS en compose.yaml).
func TestPuenteEtapa1(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	suffix := time.Now().UnixNano()
	operatorID, _ := createUserThroughAPI(t, adminToken, fmt.Sprintf("op_puente_%d", suffix), fmt.Sprintf("op.puente.%d@elcentinela.com", suffix), "OPERATOR")
	operatorToken := signedAccessToken(t, operatorID, "OPERATOR", orgID)
	permEndpoint := "/admin/users/" + operatorID + "/permissions"
	if status, body := requestValue(t, http.MethodPut, permEndpoint, adminToken, map[string]any{"permisos": []map[string]any{
		{"vmid": 101, "nivelAcceso": "FULL_ACCESS"},
		{"vmid": 102, "nivelAcceso": "READ_ONLY"},
	}}); status != http.StatusNoContent {
		t.Fatalf("precondición: asignar permisos 101 FULL_ACCESS y 102 READ_ONLY, esperado 204, recibido %d: %#v", status, body)
	}

	t.Run("BAC-21B GET /api/instances agrega campos sin romper el contrato de BAC-14", func(t *testing.T) {
		status, list := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if status != http.StatusOK || len(list) != 2 {
			t.Fatalf("GET /instances con OPERATOR: esperado 200 con 101 y 102, recibido %d: %#v", status, list)
		}
		byID := map[float64]map[string]any{}
		for _, item := range list {
			instance := item.(map[string]any)
			byID[instance["id"].(float64)] = instance
			for _, field := range []string{"id", "name", "type", "node", "status"} {
				if _, ok := instance[field]; !ok {
					t.Errorf("campo %q de BAC-14 ausente (rompe FIX-14): %#v", field, instance)
				}
			}
			for _, field := range []string{"ip", "cpuUsage", "ramUsage", "maxRam", "nivelAcceso", "activeTask"} {
				if _, ok := instance[field]; !ok {
					t.Errorf("BAC-21B: falta el campo nuevo %q en la instancia %v", field, instance["id"])
				}
			}
		}
		// Aclaración del 01/10/2026 (FIX-39): ip, cpuUsage, ramUsage, maxRam y activeTask pueden ir en
		// null (los completan BAC-22B, BAC-23B y BAC-25C); nivelAcceso lleva su valor real.
		if _, present := byID[101]["nivelAcceso"]; present && byID[101]["nivelAcceso"] != "FULL_ACCESS" {
			t.Errorf("nivelAcceso de la 101 esperado FULL_ACCESS, recibido %#v", byID[101]["nivelAcceso"])
		}
		if _, present := byID[102]["nivelAcceso"]; present && byID[102]["nivelAcceso"] != "READ_ONLY" {
			t.Errorf("nivelAcceso de la 102 esperado READ_ONLY, recibido %#v", byID[102]["nivelAcceso"])
		}
		if maxRam, ok := byID[101]["maxRam"].(float64); ok && maxRam != 4294967296 {
			t.Errorf("maxRam de la 101 debe venir de maxmem de Proxmox (4294967296), recibido %v", maxRam)
		}
	})

	// BAC-29 y contrato Etapa 1: la ruta canónica y estricta es /instances/:vmid/status/:action.
	// No se permiten rutas ambiguas o de legado.
	canonicalEnergyPath := func(vmid int, action string) string {
		return fmt.Sprintf("/instances/%d/status/%s", vmid, action)
	}

	t.Run("BAC-21B las acciones de energía exigen FULL_ACCESS y quedan en tareas_asincronas", func(t *testing.T) {
		// Esperar que cualquier tarea previa haya concluido
		waitForValue(t, "SELECT count(*) FROM tareas_asincronas WHERE estado = 'RUNNING';", "0", 5*time.Second)
		tasksBefore := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;")
		path := canonicalEnergyPath(101, "stop")
		status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil)
		if status != http.StatusAccepted {
			t.Errorf("FULL_ACCESS sobre la 101, stop (%s): esperado 202, recibido %d: %#v", path, status, body)
		} else {
			tareaID := requiredString(t, body, "tareaId")
			waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 5*time.Second)
		}
		for _, action := range []string{"start", "stop"} {
			path = canonicalEnergyPath(102, action)
			if status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil); status != http.StatusForbidden || body["errorCode"] != "INSTANCE_ACCESS_DENIED" {
				t.Errorf("READ_ONLY sobre la 102, %s (%s): esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", action, path, status, body)
			}
		}
		if after := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;"); after == tasksBefore {
			t.Errorf("las acciones de energía no registraron ninguna fila en tareas_asincronas (antes %s, después %s)", tasksBefore, after)
		}
	})

	t.Run("FIX-39 BAC-21B shutdown y reboot exigen FULL_ACCESS y quedan en tareas_asincronas", func(t *testing.T) {
		waitForValue(t, "SELECT count(*) FROM tareas_asincronas WHERE estado = 'RUNNING';", "0", 5*time.Second)
		tasksBefore := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;")
		for _, action := range []string{"shutdown", "reboot"} {
			path := canonicalEnergyPath(101, action)
			if status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil); status != http.StatusAccepted {
				t.Errorf("FULL_ACCESS sobre la 101, %s (%s): esperado 202, recibido %d: %#v", action, path, status, body)
			} else {
				tareaID := requiredString(t, body, "tareaId")
				waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 5*time.Second)
			}
			path = canonicalEnergyPath(102, action)
			if status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil); status != http.StatusForbidden || body["errorCode"] != "INSTANCE_ACCESS_DENIED" {
				t.Errorf("READ_ONLY sobre la 102, %s (%s): esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", action, path, status, body)
			}
		}
		if after := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;"); after == tasksBefore {
			t.Errorf("shutdown y reboot no registraron ninguna fila en tareas_asincronas")
		}
	})

	t.Run("BAC-21B las acciones de energía quedan en auditoria con el upid de la tarea", func(t *testing.T) {
		waitForValue(t, "SELECT count(*) FROM tareas_asincronas WHERE estado = 'RUNNING';", "0", 5*time.Second)
		status, body := requestJSON(t, http.MethodPost, canonicalEnergyPath(101, "stop"), operatorToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("stop sobre la 101: esperado 202, recibido %d: %#v", status, body)
		}
		// El stub da un UPID distinto por tarea, como Proxmox: se busca el de esta respuesta.
		upid := strings.ReplaceAll(requiredString(t, body, "upid"), "'", "''")
		query := fmt.Sprintf("SELECT count(*) FROM auditoria WHERE (instancia_id = '101' OR detalles::text LIKE '%%101%%') AND detalles::text LIKE '%%%s%%';", upid)
		for attempt := 0; attempt < 30; attempt++ {
			if queryDatabase(t, query) != "0" {
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
		t.Errorf("ninguna fila de auditoria registra el start de la 101 con su upid (las acciones de energía no se auditan)")
	})

	t.Run("BAC-21B DELETE /instances/:vmid solo ADMIN y 409 si la instancia está encendida", func(t *testing.T) {
		// El endpoint de borrado lo crea BAC-24B (Etapa 1); BAC-21B solo fija sus reglas de acceso.
		// Si todavía no existe, no hay regla que verificar y el caso se omite.
		if status, _ := requestJSON(t, http.MethodDelete, "/instances/102", adminToken, nil); status == http.StatusNotFound || status == http.StatusMethodNotAllowed {
			t.Skipf("DELETE /instances/:vmid todavía no existe (responde %d); sus reglas se verifican cuando llegue BAC-24B", status)
		}
		if status, body := requestJSON(t, http.MethodDelete, "/instances/101", operatorToken, nil); status != http.StatusForbidden {
			t.Errorf("DELETE con OPERATOR (aun con FULL_ACCESS): esperado 403, recibido %d: %#v", status, body)
		}
		if status, body := requestJSON(t, http.MethodDelete, "/instances/101", adminToken, nil); status != http.StatusConflict {
			t.Errorf("DELETE de la 101 (running) con ADMIN: esperado 409, recibido %d: %#v", status, body)
		}
	})

	t.Run("BAC-21C POST /api/events/ticket emite un ticket de un solo uso", func(t *testing.T) {
		if status, _ := requestValue(t, http.MethodPost, "/events/ticket", "", nil); status != http.StatusUnauthorized {
			t.Errorf("POST /events/ticket sin token: esperado 401, recibido %d", status)
		}
		status, body := requestJSON(t, http.MethodPost, "/events/ticket", operatorToken, nil)
		if status != http.StatusOK && status != http.StatusCreated {
			t.Fatalf("POST /events/ticket con access token: esperado 200/201, recibido %d: %#v", status, body)
		}
		ticket := requiredString(t, body, "ticket")
		if keys, _ := runCompose("exec", "-T", "redis", "redis-cli", "-a", "redis-test", "--no-auth-warning", "EXISTS", "ws_ticket:"+ticket); strings.TrimSpace(keys) != "1" {
			t.Errorf("el ticket no está guardado en Redis como ws_ticket:<uuid>")
		}
		ttlOutput, _ := runCompose("exec", "-T", "redis", "redis-cli", "-a", "redis-test", "--no-auth-warning", "TTL", "ws_ticket:"+ticket)
		if ttl, err := strconv.Atoi(strings.TrimSpace(ttlOutput)); err != nil || ttl <= 0 || ttl > 30 {
			t.Errorf("el ticket debe expirar en 30 s o menos, TTL=%q", strings.TrimSpace(ttlOutput))
		}
	})

	t.Run("BAC-21C /api/events rechaza tickets inválidos o reutilizados y corta el stream al hacer logout", func(t *testing.T) {
		user := createActiveUser(t, adminToken, "op_eventos", "OPERATOR")
		session := loginWithTOTP(t, user.Email, user.Password, user.Secret)
		client := &http.Client{Timeout: 0}

		if status := openEvents(t, client, "ticket-invalido", nil); status != http.StatusUnauthorized {
			t.Errorf("GET /events con ticket inválido: esperado 401, recibido %d", status)
		}

		status, body := requestJSON(t, http.MethodPost, "/events/ticket", session.AccessToken, nil)
		if status != http.StatusOK && status != http.StatusCreated {
			t.Fatalf("POST /events/ticket: esperado 200/201, recibido %d: %#v", status, body)
		}
		ticket := requiredString(t, body, "ticket")

		closed := make(chan struct{})
		if status := openEvents(t, client, ticket, closed); status != http.StatusOK {
			t.Fatalf("GET /events con ticket válido: esperado 200 (stream abierto), recibido %d", status)
		}
		if status := openEvents(t, client, ticket, nil); status != http.StatusUnauthorized {
			t.Errorf("GET /events reutilizando el ticket: esperado 401, recibido %d", status)
		}

		var cookies []*http.Cookie
		if session.RefreshCookie != nil {
			cookies = append(cookies, session.RefreshCookie)
		}
		if logout := requestRaw(t, http.MethodPost, "/auth/logout", session.AccessToken, nil, cookies...); logout.Status != http.StatusNoContent {
			t.Fatalf("precondición: logout esperado 204, recibido %d: %s", logout.Status, logout.RawBody)
		}
		select {
		case <-closed:
		case <-time.After(5 * time.Second):
			t.Errorf("el stream de /events sigue abierto 5 s después del logout; debe cortarse inmediatamente")
		}
	})

	t.Run("FIX-40 D2 un Proxmox que no responde da 504 PROXMOX_TIMEOUT y uno caído 502 PROXMOX_UNAVAILABLE", func(t *testing.T) {
		// backend-lento (compose.yaml) apunta a proxmox-lento, que acepta la conexión y nunca contesta:
		// el cliente corta a los 10 s (requestTimeout) y eso tiene que llegar como 504 con su propio código.
		// El cliente de la prueba espera 20 s: más que los 10 s del backend, para que corte el backend.
		request, _ := http.NewRequest(http.MethodGet, envOrDefault("BACKEND_TEST_SLOW_API_URL", "http://127.0.0.1:18082/api")+"/instances/101", nil)
		request.Header.Set("Authorization", "Bearer "+adminToken)
		started := time.Now()
		response, err := (&http.Client{Timeout: 20 * time.Second}).Do(request)
		if err != nil {
			t.Fatalf("GET /instances/101 en backend-lento falló sin respuesta HTTP: %v", err)
		}
		var body map[string]any
		_ = json.NewDecoder(response.Body).Decode(&body)
		response.Body.Close()
		if response.StatusCode != http.StatusGatewayTimeout || body["errorCode"] != "PROXMOX_TIMEOUT" {
			t.Errorf("Proxmox sin responder: esperado 504 PROXMOX_TIMEOUT, recibido %d %v (%v)", response.StatusCode, body["errorCode"], time.Since(started).Round(time.Second))
		}

		// Contraprueba: con el stub caído (la conexión falla enseguida) sigue siendo 502 PROXMOX_UNAVAILABLE.
		stopProxmoxStub(t)
		status, body := requestJSON(t, http.MethodGet, "/instances/101", adminToken, nil)
		startProxmoxStub(t)
		if status != http.StatusBadGateway || body["errorCode"] != "PROXMOX_UNAVAILABLE" {
			t.Errorf("Proxmox caído: esperado 502 PROXMOX_UNAVAILABLE, recibido %d: %#v", status, body)
		}
	})

	t.Run("INF-07B BRG-03 /api/events a través del borde Nginx entrega TASK_FINISHED sin buffer", func(t *testing.T) {
		// 1. La configuración versionada del borde (docker/nginx-edge.conf, FIX-31) es la referencia
		// del servidor: el location que atiende /api/events no debe acumular el stream ni cortarlo
		// por inactividad antes de 5 minutos.
		edgeConf := string(readFile(t, "../../docker/nginx-edge.conf"))
		block := ""
		for _, pattern := range []string{`location\s+(=\s*)?/api/events/?\s*\{([^}]*)\}`, `location\s+/api/\s*\{([^}]*)\}`} {
			if match := regexp.MustCompile(pattern).FindStringSubmatch(edgeConf); match != nil {
				block = match[len(match)-1]
				break
			}
		}
		if block == "" {
			t.Fatalf("docker/nginx-edge.conf no tiene un location /api/events ni /api/ que atienda el stream")
		}
		if !regexp.MustCompile(`proxy_http_version\s+1\.1\s*;`).MatchString(block) {
			t.Errorf("el location de /api/events debe usar proxy_http_version 1.1 (SSE con conexión persistente)")
		}
		if !regexp.MustCompile(`proxy_buffering\s+off\s*;`).MatchString(block) {
			t.Errorf("el location de /api/events debe tener proxy_buffering off: si no, Nginx retiene los eventos SSE")
		}
		if seconds := nginxDurationSeconds(regexp.MustCompile(`proxy_read_timeout\s+([0-9]+[smhd]?)\s*;`).FindStringSubmatch(block)); seconds < 300 {
			t.Errorf("proxy_read_timeout del location de /api/events debe ser de al menos 5 min (default de Nginx: 60 s), es de %d s", seconds)
		}

		// 2. Comportamiento: el circuito completo pasa por el borde HTTPS (servicio edge de compose.yaml).
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(readFile(t, "edge-tls/ca.crt")) {
			t.Fatalf("no se pudo cargar edge-tls/ca.crt")
		}
		edgeURL := envOrDefault("BACKEND_TEST_EDGE_HTTPS_URL", "https://127.0.0.1:18453") + "/api"
		transport := &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool, ServerName: "centinela.test"}}
		client := &http.Client{Timeout: 10 * time.Second, Transport: transport}
		stream := &http.Client{Timeout: 0, Transport: transport}

		user := createActiveUser(t, adminToken, "admin_borde", "ADMIN")
		session := loginWithTOTP(t, user.Email, user.Password, user.Secret)
		edgeRequest := func(method, path string) (int, map[string]any) {
			t.Helper()
			request, _ := http.NewRequest(method, edgeURL+path, nil)
			request.Header.Set("Authorization", "Bearer "+session.AccessToken)
			response, err := client.Do(request)
			if err != nil {
				t.Fatalf("%s %s a través del borde falló: %v", method, path, err)
			}
			defer response.Body.Close()
			var body map[string]any
			_ = json.NewDecoder(response.Body).Decode(&body)
			return response.StatusCode, body
		}

		status, body := edgeRequest(http.MethodPost, "/events/ticket")
		if status != http.StatusOK && status != http.StatusCreated {
			t.Fatalf("POST /api/events/ticket por el borde: esperado 200/201, recibido %d: %#v", status, body)
		}
		ticket := requiredString(t, body, "ticket")

		request, _ := http.NewRequest(http.MethodGet, edgeURL+"/events?ticket="+ticket, nil)
		request.Header.Set("Accept", "text/event-stream")
		response, err := stream.Do(request)
		if err != nil {
			t.Fatalf("GET /api/events por el borde falló: %v", err)
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK || !strings.HasPrefix(response.Header.Get("Content-Type"), "text/event-stream") {
			t.Fatalf("GET /api/events por el borde: esperado 200 text/event-stream, recibido %d %q", response.StatusCode, response.Header.Get("Content-Type"))
		}
		events := make(chan map[string]any, 16)
		go func() {
			defer close(events)
			reader := bufio.NewReader(response.Body)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				if payload, ok := strings.CutPrefix(strings.TrimRight(line, "\r\n"), "data: "); ok {
					var event map[string]any
					if json.Unmarshal([]byte(payload), &event) == nil {
						events <- event
					}
				}
			}
		}()

		// 102 (lxc) está detenida en el stub: start es válido con la ruta canónica /status/start (BAC-21B).
		status, body = edgeRequest(http.MethodPost, "/instances/102/status/start")
		if status != http.StatusAccepted {
			t.Fatalf("start de la 102 por el borde: esperado 202, recibido %d: %#v", status, body)
		}
		tareaID := requiredString(t, body, "tareaId")
		dispatched := time.Now()

		// El stub da la tarea por terminada en la primera consulta (cada 1 s): con el stream sin
		// buffer, el evento llega en pocos segundos. Con buffer, quedaría retenido en Nginx.
		deadline := time.After(8 * time.Second)
		for {
			select {
			case event, open := <-events:
				if !open {
					t.Fatalf("el borde cerró el stream de /api/events antes de entregar TASK_FINISHED")
				}
				detalles, _ := event["detalles"].(map[string]any)
				if event["tipo"] != "TASK_FINISHED" || detalles["tareaId"] != tareaID {
					continue
				}
				if event["recursoId"] != "102" {
					t.Errorf("TASK_FINISHED de la tarea %s con recursoId %v, esperado 102", tareaID, event["recursoId"])
				}
				t.Logf("TASK_FINISHED recibido por el borde %v después del 202", time.Since(dispatched).Round(time.Millisecond))
				return
			case <-deadline:
				t.Fatalf("no llegó TASK_FINISHED de la tarea %s por el borde en 8 s: Nginx retiene el stream o el backend no lo emitió", tareaID)
			}
		}
	})

	t.Run("RNF-07_bus_PubSub_Redis_multirreplica_propaga_eventos_entre_instancias_backend", func(t *testing.T) {
		smtpURL := envOrDefault("BACKEND_TEST_SMTP_API_URL", "http://127.0.0.1:18081/api")

		waitForValue(t, "SELECT count(*) FROM tareas_asincronas WHERE estado = 'RUNNING';", "0", 5*time.Second)

		status, body := requestJSON(t, http.MethodPost, "/events/ticket", adminToken, nil)
		if status != http.StatusOK {
			t.Fatalf("POST /events/ticket en backend 1: esperado 200, recibido %d: %#v", status, body)
		}
		ticket := requiredString(t, body, "ticket")

		events := make(chan map[string]any, 10)
		req, err := http.NewRequest(http.MethodGet, apiURL+"/events?ticket="+ticket, nil)
		if err != nil {
			t.Fatalf("error creando request SSE: %v", err)
		}
		req.Header.Set("Accept", "text/event-stream")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("error conectando SSE en backend 1: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("SSE backend 1: esperado 200, recibido %d", resp.StatusCode)
		}

		go func() {
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "data: ") {
					payload := strings.TrimPrefix(line, "data: ")
					var event map[string]any
					if json.Unmarshal([]byte(payload), &event) == nil {
						events <- event
					}
				}
			}
		}()

		reqAction, err := http.NewRequest(http.MethodPost, smtpURL+"/instances/102/status/start", nil)
		if err != nil {
			t.Fatalf("error creando request a backend 2: %v", err)
		}
		reqAction.Header.Set("Authorization", "Bearer "+adminToken)
		respAction, err := http.DefaultClient.Do(reqAction)
		if err != nil {
			t.Fatalf("error ejecutando acción en backend 2: %v", err)
		}
		defer respAction.Body.Close()
		bodyBytes, _ := io.ReadAll(respAction.Body)
		if respAction.StatusCode != http.StatusAccepted {
			t.Fatalf("start en backend 2: esperado 202, recibido %d: %s", respAction.StatusCode, string(bodyBytes))
		}
		var actionBody map[string]any
		json.Unmarshal(bodyBytes, &actionBody)
		tareaID := requiredString(t, actionBody, "tareaId")

		deadline := time.After(8 * time.Second)
		for {
			select {
			case ev := <-events:
				detalles, _ := ev["detalles"].(map[string]any)
				if ev["tipo"] == "TASK_FINISHED" && detalles["tareaId"] == tareaID {
					if ev["recursoId"] != "102" {
						t.Errorf("TASK_FINISHED recibido con recursoId %v, esperado 102", ev["recursoId"])
					}
					t.Logf("RNF-07 verificado: evento de tarea %s emitido por backend 2 recibido en SSE de backend 1 vía Redis Pub/Sub", tareaID)
					return
				}
			case <-deadline:
				t.Fatalf("timeout: no se recibió TASK_FINISHED de la tarea %s emitida por backend 2 en el SSE de backend 1. El bus Redis Pub/Sub multirréplica falló", tareaID)
			}
		}
	})
}

// nginxDurationSeconds convierte el valor capturado de una directiva de tiempo de Nginx
// (30, 30s, 5m, 1h, 1d) a segundos. Sin valor devuelve 60, el default de proxy_read_timeout.
func nginxDurationSeconds(match []string) int {
	if len(match) < 2 {
		return 60
	}
	value := match[1]
	unit := 1
	switch value[len(value)-1] {
	case 's':
		value = value[:len(value)-1]
	case 'm':
		value, unit = value[:len(value)-1], 60
	case 'h':
		value, unit = value[:len(value)-1], 3600
	case 'd':
		value, unit = value[:len(value)-1], 86400
	}
	number, _ := strconv.Atoi(value)
	return number * unit
}

// openEvents abre GET /api/events?ticket=... y devuelve el status. Si closed no es nil y el stream
// quedó abierto (200), una goroutine lo lee hasta que el servidor lo cierre y cierra el canal.
func openEvents(t *testing.T, client *http.Client, ticket string, closed chan struct{}) int {
	t.Helper()
	request, err := http.NewRequest(http.MethodGet, apiURL+"/events?ticket="+ticket, nil)
	if err != nil {
		t.Fatalf("no se pudo crear GET /events: %v", err)
	}
	request.Header.Set("Accept", "text/event-stream")
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("GET /events fallido: %v", err)
	}
	if response.StatusCode != http.StatusOK || closed == nil {
		response.Body.Close()
		return response.StatusCode
	}
	go func() {
		defer response.Body.Close()
		reader := bufio.NewReader(response.Body)
		for {
			if _, err := reader.ReadString('\n'); err != nil {
				close(closed)
				return
			}
		}
	}()
	return response.StatusCode
}
