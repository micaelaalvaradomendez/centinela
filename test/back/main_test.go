package back_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

const defaultAPIURL = "http://127.0.0.1:18080/api"

var (
	apiURL      string
	composeFile string
	projectName string
	setupErr    error
)

func TestMain(m *testing.M) {
	flag.Parse()
	apiURL = envOrDefault("BACKEND_TEST_API_URL", defaultAPIURL)
	composeFile = "compose.yaml"
	projectName = fmt.Sprintf("centinela-back-tests-%d", os.Getpid())

	if testing.Short() {
		os.Exit(m.Run())
	}

	_, _ = runCompose("down", "--volumes", "--remove-orphans")
	if output, err := runCompose("up", "--build", "--detach", "--wait"); err != nil {
		// En Docker, a veces el contenedor backend arranca en el milisegundo en que PostgreSQL reinicia tras init.sql.
		time.Sleep(2 * time.Second)
		if output2, err2 := runCompose("up", "--detach", "--wait"); err2 != nil {
			setupErr = fmt.Errorf("no se pudo iniciar el entorno: %w\n%s\nReintento:\n%s", err, output, output2)
		}
	}

	code := m.Run()
	if output, err := runCompose("down", "--volumes", "--remove-orphans"); err != nil {
		fmt.Fprintf(os.Stderr, "no se pudo limpiar el entorno de pruebas: %v\n%s\n", err, output)
		if code == 0 {
			code = 1
		}
	}
	os.Exit(code)
}

func requireIntegration(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("prueba de integración omitida con -short")
	}
	if setupErr != nil {
		t.Fatal(setupErr)
	}
}

func runCompose(arguments ...string) (string, error) {
	commandArguments := []string{"compose", "-p", projectName, "-f", composeFile}
	commandArguments = append(commandArguments, arguments...)
	command := exec.Command("docker", commandArguments...)
	output, err := command.CombinedOutput()
	return string(output), err
}

func queryDatabase(t *testing.T, query string) string {
	t.Helper()
	requireIntegration(t)
	output, err := runCompose(
		"exec", "-T", "db", "psql",
		"-U", "centinela_test", "-d", "centinela_test", "-Atc", query,
	)
	if err != nil {
		t.Fatalf("consulta PostgreSQL fallida: %v\n%s", err, output)
	}
	return strings.TrimSpace(output)
}

func requestValue(t *testing.T, method, path, token string, payload any) (int, any) {
	t.Helper()
	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("no se pudo codificar el body: %v", err)
		}
		body = bytes.NewReader(encoded)
	}

	request, err := http.NewRequest(method, apiURL+path, body)
	if err != nil {
		t.Fatalf("no se pudo crear la petición: %v", err)
	}
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("petición %s %s fallida: %v", method, path, err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("no se pudo leer la respuesta: %v", err)
	}
	var result any = map[string]any{}
	if len(responseBody) > 0 {
		if err := json.Unmarshal(responseBody, &result); err != nil {
			t.Fatalf("respuesta no JSON (%d): %s", response.StatusCode, responseBody)
		}
	}
	return response.StatusCode, result
}

func requestJSON(t *testing.T, method, path, token string, payload any) (int, map[string]any) {
	t.Helper()
	status, value := requestValue(t, method, path, token, payload)
	result, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("respuesta JSON no es un objeto (%d): %#v", status, value)
	}
	return status, result
}

func requestJSONArray(t *testing.T, method, path, token string) (int, []any) {
	t.Helper()
	status, value := requestValue(t, method, path, token, nil)
	result, ok := value.([]any)
	if !ok {
		t.Fatalf("respuesta JSON no es un array (%d): %#v", status, value)
	}
	return status, result
}

func login(t *testing.T, password string) map[string]any {
	t.Helper()
	status, body := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{
		"email":    "admin@elcentinela.com",
		"password": password,
	})
	if status != http.StatusOK {
		t.Fatalf("login esperado 200, recibido %d: %#v", status, body)
	}
	return body
}

func requiredString(t *testing.T, body map[string]any, field string) string {
	t.Helper()
	value, ok := body[field].(string)
	if !ok || value == "" {
		t.Fatalf("campo %q ausente o vacío en %#v", field, body)
	}
	return value
}

