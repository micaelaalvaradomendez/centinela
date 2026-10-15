# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 15/10/2026 (con los últimos commits de los submódulos en `origin/main`)  
**Alcance:** Las tareas de [documentacion/actual.md](../documentacion/actual.md) y avance de [documentacion/etapa2.md](../documentacion/etapa2.md), contrastadas tras actualizar ambos submódulos a `origin/main`. Se corrieron las suites `test/back` (Go + Docker) y `test/front` (Vitest + jsdom).

| Componente | Revisión probada (`origin/main`) | Commits nuevos incorporados |
|---|---|---|
| Backend | `44cfd16` (avanza desde `dbf7075`) | `44cfd16` (**`Telemetría`**: CPU unificada en porcentaje $0-100$ y `null` estricto con máquina detenida en WS, RRD y `/instances`), `eb14053` (docs cookies/refresh), `3b04606` (límites WS), `ad1ee86` (cierre de WS al perder acceso), `dbf7075` (**`BAC-37`**: eliminación de reservas obsoletas) |
| Frontend | `6db04d3` (en `origin/main`) | `6db04d3` / `57c739f` (reponer paneles del Dashboard con curvas, campana y alertas), `4e666a0` (**`FIX-89`**: advertencia toast en asistente estático), `757580a` (**`FIX-87`**: limpieza de datos estáticos en ficha de usuario) |

---

## 1. 🔴 Bloqueantes y fallas críticas detectadas en la integración

### 1.1 Backend: Bloqueante crítico de orden de migraciones (`FIX-81` pendiente en `centinela-back`)
El backend `44cfd16` **sigue sin arrancar contra una base de datos nueva**. Con la base limpia del stack de pruebas el contenedor aborta con código 1:

```
❌ Error fatal al iniciar la base de datos: error al clasificar usuarios eliminados:
   ERROR: relation "auditoria" does not exist (SQLSTATE 42P01)
```

- **Causa raíz:** En `backend/internal/adapters/secondary/postgres/db.go`, después de `AutoMigrate` se llama `migrarUnicidadEmail(db)` **antes** de `migrarAuditoriaParticionada(db)`. `migrarUnicidadEmail` → `clasificarEliminados` ejecuta `SELECT … FROM auditoria` (`migrar_usuarios.go`), pero la tabla `auditoria` se crea recién en `migrarAuditoriaParticionada`.
- **Detección automatizada:** Se incorporó en `test/back/fixes_acceptance_test.go` la prueba `TestFIX81_OrdenMigracionesBaseLimpia` que valida el orden de inicialización y la guarda de existencia de tabla.
- **Fix requerido:** Ejecutar `migrarUnicidadEmail(db)` **después** de `migrarAuditoriaParticionada(db)` y `asegurarParticionesFuturas(db)` en `db.go`, o añadir `if !existeTabla(db, "auditoria") { return nil }` en `clasificarEliminados`.

### 1.2 Frontend: Regresión en renderizado del Dashboard (commit `57c739f`)
El commit `57c739f` introdujo una excepción no controlada en `frontend/centinela/src/pages/Dashboard.tsx:58`:
```ts
const formatCount = (value: number | null) => value === null ? '—' : value.toLocaleString('es');
```
- **Problema:** En la tarjeta de Plantillas (línea 133), se invoca `formatCount(templates.loading ? null : templates.count)`. Si `templates.count` es `undefined` (antes de cargar o sin inicializar), la condición `value === null` es falsa y se evalúa `undefined.toLocaleString('es')`.
- **Impacto:** Arroja `TypeError: Cannot read properties of undefined (reading 'toLocaleString')`, provocando el desmontaje completo del árbol de componentes de `<Dashboard>`, haciendo fallar en cascada suites como `dashboard-node.test.tsx` (12 fallos) y `dashboard-summary.test.tsx` (2 fallos).
- **Fix requerido:** Utilizar comprobación segura: `(typeof value !== 'number' || Number.isNaN(value)) ? '—' : value.toLocaleString('es')`.

