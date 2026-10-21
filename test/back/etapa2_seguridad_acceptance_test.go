package back_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

// Pruebas de aceptación de los FIX de seguridad y hardening de backend de documentacion/actual.md:
// FIX-107, FIX-109, FIX-110, FIX-111A, FIX-112, FIX-113 (BAC-24C), FIX-114, FIX-115A, FIX-116,
// FIX-117A, FIX-118 y FIX-119.
//
// Cada prueba valida el criterio de éxito de la tarea contra el código fuente (verificaciones
// estáticas, que corren también con -short) y contra el stack real de compose (verificaciones
// de comportamiento, que requieren Docker).

// secAdmin devuelve un token de administrador firmado para el stack de pruebas.
func secAdmin(t *testing.T) (adminToken, orgID string) {
	t.Helper()
	orgID = queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	return signedAccessToken(t, adminID, "ADMIN", orgID), orgID
}

// secReadGoDir concatena el código de todos los .go de un directorio del backend (incluye los _test.go).
func secReadGoDir(t *testing.T, relative string) string {
	t.Helper()
	dir := sourcePath("backend", relative)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("no se pudo leer el directorio %s: %v", dir, err)
	}
	var sb strings.Builder
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			sb.Write(readFile(t, filepath.Join(dir, entry.Name())))
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// secFunctionBody devuelve el texto de una función Go desde su firma hasta la siguiente función de nivel superior.
func secFunctionBody(source, name string) string {
	start := regexp.MustCompile(`(?m)^func (\([^)]*\) )?` + regexp.QuoteMeta(name) + `\(`).FindStringIndex(source)
	if start == nil {
		return ""
	}
	rest := source[start[0]:]
	if next := regexp.MustCompile(`(?m)^func `).FindStringIndex(rest[1:]); next != nil {
		return rest[:next[0]+1]
	}
	return rest
}

// secTOTP genera un código TOTP vigente liberando antes el anti-replay (BAC-11) del usuario.
func secTOTP(t *testing.T, email, secret string) string {
	t.Helper()
	queryDatabase(t, fmt.Sprintf("UPDATE usuarios SET ultimo_totp_periodo = 0 WHERE email_usuario = '%s';", email))
	code, err := totp.GenerateCode(secret, time.Now().UTC())
	if err != nil {
		t.Fatalf("no se pudo generar TOTP: %v", err)
	}
	return code
}

// secWizardPayload es una creación válida del contrato wizard (RF-07); cada prueba hostil altera un campo.
func secWizardPayload(mutate func(p map[string]any)) map[string]any {
	payload := map[string]any{
		"type": "VM", "node": "pve", "name": "vm-valida-fix116",
		"cpu": map[string]any{"cores": 1}, "ram": map[string]any{"memoryMb": 512},
		"disk":    map[string]any{"storageId": "local", "sizeGb": 4},
		"network": map[string]any{"bridge": "vmbr0", "firewall": false, "ipv4": map[string]any{"type": "dhcp"}},
		"os":      map[string]any{"templateId": "local:iso/debian-12.iso"},
	}
	if mutate != nil {
		mutate(payload)
	}
	return payload
}

func secErrorCode(body map[string]any) string {
	code, _ := body["errorCode"].(string)
	return code
}

