package back_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Tareas de actual.md sin cobertura previa.
//
// BAC-26 (RF-11 parcial): el canal SSE de BAC-21C ya existe; esta tarea verifica y
// documenta la emisión de TASK_FINISHED para las acciones de ciclo de vida, con el
// payload de BAC-29 y el filtro por permiso. El UPID crudo no viaja al frontend.
func TestEtapa1Ola3BAC26(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	// OPERATOR con FULL_ACCESS sobre la 101 (el stub la reporta encendida, así que el
	// conjunto de órdenes válidas para probar el evento es stop/shutdown/reboot).
	conPermiso := createActiveUser(t, adminToken, "bac26_con", "OPERATOR")
	if status, body := requestValue(t, http.MethodPut, "/admin/users/"+conPermiso.ID+"/permissions", adminToken, map[string]any{
		"permisos": []map[string]any{{"vmid": 101, "nivelAcceso": "FULL_ACCESS"}},
	}); status != http.StatusNoContent {
		t.Fatalf("precondición BAC-26: asignar FULL_ACCESS sobre 101, esperado 204, recibido %d: %#v", status, body)
	}
	// OPERATOR sin ningún permiso: no debe recibir eventos de la 101.
	sinPermiso := createActiveUser(t, adminToken, "bac26_sin", "OPERATOR")

	sessionCon := loginWithTOTP(t, conPermiso.Email, conPermiso.Password, conPermiso.Secret)
	sessionSin := loginWithTOTP(t, sinPermiso.Email, sinPermiso.Password, sinPermiso.Secret)

	eventsCon := listenEvents(t, sessionCon.AccessToken)
	eventsSin := listenEvents(t, sessionSin.AccessToken)

	for _, action := range []string{"stop", "shutdown", "reboot"} {
		status, body := requestJSON(t, http.MethodPost, "/instances/101/status/"+action, sessionCon.AccessToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("BAC-26: POST /instances/101/status/%s: esperado 202, recibido %d: %#v", action, status, body)
		}
		tareaID := requiredString(t, body, "tareaId")

		event := waitTaskFinished(t, eventsCon, tareaID, 8*time.Second)
		if event == nil {
			t.Errorf("BAC-26: no se recibió TASK_FINISHED para la acción %q (tarea %s)", action, tareaID)
			continue
		}
		if event["tipo"] != "TASK_FINISHED" {
			t.Errorf("BAC-26: tipo esperado TASK_FINISHED, recibido %v", event["tipo"])
		}
		if fmt.Sprint(event["recursoId"]) != "101" {
			t.Errorf("BAC-26: recursoId esperado 101, recibido %v", event["recursoId"])
		}
		detalles, _ := event["detalles"].(map[string]any)
		for _, field := range []string{"tareaId", "accion", "estado", "exitstatus"} {
			if value, present := detalles[field]; !present || fmt.Sprint(value) == "" {
				t.Errorf("BAC-26: el payload de TASK_FINISHED (%s) no incluye %q (contrato BAC-29): %#v", action, field, detalles)
			}
		}
		if detalles["tareaId"] != tareaID {
			t.Errorf("BAC-26: detalles.tareaId esperado %s, recibido %v", tareaID, detalles["tareaId"])
		}
		if detalles["estado"] != "COMPLETED" {
			t.Errorf("BAC-26 (%s): detalles.estado esperado COMPLETED, recibido %v", action, detalles["estado"])
		}
	}

	// RBAC: un OPERATOR sin permiso sobre la 101 no recibe su TASK_FINISHED.
	select {
	case ev := <-eventsSin:
		detalles, _ := ev["detalles"].(map[string]any)
		if ev["tipo"] == "TASK_FINISHED" && fmt.Sprint(ev["recursoId"]) == "101" {
			t.Errorf("BAC-26: un OPERATOR sin permiso recibió TASK_FINISHED de la 101: %#v (tarea %v)", ev, detalles["tareaId"])
		}
	case <-time.After(1 * time.Second):
		// Correcto: no llegó ningún evento de la 101.
	}
}

