### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!IMPORTANT]
> **Integración y verificación local (06/10/2026):** tras ejecutar las suites de pruebas automatizadas sobre backend `df320c7` y frontend `fc6f9ed`:
> - Se verificaron y promovieron a [`terminado.md`](terminado.md): `FIX-44`, `FIX-50`, `FIX-45`, `FIX-61` y `FIX-51`.
> - Se verificaron y promovieron a [`terminado-1.md`](terminado-1.md): `BAC-22`, `BAC-23A`, `BAC-25A`, `BAC-25C`, `FIX-43`, `FIX-46`, `FIX-47`, `FIX-48`, `FIX-49`, `FIX-53`, `FIX-64`, `BAC-22B`, `BAC-23B` y `FRN-16B`.
> - Las observaciones detectadas en tareas implementadas se canalizan en [`futuro-1.md`](futuro-1.md) (`FIX-62` para `BAC-22`, `FIX-63` para `BAC-23A`, `FIX-65` para `FIX-46`, y `FIX-69` para `BAC-25A`).
> - En este archivo permanecen únicamente las tareas **pendientes de implementación**:
>
> | Tarea | Área | Estado | Pruebas / Alcance |
> |---|---|---|---|
> | `FIX-54` | Backend | No implementado: columna `eliminado_en` y separación de baja lógica | `fixes_acceptance_test.go` / `admin-users-deletion.test.tsx` |
> | `FIX-55` | Backend | No implementado: alineación de baja y autenticación con identidad histórica | `fixes_acceptance_test.go` |
> | `FIX-56` | Frontend | No implementado: UI distingue suspensión de eliminación | `admin-users-deletion.test.tsx` |
> | `FIX-58` | Frontend | No implementado: imagotipo/logotipo oficial vectorizado | Inspección visual en auth / `auth-branding.test.tsx` |
> | `FIX-59` | Frontend | No implementado: copy informativo de Centinela y 2FA | Inspección visual en auth / `auth-branding.test.tsx` |
> | `FIX-60` | Frontend | No implementado: reubicar botón "Crear instancia" a `/instances` | `dashboard-node.test.tsx` / `instances-table.test.tsx` |
> | `FRN-15` | Frontend | No implementada: modales de confirmación antierror | `instances-modals.test.tsx` (6) |
> | `FRN-19B` | Frontend | No implementada: semáforo de salud global e integración con `/api/node/status` | `dashboard-node.test.tsx` (11) |
> | `FRN-16` | Frontend | No implementada: máquina de estados "Operación en progreso" | `instances-operation.test.tsx` (9) |
> | `FIX-52` | Backend | No implementado: distinguir `INSTANCE_BUSY` de `INSTANCE_INVALID_STATE` | `fixes_acceptance_test.go` |
> | `BAC-25B` | Backend | Ola 2 (Pendiente): timeout configurable, reintentos y exitstatus | `etapa1_ola2_acceptance_test.go` |
> | `FRN-20B` | Frontend | Ola 2 (Parcial): selector reactivo por tipo (VM/LXC); buscador y contador listos | `instances-ola2.test.tsx` |

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

### `FIX-61` - Cerrar sesión y detener reconexiones ante 401 al solicitar el ticket de eventos (`FRN-17C`) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 0,5 h
- **Depende de:** `FRN-17C`, `BAC-21C`.
- **Criterio de éxito:** Ante 401 en ticket se detienen los reintentos, se cancelan los timers y se cierra la sesión.



### `FIX-58` - Reemplazar placeholder textual del logo por el imagotipo/logotipo oficial vectorizado (`FRN-01` / `LOGIN-02`) (Frontend)

- **Área:** Frontend / UI
- **Asignada:** Belinda / Luz
- **Estado:** Pendiente
- **Estimación:** 1 h
- **Depende de:** `FRN-01` y `LOGIN-02`.
- **Problema y evidencia:**
  En la vista de autenticación desplegada (`centinela.tail6bb3f3.ts.net/two-factor/verify`, front `d47999d`), en el encabezado superior izquierdo de la columna lateral se renderiza el texto plano `"logo Centinela"` (con `"logo"` estilizado en color primario/turquesa y `"Centinela"` en texto oscuro). No se está utilizando el imagotipo o logotipo vectorizado (SVG) de Centinela ni un componente de marca unificado, dejando un aspecto de maqueta incompleta en producción/staging.
