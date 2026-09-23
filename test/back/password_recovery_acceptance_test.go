package back_test

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

// Regresión de tareas de documentacion/terminado.md:
// recuperación de contraseña (BAC-19, BAC-20), entrega de credenciales (BAC-16),
// contrato de eventos (BAC-21), auditoría append-only (BAC-18) y resets administrativos
// (BAC-13, BAC-15). Las credenciales y códigos se leen del MockEmailService, el mismo
// canal por el que los recibe el usuario.
func TestHitoRecuperacionDeContrasenasYNotificaciones(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	operator := createActiveUser(t, adminToken, "op_recov", "OPERATOR")
	operatorToken := loginWithTOTP(t, operator.Email, operator.Password, operator.Secret).AccessToken

	requestRecovery := func(t *testing.T, email string) (int, map[string]any) {
		t.Helper()
		return requestJSON(t, http.MethodPost, "/auth/password/forgot", "", map[string]any{"email": email})
	}
	confirmRecovery := func(t *testing.T, email, code, password string) (int, map[string]any) {
		t.Helper()
		return requestJSON(t, http.MethodPost, "/auth/password/reset", "", map[string]any{
			"email": email, "codigo": code, "nuevaContrasena": password,
		})
	}
	loginStatus := func(t *testing.T, email, password string) int {
		t.Helper()
		status, _ := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{"email": email, "password": password})
		return status
	}

	t.Run("BAC-19 solicitud de recuperacion valida formato responde generico e invalida el codigo previo", func(t *testing.T) {
		user := createActiveUser(t, adminToken, "recupera19", "OPERATOR")

		if status, _ := requestRecovery(t, "correo-invalido-sin-arroba"); status != http.StatusBadRequest {
			t.Errorf("formato inválido: esperado 400, recibido %d", status)
		}

		existingStatus, existingBody := requestRecovery(t, user.Email)
		missingStatus, missingBody := requestRecovery(t, "noexiste.bac19@elcentinela.com")
		if existingStatus != http.StatusOK || missingStatus != existingStatus || existingBody["message"] != missingBody["message"] {
			t.Errorf("anti-enumeración: cuenta existente %d %#v vs inexistente %d %#v deben ser idénticas", existingStatus, existingBody, missingStatus, missingBody)
		}
		firstCode := mailedSecret(t, user.Email, mailRecoveryCode)
		if !regexp.MustCompile(`^\d{6}$`).MatchString(firstCode) {
			t.Fatalf("el código enviado por correo debe tener 6 dígitos, se recibió %q", firstCode)
		}

		requestRecovery(t, user.Email)
		secondCode := mailedSecret(t, user.Email, mailRecoveryCode)
		if secondCode == firstCode {
			t.Skip("los dos códigos aleatorios coincidieron (probabilidad 1e-6); no se puede verificar la invalidación")
		}
		if status, _ := confirmRecovery(t, user.Email, firstCode, "Recupera1!"); status != http.StatusBadRequest {
			t.Errorf("el código anterior debe quedar invalidado al pedir uno nuevo: esperado 400, recibido %d", status)
		}
		if status, body := confirmRecovery(t, user.Email, secondCode, "Recupera1!"); status != http.StatusOK {
			t.Errorf("el código vigente debe aceptarse: esperado 200, recibido %d: %#v", status, body)
		}
	})

	t.Run("BAC-20 confirmacion valida codigo cambia la contrasena y revoca sesiones", func(t *testing.T) {
		user := createActiveUser(t, adminToken, "recupera20", "OPERATOR")
		previousSession := loginWithTOTP(t, user.Email, user.Password, user.Secret)

		requestRecovery(t, user.Email)
		code := mailedSecret(t, user.Email, mailRecoveryCode)

		wrong := "000000"
		if code == wrong {
			wrong = "111111"
		}
		if status, body := confirmRecovery(t, user.Email, wrong, "Recupera2!"); status != http.StatusBadRequest || body["errorCode"] == nil {
			t.Errorf("código incorrecto: esperado 400 con errorCode, recibido %d: %#v", status, body)
		}
		if status, body := confirmRecovery(t, user.Email, code, "Recupera2!"); status != http.StatusOK {
			t.Fatalf("código válido: esperado 200, recibido %d: %#v", status, body)
		}
		if status, _ := confirmRecovery(t, user.Email, code, "Recupera3!"); status != http.StatusBadRequest {
			t.Errorf("código ya usado: esperado 400, recibido %d", status)
		}
		if status := loginStatus(t, user.Email, user.Password); status != http.StatusUnauthorized {
			t.Errorf("la contraseña anterior debe dejar de funcionar: esperado 401, recibido %d", status)
		}
		if status := loginStatus(t, user.Email, "Recupera2!"); status != http.StatusOK {
			t.Errorf("la contraseña nueva debe funcionar: esperado 200, recibido %d", status)
		}
		if status, body := requestJSON(t, http.MethodGet, "/account/profile", previousSession.AccessToken, nil); status != http.StatusUnauthorized {
			t.Errorf("las sesiones previas deben revocarse tras la recuperación: esperado 401, recibido %d: %#v", status, body)
		}

		requestRecovery(t, user.Email)
		expiring := mailedSecret(t, user.Email, mailRecoveryCode)
		queryDatabase(t, fmt.Sprintf("UPDATE usuarios SET expiracion_codigo = NOW() - interval '1 minute' WHERE id = '%s';", user.ID))
		if status, _ := confirmRecovery(t, user.Email, expiring, "Recupera4!"); status != http.StatusBadRequest {
			t.Errorf("código vencido: esperado 400, recibido %d", status)
		}
	})

	t.Run("BAC-16 la clave temporal se entrega solo por correo y nunca en JSON ni en texto plano", func(t *testing.T) {
		email := fmt.Sprintf("mailer.%d@elcentinela.com", time.Now().UnixNano())
		status, created := requestJSON(t, http.MethodPost, "/admin/users", adminToken, map[string]any{
			"nombreCompleto": "Usuario Mailer", "nombreUsuario": fmt.Sprintf("mailer_%d", time.Now().UnixNano()),
			"emailUsuario": email, "rol": "OPERATOR",
		})
		if status != http.StatusCreated {
			t.Fatalf("POST /api/admin/users esperado 201, recibido %d: %#v", status, created)
		}
		for field := range created {
			if strings.Contains(strings.ToLower(field), "contrasena") || strings.Contains(strings.ToLower(field), "password") {
				t.Errorf("la respuesta de alta expone el campo %q", field)
			}
		}
		temporary := mailedSecret(t, email, mailTemporaryPassword)
		if status := loginStatus(t, email, temporary); status != http.StatusOK {
			t.Errorf("la clave temporal recibida por correo debe permitir el login: recibido %d", status)
		}
		userID := requiredString(t, created, "id")
		stored := queryDatabase(t, fmt.Sprintf("SELECT contrasena_hash FROM usuarios WHERE id = '%s';", userID))
		if !strings.HasPrefix(stored, "$2") || strings.Contains(stored, temporary) {
			t.Errorf("la clave temporal no quedó hasheada con bcrypt")
		}

		resetStatus, resetBody := requestJSON(t, http.MethodPost, "/admin/users/"+userID+"/password/reset", adminToken, nil)
		if resetStatus != http.StatusOK {
			t.Fatalf("reset administrativo esperado 200, recibido %d: %#v", resetStatus, resetBody)
		}
		newTemporary := mailedSecret(t, email, mailTemporaryPassword)
		if strings.Contains(fmt.Sprint(resetBody), newTemporary) {
			t.Errorf("la respuesta del reset administrativo expone la clave temporal")
		}
	})

	t.Run("BAC-21 contrato RealtimeEvent sincronizado entre Go y TypeScript", func(t *testing.T) {
		goSource, err := os.ReadFile("../../backend/internal/core/ports/event_port.go")
		if err != nil {
			t.Fatalf("no se encontró el contrato Go: %v", err)
		}
		tsSource, err := os.ReadFile("../../frontend/centinela/src/types/notifications.ts")
		if err != nil {
			t.Fatalf("no se encontró el contrato TypeScript: %v", err)
		}

		goStruct := regexp.MustCompile(`(?s)type RealtimeEvent struct \{(.*?)\n\}`).FindSubmatch(goSource)
		tsInterface := regexp.MustCompile(`(?s)export interface RealtimeEvent \{(.*?)\n\}`).FindSubmatch(tsSource)
		if goStruct == nil || tsInterface == nil {
			t.Fatalf("no se encontró struct RealtimeEvent en Go (%v) o interface RealtimeEvent en TS (%v)", goStruct != nil, tsInterface != nil)
		}
		var goFields, tsFields []string
		for _, m := range regexp.MustCompile(`json:"([^",]+)`).FindAllSubmatch(goStruct[1], -1) {
			goFields = append(goFields, string(m[1]))
		}
		for _, m := range regexp.MustCompile(`(?m)^\s*(\w+)\??:`).FindAllSubmatch(tsInterface[1], -1) {
			tsFields = append(tsFields, string(m[1]))
		}
		sort.Strings(goFields)
		sort.Strings(tsFields)
		if strings.Join(goFields, ",") != strings.Join(tsFields, ",") {
			t.Errorf("campos distintos:\n  Go: %v\n  TS: %v", goFields, tsFields)
		}

		for _, value := range []string{"INSTANCE_STATE_CHANGED", "INSTANCE_CREATED", "RESOURCE_SATURATION", "TASK_FINISHED", "INFO", "WARNING", "CRITICAL", "VM", "LXC", "NODE"} {
			inGo := strings.Contains(string(goSource), `"`+value+`"`)
			inTS := strings.Contains(string(tsSource), `'`+value+`'`)
			if inGo != inTS {
				t.Errorf("el valor %q está en Go=%v y en TS=%v", value, inGo, inTS)
			}
		}
	})

	t.Run("BAC-18 auditoria registra acciones administrativas y es append-only en la base", func(t *testing.T) {
		email := fmt.Sprintf("auditado.%d@elcentinela.com", time.Now().UnixNano())
		auditedID, _ := createUserThroughAPI(t, adminToken, fmt.Sprintf("auditado_%d", time.Now().UnixNano()), email, "OPERATOR")

		status, page := requestJSON(t, http.MethodGet, "/admin/audit?accion=CREAR_USUARIO&pagina=1&tamano=200", adminToken, nil)
		if status != http.StatusOK {
			t.Fatalf("GET /api/admin/audit esperado 200, recibido %d: %#v", status, page)
		}
		items, ok := page["items"].([]any)
		if !ok || page["total"] == nil {
			t.Fatalf("la respuesta debe tener { total, pagina, items }: %#v", page)
		}
		foundCreation := false
		for _, raw := range items {
			item := raw.(map[string]any)
			if item["accion"] != "CREAR_USUARIO" {
				t.Errorf("el filtro accion=CREAR_USUARIO devolvió %#v", item["accion"])
			}
			if strings.Contains(fmt.Sprint(item["detalles"]), auditedID) || strings.Contains(fmt.Sprint(item["detalles"]), email) {
				foundCreation = true
			}
		}
		if !foundCreation {
			t.Errorf("la alta del usuario %s no aparece en la auditoría CREAR_USUARIO", auditedID)
		}

		if opStatus, _ := requestValue(t, http.MethodGet, "/admin/audit", operatorToken, nil); opStatus != http.StatusForbidden {
			t.Errorf("OPERATOR consultando auditoría: esperado 403, recibido %d", opStatus)
		}

		export := requestRaw(t, http.MethodGet, "/admin/audit/export?formato=csv", adminToken, nil)
		if export.Status != http.StatusOK || !strings.Contains(export.Header.Get("Content-Type"), "text/csv") {
			t.Errorf("GET /api/admin/audit/export: esperado 200 text/csv, recibido %d %q", export.Status, export.Header.Get("Content-Type"))
		}
		if opExport := requestRaw(t, http.MethodGet, "/admin/audit/export", operatorToken, nil); opExport.Status != http.StatusForbidden {
			t.Errorf("OPERATOR exportando auditoría: esperado 403, recibido %d", opExport.Status)
		}

		// Criterio de éxito de BAC-18: UPDATE/DELETE con las credenciales de la aplicación
		// (las mismas de DB_DSN en compose.yaml) deben fallar a nivel de motor.
		if _, err := execDatabase(t, "UPDATE auditoria SET resultado = 'ALTERADO' WHERE id = (SELECT id FROM auditoria LIMIT 1);"); err == nil {
			t.Errorf("la base permitió UPDATE sobre auditoria con el usuario de la aplicación; la tabla no es append-only")
		}
		if _, err := execDatabase(t, "DELETE FROM auditoria WHERE id = (SELECT id FROM auditoria LIMIT 1);"); err == nil {
			t.Errorf("la base permitió DELETE sobre auditoria con el usuario de la aplicación; la tabla no es append-only")
		}
	})

	t.Run("BAC-13 reset de 2FA invalida el secreto y revoca sesiones", func(t *testing.T) {
		user := createActiveUser(t, adminToken, "reset2fa", "OPERATOR")
		session := loginWithTOTP(t, user.Email, user.Password, user.Secret)
		path := "/admin/users/" + user.ID + "/2fa/reset"

		if status, _ := requestValue(t, http.MethodPost, path, operatorToken, nil); status != http.StatusForbidden {
			t.Errorf("OPERATOR ejecutando reset 2FA: esperado 403, recibido %d", status)
		}
		if status, body := requestJSON(t, http.MethodPost, path, adminToken, nil); status != http.StatusOK && status != http.StatusNoContent {
			t.Fatalf("reset 2FA esperado 200/204, recibido %d: %#v", status, body)
		}
		if linked := queryDatabase(t, "SELECT totp_vinculado FROM usuarios WHERE id = '"+user.ID+"';"); linked != "f" {
			t.Errorf("totp_vinculado tras el reset esperado 'f', obtenido %q", linked)
		}
		if status, _ := requestJSON(t, http.MethodGet, "/account/profile", session.AccessToken, nil); status != http.StatusUnauthorized {
			t.Errorf("las sesiones activas deben revocarse tras el reset 2FA: esperado 401, recibido %d", status)
		}

		_, relogin := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{"email": user.Email, "password": user.Password})
		if relogin["totpVinculado"] != false {
			t.Errorf("tras el reset el login debe pedir re-enrolamiento (totpVinculado=false), recibido %#v", relogin["totpVinculado"])
		}
		preAuth := requiredString(t, relogin, "jwtTemporal")
		if status, _ := requestJSON(t, http.MethodGet, "/auth/2fa/qr", preAuth, nil); status != http.StatusOK {
			t.Errorf("tras el reset debe poder generarse un QR nuevo: esperado 200, recibido %d", status)
		}
		queryDatabase(t, "UPDATE usuarios SET ultimo_totp_periodo = 0 WHERE id = '"+user.ID+"';")
		oldCode, _ := totp.GenerateCode(user.Secret, time.Now().UTC())
		if status, _ := requestJSON(t, http.MethodPost, "/auth/2fa/verify", preAuth, map[string]any{"codigo": oldCode}); status == http.StatusOK {
			t.Errorf("un código del secreto anterior fue aceptado después del reset 2FA")
		}
	})

	t.Run("BAC-15 reset de contrasena genera temporal invalida la anterior y revoca sesiones", func(t *testing.T) {
		user := createActiveUser(t, adminToken, "resetpass", "OPERATOR")
		session := loginWithTOTP(t, user.Email, user.Password, user.Secret)
		path := "/admin/users/" + user.ID + "/password/reset"

		if status, _ := requestValue(t, http.MethodPost, path, operatorToken, nil); status != http.StatusForbidden {
			t.Errorf("OPERATOR ejecutando reset de contraseña: esperado 403, recibido %d", status)
		}
		if status, body := requestJSON(t, http.MethodPost, path, adminToken, nil); status != http.StatusOK && status != http.StatusNoContent {
			t.Fatalf("reset de contraseña esperado 200/204, recibido %d: %#v", status, body)
		}
		temporary := mailedSecret(t, user.Email, mailTemporaryPassword)

		if status := loginStatus(t, user.Email, user.Password); status != http.StatusUnauthorized {
			t.Errorf("la contraseña anterior debe dejar de funcionar: esperado 401, recibido %d", status)
		}
		status, relogin := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{"email": user.Email, "password": temporary})
		if status != http.StatusOK || relogin["cambioContrasenaRequerido"] != true {
			t.Errorf("la clave temporal debe permitir el login con cambioContrasenaRequerido=true, recibido %d: %#v", status, relogin)
		}
		if status, _ := requestJSON(t, http.MethodGet, "/account/profile", session.AccessToken, nil); status != http.StatusUnauthorized {
			t.Errorf("las sesiones activas deben revocarse tras el reset de contraseña: esperado 401, recibido %d", status)
		}
	})
}