// INF-09 (BRG-06-INF, RNF-06): soporte de WebSockets en el borde Nginx.
// La telemetría de la Etapa 2 exige encabezados Upgrade/Connection y un timeout
// extendido para no cortar sesiones de monitoreo. El archivo versionado del borde
// es la referencia desplegada (FIX-31).
func TestINF09WebSocketsNginxEdge(t *testing.T) {
	conf := string(readFile(t, "../../docker/nginx-edge.conf"))

	requisitos := []struct {
		nombre string
		re     *regexp.Regexp
	}{
		{"proxy_http_version 1.1", regexp.MustCompile(`proxy_http_version\s+1\.1\s*;`)},
		{"proxy_set_header Upgrade $http_upgrade", regexp.MustCompile(`proxy_set_header\s+Upgrade\s+\$http_upgrade\s*;`)},
		{"proxy_set_header Connection upgrade", regexp.MustCompile(`proxy_set_header\s+Connection\s+"?upgrade"?\s*;|proxy_set_header\s+Connection\s+\$connection_upgrade\s*;`)},
		{"proxy_read_timeout largo (>= 3600 s)", regexp.MustCompile(`proxy_read_timeout\s+(?:3600s|1h|[4-9][0-9]{3})\s*;`)},
	}
	for _, requisito := range requisitos {
		if !requisito.re.MatchString(conf) {
			t.Errorf("INF-09: docker/nginx-edge.conf no incluye %q (soporte de WebSockets bidireccionales)", requisito.nombre)
		}
	}

	// El soporte se declara en un location de WebSockets (/api/ws o /ws) o en el catch-all /api/.
	if !regexp.MustCompile(`location\s+(?:=|\^~|~)?\s*/api/ws|location\s+(?:=|\^~|~)?\s*/ws|location\s+/api/`).MatchString(conf) {
		t.Errorf("INF-09: no hay un location que cubra las rutas de WebSockets (/api/ws o /ws) ni el catch-all /api/")
	}
	// El handshake HTTP 101 requiere que Nginx reenvíe la cabecera Upgrade.
	if !strings.Contains(conf, "upgrade") {
		t.Errorf("INF-09: la configuración del borde no menciona la cabecera de upgrade de WebSockets")
	}
}

// BAC-31 (RF-07): endpoints de inventario de storages y plantillas/ISOs.
// Entregable: GET /api/node/storages y GET /api/node/templates, ambos solo para ADMIN.
func TestEtapa1Ola3BAC31(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	operator := createActiveUser(t, adminToken, "bac31_op", "OPERATOR")
	operatorSession := loginWithTOTP(t, operator.Email, operator.Password, operator.Secret)

	for _, path := range []string{"/node/storages", "/node/templates"} {
		status, body := requestValue(t, http.MethodGet, path, adminToken, nil)
		if status != http.StatusOK {
			t.Errorf("BAC-31: GET %s con ADMIN: esperado 200, recibido %d: %#v", path, status, body)
			continue
		}
		// El listado se acepta como array o como objeto con items/data/lista.
		switch value := body.(type) {
		case []any:
			if len(value) == 0 {
				t.Errorf("BAC-31: GET %s devolvió un listado vacío: %#v", path, value)
			}
		case map[string]any:
			if _, ok := value["items"]; !ok {
				if _, ok := value["data"]; !ok {
					t.Errorf("BAC-31: GET %s no devuelve un listado reconocible: %#v", path, value)
				}
			}
		default:
			t.Errorf("BAC-31: GET %s devolvió un tipo inesperado: %#v", path, body)
		}

		// Un OPERATOR no puede consultar el inventario del hipervisor.
		if opStatus, _ := requestValue(t, http.MethodGet, path, operatorSession.AccessToken, nil); opStatus != http.StatusForbidden {
			t.Errorf("BAC-31: GET %s con OPERATOR: esperado 403, recibido %d", path, opStatus)
		}
	}
}

// BAC-33 (RF-02 / RF-07): cuotas asignables del host (headroom de seguridad).
// Entregable: GET /api/node/status (o /node/quota) expone `allocatable`
// { maxCpuCores, freeRamMb, freeStorageGb, headroomPercent } con 15 % de RAM reservada.
func TestEtapa1Ola3BAC33(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	// El contrato permite exponerlo en /node/status o en un endpoint auxiliar /node/quota.
	status, body := requestJSON(t, http.MethodGet, "/node/status", adminToken, nil)
	allocatable, _ := body["allocatable"].(map[string]any)
	if allocatable == nil {
		quotaStatus, quota := requestJSON(t, http.MethodGet, "/node/quota", adminToken, nil)
		if quotaStatus == http.StatusOK {
			allocatable, _ = quota["allocatable"].(map[string]any)
			if allocatable == nil {
				allocatable = quota
			}
		}
	}
	if status != http.StatusOK {
		t.Fatalf("BAC-33: GET /node/status: esperado 200, recibido %d: %#v", status, body)
	}
	if allocatable == nil {
		t.Fatalf("BAC-33: GET /node/status no expone 'allocatable' (ni /node/quota): %#v", body)
	}
	for _, field := range []string{"maxCpuCores", "freeRamMb", "freeStorageGb", "headroomPercent"} {
		if _, present := allocatable[field]; !present {
			t.Errorf("BAC-33: allocatable no incluye %q: %#v", field, allocatable)
		}
	}
	if headroom, ok := allocatable["headroomPercent"].(float64); ok && headroom != 15 {
		t.Errorf("BAC-33: headroomPercent esperado 15 (15%% de RAM reservada), recibido %v", allocatable["headroomPercent"])
	}
}

