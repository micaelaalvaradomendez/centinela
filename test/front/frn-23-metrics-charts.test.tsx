import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstanceDetailPage from '@/pages/InstanceDetailPage';
import { withProtectedProviders } from './app-providers';

// FRN-23: pestaña "Métricas" de InstanceDetailPage con gráficos interactivos y selector de rango.
// Entregables:
//   1. Gráficos de CPU (%), Memoria RAM (usada vs total), Tráfico de red (KB/s in/out) e I/O de disco (KB/s).
//   2. Selector de rango: "En vivo" (WebSocket), "Últimas 2 horas" y "Últimas 24 horas" (REST BAC-36).
//   3. Formateo legible de unidades (bytes a MB/GB, bps a KB/s).
//   4. Tooltips explicativos y estados de carga (skeletons).
//   5. Tolerancia a nulos: con la máquina detenida se muestra "Instancia detenida" sin NaN ni TypeError.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const GIB = 1024 ** 3;
const inventoryRunning = [
  {
    id: 101, name: 'servidor-web', type: 'vm', node: 'nodo-a1', status: 'running', ip: '192.168.1.50',
    cpuUsage: 25, ramUsage: GIB, maxRam: 4 * GIB, nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
  {
    id: 'inst_k8Pq2v', name: 'servidor-web-opaco', type: 'vm', node: 'nodo-a1', status: 'running', ip: '192.168.1.51',
    cpuUsage: 25, ramUsage: GIB, maxRam: 4 * GIB, nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
];
const inventoryStopped = inventoryRunning.map((inst) => ({ ...inst, status: 'stopped', cpuUsage: null, ramUsage: null }));

const historySeries = [
  { timestamp: 1_760_000_000, cpuPercent: 10, ramUsedBytes: GIB, ramTotalBytes: 4 * GIB, netInBps: 2048, netOutBps: 1024, diskReadBps: 4096, diskWriteBps: 2048 },
  { timestamp: 1_760_000_060, cpuPercent: 30, ramUsedBytes: 1.5 * GIB, ramTotalBytes: 4 * GIB, netInBps: 4096, netOutBps: 2048, diskReadBps: 8192, diskWriteBps: 4096 },
];
const stoppedSeries = historySeries.map((point) => ({ ...point, cpuPercent: null, ramUsedBytes: null }));

class SilentSocket {
  onopen: unknown = null; onmessage: unknown = null; onclose: unknown = null; onerror: unknown = null;
  readyState = 0;
  constructor(public url: string) {}
  addEventListener() {}
  removeEventListener() {}
  send() {}
  close() { this.readyState = 3; }
}

class FakeResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

let fetchMock: ReturnType<typeof vi.fn>;
let consoleError: ReturnType<typeof vi.spyOn>;

function stubBackend(options: { inventory?: unknown[]; series?: unknown[]; metricsPending?: boolean; stopped?: boolean } = {}) {
  const isStopped = options.stopped ?? (options.inventory === inventoryStopped);
  const seriesPayload = options.series ?? historySeries;
  fetchMock = vi.fn().mockImplementation((url: string) => {
    const target = String(url);
    if (/\/instances\/(101|inst_k8Pq2v)\/metrics/.test(target)) {
      return options.metricsPending
        ? new Promise<Response>(() => undefined)
        : Promise.resolve(json({ data: seriesPayload }));
    }
    if (/\/instances\/(101|inst_k8Pq2v)(\?|$)/.test(target)) {
      const ref = target.includes('inst_k8Pq2v') ? 'inst_k8Pq2v' : '101';
      return Promise.resolve(json({
        vmid: ref,
        nombre: 'servidor-web',
        nodo: 'nodo-a1',
        estado: isStopped ? 'stopped' : 'running',
        tipo: 'qemu',
        cpu: isStopped ? 0 : 25,
        mem: isStopped ? 0 : GIB,
        maxCpu: 2,
        maxMem: 4 * GIB,
      }));
    }
    if (/\/instances(\?|$)/.test(target)) return Promise.resolve(json(options.inventory ?? inventoryRunning));
    if (target.includes('/auth/ws-ticket') || target.includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-frn23' }));
    return Promise.resolve(json({}));
  });
  vi.stubGlobal('fetch', fetchMock);
}

function renderDetail(instanceRef = '101') {
  window.sessionStorage.setItem('centinela_access', 'access-token-frn23');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'admin-1', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
    rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));
  return render(withProtectedProviders(
    <MemoryRouter initialEntries={[`/instances/${instanceRef}`]}>
      <Routes><Route path="/instances/:vmid" element={<InstanceDetailPage />} /></Routes>
    </MemoryRouter>,
  ));
}

const metricsCalls = (ref = '101') =>
  (fetchMock.mock.calls as [string][]).map(([url]) => String(url)).filter((url) => new RegExp(`/instances/${ref}/metrics`).test(url));

// El selector puede ser tabs, botones, radios o un select: se busca por rol y nombre accesible.
async function findRangeControl(name: RegExp) {
  return waitFor(() => {
    const found = [
      ...screen.queryAllByRole('tab', { name }), ...screen.queryAllByRole('button', { name }),
      ...screen.queryAllByRole('radio', { name }), ...screen.queryAllByRole('option', { name }),
    ];
    expect(found.length, `falta el control de rango ${name}`).toBeGreaterThan(0);
    return found[0];
  }, { timeout: 3000 });
}

async function openMetricsTab(user: ReturnType<typeof userEvent.setup>) {
  const tab = await waitFor(() => {
    const found = screen.queryByRole('tab', { name: /m[eé]tricas/i }) ?? screen.queryByRole('button', { name: /m[eé]tricas/i });
    expect(found, 'la vista de detalle no tiene la pestaña "Métricas"').not.toBeNull();
    return found as HTMLElement;
  }, { timeout: 3000 });
  await user.click(tab);
}

beforeEach(() => {
  vi.stubGlobal('WebSocket', SilentSocket);
  vi.stubGlobal('EventSource', SilentSocket);
  vi.stubGlobal('ResizeObserver', FakeResizeObserver);
  consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-23 - gráficos de rendimiento y selector de rango temporal', () => {
  it('ofrece los controles "En vivo", "Últimas 2 horas" y "Últimas 24 horas"', async () => {
    stubBackend();
    renderDetail();
    const user = userEvent.setup();
    await openMetricsTab(user);

    await findRangeControl(/en vivo/i);
    await findRangeControl(/[uú]ltimas 2 horas|2 ?h\b/i);
    await findRangeControl(/[uú]ltimas 24 horas|24 ?h\b/i);
  });

  it('muestra los cuatro gráficos: CPU, memoria RAM, tráfico de red e I/O de disco', async () => {
    stubBackend();
    renderDetail();
    const user = userEvent.setup();
    await openMetricsTab(user);

    await user.click(await findRangeControl(/[uú]ltimas 2 horas|2 ?h\b/i));
    for (const title of [/cpu/i, /memoria|ram/i, /tr[aá]fico de red|red\b/i, /disco|i\/o/i]) {
      await waitFor(() => {
        const found = [...screen.queryAllByRole('heading', { name: title }), ...screen.queryAllByRole('img', { name: title }), ...screen.queryAllByRole('figure', { name: title })];
        expect(found.length, `falta el gráfico ${title}`).toBeGreaterThan(0);
      });
    }
  });

  it('los rangos históricos consultan GET /instances/101/metrics?timeframe=2h|24h (BAC-36)', async () => {
    stubBackend();
    renderDetail();
    const user = userEvent.setup();
    await openMetricsTab(user);

    await user.click(await findRangeControl(/[uú]ltimas 24 horas|24 ?h\b/i));
    await waitFor(() => expect(metricsCalls().some((url) => /timeframe=24h/.test(url))).toBe(true));

    await user.click(await findRangeControl(/[uú]ltimas 2 horas|2 ?h\b/i));
    await waitFor(() => expect(metricsCalls().some((url) => /timeframe=2h/.test(url))).toBe(true));
  });

  it('"En vivo" abre el canal WebSocket y no consulta el histórico REST', async () => {
    stubBackend();
    const opened: string[] = [];
    class TrackingSocket extends SilentSocket { constructor(url: string) { super(url); opened.push(url); } }
    vi.stubGlobal('WebSocket', TrackingSocket);
    renderDetail();
    const user = userEvent.setup();
    await openMetricsTab(user);

    await user.click(await findRangeControl(/en vivo/i));
    await waitFor(() => expect(opened.some((url) => /\/ws\/instances\/101\/metrics/.test(url))).toBe(true));
    expect(metricsCalls()).toHaveLength(0);
  });

  it('muestra un skeleton o estado de carga mientras llegan los primeros datos', async () => {
    stubBackend({ metricsPending: true });
    const { container } = renderDetail();
    const user = userEvent.setup();
    await openMetricsTab(user);
    await user.click(await findRangeControl(/[uú]ltimas 2 horas|2 ?h\b/i));

    await waitFor(() => {
      const loading = container.querySelector('[aria-busy="true"], [data-slot="skeleton"], .animate-pulse, [role="progressbar"]')
        ?? screen.queryByText(/cargando/i);
      expect(loading, 'no hay indicador de carga mientras se obtiene el histórico').not.toBeNull();
    });
  });

  it('formatea unidades legibles: bytes a MB/GB y bps a KB/s', async () => {
    stubBackend();
    renderDetail();
    const user = userEvent.setup();
    await openMetricsTab(user);
    await user.click(await findRangeControl(/[uú]ltimas 2 horas|2 ?h\b/i));

    await waitFor(() => {
      const text = document.body.textContent ?? '';
      expect(text, 'falta la unidad KB/s en red o disco').toMatch(/KB\/s/);
      expect(text, 'la memoria no se muestra en MB/GB').toMatch(/\b\d+([.,]\d+)?\s?(GB|GiB|MB|MiB)\b/);
      expect(text).not.toMatch(/\b\d{7,}\s?(bps|B\/s)\b/);
    });
  });

  it('con la máquina detenida (cpu/ram null) muestra "Instancia detenida" sin NaN ni errores', async () => {
    stubBackend({ inventory: inventoryStopped, series: stoppedSeries });
    renderDetail();
    const user = userEvent.setup();
    await openMetricsTab(user);
    await user.click(await findRangeControl(/[uú]ltimas 2 horas|2 ?h\b/i));

    expect(await screen.findByText(/instancia detenida/i)).toBeVisible();
    expect(document.body.textContent ?? '').not.toMatch(/NaN/);
    const messages = consoleError.mock.calls.map((call) => call.map(String).join(' ')).join('\n');
    expect(messages).not.toMatch(/TypeError|NaN|toFixed/);
  });

  it('cada gráfico explica sus valores con un tooltip o descripción accesible', async () => {
    stubBackend();
    renderDetail();
    const user = userEvent.setup();
    await openMetricsTab(user);
    await user.click(await findRangeControl(/[uú]ltimas 2 horas|2 ?h\b/i));

    await waitFor(() => {
      const described = document.querySelectorAll('[title], [aria-describedby], [data-slot="tooltip-trigger"], [role="tooltip"]');
      expect(described.length, 'los gráficos no tienen tooltips explicativos').toBeGreaterThan(0);
    });
    const panel = screen.getByRole('tabpanel', { hidden: true });
    expect(within(panel).queryAllByRole('img').length + within(panel).queryAllByRole('figure').length).toBeGreaterThan(0);
  });

  it('soporta referencias de instancia ofuscadas (FIX-107, ej. inst_k8Pq2v) al consultar el detalle y las métricas', async () => {
    stubBackend();
    renderDetail('inst_k8Pq2v');

    await waitFor(() => {
      expect(metricsCalls('inst_k8Pq2v').some((url) => /\/instances\/inst_k8Pq2v\/metrics\?timeframe=2h/.test(url))).toBe(true);
    });
  });
});

