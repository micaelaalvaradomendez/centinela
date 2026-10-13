# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 13/10/2026 (con los últimos commits de los submódulos en `origin/main`)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), contrastadas tras actualizar ambos submódulos a `origin/main`. Se corrieron las suites `test/back` (Go + Docker) y `test/front` (Vitest + jsdom).

| Componente | Revisión probada (`origin/main`) | Commits nuevos incorporados |
|---|---|---|
| Backend | `6afc017` | `6afc017`/`58388fc` (Nginx WS — `INF-09`), `dc018a2` (auditoría RF-08 + Swagger), `d76eb35` (persistencia de tareas), `a9b9d0f`/`c024ba3` (**separar eliminación lógica de suspensión — `FIX-54`/`FIX-55`**), `c7fd7a1` (**Redis arranque en frío — `FIX-71`**), `ec84122` (**firma de `NewSeguimientoTareas` — `FIX-72`**), `2314d38`/`6bd1cf3` (`INSTANCE_INVALID_STATE` en Swagger) |
| Frontend | `7fcc0bd` | `7fcc0bd`/`ef485be` (filtros reactivos — `FRN-20B`), `614c912` (**máquina de estados, bloqueo y mapeo de errores — `FRN-15`/`FRN-16`**), `7e4d3af` (modal de eliminar), `d0a1c52` y cadena (**distinguir suspensión y eliminación — `FIX-56`**) |

---

## 1. 🔴 Bloqueante crítico del backend (corregir en el submódulo)

El backend `6afc017` **no arranca contra una base de datos nueva**. Con la base limpia del stack de tests el contenedor sale con código 1:

```
❌ Error fatal al iniciar la base de datos: error al clasificar usuarios eliminados:
   ERROR: relation "auditoria" does not exist (SQLSTATE 42P01)
```

**Causa raíz (orden de migraciones).** En `backend/internal/adapters/secondary/postgres/db.go`, después de `AutoMigrate` se llama `migrarUnicidadEmail(db)` **antes** de `migrarAuditoriaParticionada(db)`:

- `migrarUnicidadEmail` → `clasificarEliminados` ejecuta `SELECT … FROM auditoria` (`migrar_usuarios.go`).
- La tabla `auditoria` se crea recién en `migrarAuditoriaParticionada` (más abajo, mismo archivo).

En el servidor de pruebas no se detecta porque `auditoria` ya existía; en una base nueva (deploy limpio o stack de tests) es fatal.

**Fix sugerido (mover un bloque):** ejecutar `migrarUnicidadEmail(db)` **después** de `migrarAuditoriaParticionada(db)` y `asegurarParticionesFuturas(db)`.

> **Verificación local:** con ese reordenamiento aplicado de forma temporal, el stack levanta **completo y saludable** (`docker compose up --wait` → exit 0) y toda la suite del backend corre. El cambio se revirtió para no dejar el submódulo modificado: **queda pendiente de fix en `centinela-back`.**
>
> Los resultados del backend de la §3 se obtuvieron **con ese reordenamiento aplicado temporalmente** (sin él, ninguna prueba de integración puede correr). El frontend no depende de esto.

---

## 2. Resultados de las suites

### Frontend (`7fcc0bd`) — `test/front`

| Métrica | Valor |
|---|---:|
| Archivos | 27 (17 ✅ · 8 ❌ · 2 omitidos) |
| Pruebas | 176 |
| Pasan | **144** |
| Fallan | **16** |
| Omitidas | 16 (`login04-e2e`, `int02-e2e`: se ejecutan desde `test/back`) |

| Archivo | Fallan | Tarea |
|---|---:|---|
| `instances-operation.test.tsx` | 5 | `FRN-16` (mejoró: 8 → 5) |
| `admin-users-deletion.test.tsx` | 4 | `FIX-56` (mejoró: 5 → 4) |
| `instances-modals.test.tsx` | 1 | `FRN-15` |
| `instances-ola2.test.tsx` | 1 | `FRN-20B` (mejoró: 2 → 1) |
| `audit.test.tsx` | 2 | `FRN-14` |
| `auth-branding.test.tsx` | 1 | `FIX-59` |
| `dashboard-summary.test.tsx` | 1 | `FRN-19C` (refresco por evento, pendiente) |
| `instance-detail.test.tsx` | 1 | `FRN-21` (ruta `/instances/:vmid`, pendiente) |

### Backend (`6afc017`) — `test/back` (con el reordenamiento de §1)

