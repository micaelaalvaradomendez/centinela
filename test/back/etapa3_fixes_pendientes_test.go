package back_test

import (
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Suite de aceptación de los FIX pendientes de documentacion/actual.md (FIX-123, FIX-128 vlan,
// FIX-133 a FIX-140). Metodología TDD: cada prueba sobre una tarea no implementada DEBE FALLAR
// (rojo) hasta que el código productivo exista. Son pruebas estáticas sobre el código fuente del
// submódulo backend (igual que etapa3_acceptance_test.go), salvo FIX-123 que es de infraestructura.

func backendSrc(t *testing.T, relative string) string {
	t.Helper()
	return string(readFile(t, sourcePath("backend", relative)))
}

// TestFIX123_PlantillasEISOsEnProxmoxReal verifica FIX-123 (INF-10) contra el Proxmox físico.
// Es una verificación de infraestructura: solo corre con PROXMOX_REAL_URL, PROXMOX_REAL_TOKEN_ID y
// PROXMOX_REAL_TOKEN_SECRET (y opcionalmente PROXMOX_REAL_NODE). Sin ellas se omite.
func TestFIX123_PlantillasEISOsEnProxmoxReal(t *testing.T) {
	base := os.Getenv("PROXMOX_REAL_URL")
	tokenID := os.Getenv("PROXMOX_REAL_TOKEN_ID")
	secret := os.Getenv("PROXMOX_REAL_TOKEN_SECRET")
	if base == "" || tokenID == "" || secret == "" {
		t.Skip("FIX-123: defina PROXMOX_REAL_URL, PROXMOX_REAL_TOKEN_ID y PROXMOX_REAL_TOKEN_SECRET para verificar el Proxmox físico")
	}
	node := envOrDefault("PROXMOX_REAL_NODE", "pve")
	client := &http.Client{Timeout: 15 * time.Second}

	for _, c := range []struct{ tipo, patron string }{
		{"vztmpl", `vztmpl/.+\.tar\.(zst|gz|xz)`},
		{"iso", `iso/.+\.iso`},
	} {
		req, _ := http.NewRequest(http.MethodGet, strings.TrimRight(base, "/")+"/nodes/"+node+"/storage/local/content?content="+c.tipo, nil)
		req.Header.Set("Authorization", "PVEAPIToken="+tokenID+"="+secret)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("FIX-123: no se pudo consultar storage local (%s): %v", c.tipo, err)
		}
		body := make([]byte, 1<<20)
		n, _ := resp.Body.Read(body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("FIX-123: GET /storage/local/content?content=%s respondió %d; el token necesita Datastore.Audit", c.tipo, resp.StatusCode)
			continue
		}
		if !regexp.MustCompile(c.patron).Match(body[:n]) {
			t.Errorf("FIX-123: local:%s no contiene ningún recurso verificable", c.tipo)
		}
	}
}

// TestFIX128VLAN_VlanCeroEsTraficoSinEtiquetar verifica FIX-128 (backend): vlan 0 o ausente es untagged.
func TestFIX128VLAN_VlanCeroEsTraficoSinEtiquetar(t *testing.T) {
	handler := backendSrc(t, "internal/adapters/primary/http/instance_handler.go")
	client := backendSrc(t, "internal/adapters/secondary/proxmox/client.go")

	t.Run("validarCreacionNormalizada admite vlan 0 convirtiéndola en nil", func(t *testing.T) {
		if !regexp.MustCompile(`\*input\.VLAN\s*==\s*0`).MatchString(handler) &&
			!regexp.MustCompile(`\*[a-zA-Z.]*VLAN\s*==\s*0`).MatchString(handler) {
			t.Errorf("FIX-128: instance_handler.go no trata vlan == 0 como untagged; hoy responde 400 'vlan debe estar entre 1 y 4094'")
		}
		if !regexp.MustCompile(`VLAN\s*=\s*nil`).MatchString(handler) {
			t.Errorf("FIX-128: instance_handler.go debe poner VLAN = nil cuando llega 0")
		}
	})

	t.Run("sigue rechazando vlans fuera de 1-4094", func(t *testing.T) {
		if !strings.Contains(handler, "4094") {
			t.Errorf("FIX-128: se perdió la validación del rango superior 4094")
		}
	})

	t.Run("client.go solo agrega tag= con VLAN efectiva", func(t *testing.T) {
		if !regexp.MustCompile(`if input\.VLAN != nil`).MatchString(client) {
			t.Errorf("FIX-128: client.go debe omitir ,tag= cuando VLAN es nil")
		}
	})
}

// TestFIX133_CuotaVCPUSinOvercommit verifica FIX-133 (RF-10 / BAC-40): vCPU sola no pasa por cuota de nodo.
func TestFIX133_CuotaVCPUSinOvercommit(t *testing.T) {
	cfg := backendSrc(t, "internal/core/services/configuracion_instancias.go")
	cuotas := backendSrc(t, "internal/core/services/cuotas_service.go")

	t.Run("la cuota de nodo solo se consulta cuando deltaMemoria > 0", func(t *testing.T) {
		if regexp.MustCompile(`if\s+deltaCores\s*>\s*0\s*\|\|\s*deltaMemoria\s*>\s*0`).MatchString(cfg) {
			t.Errorf("FIX-133: configuracion_instancias.go sigue delegando con deltaCores lineal")
		}
		filtraEnCfg := regexp.MustCompile(`if\s+deltaMemoria\s*>\s*0`).MatchString(cfg)
		filtraEnCuotas := regexp.MustCompile(`if\s+deltaMemoria\s*<=\s*0`).MatchString(cuotas)
		if !filtraEnCfg && !filtraEnCuotas {
			t.Errorf("FIX-133: ni configuracion_instancias.go ni cuotas_service.go omiten el bloqueo de margen de nodo cuando deltaMemoria <= 0")
		}
	})

	t.Run("AmpliarInstancia no rechaza por la suma lineal de vCPU del nodo", func(t *testing.T) {
		i := strings.Index(cuotas, "func (s *CuotasService) AmpliarInstancia")
		if i < 0 {
			t.Fatalf("FIX-133: no se encontró AmpliarInstancia")
		}
		cuerpo := cuotas[i:]
		if j := strings.Index(cuerpo, "\nfunc cantidadEnBytes"); j > 0 {
			cuerpo = cuerpo[:j]
		}
		if regexp.MustCompile(`deltaCores\s*>\s*coresDisponibles`).MatchString(cuerpo) {
			t.Errorf("FIX-133: AmpliarInstancia sigue rechazando con deltaCores > coresDisponibles (impide el CPU overcommit)")
		}
	})
}

// TestFIX134_InstanciasPausadas verifica FIX-134 (RF-03 / RF-04): paused no debe bloquear la operación.
func TestFIX134_InstanciasPausadas(t *testing.T) {
	handler := backendSrc(t, "internal/adapters/primary/http/instance_handler.go")
	client := backendSrc(t, "internal/adapters/secondary/proxmox/client.go")
	port := backendSrc(t, "internal/core/ports/proxmox_port.go")

	t.Run("stop (apagado forzado) admite running y paused", func(t *testing.T) {
		if !regexp.MustCompile(`paused`).MatchString(handler) {
			t.Fatalf("FIX-134: instance_handler.go no contempla el estado paused; stop, shutdown y delete dan 409 INSTANCE_INVALID_STATE")
		}
		if !regexp.MustCompile(`(?s)"stop".{0,200}paused|paused.{0,200}"stop"`).MatchString(handler) {
			t.Errorf("FIX-134: la acción stop debe aceptar paused además de running")
		}
	})

	t.Run("existe la acción resume en el handler", func(t *testing.T) {
		if !regexp.MustCompile(`"resume"`).MatchString(handler) {
			t.Errorf("FIX-134: CambiarEstado no admite la acción \"resume\"")
		}
	})

	t.Run("el cliente Proxmox implementa status/resume y el puerto lo declara", func(t *testing.T) {
		if !strings.Contains(client, "status/resume") && !strings.Contains(client, `"resume"`) {
			t.Errorf("FIX-134: client.go no envía POST .../status/resume a Proxmox")
		}
		if !regexp.MustCompile(`Resume\(`).MatchString(port) {
			t.Errorf("FIX-134: ports.ProxmoxPort no declara Resume")
		}
	})

	t.Run("eliminar una instancia pausada no queda en deadlock", func(t *testing.T) {
		if regexp.MustCompile(`instancia\.Estado\s*!=\s*"stopped"`).MatchString(handler) &&
			!regexp.MustCompile(`paused`).MatchString(handler) {
			t.Errorf("FIX-134: EliminarInstancia exige stopped y no ofrece salida para paused")
		}
	})
}

// TestFIX135_AlertasParadasExternas verifica FIX-135 (RF-11): parada directa en Proxmox no es instance_down.
func TestFIX135_AlertasParadasExternas(t *testing.T) {
	worker := backendSrc(t, "internal/core/services/alert_evaluator_worker.go")
	port := backendSrc(t, "internal/core/ports/proxmox_port.go")
	auditPort := backendSrc(t, "internal/core/ports/audit_port.go")
	main := backendSrc(t, "cmd/api/main.go")

	t.Run("explicadaPorCentinela consulta las tareas de Proxmox", func(t *testing.T) {
		if !regexp.MustCompile(`(?i)ListarTareas|TareasRecientes|HuboTareaExterna|tasks\?vmid`).MatchString(worker + port) {
			t.Errorf("FIX-135: no hay verificación secundaria contra GET /nodes/{node}/tasks?vmid=")
		}
		for _, tipo := range []string{"qmstop", "qmshutdown", "vzstop", "vzdump"} {
			if !strings.Contains(worker, tipo) {
				t.Errorf("FIX-135: el evaluador no reconoce la tarea externa %q", tipo)
			}
		}
	})

	t.Run("auditoría informativa SYSTEM_INSTANCE_STOPPED_EXTERNALLY", func(t *testing.T) {
		if !strings.Contains(worker+auditPort, "SYSTEM_INSTANCE_STOPPED_EXTERNALLY") {
			t.Errorf("FIX-135: falta la acción de auditoría SYSTEM_INSTANCE_STOPPED_EXTERNALLY")
		}
	})

	t.Run("variable ALERT_IGNORE_EXTERNAL_SHUTDOWN", func(t *testing.T) {
		if !strings.Contains(main+worker, "ALERT_IGNORE_EXTERNAL_SHUTDOWN") {
			t.Errorf("FIX-135: falta ALERT_IGNORE_EXTERNAL_SHUTDOWN")
		}
	})
}

// TestFIX136_UmbralesDeAlertaConfigurables verifica FIX-136 (RF-11 / RNF-07).
func TestFIX136_UmbralesDeAlertaConfigurables(t *testing.T) {
	worker := backendSrc(t, "internal/core/services/alert_evaluator_worker.go")
	main := backendSrc(t, "cmd/api/main.go")

	t.Run("existe ConfigAlertas inyectable", func(t *testing.T) {
		if !strings.Contains(worker, "type ConfigAlertas struct") {
			t.Fatalf("FIX-136: falta el struct ConfigAlertas")
		}
		for _, campo := range []string{"UmbralHostRAM", "UmbralHostCPU", "UmbralHostDisco", "UmbralInstanciaMem", "UmbralInstanciaCPU", "UmbralRecuperacion", "CiclosParaAlertar"} {
			if !strings.Contains(worker, campo) {
				t.Errorf("FIX-136: ConfigAlertas no tiene el campo %s", campo)
			}
		}
	})

	t.Run("main.go lee las variables de entorno", func(t *testing.T) {
		for _, env := range []string{"ALERT_THRESHOLD_HOST_RAM", "ALERT_THRESHOLD_HOST_CPU", "ALERT_THRESHOLD_HOST_DISK",
			"ALERT_THRESHOLD_VM_RAM", "ALERT_THRESHOLD_VM_CPU", "ALERT_THRESHOLD_RECOVERY", "ALERT_CYCLES_COUNT"} {
			if !strings.Contains(main, env) {
				t.Errorf("FIX-136: main.go no lee %s", env)
			}
		}
	})

	t.Run("los umbrales dejaron de ser constantes fijas del evaluador", func(t *testing.T) {
		if regexp.MustCompile(`umbralHostRAM\s*=\s*85\.0`).MatchString(worker) && !strings.Contains(worker, "ConfigAlertas") {
			t.Errorf("FIX-136: umbralHostRAM sigue hardcodeado")
		}
	})

	t.Run("se documentan en .env.example", func(t *testing.T) {
		env := backendSrc(t, ".env.example")
		if !strings.Contains(env, "ALERT_THRESHOLD_HOST_RAM") {
			t.Errorf("FIX-136: backend/.env.example no documenta ALERT_THRESHOLD_*")
		}
	})
}

// TestFIX137_CodigoInerteYAuditoriaDeSistema verifica FIX-137 (RF-12 / RNF-05).
func TestFIX137_CodigoInerteYAuditoriaDeSistema(t *testing.T) {
	handler := backendSrc(t, "internal/adapters/primary/http/instance_handler.go")
	audit := backendSrc(t, "internal/adapters/secondary/postgres/audit_repository.go")
	contrato := backendSrc(t, "internal/adapters/primary/http/contrato_etapa1.go")

	t.Run("se elimina el stub PausarInstancia no enrutado", func(t *testing.T) {
		if strings.Contains(handler, "func (h *InstanceHandler) PausarInstancia") {
			t.Errorf("FIX-137: sigue el stub PausarInstancia (ruta /pause no registrada)")
		}
		if strings.Contains(handler, "/instances/{instanceRef}/pause") {
			t.Errorf("FIX-137: Swagger sigue anunciando /instances/{instanceRef}/pause")
		}
	})

	t.Run("el listado de auditoría no oculta eventos de sistema sin organización", func(t *testing.T) {
		if regexp.MustCompile(`a\.organizacion_id = \? OR \(a\.organizacion_id IS NULL AND u\.organizacion_id = \?\)`).MatchString(audit) {
			t.Errorf("FIX-137: audit_repository.go conserva el filtro que oculta eventos con usuario_id y organizacion_id NULL")
		}
		if !regexp.MustCompile(`organizacion_id IS NULL`).MatchString(audit) {
			t.Errorf("FIX-137: la consulta debe incluir explícitamente los eventos de sistema (organizacion_id IS NULL)")
		}
	})

	t.Run("contrato_etapa1.go sin tipos huérfanos", func(t *testing.T) {
		for _, tipo := range regexp.MustCompile(`type\s+(\w+)\s+struct`).FindAllStringSubmatch(contrato, -1) {
			usos := 0
			for _, rel := range []string{"internal/adapters/primary/http/node_handler.go", "internal/adapters/primary/http/instance_handler.go",
				"internal/adapters/primary/http/snapshot_handler.go", "internal/adapters/primary/http/responses.go"} {
				usos += strings.Count(backendSrc(t, rel), tipo[1])
			}
			if usos == 0 {
				t.Errorf("FIX-137: el tipo %s de contrato_etapa1.go no se usa en ningún handler", tipo[1])
			}
		}
	})
}

// TestFIX138_CaracteresEspecialesNIST verifica FIX-138 (RF-01 / FIX-115A).
func TestFIX138_CaracteresEspecialesNIST(t *testing.T) {
	password := backendSrc(t, "internal/infrastructure/crypto/password.go")
	m := regexp.MustCompile("(?s)(?:CaracteresEspeciales|especiales)\\s*:?=\\s*(`[^`]*`|\"(?:[^\"\\\\]|\\\\.)*\")").FindStringSubmatch(password)
	if m == nil {
		t.Fatalf("FIX-138: no se encontró la definición de 'especiales' en password.go")
	}
	conjunto := m[1]
	for _, c := range []string{".", ",", ":", ";", "~", "(", ")", "[", "]", "/", "?"} {
		if !strings.Contains(conjunto, c) {
			t.Errorf("FIX-138: el conjunto de especiales %s rechaza el carácter %q (frente a gestores de contraseñas)", conjunto, c)
		}
	}
	// Los originales no pueden perderse.
	for _, c := range []string{"!", "@", "#", "$", "%", "^", "&", "*", "-", "_", "=", "+"} {
		if !strings.Contains(conjunto, c) {
			t.Errorf("FIX-138: se perdió el carácter especial original %q", c)
		}
	}
	t.Run("el mensaje de error refleja el conjunto ampliado o es genérico", func(t *testing.T) {
		if strings.Contains(password, "(!@#$%%^&*-_=+)") {
			t.Errorf("FIX-138: el mensaje de error sigue listando solo el conjunto restrictivo")
		}
	})
}

// TestFIX139_ReconciliacionDeTareasAlArranque verifica FIX-139 (RF-04 / RNF-04).
func TestFIX139_ReconciliacionDeTareasAlArranque(t *testing.T) {
	seg := backendSrc(t, "internal/core/services/seguimiento_tareas.go")
	main := backendSrc(t, "cmd/api/main.go")

	t.Run("existe RecuperarYReconciliarAlArranque", func(t *testing.T) {
		if !strings.Contains(seg, "RecuperarYReconciliarAlArranque") {
			t.Errorf("FIX-139: falta RecuperarYReconciliarAlArranque en seguimiento_tareas.go")
		}
	})

	t.Run("main.go la ejecuta antes de abrir el puerto HTTP", func(t *testing.T) {
		i := strings.Index(main, "RecuperarYReconciliarAlArranque")
		j := strings.Index(main, "ListenAndServe")
		if i < 0 {
			t.Fatalf("FIX-139: main.go no invoca la reconciliación inicial")
		}
		if j < 0 || i > j {
			t.Errorf("FIX-139: la reconciliación debe ejecutarse antes de srv.ListenAndServe()")
		}
	})

	t.Run("hay una prueba unitaria del reinicio con tareas RUNNING", func(t *testing.T) {
		if _, err := os.Stat(sourcePath("backend", "internal/core/services/seguimiento_tareas_test.go")); err != nil {
			t.Fatalf("FIX-139: falta seguimiento_tareas_test.go")
		}
		if !regexp.MustCompile(`(?i)RecuperarYReconciliar`).MatchString(backendSrc(t, "internal/core/services/seguimiento_tareas_test.go")) {
			t.Errorf("FIX-139: seguimiento_tareas_test.go no cubre RecuperarYReconciliarAlArranque")
		}
	})
}

// TestFIX140_ProveedorDeEmailVisible verifica FIX-140 (RF-09 / RF-11 / RF-13).
func TestFIX140_ProveedorDeEmailVisible(t *testing.T) {
	main := backendSrc(t, "cmd/api/main.go")

	t.Run("GET /api/version informa emailProvider", func(t *testing.T) {
		if !strings.Contains(main, "emailProvider") {
			t.Errorf("FIX-140: /api/version no informa el proveedor de email activo")
		}
	})

	t.Run("las respuestas que envían correo marcan X-Centinela-Email-Provider en modo mock", func(t *testing.T) {
		encontrado := false
		for _, rel := range []string{"internal/adapters/primary/http/auth_handler.go", "internal/adapters/primary/http/user_handler.go",
			"internal/adapters/primary/http/account_handler.go", "internal/adapters/primary/http/middleware/security_headers_middleware.go", "cmd/api/main.go"} {
			if strings.Contains(backendSrc(t, rel), "X-Centinela-Email-Provider") {
				encontrado = true
			}
		}
		if !encontrado {
			t.Errorf("FIX-140: ningún handler ni middleware emite la cabecera X-Centinela-Email-Provider")
		}
	})
}

// TestFIX143_VocabularioAuditoriaContratoCerrado verifica FIX-143 (RF-08 / BAC-44 / BAC-45):
// vocabulario de auditoria.accion y auditoria.resultado en un único conjunto cerrado,
// sin códigos duplicados entre caminos, publicado como enum en Swagger y en docs/contrato-etapa1.md,
// y validado por docs/contrato_etapa1_test.go.
func TestFIX143_VocabularioAuditoriaContratoCerrado(t *testing.T) {
	auditPort := backendSrc(t, "internal/core/ports/audit_port.go")
	auditHandler := backendSrc(t, "internal/adapters/primary/http/audit_handler.go")
	swaggerRaw := backendSrc(t, "docs/swagger.json")
	contratoDoc := backendSrc(t, "docs/contrato-etapa1.md")
	contratoTest := backendSrc(t, "docs/contrato_etapa1_test.go")

	reResultados := regexp.MustCompile(`Resultado\w+\s*(?:\w+)?\s*=\s*"([^"]+)"`)
	matchesRes := reResultados.FindAllStringSubmatch(auditPort, -1)
	if len(matchesRes) == 0 {
		t.Fatalf("FIX-143: no se encontraron constantes Resultado* en internal/core/ports/audit_port.go")
	}
	resultadosMap := map[string]bool{}
	var resultadosList []string
	for _, m := range matchesRes {
		if !resultadosMap[m[1]] {
			resultadosMap[m[1]] = true
			resultadosList = append(resultadosList, m[1])
		}
	}

	reAcciones := regexp.MustCompile(`Accion\w+\s*(?:\w+)?\s*=\s*"([^"]+)"`)
	matchesAcc := reAcciones.FindAllStringSubmatch(auditPort, -1)
	if len(matchesAcc) == 0 {
		t.Fatalf("FIX-143: no se encontraron constantes Accion* en internal/core/ports/audit_port.go")
	}
	accionesMap := map[string]bool{}
	var accionesList []string
	for _, m := range matchesAcc {
		if !accionesMap[m[1]] {
			accionesMap[m[1]] = true
			accionesList = append(accionesList, m[1])
		}
	}

	t.Run("resultado tiene un único conjunto cerrado sin familias duplicadas (EXITO/FALLA vs SUCCESS/FAILED)", func(t *testing.T) {
		tieneEspanol := resultadosMap["EXITO"] || resultadosMap["FALLA"]
		tieneIngles := resultadosMap["SUCCESS"] || resultadosMap["FAILED"]
		if tieneEspanol && tieneIngles {
			t.Errorf("FIX-143: ports/audit_port.go define dos familias duplicadas para auditoria.resultado (%v); debe existir un único conjunto cerrado", resultadosList)
		}
		if strings.Contains(auditHandler, `"SUCCESS"`) && strings.Contains(auditHandler, `"EXITO"`) {
			t.Errorf("FIX-143: validarConsultaAuditoria en audit_handler.go acepta ambas familias simultáneamente (SUCCESS/FAILED y EXITO/FALLA) en lugar de un único enum cerrado")
		}
	})

	t.Run("accion no contiene códigos muertos o duplicados por renombre en ports/audit_port.go", func(t *testing.T) {
		if accionesMap["CREATE_INSTANCE"] && (accionesMap["CREAR_VM"] || accionesMap["CREAR_LXC"]) {
			t.Errorf("FIX-143: ports/audit_port.go mantiene constantes duplicadas de creación (CREATE_INSTANCE junto con CREAR_VM/CREAR_LXC)")
		}
		// Toda constante Accion* declarada en audit_port.go debe tener al menos un escritor o pertenecer al enum cerrado documentado
		archivosEscritores := []string{
			"internal/core/services/audit_service.go",
			"internal/core/services/auth_service.go",
			"internal/core/services/user_service.go",
			"internal/core/services/email_change.go",
			"internal/core/services/eventos_service.go",
			"internal/core/services/configuracion_instancias.go",
			"internal/core/services/snapshot_service.go",
			"internal/core/services/seguimiento_tareas.go",
			"internal/core/services/alert_evaluator_worker.go",
			"internal/adapters/primary/http/instance_handler.go",
			"internal/adapters/primary/http/instance_config_handler.go",
			"internal/adapters/primary/http/snapshot_handler.go",
			"internal/adapters/primary/http/events_handler.go",
			"internal/adapters/primary/http/account_handler.go",
			"internal/adapters/primary/http/user_handler.go",
			"internal/adapters/primary/http/auth_handler.go",
			"internal/adapters/secondary/postgres/instance_repository.go",
			"internal/adapters/secondary/postgres/tarea_repository.go",
			"internal/adapters/secondary/postgres/user_repository.go",
		}
		codigoEscritores := ""
		for _, rel := range archivosEscritores {
			codigoEscritores += backendSrc(t, rel) + "\n"
		}
		reConstNames := regexp.MustCompile(`(Accion\w+)\s*(?:\w+)?\s*=\s*"([^"]+)"`)
		for _, m := range reConstNames.FindAllStringSubmatch(auditPort, -1) {
			constName, constVal := m[1], m[2]
			if !strings.Contains(codigoEscritores, constName) {
				t.Errorf("FIX-143: la constante ports.%s (%q) no es emitida por ningún servicio/worker activo (código muerto o duplicado en el vocabulario)", constName, constVal)
			}
		}
	})

	t.Run("el enum de auditoria.accion y auditoria.resultado se publica en Swagger y se documenta en contrato-etapa1.md", func(t *testing.T) {
		// Verificar que swagger.json declare enum para resultado y accion en el DTO o parámetros de /admin/audit
		if !regexp.MustCompile(`"resultado"\s*:\s*\{[^}]*"enum"\s*:`).MatchString(swaggerRaw) {
			t.Errorf("FIX-143: docs/swagger.json no publica el campo 'resultado' de auditoría con un 'enum' cerrado")
		}
		if !regexp.MustCompile(`"accion"\s*:\s*\{[^}]*"enum"\s*:`).MatchString(swaggerRaw) {
			t.Errorf("FIX-143: docs/swagger.json no publica el campo 'accion' de auditoría con un 'enum' cerrado")
		}
		for _, acc := range accionesList {
			if !strings.Contains(contratoDoc, acc) {
				t.Errorf("FIX-143: la acción %q de ports/audit_port.go no está documentada en docs/contrato-etapa1.md", acc)
			}
		}
		for _, res := range resultadosList {
			if !strings.Contains(contratoDoc, res) {
				t.Errorf("FIX-143: el resultado %q de ports/audit_port.go no está documentado en docs/contrato-etapa1.md", res)
			}
		}
	})

	t.Run("docs/contrato_etapa1_test.go afirma el conjunto cerrado de auditoría", func(t *testing.T) {
		if !strings.Contains(contratoTest, "/admin/audit") && !strings.Contains(contratoTest, "AuditoriaDTO") && !strings.Contains(contratoTest, "ItemAuditoriaResponse") {
			t.Errorf("FIX-143: docs/contrato_etapa1_test.go no valida el enum de auditoria.accion ni auditoria.resultado")
		}
	})
}

