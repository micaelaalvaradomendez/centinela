import { cleanup, render, screen } from '@testing-library/react';
import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, resolve } from 'node:path';
import React from 'react';
import { MemoryRouter } from 'react-router';
import { describe, expect, it, vi } from 'vitest';
import { withAppProviders } from './app-providers';

// FRN-19A (ex FRN-13A): medidores de recursos del Host (CPU / RAM / Almacenamiento).
// Entregable: componentes reutilizables en features/dashboard con barras o gauges de CPU (% y núcleos),
// RAM (GB usados sobre el total) y almacenamiento (GB o TB usados sobre el total), con datos de prueba.
// Criterio de éxito: renderizan de 0 % a 100 %; cambian de color según la saturación, normal por
// debajo del 70 % y advertencia desde el 70 % (D3); tienen pruebas de componente con los bordes 69 y 70.
//
// La tarea no fija nombres de componentes, props ni clases: se toman los componentes exportados de
// features/dashboard que, con un valor, muestran ese porcentaje. Se les pasan los nombres de prop más
// comunes para el porcentaje, y el color se compara por diferencia visual entre 69 % y 70 %, no por clase.
type ComponentModule = Record<string, unknown>;
const dashboardModules = import.meta.glob<ComponentModule>(['@/components/features/dashboard/**/*.{tsx,jsx}', '!**/*.{test,spec}.*']);

function propsFor(value: number, tipo: 'cpu' | 'ram' | 'storage' = 'cpu') {
  const usedGb = Math.round((value / 100) * 16 * 100) / 100;
  return {
    tipo, type: tipo, kind: tipo, resource: tipo, label: 'CPU', title: 'CPU', name: 'CPU',
    usagePercent: value, percent: value, percentage: value, porcentaje: value, value, uso: value, usage: value,
    cores: 8, nucleos: 8, usedGb, totalGb: 16, used: usedGb, total: 16, max: 100,
  };
}

function tryRender(Component: React.ComponentType<any>, props: Record<string, unknown>) {
  const error = vi.spyOn(console, 'error').mockImplementation(() => {});
  try {
    return render(withAppProviders(<MemoryRouter><Component {...props} /></MemoryRouter>));
  } catch {
    return null;
  } finally {
    error.mockRestore();
  }
}

async function loadGauges() {
  const gauges: { name: string; Component: React.ComponentType<any> }[] = [];
  for (const [path, load] of Object.entries(dashboardModules)) {
    const module = await load().catch(() => ({}) as ComponentModule);
    for (const [exportName, value] of Object.entries(module)) {
      const name = exportName === 'default' ? path.split('/').pop()!.replace(/\.[jt]sx$/, '') : exportName;
      if (typeof value !== 'function' || !/^[A-Z]/.test(name)) continue;
      const rendered = tryRender(value as React.ComponentType<any>, propsFor(69));
      const showsPercent = rendered !== null && screen.queryAllByText(/\b69(?:[.,]0+)?\s*%/).length > 0;
      cleanup();
      if (showsPercent) gauges.push({ name: `${name} (${path})`, Component: value as React.ComponentType<any> });
    }
  }
  if (gauges.length === 0) {
    throw new Error(`FRN-19A no implementada: ningún componente de features/dashboard muestra el porcentaje que recibe (revisados: ${Object.keys(dashboardModules).join(', ') || 'ninguno'})`);
  }
  return gauges;
}