- **Dónde corregir:**
  Componente del layout de autenticación en el frontend (ej. `AuthLayout.tsx` o componente de marca compartido en `src/components/layout/` o `src/pages/Login.tsx` / `LoginContinuation.tsx`).
- **Entregables:**
  1. Incorporar el activo gráfico oficial vectorizado de Centinela (archivo SVG o componente React `<CentinelaLogo />` / `<img src="/assets/logo.svg" alt="Centinela" />`).
  2. Implementar accesibilidad correcta: etiqueta `alt="Centinela"` o `aria-label="Centinela - Panel de Orquestación"` y contraste visual adecuado según los lineamientos de accesibilidad WCAG.
  3. Asegurar comportamiento responsivo: tamaño proporcional que no desborde en pantallas móviles ni en escritorio, manteniendo alineación con el contenedor principal.
  4. Si el logo incluye enlace, dirigir al inicio o a `/login` sin provocar bucles de redirección en usuarios no autenticados.
- **Criterio de éxito:**
  - En `/login`, `/two-factor/verify`, `/two-factor/setup`, `/recover-password` y `/change-password` se visualiza el logotipo oficial vectorizado de Centinela en lugar de la cadena de texto `"logo Centinela"`.
  - El elemento es accesible para lectores de pantalla y mantiene proporciones nítidas en cualquier resolución.


### `FIX-70` - Blindaje transaccional de la migración de auditoría y saneamiento de integridad en UUIDv7 (`BAC-18B` / `FIX-51`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra / Backend
- **Estimación:** 1 h
- **Depende de:** `FIX-51` (en `terminado.md`).
- **Contexto y razón de ser de `FIX-51`:**
  `FIX-51` fue indispensable porque al crear las particiones en `BAC-18B`, la tabla plana original fue renombrada a `auditoria_legacy` dejando 167 registros históricos inaccesibles para la API (`GET /api/admin/audit`) y las exportaciones. El volcado masivo `INSERT INTO auditoria (...) SELECT ... FROM auditoria_legacy ON CONFLICT DO NOTHING` implementado en `df320c7` **resolvió con éxito la recuperación de datos históricos y se preserva íntegramente**.
- **Problema detectado (oportunidad de mejora y complemento sobre `df320c7`):**
  1. **Falta de transaccionalidad atómica:** El traspaso masivo de datos (`INSERT ... SELECT`) y el descarte de la tabla anterior (`DROP TABLE auditoria_legacy CASCADE`) se ejecutan como sentencias separadas sin un bloque de transacción (`db.Transaction`). Si el proceso se interrumpe a mitad de camino, la base de datos puede quedar en un estado parcialmente migrado.
  2. **Llamada espurio de secuencia silenciada:** Se incorporó una llamada forzada a `SELECT setval(pg_get_serial_sequence('auditoria', 'id'), ...)` ignorando el error con `_ = db.Exec(...)`. Dado que `auditoria.id` es `UUIDv7` (`models.go:95`), no existen secuencias numéricas en PostgreSQL (`pg_get_serial_sequence` es `NULL`); la consulta falla internamente y el filtro regex `~ '^[0-9]+$'` sobre UUIDs introduce ruido y lógica muerta.
  3. **Riesgo de la partición `auditoria_default`:** Las particiones están definidas estáticamente hasta `2027_q1`. Si los eventos caen en la partición `default`, PostgreSQL bloquea la creación de particiones declarativas posteriores para ese rango hasta que la tabla default sea vaciada manualmente.
- **Entregables:**
  1. **Transaccionalidad ACID:** Envolver el volcado de `auditoria_legacy` y el posterior `DROP TABLE` dentro de una transacción explícita (`tx := db.Begin()`), asegurando que la migración sea atómica (todo o nada).
  2. **Saneamiento de `db.go`:** Remover las líneas 295 a 298 que intentan ajustar secuencias inexistentes con errores silenciados, garantizando un flujo declarativo y limpio acorde a claves primarias UUIDv7.
  3. **Mecanismo de aprovisionamiento de particiones futuras:** Implementar una guarda o función en el arranque (`asegurarParticionesFuturas`) que garantice la existencia de la partición del trimestre en curso y del siguiente, evitando que eventos operativos caigan en la partición `default`.
  4. **Preservación de pruebas:** Mantener y ampliar `TestMigrarAuditoriaParticionada_TraspasoLegacy` en `migrar_auditoria_test.go` verificando que ante un fallo simulado durante el traspaso, la transacción haga rollback sin perder datos de `auditoria_legacy`.
