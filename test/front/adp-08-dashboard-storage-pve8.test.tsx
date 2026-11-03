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
// - PVE 8 suele enviar el atributo "content" como string separado por comas ("iso,vztmpl")
//   o sin el flag 'import'. El servicio debe normalizarlo sin arrojar ApiRequestError.
// - Métricas de almacenamiento numéricas directas sin causar NaN ni divisiones por cero.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const GIB = 1024 ** 3;

// Payload de GET /api/node/storages sobre un nodo Proxmox VE 8.4:
// el backend normaliza "content" en contentTypes: string[], sin el content type "import" de PVE 9.
const pve8Storages = [
  {
    name: 'local',
    type: 'dir',
    totalBytes: 100 * GIB,
    freeBytes: 80 * GIB,
    contentTypes: ['iso', 'vztmpl', 'backup'],
  },
  {
    name: 'local-lvm',
    type: 'lvmthin',
    totalBytes: 500 * GIB,
    freeBytes: 350 * GIB,
    contentTypes: ['images', 'rootdir'],
  },
];

describe('ADP-08: Validación y Renderizado Frontend de Storages y Métricas PVE 8', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('fetchCreationStorages procesa catálogos de PVE 8 sin el content type import sin fallar con ApiRequestError', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(json(pve8Storages)));

    const storages = await fetchCreationStorages();
    expect(storages).toHaveLength(2);

    const localDir = storages.find((s) => s.value === 'local');
    expect(localDir).toBeDefined();
    expect(localDir?.type).toBe('dir');
    expect(localDir?.contentTypes).toContain('iso');
    expect(localDir?.contentTypes).toContain('vztmpl');
    expect(localDir?.contentTypes).not.toContain('import');
  });

  it('StorageMeter renderiza métricas numéricas directas de PVE 8 sin NaN ni errores', () => {
    render(<StorageMeter usedGb={170} totalGb={600} usagePercent={28.33} />);

    expect(screen.getByText(/28(\.33)?\s*%/)).toBeInTheDocument();
    expect(screen.getByText(/170/)).toBeInTheDocument();
    expect(screen.getByText(/600/)).toBeInTheDocument();
    expect(screen.queryByText(/NaN/i)).not.toBeInTheDocument();
  });
});
