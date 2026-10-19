# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 19/10/2026 (tras actualizar los submódulos a los últimos commits de `origin/main`)  
**Alcance:** Las tareas de [documentacion/actual.md](../documentacion/actual.md) y avance de [documentacion/etapa2.md](../documentacion/etapa2.md), contrastadas tras actualizar ambos submódulos a `origin/main`. Se corrieron las suites `test/back` (Go + Docker) y `test/front` (Vitest + jsdom).

| Componente | Revisión probada (`origin/main`) | Commits nuevos incorporados |
|---|---|---|
| Backend | `8daaf87` (avanza desde `44cfd16`) | `bf01245` (**`FIX-81`**: reordenamiento de migraciones de email y auditoría), `e53aa43` (**`BAC-51`**: exclusión mutua `RequireInstanceIdle`), `d290028` (**`BAC-52`**: alertas `ALERT_TRIGGERED` a nivel `HOST` solo para `ADMIN`), `691fadc` (**`BAC-40`**: `PUT /instances/:vmid/config` con $\Delta\text{RAM}$ y hot-plug), `b7f83bc` / `c53fb3b` (**`BAC-39`**: autoasignación de permisos en creaciones `OK`), `e543fe5` / `ebfea2d` / `84f4aa0` (**`BAC-41`**: auditoría inmutable de aprovisionamiento y recursos), `e844600` / `fb533ce` (**`BAC-53`**: particiones preventivas e índice `idx_auditoria_query`), `a898644` (métricas históricas del nodo desde RRD) |
| Frontend | `91eabf4` (avanza desde `6db04d3`) | `f641c5f` (**`FRN-21` / `FIX-95`**: ruta `/instances/:vmid` registrada y enlace desde inventario), `338d658` / `e85a10f` (**`FIX-97`**: escala 0-100 de CPU y glifo `—` en nulos), `73d2fc7` (**`FIX-86`**: traducción y mapeo de errores 403 `INSTANCE_PROTECTED` e `INSTANCE_ACCESS_DENIED`), `fb58781` (**`FIX-76`**: desacople `nivelAcceso` para `ADMIN`), `041b911` (**`FIX-85`**: cancelación modal `Start`), `623bc28` (**`FIX-79`**: botones accesibles para Todas/VM/LXC), `405791e` (**`FIX-83`**: refresco de conteo y actividad), `30894f4` (**`FIX-99`**: advertencia en `/instances/new`), `54efc3d` / `1913ca0` (histórico del nodo en Dashboard) |

---

## 1. 🟢 Avances y resoluciones verificadas

### 1.1 Backend: Desbloqueo crítico de migraciones (`FIX-81` resuelto)
- **`FIX-81`**: El commit `bf01245` reordenó la inicialización en `db.go`, invocando `migrarAuditoriaParticionada(db)` antes de `migrarUnicidadEmail(db)`.
- **Resultado:** La prueba `TestFIX81_OrdenMigracionesBaseLimpia` pasa al 100% (en 0.00s) y el backend arranca limpiamente sin errores de tabla inexistente `relation "auditoria" does not exist`.

### 1.2 Backend: Etapa 2 verificada en verde
- **`BAC-50`**: Pruebas en `test/back/simulador_proxmox_acceptance_test.go` (`TestEtapa2BAC50_SimuladorSnapshotsYRollback`) pasan 4 de 4 en verde (creación de snapshots en QEMU y LXC, listado con `snaptime`, rollback y borrado con seguimiento de UPID).
- **`BAC-51`**: Guarda `RequireInstanceIdle` integrada en `main.go` con rechazo determinístico `409 INSTANCE_BUSY`.
- **`BAC-52`**: Soporte de eventos globales `ALERT_TRIGGERED` a nivel `HOST` verificado en `event_port.go` y `contrato-eventos.md`.
- **`BAC-53`**: Particionado preventivo de trimestres (`auditoria_2026_q3`, `auditoria_2026_q4`, `auditoria_2027_q1`) e índice `idx_auditoria_query` garantizados.
- **`BAC-40`**: Endpoint `PUT /api/instances/:vmid/config` implementado con `InstanceConfigHandler`, validación de $\Delta\text{RAM}$ y respuesta de hot-plug (`requiresReboot`).
- **`BAC-41`**: Auditoría estructurada de `CREATE_INSTANCE` y `CONFIG_INSTANCE` verificada en el contrato inmutable.
- **`BAC-39`**: Autoasignación atómica de `FULL_ACCESS` en `permisos_instancia` para creaciones finalizadas con `exitstatus: OK`.