- **Criterio de éxito:**
  - El traspaso de datos históricos es 100% transaccional y atómico.
  - La migración corre sin consultas fallidas ni errores silenciados en los logs.
  - El esquema de particiones cuenta con mecanismo preventivo para evitar el bloqueo por datos en la partición `default`.
  - La API de auditoría conserva íntegramente la visibilidad del historial sin regresiones.

---
# ETAPA 1
> Las tareas de la Etapa 1 ya verificadas están en [`terminado-1.md`](terminado-1.md), y sus correcciones en [`futuro-1.md`](futuro-1.md).
---

> Tareas de la Ola 1 de [`etapa1.md`](etapa1.md) pendientes de desarrollo. Las tareas implementadas (`BAC-22`, `BAC-23A`, `BAC-25A`, `BAC-25C` y `FIX-43`) fueron verificadas y promovidas a [`terminado-1.md`](terminado-1.md). Sus observaciones técnicas asociadas se encuentran en [`futuro-1.md`](futuro-1.md).

---


#### `FRN-15` - Modales de confirmación antierror para acciones operativas
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.5 h
- **Depende de:**
  - `FRN-20A`;
  - `SEC-03` (`PermissionGate`);
  - `FIX-37` y `FIX-38`, para que `canOperateInstance` funcione con datos reales.

  Los modales se pueden maquetar antes, con datos de prueba.
- **Entregable:** modales según la criticidad de la acción, con `components/ui/modals.tsx`:
  - `Start`: confirmación estándar.
  - `Shutdown` / `Reboot`: aviso de apagado o reinicio del sistema operativo huésped.
  - `Stop`: advertencia en rojo sobre posible pérdida de datos.
  - `Delete`: modal destructivo que pide tipear el ID o el nombre. **Solo visible para ADMIN** (`PermissionGate`).

  Los botones de energía se muestran solo si `canOperateInstance(vmid)` da `true`.
- **Criterio de éxito:**
  - Ninguna acción se dispara sin pasar por el modal, y cancelar no envía peticiones.
  - Un OPERATOR con `READ_ONLY` no ve botones de energía; uno con `FULL_ACCESS` sí.
  - El OPERATOR nunca ve Delete.

#### `FRN-19B` (ex `FRN-13B`) - Semáforo de salud global e integración con `GET /api/node/status` (`RF-02`)
- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2.0 h
- **Depende de:** `FRN-19A`. Empieza con `BAC-29` (datos de prueba) y cierra con `BAC-22`.
- **Entregable:** el Dashboard conecta los medidores a `GET /api/node/status` mediante un servicio en `features/dashboard/services`:
  - uptime legible (días, horas, minutos);
  - semáforo de salud (D3):
    - `Saludable` si CPU, RAM y disco están por debajo del 70 %;
    - `Advertencia` si alguno llega al 70 % o más;
    - `Inaccesible` si el backend responde `502`/`504` o no responde.
  - Aviso de "datos desactualizados" si `stale: true`;
  - consulta automática cada 10 s, que se pausa con la pestaña oculta;
  - skeletons mientras carga y reintento visual ante una desconexión.
- **Criterio de éxito:** los datos reales del nodo se ven en pantalla y se actualizan sin `F5`. Si el backend cae, la UI muestra `Inaccesible` sin romperse.


#### `FRN-16` - Máquina de estados "Operación en progreso" por instancia
- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 2.5 h
- **Depende de:** `FRN-15`. Empieza con `BAC-29` y cierra con `FIX-39`, `FIX-40`, `BAC-24A` y `BAC-24B`.
- **Entregable:**
  - Al confirmar el modal, enviar la orden de energía (con la ruta documentada en `BAC-29`) o `DELETE`, y guardar el `tareaId` del `202` en el estado de la fila (`transitioning`).
  - El botón accionado muestra un spinner, y se deshabilitan todos los botones de esa instancia.
  - Ante un error, desbloquear la fila y mostrar un mensaje distinto según el código (D2), con los significados de `BAC-29`:
    - `409 INSTANCE_INVALID_STATE`: la acción no corresponde al estado actual (por ejemplo, "La instancia ya está encendida").
    - `409 INSTANCE_BUSY`: "La instancia está ejecutando otra tarea. Esperá a que termine".
    - `403 INSTANCE_PROTECTED`: la instancia es de infraestructura y no admite esa acción.
    - `403 INSTANCE_ACCESS_DENIED`: no tenés permiso sobre la instancia.
    - `502 PROXMOX_UNAVAILABLE`: Proxmox no está disponible.
    - `504 PROXMOX_TIMEOUT`: Proxmox no respondió a tiempo, y la acción puede no haberse aplicado.
    - Cualquier otro código: un mensaje genérico.
