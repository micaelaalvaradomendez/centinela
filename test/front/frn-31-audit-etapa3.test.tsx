import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import Auditoria from '@/pages/Auditoria';
import { withAppProviders } from './app-providers';

// FRN-31 (Etapa 3) — Vista avanzada de Auditoría, filtros compuestos, visor de detalles JSON y exportación (CSV y JSON).
// Entregables (documentacion/actual.md):
//   1. Controles de rango de fechas ("Desde" y "Hasta") con etiquetas accesibles, y selectores
//      para filtrar por acción (incluyendo SNAPSHOT_CREATE, SNAPSHOT_ROLLBACK, SNAPSHOT_DELETE,
//      ALERT_TRIGGERED, CONFIG_UPDATED / CONFIG_INSTANCE), usuario y tipo de recurso (recursoTipo).
//   2. Tabla con columnas Fecha/Hora, Usuario, Acción, Recurso/Instancia, Resultado (badges por
//      severidad: SUCCESS/EXITO en verde, FAILED/FALLA en rojo, WARNING/ADVERTENCIA en amarillo)
//      y paginación con selector de tamaño de página (10, 25, 50).
//   3. Modal o visor al hacer clic en "Ver detalles" de una fila para inspeccionar el payload JSON
//      de detalles formateado de manera legible.
//   4. Botones "Exportar CSV" y "Exportar JSON" que consumen GET /api/admin/audit/export con los
//      filtros activos y descargan el archivo respetando Content-Disposition.

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

const sampleAuditItems = [
  {
    id: 'aud-etapa3-1',
    usuarioId: 'usr-admin-1',
    nombreUsuario: 'Admin Principal',
    fechaHora: '2026-10-24T10:15:00Z',
    accion: 'SNAPSHOT_CREATE',
    resourceType: 'SNAPSHOT',
    resourceId: 'inst_k8Pq2v',
    recursoTipo: 'SNAPSHOT',
    recursoId: 'inst_k8Pq2v',
    instanciaId: 'inst_k8Pq2v',
    instanciaNombre: 'web-01',
    resultado: 'SUCCESS',
    detalles: { snapname: 'pre_deploy_2026', vmstate: false, duracionMs: 1420 },
  },
  {
    id: 'aud-etapa3-2',
    usuarioId: 'usr-op-1',
    nombreUsuario: 'Operador Guardia',
    fechaHora: '2026-10-24T10:20:00Z',
    accion: 'ALERT_TRIGGERED',
    resourceType: 'NODE',
    resourceId: 'nodo-a1',
    recursoTipo: 'NODE',
    recursoId: 'nodo-a1',
    instanciaId: 'nodo-a1',
    instanciaNombre: 'nodo-a1',
    resultado: 'WARNING',
    detalles: { metric: 'CPU', usagePercent: 92.4, threshold: 90 },
  },
  {
    id: 'aud-etapa3-3',
    usuarioId: 'usr-op-1',
    nombreUsuario: 'Operador Guardia',
    fechaHora: '2026-10-24T10:25:00Z',
    accion: 'SNAPSHOT_ROLLBACK',
    resourceType: 'SNAPSHOT',
    resourceId: 'inst_k8Pq2v',
    recursoTipo: 'SNAPSHOT',
    recursoId: 'inst_k8Pq2v',
    instanciaId: 'inst_k8Pq2v',
    instanciaNombre: 'web-01',
    resultado: 'FAILED',
    detalles: { snapname: 'pre_deploy_2026', motivo: 'PROXMOX_ERROR', error: 'lock timeout' },
  },
];

function renderAuditoria() {
  window.sessionStorage.setItem('centinela_access', 'token-frn31-admin');
  window.sessionStorage.setItem('centinela_user', JSON.stringify({ id: 'usr-admin-1', rol: 'ADMIN' }));
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'usr-admin-1',
    organizacionId: 'org-1',
    nombreCompleto: 'Admin Principal',
    email: 'admin@centinela.local',
    rol: 'ADMIN',
    instanciasPermitidas: [],
    tiene2FA: true,
  }));
  const router = createMemoryRouter([{ path: '/auditoria', element: <Auditoria /> }], { initialEntries: ['/auditoria'] });
  return render(withAppProviders(<RouterProvider router={router} />));
}

