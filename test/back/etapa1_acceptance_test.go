package back_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Tareas de la Etapa 1 en curso (documentacion/actual.md, Ola 1 de etapa1.md):
//   - BAC-29  Contrato HTTP y de eventos de la etapa (docs/contrato-etapa1.md y Swagger)
//   - BAC-22  GET /api/node/status con caché en Redis y último estado conocido (D1)
//   - BAC-23A Adaptador de inventario con IP por guest agent / interfaces de LXC (función interna)
//   - BAC-25A Pool acotado de workers para el seguimiento de UPID
//
// Stub de Proxmox (proxmox-stub/nginx.conf): /nodes/pve/status con CPU 25 % (8 núcleos), RAM 8/16 GiB,
// disco 25/100 GiB y uptime 266400 s; toda tarea figura terminada con exitstatus OK.
func TestEtapa1(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	t.Run("BAC-29 el contrato de la Etapa 1 está publicado en docs/contrato-etapa1.md y en Swagger", func(t *testing.T) {
		// BAC-29 es la tarea que fija el contrato: los nombres salen de su entregable. Si el equipo los
		// cambia en el documento, se actualizan la tarea y esta lista.
		contract := string(readFile(t, sourcePath("backend", "docs/contrato-etapa1.md")))
		required := map[string][]string{
			"GET /api/node/status (entregable 1)":   {"/node/status", "usagePercent", "cores", "usedGb", "totalGb", "uptimeSeconds", "instancesSummary", "stale", "fetchedAt"},
			"GET /api/instances (entregable 2)":     {"cpuUsage", "ramUsage", "maxRam", "nivelAcceso", "activeTask"},
			"acciones y DELETE (entregable 3)":      {"start", "shutdown", "stop", "reboot", "DELETE", "upid", "tareaId"},
			"TASK_FINISHED (entregable 4, D2)":      {"TASK_FINISHED", "exitstatus", "motivo", "PROXMOX_ERROR", "TIMEOUT", "COMPLETED", "FAILED"},
			"códigos de error (entregable 6, D2)":   {"INSTANCE_ACCESS_DENIED", "INSTANCE_PROTECTED", "INSTANCE_BUSY", "INSTANCE_INVALID_STATE", "INSTANCE_NOT_FOUND", "PROXMOX_UNAVAILABLE", "PROXMOX_TIMEOUT", "INVALID_ACTION"},
		}
		for section, tokens := range required {
			for _, token := range tokens {
				if !strings.Contains(contract, token) {
					t.Errorf("docs/contrato-etapa1.md, %s: no menciona %q", section, token)
				}
			}
		}
		if !regexp.MustCompile(`(?i)node/status[^\n]*(autenticad|RequireAuth|Bearer)|(autenticad|RequireAuth|Bearer)[^\n]*node/status`).MatchString(contract) {
			t.Errorf("docs/contrato-etapa1.md no indica que GET /api/node/status lo puede consultar cualquier usuario autenticado (D1)")
		}
		swagger := string(readFile(t, sourcePath("backend", "docs/swagger.json")))
		if !strings.Contains(swagger, `"/node/status"`) {
			t.Errorf("docs/swagger.json no documenta /node/status")
		}
		// FIX-08: los códigos nuevos también van al inventario de errorCode (docs/estandar_http.md).
		inventory := string(readFile(t, sourcePath("backend", "docs/estandar_http.md")))
		for _, code := range []string{"INSTANCE_INVALID_STATE", "PROXMOX_TIMEOUT"} {
			if !strings.Contains(inventory, code) {
				t.Errorf("docs/estandar_http.md (inventario de FIX-08) no incluye %s", code)
			}
		}
	})

	t.Run("BAC-22 GET /api/node/status normaliza la telemetría, la cachea en Redis y sirve el último estado conocido", func(t *testing.T) {
		// Paso 1. Sin estado conocido y con Proxmox caído: error controlado (502 o 504), nunca 200 ni 500.
		// Corre primero: el stack es nuevo en cada ejecución, así que todavía no hay nada en caché.
		stopProxmoxStub(t)
		status, body := requestJSON(t, http.MethodGet, "/node/status", adminToken, nil)
		if status == http.StatusNotFound {
			startProxmoxStub(t)
			t.Fatalf("BAC-22 no implementada: GET /api/node/status responde 404")
		}
		if (status != http.StatusBadGateway && status != http.StatusGatewayTimeout) ||
			(body["errorCode"] != "PROXMOX_UNAVAILABLE" && body["errorCode"] != "PROXMOX_TIMEOUT") {
			t.Errorf("sin estado conocido y con Proxmox caído: esperado 502 PROXMOX_UNAVAILABLE o 504 PROXMOX_TIMEOUT, recibido %d: %#v", status, body)
		}
		startProxmoxStub(t)

		// Paso 2. D1: sin token 401; un OPERATOR real (no solo el ADMIN) recibe 200.
		if status, _ := requestValue(t, http.MethodGet, "/node/status", "", nil); status != http.StatusUnauthorized {
			t.Errorf("GET /api/node/status sin token: esperado 401, recibido %d", status)
		}
		operator := createActiveUser(t, adminToken, "op_nodo", "OPERATOR")
		session := loginWithTOTP(t, operator.Email, operator.Password, operator.Secret)
		status, body = requestJSON(t, http.MethodGet, "/node/status", session.AccessToken, nil)
		if status != http.StatusOK {
			t.Fatalf("GET /api/node/status con un OPERATOR: esperado 200 (D1), recibido %d: %#v", status, body)
		}

		// Paso 3. Normalización (contrato de BAC-29): porcentaje y núcleos de CPU, GB de RAM y disco, uptime en s.
		// Se acepta GB decimal (1e9) o binario (2^30): la tarea pide "GB" sin fijar la base.
		expectNumber(t, body, 25, 0.5, "cpu", "usagePercent")
		expectNumber(t, body, 8, 0, "cpu", "cores")
		expectGB(t, body, 8, "ram", "usedGb")
		expectGB(t, body, 16, "ram", "totalGb")
		expectNumber(t, body, 50, 0.5, "ram", "usagePercent")
		expectGB(t, body, 25, "storage", "usedGb")
		expectGB(t, body, 100, "storage", "totalGb")
		expectNumber(t, body, 25, 0.5, "storage", "usagePercent")
		expectNumber(t, body, 266400, 0, "uptimeSeconds")
		if body["stale"] != false {
			t.Errorf("con Proxmox disponible stale debe ser false, recibido %#v", body["stale"])
		}
		if fetchedAt, _ := body["fetchedAt"].(string); fetchedAt == "" {
			t.Errorf("falta fetchedAt (RFC3339) en la respuesta")
		} else if _, err := time.Parse(time.RFC3339, fetchedAt); err != nil {
			t.Errorf("fetchedAt %q no es RFC3339: %v", fetchedAt, err)
		}

		// Paso 4. La telemetría queda en Redis con un TTL de 5 a 10 s.
		if !redisHasKeyWithTTL(t, 1, 10) {
			t.Errorf("después de consultar /node/status no hay ninguna clave en Redis con TTL de 10 s o menos (se esperan 5 a 10 s)")
		}

		// Paso 5. Dentro del TTL responde desde la caché: aunque Proxmox esté caído, 200 con stale false y en < 50 ms.
		stopProxmoxStub(t)
		best := time.Hour
		for range 3 {
			started := time.Now()
			status, body = requestJSON(t, http.MethodGet, "/node/status", adminToken, nil)
			if elapsed := time.Since(started); elapsed < best {
				best = elapsed
			}
			if status != http.StatusOK || body["stale"] != false {
				t.Fatalf("dentro del TTL y con Proxmox caído: esperado 200 desde la caché con stale false, recibido %d: %#v", status, body)
			}
		}
		if best >= 50*time.Millisecond {
			t.Errorf("con la caché vigente la respuesta debe tardar menos de 50 ms; la mejor de 3 tardó %v", best)
		}

		// Paso 6. Vencido el TTL y con Proxmox caído: el último estado conocido, con stale true.
		time.Sleep(11 * time.Second)
		status, body = requestJSON(t, http.MethodGet, "/node/status", adminToken, nil)
		if status != http.StatusOK || body["stale"] != true {
			t.Errorf("vencido el TTL y con Proxmox caído: esperado 200 con el último estado conocido y stale true, recibido %d: %#v", status, body)
		}
		expectNumber(t, body, 266400, 0, "uptimeSeconds")
		startProxmoxStub(t)

		// Paso 7 (FIX-62). Al restablecer Proxmox y vencer el TTL, el endpoint debe recuperarse: 200 con stale false y telemetría actualizada.
		time.Sleep(11 * time.Second)
		status, body = requestJSON(t, http.MethodGet, "/node/status", adminToken, nil)
		if status != http.StatusOK || body["stale"] != false {
			t.Errorf("tras restablecer Proxmox y vencer el TTL: esperado 200 con stale false (recuperación de telemetría), recibido %d: %#v", status, body)
		}
	})

	t.Run("BAC-23A el adaptador de inventario resuelve la IP con guest agent e interfaces de LXC y tiene pruebas unitarias", func(t *testing.T) {
		// BAC-23A es una función interna (todavía no hay endpoint que la exponga; eso es BAC-23B). Se
		// verifica que el código consulte los dos endpoints de IP y que sus pruebas unitarias pasen.
		sources, tests := goFilesContaining(t, `agent/network-get-interfaces`)
		lxcSources, lxcTests := goFilesContaining(t, `/interfaces`)
		if len(sources) == 0 {
			t.Errorf("ningún archivo del backend (fuera de cmd/proxmox-simulador) consulta qemu/{vmid}/agent/network-get-interfaces")
		}
		if len(lxcSources) == 0 {
			t.Errorf("ningún archivo del backend (fuera de cmd/proxmox-simulador) consulta lxc/{vmid}/interfaces")
		}
		if len(tests)+len(lxcTests) == 0 {
			t.Fatalf("no hay pruebas unitarias del adaptador de IP (ningún _test.go del backend menciona network-get-interfaces ni /interfaces)")
		}
		runBackendUnitTests(t, append(append(sources, tests...), append(lxcSources, lxcTests...)...))
	})

	t.Run("BAC-25A el seguimiento de UPID usa un pool acotado y no pierde tareas", func(t *testing.T) {
		// La cota de concurrencia no se observa desde la API: la verifican las pruebas unitarias del
		// backend (entregable 4). Desde afuera se verifica la configuración y que ninguna tarea se pierda.
		sources, _ := goFilesContaining(t, `UPID_WORKERS`)
		if len(sources) == 0 {
			t.Errorf("el backend no lee UPID_WORKERS (cantidad de workers del pool, por defecto 8)")
		}
		_, tests := goFilesContaining(t, `seguimiento|Seguir\(`)
		runBackendUnitTests(t, tests)

		// 20 órdenes seguidas sobre la 102 (detenida en el stub): las 20 terminan COMPLETED en tareas_asincronas.
		const total = 20
		var ids []string
		for range total {
			status, body := requestJSON(t, http.MethodPost, "/instances/102/start", adminToken, nil)
			if status == http.StatusNotFound || status == http.StatusMethodNotAllowed {
				status, body = requestJSON(t, http.MethodPost, "/instances/102/status/start", adminToken, nil)
			}
			if status != http.StatusAccepted {
				t.Fatalf("start de la 102: esperado 202, recibido %d: %#v", status, body)
			}
			tareaID := requiredString(t, body, "tareaId")
			ids = append(ids, "'"+tareaID+"'")
			waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 5*time.Second)
		}
		dispatched := time.Now()
		query := fmt.Sprintf("SELECT count(*) FROM tareas_asincronas WHERE id IN (%s) AND estado = 'COMPLETED';", strings.Join(ids, ","))
		completed := 0
		// El stub da cada tarea por terminada en la primera consulta: con el intervalo de 1 s y el
		// criterio de "< 2 s desde la finalización", 3 s alcanzan para las 20.
		for time.Since(dispatched) < 3*time.Second {
			completed, _ = strconv.Atoi(strings.TrimSpace(queryDatabase(t, query)))
			if completed == total {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		if completed != total {
			t.Errorf("de %d tareas despachadas, %d quedaron COMPLETED 3 s después; las demás se perdieron o tardaron demasiado", total, completed)
		}
	})

	t.Run("BAC-24A energía: 409 si el estado no corresponde, 403 en VMIDs protegidos, 202 en órdenes válidas, sin tráfico a Proxmox al rechazar", func(t *testing.T) {
		// Stub: 101 qemu running, 102 lxc stopped, 103 qemu running y protegida (PROXMOX_PROTECTED_VMIDS en compose.yaml).
		// Se prueban la ruta genérica /status/:action y los alias /start y /stop, que siguen montados.
		routes := func(vmid int, action string) []string {
			paths := []string{fmt.Sprintf("/instances/%d/status/%s", vmid, action)}
			if action == "start" || action == "stop" {
				paths = append(paths, fmt.Sprintf("/instances/%d/%s", vmid, action))
			}
			return paths
		}
		rejected := func(vmid int, action string, wantStatus int, wantCode string) {
			t.Helper()
			for _, path := range routes(vmid, action) {
				before := stubWrites(t, vmid)
				status, body := requestJSON(t, http.MethodPost, path, adminToken, nil)
				if status != wantStatus || body["errorCode"] != wantCode {
					t.Errorf("%s: esperado %d %s, recibido %d: %#v", path, wantStatus, wantCode, status, body)
				}
				if after := stubWrites(t, vmid); after != before {
					t.Errorf("%s: el backend envió la orden a Proxmox aunque la rechazó (%d escrituras nuevas en el stub)", path, after-before)
				}
			}
		}

		// 1. Estado previo (D2): start exige stopped; stop, shutdown y reboot exigen running.
		rejected(101, "start", http.StatusConflict, "INSTANCE_INVALID_STATE")
		for _, action := range []string{"stop", "shutdown", "reboot"} {
			rejected(102, action, http.StatusConflict, "INSTANCE_INVALID_STATE")
		}
		// 2. VMID protegido: shutdown y reboot también (stop ya lo estaba desde FIX-33).
		for _, action := range []string{"shutdown", "reboot", "stop"} {
			rejected(103, action, http.StatusForbidden, "INSTANCE_PROTECTED")
		}
		// 3. Acción desconocida en la ruta genérica.
		if status, body := requestJSON(t, http.MethodPost, "/instances/101/status/hibernar", adminToken, nil); status != http.StatusBadRequest || body["errorCode"] != "INVALID_ACTION" {
			t.Errorf("/status/hibernar: esperado 400 INVALID_ACTION, recibido %d: %#v", status, body)
		}
		// 4. Un usuario sin la instancia asignada: 403 sin tráfico hacia Proxmox.
		stranger := createActiveUser(t, adminToken, "op_sin_101", "OPERATOR")
		session := loginWithTOTP(t, stranger.Email, stranger.Password, stranger.Secret)
		before := stubWrites(t, 101)
		if status, body := requestJSON(t, http.MethodPost, "/instances/101/status/shutdown", session.AccessToken, nil); status != http.StatusForbidden {
			t.Errorf("OPERATOR sin la 101, shutdown: esperado 403, recibido %d: %#v", status, body)
		}
		if after := stubWrites(t, 101); after != before {
			t.Errorf("con un 403 por permisos el backend igual escribió en Proxmox")
		}
		// 5. Órdenes válidas: 202 con upid y tareaId.
		for _, order := range []struct {
			vmid   int
			action string
		}{{101, "shutdown"}, {101, "reboot"}, {101, "stop"}, {102, "start"}} {
			path := fmt.Sprintf("/instances/%d/status/%s", order.vmid, order.action)
			status, body := requestJSON(t, http.MethodPost, path, adminToken, nil)
			if status != http.StatusAccepted {
				t.Errorf("%s: esperado 202, recibido %d: %#v", path, status, body)
				continue
			}
			requiredString(t, body, "upid")
			tareaID := requiredString(t, body, "tareaId")
			waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 5*time.Second)
		}
		// 6. Entregable 4: códigos documentados en el inventario de FIX-08 y en Swagger.
		inventory := string(readFile(t, sourcePath("backend", "docs/estandar_http.md")))
		swagger := string(readFile(t, sourcePath("backend", "docs/swagger.json")))
		for _, code := range []string{"INSTANCE_INVALID_STATE", "INVALID_ACTION", "INSTANCE_PROTECTED"} {
			if !strings.Contains(inventory, code) || !strings.Contains(swagger, code) {
				t.Errorf("%s no está documentado en docs/estandar_http.md y docs/swagger.json", code)
			}
		}
	})

	t.Run("BAC-24B DELETE solo ADMIN, 409 si está encendida, 403 si está protegida, 202 con seguimiento, auditoría y TASK_FINISHED", func(t *testing.T) {
		if status, _ := requestValue(t, http.MethodDelete, "/instances/102", "", nil); status != http.StatusUnauthorized {
			t.Errorf("DELETE sin token: esperado 401, recibido %d", status)
		}
		// Un OPERATOR, aunque tenga FULL_ACCESS sobre la instancia, recibe 403.
		operatorID, _ := createUserThroughAPI(t, adminToken, fmt.Sprintf("op_delete_%d", time.Now().UnixNano()), fmt.Sprintf("op.delete.%d@elcentinela.com", time.Now().UnixNano()), "OPERATOR")
		if status, body := requestValue(t, http.MethodPut, "/admin/users/"+operatorID+"/permissions", adminToken, map[string]any{"permisos": []map[string]any{{"vmid": 102, "nivelAcceso": "FULL_ACCESS"}}}); status != http.StatusNoContent {
			t.Fatalf("precondición: FULL_ACCESS sobre la 102, esperado 204, recibido %d: %#v", status, body)
		}
		operatorToken := signedAccessToken(t, operatorID, "OPERATOR", orgID)
		before := map[int]int{101: stubDeletes(t, 101), 102: stubDeletes(t, 102), 103: stubDeletes(t, 103)}
		if status, body := requestJSON(t, http.MethodDelete, "/instances/102", operatorToken, nil); status != http.StatusForbidden {
			t.Errorf("DELETE con OPERATOR con FULL_ACCESS: esperado 403, recibido %d: %#v", status, body)
		}
		if status, body := requestJSON(t, http.MethodDelete, "/instances/101", adminToken, nil); status != http.StatusConflict || body["errorCode"] != "INSTANCE_INVALID_STATE" {
			t.Errorf("DELETE de la 101 (running): esperado 409 INSTANCE_INVALID_STATE, recibido %d: %#v", status, body)
		}
		if status, body := requestJSON(t, http.MethodDelete, "/instances/103", adminToken, nil); status != http.StatusForbidden || body["errorCode"] != "INSTANCE_PROTECTED" {
			t.Errorf("DELETE de la 103 (protegida): esperado 403 INSTANCE_PROTECTED, recibido %d: %#v", status, body)
		}
		for vmid, count := range before {
			if stubDeletes(t, vmid) != count {
				t.Errorf("el DELETE rechazado de la %d llegó igual a Proxmox", vmid)
			}
		}

		// Sobre la 102 (detenida): 202 con upid y tareaId, la orden queda auditada y llega TASK_FINISHED de un LXC.
		admin := createActiveUser(t, adminToken, "admin_delete", "ADMIN")
		adminSession := loginWithTOTP(t, admin.Email, admin.Password, admin.Secret)
		events := listenEvents(t, adminSession.AccessToken)
		status, body := requestJSON(t, http.MethodDelete, "/instances/102", adminSession.AccessToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("DELETE de la 102 (stopped) con ADMIN: esperado 202, recibido %d: %#v", status, body)
		}
		upid := requiredString(t, body, "upid")
		tareaID := requiredString(t, body, "tareaId")
		if stubDeletes(t, 102) != before[102]+1 {
			t.Errorf("el DELETE aceptado no llegó a Proxmox (DELETE /nodes/pve/lxc/102)")
		}
		audit := fmt.Sprintf("SELECT count(*) FROM auditoria WHERE detalles::text LIKE '%%%s%%';", strings.ReplaceAll(upid, "'", "''"))
		if queryDatabase(t, audit) == "0" {
			t.Errorf("la orden de borrado no quedó en auditoria con su upid")
		}
		event := waitTaskFinished(t, events, tareaID, 8*time.Second)
		if event != nil && (event["recursoTipo"] != "LXC" || event["recursoId"] != "102") {
			t.Errorf("TASK_FINISHED del borrado: esperado recursoTipo LXC y recursoId 102, recibido %v / %v", event["recursoTipo"], event["recursoId"])
		}
	})

	t.Run("BAC-25C al reiniciar el backend retoma las tareas RUNNING, informa activeTask y vence las de más de 3 minutos", func(t *testing.T) {
		// Con /tmp/tareas-en-curso el stub informa toda tarea como "running" (proxmox-stub/nginx.conf).
		if output, err := runCompose("exec", "-T", "proxmox", "touch", "/tmp/tareas-en-curso"); err != nil {
			t.Fatalf("no se pudo marcar las tareas como en curso en el stub: %v\n%s", err, output)
		}
		t.Cleanup(func() { _, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso") })

		fresh := queryDatabase(t, "SELECT uuid_generate_v7();")
		expired := queryDatabase(t, "SELECT uuid_generate_v7();")
		suffix := time.Now().UnixNano()
		queryDatabase(t, fmt.Sprintf(
			"INSERT INTO tareas_asincronas (id, usuario_id, instancia_id, accion, upid_proxmox, estado, fecha_creacion) VALUES "+
				"('%s', '%s', '101', 'shutdown', 'UPID:pve:%d:0:0:qmshutdown:101:root@pam:', 'RUNNING', NOW()), "+
				"('%s', '%s', '102', 'start', 'UPID:pve:%d:0:0:vzstart:102:root@pam:', 'RUNNING', NOW() - interval '10 minutes');",
			fresh, adminID, suffix, expired, adminID, suffix+1))

		restartBackend(t)

		// La de más de 3 minutos (D4): se consulta una vez y, como sigue en curso, queda FAILED.
		if !waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", expired), "FAILED", 10*time.Second) {
			t.Errorf("la tarea RUNNING de hace 10 minutos no quedó FAILED al reiniciar (D4: vence a los 3 minutos)")
		}
		// La reciente sigue RUNNING y GET /api/instances la informa en activeTask de la 101.
		status, list := requestJSONArray(t, http.MethodGet, "/instances", adminToken)
		if status != http.StatusOK {
			t.Fatalf("GET /instances: esperado 200, recibido %d", status)
		}
		var active any
		for _, item := range list {
			if instance, ok := item.(map[string]any); ok && instance["id"] == float64(101) {
				active = instance["activeTask"]
			}
		}
		// El contrato publicado (docs/contrato-etapa1.md) usa el id de la tarea; la tarea pide { tareaId, action, status }.
		switch value := active.(type) {
		case string:
			if value != fresh {
				t.Errorf("activeTask de la 101: esperado %s, recibido %s", fresh, value)
			}
		case map[string]any:
			if value["tareaId"] != fresh {
				t.Errorf("activeTask de la 101: esperado tareaId %s, recibido %#v", fresh, value)
			}
		default:
			t.Errorf("activeTask de la 101 con una tarea RUNNING: esperado el id de la tarea o { tareaId, action, status }, recibido %#v", active)
		}

		// Al terminar en Proxmox, el sondeo retomado la cierra y emite TASK_FINISHED.
		admin := createActiveUser(t, adminToken, "admin_reanuda", "ADMIN")
		adminSession := loginWithTOTP(t, admin.Email, admin.Password, admin.Secret)
		events := listenEvents(t, adminSession.AccessToken)
		_, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso")
		waitTaskFinished(t, events, fresh, 10*time.Second)
		if !waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", fresh), "COMPLETED", 10*time.Second) {
			t.Errorf("la tarea retomada no quedó COMPLETED en tareas_asincronas")
		}
	})
}

