package back_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestLOGIN04IntegracionFrontBack revalida en cada corrida la prueba integral de LOGIN-04
// ejecutando el CÓDIGO REAL DEL FRONTEND (test/front/login04-e2e.test.ts) contra el backend
// real de este stack. Complementa TestHitoLOGIN04CircuitoCompletoSeguridadYAutenticacion
// (que arma las peticiones desde Go): acá las peticiones las arma el frontend, así que detecta
// desalineaciones de contrato entre ambos lados (headers, cookies, payloads, códigos de error).
func TestLOGIN04IntegracionFrontBack(t *testing.T) {
	requireIntegration(t)
	if _, err := exec.LookPath("pnpm"); err != nil {
		t.Skip("pnpm no está instalado: no se puede ejecutar el frontend")
	}
	frontDir, _ := filepath.Abs("../front")
	if _, err := os.Stat(filepath.Join(frontDir, "node_modules")); err != nil {
		t.Skip("falta test/front/node_modules: correr `pnpm --dir test/front install`")
	}
	composePath, _ := filepath.Abs(composeFile)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	command := exec.Command("pnpm", "exec", "vitest", "run", "login04-e2e", "--reporter=verbose")
	command.Dir = frontDir
	command.Env = append(os.Environ(),
		"CENTINELA_E2E_API="+strings.TrimSuffix(apiURL, "/api"),
		"CENTINELA_E2E_ADMIN_TOKEN="+adminToken,
		"CENTINELA_E2E_COMPOSE=docker compose -p "+projectName+" -f "+composePath,
		"NO_COLOR=1", "FORCE_COLOR=0",
	)
	output, runErr := command.CombinedOutput()
	text := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(string(output), "")

	summary := regexp.MustCompile(`(?m)^\s*Tests\s+(\d.*)$`).FindStringSubmatch(text)
	if summary == nil {
		t.Fatalf("no se pudo ejecutar la suite del frontend:\n%s", tail(text, 60))
	}
	if strings.Contains(summary[1], "skipped") && !strings.Contains(summary[1], "passed") && !strings.Contains(summary[1], "failed") {
		t.Fatalf("la suite LOGIN-04 del frontend se omitió; no recibió las variables CENTINELA_E2E_*: %s", summary[1])
	}
	failedSteps := 0
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "✓") || strings.HasPrefix(trimmed, "×") {
			t.Log(trimmed)
		}
		if strings.HasPrefix(trimmed, "×") {
			failedSteps++
		}
	}
	if runErr != nil || failedSteps > 0 || strings.Contains(summary[1], "failed") {
		var failures []string
		for _, block := range regexp.MustCompile(`(?m)^ FAIL .*\n(?:.*\n){0,8}`).FindAllString(text, -1) {
			failures = append(failures, strings.TrimSpace(block))
		}
		t.Errorf("LOGIN-04 front ↔ back: %s\n\n%s", summary[1], strings.Join(failures, "\n\n"))
	}
}

func tail(text string, lines int) string {
	all := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(all) > lines {
		all = all[len(all)-lines:]
	}
	return strings.Join(all, "\n")
}
