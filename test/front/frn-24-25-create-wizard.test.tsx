import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import CreateNewInstances from '@/pages/createNewInstances';
import { withProtectedProviders } from './app-providers';

// FRN-24 / FRN-25: wizard de creación de instancias en 4 pasos.
//   Paso 1 (Tipo): tarjetas Máquina Virtual (QEMU) / Contenedor (LXC) con descripción.
//   Paso 2 (Recursos): sliders o controles numéricos de vCPU, RAM y Disco + medidor de cuota reactivo
//          (advertencia sobre el 80% del remanente; el avance se deshabilita sobre el 100%).
//   Paso 3 (Almacenamiento, imagen y red): storages de GET /api/node/storages, plantillas LXC o ISOs de
//          GET /api/node/templates según el tipo, y red DHCP / IP fija con validación.
//   Paso 4 (Resumen): ficha con todos los parámetros y botón final "Crear Instancia" solo con datos completos.
// Criterio de éxito: no se avanza al Paso 3 si la configuración excede la memoria o el disco del nodo,
// y el botón final solo se habilita con los campos requeridos completos.
// La forma exacta de los controles es libre: se buscan por rol y nombre accesible.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const MIB = 1024 ** 2;
const storages = [
  { storage: 'local-lvm', type: 'lvmthin', total: 200 * 1024 * MIB, used: 20 * 1024 * MIB, avail: 180 * 1024 * MIB, active: 1, content: 'images,rootdir', shared: 0, enabled: 1 },
  { storage: 'local', type: 'dir', total: 100 * 1024 * MIB, used: 10 * 1024 * MIB, avail: 90 * 1024 * MIB, active: 1, content: 'iso,vztmpl', shared: 0, enabled: 1 },
];
const templates = [
  { volid: 'local:vztmpl/debian-12-standard_12.2-1_amd64.tar.zst', content: 'vztmpl', format: 'tzst', size: 120 * MIB },
  { volid: 'local:iso/ubuntu-24.04.iso', content: 'iso', format: 'iso', size: 2000 * MIB },
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
    if ((init.method ?? 'GET').toUpperCase() !== 'GET') return Promise.resolve(json({ vmid: 120, tareaId: 'tarea-wizard' }, 202));
    if (/\/node\/status(\?|$)/.test(target) || /\/node\/quota(\?|$)/.test(target)) return Promise.resolve(json(status));
    if (/\/node\/storages(\?|$)/.test(target)) return Promise.resolve(json(storages));
    if (/\/node\/templates(\?|$)/.test(target)) return Promise.resolve(json(templates));
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

async function chooseType(user: ReturnType<typeof userEvent.setup>, type: 'vm' | 'lxc') {
  // Las tarjetas son botones con role="button" o <button> dentro de la sección Tipo de instancia.
  const name = type === 'vm' ? /m[aá]quina virtual/i : /contenedor/i;
  const card = await screen.findByRole('button', { name });
  await user.click(card);
}

async function toStep3(user: ReturnType<typeof userEvent.setup>, type: 'vm' | 'lxc') {
  await chooseType(user, type);
  await user.click(next());
  await waitFor(() => expect(screen.queryAllByRole('slider').length + screen.queryAllByRole('spinbutton').length).toBeGreaterThanOrEqual(3));
  await user.click(next());
}

async function openAndPick(user: ReturnType<typeof userEvent.setup>, label: RegExp, option: RegExp) {
  await user.click(await screen.findByRole('combobox', { name: label }));
  await user.click(await screen.findByRole('option', { name: option }));
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-24 - Paso 1 (tipo) y Paso 2 (recursos y cuotas en vivo)', () => {
  it('el Paso 1 ofrece tarjetas de Máquina Virtual (QEMU) y Contenedor (LXC) con su descripción', async () => {
    renderWizard();
    expect(await screen.findByRole('button', { name: /m[aá]quina virtual/i })).toBeVisible();
    expect(screen.getByRole('button', { name: /contenedor/i })).toBeVisible();
    expect(screen.getByText(/qemu|kvm|sistema operativo completo/i)).toBeVisible();
    expect(screen.getByText(/lxc|liviano|ligero|compart/i)).toBeVisible();
  });

  it('el Paso 2 expone controles para vCPU, RAM y Disco', async () => {
    renderWizard();
    const user = userEvent.setup();
    await chooseType(user, 'vm');
    await user.click(next());

    await waitFor(() => {
      const controls = [...screen.queryAllByRole('slider'), ...screen.queryAllByRole('spinbutton')];
      expect(controls.length, 'faltan controles de vCPU, RAM y Disco').toBeGreaterThanOrEqual(3);
    });
    for (const label of [/vcpu|cpu|n[uú]cleos/i, /ram|memoria/i, /disco/i]) {
      expect(screen.getAllByText(label).length).toBeGreaterThan(0);
    }
  });

  it('muestra un medidor de cuota del host que reacciona al aumentar la RAM', async () => {
    renderWizard();
    const user = userEvent.setup();
    await chooseType(user, 'vm');
    await user.click(next());
    await waitFor(() => expect(urlsCalled().some((url) => /\/node\/(status|quota)/.test(url))).toBe(true));

    const meter = await waitFor(() => {
      const found = screen.queryAllByRole('progressbar').concat(screen.queryAllByRole('meter')).find((element) => /cuota|capacidad|host|ram|memoria/i.test(element.getAttribute('aria-label') ?? ''));
      expect(found, 'no hay un medidor de cuota (progressbar/meter) del host').toBeDefined();
      return found as HTMLElement;
    });
    const before = Number(meter.getAttribute('aria-valuenow'));
    await user.click(screen.getByRole('button', { name: /aumentar (memoria|ram)/i }));
    await waitFor(() => expect(Number(meter.getAttribute('aria-valuenow'))).toBeGreaterThan(before));
  });

  it('no deja avanzar al Paso 3 si la RAM pedida excede la memoria libre del nodo', async () => {
    renderWizard({ freeRamMb: 256 });
    const user = userEvent.setup();
    await chooseType(user, 'vm');
    await user.click(next());
    await waitFor(() => expect(urlsCalled().some((url) => /\/node\/(status|quota)/.test(url))).toBe(true));

    expect(await screen.findByText(/excede|supera|insuficiente|sin capacidad|cuota/i)).toBeVisible();
    const button = next();
    if (!(button as HTMLButtonElement).disabled) await user.click(button);
    expect(screen.queryByRole('combobox', { name: /almacenamiento/i }), 'avanzó al Paso 3 con la cuota excedida').toBeNull();
    expect(next()).toBeDisabled();
  });

  it('no deja avanzar al Paso 3 si el disco pedido excede el almacenamiento libre', async () => {
    renderWizard({ freeStorageGb: 1 });
    const user = userEvent.setup();
    await chooseType(user, 'vm');
    await user.click(next());
    await waitFor(() => expect(urlsCalled().some((url) => /\/node\/(status|quota)/.test(url))).toBe(true));
    await waitFor(() => expect(next()).toBeDisabled());
  });

  it('con capacidad suficiente sí permite avanzar al Paso 3', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toStep3(user, 'vm');
    expect(await screen.findByRole('combobox', { name: /almacenamiento/i })).toBeInTheDocument();
  });
});

