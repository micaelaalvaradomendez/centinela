package back_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
)

const testJWTSecret = "test-only-jwt-secret-at-least-32-characters"

func TestTareasBackendEIntegracion(t *testing.T) {
	requireIntegration(t)

	t.Run("BAC-01 crea esquema roles y usuario de prueba", func(t *testing.T) {
		tables := queryDatabase(t, `
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;`)
		for _, expectedTable := range []string{"organizacions", "usuarios", "sesion_activas"} {
			if !containsLine(tables, expectedTable) {
				t.Errorf("falta la tabla %q; tablas encontradas:\n%s", expectedTable, tables)
			}
		}

		seed := queryDatabase(t, `
SELECT email_usuario || '|' || rol || '|' || activo::text || '|' || (contrasena_hash <> 'Admin123!')::text
FROM usuarios
WHERE email_usuario = 'admin@elcentinela.com';`)
		if seed != "admin@elcentinela.com|ADMIN|true|true" {
			t.Fatalf("usuario seed inesperado: %q", seed)
		}

		roleColumn := queryDatabase(t, `
SELECT data_type
FROM information_schema.columns
WHERE table_schema = 'public' AND table_name = 'usuarios' AND column_name = 'rol';`)
		if roleColumn == "" {
			t.Fatal("la tabla usuarios no contiene la columna rol")
		}
	})

	t.Run("BAC-02 persiste bcrypt y diferencia credenciales", func(t *testing.T) {
		hash := queryDatabase(t, `SELECT contrasena_hash FROM usuarios WHERE email_usuario = 'admin@elcentinela.com';`)
		if !strings.HasPrefix(hash, "$2") {
			t.Fatalf("la contraseña no parece un hash bcrypt: %q", hash)
		}
		if strings.Contains(hash, "Admin123!") {
			t.Fatal("la contraseña en texto plano aparece dentro del valor persistido")
		}

		status, _ := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{
			"email": "admin@elcentinela.com", "contrasena": "incorrecta",
		})
		if status != http.StatusUnauthorized {
			t.Fatalf("contraseña inválida: esperado 401, recibido %d", status)
		}
		_ = login(t, "Admin123!")
	})

	var temporaryJWT string
	t.Run("BAC-03 login PostgreSQL emite JWT firmado con identidad y rol", func(t *testing.T) {
		body := login(t, "Admin123!")
		temporaryJWT = requiredString(t, body, "jwtTemporal")
		if linked, ok := body["totpVinculado"].(bool); !ok || linked {
			t.Fatalf("totpVinculado inicial inesperado: %#v", body["totpVinculado"])
		}

		claims := decodeJWTClaims(t, temporaryJWT)
		for _, field := range []string{"sub", "jti", "rol", "org_id", "exp"} {
			if claims[field] == nil || claims[field] == "" {
				t.Errorf("claim JWT %q ausente: %#v", field, claims)
			}
		}
		if claims["rol"] != "ADMIN" || claims["tipo"] != "pre-auth" || claims["2fa_verificado"] != false {
			t.Fatalf("claims pre-auth inesperados: %#v", claims)
		}
		parsedToken, err := jwt.Parse(temporaryJWT, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				t.Fatalf("algoritmo JWT inesperado: %s", token.Method.Alg())
			}
			return []byte(testJWTSecret), nil
		})
		if err != nil || !parsedToken.Valid {
			t.Fatalf("la firma del JWT no es válida: %v", err)
		}
	})

	t.Run("BAC-04 unifica errores 400 y 401", func(t *testing.T) {
		cases := []struct {
			name       string
			payload    any
			wantStatus int
			wantCode   string
		}{
			{name: "payload incompleto", payload: map[string]any{"email": "admin@elcentinela.com"}, wantStatus: http.StatusBadRequest, wantCode: "INVALID_REQUEST"},
			{name: "credenciales inválidas", payload: map[string]any{"email": "admin@elcentinela.com", "contrasena": "incorrecta"}, wantStatus: http.StatusUnauthorized, wantCode: "AUTH_FAILED"},
		}
		for _, testCase := range cases {
			t.Run(testCase.name, func(t *testing.T) {
				status, body := requestJSON(t, http.MethodPost, "/auth/login", "", testCase.payload)
				if status != testCase.wantStatus {
					t.Fatalf("status esperado %d, recibido %d: %#v", testCase.wantStatus, status, body)
				}
				if body["errorCode"] != testCase.wantCode {
					t.Errorf("errorCode esperado %q, recibido %#v", testCase.wantCode, body["errorCode"])
				}
				if message, ok := body["message"].(string); !ok || strings.TrimSpace(message) == "" {
					t.Errorf("message ausente o vacío: %#v", body)
				}
			})
		}
	})

	var accessToken string
	t.Run("LOGIN-01 genera persiste y valida TOTP", func(t *testing.T) {
		status, qrBody := requestJSON(t, http.MethodGet, "/auth/2fa/qr", temporaryJWT, nil)
		if status != http.StatusOK {
			t.Fatalf("QR esperado 200, recibido %d: %#v", status, qrBody)
		}
		secret := requiredString(t, qrBody, "secretoManual")
		qrBase64 := requiredString(t, qrBody, "qrBase64")
		if !strings.HasPrefix(qrBase64, "data:image/png;base64,") {
			t.Fatalf("QR no está codificado como PNG data URL")
		}

		persistedSecret := queryDatabase(t, `SELECT secreto_totp_cifrado FROM usuarios WHERE email_usuario = 'admin@elcentinela.com';`)
		if persistedSecret == "" || persistedSecret == secret || strings.Contains(persistedSecret, secret) {
			t.Fatalf("el secreto TOTP no quedó cifrado correctamente")
		}
		if linked := queryDatabase(t, `SELECT totp_vinculado::text FROM usuarios WHERE email_usuario = 'admin@elcentinela.com';`); linked != "false" {
			t.Fatalf("TOTP debe permanecer sin vincular antes de confirmar, recibido %q", linked)
		}

		status, _ = requestJSON(t, http.MethodPost, "/auth/2fa/verify", temporaryJWT, map[string]any{"codigo": "123"})
		if status != http.StatusBadRequest {
			t.Fatalf("código con longitud inválida: esperado 400, recibido %d", status)
		}
		status, _ = requestJSON(t, http.MethodPost, "/auth/2fa/verify", temporaryJWT, map[string]any{"codigo": "000000"})
		if status != http.StatusUnauthorized {
			t.Fatalf("código incorrecto: esperado 401, recibido %d", status)
		}
		expiredCode, err := totp.GenerateCode(secret, time.Now().UTC().Add(-2*time.Minute))
		if err != nil {
			t.Fatalf("no se pudo generar un TOTP vencido: %v", err)
		}
		status, _ = requestJSON(t, http.MethodPost, "/auth/2fa/verify", temporaryJWT, map[string]any{"codigo": expiredCode})
		if status != http.StatusUnauthorized {
			t.Fatalf("código vencido: esperado 401, recibido %d", status)
		}

		code, err := totp.GenerateCode(secret, time.Now().UTC())
		if err != nil {
			t.Fatalf("no se pudo generar un TOTP válido: %v", err)
		}
		status, tokenBody := requestJSON(t, http.MethodPost, "/auth/2fa/verify", temporaryJWT, map[string]any{"codigo": code})
		if status != http.StatusOK {
			t.Fatalf("TOTP válido esperado 200, recibido %d: %#v", status, tokenBody)
		}
		accessToken = requiredString(t, tokenBody, "accessToken")
		_ = requiredString(t, tokenBody, "refreshToken")
		if linked := queryDatabase(t, `SELECT totp_vinculado::text FROM usuarios WHERE email_usuario = 'admin@elcentinela.com';`); linked != "true" {
			t.Fatalf("TOTP no quedó vinculado después de validarlo: %q", linked)
		}
	})

	t.Run("BAC-01 expone los roles ADMIN y OPERATOR", func(t *testing.T) {
		status, roles := requestJSONArray(t, http.MethodGet, "/roles", accessToken)
		if status != http.StatusOK {
			t.Fatalf("roles esperado 200, recibido %d: %#v", status, roles)
		}
		availableRoles := map[string]bool{}
		for _, item := range roles {
			role, ok := item.(map[string]any)
			if !ok {
				t.Fatalf("rol con formato inesperado: %#v", item)
			}
			if value, ok := role["valor"].(string); ok {
				availableRoles[value] = true
			}
		}
		for _, expectedRole := range []string{"ADMIN", "OPERATOR"} {
			if !availableRoles[expectedRole] {
				t.Errorf("falta el rol %s en %#v", expectedRole, roles)
			}
		}
	})

	t.Run("LOGIN-03 completa el recorrido backend y protege rutas", func(t *testing.T) {
		claims := decodeJWTClaims(t, accessToken)
		if claims["tipo"] != "access" || claims["2fa_verificado"] != true || claims["rol"] != "ADMIN" {
			t.Fatalf("claims de acceso inesperados: %#v", claims)
		}

		status, _ := requestJSON(t, http.MethodGet, "/account/profile", "", nil)
		if status != http.StatusUnauthorized {
			t.Fatalf("ruta protegida sin token: esperado 401, recibido %d", status)
		}
		status, _ = requestJSON(t, http.MethodGet, "/account/profile", temporaryJWT, nil)
		if status != http.StatusForbidden {
			t.Fatalf("ruta protegida con pre-auth: esperado 403, recibido %d", status)
		}
		status, body := requestJSON(t, http.MethodGet, "/account/profile", accessToken, nil)
		if status != http.StatusOK {
			t.Fatalf("ruta protegida con access token: esperado 200, recibido %d: %#v", status, body)
		}
	})
}

func containsLine(output, expected string) bool {
	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) == expected {
			return true
		}
	}
	return false
}
