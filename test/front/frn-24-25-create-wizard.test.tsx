import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import CreateNewInstances from '@/pages/createNewInstances';
import { withProtectedProviders } from './app-providers';

// FRN-24 / FRN-25: wizard de creación de instancias en 5 pasos
// (INSTANCE_CREATION_STEPS en createInstanceOptions.ts).
//   Paso 1 (General): tarjetas Máquina virtual (VM) / Contenedor (LXC) y nombre obligatorio.
//   Paso 2 (Recursos): controles numéricos de vCPU, RAM y Disco, selector de almacenamiento
//          y medidores reactivos de cuota (el avance se bloquea cuando se supera el 100%).
//   Paso 3 (Red): DHCP o IP estática (dirección IPv4 + prefijo CIDR + gateway).
//   Paso 4 (Sistema operativo): plantillas LXC (vztmpl) para contenedores e ISOs para VMs.
//   Paso 5 (Resumen): ficha completa y botón "Crear Instancia" solo tras confirmar.
// Criterio de éxito: no se avanza al Paso 3 si la configuración excede la memoria o el disco
// del nodo, y el botón final solo se habilita con los campos requeridos completos.
// Los mocks reproducen el contrato HTTP real de backend/internal/adapters/primary/http:
//   GET /api/node/storages -> [{ name, type, freeBytes, totalBytes, contentTypes }]
//   GET /api/node/templates -> [{ id, name, type: ISO|LXC_TEMPLATE, sizeBytes }]
// La forma exacta de los controles es libre: se buscan por rol y nombre accesible.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const GIB = 1024 ** 3;
const MIB = 1024 ** 2;
const storages = [
  { name: 'local-lvm', type: 'lvmthin', totalBytes: 200 * GIB, freeBytes: 180 * GIB, contentTypes: ['images', 'rootdir'] },
  { name: 'local', type: 'dir', totalBytes: 100 * GIB, freeBytes: 90 * GIB, contentTypes: ['iso', 'vztmpl'] },
];
const templates = [
  { id: 'local:vztmpl/debian-12-standard_12.2-1_amd64.tar.zst', name: 'debian-12-standard_12.2-1_amd64.tar.zst', type: 'LXC_TEMPLATE', sizeBytes: 120 * MIB },
  { id: 'local:iso/ubuntu-24.04.iso', name: 'ubuntu-24.04.iso', type: 'ISO', sizeBytes: 2000 * MIB },
];

const nodeStatus = (freeRamMb: number, freeStorageGb: number) => ({
  cpu: { usagePercent: 10, cores: 8 },
  ram: { usedGb: 2, totalGb: 16, usagePercent: 12 },
  storage: { usedGb: 30, totalGb: 300, usagePercent: 10 },
  uptimeSeconds: 3600,
  allocatable: { maxCpuCores: 8, freeRamMb, freeStorageGb, headroomPercent: 15 },
  instancesSummary: { vms: { running: 0, stopped: 0, paused: 0, total: 0 }, lxc: { running: 0, stopped: 0, paused: 0, total: 0 } },
  stale: false, fetchedAt: '2026-10-13T03:00:00Z',
});

let fetchMock: ReturnType<typeof vi.fn>;

function renderWizard(options: { freeRamMb?: number; freeStorageGb?: number } = {}) {
  const status = nodeStatus(options.freeRamMb ?? 65536, options.freeStorageGb ?? 1000);
  fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
    const target = String(url);
    if ((init.method ?? 'GET').toUpperCase() !== 'GET') return Promise.resolve(json({ vmid: 'inst_120', tareaId: 'tarea-wizard' }, 202));
    if (/\/node\/status(\?|$)/.test(target)) return Promise.resolve(json(status));
    if (/\/node\/storages(\?|$)/.test(target)) return Promise.resolve(json(storages));
    if (/\/node\/templates(\?|$)/.test(target)) return Promise.resolve(json(templates));
    if (/\/instances(\?|$)/.test(target)) return Promise.resolve(json([]));
    return Promise.resolve(json({}));
  });
  vi.stubGlobal('fetch', fetchMock);
  window.sessionStorage.setItem('centinela_access', 'access-token-wizard');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'admin-1', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
    rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));
  return render(withProtectedProviders(<MemoryRouter initialEntries={['/instances/new']}><CreateNewInstances /></MemoryRouter>));
}

