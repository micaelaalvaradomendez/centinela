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
	"strings"
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

func (s proxmoxSimulator) post(t *testing.T, path string, bodyData io.Reader) (int, []byte) {
	t.Helper()
	request, _ := http.NewRequest(http.MethodPost, s.baseURL+path, bodyData)
	request.Header.Set("Authorization", s.auth)
	if bodyData != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response, err := (&http.Client{Timeout: 5 * time.Second}).Do(request)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return response.StatusCode, body
}

func (s proxmoxSimulator) delete(t *testing.T, path string) (int, []byte) {
	t.Helper()
	request, _ := http.NewRequest(http.MethodDelete, s.baseURL+path, nil)
	request.Header.Set("Authorization", s.auth)
	response, err := (&http.Client{Timeout: 5 * time.Second}).Do(request)
	if err != nil {
		t.Fatalf("DELETE %s: %v", path, err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return response.StatusCode, body
}

// TestEtapa2BAC50_SimuladorSnapshotsYRollback verifica los entregables de BAC-50 (Etapa 2):
// Extensión del simulador de Proxmox para Snapshots y Rollback (BRG-10-BAC).
// Entregables:
//  1. Handlers simulados en backend/cmd/proxmox-simulador:
//     - POST /api2/json/nodes/:node/qemu/:vmid/snapshot y POST .../lxc/... (persiste snapshot en memoria y devuelve UPID).
//     - GET /api2/json/nodes/:node/qemu/:vmid/snapshot (lista ordenada con name, snaptime, description).
//     - POST .../rollback y DELETE .../snapshot/:snapname (devuelven UPID y simulan la demora asíncrona).
// Criterio de éxito: pruebas al simulador permiten crear, listar y revertir un snapshot simulado con seguimiento de UPID.
func TestEtapa2BAC50_SimuladorSnapshotsYRollback(t *testing.T) {
	sim := startProxmoxSimulator(t)

	t.Run("BAC-50: Crear snapshot en QEMU y LXC devuelve UPID", func(t *testing.T) {
		payloadVM := strings.NewReader(`{"snapname":"snap-test-qemu","description":"Snapshot de prueba QEMU"}`)
		status, body := sim.post(t, "/nodes/proxmox/qemu/100/snapshot", payloadVM)
		if status != http.StatusOK {
			t.Fatalf("POST qemu/100/snapshot: esperado 200, recibido %d: %s", status, body)
		}
		var resVM struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal(body, &resVM); err != nil || !strings.HasPrefix(resVM.Data, "UPID:proxmox:") {
			t.Errorf("POST qemu/100/snapshot debe devolver UPID de tarea, recibido %s", body)
		}

		payloadLXC := strings.NewReader(`{"snapname":"snap-test-lxc","description":"Snapshot de prueba LXC"}`)
		statusLXC, bodyLXC := sim.post(t, "/nodes/proxmox/lxc/101/snapshot", payloadLXC)
		if statusLXC != http.StatusOK {
			t.Fatalf("POST lxc/101/snapshot: esperado 200, recibido %d: %s", statusLXC, bodyLXC)
		}
		var resLXC struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal(bodyLXC, &resLXC); err != nil || !strings.HasPrefix(resLXC.Data, "UPID:proxmox:") {
			t.Errorf("POST lxc/101/snapshot debe devolver UPID de tarea, recibido %s", bodyLXC)
		}
	})

	t.Run("BAC-50: GET lista los snapshots con name, snaptime y description", func(t *testing.T) {
		status, body := sim.get(t, "/nodes/proxmox/qemu/100/snapshot")
		if status != http.StatusOK {
			t.Fatalf("GET qemu/100/snapshot: esperado 200, recibido %d: %s", status, body)
		}
		var res struct {
			Data []map[string]any `json:"data"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			t.Fatalf("respuesta no es JSON: %v", err)
		}
		found := false
		for _, item := range res.Data {
			if item["name"] == "snap-test-qemu" {
				found = true
				if _, ok := item["snaptime"]; !ok {
					t.Errorf("snapshot no incluye snaptime: %#v", item)
				}
				if desc, _ := item["description"].(string); !strings.Contains(desc, "QEMU") {
					t.Errorf("snapshot no incluye description esperada: %#v", item)
				}
			}
		}
		if !found {
			t.Errorf("snapshot 'snap-test-qemu' no encontrado en el listado de snapshots: %s", body)
		}
	})

	t.Run("BAC-50: Rollback de snapshot devuelve UPID", func(t *testing.T) {
		status, body := sim.post(t, "/nodes/proxmox/qemu/100/snapshot/snap-test-qemu/rollback", nil)
		if status != http.StatusOK {
			t.Fatalf("POST rollback: esperado 200, recibido %d: %s", status, body)
		}
		var res struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal(body, &res); err != nil || !strings.HasPrefix(res.Data, "UPID:proxmox:") {
			t.Errorf("POST rollback debe devolver UPID de tarea, recibido %s", body)
		}
	})

	t.Run("BAC-50: DELETE snapshot devuelve UPID", func(t *testing.T) {
		status, body := sim.delete(t, "/nodes/proxmox/qemu/100/snapshot/snap-test-qemu")
		if status != http.StatusOK {
			t.Fatalf("DELETE snapshot: esperado 200, recibido %d: %s", status, body)
		}
		var res struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal(body, &res); err != nil || !strings.HasPrefix(res.Data, "UPID:proxmox:") {
			t.Errorf("DELETE snapshot debe devolver UPID de tarea, recibido %s", body)
		}
	})
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
		"PROXMOX_SIM_DURACION_TAREA=0",
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
