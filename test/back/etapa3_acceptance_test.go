package back_test

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Pruebas de aceptación de INF-13 (cierre Etapa 2) y Etapa 3 de documentacion/actual.md:
//   - INF-13: Verificación y configuración de almacenamiento compatible con Snapshots en Proxmox.
//   - BAC-42: Adaptador Proxmox de Snapshots y endpoints de listado y creación asíncrona.
//   - BAC-43: Endpoints de Rollback y Eliminación de Snapshots con control de concurrencia.
//   - BAC-44: Interceptor global de auditoría y endpoint de consulta avanzada con filtros compuestos.
//   - BAC-45: Endpoint de exportación de auditoría en streaming (CSV y JSON) con sanitización anti-inyección.
//   - INF-11: Configuración de Nginx para producción con TLS 1.3, HSTS, CSP y soporte WSS/SSE.
//   - BAC-46: Worker evaluador de umbrales con histéresis y emisión de alertas SSE.

// TestEtapa2INF13_StorageCompatibleConSnapshots verifica INF-13 (documentacion/actual.md):
// Al menos un pool de almacenamiento activo reportado por GET /api/node/storages debe soportar
// imágenes de disco/contenedores (images / rootdir), ser de un tipo compatible con snapshots en
// caliente (lvmthin, zfspool o dir) y contar con al menos 2 GiB libres.
func TestEtapa2INF13_StorageCompatibleConSnapshots(t *testing.T) {
	t.Run("api-proxmox.md documenta compatibilidad de storages con snapshots (lvmthin, zfspool, qcow2 vs raw) sin exponer PROXMOX_TOKEN_SECRET en texto plano", func(t *testing.T) {
		clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))
		if !regexp.MustCompile(`(?i)ListarStorages|ObtenerStorages|/storage`).MatchString(clientCode) {
			t.Errorf("INF-13: client.go no implementa la consulta de storages del nodo Proxmox")
		}

		docPath := filepath.Join(envOrDefault("CENTINELA_ROOT", "../.."), "documentacion", "api-proxmox.md")
		docContent := string(readFile(t, docPath))
		for _, requiredTerm := range []string{"lvmthin", "qcow2", "raw"} {
			if !regexp.MustCompile(`(?i)\b` + requiredTerm + `\b`).MatchString(docContent) {
				t.Errorf("INF-13: documentacion/api-proxmox.md debe documentar la compatibilidad de %q con snapshots de Proxmox VE", requiredTerm)
			}
		}
		if !regexp.MustCompile(`(?i)\b(zfspool|zfs)\b`).MatchString(docContent) {
			t.Errorf("INF-13: documentacion/api-proxmox.md debe documentar zfspool/zfs como storage compatible con snapshots")
		}
		if regexp.MustCompile(`PROXMOX_TOKEN_SECRET=[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`).MatchString(docContent) {
			t.Errorf("INF-13: documentacion/api-proxmox.md expone un UUID real en PROXMOX_TOKEN_SECRET en lugar de un placeholder seguro")
		}
	})

	t.Run("GET /api/node/storages devuelve al menos un storage compatible con snapshots con >= 2 GiB libres", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)

		status, items := requestJSONArray(t, http.MethodGet, "/node/storages", adminToken)
		if status != http.StatusOK || len(items) == 0 {
			t.Fatalf("INF-13: GET /api/node/storages esperado 200 con storages, recibido %d: %#v", status, items)
		}

		const minFreeBytes = 2 * 1024 * 1024 * 1024 // 2 GiB
		compatibleTypes := map[string]bool{"lvmthin": true, "zfspool": true, "dir": true, "rbd": true, "cephfs": true}
		var foundCompatible bool

		for _, raw := range items {
			st, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			stType := strings.ToLower(fmt.Sprint(st["type"]))
			freeBytes, _ := st["freeBytes"].(float64)
			contentRaw, _ := st["contentTypes"].([]any)
			hasDiskContent := false
			for _, c := range contentRaw {
				ct := strings.ToLower(fmt.Sprint(c))
				if ct == "images" || ct == "rootdir" {
					hasDiskContent = true
					break
				}
			}
			if compatibleTypes[stType] && hasDiskContent && freeBytes >= minFreeBytes {
				foundCompatible = true
				break
			}
		}

		if !foundCompatible {
			t.Errorf("INF-13: ningún storage devuelto por GET /api/node/storages combina tipo compatible con snapshots (lvmthin/zfspool/dir), contentTypes (images/rootdir) y >= 2 GiB libres: %#v", items)
		}
	})
}

