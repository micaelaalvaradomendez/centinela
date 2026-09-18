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
		for _, expectedTable := range []string{"organizaciones", "usuarios", "sesiones_activas"} {
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

	t.Run("BAC-05 protege permisos_instancia (user_instances) contra duplicados", func(t *testing.T) {
		// actual.md pide clave primaria compuesta (user_id, instance_id); lo que esa clave
		// garantiza es que no pueda existir dos veces el mismo par (usuario, instancia).
		// Un id propio + una restricción de unicidad sobre ambas columnas cumple lo mismo,
		// así que aceptamos cualquiera de las dos formas, no solo la PK compuesta literal.
		uniqueColumns := queryDatabase(t, `
SELECT string_agg(a.attname, ',' ORDER BY a.attname)
FROM pg_index i
JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
WHERE i.indrelid = 'permisos_instancia'::regclass AND i.indisunique
GROUP BY i.indexrelid
HAVING string_agg(a.attname, ',' ORDER BY a.attname) = 'usuario_id,vmid_proxmox';`)
		if uniqueColumns != "usuario_id,vmid_proxmox" {
			t.Errorf("no hay ninguna restricción de unicidad (PK compuesta o UNIQUE) sobre (usuario_id, vmid_proxmox); la base permite pares duplicados")
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
			"email": "admin@elcentinela.com", "password": "incorrecta",
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
			{name: "credenciales inválidas", payload: map[string]any{"email": "admin@elcentinela.com", "password": "incorrecta"}, wantStatus: http.StatusUnauthorized, wantCode: "AUTH_FAILED"},
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
	var linkedSecret string
	var linkedCode string
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
		linkedSecret = secret
		linkedCode = code
		if linked := queryDatabase(t, `SELECT totp_vinculado::text FROM usuarios WHERE email_usuario = 'admin@elcentinela.com';`); linked != "true" {
			t.Fatalf("TOTP no quedó vinculado después de validarlo: %q", linked)
		}
	})

	t.Run("BAC-10 BAC-11 impiden reemplazar TOTP y validan el secreto persistido", func(t *testing.T) {
		status, relogin := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{
			"email": "admin@elcentinela.com", "password": "Admin123!",
		})
		if status != http.StatusOK {
			t.Fatalf("login posterior esperado 200, recibido %d: %#v", status, relogin)
		}
		secondJWT := requiredString(t, relogin, "jwtTemporal")
		if linked, ok := relogin["totpVinculado"].(bool); !ok || !linked {
			t.Fatalf("login posterior no informa TOTP vinculado: %#v", relogin)
		}

		// BAC-10: una cuenta ya vinculada no debe obtener otro secreto por el flujo inicial.
		qrStatus, _ := requestJSON(t, http.MethodGet, "/auth/2fa/qr", secondJWT, nil)
		if qrStatus == http.StatusOK {
			t.Fatal("una cuenta con TOTP vinculado pudo generar un QR nuevo sin pasar por reset/relink")
		}

		// BAC-11: el login posterior valida contra el secreto cifrado persistido.
		validCode, err := totp.GenerateCode(linkedSecret, time.Now().UTC())
		if err != nil {
			t.Fatalf("no se pudo generar código TOTP posterior: %v", err)
		}
		verifyStatus, tokens := requestJSON(t, http.MethodPost, "/auth/2fa/verify", secondJWT, map[string]any{"codigo": validCode})
		if verifyStatus != http.StatusOK {
			t.Fatalf("código válido posterior esperado 200, recibido %d: %#v", verifyStatus, tokens)
		}

		// BAC-11: el código usado en la vinculación no debe aceptarse como replay en otra sesión.
		replayLoginStatus, replayLogin := requestJSON(t, http.MethodPost, "/auth/login", "", map[string]any{
			"email": "admin@elcentinela.com", "password": "Admin123!",
		})
		if replayLoginStatus != http.StatusOK {
			t.Fatalf("login para probar replay esperado 200, recibido %d: %#v", replayLoginStatus, replayLogin)
		}
		replayJWT := requiredString(t, replayLogin, "jwtTemporal")
		replayStatus, _ := requestJSON(t, http.MethodPost, "/auth/2fa/verify", replayJWT, map[string]any{"codigo": linkedCode})
		if replayStatus == http.StatusOK {
			t.Fatal("el código TOTP usado en la vinculación fue aceptado nuevamente")
		}
	})

	t.Run("BAC-09 expone los roles ADMIN y OPERATOR", func(t *testing.T) {
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

	t.Run("BAC-05 BAC-06 BAC-06B gestiona usuarios y restringe OPERATOR", func(t *testing.T) {
		createdStatus, created := requestJSON(t, http.MethodPost, "/users", accessToken, map[string]any{
			"nombreCompleto": "Operador de aceptación",
			"nombreUsuario":  "operador_aceptacion",
			"emailUsuario":   "operador.aceptacion@elcentinela.com",
			"rol":            "OPERATOR",
		})
		if createdStatus != http.StatusCreated {
			t.Fatalf("crear usuario esperado 201, recibido %d: %#v", createdStatus, created)
		}
		createdID := requiredString(t, created, "id")
		if requiredString(t, created, "contrasenaTemp") == "" {
			t.Fatal("la creación debe devolver una contraseña temporal")
		}
		if created["rol"] != "OPERATOR" || created["activo"] != true {
			t.Fatalf("respuesta de creación inesperada: %#v", created)
		}
		if persisted := queryDatabase(t, `SELECT rol || '|' || activo::text || '|' || cambio_contrasena::text FROM usuarios WHERE nombre_usuario = 'operador_aceptacion';`); persisted != "OPERATOR|true|true" {
			t.Fatalf("estado inicial persistido inesperado: %q", persisted)
		}

		listStatus, listed := requestJSON(t, http.MethodGet, "/users?rol=OPERATOR&buscar=aceptacion", accessToken, nil)
		if listStatus != http.StatusOK {
			t.Fatalf("listar usuarios esperado 200, recibido %d: %#v", listStatus, listed)
		}
		users, ok := listed["users"].([]any)
		if !ok || len(users) != 1 {
			t.Fatalf("listado filtrado inesperado: %#v", listed["users"])
		}

		permissionsStatus, _ := requestValue(t, http.MethodPut, "/users/"+createdID+"/instances", accessToken, map[string]any{
			"vmids": []int{101, 102},
		})
		if permissionsStatus != http.StatusNoContent {
			t.Fatalf("asignar permisos esperado 204, recibido %d", permissionsStatus)
		}
		detailStatus, detail := requestJSON(t, http.MethodGet, "/users/"+createdID, accessToken, nil)
		if detailStatus != http.StatusOK {
			t.Fatalf("obtener usuario esperado 200, recibido %d: %#v", detailStatus, detail)
		}
		if permissions, ok := detail["instanciasPermitidas"].([]any); !ok || len(permissions) != 2 {
			t.Fatalf("permisos iniciales inesperados: %#v", detail["instanciasPermitidas"])
		}

		replaceStatus, _ := requestValue(t, http.MethodPut, "/users/"+createdID+"/instances", accessToken, map[string]any{
			"vmids": []int{102},
		})
		if replaceStatus != http.StatusNoContent {
			t.Fatalf("reemplazar permisos esperado 204, recibido %d", replaceStatus)
		}
		detailStatus, detail = requestJSON(t, http.MethodGet, "/users/"+createdID, accessToken, nil)
		if detailStatus != http.StatusOK {
			t.Fatalf("obtener usuario después del reemplazo esperado 200, recibido %d: %#v", detailStatus, detail)
		}
		permissions, ok := detail["instanciasPermitidas"].([]any)
		if !ok || len(permissions) != 1 || permissions[0] != float64(102) {
			t.Fatalf("reemplazo de permisos inesperado: %#v", detail["instanciasPermitidas"])
		}

		// Prueba funcional (no solo de esquema): si el cliente manda un VMID repetido,
		// ¿el sistema lo deduplica o persiste dos filas idénticas?
		duplicateStatus, _ := requestValue(t, http.MethodPut, "/users/"+createdID+"/instances", accessToken, map[string]any{
			"vmids": []int{201, 201},
		})
		if duplicateStatus != http.StatusNoContent {
			t.Fatalf("asignar vmids duplicados esperado 204, recibido %d", duplicateStatus)
		}
		storedRows := queryDatabase(t, `SELECT count(*) FROM permisos_instancia WHERE usuario_id = '`+createdID+`' AND vmid_proxmox = 201;`)
		if storedRows != "1" {
			t.Errorf("se enviaron 2 vmids repetidos y se persistieron %s filas para el mismo (usuario, instancia); se esperaba 1", storedRows)
		}

		operatorToken := signedAccessToken(t, createdID, "OPERATOR", decodeJWTClaims(t, accessToken)["org_id"].(string))
		for _, method := range []string{http.MethodGet, http.MethodPost} {
			payload := any(nil)
			if method == http.MethodPost {
				payload = map[string]any{
					"nombreCompleto": "No autorizado",
					"nombreUsuario":  "no_autorizado",
					"emailUsuario":   "no.autorizado@elcentinela.com",
					"rol":            "OPERATOR",
				}
			}
			status, _ := requestJSON(t, method, "/users", operatorToken, payload)
			if status != http.StatusForbidden {
				t.Fatalf("%s para OPERATOR: esperado 403, recibido %d", method, status)
			}
		}
		operatorPermissionStatus, _ := requestValue(t, http.MethodPut, "/users/"+createdID+"/instances", operatorToken, map[string]any{
			"vmids": []int{103},
		})
		if operatorPermissionStatus != http.StatusForbidden {
			t.Fatalf("asignar permisos para OPERATOR: esperado 403, recibido %d", operatorPermissionStatus)
		}

		// El hito "Gestión Administrativa de Usuarios" exige 403 en TODOS los endpoints de gestión,
		// no solo en listar/crear: also cubre consultar detalle, editar y desactivar.
		detailForbiddenStatus, _ := requestJSON(t, http.MethodGet, "/users/"+createdID, operatorToken, nil)
		if detailForbiddenStatus != http.StatusForbidden {
			t.Fatalf("consultar detalle con OPERATOR: esperado 403, recibido %d", detailForbiddenStatus)
		}
		updateForbiddenStatus, _ := requestJSON(t, http.MethodPut, "/users/"+createdID, operatorToken, map[string]any{"nombreCompleto": "Intento no autorizado"})
		if updateForbiddenStatus != http.StatusForbidden {
			t.Fatalf("editar usuario con OPERATOR: esperado 403, recibido %d", updateForbiddenStatus)
		}
		deleteForbiddenStatus, _ := requestValue(t, http.MethodDelete, "/users/"+createdID, operatorToken, nil)
		if deleteForbiddenStatus != http.StatusForbidden {
			t.Fatalf("desactivar usuario con OPERATOR: esperado 403, recibido %d", deleteForbiddenStatus)
		}

		newRole := "ADMIN"
		newEmail := "administrador.aceptacion@elcentinela.com"
		updateStatus, updated := requestJSON(t, http.MethodPut, "/users/"+createdID, accessToken, map[string]any{
			"nombreCompleto": "Administrador de aceptación",
			"emailUsuario":   newEmail,
			"rol":            newRole,
		})
		if updateStatus != http.StatusOK || updated["rol"] != newRole {
			t.Fatalf("actualizar usuario inesperado: %d %#v", updateStatus, updated)
		}
		if persistedEmail := queryDatabase(t, `SELECT email_usuario FROM usuarios WHERE id = '`+createdID+`';`); persistedEmail != newEmail {
			t.Fatalf("correo no se actualizó vía PUT: %q", persistedEmail)
		}

		// BAC-06B exige poder cambiar el estado (isActive) vía PUT, no solo vía DELETE.
		deactivateViaPutStatus, deactivated := requestJSON(t, http.MethodPut, "/users/"+createdID, accessToken, map[string]any{"activo": false})
		if deactivateViaPutStatus != http.StatusOK || deactivated["activo"] != false {
			t.Fatalf("desactivar vía PUT (activo=false) inesperado: %d %#v", deactivateViaPutStatus, deactivated)
		}
		reactivateStatus, reactivated := requestJSON(t, http.MethodPut, "/users/"+createdID, accessToken, map[string]any{"activo": true})
		if reactivateStatus != http.StatusOK || reactivated["activo"] != true {
			t.Fatalf("reactivar vía PUT (activo=true) inesperado: %d %#v", reactivateStatus, reactivated)
		}

		deleteStatus, _ := requestValue(t, http.MethodDelete, "/users/"+createdID, accessToken, nil)
		if deleteStatus != http.StatusNoContent {
			t.Fatalf("desactivar usuario esperado 204, recibido %d", deleteStatus)
		}
		if persisted := queryDatabase(t, `SELECT rol || '|' || activo::text FROM usuarios WHERE nombre_usuario = 'operador_aceptacion';`); persisted != "ADMIN|false" {
			t.Fatalf("estado final persistido inesperado: %q", persisted)
		}
	})

	t.Run("BAC-06 BAC-06B exigen el prefijo /api/admin/users documentado en actual.md", func(t *testing.T) {
		// actual.md describe GET/POST/DELETE /api/admin/users y PUT /api/admin/users/{id}.
		// En cmd/api/main.go el grupo se llama "admin" pero se monta con Group("/") en vez
		// de Group("/admin"), por lo que el CRUD queda expuesto en /api/users, sin el prefijo.
		status, _ := requestJSON(t, http.MethodGet, "/admin/users", accessToken, nil)
		if status != http.StatusOK {
			t.Errorf("GET /api/admin/users documentado por BAC-06: esperado 200, recibido %d (el grupo \"admin\" de cmd/api/main.go está montado en \"/\", no en \"/admin\")", status)
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

func signedAccessToken(t *testing.T, userID, role, orgID string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub":            userID,
		"rol":            role,
		"tipo":           "access",
		"2fa_verificado": true,
		"org_id":         orgID,
		"jti":            "acceptance-operator-session",
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("no se pudo firmar token de prueba: %v", err)
	}
	return signed
}
