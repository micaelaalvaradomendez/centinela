package back_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// Etapa 1 (documentacion/actual.md):
//   - FIX-35 el simulador de Proxmox (backend/cmd/proxmox-simulador, BAC-28) responde
//     lxc/{vmid}/interfaces igual que el Proxmox real 9.2.2 (comparado el 30/09/2026):
//     ip-addresses[].prefix es string y un contenedor apagado responde 200 {"data": null}.
//
// El simulador se compila desde el submódulo (sin modificarlo) y se levanta en un puerto libre.
func TestEtapa1SimuladorProxmox(t *testing.T) {
	sim := startProxmoxSimulator(t)

	t.Run("FIX-35 lxc interfaces informa prefix como string, igual que el Proxmox real", func(t *testing.T) {
		status, body := sim.get(t, "/nodes/proxmox/lxc/101/interfaces") // 101 "back": lxc encendido
		if status != http.StatusOK {
			t.Fatalf("GET lxc/101/interfaces (encendido): esperado 200, recibido %d: %s", status, body)
		}
		var response struct {
			Data []struct {
				Name        string           `json:"name"`
				IPAddresses []map[string]any `json:"ip-addresses"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("respuesta no es JSON: %v: %s", err, body)
		}
		checked := 0
		for _, iface := range response.Data {
			for _, address := range iface.IPAddresses {
				checked++
				if _, ok := address["prefix"].(string); !ok {
					t.Errorf("interfaz %s, %v: prefix debe ser string (\"24\") como en el Proxmox real; recibido %T %v", iface.Name, address["ip-address"], address["prefix"], address["prefix"])
				}
			}
		}
		if checked == 0 {
			t.Errorf("la respuesta no trae ip-addresses: %s", body)
		}
	})

	t.Run("FIX-35 lxc interfaces de un contenedor apagado responde 200 con data null", func(t *testing.T) {
		status, body := sim.get(t, "/nodes/proxmox/lxc/201/interfaces") // 201: lxc apagado
		if status != http.StatusOK {
			t.Fatalf("GET lxc/201/interfaces (apagado): esperado 200 {\"data\":null} como el Proxmox real, recibido %d: %s", status, body)
		}
		var response map[string]any
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("respuesta no es JSON: %v: %s", err, body)
		}
		if value, ok := response["data"]; !ok || value != nil {
			t.Errorf("esperado {\"data\": null}, recibido %s", body)
		}
	})
}

type proxmoxSimulator struct {
	baseURL string
	auth    string
}

func (s proxmoxSimulator) get(t *testing.T, path string) (int, []byte) {
	t.Helper()
	request, _ := http.NewRequest(http.MethodGet, s.baseURL+path, nil)
	request.Header.Set("Authorization", s.auth)
	response, err := (&http.Client{Timeout: 5 * time.Second}).Do(request)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return response.StatusCode, body
}

func startProxmoxSimulator(t *testing.T) proxmoxSimulator {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "proxmox-simulador")
	build := exec.Command("go", "build", "-o", binary, "./cmd/proxmox-simulador")
	build.Dir = sourcePath("backend", "")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("no se pudo compilar backend/cmd/proxmox-simulador: %v\n%s", err, output)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("no hay puerto libre: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	const tokenID, tokenSecret = "sim@pve!tests", "sim-secret"
	cmd := exec.Command(binary)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PROXMOX_SIM_PORT=%d", port),
		"PROXMOX_NODE=proxmox",
		"PROXMOX_TOKEN_ID="+tokenID,
		"PROXMOX_TOKEN_SECRET="+tokenSecret,
	)
	if err := cmd.Start(); err != nil {
		t.Fatalf("no se pudo iniciar el simulador: %v", err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() })

	sim := proxmoxSimulator{
		baseURL: fmt.Sprintf("http://127.0.0.1:%d/api2/json", port),
		auth:    fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret),
	}
	for attempt := 0; attempt < 50; attempt++ {
		if conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 200*time.Millisecond); err == nil {
			conn.Close()
			return sim
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("el simulador no abrió el puerto %d", port)
	return sim
}
