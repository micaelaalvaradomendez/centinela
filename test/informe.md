# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 23/09/2026
**Fuentes:** [documentacion/actual.md](../documentacion/actual.md) (tareas en curso), [documentacion/terminado.md](../documentacion/terminado.md) (regresión) y las suites de [test/back](back) y [test/front](front).

| Componente | Revisión probada | Cómo |
|---|---|---|
| Backend | `15032da` = `origin/main` (*fix: asegurar logout y revocación atómica de sesiones*) | `go test -v -count=1 ./...` en `test/back`, con PostgreSQL 16 y un stub de Proxmox VE en Docker Compose |
| Frontend | `3192cf4` = `origin/main` (*Merge PR #56 rama-Beli*) | `CENTINELA_FRONTEND_DIR=<export de origin/main> pnpm test` en `test/front`, sin modificar el submódulo |

> [!IMPORTANT]
> **El submódulo `frontend` de este repo no refleja el trabajo del equipo.** Está en `5b91e80`, 13 commits atrás de `origin/main`, y tiene **3 archivos modificados sin commitear que no son del equipo**: `applicationRoutes.tsx`, `apiClient.ts` y `tokenStorage.ts`. Esos parches hacían pasar pruebas que el código real no pasa. Por ejemplo, movían `/auditoria` bajo `loadAdminSession` y volvían a guardar el `refreshToken` en contra de SEC-02. Por eso el frontend se probó contra una exportación limpia de `origin/main` (`git archive`). Se recomienda descartar esos parches y actualizar el puntero del submódulo (ver [§7](#7-parches-locales-en-el-submódulo-frontend)).

> [!NOTE]
> **Reclasificación de la documentación (23/09/2026), según estos resultados:**
>
> | Situación | Tareas | Dónde quedan |
> |---|---|---|
> | Completas | BAC-14, FIX-16, BAC-17, FRN-13, SEC-02, FIX-14 | `terminado.md` |
> | Implementadas con problemas | FRN-10, FRN-14 | `terminado.md`, con **FIX-21** y **FIX-22** en `futuro.md` |
> | Retrocesos de tareas que ya estaban terminadas | BAC-18, FRN-06, FRN-06B | **FIX-23**, **FIX-24** y **FIX-25** en `futuro.md` |
> | No implementadas | SEC-01, SEC-04, FRN-11, FRN-12, SEC-03 | `actual.md` |
>
> Cada fallo de este informe corresponde a una tarea de `actual.md` o a uno de esos FIX.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`, 5 archivos) | 28 | 24 | **4** | 0 |
| Frontend (`test/front`, 12 archivos) | 72 | 55 | **17** | 0 |
| **Total** | **100** | **79** | **21** | **0** |

Las pruebas se ejecutaron **dos veces** con idéntico resultado, así que no hay tests inestables. **Los 21 fallos son brechas del producto, no del test.** Cada fallo se clasifica en [§3](#3-estado-de-las-tareas-de-actualmd) y [§4](#4-regresiones-en-tareas-de-terminadomd).

**El informe anterior (22/09) reportaba 84/84 en verde, pero no era fiable.** Buena parte de ese verde venía de falsos positivos (ver [§2](#2-qué-se-corrigió-en-las-pruebas)) y de los parches locales en el submódulo.

---

## 2. Qué se corrigió en las pruebas

Estos eran los patrones que daban verde sin verificar el criterio de la tarea:

| Patrón encontrado | Dónde | Corrección |
|---|---|---|
| `t.Logf("PENDIENTE…"); return` en lugar de fallar | SEC-01, SEC-04, pasos 6.4/6.5 de LOGIN-04 | Ahora son aserciones (`t.Errorf`/`t.Fatalf`) sobre los entregables de `actual.md` |
| Un `502` contado como "verificado" porque no había Proxmox | BAC-14 | Nuevo **stub de Proxmox VE** (`test/back/proxmox-stub/nginx.conf`) con un inventario fijo. Se verifican los campos normalizados, el filtro por operador y el mapeo de errores |
| `if (llamada) { … } else { expect(boton).toBeVisible() }` | FRN-06, FRN-06B, FRN-07, FRN-11, FRN-12 | Se exige la llamada HTTP real, con su método, ruta y payload |
| URLs inventadas (`/recover/request`, `/recover/confirm`) | FRN-12 | Se usan los endpoints reales `/auth/password/forgot` y `/auth/password/reset` |
| Tokens firmados por el test, y `logout` sin Bearer que aceptaba 200 **o** 401 | BAC-17 | Los tokens salen del circuito real login → 2FA. Se verifica la revocación en BD, `TOKEN_REVOKED`, el refresh posterior y la auditoría `LOGOUT` |
| Clave temporal forzada copiando el hash del admin; el test dependía del orden de ejecución | LOGIN-04 | La clave temporal se lee del **MockEmailService** (logs del contenedor), el mismo canal por el que la recibe el usuario |
| `t.Log("BAC-21 verificado")` sin ninguna aserción | BAC-21 | Compara los campos y enums del struct Go contra la interfaz TypeScript |
| BAC-16, BAC-19, BAC-20, BAC-13, BAC-15 solo verificaban códigos HTTP | regresión | Ahora verifican el efecto: el código anterior se invalida, las sesiones se revocan, la clave vieja deja de funcionar y el secreto viejo se rechaza |
| Tests atados al contrato anterior (`getRefreshToken`, `refreshToken` en el body) | SEC-02, FRN-13 | Se actualizaron al contrato con cookie HttpOnly |
| Falta de `window.matchMedia` en jsdom | todos los que montan el sidebar | Se agregó un polyfill en `test/front/setup.ts` |

También se agregó una **prueba de integración** que reproduce, contra el backend, las peticiones exactas que envía el frontend de `origin/main` (ver [§5](#5-hallazgo-de-integración-front--back-crítico)).

---

## 3. Estado de las tareas de `actual.md`

**Leyenda:**
- ✅ Cumplida (todas sus pruebas pasan).
- 🟡 Parcial.
- ❌ No implementada.
- 🔗 El lado propio está cumplido, pero falla la integración.

| Tarea | Área | Estado | Pruebas | Por qué falla / qué falta | Cómo proceder |
|---|---|---|---|---|---|
| `BAC-14` Inventario Proxmox | Back | ✅ | `resource_access…` BAC-14 | — | Nada. El caso `502 PROXMOX_UNAVAILABLE` no se prueba, porque el stub siempre responde. |
| `FIX-16` Guard en `/instances` | Back | ✅ | `resource_access…` FIX-16 | — | Nada. Se verifica el 403 `INSTANCE_ACCESS_DENIED` en GET/start/stop y que Proxmox no recibe la orden. |
| `BAC-17` Logout atómico | Back | ✅ 🔗 | `session_security…` BAC-17 | Con el contrato actual (refresh token en el body) revoca ambas sesiones y audita `LOGOUT`. **Pero el frontend ya no manda ese body** ([§5](#5-hallazgo-de-integración-front--back-crítico)). | Se resuelve con SEC-01. |
| `SEC-01` Refresh en cookie HttpOnly | Back | ❌ | `session_security…` SEC-01 e integración | `POST /auth/2fa/verify` no emite `Set-Cookie` y sigue devolviendo `refreshToken` en el JSON. `/auth/refresh` y `/auth/logout` exigen `refreshToken` en el body (`auth_handler.go:69-71`, `:190-192`). | Implementar los 5 entregables. **Es la tarea que desbloquea el logout y la renovación de sesión reales.** |
| `SEC-04` Niveles de acceso | Back | ❌ | `resource_access…` SEC-04 | `permisos_instancia` no tiene `nivel_acceso`; `VerificarAcceso` y `RequireInstanceAccess` no reciben el nivel. | Implementar la migración, el puerto y el guard. **El frontend ya ofrece "Solo lectura"** en el selector (`rolesAndPermissions.tsx`), pero el backend lo ignora. |
| `FRN-10` Cambio obligatorio de contraseña | Front | 🟡 | `password-change` (7/10) | Falta la opción **"Cerrar sesión"** en `/change-password` (FIX-21). Además, la validación del cliente no aplica las reglas de complejidad del backend (mayúscula, dígito y carácter especial) y envía claves que el backend rechaza con 400 (FIX-20). | FIX-21 y FIX-20. |
| `FRN-11` Acciones de recuperación | Front | ❌ | `admin-recovery` (1/4) | "Restablecer contraseña" no tiene `onClick` (`Users.tsx:354`). No existe una acción de reset de 2FA ni un diálogo de confirmación. | Implementar las dos acciones: confirmación, `POST …/password/reset` y `POST …/2fa/reset`, toast y refresco de la tabla. |
| `FRN-12` Recuperación de contraseña | Front | ❌ | `recover-password` (1/4) | `RecoverPassword.tsx` solo cambia de paso; no hace ninguna llamada HTTP. | Conectar el paso 1 a `POST /auth/password/forgot` y el paso 3 a `POST /auth/password/reset` con `{ email, codigo, nuevaContrasena }`, y mostrar `RESET_FAILED`. |
| `FRN-13` Logout en cliente | Front | ✅ 🔗 | `session-security` FRN-13 (6/6) | El lado cliente cumple: envía Bearer, limpia todas las claves aun ante error, usa `replace` y tiene listener de 401. El backend responde 400 a su logout ([§5](#5-hallazgo-de-integración-front--back-crítico)). | Se resuelve con SEC-01. |
| `SEC-02` Cliente sin refresh token | Front | ✅ 🔗 | `session-security` SEC-02 (4/4) | Cumple: no guarda ni envía el refresh token y renueva con `credentials: 'include'`. El backend todavía no lo soporta. | Se resuelve con SEC-01. |
| `FIX-14` / `FRN-07` Selector de instancias | Front | ✅ | `admin-users` FIX-14 (3/3) | Se cumple el criterio de éxito: GET con Bearer, asignadas marcadas y PUT `{ vmids }`. El entregable 3 ("atómica junto a la edición del perfil") falla por FRN-06B. | Ver FRN-06B en [§4](#4-regresiones-en-tareas-de-terminadomd). |
| `SEC-03` Contexto de permisos | Front | ❌ | `navigation` SEC-03 (1/4) | `src/context/AuthContext.js` está vacío. No existen `usePermissions` ni `PermissionGate`, y el menú no tiene acceso a Auditoría para el admin. Solo pasa el test de que el operador no ve los accesos, que hoy resuelve un `isAdmin` ad hoc en `Sidebar.tsx`. | Implementar el contexto, el gate y el enlace a Auditoría para admin. |
| `FRN-14` Pruebas de Auditoría | Front | 🟡 | `audit` (6/7) | La suite existe y cubre la carga, los filtros (acción, resultado y fechas), la paginación y el CSV. **Falla el guard**: en el código commiteado `/auditoria` está bajo `loadProtectedSession`, no bajo `loadAdminSession` (`applicationRoutes.tsx`), así que **un OPERATOR puede entrar**. | Mover la ruta dentro del grupo `loadAdminSession`. |

---

## 4. Regresiones en tareas de `terminado.md`

Estas tareas figuran como terminadas, pero las pruebas muestran que su criterio de éxito no se cumple:

| Tarea | Prueba | Problema | Cómo proceder |
|---|---|---|---|
| `BAC-18` Auditoría append-only | `password_recovery…` BAC-18 | El registro, la consulta, el filtro, el CSV y el 403 al operador funcionan. **Pero la base permite `UPDATE` y `DELETE` sobre `auditoria`** con las credenciales de la aplicación: no hay trigger, `REVOKE` ni regla. El criterio exige que fallen "a nivel de base de datos". | Crear un fix de backend: trigger `BEFORE UPDATE OR DELETE` que lance una excepción, o un rol de aplicación sin esos privilegios. |
| `FRN-06` Alta de usuarios | `admin-users` FRN-06 | Hay dos fallas. **(a)** Por BAC-16, el backend ya no devuelve `contrasenaTemp`, pero `CrearUsuarios.tsx` sigue esperándola: tras un alta exitosa la pantalla **no muestra ninguna confirmación**. **(b)** "Eliminar usuario" (`detailsUserPage.tsx:110`) no tiene `onClick` y no envía `DELETE`. | Crear un fix de frontend: toast "Usuario creado, la clave se envió por correo", quitar la caja de contraseña temporal y conectar el `DELETE` con confirmación. |
| `FRN-06B` Edición de usuario | `admin-users` FRN-06B | Editar el nombre o el correo no se guarda: "Guardar cambios" solo se habilita si cambió algún permiso de instancia, y **no existe ningún `PUT /admin/users/:id`** en el frontend. | Crear un fix de frontend: persistir `useEditableUser` con `PUT /admin/users/:id`. |

Las demás tareas de `terminado.md` pasan con pruebas reforzadas: BAC-01 a BAC-13, BAC-15, BAC-16, BAC-19 a BAC-21, LOGIN-01 a LOGIN-04, FRN-01 a FRN-05, FRN-08, FRN-09, FIX-07 y FIX-08.

---

## 5. Hallazgo de integración front ↔ back (crítico)

El frontend de `origin/main` ya implementó **SEC-02** (commits `cb208f5`, `f0f7218`, `390b4b8`): envía `POST /auth/logout` y `POST /auth/refresh` con body `{}`, confiando en la cookie. El backend todavía **no implementó SEC-01** y rechaza esas peticiones.

La prueba `SEC-01 SEC-02 integracion…` (`session_security_acceptance_test.go`) reproduce las peticiones reales del frontend:

```
POST /auth/refresh  {}          -> 400 {"errorCode":"INVALID_REQUEST","message":"Se requiere el campo refreshToken."}
POST /auth/logout   Bearer + {} -> 400 {"errorCode":"INVALID_REQUEST","message":"...Se requiere refreshToken."}
GET  /account/profile (mismo access token, tras el "logout") -> 200
```

**Impacto con ambos `main` desplegados juntos:**
1. **El logout desde la interfaz no revoca nada en el servidor.** El cliente borra su almacenamiento y redirige, pero el access token sigue siendo válido hasta que expira. Esto incumple el criterio de BAC-17 y FRN-13.
2. **La renovación silenciosa de sesión está rota.** Cuando el access token vence, el refresh falla y el usuario es expulsado al login.

**Cómo proceder:** priorizar SEC-01. Otra opción, mientras tanto, es que el backend acepte el refresh token por cookie **o** por body.

---

## 6. Documentación desactualizada detectada

- `actual.md` / FRN-10 cita `POST /api/auth/change-password`. El endpoint real (backend y frontend) es **`PUT /api/account/password`**.
- `actual.md` / BAC-17 y FIX-16: el bloque "Estado de implementación actual" describe problemas que ya se corrigieron (logout sin middleware, rutas `/instances` comentadas).
- `actual.md` / FIX-14: cita `UserDetail.tsx` con una lista estática; el componente ahora es `detailsUserPage.tsx` con `rolesAndPermissions.tsx`, conectado a la API.
- `terminado.md` / BAC-08 y FRN-07: siguen con el aviso "Implementado con fallo en pruebas", pero ya pasan (FIX-16 y FIX-14).

---

## 7. Parches locales en el submódulo frontend

`git -C frontend status` muestra estas modificaciones sin commitear, que **no son del equipo** y se hicieron durante una evaluación anterior:

| Archivo | Cambio | Efecto |
|---|---|---|
| `centinela/src/routes/applicationRoutes.tsx` | Mueve `/auditoria` bajo `loadAdminSession` | Ocultaba el fallo de FRN-14 |
| `centinela/src/services/apiClient.ts` | Acepta respuestas de refresh sin `refreshToken` y conserva el anterior | Contradice SEC-02 y entra en conflicto con `origin/main` |
| `centinela/src/storage/tokenStorage.ts` | Borra `centinela_user` de sessionStorage y localStorage | Ocultaba un test mal sembrado de FRN-13 |

Además, el puntero de los submódulos en el repo padre está sin commitear (`backend` 7cadbce → 15032da y `frontend` d088df2 → 5b91e80). **Recomendación:** que quien administre el repo descarte esos 3 cambios (`git -C frontend checkout -- centinela/src`) y actualice el submódulo a `origin/main`. La suite está preparada para ese estado.

---

## 8. Cómo reproducir

```bash
# Backend (requiere Docker; levanta PostgreSQL + backend + stub de Proxmox y los elimina al terminar)
cd test/back && go test -v -count=1 ./...

# Frontend contra el submódulo tal como está
pnpm --dir test/front test

# Frontend contra origin/main sin tocar el submódulo
mkdir -p /tmp/front-main
git -C frontend fetch && git -C frontend archive origin/main centinela | tar -x -C /tmp/front-main
ln -s "$PWD/frontend/centinela/node_modules" /tmp/front-main/centinela/node_modules
CENTINELA_FRONTEND_DIR=/tmp/front-main/centinela pnpm --dir test/front test
```

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
