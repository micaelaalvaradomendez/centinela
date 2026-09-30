# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 30/09/2026 (segunda corrida del día)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), incluidos `FIX-32` y `FIX-33`, y la regresión de [documentacion/terminado.md](../documentacion/terminado.md).

| Componente | Revisión probada | Commits nuevos desde la corrida anterior |
|---|---|---|
| Backend | `43a0b06` (último commit de `main`) | 8 commits (ver detalle abajo) |
| Frontend | `749194e` (último commit de `main`) | Ninguno |

Commits nuevos del backend:
- `9e0f43c`: 1 sesión = 1 fila en `sesiones_activas` y sesiones efímeras en Redis con TTL (**BAC-17B**).
- `546e1d5`: credenciales de Redis local unificadas y `REDIS_DB` (**INF-06A**).
- `0aac034`: puerto `KeyValueStore` con adaptador Redis y Pub/Sub (**BAC-17A**).
- `233e804`: variables SMTP en `.env.example` (**INF-08B**).
- `45ecad7`: SSE en `/api/events` con ticket efímero, bus Pub/Sub y corte en vivo (**BAC-21C**).
- `c997398`, `2c56bab` y `43a0b06`: correcciones del simulador de Proxmox (**FIX-32**) y `409 INSTANCE_BUSY` (**FIX-33**).

**Cómo se ejecutó:** cada submódulo se actualizó a su último commit de `origin/main`; ante un conflicto prevalece el remoto. Cada suite se corrió **dos veces**, con el mismo resultado y sin contenedores residuales.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 43 | **35** | **8** | 0 |
| Frontend (`test/front`) | 95 | 74 | **11** | 10 |
| **Total** | **138** | **109** | **19** | **10** |

Los 10 omitidos del frontend son la prueba integral LOGIN-04, que corre desde `test/back` contra el backend real.

**Avance: el backend completó 6 tareas** (BAC-17A, BAC-17B, BAC-21C, INF-06A, FIX-33 y FIX-32), y pasa de 30 a 35 casos aprobados. El frontend no tuvo cambios.

