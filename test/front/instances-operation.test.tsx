import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstancesPage from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { withAppProviders } from './app-providers';

// FRN-16: Máquina de estados "Operación en progreso" por instancia (RF-04 / Etapa 1).
// Entregables:
// - Al enviar una orden (start, stop, reboot, shutdown, delete), la fila entra en estado `transitioning`.
// - El botón accionado muestra un spinner y se deshabilitan todos los botones de esa instancia (antidoble-clic).
// - Ante error HTTP, desbloquea la fila y muestra un mensaje distinto según el código (D2):
//   - 409 INSTANCE_INVALID_STATE
//   - 409 INSTANCE_BUSY
//   - 403 INSTANCE_PROTECTED
//   - 403 INSTANCE_ACCESS_DENIED
//   - 502 PROXMOX_UNAVAILABLE
//   - 504 PROXMOX_TIMEOUT

const inventory = [
  {
    id: 101, name: 'servidor-web', type: 'vm', node: 'pve', status: 'running',
    ip: '192.168.1.50', cpuUsage: 30, ramUsage: 4, maxRam: 16, nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
];

let fetchMock: ReturnType<typeof vi.fn>;

function renderInstancesPage() {
  window.sessionStorage.setItem('centinela_access', 'access-token-operation');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', rol: 'ADMIN', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
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

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-16 - Máquina de estados "Operación en progreso" por instancia', () => {
  it('deshabilita los botones de la fila y muestra spinner durante una operación en curso', async () => {
    let resolvePost: (value: Response) => void;
    const postPromise = new Promise<Response>((resolve) => { resolvePost = resolve; });

    fetchMock = vi.fn().mockImplementation((url: string, init?: RequestInit) => {
      if (String(url).includes('/instances') && (!init || init.method === 'GET' || !init.method)) {
        return Promise.resolve(new Response(JSON.stringify(inventory), { status: 200 }));
      }
      return postPromise;
    });
    vi.stubGlobal('fetch', fetchMock);

    const user = userEvent.setup();
    renderInstancesPage();

    const row = await rowOf('servidor-web');
    const stopBtn = within(row).queryByRole('button', { name: /stop|detener|apagar/i });
    expect(stopBtn).not.toBeNull();

    await user.click(stopBtn!);

    // Si abre modal, confirmamos
    const confirmBtn = screen.queryByRole('button', { name: /confirmar|detener|apagar/i });
    if (confirmBtn && confirmBtn !== stopBtn) {
      await user.click(confirmBtn);
    }

    // Durante el vuelo: los botones de la fila deben quedar deshabilitados
    await waitFor(() => {
      const buttons = within(row).getAllByRole('button');
      for (const btn of buttons) {
        expect(btn).toBeDisabled();
      }
    });

    // Debe mostrarse un spinner o estado de carga
    expect(within(row).queryByRole('status') || row.querySelector('.animate-spin, [data-loading]')).not.toBeNull();

    // Completamos la petición
    resolvePost!(new Response(JSON.stringify({ upid: 'UPID:123', tareaId: 't-1' }), { status: 202 }));
  });

  it('muestra mensajes de error específicos según el código devuelto por el backend (D2)', async () => {
    const errorCodes = [
      { status: 409, code: 'INSTANCE_INVALID_STATE', match: /estado|ya est[aá]/i },
      { status: 409, code: 'INSTANCE_BUSY', match: /otra tarea|ocupad|esper[aá]/i },
      { status: 403, code: 'INSTANCE_PROTECTED', match: /protegid|infraestructura/i },
      { status: 403, code: 'INSTANCE_ACCESS_DENIED', match: /permiso|denegad/i },
      { status: 502, code: 'PROXMOX_UNAVAILABLE', match: /no est[aá] disponible|ca[ií]do/i },
      { status: 504, code: 'PROXMOX_TIMEOUT', match: /tiempo|timeout/i },
    ];

    for (const { status, code, match } of errorCodes) {
      fetchMock = vi.fn().mockImplementation((url: string, init?: RequestInit) => {
        if (String(url).includes('/instances') && (!init || init.method === 'GET' || !init.method)) {
          return Promise.resolve(new Response(JSON.stringify(inventory), { status: 200 }));
        }
        return Promise.resolve(new Response(JSON.stringify({ errorCode: code, message: `Error ${code}` }), {
          status, headers: { 'Content-Type': 'application/json' },
        }));
      });
      vi.stubGlobal('fetch', fetchMock);

      const user = userEvent.setup();
      const { unmount } = renderInstancesPage();

      const row = await rowOf('servidor-web');
      const actionBtn = within(row).getByRole('button', { name: /stop|detener|apagar/i });
      await user.click(actionBtn);

      const confirmBtn = screen.queryByRole('button', { name: /confirmar|detener|apagar/i });
      if (confirmBtn && confirmBtn !== actionBtn) {
        await user.click(confirmBtn);
      }

      // Debe aparecer el mensaje toast correspondiente al error
      const toast = await screen.findByRole('status').catch(() => screen.findByText(match));
      expect(toast.textContent).toMatch(match);

      // Y la fila debe haberse desbloqueado tras el fallo
      await waitFor(() => {
        expect(within(row).getByRole('button', { name: /stop|detener|apagar/i })).not.toBeDisabled();
      });

      unmount();
    }
  });
});