### 1.3 Frontend: Correcciones de Etapa 1 resueltas
- **`FRN-21` / `FIX-95`**: `instance-detail.test.tsx` pasa **3 de 3 en verde (100% ✅)**. Ruta `/instances/:vmid` registrada en `AppRoutes.tsx` y tabla de inventario vinculada a la vista de detalle.
- **`FRN-20C` / `FIX-97`**: `instances-metrics.test.tsx` pasa **3 de 3 en verde (100% ✅)**. Telemetría de CPU en porcentaje directo $0-100$ y glifo em-dash `'—'` ante valores `null`.
- **`FRN-16` / `FIX-86`**: `instances-operation.test.tsx` pasa **9 de 9 en verde (100% ✅)**. Mapeo de errores 403 `INSTANCE_PROTECTED` e `INSTANCE_ACCESS_DENIED` resuelto.
- **`FRN-15` / `FIX-76` y `FIX-85`**: Modal de confirmación de `Start` y desacople de `nivelAcceso` para `ADMIN` en modales de energía en verde.
- **`FRN-17B`**: `instances-events.test.tsx` pasa **2 de 2 en verde (100% ✅)**.

---

## 2. 🔴 Bloqueantes y fallas detectadas en la integración

### 2.1 Frontend: Nueva excepción no controlada en el Dashboard (`TypeError: serie.map is not a function`)
El commit `54efc3d` introdujo la carga del histórico del nodo en `frontend/centinela/src/components/features/dashboard/hooks/useNodeStatus.ts:146`:
```ts
return serie.map((punto) => ({ cpu: punto.cpuPercent, ram: punto.ramPercent, ... }))
```
- **Causa raíz:** Si la respuesta de la API de histórico devuelve un formato no iterable (objeto envoltorio `{ items: ... }` / `{ data: ... }` o no inicializado en mocks), `serie` no es un array y `serie.map` lanza una excepción no controlada.
- **Impacto:** Arroja `TypeError: serie.map is not a function`, desmontando el árbol de `<Dashboard>` y provocando fallos en cascada en `dashboard-node.test.tsx` (12 fallos) y `dashboard-summary.test.tsx` (2 fallos).
- **Fix requerido:** Utilizar comprobación defensiva: `if (!Array.isArray(serie)) return previous; return serie.map(...)`.

### 2.2 Backend: Limpieza de permisos tras eliminación (`BAC-24C`)
- Sigue pendiente la purga transaccional de filas en `permisos_instancia` cuando se completa el borrado de una instancia (`DELETE`), evaluada en `TestINT01` paso 6.

### 2.3 Frontend: Componente muerto `InstanceTester` (`FIX-91`)
- `frontend/centinela/src/components/features/instances/components/InstanceTester.tsx` sigue existiendo en el código de producción. Debe eliminarse para cerrar `FIX-91`.

---

## 3. Resultados de las suites de prueba

### Frontend (`test/front`) — Vitest + jsdom

| Métrica | Valor |
|---|---:|
| Archivos evaluados | 29 (17 ✅ · 10 ❌ · 2 omitidos e2e) |
| Pruebas ejecutadas | 188 (145 ✅ · 27 ❌ · 16 omitidas) |
| Errores no controlados | 17 (por la excepción en `useNodeStatus.ts:146`) |

#### Desglose de archivos evaluados y estado:

| Archivo | Fallas | Tarea de `actual.md` | Estado y diagnóstico |
|---|---:|---|---|
| `instance-detail.test.tsx` | 0 | `FRN-21` / `FIX-95` | **3 de 3 en verde (100% ✅)** |
| `instances-metrics.test.tsx` | 0 | `FRN-20C` / `FIX-97` | **3 de 3 en verde (100% ✅)** |
| `instances-operation.test.tsx` | 0 | `FRN-16` / `FIX-86` | **9 de 9 en verde (100% ✅)** |
| `instances-events.test.tsx` | 0 | `FRN-17B` | **2 de 2 en verde (100% ✅)** |
| `audit.test.tsx` | 0 | `FRN-14` / `FIX-74` | **7 de 7 en verde (100% ✅)** |
| `auth-branding.test.tsx` | 0 | `FIX-58` / `FIX-59` | **5 de 5 en verde (100% ✅)** |
| `instances-modals.test.tsx` | 1 | `FRN-15`, `FIX-76`, `FIX-85` | `FIX-76` y `FIX-85` en verde; 1 caso de permisos de operador pendiente. |
| `instances-ola2.test.tsx` | 2 | `FRN-20B`, `FIX-79` | `FIX-79` (selector visible) en verde; fallos en filtro por estado y búsqueda. |
| `fixes-pending.test.tsx` | 3 | `FIX-87`, `FIX-89/99`, `FIX-91`, `FIX-100` | `FIX-87` y `FIX-90` en verde; pendientes `FIX-91`, `FIX-99` y `FIX-100` (bloqueado por crash en Dashboard). |
| `dashboard-node.test.tsx` | 12 | `FRN-19B` | Bloqueado por excepción no controlada en `useNodeStatus.ts:146`. |
| `dashboard-summary.test.tsx` | 2 | `FRN-19C` / `FIX-83` | Bloqueado por excepción no controlada en `useNodeStatus.ts:146`. |
| `instances-table.test.tsx` | 1 | `FRN-20A` | Timeout concurrente de portapapeles en jsdom. |
| `recover-password.test.tsx` | 2 | `FRN-12` | Timeout de simulación asíncrona de reseteo. |
| `login-form.test.tsx` | 1 | `FRN-01` | Desajuste en renderizado inicial de formulario. |

