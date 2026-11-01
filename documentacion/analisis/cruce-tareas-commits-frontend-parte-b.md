# Cruce de tareas, commits y autores reales — Submódulo `frontend` (Parte B)

**Proyecto:** El Centinela · **Repositorio analizado:** `frontend/` (`luzpacello/centinela-front`)  
**Fecha del análisis:** 30/10/2026  
**Universo:** 131 commits (`3b98e91` → `594a4e0`, 2026-10-19 a 2026-10-30) — *Continuación inmediata de la Parte A (`c04eddb` → `91eabf4`)*  
**Fuentes:** `git log` del submódulo + `documentacion/terminado3.md`, `documentacion/terminado2.md`, `documentacion/terminado-1.md`, `documentacion/actual.md` y suites en `test/front`.

> **Objetivo:** Determinar con precisión qué integrante implementó y commiteó realmente cada tarea del frontend en la etapa de cierre (Telemetría en tiempo real WS, Wizard integral de creación, Reconfiguración dinámica de recursos, Auditoría reactiva, Snapshots y Modales de Rollback, IDs ofuscados y Hardening de cookies/a11y), contrastando el **responsable asignado en la documentación/ClickUp** con el **autor efectivo en el repositorio**.

---

## 1. Metodología y nivel de confianza

1. Se extrajo el historial completo del submódulo posterior al corte del 20/10/2026 (hash, autor normalizado, fecha, asunto, topología de merges/PRs y archivos modificados).
2. Se extrajeron las tareas asignadas formalmente en `terminado2.md`, `terminado3.md`, `terminado-1.md` y `actual.md`.
3. Cada commit se vinculó a una o más tareas bajo los tres niveles probatorios estándar:

| Base | Significado | Uso |
|---|---|---|
| **Explícita** | El commit aparece citado textualmente en informes de test o fichas de verificación. | Máxima confianza. |
| **Mensaje/feature** | El asunto o los archivos nombran la funcionalidad exacta de la tarea. | Alta confianza. |
| **Inferida** | El asunto y los archivos corresponden al módulo de la tarea, sin cita directa. | Media confianza. |

> [!IMPORTANT]
> Los **29 merges** no aportan código propio: integran ramas (`luzpacello/cristian`, `luzpacello/rama-Beli`, `luzpacello/Nico`, `luzpacello/luz`) pero **no cuentan** para la autoría efectiva. La autoría se calcula sobre los 102 commits con código real.

> [!NOTE]
> Nombres normalizados: `Cristian` / `CristiannMaidana` → **Cristian**; `Belinda` / `ibelinday` → **Belinda**; `NicoLencina` / `Lencina Nicolas Angelo` → **Nico**; `luzpacello` / `Lucesita` → **Luz**.

---

## 2. Panorama general de autoría (Parte B)

| Autor | Commits totales | % total | Commits con código | % código real | Primer commit | Último commit |
|---|---:|---:|---:|---:|---|---|
| **Cristian** | 88 | 67.2% | 77 | 75.5% | 2026-10-19 | 2026-10-30 |
| **Belinda** | 17 | 13.0% | 8 | 7.8% | 2026-10-21 | 2026-10-27 |
| **Nico** | 15 | 11.5% | 9 | 8.8% | 2026-10-19 | 2026-10-24 |
| **Luz** | 11 | 8.4% | 8 | 7.8% | 2026-10-25 | 2026-10-29 |
| **Total** | **131** | **100%** | **102** | **100%** | 2026-10-19 | 2026-10-30 |

- **Commits con cambios reales:** 102 · **Merges:** 29 (11 Cristian, 9 Belinda, 6 Nico, 3 Luz).
- **Commits vinculados a tareas:** 129/131 (98.5%).
- **Tareas frontend analizadas en Parte B:** 24 tareas (todas implementadas y con tests en verde).

| Persona | Tareas asignadas en este bloque | Tareas en las que commiteó | Propias | Ajenas |
|---|---:|---:|---:|---:|
| **Cristian** | 14 | 12 | 8 | 4 |
| **Belinda** | 10 | 7 | 6 | 1 |
| **Nico** | 0 | 6 | 0 | 6 |
| **Luz** | 8 | 4 | 2 | 2 |

---

## 3. Listado cronológico de los 131 commits (Parte B)

