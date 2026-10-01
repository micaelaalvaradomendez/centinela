
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 01/10/2026** (detalle en [test/informe.md](../test/informe.md)).
> - `BAC-16B`, `FIX-27` y `FIX-28` están completas, y `FIX-30` está implementada con problema (su corrección es `FIX-37`, en el backend). Las cuatro pasaron a [`terminado.md`](terminado.md).
> - FIX nuevos en [`futuro.md`](futuro.md): `FIX-36` (regresión del filtro "Acción" en Auditoría, `FRN-14`), `FIX-37` (`permisos` con nivel en `GET /account/profile`) y `FIX-38` (lo que ya implementó `SEC-03`).
>
> **Quedan en este archivo:**
>
> | Tarea | Área | Estado |
> |---|---|---|
> | `SEC-03` | Frontend | En proceso: `usePermissions` existe en `src/hooks/`; faltan `isOperator`, `hasRole`, `PermissionGate` y el menú (ver `FIX-38`) |
> | `FIX-29` | Frontend | En proceso, sin cambios: mayúscula en el validador |
> | `FRN-17C` | Frontend | En proceso, sin cambios: `useEvents` con ticket efímero (desbloqueada por `BAC-21C`) |
> | `BAC-18B` | Backend | No implementada: índice parcial en `sesiones_activas` (`jti_access`) y particionado de `auditoria` |
> | `BAC-21B` | Backend | No implementada: campos nuevos en `GET /instances`, `/status/:action` y `DELETE` |
> | `FIX-31` | Infraestructura | No implementada: configuración de Nginx con TLS versionada |
> | `INF-06B`, `INF-08A` | Infraestructura | Redis y SMTP en el servidor (sin prueba automatizada) |
> | `FIX-35` (Etapa 1) | Backend | No implementada: el simulador no cambió desde `43a0b06` |

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`).

---

### `SEC-03` - Contexto y sistema reactivo de permisos en Frontend (UI/UX RBAC + Resources)

- **Área:** Frontend
- **Asignados:** Cristian y Belinda
- **Estimación:** 2,5 h
- **Ventana propuesta:** A definir (Fase Base / Bloque 2).
- **Depende de:** `FRN-04`, `BAC-07` y `BAC-09`.
- **Estado verificado (01/10/2026, frontend `3d1e84a`):** en proceso, parcial. `test/front/navigation.test.tsx`: fallan los 3 casos de SEC-03 (el de `FIX-30` pasa).
  - Existe `usePermissions()` en `src/hooks/usePermissions.ts` (no en `src/context/`), con `isAdmin`, `canAccessInstance` y `canOperateInstance`.
  - Faltan `isOperator`, `hasRole` y `PermissionGate`; `src/context/AuthContext.js` sigue vacío.
  - `Sidebar.tsx` sigue mostrando "Auditoría" al OPERATOR.
  - Problemas de lo ya implementado: ver `FIX-38` en [`futuro.md`](futuro.md).
- **Estado anterior (23/09/2026):** no implementada.
  - `src/context/AuthContext.js` está vacío.
  - No existen `usePermissions` ni `PermissionGate`.
  - El menú solo filtra con un `isAdmin` ad hoc en `Sidebar.tsx` y no tiene un acceso a Auditoría para el administrador.
  - Prueba: `test/front/navigation.test.tsx` (1/4).
- **Problema:** En el frontend actual, el control de roles y permisos está disperso y acoplado únicamente a los loaders de rutas (`loadAdminSession`). No existe un mecanismo reactivo a nivel de componentes (`usePermissions` / `PermissionGate`) para ocultar o deshabilitar condicionalmente acciones según el rol (`ADMIN` vs `OPERATOR`) o según las instancias asignadas al operador. Esto genera que en vistas como `UserDetail.tsx` o `Header.tsx` se muestren controles estáticos o se dependa de que el backend rechace con 403, en lugar de ofrecer una experiencia fluida y consistente.
- **Entregable:**
  1. Crear un hook y contexto `usePermissions()` / `useAuthUser()` en `frontend/centinela/src/context/` que exponga helpers como `isAdmin`, `isOperator`, `canAccessInstance(vmid: number)` y `hasRole(role: string)`.
  2. Crear un componente wrapper `<PermissionGate requiredRole="ADMIN" fallback={...}>` para condicionar la renderización de botones y secciones administrativas (ej: botón "Eliminar usuario", accesos a auditoría, etc.).
  3. Integrar la reactividad en el menú de navegación y en el header para que las opciones no permitidas a un operador no aparezcan en la interfaz.
- **Criterio de éxito:** Si un operador inicia sesión, la interfaz no muestra accesos directos ni botones exclusivos de administrador; si intenta interactuar con un componente restringido, el helper `canAccessInstance` evalúa en memoria las instancias permitidas cargadas en la sesión; las pruebas unitarias de componentes verifican el render condicional.


---


### `BAC-18B` - Índice parcial y purga en `sesiones_activas`, y particionamiento trimestral en `auditoria` (PostgreSQL)

- **Área:** Backend
- **Asignada:** Tayra
- **Estimación:** 2 h
- **Ventana propuesta:** A definir (Cierre de Fase Base).
- **Depende de:** `BAC-17B`, `FIX-23`.
- **Problema y contexto:**
  1. Las filas históricas o con `activa = false` en `sesiones_activas` penalizan las lecturas en PostgreSQL si no existe un índice parcial ni una purga de registros vencidos.
  2. La tabla `auditoria` es *append-only* (`FIX-23`, no admite `DELETE` bajo ningún concepto) y en la Etapa 1 registrará múltiples eventos por cada operación de Proxmox VE (`PENDING` y `SUCCESS`/`FAILED`). Sin particionamiento por fechas e índices compuestos, las consultas de `GET /api/admin/audit` y la exportación CSV se degradarán progresivamente.
- **Entregable:**
  1. Crear en PostgreSQL (`init.sql` / migración) el índice parcial para `sesiones_activas`:
     ```sql
     CREATE INDEX IF NOT EXISTS idx_sesiones_activas_vigentes
       ON sesiones_activas (jti_token, usuario_id)
       WHERE activa = true;
     ```
  2. Agregar una rutina de limpieza en el backend (ticker cada 1 hora) que ejecute `DELETE FROM sesiones_activas WHERE activa = false OR fecha_expiracion < NOW();`.
  3. Configurar en PostgreSQL el **particionamiento declarativo trimestral por rango de fechas** sobre la tabla `auditoria` (`PARTITION BY RANGE (fecha_hora)`), creando las particiones trimestrales (`auditoria_2026_q3`, `auditoria_2026_q4`, `auditoria_2027_q1` y `auditoria_default`) junto con los índices compuestos `(fecha_hora DESC, accion, resultado)` y `(usuario_id, fecha_hora DESC)`.
- **Criterio de éxito:** Las sesiones muertas se purgan automáticamente de PostgreSQL; la tabla `auditoria` opera sobre particiones trimestrales manteniendo tiempos de consulta constantes (< 20 ms) e inmutabilidad append-only.


---


### `FIX-29` - Validar la mayúscula en `validatePasswordComplexity` (`FIX-20`) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** `FIX-20`.
- **Problema y evidencia:** en `frontend/centinela/src/components/features/auth/utils/validateAuthenticationFields.ts`, `validatePasswordComplexity` tiene tres condiciones. La segunda prueba `/[0-9]/` pero agrega el mensaje *"Debe contener al menos una letra mayúscula."*, y no hay ninguna condición que pruebe mayúsculas. El dígito se valida una sola vez y la mayúscula nunca, así que el cliente acepta claves como `nueva1234!` o `sinmayus1!`, que el backend rechaza con `400 PASSWORD_CHANGE_FAILED` (`crypto.ValidarComplejidadContrasena`). Además, el mensaje de largo dice *"Debe tener 8 y 12 caracteres."*.
  - Pruebas que fallan:
    - `test/front/password-change.test.tsx`: *"FIX-29 no llama a la API si la contraseña nueva no cumple la complejidad del backend (sin mayúscula)"*.
    - `test/front/recover-password.test.tsx`: *"FIX-29 el paso 3 no envía una contraseña que no cumple la complejidad del backend"*.
- **Entregable:**
  1. Separar las condiciones: `/[A-Z]/` con el mensaje de mayúscula y `/[0-9]/` con el mensaje *"Debe contener al menos un número."*.
  2. Corregir el mensaje de largo: *"Debe tener entre 8 y 12 caracteres."*.
- **Criterio de éxito:** los casos `FIX-29…` pasan, y siguen en verde los de dígito, símbolo y largo, y los casos de `FRN-12` y `FIX-21`.


### `BAC-21B` (`BRG-01`) - Alineación de esquema (`auditoria`/`tareas_asincronas`), rutas de energía (`FULL_ACCESS`), regla de `DELETE` y extensión de `GET /api/instances`

- **Área:** Backend
- **Asignados:** Tayra y Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** Previo al inicio de `BAC-23B`, `BAC-24A/B` y `BAC-27`.
- **Depende de:** `BAC-14`, `FIX-16`, `SEC-04`, `FIX-23`.
- **Problema y evidencia (análisis de cierre de la fase base):**
  1. `etapa1.md` menciona `audit_logs` y `user_instances`, cuando las tablas reales en `domain/models.go` son `auditoria`, `tareas_asincronas` y `permisos_instancia`.
  2. `GET /api/instances` ya existe (`BAC-14`) y lo consume el selector `FRN-07`. Si se reemplaza su formato en vez de extenderlo, se rompe la pantalla de permisos.
  3. `FIX-16` creó `/api/instances/:vmid/start` y `/stop`, mientras que `BAC-24A` pide `/api/instances/:vmid/status/:action`. Además, falta exigir `FULL_ACCESS` en energía y restringir `DELETE /api/instances/:vmid` solo a `ADMIN`.
- **Entregable:**
  1. Persistir la auditoría de ciclo de vida en la tabla real `auditoria` (guardando `upid`, `action` y `resource_type` en `detalles` JSONB) y el estado de ejecución en `tareas_asincronas`, consultando permisos contra `permisos_instancia`.
  2. Extender `GET /api/instances` de forma 100% retrocompatible: mantener `{ id, name, type, node, status }` y agregar `{ ip, cpuUsage, ramUsage, maxRam, nivelAcceso, activeTask, instancesSummary }`.
  3. Montar `POST /api/instances/:vmid/status/:action` manteniendo alias en `/start` y `/stop`, protegidos con `RequireInstanceAccess(repo, "vmid", "FULL_ACCESS")`. Proteger `DELETE /api/instances/:vmid` (`BAC-24B`) con `RequireRole("ADMIN")` y validación de estado `stopped` (409 Conflict si está encendida).
- **Criterio de éxito:** `GET /api/instances` responde con los campos nuevos sin romper `FIX-14`; un operador `READ_ONLY` recibe 403 en acciones de energía; un `OPERATOR` recibe 403 en `DELETE`; todo se registra en `auditoria` y `tareas_asincronas`.

---
### `FRN-17C` (`BRG-02-FRN`) - Cliente de eventos con solicitud previa de ticket efímero y reconexión segura

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 1,5 h
- **Ventana propuesta:** Junto a `FRN-17A`.
- **Depende de:** `BAC-21C` (`BRG-02-BAC`).
- **Problema y contexto:** El hook `useEvents` del frontend no puede pasar el JWT por header en `EventSource`/WebSocket ni exponer el access token largo en la URL. Debe solicitar primero el ticket efímero al backend.
- **Entregable:**
  1. En `useEvents` (`FRN-17A`), antes de abrir la conexión hacia `/api/events`, invocar `POST /api/events/ticket` con el interceptor autenticado (`Bearer`) y conectar a `/api/events?ticket=<uuid>`.
  2. Ante una desconexión de red, solicitar un nuevo ticket efímero aplicando retroceso exponencial; si `/api/events/ticket` responde `401`, disparar el cierre de sesión local y redirigir a `/login`.
- **Criterio de éxito:** El frontend se conecta a `/api/events` usando tickets de un solo uso sin exponer el JWT en la URL y se reconecta pidiendo un ticket fresco.


### `FIX-31` - Versionar y verificar el TLS de Nginx en el borde (`INF-05`) (Infraestructura)

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1 h
- **Ventana propuesta:** A definir.
- **Depende de:** `INF-04` e `INF-05`.
- **Problema y evidencia:** el criterio de `INF-05` exige que *"el tráfico hacia el sistema se sirva únicamente sobre HTTPS"*. La configuración de Nginx del servidor (CT 103) no está versionada en ningún repositorio. El deploy (`frontend/.github/workflows/deploy-front-test.yml`) solo ejecuta `nginx -t` y `reload` sobre lo que ya existe en el contenedor, y el único `nginx.conf` versionado (`docker/nginx.conf`, el del escenario de integración) escucha solo en `:80`. Por eso no se puede verificar el TLS. Prueba omitida: `test/back/cierre_fase_base_acceptance_test.go`, caso *"FIX-31 INF-05 TLS en el borde con Nginx"*.
- **Entregable:**
  1. Versionar la configuración de Nginx del borde, por ejemplo en `frontend/centinela/deploy/nginx.conf` o en `docker/`, con:
     - `listen 443 ssl` y las rutas de `ssl_certificate` y `ssl_certificate_key` (sin incluir los certificados);
     - `listen 80` que responda `return 301 https://$host$request_uri`;
     - `proxy_pass` de `/api/` al backend.
  2. Hacer que el deploy use ese archivo versionado.
