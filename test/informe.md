# Informe de estado de tareas verificado por pruebas

**Fecha de actualización:** 03/10/2026
**Alcance:** sincronización de los submódulos `backend` (`4e204f1`) y `frontend` (`907efe5`) con `origin/main` y ejecución completa de las suites de [test/back](back/README.md) y [test/front](front/package.json).

| Componente | Revisión probada | Commits nuevos desde la verificación anterior |
|---|---|---|
| Backend | `4e204f1` (`origin/main`) | 0 (se mantiene en `4e204f1`, con `FIX-37`, `FIX-39` y `BAC-18B` incorporados). |
| Frontend | `907efe5` (`origin/main`) | 2 commits nuevos (PR #78 de Belinda: **`FRN-19A`**, maquetado y medidores de recursos del Host: CPU, RAM y almacenamiento con tests de componentes). |

**Actualización de submódulos:** ambos submódulos quedaron alineados con `origin/main` (`backend: 4e204f1`, `frontend: 907efe5`).

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`, `4e204f1`) | 61 | **53** | **7** | 1 |
| Frontend (`test/front`, `907efe5`) | 113 | **90** | **13** | 10 |
| **Total** | **174** | **143** | **20** | **11** |

- **Avance clave en Frontend:**
  - **`FRN-19A` está 100% completo y verificado:** PR #78 incorporó `ResourceMeter.tsx`, integró los medidores en `Dashboard.tsx` y sumó pruebas unitarias propias (`resource-meters.test.tsx`). La suite de aceptación `dashboard-metrics.test.tsx` pasa **3/3 en verde ✅**.
  - **Fase Base Frontend:** se mantiene **100% verde** (`SEC-03`, `FIX-29`, `FIX-36`, `FIX-38`, `FIX-42`).
  - Los 13 fallos restantes en frontend son: 8 de eventos (`FRN-17C` / `FRN-17A`) y 5 de la tabla de inventario (`FRN-20A`, cubierta por `FIX-43`).
- **Avances en Backend:**
  - Se mantienen en verde `FIX-37` (permisos en `GET /account/profile`), `FIX-39` (DTO de instancias, energía y auditoría en `PENDING`), `FIX-40` (`504 PROXMOX_TIMEOUT`), flujos de seguridad y LOGIN-04 (10/10 pasos).
  - Los 7 fallos corresponden a: 2 de `BAC-18B` (índice parcial sin `jti_access` y ticker del worker, cubierto por `FIX-44`), 1 de Redis local (colisión de nombre de contenedor en ejecución) y 4 tareas pendientes de la Etapa 1 (`BAC-29`, `BAC-22`, `BAC-23A`, `BAC-25A`).
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
| `FRN-19A` Medidores del host | `dashboard-metrics.test.tsx` (3) | ✅ **Nuevo** | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (PR #78) |
| `FRN-20A` Tabla de inventario | `instances-table.test.tsx` (5) | 🟡 | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (Fix en [`futuro-1.md: FIX-43`](../documentacion/futuro-1.md)) |
| `BAC-29` Contrato de la etapa | `etapa1_acceptance…` (1) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `BAC-22` Telemetría del nodo | `etapa1_acceptance…` (1) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `BAC-23A` Adaptador de IP | `etapa1_acceptance…` (1) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `BAC-25A` Pool de UPID | `etapa1_acceptance…` (1) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |
| `FRN-17A` Consumo de eventos | `events-client.test.tsx` (4) | ❌ | Permanece en [`documentacion/actual.md`](../documentacion/actual.md) |

---

## 3. Detalle de cambios incorporados en los submódulos

1. **Frontend (`070e96b` → `907efe5`):**
   - **PR #78 (Belinda):** Implementación de `FRN-19A`. Componente reutilizable `ResourceMeter.tsx` en `features/dashboard` con barras de progreso estilizadas para CPU, RAM y almacenamiento; integración en `Dashboard.tsx`; y suite unitaria `resource-meters.test.tsx` validando valores nominales y límites de umbral al 69% y 70%.
2. **Backend (`4e204f1`):**
   - Se mantiene en `4e204f1` (particionamiento trimestral de `auditoria` en Postgres, `PurgaWorker`, DTO ampliado con soporte para `shutdown`/`reboot` y auditoría en `PENDING`, y permisos por instancia en el perfil de usuario).

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
  - Desarrollar `BAC-29`, `BAC-22`, `BAC-23A`, `BAC-25A` y `FRN-17A`.
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
