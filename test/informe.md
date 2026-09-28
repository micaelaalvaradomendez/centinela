# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 25/09/2026
**Alcance:** tareas en desarrollo de [documentacion/actual.md](../documentacion/actual.md) (SEC-01, SEC-04, FRN-11, FRN-12, SEC-03 y FIX-17 a FIX-25) y regresión de [documentacion/terminado.md](../documentacion/terminado.md).

| Componente | Revisión probada | Cambios desde el informe anterior (23/09) |
|---|---|---|
| Backend | `eb0c9af` (último commit de `main`) | 5 commits: SEC-01 cookie HttpOnly (`c719c9b`, `47bb8f2`), SEC-04 niveles de acceso (`eb0c9af`), FIX-17 (`5407599`, solo documentación) y fecha de último acceso (`fec42b7`) |
| Frontend | `7fbf969` (último commit de `main`) | 28 commits: alta y baja de usuarios con modal de confirmación, guardado del perfil (`PUT /admin/users/:id`), permisos con `nivelAcceso` y revert del rediseño del sidebar (`131ee7d`) |

**Cómo se ejecutó:** cada submódulo se actualizó a su último commit de `origin/main` (`git fetch` + `checkout -B main origin/main --force` + `reset --hard origin/main`). Ante un conflicto prevalece siempre el último commit remoto. Las suites se corrieron sobre esas carpetas, **dos veces** cada una, con el mismo resultado.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan |
|---|---:|---:|---:|
| Backend (`test/back`) | 30 | 28 | **2** |
| Frontend (`test/front`) | 80 | 62 | **18** |
| **Total** | **110** | **90** | **20** |

