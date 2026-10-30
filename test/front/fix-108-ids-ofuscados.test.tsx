import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import fs from 'node:fs';
import path from 'node:path';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { requestInstancePowerAction } from '@/components/features/instances/services/instanceService';
import { commands, findDialog, renderInstances, requireAction, rowOf } from './instances-helpers';

// FIX-108: el backend (FIX-107) expone la instancia con un identificador opaco no secuencial
// (p. ej. "inst_k8Pq2v") en lugar del VMID entero.
// Criterio de éxito: la tabla renderiza el listado con ids ofuscados y las acciones operan con ese
// identificador en la URL, sin conversiones numéricas destructivas y sin aceptar caracteres hostiles.

const OPAQUE_ID = 'inst_k8Pq2v';

const opaqueInventory = [
  {
    id: OPAQUE_ID, name: 'api-gateway', type: 'lxc', node: 'nodo-principal', status: 'stopped', ip: null,
    cpuUsage: null, ramUsage: null, maxRam: null, nivelAcceso: 'FULL_ACCESS', activeTask: null,
  },
];

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

function stubOpaqueBackend() {
  const fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
    const method = (init.method ?? 'GET').toUpperCase();
    if (method === 'GET' && /\/instances(\?|$)/.test(String(url))) return Promise.resolve(json(opaqueInventory));
    if (method === 'GET') return Promise.resolve(json({}));
    return Promise.resolve(json({ tareaId: 't-opaco-1' }, 202));
  });
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FIX-108 - identificadores de instancia ofuscados en el frontend', () => {
  it('la tabla renderiza una instancia con id opaco y "Iniciar" invoca /instances/<id opaco>/start', async () => {
    const fetchMock = stubOpaqueBackend();
    renderInstances({ rol: 'ADMIN' });

    const row = await rowOf('api-gateway');
    const user = userEvent.setup();
    await user.click(await requireAction(row, 'start'));
    const dialog = await findDialog();
    await user.click(within(dialog).getByRole('button', { name: /confirmar|iniciar|encender|aceptar/i }));

    await waitFor(() => expect(commands(fetchMock).length).toBe(1));
    const [url, init] = commands(fetchMock)[0];
    expect(String(url)).toMatch(new RegExp(`/instances/${OPAQUE_ID}/start$`));
    expect(String(url)).not.toMatch(/NaN|undefined|null/);
    expect((init?.method ?? '').toUpperCase()).toBe('POST');
  });

  it('el detalle de la fila enlaza con /instances/<id opaco> sin forzar un número', async () => {
    stubOpaqueBackend();
    renderInstances({ rol: 'ADMIN' });
    const row = await rowOf('api-gateway');
    const links = within(row).queryAllByRole('link');
    expect(links.length, 'la fila no ofrece un enlace al detalle').toBeGreaterThan(0);
    expect(links.map((link) => link.getAttribute('href')).join(' ')).toContain(`/instances/${OPAQUE_ID}`);
  });

  it('requestInstancePowerAction acepta el id opaco y arma la URL exacta', async () => {
    window.sessionStorage.setItem('centinela_access', 'access-token-fix108');
    const fetchMock = stubOpaqueBackend();
    await requestInstancePowerAction(OPAQUE_ID as never, 'shutdown');
    const [url] = commands(fetchMock)[0];
    expect(String(url)).toMatch(new RegExp(`/instances/${OPAQUE_ID}/status/shutdown$`));
  });

  it.each([
    ['path traversal', '../admin/users'],
    ['espacios', 'inst abc'],
    ['query string', 'inst_1?force=1'],
    ['separador de ruta', 'inst/1'],
    ['inyección con punto y coma', 'inst_1;rm'],
    ['cadena vacía', ''],
    ['salto de línea', 'inst_1\nX'],
  ])('rechaza identificadores con %s sin llegar a enviar la petición', async (_name, hostile) => {
    window.sessionStorage.setItem('centinela_access', 'access-token-fix108');
    const fetchMock = stubOpaqueBackend();
    await expect(requestInstancePowerAction(hostile as never, 'start')).rejects.toThrow();
    expect(commands(fetchMock)).toHaveLength(0);
  });

  it('InstanceDetailPage y useInstances no convierten el id con Number()/parseInt()', () => {
    const base = path.resolve(__dirname, '../../frontend/centinela/src');
    for (const file of ['pages/InstanceDetailPage.tsx', 'components/features/instances/hooks/useInstances.ts']) {
      const source = fs.readFileSync(path.join(base, file), 'utf-8');
      expect(source, `${file} convierte el identificador de instancia a número`).not.toMatch(/parseInt\(|Number\((vmid|id|instanceId|params)/);
    }
  });
});

