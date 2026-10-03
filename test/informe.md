# Informe de estado de tareas verificado por pruebas

**Fecha de actualización:** 02/10/2026
**Alcance:** comparación de los commits de `backend` y `frontend` con `origin/main`, y nueva ejecución de las suites de [test/back](back/README.md) y [test/front](front/package.json). El estado de tareas de las secciones siguientes corresponde a la verificación previa contra los últimos commits remotos documentados.

| Componente | Revisión probada | Commits nuevos desde la verificación anterior |
|---|---|---|
| Backend | Checkout y `origin/main` en `0167b96` | 0. No hubo commits nuevos |
| Frontend | Checkout y `origin/main` en `6fd2c7c` | 14: PR #72 corrige el rol `ADMIN`; PR #73 corrige el uso de `READ_ONLY` como rol |

**Actualización de submódulos:** backend ya estaba en `0167b96`; frontend avanzó de `4e0e7b7` a `6fd2c7c`. Ambos quedaron limpios y alineados con `origin/main`. El gitlink de frontend en el repositorio principal queda actualizado.

**Comandos ejecutados sobre esos commits:** `(cd test/back && go test -v -count=1 ./...)` y `pnpm --dir test/front test` en frontend `6fd2c7c` (segunda corrida con `--reporter=verbose` para identificar fallos).

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`, checkout `0167b96`) | 55 | **42** | **11** | 2 |
| Frontend (`test/front`, checkout `6fd2c7c`) | 113 | **86** | **17** | 10 |
| **Total** | **168** | **128** | **28** | **12** |

- En backend pasaron LOGIN-04, el recorrido front ↔ back, FIX-40 (`504 PROXMOX_TIMEOUT`), eventos SSE y el simulador. Los 10 fallos funcionales se concentran en BAC-18B, FIX-37, FIX-39 y las tareas BAC-29, BAC-22, BAC-23A y BAC-25A aún no implementadas.
- El fallo restante de backend es ambiental: INF-06A no pudo iniciar Redis porque Docker reportó el nombre `centinela-redis` ocupado por el contenedor `3832f9021f5f` (estado `Created`). No se retiró ese contenedor.
- En frontend fallaron 17 casos en 4 archivos; la salida de error detalla fallos de `instances-table.test.tsx`. La suite aprobó 86 casos y omitió los 10 de LOGIN-04, que se ejercitan desde backend. Las correcciones de permisos de PR #72/#73 ya están incorporadas.
- El backend omitió INF-08B (IP no autorizada por Brevo) y el DELETE aún no implementado.

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumplida · 🟡 parcial · ❌ no implementada.

### Fase base y FIX

| Tarea | Pruebas | Estado | Detalle |
|---|---|---|---|
| `SEC-03` Permisos reactivos | `navigation.test.tsx` | ✅ | `useAuth()`, `PermissionGate` y el menú cumplen el criterio. La regresión por `requiredRole="Admin"` fue corregida en PR #72; ya no aparece en los fallos de esta corrida. |
| `FIX-29` Complejidad de contraseña | `password-change` (3), `recover-password` (1) | ✅ | Pasan todos sus casos y los de `FRN-12` y `FIX-21` |
| `FRN-17C` Cliente de eventos | `events-client.test.tsx` (4) | ❌ | No existe `useEvents` |
| `FIX-36` Filtro "Acción" | `audit.test.tsx` (1) | ❌ | Sin cambios |
| `FIX-38` `READ_ONLY` como rol | `admin-users`, `navigation` | ✅ → `terminado.md` | PR #73 corrige el uso de `READ_ONLY` en permisos y pantallas de usuario. |
| `FIX-42` `useAuth` sin revisar el rol | `navigation` | ✅ → `terminado.md` | PR #73 corrige el tratamiento de `READ_ONLY` en `AuthContext` y la validación de autenticación. |
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

1. **Regresión de `SEC-03` en la ficha del usuario, corregida** (PR #72): se normalizó `requiredRole` a `ADMIN` en la ficha y los controles relacionados. La prueba ya no figura entre los fallos de la ejecución actual.
2. **La suite pudo llegar al Proxmox real (corregido).**
   - **Qué pasó:** con el contenedor del stub detenido, el nombre `proxmox` no resolvía en Docker. El backend lo completó con el dominio de búsqueda del host, `tail6bb3f3.ts.net` (Tailscale), y se conectó a `100.81.49.19`, el Proxmox real, con el token falso del stub.
   - **Impacto:** todas las respuestas fueron `401` y no se ejecutó ninguna acción.
   - **Corrección:** la suite ya no detiene el stub (simula la caída con un `503`), y los backends de prueba tienen `dns_search: invalid`.
3. **Helpers de permisos duplicados:** siguen existiendo `usePermissions` y `useAuth`; PR #73 corrige el tratamiento de `READ_ONLY`, pero no elimina la duplicación. Conviene unificar su responsabilidad antes de ampliar los controles de operación.
4. **Nueva ejecución del 02/10 sobre los últimos commits:** frontend `6fd2c7c` terminó con 17 fallos y backend `0167b96` con 11; FIX-40 pasó. INF-06A volvió a chocar con el contenedor existente `centinela-redis`, que se dejó intacto. Fue una ejecución por suite, no una medición de estabilidad.

---

## 5. Prueba integral LOGIN-04

**10 de 10 pasos en verde** (frontend real contra backend real).

---

## 6. Pendientes

- **Frontend:**
  - `FIX-36`;
  - `FRN-20A` sigue fallando en las pruebas de inventario;
  - `FRN-17C` con `FRN-17A`, `FRN-19A` y `FRN-20A`.
- **Backend:**
  - `BAC-18B`, `FIX-37` y `FIX-39`;
  - `BAC-29`, `BAC-22`, `BAC-23A` y `BAC-25A`.
- **Documentación (hecho el 02/10/2026):**
  - `FIX-40` pasó a `terminado.md`;
  - `FIX-29` pasó a `terminado.md`;
  - PR #72 corrige la regresión de `SEC-03`; PR #73 corrige los hallazgos `FIX-38` y `FIX-42`.
- **Infraestructura:** aplicar `docker/nginx-edge.conf` en el CT 103 (`FIX-31` e `INF-07B`).

---

## 7. Cómo reproducir

```bash
git -C backend fetch origin --prune
git -C frontend/centinela fetch origin --prune
git -C backend status --short --branch
git -C frontend/centinela status --short --branch
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)  # incluye LOGIN-04 front ↔ back y el simulador
pnpm --dir test/front test
```

El procedimiento solo actualiza referencias remotas y conserva los checkouts locales. Si se necesita probar el `origin/main` remoto en lugar del commit local, hacerlo en worktrees/copias aisladas; no usar `checkout --force` ni `reset --hard` sobre submódulos con cambios.

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