describe('FRN-25 - Paso 3 (almacenamiento, imagen y red) y Paso 4 (resumen)', () => {
  it('lista los storages reales de GET /api/node/storages', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toStep3(user, 'vm');

    await waitFor(() => expect(urlsCalled().some((url) => /\/node\/storages/.test(url))).toBe(true));
    await user.click(await screen.findByRole('combobox', { name: /almacenamiento/i }));
    expect(await screen.findByRole('option', { name: /local-lvm/i })).toBeVisible();
  });

  it('para un Contenedor ofrece solo plantillas LXC (vztmpl) y no ISOs', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toStep3(user, 'lxc');

    await waitFor(() => expect(urlsCalled().some((url) => /\/node\/templates/.test(url))).toBe(true));
    await user.click(await screen.findByRole('combobox', { name: /plantilla|imagen|template/i }));
    expect(await screen.findByRole('option', { name: /debian-12-standard/i })).toBeVisible();
    expect(screen.queryByRole('option', { name: /ubuntu-24\.04\.iso/i })).toBeNull();
  });

  it('para una Máquina Virtual ofrece solo ISOs y no plantillas LXC', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toStep3(user, 'vm');

    await user.click(await screen.findByRole('combobox', { name: /iso|imagen|plantilla|template/i }));
    expect(await screen.findByRole('option', { name: /ubuntu-24\.04\.iso/i })).toBeVisible();
    expect(screen.queryByRole('option', { name: /debian-12-standard/i })).toBeNull();
  });

  it('la red alterna entre DHCP y IP fija, y la IP fija valida dirección, máscara y gateway', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toStep3(user, 'vm');

    const dhcp = await screen.findByRole('radio', { name: /dhcp/i }).catch(() => screen.findByRole('switch', { name: /dhcp/i }));
    expect(dhcp).toBeInTheDocument();
    expect(screen.queryByLabelText(/gateway|puerta de enlace/i)).toBeNull();

    await user.click(screen.getByRole('radio', { name: /ip fija|est[aá]tica/i }));
    const ip = await screen.findByLabelText(/direcci[oó]n ip|^ip$/i);
    await user.type(ip, 'no-es-una-ip');
    await user.type(screen.getByLabelText(/m[aá]scara|prefijo/i), '24');
    await user.type(screen.getByLabelText(/gateway|puerta de enlace/i), '192.168.1.1');
    expect(await screen.findByText(/ip.*(inv[aá]lida|v[aá]lida)|formato/i)).toBeVisible();
    expect(next()).toBeDisabled();
  });

  it('el Paso 4 resume todos los parámetros y "Crear Instancia" solo se habilita con los datos completos', async () => {
    renderWizard();
    const user = userEvent.setup();
    await toStep3(user, 'vm');

    // Sin storage ni imagen elegidos no se puede avanzar a la ficha final.
    expect(next()).toBeDisabled();

    await openAndPick(user, /almacenamiento/i, /local-lvm/i);
    await openAndPick(user, /iso|imagen|plantilla|template/i, /ubuntu-24\.04\.iso/i);
    await user.click(await screen.findByRole('radio', { name: /dhcp/i }));
    await waitFor(() => expect(next()).toBeEnabled());
    await user.click(next());

    const summary = await screen.findByRole('region', { name: /resumen/i }).catch(() => document.body);
    for (const expected of [/m[aá]quina virtual|qemu/i, /local-lvm/i, /ubuntu-24\.04\.iso/i, /dhcp/i]) {
      expect(within(summary).getAllByText(expected).length).toBeGreaterThan(0);
    }
    expect(screen.getByText(/confirm|revis[aá]|advertencia|se crear[aá]/i)).toBeVisible();
    expect(screen.getByRole('button', { name: /crear instancia/i })).toBeEnabled();
  });
});

