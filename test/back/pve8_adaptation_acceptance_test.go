package back_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Suite de Aceptación para la Etapa de Adaptación a Proxmox VE 8.4.21 (documentacion/etapa-adaptacion-proxmox-8.md).
// Metodología TDD estricta: Todo test sobre una tarea pendiente DEBE FALLAR (Rojo) hasta que el equipo
// desarrolle e implemente el código productivo correspondiente.

// TestADP01_ParametrosAprovisionamientoQEMU_LXC verifica ADP-01:
// En Proxmox VE 8.4.21, client.go debe forzar ostype=l26 y agent=1 en QEMU,
// y unprivileged=1 con features: nesting=1 en LXC.
func TestADP01_ParametrosAprovisionamientoQEMU_LXC(t *testing.T) {
	clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))

	t.Run("QEMU inyecta ostype=l26 y agent=1 para canal virtio-serial en PVE 8.4", func(t *testing.T) {
		// Acepta literal de mapa ("ostype": "l26") y params.Set("ostype", "l26").
		if !regexp.MustCompile(`"ostype"(:\s*|,\s*)(\{)?"l26"`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe enviar ostype=l26 para instancias QEMU")
		}
		if !regexp.MustCompile(`"agent"(:\s*|,\s*)(\{)?"?1\b`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe enviar agent=1 para habilitar el socket virtio-serial en QEMU")
		}
	})

	t.Run("LXC inyecta unprivileged=1 y nesting=1 para compatibilidad systemd en PVE 8.4", func(t *testing.T) {
		if !regexp.MustCompile(`"unprivileged"(:\s*|,\s*)(\{)?"?1\b`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe enviar unprivileged=1 para aislamiento en contenedores LXC")
		}
		if !regexp.MustCompile(`nesting=1`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe configurar features nesting=1 para el arranque limpio de systemd en LXC")
		}
	})
}

// TestADP02_ParametrosPurgaBorrado verifica ADP-02:
// EliminarInstancia en client.go debe incorporar purge=1 y destroy-unreferenced-disks=1
// para evitar el error 500 en instancias con configuración HA y no dejar discos huérfanos.
func TestADP02_ParametrosPurgaBorrado(t *testing.T) {
	clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))

	t.Run("EliminarInstancia incluye query params de purga profunda", func(t *testing.T) {
		// Acepta "purge=1" en la URL o url.Values{"purge": {"1"}}.
		if !regexp.MustCompile(`purge=1|"purge":\s*\{"1"\}`).MatchString(clientCode) {
			t.Errorf("ADP-02: EliminarInstancia en client.go debe enviar purge=1 en la petición DELETE a Proxmox")
		}
		if !regexp.MustCompile(`destroy-unreferenced-disks=1|"destroy-unreferenced-disks":\s*\{"1"\}`).MatchString(clientCode) {
			t.Errorf("ADP-02: EliminarInstancia en client.go debe enviar destroy-unreferenced-disks=1 para evitar discos huérfanos")
		}
	})
}

// TestADP03_ParserDefensivoMemoryPrefix verifica ADP-03:
// Soporte agnóstico de tipos en enteroFlexible / prefijoFlexible y pruebas unitarias en client_test.go.
func TestADP03_ParserDefensivoMemoryPrefix(t *testing.T) {
	t.Run("client.go y client_test.go soportan memory/prefix mixtos e ip-addresses nulos o CIDR en PVE 8.4", func(t *testing.T) {
		clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))
		testCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client_test.go")))
		if !strings.Contains(clientCode, "enteroFlexible") || !strings.Contains(clientCode, "prefijoFlexible") {
			t.Errorf("ADP-03: client.go debe definir enteroFlexible y prefijoFlexible para tolerar tipos mixtos en PVE 8.4")
		}
		if !strings.Contains(testCode, "TestClient_ObtenerInterfaces_LXC") || !strings.Contains(testCode, "TestObtenerConfig_QemuYLXC") {
			t.Errorf("ADP-03: client_test.go debe incluir TestClient_ObtenerInterfaces_LXC y TestObtenerConfig_QemuYLXC")
		}
	})
}

// TestADP04_ParametrizacionVMIDsProtegidos verifica ADP-04:
// ParseVmidsProtegidos debe soportar listas vacías o personalizadas vía PROXMOX_PROTECTED_VMIDS.
func TestADP04_ParametrizacionVMIDsProtegidos(t *testing.T) {
	code := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/middleware/protected_instance.go")))
	testCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/middleware/protected_instance_test.go")))

	t.Run("ParseVmidsProtegidos soporta listas vacías y personalizadas desde entorno", func(t *testing.T) {
		if !strings.Contains(code, "ParseVmidsProtegidos") || !strings.Contains(code, `strings.Split(raw, ",")`) {
			t.Errorf("ADP-04: protected_instance.go debe implementar ParseVmidsProtegidos sobre listas separadas por comas")
		}
		if !strings.Contains(testCode, "TestRejectProtectedInstance_ListaDesdeEntorno") {
			t.Errorf("ADP-04: protected_instance_test.go debe incluir TestRejectProtectedInstance_ListaDesdeEntorno")
		}
	})
}

