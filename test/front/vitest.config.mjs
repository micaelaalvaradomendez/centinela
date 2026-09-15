import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';

const repositoryRoot = new URL('../..', import.meta.url).pathname;
const frontendSource = new URL('../../frontend/centinela/src', import.meta.url).pathname;

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
      allow: [repositoryRoot],
    },
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./setup.ts'],
    clearMocks: true,
    restoreMocks: true,
  },
});
