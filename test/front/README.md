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
| `FRN-05` | `admin-users.test.tsx` | Guard de rol, consumo de `GET /api/users` y datos reales en la tabla (actualmente en rojo) |
| `FRN-06` | `admin-users.test.tsx` | Selector de rol, `POST /api/users` y contraseña temporal visible (actualmente en rojo) |
| `FRN-06B` | `admin-users.test.tsx` | Acción de edición por usuario (actualmente en rojo) |
| `FRN-07` | `admin-users.test.tsx` | Criterios preparados para cargar instancias, seleccionar VMIDs y guardar permisos (`todo`: no hay componente aún) |
| `FRN-08` | `admin-users.test.tsx` | Criterios preparados para Bearer, respuesta visual `403` y conservación de sesión (`todo`: no hay componente aún) |

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

Los criterios de `FRN-05`, `FRN-06` y `FRN-06B` ya son pruebas reales (no `it.todo`) que renderizan [Users.jsx](../../frontend/centinela/src/pages/Users.jsx) y [CrearUsuarios.tsx](../../frontend/centinela/src/pages/CrearUsuarios.tsx). Actualmente **fallan las 7**, porque esas páginas son maquetas estáticas sin guard de rol, sin llamadas a la API y sin selector de rol ni envío de formulario. No se deben volver a declarar como `todo`: son la señal confiable de que `FRN-05`/`FRN-06`/`FRN-06B` de `documentacion/actual.md` todavía no están completas.

Los criterios de `FRN-07` y `FRN-08` sí siguen como `it.todo`: a diferencia de los anteriores, no existe ningún componente, ruta ni servicio en el frontend contra el cual escribir una aserción real (ni selector de instancias ni interceptor de errores). El hito asociado queda agrupado en el mismo archivo para probar el recorrido administrador -> asignación -> operador -> rechazo en cuanto exista una primera versión.

El resultado de referencia de la primera ejecución está documentado en [`RESULTADOS.md`](RESULTADOS.md).
