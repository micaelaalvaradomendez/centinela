import React from 'react';
import { act, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { Toaster } from '@/components/ui/toast';
import { withAppProviders } from './app-providers';

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
    // La tarea pide "disparar el cierre de sesión local y redirigir a /login". Se acepta que el cliente
    // lo haga directo o avisando a la app (API_UNAUTHORIZED_EVENT) para que lo haga el AuthProvider:
    // por eso el hook se monta con los providers de la app, sin EventsProvider (que abriría su propia
    // conexión). La redirección se ve en el historial del navegador (el router usa history.*State).
    const useEvents = await loadUseEvents();
    window.history.replaceState(null, '', '/dashboard');
    const historyCalls = [vi.spyOn(window.history, 'replaceState'), vi.spyOn(window.history, 'pushState')];
    // Si el cliente delega en el manejo global de 401 de la app (API_UNAUTHORIZED_EVENT → AuthProvider limpia
    // la sesión y navega a /login con el router), la navegación no se puede observar en jsdom: el router de la
    // app se crea al importar el módulo y en este entorno no tiene estado. En ese caso se acepta el aviso.
    const { API_UNAUTHORIZED_EVENT } = await import('@/services/apiClient');
    let unauthorizedNotified = false;
    const onUnauthorized = () => { unauthorizedNotified = true; };
    window.addEventListener(API_UNAUTHORIZED_EVENT, onUnauthorized);
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ errorCode: 'TOKEN_REVOKED', message: 'Sesión revocada' }, 401));
    vi.stubGlobal('fetch', fetchMock);
    try {
      // Como en App.tsx: Toaster por fuera del AuthProvider (el provider avisa con un toast antes de redirigir).
      renderHook(() => useEvents(), { wrapper: ({ children }) => <Toaster>{withAppProviders(children)}</Toaster> });

      await waitFor(() => expect(window.sessionStorage.getItem('centinela_access')).toBeNull(), { timeout: 5000 });
      await waitFor(() => expect(
        window.location.pathname === '/login'
          || historyCalls.some((spy) => spy.mock.calls.some((call) => String(call[2] ?? '').includes('/login')))
          || unauthorizedNotified,
        'no se redirigió a /login ni se avisó a la app de la sesión revocada (API_UNAUTHORIZED_EVENT)',
      ).toBe(true));
      expect(opened).toHaveLength(0);
    } finally {
      window.removeEventListener(API_UNAUTHORIZED_EVENT, onUnauthorized);
      window.history.replaceState(null, '', '/');
    }
  });
});

