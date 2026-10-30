# Cruce de tareas, commits y autores reales — Submódulo `backend` (Parte B)

**Proyecto:** El Centinela · **Repositorio analizado:** `backend/` (`tayraag/centinela-back`)  
**Fecha del análisis:** 30/10/2026  
**Universo:** 44 commits (`43b011e` → `908522e`, 2026-10-19 a 2026-10-29) — *Continuación inmediata de la Parte A (`cd6a83e` → `8daaf87`)*  
**Fuentes:** `git log` del submódulo + `documentacion/terminado3.md`, `documentacion/terminado2.md`, `documentacion/terminado-1.md`, `documentacion/actual.md` y suites en `test/back`.

> **Objetivo:** Determinar con precisión qué integrante implementó realmente cada una de las tareas avanzadas de la Etapa 2 y Etapa 3 (Snapshots, Alertas, Hardening de Autenticación, Auditoría particionada con streaming y Parametrización dinámica), contrastando el **responsable asignado en la documentación/ClickUp** con el **autor real de los commits en el repositorio**.

---

## 1. Metodología y nivel de confianza

1. Se extrajo el historial completo del submódulo posterior al corte del 20/10/2026 (hash, autor, correo, fecha, asunto, topología de merges y archivos modificados).
2. Se extrajeron las tareas de las etapas finales (`BAC-36`, `BAC-42` a `BAC-48`, `BAC-54`, `FIX-107` a `FIX-119`) de `terminado2.md`, `terminado3.md` y `actual.md`.
3. Cada commit se vinculó a una o más tareas bajo los tres niveles probatorios estándar:

| Base | Significado | Uso |
|---|---|---|
| **Explícita** | El hash aparece citado en el informe o en la ficha de verificación del test. | Máxima confianza. |
| **Mensaje/feature** | El asunto del commit o PR nombra el requerimiento o la funcionalidad exacta (`BAC-xx`, `FIX-xx`, `RF-xx`). | Alta confianza. |
| **Inferida** | El asunto y los archivos tocados corresponden al módulo de la tarea, sin cita directa. | Media confianza. |

> [!IMPORTANT]
> Los **6 commits de merge** no aportan código propio: integran ramas de trabajo pero **no suman** a la autoría efectiva de la persona que mergea. La autoría se computa sobre los 38 commits con código, tests y documentación técnica real.

> [!NOTE]
> Nombres normalizados por persona: `Tayra` / `tayraag` → **Tayra**; `lisandroureta` → **Lisandro**; `NicoLencina` / `Lencina Nicolas Angelo` → **Nico**; `Lucas Vela` → **Lucas**.

---

## 2. Panorama general de autoría (Parte B)

| Autor | Commits totales | % total | Commits con código | % código real | Primer commit | Último commit |
|---|---:|---:|---:|---:|---|---|
| **Tayra** | 20 | 45.5% | 20 | 52.6% | 2026-10-19 | 2026-10-27 |
| **Lisandro** | 12 | 27.3% | 12 | 31.6% | 2026-10-20 | 2026-10-29 |
| **Nico** | 7 | 15.9% | 4 | 10.5% | 2026-10-19 | 2026-10-19 |
| **Lucas** | 5 | 11.4% | 2 | 5.3% | 2026-10-19 | 2026-10-19 |
| **Total** | **44** | **100%** | **38** | **100%** | 2026-10-19 | 2026-10-29 |

- **Commits con cambios reales:** 38 · **Merges:** 6 (3 Nico, 3 Lucas).
- **Commits vinculados a tareas:** 44/44 (100%).
- **Tareas analizadas en Parte B:** 22 tareas (todas recibieron implementación y están cerradas con tests en verde).

| Persona | Tareas asignadas en este bloque | Tareas en las que commiteó | Propias | Ajenas |
|---|---:|---:|---:|---:|
| **Tayra** | 10 | 11 | 4 | 7 |
| **Lisandro** | 13 | 11 | 6 | 5 |
| **Nico** | 1 | 4 | 0 | 4 |
| **Lucas** | 0 | 2 | 0 | 2 |

