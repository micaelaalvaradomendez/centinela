# Mapa de Evidencia Empírica y Plan de Integración en el Artículo

**Artículo de Referencia:** *"El Project Manager en la Era de la IA: De Administrador a Estratega"*  
**Autora:** Micaela Alvarado Mendez (Universidad Nacional de Tierra del Fuego, Ushuaia, Argentina)  
**Proyecto Empírico:** *El Centinela* (Plataforma de orquestación y gestión de entornos virtualizados sobre Proxmox VE)  
**Fecha de Elaboración Original:** 30 de Octubre de 2026  
**Última Actualización Sincronizada con el Repositorio:** 3 de Noviembre de 2026 (Commit Orquestador `d63840e`, Backend `fb8a41c`, Frontend `877eddd`)

---

## 1. Diagnóstico del Artículo Original

El artículo redactado por Micaela Alvarado Mendez plantea una tesis sólida y de gran actualidad en la ingeniería de software contemporánea:
1. **Desplazamiento del rol operativo:** La Inteligencia Artificial Generativa y los LLMs están absorbiendo las tareas administrativas de bajo valor agregado (minutas, redacción de requerimientos, seguimiento superficial de tableros, estimaciones empíricas), reduciendo la fricción burocrática reportada por la literatura (Barbosa et al., 2025).
2. **Emergencia del AI Manager y Habilitador Estratégico:** El Project Manager evoluciona hacia un perfil enfocado en la **ingeniería de contexto** (*context engineering*), la supervisión continua *human-in-the-loop* y la gestión sociotécnica (López Hurtado, 2026).
3. **El *Judgment Gap* como núcleo del liderazgo:** Apoyándose en Herbert Simon (1978) y Antonio Nieto-Rodríguez (2026), se demuestra que la abundancia de métricas y automatizaciones no resuelve la toma de decisiones, sino que acentúa la brecha de juicio: la IA provee visibilidad analítica (*decision gap*), pero el juicio humano asume las consecuencias éticas, técnicas y de negocio (*judgment gap*).

### La debilidad identificada en el artículo original y precisión teórica
En la **Sección VII.D ("Limitaciones y agenda de investigación")**, el manuscrito original reconoce explícitamente:
> *"Segundo, el análisis es fundamentalmente conceptual: no incluye evidencia empírica primaria sobre la implementación del AI Manager en organizaciones latinoamericanas. Tercero, las tensiones identificadas en la Sección VI son hipótesis derivadas de la literatura, no hallazgos validados empíricamente en el contexto regional."*

Asimismo, en la **Sección V.C ("Verificación formal contra el código")**, el artículo utiliza un ejemplo hipotético:
> *"El PM no necesita revisar el código línea por línea: recibe un reporte sintético que dice 'PR #1234: divergencia en campo currency; falta manejo de error 403'. Con esa información, negocia con el equipo de backend..."*

Adicionalmente, tal como se documentó en [`guia-estudio.md`](./guia-estudio.md), en la **Sección III.C** conviene precisar la cita de López Hurtado (2026) respecto al **NIST AI RMF 1.0**: mientras dicha fuente menciona cuatro pilares genéricos (*gobernanza, validación, gestión del riesgo y supervisión humana*), la estructura oficial del núcleo del NIST AI RMF 1.0 se organiza en cuatro funciones: **Gobernar (*Govern*), Mapear (*Map*), Medir (*Measure*) y Gestionar (*Manage*)**.

### La oportunidad de mejora
El proyecto real **El Centinela**, desarrollado y gestionado por la autora en la cátedra de Ingeniería de Software III, proporciona **el 100% de la evidencia empírica primaria necesaria** para transformar el artículo de un ensayo conceptual a un **estudio de caso sociotécnico con respaldo cuantitativo y trazabilidad en Git**, cubriendo desde la Fase Base y las Etapas 1, 2 y 3 hasta la **Etapa de Adaptación a Proxmox VE 8.4.21** y la auditoría de despliegue en el servidor físico final (`pve-is3-g1`).

---

## 2. Mapa de Integración: Del Artículo a la Evidencia Real Actualizada

La siguiente matriz detalla en qué sección del artículo debe incorporarse cada componente de la evidencia documentada en los archivos de esta carpeta:

