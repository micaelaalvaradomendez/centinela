import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstancesPage from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { withAppProviders } from './app-providers';

// FRN-15: Modales de confirmación antierror para acciones operativas (RF-04 / Etapa 1).
// Entregables:
// - Start: confirmación estándar antes de disparar. Cancelar no envía peticiones.
// - Shutdown / Reboot: aviso de apagado o reinicio del sistema operativo huésped.
// - Stop: advertencia visual en rojo / advertencia sobre posible pérdida de datos.
// - Delete: modal destructivo que exige tipear el ID o el nombre para confirmar.
// - Control de acceso: Delete es SOLO visible para ADMIN. OPERATOR nunca ve Delete.
//   OPERATOR con READ_ONLY no puede operar acciones de energía; con FULL_ACCESS sí.
const inventoryMock = [
  {
    id: 101, name: 'servidor-web', type: 'vm', node: 'pve', status: 'running',
    ip: '192.168.1.50', cpuUsage: 25, ramUsage: 4, maxRam: 16, nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
  {
    id: 102, name: 'base-datos', type: 'lxc', node: 'pve', status: 'stopped',
    ip: null, cpuUsage: null, ramUsage: null, maxRam: null, nivelAcceso: 'READ_ONLY', activeTask: null,
  },
];

let fetchMock: ReturnType<typeof vi.fn>;

function renderInstancesPage(userRole = 'ADMIN', userPermissions: { vmid: number; nivelAcceso: string }[] = []) {
  window.sessionStorage.setItem('centinela_access', 'access-token-modales');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-test', organizacionId: 'org-1', nombreCompleto: 'Test User', email: 'test@centinela.local',
    rol: userRole, instanciasPermitidas: userPermissions.map(p => p.vmid), permisos: userPermissions, tiene2FA: true,
  }));
  return render(withAppProviders(
    <Toaster>
      <MemoryRouter initialEntries={['/instances']}>
        <InstancesPage />
      </MemoryRouter>
    </Toaster>,
  ));
}

async function rowOf(name: string) {
  const cell = await screen.findByText(new RegExp(name, 'i'));
  const row = cell.closest('tr, [role="row"]');
  expect(row, `"${name}" no está dentro de una fila de tabla`).not.toBeNull();
  return row as HTMLElement;
}

