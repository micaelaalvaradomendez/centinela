
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 01/10/2026, 2ª verificación** (backend `44a2339`, frontend `3d1e84a`; detalle en [test/informe.md](../test/informe.md)).
> - `FIX-35` (Etapa 1) está completo y pasó a [`terminado-1.md`](terminado-1.md).
> - `BAC-21B` está implementada con problema y pasó a [`terminado.md`](terminado.md). Lo que falta, junto con la aclaración del 01/10/2026, está en `FIX-39`, en [`futuro.md`](futuro.md).
> - Ninguna de las tareas que quedan en este archivo tuvo commits en esta verificación. Todas tienen pruebas automatizadas, y esas pruebas hoy fallan.
>
> | Tarea | Área | Estado | Pruebas |
> |---|---|---|---|
> | `SEC-03` | Frontend | Parcial: `usePermissions` existe en `src/hooks/`; faltan `isOperator`, `hasRole`, `PermissionGate` y el menú | `navigation.test.tsx` (3) |
> | `FIX-29` | Frontend | No implementado: la mayúscula, el mensaje de dígito y el mensaje de largo | `password-change` (3), `recover-password` (1) |
> | `FRN-17C` | Frontend | No implementada: no existe `useEvents` | `events-client.test.tsx` (4) |
> | `FIX-36` | Frontend | No implementado: el filtro "Acción" de Auditoría | `audit.test.tsx` (1) |
> | `FIX-38` | Frontend | No implementado: el selector de rol ofrece `READ_ONLY`, y el hook lo acepta como rol | `admin-users` (1), `navigation` (1) |
> | `BAC-18B` | Backend | No implementada: faltan el índice parcial, las particiones y la purga horaria. El índice `(usuario_id, fecha_hora)` ya existe (`idx_auditoria_usuario_fecha`) | `cierre_fase_base…` (2) |
> | `FIX-37` | Backend | No implementado: falta `permisos` en `GET /account/profile` | `resource_access…` (1) |

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

---

### `SEC-03` - Contexto y sistema reactivo de permisos en Frontend (UI/UX RBAC + Resources)

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 2,5 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 2).
- **Depende de:** `FRN-04`, `BAC-07` y `BAC-09`.
- **Estado verificado (01/10/2026, frontend `3d1e84a`):** en proceso, parcial. `test/front/navigation.test.tsx`: fallan los 3 casos de SEC-03 (el de `FIX-30` pasa).
  - Existe `usePermissions()` en `src/hooks/usePermissions.ts`, con `isAdmin`, `canAccessInstance` y `canOperateInstance`. La carpeta es válida: el criterio no la fija y la prueba la acepta.
  - Faltan `isOperator`, `hasRole` y `PermissionGate`.
  - `Sidebar.tsx` sigue mostrando "Auditoría" al OPERATOR.
  - `READ_ONLY` tratado como rol (en el hook y en el selector de rol de la ficha): ver `FIX-38`, más abajo en este archivo.
- **Estado anterior (23/09/2026):** no implementada.
  - `src/context/AuthContext.js` está vacío.
  - No existen `usePermissions` ni `PermissionGate`.
  - El menú solo filtra con un `isAdmin` ad hoc en `Sidebar.tsx` y no tiene un acceso a Auditoría para el administrador.
  - Prueba: `test/front/navigation.test.tsx` (1/4).
- **Problema:** En el frontend actual, el control de roles y permisos está disperso y acoplado únicamente a los loaders de rutas (`loadAdminSession`). No existe un mecanismo reactivo a nivel de componentes (`usePermissions` / `PermissionGate`) para ocultar o deshabilitar condicionalmente acciones según el rol (`ADMIN` vs `OPERATOR`) o según las instancias asignadas al operador. Esto genera que en vistas como `UserDetail.tsx` o `Header.tsx` se muestren controles estáticos o se dependa de que el backend rechace con 403, en lugar de ofrecer una experiencia fluida y consistente.
- **Entregable:**
  1. Crear un hook `usePermissions()` / `useAuthUser()` (en `src/context/` o `src/hooks/`; un contexto con Provider es opcional) que exponga helpers como `isAdmin`, `isOperator`, `canAccessInstance(vmid: number)` y `hasRole(role: string)`.
  2. Crear un componente wrapper `<PermissionGate requiredRole="ADMIN" fallback={...}>` para condicionar la renderización de botones y secciones administrativas (ej: botón "Eliminar usuario", accesos a auditoría, etc.).
  3. Integrar la reactividad en el menú de navegación y en el header para que las opciones no permitidas a un operador no aparezcan en la interfaz.