- **Criterio de éxito:** la configuración versionada sirve solo HTTPS y redirige HTTP → HTTPS. A partir de ahí, el caso `FIX-31…` deja de omitirse y verifica el archivo (443 ssl, redirección 301 y ningún `server` que sirva contenido por `:80`).


### `INF-06B` - Redis en el servidor (entornos de pruebas y estable)

- **Área:** Infraestructura
- **Asignados:** Nico y Lucas
- **Estimación:** 1 h
- **Ventana propuesta:** A definir. **No bloquea al backend**, que trabaja con `INF-06A`. Tiene que estar antes de desplegar `BAC-17B` y `BAC-21C` en el servidor.
- **Depende de:** `INF-03` e `INF-04`.
- **Entregable:**
  1. Desplegar Redis en el servidor, en el LXC del backend o en uno propio de `vmbr1`, con contraseña y sin exponerlo fuera de la red interna.
  2. Cargar `REDIS_ADDR` y `REDIS_PASSWORD` en el `.env` del backend de pruebas y del estable.
- **Criterio de éxito:** desde el LXC del backend, `redis-cli -h <host> -a <pass> PING` responde `PONG` y sin contraseña responde `NOAUTH`. Desde fuera de la red interna no se alcanza.

### `INF-08A` - Configuración del SMTP en el servidor

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir. **No bloquea al backend**, que trabaja con `INF-08B`.
- **Depende de:** `INF-03`, `INF-04` e `INF-08B`.
- **Entregable:**
  1. Cargar `EMAIL_PROVIDER` y `SMTP_*` en el `.env` del backend del servidor (pruebas y estable).
  2. Habilitar y verificar la salida de red desde el LXC del backend hacia el host SMTP, por el puerto 587 (STARTTLS) o 465 (TLS). Por ejemplo, con `openssl s_client -starttls smtp -connect <host>:587`.
