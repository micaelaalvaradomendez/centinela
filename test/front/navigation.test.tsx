import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { describe, expect, it } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';

const storedUser = {
  id: 'user-1',
  organizacionId: 'org-1',
  nombreCompleto: 'Administrador de prueba',
  email: 'admin@centinela.local',
  rol: 'ADMIN',
  instanciasPermitidas: [] as number[],
  tiene2FA: true,
};

function seedSession(user: Partial<typeof storedUser> = {}) {
  window.sessionStorage.setItem('centinela_access', 'access-token');
  window.localStorage.setItem('centinela_user', JSON.stringify({ ...storedUser, ...user }));
}

function renderApplication(path: string, authenticated: boolean, user: Partial<typeof storedUser> = {}) {
  if (authenticated) seedSession(user);
  const router = createMemoryRouter(applicationRoutes, { initialEntries: [path] });
  render(<RouterProvider router={router} />);
  return router;
}

// SEC-03 pide el hook y el gate en src/context/ sin fijar el nombre del archivo:
// se busca cualquier módulo de esa carpeta que exporte usePermissions y PermissionGate.
type PermissionsModule = {
  usePermissions?: () => {
    isAdmin: boolean;
    isOperator: boolean;
    hasRole: (role: string) => boolean;
    canAccessInstance: (vmid: number) => boolean;
    canOperateInstance?: (vmid: number) => boolean;
  };
  PermissionGate?: React.ComponentType<{ requiredRole: string; fallback?: React.ReactNode; children?: React.ReactNode }>;
  [key: string]: unknown;
};
const contextModules = import.meta.glob<PermissionsModule>('@/context/**/*.{ts,tsx,js,jsx}');

async function loadPermissionsModule(): Promise<Required<Pick<PermissionsModule, 'usePermissions' | 'PermissionGate'>> & PermissionsModule> {
  for (const load of Object.values(contextModules)) {
    const module = await load();
    if (typeof module.usePermissions === 'function' && module.PermissionGate) {
      return module as never;
    }
  }
  throw new Error(`SEC-03 no implementada: ningún módulo de src/context exporta usePermissions y PermissionGate (revisados: ${Object.keys(contextModules).join(', ') || 'ninguno'})`);
}

// FIX-30 (FRN-18) solo pide el helper canOperateInstance en usePermissions(): se acepta el hook
// en src/context o en src/hooks. La ubicación y el resto del contrato (isOperator, hasRole y
// PermissionGate en src/context) los exigen los casos de SEC-03.
const hookModules = import.meta.glob<PermissionsModule>(['@/context/**/*.{ts,tsx,js,jsx}', '@/hooks/**/*.{ts,tsx,js,jsx}']);

async function loadUsePermissionsModule(): Promise<Required<Pick<PermissionsModule, 'usePermissions'>> & PermissionsModule> {
  for (const load of Object.values(hookModules)) {
    const module = await load();
    if (typeof module.usePermissions === 'function') return module as never;
  }
  throw new Error(`FIX-30: ningún módulo de src/context ni de src/hooks exporta usePermissions (revisados: ${Object.keys(hookModules).join(', ') || 'ninguno'})`);
}

// Si el módulo exporta un Provider, se usa; si el hook lee la sesión directamente, no hace falta.
function withProvider(module: PermissionsModule, children: React.ReactNode) {
  const providerName = Object.keys(module).find((name) => /Provider$/.test(name));
  const Provider = providerName ? module[providerName] as React.ComponentType<{ children: React.ReactNode }> : null;
  return Provider ? React.createElement(Provider, null, children) : children;
}

describe('FRN-03 - navbar y rutas base', () => {
  it('redirige a login cuando no existe una sesión', async () => {
    const router = renderApplication('/dashboard', false);

    expect(await screen.findByRole('heading', { name: 'Iniciar sesión' })).toBeVisible();
    expect(router.state.location.pathname).toBe('/login');
  });

  it('permite navegar entre Dashboard e Instancias con una sesión', async () => {
    const user = userEvent.setup();
    const router = renderApplication('/dashboard', true);

    expect(await screen.findByRole('heading', { name: /Hola, Admin/ })).toBeVisible();
    await user.click(screen.getByRole('link', { name: 'Instancias' }));

    expect(await screen.findByRole('heading', { name: 'Inventario de instancias' })).toBeVisible();
    expect(router.state.location.pathname).toBe('/instances');
  });
});

