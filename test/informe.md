# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 01/10/2026 (segunda verificación del día; backend vuelto a correr después de ajustar FIX-31, INF-08B y BAC-21B, §4)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), después de completar las pruebas que les faltaban, y la regresión de [terminado.md](../documentacion/terminado.md) y [terminado-1.md](../documentacion/terminado-1.md).

| Componente | Revisión probada | Commits nuevos desde la verificación anterior |
|---|---|---|
| Backend | `44a2339` (último commit de `main`) | 1: `44a2339` "mas arreglos del simulador" (**FIX-35**) |
| Frontend | `3d1e84a` (último commit de `main`) | Ninguno |

**Cómo se ejecutó:** cada submódulo se llevó al último commit de `origin/main`; ante un conflicto prevalece el remoto. Cada suite se corrió **dos veces**, con el mismo resultado y sin contenedores residuales.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 50 | **43** | **6** | 1 |
| Frontend (`test/front`) | 100 | 76 | **14** | 10 |
| **Total** | **150** | **119** | **20** | **11** |

- Los 10 omitidos del frontend son la prueba integral LOGIN-04, que corre desde `test/back` contra el backend real (pasa 10/10). El omitido del backend es el `DELETE` de BAC-21B, que todavía no existe (§4).
- Hay **8 casos nuevos** (§3): 3 del backend y 5 del frontend. Fallan todos salvo los 2 de `FIX-35`, porque cubren partes de tareas que todavía no están implementadas.
- **Avance:** `FIX-35` (Etapa 1) está **completo** y pasó a `terminado-1.md`. Las demás tareas de `actual.md` no tuvieron commits.
- **Caso nuevo de la Etapa 1, `INF-07B`:** `/api/events` a través del Nginx del borde. **Pasa:** la configuración versionada tiene el stream sin buffer y con timeout de 1 h, y un `TASK_FINISHED` real llega por HTTPS en ~1 s. Para eso, el stub de Proxmox ahora responde el estado de las tareas. La parte del servidor (CT 103) queda para `INT-03`.
- **`FIX-31` e `INF-08B` ya pasan** (§4):
  - `FIX-31`: el borde con TLS ahora está versionado en este repositorio y se prueba en funcionamiento.
  - `INF-08B`: el rechazo de Brevo dependía de la IP de la red.

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumplida · 🟡 parcial · ❌ no implementada.

| Tarea | Pruebas | Estado | Detalle |
|---|---|---|---|
| `FIX-35` Simulador: IP de LXC (Etapa 1) | `simulador_proxmox…` (2, **nuevo**) y `go test ./cmd/proxmox-simulador` | ✅ **Nueva** → `terminado-1.md` | `prefix` como string en `lxc/{vmid}/interfaces`, y un LXC apagado responde `200 {"data": null}`, igual que el Proxmox real. El guest agent de qemu se revisó contra la especificación QAPI y mantiene `prefix` numérico. 3 pruebas unitarias nuevas en el simulador |
| `SEC-03` Permisos reactivos | `navigation.test.tsx` (3) | 🟡 | `usePermissions` en `src/hooks/` con `isAdmin`, `canAccessInstance` y `canOperateInstance`. Faltan `isOperator`, `hasRole` y `PermissionGate`, y el menú sigue mostrando "Auditoría" al OPERATOR. No hay header que integrar: está comentado en `MainLayout.tsx` |
| `FIX-29` Complejidad de contraseña | `password-change` (3, 2 **nuevos**), `recover-password` (1) | ❌ | No se valida la mayúscula. Los casos nuevos verifican los mensajes: "al menos un número" y "entre 8 y 12 caracteres" |
| `FRN-17C` Cliente de eventos | `events-client.test.tsx` (4, 1 **nuevo**) | ❌ | No existe `useEvents`. El caso nuevo verifica el retroceso exponencial entre reconexiones |
| `FIX-36` Filtro "Acción" en Auditoría | `audit.test.tsx` (1) | ❌ | Sigue sin el filtro |
| `FIX-38` `READ_ONLY` como rol | `admin-users` (1, **nuevo**), `navigation` (1, **nuevo**) | ❌ | El selector "Rol" de la ficha ofrece `ADMIN`, `OPERATOR` y `READ_ONLY`; `canAccessInstance` da `true` a un usuario con rol `READ_ONLY` |
| `BAC-18B` Índice parcial, particiones y purga | `cierre_fase_base…` (2, 1 **nuevo**; la de índices se reescribió, §4) | ❌ | Faltan el índice parcial, las particiones trimestrales y el índice `(fecha_hora, accion, resultado)`. El de `(usuario_id, fecha_hora)` ya existe. No hay purga horaria |
| `BAC-21B` Instancias extendidas | `puente_etapa1…` (5) | 🟡 → **`terminado.md`** con advertencia; lo que falta está en **`FIX-39`** (`futuro.md`) | **Cumple:** `start` y `stop` con `FULL_ACCESS` y `tareas_asincronas`. **Falta (`FIX-39`):** los campos de `GET /instances` (pueden ir en `null` salvo `nivelAcceso`), la auditoría de la orden despachada, y `shutdown`/`reboot` (caso nuevo). `DELETE`: se omite hasta `BAC-24B` |
| `FIX-37` Nivel de acceso en el perfil | `resource_access…` (1) | ❌ | `GET /account/profile` sigue sin `permisos` |

