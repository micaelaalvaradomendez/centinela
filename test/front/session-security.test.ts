import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { ApiResponseNotifier } from '@/components/common/ApiResponseNotifier';
import { useLogout } from '@/components/features/auth/hooks/useAuth';
import { logoutSession, persistSessionFromTokens } from '@/components/features/auth/services/authService';
import { mapAuthenticationError } from '@/components/features/auth/utils/mapAuthenticationError';
import { isTokenResponse } from '@/components/features/auth/utils/validateAuthenticationResponses';
import { API_UNAUTHORIZED_EVENT, ApiRequestError, apiClient } from '@/services/apiClient';
import * as tokenStorage from '@/storage/tokenStorage';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

const sessionKeys = ['centinela_access', 'centinela_refresh', 'centinela_user', 'centinela_pending_login'];

// Siembra la sesión tal como la deja el login real: access token y login pendiente en
// sessionStorage (storage/tokenStorage.ts) y perfil en localStorage (services/api.js).
// Tras SEC-02 la app ya no escribe centinela_refresh; igual se verifica que no quede.
function seedFullSession() {
  window.sessionStorage.setItem('centinela_access', 'access-123');
  window.sessionStorage.setItem('centinela_pending_login', JSON.stringify({ jwtTemporal: 'x', totpVinculado: true, cambioContrasenaRequerido: false }));
  window.localStorage.setItem('centinela_user', JSON.stringify({ id: 'u1', rol: 'ADMIN' }));
}

function remainingSessionKeys() {
  return sessionKeys.filter((key) => window.sessionStorage.getItem(key) !== null || window.localStorage.getItem(key) !== null);
}

function allStoredValues() {
  const values: string[] = [];
  for (const storage of [window.sessionStorage, window.localStorage]) {
    for (let index = 0; index < storage.length; index += 1) {
      const key = storage.key(index)!;
      values.push(`${key}=${storage.getItem(key)}`);
    }
  }
  return values.join('\n');
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-13 - Flujo integral de logout y limpieza de sesión en cliente', () => {
  beforeEach(seedFullSession);

  it('logoutSession envía POST /auth/logout con Authorization: Bearer <accessToken> y credentials: include', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 204 }));
    vi.stubGlobal('fetch', fetchMock);

    await logoutSession();

    const logoutCall = (fetchMock.mock.calls as [string, RequestInit][]).find(([url]) => url.includes('/auth/logout'));
    expect(logoutCall, 'no se llamó a /auth/logout').toBeDefined();
    const [, init] = logoutCall!;
    expect(init.method).toBe('POST');
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer access-123');
    expect(init.credentials).toBe('include');
  });

  it.each([
    ['caída de red', () => Promise.reject(new TypeError('Network error'))],
    ['error HTTP 500', () => Promise.resolve(jsonResponse({ errorCode: 'INTERNAL_ERROR', message: 'x' }, 500))],
    ['sesión ya vencida (401)', () => Promise.resolve(jsonResponse({ errorCode: 'TOKEN_REVOKED', message: 'x' }, 401))],
  ])('logoutSession limpia todas las claves de sesión ante %s', async (_case, response) => {
    vi.stubGlobal('fetch', vi.fn().mockImplementation(response));

    await logoutSession().catch(() => undefined);

    expect(remainingSessionKeys()).toEqual([]);
  });

  it('el botón de logout (useLogout) redirige a /login con replace para impedir volver con el historial', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(null, { status: 204 })));
    function LogoutButton() {
      const { logout } = useLogout();
      return React.createElement('button', { type: 'button', onClick: () => void logout() }, 'Cerrar sesión');
    }
    const router = createMemoryRouter([
      { path: '/dashboard', Component: LogoutButton },
      { path: '/login', element: React.createElement('h1', null, 'Login') },
    ], { initialEntries: ['/dashboard'] });
    render(React.createElement(RouterProvider, { router }));

    await userEvent.setup().click(screen.getByRole('button', { name: 'Cerrar sesión' }));

    await waitFor(() => expect(router.state.location.pathname).toBe('/login'));
    expect(router.state.historyAction).toBe('REPLACE');
    expect(remainingSessionKeys()).toEqual([]);
  });

  it('un 401 TOKEN_REVOKED global limpia la sesión local y redirige a /login', () => {
    const replaceMock = vi.fn();
    const originalLocation = window.location;
    Object.defineProperty(window, 'location', { configurable: true, value: { ...originalLocation, pathname: '/dashboard', replace: replaceMock } });
    try {
      render(React.createElement(ApiResponseNotifier));
      window.dispatchEvent(new CustomEvent(API_UNAUTHORIZED_EVENT, {
        detail: { errorCode: 'TOKEN_REVOKED', message: 'Token revocado.', status: 401 },
      }));

      expect(remainingSessionKeys()).toEqual([]);
      expect(replaceMock).toHaveBeenCalledWith('/login');
    } finally {
      Object.defineProperty(window, 'location', { configurable: true, value: originalLocation });
    }
  });
});

