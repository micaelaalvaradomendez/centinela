# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 28/09/2026
**Alcance:** tareas en desarrollo de [documentacion/actual.md](../documentacion/actual.md) y regresión de [documentacion/terminado.md](../documentacion/terminado.md).

| Componente | Revisión probada | Cambios desde el informe anterior (25/09) |
|---|---|---|
| Backend | `d36bc50` (último commit de `main`) | CORS con lista blanca (`c624a57`), Redis en `docker-compose.yml` (`2a943cd`), API token de Proxmox (`44feb08`) y email normalizado en la recuperación (`2bdc3b1`) |
| Frontend | `8d7ecab` (último commit de `main`) | Reset de contraseña y 2FA en `Users.tsx` (`fe158ab`), recuperación de contraseña conectada (`fe6d784`), validador de complejidad, "Cerrar sesión" en el cambio de clave (`b87ea1b`), rediseño del logout con `AuthProvider` (`e02b821` a `4062fdc`) y sidebar nuevo (`81c6424`) |

**Cómo se ejecutó:** cada submódulo se actualizó a su último commit de `origin/main` (`fetch` + `checkout -B main origin/main --force` + `reset --hard`). Ante un conflicto prevalece el remoto. Las suites se corrieron sobre esas carpetas **dos veces** cada una, con el mismo resultado.