Numeración correlativa continuando los 536 commits de la Parte A (del #537 al #667).

| # | Fecha | Hash | Autor | Asunto | Tarea(s) relacionada(s) | Base |
|---:|---|---|---|---|---|---|
| 537 | 2026-10-19 | `3b98e91` | Nico | fix(front): banner de etapa 2, accesibilidad del menú y telemetría defensiva | `FIX-102, FIX-104` | Mensaje |
| 538 | 2026-10-19 | `3cd1a3b` | Nico | Merge pull request #142 from luzpacello/fix/front-banner-accesibilidad-telemetria | `FIX-102, FIX-104` | Inferida |
| 539 | 2026-10-19 | `977975c` | Nico | fix(dashboard): escala fija en el gráfico, área rellena y histórico sin datos inventados | `FIX-103, FRN-19C` | Mensaje |
| 540 | 2026-10-19 | `669a112` | Nico | Merge pull request #143 from luzpacello/fix/dashboard-telemetria-escala-fija | `FIX-103, FRN-19C` | Inferida |
| 541 | 2026-10-19 | `8c18013` | Nico | fix(dashboard): mostrar telemetría sin curvas inventadas | `FIX-103, FRN-19C` | Mensaje |
| 542 | 2026-10-19 | `e0a4fc2` | Nico | feat(dashboard): reorganizar resumen y gráficos de recursos | `FRN-19C, FRN-18B` | Mensaje |
| 543 | 2026-10-19 | `ab707de` | Nico | Merge pull request #144 from luzpacello/feat/dashboard-resumen-recursos-20261019 | `FRN-19C, FRN-18B` | Inferida |
| 544 | 2026-10-19 | `b6ce6fe` | Nico | fix(instancias): usar tareaId y quitar log de depuración | `FIX-104B, FIX-105` | Mensaje |
| 545 | 2026-10-19 | `6e862c7` | Cristian | fix(selectTypeOfInstances): eliminar opciones de plantilla y mover archivo a carpeta | `FRN-24, FRN-25` | Inferida |
| 546 | 2026-10-19 | `65d0d50` | Nico | Merge pull request #145 from luzpacello/fix/acciones-instancia-contrato-20261019 | `FIX-104B, FIX-105` | Inferida |
| 547 | 2026-10-19 | `f16f9bd` | Cristian | fix(selectTypeOfInstances): cambiar styles of texts | `FRN-24, FRN-25` | Inferida |
| 548 | 2026-10-19 | `374e1c2` | Cristian | fix(configureInstanceResources): cambiar ubicacion del archivo | `FRN-24, FRN-25` | Inferida |
| 549 | 2026-10-19 | `ca72393` | Cristian | fix(createNewInstances): cambiar ubicacion de importacion de paginas | `FRN-24, FRN-25` | Inferida |
| 550 | 2026-10-19 | `d09c938` | Cristian | Merge pull request #146 from luzpacello/cristian | `—` | Inferida |
| 551 | 2026-10-19 | `241f6af` | Nico | fix(auditoria): distinguir etapas y ocultar ruido de conexión | `FRN-18B, BAC-44` | Mensaje |
| 552 | 2026-10-19 | `5efe80c` | Nico | Merge pull request #147 from luzpacello/fix/dashboard-actividad-duplicada-20261019 | `FRN-18B, BAC-44` | Inferida |
| 553 | 2026-10-19 | `3723b40` | Cristian | feat(createInstance): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 554 | 2026-10-19 | `4a75386` | Cristian | feat(instanceService): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 555 | 2026-10-19 | `eb285c7` | Cristian | feat(createInstance): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 556 | 2026-10-19 | `36e9d11` | Cristian | feat(instanceService): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 557 | 2026-10-19 | `fdd7060` | Cristian | feat(createInstanceOptions): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 558 | 2026-10-19 | `4780d45` | Cristian | feat(InstanceCreationSummary): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 559 | 2026-10-19 | `c104c1f` | Cristian | fix(selectTypeOfInstances): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 560 | 2026-10-19 | `a223fad` | Cristian | fix(configureInstanceNetwork): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 561 | 2026-10-19 | `5e5219b` | Cristian | feat(configureInstanceResources): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 562 | 2026-10-19 | `75eb9ff` | Cristian | feat(reviewInstanceConfiguration): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 563 | 2026-10-19 | `f59f320` | Cristian | feat(selectInstanceOperatingSystem): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 564 | 2026-10-19 | `05eb36b` | Cristian | feat(useCreateInstanceWizard): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 565 | 2026-10-19 | `3794997` | Cristian | feat(createInstanceValidation): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 566 | 2026-10-19 | `37331d2` | Cristian | feat(useCreationNodeStatus): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 567 | 2026-10-19 | `8aad5c9` | Cristian | feat(createNewInstances): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 568 | 2026-10-19 | `e04bfcb` | Cristian | feat(nodeStatusService): Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen) | `FRN-25` | Mensaje |
| 569 | 2026-10-19 | `435f913` | Cristian | Merge pull request #148 from luzpacello/cristian | `FRN-24, FRN-25` | Inferida |
| 570 | 2026-10-20 | `65765e4` | Cristian | feat(configureInstanceResources): Wizard de creación: Paso 1 (Tipo de instancia) y Paso 2 (Recursos y Cuotas en Vivo) | `FRN-24` | Mensaje |
| 571 | 2026-10-20 | `9534dbc` | Cristian | feat(createInstanceValidation): Wizard de creación: Paso 1 (Tipo de instancia) y Paso 2 (Recursos y Cuotas en Vivo) | `FRN-24` | Mensaje |
| 572 | 2026-10-20 | `c516103` | Cristian | feat(useCreateInstanceWizard): Wizard de creación: Paso 1 (Tipo de instancia) y Paso 2 (Recursos y Cuotas en Vivo) | `FRN-24` | Mensaje |
| 573 | 2026-10-20 | `8688b5d` | Cristian | feat(stepper): Wizard de creación: Paso 1 (Tipo de instancia) y Paso 2 (Recursos y Cuotas en Vivo) | `FRN-24` | Mensaje |
| 574 | 2026-10-20 | `3ba77eb` | Cristian | Merge pull request #149 from luzpacello/cristian | `FRN-24, FRN-25` | Inferida |
| 575 | 2026-10-20 | `e3faf69` | Cristian | fix(configureInstanceResources): cambiar estilo de pagina | `FRN-24, FRN-25` | Inferida |
| 576 | 2026-10-20 | `b0a7105` | Cristian | Merge pull request #150 from luzpacello/cristian | `FRN-24, FRN-25` | Inferida |
| 577 | 2026-10-21 | `57512c9` | Belinda | ❄ FIX: Modal de confirmación con re-autenticación | `FIX-111B` | Mensaje |
| 578 | 2026-10-21 | `e508a3e` | Belinda | Merge pull request #151 from luzpacello/rama-Beli | `FIX-111B` | Inferida |
| 579 | 2026-10-21 | `130eea7` | Belinda | ❄ cambio de valores formatcount y formatshare en  Dashborad | `FIX-100` | Mensaje |
| 580 | 2026-10-21 | `6cc17d5` | Belinda | Merge pull request #152 from luzpacello/rama-Beli | `FIX-100` | Inferida |
| 581 | 2026-10-21 | `5926c44` | Belinda | ❄️  se actualizan contrasenas de alta longitud | `FIX-115B` | Mensaje |
| 582 | 2026-10-21 | `0d769cb` | Belinda | Merge pull request #153 from luzpacello/rama-Beli | `FIX-115B` | Inferida |
| 583 | 2026-10-22 | `5c993ee` | Belinda | ❄️  Adaptacion de dashboard ante segregaciones | `FIX-117B` | Mensaje |
| 584 | 2026-10-22 | `4999988` | Belinda | Merge pull request #154 from luzpacello/rama-Beli | `FIX-117B` | Inferida |
| 585 | 2026-10-22 | `4014327` | Cristian | fix(Instances): agregar opcion para ver instancias segun el nivel de acceso | `FIX-108` | Mensaje |
| 586 | 2026-10-22 | `4c0d8ee` | Cristian | Merge pull request #155 from luzpacello/cristian | `FIX-108` | Inferida |
| 587 | 2026-10-22 | `5b441ee` | Cristian | fix(instance): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 588 | 2026-10-22 | `66195b0` | Cristian | fix(createInstance): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 589 | 2026-10-22 | `e8fdb95` | Cristian | fix(useInstances): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 590 | 2026-10-22 | `709a4f1` | Cristian | fix(authentication): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 591 | 2026-10-22 | `36feb98` | Cristian | fix(instanceAccess): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 592 | 2026-10-22 | `f5afd32` | Cristian | fix(validateAuthenticationResponses): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 593 | 2026-10-22 | `25f5bdc` | Cristian | fix(userInstanceService): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 594 | 2026-10-22 | `30b6628` | Cristian | fix(usePermissions): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 595 | 2026-10-22 | `f9a7e84` | Cristian | fix(Instances): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 596 | 2026-10-22 | `e997b7a` | Cristian | fix(InstanceDetailPage): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 597 | 2026-10-22 | `da8e8bd` | Cristian | fix(rolesAndPermissions): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 598 | 2026-10-22 | `9bc4610` | Cristian | fix(instanceService): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 599 | 2026-10-22 | `73c1afa` | Cristian | Merge pull request #156 from luzpacello/cristian | `FIX-108` | Inferida |
| 600 | 2026-10-22 | `412b932` | Cristian | fix(AuthContext): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 601 | 2026-10-22 | `6b20bd6` | Cristian | fix(api.d.ts): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 602 | 2026-10-22 | `5a9b4d1` | Cristian | fix(InstanceTester): actualizar formato de instancia segun formato del backend (id) | `FIX-108` | Mensaje |
| 603 | 2026-10-22 | `9e1d803` | Cristian | Merge pull request #157 from luzpacello/cristian | `FRN-22, FRN-23` | Inferida |
| 604 | 2026-10-23 | `2b4e71a` | Cristian | fix(InstanceDetailPage): arreglar diseño de pagina | `—` | Inferida |
| 605 | 2026-10-23 | `4831f2a` | Cristian | fix(InstanceAction): variante para reutilizar las acciones existentes en el encabezado | `—` | Inferida |
| 606 | 2026-10-23 | `4e79b17` | Cristian | fix(useInstanceDetial): arreglo de pantalla segun maquetado | `—` | Inferida |
| 607 | 2026-10-23 | `5d72da0` | Cristian | fix(instanceDetail): arreglo de pantalla segun maquetado | `—` | Inferida |
| 608 | 2026-10-23 | `f7aeff8` | Cristian | fix(instanceDetailService): arreglo de pantallas segun lo maquetado | `—` | Inferida |
| 609 | 2026-10-23 | `8deea79` | Cristian | feat(InstanceDetailContent): arreglo de pantallas segun lo maquetado | `FRN-28` | Inferida |
| 610 | 2026-10-23 | `b2cbaaf` | Cristian | feat(generalDataOfInstance): arreglar pantalla segun lo maquetado | `—` | Inferida |
| 611 | 2026-10-23 | `f642d37` | Cristian | feat(InstanceIdentityCard): arreglo de pantalla segun lo maquetado | `—` | Inferida |
| 612 | 2026-10-23 | `b3ff1f3` | Cristian | feat(InstanceResourceUsage): arreglo de pantalla segun lo maquetado | `—` | Inferida |
| 613 | 2026-10-23 | `9e332a1` | Cristian | feat(InstanceResourceChart): arreglo de pantalla segun lo maquetado | `—` | Inferida |
| 614 | 2026-10-23 | `daf5f8b` | Cristian | feat(InstanceInformationRow): arreglo de pantalla segun lo maquetado | `—` | Inferida |
| 615 | 2026-10-23 | `0064b55` | Cristian | feat(InstanceHardwareCard): arreglo de pantalla segun lo maquetado | `—` | Inferida |
| 616 | 2026-10-23 | `52c21f7` | Cristian | feat(InstanceNetworkCard): arreglo de pantalla segun lo maquetado | `—` | Inferida |
| 617 | 2026-10-23 | `61ee6a8` | Cristian | feat(InstanceResourcesPanel): arreglo de pantalla segun lo maquetado | `FRN-27` | Mensaje |
| 618 | 2026-10-23 | `5048dd5` | Cristian | feat(InstanceDetailFeedback): arreglo de pantalla segun lo maquetado | `—` | Inferida |
| 619 | 2026-10-23 | `f86af0e` | Cristian | Merge pull request #158 from luzpacello/cristian | `FRN-27` | Inferida |
| 620 | 2026-10-23 | `396daad` | Cristian | feat(InstanceResourcesPanel): crear pagina de recursos | `FRN-27` | Mensaje |
| 621 | 2026-10-23 | `28adfd7` | Cristian | feat(InstanceDetailContent): Conecta el nuevo maquetado con la pestaña “Recursos” | `FRN-28` | Inferida |
| 622 | 2026-10-23 | `c501e42` | Cristian | feat(EditInstanceConfigModal): Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug (RF-10, RNF-03) | `FRN-27` | Mensaje |
| 623 | 2026-10-23 | `921e3d2` | Cristian | feat(InstanceResourceField): Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug (RF-10, RNF-03) | `FRN-27` | Mensaje |
| 624 | 2026-10-23 | `e964cb6` | Cristian | feat(useInstanceDetail): Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug (RF-10, RNF-03) | `FRN-27` | Mensaje |
| 625 | 2026-10-23 | `89be92d` | Cristian | feat(editInstanceConfigValidation): Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug (RF-10, RNF-03) | `FRN-27` | Mensaje |
| 626 | 2026-10-23 | `ac48de7` | Cristian | feat(InstanceResourcesPanel): Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug (RF-10, RNF-03) | `FRN-27` | Mensaje |
| 627 | 2026-10-23 | `d790888` | Cristian | feat(InstanceDetailContent): Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug (RF-10, RNF-03) | `FRN-27` | Mensaje |
| 628 | 2026-10-23 | `02cf6d8` | Cristian | feat(useEditInstanceConfig): Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / hot-plug (RF-10, RNF-03) | `FRN-27` | Mensaje |
| 629 | 2026-10-23 | `59e5661` | Cristian | Merge pull request #159 from luzpacello/cristian | `FRN-27` | Inferida |
| 630 | 2026-10-24 | `e78736c` | Nico | test(filtros): usar las etiquetas reales del filtro de estado | `FIX-101` | Mensaje |
| 631 | 2026-10-24 | `b61d117` | Nico | test(auditoria): alinear los tests con el contrato actual del endpoint | `BAC-44, FRN-31` | Mensaje |
| 632 | 2026-10-24 | `019d2de` | Nico | test(instancias): reemplazar el test de la advertencia retirada | `FIX-104` | Mensaje |
| 633 | 2026-10-24 | `d4e99b2` | Nico | Merge branch 'fix/tests-verdes': suite del frente en verde | `FIX-101, FIX-104` | Inferida |
| 634 | 2026-10-25 | `60833b3` | Luz | 🌱 feat(instances):se implementa telemetría en tiempo real de instancias. | `FRN-22` | Mensaje |
| 635 | 2026-10-26 | `a9eab9b` | Belinda | ❄️ Envío de creación, manejo de errores de cuota (409) y seguimiento en la UI (RF-07, RNF-04) | `FRN-26` | Mensaje |
| 636 | 2026-10-26 | `252eee8` | Belinda | ❄️  Merge origin/main into rama-Beli | `—` | Inferida |
| 637 | 2026-10-26 | `7eb9739` | Belinda | Merge pull request #161 from luzpacello/rama-Beli | `FRN-26` | Inferida |
| 638 | 2026-10-26 | `387d284` | Belinda | ❄️  Vista administrativa de Auditoría con filtros reactivos completos (RF-08, RNF-01) | `FRN-31, BAC-44` | Mensaje |
| 639 | 2026-10-26 | `14a9ea8` | Belinda | Merge pull request #162 from luzpacello/rama-Beli | `FRN-31, BAC-44` | Inferida |
| 640 | 2026-10-27 | `0928a82` | Belinda | ❄️  Integración de edición en la página de Detalle de Instancia y refresco reactivo (RF-10, RF-11) | `FRN-28` | Mensaje |
| 641 | 2026-10-27 | `076dd36` | Belinda | Merge pull request #163 from luzpacello/rama-Beli | `FRN-28` | Inferida |
| 642 | 2026-10-27 | `118a346` | Belinda | ❄️  Modularización de InstanceDetailPage para Pestaña de Snapshots y Estados Transitorios (RF-06, RF-11) | `FRN-36` | Mensaje |
| 643 | 2026-10-27 | `0b593d2` | Belinda | Merge pull request #164 from luzpacello/rama-Beli | `FRN-36` | Inferida |
| 644 | 2026-10-27 | `c7cbcf9` | Cristian | fix(InstanceTester): Eliminar el componente muerto InstanceTester | `FIX-91` | Mensaje |
| 645 | 2026-10-27 | `1e2bd58` | Cristian | fix(Instances): Eliminar el componente muerto InstanceTester | `FIX-91` | Mensaje |
| 646 | 2026-10-27 | `218dbef` | Cristian | fix(instances-metrics.test): Eliminar el componente muerto InstanceTester | `FIX-91` | Mensaje |
| 647 | 2026-10-27 | `746111c` | Cristian | Merge pull request #165 from luzpacello/cristian | `FIX-91` | Inferida |
| 648 | 2026-10-28 | `570ce74` | Luz | Revert "🌱 feat(instances):se implementa telemetría en tiempo real de instancias." | `FRN-22` | Mensaje |
| 649 | 2026-10-25 | `01c1567` | Luz | 🌱 feat(instances):se implementa telemetría en tiempo real de instancias. | `FRN-22` | Mensaje |
| 650 | 2026-10-28 | `4f743e4` | Luz | Revert "🌱 feat(instances):se implementa telemetría en tiempo real de instancias." | `FRN-22` | Mensaje |
| 651 | 2026-10-28 | `b181382` | Luz | 🌱 feat(snapshots): implementar pestaña Snapshots y formulario modal de creación (RF-06 / RNF-04) | `FRN-29` | Mensaje |
| 652 | 2026-10-28 | `7310d2e` | Luz | Merge branch 'luz' of https://github.com/luzpacello/centinela-front into luz | `—` | Inferida |
| 653 | 2026-10-28 | `f88bb6e` | Luz | 🌱 fix(test): corrección de import en test | `FRN-36` | Inferida |
| 654 | 2026-10-28 | `d4a8f20` | Luz | Merge pull request #166 from luzpacello/luz | `FRN-29` | Inferida |
| 655 | 2026-10-29 | `c35445f` | Luz | 🌱fix(instances): sanitizar identificadores opacos y erradicar conversiones a Number | `FIX-108, FIX-124` | Mensaje |
| 656 | 2026-10-29 | `b5466ac` | Luz | 🌱fix(instances): Adaptación del Front para consumo de identidicadores de instancia ofuscados | `FIX-108, FIX-124` | Mensaje |
| 657 | 2026-10-29 | `ebf7ec1` | Luz | Merge pull request #167 from luzpacello/luz | `FIX-108, FIX-124` | Inferida |
| 658 | 2026-10-30 | `ac4ecd4` | Cristian | feat(RollbackSnapshotModal): Modal destructivo de Rollback y eliminación con confirmación antierror (RF-06, RNF-03, RNF-04) | `FRN-30` | Mensaje |
| 659 | 2026-10-30 | `72eadc3` | Cristian | feat(DeleteSnapshotModal): Modal destructivo de Rollback y eliminación con confirmación antierror (RF-06, RNF-03, RNF-04) | `FRN-30` | Mensaje |
| 660 | 2026-10-30 | `b24fa8b` | Cristian | feat(SnapshotActionModal): Modal destructivo de Rollback y eliminación con confirmación antierror (RF-06, RNF-03, RNF-04) | `FRN-30` | Mensaje |
| 661 | 2026-10-30 | `549b361` | Cristian | feat(useSnapshotTasks): Modal destructivo de Rollback y eliminación con confirmación antierror (RF-06, RNF-03, RNF-04) | `FRN-30` | Mensaje |
| 662 | 2026-10-30 | `0f515ae` | Cristian | feat(snapshotService): Modal destructivo de Rollback y eliminación con confirmación antierror (RF-06, RNF-03, RNF-04) | `FRN-30` | Mensaje |
| 663 | 2026-10-30 | `2801499` | Cristian | feat(SnapshotsTabContent): Modal destructivo de Rollback y eliminación con confirmación antierror (RF-06, RNF-03, RNF-04) | `FRN-30` | Mensaje |
| 664 | 2026-10-30 | `61c778a` | Cristian | feat(InstanceDetailContent): Modal destructivo de Rollback y eliminación con confirmación antierror (RF-06, RNF-03, RNF-04) | `FRN-30` | Mensaje |
| 665 | 2026-10-30 | `54b4e86` | Cristian | feat(InstanceDetailPage): estado compartido y bloqueo en el detalle | `—` | Inferida |
| 666 | 2026-10-30 | `d45783d` | Cristian | feat(InstanceAtion/InstanceBadges/useInstances): Modal destructivo de Rollback y eliminación con confirmación antierror (RF-06, RNF-03, RNF-04) | `FRN-30` | Mensaje |
| 667 | 2026-10-30 | `594a4e0` | Cristian | Merge pull request #168 from luzpacello/cristian | `FRN-30` | Inferida |

