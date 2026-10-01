import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';
import { Toaster } from '@/components/ui/toast';
import { ApiResponseNotifier } from '@/components/common/ApiResponseNotifier';
import { apiClient } from '@/services/apiClient';
import Users from '@/pages/Users';
import CrearUsuarios from '@/pages/CrearUsuarios';
import UserDetail from '@/pages/detailsUserPage';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

type FetchCall = [string, RequestInit];
const callsOf = (fetchMock: ReturnType<typeof vi.fn>) => fetchMock.mock.calls as FetchCall[];
const bodyOf = (call: FetchCall) => JSON.parse(String(call[1].body ?? '{}'));

// Contrato real de GET /api/admin/users/:id (UsuarioDetalleDTO).
function defaultUserDetails(userId = 'u2') {
  return {
    id: userId,
    nombreCompleto: 'Ada Lovelace',
    nombreUsuario: 'alovelace',
    emailUsuario: 'ada@example.com',
    organizacionId: 'org-1',
    rol: 'OPERATOR',
    activo: true,
    totpVinculado: true,
    cambioContrasenaRequerido: false,
    fechaCreacion: '2026-09-01T00:00:00Z',
    fechaUltimoAcceso: null,
    instanciasPermitidas: [101],
  };
}

// Contrato real de GET /api/instances (BAC-14): { id, name, type: vm|lxc, node, status }.
const inventory = [
  { id: 101, name: 'Ubuntu Server', type: 'vm', node: 'pve', status: 'running' },
  { id: 102, name: 'Debian 12', type: 'lxc', node: 'pve', status: 'stopped' },
];

// Backend simulado para la vista de detalle. Las mutaciones responden como el backend real.
// Permisos (SEC-04): GET/PUT /admin/users/:id/permissions -> { permisos: [{ vmid, nivelAcceso }] }.
function detailBackend({ profilePut, permisos = [{ vmid: 101, nivelAcceso: 'FULL_ACCESS' }] }: {
  profilePut?: (init: RequestInit) => Response;
  permisos?: { vmid: number; nivelAcceso: string }[];
} = {}) {
  return vi.fn().mockImplementation((url: string, init?: RequestInit) => {
    const method = init?.method ?? 'GET';
    if (url.includes('/permissions')) {
      return Promise.resolve(method === 'PUT'
        ? new Response(null, { status: 204 })
        : jsonResponse({ permisos }));
    }
    if (url.includes('/instances')) return Promise.resolve(jsonResponse(inventory));
    if (method === 'DELETE') return Promise.resolve(new Response(null, { status: 204 }));
    if (method === 'PUT') {
      return Promise.resolve(profilePut ? profilePut(init!) : jsonResponse({ ...defaultUserDetails('u2'), ...bodyOf([url, init!]) }));
    }
    return Promise.resolve(jsonResponse(defaultUserDetails('u2')));
  });
}

function renderUserDetail(userId = 'u2') {
  window.sessionStorage.setItem('centinela_access', 'access-token-admin');
  return render(
    <Toaster>
      <MemoryRouter initialEntries={[`/users/${userId}`]}>
        <Routes>
          <Route path="/users/:userId" element={<UserDetail />} />
        </Routes>
      </MemoryRouter>
    </Toaster>,
  );
}

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
  render(<Toaster><RouterProvider router={router} /></Toaster>);
  return router;
}

async function fillCreateUserForm(user: ReturnType<typeof userEvent.setup>) {
  await user.type(screen.getByLabelText(/nombre completo/i), 'Ada Lovelace');
  await user.type(screen.getByLabelText(/nombre de usuario/i), 'ada');
  await user.type(screen.getByLabelText(/^correo electrónico/i), 'ada@centinela.local');
  const confirmation = screen.queryByLabelText(/confirmar correo/i);
  if (confirmation) await user.type(confirmation, 'ada@centinela.local');
  await user.click(screen.getByRole('button', { name: /crear usuario/i }));
}