> [!NOTE]
> **Actualización del 28/09/2026 (segunda revisión, sin commits nuevos en los submódulos, sin cambios en las pruebas; los resultados de la corrida son los mismos):**
> - `FRN-11`, `FIX-19`, `FRN-12`, `FIX-18`, `FIX-21` e `INF-06` pasaron a `terminado.md`.
> - Se revisó una por una cada tarea "no implementada" en busca de una implementación alternativa en `main` o en otra rama. No hay ninguna (ver [§3.1](#31-revisión-de-las-tareas-no-implementadas)).
> - Se confirmó la regresión de `FRN-13` y se registró como **`FIX-28`** en `futuro.md` (ver [§4](#4-regresiones-en-tareas-de-terminadomd)).
> - El merge pendiente del repo principal ya está cerrado.
> - **`FIX-26` descartado:** se adoptó el formato nuevo `{ permisos }` como contrato oficial. SEC-04 queda completa, su entregable 4 se corrigió en `terminado.md` y la prueba ahora verifica que `{ vmids }` se rechace con 400. El backend pasa de 7 a 6 fallos.
> - `FIX-20`, `FRN-18` e `INF-05` están implementadas con un problema cada una: pasaron a `terminado.md` con aviso, y sus correcciones son **`FIX-29`** (mayúscula en el validador), **`FIX-30`** (`canOperateInstance`) y **`FIX-31`** (TLS de Nginx versionado) en `futuro.md`. Sus pruebas llevan el ID del FIX.
>
> **Clasificación resultante:**
>
> | Situación | Tareas | Dónde quedan |
> |---|---|---|
> | Implementadas y completas | FRN-11, FIX-19, FRN-12, FIX-18, FIX-21, INF-06 | `terminado.md` |
> | Implementadas con problema | FIX-20, FRN-18, INF-05 | `terminado.md` con aviso, y `FIX-29`, `FIX-30` y `FIX-31` en `futuro.md` |
> | Regresión de una tarea terminada | FRN-13 | Aviso en `terminado.md` y `FIX-28` en `futuro.md` |
> | Sin implementar | SEC-03, FIX-22, FIX-23, INF-08, BAC-16B, BAC-17B, BAC-18B | `actual.md` |

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 38 | 31 | **6** | 1 |
| Frontend (`test/front`) | 82 | 73 | **9** | 0 |
| **Total** | **120** | **104** | **15** | **1** |

El caso omitido es el TLS de INF-05: la configuración de Nginx del servidor (CT 103) no está versionada en ningún repositorio, así que no se puede verificar de forma automática.

**Avances:**
- Pasan a estar completas: **FRN-11 / FIX-19**, **FRN-12 / FIX-18**, **FIX-21** e **INF-06**.
- El CORS de **INF-05** también funciona.

**Regresión crítica en una tarea terminada (FRN-13):** desde `deb59cb`, el logout del frontend ya no envía el `Authorization: Bearer`. El backend lo rechaza con `401 MISSING_TOKEN` y **la sesión queda activa en el servidor** (ver [§4](#4-regresiones-en-tareas-de-terminadomd)).

**Todos los fallos son del producto.** Los problemas que eran del test se corrigieron (ver [§6](#6-cambios-en-las-pruebas-en-esta-revisión)).

---

## 2. Cobertura: cada tarea de `actual.md` tiene pruebas

| Tarea | Pruebas |
|---|---|
| `FRN-11` / `FIX-19` | `test/front/admin-recovery.test.tsx` (4) |
| `FRN-12` / `FIX-18` | `test/front/recover-password.test.tsx` (7) |
| `SEC-03` | `test/front/navigation.test.tsx`, bloque SEC-03 (4) |
| `FIX-20` | `password-change.test.tsx` (3 casos de complejidad) y `recover-password.test.tsx` (paso 3) |
| `FIX-21` | `password-change.test.tsx`: `la pantalla de cambio ofrece cerrar sesión…` |
| `FIX-22` | `audit.test.tsx`: `si un operador intenta entrar a /auditoria…` |
| `FIX-23` | `test/back/password_recovery_acceptance_test.go`: `BAC-18 … append-only` (UPDATE, DELETE y TRUNCATE) |
| `INF-05` | `test/back/cierre_fase_base_acceptance_test.go`: `INF-05 CORS…` y `INF-05 TLS…` (omitido) — **nuevo** |
| `INF-06` | `cierre_fase_base…`: `INF-06 Redis…`. Levanta el servicio `redis` del `docker-compose.yml` del equipo y lo prueba en vivo — **nuevo** |
| `INF-08` | `cierre_fase_base…`: `INF-08 variables SMTP…` — **nuevo** |
| `BAC-16B` | `cierre_fase_base…`: `BAC-16B …` contra un segundo backend con `EMAIL_PROVIDER=smtp` y el servidor SMTP de prueba **Mailpit** — **nuevo** |
| `BAC-17B` | `cierre_fase_base…`: `BAC-17B …` (login + 2FA + 10 refresh + logout) — **nuevo** |
| `BAC-18B` | `cierre_fase_base…`: `BAC-18B …` (índice parcial, particiones e índices compuestos) — **nuevo** |
| `FRN-18` | `admin-users.test.tsx`: 2 casos `FRN-18…` (leer y enviar `nivelAcceso`) y `navigation.test.tsx`: `FRN-18 canOperateInstance…` — **nuevo** |

**Lo que queda sin verificación automática, por el alcance de las tareas:**
- **INF-05:** el TLS de Nginx en el servidor.
- **INF-06:** la red `vmbr1` del servidor.
- **INF-08:** la conectividad saliente real hacia el proveedor SMTP.
- **BAC-18B:** la purga horaria (requiere esperar el ticker) y el tiempo de consulta menor a 20 ms (requiere un volumen de datos real).

---

## 3. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumplida · 🟡 implementada con problema o parcial · ❌ no implementada.

| Tarea | Estado | Resultado | Qué falta / por qué falla | Cómo proceder |
|---|---|---|---|---|
| `FRN-11` / `FIX-19` Reset de contraseña y 2FA | ✅ | 4/4 | — | `Users.tsx` tiene las dos acciones con `ConfirmUserAction`, llama a `POST …/password/reset` y `…/2fa/reset` con Bearer, muestra el toast y la tabla pasa a "Desactivado". **Movida a `terminado.md`.** |
| `FRN-12` / `FIX-18` Recuperación de contraseña | ✅ | 6/7 | El único fallo es de `FIX-20` (paso 3, ver abajo). Todo lo de FRN-12 y FIX-18 pasa: forgot y reset con los payloads reales, validación del código de 6 dígitos, error `RESET_FAILED` y redirección a `/login`. | **Movida a `terminado.md`.** El fallo restante lo cubre FIX-20. |
| `FIX-21` Cerrar sesión en el cambio de clave | ✅ | 1/1 | — | **Movida a `terminado.md`.** Su revocación en el servidor depende de `FIX-28`. |
| `INF-06` Redis | ✅ | 1/1 | — | `redis:7-alpine` con `requirepass`: responde `PONG` con clave y `NOAUTH` sin ella, con `maxmemory 256mb` y `volatile-lru`. `REDIS_ADDR` y `REDIS_PASSWORD` están en `.env.example`. **Movida a `terminado.md`**; la red `vmbr1` se verifica en el servidor. |
| `INF-05` CORS y TLS | 🟡 | 1/1 + 1 omitido | **CORS cumple:** el origen permitido recibe 204 con `Allow-Origin` exacto y `Allow-Credentials: true`; uno no autorizado recibe 403 sin cabeceras. **TLS no se puede verificar:** la configuración de Nginx del CT 103 no está versionada. | Infraestructura: versionar la configuración de Nginx (por ejemplo en `docker/` o en el repo de frontend) o validarla a mano. |
| `FIX-20` Complejidad de contraseñas | 🟡 | 2/4 | **Bug:** `validatePasswordComplexity` (`components/features/auth/utils/validateAuthenticationFields.ts`) **no verifica mayúsculas**. La segunda condición prueba `/[0-9]/` con el mensaje de mayúscula, así que el dígito se valida dos veces y la mayúscula nunca. Una clave como `nueva1234!` o `sinmayus1!` pasa el cliente y el backend la rechaza con 400. | Reemplazar esa condición por `/[A-Z]/` y agregar un chequeo aparte de dígito. **Propuesta: pasar a `terminado.md` y crear un FIX.** |
| `FRN-18` Niveles de acceso en frontend | 🟡 | 2/3 | **Entregable 1 cumple:** lee `nivelAcceso` (READ_ONLY → "Solo lectura") y lo envía. **Entregable 2 no:** no existe `canOperateInstance`, porque depende de SEC-03. | Implementarlo junto con SEC-03. |
| `SEC-03` Contexto de permisos | ❌ | 1/4 | `context/AuthContext.js` ahora existe, pero solo maneja la expulsión por `TOKEN_REVOKED`: no exporta `usePermissions` ni `PermissionGate`. Además, **el sidebar nuevo (`81c6424`) muestra "Auditoría" a todos los roles**: solo "Usuarios" está bajo `rol === 'ADMIN'`. | Implementar el hook y el gate, y poner el enlace a Auditoría bajo el gate de ADMIN. |
| `FIX-22` Guard de `/auditoria` | ❌ | 0/1 | Sigue fuera de `loadAdminSession` (`applicationRoutes.tsx:50`). Sumado al enlace visible en el sidebar, **un OPERATOR ve el acceso y entra**. | Mover la ruta al grupo `loadAdminSession`. Prioridad alta. |
| `FIX-23` Auditoría append-only | ❌ | 0/1 | Sin trigger ni `REVOKE`: se permiten UPDATE, DELETE y TRUNCATE. | Según el FIX. |
| `INF-08` Variables SMTP | ❌ | 0/1 | `.env.example` y `docker-compose.yml` del backend no documentan `EMAIL_PROVIDER`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` ni `SMTP_FROM`. | Agregarlas. |
| `BAC-16B` Adaptador SMTP | ❌ | 0/1 | Con `EMAIL_PROVIDER=smtp` el backend sigue usando `MockEmailService` (`main.go`: `emailService := email.NewMockEmailService()`) y Mailpit no recibe ningún correo. | Implementar `SmtpEmailService` y elegirlo según `EMAIL_PROVIDER`. **Nota:** la tarea cita `EnviarCredencialesTemporales`, pero el puerto real se llama `EnviarContrasenaTemporal` (`ports/email_port.go`). |
| `BAC-17B` Una sesión, un registro | ❌ | 0/1 | Confirmado el problema que describe la tarea: login + 2FA + 10 refresh crean **13 filas activas** en `sesiones_activas`, y después del logout **quedan 11 activas**. `fecha_ultimo_acceso` sí se actualiza. | Según la tarea. |
| `BAC-18B` Índice parcial y particionado | ❌ | 0/1 | No existen `idx_sesiones_activas_vigentes` ni las particiones (`auditoria` es una tabla común, `relkind = r`), ni los índices compuestos. | Según la tarea. Depende de `FIX-23`. |

### 3.1 Revisión de las tareas "no implementadas"

Se buscó en todo `main` de ambos submódulos cualquier implementación con otro nombre, ubicación o mecanismo, y también en las ramas remotas no integradas. En el backend no hay ninguna; en el frontend solo existe `feat/stamp-deploy-panel`, del 15/09, que no tiene relación con estas tareas.

| Tarea | Qué se buscó | Resultado |
|---|---|---|
| `SEC-03` | `usePermissions`, `PermissionGate`, `canAccess*`, `hasRole`, `isAdmin`, `useAuthUser`, `RoleGate`, `ProtectedRoute` y todos los `rol ===` del front | **No implementada.** `context/AuthContext.js` (commit `e02b821`) guarda el usuario en un contexto, pero no exporta el contexto ni ningún hook. Su única función es la expulsión por `TOKEN_REVOKED`. El control de rol sigue siendo ad hoc (`Sidebar.tsx:88`, `sessionGuard.ts:22`). |
| `FIX-22` | Guard alternativo dentro de `Auditoria.tsx` (redirect por rol) | **No implementada.** La vista no controla el rol y la ruta sigue bajo `loadProtectedSession` (`applicationRoutes.tsx:50`). |
| `FIX-23` | `trigger`, `REVOKE`, `append_only` en Go y SQL, incluido `scripts/init.sql` | **No implementada.** "append-only" solo aparece en comentarios. |
| `INF-08` / `BAC-16B` | `smtp`, `EMAIL_PROVIDER`, `MAIL_HOST`, `net/smtp`, `gomail`, `mailer` | **No implementadas.** Solo existe `email/mock_email.go`, y `main.go:85` siempre lo instancia. |
| `BAC-17B` | Cambios en `auth_service.go` posteriores al 25/09 | **No implementada.** El último cambio es `fec42b7` (fecha de último acceso). El test confirma 13 filas por sesión. |
| `BAC-18B` | Índice parcial, `PARTITION`, índices compuestos | **No implementada.** Los únicos `fecha_hora DESC` son `ORDER BY` en las consultas. |

**Conclusión:** ninguna de estas tareas se implementó de otra forma. Los tests no necesitan correcciones.

### 3.2 Tareas implementadas que no pasan: ¿problema del test o de la implementación?

| Tarea | ¿El test es correcto? | ¿La implementación cumple la tarea? |
|---|---|---|
| `FIX-20` | **Sí.** Prueba claves que violan una sola regla a la vez y tienen un largo válido. En esta revisión se corrigió un caso que ocultaba el bug con una clave demasiado larga. | **No.** `validatePasswordComplexity` no incluye la condición de mayúscula: la segunda condición prueba `/[0-9]/` con el mensaje "al menos una letra mayúscula". Además, el mensaje de largo dice "Debe tener 8 y 12 caracteres" (le falta "entre"). Las reglas de dígito y símbolo sí funcionan. |
| `FRN-18` | **Sí.** Los 2 casos del entregable 1 pasan. El caso de `canOperateInstance` asume que la sesión del operador trae los niveles por instancia (`permisos: [{ vmid, nivelAcceso }]`), porque la tarea no define de dónde los lee el hook; habrá que confirmarlo cuando se implemente SEC-03. | **Parcialmente.** El entregable 1 (leer y enviar `nivelAcceso`) cumple. El entregable 2 (`canOperateInstance`) no existe. |
| `INF-05` | **Sí.** El CORS se prueba contra el middleware real con `ALLOWED_ORIGINS`. El TLS se omite con el motivo explícito, en lugar de darse por aprobado. | **CORS sí; TLS no verificable.** El criterio de HTTPS depende de la configuración de Nginx del CT 103, que no está versionada. |

**Hito "Gestión Administrativa de Usuarios":** sigue abierto solo por `FIX-27` (mensaje del 502), que está en `futuro.md`.

---

## 4. Regresiones en tareas de `terminado.md`

| Tarea | Pruebas que fallan | Problema | Cómo proceder |
|---|---|---|---|
| **`FRN-13`** Logout en cliente — **confirmada, `FIX-28`** | `session-security.test.ts`: `logoutSession envía … Authorization: Bearer`, y `test/back/session_security…`: `FRN-13 integracion…` | **Confirmada con evidencia de los dos lados:**<br>• **Frontend:** el commit `deb59cb` (Cristian, 25/09, *"Se limpia implementacion vieja de endpoint donde se utiliza body y header"*) cambió `logoutSession()` a `sendJsonPostRequest('/auth/logout', undefined, { skipAuthorization: true })`. Fue un cambio intencional que quitó el Bearer.<br>• **Backend:** `cmd/api/main.go:157` sigue montando `auth.POST("/logout", middleware.RequireAuth(authRepo), …)`, así que sin el Bearer responde `401 MISSING_TOKEN` antes del handler.<br>• **Efecto medido:** después del "logout", el mismo access token sigue respondiendo 200 en `/account/profile`. El usuario no lo nota, porque `useLogout()` limpia el almacenamiento y redirige igual.<br>• No hay ningún otro mecanismo en el frontend que revoque la sesión (`AuthProvider` solo reacciona a un `TOKEN_REVOKED`). | **`FIX-28`** (`futuro.md`, prioridad alta): quitar `skipAuthorization: true` del logout. Aceptar un logout solo con la cookie en el backend no alcanza, porque sin el access token no se puede revocar su JTI (BAC-17). |
| ~~`FIX-26`~~ (SEC-04) | — | **Descartado el 28/09/2026:** el contrato oficial es `{ permisos: [{ vmid, nivelAcceso }] }`. El caso de prueba se reemplazó por `SEC-04 el payload anterior { vmids } se rechaza…`, que pasa: responde 400 `INVALID_REQUEST` sin modificar los permisos. | SEC-04 queda completa en `terminado.md`. |
| `FIX-27` (FIX-24, en `futuro.md`) | `FIX-27 si el correo no se pudo enviar…` | Sin cambios: el 502 muestra el mensaje genérico. | Pendiente de asignar. |

---

## 5. Integración front ↔ back

| Flujo | 25/09 | 28/09 |
|---|---|---|
| Logout desde la interfaz | ✅ | ❌ **Regresión:** `401 MISSING_TOKEN` y la sesión sigue activa (FRN-13) |
| Renovación silenciosa (`/auth/refresh` con cookie) | ✅ | ✅ |
| CORS para el frontend (`Allow-Credentials` para la cookie) | — | ✅ Con el origen en `ALLOWED_ORIGINS` |
| Recuperación de contraseña (`/auth/password/forgot` y `/reset`) | ❌ | ✅ Contrato alineado |
| Reset administrativo (`/admin/users/:id/password/reset` y `/2fa/reset`) | ❌ | ✅ |
| Complejidad de contraseña | ❌ | 🟡 Se frenan las claves sin dígito o sin símbolo; las claves sin mayúscula siguen llegando al backend |

---

## 6. Cambios en las pruebas en esta revisión

| Cambio | Motivo |
|---|---|
| Archivo nuevo `test/back/cierre_fase_base_acceptance_test.go` (INF-05, INF-06, INF-08, BAC-16B, BAC-17B y BAC-18B) | Tareas nuevas de `actual.md` sin cobertura |
| `compose.yaml`: `ALLOWED_ORIGINS`, `EMAIL_PROVIDER=mock`, servicio `backend-smtp` (misma imagen, con `EMAIL_PROVIDER=smtp`) y `mailpit` | Probar CORS y el adaptador SMTP de punta a punta sin un proveedor externo |
| FRN-18: casos de lectura del nivel guardado y de `canOperateInstance` | Entregables sin cobertura |
| FRN-13: la limpieza se prueba con el botón de logout real (204, red caída, 500 y 401) y la expulsión con un 401 `TOKEN_REVOKED` recibido por `apiClient` con `AuthProvider` montado | El equipo movió esas responsabilidades a `useLogout` y a `AuthProvider`. Se exige lo mismo que antes: sesión vacía y `/login` con `replace` |
| FRN-13 (backend): caso de integración con el logout exacto del frontend | Medir el efecto real de la regresión |
| Recuperación de contraseña: textos de la vista nueva ("Verifica tu identidad", "Creá tu nueva contraseña", "Confirmar nueva contraseña", "Contraseña restablecida…", "El código ingresado es incorrecto…") | Cambió el texto de la interfaz; lo que se verifica es lo mismo |
| FIX-20, paso 3: clave de prueba `sinmayus1!` (10 caracteres) en lugar de `sinmayuscula1!` (14) | Con 14 caracteres la frenaba la regla de largo y **ocultaba** el bug de mayúsculas |
| Cambio de contraseña: el mensaje de largo se valida con `/8 y 12 caracteres/` y el caso de error del backend usa una clave que el cliente acepta | Cambió el texto del mensaje, y el cliente ahora frena antes las claves sin símbolo |
| FRN-11: se acepta el toast duplicado en el DOM | El toast se renderiza dos veces; el aviso sí aparece |

---

## 7. Recomendaciones para la documentación

- ~~Terminar el merge del repo principal.~~ **Hecho.**
- ~~Pasar a `terminado.md`: FRN-11, FIX-19, FRN-12, FIX-18, FIX-21 e INF-06.~~ **Hecho el 28/09/2026.**
- ~~Registrar la regresión de FRN-13.~~ **Hecho: `FIX-28` en `futuro.md` y aviso en FRN-13 dentro de `terminado.md`.**
- ~~FIX-20, FRN-18 e INF-05 a `terminado.md` con su FIX.~~ **Hecho: `FIX-29`, `FIX-30` y `FIX-31`.**
- **BAC-16B:** corregir el nombre del método del puerto (`EnviarContrasenaTemporal`, no `EnviarCredencialesTemporales`).

---

## 8. Cómo reproducir

```bash
# 1. Traer el último commit de cada submódulo (prevalece el remoto ante conflictos)
for s in backend frontend; do
  git -C $s fetch origin --prune
  git -C $s checkout -B main origin/main --force
  git -C $s reset --hard origin/main
done

# 2. Backend (Docker: PostgreSQL, 2 backends, stub de Proxmox y Mailpit; ~2 min)
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)

# 3. Frontend
pnpm --dir test/front test
```

Usar `| tee` y no `> archivo` para capturar la salida del backend: en este entorno, redirigir con `>` cortó la corrida y dejó contenedores huérfanos. Los puertos que se usan son `15433`, `18080`, `18081` y `18025`. El comando de limpieza está en `test/back/README.md`.

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
