import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import fs from 'node:fs';
import path from 'node:path';
import UserDetail from '@/pages/detailsUserPage';
import DashboardPage from '@/pages/Dashboard';
import CreateNewInstances from '@/pages/createNewInstances';
import { withAppProviders, withProtectedProviders } from './app-providers';

// Pruebas de aceptación para los FIX pendientes de Frontend en documentacion/actual.md:
// - FIX-87: Quitar los datos hardcodeados de la ficha de usuario (detailsUserPage.tsx).
// - FIX-89 / FIX-99: Advertencia visible 'Próximamente' en vista directa /instances/new (createNewInstances.tsx).
// - FIX-90: Quitar los placeholders muertos del Dashboard (Dashboard.tsx).
// - FIX-91: Eliminar el componente muerto InstanceTester (InstanceTester.tsx).
// - FIX-100: Manejo seguro de valores undefined en formatCount del Dashboard (Dashboard.tsx).

afterEach(() => {
  vi.unstubAllGlobals();
});

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

const mockUserDetails = {
  id: 'usr-10',
  organizacionId: 'org-1',
  nombreCompleto: 'Usuario de Prueba',
  nombreUsuario: 'usuario_prueba',
  emailUsuario: 'usuario@centinela.local',
  rol: 'OPERATOR',
  activo: true,
  totpVinculado: true,
  cambioContrasenaRequerido: false,
  fechaCreacion: '2026-10-01T12:00:00Z',
  fechaUltimoAcceso: '2026-10-05T15:30:00Z',
  eliminadoEn: null,
  instanciasPermitidas: [],
};

function renderUserDetail(userId = 'usr-10') {
  window.sessionStorage.setItem('centinela_access', 'access-token-fix87');
  window.sessionStorage.setItem('centinela_user', JSON.stringify({ id: 'user-admin', rol: 'ADMIN' }));
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-admin',
    organizacionId: 'org-1',
    nombreCompleto: 'Admin Test',
    email: 'admin@centinela.local',
    rol: 'ADMIN',
    instanciasPermitidas: [],
    tiene2FA: true,
  }));

  const fetchMock = vi.fn().mockImplementation((url: string) => {
    if (String(url).includes('/permissions')) return Promise.resolve(jsonResponse({ items: [] }));
    if (String(url).includes('/instances')) return Promise.resolve(jsonResponse([]));
    if (String(url).includes('/admin/audit')) return Promise.resolve(jsonResponse({ total: 0, pagina: 1, items: [] }));
    return Promise.resolve(jsonResponse(mockUserDetails));
  });
  vi.stubGlobal('fetch', fetchMock);

  return render(
    withAppProviders(
      <MemoryRouter initialEntries={[`/users/${userId}`]}>
        <Routes>
          <Route path="/users/:userId" element={<UserDetail />} />
        </Routes>
      </MemoryRouter>
    )
  );
}

describe('FIX-87 - Quitar los datos hardcodeados de la ficha de usuario', () => {
  it('no muestra valores inventados de cambio de contraseña ni actividad reciente ficticia en la UI', async () => {
    renderUserDetail('usr-10');

    await screen.findByDisplayValue('Usuario de Prueba');

    // 1. Fecha hardcodeada de cambio de contraseña inventada
    expect(screen.queryByText(/02\/05\/2024/i), 'se encontró la fecha hardcodeada "02/05/2024"').toBeNull();

    // 2. Actividades ficticias de placeholders
    expect(screen.queryByText(/Tomó snapshot de Ubuntu Server/i), 'persiste la actividad ficticia de snapshot').toBeNull();
    expect(screen.queryByText(/Inició instancia Windows 11/i), 'persiste la actividad ficticia de Windows 11').toBeNull();
    expect(screen.queryByText(/192\.168\.1\.77/i), 'persiste la IP hardcodeada 192.168.1.77').toBeNull();
  });

  it('el archivo detailsUserPage.tsx no define arrays estáticos con datos ficticios', () => {
    const filePath = path.resolve(__dirname, '../../frontend/centinela/src/pages/detailsUserPage.tsx');
    expect(fs.existsSync(filePath), 'detailsUserPage.tsx no existe').toBe(true);

    const content = fs.readFileSync(filePath, 'utf-8');
    expect(content).not.toContain('02/05/2024');
    expect(content).not.toContain('Tomó snapshot de Ubuntu Server');
    expect(content).not.toContain('Inició instancia Windows 11');
    expect(content).not.toContain('192.168.1.77');
  });
});

describe('FIX-89 / FIX-99 - Advertencia visible "Próximamente" en vista directa /instances/new', () => {
  it('el asistente estático debe advertir explícitamente que es próximamente o estar deshabilitado', () => {
    render(
      withProtectedProviders(
        <MemoryRouter initialEntries={['/instances/new']}>
          <CreateNewInstances />
        </MemoryRouter>
      )
    );

    // Si la pantalla de asistente estático se renderiza sin integración con el backend,
    // debe advertir al usuario con un rótulo o aviso explícito ("próximamente", "en desarrollo", "no disponible").
    const aviso = screen.queryByText(/pr[oó]ximamente|en desarrollo|no disponible|fase 2|etapa 2/i);
    expect(
      aviso,
      'el asistente de /instances/new no advierte de forma explícita que está en desarrollo / próximamente',
    ).not.toBeNull();
  });
});

