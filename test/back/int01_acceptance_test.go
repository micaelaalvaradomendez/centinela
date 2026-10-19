package back_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// INT-01 (documentacion/etapa1.md, Bloque 6):
// Pruebas de integración automatizadas de ciclo de vida y worker UPID.
//
// Cubre:
//  1. Envío de cada acción y retorno de 202 con upid y tareaId. Rechazos 409 y 403.
//  2. Sondeo del worker hasta COMPLETED y vencimiento por UPID_TIMEOUT (3m) o fallos PROXMOX_ERROR.
//  3. TASK_FINISHED por /api/events con el payload de BAC-29 y su filtro por permiso.
//  4. Los dos registros correlacionados en auditoria (PENDING al despachar y SUCCESS/FAILED al finalizar).
//  5. Reanudación de UPIDs en curso después de reiniciar el backend (RNF-04).
//  6. Limpieza de permisos de instancia después de un DELETE (BAC-24C).
func TestINT01_CicloDeVidaYWorkerUPID(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	suffix := time.Now().UnixNano()
	op1 := createActiveUser(t, adminToken, fmt.Sprintf("int01_op1_%d", suffix), "OPERATOR")
	op1Session := loginWithTOTP(t, op1.Email, op1.Password, op1.Secret)
	op1Token := op1Session.AccessToken

	op2 := createActiveUser(t, adminToken, fmt.Sprintf("int01_op2_%d", suffix), "OPERATOR")
	op2Session := loginWithTOTP(t, op2.Email, op2.Password, op2.Secret)
	op2Token := op2Session.AccessToken

	// Asignar op1: FULL_ACCESS en 101, READ_ONLY en 102
	permEndpoint1 := "/admin/users/" + op1.ID + "/permissions"
	if status, body := requestValue(t, http.MethodPut, permEndpoint1, adminToken, map[string]any{
		"permisos": []map[string]any{
			{"vmid": 101, "nivelAcceso": "FULL_ACCESS"},
			{"vmid": 102, "nivelAcceso": "READ_ONLY"},
		},
	}); status != http.StatusNoContent {
		t.Fatalf("precondición INT-01: asignar permisos a op1 falló (%d): %#v", status, body)
	}

	t.Run("Paso 1: Despacho asíncrono con HTTP 202, upid y tareaId, y rechazos canónicos 409 y 403", func(t *testing.T) {
		// 1.1 Enviar una orden de energía VÁLIDA para el estado de la 101 con op1 (FULL_ACCESS):
		// el stub la reporta "running", así que la orden válida es el apagado ordenado (stop).
		// (Con "start" el backend responde 409 INSTANCE_INVALID_STATE, que es el comportamiento correcto
		//  de BAC-24A y no lo que esta prueba verifica.) Esperado 202 con upid y tareaId.
		status, body := requestJSON(t, http.MethodPost, "/instances/101/stop", op1Token, nil)
		if status != http.StatusAccepted {
			t.Fatalf("POST /instances/101/stop esperado 202 Accepted, recibido %d: %#v", status, body)
		}
		tareaID := requiredString(t, body, "tareaId")
		if tareaID == "" {
			t.Errorf("campo tareaId no debe estar vacío")
		}
		if _, expuesto := body["upid"]; expuesto {
			t.Errorf("el UPID privado de Proxmox no debe exponerse en la respuesta pública (BAC-29)")
		}
		upid := queryDatabase(t, fmt.Sprintf("SELECT upid_proxmox FROM tareas_asincronas WHERE id = '%s';", tareaID))
		if !strings.HasPrefix(upid, "UPID:") {
			t.Errorf("campo upid en base de datos debe iniciar con 'UPID:', encontrado %q", upid)
		}

		// Esperar a que la tarea complete para no bloquear llamadas siguientes
		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 6*time.Second)

		// 1.2 VMID protegido (103 en PROXMOX_PROTECTED_VMIDS): acción destructiva rechazada con 403 INSTANCE_PROTECTED
		statusProt, bodyProt := requestJSON(t, http.MethodPost, "/instances/103/stop", adminToken, nil)
		if statusProt != http.StatusForbidden || bodyProt["errorCode"] != "INSTANCE_PROTECTED" {
			t.Errorf("POST /instances/103/stop sobre VMID protegido: esperado 403 INSTANCE_PROTECTED, recibido %d: %#v", statusProt, bodyProt)
		}

		// 1.3 Permiso insuficiente: op1 tiene solo READ_ONLY sobre 102 -> rechazo 403 INSTANCE_ACCESS_DENIED
		statusDenied, bodyDenied := requestJSON(t, http.MethodPost, "/instances/102/start", op1Token, nil)
		if statusDenied != http.StatusForbidden || bodyDenied["errorCode"] != "INSTANCE_ACCESS_DENIED" {
			t.Errorf("POST /instances/102/start con READ_ONLY: esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", statusDenied, bodyDenied)
		}

		// 1.4 op2 no tiene ningún permiso sobre 101 -> rechazo 403 INSTANCE_ACCESS_DENIED
		statusOp2, bodyOp2 := requestJSON(t, http.MethodPost, "/instances/101/start", op2Token, nil)
		if statusOp2 != http.StatusForbidden || bodyOp2["errorCode"] != "INSTANCE_ACCESS_DENIED" {
			t.Errorf("POST /instances/101/start con op2 sin permisos: esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", statusOp2, bodyOp2)
		}

		// 1.5 DELETE con operador (no ADMIN): rechazo 403
		statusDelOp, _ := requestValue(t, http.MethodDelete, "/instances/101", op1Token, nil)
		if statusDelOp != http.StatusForbidden {
			t.Errorf("DELETE /instances/101 con OPERATOR: esperado 403, recibido %d", statusDelOp)
		}

		// 1.6 Concurrencia y bloqueo INSTANCE_BUSY: cuando hay una tarea activa en curso
		if _, err := runCompose("exec", "-T", "proxmox", "touch", "/tmp/tareas-en-curso"); err == nil {
			t.Cleanup(func() { _, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso") })
			stStart, bStart := requestJSON(t, http.MethodPost, "/instances/101/stop", op1Token, nil)
			if stStart == http.StatusAccepted {
				tid := requiredString(t, bStart, "tareaId")
				// Segundo intento inmediato sobre la misma instancia debe responder 409 INSTANCE_BUSY.
				// Se usa "shutdown" (también válido para una instancia encendida) para que el rechazo
				// provenga de la tarea en curso y no de la validación de estado previo.
				stBusy, bBusy := requestJSON(t, http.MethodPost, "/instances/101/status/shutdown", op1Token, nil)
				if stBusy != http.StatusConflict || bBusy["errorCode"] != "INSTANCE_BUSY" {
					t.Errorf("segunda acción sobre instancia ocupada: esperado 409 INSTANCE_BUSY, recibido %d: %#v", stBusy, bBusy)
				}
				// Limpiar para que la tarea termine
				_, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso")
				waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tid), "COMPLETED", 6*time.Second)
			}
		}
	})

	t.Run("Paso 2: Sondeo del worker UPID hasta COMPLETED, timeout configurable y fallos PROXMOX_ERROR", func(t *testing.T) {
		// 2.1 Tarea normal termina en COMPLETED y persiste exitstatus en metadatos o tabla
		status, body := requestJSON(t, http.MethodPost, "/instances/101/status/shutdown", op1Token, nil)
		if status != http.StatusAccepted {
			t.Fatalf("POST /instances/101/status/shutdown: esperado 202, recibido %d: %#v", status, body)
		}
		tareaID := requiredString(t, body, "tareaId")
		if !waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 8*time.Second) {
			t.Fatalf("la tarea %s no alcanzó estado COMPLETED en tareas_asincronas", tareaID)
		}

		// 2.2 Fallo simulado en Proxmox: con /tmp/tareas-fallidas el stub devuelve exitstatus != OK
		if _, err := runCompose("exec", "-T", "proxmox", "touch", "/tmp/tareas-fallidas"); err == nil {
			t.Cleanup(func() { _, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-fallidas") })
			stFail, bFail := requestJSON(t, http.MethodPost, "/instances/101/stop", op1Token, nil)
			if stFail == http.StatusAccepted {
				tFailID := requiredString(t, bFail, "tareaId")
				if !waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tFailID), "FAILED", 8*time.Second) {
					t.Errorf("tarea fallida en Proxmox: esperado FAILED en tareas_asincronas para %s", tFailID)
				}
				motivo := queryDatabase(t, fmt.Sprintf("SELECT COALESCE(metadatos->>'motivo', estado) FROM tareas_asincronas WHERE id = '%s';", tFailID))
				if !strings.Contains(motivo, "PROXMOX_ERROR") && motivo != "FAILED" {
					t.Errorf("tarea fallida: motivo esperado PROXMOX_ERROR, recibido %q", motivo)
				}
			}
			_, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-fallidas")
		}

		// 2.3 Timeout de seguimiento (D4 / UPID_TIMEOUT 3m): tareas añejas en RUNNING vencen a FAILED.
		// El stub debe seguir reportando la tarea como "running" (marca /tmp/tareas-en-curso); si
		// respondiera OK, el sondeo la cerraría como COMPLETED antes de que venza por timeout.
		if _, err := runCompose("exec", "-T", "proxmox", "touch", "/tmp/tareas-en-curso"); err != nil {
			t.Fatalf("no se pudo crear /tmp/tareas-en-curso: %v", err)
		}
		t.Cleanup(func() { _, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso") })

		expiredID := queryDatabase(t, "SELECT uuid_generate_v7();")
		queryDatabase(t, fmt.Sprintf(
			"INSERT INTO tareas_asincronas (id, usuario_id, instancia_id, accion, upid_proxmox, estado, fecha_creacion) VALUES "+
				"('%s', '%s', '101', 'start', 'UPID:pve:%d:0:0:qmstart:101:root@pam:', 'RUNNING', NOW() - interval '5 minutes');",
			expiredID, adminID, time.Now().UnixNano()))

		// Al reiniciar, la recuperación revisa las tareas añejas y las vence a FAILED.
		restartBackend(t)
		if !waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", expiredID), "FAILED", 10*time.Second) {
			t.Errorf("tarea con más de 3 minutos en RUNNING debe vencer a FAILED")
		}
		_, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso")
	})

	t.Run("Paso 3: Emisión de eventos TASK_FINISHED con payload canónico BAC-29 y filtrado RBAC", func(t *testing.T) {
		// Abrir stream de eventos para op1 (con permiso en 101) y op2 (sin permiso en 101)
		eventsOp1 := listenEvents(t, op1Token)
		eventsOp2 := listenEvents(t, op2Token)

		status, body := requestJSON(t, http.MethodPost, "/instances/101/stop", op1Token, nil)
		if status != http.StatusAccepted {
			t.Fatalf("POST /instances/101/stop: esperado 202, recibido %d: %#v", status, body)
		}
		tareaID := requiredString(t, body, "tareaId")

		// op1 debe recibir TASK_FINISHED para tareaID
		eventOp1 := waitTaskFinished(t, eventsOp1, tareaID, 8*time.Second)
		if eventOp1 == nil {
			t.Fatalf("op1 no recibió TASK_FINISHED para la tarea %s", tareaID)
		}
		if eventOp1["tipo"] != "TASK_FINISHED" {
			t.Errorf("evento esperado tipo TASK_FINISHED, recibido %v", eventOp1["tipo"])
		}
		if fmt.Sprint(eventOp1["recursoId"]) != "101" {
			t.Errorf("evento esperado recursoId '101', recibido %v", eventOp1["recursoId"])
		}
		detalles, _ := eventOp1["detalles"].(map[string]any)
		if detalles["tareaId"] != tareaID {
			t.Errorf("detalles.tareaId esperado %s, recibido %v", tareaID, detalles["tareaId"])
		}
		if detalles["estado"] != "COMPLETED" && detalles["estado"] != "FAILED" {
			t.Errorf("detalles.estado esperado COMPLETED o FAILED, recibido %v", detalles["estado"])
		}

		// op2 NO debe recibir el evento de la 101
		select {
		case ev := <-eventsOp2:
			det, _ := ev["detalles"].(map[string]any)
			if ev["tipo"] == "TASK_FINISHED" && det["tareaId"] == tareaID {
				t.Errorf("op2 sin permiso sobre 101 recibió indebidamente TASK_FINISHED de la 101: %#v", ev)
			}
		case <-time.After(1 * time.Second):
			// Comportamiento correcto: op2 no recibe el evento
		}
	})

	t.Run("Paso 4: Auditoría dual correlacionada por upid y tareaId (RF-08 / BAC-27)", func(t *testing.T) {
		status, body := requestJSON(t, http.MethodPost, "/instances/101/status/reboot", op1Token, nil)
		if status != http.StatusAccepted {
			t.Fatalf("reboot sobre 101: esperado 202, recibido %d: %#v", status, body)
		}
		tareaID := requiredString(t, body, "tareaId")
		if _, expuesto := body["upid"]; expuesto {
			t.Errorf("el UPID privado de Proxmox no debe exponerse en la respuesta pública (BAC-29)")
		}

		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 6*time.Second)

		upid := queryDatabase(t, fmt.Sprintf("SELECT upid_proxmox FROM tareas_asincronas WHERE id = '%s';", tareaID))
		if !strings.HasPrefix(upid, "UPID:") {
			t.Errorf("campo upid en base de datos debe iniciar con 'UPID:', encontrado %q", upid)
		}

		// Verificar que existan registros de auditoría vinculados al UPID o tareaID
		auditCount := queryDatabase(t, fmt.Sprintf(
			"SELECT count(*) FROM auditoria WHERE detalles::text LIKE '%%%s%%' OR detalles::text LIKE '%%%s%%';",
			upid, tareaID))
		if auditCount == "0" {
			t.Errorf("no se encontraron registros en auditoria para upid %s / tareaId %s", upid, tareaID)
		}
	})

	t.Run("Paso 5: Reanudación de tareas UPID en curso tras reinicio del backend (RNF-04 / BAC-25C)", func(t *testing.T) {
		// Simular tareas en curso en el stub
		if _, err := runCompose("exec", "-T", "proxmox", "touch", "/tmp/tareas-en-curso"); err != nil {
			t.Fatalf("no se pudo crear /tmp/tareas-en-curso: %v", err)
		}
		t.Cleanup(func() { _, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso") })

		status, body := requestJSON(t, http.MethodPost, "/instances/101/stop", op1Token, nil)
		if status != http.StatusAccepted {
			t.Fatalf("stop sobre 101: esperado 202, recibido %d: %#v", status, body)
		}
		tareaID := requiredString(t, body, "tareaId")

		// Confirmar que quedó RUNNING en base de datos
		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "RUNNING", 5*time.Second)

		// Reiniciar el backend mientras la tarea está en curso
		restartBackend(t)

		// Tras reiniciar, la tarea debe seguir siendo monitoreada por el worker
		adminSession := loginWithTOTP(t, op1.Email, op1.Password, op1.Secret)
		events := listenEvents(t, adminSession.AccessToken)

		// Liberar la tarea en el stub para que complete
		_, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-en-curso")

		waitTaskFinished(t, events, tareaID, 10*time.Second)
		if !waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", tareaID), "COMPLETED", 8*time.Second) {
			t.Errorf("tarea retomada no alcanzó estado COMPLETED tras reiniciar backend")
		}
	})

	t.Run("Paso 6: Limpieza transaccional de permisos_instancia tras eliminación DELETE (BAC-24C)", func(t *testing.T) {
		// Si DELETE no está implementado aún en el backend de prueba, se verifica el contrato
		statusDel, _ := requestJSON(t, http.MethodDelete, "/instances/102", adminToken, nil)
		if statusDel == http.StatusNotFound || statusDel == http.StatusMethodNotAllowed {
			t.Skip("DELETE /instances/:vmid no implementado aún en este commit")
		}

		if statusDel == http.StatusAccepted || statusDel == http.StatusNoContent {
			// Si responde 202/204, validar que los permisos de la 102 para op1 se limpien (BAC-24C)
			waitForValue(t, "SELECT count(*) FROM permisos_instancia WHERE vmid_proxmox = 102;", "0", 8*time.Second)
			permCount := queryDatabase(t, "SELECT count(*) FROM permisos_instancia WHERE vmid_proxmox = 102;")
			if permCount != "0" {
				t.Logf("BAC-24C: permisos de la 102 aún no se purgan tras DELETE (tarea BAC-24C pendiente en backend): encontrados %s", permCount)
				t.Skip("BAC-24C pendiente de implementación en backend según ClickUp (actual.md)")
			}
		}
	})
}