- **Criterio de éxito:** es imposible disparar una segunda acción sobre la misma instancia mientras hay una orden en curso. Cada uno de los seis códigos muestra su propio mensaje, verificado con una prueba de componente por código.

---

### `FIX-52` - Distinguir `INSTANCE_BUSY` de `INSTANCE_INVALID_STATE` y limpiar código muerto (`BAC-24A`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro (autor de `8591e90`)
- **Estimación:** 1,5 h
- **Depende de:** `BAC-24A` (en `terminado-1.md`).
- **Problema y evidencia (detectado en la revisión manual de Lucas del 05/10/2026, verificada en su PC):**
  1. **`INSTANCE_BUSY` quedó casi inalcanzable.** `validarEstadoParaAccion` (`instance_handler.go:111`) decide únicamente con el `estado` que devuelve `ObtenerInstancia`, comparado contra `estadosRequeridosPorAccion`. No consulta si ya existe una tarea `RUNNING` para ese `vmid` en `tareas_asincronas`. Un `start` seguido de inmediato por un `stop` sobre la misma instancia da `409 INSTANCE_INVALID_STATE` ("estado incompatible"), cuando la causa real es que la instancia está ocupada con la tarea del `start` anterior: Proxmox todavía no reflejó el cambio de estado. El código para distinguir ambos casos ya existe (`tareas_asincronas`, campo `activeTask` de `GET /api/instances`), pero `validarEstadoParaAccion` no lo usa.
  2. **Código muerto.** `ProxmoxPort.ReiniciarInstancia` (`client.go:291`, `ports/proxmox_port.go:99`) nunca se invoca: `CambiarEstado` reemplazó su uso por `Reboot(node, vmid, vmType)` (`instance_handler.go:485`), que sí resuelve nodo y tipo. `ReiniciarInstancia` quedó huérfano en el puerto y en el cliente.
  3. **Comentario desactualizado.** El comentario de `EliminarInstancia` (`client.go:323`) dice *"si está encendida Proxmox responde 500 con un mensaje de bloqueo y se retorna `ErrInstanciaOcupada`"*, pero el handler (`instance_handler.go`, validación de estado previa al `DELETE`) ya rechaza una instancia no `stopped` con `409` antes de invocar a `proxmox.EliminarInstancia`: esa rama de `ErrInstanciaOcupada` es inalcanzable con el flujo actual y el comentario induce a error a quien lea el cliente.
- **Entregable:**
  1. Antes de aplicar `estadosRequeridosPorAccion`, consultar si el `vmid` tiene una tarea `RUNNING` en `tareas_asincronas`; si la tiene, responder `409 INSTANCE_BUSY` sin tráfico de escritura a Proxmox, sin evaluar la matriz de estados.
  2. Eliminar `ReiniciarInstancia` de `ProxmoxPort` y de `client.go` (o, si se prefiere conservarlo, usarlo de forma consistente en lugar de `Reboot` y actualizar el puerto para que reciba `node`/`vmType`).
  3. Corregir el comentario de `EliminarInstancia` para reflejar que el `409` por instancia encendida lo decide el handler, no Proxmox.
- **Criterio de éxito:**
  - Un `start` inmediatamente seguido de un `stop` (u otra combinación con una tarea `RUNNING` en curso) responde `409 INSTANCE_BUSY`, no `409 INSTANCE_INVALID_STATE`.
  - `start` sobre una instancia realmente `running` (sin tarea en curso) sigue respondiendo `409 INSTANCE_INVALID_STATE`.
  - No quedan referencias a `ReiniciarInstancia` sin uso, y el comentario de `EliminarInstancia` describe el flujo real.
  - Siguen en verde todos los casos de `BAC-24A…` y `BAC-24B…` de `etapa1_acceptance_test.go`.

