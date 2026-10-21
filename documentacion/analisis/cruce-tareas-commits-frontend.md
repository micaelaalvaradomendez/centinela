# Cruce de tareas, commits y autores reales — Submódulo `frontend`

**Proyecto:** El Centinela · **Repositorio analizado:** `frontend/` (`luzpacello/centinela-front`)  
**Fecha del análisis:** 20/10/2026  
**Universo:** 536 commits (desde `c04eddb` hasta `91eabf4`, 2026-08-17 a 2026-10-18)  
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
| **Mensaje/feature** | El asunto o los archivos nombran la funcionalidad exacta de la tarea. | Alta confianza. |
| **Inferida** | El asunto y los archivos corresponden al módulo de la tarea, sin cita directa. | Media; revisar si se usa para evaluar desempeño. |

> [!IMPORTANT]
> Los **158 merges** no aportan código propio: se listan y se asocian a las tareas de la rama integrada, pero **no cuentan** para la autoría efectiva. La autoría se calcula sobre los 378 commits con cambios reales.

> [!NOTE]
> Nombres normalizados: `Cristian`/`CristiannMaidana` → **Cristian**; `ibelinday`/`Belinda` → **Belinda**; `luzpacello`/`Lucesita` → **Luz**; `NicoLencina`/`Lencina Nicolas Angelo` → **Nico**. En asignaciones se toma el nombre principal antes del paréntesis (p. ej. “Luz (implementación inicial de Cristian)” → Luz).

---

## 2. Panorama general de autoría

| Autor | Commits totales | % | Commits con código | % | Primer commit | Último commit |
|---|---:|---:|---:|---:|---|---|
| **Cristian** | 340 | 63.4% | 271 | 71.7% | 2026-08-18 | 2026-10-18 |
| **Belinda** | 84 | 15.7% | 36 | 9.5% | 2026-08-18 | 2026-10-16 |
| **Nico** | 77 | 14.4% | 50 | 13.2% | 2026-09-13 | 2026-10-18 |
| **Luz** | 35 | 6.5% | 21 | 5.6% | 2026-08-17 | 2026-10-18 |
| **Total** | **536** | **100%** | **378** | **100%** | 2026-08-17 | 2026-10-18 |

- **Merges:** 158 · **Commits con cambios reales:** 378
- **Commits vinculados a tareas:** 501/536 (93.5%). De los commits reales, 369/378 (97.6%).
- **Vinculación explícita:** 94 commits; el resto se resolvió por mensaje/feature o inferencia de módulo.
- **Tareas frontend analizadas:** 103 (con y sin commits).

| Persona | Tareas frontend asignadas | Tareas en las que commiteó | Propias | Ajenas |
|---|---:|---:|---:|---:|
| **Cristian** | 48 | 42 | 20 | 22 |
| **Nico** | 3 | 28 | 3 | 25 |
| **Belinda** | 45 | 17 | 12 | 5 |
| **Luz** | 40 | 19 | 9 | 10 |

> _Una tarea con dos asignados cuenta para ambos. “Ajenas” = tareas ejecutadas sin estar asignada. Incluye tareas cross-area e infraestructura referenciadas (ver 5.5)._

---

## 3. Listado y análisis de los 536 commits

Orden cronológico. **Base:** Explícita (hash citado en docs/informe), Mensaje (el asunto nombra la feature) o Inferida (módulo/archivos). Los merges figuran como integración y no suman autoría.