// stubWrites y stubDeletes cuentan, en el log de acceso del stub, las órdenes que el backend envió a Proxmox
// para un VMID (POST .../status/<acción> y DELETE .../<vmid>).
func stubWrites(t *testing.T, vmid int) int {
	return stubLogCount(t, fmt.Sprintf(`"POST /api2/json/nodes/[^/]+/(qemu|lxc)/%d/status/`, vmid))
}

func stubDeletes(t *testing.T, vmid int) int {
	return stubLogCount(t, fmt.Sprintf(`"DELETE /api2/json/nodes/[^/]+/(qemu|lxc)/%d[ ?]`, vmid))
}

func stubLogCount(t *testing.T, expression string) int {
	t.Helper()
	logs, err := runCompose("logs", "--no-color", "proxmox")
	if err != nil {
		t.Fatalf("no se pudo leer el log del stub de Proxmox: %v", err)
	}
	return len(regexp.MustCompile(expression).FindAllStringIndex(logs, -1))
}

// restartBackend reinicia el contenedor del backend y espera a que vuelva a responder.
func restartBackend(t *testing.T) {
	t.Helper()
	if output, err := runCompose("restart", "backend"); err != nil {
		t.Fatalf("no se pudo reiniciar el backend: %v\n%s", err, output)
	}
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		if response, err := (&http.Client{Timeout: 2 * time.Second}).Get(apiURL + "/version"); err == nil {
			response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("el backend no volvió a responder 60 s después de reiniciarlo")
}

func waitForValue(t *testing.T, query, want string, timeout time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if strings.TrimSpace(queryDatabase(t, query)) == want {
			return true
		}
		time.Sleep(300 * time.Millisecond)
	}
	return false
}

