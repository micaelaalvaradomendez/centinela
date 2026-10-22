# El Centinela: informe de evolución (Etapa Base → mitad de la Etapa 1)

**Versión 2 · 06/10/2026**
**Punto de partida:** `Resumen_Fase_Base.pdf` (previo al inicio de la Etapa 1)
**Fuentes:** `terminado.md`, `terminado-1.md`, `actual.md`, `requerimientos.md`, `centinela.md`; contexto adicional de `etapa1.md`, `futuro.md`, `futuro-1.md` y `test/informe.md`.

> **Cómo leer este informe.**
> - **Verificado** = está en `terminado.md` o `terminado-1.md`. Cuenta como avance aunque tenga un FIX pendiente: "implementada con problema" significa que la funcionalidad existe y se probó; lo que falta suele ser documentación, un código de error o un detalle de interfaz.
> - **Proyectado** = está en `actual.md`. Por pedido, el informe describe el sistema **como si `actual.md` estuviera implementado** (cada tarea cumpliendo su criterio de éxito). Esa parte es un escenario de planificación, no un resultado de pruebas.
> - **Limitación:** en esta sesión no tengo acceso al código de los submódulos `frontend/` y `backend/`, solo a la documentación y a las pruebas de aceptación que ella describe. La "revisión de la implementación" se apoya en lo que esos documentos registran como verificado (commits, casos de prueba, resultados), no en una lectura directa del código.
> - **Corrección respecto de la v1:** allí traté "implementada con problema" como casi no hecha, y eso subestimaba el avance. Acá se cuenta como lo que es: trabajo sustancial ya entregado.

---

## 1. Resumen ejecutivo

1. **Etapa Base cerrada.** Todo lo que el PDF dejaba pendiente se completó y verificó. `LOGIN-04` pasó **10/10** con el frontend real contra el backend real (01/10/2026).
2. **Etapa 1: avance verificado de 12,5 h de 50,5 (≈25 %)** en 8 de 29 tareas, **más** todo el trabajo previo a la etapa (simulador de Proxmox, `FIX-32`/`33`/`35`, `BAC-28`) y los FIX transversales (`FIX-37`, `FIX-39`, `FIX-40`).
3. **Con `actual.md` implementado, el avance asciende a 28 h de 50,5 (≈55 %), 15 de 29 tareas.** Ese es el punto "mitad de la etapa": la **Ola 1 completa**, más energía (`BAC-24A/B`), modales y estado "operación en progreso", semáforo del nodo y reanudación de tareas.
4. **Se destraba el camino crítico.** Con `FIX-39`, `FIX-37`, `SEC-03`, `FIX-38` y `FIX-40` cerrados, el frontend de energía ya no depende de la base.
5. **Lo que queda es la segunda mitad:** enriquecer telemetría e inventario (IP, CPU, RAM), timeout y reintentos de UPID, auditoría del resultado, limpieza de permisos al borrar, filtros, resincronización tras F5 y la integración final local y en servidor (`INT-01` a `INT-03`).

---

## 2. Etapa Base: lo pendiente del PDF, hoy

| Pendiente según el PDF | Estado | Fuente |
|---|---|---|
| Recuperación de contraseña (SEC-01, FRN-11, FRN-12), cambio obligatorio | ✅ | `terminado.md` |
| Niveles `READ_ONLY`/`FULL_ACCESS` (SEC-03, SEC-04, FRN-18) | ✅ (con `FIX-30`, `FIX-37`, `FIX-38`, `FIX-41`, `FIX-42`) | `terminado.md` |
| Redis (`INF-06`) | ✅ en local (`INF-06A`, `BAC-17A`); servidor: `INF-06B` | `terminado.md` |
| SMTP real (`INF-08`, `BAC-16B`) | ✅ en local (`INF-08B`, `BAC-16B`, verificado con Mailpit en STARTTLS y autenticación contra el relay); servidor: `INF-08A` | `terminado.md` |
| CORS y TLS (`INF-05`, `FIX-31`) | ✅ CORS y Nginx del borde versionado (`docker/nginx-edge.conf`) | `terminado.md` |
| Sesión única y Redis (`BAC-17B`) | ✅ 1 sesión = 1 fila, probado con 10 refresh | `terminado.md` |
| Auditoría append-only (`BAC-18`, `FIX-23`) | ✅ trigger + `REVOKE`; `UPDATE`/`DELETE`/`TRUNCATE` rechazados por el motor | `terminado.md` |
| Particionado de auditoría y purga (`BAC-18B`) | ✅ particionado trimestral y purga horaria; **proyectado**: índice con `jti_access`, apagado ordenado e historial legacy accesible (`FIX-44`, `FIX-50`, `FIX-51`) | `terminado.md` + `actual.md` |
| `LOGIN-04` | ✅ 10/10 front↔back | `terminado.md` |

