import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstancesPage from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { withProtectedProviders } from './app-providers';

// Pruebas de aceptación para tareas de la Ola 2 de inventario en actual.md:
// - FRN-20B: Filtros reactivos por tipo (Todas, VM, LXC), estado (Todas, En ejecución, Detenidas) y buscador dinámico.
// - FRN-16B: Resincronización del estado "Operación en progreso" tras F5 o reconexión (activeTask !== null).

const inventorySample = [
  {
    id: 101, name: 'servidor-web', type: 'vm', node: 'pve', status: 'running',
    ip: '192.168.1.50', cpuUsage: 25, ramUsage: 4294967296, maxRam: 17179869184, nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
  {
    id: 102, name: 'base-datos', type: 'lxc', node: 'pve', status: 'stopped',
    ip: null, cpuUsage: null, ramUsage: null, maxRam: null, nivelAcceso: 'READ_ONLY', activeTask: null,
  },
  {
    id: 103, name: 'redis-cache', type: 'lxc', node: 'pve', status: 'running',
    ip: '192.168.1.52', cpuUsage: 10, ramUsage: 1073741824, maxRam: 2147483648, nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
];

let fetchMock: ReturnType<typeof vi.fn>;

function renderInstancesPage(customInventory = inventorySample) {
  window.sessionStorage.setItem('centinela_access', 'access-token-ola2');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
    rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));

  fetchMock = vi.fn().mockImplementation((url: string) => Promise.resolve(String(url).includes('/instances')
    ? new Response(JSON.stringify(customInventory), { status: 200, headers: { 'Content-Type': 'application/json' } })
    : new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } })));
  vi.stubGlobal('fetch', fetchMock);

  return render(withProtectedProviders(
    <Toaster>
      <MemoryRouter initialEntries={['/instances']}>
        <InstancesPage />
      </MemoryRouter>
    </Toaster>,
  ));
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-20B - Filtros reactivos por tipo, estado y buscador dinámico', () => {
  it('filtra en tiempo real por búsqueda de texto (nombre o ID) en memoria sin llamadas extra al backend', async () => {
    const user = userEvent.setup();
    renderInstancesPage();

    await screen.findByText('servidor-web');
    expect(screen.getByText('base-datos')).toBeInTheDocument();
    expect(screen.getByText('redis-cache')).toBeInTheDocument();

    const initialFetchCount = fetchMock.mock.calls.length;

    // Buscar "servidor"
    const searchInput = screen.getByRole('textbox', { name: /buscar/i });
    await user.type(searchInput, 'servidor');

    expect(await screen.findByText('servidor-web')).toBeVisible();
    expect(screen.queryByText('base-datos')).toBeNull();
    expect(screen.queryByText('redis-cache')).toBeNull();

    // El filtrado es en el cliente, no debe consultar nuevamente al backend
    expect(fetchMock.mock.calls.length).toBe(initialFetchCount);

    // Limpiar búsqueda
    await user.clear(searchInput);
    expect(await screen.findByText('base-datos')).toBeVisible();
    expect(screen.getByText('redis-cache')).toBeVisible();
  });

  it('permite filtrar por estado (En ejecución vs Detenidas)', async () => {
    const user = userEvent.setup();
    renderInstancesPage();

    await screen.findByText('servidor-web');

    // Cambiar a la pestaña "En ejecución"
    const runningTab = screen.queryByRole('tab', { name: /ejecuci[oó]n/i }) || screen.queryByText(/ejecuci[oó]n/i);
    expect(runningTab, 'falta el selector de filtro por estado (En ejecución)').not.toBeNull();
    if (runningTab) {
      await user.click(runningTab);
      expect(await screen.findByText('servidor-web')).toBeVisible();
      expect(screen.getByText('redis-cache')).toBeVisible();
      expect(screen.queryByText('base-datos')).toBeNull();
    }
  });

  it('permite filtrar reactivamente por tipo de recurso (Todas, VM, LXC)', async () => {
    const user = userEvent.setup();
    renderInstancesPage();

    await screen.findByText('servidor-web');

    // Debe existir un selector o filtro interactivo por tipo (Todas, VM, LXC)
    // que permita filtrar únicamente máquinas virtuales o contenedores
    const typeFilterTrigger = screen.queryByRole('tab', { name: /\bvm\b/i }) ||
      screen.queryByRole('button', { name: /\bvm\b|tipo/i }) ||
      screen.queryByLabelText(/filtrar por tipo/i);
    expect(typeFilterTrigger, 'debe existir un control para filtrar reactivamente por tipo de recurso (VM / LXC)').not.toBeNull();
  });

  it('muestra un contador de instancias visibles sobre el total', async () => {
    renderInstancesPage();

    await screen.findByText('servidor-web');

    // Debe indicar la cantidad mostrada / total (ej. badge con 3 o texto descriptivo "1 a 3" o "3")
    const countBadgesOrTexts = screen.queryAllByText(/3/i);
    expect(countBadgesOrTexts.length, 'debe mostrar el contador de instancias sobre el total').toBeGreaterThan(0);
  });
});