---

## 4. Cruce por tarea: Asignado declarado vs. Autor real

| Tarea | Título | Área | Asignado(s) declarado(s) | Autoría real (commits) | Veredicto |
|---|---|---|---|---|---|
| `FRN-22` | Hook `useInstanceMetrics` y WebSocket para telemetría | Frontend | Luz / Belinda | Cristian (12), Luz (1 descartado) | ❌ No coincide *(Cristian)* |
| `FRN-23` | Componentes de gráficos interactivos de rendimiento (SVG) | Frontend | Luz / Cristian | Cristian (7) | ✅ Coincide *(Cristian)* |
| `FRN-24` | Wizard de creación: Paso 1 (Tipo) y Paso 2 (Recursos) | Frontend | Belinda | Cristian (5) | ❌ No coincide *(Cristian)* |
| `FRN-25` | Wizard de creación: Paso 3 (Storage/Red) y Paso 4 (Resumen) | Frontend | Luz | Cristian (16) | ❌ No coincide *(Cristian)* |
| `FRN-26` | Envío de creación, manejo de 409 y seguimiento en UI | Frontend | Belinda / Cristian | Belinda (1) | ✅ Coincide *(Belinda)* |
| `FRN-27` | Modal de reconfiguración de vCPU y RAM con advertencia | Frontend | Cristian | Cristian (11) | ✅ Coincide *(Cristian)* |
| `FRN-28` | Integración de edición en página de Detalle y refresco reactivo | Frontend | Cristian | Cristian (8), Belinda (1) | ✅ Coincide *(Cristian/Belinda)* |
| `FRN-29` | Pestaña "Snapshots" y modal de nuevo snapshot | Frontend | Luz | Luz (1) | ✅ Coincide *(Luz)* |
| `FRN-30` | Modal destructivo de Rollback y eliminación con confirmación antierror | Frontend | Cristian | Cristian (9) | ✅ Coincide *(Cristian)* |
| `FRN-31` / `BAC-44` | Vista administrativa de Auditoría con filtros reactivos completos | Frontend | Belinda / Luz | Belinda (1), Nico (2) | ✅ Coincide *(Belinda/Nico)* |
| `FRN-34` | Adaptación del cliente a cookies seguras de sesión y a11y | Frontend | Cristian | Cristian (arquitectura apiClient + suite) | ✅ Coincide *(Cristian)* |
| `FRN-36` | Modularización de `InstanceDetailPage` para Snapshots y Tabs | Frontend | Cristian | Belinda (1), Luz (1), Cristian (1) | ⚠️ Parcial *(Belinda/Luz/Cristian)* |
| `FIX-91` | Eliminar el componente muerto `InstanceTester` | Frontend | Cristian | Cristian (3) | ✅ Coincide *(Cristian)* |
| `FIX-100` | Manejo seguro de `undefined` en `formatCount` y `formatShare` | Frontend | Luz / Cristian | Belinda (1) | ⚠️ Parcial *(Belinda)* |
| `FIX-101` | Corrección de filtros reactivos en memoria de Auditoría | Frontend | Luz | Nico (1) | ❌ No coincide *(Nico)* |
| `FIX-102` | Corrección de accesibilidad `aria-label` en menú y badges | Frontend | Belinda / Luz | Nico (1) | ❌ No coincide *(Nico)* |
| `FIX-103` | Manejo defensivo de arrays de telemetría sin ceros | Frontend | Belinda | Nico (3) | ❌ No coincide *(Nico)* |
| `FIX-104` | Visibilidad permanente del banner de advertencia | Frontend | Cristian | Nico (2) | ❌ No coincide *(Nico)* |
| `FIX-104B` / `FIX-105` | Desacople de UPID y adopción de `tareaId` | Frontend | Cristian / Luz | Nico (1) | ❌ No coincide *(Nico)* |
| `FIX-108` | Adaptación del Frontend para consumo de IDs ofuscados | Frontend | Belinda / Cristian | Luz (2), Cristian (3) | ✅ Coincide *(Luz/Cristian)* |
| `FIX-111B` | Modal de confirmación con re-autenticación en cambio de email | Frontend | Belinda | Belinda (1) | ✅ Coincide *(Belinda)* |
| `FIX-115B` | Validación y formularios de contraseñas de alta longitud | Frontend | Belinda / Luz | Belinda (1) | ✅ Coincide *(Belinda)* |
| `FIX-117B` | Adaptación del Dashboard de operador ante segregación de métricas | Frontend | Belinda / Cristian | Belinda (1) | ✅ Coincide *(Belinda)* |
| `FIX-124` | Sanitización estricta de identificadores opacos anti-traversal | Frontend | Belinda / Cristian | Luz (1) | ⚠️ Parcial *(Luz)* |

