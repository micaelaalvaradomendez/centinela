### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Integración y verificación local (09/10/2026):** tras ejecutar las suites de pruebas automatizadas sobre backend `4f68e44` y frontend `ee4644c`:
> - Se verificaron y promovieron a [`terminado.md`](terminado.md): `FIX-58`, `FIX-59` (con observación en `FIX-80`), `FIX-60`, `FIX-61` y `FIX-70`.
> - Se verificaron y promovieron a [`terminado-1.md`](terminado-1.md): `FRN-15`, `FRN-19B` (con observación en `FIX-78`), `FRN-16` (con observación en `FIX-77`), `FIX-52`, `BAC-25B`, `FRN-20B` (con observación en `FIX-79`), `FIX-48` y `FIX-62`.
> - Las observaciones detectadas en tareas implementadas se canalizan en [`futuro.md`](futuro.md) (`FIX-80`) y [`futuro-1.md`](futuro-1.md) (`FIX-77`, `FIX-78`, `FIX-79`).
> - En este archivo permanecen únicamente las tareas **pendientes de implementación**:
>
> | Tarea | Área | Estado | Pruebas / Alcance |
> |---|---|---|---|
> | `FIX-54` | Backend | No implementado: columna `eliminado_en` y separación de baja lógica | `fixes_acceptance_test.go` / `admin-users-deletion.test.tsx` |
> | `FIX-55` | Backend | No implementado: alineación de baja y autenticación con identidad histórica | `fixes_acceptance_test.go` |
> | `FIX-56` | Frontend | No implementado: UI distingue suspensión de eliminación | `admin-users-deletion.test.tsx` |

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

> **Cobertura incorporada del PR #3:** las pruebas nuevas comprueban el comportamiento actual: conflicto al reutilizar correo o username de un usuario dado de baja, y alta con datos distintos. No implementan ni validan la reutilización del correo requerida por `FIX-54` a `FIX-57`; esa regresión debe cambiar cuando se implemente el nuevo contrato.