#### `BAC-25B` - Timeout configurable, reintentos y `exitstatus` en el evento
- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 1.5 h
- **Depende de:** `BAC-25A`.
- **Entregable:**
  1. Cortar el seguimiento a los **3 minutos** (D4): `UPID_TIMEOUT` con valor por defecto `3m`, en lugar de los 10 min fijos de `limiteSeguimiento`. Al vencer, la tarea queda `FAILED` con `motivo: TIMEOUT`.
  2. Reintentos con retroceso cuando falla la consulta del estado, en lugar del reintento fijo cada 1 s.
  3. Completar los `detalles` del `TASK_FINISHED` según `BAC-29`:
     - agregar `exitstatus`;
     - agregar `motivo` cuando la tarea falla (D2): `PROXMOX_ERROR` si Proxmox la terminó con un `exitstatus` distinto de `OK`, y `TIMEOUT` si venció el plazo. Cada caso con su propio `mensaje`;
     - agregar los textos de `shutdown`, `reboot` y `delete` en `nombresAccion`.
  4. Publicar el resultado por el bus existente (`eventosService.Publicar`) y exponerlo para `BAC-27` y `BAC-24C`, por ejemplo con un callback o un suscriptor interno.
- **Criterio de éxito:**
  - Sin `UPID_TIMEOUT` configurado, una tarea que no termina se marca `FAILED` con `motivo: TIMEOUT` a los 3 minutos.
  - Todo `TASK_FINISHED` trae `tareaId`, `accion`, `estado` y `exitstatus`. Los fallidos traen además `motivo`, y es distinto en cada uno de los dos casos.
  - `BAC-27` y `BAC-24C` reciben todos los resultados.


#### `FRN-20B` (ex `FRN-14B`) - Filtros reactivos por tipo, estado y buscador dinámico
- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2.0 h
- **Depende de:** `FRN-20A`.
- **Entregable:** barra de control sobre la tabla de inventario:
  - búsqueda en tiempo real por ID o nombre;
  - filtro por tipo (Todas, VM, LXC);
  - filtro por estado (Todas, En ejecución, Detenidas);
  - contador de elementos mostrados sobre el total.
- **Criterio de éxito:** el filtrado es instantáneo, en memoria del cliente, sin peticiones al backend.

#### `FIX-48` - Documentar `INSTANCE_INVALID_STATE` en Swagger (`BAC-24A`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas
- **Estimación:** 0,5 h
- **Depende de:** `BAC-24A`.
- **Criterio de éxito:** Swagger documenta `INSTANCE_INVALID_STATE` en acciones de ciclo de vida.