**Hitos de cierre después del PDF** (de `terminado.md`): 25/09 SEC-01/02, FIX-17, FIX-25, SEC-04 · 28/09 FRN-11/12, FIX-18/19/20/21, INF-05 · 29/09 FIX-22, FIX-23, FRN-14 · 30/09 BAC-17A/B, BAC-21C, INF-06A, INF-08B · 01/10 BAC-16B, FIX-27/28/30/31, LOGIN-04 · 02/10 FIX-29/38/40/41/42, SEC-03 · 03/10 FIX-36/37/39, BAC-18B.

**Todavía abierto de la base** (no es de `actual.md`): `INF-06B`, `INF-08A` (servidor) y `FIX-54` a `FIX-57` (reutilización de correo tras eliminar un usuario; revisión estática, sin ejecutar).

---

## 3. Etapa 1: tablero

**Leyenda:** ✅ verificada · 🟡 verificada con ajuste pendiente (implementada) · 🔮 proyectada (`actual.md`) · ⏳ planificada (olas 2 a 5)

### 3.1 Antes de la etapa y transversales (no cuentan en las 29 tareas)

| Pieza | Estado |
|---|---|
| Simulador de Proxmox (`BAC-28`, `FIX-32`, `FIX-33`, `FIX-35`) | ✅ DELETE, IP VM/LXC, `nextid`, tareas, `409 INSTANCE_BUSY`, formatos idénticos al Proxmox real |
| `FIX-39` (campos de `GET /api/instances`, `shutdown`/`reboot`, auditoría de la orden despachada) | ✅ |
| `FIX-37` (`permisos` en `/account/profile`) | ✅ |
| `FIX-40` (`504 PROXMOX_TIMEOUT` vs `502 PROXMOX_UNAVAILABLE`) | ✅ |

### 3.2 Las 29 tareas

| Ola | Tarea | h | Estado | Qué hay |
|---|---|---:|---|---|
| 1 | `INF-07` token, red y VMIDs protegidos | 1,0 | 🟡 | Token verificado (privilegios `Sys.Audit`, `VM.PowerMgmt`, etc.). Pendiente de servidor: rotar el secreto |
| 1 | `INF-07B` Nginx para SSE | 1,0 | ✅ | Referencia local; evento en ~1 s sin buffer. Aplicación en CT 103 se valida en `INT-03` |
| 1 | `BAC-29` contrato de la etapa | 1,0 | 🟡 | `contrato-etapa1.md`, versión `.docx` y Swagger de `/node/status`. Ajuste: `FIX-46` (`activeTask`, `DELETE`, `exitstatus`) |
| 1 | `FRN-19A` medidores del host | 2,0 | ✅ | `ResourceMeter.tsx`, umbral 70 % (D3), 3/3 |
| 1 | `FRN-20A` tabla de inventario | 2,5 | 🟡 | Integrada a `GET /instances` con Bearer y permisos. Ajuste: `FIX-43` (maquetado, IP, badges) |
| 1 | `FRN-17A` consumo de eventos | 1,5 | 🟡 | Una sola conexión (`EventsProvider`), descarte de mensajes inválidos. Ajuste: `FIX-47` (deduplicar) |
| 2 | `BAC-24A` estado previo y protegidos | 1,5 | 🟡 | **Todo el comportamiento cumple**: 409/403/400/202 en ruta genérica y alias, sin tráfico a Proxmox al rechazar. Ajuste: Swagger (`FIX-48`) |
| 2 | `BAC-24B` `DELETE /api/instances/:vmid` | 2,0 | 🟡 | Existe: solo `ADMIN`, protegidos con 403, sin llegar a Proxmox al rechazar. Ajuste: `FIX-49` (202 con seguimiento, `INSTANCE_INVALID_STATE`) |
| 1 | `BAC-22` telemetría con Redis | 2,5 | 🔮 | `GET /api/node/status`: normalización, caché TTL 5-10 s, último estado con `stale`, D1 |
| 1 | `BAC-23A` adaptador de inventario e IP | 2,5 | 🔮 | QEMU/LXC unificados, IP por guest agent/interfaces LXC, timeout ≤ 2 s |
| 1 | `BAC-25A` pool acotado de UPID | 2,0 | 🔮 | `UPID_WORKERS` (8 por defecto), sin tareas perdidas |
| 2 | `FRN-15` modales antierror | 2,5 | 🔮 | Start, Shutdown/Reboot, Stop en rojo, Delete con tipeo (solo ADMIN) |
| 2 | `FRN-19B` semáforo del nodo | 2,0 | 🔮 | Saludable / Advertencia / Inaccesible, polling 10 s, `stale` |
| 3 | `FRN-16` operación en progreso | 2,5 | 🔮 | Spinner, fila bloqueada, seis mensajes de D2 |
| 3 | `BAC-25C` reanudación de UPID | 1,5 | 🔮 | Retoma `RUNNING` al reiniciar, vence a los 3 min (D4), `activeTask` |
| 2 | `BAC-22B`, `BAC-23B`, `BAC-25B`, `FRN-20B` | 6,5 | ⏳ | Resumen por estado, IP/CPU/RAM en inventario, timeout y `exitstatus`, filtros |
| 3 | `BAC-24C`, `BAC-26`, `BAC-27`, `FRN-19C`, `FRN-20C` | 6,5 | ⏳ | Limpieza de permisos, `TASK_FINISHED` de todas las acciones, auditoría del resultado, tarjetas y columnas |
| 4 | `FRN-16B`, `FRN-17B`, `INT-01` | 6,0 | ⏳ | Resincronización tras F5, toasts, suite de integración |
| 5 | `INT-02`, `INT-03` | 3,5 | ⏳ | Smoke test local y validación en el servidor real |

