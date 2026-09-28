# Pruebas de aceptación del frontend

Estas pruebas contrastan el frontend (`frontend/centinela`, React + TS) con los criterios de éxito de `documentacion/actual.md` y funcionan como prueba de regresión de `documentacion/terminado.md`. Usan Vitest, React Testing Library y jsdom, sin agregar dependencias ni archivos al submódulo.

## Principios

- **Una prueba que pasa verifica el criterio de éxito.** Se exige la llamada HTTP real (método, ruta, cabeceras, payload) y su efecto en la interfaz. No se usan patrones del tipo `if (llamada) … else expect(boton).toBeVisible()`.
- **Los mocks de `fetch` reproducen el contrato real del backend**, tomado de `backend/internal/adapters/primary/http`. Ejemplos: `/auth/password/forgot`, `{ vmids }` en permisos, `GET /instances` con `{ id, name, type: vm|lxc }` y la ausencia de `contrasenaTemp` por BAC-16. Si el frontend espera otro contrato, la prueba falla: es un problema de integración.
- **La sesión se siembra donde la guarda la app.** `centinela_access` y `centinela_pending_login` van en sessionStorage; `centinela_user` va en localStorage.

## Cobertura

| Archivo | Tareas |
|---|---|
| `login-form.test.tsx` | FRN-01, FRN-02 |
| `authentication-contract.test.ts` | FRN-04, LOGIN-02, LOGIN-03 |
| `two-factor-form.test.tsx` | LOGIN-02 |
| `two-factor-enrollment.test.tsx` | FRN-09 |
| `api-client.test.ts` | FIX-07 |
| `navigation.test.tsx` | FRN-03, **SEC-03** |
| `admin-users.test.tsx` | FRN-05, FRN-06, FRN-06B, **FIX-14/FRN-07**, FRN-08 |
| `admin-recovery.test.tsx` | **FRN-11** |
| `recover-password.test.tsx` | **FRN-12** |
| `password-change.test.tsx` | **FRN-10** |
| `session-security.test.ts` | **FRN-13**, **SEC-02**, FIX-08 |
| `audit.test.tsx` | **FRN-14** |

## Ejecutar

```bash
pnpm --dir test/front install
pnpm --dir test/front test                 # contra el submódulo tal como está checkouteado
```

Antes de correr la suite se trae el último commit del submódulo; ante un conflicto prevalece el remoto:

```bash
git -C frontend fetch origin --prune
git -C frontend checkout -B main origin/main --force
git -C frontend reset --hard origin/main
```

`vitest.config.mjs` también acepta `CENTINELA_FRONTEND_DIR` para apuntar a otra copia de `centinela/`, pero el procedimiento estándar es correr la suite sobre el submódulo actualizado.

Si cambian las dependencias del frontend, primero hay que correr `pnpm install` en `frontend/centinela`.

## Entorno

`setup.ts` agrega los polyfills que jsdom no trae: `document.elementFromPoint` y `window.matchMedia` (este último lo usa `hooks/use-mobile.ts` del sidebar). Después de cada prueba limpia sessionStorage y localStorage.

Resultado vigente: [`RESULTADOS.md`](RESULTADOS.md). Informe consolidado: [`../informe.md`](../informe.md).
