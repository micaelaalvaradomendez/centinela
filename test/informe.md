# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 02/10/2026
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md): las de la fase base, más `FIX-39`, `FIX-40` y la Ola 1 de la Etapa 1, con las decisiones D1 a D4 de `etapa1.md`. Se revisó si cada tarea tenía pruebas que evalúen su criterio de éxito, y se agregaron o corrigieron las que faltaban o estaban mal. También se corrió la regresión de [terminado.md](../documentacion/terminado.md) y [terminado-1.md](../documentacion/terminado-1.md).

| Componente | Revisión probada | Commits nuevos desde la verificación anterior |
|---|---|---|
| Backend | `0167b96` (último commit de `main`) | 1: `0167b96` "fix: distinguiendo timeout de proxmox" (**FIX-40**) |
| Frontend | `4e0e7b7` (último commit de `main`) | 6 (PR #71: `SEC-03` con `AuthContext`, `PermissionGate` e `InstanceGate`; `FIX-29`) |

**Cómo se ejecutó:** cada submódulo se llevó al último commit de `origin/main`; ante un conflicto prevalece el remoto. El backend se corrió 3 veces, con el mismo resultado y sin contenedores residuales. El frontend se corrió 7 veces: 3 dieron el resultado de §1, y las otras 4, con la máquina sobrecargada, sumaron fallos por timeout en casos sin cambios (ver §4).

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 55 | **43** | **10** | 2 |
| Frontend (`test/front`) | 113 | **79** | **24** | 10 |
| **Total** | **168** | **122** | **34** | **12** |

- **Avance del frontend:**
  - `FIX-29` pasa.
  - `SEC-03` cumple su criterio: hook, `PermissionGate` y menú.
  - **Pero `SEC-03` rompió `FRN-07` y `FRN-18`**, que ya estaban terminadas: la ficha del usuario usa `requiredRole="Admin"` en lugar de `"ADMIN"` (§2).
- **Las 9 tareas de la Ola 1 y `FIX-40` ya tienen pruebas.** `INF-07B` pasa; las demás fallan porque todavía no están implementadas.
- **Se corrigieron pruebas que no evaluaban bien** (§3):
  - tres casos de `FRN-17A` pasaban sin verificar nada;
  - `FRN-19A` y `FRN-20A` fijaban nombres que la tarea no pide;
  - las de `BAC-22`, `BAC-23A`, `BAC-25A` y `FIX-40` buscaban texto en archivos fijos.
- **Hallazgo de seguridad en la suite, corregido:** en una corrida, el backend de pruebas llegó al Proxmox real a través del DNS de Tailscale. Solo recibió `401`; no se ejecutó ninguna acción (§4).

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumplida · 🟡 parcial · ❌ no implementada.

### Fase base y FIX

| Tarea | Pruebas | Estado | Detalle |
|---|---|---|---|
| `SEC-03` Permisos reactivos | `navigation.test.tsx` (4) | 🟡 | **Cumple su criterio:**<br>• `useAuth()` (en `AuthContext.ts`, con `AuthProvider`) expone `isAdmin()`, `isOperator()`, `hasRole` y `canAccessInstance`;<br>• `PermissionGate` funciona;<br>• el OPERATOR no ve Usuarios ni Auditoría en el menú.<br>**Regresión:** `detailsUserPage.tsx:224`, `:231` y `:256` usan `<PermissionGate requiredRole="Admin">`, así que el ADMIN tampoco ve "Roles y permisos" ni "Gestionar acceso". Fallan 4 casos de `FRN-07`/`FRN-18` en `admin-users.test.tsx` |
| `FIX-29` Complejidad de contraseña | `password-change` (3), `recover-password` (1) | ✅ | Pasan todos sus casos y los de `FRN-12` y `FIX-21` |
| `FRN-17C` Cliente de eventos | `events-client.test.tsx` (4) | ❌ | No existe `useEvents` |
| `FIX-36` Filtro "Acción" | `audit.test.tsx` (1) | ❌ | Sin cambios |
| `FIX-38` `READ_ONLY` como rol | `admin-users` (1), `navigation` (1) | ❌ | El selector "Rol" sigue ofreciendo `READ_ONLY`, y `usePermissions` lo acepta como rol |
| `FIX-42` (nuevo, en `futuro.md`) `useAuth` sin revisar el rol | `navigation` (1, **nuevo**) | ❌ | `useAuth().canAccessInstance`, de `SEC-03`, solo mira `instanciasPermitidas` |
| `BAC-18B` Índices, particiones y purga | `cierre_fase_base…` (2) | ❌ | Sin cambios |
| `FIX-37` Nivel de acceso en el perfil | `resource_access…` (1) | ❌ | Sin cambios |
| `FIX-39` Completar `BAC-21B` | `puente_etapa1…` (3) | ❌ | Faltan los campos, `shutdown`/`reboot` (404) y la auditoría de energía |
| `FIX-40` `504 PROXMOX_TIMEOUT` | `puente_etapa1…` (1, **nuevo**) | ✅ → `terminado.md` | Backend `0167b96`: `504 PROXMOX_TIMEOUT` ante un timeout; `502` si Proxmox está caído. `BAC-14`, `SEC-04` y `LOGIN-04` siguen en verde |

### Etapa 1 (Ola 1)

| Tarea | Pruebas | Estado | Detalle |
|---|---|---|---|
| `INF-07B` Nginx para SSE | `puente_etapa1…` (1) | ✅ (referencia local) | `TASK_FINISHED` por HTTPS en ~1 s. Falta aplicarlo en el CT 103 (`INT-03`) |
| `BAC-29` Contrato de la etapa | `etapa1_acceptance…` (1, **nuevo**) | ❌ | No existe `backend/docs/contrato-etapa1.md` |
| `BAC-22` Telemetría del nodo | `etapa1_acceptance…` (1, **nuevo**) | ❌ | `GET /api/node/status` responde 404 |
| `BAC-23A` Adaptador de IP | `etapa1_acceptance…` (1, **nuevo**) | ❌ | No se consultan los endpoints de IP |
| `BAC-25A` Pool de UPID | `etapa1_acceptance…` (1, **nuevo**) | ❌ | Falta `UPID_WORKERS`. No se pierden tareas: de 20 órdenes, las 20 quedan `COMPLETED` |
| `FRN-19A` Medidores del host | `dashboard-metrics.test.tsx` (3, **reescrito**) | ❌ | `features/dashboard` solo tiene `DashCard.jsx`, vacío |
| `FRN-20A` Tabla de inventario | `instances-table.test.tsx` (5, **reescrito**) | ❌ | `Instances.tsx` es una maqueta que no consulta la API |
| `FRN-17A` Consumo de eventos | `events-client.test.tsx` (4, **reescrito**) | ❌ | No existe `useEvents` |

---

## 3. Pruebas revisadas: qué faltaba y qué estaba mal

| Tarea | Problema | Corrección |
|---|---|---|
| `SEC-03` (`navigation`) | Buscaba solo `usePermissions` con booleanos y `PermissionGate` como export con nombre; además renderizaba sin el `AuthProvider` que `App.tsx` sí monta. El frontend usó `useAuth()` con funciones y un export por defecto, y la tarea lo permite | Acepta `usePermissions`, `useAuthUser` o `useAuth`, con booleanos o funciones, y `PermissionGate` con export por defecto. Nuevo `app-providers.tsx` con los providers de `src/context`. `FIX-38` evalúa `usePermissions` y el caso nuevo `FIX-42`, los demás hooks |
| `FRN-17A` (`events-client`) | 3 casos con `if (subscribe) … else expect(definido)`: pasaban sin verificar | Verifican la entrega y el descarte, la deduplicación, una sola conexión con dos consumidores, y la limpieza de `useWebSocket.js` y `notifications.ts` |
| `FRN-20A` (`instances-table`) | Tipo `qemu` en el mock (la API devuelve `vm`), un `if (copyButton)`, `AuthProvider` importado de forma fija y un nombre accesible exigido para la tabla. No verificaba los colores, la columna Tipo, el `Bearer` ni la botonera | Reescrita contra el contrato de `BAC-14` y el entregable |
| `FRN-19A` (`dashboard-metrics`) | Exigía `CpuGauge`, `usagePercent` y clases como `.text-warning` | Detecta los medidores por lo que renderizan; compara el aspecto en 20, 69 y 70 % (D3); busca las pruebas del frontend con 69 y 70 |
| `BAC-22`, `BAC-23A`, `BAC-25A`, `BAC-29`, `FIX-40` (`etapa1_core…`, eliminado) | Buscaban texto en archivos fijos (`"lxc"`, `"workers"`, una expresión regular); `BAC-22` pedía `instancesSummary` (de `BAC-22B`) y no probaba ni la caché ni `stale` | `etapa1_acceptance_test.go` y el caso `FIX-40` de `puente…`, por comportamiento (detalle en [`back/RESULTADOS.md`](back/RESULTADOS.md)) |
| `BAC-21B` (auditoría) | Buscaba un UPID fijo del stub | Usa el `upid` de la respuesta |

**Contrapruebas del frontend:** contra una copia temporal con implementaciones mínimas correctas pasan 16 de 16. Con 8 defectos introducidos (umbral corrido, sin deduplicar, una conexión por consumidor, badges sin color, entre otros), cada vez falla exactamente la prueba que corresponde. Las del backend no se pudieron contrastar contra una implementación, porque las tareas no existen todavía.

---

## 4. Hallazgos

1. **Regresión de `SEC-03` en la ficha del usuario** (frontend `1156eaa`): `requiredRole="Admin"` en `detailsUserPage.tsx:224`, `:231` y `:256`. El rol del sistema es `ADMIN`, y `hasRole` compara exacto. Se puede registrar como FIX o resolver dentro de `SEC-03`, que sigue en `actual.md`.
2. **La suite pudo llegar al Proxmox real (corregido).**
   - **Qué pasó:** con el contenedor del stub detenido, el nombre `proxmox` no resolvía en Docker. El backend lo completó con el dominio de búsqueda del host, `tail6bb3f3.ts.net` (Tailscale), y se conectó a `100.81.49.19`, el Proxmox real, con el token falso del stub.
   - **Impacto:** todas las respuestas fueron `401` y no se ejecutó ninguna acción.
   - **Corrección:** la suite ya no detiene el stub (simula la caída con un `503`), y los backends de prueba tienen `dns_search: invalid`.
3. **Duplicación en el frontend:** hay dos helpers de permisos, `usePermissions` (con `canOperateInstance`) y `useAuth` (con `hasRole` y su propio `canAccessInstance`), y deciden distinto. `FIX-38` corrige `usePermissions` y `FIX-42`, `useAuth`. Conviene que el equipo elija uno antes de `FRN-15`/`FRN-16`.
4. **Estabilidad del frontend:** con la máquina sin memoria disponible y carga de 15 a 41 (procesos ajenos: `k3s server` y un `docker buildx` de otro proyecto), algunas corridas sumaron de 2 a 13 timeouts en casos sin cambios. El resultado de §1 se repitió en 3 corridas.

---

## 5. Prueba integral LOGIN-04

**10 de 10 pasos en verde** (frontend real contra backend real).

---

## 6. Pendientes

- **Frontend:**
  - `FIX-41`, la regresión de `SEC-03` (`requiredRole="Admin"`);
  - `FIX-38` y `FIX-42` (el mismo defecto en `usePermissions` y en `useAuth`);
  - `FIX-36`;
  - `FRN-17C` con `FRN-17A`, `FRN-19A` y `FRN-20A`.
- **Backend:**
  - `BAC-18B`, `FIX-37` y `FIX-39`;
  - `BAC-29`, `BAC-22`, `BAC-23A` y `BAC-25A`.
- **Documentación (hecho el 02/10/2026):**
  - `FIX-40` pasó a `terminado.md`;
  - `FIX-29` pasó a `terminado.md`;
  - `SEC-03` pasó a `terminado.md` como implementada con problema, y su regresión es `FIX-41`, en `futuro.md`.
- **Infraestructura:** aplicar `docker/nginx-edge.conf` en el CT 103 (`FIX-31` e `INF-07B`).

---

## 7. Cómo reproducir

```bash
for s in backend frontend; do
  git -C $s fetch origin --prune && git -C $s checkout -B main origin/main --force && git -C $s reset --hard origin/main
done
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)   # incluye LOGIN-04 front ↔ back y el simulador
pnpm --dir test/front test                                        # con la máquina cargada: pnpm --dir test/front vitest run --no-file-parallelism
```

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
