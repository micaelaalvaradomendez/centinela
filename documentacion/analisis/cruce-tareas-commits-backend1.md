# Cruce de tareas, commits y autores reales — Submódulo `backend`

**Proyecto:** El Centinela · **Repositorio analizado:** `backend/` (`tayraag/centinela-back`)  
**Fecha del análisis:** 20/10/2026  
**Universo:** 218 commits (`cd6a83e` → `8daaf87`, 2026-09-06 a 2026-10-18)  
**Fuentes:** `git log` del submódulo + `documentacion/actual.md`, `terminado.md`, `terminado-1.md`, `terminado2.md`, `etapa1..3.md` y `test/informe.md`.

> **Objetivo:** determinar quién escribió realmente cada tarea, comparando al **asignado declarado en la documentación/ClickUp** con el **autor de los commits** que implementan esa tarea.

---

## 1. Metodología y nivel de confianza

1. Se extrajo el historial completo del submódulo (hash, autor, correo, fecha, asunto y archivos tocados).
2. Se extrajeron las tareas y su campo **Asignado/Asignada** de la documentación (espejo del tablero ClickUp).
3. Cada commit se vinculó a una o más tareas con tres niveles de evidencia:

| Base | Significado | Uso |
|---|---|---|
| **Explícita** | El hash aparece citado en el informe o en la ficha de la tarea. | Máxima confianza. |
| **Mensaje/feature** | El asunto del commit nombra el requerimiento o la funcionalidad exacta (`RF-xx`, `BAC-xx`, `FIX-xx`, etc.). | Alta confianza. |
| **Inferida** | El asunto y los archivos tocados corresponden al módulo de la tarea, pero no hay cita directa. | Media; revisar si se usa para evaluar desempeño. |

> [!IMPORTANT]
> Los commits de **merge** no aportan código propio: se listan y se asocian a las tareas de la rama integrada, pero **no cuentan** para la autoría efectiva.

> [!NOTE]
> Nombres normalizados por persona: `tayraag`/`Tayra` → **Tayra**; `lisandroureta`/`Lisandro Ureta` → **Lisandro**; `NicoLencina`/`Lencina Nicolas Angelo` → **Nico**; `Lucas Vela` → **Lucas**.

---

## 2. Panorama general de autoría

| Autor | Commits | % del total | Primer commit | Último commit |
|---|---:|---:|---|---|
| **Tayra** | 79 | 36.2% | 2026-09-09 | 2026-10-18 |
| **Lisandro** | 76 | 34.9% | 2026-09-06 | 2026-10-16 |
| **Nico** | 45 | 20.6% | 2026-09-13 | 2026-10-18 |
| **Lucas** | 18 | 8.3% | 2026-09-21 | 2026-10-12 |
| **Total** | **218** | **100%** | 2026-09-06 | 2026-10-18 |

- **Commits con cambios reales:** 182 · **Merges:** 36
- **Commits vinculados a al menos una tarea:** 199/218 (91.3%).
- **Commits sin tarea asociada** (bootstrap, CI, documentación o merges): 19.
- **Vinculación explícita:** 54; el resto se resolvió por mensaje/feature o por inferencia de módulo.

| Persona | Tareas backend asignadas | Tareas en las que commiteó | Propias | Ajenas |
|---|---:|---:|---:|---:|
| **Tayra** | 60 | 48 | 26 | 22 |
| **Lisandro** | 68 | 47 | 26 | 21 |
| **Nico** | 10 | 21 | 4 | 17 |
| **Lucas** | 7 | 11 | 1 | 10 |

> _Tareas “propias” = asignadas a la persona y con al menos un commit suyo. “Ajenas” = ejecutadas sin estar asignada. Una tarea con dos asignados cuenta para ambos. No se listan Cristian, Belinda ni Luz porque no tienen commits en `backend/`. Incluye tareas cross-area (ver 5.5)._

---

## 3. Listado y análisis de los 218 commits

Orden cronológico (del más antiguo al más reciente). La columna **Base** indica cómo se estableció la relación con la tarea.

