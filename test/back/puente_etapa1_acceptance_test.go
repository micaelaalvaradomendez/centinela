package back_test

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
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

	// BAC-21B se verifica contra su criterio de éxito, no contra la forma sugerida en el entregable:
	// la tarea propone POST /instances/:vmid/status/:action con alias /start y /stop, pero el criterio
	// pide que las acciones de energía exijan FULL_ACCESS y queden registradas. Se acepta cualquiera
	// de las dos rutas.
	energyPath := func(t *testing.T, vmid int, action string) string {
		t.Helper()
		unified := fmt.Sprintf("/instances/%d/status/%s", vmid, action)
		if status, _ := requestJSON(t, http.MethodPost, unified, adminToken, nil); status != http.StatusNotFound && status != http.StatusMethodNotAllowed {
			return unified
		}
		return fmt.Sprintf("/instances/%d/%s", vmid, action)
	}

	t.Run("BAC-21B las acciones de energía exigen FULL_ACCESS y quedan en tareas_asincronas", func(t *testing.T) {
		tasksBefore := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;")
		for _, action := range []string{"start", "stop"} {
			path := energyPath(t, 101, action)
			if status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil); status != http.StatusAccepted {
				t.Errorf("FULL_ACCESS sobre la 101, %s (%s): esperado 202, recibido %d: %#v", action, path, status, body)
			}
			path = energyPath(t, 102, action)
			if status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil); status != http.StatusForbidden || body["errorCode"] != "INSTANCE_ACCESS_DENIED" {
				t.Errorf("READ_ONLY sobre la 102, %s (%s): esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", action, path, status, body)
			}
		}
		if after := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;"); after == tasksBefore {
			t.Errorf("las acciones de energía no registraron ninguna fila en tareas_asincronas (antes %s, después %s)", tasksBefore, after)
		}
	})

	t.Run("FIX-39 BAC-21B shutdown y reboot exigen FULL_ACCESS y quedan en tareas_asincronas", func(t *testing.T) {
		tasksBefore := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;")
		for _, action := range []string{"shutdown", "reboot"} {
			path := energyPath(t, 101, action)
			if status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil); status != http.StatusAccepted {
				t.Errorf("FULL_ACCESS sobre la 101, %s (%s): esperado 202, recibido %d: %#v", action, path, status, body)
			}
			path = energyPath(t, 102, action)
			if status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil); status != http.StatusForbidden || body["errorCode"] != "INSTANCE_ACCESS_DENIED" {
				t.Errorf("READ_ONLY sobre la 102, %s (%s): esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", action, path, status, body)
			}
		}
		if after := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;"); after == tasksBefore {
			t.Errorf("shutdown y reboot no registraron ninguna fila en tareas_asincronas")
		}
	})

	t.Run("BAC-21B las acciones de energía quedan en auditoria con el upid de la tarea", func(t *testing.T) {
		// Criterio: "todo se registra en auditoria". La forma (columnas accion/instancia_id o claves de
		// detalles) es libre; se exige una fila de la 101 que lleve el upid de la acción. Se espera hasta
		// 15 s por si se audita al terminar la tarea (TASK_FINISHED).
		if status, body := requestJSON(t, http.MethodPost, energyPath(t, 101, "start"), operatorToken, nil); status != http.StatusAccepted {
			t.Fatalf("start sobre la 101: esperado 202, recibido %d: %#v", status, body)
		}
		query := "SELECT count(*) FROM auditoria WHERE (instancia_id = '101' OR detalles::text LIKE '%101%') AND detalles::text LIKE '%UPID:pve:stub:start:101%';"
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

		// 102 (lxc) está detenida en el stub: start es válido aunque el backend valide el estado previo (BAC-24A).
		status, body = edgeRequest(http.MethodPost, "/instances/102/start")
		if status == http.StatusNotFound || status == http.StatusMethodNotAllowed {
			status, body = edgeRequest(http.MethodPost, "/instances/102/status/start")
		}
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