**Avance desde el 23/09:**
- Pasan a estar cumplidas: **SEC-01**, **FIX-17** y **FIX-25**.
- Quedan casi cumplidas: **SEC-04** y **FIX-24**.
- El hallazgo crítico de integración front ↔ back quedó **resuelto**: el logout y el refresh que envía el frontend ahora funcionan contra el backend (ver [§4](#4-integración-front--back)).

**Todos los fallos son del producto.** En esta revisión también se corrigieron problemas del test (ver [§5](#5-cambios-en-las-pruebas-en-esta-revisión)), cada uno en su lugar.

---

## 2. Cobertura: cada tarea de `actual.md` tiene pruebas

| Tarea | Área | Pruebas que la verifican |
|---|---|---|
| `SEC-01` | Back | `test/back/session_security_acceptance_test.go`: `SEC-01 refresh token solo en cookie HttpOnly` y `SEC-01 SEC-02 integracion…` |
| `SEC-04` → `FIX-26` | Back | `test/back/resource_access_acceptance_test.go`: `SEC-04 niveles de acceso…` y `FIX-26 SEC-04 retrocompatibilidad…` |
| `FRN-11` | Front | `test/front/admin-recovery.test.tsx` (4) |
| `FRN-12` | Front | `test/front/recover-password.test.tsx` (7) |
| `SEC-03` | Front | `test/front/navigation.test.tsx`, bloque SEC-03 (4) |
| `FIX-17` | Back | `test/back/password_recovery_acceptance_test.go`: `FIX-17 contratos canonicos…` |
| `FIX-18` | Front | `recover-password.test.tsx` (7). Incluye la redirección a `/login` y la validación del código en el paso 2 |
| `FIX-19` | Front | `admin-recovery.test.tsx` (4), las mismas pruebas que FRN-11 |
| `FIX-20` | Front | `password-change.test.tsx` (3 casos de complejidad) y `recover-password.test.tsx` (complejidad en el paso 3) |
| `FIX-21` | Front | `password-change.test.tsx`: `la pantalla de cambio ofrece cerrar sesión…` |
| `FIX-22` | Front | `audit.test.tsx`: `si un operador intenta entrar a /auditoria…` |
| `FIX-23` | Back | `password_recovery_acceptance_test.go`: `BAC-18 … append-only`, que verifica UPDATE, DELETE y TRUNCATE |
| `FIX-24` → `FIX-27` | Front | `admin-users.test.tsx`, bloque FRN-06: alta, confirmación, formulario sin contraseña, baja con `DELETE` y los dos casos `FIX-27…` (mensaje del 502 y conservación del formulario) |
| `FIX-25` | Front | `admin-users.test.tsx`: bloque FRN-06B (PUT del perfil) y bloque FIX-25 (409 `USER_CONFLICT`) |

---

## 3. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumplida · 🟡 parcial · ❌ no implementada.

| Tarea | Estado | Resultado | Qué falta / por qué falla | Cómo proceder |
|---|---|---|---|---|
| `SEC-01` Refresh token en cookie HttpOnly | ✅ | 2/2 | — | Se verifica lo siguiente (`auth_handler.go`, `setRefreshCookie`):<br>• `centinela_refresh` con `HttpOnly`, `Secure` por defecto, `SameSite=Strict` y `Path=/api/auth`.<br>• El token no sale en el JSON (`TokenResult.RefreshToken` con `json:"-"`).<br>• Refresh con rotación y 401 ante cookie inválida.<br>• Logout que borra la cookie y el token no aparece en logs.<br>**Se puede pasar a `terminado.md`.** |
| `SEC-04` Niveles de acceso | 🟡 | 1/2 | Funcionan la columna `nivel_acceso` con default y CHECK, el contrato `{ permisos: [{ vmid, nivelAcceso }] }` en PUT y GET, y el guard: READ_ONLY puede hacer GET pero recibe 403 en start/stop, y FULL_ACCESS puede hacer start.<br>**Falla el entregable 4 (retrocompatibilidad):** `PUT /permissions` con el payload anterior `{ "vmids": [103] }` responde `400 INVALID_REQUEST`. Motivo: `asignarPermisosRequest` (`user_handler.go`) solo acepta `permisos`. | **Backend:** aceptar también `vmids` y convertirlo a `permisos` con `FULL_ACCESS`, **o** redefinir el entregable 4 si el equipo decide romper el contrato. El frontend ya usa el contrato nuevo. |
| `FRN-11` Acciones de recuperación | ❌ | 1/4 | "Restablecer contraseña" sigue sin `onClick` (`Users.tsx:347`). No existe la acción "Restablecer 2FA". | Implementar según `FIX-19`. Ya existe un modal reutilizable, `ConfirmUserAction`. |
| `FIX-19` | ❌ | 1/4 | Ídem FRN-11: son las mismas pruebas. | Ídem. Conviene cerrar FRN-11 y FIX-19 juntos. |
| `FRN-12` Recuperación de contraseña | ❌ | 1/7 | `RecoverPassword.tsx` no cambió: solo hace `setCurrentStep`, sin llamadas HTTP, sin validar el código y sin redirigir. | Implementar según `FIX-18`. |
| `FIX-18` | ❌ | 1/7 | Ídem FRN-12. Además falla la redirección a `/login` (la vista se queda en `/recover-password`) y el paso 2 avanza sin código. | Ídem. |
| `SEC-03` Contexto de permisos | ❌ | 1/4 | `src/context/AuthContext.js` sigue vacío. No existen `usePermissions` ni `PermissionGate`, y el menú no tiene un enlace a Auditoría para el ADMIN. | Implementar los 3 entregables. |
| `FIX-17` Contratos de contraseñas | ✅ | 1/1 | — | El commit `5407599` solo actualizó la documentación, pero el comportamiento ya cumple:<br>• Swagger documenta las 4 rutas canónicas y ninguna histórica.<br>• `/auth/change-password` y `/admin/users/:id/reset-password` responden 404.<br>• Códigos de error `PASSWORD_CHANGE_FAILED`, `INVALID_REQUEST` y `RESET_FAILED`.<br>• La recuperación pública deja `cambio_contrasena=false` y revoca las sesiones.<br>**Se puede pasar a `terminado.md`.** |
| `FIX-20` Complejidad de contraseñas | ❌ | 0/4 | `ChangePassword.tsx` sigue validando solo el largo y envía `nueva1234!`, `NuevaClave!` y `Nueva12345`. En `RecoverPassword.tsx` no hay validación porque no hay formulario conectado. | Implementar `validatePasswordComplexity` según el FIX. |
| `FIX-21` Cerrar sesión en cambio de clave | ❌ | 0/1 | `ChangePassword.tsx` sigue sin botón "Cerrar sesión". | Según el FIX. |
| `FIX-22` Guard de `/auditoria` | ❌ | 0/1 | `/auditoria` sigue bajo `loadProtectedSession` (`applicationRoutes.tsx:50`): **un OPERATOR puede entrar a la vista**. | Moverla al grupo `loadAdminSession`. Es un cambio de 1 línea con prioridad alta, porque es un control de acceso. |
| `FIX-23` Auditoría append-only | ❌ | 0/1 | La base sigue permitiendo `UPDATE`, `DELETE` y `TRUNCATE` sobre `auditoria`: no hay ningún trigger ni `REVOKE` en el backend. | Aplicar el trigger del FIX. |
| `FIX-24` Alta sin contraseña y baja | 🟡 | 5/6 | **Implementado:**<br>• Alta con confirmación de correo y modal (`ConfirmUserAction`).<br>• Toast "Usuario creado exitosamente… enviada por correo".<br>• Formulario sin campos de contraseña.<br>• "Eliminar usuario" con modal que envía `DELETE /admin/users/:id`.<br>**Falla el entregable 2:** ante `502 EMAIL_DELIVERY_FAILED`, `createUserService.ts` arma el mensaje correcto, pero `useCreateUser.ts` lo descarta y muestra el genérico *"No se pudo completar la creación del usuario"*. | En el `catch` de `useCreateUser`, usar `error.message` en la descripción del toast cuando `errorCode === 'EMAIL_DELIVERY_FAILED'`. |
| `FIX-25` Persistir edición del perfil | ✅ | 3/3 | — | Se verifican el `PUT /admin/users/:id` con los campos modificados (`userDetailsService.updateUserDetails`) y el 409 mostrado junto al campo de correo. Los 4 casos de FIX-14 siguen en verde, incluido "Solo lectura", que envía `READ_ONLY`. **Se puede pasar a `terminado.md`.** |

**Hito "Gestión Administrativa de Usuarios":** el recorrido completo ya funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al operador. Solo queda el mensaje del 502 (`FIX-27`).

> [!NOTE]
> **Reclasificación (25/09/2026):**
> - `SEC-01`, `FIX-17` y `FIX-25` pasaron todas sus pruebas y se movieron a `terminado.md`.
> - `SEC-04` y `FIX-24` están implementadas con un entregable incumplido cada una. Pasaron a `terminado.md` con aviso, y sus correcciones quedaron en `futuro.md`:
>   - **`FIX-26`**: retrocompatibilidad `{ vmids }` en backend.
>   - **`FIX-27`**: mensaje del 502 en el alta.
>
>   Sus pruebas llevan el ID del FIX. `FIX-27` tiene además un caso que verifica que el formulario conserva los datos tras el 502; hoy pasa y queda como control de regresión.

---

## 4. Integración front ↔ back

| Flujo | 23/09 | 25/09 |
|---|---|---|
| Logout desde la interfaz (`POST /auth/logout`, Bearer + body `{}`) | ❌ 400 `INVALID_REQUEST`, la sesión seguía activa | ✅ 204, revoca access y refresh (cookie `centinela_refresh`) |
| Renovación silenciosa (`POST /auth/refresh` con body `{}`) | ❌ 400 | ✅ 200, con rotación de la cookie |
| Permisos por instancia (`{ permisos: [{ vmid, nivelAcceso }] }`) | — | ✅ El frontend (`userInstanceService.ts`) y el backend usan el mismo contrato; "Solo lectura" llega como `READ_ONLY` y el guard lo aplica |
| Complejidad de contraseña | ❌ | ❌ El frontend sigue enviando claves que el backend rechaza con 400 (FIX-20) |

---

## 5. Cambios en las pruebas en esta revisión

| Cambio | Motivo |
|---|---|
| Permisos: `{ vmids }` → `{ permisos: [{ vmid, nivelAcceso }] }` en BAC-05/06, BAC-07, LOGIN-04, SEC-04 y FIX-14 | El contrato canónico cambió con SEC-04. La retrocompatibilidad con `{ vmids }` se verifica aparte, como entregable 4 de SEC-04 |
| SEC-04 asigna los niveles por API en lugar de por SQL, y verifica el `nivelAcceso` en el GET | Ahora existe el contrato |
| SEC-01: se agregó la verificación de `Secure` por defecto | Entregable 4 |
| Tests nuevos: FIX-17 y el `TRUNCATE` de FIX-23 | Tareas nuevas en `actual.md` sin cobertura |
| Tests nuevos: FIX-18 (redirección y código vacío), FIX-20 (paso 3), FIX-24 (502 y formulario sin contraseña), FIX-25 (409) y FIX-14 (Solo lectura) | Entregables sin cobertura |
| El alta de FRN-06 completa "Confirmar correo" y confirma el modal | Flujo nuevo del frontend |
| El menú del ADMIN ya no exige el enlace "Crear usuario" | Se quitó a propósito (`a387e10`); el alta se inicia desde `/users` |
| El test de FIX-20 en el paso 3 incluye una contraprueba (una clave válida sí se envía) | Sin ella daba verde con una vista que nunca llama a la API |
| Timeouts de 15 s en los flujos de varios pasos (alta y recuperación) | Tipear 5 campos superaba los 5 s por defecto: era un fallo del test |
| `setup.ts`: polyfill de `HTMLDialogElement.showModal/close` | jsdom no lo implementa y los modales nuevos lo usan |
| `compose.yaml` y `sourcePath()` aceptan `CENTINELA_ROOT` | Para construir el backend de `origin/main` sin tocar el submódulo |

---

## 6. Recomendaciones para la documentación

- ~~Pasar a `terminado.md`: `SEC-01`, `FIX-17` y `FIX-25`.~~ **Hecho el 25/09/2026.** También `SEC-04` y `FIX-24`, con sus correcciones `FIX-26` y `FIX-27`.
- **`FIX-26`:** decidir si se mantiene la retrocompatibilidad con `{ vmids }`. Si no se mantiene, corregir el entregable 4 de `SEC-04` y descartar el FIX.
- **`FIX-21` y `FIX-22`:** sus criterios citan "7/7", pero hoy esos archivos tienen 10 y 7 casos. Conviene decir "los casos de FIX-21/FIX-22 pasan" en lugar de un número fijo.
- **Duplicados:** `FRN-11`/`FIX-19` y `FRN-12`/`FIX-18` describen el mismo trabajo con las mismas pruebas. Conviene dejar una sola tarea de cada par.

---

## 7. Cómo reproducir

```bash
# 1. Traer el último commit de cada submódulo (prevalece el remoto ante conflictos)
for s in backend frontend; do
  git -C $s fetch origin --prune
  git -C $s checkout -B main origin/main --force
  git -C $s reset --hard origin/main
done

# 2. Backend (Docker: PostgreSQL + backend + stub de Proxmox)
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)

# 3. Frontend
pnpm --dir test/front test
```

Usar `| tee` y no `> archivo` para capturar la salida del backend: en este entorno, redirigir con `>` cortó la corrida y dejó contenedores huérfanos. El comando de limpieza está en `test/back/README.md`.

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
