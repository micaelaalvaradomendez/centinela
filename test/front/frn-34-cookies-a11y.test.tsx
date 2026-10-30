import fs from 'node:fs';
import path from 'node:path';
import { describe, expect, it } from 'vitest';

// FRN-34 (Etapa 3) — Adaptación del cliente a cookies seguras de sesión y auditoría de accesibilidad.
// Metodología TDD estricta: Valida que los modales críticos de Etapa 3 (Snapshots Rollback / Confirmaciones)
// cuenten con atributos de accesibilidad WCAG 2.1 AA (role="alertdialog" / role="dialog", aria-labelledby, aria-describedby)
// y que el flujo de autenticación use cookies HttpOnly sin persistir tokens de refresco en web storage.

describe('FRN-34 - Cookies seguras de sesión y accesibilidad en Frontend', () => {
  const apiClientPath = path.resolve(__dirname, '../../frontend/centinela/src/services/apiClient.ts');
  const tokenStoragePath = path.resolve(__dirname, '../../frontend/centinela/src/storage/tokenStorage.ts');
  const snapshotsDir = path.resolve(__dirname, '../../frontend/centinela/src/components/features/instances');

  it('apiClient.ts configura credentials: "include" para transmitir cookies HttpOnly', () => {
    let content = '';
    if (fs.existsSync(apiClientPath)) {
      content = fs.readFileSync(apiClientPath, 'utf-8');
    }
    expect(content, 'apiClient debe configurar credentials: "include"').toMatch(/credentials:\s*['"]include['"]/);
  });

  it('el cliente no persiste refreshToken en sessionStorage ni localStorage', () => {
    if (fs.existsSync(tokenStoragePath)) {
      const storageContent = fs.readFileSync(tokenStoragePath, 'utf-8');
      expect(storageContent, 'tokenStorage.ts no debe almacenar el refreshToken en storage del navegador').not.toMatch(
        /setItem\([^)]*refresh/i,
      );
    }
  });

  it('los modales de snapshots y acciones destructivas de Etapa 3 implementan accesibilidad WCAG (alertdialog / aria-labelledby)', () => {
    // Busca archivos de modal de snapshots
    let rollbackModalFound = false;
    let hasWcagDialogAttributes = false;

    function walkDir(dir: string) {
      if (!fs.existsSync(dir)) return;
      for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
        const fullPath = path.join(dir, entry.name);
        if (entry.isDirectory()) {
          walkDir(fullPath);
        } else if (entry.isFile() && /snapshot.*modal|rollback.*modal/i.test(entry.name)) {
          rollbackModalFound = true;
          const text = fs.readFileSync(fullPath, 'utf-8');
          if (/role=["']alertdialog["']|role=["']dialog["']|aria-labelledby|aria-describedby/.test(text)) {
            hasWcagDialogAttributes = true;
          }
        }
      }
    }

    walkDir(snapshotsDir);

    expect(rollbackModalFound, 'FRN-34: Debe existir el componente modal de rollback/snapshots').toBe(true);
    expect(hasWcagDialogAttributes, 'FRN-34: El modal destructivo debe declarar atributos de accesibilidad WCAG (role="alertdialog" o aria-labelledby)').toBe(true);
  });
});
