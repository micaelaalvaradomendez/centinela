package back_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Acceptance tests para las tareas de la Etapa 1 en backend (documentacion/actual.md):
//   - BAC-29: Contrato HTTP y de eventos de la Etapa 1 (backend/docs/contrato-etapa1.md y Swagger).
//   - FIX-40: Distinguir 504 PROXMOX_TIMEOUT de 502 PROXMOX_UNAVAILABLE (decisión D2 de etapa1.md).
//   - BAC-22: Adaptador de telemetría del nodo con caché en Redis (GET /api/node/status, RF-02).
//   - BAC-23A: Adaptador y normalización de inventario Proxmox (QEMU / LXC / IP).
//   - BAC-25A: Worker pool acotado para el seguimiento de UPID (UPID_WORKERS).

func TestBAC29_ContratoEtapa1Documentado(t *testing.T) {
	contratoPath := filepath.Clean(sourcePath("backend", "docs/contrato-etapa1.md"))
	content, err := os.ReadFile(contratoPath)
	if err != nil {
		t.Fatalf("BAC-29 no implementada: falta el archivo de contrato %s: %v", contratoPath, err)
	}
	text := string(content)

	t.Run("documenta esquema de GET /api/node/status", func(t *testing.T) {
		for _, required := range []string{"/api/node/status", "usagePercent", "cores", "usedGb", "totalGb", "uptimeSeconds", "instancesSummary", "stale", "fetchedAt"} {
			if !strings.Contains(text, required) {
				t.Errorf("contrato-etapa1.md no documenta el campo o ruta %q de GET /api/node/status", required)
			}
		}
	})

	t.Run("documenta campos extendidos de GET /api/instances", func(t *testing.T) {
		for _, required := range []string{"/api/instances", "ip", "cpuUsage", "ramUsage", "maxRam", "nivelAcceso", "activeTask"} {
			if !strings.Contains(text, required) {
				t.Errorf("contrato-etapa1.md no documenta el campo o ruta %q de GET /api/instances", required)
			}
		}
	})

	t.Run("documenta rutas de energía y respuesta 202 con upid y tareaId", func(t *testing.T) {
		for _, required := range []string{"start", "shutdown", "stop", "reboot", "upid", "tareaId", "202"} {
			if !strings.Contains(text, required) {
				t.Errorf("contrato-etapa1.md no documenta la acción o respuesta esperada %q", required)
			}
		}
	})

	t.Run("documenta formato y detalles de TASK_FINISHED", func(t *testing.T) {
		for _, required := range []string{"TASK_FINISHED", "COMPLETED", "FAILED"} {
			if !strings.Contains(text, required) {
				t.Errorf("contrato-etapa1.md no documenta el evento o estado %q", required)
			}
		}
	})

	t.Run("documenta códigos de error normalizados de la Etapa 1", func(t *testing.T) {
		for _, code := range []string{
			"INSTANCE_ACCESS_DENIED",
			"INSTANCE_PROTECTED",
			"INSTANCE_BUSY",
			"INSTANCE_INVALID_STATE",
			"INSTANCE_NOT_FOUND",
			"PROXMOX_UNAVAILABLE",
			"PROXMOX_TIMEOUT",
			"INVALID_ACTION",
		} {
			if !strings.Contains(text, code) {
				t.Errorf("contrato-etapa1.md no documenta el código de error %q", code)
			}
		}
	})
}

func TestFIX40_DistinguirProxmoxTimeoutDeUnavailable(t *testing.T) {
	handlerPath := filepath.Clean(sourcePath("backend", "internal/adapters/primary/http/instance_handler.go"))
	content, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("no se pudo leer instance_handler.go: %v", err)
	}
	text := string(content)

	// Criterio de éxito FIX-40:
	// 1. ErrProxmoxTimeout mapea a 504 con PROXMOX_TIMEOUT
	// 2. ErrProxmoxNoDisponible / ErrProxmoxCredenciales mapean a 502 con PROXMOX_UNAVAILABLE
	if !strings.Contains(text, "PROXMOX_TIMEOUT") {
		t.Errorf("FIX-40 no implementado: instance_handler.go aún no define ni devuelve errorCode PROXMOX_TIMEOUT")
	}

	reTimeout := regexp.MustCompile(`(?s)errors\.Is\(err,\s*ports\.ErrProxmoxTimeout\).*?SendError\([^,]+,\s*(?:http\.)?StatusGatewayTimeout,\s*"([^"]+)"`)
	match := reTimeout.FindStringSubmatch(text)
	if match == nil {
		t.Errorf("FIX-40: no se encontró mapeo de ErrProxmoxTimeout a StatusGatewayTimeout en instance_handler.go")
	} else if match[1] != "PROXMOX_TIMEOUT" {
		t.Errorf("FIX-40: ErrProxmoxTimeout mapea a %q, esperado \"PROXMOX_TIMEOUT\"", match[1])
	}

	reUnavailable := regexp.MustCompile(`(?s)errors\.Is\(err,\s*ports\.ErrProxmoxNoDisponible\).*?SendError\([^,]+,\s*(?:http\.)?StatusBadGateway,\s*"([^"]+)"`)
	matchUn := reUnavailable.FindStringSubmatch(text)
	if matchUn == nil || matchUn[1] != "PROXMOX_UNAVAILABLE" {
		t.Errorf("FIX-40: ErrProxmoxNoDisponible debe mantener StatusBadGateway con PROXMOX_UNAVAILABLE")
	}
}

