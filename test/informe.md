# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 29/09/2026
**Alcance:** tareas en desarrollo de [documentacion/actual.md](../documentacion/actual.md), correcciones pendientes de [documentacion/futuro.md](../documentacion/futuro.md) y regresión de [documentacion/terminado.md](../documentacion/terminado.md).

| Componente | Revisión probada | Cambios desde el informe anterior (28/09) |
|---|---|---|
| Backend | `9554efa` (último commit de `main`) | `ad52485` y `9554efa`: trigger de inmutabilidad `trg_auditoria_inmutable` (`BEFORE UPDATE OR DELETE OR TRUNCATE`) y `REVOKE` sobre `auditoria`. La FK de `auditoria.usuario_id` pasa a `ON DELETE RESTRICT` (**FIX-23**) |
| Frontend | `749194e` (último commit de `main`) | `5c3d780`: `/auditoria` se movió dentro del grupo `loadAdminSession` (**FIX-22**). `bd0c42f`: helper `services/request.ts` (hoy no se usa: solo aparece en código comentado de `Auditoria.tsx`) |

**Cómo se ejecutó:** cada submódulo se actualizó a su último commit de `origin/main` (`fetch` + `checkout -B main origin/main --force` + `reset --hard`). Ante un conflicto prevalece el remoto. Las suites se corrieron sobre esas carpetas **dos veces** cada una, con el mismo resultado. Las pruebas no necesitaron cambios en esta revisión.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 38 | 32 | **5** | 1 |
| Frontend (`test/front`) | 82 | 74 | **8** | 0 |
| **Total** | **120** | **106** | **13** | **1** |

**Avances desde el 28/09:**
- **FIX-22** (`/auditoria` protegida) y **FIX-23** (auditoría append-only) pasan todas sus pruebas.
- Con eso quedan resueltos los problemas pendientes de **FRN-14** y **BAC-18**, dos tareas que ya estaban en `terminado.md`.

**Sigue abierta la regresión de FRN-13 (`FIX-28`):** el logout del frontend sale sin Bearer y la sesión sigue activa en el servidor.

**Todos los fallos son del producto.** El caso omitido es el TLS de INF-05 (`FIX-31`), porque la configuración de Nginx del servidor no está versionada.

---

## 2. Tareas de `actual.md`

**Leyenda:** ✅ cumplida · ❌ no implementada.

