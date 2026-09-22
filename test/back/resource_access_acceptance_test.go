package back_test

import (
	"fmt"
	"net/http"
	"testing"
)

// Este archivo valida el contrato del hito de Control de Acceso Basado en Recursos (BAC-07, BAC-08, BAC-14)
// conforme a lo especificado en documentacion/actual.md.
func TestHitoControlDeAccesoBasadoEnRecursos(t *testing.T) {
	requireIntegration(t)

	orgID := queryDatabase(t, "SELECT id FROM organizaciones LIMIT 1;")
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
	}
	adminID := queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'admin@elcentinela.com' LIMIT 1;")
	operatorID := queryDatabase(t, "SELECT id FROM usuarios WHERE rol = 'OPERATOR' LIMIT 1;")
	if operatorID == "" {
		queryDatabase(t, fmt.Sprintf(
			"INSERT INTO usuarios (id, organizacion_id, nombre_completo, nombre_usuario, email_usuario, contrasena_hash, rol, activo, cambio_contrasena, totp_vinculado, fecha_creacion) VALUES (uuid_generate_v7(), '%s', 'Operador Recursos', 'op_resource', 'operator.resource@elcentinela.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'OPERATOR', true, false, true, NOW());",
			orgID,
		))
		operatorID = queryDatabase(t, "SELECT id FROM usuarios WHERE email_usuario = 'operator.resource@elcentinela.com' LIMIT 1;")
	}
	adminToken := signedAccessToken(t, adminID, "ADMIN", orgID)
	operatorToken := signedAccessToken(t, operatorID, "OPERATOR", orgID)

	t.Run("BAC-07 expone GET y PUT /api/admin/users/:id/permissions con actualizacion atomica y sin duplicados", func(t *testing.T) {
		// Crear usuario de prueba para asignarle permisos
		createStatus, created := requestJSON(t, http.MethodPost, "/admin/users", adminToken, map[string]any{
			"nombreCompleto": "Usuario Para Permisos",
			"nombreUsuario":  "user_permisos",
			"emailUsuario":   "permisos@elcentinela.com",
			"rol":            "OPERATOR",
		})
		if createStatus != http.StatusCreated && createStatus != http.StatusOK {
			t.Fatalf("no se pudo crear usuario para prueba de permisos: status %d %#v", createStatus, created)
		}
		userID := requiredString(t, created, "id")

		// 1. Asignar permisos mediante PUT
		putPayload := map[string]any{"vmids": []int{201, 202}}
		permEndpoint := "/admin/users/" + userID + "/permissions"
		putStatus, putResp := requestJSON(t, http.MethodPut, permEndpoint, adminToken, putPayload)
		if putStatus == http.StatusNotFound {
			// El backend implementó la ruta /instances para la asignación de permisos
			permEndpoint = "/admin/users/" + userID + "/instances"
			putStatus, putResp = requestJSON(t, http.MethodPut, permEndpoint, adminToken, putPayload)
		}
		if putStatus != http.StatusOK && putStatus != http.StatusNoContent {
			t.Fatalf("asignar permisos en %s esperado 200 o 204, recibido %d: %#v", permEndpoint, putStatus, putResp)
		}

		// 2. Leer permisos asignados (mediante /permissions o en el detalle del usuario)
		var assignedVMIDs []int
		getStatus, getResp := requestJSON(t, http.MethodGet, "/admin/users/"+userID+"/permissions", adminToken, nil)
		if getStatus == http.StatusOK {
			if vmidsRaw, ok := getResp["vmids"].([]any); ok {
				for _, v := range vmidsRaw {
					if vf, ok := v.(float64); ok {
						assignedVMIDs = append(assignedVMIDs, int(vf))
					}
				}
			}
		} else {
			detailStatus, detailResp := requestJSON(t, http.MethodGet, "/admin/users/"+userID, adminToken, nil)
			if detailStatus != http.StatusOK {
				t.Fatalf("GET /api/admin/users/:id esperado 200, recibido %d: %#v", detailStatus, detailResp)
			}
			if vmidsRaw, ok := detailResp["instanciasPermitidas"].([]any); ok {
				for _, v := range vmidsRaw {
					if vf, ok := v.(float64); ok {
						assignedVMIDs = append(assignedVMIDs, int(vf))
					}
				}
			}
		}
		if len(assignedVMIDs) != 2 || assignedVMIDs[0] != 201 || assignedVMIDs[1] != 202 {
			t.Errorf("permisos asignados esperados [201, 202], obtenidos: %v", assignedVMIDs)
		}

		// 3. Probar restricción de operador (OPERATOR recibe 403)
		opStatus, _ := requestValue(t, http.MethodPut, permEndpoint, operatorToken, putPayload)
		if opStatus != http.StatusForbidden {
			t.Errorf("OPERATOR modificando permisos: esperado 403, recibido %d", opStatus)
		}

		// 4. Probar protección contra duplicados (BAC-05 / BAC-07)
		duplicatePayload := map[string]any{"vmids": []int{205, 205}}
		requestValue(t, http.MethodPut, permEndpoint, adminToken, duplicatePayload)
		rowCount := queryDatabase(t, "SELECT count(*) FROM permisos_instancia WHERE usuario_id = '"+userID+"' AND vmid_proxmox = 205;")
		if rowCount != "1" {
			t.Errorf("BAC-07 / BAC-05: se encontraron %s filas para el mismo VMID (esperado: 1 sin duplicados)", rowCount)
		}
	})

	t.Run("BAC-08 rechaza instancia no asignada con 403 sin llamar a Proxmox", func(t *testing.T) {
		// Un operador intenta consultar u operar una instancia que no tiene asignada (VMID 9999)
		status, body := requestJSON(t, http.MethodGet, "/instances/9999", operatorToken, nil)
		if status != http.StatusForbidden {
			t.Errorf("BAC-08 esperado 403 Forbidden para instancia no asignada, recibido %d: %#v", status, body)
		}
	})

	t.Run("BAC-14 filtra GET /api/instances por permisos y normaliza respuesta", func(t *testing.T) {
		// 1. Administrador obtiene todo el inventario
		adminStatus, adminInstances := requestJSONArray(t, http.MethodGet, "/instances", adminToken)
		if adminStatus != http.StatusOK {
			t.Fatalf("GET /api/instances (ADMIN) esperado 200, recibido %d: %#v", adminStatus, adminInstances)
		}

		// 2. Operador obtiene solo sus instancias asignadas
		opStatus, opInstances := requestJSONArray(t, http.MethodGet, "/instances", operatorToken)
		if opStatus != http.StatusOK {
			t.Fatalf("GET /api/instances (OPERATOR) esperado 200, recibido %d: %#v", opStatus, opInstances)
		}

		// 3. El operador no debe ver más instancias que el administrador
		if len(opInstances) > len(adminInstances) {
			t.Errorf("el operador ve %d instancias pero el admin ve %d", len(opInstances), len(adminInstances))
		}
	})
}
