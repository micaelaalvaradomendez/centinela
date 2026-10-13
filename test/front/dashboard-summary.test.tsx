import { act, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DashboardPage from '@/pages/Dashboard';
import { withProtectedProviders } from './app-providers';

// FRN-19C (BRG-05-FRN, RF-02): tarjetas de conteo de VMs/LXC en el Dashboard.
// Entregable: "En ejecución", "Detenidas o pausadas" y "Total" con los datos de
// `instancesSummary` que devuelve GET /api/node/status (BAC-22B / BAC-29).
// Criterio de éxito: al cambiar el estado de una instancia (evento TASK_FINISHED),
// el conteo se actualiza sin F5 en menos de 2 s.

const summaryBase = {
  vms: { running: 2, stopped: 1, paused: 0, total: 3 },
  lxc: { running: 1, stopped: 0, paused: 0, total: 1 },
};

const status = (summary = summaryBase) => ({
  cpu: { usagePercent: 40, cores: 8 },
  ram: { usedGb: 8, totalGb: 16, usagePercent: 50 },
  storage: { usedGb: 30, totalGb: 100, usagePercent: 30 },
  uptimeSeconds: 266400,
  instancesSummary: summary,
  stale: false,
  fetchedAt: '2026-10-13T03:00:00Z',
});

const json = (body: unknown, code = 200) =>
  new Response(JSON.stringify(body), { status: code, headers: { 'Content-Type': 'application/json' } });

// Conexión SSE simulada: permite emitir un evento TASK_FINISHED como lo haría el backend.
class FakeSse {
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onopen: ((event: Event) => void) | null = null;
  readyState = 1;
  private listeners: Record<string, ((event: Event) => void)[]> = {};
  constructor(public url: string) {}
  addEventListener(type: string, listener: (event: Event) => void) { (this.listeners[type] ??= []).push(listener); }
  removeEventListener() {}
  close() { this.readyState = 2; }
  emitMessage(data: unknown) {
    const event = new MessageEvent('message', { data: typeof data === 'string' ? data : JSON.stringify(data) });
    this.onmessage?.(event);
    this.listeners.message?.forEach((listener) => listener(event));
  }
}

let latestSse: FakeSse | null = null;

beforeEach(() => {
  latestSse = null;
  // El doble captura la instancia para poder emitirle eventos desde el test.
  class CapturingSse extends FakeSse {
    constructor(url: string) { super(url); latestSse = this; }
  }
  vi.stubGlobal('EventSource', CapturingSse);
  vi.stubGlobal('WebSocket', CapturingSse);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

function renderDashboard(respondStatus: () => Promise<Response>) {
  const fetchMock = vi.fn().mockImplementation((url: string) => {
    if (/\/node\/status(\?|$)/.test(String(url))) return respondStatus();
    if (String(url).includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-dashboard' }));
    return Promise.resolve(json({}));
  });
  vi.stubGlobal('fetch', fetchMock);
  window.sessionStorage.setItem('centinela_access', 'access-token-resumen');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
    rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));
  const utils = render(withProtectedProviders(
    <MemoryRouter initialEntries={['/dashboard']}><DashboardPage /></MemoryRouter>,
  ));
  return { fetchMock, ...utils };
}

const nodeCalls = (fetchMock: ReturnType<typeof vi.fn>) =>
  (fetchMock.mock.calls as [string, RequestInit?][]).filter(([url]) => /\/node\/status(\?|$)/.test(String(url)));

// El valor de la tarjeta vive en el mismo contenedor que su título.
// Se ubica por el encabezado (h2) para no confundirlo con el pie de la tarjeta,
// que repite "detenidas o pausadas".
function cardValue(title: RegExp): string {
  const heading = screen.getByRole('heading', { name: title });
  return heading.closest('div')?.textContent ?? '';
}

// La métrica es el primer número del contenedor (los títulos son texto).
function leadNumber(value: string): string | undefined {
  return value.match(/^\s*(\d+)/)?.[1];
}

describe('FRN-19C - Tarjetas de conteo de instancias con auto-actualización', () => {
  it('muestra el total, en ejecución y detenidas/pausadas desde instancesSummary', async () => {
    const { fetchMock } = renderDashboard(() => Promise.resolve(json(status())));

    await waitFor(() => expect(nodeCalls(fetchMock).length).toBeGreaterThan(0));

    // Total = 3 VMs + 1 LXC = 4; en ejecución = 2 + 1 = 3; inactivas = 1 + 0 = 1.
    // El contenedor de la tarjeta concatena la métrica y su título ("4Instancias totales").
    expect(leadNumber(cardValue(/instancias totales/i))).toBe('4');
    expect(leadNumber(cardValue(/en ejecuci[oó]n/i))).toBe('3');
    expect(leadNumber(cardValue(/detenidas o pausadas/i))).toBe('1');
  });

  it('vuelve a consultar /node/status al recibir TASK_FINISHED para actualizar el conteo sin F5', async () => {
    const { fetchMock } = renderDashboard(() => Promise.resolve(json(status())));
    await waitFor(() => expect(latestSse, 'no se abrió la conexión de eventos').not.toBeNull());
    await waitFor(() => expect(nodeCalls(fetchMock).length).toBeGreaterThan(0));

    const before = nodeCalls(fetchMock).length;
    act(() => {
      latestSse!.emitMessage({
        id: 'ev-1', tipo: 'TASK_FINISHED', severidad: 'INFO', recursoTipo: 'VM', recursoId: '101',
        mensaje: 'finalizó', fechaHora: '2026-10-13T03:00:00Z',
        detalles: { tareaId: 't-1', accion: 'start', estado: 'COMPLETED', exitstatus: 'OK' },
      });
    });

    // Criterio de éxito de FRN-19C: el conteo se refresca por el evento, no solo por el polling de 10 s.
    await waitFor(() => expect(nodeCalls(fetchMock).length).toBeGreaterThan(before), { timeout: 2000 });
  });
});