beforeEach(() => {
  fetchMock = vi.fn().mockImplementation((url: string, init?: RequestInit) => {
    const urlStr = String(url);
    if (urlStr.includes('/instances') && (!init || init.method === 'GET' || !init.method)) {
      return Promise.resolve(new Response(JSON.stringify(inventoryMock), {
        status: 200, headers: { 'Content-Type': 'application/json' },
      }));
    }
    return Promise.resolve(new Response(JSON.stringify({ upid: 'UPID:pve:123', tareaId: '0192-task' }), {
      status: 202, headers: { 'Content-Type': 'application/json' },
    }));
  });
  vi.stubGlobal('fetch', fetchMock);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-15 - Modales de confirmación antierror para acciones operativas', () => {
  it('Start abre un modal de confirmación y cancelar no envía peticiones al backend', async () => {
    const user = userEvent.setup();
    renderInstancesPage('ADMIN');

    const dbRow = await rowOf('base-datos');
    const startBtn = within(dbRow).queryByRole('button', { name: /start|iniciar|encender/i });
    expect(startBtn, 'no se encontró botón para iniciar la instancia detenida').not.toBeNull();

    await user.click(startBtn!);

    // Debe abrirse el diálogo/modal de confirmación
    const dialog = await screen.findByRole('dialog');
    expect(dialog).toBeInTheDocument();
    expect(within(dialog).getByText(/iniciar|encender|confirmar/i)).toBeInTheDocument();

    // Cancelar no dispara peticiones POST hacia la API
    const cancelBtn = within(dialog).getByRole('button', { name: /cancelar|volver/i });
    await user.click(cancelBtn);

    const postCalls = (fetchMock.mock.calls as [string, RequestInit?][]).filter(
      ([, init]) => init?.method === 'POST' || init?.method === 'DELETE',
    );
    expect(postCalls.length, 'cancelar el modal envió una petición no deseada').toBe(0);
  });

  it('Stop muestra un modal de advertencia destacando riesgo de pérdida de datos', async () => {
    const user = userEvent.setup();
    renderInstancesPage('ADMIN');

    const webRow = await rowOf('servidor-web');
    const stopBtn = within(webRow).queryByRole('button', { name: /stop|detener|apagar/i });
    expect(stopBtn, 'no se encontró botón para detener la instancia en ejecución').not.toBeNull();

    await user.click(stopBtn!);

    const dialog = await screen.findByRole('dialog');
    expect(dialog).toBeInTheDocument();
    // Debe advertir sobre pérdida de datos o apagado forzoso
    const warningText = dialog.textContent || '';
    expect(warningText).toMatch(/p[eé]rdida|datos|forzoso|inmediat|riesgo/i);
  });

  it('Shutdown y Reboot advierten sobre el apagado/reinicio del sistema huésped', async () => {
    const user = userEvent.setup();
    renderInstancesPage('ADMIN');

    const webRow = await rowOf('servidor-web');
    const actionBtn = within(webRow).queryByRole('button', { name: /reboot|reiniciar|shutdown|apagar/i });
    if (!actionBtn) {
      throw new Error('FRN-15 no implementada: la botonera no expone acciones de shutdown o reboot');
    }

    await user.click(actionBtn);
    const dialog = await screen.findByRole('dialog');
    expect(dialog).toBeInTheDocument();
    expect(dialog.textContent).toMatch(/sistema|hu[eé]sped|guest|ordenad/i);
  });

  it('Delete es un modal destructivo que exige ingresar el nombre o ID, y es exclusivo de ADMIN', async () => {
    const user = userEvent.setup();

    // 1. Un OPERATOR nunca debe ver la acción Delete
    const { unmount } = renderInstancesPage('OPERATOR', [{ vmid: 101, nivelAcceso: 'FULL_ACCESS' }]);
    const opRow = await rowOf('servidor-web');
    const opDeleteBtn = within(opRow).queryByRole('button', { name: /delete|eliminar|borrar/i });
    expect(opDeleteBtn, 'un OPERATOR no debe ver el botón Delete').toBeNull();
    unmount();

    // 2. El ADMIN sí tiene la acción Delete y exige confirmación tipeada
    renderInstancesPage('ADMIN');
    const adminRow = await rowOf('servidor-web');
    const deleteBtn = within(adminRow).queryByRole('button', { name: /delete|eliminar|borrar/i });
    expect(deleteBtn, 'el ADMIN debe contar con la opción Delete').not.toBeNull();

    await user.click(deleteBtn!);
    const dialog = await screen.findByRole('dialog');
    expect(dialog).toBeInTheDocument();

    // Requiere tipear el nombre o el ID
    const inputConfirm = within(dialog).getByRole('textbox');
    expect(inputConfirm).toBeInTheDocument();

    const confirmDeleteBtn = within(dialog).getByRole('button', { name: /eliminar|borrar|confirmar/i });
    expect(confirmDeleteBtn).toBeDisabled();

    // Al tipear el nombre correcto se habilita el botón
    await user.type(inputConfirm, 'servidor-web');
    expect(confirmDeleteBtn).not.toBeDisabled();
  });

  it('OPERATOR con READ_ONLY no puede operar energía; con FULL_ACCESS sí', async () => {
    renderInstancesPage('OPERATOR', [
      { vmid: 101, nivelAcceso: 'FULL_ACCESS' },
      { vmid: 102, nivelAcceso: 'READ_ONLY' },
    ]);

    const webRow = await rowOf('servidor-web');
    const dbRow = await rowOf('base-datos');

    // 101 (FULL_ACCESS): puede operar
    const webActionBtn = within(webRow).queryByRole('button', { name: /stop|detener|reboot|reiniciar/i });
    expect(webActionBtn, 'un operador con FULL_ACCESS debe tener habilitadas las acciones de energía').toBeEnabled();

    // 102 (READ_ONLY): botones de energía deshabilitados o ausentes
    const dbActionBtn = within(dbRow).queryByRole('button', { name: /start|iniciar|encender/i });
    if (dbActionBtn) {
      expect(dbActionBtn, 'un operador con READ_ONLY no debe tener botón de start habilitado').toBeDisabled();
    }
  });
});