func TestBAC22_EstadoNodoTelemetryEndpoint(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	suffix := time.Now().UnixNano()
	operatorID, _ := createUserThroughAPI(t, adminToken, fmt.Sprintf("op_telemetry_%d", suffix), fmt.Sprintf("op.telem.%d@elcentinela.com", suffix), "OPERATOR")
	operatorToken := signedAccessToken(t, operatorID, "OPERATOR", orgID)

	t.Run("sin autenticación responde 401 UNAUTHORIZED", func(t *testing.T) {
		status, _ := requestValue(t, http.MethodGet, "/node/status", "", nil)
		if status != http.StatusUnauthorized {
			t.Errorf("GET /node/status sin token esperado 401, recibido %d", status)
		}
	})

	t.Run("con rol OPERATOR o ADMIN responde 200 con telemetría normalizada", func(t *testing.T) {
		status, body := requestValue(t, http.MethodGet, "/node/status", operatorToken, nil)
		if status != http.StatusOK {
			t.Fatalf("BAC-22 no implementada: GET /node/status esperado 200, recibido %d: %#v", status, body)
		}

		data, ok := body.(map[string]any)
		if !ok {
			t.Fatalf("respuesta no es un objeto JSON: %#v", body)
		}

		// Validar campos de telemetría según BAC-29 / BAC-22
		for _, section := range []string{"cpu", "ram", "storage", "uptimeSeconds", "instancesSummary", "stale", "fetchedAt"} {
			if _, exists := data[section]; !exists {
				t.Errorf("campo %q faltante en GET /node/status: %#v", section, data)
			}
		}

		if cpu, ok := data["cpu"].(map[string]any); ok {
			if _, hasUsage := cpu["usagePercent"]; !hasUsage {
				t.Errorf("cpu.usagePercent faltante en telemetría: %#v", cpu)
			}
			if _, hasCores := cpu["cores"]; !hasCores {
				t.Errorf("cpu.cores faltante en telemetría: %#v", cpu)
			}
		} else {
			t.Errorf("sección cpu no es objeto: %#v", data["cpu"])
		}

		if ram, ok := data["ram"].(map[string]any); ok {
			for _, f := range []string{"usedGb", "totalGb", "usagePercent"} {
				if _, has := ram[f]; !has {
					t.Errorf("ram.%s faltante en telemetría: %#v", f, ram)
				}
			}
		} else {
			t.Errorf("sección ram no es objeto: %#v", data["ram"])
		}
	})
}

func TestBAC23A_NormalizacionInventarioQemuLxc(t *testing.T) {
	clientPath := filepath.Clean(sourcePath("backend", "internal/adapters/secondary/proxmox/client.go"))
	content, err := os.ReadFile(clientPath)
	if err != nil {
		t.Fatalf("no se pudo leer client.go de Proxmox: %v", err)
	}
	text := string(content)

	// BAC-23A exige consultar interfaces de qemu y lxc para resolver la IP
	hasQemuAgent := strings.Contains(text, "network-get-interfaces") || strings.Contains(text, "agent/network-get-interfaces")
	hasLxcInterfaces := strings.Contains(text, "/interfaces") || strings.Contains(text, "lxc")

	if !hasQemuAgent || !hasLxcInterfaces {
		t.Errorf("BAC-23A no implementada: proxmox/client.go aún no implementa resolución de interfaces para QEMU y LXC")
	}
}

func TestBAC25A_WorkerPoolAcotadoParaUPID(t *testing.T) {
	servicePath := filepath.Clean(sourcePath("backend", "internal/core/services/seguimiento_tareas.go"))
	content, err := os.ReadFile(servicePath)
	if err != nil {
		t.Fatalf("no se pudo leer seguimiento_tareas.go: %v", err)
	}
	text := string(content)

	// BAC-25A exige un worker pool acotado (UPID_WORKERS / canal con buffer / pool de workers)
	hasWorkers := strings.Contains(text, "UPID_WORKERS") || strings.Contains(text, "workers") || strings.Contains(text, "workerPool") || strings.Contains(text, "WorkerPool")
	if !hasWorkers {
		t.Errorf("BAC-25A no implementada: seguimiento_tareas.go aún no define un pool acotado de workers para seguimiento de UPID")
	}
}