// Si la acción abre un diálogo de confirmación, lo confirma; si no, no hace nada.
async function confirmIfAsked(user: ReturnType<typeof userEvent.setup>) {
  const dialog = await screen.findByRole('alertdialog', {}, { timeout: 500 }).catch(() => screen.queryByRole('dialog'));
  if (!dialog) return;
  const buttons = within(dialog).getAllByRole('button');
  const confirm = buttons.find((button) => !/cancelar|cerrar|volver/i.test(button.textContent ?? '')) ?? buttons[buttons.length - 1];
  await user.click(confirm);
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-05 - panel de gestión de usuarios', () => {
  it('un usuario con rol OPERATOR no puede ver el panel de administración', async () => {
    const router = renderAppAs('OPERATOR', '/users');

    await waitFor(() => expect(router.state.location.pathname).toBe('/dashboard'));
    expect(screen.queryByRole('heading', { name: 'Gestión de usuarios' })).not.toBeInTheDocument();
  });

  it('consulta GET /api/admin/users con Authorization Bearer al entrar a /users', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ summary: { total: 0, admins: 0, operators: 0 }, users: [] }));
    vi.stubGlobal('fetch', fetchMock);

    renderAppAs('ADMIN', '/users');

    await waitFor(() => expect(callsOf(fetchMock).some(([url]) => url.includes('/admin/users'))).toBe(true));
    const [, init] = callsOf(fetchMock).find(([url]) => url.includes('/admin/users'))!;
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer access-token');
  });

  it('renderiza los usuarios reales devueltos por el backend', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({
      summary: { total: 1, admins: 1, operators: 0 },
      users: [{ id: 'u1', nombreCompleto: 'Ada Lovelace', rol: 'ADMIN', activo: true, totpVinculado: true }],
    })));

    renderAppAs('ADMIN', '/users');

    const table = await screen.findByRole('table');
    expect(await within(table).findByText('Ada Lovelace')).toBeVisible();
  });
});

describe('FRN-06 - alta y desactivación de usuarios', () => {
  it('incluye un selector de rol ADMIN u OPERATOR en el alta', () => {
    render(<MemoryRouter><CrearUsuarios /></MemoryRouter>);
    const roleSelect = screen.getByRole('combobox', { name: /rol/i });
    const options = within(roleSelect).getAllByRole('option').map((option) => (option as HTMLOptionElement).value);
    expect(options).toEqual(expect.arrayContaining(['ADMIN', 'OPERATOR']));
  });

  it('envía POST /api/admin/users con el contrato del backend al confirmar el alta', { timeout: 15000 }, async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ id: 'u2', rol: 'OPERATOR', activo: true }, 201));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    render(<Toaster><MemoryRouter><CrearUsuarios /></MemoryRouter></Toaster>);

    await fillCreateUserForm(user);
    await confirmIfAsked(user);

    await waitFor(() => expect(callsOf(fetchMock).some(([url, init]) => url.includes('/admin/users') && init.method === 'POST')).toBe(true));
    const call = callsOf(fetchMock).find(([url, init]) => url.includes('/admin/users') && init.method === 'POST')!;
    expect(bodyOf(call)).toMatchObject({ nombreCompleto: 'Ada Lovelace', nombreUsuario: 'ada', emailUsuario: 'ada@centinela.local', rol: expect.stringMatching(/^(ADMIN|OPERATOR)$/) });
  });

  it('tras el alta informa el resultado y que la clave temporal se envió por correo (BAC-16), sin mostrar ninguna contraseña', { timeout: 15000 }, async () => {
    // BAC-16 (terminado.md): el backend ya NO devuelve contrasenaTemp; la entrega por correo.
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({ id: 'u2', rol: 'OPERATOR', activo: true }, 201)));
    const user = userEvent.setup();
    render(<Toaster><MemoryRouter><CrearUsuarios /></MemoryRouter></Toaster>);

    await fillCreateUserForm(user);
    await confirmIfAsked(user);

    expect((await screen.findAllByText(/usuario creado|creado correctamente|se envi[óo].*correo|correo.*enviad/i, {}, { timeout: 2000 })).length).toBeGreaterThan(0);
    expect(screen.queryByText(/contraseña temporal generada/i)).not.toBeInTheDocument();
  });

  it('el formulario de alta ya no pide contraseña: la genera el backend y viaja por correo (FIX-24)', () => {
    render(<MemoryRouter><CrearUsuarios /></MemoryRouter>);
    expect(screen.queryByLabelText(/contraseña/i)).not.toBeInTheDocument();
  });

  it('FIX-27 si el correo no se pudo enviar (502 EMAIL_DELIVERY_FAILED) informa que el usuario no fue creado', { timeout: 15000 }, async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse(
      { errorCode: 'EMAIL_DELIVERY_FAILED', message: 'Usuario no creado. El servidor de correo no está disponible.' }, 502,
    )));
    const user = userEvent.setup();
    render(<Toaster><MemoryRouter><CrearUsuarios /></MemoryRouter></Toaster>);

    await fillCreateUserForm(user);
    await confirmIfAsked(user);

    expect((await screen.findAllByText(/no se pudo enviar el correo|servidor de correo no est[áa] disponible/i, {}, { timeout: 2000 })).length).toBeGreaterThan(0);
  });

  it('FIX-27 ante un 502 del correo permanece en el alta con los datos cargados para reintentar', { timeout: 15000 }, async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse(
      { errorCode: 'EMAIL_DELIVERY_FAILED', message: 'Usuario no creado. El servidor de correo no está disponible.' }, 502,
    )));
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token-admin');
    const router = createMemoryRouter([
      { path: '/users/new', Component: CrearUsuarios },
      { path: '/users', element: <h1>Gestión de usuarios</h1> },
    ], { initialEntries: ['/users/new'] });
    render(<Toaster><RouterProvider router={router} /></Toaster>);

    await fillCreateUserForm(user);
    await confirmIfAsked(user);

    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    expect(router.state.location.pathname).toBe('/users/new');
    expect(screen.getByLabelText(/nombre completo/i)).toHaveValue('Ada Lovelace');
    expect(screen.getByLabelText(/^correo electrónico/i)).toHaveValue('ada@centinela.local');
  });

  it('el botón "Eliminar usuario" solicita DELETE /api/admin/users/:id', async () => {
    const fetchMock = detailBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderUserDetail('u2');
    await user.click(await screen.findByRole('button', { name: /eliminar usuario|desactivar/i }));
    await confirmIfAsked(user);

    await waitFor(() => {
      const deleteCall = callsOf(fetchMock).find(([, init]) => init?.method === 'DELETE');
      expect(deleteCall, 'la acción no envió DELETE a la API').toBeDefined();
      expect(deleteCall![0]).toMatch(/\/admin\/users\/u2$/);
    });
  });
});

