package back_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Tareas de cierre de la fase base (documentacion/actual.md):
//   - INF-05  CORS con lista blanca (el TLS de Nginx vive en el servidor y no está versionado)
//   - INF-06  Redis en backend/docker-compose.yml con contraseña y límite de memoria
//   - INF-08  Variables SMTP documentadas para el correo saliente
//   - BAC-16B Adaptador SmtpEmailService seleccionado por EMAIL_PROVIDER
//   - BAC-17B 1 sesión de usuario = 1 registro activo en sesiones_activas
//   - BAC-18B Índice parcial en sesiones_activas y particionamiento trimestral de auditoria
func TestCierreFaseBase(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	t.Run("INF-05 CORS acepta solo los origenes de la lista blanca", func(t *testing.T) {
		const allowed = "http://front.centinela.test" // ALLOWED_ORIGINS en compose.yaml
		preflight := func(origin string) *http.Response {
			request, _ := http.NewRequest(http.MethodOptions, apiURL+"/auth/login", nil)
			request.Header.Set("Origin", origin)
			request.Header.Set("Access-Control-Request-Method", http.MethodPost)
			request.Header.Set("Access-Control-Request-Headers", "Content-Type")
			response, err := (&http.Client{Timeout: 10 * time.Second}).Do(request)
			if err != nil {
				t.Fatalf("preflight desde %s fallido: %v", origin, err)
			}
			response.Body.Close()
			return response
		}

		ok := preflight(allowed)
		if ok.StatusCode != http.StatusNoContent {
			t.Errorf("preflight desde el origen permitido: esperado 204, recibido %d", ok.StatusCode)
		}
		if got := ok.Header.Get("Access-Control-Allow-Origin"); got != allowed {
			t.Errorf("Access-Control-Allow-Origin esperado %q (nunca *), recibido %q", allowed, got)
		}
		if ok.Header.Get("Access-Control-Allow-Credentials") != "true" {
			t.Errorf("el origen permitido debe recibir Access-Control-Allow-Credentials: true (cookie HttpOnly de SEC-01)")
		}

		denied := preflight("http://atacante.example")
		if denied.StatusCode != http.StatusForbidden {
			t.Errorf("preflight desde un origen no autorizado: esperado 403, recibido %d", denied.StatusCode)
		}
		if got := denied.Header.Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("un origen no autorizado no debe recibir Access-Control-Allow-Origin, recibió %q", got)
		}

		request, _ := http.NewRequest(http.MethodGet, apiURL+"/version", nil)
		request.Header.Set("Origin", "http://atacante.example")
		response, err := (&http.Client{Timeout: 10 * time.Second}).Do(request)
		if err != nil {
			t.Fatalf("GET /version con Origin no autorizado fallido: %v", err)
		}
		response.Body.Close()
		if got := response.Header.Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("una petición simple desde un origen no autorizado no debe recibir Access-Control-Allow-Origin, recibió %q", got)
		}
	})

	t.Run("FIX-31 INF-05 TLS en el borde con Nginx", func(t *testing.T) {
		// El criterio ("el tráfico se sirve únicamente sobre HTTPS") depende de la
		// configuración de Nginx del contenedor CT 103, que no está versionada en
		// ningún repositorio: el deploy (frontend/.github/workflows/deploy-front-test.yml)
		// solo ejecuta `nginx -t` y `reload` sobre la configuración del servidor.
		t.Skip("INF-05 TLS no verificable automáticamente: la configuración de Nginx del servidor no está versionada; requiere verificación manual o versionar el archivo")
	})

	t.Run("INF-06 Redis con contraseña y limite de memoria en docker-compose del backend", func(t *testing.T) {
		composePath := sourcePath("backend", "docker-compose.yml")
		env := string(readFile(t, sourcePath("backend", ".env.example")))
		for _, variable := range []string{"REDIS_ADDR=", "REDIS_PASSWORD="} {
			if !strings.Contains(env, variable) {
				t.Errorf("backend/.env.example no documenta %s para el backend", strings.TrimSuffix(variable, "="))
			}
		}

		// Se levanta únicamente el servicio redis del compose del equipo para verificarlo en vivo.
		project := fmt.Sprintf("centinela-inf06-%d", os.Getpid())
		password := "inf06-test-pass"
		compose := func(arguments ...string) (string, error) {
			command := exec.Command("docker", append([]string{"compose", "-p", project, "-f", composePath}, arguments...)...)
			command.Env = append(os.Environ(), "REDIS_PASSWORD="+password)
			output, err := command.CombinedOutput()
			return string(output), err
		}
		if output, err := compose("up", "--detach", "--wait", "redis"); err != nil {
			t.Fatalf("no se pudo levantar el servicio redis de backend/docker-compose.yml: %v\n%s", err, output)
		}
		t.Cleanup(func() { _, _ = compose("down", "--volumes", "--remove-orphans") })

		redis := func(arguments ...string) string {
			output, _ := compose(append([]string{"exec", "-T", "redis", "redis-cli"}, arguments...)...)
			return strings.TrimSpace(output)
		}
		if got := redis("-a", password, "--no-auth-warning", "PING"); got != "PONG" {
			t.Errorf("redis-cli PING con contraseña: esperado PONG, recibido %q", got)
		}
		if got := redis("PING"); !strings.Contains(got, "NOAUTH") {
			t.Errorf("redis-cli PING sin contraseña debe rechazarse con NOAUTH, recibido %q", got)
		}
		if got := redis("-a", password, "--no-auth-warning", "CONFIG", "GET", "maxmemory"); !strings.Contains(got, "268435456") {
			t.Errorf("maxmemory esperado 256mb (268435456), recibido %q", got)
		}
		if got := redis("-a", password, "--no-auth-warning", "CONFIG", "GET", "maxmemory-policy"); !strings.Contains(got, "volatile-lru") {
			t.Errorf("maxmemory-policy esperado volatile-lru, recibido %q", got)
		}
	})

	t.Run("INF-08 variables SMTP documentadas para el correo saliente", func(t *testing.T) {
		env := string(readFile(t, sourcePath("backend", ".env.example")))
		compose := string(readFile(t, sourcePath("backend", "docker-compose.yml")))
		for _, variable := range []string{"EMAIL_PROVIDER", "SMTP_HOST", "SMTP_PORT", "SMTP_USER", "SMTP_PASS", "SMTP_FROM"} {
			if !strings.Contains(env, variable+"=") && !strings.Contains(compose, variable) {
				t.Errorf("%s no está documentada en backend/.env.example ni en backend/docker-compose.yml", variable)
			}
		}
	})

	t.Run("BAC-16B con EMAIL_PROVIDER=smtp las credenciales y el codigo OTP llegan por SMTP", func(t *testing.T) {
		useBackend(t, envOrDefault("BACKEND_TEST_SMTP_API_URL", "http://127.0.0.1:18081/api"))
		email := fmt.Sprintf("smtp.%d@elcentinela.com", time.Now().UnixNano())

		status, created := requestJSON(t, http.MethodPost, "/admin/users", adminToken, map[string]any{
			"nombreCompleto": "Usuario SMTP", "nombreUsuario": fmt.Sprintf("smtp_%d", time.Now().UnixNano()),
			"emailUsuario": email, "rol": "OPERATOR",
		})
		if status != http.StatusCreated {
			t.Fatalf("POST /api/admin/users (EMAIL_PROVIDER=smtp) esperado 201, recibido %d: %#v", status, created)
		}
		body := waitForMail(t, email, "")
		if body == "" {
			t.Fatalf("con EMAIL_PROVIDER=smtp el alta no envió ningún correo SMTP a %s (¿se sigue usando MockEmailService?)", email)
		}
		temporary := ""
		for _, candidate := range passwordCandidates(body) {
			if loginStatusAt(t, email, candidate) == http.StatusOK {
				temporary = candidate
				break
			}
		}
		if temporary == "" {
			t.Errorf("el correo de alta llegó pero ninguna palabra del cuerpo funciona como contraseña temporal:\n%s", body)
		}

		if status, _ := requestJSON(t, http.MethodPost, "/auth/password/forgot", "", map[string]any{"email": email}); status != http.StatusOK {
			t.Fatalf("POST /auth/password/forgot esperado 200, recibido %d", status)
		}
		if code := regexp.MustCompile(`\b\d{6}\b`).FindString(waitForMail(t, email, body)); code == "" {
			t.Errorf("el código de recuperación de 6 dígitos no llegó por SMTP a %s", email)
		}

		// Ante un fallo de entrega la operación se aborta con 502 y no se persiste el usuario.
		if output, err := runCompose("stop", "mailpit"); err != nil {
			t.Fatalf("no se pudo detener mailpit: %v\n%s", err, output)
		}
		t.Cleanup(func() { _, _ = runCompose("start", "mailpit") })
		failing := fmt.Sprintf("smtp.falla.%d@elcentinela.com", time.Now().UnixNano())
		status, response := requestJSON(t, http.MethodPost, "/admin/users", adminToken, map[string]any{
			"nombreCompleto": "Usuario Sin Correo", "nombreUsuario": fmt.Sprintf("smtp_falla_%d", time.Now().UnixNano()),
			"emailUsuario": failing, "rol": "OPERATOR",
		})
		if status != http.StatusBadGateway || response["errorCode"] != "EMAIL_DELIVERY_FAILED" {
			t.Errorf("con el servidor SMTP caído el alta debe responder 502 EMAIL_DELIVERY_FAILED, recibido %d: %#v", status, response)
		}
		if rows := queryDatabase(t, "SELECT count(*) FROM usuarios WHERE email_usuario = '"+failing+"';"); rows != "0" {
			t.Errorf("con el correo fallido el usuario no debe persistirse; se encontraron %s filas", rows)
		}
	})

	t.Run("BAC-17B una sesion de usuario ocupa un unico registro en sesiones_activas", func(t *testing.T) {
		user := createActiveUser(t, adminToken, "sesion_unica", "OPERATOR")
		since := queryDatabase(t, "SELECT now();")

		session := loginWithTOTP(t, user.Email, user.Password, user.Secret)
		cookie := session.RefreshCookie
		if cookie == nil {
			t.Fatalf("2fa/verify no emitió la cookie de refresh (SEC-01)")
		}
		for i := 0; i < 10; i++ {
			refresh := requestRaw(t, http.MethodPost, "/auth/refresh", "", nil, cookie)
			if refresh.Status != http.StatusOK {
				t.Fatalf("refresh %d: esperado 200, recibido %d: %s", i+1, refresh.Status, refresh.RawBody)
			}
			if next := refresh.cookie("refresh"); next != nil && next.Value != "" {
				cookie = next
			}
			session.AccessToken = requiredString(t, refresh.Body, "accessToken")
		}

		created := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM sesiones_activas WHERE usuario_id = '%s' AND fecha_creacion >= '%s';", user.ID, since))
		active := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM sesiones_activas WHERE usuario_id = '%s' AND activa = true AND fecha_creacion >= '%s';", user.ID, since))
		if created != "1" || active != "1" {
			t.Errorf("login + 2FA + 10 refresh debe dejar exactamente 1 registro (activo); se crearon %s filas y hay %s activas", created, active)
		}
		if lastAccess := queryDatabase(t, fmt.Sprintf("SELECT coalesce(fecha_ultimo_acceso >= '%s', false) FROM usuarios WHERE id = '%s';", since, user.ID)); lastAccess != "t" {
			t.Errorf("fecha_ultimo_acceso debe actualizarse al completar el 2FA")
		}

		if logout := requestRaw(t, http.MethodPost, "/auth/logout", session.AccessToken, nil, cookie); logout.Status != http.StatusNoContent {
			t.Fatalf("logout esperado 204, recibido %d: %s", logout.Status, logout.RawBody)
		}
		if left := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM sesiones_activas WHERE usuario_id = '%s' AND activa = true AND fecha_creacion >= '%s';", user.ID, since)); left != "0" {
			t.Errorf("tras el logout quedaron %s registros activos de esa sesión", left)
		}
	})

	t.Run("BAC-18B indice parcial en sesiones_activas y auditoria particionada por trimestre", func(t *testing.T) {
		index := queryDatabase(t, "SELECT coalesce(max(indexdef), '') FROM pg_indexes WHERE tablename = 'sesiones_activas' AND indexname = 'idx_sesiones_activas_vigentes';")
		if index == "" {
			t.Errorf("falta el índice parcial idx_sesiones_activas_vigentes en sesiones_activas")
		} else if !strings.Contains(index, "WHERE (activa = true)") || !strings.Contains(index, "jti_token") {
			t.Errorf("idx_sesiones_activas_vigentes debe ser (jti_token, usuario_id) WHERE activa = true; definición: %s", index)
		}

		if kind := queryDatabase(t, "SELECT relkind FROM pg_class WHERE relname = 'auditoria';"); kind != "p" {
			t.Errorf("auditoria debe ser una tabla particionada (PARTITION BY RANGE (fecha_hora)); relkind actual %q", kind)
		}
		partitions := queryDatabase(t, "SELECT coalesce(string_agg(c.relname, ','), '') FROM pg_inherits i JOIN pg_class c ON c.oid = i.inhrelid JOIN pg_class p ON p.oid = i.inhparent WHERE p.relname = 'auditoria';")
		for _, name := range []string{"auditoria_2026_q3", "auditoria_2026_q4", "auditoria_2027_q1", "auditoria_default"} {
			if !strings.Contains(partitions, name) {
				t.Errorf("falta la partición %s de auditoria (particiones actuales: %q)", name, partitions)
			}
		}
		indexes := queryDatabase(t, "SELECT coalesce(string_agg(indexdef, E'\\n'), '') FROM pg_indexes WHERE tablename = 'auditoria';")
		for _, columns := range []string{"(fecha_hora DESC, accion, resultado)", "(usuario_id, fecha_hora DESC)"} {
			if !strings.Contains(indexes, columns) {
				t.Errorf("falta el índice compuesto %s en auditoria", columns)
			}
		}
		// La purga horaria (DELETE de sesiones inactivas o vencidas) y el tiempo de consulta
		// < 20 ms no se verifican aquí: requieren esperar el ticker o un volumen de datos real.
	})
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("no se pudo leer %s: %v", path, err)
	}
	return content
}