### 3.3 Avance en cifras

| Corte | Tareas | Horas | % de 50,5 h |
|---|---:|---:|---:|
| Verificado (`terminado*.md`) | 8 de 29 | 12,5 | ≈ 25 % |
| Proyectado (`actual.md` implementado) | +7 = 15 de 29 | +15,5 = 28,0 | ≈ 55 % |
| Restante (olas 2 a 5 no cargadas) | 14 | 22,5 | ≈ 45 % |

*Las horas son las estimadas en `etapa1.md`. El porcentaje es mío.*

### 3.4 FIX de `actual.md` incluidos en el escenario proyectado

| Grupo | FIX | Qué cierra |
|---|---|---|
| Contrato y eventos | `FIX-45`, `FIX-46`, `FIX-47` | Ruta `/api/events` del stream, `activeTask` unificado, deduplicación por `id` |
| Energía y borrado | `FIX-48`, `FIX-49`, `FIX-52`, `FIX-53` | Swagger, `DELETE` con 202 y seguimiento, `INSTANCE_BUSY` vs `INSTANCE_INVALID_STATE`, auditoría con códigos unificados |
| Inventario | `FIX-43` | Nombre/VMID, IP con copia, badges de estado |
| Base | `FIX-44`, `FIX-50`, `FIX-51` | Índice con `jti_access`, apagado ordenado del servidor, historial legacy de auditoría |

Suman ≈ 10 h de ajustes, la mayoría de 0,5 a 1,5 h cada uno.

---

## 4. Estado de los requerimientos (`requerimientos.md`)

| Requerimiento | Verificado hoy | Con `actual.md` implementado |
|---|---|---|
| RF-01 Autenticación, 2FA, roles | ✅ | ✅ |
| RF-09 Usuarios y permisos | ✅ (brecha de correo reutilizable: `FIX-54..57`) | igual |
| RF-13 Recuperación de contraseña | ✅ | ✅ |
| RF-08 Auditoría | 🟡 append-only, particionada, consulta y export CSV | ✅ historial legacy visible; falta el resultado de acciones de Proxmox (`BAC-27`) |
| RF-11 Notificaciones | 🟡 contrato, canal SSE con ticket, Pub/Sub, cliente | ✅ cliente y ruta alineados, sin duplicados; falta `TASK_FINISHED` de todas las acciones y toasts (`BAC-26`, `FRN-17B`) |
| RF-02 Dashboard del nodo | 🟡 medidores | ✅ datos reales, semáforo D3, `stale`; falta resumen por estado (`BAC-22B`, `FRN-19C`) |
| RF-03 Inventario | 🟡 listado con RBAC y tabla integrada | ✅ tabla corregida; IP/CPU/RAM reales llegan con `BAC-23B`/`BAC-22B` |
| RF-04 Ciclo de vida con UPID | 🟡 órdenes 202, validación, protegidos, `DELETE`, seguimiento básico | ✅ pool acotado, reanudación, modales, estado en progreso; falta timeout/reintentos (`BAC-25B`) y resincronización (`FRN-16B`) |
| RF-05, 06, 07, 10 | ⏳ Etapas 2 y 3 | ⏳ |
| RF-12, RNF-07 | ⛔ fuera del MVP | ⛔ |
| RNF-01 Seguridad | 🟡 | 🟡 ver §6 |
| RNF-04 Confiabilidad | 🟡 errores de Proxmox diferenciados | ✅ con `BAC-25A/C`; completo con `BAC-25B` |
| RNF-05 Mantenibilidad | 🟡 | ✅ Swagger alineado (`FIX-46`, `FIX-48`, `FIX-49`) |

