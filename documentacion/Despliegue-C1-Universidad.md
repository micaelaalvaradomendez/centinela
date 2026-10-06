# Despliegue completo de Centinela en el servidor **C1** (Universidad)

> **Qué es esto**: el plan completo, paso a paso, para replicar **todo** Centinela
> (proxy, frontend, backend, base de datos, Redis y el relay `centinela-smtp-brevo`)
> en el servidor **C1** de la universidad, con **un contenedor (LXC) por servicio**,
> igual que en casa.
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

## 1. Arquitectura objetivo en C1

Se mantiene el mismo esquema que ya usamos en casa: **un contenedor por servicio**.

> **Decisión tomada**: se evaluó consolidar todo con Docker Compose en un solo nodo,
> y **se descartó**. Se replica el esquema actual (LXC por servicio) porque el equipo
> ya lo conoce, se opera igual y C1 lo aguanta (probado). El costo es un poco más de
> RAM (cada LXC tiene su propio SO), y se asume.

```
                        Internet
                           │
                           v
                 ┌───────────────────────┐
                 │   C1 (universidad)    │  IP pública FIJA
                 │                       │
                 │  CT-proxy  (nginx)    │  TLS + /api + /swagger
                 │     │           │     │
                 │     v           v     │
                 │  CT-front     CT-api  │
                 │  (nginx SPA)  (Go)    │
                 │                 │     │
                 │        ┌────────┼──────────┐
                 │        v        v          v
                 │    CT-db      CT-redis   CT-relay
                 │   (postgres)             (postfix)
                 │                            │
                 └────────────────────────────┼───────┘
                                              v
                                            Brevo ──> destinatarios
```

### Contenedores a crear en C1

| CT      | Nombre                   | Servicio              | Rol                          |
|---------|--------------------------|-----------------------|------------------------------|
| CT200   | `centinela-proxy`        | nginx                 | Entrada, TLS, proxy          |
| CT201   | `centinela-db`           | PostgreSQL            | Base de datos                |
| CT202   | `centinela-api`          | Go (systemd)          | API                          |
| CT203   | `centinela-front`        | nginx                 | Sirve el SPA                 |
| CT204   | `centinela-redis`        | Redis                 | Caché / tickets              |
| CT205   | `centinela-smtp-brevo`   | postfix               | Relay de correo a Brevo      |

*(Los IDs son una sugerencia; se ajustan a los que estén libres en C1.)*

---

## 2. Inventario: qué hay hoy y cómo se replica

| Componente      | Hoy (casa)                    | En C1                            |
|-----------------|-------------------------------|----------------------------------|
| Proxy / TLS     | CT100 nginx                   | CT200 nginx                      |
| Frontend        | CT103 nginx + SPA             | CT203 nginx + SPA                |
| Backend         | CT102 Go + systemd            | CT202 Go + systemd               |
| Base de datos   | CT101 PostgreSQL              | CT201 PostgreSQL                 |
| Redis           | CT106 Redis                   | CT204 Redis                      |
| Relay correo    | CT107 postfix                 | CT205 postfix                    |
| Red interna     | vmbr0 `10.10.10.0/24`         | bridge de C1 (a definir)         |
| Acceso          | `centinela.tail6bb3f3.ts.net` | dominio/Tailscale de C1          |

**Repositorios que se usan**

- `orchestrator` (general/infra): `compose.yaml`, `docker/`, `documentacion/`, scripts de deploy.
- `backend` (submódulo, `tayraag/centinela-back`): código Go.
- `frontend` (submódulo, `luzpacello/centinela-front`): SPA React.

---

## 3. Presupuesto de recursos (LXC por servicio)

Cada LXC tiene el costo de su propio sistema operativo (~80–120 MB en reposo) más el servicio.

| CT               | RAM asignada    | Disco         | Nota                            |
|------------------|-----------------|---------------|---------------------------------|
| `proxy`          | 256 MB          | 2 GB          | nginx es liviano                |
| `db`             | 512–1024 MB     | 8–10 GB       | el más pesado; guarda estado    |
| `api`            | 256 MB          | 4 GB          | binario Go + logs               |
| `front`          | 256 MB          | 2 GB          | nginx estático                  |
| `redis`          | 256 MB          | 2 GB          | caché                           |
| `relay`          | 256 MB          | 2 GB          | postfix                         |
| **Total**        | **~1.8–2.3 GB** | **~20–22 GB** | + lo que use el host            |

Con **4 GB de RAM** en C1 va cómodo (probado: "se la banca").
Si aprieta, se puede bajar cada CT a 192 MB y la db a 512 MB.

---

## 4. Requisitos previos (antes de tocar C1)

