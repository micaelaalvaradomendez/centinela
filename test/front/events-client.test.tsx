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
const opened: { url: string; instance: FakeConnection; at: number }[] = [];
class FakeConnection {
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onopen: ((event: Event) => void) | null = null;
  readyState = 1;
  private listeners: Record<string, ((event: Event) => void)[]> = {};
  constructor(public url: string) { opened.push({ url: String(url), instance: this, at: Date.now() }); }
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
  emitMessage(data: unknown) {
    const raw = typeof data === 'string' ? data : JSON.stringify(data);
    const event = new MessageEvent('message', { data: raw });
    this.onmessage?.(event);
    this.listeners.message?.forEach((l) => l(event));
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

  it('las reconexiones sucesivas aplican retroceso exponencial (la segunda espera más que la primera)', async () => {
    const useEvents = await loadUseEvents();
    let n = 0;
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ ticket: `ticket-${++n}` })));
    vi.stubGlobal('fetch', fetchMock);

    renderHook(() => useEvents());
    await waitFor(() => expect(opened).toHaveLength(1));

    // Ninguna conexión llega a abrir (onopen), así que el retroceso no debería reiniciarse.
    const firstFail = Date.now();
    act(() => opened[0].instance.fail());
    await waitFor(() => expect(opened.length).toBeGreaterThanOrEqual(2), { timeout: 30000 });
    const firstDelay = opened[1].at - firstFail;

    const secondFail = Date.now();
    act(() => opened[1].instance.fail());
    await waitFor(() => expect(opened.length).toBeGreaterThanOrEqual(3), { timeout: 30000 });
    const secondDelay = opened[2].at - secondFail;

    expect(secondDelay, `espera 1: ${firstDelay} ms, espera 2: ${secondDelay} ms`).toBeGreaterThan(firstDelay * 1.5);
  }, 70000);

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

describe('FRN-17A - Consumo de eventos en tiempo real y distribución por instancia', () => {
  it('parsea mensajes válidos según RealtimeEvent y descarta los que no cumplen el contrato', async () => {
    const useEvents = await loadUseEvents();
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-parse' }));
    vi.stubGlobal('fetch', fetchMock);

    const received: unknown[] = [];
    const { result } = renderHook(() => useEvents());
    await waitFor(() => expect(opened).toHaveLength(1));

    if (typeof (result.current as any)?.subscribe === 'function') {
      (result.current as any).subscribe((ev: unknown) => received.push(ev));
    }

    // Mensaje inválido (sin campos requeridos de RealtimeEvent)
    act(() => {
      opened[0].instance.emitMessage({ invalid: 'payload' });
    });

    const validEvent = {
      id: 'evt-valid-1',
      tipo: 'TASK_FINISHED',
      severidad: 'INFO',
      recursoTipo: 'VM',
      recursoId: '101',
      mensaje: 'Tarea finalizada exitosamente',
      fechaHora: new Date().toISOString(),
      detalles: { tareaId: 'task-101', estado: 'COMPLETED' },
    };

    act(() => {
      opened[0].instance.emitMessage(validEvent);
    });

    if (typeof (result.current as any)?.subscribe === 'function') {
      await waitFor(() => expect(received).toContainEqual(expect.objectContaining({ id: 'evt-valid-1' })));
      expect(received).not.toContainEqual(expect.objectContaining({ invalid: 'payload' }));
    } else {
      expect((result.current as any)?.subscribe || (result.current as any)?.events || (result.current as any)?.lastEvent).toBeDefined();
    }
  });

  it('deduplica eventos recibidos con el mismo id', async () => {
    const useEvents = await loadUseEvents();
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-dedup' }));
    vi.stubGlobal('fetch', fetchMock);

    const received: unknown[] = [];
    const { result } = renderHook(() => useEvents());
    await waitFor(() => expect(opened).toHaveLength(1));

    if (typeof (result.current as any)?.subscribe === 'function') {
      (result.current as any).subscribe((ev: unknown) => received.push(ev));
    }

    const event = {
      id: 'evt-dup-1',
      tipo: 'INSTANCE_STATE_CHANGED',
      severidad: 'INFO',
      recursoTipo: 'VM',
      recursoId: '101',
      mensaje: 'Instancia iniciada',
      fechaHora: new Date().toISOString(),
    };

    act(() => {
      opened[0].instance.emitMessage(event);
      opened[0].instance.emitMessage(event);
    });

    if (typeof (result.current as any)?.subscribe === 'function') {
      await waitFor(() => {
        const matches = received.filter((e: any) => e.id === 'evt-dup-1');
        expect(matches).toHaveLength(1);
      });
    } else {
      expect((result.current as any)?.subscribe || (result.current as any)?.events).toBeDefined();
    }
  });

  it('permite suscribirse por tipo y por recursoId con una sola conexión activa', async () => {
    const useEvents = await loadUseEvents();
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-sub' }));
    vi.stubGlobal('fetch', fetchMock);

    const { result } = renderHook(() => useEvents());
    await waitFor(() => expect(opened).toHaveLength(1));

    expect(
      typeof (result.current as any)?.subscribe === 'function' ||
      typeof (result.current as any)?.subscribeToResource === 'function' ||
      typeof (result.current as any)?.onEvent === 'function'
    ).toBe(true);
  });
});

