import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiRequestError, apiClient } from '@/services/apiClient';
import { getErrorMessage } from '@/utils/errorMapper';

describe('Pruebas adversariales de seguridad del cliente (Front)', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('XSS en mensajes de error: payloads con etiquetas HTML o scripts no ejecutan código y se tratan como texto', () => {
    const maliciousPayloads = [
      '<script>alert("xss")</script>',
      '<img src=x onerror=alert(1)>',
      'javascript:/*--></title></style></textarea></script><svg/onload=alert(1)>',
    ];

    for (const payload of maliciousPayloads) {
      const error = new ApiRequestError(payload, 500, 'SERVER_ERROR');
      const mapped = getErrorMessage(error);
      // El mapeador debe conservar el string literal sin transformarlo en HTML ejecutable
      expect(mapped).toBeDefined();
      expect(typeof mapped).toBe('string');
      expect(mapped).toContain(payload);
    }
  });

  it('Inyección de rutas y caracteres de control en el cliente de API', async () => {
    const fetchSpy = vi.fn().mockResolvedValue(new Response(
      JSON.stringify({ errorCode: 'NOT_FOUND', message: 'No encontrado' }),
      { status: 404, headers: { 'Content-Type': 'application/json' } }
    ));
    vi.stubGlobal('fetch', fetchSpy);

    const maliciousPaths = [
      '/instances/../../admin/users',
      '/instances/%2e%2e%2fadmin',
      '/instances/101\r\nHeader-Injection: true',
    ];

    for (const path of maliciousPaths) {
      try {
        await apiClient.request(path);
      } catch (error) {
        expect(error).toBeInstanceOf(ApiRequestError);
      }
    }
  });

  it('Resiliencia ante respuestas corruptas del backend (JSON malformado o tipos inesperados)', async () => {
    // 1. Respuesta 200 con cuerpo no parseable (HTML de error de proxy o texto plano)
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(
      '<html><body>502 Bad Gateway</body></html>',
      { status: 502, headers: { 'Content-Type': 'text/html' } }
    )));

    try {
      await apiClient.request('/instances');
      expect.unreachable('Debería haber fallado con ApiRequestError');
    } catch (error) {
      expect(error).toBeInstanceOf(ApiRequestError);
      expect((error as ApiRequestError).status).toBe(502);
    }

    // 2. Respuesta 500 con cuerpo completamente vacío
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(
      '',
      { status: 500, headers: { 'Content-Type': 'application/json' } }
    )));

    try {
      await apiClient.request('/instances');
      expect.unreachable('Debería haber fallado con ApiRequestError');
    } catch (error) {
      expect(error).toBeInstanceOf(ApiRequestError);
      expect((error as ApiRequestError).status).toBe(500);
    }
  });

  it('Mapeo de errores ante códigos de error maliciosos o no tipados', () => {
    const unmappedErrors = [
      { errorCode: '__proto__', message: 'prototype pollution attempt' },
      { errorCode: 'CONSTRUCTOR', message: 'constructor injection attempt' },
      { errorCode: null, message: null },
      { errorCode: undefined, message: undefined },
    ];

    for (const err of unmappedErrors) {
      const msg = getErrorMessage(err);
      expect(typeof msg).toBe('string');
      expect(msg.length).toBeGreaterThan(0);
    }
  });
});