// TestEtapa3BAC42_AdaptadorYEndpointsSnapshotsCreacionYListado verifica BAC-42 (documentacion/actual.md):
//   1. Métodos en puerto y cliente Proxmox: ListarSnapshots, CrearSnapshot, RevertirSnapshot, EliminarSnapshot
//      filtrando el nodo ficticio "current" que devuelve Proxmox VE.
//   2. GET /api/instances/:instanceRef/snapshots con validación RBAC (ADMIN y OPERATOR con acceso).
//   3. POST /api/instances/:instanceRef/snapshots (ADMIN y OPERATOR con FULL_ACCESS; 403 para READ_ONLY),
//      validación de snapname (/^[a-zA-Z0-9_\-]+$/, máx. 40 caracteres) y description (<= 250),
//      respuesta 202 Accepted con { tareaId, snapname } (sin upid según FIX-110) y notificación SSE TASK_FINISHED.
func TestEtapa3BAC42_AdaptadorYEndpointsSnapshotsCreacionYListado(t *testing.T) {
	t.Run("puerto y cliente Proxmox definen ListarSnapshots, CrearSnapshot, RevertirSnapshot y EliminarSnapshot", func(t *testing.T) {
		portsCode := string(readFile(t, sourcePath("backend", "internal/core/ports/proxmox_port.go")))
		clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))
		clientTestCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client_test.go")))

		for _, method := range []string{"ListarSnapshots", "CrearSnapshot", "RevertirSnapshot", "EliminarSnapshot"} {
			if !strings.Contains(portsCode, method) && !strings.Contains(clientCode, method) {
				t.Errorf("BAC-42: falta el método %s en internal/core/ports/proxmox_port.go o internal/adapters/secondary/proxmox/client.go", method)
			}
		}
		if !regexp.MustCompile(`(?i)"current"`).MatchString(clientCode) {
			t.Errorf("BAC-42: ListarSnapshots en client.go debe filtrar el pseudo-snapshot \"current\" de Proxmox VE")
		}
		if !regexp.MustCompile(`(?i)Snapshot`).MatchString(clientTestCode) {
			t.Errorf("BAC-42: client_test.go no contiene pruebas unitarias para los métodos de snapshots")
		}
	})

	t.Run("GET y POST /api/instances/:instanceRef/snapshots aplican RBAC, validan snapname y despachan tarea asíncrona 202", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)
		ref101 := instanceRefByVMID(t, adminToken, 101)
		ref102 := instanceRefByVMID(t, adminToken, 102)

		opReadOnly := createActiveUser(t, adminToken, "bac42ro", "OPERATOR")
		opFull := createActiveUser(t, adminToken, "bac42full", "OPERATOR")

		// opReadOnly tiene READ_ONLY sobre 101 y ningún permiso sobre 102
		if status, body := requestValue(t, http.MethodPut, "/admin/users/"+opReadOnly.ID+"/permissions", adminToken, map[string]any{
			"permisos": []map[string]any{{"vmid": ref101, "nivelAcceso": "READ_ONLY"}},
		}); status != http.StatusNoContent {
			t.Fatalf("precondición BAC-42: asignar READ_ONLY falló (%d): %#v", status, body)
		}
		// opFull tiene FULL_ACCESS sobre 101
		if status, body := requestValue(t, http.MethodPut, "/admin/users/"+opFull.ID+"/permissions", adminToken, map[string]any{
			"permisos": []map[string]any{{"vmid": ref101, "nivelAcceso": "FULL_ACCESS"}},
		}); status != http.StatusNoContent {
			t.Fatalf("precondición BAC-42: asignar FULL_ACCESS falló (%d): %#v", status, body)
		}

		roToken := loginWithTOTP(t, opReadOnly.Email, opReadOnly.Password, opReadOnly.Secret).AccessToken
		fullToken := loginWithTOTP(t, opFull.Email, opFull.Password, opFull.Secret).AccessToken

		// 1. OPERATOR sin permiso sobre 102 recibe 403 en GET /snapshots
		if status, _ := requestValue(t, http.MethodGet, "/instances/"+ref102+"/snapshots", roToken, nil); status != http.StatusForbidden {
			t.Errorf("BAC-42: GET /instances/%s/snapshots sin permiso esperado 403, recibido %d", ref102, status)
		}

		// 2. OPERATOR con READ_ONLY sobre 101 puede listar snapshots (200) pero no crear (403)
		statusList, listVal := requestValue(t, http.MethodGet, "/instances/"+ref101+"/snapshots", roToken, nil)
		if statusList != http.StatusOK {
			t.Errorf("BAC-42: GET /instances/%s/snapshots con READ_ONLY esperado 200, recibido %d: %#v", ref101, statusList, listVal)
		} else {
			encoded, _ := json.Marshal(listVal)
			if regexp.MustCompile(`"name":\s*"current"`).Match(encoded) {
				t.Errorf("BAC-42: GET /instances/%s/snapshots no debe incluir el nodo ficticio 'current': %s", ref101, encoded)
			}
		}

		rawROPost := requestRaw(t, http.MethodPost, "/instances/"+ref101+"/snapshots", roToken, map[string]any{
			"snapname":    "snap_prohibido",
			"description": "No debería permitirse con READ_ONLY",
		})
		if rawROPost.Status != http.StatusForbidden {
			t.Errorf("BAC-42: POST /instances/%s/snapshots con READ_ONLY esperado 403, recibido %d: %s", ref101, rawROPost.Status, rawROPost.RawBody)
		}

		// 3. Validaciones de formato de snapname (/^[a-zA-Z0-9_\-]+$/, <= 40 chars) y description (<= 250 chars) -> 400 INVALID_REQUEST
		invalidPayloads := map[string]map[string]any{
			"snapname vacío":                  {"snapname": "", "description": "ok"},
			"snapname con caracteres inválidos": {"snapname": "snap invalido!", "description": "ok"},
			"snapname con barras":             {"snapname": "snap/../../etc", "description": "ok"},
			"snapname > 40 caracteres":        {"snapname": strings.Repeat("a", 41), "description": "ok"},
			"description > 250 chars":         {"snapname": "snap_desc_larga", "description": strings.Repeat("x", 251)},
		}
		for caseName, payload := range invalidPayloads {
			raw := requestRaw(t, http.MethodPost, "/instances/"+ref101+"/snapshots", fullToken, payload)
			if raw.Status != http.StatusBadRequest {
				t.Errorf("BAC-42: %s esperado 400 INVALID_REQUEST, recibido %d: %s", caseName, raw.Status, raw.RawBody)
			}
		}

		// 4. Creación válida con FULL_ACCESS devuelve 202 Accepted con { tareaId, snapname }, sin upid (FIX-110), y emite TASK_FINISHED
		events := listenEvents(t, fullToken)
		snapName := fmt.Sprintf("snap_bac42_%d", time.Now().Unix()%100000)
		statusCreate, createBody := requestJSON(t, http.MethodPost, "/instances/"+ref101+"/snapshots", fullToken, map[string]any{
			"snapname":    snapName,
			"description": "Snapshot creado en prueba de aceptación BAC-42",
			"vmstate":     false,
		})
		if statusCreate != http.StatusAccepted {
			t.Fatalf("BAC-42: POST /instances/%s/snapshots esperado 202 Accepted, recibido %d: %#v", ref101, statusCreate, createBody)
		}
		tareaID := requiredString(t, createBody, "tareaId")
		if gotSnap := fmt.Sprint(createBody["snapname"]); gotSnap != snapName {
			t.Errorf("BAC-42: respuesta 202 debe incluir snapname=%q, recibido %#v", snapName, createBody["snapname"])
		}
		if _, hasUPID := createBody["upid"]; hasUPID {
			t.Errorf("BAC-42 / FIX-110: respuesta 202 no debe exponer 'upid' de Proxmox: %#v", createBody)
		}
		if event := waitTaskFinished(t, events, tareaID, 15*time.Second); event == nil {
			t.Errorf("BAC-42: no llegó TASK_FINISHED por SSE tras crear el snapshot %s (tarea %s)", snapName, tareaID)
		}
	})
}