describe('SEC-02 - Cliente frontend compatible con refresh token HttpOnly', () => {
  it('tokenStorage ya no expone lectura ni escritura del refresh token', () => {
    expect(Object.keys(tokenStorage)).not.toContain('getRefreshToken');
    tokenStorage.storeAuthTokens({ accessToken: 'solo-access', refreshToken: 'rt-no-debe-guardarse' } as never);
    expect(tokenStorage.getAccessToken()).toBe('solo-access');
    expect(allStoredValues()).not.toContain('rt-no-debe-guardarse');
  });

  it('acepta la respuesta de 2FA/refresh sin refreshToken y no persiste el que llegue en el body', async () => {
    expect(isTokenResponse({ accessToken: 'a', expiresIn: 3600 })).toBe(true);

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({
      id: 'u1', organizacionId: 'org-1', nombreCompleto: 'Ada Lovelace', nombreUsuario: 'ada',
      emailUsuario: 'ada@centinela.local', rol: 'ADMIN', totpVinculado: true, instanciasPermitidas: [],
    })));
    await persistSessionFromTokens({ accessToken: 'access-nuevo', refreshToken: 'rt-legacy-body', expiresIn: 3600 } as never);

    expect(tokenStorage.getAccessToken()).toBe('access-nuevo');
    expect(allStoredValues()).not.toContain('rt-legacy-body');
  });

  it('el interceptor renueva la sesión con la cookie: POST /auth/refresh sin refreshToken en el body y reintenta', async () => {
    window.sessionStorage.setItem('centinela_access', 'expired-access');
    let protectedCalls = 0;
    const fetchMock = vi.fn().mockImplementation((url: string) => {
      if (url.includes('/auth/refresh')) return Promise.resolve(jsonResponse({ accessToken: 'nuevo-access', expiresIn: 3600 }));
      protectedCalls += 1;
      return Promise.resolve(protectedCalls === 1
        ? jsonResponse({ errorCode: 'INVALID_TOKEN', message: 'Expirado' }, 401)
        : jsonResponse({ data: 'exito' }));
    });
    vi.stubGlobal('fetch', fetchMock);

    const result = await apiClient.get<{ data: string }>('/recurso-protegido');

    expect(result.data).toBe('exito');
    expect(tokenStorage.getAccessToken()).toBe('nuevo-access');
    const refreshCall = (fetchMock.mock.calls as [string, RequestInit][]).find(([url]) => url.includes('/auth/refresh'));
    expect(refreshCall, 'no se intentó renovar la sesión').toBeDefined();
    expect(refreshCall![1].credentials).toBe('include');
    expect(String(refreshCall![1].body ?? '')).not.toMatch(/refresh/i);
  });

  it('logout no envía el refresh token en el body', async () => {
    seedFullSession();
    const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 204 }));
    vi.stubGlobal('fetch', fetchMock);

    await logoutSession();

    const [, init] = (fetchMock.mock.calls as [string, RequestInit][]).find(([url]) => url.includes('/auth/logout'))!;
    expect(String(init.body ?? '')).not.toMatch(/refresh/i);
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