- **Criterio de éxito:** Si un operador inicia sesión, la interfaz no muestra accesos directos ni botones exclusivos de administrador; si intenta interactuar con un componente restringido, el helper `canAccessInstance` evalúa en memoria las instancias permitidas cargadas en la sesión; las pruebas unitarias de componentes verifican el render condicional.


### `BAC-18B` - Índice parcial y purga en `sesiones_activas`, y particionamiento trimestral en `auditoria` (PostgreSQL)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `BAC-17B`, `FIX-23`.
- **Estado verificado (01/10/2026, backend `44a2339`):** no implementada. Las pruebas verifican el criterio y no los nombres sugeridos:
  - un índice parcial cualquiera que cubra `jti_access` con `WHERE activa…`;
  - `auditoria` particionada por `RANGE (fecha_hora)`, con tramos trimestrales para el trimestre actual y el siguiente, y una partición `DEFAULT`;
  - un índice compuesto que empiece por `fecha_hora` e incluya `accion` y `resultado`, y otro `(usuario_id, fecha_hora…)`;
  - una purga horaria de sesiones inactivas o vencidas.

  Ya existe el índice `(usuario_id, fecha_hora)` (`idx_auditoria_usuario_fecha`, del modelo GORM); falta todo lo demás.
- **Problema y contexto:**
  1. Las filas históricas o con `activa = false` en `sesiones_activas` penalizan las lecturas en PostgreSQL si no existe un índice parcial ni una purga de registros vencidos.
  2. La tabla `auditoria` es *append-only* (`FIX-23`, no admite `DELETE` bajo ningún concepto) y en la Etapa 1 registrará múltiples eventos por cada operación de Proxmox VE (`PENDING` y `SUCCESS`/`FAILED`). Sin particionamiento por fechas e índices compuestos, las consultas de `GET /api/admin/audit` y la exportación CSV se degradarán progresivamente.
- **Entregable:**
  1. Crear en PostgreSQL (`init.sql` / migración) el índice parcial para `sesiones_activas`:
     ```sql
     CREATE INDEX IF NOT EXISTS idx_sesiones_activas_vigentes
       ON sesiones_activas (jti_token, usuario_id)
       WHERE activa = true;
     ```
  2. Agregar una rutina de limpieza en el backend (ticker cada 1 hora) que ejecute `DELETE FROM sesiones_activas WHERE activa = false OR fecha_expiracion < NOW();`.
  3. Configurar en PostgreSQL el **particionamiento declarativo trimestral por rango de fechas** sobre la tabla `auditoria` (`PARTITION BY RANGE (fecha_hora)`), creando las particiones trimestrales (`auditoria_2026_q3`, `auditoria_2026_q4`, `auditoria_2027_q1` y `auditoria_default`) junto con los índices compuestos `(fecha_hora DESC, accion, resultado)` y `(usuario_id, fecha_hora DESC)`.
- **Criterio de éxito:** Las sesiones muertas se purgan automáticamente de PostgreSQL; la tabla `auditoria` opera sobre particiones trimestrales manteniendo tiempos de consulta constantes (< 20 ms) e inmutabilidad append-only.

### `FIX-29` - Validar la mayúscula en `validatePasswordComplexity` (`FIX-20`) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** `FIX-20`.
- **Problema y evidencia:** en `frontend/centinela/src/components/features/auth/utils/validateAuthenticationFields.ts`, `validatePasswordComplexity` tiene tres condiciones. La segunda prueba `/[0-9]/` pero agrega el mensaje *"Debe contener al menos una letra mayúscula."*, y no hay ninguna condición que pruebe mayúsculas. El dígito se valida una sola vez y la mayúscula nunca, así que el cliente acepta claves como `nueva1234!` o `sinmayus1!`, que el backend rechaza con `400 PASSWORD_CHANGE_FAILED` (`crypto.ValidarComplejidadContrasena`). Además, el mensaje de largo dice *"Debe tener 8 y 12 caracteres."*.
  - Pruebas que fallan:
    - `test/front/password-change.test.tsx`: *"FIX-29 no llama a la API si la contraseña nueva no cumple la complejidad del backend (sin mayúscula)"*.
    - `test/front/recover-password.test.tsx`: *"FIX-29 el paso 3 no envía una contraseña que no cumple la complejidad del backend"*.
