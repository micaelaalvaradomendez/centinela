import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { describe, expect, it, vi } from 'vitest';
import LoginPage from '@/pages/Login';

function renderLogin(action = vi.fn()) {
  const router = createMemoryRouter([
    { path: '/login', Component: LoginPage, action },
  ], { initialEntries: ['/login'] });

  render(<RouterProvider router={router} />);
  return { action, router };
}

describe('FRN-01 - formulario de login', () => {
  it('muestra credenciales, recordarme, envío y estado inicial', async () => {
    renderLogin();

    expect(await screen.findByRole('heading', { name: 'Iniciar sesión' })).toBeVisible();
    expect(screen.getByLabelText('Correo electrónico')).toHaveAttribute('type', 'email');
    expect(screen.getByLabelText('Contraseña')).toHaveAttribute('type', 'password');
    expect(screen.getByRole('checkbox', { name: 'Recordarme' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Iniciar sesión' })).toBeEnabled();
  });

  it('muestra carga mientras se procesa un formulario válido', async () => {
    let finishAction: ((value: null) => void) | undefined;
    const action = vi.fn(() => new Promise<null>((resolve) => {
      finishAction = resolve;
    }));
    const user = userEvent.setup();
    renderLogin(action);

    await user.type(await screen.findByLabelText('Correo electrónico'), 'admin@centinela.local');
    await user.type(screen.getByLabelText('Contraseña'), 'Admin123!');
    await user.click(screen.getByRole('button', { name: 'Iniciar sesión' }));

    expect(await screen.findByRole('button', { name: 'Iniciando sesión…' })).toBeDisabled();
    finishAction?.(null);
  });
});

describe('FRN-02 - validación del login', () => {
  it('no envía campos vacíos y muestra los errores junto a cada campo', async () => {
    const user = userEvent.setup();
    const { action } = renderLogin();

    await user.click(await screen.findByRole('button', { name: 'Iniciar sesión' }));

    expect(action).not.toHaveBeenCalled();
    expect(screen.getByText('Ingresá tu correo electrónico.')).toBeVisible();
    expect(screen.getByText('Ingresá tu contraseña.')).toBeVisible();
    expect(screen.getByLabelText('Correo electrónico')).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByLabelText('Contraseña')).toHaveAttribute('aria-invalid', 'true');
  });

  it('rechaza un correo con formato inválido antes de enviarlo', async () => {
    const user = userEvent.setup();
    const { action } = renderLogin();

    await user.type(await screen.findByLabelText('Correo electrónico'), 'correo-invalido');
    await user.type(screen.getByLabelText('Contraseña'), 'Admin123!');
    await user.click(screen.getByRole('button', { name: 'Iniciar sesión' }));

    expect(action).not.toHaveBeenCalled();
    expect(screen.getByText(/Ingresá un correo completo y válido/)).toBeVisible();
  });
});