---

## 3. Pruebas que faltaban (creadas en esta verificación)

Se revisó cada tarea de `actual.md` contra su entregable y su criterio de éxito. Faltaban estas pruebas:

| Tarea | Prueba nueva | Qué verifica |
|---|---|---|
| `FIX-35` | `test/back/simulador_proxmox_acceptance_test.go` (2 casos) | Compila el simulador desde el submódulo (sin modificarlo), lo levanta en un puerto libre y consulta `lxc/101/interfaces` (`prefix` string) y `lxc/201/interfaces` apagado (`200 {"data": null}`) |
| `BAC-18B` | `cierre_fase_base…`, caso *"rutina horaria que purga…"* | Que el backend tenga el borrado de `sesiones_activas` por `activa = false` o `fecha_expiracion <`, con SQL o con GORM, y un ticker de 1 hora. Se revisa en el código porque la suite no puede esperar una hora |
| `BAC-21B` | `puente_etapa1…`, reescritas en §4 | Criterio de éxito, sin atarse a la forma sugerida |
| `FIX-29` | `password-change.test.tsx`, 2 casos | El mensaje de dígito ("al menos un número") y el de largo ("entre 8 y 12 caracteres"). Antes la prueba de largo aceptaba el texto incorrecto "8 y 12" |
| `FRN-17C` | `events-client.test.tsx`, caso *"retroceso exponencial"* | Dos desconexiones seguidas: la segunda espera es al menos 1,5 veces la primera |
| `FIX-38` | `admin-users.test.tsx` y `navigation.test.tsx` | El selector "Rol" ofrece solo `ADMIN` y `OPERATOR`; un rol `READ_ONLY` no da acceso a instancias |

Las demás tareas ya tenían pruebas que cubren su criterio. La duda sobre `instancesSummary` quedó resuelta con la aclaración del 01/10/2026: va en `GET /api/node/status` (`BAC-22B`), no en cada instancia.

---

## 4. Ajustes posteriores: FIX-31, INF-08B y BAC-21B