- **Entregable:**
  1. Separar las condiciones: `/[A-Z]/` con el mensaje de mayúscula y `/[0-9]/` con el mensaje *"Debe contener al menos un número."*.
  2. Corregir el mensaje de largo: *"Debe tener entre 8 y 12 caracteres."*.
- **Criterio de éxito:** los casos `FIX-29…` pasan, y siguen en verde los de dígito, símbolo y largo, y los casos de `FRN-12` y `FIX-21`.


### `FRN-17C` (`BRG-02-FRN`) - Cliente de eventos con solicitud previa de ticket efímero y reconexión segura

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1,5 h
- **Ventana propuesta:** Junto a `FRN-17A`.
- **Depende de:** `BAC-21C` (`BRG-02-BAC`).
- **Problema y contexto:** El hook `useEvents` del frontend no puede pasar el JWT por header en `EventSource`/WebSocket ni exponer el access token largo en la URL. Debe solicitar primero el ticket efímero al backend.
- **Entregable:**
  1. En `useEvents` (`FRN-17A`), antes de abrir la conexión hacia `/api/events`, invocar `POST /api/events/ticket` con el interceptor autenticado (`Bearer`) y conectar a `/api/events?ticket=<uuid>`.
  2. Ante una desconexión de red, solicitar un nuevo ticket efímero aplicando retroceso exponencial; si `/api/events/ticket` responde `401`, disparar el cierre de sesión local y redirigir a `/login`.
- **Criterio de éxito:** El frontend se conecta a `/api/events` usando tickets de un solo uso sin exponer el JWT en la URL y se reconecta pidiendo un ticket fresco.


### `FIX-36` - Regresión: restaurar el filtro "Acción" en Auditoría (`FRN-14` / RF-08) (Frontend)

- **Área:** Frontend
- **Asignado:** Cristian (autor de `56b88f5`)
- **Estimación:** 0,5 h
- **Depende de:** `FRN-14` y `BAC-18`.
- **Problema y evidencia:** el commit `56b88f5` (*"eliminar columnas que no serán utilizadas"*) quitó de `frontend/centinela/src/pages/Auditoria.tsx` dos cosas: las columnas *Nodo* e *IP origen*, y también el estado `accion`, el campo `aria-label="Acción"` y el `params.set('accion', …)` de la consulta. Las columnas eran opcionales; el filtro no:
  - RF-08 exige *"filtros por fecha, usuario, tipo de acción y resultado"* (`requerimientos.md:193`).
  - El backend sigue aceptando `accion` en `GET /api/admin/audit` (BAC-18 pasa).
  - Prueba que falla: `test/front/audit.test.tsx`, caso *"actualiza los query parameters de filtro al cambiar los selectores reactivos"* (`Unable to find a label with the text of: /acción/i`). Antes pasaba 7/7; ahora 6/7.
- **Entregable:**
  1. Restaurar en `Auditoria.tsx` el filtro "Acción" (`aria-label="Acción"`), con el estado `accion`, su `updateFilter`, la dependencia del `useEffect` y `params.set('accion', accion)`. La exportación CSV tiene que usar los mismos filtros.
  2. Si el equipo decide que el filtro no va, primero hay que modificar RF-08 y avisar para ajustar la prueba. No se quita en silencio.
- **Criterio de éxito:** `audit.test.tsx` vuelve a pasar 7/7 y al elegir una acción la consulta incluye `accion=<valor>`.

### `FIX-37` - Informar el nivel de acceso por instancia en `GET /account/profile` (`FIX-30` / `FRN-18` / `SEC-04`) (Backend)

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 1 h
- **Depende de:** `SEC-04` y `FIX-30`.
- **Problema y evidencia:**
  - `FIX-30` (frontend) arma `canOperateInstance(vmid)` con `perfil.permisos: [{ vmid, nivelAcceso }]` de `GET /account/profile`, como proponía su entregable 2. El backend no envía ese campo: `ports.UsuarioDetalleDTO` (`internal/core/ports/user_port.go:75`) solo tiene `instanciasPermitidas: []int`.
  - Con el backend real, `mapPerfilToUserSession` guarda `permisos: []`, y **`canOperateInstance` da `false` para todo OPERATOR, aunque tenga `FULL_ACCESS`**. Cuando la UI de la Etapa 1 use el helper para habilitar encender, apagar o reiniciar, ningún operador va a poder operar.
  - Prueba que falla: `test/back/resource_access_acceptance_test.go`, caso *"FIX-37 FRN-18 GET /account/profile expone el nivel de acceso por instancia…"*. Las claves que se reciben no incluyen `permisos`.