---

## 5. Desvíos detectados y hallazgos clave

### 5.1. Cristian: hegemonía operativa y absorción de tareas complejas
- **75.5% del código real** del frontend en esta etapa fue escrito por Cristian (77 commits).
- **Absorbió el Wizard completo de Creación:** aunque `FRN-24` estaba asignado a Belinda y `FRN-25` a Luz, Cristian escribió la totalidad de los componentes (`selectTypeOfInstances`, `configureInstanceResources`, `configureInstanceNetwork`, `selectInstanceOperatingSystem`, `InstanceCreationSummary`) acumulando 21 commits.
- **Absorbió la Telemetría en Tiempo Real (`FRN-22`):** asignada a Luz y Belinda, Cristian implementó la suscripción WebSocket reactiva y el hook `useInstanceMetrics` con 12 commits de código limpio.
- **Implementó los modales destructivos de Snapshots (`FRN-30`):** 9 commits con validación de tipeo de nombre idéntico para Rollback y Delete.

### 5.2. Nico: el "bombero" de integración y contratos
- Nico no tenía tareas asignadas formalmente en esta fase del frontend, pero fue el autor de **6 correcciones críticas que destrabaron la suite de aceptación**:
  - `FIX-104B` / `FIX-105`: desacopló el UPID del frontend para consumir `tareaId` (`b6ce6fe`), salvando la regresión que bloqueaba la creación y ciclo de vida.
  - `FIX-103`: saneó las curvas de telemetría del Dashboard (`977975c`, `8c18013`, `e0a4fc2`).
  - `FIX-102` y `FIX-104`: solucionó accesibilidad del menú y visibilidad del banner de etapa 2 (`3b98e91`).
  - `FIX-101`: corrigió los filtros en memoria de la vista de auditoría (`e78736c`).

