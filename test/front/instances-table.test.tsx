import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstancesPage from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { withProtectedProviders } from './app-providers';

// FRN-20A (ex FRN-14A): tabla de inventario con badges de estado e IP (RF-03).
// Criterio de éxito: VMs y LXC en la misma tabla; IP null → "No detectada"; tolera los campos nuevos
// en null. Entregable: columnas ID, Nombre, Tipo (VM / LXC), Estado (Running verde, Stopped gris),
// IP con botón para copiar y botonera de acciones (solo la estructura).
// El inventario respeta el contrato vigente: BAC-14 ({ id, name, type: "vm" | "lxc", node, status })
// más los campos de FIX-39 / BAC-29, que por ahora pueden venir en null.
const inventory = [
  {
    id: 101, name: 'servidor-web', type: 'vm', node: 'pve', status: 'running',
    ip: '192.168.1.50', cpuUsage: null, ramUsage: null, maxRam: null, nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
  {
    id: 102, name: 'base-datos', type: 'lxc', node: 'pve', status: 'stopped',
    ip: null, cpuUsage: null, ramUsage: null, maxRam: null, nivelAcceso: 'READ_ONLY', activeTask: null,
  },
];

let fetchMock: ReturnType<typeof vi.fn>;

function renderInstancesPage() {
  window.sessionStorage.setItem('centinela_access', 'access-token-inventario');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
    rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));
  return render(withProtectedProviders(
    <Toaster>
      <MemoryRouter initialEntries={['/instances']}>
        <InstancesPage />
      </MemoryRouter>
    </Toaster>,
  ));
}

// Fila de la tabla que contiene el nombre de la instancia.
async function rowOf(name: string) {
  const cell = await screen.findByText(name);
  const row = cell.closest('tr, [role="row"]');
  expect(row, `"${name}" no está dentro de una fila de tabla`).not.toBeNull();
  return row as HTMLElement;
}

// Clases y estilos del elemento y de sus ancestros dentro de la fila: el color del badge puede estar
// en el propio texto o en el contenedor que lo envuelve.
function visualOf(element: HTMLElement, row: HTMLElement) {
  const parts: string[] = [];
  for (let node: HTMLElement | null = element; node && node !== row; node = node.parentElement) {
    parts.push(node.className?.toString() ?? '', node.getAttribute('style') ?? '', node.getAttribute('data-variant') ?? '', node.getAttribute('data-state') ?? '');
  }
  return parts.join(' ');
}

beforeEach(() => {
  fetchMock = vi.fn().mockImplementation((url: string) => Promise.resolve(String(url).includes('/instances')
    ? new Response(JSON.stringify(inventory), { status: 200, headers: { 'Content-Type': 'application/json' } })
    : new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } })));
  vi.stubGlobal('fetch', fetchMock);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-20A - Tabla interactiva de inventario con badges de estado e IP', () => {
  it('consulta GET /api/instances con Bearer y muestra VMs y LXC en la misma tabla', async () => {
    renderInstancesPage();

    const web = await rowOf('servidor-web');
    const db = await rowOf('base-datos');
    expect(web.closest('table, [role="table"], [role="grid"]')).toBe(db.closest('table, [role="table"], [role="grid"]'));
    expect(within(web).getByText('101')).toBeInTheDocument();
    expect(within(db).getByText('102')).toBeInTheDocument();
    // Tipo: VM para type "vm" y LXC para type "lxc".
    expect(within(web).getByText(/^\s*VM\s*$/)).toBeInTheDocument();
    expect(within(db).getByText(/^\s*LXC\s*$/)).toBeInTheDocument();

    const call = (fetchMock.mock.calls as [string, RequestInit?][]).find(([url]) => /\/instances(\?|$)/.test(String(url)));
    expect(call, 'la vista no consultó GET /api/instances').toBeDefined();
    expect(new Headers(call![1]?.headers).get('Authorization')).toBe('Bearer access-token-inventario');
  });

  it('muestra badges de estado diferenciados: Running en verde y Stopped en gris', async () => {
    renderInstancesPage();

    const web = await rowOf('servidor-web');
    const db = await rowOf('base-datos');
    const running = within(web).getByText(/running|en ejecuci[oó]n|encendida/i);
    const stopped = within(db).getByText(/stopped|detenida|apagada/i);

    expect(visualOf(running, web), 'el badge Running no tiene un estilo verde').toMatch(/green|emerald|success|#16a34a|#22c55e|rgb\(22, 163, 74\)/i);
    expect(visualOf(stopped, db), 'el badge Stopped no tiene un estilo gris').toMatch(/gray|grey|slate|zinc|neutral|stone|muted|secondary/i);
  });

  it('muestra "No detectada" si la IP es null y copia la IP al portapapeles cuando existe', async () => {
    const user = userEvent.setup();
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } });

    renderInstancesPage();

    const db = await rowOf('base-datos');
    expect(within(db).getByText(/no detectada/i)).toBeInTheDocument();

    const web = await rowOf('servidor-web');
    expect(within(web).getByText('192.168.1.50')).toBeInTheDocument();
    await user.click(within(web).getByRole('button', { name: /copiar/i }));
    await waitFor(() => expect(writeText).toHaveBeenCalledWith('192.168.1.50'));
  });

  it('cada fila tiene su botonera de acciones, aunque no tenga IP para copiar', async () => {
    renderInstancesPage();

    // La 102 no tiene IP (no hay botón de copiar): cualquier botón de la fila es de la botonera.
    const db = await rowOf('base-datos');
    expect(within(db).queryAllByRole('button').length, 'la fila de la 102 no tiene botonera de acciones').toBeGreaterThan(0);
  });

  it('tolera que los campos nuevos de la Etapa 1 vengan en null sin romper la tabla', async () => {
    renderInstancesPage();

    await rowOf('servidor-web');
    // Una fila de datos por instancia: ninguna se pierde por los null.
    const table = (await rowOf('base-datos')).closest('table, [role="table"], [role="grid"]') as HTMLElement;
    const dataRows = within(table).getAllByRole('row').filter((row) => within(row).queryAllByRole('cell').length > 0);
    expect(dataRows).toHaveLength(inventory.length);
  });
});
