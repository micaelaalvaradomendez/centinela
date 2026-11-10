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

class FakeSse {
  listeners: Record<string, ((event: Event) => void)[]> = {};
  constructor(public url: string) {
    setTimeout(() => {
      this.listeners.open?.forEach((l) => l(new Event('open')));
    }, 0);
  }
  addEventListener(type: string, listener: (event: Event) => void) {
    (this.listeners[type] ??= []).push(listener);
  }
  removeEventListener() {}
  close() {}
  emitMessage(data: unknown) {
    const event = new MessageEvent('message', { data: typeof data === 'string' ? data : JSON.stringify(data) });
    this.listeners.message?.forEach((listener) => listener(event));
  }
}

let latestSse: FakeSse | null = null;

describe('FRN-33 - Centro de notificaciones y banners contextuales de alerta', () => {
  beforeEach(() => {
    latestSse = null;
    class CapturingSse extends FakeSse {
      constructor(url: string) {
        super(url);
        latestSse = this;
      }
    }
    vi.stubGlobal('EventSource', CapturingSse);
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

    const hasAlertTriggered = searchPattern(/ALERT_TRIGGERED/) && searchPattern(/ALERT_RESOLVED/);
    // Un componente real: importa el icono Bell de lucide-react y es un archivo de componente.
    const hasNotificationBell = searchPattern(/import\s*\{[^}]*\bBell\b[^}]*\}\s*from\s*'lucide-react'/);
    const hasNotificationComponent = fs
      .readdirSync(path.join(srcDir, 'components'), { recursive: true })
      .some((f) => /(Notification|Alert).*\.tsx$/i.test(String(f)));

    expect(hasAlertTriggered, 'El frontend no procesa ALERT_TRIGGERED y ALERT_RESOLVED').toBe(true);
    expect(hasNotificationBell, 'Falta la campana (icono Bell) del centro de notificaciones').toBe(true);
    expect(hasNotificationComponent, 'Falta un componente Notification*/Alert*.tsx en src/components').toBe(true);
  });

  it('el Dashboard muestra el banner de alerta cuando se reporta saturación crítica de recursos', async () => {
    const fetchMock = vi.fn().mockImplementation((url: string) => {
      const u = String(url);
      if (u.includes('/events/ticket') || u.includes('/ws-ticket')) {
        return Promise.resolve(json({ ticket: 'test-ticket-frn33' }));
      }
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
      if (u.includes('/instances')) {
        return Promise.resolve(json([]));
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

    await waitFor(() => {
      expect(latestSse).not.toBeNull();
    });

    latestSse?.emitMessage({
      id: 'alert-crit-1',
      tipo: 'ALERT_TRIGGERED',
      severidad: 'CRITICAL',
      recursoTipo: 'HOST',
      recursoId: null,
      mensaje: 'Saturación crítica del procesador: 95%',
      fechaHora: new Date().toISOString(),
      detalles: { alertId: 'alert-crit-1', metrica: 'CPU' },
    });

    // Banner rojo persistente (role="alert") con el texto de saturación crítica.
    await waitFor(() => {
      const banners = screen.queryAllByRole('alert').filter(
        (el) => /saturaci[oó]n|cr[ií]tic/i.test(el.textContent ?? '') && /red|rose/.test(el.className),
      );
      expect(banners.length, 'FRN-33: no hay un banner rojo role="alert" de saturación crítica').toBeGreaterThan(0);
    });
  });
});