> *Nota: Una tarea co-asignada (p. ej. `Lisandro / Tayra`) cuenta como propia para ambas si commitearon en ella. Lucas y Nico no tenían asignaciones directas en estas tareas de backend pero aportaron fixes y soporte.*

---

## 3. Listado cronológico de los 44 commits (Parte B)

Numeración correlativa continuando los 218 commits de la Parte A (del #219 al #262).

| # | Fecha | Hash | Autor | Asunto | Tarea(s) relacionada(s) | Base |
|---:|---|---|---|---|---|---|
| 219 | 2026-10-19 | `43b011e` | Nico | fix(nodo): no publicar RAM ni almacenamiento en el histórico cuando falta el total | `FRN-19C`, `FIX-103` | Mensaje |
| 220 | 2026-10-19 | `9cdb76a` | Nico | Merge pull request #23 from tayraag/fix/nodo-historial-sin-ceros *(merge)* | `FRN-19C`, `FIX-103` | Inferida |
| 221 | 2026-10-19 | `49065dd` | Nico | fix(nodo): recuperar histórico de RAM y disco en Proxmox | `FRN-19C`, `FIX-103` | Mensaje |
| 222 | 2026-10-19 | `4240edf` | Lucas | Limpieza de permisos al eliminar una instancia (RF-03 / RF-09) | `BAC-24C`, `FIX-113` | Mensaje |
| 223 | 2026-10-19 | `26418a2` | Lucas | Merge branch 'main' into Limpieza-de-permisos-al-eliminar-una-instancia-(RF-03-/-RF-09) *(merge)* | `BAC-24C`, `FIX-113` | Inferida |
| 224 | 2026-10-19 | `846ad66` | Lucas | Merge pull request #24 from tayraag/Limpieza-de-permisos-al-eliminar-una-instancia-(RF-03-/-RF-09) *(merge)* | `BAC-24C`, `FIX-113` | Inferida |
| 225 | 2026-10-19 | `5b34a94` | Lucas | Se corrigió el audit port que quedo falaldo durante el merge | `BAC-18`, `BAC-41` | Mensaje |
| 226 | 2026-10-19 | `f8cd840` | Lucas | Merge pull request #25 from tayraag/Correción-de-audit_port *(merge)* | `BAC-18`, `BAC-41` | Inferida |
| 227 | 2026-10-19 | `6925281` | Nico | fix(auditoria): registrar advertencias operativas como éxito | `BAC-27`, `FIX-53`, `BAC-41` | Mensaje |
| 228 | 2026-10-19 | `99dff50` | Nico | Merge pull request #26 from tayraag/fix/acciones-instancia-contrato-20261019 *(merge)* | `BAC-27`, `FIX-53`, `BAC-41` | Inferida |
| 229 | 2026-10-19 | `9b3f32d` | Nico | fix(eventos): mantener vivo el stream SSE | `BAC-21`, `BAC-52` | Mensaje |
| 230 | 2026-10-19 | `248163d` | Nico | Merge pull request #27 from tayraag/fix/sse-latido-auditoria-20261019 *(merge)* | `BAC-21`, `BAC-52` | Inferida |
| 231 | 2026-10-19 | `dc50c0a` | Tayra | fix(cuenta): blindar el cambio de correo con reautenticación y confirmación | `FIX-111A` | Mensaje |
| 232 | 2026-10-19 | `526d0f2` | Tayra | test(cuenta): cubrir las protecciones del cambio de correo | `FIX-111A` | Mensaje |
| 233 | 2026-10-19 | `10e2c84` | Tayra | docs(api): actualizar Swagger del cambio seguro de correo | `FIX-111A` | Mensaje |
| 234 | 2026-10-19 | `b7cdaf0` | Tayra | fix(auth): Ampliación de cota máxima de contraseñas a 128 caracteres en Backend y DTOs | `FIX-115A` | Mensaje |
| 235 | 2026-10-19 | `638bb5e` | Tayra | fix(auth): revocar sesiones concurrentes al cambiar contraseña | `FIX-114` | Mensaje |
| 236 | 2026-10-19 | `a0f5465` | Tayra | fix(realtime): cerrar canales de sesiones revocadas | `FIX-114` | Mensaje |
| 237 | 2026-10-19 | `c99c9f5` | Tayra | docs(auth): documentar revocación tras cambio de contraseña | `FIX-114` | Mensaje |
| 238 | 2026-10-19 | `f738f1c` | Tayra | fix(auth): limitar intentos distribuidos en rutas críticas | `FIX-112`, `BAC-48` | Mensaje |
| 239 | 2026-10-19 | `ef53022` | Tayra | fix(auth): serializar intentos de recuperación | `FIX-112` | Mensaje |
| 240 | 2026-10-19 | `c6c97f2` | Tayra | docs(auth): documentar límites y modo seguro | `FIX-112` | Mensaje |
| 241 | 2026-10-20 | `c51c77d` | Lisandro | fix(instances): ofuscar VMID con referencias públicas contra IDOR y enumeración | `FIX-107` | Mensaje |
| 242 | 2026-10-20 | `1d84e7b` | Lisandro | fix(topology): enmascarar nodos y storages para usuarios no ADMIN | `FIX-109` | Mensaje |
| 243 | 2026-10-20 | `c9deb62` | Lisandro | fix(errors): sanear errores HTTP y eventos SSE de trazas de Proxmox | `FIX-110` | Mensaje |
| 244 | 2026-10-20 | `16a7ba2` | Tayra | feat(instances): servir metricas historicas con cache de 30 segundos | `BAC-36` | Mensaje |
| 245 | 2026-10-20 | `ad6565f` | Tayra | docs(swagger): documentar metricas historicas de instancia | `BAC-36` | Mensaje |
| 246 | 2026-10-20 | `5c72bbb` | Tayra | FIX: Validación de existencia de VMIDs y bloqueo de instancias protegidas en asignación de permisos (AsignarPermisos) | `FIX-119` | Mensaje |
| 247 | 2026-10-21 | `94e95a2` | Lisandro | fix(node): segregar telemetría del host por rol para operadores | `FIX-117A` | Mensaje |
| 248 | 2026-10-21 | `a25d1b3` | Lisandro | fix(instances): validar parámetros de creación contra inyección de directivas | `FIX-116` | Mensaje |
| 249 | 2026-10-21 | `1801e33` | Lisandro | fix(db): quitar credenciales por defecto del sembrado del administrador | `FIX-118` | Mensaje |
| 250 | 2026-10-23 | `c08af54` | Tayra | fix: purgar permisos al destruir instancias | `BAC-24C`, `FIX-113` | Mensaje |
| 251 | 2026-10-26 | `7c922c6` | Lisandro | feat(snapshots): listar y crear snapshots de instancias | `BAC-42` | Mensaje |
| 252 | 2026-10-27 | `051859b` | Lisandro | feat(snapshots): revertir y eliminar snapshots con exclusión mutua | `BAC-43` | Mensaje |
| 253 | 2026-10-27 | `5f98ff0` | Lisandro | feat(alertas): evaluador de umbrales con histéresis y deduplicación | `BAC-46` | Mensaje |
| 254 | 2026-10-27 | `9867324` | Lisandro | feat(alertas): despacho por correo y auditoría de alertas | `BAC-47` | Mensaje |
| 255 | 2026-10-27 | `71dd4f2` | Tayra | feat(auditoria): completar registro y consultas administrativas RF-08 | `BAC-44` | Mensaje |
| 256 | 2026-10-27 | `53e6d0c` | Tayra | test(auditoria): verificar trazabilidad, aislamiento y rendimiento RF-08 | `BAC-44` | Mensaje |
| 257 | 2026-10-27 | `a56af59` | Tayra | docs(auditoria): actualizar Swagger y guía de verificación RF-08 | `BAC-44` | Mensaje |
| 258 | 2026-10-27 | `273b617` | Tayra | feat(auditoria): transmitir exportaciones CSV y JSON por streaming | `BAC-45` | Mensaje |
| 259 | 2026-10-27 | `f6398a8` | Tayra | test(auditoria): verificar descargas masivas y cancelación | `BAC-45` | Mensaje |
| 260 | 2026-10-27 | `a3be3d9` | Tayra | docs(auditoria): actualizar contrato de exportación por streaming | `BAC-45` | Mensaje |
| 261 | 2026-10-29 | `651e8de` | Lisandro | feat(auth): cookie alineada al refresh, límite por IP y saneamiento de entradas | `BAC-48` | Mensaje |
| 262 | 2026-10-29 | `908522e` | Lisandro | feat(proxmox): nodo, VMIDs protegidos y versión de PVE desde el .env | `BAC-54` | Mensaje |

---

## 4. Cruce por tarea: Asignado declarado vs. Autor real

| Tarea | Título | Área | Asignado(s) declarado(s) | Autoría real (commits) | Veredicto |
|---|---|---|---|---|---|
| `BAC-24C` / `FIX-113` | Limpieza de permisos al eliminar una instancia | Backend | Tayra / Lisandro | Lucas (1), Tayra (1) | ⚠️ Parcial |
| `BAC-36` | Endpoint REST de métricas históricas de instancia (`timeframe=2h\|24h`) | Backend | Tayra | Tayra (2) | ✅ Coincide |
| `BAC-42` | Adaptador Proxmox Snapshots, creación y listado | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `BAC-43` | Reversión y eliminación de Snapshots con exclusión mutua | Backend | Tayra | Lisandro (1) | ❌ No coincide *(Lisandro)* |
| `BAC-44` | Interceptor global y consultas complejas sobre auditoría | Backend | Tayra | Tayra (3), Nico (1) | ⚠️ Parcial *(mayoría Tayra)* |
| `BAC-45` | Streaming de exportación estructurada de auditoría CSV/JSON | Backend | Lisandro | Tayra (3) | ❌ No coincide *(Tayra)* |
| `BAC-46` | Worker evaluador de umbrales con histéresis y deduplicación | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `BAC-47` | Despacho multicanal de alertas (SSE y correo estructurado) | Backend | Tayra | Lisandro (1) | ❌ No coincide *(Lisandro)* |
| `BAC-48` | Hardening de autenticación: Cookies HttpOnly, Rate Limiting y Sanitización | Backend | Tayra | Lisandro (1), Tayra (1) | ⚠️ Parcial *(Lisandro/Tayra)* |
| `BAC-54` | Parametrización en `.env` para Servidores Heterogéneos | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `FIX-107` | Cifrado y ofuscación de VMIDs contra IDOR | Backend | Tayra | Lisandro (1) | ❌ No coincide *(Lisandro)* |
| `FIX-109` | Enmascaramiento de host, nodo y pools para no ADMIN | Backend | Lisandro / Tayra | Lisandro (1) | ✅ Coincide |
| `FIX-110` | Saneamiento de errores y SSE de trazas de hipervisor | Backend | Tayra | Lisandro (1) | ❌ No coincide *(Lisandro)* |
| `FIX-111A` | Blindaje contra secuestro de cuenta en cambio de correo | Backend | Lisandro | Tayra (3) | ❌ No coincide *(Tayra)* |
| `FIX-112` | Middleware de Rate Limiting y fuerza bruta | Backend | Lisandro / Nico | Tayra (3) | ❌ No coincide *(Tayra)* |
| `FIX-114` | Revocación de sesiones activas en cambio de contraseña | Backend | Lisandro | Tayra (3) | ❌ No coincide *(Tayra)* |
| `FIX-115A` | Cota máxima de contraseñas a 128 caracteres | Backend | Lisandro | Tayra (1) | ❌ No coincide *(Tayra)* |
| `FIX-116` | Sanitización de directivas e inyección en Proxmox | Backend | Lisandro / Tayra | Lisandro (1) | ✅ Coincide |
| `FIX-117A` | Segregación de telemetría de host para operadores | Backend | Tayra | Lisandro (1) | ❌ No coincide *(Lisandro)* |
| `FIX-118` | Eliminación de credenciales por defecto en migraciones | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `FIX-119` | Validación de existencia y bloqueo de VMIDs protegidos | Backend | Lisandro / Tayra | Tayra (1) | ✅ Coincide |
| `FRN-19C` / `FIX-103` | Soporte backend del histórico del host sin ceros | Backend / Dashboard | Nico / Belinda | Nico (2) | ✅ Coincide *(Nico)* |

---

## 5. Desvíos detectados y hallazgos clave

### 5.1. El gran intercambio entre Lisandro y Tayra (Inversión de Roles)
En este bloque final se produjo un fenómeno sistemático de intercambio cruzado de módulos entre los dos desarrolladores principales del backend:

1. **Lisandro asumió tareas asignadas formalmente a Tayra:**
   - `BAC-43` (Reversión y Eliminación de Snapshots).
   - `BAC-47` (Despacho multicanal de alertas vía SSE y correo).
   - `BAC-48` (Hardening de autenticación y cookie HttpOnly en `651e8de`).
   - `FIX-107` (Ofuscación de VMIDs y referencias públicas).
   - `FIX-110` (Saneamiento de errores HTTP y trazas Proxmox).
   - `FIX-117A` (Segregación de telemetría del host por rol).

2. **Tayra asumió tareas asignadas formalmente a Lisandro:**
   - `BAC-45` (Streaming de exportación masiva de auditoría CSV/JSON).
   - `FIX-111A` (Blindaje de cambio de correo electrónico con reautenticación).
   - `FIX-112` (Protección contra ataques distribuidos de fuerza bruta y rate limiting).
   - `FIX-114` (Revocación obligatoria de sesiones concurrentes al cambiar contraseña).
   - `FIX-115A` (Ampliación de cotas de contraseñas a 128 caracteres).

### 5.2. Intervención de Lucas en `BAC-24C` / `FIX-113`
- La tarea de purga de permisos huérfanos al eliminar una instancia (`BAC-24C` / `FIX-113`) figuraba asignada a Tayra y Lisandro.
- Sin embargo, la implementación estructural de la purga fue realizada por **Lucas Vela** en el commit `4240edf`, resolviendo además conflictos en el puerto de auditoría (`5b34a94`). Posteriormente Tayra consolidó un fix de cierre en `c08af54`.

### 5.3. Intervención transversal de Nico
- Nico aportó correcciones determinantes en la telemetría del nodo (`43b011e`, `49065dd` para evitar promedios en cero que rompían los gráficos del frontend), en el feed de auditoría (`6925281`) y en el latido (keep-alive) del canal SSE (`9b3f32d`).

---

## 6. Conclusiones y balance final del backend

- En la Parte A se habían detectado **26 tareas de backend sin commits** (pendientes). En esta Parte B, **el 100% de esas tareas pendientes fue implementado, testeado y cerrado**.
- **Desvío de autoría:** de las 22 tareas analizadas, **8 coinciden totalmente**, **3 son parciales/compartidas** y **11 no coinciden con la asignación original (50% de desvío)**.
- El desvío no representó retraso ni falta de trabajo, sino una redistribución táctica: Tayra se especializó en el módulo de seguridad de cuentas y streaming de auditoría, mientras Lisandro resolvió la lógica de snapshots, el motor de alertas y la ofuscación de hipervisor.

---

## 7. Commits críticos de la Parte B

| Commit | Autor | Tarea | Impacto técnico |
|---|---|---|---|
| `c51c77d` | Lisandro | `FIX-107` | Implementa el token criptográfico opaco para reemplazar el VMID en URLs públicas, mitigando IDOR en todo el sistema. |
| `dc50c0a` | Tayra | `FIX-111A` | Exige contraseña actual y token OTP antes de permitir la actualización del correo de la cuenta. |
| `051859b` | Lisandro | `BAC-43` | Añade la guarda `INSTANCE_BUSY` (409) impidiendo alterar snapshots si hay tareas de energía en vuelo. |
| `273b617` | Tayra | `BAC-45` | Transmisión por chunks (`http.Flusher`) de auditoría con cursor SQL `.Rows()` sin saturar RAM. |
| `651e8de` | Lisandro | `BAC-48` | Emite el refreshToken exclusivamente como cookie HTTP segura con atributo `SameSite=Strict`. |
| `908522e` | Lisandro | `BAC-54` | Parametriza nodo y VMIDs protegidos en variables de entorno para compatibilidad multi-entorno PVE 8/9. |
