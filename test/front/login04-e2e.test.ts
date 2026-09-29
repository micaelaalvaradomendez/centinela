// LOGIN-04 — Prueba integral de autenticación y autorización, de punta a punta.
//
// Ejecuta el CÓDIGO REAL DEL FRONTEND (servicios, apiClient, interceptor y almacenamiento de
// sesión) contra el BACKEND REAL levantado por test/back. No hay mocks de la API: un fetch
// que se comporta como el navegador resuelve las rutas relativas contra el backend y maneja
// cookies con Path (la cookie HttpOnly de refresh de SEC-01).
//
// La orquesta test/back (TestLOGIN04IntegracionFrontBack), que levanta el stack y define:
//   CENTINELA_E2E_API          URL del backend, p. ej. http://127.0.0.1:18080
//   CENTINELA_E2E_ADMIN_TOKEN  access token de un ADMIN
//   CENTINELA_E2E_COMPOSE      comando docker compose del stack (para leer el correo simulado)
// Sin esas variables (p. ej. `pnpm test` solo), la suite se omite.
import { execSync } from 'node:child_process';
import { createHmac } from 'node:crypto';
import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import {
  PasswordChangeRequiredError,
  fetchUserProfile,
  logoutSession,
  persistSessionFromTokens,
  submitLoginCredentials,
} from '@/components/features/auth/services/authService';
import { twoFactorService } from '@/components/features/2fa/services/2fa.service';
import { createUser } from '@/components/features/createuser/services/createUserService';
import { assignUserInstances, fetchUserInstanceInventory } from '@/components/features/users/services/userInstanceService';
import { ApiRequestError, apiClient } from '@/services/apiClient';
import { getAccessToken } from '@/storage/tokenStorage';

const API = process.env.CENTINELA_E2E_API ?? '';
const ADMIN_TOKEN = process.env.CENTINELA_E2E_ADMIN_TOKEN ?? '';
const COMPOSE = process.env.CENTINELA_E2E_COMPOSE ?? '';

// ---------- Navegador simulado: fetch con base URL y cookie jar ----------
type Cookie = { name: string; value: string; path: string };
let jar: Cookie[] = [];
const realFetch = globalThis.fetch;
const origin = API.replace(/\/api\/?$/, '');

function browserFetch(input: RequestInfo | URL, init: RequestInit = {}): Promise<Response> {
  const url = new URL(String(input), origin + '/');
  const headers = new Headers(init.headers);
  if (init.credentials === 'include') {
    const cookies = jar.filter((c) => url.pathname.startsWith(c.path));
    if (cookies.length) headers.set('Cookie', cookies.map((c) => `${c.name}=${c.value}`).join('; '));
  }
  return realFetch(url, { ...init, headers }).then((response) => {
    for (const raw of response.headers.getSetCookie?.() ?? []) {
      const [pair, ...attrs] = raw.split(';').map((part) => part.trim());
      const [name, ...rest] = pair.split('=');
      const value = rest.join('=');
      const path = attrs.find((a) => /^path=/i.test(a))?.split('=')[1] ?? '/';
      const expired = !value || attrs.some((a) => /^max-age=(0|-)/i.test(a));
      jar = jar.filter((c) => !(c.name === name && c.path === path));
      if (!expired) jar.push({ name, value, path });
    }
    return response;
  });
}

// ---------- Utilidades ----------
function totp(secretBase32: string, offsetSteps = 0): string {
  const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';
  const clean = secretBase32.replace(/=+$/, '').toUpperCase();
  let bits = '';
  for (const char of clean) bits += alphabet.indexOf(char).toString(2).padStart(5, '0');
  const key = Buffer.from(bits.match(/.{8}/g)!.map((byte) => parseInt(byte, 2)));
  const counter = Buffer.alloc(8);
  counter.writeBigUInt64BE(BigInt(Math.floor(Date.now() / 30000) + offsetSteps));
  const hmac = createHmac('sha1', key).update(counter).digest();
  const offset = hmac[hmac.length - 1] & 0xf;
  const code = (hmac.readUInt32BE(offset) & 0x7fffffff) % 1_000_000;
  return String(code).padStart(6, '0');
}

// Bandeja del MockEmailService: último secreto enviado a `email` con el marcador dado.
function mailed(email: string, marker: string): string {
  for (let attempt = 0; attempt < 20; attempt += 1) {
    const logs = execSync(`${COMPOSE} logs --no-color backend`, { encoding: 'utf8', maxBuffer: 64 * 1024 * 1024 });
    let current = false;
    let found = '';
    for (const line of logs.split('\n')) {
      if (line.includes('Destinatario: ')) { current = line.includes(`Destinatario: ${email}`); continue; }
      if (current && line.includes(marker)) {
        const match = line.match(/\*\* (\S+) \*\*/);
        if (match) found = match[1];
        current = false;
      }
    }
    if (found) return found;
    execSync('sleep 0.3');
  }
  throw new Error(`no llegó ningún correo "${marker}" a ${email}`);
}