- **Entregable:**
  1. Agregar a la respuesta de `GET /account/profile` el campo `permisos: [{ vmid, nivelAcceso: "FULL_ACCESS" | "READ_ONLY" }]`, con el mismo formato de `GET /api/admin/users/:id/permissions`, leído de `permisos_instancia`.
  2. Mantener `instanciasPermitidas` para no romper a quienes ya lo usan.
- **Criterio de éxito:**
  - El caso `FIX-37…` pasa: con la 101 en `FULL_ACCESS` y la 102 en `READ_ONLY`, el perfil del OPERATOR informa esos niveles.
  - En el frontend, `canOperateInstance(101)` da `true` y `canOperateInstance(102)` da `false` con datos reales.
  - BAC-14, SEC-04 y LOGIN-04 siguen en verde.

### `FIX-38` - "Solo lectura" (`READ_ONLY`) usado como rol de usuario (`SEC-03` / `FRN-18` / `BAC-09`) (Frontend)

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 0,5 h
- **Depende de:** `BAC-09`, `SEC-04` y `FRN-18`.
- **Revisión (01/10/2026):** la versión anterior de este FIX también pedía mover `usePermissions()` a `src/context/`. **Ese punto se descartó: era un problema de la prueba.** El criterio de éxito de `SEC-03` es de comportamiento y no fija carpeta, y un hook que lee la sesión sin Provider (`useSyncExternalStore`) es válido. `navigation.test.tsx` ahora acepta el hook en `src/context/` o en `src/hooks/`, y `PermissionGate` también en `src/components/`.
- **Problema y evidencia:** los roles de usuario son solo `ADMIN` y `OPERATOR` (`RF-01`, `BAC-09`). `READ_ONLY` es un **nivel de acceso por instancia** (`SEC-04`); `Users.tsx` incluso lo aclara en un comentario. En tres lugares del frontend se lo trata como rol:
  1. `components/features/users/components/informationOfUser.tsx:56`: el selector **"Rol"** de la ficha del usuario ofrece `<option value="READ_ONLY">Solo lectura</option>`. Si el administrador lo elige y guarda, `PUT /api/admin/users/:id` envía `rol: "READ_ONLY"`, y el backend lo rechaza con `400`, porque valida `binding:"omitempty,oneof=ADMIN OPERATOR"` (`backend/internal/core/ports/user_port.go:119`). La UI ofrece una opción que nunca puede guardarse.
  2. `pages/detailsUserPage.tsx:348`, en `formatRole`, traduce el rol `READ_ONLY` a "Solo lectura".
  3. `hooks/usePermissions.ts:26`, en `canAccessInstance`, acepta `user.rol` en `['OPERATOR', 'READ_ONLY']`. Hoy es código muerto, porque el backend nunca emite ese rol, pero mantiene la confusión entre rol y nivel.
- **Entregable:**
  1. Quitar la opción `READ_ONLY` del selector de rol de `informationOfUser.tsx`. "Solo lectura" se asigna **por instancia**, en `rolesAndPermissions.tsx`, que ya lo hace bien.
  2. Quitar el caso `READ_ONLY` de `formatRole` y la condición de rol `READ_ONLY` de `canAccessInstance`. Solo `ADMIN` y `OPERATOR` son roles; el nivel se evalúa en `permisos`.
- **Criterio de éxito:** el selector de rol muestra solo "Operador" y "Administrador"; "Solo lectura" sigue disponible por instancia; un usuario con un rol desconocido no obtiene acceso a ninguna instancia. Los casos de `FRN-05`, `FRN-06B`, `FRN-18` y `FIX-30` siguen en verde.


---
# ETAPA 1
> Las tareas de la Etapa 1 ya verificadas están en [`terminado-1.md`](terminado-1.md), y sus correcciones en [`futuro-1.md`](futuro-1.md).
---

> `FIX-35` se completó el 01/10/2026 y pasó a [`terminado-1.md`](terminado-1.md). No quedan tareas de la Etapa 1 en curso en este archivo.