// FIX-107: los VMID viajan como tokens opacos y reversibles, no como enteros secuenciales.
func TestFIX107_OfuscacionDeVMIDs(t *testing.T) {
	cryptoCode := secReadGoDir(t, "internal/infrastructure/crypto")
	mainCode := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))
	handlersCode := secReadGoDir(t, "internal/adapters/primary/http")

	t.Run("servicio de ofuscación bidireccional con clave secreta y pruebas unitarias", func(t *testing.T) {
		for _, fn := range []string{"OfuscarVMID", "DesofuscarVMID"} {
			if !regexp.MustCompile(`func (\([^)]*\) )?` + fn + `\(`).MatchString(cryptoCode) {
				t.Errorf("FIX-107: internal/infrastructure/crypto no define %s", fn)
			}
			if !regexp.MustCompile(`(?s)func Test\w*.*` + fn).MatchString(cryptoCode) {
				t.Errorf("FIX-107: no hay pruebas unitarias en internal/infrastructure/crypto que ejerciten %s", fn)
			}
		}
		if !strings.Contains(cryptoCode, "INSTANCE_ID_SALT") && !strings.Contains(mainCode, "INSTANCE_ID_SALT") {
			t.Errorf("FIX-107: la clave de ofuscación debe provenir de la variable de entorno INSTANCE_ID_SALT")
		}
	})

	t.Run("las rutas de instancias aceptan :instanceRef", func(t *testing.T) {
		if !strings.Contains(mainCode, ":instanceRef") && !strings.Contains(handlersCode, ":instanceRef") {
			t.Errorf("FIX-107: las rutas /instances/:instanceRef no están declaradas en cmd/api/main.go")
		}
	})

	t.Run("GET /api/instances devuelve ids opacos y las acciones aceptan el token ofuscado", func(t *testing.T) {
		adminToken, _ := secAdmin(t)
		status, items := requestJSONArray(t, http.MethodGet, "/instances", adminToken)
		if status != http.StatusOK || len(items) == 0 {
			t.Fatalf("GET /api/instances esperado 200 con instancias, recibido %d: %#v", status, items)
		}
		refByName := map[string]string{}
		for _, raw := range items {
			item := raw.(map[string]any)
			ref, isString := item["id"].(string)
			if !isString || ref == "" {
				t.Fatalf("FIX-107: el id público debe ser un string opaco, no %T(%v)", item["id"], item["id"])
			}
			if regexp.MustCompile(`^\d+$`).MatchString(ref) {
				t.Errorf("FIX-107: el id %q es un VMID numérico secuencial sin ofuscar", ref)
			}
			refByName[fmt.Sprint(item["name"])] = ref
		}
		ref := refByName["web-01"]
		if ref == "" {
			t.Fatalf("la instancia web-01 (VMID 101) no figura en el inventario: %#v", refByName)
		}

		// El token ofuscado opera la máquina correcta sin exponer el VMID.
		if status, body := requestJSON(t, http.MethodPost, "/instances/"+ref+"/stop", adminToken, nil); status != http.StatusAccepted {
			t.Errorf("FIX-107: POST /instances/%s/stop esperado 202, recibido %d: %#v", ref, status, body)
		}

		// Sin oráculo: un token inválido y uno inexistente responden idéntico.
		invalidStatus, invalidBody := requestJSON(t, http.MethodPost, "/instances/inst_invalido0/start", adminToken, nil)
		missingStatus, missingBody := requestJSON(t, http.MethodPost, "/instances/9999/start", adminToken, nil)
		if invalidStatus != http.StatusNotFound || missingStatus != http.StatusNotFound {
			t.Errorf("FIX-107: tokens inválidos/inexistentes deben dar 404, recibido %d y %d", invalidStatus, missingStatus)
		}
		if secErrorCode(invalidBody) != secErrorCode(missingBody) {
			t.Errorf("FIX-107: el error distingue un token inválido (%q) de un inexistente (%q): oráculo de descifrado", secErrorCode(invalidBody), secErrorCode(missingBody))
		}
	})
}

// FIX-109: el frontend nunca recibe hostnames reales del hipervisor ni nombres de storage.
func TestFIX109_EnmascaramientoDeTopologia(t *testing.T) {
	requireIntegration(t)
	adminToken, _ := secAdmin(t)

	status, items := requestJSONArray(t, http.MethodGet, "/instances", adminToken)
	if status != http.StatusOK {
		t.Fatalf("GET /api/instances esperado 200, recibido %d", status)
	}
	for _, raw := range items {
		item := raw.(map[string]any)
		if node := fmt.Sprint(item["node"]); node == "pve" {
			t.Errorf("FIX-109: la instancia %v expone el nombre real del nodo %q", item["name"], node)
		}
	}

	for _, path := range []string{"/instances", "/node/status", "/node/storages", "/node/templates"} {
		status, body := requestValue(t, http.MethodGet, path, adminToken, nil)
		if status != http.StatusOK {
			continue
		}
		encoded, _ := json.Marshal(body)
		for _, leak := range []string{`"pve"`, `"local-lvm"`, `"local-zfs"`} {
			if strings.Contains(string(encoded), leak) {
				t.Errorf("FIX-109: GET %s expone el identificador físico %s: %s", path, leak, encoded)
			}
		}
	}
}