// TestADP05_ModoDualSimulador verifica ADP-05:
// cmd/proxmox-simulador soporta PROXMOX_SIM_TARGET_VERSION=8.4.21.
func TestADP05_ModoDualSimulador(t *testing.T) {
	simCode := string(readFile(t, sourcePath("backend", "cmd/proxmox-simulador/simulador.go")))

	t.Run("simulador lee PROXMOX_SIM_TARGET_VERSION", func(t *testing.T) {
		if !regexp.MustCompile(`(?i)PROXMOX_SIM_TARGET_VERSION|TargetVersion`).MatchString(simCode) {
			t.Errorf("ADP-05: cmd/proxmox-simulador/simulador.go no implementa PROXMOX_SIM_TARGET_VERSION")
		}
	})
}

// TestADP06_CatalogoDebian12Simulador verifica ADP-06:
// storage.go no debe hardcodear Debian 13 en el simulador para PVE 8.4 y debe ofrecer ISO oficial de Debian 12.
func TestADP06_CatalogoDebian12Simulador(t *testing.T) {
	storageCode := string(readFile(t, sourcePath("backend", "cmd/proxmox-simulador/storage.go")))

	t.Run("storage.go no expone Debian 13 y contiene ISO oficial de Debian 12 Bookworm", func(t *testing.T) {
		if strings.Contains(storageCode, "debian-13") {
			t.Errorf("ADP-06: cmd/proxmox-simulador/storage.go contiene referencias hardcodeadas a debian-13; debe adaptarse a Debian 12 Bookworm para PVE 8.4")
		}
		if !strings.Contains(storageCode, "local:iso/debian-12.7.0-amd64-netinst.iso") {
			t.Errorf("ADP-06: storage.go debe ofrecer local:iso/debian-12.7.0-amd64-netinst.iso")
		}
		for _, plantilla := range []string{"local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst", "local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst"} {
			if !strings.Contains(storageCode, plantilla) {
				t.Errorf("ADP-06: storage.go debe ofrecer la plantilla %s", plantilla)
			}
		}
	})

	t.Run("el storage local no anuncia el content type import (inexistente en PVE 8.4)", func(t *testing.T) {
		if regexp.MustCompile(`(?i)content[^\n]*\bimport\b|"import"`).MatchString(storageCode) {
			t.Errorf("ADP-06: storage.go sigue anunciando el content type 'import' en local; PVE 8.4 no lo tiene")
		}
	})
}

// TestADP09_ADP10_InfraestructuraPVE8 verifica ADP-09 y ADP-10:
// Requiere la presencia del script de auditoría y validación de host real scripts/pve8-audit-probe.sh.
func TestADP09_ADP10_InfraestructuraPVE8(t *testing.T) {
	scriptPath := filepath.Join(envOrDefault("CENTINELA_ROOT", "../.."), "scripts", "pve8-audit-probe.sh")

	t.Run("script de auditoría scripts/pve8-audit-probe.sh existe y verifica permisos y storage", func(t *testing.T) {
		info, err := os.Stat(scriptPath)
		if os.IsNotExist(err) {
			t.Fatalf("ADP-09/ADP-10: falta el script de validación de sonda en host real %q", scriptPath)
		}
		if info.Mode()&0111 == 0 {
			t.Errorf("ADP-09/ADP-10: el script %q debe tener permisos de ejecución (+x)", scriptPath)
		}
		content := string(readFile(t, scriptPath))
		// ADP-09: los once privilegios del rol CentinelaAdmin.
		for _, requiredPriv := range []string{"VM.Allocate", "Datastore.AllocateSpace", "VM.Config.CPU", "VM.Config.Memory", "VM.Config.Disk",
			"VM.Config.Network", "VM.PowerMgmt", "VM.Audit", "VM.Snapshot", "VM.Snapshot.Rollback", "Datastore.Audit"} {
			if !strings.Contains(content, requiredPriv) {
				t.Errorf("ADP-09: scripts/pve8-audit-probe.sh debe verificar el privilegio %q", requiredPriv)
			}
		}
		// ADP-10: storage apto para snapshots y plantilla Debian 12.
		for _, requiredStorage := range []string{"local-lvm", "lvmthin", "debian-12-standard", "vztmpl"} {
			if !strings.Contains(content, requiredStorage) {
				t.Errorf("ADP-10: scripts/pve8-audit-probe.sh debe verificar %q", requiredStorage)
			}
		}
	})

	t.Run("la documentación api-proxmox.md describe el procedimiento PVE 8.4.21", func(t *testing.T) {
		doc := string(readFile(t, filepath.Join(envOrDefault("CENTINELA_ROOT", "../.."), "documentacion", "api-proxmox.md")))
		for _, req := range []string{"centi-api@pve", "CentinelaAdmin", "pveum user token add", "pveam download local debian-12-standard"} {
			if !strings.Contains(doc, req) {
				t.Errorf("ADP-09/ADP-10: documentacion/api-proxmox.md debe documentar %q", req)
			}
		}
	})
}