- [ ] Acceso SSH al Proxmox de C1 (usuario + clave).
- [ ] **IP pública fija** confirmada (la que ve internet).
- [ ] Template LXC disponible (`debian-13-standard` o el que haya).
- [ ] Storage para los rootfs (`local-lvm` o equivalente).
- [ ] Bridge de red + gateway + DNS definidos.
- [ ] Puerto 80/443 accesibles (o el que asigne la facultad).
- [ ] DNS o Tailscale para el acceso con nombre.
- [ ] Credenciales de Brevo: **SMTP key** (`SMTP_USER` / `SMTP_PASS`) y acceso al **panel**.
- [ ] El **dump** de la base actual (para migrar usuarios y auditoría).
- [ ] Los repos clonados (o los archivos del relay del PR #1).

---

## 5. Runbook (paso a paso)

### Fase 0 — Preparar el host C1

1. Verificar template, storage y red disponibles:
   ```bash
   pveam list local | grep debian
   pvesm status
   grep -E "iface vmbr|address|bridge" /etc/network/interfaces
   ```
2. Anotar: **bridge**, **gateway** y **DNS** (los que use C1).
3. Verificar la IP pública: `curl -s https://api.ipify.org` (es la que va a Brevo).
4. **DNS**: si el resolver del host es un stub local (`127.0.0.53`), los contenedores
   no lo alcanzan. A los LXC hay que darles un **nameserver alcanzable** (ver Gotcha 1).

### Fase 1 — Crear los 6 contenedores

Para cada uno (ajustando ID, IP y tamaño):

```bash
pct create 20X local:vztmpl/debian-13-standard_13.6-1_amd64.tar.zst \
  --hostname centinela-<servicio> \
  --cores 1 --memory <RAM> --swap 256 \
  --rootfs local-lvm:<GB> \
  --net0 name=eth0,bridge=<bridge>,ip=<IP>/24,gw=<gw>,firewall=1 \
  --nameserver <DNS-alcanzable> \
  --unprivileged 1 --onboot 1 \
  --description "<servicio> de Centinela"
pct start 20X
```

Verificación rápida por CT: `pct exec 20X -- ip -brief a` y `getent hosts deb.debian.org`.

### Fase 2 — Base de datos (CT201) y Redis (CT204)

1. **Postgres**: instalar, crear la base y el usuario de la app.
2. **Migrar los datos actuales**:
   - En casa: `pg_dump` de `centinela_test`.
   - Copiar el dump a C1 y restaurarlo.
3. **Redis**: instalar, setear password y `bind` a la red interna.
4. Verificar: `\dt`, conteo de usuarios, y que el trigger de auditoría exista.

### Fase 3 — Backend (CT202)

1. Instalar el toolchain Go (o compilar afuera y copiar el binario).
2. Crear el `.env` real (permisos `600`) con TODAS las variables (ver Anexo A).
3. Apuntar `SMTP_HOST` al **relay** (IP de CT205) y `SMTP_PORT=25`.
4. Crear el servicio systemd `centinela-api` con `EnvironmentFile=/etc/centinela-api.env`.
5. Arrancar y verificar `GET /api/version`.

### Fase 4 — Frontend (CT203) y Proxy (CT200)

1. Build del SPA (`pnpm build`) con `VITE_API_BASE_URL=/api`.
2. CT203: nginx sirviendo el SPA (fallback `try_files ... /index.html`).
3. CT200: nginx de entrada que:
   - sirve `/` → CT203,
   - proxipea `/api/` → CT202:8080,
   - normaliza `/swagger` → CT202,
   - tiene el TLS (cert real: Tailscale cert o Let's Encrypt de la facultad).

### Fase 5 — Relay `centinela-smtp-brevo` (CT205)

Los archivos ya están en el repo (PR #1): `docker/relay/` + `compose.relay.yaml`.

1. Instalar `postfix`, `libsasl2-modules`, `ca-certificates`.
2. Generar el certificado del relay (SAN = su IP interna) para el STARTTLS.
3. Configurar postfix como relay a Brevo (STARTTLS + clave SMTP) — mismo `main.cf`
   que en casa.
4. Desactivar chroot y copiar el `resolv.conf` (ver Gotcha 2).
5. Probar: cliente → relay → Brevo, con `status=sent` en el log.

### Fase 6 — Brevo

1. Panel → **Settings → Security → Authorized IPs**.
2. **Agregar la IP pública de C1** (dejar la de casa durante la transición).
3. **SMTP keys → Blocking: activado**.
4. Verificar el **remitente** (`centinelauntdf@gmail.com`) activo.
5. **Higiene**: rotar la SMTP key (la actual estuvo expuesta en ClickUp).
6. Recordar el límite: **free = 300 envíos/día** (compartido por toda la cuenta).

### Fase 7 — Verificación end-to-end

1. Alta de un usuario de prueba → debe llegar el correo con la clave temporal.
2. "Olvidé mi contraseña" → debe llegar el OTP.
3. En Brevo: confirmar eventos `delivered`.
4. Integridad: retener un mensaje con `defer_transports`, `postcat` para ver el cuerpo
   intacto (con acentos), y liberarlo.

---

## 6. Gotchas (lo que aprendimos a golpes)

Estos problemas **ya nos costaron tiempo** en casa. En C1 reaparecen si no se tienen en cuenta.

### Gotcha 1 — DNS dentro de los contenedores
- **Síntoma**: `getent hosts smtp-relay.brevo.com` vacío; el correo queda en cola con
  `Host or domain name not found`.
- **Causa**: el contenedor no puede usar el resolver del host (`127.0.0.53`, systemd-resolved).
- **Solución**: darle un `nameserver` alcanzable (el DNS del lab o uno público permitido).

### Gotcha 2 — postfix en chroot no ve el DNS
- **Síntoma**: postfix resuelve bien desde el shell pero al enviar falla.
- **Causa**: postfix corre enjaulado (chroot) y adentro no ve `/etc/resolv.conf`.
- **Solución**: `postconf -F '*/*/chroot = n'` + copiar el `resolv.conf` al jaque.

### Gotcha 3 — Falta `libsasl2-modules`
- **Síntoma**: `SASL authentication failure: No worthy mechs found`.
- **Causa**: falta el paquete con los mecanismos SASL (LOGIN/PLAIN).
- **Solución**: instalar `libsasl2-modules` en el CT del relay.

### Gotcha 4 — Egreso restringido
- **Síntoma**: `apt` no llega a `deb.debian.org:80` (timeout).
- **Causa**: la red bloquea puertos/hosts; solo pasan algunos (ej. 587 a Brevo).
- **Solución**: bajar el `.deb` en el host (que sí tiene internet) y pasarlo al CT con `pct push`.

### Gotcha 5 — Sin ruta IPv6
- **Síntoma**: `apt` intenta IPv6 y falla (`Network is unreachable`).
- **Solución**: forzar IPv4 (`-o Acquire::ForceIPv4=true`).

### Gotcha 6 — STARTTLS interno y certificados
- **Síntoma**: el backend rechaza el STARTTLS del relay (cert no confiable).
- **Causa**: el certificado es self-signed y el cliente valida contra el nombre/IP.
- **Solución**: generar el cert con **SAN** correcto (IP del relay) e **instalarlo como
  confiable** en el backend (`update-ca-certificates`).

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

- [ ] Host C1 con template, storage y red listos.
- [ ] 6 LXC creados y con DNS funcionando.
- [ ] Base migrada (dump/restore) y verificada.
- [ ] Redis arriba con password.
- [ ] Backend arriba con el `.env` completo (`600`) y systemd.
- [ ] Frontend servido y `/api` proxiado por el proxy.
- [ ] Relay arriba, con STARTTLS y cert confiado por el backend.
- [ ] IP de C1 autorizada en Brevo + bloqueo activo.
- [ ] SMTP key rotada.
- [ ] Alta de usuario de prueba → correo `delivered`.
- [ ] Recuperación de contraseña → OTP `delivered`.
- [ ] Rollback anotado (sección 7) por si falla.

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
| `SMTP_HOST`                  | **El relay** (IP del CT205)                 |
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

# Ver logs del relay
pct exec 205 -- tail -f /var/log/postfix-relay.log

# Probar el relay (cliente -> relay -> Brevo)
python3 -c "
import smtplib
from email.message import EmailMessage
m=EmailMessage(); m['From']='centinelauntdf@gmail.com'; m['To']='centinelauntdf@gmail.com'
m['Subject']='Prueba C1'; m.set_content('hola')
s=smtplib.SMTP('<IP-del-relay>',25,timeout=20); s.send_message(m); print('OK')
"
```

### Anexo D — Certificado del relay (STARTTLS interno)

```bash
# Generar (SAN = IP interna del relay)
openssl req -x509 -newkey rsa:2048 -nodes -keyout relay.key -out relay.crt -days 3650 \
  -subj "/CN=centinela-smtp-brevo" -addext "subjectAltName=IP:<IP-del-relay>"

# En el backend: instalar como CA confiable
cp relay.crt /usr/local/share/ca-certificates/centinela-smtp-brevo.crt
update-ca-certificates
```

---

## 10. Diferencias respecto de casa (resumen)

| Tema            | Casa (hoy)                         | C1 (objetivo)                          |
|-----------------|------------------------------------|----------------------------------------|
| Contenedores    | 6 LXC separados                    | 6 LXC separados (igual)                |
| IP de salida    | domiciliaria (dinámica)            | **fija** (universidad)                 |
| Relay           | CT107 (`10.10.10.60`)              | CT205                                  |
| Brevo IPs       | 1 autorizada (casa)                | + IP de C1 autorizada                  |
| Acceso          | Tailscale `centinela.tail6bb3f3`   | dominio/Tailscale de C1                |
| Deploy          | scripts `deploy-back/front.sh`     | los mismos scripts (cambiando host/CT) |
