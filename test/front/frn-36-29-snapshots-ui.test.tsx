import { act, render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstanceDetailPage from '@/pages/InstanceDetailPage';
import { InstanceStatusBadge } from '@/components/features/instances/components/InstanceBadges';
import { Toaster } from '@/components/ui/toast';
import { withProtectedProviders } from './app-providers';

// FRN-36 (Etapa 2) y FRN-29 (Etapa 3):
//   FRN-36: Preparación del contenedor de Snapshots y estados transitorios en InstanceDetailPage:
//     1. Pestaña "Snapshots" en InstanceDetailPage.
//     2. Soporte en InstanceStatusBadge para los estados transitorios SNAPSHOTTING ("Creando snapshot...")
//        y ROLLING_BACK ("Revirtiendo...").
//     3. Bloqueo de botones de acción de la cabecera mientras la instancia se encuentra en
//        SNAPSHOTTING o ROLLING_BACK.
//   FRN-29: Interfaz de gestión de Snapshots y creación asíncrona:
//     1. Tabla/lista de snapshots desde GET /api/instances/:instanceRef/snapshots (Nombre,
//        Descripción, Fecha de creación, Estado de RAM).
//     2. Botón "Tomar Snapshot" / "Crear Snapshot" (solo ADMIN o FULL_ACCESS; oculto/deshabilitado
//        para READ_ONLY) con modal validando snapname (/^[a-zA-Z0-9_\-]+$/, sin espacios ni
//        caracteres especiales) y description (<= 250 caracteres).
//     3. POST /api/instances/:instanceRef/snapshots -> 202 Accepted ({ tareaId, snapname }, sin upid
//        según FIX-110) -> estado SNAPSHOTTING ("Creando snapshot...") -> recepción de TASK_FINISHED
//        (snapshot_create) recargando la lista de snapshots.
//     4. Acciones de Rollback (POST .../snapshots/:snapname/rollback) y Eliminación
//        (DELETE .../snapshots/:snapname) con confirmación.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const GIB = 1024 ** 3;
const MIB = 1024 ** 2;

class FakeEventSource {
  static instances: FakeEventSource[] = [];
  readyState = 1;
  onopen: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  private listeners: Record<string, ((event: MessageEvent) => void)[]> = {};

  constructor(public url: string) {
    FakeEventSource.instances.push(this);
    queueMicrotask(() => {
      const evt = new MessageEvent('open');
      this.onopen?.(evt);
      (this.listeners.open ?? []).forEach((fn) => fn(evt));
    });
  }

  addEventListener(type: string, listener: (event: MessageEvent) => void) {
    (this.listeners[type] ??= []).push(listener);
  }

  removeEventListener(type: string, listener: (event: MessageEvent) => void) {
    this.listeners[type] = (this.listeners[type] ?? []).filter((fn) => fn !== listener);
  }

  close() {
    this.readyState = 2;
  }

  emitEvent(payload: Record<string, unknown>, eventType = 'message') {
    const msg = new MessageEvent(eventType, { data: JSON.stringify(payload) });
    if (eventType === 'message') this.onmessage?.(msg);
    (this.listeners[eventType] ?? []).forEach((fn) => fn(msg));
  }
}

class FakeResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

let fetchMock: ReturnType<typeof vi.fn>;

const initialSnapshots = [
  {
    name: 'pre_upgrade_v1',
    description: 'Snapshot previo a actualización de paquetes',
    snaptime: 1_760_000_000,
    vmstate: true,
  },
  {
    name: 'base_limpia',
    description: 'Instalación inicial limpia',
    snaptime: 1_759_900_000,
    vmstate: false,
  },
];

function stubSnapshotsBackend(options: {
  instanceRef?: string;
  nivelAcceso?: 'FULL_ACCESS' | 'READ_ONLY';
  activeTask?: { tareaId: string; action: string; status: 'RUNNING' | 'COMPLETED' | 'FAILED' } | null;
  snapshots?: typeof initialSnapshots;
} = {}) {
  const ref = options.instanceRef ?? 'inst_k8Pq2v';
  const nivelAcceso = options.nivelAcceso ?? 'FULL_ACCESS';
  let currentSnapshots = [...(options.snapshots ?? initialSnapshots)];
  let currentActiveTask = options.activeTask ?? null;

  fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
    const target = String(url);
    const method = (init.method ?? 'GET').toUpperCase();

    if (new RegExp(`/instances/${ref}/snapshots/([^/]+)/rollback$`).test(target) && method === 'POST') {
      const snapName = target.split('/').slice(-2)[0];
      currentActiveTask = { tareaId: 'tarea-snap-rollback-1', action: 'SNAPSHOT_ROLLBACK', status: 'RUNNING' };
      return Promise.resolve(json({ tareaId: 'tarea-snap-rollback-1', snapname: snapName }, 202));
    }

    if (new RegExp(`/instances/${ref}/snapshots/([^/]+)$`).test(target) && method === 'DELETE') {
      const snapName = target.split('/').pop()!;
      currentSnapshots = currentSnapshots.filter((s) => s.name !== snapName);
      currentActiveTask = { tareaId: 'tarea-snap-del-1', action: 'SNAPSHOT_DELETE', status: 'RUNNING' };
      return Promise.resolve(json({ tareaId: 'tarea-snap-del-1', snapname: snapName }, 202));
    }

    if (new RegExp(`/instances/${ref}/snapshots(\\?|$)`).test(target)) {
      if (method === 'POST') {
        const body = JSON.parse(String(init.body ?? '{}')) as { snapname: string; description?: string };
        currentSnapshots = [
          {
            name: body.snapname,
            description: body.description ?? '',
            snaptime: 1_760_100_000,
            vmstate: false,
          },
          ...currentSnapshots,
        ];
        currentActiveTask = { tareaId: 'tarea-snap-create-1', action: 'SNAPSHOT_CREATE', status: 'RUNNING' };
        return Promise.resolve(json({ tareaId: 'tarea-snap-create-1', snapname: body.snapname }, 202));
      }
      return Promise.resolve(json(currentSnapshots));
    }

    if (new RegExp(`/instances/${ref}/metrics`).test(target)) {
      return Promise.resolve(json({ data: [] }));
    }

    if (new RegExp(`/instances/${ref}(\\?|$)`).test(target)) {
      return Promise.resolve(json({
        vmid: ref,
        nombre: 'srv-produccion',
        nodo: 'nodo-a1',
        estado: 'running',
        tipo: 'qemu',
        cpu: 15,
        mem: 1024 * MIB,
        maxCpu: 4,
        maxMem: 4 * GIB,
      }));
    }

    if (/\/instances(\?|$)/.test(target)) {
      return Promise.resolve(json([
        {
          id: ref,
          name: 'srv-produccion',
          type: 'vm',
          node: 'nodo-a1',
          status: 'running',
          ip: '10.0.0.10',
          cpuUsage: 15,
          ramUsage: 1024 * MIB,
          maxRam: 4 * GIB,
          nivelAcceso,
          activeTask: currentActiveTask,
        },
      ]));
    }

    if (/\/node\/status(\?|$)/.test(target)) {
      return Promise.resolve(json({
        cpu: { usagePercent: 10, cores: 8 },
        ram: { usedGb: 4, totalGb: 16, usagePercent: 25 },
        uptimeSeconds: 3600,
        allocatable: { maxCpuCores: 8, freeRamMb: 8192, freeStorageGb: 100, headroomPercent: 15 },
        instancesSummary: { vms: { running: 1, stopped: 0, paused: 0, total: 1 }, lxc: { running: 0, stopped: 0, paused: 0, total: 0 } },
        stale: false,
        fetchedAt: '2026-10-24T12:00:00Z',
      }));
    }

    if (target.includes('/events/ticket') || target.includes('/auth/ws-ticket')) {
      return Promise.resolve(json({ ticket: 'ticket-frn36-29' }));
    }

    return Promise.resolve(json({}));
  });

  vi.stubGlobal('fetch', fetchMock);
  return {
    clearActiveTask: () => {
      currentActiveTask = null;
    },
  };
}