| Resultado | Test |
|---|---|
| ✅ PASS | `TestCierreFaseBase`, `TestEtapa1` (**`BAC-25A`/`FIX-72` arreglado**), `TestEtapa1Ola2`, `TestFixesPendientes` (**`FIX-54`/`FIX-55`/`FIX-71` arreglados**), `TestEtapa1Ola3BAC26`, `TestFIX71RedisArranqueEnFrio`, `TestINF09WebSocketsNginxEdge`, `TestPuenteEtapa1`, `TestINT02` (7/7), `TestINT02IntegracionFrontBack` (6/6), `TestINT03` (5/5), `TestHitoLOGIN04…`, `TestLOGIN04IntegracionFrontBack` (10/10), `TestHitoRecuperacion…`, `TestHitoControlDeAcceso…`, `TestHitoSeguridadDeSesiones…`, `TestEtapa1SimuladorProxmox` |
| ❌ FAIL | `TestTareasBackendEIntegracion` (subtest `BAC-05/06` — **test desactualizado, ver §4**), `TestINT01` (Paso 6 = `BAC-24C`, pendiente), `TestEtapa1Ola3BAC30` (solo `storage`), `BAC-31`, `BAC-32`, `BAC-33`, `INF-07C` (pendientes de Etapa 2) |

---

## 3. Tareas resueltas por los commits nuevos (verificadas)

- **`FIX-72`** (`BAC-25A`): las pruebas unitarias de `internal/core/services` compilan de nuevo → `TestEtapa1` **pasa**.
- **`FIX-54` / `FIX-55`**: `eliminado_en`, índice parcial del correo y alineación de endpoints → `TestFixesPendientes` **pasa** (`FIX-55` ya no se omite).
- **`FIX-71`**: el backend tolera el arranque en frío de Redis → **verde**.
- **`INF-09`**: Nginx con `Upgrade`/`Connection` y timeout largo → **verde**.
- **Frontend `FIX-56`** (`admin-users-deletion` 5 → 4), **`FRN-15`/`FRN-16`** (`instances-operation` 8 → 5) y **`FRN-20B`** (`instances-ola2` 2 → 1): avances parciales.

## 4. Cambios a la suite de este repositorio en esta corrida

- **`test/back/backend_acceptance_test.go` (`BAC-05/06`):** la prueba esperaba `409` al reutilizar el correo de un usuario eliminado. Con **`FIX-54` ya implementado** el correo queda libre (`201`) y sólo el **nombre de usuario** sigue reservado. Se alineó la aserción al contrato nuevo.
- **`test/back`**: se agregó el servicio **`proxmox-sim`** (`proxmox-sim.Dockerfile` + `compose.yaml`) para correr el **simulador real** y la suite `etapa1_ola3_acceptance_test.go` (`BAC-26/30/31/32/33`, `FIX-71`, `INF-07C`, `INF-09`).
- **`test/front`**: se agregaron `instances-metrics` (`FRN-20C`), `instances-events` (`FRN-17B`), `dashboard-summary` (`FRN-19C`) e `instance-detail` (`FRN-21`).
- **`test/front/admin-users.test.tsx` (corregido, 23/23 ✅):** los 4 fallos no eran del producto sino de la suite. (a) Los mocks usaban `vi.fn().mockResolvedValue(jsonResponse(...))`, que devuelve **la misma `Response`** en cada llamada; con el nuevo `useDeletedUsers` (pestaña "Eliminados") la vista hace `GET /admin/audit` y consume el cuerpo antes que la lista → `apiClient` devolvía "El servidor devolvió una respuesta inesperada." Se cambió a `mockImplementation(() => Promise.resolve(jsonResponse(...)))`. (b) Se usaba el código obsoleto `USER_CONFLICT`: el backend `6afc017` emite `USER_EMAIL_ALREADY_EXISTS` / `USER_USERNAME_ALREADY_EXISTS`. Se alinearon los mocks y las aserciones al contrato vigente y al rediseño de `FIX-56`.
- **Correcciones de la propia suite** (no del producto): acciones válidas por estado de la 101, columna `vmid_proxmox`, contrato paginado de `/admin/audit`, marca `/tmp/tareas-en-curso` para el timeout y contraseña ≤12 en el e2e.

---

## 5. Próximos pasos

1. **Backend (crítico):** aplicar el reordenamiento de migraciones de §1 en `centinela-back` para que el backend arranque en una base nueva.
2. **Frontend:** cerrar `FRN-15`/`FRN-16` (`instances-operation`, 5 fallos), `FIX-56` (`admin-users-deletion`, 4), `FRN-14` (`audit`, 2), `FIX-59` (`auth-branding`, 1) y `FRN-20B` (`instances-ola2`, 1).
3. **Etapa 2 (backend):** `BAC-30` (storage/contenido), `BAC-31` (`/node/storages`, `/node/templates`), `BAC-32` (worker `create` + autoasignación), `BAC-33` (`allocatable`) e `INF-07C`; y `BAC-24C` (purga de permisos al eliminar).
4. **Frontend (Etapa 1):** `FRN-19C` (refresco del conteo por `TASK_FINISHED`) y `FRN-21` (ruta `/instances/:vmid`).