// FIX-110: ningún error ni evento filtra rutas del sistema, trazas Perl ni el UPID de Proxmox.
func TestFIX110_SaneamientoDeErrores(t *testing.T) {
	handlersDir := "internal/adapters/primary/http"
	t.Run("los handlers de infraestructura no devuelven err.Error() al cliente en errores 5xx", func(t *testing.T) {
		// Los 400 de validación propia (p. ej. validarCreacionNormalizada) sí pueden describir el campo inválido;
		// lo que no puede salir es el texto de un error del hipervisor o del sistema.
		leak := regexp.MustCompile(`(SendError\(c, http\.Status(InternalServerError|BadGateway|GatewayTimeout|ServiceUnavailable)|c\.JSON\(http\.Status(InternalServerError|BadGateway|GatewayTimeout|ServiceUnavailable))[^\n]*\berr(or)?\.Error\(\)`)
		for _, file := range []string{"instance_handler.go", "instance_config_handler.go", "node_handler.go", "ws_metrics_handler.go", "events_handler.go"} {
			source := string(readFile(t, sourcePath("backend", handlersDir+"/"+file)))
			if loc := leak.FindString(source); loc != "" {
				t.Errorf("FIX-110: %s envía el error crudo al cliente: %s", file, strings.TrimSpace(loc))
			}
		}
	})

	t.Run("mapearErrorProxmox solo emite mensajes canónicos y registra el error crudo en el log", func(t *testing.T) {
		body := secFunctionBody(string(readFile(t, sourcePath("backend", handlersDir+"/instance_handler.go"))), "mapearErrorProxmox")
		if body == "" {
			t.Fatalf("FIX-110: no se encontró mapearErrorProxmox")
		}
		if regexp.MustCompile(`SendError\([^\n]*err\b`).MatchString(body) {
			t.Errorf("FIX-110: mapearErrorProxmox interpola el error de Proxmox en la respuesta")
		}
		if !strings.Contains(body, "log.Printf") {
			t.Errorf("FIX-110: mapearErrorProxmox debe registrar el error crudo solo en el log del servidor")
		}
	})

	t.Run("TASK_FINISHED y las métricas en vivo nunca incluyen el upid", func(t *testing.T) {
		tracking := string(readFile(t, sourcePath("backend", "internal/core/services/seguimiento_tareas.go")))
		if strings.Contains(tracking, `"upid"`) {
			t.Errorf("FIX-110: seguimiento_tareas.go publica la clave \"upid\" en los detalles del evento")
		}
		ws := string(readFile(t, sourcePath("backend", handlersDir+"/ws_metrics_handler.go")))
		if strings.Contains(strings.ToLower(ws), `json:"upid`) {
			t.Errorf("FIX-110: el canal WebSocket de métricas expone el campo upid")
		}
	})

	t.Run("ninguna respuesta de error filtra rutas de Proxmox, módulos Perl ni UPID", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)
		patterns := regexp.MustCompile(`(?i)/usr/share|\.pm line|perl|UPID:|unable to create|pvedaemon|pveproxy`)
		for _, probe := range []struct{ method, path string }{
			{http.MethodPost, "/instances/9999/start"},
			{http.MethodDelete, "/instances/9999"},
			{http.MethodGet, "/instances/9999/metrics"},
			{http.MethodPost, "/instances/../../etc/passwd/start"},
			{http.MethodPut, "/instances/9999/config"},
		} {
			raw := requestRaw(t, probe.method, probe.path, adminToken, nil)
			if loc := patterns.FindString(raw.RawBody); loc != "" {
				t.Errorf("FIX-110: %s %s filtra %q: %s", probe.method, probe.path, loc, raw.RawBody)
			}
		}
	})
}

// FIX-111A: el cambio de email exige reautenticación y confirmación de doble canal.
func TestFIX111A_CambioDeEmailSeguro(t *testing.T) {
	requireIntegration(t)
	adminToken, _ := secAdmin(t)

	user := createActiveUser(t, adminToken, "fix111a", "OPERATOR")
	session := loginWithTOTP(t, user.Email, user.Password, user.Secret)
	newEmail := fmt.Sprintf("nuevo.fix111a.%d@elcentinela.com", time.Now().UnixNano())
	currentEmail := func() string {
		return queryDatabase(t, fmt.Sprintf("SELECT email_usuario FROM usuarios WHERE id = '%s';", user.ID))
	}

	t.Run("sin contraseña actual ni código 2FA se rechaza y el email no cambia", func(t *testing.T) {
		status, body := requestJSON(t, http.MethodPut, "/account/profile", session.AccessToken, map[string]any{"emailUsuario": newEmail})
		if status != http.StatusBadRequest && status != http.StatusUnauthorized {
			t.Errorf("FIX-111A: cambio de email sin credenciales esperado 400/401, recibido %d: %#v", status, body)
		}
		if got := currentEmail(); got != user.Email {
			t.Errorf("FIX-111A: el email mutó sin reautenticación: %s", got)
		}
	})

	t.Run("con contraseña o TOTP incorrectos se rechaza", func(t *testing.T) {
		status, _ := requestJSON(t, http.MethodPut, "/account/profile", session.AccessToken, map[string]any{
			"emailUsuario": newEmail, "contrasenaActual": "Incorrecta1!", "codigo2fa": "000000",
		})
		if status != http.StatusBadRequest && status != http.StatusUnauthorized {
			t.Errorf("FIX-111A: credenciales inválidas esperado 400/401, recibido %d", status)
		}
		if got := currentEmail(); got != user.Email {
			t.Errorf("FIX-111A: el email mutó con credenciales inválidas: %s", got)
		}
	})

	t.Run("con reautenticación válida el cambio queda pendiente de confirmación y auditado", func(t *testing.T) {
		status, body := requestJSON(t, http.MethodPut, "/account/profile", session.AccessToken, map[string]any{
			"emailUsuario": newEmail, "contrasenaActual": user.Password, "codigo2fa": secTOTP(t, user.Email, user.Secret),
		})
		if status != http.StatusOK && status != http.StatusAccepted {
			t.Fatalf("FIX-111A: cambio con reautenticación esperado 200/202, recibido %d: %#v", status, body)
		}
		if got := currentEmail(); got != user.Email {
			t.Errorf("FIX-111A: el email debe mutar solo al confirmar por el nuevo canal, pero ya vale %s", got)
		}
		if n := auditCount(t, user.ID, "UPDATE_PROFILE_EMAIL_REQUESTED"); n < 1 {
			t.Errorf("FIX-111A: falta el registro UPDATE_PROFILE_EMAIL_REQUESTED en auditoria")
		}
	})

	t.Run("editar solo el nombre sigue funcionando sin reautenticación", func(t *testing.T) {
		status, body := requestJSON(t, http.MethodPut, "/account/profile", session.AccessToken, map[string]any{"nombreCompleto": "Nombre Actualizado"})
		if status != http.StatusOK {
			t.Errorf("FIX-111A: cambiar solo el nombre esperado 200, recibido %d: %#v", status, body)
		}
	})
}