// BAC-30 (RF-05 / RF-07): el simulador de Proxmox debe soportar creación, edición de
// config, storages/contenido y rrddata para desarrollar la Etapa 2 sin el hipervisor real.
// Se prueba contra el simulador real (servicio proxmox-sim de compose.yaml).
func TestEtapa1Ola3BAC30(t *testing.T) {
	requireIntegration(t)

	base := envOrDefault("BACKEND_TEST_SIM_URL", "http://127.0.0.1:18083/api2/json")
	client := &http.Client{Timeout: 10 * time.Second}
	const token = "PVEAPIToken=stub@pve!tests=stub-secret"

	call := func(method, path string, payload any) (int, map[string]any) {
		t.Helper()
		var body io.Reader
		if payload != nil {
			encoded, _ := json.Marshal(payload)
			body = bytes.NewReader(encoded)
		}
		request, err := http.NewRequest(method, base+path, body)
		if err != nil {
			t.Fatalf("BAC-30: no se pudo crear %s %s: %v", method, path, err)
		}
		request.Header.Set("Authorization", token)
		if payload != nil {
			request.Header.Set("Content-Type", "application/json")
		}
		response, err := client.Do(request)
		if err != nil {
			t.Fatalf("BAC-30: %s %s contra el simulador falló: %v", method, path, err)
		}
		defer response.Body.Close()
		raw, _ := io.ReadAll(response.Body)
		decoded := map[string]any{}
		_ = json.Unmarshal(raw, &decoded)
		return response.StatusCode, decoded
	}

	// 1. Reserva de VMID libre.
	if status, body := call(http.MethodGet, "/cluster/nextid", nil); status != http.StatusOK {
		t.Errorf("BAC-30: GET /cluster/nextid: esperado 200, recibido %d: %#v", status, body)
	}

	// 2. Creación de un contenedor (vmid >= 106, fuera del rango protegido).
	const vmid = 161
	status, created := call(http.MethodPost, "/nodes/pve/lxc", map[string]any{
		"vmid": vmid, "hostname": "bac30-ct", "memory": 512, "cores": 1,
		"ostemplate": "local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst",
		"rootfs":     "local-lvm:8",
	})
	if status != http.StatusOK {
		t.Errorf("BAC-30: POST /nodes/pve/lxc: esperado 200 con UPID, recibido %d: %#v", status, created)
	} else if upid, _ := created["data"].(string); !strings.HasPrefix(upid, "UPID:") || !strings.Contains(upid, "create") {
		t.Errorf("BAC-30: el UPID de creación debe ser del tipo vzcreate, recibido %v", created["data"])
	}

	// 3. Edición de recursos (PUT config) y relectura.
	if status, body := call(http.MethodPut, fmt.Sprintf("/nodes/pve/lxc/%d/config", vmid), map[string]any{"memory": 1024, "cores": 2}); status != http.StatusOK {
		t.Errorf("BAC-30: PUT /nodes/pve/lxc/%d/config: esperado 200, recibido %d: %#v", vmid, status, body)
	}
	if status, body := call(http.MethodGet, fmt.Sprintf("/nodes/pve/lxc/%d/config", vmid), nil); status != http.StatusOK {
		t.Errorf("BAC-30: GET config tras editar: esperado 200, recibido %d", status)
	} else if config, _ := body["data"].(map[string]any); fmt.Sprint(config["memory"]) != "1024" {
		t.Errorf("BAC-30: la memoria editada debería ser 1024, config devuelve %#v", config["memory"])
	}

	// 4. Métricas rrdata (RF-05): array de puntos.
	if status, rrd := call(http.MethodGet, fmt.Sprintf("/nodes/pve/lxc/%d/rrddata?timeframe=hour", vmid), nil); status != http.StatusOK {
		t.Errorf("BAC-30: GET rrddata: esperado 200, recibido %d: %#v", status, rrd)
	} else if _, ok := rrd["data"].([]any); !ok {
		t.Errorf("BAC-30: rrddata debe devolver un array de puntos: %#v", rrd["data"])
	}

	// 5. Storages y contenido (plantillas/ISOs) para el asistente de creación (RF-07).
	for _, path := range []string{"/nodes/pve/storage", "/nodes/pve/storage/local/content"} {
		if status, body := call(http.MethodGet, path, nil); status != http.StatusOK {
			t.Errorf("BAC-30: GET %s: esperado 200, recibido %d: %#v", path, status, body)
		}
	}
}

