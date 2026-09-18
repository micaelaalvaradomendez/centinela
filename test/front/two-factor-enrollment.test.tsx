import { render, screen } from '@testing-library/react';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { LoginContinuation } from '@/components/features/auth/components/LoginContinuation';

const pendingSession = {
  jwtTemporal: 'jwt-temporal',
  totpVinculado: false,
  cambioContrasenaRequerido: false,
};

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

describe('FRN-09 - vinculación 2FA mediante QR', () => {
  it('muestra QR, clave manual y formulario de seis dígitos para una cuenta sin 2FA', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      qrBase64: 'data:image/png;base64,AAAA',
      secretoManual: 'JBSWY3DPEHPK3PXP',
    }));
    vi.stubGlobal('fetch', fetchMock);

    const router = createMemoryRouter([
      {
        path: '/two-factor/setup',
        Component: LoginContinuation,
        loader: () => pendingSession,
      },
    ], { initialEntries: ['/two-factor/setup'] });
    render(<RouterProvider router={router} />);

    expect(await screen.findByRole('img', { name: 'Código QR de configuración' })).toHaveAttribute(
      'src',
      'data:image/png;base64,AAAA',
    );
    expect(screen.getByText('JBSW Y3DP EHPK 3PXP')).toBeVisible();
    expect(screen.getByRole('textbox')).toBeVisible();
    expect(screen.getByRole('button', { name: 'Verificar código' })).toBeDisabled();
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/auth/2fa/qr',
      expect.objectContaining({
        method: 'GET',
        headers: expect.objectContaining({ Authorization: 'Bearer jwt-temporal' }),
      }),
    );
  });
});

afterEach(() => {
  vi.unstubAllGlobals();
});
