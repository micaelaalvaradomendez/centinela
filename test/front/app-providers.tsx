import React from 'react';

// La aplicación monta los providers de src/context en dos lugares: AuthProvider en App.tsx (SEC-03) y
// EventsProvider en ProtectedLayout (FRN-17A), que envuelve a todas las páginas protegidas.
// Las pruebas que renderizan rutas o páginas usan los mismos: si no, un componente que lee el contexto
// (Sidebar, PermissionGate, useInstances) falla por falta de Provider, cosa que en la app no pasa.
// No se fija el archivo: se toma todo export `*Provider` de src/context.
type ContextModule = Record<string, unknown>;
const contextModules = import.meta.glob<ContextModule>(
  ['@/context/**/*.{ts,tsx,js,jsx}', '@/components/**/context/**/*.{ts,tsx,js,jsx}'],
  { eager: true },
);

const providers = Object.values(contextModules).flatMap((module) =>
  Object.entries(module)
    .filter(([name, value]) => /Provider$/.test(name) && typeof value === 'function')
    .map(([name, value]) => ({ name, Provider: value as React.ComponentType<{ children?: React.ReactNode }> })),
);

function wrap(children: React.ReactNode, include: (name: string) => boolean): React.ReactElement {
  const active = providers.filter(({ name }) => include(name));
  active.sort((a, b) => {
    if (/Auth/.test(a.name)) return -1;
    if (/Auth/.test(b.name)) return 1;
    if (/Event/.test(a.name)) return -1;
    if (/Event/.test(b.name)) return 1;
    if (/PendingInstanceCreations/.test(a.name)) return 1;
    if (/PendingInstanceCreations/.test(b.name)) return -1;
    return 0;
  });
  return active.reduceRight<React.ReactElement>(
    (inner, { Provider }) => React.createElement(Provider, null, inner),
    React.createElement(React.Fragment, null, children),
  );
}

// Lo que monta App.tsx (todos los providers salvo los que monta ProtectedLayout). Para renderizar rutas completas:
// ProtectedLayout ya monta EventsProvider y PendingInstanceCreationsProvider dentro de las rutas protegidas.
export function withAppProviders(children: React.ReactNode): React.ReactElement {
  return wrap(children, (name) => !/Event|PendingInstanceCreations/.test(name));
}

// App.tsx más ProtectedLayout (con EventsProvider). Para renderizar una página protegida suelta
// (Instances, Dashboard), que en la app siempre está debajo de ProtectedLayout.
export function withProtectedProviders(children: React.ReactNode): React.ReactElement {
  return wrap(children, () => true);
}