### `FIX-62` - Heurística de circuit breaker y recuperación de telemetría en `nodo_service` (`BAC-22`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra / Lisandro
- **Estado:** Pendiente en el submódulo `backend`. La suite de tests de este repositorio ya adecuó la tolerancia de sincronización de recuperación, pero el servicio en el backend presenta una anomalía de diseño cuando no existe caché previa.
- **Estimación:** 1 h
- **Depende de:** `BAC-22` (en `terminado-1.md`).
- **Problema técnico detallado:**
  En `backend/internal/core/services/nodo_service.go:24`, se definió `pausaTrasFallaNodo = 5 * time.Second` y en la línea 64:
  ```go
  if err := s.fallaReciente(); err != nil {
      return s.respaldo(ctx, err)
  }
  ```
  `s.respaldo` busca la clave `claveNodoUltimo` (`"node:status:last_known"`) en Redis. Si dicha clave no existe (por ejemplo, arranque en frío o caída de Proxmox antes de registrar lecturas exitosas), `s.respaldo` retorna `nil, false, errProxmox` (HTTP 502 inmediato). Esto genera que durante 5 segundos continuos cualquier petición sea rechazada a ciegas sin siquiera intentar consultar Proxmox, incluso si Proxmox ya se recuperó de inmediato. Un circuit breaker solo debe evitar llamadas si tiene datos *stale* que servir como fallback, o debe implementar un estado *half-open* que permita sondear la recuperación.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/core/services/nodo_service.go`, en la función `ObtenerEstado`, condicionar la llamada a `fallaReciente()`: si no existe lectura previa en `claveNodoUltimo`, **no bloquear las peticiones entrantes durante 5 segundos a ciegas**; permitir el intento hacia Proxmox o reducir la penalización sin caché.
  2. Implementar un mecanismo de sondeo o reintento inmediato tras la expiración del cooldown que limpie `ultimaFalla` y `errFalla` al obtener una respuesta exitosa (`200 OK`).
  3. Asegurar que `TestNodoService` en `nodo_service_test.go` cubra el escenario de recuperación tras fallo sin caché previa y con caché previa.
- **Criterio de éxito:**
  - El backend se recupera inmediatamente tras el restablecimiento del hipervisor sin quedar en un blackout de 5 segundos cuando no hay caché previa disponible.
  - Pasan las pruebas unitarias de `nodo_service_test.go` y la prueba de integración `TestEtapa1/BAC-22`.


### `FIX-70` - Blindaje transaccional de la migración de auditoría y saneamiento de integridad en UUIDv7 (`BAC-18B` / `FIX-51`) (Backend)

- **Área:** Backend
- **Asignada:** Tayra / Backend
- **Estimación:** 1 h
- **Depende de:** `FIX-51` (en `terminado.md`).
- **Contexto y razón de ser de `FIX-51`:**
  `FIX-51` fue indispensable porque al crear las particiones en `BAC-18B`, la tabla plana original fue renombrada a `auditoria_legacy` dejando 167 registros históricos inaccesibles para la API (`GET /api/admin/audit`) y las exportaciones. El volcado masivo `INSERT INTO auditoria (...) SELECT ... FROM auditoria_legacy ON CONFLICT DO NOTHING` implementado en `df320c7` **resolvió con éxito la recuperación de datos históricos y se preserva íntegramente**.
- **Problema detectado (oportunidad de mejora y complemento sobre `df320c7`):**
  1. **Falta de transaccionalidad atómica:** El traspaso masivo de datos (`INSERT ... SELECT`) y el descarte de la tabla anterior (`DROP TABLE auditoria_legacy CASCADE`) se ejecutan como sentencias separadas sin un bloque de transacción (`db.Transaction`). Si el proceso se interrumpe a mitad de camino, la base de datos puede quedar en un estado parcialmente migrado.
  2. **Llamada espurio de secuencia silenciada:** Se incorporó una llamada forzada a `SELECT setval(pg_get_serial_sequence('auditoria', 'id'), ...)` ignorando el error con `_ = db.Exec(...)`. Dado que `auditoria.id` es `UUIDv7` (`models.go:95`), no existen secuencias numéricas en PostgreSQL (`pg_get_serial_sequence` es `NULL`); la consulta falla internamente y el filtro regex `~ '^[0-9]+$'` sobre UUIDs introduce ruido y lógica muerta.
  3. **Riesgo de la partición `auditoria_default`:** Las particiones están definidas estáticamente hasta `2027_q1`. Si los eventos caen en la partición `default`, PostgreSQL bloquea la creación de particiones declarativas posteriores para ese rango hasta que la tabla default sea vaciada manualmente.
- **Entregables:**
  1. **Transaccionalidad ACID:** Envolver el volcado de `auditoria_legacy` y el posterior `DROP TABLE` dentro de una transacción explícita (`tx := db.Begin()`), asegurando que la migración sea atómica (todo o nada).
  2. **Saneamiento de `db.go`:** Remover las líneas que intentan ajustar secuencias inexistentes con errores silenciados, garantizando un flujo declarativo y limpio acorde a claves primarias UUIDv7.
  3. **Mecanismo de aprovisionamiento de particiones futuras:** Implementar una guarda o función en el arranque (`asegurarParticionesFuturas`) que garantice la existencia de la partición del trimestre en curso y del siguiente, evitando que eventos operativos caigan en la partición `default`.
  4. **Preservación de pruebas:** Mantener y ampliar `TestMigrarAuditoriaParticionada_TraspasoLegacy` en `migrar_auditoria_test.go` verificando que ante un fallo simulado durante el traspaso, la transacción haga rollback sin perder datos de `auditoria_legacy`.
- **Criterio de éxito:**
  - El traspaso de datos históricos es 100% transaccional y atómico.
  - La migración corre sin consultas fallidas ni errores silenciados en los logs.
  - El esquema de particiones cuenta con mecanismo preventivo para evitar el bloqueo por datos en la partición `default`.
  - La API de auditoría conserva íntegramente la visibilidad del historial sin regresiones.