| Sección del Artículo Original | Contenido Original | Evidencia Empírica de *El Centinela* a Incorporar (Estado al 03/11/2026) | Archivo Fuente Generado |
|---|---|---|---|
| **Sección II: De la tarea administrativa a la ingeniería de contexto** | Discusión sobre tableros, historias de usuario y prompts de contexto. | Arquitectura del repositorio orquestador, gestión de submódulos (`backend` y `frontend`) con punteros inmutables, desacople *Local-First* con Docker Compose, simulador dual de Proxmox VE (9.2.2 y 8.4.21) y topología real LXC de 3 capas (`vmbr1`, `vmbr2`, `vmbr3`) con bastión Tailscale (`centi-mica`). | [`01_estructura_repo_readme_submodulos.md`](./01_estructura_repo_readme_submodulos.md) |
| **Sección III: El paradigma del AI Manager en el ciclo de entrega** | Supervisión continua, generación de especificaciones y gobernanza algorítmica. | Configuración de CI (`integration-tests.yml`), automatización por `repository_dispatch`, taxonomía atómica en ClickUp (`BAC`, `FRN`, `INF`, `SEC`, `BRG`, `ADP`) y creación sistemática de **140 tickets de corrección (`FIX-01` a `FIX-140`)**. | [`02_pipeline_ci_generacion_tareas_y_fixes.md`](./02_pipeline_ci_generacion_tareas_y_fixes.md) |
| **Sección V.A y V.B: Identificación de cuellos de botella y tests de aceptación** | Detección de bloqueos mediante métricas algorítmicas y arneses de pruebas. | Reporte de suites independientes (`/test`), orquestación cruzada en Go de tests E2E del Frontend (`int02-e2e.test.ts` y `login04-e2e.test.ts`), verificación de los 12 Casos de Uso (`CU-01..12`) y nuevas suites TDD de Etapa 3 y PVE 8.4 (`etapa3_fixes_pendientes_test.go`, `pve8_adaptation_acceptance_test.go`, `etapa3-pendientes.test.tsx`). | [`03_reportes_suite_tests_integracion_front_back.md`](./03_reportes_suite_tests_integracion_front_back.md) |
| **Sección V.C: Verificación contra criterios de aceptación y contratos** | Tabla de tipos de auditoría automatizada y el ejemplo hipotético del "PR #1234". | Sustitución del "PR #1234" por episodios reales: discrepancia de código HTTP en `BAC-24B` (`409 INSTANCE_NOT_STOPPED` vs `INSTANCE_INVALID_STATE`), validación de identificadores en `FIX-108`/`FIX-124`, rechazo de tráfico *untagged* `vlan: 0` en `FIX-128`, y quema indebida de OTP ante `400 SAME_PASSWORD` en `FIX-132`. | [`06_episodio_divergencia_real_vs_pr1234.md`](./06_episodio_divergencia_real_vs_pr1234.md) |
| **Sección V.C y Cuadro XII: Auditorías de seguridad sobre el código e infraestructura** | Mención teórica de patrones inseguros, secretos expuestos y contratos. | Auditoría de seguridad estructural: 9 vulnerabilidades críticas en aplicación (CWE-287, CWE-307, CWE-284, CWE-613, CWE-521, CWE-88, CWE-200, CWE-798, CWE-285; `FIX-111` a `FIX-119`), blindaje anti-IDOR (`FIX-105..110`), más hallazgos de seguridad de red en el servidor final (`centinela-redis` expuesto en `vmbr1` y egreso SMTP acotado con `iptables` hacia Brevo). | [`04_reportes_auditoria_seguridad.md`](./04_reportes_auditoria_seguridad.md) |
| **Sección V.D y VI.A: Desbalance de productividad y métricas concretas** | Citas de Barbosa et al. sobre 70% de eficiencia y brecha senior/junior. | Métricas cuantitativas consolidadas al 03/11/2026: **276 micro-tareas estructuradas**, **1.082 commits auditados** (`101` en orquestador, `273` en backend, `708` en frontend), **140 defectos y desvíos `FIX` interceptados**, **72.2% de tareas con fallas pre-merge** en el corte del 24/10, y desvío de autoría real vs. tablero (50% en backend, 33.3% en frontend). | [`05_metricas_y_datos_cuantitativos.md`](./05_metricas_y_datos_cuantitativos.md) |
| **Sección IV y VI.B-D: Superación del *Judgment Gap* y límites de la IA** | Teoría sobre el juicio humano, la asunción de riesgos y la supervisión ética. | **8 casos empíricos** donde la automatización o el diseño estático fallaron y requirieron juicio humano: rigidez sintáctica (`BAC-18B`), colisión de puerto (`INF-06A`), crash en frío de migraciones (`FIX-81`), bloqueo SMTP/NAT (`INF-08B` y `SMTP-y-Brevo.md`), falsa granularidad RBAC en UI, bloqueo de *CPU Overcommit* (`FIX-133`), deadlock en estado `paused` (`FIX-134`) y falsos positivos de caída ante tareas directas en Proxmox (`FIX-135`). | [`07_casos_fallo_automatizacion_y_judgment_gap.md`](./07_casos_fallo_automatizacion_y_judgment_gap.md) |
| **Sección VII.D: Limitaciones y agenda de investigación** | Declaración de falta de validación empírica. | Reescritura de las limitaciones: presentar a *El Centinela* como caso de estudio empírico reproducible que incorpora la verificación automatizada en la *Definition of Done* (DoD), reconociendo la limitación de ser un caso único académico y abriendo la agenda a mediciones multi-equipo. | [`00_analisis_articulo_y_mapa_evidencia.md`](./00_analisis_articulo_y_mapa_evidencia.md) |