// FIX-112: rate limiting y protección contra fuerza bruta en login, 2FA y recuperación.
func TestFIX112_RateLimitingAutenticacion(t *testing.T) {
	t.Run("ConfirmarRecuperacionContrasena bloquea la fila con FOR UPDATE", func(t *testing.T) {
		source := secReadGoDir(t, "internal/core/services")
		body := secFunctionBody(source, "ConfirmarRecuperacionContrasena")
		if body == "" {
			t.Fatalf("FIX-112: no se encontró ConfirmarRecuperacionContrasena")
		}
		if !regexp.MustCompile(`FOR UPDATE|clause\.Locking|ForUpdate`).MatchString(body) {
			t.Errorf("FIX-112: IntentosRecuperacion se lee sin bloqueo pesimista (SELECT FOR UPDATE): condición de carrera")
		}
	})

	t.Run("login: el sexto intento fallido del mismo par IP/email responde 429 con Retry-After", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)
		user := createActiveUser(t, adminToken, "fix112login", "OPERATOR")
		for i := 1; i <= 5; i++ {
			raw := requestRaw(t, http.MethodPost, "/auth/login", "", map[string]any{"email": user.Email, "password": "Equivocada1!"})
			if raw.Status != http.StatusUnauthorized {
				t.Fatalf("FIX-112: intento fallido %d esperado 401, recibido %d: %s", i, raw.Status, raw.RawBody)
			}
		}
		raw := requestRaw(t, http.MethodPost, "/auth/login", "", map[string]any{"email": user.Email, "password": user.Password})
		if raw.Status != http.StatusTooManyRequests {
			t.Fatalf("FIX-112: el sexto intento esperado 429, recibido %d: %s", raw.Status, raw.RawBody)
		}
		if got := raw.Header.Get("Retry-After"); got != "900" {
			t.Errorf("FIX-112: Retry-After esperado 900, recibido %q", got)
		}
	})

	t.Run("2FA: el cuarto intento revoca el token pre-auth con 401 TOKEN_REVOKED", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)
		user := createActiveUser(t, adminToken, "fix112tfa", "OPERATOR")
		status, body := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{"email": user.Email, "password": user.Password})
		if status != http.StatusOK {
			t.Fatalf("login esperado 200, recibido %d: %#v", status, body)
		}
		preAuth := requiredString(t, body, "jwtTemporal")
		for i := 1; i <= 3; i++ {
			raw := requestRaw(t, http.MethodPost, "/auth/2fa/verify", preAuth, map[string]any{"codigo": "000000"})
			if raw.Status == http.StatusOK || raw.Status == http.StatusTooManyRequests {
				t.Fatalf("FIX-112: código incorrecto %d no debe devolver %d", i, raw.Status)
			}
		}
		// Aunque el código sea correcto, el token temporal ya fue destruido.
		raw := requestRaw(t, http.MethodPost, "/auth/2fa/verify", preAuth, map[string]any{"codigo": secTOTP(t, user.Email, user.Secret)})
		if raw.Status != http.StatusUnauthorized || secErrorCode(raw.Body) != "TOKEN_REVOKED" {
			t.Errorf("FIX-112: el cuarto intento esperado 401 TOKEN_REVOKED, recibido %d %s", raw.Status, raw.RawBody)
		}
	})

	t.Run("recuperación: máximo 3 solicitudes por correo cada 60 minutos", func(t *testing.T) {
		requireIntegration(t)
		email := fmt.Sprintf("fix112.forgot.%d@elcentinela.com", time.Now().UnixNano())
		for i := 1; i <= 3; i++ {
			if status, _ := requestJSON(t, http.MethodPost, "/auth/password/forgot", "", map[string]any{"email": email}); status != http.StatusOK {
				t.Fatalf("FIX-112: solicitud %d esperado 200, recibido %d", i, status)
			}
		}
		raw := requestRaw(t, http.MethodPost, "/auth/password/forgot", "", map[string]any{"email": email})
		if raw.Status != http.StatusTooManyRequests {
			t.Errorf("FIX-112: la cuarta solicitud esperado 429, recibido %d: %s", raw.Status, raw.RawBody)
		}
	})
}

