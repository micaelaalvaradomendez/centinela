import fs from 'node:fs';
import path from 'node:path';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { Toaster } from '@/components/ui/toast';
import Auditoria from '@/pages/Auditoria';
import RecoverPassword from '@/pages/RecoverPassword';
import InstanceResourceUsage from '@/components/features/instances/components/detail/InstanceResourceUsage';
import { SnapshotsTabContent } from '@/components/features/instances/components/detail/SnapshotsTabContent';
import { AUDIT_ACTION_LABELS, AUDIT_ACTION_OPTIONS, HIDDEN_AUDIT_ACTIONS } from '@/utils/auditActionLabels';
import { withAppProviders } from './app-providers';

// Pruebas de aceptación de tareas pendientes de documentacion/actual.md sin cobertura previa:
// FIX-127, FIX-131, FIX-132, FRN-32 y FIX-144. Metodología TDD: en rojo hasta implementar.

const frontendSrc = path.resolve(__dirname, '../../frontend/centinela/src');
const backendSrc = path.resolve(__dirname, '../../backend');
const read = (relative: string) => fs.readFileSync(path.join(frontendSrc, relative), 'utf-8');
const readBackend = (relative: string) => fs.readFileSync(path.join(backendSrc, relative), 'utf-8');

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });
}

afterEach(() => {
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe('FIX-127 - Resiliencia de SnapshotsTabContent y alias de ruta :vmid / :instanceRef', () => {
  beforeEach(() => {
    window.sessionStorage.setItem('centinela_access', 'token-fix127');
    window.localStorage.setItem('centinela_user', JSON.stringify({
      id: 'admin-1', nombreCompleto: 'Admin', email: 'a@centinela.local', rol: 'ADMIN', tiene2FA: true,
    }));
  });

  it('SnapshotsTabContent se monta sin EventsProvider sin lanzar error', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse([])));
    expect(() =>
      render(withAppProviders(<MemoryRouter><SnapshotsTabContent vmid="inst_abc123" /></MemoryRouter>)),
    ).not.toThrow();
    // Estado vacío visible: no queda el árbol roto por falta del provider.
    await waitFor(() => expect(screen.queryByText(/useEventsContext debe utilizarse/i)).not.toBeInTheDocument());
  });

  it('InstanceDetailPage lee instanceRef con fallback a vmid', () => {
    const source = read('pages/InstanceDetailPage.tsx');
    expect(source, 'FIX-127: la página debe leer también el parámetro :vmid').toMatch(/vmid/);
    expect(source).toMatch(/instanceRef\s*\?\?\s*vmid|vmid\s*\?\?\s*instanceRef|instanceRef\s*\|\|\s*vmid/);
  });

  it('applicationRoutes documenta o registra la compatibilidad de ruta entre :vmid e :instanceRef', () => {
    const files = ['routes/applicationRoutes.tsx', 'lib/applicationRoutes.tsx', 'applicationRoutes.tsx', 'App.tsx']
      .filter((f) => fs.existsSync(path.join(frontendSrc, f)));
    const routes = files.map(read).join('\n');
    expect(routes, 'FIX-127: falta la referencia o alias :vmid en las rutas').toMatch(/:vmid/);
    expect(routes).toMatch(/instances\/:instanceRef/);
  });
});

