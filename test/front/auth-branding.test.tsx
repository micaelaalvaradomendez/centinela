import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, describe, expect, it, vi } from 'vitest';
import MainLayoutAuth from '@/components/layout_auth/MainLayoutAuth';
import DashboardPage from '@/pages/Dashboard';
import InstancesPage from '@/pages/Instances';
import { withProtectedProviders } from './app-providers';
import { Toaster } from '@/components/ui/toast';

// Pruebas de aceptación y consistencia visual para tareas de interfaz pendientes en actual.md:
// - FIX-58: Imagotipo/logotipo oficial vectorizado en layouts de autenticación.
// - FIX-59: Copy amigable sobre Centinela y 2FA en onboarding / login (sin texto placeholder "imagenes y informacion random").
// - FIX-60: Reubicación del botón "Crear instancia" del Dashboard a la vista de inventario /instances (con control RBAC ADMIN vs OPERATOR).

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FIX-58 - Imagotipo y logotipo oficial vectorizado en vistas de autenticación', () => {
  it('no muestra el placeholder textual "logo" o "logo Centinela" y cuenta con elemento gráfico accesible', () => {
    render(
      <MemoryRouter>
        <MainLayoutAuth>
          <div>Contenido de prueba</div>
        </MainLayoutAuth>
      </MemoryRouter>
    );

    // 1. No debe existir texto plano "logo" suelto como placeholder junto al título.
    const textNodes = screen.queryAllByText(/^\s*logo\s*$/i);
    expect(textNodes, 'persiste el texto plano placeholder "logo" en el layout de autenticación').toHaveLength(0);

    // 2. Debe renderizarse un imagotipo/logotipo accesible (svg o img con alt o aria-label).
    const logoGraphic = screen.queryByRole('img', { name: /centinela/i }) ||
      document.querySelector('aside svg[aria-label*="Centinela" i], aside img[alt*="Centinela" i], aside [data-logo="centinela"]');
    expect(logoGraphic, 'no se encontró un imagotipo o logotipo vectorizado oficial de Centinela con accesibilidad').not.toBeNull();
  });

  it('el enlace del logotipo en vistas de autenticación dirige a inicio o a /login sin provocar bucles', () => {
    render(
      <MemoryRouter>
        <MainLayoutAuth>
          <div>Contenido</div>
        </MainLayoutAuth>
      </MemoryRouter>
    );

    const logoLink = document.querySelector('aside a[href*="/login"], aside a[href="/"]');
    if (logoLink) {
      const href = logoLink.getAttribute('href');
      expect(href, 'el enlace del logo debe apuntar a la raíz o a login').toMatch(/^(\/|\/login)$/);
    }
  });
});

describe('FIX-59 - Sustituir placeholder de desarrollo por mensaje informativo de Centinela y 2FA', () => {
  it('no contiene el texto "imagenes y informacion random" y explica el propósito de la plataforma y 2FA', () => {
    render(
      <MemoryRouter>
        <MainLayoutAuth>
          <div>Contenido de prueba</div>
        </MainLayoutAuth>
      </MemoryRouter>
    );

    // 1. El placeholder informal no debe existir bajo ningún concepto.
    expect(screen.queryByText(/imagenes y informacion random/i)).toBeNull();

    // 2. Debe explicar la propuesta de valor de Centinela (orquestación / infraestructura / Proxmox).
    const plataformaCopy = screen.queryAllByText(/infraestructura|proxmox|orquestaci[oó]n|virtualizad|m[aá]quinas virtuales|servidores/i);
    expect(plataformaCopy.length, 'falta el copy explicativo sobre el propósito de la plataforma Centinela').toBeGreaterThan(0);

    // 3. Debe fundamentar la seguridad y el valor del doble factor (2FA).
    const seguridadCopy = screen.queryAllByText(/2fa|doble factor|dos pasos|seguridad|protecci[oó]n/i);
    expect(seguridadCopy.length, 'falta el copy explicativo sobre la importancia del segundo factor de autenticación (2FA)').toBeGreaterThan(0);

    // 4. Debe contener los títulos estructurados aprobados en actual.md
    expect(screen.getByText(/Tu infraestructura virtual, simplificada y bajo control/i)).toBeInTheDocument();
    expect(screen.getByText(/Protección de infraestructura con doble factor/i)).toBeInTheDocument();
  });
});