// TestEtapa3BAC43_RollbackYEliminacionSnapshotsConcurrencia verifica BAC-43 (documentacion/actual.md):
//   1. POST /api/instances/:instanceRef/snapshots/:snapname/rollback -> 202 Accepted ({ tareaId }).
//   2. DELETE /api/instances/:instanceRef/snapshots/:snapname -> 202 Accepted ({ tareaId }).
//   3. Control de concurrencia: rechaza con 409 INSTANCE_BUSY si hay otra tarea asíncrona RUNNING.
//   4. RBAC: requiere ADMIN u OPERATOR con FULL_ACCESS (403 para READ_ONLY).
func TestEtapa3BAC43_RollbackYEliminacionSnapshotsConcurrencia(t *testing.T) {
	t.Run("rutas de rollback y delete de snapshots declaradas en router y handlers", func(t *testing.T) {
		mainCode := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))
		handlersCode := secReadGoDir(t, "internal/adapters/primary/http")
		if !strings.Contains(mainCode, "/snapshots/:snapname/rollback") && !strings.Contains(handlersCode, "/snapshots/:snapname/rollback") {
			t.Errorf("BAC-43: no está declarada la ruta POST /instances/:instanceRef/snapshots/:snapname/rollback")
		}
		if !strings.Contains(mainCode, "/snapshots/:snapname") && !strings.Contains(handlersCode, "/snapshots/:snapname") {
			t.Errorf("BAC-43: no está declarada la ruta DELETE /instances/:instanceRef/snapshots/:snapname")
		}
	})

	t.Run("rollback y eliminación validan RBAC, bloquean con 409 INSTANCE_BUSY y despachan con 202 Accepted", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)
		ref101 := instanceRefByVMID(t, adminToken, 101)

		opReadOnly := createActiveUser(t, adminToken, "bac43ro", "OPERATOR")
		if status, body := requestValue(t, http.MethodPut, "/admin/users/"+opReadOnly.ID+"/permissions", adminToken, map[string]any{
			"permisos": []map[string]any{{"vmid": ref101, "nivelAcceso": "READ_ONLY"}},
		}); status != http.StatusNoContent {
			t.Fatalf("precondición BAC-43: asignar READ_ONLY falló (%d): %#v", status, body)
		}
		roToken := loginWithTOTP(t, opReadOnly.Email, opReadOnly.Password, opReadOnly.Secret).AccessToken

		// Crear un snapshot previo como ADMIN para probar rollback y delete
		snapName := fmt.Sprintf("snap_bac43_%d", time.Now().Unix()%100000)
		statusCreate, createBody := requestJSON(t, http.MethodPost, "/instances/"+ref101+"/snapshots", adminToken, map[string]any{
			"snapname":    snapName,
			"description": "Snapshot para rollback y delete",
		})
		if statusCreate != http.StatusAccepted {
			t.Fatalf("BAC-43: precondición crear snapshot esperado 202, recibido %d: %#v", statusCreate, createBody)
		}
		createTaskID := requiredString(t, createBody, "tareaId")
		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", createTaskID), "COMPLETED", 15*time.Second)

		// 1. READ_ONLY recibe 403 en rollback y delete
		if status, _ := requestValue(t, http.MethodPost, fmt.Sprintf("/instances/%s/snapshots/%s/rollback", ref101, snapName), roToken, nil); status != http.StatusForbidden {
			t.Errorf("BAC-43: POST rollback con READ_ONLY esperado 403, recibido %d", status)
		}
		if status, _ := requestValue(t, http.MethodDelete, fmt.Sprintf("/instances/%s/snapshots/%s", ref101, snapName), roToken, nil); status != http.StatusForbidden {
			t.Errorf("BAC-43: DELETE snapshot con READ_ONLY esperado 403, recibido %d", status)
		}

		// 2. Control de concurrencia: mientras existe una tarea EN_PROGRESO (ROLLING_BACK / SNAPSHOTTING) sobre 101,
		// cualquier intento concurrente de rollback, delete o acción de energía debe responder 409 INSTANCE_BUSY.
		busyTaskID := fmt.Sprintf("00000000-0000-4000-8000-%012d", time.Now().UnixNano()%1_000_000_000_000)
		queryDatabase(t, fmt.Sprintf(
			"INSERT INTO tareas_asincronas (id, vmid, nodo, upid, accion, estado, usuario_id) VALUES ('%s', 101, 'pve-node1', 'UPID:pve-node1:0000BAC43', 'SNAPSHOT_ROLLBACK', 'RUNNING', '%s');",
			busyTaskID, opReadOnly.ID,
		))
		defer queryDatabase(t, fmt.Sprintf("DELETE FROM tareas_asincronas WHERE id = '%s';", busyTaskID))

		rawBusyRB := requestRaw(t, http.MethodPost, fmt.Sprintf("/instances/%s/snapshots/%s/rollback", ref101, snapName), adminToken, nil)
		if rawBusyRB.Status != http.StatusConflict || !strings.Contains(rawBusyRB.RawBody, "INSTANCE_BUSY") {
			t.Errorf("BAC-43: rollback con tarea en curso esperado 409 INSTANCE_BUSY, recibido %d: %s", rawBusyRB.Status, rawBusyRB.RawBody)
		}
		rawBusyDel := requestRaw(t, http.MethodDelete, fmt.Sprintf("/instances/%s/snapshots/%s", ref101, snapName), adminToken, nil)
		if rawBusyDel.Status != http.StatusConflict || !strings.Contains(rawBusyDel.RawBody, "INSTANCE_BUSY") {
			t.Errorf("BAC-43: delete snapshot con tarea en curso esperado 409 INSTANCE_BUSY, recibido %d: %s", rawBusyDel.Status, rawBusyDel.RawBody)
		}
		rawBusyPower := requestRaw(t, http.MethodPost, fmt.Sprintf("/instances/%s/status/reboot", ref101), adminToken, nil)
		if rawBusyPower.Status != http.StatusConflict || !strings.Contains(rawBusyPower.RawBody, "INSTANCE_BUSY") {
			t.Errorf("BAC-43: acción de energía durante rollback esperado 409 INSTANCE_BUSY, recibido %d: %s", rawBusyPower.Status, rawBusyPower.RawBody)
		}
		queryDatabase(t, fmt.Sprintf("DELETE FROM tareas_asincronas WHERE id = '%s';", busyTaskID))

		// 3. Rollback con ADMIN devuelve 202 Accepted y completa tarea
		statusRB, rbBody := requestJSON(t, http.MethodPost, fmt.Sprintf("/instances/%s/snapshots/%s/rollback", ref101, snapName), adminToken, nil)
		if statusRB != http.StatusAccepted {
			t.Fatalf("BAC-43: POST rollback esperado 202 Accepted, recibido %d: %#v", statusRB, rbBody)
		}
		rbTaskID := requiredString(t, rbBody, "tareaId")
		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", rbTaskID), "COMPLETED", 15*time.Second)

		// 4. Delete con ADMIN devuelve 202 Accepted y elimina el snapshot
		statusDel, delBody := requestJSON(t, http.MethodDelete, fmt.Sprintf("/instances/%s/snapshots/%s", ref101, snapName), adminToken, nil)
		if statusDel != http.StatusAccepted {
			t.Fatalf("BAC-43: DELETE snapshot esperado 202 Accepted, recibido %d: %#v", statusDel, delBody)
		}
		delTaskID := requiredString(t, delBody, "tareaId")
		waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", delTaskID), "COMPLETED", 15*time.Second)
	})
}