// useBackend dirige las peticiones de la prueba a otra instancia del backend.
func useBackend(t *testing.T, url string) {
	t.Helper()
	previous := apiURL
	apiURL = url
	t.Cleanup(func() { apiURL = previous })
}

func loginStatusAt(t *testing.T, email, password string) int {
	t.Helper()
	status, _ := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{"email": email, "password": password})
	return status
}

// waitForMail consulta la API de Mailpit y devuelve el cuerpo del último correo a destinatario
// distinto de previous. Devuelve "" si no llega ninguno en 10 segundos.
func waitForMail(t *testing.T, destinatario, previous string) string {
	t.Helper()
	base := envOrDefault("BACKEND_TEST_MAILPIT_URL", "http://127.0.0.1:18025")
	client := &http.Client{Timeout: 5 * time.Second}
	for attempt := 0; attempt < 20; attempt++ {
		response, err := client.Get(base + "/api/v1/search?query=" + "to:" + destinatario)
		if err == nil {
			var search struct {
				Messages []struct {
					ID string `json:"ID"`
				} `json:"messages"`
			}
			_ = json.NewDecoder(response.Body).Decode(&search)
			response.Body.Close()
			if len(search.Messages) > 0 {
				message, err := client.Get(base + "/api/v1/message/" + search.Messages[0].ID)
				if err == nil {
					var detail struct {
						Text string `json:"Text"`
						HTML string `json:"HTML"`
					}
					raw, _ := io.ReadAll(message.Body)
					message.Body.Close()
					_ = json.Unmarshal(raw, &detail)
					body := detail.Text + "\n" + detail.HTML
					if body != previous {
						return body
					}
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return ""
}

// passwordCandidates devuelve las palabras del correo que podrían ser la clave temporal
// (el formato del cuerpo lo define BAC-16B; la prueba no lo asume).
func passwordCandidates(body string) []string {
	cleaner := strings.NewReplacer("<", " ", ">", " ", "*", " ", "\"", " ", "'", " ", "`", " ")
	var candidates []string
	for _, word := range strings.Fields(cleaner.Replace(body)) {
		word = strings.Trim(word, ".,:;()[]")
		if len(word) >= 8 && len(word) <= 16 && !strings.Contains(word, "@") && !strings.Contains(word, "=") {
			candidates = append(candidates, word)
		}
	}
	return candidates
}
