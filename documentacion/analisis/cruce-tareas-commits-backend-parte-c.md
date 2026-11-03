# Cruce de Tareas vs Commits — Submódulo Backend (Parte C: 30/10/2026 al 03/11/2026)

**Fecha de generación:** 03/11/2026  
**Submódulo:** `backend/`  
**Rango analizado (Parte C):** `61b973d` (2026-10-30, commit #263) → `fb8a41c` (2026-11-02, commit #273, HEAD actual)  
**Corte previo (Parte B):** `908522e` (2026-10-29, commit #262 — `BAC-54`)  
**Total de commits en el tramo Parte C:** 11 commits (0 merges, 11 commits de código reales)  
**Total acumulado en el repositorio:** 273 commits (`cd6a83e` → `fb8a41c`)

---

## 1. Metodología y Continuidad con Partes A y B

Este documento continúa la trazabilidad histórica del submódulo `backend/`:
1. **Parte A (`cruce-tareas-commits-backend.md` / `cruce-tareas-commits-backend1.md` — 20/10/2026):** Cubrió los primeros **218 commits** (`#1 cd6a83e` al `#218 8daaf87`, del 06/09/2026 al 18/10/2026), abarcando `BAC-01` a `BAC-45` y `FIX-01` a `FIX-82`.
2. **Parte B (`cruce-tareas-commits-backend-parte-b.md` — 30/10/2026):** Cubrió **44 commits** (`#219 43b011e` al `#262 908522e`, del 19/10/2026 al 29/10/2026), abarcando `BAC-46` a `BAC-54` y los fixes `FIX-83` a `FIX-125`.
3. **Parte C (este informe — 03/11/2026):** Analiza los **11 commits posteriores** (`#263 61b973d` al `#273 fb8a41c`, del 30/10/2026 al 02/11/2026) y los cruza contra el bloque de adaptación al servidor real Proxmox VE 8.4 (`ADP-01` a `ADP-06` / `BAC-55`) y las tareas de corrección de Etapa 3 y Auditoría de Seguridad (`FIX-128` y `FIX-133` a `FIX-140`) documentadas en `documentacion/actual.md`, `documentacion/tareas-ajustes-etapa3.md` y `documentacion/servidorfinal/01_tareas_adaptacion_pve8.md`.

---

## 2. Panorama General de Commits por Integrante

### 2.1 Tramo Parte C (`#263` al `#273` — 30/10/2026 al 02/11/2026)

| Integrante | Usuario / Email en Git | Commits Totales | Merges | Commits de Código Real | % del Tramo (Código) |
|---|---|---:|---:|---:|---:|
| **Lisandro** | `lisandroureta` (`lisandroureta04@gmail.com`) | 6 | 0 | 6 | 54.5% |
| **Tayra** | `tayraag` (`tayraaguirre2001@gmail.com`) | 5 | 0 | 5 | 45.5% |
| **Nicolás** | `Nico` (`nicolasbutterfield02@gmail.com`) | 0 | 0 | 0 | 0.0% |
| **Lucas** | `LucasC19` (`lucascassano19@gmail.com`) | 0 | 0 | 0 | 0.0% |
| **Total Parte C** | — | **11** | **0** | **11** | **100.0%** |

### 2.2 Acumulado Histórico Completo (`#1` al `#273` — 06/09/2026 al 02/11/2026)

| Integrante | Parte A (#1–#218) | Parte B (#219–#262) | Parte C (#263–#273) | Total Commits | Merges Totales | Código Real Total | % Código Real |
|---|---:|---:|---:|---:|---:|---:|---:|
| **Lisandro** | 89 | 38 | 6 | **133** | 21 | **112** | 45.2% |
| **Tayra** | 65 | 4 | 5 | **74** | 0 | **74** | 29.8% |
| **Nicolás** | 38 | 2 | 0 | **40** | 0 | **40** | 16.1% |
| **Lucas** | 26 | 0 | 0 | **26** | 4 | **22** | 8.9% |
| **Total Histórico** | **218** | **44** | **11** | **273** | **25** | **248** | **100.0%** |

---

## 3. Lista Cronológica de Commits del Tramo Parte C (#263 al #273)

| # | Fecha | Hash | Autor | Mensaje del Commit | Archivos Principales Modificados | Tarea Asociada |
|---:|---|---|---|---|---|---|
| 263 | 2026-10-30 | `61b973d` | Lisandro | `feat(simulador): perfil de versión PVE 9.2.2 u 8.4.21 con PROXMOX_SIM_TARGET_VERSION` | `simulador-proxmox/main.go`, `simulador-proxmox/main_test.go`, `.env.example`, `docker-compose.yml` | **BAC-55 / ADP-05** |
| 264 | 2026-10-30 | `c09a867` | Tayra | `feat(proxmox): parametrización dinámica de VMIDs protegidos por entorno` | `internal/infraestructura/proxmox/client.go`, `client_test.go`, `sincronizador_test.go`, `cmd/api/main.go` | **ADP-04 (BAC-54)** |
| 265 | 2026-10-30 | `7505aa5` | Tayra | `feat(proxmox): ajustar parámetros de aprovisionamiento QEMU y LXC` | `internal/infraestructura/proxmox/client.go`, `client_test.go`, `simulador-proxmox/main.go` | **ADP-01** |
| 266 | 2026-10-30 | `084ad9a` | Tayra | `feat(proxmox): agregar purga profunda al eliminar instancias` | `internal/infraestructura/proxmox/client.go`, `client_test.go`, `simulador-proxmox/main.go` | **ADP-02** |
| 267 | 2026-10-30 | `166c108` | Tayra | `feat(proxmox): flexibilizar parser de interfaces y configuración` | `internal/infraestructura/proxmox/client.go`, `client_test.go` | **ADP-03** |
| 268 | 2026-10-30 | `ec5922c` | Tayra | `feat(simulador): completar modo dual de versiones PVE` | `simulador-proxmox/main.go`, `simulador-proxmox/main_test.go`, `.env.example` | **ADP-05 / BAC-55** |
| 269 | 2026-11-02 | `aee232a` | Lisandro | `feat(simulador): catálogo de plantillas e ISOs a Debian 12 Bookworm` | `simulador-proxmox/main.go`, `simulador-proxmox/main_test.go`, `internal/infraestructura/proxmox/client.go` | **ADP-06** |
| 270 | 2026-11-02 | `b2d5fec` | Lisandro | `fix(instancias): aceptar vlan 0 o null como red sin etiquetar` | `internal/dominio/instancia/modelo.go`, `internal/infraestructura/http/manejadores/instancia_handler.go`, `instancia_handler_test.go` | **FIX-128** (Back) |
| 271 | 2026-11-02 | `3fe055b` | Lisandro | `fix(auth): ampliar caracteres especiales de contraseñas según NIST SP 800-63B` | `internal/dominio/usuario/validador.go`, `validador_test.go` | **FIX-138** |
| 272 | 2026-11-02 | `4e4d2c4` | Lisandro | `fix(email): exponer el proveedor de correo simulado en /api/version, cabecera y log` | `cmd/api/main.go`, `internal/infraestructura/notificaciones/email.go` | **FIX-140** |
| 273 | 2026-11-02 | `fb8a41c` | Lisandro | `fix(cuotas): permitir CPU overcommit y validar solo el aumento de RAM` | `internal/aplicacion/servicios/cuota_servicio.go`, `cuota_servicio_test.go`, `instancia_servicio.go` | **FIX-133** |

---

## 4. Cruce Detallado por Tarea (Adaptación PVE 8.4 y Fixes Etapa 3 / Seguridad)

### 4.1 Bloque de Adaptación al Servidor Real Proxmox VE 8.4 (`ADP-01` a `ADP-06` y `BAC-55`)

| Código | Descripción de la Tarea | Responsable Asignado (`01_tareas_adaptacion_pve8.md` / `actual.md`) | Autor(es) en Git | Commits Asociados | Coincidencia de Autoría | Estado Técnico en Suite (`test/back`) |
|---|---|---|---|---|---|---|
| **ADP-01** | Compatibilidad de parámetros de aprovisionamiento QEMU y LXC con PVE 8.4 (`ostype=l26`, `scsihw=virtio-scsi-pci`, `balloon=0`, `cmode=tty`, `onboot=0`) | **Lisandro** | **Tayra** | `7505aa5` (#265) | ❌ **No coincide** *(Tayra)* | ✅ **Pasa** (`TestADP01_ParametrosAprovisionamientoQEMUyLXC`) |
| **ADP-02** | Purga profunda al eliminar instancias en Proxmox (`purge=1` y `destroy-unreferenced-disks=1`) | **Tayra** | **Tayra** | `084ad9a` (#266) | ✅ **Coincide** *(Tayra)* | ✅ **Pasa** (`TestADP02_PurgaProfundaAlEliminarInstancia`) |
| **ADP-03** | Resiliencia del parser de interfaces de red (`net0..net9`) y tipos dinámicos JSON en `/config` | **Lisandro** | **Tayra** | `166c108` (#267) | ❌ **No coincide** *(Tayra)* | ⚠️ **Incompleto en test externo**: El código en `client.go` y `client_test.go` fue implementado, pero el test de aceptación `TestADP03` exige explícitamente el archivo separado `internal/infraestructura/proxmox/client_pve8_test.go`, el cual no fue creado. |
| **ADP-04** (`BAC-54`) | Parametrización dinámica de `PROXMOX_PROTECTED_VMIDS` y `PROXMOX_DEFAULT_NODE` por variables de entorno | **Tayra** *(en `ADP-04`)* / **Lisandro** *(en `BAC-54`)* | **Lisandro** (`908522e` en Parte B) + **Tayra** (`c09a867` #264) | `c09a867` (#264) | ✅ **Coincide** *(Tayra/Lisandro)* | ⚠️ **Incompleto en test externo**: `ParseVmidsProtegidos` parsea listas separadas por comas (`100,101`), pero aún no soporta rangos con guion (`100-105`) exigidos por `TestADP04_ParametrizacionVMIDsProtegidosPorEntorno`. |
| **ADP-05 / BAC-55** | Modo Dual de Versiones en `simulador-proxmox` (`PROXMOX_SIM_TARGET_VERSION=8.4.21` vs `9.2.2`) | **Lisandro** | **Lisandro** (`61b973d` #263) + **Tayra** (`ec5922c` #268) | `61b973d` (#263), `ec5922c` (#268) | ⚠️ **Compartida** *(Lisandro + Tayra)* | ✅ **Pasa** (`TestADP05_SimuladorProxmoxModoDualVersiones`) |
| **ADP-06** | Unificación del catálogo de plantillas LXC e ISOs a **Debian 12 Bookworm** (`debian-12-standard_12.7-1_amd64.tar.zst`) | **Tayra** | **Lisandro** | `aee232a` (#269) | ❌ **No coincide** *(Lisandro)* | ✅ **Pasa** (`TestADP06_CatalogoPlantillasDebian12Bookworm`) |

### 4.2 Bloque de Correcciones de Etapa 3 y Auditoría de Seguridad (`FIX-128` y `FIX-133` a `FIX-140`)

| Código | Descripción de la Tarea | Responsable Asignado (`actual.md`) | Autor(es) en Git | Commits Asociados | Coincidencia de Autoría | Estado Técnico |
|---|---|---|---|---|---|---|
| **FIX-128 (Back)** | Aceptar `vlan: 0` o `null` como red sin etiquetar (untagged) en creación/edición de instancias | **Lisandro / Tayra** | **Lisandro** | `b2d5fec` (#270) | ✅ **Coincide** *(Lisandro)* | ✅ **Resuelto y verificado** |
| **FIX-133** | Sobreasignación de CPU (CPU Overcommit) y validación de cuota de RAM únicamente sobre el incremento neto ($\Delta \text{RAM} > 0$) | **Lucas / Micaela** | **Lisandro** | `fb8a41c` (#273) | ❌ **No coincide** *(Lisandro)* | ✅ **Resuelto y verificado** |
| **FIX-138** | Ampliar caracteres especiales permitidos en contraseñas según NIST SP 800-63B (incluyendo `_`, `-`, `.`, espacios e imprimibles ASCII) | **Lucas / Micaela** | **Lisandro** | `3fe055b` (#271) | ❌ **No coincide** *(Lisandro)* | ✅ **Resuelto y verificado** |
| **FIX-140** | Exponer proveedor de correo simulado (`email_provider: "mock-log"`) en `/api/version`, cabecera `X-Centinela-Email-Provider` y banner de consola | **Lucas / Micaela** | **Lisandro** | `4e4d2c4` (#272) | ❌ **No coincide** *(Lisandro)* | ✅ **Resuelto y verificado** |
| **FIX-134** | Paginación por defecto (`page=1`, `page_size=50`, `max=200`) y metadatos `X-Total-Count` en `GET /api/instances`, `/api/users` y `/api/quotas` | **Lucas / Nico** | *Sin commits* | — | ⏳ **Pendiente** | ❌ **Sin implementar al 03/11/2026** |
| **FIX-135** | Bloqueo de eliminación de usuario con instancias asignadas (`409 Conflict`) o reasignación atómica (`?transfer_to={admin_id}`) | **Lucas / Belinda** | *Sin commits* | — | ⏳ **Pendiente** | ❌ **Sin implementar al 03/11/2026** |
| **FIX-136** | Rate Limiting por IP y por cuenta en endpoints críticos de autenticación (`/api/auth/login`, `/recover`, `/reset`, `/refresh`) | **Nico / Lucas** | *Sin commits* | — | ⏳ **Pendiente** | ❌ **Sin implementar al 03/11/2026** |
| **FIX-137** | Auditoría del intento de cambio de contraseña fallido (`PASSWORD_CHANGE_FAILED`) en `audit_logs` | **Lucas / Micaela** | *Sin commits* | — | ⏳ **Pendiente** | ❌ **Sin implementar al 03/11/2026** |
| **FIX-139** | Cabeceras HTTP de defensa en profundidad (`X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, `CSP`) y validación de entropía de `JWT_SECRET` ($\ge 32$ bytes) | **Lucas / Nico** | *Sin commits* | — | ⏳ **Pendiente** | ❌ **Sin implementar al 03/11/2026** |

---

## 5. Hallazgos Clave del Tramo Parte C (Backend)

1. **Concentración total en 2 integrantes (Lisandro y Tayra):**
   - El 100% de los 11 commits del tramo Parte C fue ejecutado entre **Lisandro** (6 commits, 54.5%) y **Tayra** (5 commits, 45.5%).
   - **Lucas** y **Nicolás** mantienen **0 commits** en el repositorio `backend` durante este período (Lucas no registra commits en `backend` desde el 14/10/2026 en la Parte A, y Nicolás se concentró íntegramente en destrabar PRs del `frontend`).

2. **Inversión cruzada de asignaciones en el plan `ADP-01` a `ADP-06`:**
   - En `documentacion/servidorfinal/01_tareas_adaptacion_pve8.md`, `ADP-01` y `ADP-03` estaban asignadas a **Lisandro**, mientras que `ADP-06` estaba asignada a **Tayra**. En la práctica ocurrió exactamente a la inversa: **Tayra** implementó el bloque de cliente Proxmox (`ADP-01` en `7505aa5`, `ADP-02` en `084ad9a`, `ADP-03` en `166c108`, `ADP-04` en `c09a867` y el cierre de `ADP-05` en `ec5922c` el 30/10), y **Lisandro** implementó la apertura del simulador dual (`ADP-05` en `61b973d`) y el catálogo Debian 12 Bookworm (`ADP-06` en `aee232a` el 02/11).

3. **Absorción por Lisandro de las tareas de Seguridad y Cuotas asignadas a Lucas:**
   - Las tareas `FIX-133` (CPU overcommit y delta de RAM), `FIX-138` (contraseñas NIST SP 800-63B) y `FIX-140` (transparencia del proveedor SMTP mock), que figuraban en `documentacion/actual.md` asignadas a **Lucas**, fueron implementadas íntegramente por **Lisandro** el 02/11/2026 (`fb8a41c`, `3fe055b`, `4e4d2c4`).

4. **Detalles pendientes para el 100% en la suite de aceptación PVE 8.4 (`test/back`):**
   - Aunque la funcionalidad de adaptación a PVE 8.4 ya está codificada en `internal/infraestructura/proxmox/client.go`, la suite externa `pve8_adaptation_acceptance_test.go` aún marca 2 fallos puntuales de contrato:
     - **ADP-03:** Mover o exponer los tests de resiliencia PVE 8.4 en el archivo `internal/infraestructura/proxmox/client_pve8_test.go`.
     - **ADP-04:** Agregar soporte de rangos con guion (`100-105`) en la función `ParseVmidsProtegidos`.
   - Adicionalmente, restan implementar en `backend` los **5 fixes de endurecimiento y seguridad** (`FIX-134`, `FIX-135`, `FIX-136`, `FIX-137` y `FIX-139`).