// listenEvents abre /api/events con un ticket del usuario y entrega cada evento JSON por el canal.
func listenEvents(t *testing.T, accessToken string) <-chan map[string]any {
	t.Helper()
	status, body := requestJSON(t, http.MethodPost, "/events/ticket", accessToken, nil)
	if status != http.StatusOK && status != http.StatusCreated {
		t.Fatalf("POST /events/ticket: esperado 200/201, recibido %d: %#v", status, body)
	}
	request, _ := http.NewRequest(http.MethodGet, apiURL+"/events?ticket="+requiredString(t, body, "ticket"), nil)
	request.Header.Set("Accept", "text/event-stream")
	response, err := (&http.Client{Timeout: 0}).Do(request)
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("GET /events: no se pudo abrir el stream (%v)", err)
	}
	t.Cleanup(func() { response.Body.Close() })
	events := make(chan map[string]any, 64)
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
	return events
}

// waitTaskFinished espera el TASK_FINISHED de la tarea indicada (detalles.tareaId).
func waitTaskFinished(t *testing.T, events <-chan map[string]any, tareaID string, timeout time.Duration) map[string]any {
	t.Helper()
	deadline := time.After(timeout)
	for {
		select {
		case event, open := <-events:
			if !open {
				t.Errorf("el stream de /api/events se cerró antes del TASK_FINISHED de la tarea %s", tareaID)
				return nil
			}
			detalles, _ := event["detalles"].(map[string]any)
			if event["tipo"] == "TASK_FINISHED" && detalles["tareaId"] == tareaID {
				return event
			}
		case <-deadline:
			t.Errorf("no llegó TASK_FINISHED de la tarea %s en %v", tareaID, timeout)
			return nil
		}
	}
}