**Todos los fallos son del producto.** Los ajustes que hubo que hacer en las pruebas se explican en la [sección 5](#5-ajustes-a-las-pruebas-en-esta-corrida).

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumplida · 🟡 parcial o implementada con problema · ❌ no implementada.

| Tarea | Pruebas | Estado | Detalle |
|---|---|---|---|
| `BAC-17A` Adaptador de Redis | `cierre_fase_base…` | ✅ **Nueva** | `go-redis`, adaptador `secondary/redis`, puerto `KeyValueStore` con `GetDel`/`Publish`/`Subscribe`, y el backend conectado a Redis al arrancar. Si Redis no está, arranca en modo degradado en memoria (`docs/redis.md`) |
| `BAC-17B` Una sesión = un registro | `cierre_fase_base…` | ✅ **Nueva** | Login + 2FA + 10 refresh = 1 fila; después del logout no quedan filas activas. Además, cada refresh rota el access token y el anterior deja de valer |
| `BAC-21C` Tickets y `/api/events` | `puente_etapa1…` (2) | ✅ **Nueva** | El ticket vive en Redis con TTL de 30 s o menos y es de un solo uso (un ticket inválido o reutilizado recibe 401). El stream SSE se corta al hacer logout |
| `INF-06A` Redis local | `cierre_fase_base…` | ✅ **Nueva** | Publicado en `127.0.0.1:6379`, `REDIS_ADDR=localhost:6379` y la misma contraseña en el compose y en `.env.example` |
| `FIX-33` Bloqueo → 409 | Verificación contra el simulador (§4) | ✅ **Nueva** | Un `stop` inmediatamente después de un `start` responde `409 INSTANCE_BUSY` ("La instancia se encuentra ejecutando otra tarea") |
| `FIX-32` Fidelidad del simulador | Comparación contra el Proxmox real (§4) | ✅ **Nueva** | `unprivileged` numérico, `ha` en `status/current`, `/cluster/nextid`, `/nodes/{node}/tasks` con los mismos campos que el real y 401 sin cuerpo. `tasks` lista solo tareas terminadas, igual que el real; la primera medición la vi vacía porque consulté con la tarea en curso |
| `BAC-28` Simulador para la Etapa 1 | Verificación del simulador y comparación del endpoint de IP con el real | 🟡 | `DELETE` (UPID `qmdestroy`, 500 si está encendida) y guest agent funcionan. En `lxc/{vmid}/interfaces`, `prefix` es número (en el real, string) y un LXC apagado responde 500 (en el real, `200 {"data":null}`) → **`FIX-35`** (`futuro-1.md`) |
| `INF-08B` Credenciales SMTP → `FIX-34` | `cierre_fase_base…` | 🟡 | Las 6 variables ya están en `.env.example` (`233e804`), pero `SMTP_USER` y `SMTP_PASS` tienen **valores de ejemplo** (`tu_correo@ejemplo.com`, `tu_clave_secreta_aqui`). Con eso no se puede autenticar contra `smtp-relay.brevo.com` |
| `BAC-16B` Adaptador SMTP | `cierre_fase_base…` | ❌ | Con `EMAIL_PROVIDER=smtp` sigue sin salir ningún correo por SMTP |
| `BAC-18B` Índice parcial y particiones | `cierre_fase_base…` | ❌ | No hay `idx_sesiones_activas_vigentes` (ahora sobre `jti_access`), ni particiones, ni índices compuestos |
| `BAC-21B` Instancias extendidas | `puente_etapa1…` (3) | ❌ | `GET /instances` no trae `ip`/`cpuUsage`/`ramUsage`/`maxRam`/`nivelAcceso`/`activeTask`; `/status/:action` → 404; `DELETE` → 405. **Nota:** `start` ya devuelve `tareaId` además del `upid` |
| `FIX-31` TLS de Nginx versionado | `cierre_fase_base…` | ❌ | No hay ninguna configuración de Nginx con `listen 443 ssl` |
| `FIX-28` Logout sin Bearer | `session-security`, `session_security…` y LOGIN-04 paso 7 | ❌ | El frontend no cambió: la sesión sigue activa en el servidor después de "Cerrar sesión" |
| `SEC-03`, `FIX-30` | `navigation.test.tsx` | ❌ | Frontend sin cambios |
| `FIX-27`, `FIX-29` | `admin-users`, `password-change` y `recover-password` | ❌ | Frontend sin cambios |
| `FRN-17C` Cliente de eventos | `events-client.test.tsx` (3) | ❌ | No existe `useEvents`. **El backend ya expone `/api/events/ticket` y `/api/events` (BAC-21C), así que FRN-17C está desbloqueada** |
| `INF-06B`, `INF-08A` | Sin prueba automatizada | — | Configuración del servidor |

**Reclasificación (30/09/2026):**
- Pasaron a `terminado.md`: `BAC-17A`, `BAC-17B`, `BAC-21C` e `INF-06A` (completas), e `INF-08B` (con problema; su corrección es `FIX-34` en `futuro.md`).
- Pasaron a `terminado-1.md`: `FIX-32` y `FIX-33` (completas), y `BAC-28` (con problema; su corrección es `FIX-35` en `futuro-1.md`).

---

## 3. Prueba integral LOGIN-04 (frontend real contra backend real)

9 de 10 pasos en verde, igual que antes, ahora también con las sesiones nuevas de BAC-17B. **Sigue fallando el paso 7** (`FIX-28`): después de "Cerrar sesión", el access token sigue válido.

---

## 4. Simulador de Proxmox (FIX-32 y FIX-33)

Se compararon el simulador de `43a0b06` y el Proxmox real, con las credenciales de `api-proxmox.md`. Sobre el real se hicieron solo lecturas.

| Endpoint | 30/09 (1ª corrida) | Ahora |
|---|---|---|
| `cluster/resources`, `nodes/{node}/status`, `lxc`, `rrddata` | ✅ | ✅ |
| `lxc/{vmid}/status/current` | ❌ faltaba `ha` | ✅ |
| `lxc/{vmid}/config` | ❌ `unprivileged` era texto | ✅ Mismos tipos (las claves que difieren son propias de cada contenedor) |
| `cluster/nextid` | ❌ 501 | ✅ |
| `nodes/{node}/tasks` | ❌ 501 | ✅ `{total, data[]}` con los mismos campos que el real; lista las tareas terminadas |
| `lxc/{vmid}/interfaces` (BAC-28) | — | ❌ `prefix` es número (en el real, string); LXC apagado → 500 (en el real, `200 {"data":null}`) → `FIX-35` |
| 401 con token inválido | 🟡 con cuerpo JSON | ✅ `401 Authentication failed!` sin cuerpo, igual que el real |

**Backend contra el simulador:** el listado, la protección de VMIDs y el `start` con UPID y `tareaId` funcionan. Un segundo comando mientras hay una tarea en curso ahora responde `409 INSTANCE_BUSY` (**FIX-33 ✅**).

---

## 5. Ajustes a las pruebas en esta corrida

Hicieron falta por el cambio de esquema de BAC-17B y no ocultan fallos:

| Cambio | Motivo |
|---|---|
| `signedAccessToken` crea la sesión con `id = sid`, `jti_access` y `jti_refresh`, y agrega el claim `sid` | `sesiones_activas` ya no tiene `jti_token`. La sesión se identifica por el claim `sid` y el middleware valida `jti_access` |
| BAC-17 busca por `jti_access` | Cambió el nombre de la columna |
| BAC-18B espera el índice parcial sobre `jti_access` | La tarea nombraba `jti_token`, que se renombró en BAC-17B |
| La integración SEC-01/SEC-02 usa el access token nuevo y la cookie rotada después del refresh, y verifica que el access anterior quede invalidado | Es lo que hace el frontend real. BAC-17B rota el access en cada refresh |

---

## 6. Pendientes y observaciones

- **Frontend:** no hubo commits. Lo más urgente sigue siendo `FIX-28` (media hora), y ahora se puede avanzar con `FRN-17C`.
- **INF-08B:** falta reemplazar los valores de ejemplo por las credenciales reales del relay (Brevo). Cuando estén, el test se autentica solo contra el SMTP.
- **Seguridad:** el secreto del token de Proxmox sigue en texto plano en `documentacion/api-proxmox.md`.

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
