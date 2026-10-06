import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { Toaster } from '@/components/ui/toast';
import Users from '@/pages/Users';
import UserDetail from '@/pages/detailsUserPage';
import CrearUsuarios from '@/pages/CrearUsuarios';
import { AuthProvider } from '@/context/AuthContext';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

type FetchCall = [string, RequestInit];
const callsOf = (fetchMock: ReturnType<typeof vi.fn>) => fetchMock.mock.calls as FetchCall[];

// Usuario activo estándar.
function activeUser(id = 'u1') {
  return { id, nombreCompleto: 'Ada Lovelace', rol: 'OPERATOR', activo: true, eliminadoEn: null, totpVinculado: true };
}

// Usuario suspendido (activo=false, sin fecha de eliminación).
function suspendedUser(id = 'u2') {
  return { id, nombreCompleto: 'Charles Babbage', rol: 'OPERATOR', activo: false, eliminadoEn: null, totpVinculado: false };
}

// Usuario eliminado (eliminadoEn con timestamp).
function deletedUser(id = 'u3') {
  return { id, nombreCompleto: 'Grace Hopper', rol: 'ADMIN', activo: false, eliminadoEn: '2026-09-15T12:00:00Z', totpVinculado: true };
}

// Detalle completo de un usuario eliminado (UsuarioDetalleDTO con campos de FIX-54).
function deletedUserDetails(userId = 'u3') {
  return {
    id: userId,
    nombreCompleto: 'Grace Hopper',
    nombreUsuario: 'ghopper',
    emailUsuario: 'grace@example.com',
    organizacionId: 'org-1',
    rol: 'ADMIN',
    activo: false,
    eliminadoEn: '2026-09-15T12:00:00Z',
    totpVinculado: true,
    cambioContrasenaRequerido: false,
    fechaCreacion: '2026-01-01T00:00:00Z',
    fechaUltimoAcceso: '2026-09-14T08:00:00Z',
    instanciasPermitidas: [],
  };
}

function setUpAdminSession() {
  window.sessionStorage.setItem('centinela_access', 'access-token-admin');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin',
    organizacionId: 'org-1',
    nombreCompleto: 'Admin',
    email: 'admin@centinela.local',
    rol: 'ADMIN',
    instanciasPermitidas: [],
    tiene2FA: true,
  }));
}

function renderUsersList(fetchMock: ReturnType<typeof vi.fn>) {
  vi.stubGlobal('fetch', fetchMock);
  setUpAdminSession();
  return render(
    <AuthProvider>
      <Toaster>
        <MemoryRouter initialEntries={['/users']}>
          <Routes>
            <Route path="/users" element={<Users />} />
          </Routes>
        </MemoryRouter>
      </Toaster>
    </AuthProvider>,
  );
}

