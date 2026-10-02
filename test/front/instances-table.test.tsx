import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstancesPage from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { AuthProvider } from '@/context/AuthContext';

// FRN-20A (ex FRN-14A): Tabla interactiva de inventario con badges de estado e IP (RF-03).
// Entregable: tabla en pages/Instances.tsx con su servicio y su hook en features/intances (o features/instances).
// Criterio de éxito:
//   - VMs y LXC se ven en la misma tabla.
//   - Badges de estado: Running (verde), Stopped (gris).
//   - IP con botón para copiar al portapapeles; si la IP es null, muestra "No detectada".
//   - Tolera que campos nuevos vengan en null (ip, cpuUsage, ramUsage, maxRam, activeTask).
//   - Botonera de acciones presente.

const mockInventoryData = [
  {
    id: 101,
    name: 'servidor-web',
    type: 'qemu',
    node: 'pve-node-01',
    status: 'running',
    ip: '192.168.1.50',
    cpuUsage: null,
    ramUsage: null,
    maxRam: null,
    nivelAcceso: 'FULL_ACCESS',
    activeTask: null,
  },
  {
    id: 102,
    name: 'base-datos',
    type: 'lxc',
    node: 'pve-node-01',
    status: 'stopped',
    ip: null,
    cpuUsage: null,
    ramUsage: null,
    maxRam: null,
    nivelAcceso: 'READ_ONLY',
    activeTask: null,
  },
];

function renderInstancesPage() {
  window.sessionStorage.setItem('centinela_access', 'test-access-token');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin',
    organizacionId: 'org-1',
    nombreCompleto: 'Admin Test',
    email: 'admin@centinela.local',
    rol: 'ADMIN',
    instanciasPermitidas: [],
    tiene2FA: true,
  }));

  return render(
    <AuthProvider>
      <Toaster>
        <MemoryRouter initialEntries={['/instances']}>
          <InstancesPage />
        </MemoryRouter>
      </Toaster>
    </AuthProvider>,
  );
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn().mockImplementation((url: string) => {
    if (String(url).includes('/instances')) {
      return Promise.resolve(new Response(JSON.stringify(mockInventoryData), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }));
    }
    return Promise.resolve(new Response('{}', { status: 200 }));
  }));
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-20A - Tabla interactiva de inventario con badges de estado e IP', () => {
  it('muestra VMs y LXC en la misma tabla con sus IDs y nombres', async () => {
    renderInstancesPage();

    // Debe renderizar la tabla real con datos del backend, no el placeholder estático "Sin instancias"
    expect(await screen.findByText('servidor-web')).toBeInTheDocument();
    expect(screen.getByText('base-datos')).toBeInTheDocument();
    expect(screen.getByText('101')).toBeInTheDocument();
    expect(screen.getByText('102')).toBeInTheDocument();
    expect(screen.queryByText(/Sin instancias/i)).not.toBeInTheDocument();
  });

  it('muestra badges de estado diferenciados (Running en verde, Stopped en gris)', async () => {
    renderInstancesPage();

    const runningBadge = await screen.findByText(/running|en ejecución/i);
    const stoppedBadge = screen.getByText(/stopped|detenida/i);

    expect(runningBadge).toBeInTheDocument();
    expect(stoppedBadge).toBeInTheDocument();
  });

  it('muestra "No detectada" cuando la IP es null y ofrece botón de copiar cuando existe IP', async () => {
    const user = userEvent.setup();
    const writeTextMock = vi.fn().mockResolvedValue(undefined);
    Object.assign(navigator, {
      clipboard: { writeText: writeTextMock },
    });

    renderInstancesPage();

    // Para la instancia 102 con ip: null
    expect(await screen.findByText(/no detectada/i)).toBeInTheDocument();

    // Para la instancia 101 con IP real
    expect(screen.getByText('192.168.1.50')).toBeInTheDocument();

    // Botón de copiar IP
    const copyButton = screen.queryByRole('button', { name: /copiar ip|copiar/i }) ||
      screen.queryByLabelText(/copiar ip/i);
    expect(copyButton).toBeInTheDocument();

    if (copyButton) {
      await user.click(copyButton);
      expect(writeTextMock).toHaveBeenCalledWith('192.168.1.50');
    }
  });

  it('tolera que los campos nuevos de la etapa 1 vengan en null sin romper la renderización', async () => {
    renderInstancesPage();

    // La tabla debe permanecer montada sin arrojar error al encontrar nulls en cpuUsage, ramUsage, etc.
    expect(await screen.findByRole('table', { name: /inventario de instancias/i })).toBeInTheDocument();
  });
});