describe('SEC-03 - Contexto y sistema reactivo de permisos en Frontend', () => {
  it('usePermissions expone isAdmin, isOperator, hasRole y canAccessInstance según la sesión', async () => {
    const module = await loadPermissionsModule();
    seedSession({ rol: 'OPERATOR', instanciasPermitidas: [101] });
    let permissions: ReturnType<NonNullable<PermissionsModule['usePermissions']>> | undefined;
    function Probe() {
      permissions = module.usePermissions();
      return null;
    }

    render(<>{withProvider(module, <Probe />)}</>);

    expect(permissions).toMatchObject({ isAdmin: false, isOperator: true });
    expect(permissions!.hasRole('OPERATOR')).toBe(true);
    expect(permissions!.hasRole('ADMIN')).toBe(false);
    expect(permissions!.canAccessInstance(101)).toBe(true);
    expect(permissions!.canAccessInstance(999)).toBe(false);
  });

  it('FIX-30 canOperateInstance distingue FULL_ACCESS de READ_ONLY (FRN-18)', async () => {
    const module = await loadUsePermissionsModule();
    // La sesión del operador informa el nivel por instancia (contrato de GET /permissions).
    seedSession({
      rol: 'OPERATOR',
      instanciasPermitidas: [101, 102],
      permisos: [{ vmid: 101, nivelAcceso: 'FULL_ACCESS' }, { vmid: 102, nivelAcceso: 'READ_ONLY' }],
    } as never);
    let permissions: ReturnType<NonNullable<PermissionsModule['usePermissions']>> | undefined;
    function Probe() {
      permissions = module.usePermissions();
      return null;
    }

    render(<>{withProvider(module, <Probe />)}</>);

    expect(typeof permissions?.canOperateInstance, 'usePermissions no expone canOperateInstance').toBe('function');
    expect(permissions!.canAccessInstance(102)).toBe(true);
    expect(permissions!.canOperateInstance!(101)).toBe(true);
    expect(permissions!.canOperateInstance!(102)).toBe(false);
    expect(permissions!.canOperateInstance!(999)).toBe(false);
  });

  it('PermissionGate muestra el contenido al ADMIN y el fallback al OPERATOR', async () => {
    const { PermissionGate, ...module } = await loadPermissionsModule();
    const gate = () => withProvider(module, (
      <PermissionGate requiredRole="ADMIN" fallback={<p>Sin permiso</p>}>
        <button type="button">Eliminar usuario</button>
      </PermissionGate>
    ));

    seedSession({ rol: 'OPERATOR' });
    const { unmount } = render(<>{gate()}</>);
    expect(screen.queryByRole('button', { name: 'Eliminar usuario' })).not.toBeInTheDocument();
    expect(screen.getByText('Sin permiso')).toBeInTheDocument();
    unmount();

    seedSession({ rol: 'ADMIN' });
    render(<>{gate()}</>);
    expect(screen.getByRole('button', { name: 'Eliminar usuario' })).toBeInTheDocument();
  });

  it('un OPERATOR no ve accesos administrativos en el menú (Usuarios, Crear usuario, Auditoría)', async () => {
    renderApplication('/dashboard', true, { rol: 'OPERATOR', nombreCompleto: 'Operador Test' });

    expect(await screen.findByRole('heading', { name: /Hola/i })).toBeVisible();
    for (const name of ['Usuarios', 'Crear usuario', 'Auditoría']) {
      expect(screen.queryByRole('link', { name })).not.toBeInTheDocument();
    }
  });

  // "Crear usuario" dejó de ser un enlace del menú (commit a387e10): el alta se inicia
  // desde el botón de /users. El menú administrativo queda en Usuarios y Auditoría.
  it('un ADMIN ve los accesos administrativos en el menú: Usuarios y Auditoría', async () => {
    renderApplication('/dashboard', true);

    expect(await screen.findByRole('heading', { name: /Hola, Admin/ })).toBeVisible();
    for (const name of ['Usuarios', 'Auditoría']) {
      expect(screen.getByRole('link', { name })).toBeVisible();
    }
  });
});
