package back_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// Contrato de seguridad de sesiones (documentacion/actual.md):
//   - BAC-17: logout con revocación atómica de access + refresh y auditoría LOGOUT
//   - SEC-01: refresh token exclusivamente en cookie HttpOnly
//   - FIX-08: contrato uniforme { errorCode, message } (terminado.md, regresión)
//
// Todos los tokens se obtienen del circuito real login -> 2FA; la prueba no firma tokens propios.
func TestHitoSeguridadDeSesionesYCookies(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)
	user := createActiveUser(t, adminToken, "sesiones", "OPERATOR")

	// refreshCredential devuelve el refresh token por el canal que el backend use hoy:
	// body JSON (contrato previo) o cookie (SEC-01).
	refreshCredential := func(session loginResult) (map[string]any, []*http.Cookie) {
		if session.RefreshCookie != nil {
			return nil, []*http.Cookie{session.RefreshCookie}
		}
		return map[string]any{"refreshToken": session.RefreshToken}, nil
	}

	t.Run("BAC-17 logout revoca atomicamente access y refresh y registra LOGOUT", func(t *testing.T) {
		session := loginWithTOTP(t, user.Email, user.Password, user.Secret)
		accessJTI, _ := decodeJWTClaims(t, session.AccessToken)["jti"].(string)
		if session.RefreshToken == "" && session.RefreshCookie == nil {
			t.Fatalf("2fa/verify no entregó refresh token ni en JSON ni en cookie: %s", session.Verify.RawBody)
		}

		if status, body := requestJSON(t, http.MethodGet, "/account/profile", session.AccessToken, nil); status != http.StatusOK {
			t.Fatalf("precondición: el access token recién emitido debe funcionar, recibido %d: %#v", status, body)
		}
		auditBefore := auditCount(t, user.ID, "LOGOUT")

		payload, cookies := refreshCredential(session)
		if status := requestRaw(t, http.MethodPost, "/auth/logout", "", payload, cookies...).Status; status != http.StatusUnauthorized {
			t.Errorf("POST /auth/logout sin Authorization debe exigir el access token (401), recibido %d", status)
		}

		logout := requestRaw(t, http.MethodPost, "/auth/logout", session.AccessToken, payload, cookies...)
		if logout.Status != http.StatusNoContent {
			t.Fatalf("POST /auth/logout esperado 204, recibido %d: %s", logout.Status, logout.RawBody)
		}

		status, body := requestJSON(t, http.MethodGet, "/account/profile", session.AccessToken, nil)
		if status != http.StatusUnauthorized || body["errorCode"] != "TOKEN_REVOKED" {
			t.Errorf("access token tras logout: esperado 401 TOKEN_REVOKED, recibido %d: %#v", status, body)
		}

		refreshAfter := requestRaw(t, http.MethodPost, "/auth/refresh", "", payload, cookies...)
		if refreshAfter.Status != http.StatusUnauthorized {
			t.Errorf("refresh token tras logout: esperado 401, recibido %d: %s", refreshAfter.Status, refreshAfter.RawBody)
		}

		inactive := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM sesiones_activas WHERE usuario_id = '%s' AND jti_token = '%s' AND activa = false;", user.ID, accessJTI))
		if inactive != "1" {
			t.Errorf("la sesión del access token (jti %s) no quedó con activa=false", accessJTI)
		}
		stillActive := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM sesiones_activas WHERE usuario_id = '%s' AND activa = true AND fecha_creacion >= (SELECT fecha_creacion FROM sesiones_activas WHERE jti_token = '%s');", user.ID, accessJTI))
		if stillActive != "0" {
			t.Errorf("tras el logout quedaron %s sesiones activas de este login (se esperaba revocar access y refresh)", stillActive)
		}

		if after := auditCount(t, user.ID, "LOGOUT"); after != auditBefore+1 {
			t.Errorf("auditoría LOGOUT esperada %d, encontrada %d", auditBefore+1, after)
		}
	})

	t.Run("SEC-01 refresh token solo en cookie HttpOnly", func(t *testing.T) {
		session := loginWithTOTP(t, user.Email, user.Password, user.Secret)

		if _, exposed := session.Verify.Body["refreshToken"]; exposed {
			t.Errorf("POST /auth/2fa/verify todavía devuelve refreshToken en el JSON (entregable 3)")
		}
		cookie := session.RefreshCookie
		if cookie == nil {
			t.Fatalf("SEC-01 no implementada: POST /auth/2fa/verify no emite Set-Cookie con el refresh token (cabeceras: %v)", session.Verify.Header.Values("Set-Cookie"))
		}
		if !cookie.HttpOnly {
			t.Errorf("la cookie %q no tiene HttpOnly", cookie.Name)
		}
		// compose.yaml no define ninguna variable de cookies: la configuración por defecto
		// debe ser la segura (Secure), y relajarla solo por configuración explícita.
		if !cookie.Secure {
			t.Errorf("con la configuración por defecto la cookie %q debe ser Secure", cookie.Name)
		}
		if cookie.SameSite == http.SameSiteDefaultMode || cookie.SameSite == http.SameSiteNoneMode && !cookie.Secure {
			t.Errorf("la cookie %q debe declarar SameSite (Lax/Strict, o None solo con Secure); Set-Cookie: %v", cookie.Name, session.Verify.Header.Values("Set-Cookie"))
		}
		if cookie.Path == "" || cookie.Path == "/" {
			t.Errorf("la cookie %q debe restringir Path a las rutas de auth (p. ej. /api/auth), recibido %q", cookie.Name, cookie.Path)
		}

		refresh := requestRaw(t, http.MethodPost, "/auth/refresh", "", nil, cookie)
		if refresh.Status != http.StatusOK {
			t.Fatalf("POST /auth/refresh solo con la cookie: esperado 200, recibido %d: %s", refresh.Status, refresh.RawBody)
		}
		newAccess := requiredString(t, refresh.Body, "accessToken")
		if _, exposed := refresh.Body["refreshToken"]; exposed {
			t.Errorf("POST /auth/refresh devuelve refreshToken en el JSON")
		}
		if next := refresh.cookie("refresh"); next != nil && next.Value != "" {
			cookie = next
		}

		invalid := requestRaw(t, http.MethodPost, "/auth/refresh", "", nil, &http.Cookie{Name: cookie.Name, Value: "cookie-invalida"})
		if invalid.Status != http.StatusUnauthorized || invalid.Body["errorCode"] == nil {
			t.Errorf("cookie inválida: esperado 401 con errorCode, recibido %d: %s", invalid.Status, invalid.RawBody)
		}

		logout := requestRaw(t, http.MethodPost, "/auth/logout", newAccess, nil, cookie)
		if logout.Status != http.StatusNoContent {
			t.Fatalf("POST /auth/logout solo con cookie + Bearer: esperado 204, recibido %d: %s", logout.Status, logout.RawBody)
		}
		cleared := logout.cookie("refresh")
		if cleared == nil || (cleared.MaxAge >= 0 && cleared.Value != "") {
			t.Errorf("el logout debe eliminar la cookie (Max-Age=0 o valor vacío); Set-Cookie: %v", logout.Header.Values("Set-Cookie"))
		}
		if again := requestRaw(t, http.MethodPost, "/auth/refresh", "", nil, cookie); again.Status != http.StatusUnauthorized {
			t.Errorf("la cookie revocada no debe renovar la sesión: esperado 401, recibido %d", again.Status)
		}

		logs, _ := runCompose("logs", "--no-color", "backend")
		if strings.Contains(logs, cookie.Value) {
			t.Errorf("el valor del refresh token aparece en los logs del backend")
		}
	})

	// Integración front <-> back: el frontend en origin/main (SEC-02, commits cb208f5 / f0f7218)
	// ya envía POST /auth/logout y POST /auth/refresh con body {} y credentials: 'include',
	// confiando en la cookie de SEC-01. Esta prueba reproduce exactamente esas peticiones.
	t.Run("SEC-01 SEC-02 integracion el backend acepta logout y refresh tal como los envia el frontend", func(t *testing.T) {
		session := loginWithTOTP(t, user.Email, user.Password, user.Secret)
		var cookies []*http.Cookie
		if session.RefreshCookie != nil {
			cookies = append(cookies, session.RefreshCookie)
		}

		refresh := requestRaw(t, http.MethodPost, "/auth/refresh", "", map[string]any{}, cookies...)
		if refresh.Status != http.StatusOK {
			t.Errorf("POST /auth/refresh con body {} (petición real del frontend): esperado 200, recibido %d: %s", refresh.Status, refresh.RawBody)
		}

		logout := requestRaw(t, http.MethodPost, "/auth/logout", session.AccessToken, map[string]any{}, cookies...)
		if logout.Status != http.StatusNoContent {
			t.Errorf("POST /auth/logout con Bearer y body {} (petición real del frontend): esperado 204, recibido %d: %s", logout.Status, logout.RawBody)
		}
		if status, _ := requestJSON(t, http.MethodGet, "/account/profile", session.AccessToken, nil); status != http.StatusUnauthorized {
			t.Errorf("tras el logout iniciado desde el frontend el access token debe quedar revocado: esperado 401, recibido %d", status)
		}
	})

	t.Run("FIX-08 backend emite estructura estandar de error con errorCode y message", func(t *testing.T) {
		casos := []struct {
			nombre       string
			metodo       string
			path         string
			payload      any
			statusEsper  int
			codeEsperado string
		}{
			{"payload incompleto en login", http.MethodPost, "/auth/login", map[string]any{"email": "solo-email@elcentinela.com"}, http.StatusBadRequest, "INVALID_REQUEST"},
			{"credenciales inválidas en login", http.MethodPost, "/auth/login", map[string]any{"email": "admin@elcentinela.com", "password": "PasswordInvalida123!"}, http.StatusUnauthorized, "AUTH_FAILED"},
			{"ruta inexistente devuelve NOT_FOUND", http.MethodGet, "/ruta-completamente-inexistente-xyz", nil, http.StatusNotFound, "NOT_FOUND"},
			{"ruta protegida sin token devuelve MISSING_TOKEN", http.MethodGet, "/account/profile", nil, http.StatusUnauthorized, "MISSING_TOKEN"},
		}
		for _, tc := range casos {
			t.Run(tc.nombre, func(t *testing.T) {
				status, resp := requestJSON(t, tc.metodo, tc.path, "", tc.payload)
				if status != tc.statusEsper {
					t.Errorf("status esperado %d, recibido %d", tc.statusEsper, status)
				}
				if resp["errorCode"] != tc.codeEsperado {
					t.Errorf("errorCode esperado %q, recibido %#v", tc.codeEsperado, resp["errorCode"])
				}
				if message, ok := resp["message"].(string); !ok || strings.TrimSpace(message) == "" {
					t.Errorf("message esperado no vacío en respuesta: %#v", resp)
				}
			})
		}
	})
}
