import React from 'react';

// App.tsx monta los providers de src/context alrededor del router (desde SEC-03, AuthProvider).
// Las pruebas que renderizan rutas o páginas tienen que usar los mismos: si no, un componente que
// lee el contexto (Sidebar, PermissionGate, informationOfUser) falla por falta de Provider, cosa
// que en la aplicación real no pasa. No se fija el archivo: se toma todo export `*Provider`.
type ContextModule = Record<string, unknown>;
const contextModules = import.meta.glob<ContextModule>('@/context/**/*.{ts,tsx,js,jsx}', { eager: true });

const providers = Object.values(contextModules).flatMap((module) =>
  Object.entries(module)
    .filter(([name, value]) => /Provider$/.test(name) && typeof value === 'function')
    .map(([, value]) => value as React.ComponentType<{ children?: React.ReactNode }>),
);

export function withAppProviders(children: React.ReactNode): React.ReactElement {
  return providers.reduceRight<React.ReactElement>(
    (inner, Provider) => React.createElement(Provider, null, inner),
    React.createElement(React.Fragment, null, children),
  );
}
