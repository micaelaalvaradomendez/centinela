package back_test

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Hito "Control de Acceso Basado en Recursos" (documentacion/actual.md):
//   - BAC-07 asignación atómica de permisos por instancia (terminado.md, regresión)
//   - FIX-16 / BAC-08 guard RequireInstanceAccess montado en /instances/:vmid
//   - BAC-14 GET /api/instances normalizado y filtrado por rol
//   - SEC-04 niveles de acceso FULL_ACCESS / READ_ONLY
//
// Proxmox VE se reemplaza por proxmox-stub (nginx) con un inventario fijo:
// 101 qemu "web-01" running, 102 lxc "db-01" stopped, 103 qemu "ci-runner" running,
// más una entrada de nodo y otra de storage que el backend debe descartar.
func TestHitoControlDeAccesoBasadoEnRecursos(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	suffix := time.Now().UnixNano()
	operatorID, _ := createUserThroughAPI(t, adminToken, fmt.Sprintf("op_recursos_%d", suffix), fmt.Sprintf("op.recursos.%d@elcentinela.com", suffix), "OPERATOR")
	operatorToken := signedAccessToken(t, operatorID, "OPERATOR", orgID)
	permEndpoint := "/admin/users/" + operatorID + "/permissions"

	ref101 := instanceRefByVMID(t, adminToken, 101)
	ref102 := instanceRefByVMID(t, adminToken, 102)
	ref103 := instanceRefByVMID(t, adminToken, 103)
	vmidToRef := map[int]string{101: ref101, 102: ref102, 103: ref103}
	refToVmid := map[string]int{ref101: 101, ref102: 102, ref103: 103}

	// Contrato canónico desde SEC-04: { permisos: [{ vmid, nivelAcceso? }] }.
	assign := func(t *testing.T, vmids ...int) {
		t.Helper()
		permisos := make([]map[string]any, 0, len(vmids))
		for _, vmid := range vmids {
			vRef := vmidToRef[vmid]
			if vRef == "" {
				vRef = fmt.Sprint(vmid)
			}
			permisos = append(permisos, map[string]any{"vmid": vRef})
		}
		status, body := requestValue(t, http.MethodPut, permEndpoint, adminToken, map[string]any{"permisos": permisos})
		if status != http.StatusNoContent && status != http.StatusOK {
			t.Fatalf("PUT %s esperado 204, recibido %d: %#v", permEndpoint, status, body)
		}
	}
	readAssigned := func(t *testing.T) []int {
		t.Helper()
		status, body := requestJSON(t, http.MethodGet, permEndpoint, adminToken, nil)
		if status != http.StatusOK {
			t.Fatalf("GET %s esperado 200, recibido %d: %#v", permEndpoint, status, body)
		}
		levels := permissionLevels(t, body, refToVmid)
		vmids := make([]int, 0, len(levels))
		for vmid := range levels {
			vmids = append(vmids, vmid)
		}
		sort.Ints(vmids)
		return vmids
	}

	t.Run("BAC-07 GET y PUT /api/admin/users/:id/permissions reemplazan el conjunto de forma atomica", func(t *testing.T) {
		auditBefore := auditCount(t, adminID, "ASIGNAR_PERMISOS")

		assign(t, 101, 102)
		if got := readAssigned(t); fmt.Sprint(got) != "[101 102]" {
			t.Errorf("tras asignar [101 102] GET devolvió %v", got)
		}

		assign(t, 102)
		if got := readAssigned(t); fmt.Sprint(got) != "[102]" {
			t.Errorf("el PUT debe REEMPLAZAR el conjunto: tras asignar [102] GET devolvió %v", got)
		}

		assign(t, 101, 101)
		if rows := queryDatabase(t, "SELECT count(*) FROM permisos_instancia WHERE usuario_id = '"+operatorID+"' AND vmid_proxmox = 101;"); rows != "1" {
			t.Errorf("asignar [101 101] duplicó la fila en la base; filas encontradas: %s", rows)
		}

		for _, method := range []string{http.MethodGet, http.MethodPut} {
			status, _ := requestValue(t, method, permEndpoint, operatorToken, permissionsPayload(101))
			if status != http.StatusForbidden {
				t.Errorf("%s %s con OPERATOR: esperado 403, recibido %d", method, permEndpoint, status)
			}
		}

		if after := auditCount(t, adminID, "ASIGNAR_PERMISOS"); after < auditBefore+3 {
			t.Errorf("BAC-18: se hicieron 3 asignaciones pero la auditoría ASIGNAR_PERMISOS pasó de %d a %d", auditBefore, after)
		}
	})

	t.Run("FIX-16 BAC-08 rechaza instancias no asignadas con 403 INSTANCE_ACCESS_DENIED sin llamar a Proxmox", func(t *testing.T) {
		assign(t, 101)

		requests := []struct{ method, path string }{
			{http.MethodGet, "/instances/" + ref103},
			{http.MethodPost, "/instances/" + ref103 + "/start"},
			{http.MethodPost, "/instances/" + ref103 + "/stop"},
		}
		for _, r := range requests {
			status, body := requestJSON(t, r.method, r.path, operatorToken, nil)
			if status != http.StatusForbidden {
				t.Errorf("%s %s (no asignada) con OPERATOR: esperado 403, recibido %d: %#v", r.method, r.path, status, body)
				continue
			}
			if body["errorCode"] != "INSTANCE_ACCESS_DENIED" {
				t.Errorf("%s %s: errorCode esperado INSTANCE_ACCESS_DENIED, recibido %#v", r.method, r.path, body["errorCode"])
			}
		}

		// El guard debe cortar ANTES de Proxmox: el stub no puede haber recibido la orden sobre la 103.
		stubLog, _ := runCompose("logs", "--no-color", "proxmox")
		if strings.Contains(stubLog, "/qemu/103/status/") {
			t.Errorf("Proxmox recibió una orden sobre la VM 103 aunque el OPERATOR no la tiene asignada")
		}

		// El guard no debe bloquear de más: la instancia asignada y el ADMIN pasan.
		status, body := requestJSON(t, http.MethodGet, "/instances/"+ref101, operatorToken, nil)
		if status != http.StatusOK || (body["vmid"] != float64(101) && body["vmid"] != ref101) {
			t.Errorf("GET /instances/%s (asignada) con OPERATOR: esperado 200 con vmid 101, recibido %d: %#v", ref101, status, body)
		}
		status, body = requestJSON(t, http.MethodGet, "/instances/"+ref103, adminToken, nil)
		if status != http.StatusOK {
			t.Errorf("GET /instances/%s con ADMIN: esperado 200, recibido %d: %#v", ref103, status, body)
		}

		// VMIDs de infraestructura (commit 725d436, PROXMOX_PROTECTED_VMIDS=103 en compose.yaml):
		// ni el ADMIN puede apagarlos, y Proxmox no recibe la orden.
		status, body = requestJSON(t, http.MethodPost, "/instances/"+ref103+"/stop", adminToken, nil)
		if status != http.StatusForbidden || body["errorCode"] != "INSTANCE_PROTECTED" {
			t.Errorf("POST /instances/%s/stop (VMID protegido) con ADMIN: esperado 403 INSTANCE_PROTECTED, recibido %d: %#v", ref103, status, body)
		}
		if stubLog, _ := runCompose("logs", "--no-color", "proxmox"); strings.Contains(stubLog, "/qemu/103/status/stop") {
			t.Errorf("Proxmox recibió la orden de apagar la VM protegida 103")
		}
	})

	t.Run("BAC-14 GET /api/instances normaliza el inventario y lo filtra por permisos", func(t *testing.T) {
		status, _ := requestValue(t, http.MethodGet, "/instances", "", nil)
		if status != http.StatusUnauthorized {
			t.Errorf("GET /api/instances sin token: esperado 401, recibido %d", status)
		}

		adminStatus, adminList := requestJSONArray(t, http.MethodGet, "/instances", adminToken)
		if adminStatus != http.StatusOK {
			t.Fatalf("GET /api/instances con ADMIN: esperado 200, recibido %d: %#v", adminStatus, adminList)
		}
		byID := map[int]map[string]any{}
		for _, item := range adminList {
			instance, _ := item.(map[string]any)
			switch idVal := instance["id"].(type) {
			case float64:
				byID[int(idVal)] = instance
			case string:
				if vmid, ok := refToVmid[idVal]; ok {
					byID[vmid] = instance
				}
			}
		}
		if len(adminList) != 3 || byID[101] == nil || byID[102] == nil || byID[103] == nil {
			t.Fatalf("ADMIN debe recibir exactamente las 3 instancias (101, 102, 103) sin nodos ni storages; recibió %#v", adminList)
		}
		expected := map[int]map[string]any{
			101: {"name": "web-01", "type": "vm", "node": "pve", "status": "running"},
			102: {"name": "db-01", "type": "lxc", "node": "pve", "status": "stopped"},
		}
		for id, fields := range expected {
			for field, want := range fields {
				if byID[id][field] != want {
					t.Errorf("instancia %d: campo %q esperado %q, recibido %#v", id, field, want, byID[id][field])
				}
			}
		}

		assign(t, 101)
		_, _ = execDatabase(t, fmt.Sprintf("INSERT INTO permisos_instancia (id, organizacion_id, usuario_id, vmid_proxmox, nivel_acceso) VALUES (gen_random_uuid(), '%s', '%s', 9999, 'FULL_ACCESS') ON CONFLICT DO NOTHING;", orgID, operatorID))
		opStatus, opList := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if opStatus != http.StatusOK {
			t.Fatalf("GET /api/instances con OPERATOR: esperado 200, recibido %d: %#v", opStatus, opList)
		}
		if len(opList) != 1 || (opList[0].(map[string]any)["id"] != float64(101) && opList[0].(map[string]any)["id"] != ref101) {
			t.Errorf("OPERATOR con permisos [101 9999] debe ver solo la 101 (9999 no existe en Proxmox); recibió %#v", opList)
		}

		assign(t)
		_, emptyList := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if len(emptyList) != 0 {
			t.Errorf("OPERATOR sin permisos no debe descubrir instancias; recibió %#v", emptyList)
		}

		// Los errores de Proxmox se traducen a respuestas HTTP controladas.
		nfStatus, nfBody := requestJSON(t, http.MethodGet, "/instances/424242", adminToken, nil)
		if nfStatus != http.StatusNotFound || (nfBody["errorCode"] != "INSTANCE_NOT_FOUND" && nfBody["errorCode"] != "RESOURCE_NOT_FOUND") {
			t.Errorf("VMID inexistente en Proxmox: esperado 404 INSTANCE_NOT_FOUND o RESOURCE_NOT_FOUND, recibido %d: %#v", nfStatus, nfBody)
		}
	})

	t.Run("SEC-04 niveles de acceso FULL_ACCESS y READ_ONLY sobre instancias", func(t *testing.T) {
		column := queryDatabase(t, `SELECT coalesce(column_default, '') FROM information_schema.columns WHERE table_name = 'permisos_instancia' AND column_name = 'nivel_acceso';`)
		if column == "" {
			t.Fatalf("SEC-04 no implementada: la tabla permisos_instancia no tiene la columna nivel_acceso (entregable 1)")
		}
		if !strings.Contains(column, "FULL_ACCESS") {
			t.Errorf("nivel_acceso debe tener default 'FULL_ACCESS', tiene %q", column)
		}

		// Sin nivel explícito se asume FULL_ACCESS.
		assign(t, 101)
		if levels := queryDatabase(t, "SELECT string_agg(DISTINCT nivel_acceso, ',') FROM permisos_instancia WHERE usuario_id = '"+operatorID+"';"); levels != "FULL_ACCESS" {
			t.Errorf("un permiso sin nivelAcceso debe persistirse como FULL_ACCESS; niveles persistidos: %q", levels)
		}
		if _, err := execDatabase(t, "UPDATE permisos_instancia SET nivel_acceso = 'SUPERUSER' WHERE usuario_id = '"+operatorID+"';"); err == nil {
			t.Errorf("la base aceptó nivel_acceso = 'SUPERUSER'; se esperaba CHECK/enum con FULL_ACCESS y READ_ONLY")
		}

		// Nivel explícito por API y lectura del nivel en GET.
		status, body := requestValue(t, http.MethodPut, permEndpoint, adminToken, map[string]any{"permisos": []map[string]any{
			{"vmid": ref101, "nivelAcceso": "FULL_ACCESS"},
			{"vmid": ref102, "nivelAcceso": "READ_ONLY"},
		}})
		if status != http.StatusNoContent {
			t.Fatalf("PUT permisos con niveles explícitos: esperado 204, recibido %d: %#v", status, body)
		}
		_, got := requestJSON(t, http.MethodGet, permEndpoint, adminToken, nil)
		if levels := permissionLevels(t, got, refToVmid); levels[101] != "FULL_ACCESS" || levels[102] != "READ_ONLY" {
			t.Errorf("GET permissions debe devolver el nivel de cada vmid; recibido %#v", got)
		}

		// El guard aplica el nivel: READ_ONLY consulta pero no opera; FULL_ACCESS opera.
		if status, body := requestJSON(t, http.MethodGet, "/instances/"+ref102, operatorToken, nil); status != http.StatusOK {
			t.Errorf("READ_ONLY sobre la 102 debe permitir GET /instances/%s: esperado 200, recibido %d: %#v", ref102, status, body)
		}
		for _, action := range []string{"start", "stop"} {
			if status, body := requestJSON(t, http.MethodPost, "/instances/"+ref102+"/"+action, operatorToken, nil); status != http.StatusForbidden || body["errorCode"] != "INSTANCE_ACCESS_DENIED" {
				t.Errorf("READ_ONLY sobre la 102 no debe permitir %s: esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", action, status, body)
			}
		}
		// La 101 está running en el stub: desde BAC-24A la orden válida es stop (start daría 409 por estado).
		if status, body := requestJSON(t, http.MethodPost, "/instances/"+ref101+"/stop", operatorToken, nil); status != http.StatusAccepted {
			t.Errorf("FULL_ACCESS sobre la 101 debe permitir operar (stop): esperado 202, recibido %d: %#v", status, body)
		}
		// El inventario del operador incluye también las instancias de solo lectura.
		_, list := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if len(list) != 2 {
			t.Errorf("el OPERATOR con permisos sobre 101 (FULL) y 102 (READ_ONLY) debe ver ambas en el inventario; recibió %#v", list)
		}
	})

	t.Run("SEC-04 el payload anterior { vmids } se rechaza y el contrato es { permisos }", func(t *testing.T) {
		// Decisión del 28/09/2026: el contrato oficial es { permisos: [{ vmid, nivelAcceso? }] };
		// no se mantiene la retrocompatibilidad con { vmids } (FIX-26 descartado).
		assign(t, 101)

		status, body := requestJSON(t, http.MethodPut, permEndpoint, adminToken, map[string]any{"vmids": []int{103}})
		if status != http.StatusBadRequest || body["errorCode"] != "INVALID_REQUEST" {
			t.Errorf("PUT con el payload anterior { vmids }: esperado 400 INVALID_REQUEST, recibido %d: %#v", status, body)
		}
		_, got := requestJSON(t, http.MethodGet, permEndpoint, adminToken, nil)
		if levels := permissionLevels(t, got, refToVmid); len(levels) != 1 || levels[101] != "FULL_ACCESS" {
			t.Errorf("un PUT rechazado no debe modificar los permisos existentes; se esperaba solo 101 FULL_ACCESS, recibido %#v", got)
		}

		if status, body := requestValue(t, http.MethodPut, permEndpoint, adminToken, map[string]any{"permisos": []any{}}); status != http.StatusNoContent {
			t.Errorf("PUT { permisos: [] } debe quitar todos los permisos con 204, recibido %d: %#v", status, body)
		} else if rows := queryDatabase(t, "SELECT count(*) FROM permisos_instancia WHERE usuario_id = '"+operatorID+"';"); rows != "0" {
			t.Errorf("PUT { permisos: [] } dejó %s permisos; se esperaba 0", rows)
		}
	})

	t.Run("FIX-37 FRN-18 GET /account/profile expone el nivel de acceso por instancia para canOperateInstance", func(t *testing.T) {
		// El frontend (FIX-30) arma canOperateInstance con perfil.permisos: [{ vmid, nivelAcceso }].
		if status, body := requestValue(t, http.MethodPut, permEndpoint, adminToken, map[string]any{"permisos": []map[string]any{
			{"vmid": ref101, "nivelAcceso": "FULL_ACCESS"},
			{"vmid": ref102, "nivelAcceso": "READ_ONLY"},
		}}); status != http.StatusNoContent {
			t.Fatalf("PUT permisos con niveles explícitos: esperado 204, recibido %d: %#v", status, body)
		}
		status, profile := requestJSON(t, http.MethodGet, "/account/profile", operatorToken, nil)
		if status != http.StatusOK {
			t.Fatalf("GET /account/profile del OPERATOR: esperado 200, recibido %d: %#v", status, profile)
		}
		if _, ok := profile["permisos"]; !ok {
			t.Fatalf("GET /account/profile no incluye permisos [{ vmid, nivelAcceso }]: el frontend no puede distinguir FULL_ACCESS de READ_ONLY y canOperateInstance da false para todo OPERATOR (claves recibidas: %v)", mapKeys(profile))
		}
		if levels := permissionLevels(t, profile, refToVmid); levels[101] != "FULL_ACCESS" || levels[102] != "READ_ONLY" {
			t.Errorf("GET /account/profile debe informar 101 FULL_ACCESS y 102 READ_ONLY; recibido %#v", profile["permisos"])
		}
	})
}

