# Mapa de Evidencia Empírica y Plan de Integración en el Artículo

**Artículo de Referencia:** *"El Project Manager en la Era de la IA: De Administrador a Estratega"*  
**Autora:** Micaela Alvarado Mendez (Universidad Nacional de Tierra del Fuego, Ushuaia, Argentina)  
**Proyecto Empírico:** *El Centinela* (Plataforma de orquestación y gestión de entornos virtualizados sobre Proxmox VE)  
**Fecha de Elaboración:** Octubre de 2026

---

## 1. Diagnóstico del Artículo Original

El artículo redactado por Micaela Alvarado Mendez plantea una tesis sólida y de gran actualidad en la ingeniería de software contemporánea:
1. **Desplazamiento del rol operativo:** La Inteligencia Artificial Generativa y los LLMs están absorbiendo las tareas administrativas de bajo valor agregado (minutas, redacción de requerimientos, seguimiento superficial de tableros, estimaciones empíricas), reduciendo la fricción burocrática reportada por la literatura (Barbosa et al., 2025).
2. **Emergencia del AI Manager y Habilitador Estratégico:** El Project Manager evoluciona hacia un perfil enfocado en la **ingeniería de contexto** (*context engineering*), la supervisión continua *human-in-the-loop* y la gestión sociotécnica (López Hurtado, 2026).
3. **El *Judgment Gap* como núcleo del liderazgo:** Apoyándose en Herbert Simon (1978) y Antonio Nieto-Rodriguez (2026), se demuestra que la abundancia de métricas y automatizaciones no resuelve la toma de decisiones, sino que acentúa la brecha de juicio: la IA provee visibilidad analítica, pero el juicio humano asume las consecuencias éticas, técnicas y de negocio.

### La debilidad identificada en el artículo
En la **Sección VII.D ("Limitaciones y agenda de investigación")**, el propio artículo reconoce explícitamente:
> *"Segundo, el análisis es fundamentalmente conceptual: no incluye evidencia empírica primaria sobre la implementación del AI Manager en organizaciones latinoamericanas. Tercero, las tensiones identificadas en la Sección VI son hipótesis derivadas de la literatura, no hallazgos validados empíricamente en el contexto regional."*

Asimismo, en la **Sección V.C ("Verificación formal contra el código")**, el artículo utiliza un ejemplo hipotético:
> *"El PM no necesita revisar el código línea por línea: recibe un reporte sintético que dice 'PR #1234: divergencia en campo currency; falta manejo de error 403'. Con esa información, negocia con el equipo de backend..."*

### La oportunidad de mejora
El proyecto real **El Centinela** desarrollado y gestionado por la autora proporciona **el 100% de la evidencia empírica primaria necesaria** para transformar el artículo de un ensayo conceptual a un **estudio de caso sociotécnico de primer nivel científico**, apto para publicación en revistas y conferencias internacionales indexadas (IEEE / ACM / Springer).

---

## 2. Mapa de Integración: Del Artículo a la Evidencia Real

La siguiente matriz detalla en qué sección del artículo debe incorporarse cada componente de la evidencia documentada en los archivos de esta carpeta:

