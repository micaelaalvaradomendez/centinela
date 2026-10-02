# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 01/10/2026 (tercera verificación; cobertura integral de Etapa 1 agregada a las suites de aceptación)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md) (Fase Base y Etapa 1), tras implementar las pruebas de aceptación que faltaban, y la regresión de [terminado.md](../documentacion/terminado.md) y [terminado-1.md](../documentacion/terminado-1.md).

| Componente | Revisión probada | Commits nuevos desde la verificación anterior |
|---|---|---|
| Backend | `44a2339` (último commit de `main`) | Ninguno (`44a2339` "mas arreglos del simulador", **FIX-35**) |
| Frontend | `4e0e7b7` (último commit de `main`) | 4 commits (PR #71 de Luz: refactor TypeScript de `AuthContext` y validación de `isAdmin` en `informationOfUser`) |

**Cómo se ejecutó:** submódulos sincronizados con `origin/main`. Se agregaron las pruebas de aceptación automatizadas faltantes para validar todo el backlog activo de `actual.md`. Ambas suites se ejecutaron de punta a punta sin contenedores residuales.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 58 | **44** | **12** | 2 |
| Frontend (`test/front`) | 111 | **75** | **26** | 10 |
| **Total** | **169** | **119** | **38** | **12** |

- Los 10 omitidos del frontend son la prueba integral LOGIN-04, que corre desde `test/back` contra el backend real (pasa **10/10**).
- Los 2 omitidos del backend son el `DELETE` de `BAC-21B` (espera a `BAC-24B`) y el relay SMTP Brevo de `INF-08B` (dependencia de red).
- Se agregaron **19 casos de prueba nuevos** para cubrir el 100% de las tareas activas de `actual.md` que carecían de aceptación (§3).
- **Hallazgo de integración:** el PR #71 de Frontend incorporó `useAuth()` dentro de `informationOfUser.tsx` e introdujo un bug de casing en `detailsUserPage.tsx` (`requiredRole="Admin"` con minúscula en lugar de `"ADMIN"`, lo que oculta las pestañas de administración).

---

## 2. Estado de las tareas de `actual.md`

**Leyenda:** ✅ cumplida · 🟡 parcial · ❌ no implementada.

### Fase Base

| Tarea | Pruebas | Estado | Detalle |
|---|---|---|---|
| `SEC-03` Permisos reactivos | `navigation.test.tsx` (3) | 🟡 | `usePermissions` existe en `src/hooks/` con `isAdmin`, `canAccessInstance` y `canOperateInstance`. Faltan `isOperator`, `hasRole` y `PermissionGate`, y el menú sigue mostrando "Auditoría" al OPERATOR |
| `FIX-29` Complejidad de contraseña | `password-change` (3), `recover-password` (1) | ✅ | Complejidad de contraseña con mayúscula validada y mensajes canónicos |
| `FRN-17C` Cliente de eventos (ticket efímero) | `events-client.test.tsx` (4) | ❌ | No existe `useEvents` que pida ticket efímero ni aplique retroceso exponencial |
| `FIX-36` Filtro "Acción" en Auditoría | `audit.test.tsx` (1) | ❌ | Falta restaurar el selector y parámetro reactivo de `accion` en `Auditoria.tsx` |
| `FIX-38` `READ_ONLY` como rol | `admin-users` (1), `navigation` (1) | ❌ | El selector "Rol" de la ficha sigue ofreciendo `READ_ONLY` y el hook lo acepta como rol |
| `BAC-18B` Índice parcial, particiones y purga | `cierre_fase_base…` (2) | ❌ | Faltan el índice parcial, las particiones trimestrales y la rutina de purga horaria |
| `BAC-21B` / `FIX-39` Instancias extendidas y energía | `puente_etapa1…` (5) | 🟡 | Pasa `start`/`stop` con `FULL_ACCESS` y `tareas_asincronas`. Faltan campos en `GET /instances`, auditoría de órdenes de energía y `shutdown`/`reboot` |
| `FIX-37` Nivel de acceso en el perfil | `resource_access…` (1) | ❌ | `GET /account/profile` sigue sin el campo `permisos: [{ vmid, nivelAcceso }]` |
| `FIX-40` Distinguir 504 de 502 en Proxmox | `etapa1_core…` (1, **nuevo**) | ❌ | `instance_handler.go` sigue respondiendo `PROXMOX_UNAVAILABLE` en timeouts 504 en vez de `PROXMOX_TIMEOUT` |

### Etapa 1

| Tarea | Pruebas | Estado | Detalle |
|---|---|---|---|
| `INF-07B` Nginx SSE en el borde | `puente_etapa1…` (1) | ✅ | El Nginx de borde entrega `TASK_FINISHED` por HTTPS sin buffer ni cortes |
| `BAC-29` Contrato HTTP y eventos Etapa 1 | `etapa1_core…` (1, **nuevo**) | ❌ | No existe aún `backend/docs/contrato-etapa1.md` ni sus anotaciones Swagger |
| `BAC-22` Telemetría del nodo con caché Redis | `etapa1_core…` (1, **nuevo**) | ❌ | No existe el endpoint `GET /api/node/status` (responde 404) |
| `BAC-23A` Normalización inventario QEMU/LXC/IP | `etapa1_core…` (1, **nuevo**) | ❌ | Falta resolución de interfaces y consolidación de VMs/LXC en `proxmox/client.go` |
| `BAC-25A` Worker pool acotado para UPID | `etapa1_core…` (1, **nuevo**) | ❌ | `seguimiento_tareas.go` sigue usando goroutines libres sin worker pool (`UPID_WORKERS`) |
| `FRN-19A` Medidores de recursos del Host | `dashboard-metrics.test.tsx` (4, **nuevo**) | ❌ | Faltan componentes en `features/dashboard` con gauges de CPU/RAM/Disco y umbral al 70% |
| `FRN-20A` Tabla interactiva con badges e IP | `instances-table.test.tsx` (4, **nuevo**) | ❌ | `Instances.tsx` sigue siendo una maqueta estática ("Sin instancias") |
| `FRN-17A` Consumo de eventos y distribución | `events-client.test.tsx` (3, **nuevo**) | ❌ | Falta parseo tipado de `RealtimeEvent`, deduplicación por ID y suscripción por recurso |

---

## 3. Pruebas que faltaban (creadas en esta verificación)

Se construyeron e incorporaron pruebas de aceptación formales para todas las tareas que no tenían cobertura en `actual.md`:

| Área | Archivo | Casos nuevos | Qué valida |
|---|---|---:|---|
| Backend | `test/back/etapa1_core_acceptance_test.go` | 5 | **BAC-29:** existencia de `contrato-etapa1.md` con esquemas de status, instances, energía y códigos de error.<br>**FIX-40:** mapeo de `ErrProxmoxTimeout` a `504 PROXMOX_TIMEOUT` vs `502 PROXMOX_UNAVAILABLE`.<br>**BAC-22:** `GET /api/node/status` exige auth (401 sin token, 200 con OPERATOR/ADMIN) y devuelve telemetría normalizada.<br>**BAC-23A:** resolución de interfaces QEMU y LXC con IP unificada.<br>**BAC-25A:** worker pool acotado para seguimiento de tareas UPID |
| Frontend | `test/front/dashboard-metrics.test.tsx` | 4 | **FRN-19A:** componentes de métricas de host en `features/dashboard`, comportamiento responsive 0-100%, y cambio de estado a advertencia exactamente a partir del 70% (valores de borde 69% y 70%) |
| Frontend | `test/front/instances-table.test.tsx` | 4 | **FRN-20A:** tabla dinámica de inventario con VMs y LXC en la misma grilla, badges de estado (`Running` verde / `Stopped` gris), IP con botón de copia o "No detectada", y tolerancia a campos en null |
| Frontend | `test/front/events-client.test.tsx` | 3 | **FRN-17A:** parseo y validación de `RealtimeEvent` descartando payloads inválidos, deduplicación por `id`, y suscripción selectiva por tipo y `recursoId` |

---

## 4. Adaptaciones y Hallazgos de Integración

1. **Adaptación a PR #71 (`AuthContext` tipado en Frontend):**
   - El commit `952b43c` introdujo `const { isAdmin } = useAuth()` en `informationOfUser.tsx`.
   - Se actualizó el fixture de prueba `renderUserDetail` en `test/front/admin-users.test.tsx` envolviéndolo con `<AuthProvider>` y poblando la sesión en `localStorage` (`centinela_user`), permitiendo que el componente monte limpiamente sin violar el contexto.

2. **Detección de Defecto en `detailsUserPage.tsx`:**
   - En el commit `1156eaa0` de Luz, las pestañas de "Roles y permisos" en `detailsUserPage.tsx` fueron envueltas con `<PermissionGate requiredRole="Admin">`.
   - Dado que los roles del sistema son estrictamente en mayúsculas (`ADMIN` / `OPERATOR`), la comparación `user.rol === role` evalúa a `false`, provocando que las pestañas queden ocultas incluso para administradores. Queda documentado para su corrección por el equipo de frontend.

---

## 5. Prueba integral LOGIN-04

**10 de 10 pasos en verde** (frontend real ejecutado por Vitest contra el backend real en Docker Compose).

---

## 6. Cómo reproducir

```bash
# 1. Asegurar últimas revisiones de los submódulos
git submodule update --init --recursive

# 2. Correr suite completa de Backend (incluye Docker Compose, LOGIN-04 front↔back y simulador Proxmox)
cd test/back && go test -v -count=1 ./...

# 3. Correr suite completa de Frontend
pnpm --dir test/front test --run
```
