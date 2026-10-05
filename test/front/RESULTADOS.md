# Resultados de las pruebas frontend

**Fecha:** 05/10/2026. **Frontend:** submódulo en `d47999d` (último commit de `main`; nuevos: `FIX-36`, `FRN-19A`, `FRN-17C` y `FRN-17A`). **Comando:** `pnpm vitest run --no-file-parallelism`. 2 corridas con el mismo resultado.

| Métrica | Valor |
|---|---:|
| Pruebas | 139 |
| Aprueban | 95 |
| Fallan | **34** |
| Omitidas | 10 (`login04-e2e.test.ts`: se ejecuta desde `test/back`, donde pasa 10/10) |

## Tareas de `actual.md`

| Tarea | Archivo | Fallan | Causa |
|---|---|---:|---|
| `FRN-17C` cliente con ticket | `events-client.test.tsx` | 1 de 4 | **El frontend abre el stream en `/api/events/stream?ticket=…`** (`EVENTS_STREAM_PATH = '/events/stream'`), pero el backend lo sirve en `GET /api/events`: contra el backend real respondería 404. El ticket con `Bearer`, la reconexión con ticket nuevo, el retroceso exponencial y el `401` cumplen |
| `FRN-17A` consumo de eventos | `events-client.test.tsx` | 2 de 4 | La conexión única (`EventsProvider` + `useEventsContext`) y el descarte de mensajes inválidos cumplen. **No deduplica:** el cliente descarta el `id` del evento y un evento repetido se entrega dos veces. Además sigue existiendo `hooks/useWebSocket.js` |
| `FIX-43` tabla de instancias | `instances-table.test.tsx` | 5 de 5 | Sin cambios: nombre e ID en el mismo texto, IP fija en `—`, estado sin badge |
| `FRN-15` modales | `instances-modals.test.tsx` (**reescrito**) | 6 de 6 | La tabla no tiene botonera ni modales |
| `FRN-16` operación en progreso | `instances-operation.test.tsx` (**reescrito**) | 9 de 9 | Depende de `FRN-15`: no hay acciones que disparar |
| `FRN-19B` semáforo del nodo | `dashboard-node.test.tsx` (**reescrito**) | 11 de 11 | El Dashboard no consulta `GET /api/node/status` |

## Regresión

Pasan `FRN-19A` (`dashboard-metrics`), `FIX-36` (`audit`), `SEC-03`, `FIX-38`, `FIX-41`, `FIX-42`, `FIX-29` y el resto de las tareas terminadas.

## Cambios en la suite (05/10/2026)

- **`app-providers.tsx`:**
  - `withAppProviders` monta lo mismo que `App.tsx` (`AuthProvider`). Se usa en las pruebas que renderizan rutas, porque `ProtectedLayout` ya monta su `EventsProvider`.
  - `withProtectedProviders` (nuevo) agrega el `EventsProvider`. Se usa en las pruebas que renderizan una página protegida suelta.
  - Incluir el `EventsProvider` en todas las pruebas rompía `password-change`: el mock de `fetch` reutiliza la misma `Response`, y el pedido de ticket consumía su cuerpo.
- **`instances-helpers.tsx` (nuevo):** inventario de prueba y búsqueda de acciones de la fila. Acepta botones directos o un menú "Acciones", y modales con `role="dialog"` o `"alertdialog"`. Solo cuenta como órdenes las peticiones a `/instances/`.
- **FRN-15 (reescrito):**
  - el caso de `READ_ONLY` tenía un `if (botón)` que lo dejaba pasar sin verificar;
  - de shutdown y reboot probaba solo uno;
  - no verificaba que cancelar no envíe nada en todas las acciones, ni el rojo de Stop, ni que un texto incorrecto mantenga deshabilitado el borrado.
- **FRN-16 (reescrito):**
  - confirmaba el modal con un `if`;
  - tomaba como mensaje cualquier `role="status"`, que podía ser el propio spinner, y sus expresiones regulares coincidían con texto de la tabla;
  - ahora verifica la ruta, el `Bearer`, el bloqueo de toda la fila, que no se pueda enviar una segunda orden, un mensaje propio por cada uno de los seis códigos de D2, que los seis sean distintos, y el mensaje genérico.
- **FRN-19B (reescrito):** ahora prueba:
  - los bordes 69 % y 70 % de D3 en CPU, RAM y disco;
  - "Inaccesible" con `502`, `504` y error de red;
  - el aviso de `stale`, que solo aparece cuando corresponde;
  - el skeleton mientras se espera `/node/status`;
  - la consulta cada 10 s y la pausa con la pestaña oculta.
- **FRN-17A (reescrito):** acepta el diseño implementado: `EventsProvider` montado una vez, consumidores con `useEventsContext()`, y el evento recibido en `ultimoMensaje`. Identifica cada evento por `detalles.tareaId`, porque el cliente no conserva el `id`.
- **FRN-17C, caso del `401`:**
  - reemplazaba `window.location` por un objeto armado con `{...location}`, que pierde `href` y `origin`;
  - ahora monta los providers de la app con su `Toaster`, y acepta la redirección o el aviso `API_UNAUTHORIZED_EVENT`, que es como el cliente delega en el `AuthProvider`;
  - en jsdom no se puede observar el `navigate` del router de la app.

### Contrapruebas

Contra una copia temporal del frontend con implementaciones mínimas correctas pasan **34 de 34**, en `instances-modals`, `instances-operation`, `dashboard-node` y `events-client`. La copia incluye el cliente de eventos real con la ruta, la deduplicación y la limpieza corregidas. Con 13 defectos introducidos, uno por vez, cada uno hace fallar exactamente la prueba que corresponde:

| Defecto | Lo detecta |
|---|---|
| Start sin modal, Stop sin rojo, Delete que acepta cualquier texto, READ_ONLY con energía, OPERATOR con Delete | FRN-15, el caso de cada uno |
| Fila sin bloqueo, sin desbloqueo tras error, dos códigos con el mismo mensaje | FRN-16, el caso de cada uno |
| Umbral en 80 %, red caída sin "Inaccesible", sin polling, sin pausa con pestaña oculta, aviso de `stale` siempre visible | FRN-19B, el caso de cada uno |