---

## 5. Cómo evolucionó el método

- **Suite de aceptación independiente** (`test/back`, `test/front`): al 05/10, backend 58 casos (49 en verde) y frontend 139 (95 en verde), con todos los fallos atribuidos al producto. Siete pruebas defectuosas se reescribieron y las del frontend se validaron con contrapruebas (34/34 con implementación correcta; 13 defectos inyectados detectados uno a uno).
- **Contrato primero** (`BAC-29`) y simulador de Proxmox completo: el frontend desarrolla sin esperar al backend ni al servidor.
- **Local-First**: infraestructura de servidor (`INF-06B`, `INF-08A`, `INF-07`) replica lo ya validado en local y no frena el desarrollo.
- **Decisiones de producto D1 a D4** registradas (acceso a `/node/status`, un código por resultado, umbral 70 %, timeout de 3 min).
- **Estado explícito "implementada con problema"** con su FIX, en lugar de ocultar tareas a medias.

---

## 6. Riesgos y observaciones (siguen vigentes)

1. **Secretos en el repositorio:** `test/back/smtp-brevo.env` tiene credenciales SMTP reales y `documentacion/api-proxmox.md` el secreto del token de Proxmox en texto plano. Conviene rotarlos (`INF-07` ya lo pide para el token).
2. **Servidor sin validar:** `INF-06B`, `INF-08A`, la aplicación de `nginx-edge.conf` en el CT 103 y la rotación del token son requisitos de `INT-03`.
3. **Baja de usuarios** (`FIX-54` a `FIX-57`): propuesta de revisión estática, sin implementar ni probar; toca el esquema y conviene decidir su alcance pronto.
4. **Carga concentrada:** Lisandro (≈ 12 h de etapa + FIX de energía) y Tayra (≈ 9,5 h + `INT-01` + FIX de base).
5. **Inconsistencias documentales:** `siguientesetapas.md` dice SMTP real descartado; `requerimientos.md` conserva Kafka y `RF-12`; `centinela.md` habla de un wizard de 3 pasos y la Etapa 2 usa 4; `QA-11` figura en el PDF y no en estos documentos.
6. **Lo proyectado hay que confirmarlo:** el 55 % depende de que cada tarea de `actual.md` cumpla su criterio de éxito en la suite de aceptación.

---

## 7. Próximos pasos

1. Ejecutar y cerrar `actual.md` (≈ 25,5 h de tareas + ≈ 10 h de FIX), empezando por los FIX de contrato (`FIX-45`, `FIX-46`, `FIX-49`), que destraban `FRN-16` y `FRN-17B`.
2. Cargar la **Ola 2 y 3** en ClickUp: `BAC-22B`, `BAC-23B`, `BAC-25B`, `FRN-20B`, y luego `BAC-24C`, `BAC-26`, `BAC-27`, `FRN-19C`, `FRN-20C` (necesita `BAC-18B` ya mergeado).
3. Coordinar con infraestructura `INF-06B`, `INF-08A` e `INF-07` antes de `INT-03`.
4. Mantener la regla: una tarea pasa a `terminado*.md` con su criterio verde; si queda un ajuste, se registra su FIX.

---

## 8. Conclusión

Entre el cierre del PDF y hoy el sistema pasó de "base lista para enchufar Proxmox" a **un núcleo operativo funcionando sobre un simulador fiel**: contrato publicado, energía y borrado con reglas de acceso y de estado, canal de eventos en tiempo real, tabla de inventario integrada y medidores del host, además de una base de autenticación completa con `LOGIN-04` en verde. Con `actual.md` cumplido, la Etapa 1 llega a **más de la mitad (≈ 55 % de las horas)** con telemetría real, seguimiento robusto de tareas y la interfaz operativa de energía. La segunda mitad se concentra en completar datos (IP, CPU, RAM), cerrar el ciclo de auditoría y notificaciones, y validar todo en el servidor real.