| Tarea | Pruebas | Resultado | Estado | Detalle |
|---|---|---|---|---|
| `FIX-22` Guard de `/auditoria` | `audit.test.tsx`: `si un operador intenta entrar a /auditoria…` | 1/1 | ✅ | La ruta está bajo `loadAdminSession` y el OPERATOR vuelve a `/dashboard`. Los otros 6 casos de FRN-14 siguen en verde. **Movida a `terminado.md`.** |
| `FIX-23` Auditoría append-only | `password_recovery…`: `BAC-18 … append-only` | 1/1 | ✅ | `UPDATE`, `DELETE` y `TRUNCATE` sobre `auditoria` fallan en PostgreSQL con las credenciales de la aplicación. El registro, la consulta, los filtros y el CSV siguen funcionando, y ningún otro caso se rompió por el trigger ni por la FK `RESTRICT`. **Movida a `terminado.md`.** |
| `SEC-03` Contexto de permisos | `navigation.test.tsx`, bloque SEC-03 (4) | 1/4 | ❌ | No existen `usePermissions` ni `PermissionGate` (`context/AuthContext.js` solo maneja `TOKEN_REVOKED`). El sidebar le sigue mostrando **el enlace "Auditoría" al OPERATOR** (`Sidebar.tsx`, fuera de `rol === 'ADMIN'`). Desde FIX-22 el operador ya no puede entrar a la vista, pero el acceso sigue visible. |
| `INF-08` Variables SMTP | `cierre_fase_base…`: `INF-08…` | 0/1 | ❌ | `EMAIL_PROVIDER`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` y `SMTP_FROM` no figuran en `backend/.env.example` ni en `backend/docker-compose.yml`. |
| `BAC-16B` Adaptador SMTP | `cierre_fase_base…`: `BAC-16B…` (backend con `EMAIL_PROVIDER=smtp` + Mailpit) | 0/1 | ❌ | El alta no envía ningún correo SMTP: `cmd/api/main.go` siempre instancia `email.NewMockEmailService()`. |
| `BAC-17B` Una sesión, un registro | `cierre_fase_base…`: `BAC-17B…` | 0/1 | ❌ | Login + 2FA + 10 refresh crean **13 filas activas**, y después del logout **quedan 11**. |
| `BAC-18B` Índice parcial y particionado | `cierre_fase_base…`: `BAC-18B…` | 0/1 | ❌ | Faltan `idx_sesiones_activas_vigentes`, el particionado de `auditoria` (sigue siendo `relkind = r`, sin particiones) y los índices compuestos. Su dependencia `FIX-23` ya está resuelta. |

---

## 3. Correcciones pendientes de `futuro.md`

| FIX | Tarea de origen | Pruebas | Resultado | Detalle |
|---|---|---|---|---|
| **`FIX-28`** (prioridad alta) | FRN-13 (regresión) | `session-security.test.ts`: `FIX-28 logoutSession…` y `session_security…`: `FIX-28 FRN-13 integracion…` | ❌ 0/2 | `logoutSession()` sigue usando `skipAuthorization: true` (commit `deb59cb`). El backend responde `401 MISSING_TOKEN` y el access token sigue válido después de "Cerrar sesión". |
| `FIX-27` | FIX-24 / FRN-06 | `admin-users.test.tsx`: `FIX-27 si el correo no se pudo enviar…` | ❌ 1/2 | Ante un `502 EMAIL_DELIVERY_FAILED` se sigue mostrando el mensaje genérico. El otro caso (el formulario conserva los datos) pasa. |
| `FIX-29` | FIX-20 / FRN-10 | `password-change…` (sin mayúscula) y `recover-password…` (paso 3) | ❌ | `validatePasswordComplexity` sigue sin verificar la mayúscula. Los casos sin dígito, sin símbolo y fuera del largo pasan. |
| `FIX-30` | FRN-18 | `navigation.test.tsx`: `FIX-30 canOperateInstance…` | ❌ 0/1 | Depende de `SEC-03`, que no está implementada. |
| `FIX-31` | INF-05 | `cierre_fase_base…`: `FIX-31 INF-05 TLS…` | ⏭️ omitido | La configuración de Nginx del servidor sigue sin versionar. |

---

## 4. Regresión de `terminado.md`

El resto de las tareas de `terminado.md` pasa todas sus pruebas. Las implementadas con problema tienen su corrección en la sección 3, y FRN-14 y BAC-18 quedaron completas con FIX-22 y FIX-23.

| Tarea | Estado al 29/09 | Nota para la documentación |
|---|---|---|
| `FRN-14` | ✅ 7/7 | El aviso de `terminado.md` (el guard faltante → FIX-22) quedó desactualizado. |
| `BAC-18` | ✅ | El aviso de `terminado.md` (no es append-only → FIX-23) quedó desactualizado. |
| `FRN-13` | ❌ regresión | `FIX-28`, sin cambios. |
| `INF-03` e `INF-04` | Sin prueba automatizada | Son configuración del servidor sin artefactos versionados. |

---

## 5. Integración front ↔ back

| Flujo | 28/09 | 29/09 |
|---|---|---|
| Logout desde la interfaz | ❌ `401 MISSING_TOKEN`, la sesión sigue activa | ❌ Sin cambios (`FIX-28`) |
| Renovación silenciosa con la cookie `centinela_refresh` | ✅ | ✅ |
| CORS (`ALLOWED_ORIGINS`, `Allow-Credentials`) | ✅ | ✅ |
| Recuperación y reset administrativo de contraseña y 2FA | ✅ | ✅ |
| Permisos por instancia `{ permisos: [{ vmid, nivelAcceso }] }` | ✅ | ✅ |
| Acceso del OPERATOR a Auditoría | ❌ Entraba a la vista | ✅ La ruta lo redirige (el backend ya respondía 403). 🟡 El enlace del menú sigue visible (`SEC-03`) |
| Complejidad de contraseña | 🟡 | 🟡 Las claves sin mayúscula siguen llegando al backend (`FIX-29`) |

---

## 6. Recomendaciones para la documentación

- ~~Pasar `FIX-22` y `FIX-23` a `terminado.md` y actualizar los avisos de `FRN-14` y `BAC-18`.~~ **Hecho el 29/09/2026:** las cuatro figuran como completadas en `terminado.md`.
- **Prioridad de asignación:** `FIX-28` sigue siendo la corrección más urgente (media hora, frontend).
- **`BAC-18B`:** ya tiene resuelta su dependencia `FIX-23`. Al particionar `auditoria`, el trigger de inmutabilidad tiene que aplicarse también a las particiones; el test de BAC-18 lo va a verificar.

---

## 7. Lo que la suite no cubre

- `502 PROXMOX_UNAVAILABLE`: el stub de Proxmox siempre responde.
- **INF-05:** el TLS de Nginx en el servidor (`FIX-31`).
- **INF-06:** la red `vmbr1`.
- **INF-08:** la conectividad hacia un SMTP real.
- **BAC-18B:** la purga horaria y el tiempo de consulta menor a 20 ms.
- **INF-03 e INF-04:** la configuración del servidor.

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

Usar `| tee` y no `> archivo` para capturar la salida del backend. Los puertos que se usan son `15433`, `18080`, `18081` y `18025`. El comando de limpieza está en `test/back/README.md`.

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