describe('FIX-90 - Quitar los placeholders muertos del Dashboard', () => {
  it('no contiene bloques de texto estáticos muertos ("Sin instancias", "Sin novedades", "Sin alertas") sin funcionalidad', async () => {
    window.sessionStorage.setItem('centinela_access', 'token-dashboard-test');
    window.localStorage.setItem('centinela_user', JSON.stringify({
      id: 'admin-1', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
      rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
    }));

    const mockStatus = {
      cpu: { usagePercent: 20, cores: 4 },
      ram: { usedGb: 4, totalGb: 16, usagePercent: 25 },
      storage: { usedGb: 10, totalGb: 100, usagePercent: 10 },
      uptimeSeconds: 12000,
      instancesSummary: {
        vms: { running: 1, stopped: 0, paused: 0, total: 1 },
        lxc: { running: 0, stopped: 0, paused: 0, total: 0 },
      },
      stale: false,
      fetchedAt: '2026-10-14T08:00:00Z',
    };

    vi.stubGlobal('fetch', vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/node/status')) return Promise.resolve(jsonResponse(mockStatus));
      if (String(url).includes('/events/ticket')) return Promise.resolve(jsonResponse({ ticket: 'dummy' }));
      return Promise.resolve(jsonResponse({}));
    }));

    render(
      withProtectedProviders(
        <MemoryRouter initialEntries={['/dashboard']}>
          <DashboardPage />
        </MemoryRouter>
      )
    );

    await screen.findByText(/resumen general de tus instancias/i);

    // FIX-90: No deben quedar los estados vacíos planos ni botones muertos
    expect(screen.queryByText(/^Sin instancias$/i), 'persiste el placeholder muerto "Sin instancias"').toBeNull();
    expect(screen.queryByText(/^Sin novedades$/i), 'persiste el placeholder muerto "Sin novedades"').toBeNull();
    expect(screen.queryByText(/^Sin alertas$/i), 'persiste el placeholder muerto "Sin alertas"').toBeNull();
  });

  it('no contiene botones sin acción (dead buttons) en el Dashboard', async () => {
    const filePath = path.resolve(__dirname, '../../frontend/centinela/src/pages/Dashboard.tsx');
    expect(fs.existsSync(filePath), 'Dashboard.tsx no existe').toBe(true);

    const content = fs.readFileSync(filePath, 'utf-8');
    // En Dashboard.tsx no deben existir botones sin onClick como <button>Ver todas</button> o <button>Ver todo</button>
    const deadButtons = content.match(/<button[^>]*>(?:Ver todas|Ver todo)<\/button>/g);
    expect(deadButtons, 'se encontraron botones sin handler onClick en Dashboard.tsx: ' + deadButtons?.join(', ')).toBeNull();
  });
});

describe('FIX-91 - Eliminar el componente muerto InstanceTester', () => {
  it('el componente muerto InstanceTester.tsx no existe en el código de producción de componentes', () => {
    const componentPath = path.resolve(
      __dirname,
      '../../frontend/centinela/src/components/features/instances/components/InstanceTester.tsx'
    );
    expect(
      fs.existsSync(componentPath),
      'InstanceTester.tsx sigue existiendo en el código de componentes de producción; debe ser eliminado o movido a test'
    ).toBe(false);
  });
});

describe('FIX-100 - Manejo seguro de valores undefined en formatCount del Dashboard', () => {
  it('el Dashboard no arroja TypeError y tolera templates.count undefined o nulo', async () => {
    window.sessionStorage.setItem('centinela_access', 'token-dashboard-fix100');
    window.localStorage.setItem('centinela_user', JSON.stringify({
      id: 'admin-1', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
      rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
    }));

    const mockStatus = {
      cpu: { usagePercent: 20, cores: 4 },
      ram: { usedGb: 4, totalGb: 16, usagePercent: 25 },
      storage: { usedGb: 10, totalGb: 100, usagePercent: 10 },
      uptimeSeconds: 12000,
      instancesSummary: {
        vms: { running: 1, stopped: 0, paused: 0, total: 1 },
        lxc: { running: 0, stopped: 0, paused: 0, total: 0 },
      },
      stale: false,
      fetchedAt: '2026-10-14T08:00:00Z',
    };

    vi.stubGlobal('fetch', vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/node/status')) return Promise.resolve(jsonResponse(mockStatus));
      if (String(url).includes('/events/ticket')) return Promise.resolve(jsonResponse({ ticket: 'dummy' }));
      if (String(url).includes('/node/templates')) return Promise.resolve(jsonResponse({ items: [] }));
      return Promise.resolve(jsonResponse({}));
    }));

    render(
      withProtectedProviders(
        <MemoryRouter initialEntries={['/dashboard']}>
          <DashboardPage />
        </MemoryRouter>
      )
    );

    // Debe renderizarse sin desmontar el árbol por TypeError
    expect(await screen.findByText(/resumen general de tus instancias/i)).toBeInTheDocument();
  });
});