| # | Fecha | Hash | Autor | Asunto | Tarea(s) | Base |
|---:|---|---|---|---|---|---|
| 1 | 2026-08-17 | `c04eddb` | Luz | Initial commit | — | — |
| 2 | 2026-08-18 | `08a7fcf` | Luz | 🌱 vite + react instalado con pnpm | `FRN-03` | Inferida |
| 3 | 2026-08-18 | `79f03d3` | Belinda | Merge pull request #1 from luzpacello/luz _(merge)_ | `FRN-03` | Integración |
| 4 | 2026-08-18 | `98f21df` | Cristian | docs(README): pasos para instalar | — | — |
| 5 | 2026-08-19 | `ed1d3d8` | Luz | 🌱 se agrega las carpetas faltantes | `FRN-03` | Inferida |
| 6 | 2026-08-19 | `923649f` | Cristian | Merge pull request #2 from luzpacello/cristian _(merge)_ | — | — |
| 7 | 2026-08-25 | `d009ac1` | Cristian | Merge pull request #3 from luzpacello/luz _(merge)_ | `FRN-03` | Integración |
| 8 | 2026-08-26 | `2bad8ab` | Belinda | ❄️ Maqueta la estructura base y estilos del layout principal | `FRN-03` | Inferida |
| 9 | 2026-08-26 | `9a2bf05` | Luz | Merge pull request #4 from luzpacello/rama-Beli _(merge)_ | `FRN-03` | Integración |
| 10 | 2026-09-01 | `851a18e` | Cristian | build(tsconfig.app.json): cambiar archivo para que typescript sepa que '@' representa 'src' en… | `FRN-01` | Inferida |
| 11 | 2026-09-01 | `452c492` | Cristian | build(vite.config.ts): cambiar archivo para que sepa que cuando hay '@' representa 'src' | `FRN-01` | Inferida |
| 12 | 2026-09-02 | `d0cb6e6` | Cristian | build(tsconfig.json): agregar ruta de path para shadcn sepa donde escribir | `FRN-01` | Inferida |
| 13 | 2026-09-02 | `f8cdc29` | Cristian | Merge pull request #5 from luzpacello/cristian _(merge)_ | `FRN-01` | Integración |
| 14 | 2026-09-06 | `79dacbd` | Cristian | build(package.json): instalar dependencias de shadcn/ui components | `FRN-01` | Inferida |
| 15 | 2026-09-06 | `366004a` | Cristian | build(pnpm-lock.yaml): instalar dependencias de shadcn/ui component | `FRN-01` | Inferida |
| 16 | 2026-09-06 | `4bb8c95` | Cristian | feat(index.css): agregar etilo necesarios luego de instalacion de dependencia de shadcn/ui com… | `FRN-01` | Inferida |
| 17 | 2026-09-06 | `2117570` | Cristian | feat(components.json): se creo archivo luego de instalar dependencia shadcn/ui componentes | `FRN-01` | Inferida |
| 18 | 2026-09-06 | `26aa8eb` | Cristian | feat(utils): agrega utilidad cn para componentes de shadcn | `FRN-01` | Inferida |
| 19 | 2026-09-06 | `2b81804` | Cristian | feat(button): creacion de button | `FRN-01` | Inferida |
| 20 | 2026-09-06 | `427714f` | Cristian | feat(input): creacion de input | `FRN-01` | Inferida |
| 21 | 2026-09-07 | `4d2c8fe` | Belinda | Merge pull request #6 from luzpacello/cristian _(merge)_ | `FRN-01` | Integración |
| 22 | 2026-09-09 | `0144f42` | Belinda | ❄️ Configura tipografía global de Penpot en Tailwind | `FRN-01` | Inferida |
| 23 | 2026-09-09 | `15feecd` | Luz | 🌱se agrega el componente stepper y Layout para auth y 2fa | `FRN-01` | Inferida |
| 24 | 2026-09-09 | `06af276` | Belinda | ❄️ Configura tipografía global de Penpot en Tailwind | `FRN-01` | Inferida |
| 25 | 2026-09-09 | `1a512b9` | Belinda | Merge pull request #7 from luzpacello/rama-Beli _(merge)_ | `FRN-01` | Integración |
| 26 | 2026-09-09 | `4ba0083` | Luz | 🌱 Restauración de un button | `FRN-01` | Inferida |
| 27 | 2026-09-09 | `30119b3` | Cristian | Merge pull request #8 from luzpacello/luz _(merge)_ | `FRN-01` | Integración |
| 28 | 2026-09-09 | `f2c3955` | Cristian | feat(checkbox): crear componente reutilizable checkbox con estilo | `FRN-01` | Inferida |
| 29 | 2026-09-09 | `bde4238` | Cristian | feat(Login): crear pagina login | `FRN-01`, `FRN-02` | Inferida |
| 30 | 2026-09-09 | `393012f` | Cristian | refactor(Login): cambiar pregunta | `FRN-01`, `FRN-02` | Inferida |
| 31 | 2026-09-09 | `55e3db6` | Cristian | feat(SignUp): crear pagina de registro | `FRN-01`, `FRN-02` | Inferida |
| 32 | 2026-09-10 | `0bfaa81` | Cristian | feat(separator): crear componente separator | `FRN-01` | Inferida |
| 33 | 2026-09-10 | `27c26cc` | Cristian | feat(label): crear componente label necesario para componente field | `FRN-01` | Inferida |
| 34 | 2026-09-10 | `db00e97` | Cristian | feat(field): crear componente field para estructurar los formularios | `FRN-01` | Inferida |
| 35 | 2026-09-10 | `d8ff411` | Belinda | Merge pull request #9 from luzpacello/cristian _(merge)_ | `FRN-01`, `FRN-02` | Integración |
| 36 | 2026-09-10 | `ea92e52` | Luz | 🌱 feat(2fa): agrego TwoFactorPage y TwoFactoForm con validación interna | `LOGIN-02` | Inferida |
| 37 | 2026-09-10 | `09547ca` | Cristian | Merge pull request #10 from luzpacello/luz _(merge)_ | `LOGIN-02` | Integración |
| 38 | 2026-09-10 | `34f631b` | Cristian | build(package): instalacion del paquete react-router | `FRN-03` | Inferida |
| 39 | 2026-09-10 | `672fed1` | Cristian | build(pnpm-lock): instalacion del paquete react-router | `FRN-03` | Inferida |
| 40 | 2026-09-11 | `4f2763e` | Belinda | ❄️ pantalla de recuperacion de contrasena | `FRN-12` | Inferida |
| 41 | 2026-09-11 | `da65fb6` | Belinda | Merge branch 'main' into rama-Beli _(merge)_ | — | — |
| 42 | 2026-09-11 | `30cecfb` | Luz | Merge pull request #11 from luzpacello/rama-Beli _(merge)_ | `FRN-12` | Integración |
| 43 | 2026-09-11 | `a9a8929` | Cristian | refactor(card): eliminar importacion de icons no utilizada | `FRN-01` | Inferida |
| 44 | 2026-09-11 | `9601e47` | Cristian | refactor(inputOTPSeparate): simplificar importacion | `FRN-01` | Inferida |
| 45 | 2026-09-11 | `3074fe6` | Cristian | refactor(label): simplificar estilo | `FRN-01` | Inferida |
| 46 | 2026-09-11 | `b4785dc` | Cristian | feat: add authentication foundation | `FRN-04` | Inferida |
| 47 | 2026-09-11 | `64460c8` | Belinda | Merge pull request #13 from luzpacello/feature/auth-foundation _(merge)_ | `FRN-04` | Integración |
| 48 | 2026-09-11 | `5cf2254` | Belinda | Merge pull request #12 from luzpacello/cristian _(merge)_ | `FRN-01`, `FRN-03` | Integración |
| 49 | 2026-09-11 | `368a4ce` | Cristian | feat: add login flow | `FRN-02`, `FRN-04` | Inferida |
| 50 | 2026-09-11 | `cd7794b` | Belinda | Merge pull request #14 from luzpacello/feature/login-flow _(merge)_ | `FRN-02`, `FRN-04` | Integración |
| 51 | 2026-09-11 | `6d63cf7` | Cristian | feat: add organization registration flow | `FRN-01` | Inferida |
| 52 | 2026-09-11 | `f75a536` | Belinda | Merge pull request #15 from luzpacello/feature/organization-registration _(merge)_ | `FRN-01` | Integración |
| 53 | 2026-09-11 | `35d22d9` | Cristian | feat: integrate application router | `FRN-03` | Inferida |
| 54 | 2026-09-11 | `1c4dab9` | Belinda | Merge pull request #16 from luzpacello/feature/data-router-integration _(merge)_ | `FRN-03` | Integración |
| 55 | 2026-09-13 | `33c3ed4` | Nico | necesario para comunicacion Nico | `FRN-03` | Inferida |
| 56 | 2026-09-13 | `17c5a8a` | Cristian | Merge pull request #17 from luzpacello/feat/componente-comunicacion-nico _(merge)_ | `FRN-03` | Integración |
| 57 | 2026-09-15 | `a1e85eb` | Cristian | feat(MainLayoutAuth): cambio de texto para testear conexion con infra | `FRN-03` | Inferida |
| 58 | 2026-09-15 | `9001530` | Cristian | Merge pull request #18 from luzpacello/cristian _(merge)_ | `FRN-03` | Integración |
| 59 | 2026-09-15 | `3bfbccc` | Nico | ci: deploy automatico del front a CT 103 (PRUEBAS) en cada push a main | — | — |
| 60 | 2026-09-15 | `86b9a02` | Luz | Merge pull request #19 from luzpacello/ci/deploy-automatico-front _(merge)_ | — | — |
| 61 | 2026-09-15 | `efb495a` | Nico | fix(auth): alinear la capa de autenticacion al contrato real del backend | `FRN-04` | Inferida |
| 62 | 2026-09-15 | `6dcdeb5` | Belinda | Merge pull request #20 from luzpacello/fix/auth-contrato-real-backend _(merge)_ | `FRN-04` | Integración |
| 63 | 2026-09-15 | `b803b20` | Cristian | feat(MainLayoutAuth): cambio de texto (solo prueba) | `FRN-03` | Inferida |
| 64 | 2026-09-15 | `ee4ba82` | Cristian | feat(Dashboard): creacion de pagina, solo ubicacion de elementos primera version | `FRN-03` | Inferida |
| 65 | 2026-09-15 | `ac5abdb` | Cristian | feat(Instances): creacion de pagina, solo ubicacion de elementos primera version | `FRN-03` | Inferida |
| 66 | 2026-09-15 | `eaea038` | Cristian | feat(NotFound): creacion de pagina, solo ubicacion de elementos primera version | `FRN-03` | Inferida |
| 67 | 2026-09-15 | `873d183` | Cristian | Merge branch 'main' into cristian _(merge)_ | — | — |
| 68 | 2026-09-15 | `8965611` | Belinda | Merge pull request #21 from luzpacello/cristian _(merge)_ | `FRN-03` | Integración |
| 69 | 2026-09-15 | `d8a36a5` | Nico | feat(layout_auth): mostrar fecha y hora del deploy en el panel | — | — |
| 70 | 2026-09-15 | `70e59a2` | Nico | ci(deploy): el deploy del front solo se dispara desde main | — | — |
| 71 | 2026-09-15 | `7a05a16` | Belinda | Merge pull request #22 from luzpacello/feat/stamp-deploy-panel _(merge)_ | — | — |
| 72 | 2026-09-15 | `566b5ab` | Belinda | Merge pull request #23 from luzpacello/ci/solo-main-despliega _(merge)_ | — | — |
| 73 | 2026-09-15 | `2f779e7` | Nico | infra(deploy-check): bloque temporal que muestra front y back juntos | — | — |
| 74 | 2026-09-15 | `26a927e` | Belinda | Merge pull request #24 from luzpacello/infra/deploy-check-front-back _(merge)_ | — | — |
| 75 | 2026-09-15 | `c259f20` | Nico | infra(deploy-check): el sello del front incluye el commit | — | — |
| 76 | 2026-09-15 | `dc605af` | Belinda | Merge pull request #25 from luzpacello/infra/front-stamp-con-commit _(merge)_ | — | — |
| 77 | 2026-09-16 | `45db606` | Cristian | feat(Sidebar): agregar nuevos links a paginas | `FRN-03` | Inferida |
| 78 | 2026-09-16 | `c674703` | Cristian | feat(applicationRoutes): crear nuevos links de pagina | `FRN-03` | Inferida |
| 79 | 2026-09-16 | `0dc7d8b` | Cristian | feat(Users): crear pagina de gestion de usuarios | `FRN-05` | Inferida |
| 80 | 2026-09-16 | `7fff0a0` | Cristian | feat(CrearUsuarios): crear pagina de creacion de usuarios | `FRN-06` | Inferida |
| 81 | 2026-09-16 | `766643e` | Cristian | Merge pull request #26 from luzpacello/cristian _(merge)_ | `FRN-03`, `FRN-05`, `FRN-06` | Integración |
| 82 | 2026-09-16 | `d315882` | Nico | fix(routes): quitar import y ruta huerfanos a ComponentPlayground | `FRN-03` | Inferida |
| 83 | 2026-09-16 | `07df625` | Nico | Merge pull request #27 from luzpacello/fix/routes-import-huerfano-componentplayground _(merge)_ | `FRN-03` | Integración |
| 84 | 2026-09-16 | `22d498a` | Cristian | feat(LoginForm): agregar link para navegacion a registro | `FRN-02`, `FRN-03` | Inferida |
| 85 | 2026-09-16 | `25d13a1` | Cristian | feat(authenticationActions): crear funcion para la validacion de datos de creacion de organiza… | `FRN-02`, `FRN-03` | Inferida |
| 86 | 2026-09-16 | `15f5a22` | Cristian | feat(applicationRoutes): crear link para pagina de registro de organizacion | `FRN-02`, `FRN-03` | Inferida |
| 87 | 2026-09-16 | `0d186da` | Cristian | Merge pull request #28 from luzpacello/cristian _(merge)_ | `FRN-02`, `FRN-03` | Integración |
| 88 | 2026-09-16 | `b15c588` | Cristian | refactor(applicationRoutes): eliminar importacion de pagina no existente | `FRN-03` | Inferida |
| 89 | 2026-09-16 | `86efc92` | Cristian | merge: sincronizar cristian con main _(merge)_ | — | — |
| 90 | 2026-09-16 | `ab434da` | Cristian | Merge pull request #29 from luzpacello/cristian _(merge)_ | `FRN-03` | Integración |
| 91 | 2026-09-16 | `feaaf27` | Belinda | acepta solo numeros | `LOGIN-02` | Inferida |
| 92 | 2026-09-16 | `0dabc28` | Belinda | Merge pull request #30 from luzpacello/rama-Beli _(merge)_ | `LOGIN-02` | Integración |
| 93 | 2026-09-16 | `c784f04` | Belinda | qr de autenticacion | `FRN-09` | Inferida |
| 94 | 2026-09-16 | `78659d0` | Belinda | Merge pull request #31 from luzpacello/rama-Beli _(merge)_ | `FRN-09` | Integración |
| 95 | 2026-09-17 | `51e7106` | Cristian | feat(nativeSelected): crear componente y agregarle colores segun diseño | `FRN-01` | Inferida |
| 96 | 2026-09-17 | `5d36e22` | Cristian | feat(tabs): crear componente y agregarle colores segun diseño | `FRN-01` | Inferida |
| 97 | 2026-09-17 | `6b1ccc4` | Cristian | feat(bardge): crear componente | `FRN-01` | Inferida |
| 98 | 2026-09-17 | `bc51bdb` | Cristian | feat(toast): crear componente y agregarle color segun diseño | `FRN-01` | Inferida |
| 99 | 2026-09-17 | `7775b62` | Cristian | Merge pull request #32 from luzpacello/cristian _(merge)_ | `FRN-01` | Integración |
| 100 | 2026-09-18 | `1f62bfe` | Cristian | feat(tokenStorage): crear archivo sessionStorage para los tokens de acceso, temporal, de sesio… | `FRN-04`, `SEC-02` | Inferida |
| 101 | 2026-09-18 | `52ea961` | Cristian | refactor(pendingLoginSession): ahora almacena el sessionStorage para poder seguir con la sesio… | `FRN-04`, `SEC-02` | Inferida |
| 102 | 2026-09-18 | `432d747` | Cristian | refactor(api): cambio localStorage por sessionStorage, conservando mismo nombre de variables | `FRN-04`, `SEC-02` | Inferida |
| 103 | 2026-09-19 | `9e8b45c` | Cristian | Merge pull request #33 from luzpacello/cristian _(merge)_ | `FRN-04`, `SEC-02` | Integración |
| 104 | 2026-09-19 | `cbf7481` | Belinda | ❄️pantalla auditoria | `FRN-14` | Inferida |
| 105 | 2026-09-19 | `222fac5` | Belinda | Merge pull request #34 from luzpacello/rama-Beli _(merge)_ | `FRN-14` | Integración |
| 106 | 2026-09-19 | `daa8bd5` | Cristian | feat(ApiResponseNotifier): crea el mensaje del toast, y asigna paginas para cuando no se aplic… | `FRN-08` | Inferida |
| 107 | 2026-09-19 | `2711505` | Cristian | feat(apiClient): crea las conexiones fetch http | `FRN-08` | Inferida |
| 108 | 2026-09-19 | `84afd58` | Cristian | refactor(api): usa llamados de apiClient | `FRN-08` | Inferida |
| 109 | 2026-09-19 | `acc5ecf` | Cristian | feat(App): envuelve el router, para el manejo de notificaciones | `FRN-08` | Inferida |
| 110 | 2026-09-19 | `12b0779` | Cristian | Merge pull request #35 from luzpacello/cristian _(merge)_ | `FRN-08` | Integración |
| 111 | 2026-09-20 | `6b08566` | Belinda | se agregan pantallas auditoria, crear usuario, user Detail, users | `FIX-08` | Explícita |
| 112 | 2026-09-20 | `c7a8800` | Belinda | Merge branch 'main' into rama-Beli _(merge)_ | — | — |
| 113 | 2026-09-20 | `68bd23f` | Belinda | Merge pull request #36 from luzpacello/rama-Beli _(merge)_ | `FIX-08` | Integración |
| 114 | 2026-09-21 | `231779f` | Nico | fix(86e3bhj61): restaurar guards de sesion en las rutas privadas | `FIX-08` | Explícita |
| 115 | 2026-09-21 | `64cad86` | Nico | fix(86e3bh90p): panel de usuarios contra la API real y guard por rol | `FRN-05`, `FIX-12` | Inferida |
| 116 | 2026-09-21 | `8a639a9` | Nico | Merge pull request #37 from luzpacello/Nico _(merge)_ | `FIX-08`, `FIX-12`, `FRN-05` | Integración |
| 117 | 2026-09-21 | `63dc6de` | Nico | 86e3bh9d7 - mensajes de error alineados con los códigos del backend | `FIX-07` | Inferida |
| 118 | 2026-09-21 | `a815bd6` | Nico | Merge pull request #38 from luzpacello/Nico _(merge)_ | `FIX-07` | Integración |
| 119 | 2026-09-21 | `e4a9798` | Nico | 86e3bh95w - alta y edición de usuarios conectadas a la API | `FRN-06`, `FIX-06` | Inferida |
| 120 | 2026-09-21 | `dbb4758` | Nico | Merge pull request #39 from luzpacello/Nico _(merge)_ | `FIX-06`, `FRN-06` | Integración |
| 121 | 2026-09-21 | `6145740` | Luz | 🌱 refactor(2fa): Se cambia de LoginContinue a TwoFactorPage | `LOGIN-02` | Inferida |
| 122 | 2026-09-21 | `3e330fc` | Belinda | Merge pull request #40 from luzpacello/luz _(merge)_ | `LOGIN-02` | Integración |
| 123 | 2026-09-21 | `9ccef0e` | Cristian | refactor(Users): eliminar panel lateral de detalle de usuarios | `FRN-05` | Inferida |
| 124 | 2026-09-21 | `3112eea` | Belinda | se cambian extensiones | `FRN-03` | Inferida |
| 125 | 2026-09-21 | `e36684d` | Belinda | Merge pull request #41 from luzpacello/rama-Beli _(merge)_ | `FRN-03` | Integración |
| 126 | 2026-09-21 | `c8557b6` | Nico | 86e3bpfw3 - Tipos del contrato de eventos compartidos con el backend | `FRN-17A` | Inferida |
| 127 | 2026-09-21 | `0b50cd6` | Belinda | corrige el vite | `FRN-03` | Inferida |
| 128 | 2026-09-21 | `da396a1` | Belinda | corrige import sin uso | `FRN-03` | Inferida |
| 129 | 2026-09-21 | `8bcb893` | Belinda | Merge pull request #42 from luzpacello/rama-Beli _(merge)_ | `FRN-03` | Integración |
| 130 | 2026-09-21 | `03209be` | Cristian | refactor(tabs): cambiar estilo de componente | `FRN-05` | Inferida |
| 131 | 2026-09-21 | `e00a1ab` | Cristian | feat(informationOfUser): crear componente | `FRN-06B` | Inferida |
| 132 | 2026-09-21 | `6f9cdfa` | Cristian | refactor(informationOfUser): cambiar tipo de opcion de rol | `FRN-06B` | Inferida |
| 133 | 2026-09-21 | `dace934` | Cristian | feat(index): crear 2 estilos de texto | `FRN-05` | Inferida |
| 134 | 2026-09-21 | `fdab248` | Cristian | feat(table): crear componente table | `FRN-05` | Inferida |
| 135 | 2026-09-21 | `1832e85` | Nico | Merge pull request #43 from luzpacello/Nico _(merge)_ | `FRN-17A` | Integración |
| 136 | 2026-09-21 | `8b3a693` | Cristian | feat(rolesAndPermissions): crear componente | `FRN-07` | Inferida |
| 137 | 2026-09-21 | `8687564` | Cristian | refactor(badge): cambiar nombre de componente | `FRN-06B` | Inferida |
| 138 | 2026-09-21 | `dc84dfd` | Cristian | feat(detailsUserPage): crear marco de pagina de detalle de usuarios | `FRN-06B` | Inferida |
| 139 | 2026-09-21 | `7a520c1` | Nico | fix(users): arreglo la carga del usuario real en UserDetail | `FRN-05` | Inferida |
| 140 | 2026-09-21 | `5ad9b88` | Cristian | feat(Users): agregar funcionalidad a boton | `FRN-06` | Inferida |
| 141 | 2026-09-21 | `e47e041` | Nico | fix(users): arreglo el guard de rol ADMIN en las rutas de usuarios | `FIX-12` | Inferida |
| 142 | 2026-09-21 | `ca52e47` | Nico | fix(2fa): arreglo el mensaje de error de verificación TOTP | `LOGIN-02` | Inferida |
| 143 | 2026-09-21 | `e134383` | Nico | fix(auditoria): arreglo la vista de auditoría para consumir la API real | `FRN-14` | Inferida |
| 144 | 2026-09-21 | `27b4467` | Nico | fix(rutas): arreglo los imports con extensión .jsx | `FRN-03` | Inferida |
| 145 | 2026-09-21 | `d9e0339` | Nico | chore(limpieza): arreglo el código muerto del front | `FRN-03` | Inferida |
| 146 | 2026-09-21 | `efeb4a5` | Nico | fix(errores): arreglo el copy de los códigos de error que faltaban | `FIX-07` | Inferida |
| 147 | 2026-09-21 | `cc8a8a3` | Cristian | feat(switch): crear componente | `FRN-06B` | Inferida |
| 148 | 2026-09-21 | `e5da01b` | Cristian | refactor(rolesAndPermissions): cambiar nombre de importacion | `FRN-06B` | Inferida |
| 149 | 2026-09-21 | `9c1f768` | Cristian | feat(UserDetail): eliminar pagina | `FRN-05` | Inferida |
| 150 | 2026-09-21 | `d9cd6c1` | Nico | ci(front): agrego workflow de build en pull_request | — | — |
| 151 | 2026-09-21 | `a4d3f65` | Cristian | feat(applicationRoutes): crear ruta para nueva pagina | `FRN-03` | Inferida |
| 152 | 2026-09-22 | `924d65f` | Nico | fix(auditoria): arreglo las métricas mock de la vista de auditoría | `FRN-14` | Inferida |
| 153 | 2026-09-22 | `ea5f963` | Nico | Merge pull request #44 from luzpacello/Nico _(merge)_ | `FIX-07`, `FIX-12`, `FRN-03`, `FRN-05`, `FRN-14`, `LOGIN-02` | Integración |
| 154 | 2026-09-22 | `7bc6dc5` | Nico | fix(menu): arreglo los enlaces muertos del menu lateral | `FRN-03` | Inferida |
| 155 | 2026-09-22 | `a7e7eb3` | Cristian | feat(user): crear persistencia de los datos de informacion del usuario | `FIX-25` | Inferida |
| 156 | 2026-09-22 | `ff0bf64` | Cristian | feat(userDetailsServices): crear peticion de busqueda de datos del usuario | `FIX-25` | Inferida |
| 157 | 2026-09-22 | `f89ab3c` | Cristian | feat(useUserDetails): crear hook para el manejo de la peticion de detalle de usuario | `FIX-25` | Inferida |
| 158 | 2026-09-22 | `24a97e7` | Cristian | feat(useEditableUser): crear hook para el manejo de datos editables del usuario | `FIX-25` | Inferida |
| 159 | 2026-09-22 | `4c3fcf3` | Cristian | feat(Users): cambiar boton de navegacion | `FRN-05` | Inferida |
| 160 | 2026-09-22 | `66eb5be` | Nico | Merge pull request #45 from luzpacello/Nico _(merge)_ | `FRN-03` | Integración |
| 161 | 2026-09-22 | `6f42344` | Cristian | feat(detailUserPage): cargar los datos reales del usuario | `FIX-25` | Inferida |
| 162 | 2026-09-22 | `98c27c7` | Cristian | feat(informationOfUser): cambiar los valores estaticos de usuario por los reales | `FIX-25` | Inferida |
| 163 | 2026-09-22 | `06c379d` | Cristian | feat(rolesAndPermissions): cambiar valor estatico de rol del usuario por el real que tiene | `FIX-25` | Inferida |
| 164 | 2026-09-22 | `893bf09` | Nico | fix(auth): arreglo el cambio obligatorio de contrasena para que el usuario pueda entrar | `FRN-10` | Inferida |
| 165 | 2026-09-22 | `266fc1d` | Cristian | feat(userUserDetails): cambiar modo de alerta de error ahora con toast | `FRN-08` | Inferida |
| 166 | 2026-09-22 | `d6d2005` | Nico | Merge pull request #46 from luzpacello/Nico _(merge)_ | `FRN-10` | Integración |
| 167 | 2026-09-22 | `541bf1c` | Cristian | feat(ApiResponseNotifier): crear alertas para diferentes errores | `FRN-08` | Inferida |
| 168 | 2026-09-22 | `c4ac90c` | Cristian | feat(userInstanceService): crear conexion a endpoint para obtener la lista de isntancias | `FRN-07`, `FIX-14` | Inferida |
| 169 | 2026-09-22 | `cdca5f4` | Cristian | feat(useUserInstanceAccess): crear el hook para el manejo de la instancias al usuario | `FRN-07`, `FIX-14` | Inferida |
| 170 | 2026-09-22 | `d6feb7b` | Cristian | feat(deatilUserPage): importar hook | `FRN-07`, `FIX-14` | Inferida |
| 171 | 2026-09-22 | `1e8e8f3` | Cristian | feat(rolesAndPermissions): cambiar datos estaticos de tabla por dinamicos de endpoints | `FRN-07`, `FIX-14` | Inferida |
| 172 | 2026-09-22 | `ab5e0f9` | Belinda | ❄️ se unifica el header, y se rehubica cerrar sesion | `FRN-03` | Inferida |
| 173 | 2026-09-22 | `d088df2` | Belinda | Merge pull request #47 from luzpacello/rama-Beli _(merge)_ | `FRN-03` | Integración |
| 174 | 2026-09-22 | `635c8e7` | Belinda | ❄️  minimiza el alto del header | `FRN-03` | Inferida |
| 175 | 2026-09-22 | `d539253` | Belinda | Merge pull request #48 from luzpacello/main _(merge)_ | — | — |
| 176 | 2026-09-22 | `94f2d7b` | Cristian | merge: sincronizar cristian con main _(merge)_ | — | — |
| 177 | 2026-09-22 | `998da0c` | Cristian | feat(useSafeNavigate): eliminar hook de navegacion | `FRN-03` | Inferida |
| 178 | 2026-09-22 | `65e1dc1` | Cristian | feat(nativeSelected): agregar componente | `FRN-03` | Inferida |
| 179 | 2026-09-22 | `995b2d6` | Cristian | feat(useSafeNavigate): agregar archivo eliminado | `FRN-03` | Inferida |
| 180 | 2026-09-22 | `80fa862` | Cristian | Merge pull request #49 from luzpacello/cristian _(merge)_ | `FIX-14`, `FIX-25`, `FRN-03`, `FRN-05`, `FRN-06`, `FRN-06B`, `FRN-07`, `FRN-08` | Integración |
| 181 | 2026-09-22 | `bf82107` | Cristian | feat(userInstanceService): cambiar urls de las funciones acorde a enpoints actualizados | `FRN-07`, `FIX-14` | Inferida |
| 182 | 2026-09-22 | `791cb81` | Cristian | feat(userUserInstanceAccess): actualizar funcionalidad con endpoints actualizados | `FRN-07`, `FIX-14` | Inferida |
| 183 | 2026-09-22 | `f115fee` | Cristian | feat(rolesAndPermissions): actualizar funcionalidad para el funcionamiento de la asignacion de… | `FRN-07`, `FIX-14` | Inferida |
| 184 | 2026-09-22 | `10c723d` | Cristian | feat(detailsUserPage): actualizar los datos de la instancia por endpoint actualizado | `FRN-07`, `FIX-14` | Inferida |
| 185 | 2026-09-22 | `24695fa` | Cristian | merge: resolver sincronizacion pendiente _(merge)_ | — | — |
| 186 | 2026-09-22 | `de6b25b` | Cristian | Merge pull request #50 from luzpacello/cristian _(merge)_ | `FIX-14`, `FRN-07` | Integración |
| 187 | 2026-09-22 | `f07da88` | Belinda | Merge branch 'main' of https://github.com/luzpacello/centinela-front into rama-Beli _(merge)_ | — | — |
| 188 | 2026-09-22 | `76e7ebe` | Belinda | ❄️ se cambia posicion de cerrar sesion y se elimina el header | `FRN-03` | Inferida |
| 189 | 2026-09-22 | `4e5fbab` | Belinda | ❄️Merge de actualizaciones remotas _(merge)_ | — | — |
| 190 | 2026-09-22 | `61b2854` | Belinda | Merge pull request #51 from luzpacello/rama-Beli _(merge)_ | `FRN-03` | Integración |
| 191 | 2026-09-22 | `b50057e` | Belinda | ❄️  se olvido comentar | `FRN-03` | Inferida |
| 192 | 2026-09-22 | `44e57e6` | Belinda | Merge pull request #52 from luzpacello/main _(merge)_ | — | — |
| 193 | 2026-09-22 | `b1d5acc` | Belinda | ❄️  menu lateral con theme del sistema | `FRN-03` | Inferida |
| 194 | 2026-09-22 | `5b91e80` | Belinda | Merge pull request #53 from luzpacello/rama-Beli _(merge)_ | `FRN-03` | Integración |
| 195 | 2026-09-22 | `08d44ee` | Cristian | feat(authentication): eliminar propiedad refreshToken | `SEC-02` | Explícita |
| 196 | 2026-09-22 | `ba9b5cf` | Cristian | feat(tokenStorage): eliminar constante de 'refreshKey' y de funcion 'storeAuthTokens' | `SEC-02` | Explícita |
| 197 | 2026-09-22 | `4b7967f` | Cristian | feat(validateAuthenticationResponses): eliminar validacion de 'refreshToken' en metodo 'isToke… | `SEC-02` | Explícita |
| 198 | 2026-09-22 | `cb208f5` | Cristian | feat(apiClient): eliminar importacion y relacionado a 'refreshToken', actualizar validacion y … | `SEC-02` | Explícita |
| 199 | 2026-09-22 | `509a159` | Cristian | feat(api): eliminar importacion de metodo relacionado a 'refreshToken' | `SEC-02` | Explícita |
| 200 | 2026-09-22 | `f0f7218` | Cristian | feat(authService): eliminar importacion y uso de 'refreshToken' | `SEC-02` | Explícita |
| 201 | 2026-09-22 | `390b4b8` | Cristian | feat(api.d.ts): eliminar uso de 'refreshToken' por sessionStorage | `SEC-02` | Explícita |
| 202 | 2026-09-22 | `781db6c` | Cristian | Merge pull request #54 from luzpacello/cristian _(merge)_ | `SEC-02` | Integración |
| 203 | 2026-09-22 | `c75acce` | Cristian | refactor(Users): cambiar estilo de boton principal | `FRN-05` | Inferida |
| 204 | 2026-09-22 | `360d479` | Cristian | Merge pull request #55 from luzpacello/cristian _(merge)_ | `FRN-05` | Integración |
| 205 | 2026-09-23 | `c8e5592` | Belinda | ❄️  se prueba logo, sideBar estilo | `FRN-03` | Inferida |
| 206 | 2026-09-23 | `a2da502` | Belinda | ❄️  se arregla error | `FRN-03` | Inferida |
| 207 | 2026-09-23 | `3192cf4` | Belinda | Merge pull request #56 from luzpacello/rama-Beli _(merge)_ | `FRN-03` | Integración |
| 208 | 2026-09-23 | `131ee7d` | Belinda | Revert "Merge pull request #56 from luzpacello/rama-Beli" | `FRN-03` | Inferida |
| 209 | 2026-09-23 | `a387e10` | Cristian | feat(Sidebar): eliminar link a pagina crear usuarios | `FRN-03`, `FRN-06` | Inferida |
| 210 | 2026-09-23 | `684b09c` | Cristian | feat(Users): agregar funcionalidad de navegacion a boton, navega a pagina CrearUsuarios | `FRN-03`, `FRN-06` | Inferida |
| 211 | 2026-09-23 | `c92009e` | Cristian | feat(CrearUsuarios): agregar funcionalidad de navegacion en pagina | `FRN-03`, `FRN-06` | Inferida |
| 212 | 2026-09-23 | `bd40dc9` | Cristian | Merge pull request #57 from luzpacello/cristian _(merge)_ | `FRN-03`, `FRN-06` | Integración |
| 213 | 2026-09-23 | `bac88cc` | Cristian | refactor(Users): cambiar estilos y funcionalidad | `FRN-05` | Inferida |
| 214 | 2026-09-23 | `f2d6d9f` | Cristian | feat(userDetailService): agregar funcion para actualizar los datos del usuario | `FIX-25` | Explícita |
| 215 | 2026-09-23 | `60f81e5` | Cristian | feat(informationOfUser): agregar error a mail para validar que no haya 2 iguales | `FIX-25` | Explícita |
| 216 | 2026-09-23 | `5b3e073` | Cristian | feat(useUserInstanceAccess): agregar posible error | `FIX-25` | Explícita |
| 217 | 2026-09-23 | `0757bc8` | Cristian | feat(useEditableUser): agregar constante en funcion | `FIX-25` | Explícita |
| 218 | 2026-09-23 | `8164670` | Cristian | feat(detailsUserPage): agregar funcionalidad para la activacion y actualizacion del guardado d… | `FIX-25` | Explícita |
| 219 | 2026-09-23 | `e8a2e27` | Cristian | Merge pull request #58 from luzpacello/cristian _(merge)_ | `FIX-25`, `FRN-05` | Integración |
| 220 | 2026-09-23 | `a36a8c9` | Cristian | feat(modals): crear componente de modals | `FRN-06` | Inferida |
| 221 | 2026-09-23 | `6198cf9` | Cristian | feat(userDesactivationService): conexion para eliminar usuario | `FRN-06`, `FIX-06` | Inferida |
| 222 | 2026-09-23 | `bd23180` | Cristian | feat(toast): agregar propiedad forceExpanded para si es mas de una notificacion expandir y que… | `FRN-06`, `FIX-06` | Inferida |
| 223 | 2026-09-23 | `fa6cbb7` | Cristian | feat(useCreateUser): crear funcion para crear usuario, agregar notificaciones de exito, alerta… | `FRN-06`, `FIX-06` | Inferida |
| 224 | 2026-09-23 | `1fb8bfd` | Cristian | feat(createUser): crear persistencia de datos de creacion de usuario y respuesta | `FRN-06`, `FIX-06` | Inferida |
| 225 | 2026-09-23 | `8db47cf` | Cristian | refactor(useCreateUser): arreglar importacion de type de user | `FRN-06`, `FIX-06` | Inferida |
| 226 | 2026-09-23 | `597ba06` | Cristian | feat(createUserService): crear conexion con enpoint y manejo de respuestas de error | `FRN-06`, `FIX-06` | Inferida |
| 227 | 2026-09-23 | `f290dad` | Cristian | feat(ConfirmUserAction): generador de modal para cada tipo de confirmacion y situacion del usu… | `FRN-06`, `FIX-06` | Inferida |
| 228 | 2026-09-24 | `b13931c` | Cristian | feat(deatilsUserPage): agregar funcionalidad para eliminar usuarios | `FRN-06`, `FIX-06` | Inferida |
| 229 | 2026-09-24 | `58d15c0` | Cristian | feat(Users): agregar funcionalidad para eliminar usuarios | `FRN-06`, `FIX-06` | Inferida |
| 230 | 2026-09-24 | `21610fb` | Cristian | feat(CrearUsuraios): agregar funcionalida para la creacion de usuarios | `FRN-06`, `FIX-06` | Inferida |
| 231 | 2026-09-24 | `b82b6a5` | Cristian | Merge pull request #59 from luzpacello/cristian _(merge)_ | `FIX-06`, `FRN-06` | Integración |
| 232 | 2026-09-24 | `aa339d8` | Cristian | feat(userInstanceService): actualizar en base a endpoint | `FRN-07`, `FIX-14` | Inferida |
| 233 | 2026-09-24 | `96c56c6` | Cristian | feat(useUserInstanceAccess): actualizar en base a endpoint | `FRN-07`, `FIX-14` | Inferida |
| 234 | 2026-09-24 | `b865771` | Cristian | feat(rolesAndPermissions): actualizar en base a endpoint | `FRN-07`, `FIX-14` | Inferida |
| 235 | 2026-09-24 | `7fbf969` | Cristian | Merge pull request #60 from luzpacello/cristian _(merge)_ | `FIX-14`, `FRN-07` | Integración |
| 236 | 2026-09-25 | `e02b821` | Cristian | feat(AuthContext): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Explícita |
| 237 | 2026-09-25 | `4556331` | Cristian | feat(App): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Explícita |
| 238 | 2026-09-25 | `607e455` | Cristian | feat(useAuth): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Explícita |
| 239 | 2026-09-25 | `deb59cb` | Cristian | feat(authService): implementar flujo de logout y limpieza de sesion en cliente | `FIX-28`, `FRN-13` | Explícita |
| 240 | 2026-09-25 | `b900d19` | Cristian | feat(apiClient): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Explícita |
| 241 | 2026-09-25 | `72e205d` | Cristian | feat(api): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Explícita |
| 242 | 2026-09-25 | `300a0bd` | Cristian | feat(api.d): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Explícita |
| 243 | 2026-09-25 | `cce5a55` | Cristian | feat(tokenStorage): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Explícita |
| 244 | 2026-09-25 | `4062fdc` | Cristian | feat(ApiResponseNotifier): implementar flujo de logout y limpieza de sesion en cliente | `FRN-13` | Explícita |
| 245 | 2026-09-25 | `b8469ea` | Cristian | Merge pull request #61 from luzpacello/cristian _(merge)_ | `FIX-28`, `FRN-13` | Integración |
| 246 | 2026-09-25 | `b87ea1b` | Cristian | fix(ChangePassword): opción de cerrar sesión en el cambio obligatorio de contraseña | `FIX-21` | Explícita |
| 247 | 2026-09-25 | `b464c72` | Cristian | Merge pull request #62 from luzpacello/cristian _(merge)_ | `FIX-21` | Integración |
| 248 | 2026-09-26 | `fe158ab` | Luz | 🌱 feat(users): conectar restablecimiento de contraseña y 2FA | `FIX-19`, `FRN-11` | Explícita |
| 249 | 2026-09-26 | `6ad0c14` | Luz | Merge pull request #63 from luzpacello/luz _(merge)_ | `FIX-19`, `FRN-11` | Integración |
| 250 | 2026-09-27 | `fe6d784` | Luz | 🌱 Sincronización de políticas de complejidad y distinción UX entre Cambio y Restablecimiento | `FIX-18`, `FIX-20`, `FRN-12` | Explícita |
| 251 | 2026-09-27 | `3da3e78` | Luz | Merge pull request #64 from luzpacello/luz _(merge)_ | `FIX-18`, `FIX-20`, `FRN-12` | Integración |
| 252 | 2026-09-28 | `81c6424` | Belinda | ❄️; Se cambia Sidebar | `FRN-03` | Inferida |
| 253 | 2026-09-28 | `8d7ecab` | Belinda | Merge pull request #65 from luzpacello/rama-Beli _(merge)_ | `FRN-03` | Integración |
| 254 | 2026-09-29 | `5c3d780` | Belinda | ❄️; fix ruta auditoria | `FIX-22` | Explícita |
| 255 | 2026-09-29 | `bd0c42f` | Belinda | ❄️ fix request | `FIX-22` | Explícita |
| 256 | 2026-09-29 | `749194e` | Belinda | Merge pull request #66 from luzpacello/rama-Beli _(merge)_ | `FIX-22` | Integración |
| 257 | 2026-09-30 | `56b88f5` | Cristian | refactor(Auditoria): eliminar columnas que no seran utilizadas | `FIX-36`, `FRN-14` | Explícita |
| 258 | 2026-09-30 | `e6ad7ad` | Cristian | refactor(Auditoria): cambiar textos | `FRN-14` | Explícita |
| 259 | 2026-09-30 | `87d8926` | Cristian | Merge pull request #67 from luzpacello/cristian _(merge)_ | `FIX-36`, `FRN-14` | Integración |
| 260 | 2026-09-30 | `3ea7fa4` | Cristian | fix(useCreatUser): arreglar notificacion por fallo de mail no enviado a la hora de crear usuar… | `FIX-24`, `FIX-27` | Explícita |
| 261 | 2026-09-30 | `f6090c3` | Cristian | Merge pull request #68 from luzpacello/cristian _(merge)_ | `FIX-24`, `FIX-27` | Integración |
| 262 | 2026-09-30 | `cfb88f7` | Cristian | fix(authService): el logout ahora envia Bearer y cookie HttpOnly para el correcto logout | `FIX-28`, `FRN-13` | Explícita |
| 263 | 2026-09-30 | `bdc3116` | Cristian | Merge pull request #69 from luzpacello/cristian _(merge)_ | `FIX-28`, `FRN-13` | Integración |
| 264 | 2026-09-30 | `54b397e` | Cristian | fix(usePermissions): helper canOperateInstance en usePermissions() | `FIX-30` | Explícita |
| 265 | 2026-09-30 | `d6ac849` | Cristian | fix(authentication): helper canOperateInstance en usePermissions() | `FIX-30` | Explícita |
| 266 | 2026-09-30 | `1a61fae` | Cristian | fix(validateAuthenticationResponses): helper canOperateInstance en usePermissions() | `FIX-30` | Explícita |
| 267 | 2026-09-30 | `d1c2742` | Cristian | fix(api): helper canOperateInstance en usePermissions() | `FIX-30` | Explícita |
| 268 | 2026-09-30 | `0b632f1` | Cristian | fix(api.d): helper canOperateInstance en usePermissions() | `FIX-30` | Explícita |
| 269 | 2026-09-30 | `3d1e84a` | Cristian | Merge pull request #70 from luzpacello/cristian _(merge)_ | `FIX-30` | Integración |
| 270 | 2026-10-01 | `1156eaa` | Luz | 🌱 Contexto y sistema reactivo de permisos en Frontend (UI/UX RBAC + Resources) | `FIX-29`, `FIX-41`, `FIX-42`, `SEC-03` | Explícita |
| 271 | 2026-10-01 | `92c7286` | Luz | Merge branch 'main' into luz _(merge)_ | — | — |
| 272 | 2026-10-01 | `353f0f1` | Luz | 🌱 Solución al problema de merge | `SEC-03` | Explícita |
| 273 | 2026-10-01 | `952b43c` | Luz | 🌱 fix(informationOfUser): modifico la validación de administrador con useAuth | `SEC-03` | Explícita |
| 274 | 2026-10-01 | `47e5c19` | Luz | 🌱 refactor(auth): tipar AuthContext con TypeScript | `FIX-42`, `SEC-03` | Explícita |
| 275 | 2026-10-01 | `4e0e7b7` | Luz | Merge pull request #71 from luzpacello/luz _(merge)_ | `FIX-29`, `FIX-41`, `FIX-42`, `SEC-03` | Integración |
| 276 | 2026-10-02 | `3269dbd` | Cristian | fix(detailsUserPage): Regresión: requiredRole='Admin' oculta 'Roles y permisos' al administrad… | `FIX-41` | Explícita |
| 277 | 2026-10-02 | `f5e3c61` | Cristian | fix(authentication): Regresión: requiredRole='Admin' oculta 'Roles y permisos' al administrador | `FIX-41` | Explícita |
| 278 | 2026-10-02 | `a351564` | Cristian | fix(AuthContext): Regresión: requiredRole='Admin' oculta 'Roles y permisos' al administrador | `FIX-41` | Explícita |
| 279 | 2026-10-02 | `88f31b6` | Cristian | fix(PermissionGate): Regresión: requiredRole='Admin' oculta 'Roles y permisos' al administrador | `FIX-41` | Explícita |
| 280 | 2026-10-02 | `b65e8b2` | Cristian | Merge pull request #72 from luzpacello/cristian _(merge)_ | `FIX-41` | Integración |
| 281 | 2026-10-02 | `369fc05` | Cristian | fix(detailUserPage): 'Solo lectura' (READ_ONLY) usado erróneamente como rol de usuario | `FIX-38` | Explícita |
| 282 | 2026-10-02 | `b2a5161` | Cristian | fix(CrearUsuarios): 'Solo lectura' (READ_ONLY) usado erróneamente como rol de usuario | `FIX-38` | Explícita |
| 283 | 2026-10-02 | `f65fa67` | Cristian | fix(informationOfUser): 'Solo lectura' (READ_ONLY) usado erróneamente como rol de usuario | `FIX-38` | Explícita |
| 284 | 2026-10-02 | `1da04a1` | Cristian | fix(rolesAndPermissions): 'Solo lectura' (READ_ONLY) usado erróneamente como rol de usuario | `FIX-38` | Explícita |
| 285 | 2026-10-02 | `7d9e33a` | Cristian | fix(usePermissions): 'Solo lectura' (READ_ONLY) usado erróneamente como rol de usuario | `FIX-38` | Explícita |
| 286 | 2026-10-02 | `fb60cb2` | Cristian | fix(AuthContext): 'Solo lectura' (READ_ONLY) usado erróneamente como rol de usuario | `FIX-38`, `FIX-42` | Explícita |
| 287 | 2026-10-02 | `8a111cb` | Cristian | fix(validateAuthenticationResponses): 'Solo lectura' (READ_ONLY) usado erróneamente como rol d… | `FIX-38` | Explícita |
| 288 | 2026-10-02 | `6377790` | Cristian | fix(userDetailsService): 'Solo lectura' (READ_ONLY) usado erróneamente como rol de usuario | `FIX-38` | Explícita |
| 289 | 2026-10-02 | `6fd2c7c` | Cristian | Merge pull request #73 from luzpacello/cristian _(merge)_ | `FIX-38`, `FIX-42` | Integración |
| 290 | 2026-10-02 | `252f086` | Cristian | feat(Instances): Integración en Frontend de niveles de acceso por instancia (READ_ONLY / FULL_… | `FRN-18` | Inferida |
| 291 | 2026-10-02 | `1a7d9b6` | Cristian | feat(useInstances): Integración en Frontend de niveles de acceso por instancia (READ_ONLY / FU… | `FRN-18` | Inferida |
| 292 | 2026-10-02 | `1a78041` | Cristian | feat(instanceService): Integración en Frontend de niveles de acceso por instancia (READ_ONLY /… | `FRN-18` | Inferida |
| 293 | 2026-10-02 | `c941f5d` | Cristian | feat(InstanceAction): Integración en Frontend de niveles de acceso por instancia (READ_ONLY / … | `FRN-18` | Inferida |
| 294 | 2026-10-02 | `0f1434a` | Cristian | feat(AuthContext): Integración en Frontend de niveles de acceso por instancia (READ_ONLY / FUL… | `FRN-18` | Inferida |
| 295 | 2026-10-02 | `52cd282` | Cristian | Merge pull request #74 from luzpacello/cristian _(merge)_ | `FRN-18` | Integración |
| 296 | 2026-10-02 | `31f468e` | Cristian | fix(instanceAccess): useAuth().canAccessInstance no revisa el rol | `FIX-42` | Explícita |
| 297 | 2026-10-02 | `d32b218` | Cristian | fix(AuthContext): useAuth().canAccessInstance no revisa el rol | `FIX-42` | Explícita |
| 298 | 2026-10-02 | `ee60382` | Cristian | fix(usePermissions): useAuth().canAccessInstance no revisa el rol | `FIX-42` | Explícita |
| 299 | 2026-10-02 | `5a274e9` | Cristian | Merge pull request #75 from luzpacello/cristian _(merge)_ | `FIX-42` | Integración |
| 300 | 2026-10-02 | `2057fc0` | Belinda | ❄️  se agregar implmentan lo test y fix en audotoria | `FIX-36` | Explícita |
| 301 | 2026-10-02 | `cdbc8e9` | Belinda | Merge branch 'main' into rama-Beli _(merge)_ | — | — |
| 302 | 2026-10-02 | `1cdcaf7` | Belinda | Merge pull request #76 from luzpacello/rama-Beli _(merge)_ | `FIX-36` | Integración |
| 303 | 2026-10-02 | `fc27ba2` | Belinda | ❄ ️  fix restaurar filtro accion de auditoria | `FIX-36` | Explícita |
| 304 | 2026-10-02 | `070e96b` | Belinda | Merge pull request #77 from luzpacello/rama-Beli _(merge)_ | `FIX-36` | Integración |
| 305 | 2026-10-03 | `6572011` | Belinda | ❄ ️Maquetado y medidores de recursos del Host (CPU / RAM / Almacenamiento) | `FRN-19A` | Explícita |
| 306 | 2026-10-03 | `907efe5` | Belinda | Merge pull request #78 from luzpacello/rama-Beli _(merge)_ | `FRN-19A` | Integración |
| 307 | 2026-10-03 | `2a2db82` | Nico | feat: cliente SSE con tickets efímeros y backoff | `FRN-17C`, `FIX-61` | Inferida |
| 308 | 2026-10-03 | `88c5cb6` | Nico | Merge pull request #79 from luzpacello/Nico _(merge)_ | `FIX-61`, `FRN-17C` | Integración |
| 309 | 2026-10-03 | `8739022` | Nico | feat: events provider y consumo en tiempo real de instancias | `FRN-17A` | Inferida |
| 310 | 2026-10-03 | `15931cb` | Nico | Merge pull request #80 from luzpacello/Nico _(merge)_ | `FRN-17A` | Integración |
| 311 | 2026-10-03 | `0407fbd` | Nico | fix: alinear badge de estado de eventos con los botones del header | `FRN-19A` | Inferida |
| 312 | 2026-10-03 | `d47999d` | Nico | Merge pull request #81 from luzpacello/Nico _(merge)_ | `FRN-19A` | Integración |
| 313 | 2026-10-05 | `7f945ec` | Nico | Corregir la ruta del canal de eventos en el frontend (86e3k5204) | `FIX-45`, `FRN-17C` | Explícita |
| 314 | 2026-10-05 | `1ce64f6` | Nico | Evitar avisos repetidos en el canal de eventos y limpiar restos del canal anterior (86e3k57af) | `FIX-47`, `FRN-17A` | Explícita |
| 315 | 2026-10-05 | `c185453` | Cristian | refactor(InstanceStatus/Card/List): eliminar archivos vacios | `FRN-20A` | Inferida |
| 316 | 2026-10-05 | `0901682` | Cristian | feat(Instances): Tabla interactiva de inventario con badges de estado e IP (RF-03) | `FIX-43` | Explícita |
| 317 | 2026-10-05 | `fbd42ff` | Nico | Merge pull request #82 from luzpacello/Nico _(merge)_ | `FIX-45`, `FIX-47`, `FRN-17A`, `FRN-17C` | Integración |
| 318 | 2026-10-05 | `c8f218b` | Cristian | feat(instanceService): Tabla interactiva de inventario con badges de estado e IP (RF-03) | `FRN-20A` | Explícita |
| 319 | 2026-10-05 | `7472107` | Cristian | feat(useInstances): Tabla interactiva de inventario con badges de estado e IP (RF-03) | `FRN-20A` | Explícita |
| 320 | 2026-10-05 | `73b2b99` | Cristian | feat(InstanceAction): Tabla interactiva de inventario con badges de estado e IP (RF-03) | `FRN-20A` | Explícita |
| 321 | 2026-10-05 | `d9fcdbb` | Cristian | feat(instance): Tabla interactiva de inventario con badges de estado e IP (RF-03) | `FRN-20A` | Explícita |
| 322 | 2026-10-05 | `f69c468` | Cristian | feat(InstanceBadges): Tabla interactiva de inventario con badges de estado e IP (RF-03) | `FRN-20A` | Explícita |
| 323 | 2026-10-05 | `98deff1` | Cristian | feat(InstanceIpAddress): Tabla interactiva de inventario con badges de estado e IP (RF-03) | `FRN-20A` | Explícita |
| 324 | 2026-10-05 | `189ab69` | Cristian | feat(useCopyInstanceIp): Tabla interactiva de inventario con badges de estado e IP (RF-03) | `FRN-20A` | Explícita |
| 325 | 2026-10-05 | `04987e0` | Cristian | feat(formatInstanceTelemetry): Tabla interactiva de inventario con badges de estado e IP (RF-0… | `FRN-20A` | Explícita |
| 326 | 2026-10-05 | `7234033` | Cristian | merge: sincronizar cristian con main _(merge)_ | — | — |
| 327 | 2026-10-05 | `06458dd` | Cristian | Merge pull request #83 from luzpacello/cristian _(merge)_ | `FIX-43`, `FRN-20A` | Integración |
| 328 | 2026-10-06 | `743e4cc` | Cristian | feat(InstanceAction): cambiar estilo de acciones en instancia | `FIX-43` | Explícita |
| 329 | 2026-10-06 | `56d67bb` | Cristian | feat(Instances): cambiar estilo de tabla | `FIX-43` | Explícita |
| 330 | 2026-10-06 | `5a86dce` | Cristian | Merge pull request #84 from luzpacello/cristian _(merge)_ | `FIX-43` | Integración |
| 331 | 2026-10-06 | `bca9dd4` | Cristian | Fix(Instances): Corregir maquetado, visualización de IP, badges de estado y columnas en tabla … | `FIX-43` | Explícita |
| 332 | 2026-10-06 | `825f9c2` | Cristian | Merge pull request #85 from luzpacello/cristian _(merge)_ | `FIX-43` | Integración |
| 333 | 2026-10-06 | `dfb7800` | Cristian | feat(dropDownMenu): crear nuevo componente | `FRN-20A` | Inferida |
| 334 | 2026-10-06 | `8160c0f` | Cristian | fix(Dashboard): Reubicar botón 'Crear instancia' del Dashboard a la página de Instancias | `FIX-60` | Explícita |
| 335 | 2026-10-06 | `612cb68` | Cristian | fix(Instances): Reubicar botón 'Crear instancia' del Dashboard a la página de Instancias | `FIX-60` | Explícita |
| 336 | 2026-10-06 | `56a78df` | Cristian | Merge pull request #86 from luzpacello/cristian _(merge)_ | `FIX-60`, `FRN-20A` | Integración |
| 337 | 2026-10-06 | `50f7edb` | Cristian | fix(Dashboard): eliminar icon inutilizado | `FRN-19A` | Inferida |
| 338 | 2026-10-06 | `013badc` | Cristian | Merge pull request #87 from luzpacello/cristian _(merge)_ | `FRN-19A` | Integración |
| 339 | 2026-10-06 | `49ca405` | Cristian | feat(nativeSelected): eliminar componente | `FRN-24` | Inferida |
| 340 | 2026-10-06 | `e6a344e` | Cristian | feat(combobox): crear componente y agregarle estilo segun diseño | `FRN-24` | Inferida |
| 341 | 2026-10-06 | `741b4a0` | Cristian | refactor(combobox): agregar padding 10px a cada opcion de celda | `FRN-24` | Inferida |
| 342 | 2026-10-06 | `9407479` | Cristian | feat(input-group): necesario para el funcionamiento de componente combobox | `FRN-24` | Inferida |
| 343 | 2026-10-06 | `cfdc659` | Cristian | refactor(CrearUsuario): cambiar componente 'nativeSelected' por 'combobox' | `FRN-06` | Inferida |
| 344 | 2026-10-06 | `62c5c4b` | Cristian | refactor(rolesAndPermissions): cambiar componente 'nativeSelected' por 'combobox' | `FRN-07` | Inferida |
| 345 | 2026-10-06 | `3009fe5` | Cristian | refactor(informationOfUser): cambiar componente 'nativeSelected' por 'combobox' | `FRN-06B` | Inferida |
| 346 | 2026-10-06 | `2aa0787` | Cristian | feat(textarea): componenete necesario para funcionamiento de 'combobox' | `FRN-24` | Inferida |
| 347 | 2026-10-06 | `5be6aa9` | Cristian | Merge pull request #88 from luzpacello/cristian _(merge)_ | `FRN-06`, `FRN-06B`, `FRN-07`, `FRN-24` | Integración |
| 348 | 2026-10-06 | `493fe8a` | Cristian | refactor(combobox): cambiar color de input cuando esta seleccionado | `FRN-24` | Inferida |
| 349 | 2026-10-06 | `2111a01` | Cristian | feat(pagination): crear componente pagination con estilo segun diseño | `FRN-20A` | Inferida |
| 350 | 2026-10-06 | `d98db16` | Cristian | refactor(Instances): cambiar pagination por componente | `FRN-20A` | Inferida |
| 351 | 2026-10-06 | `5a1001c` | Cristian | refactor(Auditoria): cambiar pagination estatico por componente | `FRN-14` | Inferida |
| 352 | 2026-10-06 | `2daf259` | Cristian | refactor(Users): cambiar pagination por componente | `FRN-05` | Inferida |
| 353 | 2026-10-06 | `5258908` | Cristian | Merge pull request #89 from luzpacello/cristian _(merge)_ | `FRN-05`, `FRN-14`, `FRN-20A`, `FRN-24` | Integración |
| 354 | 2026-10-06 | `5f391f5` | Cristian | feat(Instances): actualizar actualizacion de badges | `FRN-16` | Inferida |
| 355 | 2026-10-06 | `dfbee12` | Cristian | feat(InstanceAction): actualizar actualizacion de badges segun accion de ejecucion | `FRN-16` | Inferida |
| 356 | 2026-10-06 | `297aec0` | Cristian | feat(InstanceBadges): actualizar actualizacion de badges segun accion de ejecucion | `FRN-16` | Inferida |
| 357 | 2026-10-06 | `f7338bc` | Cristian | feat(useInstances): actualizar actualizacion de badges segun accion de ejecucion | `FRN-16` | Inferida |
| 358 | 2026-10-06 | `6a859d7` | Cristian | Merge pull request #90 from luzpacello/cristian _(merge)_ | `FRN-16` | Integración |
| 359 | 2026-10-06 | `878bc39` | Cristian | feat(InstanceBadges): mostrar el cambio de estado una vez terminada la solicitud | `FRN-16` | Inferida |
| 360 | 2026-10-06 | `08fde7c` | Cristian | feat(useInstances): mostrar el cambio de estado una vez terminada la solicitud | `FRN-16` | Inferida |
| 361 | 2026-10-06 | `fc6f9ed` | Cristian | Merge pull request #91 from luzpacello/cristian _(merge)_ | `FRN-16` | Integración |
| 362 | 2026-10-06 | `f9b2c0a` | Nico | fix(instancias): resincronizar estado en progreso tras F5 y reconexión (86e3m0bw2) | `FRN-16B` | Inferida |
| 363 | 2026-10-06 | `a3a86de` | Nico | Mantener una única conexión de eventos por pestaña con suscripción selectiva (86e3kzrqy) | `FIX-64` | Inferida |
| 364 | 2026-10-06 | `e7ee40f` | Nico | Reemplazar el texto provisorio de bienvenida por la información de Centinela y del 2FA (86e3kz… | `FIX-59` | Explícita |
| 365 | 2026-10-06 | `ea64514` | Nico | Merge pull request #92 from luzpacello/Nico _(merge)_ | `FIX-59`, `FIX-64`, `FRN-16B` | Integración |
| 366 | 2026-10-06 | `3875926` | Nico | Quitar el bloque temporal de verificación de deploy del login y declarar el sitio en español | — | — |
| 367 | 2026-10-06 | `7367405` | Nico | Merge pull request #93 from luzpacello/Nico _(merge)_ | — | — |
| 368 | 2026-10-06 | `f9497e4` | Cristian | feat(createNewInstances): crear nueva pagina | `FRN-24` | Inferida |
| 369 | 2026-10-06 | `317e247` | Cristian | feat(stepper): agregar descripcion opcion | `FRN-24` | Inferida |
| 370 | 2026-10-06 | `79865c9` | Cristian | feat(applicationRoutes): crear ruta para pagina crear nueva instancia | `FRN-24` | Inferida |
| 371 | 2026-10-06 | `24ac662` | Cristian | feat(Instances): agregar funcionalidad a boton principal | `FRN-24` | Inferida |
| 372 | 2026-10-06 | `67b85c3` | Cristian | merge: sincronizar cristian con main _(merge)_ | — | — |
| 373 | 2026-10-06 | `562bdd6` | Cristian | Merge pull request #94 from luzpacello/cristian _(merge)_ | `FRN-24` | Integración |
| 374 | 2026-10-06 | `31e7cfd` | Cristian | fix(stepper): eliminar tone no utilizable | `FRN-24` | Inferida |
| 375 | 2026-10-06 | `448e415` | Cristian | Merge pull request #95 from luzpacello/cristian _(merge)_ | `FRN-24` | Integración |
| 376 | 2026-10-07 | `24315e6` | Cristian | refactor(Users): cambiar cards por tabs con funcionalidad y agregar opciones al boton 'filtro' | `FRN-05` | Inferida |
| 377 | 2026-10-07 | `83d3e04` | Cristian | refactor(Auditoria): cambiar cards por tabs con funcionalidad y agregar opciones al boton 'fil… | `FRN-14` | Inferida |
| 378 | 2026-10-07 | `9751d32` | Cristian | Merge pull request #96 from luzpacello/cristian _(merge)_ | `FRN-05`, `FRN-14` | Integración |
| 379 | 2026-10-07 | `fe16eaa` | Nico | feat(auth): nuevo diseño del login y unificación visual del flujo de autenticación (86e3kz4du) | `FIX-59` | Inferida |
| 380 | 2026-10-07 | `902878e` | Nico | refactor(auth): restaurar la clase text-info en los enlaces del flujo de autenticación (86e3kz… | `FIX-59` | Inferida |
| 381 | 2026-10-07 | `f4e89e1` | Nico | Merge pull request #97 from luzpacello/Nico _(merge)_ | `FIX-59` | Integración |
| 382 | 2026-10-07 | `e45645c` | Nico | fix(auth): tarjetas alcanzables en pantallas bajas y puntos un poco más marcados (86e3kz4du) | `FIX-59` | Inferida |
| 383 | 2026-10-07 | `14d7df9` | Nico | Merge pull request #98 from luzpacello/Nico _(merge)_ | `FIX-59` | Integración |
| 384 | 2026-10-07 | `87885e1` | Nico | fix(auth): puntos del fondo más marcados (86e3kz4du) | `FIX-58` | Inferida |
| 385 | 2026-10-07 | `2dc72ed` | Nico | Merge pull request #99 from luzpacello/Nico _(merge)_ | `FIX-58` | Integración |
| 386 | 2026-10-07 | `82d89a6` | Nico | fix(auth): subir el contraste de los puntos del fondo (86e3kz4du) | `FIX-58` | Inferida |
| 387 | 2026-10-07 | `e91b70f` | Nico | feat(auth): contorno con sombra en los iconos de beneficios del login (86e3kz4du) | `FIX-58` | Inferida |
| 388 | 2026-10-07 | `57fbf75` | Nico | Merge pull request #100 from luzpacello/Nico _(merge)_ | `FIX-58` | Integración |
| 389 | 2026-10-07 | `12d7bfa` | Nico | feat(auth): intro de video de bienvenida antes del login | `FIX-58` | Inferida |
| 390 | 2026-10-07 | `302eff7` | Nico | Merge pull request #101 from luzpacello/Nico _(merge)_ | `FIX-58` | Integración |
| 391 | 2026-10-07 | `226539e` | Nico | feat(auth): puntos del fondo con reacción al mouse y logo pulpo en el login | `FIX-58` | Inferida |
| 392 | 2026-10-07 | `218c1ad` | Nico | Merge pull request #102 from luzpacello/Nico _(merge)_ | `FIX-58` | Integración |
| 393 | 2026-10-07 | `eb6da06` | Nico | feat(sidebar): avatar por defecto según el rol y logo de la mascota | `FRN-03` | Inferida |
| 394 | 2026-10-07 | `bd11531` | Nico | Merge pull request #103 from luzpacello/Nico _(merge)_ | `FRN-03` | Integración |
| 395 | 2026-10-08 | `95e8c1b` | Belinda | ❄ integracion de node/status | `FRN-19B` | Inferida |
| 396 | 2026-10-08 | `fca4d18` | Belinda | Merge pull request #104 from luzpacello/rama-Beli _(merge)_ | `FRN-19B` | Integración |
| 397 | 2026-10-08 | `1d429f5` | Luz | 🌱 feat(instances): implementación de confirmación para las acciones operativas de las instanci… | `FRN-15` | Inferida |
| 398 | 2026-10-08 | `65a195e` | Luz | Merge pull request #105 from luzpacello/luz _(merge)_ | `FRN-15` | Integración |
| 399 | 2026-10-09 | `260e317` | Belinda | ❄ test(dashboard): cubrir errores 502/504 en telemetría del nodo | `FRN-19B` | Inferida |
| 400 | 2026-10-09 | `ee4644c` | Belinda | Merge pull request #106 from luzpacello/rama-Beli _(merge)_ | `FRN-19B` | Integración |
| 401 | 2026-10-09 | `acf79a9` | Cristian | fix(user): distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 402 | 2026-10-09 | `c3d3690` | Cristian | fix(Users): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 403 | 2026-10-09 | `b014a92` | Cristian | fix(detailsUserPage): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 404 | 2026-10-09 | `72a88bd` | Cristian | fix(informationOfUser): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 405 | 2026-10-09 | `1202028` | Cristian | fix(CrearUsuarios): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 406 | 2026-10-09 | `0ebccaa` | Cristian | fix(useCreateUser): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 407 | 2026-10-09 | `36e9a36` | Cristian | fix(userDetailsService): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 408 | 2026-10-09 | `6571be9` | Cristian | fix(userDesactivationService): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 409 | 2026-10-09 | `5ed2f43` | Cristian | fix(useUserInstanceAccess): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 410 | 2026-10-09 | `83adcc8` | Cristian | fix(userAccountMessages): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 411 | 2026-10-09 | `65482b7` | Cristian | fix(useDeletedUsers): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 412 | 2026-10-09 | `d0a1c52` | Cristian | fix(deletedUsersService): Distinguir suspensión y eliminación en contrato e interfaz | `FIX-56` | Inferida |
| 413 | 2026-10-09 | `0367985` | Cristian | Merge pull request #107 from luzpacello/cristian _(merge)_ | `FIX-56` | Integración |
| 414 | 2026-10-09 | `7e4d3af` | Luz | 🌱 fix(instances): arreglo posicionamiento del modal para la acción de eliminar instancia | `FRN-15` | Inferida |
| 415 | 2026-10-09 | `84662bd` | Luz | Merge pull request #108 from luzpacello/luz _(merge)_ | `FRN-15` | Integración |
| 416 | 2026-10-09 | `614c912` | Luz | 🌱 feat(instances): implementado la máquina de estados, bloqueo por instancia y mapeo de errores | `FRN-16`, `FIX-77` | Inferida |
| 417 | 2026-10-09 | `9679081` | Luz | Merge pull request #109 from luzpacello/luz _(merge)_ | `FIX-77`, `FRN-16` | Integración |
| 418 | 2026-10-09 | `ef485be` | Belinda | ❄ se completa los filtros reactivos RF-03 | `FRN-20B` | Inferida |
| 419 | 2026-10-09 | `7fcc0bd` | Belinda | Merge pull request #110 from luzpacello/rama-Beli _(merge)_ | `FRN-20B` | Integración |
| 420 | 2026-10-13 | `b5da499` | Nico | fix(auth): alinear textos aprobados de autenticación (86e3nrvj2) | `FIX-80` | Inferida |
| 421 | 2026-10-14 | `e7b2d1d` | Nico | fix(dashboard): quitar tarjetas simuladas y controles inertes (86e3pu2e7) | `FIX-90` | Inferida |
| 422 | 2026-10-14 | `4d0b4a7` | Nico | Merge pull request #111 from luzpacello/fix/86e3nrvj2-copy-autenticacion _(merge)_ | `FIX-80` | Integración |
| 423 | 2026-10-14 | `add947f` | Nico | Merge pull request #112 from luzpacello/fix/86e3pu2e7-limpiar-dashboard _(merge)_ | `FIX-90` | Integración |
| 424 | 2026-10-14 | `56342c1` | Cristian | feat(formatInstanceTelemetry): Columnas de CPU/RAM y estado en vivo en el Inventario (RF-03) | `FRN-20C` | Inferida |
| 425 | 2026-10-14 | `08f1f74` | Cristian | feat(useInstances): Columnas de CPU/RAM y estado en vivo en el Inventario (RF-03) | `FRN-20C` | Inferida |
| 426 | 2026-10-14 | `afd78e5` | Cristian | feat(Instances): Columnas de CPU/RAM y estado en vivo en el Inventario (RF-03) | `FRN-20C` | Inferida |
| 427 | 2026-10-14 | `9350a65` | Cristian | Merge pull request #113 from luzpacello/cristian _(merge)_ | `FRN-20C` | Integración |
| 428 | 2026-10-14 | `6b38a64` | Cristian | feat(popover): crear nuevo componente para date-picker | `FRN-31` | Inferida |
| 429 | 2026-10-14 | `d9fa6ce` | Cristian | feat(calendar): crear nuevo componente para date-picker | `FRN-31` | Inferida |
| 430 | 2026-10-14 | `0361380` | Cristian | feat(date-picker): crear nuevo componente | `FRN-31` | Inferida |
| 431 | 2026-10-14 | `4a51f64` | Cristian | build(package): instalacion de dependencias para componente calendar | `FRN-31` | Inferida |
| 432 | 2026-10-14 | `9b97227` | Cristian | build(pnpm-lock): instalacion de dependencias para componente calendar | `FRN-31` | Inferida |
| 433 | 2026-10-14 | `a8073f7` | Cristian | fix(Auditoria): Reactividad de filtros y etiquetas accesibles en vista de Auditoría | `FIX-74` | Explícita |
| 434 | 2026-10-14 | `d9c2825` | Cristian | fix(date-picker): Reactividad de filtros y etiquetas accesibles en vista de Auditoría | `FIX-74` | Explícita |
| 435 | 2026-10-14 | `7e3f43d` | Cristian | Merge pull request #114 from luzpacello/cristian _(merge)_ | `FIX-74`, `FRN-31` | Integración |
| 436 | 2026-10-14 | `5694c7a` | Cristian | fix(date-picker): cambiar nombre sin seleccionar | `FRN-31` | Inferida |
| 437 | 2026-10-14 | `002ead3` | Cristian | fix(Users): cambiar ancho del dropdowmenu | `FRN-05` | Inferida |
| 438 | 2026-10-14 | `0aa6f8a` | Cristian | Merge pull request #115 from luzpacello/cristian _(merge)_ | `FRN-05`, `FRN-31` | Integración |
| 439 | 2026-10-14 | `576f65a` | Cristian | feat(badge): Agregar variantes a badge | `FRN-20C` | Inferida |
| 440 | 2026-10-14 | `e85cfd3` | Cristian | fix(useUserRecentActivity): Quitar los datos hardcodeados de la ficha de usuario | `FIX-87` | Explícita |
| 441 | 2026-10-14 | `95b542d` | Cristian | fix(userRecentActivityService): Quitar los datos hardcodeados de la ficha de usuario | `FIX-87` | Explícita |
| 442 | 2026-10-14 | `232e317` | Cristian | fix(detailsUserPages): Quitar los datos hardcodeados de la ficha de usuario | `FIX-87` | Explícita |
| 443 | 2026-10-14 | `757580a` | Cristian | fix(Auditoria): Quitar los datos hardcodeados de la ficha de usuario | `FIX-87` | Explícita |
| 444 | 2026-10-14 | `44d1e3a` | Cristian | Merge pull request #116 from luzpacello/cristian _(merge)_ | `FIX-87`, `FRN-20C` | Integración |
| 445 | 2026-10-14 | `fdcfc22` | Cristian | fix(userInstanceService): agregar persistencia de datos para opciones null | `FIX-87` | Inferida |
| 446 | 2026-10-14 | `7dca556` | Cristian | feat(detailsUserPage): crear tabla y opcion visual si no cuenta con instancias asignadas | `FIX-87` | Inferida |
| 447 | 2026-10-14 | `9697a9f` | Cristian | Merge pull request #117 from luzpacello/cristian _(merge)_ | `FIX-87` | Integración |
| 448 | 2026-10-14 | `4e666a0` | Cristian | fix(Instances):  Ruta  con asistente estático sin lógica | `FIX-89` | Inferida |
| 449 | 2026-10-14 | `bfcc5f4` | Cristian | Merge pull request #118 from luzpacello/cristian _(merge)_ | `FIX-89` | Integración |
| 450 | 2026-10-15 | `57c739f` | Nico | fix(dashboard): reponer paneles y sumar curvas, campana, actividad y alertas (86e3pu2e7) | `FIX-100` | Explícita |
| 451 | 2026-10-15 | `6db04d3` | Nico | Merge pull request #119 from luzpacello/fix/86e3pu2e7-alinear-dashboard _(merge)_ | `FIX-100` | Integración |
| 452 | 2026-10-15 | `6c8ecc9` | Nico | feat(dashboard): desglosar el conteo de instancias por tipo de recurso (86e3nrn90) | `FRN-19C` | Inferida |
| 453 | 2026-10-15 | `405791e` | Nico | fix(dashboard): refrescar el conteo y la actividad en vivo, y alinear el pie y las curvas (86e… | `FIX-83` | Explícita |
| 454 | 2026-10-15 | `3466a46` | Nico | Merge pull request #120 from luzpacello/fix/86e3ptewk-refresco-reactivo _(merge)_ | `FIX-83`, `FRN-19C` | Integración |
| 455 | 2026-10-15 | `2b120f4` | Luz | 🌱 feat(instances): implementación de desbloqueo reactivo y notificaciones toast adoptativas (R… | `FRN-17B` | Inferida |
| 456 | 2026-10-15 | `3c9a90b` | Luz | Merge pull request #121 from luzpacello/luz _(merge)_ | `FRN-17B` | Integración |
| 457 | 2026-10-15 | `34cf00f` | Cristian | fix(SideBar): cambiar iconos de instancias 'Dashboard' y 'Auditoria' | `FRN-03` | Inferida |
| 458 | 2026-10-15 | `b4a34bc` | Cristian | feat(InstanceCreationSummary): crear componente de resumen de creacion de instancia | `FRN-24` | Inferida |
| 459 | 2026-10-15 | `6b85a79` | Cristian | feat(createInstanceOptions): crear datos estaticos para prueba de steps | `FRN-24` | Inferida |
| 460 | 2026-10-15 | `4029e64` | Cristian | feat(useCreateInstanceWizard): Agregar funcionalidades al Wizard | `FRN-24` | Inferida |
| 461 | 2026-10-15 | `944d0d2` | Cristian | feat(createInstance): crear estructura para creacion de instancia | `FRN-24` | Inferida |
| 462 | 2026-10-15 | `fda22d7` | Cristian | feat(selectTypeOfInstances): crear contenido de tipo de instancia | `FRN-24` | Inferida |
| 463 | 2026-10-15 | `38f215e` | Cristian | feat(createNewInstances): agregarle funcionalidad a la pagina y componentes | `FRN-24` | Inferida |
| 464 | 2026-10-15 | `bf1d581` | Cristian | Merge pull request #122 from luzpacello/cristian _(merge)_ | `FRN-03`, `FRN-24` | Integración |
| 465 | 2026-10-15 | `7e40c4f` | Belinda | ❄ se agrega filtro por IP y reubica filtros | `FRN-20B` | Inferida |
| 466 | 2026-10-15 | `1a0cdba` | Belinda | Merge pull request #123 from luzpacello/rama-Beli _(merge)_ | `FRN-20B` | Integración |
| 467 | 2026-10-15 | `f61374a` | Belinda | ❄ alinear filtros | `FRN-20B` | Inferida |
| 468 | 2026-10-15 | `cba3327` | Cristian | feat(configureInstanceResources): crear contenido para la seleccion de recursos para la instan… | `FRN-24` | Inferida |
| 469 | 2026-10-15 | `91d1676` | Cristian | feat(createInstanceOptions): agregar texto estatico para simular carga de recursos y eliminar … | `FRN-24` | Inferida |
| 470 | 2026-10-15 | `93937be` | Cristian | feat(InstanceCreationSummary): agregar funcion para calcular memoria en GB | `FRN-24` | Inferida |
| 471 | 2026-10-15 | `c9fdb6b` | Cristian | feat(createInstance): agregar persistencia de datos de almacenamiento | `FRN-24` | Inferida |
| 472 | 2026-10-15 | `3fa0277` | Cristian | feat(useCreateInstanceWizard): agregar paso 2 | `FRN-24` | Inferida |
| 473 | 2026-10-15 | `385a52b` | Cristian | merge: sincronizar cristian con main _(merge)_ | — | — |
| 474 | 2026-10-15 | `1fc14f4` | Cristian | fix(selectedTypeOfInstances): envolver formulario en card | `FRN-24` | Inferida |
| 475 | 2026-10-15 | `049f16d` | Belinda | Merge branch 'main' of https://github.com/luzpacello/centinela-front into rama-beli _(merge)_ | — | — |
| 476 | 2026-10-15 | `99aa5dc` | Cristian | fix(createNewInstances): eliminar fondo blanco para todas las opciones | `FRN-24` | Inferida |
| 477 | 2026-10-15 | `8a28b88` | Cristian | Merge pull request #124 from luzpacello/cristian _(merge)_ | `FRN-24` | Integración |
| 478 | 2026-10-15 | `e966818` | Cristian | fix(steppers): cambiar estilo para que entre el texto largo del subtitulo | `FRN-24` | Inferida |
| 479 | 2026-10-15 | `8116e6f` | Cristian | fix(createNewInstances): agregar pagina a paso 2 | `FRN-24` | Inferida |
| 480 | 2026-10-15 | `597e35b` | Cristian | Merge pull request #125 from luzpacello/cristian _(merge)_ | `FRN-24` | Integración |
| 481 | 2026-10-15 | `ee15fb6` | Belinda | ❄ mensajes traducidos | `FIX-80` | Inferida |
| 482 | 2026-10-15 | `6147bc6` | Belinda | Merge pull request #126 from luzpacello/rama-Beli _(merge)_ | `FIX-80`, `FRN-20B` | Integración |
| 483 | 2026-10-15 | `7271e1a` | Belinda | ❄ se corrige duplicidad de texto saludable inaccesible | `FIX-78` | Inferida |
| 484 | 2026-10-15 | `fd16cce` | Belinda | Merge pull request #127 from luzpacello/rama-Beli _(merge)_ | `FIX-78` | Integración |
| 485 | 2026-10-16 | `623bc28` | Belinda | ❄ se Renderiza botones  accesibles visibles en la cabecera | `FIX-101`, `FIX-79` | Explícita |
| 486 | 2026-10-16 | `f57d46f` | Belinda | Merge pull request #128 from luzpacello/rama-Beli _(merge)_ | `FIX-101`, `FIX-79` | Integración |
| 487 | 2026-10-16 | `fb58781` | Luz | 🌱 fix(instances): desacoplar guarda de nivelAcceso para rol ADMIN en acciones de energía | `FIX-76` | Explícita |
| 488 | 2026-10-16 | `c01dd36` | Luz | 🌱 test(instances): integro pruebas para modales | `FIX-76` | Explícita |
| 489 | 2026-10-16 | `6bb563f` | Luz | Merge pull request #129 from luzpacello/luz _(merge)_ | `FIX-76` | Integración |
| 490 | 2026-10-16 | `73d2fc7` | Belinda | ❄ FIX: Traducción de errores 403 de instancia | `FIX-86` | Explícita |
| 491 | 2026-10-16 | `b005f16` | Belinda | Merge branch 'main' into rama-Beli _(merge)_ | — | — |
| 492 | 2026-10-16 | `a3c5e91` | Belinda | ❄ chore: actualizar instancias, tests y config | `FIX-86` | Inferida |
| 493 | 2026-10-16 | `2731d64` | Belinda | ❄ fix(instances): corregir tipado de getErrorMessage al eliminar instancia | `FIX-86` | Inferida |
| 494 | 2026-10-16 | `28a4513` | Belinda | Merge pull request #130 from luzpacello/rama-Beli _(merge)_ | `FIX-86` | Integración |
| 495 | 2026-10-16 | `b6e94da` | Cristian | fix(Auditoria): cambiar boton de exportacion a tipo primario y cambiar de posicion | `FRN-31` | Inferida |
| 496 | 2026-10-16 | `a42ecf1` | Cristian | fix(Instances): crear ruta para crear una nueva instancia | `FIX-89` | Inferida |
| 497 | 2026-10-16 | `30894f4` | Cristian | fix(createNewInstance): Advertencia visible Próximamente en vista directa /instances/new | `FIX-104`, `FIX-99` | Explícita |
| 498 | 2026-10-16 | `68fc4b8` | Cristian | Merge pull request #131 from luzpacello/cristian _(merge)_ | `FIX-104`, `FIX-89`, `FIX-99`, `FRN-31` | Integración |
| 499 | 2026-10-16 | `608d04b` | Cristian | fix(Instances): cambiar tabs y filtro | `FIX-101` | Inferida |
| 500 | 2026-10-16 | `2b1720c` | Cristian | fix(InstanceFiltersBar): cambiar tabs y filtro | `FIX-101` | Inferida |
| 501 | 2026-10-16 | `10ce750` | Cristian | Merge pull request #132 from luzpacello/cristian _(merge)_ | `FIX-101` | Integración |
| 502 | 2026-10-16 | `b464b7e` | Cristian | fix(instance-filters-bar): acomodar test con nuevo texto de filtros | `FIX-101` | Inferida |
| 503 | 2026-10-16 | `15245a0` | Cristian | Merge pull request #133 from luzpacello/cristian _(merge)_ | `FIX-101` | Integración |
| 504 | 2026-10-16 | `773c4b0` | Cristian | fix(Auditoria): cambiar diseño de tabla y representacion de contador de eventos exitosos/ fall… | `FRN-31` | Inferida |
| 505 | 2026-10-16 | `143ddda` | Cristian | Merge pull request #134 from luzpacello/cristian _(merge)_ | `FRN-31` | Integración |
| 506 | 2026-10-16 | `e85a10f` | Cristian | fix(formatInstanceTelemetry): Normalización de telemetría de inventario (escala 0-100, desacop… | `FIX-97` | Explícita |
| 507 | 2026-10-16 | `338d658` | Cristian | fix(Instances): Normalización de telemetría de inventario (escala 0-100, desacople de stopped … | `FIX-97` | Explícita |
| 508 | 2026-10-16 | `cd8a1f4` | Cristian | Merge pull request #135 from luzpacello/cristian _(merge)_ | `FIX-97` | Integración |
| 509 | 2026-10-16 | `041b911` | Luz | 🌱 fix (instance): corregir cancelación de la confirmación de start en los modales de energía | `FIX-85` | Explícita |
| 510 | 2026-10-16 | `46cd4b3` | Luz | Merge pull request #136 from luzpacello/luz _(merge)_ | `FIX-85` | Integración |
| 511 | 2026-10-18 | `f641c5f` | Luz | 🌱feat(instance): integración de vista de detalle y navegación | `FRN-21` | Explícita |
| 512 | 2026-10-18 | `94f4db8` | Cristian | fea(NetworkModeOption): crear paso de red en creacion de instancias | `FRN-25` | Inferida |
| 513 | 2026-10-18 | `8a8bb78` | Cristian | feat(configureInstanceNetwork): crear paso de Red, en creacion de instancias | `FRN-25` | Inferida |
| 514 | 2026-10-18 | `39f5c3a` | Cristian | feat(createInstance): crear paso de Red, en creacion de instancias | `FRN-25` | Inferida |
| 515 | 2026-10-18 | `4a59cbb` | Luz | Merge pull request #137 from luzpacello/luz _(merge)_ | `FRN-21` | Integración |
| 516 | 2026-10-18 | `4d603aa` | Cristian | feat(createInstanceOptions): crear paso de Red, en creacion de instancias | `FRN-25` | Inferida |
| 517 | 2026-10-18 | `74cb5d1` | Cristian | feat(useCreateInstanceWizard): crear paso de Red, en creacion de instancias | `FRN-25` | Inferida |
| 518 | 2026-10-18 | `21b8522` | Cristian | feat(createNewInstances): crear paso de Red, en creacion de instancias | `FRN-25` | Inferida |
| 519 | 2026-10-18 | `ea4f1c9` | Cristian | Merge pull request #138 from luzpacello/cristian _(merge)_ | `FRN-25` | Integración |
| 520 | 2026-10-18 | `c51dc9d` | Cristian | feat(selectInstanceOperatingSystem): crear paso de SO, en creacion de instancias | `FRN-25` | Inferida |
| 521 | 2026-10-18 | `6cd2797` | Cristian | feat(creteInstanceOptions): crear paso de SO, en creacion de instancias | `FRN-25` | Inferida |
| 522 | 2026-10-18 | `dd203f6` | Cristian | feat(createInstance): crear paso de SO, en creacion de instancias | `FRN-25` | Inferida |
| 523 | 2026-10-18 | `332e177` | Cristian | feat(createNewInstances): crear paso de SO, en creacion de instancias | `FRN-25` | Inferida |
| 524 | 2026-10-18 | `ce676cc` | Cristian | Merge pull request #139 from luzpacello/cristian _(merge)_ | `FRN-25` | Integración |
| 525 | 2026-10-18 | `bd62aa9` | Cristian | feat(reviewInstanceConfiguration): crear paso de Resumen, en creacion de instancias | `FRN-25` | Inferida |
| 526 | 2026-10-18 | `7fe3218` | Cristian | feat(createInstance): crear paso de Resumen, en creacion de instancias | `FRN-25` | Inferida |
| 527 | 2026-10-18 | `a4191b9` | Cristian | feat(useCreateInstanceWizard): crear paso de Resumen, en creacion de instancias | `FRN-25` | Inferida |
| 528 | 2026-10-18 | `0664b94` | Cristian | feat(instanceService): crear paso de Resumen, en creacion de instancias | `FRN-25` | Inferida |
| 529 | 2026-10-18 | `3510240` | Cristian | feat(createNewInstances): crear paso de Resumen, en creacion de instancias | `FRN-25` | Inferida |
| 530 | 2026-10-18 | `9e0683b` | Cristian | Merge branch 'main' into cristian _(merge)_ | — | — |
| 531 | 2026-10-18 | `3d840ad` | Cristian | Merge pull request #140 from luzpacello/cristian _(merge)_ | `FRN-25` | Integración |
| 532 | 2026-10-18 | `48a4406` | Nico | fix(dashboard): sacar el ruido del canal de eventos de la actividad reciente | `FRN-19C` | Inferida |
| 533 | 2026-10-18 | `54efc3d` | Nico | feat(dashboard): cargar el histórico del nodo al entrar en el dashboard | `FIX-103`, `FRN-19C` | Explícita |
| 534 | 2026-10-18 | `d3643e0` | Nico | Merge remote-tracking branch 'origin/main' into fix/actividad-reciente-sin-ruido _(merge)_ | — | — |
| 535 | 2026-10-18 | `1913ca0` | Nico | fix(dashboard): alinear los detalles de auditoría con el contrato en objeto | `FRN-31` | Inferida |
| 536 | 2026-10-18 | `91eabf4` | Nico | Merge pull request #141 from luzpacello/fix/actividad-reciente-sin-ruido _(merge)_ | `FIX-103`, `FRN-19C`, `FRN-31` | Integración |

---

## 4. Cruce por tarea: asignado declarado vs. autor real

Incluye las tareas frontend documentadas (FRN/SEC/LOGIN) y las referenciadas por commits.  
**Veredicto:** ✅ Coincide · ⚠️ Parcial (hay commits de asignados y de terceros) · ❌ No coincide (ningún commit del asignado) · ➖ Sin commits · ❓ Sin asignación.

| Tarea | Título | Área | Asignado(s) | Autoría real | Veredicto |
|---|---|---|---|---|---|
| `LOGIN-02` | Pantalla y validación del código TOTP | Frontend | Belinda y Luz | Luz (2), Belinda (1), Nico (1) | ⚠️ Parcial (3/4) |
| `LOGIN-03` | Integración y prueba del prototipo | Frontend / Backend | Cristian, Tayra y Lisandro | — | ➖ Sin commits |
| `LOGIN-04` | Prueba integral de autenticación y autorización | Frontend / Backend | Cristian, Tayra y Lisandro | — | ➖ Sin commits |
| `FRN-01` | Maquetado del formulario de login | Frontend | Belinda | Cristian (25), Belinda (2), Luz (2) | ⚠️ Parcial (2/29) |
| `FRN-02` | Validación en cliente del login | Frontend | Belinda | Cristian (7) | ❌ No coincide |
| `FRN-03` | Navbar y rutas base | Frontend | Cristian y Luz | Cristian (22), Belinda (13), Nico (6), Luz (2) | ⚠️ Parcial (24/43) |
| `FRN-04` | Conexión del login con localhost | Frontend | Cristian | Cristian (5), Nico (1) | ⚠️ Parcial (5/6) |
| `FRN-05` | Panel de gestión de usuarios | Frontend | Belinda | Cristian (12), Nico (2) | ❌ No coincide |
| `FRN-06` | Modal de creación y desactivación de usuarios | Frontend | Luz | Cristian (17), Nico (1) | ❌ No coincide |
| `FRN-06B` | Modal de edición de usuario y cambio de rol | Frontend | Luz | Cristian (7) | ❌ No coincide |
| `FRN-07` | Selector de asignación de instancias | Frontend | Cristian | Cristian (13) | ✅ Coincide |
| `FRN-08` | Interceptor HTTP y manejo de `403 Forbidden` | Frontend | Cristian | Cristian (6) | ✅ Coincide |
| `FRN-09` | Vinculación 2FA mediante QR | Frontend | Belinda y Luz | Belinda (1) | ✅ Coincide |
| `FRN-10` | Cambio obligatorio de contraseña temporal | Frontend | Belinda | Nico (1) | ❌ No coincide |
| `FRN-11` | Acciones administrativas de recuperación | Frontend | Luz | Luz (1) | ✅ Coincide |
| `FRN-12` | Vistas de recuperación de contraseña (RF-13) | Frontend | Belinda | Belinda (1), Luz (1) | ⚠️ Parcial (1/2) |
| `FRN-13` | Flujo integral de logout y limpieza de sesión en cliente (Frontend) | Frontend | Cristian y Belinda | Cristian (10) | ✅ Coincide |
| `FRN-14` | Suite de pruebas unitarias y de integración para la vista de Auditorí… | Frontend | Belinda / Luz | Cristian (4), Nico (2), Belinda (1) | ⚠️ Parcial (1/7) |
| `FRN-15` | Modales de confirmación antierror para acciones operativas | Frontend | Belinda | Luz (2) | ❌ No coincide |
| `FRN-16` | Máquina de estados "Operación en progreso" por instancia | Frontend | Cristian | Cristian (6), Luz (1) | ⚠️ Parcial (6/7) |
| `FRN-16B` | (`BRG-04-FRN`) - Resincronización del estado "Operación en progreso" … | Frontend | Cristian | Nico (1) | ❌ No coincide |
| `FRN-17A` | Consumo de eventos en tiempo real y distribución por instancia | Frontend | Cristian | Nico (3) | ❌ No coincide |
| `FRN-17B` | Desbloqueo reactivo de instancia y notificación Toast adaptativa | Frontend | Cristian | Luz (1) | ❌ No coincide |
| `FRN-17C` | (`BRG-02-FRN`) - Cliente de eventos con solicitud previa de ticket ef… | Frontend | Cristian | Nico (2) | ❌ No coincide |
| `FRN-18` | (`SEC-04B`) - Integración en Frontend de niveles de acceso por instan… | Frontend | Cristian y Belinda | Cristian (5) | ✅ Coincide |
| `FRN-19A` | (ex `FRN-13A`) - Maquetado y medidores de recursos del Host (CPU / RA… | Frontend | Belinda (PR #78) | Belinda (1), Nico (1), Cristian (1) | ⚠️ Parcial (1/3) |
| `FRN-19B` | (ex `FRN-13B`) - Semáforo de salud global e integración con `GET /api… | Frontend | Belinda | Belinda (2) | ✅ Coincide |
| `FRN-19C` | & `FIX-83` - Tarjetas de conteo de VMs/LXC y refresco ante `TASK_FINI… | Frontend | Belinda | Nico (3) | ❌ No coincide |
| `FRN-20A` | (ex `FRN-14A`) - Tabla interactiva de inventario con badges de estado… | Frontend | Luz (implementación inicial de Cristian en PR #74) | Cristian (12) | ❌ No coincide |
| `FRN-20B` | (ex `FRN-14B`) - Filtros reactivos por tipo, estado y buscador dinámi… | Frontend | Luz | Belinda (3) | ❌ No coincide |
| `FRN-20C` | (`BRG-05-FRN`, parte Inventario) - Columnas de CPU/RAM y estado en vi… | Frontend | Luz | Cristian (4) | ❌ No coincide |
| `FRN-21` | (`BRG-06-FRN`) / `FIX-95` - Maquetado de la vista Detalle de Instanci… | Frontend | Cristian / Belinda | Luz (1) | ❌ No coincide |
| `FRN-22` | Hook `useInstanceMetrics` y cliente WebSocket para telemetría en tiem… | Frontend | Belinda | — | ➖ Sin commits |
| `FRN-23` | Componentes de gráficos interactivos de rendimiento y selector de ran… | Frontend | Luz | — | ➖ Sin commits |
| `FRN-24` | Wizard de creación: Paso 1 (Tipo de instancia) y Paso 2 (Recursos y C… | Frontend | Belinda | Cristian (26) | ❌ No coincide |
| `FRN-25` | Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Pas… | Frontend | Luz | Cristian (15) | ❌ No coincide |
| `FRN-26` | Envío de creación, manejo de errores de cuota (`409`) y seguimiento e… | Frontend | Cristian | — | ➖ Sin commits |
| `FRN-27` | Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / … | Frontend | Cristian | — | ➖ Sin commits |
| `FRN-28` | Integración de edición en la página de Detalle de Instancia y refresc… | Frontend | Cristian | — | ➖ Sin commits |
| `FRN-29` | Pestaña "Snapshots" en Detalle de Instancia y formulario de nuevo sna… | Frontend | Luz | — | ➖ Sin commits |
| `FRN-30` | Modal destructivo de Rollback y eliminación con confirmación antierror | Frontend | Cristian | — | ➖ Sin commits |
| `FRN-31` | Vista administrativa de Auditoría con filtros reactivos completos | Frontend | Belinda | Cristian (8), Nico (1) | ❌ No coincide |
| `FRN-32` | Componente y flujo de exportación de reportes de auditoría | Frontend | Luz | — | ➖ Sin commits |
| `FRN-33` | Centro de Notificaciones y banners contextuales de alerta por saturac… | Frontend | Belinda | — | ➖ Sin commits |
| `FRN-34` | Adaptación del cliente a cookies seguras de sesión y auditoría de acc… | Frontend | Cristian | — | ➖ Sin commits |
| `FRN-36` | (`BRG-10-FRN`) - Modularización de `InstanceDetailPage` para Pestaña … | Frontend | Cristian | — | ➖ Sin commits |
| `SEC-02` | Cliente frontend compatible con refresh token HttpOnly | Frontend | Cristian | Cristian (10) | ✅ Coincide |
| `SEC-03` | Contexto y sistema reactivo de permisos en Frontend (UI/UX RBAC + Res… | Frontend | Cristian y Belinda | Luz (4) | ❌ No coincide |
| `FIX-06` | Modal de alta y edición de usuarios sin funcionalidad real (`FRN-06` … | Frontend | Luz | Cristian (10), Nico (1) | ❌ No coincide |
| `FIX-07` | El frontend lee `body.code` en vez de `errorCode` | Frontend | Cristian | Nico (2) | ❌ No coincide |
| `FIX-08` | Restaurar protección de rutas y guards de sesión (`FRN-03` / `FRN-05`) | Frontend | Belinda | Belinda (1), Nico (1) | ⚠️ Parcial (1/2) |
| `FIX-12` | Falta guard de rol administrativo en rutas de gestión de usuarios (`F… | Frontend | Belinda / Cristian | Nico (2) | ❌ No coincide |
| `FIX-14` | Integración frontend del selector de instancias (`FRN-07`) | Frontend | Cristian | Cristian (11) | ✅ Coincide |
| `FIX-18` | Conexión de `RecoverPassword.tsx` a la API y alineación de tests (`FR… | Frontend | Belinda | Luz (1) | ❌ No coincide |
| `FIX-19` | Acciones de Restablecimiento Administrativo de Contraseña y 2FA en `U… | Frontend | Luz / Cristian | Luz (1) | ✅ Coincide |
| `FIX-20` | Sincronización de políticas de complejidad y distinción UX entre Camb… | Frontend | Belinda | Luz (1) | ❌ No coincide |
| `FIX-21` | Opción de cerrar sesión en el cambio obligatorio de contraseña (`FRN-… | Frontend | Belinda | Cristian (1) | ❌ No coincide |
| `FIX-22` | Guard administrativo en la ruta `/auditoria` (`FRN-14` / `BAC-18`) (F… | Frontend | Belinda / Luz | Belinda (2) | ✅ Coincide |
| `FIX-24` | Confirmación del alta sin contraseña temporal y baja de usuario (`FRN… | Frontend | Luz | Cristian (1) | ❌ No coincide |
| `FIX-25` | Persistir la edición del perfil de usuario (`FRN-06B` / `FIX-14`) (Fr… | Frontend | Luz / Cristian | Cristian (12) | ✅ Coincide |
| `FIX-27` | Mensaje específico cuando falla el envío del correo en el alta (`FIX-… | Frontend | Cristian | Cristian (1) | ✅ Coincide |
| `FIX-28` | Regresión: el logout del frontend no envía el access token y la sesió… | Frontend | Cristian | Cristian (2) | ✅ Coincide |
| `FIX-29` | Validar la mayúscula en `validatePasswordComplexity` (`FIX-20`) (Fron… | Frontend | Luz | Luz (1) | ✅ Coincide |
| `FIX-30` | Helper `canOperateInstance` en `usePermissions()` (`FRN-18`) (Fronten… | Frontend | Cristian y Belinda | Cristian (5) | ✅ Coincide |
| `FIX-36` | Regresión: restaurar el filtro "Acción" en Auditoría (`FRN-14` / RF-0… | Frontend | Belinda y Luz (PR #77) | Belinda (2), Cristian (1) | ⚠️ Parcial (2/3) |
| `FIX-38` | "Solo lectura" (`READ_ONLY`) usado como rol de usuario (`SEC-03` / `F… | Frontend | Cristian y Belinda | Cristian (8) | ✅ Coincide |
| `FIX-41` | Regresión: `requiredRole="Admin"` oculta "Roles y permisos" al admini… | Frontend | Luz (autora de `1156eaa`) | Cristian (4), Luz (1) | ⚠️ Parcial (1/5) |
| `FIX-42` | `useAuth().canAccessInstance` no revisa el rol (`SEC-03`) (Frontend) | Frontend | Luz (autora de `1156eaa` y `47e5c19`) | Cristian (4), Luz (2) | ⚠️ Parcial (2/6) |
| `FIX-43` | Corregir maquetado, visualización de IP, badges de estado y columnas … | Frontend | Luz / Cristian | Cristian (4) | ✅ Coincide |
| `FIX-45` | El cliente de eventos abre el stream en una ruta que el backend no si… | Frontend | Nico | Nico (1) | ✅ Coincide |
| `FIX-47` | Deduplicar eventos por `id` y limpiar restos del canal anterior (`FRN… | Frontend | Nico | Nico (1) | ✅ Coincide |
| `FIX-56` | Distinguir suspensión y eliminación en contrato e interfaz | Frontend | — | Cristian (12) | ❓ Sin asignación |
| `FIX-58` | Reemplazar placeholder textual del logo por el imagotipo/logotipo ofi… | Frontend / UI | Belinda / Luz | Nico (5) | ❌ No coincide |
| `FIX-59` | Sustituir el placeholder "imagenes y informacion random" por mensaje … | Frontend / UX Copy | Luz / Belinda | Nico (4) | ❌ No coincide |
| `FIX-60` | Reubicar el botón "Crear instancia" del Dashboard a la página de Inst… | Frontend / UI | Cristian / Luz | Cristian (2) | ✅ Coincide |
| `FIX-61` | Cerrar sesión y detener reconexiones ante 401 al solicitar el ticket … | Frontend | Cristian | Nico (1) | ❌ No coincide |
| `FIX-64` | Conexión SSE compartida y suscripción selectiva (`FRN-17A`) (Frontend) | Frontend | Cristian | Nico (1) | ❌ No coincide |
| `FIX-74` | Reactividad de filtros y etiquetas accesibles en vista de Auditoría (… | Frontend | Belinda / Luz | Cristian (2) | ❌ No coincide |
| `FIX-76` | & `FIX-85` - Modales de energía y confirmación `Start` desacoplada pa… | Frontend | Belinda / Luz | Luz (2) | ✅ Coincide |
| `FIX-77` | Mapeo de errores HTTP de ciclo de vida (`BAC-29`) y resiliencia de la… | Frontend | Cristian / Belinda | Luz (1) | ❌ No coincide |
| `FIX-78` | Desduplicación de indicadores de salud del nodo y accesibilidad en Da… | Frontend | Belinda / Luz | Belinda (1) | ✅ Coincide |
| `FIX-79` | Exposición directa de controles para selector reactivo por tipo de re… | Frontend | Luz | Belinda (1) | ❌ No coincide |
| `FIX-80` | Alineación textual y normalización de copy en vistas de autenticación… | Frontend / UI | Belinda / Nico | Nico (1), Belinda (1) | ✅ Coincide |
| `FIX-83` | Tarjetas de conteo de VMs/LXC y refresco ante TASK_FINISHED | Frontend | Belinda | Nico (1) | ❌ No coincide |
| `FIX-85` | Modales de energía y confirmación Start desacoplada para ADMIN | Frontend | Belinda / Luz | Luz (1) | ✅ Coincide |
| `FIX-86` | Completar la máquina de estados y el mapeo D2 de operaciones (`FIX-77… | Frontend | Cristian | Belinda (3) | ❌ No coincide |
| `FIX-87` | Quitar los datos hardcodeados de la ficha de usuario (`FRN-05` / `FRN… | Frontend | Cristian / Belinda | Cristian (6) | ✅ Coincide |
| `FIX-89` | Ruta `/instances/new` con asistente estático sin lógica (`FIX-60`) (F… | Frontend | Luz / Belinda | Cristian (2) | ❌ No coincide |
| `FIX-90` | Quitar los placeholders muertos del Dashboard (`FRN-19A` / `FRN-19B`)… | Frontend | Belinda | Nico (1) | ❌ No coincide |
| `FIX-91` | Eliminar el componente muerto `InstanceTester` (`FRN-15` / `FRN-16`) … | Frontend | Cristian | — | ➖ Sin commits |
| `FIX-97` | Normalización de telemetría de inventario (escala porcentual 0-100 y … | Frontend | Cristian / Luz | Cristian (2) | ✅ Coincide |
| `FIX-99` | Advertencia visible "Próximamente" en vista directa `/instances/new` … | Frontend | Cristian | Cristian (1) | ✅ Coincide |
| `FIX-100` | Manejo seguro de valores `undefined` en `formatCount` del Dashboard (… | Frontend | Luz / Cristian | Nico (1) | ❌ No coincide |
| `FIX-101` | Corrección de filtros reactivos en memoria (búsqueda de texto y alter… | Frontend | Luz | Cristian (3), Belinda (1) | ❌ No coincide |
| `FIX-102` | Corrección de atributo de accesibilidad `aria-label` en `InstanceActi… | Frontend | Belinda / Luz | — | ➖ Sin commits |
| `FIX-103` | Manejo defensivo de arrays de telemetría de nodo en `useNodeStatus.ts… | Frontend | Belinda | Nico (1) | ❌ No coincide |
| `FIX-104` | Visibilidad permanente del banner de advertencia "Próximamente" en to… | Frontend | Cristian | Cristian (1) | ✅ Coincide |
| `FIX-104B` | Desacople de UPID y adopción del contrato seguro `tareaId` en Fronten… | Frontend | Cristian / Luz | — | ➖ Sin commits |
| `FIX-105` | Desacople de UPID y adopción del contrato seguro `tareaId` en Fronten… | Frontend | Cristian / Luz | — | ➖ Sin commits |
| `FIX-108` | Adaptación del Frontend para consumo de identificadores de instancia … | Frontend | Belinda / Cristian | — | ➖ Sin commits |
| `FIX-111B` | Modal de confirmación con re-autenticación (Password + 2FA) en cambio… | Frontend | Cristian | — | ➖ Sin commits |
| `FIX-115B` | Actualización de validación y formularios de contraseñas de alta long… | Frontend | Luz | — | ➖ Sin commits |
| `FIX-117B` | Adaptación del Dashboard de operador ante segregación de métricas glo… | Frontend | Belinda / Cristian | — | ➖ Sin commits |

### 4.1. Resumen del cruce

| Resultado | Tareas | % |
|---|---:|---:|
| ✅ Coincide | 29 | 28.2% |
| ⚠️ Parcial | 12 | 11.7% |
| ❌ No coincide | 40 | 38.8% |
| ➖ Sin commits | 21 | 20.4% |
| ❓ Sin asignación | 1 | 1.0% |
| **Total** | **103** | **100%** |

---

## 5. Discrepancias detectadas

### 5.1. Tareas donde el asignado no escribió ningún commit (❌)

| Tarea | Título | Asignado(s) | Autor real | Commits | Evidencia (hash) |
|---|---|---|---|---:|---|
| `FRN-24` | Wizard de creación: Paso 1 (Tipo de instancia) y Paso 2 (Re… | Belinda | Cristian (26) | 26 | `49ca405`, `e6a344e`, `741b4a0`, `9407479`, `2aa0787`, `493fe8a`, `f9497e4`, `317e247`, `79865c9`, `24ac662`, `31e7cfd`, `b4a34bc`, `6b85a79`, `4029e64`, `944d0d2`, `fda22d7`, `38f215e`, `cba3327`, `91d1676`, `93937be`, `c9fdb6b`, `3fa0277`, `1fc14f4`, `99aa5dc`, `e966818`, `8116e6f` |
| `FRN-06` | Modal de creación y desactivación de usuarios | Luz | Cristian (17), Nico (1) | 18 | `7fff0a0`, `e4a9798`, `5ad9b88`, `a387e10`, `684b09c`, `c92009e`, `a36a8c9`, `6198cf9`, `bd23180`, `fa6cbb7`, `1fb8bfd`, `8db47cf`, `597ba06`, `f290dad`, `b13931c`, `58d15c0`, `21610fb`, `cfdc659` |
| `FRN-25` | Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y … | Luz | Cristian (15) | 15 | `94f4db8`, `8a8bb78`, `39f5c3a`, `4d603aa`, `74cb5d1`, `21b8522`, `c51dc9d`, `6cd2797`, `dd203f6`, `332e177`, `bd62aa9`, `7fe3218`, `a4191b9`, `0664b94`, `3510240` |
| `FRN-05` | Panel de gestión de usuarios | Belinda | Cristian (12), Nico (2) | 14 | `0dc7d8b`, `64cad86`, `9ccef0e`, `03209be`, `dace934`, `fdab248`, `7a520c1`, `9c1f768`, `4c3fcf3`, `c75acce`, `bac88cc`, `2daf259`, `24315e6`, `002ead3` |
| `FRN-20A` | (ex `FRN-14A`) - Tabla interactiva de inventario con badges… | Luz (implementación inicial de Cristian en PR #74) | Cristian (12) | 12 | `c185453`, `c8f218b`, `7472107`, `73b2b99`, `d9fcdbb`, `f69c468`, `98deff1`, `189ab69`, `04987e0`, `dfb7800`, `2111a01`, `d98db16` |
| `FIX-06` | Modal de alta y edición de usuarios sin funcionalidad real … | Luz | Cristian (10), Nico (1) | 11 | `e4a9798`, `6198cf9`, `bd23180`, `fa6cbb7`, `1fb8bfd`, `8db47cf`, `597ba06`, `f290dad`, `b13931c`, `58d15c0`, `21610fb` |
| `FRN-31` | Vista administrativa de Auditoría con filtros reactivos com… | Belinda | Cristian (8), Nico (1) | 9 | `6b38a64`, `d9fa6ce`, `0361380`, `4a51f64`, `9b97227`, `5694c7a`, `b6e94da`, `773c4b0`, `1913ca0` |
| `FRN-02` | Validación en cliente del login | Belinda | Cristian (7) | 7 | `bde4238`, `393012f`, `55e3db6`, `368a4ce`, `22d498a`, `25d13a1`, `15f5a22` |
| `FRN-06B` | Modal de edición de usuario y cambio de rol | Luz | Cristian (7) | 7 | `e00a1ab`, `6f9cdfa`, `8687564`, `dc84dfd`, `cc8a8a3`, `e5da01b`, `3009fe5` |
| `FIX-58` | Reemplazar placeholder textual del logo por el imagotipo/lo… | Belinda / Luz | Nico (5) | 5 | `87885e1`, `82d89a6`, `e91b70f`, `12d7bfa`, `226539e` |
| `FIX-101` | Corrección de filtros reactivos en memoria (búsqueda de tex… | Luz | Cristian (3), Belinda (1) | 4 | `623bc28`, `608d04b`, `2b1720c`, `b464b7e` |
| `FIX-59` | Sustituir el placeholder "imagenes y informacion random" po… | Luz / Belinda | Nico (4) | 4 | `e7ee40f`, `fe16eaa`, `902878e`, `e45645c` |
| `FRN-20C` | (`BRG-05-FRN`, parte Inventario) - Columnas de CPU/RAM y es… | Luz | Cristian (4) | 4 | `56342c1`, `08f1f74`, `afd78e5`, `576f65a` |
| `SEC-03` | Contexto y sistema reactivo de permisos en Frontend (UI/UX … | Cristian y Belinda | Luz (4) | 4 | `1156eaa`, `353f0f1`, `952b43c`, `47e5c19` |
| `FIX-86` | Completar la máquina de estados y el mapeo D2 de operacione… | Cristian | Belinda (3) | 3 | `73d2fc7`, `a3c5e91`, `2731d64` |
| `FRN-17A` | Consumo de eventos en tiempo real y distribución por instan… | Cristian | Nico (3) | 3 | `c8557b6`, `8739022`, `1ce64f6` |
| `FRN-19C` | & `FIX-83` - Tarjetas de conteo de VMs/LXC y refresco ante … | Belinda | Nico (3) | 3 | `6c8ecc9`, `48a4406`, `54efc3d` |
| `FRN-20B` | (ex `FRN-14B`) - Filtros reactivos por tipo, estado y busca… | Luz | Belinda (3) | 3 | `ef485be`, `7e40c4f`, `f61374a` |
| `FIX-07` | El frontend lee `body.code` en vez de `errorCode` | Cristian | Nico (2) | 2 | `63dc6de`, `efeb4a5` |
| `FIX-12` | Falta guard de rol administrativo en rutas de gestión de us… | Belinda / Cristian | Nico (2) | 2 | `64cad86`, `e47e041` |
| `FIX-74` | Reactividad de filtros y etiquetas accesibles en vista de A… | Belinda / Luz | Cristian (2) | 2 | `a8073f7`, `d9c2825` |
| `FIX-89` | Ruta `/instances/new` con asistente estático sin lógica (`F… | Luz / Belinda | Cristian (2) | 2 | `4e666a0`, `a42ecf1` |
| `FRN-15` | Modales de confirmación antierror para acciones operativas | Belinda | Luz (2) | 2 | `1d429f5`, `7e4d3af` |
| `FRN-17C` | (`BRG-02-FRN`) - Cliente de eventos con solicitud previa de… | Cristian | Nico (2) | 2 | `2a2db82`, `7f945ec` |
| `FIX-100` | Manejo seguro de valores `undefined` en `formatCount` del D… | Luz / Cristian | Nico (1) | 1 | `57c739f` |
| `FIX-103` | Manejo defensivo de arrays de telemetría de nodo en `useNod… | Belinda | Nico (1) | 1 | `54efc3d` |
| `FIX-18` | Conexión de `RecoverPassword.tsx` a la API y alineación de … | Belinda | Luz (1) | 1 | `fe6d784` |
| `FIX-20` | Sincronización de políticas de complejidad y distinción UX … | Belinda | Luz (1) | 1 | `fe6d784` |
| `FIX-21` | Opción de cerrar sesión en el cambio obligatorio de contras… | Belinda | Cristian (1) | 1 | `b87ea1b` |
| `FIX-24` | Confirmación del alta sin contraseña temporal y baja de usu… | Luz | Cristian (1) | 1 | `3ea7fa4` |
| `FIX-61` | Cerrar sesión y detener reconexiones ante 401 al solicitar … | Cristian | Nico (1) | 1 | `2a2db82` |
| `FIX-64` | Conexión SSE compartida y suscripción selectiva (`FRN-17A`)… | Cristian | Nico (1) | 1 | `a3a86de` |
| `FIX-77` | Mapeo de errores HTTP de ciclo de vida (`BAC-29`) y resilie… | Cristian / Belinda | Luz (1) | 1 | `614c912` |
| `FIX-79` | Exposición directa de controles para selector reactivo por … | Luz | Belinda (1) | 1 | `623bc28` |
| `FIX-83` | Tarjetas de conteo de VMs/LXC y refresco ante TASK_FINISHED | Belinda | Nico (1) | 1 | `405791e` |
| `FIX-90` | Quitar los placeholders muertos del Dashboard (`FRN-19A` / … | Belinda | Nico (1) | 1 | `e7b2d1d` |
| `FRN-10` | Cambio obligatorio de contraseña temporal | Belinda | Nico (1) | 1 | `893bf09` |
| `FRN-16B` | (`BRG-04-FRN`) - Resincronización del estado "Operación en … | Cristian | Nico (1) | 1 | `f9b2c0a` |
| `FRN-17B` | Desbloqueo reactivo de instancia y notificación Toast adapt… | Cristian | Luz (1) | 1 | `2b120f4` |
| `FRN-21` | (`BRG-06-FRN`) / `FIX-95` - Maquetado de la vista Detalle d… | Cristian / Belinda | Luz (1) | 1 | `f641c5f` |

### 5.2. Tareas con autoría compartida / parcial (⚠️)

| Tarea | Título | Asignado(s) | Autoría real | % del asignado | Evidencia (hash) |
|---|---|---|---|---:|---|
| `FRN-01` | Maquetado del formulario de login | Belinda | Cristian (25), Belinda (2), Luz (2) | 2/29 (7%) | `851a18e`, `452c492`, `d0cb6e6`, `79dacbd`, `366004a`, `4bb8c95`, `2117570`, `26aa8eb`, `2b81804`, `427714f`, `0144f42`, `15feecd`, `06af276`, `4ba0083`, `f2c3955`, `bde4238`, `393012f`, `55e3db6`, `0bfaa81`, `27c26cc`, `db00e97`, `a9a8929`, `9601e47`, `3074fe6`, `6d63cf7`, `51e7106`, `5d36e22`, `6b1ccc4`, `bc51bdb` |
| `FRN-14` | Suite de pruebas unitarias y de integración para la vista d… | Belinda / Luz | Cristian (4), Nico (2), Belinda (1) | 1/7 (14%) | `cbf7481`, `e134383`, `924d65f`, `56b88f5`, `e6ad7ad`, `5a1001c`, `83d3e04` |
| `FIX-41` | Regresión: `requiredRole="Admin"` oculta "Roles y permisos"… | Luz (autora de `1156eaa`) | Cristian (4), Luz (1) | 1/5 (20%) | `1156eaa`, `3269dbd`, `f5e3c61`, `a351564`, `88f31b6` |
| `FIX-42` | `useAuth().canAccessInstance` no revisa el rol (`SEC-03`) (… | Luz (autora de `1156eaa` y `47e5c19`) | Cristian (4), Luz (2) | 2/6 (33%) | `1156eaa`, `47e5c19`, `fb60cb2`, `31f468e`, `d32b218`, `ee60382` |
| `FRN-19A` | (ex `FRN-13A`) - Maquetado y medidores de recursos del Host… | Belinda (PR #78) | Belinda (1), Nico (1), Cristian (1) | 1/3 (33%) | `6572011`, `0407fbd`, `50f7edb` |
| `FIX-08` | Restaurar protección de rutas y guards de sesión (`FRN-03` … | Belinda | Belinda (1), Nico (1) | 1/2 (50%) | `6b08566`, `231779f` |
| `FRN-12` | Vistas de recuperación de contraseña (RF-13) | Belinda | Belinda (1), Luz (1) | 1/2 (50%) | `4f2763e`, `fe6d784` |
| `FRN-03` | Navbar y rutas base | Cristian y Luz | Cristian (22), Belinda (13), Nico (6), Luz (2) | 24/43 (56%) | `08a7fcf`, `ed1d3d8`, `2bad8ab`, `34f631b`, `672fed1`, `35d22d9`, `33c3ed4`, `a1e85eb`, `b803b20`, `ee4ba82`, `ac5abdb`, `eaea038`, `45db606`, `c674703`, `d315882`, `22d498a`, `25d13a1`, `15f5a22`, `b15c588`, `3112eea`, `0b50cd6`, `da396a1`, `27b4467`, `d9e0339`, `a4d3f65`, `7bc6dc5`, `ab5e0f9`, `635c8e7`, `998da0c`, `65e1dc1`, `995b2d6`, `76e7ebe`, `b50057e`, `b1d5acc`, `c8e5592`, `a2da502`, `131ee7d`, `a387e10`, `684b09c`, `c92009e`, `81c6424`, `eb6da06`, `34cf00f` |
| `FIX-36` | Regresión: restaurar el filtro "Acción" en Auditoría (`FRN-… | Belinda y Luz (PR #77) | Belinda (2), Cristian (1) | 2/3 (67%) | `56b88f5`, `2057fc0`, `fc27ba2` |
| `LOGIN-02` | Pantalla y validación del código TOTP | Belinda y Luz | Luz (2), Belinda (1), Nico (1) | 3/4 (75%) | `ea92e52`, `feaaf27`, `6145740`, `ca52e47` |
| `FRN-04` | Conexión del login con localhost | Cristian | Cristian (5), Nico (1) | 5/6 (83%) | `b4785dc`, `368a4ce`, `efb495a`, `1f62bfe`, `52ea961`, `432d747` |
| `FRN-16` | Máquina de estados "Operación en progreso" por instancia | Cristian | Cristian (6), Luz (1) | 6/7 (86%) | `5f391f5`, `dfbee12`, `297aec0`, `f7338bc`, `878bc39`, `08fde7c`, `614c912` |

### 5.3. Tareas frontend sin commits en el submódulo (➖)

Pueden estar pendientes o implementadas fuera de `frontend/` (por ejemplo, pruebas en `test/front`).

| Tarea | Título | Área | Asignado(s) | Estado inferido |
|---|---|---|---|---|
| `LOGIN-03` | Integración y prueba del prototipo | Frontend / Backend | Cristian, Tayra y Lisandro | Se ejecuta en `test/front` (fuera del submódulo) |
| `LOGIN-04` | Prueba integral de autenticación y autorización | Frontend / Backend | Cristian, Tayra y Lisandro | Se ejecuta en `test/front` (fuera del submódulo) |
| `FRN-22` | Hook `useInstanceMetrics` y cliente WebSocket para telemetría en tiem… | Frontend | Belinda | Pendiente / sin implementación en `frontend/` |
| `FRN-23` | Componentes de gráficos interactivos de rendimiento y selector de ran… | Frontend | Luz | Pendiente / sin implementación en `frontend/` |
| `FRN-26` | Envío de creación, manejo de errores de cuota (`409`) y seguimiento e… | Frontend | Cristian | Pendiente / sin implementación en `frontend/` |
| `FRN-27` | Modal de reconfiguración de vCPU y RAM con advertencia de reinicio / … | Frontend | Cristian | Pendiente / sin implementación en `frontend/` |
| `FRN-28` | Integración de edición en la página de Detalle de Instancia y refresc… | Frontend | Cristian | Pendiente / sin implementación en `frontend/` |
| `FRN-29` | Pestaña "Snapshots" en Detalle de Instancia y formulario de nuevo sna… | Frontend | Luz | Pendiente / sin implementación en `frontend/` |
| `FRN-30` | Modal destructivo de Rollback y eliminación con confirmación antierror | Frontend | Cristian | Pendiente / sin implementación en `frontend/` |
| `FRN-32` | Componente y flujo de exportación de reportes de auditoría | Frontend | Luz | Pendiente / sin implementación en `frontend/` |
| `FRN-33` | Centro de Notificaciones y banners contextuales de alerta por saturac… | Frontend | Belinda | Pendiente / sin implementación en `frontend/` |
| `FRN-34` | Adaptación del cliente a cookies seguras de sesión y auditoría de acc… | Frontend | Cristian | Pendiente / sin implementación en `frontend/` |
| `FRN-36` | (`BRG-10-FRN`) - Modularización de `InstanceDetailPage` para Pestaña … | Frontend | Cristian | Pendiente / sin implementación en `frontend/` |
| `FIX-91` | Eliminar el componente muerto `InstanceTester` (`FRN-15` / `FRN-16`) … | Frontend | Cristian | Pendiente / sin implementación en `frontend/` |
| `FIX-102` | Corrección de atributo de accesibilidad `aria-label` en `InstanceActi… | Frontend | Belinda / Luz | Pendiente / sin implementación en `frontend/` |
| `FIX-104B` | Desacople de UPID y adopción del contrato seguro `tareaId` en Fronten… | Frontend | Cristian / Luz | Pendiente / sin implementación en `frontend/` |
| `FIX-105` | Desacople de UPID y adopción del contrato seguro `tareaId` en Fronten… | Frontend | Cristian / Luz | Pendiente / sin implementación en `frontend/` |
| `FIX-108` | Adaptación del Frontend para consumo de identificadores de instancia … | Frontend | Belinda / Cristian | Pendiente / sin implementación en `frontend/` |
| `FIX-111B` | Modal de confirmación con re-autenticación (Password + 2FA) en cambio… | Frontend | Cristian | Pendiente / sin implementación en `frontend/` |
| `FIX-115B` | Actualización de validación y formularios de contraseñas de alta long… | Frontend | Luz | Pendiente / sin implementación en `frontend/` |
| `FIX-117B` | Adaptación del Dashboard de operador ante segregación de métricas glo… | Frontend | Belinda / Cristian | Pendiente / sin implementación en `frontend/` |

### 5.4. Tareas sin asignación declarada (❓)

- `FIX-56` — Distinguir suspensión y eliminación en contrato e interfaz · autor real: **Cristian (12)** · evidencia: `acf79a9`, `c3d3690`, `b014a92`, `72a88bd`, `1202028`, `0ebccaa`, `36e9a36`, `6571be9`, `5ed2f43`, `83adcc8`, `65482b7`, `d0a1c52`

### 5.5. Observaciones de metadata que afectan el cruce

- **`FRN-20A`** está asignada a **Luz**, pero la ficha aclara “implementación inicial de Cristian en PR #74”. Se contabiliza como asignada solo a Luz; los 12 commits son de Cristian.
- **`FRN-21` / `FIX-95`**: la vista de detalle fue implementada por **Luz** (`f641c5f`), aunque la tarea figura asignada a Cristian/Belinda.
- **`FRN-17A/B/C`, `FRN-19C`, `FIX-83`, `FIX-90`, `FIX-103`**: el trabajo de eventos y dashboard lo hizo **Nico**, mientras la documentación lo asigna a Cristian/Belinda.
- **`FRN-24` y `FRN-25`** (wizard) figuran asignadas a Belinda y Luz, pero los 41 commits del wizard son de **Cristian**.
- **`FIX-56`** no tiene asignado en la documentación, aunque **Cristian** escribió 12 commits de la interfaz de suspensión/eliminación.
- **`FIX-14` y `FIX-61`** no tienen commits propios: su implementación viaja dentro de los commits de `FRN-07` y del cliente SSE (`FRN-17C`), respectivamente. Se mapearon allí para no reportarlas como pendientes.
- **`FRN-32`** no se marca implementada: el único commit relacionado (`b6e94da`) solo cambia el botón de exportación dentro de la vista de Auditoría (`FRN-31`).
- Algunas tareas figuran como “Frontend” pero su verificación vive en `test/front` (`LOGIN-03`, `LOGIN-04`), por lo que no tienen commits en este submódulo.

---

## 6. Conclusiones

- De **103 tareas frontend** analizadas, **29** coinciden exactamente con el asignado, **12** son parciales, **40** no tienen ningún commit del asignado, **21** no registran commits en el submódulo y **1** sin asignación declarada.
- **52 tareas (50%)** tienen al menos un commit de una persona distinta a la asignada.

**Commits en tareas ajenas:**

| Persona | Tareas ajenas | Commits en tareas ajenas |
|---|---:|---:|
| **Cristian** | 21 | 166 |
| **Nico** | 25 | 46 |
| **Belinda** | 5 | 21 |
| **Luz** | 10 | 15 |

**Casos más relevantes (tareas 100% ajenas y con volumen):**

- **`FRN-24`** — Wizard de creación: Paso 1 (Tipo de instancia) y Paso 2 (Recursos y Cuotas en Vivo): asignada a **Belinda**, ejecutada por **Cristian (26)** (26 commits).
- **`FRN-06`** — Modal de creación y desactivación de usuarios: asignada a **Luz**, ejecutada por **Cristian (17), Nico (1)** (18 commits).
- **`FRN-25`** — Wizard de creación: Paso 3 (Almacenamiento, Template/ISO y Red) y Paso 4 (Resumen): asignada a **Luz**, ejecutada por **Cristian (15)** (15 commits).
- **`FRN-05`** — Panel de gestión de usuarios: asignada a **Belinda**, ejecutada por **Cristian (12), Nico (2)** (14 commits).
- **`FRN-20A`** — (ex `FRN-14A`) - Tabla interactiva de inventario con badges de estado e IP (`RF-03`) (Frontend): asignada a **Luz**, ejecutada por **Cristian (12)** (12 commits).
- **`FIX-06`** — Modal de alta y edición de usuarios sin funcionalidad real (`FRN-06` / `FRN-06B`): asignada a **Luz**, ejecutada por **Cristian (10), Nico (1)** (11 commits).
- **`FRN-31`** — Vista administrativa de Auditoría con filtros reactivos completos: asignada a **Belinda**, ejecutada por **Cristian (8), Nico (1)** (9 commits).
- **`FRN-02`** — Validación en cliente del login: asignada a **Belinda**, ejecutada por **Cristian (7)** (7 commits).
- **`FRN-06B`** — Modal de edición de usuario y cambio de rol: asignada a **Luz**, ejecutada por **Cristian (7)** (7 commits).
- **`FIX-58`** — Reemplazar placeholder textual del logo por el imagotipo/logotipo oficial vectorizado (`FRN-01` / `LOGIN-02`) (Frontend): asignada a **Belinda, Luz**, ejecutada por **Nico (5)** (5 commits).
- **`FIX-101`** — Corrección de filtros reactivos en memoria (búsqueda de texto y alternancia en ejecución / detenidas) (`FIX-79` / `FRN-20B`) (Frontend): asignada a **Luz**, ejecutada por **Cristian (3), Belinda (1)** (4 commits).
- **`FIX-59`** — Sustituir el placeholder "imagenes y informacion random" por mensaje informativo y amigable sobre Centinela y 2FA (`FRN-01` / `LOGIN-02`) (Frontend / UX Copy): asignada a **Belinda, Luz**, ejecutada por **Nico (4)** (4 commits).

**Lectura por persona:**

- **Cristian** es el autor dominante del frontend (271 de 378 commits reales) y ejecutó tareas asignadas a Belinda y Luz: `FRN-01`, `FRN-05`, `FRN-06`, `FRN-06B`, `FRN-20A`, `FRN-24`, `FRN-25`, `FRN-31`.
- **Belinda** hizo el layout/base y buena parte de Auditoría, Dashboard y filtros, pero varias de sus tareas asignadas (`FRN-05`, `FRN-06B`, `FRN-19C`, `FIX-83`, `FIX-90`, `FIX-103`) las terminó Cristian o Nico.
- **Nico** absorbió eventos SSE y dashboard (`FRN-17A/B/C`, `FRN-19C`, `FIX-83`, `FIX-90`, `FIX-103`), además de `FIX-58/59` y `FIX-64`, tareas de Cristian/Belinda/Luz.
- **Luz** implementó `FRN-21`, `SEC-03` y parte de `FIX-41/42`, tareas asignadas a Cristian/Belinda, y dejó su huella principal en `FRN-13`, `FRN-17B` y `FRN-15`.

**Recomendaciones:**

1. **Reasignar en ClickUp** las tareas de 5.1 o convertirlas en co-asignadas, especialmente el bloque wizard (`FRN-24`/`FRN-25`) y el dashboard/eventos (`FRN-17x`, `FRN-19C`).
2. **Convención de commits con ID de tarea** (`FRN-xx`, `FIX-xx`) en el asunto o pie, como ya se usó en varios commits (`FIX-76`, `FIX-101`, `FRN-19C`).
3. **Separar tareas de diseño de tareas de implementación**: hoy `FRN-01` mezcla 29 commits de Cristian, Belinda y Luz bajo una asignación única a Belinda.
4. **Reconciliar pendientes reales** antes de darlas por terminadas: `FRN-22`, `FRN-23`, `FRN-26`–`FRN-30`, `FRN-32`–`FRN-34`, `FRN-36`, `FIX-91`, `FIX-102`, `FIX-104B`, `FIX-105`, `FIX-108`, `FIX-111B`, `FIX-115B` y `FIX-117B` no tienen implementación en `frontend/`.
5. **Definir ownership por módulo** (auth, usuarios, dashboard, instancias/wizard, auditoría, eventos) para reducir el traspaso implícito de tareas.

---

## 7. Commits críticos y su impacto

| Commit | Autor | Tarea | Por qué es crítico |
|---|---|---|---|
| `bde4238`, `55e3db6` | Cristian | `FRN-01/FRN-02` | Primeras páginas de Login y SignUp. |
| `b4785dc`, `368a4ce`, `35d22d9` | Cristian | `FRN-02/FRN-03/FRN-04` | Base de autenticación, login y router del frontend. |
| `e02b821`, `4556331`, `607e455`, `deb59cb`, `b900d19`, `72e205d`, `300a0bd`, `cce5a55`, `4062fdc` | Cristian | `FRN-13/FIX-28` | Flujo integral de logout y limpieza de sesión. |
| `1156eaa` | Luz | `SEC-03/FIX-41/FIX-42` | Contexto y sistema reactivo de permisos; origen de las regresiones FIX-41/42. |
| `2a2db82` | Nico | `FRN-17C/FIX-61` | Cliente SSE con tickets efímeros, backoff y corte por 401. |
| `7f945ec` | Nico | `FIX-45` | Corrige la ruta del canal de eventos que el backend no servía. |
| `0901682`, `743e4cc`, `56d67bb`, `bca9dd4` | Cristian | `FRN-20A/FIX-43` | Tabla de inventario con badges, IP y estilos. |
| `8160c0f`, `612cb68` | Cristian | `FIX-60` | Reubica el botón Crear instancia del Dashboard a Instancias. |
| `e7ee40f`, `fe16eaa`, `14d7df9`, `12d7bfa`, `226539e` | Nico | `FIX-58/FIX-59` | Rediseño, copy, logo y fondo del login. |
| `614c912` | Luz | `FRN-16/FIX-77` | Máquina de estados y mapeo de errores de operaciones. |
| `57c739f`, `405791e` | Nico | `FIX-100/FIX-83` | Reposición de paneles y refresco reactivo del Dashboard. |
| `30894f4` | Cristian | `FIX-99/FIX-104` | Banner Próximamente en /instances/new. |
| `f641c5f` | Luz | `FRN-21` | Vista de detalle de instancia y navegación desde inventario. |
| `94f4db8`, `8a8bb78`, `39f5c3a`, `4d603aa`, `74cb5d1`, `21b8522`, `c51dc9d`, `6cd2797`, `dd203f6`, `332e177`, `bd62aa9`, `7fe3218`, `a4191b9`, `0664b94`, `3510240` | Cristian | `FRN-25` | Pasos de Red, SO y Resumen del wizard (tarea asignada a Luz). |
| `54efc3d` | Nico | `FRN-19C/FIX-103` | Histórico del nodo en Dashboard; introduce la regresión serie.map. |

---

## Anexo A. Commits sin tarea asociada

Son bootstrap, CI/CD, documentación o merges de sincronización sin tarea propia.

| # | Fecha | Hash | Autor | Asunto | Motivo |
|---:|---|---|---|---|---|
| 1 | 2026-08-17 | `c04eddb` | Luz | Initial commit | Documentación |
| 4 | 2026-08-18 | `98f21df` | Cristian | docs(README): pasos para instalar | Documentación |
| 6 | 2026-08-19 | `923649f` | Cristian | Merge pull request #2 from luzpacello/cristian | Merge / sincronización |
| 41 | 2026-09-11 | `da65fb6` | Belinda | Merge branch 'main' into rama-Beli | Merge / sincronización |
| 59 | 2026-09-15 | `3bfbccc` | Nico | ci: deploy automatico del front a CT 103 (PRUEBAS) en cada push a main | CI/CD del frontend |
| 60 | 2026-09-15 | `86b9a02` | Luz | Merge pull request #19 from luzpacello/ci/deploy-automatico-front | Merge / sincronización |
| 67 | 2026-09-15 | `873d183` | Cristian | Merge branch 'main' into cristian | Merge / sincronización |
| 69 | 2026-09-15 | `d8a36a5` | Nico | feat(layout_auth): mostrar fecha y hora del deploy en el panel | CI/CD del frontend |
| 70 | 2026-09-15 | `70e59a2` | Nico | ci(deploy): el deploy del front solo se dispara desde main | CI/CD del frontend |
| 71 | 2026-09-15 | `7a05a16` | Belinda | Merge pull request #22 from luzpacello/feat/stamp-deploy-panel | Merge / sincronización |
| 72 | 2026-09-15 | `566b5ab` | Belinda | Merge pull request #23 from luzpacello/ci/solo-main-despliega | Merge / sincronización |
| 73 | 2026-09-15 | `2f779e7` | Nico | infra(deploy-check): bloque temporal que muestra front y back juntos | CI/CD del frontend |
| 74 | 2026-09-15 | `26a927e` | Belinda | Merge pull request #24 from luzpacello/infra/deploy-check-front-back | Merge / sincronización |
| 75 | 2026-09-15 | `c259f20` | Nico | infra(deploy-check): el sello del front incluye el commit | CI/CD del frontend |
| 76 | 2026-09-15 | `dc605af` | Belinda | Merge pull request #25 from luzpacello/infra/front-stamp-con-commit | Merge / sincronización |
| 89 | 2026-09-16 | `86efc92` | Cristian | merge: sincronizar cristian con main | Merge / sincronización |
| 112 | 2026-09-20 | `c7a8800` | Belinda | Merge branch 'main' into rama-Beli | Merge / sincronización |
| 150 | 2026-09-21 | `d9cd6c1` | Nico | ci(front): agrego workflow de build en pull_request | CI/CD del frontend |
| 175 | 2026-09-22 | `d539253` | Belinda | Merge pull request #48 from luzpacello/main | Merge / sincronización |
| 176 | 2026-09-22 | `94f2d7b` | Cristian | merge: sincronizar cristian con main | Merge / sincronización |
| 185 | 2026-09-22 | `24695fa` | Cristian | merge: resolver sincronizacion pendiente | Merge / sincronización |
| 187 | 2026-09-22 | `f07da88` | Belinda | Merge branch 'main' of https://github.com/luzpacello/centinela-front into rama-… | Merge / sincronización |
| 189 | 2026-09-22 | `4e5fbab` | Belinda | ❄️Merge de actualizaciones remotas | Merge / sincronización |
| 192 | 2026-09-22 | `44e57e6` | Belinda | Merge pull request #52 from luzpacello/main | Merge / sincronización |
| 271 | 2026-10-01 | `92c7286` | Luz | Merge branch 'main' into luz | Merge / sincronización |
| 301 | 2026-10-02 | `cdbc8e9` | Belinda | Merge branch 'main' into rama-Beli | Merge / sincronización |
| 326 | 2026-10-05 | `7234033` | Cristian | merge: sincronizar cristian con main | Merge / sincronización |
| 366 | 2026-10-06 | `3875926` | Nico | Quitar el bloque temporal de verificación de deploy del login y declarar el sit… | CI/CD del frontend |
| 367 | 2026-10-06 | `7367405` | Nico | Merge pull request #93 from luzpacello/Nico | Merge / sincronización |
| 372 | 2026-10-06 | `67b85c3` | Cristian | merge: sincronizar cristian con main | Merge / sincronización |
| 473 | 2026-10-15 | `385a52b` | Cristian | merge: sincronizar cristian con main | Merge / sincronización |
| 475 | 2026-10-15 | `049f16d` | Belinda | Merge branch 'main' of https://github.com/luzpacello/centinela-front into rama-… | Merge / sincronización |
| 491 | 2026-10-16 | `b005f16` | Belinda | Merge branch 'main' into rama-Beli | Merge / sincronización |
| 530 | 2026-10-18 | `9e0683b` | Cristian | Merge branch 'main' into cristian | Merge / sincronización |
| 534 | 2026-10-18 | `d3643e0` | Nico | Merge remote-tracking branch 'origin/main' into fix/actividad-reciente-sin-ruido | Merge / sincronización |

## Anexo B. Limitaciones

1. **Fuente de asignación:** se usó el campo `Asignado/Asignada` de la documentación del repositorio (espejo declarado de ClickUp). Si ClickUp cambió después del 19/10/2026, el cruce debe reejecutarse contra su API.
2. **Inferencia:** los commits marcados como `Inferida` no citan la tarea. La relación se basa en asunto, archivos y módulo; pueden existir errores de frontera entre tareas vecinas (p. ej. `FRN-05` vs `FRN-06B`, `FRN-14` vs `FRN-31`).
3. **Solapamiento:** una misma implementación puede servir a dos tareas (p. ej. `FRN-06` y `FIX-06`); el commit se cuenta en ambas y por eso la suma de commits por tarea supera el total.
4. **Merges:** se les asigna la unión de tareas de la rama integrada, pero no suman a la autoría efectiva.
5. **Autoría de Git:** no detecta pair programming, código copiado de otra rama ni commits reescritos por merges.
6. **Trabajo fuera del submódulo:** tareas frontend pueden haberse validado en `test/front`; esos cambios no están en `frontend/`.

---

_Informe generado automáticamente a partir de `git log` de `frontend/` y de la documentación del repositorio. Los hashes cortos permiten localizar el commit completo con `git show <hash>`._
