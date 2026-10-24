import { act, render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import CreateNewInstances from '@/pages/createNewInstances';
import Instances from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { withProtectedProviders } from './app-providers';

// FRN-26: Envío de creación, manejo de colisión de cuotas (409) y transición optimista por SSE.
// Entregables (documentacion/actual.md):
//   1. Envío a POST /api/instances con estado de carga en el botón ("Creando instancia...", deshabilitado).
//   2. Manejo de error 409 Conflict (RESOURCE_QUOTA_EXCEEDED): si otro usuario consumió recursos
//      del host mientras se completaba el wizard, mostrar un banner/alerta con el mensaje exacto
//      del backend y permitir volver al Paso 2 sin perder el resto del formulario.
//   3. Ante respuesta 202 Accepted ({ vmid, tareaId }), cerrar el wizard / navegar a /instances
//      y mostrar la nueva instancia en la tabla con estado transitorio "Creando..." (spinner).
//   4. Suscripción al evento SSE TASK_FINISHED (accion: "create" / "CREATE"): al completarse,
//      actualizar la fila a "Stopped" y mostrar toast de éxito; si falla, mostrar toast de error
//      con el motivo reportado.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const GIB = 1024 ** 3;
const MIB = 1024 ** 2;

const storages = [
  { name: 'pool-rapido', type: 'lvmthin', totalBytes: 200 * GIB, freeBytes: 180 * GIB, contentTypes: ['images', 'rootdir'] },
];

const templates = [
  { id: 'local:iso/ubuntu-24.04.iso', name: 'ubuntu-24.04.iso', type: 'ISO', sizeBytes: 2000 * MIB },
  { id: 'local:vztmpl/debian-12-standard_12.2-1_amd64.tar.zst', name: 'debian-12-standard_12.2-1_amd64.tar.zst', type: 'LXC_TEMPLATE', sizeBytes: 120 * MIB },
];

const nodeStatus = {
  cpu: { usagePercent: 10, cores: 8 },
  ram: { usedGb: 2, totalGb: 16, usagePercent: 12 },
  storage: { usedGb: 30, totalGb: 300, usagePercent: 10 },
  uptimeSeconds: 3600,
  allocatable: { maxCpuCores: 8, freeRamMb: 12288, freeStorageGb: 180, headroomPercent: 15 },
  instancesSummary: { vms: { running: 0, stopped: 0, paused: 0, total: 0 }, lxc: { running: 0, stopped: 0, paused: 0, total: 0 } },
  stale: false,
  fetchedAt: '2026-10-24T12:00:00Z',
};

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
      const openEvt = new MessageEvent('open');
      this.onopen?.(openEvt);
      (this.listeners.open ?? []).forEach((fn) => fn(openEvt));
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
    (this.listeners.message ?? []).forEach((fn) => {
      if (eventType !== 'message') fn(msg);
    });
  }
}

let fetchMock: ReturnType<typeof vi.fn>;

function setupAdminSession() {
  window.sessionStorage.setItem('centinela_access', 'access-token-frn26');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'admin-1',
    organizacionId: 'org-1',
    nombreCompleto: 'Admin Test',
    email: 'admin@centinela.local',
    rol: 'ADMIN',
    instanciasPermitidas: [],
    tiene2FA: true,
  }));
}

