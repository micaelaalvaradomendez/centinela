import { act, render, screen, waitFor, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstancesPage from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { withProtectedProviders } from './app-providers';
import { findOutsideTable } from './instances-helpers';

// FRN-17B (RF-11 / RF-04): al recibir un TASK_FINISHED cuyo detalles.tareaId coincide con una
// fila en transición, la fila se desbloquea, se actualiza su estado y se muestra un toast
// adaptativo (verde si COMPLETED; rojo con el motivo de D2 si FAILED).
const inventoryWithActiveTask = [
  {
    id: 101, name: 'servidor-web', type: 'vm', node: 'pve', status: 'running',
    ip: '192.168.1.50', cpuUsage: 25, ramUsage: 4294967296, maxRam: 17179869184, nivelAcceso: 'FULL_ACCESS',
    activeTask: { tareaId: 't-101', upid: 'UPID:pve:1:stop:101:', action: 'stop', status: 'RUNNING' },
  },
];

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

// Conexión SSE simulada: permite emitir TASK_FINISHED como lo haría el backend.
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

let latest: FakeSse | null = null;

beforeEach(() => {
  latest = null;
  class CapturingSse extends FakeSse {
    constructor(url: string) { super(url); latest = this; }
  }
  vi.stubGlobal('EventSource', CapturingSse);
  vi.stubGlobal('WebSocket', CapturingSse);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

function renderWithEvents() {
  const fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
    const method = (init.method ?? 'GET').toUpperCase();
    if (method === 'GET' && /\/instances(\?|$)/.test(String(url))) return Promise.resolve(json(inventoryWithActiveTask));
    if (String(url).includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-frn17b' }));
    return Promise.resolve(json({}));
  });
  vi.stubGlobal('fetch', fetchMock);
  window.sessionStorage.setItem('centinela_access', 'access-token-eventos');
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

async function rowOf(name: string) {
  const cell = await screen.findByText(name, {}, { timeout: 3000 });
  const row = cell.closest('tr, [role="row"]');
  expect(row, `"${name}" no está dentro de una fila`).not.toBeNull();
  return row as HTMLElement;
}

const actionButtons = (row: HTMLElement) => {
  const group = within(row).queryByRole('group', { name: /acciones/i });
  return within(group ?? row).queryAllByRole('button');
};

describe('FRN-17B - Desbloqueo reactivo y toast ante TASK_FINISHED', () => {
  it('libera la fila en transición y muestra un toast cuando la tarea termina COMPLETED', async () => {
    renderWithEvents();

    const row = await rowOf('servidor-web');
    await waitFor(() => expect(latest, 'no se abrió la conexión de eventos').not.toBeNull());

    // Con activeTask RUNNING la fila arranca bloqueada.
    await waitFor(() => {
      const actions = actionButtons(row);
      expect(actions.length).toBeGreaterThan(0);
      expect(actions.every((button) => button.hasAttribute('disabled'))).toBe(true);
    });

    act(() => {
      latest!.emitMessage({
        id: 'ev-101', tipo: 'TASK_FINISHED', severidad: 'INFO', recursoTipo: 'VM', recursoId: '101',
        mensaje: 'finalizó', fechaHora: '2026-10-13T03:00:00Z',
        detalles: { tareaId: 't-101', accion: 'stop', estado: 'COMPLETED', exitstatus: 'OK' },
      });
    });

    // Toast de resultado (verde: COMPLETED).
    await findOutsideTable(/orden completada|finaliz[oó] correctamente|completad/i);

    // La fila se libera sin F5.
    await waitFor(() => {
      const actions = actionButtons(row);
      expect(actions.some((button) => !button.hasAttribute('disabled')), 'la fila siguió bloqueada tras el TASK_FINISHED').toBe(true);
    });
  });

  it('muestra un toast de error con el motivo cuando la tarea termina FAILED', async () => {
    renderWithEvents();

    const row = await rowOf('servidor-web');
    await waitFor(() => expect(latest).not.toBeNull());

    act(() => {
      latest!.emitMessage({
        id: 'ev-101-fail', tipo: 'TASK_FINISHED', severidad: 'WARNING', recursoTipo: 'VM', recursoId: '101',
        mensaje: 'falló', fechaHora: '2026-10-13T03:00:00Z',
        detalles: { tareaId: 't-101', accion: 'stop', estado: 'FAILED', exitstatus: 'command failed', motivo: 'TIMEOUT' },
      });
    });

    // Toast de error: menciona el fallo (motivo TIMEOUT → tiempo de espera).
    await findOutsideTable(/fall[oó]|tiempo de espera|no se pudo|no pudo/i);

    await waitFor(() => {
      const actions = actionButtons(row);
      expect(actions.some((button) => !button.hasAttribute('disabled')), 'la fila siguió bloqueada tras el TASK_FINISHED fallido').toBe(true);
    });
  });
});
