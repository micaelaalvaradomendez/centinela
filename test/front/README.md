# Pruebas de aceptación del frontend

Esta suite contrasta las tareas frontend de `documentacion/actual.md` con el código actual de `frontend/centinela`, sin agregar dependencias ni archivos al submódulo.

## Cobertura

| Tarea | Archivo | Comportamiento validado |
|---|---|---|
| `FRN-01` | `login-form.test.tsx` | Campos, controles y estado de carga del login |
| `FRN-02` | `login-form.test.tsx` | Validación previa al envío y errores por campo |
| `FRN-03` | `navigation.test.tsx` | Protección de rutas y navegación Dashboard/Instancias |
| `FRN-04` | `authentication-contract.test.ts` | Payload y respuesta del login según la API Go |
| `LOGIN-02` | `two-factor-form.test.tsx`, `authentication-contract.test.ts` | Código de seis dígitos, errores y contrato 2FA |
| `LOGIN-03` | `authentication-contract.test.ts` | Continuidad desde login válido hacia enrolamiento 2FA |

## Ejecutar

Desde la raíz del repositorio:

```bash
corepack pnpm --dir test/front install
corepack pnpm --dir test/front test
```

Para desarrollo interactivo:

```bash
corepack pnpm --dir test/front test:watch
```

## Interpretación

Una prueba fallida indica que el criterio de la tarea no está satisfecho por el frontend actual. En particular, las pruebas de contrato usan los nombres, rutas, métodos y respuestas expuestos por el backend Go auditado en `main`; no adaptan ni simulan el contrato alternativo que actualmente espera el frontend.

El resultado de referencia de la primera ejecución está documentado en [`RESULTADOS.md`](RESULTADOS.md).