| # | Fecha | Hash | Autor | Asunto | Tarea(s) relacionada(s) | Base |
|---:|---|---|---|---|---|---|
| 1 | 2026-09-06 | `cd6a83e` | Lisandro | Initial commit: estructura base del backend y documentacion | — | Inferida |
| 2 | 2026-09-06 | `ebee84e` | Lisandro | Agrego md a todas las carpetas | — | Inferida |
| 3 | 2026-09-06 | `b67468c` | Lisandro | feat: agregar persistencia postgres con docker, modelos de dominio y migraciones gorm | `BAC-01`, `BAC-02`, `BAC-05` | Inferida |
| 4 | 2026-09-09 | `aaaf3d3` | Tayra | primer intento de login y auth | `BAC-03`, `LOGIN-01` | Mensaje |
| 5 | 2026-09-12 | `807a523` | Lisandro | feat(auth): implementar autenticación, 2FA (TOTP), JWT y control de acceso (RF-01) | `BAC-02`, `BAC-03`, `LOGIN-01`, `BAC-10` | Mensaje |
| 6 | 2026-09-13 | `538b494` | Lisandro | feat: implemento gestion de usuarios | `BAC-05`, `BAC-06`, `BAC-09` | Mensaje |
| 7 | 2026-09-13 | `a0d6b6c` | Nico | ci: add deploy-back workflow + deploy script | `INF-03` | Inferida |
| 8 | 2026-09-13 | `8eac6f8` | Nico | ci: set GOPATH for self-hosted runner build | `INF-03` | Inferida |
| 9 | 2026-09-13 | `513b67c` | Nico | ci: declare GOCACHE for self-hosted runner | `INF-03` | Inferida |
| 10 | 2026-09-13 | `ac83887` | Nico | ci: atomic mv for running binary deploy | `INF-03` | Inferida |
| 11 | 2026-09-13 | `bee5044` | Nico | ci: wait for service active with retry loop | `INF-03` | Inferida |
| 12 | 2026-09-13 | `f69a6f4` | Lisandro | entra swagger epicamente | `BAC-09` | Inferida |
| 13 | 2026-09-14 | `92b7dcb` | Lisandro | Revert "primer intento de login y auth" | `BAC-03` | Mensaje |
| 14 | 2026-09-14 | `6beb9a2` | Lisandro | Merge branch 'login' _(merge)_ | `BAC-03` | Mensaje |
| 15 | 2026-09-14 | `8d7bce5` | Lisandro | feat: estandarizacion de errores HTTP y refactorizacion de entorno | `BAC-04` | Mensaje |
| 16 | 2026-09-15 | `2e8b02a` | Nico | fix(db): aceptar DATABASE_URL como alternativa a DB_DSN | `BAC-01`, `INF-03` | Inferida |
| 17 | 2026-09-15 | `a30bf35` | Nico | Merge pull request #1 from tayraag/fix/db-dsn-compatibilidad _(merge)_ | `INF-03` | Inferida |
| 18 | 2026-09-15 | `1b9456c` | Lisandro | se corrigen relaciones de las entidades y nombres de las tablas | `BAC-05` | Inferida |
| 19 | 2026-09-15 | `4919a88` | Lisandro | Merge branch 'main' of https://github.com/tayraag/centinela-back _(merge)_ | — | Inferida |
| 20 | 2026-09-15 | `4e8e9c4` | Lisandro | fix: ajustar JSON de login y definir foreign keys en modelos | `BAC-03`, `BAC-05` | Inferida |
| 21 | 2026-09-15 | `e70701d` | Nico | fix(auth): corregir mensaje de login y agregar endpoint de logout | `BAC-03`, `BAC-17` | Mensaje |
| 22 | 2026-09-15 | `7debb72` | Nico | infra: GET /api/version para verificar deploy + deploy solo desde main | `INF-03` | Inferida |
| 23 | 2026-09-15 | `ca15ad9` | Nico | Merge pull request #2 from tayraag/infra/verificacion-deploy-version _(merge)_ | `INF-03` | Inferida |
| 24 | 2026-09-15 | `aea6f63` | Nico | infra: comentario de prueba para verificar el deploy del back | `INF-03` | Inferida |
| 25 | 2026-09-15 | `4536776` | Nico | Merge pull request #3 from tayraag/infra/test-comentario-deploy _(merge)_ | `INF-03` | Inferida |
| 26 | 2026-09-17 | `31e9bd9` | Lisandro | fix(swagger): remove hardcoded localhost host to allow dynamic remote access | — | Inferida |
| 27 | 2026-09-17 | `0955134` | Nico | fix(swagger): quitar @host para que Try it out apunte al servidor | — | Inferida |
| 28 | 2026-09-17 | `8424bfd` | Nico | Merge pull request #4 from tayraag/fix/swagger-host _(merge)_ | — | Inferida |
| 29 | 2026-09-17 | `a37160c` | Lisandro | chore: format files and save local changes | — | Inferida |
| 30 | 2026-09-17 | `0d3e9ca` | Lisandro | Merge branch 'main' - resolve conflict in swagger annotations _(merge)_ | — | Inferida |
| 31 | 2026-09-17 | `addb9b6` | Lisandro | fix(auth): remove duplicate CerrarSesion implementations introduced during merge | `BAC-17` | Inferida |
| 32 | 2026-09-17 | `cd04708` | Lisandro | fix(auth): remove dangling syntax in auth_handler.go | `BAC-17` | Inferida |
| 33 | 2026-09-17 | `dbe4079` | Lisandro | limpieza: eliminar endpoint redundante de cierre de sesion | `BAC-17` | Inferida |
| 34 | 2026-09-17 | `9b12618` | Lisandro | fix: eliminar import sin uso en account_handler.go | — | Inferida |
| 35 | 2026-09-17 | `0a5faad` | Lisandro | feat(auth): normalizar correos electronicos a minusculas para login case-insensitive | `BAC-03` | Mensaje |
| 36 | 2026-09-17 | `caf71c2` | Tayra | actulizacion del README | — | Inferida |
| 37 | 2026-09-18 | `3052a02` | Lisandro | fix(users): evitar permisos de instancia duplicados con restriccion unica y deduplicacion | `BAC-07` | Inferida |
| 38 | 2026-09-18 | `adcf261` | Lisandro | fix(users): montar CRUD de usuarios bajo /admin/users y corregir endpoints en Swagger | `BAC-06` | Inferida |
| 39 | 2026-09-19 | `52ebce1` | Lisandro | chore(docs): actualizar anotaciones de Swagger pendientes para /admin/users | `BAC-06` | Inferida |
| 40 | 2026-09-20 | `924641f` | Tayra | fix: rechazar generacion de QR si el 2FA ya esta vinculado (409 TOTP_ALREADY_LINKED) | `BAC-10`, `FIX-03` | Mensaje |
| 41 | 2026-09-20 | `9df22e6` | Tayra | feat: implementar proteccion anti-replay para codigos TOTP en ventana de 30s | `BAC-11` | Mensaje |
| 42 | 2026-09-21 | `019017b` | Tayra | refactor(auth): eliminar endpoint legacy relink y estandarizar reset administrativo | `BAC-10`, `BAC-13`, `BAC-15` | Inferida |
| 43 | 2026-09-21 | `a4c5215` | Lisandro | chore(docs): corregir SETUP.md para que coincida con el backend actual | — | Inferida |
| 44 | 2026-09-21 | `de91ca3` | Lisandro | fix(domain): eliminar colision del indice idx_usuario_vmid entre modelos | `BAC-05`, `BAC-07`, `FIX-01` | Inferida |
| 45 | 2026-09-21 | `680cc13` | Tayra | fix(auth): corregir código de error para protección de revinculación 2FA | `FIX-03` | Inferida |
| 46 | 2026-09-21 | `d192003` | Tayra | feat(auth): implementar guard y validacion de cambio obligatorio de contrasena | `BAC-12` | Mensaje |
| 47 | 2026-09-21 | `eb84a58` | Lisandro | feat: middleware de autorizacion por recurso (instance guard) | `BAC-08` | Mensaje |
| 48 | 2026-09-21 | `d11426e` | Lisandro | docs: actualizar README con instance_guard, instance_port e instance_repository | `BAC-08` | Inferida |
| 49 | 2026-09-21 | `5a7abe1` | Tayra | docs(admin): actualizar swagger del restablecimiento administrativo de contrasena | `BAC-15` | Inferida |
| 50 | 2026-09-21 | `e363afa` | Tayra | feat: implementar recuperacion de contraseña (RF-13) y mailer service simulado | `BAC-16`, `BAC-19`, `BAC-20` | Mensaje |
| 51 | 2026-09-21 | `c89ddcb` | Tayra | docs: actualizar swagger con nuevos endpoints de auth | `BAC-19` | Inferida |
| 52 | 2026-09-21 | `0e68371` | Tayra | fix: aplicar seguridad estricta crypto y mitigar timing attacks en recuperacion | `BAC-19` | Inferida |
| 53 | 2026-09-21 | `d6fdcb1` | Lucas | Se termina de hacer el requerimiento de la auditoria con registros básicos, el esqueleto, aun … | `BAC-18` | Mensaje |
| 54 | 2026-09-21 | `d15a8fe` | Lucas | Se crean test para probar auditoria (creacion y recuperación de un evento) | `BAC-18` | Inferida |
| 55 | 2026-09-21 | `732e65d` | Tayra | feat(auth): implementar revocacion global de JWTs y alinear diseño endpoints | `BAC-17` | Mensaje |
| 56 | 2026-09-21 | `e405566` | Lucas | Merge pull request #5 from tayraag/Base-transversal-de-auditoría _(merge)_ | `BAC-18` | Inferida |
| 57 | 2026-09-21 | `348dd59` | Tayra | feat(auth): implementar limites OTP, validar password historico y sync docs | `BAC-11`, `BAC-12` | Mensaje |
| 58 | 2026-09-21 | `d509d18` | Nico | 86e3bpfw3 - Contrato base del canal de eventos (RF-11) | `BAC-21` | Mensaje |
| 59 | 2026-09-21 | `6c75069` | Nico | Merge pull request #6 from tayraag/Nico _(merge)_ | `BAC-21` | Inferida |
| 60 | 2026-09-21 | `fedd4ba` | Tayra | Merge branch 'main' of https://github.com/tayraag/centinela-back into feat/gestion-credenciales _(merge)_ | `BAC-15`, `BAC-19` | Inferida |
| 61 | 2026-09-21 | `3dadc48` | Tayra | Resolucion de conflictos y correccion de build tras el merge _(merge)_ | — | Inferida |
| 62 | 2026-09-21 | `7a2403c` | Lisandro | feat(users): agregar GET de permisos y cambiar rutas de asignación de instancias a /permissions | `BAC-07` | Inferida |
| 63 | 2026-09-22 | `fea91a8` | Nico | fix(api): arreglo el 404 y el 405 para que respeten el contrato de errores | `BAC-04` | Inferida |
| 64 | 2026-09-22 | `5d8a7fd` | Nico | fix(auth): arreglo que el detalle interno del refresh token llegue al cliente | `SEC-01`, `BAC-17` | Inferida |
| 65 | 2026-09-22 | `3d06d7c` | Nico | fix(auditoria): arreglo que los registros sin usuario se vieran desde cualquier organizacion | `BAC-18` | Inferida |
| 66 | 2026-09-22 | `0e97f93` | Nico | fix(api): arreglo el swagger expuesto y agrego los headers de seguridad | `BAC-04` | Inferida |
| 67 | 2026-09-22 | `272a16e` | Nico | Merge remote-tracking branch 'origin/main' into Nico _(merge)_ | — | Inferida |
| 68 | 2026-09-22 | `e4a709d` | Nico | Merge pull request #7 from tayraag/Nico _(merge)_ | `BAC-04`, `BAC-18` | Inferida |
| 69 | 2026-09-22 | `7cadbce` | Lisandro | Merge branch 'feature/user-resource-authorization' _(merge)_ | `BAC-08` | Inferida |
| 70 | 2026-09-22 | `fd1f040` | Lisandro | fix(api): montar rutas de instancias con el guard de autorizacion por recurso | `FIX-16` | Inferida |
| 71 | 2026-09-22 | `f3c659a` | Lisandro | feat(proxmox): implementar cliente real de Proxmox VE | `BAC-14` | Inferida |
| 72 | 2026-09-22 | `60a6403` | Lisandro | feat(instances): agregar lectura minima del inventario de Proxmox | `BAC-14` | Mensaje |
| 73 | 2026-09-22 | `15032da` | Tayra | fix: asegurar logout y revocación atómica de sesiones | `BAC-17` | Explícita |
| 74 | 2026-09-23 | `5407599` | Tayra | FIX: Alinear y formalizar endpoints y contratos de contraseñas | `FIX-17`, `FIX-13` | Explícita |
| 75 | 2026-09-23 | `fec42b7` | Tayra | fix(auth): registrar fecha de ultimo acceso en login exitoso | `BAC-03` | Inferida |
| 76 | 2026-09-23 | `c719c9b` | Tayra | feat(auth): migrar refreshToken a cookie HttpOnly para mitigar XSS | `SEC-01` | Explícita |
| 77 | 2026-09-23 | `47bb8f2` | Tayra | fix(auth): corregir metodo SetSameSite de gin | `SEC-01` | Explícita |
| 78 | 2026-09-24 | `eb0c9af` | Lisandro | feat(instances): agregar niveles de acceso FULL_ACCESS/READ_ONLY a permisos_instancia | `SEC-04` | Explícita |
| 79 | 2026-09-26 | `c624a57` | Lucas | Implementado Cors middleware junto a sus tests, además de reiniciada la api del proxmox luego … | `INF-05` | Explícita |
| 80 | 2026-09-26 | `0cbf706` | Lucas | Merge branch 'main' into Base-transversal-de-auditoría _(merge)_ | `BAC-18` | Inferida |
| 81 | 2026-09-26 | `40e7a55` | Lucas | Merge pull request #8 from tayraag/Base-transversal-de-auditoría _(merge)_ | `BAC-18` | Inferida |
| 82 | 2026-09-27 | `2bdc3b1` | Nico | fix(auth): normalizar el email en recuperación de contraseña (mayúsculas/espacios rompían forg… | `BAC-19` | Inferida |
| 83 | 2026-09-27 | `44feb08` | Lucas | soporte de api token proxmox, filtrado rbac y tests de integracion | `BAC-08`, `BAC-14`, `INF-07` | Mensaje |
| 84 | 2026-09-27 | `c16bbee` | Lucas | Merge branch 'main' of https://github.com/tayraag/centinela-back _(merge)_ | — | Inferida |
| 85 | 2026-09-27 | `2a943cd` | Lucas | Agregado redis al docker compose | `BAC-17A`, `INF-06A` | Inferida |
| 86 | 2026-09-27 | `d36bc50` | Lucas | Merge pull request #9 from tayraag/Agregar-redis-a-docker _(merge)_ | `BAC-17A`, `INF-06A` | Inferida |
| 87 | 2026-09-28 | `ad52485` | Tayra | feat(seguridad): garantizar inmutabilidad física y legal de la auditoría (RF-08) | `BAC-18`, `FIX-23` | Explícita |
| 88 | 2026-09-28 | `9554efa` | Tayra | fix: arrreglo funcion de inmutabilidad en db | `FIX-23` | Explícita |
| 89 | 2026-09-29 | `725d436` | Lisandro | fix(instances): robustecer cliente Proxmox y proteger VMIDs de infraestructura | `BAC-14` | Inferida |
| 90 | 2026-09-29 | `82add96` | Lisandro | docs(proxmox): usar configuración oficial de la PM en .env.example y docs | `INF-07` | Inferida |
| 91 | 2026-09-29 | `de407a5` | Lisandro | feat(proxmox): simulador de Proxmox para desarrollo local, 504 en timeouts y logger sin datos … | `BAC-28`, `FIX-32`, `FIX-33` | Explícita |
| 92 | 2026-09-29 | `e6e7dc5` | Lisandro | chore(env): agregar ENABLE_SWAGGER a .env.example | — | Inferida |
| 93 | 2026-09-30 | `9e0f43c` | Lisandro | refactor(auth): 1 sesión = 1 fila en sesiones_activas y sesiones efímeras en Redis con TTL | `BAC-17B` | Explícita |
| 94 | 2026-09-30 | `546e1d5` | Lisandro | chore(redis): unificar credenciales de Redis local y agregar REDIS_DB | `BAC-17A` | Explícita |
| 95 | 2026-09-30 | `0aac034` | Lisandro | feat(redis): puerto KeyValueStore con adaptador Redis (clave-valor y Pub/Sub) | `BAC-17A` | Explícita |
| 96 | 2026-09-30 | `233e804` | Tayra | chore(config): definir variables de entorno para SMTP en env.example | `INF-08B` | Explícita |
| 97 | 2026-09-30 | `45ecad7` | Lisandro | feat(eventos): stream SSE en /api/events con ticket efímero, bus Pub/Sub y corte en vivo | `BAC-21`, `BAC-29` | Mensaje |
| 98 | 2026-09-30 | `c997398` | Lisandro | corrigiendo el simulador de proxmox | `FIX-32` | Explícita |
| 99 | 2026-09-30 | `2c56bab` | Lisandro | enter simulator | `FIX-32` | Explícita |
| 100 | 2026-09-30 | `43a0b06` | Lisandro | dejando ready el simulador de proxmox | `FIX-32` | Explícita |
| 101 | 2026-09-30 | `860b3c9` | Tayra | feat(email): adaptador SmtpEmailService para envio real de credenciales y codigos OTP | `BAC-16`, `BAC-16B`, `INF-08B` | Explícita |
| 102 | 2026-10-01 | `44a2339` | Lisandro | mas arreglos del simulador | `FIX-32` | Inferida |
| 103 | 2026-10-02 | `0167b96` | Lisandro | fix: distinguiendo timeout de proxmox | `FIX-40` | Explícita |
| 104 | 2026-10-02 | `96106a6` | Tayra | feat(db): indices parciales, purga de sesiones y particionamiento de auditoria | `BAC-18B` | Explícita |
| 105 | 2026-10-02 | `b8d2631` | Tayra | feat: Alineación de esquema, auditoría y operaciones de ciclo de vida de instancias | `FIX-39`, `FIX-53`, `BAC-18B`, `BAC-21B` | Explícita |
| 106 | 2026-10-02 | `573f0c3` | Tayra | docs: Actualizar documento de diseño de front para reflejar endpoints de auth, sse y auditoria… | — | Inferida |
| 107 | 2026-10-03 | `42e94ca` | Tayra | fix: Consolidar DTO de instancias, agregar soporte para shutdown/reboot y auditar acciones en … | `FIX-39`, `BAC-21B` | Explícita |
| 108 | 2026-10-03 | `4e204f1` | Tayra | fix: Informar el nivel de acceso por instancia en GET /account/profile y actualizar docs | `FIX-37` | Explícita |
| 109 | 2026-10-03 | `5787179` | Nico | docs: node status and instance inventory schemas (T01) | `BAC-29`, `FIX-46` | Explícita |
| 110 | 2026-10-03 | `21b1332` | Nico | docs: action schemas and explicit error catalog (T02) | `BAC-29`, `FIX-46` | Explícita |
| 111 | 2026-10-03 | `4ef36e3` | Nico | docs: SSE payload schemas and verification (T03) | `BAC-29`, `FIX-46` | Explícita |
| 112 | 2026-10-04 | `1b0486e` | Nico | Merge pull request #11 from tayraag/docs/http-events-contract-stage1 _(merge)_ | `BAC-29` | Inferida |
| 113 | 2026-10-04 | `8ba0c59` | Nico | docs: agregar documento de contrato en formato docx para OnlyOffice | `BAC-29` | Explícita |
| 114 | 2026-10-04 | `aade54c` | Nico | Merge pull request #12 from tayraag/docs/http-events-contract-stage1 _(merge)_ | `BAC-29` | Inferida |
| 115 | 2026-10-04 | `8591e90` | Lucas | Validación de estado previo y VMIDs protegidos en acciones de energía | `BAC-24A` | Explícita |
| 116 | 2026-10-04 | `e1f2df4` | Lucas | Merge pull request #13 from tayraag/Validación-de-estado-previo-y-VMIDs-protegidos-en-acciones… _(merge)_ | `BAC-24A` | Inferida |
| 117 | 2026-10-05 | `13f9c35` | Lisandro | feat(nodo): GET /api/node/status con caché en Redis y lectura stale | `BAC-22` | Explícita |
| 118 | 2026-10-05 | `3988546` | Lisandro | feat(tareas): worker pool acotado para el seguimiento de UPID, reconciliador y apagado ordenado | `BAC-25A` | Explícita |
| 119 | 2026-10-05 | `d1dec4f` | Lisandro | feat(inventario): consolidación de VMs y LXC con resolución concurrente de IPs | `BAC-23A` | Explícita |
| 120 | 2026-10-05 | `0850ec7` | Lisandro | fix(contrato): activeTask como objeto, DELETE asíncrono y TASK_FINISHED completo | `BAC-29`, `FIX-46`, `FIX-49` | Explícita |
| 121 | 2026-10-05 | `f17295c` | Lisandro | feat(tareas): recuperación de tareas en curso al arrancar con ventana de 3 minutos | `BAC-25C` | Mensaje |
| 122 | 2026-10-05 | `b476636` | Nico | fix(auditoria): unificar códigos de acción/resultado y preservar metadatos en acciones de inst… | `FIX-53` | Explícita |
| 123 | 2026-10-05 | `eec77ff` | Nico | Merge pull request #14 from tayraag/fix/86e3k5fqp-codigos-auditoria _(merge)_ | `FIX-53` | Inferida |
| 124 | 2026-10-06 | `c02fd29` | Tayra | FIX: Corregir índice parcial en sesiones_activas y alinear worker de purga | `FIX-44` | Explícita |
| 125 | 2026-10-06 | `088e175` | Tayra | feat: implementar endpoint destructivo DELETE /api/instances/:vmid | `BAC-24B` | Explícita |
| 126 | 2026-10-06 | `f06640f` | Tayra | feat: implementar endpoint destructivo DELETE /api/instances/:vmid | `BAC-24B`, `FIX-49` | Explícita |
| 127 | 2026-10-06 | `59148b1` | Tayra | Merge branch 'feat/delete-instance' of https://github.com/tayraag/centinela-back into feat/del… _(merge)_ | `BAC-24B`, `FIX-49` | Inferida |
| 128 | 2026-10-06 | `a25b7e6` | Nico | docs(contrato): explicito el control de acceso D1 de /node/status (86e3m2krc) | `FIX-65` | Mensaje |
| 129 | 2026-10-06 | `973ee0b` | Nico | test(proxmox): cubro ObtenerInterfaces con tests unitarios (86e3m2hgg) | `FIX-63` | Mensaje |
| 130 | 2026-10-06 | `9135553` | Nico | Merge pull request #15 from tayraag/fix/86e3m2krc-86e3m2hgg _(merge)_ | `FIX-63`, `FIX-65` | Inferida |
| 131 | 2026-10-06 | `df320c7` | Tayra | fix: recuperar y consolidar historial de auditoria en tabla particionada | `FIX-51` | Explícita |
| 132 | 2026-10-06 | `464b494` | Nico | fix(contrato): restauro los esquemas canónicos http.* y el path pause en Swagger (86e3m2p6y) | `FIX-69` | Inferida |
| 133 | 2026-10-06 | `dc14312` | Tayra | fix: implementar graceful shutdown en servidor HTTP | `FIX-50` | Mensaje |
| 134 | 2026-10-06 | `b19d9b9` | Nico | Merge pull request #16 from tayraag/fix/86e3m2p6y-swagger-canonico _(merge)_ | `FIX-69` | Inferida |
| 135 | 2026-10-06 | `1a6015e` | Nico | fix(tareas): tratar los WARNINGS de Proxmox como tarea completada | `BAC-25C`, `FIX-73` | Explícita |
| 136 | 2026-10-06 | `b630b2f` | Nico | Merge pull request #17 from tayraag/Nico _(merge)_ | `BAC-25C` | Inferida |
| 137 | 2026-10-08 | `08ca9dd` | Tayra | fix: blindaje transaccional de migración de auditoría y saneamiento de particiones | `FIX-70` | Explícita |
| 138 | 2026-10-08 | `eca422d` | Tayra | fix: heurística de circuit breaker y recuperación en nodo_service | `FIX-62` | Explícita |
| 139 | 2026-10-08 | `08762f4` | Lucas | FIX: Distinguir INSTANCE_BUSY de INSTANCE_INVALID_STATE y limpiar código muerto | `FIX-52` | Mensaje |
| 140 | 2026-10-08 | `222760d` | Lucas | Merge pull request #18 from tayraag/FIX--Distinguir-INSTANCE_BUSY-de-INSTANCE_INVALID_STATE-y-… _(merge)_ | `FIX-52` | Inferida |
| 141 | 2026-10-09 | `8f46c3d` | Tayra | feat: conteo de instancias (RF-02) y metricas de recursos (RF-03) | `BAC-22`, `BAC-29` | Mensaje |
| 142 | 2026-10-09 | `4f68e44` | Tayra | feat: implementar timeout configurable y backoff exponencial (RF-04) | `BAC-25B` | Explícita |
| 143 | 2026-10-09 | `c108594` | Lisandro | test(tareas): adaptar tests del seguimiento al OnTaskFinished | `FIX-72` | Inferida |
| 144 | 2026-10-09 | `a9b9d0f` | Lisandro | fix(usuarios): separar eliminación lógica de suspensión y unicidad parcial del correo | `FIX-54` | Mensaje |
| 145 | 2026-10-09 | `c024ba3` | Lisandro | cosas de usuarios | `FIX-55` | Inferida |
| 146 | 2026-10-10 | `c7fd7a1` | Tayra | fix: Robustecer reconexión y reintentos ante arranque en frío de Redis | `FIX-71` | Mensaje |
| 147 | 2026-10-10 | `ec84122` | Tayra | fix: Sincronizar firma de NewSeguimientoTareas en pruebas unitarias de servicios | `FIX-72` | Inferida |
| 148 | 2026-10-10 | `6bd1cf3` | Lucas | FIX: Documentar INSTANCE_INVALID_STATE en Swagger | `FIX-48` | Mensaje |
| 149 | 2026-10-10 | `2314d38` | Lucas | FIX: Documentar INSTANCE_INVALID_STATE en Swagger _(merge)_ | `FIX-48` | Inferida |
| 150 | 2026-10-10 | `d76eb35` | Tayra | Persistencia y seguimiento de tareas RF-08 | `BAC-27` | Mensaje |
| 151 | 2026-10-10 | `dc018a2` | Tayra | API de auditoría, filtros y actualizacion de swagger | `BAC-27` | Mensaje |
| 152 | 2026-10-10 | `7acbb56` | Tayra | Arreglo de tests | `BAC-26`, `BAC-27` | Inferida |
| 153 | 2026-10-10 | `a5e2110` | Tayra | Actualizacion de documentacion | — | Inferida |
| 154 | 2026-10-11 | `3ecd1c9` | Tayra | Merge branch 'main' of https://github.com/tayraag/centinela-back into feat/rf-08-auditoria-cic… _(merge)_ | `BAC-27` | Inferida |
| 155 | 2026-10-12 | `58388fc` | Lucas | Configuración de Nginx para WebSockets y Streaming bidireccional (RNF-06) | `INF-09` | Mensaje |
| 156 | 2026-10-12 | `6afc017` | Lucas | Merge pull request #20 from tayraag/Configuración-de-Nginx-para-WebSockets-y-Streaming-bidirec… _(merge)_ | `INF-09` | Inferida |
| 157 | 2026-10-13 | `8d6ce85` | Lisandro | feat(simulador): agregar endpoints de storage y contenidos para RF-05 y RF-07 | `BAC-30`, `FIX-84` | Mensaje |
| 158 | 2026-10-13 | `2d2495f` | Lisandro | feat(api): implementar endpoints de inventario de storages y templates (rf-07) | `BAC-31` | Mensaje |
| 159 | 2026-10-13 | `b307a63` | Tayra | fix: Evitar que `EMAIL_PROVIDER` caiga en silencio en el correo simulado | `FIX-88` | Mensaje |
| 160 | 2026-10-13 | `0022e18` | Tayra | feat(tareas): reconocer aprovisionamiento en worker | `BAC-32` | Mensaje |
| 161 | 2026-10-13 | `e37ed4d` | Tayra | feat(permisos): autoasignar acceso al completar creación | `BAC-32`, `BAC-39` | Mensaje |
| 162 | 2026-10-13 | `882c828` | Tayra | fix(eventos): normalizar accion CREATE | `BAC-39` | Mensaje |
| 163 | 2026-10-13 | `9891378` | Tayra | fix(permisos): asignar acceso en creaciones con warnings | `BAC-39` | Mensaje |
| 164 | 2026-10-13 | `468e452` | Tayra | feat(api): conectar aprovisionamiento con seguimiento | `BAC-38`, `BAC-39` | Mensaje |
| 165 | 2026-10-14 | `d1b3aeb` | Tayra | feat(aprovisionamiento): coordinar capacidad con redis | `BAC-37` | Mensaje |
| 166 | 2026-10-14 | `a363c9a` | Tayra | docs(api): documentar validacion de capacidad | `BAC-37` | Inferida |
| 167 | 2026-10-14 | `f2215af` | Lisandro | test(eventos): verificar TASK_FINISHED para 5 acciones canonicas, fallos y matriz rbac (rf-11) | `BAC-26` | Mensaje |
| 168 | 2026-10-14 | `9ab133e` | Lisandro | fix(simulador): completar storages y catalogo de contenidos canonicos (rf-07) | `BAC-30`, `FIX-84` | Mensaje |
| 169 | 2026-10-14 | `208907d` | Lisandro | docs(simulador): documentar endpoints de storage y catalogo de contenido | `FIX-84` | Inferida |
| 170 | 2026-10-14 | `d57bbe7` | Lisandro | feat(proxmox): implementar adaptador RRDdata y DTO PuntoMetricaDTO (rf-05) | `BAC-34` | Mensaje |
| 171 | 2026-10-14 | `394b25b` | Lisandro | entra WebSockets | `BAC-35` | Mensaje |
| 172 | 2026-10-14 | `acec121` | Tayra | feat(node): exponer cuotas asignables seguras | `BAC-33`, `BAC-37` | Mensaje |
| 173 | 2026-10-14 | `a97dbf0` | Tayra | docs(api): documentar cuotas asignables del nodo | `BAC-33` | Inferida |
| 174 | 2026-10-14 | `993c86f` | Lisandro | docs(ws): documentar endpoints de websocket y ws-ticket en swagger y redis (rf-05, rnf-06) | `BAC-35` | Inferida |
| 175 | 2026-10-14 | `d4ee23b` | Tayra | feat(instances): validar cuotas con bloqueo distribuido | `BAC-37` | Mensaje |
| 176 | 2026-10-14 | `ff92dd3` | Tayra | docs(api): documentar validacion distribuida de cuotas | `BAC-37` | Inferida |
| 177 | 2026-10-14 | `854c531` | Tayra | Merge remote-tracking branch 'origin/main' into feat/rf-02-rf-07-cuotas-host _(merge)_ | `BAC-37` | Inferida |
| 178 | 2026-10-14 | `148b44c` | Lisandro | mas documentacion actualizada | — | Inferida |
| 179 | 2026-10-14 | `083292e` | Nico | fix(aprovisionamiento): posibles arreglos a concurrencia | `BAC-37`, `BAC-38` | Mensaje |
| 180 | 2026-10-14 | `d8ed44f` | Nico | Merge remote-tracking branch 'origin/main' into fix/reservas-aprovisionamiento _(merge)_ | `BAC-37` | Inferida |
| 181 | 2026-10-14 | `84e8692` | Nico | Merge pull request #21 from tayraag/fix/reservas-aprovisionamiento _(merge)_ | `BAC-37` | Inferida |
| 182 | 2026-10-15 | `eaf7303` | Tayra | feat(instances): completar endpoint unificado de creacion | `BAC-38` | Mensaje |
| 183 | 2026-10-15 | `00e115c` | Tayra | docs(api): documentar creacion unificada de instancias | `BAC-38` | Inferida |
| 184 | 2026-10-15 | `c28bc01` | Tayra | merge(main): integrar reservas seguras de aprovisionamiento _(merge)_ | `BAC-37` | Inferida |
| 185 | 2026-10-15 | `6325b9b` | Tayra | docs(cuotas): alinear nombres y reservas persistentes | `BAC-37` | Inferida |
| 186 | 2026-10-15 | `dbf7075` | Tayra | refactor(cuotas): eliminar lógica obsoleta de vencimiento de reservas | `BAC-37` | Explícita |
| 187 | 2026-10-15 | `ad1ee86` | Lisandro | fix(ws): cortar el WebSocket de métricas al perder el acceso | `BAC-35` | Mensaje |
| 188 | 2026-10-15 | `3b04606` | Lisandro | fix(ws): errores explícitos de Proxmox, límites de conexión y control de origen | `BAC-35` | Mensaje |
| 189 | 2026-10-15 | `eb14053` | Lisandro | docs: refresh y logout con cookie en SETUP y README sin RRD | `SEC-01` | Inferida |
| 190 | 2026-10-15 | `44cfd16` | Lisandro | fix(metricas): CPU en porcentaje y null con la instancia apagada en todas las salidas | `BAC-34`, `BAC-35` | Mensaje |
| 191 | 2026-10-16 | `bf01245` | Lisandro | fix(db): migrar unicidad del correo después de crear auditoria | `FIX-81` | Explícita |
| 192 | 2026-10-16 | `e6aacfb` | Lisandro | test(db): quitar test redundante de la guarda de clasificarEliminados | `FIX-81` | Inferida |
| 193 | 2026-10-16 | `4f7686e` | Lisandro | test(simulador): ciclo completo de snapshots para QEMU y LXC | `BAC-50` | Mensaje |
| 194 | 2026-10-16 | `d290028` | Lisandro | feat(eventos): alertas del host ALERT_TRIGGERED solo para ADMIN | `BAC-52` | Explícita |
| 195 | 2026-10-16 | `e53aa43` | Lisandro | feat(instancias): exclusión mutua de operaciones por instancia (RequireInstanceIdle) | `BAC-51` | Explícita |
| 196 | 2026-10-16 | `bac10a5` | Lisandro | docs: README, SETUP y diseño del front al día con guardas, eventos y métricas | — | Inferida |
| 197 | 2026-10-16 | `691fadc` | Lisandro | feat(instancias): PUT /config con DeltaRAM contra el nodo y detección de hot-plug | `BAC-40` | Explícita |
| 198 | 2026-10-16 | `4df0f2c` | Lisandro | docs: contrato real de PUT /config en el diseño del front | `BAC-40` | Inferida |
| 199 | 2026-10-16 | `5dbf0e3` | Lisandro | fix(db): transacción anidable en la migración de auditoría legacy | `FIX-70` | Mensaje |
| 200 | 2026-10-16 | `c53fb3b` | Tayra | fix(tareas): cerrar creaciones solo con exitstatus OK | `BAC-39` | Explícita |
| 201 | 2026-10-16 | `b7f83bc` | Tayra | fix(permisos): autoasignar acceso solo en creaciones exitosas | `BAC-39` | Explícita |
| 202 | 2026-10-16 | `f0b0313` | Tayra | test(aprovisionamiento): cubrir lock y nivel de acceso | `BAC-39`, `BAC-51` | Mensaje |
| 203 | 2026-10-16 | `982ed12` | Tayra | docs(swagger): actualizar contrato de creación | `BAC-38` | Inferida |
| 204 | 2026-10-17 | `e543fe5` | Tayra | feat(auditoria): extender contrato estructurado e inmutable | `BAC-41` | Explícita |
| 205 | 2026-10-17 | `ebfea2d` | Tayra | feat(instancias): auditar el despacho de creaciones | `BAC-41` | Explícita |
| 206 | 2026-10-17 | `b3986aa` | Tayra | test(auditoria): validar correlacion transaccional por UPID | `BAC-41` | Mensaje |
| 207 | 2026-10-17 | `84f4aa0` | Tayra | feat(instancias): auditar reconfiguraciones sin ocultar fallos | `BAC-41` | Explícita |
| 208 | 2026-10-17 | `1eb3d4f` | Tayra | feat(auditoria): exponer filtros y detalles JSON al administrador | `BAC-41` | Mensaje |
| 209 | 2026-10-17 | `39bb224` | Tayra | docs(swagger): documentar auditoria RF-08 | `BAC-41` | Inferida |
| 210 | 2026-10-18 | `bd22c3d` | Tayra | fix(seguridad): ocultar UPID tras el identificador de tarea | `FIX-105` | Explícita |
| 211 | 2026-10-18 | `1f63959` | Tayra | docs(api): documentar tareaId como referencia publica | `FIX-105` | Inferida |
| 212 | 2026-10-18 | `e844600` | Tayra | feat(auditoria): agregar migracion preventiva de particiones | `BAC-53` | Explícita |
| 213 | 2026-10-18 | `fb533ce` | Tayra | feat(auditoria): garantizar indice y limites UTC al iniciar | `BAC-53` | Explícita |
| 214 | 2026-10-18 | `a61f451` | Tayra | test(auditoria): validar particiones pruning y migracion segura | `BAC-53` | Mensaje |
| 215 | 2026-10-18 | `30d89d6` | Nico | fix(auditoria): permitir excluir acciones del listado para el feed del dashboard | `BAC-44` | Mensaje |
| 216 | 2026-10-18 | `a898644` | Nico | feat(nodo): histórico del host desde el RRD y endpoint para el gráfico | `FRN-19C` | Explícita |
| 217 | 2026-10-18 | `43240ca` | Nico | Merge remote-tracking branch 'origin/main' into Nico _(merge)_ | `FRN-19C`, `BAC-44` | Inferida |
| 218 | 2026-10-18 | `8daaf87` | Nico | Merge pull request #22 from tayraag/Nico _(merge)_ | `FRN-19C`, `BAC-44` | Inferida |

---

## 4. Cruce por tarea: asignado declarado vs. autor real

Solo se incluyen tareas con relación con el backend (área Backend/Base de Datos/Infraestructura, tareas BAC/INF/INT y tareas referenciadas por algún commit).  
**Veredicto:** ✅ Coincide = todos los commits son de personas asignadas · ⚠️ Parcial = hay commits de asignados y de terceros · ❌ No coincide = ningún commit del asignado · ➖ Sin commits = no hay implementación en el submódulo.

| Tarea | Título | Área | Asignado(s) | Autoría real (commits) | Veredicto |
|---|---|---|---|---|---|
| `LOGIN-01` | Flujo 2FA/TOTP en memoria | Backend | Tayra y Lisandro | Tayra (1), Lisandro (1) | ✅ Coincide |
| `LOGIN-03` | Integración y prueba del prototipo | Frontend / Backend | Cristian, Tayra y Lisandro | — | ➖ Sin commits |
| `LOGIN-04` | Prueba integral de autenticación y autorización | Frontend / Backend | Cristian, Tayra y Lisandro | — | ➖ Sin commits |
| `BAC-01` | Script Docker Compose de base de datos local | Backend / Infraestructura | Tayra y Lucas | Lisandro (1), Nico (1) | ❌ No coincide |
| `BAC-02` | Modelo y hashing de contraseñas | Backend | Lisandro | Lisandro (2) | ✅ Coincide |
| `BAC-03` | Endpoint `POST /api/auth/login` | Backend | Tayra | Lisandro (4), Tayra (2), Nico (1) | ⚠️ Parcial (2/7) |
| `BAC-04` | Manejo de errores HTTP | Backend | Lisandro | Nico (2), Lisandro (1) | ⚠️ Parcial (1/3) |
| `BAC-05` | Migración y esquema relacional de usuarios y RBAC | Backend | Tayra y Lisandro | Lisandro (5) | ✅ Coincide |
| `BAC-06` | CRUD de usuarios para administradores | Backend | Lisandro | Lisandro (3) | ✅ Coincide |
| `BAC-07` | Asignación de permisos por recurso | Backend | Tayra | Lisandro (3) | ❌ No coincide |
| `BAC-08` | Middleware de autorización por recurso | Backend | Lisandro | Lisandro (2), Lucas (1) | ⚠️ Parcial (2/3) |
| `BAC-09` | Endpoint de roles del sistema (GET /api/roles) | Backend | Tayra | Lisandro (2) | ❌ No coincide |
| `BAC-10` | Enrolamiento 2FA y generación de QR | Backend | Tayra | Tayra (2), Lisandro (1) | ⚠️ Parcial (2/3) |
| `BAC-11` | Validación y persistencia del secreto TOTP | Backend | Lisandro | Tayra (2) | ❌ No coincide |
| `BAC-12` | Cambio obligatorio de contraseña temporal | Backend | Lisandro | Tayra (2) | ❌ No coincide |
| `BAC-13` | Restablecimiento administrativo de 2FA | Backend | Tayra | Tayra (1) | ✅ Coincide |
| `BAC-14` | Lectura mínima del inventario de Proxmox | Backend | Tayra | Lisandro (3), Lucas (1) | ❌ No coincide |
| `BAC-15` | Restablecimiento administrativo de contraseña | Backend | Lisandro | Tayra (2) | ❌ No coincide |
| `BAC-16` | Entrega segura de credenciales temporales (Mailer Service RF-13) | Backend / Infraestructura | Lisandro y Nico | Tayra (2) | ❌ No coincide |
| `BAC-16B` | Adaptador Go `SmtpEmailService` para envío real de credenciales y cód… | Backend | Lisandro | Tayra (1) | ❌ No coincide |
| `BAC-17` | Logout y revocación de sesión/JWT | Backend | Lisandro | Lisandro (3), Nico (2), Tayra (2) | ⚠️ Parcial (3/7) |
| `BAC-17A` | Adaptador base de Redis en el backend (conexión, configuración y puer… | Backend | Tayra | Lisandro (2), Lucas (1) | ❌ No coincide |
| `BAC-17B` | Corrección de lógica de inserción en `sesiones_activas` (1 sesión = 1… | Backend | Tayra | Lisandro (1) | ❌ No coincide |
| `BAC-18` | Base transversal de auditoría (append-only) | Backend | Tayra | Lucas (2), Nico (1), Tayra (1) | ⚠️ Parcial (1/4) |
| `BAC-18B` | Índice parcial y purga en `sesiones_activas`, y particionamiento trim… | Backend / Base de Datos | Tayra | Tayra (2) | ✅ Coincide |
| `BAC-19` | Solicitud de recuperación de contraseña (RF-13) | Backend | Lisandro | Tayra (3), Nico (1) | ❌ No coincide |
| `BAC-20` | Confirmación de recuperación de contraseña (RF-13) | Backend | Tayra | Tayra (1) | ✅ Coincide |
| `BAC-21` | Contrato base del canal de eventos/notificaciones (RF-11) | Backend / Frontend | Tayra y Lisandro | Nico (1), Lisandro (1) | ⚠️ Parcial (1/2) |
| `BAC-21B` | Alineación de esquema, rutas de energía, DELETE y GET /instances | Backend | Tayra y Lisandro | Tayra (2) | ✅ Coincide |
| `BAC-22` | Adaptador de telemetría del nodo con caché en Redis (`RF-02`) | Backend | Tayra | Lisandro (1), Tayra (1) | ⚠️ Parcial (1/2) |
| `BAC-23A` | Adaptador y normalización de inventario Proxmox (QEMU / LXC / IP) | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `BAC-24A` | Validación de estado previo y VMIDs protegidos en acciones de energía… | Backend | Lisandro | Lucas (1) | ❌ No coincide |
| `BAC-24B` | Endpoint de eliminación `DELETE /api/instances/:vmid` | Backend | Lisandro | Tayra (2) | ❌ No coincide |
| `BAC-24C` | Limpieza de permisos al eliminar una instancia (nueva) | Backend | Tayra | — | ➖ Sin commits |
| `BAC-25A` | Worker pool acotado para el seguimiento de UPID | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `BAC-25B` | Timeout configurable, reintentos y `exitstatus` en el evento | Frontend | Luz | Tayra (1) | ❌ No coincide |
| `BAC-25C` | Reanudación de UPIDs en curso al arrancar el backend (RNF-04) | Backend | Lisandro | Lisandro (1), Nico (1) | ⚠️ Parcial (1/2) |
| `BAC-26` | Verificación de `TASK_FINISHED` por `/api/events` para todas las acci… | Backend | Tayra | Tayra (1), Lisandro (1) | ⚠️ Parcial (1/2) |
| `BAC-27` | Auditoría del resultado de las acciones de ciclo de vida en `auditori… | Infraestructura | Nico | Tayra (3) | ❌ No coincide |
| `BAC-28` | Completar el simulador de Proxmox para la Etapa 1 | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `BAC-29` | Contrato HTTP y de eventos de la Etapa 1 (nueva) | Backend | Lisandro | Nico (4), Lisandro (2), Tayra (1) | ⚠️ Parcial (2/7) |
| `BAC-30` | Extensión del simulador para aprovisionamiento, storages y RRDdata | Backend / Testing | Lisandro | Lisandro (2) | ✅ Coincide |
| `BAC-31` | Endpoints de inventario de almacenamiento, plantillas e ISOs (RF-07) | Backend | Tayra | Lisandro (1) | ❌ No coincide |
| `BAC-32` | Extensión del worker UPID para creación y autoasignación de permisos | Backend | Lisandro | Tayra (2) | ❌ No coincide |
| `BAC-33` | Cálculo y exposición de cuotas asignables del host | Backend | Tayra | Tayra (2) | ✅ Coincide |
| `BAC-34` | Adaptador de Proxmox RRDdata en el cliente y servicio de telemetría d… | Backend | Lisandro | Lisandro (2) | ✅ Coincide |
| `BAC-35` | Canal WebSocket de métricas en tiempo real con ticket efímero | Backend | Tayra | Lisandro (5) | ❌ No coincide |
| `BAC-36` | Endpoint REST de métricas históricas de instancia (`timeframe=2h\|24h`) | Backend | Tayra | — | ➖ Sin commits |
| `BAC-37` | Servicio de validación de cuotas del host y bloqueo distribuido | Backend | Tayra | Tayra (7), Nico (1) | ⚠️ Parcial (7/8) |
| `BAC-38` | Endpoint `POST /api/instances` para creación de VMs y Contenedores LXC | Backend | Lisandro | Tayra (4), Nico (1) | ❌ No coincide |
| `BAC-39` | Seguimiento de creación, autoasignación de permisos y eventos (`RF-07… | Backend | Lisandro | Tayra (7) | ❌ No coincide |
| `BAC-40` | Endpoint `PUT /api/instances/:vmid/config` con cálculo diferencial ($… | Backend | Tayra | Lisandro (2) | ❌ No coincide |
| `BAC-41` | Auditoría inmutable de aprovisionamiento y cambios de recursos (`RF-0… | Backend / Base de Datos | Lisandro | Tayra (6) | ❌ No coincide |
| `BAC-42` | Adaptador de Proxmox para Snapshots, endpoints de Creación y Listado | Backend | Lisandro | — | ➖ Sin commits |
| `BAC-43` | Reversión (Rollback) y Eliminación de Snapshots con bloqueo de concur… | Backend | Tayra | — | ➖ Sin commits |
| `BAC-44` | Interceptor global y consultas complejas sobre auditoría particionada | Backend | Tayra | Nico (1) | ❌ No coincide |
| `BAC-45` | Streaming de exportación estructurada de auditoría a CSV y JSON | Backend | Lisandro | — | ➖ Sin commits |
| `BAC-46` | Background worker evaluador de umbrales con histéresis y deduplicación | Backend | Lisandro | — | ➖ Sin commits |
| `BAC-47` | Despacho multicanal de alertas (SSE en tiempo real y correo estructur… | Backend | Tayra | — | ➖ Sin commits |
| `BAC-48` | Hardening de autenticación: Cookies `HttpOnly`, Rate Limiting y Sanit… | Backend | Tayra | — | ➖ Sin commits |
| `BAC-50` | Extensión del simulador para snapshots y rollback | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `BAC-51` | Guarda centralizada de concurrencia operativa (409 INSTANCE_BUSY) | Backend | Tayra | Lisandro (1), Tayra (1) | ⚠️ Parcial (1/2) |
| `BAC-52` | Generalización del broker SSE para alertas globales | Backend | Tayra | Lisandro (1) | ❌ No coincide |
| `BAC-53` | Particionado preventivo e indexación compuesta de auditoría | Backend | Lisandro | Tayra (3) | ❌ No coincide |
| `SEC-01` | Refresh token en cookie HttpOnly en backend | Backend | Lisandro | Tayra (2), Nico (1), Lisandro (1) | ⚠️ Parcial (1/4) |
| `SEC-04` | Esquema extensible de niveles de acceso a recursos en Backend y Model… | Backend | Lisandro / Tayra | Lisandro (1) | ✅ Coincide |
| `FIX-01` | Colisión de índice único rompe el login repetido (CRÍTICA) | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `FIX-03` | El guard de re-enrolamiento 2FA no rechaza reemplazar un secreto ya v… | Backend | Tayra | Tayra (2) | ✅ Coincide |
| `FIX-13` | Desalineación de endpoints y PR pendiente en recuperación de contrase… | Backend | Lisandro y Tayra | Tayra (1) | ✅ Coincide |
| `FIX-16` | Montar middleware de autorización por recurso en rutas de instancias … | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `FIX-17` | Alinear y formalizar endpoints y contratos de contraseñas (Backend) | Backend | Lisandro / Tayra | Tayra (1) | ✅ Coincide |
| `FIX-23` | Hacer append-only la tabla `auditoria` a nivel de base de datos (`BAC… | Backend | Tayra | Tayra (2) | ✅ Coincide |
| `FIX-31` | Versionar y verificar el TLS de Nginx en el borde (`INF-05`) (Infraes… | Backend | Tayra y Lisandro | — | ➖ Sin commits |
| `FIX-32` | Fidelidad del simulador de Proxmox con la API real (Backend) | Backend | Lisandro | Lisandro (5) | ✅ Coincide |
| `FIX-33` | Conflicto de bloqueo de Proxmox devuelto como `502 PROXMOX_UNAVAILABL… | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `FIX-37` | Informar el nivel de acceso por instancia en `GET /account/profile` (… | Backend | Tayra | Tayra (1) | ✅ Coincide |
| `FIX-39` | Completar `BAC-21B`: campos de `GET /api/instances`, auditoría de ene… | Backend | Tayra y Lisandro | Tayra (2) | ✅ Coincide |
| `FIX-40` | Distinguir `504 PROXMOX_TIMEOUT` de `502 PROXMOX_UNAVAILABLE` (decisi… | Backend | Lisandro | Lisandro (1) | ✅ Coincide |
| `FIX-44` | Corregir índice parcial en `sesiones_activas` para incluir `jti_acces… | Backend | Tayra | Tayra (1) | ✅ Coincide |
| `FIX-46` | Completar el contrato de la Etapa 1 y alinear `activeTask` (`BAC-29`)… | Backend | Nico | Nico (3), Lisandro (1) | ⚠️ Parcial (3/4) |
| `FIX-48` | Documentar `INSTANCE_INVALID_STATE` en Swagger (`BAC-24A`) (Backend) | Backend | Lucas | Lucas (1) | ✅ Coincide |
| `FIX-49` | `DELETE /api/instances/:vmid`: código de estado (D2) y `202` con segu… | Backend | Tayra | Lisandro (1), Tayra (1) | ⚠️ Parcial (1/2) |
| `FIX-50` | Apagado correcto del servidor HTTP ante `SIGINT`/`SIGTERM` (`BAC-18B`… | Backend | Tayra | Tayra (1) | ✅ Coincide |
| `FIX-51` | El historial de `auditoria` previo a la migración queda inaccesible d… | Backend | Tayra | Tayra (1) | ✅ Coincide |
| `FIX-52` | Distinguir `INSTANCE_BUSY` de `INSTANCE_INVALID_STATE` y limpiar códi… | Backend | Lisandro | Lucas (1) | ❌ No coincide |
| `FIX-53` | Unificar códigos de acción y resultado en la auditoría de instancias … | Backend | Lisandro | Tayra (1), Nico (1) | ❌ No coincide |
| `FIX-54` | Separar eliminación lógica y suspensión; migrar unicidad del correo (… | Backend | — | Lisandro (1) | ❓ Sin asignación |
| `FIX-55` | Alinear baja, validación y autenticación con identidad histórica (Bac… | Backend | — | Lisandro (1) | ❓ Sin asignación |
| `FIX-62` | Heurística de circuit breaker y recuperación de telemetría en `nodo_s… | Backend | Tayra / Lisandro | Tayra (1) | ✅ Coincide |
| `FIX-63` | Cobertura de pruebas unitarias para `ObtenerInterfaces` en el cliente… | Backend / pruebas | Lisandro | Nico (1) | ❌ No coincide |
| `FIX-65` | Explicitar en `docs/contrato-etapa1.md` el acceso de cualquier usuari… | Backend / documentación | Nico | Nico (1) | ✅ Coincide |
| `FIX-69` | Restaurar definiciones SSE, esquemas de eventos y contratos en `swagg… | Backend | Lisandro / Lucas | Nico (1) | ❌ No coincide |
| `FIX-70` | Blindaje transaccional de la migración de auditoría y saneamiento de … | Backend | Tayra | Tayra (1), Lisandro (1) | ⚠️ Parcial (1/2) |
| `FIX-71` | Robustecer reconexión y reintentos ante arranque en frío de Redis (`I… | Backend | Lisandro / Lucas | Tayra (1) | ❌ No coincide |
| `FIX-72` | Sincronizar firma de `NewSeguimientoTareas` en pruebas unitarias de s… | Backend / pruebas | Lisandro / Lucas | Lisandro (1), Tayra (1) | ⚠️ Parcial (1/2) |
| `FIX-73` | Sincronización y manejo determinista de tareas consecutivas en accion… | Infraestructura | Nico / Lucas | Nico (1) | ✅ Coincide |
| `FIX-81` | Corregir el orden de migraciones: `migrarUnicidadEmail` antes de crea… | Backend | Lisandro / Tayra | Lisandro (2) | ✅ Coincide |
| `FIX-84` | Completar storages y contenido en el simulador de Proxmox (`BAC-30`) … | Backend / Testing | Lisandro | Lisandro (3) | ✅ Coincide |
| `FIX-88` | Evitar que `EMAIL_PROVIDER` caiga en silencio en el correo simulado (… | Backend | Tayra / Lisandro | Tayra (1) | ✅ Coincide |
| `FIX-105` | Desacople de UPID y adopción del contrato seguro `tareaId` en Fronten… | Frontend | Cristian / Luz | Tayra (2) | ❌ No coincide |
| `FIX-107` | Cifrado y ofuscación reversible de VMIDs de Proxmox contra IDOR y enu… | Backend | Tayra | — | ➖ Sin commits |
| `FIX-109` | Enmascaramiento de nombres de host, nodos y pools de almacenamiento (… | Backend | Lisandro / Tayra | — | ➖ Sin commits |
| `FIX-110` | Saneamiento estricto de mensajes de error y canales en tiempo real (I… | Backend | Tayra | — | ➖ Sin commits |
| `FIX-111A` | Blindaje contra secuestro de cuenta (Account Takeover) en cambio de e… | Backend | Lisandro | — | ➖ Sin commits |
| `FIX-112` | Middleware de Rate Limiting y protección contra fuerza bruta en auten… | Backend | Lisandro / Nico | — | ➖ Sin commits |
| `FIX-113` | Purga transaccional de permisos huérfanos al eliminar instancias Prox… | Backend | Tayra / Lisandro | — | ➖ Sin commits |
| `FIX-114` | Revocación obligatoria de sesiones activas en cambio voluntario de co… | Backend | Lisandro | — | ➖ Sin commits |
| `FIX-115A` | Ampliación de cota máxima de contraseñas a 128 caracteres en Backend … | Backend | Lisandro | — | ➖ Sin commits |
| `FIX-116` | Sanitización estricta y prevención de inyección de parámetros en Prox… | Backend | Lisandro / Tayra | — | ➖ Sin commits |
| `FIX-117A` | Restricción de acceso y segregación de telemetría global del host par… | Backend | Tayra | — | ➖ Sin commits |
| `FIX-118` | Eliminación de credenciales por defecto hardcodeadas en migraciones d… | Backend | Lisandro | — | ➖ Sin commits |
| `FIX-119` | Validación de existencia de VMIDs y bloqueo de instancias protegidas … | Backend | Lisandro / Tayra | — | ➖ Sin commits |
| `INF-03` | PostgreSQL persistente en el servidor | Infraestructura | Lucas y Nico | Nico (8) | ✅ Coincide |
| `INF-05` | CORS y TLS en el borde | Infraestructura | Nico | Lucas (1) | ❌ No coincide |
| `INF-06A` | Redis local para desarrollo en el repositorio del backend | Backend (en el repositorio del backend) | Tayra | Lucas (1) | ❌ No coincide |
| `INF-07` | Token de Proxmox, conectividad y VMIDs protegidos en el servidor | Frontend | Nico | Lucas (1), Lisandro (1) | ❌ No coincide |
| `INF-08B` | Credenciales SMTP en `.env.example` y `.env` del repositorio del back… | Backend (en el repositorio del backend) | Lisandro | Tayra (2) | ❌ No coincide |
| `INF-09` | Nginx para WebSockets y streaming bidireccional (RNF-06) | Infraestructura | Nico | Lucas (1) | ❌ No coincide |
| `INT-01` | Pruebas de integración automatizadas de ciclo de vida y worker UPID | Backend / Testing | Tayra y Cristian | — | ➖ Sin commits |
| `INT-04` | Suite automatizada de pruebas de integración de la Etapa 2 | Backend / Testing | Tayra y Lisandro | — | ➖ Sin commits |
| `INT-07` | Suite automatizada de pruebas de integración de la Etapa 3 | Backend / Testing | Tayra y Lisandro | — | ➖ Sin commits |
| `FRN-19C` | Soporte backend del histórico del host para el Dashboard (tarea front… | Frontend | Belinda | Nico (1) | ❌ No coincide |

### 4.1. Resumen del cruce

| Resultado | Tareas | % |
|---|---:|---:|
| ✅ Coincide | 38 | 31.1% |
| ⚠️ Parcial | 18 | 14.8% |
| ❌ No coincide | 38 | 31.1% |
| ➖ Sin commits | 26 | 21.3% |
| ❓ Sin asignación | 2 | 1.6% |
| **Total** | **122** | **100%** |

---

## 5. Discrepancias detectadas

### 5.1. Tareas donde el asignado no escribió ningún commit (❌)

| Tarea | Título | Asignado(s) | Autor real | Commits | Evidencia (hash) |
|---|---|---|---|---:|---|
| `BAC-39` | Seguimiento de creación, autoasignación de permisos y event… | Lisandro | Tayra (7) | 7 | `e37ed4d`, `882c828`, `9891378`, `468e452`, `c53fb3b`, `b7f83bc`, `f0b0313` |
| `BAC-41` | Auditoría inmutable de aprovisionamiento y cambios de recur… | Lisandro | Tayra (6) | 6 | `e543fe5`, `ebfea2d`, `b3986aa`, `84f4aa0`, `1eb3d4f`, `39bb224` |
| `BAC-35` | Canal WebSocket de métricas en tiempo real con ticket efíme… | Tayra | Lisandro (5) | 5 | `394b25b`, `993c86f`, `ad1ee86`, `3b04606`, `44cfd16` |
| `BAC-38` | Endpoint `POST /api/instances` para creación de VMs y Conte… | Lisandro | Tayra (4), Nico (1) | 5 | `468e452`, `083292e`, `eaf7303`, `00e115c`, `982ed12` |
| `BAC-14` | Lectura mínima del inventario de Proxmox | Tayra | Lisandro (3), Lucas (1) | 4 | `f3c659a`, `60a6403`, `44feb08`, `725d436` |
| `BAC-19` | Solicitud de recuperación de contraseña (RF-13) | Lisandro | Tayra (3), Nico (1) | 4 | `e363afa`, `c89ddcb`, `0e68371`, `2bdc3b1` |
| `BAC-07` | Asignación de permisos por recurso | Tayra | Lisandro (3) | 3 | `3052a02`, `de91ca3`, `7a2403c` |
| `BAC-17A` | Adaptador base de Redis en el backend (conexión, configurac… | Tayra | Lisandro (2), Lucas (1) | 3 | `2a943cd`, `546e1d5`, `0aac034` |
| `BAC-27` | Auditoría del resultado de las acciones de ciclo de vida en… | Nico | Tayra (3) | 3 | `d76eb35`, `dc018a2`, `7acbb56` |
| `BAC-53` | Particionado preventivo e indexación compuesta de auditoría | Lisandro | Tayra (3) | 3 | `e844600`, `fb533ce`, `a61f451` |
| `BAC-01` | Script Docker Compose de base de datos local | Tayra y Lucas | Lisandro (1), Nico (1) | 2 | `b67468c`, `2e8b02a` |
| `BAC-09` | Endpoint de roles del sistema (GET /api/roles) | Tayra | Lisandro (2) | 2 | `538b494`, `f69a6f4` |
| `BAC-11` | Validación y persistencia del secreto TOTP | Lisandro | Tayra (2) | 2 | `9df22e6`, `348dd59` |
| `BAC-12` | Cambio obligatorio de contraseña temporal | Lisandro | Tayra (2) | 2 | `d192003`, `348dd59` |
| `BAC-15` | Restablecimiento administrativo de contraseña | Lisandro | Tayra (2) | 2 | `019017b`, `5a7abe1` |
| `BAC-16` | Entrega segura de credenciales temporales (Mailer Service R… | Lisandro y Nico | Tayra (2) | 2 | `e363afa`, `860b3c9` |
| `BAC-24B` | Endpoint de eliminación `DELETE /api/instances/:vmid` | Lisandro | Tayra (2) | 2 | `088e175`, `f06640f` |
| `BAC-32` | Extensión del worker UPID para creación y autoasignación de… | Lisandro | Tayra (2) | 2 | `0022e18`, `e37ed4d` |
| `BAC-40` | Endpoint `PUT /api/instances/:vmid/config` con cálculo dife… | Tayra | Lisandro (2) | 2 | `691fadc`, `4df0f2c` |
| `FIX-105` | Desacople de UPID y adopción del contrato seguro `tareaId` … | Cristian / Luz | Tayra (2) | 2 | `bd22c3d`, `1f63959` |
| `FIX-53` | Unificar códigos de acción y resultado en la auditoría de i… | Lisandro | Tayra (1), Nico (1) | 2 | `b8d2631`, `b476636` |
| `INF-07` | Token de Proxmox, conectividad y VMIDs protegidos en el ser… | Nico | Lucas (1), Lisandro (1) | 2 | `44feb08`, `82add96` |
| `INF-08B` | Credenciales SMTP en `.env.example` y `.env` del repositori… | Lisandro | Tayra (2) | 2 | `233e804`, `860b3c9` |
| `BAC-16B` | Adaptador Go `SmtpEmailService` para envío real de credenci… | Lisandro | Tayra (1) | 1 | `860b3c9` |
| `BAC-17B` | Corrección de lógica de inserción en `sesiones_activas` (1 … | Tayra | Lisandro (1) | 1 | `9e0f43c` |
| `BAC-24A` | Validación de estado previo y VMIDs protegidos en acciones … | Lisandro | Lucas (1) | 1 | `8591e90` |
| `BAC-25B` | Timeout configurable, reintentos y `exitstatus` en el evento | Luz | Tayra (1) | 1 | `4f68e44` |
| `BAC-31` | Endpoints de inventario de almacenamiento, plantillas e ISO… | Tayra | Lisandro (1) | 1 | `2d2495f` |
| `BAC-44` | Interceptor global y consultas complejas sobre auditoría pa… | Tayra | Nico (1) | 1 | `30d89d6` |
| `BAC-52` | Generalización del broker SSE para alertas globales | Tayra | Lisandro (1) | 1 | `d290028` |
| `FIX-52` | Distinguir `INSTANCE_BUSY` de `INSTANCE_INVALID_STATE` y li… | Lisandro | Lucas (1) | 1 | `08762f4` |
| `FIX-63` | Cobertura de pruebas unitarias para `ObtenerInterfaces` en … | Lisandro | Nico (1) | 1 | `973ee0b` |
| `FIX-69` | Restaurar definiciones SSE, esquemas de eventos y contratos… | Lisandro / Lucas | Nico (1) | 1 | `464b494` |
| `FIX-71` | Robustecer reconexión y reintentos ante arranque en frío de… | Lisandro / Lucas | Tayra (1) | 1 | `c7fd7a1` |
| `FRN-19C` | Soporte backend del histórico del host para el Dashboard (t… | Belinda | Nico (1) | 1 | `a898644` |
| `INF-05` | CORS y TLS en el borde | Nico | Lucas (1) | 1 | `c624a57` |
| `INF-06A` | Redis local para desarrollo en el repositorio del backend | Tayra | Lucas (1) | 1 | `2a943cd` |
| `INF-09` | Nginx para WebSockets y streaming bidireccional (RNF-06) | Nico | Lucas (1) | 1 | `58388fc` |

### 5.2. Tareas con autoría compartida / parcial (⚠️)

El asignado participó, pero otra persona escribió parte de los commits. El porcentaje indica la proporción de commits del/los asignados.

| Tarea | Título | Asignado(s) | Autoría real | % del asignado | Evidencia (hash) |
|---|---|---|---|---:|---|
| `BAC-18` | Base transversal de auditoría (append-only) | Tayra | Lucas (2), Nico (1), Tayra (1) | 1/4 (25%) | `d6fdcb1`, `d15a8fe`, `3d06d7c`, `ad52485` |
| `SEC-01` | Refresh token en cookie HttpOnly en backend | Lisandro | Tayra (2), Nico (1), Lisandro (1) | 1/4 (25%) | `5d8a7fd`, `c719c9b`, `47bb8f2`, `eb14053` |
| `BAC-03` | Endpoint `POST /api/auth/login` | Tayra | Lisandro (4), Tayra (2), Nico (1) | 2/7 (29%) | `aaaf3d3`, `807a523`, `92b7dcb`, `4e8e9c4`, `e70701d`, `0a5faad`, `fec42b7` |
| `BAC-29` | Contrato HTTP y de eventos de la Etapa 1 (nueva) | Lisandro | Nico (4), Lisandro (2), Tayra (1) | 2/7 (29%) | `45ecad7`, `5787179`, `21b1332`, `4ef36e3`, `8ba0c59`, `0850ec7`, `8f46c3d` |
| `BAC-04` | Manejo de errores HTTP | Lisandro | Nico (2), Lisandro (1) | 1/3 (33%) | `8d7bce5`, `fea91a8`, `0e97f93` |
| `BAC-17` | Logout y revocación de sesión/JWT | Lisandro | Lisandro (3), Nico (2), Tayra (2) | 3/7 (43%) | `e70701d`, `addb9b6`, `cd04708`, `dbe4079`, `732e65d`, `5d8a7fd`, `15032da` |
| `BAC-21` | Contrato base del canal de eventos/notificaciones (RF-11) | Tayra y Lisandro | Nico (1), Lisandro (1) | 1/2 (50%) | `d509d18`, `45ecad7` |
| `BAC-22` | Adaptador de telemetría del nodo con caché en Redis (`RF-02… | Tayra | Lisandro (1), Tayra (1) | 1/2 (50%) | `13f9c35`, `8f46c3d` |
| `BAC-25C` | Reanudación de UPIDs en curso al arrancar el backend (RNF-0… | Lisandro | Lisandro (1), Nico (1) | 1/2 (50%) | `f17295c`, `1a6015e` |
| `BAC-26` | Verificación de `TASK_FINISHED` por `/api/events` para toda… | Tayra | Tayra (1), Lisandro (1) | 1/2 (50%) | `7acbb56`, `f2215af` |
| `BAC-51` | Guarda centralizada de concurrencia operativa (409 INSTANCE… | Tayra | Lisandro (1), Tayra (1) | 1/2 (50%) | `e53aa43`, `f0b0313` |
| `FIX-49` | `DELETE /api/instances/:vmid`: código de estado (D2) y `202… | Tayra | Lisandro (1), Tayra (1) | 1/2 (50%) | `0850ec7`, `f06640f` |
| `FIX-70` | Blindaje transaccional de la migración de auditoría y sanea… | Tayra | Tayra (1), Lisandro (1) | 1/2 (50%) | `08ca9dd`, `5dbf0e3` |
| `FIX-72` | Sincronizar firma de `NewSeguimientoTareas` en pruebas unit… | Lisandro / Lucas | Lisandro (1), Tayra (1) | 1/2 (50%) | `c108594`, `ec84122` |
| `BAC-08` | Middleware de autorización por recurso | Lisandro | Lisandro (2), Lucas (1) | 2/3 (67%) | `eb84a58`, `d11426e`, `44feb08` |
| `BAC-10` | Enrolamiento 2FA y generación de QR | Tayra | Tayra (2), Lisandro (1) | 2/3 (67%) | `807a523`, `924641f`, `019017b` |
| `FIX-46` | Completar el contrato de la Etapa 1 y alinear `activeTask` … | Nico | Nico (3), Lisandro (1) | 3/4 (75%) | `5787179`, `21b1332`, `4ef36e3`, `0850ec7` |
| `BAC-37` | Servicio de validación de cuotas del host y bloqueo distrib… | Tayra | Tayra (7), Nico (1) | 7/8 (88%) | `d1b3aeb`, `a363c9a`, `acec121`, `d4ee23b`, `ff92dd3`, `083292e`, `6325b9b`, `dbf7075` |

### 5.3. Tareas de backend sin ningún commit en el submódulo (➖)

Pueden estar pendientes, o ejecutadas fuera del backend (frontend/infraestructura). Se listan para no confundirlas con trabajo entregado.

| Tarea | Título | Área | Asignado(s) | Estado inferido |
|---|---|---|---|---|
| `LOGIN-03` | Integración y prueba del prototipo | Frontend / Backend | Cristian, Tayra y Lisandro | Se ejecuta en `test/back` / `test/front` (fuera del submódulo) |
| `LOGIN-04` | Prueba integral de autenticación y autorización | Frontend / Backend | Cristian, Tayra y Lisandro | Se ejecuta en `test/back` / `test/front` (fuera del submódulo) |
| `BAC-24C` | Limpieza de permisos al eliminar una instancia (nueva) | Backend | Tayra | Pendiente / sin implementación en `backend/` |
| `BAC-36` | Endpoint REST de métricas históricas de instancia (`timeframe=2h\|24h`) | Backend | Tayra | Pendiente / sin implementación en `backend/` |
| `BAC-42` | Adaptador de Proxmox para Snapshots, endpoints de Creación y Listado | Backend | Lisandro | Pendiente / sin implementación en `backend/` |
| `BAC-43` | Reversión (Rollback) y Eliminación de Snapshots con bloqueo de concur… | Backend | Tayra | Pendiente / sin implementación en `backend/` |
| `BAC-45` | Streaming de exportación estructurada de auditoría a CSV y JSON | Backend | Lisandro | Pendiente / sin implementación en `backend/` |
| `BAC-46` | Background worker evaluador de umbrales con histéresis y deduplicación | Backend | Lisandro | Pendiente / sin implementación en `backend/` |
| `BAC-47` | Despacho multicanal de alertas (SSE en tiempo real y correo estructur… | Backend | Tayra | Pendiente / sin implementación en `backend/` |
| `BAC-48` | Hardening de autenticación: Cookies `HttpOnly`, Rate Limiting y Sanit… | Backend | Tayra | Pendiente / sin implementación en `backend/` |
| `FIX-31` | Versionar y verificar el TLS de Nginx en el borde (`INF-05`) (Infraes… | Backend | Tayra y Lisandro | Infraestructura (Nginx/TLS), fuera de `backend/` |
| `FIX-107` | Cifrado y ofuscación reversible de VMIDs de Proxmox contra IDOR y enu… | Backend | Tayra | Pendiente / sin implementación en `backend/` |
| `FIX-109` | Enmascaramiento de nombres de host, nodos y pools de almacenamiento (… | Backend | Lisandro / Tayra | Pendiente / sin implementación en `backend/` |
| `FIX-110` | Saneamiento estricto de mensajes de error y canales en tiempo real (I… | Backend | Tayra | Pendiente / sin implementación en `backend/` |
| `FIX-111A` | Blindaje contra secuestro de cuenta (Account Takeover) en cambio de e… | Backend | Lisandro | Pendiente / sin implementación en `backend/` |
| `FIX-112` | Middleware de Rate Limiting y protección contra fuerza bruta en auten… | Backend | Lisandro / Nico | Pendiente / sin implementación en `backend/` |
| `FIX-113` | Purga transaccional de permisos huérfanos al eliminar instancias Prox… | Backend | Tayra / Lisandro | Pendiente / sin implementación en `backend/` |
| `FIX-114` | Revocación obligatoria de sesiones activas en cambio voluntario de co… | Backend | Lisandro | Pendiente / sin implementación en `backend/` |
| `FIX-115A` | Ampliación de cota máxima de contraseñas a 128 caracteres en Backend … | Backend | Lisandro | Pendiente / sin implementación en `backend/` |
| `FIX-116` | Sanitización estricta y prevención de inyección de parámetros en Prox… | Backend | Lisandro / Tayra | Pendiente / sin implementación en `backend/` |
| `FIX-117A` | Restricción de acceso y segregación de telemetría global del host par… | Backend | Tayra | Pendiente / sin implementación en `backend/` |
| `FIX-118` | Eliminación de credenciales por defecto hardcodeadas en migraciones d… | Backend | Lisandro | Pendiente / sin implementación en `backend/` |
| `FIX-119` | Validación de existencia de VMIDs y bloqueo de instancias protegidas … | Backend | Lisandro / Tayra | Pendiente / sin implementación en `backend/` |
| `INT-01` | Pruebas de integración automatizadas de ciclo de vida y worker UPID | Backend / Testing | Tayra y Cristian | Se ejecuta en `test/back` / `test/front` (fuera del submódulo) |
| `INT-04` | Suite automatizada de pruebas de integración de la Etapa 2 | Backend / Testing | Tayra y Lisandro | Se ejecuta en `test/back` / `test/front` (fuera del submódulo) |
| `INT-07` | Suite automatizada de pruebas de integración de la Etapa 3 | Backend / Testing | Tayra y Lisandro | Se ejecuta en `test/back` / `test/front` (fuera del submódulo) |

### 5.4. Tareas sin asignación declarada (❓)

- `FIX-54` — Separar eliminación lógica y suspensión; migrar unicidad del correo (Backend) · autor real: **Lisandro (1)** · evidencia: `a9b9d0f`
- `FIX-55` — Alinear baja, validación y autenticación con identidad histórica (Backend) · autor real: **Lisandro (1)** · evidencia: `c024ba3`

### 5.5. Observaciones de metadata que afectan el cruce

- **`FIX-81`**: la ficha lo rotula como Frontend en algunos índices, pero el trabajo es backend y la sección real lo asigna a **Lisandro / Tayra**. Los commits son de Lisandro. Se usó la ficha real, no el índice.
- **`FIX-105`**: la tarea documentada es frontend (Cristian/Luz), pero `bd22c3d` implementa el **lado backend** del ocultamiento del UPID. No existe una subtarea backend separada: por eso figura como “No coincide”, aunque en rigor es un trabajo backend no asignado explícitamente.
- **`BAC-25B`**: tarea frontend (Luz) que incluye “timeout configurable”; el commit `4f68e44` (Tayra) implementa la parte backend. Se contabiliza como trabajo backend de una tarea cross-area.
- **`FRN-19C`**: el commit `a898644` (Nico) es el soporte backend del histórico del host para una tarea de frontend. No se asignó a `BAC-36` porque `BAC-36` refiere a métricas históricas **de instancia**, no del host; `BAC-36` sigue sin commits.
- **`BAC-29`** (contrato) aparece como tarea de Lisandro, pero la mayor parte de los commits de contrato los escribió **Nico** (`5787179`, `21b1332`, `4ef36e3`, `8ba0c59`).

---

## 6. Conclusiones

- De **122 tareas** analizadas, **38** tienen coincidencia total entre asignado y autor, **18** son parciales, **38** no tienen ningún commit del asignado y **26** no registran commits en el submódulo.
- Es decir, **56 tareas (46%)** tienen al menos un commit de una persona distinta a la asignada. Ese es el desvío que se buscaba detectar.

**Trabajo realizado en tareas ajenas (commits de una persona en tareas que no le fueron asignadas):**

| Persona | Tareas ajenas en las que commiteó | Commits en tareas ajenas |
|---|---:|---:|
| **Tayra** | 22 | 52 |
| **Lisandro** | 19 | 33 |
| **Nico** | 17 | 22 |
| **Lucas** | 10 | 11 |

**Casos más relevantes (tareas 100% ajenas y con volumen):**

- **`BAC-39`** — Seguimiento de creación, autoasignación de permisos y eventos (`RF-07`, `RF-11`): asignada a **Lisandro**, ejecutada por **Tayra (7)** (7 commits).
- **`BAC-41`** — Auditoría inmutable de aprovisionamiento y cambios de recursos (`RF-08`): asignada a **Lisandro**, ejecutada por **Tayra (6)** (6 commits).
- **`BAC-35`** — Canal WebSocket de métricas en tiempo real con ticket efímero: asignada a **Tayra**, ejecutada por **Lisandro (5)** (5 commits).
- **`BAC-38`** — Endpoint `POST /api/instances` para creación de VMs y Contenedores LXC: asignada a **Lisandro**, ejecutada por **Tayra (4), Nico (1)** (5 commits).
- **`BAC-14`** — Lectura mínima del inventario de Proxmox: asignada a **Tayra**, ejecutada por **Lisandro (3), Lucas (1)** (4 commits).
- **`BAC-19`** — Solicitud de recuperación de contraseña (RF-13): asignada a **Lisandro**, ejecutada por **Tayra (3), Nico (1)** (4 commits).
- **`BAC-07`** — Asignación de permisos por recurso: asignada a **Tayra**, ejecutada por **Lisandro (3)** (3 commits).
- **`BAC-17A`** — Adaptador base de Redis en el backend (conexión, configuración y puerto): asignada a **Tayra**, ejecutada por **Lisandro (2), Lucas (1)** (3 commits).
- **`BAC-27`** — Auditoría del resultado de las acciones de ciclo de vida en `auditoria` (`RF-08`): asignada a **Nico**, ejecutada por **Tayra (3)** (3 commits).
- **`BAC-53`** — Particionado preventivo e indexación compuesta de auditoría: asignada a **Lisandro**, ejecutada por **Tayra (3)** (3 commits).
- **`BAC-01`** — Script Docker Compose de base de datos local: asignada a **Tayra, Lucas**, ejecutada por **Lisandro (1), Nico (1)** (2 commits).
- **`BAC-09`** — Endpoint de roles del sistema (GET /api/roles): asignada a **Tayra**, ejecutada por **Lisandro (2)** (2 commits).

**Lectura por persona (backend):**

- **Tayra** concentra la mayor cantidad de commits y absorbió tareas originalmente asignadas a Lisandro en aprovisionamiento y auditoría (`BAC-38`, `BAC-39`, `BAC-41`, `BAC-53`), además de `BAC-27` (asignada a Nico).
- **Lisandro** escribió la mayor parte de la base de auth/usuarios y, a la inversa, resolvió tareas asignadas a Tayra (`BAC-22`, `BAC-35`, `BAC-40`).
- **Nico** aportó contratos, tests y fixes en tareas asignadas a otros (`BAC-29`, `BAC-44`, `FIX-53`, `FIX-63`, `FIX-69`) y el soporte backend del histórico del host (`a898644`); en cambio, `BAC-27`, que tenía asignada, la ejecutó Tayra.
- **Lucas** resolvió `BAC-24A`, `FIX-52`, `INF-05`, `INF-06A` e `INF-09`, todas asignadas a otras personas en la documentación.

**Recomendaciones:**

1. **Reasignar en ClickUp** las tareas listadas en 5.1 o convertirlas en co-asignadas, para que el tablero refleje la autoría real.
2. **Convención de commits con ID de tarea** (`BAC-xx`, `FIX-xx`) en el asunto o en el pie, como ya se hizo en varios commits (`RF-08`, `RF-05`, `FIX-81`). Reduce la inferencia y hace auditable el cruce.
3. **Separar subtareas backend/frontend** cuando una tarea cruza capas (`FIX-105`, `BAC-25B`), con asignado propio para cada lado.
4. **Reconciliar las tareas sin commits** (5.3) antes de darlas por terminadas: `BAC-36`, `BAC-42`, `BAC-43`, `BAC-45`–`BAC-48` y las tareas backend de hardening (`FIX-107`, `FIX-109`, `FIX-110`, `FIX-111A`, `FIX-112`–`FIX-114`, `FIX-115A`, `FIX-116`–`FIX-119`) siguen sin implementación en `backend/`.
5. **Definir code ownership** por módulo (auth/usuarios, instancias/UPID, auditoría, telemetría/WS, aprovisionamiento) para que las revisiones y reasignaciones no dependan de quién tuvo tiempo.

---

## 7. Commits críticos y su impacto

| Commit | Autor | Tarea | Por qué es crítico |
|---|---|---|---|
| `92b7dcb` | Lisandro | `BAC-03` | Revierte el primer intento de login/auth de Tayra y reinicia la implementación sobre otra base. |
| `807a523` | Lisandro | `BAC-02/BAC-03/LOGIN-01/BAC-10` | Implementación real de autenticación, 2FA TOTP, JWT y control de acceso (RF-01). |
| `d6fdcb1` | Lucas | `BAC-18` | Crea el esqueleto de auditoría; la tarea estaba asignada a Tayra. |
| `ad52485 / 9554efa` | Tayra | `BAC-18/FIX-23` | Inmutabilidad física y legal de `auditoria` (trigger y permisos). |
| `44a2339` | Lisandro | `FIX-32` | Simulador local + 504 en timeouts: habilita el desarrollo sin Proxmox real. |
| `0850ec7` | Lisandro | `BAC-29/FIX-46/FIX-49` | Cambia el contrato: `activeTask` como objeto, DELETE asíncrono y `TASK_FINISHED` completo. Impacta al frontend. |
| `088e175 / f06640f` | Tayra | `BAC-24B/FIX-49` | Implementa `DELETE /api/instances/:vmid` asíncrono; la tarea estaba asignada a Lisandro. |
| `e543fe5` | Tayra | `BAC-41` | Extiende el contrato estructurado e inmutable de auditoría; tarea asignada a Lisandro. |
| `bd22c3d` | Tayra | `FIX-105` | Oculta el UPID tras `tareaId`. Es el cambio que rompió el frontend y originó `FIX-104B`/`FIX-105`. |
| `e844600 / fb533ce` | Tayra | `BAC-53` | Migración preventiva de particiones e índice de auditoría; tarea asignada a Lisandro. |
| `a898644` | Nico | `FRN-19C` | Agrega el histórico del host desde RRD y el endpoint para el gráfico del Dashboard. |
| `691fadc` | Lisandro | `BAC-40` | PUT /config con DeltaRAM y detección de hot-plug; tarea asignada a Tayra. |

---

## Anexo A. Commits sin tarea asociada

Son bootstrap, CI, documentación o merges sin código propio. No se les asignó tarea para no forzar relaciones inexistentes.

| # | Fecha | Hash | Autor | Asunto | Motivo |
|---:|---|---|---|---|---|
| 1 | 2026-09-06 | `cd6a83e` | Lisandro | Initial commit: estructura base del backend y documentacion | Documentación / Swagger |
| 2 | 2026-09-06 | `ebee84e` | Lisandro | Agrego md a todas las carpetas | Bootstrap / chore |
| 19 | 2026-09-15 | `4919a88` | Lisandro | Merge branch 'main' of https://github.com/tayraag/centinela-back | Merge |
| 26 | 2026-09-17 | `31e9bd9` | Lisandro | fix(swagger): remove hardcoded localhost host to allow dynamic remote access | Documentación / Swagger |
| 27 | 2026-09-17 | `0955134` | Nico | fix(swagger): quitar @host para que Try it out apunte al servidor | Documentación / Swagger |
| 28 | 2026-09-17 | `8424bfd` | Nico | Merge pull request #4 from tayraag/fix/swagger-host | Merge |
| 29 | 2026-09-17 | `a37160c` | Lisandro | chore: format files and save local changes | Bootstrap / chore |
| 30 | 2026-09-17 | `0d3e9ca` | Lisandro | Merge branch 'main' - resolve conflict in swagger annotations | Merge |
| 34 | 2026-09-17 | `9b12618` | Lisandro | fix: eliminar import sin uso en account_handler.go | Bootstrap / chore |
| 36 | 2026-09-17 | `caf71c2` | Tayra | actulizacion del README | Documentación |
| 43 | 2026-09-21 | `a4c5215` | Lisandro | chore(docs): corregir SETUP.md para que coincida con el backend actual | Documentación |
| 61 | 2026-09-21 | `3dadc48` | Tayra | Resolucion de conflictos y correccion de build tras el merge | Merge |
| 67 | 2026-09-22 | `272a16e` | Nico | Merge remote-tracking branch 'origin/main' into Nico | Merge |
| 84 | 2026-09-27 | `c16bbee` | Lucas | Merge branch 'main' of https://github.com/tayraag/centinela-back | Merge |
| 92 | 2026-09-29 | `e6e7dc5` | Lisandro | chore(env): agregar ENABLE_SWAGGER a .env.example | Documentación / Swagger |
| 106 | 2026-10-02 | `573f0c3` | Tayra | docs: Actualizar documento de diseño de front para reflejar endpoints de auth, … | Documentación |
| 153 | 2026-10-10 | `a5e2110` | Tayra | Actualizacion de documentacion | Documentación |
| 178 | 2026-10-14 | `148b44c` | Lisandro | mas documentacion actualizada | Documentación |
| 196 | 2026-10-16 | `bac10a5` | Lisandro | docs: README, SETUP y diseño del front al día con guardas, eventos y métricas | Documentación |

## Anexo B. Limitaciones

1. **Fuente de asignación:** se usó el campo `Asignado/Asignada` de la documentación del repositorio, que declara ser espejo de ClickUp. Si ClickUp cambió después del 19/10/2026, el cruce debe reejecutarse contra la API de ClickUp.
2. **Inferencia:** los commits marcados como `Inferida` no citan la tarea. La relación se basa en el asunto, los archivos y el módulo; pueden existir errores de frontera entre tareas vecinas (p. ej. `BAC-41` vs `BAC-44`, `FIX-69` vs `FIX-46`).
3. **Autoría vs. autoría de Git:** se usa el autor del commit. No detecta pair programming, código copiado de otra rama ni commits reescritos por un merge.
4. **Merges:** se les asigna la tarea de la rama integrada, pero no suman a la autoría efectiva.
5. **Trabajo fuera del submódulo:** tareas de backend pueden haberse validado o completado en `test/back`, `docker/` o el servidor; esos cambios no siempre están en `backend/`.

---

_Informe generado automáticamente a partir de `git log` de `backend/` y de la documentación del repositorio. Los hashes cortos permiten localizar el commit completo con `git show <hash>`._
