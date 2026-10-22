import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import UserDetail from '@/pages/detailsUserPage';
import { queryDialog } from './instances-helpers';
import { withAppProviders } from './app-providers';

// FIX-111B (depende de FIX-111A): cambiar el correo exige reautenticar con contraseña actual y
// código TOTP en un modal de confirmación. Sin ese paso el backend rechaza el PUT con 400/401.
// Criterio de éxito:
//   - La interfaz no permite enviar el cambio de correo sin pedir contraseña actual y código 2FA.
//   - Se informa al usuario que se envió un correo de confirmación a la nueva dirección.
// La forma del modal es libre (role="dialog" o "alertdialog"); no se fija el endpoint, solo que
// ninguna petición de escritura salga antes de confirmar y que la que salga lleve las credenciales.

const userDetails = {
  id: 'usr-10', organizacionId: 'org-1', nombreCompleto: 'Usuario de Prueba', nombreUsuario: 'usuario_prueba',
  emailUsuario: 'usuario@centinela.local', rol: 'OPERATOR', activo: true, totpVinculado: true,
  cambioContrasenaRequerido: false, fechaCreacion: '2026-10-01T12:00:00Z', fechaUltimoAcceso: '2026-10-05T15:30:00Z',
  eliminadoEn: null, instanciasPermitidas: [],
};

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

function renderUserDetail() {
  window.sessionStorage.setItem('centinela_access', 'access-token-fix111b');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
    rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
  }));
  const fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
    const method = (init.method ?? 'GET').toUpperCase();
    if (method !== 'GET') return Promise.resolve(json({ ...userDetails, emailUsuario: 'nuevo@centinela.local' }, 200));
    if (String(url).includes('/permissions')) return Promise.resolve(json({ items: [] }));
    if (String(url).includes('/instances')) return Promise.resolve(json([]));
    if (String(url).includes('/admin/audit')) return Promise.resolve(json({ total: 0, pagina: 1, items: [] }));
    return Promise.resolve(json(userDetails));
  });
  vi.stubGlobal('fetch', fetchMock);
  render(withAppProviders(
    <MemoryRouter initialEntries={[{ pathname: '/users/usr-10', state: { isCurrentUser: true } }]}>
      <Routes><Route path="/users/:userId" element={<UserDetail />} /></Routes>
    </MemoryRouter>,
  ));
  return fetchMock;
}

const writes = (fetchMock: ReturnType<typeof vi.fn>) =>
  (fetchMock.mock.calls as [string, RequestInit?][]).filter(([, init]) => ['PUT', 'PATCH', 'POST', 'DELETE'].includes((init?.method ?? 'GET').toUpperCase()));

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FIX-111B - re-autenticación (contraseña + 2FA) al cambiar el correo', () => {
  it('al guardar un correo nuevo abre un modal que pide contraseña actual y código de 6 dígitos, sin enviar nada antes', async () => {
    const fetchMock = renderUserDetail();
    const user = userEvent.setup();
    const email = await screen.findByLabelText(/correo electr[oó]nico/i);
    await screen.findByDisplayValue('Usuario de Prueba');

    await user.clear(email);
    await user.type(email, 'nuevo@centinela.local');
    await user.click(screen.getByRole('button', { name: /guardar cambios/i }));

    const dialog = await waitFor(() => {
      const found = queryDialog();
      expect(found, 'no se abrió el modal de confirmación de seguridad al cambiar el correo').not.toBeNull();
      return found as HTMLElement;
    });
    expect(within(dialog).getByLabelText(/contrase[ñn]a actual|contrase[ñn]a/i)).toBeInTheDocument();
    expect(within(dialog).getByLabelText(/c[oó]digo|2fa|autenticador|totp/i)).toBeInTheDocument();
    expect(writes(fetchMock), 'se envió una petición de escritura antes de reautenticar').toHaveLength(0);
  });

  it('confirmar envía contrasenaActual y codigo2fa y muestra el aviso del correo de confirmación', async () => {
    const fetchMock = renderUserDetail();
    const user = userEvent.setup();
    const email = await screen.findByLabelText(/correo electr[oó]nico/i);
    await screen.findByDisplayValue('Usuario de Prueba');
    await user.clear(email);
    await user.type(email, 'nuevo@centinela.local');
    await user.click(screen.getByRole('button', { name: /guardar cambios/i }));

    const dialog = await waitFor(() => {
      const found = queryDialog();
      expect(found).not.toBeNull();
      return found as HTMLElement;
    });
    await user.type(within(dialog).getByLabelText(/contrase[ñn]a actual|contrase[ñn]a/i), 'Actual123!');
    await user.type(within(dialog).getByLabelText(/c[oó]digo|2fa|autenticador|totp/i), '123456');
    await user.click(within(dialog).getByRole('button', { name: /confirmar|continuar|verificar|guardar|enviar/i }));

    await waitFor(() => expect(writes(fetchMock).length).toBeGreaterThan(0));
    const body = JSON.parse(String(writes(fetchMock)[0][1]?.body ?? '{}'));
    expect(body).toMatchObject({ contrasenaActual: 'Actual123!', codigo2fa: '123456' });
    expect(body.emailUsuario).toBe('nuevo@centinela.local');
    expect(await screen.findByText(/correo de confirmaci[oó]n|enviamos un correo|revis[aá] tu correo|nueva direcci[oó]n/i)).toBeVisible();
  });

  it('cancelar el modal no envía ninguna petición ni cambia el correo guardado', async () => {
    const fetchMock = renderUserDetail();
    const user = userEvent.setup();
    const email = await screen.findByLabelText(/correo electr[oó]nico/i);
    await screen.findByDisplayValue('Usuario de Prueba');
    await user.clear(email);
    await user.type(email, 'nuevo@centinela.local');
    await user.click(screen.getByRole('button', { name: /guardar cambios/i }));
    const dialog = await waitFor(() => {
      const found = queryDialog();
      expect(found).not.toBeNull();
      return found as HTMLElement;
    });
    await user.click(within(dialog).getByRole('button', { name: /cancelar|volver|cerrar/i }));
    await waitFor(() => expect(queryDialog()).toBeNull());
    expect(writes(fetchMock)).toHaveLength(0);
  });

  it('editar solo el nombre completo no exige reautenticación', async () => {
    const fetchMock = renderUserDetail();
    const user = userEvent.setup();
    const name = await screen.findByDisplayValue('Usuario de Prueba');
    await user.clear(name);
    await user.type(name, 'Otro Nombre');
    await user.click(screen.getByRole('button', { name: /guardar cambios/i }));

    await waitFor(() => expect(writes(fetchMock).length).toBeGreaterThan(0));
    expect(queryDialog(), 'se pidió reautenticación para un cambio que no toca el correo').toBeNull();
    const body = JSON.parse(String(writes(fetchMock)[0][1]?.body ?? '{}'));
    expect(body).not.toHaveProperty('codigo2fa');
  });
});

