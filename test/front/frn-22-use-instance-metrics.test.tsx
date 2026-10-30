import { act, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

// FRN-22: hook `useInstanceMetrics(vmid)` y cliente WebSocket de telemetría en tiempo real.
// Entregables:
//   1. Servicio de ticket efímero `getWsTicket()` (POST /auth/ws-ticket).
//   2. Hook que gestiona la conexión, su estado (conectando, conectado, reconectando, error) y un
//      búfer de los últimos N puntos de telemetría.
//   3. Reintento con backoff exponencial ante una desconexión accidental.
//   4. Cierre explícito del socket al desmontar.
//   5. Contrato de telemetría: CPU en escala 0-100 (sin multiplicar por 100) y `cpu: null` / `mem: null`
//      con la instancia apagada, sin invocar `.toFixed()` sobre null.
//
// El nombre de los campos que devuelve el hook no está fijado por la tarea: se aceptan los alias
// habituales (`status|estado`, `points|puntos|series|buffer|data`).

type HookModule = Record<string, unknown>;
const hookFiles = import.meta.glob<HookModule>('@/**/useInstanceMetrics.{ts,tsx,js,jsx}');
const sources = import.meta.glob<string>('@/**/*.{ts,tsx,js,jsx}', { query: '?raw', import: 'default', eager: true });

async function loadHook(): Promise<(vmid: number | string) => Record<string, unknown>> {
  const loaders = Object.values(hookFiles);
  expect(loaders.length, 'no existe el hook useInstanceMetrics en el frontend (FRN-22)').toBeGreaterThan(0);
  const module = await loaders[0]();
  const hook = (module.useInstanceMetrics ?? module.default) as ((vmid: number | string) => Record<string, unknown>) | undefined;
  expect(typeof hook, 'useInstanceMetrics debe exportarse (nombrado o por defecto)').toBe('function');
  return hook!;
}

class FakeWebSocket {
  static instances: FakeWebSocket[] = [];
  static OPEN = 1; static CLOSED = 3; static CONNECTING = 0; static CLOSING = 2;
  readyState = 0;
  onopen: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  closeCalls = 0;
  private listeners: Record<string, ((event: Event) => void)[]> = {};
  constructor(public url: string) { FakeWebSocket.instances.push(this); }
  addEventListener(type: string, listener: (event: Event) => void) { (this.listeners[type] ??= []).push(listener); }
  removeEventListener(type: string, listener: (event: Event) => void) { this.listeners[type] = (this.listeners[type] ?? []).filter((l) => l !== listener); }
  send() {}
  close() { this.closeCalls += 1; this.readyState = 3; }
  private emit(type: string, event: Event) {
    (this as unknown as Record<string, ((e: Event) => void) | null>)[`on${type}`]?.(event);
    (this.listeners[type] ?? []).forEach((listener) => listener(event));
  }
  serverOpen() { this.readyState = 1; this.emit('open', new Event('open')); }
  serverMessage(payload: unknown) { this.emit('message', new MessageEvent('message', { data: JSON.stringify(payload) })); }
  serverClose(code = 1006, reason = '') { this.readyState = 3; this.emit('close', new CloseEvent('close', { code, reason })); }
}

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

let fetchMock: ReturnType<typeof vi.fn>;
let ticketCounter = 0;

const sample = (cpu: number | null, mem: number | null = 1_000_000_000) => ({
  timestamp: Date.now(), vmid: 101, status: cpu === null ? 'stopped' : 'running', cpu, cpus: 2, mem, maxmem: 4_000_000_000,
  disk: 0, maxdisk: 10_000_000_000, netin: 2048, netout: 1024, diskread: 512, diskwrite: 256, uptime: 120,
});

const statusOf = (result: Record<string, unknown>) => String(result.status ?? result.estado ?? result.connectionState ?? '');
const pointsOf = (result: Record<string, unknown>): Record<string, unknown>[] => {
  for (const key of ['points', 'puntos', 'series', 'buffer', 'data', 'metrics', 'history']) {
    if (Array.isArray(result[key])) return result[key] as Record<string, unknown>[];
  }
  return [];
};
const cpuOf = (point: Record<string, unknown> | undefined) => (point ? (point.cpu ?? point.cpuPercent) : undefined);

const flush = async (ms = 0) => { await act(async () => { await vi.advanceTimersByTimeAsync(ms); }); };

beforeEach(() => {
  vi.useFakeTimers();
  FakeWebSocket.instances = [];
  ticketCounter = 0;
  window.sessionStorage.setItem('centinela_access', 'access-token-frn22');
  fetchMock = vi.fn().mockImplementation((url: string) => {
    if (String(url).includes('/auth/ws-ticket')) return Promise.resolve(json({ ticket: `ticket-${++ticketCounter}` }));
    return Promise.resolve(json({}));
  });
  vi.stubGlobal('fetch', fetchMock);
  vi.stubGlobal('WebSocket', FakeWebSocket);
});

afterEach(() => {
  vi.useRealTimers();
  vi.unstubAllGlobals();
});

describe('FRN-22 - useInstanceMetrics y cliente WebSocket', () => {
  it('existe getWsTicket() que solicita POST /auth/ws-ticket', () => {
    const exported = Object.values(sources).some((source) => /export\s+(async\s+)?(function|const)\s+getWsTicket\b/.test(source));
    expect(exported, 'no hay un servicio exportado getWsTicket() en el frontend').toBe(true);
    const callsEndpoint = Object.values(sources).some((source) => /getWsTicket[\s\S]{0,400}\/auth\/ws-ticket|\/auth\/ws-ticket[\s\S]{0,400}getWsTicket/.test(source));
    expect(callsEndpoint, 'getWsTicket no consulta /auth/ws-ticket').toBe(true);
  });

  it('pide el ticket antes de abrir el socket, lo pasa en la URL y recorre conectando -> conectado', async () => {
    const useInstanceMetrics = await loadHook();
    const { result } = renderHook(() => useInstanceMetrics(101));

    expect(statusOf(result.current)).toMatch(/conectando|connecting/i);
    await flush();

    const ticketCalls = (fetchMock.mock.calls as [string, RequestInit][]).filter(([url]) => String(url).includes('/auth/ws-ticket'));
    expect(ticketCalls).toHaveLength(1);
    expect(ticketCalls[0][1].method).toBe('POST');
    expect(FakeWebSocket.instances).toHaveLength(1);
    expect(FakeWebSocket.instances[0].url).toMatch(/\/ws\/instances\/101\/metrics\?.*ticket=ticket-1/);

    act(() => FakeWebSocket.instances[0].serverOpen());
    expect(statusOf(result.current)).toMatch(/conectado|connected|open/i);
  });

  it('acumula los puntos en un búfer acotado y conserva la CPU en escala 0-100', async () => {
    const useInstanceMetrics = await loadHook();
    const { result } = renderHook(() => useInstanceMetrics(101));
    await flush();
    const socket = FakeWebSocket.instances[0];
    act(() => socket.serverOpen());

    act(() => socket.serverMessage(sample(42.5)));
    expect(pointsOf(result.current)).toHaveLength(1);
    expect(cpuOf(pointsOf(result.current)[0]), 'la CPU debe quedar en escala 0-100, sin multiplicar por 100').toBe(42.5);

    act(() => { for (let i = 0; i < 2000; i += 1) socket.serverMessage(sample(i % 100)); });
    const points = pointsOf(result.current);
    expect(points.length, 'el búfer crece sin límite').toBeLessThan(2000);
    expect(points.length).toBeGreaterThan(10);
    expect(cpuOf(points[points.length - 1])).toBe(1999 % 100);
  });

  it('tolera cpu: null y mem: null (instancia apagada) sin excepciones ni NaN', async () => {
    const useInstanceMetrics = await loadHook();
    const { result } = renderHook(() => useInstanceMetrics(101));
    await flush();
    const socket = FakeWebSocket.instances[0];
    act(() => socket.serverOpen());

    expect(() => act(() => socket.serverMessage(sample(null, null)))).not.toThrow();
    const points = pointsOf(result.current);
    expect(points.length).toBeGreaterThan(0);
    const last = points[points.length - 1];
    expect(cpuOf(last) ?? null).toBeNull();
    const numeric = Object.values(last).filter((value): value is number => typeof value === 'number');
    expect(numeric.some((value) => Number.isNaN(value)), 'un punto de instancia apagada contiene NaN').toBe(false);
  });

  it('reintenta con backoff exponencial ante una desconexión accidental y pide un ticket nuevo', async () => {
    const useInstanceMetrics = await loadHook();
    const { result } = renderHook(() => useInstanceMetrics(101));
    await flush();

    const delays: number[] = [];
    for (let attempt = 0; attempt < 3; attempt += 1) {
      const socket = FakeWebSocket.instances[FakeWebSocket.instances.length - 1];
      act(() => socket.serverOpen());
      const before = FakeWebSocket.instances.length;
      act(() => socket.serverClose(1006));
      expect(statusOf(result.current)).toMatch(/reconectando|reconnecting/i);

      let waited = 0;
      while (FakeWebSocket.instances.length === before && waited < 120_000) {
        await flush(100);
        waited += 100;
      }
      expect(FakeWebSocket.instances.length, `no hubo reintento tras la caída ${attempt + 1}`).toBe(before + 1);
      delays.push(waited);
    }
    // Sin abrir el nuevo socket (no hay reinicio del contador), cada espera crece.
    expect(delays[1], `esperas ${delays.join(', ')} ms: no crecen`).toBeGreaterThan(delays[0]);
    expect(delays[2]).toBeGreaterThan(delays[1]);
    const tickets = (fetchMock.mock.calls as [string][]).filter(([url]) => String(url).includes('/auth/ws-ticket'));
    expect(tickets.length, 'cada intento debe usar un ticket nuevo (son de un solo uso)').toBeGreaterThanOrEqual(4);
  });

  it('no reintenta ante un cierre terminal 1008 (INSTANCE_ACCESS_DENIED) y pasa a estado error', async () => {
    const useInstanceMetrics = await loadHook();
    const { result } = renderHook(() => useInstanceMetrics(101));
    await flush();
    const socket = FakeWebSocket.instances[0];
    act(() => socket.serverOpen());
    act(() => socket.serverClose(1008, 'INSTANCE_ACCESS_DENIED'));
    await flush(120_000);

    expect(FakeWebSocket.instances).toHaveLength(1);
    expect(statusOf(result.current)).toMatch(/error/i);
  });

  it('cierra explícitamente el socket al desmontar y no deja reintentos pendientes', async () => {
    const useInstanceMetrics = await loadHook();
    const { unmount } = renderHook(() => useInstanceMetrics(101));
    await flush();
    const socket = FakeWebSocket.instances[0];
    act(() => socket.serverOpen());

    unmount();
    expect(socket.closeCalls, 'el socket no se cerró al desmontar').toBeGreaterThan(0);

    const ticketsBefore = (fetchMock.mock.calls as [string][]).filter(([url]) => String(url).includes('/auth/ws-ticket')).length;
    act(() => socket.serverClose(1006));
    await flush(120_000);
    const ticketsAfter = (fetchMock.mock.calls as [string][]).filter(([url]) => String(url).includes('/auth/ws-ticket')).length;
    expect(ticketsAfter, 'se reintentó la conexión después de desmontar').toBe(ticketsBefore);
    expect(FakeWebSocket.instances).toHaveLength(1);
  });

  it('al cambiar de instancia cierra el socket anterior y abre uno para la nueva', async () => {
    const useInstanceMetrics = await loadHook();
    const { rerender } = renderHook(({ vmid }) => useInstanceMetrics(vmid), { initialProps: { vmid: 101 } });
    await flush();
    const first = FakeWebSocket.instances[0];
    act(() => first.serverOpen());

    rerender({ vmid: 102 });
    await flush();
    expect(first.closeCalls).toBeGreaterThan(0);
    expect(FakeWebSocket.instances).toHaveLength(2);
    expect(FakeWebSocket.instances[1].url).toMatch(/\/ws\/instances\/102\/metrics/);
  });

  it('soporta referencias de instancia ofuscadas (FIX-107, ej. inst_k8Pq2v) en la URL del WebSocket', async () => {
    const useInstanceMetrics = await loadHook();
    renderHook(() => useInstanceMetrics('inst_k8Pq2v'));
    await flush();

    expect(FakeWebSocket.instances).toHaveLength(1);
    expect(FakeWebSocket.instances[0].url).toMatch(/\/ws\/instances\/inst_k8Pq2v\/metrics\?.*ticket=ticket-1/);
  });
});

