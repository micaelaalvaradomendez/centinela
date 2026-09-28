import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { Toaster } from '@/components/ui/toast';
import Users from '@/pages/Users';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

type FetchCall = [string, RequestInit];

// Backend simulado con los contratos reales:
//   POST /api/admin/users/:id/password/reset -> 200 { message }  (la clave viaja por correo, BAC-16)
//   POST /api/admin/users/:id/2fa/reset      -> 200 { message }
function usersBackend() {
  let totpVinculado = true;
  return vi.fn().mockImplementation((url: string, init?: RequestInit) => {
    if (init?.method === 'POST' && url.endsWith('/admin/users/u1/password/reset')) {
      return Promise.resolve(jsonResponse({ message: 'Contraseña restablecida. Se ha enviado un correo al usuario.' }));
    }
    if (init?.method === 'POST' && url.endsWith('/admin/users/u1/2fa/reset')) {
      totpVinculado = false;
      return Promise.resolve(jsonResponse({ message: 'TOTP reseteado. El usuario deberá vincularlo en su próximo acceso.' }));
    }
    return Promise.resolve(jsonResponse({
      summary: { total: 1, admins: 0, operators: 1 },
      users: [{ id: 'u1', nombreCompleto: 'Ada Lovelace', email: 'ada@centinela.local', rol: 'OPERATOR', activo: true, totpVinculado }],
    }));
  });
}

const mutationCalls = (fetchMock: ReturnType<typeof vi.fn>, fragment: string) =>
  (fetchMock.mock.calls as FetchCall[]).filter(([url, init]) => url.includes(fragment) && init?.method === 'POST');

function renderUsers() {
  window.sessionStorage.setItem('centinela_access', 'admin-token-123');
  return render(<Toaster><MemoryRouter><Users /></MemoryRouter></Toaster>);
}

async function openActions(user: ReturnType<typeof userEvent.setup>) {
  await user.click(await screen.findByRole('button', { name: /acciones para ada lovelace/i }));
}

// Acepta tanto <AlertDialog> (role=alertdialog) como <Dialog> (role=dialog).
function findConfirmationDialog() {
  return waitFor(() => screen.queryByRole('alertdialog') ?? screen.getByRole('dialog'), { timeout: 1000 });
}

async function confirmDialog(user: ReturnType<typeof userEvent.setup>) {
  const dialog = await findConfirmationDialog();
  const confirm = within(dialog).getAllByRole('button').find((button) => !/cancelar|cerrar|volver/i.test(button.textContent ?? ''));
  expect(confirm, 'el diálogo de confirmación no tiene botón para confirmar').toBeDefined();
  await user.click(confirm!);
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-11 - Acciones administrativas de recuperación', () => {
  it('restablecer contraseña pide confirmación explícita antes de llamar a la API', async () => {
    const fetchMock = usersBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderUsers();

    await openActions(user);
    await user.click(await screen.findByRole('button', { name: /restablecer contraseña/i }));

    expect(await findConfirmationDialog()).toBeInTheDocument();
    expect(mutationCalls(fetchMock, '/password/reset')).toHaveLength(0);
  });

  it('al confirmar envía POST /api/admin/users/:id/password/reset con Bearer y notifica el resultado', async () => {
    const fetchMock = usersBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderUsers();

    await openActions(user);
    await user.click(await screen.findByRole('button', { name: /restablecer contraseña/i }));
    await confirmDialog(user);

    await waitFor(() => expect(mutationCalls(fetchMock, '/admin/users/u1/password/reset')).toHaveLength(1));
    const [, init] = mutationCalls(fetchMock, '/admin/users/u1/password/reset')[0];
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer admin-token-123');
    expect((await screen.findAllByText(/contraseña restablecida|correo/i)).length).toBeGreaterThan(0);
    expect(mutationCalls(fetchMock, '/2fa/reset')).toHaveLength(0);
  });

  it('restablecer 2FA es una acción separada, confirmada, que llama a /2fa/reset y la tabla refleja el 2FA desvinculado', async () => {
    const fetchMock = usersBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderUsers();

    const table = await screen.findByRole('table');
    expect(await within(table).findByText('Activado')).toBeInTheDocument();

    await openActions(user);
    await user.click(await screen.findByRole('button', { name: /restablecer 2fa|desvincular 2fa|resetear 2fa/i }));
    expect(mutationCalls(fetchMock, '/2fa/reset')).toHaveLength(0);
    await confirmDialog(user);

    await waitFor(() => expect(mutationCalls(fetchMock, '/admin/users/u1/2fa/reset')).toHaveLength(1));
    expect(mutationCalls(fetchMock, '/password/reset')).toHaveLength(0);
    await waitFor(() => expect(within(screen.getByRole('table')).queryByText('Activado')).not.toBeInTheDocument());
  });

  it('nunca muestra secretos TOTP ni contraseñas hasheadas en la tabla', async () => {
    window.sessionStorage.setItem('centinela_access', 'admin-token-123');
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({
      summary: { total: 1, admins: 0, operators: 1 },
      users: [{
        id: 'u1', nombreCompleto: 'Ada Lovelace', email: 'ada@centinela.local', rol: 'OPERATOR', activo: true, totpVinculado: true,
        contrasenaHash: '$2a$10$abcdefghijklmnopqrstuvwxyz123456', totpSecret: 'JBSWY3DPEHPK3PXP',
      }],
    })));

    render(<MemoryRouter><Users /></MemoryRouter>);

    const table = await screen.findByRole('table');
    expect(within(table).queryByText(/JBSWY3DPEHPK3PXP/i)).not.toBeInTheDocument();
    expect(within(table).queryByText(/\$2a\$10\$/i)).not.toBeInTheDocument();
  });
});
