import { render, screen, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import InstancesPage from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { withProtectedProviders } from './app-providers';

// FRN-20C (RF-03): columnas de CPU (%) y RAM (usado / total) en el inventario de instancias.
// Criterio de éxito: los valores llegan de GET /api/instances (contrato BAC-29:
// cpuUsage, ramUsage y maxRam) y se muestra "—" cuando vienen en null.
// La tarea no fija el formato exacto; se acepta cualquier presentación que:
//   - refleje el valor real de cpuUsage como porcentaje (p. ej. "25 %"), y
//   - muestre la RAM usada y su tope (p. ej. "4 GiB / 16 GiB").
const inventory = [
  {
    id: 101, name: 'servidor-web', type: 'vm', node: 'pve', status: 'running',
    ip: '192.168.1.50', cpuUsage: 25, ramUsage: 4294967296, maxRam: 17179869184,
    nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
  {
    id: 102, name: 'base-datos', type: 'lxc', node: 'pve', status: 'stopped',
    ip: null, cpuUsage: null, ramUsage: null, maxRam: null,
    nivelAcceso: 'READ_ONLY', activeTask: null,
  },
];

function renderInstances() {
  window.sessionStorage.setItem('centinela_access', 'access-token-metricas');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
    rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));
  vi.stubGlobal('fetch', vi.fn().mockImplementation((url: string) => Promise.resolve(String(url).includes('/instances')
    ? new Response(JSON.stringify(inventory), { status: 200, headers: { 'Content-Type': 'application/json' } })
    : new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))));
  return render(withProtectedProviders(
    <Toaster>
      <MemoryRouter initialEntries={['/instances']}>
        <InstancesPage />
      </MemoryRouter>
    </Toaster>,
  ));
}

async function rowOf(name: string) {
  const cell = await screen.findByText(name, {}, { timeout: 3000 });
  const row = cell.closest('tr, [role="row"]');
  expect(row, `"${name}" no está dentro de una fila de tabla`).not.toBeNull();
  return row as HTMLElement;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-20C - Columnas de CPU/RAM en el inventario', () => {
  it('muestra el consumo de CPU (%) y de RAM (usado / total) de cada instancia', async () => {
    renderInstances();

    const row = await rowOf('servidor-web');
    // CPU en porcentaje: el contrato envía 25 (escala 0 a 100), debe renderizarse como 25 %.
    expect(within(row).getByText(/25(?:[.,]0+)?\s*%/)).toBeInTheDocument();
    // RAM: 4 GiB usados de 16 GiB de tope (los bytes del contrato, permitiendo decimales opcionales).
    expect(within(row).getByText(/4(?:\.0+)?\s*(?:GiB|GB)\s*\/\s*16(?:\.0+)?\s*(?:GiB|GB)/i)).toBeInTheDocument();
  });

  it('usa "—" cuando cpuUsage y ramUsage llegan en null', async () => {
    renderInstances();

    const row = await rowOf('base-datos');
    // CPU y RAM en null → dos celdas con marcador de ausencia.
    expect(within(row).getAllByText('—').length).toBeGreaterThanOrEqual(2);
  });

  it('la cabecera de la tabla identifica las columnas de CPU y RAM', async () => {
    renderInstances();

    const row = await rowOf('servidor-web');
    const table = row.closest('table, [role="table"], [role="grid"]') as HTMLElement;
    expect(within(table).getAllByText(/cpu/i).length).toBeGreaterThan(0);
    expect(within(table).getAllByText(/ram|memoria/i).length).toBeGreaterThan(0);
  });
});