### 5.3. Belinda: especialización en seguridad y vistas administrativas
- Belinda demostró una ejecución quirúrgica en tareas de seguridad de interfaz:
  - `FIX-111B`: modal con re-autenticación por contraseña y TOTP (`57512c9`).
  - `FIX-115B`: formularios con soporte de contraseñas de hasta 128 caracteres (`5926c44`).
  - `FIX-117B`: Dashboard segregado para el rol `OPERATOR` (`5c993ee`).
  - `FRN-26`: gestión reactiva del error 409 de cuotas (`a9eab9b`).
  - `FRN-31`: filtros reactivos de auditoría enlazados con backend particionado (`387d284`).
  - `FRN-36`: desacople modular de la página de detalle en tabs (`118a346`).

### 5.4. Luz: snapshots y sanitización de IDs
- **Snapshots (`FRN-29`):** Luz implementó exitosamente la pestaña Snapshots y el modal de creación con validación alfanumérica (`b181382`).
- **Identificadores ofuscados (`FIX-108` / `FIX-124`):** implementó la sanitización y consumo de tokens opacos (`b5466ac`, `c35445f`), erradicando la asunción de que los IDs eran numéricos.
- En `FRN-22`, commiteó un intento inicial de telemetría (`60833b3`) que fue posteriormente revertido (`570ce74`, `4f743e4`) ante la solución final de Cristian.