// FIX-113 / BAC-24C: al completarse un borrado se purgan los permisos del VMID (y solo entonces).
func TestFIX113_PurgaDePermisosAlEliminarInstancia(t *testing.T) {
	t.Run("seguimiento_tareas.go purga permisos y audita INSTANCE_PERMISSIONS_PURGED", func(t *testing.T) {
		source := secReadGoDir(t, "internal/core/services")
		for _, marker := range []string{"permisos_instancia", "INSTANCE_PERMISSIONS_PURGED"} {
			if !strings.Contains(source, marker) {
				t.Errorf("FIX-113: internal/core/services no contiene %q", marker)
			}
		}
	})

	t.Run("un borrado fallido conserva los permisos y uno exitoso los elimina", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)
		operator := createActiveUser(t, adminToken, "fix113", "OPERATOR")
		operatorToken := loginWithTOTP(t, operator.Email, operator.Password, operator.Secret).AccessToken
		permissions := "/admin/users/" + operator.ID + "/permissions"
		if status, body := requestValue(t, http.MethodPut, permissions, adminToken, map[string]any{
			"permisos": []map[string]any{{"vmid": 102, "nivelAcceso": "READ_ONLY"}},
		}); status != http.StatusNoContent {
			t.Fatalf("precondición: asignar permiso sobre 102 falló (%d): %#v", status, body)
		}
		count := func() string {
			return queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM permisos_instancia WHERE usuario_id = '%s' AND vmid_proxmox = 102;", operator.ID))
		}

		// 1. Tarea fallida en Proxmox: no se toca nada.
		if _, err := runCompose("exec", "-T", "proxmox", "touch", "/tmp/tareas-fallidas"); err != nil {
			t.Fatalf("no se pudo simular el fallo en el stub de Proxmox: %v", err)
		}
		t.Cleanup(func() { _, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-fallidas") })
		status, body := requestJSON(t, http.MethodDelete, "/instances/102", adminToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("DELETE /instances/102 esperado 202, recibido %d: %#v", status, body)
		}
		failedID := requiredString(t, body, "tareaId")
		if !waitForValue(t, fmt.Sprintf("SELECT estado FROM tareas_asincronas WHERE id = '%s';", failedID), "FAILED", 15*time.Second) {
			t.Fatalf("la tarea de borrado fallida no llegó a FAILED")
		}
		if got := count(); got != "1" {
			t.Errorf("FIX-113: un borrado fallido no debe tocar permisos_instancia, quedan %s filas", got)
		}

		// 2. Tarea exitosa: se purgan los permisos y queda auditado.
		_, _ = runCompose("exec", "-T", "proxmox", "rm", "-f", "/tmp/tareas-fallidas")
		status, body = requestJSON(t, http.MethodDelete, "/instances/102", adminToken, nil)
		if status != http.StatusAccepted {
			t.Fatalf("DELETE /instances/102 esperado 202, recibido %d: %#v", status, body)
		}
		if !waitForValue(t, fmt.Sprintf("SELECT count(*) FROM permisos_instancia WHERE usuario_id = '%s' AND vmid_proxmox = 102;", operator.ID), "0", 15*time.Second) {
			t.Fatalf("FIX-113: los permisos de la 102 no se purgaron tras el borrado exitoso (quedan %s)", count())
		}
		status, value := requestValue(t, http.MethodGet, permissions, adminToken, nil)
		encoded, _ := json.Marshal(value)
		if status != http.StatusOK || regexp.MustCompile(`"vmid":\s*102\b`).Match(encoded) {
			t.Errorf("FIX-113: GET %s sigue listando la 102 (%d): %s", permissions, status, encoded)
		}
		if n := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM auditoria WHERE accion = 'INSTANCE_PERMISSIONS_PURGED' AND detalles::text LIKE '%%%s%%';", operator.ID)); n == "0" {
			t.Errorf("FIX-113: falta INSTANCE_PERMISSIONS_PURGED con el usuario afectado %s en auditoria", operator.ID)
		}

		// 3. Un VMID reciclado no hereda acceso: el operador no puede operar la instancia.
		if status, _ := requestJSON(t, http.MethodGet, "/instances/102/metrics", operatorToken, nil); status != http.StatusForbidden {
			t.Errorf("FIX-113: el operador debe recibir 403 sobre el VMID reciclado, recibido %d", status)
		}
	})
}