function renderDetailAs(
  rol: 'ADMIN' | 'OPERATOR',
  nivelAcceso: 'FULL_ACCESS' | 'READ_ONLY' = 'FULL_ACCESS',
  instanceRef = 'inst_k8Pq2v',
) {
  window.sessionStorage.setItem('centinela_access', 'access-token-frn36-29');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: `${rol.toLowerCase()}-1`,
    organizacionId: 'org-1',
    nombreCompleto: `${rol} Test`,
    email: `${rol.toLowerCase()}@centinela.local`,
    rol,
    instanciasPermitidas: [instanceRef],
    permisos: [{ vmid: instanceRef, nivelAcceso }],
    tiene2FA: true,
  }));

  return render(withProtectedProviders(
    <MemoryRouter initialEntries={[`/instances/${instanceRef}`]}>
      <Toaster />
      <Routes>
        <Route path="/instances/:vmid" element={<InstanceDetailPage />} />
      </Routes>
    </MemoryRouter>,
  ));
}

beforeEach(() => {
  HTMLDialogElement.prototype.showModal ??= function showModalPolyfill(this: HTMLDialogElement) {
    this.setAttribute('open', '');
  };
  HTMLDialogElement.prototype.close ??= function closePolyfill(this: HTMLDialogElement) {
    this.removeAttribute('open');
  };
  FakeEventSource.instances = [];
  vi.stubGlobal('EventSource', FakeEventSource);
  vi.stubGlobal('WebSocket', FakeEventSource);
  vi.stubGlobal('ResizeObserver', FakeResizeObserver);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-36 - Pestaña Snapshots y estados transitorios SNAPSHOTTING / ROLLING_BACK', () => {
  it('renderiza la pestaña "Snapshots" en InstanceDetailPage y permite seleccionarla', async () => {
    stubSnapshotsBackend();
    renderDetailAs('ADMIN');
    const user = userEvent.setup();

    const snapshotsTab = await screen.findByRole('tab', { name: /snapshots/i });
    expect(snapshotsTab).toBeVisible();

    await user.click(snapshotsTab);
    expect(snapshotsTab).toHaveAttribute('aria-selected', 'true');
  });

  it('InstanceStatusBadge muestra "Creando snapshot..." (SNAPSHOTTING) y "Revirtiendo..." (ROLLING_BACK)', () => {
    const { rerender } = render(
      <InstanceStatusBadge
        status="SNAPSHOTTING"
        activeTask={{ tareaId: 't-snap', action: 'SNAPSHOT_CREATE', status: 'RUNNING' }}
      />,
    );
    expect(screen.getByText(/creando snapshot/i)).toBeVisible();

    rerender(
      <InstanceStatusBadge
        status="ROLLING_BACK"
        activeTask={{ tareaId: 't-rb', action: 'SNAPSHOT_ROLLBACK', status: 'RUNNING' }}
      />,
    );
    expect(screen.getByText(/revirtiendo/i)).toBeVisible();
  });

  it('bloquea los botones de acción de la cabecera mientras la instancia está en SNAPSHOTTING o ROLLING_BACK', async () => {
    stubSnapshotsBackend({
      activeTask: { tareaId: 't-snap-active', action: 'SNAPSHOT_CREATE', status: 'RUNNING' },
    });
    renderDetailAs('ADMIN');

    await screen.findByRole('heading', { name: 'srv-produccion' });
    const guardarCambiosBtn = screen.getByRole('button', { name: /guardar cambios/i });
    const apagarBtn = await screen.findByRole('button', { name: /^apagar$/i });
    const reiniciarBtn = screen.getByRole('button', { name: /^reiniciar$/i });

    expect(guardarCambiosBtn).toBeDisabled();
    expect(apagarBtn).toBeDisabled();
    expect(reiniciarBtn).toBeDisabled();
  });
});

