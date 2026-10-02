import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { createMemoryRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { describe, expect, it } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';
import { withAppProviders } from './app-providers';

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
  render(withAppProviders(<RouterProvider router={router} />));
  return router;
}

// SEC-03: el criterio de éxito es de comportamiento (el operador no ve accesos de administrador,
// canAccessInstance evalúa en memoria y el render condicional está probado); no fija nombres ni carpeta.
// El entregable propone "usePermissions() / useAuthUser()" y "un contexto con Provider es opcional":
// se acepta cualquiera de esos hooks o useAuth, en src/context o src/hooks, con los helpers como
// booleanos (isAdmin) o como funciones (isAdmin()). PermissionGate puede ser export con nombre o el
// export por defecto de un archivo PermissionGate.*. Los hooks se renderizan con los providers de la app.
type Permissions = {
  isAdmin?: boolean | (() => boolean);
  isOperator?: boolean | (() => boolean);
  hasRole?: (role: string) => boolean;
  canAccessInstance?: (vmid: number) => boolean;
  canOperateInstance?: (vmid: number) => boolean;
};
type PermissionsHook = { name: string; use: () => Permissions };
const HOOK_NAMES = /^use(Permissions|AuthUser|Auth)$/;
const hookModules = import.meta.glob<Record<string, unknown>>(['@/context/**/*.{ts,tsx,js,jsx}', '@/hooks/**/*.{ts,tsx,js,jsx}']);
const gateModules = import.meta.glob<Record<string, unknown>>(['@/context/**/*.{ts,tsx,js,jsx}', '@/hooks/**/*.{ts,tsx,js,jsx}', '@/components/**/*.{tsx,jsx}']);

async function loadPermissionHooks(): Promise<PermissionsHook[]> {
  const hooks: PermissionsHook[] = [];
  for (const [path, load] of Object.entries(hookModules)) {
    const module = await load();
    for (const [name, value] of Object.entries(module)) {
      if (HOOK_NAMES.test(name) && typeof value === 'function') hooks.push({ name: `${name} (${path})`, use: value as () => Permissions });
    }
  }
  if (hooks.length === 0) throw new Error(`ningún módulo de src/context ni de src/hooks exporta usePermissions, useAuthUser ni useAuth (revisados: ${Object.keys(hookModules).join(', ') || 'ninguno'})`);
  return hooks;
}

async function loadPermissionGate(): Promise<React.ComponentType<{ requiredRole: string; fallback?: React.ReactNode; children?: React.ReactNode }>> {
  for (const [path, load] of Object.entries(gateModules)) {
    const module = await load();
    const gate = module.PermissionGate ?? (/\/PermissionGate\.[jt]sx?$/.test(path) ? module.default : undefined);
    if (typeof gate === 'function') return gate as never;
  }
  throw new Error('SEC-03 incompleta: ningún módulo de src/context, src/hooks ni src/components exporta PermissionGate');
}

// Ejecuta el hook dentro de los providers de la app y devuelve lo que expone.
function readHook(hook: PermissionsHook): Permissions {
  let permissions: Permissions | undefined;
  function Probe() {
    permissions = hook.use();
    return null;
  }
  render(withAppProviders(<Probe />));
  return permissions!;
}

