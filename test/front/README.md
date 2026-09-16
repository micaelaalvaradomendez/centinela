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
| `FRN-05` | `admin-users.test.tsx` | Criterios preparados para ruta protegida, listado, filtros y RBAC visual |
| `FRN-06` | `admin-users.test.tsx` | Criterios preparados para alta, contraseña temporal, desactivación y errores |
| `FRN-06B` | `admin-users.test.tsx` | Criterios preparados para edición de usuario, rol, estado y confirmación |
| `FRN-07` | `admin-users.test.tsx` | Criterios preparados para cargar instancias, seleccionar VMIDs y guardar permisos |
| `FRN-08` | `admin-users.test.tsx` | Criterios preparados para Bearer, respuesta visual `403` y conservación de sesión |

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

Los criterios de `FRN-05`, `FRN-06` y `FRN-06B` están declarados como `it.todo` porque la ruta `/admin/users` y su servicio todavía no existen en el frontend. Al implementarlos, cada pendiente debe convertirse en una prueba ejecutable conservando el contrato actual del backend: `/api/users` y `/api/roles`.

Los criterios de `FRN-07` y `FRN-08` también están declarados como `it.todo`: requieren `GET /api/instances`, el endpoint definitivo de permisos y una superficie visual para el error `403`. El hito asociado queda agrupado en el mismo archivo para probar el recorrido administrador -> asignación -> operador -> rechazo.

El resultado de referencia de la primera ejecución está documentado en [`RESULTADOS.md`](RESULTADOS.md).
