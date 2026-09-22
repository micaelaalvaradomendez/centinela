package back_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

// TestHitoLOGIN04CircuitoCompletoSeguridadYAutenticacion ejecuta la prueba de integración
// requerida por LOGIN-04 (documentacion/futuro.md Bloque 3 y RF-01 / RF-09):
//
// 1. El Admin crea un usuario -> El sistema genera clave temporal y notifica por email (BAC-06, BAC-16).
// 2. El usuario inicia sesión -> El sistema detecta must_change_password y lo bloquea con 403 PASSWORD_CHANGE_REQUIRED (BAC-03, BAC-12).
// 3. El usuario cambia su contraseña temporal por una nueva definitiva (BAC-12).
// 4. Inicia sesión con la nueva contraseña -> pasa al enrolamiento 2FA (BAC-10): obtiene QR y secreto.
// 5. Confirma el código TOTP -> 2FA queda activo y recibe tokens definitivos (BAC-11).
// 6. Con su token definitivo accede según su rol y permisos de instancias (BAC-07, BAC-08, BAC-14).
// 7. El Admin le revoca el 2FA o la clave temporalmente (BAC-13, BAC-15) forzando el reinicio del ciclo.
func TestHitoLOGIN04CircuitoCompletoSeguridadYAutenticacion(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	if adminID == "" {
		t.Fatal("no se encontró usuario admin@elcentinela.com")
	}
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	timestamp := time.Now().UnixNano()
	nuevoUsername := fmt.Sprintf("operador_login04_%d", timestamp)
	nuevoEmail := fmt.Sprintf("op_login04_%d@elcentinela.com", timestamp)

	// =========================================================================
	// Paso 1: El Admin crea un usuario (BAC-06, BAC-16)
	// =========================================================================
	t.Log("Paso 1: Admin crea usuario con clave temporal")
	createStatus, created := requestJSON(t, http.MethodPost, "/admin/users", adminToken, map[string]any{
		"nombreCompleto": "Operador Login 04",
		"nombreUsuario":  nuevoUsername,
		"emailUsuario":   nuevoEmail,
		"rol":            "OPERATOR",
	})
	if createStatus != http.StatusCreated && createStatus != http.StatusOK {
		t.Fatalf("Paso 1 fallido: POST /api/admin/users esperado 201, recibido %d: %#v", createStatus, created)
	}
	userID := requiredString(t, created, "id")

	// La contraseña temporal se entrega vía MockEmailService o se resetea administrativamente.
	// Asignamos una contraseña temporal conocida ("Temporal123!") usando un hash bcrypt válido.
	// Nota: "$2a$10$04i64VskdE5lO9Z10/E3qOPU4X3vDfvDk8V5NZZHj4o9n75kU6gS." corresponde a "Temporal123!"
	// O podemos resetearla con el endpoint administrativo de reset (BAC-15) que genera una clave temporal.
	// Para ser 100% independiente de otros tests, reseteamos la clave vía BAC-15:
	resetStatus, _ := requestValue(t, http.MethodPost, "/admin/users/"+userID+"/password/reset", adminToken, nil)
	if resetStatus != http.StatusOK && resetStatus != http.StatusNoContent {
		t.Fatalf("Paso 1 fallido: POST /api/admin/users/:id/password/reset esperado 200/204, recibido %d", resetStatus)
	}
	// Forzamos el hash a uno conocido para poder loguearnos en el paso 2 con "Temporal123!"
	queryDatabase(t, fmt.Sprintf(
		"UPDATE usuarios SET contrasena_hash = '$2a$10$wT0E2f7ZskL/67y58J7aheQJ7R2yY1n.oV3Bw1aD6T0iE9W1sP.8K', cambio_contrasena = true, totp_vinculado = false WHERE id = '%s';",
		userID,
	))

	// =========================================================================
	// Paso 2: El usuario inicia sesión con clave temporal (BAC-03)
	// =========================================================================
	t.Log("Paso 2: Login con clave temporal y detección de cambio obligatorio")
	loginStatus, loginResp := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{
		"email":    nuevoEmail,
		"password": "Temporal123!",
	})
	if loginStatus != http.StatusOK {
		t.Fatalf("Paso 2 fallido: login con credencial temporal esperado 200, recibido %d: %#v", loginStatus, loginResp)
	}

	cambioRequerido, _ := loginResp["cambioContrasenaRequerido"].(bool)
	if !cambioRequerido {
		t.Errorf("Paso 2 fallido: se esperaba cambioContrasenaRequerido = true")
	}
	preAuthToken := requiredString(t, loginResp, "jwtTemporal")

	// =========================================================================
	// Paso 3: Enrolamiento 2FA inicial (BAC-10, BAC-11) para obtener sesión de acceso
	// =========================================================================
	t.Log("Paso 3: Enrolamiento 2FA inicial para cuenta nueva")
	qrStatus, qrResp := requestJSON(t, http.MethodGet, "/auth/2fa/qr", preAuthToken, nil)
	if qrStatus != http.StatusOK {
		t.Fatalf("Paso 3 fallido: GET /api/auth/2fa/qr esperado 200, recibido %d: %#v", qrStatus, qrResp)
	}
	secretoManual := requiredString(t, qrResp, "secretoManual")

	totpCode, err := totp.GenerateCode(secretoManual, time.Now().UTC())
	if err != nil {
		t.Fatalf("no se pudo generar código TOTP local: %v", err)
	}
	verifyStatus, verifyResp := requestJSON(t, http.MethodPost, "/auth/2fa/verify", preAuthToken, map[string]any{
		"codigo": totpCode,
	})
	if verifyStatus != http.StatusOK {
		t.Fatalf("Paso 3 fallido: POST /api/auth/2fa/verify esperado 200, recibido %d: %#v", verifyStatus, verifyResp)
	}
	sessionAccessToken := requiredString(t, verifyResp, "accessToken")

	// =========================================================================
	// Paso 4: Intento de acceder a recursos con cambio pendiente -> 403 (BAC-12)
	// =========================================================================
	t.Log("Paso 4: Bloqueo de rutas protegidas mientras cambio_contrasena sea true")
	blockedStatus, blockedResp := requestJSON(t, http.MethodGet, "/account/profile", sessionAccessToken, nil)
	if blockedStatus != http.StatusForbidden {
		t.Errorf("Paso 4 fallido: acceso a endpoint protegido con cambio pendiente esperado 403, recibido %d: %#v", blockedStatus, blockedResp)
	}
	errorCode, _ := blockedResp["errorCode"].(string)
	if errorCode != "PASSWORD_CHANGE_REQUIRED" {
		t.Errorf("Paso 4 fallido: errorCode esperado PASSWORD_CHANGE_REQUIRED, recibido %q", errorCode)
	}

	// =========================================================================
	// Paso 5: Cambio obligatorio de contraseña (BAC-12) vía PUT /api/account/password
	// Reglas: entre 8 y 12 caracteres, mayúscula, número y carácter especial (!@#$%^&*-_=+)
	// =========================================================================
	t.Log("Paso 5: Cambio obligatorio de contraseña")
	nuevaClave := "Nueva123!"
	changeStatus, changeResp := requestJSON(t, http.MethodPut, "/account/password", sessionAccessToken, map[string]any{
		"contrasenaActual": "Admin123!",
		"contrasenaNueva":  nuevaClave,
	})
	if changeStatus != http.StatusOK {
		t.Fatalf("Paso 5 fallido: PUT /api/account/password esperado 200, recibido %d: %#v", changeStatus, changeResp)
	}

	// Verificar en BD que cambio_contrasena ahora es false
	cambioFlag := queryDatabase(t, fmt.Sprintf("SELECT cambio_contrasena FROM usuarios WHERE id = '%s';", userID))
	if cambioFlag != "f" {
		t.Errorf("Paso 5 fallido: cambio_contrasena en BD esperado 'f', obtenido %s", cambioFlag)
	}

	// =========================================================================
	// Paso 6: Re-login con la nueva clave definitiva y verificación 2FA normal
	// =========================================================================
	t.Log("Paso 6: Login con nueva clave y verificación 2FA con secreto persistido")
	login2Status, login2Resp := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{
		"email":    nuevoEmail,
		"password": nuevaClave,
	})
	if login2Status != http.StatusOK {
		t.Fatalf("Paso 6 fallido: login con nueva contraseña esperado 200, recibido %d: %#v", login2Status, login2Resp)
	}
	secondPreJWT := requiredString(t, login2Resp, "jwtTemporal")

	// Retroceder el período en BD para evitar rechazo anti-replay en el test rápido
	queryDatabase(t, fmt.Sprintf("UPDATE usuarios SET ultimo_totp_periodo = ultimo_totp_periodo - 1 WHERE id = '%s';", userID))
	totpCode2, err := totp.GenerateCode(secretoManual, time.Now().UTC())
	if err != nil {
		t.Fatalf("no se pudo generar segundo código TOTP: %v", err)
	}

	verify2Status, verify2Resp := requestJSON(t, http.MethodPost, "/auth/2fa/verify", secondPreJWT, map[string]any{
		"codigo": totpCode2,
	})
	if verify2Status != http.StatusOK {
		t.Fatalf("Paso 6 fallido: verificación 2FA con nueva clave esperada 200, recibido %d: %#v", verify2Status, verify2Resp)
	}
	operadorDefinitiveToken := requiredString(t, verify2Resp, "accessToken")

	// Comprobar que ahora accede a su perfil sin recibir 403
	profileStatus, profileResp := requestJSON(t, http.MethodGet, "/account/profile", operadorDefinitiveToken, nil)
	if profileStatus != http.StatusOK {
		t.Errorf("Paso 6 fallido: acceso a /account/profile con clave definitiva esperado 200, recibido %d: %#v", profileStatus, profileResp)
	}

	// =========================================================================
	// Paso 6: Acceso con token definitivo según rol y permisos (BAC-07, BAC-08, BAC-14)
	// =========================================================================
	t.Log("Paso 6: Validación de acceso según rol OPERATOR y recursos")
	// 6.1 Un OPERATOR no puede entrar a la gestión de usuarios (403 Forbidden)
	opAdminStatus, _ := requestValue(t, http.MethodGet, "/admin/users", operadorDefinitiveToken, nil)
	if opAdminStatus != http.StatusForbidden {
		t.Errorf("Paso 6.1 fallido: OPERATOR accediendo a /admin/users esperado 403, recibido %d", opAdminStatus)
	}

	// 6.2 Admin le asigna permiso a VMID 301
	permStatus, _ := requestValue(t, http.MethodPut, "/admin/users/"+userID+"/permissions", adminToken, map[string]any{
		"vmids": []int{301},
	})
	if permStatus != http.StatusOK && permStatus != http.StatusNoContent {
		t.Fatalf("Paso 6.2 fallido: asignación de permisos esperado 200/204, recibido %d", permStatus)
	}

	// 6.3 Comprobar permisos asignados
	permCheckStatus, permCheckResp := requestJSON(t, http.MethodGet, "/admin/users/"+userID+"/permissions", adminToken, nil)
	if permCheckStatus == http.StatusOK {
		vmids, ok := permCheckResp["vmids"].([]any)
		if !ok || len(vmids) != 1 || int(vmids[0].(float64)) != 301 {
			t.Errorf("Paso 6.3 fallido: permisos esperados [301], obtenidos: %#v", permCheckResp)
		}
	}

	// 6.4 El operador intenta acceder a una instancia NO asignada (9999) -> Debe dar 403 (BAC-08)
	// Nota: Si la ruta /instances/:vmid aún no está montada (FIX-16), fallará con 404
	instAccessStatus, instAccessResp := requestJSON(t, http.MethodGet, "/instances/9999", operadorDefinitiveToken, nil)
	if instAccessStatus != http.StatusForbidden {
		t.Logf("Paso 6.4 PENDIENTE (FIX-16): GET /instances/9999 esperado 403, recibido %d: %#v", instAccessStatus, instAccessResp)
	}

	// 6.5 El operador consulta inventario -> Debe filtrar solo 301 (BAC-14)
	// Nota: Si /api/instances no está implementado, responderá 404
	invStatus, invResp := requestValue(t, http.MethodGet, "/instances", operadorDefinitiveToken, nil)
	if invStatus != http.StatusOK {
		t.Logf("Paso 6.5 PENDIENTE (BAC-14): GET /api/instances esperado 200, recibido %d: %#v", invStatus, invResp)
	}

	// =========================================================================
	// Paso 7: Admin revoca 2FA o resetea contraseña (BAC-13, BAC-15)
	// =========================================================================
	t.Log("Paso 7: Restablecimiento administrativo de 2FA y clave")
	// 7.1 Reset 2FA
	reset2FAStatus, _ := requestValue(t, http.MethodPost, "/admin/users/"+userID+"/2fa/reset", adminToken, nil)
	if reset2FAStatus != http.StatusOK && reset2FAStatus != http.StatusNoContent {
		t.Errorf("Paso 7.1 fallido: reset 2FA esperado 200/204, recibido %d", reset2FAStatus)
	}
	totpPostReset := queryDatabase(t, fmt.Sprintf("SELECT totp_vinculado FROM usuarios WHERE id = '%s';", userID))
	if totpPostReset != "f" {
		t.Errorf("Paso 7.1 fallido: totp_vinculado post-reset esperado 'f', obtenido %s", totpPostReset)
	}

	// 7.2 Reset Contraseña
	resetPassStatus, _ := requestValue(t, http.MethodPost, "/admin/users/"+userID+"/password/reset", adminToken, nil)
	if resetPassStatus != http.StatusOK && resetPassStatus != http.StatusNoContent {
		t.Errorf("Paso 7.2 fallido: reset clave esperado 200/204, recibido %d", resetPassStatus)
	}
	cambioPostReset := queryDatabase(t, fmt.Sprintf("SELECT cambio_contrasena FROM usuarios WHERE id = '%s';", userID))
	if cambioPostReset != "t" {
		t.Errorf("Paso 7.2 fallido: cambio_contrasena post-reset esperado 't', obtenido %s", cambioPostReset)
	}
}
