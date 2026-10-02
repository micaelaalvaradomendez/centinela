import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';
import { withAppProviders } from './app-providers';
import Auditoria from '@/pages/Auditoria';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

function renderAppAs(rol: 'ADMIN' | 'OPERATOR', path: string) {
  window.sessionStorage.setItem('centinela_access', 'access-token-audit-123');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-1',
    organizacionId: 'org-1',
    nombreCompleto: rol === 'ADMIN' ? 'Administradora de prueba' : 'Operador de prueba',
    email: 'sesion@centinela.local',
    rol,
    instanciasPermitidas: [],
    tiene2FA: true,
  }));
  const router = createMemoryRouter(applicationRoutes, { initialEntries: [path] });
  render(withAppProviders(<RouterProvider router={router} />));
  return router;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-14 - Suite de pruebas para vista de Auditoría (Auditoria.tsx)', () => {
  it('realiza la petición GET /api/admin/audit?pagina=1&tamano=10 con cabecera Authorization: Bearer', async () => {
    window.sessionStorage.setItem('centinela_access', 'token-bearer-test');
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      total: 0,
      pagina: 1,
      items: [],
    }));
    vi.stubGlobal('fetch', fetchMock);

    render(<Auditoria />);

    await waitFor(() => expect(fetchMock).toHaveBeenCalled());
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toContain('/admin/audit');
    expect(url).toContain('pagina=1');
    expect(url).toContain('tamano=10');
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer token-bearer-test');
  });

  it('renderiza registros devueltos por la API real en lugar de datos estáticos', async () => {
    window.sessionStorage.setItem('centinela_access', 'token-bearer-test');
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      total: 1,
      pagina: 1,
      items: [{
        id: 'aud-1',
        usuarioId: 'usr-1',
        nombreUsuario: 'Ada Lovelace',
        fechaHora: '2026-09-21T12:00:00Z',
        accion: 'LOGIN',
        instanciaId: '',
        instanciaNombre: '',
        resultado: 'EXITO',
        detalles: '{"ip":"192.168.1.10"}',
      }],
    }));
    vi.stubGlobal('fetch', fetchMock);

    render(<Auditoria />);

    const table = await screen.findByRole('table');
    expect(await within(table).findByText('Ada Lovelace')).toBeVisible();
    expect(within(table).getByText('LOGIN')).toBeVisible();
    expect(within(table).getByText(/éxito|exito/i)).toBeVisible();
  });

  it('actualiza los query parameters de filtro al cambiar los selectores reactivos', async () => {
    window.sessionStorage.setItem('centinela_access', 'token-bearer-test');
    const fetchMock = vi.fn().mockImplementation((_url: string) => {
      return Promise.resolve(jsonResponse({
        total: 0,
        pagina: 1,
        items: [],
      }));
    });
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    render(<Auditoria />);
    await waitFor(() => expect(fetchMock).toHaveBeenCalled());

    // Cambiar filtro de acción
    const accionInput = screen.getByLabelText(/acción/i);
    await user.type(accionInput, 'LOGIN');

    await waitFor(() => {
      const calls = fetchMock.mock.calls as [string, RequestInit][];
      const filteredCall = calls.find(([url]) => url.includes('accion=LOGIN'));
      expect(filteredCall).toBeDefined();
    });
  });

  it('el filtro de resultado y el rango de fechas se envían como query params', async () => {
    window.sessionStorage.setItem('centinela_access', 'token-bearer-test');
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ total: 0, pagina: 1, items: [] }));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    render(<Auditoria />);
    await waitFor(() => expect(fetchMock).toHaveBeenCalled());

    await user.selectOptions(screen.getByLabelText('Resultado'), 'FALLA');
    await waitFor(() => expect((fetchMock.mock.calls as [string][]).some(([url]) => url.includes('resultado=FALLA'))).toBe(true));

    fireEvent.change(screen.getByLabelText('Desde'), { target: { value: '2026-09-01' } });
    fireEvent.change(screen.getByLabelText('Hasta'), { target: { value: '2026-09-30' } });
    await waitFor(() => {
      const last = String((fetchMock.mock.calls as [string][]).at(-1)?.[0]);
      expect(last).toContain('resultado=FALLA');
      expect(last).toMatch(/desde=2026-09-01/);
      expect(last).toMatch(/hasta=2026-09-30/);
    });
  });

  it('la paginación solicita la página siguiente al backend', async () => {
    window.sessionStorage.setItem('centinela_access', 'token-bearer-test');
    const fetchMock = vi.fn().mockImplementation((url: string) => Promise.resolve(jsonResponse({
      total: 25,
      pagina: url.includes('pagina=2') ? 2 : 1,
      items: [{ id: `aud-${url.includes('pagina=2') ? 2 : 1}`, usuarioId: 'usr-1', nombreUsuario: 'Ada Lovelace', fechaHora: '2026-09-21T12:00:00Z', accion: 'LOGIN', resultado: 'EXITO' }],
    })));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    render(<Auditoria />);
    await screen.findByRole('table');
    await user.click(screen.getByRole('button', { name: 'Página siguiente' }));

    await waitFor(() => expect((fetchMock.mock.calls as [string][]).some(([url]) => url.includes('pagina=2'))).toBe(true));
  });

  it('el botón Exportar CSV dispara la llamada a /admin/audit/export?formato=csv con Bearer token', async () => {
    window.sessionStorage.setItem('centinela_access', 'token-bearer-test');
    const fetchMock = vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/export')) {
        return Promise.resolve(new Response('col1,col2\nval1,val2', {
          status: 200,
          headers: { 'Content-Type': 'text/csv' },
        }));
      }
      return Promise.resolve(jsonResponse({ total: 0, pagina: 1, items: [] }));
    });
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    render(<Auditoria />);
    await waitFor(() => expect(fetchMock).toHaveBeenCalled());

    const exportButton = await screen.findByRole('button', { name: /exportar/i });
    await user.click(exportButton);

    await waitFor(() => {
      const calls = fetchMock.mock.calls as [string, RequestInit][];
      const exportCall = calls.find(([url]) => url.includes('/admin/audit/export') && url.includes('formato=csv'));
      expect(exportCall).toBeDefined();
      expect((exportCall?.[1].headers as Record<string, string>)?.Authorization).toBe('Bearer token-bearer-test');
    });
  });

  it('si un operador intenta entrar a /auditoria, el guard administrativo lo redirige a /dashboard', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({ total: 0, pagina: 1, items: [] })));
    const router = renderAppAs('OPERATOR', '/auditoria');

    await waitFor(() => expect(router.state.location.pathname).toBe('/dashboard'));
    expect(screen.queryByRole('heading', { name: 'Auditoría' })).not.toBeInTheDocument();
  });
});