func expectNumber(t *testing.T, body map[string]any, want, tolerance float64, path ...string) {
	t.Helper()
	value, ok := lookup(body, path...).(float64)
	if !ok {
		t.Errorf("%s: falta o no es numérico (%#v)", strings.Join(path, "."), lookup(body, path...))
		return
	}
	if math.Abs(value-want) > tolerance {
		t.Errorf("%s: esperado %v, recibido %v", strings.Join(path, "."), want, value)
	}
}

// expectGB acepta el valor en GiB (bytes / 2^30) o en GB (bytes / 1e9), con dos decimales de tolerancia.
func expectGB(t *testing.T, body map[string]any, gib float64, path ...string) {
	t.Helper()
	value, ok := lookup(body, path...).(float64)
	if !ok {
		t.Errorf("%s: falta o no es numérico (%#v)", strings.Join(path, "."), lookup(body, path...))
		return
	}
	decimal := gib * (1 << 30) / 1e9
	if math.Abs(value-gib) > 0.05 && math.Abs(value-decimal) > 0.05 {
		t.Errorf("%s: esperado %.2f (GiB) o %.2f (GB), recibido %v", strings.Join(path, "."), gib, decimal, value)
	}
}

func lookup(body map[string]any, path ...string) any {
	var current any = body
	for _, key := range path {
		object, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = object[key]
	}
	return current
}

