import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';
import Users from '@/pages/Users.jsx';
import CrearUsuarios from '@/pages/CrearUsuarios';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

// Igual que en navigation.test.tsx: siembra una sesión válida antes de navegar.
function renderAppAs(rol: 'ADMIN' | 'OPERATOR', path: string) {
  window.localStorage.setItem('centinela_access', 'access-token');
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

// Estos tests corren contra el código real (no son it.todo): documentan qué falta
// para dar por cumplidos los criterios de éxito de actual.md, no solo si compila.
describe('FRN-05 - panel de gestión de usuarios', () => {
  it('un usuario con rol OPERATOR no debe poder ver el panel de administración', async () => {
    renderAppAs('OPERATOR', '/users');
    await screen.findByText(/./);
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

    expect(await screen.findByText('Ada Lovelace')).toBeVisible();
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
});

// FRN-07 y FRN-08 siguen en todo real: no existe ningún componente, ruta ni
// servicio en el frontend contra el cual afirmar o refutar el criterio.
describe('FRN-07 - selector de asignación de instancias', () => {
  it.todo('consulta GET /api/instances con Authorization Bearer');
  it.todo('muestra instancias disponibles y permisos actuales del usuario');
  it.todo('envía el array de VMIDs seleccionado al endpoint de permisos');
  it.todo('recarga la selección guardada al abrir nuevamente el usuario');
});

describe('FRN-08 - manejo de 403 en recursos protegidos', () => {
  it.todo('envía el JWT Bearer en las peticiones de instancias y permisos');
  it.todo('muestra un error amigable ante 403 sin borrar la sesión local');
  it.todo('mantiene al usuario en la vista después de un 403');
});
