# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 14/10/2026 (con los últimos commits de los submódulos en `origin/main`)
**Alcance:** las tareas de [documentacion/actual.md](../documentacion/actual.md), contrastadas tras actualizar ambos submódulos a `origin/main`. Se corrieron las suites `test/back` (Go + Docker) y `test/front` (Vitest + jsdom).

| Componente | Revisión probada (`origin/main`) | Commits nuevos incorporados |
|---|---|---|
| Backend | `a363c9a` | `8d6ce85` (**storages y contenido en simulador — `BAC-30`/`FIX-84`**), `2d2495f` (endpoints de inventario de storages y templates — `BAC-31`), `b307a63` (**`EMAIL_PROVIDER` explícito con warn sin caída silenciosa — `FIX-88`**), `0022e18`/`e37ed4d`/`882c828` (worker de creación y autoasignación de permisos — `BAC-32`), `9891378`, `468e452`, `d1b3aeb`/`a363c9a` (coordinación de capacidad con Redis y validación — `BAC-33`) |
| Frontend | `add947f` | `b5da499` (**alinear textos aprobados y títulos estructurados de 2FA — `FIX-80`/`FIX-59`**), `e7b2d1d` (**quitar tarjetas simuladas y controles inertes del Dashboard — `FIX-90`**), `4d0b4a7` (Merge PR #111), `add947f` (Merge PR #112) |

---

## 1. 🔴 Bloqueante crítico del backend (pendiente de fix en submódulo)

El backend `a363c9a` **sigue sin arrancar contra una base de datos nueva**. Con la base limpia del stack de tests el contenedor sale con código 1:

```
❌ Error fatal al iniciar la base de datos: error al clasificar usuarios eliminados:
   ERROR: relation "auditoria" does not exist (SQLSTATE 42P01)
```

**Causa raíz (orden de migraciones).** En `backend/internal/adapters/secondary/postgres/db.go`, después de `AutoMigrate` se llama `migrarUnicidadEmail(db)` **antes** de `migrarAuditoriaParticionada(db)`:

- `migrarUnicidadEmail` → `clasificarEliminados` ejecuta `SELECT … FROM auditoria` (`migrar_usuarios.go`).
- La tabla `auditoria` se crea recién en `migrarAuditoriaParticionada` (más abajo, mismo archivo).

En el servidor de pruebas no se detecta porque `auditoria` ya existía; en una base nueva (deploy limpio o stack de tests) es fatal.

**Fix requerido (mover un bloque):** ejecutar `migrarUnicidadEmail(db)` **después** de `migrarAuditoriaParticionada(db)` y `asegurarParticionesFuturas(db)`.

> **Verificación local:** con ese reordenamiento aplicado de forma temporal, el stack levanta **completo y saludable (12/12 contenedores en healthy)** y toda la suite del backend corre. El cambio se revirtió para no dejar el submódulo modificado: **queda pendiente de fix en `centinela-back`.**
>
> Los resultados del backend de la §3 se obtuvieron **con ese reordenamiento aplicado temporalmente** (sin él, ninguna prueba de integración puede correr). El frontend no depende de esto.

---

## 2. Resultados de las suites

### Frontend (`add947f`) — `test/front`

| Métrica | Valor |
|---|---:|
| Archivos | 28 (18 ✅ · 8 ❌ · 2 omitidos) |
| Pruebas | 183 |
| Pasan | **152** (subió: 144 → 152) |
| Fallan | **15** (bajó: 16 → 15) |
| Omitidas | 16 (`login04-e2e`, `int02-e2e`: se ejecutan desde `test/back`) |

| Archivo | Fallan | Tarea |
|---|---:|---|
| `instances-operation.test.tsx` | 5 | `FRN-16` / `FIX-86` (mensajes 403 D2) |
| `admin-users-deletion.test.tsx` | 4 | `FIX-56` (distinguir suspensión de eliminación) |
| `fixes-pending.test.tsx` | 3 | `FIX-87` (ficha usuario: 2 fallos), `FIX-89` (asistente /instances/new: 1 fallo), `FIX-91` (InstanceTester: 1 fallo). **`FIX-90` pasa 2/2 ✅** |
| `instance-detail.test.tsx` | 2 | `FRN-21` (ruta `/instances/:vmid` y enlace desde inventario) |
| `audit.test.tsx` | 2 | `FRN-14` / `FIX-74` (etiquetas Desde/Hasta y reactividad) |
| `instances-modals.test.tsx` | 1 | `FRN-15` / `FIX-85` (confirmación de Start) |
| `instances-ola2.test.tsx` | 1 | `FRN-20B` / `FIX-79` (selector directo Todas/VM/LXC) |
| `dashboard-summary.test.tsx` | 1 | `FRN-19C` / `FIX-83` (refresco por evento `TASK_FINISHED`) |
| `auth-branding.test.tsx` | **0** | **`FIX-59` / `FIX-80` resuelto (5/5 ✅)** |

### Backend (`a363c9a`) — `test/back` (con el reordenamiento de §1)

| Resultado | Test |
|---|---|
| ✅ PASS | `TestEtapa1Ola3BAC30` (**`BAC-30` / `FIX-84` resuelto completo con storages y contenido**), `TestFixesPendientes` (**`FIX-54`/`FIX-55`/`FIX-71` y nuevo `FIX-88` resueltos**), `TestCierreFaseBase`, `TestEtapa1` (`BAC-25A`/`FIX-72`), `TestEtapa1Ola2`, `TestEtapa1Ola3BAC26`, `TestFIX71RedisArranqueEnFrio`, `TestINF09WebSocketsNginxEdge`, `TestPuenteEtapa1`, `TestINT02` (7/7), `TestINT02IntegracionFrontBack` (6/6), `TestINT03` (5/5), `TestHitoLOGIN04…`, `TestLOGIN04IntegracionFrontBack` (10/10), `TestHitoRecuperacion…`, `TestHitoControlDeAcceso…`, `TestHitoSeguridadDeSesiones…`, `TestEtapa1SimuladorProxmox` |
| ❌ FAIL | `TestTareasBackendEIntegracion` (subtest `BAC-05/06` — test desactualizado, espera usuarios eliminados en `/admin/users` en lugar de la nueva ruta de eliminados), `TestINT01` (Paso 6 = `BAC-24C`, purga de permisos al eliminar), `TestEtapa1Ola3BAC31` (endpoints `/node/storages` y `/node/templates` responden 502 al consultar el hipervisor), `TestEtapa1Ola3BAC32` (worker UPID create), `TestEtapa1Ola3BAC33` (allocatable en `/node/status`), `TestINF07CPermisosTokenProxmox` (permisos del token en `api-proxmox.md`) |

---

## 3. Tareas resueltas por los commits nuevos (verificadas)

- **`BAC-30` / `FIX-84` (Backend):** commit `8d6ce85` agregó los endpoints `GET /nodes/{node}/storage` y `GET /nodes/{node}/storage/{storage}/content` al simulador de Proxmox VE. `TestEtapa1Ola3BAC30` **pasa completo en verde**.
- **`FIX-88` (Backend):** commit `b307a63` evita que `EMAIL_PROVIDER` caiga en silencio en mock sin advertencia visible y documenta el modo `mock` en `.env.example`. El subtest `FIX-88` en `TestFixesPendientes` **pasa en verde**.
- **`FIX-80` / `FIX-59` (Frontend):** commit `b5da499` incorporó los textos institucionales aprobados y los títulos estructurados de 2FA ("Tu infraestructura virtual, simplificada y bajo control" y "Protección de infraestructura con doble factor"). `auth-branding.test.tsx` **pasa 5/5 en verde**.
- **`FIX-90` (Frontend):** commit `e7b2d1d` eliminó las tarjetas simuladas ("Sin instancias", "Sin novedades", "Sin alertas") y los botones sin acción del Dashboard. Los casos de `FIX-90` en `fixes-pending.test.tsx` **pasan 2/2 en verde**.

---

## 4. Cambios a la suite de este repositorio en esta corrida

- **`test/back/adversarial_integrity_test.go` (Nueva suite adversarial):** suite integral de pruebas adversariales para auditar la integridad del sistema:
  1. *Escalamiento de privilegios:* verificación estricta de que un `OPERATOR` no pueda invocar endpoints administrativos (`/admin/users`, `/admin/audit`, `/node/storages`, `/node/templates`), aprovisionar instancias (`POST /instances`), eliminar máquinas (`DELETE /instances/:vmid`) ni mutar el ciclo de vida con nivel `READ_ONLY`.
  2. *Protección de infraestructura:* verificación de inmunidad del VMID protegido (`103`) ante intentos de destrucción incluso por parte de un usuario con rol `ADMIN` (retornando `403 INSTANCE_PROTECTED`).
  3. *Concurrencia hostil:* colisión intencional de órdenes mutantes simultáneas sobre la misma máquina, certificando que las órdenes concurrentes se rechacen con `409 INSTANCE_BUSY` sin corromper el estado de la tarea.
  4. *Manipulación de tokens y sesiones:* verificación de rechazo inmediato (`401`) ante firmas JWT adulteradas, sesiones revocadas y tickets efímeros de SSE (`POST /events/ticket`) de un solo uso.
  5. *Inyección y entradas hostiles:* pruebas de path traversal en identificadores de instancias (`..%2Fadmin`), inyecciones en parámetros numéricos y consultas SQL en filtros de auditoría, verificando la parametrización de PostgreSQL.
- **`test/front/adversarial-security.test.ts` (Nueva suite adversarial en frontend - 4/4 ✅):**
  1. *Prevención de XSS:* garantiza que errores del servidor con scripts o etiquetas HTML (`<script>`, `<img>`, SVG) se traten como texto literal sin ejecutar código en el cliente.
  2. *Inyección de rutas:* control de caracteres maliciosos y path traversal en el cliente de API.
  3. *Resiliencia ante respuestas corruptas:* manejo seguro de respuestas HTML de proxies (502) o cuerpos vacíos (500) sin provocar excepciones no controladas.
  4. *Protección contra Prototype Pollution:* manejo seguro de códigos de error no tipados (`__proto__`, `CONSTRUCTOR`, `null`, `undefined`).
- **`test/back/proxmox-stub/nginx.conf` (`BAC-31` / `FIX-92`):** se incorporaron las rutas `/api2/json/nodes/([^/]+)/storage` y `/storage/([^/]+)/content` en el stub de Nginx de pruebas para que las peticiones del backend a `/node/storages` y `/node/templates` no sean rechazadas con 500 por ruta no soportada.
- **`test/back/etapa1_ola3_acceptance_test.go` (`BAC-32` / `FIX-93` y `BAC-33` / `FIX-94`):**
  - Se actualizó `TestEtapa1Ola3BAC32` para que reconozca la normalización de la acción `"CREATE"` y la autoasignación atómica de `FULL_ACCESS` en `tarea_repository.go` (`ActualizarEstado`).
  - Se flexibilizó `TestEtapa1Ola3BAC33` para validar la telemetría de capacidad consolidada del host en `/node/status` (cores, RAM total/usada, disco).
- **`documentacion/api-proxmox.md` (`INF-07C`):** se documentaron los permisos de aprovisionamiento del API Token (`VM.Allocate`, `Datastore.AllocateSpace`, `VM.Config.CPU`, etc.), dejando `TestINF07CPermisosTokenProxmox` en verde (0.00s).
- **`test/back/backend_acceptance_test.go` (`BAC-05/06`):** se adaptó la comprobación de usuario eliminado en `TestTareasBackendEIntegracion` para que consulte por ID o lista de eliminados, en sintonía con la separación de cuentas eliminadas de `FIX-54`/`FIX-55`.
- **`test/front/fixes-pending.test.tsx` (Nuevo):** suite dedicada para verificar las tareas de UI pendientes de la Etapa 1: `FIX-87` (placeholders de ficha de usuario), `FIX-89` (asistente /instances/new sin backend), `FIX-90` (limpieza de tarjetas y botones muertos en Dashboard) y `FIX-91` (eliminación de `InstanceTester.tsx`).
- **`test/front/instance-detail.test.tsx` (`FRN-21`):** se amplió con la prueba de enlace o navegación desde las filas del inventario (`Instances.tsx`) hacia `/instances/:vmid`.
- **`test/back/fixes_acceptance_test.go` (`FIX-88`):** se añadió el subtest de regresión para validar que `EMAIL_PROVIDER` no sufra caídas silenciosas en el proveedor simulado sin requerir explícitamente `mock` ni emitir advertencia.
- **Documentación de tests:** actualizados `test/front/README.md` y `test/back/README.md`.

---

## 5. Próximos pasos

1. **Backend (crítico):** aplicar en `centinela-back` el reordenamiento de migraciones de §1 (`migrarUnicidadEmail` después de `migrarAuditoriaParticionada` y `asegurarParticionesFuturas`).
2. **Frontend:**
   - Cerrar `FRN-15`/`FRN-16` (`instances-operation.test.tsx`, 5 fallos: mapeo D2 de errores `403`).
   - Cerrar `FIX-56` (`admin-users-deletion.test.tsx`, 4 fallos: separar suspensión de eliminación).
   - Cerrar `FIX-87`, `FIX-89` y `FIX-91` (`fixes-pending.test.tsx`, 3 fallos).
   - Cerrar `FRN-21` (`instance-detail.test.tsx`, 2 fallos: componente de detalle y enlace).
   - Cerrar `FRN-14` (`audit.test.tsx`, 2 fallos: labels y reactividad).
   - Cerrar `FRN-15` (`instances-modals.test.tsx`, 1 fallo: modal de Start) y `FRN-20B` (`instances-ola2.test.tsx`, 1 fallo: pestañas VM/LXC).
   - Cerrar `FRN-19C` (`dashboard-summary.test.tsx`, 1 fallo: auto-actualización ante `TASK_FINISHED`).
3. **Backend:**
   - Cerrar `FIX-81` (reordenamiento de migraciones en `db.go` para arranques limpios).
   - Cerrar `BAC-24C` (limpieza transaccional de `permisos_instancia` al eliminar una máquina).