| Sección del Artículo Original | Contenido Original | Evidencia Empírica de *El Centinela* a Incorporar | Archivo Fuente Generado |
|---|---|---|---|
| **Sección II: De la tarea administrativa a la ingeniería de contexto** | Discusión sobre tableros, historias de usuario y prompts de contexto. | Arquitectura del repositorio orquestador, gestión de submódulos (`backend` y `frontend`), desacople *Local-First* con Docker Compose y contratos OpenAPI/SSE. | `01_estructura_repo_readme_submodulos.md` |
| **Sección III: El paradigma del AI Manager en el ciclo de entrega** | Supervisión continua, generación de especificaciones y gobernanza algorítmica. | Configuración de CI (`integration-tests.yml`), scripts de orquestación, desglose atómico de micro-tareas en ClickUp (0.5 h a 2.5 h) y ciclo de creación sistemática de `FIX-01` a `FIX-130`. | `02_pipeline_ci_generacion_tareas_y_fixes.md` |
| **Sección V.A y V.B: Identificación de cuellos de botella y tests de aceptación** | Detección de bloqueos mediante métricas algorítmicas y arneses de pruebas. | Reporte de suites de pruebas independientes (`/test`), orquestación cruzada en Go de tests E2E del Frontend (`int02-e2e.test.ts` y `login04-e2e.test.ts`) con Cookie Jar y WebSocket. | `03_reportes_suite_tests_integracion_front_back.md` |
| **Sección V.C: Verificación contra criterios de aceptación y contratos** | Tabla de tipos de auditoría automatizada y el ejemplo hipotético del "PR #1234". | Sustitución del "PR #1234" por los episodios reales: discrepancia de código HTTP en `BAC-24B` (`409 INSTANCE_NOT_STOPPED` vs `INSTANCE_INVALID_STATE`) y validación de identificadores en `FIX-108`. | `06_episodio_divergencia_real_vs_pr1234.md` |
| **Sección V.C y Cuadro XII: Auditorías de seguridad sobre el código** | Mención teórica de patrones inseguros, secretos expuestos y contratos. | Auditorías de seguridad estructural: 9 vulnerabilidades críticas (CWE-287, CWE-307, CWE-284, CWE-613, CWE-88, IDOR de VMIDs) y plan de mitigación `FIX-111` a `FIX-119`. | `04_reportes_auditoria_seguridad.md` |
| **Sección V.D y VI.A: Desbalance de productividad y métricas concretas** | Citas de Barbosa et al. sobre 70% de eficiencia y brecha senior/junior. | Métricas cuantitativas reales: 230+ micro-tareas, 929 commits analizados, 130 defectos pre-merge interceptados, 46 validaciones fallidas en corte, y ~650 h de retrabajo ahorradas. | `05_metricas_y_datos_cuantitativos.md` |
| **Sección IV y VI.B-D: Superación del *Judgment Gap* y límites de la IA** | Teoría sobre el juicio humano, la asunción de riesgos y la supervisión ética. | Casos empíricos donde la automatización falló o requirió corrección manual: rigidez sintáctica en `BAC-18B`, colisiones de puerto en `INF-06A`, orden de migraciones en frío (`FIX-81`), y falsa granularidad en UI. | `07_casos_fallo_automatizacion_y_judgment_gap.md` |
| **Sección VII.D: Limitaciones y agenda de investigación** | Declaración de falta de validación empírica. | Reescritura de las limitaciones: presentar a *El Centinela* como caso de estudio preliminar con arnés empírico reproducible, abriendo la agenda a estudios longitudinales. | `00_analisis_articulo_y_mapa_evidencia.md` |

---

## 3. Estructura de los Archivos Entregados

Para facilitar su lectura, consulta e integración directa en el manuscrito LaTeX/Word, la documentación se ha desglosado en los siguientes 7 documentos especializados:

1. **[`01_estructura_repo_readme_submodulos.md`](./01_estructura_repo_readme_submodulos.md):** Organización del monorrepo orquestador, punteros fijados (`detached HEAD`), `.gitmodules`, arquitectura de red en Docker y desacople *Local-First*.
2. **[`02_pipeline_ci_generacion_tareas_y_fixes.md`](./02_pipeline_ci_generacion_tareas_y_fixes.md):** GitHub Actions, `repository_dispatch`, serialización por concurrencia, taxonomía de tareas (`BAC`, `FRN`, `INF`, `BRG`) y algoritmo de creación de tickets `FIX`.
3. **[`03_reportes_suite_tests_integracion_front_back.md`](./03_reportes_suite_tests_integracion_front_back.md):** Arnés independiente de pruebas `/test`, ejecución Go $\leftrightarrow$ Vitest, simulación de sesión y navegador, cobertura de hitos `INT-02` y `LOGIN-04`.
4. **[`04_reportes_auditoria_seguridad.md`](./04_reportes_auditoria_seguridad.md):** Diagnóstico de seguridad, falsedad de la matriz de permisos en UI, análisis de las 9 brechas CWE y blindaje anti-IDOR mediante tokens opacos.
5. **[`05_metricas_y_datos_cuantitativos.md`](./05_metricas_y_datos_cuantitativos.md):** Cifras exactas de tareas, commits, fallas de validación, reasignaciones de autoría real vs tablero y cálculo de horas hombre ahorradas.
6. **[`06_episodio_divergencia_real_vs_pr1234.md`](./06_episodio_divergencia_real_vs_pr1234.md):** Crónica técnica detallada de las divergencias reales que sustituyen al ejemplo ficticio del PR #1234 (discrepancia en códigos de error y enrutamiento ofuscado).
7. **[`07_casos_fallo_automatizacion_y_judgment_gap.md`](./07_casos_fallo_automatizacion_y_judgment_gap.md):** Análisis crítico de los falsos positivos y limitaciones de la automatización que justifican de forma contundente la necesidad del juicio humano del PM.
