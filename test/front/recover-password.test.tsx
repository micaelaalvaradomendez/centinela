import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { Toaster } from '@/components/ui/toast';
import RecoverPassword from '@/pages/RecoverPassword';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

type FetchCall = [string, RequestInit];
const callsTo = (fetchMock: ReturnType<typeof vi.fn>, path: string) =>
  (fetchMock.mock.calls as FetchCall[]).filter(([url]) => url.endsWith(path));

// Contratos reales del backend (BAC-19 / BAC-20, auth_handler.go):
//   POST /api/auth/password/forgot { email }                          -> 200 { message } (genérico)
//   POST /api/auth/password/reset  { email, codigo, nuevaContrasena } -> 200 { message } | 400 RESET_FAILED
function recoveryBackend(resetResponse = () => jsonResponse({ message: 'Contraseña actualizada exitosamente. Ya puedes iniciar sesión.' })) {
  return vi.fn().mockImplementation((url: string) => {
    if (url.endsWith('/auth/password/forgot')) {
      return Promise.resolve(jsonResponse({ message: 'Si el correo está registrado, recibirás un código de recuperación en unos minutos.' }));
    }
    if (url.endsWith('/auth/password/reset')) return Promise.resolve(resetResponse());
    return Promise.resolve(jsonResponse({}));
  });
}

function renderRecover() {
  const router = createMemoryRouter([
    { path: '/recover-password', Component: RecoverPassword },
    { path: '/login', element: <h1>Iniciar sesión</h1> },
  ], { initialEntries: ['/recover-password'] });
  render(<Toaster><RouterProvider router={router} /></Toaster>);
  return router;
}

async function completeSteps(user: ReturnType<typeof userEvent.setup>, code = '123456', password = 'Nueva2026!') {
  await user.type(screen.getByLabelText(/correo electrónico/i), 'usuario@centinela.local');
  await user.click(screen.getByRole('button', { name: /enviar código de recuperación/i }));
  expect(await screen.findByText(/verifica(ción de| tu) identidad/i)).toBeInTheDocument();
  await user.type(screen.getByLabelText(/código de verificación/i), code);
  await user.click(screen.getByRole('button', { name: /verificar código/i }));
  expect(await screen.findByRole('heading', { name: /nueva contraseña/i })).toBeInTheDocument();
  await user.type(screen.getByLabelText(/^nueva contraseña$/i), password);
  await user.type(screen.getByLabelText(/repetir contraseña|confirmar nueva contraseña/i), password);
  await user.click(screen.getByRole('button', { name: /restablecer contraseña/i }));
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-12 - Vistas de recuperación de contraseña (RF-13)', () => {
  it('el paso 1 envía POST /api/auth/password/forgot con el correo y avanza a la verificación', { timeout: 15000 }, async () => {
    const fetchMock = recoveryBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderRecover();

    await user.type(screen.getByLabelText(/correo electrónico/i), 'usuario@centinela.local');
    await user.click(screen.getByRole('button', { name: /enviar código de recuperación/i }));

    expect(await screen.findByText(/verifica(ción de| tu) identidad/i)).toBeInTheDocument();
    await waitFor(() => expect(callsTo(fetchMock, '/auth/password/forgot')).toHaveLength(1));
    const [, init] = callsTo(fetchMock, '/auth/password/forgot')[0];
    expect(init.method).toBe('POST');
    expect(JSON.parse(String(init.body))).toEqual({ email: 'usuario@centinela.local' });
  });

  it('el paso 1 no llama a la API con un correo mal formado', { timeout: 15000 }, async () => {
    const fetchMock = recoveryBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderRecover();

    await user.type(screen.getByLabelText(/correo electrónico/i), 'correo-invalido');
    await user.click(screen.getByRole('button', { name: /enviar código de recuperación/i }));

    expect(callsTo(fetchMock, '/auth/password/forgot')).toHaveLength(0);
    expect(screen.queryByText(/verifica(ción de| tu) identidad/i)).not.toBeInTheDocument();
  });

  it('el paso final envía POST /api/auth/password/reset con email, código y nueva contraseña, y confirma el éxito', { timeout: 15000 }, async () => {
    const fetchMock = recoveryBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderRecover();

    await completeSteps(user, '654321');

    await waitFor(() => expect(callsTo(fetchMock, '/auth/password/reset')).toHaveLength(1));
    const [, init] = callsTo(fetchMock, '/auth/password/reset')[0];
    expect(init.method).toBe('POST');
    expect(JSON.parse(String(init.body))).toEqual({
      email: 'usuario@centinela.local',
      codigo: '654321',
      nuevaContrasena: 'Nueva2026!',
    });
    expect((await screen.findAllByText(/contraseña (actualizada|restablecida)|ya pod[ée]s iniciar sesión|ya puedes iniciar sesión/i))[0]).toBeInTheDocument();
  });

  it('un código inválido o vencido (400 RESET_FAILED) se informa en la vista sin perder el flujo', { timeout: 15000 }, async () => {
    const fetchMock = recoveryBackend(() => jsonResponse({ errorCode: 'RESET_FAILED', message: 'el código ha expirado' }, 400));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderRecover();

    await completeSteps(user);

    await waitFor(() => expect(callsTo(fetchMock, '/auth/password/reset')).toHaveLength(1));
    expect((await screen.findAllByText(/código.*(inválido|incorrecto|vencido|expir)/i))[0]).toBeVisible();
    expect(screen.queryByText(/contraseña (actualizada|restablecida)/i)).not.toBeInTheDocument();
  });

  it('al terminar con éxito redirige a /login (FIX-18)', { timeout: 15000 }, async () => {
    vi.stubGlobal('fetch', recoveryBackend());
    const user = userEvent.setup();
    const router = renderRecover();

    await completeSteps(user);

    await waitFor(() => expect(router.state.location.pathname).toBe('/login'), { timeout: 3000 });
  });

  it('el paso 2 no avanza sin un código de 6 dígitos (FIX-18)', { timeout: 15000 }, async () => {
    vi.stubGlobal('fetch', recoveryBackend());
    const user = userEvent.setup();
    renderRecover();

    await user.type(screen.getByLabelText(/correo electrónico/i), 'usuario@centinela.local');
    await user.click(screen.getByRole('button', { name: /enviar código de recuperación/i }));
    expect(await screen.findByText(/verifica(ción de| tu) identidad/i)).toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: /verificar código/i }));

    expect(screen.queryByRole('heading', { name: /nueva contraseña/i })).not.toBeInTheDocument();
  });

  it('el paso 3 no envía una contraseña que no cumple la complejidad del backend (FIX-20)', { timeout: 15000 }, async () => {
    const fetchMock = recoveryBackend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderRecover();

    await completeSteps(user, '123456', 'sinmayuscula1!');
    expect(callsTo(fetchMock, '/auth/password/reset')).toHaveLength(0);

    // Contraprueba: con una clave válida el mismo formulario sí la envía. Sin esto, la prueba
    // pasaría con una vista que nunca llama a la API.
    const password = screen.getByLabelText(/^nueva contraseña$/i);
    const confirmation = screen.getByLabelText(/repetir contraseña|confirmar nueva contraseña/i);
    await user.clear(password);
    await user.type(password, 'Valida2026!');
    await user.clear(confirmation);
    await user.type(confirmation, 'Valida2026!');
    await user.click(screen.getByRole('button', { name: /restablecer contraseña/i }));
    await waitFor(() => expect(callsTo(fetchMock, '/auth/password/reset')).toHaveLength(1));
  });
});
