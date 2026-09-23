Sí, quedan huecos. Crucé `etapa1.md` con lo que realmente quedó implementado (según la evidencia de `terminado.md`, `actual.md` y `futuro.md`) y con `requerimientos.md`. Los agrupo de más a menos grave.

## 1. La Etapa 1 está escrita sobre nombres que no coinciden con lo implementado

Es lo más peligroso, porque varias tareas tienen criterios de éxito que dependen de reutilizar cosas "sin migrar", y lo que van a encontrar en el código es distinto.

- **La tabla de auditoría no se llama `audit_logs`.** La evidencia de `FIX-23` muestra que la tabla real es `auditoria`. `BAC-27` exige "sin migraciones de esquema adicionales" y nombra columnas en inglés (`user_id`, `resource_type`, `upid`). Alguien tiene que confirmar qué columnas existen de verdad antes de que Tayra arranque; si no, el criterio de `BAC-27` es imposible.
- **La tabla de permisos no es `user_instances`.** Según `SEC-04`, es `permisos_instancia (usuario_id, vmid_proxmox)`. `BAC-23B` y `BAC-24A` hablan de `user_instances`.
- **`GET /api/instances` ya existe** (`BAC-14`), con la forma `{ id, name, type, node, status }`, y lo consume el selector de `FRN-07`/`FIX-14`. `BAC-23B` lo plantea como si se creara de cero. Tiene que ser una extensión compatible (agregar IP y métricas), no un reemplazo, o se rompe el selector de permisos.
- **Ya existen rutas de energía.** `FIX-16` montó `POST /instances/:vmid/start` y `/stop` con el guard, y las pruebas de `resource_access_acceptance_test.go` las usan. `BAC-24A` define otra forma: `POST /api/instances/{id}/status/{action}`. Nadie decide si las viejas se eliminan, se redirigen o conviven, ni quién actualiza esas pruebas.
- **Hay IDs repetidos.** `FRN-13` y `FRN-14` ya existen en la base (logout y suite de auditoría), y la Etapa 1 usa `FRN-13A/B` y `FRN-14A/B` para cosas totalmente distintas. En ClickUp va a generar confusión; conviene renumerarlas (por ejemplo, a partir de `FRN-18`).

## 2. Trabajo de la fase base que la Etapa 1 no recoge

- **`SEC-04` (niveles `READ_ONLY` / `FULL_ACCESS`) no aparece en ninguna tarea de la Etapa 1.** `BAC-24A` y `BAC-24B` solo validan que la instancia esté asignada, no el nivel. Así, un operador de solo lectura podría apagar máquinas, que es justo lo que `SEC-04` quiere impedir. Además, `futuro.md` dice en "Decisiones de alcance" que los permisos granulares por acción están fuera de alcance, y `SEC-04` es exactamente eso. Hay que resolver esa contradicción y, si `SEC-04` sigue, agregar el nivel requerido en `BAC-24A`/`24B`, en el filtro de eventos de `BAC-26` y en la botonera de `FRN-15`/`FRN-16`.
- **`SEC-03` (`usePermissions`, `canAccessInstance`, `PermissionGate`) no figura como dependencia** de ninguna tarea de frontend de la Etapa 1. Sin embargo, `FRN-14A`, `FRN-15` y `FRN-16` necesitan saber qué botones mostrar según el rol y el nivel de acceso.
- **`FIX-23` (trigger append-only) tendría que ser un requisito explícito de `BAC-27`.** Si no, la Etapa 1 empieza a auditar operaciones sobre una tabla que todavía permite `UPDATE` y `DELETE`.
- **`INF-05` (TLS/CORS) no está en la "dependencia previa" de la Etapa 1**, que solo menciona `LOGIN-04`. Sin HTTPS, la cookie `Secure` de `SEC-01` no funciona, y un WebSocket sin TLS (`ws://` en lugar de `wss://`) desde una página HTTPS lo bloquea el navegador.

## 3. Cosas que no cubre ninguna de las dos fases

**Autenticación del canal en tiempo real.** Es el hueco técnico más grande. `BAC-26` y `FRN-17A` no definen cómo se autentica `/api/events`:
- `EventSource` (SSE) no permite enviar el header `Authorization`, y en los WebSockets tampoco es práctico.
- La cookie de `SEC-01` es la del refresh token, y lo normal es que esté restringida a `/api/auth`.

