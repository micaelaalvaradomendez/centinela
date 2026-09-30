# Informe de estado de tareas verificado por pruebas

**Fecha de ejecución:** 30/09/2026
**Alcance:** las 17 tareas de [documentacion/actual.md](../documentacion/actual.md), regresión de [documentacion/terminado.md](../documentacion/terminado.md) y verificación del simulador de Proxmox del backend contra el Proxmox real.

| Componente | Revisión probada | Commits nuevos desde el 29/09 |
|---|---|---|
| Backend | `e6e7dc5` (último commit de `main`) | Ninguno desde la corrida anterior. Entre el 29 y el 30/09 llegaron `725d436`/`82add96` (protección de VMIDs de infraestructura y API Token único), `de407a5` (simulador de Proxmox para desarrollo local) y `e6e7dc5` (`ENABLE_SWAGGER` en `.env.example`) |
| Frontend | `749194e` (último commit de `main`) | Ninguno |

**Cómo se ejecutó:** cada submódulo se actualizó a su último commit de `origin/main` (`fetch` + `checkout -B main origin/main --force` + `reset --hard`); ante un conflicto prevalece el remoto. Las suites se corrieron sobre esas carpetas, **dos veces**, con el mismo resultado y sin contenedores residuales.

---

## 1. Resumen

| Suite | Casos | Aprueban | Fallan | Omitidos |
|---|---:|---:|---:|---:|
| Backend (`test/back`) | 43 | 30 | **13** | 0 |
| Frontend (`test/front`) | 95 | 74 | **11** | 10 |
| **Total** | **138** | **104** | **24** | **10** |

- Los 10 casos omitidos del frontend son la prueba integral **LOGIN-04** (`login04-e2e.test.ts`). Corre desde `test/back` contra el backend real, y ahí se cuenta como un caso: `TestLOGIN04IntegracionFrontBack`, 9/10 pasos.
- **Avance desde el 29/09: ninguno.** No hubo commits nuevos en los submódulos y todas las tareas de `actual.md` siguen igual. Los fallos nuevos respecto al informe anterior son de pruebas agregadas para tareas que no tenían cobertura: FRN-17C, BAC-17A y FIX-31.
- **Regresión:** todo lo que está en `terminado.md` sigue en verde, salvo la regresión conocida de `FRN-13` (`FIX-28`).
- **Todos los fallos son del producto.** El único cambio de infraestructura de las pruebas fue adaptar el stub de Proxmox al API Token obligatorio (`725d436`).

---

## 2. Cobertura y estado de las tareas de `actual.md`

| Tarea | Pruebas | Resultado | Qué falta |
|---|---|---|---|
| `SEC-03` | `navigation.test.tsx` (4) | ❌ 1/4 | No existen `usePermissions` ni `PermissionGate`, y el menú le muestra "Auditoría" al OPERATOR |
| `BAC-16B` | `cierre_fase_base…` (backend con `EMAIL_PROVIDER=smtp` + Mailpit) | ❌ | `main.go` siempre usa `MockEmailService` |
| `BAC-17A` | `cierre_fase_base…` **(nuevo)** | ❌ | Falta `go-redis` en `go.mod`, el adaptador `adapters/secondary/redis` y los puertos `GetDel`/`Publish`/`Subscribe`. El backend no abre ninguna conexión a Redis |
| `BAC-17B` | `cierre_fase_base…` | ❌ | Login + 2FA + 10 refresh crean 13 filas activas; después del logout quedan 11 |
| `BAC-18B` | `cierre_fase_base…` | ❌ | Faltan el índice parcial, las particiones y los índices compuestos |
| `BAC-21B` | `puente_etapa1…` (3) | ❌ 0/3 | `GET /instances` no trae los campos nuevos; `/status/:action` → 404; `DELETE` → 405 |
| `BAC-21C` | `puente_etapa1…` (2) | ❌ 0/2 | `/events/ticket` y `/events` → 404 |
| `FRN-17C` | `events-client.test.tsx` (3) **(nuevo)** | ❌ 0/3 | No existe `useEvents` |
| `FIX-27` | `admin-users.test.tsx` (2) | ❌ 1/2 | El 502 del correo muestra un mensaje genérico |
| `FIX-28` | `session-security.test.ts`, `session_security…` y LOGIN-04 paso 7 | ❌ | El logout sale sin Bearer y la sesión sigue activa en el servidor |
| `FIX-29` | `password-change` y `recover-password` | ❌ | El validador no chequea la mayúscula |
| `FIX-30` | `navigation.test.tsx` | ❌ | Depende de `SEC-03` |
| `FIX-31` | `cierre_fase_base…` **(ahora verifica, antes se omitía)** | ❌ | No hay ninguna configuración de Nginx versionada con `listen 443 ssl` (se revisó `docker/nginx.conf` y los dos repos) |
| `INF-06A` | `cierre_fase_base…` | ❌ | Redis no está publicado en `127.0.0.1:6379`, `REDIS_ADDR=redis:6379`, y la contraseña del compose es distinta a la de `.env.example` |
| `INF-08B` | `cierre_fase_base…` (valida las variables y **se autentica** contra el SMTP) | ❌ | Las 6 variables no están en `backend/.env.example` (falta el `push`) |
| `INF-06B` | Sin prueba automatizada | — | Redis en `vmbr1`: no se alcanza desde el repo |
| `INF-08A` | Sin prueba automatizada | — | SMTP del servidor: no se alcanza desde el repo |

