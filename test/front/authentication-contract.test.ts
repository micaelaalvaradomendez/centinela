import { afterEach, describe, expect, it, vi } from 'vitest';
import { submitLoginAction } from '@/components/features/auth/routes/authenticationActions';
import { submitLoginCredentials } from '@/components/features/auth/services/authService';
import { twoFactorService } from '@/components/features/2fa/services/2fa.service';

const backendLoginResponse = {
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

function requestDetails(fetchMock: ReturnType<typeof vi.fn>) {
  const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit];
  return {
    url,
    method: init.method,
    headers: init.headers as Record<string, string>,
    body: init.body ? JSON.parse(String(init.body)) : undefined,
  };
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-04 - contrato de login con el backend', () => {
  it('envía email y password con los nombres aceptados por Go', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(backendLoginResponse));
    vi.stubGlobal('fetch', fetchMock);

    await submitLoginCredentials({
      email: ' admin@centinela.local ',
      password: 'Admin123!',
      recordarSesion: false,
    }).catch(() => undefined);

    expect(requestDetails(fetchMock)).toMatchObject({
      url: '/api/auth/login',
      method: 'POST',
      body: {
        email: 'admin@centinela.local',
        password: 'Admin123!',
      },
    });
  });

  it('acepta la respuesta pre-2FA emitida actualmente por el backend', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse(backendLoginResponse)));

    await expect(submitLoginCredentials({
      email: 'admin@centinela.local',
      password: 'Admin123!',
      recordarSesion: false,
    })).resolves.toMatchObject(backendLoginResponse);
  });
});

describe('LOGIN-02 - contrato de 2FA con el backend', () => {
  it('obtiene el QR con GET y el JWT temporal en Authorization', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      qrBase64: 'data:image/png;base64,AAAA',
      secretoManual: 'JBSWY3DPEHPK3PXP',
    }));
    vi.stubGlobal('fetch', fetchMock);

    await twoFactorService.qr('jwt-temporal').catch(() => undefined);

    expect(requestDetails(fetchMock)).toMatchObject({
      url: '/api/auth/2fa/qr',
      method: 'GET',
      headers: { Authorization: 'Bearer jwt-temporal' },
    });
  });

  it('verifica codigo con el JWT temporal y acepta los tokens definitivos', async () => {
    const tokenResponse = {
      accessToken: 'access-token',
      refreshToken: 'refresh-token',
      expiresIn: 28800,
    };
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(tokenResponse));
    vi.stubGlobal('fetch', fetchMock);

    await twoFactorService.verify('jwt-temporal', '123456').catch(() => undefined);

    expect(requestDetails(fetchMock)).toMatchObject({
      url: '/api/auth/2fa/verify',
      method: 'POST',
      headers: { Authorization: 'Bearer jwt-temporal' },
      body: { codigo: '123456' },
    });
  });
});

describe('LOGIN-03 - continuidad del flujo de autenticación', () => {
  it('redirige al enrolamiento 2FA después de un login válido sin TOTP', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse(backendLoginResponse)));
    const request = new Request('http://localhost/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: 'admin@centinela.local',
        password: 'Admin123!',
        recordarSesion: false,
      }),
    });

    const result = await submitLoginAction({ request });

    expect(result).toBeInstanceOf(Response);
    expect((result as Response).status).toBe(302);
    expect((result as Response).headers.get('Location')).toBe('/two-factor/setup');
  });
});
