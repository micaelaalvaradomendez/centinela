import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { describe, expect, it } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';

const storedUser = {
  id: 'user-1',
  organizacionId: 'org-1',
  nombreCompleto: 'Administrador de prueba',
  email: 'admin@centinela.local',
  rol: 'ADMIN',
  instanciasPermitidas: [],
  tiene2FA: true,
};

function renderApplication(path: string, authenticated: boolean) {
  if (authenticated) {
    window.localStorage.setItem('centinela_access', 'access-token');
    window.localStorage.setItem('centinela_user', JSON.stringify(storedUser));
  }
  const router = createMemoryRouter(applicationRoutes, { initialEntries: [path] });
  render(<RouterProvider router={router} />);
  return router;
}

describe('FRN-03 - navbar y rutas base', () => {
  it('redirige a login cuando no existe una sesión', async () => {
    const router = renderApplication('/dashboard', false);

    expect(await screen.findByRole('heading', { name: 'Iniciar sesión' })).toBeVisible();
    expect(router.state.location.pathname).toBe('/login');
  });

  it('permite navegar entre Dashboard e Instancias con una sesión', async () => {
    const user = userEvent.setup();
    const router = renderApplication('/dashboard', true);

    expect(await screen.findByRole('heading', { name: 'Dashboard' })).toBeVisible();
    await user.click(screen.getByRole('link', { name: 'Instancias' }));

    expect(await screen.findByRole('heading', { name: 'Instancias' })).toBeVisible();
    expect(router.state.location.pathname).toBe('/instances');
  });
});