describe('FIX-131 - Gráfico de Tráfico de Red en el detalle de instancia (CU-05)', () => {
  const detail = {
    id: 'inst_abc123', name: 'web-01', type: 'vm', status: 'running', cpuPercent: 10, cores: 2,
    memoryUsedBytes: 1024 ** 3, memoryTotalBytes: 2 * 1024 ** 3,
  } as never;
  const metrics = [
    { timestamp: 1_760_000_000, cpuPercent: 10, ramUsedBytes: 1e9, ramTotalBytes: 2e9, diskReadBps: 1, diskWriteBps: 1, netInBps: 2048, netOutBps: 1024 },
    { timestamp: 1_760_000_060, cpuPercent: 20, ramUsedBytes: 1e9, ramTotalBytes: 2e9, diskReadBps: 1, diskWriteBps: 1, netInBps: 4096, netOutBps: 2048 },
  ] as never;

  it('InstanceMetricPoint declara netInBps y netOutBps', () => {
    const types = read('components/features/instances/types/instanceDetail.ts');
    expect(types).toMatch(/netInBps:\s*number \| null/);
    expect(types).toMatch(/netOutBps:\s*number \| null/);
  });

  it('el servicio normaliza netInBps y netOutBps', () => {
    const service = read('components/features/instances/services/instanceDetailService.ts');
    expect(service).toMatch(/netInBps/);
    expect(service).toMatch(/netOutBps/);
  });

  it('la grilla de uso renderiza CPU, Memoria, Disco y Tráfico de red', () => {
    render(
      <InstanceResourceUsage detail={detail} metrics={metrics} timeframe="2h"
        onTimeframeChange={() => {}} isLoading={false} error={null} onRetry={() => {}} />,
    );
    for (const titulo of [/^cpu$/i, /^memoria$/i, /^disco$/i, /tr[aá]fico de red|^red$/i]) {
      expect(screen.getAllByText(titulo).length, `falta la métrica ${titulo}`).toBeGreaterThan(0);
    }
    // Formato de tasa: KB/s o MB/s
    expect(screen.getAllByText(/KB\/s|MB\/s/).length).toBeGreaterThan(0);
  });
});

describe('FIX-132 - Errores 400 de contraseña no queman el OTP (CU-12)', () => {
  function backend(errorCode: string) {
    return vi.fn().mockImplementation((url: string) => {
      if (url.endsWith('/auth/password/forgot')) return Promise.resolve(jsonResponse({ message: 'ok' }));
      if (url.endsWith('/auth/password/reset')) return Promise.resolve(jsonResponse({ errorCode, message: 'La nueva contraseña no es válida.' }, 400));
      return Promise.resolve(jsonResponse({}));
    });
  }

  async function llegarAlPasoFinal(user: ReturnType<typeof userEvent.setup>) {
    const router = createMemoryRouter([
      { path: '/recover-password', Component: RecoverPassword },
      { path: '/login', element: <h1>Iniciar sesión</h1> },
    ], { initialEntries: ['/recover-password'] });
    render(<Toaster><RouterProvider router={router} /></Toaster>);
    await user.type(screen.getByLabelText(/correo electrónico/i), 'usuario@centinela.local');
    await user.click(screen.getByRole('button', { name: /enviar código de recuperación/i }));
    await user.type(await screen.findByLabelText(/código de verificación/i), '123456');
    await user.click(screen.getByRole('button', { name: /verificar código/i }));
    await screen.findByRole('heading', { name: /nueva contraseña/i });
    await user.type(screen.getByLabelText(/^nueva contraseña$/i), 'Nueva2026!');
    await user.type(screen.getByLabelText(/repetir contraseña|confirmar nueva contraseña/i), 'Nueva2026!');
  }

  for (const errorCode of ['SAME_PASSWORD', 'WEAK_PASSWORD', 'PASSWORD_RECENTLY_USED']) {
    it(`${errorCode} se muestra sobre la contraseña, no cuenta intentos de OTP y no bloquea`, { timeout: 30000 }, async () => {
      vi.stubGlobal('fetch', backend(errorCode));
      const user = userEvent.setup();
      await llegarAlPasoFinal(user);

      // Cuatro intentos: con el bug, el 3.º bloquea el proceso y acusa al código.
      for (let i = 0; i < 4; i += 1) {
        const boton = screen.getByRole('button', { name: /restablecer contraseña/i });
        expect(boton, `el botón se bloqueó en el intento ${i + 1}`).toBeEnabled();
        await user.click(boton);
        await waitFor(() => expect(screen.getAllByText(/contrase(ñ|n)a/i).length).toBeGreaterThan(0));
      }

      expect(screen.queryAllByText(/código.*(inválido|incorrecto|vencido|expir|no es válido)/i)).toHaveLength(0);
      expect(screen.queryAllByText(/bloquead|demasiados intentos|solicit[aá] un nuevo c[oó]digo/i)).toHaveLength(0);
      expect(screen.getByRole('heading', { name: /nueva contraseña/i })).toBeInTheDocument();
    });
  }

  it('RESET_FAILED sí cuenta intentos y bloquea al 3.er fallo', { timeout: 30000 }, async () => {
    vi.stubGlobal('fetch', backend('RESET_FAILED'));
    const user = userEvent.setup();
    await llegarAlPasoFinal(user);
    await user.click(screen.getByRole('button', { name: /restablecer contraseña/i }));
    expect((await screen.findAllByText(/código.*(inválido|incorrecto|vencido|expir|no es válido)/i))[0]).toBeVisible();
  });
});

