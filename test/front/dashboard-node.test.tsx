import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DashboardPage from '@/pages/Dashboard';
import { withAppProviders } from './app-providers';

// FRN-19B (ex FRN-13B): Semáforo de salud global e integración con GET /api/node/status (RF-02 / Etapa 1).
// Entregables:
// - Conecta el Dashboard a GET /api/node/status mediante un servicio en features/dashboard/services.
// - Semáforo de salud (D3):
//   - "Saludable" si CPU, RAM y disco < 70%.
//   - "Advertencia" si alguno llega a 70% o más.
//   - "Inaccesible" si el backend responde 502/504 o la petición falla.
// - Uptime legible (convierte uptimeSeconds a días, horas y minutos).
// - Aviso o banner de "datos desactualizados" si stale: true.
// - Polling periódico cada 10 s y skeletons de carga mientras se obtienen los datos.

const normalNodeStatus = {
  cpu: { usagePercent: 45, cores: 8 },
  ram: { usedGb: 8.0, totalGb: 16.0, usagePercent: 50 },
  storage: { usedGb: 30.0, totalGb: 100.0, usagePercent: 30 },
  uptimeSeconds: 266400, // 3 días y 2 horas
  instancesSummary: { vms: { running: 2, stopped: 1, paused: 0, total: 3 }, lxc: { running: 1, stopped: 0, paused: 0, total: 1 } },
  stale: false,
  fetchedAt: '2026-10-04T03:00:00Z',
};

const warningNodeStatus = {
  ...normalNodeStatus,
  cpu: { usagePercent: 75, cores: 8 },
};

const staleNodeStatus = {
  ...normalNodeStatus,
  stale: true,
};

let fetchMock: ReturnType<typeof vi.fn>;

function renderDashboard() {
  window.sessionStorage.setItem('centinela_access', 'access-token-dashboard');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', rol: 'ADMIN', nombreCompleto: 'Admin Test',
  }));
  return render(withAppProviders(
    <MemoryRouter initialEntries={['/dashboard']}>
      <DashboardPage />
    </MemoryRouter>,
  ));
}

afterEach(() => {
  vi.unstubAllGlobals();
  vi.useRealTimers();
});

describe('FRN-19B - Semáforo de salud global e integración con GET /api/node/status', () => {
  it('consulta GET /api/node/status con Bearer y muestra el semáforo en estado Saludable', async () => {
    fetchMock = vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/node/status')) {
        return Promise.resolve(new Response(JSON.stringify(normalNodeStatus), {
          status: 200, headers: { 'Content-Type': 'application/json' },
        }));
      }
      return Promise.resolve(new Response('{}', { status: 200 }));
    });
    vi.stubGlobal('fetch', fetchMock);

    renderDashboard();

    // 1. Debe haber consultado el endpoint con el Bearer token
    await waitFor(() => {
      const call = (fetchMock.mock.calls as [string, RequestInit?][]).find(([url]) => /\/node\/status(\?|$)/.test(String(url)));
      expect(call, 'Dashboard no consultó GET /api/node/status').toBeDefined();
      expect(new Headers(call![1]?.headers).get('Authorization')).toBe('Bearer access-token-dashboard');
    });

    // 2. Semáforo saludable (D3)
    const statusBadge = await screen.findByText(/saludable|healthy|normal/i);
    expect(statusBadge).toBeInTheDocument();
  });

  it('conmuta el semáforo a Advertencia cuando la saturación alcanza o supera el 70%', async () => {
    fetchMock = vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/node/status')) {
        return Promise.resolve(new Response(JSON.stringify(warningNodeStatus), {
          status: 200, headers: { 'Content-Type': 'application/json' },
        }));
      }
      return Promise.resolve(new Response('{}', { status: 200 }));
    });
    vi.stubGlobal('fetch', fetchMock);

    renderDashboard();

    const warningBadge = await screen.findByText(/advertencia|warning/i);
    expect(warningBadge).toBeInTheDocument();
  });

  it('muestra Inaccesible ante errores 502/504 o caída de conexión sin romper la UI', async () => {
    fetchMock = vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/node/status')) {
        return Promise.resolve(new Response(JSON.stringify({ errorCode: 'PROXMOX_UNAVAILABLE' }), {
          status: 502, headers: { 'Content-Type': 'application/json' },
        }));
      }
      return Promise.resolve(new Response('{}', { status: 200 }));
    });
    vi.stubGlobal('fetch', fetchMock);

    renderDashboard();

    const inaccessibleBadge = await screen.findByText(/inaccesible|inaccessible|desconectado|sin conexi[oó]n/i);
    expect(inaccessibleBadge).toBeInTheDocument();
  });

  it('muestra el uptime formateado en días y horas, y avisa cuando los datos son desactualizados (stale: true)', async () => {
    fetchMock = vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/node/status')) {
        return Promise.resolve(new Response(JSON.stringify(staleNodeStatus), {
          status: 200, headers: { 'Content-Type': 'application/json' },
        }));
      }
      return Promise.resolve(new Response('{}', { status: 200 }));
    });
    vi.stubGlobal('fetch', fetchMock);

    renderDashboard();

    // 266400 segundos = 3 días y 2 horas
    const uptimeText = await screen.findByText(/3\s*d[ií]as/i);
    expect(uptimeText).toBeInTheDocument();

    // stale: true debe mostrar advertencia de datos desactualizados o en caché
    const staleNotice = await screen.findByText(/desactualizad|stale|cach[eé]/i);
    expect(staleNotice).toBeInTheDocument();
  });
});
