package back_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

// TestFixesPendientes verifica las tareas y correcciones identificadas en documentacion/actual.md
// que requieren pruebas de aceptación y regresión dedicadas:
//   - FIX-50: Apagado correcto del servidor HTTP ante SIGINT/SIGTERM (BAC-18B).
//   - FIX-51: Historial de auditoria previo a la migración particionada accesible desde la API (BAC-18B).
//   - FIX-52: Distinguir INSTANCE_BUSY de INSTANCE_INVALID_STATE y limpiar código muerto (BAC-24A).
//   - FIX-53: Unificar códigos de acción y resultado en la auditoría de instancias (BAC-24A).
//   - FIX-54: Separar eliminación lógica y suspensión; migrar unicidad del correo (Backend).
//   - FIX-55: Alinear baja, validación y autenticación con identidad histórica (Backend).
func TestFixesPendientes(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)

	t.Run("FIX-50 el servidor HTTP implementa apagado ordenado con http.Server y Shutdown ante senales", func(t *testing.T) {
		mainContent := string(readFile(t, sourcePath("backend", "cmd/api/main.go")))

		// 1. router.Run no debe usarse porque bloquea sin recibir el contexto de apagado.
		if regexp.MustCompile(`router\.Run\(`).MatchString(mainContent) {
			t.Errorf("cmd/api/main.go usa router.Run(), lo que bloquea sin observar el contexto de apagado")
		}

		// 2. Debe instanciarse un http.Server explícito con Addr y Handler.
		if !regexp.MustCompile(`&?http\.Server\s*\{[^}]*Handler:\s*router`).MatchString(mainContent) &&
			!regexp.MustCompile(`&?http\.Server\s*\{[^}]*Addr:`).MatchString(mainContent) {
			t.Errorf("cmd/api/main.go no crea una instancia explícita de http.Server{Addr, Handler: router}")
		}

		// 3. Debe invocarse Shutdown(ctx) ante la cancelación del contexto de señales.
		if !regexp.MustCompile(`(?i)\.Shutdown\(`).MatchString(mainContent) {
			t.Errorf("cmd/api/main.go no invoca srv.Shutdown(...) al recibir la señal de terminación")
		}

		// 4. Debe escucharse <-ctx.Done() conectado a NotifyContext o signal.Notify.
		if !regexp.MustCompile(`<-ctx\.Done\(\)`).MatchString(mainContent) {
			t.Errorf("cmd/api/main.go no espera <-ctx.Done() para iniciar el apagado ordenado")
		}
	})

	t.Run("FIX-51 el historial de auditoria previo a la migracion particionada es visible por la API", func(t *testing.T) {
		// Entregable: los registros de auditoria anteriores al particionamiento (conservados en
		// auditoria_legacy o migrados) deben seguir siendo visibles en GET /api/admin/audit y exportaciones.
		dbCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/postgres/db.go")))
		repoCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/postgres/audit_repository.go")))

		// Verificación estática: db.go debe migrar los datos a la tabla particionada (INSERT INTO auditoria SELECT ... FROM auditoria_legacy)
		// o el repositorio de auditoría debe incluir auditoria_legacy con UNION ALL.
		migratesLegacy := regexp.MustCompile(`(?i)INSERT\s+INTO\s+auditoria[^(]*SELECT\s+.*FROM\s+auditoria_legacy`).MatchString(dbCode)
		queriesLegacy := strings.Contains(repoCode, "auditoria_legacy")

		if !migratesLegacy && !queriesLegacy {
			t.Errorf("FIX-51 no implementada: db.go no migra datos de auditoria_legacy a auditoria particionada, ni audit_repository.go consulta auditoria_legacy con UNION ALL")
		}

		// Verificación de base de datos si existe auditoria_legacy:
		hasLegacyTable := queryDatabase(t, "SELECT count(*) FROM information_schema.tables WHERE table_name = 'auditoria_legacy';")
		if hasLegacyTable == "1" {
			// Si la tabla legacy existe con filas, verificar que no queden huérfanas.
			legacyCount := queryDatabase(t, "SELECT count(*) FROM auditoria_legacy;")
			if legacyCount != "0" && !queriesLegacy {
				t.Errorf("existen %s filas en auditoria_legacy que no son consultadas por la API", legacyCount)
			}
		}
	})

	t.Run("FIX-52 distingue INSTANCE_BUSY de INSTANCE_INVALID_STATE y limpia codigo muerto", func(t *testing.T) {
		// 1. Limpieza de código muerto: ProxmoxPort.ReiniciarInstancia huérfano.
		portsCode := string(readFile(t, sourcePath("backend", "internal/core/ports/proxmox_port.go")))
		clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))

		hasReiniciarInPort := regexp.MustCompile(`(?m)^\s*ReiniciarInstancia\(`).MatchString(portsCode)
		if hasReiniciarInPort {
			// Si existe en el puerto, verificar que se use de forma consistente o se haya removido a favor de Reboot.
			handlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/instance_handler.go")))
			if !strings.Contains(handlerCode, "ReiniciarInstancia") {
				t.Errorf("ProxmoxPort.ReiniciarInstancia está declarado pero no se invoca en los handlers (código muerto)")
			}
		}

		// 2. Comentario desactualizado en EliminarInstancia.
		if strings.Contains(clientCode, "si está encendida Proxmox responde 500 con un mensaje de bloqueo") {
			t.Errorf("el comentario de EliminarInstancia en client.go sigue afirmando que Proxmox responde 500 por instancia encendida; esa rama fue reemplazada por la validación 409 del handler")
		}

		// 3. Comportamiento en vivo: acción rechazada por tarea RUNNING en curso debe dar 409 INSTANCE_BUSY,
		// mientras que acción rechazada por estado incompatible sin tarea en curso debe dar 409 INSTANCE_INVALID_STATE.
		// En el stub, 101 está running y sin tareas en curso: start da INSTANCE_INVALID_STATE.
		status, body := requestJSON(t, http.MethodPost, "/instances/101/status/start", adminToken, nil)
		if status != http.StatusConflict || body["errorCode"] != "INSTANCE_INVALID_STATE" {
			t.Errorf("start sobre instancia 101 (running, sin tarea en curso): esperado 409 INSTANCE_INVALID_STATE, recibido %d: %#v", status, body)
		}

		// Simulamos una tarea RUNNING en tareas_asincronas para la 102 (detenida).
		tareaUUID := queryDatabase(t, "SELECT uuid_generate_v7();")
		queryDatabase(t, fmt.Sprintf(
			"INSERT INTO tareas_asincronas (id, usuario_id, instancia_id, accion, upid_proxmox, estado, fecha_creacion) VALUES ('%s', '%s', '102', 'start', 'UPID:pve:dummy:0:0:vzstart:102:root@pam:', 'RUNNING', NOW());",
			tareaUUID, adminID,
		))
		t.Cleanup(func() {
			queryDatabase(t, fmt.Sprintf("DELETE FROM tareas_asincronas WHERE id = '%s';", tareaUUID))
		})

		// Con tarea en curso sobre la 102, intentar otra acción (stop) debe responder 409 INSTANCE_BUSY.
		busyStatus, busyBody := requestJSON(t, http.MethodPost, "/instances/102/status/stop", adminToken, nil)
		if busyStatus != http.StatusConflict || busyBody["errorCode"] != "INSTANCE_BUSY" {
			t.Errorf("acción sobre instancia con tarea RUNNING en curso: esperado 409 INSTANCE_BUSY, recibido %d: %#v", busyStatus, busyBody)
		}
	})

	t.Run("FIX-53 unificacion de codigos de accion y resultado en la auditoria de instancias", func(t *testing.T) {
		auditPortCode := string(readFile(t, sourcePath("backend", "internal/core/ports/audit_port.go")))
		handlerCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/instance_handler.go")))

		// 1. Constante ResultadoPendiente en audit_port.go.
		if !regexp.MustCompile(`ResultadoPendiente\s*=\s*"[^"]+"`).MatchString(auditPortCode) {
			t.Errorf("internal/core/ports/audit_port.go no define la constante ResultadoPendiente")
		}

		// 2. Literal "PENDING" suelto no debe usarse en instance_handler.go; debe usarse ports.ResultadoPendiente.
		pendingLiterals := regexp.MustCompile(`"PENDING"`).FindAllStringIndex(handlerCode, -1)
		// Si se usa ports.ResultadoPendiente, los literales no deben aparecer sueltos en llamadas de auditoría.
		if len(pendingLiterals) > 1 && !strings.Contains(handlerCode, "ports.ResultadoPendiente") {
			t.Errorf("instance_handler.go contiene literales sueltos \"PENDING\" en lugar de usar ports.ResultadoPendiente")
		}

		// 3. IniciarInstancia y DetenerInstancia deben auditar el nombre y tipo de la instancia, no vacíos ni fijos.
		if strings.Contains(handlerCode, `instanciaNombre: ""`) || strings.Contains(handlerCode, `"vm_or_lxc"`) {
			t.Errorf("IniciarInstancia o DetenerInstancia auditan instanciaNombre vacío o resource_type \"vm_or_lxc\"")
		}

		// 4. Códigos de acción unificados entre alias (/start, /stop) y genérico (/status/:action).
		// Ambas rutas deben registrar el mismo código de acción en auditoría.
		status, body := requestJSON(t, http.MethodPost, "/instances/101/status/stop", adminToken, nil)
		if status == http.StatusAccepted {
			upid := strings.ReplaceAll(requiredString(t, body, "upid"), "'", "''")
			accionAuditoria := queryDatabase(t, fmt.Sprintf("SELECT accion FROM auditoria WHERE detalles::text LIKE '%%%s%%' LIMIT 1;", upid))
			if accionAuditoria != "STOP" && accionAuditoria != "DETENER_VM" {
				t.Errorf("accion de auditoría inesperada para stop: %q", accionAuditoria)
			}
		}
	})

	t.Run("FIX-54 separar eliminacion logica y suspension con migracion de unicidad de correo", func(t *testing.T) {
		// 1. La tabla usuarios debe incluir la columna eliminado_en timestamptz nullable.
		columnaEliminado := queryDatabase(t, "SELECT count(*) FROM information_schema.columns WHERE table_name = 'usuarios' AND column_name = 'eliminado_en';")
		if columnaEliminado != "1" {
			t.Errorf("la tabla usuarios no tiene la columna eliminado_en (timestamptz nullable)")
		}

		// 2. El índice único de correo debe ser parcial (WHERE eliminado_en IS NULL).
		indiceParcial := queryDatabase(t, "SELECT coalesce(string_agg(indexdef, E'\\n'), '') FROM pg_indexes WHERE tablename = 'usuarios' AND indexdef ILIKE '%email_usuario%' AND indexdef ILIKE '%WHERE%eliminado_en%IS NULL%';")
		if indiceParcial == "" {
			t.Errorf("falta un índice único parcial en usuarios(email_usuario) con WHERE eliminado_en IS NULL")
		}

		// 3. No debe haber un índice único incondicional en email_usuario que impida reutilizar correos de eliminados.
		indiceIncondicional := queryDatabase(t, "SELECT coalesce(string_agg(indexdef, E'\\n'), '') FROM pg_indexes WHERE tablename = 'usuarios' AND indexdef ILIKE '%UNIQUE%' AND indexdef ILIKE '%email_usuario%' AND indexdef NOT ILIKE '%WHERE%';")
		if indiceIncondicional != "" {
			t.Errorf("persiste un índice único incondicional sobre email_usuario: %s", indiceIncondicional)
		}
	})

	t.Run("FIX-55 alinear baja, validacion y autenticacion con identidad historica", func(t *testing.T) {
		columnaEliminado := queryDatabase(t, "SELECT count(*) FROM information_schema.columns WHERE table_name = 'usuarios' AND column_name = 'eliminado_en';")
		if columnaEliminado != "1" {
			t.Skip("eliminado_en no está implementado aún en la base de datos (FIX-54)")
		}

		suffix := time.Now().UnixNano()
		email := fmt.Sprintf("historico.%d@elcentinela.com", suffix)
		username := fmt.Sprintf("hist_%d", suffix)

		// 1. Crear usuario inicial.
		userID, tempPass := createUserThroughAPI(t, adminToken, username, email, "OPERATOR")
		userSession := loginWithTOTP(t, email, tempPass, "")

		// 2. DELETE /api/admin/users/:id debe marcarlo como eliminado (eliminado_en != null y activo = false).
		delStatus, _ := requestValue(t, http.MethodDelete, "/admin/users/"+userID, adminToken, nil)
		if delStatus != http.StatusNoContent && delStatus != http.StatusOK {
			t.Fatalf("DELETE /admin/users/%s: esperado 204 o 200, recibido %d", userID, delStatus)
		}

		eliminadoEn := queryDatabase(t, fmt.Sprintf("SELECT coalesce(eliminado_en::text, '') FROM usuarios WHERE id = '%s';", userID))
		activo := queryDatabase(t, fmt.Sprintf("SELECT activo FROM usuarios WHERE id = '%s';", userID))
		if eliminadoEn == "" || activo != "f" {
			t.Errorf("tras DELETE: esperado eliminado_en con valor y activo = false, recibido eliminado_en=%q, activo=%s", eliminadoEn, activo)
		}

		// 3. Las credenciales previas quedan completamente invalidadas: login rechazado.
		loginStatus, _ := requestValue(t, http.MethodPost, "/auth/login", "", map[string]any{"email": email, "password": tempPass})
		if loginStatus != http.StatusUnauthorized && loginStatus != http.StatusForbidden {
			t.Errorf("login con credenciales de usuario eliminado debe ser rechazado (401 o 403), recibido %d", loginStatus)
		}

		// 4. Token previo revocado.
		tokenStatus, _ := requestValue(t, http.MethodGet, "/account/profile", userSession.AccessToken, nil)
		if tokenStatus != http.StatusUnauthorized {
			t.Errorf("token de sesión de usuario eliminado debe ser rechazado con 401, recibido %d", tokenStatus)
		}

		// 5. No se puede reactivar un usuario eliminado con PUT activo=true.
		reactivateStatus, _ := requestValue(t, http.MethodPut, "/admin/users/"+userID, adminToken, map[string]any{"activo": true})
		if reactivateStatus == http.StatusOK || reactivateStatus == http.StatusNoContent {
			t.Errorf("PUT /admin/users/:id con activo=true sobre usuario eliminado no debe permitirse, recibido %d", reactivateStatus)
		}

		// 6. Crear una cuenta nueva con el MISMO correo debe ser permitido y tener un UUID diferente.
		newUserID, _ := createUserThroughAPI(t, adminToken, fmt.Sprintf("nuevo_%d", suffix), email, "OPERATOR")
		if newUserID == userID {
			t.Errorf("la nueva cuenta creada con el mismo correo debe recibir un nuevo UUID, recibido el mismo: %s", newUserID)
		}

		// 7. El historial y auditoría previos siguen asociados al UUID original, sin fusionar identidades.
		auditRowsOld := queryDatabase(t, fmt.Sprintf("SELECT count(*) FROM auditoria WHERE usuario_id = '%s';", userID))
		if auditRowsOld == "0" {
			t.Errorf("los registros de auditoría del usuario original (%s) se perdieron o fueron reasignados", userID)
		}
	})

	t.Run("FIX-68 integracion del cliente de eventos con la ruta real servida por el backend", func(t *testing.T) {
		clientCode := string(readFile(t, sourcePath("frontend", "centinela/src/services/eventsClient.ts")))
		// 1. Contrato del frontend: EVENTS_STREAM_PATH debe ser '/events' y no '/events/stream'.
		if strings.Contains(clientCode, "EVENTS_STREAM_PATH = '/events/stream'") {
			t.Errorf("frontend/centinela/src/services/eventsClient.ts define EVENTS_STREAM_PATH = '/events/stream', ruta inexistente en el backend")
		}
		if !strings.Contains(clientCode, "EVENTS_STREAM_PATH = '/events'") {
			t.Errorf("frontend/centinela/src/services/eventsClient.ts no define EVENTS_STREAM_PATH = '/events'")
		}

		// 2. Comprobación en vivo del backend: la ruta /api/events responde con SSE para un ticket válido,
		// mientras que /api/events/stream responde 404 (confirmando que una desalineación rompería el cliente).
		statusWrong, _ := requestValue(t, http.MethodGet, "/events/stream", "", nil)
		if statusWrong != http.StatusNotFound {
			t.Errorf("GET /api/events/stream debe responder 404 Not Found, recibido %d", statusWrong)
		}

		// 3. Obtener ticket efímero con token de admin y verificar conexión a /events?ticket=...
		ticketStatus, ticketBody := requestJSON(t, http.MethodPost, "/events/ticket", adminToken, nil)
		if ticketStatus != http.StatusOK {
			t.Fatalf("POST /api/events/ticket falló con status %d: %#v", ticketStatus, ticketBody)
		}
		ticket := requiredString(t, ticketBody, "ticket")

		req, err := http.NewRequest(http.MethodGet, apiURL+"/events?ticket="+ticket, nil)
		if err != nil {
			t.Fatalf("fallo al crear request SSE: %v", err)
		}
		req.Header.Set("Accept", "text/event-stream")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("fallo al conectar a /api/events: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /api/events?ticket=...: esperado 200 OK, recibido %d", resp.StatusCode)
		}
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/event-stream") {
			t.Errorf("Content-Type de /api/events debe ser text/event-stream, recibido %q", contentType)
		}
	})
}

