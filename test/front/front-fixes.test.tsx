import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';
import Auditoria from '@/pages/Auditoria';
import { TwoFactorPage } from '@/pages/TwoFactor';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

// Siembra una sesión válida antes de navegar (mismo patrón que admin-users.test.tsx).
function renderAppAs(rol: 'ADMIN' | 'OPERATOR', path: string) {
  window.sessionStorage.setItem('centinela_access', 'access-token');
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
  render(<RouterProvider router={router} />);
  return router;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('F1 - UserDetail carga el usuario real', () => {
  it('carga el usuario por GET /admin/users/:id y guarda los datos reales en el PUT', async () => {
    const fetchMock = vi.fn().mockImplementation((_url: string, init?: RequestInit) => {
      if (init?.method === 'PUT') {
        return Promise.resolve(jsonResponse({ id: 'u2', nombreCompleto: 'Ada Lovelace' }));
      }
      return Promise.resolve(jsonResponse({
        id: 'u2',
        nombreCompleto: 'Ada Lovelace',
        nombreUsuario: 'ada',
        emailUsuario: 'ada@centinela.local',
        rol: 'ADMIN',
        activo: true,
      }));
    });
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderAppAs('ADMIN', '/users/u2');

    expect(await screen.findByDisplayValue('Ada Lovelace')).toBeVisible();
    expect(screen.queryByDisplayValue('Usuario Dos')).not.toBeInTheDocument();

    const saveButton = screen.getByRole('button', { name: /guardar cambios/i });
    await waitFor(() => expect(saveButton).toBeEnabled());
    await user.click(saveButton);

    await waitFor(() => {
      const calls = fetchMock.mock.calls as [string, RequestInit][];
      const putCall = calls.find(([, init]) => init.method === 'PUT');
      expect(putCall).toBeDefined();
      const payload = JSON.parse(String(putCall?.[1].body));
      expect(payload).toMatchObject({
        nombreCompleto: 'Ada Lovelace',
        emailUsuario: 'ada@centinela.local',
        rol: 'ADMIN',
      });
      expect(JSON.stringify(payload)).not.toContain('Usuario Dos');
    });
  });
});

describe('F2 - guard ADMIN en las rutas de usuarios', () => {
  it('redirige a /dashboard cuando un OPERATOR intenta entrar a /users', async () => {
    const router = renderAppAs('OPERATOR', '/users');

    await waitFor(() => expect(router.state.location.pathname).toBe('/dashboard'));
    expect(screen.queryByRole('heading', { name: 'Gestión de usuarios' })).not.toBeInTheDocument();
  });

  it('renderiza el panel de gestión cuando el usuario es ADMIN', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      summary: { total: 0, admins: 0, operators: 0 },
      users: [],
    }));
    vi.stubGlobal('fetch', fetchMock);

    const router = renderAppAs('ADMIN', '/users');

    expect(await screen.findByRole('heading', { name: 'Gestión de usuarios' })).toBeVisible();
    expect(router.state.location.pathname).toBe('/users');
  });
});

describe('F3 - mensaje de error de verificación 2FA', () => {
  it('muestra el mensaje mapeado ante un verify fallido, no el genérico', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(
      { errorCode: 'TOTP_FAILED', message: 'mensaje crudo del backend' },
      401,
    ));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    const router = createMemoryRouter([
      {
        path: '/two-factor/verify',
        Component: TwoFactorPage,
        loader: () => ({ jwtTemporal: 'jwt-temporal', totpVinculado: true, cambioContrasenaRequerido: false }),
      },
    ], { initialEntries: ['/two-factor/verify'] });
    render(<RouterProvider router={router} />);

    await user.type(await screen.findByRole('textbox'), '123456');
    await user.click(screen.getByRole('button', { name: 'Verificar código' }));

    const alert = await screen.findByRole('alert');
    expect(alert).toHaveTextContent('El código de verificación es incorrecto o ya se usó.');
    expect(screen.queryByText('No se pudo completar la operación. Intentá nuevamente.')).not.toBeInTheDocument();
  });
});

describe('F5 - Auditoría consume la API real', () => {
  it('renderiza los registros que devuelve GET /admin/audit en vez de datos hardcodeados', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      total: 1,
      pagina: 1,
      items: [{
        id: 'a1',
        usuarioId: 'u1',
        nombreUsuario: 'Ada Lovelace',
        fechaHora: '2026-09-21T12:00:00Z',
        accion: 'LOGIN',
        instanciaId: '',
        instanciaNombre: '',
        resultado: 'EXITO',
        detalles: '{}',
      }],
    }));
    vi.stubGlobal('fetch', fetchMock);

    render(<Auditoria />);

    const table = await screen.findByRole('table');
    expect(await within(table).findByText('Ada Lovelace')).toBeVisible();
    expect(within(table).getByText('LOGIN')).toBeVisible();
    expect(screen.queryByText('AD Admin')).not.toBeInTheDocument();
    expect(screen.queryByText('Tomo snapshot')).not.toBeInTheDocument();

    const [url] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toContain('/admin/audit');
    expect(url).toContain('pagina=1');
    expect(url).toContain('tamano=10');
  });
});
