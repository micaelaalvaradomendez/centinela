# Informe de estado de tareas verificado por pruebas

**Fecha de actualización:** 06/10/2026
**Alcance:** sincronización de los submódulos `backend` (`eec77ff`) y `frontend` (`5a86dce`) con `origin/main` y ejecución completa de las suites de [test/back](back/README.md) y [test/front](front/package.json).

| Componente | Revisión probada | Commits nuevos desde la verificación anterior |
|---|---|---|
| Backend | `eec77ff` (`origin/main`) | Implementación de `BAC-25A` (pool acotado de workers UPID), `BAC-29` (contrato Etapa 1), `BAC-22` (telemetría con Redis) y `BAC-23A` (adaptador de inventario e IPs). |
| Frontend | `5a86dce` (`origin/main`) | Implementación de `FRN-17C` y `FRN-17A` (cliente de eventos con tickets efímeros y contexto) y resolución de `FIX-43` / `FRN-20A` (maquetado de inventario, IP con copia y badges). |

**Actualización de submódulos:** ambos submódulos quedaron alineados con `origin/main` (`backend: eec77ff`, `frontend: 5a86dce`).

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`, `eec77ff`) | 61 | **54** | **6** | 1 |
| Frontend (`test/front`, `5a86dce`) | 113 | **99** | **4** | 10 |
| **Total** | **174** | **153** | **10** | **11** |

- **Avances clave en Frontend:**
  - **`FRN-20A` / `FIX-43` está 100% verificado y en verde:** PR #83 y #84 incorporaron la separación visual de nombre y VMID, columna IP con botón de copia al portapapeles, badges de estado y diseño responsivo. `instances-table.test.tsx` pasa **5/5 en verde ✅**.
  - **`FRN-17C` y `FRN-17A` implementados:** PR #80, #81 y #82 implementaron la solicitud de tickets efímeros y el contexto de eventos en tiempo real. 4 pruebas de `events-client.test.tsx` pasan; restan 4 observaciones derivadas a `FIX-45` (en `futuro.md`) y `FIX-49` (en `futuro-1.md`).
- **Avances en Backend:**
  - **`BAC-25A` está 100% completo y verificado:** `seguimiento_tareas.go` implementa el pool acotado de workers (`UPID_WORKERS`). La prueba de aceptación `TestEtapa1/BAC-25A_...` pasa **100% en verde ✅**.
  - **`BAC-29`, `BAC-22` y `BAC-23A` implementados:** se publicaron el contrato `contrato-etapa1.md`, el endpoint de telemetría `GET /api/node/status` con Redis y el servicio de inventario unificado. Sus detalles pendientes quedaron derivados a `FIX-46`, `FIX-47` y `FIX-48` en `futuro-1.md`.
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
| `FRN-17C` Ticket efímero para eventos | `events-client.test.tsx` (4) | 🟡 | [`documentacion/terminado.md`](../documentacion/terminado.md) (Fix en [`futuro.md: FIX-45`](../documentacion/futuro.md)) |

### Etapa 1 (Ola 1)

| Tarea | Pruebas | Estado | Destino Documental |
|---|---|---|---|
| `INF-07B` Nginx para SSE | `puente_etapa1…` (1) | ✅ (ref. local) | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (pendiente CT 103 en `INT-03`) |
| `FRN-19A` Medidores del host | `dashboard-metrics.test.tsx` (3) | ✅ | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (PR #78) |
| `FRN-20A` Tabla de inventario | `instances-table.test.tsx` (5) | ✅ | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (Resuelto con `FIX-43`, PR #83/#84) |
| `BAC-25A` Pool de UPID | `etapa1_acceptance…` (1) | ✅ **Nuevo** | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (commit `3988546`) |
| `BAC-29` Contrato de la etapa | `etapa1_acceptance…` (1) | 🟡 | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (Fix en [`futuro-1.md: FIX-46`](../documentacion/futuro-1.md)) |
| `BAC-22` Telemetría del nodo | `etapa1_acceptance…` (1) | 🟡 | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (Fix en [`futuro-1.md: FIX-47`](../documentacion/futuro-1.md)) |
| `BAC-23A` Adaptador de IP | `etapa1_acceptance…` (1) | 🟡 | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (Fix en [`futuro-1.md: FIX-48`](../documentacion/futuro-1.md)) |
| `FRN-17A` Consumo de eventos | `events-client.test.tsx` (4) | 🟡 | [`documentacion/terminado-1.md`](../documentacion/terminado-1.md) (Fix en [`futuro-1.md: FIX-49`](../documentacion/futuro-1.md)) |

---

## 3. Detalle de cambios incorporados en los submódulos

1. **Frontend (`907efe5` → `5a86dce`):**
   - **PR #80, #81, #82 (Cristian):** Implementación de `FRN-17C` y `FRN-17A`. Añadido `eventsClient.ts` con llamada a `POST /api/events/ticket` e inicialización de EventSource/SSE. Creación de `EventsContext.tsx` y `useEvents.ts`.
   - **PR #83, #84 (Luz / Cristian):** Implementación de `FIX-43` (`FRN-20A`). Separación de nombre e ID de instancia, columna IP con botón interactivo de copia al portapapeles (`useCopyInstanceIp`), badges de estado diferenciados (`InstanceBadges.tsx`) y estilización de tabla.
2. **Backend (`4e204f1` → `eec77ff`):**
   - **`3988546` (Lisandro):** Implementación de `BAC-25A`. Creación de pool de workers acotado (`UPID_WORKERS`) en `seguimiento_tareas.go` para monitoreo de tareas de Proxmox con canal amortiguado y cierre ordenado.
   - **`5787179`, `21b1332`, `4ef36e3`, `0850ec7` (Lisandro):** Implementación de `BAC-29`. Publicación de `backend/docs/contrato-etapa1.md` y anotaciones Swagger.
   - **`13f9c35` (Tayra):** Implementación de `BAC-22`. Endpoint `GET /api/node/status` con normalización de métricas de host y caché Redis de telemetría con lectura stale.
   - **`d1dec4f` (Lisandro):** Implementación de `BAC-23A`. Adaptador de inventario y resolución de interfaces IP para VMs y contenedores LXC.

---

## 4. Prueba integral LOGIN-04

**10 de 10 pasos en verde** (frontend real ejecutado por Vitest contra el backend real en Docker Compose).

---

## 5. Pendientes actuales

- **Fase Base:**
  - Corregir `FIX-44` (`BAC-18B`: incluir `jti_access` en el índice parcial de Postgres).
  - Corregir `FIX-45` (`FRN-17C`: redirección inmediata a `/login` y cierre de sesión local ante 401 en ticket efímero).
- **Etapa 1:**
  - Corregir `FIX-46` (`BAC-29`: documentar permisos de `GET /api/node/status` en bloque del endpoint).
  - Corregir `FIX-47` (`BAC-22`: ciclo de recuperación y respuesta 200 tras caída de Proxmox en telemetría).
  - Corregir `FIX-48` (`BAC-23A`: incluir literales de rutas de red en pruebas unitarias del cliente Proxmox).
  - Corregir `FIX-49` (`FRN-17A`: deduplicación por ID, conexión SSE única compartida y suscripción por tipo/recurso).
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
