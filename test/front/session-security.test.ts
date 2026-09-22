import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { render } from '@testing-library/react';
import React from 'react';
import { ApiResponseNotifier } from '@/components/common/ApiResponseNotifier';
import { logoutSession } from '@/components/features/auth/services/authService';
import { mapAuthenticationError } from '@/components/features/auth/utils/mapAuthenticationError';
import {
  API_UNAUTHORIZED_EVENT,
  ApiRequestError,
  apiClient,
} from '@/services/apiClient';
import {
  clearAuthTokens,
  getAccessToken,
  getRefreshToken,
  storeAuthTokens,
} from '@/storage/tokenStorage';

function jsonResponse(body: unknown, status = 200, headers: Record<string, string> = {}) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json', ...headers },
  });
}

describe('FRN-13 - Flujo integral de logout y limpieza de sesión en cliente', () => {
  beforeEach(() => {
    window.sessionStorage.clear();
    storeAuthTokens({ accessToken: 'access-123', refreshToken: 'refresh-abc' });
    window.sessionStorage.setItem('centinela_user', JSON.stringify({ id: 'u1', email: 'test@example.com' }));
    window.sessionStorage.setItem('centinela_pending_login', 'pending-state');
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    window.sessionStorage.clear();
  });

  it('logoutSession envía la cabecera Authorization: Bearer <accessToken> junto con { refreshToken } para revocación atómica', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 204 }));
    vi.stubGlobal('fetch', fetchMock);

    await logoutSession();

    expect(fetchMock).toHaveBeenCalledTimes(1);
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toContain('/auth/logout');
    expect(init.method).toBe('POST');

    // Debe enviar el accessToken en la cabecera Authorization para que el backend pueda revocar ambos tokens
    const headers = init.headers as Record<string, string>;
    expect(headers?.Authorization).toBe('Bearer access-123');

    // Y el refreshToken en el body
    const body = init.body ? JSON.parse(String(init.body)) : {};
    expect(body.refreshToken).toBe('refresh-abc');
  });

  it('logoutSession limpia exhaustivamente sessionStorage ante fallos de red o errores HTTP', async () => {
    // Simular falla de red (rejected promise)
    const fetchMock = vi.fn().mockRejectedValue(new Error('Network error'));
    vi.stubGlobal('fetch', fetchMock);

    await logoutSession().catch(() => undefined);

    // Debe limpiar exhaustivamente las claves de sesión local
    expect(getAccessToken()).toBeNull();
    expect(getRefreshToken()).toBeNull();
    expect(window.sessionStorage.getItem('centinela_user')).toBeNull();
    expect(window.sessionStorage.getItem('centinela_pending_login')).toBeNull();
  });

  it('evento centinela:api-unauthorized con TOKEN_REVOKED limpia la sesión local y redirige a /login', async () => {
    const replaceMock = vi.fn();
    delete (window as unknown as { location: unknown }).location;
    window.location = { pathname: '/dashboard', replace: replaceMock } as unknown as Location;

    render(React.createElement(ApiResponseNotifier));

    window.dispatchEvent(new CustomEvent(API_UNAUTHORIZED_EVENT, {
      detail: { errorCode: 'TOKEN_REVOKED', message: 'Token revocado.', status: 401 },
    }));

    expect(getAccessToken()).toBeNull();
    expect(getRefreshToken()).toBeNull();
    expect(replaceMock).toHaveBeenCalledWith('/login');
  });
});

describe('SEC-02 - Cliente frontend compatible con refresh token HttpOnly', () => {
  beforeEach(() => {
    window.sessionStorage.clear();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    window.sessionStorage.clear();
  });

  it('apiClient envía credentials: include en todas las peticiones a la API para soportar cookies HttpOnly', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ok: true }));
    vi.stubGlobal('fetch', fetchMock);

    await apiClient.get('/test-secure');

    expect(fetchMock).toHaveBeenCalledTimes(1);
    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(init.credentials).toBe('include');
  });

  it('renovación de sesión es compatible con respuestas del backend que no incluyen refreshToken en el cuerpo JSON', async () => {
    // En modo cookie HttpOnly, el backend responde solo { accessToken: 'nuevo-access' } y el refreshToken viene en Set-Cookie
    storeAuthTokens({ accessToken: 'expired-access', refreshToken: 'http-only-cookie-managed' });

    const fetchMock = vi.fn().mockImplementation((url: string) => {
      if (url.includes('/auth/refresh')) {
        return Promise.resolve(jsonResponse({ accessToken: 'nuevo-access-token' }));
      }
      if (url.includes('/recurso-protegido')) {
        const calls = fetchMock.mock.calls as [string, RequestInit][];
        // Primera llamada falla con 401 para disparar la renovación
        if (calls.filter(([u]) => u.includes('/recurso-protegido')).length === 1) {
          return Promise.resolve(jsonResponse({ errorCode: 'TOKEN_EXPIRED', message: 'Expirado' }, 401));
        }
        return Promise.resolve(jsonResponse({ data: 'exito' }));
      }
      return Promise.resolve(jsonResponse({}));
    });
    vi.stubGlobal('fetch', fetchMock);

    const result = await apiClient.get<{ data: string }>('/recurso-protegido');
    expect(result.data).toBe('exito');
    expect(getAccessToken()).toBe('nuevo-access-token');
  });

  it('el cliente permite operar conservando únicamente accessToken en sessionStorage cuando se usan cookies seguras', () => {
    clearAuthTokens();
    storeAuthTokens({ accessToken: 'access-cookie-only', refreshToken: '' });

    expect(getAccessToken()).toBe('access-cookie-only');
    expect(getRefreshToken()).toBeFalsy();
  });
});

describe('FIX-08 - Auditoría del contrato de códigos de error en cliente', () => {
  it('mapAuthenticationError traduce códigos estándar del backend a mensajes legibles sin mostrar undefined', () => {
    const codigosBackend = [
      'INVALID_REQUEST',
      'TOTP_FAILED',
      'TWO_FACTOR_ALREADY_ENABLED',
      'QR_ERROR',
      'USER_CONFLICT',
      'PASSWORD_CHANGE_REQUIRED',
      'REFRESH_FAILED',
      'USER_NOT_FOUND',
      'PASSWORD_CHANGE_FAILED',
    ];

    for (const errorCode of codigosBackend) {
      const error = new ApiRequestError('Mensaje crudo de backend', 400, errorCode);
      const mapped = mapAuthenticationError(error, 'login');

      expect(mapped.errorCode).toBe(errorCode);
      expect(mapped.message).toBeDefined();
      expect(mapped.message.length).toBeGreaterThan(5);
      expect(mapped.message).not.toContain('undefined');
    }
  });

  it('mapAuthenticationError maneja AUTH_FAILED preservando el mensaje descriptivo del servidor', () => {
    const error = new ApiRequestError('Usuario o contraseña incorrectos.', 401, 'AUTH_FAILED');
    const mapped = mapAuthenticationError(error, 'login');

    expect(mapped.errorCode).toBe('AUTH_FAILED');
    expect(mapped.message).toBe('Usuario o contraseña incorrectos.');
  });
});