> **Nota de alcance (06/10/2026, revisión estática):** “recorrido completo” se refiere a los casos anteriores, no a crear una cuenta nueva con el correo de una eliminada y preservar la identidad histórica. Esta brecha de la fase base queda pendiente en `FIX-54` a `FIX-57` de [`futuro.md`](futuro.md#fix-54---separar-eliminación-lógica-y-suspensión-migrar-unicidad-del-correo-backend). No se ejecutaron suites ni se verificó el esquema desplegado.

---

### `FIX-54` - Separar eliminación lógica y suspensión; migrar unicidad del correo (Backend)

- **Área:** Backend
- **Estado:** Pendiente; no implementado ni probado.
- **Estimación:** 2 h
- **Depende de:** `BAC-05`, `BAC-06` y la clasificación controlada de datos históricos descrita arriba; coordinar disponibilidad de evidencia con `FIX-51`, sin mezclar su solución.
- **Problema:** la baja conserva correctamente la identidad, pero la unicidad incondicional reserva su correo; `activo` no distingue eliminación de suspensión.
- **Entregables:**
  1. Agregar `eliminado_en timestamptz` nullable y definir la invariante de eliminado siempre inactivo, conservando identidad y relaciones. Aplicar el plan seguro de clasificación histórica; no inferir eliminación solo desde `activo=false`.
  2. Quitar `uniqueIndex` incondicional de `EmailUsuario` del modelo e inspeccionar el nombre real del índice/constraint existente antes de reemplazarlo explícitamente por `UNIQUE (lower(btrim(email_usuario))) WHERE eliminado_en IS NULL`. Evitar que `AutoMigrate` recree la unicidad anterior; mantener la del username.
  3. Prevalidar duplicados normalizados entre no eliminados sin fusionar identidades ni modificar auditoría; definir resolución controlada y transacción/orden seguro de migración. Documentar base nueva, base existente, reinicio e idempotencia.
  4. Corregir el comentario `SET NULL` para reflejar `RESTRICT`, sin cambiar FK ni aplicar scopes de borrado a JOIN de auditoría o actividad histórica.
- **Criterios de aceptación pendientes:** múltiples eliminados pueden conservar el mismo correo; solo una fila no eliminada puede reservarlo, incluso con mayúsculas/espacios; suspendidos lo reservan. Migración repetible y reinicio no recrean el índice anterior ni pierden filas/relaciones/auditoría.

### `FIX-55` - Alinear baja, validación y autenticación con identidad histórica (Backend)

- **Área:** Backend
- **Estado:** Pendiente; no implementado ni probado.
- **Estimación:** 2,5 h
- **Depende de:** `FIX-54`, `BAC-06`, `BAC-06B`, `BAC-13`, `BAC-15` y los mecanismos existentes de autenticación/revocación.
- **Problema:** cambiar solo la unicidad dejaría lookup ambiguo, reactivación de eliminados y flujos con estado temporal anterior a la baja.
- **Entregables:**
  1. `DELETE` fija la marca y desactiva conservando identidad; `PUT activo=false` solo suspende y `PUT activo=true` no restaura eliminados. Bloquear edición, cambios de permisos y resets administrativos sobre eliminados; conservar consultas históricas explícitas por UUID.
  2. `ExisteEmailEnOrg` y validaciones de crear/editar/perfil excluyen eliminados, incluyen suspendidos y comparten normalización con DB. El lookup de autenticación por email excluye eliminados, mantiene el rechazo de suspendidos y no usa `First` ambiguo entre generaciones.
  3. Incorporar guardas de eliminado por UUID en login, 2FA, refresh y recuperación y revisar los caminos que emiten/usan acceso. Invalidar OTP de recuperación y expiración, preauth, tickets, sesiones PostgreSQL/Redis y streams con los mecanismos existentes. Definir consistencia/rollback o estrategia fail-closed: un error de revocación no debe producir éxito engañoso ni dejar acceso habilitado.
  4. Crear con UUID y clave nuevos, cambio obligatorio y 2FA nuevo, sin permisos/sesiones heredados; mantener relaciones históricas en el UUID anterior. Conservar la unicidad DB como última barrera ante carreras y mapear su violación de correo a `409` estructurado, con código acordado en contrato, sin SQL crudo y sin convertir cualquier fallo DB en conflicto.
  5. Revisar el envío SMTP previo al INSERT al probar altas concurrentes y documentar su riesgo; no prometer una transacción distribuida email/DB ni incorporar un refactor SMTP ajeno como requisito obligatorio.
- **Criterios de aceptación pendientes:** nueva cuenta accesible sin seleccionar la histórica; eliminado no puede reactivarse, editarse ni acceder mediante credenciales, tokens, preauth u OTP previos. Suspensión sigue siendo reversible y reserva correo. Conflictos y fallos de revocación tienen respuestas controladas y coherentes con el estado persistido.

### `FIX-56` - Distinguir suspensión y eliminación en contrato e interfaz (Frontend)

- **Área:** Frontend
- **Estado:** Pendiente; no implementado ni probado.
- **Estimación:** 1,5 h
- **Depende de:** DTO/contrato de `FIX-54`/`FIX-55`, `FRN-05`, `FRN-06`, `FRN-06B` y `FIX-24`.
- **Problema:** la interfaz llama eliminación a una desactivación y no dispone de estado separado; no es la causa del bloqueo de unicidad backend.
- **Entregables:**
  1. Diferenciar “Inactivo” de “Eliminado” según DTO; deshabilitar reactivar, editar y resetear eliminados. Ocultarlos del listado operativo por defecto manteniendo acceso histórico explícito y auditoría.
  2. Modal de eliminación: explicar conservación de historial y liberación del correo; suspensión: explicar reversibilidad y reserva del correo. Alinear acciones, mensajes y navegación con el contrato.
  3. Mostrar el conflicto de correo de cuentas no eliminadas con el código acordado por backend; mantener alta `201`, baja `204` y restricción `403` salvo cambio contractual explícito documentado.
- **Criterios de aceptación pendientes:** mensajes y acciones distinguen ambos estados; un eliminado no ofrece acciones operativas; historial accesible y errores de correo claros, sin atribuir la restricción DB al frontend.