async function advanceWizardToReviewStep(user: ReturnType<typeof userEvent.setup>, vmName = 'vm-frn26') {
  // Paso 1: Tipo y nombre
  await user.click(await screen.findByRole('button', { name: /m[aá]quina virtual/i }));
  const nameInput = screen.getByLabelText(/nombre de la instancia/i);
  await user.type(nameInput, vmName);
  await user.click(screen.getByRole('button', { name: /siguiente/i }));

  // Paso 2: Recursos y almacenamiento
  await waitFor(() => expect(screen.queryAllByRole('spinbutton').length).toBeGreaterThanOrEqual(3));
  await user.click(await screen.findByRole('combobox', { name: /almacenamiento de destino/i }));
  await user.click(await screen.findByRole('option', { name: /pool-rapido/i }));
  await waitFor(() => expect(screen.getByRole('button', { name: /siguiente/i })).toBeEnabled());
  await user.click(screen.getByRole('button', { name: /siguiente/i }));

  // Paso 3: Red (DHCP por defecto)
  await screen.findByRole('radiogroup', { name: /modo de asignaci[oó]n ipv4/i });
  await user.click(screen.getByRole('button', { name: /siguiente/i }));

  // Paso 4: Sistema Operativo
  await screen.findByRole('radiogroup', { name: /template del sistema operativo/i });
  await user.click(await screen.findByRole('radio', { name: /seleccionar ubuntu-24\.04\.iso/i }));
  await waitFor(() => expect(screen.getByRole('button', { name: /siguiente/i })).toBeEnabled());
  await user.click(screen.getByRole('button', { name: /siguiente/i }));

  // Paso 5: Resumen y confirmación
  await screen.findByRole('heading', { name: /resumen de la configuraci[oó]n/i });
  const confirmCheckbox = screen.getByRole('checkbox');
  await user.click(confirmCheckbox);
}