---

## 3. Prueba integral LOGIN-04 (código real del frontend contra el backend real)

`TestLOGIN04IntegracionFrontBack`: 9 de 10 pasos en verde. Funcionan de punta a punta:
- el alta con la clave por correo;
- el 2FA con la cookie HttpOnly;
- el cambio y la recuperación de contraseña;
- los permisos por instancia;
- la renovación silenciosa;
- los resets administrativos;
- la auditoría.

**Falla el paso 7:** después de "Cerrar sesión", el access token sigue respondiendo 200 en el servidor (`FIX-28`). **LOGIN-04 no se puede aprobar hasta que se corrija `FIX-28`.**

---

## 4. Simulador de Proxmox del backend (`cmd/proxmox-simulador`, commit `de407a5`)

### 4.1 Flujo del backend contra el simulador
Se levantó el backend de `main` con `PROXMOX_URL` apuntando al simulador, con las mismas credenciales de `documentacion/api-proxmox.md`:

| Acción | Resultado |
|---|---|
| `GET /api/instances` | ✅ 200 con las 8 instancias del simulador, con el contrato de BAC-14 |
| `GET /api/instances/110` | ✅ 200 |
| `POST /api/instances/100/stop` (VMID protegido) | ✅ 403 `INSTANCE_PROTECTED`, sin llegar a Proxmox |
| `POST /api/instances/110/start` | ✅ 202 con UPID; la 110 pasa de `stopped` a `running` cuando termina la tarea |
| `GET /api/instances/424242` | ✅ 404 `INSTANCE_NOT_FOUND` |
| `stop` inmediatamente después de un `start` | ❌ **502 `PROXMOX_UNAVAILABLE`**. Proxmox responde `500 can't lock file…` (instancia ocupada) y el backend lo reporta como "servidor caído" → **`FIX-33`** |

### 4.2 Simulador vs Proxmox real (mismas credenciales, solo lecturas sobre el real)

| Endpoint | Comparación |
|---|---|
| `cluster/resources` (con y sin `?type=vm`) | ✅ Misma estructura. `hastate` aparece solo en instancias con HA, igual que en el real |
| `nodes/{node}/status` | ✅ Igual |
| `nodes/{node}/lxc` | ✅ Igual |
| `lxc/{vmid}/rrddata` | ✅ Igual |
| `lxc/{vmid}/snapshot` | ✅ Igual (los campos dependen de si hay snapshots) |
| Errores: VMID inexistente, nodo inexistente, ruta inexistente y token inválido | ✅ Mismos códigos (500, 500, 501 y 401) y los mismos mensajes |
| `lxc/{vmid}/config` | ❌ `unprivileged` es número en el real y string en el simulador |
| `lxc/{vmid}/status/current` | ❌ Falta `ha: {managed: 0}` en el simulador |
| `cluster/nextid` | ❌ El simulador responde 501; el real devuelve `"106"` (lo necesita RF-07) |
| `nodes/{node}/tasks` | ❌ El simulador responde 501 (útil para BAC-25C y RNF-04) |
| Cuerpo del 401 | 🟡 El real responde sin cuerpo; el simulador devuelve un JSON. No afecta al backend |
| `nodes/{node}/qemu` | ➖ No comparable: el servidor real no tiene VMs |

**Conclusión:** el simulador está bien hecho en lo que el backend usa hoy (inventario, acciones, UPID y errores). Las diferencias quedaron registradas como **`FIX-32`** en `documentacion/etapa1.md`, junto con **`FIX-33`**, que es un problema del backend y no del simulador.

---

## 5. Estado del repositorio y configuración de las pruebas

- Los submódulos están en el último commit y sin cambios locales.
- **Cambios en las pruebas:**
  - el stub de Proxmox exige API Token;
  - `compose.yaml` suma Redis y `PROXMOX_PROTECTED_VMIDS=103`;
  - nuevos tests de BAC-21B, BAC-21C, BAC-17A, FRN-17C y FIX-31, y la prueba integral LOGIN-04.
- **`backend/.env.example` cambió el ejemplo de Proxmox:** usa la IP interna `10.10.20.1` y el token `centinela-api@pve!backend-token`. Las credenciales que funcionan desde esta máquina (Tailscale) son las de `api-proxmox.md` (`100.81.49.19`, `centi-api@pve!backend-token`). Conviene que el equipo confirme cuál es el vigente.
- **Seguridad:** `documentacion/api-proxmox.md` tiene el secreto real del token de Proxmox versionado. Hay que regenerarlo y sacarlo del documento.

---

## 6. Cómo reproducir

```bash
for s in backend frontend; do
  git -C $s fetch origin --prune && git -C $s checkout -B main origin/main --force && git -C $s reset --hard origin/main
done
(cd test/back && go test -v -count=1 ./... | tee /tmp/back.log)   # incluye LOGIN-04 front ↔ back
pnpm --dir test/front test
```

Detalle por suite: [test/back/RESULTADOS.md](back/RESULTADOS.md) y [test/front/RESULTADOS.md](front/RESULTADOS.md).