describe('FRN-17A - Consumo de eventos en tiempo real y distribución por instancia', () => {
  // Criterio de éxito: el Dashboard y la tabla reciben los TASK_FINISHED en tiempo real con UNA sola
  // conexión por pestaña; un evento repetido se procesa una sola vez. La tarea sugiere "un provider montado
  // una sola vez": si el frontend exporta un *Provider con un hook consumidor de eventos, los consumidores
  // usan ese hook dentro del provider (como en la app); si no, se usa useEvents directamente.
  // Cómo llegan los eventos al consumidor es libre: una función de suscripción (subscribe, onEvent,
  // addListener, listen), una lista (events) o el último evento (ultimoMensaje, ultimoEvento, lastEvent,
  // lastMessage). Cada evento se identifica por su detalles.tareaId, que el cliente conserva aunque
  // descarte otros campos.
  const SUBSCRIBE = ['subscribe', 'onEvent', 'addListener', 'listen'];
  const LAST = ['ultimoMensaje', 'ultimoEvento', 'lastEvent', 'lastMessage'];

  async function loadConsumer() {
    for (const load of Object.values(candidates)) {
      const module = await load().catch(() => ({}) as EventsModule);
      const providerName = Object.keys(module).find((name) => /Provider$/.test(name) && /event/i.test(name));
      const hookName = Object.keys(module).find((name) => /^use[A-Z]/.test(name) && /event/i.test(name) && typeof module[name] === 'function');
      if (providerName && hookName) {
        const Provider = module[providerName] as React.ComponentType<{ children?: React.ReactNode }>;
        return { name: `${hookName} con ${providerName}`, use: module[hookName] as () => unknown, wrapper: ({ children }: { children: React.ReactNode }) => <Provider>{children}</Provider> };
      }
    }
    return { name: 'useEvents', use: (await loadUseEvents()) as () => unknown, wrapper: undefined };
  }

  // Monta un consumidor y devuelve la lista de eventos que le llegaron (por identidad del objeto).
  async function mountConsumer() {
    const consumer = await loadConsumer();
    const snapshots: Record<string, unknown>[] = [];
    const received: Record<string, unknown>[] = [];
    const { result } = renderHook(() => {
      const value = consumer.use() as Record<string, unknown>;
      snapshots.push(value ?? {});
      return value;
    }, { wrapper: consumer.wrapper });
    await waitFor(() => expect(opened).toHaveLength(1));
    const value = result.current ?? {};
    const subscribe = SUBSCRIBE.map((key) => value[key]).find((fn) => typeof fn === 'function') as ((cb: (e: Record<string, unknown>) => void) => unknown) | undefined;
    const lastKey = LAST.find((key) => key in value);
    if (!subscribe && !Array.isArray(value.events) && !lastKey) {
      throw new Error(`FRN-17A: ${consumer.name} no expone cómo recibir los eventos (${[...SUBSCRIBE, 'events', ...LAST].join(', ')}); expone: ${Object.keys(value).join(', ') || 'nada'}`);
    }
    if (subscribe) act(() => { subscribe((event) => received.push(event)); });
    const read = (): Record<string, unknown>[] => {
      if (subscribe) return received;
      const last = snapshots[snapshots.length - 1];
      if (Array.isArray(last.events)) return last.events as Record<string, unknown>[];
      const seen = new Set<unknown>();
      return snapshots.map((snapshot) => snapshot[lastKey!]).filter((event) => event && !seen.has(event) && seen.add(event)) as Record<string, unknown>[];
    };
    return read;
  }

  const tareaOf = (event: Record<string, unknown>) => (event.detalles as Record<string, unknown> | undefined)?.tareaId;
  const taskFinished = (id: string, recursoId = '101') => ({
    id, tipo: 'TASK_FINISHED', severidad: 'INFO', recursoTipo: 'VM', recursoId, mensaje: 'La tarea de encendido finalizó correctamente',
    fechaHora: new Date().toISOString(), detalles: { tareaId: `tarea-${id}`, accion: 'start', estado: 'COMPLETED', exitstatus: 'OK' },
  });

  it('entrega los mensajes que cumplen RealtimeEvent y descarta los que no', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-parse' })));
    const read = await mountConsumer();

    act(() => {
      opened[0].instance.emitMessage({ invalid: 'payload', detalles: { tareaId: 'tarea-invalido' } });
      opened[0].instance.emitMessage('esto no es JSON');
      opened[0].instance.emitMessage({ ...taskFinished('sin-tipo'), tipo: undefined });
      opened[0].instance.emitMessage(taskFinished('valido'));
    });

    await waitFor(() => expect(read().map(tareaOf)).toContain('tarea-valido'));
    expect(read().map(tareaOf).filter((tarea) => tarea !== 'tarea-valido'), 'llegaron mensajes que no cumplen el contrato RealtimeEvent').toHaveLength(0);
  });

  it('un evento repetido (mismo id) se procesa una sola vez', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-dedup' })));
    const read = await mountConsumer();

    act(() => { opened[0].instance.emitMessage(taskFinished('dup')); });
    act(() => { opened[0].instance.emitMessage(taskFinished('dup')); });
    act(() => { opened[0].instance.emitMessage(taskFinished('otro', '102')); });

    await waitFor(() => expect(read().map(tareaOf)).toContain('tarea-otro'));
    expect(read().filter((event) => tareaOf(event) === 'tarea-dup'), 'el evento repetido se entregó dos veces').toHaveLength(1);
  });

  it('dos consumidores en la misma pestaña (Dashboard y tabla) comparten una sola conexión', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ticket: 'ticket-compartido' }));
    vi.stubGlobal('fetch', fetchMock);
    const consumer = await loadConsumer();

    renderHook(() => [consumer.use(), consumer.use()], { wrapper: consumer.wrapper });

    await waitFor(() => expect(opened.length).toBeGreaterThanOrEqual(1));
    await new Promise((resolve) => setTimeout(resolve, 300));
    expect(opened, `con dos consumidores (${consumer.name}) se abrió más de una conexión a /api/events`).toHaveLength(1);
    expect(ticketCalls(fetchMock)).toHaveLength(1);
  });

  it('se eliminó el hook vacío useWebSocket.js y notifications.ts ya no dice que el servidor no existe', () => {
    expect(Object.keys(import.meta.glob('@/hooks/useWebSocket.*')), 'sigue existiendo hooks/useWebSocket').toHaveLength(0);
    const notifications = Object.values(import.meta.glob<string>('@/types/notifications.ts', { query: '?raw', import: 'default', eager: true }))[0];
    expect(notifications).not.toMatch(/todav[ií]a no existe/i);
  });
});
