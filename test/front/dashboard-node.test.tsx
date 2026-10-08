import { act, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import DashboardPage from '@/pages/Dashboard';
import { withProtectedProviders } from './app-providers';

// FRN-19B (ex FRN-13B): semáforo de salud global e integración con GET /api/node/status (RF-02).
// Criterio de éxito: los datos reales del nodo se ven en pantalla y se actualizan sin F5; si el backend
// cae, la UI muestra "Inaccesible" sin romperse.
// Entregable: uptime legible; semáforo D3 (Saludable < 70 %, Advertencia si alguno llega al 70 %,
// Inaccesible ante 502/504 o sin respuesta); aviso si stale: true; consulta cada 10 s que se pausa con la
// pestaña oculta; skeletons mientras carga. El payload es el del contrato de BAC-29.

const status = (overrides: Record<string, unknown> = {}) => ({
  cpu: { usagePercent: 69, cores: 8 },
  ram: { usedGb: 8, totalGb: 16, usagePercent: 50 },
  storage: { usedGb: 30, totalGb: 100, usagePercent: 30 },
  uptimeSeconds: 266400, // 3 días y 2 horas
  instancesSummary: { vms: { running: 2, stopped: 1, paused: 0, total: 3 }, lxc: { running: 1, stopped: 0, paused: 0, total: 1 } },
  stale: false,
  fetchedAt: '2026-10-04T03:00:00Z',
  ...overrides,
});

const json = (body: unknown, code = 200) => new Response(JSON.stringify(body), { status: code, headers: { 'Content-Type': 'application/json' } });

function stubNodeStatus(respond: () => Promise<Response>) {
  const fetchMock = vi.fn().mockImplementation((url: string) =>
    /\/node\/status(\?|$)/.test(String(url)) ? respond() : Promise.resolve(json({})));
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

const nodeCalls = (fetchMock: ReturnType<typeof vi.fn>) =>
  (fetchMock.mock.calls as [string, RequestInit?][]).filter(([url]) => /\/node\/status(\?|$)/.test(String(url)));

function renderDashboard() {
  window.sessionStorage.setItem('centinela_access', 'access-token-dashboard');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local', rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));
  return render(withProtectedProviders(<MemoryRouter initialEntries={['/dashboard']}><DashboardPage /></MemoryRouter>));
}

function setVisibility(state: 'visible' | 'hidden') {
  Object.defineProperty(document, 'visibilityState', { configurable: true, get: () => state });
  Object.defineProperty(document, 'hidden', { configurable: true, get: () => state === 'hidden' });
  document.dispatchEvent(new Event('visibilitychange'));
}

afterEach(() => {
  vi.unstubAllGlobals();
  vi.useRealTimers();
  setVisibility('visible');
});

describe('FRN-19B - Semáforo de salud global e integración con GET /api/node/status', () => {
  it('consulta GET /api/node/status con Bearer y muestra los datos reales del nodo', async () => {
    const fetchMock = stubNodeStatus(() => Promise.resolve(json(status({ cpu: { usagePercent: 37, cores: 12 } }))));
    renderDashboard();

    await waitFor(() => expect(nodeCalls(fetchMock).length).toBeGreaterThan(0));
    expect(new Headers(nodeCalls(fetchMock)[0][1]?.headers).get('Authorization')).toBe('Bearer access-token-dashboard');
    // Valores que solo pueden venir de la respuesta (no de una maqueta).
    expect(await screen.findByText(/\b37(?:[.,]0+)?\s*%/)).toBeInTheDocument();
    expect(screen.getByText(/\b12\b/)).toBeInTheDocument();
  });

  it('D3: con todo por debajo del 70 % (borde 69 %) el semáforo dice Saludable', async () => {
    stubNodeStatus(() => Promise.resolve(json(status())));
    renderDashboard();
    expect(await screen.findByText(/saludable/i)).toBeInTheDocument();
    expect(screen.queryByText(/^\s*inaccesible\s*$/i)).not.toBeInTheDocument();
  });

  it.each([
    ['CPU', { cpu: { usagePercent: 70, cores: 8 } }],
    ['RAM', { ram: { usedGb: 11.2, totalGb: 16, usagePercent: 70 } }],
    ['disco', { storage: { usedGb: 70, totalGb: 100, usagePercent: 70 } }],
  ])('D3: si %s llega al 70 % el semáforo pasa a Advertencia', async (_, overrides) => {
    stubNodeStatus(() => Promise.resolve(json(status({ cpu: { usagePercent: 20, cores: 8 }, ...overrides }))));
    renderDashboard();
    expect(await screen.findAllByText(/advertencia/i)).not.toHaveLength(0);
    expect(screen.queryByText(/saludable/i), 'con un recurso al 70 % el semáforo sigue diciendo Saludable').not.toBeInTheDocument();
  });

  it.each([
    ['502 PROXMOX_UNAVAILABLE', () => Promise.resolve(json({ errorCode: 'PROXMOX_UNAVAILABLE', message: 'x' }, 502))],
    ['504 PROXMOX_TIMEOUT', () => Promise.resolve(json({ errorCode: 'PROXMOX_TIMEOUT', message: 'x' }, 504))],
    ['sin respuesta (error de red)', () => Promise.reject(new TypeError('Failed to fetch'))],
  ])('si el backend responde %s el semáforo dice Inaccesible y la página no se rompe', async (_, respond) => {
    stubNodeStatus(respond as () => Promise<Response>);
    renderDashboard();
    expect(await screen.findByText(/inaccesible/i)).toBeInTheDocument();
    expect(screen.queryByText(/saludable/i)).not.toBeInTheDocument();
    expect(screen.getAllByRole('heading').length, 'la página dejó de renderizarse').toBeGreaterThan(0);
  });

  it('muestra el uptime legible (3 días, 2 horas) y avisa solo cuando los datos están desactualizados', async () => {
    let stale = false;
    stubNodeStatus(() => Promise.resolve(json(status({ stale }))));
    const { unmount } = renderDashboard();
    expect(await screen.findByText(/3\s*d(í|i)as?\b.*2\s*h|3\s*d\b.*2\s*h/i)).toBeInTheDocument();
    expect(screen.queryByText(/desactualizad/i), 'con stale: false no debe haber aviso').not.toBeInTheDocument();
    unmount();

    stale = true;
    renderDashboard();
    expect(await screen.findByText(/desactualizad/i)).toBeInTheDocument();
  });

  it('muestra un estado de carga mientras espera /node/status y lo saca cuando llega la respuesta', async () => {
    let release!: (response: Response) => void;
    const fetchMock = stubNodeStatus(() => new Promise<Response>((resolve) => { release = resolve; }));
    const { container } = renderDashboard();
    const loading = () => container.querySelector('[aria-busy="true"], .animate-pulse, [data-loading="true"], [data-slot="skeleton"]');

    await waitFor(() => expect(nodeCalls(fetchMock).length, 'el Dashboard no consultó /node/status').toBeGreaterThan(0));
    expect(loading(), 'no hay skeleton ni estado de carga mientras se espera /node/status').not.toBeNull();

    release(json(status()));
    expect(await screen.findByText(/saludable/i)).toBeInTheDocument();
    await waitFor(() => expect(loading(), 'el estado de carga sigue visible después de recibir los datos').toBeNull());
  });

  it('se actualiza sola cada 10 s y deja de consultar con la pestaña oculta', async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    let cpu = 31;
    const fetchMock = stubNodeStatus(() => Promise.resolve(json(status({ cpu: { usagePercent: cpu, cores: 8 } }))));
    renderDashboard();
    expect(await screen.findByText(/\b31(?:[.,]0+)?\s*%/)).toBeInTheDocument();

    cpu = 42;
    await act(async () => { await vi.advanceTimersByTimeAsync(10_500); });
    expect(await screen.findByText(/\b42(?:[.,]0+)?\s*%/), 'a los 10 s no se volvió a consultar ni se actualizó el valor').toBeInTheDocument();

    setVisibility('hidden');
    const before = nodeCalls(fetchMock).length;
    await act(async () => { await vi.advanceTimersByTimeAsync(35_000); });
    expect(nodeCalls(fetchMock).length, 'con la pestaña oculta siguió consultando /node/status').toBe(before);

    setVisibility('visible');
    await act(async () => { await vi.advanceTimersByTimeAsync(10_500); });
    expect(nodeCalls(fetchMock).length, 'al volver a la pestaña no retomó las consultas').toBeGreaterThan(before);
  });

  it('FIX-60: el Dashboard no renderiza el botón "Crear instancia"', async () => {
    stubNodeStatus(() => Promise.resolve(json(status())));
    renderDashboard();
    await screen.findByText(/saludable/i);
    expect(screen.queryByRole('button', { name: /crear instancia|nueva instancia/i })).toBeNull();
  });
});
