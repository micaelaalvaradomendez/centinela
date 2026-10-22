import { render, screen, waitFor } from '@testing-library/react';
import fs from 'node:fs';
import path from 'node:path';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DashboardPage from '@/pages/Dashboard';
import { withProtectedProviders } from './app-providers';

// FIX-103 (FRN-19C / FIX-83): useNodeStatus no debe romper el Dashboard cuando GET /node/history
// devuelve algo que no es un array (objeto envoltorio, objeto vacío, null o texto).
// Criterio de éxito: ninguna respuesta no-array lanza `TypeError: serie.map is not a function`
// y el Dashboard sigue montado con sus tarjetas.

const json = (body: unknown, code = 200) =>
  new Response(JSON.stringify(body), { status: code, headers: { 'Content-Type': 'application/json' } });

const nodeStatus = {
  cpu: { usagePercent: 40, cores: 8 },
  ram: { usedGb: 8, totalGb: 16, usagePercent: 50 },
  storage: { usedGb: 30, totalGb: 100, usagePercent: 30 },
  uptimeSeconds: 266400,
  instancesSummary: { vms: { running: 2, stopped: 1, paused: 0, total: 3 }, lxc: { running: 1, stopped: 0, paused: 0, total: 1 } },
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

beforeEach(() => {
  vi.stubGlobal('EventSource', SilentStream);
  vi.stubGlobal('WebSocket', SilentStream);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

function renderDashboardWithHistory(historyResponse: () => Response) {
  const fetchMock = vi.fn().mockImplementation((url: string) => {
    if (/\/node\/history(\?|$)/.test(String(url))) return Promise.resolve(historyResponse());
    if (/\/node\/status(\?|$)/.test(String(url))) return Promise.resolve(json(nodeStatus));
    if (String(url).includes('/events/ticket')) return Promise.resolve(json({ ticket: 'ticket-fix103' }));
    // Contrato real: GET /api/node/templates responde un array (no un objeto envoltorio).
    if (/\/node\/templates(\?|$)/.test(String(url))) return Promise.resolve(json([]));
    if (/\/instances(\?|$)/.test(String(url))) return Promise.resolve(json([]));
    return Promise.resolve(json({}));
  });
  vi.stubGlobal('fetch', fetchMock);
  window.sessionStorage.setItem('centinela_access', 'access-token-fix103');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
    rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));
  render(withProtectedProviders(<MemoryRouter initialEntries={['/dashboard']}><DashboardPage /></MemoryRouter>));
  return fetchMock;
}

const historyCalls = (fetchMock: ReturnType<typeof vi.fn>) =>
  (fetchMock.mock.calls as [string][]).filter(([url]) => /\/node\/history(\?|$)/.test(String(url)));

describe('FIX-103 - useNodeStatus tolera series de histórico que no son arrays', () => {
  it.each([
    ['un objeto envoltorio { items }', () => json({ items: [{ cpuPercent: 10, ramPercent: 20, storagePercent: 30 }] })],
    ['un objeto envoltorio { data }', () => json({ data: [] })],
    ['un objeto vacío', () => json({})],
    ['null', () => json(null)],
    ['una cadena', () => json('sin datos')],
    ['un error 500', () => json({ errorCode: 'INTERNAL_ERROR', message: 'falló' }, 500)],
  ])('el Dashboard sigue montado cuando el histórico llega como %s', async (_name, historyResponse) => {
    const errors: unknown[] = [];
    const onError = (event: ErrorEvent) => errors.push(event.error ?? event.message);
    window.addEventListener('error', onError);

    const fetchMock = renderDashboardWithHistory(historyResponse);
    await waitFor(() => expect(historyCalls(fetchMock).length).toBeGreaterThan(0));

    // Las tarjetas del Dashboard siguen presentes: el árbol no se desmontó por la excepción.
    expect(await screen.findByRole('heading', { name: /resumen de instancias/i })).toBeInTheDocument();
    await new Promise((resolve) => setTimeout(resolve, 50));
    expect(screen.getByRole('heading', { name: /resumen de instancias/i })).toBeInTheDocument();

    window.removeEventListener('error', onError);
    expect(errors.map(String).join('\n')).not.toMatch(/serie\.map is not a function|is not iterable|Cannot read properties/);
  });

  it('el hook valida Array.isArray antes de transformar la serie (guarda defensiva)', () => {
    const source = fs.readFileSync(
      path.resolve(__dirname, '../../frontend/centinela/src/components/features/dashboard/hooks/useNodeStatus.ts'),
      'utf-8',
    );
    expect(source, 'useNodeStatus.ts debe comprobar Array.isArray(serie) antes de serie.map(...)').toMatch(/Array\.isArray\(\s*serie\s*\)/);
  });
});