func decodeJWTClaims(t *testing.T, token string) map[string]any {
	t.Helper()
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("JWT inválido: se esperaban tres segmentos")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("payload JWT inválido: %v", err)
	}
	claims := map[string]any{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("claims JWT inválidos: %v", err)
	}
	return claims
}

// execDatabase ejecuta una sentencia y devuelve el error en lugar de abortar el test.
// Se usa para verificar que la base RECHACE una operación (por ejemplo UPDATE sobre auditoría).
func execDatabase(t *testing.T, query string) (string, error) {
	t.Helper()
	requireIntegration(t)
	output, err := runCompose(
		"exec", "-T", "db", "psql", "-v", "ON_ERROR_STOP=1",
		"-U", "centinela_test", "-d", "centinela_test", "-Atc", query,
	)
	return strings.TrimSpace(output), err
}

// rawResponse conserva cabeceras y cookies, que requestJSON descarta.
type rawResponse struct {
	Status  int
	Header  http.Header
	Cookies []*http.Cookie
	Body    map[string]any
	RawBody string
}

func (r rawResponse) cookie(nameFragment string) *http.Cookie {
	for _, c := range r.Cookies {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(nameFragment)) {
			return c
		}
	}
	return nil
}

// requestRaw envía una petición con cookies opcionales y devuelve la respuesta completa.
func requestRaw(t *testing.T, method, path, token string, payload any, cookies ...*http.Cookie) rawResponse {
	t.Helper()
	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("no se pudo codificar el body: %v", err)
		}
		body = bytes.NewReader(encoded)
	}
	request, err := http.NewRequest(method, apiURL+path, body)
	if err != nil {
		t.Fatalf("no se pudo crear la petición: %v", err)
	}
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	for _, c := range cookies {
		request.AddCookie(c)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("petición %s %s fallida: %v", method, path, err)
	}
	defer response.Body.Close()
	responseBody, _ := io.ReadAll(response.Body)
	result := rawResponse{
		Status:  response.StatusCode,
		Header:  response.Header,
		Cookies: response.Cookies(),
		Body:    map[string]any{},
		RawBody: string(responseBody),
	}
	if len(responseBody) > 0 {
		_ = json.Unmarshal(responseBody, &result.Body)
	}
	return result
}

// Marcadores que imprime backend/internal/adapters/secondary/email/mock_email.go.
const (
	mailTemporaryPassword = "clave provisoria es:"
	mailRecoveryCode      = "código de seguridad temporal es:"
)

var mailSecretPattern = regexp.MustCompile(`\*\* (\S+) \*\*`)

// mailedSecret lee la "bandeja de entrada" del MockEmailService (logs del contenedor backend)
// y devuelve el último secreto enviado a destinatario que contenga el marcador indicado.
// Así las pruebas consumen la credencial por el mismo canal que el usuario real (BAC-16).
func mailedSecret(t *testing.T, destinatario, marker string) string {
	t.Helper()
	var found string
	for attempt := 0; attempt < 10 && found == ""; attempt++ {
		output, err := runCompose("logs", "--no-color", "backend")
		if err != nil {
			t.Fatalf("no se pudieron leer los logs del backend: %v\n%s", err, output)
		}
		current := false
		for _, line := range strings.Split(output, "\n") {
			if strings.Contains(line, "Destinatario: ") {
				current = strings.Contains(line, "Destinatario: "+destinatario)
				continue
			}
			if current && strings.Contains(line, marker) {
				if match := mailSecretPattern.FindStringSubmatch(line); match != nil {
					found = match[1]
				}
				current = false
			}
		}
		if found == "" {
			time.Sleep(200 * time.Millisecond)
		}
	}
	if found == "" {
		t.Fatalf("MockEmailService no registró ningún correo %q para %s", marker, destinatario)
	}
	return found
}

// createUserThroughAPI da de alta un usuario como lo haría el panel y devuelve su id
// junto con la contraseña temporal recibida por correo.
func createUserThroughAPI(t *testing.T, adminToken, username, email, role string) (string, string) {
	t.Helper()
	status, created := requestJSON(t, http.MethodPost, "/admin/users", adminToken, map[string]any{
		"nombreCompleto": "Usuario " + username,
		"nombreUsuario":  username,
		"emailUsuario":   email,
		"rol":            role,
	})
	if status != http.StatusCreated {
		t.Fatalf("POST /api/admin/users esperado 201, recibido %d: %#v", status, created)
	}
	return requiredString(t, created, "id"), mailedSecret(t, email, mailTemporaryPassword)
}

