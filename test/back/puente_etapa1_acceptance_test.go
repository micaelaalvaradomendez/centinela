package back_test

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Tareas puente entre la fase base y la Etapa 1 (documentacion/actual.md):
//   - BAC-21B: GET /api/instances extendido de forma retrocompatible, POST /instances/:vmid/status/:action
//     con alias /start y /stop, DELETE solo ADMIN con 409 si está encendida, y registro en auditoria y
//     tareas_asincronas.
//   - BAC-21C: tickets efímeros en Redis para /api/events y corte del stream al revocar la sesión.
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
		if byID[101]["nivelAcceso"] != nil && byID[101]["nivelAcceso"] != "FULL_ACCESS" {
			t.Errorf("nivelAcceso de la 101 esperado FULL_ACCESS, recibido %#v", byID[101]["nivelAcceso"])
		}
		if byID[102]["nivelAcceso"] != nil && byID[102]["nivelAcceso"] != "READ_ONLY" {
			t.Errorf("nivelAcceso de la 102 esperado READ_ONLY, recibido %#v", byID[102]["nivelAcceso"])
		}
		if maxRam, ok := byID[101]["maxRam"].(float64); ok && maxRam != 4294967296 {
			t.Errorf("maxRam de la 101 debe venir de maxmem de Proxmox (4294967296), recibido %v", maxRam)
		}
	})

	t.Run("BAC-21B POST /instances/:vmid/status/:action con FULL_ACCESS, alias y 403 a READ_ONLY", func(t *testing.T) {
		auditBefore := queryDatabase(t, "SELECT count(*) FROM auditoria;")
		tasksBefore := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;")

		status, body := requestJSON(t, http.MethodPost, "/instances/101/status/start", operatorToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("POST /instances/101/status/start con FULL_ACCESS: esperado 202, recibido %d: %#v", status, body)
		}
		for _, path := range []string{"/instances/101/start", "/instances/101/stop"} {
			if status, body := requestJSON(t, http.MethodPost, path, operatorToken, nil); status != http.StatusAccepted {
				t.Errorf("alias %s: esperado 202, recibido %d: %#v", path, status, body)
			}
		}
		for _, action := range []string{"start", "stop"} {
			status, body := requestJSON(t, http.MethodPost, "/instances/102/status/"+action, operatorToken, nil)
			if status != http.StatusForbidden || body["errorCode"] != "INSTANCE_ACCESS_DENIED" {
				t.Errorf("READ_ONLY sobre la 102, %s: esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", action, status, body)
			}
		}

		if after := queryDatabase(t, "SELECT count(*) FROM tareas_asincronas;"); after == tasksBefore {
			t.Errorf("las acciones de energía no registraron ninguna fila en tareas_asincronas (antes %s, después %s)", tasksBefore, after)
		}
		if upid := queryDatabase(t, "SELECT count(*) FROM auditoria WHERE detalles::text LIKE '%UPID:pve:stub:start:101%';"); upid == "0" {
			t.Errorf("auditoria no guarda el upid de la acción start sobre la 101 en detalles (filas antes: %s)", auditBefore)
		}
	})

	t.Run("BAC-21B DELETE /instances/:vmid solo ADMIN y 409 si la instancia está encendida", func(t *testing.T) {
		if status, body := requestJSON(t, http.MethodDelete, "/instances/101", operatorToken, nil); status != http.StatusForbidden {
			t.Errorf("DELETE con OPERATOR (aun con FULL_ACCESS): esperado 403, recibido %d: %#v", status, body)
		}
		if status, body := requestJSON(t, http.MethodDelete, "/instances/101", adminToken, nil); status != http.StatusConflict {
			t.Errorf("DELETE de la 101 (running) con ADMIN: esperado 409, recibido %d: %#v", status, body)
		}
		// Contraprueba: una instancia detenida sí se puede eliminar.
		if status, body := requestJSON(t, http.MethodDelete, "/instances/102", adminToken, nil); status != http.StatusAccepted && status != http.StatusOK && status != http.StatusNoContent {
			t.Errorf("DELETE de la 102 (stopped) con ADMIN: esperado 2xx, recibido %d: %#v", status, body)
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