// FIX-114: el cambio voluntario de contraseña revoca las demás sesiones activas.
func TestFIX114_RevocacionDeSesionesAlCambiarContrasena(t *testing.T) {
	requireIntegration(t)
	adminToken, _ := secAdmin(t)
	user := createActiveUser(t, adminToken, "fix114", "OPERATOR")

	sessionA := loginWithTOTP(t, user.Email, user.Password, user.Secret)
	sessionB := loginWithTOTP(t, user.Email, user.Password, user.Secret)
	if status, _ := requestJSON(t, http.MethodGet, "/account/profile", sessionB.AccessToken, nil); status != http.StatusOK {
		t.Fatalf("precondición: la sesión B debe estar operativa, recibió %d", status)
	}

	status, body := requestJSON(t, http.MethodPut, "/account/password", sessionA.AccessToken, map[string]any{
		"contrasenaActual": user.Password, "contrasenaNueva": "Cambiada123!",
	})
	if status != http.StatusOK {
		t.Fatalf("PUT /account/password esperado 200, recibido %d: %#v", status, body)
	}

	raw := requestRaw(t, http.MethodGet, "/account/profile", sessionB.AccessToken, nil)
	if raw.Status != http.StatusUnauthorized || secErrorCode(raw.Body) != "TOKEN_REVOKED" {
		t.Errorf("FIX-114: la otra sesión debe recibir 401 TOKEN_REVOKED, recibió %d %s", raw.Status, raw.RawBody)
	}
	if n := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM sesiones_activas WHERE usuario_id = '%s' AND activa = true;", user.ID)); n == "2" {
		t.Errorf("FIX-114: las dos sesiones siguen activas en sesiones_activas")
	}
	if n := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM auditoria WHERE usuario_id = '%s' AND detalles::text LIKE '%%sesiones_revocadas%%';", user.ID)); n == "0" {
		t.Errorf("FIX-114: la auditoría no registra sesiones_revocadas: true")
	}
}

// FIX-115A: las contraseñas aceptan hasta 128 caracteres (NIST SP 800-63B).
func TestFIX115A_ContrasenasDeHasta128Caracteres(t *testing.T) {
	t.Run("el código ya no limita a 12 caracteres", func(t *testing.T) {
		passwordCode := string(readFile(t, sourcePath("backend", "internal/infrastructure/crypto/password.go")))
		if regexp.MustCompile(`len\(contrasena\)\s*>\s*12\b`).MatchString(passwordCode) {
			t.Errorf("FIX-115A: password.go sigue rechazando contraseñas de más de 12 caracteres")
		}
		if !regexp.MustCompile(`len\(contrasena\)\s*>\s*128\b`).MatchString(passwordCode) {
			t.Errorf("FIX-115A: password.go no impone la cota de 128 caracteres")
		}
		for _, file := range []string{
			"internal/adapters/primary/http/auth_handler.go", "internal/adapters/primary/http/user_handler.go",
			"internal/adapters/primary/http/account_handler.go", "internal/core/ports/user_port.go",
		} {
			if strings.Contains(string(readFile(t, sourcePath("backend", file))), "max=12") {
				t.Errorf("FIX-115A: %s conserva el tag de validación max=12", file)
			}
		}
	})

	t.Run("PUT /account/password acepta 20, 64 y 128 caracteres y rechaza 129", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)
		user := createActiveUser(t, adminToken, "fix115a", "OPERATOR")
		current := user.Password
		change := func(newPassword string) int {
			token := loginWithTOTP(t, user.Email, current, user.Secret).AccessToken
			status, _ := requestJSON(t, http.MethodPut, "/account/password", token, map[string]any{
				"contrasenaActual": current, "contrasenaNueva": newPassword,
			})
			if status == http.StatusOK {
				current = newPassword
			}
			return status
		}
		long := func(n int) string { return "Aa1!" + strings.Repeat("x", n-4) }

		if status := change(long(129)); status != http.StatusBadRequest {
			t.Errorf("FIX-115A: 129 caracteres esperado 400, recibido %d", status)
		}
		for _, n := range []int{20, 64, 128} {
			if status := change(long(n)); status != http.StatusOK {
				t.Errorf("FIX-115A: %d caracteres esperado 200, recibido %d", n, status)
			}
		}
	})
}

