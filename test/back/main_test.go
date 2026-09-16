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
	"strings"
	"testing"
	"time"
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
		setupErr = fmt.Errorf("no se pudo iniciar el entorno: %w\n%s", err, output)
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

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
