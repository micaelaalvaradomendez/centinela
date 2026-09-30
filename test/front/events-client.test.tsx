import { act, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

// FRN-17C (BRG-02-FRN): el cliente de eventos pide un ticket efímero antes de conectarse.
// La tarea no fija el archivo: se busca en src/ cualquier módulo que exporte `useEvents` (FRN-17A).
type EventsModule = { useEvents?: (...args: unknown[]) => unknown; [key: string]: unknown };
const candidates = import.meta.glob<EventsModule>([
  '@/hooks/**/*.{ts,tsx,js,jsx}',
  '@/components/**/hooks/**/*.{ts,tsx,js,jsx}',
  '@/services/**/*.{ts,tsx,js,jsx}',
  '@/context/**/*.{ts,tsx,js,jsx}',
]);

async function loadUseEvents() {
  for (const load of Object.values(candidates)) {
    const module = await load().catch(() => ({}) as EventsModule);
    if (typeof module.useEvents === 'function') return module.useEvents;
  }
  throw new Error('FRN-17C no implementada: ningún módulo de hooks/, services/ o context/ exporta useEvents');
}

// Sustitutos de EventSource y WebSocket: registran la URL con la que se abre cada conexión.
const opened: { url: string; instance: FakeConnection }[] = [];
class FakeConnection {
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onopen: ((event: Event) => void) | null = null;
  readyState = 1;
  private listeners: Record<string, ((event: Event) => void)[]> = {};
  constructor(public url: string) { opened.push({ url: String(url), instance: this }); }
  addEventListener(type: string, listener: (event: Event) => void) { (this.listeners[type] ??= []).push(listener); }
  removeEventListener() {}
  close() { this.readyState = 2; }
  send() {}
  fail() {
    this.readyState = 2;
    const event = new Event('error');
    this.onerror?.(event); this.listeners.error?.forEach((l) => l(event));
    const close = new Event('close');
    this.onclose?.(close); this.listeners.close?.forEach((l) => l(close));
  }
}

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });
}

const ticketCalls = (fetchMock: ReturnType<typeof vi.fn>) =>
  (fetchMock.mock.calls as [string, RequestInit][]).filter(([url]) => String(url).includes('/events/ticket'));

beforeEach(() => {
  opened.length = 0;
  vi.stubGlobal('EventSource', FakeConnection);
  vi.stubGlobal('WebSocket', FakeConnection);
  window.sessionStorage.setItem('centinela_access', 'access-token-largo');
});

afterEach(() => {
  vi.useRealTimers();
  vi.unstubAllGlobals();
});

describe('FRN-17C - Cliente de eventos con ticket efímero y reconexión segura', () => {
  it('pide POST /api/events/ticket con Bearer y se conecta a /api/events?ticket=<uuid> sin exponer el JWT', async () => {
    const useEvents = await loadUseEvents();
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-uno' }));
    vi.stubGlobal('fetch', fetchMock);

    renderHook(() => useEvents());

    await waitFor(() => expect(ticketCalls(fetchMock)).toHaveLength(1));
    const [, init] = ticketCalls(fetchMock)[0];
    expect(init.method).toBe('POST');
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer access-token-largo');

    await waitFor(() => expect(opened).toHaveLength(1));
    expect(opened[0].url).toMatch(/\/api\/events\?ticket=ticket-uno$/);
    expect(opened[0].url).not.toContain('access-token-largo');
  });

  it('ante una desconexión pide un ticket nuevo (no reutiliza el anterior) y se reconecta', async () => {
    const useEvents = await loadUseEvents();
    let n = 0;
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ ticket: `ticket-${++n}` })));
    vi.stubGlobal('fetch', fetchMock);

    renderHook(() => useEvents());
    await waitFor(() => expect(opened).toHaveLength(1));

    act(() => opened[0].instance.fail());

    await waitFor(() => expect(opened.length).toBeGreaterThanOrEqual(2), { timeout: 20000 });
    expect(ticketCalls(fetchMock).length).toBeGreaterThanOrEqual(2);
    expect(opened[1].url).toMatch(/ticket=ticket-2$/);
  }, 30000);

  it('si /api/events/ticket responde 401 cierra la sesión local y lleva a /login', async () => {
    const useEvents = await loadUseEvents();
    const replace = vi.fn();
    const originalLocation = window.location;
    Object.defineProperty(window, 'location', { configurable: true, value: { ...originalLocation, pathname: '/dashboard', replace, assign: replace } });
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ errorCode: 'TOKEN_REVOKED', message: 'Sesión revocada' }, 401));
    vi.stubGlobal('fetch', fetchMock);
    try {
      renderHook(() => useEvents());

      await waitFor(() => expect(window.sessionStorage.getItem('centinela_access')).toBeNull(), { timeout: 5000 });
      expect(opened).toHaveLength(0);
    } finally {
      Object.defineProperty(window, 'location', { configurable: true, value: originalLocation });
    }
  });
});
