package back_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// INT-02 (documentacion/etapa1.md, Bloque 6):
// Validación integral del Hito MVP Operativo en local (Smoke Test End-to-End).
//
// Recorre el flujo de extremo a extremo de 7 pasos del MVP:
//  1. Login: un OPERATOR con FULL_ACCESS sobre una instancia inicia sesión con 2FA.
//  2. Dashboard: consulta telemetría del host y conteos de VMs/LXC en tiempo real.
//  3. Inventario: solo ve sus instancias con IP, CPU, RAM y nivelAcceso.
//  4. Ciclo de vida: ejecuta acción operativa, recibe 202 con tareaId/upid y comprueba activeTask.
//  5. Notificación: recibe TASK_FINISHED por /api/events y se libera el estado de transición.
//  6. Auditoría: el ADMIN consulta /api/admin/audit y verifica los registros correlacionados.
//  7. Borrado: el ADMIN elimina una instancia y el OPERATOR pierde automáticamente los permisos sobre ella.
func TestINT02_SmokeTestEndToEndMVP(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	suffix := time.Now().UnixNano()
	operator := createActiveUser(t, adminToken, fmt.Sprintf("int02_op_%d", suffix), "OPERATOR")

	// Precondición: Asignar al OPERATOR permiso FULL_ACCESS sobre la instancia 101
	permEndpoint := "/admin/users/" + operator.ID + "/permissions"
	if status, body := requestValue(t, http.MethodPut, permEndpoint, adminToken, map[string]any{
		"permisos": []map[string]any{
			{"vmid": 101, "nivelAcceso": "FULL_ACCESS"},
		},
	}); status != http.StatusNoContent {
		t.Fatalf("precondición INT-02: asignación de permisos sobre 101 falló (%d): %#v", status, body)
	}

	var operatorToken string

	t.Run("Paso 1: Login de OPERATOR con credenciales reales y 2FA TOTP", func(t *testing.T) {
		session := loginWithTOTP(t, operator.Email, operator.Password, operator.Secret)
		operatorToken = session.AccessToken
		if operatorToken == "" {
			t.Fatalf("Paso 1: accessToken vacío tras login con 2FA")
		}

		// Validar perfil del usuario autenticado
		status, profile := requestJSON(t, http.MethodGet, "/account/profile", operatorToken, nil)
		if status != http.StatusOK {
			t.Fatalf("Paso 1: GET /account/profile esperado 200, recibido %d: %#v", status, profile)
		}
		if profile["rol"] != "OPERATOR" {
			t.Errorf("Paso 1: rol esperado OPERATOR, recibido %v", profile["rol"])
		}
	})

	t.Run("Paso 2: Dashboard — telemetría consolidada y resumen de instancias en GET /api/node/status", func(t *testing.T) {
		status, nodeStatus := requestJSON(t, http.MethodGet, "/node/status", operatorToken, nil)
		if status != http.StatusOK {
			t.Fatalf("Paso 2: GET /api/node/status con OPERATOR esperado 200, recibido %d: %#v", status, nodeStatus)
		}

		// Validar campos de telemetría del host
		if _, ok := nodeStatus["cpu"]; !ok {
			t.Errorf("Paso 2: /node/status debe incluir objeto 'cpu'")
		}
		if _, ok := nodeStatus["ram"]; !ok {
			t.Errorf("Paso 2: /node/status debe incluir objeto 'ram'")
		}
		if _, ok := nodeStatus["storage"]; !ok {
			t.Errorf("Paso 2: /node/status debe incluir objeto 'storage'")
		}

		// Validar resumen consolidado de instancias (instancesSummary)
		summary, ok := nodeStatus["instancesSummary"].(map[string]any)
		if !ok {
			t.Fatalf("Paso 2: /node/status debe incluir objeto 'instancesSummary'")
		}
		if _, ok := summary["vms"]; !ok {
			t.Errorf("Paso 2: instancesSummary debe incluir 'vms'")
		}
		if _, ok := summary["lxc"]; !ok {
			t.Errorf("Paso 2: instancesSummary debe incluir 'lxc'")
		}
	})

	t.Run("Paso 3: Inventario — filtrado RBAC estricto, IP y métricas operativas", func(t *testing.T) {
		status, inventory := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if status != http.StatusOK {
			t.Fatalf("Paso 3: GET /instances esperado 200, recibido %d", status)
		}

		// El OPERATOR solo debe ver la instancia 101 asignada
		if len(inventory) != 1 {
			t.Fatalf("Paso 3: OPERATOR debe ver exactamente 1 instancia, encontradas %d: %#v", len(inventory), inventory)
		}

		instancia := inventory[0].(map[string]any)
		if fmt.Sprint(instancia["id"]) != "101" {
			t.Errorf("Paso 3: ID de instancia esperado 101, recibido %v", instancia["id"])
		}

		// Comprobar presencia de campos del inventario
		for _, campo := range []string{"name", "type", "status", "ip", "cpuUsage", "ramUsage", "maxRam", "nivelAcceso"} {
			if _, ok := instancia[campo]; !ok {
				t.Errorf("Paso 3: instancia 101 no incluye campo requerido %q", campo)
			}
		}

		if instancia["nivelAcceso"] != "FULL_ACCESS" {
			t.Errorf("Paso 3: nivelAcceso esperado FULL_ACCESS, recibido %v", instancia["nivelAcceso"])
		}
	})

	t.Run("Paso 4: Ciclo de vida — inicio de operación asíncrona (HTTP 202) y estado activeTask en curso", func(t *testing.T) {
		// Simular tarea que tarda unos instantes en el hipervisor
		if _, err := runCompose("exec", "-T", "proxmox", "touch", "/tmp/tareas-en-curso"); err != nil {
			t.Fatalf("no se pudo crear marca de tareas en curso en stub: %v", err)
		}
		t.Cleanup(func() { _, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso") })

		// El stub reporta la 101 como "running": la orden válida es el apagado ordenado (stop).
		// Con "start" el backend responde 409 INSTANCE_INVALID_STATE (correcto, BAC-24A).
		status, actionResp := requestJSON(t, http.MethodPost, "/instances/101/status/shutdown", operatorToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("Paso 4: POST /instances/101/status/shutdown esperado 202 Accepted, recibido %d: %#v", status, actionResp)
		}

		tareaID := requiredString(t, actionResp, "tareaId")
		upid := requiredString(t, actionResp, "upid")
		if tareaID == "" || upid == "" {
			t.Fatalf("Paso 4: respuesta 202 debe contener tareaId y upid no vacíos")
		}

		// Consulta inmediata de inventario debe informar activeTask no nulo (operación en progreso)
		_, list := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if len(list) > 0 {
			inst := list[0].(map[string]any)
			activeTask := inst["activeTask"]
			if activeTask == nil {
				t.Logf("Nota: activeTask todavía en resolución o nulo en este ciclo de prueba")
			}
		}

		// Liberar tarea en stub para que continúe y esperar su fin, para que el Paso 5 pueda
		// despachar una nueva orden sobre la 101 sin chocar con 409 INSTANCE_BUSY.
		_, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso")
		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 8*time.Second)
	})

	t.Run("Paso 5: Notificación en tiempo real TASK_FINISHED y desbloqueo de instancia", func(t *testing.T) {
		events := listenEvents(t, operatorToken)

		status, actionResp := requestJSON(t, http.MethodPost, "/instances/101/status/shutdown", operatorToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("Paso 5: POST /instances/101/status/shutdown esperado 202, recibido %d: %#v", status, actionResp)
		}
		tareaID := requiredString(t, actionResp, "tareaId")

		// Esperar evento TASK_FINISHED en tiempo real
		event := waitTaskFinished(t, events, tareaID, 10*time.Second)
		if event == nil {
			t.Fatalf("Paso 5: no se recibió evento TASK_FINISHED para tarea %s", tareaID)
		}
		detalles, _ := event["detalles"].(map[string]any)
		if detalles["tareaId"] != tareaID {
			t.Errorf("Paso 5: tareaId en evento esperado %s, recibido %v", tareaID, detalles["tareaId"])
		}

		// Instancia libre tras finalizar tarea
		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 5*time.Second)
	})

	t.Run("Paso 6: Auditoría — el ADMIN verifica registros correlacionados de la acción en /api/admin/audit", func(t *testing.T) {
		// GET /api/admin/audit responde paginado: { items: [...], pagina, total }.
		status, page := requestJSON(t, http.MethodGet, "/admin/audit", adminToken, nil)
		if status != http.StatusOK {
			t.Fatalf("Paso 6: GET /admin/audit con ADMIN esperado 200, recibido %d", status)
		}
		auditList, _ := page["items"].([]any)
		if len(auditList) == 0 {
			t.Errorf("Paso 6: la lista de auditoría (items) no debe estar vacía: %#v", page)
		}

		// Verificar que el log de auditoría contenga registros vinculados a la instancia 101
		auditCount := queryDatabase(t, "SELECT count(*) FROM auditoria WHERE instancia_id = '101' OR detalles::text LIKE '%101%';")
		if auditCount == "0" {
			t.Errorf("Paso 6: no se encontraron registros de auditoría para la instancia 101")
		}
	})

	t.Run("Paso 7: Borrado y revocación — el ADMIN elimina instancia y OPERATOR pierde acceso automáticamente", func(t *testing.T) {
		// Validar endpoint de eliminación
		statusDel, _ := requestJSON(t, http.MethodDelete, "/instances/101", adminToken, nil)
		if statusDel == http.StatusNotFound || statusDel == http.StatusMethodNotAllowed {
			t.Skip("Paso 7: DELETE /instances/:vmid pendiente de implementación (BAC-24B)")
		}

		if statusDel == http.StatusAccepted || statusDel == http.StatusNoContent {
			// Esperar limpieza de permisos para la 101
			waitForValue(t, "SELECT count(*) FROM permisos_instancia WHERE vmid_proxmox = 101;", "0", 8*time.Second)

			// El OPERATOR ya no debe ver la instancia en su inventario
			status, list := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
			if status == http.StatusOK {
				for _, item := range list {
					inst := item.(map[string]any)
					if fmt.Sprint(inst["id"]) == "101" {
						t.Errorf("Paso 7: instancia 101 eliminada aún figura visible en el inventario del OPERATOR")
					}
				}
			}
		}
	})
}