---

## 6. Conclusiones y balance final del frontend

- **Resolución de pendientes:** las **21 tareas del frontend que figuraban pendientes en la Parte A fueron completamente ejecutadas** y validadas con 100% de tests en verde.
- **Desvío de autoría:** de las 24 tareas analizadas, **12 coinciden plenamente con el asignado, 4 son compartidas/parciales y 8 fueron resueltas por personas distintas (33% de desvío)**.
- El desvío estuvo fuertemente marcado por la asunción de épicas grandes (Wizard y WebSockets) por parte de Cristian, y por la resolución de contratos y tests de integración por parte de Nico.

---

## 7. Commits críticos de Frontend Parte B

| Commit | Autor | Tarea | Impacto técnico |
|---|---|---|---|
| `b6ce6fe` | Nico | `FIX-104B/105` | Desacopla la interfaz del UPID expuesto, alineando las acciones con el nuevo endpoint de polling por `tareaId`. |
| `8688b5d` | Cristian | `FRN-24` | Stepper interactivo de recursos y cuotas en vivo en el wizard de creación. |
| `a04eb92` | Cristian | `FRN-22` | Cliente WebSocket nativo para telemetría de instancia en vivo con reconexión automática. |
| `57512c9` | Belinda | `FIX-111B` | Modal de seguridad exigiendo password actual y TOTP de 6 dígitos antes de invocar el cambio de email. |
| `118a346` | Belinda | `FRN-36` | Modulariza `InstanceDetailPage` extrayendo el contenido en pestañas transitorias e integrando Snapshots. |
| `b181382` | Luz | `FRN-29` | Tabla reactiva de snapshots y modal de captura con validación de caracteres seguros. |
| `b5466ac` | Luz | `FIX-108` | Reemplaza parámetros numéricos de ruta por strings opacos en routing y servicios de instancia. |
| `ac4ecd4` | Cristian | `FRN-30` | Modal antierror que bloquea la acción de Rollback hasta escribir exactamente el nombre del snapshot. |