describe('FIX-60 - Reubicación del botón "Crear instancia" del Dashboard a la página de Instancias', () => {
  it('el Dashboard no contiene el botón "Crear instancia"', () => {
    window.sessionStorage.setItem('centinela_access', 'access-token-dashboard');
    window.localStorage.setItem('centinela_user', JSON.stringify({
      id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
      rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
    }));

    vi.stubGlobal('fetch', vi.fn().mockImplementation((url: string) => {
      if (String(url).includes('/node/status')) {
        return Promise.resolve(new Response(JSON.stringify({
          cpu: { usagePercent: 10, cores: 4 },
          ram: { usedGb: 2, totalGb: 8, usagePercent: 25 },
          storage: { usedGb: 10, totalGb: 50, usagePercent: 20 },
          uptimeSeconds: 1000,
          instancesSummary: { vms: { running: 0, stopped: 0, paused: 0, total: 0 }, lxc: { running: 0, stopped: 0, paused: 0, total: 0 } },
          stale: false,
          fetchedAt: '2026-10-14T00:00:00Z',
        }), { status: 200, headers: { 'Content-Type': 'application/json' } }));
      }
      return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }));
    }));

    render(withProtectedProviders(
      <MemoryRouter initialEntries={['/dashboard']}>
        <DashboardPage />
      </MemoryRouter>
    ));

    const createButton = screen.queryByRole('button', { name: /crear instancia|nueva instancia/i });
    expect(createButton, 'el botón "Crear instancia" sigue presente indebidamente en el Dashboard').toBeNull();
  });

  it('la página de inventario /instances muestra el botón "Crear instancia" solo para ADMIN y no para OPERATOR', async () => {
    const mockInventory = [
      { id: 101, name: 'srv-1', type: 'vm', node: 'pve', status: 'running', ip: null, cpuUsage: null, ramUsage: null, maxRam: null, nivelAcceso: 'FULL_ACCESS', activeTask: null },
    ];

    vi.stubGlobal('fetch', vi.fn().mockImplementation((url: string) => Promise.resolve(String(url).includes('/instances')
      ? new Response(JSON.stringify(mockInventory), { status: 200, headers: { 'Content-Type': 'application/json' } })
      : new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))));

    // 1. ADMIN en /instances: debe ver el botón.
    window.sessionStorage.setItem('centinela_access', 'access-token-admin');
    window.localStorage.setItem('centinela_user', JSON.stringify({
      id: 'user-admin', organizacionId: 'org-1', nombreCompleto: 'Admin Test', email: 'admin@centinela.local',
      rol: 'ADMIN', instanciasPermitidas: [], tiene2FA: true,
    }));

    const { unmount } = render(withProtectedProviders(
      <Toaster>
        <MemoryRouter initialEntries={['/instances']}>
          <InstancesPage />
        </MemoryRouter>
      </Toaster>
    ));

    expect(await screen.findByRole('button', { name: /crear instancia|nueva instancia/i })).toBeVisible();
    unmount();

    // 2. OPERATOR en /instances: NO debe ver el botón.
    window.sessionStorage.setItem('centinela_access', 'access-token-op');
    window.localStorage.setItem('centinela_user', JSON.stringify({
      id: 'user-op', organizacionId: 'org-1', nombreCompleto: 'Op Test', email: 'op@centinela.local',
      rol: 'OPERATOR', instanciasPermitidas: [101], tiene2FA: true,
    }));

    render(withProtectedProviders(
      <Toaster>
        <MemoryRouter initialEntries={['/instances']}>
          <InstancesPage />
        </MemoryRouter>
      </Toaster>
    ));

    await screen.findByText('srv-1');
    expect(screen.queryByRole('button', { name: /crear instancia|nueva instancia/i })).toBeNull();
  });
});

