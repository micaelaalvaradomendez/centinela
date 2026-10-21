import { render, screen, waitFor } from '@testing-library/react';
import fs from 'node:fs';
import path from 'node:path';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DashboardPage from '@/pages/Dashboard';
import { withProtectedProviders } from './app-providers';

// FIX-117B (depende de FIX-117A): con la segregación de telemetría, GET /node/status devuelve a un
// OPERATOR una respuesta acotada: sin `storage` global y con `instancesSummary` limitado a sus
// instancias asignadas. El Dashboard no debe lanzar `TypeError: undefined reading ...` ni mostrar NaN.
// Criterio de éxito: el Dashboard de un operador renderiza limpiamente, sin errores de consola.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

// Respuesta acotada para OPERATOR: solo cuenta lo que tiene asignado (1 VM encendida).
const scopedStatus = {
  cpu: { usagePercent: 12, cores: 8 },
  ram: { usedGb: 2, totalGb: 16, usagePercent: 12 },
  uptimeSeconds: 3600,
  instancesSummary: { vms: { running: 1, stopped: 0, paused: 0, total: 1 }, lxc: { running: 0, stopped: 0, paused: 0, total: 0 } },
  stale: false,
  fetchedAt: '2026-10-13T03:00:00Z',
};

class SilentStream {
  onerror: unknown = null; onclose: unknown = null; onmessage: unknown = null; onopen: unknown = null;
  readyState = 1;
  constructor(public url: string) {}
  addEventListener() {}
  removeEventListener() {}
  close() { this.readyState = 2; }
}

let consoleError: ReturnType<typeof vi.spyOn>;

beforeEach(() => {
  vi.stubGlobal('EventSource', SilentStream);
  vi.stubGlobal('WebSocket', SilentStream);
  consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

function renderDashboardAs(rol: 'ADMIN' | 'OPERATOR', statusBody: unknown) {
  const fetchMock = vi.fn().mockImplementation((url: string) => {
    if (/\/node\/status(\?|$)/.test(String(url))) return Promise.resolve(json(statusBody));
    if (/\/node\/history(\?|$)/.test(String(url))) return Promise.resolve(json([]));
    if (String(url).includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-fix117b' }));
    return Promise.resolve(json({}));
  });
  vi.stubGlobal('fetch', fetchMock);
  window.sessionStorage.setItem('centinela_access', 'access-token-fix117b');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-1', organizacionId: 'org-1', nombreCompleto: 'Usuario Test', email: 'usuario@centinela.local',
    rol, instanciasPermitidas: [101], permisos: [{ vmid: 101, nivelAcceso: 'FULL_ACCESS' }], tiene2FA: true,
  }));
  render(withProtectedProviders(<MemoryRouter initialEntries={['/dashboard']}><DashboardPage /></MemoryRouter>));
  return fetchMock;
}

describe('FIX-117B - Dashboard de operador ante métricas globales segregadas', () => {
  it('renderiza las tarjetas de instancias acotadas sin storage global y sin errores de consola', async () => {
    const fetchMock = renderDashboardAs('OPERATOR', scopedStatus);
    await waitFor(() => expect(
      (fetchMock.mock.calls as [string][]).some(([url]) => /\/node\/status/.test(String(url))),
    ).toBe(true));

    const total = await screen.findByRole('heading', { name: /instancias totales/i });
    expect(total.closest('div')?.textContent ?? '').toMatch(/^\s*1/);

    expect(document.body.textContent ?? '').not.toMatch(/NaN|undefined/);
    const messages = consoleError.mock.calls.map((call) => call.map(String).join(' ')).join('\n');
    expect(messages, 'el Dashboard de operador escribió errores en consola').not.toMatch(/TypeError|Cannot read properties|undefined/);
  });

  it('no estima la salud del nodo con un storage inexistente ni pinta el medidor de almacenamiento global', async () => {
    renderDashboardAs('OPERATOR', scopedStatus);
    await screen.findByRole('heading', { name: /instancias totales/i });

    expect(screen.queryByText(/inaccesible|no se pudo obtener/i), 'el Dashboard cayó al estado de error al faltar storage').toBeNull();
    expect(screen.queryByRole('meter', { name: /almacenamiento/i })).toBeNull();
    expect(screen.queryByText(/almacenamiento del (host|nodo)/i)).toBeNull();
  });

  it('el ADMIN conserva el medidor de almacenamiento del host', async () => {
    renderDashboardAs('ADMIN', { ...scopedStatus, storage: { usedGb: 30, totalGb: 100, usagePercent: 30 } });
    await screen.findByRole('heading', { name: /instancias totales/i });
    expect(await screen.findAllByText(/almacenamiento/i)).not.toHaveLength(0);
  });

  it('useNodeStatus no accede a data.storage sin comprobar que exista', () => {
    const source = fs.readFileSync(
      path.resolve(__dirname, '../../frontend/centinela/src/components/features/dashboard/hooks/useNodeStatus.ts'),
      'utf-8',
    );
    expect(source, 'useNodeStatus.ts usa data.storage.usagePercent sin optional chaining').not.toMatch(/data\.storage\.usagePercent/);
  });
});