describe('FRN-06B - edición de usuario y cambio de rol', () => {
  it('ofrece una acción de edición por usuario en la tabla', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({
      summary: { total: 1, admins: 1, operators: 0 },
      users: [{ id: 'u1', nombreCompleto: 'Ada Lovelace', rol: 'ADMIN', activo: true, totpVinculado: true }],
    })));
    const user = userEvent.setup();

    render(<MemoryRouter><Users /></MemoryRouter>);

    await user.click(await screen.findByRole('button', { name: /acciones para/i }));
    expect(await screen.findByRole('button', { name: /editar usuario/i })).toBeInTheDocument();
  });

  it('editar el nombre y guardar envía PUT /api/admin/users/:id con los datos modificados', async () => {
    const fetchMock = detailBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderUserDetail('u2');
    const nameInput = await screen.findByLabelText(/nombre completo/i);
    await user.clear(nameInput);
    await user.type(nameInput, 'Ada Byron');
    await user.click(screen.getByRole('button', { name: /guardar cambios/i }));

    await waitFor(() => {
      const putCall = callsOf(fetchMock).find(([url, init]) => init?.method === 'PUT' && /\/admin\/users\/u2$/.test(url));
      expect(putCall, 'guardar no envió PUT /admin/users/u2').toBeDefined();
      expect(bodyOf(putCall!)).toMatchObject({ nombreCompleto: 'Ada Byron' });
    });
  });
});

describe('FIX-38 - "Solo lectura" no es un rol de usuario (BAC-09 / SEC-04)', () => {
  it('el selector "Rol" de la ficha ofrece solo ADMIN y OPERATOR (el backend rechaza otro rol con 400)', async () => {
    vi.stubGlobal('fetch', detailBackend());

    renderUserDetail('u2');
    const roleSelect = await screen.findByRole('combobox', { name: /^rol$/i });
    const options = within(roleSelect).getAllByRole('option').map((option) => (option as HTMLOptionElement).value);
    expect(options.sort()).toEqual(['ADMIN', 'OPERATOR']);
  });
});

describe('FIX-25 - errores al guardar el perfil', () => {
  it('un correo duplicado (409 USER_CONFLICT) se muestra junto al campo de correo', async () => {
    vi.stubGlobal('fetch', detailBackend({
      profilePut: () => jsonResponse({ errorCode: 'USER_CONFLICT', message: 'el email ya está en uso' }, 409),
    }));
    const user = userEvent.setup();

    renderUserDetail('u2');
    const emailInput = await screen.findByLabelText(/correo electrónico/i);
    await user.clear(emailInput);
    await user.type(emailInput, 'otro@example.com');
    await user.click(screen.getByRole('button', { name: /guardar cambios/i }));

    await waitFor(() => expect(emailInput).toHaveAttribute('aria-invalid', 'true'));
    expect(screen.getByText(/correo.*(ya|pertenece|uso|existe)/i)).toBeVisible();
  });
});