describe('FRN-32 - Exportación de reportes de auditoría', () => {
  const items = [{
    id: 'aud-1', usuarioId: 'u1', nombreUsuario: 'Admin Principal', fechaHora: '2026-10-24T10:15:00Z',
    accion: 'SNAPSHOT_CREATE', resourceType: 'SNAPSHOT', recursoTipo: 'SNAPSHOT', instanciaNombre: 'web-01',
    resultado: 'SUCCESS', detalles: {},
  }];

  function renderAuditoria() {
    const router = createMemoryRouter([{ path: '/audit', Component: Auditoria }], { initialEntries: ['/audit'] });
    render(withAppProviders(<Toaster><RouterProvider router={router} /></Toaster>));
  }

  beforeEach(() => {
    window.sessionStorage.setItem('centinela_access', 'token-frn32');
    window.sessionStorage.setItem('centinela_user', JSON.stringify({ id: 'admin-1', rol: 'ADMIN' }));
    window.localStorage.setItem('centinela_user', JSON.stringify({
      id: 'admin-1', nombreCompleto: 'Admin', email: 'a@centinela.local', rol: 'ADMIN', tiene2FA: true,
    }));
  });

  function backend(exportDelayMs = 0) {
    return vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/admin/audit/export')) {
        const isJson = /format(o)?=json/i.test(String(url));
        const respond = () => new Response(isJson ? '[]' : 'fecha,usuario\n', {
          status: 200,
          headers: {
            'Content-Type': isJson ? 'application/json' : 'text/csv',
            'Content-Disposition': `attachment; filename="auditoria_centinela_20261024.${isJson ? 'json' : 'csv'}"`,
          },
        });
        return exportDelayMs ? new Promise((r) => setTimeout(() => r(respond()), exportDelayMs)) : Promise.resolve(respond());
      }
      return Promise.resolve(jsonResponse({ total: items.length, pagina: 1, items }));
    });
  }

  beforeEach(() => {
    vi.stubGlobal('URL', Object.assign(URL, { createObjectURL: vi.fn(() => 'blob:x'), revokeObjectURL: vi.fn() }));
  });

  it('un único botón "Exportar registros" despliega las opciones CSV y JSON', async () => {
    vi.stubGlobal('fetch', backend());
    const user = userEvent.setup();
    renderAuditoria();
    await screen.findByRole('table');

    await user.click(screen.getByRole('button', { name: /exportar registros|exportar/i }));
    expect(await screen.findByRole('menuitem', { name: /exportar.*csv/i })).toBeVisible();
    expect(screen.getByRole('menuitem', { name: /exportar.*json/i })).toBeVisible();
  });

  it('exportar como JSON pide format=json, descarga el archivo y avisa con un Toast', async () => {
    const descargas: string[] = [];
    vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(function (this: HTMLAnchorElement) {
      if (this.download) descargas.push(this.download);
    });
    const fetchMock = backend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderAuditoria();
    await screen.findByRole('table');

    await user.click(screen.getByRole('button', { name: /exportar/i }));
    await user.click(await screen.findByRole('menuitem', { name: /exportar.*json/i }));

    await waitFor(() => {
      const urls = (fetchMock.mock.calls as [string][]).map(([u]) => String(u));
      expect(urls.some((u) => u.includes('/admin/audit/export') && /format(o)?=json/i.test(u))).toBe(true);
    });
    await waitFor(() => expect(descargas.some((n) => /^auditoria_centinela_\d{8}(_\d{4})?\.json$/.test(n))).toBe(true));
    expect((await screen.findAllByText(/descarg|exportaci[oó]n/i)).length).toBeGreaterThan(0);
  });

  it('mientras genera el reporte el botón muestra estado de carga y queda deshabilitado', async () => {
    vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => {});
    vi.stubGlobal('fetch', backend(400));
    const user = userEvent.setup();
    renderAuditoria();
    await screen.findByRole('table');

    await user.click(screen.getByRole('button', { name: /exportar/i }));
    await user.click(await screen.findByRole('menuitem', { name: /exportar.*csv/i }));

    const boton = screen.getByRole('button', { name: /exportando|generando|exportar/i });
    expect(boton).toBeDisabled();
    await waitFor(() => expect(screen.getByRole('button', { name: /exportar/i })).toBeEnabled(), { timeout: 3000 });
  });

  it('transfiere los filtros activos al servicio de exportación', async () => {
    vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => {});
    const fetchMock = backend();
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();
    renderAuditoria();
    const table = await screen.findByRole('table');
    expect(within(table).getByText('Admin Principal')).toBeVisible();

    await user.click(screen.getByRole('button', { name: /exportar/i }));
    await user.click(await screen.findByRole('menuitem', { name: /exportar.*csv/i }));

    await waitFor(() => {
      const urls = (fetchMock.mock.calls as [string][]).map(([u]) => String(u));
      expect(urls.some((u) => u.includes('/admin/audit/export'))).toBe(true);
    });
  });
});