function renderUserDetail(userId = 'u3') {
  setUpAdminSession();
  return render(
    <AuthProvider>
      <Toaster>
        <MemoryRouter initialEntries={[`/users/${userId}`]}>
          <Routes>
            <Route path="/users/:userId" element={<UserDetail />} />
          </Routes>
        </MemoryRouter>
      </Toaster>
    </AuthProvider>,
  );
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FIX-56 - distinguir suspensión y eliminación en la interfaz', () => {
  it('muestra etiquetas diferenciadas para usuarios activos, suspendidos y eliminados en la tabla', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      summary: { total: 3, admins: 1, operators: 2 },
      users: [activeUser(), suspendedUser(), deletedUser()],
    }));

    renderUsersList(fetchMock);
    const table = await screen.findByRole('table');

    // Usuario activo → etiqueta "Activo"
    const activeRow = (await within(table).findByText('Ada Lovelace')).closest('tr')!;
    expect(within(activeRow).getByText('Activo')).toBeVisible();

    // Usuario suspendido → etiqueta "Inactivo" o "Suspendido" (distinta de "Eliminado")
    const suspendedRow = (await within(table).findByText('Charles Babbage')).closest('tr')!;
    const suspendedBadge = within(suspendedRow).getByText(/inactivo|suspendido/i);
    expect(suspendedBadge).toBeVisible();
    expect(suspendedBadge.textContent).not.toMatch(/eliminado/i);

    // Usuario eliminado → etiqueta "Eliminado" (diferente a la de suspendido)
    const deletedRow = (await within(table).findByText('Grace Hopper')).closest('tr')!;
    expect(within(deletedRow).getByText(/eliminado/i)).toBeVisible();
  });

  it('oculta usuarios eliminados en la vista por defecto de la lista', async () => {
    // La lista debería filtrar los eliminados por defecto, o el backend no los envía.
    // Si la UI recibe eliminados, debería ocultarlos a menos que se active un filtro.
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      summary: { total: 3, admins: 1, operators: 2 },
      users: [activeUser(), suspendedUser(), deletedUser()],
    }));

    renderUsersList(fetchMock);
    const table = await screen.findByRole('table');

    // Los usuarios activos y suspendidos deberían ser visibles.
    expect(await within(table).findByText('Ada Lovelace')).toBeVisible();
    expect(within(table).getByText('Charles Babbage')).toBeVisible();

    // El usuario eliminado NO debería aparecer en la vista por defecto.
    expect(within(table).queryByText('Grace Hopper')).not.toBeInTheDocument();
  });

  it('deshabilita o esconde acciones como Editar, Reactivar y Resetear contraseña para un usuario eliminado', async () => {
    const fetchMock = vi.fn().mockImplementation((url: string, init?: RequestInit) => {
      const method = init?.method ?? 'GET';
      if (url.includes('/permissions')) return Promise.resolve(jsonResponse({ permisos: [] }));
      if (url.includes('/instances')) return Promise.resolve(jsonResponse([]));
      if (method === 'GET') return Promise.resolve(jsonResponse(deletedUserDetails('u3')));
      return Promise.resolve(new Response(null, { status: 204 }));
    });
    vi.stubGlobal('fetch', fetchMock);

    renderUserDetail('u3');

    // Esperar a que se cargue la ficha del usuario eliminado.
    await screen.findByText('Grace Hopper');

    // Los botones de acciones destructivas o de edición deben estar deshabilitados o ausentes.
    const editButton = screen.queryByRole('button', { name: /editar usuario/i });
    const reactivateButton = screen.queryByRole('button', { name: /reactivar/i });
    const resetPasswordButton = screen.queryByRole('button', { name: /restablecer contraseña|resetear contraseña/i });

    // Cada botón debe estar ausente o, si presente, deshabilitado.
    if (editButton) expect(editButton).toBeDisabled();
    if (reactivateButton) expect(reactivateButton).toBeDisabled();
    if (resetPasswordButton) expect(resetPasswordButton).toBeDisabled();

    // Al menos uno de los botones principales (guardar cambios) debería estar deshabilitado.
    const saveButton = screen.queryByRole('button', { name: /guardar cambios/i });
    if (saveButton) expect(saveButton).toBeDisabled();
  });

  it('el modal de eliminación menciona que el historial se conserva y el correo queda libre', async () => {
    const fetchMock = vi.fn().mockImplementation((url: string, init?: RequestInit) => {
      const method = init?.method ?? 'GET';
      if (url.includes('/permissions')) return Promise.resolve(jsonResponse({ permisos: [] }));
      if (url.includes('/instances')) return Promise.resolve(jsonResponse([]));
      if (method === 'DELETE') return Promise.resolve(new Response(null, { status: 204 }));
      return Promise.resolve(jsonResponse({
        id: 'u4',
        nombreCompleto: 'Alan Turing',
        nombreUsuario: 'aturing',
        emailUsuario: 'alan@example.com',
        organizacionId: 'org-1',
        rol: 'OPERATOR',
        activo: true,
        eliminadoEn: null,
        totpVinculado: true,
        cambioContrasenaRequerido: false,
        fechaCreacion: '2026-01-01T00:00:00Z',
        fechaUltimoAcceso: null,
        instanciasPermitidas: [],
      }));
    });
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderUserDetail('u4');

    // Esperar a que cargue la ficha y hacer click en "Eliminar usuario".
    const deleteButton = await screen.findByRole('button', { name: /eliminar usuario/i });
    await user.click(deleteButton);

    // El diálogo de confirmación para eliminación debería mencionar historial y correo libre.
    const dialog = await screen.findByRole('alertdialog').catch(() => screen.findByRole('dialog'));
    expect(within(dialog).getByText(/historial.*conserv|conserv.*historial|registros.*mantien/i)).toBeVisible();
    expect(within(dialog).getByText(/correo.*libera|email.*libera|correo.*disponible|email.*disponible/i)).toBeVisible();
  });

  it('el modal de suspensión menciona que la acción es reversible y el correo permanece reservado', async () => {
    // Renderizar la tabla con un usuario activo para acceder a la acción de suspender.
    const fetchMock = vi.fn().mockImplementation((url: string, init?: RequestInit) => {
      const method = init?.method ?? 'GET';
      if (url.includes('/admin/users') && method === 'GET') {
        return Promise.resolve(jsonResponse({
          summary: { total: 1, admins: 0, operators: 1 },
          users: [{ ...activeUser('u1'), emailUsuario: 'ada@centinela.local', nombreUsuario: 'alovelace' }],
        }));
      }
      if (method === 'PUT') return Promise.resolve(jsonResponse({ ...activeUser('u1'), activo: false }));
      return Promise.resolve(jsonResponse({}, 200));
    });
    vi.stubGlobal('fetch', fetchMock);
    setUpAdminSession();
    const user = userEvent.setup();

    render(
      <AuthProvider>
        <Toaster>
          <MemoryRouter initialEntries={['/users']}>
            <Routes>
              <Route path="/users" element={<Users />} />
            </Routes>
          </MemoryRouter>
        </Toaster>
      </AuthProvider>,
    );

    // Abrir el menú de acciones del usuario y hacer click en "Desactivar" o "Suspender".
    const actionButton = await screen.findByRole('button', { name: /acciones para/i });
    await user.click(actionButton);
    const suspendButton = await screen.findByRole('button', { name: /desactivar usuario|suspender/i });
    await user.click(suspendButton);

    // El diálogo de confirmación para suspensión debería explicar reversibilidad y reserva de correo.
    const dialog = await screen.findByRole('alertdialog').catch(() => screen.findByRole('dialog'));
    expect(within(dialog).getByText(/reversible|reactivar|volver a activar/i)).toBeVisible();
    expect(within(dialog).getByText(/correo.*reserv|email.*reserv|correo.*no.*libera|email.*no.*libera/i)).toBeVisible();
  });

  it('crear un usuario con email de un usuario activo o suspendido (no eliminado) muestra error 409 claro', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse(
      { errorCode: 'USER_CONFLICT', message: 'El correo electrónico ya está registrado por un usuario activo o suspendido' }, 409,
    )));
    setUpAdminSession();
    const user = userEvent.setup();
    const router = createMemoryRouter([
      { path: '/users/new', Component: CrearUsuarios },
      { path: '/users', element: <h1>Gestión de usuarios</h1> },
    ], { initialEntries: ['/users/new'] });
    render(<Toaster><RouterProvider router={router} /></Toaster>);

    // Completar el formulario de alta.
    await user.type(screen.getByLabelText(/nombre completo/i), 'Nuevo Usuario');
    await user.type(screen.getByLabelText(/nombre de usuario/i), 'nuevo');
    await user.type(screen.getByLabelText(/^correo electrónico/i), 'duplicado@centinela.local');
    const confirmation = screen.queryByLabelText(/confirmar correo/i);
    if (confirmation) await user.type(confirmation, 'duplicado@centinela.local');
    await user.click(screen.getByRole('button', { name: /crear usuario/i }));

    // Si se abre un diálogo de confirmación, confirmarlo.
    const dialog = await screen.findByRole('alertdialog', {}, { timeout: 500 }).catch(() => screen.queryByRole('dialog'));
    if (dialog) {
      const buttons = within(dialog).getAllByRole('button');
      const confirm = buttons.find((button) => !/cancelar|cerrar|volver/i.test(button.textContent ?? '')) ?? buttons[buttons.length - 1];
      await user.click(confirm);
    }

    // Debería mostrar un mensaje de conflicto que indique que el correo pertenece a un usuario existente (activo/suspendido).
    await waitFor(() => {
      const notices = screen.getAllByText(/ya.*registrad|datos ya registrados|correo.*(ya|existe|pertenece|uso)/i);
      expect(notices.length).toBeGreaterThan(0);
    });
    expect(router.state.location.pathname).toBe('/users/new');
  });
});
