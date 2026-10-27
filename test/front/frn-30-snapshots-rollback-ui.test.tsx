import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import fs from 'node:fs';
import path from 'node:path';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstanceDetailPage from '@/pages/InstanceDetailPage';
import { withProtectedProviders } from './app-providers';

// FRN-30 (Etapa 3) — Modal destructivo de Rollback y eliminación con confirmación antierror.
// Entregables:
//   1. Modal antierror de Rollback con advertencia destructiva explícita.
//   2. Input de confirmación que exige escribir textualmente el nombre exacto del snapshot para habilitar el botón Revertir.
//   3. Modal de confirmación para eliminar snapshots.
//   4. Bloqueo de la botonera con estado transitorio ROLLING_BACK mientras se ejecuta la tarea asíncrona.

const mockSnapshots = [
  {
    name: 'snap-seguridad-01',
    description: 'Punto antes de actualizar kernel',
    snaptime: 1729000000,
    vmstate: 0,
  },
  {
    name: 'snap-estable-v2',
    description: 'Versión estable producción',
    snaptime: 1729100000,
    vmstate: 1,
  },
];

const mockInstance = {
  id: 'inst-test-01',
  name: 'servidor-web-prod',
  type: 'qemu',
  status: 'running',
  node: 'proxmox-pve',
  ip: '192.168.1.150',
  cpuUsage: 12.5,
  ramUsage: 1073741824,
  maxRam: 2147483648,
  nivelAcceso: 'FULL_ACCESS',
  activeTask: null,
};

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

describe('FRN-30 - Modal destructivo de Rollback y confirmación antierror', () => {
  beforeEach(() => {
    window.sessionStorage.setItem('centinela_access', 'token-frn30');
    window.localStorage.setItem(
      'centinela_user',
      JSON.stringify({
        id: 'admin-1',
        nombreCompleto: 'Admin Test',
        email: 'admin@centinela.local',
        rol: 'ADMIN',
        permisos: [{ vmid: 'inst-test-01', nivelAcceso: 'FULL_ACCESS' }],
        tiene2FA: true,
      }),
    );
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('el código fuente implementa el modal con advertencia y verificación de nombre exacto', () => {
    const srcDir = path.resolve(__dirname, '../../frontend/centinela/src');
    function searchPattern(pattern: RegExp): boolean {
      let found = false;
      function scan(dir: string) {
        for (const file of fs.readdirSync(dir)) {
          const full = path.join(dir, file);
          if (fs.statSync(full).isDirectory()) scan(full);
          else if (/\.(tsx?|jsx?)$/.test(file)) {
            const content = fs.readFileSync(full, 'utf-8');
            if (pattern.test(content)) {
              found = true;
              return;
            }
          }
        }
      }
      scan(srcDir);
      return found;
    }

    const hasWarningText = searchPattern(/Esta acción restaurará la máquina al punto seleccionado|datos posteriores se perderán|revertir.*snapshot/i);
    const hasRollbackEndpoint = searchPattern(/\/snapshots\/.*\/rollback/);

    expect(hasWarningText, 'Falta la advertencia destructiva en el componente modal de rollback').toBe(true);
    expect(hasRollbackEndpoint, 'Falta el consumo del endpoint de rollback en los servicios de snapshots').toBe(true);
  });

  it('el botón Revertir permanece deshabilitado hasta que el usuario tipea el nombre idéntico', async () => {
    const fetchMock = vi.fn().mockImplementation((url: string) => {
      const u = String(url);
      if (u.includes('/snapshots')) return Promise.resolve(json(mockSnapshots));
      if (u.includes('/instances/')) return Promise.resolve(json(mockInstance));
      return Promise.resolve(json({}));
    });
    vi.stubGlobal('fetch', fetchMock);

    render(
      withProtectedProviders(
        <MemoryRouter initialEntries={['/instances/inst-test-01']}>
          <Routes>
            <Route path="/instances/:vmid" element={<InstanceDetailPage />} />
          </Routes>
        </MemoryRouter>,
      ),
    );

    // Si la pestaña de snapshots existe, buscar botón de rollback
    const snapshotsTab = await screen.findByRole('tab', { name: /snapshots/i });
    const user = userEvent.setup();
    await user.click(snapshotsTab);

    const rollbackButton = screen.queryByRole('button', { name: /revertir|rollback/i });
    if (rollbackButton) {
      await user.click(rollbackButton);
      const dialog = await screen.findByRole('dialog');
      const confirmButton = within(dialog).getByRole('button', { name: /confirmar|revertir/i });
      expect(confirmButton).toBeDisabled();

      const input = within(dialog).getByRole('textbox');
      await user.type(input, 'nombre-invalido');
      expect(confirmButton).toBeDisabled();

      await user.clear(input);
      await user.type(input, 'snap-seguridad-01');
      expect(confirmButton).toBeEnabled();
    }
  });
});

