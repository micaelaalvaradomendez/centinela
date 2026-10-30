import fs from 'node:fs';
import path from 'node:path';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import {
  createInstance,
  deleteInstance,
  requestInstancePowerAction,
} from '@/components/features/instances/services/instanceService';

// FIX-104B / FIX-105: el backend (commit bd22c3d, BAC-29) dejó de exponer `upid` en el 202 Accepted
// y devuelve solo el identificador opaco { tareaId }. El servicio de instancias debe seguir la
// orden con `tareaId` y no exigir `upid`.
// Criterio de éxito: encender, apagar, reiniciar, borrar y crear funcionan sin `upid`.

const json = (body: unknown, status = 202) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

let fetchMock: ReturnType<typeof vi.fn>;

function respondWith(body: unknown) {
  fetchMock = vi.fn().mockImplementation(() => Promise.resolve(json(body)));
  vi.stubGlobal('fetch', fetchMock);
}

beforeEach(() => {
  window.sessionStorage.setItem('centinela_access', 'access-token-tareaid');
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FIX-104B / FIX-105 - contrato seguro tareaId en instanceService', () => {
  it.each(['start', 'stop', 'shutdown', 'reboot'] as const)(
    'requestInstancePowerAction(%s) acepta un 202 con solo { tareaId }',
    async (action) => {
      respondWith({ tareaId: 'tarea-uuid-1' });
      const result = await requestInstancePowerAction(101, action);
      expect(result.tareaId).toBe('tarea-uuid-1');
      expect((result as { upid?: string }).upid).toBeUndefined();
    },
  );

  it.each([
    ['sin tareaId', {}],
    ['con tareaId vacío', { tareaId: '   ' }],
    ['con tareaId no textual', { tareaId: 42 }],
    ['con solo el upid heredado (sin tareaId)', { upid: 'UPID:pve:0001:start:101:' }],
  ])('requestInstancePowerAction rechaza un 202 %s', async (_name, body) => {
    respondWith(body);
    await expect(requestInstancePowerAction(101, 'start')).rejects.toThrow();
  });

  it('deleteInstance acepta { tareaId } sin upid y rechaza el 202 sin tareaId', async () => {
    respondWith({ tareaId: 'tarea-uuid-2' });
    await expect(deleteInstance(102)).resolves.toMatchObject({ tareaId: 'tarea-uuid-2' });

    respondWith({ upid: 'UPID:pve:0001:vzdestroy:102:' });
    await expect(deleteInstance(102)).rejects.toThrow();
  });

  it('createInstance acepta { vmid, tareaId } sin upid y rechaza si falta tareaId', async () => {
    const payload = {
      cores: 1, descripcion: '', discoGb: 4, nodo: 'pve', nombre: 'vm-nueva', password: '', ramMb: 512,
      red: { bridge: 'vmbr0', firewall: false, gateway: '', ip: '', mascara: '', modo: 'dhcp', vlan: 0 },
      storageId: 'local', templateOso: 'local:iso/debian-12.iso', tipo: 'vm',
    } as never;

    respondWith({ vmid: 120, tareaId: 'tarea-uuid-3' });
    await expect(createInstance(payload)).resolves.toMatchObject({ vmid: 120, tareaId: 'tarea-uuid-3' });

    respondWith({ vmid: 120 });
    await expect(createInstance(payload)).rejects.toThrow();
  });

  it('los tipos públicos no declaran `upid` como propiedad obligatoria', () => {
    const typesDir = path.resolve(__dirname, '../../frontend/centinela/src/components/features/instances/types');
    for (const file of ['instance.ts', 'createInstance.ts']) {
      const source = fs.readFileSync(path.join(typesDir, file), 'utf-8');
      expect(source, `${file} declara \`upid: string\` obligatorio en el contrato público`).not.toMatch(/\bupid:\s*string/);
    }
  });

  it('el servicio ya no valida response.upid ni creation.upid', () => {
    const source = fs.readFileSync(
      path.resolve(__dirname, '../../frontend/centinela/src/components/features/instances/services/instanceService.ts'),
      'utf-8',
    );
    expect(source).not.toMatch(/response\.upid|creation\.upid/);
  });
});

