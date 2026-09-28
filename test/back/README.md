# Pruebas de aceptación del backend

Estas pruebas contrastan el backend (`backend/`, Go + Gin) con los criterios de éxito de `documentacion/actual.md` y funcionan como prueba de regresión de `documentacion/terminado.md`. Se ejecutan contra el stack real en Docker Compose:

| Servicio | Imagen | Uso |
|---|---|---|
| `db` | `postgres:16-alpine` | Base aislada, inicializada con `backend/scripts/init.sql` |
| `backend` | `docker/backend.Dockerfile` (compila `backend/`) | API real en `127.0.0.1:18080` |
| `proxmox` | `nginx:1.27-alpine` + [`proxmox-stub/nginx.conf`](proxmox-stub/nginx.conf) | Stub de la API de Proxmox VE con un inventario fijo: `101` qemu, `102` lxc, `103` qemu, más un nodo y un storage que el backend debe descartar |

## Principios

- **Una prueba que pasa verifica el criterio de éxito.** No se admiten `t.Logf` + `return` ni aserciones que acepten "cualquiera de varios estados".
- **Las credenciales se obtienen igual que un usuario real.** Las claves temporales y los códigos de recuperación se leen del `MockEmailService`, es decir, de los logs del contenedor `backend` (helper `mailedSecret`). Los tokens salen del circuito `login → 2fa/qr → 2fa/verify` (helpers `loginWithTOTP` y `createActiveUser`). Solo se firman tokens en la prueba (`signedAccessToken`) para el actor administrador de preparación.
- **Integración front ↔ back.** Cuando el frontend cambia un contrato, se agrega una prueba que reproduce su petición exacta. Ejemplo: `SEC-01 SEC-02 integracion…`.

## Cobertura

| Archivo | Tareas |
|---|---|
| `backend_acceptance_test.go` | BAC-01, BAC-02, BAC-03, BAC-04, BAC-05, BAC-06, BAC-06B, BAC-09, BAC-10, BAC-11, BAC-12, LOGIN-01, LOGIN-03 |
| `login04_acceptance_test.go` | LOGIN-04: circuito alta → clave por correo → 2FA → cambio obligatorio → rol e instancias → resets |
| `password_recovery_acceptance_test.go` | BAC-13, BAC-15, BAC-16, BAC-18, BAC-19, BAC-20, BAC-21 |
| `resource_access_acceptance_test.go` | BAC-07, FIX-16/BAC-08, BAC-14, **SEC-04** |
| `session_security_acceptance_test.go` | **BAC-17**, **SEC-01**, integración SEC-01/SEC-02, FIX-08 |

## Requisitos

- Docker con `docker compose` (la primera corrida descarga `postgres:16-alpine` y `nginx:1.27-alpine`).
- Go 1.26 o superior.
- Puertos locales `15433` y `18080` libres. Se pueden cambiar con `BACKEND_TEST_DB_PORT`, `BACKEND_TEST_API_PORT` y `BACKEND_TEST_API_URL`.

## Ejecutar

```bash
cd test/back
go test -v -count=1 ./...        # levanta el stack, prueba y lo elimina (~75 s)
go test -short ./...             # solo compila, sin Docker
```

Antes de correr la suite se trae el último commit del submódulo (`git -C backend fetch origin --prune && git -C backend checkout -B main origin/main --force && git -C backend reset --hard origin/main`); ante un conflicto prevalece el remoto. Guardá la salida con `| tee archivo`, no con `> archivo`. `CENTINELA_ROOT` y `CENTINELA_FRONTEND_DIR` permiten apuntar a otra copia, pero no hacen falta en el procedimiento estándar.

Si una corrida se interrumpe y deja contenedores levantados, los puertos quedan ocupados. Para limpiar:

```bash
docker ps -a --format '{{.Names}}' | grep -o 'centinela-back-tests-[0-9]*' | sort -u \
  | xargs -I{} docker compose -p {} -f compose.yaml down --volumes --remove-orphans
```

Resultado vigente: [`RESULTADOS.md`](RESULTADOS.md). Informe consolidado: [`../informe.md`](../informe.md).