// Firma visual: clases, colores y atributos de estado de todo el árbol, sin lo que depende del valor
// (el texto, el ancho de la barra, los valores ARIA numéricos ni las clases con el número).
function visualSignature(container: HTMLElement, compared: number[]) {
  const numbers = compared.map(String);
  const mentions = (text: string) => numbers.some((number) => text.includes(number));
  const parts: string[] = [];
  for (const element of Array.from(container.querySelectorAll('*'))) {
    const classes = (element.getAttribute('class') ?? '').split(/\s+/)
      .filter((token) => token && !mentions(token) && !/\[.*\]/.test(token)).sort().join('.');
    const style = (element.getAttribute('style') ?? '').split(';').map((rule) => rule.trim())
      .filter((rule) => rule && !/^(width|height|max-width|min-width|transform|inset|left|right|top|bottom|stroke-dash\w*|--[\w-]+)\s*:/i.test(rule) && !mentions(rule))
      .sort().join(';');
    const attributes = Array.from(element.attributes)
      .filter((attribute) => /^(data-|fill$|stroke$|color$)/.test(attribute.name) && !mentions(attribute.value))
      .map((attribute) => `${attribute.name}=${attribute.value}`).sort().join(',');
    parts.push(`${element.tagName}{${classes}|${style}|${attributes}}`);
  }
  return parts.join('');
}

// Se excluyen de las dos firmas los números de ambos valores comparados, para que el filtro sea simétrico.
function signatureAt(Component: React.ComponentType<any>, value: number, compared: number[]) {
  const rendered = tryRender(Component, propsFor(value));
  expect(rendered, `el medidor no se renderiza con ${value} %`).not.toBeNull();
  const signature = visualSignature(rendered!.container, compared);
  cleanup();
  return signature;
}

function frontendTestFiles(): string[] {
  // Misma raíz que vitest.config.mjs: el submódulo, salvo que CENTINELA_FRONTEND_DIR apunte a otra copia.
  const root = process.env.CENTINELA_FRONTEND_DIR ?? resolve(process.cwd(), '../../frontend/centinela');
  const files: string[] = [];
  const walk = (directory: string) => {
    for (const entry of readdirSync(directory)) {
      if (entry === 'node_modules' || entry === 'dist' || entry.startsWith('.')) continue;
      const path = join(directory, entry);
      if (statSync(path).isDirectory()) walk(path);
      else if (/\.(test|spec)\.[jt]sx?$/.test(entry)) files.push(path);
    }
  };
  walk(root);
  return files;
}

describe('FRN-19A - Medidores de recursos del Host (CPU / RAM / Almacenamiento)', () => {
  it('features/dashboard tiene medidores que muestran el porcentaje de 0 % a 100 %', async () => {
    const gauges = await loadGauges();
    for (const { name, Component } of gauges) {
      for (const value of [0, 100]) {
        expect(tryRender(Component, propsFor(value)), `${name} no se renderiza con ${value} %`).not.toBeNull();
        expect(screen.queryAllByText(new RegExp(`\\b${value}(?:[.,]0+)?\\s*%`)).length, `${name} no muestra ${value} %`).toBeGreaterThan(0);
        cleanup();
      }
    }
  });

  it('D3: el aspecto es el mismo por debajo del 70 % (20 y 69) y cambia al llegar al 70 % (bordes 69 y 70)', async () => {
    const gauges = await loadGauges();
    for (const { name, Component } of gauges) {
      expect(signatureAt(Component, 69, [20, 69]), `${name}: con 69 % el medidor tiene que verse "normal", igual que con 20 %`)
        .toBe(signatureAt(Component, 20, [20, 69]));
      expect(signatureAt(Component, 70, [69, 70]), `${name}: con 70 % el medidor tiene que pasar a "advertencia" (D3), pero se ve igual que con 69 %`)
        .not.toBe(signatureAt(Component, 69, [69, 70]));
    }
  });

  it('el frontend tiene pruebas de componente de los medidores con los valores de borde 69 y 70', () => {
    const files = frontendTestFiles();
    const covering = files.filter((file) => {
      const source = readFileSync(file, 'utf8');
      return /features\/dashboard|Gauge|Medidor|Meter/i.test(source) && /\b69\b/.test(source) && /\b70\b/.test(source);
    });
    expect(covering, `ninguna prueba del frontend (${files.length} archivos *.test/*.spec) cubre los medidores con 69 y 70`).not.toHaveLength(0);
  });
});
