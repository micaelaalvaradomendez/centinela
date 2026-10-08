# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 08/10/2026 (corrida con últimos commits de submódulos)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), contrastadas tras la incorporación de los nuevos commits en el backend (`b630b2f`) y frontend (`bd11531`). Se corrieron las suites de pruebas de backend (`test/back`) y frontend (`test/front`), incorporando las nuevas pruebas de aceptación para las tareas pendientes y la verificación de regresión.

| Componente | Revisión probada | Commits nuevos incorporados |
|---|---|---|
| Backend | `b630b2f` (último commit de `main`) | `b630b2f` (Merge PR #17), `1a6015e` (tratar WARNINGS de Proxmox como tarea completada), `dc14312` (graceful shutdown con `http.Server` y `Shutdown(ctx)` para `FIX-50`), `464b494` (restauración de esquemas canónicos http.* y path pause en Swagger para `FIX-69`/`BAC-25A`) |
| Frontend | `bd11531` (último commit de `main`) | `bd11531` (Merge PR #103), `eb6da06` (avatar por rol y logo de mascota pulpo), `226539e` (fondo reactivo y logo pulpo), `fe16eaa` / `e91b70f` / `82d89a6` (unificación visual de auth y reemplazo de copy informal por plataforma y 2FA para `FIX-59`), `83d3e04` y `24315e6` (refactor de cards a tabs con filtros en Auditoría y Usuarios) |

---

## 1. Resumen Ejecutivo

- **Backend (`b630b2f`):**
  - `FIX-50` (apagado ordenado HTTP): **Completado y en verde.** Implementado en `cmd/api/main.go` con `http.Server{Addr, Handler: router}`, `srv.Shutdown(shutdownCtx)` y `<-ctx.Done()`.
  - `BAC-25A` / `FIX-69` (Swagger y pool de workers): **Completado y en verde.** Restaurados los esquemas canónicos en `docs/swagger.json` y `docs/swagger.yaml`.
  - `FIX-51` (auditoría legacy), `BAC-23B` (IP real y RBAC) y `BAC-22B` (telemetría y resumen): Se mantienen **100% en verde**.
  - **Pendientes en backend:**
    - `FIX-52`: Falta consultar tareas en estado `RUNNING` antes de evaluar estados en `instance_handler.go`, y limpiar `ReiniciarInstancia` en `proxmox_port.go`.
    - `FIX-54`: Falta agregar `eliminado_en timestamptz` en tabla `usuarios`, reemplazar el índice incondicional por parcial `WHERE eliminado_en IS NULL` normalizado con `lower(btrim(...))`, y corregir comentario en `models.go`.
    - `FIX-55`: Bloqueada por dependencia de `FIX-54`.
    - `BAC-25B`: Falta parametrizar `UPID_TIMEOUT` (3m) en `seguimiento_tareas.go`, agregar backoff en reintentos, declarar `MotivoTimeout` y completar `nombresAccion`.

- **Frontend (`bd11531`):**
  - `FIX-59` (copy de plataforma y 2FA): **Implementado en `fe16eaa`.** Se erradicó el texto informal `"imagenes y informacion random"` y se integraron las secciones descriptivas sobre la plataforma Centinela y la importancia del 2FA.
  - `FRN-16B` (resincronización de tareas activas al montar la tabla): Se mantiene **100% en verde**.
  - `FRN-20B` (filtros reactivos): **Parcial.** La búsqueda en memoria por nombre o ID numérico y el contador de instancias están implementados; resta añadir el selector por tipo de recurso (Todas, VM, LXC).
  - **Pendientes en frontend:**
    - `FIX-58`: En `MainLayoutAuth.tsx` persiste el placeholder en texto plano `<span>logo</span>` junto a `<h1> Centinela </h1>` en lugar del imagotipo/logotipo SVG oficial accesible.
    - `FIX-60`: El botón "Crear instancia" sigue ubicado en `/dashboard` en lugar de residir en la cabecera de `/instances` protegido por rol `ADMIN`.
    - `FIX-56`: Falta distinguir eliminación de suspensión en la tabla de usuarios y en los diálogos de acción.
    - `FRN-15`: Falta implementar los modales de confirmación para acciones operativas (Start, Shutdown/Reboot, Stop en rojo, Delete con tipeo).
    - `FRN-16`: Falta la máquina de estados con spinner y mapeo de los 6 códigos de error de D2.
    - `FRN-19B`: El Dashboard no consume `GET /api/node/status` ni renderiza el semáforo con estados D3 (Saludable, Advertencia, Inaccesible).

---

## 2. Estado Detallado de Tareas de `documentacion/actual.md`

**Leyenda:** ✅ Cumple · 🟡 Implementación parcial / en progreso · ❌ No implementada.

| Tarea | Área | Estado | Detalle y Cobertura de Pruebas |
|---|---|---|---|
| `FIX-50` apagado ordenado | Backend | ✅ | **Pasa 100% en verde.** `cmd/api/main.go` implementa graceful shutdown con `http.Server`, captura de señales con `signal.NotifyContext` y `srv.Shutdown(ctx)`. Validado en `fixes_acceptance_test.go:27`. |
| `FIX-51` auditoría legacy | Backend | ✅ | **Pasa 100% en verde.** Migración de datos históricos en `db.go:286` y tabla particionada activa. Validado en `fixes_acceptance_test.go:52`. |
| `BAC-23B` IP real y RBAC | Backend | ✅ | **Pasa 100% en verde.** Resolución de IP con `ResolverIPs` posterior al filtrado de permisos. Validado en `etapa1_ola2_acceptance_test.go:28`. |
| `BAC-22B` telemetría y resumen | Backend | ✅ | **Pasa 100% en verde.** `instancesSummary` (vms/lxc) en `/node/status` y métricas `cpuUsage`, `ramUsage`, `maxRam` en `/instances`. Validado en `etapa1_ola2_acceptance_test.go:74`. |
| `BAC-25A` pool de UPID y Swagger | Backend | ✅ | **Pasa 100% en verde.** Esquemas canónicos restaurados en `docs/swagger.json` y `docs/swagger.yaml` tras `464b494` (`FIX-69`). |
| `FRN-16B` resincronización F5 | Frontend | ✅ | **Pasa 100% en verde.** Filas con `activeTask !== null` inician bloqueadas en estado `transitioning`. Validado en `instances-ola2.test.tsx:125`. |
| `FIX-59` copy amigable auth | Frontend | ✅ | **Pasa 100% en verde.** Reemplazado `"imagenes y informacion random"` por las secciones descriptivas de Centinela y 2FA en `MainLayoutAuth.tsx`. Validado en `auth-branding.test.tsx:42`. |
| `FRN-20B` filtros reactivos | Frontend | 🟡 | Búsqueda reactiva en memoria por ID/nombre y contador de elementos pasan 100%. Pendiente agregar el selector reactivo por tipo de recurso (Todas, VM, LXC). Validado en `instances-ola2.test.tsx:55`. |
| `BAC-25B` timeout configurable | Backend | ❌ | Pendiente en backend: `limiteSeguimiento` fijo en 10 min sin soportar `UPID_TIMEOUT` (3m), sin backoff, sin `MotivoTimeout` y faltan textos de shutdown/reboot/delete en `nombresAccion`. Verificado en `etapa1_ola2_acceptance_test.go:133`. |
| `FIX-52` INSTANCE_BUSY | Backend | ❌ | Pendiente en backend: `validarEstadoParaAccion` no consulta si existe tarea `RUNNING` antes de evaluar estados, y persiste código muerto de `ReiniciarInstancia`. Verificado en `fixes_acceptance_test.go:78`. |
| `FIX-54` eliminación lógica | Backend | ❌ | Pendiente en backend: falta columna `eliminado_en`, reemplazo del índice único incondicional por parcial `WHERE eliminado_en IS NULL` normalizado con `lower(btrim(...))` y corrección de comentario en `models.go`. Verificado en `fixes_acceptance_test.go:155`. |
| `FIX-55` identidad histórica | Backend | ❌ | Pendiente en backend: bloqueada por dependencia con `FIX-54`. Verificado en `fixes_acceptance_test.go:175`. |
| `FIX-56` distinción UI baja | Frontend | ❌ | Pendiente en frontend: la interfaz no distingue eliminación de suspensión en badges ni modal. Verificado en `admin-users-deletion.test.tsx:104`. |
| `FIX-58` logotipo oficial | Frontend | ❌ | Pendiente en frontend: layout de autenticación renderiza texto plano `<span>logo</span>` en lugar del imagotipo/logotipo SVG accesible. Verificado en `auth-branding.test.tsx:19`. |
| `FIX-60` reubicar botón crear | Frontend | ❌ | Pendiente en frontend: `/dashboard` contiene el botón "Crear instancia", el cual debe ubicarse en `/instances` protegido por rol ADMIN. Verificado en `dashboard-node.test.tsx:147` e `instances-table.test.tsx:136`. |
| `FRN-15` modales antierror | Frontend | ❌ | Pendiente en frontend: faltan los modales de confirmación para las acciones operativas en el inventario. Verificado en `instances-modals.test.tsx:41`. |
| `FRN-16` operación en progreso | Frontend | ❌ | Pendiente en frontend: falta la máquina de estados reactiva con spinner y mapeo de errores D2. Verificado en `instances-operation.test.tsx:35`. |
| `FRN-19B` semáforo de nodo | Frontend | ❌ | Pendiente en frontend: Dashboard no consume `/node/status` ni renderiza el semáforo con estados D3. Verificado en `dashboard-node.test.tsx:57`. |

---

## 3. Pruebas Creadas y Complementadas

Se incorporaron casos de prueba faltantes para completar la cobertura de todos los entregables de `actual.md`:

1. **`test/back/fixes_acceptance_test.go`:**
   - **`FIX-52`:** Se añadió verificación de `DELETE /instances/:id` rechazado con `409 INSTANCE_BUSY` ante tarea `RUNNING` en curso, y liberación de la instancia una vez completada la tarea (`COMPLETED`).
   - **`FIX-54`:** Se agregaron aserciones de normalización `lower(btrim(email_usuario))` en el índice parcial único, verificación de unicidad incondicional sobre `nombre_usuario`, corrección del comentario `RESTRICT` en `models.go` y soporte de múltiples filas eliminadas con el mismo correo en base de datos.
   - **`FIX-55`:** Se agregaron pruebas de rechazo de edición (`PUT /admin/users/:id`) y permisos (`PUT /admin/users/:id/permissions`) sobre usuarios eliminados, consulta histórica explícita por UUID (`GET /admin/users/:id`) con `eliminadoEn`, y verificación de reversibilidad y reserva de correo de la suspensión.

2. **`test/back/etapa1_ola2_acceptance_test.go`:**
   - **`BAC-25B`:** Se agregaron comprobaciones de la constante `MotivoTimeout` / `"TIMEOUT"` en `ports/event_port.go`, presencia de textos para `"shutdown"`, `"reboot"` y `"delete"` en `nombresAccion`, y validación del esquema de `detalles` del evento `TASK_FINISHED` (`tareaId`, `accion`, `estado`, `exitstatus`, `motivo`, `error`).

3. **`test/front/auth-branding.test.tsx`:**
   - **`FIX-58`:** Se validó que el enlace del imagotipo vectorizado oficial dirija a `/` o `/login` sin crear bucles de redirección.
   - **`FIX-59`:** Se validó la presencia de los textos aprobados para los bloques de Plataforma (*"Tu infraestructura virtual, simplificada y bajo control"*) y 2FA (*"Protección de infraestructura con doble factor"*).

4. **`test/front/instances-ola2.test.tsx`:**
   - **`FRN-20B`:** Se añadieron pruebas de búsqueda reactiva por ID numérico (`102`), interacción con el selector de tipo de recurso (filtrado dinámico de `VM`, `LXC` y `Todas`), y actualización dinámica del contador de elementos sobre el total.

5. **`test/front/dashboard-node.test.tsx` e `instances-table.test.tsx`:**
   - **`FIX-60`:** Se añadió la comprobación en `dashboard-node.test.tsx` garantizando que el Dashboard no renderice el botón "Crear instancia", y en `instances-table.test.tsx` asegurando que la cabecera del inventario incluya el botón para usuarios con rol `ADMIN`.

---

## 4. Próximos Pasos para los Equipos de Desarrollo

1. **Promover tareas resueltas a `documentacion/terminado.md` / `documentacion/terminado-1.md`:**
   - `FIX-50` (Backend): completamente verificado y pasando en verde con graceful shutdown.
   - `FIX-59` (Frontend): completamente verificado tras el rediseño de autenticación en `fe16eaa`.
   - `BAC-25A` (Backend): esquemas restaurados en Swagger.

2. **Backend:**
   - Implementar `BAC-25B` en `seguimiento_tareas.go`: configurar `UPID_TIMEOUT` con valor default de 3 minutos, retroceso en reintentos y constantes de motivo.
   - Implementar `FIX-52` en `instance_handler.go`: consultar tareas activas `RUNNING` antes de evaluar compatibilidad de estado, y remover `ReiniciarInstancia` del puerto.
   - Implementar `FIX-54` y `FIX-55`: migrar esquema agregando `eliminado_en`, reemplazar índice por condicional y alinear endpoints de usuario.

3. **Frontend:**
   - Resolver `FIX-58`: reemplazar `<span>logo</span>` en `MainLayoutAuth.tsx` por el SVG oficial de Centinela.
   - Resolver `FIX-60`: reubicar el botón "Crear instancia" del Dashboard a la barra de inventario `/instances` condicionado por rol `ADMIN`.
   - Completar el selector por tipo de recurso en `Instances.tsx` (`FRN-20B`).
   - Implementar los modales antierror (`FRN-15`), máquina de estados (`FRN-16`) y semáforo del Dashboard (`FRN-19B`).