Hay que decidir el mecanismo; por ejemplo, un ticket de corta duración emitido por un endpoint autenticado. También hay que definir qué pasa con una conexión abierta cuando la sesión se revoca (`BAC-17`) o cuando el admin le quita una instancia al operador. Hoy nada cierra esa conexión ni recalcula el filtro.

**Nginx para WebSocket/SSE.** `INF-04` solo enruta `/api/*`. Faltan los headers `Upgrade`/`Connection` para WebSocket, o `proxy_buffering off` y timeouts largos para SSE. `RNF-06` lo pide explícitamente y no hay tarea de infraestructura para eso en la Etapa 1.

**Permisos del token de Proxmox.** `INF-07` pide `Sys.Audit`, `VM.Audit` y `VM.PowerMgmt`, pero eso no alcanza para todo:
- Eliminar una VM o un LXC (`BAC-24B`) requiere `VM.Allocate`; con esos tres permisos, Proxmox va a rechazar el borrado.
- Leer la IP por QEMU Guest Agent (`BAC-23A`) necesita un privilegio de monitoreo del agente, cuyo nombre cambió entre versiones de Proxmox.

Conviene que Nico lo confirme contra la versión instalada.

**Quién puede eliminar instancias.** `RF-04` no incluye el borrado; lo agrega la Etapa 1. Pero `BAC-24B` no dice si es solo para `ADMIN` o si un operador puede borrar sus propias máquinas. Es la acción más destructiva del sistema y no tiene una regla de autorización definida.

**Qué pasa si se pierde un evento.** El worker de UPID vive en memoria, y hay tres escenarios sin cubrir:
- Si el backend se reinicia con tareas en curso, esas tareas quedan huérfanas, la auditoría queda en `PENDING` para siempre y la fila de la interfaz queda bloqueada con el spinner.
- Si el usuario recarga la página o pierde la conexión justo cuando termina la tarea, `FRN-17B` nunca recibe el `TASK_FINISHED`.

Falta un mecanismo de recuperación: un `GET` del estado de la tarea o del estado de la instancia al reconectar, y la reanudación de UPIDs pendientes al arrancar. `RNF-04` lo exige.

**Contradicción con la escalabilidad.** La decisión 4 de la Etapa 1 justifica Redis "para escalabilidad horizontal", pero `BAC-25B` usa un bus interno en memoria (un canal de Go). Con dos réplicas del backend, el evento solo llega a los clientes conectados a la réplica que hizo el sondeo. Si ya van a tener Redis, lo lógico es usar su pub/sub como bus, o bien declarar explícitamente que el MVP funciona con una sola réplica.

**Requerimientos parcialmente cubiertos:**
- `RF-02` pide la cantidad de VMs y LXC agrupadas por estado en el dashboard, y ni `BAC-22` ni `FRN-13B` lo incluyen. También pide que el dashboard se actualice sin recargar, y `FRN-13B` no define ni polling ni suscripción.
- `RF-03` pide un resumen de CPU y RAM por instancia, y no está en `BAC-23A` ni en las columnas de `FRN-14A`. También pide sincronizar el inventario con cambios de estado: con el diseño actual solo se actualiza lo que dispara el propio usuario, no lo que otro usuario (o alguien desde Proxmox) haga.
- `RF-11` pide que cada notificación enlace al detalle de la instancia, y en la Etapa 1 no existe esa vista. Se puede posponer, pero conviene dejarlo escrito.

**Contrato de errores.** Van a aparecer códigos nuevos (409 por estado inválido, timeout de tarea, Proxmox inaccesible), y ninguna tarea los agrega al inventario de `errorCode` de `FIX-08` ni a Swagger (`RNF-05`).

## 4. El criterio de cierre de la fase base es débil

La Etapa 1 exige "`LOGIN-04` aprobado", pero `LOGIN-04` no depende de `SEC-01`, `SEC-03`, `SEC-04`, `INF-05` ni de los `FIX-20` a `FIX-25`. Se podría aprobar `LOGIN-04` y pasar de fase con el logout real roto (`SEC-01`), la auditoría modificable (`FIX-23`) y un operador que ve la pantalla de auditoría (`FIX-22`).

Sugiero definir la salida de la fase base como `LOGIN-04` más todas las tareas de `actual.md` y los `FIX` de `futuro.md` cerrados, o bien declarar explícitamente cuáles se aceptan como deuda.

---

Me expliqué mal. La base **no** está mal hecha: casi todo lo que marqué son ajustes al documento de la Etapa 1, no tareas pendientes de la fase base.

## Sobre los permisos granulares