// TestEtapa3BAC44_InterceptorGlobalYConsultasAuditoria verifica BAC-44 (documentacion/actual.md):
//   1. Catálogo de acciones de Etapa 3 en puertos de dominio: SNAPSHOT_CREATE, SNAPSHOT_ROLLBACK,
//      SNAPSHOT_DELETE, ALERT_TRIGGERED, CONFIG_UPDATED.
//   2. Consulta paginada GET /api/admin/audit con soporte de filtros compuestos: desde, hasta,
//      accion, usuarioId, recursoTipo, recursoId, resultado, page y limit.
//   3. Validación de rango temporal: si desde > hasta responde 400 INVALID_REQUEST.
//   4. Orden cronológico inverso (fecha_hora DESC).
func TestEtapa3BAC44_InterceptorGlobalYConsultasAuditoria(t *testing.T) {
	t.Run("puertos, repositorio y handler de auditoría declaran acciones Etapa 3, filtros compuestos (recursoTipo, recursoId, page, limit) y orden DESC", func(t *testing.T) {
		portsCode := secReadGoDir(t, "internal/core/ports")
		auditHandler := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/audit_handler.go")))
		auditRepo := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/postgres/audit_repository.go")))
		combined := portsCode + "\n" + auditHandler + "\n" + auditRepo

		for _, action := range []string{"SNAPSHOT_CREATE", "SNAPSHOT_ROLLBACK", "SNAPSHOT_DELETE", "ALERT_TRIGGERED", "CONFIG_UPDATED"} {
			if !strings.Contains(portsCode, action) {
				t.Errorf("BAC-44: internal/core/ports no declara la constante de acción de auditoría %q", action)
			}
		}

		for _, param := range []string{"desde", "hasta", "accion", "usuarioId", "recursoTipo", "recursoId", "resultado", "page", "limit"} {
			if !strings.Contains(auditHandler, param) {
				t.Errorf("BAC-44: audit_handler.go no procesa el query param %q en GET /api/admin/audit", param)
			}
		}
		if !strings.Contains(auditRepo, "fecha_hora") || !strings.Contains(combined, "DESC") {
			t.Errorf("BAC-44: audit_repository.go / audit_handler.go debe ordenar los registros por fecha_hora DESC")
		}
	})

	t.Run("GET /api/admin/audit filtra por fecha, acción, usuario, recursoTipo, recursoId, page y limit y rechaza desde > hasta", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)

		// Rango invertido desde > hasta debe responder 400
		rawInverted := requestRaw(t, http.MethodGet, "/admin/audit?desde=2026-12-31&hasta=2026-01-01", adminToken, nil)
		if rawInverted.Status != http.StatusBadRequest {
			t.Errorf("BAC-44: GET /admin/audit con desde > hasta esperado 400, recibido %d: %s", rawInverted.Status, rawInverted.RawBody)
		}

		// Consulta filtrada válida con page, limit, recursoTipo y recursoId
		status, page := requestJSON(t, http.MethodGet, "/admin/audit?page=1&limit=10&recursoTipo=INSTANCE&recursoId=101&resultado=EXITO", adminToken, nil)
		if status != http.StatusOK {
			t.Fatalf("BAC-44: GET /admin/audit esperado 200, recibido %d: %#v", status, page)
		}
		if _, ok := page["items"].([]any); !ok {
			t.Errorf("BAC-44: GET /admin/audit debe devolver el arreglo 'items': %#v", page)
		}
		if _, ok := page["total"].(float64); !ok {
			t.Errorf("BAC-44: GET /admin/audit debe devolver el conteo 'total': %#v", page)
		}
	})
}