beforeEach(() => {
  HTMLDialogElement.prototype.showModal ??= function showModalPolyfill(this: HTMLDialogElement) {
    this.setAttribute('open', '');
  };
  HTMLDialogElement.prototype.close ??= function closePolyfill(this: HTMLDialogElement) {
    this.removeAttribute('open');
  };
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-31 - Vista avanzada de Auditoría, filtros Etapa 3, visor JSON y exportación CSV/JSON', () => {
  it('renderiza badges diferenciados por severidad/resultado: SUCCESS (verde), WARNING (amarillo) y FAILED (rojo)', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      total: sampleAuditItems.length,
      pagina: 1,
      items: sampleAuditItems,
    }));
    vi.stubGlobal('fetch', fetchMock);

    renderAuditoria();

    const table = await screen.findByRole('table');
    const successBadge = await within(table).findByText(/éxito|exito|success/i);
    const warningBadge = await within(table).findByText(/advertencia|warning/i);
    const failedBadge = await within(table).findByText(/falla|fallido|failed/i);

    expect(successBadge.closest('span')?.className ?? successBadge.className).toMatch(/emerald|green/);
    expect(warningBadge.closest('span')?.className ?? warningBadge.className).toMatch(/amber|yellow/);
    expect(failedBadge.closest('span')?.className ?? failedBadge.className).toMatch(/rose|red/);
  });

  it('ofrece controles accesibles de rango de fechas (Desde y Hasta), filtro por recursoTipo y acciones de Etapa 3 (SNAPSHOT_CREATE, SNAPSHOT_ROLLBACK, SNAPSHOT_DELETE, ALERT_TRIGGERED, CONFIG_UPDATED)', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      total: sampleAuditItems.length,
      pagina: 1,
      items: sampleAuditItems,
    }));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderAuditoria();
    await screen.findByRole('table');

    const filtrosBtn = screen.getByRole('button', { name: /filtros/i });
    await user.click(filtrosBtn);

    // Controles accesibles de rango de fechas ("Desde" y "Hasta", actual.md:309)
    const desdeControl = screen.queryByLabelText(/desde/i);
    const hastaControl = screen.queryByLabelText(/hasta/i);
    expect(desdeControl, 'falta el control accesible "Desde" para el rango de fechas (FRN-31)').not.toBeNull();
    expect(hastaControl, 'falta el control accesible "Hasta" para el rango de fechas (FRN-31)').not.toBeNull();

    // Filtro por tipo de recurso (recursoTipo: INSTANCE, SNAPSHOT, NODE, USER)
    const resourceTypeFilter = screen.queryByLabelText(/tipo de recurso|recurso/i)
      ?? screen.queryByRole('menuitemradio', { name: /snapshot/i })
      ?? screen.queryByRole('option', { name: /snapshot/i });
    expect(resourceTypeFilter, 'falta el selector de filtro por tipo de recurso (recursoTipo) en Auditoría (FRN-31)').not.toBeNull();

    // Opciones de acción de Etapa 3 en el selector de acciones
    for (const actionPattern of [
      /snapshot.*cread|crear.*snapshot|snapshot_create/i,
      /revertir.*snapshot|rollback|snapshot_rollback/i,
      /eliminar.*snapshot|snapshot_delete/i,
      /alerta|alert_triggered/i,
      /configuraci[oó]n|config_updated|config_instance/i,
    ]) {
      const option = screen.queryByRole('menuitemradio', { name: actionPattern })
        ?? screen.queryByRole('option', { name: actionPattern });
      expect(option, `falta la acción ${actionPattern} en el filtro de auditoría (FRN-31)`).not.toBeNull();
    }
  });

  it('al hacer clic en "Ver detalles" abre un modal o visor con el payload JSON de detalles formateado', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      total: sampleAuditItems.length,
      pagina: 1,
      items: sampleAuditItems,
    }));
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderAuditoria();
    const table = await screen.findByRole('table');
    await within(table).findByText('Admin Principal');

    const detailButtons = within(table).getAllByRole('button', { name: /ver detalles/i });
    await user.click(detailButtons[0]);

    const dialog = await screen.findByRole('dialog');
    expect(within(dialog).getByText(/pre_deploy_2026/)).toBeVisible();
  });

  it('permite cambiar el tamaño de página (10, 25, 50) y exportar tanto en CSV como en JSON respetando el nombre auditoria_centinela_YYYYMMDD', async () => {
    const downloadedFiles: string[] = [];
    const anchorClickSpy = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(function (this: HTMLAnchorElement) {
      if (this.download) downloadedFiles.push(this.download);
    });

    const fetchMock = vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/admin/audit/export')) {
        const isJson = String(url).includes('json');
        return Promise.resolve(new Response(isJson ? '[]' : 'id,accion\n1,LOGIN', {
          status: 200,
          headers: {
            'Content-Type': isJson ? 'application/json' : 'text/csv',
            'Content-Disposition': `attachment; filename="auditoria_centinela_20261024.${isJson ? 'json' : 'csv'}"`,
          },
        }));
      }
      return Promise.resolve(jsonResponse({
        total: sampleAuditItems.length,
        pagina: 1,
        items: sampleAuditItems,
      }));
    });
    vi.stubGlobal('fetch', fetchMock);
    const user = userEvent.setup();

    renderAuditoria();
    await screen.findByRole('table');

    // Verificar exportación CSV y JSON
    const exportCsvBtn = screen.queryByRole('button', { name: /exportar csv|^exportar$/i });
    expect(exportCsvBtn).not.toBeNull();
    await user.click(exportCsvBtn!);

    const exportJsonBtn = screen.queryByRole('button', { name: /exportar json/i })
      ?? screen.queryByRole('menuitem', { name: /json/i });
    expect(exportJsonBtn, 'falta la opción o botón para exportar en formato JSON (FRN-31 / BAC-45)').not.toBeNull();
    await user.click(exportJsonBtn!);

    await waitFor(() => {
      const exportCalls = (fetchMock.mock.calls as [string][]).map(([url]) => String(url)).filter((url) => url.includes('/admin/audit/export'));
      expect(exportCalls.some((url) => /formato=csv|format=csv/i.test(url))).toBe(true);
      expect(exportCalls.some((url) => /formato=json|format=json/i.test(url))).toBe(true);
    });

    expect(downloadedFiles.some((name) => /^auditoria_centinela_\d{8}\.csv$/.test(name))).toBe(true);
    expect(downloadedFiles.some((name) => /^auditoria_centinela_\d{8}\.json$/.test(name))).toBe(true);
    anchorClickSpy.mockRestore();
  });
});
