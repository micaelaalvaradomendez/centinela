package back_test

import (
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
			ids = append(ids, "'"+requiredString(t, body, "tareaId")+"'")
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