// TestEtapa3BAC45_StreamingExportacionAuditoriaCSVyJSON verifica BAC-45 (documentacion/actual.md):
//   1. GET /api/admin/audit/export soporta format=csv y format=json (además de formato=csv|json).
//   2. Escritura por streaming (cursor SQL Rows()/ScanRows o c.Stream) sin cargar 100.000 filas en memoria con Find(&rows).
//   3. Cabeceras Content-Type (text/csv o application/json) y Content-Disposition con nombre
//      auditoria_centinela_YYYYMMDD.csv|json.
//   4. Sanitización contra CSV Formula Injection (prefijando celdas que comienzan con =, +, -, @).
//   5. Registro en auditoría de la propia exportación (AUDIT_EXPORT_CSV / AUDIT_EXPORT_JSON / EXPORTAR_AUDITORIA).
func TestEtapa3BAC45_StreamingExportacionAuditoriaCSVyJSON(t *testing.T) {
	t.Run("audit_handler.go y audit_repository.go implementan streaming real, nombre auditoria_centinela_YYYYMMDD y sanitización anti-inyección CSV", func(t *testing.T) {
		auditHandler := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/audit_handler.go")))
		auditRepo := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/postgres/audit_repository.go")))
		auditService := secReadGoDir(t, "internal/core/services")
		auditCode := auditHandler + "\n" + auditRepo + "\n" + auditService

		if !regexp.MustCompile(`auditoria_centinela_`).MatchString(auditCode) {
			t.Errorf("BAC-45: la exportación de auditoría debe nombrar el archivo 'auditoria_centinela_YYYYMMDD.csv|json' en Content-Disposition")
		}
		if !regexp.MustCompile(`\.Rows\(\)|ScanRows|c\.Stream`).MatchString(auditCode) {
			t.Errorf("BAC-45: la exportación debe transmitir registros mediante cursor/streaming (.Rows(), ScanRows o c.Stream) sin cargar toda la tabla en memoria con Find")
		}
		if !regexp.MustCompile(`['"=+\-@]`).MatchString(auditCode) || !regexp.MustCompile(`(?i)sanitiz|formula|escape|prefij`).MatchString(auditCode) {
			// Verificar presencia explícita de los 4 disparadores de inyección de fórmulas (=, +, -, @)
			if !(strings.Contains(auditCode, "'='") && strings.Contains(auditCode, "'+'") && strings.Contains(auditCode, "'-'") && strings.Contains(auditCode, "'@'")) &&
				!(strings.Contains(auditCode, "\"=\"") && strings.Contains(auditCode, "\"+\"") && strings.Contains(auditCode, "\"-\"") && strings.Contains(auditCode, "\"@\"")) {
				t.Errorf("BAC-45: falta sanitización explícita contra inyección de fórmulas CSV (=, +, -, @)")
			}
		}
	})

	t.Run("GET /api/admin/audit/export descarga CSV y JSON válidos con format=csv|json y registra la acción en auditoría", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)

		filenamePatternCSV := regexp.MustCompile(`auditoria_centinela_\d{8}\.csv`)
		filenamePatternJSON := regexp.MustCompile(`auditoria_centinela_\d{8}\.json`)

		// 1. Exportación CSV (?format=csv)
		csvResp := requestRaw(t, http.MethodGet, "/admin/audit/export?format=csv", adminToken, nil)
		if csvResp.Status != http.StatusOK {
			t.Fatalf("BAC-45: GET /admin/audit/export?format=csv esperado 200, recibido %d: %s", csvResp.Status, csvResp.RawBody)
		}
		if !strings.Contains(strings.ToLower(csvResp.Header.Get("Content-Type")), "text/csv") {
			t.Errorf("BAC-45: Content-Type CSV inesperado: %q", csvResp.Header.Get("Content-Type"))
		}
		if !filenamePatternCSV.MatchString(csvResp.Header.Get("Content-Disposition")) {
			t.Errorf("BAC-45: Content-Disposition CSV debe incluir auditoria_centinela_YYYYMMDD.csv, recibido: %q", csvResp.Header.Get("Content-Disposition"))
		}
		reader := csv.NewReader(strings.NewReader(csvResp.RawBody))
		records, err := reader.ReadAll()
		if err != nil || len(records) == 0 {
			t.Errorf("BAC-45: el cuerpo devuelto por /admin/audit/export?format=csv no es un CSV válido (%v): %s", err, csvResp.RawBody)
		}

		// 2. Exportación JSON (?format=json)
		jsonResp := requestRaw(t, http.MethodGet, "/admin/audit/export?format=json", adminToken, nil)
		if jsonResp.Status != http.StatusOK {
			t.Errorf("BAC-45: GET /admin/audit/export?format=json esperado 200, recibido %d: %s", jsonResp.Status, jsonResp.RawBody)
		} else {
			if !strings.Contains(strings.ToLower(jsonResp.Header.Get("Content-Type")), "application/json") {
				t.Errorf("BAC-45: Content-Type JSON inesperado: %q", jsonResp.Header.Get("Content-Type"))
			}
			if !filenamePatternJSON.MatchString(jsonResp.Header.Get("Content-Disposition")) {
				t.Errorf("BAC-45: Content-Disposition JSON debe incluir auditoria_centinela_YYYYMMDD.json, recibido: %q", jsonResp.Header.Get("Content-Disposition"))
			}
			var exportedItems []any
			if err := json.Unmarshal([]byte(jsonResp.RawBody), &exportedItems); err != nil {
				var wrapper map[string]any
				if errObj := json.Unmarshal([]byte(jsonResp.RawBody), &wrapper); errObj != nil {
					t.Errorf("BAC-45: la exportación JSON no devolvió un JSON válido: %v", err)
				}
			}
		}

		// 3. Debe quedar auditada la exportación
		if n := queryDatabase(t, "SELECT count(*) FROM auditoria WHERE accion IN ('AUDIT_EXPORT_CSV', 'AUDIT_EXPORT_JSON', 'EXPORTAR_AUDITORIA', 'AUDIT_EXPORT');"); n == "0" {
			t.Errorf("BAC-45: la exportación de auditoría no dejó rastro en la tabla auditoria")
		}
	})
}

// TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS verifica INF-11 (documentacion/actual.md):
//   1. Redirección HTTP (80) a HTTPS (443) con 301.
//   2. TLS 1.2 y TLS 1.3 habilitados con cifrados seguros.
//   3. Cabeceras HSTS (Strict-Transport-Security con includeSubDomains), X-Frame-Options,
//      X-Content-Type-Options, Referrer-Policy y Content-Security-Policy (CSP con default-src 'self' y wss:).
//   4. Soporte de proxy reverso para WebSockets (/api/ws/ con Upgrade y Connection) y SSE (/api/events con proxy_buffering off).
func TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS(t *testing.T) {
	nginxPath := filepath.Join(envOrDefault("CENTINELA_ROOT", "../.."), "docker", "nginx-edge.conf")
	config := string(readFile(t, nginxPath))

	t.Run("redirección 80 -> 443 y protocolos TLS 1.2 / TLS 1.3", func(t *testing.T) {
		if !regexp.MustCompile(`listen\s+80\b`).MatchString(config) || !regexp.MustCompile(`return\s+301\s+https://`).MatchString(config) {
			t.Errorf("INF-11: docker/nginx-edge.conf no redirige el puerto 80 a HTTPS con 301")
		}
		if !regexp.MustCompile(`ssl_protocols\s+[^;]*TLSv1\.3`).MatchString(config) {
			t.Errorf("INF-11: docker/nginx-edge.conf debe habilitar explícitamente TLSv1.3 en ssl_protocols")
		}
	})

	t.Run("cabeceras de endurecimiento HSTS (includeSubDomains), CSP (con wss:), X-Frame-Options y X-Content-Type-Options", func(t *testing.T) {
		requiredHeaders := map[string]*regexp.Regexp{
			"Strict-Transport-Security": regexp.MustCompile(`(?i)add_header\s+Strict-Transport-Security\s+"max-age=\d+;\s*includeSubDomains`),
			"X-Frame-Options":           regexp.MustCompile(`(?i)add_header\s+X-Frame-Options\s+"?(DENY|SAMEORIGIN)"?`),
			"X-Content-Type-Options":    regexp.MustCompile(`(?i)add_header\s+X-Content-Type-Options\s+"?nosniff"?`),
			"Content-Security-Policy":   regexp.MustCompile(`(?i)add_header\s+Content-Security-Policy\s+"[^"]*default-src\s+'self'[^"]*wss:`),
		}
		for header, pattern := range requiredHeaders {
			if !pattern.MatchString(config) {
				t.Errorf("INF-11: falta o es incompleta la cabecera %s en docker/nginx-edge.conf", header)
			}
		}
	})

	t.Run("soporte de Upgrade WebSocket (/api/ws/) y streaming SSE sin buffering (/api/events)", func(t *testing.T) {
		if !regexp.MustCompile(`(?i)proxy_set_header\s+Upgrade\s+\$http_upgrade`).MatchString(config) {
			t.Errorf("INF-11: docker/nginx-edge.conf no reenvía la cabecera Upgrade para WebSockets (/api/ws/)")
		}
		if !regexp.MustCompile(`(?i)proxy_buffering\s+off`).MatchString(config) {
			t.Errorf("INF-11: docker/nginx-edge.conf no deshabilita proxy_buffering para SSE (/api/events)")
		}
	})
}