### 1.3 Frontend: Ruptura de telemetría de CPU (`2.500 %`) y glifo de ausencia (`FIX-97`)
A raíz de la unificación del Backend en `44cfd16` (CPU enviada como porcentaje $0-100$, ej. `25`):
- **Problema 1 (Bug $\times 100$):** `formatCpuUsage` en `formatInstanceTelemetry.ts` utiliza `Intl.NumberFormat('es-AR', { style: 'percent' })`, el cual multiplica internamente por 100 la entrada esperando una fracción $0..1$. Al recibir `25`, renderiza erróneamente **`"2.500 %"`** en la tabla.
- **Problema 2 (Glifo de ausencia):** Ante instancias apagadas (`null`), tanto `formatCpuUsage` como `formatRamUsage` y la fila de `Instances.tsx:144` retornan un guion simple `'-'` (`U+002D`) en vez del em-dash `'—'` (`U+2014`).
- **Impacto:** `test/front/instances-metrics.test.tsx` falla 2 de 3 pruebas (`muestra el consumo de CPU (%)` y `usa "—" cuando vienen en null`).

---

## 2. Resultados de las suites de prueba

### Frontend (`test/front`) — Vitest + jsdom

| Métrica | Valor |
|---|---:|
| Archivos evaluados | 29 (16 ✅ · 11 ❌ · 2 omitidos e2e) |
| Pruebas ejecutadas | 171 (142 ✅ · 29 ❌ · 16 omitidas) |
| Errores no controlados | 16 (por la excepción en `Dashboard.tsx:58`) |

#### Desglose de archivos que fallan y causas técnicas:

| Archivo | Fallas | Tarea de `actual.md` / `futuro-1.md` | Causa técnica detallada |
|---|---:|---|---|
| `dashboard-node.test.tsx` | 12 | `FRN-19B` | Bloqueado por excepción no controlada en `Dashboard.tsx:58` (`undefined.toLocaleString`). |
| `instances-metrics.test.tsx` | 2 | `FRN-20C` / `FIX-97` | **1)** CPU muestra `2.500 %` en vez de `25 %`. **2)** Celdas nulas muestran `'-'` en vez de `'—'`. |
| `instances-operation.test.tsx` | 3 | `FRN-16` / `FIX-86` | Mapeo de errores 403 D2: `INSTANCE_PROTECTED` e `INSTANCE_ACCESS_DENIED` se delegan a mensaje genérico. |
| `dashboard-summary.test.tsx` | 2 | `FRN-19C` / `FIX-83` | Conteo no se refresca en menos de 2 s ante evento `TASK_FINISHED` (solo polling cada 10 s). |
| `instances-ola2.test.tsx` | 1 | `FRN-20B` / `FIX-79` | Selector de tipo de recurso (Todas, VM, LXC) oculto en dropdown en vez de pestañas/botones visibles. |
| `instance-detail.test.tsx` | 2 | `FRN-21` / `FIX-95` | Ruta `/instances/:vmid` no registrada en `AppRoutes.tsx` e inventario sin enlaces hacia el detalle. |
| `fixes-pending.test.tsx` | 1 | `FIX-89` / `FIX-99` | Acceso directo por URL a `/instances/new` carece de banner visible de "Próximamente en Etapa 2". |
| `instances-modals.test.tsx` | 1 | `FRN-15` / `FIX-85` | Modal de confirmación de acción `Start` ("Iniciar") no confirma adecuadamente o ejecuta sin modal. |
| `instances-table.test.tsx` | 1 | `FRN-20A` | Timeout concurrente al interactuar con el portapapeles en el entorno de pruebas jsdom. |
| `recover-password.test.tsx` | 1 | `FRN-12` | Timeout de simulación asíncrona de envío de formulario de reseteo. |
| `auth-branding.test.tsx` | 0 | `FIX-59` / `FIX-80` | **5 de 5 en verde (100% ✅)**. |
| `audit.test.tsx` | 0 | `FRN-14` / `FIX-74` | **7 de 7 en verde (100% ✅)** (adaptado a la UI real de pestañas y dropdown). |

---

### Backend (`test/back`) — Go + Docker