---

## 3. Estructura de los Archivos Entregados y Actualizados

Para facilitar su lectura, defensa oral y consulta directa, la documentación de evidencia empírica se organiza en los siguientes 7 documentos especializados (acompañados por [`guia-estudio.md`](./guia-estudio.md) y [`presentacion.html`](./presentacion.html)):

1. **[`01_estructura_repo_readme_submodulos.md`](./01_estructura_repo_readme_submodulos.md):** Organización del repositorio orquestador, punteros fijados (`detached HEAD`), `.gitmodules`, arquitectura *Local-First* con simulador dual PVE 9.2.2 / 8.4.21 y topología real de despliegue en contenedores LXC sobre el nodo físico `pve-is3-g1`.
2. **[`02_pipeline_ci_generacion_tareas_y_fixes.md`](./02_pipeline_ci_generacion_tareas_y_fixes.md):** GitHub Actions, `repository_dispatch`, taxonomía completa de tareas (`BAC`, `FRN`, `INF`, `SEC`, `BRG`, `ADP`) y evolución del algoritmo de creación de tickets desde `FIX-01` hasta `FIX-140`.
3. **[`03_reportes_suite_tests_integracion_front_back.md`](./03_reportes_suite_tests_integracion_front_back.md):** Arnés independiente `/test`, ejecución cruzada Go $\leftrightarrow$ Vitest (`INT-02`, `LOGIN-04`), auditoría de los 12 Casos de Uso y estado actual de las suites de Etapa 3 y Adaptación a Proxmox VE 8.4.21 al 03/11/2026.
4. **[`04_reportes_auditoria_seguridad.md`](./04_reportes_auditoria_seguridad.md):** Diagnóstico de seguridad, falsedad de la matriz de permisos en UI, análisis de las 9 brechas CWE (`FIX-111..119`), blindaje anti-IDOR (`FIX-105..110`, `FIX-124`), ampliación NIST SP 800-63B (`FIX-138`) y auditoría de red e infraestructura del servidor final (`vmbr1..3`, `centi-mica` y egreso SMTP Brevo).
5. **[`05_metricas_y_datos_cuantitativos.md`](./05_metricas_y_datos_cuantitativos.md):** Cifras actualizadas al 03/11/2026 de tareas (276), commits (1.082), tasas de rechazo en cortes de validación, desvío de autoría real vs. ClickUp y estimación de retrabajo evitado.
6. **[`06_episodio_divergencia_real_vs_pr1234.md`](./06_episodio_divergencia_real_vs_pr1234.md):** Crónica técnica detallada de las 4 divergencias reales que sustituyen al ejemplo ficticio del PR #1234 (`BAC-24B`, `FIX-108`/`FIX-124`, `FIX-128` VLAN 0 y `FIX-132` OTP en recuperación de clave).
7. **[`07_casos_fallo_automatizacion_y_judgment_gap.md`](./07_casos_fallo_automatizacion_y_judgment_gap.md):** Análisis crítico de los 8 episodios empíricos (5 de pipeline/entorno y 3 de lógica de dominio de virtualización en Etapa 3) que demuestran por qué las reglas automatizadas sin contexto no reemplazan el juicio humano del PM.
