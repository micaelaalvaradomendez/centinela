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
		if !regexp.MustCompile(`"ostype":\s*"l26"`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe enviar ostype=l26 para instancias QEMU")
		}
		if !regexp.MustCompile(`"agent":\s*"1"`).MatchString(clientCode) && !regexp.MustCompile(`"agent":\s*1\b`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe enviar agent=1 para habilitar el socket virtio-serial en QEMU")
		}
	})

	t.Run("LXC inyecta unprivileged=1 y nesting=1 para compatibilidad systemd en PVE 8.4", func(t *testing.T) {
		if !regexp.MustCompile(`"unprivileged":\s*"1"`).MatchString(clientCode) && !regexp.MustCompile(`"unprivileged":\s*1\b`).MatchString(clientCode) {
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
		if !regexp.MustCompile(`(?i)DELETE.*purge=1`).MatchString(clientCode) &&
			!strings.Contains(clientCode, "purge=1") {
			t.Errorf("ADP-02: EliminarInstancia en client.go debe enviar purge=1 en la petición DELETE a Proxmox")
		}
		if !strings.Contains(clientCode, "destroy-unreferenced-disks=1") {
			t.Errorf("ADP-02: EliminarInstancia en client.go debe enviar destroy-unreferenced-disks=1 para evitar discos huérfanos")
		}
	})
}

// TestADP03_ParserDefensivoMemoryPrefix verifica ADP-03:
// Soporte agnóstico de tipos en json.Number / prefijoFlexible y existencia de suite unitaria dedicada client_pve8_test.go.
func TestADP03_ParserDefensivoMemoryPrefix(t *testing.T) {
	t.Run("archivo de pruebas unitarias client_pve8_test.go presente en proxmox adapter", func(t *testing.T) {
		unitTestPath := sourcePath("backend", "internal/adapters/secondary/proxmox/client_pve8_test.go")
		if _, err := os.Stat(unitTestPath); os.IsNotExist(err) {
			t.Fatalf("ADP-03: falta la suite unitaria internal/adapters/secondary/proxmox/client_pve8_test.go con payloads crudos de PVE 8.4")
		}
	})
}

// TestADP04_ParametrizacionVMIDsProtegidos verifica ADP-04:
// ParseVmidsProtegidos debe soportar rangos numéricos con guion (ej. "100-105, 200, 300-302").
func TestADP04_ParametrizacionVMIDsProtegidos(t *testing.T) {
	code := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/middleware/protected_instance.go")))
	testCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/middleware/protected_instance_test.go")))

	t.Run("ParseVmidsProtegidos soporta rangos con guion y listas mixtas", func(t *testing.T) {
		if !regexp.MustCompile(`strings\.Split.*-|\bstrings\.Cut.*-|\bstrings\.Contains.*-|\bSplit\(.*"-"`).MatchString(code) {
			t.Errorf("ADP-04: ParseVmidsProtegidos en protected_instance.go no implementa soporte para rangos numéricos con guion (ej. '100-105')")
		}
		if !regexp.MustCompile(`(?i)Test.*Rango|100-105|100-103`).MatchString(testCode) {
			t.Errorf("ADP-04: protected_instance_test.go debe incluir casos de prueba unitaria para rangos con guion")
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
		if !strings.Contains(storageCode, "debian-12") || !strings.Contains(storageCode, ".iso") {
			t.Errorf("ADP-06: cmd/proxmox-simulador/storage.go debe incluir ISO oficial de Debian 12 (ej. debian-12.5.0-amd64-netinst.iso)")
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
		for _, requiredPriv := range []string{"VM.Allocate", "Datastore.AllocateSpace", "local-lvm"} {
			if !strings.Contains(content, requiredPriv) {
				t.Errorf("ADP-09/ADP-10: scripts/pve8-audit-probe.sh debe verificar %q", requiredPriv)
			}
		}
	})
}