describe('FIX-14 / FRN-07 - selector de asignación de instancias', () => {
  it('consulta GET /api/instances con Authorization Bearer', async () => {
    const fetchMock = detailBackend();
    vi.stubGlobal('fetch', fetchMock);

    renderUserDetail('u2');

    await waitFor(() => {
      const instanceCall = callsOf(fetchMock).find(([url]) => /\/instances$/.test(url));
      expect(instanceCall).toBeDefined();
      expect((instanceCall![1].headers as Record<string, string>).Authorization).toBe('Bearer access-token-admin');
    });
  });

  it('muestra marcadas las instancias ya asignadas (GET /permissions) y sin marcar las demás', async () => {
    vi.stubGlobal('fetch', detailBackend());
    const user = userEvent.setup();

    renderUserDetail('u2');
    await user.click(await screen.findByRole('tab', { name: /roles y permisos/i }));

    expect(await screen.findByRole('checkbox', { name: /ubuntu server/i })).toBeChecked();
    expect(screen.getByRole('checkbox', { name: /debian 12/i })).not.toBeChecked();
  });

  it('seleccionar una instancia y guardar envía PUT /api/admin/users/:id/permissions con { permisos }', async () => {
    const fetchMock = detailBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderUserDetail('u2');
    await user.click(await screen.findByRole('tab', { name: /roles y permisos/i }));
    await user.click(await screen.findByRole('checkbox', { name: /debian 12/i }));
    await user.click(screen.getByRole('button', { name: /guardar cambios/i }));

    await waitFor(() => {
      const permCall = callsOf(fetchMock).find(([url, init]) => init?.method === 'PUT' && url.endsWith('/admin/users/u2/permissions'));
      expect(permCall, 'no se envió PUT /admin/users/u2/permissions').toBeDefined();
      const permisos = bodyOf(permCall!).permisos as { vmid: number; nivelAcceso?: string }[];
      expect(permisos.map((permiso) => permiso.vmid).sort()).toEqual([101, 102]);
      expect(permisos.every((permiso) => (permiso.nivelAcceso ?? 'FULL_ACCESS') === 'FULL_ACCESS')).toBe(true);
    });
  });

  it('FRN-18 al abrir la ficha muestra el nivel guardado: READ_ONLY como "Solo lectura" y FULL_ACCESS como acceso completo', async () => {
    vi.stubGlobal('fetch', detailBackend({ permisos: [
      { vmid: 101, nivelAcceso: 'FULL_ACCESS' },
      { vmid: 102, nivelAcceso: 'READ_ONLY' },
    ] }));
    const user = userEvent.setup();

    renderUserDetail('u2');
    await user.click(await screen.findByRole('tab', { name: /roles y permisos/i }));

    expect(await screen.findByRole('combobox', { name: /acceso para debian 12/i })).toHaveValue('Solo lectura');
    expect(screen.getByRole('combobox', { name: /acceso para ubuntu server/i })).toHaveValue('Acceso completo');
  });

  it('FRN-18 elegir "Solo lectura" envía nivelAcceso READ_ONLY para esa instancia (integración con SEC-04)', async () => {
    const fetchMock = detailBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderUserDetail('u2');
    await user.click(await screen.findByRole('tab', { name: /roles y permisos/i }));
    await user.selectOptions(await screen.findByRole('combobox', { name: /acceso para debian 12/i }), 'Solo lectura');
    await user.click(screen.getByRole('button', { name: /guardar cambios/i }));

    await waitFor(() => {
      const permCall = callsOf(fetchMock).find(([url, init]) => init?.method === 'PUT' && url.endsWith('/admin/users/u2/permissions'));
      expect(permCall, 'no se envió PUT /admin/users/u2/permissions').toBeDefined();
      const permisos = bodyOf(permCall!).permisos as { vmid: number; nivelAcceso?: string }[];
      expect(permisos).toEqual(expect.arrayContaining([
        expect.objectContaining({ vmid: 102, nivelAcceso: 'READ_ONLY' }),
        expect.objectContaining({ vmid: 101 }),
      ]));
    });
  });
});

describe('FRN-08 - manejo de 403 en recursos protegidos', () => {
  it('el interceptor envía el Bearer y ante un 403 muestra un aviso sin cerrar la sesión', async () => {
    window.sessionStorage.setItem('centinela_access', 'bearer-token-123');
    window.localStorage.setItem('centinela_user', JSON.stringify({ id: 'u1', rol: 'OPERATOR' }));
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ errorCode: 'INSTANCE_ACCESS_DENIED', message: 'No tenés acceso a esta instancia.' }, 403));
    vi.stubGlobal('fetch', fetchMock);
    render(<Toaster><ApiResponseNotifier /></Toaster>);

    await expect(apiClient.get('/instances/103')).rejects.toMatchObject({ status: 403 });

    expect((callsOf(fetchMock)[0][1].headers as Record<string, string>).Authorization).toBe('Bearer bearer-token-123');
    expect((await screen.findAllByText('No tenés acceso a esta instancia.')).length).toBeGreaterThan(0);
    expect(window.sessionStorage.getItem('centinela_access')).toBe('bearer-token-123');
    expect(window.localStorage.getItem('centinela_user')).not.toBeNull();
  });
});
