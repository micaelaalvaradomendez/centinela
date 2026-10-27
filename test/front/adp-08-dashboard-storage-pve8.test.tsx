import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import Dashboard from '@/pages/Dashboard';
import { StorageMeter } from '@/components/features/dashboard/components/ResourceMeter';
import { fetchCreationStorages } from '@/components/features/instances/services/instanceService';
import { withProtectedProviders } from './app-providers';

// ADP-08: Validación y Renderizado Frontend de Storages y Métricas PVE 8.
// Probar que el Dashboard (Dashboard, StorageMeter) y los selectores de almacenamiento
// manejan correctamente las respuestas de Proxmox VE 8.4.21:
// - Storages tipo "dir" y "lvmthin" sin flag 'import' en contentTypes (ej. ['iso', 'vztmpl'] o ['images', 'rootdir']).
// - Métricas de almacenamiento numéricas directas sin causar NaN ni divisiones por cero.
// - Renderizado correcto de gauges y barras de saturación.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const GIB = 1024 ** 3;

const pve8StoragesData = [
  {
    name: 'local',
    type: 'dir',
    totalBytes: 100 * GIB,
    freeBytes: 80 * GIB,
    contentTypes: ['iso', 'vztmpl', 'backup'], // Sin flag 'import' (específico de PVE 8.4)
  },
  {
    name: 'local-lvm',
    type: 'lvmthin',
    totalBytes: 500 * GIB,
    freeBytes: 350 * GIB,
    contentTypes: ['images', 'rootdir'],
  },
];

const pve8NodeStatus = {
  cpu: { usagePercent: 22, cores: 16 },
  ram: { usedGb: 8.5, totalGb: 32, usagePercent: 26.56 },
  storage: { usedGb: 170, totalGb: 600, usagePercent: 28.33 },
  uptimeSeconds: 86400,
  instancesSummary: {
    vms: { running: 2, stopped: 1, paused: 0, total: 3 },
    lxc: { running: 4, stopped: 0, paused: 0, total: 4 },
  },
  stale: false,
  fetchedAt: '2026-10-27T10:00:00Z',
};

describe('ADP-08: Validación y Renderizado Frontend de Storages y Métricas PVE 8', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('fetchCreationStorages parsea correctamente almacenamientos PVE 8 sin flag import', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(json(pve8StoragesData)));

    const storages = await fetchCreationStorages();
    expect(storages).toHaveLength(2);

    const localDir = storages.find((s) => s.value === 'local');
    expect(localDir).toBeDefined();
    expect(localDir?.type).toBe('dir');
    expect(localDir?.contentTypes).toContain('iso');
    expect(localDir?.contentTypes).toContain('vztmpl');
    expect(localDir?.contentTypes).not.toContain('import');
    expect(localDir?.availableGb).toBeCloseTo(80, 1);
    expect(localDir?.totalGb).toBeCloseTo(100, 1);

    const localLvm = storages.find((s) => s.value === 'local-lvm');
    expect(localLvm).toBeDefined();
    expect(localLvm?.type).toBe('lvmthin');
    expect(localLvm?.contentTypes).toContain('images');
  });

  it('StorageMeter renderiza métricas de almacenamiento PVE 8 sin NaN ni errores', () => {
    render(<StorageMeter usedGb={170} totalGb={600} usagePercent={28.33} />);

    // Verifica que se renderiza el porcentaje y los valores en GB sin NaN
    expect(screen.getByText(/28(\.33)?\s*%/)).toBeInTheDocument();
    expect(screen.getByText(/170/)).toBeInTheDocument();
    expect(screen.getByText(/600/)).toBeInTheDocument();
    expect(screen.queryByText(/NaN/i)).not.toBeInTheDocument();
  });

  it('StorageMeter maneja valores límite (0% y >70%) correctamente', () => {
    const { rerender } = render(<StorageMeter usedGb={0} totalGb={500} usagePercent={0} />);
    expect(screen.getByText(/0\s*%/)).toBeInTheDocument();

    rerender(<StorageMeter usedGb={400} totalGb={500} usagePercent={80} />);
    expect(screen.getByText(/80\s*%/)).toBeInTheDocument();
    expect(screen.queryByText(/NaN/i)).not.toBeInTheDocument();
  });

  it('DashboardPage carga e integra telemetría de nodo y storage PVE 8', async () => {
    window.sessionStorage.setItem('centinela_access', 'access-token-dashboard-pve8');
    window.localStorage.setItem(
      'centinela_user',
      JSON.stringify({
        id: 'user-admin',
        organizacionId: 'org-1',
        nombreCompleto: 'Admin Centinela',
        email: 'admin@centinela.local',
        rol: 'ADMIN',
        instanciasPermitidas: [],
        tiene2FA: true,
      })
    );

    vi.stubGlobal(
      'fetch',
      vi.fn().mockImplementation((url: string) => {
        const target = String(url);
        if (/\/node\/status(\?|$)/.test(target)) return Promise.resolve(json(pve8NodeStatus));
        if (/\/node\/templates(\?|$)/.test(target)) return Promise.resolve(json([]));
        if (/\/instances(\?|$)/.test(target)) return Promise.resolve(json([]));
        if (/\/audit(\?|$)/.test(target)) return Promise.resolve(json([]));
        return Promise.resolve(json({}));
      })
    );

    render(
      withProtectedProviders(
        <MemoryRouter initialEntries={['/dashboard']}>
          <Dashboard />
        </MemoryRouter>
      )
    );

    // Esperar a que se carguen los datos del nodo
    const storageElements = await screen.findAllByText(/170/);
    expect(storageElements.length).toBeGreaterThan(0);
    expect(screen.getAllByText(/600/).length).toBeGreaterThan(0);
    expect(screen.queryByText(/NaN/i)).not.toBeInTheDocument();
  });
});
