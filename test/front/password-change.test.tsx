import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

// El layout de autenticación monta InfraDeployCheck, que pide GET /api/version.
// Los tests filtran las llamadas por URL para no confundirse con esa sonda.
function callsTo(fetchMock: ReturnType<typeof vi.fn>, fragment: string) {
  return (fetchMock.mock.calls as [string, RequestInit][]).filter(([url]) => String(url).includes(fragment));
}

function seedPendingTwoFactor() {
  window.sessionStorage.setItem('centinela_pending_login', JSON.stringify({
    jwtTemporal: 'jwt-temporal',
    totpVinculado: true,
    cambioContrasenaRequerido: true,
  }));
}

function renderApp(path: string) {
  const router = createMemoryRouter(applicationRoutes, { initialEntries: [path] });
  render(<RouterProvider router={router} />);
  return router;
}

async function fillChangeForm(user: ReturnType<typeof userEvent.setup>, nueva: string, confirmacion = nueva) {
  await user.type(await screen.findByLabelText('Contraseña actual'), 'Temp123!');
  await user.type(screen.getByLabelText('Contraseña nueva'), nueva);
  await user.type(screen.getByLabelText('Confirmar contraseña nueva'), confirmacion);
  await user.click(screen.getByRole('button', { name: 'Cambiar contraseña' }));
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-10 - Cambio obligatorio de contraseña temporal', () => {
  it('deriva a /change-password cuando el perfil responde 403 PASSWORD_CHANGE_REQUIRED tras el 2FA', { timeout: 10000 }, async () => {
    const fetchMock = vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/account/profile')) {
        return Promise.resolve(jsonResponse(
          { errorCode: 'PASSWORD_CHANGE_REQUIRED', message: 'Debe cambiar su contraseña temporal antes de continuar.' },
          403,
        ));
      }
      if (String(url).includes('/auth/2fa/verify')) {
        return Promise.resolve(jsonResponse({ accessToken: 'access-token', expiresIn: 3600 }));
      }
      return Promise.resolve(jsonResponse({}));
    });
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    seedPendingTwoFactor();

    const router = renderApp('/two-factor/verify');

    await user.type(await screen.findByRole('textbox'), '123456');
    await user.click(screen.getByRole('button', { name: 'Verificar código' }));

    await waitFor(() => expect(router.state.location.pathname).toBe('/change-password'), { timeout: 5000 });
    expect(await screen.findByRole('heading', { name: 'Cambio obligatorio de contraseña' }, { timeout: 5000 })).toBeVisible();
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(window.sessionStorage.getItem('centinela_access')).toBe('access-token');
  });

  it('con el cambio pendiente la navegación general queda bloqueada', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({})));
    // Estado real tras el 403: hay access token pero el perfil (centinela_user) no se guardó.
    window.sessionStorage.setItem('centinela_access', 'access-token');

    for (const path of ['/dashboard', '/instances', '/users']) {
      const router = renderApp(path);
      await waitFor(() => expect(router.state.location.pathname).not.toBe(path));
    }
  });

  it('la pantalla de cambio ofrece cerrar sesión como única alternativa', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 204 }));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    const router = renderApp('/change-password');

    await user.click(await screen.findByRole('button', { name: /cerrar sesión/i }));

    await waitFor(() => expect(router.state.location.pathname).toBe('/login'));
    expect(callsTo(fetchMock, '/auth/logout')).toHaveLength(1);
    expect(window.sessionStorage.getItem('centinela_access')).toBeNull();
  });

  it('envía PUT /api/account/password y, al terminar, vuelve al login para completar el 2FA con la clave nueva', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ message: 'Contraseña actualizada correctamente.' })));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    const router = renderApp('/change-password');
    await fillChangeForm(user, 'Nueva1234!');

    await waitFor(() => expect(router.state.location.pathname).toBe('/login'), { timeout: 5000 });
    const [url, init] = callsTo(fetchMock, '/account/password')[0];
    expect(url).toContain('/account/password');
    expect(init.method).toBe('PUT');
    expect(JSON.parse(String(init.body))).toEqual({ contrasenaActual: 'Temp123!', contrasenaNueva: 'Nueva1234!' });
    expect(window.sessionStorage.getItem('centinela_access')).toBeNull();
  });

  it('muestra el error del backend si rechaza la contraseña nueva', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse(
      { errorCode: 'PASSWORD_CHANGE_FAILED', message: 'La contraseña debe contener al menos un carácter especial.' },
      400,
    )));
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    const router = renderApp('/change-password');
    await fillChangeForm(user, 'Nueva1234!');

    expect(await screen.findByText('La contraseña debe contener al menos un carácter especial.')).toBeVisible();
    expect(router.state.location.pathname).toBe('/change-password');
  });

  it('no llama a la API si las contraseñas nuevas no coinciden', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ message: 'ok' })));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    renderApp('/change-password');
    await fillChangeForm(user, 'Nueva1234!', 'Distinta99!');

    expect(await screen.findByText('Las contraseñas nuevas no coinciden.')).toBeVisible();
    expect(callsTo(fetchMock, '/account/password')).toHaveLength(0);
  });

  // Mismas reglas que backend/internal/infrastructure/crypto/password.go (ValidarComplejidadContrasena).
  it.each([
    ['sin mayúscula', 'nueva1234!'],
    ['sin dígito', 'NuevaClave!'],
    ['sin carácter especial', 'Nueva12345'],
  ])('FIX-29 no llama a la API si la contraseña nueva no cumple la complejidad del backend (%s)', async (_case, password) => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ message: 'ok' })));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    renderApp('/change-password');
    await fillChangeForm(user, password);

    await waitFor(() => expect(screen.getByRole('button', { name: 'Cambiar contraseña' })).toBeEnabled());
    expect(callsTo(fetchMock, '/account/password')).toHaveLength(0);
  });

  it('no llama a la API si la contraseña nueva queda fuera de 8-12 caracteres', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ message: 'ok' })));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    renderApp('/change-password');
    await fillChangeForm(user, 'Corta1!');

    expect(await screen.findByText(/8 y 12 caracteres/i)).toBeVisible();
    expect(callsTo(fetchMock, '/account/password')).toHaveLength(0);
  });
});
