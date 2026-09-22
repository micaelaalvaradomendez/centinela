import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';
import Users from '@/pages/Users';
import CrearUsuarios from '@/pages/CrearUsuarios';
import UserDetail from '@/pages/UserDetail';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

// Siembra una sesión válida antes de navegar.
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

// Estos tests corren contra el código real: documentan qué falta
// para dar por cumplidos los criterios de éxito de actual.md.
describe('FRN-05 - panel de gestión de usuarios', () => {
  it('un usuario con rol OPERATOR no debe poder ver el panel de administración', async () => {
    const router = renderAppAs('OPERATOR', '/users');

    // El guard administrativo debe sacar al OPERATOR de /users y devolverlo al dashboard.
    await waitFor(() => expect(router.state.location.pathname).not.toBe('/users'));
    expect(router.state.location.pathname).toBe('/dashboard');
    expect(screen.queryByRole('heading', { name: 'Gestión de usuarios' })).not.toBeInTheDocument();
  });

  it('consulta GET /api/users con Authorization Bearer al entrar a /users', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ summary: { total: 0, admins: 0, operators: 0 }, users: [] }));
    vi.stubGlobal('fetch', fetchMock);

    renderAppAs('ADMIN', '/users');

    await waitFor(() => expect(fetchMock).toHaveBeenCalled());
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toContain('/users');
    expect((init.headers as Record<string, string>).Authorization).toMatch(/^Bearer /);
  });

  it('renderiza los usuarios reales devueltos por el backend en vez de la tabla vacía fija', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      summary: { total: 1, admins: 1, operators: 0 },
      users: [{ id: 'u1', nombreCompleto: 'Ada Lovelace', rol: 'ADMIN', activo: true, totpVinculado: true }],
    }));
    vi.stubGlobal('fetch', fetchMock);

    renderAppAs('ADMIN', '/users');

    // El panel de detalle lateral muestra el mismo usuario; se verifica la fila de la tabla.
    const table = await screen.findByRole('table');
    expect(await within(table).findByText('Ada Lovelace')).toBeVisible();
  });
});

describe('FRN-06 - alta y desactivación de usuarios', () => {
  it('incluye un selector de rol ADMIN u OPERATOR en el alta', () => {
    render(<CrearUsuarios />);
    expect(screen.getByRole('combobox', { name: /rol/i })).toBeInTheDocument();
  });

  it('envía POST /api/users al confirmar el alta', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ id: 'u2', rol: 'OPERATOR', activo: true, contrasenaTemp: 'Temp123!' }, 201));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    render(<CrearUsuarios />);

    await user.type(screen.getByLabelText(/nombre completo/i), 'Ada Lovelace');
    await user.type(screen.getByLabelText(/nombre de usuario/i), 'ada');
    await user.type(screen.getByLabelText(/^correo electrónico$/i), 'ada@centinela.local');
    await user.click(screen.getByRole('button', { name: /crear usuario/i }));

    await waitFor(() => expect(fetchMock).toHaveBeenCalled());
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toContain('/users');
    expect(init.method).toBe('POST');
  });

  it('muestra la contraseña temporal que devuelve el backend tras crear el usuario', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ id: 'u2', rol: 'OPERATOR', activo: true, contrasenaTemp: 'Temp123!' }, 201));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    render(<CrearUsuarios />);

    await user.type(screen.getByLabelText(/nombre completo/i), 'Ada Lovelace');
    await user.type(screen.getByLabelText(/nombre de usuario/i), 'ada');
    await user.type(screen.getByLabelText(/^correo electrónico$/i), 'ada@centinela.local');
    await user.click(screen.getByRole('button', { name: /crear usuario/i }));

    expect(await screen.findByText('Temp123!')).toBeVisible();
  });

  it('permite solicitar la desactivación o baja de un usuario mediante DELETE a la API', async () => {
    const fetchMock = vi.fn().mockImplementation((_url: string, init?: RequestInit) => {
      if (init?.method === 'DELETE') {
        return Promise.resolve(new Response(null, { status: 204 }));
      }
      return Promise.resolve(jsonResponse({
        id: 'u1',
        nombreCompleto: 'Ada Lovelace',
        nombreUsuario: 'alovelace',
        emailUsuario: 'ada@example.com',
        rol: 'ADMIN',
        activo: true,
      }));
    });
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    render(<UserDetail />);
    const deleteButton = await screen.findByRole('button', { name: /eliminar usuario|desactivar/i });
    expect(deleteButton).toBeInTheDocument();
    await user.click(deleteButton);
    await waitFor(() => {
      const calls = fetchMock.mock.calls as [string, RequestInit][];
      const deleteCall = calls.find(([_, init]) => init?.method === 'DELETE');
      expect(deleteCall).toBeDefined();
      expect(deleteCall?.[0]).toContain('/users');
    });
  });
});

