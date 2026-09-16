package back_test

import "testing"

// Este archivo reserva el contrato del hito de Control de Acceso Basado en Recursos.
// Las pruebas se activan cuando el backend publique las rutas de permisos,
// inventario y operaciones sobre instancias descritas en documentacion/actual.md.
func TestHitoControlDeAccesoBasadoEnRecursos(t *testing.T) {
	t.Run("BAC-07 expone GET y PUT /api/admin/users/:id/permissions", func(t *testing.T) {
		t.Skip("pendiente: el backend actual expone PUT /api/users/:id/instances y no GET /permissions")
	})

	t.Run("BAC-08 rechaza instancia no asignada sin llamar a Proxmox", func(t *testing.T) {
		t.Skip("pendiente: falta endpoint de instancia protegido y un doble de Proxmox observable")
	})

	t.Run("BAC-14 filtra GET /api/instances por permisos y normaliza Proxmox", func(t *testing.T) {
		t.Skip("pendiente: falta GET /api/instances y contrato de adaptador Proxmox para la suite")
	})
}
