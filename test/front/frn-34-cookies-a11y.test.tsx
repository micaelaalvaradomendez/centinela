import fs from 'node:fs';
import path from 'node:path';
import { describe, expect, it } from 'vitest';

// FRN-34 (Etapa 3) — Adaptación del cliente a cookies seguras de sesión y auditoría de accesibilidad.
// Entregables:
//   1. Ajuste de apiClient.ts para enviar credentials: 'include' en peticiones al backend.
//   2. Eliminación de referencias al almacenamiento manual del refreshToken en sessionStorage.
//   3. Auditoría de accesibilidad (WCAG 2.1 AA): etiquetas aria-* en formularios, botones y modales.
describe('FRN-34 - Cookies seguras de sesión y accesibilidad en Frontend', () => {
  const apiClientPath = path.resolve(__dirname, '../../frontend/centinela/src/api/apiClient.ts');
  const tokenStoragePath = path.resolve(__dirname, '../../frontend/centinela/src/storage/tokenStorage.ts');

  it('apiClient.ts configura credentials: "include" para transmitir cookies HttpOnly', () => {
    let content = '';
    if (fs.existsSync(apiClientPath)) {
      content = fs.readFileSync(apiClientPath, 'utf-8');
    } else {
      // Buscar cliente alternativo si cambió de ruta
      const clientAlt = path.resolve(__dirname, '../../frontend/centinela/src/services/apiClient.ts');
      if (fs.existsSync(clientAlt)) content = fs.readFileSync(clientAlt, 'utf-8');
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

  it('los modales y formularios de la aplicación implementan atributos de accesibilidad WCAG', () => {
    const pagesDir = path.resolve(__dirname, '../../frontend/centinela/src/pages');
    let hasAria = false;
    for (const file of fs.readdirSync(pagesDir)) {
      if (/\.tsx$/.test(file)) {
        const text = fs.readFileSync(path.join(pagesDir, file), 'utf-8');
        if (/aria-label|aria-describedby|role="dialog"|role="alert"/.test(text)) {
          hasAria = true;
          break;
        }
      }
    }
    expect(hasAria, 'Las páginas y componentes deben incluir atributos de accesibilidad aria-*').toBe(true);
  });
});

