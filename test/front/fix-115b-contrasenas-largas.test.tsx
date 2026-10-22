import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import fs from 'node:fs';
import path from 'node:path';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { validatePasswordComplexity } from '@/components/features/auth/utils/validateAuthenticationFields';
import { applicationRoutes } from '@/routes/applicationRoutes';
import { withAppProviders } from './app-providers';

// FIX-115B (NIST SP 800-63B, depende de FIX-115A): el frontend acepta contraseñas de hasta 128
// caracteres: la validación, los formularios (maxLength) y los textos de ayuda ya no hablan de 12.
// Criterio de éxito: se valida y envía una contraseña larga sin bloquear el botón ni mostrar errores espurios.

const jsonResponse = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

// Cumple mayúscula, número y símbolo, con la longitud pedida.
const password = (length: number) => 'Aa1!' + 'x'.repeat(length - 4);

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FIX-115B - validatePasswordComplexity admite hasta 128 caracteres', () => {
  it.each([8, 12, 13, 20, 32, 64, 128])('acepta una contraseña de %i caracteres', (length) => {
    const result = validatePasswordComplexity(password(length));
    expect(result.missingConditions, `${length} caracteres`).toEqual([]);
    expect(result.isValid).toBe(true);
  });

  it('rechaza 129 caracteres (cota anti-DoS de bcrypt) y 7 caracteres', () => {
    expect(validatePasswordComplexity(password(129)).isValid).toBe(false);
    expect(validatePasswordComplexity(password(7)).isValid).toBe(false);
  });

  it('el mensaje de longitud ya no menciona el máximo de 12', () => {
    const { missingConditions } = validatePasswordComplexity('abc');
    expect(missingConditions.join(' ')).not.toMatch(/\b12\b/);
  });
});

describe('FIX-115B - formularios de contraseña', () => {
  it('ChangePassword envía una contraseña de 64 caracteres sin errores de validación', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(jsonResponse({ message: 'Contraseña actualizada correctamente.' })));
    vi.stubGlobal('fetch', fetchMock);
    window.sessionStorage.setItem('centinela_access', 'access-token');
    const user = userEvent.setup();

    const router = createMemoryRouter(applicationRoutes, { initialEntries: ['/change-password'] });
    render(withAppProviders(<RouterProvider router={router} />));

    const larga = password(64);
    const nueva = await screen.findByLabelText('Contraseña nueva');
    expect(Number(nueva.getAttribute('maxlength') ?? 128), 'el input limita por debajo de 128').toBeGreaterThanOrEqual(128);
    await user.type(await screen.findByLabelText('Contraseña actual'), 'Temp123!');
    await user.type(nueva, larga);
    await user.type(screen.getByLabelText('Confirmar contraseña nueva'), larga);
    await user.click(screen.getByRole('button', { name: 'Cambiar contraseña' }));

    await waitFor(() => {
      const calls = (fetchMock.mock.calls as [string, RequestInit][]).filter(([url]) => String(url).includes('/account/password'));
      expect(calls, 'no se envió PUT /account/password con la contraseña larga').toHaveLength(1);
      expect(JSON.parse(String(calls[0][1].body)).contrasenaNueva).toBe(larga);
    });
    expect(screen.queryByText(/entre 8 y 12|m[aá]ximo de 12/i)).toBeNull();
  });

  it('los textos y reglas de ChangePassword, RecoverPassword y la validación no fijan el máximo en 12', () => {
    const src = path.resolve(__dirname, '../../frontend/centinela/src');
    for (const file of [
      'pages/ChangePassword.tsx', 'pages/RecoverPassword.tsx',
      'components/features/auth/utils/validateAuthenticationFields.ts',
    ]) {
      const source = fs.readFileSync(path.join(src, file), 'utf-8');
      expect(source, `${file} menciona "entre 8 y 12 caracteres"`).not.toMatch(/entre 8 y 12\b/i);
      expect(source, `${file} compara la longitud contra 12`).not.toMatch(/length\s*(<=|>)\s*12\b/);
    }
  });
});

