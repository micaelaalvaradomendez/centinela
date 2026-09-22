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

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-09 - cambio obligatorio de contraseña', () => {
  it('deriva a /change-password cuando el perfil responde 403 PASSWORD_CHANGE_REQUIRED', async () => {
    const fetchMock = vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/account/profile')) {
        return Promise.resolve(jsonResponse(
          { errorCode: 'PASSWORD_CHANGE_REQUIRED', message: 'Debe cambiar su contraseña temporal antes de continuar.' },
          403,
        ));
      }
      if (String(url).includes('/auth/2fa/verify')) {
        return Promise.resolve(jsonResponse({ accessToken: 'access-token', refreshToken: 'refresh-token', expiresIn: 3600 }));
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
    // No debe quedar la pantalla de error del 2FA.
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    // Los tokens quedaron guardados para que la pantalla de cambio pueda autenticarse.
    expect(window.sessionStorage.getItem('centinela_access')).toBe('access-token');
  });

  it('envía PUT /account/password y navega a /login al actualizar la contraseña', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ message: 'Contraseña actualizada correctamente.' })));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    const router = renderApp('/change-password');

    await user.type(await screen.findByLabelText('Contraseña actual'), 'Temp123!');
    await user.type(screen.getByLabelText('Contraseña nueva'), 'Nueva1234');
    await user.type(screen.getByLabelText('Confirmar contraseña nueva'), 'Nueva1234');
    await user.click(screen.getByRole('button', { name: 'Cambiar contraseña' }));

    await waitFor(() => expect(router.state.location.pathname).toBe('/login'), { timeout: 5000 });

    const [url, init] = callsTo(fetchMock, '/account/password')[0];
    expect(url).toContain('/account/password');
    expect(init.method).toBe('PUT');
    expect(JSON.parse(String(init.body))).toEqual({
      contrasenaActual: 'Temp123!',
      contrasenaNueva: 'Nueva1234',
    });
    expect(window.sessionStorage.getItem('centinela_access')).toBeNull();
  });

  it('no llama a la API si las contraseñas nuevas no coinciden', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ message: 'ok' })));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    renderApp('/change-password');

    await user.type(await screen.findByLabelText('Contraseña actual'), 'Temp123!');
    await user.type(screen.getByLabelText('Contraseña nueva'), 'Nueva1234');
    await user.type(screen.getByLabelText('Confirmar contraseña nueva'), 'Distinta99');
    await user.click(screen.getByRole('button', { name: 'Cambiar contraseña' }));

    expect(await screen.findByText('Las contraseñas nuevas no coinciden.')).toBeVisible();
    expect(callsTo(fetchMock, '/account/password')).toHaveLength(0);
  });

  it('no llama a la API si la contraseña nueva queda fuera de 8-12 caracteres', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ message: 'ok' })));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    window.sessionStorage.setItem('centinela_access', 'access-token');

    renderApp('/change-password');

    await user.type(await screen.findByLabelText('Contraseña actual'), 'Temp123!');
    await user.type(screen.getByLabelText('Contraseña nueva'), 'Corta1');
    await user.type(screen.getByLabelText('Confirmar contraseña nueva'), 'Corta1');
    await user.click(screen.getByRole('button', { name: 'Cambiar contraseña' }));

    expect(await screen.findByText('La contraseña nueva debe tener entre 8 y 12 caracteres.')).toBeVisible();
    expect(callsTo(fetchMock, '/account/password')).toHaveLength(0);
  });
});
