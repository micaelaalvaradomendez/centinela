import React from 'react';
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
  // La tarea no fija cómo se consumen los eventos ("por ejemplo con un provider"). Se acepta que el
  // valor de useEvents exponga una función de suscripción (subscribe, onEvent, addListener, listen o
  // subscribeToResource) que recibe un callback, o el estado `events` (lista) / `lastEvent`. Si no
  // expone ninguna de esas formas, la prueba falla: no hay manera de que el Dashboard o la tabla reciban.
  const SUBSCRIBE = ['subscribe', 'onEvent', 'addListener', 'listen', 'subscribeToResource'];

  function collect(result: { current: unknown }) {
    const received: Record<string, unknown>[] = [];
    const value = result.current as Record<string, unknown> | null;
    const subscribe = SUBSCRIBE.map((name) => value?.[name]).find((fn) => typeof fn === 'function') as
      | ((callback: (event: Record<string, unknown>) => void) => unknown)
      | undefined;
    if (subscribe) act(() => { subscribe((event) => received.push(event)); });
    const read = () => {
      if (subscribe) return received;
      const current = result.current as { events?: unknown; lastEvent?: unknown } | null;
      if (Array.isArray(current?.events)) return current.events as Record<string, unknown>[];
      if (current && 'lastEvent' in current) return current.lastEvent ? [current.lastEvent as Record<string, unknown>] : [];
      throw new Error(`FRN-17A: useEvents no expone cómo recibir los eventos (${SUBSCRIBE.join(', ')}, events o lastEvent); expone: ${Object.keys(current ?? {}).join(', ') || 'nada'}`);
    };
    return read;
  }

  const taskFinished = (id: string, recursoId = '101') => ({
    id, tipo: 'TASK_FINISHED', severidad: 'INFO', recursoTipo: 'VM', recursoId, mensaje: 'La tarea de encendido finalizó correctamente',
    fechaHora: new Date().toISOString(), detalles: { tareaId: `tarea-${id}`, accion: 'start', estado: 'COMPLETED', exitstatus: 'OK' },
  });

  it('entrega los mensajes que cumplen RealtimeEvent y descarta los que no', async () => {
    const useEvents = await loadUseEvents();
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-parse' })));

    const { result } = renderHook(() => useEvents());
    await waitFor(() => expect(opened).toHaveLength(1));
    const read = collect(result);

    act(() => {
      opened[0].instance.emitMessage({ invalid: 'payload' });
      opened[0].instance.emitMessage('esto no es JSON');
      opened[0].instance.emitMessage({ ...taskFinished('evt-sin-tipo'), tipo: undefined });
      opened[0].instance.emitMessage(taskFinished('evt-valid-1'));
    });

    await waitFor(() => expect(read()).toContainEqual(expect.objectContaining({ id: 'evt-valid-1', tipo: 'TASK_FINISHED' })));
    expect(read().filter((event) => event.id !== 'evt-valid-1'), 'llegaron mensajes que no cumplen el contrato RealtimeEvent').toHaveLength(0);
  });

  it('un evento repetido (mismo id) se procesa una sola vez', async () => {
    const useEvents = await loadUseEvents();
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-dedup' })));

    const { result } = renderHook(() => useEvents());
    await waitFor(() => expect(opened).toHaveLength(1));
    const read = collect(result);

    act(() => {
      opened[0].instance.emitMessage(taskFinished('evt-dup-1'));
      opened[0].instance.emitMessage(taskFinished('evt-dup-1'));
      opened[0].instance.emitMessage(taskFinished('evt-otro', '102'));
    });

    await waitFor(() => expect(read()).toContainEqual(expect.objectContaining({ id: 'evt-otro' })));
    expect(read().filter((event) => event.id === 'evt-dup-1')).toHaveLength(1);
  });

  it('dos consumidores en la misma pestaña (Dashboard y tabla) comparten una sola conexión', async () => {
    const useEvents = await loadUseEvents();
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-compartido' }));
    vi.stubGlobal('fetch', fetchMock);

    // Si la conexión vive en un provider (como sugiere la tarea), se monta igual que en la app.
    const providers = Object.values(await Promise.all(Object.values(candidates).map((load) => load().catch(() => ({}) as EventsModule))))
      .flatMap((module) => Object.entries(module).filter(([name, value]) => /Provider$/.test(name) && typeof value === 'function'))
      .map(([, value]) => value as React.ComponentType<{ children?: React.ReactNode }>);
    const wrapper = ({ children }: { children: React.ReactNode }) =>
      providers.reduceRight<React.ReactElement>((inner, Provider) => React.createElement(Provider, null, inner), React.createElement(React.Fragment, null, children));

    renderHook(() => [useEvents(), useEvents()], { wrapper });

    await waitFor(() => expect(opened.length).toBeGreaterThanOrEqual(1));
    await new Promise((resolve) => setTimeout(resolve, 300));
    expect(opened, 'cada consumidor abrió su propia conexión a /api/events').toHaveLength(1);
    expect(ticketCalls(fetchMock)).toHaveLength(1);
  });

  it('se eliminó el hook vacío useWebSocket.js y notifications.ts ya no dice que el servidor no existe', () => {
    expect(Object.keys(import.meta.glob('@/hooks/useWebSocket.*')), 'sigue existiendo hooks/useWebSocket').toHaveLength(0);
    const notifications = Object.values(import.meta.glob<string>('@/types/notifications.ts', { query: '?raw', import: 'default', eager: true }))[0];
    expect(notifications).not.toMatch(/todav[ií]a no existe/i);
  });
});
