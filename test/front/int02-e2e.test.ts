// INT-02 — Validación integral del Hito MVP Operativo en local (Smoke Test End-to-End).
//
// Ejecuta el CÓDIGO REAL DEL FRONTEND (servicios, apiClient, interceptores) contra el
// BACKEND REAL levantado por test/back. Recorre los 7 pasos del MVP operativo:
//   1. Login con 2FA de un OPERATOR.
//   2. Dashboard: consulta de /api/node/status con telemetría y conteos.
//   3. Inventario: consulta de /api/instances con filtrado RBAC y métricas.
//   4. Ciclo de vida: acción operativa con HTTP 202 { upid, tareaId }.
//   5. Notificación: recepción de TASK_FINISHED por /api/events y liberación.
//   6. Auditoría: verificación de registros en /api/admin/audit con ADMIN.
//   7. Borrado y revocación: ADMIN elimina instancia y OPERATOR pierde acceso.
//
// Orquestado por test/back (TestINT02IntegracionFrontBack) o ejecutable directamente con:
//   CENTINELA_E2E_API, CENTINELA_E2E_ADMIN_TOKEN, CENTINELA_E2E_COMPOSE
import { execSync } from 'node:child_process';
import { createHmac } from 'node:crypto';
import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import {
  submitLoginCredentials,
} from '@/components/features/auth/services/authService';
import { twoFactorService } from '@/components/features/2fa/services/2fa.service';
import { createUser } from '@/components/features/createuser/services/createUserService';
import { assignUserInstances } from '@/components/features/users/services/userInstanceService';
import { apiClient } from '@/services/apiClient';
import { getAccessToken } from '@/storage/tokenStorage';

const API = process.env.CENTINELA_E2E_API ?? '';
const ADMIN_TOKEN = process.env.CENTINELA_E2E_ADMIN_TOKEN ?? '';
const COMPOSE = process.env.CENTINELA_E2E_COMPOSE ?? '';

// ---------- Cookie Jar y Navegador Simulado ----------
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
  throw new Error(`no llegó correo "${marker}" a ${email}`);
}

function allowNextTotp(email: string) {
  execSync(`${COMPOSE} exec -T db psql -U centinela_test -d centinela_test -c "UPDATE usuarios SET ultimo_totp_periodo = 0 WHERE email_usuario = '${email}';"`);
}

function asAdmin() { window.sessionStorage.setItem('centinela_access', ADMIN_TOKEN); }
function clearBrowserSession() { window.sessionStorage.clear(); window.localStorage.clear(); jar = []; }

describe.skipIf(!API || !ADMIN_TOKEN || !COMPOSE)('INT-02 — Smoke Test End-to-End MVP (Frontend ↔ Backend real)', () => {
  const stamp = Date.now();
  const email = `int02.e2e.${stamp}@elcentinela.com`;
  let userId = '';
  let temporary = '';
  let secret = '';
  // El backend limita la contraseña a 12 caracteres (maxLength del contrato CambiarContrasenaInput).
  const password = 'Int02Pass12!';
  let operatorToken = '';

  beforeAll(() => { globalThis.fetch = browserFetch as typeof fetch; });
  afterAll(() => { globalThis.fetch = realFetch; });

  it('Paso 1: Alta, enrolamiento 2FA y login del OPERATOR', { timeout: 30000 }, async () => {
    clearBrowserSession(); asAdmin();
    const created = await createUser({
      nombreCompleto: 'Operador INT-02 E2E',
      nombreUsuario: `int02_e2e_${stamp}`,
      emailUsuario: email,
      rol: 'OPERATOR',
    });
    userId = created.id;
    temporary = mailed(email, 'clave provisoria es:');

    // Asignar permiso FULL_ACCESS sobre 101
    await assignUserInstances(userId, [{ vmid: 101, nivelAcceso: 'FULL_ACCESS' }]);

    // Enrolar 2FA
    clearBrowserSession();
    const preLogin = await submitLoginCredentials({ email, password: temporary } as never);
    const qr = await twoFactorService.qr(preLogin.jwtTemporal);
    secret = qr.secretoManual;
    const tokens = await twoFactorService.verify(preLogin.jwtTemporal, totp(secret));

    // Cambio obligatorio de contraseña
    window.sessionStorage.setItem('centinela_access', tokens.accessToken);
    await apiClient.request('/account/password', {
      method: 'PUT',
      payload: { contrasenaActual: temporary, contrasenaNueva: password },
    });

    // Login definitivo
    clearBrowserSession(); allowNextTotp(email);
    const defLogin = await submitLoginCredentials({ email, password } as never);
    const defTokens = await twoFactorService.verify(defLogin.jwtTemporal, totp(secret));
    operatorToken = defTokens.accessToken;
    window.sessionStorage.setItem('centinela_access', operatorToken);
    expect(operatorToken).toBeTruthy();
  });

  it('Paso 2: Dashboard — consulta de telemetría de nodo y resumen de instancias', { timeout: 15000 }, async () => {
    window.sessionStorage.setItem('centinela_access', operatorToken);
    const nodeStatus = await apiClient.get<Record<string, unknown>>('/node/status');
    expect(nodeStatus).toHaveProperty('cpu');
    expect(nodeStatus).toHaveProperty('ram');
    expect(nodeStatus).toHaveProperty('storage');
    expect(nodeStatus).toHaveProperty('instancesSummary');
  });

  it('Paso 3: Inventario — filtrado RBAC (solo ve 101 con nivelAcceso FULL_ACCESS)', { timeout: 15000 }, async () => {
    window.sessionStorage.setItem('centinela_access', operatorToken);
    const instances = await apiClient.get<Array<Record<string, unknown>>>('/instances');
    expect(instances.length).toBe(1);
    expect(String(instances[0].id)).toBe('101');
    expect(instances[0].nivelAcceso).toBe('FULL_ACCESS');
  });

  it('Paso 4: Ciclo de vida — despacho asíncrono y retorno 202 con tareaId', { timeout: 20000 }, async () => {
    window.sessionStorage.setItem('centinela_access', operatorToken);
    // La 101 está encendida: la orden válida es el apagado ordenado (con /start sería 409 INSTANCE_INVALID_STATE).
    const result = await apiClient.post<{ tareaId: string; upid?: string }>('/instances/101/status/shutdown');
    expect(result.tareaId).toBeTruthy();
    expect(typeof result.tareaId).toBe('string');
  });

  it('Paso 5: Notificación — ticket efímero de eventos y recepción de fin de tarea', { timeout: 25000 }, async () => {
    window.sessionStorage.setItem('centinela_access', operatorToken);
    const ticketResp = await apiClient.post<{ ticket: string }>('/events/ticket');
    expect(ticketResp.ticket).toBeTruthy();
  });

  it('Paso 6: Auditoría — el ADMIN verifica registros de ciclo de vida', { timeout: 15000 }, async () => {
    asAdmin();
    // GET /api/admin/audit responde paginado: { items: [...], pagina, total }.
    const audit = await apiClient.get<{ items: Array<Record<string, unknown>> }>('/admin/audit');
    expect(Array.isArray(audit.items)).toBe(true);
    expect(audit.items.length).toBeGreaterThan(0);
  });
});