// FIX-71 (INF-06 / BAC-17A): el backend debe tolerar el arranque en frío de Redis.
// Se reinicia el contenedor de Redis y se verifica que una operación que depende de él
// (el ticket efímero de /api/events vive en Redis, BAC-21C) se recupera en pocos segundos.
func TestFIX71RedisArranqueEnFrio(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	ticketOK := func() bool {
		status, _ := requestJSON(t, http.MethodPost, "/events/ticket", adminToken, nil)
		return status == http.StatusOK || status == http.StatusCreated
	}
	if !ticketOK() {
		t.Fatalf("FIX-71: precondición: POST /events/ticket debe funcionar con Redis arriba")
	}

	if output, err := runCompose("restart", "redis"); err != nil {
		t.Fatalf("FIX-71: no se pudo reiniciar el contenedor redis: %v\n%s", err, output)
	}

	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if ticketOK() {
			t.Logf("FIX-71: el backend se recompuso tras el arranque en frío de Redis")
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Errorf("FIX-71: el backend no volvió a responder operaciones con Redis tras el reinicio (arranque en frío no tolerado)")
}

// INF-07C (RF-07 / RF-10): la documentación de la API de Proxmox debe reflejar los
// permisos que el token necesita para aprovisionar y reconfigurar en la Etapa 2.
func TestINF07CPermisosTokenProxmox(t *testing.T) {
	doc := string(readFile(t, "../../documentacion/api-proxmox.md"))
	for _, permiso := range []string{
		"VM.Allocate",
		"Datastore.AllocateSpace",
		"VM.Config.CPU",
		"VM.Config.Memory",
		"VM.Config.Disk",
		"VM.Config.Network",
	} {
		if !strings.Contains(doc, permiso) {
			t.Errorf("INF-07C: documentacion/api-proxmox.md no documenta el permiso %q requerido para el aprovisionamiento", permiso)
		}
	}
}

// BAC-32 (RF-04 / RF-07): el worker de UPID debe contemplar las tareas de creación
// (qmcreate/vzcreate) y, al completarse, autoasignar FULL_ACCESS al creador.
// No hay endpoint de creación todavía, así que se verifica el worker a nivel de fuente
// (mismo criterio que la verificación estática de BAC-23B / BAC-25B en esta suite).
func TestEtapa1Ola3BAC32(t *testing.T) {
	worker := string(readFile(t, sourcePath("backend", "internal/core/services/seguimiento_tareas.go")))

	// 1. La acción "create" (qmcreate / vzcreate) debe estar contemplada por el worker
	//    (traducción de nombres y/o manejo del resultado).
	if !strings.Contains(worker, `"create"`) {
		t.Errorf("BAC-32: seguimiento_tareas.go no contempla la acción \"create\" (qmcreate/vzcreate) del worker UPID")
	}

	// 2. Al completarse una creación debe asignarse el permiso FULL_ACCESS al creador:
	//    el worker tiene que integrar el repositorio de permisos_instancia.
	if !regexp.MustCompile(`(?i)permiso`).MatchString(worker) {
		t.Errorf("BAC-32: seguimiento_tareas.go no integra el repositorio de permisos para autoasignar FULL_ACCESS al creador")
	}

	// 3. El contrato de eventos debe poder reportar la acción de creación.
	if tipos := string(readFile(t, sourcePath("backend", "internal/core/ports/event_port.go"))); !strings.Contains(tipos, "create") {
		t.Errorf("BAC-32: el contrato de eventos no referencia la acción de creación")
	}
}
