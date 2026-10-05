# Resultados de las pruebas backend

**Fecha:** 05/10/2026. **Backend:** submódulo en `e1f2df4` (último commit de `main`; nuevos: contrato de la Etapa 1, `FIX-37`, `FIX-39`, `BAC-24A` y parte de `BAC-24B`). **Comando:** `go test -v -count=1 ./...`. 2 corridas completas con el mismo resultado y sin contenedores residuales.

| Métrica | Valor |
|---|---:|
| Casos | 58 |
| Aprueban | **49** |
| Fallan | **8** |
| Omitidos | 1: INF-08B (Brevo rechaza la IP pública con `525 Unauthorized IP`; la prueba se omite a propósito en ese caso) |

## Tareas de `actual.md`

| Tarea | Caso | Resultado | Causa |
|---|---|---|---|
| `FIX-44` (índice y purga, `BAC-18B`) | `cierre_fase_base…` BAC-18B (2) | ❌ índice / ✅ purga | Falta el índice parcial con `jti_access`. La purga horaria **cumple** (ver §Cambios) |
| `BAC-29` contrato | `etapa1_acceptance…` | ❌ | `docs/contrato-etapa1.md` existe, pero no menciona `DELETE`, `exitstatus`, `PROXMOX_ERROR` ni `INSTANCE_INVALID_STATE`, ni que `/node/status` lo puede consultar cualquier usuario autenticado (D1) |
| `BAC-22` telemetría | `etapa1_acceptance…` | ❌ | `GET /api/node/status` responde 404 |
| `BAC-23A` adaptador de IP | `etapa1_acceptance…` | ❌ | No se consultan los endpoints de IP |
| `BAC-25A` pool de UPID | `etapa1_acceptance…` | ❌ | Falta `UPID_WORKERS` |
| `BAC-24A` estado previo y protegidos | `etapa1_acceptance…` (**reescrito**) | ❌ solo por documentación | Cumple todo el comportamiento:<br>• `409 INSTANCE_INVALID_STATE` en las cuatro acciones y en las dos rutas;<br>• `403 INSTANCE_PROTECTED` en shutdown, reboot y stop;<br>• `400 INVALID_ACTION`;<br>• `403` sin permiso;<br>• `202 { upid, tareaId }`;<br>• ninguna escritura en Proxmox al rechazar.<br>Falta `INSTANCE_INVALID_STATE` en `docs/swagger.json` (entregable 4) |
| `BAC-24B` DELETE | `etapa1_acceptance…` (**reescrito**) | ❌ | Ya existe y cumple el `401`, el `403` al OPERATOR con `FULL_ACCESS` y el `403 INSTANCE_PROTECTED`. Pero responde `409 INSTANCE_NOT_STOPPED` en lugar de `INSTANCE_INVALID_STATE` (D2), y `204` sin cuerpo en lugar de `202 { upid, tareaId }`, así que no hay seguimiento ni `TASK_FINISHED` |
| `BAC-25C` reanudación | `etapa1_acceptance…` (**reescrito**) | ❌ | `activeTask` ya se informa (parte 2 cumple). Al reiniciar no se retoma el sondeo: la tarea vencida no queda `FAILED` y la reciente no emite `TASK_FINISHED` |

## Regresión

Pasan todas las demás tareas terminadas, incluidas `LOGIN-04` (desde Go y front ↔ back, 10/10), `FIX-40`, `INF-07B`, `BAC-21C`, `SEC-04`, `FIX-37` y `FIX-39`.

## Cambios en la suite (05/10/2026)

- **BAC-24A (reescrito):** la versión anterior usaba el VMID 100 como protegido, pero en la suite el protegido es el 103 y el 100 no existe en el stub. Tampoco verificaba el `202` ni la ausencia de tráfico hacia Proxmox. Ahora cubre todo el criterio y el entregable 4, y cuenta en el log del stub las escrituras que llegaron.
- **BAC-24B (reescrito):** tenía el mismo VMID 100, y su OPERATOR no tenía `FULL_ACCESS`. Ahora verifica también la auditoría de la orden y el `TASK_FINISHED` con `recursoTipo` `LXC`.
- **BAC-25C (reescrito):**
  - antes no probaba el reinicio y leía `body["data"]`, pero `GET /instances` devuelve una lista;
  - ahora inserta una tarea `RUNNING` reciente y otra de hace 10 min, reinicia el backend y verifica la vencida `FAILED`, la reciente en `activeTask`, y el `TASK_FINISHED` y `COMPLETED` cuando Proxmox la termina;
  - el stub tiene un interruptor nuevo, `/tmp/tareas-en-curso`, que deja todas las tareas en `running`.
- **BAC-18B purga (reescrito):** exigía escribir literal `time.NewTicker(time.Hour)`. El backend usa `NuevoPurgaWorker(db, 1*time.Hour)`, que cumple y además permite acortar el intervalo en las pruebas unitarias. Ahora la prueba:
  - inserta sesiones inactiva, vencida y vigente;
  - reinicia el backend (la purga corre al arrancar) y verifica que solo queden las vigentes;
  - acepta la hora tanto en el ticker como en la llamada que crea la rutina.
  **Pasa.** Por eso el entregable 2 de `FIX-44` no hace falta.
- **Pruebas viejas adaptadas a `BAC-24A`:** `BAC-21B` (energía y auditoría) y `SEC-04` mandaban `start` a la 101, que en el stub está encendida. Con la validación de estado previo ahora reciben `409`, así que usan `stop`.

## Lo que la suite no cubre

- BAC-18B: el tiempo de consulta (< 20 ms).
- BAC-23A: el timeout de 2 s y la concurrencia acotada de las consultas de IP (se verán con `BAC-23B`).
- BAC-25A: la cota de N consultas en paralelo (queda a cargo de las pruebas unitarias del backend).
- BAC-24B: que el tipo de recurso se guarde **antes** de borrar. El stub sigue listando la instancia después del `DELETE`, así que no se puede detectar el error de averiguarlo al final.
- BAC-25C: el registro final en `auditoria`. La entrada de resolución la define `BAC-27`, que todavía no está cargada (ver `test/informe.md`).
- FIX-40: el `502` con el token rechazado (lo cubren las pruebas unitarias del backend).