// redisHasKeyWithTTL indica si hay alguna clave (que no sea un ticket de /api/events) con TTL entre min y max s.
func redisHasKeyWithTTL(t *testing.T, minSeconds, maxSeconds int) bool {
	t.Helper()
	keys, err := runCompose("exec", "-T", "redis", "redis-cli", "-a", "redis-test", "--no-auth-warning", "--scan")
	if err != nil {
		t.Fatalf("no se pudo listar las claves de Redis: %v\n%s", err, keys)
	}
	for _, key := range strings.Fields(keys) {
		if strings.HasPrefix(key, "ws_ticket:") {
			continue
		}
		output, _ := runCompose("exec", "-T", "redis", "redis-cli", "-a", "redis-test", "--no-auth-warning", "TTL", key)
		if ttl, err := strconv.Atoi(strings.TrimSpace(output)); err == nil && ttl >= minSeconds && ttl <= maxSeconds {
			return true
		}
	}
	return false
}

// stopProxmoxStub simula la caída de Proxmox: el stub responde 503 a todo (el backend lo traduce a 502).
// No se detiene el contenedor: con el contenedor detenido, el nombre "proxmox" deja de resolver en
// Docker y la búsqueda DNS del host podría llevar al Proxmox real. startProxmoxStub lo restablece y
// espera a que el backend vuelva a alcanzarlo; la limpieza garantiza que el resto de la suite lo tenga.
func stopProxmoxStub(t *testing.T) {
	t.Helper()
	if output, err := runCompose("exec", "-T", "proxmox", "touch", "/tmp/proxmox-caido"); err != nil {
		t.Fatalf("no se pudo simular la caída del stub de Proxmox: %v\n%s", err, output)
	}
	t.Cleanup(func() { _, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/proxmox-caido") })
}

func startProxmoxStub(t *testing.T) {
	t.Helper()
	if output, err := runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/proxmox-caido"); err != nil {
		t.Fatalf("no se pudo restablecer el stub de Proxmox: %v\n%s", err, output)
	}
	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	token := signedAccessToken(t, adminID, "ADMIN", orgID)
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if status, _ := requestValue(t, http.MethodGet, "/instances", token, nil); status == http.StatusOK {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatalf("el backend no volvió a alcanzar el stub de Proxmox 15 s después de restablecerlo")
}

// goFilesContaining busca en backend/ (sin el simulador ni vendor) los .go que coinciden con la expresión,
// separados en código y pruebas. Las rutas son relativas a backend/.
func goFilesContaining(t *testing.T, expression string) (sources, tests []string) {
	t.Helper()
	pattern := regexp.MustCompile(expression)
	root := sourcePath("backend", ".")
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() && (entry.Name() == "proxmox-simulador" || entry.Name() == "vendor" || entry.Name() == ".git") {
			return filepath.SkipDir
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil || !pattern.Match(content) {
			return nil
		}
		relative, _ := filepath.Rel(root, path)
		if strings.HasSuffix(path, "_test.go") {
			tests = append(tests, relative)
		} else {
			sources = append(sources, relative)
		}
		return nil
	})
	return sources, tests
}

// runBackendUnitTests corre `go test` en los paquetes de los archivos dados, dentro de backend/
// (sin modificarlo: -mod=readonly y la caché de compilación fuera del repositorio).
func runBackendUnitTests(t *testing.T, files []string) {
	t.Helper()
	packages := map[string]bool{}
	for _, file := range files {
		packages["./"+filepath.ToSlash(filepath.Dir(file))] = true
	}
	if len(packages) == 0 {
		t.Errorf("no hay paquetes con pruebas unitarias para correr")
		return
	}
	arguments := []string{"test", "-count=1", "-mod=readonly"}
	for name := range packages {
		arguments = append(arguments, name)
	}
	command := exec.Command("go", arguments...)
	command.Dir = sourcePath("backend", ".")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Errorf("las pruebas unitarias del backend fallan (go %s):\n%s", strings.Join(arguments, " "), output)
	}
}
