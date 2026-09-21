import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiRequestError, apiClient } from '@/services/apiClient';

describe('FIX-07 - parser de errores de la API (errorCode vs code)', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('extrae correctamente el campo errorCode del cuerpo de respuesta del backend', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(
      JSON.stringify({ errorCode: 'AUTH_FAILED', message: 'Credenciales inválidas' }),
      { status: 401, headers: { 'Content-Type': 'application/json' } }
    )));

    try {
      await apiClient.request('/test-endpoint');
      expect.unreachable('Debería haber lanzado un ApiRequestError');
    } catch (error) {
      expect(error).toBeInstanceOf(ApiRequestError);
      const apiError = error as ApiRequestError;
      expect(apiError.status).toBe(401);
      expect(apiError.errorCode).toBe('AUTH_FAILED');
      expect(apiError.message).toBe('Credenciales inválidas');
    }
  });

  it('mantiene compatibilidad con respuestas que incluyan el campo legacy code', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(
      JSON.stringify({ code: 'LEGACY_ERROR', message: 'Error de formato' }),
      { status: 400, headers: { 'Content-Type': 'application/json' } }
    )));

    try {
      await apiClient.request('/test-endpoint');
      expect.unreachable('Debería haber lanzado un ApiRequestError');
    } catch (error) {
      expect(error).toBeInstanceOf(ApiRequestError);
      const apiError = error as ApiRequestError;
      expect(apiError.status).toBe(400);
      expect(apiError.errorCode).toBe('LEGACY_ERROR');
    }
  });

  it('asigna undefined a errorCode cuando el cuerpo no especifica ningún código de error', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(
      JSON.stringify({ message: 'Error sin código' }),
      { status: 500, headers: { 'Content-Type': 'application/json' } }
    )));

    try {
      await apiClient.request('/test-endpoint');
      expect.unreachable('Debería haber lanzado un ApiRequestError');
    } catch (error) {
      expect(error).toBeInstanceOf(ApiRequestError);
      const apiError = error as ApiRequestError;
      expect(apiError.status).toBe(500);
      expect(apiError.errorCode).toBeUndefined();
    }
  });
});