// FIX-116: sanitización de parámetros que Proxmox interpreta como listas separadas por comas.
func TestFIX116_InyeccionDeParametrosProxmox(t *testing.T) {
	t.Run("client.go valida que no haya comas en los parámetros concatenados", func(t *testing.T) {
		client := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))
		if !regexp.MustCompile(`strings\.Contains(Any)?\([^)]*","`).MatchString(client) && !regexp.MustCompile(`ContainsRune\([^)]*','`).MatchString(client) {
			t.Errorf("FIX-116: client.go no verifica que Bridge/StorageID/TemplateID carezcan de comas antes de armar net0")
		}
	})

	t.Run("POST /api/instances rechaza con 400 INVALID_REQUEST valores hostiles antes de llegar a Proxmox", func(t *testing.T) {
		requireIntegration(t)
		adminToken, _ := secAdmin(t)
		hostile := map[string]func(p map[string]any){
			"bridge con directivas extra": func(p map[string]any) { p["network"].(map[string]any)["bridge"] = "vmbr0,rate=100,tag=20" },
			"bridge con espacios":         func(p map[string]any) { p["network"].(map[string]any)["bridge"] = "vmbr0 firewall=0" },
			"bridge con salto de línea":   func(p map[string]any) { p["network"].(map[string]any)["bridge"] = "vmbr0\nrate=1" },
			"hostname con guion bajo":     func(p map[string]any) { p["name"] = "mi_servidor" },
			"hostname con coma":           func(p map[string]any) { p["name"] = "srv,cores=64" },
			"hostname que empieza con -":  func(p map[string]any) { p["name"] = "-srv" },
			"hostname de más de 63":       func(p map[string]any) { p["name"] = strings.Repeat("a", 64) },
			"storageId con coma":          func(p map[string]any) { p["disk"].(map[string]any)["storageId"] = "local,size=999" },
			"templateId con coma":         func(p map[string]any) { p["os"].(map[string]any)["templateId"] = "local:iso/x.iso,media=disk" },
			"templateId con control":      func(p map[string]any) { p["os"].(map[string]any)["templateId"] = "local:iso/x.iso\x00" },
		}
		for name, mutate := range hostile {
			raw := requestRaw(t, http.MethodPost, "/instances", adminToken, secWizardPayload(mutate))
			if raw.Status != http.StatusBadRequest || secErrorCode(raw.Body) != "INVALID_REQUEST" {
				t.Errorf("FIX-116: %s esperado 400 INVALID_REQUEST, recibido %d %s", name, raw.Status, raw.RawBody)
			}
		}
		// Control: una creación bien formada no se rechaza por validación.
		if raw := requestRaw(t, http.MethodPost, "/instances", adminToken, secWizardPayload(nil)); raw.Status == http.StatusBadRequest {
			t.Errorf("FIX-116: una creación válida no debe dar 400: %s", raw.RawBody)
		}
	})
}

// FIX-117A: un OPERATOR no ve la telemetría global del host.
func TestFIX117A_SegregacionDeTelemetriaGlobal(t *testing.T) {
	requireIntegration(t)
	adminToken, _ := secAdmin(t)
	operator := createActiveUser(t, adminToken, "fix117a", "OPERATOR")
	operatorToken := loginWithTOTP(t, operator.Email, operator.Password, operator.Secret).AccessToken
	if status, body := requestValue(t, http.MethodPut, "/admin/users/"+operator.ID+"/permissions", adminToken, map[string]any{
		"permisos": []map[string]any{{"vmid": 101, "nivelAcceso": "FULL_ACCESS"}},
	}); status != http.StatusNoContent {
		t.Fatalf("precondición: asignar permisos falló (%d): %#v", status, body)
	}

	t.Run("ADMIN conserva la telemetría completa del clúster", func(t *testing.T) {
		status, body := requestJSON(t, http.MethodGet, "/node/status", adminToken, nil)
		if status != http.StatusOK {
			t.Fatalf("GET /node/status esperado 200, recibido %d", status)
		}
		if body["instancesSummary"] == nil || body["storage"] == nil {
			t.Errorf("FIX-117A: el ADMIN debe recibir instancesSummary y storage: %#v", body)
		}
	})

	t.Run("OPERATOR no recibe el inventario global ni la capacidad de los discos", func(t *testing.T) {
		status, body := requestJSON(t, http.MethodGet, "/node/status", operatorToken, nil)
		if status != http.StatusOK {
			t.Fatalf("GET /node/status esperado 200 para OPERATOR, recibido %d: %#v", status, body)
		}
		if body["storage"] != nil {
			t.Errorf("FIX-117A: el OPERATOR recibe la capacidad de almacenamiento del host: %#v", body["storage"])
		}
		if summary, ok := body["instancesSummary"].(map[string]any); ok {
			total := 0.0
			for _, kind := range []string{"vms", "lxc"} {
				if group, ok := summary[kind].(map[string]any); ok {
					if n, ok := group["total"].(float64); ok {
						total += n
					}
				}
			}
			if total > 1 {
				t.Errorf("FIX-117A: el conteo del OPERATOR (%v) excede sus instancias asignadas (1)", total)
			}
		}
	})

	t.Run("OPERATOR no recibe discos globales en el histórico del nodo", func(t *testing.T) {
		raw := requestRaw(t, http.MethodGet, "/node/history", operatorToken, nil)
		if raw.Status == http.StatusOK && strings.Contains(strings.ToLower(raw.RawBody), "storage") {
			t.Errorf("FIX-117A: /node/history expone storage al OPERATOR: %s", raw.RawBody)
		}
	})
}