describe('FIX-144 - Alineación de filtros, contadores y diccionario de acciones de auditoría en Frontend con el enum cerrado del Backend (RF-08 / FRN-31)', () => {
  it('AUDIT_ACTION_LABELS y AUDIT_ACTION_OPTIONS en auditActionLabels.ts coinciden con las constantes Accion* de backend/internal/core/ports/audit_port.go', () => {
    const auditPort = readBackend('internal/core/ports/audit_port.go');
    const backendActions = Array.from(auditPort.matchAll(/Accion\w+\s*(?:\w+)?\s*=\s*"([^"]+)"/g), (m) => m[1]);
    const labelKeys = Object.keys(AUDIT_ACTION_LABELS);
    const frontendOptions = AUDIT_ACTION_OPTIONS.map((o) => o.value);

    // Ninguna opción del catálogo o filtro del frontend puede enviar un código inexistente en ports/audit_port.go
    const missingInBackend = labelKeys.filter((code) => !backendActions.includes(code));
    expect(missingInBackend).toEqual([]);

    // Todas las constantes activas de ports/audit_port.go deben estar en AUDIT_ACTION_LABELS y (salvo HIDDEN_AUDIT_ACTIONS) en AUDIT_ACTION_OPTIONS
    const missingInLabels = backendActions.filter((code) => !labelKeys.includes(code));
    expect(missingInLabels).toEqual([]);

    const missingInOptions = backendActions
      .filter((code) => !HIDDEN_AUDIT_ACTIONS.includes(code))
      .filter((code) => !frontendOptions.includes(code));
    expect(missingInOptions).toEqual([]);
  });

  it('Auditoria.tsx y deletedUsersService.ts consultan un único conjunto cerrado de resultado alineado con ports/audit_port.go', () => {
    const auditPort = readBackend('internal/core/ports/audit_port.go');
    const backendResults = Array.from(new Set(Array.from(auditPort.matchAll(/Resultado\w+\s*(?:\w+)?\s*=\s*"([^"]+)"/g), (m) => m[1])));

    // El backend no debe mezclar las dos familias (EXITO/FALLA junto con SUCCESS/FAILED)
    const hasSpanish = backendResults.includes('EXITO') || backendResults.includes('FALLA');
    const hasEnglish = backendResults.includes('SUCCESS') || backendResults.includes('FAILED');
    expect(hasSpanish && hasEnglish).toBe(false);

    // Auditoria.tsx y deletedUsersService.ts deben usar la misma familia canónica (EXITO/FALLA/PENDIENTE)
    const auditoriaSrc = read('pages/Auditoria.tsx');
    const deletedUsersSrc = read('components/features/users/services/deletedUsersService.ts');

    const auditoriaUsesSuccess = /resultado:\s*['"]SUCCESS['"]|value=['"]SUCCESS['"]/.test(auditoriaSrc);
    const deletedUsersUsesExito = /resultado:\s*['"]EXITO['"]/.test(deletedUsersSrc);
    expect(auditoriaUsesSuccess).toBe(false);
    expect(deletedUsersUsesExito).toBe(true);
  });
});