// TestEtapa3BAC46_WorkerEvaluadorUmbralesHisteresis verifica BAC-46 (documentacion/actual.md):
//   1. Worker en segundo plano que evalúa umbrales de CPU, RAM y Almacenamiento del host.
//   2. Histéresis (múltiples muestras consecutivas por encima del umbral) y período de enfriamiento
//      (cooldown en Redis) para evitar tormentas de alertas.
//   3. Emisión de evento SSE ALERT_TRIGGERED y registro en auditoría con acción ALERT_TRIGGERED.
func TestEtapa3BAC46_WorkerEvaluadorUmbralesHisteresis(t *testing.T) {
	servicesCode := secReadGoDir(t, "internal/core/services")
	mainCode := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))
	combined := servicesCode + "\n" + mainCode

	for _, marker := range []string{"ALERT_TRIGGERED", "ALERT_CPU_THRESHOLD", "ALERT_RAM_THRESHOLD", "ALERT_STORAGE_THRESHOLD"} {
		if !strings.Contains(combined, marker) {
			t.Errorf("BAC-46: el backend no contiene %q para el worker evaluador de umbrales con histéresis", marker)
		}
	}

	_, testFiles := goFilesContaining(t, `ALERT_CPU_THRESHOLD|ALERT_RAM_THRESHOLD|ALERT_STORAGE_THRESHOLD|Histeresis`)
	if len(testFiles) == 0 {
		t.Errorf("BAC-46: no se encontraron pruebas unitarias del evaluador de umbrales e histéresis en el backend")
	} else {
		runBackendUnitTests(t, testFiles)
	}
}
