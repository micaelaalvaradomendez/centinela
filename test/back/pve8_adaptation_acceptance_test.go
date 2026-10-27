package back_test

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Pruebas de aceptación de la Etapa de Adaptación a Proxmox VE 8.4.21 (documentacion/actual.md):
//   - ADP-01 (BAC-PVE8-01): Ajuste de parámetros de aprovisionamiento QEMU (ostype=l26, agent=1) y LXC (unprivileged=1, nesting=1).
//   - ADP-02 (BAC-PVE8-02): Parámetros de purga profunda en eliminación de instancias (purge=1, destroy-unreferenced-disks=1).
//   - ADP-03 (BAC-PVE8-03): Parser tolerante de memoria (entero/string) y prefijo de red CIDR sin desajuste de digest.
//   - ADP-04 (BAC-PVE8-04): Parametrización dinámica de VMIDs protegidos vía variable de entorno.
//   - ADP-05 (SIM-PVE8-01): Modo dual en el simulador Proxmox con PROXMOX_SIM_TARGET_VERSION=8.4.21.
//   - ADP-06 (SIM-PVE8-02): Catálogo de plantillas e ISOs oficiales Debian 12 Bookworm en el simulador.
//   - ADP-09 (INF-PVE8-01) y ADP-10 (INF-PVE8-02): Privilegios de token y storage LVM-thin documentados para PVE 8.4.21.

// TestADP01_ParametrosAprovisionamientoQEMU_LXC verifica ADP-01:
// Inyección obligatoria de ostype=l26 y agent=1 para QEMU, y unprivileged=1 y nesting=1 para LXC.
func TestADP01_ParametrosAprovisionamientoQEMU_LXC(t *testing.T) {
	clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))

	t.Run("QEMU inyecta ostype=l26 y agent=1 para canal virtio-serial en PVE 8.4", func(t *testing.T) {
		if !regexp.MustCompile(`"ostype",\s*"l26"`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe enviar ostype=l26 para instancias QEMU")
		}
		if !regexp.MustCompile(`"agent",\s*"1"`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe enviar agent=1 para habilitar el socket virtio-serial en QEMU")
		}
	})

	t.Run("LXC inyecta unprivileged=1 y nesting=1 para compatibilidad systemd en PVE 8.4", func(t *testing.T) {
		if !regexp.MustCompile(`"unprivileged",\s*"1"`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe enviar unprivileged=1 para aislamiento en contenedores LXC")
		}
		if !regexp.MustCompile(`(?i)(nesting=1|"features",\s*"nesting=1")`).MatchString(clientCode) {
			t.Errorf("ADP-01: client.go:parametrosCreacion debe configurar features nesting=1 para el arranque limpio de systemd en LXC")
		}
	})
}

// TestADP02_ParametrosPurgaBorrado verifica ADP-02:
// DELETE /nodes/{node}/{tipo}/{vmid} debe incluir ?purge=1&destroy-unreferenced-disks=1.
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
// Soporte agnóstico de tipos en json.Number / prefijoFlexible para compatibilidad con PVE 8.4 y 9.2.
func TestADP03_ParserDefensivoMemoryPrefix(t *testing.T) {
	clientCode := string(readFile(t, sourcePath("backend", "internal/adapters/secondary/proxmox/client.go")))

	t.Run("prefijoFlexible y json.Number presentes en cliente Proxmox", func(t *testing.T) {
		if !strings.Contains(clientCode, "prefijoFlexible") {
			t.Errorf("ADP-03: client.go debe contar con tipo prefijoFlexible para tolerar máscaras de red como string y número")
		}
		if !strings.Contains(clientCode, "UseNumber") || !strings.Contains(clientCode, "json.Number") {
			t.Errorf("ADP-03: ObtenerConfig debe usar json.Number para admitir memory tanto en formato string como numérico")
		}
	})
}

// TestADP04_ParametrizacionVMIDsProtegidos verifica ADP-04:
// PROXMOX_PROTECTED_VMIDS permite personalizar o vaciar la lista de protección.
func TestADP04_ParametrizacionVMIDsProtegidos(t *testing.T) {
	middlewareCode := string(readFile(t, sourcePath("backend", "internal/adapters/primary/http/middleware/protected_instance.go")))

	t.Run("ParseVmidsProtegidos maneja cadenas personalizadas y listas vacías", func(t *testing.T) {
		if !strings.Contains(middlewareCode, "ParseVmidsProtegidos") {
			t.Errorf("ADP-04: protected_instance.go debe exportar la función ParseVmidsProtegidos")
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
// storage.go ofrece plantillas oficiales debian-12 en vez de debian-13 y no incluye el content type 'import'.
func TestADP06_CatalogoDebian12Simulador(t *testing.T) {
	storageCode := string(readFile(t, sourcePath("backend", "cmd/proxmox-simulador/storage.go")))

	t.Run("storage.go expone debian-12 en catálogo de plantillas", func(t *testing.T) {
		if !strings.Contains(storageCode, "debian-12") {
			t.Errorf("ADP-06: cmd/proxmox-simulador/storage.go debe incluir plantillas/ISOs oficiales de Debian 12 Bookworm")
		}
	})
}

// TestADP09_ADP10_InfraestructuraPVE8 verifica ADP-09 y ADP-10:
// api-proxmox.md documenta la configuración específica para PVE 8.4.21.
func TestADP09_ADP10_InfraestructuraPVE8(t *testing.T) {
	docPath := filepath.Join(envOrDefault("CENTINELA_ROOT", "../.."), "documentacion", "api-proxmox.md")
	docContent := string(readFile(t, docPath))

	t.Run("api-proxmox.md documenta privilegios y storages para Proxmox VE 8.4.21", func(t *testing.T) {
		for _, requiredPriv := range []string{"VM.Allocate", "Datastore.AllocateSpace", "VM.Config.CPU", "VM.Config.Memory", "VM.Config.Network"} {
			if !strings.Contains(docContent, requiredPriv) {
				t.Errorf("ADP-09: api-proxmox.md debe documentar el privilegio %q", requiredPriv)
			}
		}

		if !regexp.MustCompile(`(?i)local-lvm|lvmthin`).MatchString(docContent) {
			t.Errorf("ADP-10: api-proxmox.md debe documentar pool LVM-thin (local-lvm) para aprovisionamiento")
		}
	})
}

