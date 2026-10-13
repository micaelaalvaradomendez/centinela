import { describe, expect, it } from 'vitest';
import { applicationRoutes } from '@/routes/applicationRoutes';

// FRN-21 (BRG-06-FRN, RF-03 / RF-11): vista "Detalle de Instancia" y enlace desde el inventario.
// Entregable:
//   1. Ruta /instances/:vmid bajo ProtectedLayout con guard de sesión y permiso de acceso.
//   2. Componente de detalle con pestañas base (General, Métricas, Configuración) y botón de retorno.
//   3. Enlace desde la tabla de inventario y desde los toasts de notificación.
// Se verifica el contrato de la ruta de forma determinista, sin acoplar a la implementación interna.

type AnyRoute = { path?: string; index?: boolean; children?: AnyRoute[] };

// Recorre el árbol de RouteObject y devuelve todas las rutas absolutas declaradas.
function collectPaths(route: AnyRoute, parent = ''): string[] {
  const absolute = route.path
    ? route.path.startsWith('/')
      ? route.path
      : `${parent.replace(/\/$/, '')}/${route.path}`
    : parent;
  const own = route.path ? [absolute] : [];
  return [...own, ...(route.children ?? []).flatMap((child) => collectPaths(child, absolute))];
}

const declaredPaths = applicationRoutes.flatMap((route) => collectPaths(route));

describe('FRN-21 - Vista Detalle de Instancia (/instances/:vmid)', () => {
  it('declara la ruta /instances/:vmid dentro del árbol de rutas de la aplicación', () => {
    expect(
      declaredPaths,
      'falta la ruta /instances/:vmid (FRN-21). Rutas declaradas: ' + declaredPaths.join(', '),
    ).toContain('/instances/:vmid');
  });

  it('mantiene la ruta de inventario /instances como padre del detalle', () => {
    expect(declaredPaths).toContain('/instances');
  });
});