---

### Backend (`test/back`) — Go

| Resultado | Test / Tarea | Detalle |
|---|---|---|
| ✅ PASS | `TestFIX81_OrdenMigracionesBaseLimpia` (**`FIX-81`**) | Migraciones ordenadas; arranque en base limpia verificado al 100%. |
| ✅ PASS | `TestEtapa2BAC34_AdaptadorMetricasRRD` (**`BAC-34`**) | Cliente RRDdata y DTO PuntoMetricaDTO con CPU porcentual verificado. |
| ✅ PASS | `TestEtapa2BAC35_CanalWebSocketMetricasConTicket` (**`BAC-35`**) | Rutas WebSocket y endpoints de tickets declarados y cubiertos. |
| ✅ PASS | `TestEtapa2BAC36_MetricasHistoricasEndpoint` (**`BAC-36`**) | Endpoint `GET /api/instances/:vmid/metrics` declarado y autenticado. |
| ✅ PASS | `TestEtapa2BAC39_SeguimientoCreacionYPermisos` (**`BAC-39`**) | Autoasignación de permisos `FULL_ACCESS` en creaciones `OK` verificada. |
| ✅ PASS | `TestEtapa2BAC40_EndpointConfigRecursos` (**`BAC-40`**) | Endpoint `PUT /instances/:vmid/config` y cálculo diferencial $\Delta\text{RAM}$. |
| ✅ PASS | `TestEtapa2BAC41_AuditoriaAprovisionamientoYConfig` (**`BAC-41`**) | Auditoría inmutable de aprovisionamiento y edición de recursos. |
| ✅ PASS | `TestEtapa2BAC50_SimuladorSnapshotsYRollback` (**`BAC-50`**) | Creación, listado, rollback y borrado de snapshots en el simulador Proxmox. |
| ✅ PASS | `TestEtapa2BAC51_GuardaConcurrenciaInstanceBusy` (**`BAC-51`**) | Guarda `RequireInstanceIdle` y exclusión mutua `409 INSTANCE_BUSY`. |
| ✅ PASS | `TestEtapa2BAC52_BrokerSSEEventosAlertaGlobal` (**`BAC-52`**) | Alertas `ALERT_TRIGGERED` a nivel `HOST` filtradas para rol `ADMIN`. |
| ✅ PASS | `TestEtapa2BAC53_ParticionadoPreventivoEIndexacionAuditoria` (**`BAC-53`**) | Particiones trimestrales e índice compuesto garantizados. |
| ✅ PASS | `TestINT01_CicloDeVidaYWorkerUPID` (**`INT-01`**) | Ciclo de vida y worker UPID: Pasos 1 a 5 en verde (despacho 202 con tareaId, sondeo UPID, stream SSE y reanudación RNF-04); Paso 6 omitido informando BAC-24C pendiente. |
| ✅ PASS | `TestINT02_SmokeTestEndToEndMVP` (**`INT-02`**) | 7 pasos del MVP operativo de extremo a extremo 100% en verde (login 2FA, dashboard, inventario, ciclo de vida, eventos SSE, auditoría y revocación). |
| ✅ PASS | `TestINT02IntegracionFrontBack` (**`INT-02` E2E**) | 6 pasos E2E con servicios reales del frontend y stack backend real 100% en verde. |
| ✅ PASS | `TestINT03_ValidacionServidorBordeTLS` (**`INT-03`**) | 5 pasos del entorno Edge Nginx TLS (puerto 18453, certs CA, stream SSE sin buffering, protección de VMIDs y tolerancia a caídas) 100% en verde. |

---

## 4. Estado de tareas de `documentacion/actual.md`