- **Criterio de éxito:** desde el LXC del backend se establece la conexión TLS con el servidor SMTP, y el backend desplegado envía el correo de alta de usuario.

---
# ETAPA 1
---

> Las tareas de la Etapa 1 ya verificadas están en [`terminado-1.md`](terminado-1.md), y sus correcciones en [`futuro-1.md`](futuro-1.md).

### `FIX-35` - Formato de las IP de contenedores en el simulador de Proxmox (`BAC-28`) (Backend)

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 0,5 h
- **Depende de:** `BAC-28`.
- **Problema y evidencia:** se comparó `GET /nodes/proxmox/lxc/{vmid}/interfaces` del simulador con el del Proxmox real 9.2.2 (lecturas del 30/09/2026):
  1. **Tipo de dato:** `ip-addresses[].prefix` es **string** en el real (`"prefix":"24"`) y **número** en el simulador (`"prefix":24`). `BAC-23A` va a leer la IP de ahí, y un decodificador escrito contra el simulador falla contra el real, o al revés.
  2. **Contenedor apagado:** el real responde **`200 {"data":null}`** (verificado con la 104, apagada); el simulador responde `500 CT 201 not running`. Con el simulador, `BAC-23A` trataría como error algo que en el real es "sin IP".
- **Entregable:**
  1. Devolver `prefix` como string en `lxc/{vmid}/interfaces`. Revisar también `agent/network-get-interfaces` contra la documentación de Proxmox, porque no se pudo comparar: el servidor no tiene VMs.
  2. Para un LXC apagado, responder `200` con `{"data": null}`.
  3. Agregar ambos casos a `cmd/proxmox-simulador/simulador_test.go`.
- **Criterio de éxito:** la comparación contra el Proxmox real de `lxc/{vmid}/interfaces` no muestra diferencias de tipos, y un contenedor apagado responde igual en los dos.