// FIX-118: no existen credenciales por defecto hardcodeadas en producción.
func TestFIX118_SinCredencialesPorDefecto(t *testing.T) {
	source := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/postgres/db.go")))
	body := secFunctionBody(source, "seedAdminPorDefecto") + "\n" + secFunctionBody(source, "credencialesIniciales")
	if strings.TrimSpace(body) == "" {
		t.Skip("seedAdminPorDefecto ya no existe: el seed por defecto fue eliminado")
	}
	if strings.Contains(source, "Admin123!") {
		t.Errorf("FIX-118: seedAdminPorDefecto conserva la contraseña estática Admin123!")
	}
	if !strings.Contains(body, "APP_ENV") || !strings.Contains(body, "production") {
		t.Errorf("FIX-118: seedAdminPorDefecto no se desactiva cuando APP_ENV=production")
	}
	if !strings.Contains(body, "GenerarContrasenaTemp") {
		t.Errorf("FIX-118: fuera de producción la contraseña debe generarse con crypto.GenerarContrasenaTemp()")
	}
	if !strings.Contains(source, "INITIAL_ADMIN_EMAIL") || !strings.Contains(source, "INITIAL_ADMIN_PASSWORD") {
		t.Errorf("FIX-118: el administrador inicial de producción debe parametrizarse con INITIAL_ADMIN_EMAIL / INITIAL_ADMIN_PASSWORD")
	}
}

// FIX-119: AsignarPermisos valida existencia de VMIDs y bloquea las instancias protegidas.
func TestFIX119_ValidacionDeVMIDsEnAsignacionDePermisos(t *testing.T) {
	requireIntegration(t)
	adminToken, _ := secAdmin(t)
	operator := createActiveUser(t, adminToken, "fix119", "OPERATOR")
	endpoint := "/admin/users/" + operator.ID + "/permissions"
	assigned := func() string {
		return queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM permisos_instancia WHERE usuario_id = '%s';", operator.ID))
	}

	t.Run("un VMID inexistente en Proxmox responde 400 INVALID_REQUEST", func(t *testing.T) {
		raw := requestRaw(t, http.MethodPut, endpoint, adminToken, map[string]any{
			"permisos": []map[string]any{{"vmid": 99999, "nivelAcceso": "READ_ONLY"}},
		})
		if raw.Status != http.StatusBadRequest || secErrorCode(raw.Body) != "INVALID_REQUEST" {
			t.Errorf("FIX-119: VMID inexistente esperado 400 INVALID_REQUEST, recibido %d %s", raw.Status, raw.RawBody)
		}
		if !strings.Contains(raw.RawBody, "99999") {
			t.Errorf("FIX-119: el error debe indicar el VMID inexistente: %s", raw.RawBody)
		}
	})

	t.Run("un VMID protegido responde 403 PROTECTED_INSTANCE_CANNOT_BE_ASSIGNED y no persiste nada", func(t *testing.T) {
		raw := requestRaw(t, http.MethodPut, endpoint, adminToken, map[string]any{
			"permisos": []map[string]any{{"vmid": 101, "nivelAcceso": "FULL_ACCESS"}, {"vmid": 103, "nivelAcceso": "READ_ONLY"}},
		})
		if raw.Status != http.StatusForbidden || secErrorCode(raw.Body) != "PROTECTED_INSTANCE_CANNOT_BE_ASSIGNED" {
			t.Errorf("FIX-119: VMID protegido esperado 403 PROTECTED_INSTANCE_CANNOT_BE_ASSIGNED, recibido %d %s", raw.Status, raw.RawBody)
		}
		if got := assigned(); got != "0" {
			t.Errorf("FIX-119: la asignación rechazada persistió %s filas en permisos_instancia", got)
		}
	})

	t.Run("VMIDs existentes y operables se asignan con 204", func(t *testing.T) {
		status, body := requestValue(t, http.MethodPut, endpoint, adminToken, map[string]any{
			"permisos": []map[string]any{{"vmid": 101, "nivelAcceso": "FULL_ACCESS"}, {"vmid": 102, "nivelAcceso": "READ_ONLY"}},
		})
		if status != http.StatusNoContent {
			t.Errorf("FIX-119: asignación válida esperado 204, recibido %d: %#v", status, body)
		}
		if got := assigned(); got != "2" {
			t.Errorf("FIX-119: se esperaban 2 filas en permisos_instancia, hay %s", got)
		}
	})
}
