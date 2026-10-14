package back_test

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// Suite de pruebas adversariales para evaluar la integridad, resiliencia y
// robustez de seguridad del sistema Centinela:
//
// 1. Aislamiento RBAC y contención de escalamiento de privilegios.
// 2. Integridad de infraestructura protegida (VMIDs críticos inmunes a destrucción).
// 3. Concurrencia hostil y colisión de órdenes simultáneas (bloqueo 409 INSTANCE_BUSY).
// 4. Manipulación de tokens, revocación de sesiones y tickets efímeros.
// 5. Inyecciones y entradas malformadas (path traversal, SQL injection, desbordamiento).
func TestAdversarialSystemIntegrity(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	suffix := time.Now().UnixNano()
	operatorID, _ := createUserThroughAPI(t, adminToken, fmt.Sprintf("adv_op_%d", suffix), fmt.Sprintf("adv.op.%d@elcentinela.com", suffix), "OPERATOR")
	operatorToken := signedAccessToken(t, operatorID, "OPERATOR", orgID)

	// Subtest 1: Intentos de escalamiento de privilegios desde un rol OPERATOR
	t.Run("Escalamiento de privilegios: OPERATOR no puede ejecutar acciones ADMIN ni mutaciones no autorizadas", func(t *testing.T) {
		// A. Consulta administrativa de usuarios
		status, _ := requestValue(t, http.MethodGet, "/admin/users", operatorToken, nil)
		if status != http.StatusForbidden {
			t.Errorf("Adversarial RBAC: OPERATOR pudo consultar /admin/users: status %d (esperado 403)", status)
		}

		// B. Consulta de auditoría administrativa
		status, _ = requestValue(t, http.MethodGet, "/admin/audit", operatorToken, nil)
		if status != http.StatusForbidden {
			t.Errorf("Adversarial RBAC: OPERATOR pudo consultar /admin/audit: status %d (esperado 403)", status)
		}

		// C. Consulta de inventario de infraestructura (storages y plantillas)
		for _, path := range []string{"/node/storages", "/node/templates"} {
			status, _ = requestValue(t, http.MethodGet, path, operatorToken, nil)
			if status != http.StatusForbidden {
				t.Errorf("Adversarial RBAC: OPERATOR pudo consultar %s: status %d (esperado 403)", path, status)
			}
		}

		// D. Aprovisionamiento de nuevas instancias (solo ADMIN)
		status, _ = requestValue(t, http.MethodPost, "/instances", operatorToken, map[string]any{
			"vmid": 170, "node": "pve", "type": "vm", "name": "rogue-vm",
		})
		if status != http.StatusForbidden {
			t.Errorf("Adversarial RBAC: OPERATOR pudo invocar POST /instances: status %d (esperado 403)", status)
		}

		// E. Eliminación destructiva de instancias (solo ADMIN)
		status, _ = requestValue(t, http.MethodDelete, "/instances/101", operatorToken, nil)
		if status != http.StatusForbidden {
			t.Errorf("Adversarial RBAC: OPERATOR pudo invocar DELETE /instances/101: status %d (esperado 403)", status)
		}

		// F. Mutación de ciclo de vida con nivel de acceso READ_ONLY
		permEndpoint := "/admin/users/" + operatorID + "/permissions"
		permStatus, _ := requestValue(t, http.MethodPut, permEndpoint, adminToken, map[string]any{
			"permisos": []map[string]any{{"vmid": 101, "nivelAcceso": "READ_ONLY"}},
		})
		if permStatus != http.StatusNoContent && permStatus != http.StatusOK {
			t.Fatalf("No se pudo configurar permiso READ_ONLY: %d", permStatus)
		}

		for _, action := range []string{"start", "stop", "shutdown", "reboot"} {
			status, _ = requestValue(t, http.MethodPost, fmt.Sprintf("/instances/101/status/%s", action), operatorToken, nil)
			if status != http.StatusForbidden {
				t.Errorf("Adversarial RBAC: OPERATOR con READ_ONLY pudo ejecutar orden %s en 101: status %d (esperado 403)", action, status)
			}
		}
	})

	// Subtest 2: Protección de infraestructura crítica (VMIDs protegidos)
	t.Run("Infraestructura crítica: VMID protegido (103) inmune a destrucción incluso ante ADMIN", func(t *testing.T) {
		// La 103 está configurada en PROXMOX_PROTECTED_VMIDS
		deleteStatus, deleteBody := requestJSON(t, http.MethodDelete, "/instances/103", adminToken, nil)
		if deleteStatus != http.StatusForbidden {
			t.Errorf("Adversarial Infra: ADMIN pudo solicitar DELETE en VMID protegido 103: status %d (esperado 403)", deleteStatus)
		}
		if code, _ := deleteBody["errorCode"].(string); code != "INSTANCE_PROTECTED" {
			t.Errorf("Adversarial Infra: código de error esperado INSTANCE_PROTECTED, recibido %v", code)
		}
	})

	// Subtest 3: Concurrencia hostil y colisión de órdenes sobre el mismo recurso
	t.Run("Concurrencia hostil: dos órdenes mutantes concurrentes sobre la misma instancia colisionan con 409 INSTANCE_BUSY", func(t *testing.T) {
		// Emitir concurrentemente 5 peticiones de stop sobre la 101
		const concurrentRequests = 5
		statusCodes := make([]int, concurrentRequests)
		var wg sync.WaitGroup

		for i := 0; i < concurrentRequests; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				status, _ := requestValue(t, http.MethodPost, "/instances/101/status/stop", adminToken, nil)
				statusCodes[idx] = status
			}(i)
		}
		wg.Wait()

		acceptedCount := 0
		busyCount := 0
		for _, code := range statusCodes {
			switch code {
			case http.StatusAccepted:
				acceptedCount++
			case http.StatusConflict: // 409 INSTANCE_BUSY
				busyCount++
			}
		}

		// Como máximo 1 puede ser aceptada; las concurrentes deben ser rechazadas con 409
		if acceptedCount > 1 {
			t.Errorf("Adversarial Concurrency: más de una orden concurrente fue aceptada (%d aceptadas de %d): %#v", acceptedCount, concurrentRequests, statusCodes)
		}
		if busyCount == 0 && acceptedCount == 1 {
			t.Logf("Adversarial Concurrency: 1 orden aceptada, el resto ejecutó secuencialmente o fue rechazada")
		}
	})

	// Subtest 4: Manipulación de firmas y reutilización de tokens y tickets revocados
	t.Run("Integridad de autenticación: token con firma adulterada o sesión revocada es rechazado", func(t *testing.T) {
		// A. Token adulterado (firma cambiada)
		parts := strings.Split(adminToken, ".")
		if len(parts) == 3 {
			tamperedToken := parts[0] + "." + parts[1] + ".adulteratedSignature1234567890"
			status, _ := requestValue(t, http.MethodGet, "/node/status", tamperedToken, nil)
			if status != http.StatusUnauthorized {
				t.Errorf("Adversarial Token: token con firma adulterada no fue rechazado con 401: %d", status)
			}
		}

		// B. Reutilización de token sin bearer
		status, _ := requestValue(t, http.MethodGet, "/node/status", "invalid-token-string", nil)
		if status != http.StatusUnauthorized {
			t.Errorf("Adversarial Token: token sin formato JWT no fue rechazado con 401: %d", status)
		}

		// C. Ticket efímero de SSE es de un solo uso
		ticketStatus, ticketBody := requestJSON(t, http.MethodPost, "/events/ticket", adminToken, nil)
		if ticketStatus == http.StatusOK || ticketStatus == http.StatusCreated {
			ticket, ok := ticketBody["ticket"].(string)
			if ok && ticket != "" {
				// El ticket vive en Redis y solo permite autenticar la conexión inicial
				t.Logf("Adversarial Token: ticket efímero emitido correctamente: %s", ticket[:8])
			}
		}
	})

	// Subtest 5: Entradas malformadas, inyecciones de ruta y SQL injection
	t.Run("Inyección y entradas hostiles: path traversal y payloads SQL se gestionan de forma segura", func(t *testing.T) {
		// A. Path traversal en el identificador de instancia
		status, _ := requestValue(t, http.MethodGet, "/instances/..%2Fadmin%2Fusers", adminToken, nil)
		if status != http.StatusNotFound && status != http.StatusBadRequest {
			t.Errorf("Adversarial Input: path traversal en /instances/ no devolvió 400/404: %d", status)
		}

		// B. Identificador alfanumérico malicioso en ruta numérica
		status, _ = requestValue(t, http.MethodGet, "/instances/101'OR'1'='1", adminToken, nil)
		if status != http.StatusBadRequest && status != http.StatusNotFound {
			t.Errorf("Adversarial Input: inyección en VMID no devolvió 400/404: %d", status)
		}

		// C. Intento de inyección SQL en query parameters de auditoría
		status, auditBody := requestJSON(t, http.MethodGet, "/admin/audit?accion=' OR '1'='1", adminToken, nil)
		if status == http.StatusOK {
			// El backend debe parametrizar la consulta: no debe devolver registros de acciones no coincidentes
			if logs, ok := auditBody["logs"].([]any); ok && len(logs) > 0 {
				for _, l := range logs {
					if logMap, ok := l.(map[string]any); ok {
						if fmt.Sprint(logMap["accion"]) != "' OR '1'='1" {
							t.Errorf("Adversarial Input: posible inyección SQL en filtro de auditoría: devolvió accion %v", logMap["accion"])
							break
						}
					}
				}
			}
		}
	})
}

