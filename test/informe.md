# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 01/10/2026
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md) y la regresión de [terminado.md](../documentacion/terminado.md) y [terminado-1.md](../documentacion/terminado-1.md).

| Componente | Revisión probada | Commits nuevos desde la corrida anterior (30/09) |
|---|---|---|
| Backend | `860b3c9` (último commit de `main`) | 1 |
| Frontend | `3d1e84a` (último commit de `main`) | 4 merges (PR #67 a #70) |

**Commits nuevos:**
- **Backend** `860b3c9`: adaptador `SmtpEmailService` con STARTTLS/TLS, elegido con `EMAIL_PROVIDER=smtp` (**BAC-16B**).
- **Frontend** `cfb88f7`: el logout envía Bearer y la cookie HttpOnly (**FIX-28**).
- **Frontend** `3ea7fa4`: mensaje correcto cuando falla el correo al crear un usuario (**FIX-27**).
- **Frontend** `54b397e` … `0b632f1`: hook `usePermissions()` con `canOperateInstance` (**FIX-30**).
- **Frontend** `56b88f5` y `e6ad7ad`: refactor de Auditoría. Se quitaron las columnas Nodo e IP origen, y también **el filtro por acción** (ver §2.1).

**Cómo se ejecutó:** cada submódulo se llevó al último commit de `origin/main`; ante un conflicto prevalece el remoto. El backend se corrió **dos veces completo**, con el mismo resultado y sin contenedores residuales. El frontend también se corrió dos veces, con el mismo resultado.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 44 | **37** (antes 35) | **7** | 0 |
| Frontend (`test/front`) | 95 | **76** (antes 74) | **9** | 10 |
| **Total** | **139** | **113** | **16** | **10** |

Los 10 omitidos del frontend son la prueba integral LOGIN-04, que se ejecuta desde `test/back` contra el backend real. **Esta corrida pasa 10/10.**

**Avance:**
- **Completas:** `BAC-16B`, `FIX-27` y `FIX-28`. `FIX-30` está hecha en el frontend, pero el backend no le da los datos: `FIX-37`, caso nuevo que falla (§2).
- **La prueba integral LOGIN-04 pasa completa** por primera vez: el paso 7 (logout) ya revoca la sesión.
- **Regresión nueva:** la pantalla de Auditoría ya no tiene el filtro por acción (`FRN-14`).

Todos los fallos que quedan son del producto. Algunos tests se tuvieron que adaptar; se explica en la [sección 5](#5-ajustes-a-las-pruebas-en-esta-corrida).

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumplida · 🟡 parcial o implementada con problema · ❌ no implementada.

| Tarea | Pruebas | Estado | Detalle |
|---|---|---|---|
| `BAC-16B` Adaptador SMTP | `cierre_fase_base…` | ✅ **Nueva** | Con `EMAIL_PROVIDER=smtp` y STARTTLS obligatorio, la clave temporal llega por SMTP y permite iniciar sesión, y también llega el código de 6 dígitos. Con el servidor SMTP caído, el alta responde `502 EMAIL_DELIVERY_FAILED` y no se guarda el usuario. Sin `EMAIL_PROVIDER` se sigue usando el mock |
| `FIX-28` Logout con Bearer | `session-security`, `session_security…` y LOGIN-04 paso 7 | ✅ **Nueva** | El frontend envía `Authorization: Bearer` y la cookie (`cfb88f7`). En la prueba integral, el access token queda revocado después de "Cerrar sesión" |
| `FIX-27` Mensaje ante fallo de correo | `admin-users.test.tsx` | ✅ **Nueva** | El mensaje del 502 se muestra tal como lo envía el backend (`3ea7fa4`) |
| `FIX-30` `canOperateInstance` | `navigation.test.tsx` y `resource_access…` (`FIX-37`, nuevo) | 🟡 | En el frontend cumple: `canOperateInstance(vmid)` responde `true` con `FULL_ACCESS` y `false` con `READ_ONLY` o sin asignación, leyendo `perfil.permisos`. **Pero `GET /account/profile` no devuelve `permisos`**, solo `instanciasPermitidas`. Con datos reales, el helper da `false` para todo OPERATOR → **`FIX-37`** (backend) |
| `SEC-03` Permisos reactivos | `navigation.test.tsx` (3) | 🟡 Parcial | Existe `usePermissions()`, pero en `src/hooks/` y no en `src/context/`, como pide el entregable. Le faltan `isOperator` y `hasRole`; no existe `PermissionGate`, y el menú sigue mostrando "Auditoría" al OPERATOR |
| `FIX-29` Complejidad de contraseña | `password-change`, `recover-password` | ❌ | Sin cambios: se valida `/[0-9]/` donde corresponde validar mayúsculas |
| `FRN-17C` Cliente de eventos | `events-client.test.tsx` (3) | ❌ | No existe `useEvents`. El backend ya ofrece `/api/events/ticket` y `/api/events` (BAC-21C) |
| `INF-08B` → `FIX-34` | `cierre_fase_base…` | 🟡 (ya en `terminado.md`) | `SMTP_USER` y `SMTP_PASS` siguen con valores de ejemplo |
| `BAC-18B` Índice parcial y particiones | `cierre_fase_base…` | ❌ | Sin índice parcial sobre `jti_access`, sin particiones y sin índices compuestos |
| `BAC-21B` Instancias extendidas | `puente_etapa1…` (3) | ❌ | Faltan `ip`, `cpuUsage`, `ramUsage`, `maxRam`, `nivelAcceso` y `activeTask`. `/status/:action` responde 404 y `DELETE` responde 405 |
| `FIX-31` TLS de Nginx versionado | `cierre_fase_base…` | ❌ | Ninguna configuración de Nginx versionada tiene `listen 443 ssl` |
| `INF-06B`, `INF-08A` | Sin prueba automatizada | — | Configuración del servidor |

### 2.1 Regresión: filtro por acción en Auditoría (`FRN-14`, en `terminado.md`)

`56b88f5` ("eliminar columnas que no serán utilizadas") sacó de `src/pages/Auditoria.tsx`:
- las columnas *Nodo* e *IP origen*;
- el campo **"Acción"** y el parámetro `accion` de la consulta.

Las columnas eran opcionales, pero el filtro lo exige **RF-08**: [requerimientos.md:193](../documentacion/requerimientos.md#L193) dice "filtros por fecha, usuario, tipo de acción y resultado". El backend lo sigue aceptando (BAC-18 pasa).

**Prueba que falla:** `audit.test.tsx`, *"actualiza los query parameters de filtro al cambiar los selectores reactivos"*. Antes pasaba (7/7) y ahora da 6/7.

**Cómo seguir:** restaurar el filtro "Acción". Si se decide que no va, primero hay que cambiar RF-08. Corresponde agregar un **FIX** en `futuro.md` y una advertencia en `FRN-14`.

---

## 3. Prueba integral LOGIN-04 (frontend real contra backend real)

**10 de 10 pasos en verde.** El paso 7 ya funciona: "Cerrar sesión" revoca la sesión en el servidor (`FIX-28`). El paso 8 (recuperación de contraseña) también pasa, después de adaptar el test al texto nuevo del correo simulado (§5).

---

## 4. Simulador de Proxmox

El simulador no cambió desde la corrida anterior (`860b3c9` solo toca el correo), así que no repetí la comparación con el Proxmox real.

**Sigue vigente:**
- `FIX-32` y `FIX-33` completas (`terminado-1.md`).
- `BAC-28` con diferencias en `lxc/{vmid}/interfaces`: `prefix` es número y un LXC apagado responde 500. Su corrección es `FIX-35`, en `futuro-1.md`.

---

## 5. Ajustes a las pruebas en esta corrida

Ninguno oculta un fallo. Todos alinean el test con lo que pide la tarea:

| Cambio | Motivo |
|---|---|
| Mailpit (`compose.yaml`) con certificado de prueba y **STARTTLS obligatorio** (`MP_SMTP_REQUIRE_STARTTLS`). `backend-smtp` confía en esa CA (`SSL_CERT_FILE`, `test/back/mailpit-tls/`) | El adaptador, como corresponde, no manda credenciales por una conexión sin cifrar, y el Mailpit anterior no ofrecía TLS: el alta daba 502. El relay real (puerto 587) usa STARTTLS. Ahora el test además comprueba que el envío vaya cifrado |
| `passwordCandidates` toma la clave de "contraseña … es:" en el cuerpo sin limpiar | La clave temporal puede tener `@`, `*`, `<`, `=`, etc. El test la descartaba o la recortaba, así que fallaba según la clave que tocara |
| La marca del código de recuperación pasa a `"código de seguridad es:"` (`main_test.go` y `login04-e2e.test.ts`) | `860b3c9` cambió el texto que imprime el `MockEmailService`. Es un detalle interno del mock, no un contrato. BAC-19, BAC-20, FIX-17 y el paso 8 de LOGIN-04 fallaban solo por eso |
| La integración FIX-28 (`session_security…`) envía el logout con Bearer y cookie | Simulaba el frontend anterior (sin Bearer). Ahora replica lo que hace `cfb88f7`; que el frontend mande el header lo verifica `session-security.test.ts` |
| El caso FIX-30 busca `usePermissions` en `src/context` **o** `src/hooks`, sin exigir `PermissionGate` | FIX-30 solo pide el helper `canOperateInstance`. La ubicación en `src/context`, `isOperator`, `hasRole` y `PermissionGate` los siguen exigiendo los casos de SEC-03 |
| Caso nuevo `FIX-37…` en `resource_access…`: `GET /account/profile` debe informar `permisos [{ vmid, nivelAcceso }]` | La prueba de FIX-30 siembra la sesión a mano; este caso verifica que el backend real entregue el dato. Hoy falla |

---

## 6. Pendientes y observaciones

- **Reclasificación (01/10/2026):**
  - Pasaron a `terminado.md`: `BAC-16B`, `FIX-27` y `FIX-28` (completas) y `FIX-30` (con problema).
  - FIX nuevos en `futuro.md`: `FIX-36` (filtro "Acción", §2.1), `FIX-37` (`permisos` en el perfil) y `FIX-38` (lo ya hecho de `SEC-03`: ubicación del hook y "READ_ONLY" tratado como rol).
  - `FIX-35` (Etapa 1) sigue en `actual.md`: el simulador no cambió.
- **Frontend:** quedan `SEC-03` (ubicación, `isOperator`, `hasRole`, `PermissionGate` y menú), `FIX-29` y `FRN-17C`.
- **Backend:** quedan `BAC-18B`, `BAC-21B`, `FIX-31` y `FIX-34` (credenciales reales en `.env.example`).
- **Seguridad:**
  - El secreto del token de Proxmox sigue en texto plano en `documentacion/api-proxmox.md`.
  - Cuando se cierre `FIX-34`, las credenciales SMTP quedarán versionadas en el repositorio del backend.
  - El certificado de `test/back/mailpit-tls/` es solo para pruebas: se creó con CN `mailpit`, sin la clave de la CA y no vale fuera de este entorno.

---

## 7. Cómo reproducir

```bash
for s in backend frontend; do
  git -C $s fetch origin --prune && git -C $s checkout -B main origin/main --force && git -C $s reset --hard origin/main
done
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)   # incluye LOGIN-04 front ↔ back
pnpm --dir test/front test
```

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