// El backend rechaza reutilizar el período TOTP (anti-replay, BAC-11). Para no esperar 30 s
// entre logins, se retrocede el período guardado, igual que en test/back.
function allowNextTotp(email: string) {
  execSync(`${COMPOSE} exec -T db psql -U centinela_test -d centinela_test -c "UPDATE usuarios SET ultimo_totp_periodo = 0 WHERE email_usuario = '${email}';"`);
}

function asAdmin() { window.sessionStorage.setItem('centinela_access', ADMIN_TOKEN); }
function clearBrowserSession() { window.sessionStorage.clear(); window.localStorage.clear(); jar = []; }

async function expectApiError(promise: Promise<unknown>, status: number, errorCode?: string) {
  const error = await promise.then(() => null, (e: unknown) => e);
  expect(error, `se esperaba un error HTTP ${status}${errorCode ? ` ${errorCode}` : ''}`).toBeInstanceOf(ApiRequestError);
  expect((error as ApiRequestError).status).toBe(status);
  if (errorCode) expect((error as ApiRequestError).errorCode).toBe(errorCode);
}

// ---------- Recorrido ----------
describe.skipIf(!API || !ADMIN_TOKEN || !COMPOSE)('LOGIN-04 - circuito integral front ↔ back con el código real del frontend', () => {
  const stamp = Date.now();
  const email = `login04.e2e.${stamp}@elcentinela.com`;
  let userId = '';
  let temporary = '';
  let secret = '';
  let password = 'Login04e2e!';
  let revokedAccessToken = '';
  // El setup global limpia el almacenamiento después de cada caso: el recorrido guarda los
  // tokens acá y cada paso restaura la sesión que necesita. El cookie jar sí persiste.
  let changePasswordToken = '';
  let operatorToken = '';
  const restore = (token: string) => window.sessionStorage.setItem('centinela_access', token);

  beforeAll(() => { globalThis.fetch = browserFetch as typeof fetch; });
  afterAll(() => { globalThis.fetch = realFetch; });

  it('1. El ADMIN da de alta un usuario desde el panel y la clave temporal llega solo por correo (BAC-06, BAC-16)', { timeout: 30000 }, async () => {
    clearBrowserSession(); asAdmin();
    const created = await createUser({ nombreCompleto: 'Operador LOGIN-04 E2E', nombreUsuario: `login04_e2e_${stamp}`, emailUsuario: email, rol: 'OPERATOR' });
    userId = created.id;
    expect(created.rol).toBe('OPERATOR');
    expect(JSON.stringify(created)).not.toMatch(/contrasena|password/i);
    temporary = mailed(email, 'clave provisoria es:');
  });

  it('2. Login con la clave temporal: pide cambio de contraseña y enrolamiento 2FA (BAC-03, BAC-12)', { timeout: 30000 }, async () => {
    clearBrowserSession();
    const login = await submitLoginCredentials({ email, password: temporary } as never);
    expect(login.cambioContrasenaRequerido).toBe(true);
    expect(login.totpVinculado).toBe(false);

    // Omitir pasos: el JWT temporal no da acceso a rutas protegidas.
    await expectApiError(apiClient.get('/account/profile', { bearer: login.jwtTemporal } as never), 403);

    const qr = await twoFactorService.qr(login.jwtTemporal);
    secret = qr.secretoManual;
    await expectApiError(twoFactorService.verify(login.jwtTemporal, '000000'), 401);
    const tokens = await twoFactorService.verify(login.jwtTemporal, totp(secret));
    expect(JSON.stringify(tokens)).not.toContain('refreshToken');
    expect(jar.map((c) => c.name), 'SEC-01: el refresh token debe llegar como cookie HttpOnly').toContain('centinela_refresh');

    // Con el cambio pendiente, el perfil responde 403 PASSWORD_CHANGE_REQUIRED y el front deriva.
    await expect(persistSessionFromTokens(tokens)).rejects.toBeInstanceOf(PasswordChangeRequiredError);
    expect(getAccessToken()).toBe(tokens.accessToken);
    changePasswordToken = tokens.accessToken;
  });

  it('3. Cambio obligatorio de contraseña (PUT /api/account/password) y la temporal deja de servir (FRN-10)', { timeout: 30000 }, async () => {
    restore(changePasswordToken);
    await apiClient.request('/account/password', { method: 'PUT', payload: { contrasenaActual: temporary, contrasenaNueva: password } });
    clearBrowserSession();
    await expectApiError(submitLoginCredentials({ email, password: temporary } as never), 401);
  });

  it('4. Login definitivo con 2FA persistido: sesión de OPERATOR y rutas de ADMIN rechazadas sin cerrar la sesión (BAC-11, FRN-08)', { timeout: 30000 }, async () => {
    clearBrowserSession(); allowNextTotp(email);
    const login = await submitLoginCredentials({ email, password } as never);
    expect(login.totpVinculado).toBe(true);
    await expectApiError(twoFactorService.qr(login.jwtTemporal), 409); // no se puede re-enrolar
    const user = await persistSessionFromTokens(await twoFactorService.verify(login.jwtTemporal, totp(secret)));
    expect(user.rol).toBe('OPERATOR');

    await expectApiError(apiClient.get('/admin/users'), 403);
    await expectApiError(apiClient.get('/admin/audit'), 403);
    expect(getAccessToken(), 'un 403 no debe cerrar la sesión').toBeTruthy();
    operatorToken = getAccessToken()!;
  });

  it('5. El ADMIN asigna la 101 y el OPERATOR solo ve y accede a esa instancia (BAC-07, BAC-08, BAC-14, FRN-07)', { timeout: 30000 }, async () => {
    asAdmin();
    await assignUserInstances(userId, [{ vmid: 101, nivelAcceso: 'FULL_ACCESS' }]);

    window.sessionStorage.setItem('centinela_access', operatorToken);
    const inventory = await fetchUserInstanceInventory();
    expect(inventory.map((i) => i.id)).toEqual(['101']);
    await expectApiError(apiClient.get('/instances/102'), 403, 'INSTANCE_ACCESS_DENIED');
  });

  it('6. Renovación silenciosa: con el access token inválido el interceptor renueva con la cookie y reintenta (SEC-01, SEC-02)', { timeout: 30000 }, async () => {
    restore(operatorToken);
    const profileBefore = await fetchUserProfile();
    window.sessionStorage.setItem('centinela_access', 'token-vencido-o-invalido');
    const profile = await apiClient.get<{ id: string }>('/account/profile');
    expect(profile.id).toBe(profileBefore.id);
    expect(getAccessToken()).not.toBe('token-vencido-o-invalido');
    operatorToken = getAccessToken()!;
  });

  it('7. "Cerrar sesión" revoca la sesión en el servidor: el access token viejo ya no sirve (FRN-13, BAC-17)', { timeout: 30000 }, async () => {
    restore(operatorToken);
    revokedAccessToken = operatorToken;
    await logoutSession().catch(() => undefined);
    const response = await browserFetch('/api/account/profile', { headers: { Authorization: `Bearer ${revokedAccessToken}` } });
    expect(response.status, 'el access token sigue válido en el servidor después del logout (la sesión no se revocó)').toBe(401);
    expect((await response.json()).errorCode).toBe('TOKEN_REVOKED');
  });

  it('8. Recuperación de contraseña por el propio usuario, con el código enviado por correo (FRN-12, BAC-19, BAC-20)', { timeout: 30000 }, async () => {
    clearBrowserSession();
    await apiClient.post('/auth/password/forgot', { email });
    const code = mailed(email, 'código de seguridad temporal es:');
    await expectApiError(apiClient.post('/auth/password/reset', { email, codigo: code === '000000' ? '111111' : '000000', nuevaContrasena: 'Recupera04!' }), 400);
    await apiClient.post('/auth/password/reset', { email, codigo: code, nuevaContrasena: 'Recupera04!' });
    await expectApiError(submitLoginCredentials({ email, password } as never), 401); // credencial anterior
    password = 'Recupera04!';
    const login = await submitLoginCredentials({ email, password } as never);
    expect(login.jwtTemporal).toBeTruthy();
  });

  it('9. Resets administrativos de 2FA y contraseña reinician el ciclo (FRN-11, BAC-13, BAC-15)', { timeout: 30000 }, async () => {
    clearBrowserSession(); asAdmin();
    await apiClient.post(`/admin/users/${userId}/2fa/reset`, {});
    await apiClient.post(`/admin/users/${userId}/password/reset`, {});
    const newTemporary = mailed(email, 'clave provisoria es:');
    expect(newTemporary).not.toBe(temporary);

    clearBrowserSession();
    await expectApiError(submitLoginCredentials({ email, password } as never), 401);
    const login = await submitLoginCredentials({ email, password: newTemporary } as never);
    expect(login.totpVinculado).toBe(false);
    expect(login.cambioContrasenaRequerido).toBe(true);
  });

  it('10. Las acciones administrativas del recorrido quedaron auditadas y el ADMIN las ve (BAC-18)', { timeout: 30000 }, async () => {
    clearBrowserSession(); asAdmin();
    const page = await apiClient.get<{ items: { accion: string; detalles?: string }[] }>(`/admin/audit?pagina=1&tamano=200`);
    const acciones = new Set(page.items.filter((item) => String(item.detalles ?? '').includes(userId) || String(item.detalles ?? '').includes(email)).map((item) => item.accion));
    for (const accion of ['CREAR_USUARIO', 'ASIGNAR_PERMISOS', 'RESETEAR_TOTP', 'RESETEAR_CONTRASENA']) {
      expect([...acciones], `falta la acción ${accion} en la auditoría`).toContain(accion);
    }
  });
});
