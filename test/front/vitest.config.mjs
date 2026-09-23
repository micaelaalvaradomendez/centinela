import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';

const repositoryRoot = new URL('../..', import.meta.url).pathname;
// Por defecto se prueba el submódulo tal como está checkouteado. Para verificar otra
// revisión del equipo (por ejemplo origin/main exportado con `git archive`) sin tocar
// el submódulo: CENTINELA_FRONTEND_DIR=/ruta/a/centinela pnpm test
const frontendDirectory = process.env.CENTINELA_FRONTEND_DIR
  ?? new URL('../../frontend/centinela', import.meta.url).pathname;
const frontendSource = `${frontendDirectory.replace(/\/$/, '')}/src`;

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': frontendSource,
    },
    dedupe: ['react', 'react-dom', 'react-router'],
  },
  server: {
    fs: {
      allow: [repositoryRoot, frontendDirectory],
    },
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./setup.ts'],
    clearMocks: true,
    restoreMocks: true,
  },
});
