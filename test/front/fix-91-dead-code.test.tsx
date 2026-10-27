import fs from 'node:fs';
import path from 'node:path';
import { describe, expect, it } from 'vitest';

// FIX-91 (documentacion/actual.md):
// Eliminar el componente muerto InstanceTester (FRN-15 / FRN-16).
// Problema: InstanceTester.tsx define INITIAL_INSTANCES con datos mock (ids 1..3, "servidor-web-produccion")
// y no se utiliza en producción. Debe eliminarse para que no quede código muerto.
describe('FIX-91 - eliminación de componente muerto InstanceTester', () => {
  const instanceTesterPath = path.resolve(
    __dirname,
    '../../frontend/centinela/src/components/features/instances/components/InstanceTester.tsx',
  );

  it('el archivo InstanceTester.tsx ha sido eliminado del submódulo frontend', () => {
    const exists = fs.existsSync(instanceTesterPath);
    expect(exists, 'InstanceTester.tsx todavía existe como código muerto en el frontend').toBe(false);
  });

  it('ningún archivo de producción importa InstanceTester ni utiliza datos falsos de mock local', () => {
    const srcDir = path.resolve(__dirname, '../../frontend/centinela/src');
    function scanDir(dir: string): string[] {
      let results: string[] = [];
      const list = fs.readdirSync(dir);
      for (const file of list) {
        const fullPath = path.join(dir, file);
        const stat = fs.statSync(fullPath);
        if (stat.isDirectory()) {
          results = results.concat(scanDir(fullPath));
        } else if (/\.(tsx?|jsx?)$/.test(file)) {
          const content = fs.readFileSync(fullPath, 'utf-8');
          if (content.includes('InstanceTester') || content.includes('Probador de Instancias (Mock Local)')) {
            results.push(fullPath);
          }
        }
      }
      return results;
    }

    const offendingFiles = scanDir(srcDir).filter((file) => !file.endsWith('InstanceTester.tsx'));
    expect(offendingFiles, 'Hay archivos de código fuente importando o referenciando InstanceTester').toEqual([]);
  });
});