| Tarea | Qué pasaba | Qué se hizo | Resultado |
|---|---|---|---|
| `FIX-31` (en `terminado.md`) | La prueba buscaba la configuración de Nginx en los submódulos, donde nunca va a estar | Se versionó el borde en este repositorio, **`docker/nginx-edge.conf`**: 443 con TLS 1.2 y 1.3, HSTS, 80 → 301 y `/api/` hacia el backend sin buffer, para el SSE. Se agregó el servicio `edge` a `test/back/compose.yaml`, con un certificado de prueba (`edge-tls/`). La prueba revisa el archivo **y** lo usa: HTTPS 200 desde el backend, HTTP 301 hacia `https://` y TLS 1.1 rechazado | ✅ |
| `INF-08B` (en `terminado.md`) | `525 Unauthorized IP address` en una corrida. La red de esta máquina sale por DHCP/NAT y su IP pública puede cambiar; la cuenta de Brevo restringe por IP | Con la IP actual (`179.238.41.248`), 5 de 5 autenticaciones dieron OK. **Era un problema del test:** fallaba por la red y no por las credenciales. Ahora reintenta ante errores de red o DNS y, ante un `525`, omite el caso con el motivo | ✅ |
| `BAC-21B` (en `actual.md`) | La prueba exigía la forma sugerida en el entregable (`/status/:action`, `action`/`resource_type` en `detalles`, `DELETE`) | Se reescribió contra el **criterio de éxito**: la energía se acepta por `/status/:action` **o** por `/start` y `/stop`; la auditoría solo exige una fila de la instancia con el `upid`, con cualquier forma; el `DELETE` se omite mientras no exista, porque lo crea `BAC-24B` | 1 ✅, 2 ❌ (campos de `GET` y auditoría de energía), 1 omitido |

| `BAC-18B` (en `actual.md`) | La prueba exigía los nombres sugeridos (`idx_sesiones_activas_vigentes`, `auditoria_2026_q3`…) y el orden exacto de columnas | Se reescribió contra el criterio: cualquier índice parcial que cubra `jti_access`; particionado `RANGE (fecha_hora)` con tramos de 89 a 93 días para el trimestre actual y el siguiente, más un `DEFAULT`; índices compuestos que empiecen por `fecha_hora` (con `accion` y `resultado`) y por `usuario_id` (con `fecha_hora`). Se validó con un esquema correcto en un Postgres temporal | ❌ (el producto no lo implementó; el índice por usuario ya existe) |

**Movimiento de `BAC-21B`:** pasó a `terminado.md` como implementada con problema. Lo que falta, junto con la aclaración del 01/10/2026, quedó en **`FIX-39`**, en `futuro.md`:
- campos en `null` salvo `nivelAcceso`;
- `instancesSummary` en `GET /api/node/status` (`BAC-22B`);
- `shutdown` y `reboot`;
- la auditoría de la orden despachada (`PENDING`).

Se sumó el caso *"FIX-39 … shutdown y reboot exigen FULL_ACCESS…"*, y el de `GET` ahora exige el valor real de `nivelAcceso`.

**Lo que el backend hizo en BAC-21B, de otra forma:** `start` y `stop` con `FULL_ACCESS`, seguimiento en `tareas_asincronas` con `tareaId` y `TASK_FINISHED`, y permisos sobre `permisos_instancia`.

**Lo que no hizo:**
- `GET /instances` sigue con los 5 campos de BAC-14.
- Ninguna acción de energía se escribe en `auditoria`; solo se auditan la apertura y el cierre del stream de eventos.

---

## 5. Prueba integral LOGIN-04

**10 de 10 pasos en verde** (frontend real contra backend real), igual que en la verificación anterior.

---

## 6. Pendientes

- **Frontend** (sin commits): `SEC-03`, `FIX-29`, `FRN-17C`, `FIX-36` y `FIX-38`.
- **Backend:** `BAC-18B`, `FIX-37` y `FIX-39` (lo que falta de `BAC-21B`).
- **Infraestructura:** tomar `docker/nginx-edge.conf` como referencia para el CT 103 (`FIX-31` e `INF-07B`).
- **Seguridad:**
  - El secreto del token de Proxmox sigue en `documentacion/api-proxmox.md`.
  - Las credenciales SMTP están versionadas en `test/back/smtp-brevo.env`, por decisión del equipo; para la entrega se reemplazan.

---

## 7. Cómo reproducir

```bash
for s in backend frontend; do
  git -C $s fetch origin --prune && git -C $s checkout -B main origin/main --force && git -C $s reset --hard origin/main
done
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)   # incluye LOGIN-04 front ↔ back y el simulador
pnpm --dir test/front test
```

Si el puerto 18080 está ocupado: `BACKEND_TEST_API_PORT=18090 BACKEND_TEST_API_URL=http://127.0.0.1:18090/api`.

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