// loginResult agrupa lo que devuelve el circuito login -> 2FA.
type loginResult struct {
	PreAuth       string
	Secret        string
	Verify        rawResponse
	AccessToken   string
	RefreshToken  string // vacío si el backend ya no lo expone en JSON (SEC-01)
	RefreshCookie *http.Cookie
}

// loginWithTOTP recorre POST /auth/login -> (GET /auth/2fa/qr si hace falta) -> POST /auth/2fa/verify.
// Si secret está vacío, enrola el 2FA y devuelve el secreto obtenido.
func loginWithTOTP(t *testing.T, email, password, secret string) loginResult {
	t.Helper()
	status, body := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{"email": email, "password": password})
	if status != http.StatusOK {
		t.Fatalf("login de %s esperado 200, recibido %d: %#v", email, status, body)
	}
	result := loginResult{PreAuth: requiredString(t, body, "jwtTemporal"), Secret: secret}
	if result.Secret == "" {
		qrStatus, qr := requestJSON(t, http.MethodGet, "/auth/2fa/qr", result.PreAuth, nil)
		if qrStatus != http.StatusOK {
			t.Fatalf("GET /auth/2fa/qr esperado 200, recibido %d: %#v", qrStatus, qr)
		}
		result.Secret = requiredString(t, qr, "secretoManual")
	}
	// El anti-replay (BAC-11) rechaza reusar el período TOTP; las pruebas encadenan logins
	// dentro de la misma ventana de 30 s, así que se libera el período antes de verificar.
	queryDatabase(t, fmt.Sprintf("UPDATE usuarios SET ultimo_totp_periodo = 0 WHERE email_usuario = '%s';", email))
	code, err := totp.GenerateCode(result.Secret, time.Now().UTC())
	if err != nil {
		t.Fatalf("no se pudo generar TOTP: %v", err)
	}
	result.Verify = requestRaw(t, http.MethodPost, "/auth/2fa/verify", result.PreAuth, map[string]any{"codigo": code})
	if result.Verify.Status != http.StatusOK {
		t.Fatalf("POST /auth/2fa/verify esperado 200, recibido %d: %s", result.Verify.Status, result.Verify.RawBody)
	}
	result.AccessToken = requiredString(t, result.Verify.Body, "accessToken")
	result.RefreshToken, _ = result.Verify.Body["refreshToken"].(string)
	result.RefreshCookie = result.Verify.cookie("refresh")
	return result
}

// activeUser es una cuenta que ya completó alta -> 2FA -> cambio obligatorio de contraseña.
type activeUser struct {
	ID       string
	Email    string
	Password string
	Secret   string
}

// createActiveUser recorre el circuito real de LOGIN-04 para obtener una cuenta utilizable
// con tokens emitidos por el backend (no firmados por la prueba).
func createActiveUser(t *testing.T, adminToken, prefix, role string) activeUser {
	t.Helper()
	suffix := time.Now().UnixNano()
	user := activeUser{Email: fmt.Sprintf("%s.%d@elcentinela.com", prefix, suffix), Password: "Activa123!"}
	var temporary string
	user.ID, temporary = createUserThroughAPI(t, adminToken, fmt.Sprintf("%s_%d", prefix, suffix), user.Email, role)
	first := loginWithTOTP(t, user.Email, temporary, "")
	user.Secret = first.Secret
	status, body := requestJSON(t, http.MethodPut, "/account/password", first.AccessToken, map[string]any{
		"contrasenaActual": temporary,
		"contrasenaNueva":  user.Password,
	})
	if status != http.StatusOK {
		t.Fatalf("PUT /account/password para %s esperado 200, recibido %d: %#v", user.Email, status, body)
	}
	return user
}

// auditCount cuenta registros de auditoría de un usuario para una acción.
func auditCount(t *testing.T, userID, accion string) int {
	t.Helper()
	var count int
	fmt.Sscan(queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM auditoria WHERE usuario_id = '%s' AND accion = '%s';", userID, accion)), &count)
	return count
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