const urlsCalled = () => (fetchMock.mock.calls as [string][]).map(([url]) => String(url));
const next = () => screen.getByRole('button', { name: /siguiente/i });
const createButtons = () => screen.getAllByRole('button', { name: /crear instancia/i });

async function chooseType(user: ReturnType<typeof userEvent.setup>, type: 'vm' | 'lxc') {
  const name = type === 'vm' ? /m[aá]quina virtual/i : /contenedor/i;
  await user.click(await screen.findByRole('button', { name }));
}

// El Paso 1 exige tipo y nombre: sin nombre el botón "Siguiente" queda deshabilitado.
async function toResourcesStep(user: ReturnType<typeof userEvent.setup>, type: 'vm' | 'lxc') {
  await chooseType(user, type);
  const nameInput = screen.getByLabelText(type === 'vm' ? /nombre de la instancia/i : /nombre del contenedor/i);
  await user.type(nameInput, type === 'vm' ? 'vm-prueba' : 'lxc-prueba');
  await user.click(next());
  await waitFor(() => expect(screen.queryAllByRole('spinbutton').length).toBeGreaterThanOrEqual(3));
}

async function pickStorage(user: ReturnType<typeof userEvent.setup>, option = /local-lvm/i) {
  await user.click(await screen.findByRole('combobox', { name: /almacenamiento de destino/i }));
  await user.click(await screen.findByRole('option', { name: option }));
  await waitFor(() => expect(next()).toBeEnabled());
}

async function toNetworkStep(user: ReturnType<typeof userEvent.setup>, type: 'vm' | 'lxc') {
  await toResourcesStep(user, type);
  await pickStorage(user);
  await user.click(next());
  await screen.findByRole('radiogroup', { name: /modo de asignaci[oó]n ipv4/i });
}

async function toOperatingSystemStep(user: ReturnType<typeof userEvent.setup>, type: 'vm' | 'lxc') {
  await toNetworkStep(user, type);
  await user.click(next());
  await screen.findByRole('radiogroup', { name: /template del sistema operativo/i });
}

async function toReviewStep(user: ReturnType<typeof userEvent.setup>, type: 'vm' | 'lxc') {
  await toOperatingSystemStep(user, type);
  const image = type === 'vm' ? /seleccionar ubuntu-24\.04\.iso/i : /seleccionar debian-12-standard/i;
  await user.click(await screen.findByRole('radio', { name: image }));
  await waitFor(() => expect(next()).toBeEnabled());
  await user.click(next());
  await screen.findByRole('heading', { name: /resumen de la configuraci[oó]n/i });
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-24 - Paso 1 (general) y Paso 2 (recursos y cuotas en vivo)', () => {
  it('el Paso 1 ofrece tarjetas de Máquina virtual (VM) y Contenedor (LXC) con su descripción', async () => {
    renderWizard();
    expect(await screen.findByRole('button', { name: /m[aá]quina virtual/i })).toBeVisible();
    expect(screen.getByRole('button', { name: /contenedor/i })).toBeVisible();
    expect(screen.getAllByText(/sistema operativo/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText(/ligero/i).length).toBeGreaterThan(0);
  });

  it('el Paso 2 expone controles para vCPU, RAM y Disco', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toResourcesStep(user, 'vm');

    const controls = screen.queryAllByRole('spinbutton');
    expect(controls.length, 'faltan controles de vCPU, RAM y Disco').toBeGreaterThanOrEqual(3);
    for (const label of [/cpu|n[uú]cleos/i, /ram|memoria/i, /disco/i]) {
      expect(screen.getAllByText(label).length).toBeGreaterThan(0);
    }
  });

  it('muestra un medidor de cuota del host que reacciona al aumentar la RAM', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toResourcesStep(user, 'vm');
    await waitFor(() => expect(urlsCalled().some((url) => /\/node\/status/.test(url))).toBe(true));

    const ramMeter = () => {
      const meter = screen.queryAllByRole('progressbar')
        .find((element) => /cuota de memoria ram/i.test(element.getAttribute('aria-label') ?? ''));
      expect(meter, 'no hay un medidor de cuota de RAM (progressbar) del host').toBeDefined();
      return meter as HTMLElement;
    };
    const before = Number(ramMeter().getAttribute('aria-valuenow'));

    await user.click(screen.getByRole('button', { name: /aumentar memoria ram/i }));
    await waitFor(() => expect(Number(ramMeter().getAttribute('aria-valuenow'))).toBeGreaterThan(before));
  });

  it('no deja avanzar al Paso 3 si la RAM pedida excede la memoria libre del nodo', async () => {
    renderWizard({ freeRamMb: 256 });
    const user = userEvent.setup();
    await toResourcesStep(user, 'vm');

    expect((await screen.findAllByText(/excede|supera|insuficiente|sin capacidad|cuota/i)).length).toBeGreaterThan(0);
    expect(next()).toBeDisabled();
    expect(screen.queryByRole('radiogroup', { name: /modo de asignaci[oó]n ipv4/i })).toBeNull();
  });

  it('no deja avanzar al Paso 3 si el disco pedido excede el almacenamiento libre', async () => {
    renderWizard({ freeStorageGb: 1 });
    const user = userEvent.setup();
    await toResourcesStep(user, 'vm');

    expect((await screen.findAllByText(/excede|supera|insuficiente|sin capacidad|cuota/i)).length).toBeGreaterThan(0);
    expect(next()).toBeDisabled();
  });

  it('con capacidad suficiente sí permite avanzar al Paso 3 (Red)', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toNetworkStep(user, 'vm');

    expect(screen.getByRole('radiogroup', { name: /modo de asignaci[oó]n ipv4/i })).toBeInTheDocument();
  });
});