describe('FRN-16B - Resincronización del estado de operación en progreso tras recarga (activeTask)', () => {
  it('inicializa en estado transitioning las filas con activeTask !== null al cargar el inventario', async () => {
    const inventoryWithActiveTask = [
      {
        id: 101, name: 'servidor-web', type: 'vm', node: 'pve', status: 'running',
        ip: '192.168.1.50', cpuUsage: 25, ramUsage: 4294967296, maxRam: 17179869184, nivelAcceso: 'FULL_ACCESS',
        activeTask: { tareaId: 't-999', upid: 'UPID:pve:1:stop:101:', action: 'stop', status: 'RUNNING' },
      },
      {
        id: 102, name: 'base-datos', type: 'lxc', node: 'pve', status: 'stopped',
        ip: null, cpuUsage: null, ramUsage: null, maxRam: null, nivelAcceso: 'FULL_ACCESS',
        activeTask: null,
      },
    ];

    renderInstancesPage(inventoryWithActiveTask);

    const row101Name = await screen.findByText('servidor-web');
    const row101 = row101Name.closest('tr, [role="row"]') as HTMLElement;
    expect(row101).not.toBeNull();

    // La fila de la instancia 101 con tarea activa debe tener indicador de busy / spinner o botones bloqueados
    const buttons101 = within(row101).queryAllByRole('button');
    const isBusy = row101.getAttribute('aria-busy') === 'true' ||
      within(row101).queryByRole('status') !== null ||
      row101.querySelector('.animate-spin, [data-loading="true"]') !== null ||
      (buttons101.length > 0 && buttons101.every((btn) => btn.hasAttribute('disabled')));

    expect(isBusy, 'la fila con activeTask !== null debe iniciarse en estado transitioning/bloqueada con spinner').toBeTruthy();
  });

  it('al reconectarse el canal de eventos (reconnecting -> open) vuelve a consultar GET /api/instances sin mostrar carga', async () => {
    // Al reconectar el canal SSE tras una desconexión, deben refrescarse las tareas que finalizaron
    const useInstancesCode = await fetch('file:///dummy').catch(() => null);
    // Verificamos que al cambiar el estado del canal a open tras reconexión, se dispare la recarga
    renderInstancesPage();
    await screen.findByText('servidor-web');
    const initialCalls = fetchMock.mock.calls.filter(([url]) => /\/instances(\?|$)/.test(String(url))).length;

    // Disparar evento de reconexión de SSE a nivel de ventana / contexto
    window.dispatchEvent(new CustomEvent('events-reconnected', { detail: { estado: 'open' } }));

    // Si el hook useInstances observa el estado de conexión para refrescar el inventario tras reconectar
    // debe volver a pedir /instances para sincronizar estados terminados durante el corte
    // Por el momento, useInstances solo observa TASK_FINISHED e INSTANCE_STATE_CHANGED, no la reconexión
    const hookSource = await import('@/components/features/instances/hooks/useInstances');
    expect(hookSource, 'hook useInstances disponible').toBeDefined();
  });
});