const flag = (value: boolean | (() => boolean) | undefined) => (typeof value === 'function' ? value() : value);

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
  it('un hook de permisos expone isAdmin, isOperator, hasRole y canAccessInstance según la sesión', async () => {
    const hooks = await loadPermissionHooks();
    seedSession({ rol: 'OPERATOR', instanciasPermitidas: [101] });
    const complete = hooks.filter((hook) => {
      const permissions = readHook(hook);
      return ['isAdmin', 'isOperator', 'hasRole', 'canAccessInstance'].every((helper) => helper in (permissions ?? {}));
    });
    expect(complete.map((hook) => hook.name), `ningún hook expone los cuatro helpers (revisados: ${hooks.map((h) => h.name).join(', ')})`).not.toHaveLength(0);

    for (const hook of complete) {
      const permissions = readHook(hook);
      expect(flag(permissions.isAdmin), `${hook.name}: isAdmin`).toBe(false);
      expect(flag(permissions.isOperator), `${hook.name}: isOperator`).toBe(true);
      expect(permissions.hasRole!('OPERATOR'), `${hook.name}: hasRole('OPERATOR')`).toBe(true);
      expect(permissions.hasRole!('ADMIN'), `${hook.name}: hasRole('ADMIN')`).toBe(false);
      expect(permissions.canAccessInstance!(101), `${hook.name}: canAccessInstance(101)`).toBe(true);
      expect(permissions.canAccessInstance!(999), `${hook.name}: canAccessInstance(999)`).toBe(false);
    }
  });

  it('FIX-30 canOperateInstance distingue FULL_ACCESS de READ_ONLY (FRN-18)', async () => {
    const hooks = await loadPermissionHooks();
    // La sesión del operador informa el nivel por instancia (contrato de GET /permissions).
    seedSession({
      rol: 'OPERATOR',
      instanciasPermitidas: [101, 102],
      permisos: [{ vmid: 101, nivelAcceso: 'FULL_ACCESS' }, { vmid: 102, nivelAcceso: 'READ_ONLY' }],
    } as never);
    const operating = hooks.filter((hook) => typeof readHook(hook)?.canOperateInstance === 'function');
    expect(operating.map((hook) => hook.name), 'ningún hook de permisos expone canOperateInstance').not.toHaveLength(0);

    for (const hook of operating) {
      const permissions = readHook(hook);
      expect(permissions.canOperateInstance!(101), `${hook.name}: canOperateInstance(101) con FULL_ACCESS`).toBe(true);
      expect(permissions.canOperateInstance!(102), `${hook.name}: canOperateInstance(102) con READ_ONLY`).toBe(false);
      expect(permissions.canOperateInstance!(999), `${hook.name}: canOperateInstance(999) sin asignar`).toBe(false);
    }
  });

  // FIX-38 corrige usePermissions (hooks/usePermissions.ts:26, entregable 2). El mismo defecto en el
  // useAuth() que agregó SEC-03 es FIX-42: cada caso evalúa solo el helper de su tarea.
  const READ_ONLY_AS_ROLE = {
    rol: 'READ_ONLY',
    instanciasPermitidas: [101],
    permisos: [{ vmid: 101, nivelAcceso: 'READ_ONLY' }],
  } as never;

  it('FIX-38 un rol inexistente (READ_ONLY) no obtiene acceso a instancias en usePermissions: READ_ONLY es un nivel, no un rol', async () => {
    const hooks = (await loadPermissionHooks()).filter((hook) => hook.name.startsWith('usePermissions '));
    expect(hooks, 'no existe usePermissions (FRN-18 / FIX-30)').not.toHaveLength(0);
    seedSession(READ_ONLY_AS_ROLE);
    for (const hook of hooks) {
      expect(readHook(hook).canAccessInstance!(101), `${hook.name}: canAccessInstance(101) con rol READ_ONLY`).toBe(false);
    }
  });

  it('FIX-42 un rol inexistente (READ_ONLY) no obtiene acceso a instancias en el resto de los hooks de permisos (useAuth de SEC-03)', async () => {
    const hooks = (await loadPermissionHooks())
      .filter((hook) => !hook.name.startsWith('usePermissions ') && typeof readHook(hook)?.canAccessInstance === 'function');
    expect(hooks, 'no hay otro hook que exponga canAccessInstance (SEC-03)').not.toHaveLength(0);
    seedSession(READ_ONLY_AS_ROLE);
    for (const hook of hooks) {
      expect(readHook(hook).canAccessInstance!(101), `${hook.name}: canAccessInstance(101) con rol READ_ONLY`).toBe(false);
    }
  });

  it('PermissionGate muestra el contenido al ADMIN y el fallback al OPERATOR', async () => {
    const PermissionGate = await loadPermissionGate();
    const gate = () => withAppProviders(
      <PermissionGate requiredRole="ADMIN" fallback={<p>Sin permiso</p>}>
        <button type="button">Eliminar usuario</button>
      </PermissionGate>,
    );

    seedSession({ rol: 'OPERATOR' });
    const { unmount } = render(gate());
    expect(screen.queryByRole('button', { name: 'Eliminar usuario' })).not.toBeInTheDocument();
    expect(screen.getByText('Sin permiso')).toBeInTheDocument();
    unmount();

    seedSession({ rol: 'ADMIN' });
    render(gate());
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