describe('FRN-25 - Pasos 3 (red), 4 (sistema operativo) y 5 (resumen)', () => {
  it('lista los storages reales de GET /api/node/storages', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toResourcesStep(user, 'vm');

    await waitFor(() => expect(urlsCalled().some((url) => /\/node\/storages/.test(url))).toBe(true));
    await user.click(await screen.findByRole('combobox', { name: /almacenamiento de destino/i }));
    expect(await screen.findByRole('option', { name: /local-lvm/i })).toBeVisible();
  });

  it('para un Contenedor ofrece solo plantillas LXC (vztmpl) y no ISOs', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toOperatingSystemStep(user, 'lxc');

    await waitFor(() => expect(urlsCalled().some((url) => /\/node\/templates/.test(url))).toBe(true));
    expect(await screen.findByRole('radio', { name: /seleccionar debian-12-standard/i })).toBeVisible();
    expect(screen.queryByRole('radio', { name: /seleccionar ubuntu-24\.04\.iso/i })).toBeNull();
  });

  it('para una Máquina virtual ofrece solo ISOs y no plantillas LXC', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toOperatingSystemStep(user, 'vm');

    expect(await screen.findByRole('radio', { name: /seleccionar ubuntu-24\.04\.iso/i })).toBeVisible();
    expect(screen.queryByRole('radio', { name: /seleccionar debian-12-standard/i })).toBeNull();
  });

  it('la red alterna entre DHCP e IP estática, y la IP estática valida dirección, prefijo y gateway', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toNetworkStep(user, 'vm');

    const dhcp = screen.getByRole('radio', { name: /dhcp/i });
    expect(dhcp).toHaveAttribute('aria-checked', 'true');
    // En modo DHCP los campos de IP fija existen pero están deshabilitados.
    expect(screen.getByLabelText(/direcci[oó]n ip/i)).toBeDisabled();
    expect(screen.getByLabelText(/gateway|puerta de enlace/i)).toBeDisabled();

    await user.click(screen.getByRole('radio', { name: /ip est[aá]tica/i }));
    const ip = await screen.findByLabelText(/direcci[oó]n ip/i);
    await user.type(ip, 'no-es-una-ip');
    await user.type(screen.getByLabelText(/prefijo/i), '24');
    await user.type(screen.getByLabelText(/gateway|puerta de enlace/i), '192.168.1.1');

    expect(await screen.findByRole('alert')).toHaveTextContent(/ip ipv4 v[aá]lida/i);
    expect(next()).toBeDisabled();
  });

  it('el Paso 5 resume todos los parámetros y "Crear Instancia" solo se habilita al confirmar', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toReviewStep(user, 'vm');

    for (const expected of [/m[aá]quina virtual/i, /local-lvm/i, /ubuntu-24\.04\.iso/i, /dhcp/i]) {
      expect(screen.getAllByText(expected).length).toBeGreaterThan(0);
    }
    expect(screen.getByText(/confirmo la configuraci[oó]n de la instancia/i)).toBeVisible();
    expect(createButtons().every((button) => (button as HTMLButtonElement).disabled)).toBe(true);

    await user.click(screen.getByRole('checkbox'));
    await waitFor(() => expect(createButtons().every((button) => !(button as HTMLButtonElement).disabled)).toBe(true));
  });
});

