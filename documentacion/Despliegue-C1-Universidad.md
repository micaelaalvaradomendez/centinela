# Despliegue completo de Centinela en el servidor **C1** (Universidad)

> **Qué es esto**: el plan completo, paso a paso, para replicar **todo** Centinela
> (frontend, backend, base de datos, Redis y el relay de correo `centinela-smtp-brevo`)
> en el servidor **C1** de la universidad, con recursos limitados.
>
> **Estado**: plan listo para ejecutar. Cuando la facultad habilite C1, se sigue
> esta guía de punta a punta y no hay que pensar nada nuevo.

---

## 0. Objetivo

Que C1 corra **el sistema completo**, con **una sola IP pública fija** autorizada en
Brevo (la de C1), y que el correo salga **cifrado** en los dos tramos.

Resultado esperado:

```
Backend --STARTTLS--> relay centinela-smtp-brevo --STARTTLS + IP autorizada--> Brevo --> destinatario
```

---

## 1. Arquitectura objetivo en C1 (consolidada)

Hoy en casa corremos **6 contenedores LXC** (proxy, db, api, front, redis, relay).
Cada LXC es un SO completo: eso es el lujo que con pocos recursos no nos podemos dar.

En C1 **consolidamos todo en un solo nodo con Docker Compose**: 5 servicios livianos
en vez de 6 sistemas operativos.

```
                            Internet
                               │
                               v
                    ┌─────────────────────────┐
                    │   C1 (universidad)      │  IP pública FIJA
                    │                         │
                    │   nginx  (80/443)       │  SPA + /api + /swagger
                    │     │           │       │
                    │     v           v       │
                    │  front(SPA)  backend(Go)│
                    │                 │       │
                    │        ┌────────┼─────────────┐
                    │        v        v             v
                    │   postgres     redis      relay(postfix)
                    │                              │
                    └──────────────────────────────┼──────────┘
                                                   v
                                                 Brevo ──> destinatarios
```

**Decisiones de diseño**

- **Un solo nginx** sirve el SPA **y** hace de proxy de `/api` (fusiona los viejos CT100 + CT103).
- **Postgres y Redis** como contenedores con volumen persistente.
- **El relay** corre en el mismo nodo: su salida es la **IP fija de C1**, que es lo que Brevo autoriza.
- Todo con **Docker Compose**, un solo `up -d` levanta el sistema.

---

## 2. Inventario: qué hay hoy y cómo se replica

| Componente      | Hoy (casa)                    | En C1 (consolidado)              |
|-----------------|-------------------------------|----------------------------------|
| Proxy / TLS     | CT100 nginx                   | nginx del stack (80/443)         |
| Frontend        | CT103 nginx + SPA             | servicio `frontend` (nginx)      |
| Backend         | CT102 Go + systemd            | servicio `backend` (Go)          |
| Base de datos   | CT101 PostgreSQL              | servicio `db` (postgres)         |
| Redis           | CT106 Redis                   | servicio `redis`                 |
| Relay correo    | CT107 postfix                 | servicio `smtp-relay`            |
| Red interna     | vmbr0 `10.10.10.0/24`         | red Docker interna               |
| Acceso          | `centinela.tail6bb3f3.ts.net` | dominio/Tailscale de C1          |

**Repositorios que se usan**

- `orchestrator` (general/infra): `compose.yaml`, `docker/`, `documentacion/`.
- `backend` (submódulo, `tayraag/centinela-back`): código Go.
- `frontend` (submódulo, `luzpacello/centinela-front`): SPA React.

---

## 3. Presupuesto de recursos

| Servicio     | RAM sugerida   | Rol                   |
|--------------|----------------|-----------------------|
| Postgres     | 256–512 MB     | Base de datos         |
| Redis        | 64 MB          | Caché / tickets       |
| Backend Go   | 64–128 MB      | API                   |
| Nginx        | 32 MB          | SPA + proxy /api      |
| Postfix      | 64 MB          | Relay a Brevo         |
| Total        | ~0.5–1 GB      | 1 vCPU, ~8 GB disco   |

Con **1 vCPU / 1 GB / 8 GB** corre justo. Con **2 vCPU / 2 GB** va cómodo.
Si C1 es más chico, se puede bajar Postgres a 256 MB y Redis a 64 MB.