| Resultado | Test / Tarea | Detalle |
|---|---|---|
| ❌ FAIL | `TestFIX81_OrdenMigracionesBaseLimpia` (**`FIX-81`**) | Falla el arranque en base limpia: `relation "auditoria" does not exist`. |
| ❌ FAIL | `TestINT01` (Paso 6 = **`BAC-24C`**) | Falla porque el backend no purga las filas de `permisos_instancia` al eliminar una VM (`DELETE`). |
| ❌ FAIL | `TestEtapa2INF10_DisponibilidadPlantillasEIsos` (**`INF-10`**) | El catálogo real de `/node/templates` aún no expone ISOs y plantillas validadas en Proxmox. |
| ⏳ PEND | `BAC-36`, `BAC-39`, `BAC-40`, `BAC-41`, `BAC-50`, `BAC-51` | Tareas de `actual.md` pendientes de desarrollo; aún no cuentan con endpoints en el router. |
| ✅ PASS | `TestEtapa1Ola3BAC30` (**`BAC-30` / `FIX-84`**) | Storages y catálogo simulado canónico en verde. |
| ✅ PASS | `TestEtapa1Ola3BAC26` (**`BAC-26`**) | Verificación de 5 acciones operativas emitiendo `TASK_FINISHED` canónico con `exitstatus`. |
| ✅ PASS | `TestEtapa2BAC34_AdaptadorMetricasRRD` (**`BAC-34`**) | Cliente RRDdata y DTO `PuntoMetricaDTO` con CPU porcentual verificado. |
| ✅ PASS | `TestEtapa2BAC35_CanalWebSocketMetricasConTicket` (**`BAC-35`**) | Handshake con ticket efímero Redis y WebSocket de telemetría en verde. |
| ✅ PASS | `TestFixesPendientes` | `FIX-54`, `FIX-55`, `FIX-71` y `FIX-88` verificados en verde. |

---

## 3. Estado de tareas de `documentacion/actual.md`

| Identificador | Área | Responsable | Estado real según pruebas | Ubicación objetivo |
|---|---|---|---|---|
| `FIX-81` | Backend | Lisandro / Tayra | ❌ Falla (bloquea base limpia) | `futuro.md` |
| `BAC-24C` | Backend | Tayra | ❌ Falla (`TestINT01` paso 6) | `actual.md` |
| `BAC-36` | Backend | Tayra | ⏳ No implementada | `actual.md` / `etapa2.md` |
| `BAC-39` | Backend | Lisandro | ⏳ No implementada | `actual.md` / `etapa2.md` |
| `BAC-40` | Backend | Tayra | ⏳ No implementada | `actual.md` / `etapa2.md` |
| `BAC-41` | Backend | Tayra | ⏳ No implementada | `actual.md` / `etapa2.md` |
| `BAC-50`..`53` | Backend | Lisandro / Tayra | ⏳ No implementadas | `actual.md` / `etapa2.md` |
| `FRN-20C` / `FIX-97` | Frontend | Cristian / Luz | ❌ Falla (CPU $2500\%$ y glifo `'-'`) | `futuro-1.md` |
| `FRN-21` / `FIX-95` | Frontend | Cristian / Belinda | ❌ Falla (ruta `/instances/:vmid` no registrada) | `futuro-1.md` |
| `FRN-19C` / `FIX-83` | Frontend | Belinda | ❌ Falla (no refresca ante `TASK_FINISHED`) | `futuro-1.md` |
| `FRN-16` / `FIX-86` | Frontend | Cristian | ❌ Falla (mapeo de errores 403 D2) | `futuro-1.md` |
| `FRN-15` / `FIX-85` | Frontend | Luz | ❌ Falla (confirmación de Start en modal) | `futuro-1.md` |
| `FRN-20B` / `FIX-79` | Frontend | Luz | ❌ Falla (selector Todas/VM/LXC oculto en menú) | `futuro-1.md` |
| `FIX-89` / `FIX-99` | Frontend | Cristian | ❌ Falla (falta banner en `/instances/new`) | `futuro-1.md` |
| `FIX-74` | Frontend | Luz | ✅ Resuelta y verde (7/7 en `audit.test.tsx`) | `terminado-1.md` |
| `FIX-87` | Frontend | Cristian | ✅ Resuelta y verde (2/2 en `fixes-pending.test.tsx`) | `terminado-1.md` |

---

## 4. Próximos pasos recomendados para los equipos

1. **Backend (`centinela-back`):**
   - **Bloqueante 1:** Resolver `FIX-81` en `backend/internal/adapters/secondary/postgres/db.go` para permitir el despliegue en bases de datos limpias.
   - **Tarea inmediata:** Implementar `BAC-24C` (limpieza atómica de permisos tras borrado de instancia).
2. **Frontend (`centinela-front`):**
   - **Bloqueante 1 (Regresión crítica):** Corregir en `Dashboard.tsx:58` el manejo de `undefined` en `formatCount` para destrabar el renderizado del Dashboard.
   - **Bloqueante 2 (`FIX-97`):** Adaptar `formatCpuUsage` para que no multiplique por 100 la CPU enviada por el backend (`44cfd16`) y retorne `'—'` ante valores `null`.
   - Registrar la ruta `/instances/:vmid` en `AppRoutes.tsx` (`FIX-95`).
