package back_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Tareas de la Etapa 1 Ola 2 en documentacion/actual.md:
//   - BAC-23B: IP real en GET /api/instances y verificación del filtrado RBAC (RF-03).
//   - BAC-22B: Conteo de instancias por estado en GET /api/node/status (RF-02) y métricas por instancia (RF-03).
//   - BAC-25B: Timeout configurable (UPID_TIMEOUT 3m), reintentos con backoff y exitstatus en TASK_FINISHED.
func TestEtapa1Ola2(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	suffix := time.Now().UnixNano()
	operatorID, _ := createUserThroughAPI(t, adminToken, fmt.Sprintf("op_ola2_%d", suffix), fmt.Sprintf("op.ola2.%d@elcentinela.com", suffix), "OPERATOR")
	operatorToken := signedAccessToken(t, operatorID, "OPERATOR", orgID)

	t.Run("BAC-23B IP real en GET /api/instances y filtrado RBAC estricto", func(t *testing.T) {
		// 1. Sin token -> 401
		if status, _ := requestValue(t, http.MethodGet, "/instances", "", nil); status != http.StatusUnauthorized {
			t.Errorf("GET /instances sin token: esperado 401, recibido %d", status)
		}

		// 2. Con OPERATOR sin permisos -> lista vacía
		statusEmpty, listEmpty := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if statusEmpty != http.StatusOK || len(listEmpty) != 0 {
			t.Errorf("GET /instances con OPERATOR sin permisos: esperado 200 con lista vacía, recibido %d con %d elementos", statusEmpty, len(listEmpty))
		}

		// 3. Asignar permiso solo sobre 101 al OPERATOR
		permEndpoint := "/admin/users/" + operatorID + "/permissions"
		if status, body := requestValue(t, http.MethodPut, permEndpoint, adminToken, map[string]any{"permisos": []map[string]any{
			{"vmid": 101, "nivelAcceso": "FULL_ACCESS"},
		}}); status != http.StatusNoContent {
			t.Fatalf("asignar permiso sobre 101 a OPERATOR falló con status %d: %#v", status, body)
		}

		// 4. Con permiso sobre 101, OPERATOR solo ve 101 y ve su IP resuelta (o null si apagada)
		status, list := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if status != http.StatusOK || len(list) != 1 {
			t.Fatalf("GET /instances con OPERATOR con permiso en 101: esperado 1 instancia, recibido %d: %#v", status, list)
		}
		instancia101 := list[0].(map[string]any)
		if id, ok := instancia101["id"].(float64); !ok || id != 101 {
			t.Errorf("instancia devuelta debe ser ID 101, recibido %v", instancia101["id"])
		}

		// 5. La IP debe estar presente como campo string o null
		if _, ok := instancia101["ip"]; !ok {
			t.Errorf("la respuesta de GET /instances no contiene el campo 'ip'")
		}

		// 6. Análisis estático de RBAC: ResolverIPs debe invocarse solo DESPUÉS del filtro de instancias del usuario
		handlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/instance_handler.go")))
		idxFilter := strings.Index(handlerCode, "instancias = filtradas")
		idxResolver := strings.Index(handlerCode, "ResolverIPs")
		if idxResolver != -1 && idxFilter != -1 && idxResolver < idxFilter {
			t.Errorf("instance_handler.go invoca ResolverIPs ANTES de filtrar por permisos del usuario (fuga de consultas a máquinas ajenas)")
		}
	})

	t.Run("BAC-22B conteo de instancias en GET /api/node/status y telemetria por instancia", func(t *testing.T) {
		// 1. GET /api/node/status debe incluir instancesSummary con vms y lxc
		status, body := requestJSON(t, http.MethodGet, "/node/status", adminToken, nil)
		if status != http.StatusOK {
			t.Fatalf("GET /api/node/status falló con %d: %#v", status, body)
		}

		summaryRaw, ok := body["instancesSummary"]
		if !ok || summaryRaw == nil {
			t.Fatalf("GET /api/node/status no incluye el campo 'instancesSummary'")
		}
		summary, ok := summaryRaw.(map[string]any)
		if !ok {
			t.Fatalf("instancesSummary no es un objeto: %#v", summaryRaw)
		}

		vms, okVMs := summary["vms"].(map[string]any)
		lxc, okLXC := summary["lxc"].(map[string]any)
		if !okVMs || !okLXC {
			t.Fatalf("instancesSummary debe contener 'vms' y 'lxc', recibido: %#v", summary)
		}

		for _, campo := range []string{"running", "stopped", "paused", "total"} {
			if _, ok := vms[campo]; !ok {
				t.Errorf("instancesSummary.vms no contiene el campo %q", campo)
			}
			if _, ok := lxc[campo]; !ok {
				t.Errorf("instancesSummary.lxc no contiene el campo %q", campo)
			}
		}

		// En el stub (101 running VM, 102 stopped LXC, 103 running VM):
		if totalVMs, ok := vms["total"].(float64); !ok || totalVMs < 2 {
			t.Errorf("instancesSummary.vms.total esperado >= 2, recibido %v", vms["total"])
		}
		if totalLXC, ok := lxc["total"].(float64); !ok || totalLXC < 1 {
			t.Errorf("instancesSummary.lxc.total esperado >= 1, recibido %v", lxc["total"])
		}

		// 2. GET /api/instances devuelve cpuUsage, ramUsage y maxRam
		_, list := requestJSONArray(t, http.MethodGet, "/instances", adminToken)
		var foundRunning bool
		for _, item := range list {
			inst := item.(map[string]any)
			if inst["status"] == "running" {
				foundRunning = true
				if _, ok := inst["cpuUsage"]; !ok {
					t.Errorf("instancia %v (running) no incluye cpuUsage", inst["id"])
				}
				if _, ok := inst["ramUsage"]; !ok {
					t.Errorf("instancia %v (running) no incluye ramUsage", inst["id"])
				}
				if _, ok := inst["maxRam"]; !ok {
					t.Errorf("instancia %v no incluye maxRam", inst["id"])
				}
			}
		}
		if !foundRunning {
			t.Errorf("no se encontraron instancias en ejecución para verificar métricas")
		}
	})

	t.Run("BAC-25B timeout configurable UPID_TIMEOUT y detalles en TASK_FINISHED", func(t *testing.T) {
		seguimientoCode := string(readFile(t, sourcePath("backend", "internal/core/services/seguimiento_tareas.go")))

		// 1. D4: Cortar seguimiento a los 3 minutos (UPID_TIMEOUT con default 3m en vez de 10 min fijos).
		usesTimeoutConfig := strings.Contains(seguimientoCode, "UPID_TIMEOUT") ||
			strings.Contains(seguimientoCode, "3*time.Minute") ||
			strings.Contains(seguimientoCode, "3 * time.Minute")

		hasFixed10Minutes := regexp.MustCompile(`limiteSeguimiento\s*=\s*10\s*\*\s*time\.Minute`).MatchString(seguimientoCode)
		if hasFixed10Minutes && !usesTimeoutConfig {
			t.Errorf("BAC-25B pendiente: seguimiento_tareas.go conserva limiteSeguimiento fijo en 10 minutos sin soportar UPID_TIMEOUT (3m)")
		}

		// 2. Reintentos con retroceso o backoff ante fallas temporales
		hasBackoff := strings.Contains(seguimientoCode, "backoff") ||
			strings.Contains(seguimientoCode, "reintento") ||
			strings.Contains(seguimientoCode, "time.Sleep")

		if !hasBackoff {
			t.Errorf("seguimiento_tareas.go no implementa estrategia de reintento o retroceso")
		}

		// 3. D2 y BAC-29: TASK_FINISHED con exitstatus y motivo cuando falla
		eventsCode := string(readFile(t, sourcePath("backend", "internal/core/ports/event_port.go")))
		if !strings.Contains(eventsCode, "MotivoProxmoxError") || !strings.Contains(eventsCode, "PROXMOX_ERROR") {
			t.Errorf("ports/event_port.go no incluye constantes de motivo de error (PROXMOX_ERROR)")
		}
		if !strings.Contains(eventsCode, "MotivoTimeout") || !strings.Contains(eventsCode, "TIMEOUT") {
			t.Errorf("ports/event_port.go no incluye constante de motivo de timeout (MotivoTimeout / TIMEOUT)")
		}

		// 4. nombresAccion debe incluir textos para shutdown, reboot y delete
		hasShutdown := strings.Contains(seguimientoCode, `"shutdown"`)
		hasReboot := strings.Contains(seguimientoCode, `"reboot"`)
		hasDelete := strings.Contains(seguimientoCode, `"delete"`)
		if !hasShutdown || !hasReboot || !hasDelete {
			t.Errorf("nombresAccion en seguimiento_tareas.go debe incluir shutdown, reboot y delete")
		}

		// 5. detalles de TASK_FINISHED incluye exitstatus, motivo y error según BAC-29
		hasDetallesExitStatus := strings.Contains(seguimientoCode, `"exitstatus"`)
		hasDetallesMotivo := strings.Contains(seguimientoCode, `"motivo"`)
		if !hasDetallesExitStatus || !hasDetallesMotivo {
			t.Errorf("detalles del TASK_FINISHED en seguimiento_tareas.go debe incluir exitstatus y motivo")
		}
	})
}
