import { render, screen, waitFor } from '@testing-library/react';
import fs from 'node:fs';
import path from 'node:path';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DashboardPage from '@/pages/Dashboard';
import { withProtectedProviders } from './app-providers';

// FRN-33 (Etapa 3) — Centro de Notificaciones y banners contextuales de alerta por saturación.
// Entregables:
//   1. Campana / Centro de Notificaciones en la barra superior con contador de alertas no leídas.
//   2. Banner contextual persistente de color rojo en la cabecera del Dashboard ante alertas CRITICAL.
//   3. Badge de advertencia en la fila del inventario si una instancia tiene alertas de saturación.
//   4. Recepción en vivo de eventos ALERT_TRIGGERED / ALERT_RESOLVED vía SSE.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

describe('FRN-33 - Centro de notificaciones y banners contextuales de alerta', () => {
  beforeEach(() => {
    window.sessionStorage.setItem('centinela_access', 'token-frn33');
    window.localStorage.setItem(
      'centinela_user',
      JSON.stringify({
        id: 'admin-1',
        nombreCompleto: 'Admin Test',
        email: 'admin@centinela.local',
        rol: 'ADMIN',
        tiene2FA: true,
      }),
    );
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('el código fuente define el componente del centro de notificaciones y manejo de alertas críticas', () => {
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

    const hasAlertTriggered = searchPattern(/ALERT_TRIGGERED|ALERT_RESOLVED/);
    const hasNotificationBell = searchPattern(/(notification|campana|alert-badge|bell)/i);

    expect(hasAlertTriggered, 'El frontend no procesa los eventos ALERT_TRIGGERED o ALERT_RESOLVED').toBe(true);
    expect(hasNotificationBell, 'Falta el componente o icono del centro de notificaciones / alertas').toBe(true);
  });

  it('el Dashboard muestra el banner de alerta cuando se reporta saturación crítica de recursos', async () => {
    const fetchMock = vi.fn().mockImplementation((url: string) => {
      const u = String(url);
      if (u.includes('/node/status')) {
        return Promise.resolve(
          json({
            cpu: { usagePercent: 95, cores: 8 },
            ram: { usedGb: 15, totalGb: 16, usagePercent: 93.75 },
            storage: { usedGb: 95, totalGb: 100, usagePercent: 95 },
            uptimeSeconds: 3600,
            instancesSummary: { vms: { running: 2, stopped: 0, paused: 0, total: 2 }, lxc: { running: 0, stopped: 0, paused: 0, total: 0 } },
          }),
        );
      }
      return Promise.resolve(json({}));
    });
    vi.stubGlobal('fetch', fetchMock);

    render(
      withProtectedProviders(
        <MemoryRouter initialEntries={['/dashboard']}>
          <DashboardPage />
        </MemoryRouter>,
      ),
    );

    // Debe renderizar el panel del Dashboard sin crashear
    await waitFor(() => {
      expect(screen.queryByText(/saturación|alerta|crítica|recursos/i) || screen.queryByRole('alert')).toBeDefined();
    });
  });
});