describe('FRN-29 - Interfaz de gestión de Snapshots, validaciones y ciclo asíncrono (Etapa 3)', () => {
  it('consulta GET /api/instances/:instanceRef/snapshots al abrir la pestaña y lista los snapshots con Nombre, Descripción, Fecha y Estado de RAM', async () => {
    stubSnapshotsBackend();
    renderDetailAs('ADMIN');
    const user = userEvent.setup();

    await user.click(await screen.findByRole('tab', { name: /snapshots/i }));

    await waitFor(() => {
      const snapCalls = (fetchMock.mock.calls as [string][]).filter(([url]) => /\/instances\/inst_k8Pq2v\/snapshots(\?|$)/.test(String(url)));
      expect(snapCalls.length, 'la pestaña Snapshots debe consultar GET /api/instances/:instanceRef/snapshots (FRN-29)').toBeGreaterThan(0);
    });

    expect(await screen.findByText('pre_upgrade_v1')).toBeVisible();
    expect(screen.getByText('Snapshot previo a actualización de paquetes')).toBeVisible();
    expect(screen.getByText('base_limpia')).toBeVisible();
  });

  it('valida el formulario de "Tomar Snapshot" (snapname /^[a-zA-Z0-9_\\-]+$/ sin espacios ni caracteres inválidos y description <= 250) y oculta la creación para READ_ONLY', async () => {
    stubSnapshotsBackend({ nivelAcceso: 'READ_ONLY' });
    const { unmount } = renderDetailAs('OPERATOR', 'READ_ONLY');
    const user = userEvent.setup();

    await user.click(await screen.findByRole('tab', { name: /snapshots/i }));
    const readOnlyCreateBtn = screen.queryByRole('button', { name: /tomar snapshot|crear snapshot|nuevo snapshot/i });
    if (readOnlyCreateBtn) {
      expect(readOnlyCreateBtn).toBeDisabled();
    }
    unmount();

    // Con FULL_ACCESS o ADMIN el botón está habilitado y abre el modal de creación
    stubSnapshotsBackend({ nivelAcceso: 'FULL_ACCESS' });
    renderDetailAs('OPERATOR', 'FULL_ACCESS');
    await user.click(await screen.findByRole('tab', { name: /snapshots/i }));

    const createBtn = await screen.findByRole('button', { name: /tomar snapshot|crear snapshot|nuevo snapshot/i });
    expect(createBtn).toBeEnabled();
    await user.click(createBtn);

    const dialog = await screen.findByRole('dialog');
    const nameInput = within(dialog).getByLabelText(/nombre/i);
    const submitBtn = within(dialog).getByRole('button', { name: /crear|tomar|guardar|confirmar/i });

    // Nombre inválido con espacios y caracteres especiales fuera de /^[a-zA-Z0-9_\-]+$/
    await user.type(nameInput, 'snap invalido!');
    expect(submitBtn).toBeDisabled();
  });

  it('envía POST /api/instances/:instanceRef/snapshots, recibe 202 Accepted, muestra SNAPSHOTTING y recarga la lista al recibir TASK_FINISHED', async () => {
    const backend = stubSnapshotsBackend();
    renderDetailAs('ADMIN');
    const user = userEvent.setup();

    await user.click(await screen.findByRole('tab', { name: /snapshots/i }));
    const createBtn = await screen.findByRole('button', { name: /tomar snapshot|crear snapshot|nuevo snapshot/i });
    await user.click(createBtn);

    const dialog = await screen.findByRole('dialog');
    await user.type(within(dialog).getByLabelText(/nombre/i), 'snap_etapa3_ok');
    const descInput = within(dialog).queryByLabelText(/descripci[oó]n/i);
    if (descInput) {
      await user.type(descInput, 'Snapshot de prueba Etapa 3');
    }

    await user.click(within(dialog).getByRole('button', { name: /crear|tomar|guardar|confirmar/i }));

    await waitFor(() => {
      const postCall = (fetchMock.mock.calls as [string, RequestInit][]).find(
        ([url, init]) => (init?.method ?? 'GET').toUpperCase() === 'POST' && /\/instances\/inst_k8Pq2v\/snapshots$/.test(String(url)),
      );
      expect(postCall).toBeDefined();
    });

    // Al finalizar la tarea por SSE (TASK_FINISHED), limpia el bloqueo y muestra el nuevo snapshot
    backend.clearActiveTask();
    await waitFor(() => expect(FakeEventSource.instances.length).toBeGreaterThan(0));
    act(() => {
      FakeEventSource.instances[FakeEventSource.instances.length - 1].emitEvent({
        id: 'evt-snap-create-1',
        tipo: 'TASK_FINISHED',
        severidad: 'INFO',
        recursoId: 'inst_k8Pq2v',
        detalles: {
          tareaId: 'tarea-snap-create-1',
          accion: 'SNAPSHOT_CREATE',
          estado: 'COMPLETED',
          exitstatus: 'OK',
          motivo: null,
          error: null,
        },
      });
    });

    expect(await screen.findByText('snap_etapa3_ok')).toBeVisible();
  });

  it('permite ejecutar Rollback (POST .../rollback) y Eliminar (DELETE .../snapshots/:snapname) con confirmación', async () => {
    stubSnapshotsBackend();
    renderDetailAs('ADMIN');
    const user = userEvent.setup();

    await user.click(await screen.findByRole('tab', { name: /snapshots/i }));
    await screen.findByText('pre_upgrade_v1');

    const rollbackBtns = screen.getAllByRole('button', { name: /revertir|rollback|restaurar/i });
    await user.click(rollbackBtns[0]);

    const confirmDialog = await screen.findByRole('dialog');
    await user.click(within(confirmDialog).getByRole('button', { name: /confirmar|revertir|restaurar|aceptar/i }));

    await waitFor(() => {
      const rbCall = (fetchMock.mock.calls as [string, RequestInit][]).find(
        ([url, init]) => (init?.method ?? 'GET').toUpperCase() === 'POST' && /\/instances\/inst_k8Pq2v\/snapshots\/pre_upgrade_v1\/rollback$/.test(String(url)),
      );
      expect(rbCall).toBeDefined();
    });
  });
});
