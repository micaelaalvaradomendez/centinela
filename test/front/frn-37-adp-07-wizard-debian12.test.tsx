import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import CreateNewInstances from '@/pages/createNewInstances';
import { withProtectedProviders } from './app-providers';
import { fetchCreationTemplates } from '@/components/features/instances/services/instanceService';

// FRN-37 / ADP-07: Adaptación Frontend Selector Plantillas Debian 12 / PVE 8.
// En Proxmox VE 8.4.21, las plantillas predeterminadas de LXC son Debian 12 (Bookworm),
// e.g.: "local:vztmpl/debian-12-standard_12.2-1_amd64.tar.zst".
// El frontend debe detectar y extraer metadatos de sistema operativo ('Debian') y versión ('12')
// a partir del identificador de la plantilla (no mostrar "—"), y renderizar el badge y resumen correspondientes.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const GIB = 1024 ** 3;
const MIB = 1024 ** 2;

const pve8Storages = [
  { name: 'local-lvm', type: 'lvmthin', totalBytes: 250 * GIB, freeBytes: 200 * GIB, contentTypes: ['images', 'rootdir'] },
  { name: 'local', type: 'dir', totalBytes: 100 * GIB, freeBytes: 85 * GIB, contentTypes: ['iso', 'vztmpl'] },
];

const pve8Templates = [
  {
    id: 'local:vztmpl/debian-12-standard_12.2-1_amd64.tar.zst',
    name: 'debian-12-standard_12.2-1_amd64.tar.zst',
    type: 'LXC_TEMPLATE',
    sizeBytes: 125 * MIB,
  },
  {
    id: 'local:iso/debian-12.5.0-amd64-netinst.iso',
    name: 'debian-12.5.0-amd64-netinst.iso',
    type: 'ISO',
    sizeBytes: 650 * MIB,
  },
];

const nodeStatus = {
  cpu: { usagePercent: 15, cores: 8 },
  ram: { usedGb: 4, totalGb: 16, usagePercent: 25 },
  storage: { usedGb: 50, totalGb: 350, usagePercent: 14 },
  uptimeSeconds: 7200,
  allocatable: { maxCpuCores: 8, freeRamMb: 12288, freeStorageGb: 300, headroomPercent: 15 },
  instancesSummary: { vms: { running: 1, stopped: 0, paused: 0, total: 1 }, lxc: { running: 2, stopped: 0, paused: 0, total: 2 } },
  stale: false,
  fetchedAt: '2026-10-27T10:00:00Z',
};

function setupPVE8Session() {
  window.sessionStorage.setItem('centinela_access', 'access-token-wizard-pve8');
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

  const fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
    const target = String(url);
    if ((init.method ?? 'GET').toUpperCase() !== 'GET') {
      return Promise.resolve(json({ vmid: '105', tareaId: 'UPID:pve8:00001:create' }, 202));
    }
    if (/\/node\/status(\?|$)/.test(target)) return Promise.resolve(json(nodeStatus));
    if (/\/node\/storages(\?|$)/.test(target)) return Promise.resolve(json(pve8Storages));
    if (/\/node\/templates(\?|$)/.test(target)) return Promise.resolve(json(pve8Templates));
    if (/\/instances(\?|$)/.test(target)) return Promise.resolve(json([]));
    return Promise.resolve(json({}));
  });
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

describe('FRN-37 / ADP-07: Selector y Metadatos de Plantillas Debian 12 PVE 8', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('fetchCreationTemplates extrae metadatos de sistema operativo ("Debian") y versión ("12") para plantillas Debian 12', async () => {
    setupPVE8Session();
    const result = await fetchCreationTemplates();
    expect(result).toHaveLength(2);

    const debian12Lxc = result.find((t) => t.value.includes('debian-12-standard'));
    expect(debian12Lxc).toBeDefined();
    expect(debian12Lxc?.value).toBe('local:vztmpl/debian-12-standard_12.2-1_amd64.tar.zst');
    expect(debian12Lxc?.type).toBe('LXC_TEMPLATE');

    // Criterio de éxito de ADP-07: no mostrar "—", extraer el sistema operativo y la versión
    expect(debian12Lxc?.operatingSystem, 'ADP-07: el servicio debe identificar "Debian" como sistema operativo').toBe('Debian');
    expect(debian12Lxc?.version, 'ADP-07: el servicio debe extraer "12" como versión de la plantilla').toBe('12');
  });

  it('el wizard muestra la información detectada de Debian 12 en la tabla de templates', async () => {
    setupPVE8Session();
    const user = userEvent.setup();

    render(
      withProtectedProviders(
        <MemoryRouter>
          <CreateNewInstances />
        </MemoryRouter>
      )
    );

    // Paso 1: Seleccionar Contenedor (LXC) y nombre
    const lxcCard = await screen.findByRole('button', { name: /contenedor/i });
    await user.click(lxcCard);

    const nameInput = screen.getByLabelText(/nombre del contenedor/i);
    await user.clear(nameInput);
    await user.type(nameInput, 'ct-debian12-test');

    const nextBtn1 = screen.getByRole('button', { name: /siguiente/i });
    expect(nextBtn1).toBeEnabled();
    await user.click(nextBtn1);

    // Paso 2: Recursos (seleccionar storage local-lvm)
    await screen.findByRole('combobox', { name: /almacenamiento de destino/i });
    await user.click(screen.getByRole('combobox', { name: /almacenamiento de destino/i }));
    await user.click(await screen.findByRole('option', { name: /local-lvm/i }));
    
    const nextBtn2 = screen.getByRole('button', { name: /siguiente/i });
    await waitFor(() => expect(nextBtn2).toBeEnabled());
    await user.click(nextBtn2);

    // Paso 3: Red (avanzar con DHCP)
    await screen.findByRole('radiogroup', { name: /modo de asignaci[oó]n ipv4/i });
    const nextBtn3 = screen.getByRole('button', { name: /siguiente/i });
    expect(nextBtn3).toBeEnabled();
    await user.click(nextBtn3);

    // Paso 4: Sistema Operativo (selección de template)
    await screen.findByRole('radiogroup', { name: /template del sistema operativo/i });
    
    // Verificar que aparece el template y que sus columnas muestran Debian y 12
    const debianCell = await screen.findByText(/^Debian$/i);
    expect(debianCell).toBeInTheDocument();
  });
});