Sí se están definiendo: eso es exactamente `SEC-04`, que está en `actual.md` (niveles `READ_ONLY` y `FULL_ACCESS`). Lo que quise señalar son dos cosas:

1. En `futuro.md`, la sección "Decisiones de alcance" todavía dice que los permisos granulares están fuera de alcance. Es un texto viejo que quedó desactualizado cuando se agregó `SEC-04`, y hay que corregirlo.
2. El documento de la Etapa 1 se escribió como si `SEC-04` no existiera, así que sus endpoints de energía no usan el nivel de acceso.

## Qué le falta a la fase base

Una sola tarea nueva y dos correcciones de documentación.

**Tarea nueva: conectar el nivel de acceso en el frontend.** `SEC-04` es solo backend, y su criterio acepta que el frontend siga mandando `{ vmids: [101] }`, que se guarda como `FULL_ACCESS`. El selector ya muestra "Solo lectura", pero nadie tiene asignado:
- mandar el nivel en `PUT /permissions`;
- leerlo de vuelta al abrir el usuario;
- sumar al contexto de `SEC-03` algo como `canOperateInstance(vmid)`, que responda "puede operar", distinto de "puede ver".

Sin esto, cuando termine `SEC-04`, el administrador no va a tener forma de asignar "solo lectura" desde la interfaz.

**Correcciones de documentación:**
- Actualizar la "Decisión de alcance" de `futuro.md` sobre los permisos granulares.
- Definir qué cierra la fase base. Hoy es solo `LOGIN-04`, que no incluye `SEC-01`, `SEC-03`, `SEC-04`, `INF-05` ni los `FIX-20` a `FIX-25`. Lo razonable es que la base cierre cuando esté todo eso más `LOGIN-04`.

## Qué le falta a la Etapa 1

Todo lo demás del mensaje anterior pertenece a la Etapa 1, porque son cosas que aparecen recién cuando existen el canal de eventos, el worker de UPID y las acciones de energía.

**Tareas nuevas:**
- **Autenticación del canal `/api/events`:** cómo se autentica la conexión y qué pasa si se revoca la sesión o se quita una instancia mientras está abierta.
- **Nginx para WebSocket/SSE:** headers de upgrade, buffering y timeouts.
- **Resiliencia del worker:** retomar los UPIDs pendientes cuando el backend se reinicia, y que el frontend se resincronice al reconectar o al recargar la página.

**Ajustes a tareas que ya están escritas:**
- **Alineación de nombres:** usar `auditoria` y `permisos_instancia`, extender el `GET /api/instances` existente sin romper el selector, decidir qué pasa con las rutas `start`/`stop` actuales y renumerar `FRN-13A` a `FRN-14B`.
- **`INF-07`:** agregar el permiso de borrado (`VM.Allocate`) y el de lectura de IP por guest agent.
- **`BAC-24A` y `BAC-24B`:** exigir `FULL_ACCESS` y definir si eliminar instancias es solo para `ADMIN`.
- **`FRN-15` y `FRN-16`:** depender de `SEC-03` para mostrar u ocultar botones.
- **`BAC-27`:** depender de `FIX-23`.
- **`BAC-22` y `FRN-13B`:** agregar el conteo de instancias por estado y definir cómo se actualiza el dashboard.
- **`BAC-23A` y `FRN-14A`:** agregar el resumen de CPU y RAM por instancia, que pide `RF-03`.
- **`BAC-25B`:** usar Redis pub/sub como bus, o dejar escrito que el MVP corre con una sola réplica.

---

En resumen, a la base le falta una tarea chica de frontend más dos correcciones de texto; el resto es mejorar el documento de la Etapa 1 antes de cargarlo en ClickUp.

¿Querés que te redacte la tarea nueva de la base y las tres de la Etapa 1 en el formato que usan en los `.md`?

Para arrancar la Etapa 1 sin sorpresas, lo mínimo que haría es:
1. Una tarea corta de alineación de contratos: nombres reales de tablas y columnas, rutas de energía nuevas contra las viejas, y renumeración de IDs.
2. Una tarea de backend para la autenticación del canal de eventos, más una de infraestructura para Nginx con WebSocket/SSE.
3. Ampliar `INF-07` con los permisos de borrado y de lectura de IP.
4. Una tarea de resiliencia: reanudar UPIDs al arrancar y resincronizar el frontend al reconectar.

Si querés, te redacto esas tareas en el mismo formato de ClickUp que usan en los `.md`.
