import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import React from 'react';

// FRN-19A (ex FRN-13A): Maquetado y medidores de recursos del Host (CPU / RAM / Almacenamiento).
// Entregable: componentes reutilizables en features/dashboard con barras o medidores de:
//   - CPU (% y núcleos)
//   - RAM (GB usados sobre el total)
//   - Almacenamiento (GB o TB usados sobre el total)
// Criterio de éxito:
//   - Renderizan valores de 0 % a 100 %.
//   - Cambio de color según saturación: normal por debajo del 70 % y advertencia desde el 70 % (D3).
//   - Pruebas de componente, incluidos los valores de borde 69 % y 70 %.

type DashboardModule = {
  CpuGauge?: React.ComponentType<any>;
  RamGauge?: React.ComponentType<any>;
  StorageGauge?: React.ComponentType<any>;
  HostResourceGauges?: React.ComponentType<any>;
  ResourceGauge?: React.ComponentType<any>;
  [key: string]: unknown;
};

const dashboardCandidateModules = import.meta.glob<DashboardModule>([
  '@/components/features/dashboard/**/*.{ts,tsx,js,jsx}',
  '@/pages/Dashboard.{ts,tsx,js,jsx}',
]);

async function loadResourceGauges() {
  for (const load of Object.values(dashboardCandidateModules)) {
    const mod = await load().catch(() => ({}) as DashboardModule);
    const candidate = mod.CpuGauge || mod.RamGauge || mod.StorageGauge || mod.HostResourceGauges || mod.ResourceGauge;
    if (candidate) return { candidate, mod };
  }
  throw new Error('FRN-19A no implementada: no se encontraron componentes de medidores de recursos (CPU/RAM/Almacenamiento) en features/dashboard');
}

describe('FRN-19A - Medidores de recursos del Host (CPU / RAM / Almacenamiento)', () => {
  it('los componentes de medidores existen en features/dashboard', async () => {
    const { candidate } = await loadResourceGauges();
    expect(candidate).toBeDefined();
  });

  it('renderiza valores normales por debajo del 70% (valor de borde: 69%)', async () => {
    const { candidate: Gauge } = await loadResourceGauges();
    const { container } = render(
      <Gauge
        tipo="cpu"
        usagePercent={69}
        cores={8}
        usedGb={11}
        totalGb={16}
        title="Uso de CPU"
      />
    );

    expect(screen.getByText(/69\s*%/)).toBeInTheDocument();
    // En 69%, no debe tener clases ni atributos de advertencia/warning
    const warningIndicator = container.querySelector('[data-state="warning"], .text-warning, .border-warning, .bg-warning, [aria-label*="advertencia" i]');
    expect(warningIndicator).toBeNull();
  });

  it('activa el estado de advertencia desde el 70% inclusive (valor de borde: 70%)', async () => {
    const { candidate: Gauge } = await loadResourceGauges();
    const { container } = render(
      <Gauge
        tipo="cpu"
        usagePercent={70}
        cores={8}
        usedGb={11.2}
        totalGb={16}
        title="Uso de CPU"
      />
    );

    expect(screen.getByText(/70\s*%/)).toBeInTheDocument();
    // En 70% o más, debe indicar advertencia (clase, data-state o estilo)
    const warningIndicator = container.querySelector('[data-state="warning"], .text-warning, .border-warning, .bg-warning, [aria-label*="advertencia" i], .text-amber-500, .text-red-500');
    expect(warningIndicator).not.toBeNull();
  });

  it('renderiza extremos 0% y 100% y núcleos/GB correctamente', async () => {
    const { candidate: Gauge } = await loadResourceGauges();
    const { rerender } = render(
      <Gauge tipo="ram" usagePercent={0} usedGb={0} totalGb={32} />
    );
    expect(screen.getByText(/0\s*%/)).toBeInTheDocument();

    rerender(<Gauge tipo="ram" usagePercent={100} usedGb={32} totalGb={32} />);
    expect(screen.getByText(/100\s*%/)).toBeInTheDocument();
    expect(screen.getByText(/32/)).toBeInTheDocument();
  });
});