beforeEach(() => {
  FakeEventSource.instances = [];
  vi.stubGlobal('EventSource', FakeEventSource);
  setupAdminSession();
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-26 - Envío de creación, manejo de 409 RESOURCE_QUOTA_EXCEEDED y transición SSE', () => {
  it('envía POST /api/instances y deshabilita el botón mostrando "Creando instancia..." mientras está en vuelo', async () => {
    let resolveCreation!: (res: Response) => void;
    const pendingPost = new Promise<Response>((resolve) => {
      resolveCreation = resolve;
    });

    fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
      const target = String(url);
      const method = (init.method ?? 'GET').toUpperCase();
      if (method === 'POST' && /\/instances(\?|$)/.test(target)) return pendingPost;
      if (/\/node\/status(\?|$)/.test(target)) return Promise.resolve(json(nodeStatus));
      if (/\/node\/storages(\?|$)/.test(target)) return Promise.resolve(json(storages));
      if (/\/node\/templates(\?|$)/.test(target)) return Promise.resolve(json(templates));
      if (target.includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-frn26' }));
      if (/\/instances(\?|$)/.test(target)) return Promise.resolve(json([]));
      return Promise.resolve(json({}));
    });
    vi.stubGlobal('fetch', fetchMock);

    render(withProtectedProviders(
      <MemoryRouter initialEntries={['/instances/new']}>
        <Toaster />
        <Routes>
          <Route path="/instances/new" element={<CreateNewInstances />} />
          <Route path="/instances" element={<Instances />} />
        </Routes>
      </MemoryRouter>,
    ));

    const user = userEvent.setup();
    await advanceWizardToReviewStep(user, 'vm-en-vuelo');

    const submitButtons = screen.getAllByRole('button', { name: /^crear instancia$/i });
    await user.click(submitButtons[0]);

    // Mientras la petición está en vuelo, el botón debe indicar "Creando instancia..." y estar deshabilitado
    const loadingButtons = await screen.findAllByRole('button', { name: /creando instancia/i });
    expect(loadingButtons.length).toBeGreaterThan(0);
    for (const btn of loadingButtons) {
      expect(btn).toBeDisabled();
    }

    // Verificar payload despachado a POST /instances
    const postCall = (fetchMock.mock.calls as [string, RequestInit][]).find(
      ([url, init]) => (init?.method ?? 'GET').toUpperCase() === 'POST' && /\/instances(\?|$)/.test(String(url)),
    );
    expect(postCall).toBeDefined();
    const sentBody = JSON.parse(String(postCall![1].body)) as Record<string, unknown>;
    expect(sentBody.nombre ?? sentBody.name).toBe('vm-en-vuelo');
    expect(String(sentBody.tipo ?? sentBody.type).toLowerCase()).toBe('vm');

    await act(async () => {
      resolveCreation(json({ vmid: 'inst_120', tareaId: 'tarea-120' }, 202));
    });
  });

  it('ante 409 Conflict (RESOURCE_QUOTA_EXCEEDED) muestra el mensaje exacto del backend y permite regresar al Paso 2 sin perder el formulario', async () => {
    const backendQuotaMessage = 'La memoria solicitada (4096 MB) excede la disponible en el nodo (1024 MB libres tras el margen de seguridad del 15%).';

    fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
      const target = String(url);
      const method = (init.method ?? 'GET').toUpperCase();
      if (method === 'POST' && /\/instances(\?|$)/.test(target)) {
        return Promise.resolve(json({
          errorCode: 'RESOURCE_QUOTA_EXCEEDED',
          message: backendQuotaMessage,
        }, 409));
      }
      if (/\/node\/status(\?|$)/.test(target)) return Promise.resolve(json(nodeStatus));
      if (/\/node\/storages(\?|$)/.test(target)) return Promise.resolve(json(storages));
      if (/\/node\/templates(\?|$)/.test(target)) return Promise.resolve(json(templates));
      if (target.includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-frn26' }));
      return Promise.resolve(json({}));
    });
    vi.stubGlobal('fetch', fetchMock);

    render(withProtectedProviders(
      <MemoryRouter initialEntries={['/instances/new']}>
        <Routes>
          <Route path="/instances/new" element={<CreateNewInstances />} />
          <Route path="/instances" element={<Instances />} />
        </Routes>
      </MemoryRouter>,
    ));

    const user = userEvent.setup();
    await advanceWizardToReviewStep(user, 'vm-colision-cuota');

    const submitButtons = screen.getAllByRole('button', { name: /^crear instancia$/i });
    await user.click(submitButtons[0]);

    // Debe mostrarse un banner/alerta de error en línea dentro del propio wizard (no solo un toast global) y permanecer en el wizard
    const wizardSection = screen.getByRole('heading', { name: /crear nueva instancia/i }).closest('section')!;
    const inlineAlert = await within(wizardSection).findByRole('alert');
    expect(inlineAlert).toHaveTextContent(backendQuotaMessage);
    expect(screen.getByRole('heading', { name: /resumen de la configuraci[oó]n/i })).toBeVisible();

    // Regresar al Paso 2 ("Anterior" tres veces) sin perder los datos cargados
    const prevButton = screen.getByRole('button', { name: /anterior/i });
    await user.click(prevButton); // Paso 4
    await user.click(prevButton); // Paso 3
    await user.click(prevButton); // Paso 2

    await waitFor(() => expect(screen.queryAllByRole('spinbutton').length).toBeGreaterThanOrEqual(3));
    // El resumen lateral conserva el nombre 'vm-colision-cuota' y el almacenamiento 'pool-rapido'
    expect(screen.getAllByText(/vm-colision-cuota/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText(/pool-rapido/i).length).toBeGreaterThan(0);
  });

  it('ante 202 Accepted navega a /instances y al recibir TASK_FINISHED (CREATE COMPLETED) muestra la instancia en estado "Stopped" con toast de éxito', async () => {
    let taskCompleted = false;

    const creatingInstance = {
      id: 'inst_120',
      name: 'vm-optimista',
      type: 'vm',
      node: 'nodo-a1',
      status: 'stopped',
      ip: null,
      cpuUsage: null,
      ramUsage: null,
      maxRam: 2 * GIB,
      nivelAcceso: 'FULL_ACCESS',
      activeTask: {
        tareaId: 'tarea-create-120',
        action: 'CREATE',
        status: 'RUNNING',
      },
    };

    const completedInstance = {
      ...creatingInstance,
      activeTask: null,
    };

    fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
      const target = String(url);
      const method = (init.method ?? 'GET').toUpperCase();
      if (method === 'POST' && /\/instances(\?|$)/.test(target)) {
        return Promise.resolve(json({ vmid: 'inst_120', tareaId: 'tarea-create-120' }, 202));
      }
      if (/\/node\/status(\?|$)/.test(target)) return Promise.resolve(json(nodeStatus));
      if (/\/node\/storages(\?|$)/.test(target)) return Promise.resolve(json(storages));
      if (/\/node\/templates(\?|$)/.test(target)) return Promise.resolve(json(templates));
      if (target.includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-frn26' }));
      if (/\/instances(\?|$)/.test(target)) {
        return Promise.resolve(json([taskCompleted ? completedInstance : creatingInstance]));
      }
      return Promise.resolve(json({}));
    });
    vi.stubGlobal('fetch', fetchMock);

    render(withProtectedProviders(
      <MemoryRouter initialEntries={['/instances/new']}>
        <Toaster />
        <Routes>
          <Route path="/instances/new" element={<CreateNewInstances />} />
          <Route path="/instances" element={<Instances />} />
        </Routes>
      </MemoryRouter>,
    ));

    const user = userEvent.setup();
    await advanceWizardToReviewStep(user, 'vm-optimista');

    const submitButtons = screen.getAllByRole('button', { name: /^crear instancia$/i });
    await user.click(submitButtons[0]);

    // Debe cerrar el wizard y navegar a /instances mostrando la fila 'vm-optimista' en estado "Creando..."
    expect(await screen.findByRole('heading', { name: /inventario de instancias/i })).toBeVisible();
    expect(await screen.findByText('vm-optimista')).toBeVisible();
    expect(await screen.findByText(/creando/i)).toBeVisible();

    // Al recibir el evento SSE TASK_FINISHED (accion: 'CREATE', estado: 'COMPLETED'), actualiza a 'Stopped' y muestra toast de éxito
    taskCompleted = true;
    await waitFor(() => expect(FakeEventSource.instances.length).toBeGreaterThan(0));
    act(() => {
      FakeEventSource.instances[FakeEventSource.instances.length - 1].emitEvent({
        id: 'evt-create-120',
        tipo: 'TASK_FINISHED',
        severidad: 'INFO',
        recursoId: 'inst_120',
        detalles: {
          tareaId: 'tarea-create-120',
          accion: 'CREATE',
          estado: 'COMPLETED',
          exitstatus: 'OK',
          motivo: null,
          error: null,
        },
      });
    });

    expect(await screen.findByText(/^stopped$/i)).toBeVisible();
    expect((await screen.findAllByText(/orden completada en vm-optimista/i)).length).toBeGreaterThan(0);
  });

  it('muestra el estado transitorio "Creando..." con spinner mientras la tarea CREATE está en estado RUNNING', async () => {
    fetchMock = vi.fn().mockImplementation((url: string) => {
      const target = String(url);
      if (target.includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-frn26' }));
      if (/\/instances(\?|$)/.test(target)) {
        return Promise.resolve(json([{
          id: 'inst_120',
          name: 'vm-en-creacion',
          type: 'vm',
          node: 'nodo-a1',
          status: 'stopped',
          ip: null,
          cpuUsage: null,
          ramUsage: null,
          maxRam: 2 * GIB,
          nivelAcceso: 'FULL_ACCESS',
          activeTask: {
            tareaId: 'tarea-create-120',
            action: 'CREATE',
            status: 'RUNNING',
          },
        }]));
      }
      return Promise.resolve(json({}));
    });
    vi.stubGlobal('fetch', fetchMock);

    render(withProtectedProviders(
      <MemoryRouter initialEntries={['/instances']}>
        <Routes>
          <Route path="/instances" element={<Instances />} />
        </Routes>
      </MemoryRouter>,
    ));

    expect(await screen.findByText('vm-en-creacion')).toBeVisible();
    expect(await screen.findByText(/creando/i)).toBeVisible();
  });
});