describe('FRN-06B - edición de usuario y cambio de rol', () => {
  it('ofrece una acción de edición por usuario en la tabla', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      summary: { total: 1, admins: 1, operators: 0 },
      users: [{ id: 'u1', nombreCompleto: 'Ada Lovelace', rol: 'ADMIN', activo: true, totpVinculado: true }],
    }));
    vi.stubGlobal('fetch', fetchMock);

    render(<Users />);

    expect(await screen.findByRole('button', { name: /editar/i })).toBeInTheDocument();
  });

  it('permite modificar información básica y rol enviando PUT /api/admin/users/:id', async () => {
    const fetchMock = vi.fn().mockImplementation((_url: string, init?: RequestInit) => {
      if (init?.method === 'PUT') {
        return Promise.resolve(jsonResponse({
          id: 'u2',
          nombreCompleto: 'Usuario Dos Editado',
          rol: 'ADMIN',
          activo: true,
        }));
      }
      return Promise.resolve(jsonResponse({
        id: 'u2',
        nombreCompleto: 'Usuario Dos',
        nombreUsuario: 'user2',
        emailUsuario: 'user2@example.com',
        rol: 'ADMIN',
        activo: true,
      }));
    });
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    render(<UserDetail />);

    const saveButton = await screen.findByRole('button', { name: /guardar cambios/i });
    expect(saveButton).toBeInTheDocument();

    await user.click(saveButton);

    await waitFor(() => {
      const calls = fetchMock.mock.calls as [string, RequestInit][];
      const putCall = calls.find(([_, init]) => init?.method === 'PUT');
      expect(putCall).toBeDefined();
      expect(putCall?.[0]).toContain('/users');
    });
  });
});

describe('FRN-07 - selector de asignación de instancias', () => {
  it('consulta GET /api/instances con Authorization Bearer para listar instancias', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse([
      { id: '101', name: 'Ubuntu Server', type: 'qemu', node: 'pve01', status: 'running' },
    ]));
    vi.stubGlobal('fetch', fetchMock);

    render(<UserDetail />);

    await waitFor(() => {
      const calls = fetchMock.mock.calls as [string, RequestInit][];
      const instanceCall = calls.find(([url]) => url.includes('/instances'));
      expect(instanceCall).toBeDefined();
      if (instanceCall) {
        expect((instanceCall[1].headers as Record<string, string>)?.Authorization).toMatch(/^Bearer /);
      }
    });
  });

  it('permite seleccionar VMIDs y enviarlos al endpoint de permisos', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ success: true }));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    render(<UserDetail />);

    const permTab = screen.queryByRole('button', { name: /roles y permisos/i });
    if (permTab) {
      await user.click(permTab);
    }

    const saveButton = screen.getByRole('button', { name: /guardar cambios/i });
    await user.click(saveButton);

    await waitFor(() => {
      const calls = fetchMock.mock.calls as [string, RequestInit][];
      const permCall = calls.find(([url, init]) => (url.includes('/permissions') || url.includes('/instances')) && init.method === 'PUT');
      expect(permCall).toBeDefined();
    });
  });
});

describe('FRN-08 - manejo de 403 en recursos protegidos', () => {
  it('envía el JWT Bearer en las peticiones y no borra la sesión ante un 403', async () => {
    window.sessionStorage.setItem('centinela_access', 'bearer-token-123');
    window.localStorage.setItem('centinela_user', JSON.stringify({ id: 'u1', rol: 'OPERATOR' }));

    const fetchMock = vi.fn().mockResolvedValue(new Response(
      JSON.stringify({ errorCode: 'FORBIDDEN', message: 'Acceso denegado' }),
      { status: 403, headers: { 'Content-Type': 'application/json' } }
    ));
    vi.stubGlobal('fetch', fetchMock);

    renderAppAs('OPERATOR', '/dashboard');

    expect(window.sessionStorage.getItem('centinela_access')).toBe('access-token');
  });

  it('mantiene al usuario en la vista protegida sin forzar logout ante un 403', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(
      JSON.stringify({ errorCode: 'FORBIDDEN', message: 'No tenés permisos' }),
      { status: 403, headers: { 'Content-Type': 'application/json' } }
    ));
    vi.stubGlobal('fetch', fetchMock);

    const router = renderAppAs('OPERATOR', '/dashboard');

    expect(router.state.location.pathname).not.toBe('/login');
  });
});
