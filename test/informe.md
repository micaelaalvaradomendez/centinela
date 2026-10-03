# Informe de estado de tareas verificado por pruebas

**Fecha de actualización:** 03/10/2026
**Alcance:** sincronización de los submódulos `backend` (`4e204f1`) y `frontend` (`070e96b`) con `origin/main` y ejecución completa de las suites de [test/back](back/README.md) y [test/front](front/package.json).

| Componente | Revisión probada | Commits nuevos desde la verificación anterior |
|---|---|---|
| Backend | `4e204f1` (`origin/main`) | 5 commits nuevos de Tayra (`96106a6`, `b8d2631`, `573f0c3`, `42e94ca`, `4e204f1`): índices parciales y particionamiento de auditoría (`BAC-18B`), DTO consolidado y acciones de energía auditadas (`FIX-39`), y permisos en el perfil (`FIX-37`). |
| Frontend | `070e96b` (`origin/main`) | 12 commits nuevos (PR #74 y PR #75 de Cristian: niveles de acceso por instancia y maqueta viva de `Instances.tsx`; PR #76 y PR #77 de Belinda y Luz: restauración del filtro "Acción" en Auditoría, **`FIX-36`**). |

**Actualización de submódulos:** ambos submódulos quedaron alineados con `origin/main` (`backend: 4e204f1`, `frontend: 070e96b`).

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`, `4e204f1`) | 61 | **53** | **7** | 1 |
| Frontend (`test/front`, `070e96b`) | 113 | **87** | **16** | 10 |
| **Total** | **174** | **140** | **23** | **11** |

- **Avances clave en Backend:**
  - **`FIX-37` está 100% completo y verificado:** `GET /account/profile` ahora expone `permisos: [{ vmid, nivelAcceso }]`.
  - **`FIX-39` está 100% completo y verificado:** `GET /api/instances` suma los campos nuevos de la etapa, `shutdown` y `reboot` quedan integrados con `FULL_ACCESS`, y las órdenes se auditan en `PENDING` con el `upid`.
  - **`BAC-18B` implementado con observaciones:** particionamiento trimestral de `auditoria` en Postgres completado y `PurgaWorker` activo. Falta incorporar `jti_access` en el índice parcial (se corrige con `FIX-44`).
- **Avances clave en Frontend:**
  - **`FIX-36` está 100% completo y verificado:** `Auditoria.tsx` restauró el filtro reactivo y exportación CSV de "Acción" (7/7 verde).
  - **`FRN-20A` implementado con observaciones:** se reemplazó la maqueta estática por integración reactiva contra `GET /instances` y modales de confirmación en `Instances.tsx`, pero fallan 5 casos por discrepancias de maquetado en IP, badges y nombres (se corrige con `FIX-43`).
- **Pendientes no implementados:**
  - Frontend: `FRN-17C` / `FRN-17A` (cliente/eventos) y `FRN-19A` (gauges de dashboard).
  - Backend: `BAC-29` (contrato), `BAC-22` (telemetría), `BAC-23A` (IPs) y `BAC-25A` (worker pool).
- **Omitidos:** 10 de frontend corresponden a LOGIN-04 (corre integrado contra backend real); 1 de backend corresponde a `INF-08B` (IP pública rechazada por Brevo).

---

## 2. Estado de las tareas y correspondencia documental

**Leyenda:** ✅ cumplida · 🟡 implementada con observaciones/problema · ❌ no implementada.

### Fase Base y FIXes

| Tarea | Pruebas | Estado | Destino Documental |
|---|---|---|---|
| `SEC-03` Permisos reactivos | `navigation.test.tsx` (9) | ✅ | [`documentacion/terminado.md`](../documentacion/terminado.md) |
| `FIX-29` Complejidad de contraseña | `password-change` (3), `recover-password` (1) | ✅ | [`documentacion/terminado.md`](../documentacion/terminado.md) |
| `FIX-36` Filtro "Acción" en Auditoría | `audit.test.tsx` (7) | ✅ | [`documentacion/terminado.md`](../documentacion/terminado.md) |
| `FIX-38` `READ_ONLY` como rol | `admin-users` (20), `navigation` (9) | ✅ | [`documentacion/terminado.md`](../documentacion/terminado.md) |
| `FIX-42` `useAuth` sin revisar rol | `navigation` (9) | ✅ | [`documentacion/terminado.md`](../documentacion/terminado.md) |
| `FIX-40` `504 PROXMOX_TIMEOUT` | `puente_etapa1…` (1) | ✅ | [`documentacion/terminado.md`](../documentacion/terminado.md) |
| `FIX-37` Nivel de acceso en el perfil | `resource_access…` (1) | ✅ | [`documentacion/terminado.md`](../documentacion/terminado.md) |
| `FIX-39` Completar `BAC-21B` | `puente_etapa1…` (3) | ✅ | [`documentacion/terminado.md`](../documentacion/terminado.md) |
| `BAC-18B` Índices, particiones y purga | `cierre_fase_base…` (2) | 🟡 | [`documentacion/terminado.md`](../documentacion/terminado.md) (Fix en [`futuro.md: FIX-44`](../documentacion/futuro.md)) |
| `FRN-17C` Ticket efímero para eventos | `events-client.test.tsx` (4) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |

### Etapa 1 (Ola 1)

| Tarea | Pruebas | Estado | Destino Documental |
|---|---|---|---|
| `INF-07B` Nginx para SSE | `puente_etapa1…` (1) | ✅ (ref. local) | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (pendiente CT 103 en `INT-03`) |
| `FRN-20A` Tabla de inventario | `instances-table.test.tsx` (5) | 🟡 | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (Fix en [`futuro-1.md: FIX-43`](../documentacion/futuro-1.md)) |
| `BAC-29` Contrato de la etapa | `etapa1_acceptance…` (1) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `BAC-22` Telemetría del nodo | `etapa1_acceptance…` (1) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `BAC-23A` Adaptador de IP | `etapa1_acceptance…` (1) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `BAC-25A` Pool de UPID | `etapa1_acceptance…` (1) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `FRN-19A` Medidores del host | `dashboard-metrics.test.tsx` (3) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `FRN-17A` Consumo de eventos | `events-client.test.tsx` (4) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |

---

## 3. Detalle de cambios incorporados en los submódulos

1. **Backend (`0167b96` → `4e204f1`):**
   - **`96106a6`:** Particionamiento declarativo por rango trimestral de `auditoria` en Postgres, índices compuestos y `PurgaWorker` periódico para `sesiones_activas` (`BAC-18B`).
   - **`b8d2631` y `42e94ca`:** `InstanciaListadaDTO` ampliado con campos nuevos de telemetría y `nivelAcceso`, incorporación de `shutdown` y `reboot` en `ProxmoxPort`, y auditoría de acciones de energía en estado `PENDING` (`FIX-39`).
   - **`4e204f1`:** Incorporación de `permisos: [{ vmid, nivelAcceso }]` en `UsuarioDetalleDTO` para `GET /account/profile` (`FIX-37`).
2. **Frontend (`4e0e7b7` → `070e96b`):**
   - **PR #74 y #75 (Cristian):** Integración reactiva del inventario en `Instances.tsx` con `useInstances.ts`, `instanceService.ts`, `InstanceAction.tsx` y filtrado por `canAccessInstance`.
   - **PR #76 y #77 (Belinda y Luz):** Restauración del selector de filtro por "Acción" en `Auditoria.tsx` (`FIX-36`), query parameter reactivo y exportación CSV.

---

## 4. Prueba integral LOGIN-04

**10 de 10 pasos en verde** (frontend real ejecutado por Vitest contra el backend real en Docker Compose).

---

## 5. Pendientes actuales

- **Fase Base:**
  - Implementar `FRN-17C` (ticket efímero en cliente de eventos).
  - Corregir `FIX-44` (`BAC-18B`: incluir `jti_access` en el índice parcial de Postgres).
- **Etapa 1:**
  - Corregir `FIX-43` (`FRN-20A`: columna IP, badges y formato de celdas en `Instances.tsx`).
  - Desarrollar `BAC-29`, `BAC-22`, `BAC-23A`, `BAC-25A`, `FRN-19A` y `FRN-17A`.
  - Desplegar `docker/nginx-edge.conf` en el contenedor CT 103 de Proxmox (`INT-03` / `INF-07B`).

---

## 6. Cómo reproducir

```bash
# 1. Asegurar últimas revisiones de submódulos
git -C backend pull origin main
git -C frontend pull origin main

# 2. Correr suite completa de Backend (Docker Compose + LOGIN-04 front↔back + simulador)
(cd test/back && go test -v -count=1 ./...)

# 3. Correr suite completa de Frontend
pnpm --dir test/front test --run
```
