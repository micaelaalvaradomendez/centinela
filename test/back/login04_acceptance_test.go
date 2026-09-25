package back_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestHitoLOGIN04CircuitoCompletoSeguridadYAutenticacion recorre de punta a punta el circuito
// de LOGIN-04 (RF-01 / RF-09) con credenciales reales, sin alterar hashes en la base:
//
//  1. El Admin crea un usuario; la clave temporal llega por MockEmailService (BAC-06, BAC-16).
//  2. Login con la clave temporal -> cambioContrasenaRequerido = true (BAC-03, BAC-12).
//  3. Enrolamiento 2FA inicial con QR (BAC-10, BAC-11).
//  4. Con cambio pendiente, las rutas protegidas responden 403 PASSWORD_CHANGE_REQUIRED (BAC-12).
//  5. Cambio obligatorio de contraseña (BAC-12).
//  6. Re-login con 2FA persistido; acceso según rol e instancias (BAC-07, BAC-08, BAC-14).
//  7. El Admin resetea 2FA y contraseña, reiniciando el ciclo (BAC-13, BAC-15).
func TestHitoLOGIN04CircuitoCompletoSeguridadYAutenticacion(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	suffix := time.Now().UnixNano()
	email := fmt.Sprintf("op_login04_%d@elcentinela.com", suffix)

	t.Log("Paso 1: alta con clave temporal entregada por correo")
	userID, temporary := createUserThroughAPI(t, adminToken, fmt.Sprintf("operador_login04_%d", suffix), email, "OPERATOR")

	t.Log("Paso 2: login con clave temporal")
	status, loginResp := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{"email": email, "password": temporary})
	if status != http.StatusOK {
		t.Fatalf("Paso 2: login con clave temporal esperado 200, recibido %d: %#v", status, loginResp)
	}
	if loginResp["cambioContrasenaRequerido"] != true {
		t.Errorf("Paso 2: se esperaba cambioContrasenaRequerido = true, recibido %#v", loginResp["cambioContrasenaRequerido"])
	}
	if loginResp["totpVinculado"] != false {
		t.Errorf("Paso 2: una cuenta nueva debe tener totpVinculado = false, recibido %#v", loginResp["totpVinculado"])
	}

	t.Log("Paso 3: enrolamiento 2FA")
	first := loginWithTOTP(t, email, temporary, "")

	t.Log("Paso 4: bloqueo mientras el cambio de clave esté pendiente")
	status, blocked := requestJSON(t, http.MethodGet, "/account/profile", first.AccessToken, nil)
	if status != http.StatusForbidden || blocked["errorCode"] != "PASSWORD_CHANGE_REQUIRED" {
		t.Errorf("Paso 4: esperado 403 PASSWORD_CHANGE_REQUIRED, recibido %d: %#v", status, blocked)
	}

	t.Log("Paso 5: cambio obligatorio de contraseña")
	finalPassword := "Nueva123!"
	status, weak := requestJSON(t, http.MethodPut, "/account/password", first.AccessToken, map[string]any{"contrasenaActual": temporary, "contrasenaNueva": "debil"})
	if status == http.StatusOK {
		t.Errorf("Paso 5: una contraseña débil fue aceptada: %#v", weak)
	}
	status, changed := requestJSON(t, http.MethodPut, "/account/password", first.AccessToken, map[string]any{"contrasenaActual": temporary, "contrasenaNueva": finalPassword})
	if status != http.StatusOK {
		t.Fatalf("Paso 5: PUT /api/account/password esperado 200, recibido %d: %#v", status, changed)
	}
	if flag := queryDatabase(t, fmt.Sprintf("SELECT cambio_contrasena FROM usuarios WHERE id = '%s';", userID)); flag != "f" {
		t.Errorf("Paso 5: cambio_contrasena esperado 'f', obtenido %q", flag)
	}
	if status, _ := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{"email": email, "password": temporary}); status != http.StatusUnauthorized {
		t.Errorf("Paso 5: la clave temporal debe dejar de funcionar, recibido %d", status)
	}

	t.Log("Paso 6: re-login con el secreto 2FA persistido y acceso según rol")
	session := loginWithTOTP(t, email, finalPassword, first.Secret)
	operatorToken := session.AccessToken
	if status, body := requestJSON(t, http.MethodGet, "/account/profile", operatorToken, nil); status != http.StatusOK {
		t.Errorf("Paso 6: /account/profile con clave definitiva esperado 200, recibido %d: %#v", status, body)
	}
	if status, _ := requestValue(t, http.MethodGet, "/admin/users", operatorToken, nil); status != http.StatusForbidden {
		t.Errorf("Paso 6.1: OPERATOR en /admin/users esperado 403, recibido %d", status)
	}

	if status, _ := requestValue(t, http.MethodPut, "/admin/users/"+userID+"/permissions", adminToken, permissionsPayload(101)); status != http.StatusNoContent {
		t.Fatalf("Paso 6.2: asignación de permisos esperado 204, recibido %d", status)
	}
	if status, body := requestJSON(t, http.MethodGet, "/instances/103", operatorToken, nil); status != http.StatusForbidden || body["errorCode"] != "INSTANCE_ACCESS_DENIED" {
		t.Errorf("Paso 6.3: instancia no asignada esperado 403 INSTANCE_ACCESS_DENIED, recibido %d: %#v", status, body)
	}
	status, inventory := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
	if status != http.StatusOK || len(inventory) != 1 || inventory[0].(map[string]any)["id"] != float64(101) {
		t.Errorf("Paso 6.4: el inventario del OPERATOR debe contener solo la 101, recibido %d: %#v", status, inventory)
	}

	t.Log("Paso 7: resets administrativos")
	if status, _ := requestValue(t, http.MethodPost, "/admin/users/"+userID+"/2fa/reset", adminToken, nil); status != http.StatusOK && status != http.StatusNoContent {
		t.Errorf("Paso 7.1: reset 2FA esperado 200/204, recibido %d", status)
	}
	if linked := queryDatabase(t, fmt.Sprintf("SELECT totp_vinculado FROM usuarios WHERE id = '%s';", userID)); linked != "f" {
		t.Errorf("Paso 7.1: totp_vinculado esperado 'f', obtenido %q", linked)
	}
	if status, _ := requestValue(t, http.MethodPost, "/admin/users/"+userID+"/password/reset", adminToken, nil); status != http.StatusOK && status != http.StatusNoContent {
		t.Errorf("Paso 7.2: reset de clave esperado 200/204, recibido %d", status)
	}
	if flag := queryDatabase(t, fmt.Sprintf("SELECT cambio_contrasena FROM usuarios WHERE id = '%s';", userID)); flag != "t" {
		t.Errorf("Paso 7.2: cambio_contrasena esperado 't', obtenido %q", flag)
	}
	if status, _ := requestJSON(t, http.MethodGet, "/account/profile", operatorToken, nil); status != http.StatusUnauthorized {
		t.Errorf("Paso 7: la sesión previa debe quedar revocada tras los resets, recibido %d", status)
	}
}