func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// permissionsPayload arma el body canónico de PUT /admin/users/:id/permissions sin nivel
// explícito (el backend asume FULL_ACCESS).
func permissionsPayload(vmids ...int) map[string]any {
	permisos := make([]map[string]any, 0, len(vmids))
	for _, vmid := range vmids {
		permisos = append(permisos, map[string]any{"vmid": vmid})
	}
	return map[string]any{"permisos": permisos}
}

// permissionLevels interpreta GET /admin/users/:id/permissions -> { permisos: [{ vmid, nivelAcceso }] }.
func permissionLevels(t *testing.T, body map[string]any, refToVmid ...map[string]int) map[int]string {
	t.Helper()
	raw, ok := body["permisos"].([]any)
	if !ok {
		t.Fatalf("GET permissions debe devolver { permisos: [{ vmid, nivelAcceso }] }, recibido %#v", body)
	}
	var mapper map[string]int
	if len(refToVmid) > 0 {
		mapper = refToVmid[0]
	}
	levels := map[int]string{}
	for _, item := range raw {
		permission, _ := item.(map[string]any)
		level, _ := permission["nivelAcceso"].(string)
		switch v := permission["vmid"].(type) {
		case float64:
			levels[int(v)] = level
		case string:
			if mapper != nil {
				if vmid, ok := mapper[v]; ok {
					levels[vmid] = level
					continue
				}
			}
			if num, err := strconv.Atoi(v); err == nil {
				levels[num] = level
			}
		}
	}
	return levels
}