---

## 4. Requisitos previos (antes de tocar C1)

- [ ] Acceso SSH a C1 (usuario + clave).
- [ ] **IP pública fija** confirmada (la que ve internet).
- [ ] Docker + `docker compose` disponibles (o LXC con `nesting=1` si C1 es un LXC).
- [ ] Puerto 80/443 accesibles (o el que asigne la facultad).
- [ ] DNS o Tailscale para el acceso con nombre.
- [ ] Credenciales de Brevo: **SMTP key** (`SMTP_USER` / `SMTP_PASS`) y acceso al **panel**.
- [ ] El **dump** de la base actual (para migrar los usuarios y la auditoría).
- [ ] Los repos clonados (o los archivos del relay del PR #1).

---

## 5. Runbook (paso a paso)

### Fase 0 — Preparar C1

1. Actualizar el sistema e instalar Docker + Compose.
2. Verificar la IP pública: `curl -s https://api.ipify.org` (anotarla: es la que va a Brevo).
3. **DNS**: si el resolver del host es un stub local (`127.0.0.53`), los contenedores
   no lo alcanzan. Hay que darle DNS explícito al stack (ver Gotcha 1).
4. Si C1 es un LXC: habilitar `nesting=1` para poder correr Docker adentro.

### Fase 1 — Datos (Postgres + Redis)

1. Levantar `db` y `redis` con volumen persistente.
2. Crear la base y el usuario de la app.
3. **Migrar los datos actuales**:
   - En casa: `pg_dump` de `centinela_test`.
   - Copiar el dump a C1 y restaurarlo.
4. Verificar: `\dt` (tablas), conteo de usuarios, y que el trigger de auditoría exista.

### Fase 2 — Backend

1. Build del backend (Go) — por Docker o binario + systemd.
2. Crear el `.env` real (permisos `600`) con TODAS las variables (ver Anexo A).
3. Apuntar `SMTP_HOST` al relay (`smtp-relay`, puerto `25`) y `SMTP_PORT=25`.
4. Arrancar y verificar `GET /api/version`.

### Fase 3 — Frontend + nginx

1. Build del SPA (`pnpm build`) con `VITE_API_BASE_URL=/api`.
2. Un **solo nginx**: sirve el SPA y hace proxy de `/api` al backend.
3. TLS: certificado real (Tailscale cert o Let's Encrypt de la facultad).

### Fase 4 — Relay `centinela-smtp-brevo`

Los archivos ya están en el repo (PR #1): `docker/relay/` + `compose.relay.yaml`.

1. Generar el certificado self-signed del relay con **SAN = nombre del servicio**
   (`smtp-relay`) y su IP interna, para el STARTTLS backend→relay.
2. Levantar el relay con `compose.relay.yaml` (env: `BREVO_HOST/USER/PASS`).
3. Que el backend **confíe** ese certificado (instalarlo en su trust store).
4. Probar: `swaks`/`python smtplib` al relay → ver `status=sent` en los logs.

### Fase 5 — Brevo

1. Panel → **Settings → Security → Authorized IPs**.
2. **Agregar la IP pública de C1** (y dejar la de casa durante la transición).
3. **SMTP keys → Blocking: activado** (si no lo está).
4. Verificar el **remitente** (`centinelauntdf@gmail.com`) activo.
5. **Higiene**: rotar la SMTP key (la actual estuvo expuesta en un comentario de ClickUp).
6. Recordar el límite: **free = 300 envíos/día** (compartido por toda la cuenta).

### Fase 6 — Verificación end-to-end

1. Alta de un usuario de prueba → debe llegar el correo con la clave temporal.
2. "Olvidé mi contraseña" → debe llegar el OTP.
3. En Brevo: confirmar eventos `delivered`.
4. Integridad: retener un mensaje con `defer_transports`, `postcat` para ver el cuerpo
   intacto (con acentos), y liberarlo.

---

## 6. Gotchas (lo que aprendimos a golpes)

Estos son los problemas que **ya nos costaron tiempo** en casa. En C1 van a reaparecer
si no se tienen en cuenta.

### Gotcha 1 — DNS dentro de contenedores
- **Síntoma**: `getent hosts smtp-relay.brevo.com` vacío; el correo queda en cola con
  `Host or domain name not found`.
- **Causa**: el contenedor usa el resolver del host (`127.0.0.53`, systemd-resolved),
  que **no es alcanzable desde adentro**.
- **Solución**: dar DNS explícito al servicio (`dns: [1.1.1.1, 8.8.8.8]` en Compose) o
  el DNS que provea la facultad.

### Gotcha 2 — postfix en chroot no ve el DNS
- **Síntoma**: postfix resuelve bien desde el shell pero al enviar falla.
- **Causa**: postfix corre enjaulado (chroot) y adentro no ve `/etc/resolv.conf`.
- **Solución**: `postconf -F '*/*/chroot = n'` + copiar el `resolv.conf` al jaque.

### Gotcha 3 — Falta `libsasl2-modules`
- **Síntoma**: `SASL authentication failure: No worthy mechs found`.
- **Causa**: falta el paquete con los mecanismos SASL (LOGIN/PLAIN).
- **Solución**: instalar `libsasl2-modules` en el contenedor del relay.

### Gotcha 4 — Egreso restringido
- **Síntoma**: `apt` no llega a `deb.debian.org:80` (timeout).
- **Causa**: la red bloquea puertos/hosts; solo pasan algunos (ej. 587 a Brevo).
- **Solución**: bajar el `.deb` en el host (que sí tiene internet) y pasarlo al contenedor.

### Gotcha 5 — Sin ruta IPv6
- **Síntoma**: `apt` intenta IPv6 y falla (`Network is unreachable`).
- **Solución**: forzar IPv4 (`-o Acquire::ForceIPv4=true`).

### Gotcha 6 — STARTTLS interno y certificados
- **Síntoma**: el backend rechaza el STARTTLS del relay (cert no confiable).
- **Causa**: el certificado es self-signed y el cliente valida contra el nombre/IP.
- **Solución**: generar el cert con **SAN** correcto (el hostname del relay) e
  **instalarlo como confiable** en el backend (`update-ca-certificates`).

### Gotcha 7 — Cómo se comporta el cliente SMTP (go-mail)
- STARTTLS es **oportunista**: si el server no lo ofrece, sigue en texto plano.
- AUTH solo se intenta **si el servidor lo anuncia**.
- Por eso apuntar el backend a un relay sin TLS/AUTH funciona, y con TLS funciona igual.

### Gotcha 8 — Brevo bloquea por IP de origen (`525`)
- **Síntoma**: `525 "5.7.1 Unauthorized IP address"` y las altas/recuperaciones fallan.
- **Causa**: el bloqueo de IPs para SMTP keys está activo y la IP de salida no está autorizada.
- **Solución**: el relay (IP fija) + autorizar esa única IP. Nunca "abrir" el bloqueo.

### Gotcha 9 — La auditoría es inmutable
- **Síntoma**: `DELETE` de un usuario falla por FK `RESTRICT` + trigger.
- **Causa**: `auditoria` tiene un trigger que prohíbe `DELETE/UPDATE/TRUNCATE`.
- **Consecuencia**: un usuario con auditoría **no se puede borrar físicamente**.
  Para "dar de baja" se usa soft-delete + anonimización (no borrado físico).

---

## 7. Plan B / rollback

Si algo sale mal en C1, **volver atrás es simple**:

1. En el backend de C1, devolver `SMTP_HOST=smtp-relay.brevo.com`, `SMTP_PORT=587`
   (y `SMTP_USER/PASS/FROM` de Brevo).
2. En Brevo, **desactivar** el bloqueo de SMTP keys (o autorizar la IP que corresponda).
3. El resto del stack no depende del relay: se puede dejar como está.

El sistema sigue andando en casa (PRUEBAS) mientras tanto.

---

## 8. Checklist del día de la migración

- [ ] C1 con Docker + IP fija confirmada.
- [ ] DNS del stack configurado (Gotcha 1).
- [ ] Base migrada (dump/restore) y verificada.
- [ ] Backend arriba con el `.env` completo (`600`).
- [ ] Frontend servido y `/api` proxiado.
- [ ] Relay arriba, con STARTTLS y cert confiado por el backend.
- [ ] IP de C1 autorizada en Brevo + bloqueo activo.
- [ ] SMTP key rotada.
- [ ] Alta de usuario de prueba → correo `delivered`.
- [ ] Recuperación de contraseña → OTP `delivered`.
- [ ] Rollback anotado (Fase 7) por si falla.

---

## 9. Anexos

### Anexo A — Variables del backend (`.env` real, permisos `600`)

| Variable                     | Para qué                                    |
|------------------------------|---------------------------------------------|
| `AUTH_MODE`                  | Modo de autenticación (`postgres`)          |
| `DATABASE_URL` / `DB_DSN`    | Conexión a PostgreSQL                       |
| `HTTP_ADDR`                  | Puerto de la API (`:8080`)                  |
| `JWT_SECRET`                 | Firma de tokens                             |
| `TOTP_ENCRYPTION_KEY`        | Cifrado de los secretos 2FA                 |
| `ENABLE_SWAGGER`             | Expone `/swagger`                           |
| `ALLOWED_ORIGINS`            | CORS                                        |
| `PROXMOX_*`                  | Integración con Proxmox                     |
| `REDIS_*`                    | Conexión a Redis                            |
| `EMAIL_PROVIDER=smtp`        | Proveedor de correo real                    |
| `SMTP_HOST`                  | **El relay** (`smtp-relay`)                 |
| `SMTP_PORT`                  | `25` (interno)                              |
| `SMTP_USER/PASS`             | Credenciales SMTP de Brevo                  |
| `SMTP_FROM`                  | Remitente (`centinelauntdf@gmail.com`)      |

### Anexo B — Archivos del relay (ya versionados)

- `docker/relay/Dockerfile`
- `docker/relay/main.cf.tpl`
- `docker/relay/entrypoint.sh`
- `docker/relay/README.md`
- `compose.relay.yaml`
- `.env.relay.example`

### Anexo C — Comandos clave

```bash
# IP pública (lo que ve Brevo)
curl -s https://api.ipify.org

# Migrar la base (desde casa)
ssh root@proxmox 'pct exec 101 -- su - postgres -c "pg_dump -d centinela_test -Fc -f /var/lib/postgresql/centinela_test.dump"'
# copiar el dump y restaurar en C1:
pg_restore -d centinela_test --clean --if-exists /ruta/centinela_test.dump

# Levantar el stack
docker compose up -d --build

# Ver logs del relay
docker logs -f centinela-smtp-brevo

# Probar el relay (cliente -> relay -> Brevo)
python3 -c "
import smtplib
from email.message import EmailMessage
m=EmailMessage(); m['From']='centinelauntdf@gmail.com'; m['To']='centinelauntdf@gmail.com'
m['Subject']='Prueba C1'; m.set_content('hola')
s=smtplib.SMTP('smtp-relay',25,timeout=20); s.send_message(m); print('OK')
"
```

### Anexo D — Certificado del relay (STARTTLS interno)

```bash
# Generar (SAN = nombre del servicio en la red Docker)
openssl req -x509 -newkey rsa:2048 -nodes -keyout relay.key -out relay.crt -days 3650 \
  -subj "/CN=smtp-relay" -addext "subjectAltName=DNS:smtp-relay"

# En el backend: instalar como CA confiable
cp relay.crt /usr/local/share/ca-certificates/centinela-smtp-brevo.crt
update-ca-certificates
```

---

## 10. Diferencias respecto de casa (resumen)

| Tema              | Casa (hoy)                          | C1 (objetivo)                       |
|-------------------|-------------------------------------|-------------------------------------|
| Contenedores      | 6 LXC separados                     | 1 nodo con 5 servicios Docker       |
| IP de salida      | domiciliaria (dinámica)             | **fija** (universidad)              |
| Relay             | CT107 (`10.10.10.60`)               | servicio `smtp-relay`               |
| Brevo IPs         | 1 autorizada (casa)                 | + IP de C1 autorizada               |
| Acceso            | Tailscale `centinela.tail6bb3f3`    | dominio/Tailscale de C1             |
| Deploy            | scripts `deploy-back/front.sh`      | `docker compose up -d --build`      |