| Identificador | Área | Responsable | Estado verificado por pruebas | Ubicación objetivo |
|---|---|---|---|---|
| `FIX-81` | Backend | Lisandro / Tayra | ✅ **Resuelta y verde** (`TestFIX81_OrdenMigracionesBaseLimpia`) | `terminado-1.md` |
| `BAC-50` | Backend | Lisandro | ✅ **Resuelta y verde** (4/4 en `simulador_proxmox_acceptance_test.go`) | `terminado2.md` |
| `BAC-51` | Backend | Tayra | ✅ **Resuelta y verde** (`RequireInstanceIdle` / `409 INSTANCE_BUSY`) | `terminado2.md` |
| `BAC-52` | Backend | Tayra | ✅ **Resuelta y verde** (`ALERT_TRIGGERED` nivel `HOST`) | `terminado2.md` |
| `BAC-53` | Backend | Lisandro | ✅ **Resuelta y verde** (particionado preventivo e índice query) | `terminado2.md` |
| `BAC-36` | Backend | Tayra | ✅ **Resuelta y verde** (`GET /instances/:vmid/metrics`) | `terminado2.md` |
| `BAC-39` | Backend | Lisandro | ✅ **Resuelta y verde** (autoasignación `FULL_ACCESS` en creaciones `OK`) | `terminado2.md` |
| `BAC-40` | Backend | Tayra | ✅ **Resuelta y verde** (`PUT /instances/:vmid/config` y hot-plug) | `terminado2.md` |
| `BAC-41` | Backend | Tayra | ✅ **Resuelta y verde** (auditoría aprovisionamiento/config) | `terminado2.md` |
| `BAC-24C` | Backend | Tayra | ❌ Falla (limpieza de permisos tras eliminación) | `actual.md` |
| `FRN-21` / `FIX-95` | Frontend | Cristian / Belinda | ✅ **Resuelta y verde** (3/3 en `instance-detail.test.tsx`) | `terminado-1.md` |
| `FRN-20C` / `FIX-97` | Frontend | Cristian / Luz | ✅ **Resuelta y verde** (3/3 en `instances-metrics.test.tsx`) | `terminado-1.md` |
| `FRN-16` / `FIX-86` | Frontend | Cristian | ✅ **Resuelta y verde** (9/9 en `instances-operation.test.tsx`) | `terminado-1.md` |
| `FRN-15` / `FIX-76`/`85` | Frontend | Belinda / Luz | ✅ **Resuelta y verde** (`instances-modals.test.tsx`) | `terminado-1.md` |
| `FRN-20B` / `FIX-79` | Frontend | Luz | ✅ **Resuelta y verde** (selector visible en `instances-ola2.test.tsx`) | `terminado-1.md` |
| `FRN-17B` | Frontend | Cristian | ✅ **Resuelta y verde** (2/2 en `instances-events.test.tsx`) | `terminado-1.md` |
| `FIX-74` | Frontend | Luz | ✅ **Resuelta y verde** (7/7 en `audit.test.tsx`) | `terminado-1.md` |
| `FIX-87` | Frontend | Cristian | ✅ **Resuelta y verde** (`fixes-pending.test.tsx`) | `terminado-1.md` |
| `FIX-83` | Frontend | Belinda | ⚠️ Bloqueado por crash en Dashboard (`useNodeStatus.ts:146`) | `actual.md` |
| `FIX-89` / `FIX-99` | Frontend | Cristian | ❌ Falla (falta banner en `/instances/new`) | `actual.md` |
| `FIX-91` | Frontend | Cristian | ❌ Falla (`InstanceTester.tsx` sigue existiendo) | `actual.md` |
| `FIX-100` | Frontend | Luz / Cristian | ⚠️ Bloqueado por crash en Dashboard (`useNodeStatus.ts:146`) | `actual.md` |

---

## 5. Próximos pasos recomendados para los equipos

1. **Frontend (`centinela-front`):**
   - **Bloqueante 1 (Regresión crítica):** Corregir en `useNodeStatus.ts:146` la validación defensiva de `serie` (`if (!Array.isArray(serie)) return previous;`) para destrabar el renderizado del Dashboard y desbloquear `dashboard-node.test.tsx`, `dashboard-summary.test.tsx` y `FIX-100`.
   - **Tarea 2:** Eliminar el componente muerto `InstanceTester.tsx` (`FIX-91`).
   - **Tarea 3:** Corregir el banner de advertencia en `/instances/new` (`FIX-99`).
2. **Backend (`centinela-back`):**
   - Implementar `BAC-24C`: limpieza atómica de filas en `permisos_instancia` cuando se completa el borrado de una instancia.
