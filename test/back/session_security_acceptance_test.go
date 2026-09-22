package back_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signedRefreshToken(t *testing.T, userID, orgID string) (string, string) {
	t.Helper()
	jti := fmt.Sprintf("test-refresh-%d", time.Now().UnixNano())
	claims := jwt.MapClaims{
		"sub":    userID,
		"tipo":   "refresh",
		"org_id": orgID,
		"jti":    jti,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("no se pudo firmar refresh token de prueba: %v", err)
	}
	queryDatabase(t, fmt.Sprintf(
		"INSERT INTO sesiones_activas (id, usuario_id, jti_token, activa, fecha_expiracion, fecha_creacion) VALUES (uuid_generate_v7(), '%s', '%s', true, NOW() + interval '7 days', NOW());",
		userID, jti,
	))
	return signed, jti
}

// Este archivo valida el contrato de las tareas, fixes y mejoras de seguridad:
// - BAC-17: Cierre de sesión y revocación atómica de sesiones (Backend)
// - SEC-01: Refresh token en cookie HttpOnly en backend (Backend)
// - FIX-08: Auditoría y unificación del contrato de códigos de error (Backend)
func TestHitoSeguridadDeSesionesYCookies(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	if adminID == "" {
		t.Fatal("no se encontró usuario admin@elcentinela.com en base de datos")
	}

	t.Run("BAC-17 revocacion atomica de accessToken y refreshToken en logout", func(t *testing.T) {
		// 1. Iniciar sesión para obtener tokens reales emitidos por el backend
		adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)
		claims := decodeJWTClaims(t, adminToken)
		accessJTI, _ := claims["jti"].(string)

		// Crear una sesión de refresh en base de datos asociada a este usuario
		refreshTokenStr, refreshJTI := signedRefreshToken(t, adminID, orgID)

		// 2. Ejecutar logout enviando Authorization: Bearer <accessToken> y body con refreshToken
		logoutPayload := map[string]any{
			"refreshToken": refreshTokenStr,
		}
		status, _ := requestValue(t, http.MethodPost, "/auth/logout", adminToken, logoutPayload)
		if status != http.StatusNoContent && status != http.StatusOK {
			t.Fatalf("POST /api/auth/logout esperado 204 o 200, recibido %d", status)
		}

		// 3. El accessToken revocado DEBE ser rechazado inmediatamente con 401 en cualquier llamada protegida
		profStatus, profBody := requestJSON(t, http.MethodGet, "/account/profile", adminToken, nil)
		if profStatus != http.StatusUnauthorized {
			t.Errorf("BAC-17 acceso con accessToken post-logout: esperado 401 (TOKEN_REVOKED), recibido %d: %#v", profStatus, profBody)
		} else {
			errorCode, _ := profBody["errorCode"].(string)
			if errorCode != "TOKEN_REVOKED" && errorCode != "AUTH_FAILED" {
				t.Errorf("BAC-17 errorCode esperado TOKEN_REVOKED o AUTH_FAILED, recibido: %q", errorCode)
			}
		}

		// 4. Verificar en PostgreSQL que las sesiones activas asociadas quedaron inactivas (activa = false)
		if accessJTI != "" {
			accessActiva := queryDatabase(t, fmt.Sprintf("SELECT activa FROM sesiones_activas WHERE jti_token = '%s';", accessJTI))
			if accessActiva == "t" {
				t.Errorf("BAC-17 el accessToken (jti: %s) sigue con activa=true en sesiones_activas tras logout", accessJTI)
			}
		}
		refreshActiva := queryDatabase(t, fmt.Sprintf("SELECT activa FROM sesiones_activas WHERE jti_token = '%s';", refreshJTI))
		if refreshActiva == "t" {
			t.Errorf("BAC-17 el refreshToken (jti: %s) sigue con activa=true en sesiones_activas tras logout", refreshJTI)
		}

		// 5. Verificar que la tabla de auditoría registró formalmente la acción LOGOUT
		auditRows := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM auditoria WHERE usuario_id = '%s' AND accion = 'LOGOUT';", adminID))
		if auditRows == "0" {
			t.Errorf("BAC-17 no se encontró registro de auditoría con accion 'LOGOUT' para el usuario %s", adminID)
		}
	})

	t.Run("SEC-01 emision de refreshToken en cookie HttpOnly y soporte de rotacion", func(t *testing.T) {
		// 1. Simular petición de verificación 2FA para comprobar si Set-Cookie emite HttpOnly
		verifyPayload, _ := json.Marshal(map[string]any{
			"codigo": "123456",
		})
		req, err := http.NewRequest(http.MethodPost, apiURL+"/auth/2fa/verify", bytes.NewReader(verifyPayload))
		if err != nil {
			t.Fatalf("no se pudo construir petición 2fa verify: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer dummy-pre-auth-token")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("falló llamada a 2fa/verify: %v", err)
		}
		defer resp.Body.Close()

		// Inspeccionar cookies emitidas
		var foundHttpOnlyRefreshCookie bool
		for _, cookie := range resp.Cookies() {
			if strings.Contains(strings.ToLower(cookie.Name), "refresh") {
				if cookie.HttpOnly {
					foundHttpOnlyRefreshCookie = true
				}
			}
		}

		// SEC-01 exige que el refresh token se emita con Set-Cookie HttpOnly
		if !foundHttpOnlyRefreshCookie {
			t.Errorf("SEC-01 esperado Set-Cookie HttpOnly para refreshToken en POST /api/auth/2fa/verify, recibido cookies: %v", resp.Cookies())
		}

		// 2. Probar si POST /api/auth/refresh acepta cookie sin requerir body JSON
		refreshReq, err := http.NewRequest(http.MethodPost, apiURL+"/auth/refresh", nil)
		if err != nil {
			t.Fatalf("no se pudo construir petición /auth/refresh: %v", err)
		}
		refreshReq.AddCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    "dummy-cookie-value",
			HttpOnly: true,
		})
		refreshResp, err := client.Do(refreshReq)
		if err != nil {
			t.Fatalf("falló llamada /auth/refresh con cookie: %v", err)
		}
		defer refreshResp.Body.Close()

		// Si responde 400 INVALID_REQUEST es porque exige body y no soporta lectura por cookie
		if refreshResp.StatusCode == http.StatusBadRequest {
			t.Errorf("SEC-01 POST /api/auth/refresh devolvió 400: no lee refreshToken desde cookie HttpOnly")
		}

		// 3. Probar si POST /api/auth/logout invalida la cookie enviando Max-Age=0 o fecha expirada
		logoutReq, err := http.NewRequest(http.MethodPost, apiURL+"/auth/logout", nil)
		if err != nil {
			t.Fatalf("no se pudo construir petición /auth/logout: %v", err)
		}
		logoutReq.AddCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    "dummy-cookie-value",
			HttpOnly: true,
		})
		logoutResp, err := client.Do(logoutReq)
		if err != nil {
			t.Fatalf("falló llamada /auth/logout con cookie: %v", err)
		}
		defer logoutResp.Body.Close()

		var cookieExpired bool
		for _, c := range logoutResp.Cookies() {
			if strings.Contains(strings.ToLower(c.Name), "refresh") {
				if c.MaxAge <= 0 {
					cookieExpired = true
				}
			}
		}
		if !cookieExpired {
			t.Errorf("SEC-01 POST /api/auth/logout esperado Set-Cookie con Max-Age <= 0 para revocar cookie de refresh")
		}
	})

	t.Run("FIX-08 backend emite estructura estandar de error con errorCode y message", func(t *testing.T) {
		casos := []struct {
			nombre       string
			metodo       string
			path         string
			token        string
			payload      any
			statusEsper  int
			codeEsperado string
		}{
			{
				nombre:       "payload incompleto en login",
				metodo:       http.MethodPost,
				path:         "/auth/login",
				token:        "",
				payload:      map[string]any{"email": "solo-email@elcentinela.com"},
				statusEsper:  http.StatusBadRequest,
				codeEsperado: "INVALID_REQUEST",
			},
			{
				nombre:       "credenciales inválidas en login",
				metodo:       http.MethodPost,
				path:         "/auth/login",
				token:        "",
				payload:      map[string]any{"email": "admin@elcentinela.com", "password": "PasswordInvalida123!"},
				statusEsper:  http.StatusUnauthorized,
				codeEsperado: "AUTH_FAILED",
			},
			{
				nombre:       "ruta inexistente devuelve NOT_FOUND",
				metodo:       http.MethodGet,
				path:         "/ruta-completamente-inexistente-xyz",
				token:        "",
				payload:      nil,
				statusEsper:  http.StatusNotFound,
				codeEsperado: "NOT_FOUND",
			},
		}

		for _, tc := range casos {
			t.Run(tc.nombre, func(t *testing.T) {
				status, resp := requestJSON(t, tc.metodo, tc.path, tc.token, tc.payload)
				if status != tc.statusEsper {
					t.Errorf("status esperado %d, recibido %d", tc.statusEsper, status)
				}
				errorCode, ok := resp["errorCode"].(string)
				if !ok || errorCode != tc.codeEsperado {
					t.Errorf("errorCode esperado %q, recibido %#v en respuesta: %#v", tc.codeEsperado, resp["errorCode"], resp)
				}
				message, ok := resp["message"].(string)
				if !ok || strings.TrimSpace(message) == "" {
					t.Errorf("message esperado no vacío en respuesta: %#v", resp)
				}
			})
		}
	})
}
