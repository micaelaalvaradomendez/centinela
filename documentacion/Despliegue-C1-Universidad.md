# Despliegue completo de Centinela en el servidor **C1** (Universidad)

> **Qué es esto**: el plan completo, paso a paso, para replicar **todo** Centinela
> (proxy, frontend, backend, base de datos, Redis y el relay `centinela-smtp-brevo`)
> en **C1**, el nodo que la cátedra asigna al equipo dentro del **laboratorio ING3**,
> con **un contenedor (LXC) por servicio**, igual que en casa.
>
> **Estado**: plan definido y **verificado en lo esencial** sin root. Para ejecutar:
> autorizar la IP de salida en Brevo, probar el firewall con un contenedor y, cuando la
> cátedra habilite C1, replicar.
>
> **Revisiones**:
> - **Rev 2 — 06/10/2026**: re-análisis con el entorno real del laboratorio.
> - **Rev 3 — 06/10/2026**: verificación en vivo (solo lectura) del nodo y **camino sin
>   root** para el correo (el equipo no tiene root; el egreso del diseño lo requiere).
>   Evidencia en el Anexo E.

---

## 0. Objetivo

Que C1 corra **el sistema completo**, con el **correo saliendo por el reenviador de usuario
del host hacia Brevo**, y que el correo salga **cifrado** en los dos tramos.

Resultado esperado:

```
Backend --STARTTLS--> relay centinela-smtp-brevo --STARTTLS--> Brevo --> destinatario
```

---

## 1. El entorno real (C1 en ING3)

Cómo es el nodo (datos de referencia del laboratorio, verificados en C1):

| Recurso          | Valor                          |
|------------------|--------------------------------|
| Nodo             | Proxmox VE 8.4.21              |
| CPU              | 2 vCPU                         |
| Memoria          | ~4,1 GB de RAM                 |
| Disco            | ~86 GB (37,6 + 48,6)           |
| Acceso           | Solo por Tailscale             |
| Cuenta del grupo | `infra@pve`                    |
| Root             | **No disponible** (verificado) |

**Detalle:**

- Es un nodo **independiente** (no es clúster): se administra desde su panel web.
- Se entra **solo por Tailscale**: sin Tailscale encendido y conectado, el nodo no responde.
- La cuenta del grupo crea y administra máquinas y **solo** las conecta a las redes internas.
- **Sin root**: el equipo nunca tuvo root en el nodo (`sudo` pide contraseña). La
  administración es por panel + usuario SSH (`nicolas.lencina`).

### Las redes

| Red       | Bridge | Rango        | Salida       |
|-----------|--------|--------------|--------------|
| Campus    | vmbr0  | 10.7.93.0/24 | solo el host |
| Interna A | vmbr1  | 10.10.1.0/24 | sin salida   |
| Interna B | vmbr2  | 10.10.2.0/24 | sin salida   |
| Interna C | vmbr3  | 10.10.3.0/24 | sin salida   |

**Detalle:**

- Las internas **no tienen uplink, ni gateway, ni NAT** hoy: los huéspedes no navegan.
- El **host sí tiene IP en cada bridge interna** (`10.10.1.1`, `10.10.2.1`, `10.10.3.1`):
  ahí puede escuchar un **proceso de usuario** (sin root) que dé salida puntual.
- El nodo **no rutea entre redes internas**: todo el stack va en **una sola bridge** (`vmbr1`).
- El **host** sale a internet por `vmbr0` (la salida del campus) y por ahí corre Tailscale.

### La consecuencia principal (el "servicio de msj")

El diseño original del lab (doc 10) prevé "red interna **con salida a internet** y sin
acceso al campus", pero esa etapa (**NAT**) **requiere root** y **no está aplicada**.
Sin root no se puede aplicar; entonces el correo se resuelve con un **reenviador de usuario
en el host** (sección 4) y la instalación se hace **preparando los CT afuera**.

---

## 2. Arquitectura objetivo en C1

Se mantiene el esquema de casa: **un contenedor por servicio**. Se evaluó consolidar
con Docker Compose y **se descartó** (menos componentes a gestionar).

```
                      Notebooks del equipo (Tailscale)
                                 │
                                 v   (publicación por el host, Fase 4)
                ┌────────────────────────────────┐
                │            C1 (ING3)           │
                │                                │
                │  CT-proxy  (nginx)  10.10.1.10 │
                │     │           │              │
                │     v           v              │
                │  CT-front     CT-api           │
                │  10.10.1.13   10.10.1.12       │
                │                 │              │
                │        ┌────────┼─────────┐    │
                │        v        v         v    │
                │    CT-db      CT-redis   CT-relay
                │  10.10.1.11  10.10.1.14  10.10.1.15
                │                             │  │
                └─────────────────────────────┼──┘
                                              v  10.10.1.1:2587 (reenviador de usuario)
                                            Brevo (587) ──> destinatarios
```

### Contenedores a crear

| CT  | Nombre                 | Servicio     | IP sugerida |
|-----|------------------------|--------------|-------------|
| 200 | `centinela-proxy`      | nginx        | 10.10.1.10  |
| 201 | `centinela-db`         | PostgreSQL   | 10.10.1.11  |
| 202 | `centinela-api`        | Go (systemd) | 10.10.1.12  |
| 203 | `centinela-front`      | nginx        | 10.10.1.13  |
| 204 | `centinela-redis`      | Redis        | 10.10.1.14  |
| 205 | `centinela-smtp-brevo` | postfix      | 10.10.1.15  |

*(IDs e IPs son una sugerencia; se ajustan a lo que esté libre.)*

---

## 3. Presupuesto de recursos (LXC por servicio)

| CT        | RAM asignada    | Disco         | Nota                         |
|-----------|-----------------|---------------|------------------------------|
| `proxy`   | 256 MB          | 2 GB          | nginx es liviano             |
| `db`      | 512–1024 MB     | 8–10 GB       | el más pesado; guarda estado |
| `api`     | 256 MB          | 4 GB          | binario Go + logs            |
| `front`   | 256 MB          | 2 GB          | nginx estático               |
| `redis`   | 256 MB          | 2 GB          | caché                        |
| `relay`   | 256 MB          | 2 GB          | postfix                      |
| **Total** | **~1,8–2,3 GB** | **~20–22 GB** | + lo que use el host         |

**Detalle:**

- Entra cómodo en **4,1 GB / ~86 GB**.
- **La CPU es el recurso fino**: 2 vCPU para todo el nodo. Operaciones de a una.
- **Nada de runners de CI en el nodo**.

---

## 4. El camino del correo (sin root)

**Cómo sale** (verificado el 06/10, solo lectura):

1. El **host ya tiene internet** (la salida del campus). No hace falta NAT para esto.
2. Un **proceso de usuario** (no root) en el host escucha en `10.10.1.1:2587` y reenvía a
   `smtp-relay.brevo.com:587`:

```bash
socat TCP-LISTEN:2587,bind=10.10.1.1,fork,reuseaddr TCP:smtp-relay.brevo.com:587
```

3. El **relay** (CT205) entrega a `[10.10.1.1]:2587`.
4. Brevo recibe desde la **IP pública del campus**: `170.210.70.248` (medida, estable en
   las muestras; sin IPv6).

**Persistencia sin root:** `cron` de usuario (`@reboot`) con un watchdog simple; `tmux`
también está disponible.

**Pendientes antes de ejecutar:**

| #   | Pendiente                          | Quién      |
|-----|------------------------------------|------------|
| 1   | Autorizar la IP de egreso en Brevo | Nicolás    |
| 2   | Rotar la clave SMTP                | Nicolás    |
| 3   | Probar el firewall (CT → host)     | Despliegue |
| 4   | Probar el reenviador con un CT     | Despliegue |

**Detalle:**

- **1** — `170.210.70.248` en Authorized IPs, con el bloqueo activo.
- **2** — la clave actual estuvo expuesta.
- **3** — CT → `10.10.1.1:2587`; si el firewall corta, la regla se agrega desde el panel.
- **4** — extremo a extremo: cliente → relay → Brevo.

**Plan B (fuerza mayor):** si el camino directo no funcionara, el correo sale por el
**relay de casa** con un puente Tailscale armado del lado nuestro (los equipos compartidos
están en cuarentena: reciben conexiones, no las inician; el puente no puede depender del
nodo). No es el plan principal.

**Si algún día hay root:** se aplica el egreso documentado del lab (NAT + bloqueos) y el
relay pasa a Brevo directo, sin cambiar nada más.

**Instalación (sin salida general):** los CT se **preparan afuera** (en casa) y se traen
con `vzdump`/restore; alternativa futura: un proxy de usuario para `apt`.

### Verificado en casa (30/10/2026)

El relay de casa salía a Brevo con la IP pública autorizada. La cola se trababa porque el
firewall del host solo dejaba salir al proxy: las reglas de 587/465 habían quedado en un
snapshot viejo. El fix que funcionó fue **acotar** el egreso al relay, sin abrirlo a todos:

```bash
iptables -t nat -A POSTROUTING -s <IP-RELAY>/32 -o <UPLINK> -p tcp -m multiport --dports 587,465 -j MASQUERADE
iptables -A FORWARD -i <BRIDGE> -o <UPLINK> -s <IP-RELAY>/32 -p tcp -m multiport --dports 587,465 -j ACCEPT
```

Resultado: TCP 587 + STARTTLS TLS 1.3 con el certificado de Brevo, cola vacía y 8 mensajes
con `250 2.0.0 OK`. En C1 no hace falta: el relay vive en `vmbr1` (con SNAT) y Brevo autoriza
`170.210.70.248`. Si el 587 estuviera bloqueado, probar 465; si todo el SMTP saliente está
bloqueado, relay en VPS por 443 o API HTTP de Brevo.

---

## 5. Runbook (paso a paso)

### Fase 0 — Preparar

1. Pendientes de la sección 4 (Brevo + reenviador + firewall).
2. Verificar que el reenviador esté arriba: `ss -tln | grep 2587`.
3. Anotar de C1: IDs libres, almacenamiento (`local`, `local-lvm`) y bridge (`vmbr1`).
4. Certificado: cert del nodo (Tailscale) o self-signed.

### Fase 1 — Crear los 6 contenedores

Desde el panel con `infra@pve` (o por `pct` si se hace por SSH):

```bash
pct create 20X local:vztmpl/debian-13-standard_...tar.zst \
  --hostname centinela-<servicio> --cores 1 --memory <RAM> --swap 256 \
  --rootfs local-lvm:<GB> \
  --net0 name=eth0,bridge=vmbr1,ip=10.10.1.1X/24,firewall=1 \
  --unprivileged 1 --onboot 1 \
  --description "<servicio> de Centinela"
pct start 20X
```

**Detalle:** **sin gateway** (la red no tiene salida general); para llegar al host alcanza
la misma subred (`10.10.1.1`). Nombres internos por `/etc/hosts`.

### Fase 2 — Base de datos (201) y Redis (204)

1. Postgres y Redis: instalados **desde la imagen preparada** (o vía proxy).
2. **Migrar los datos**: `pg_dump` en casa y traer el dump por `pct push`/restore.
3. Verificar: `\dt`, usuarios y trigger de auditoría.

### Fase 3 — Backend (202)

1. Binario Go (o toolchain si hubiera salida).
2. `.env` real (permisos `600`) — Anexo A. `SMTP_HOST` = relay (`10.10.1.15`),
   `SMTP_PORT=25`.
3. Servicio systemd `centinela-api`; verificar `GET /api/version`.
4. La app alcanza la **API de Proxmox** del host por la bridge interna (confirmar IP).

### Fase 4 — Frontend (203) y Proxy (200)

1. Build del SPA y nginx con TLS.
2. **Publicar al equipo** con el mismo patrón del correo: un proceso de usuario en el host
   escuchando en la IP Tailscale (`100.116.49.111`) en un **puerto no privilegiado** y
   reenviando al CT200:

```bash
socat TCP-LISTEN:8443,bind=100.116.49.111,fork,reuseaddr TCP:10.10.1.10:443
```

El equipo entra por `https://100.116.49.111:8443/`.

### Fase 5 — Relay `centinela-smtp-brevo` (205)

1. Postfix nativo, igual que en casa (no Docker).
2. `relayhost = [10.10.1.1]:2587`. El STARTTLS se negocia extremo a extremo con Brevo
   (el reenviador es transparente). Si algún día se activa validación de certificado,
   usar `smtp_tls_policy_maps` con `match=smtp-relay.brevo.com`.
3. Probar: cliente → relay → Brevo con `status=sent`.

### Fase 6 — Brevo

1. Autorizar `170.210.70.248` (Authorized IPs) con el **bloqueo activo**.
2. **Rotar la SMTP key** (la actual estuvo expuesta).
3. Recordar: free = **300 envíos/día**.

### Fase 7 — Verificación end-to-end

1. Alta de usuario de prueba → correo con clave temporal.
2. Recuperación → OTP.
3. Brevo: eventos `delivered`.
4. Integridad con `postcat`.

---

## 6. Gotchas

Los de casa (siguen valiendo): DNS del contenedor; chroot de postfix; faltante de
`libsasl2-modules`; egreso restringido; sin ruta IPv6; STARTTLS y certificados; go-mail
oportunista; bloqueo de Brevo (`525`); auditoría inmutable.

De C1:

### Gotcha 10 — Sin salida general
- Los CT no navegan: `apt` no funciona adentro. Los CT se preparan afuera.

### Gotcha 11 — Una sola bridge
- El nodo no rutea entre `vmbr1/2/3`: todo el stack en `vmbr1`.

### Gotcha 12 — La IP que ve Brevo es la del campus
- `170.210.70.248`; puede ser compartida. Mantener el bloqueo y la clave rotada.

### Gotcha 13 — 2 vCPU
- Nada de runners de CI; operaciones de a una.

### Gotcha 14 — Sin root
- El egreso (NAT) del diseño no se puede aplicar: el correo va por el reenviador de
  usuario; los cambios de host quedan para quien tenga root.

### Gotcha 15 — Firewall del host activo
- `pve-firewall` está activo: probar CT → host con el primer CT; si corta, agregar la
  regla desde el panel.

---

## 7. Rollback

- Apagar el reenviador del host si hiciera falta (matar solo ese `socat`).
- En emergencia de correo: apuntar el relay al **relay de casa** (plan B).
- El resto del stack no depende del correo; casa sigue andando (PRUEBAS).

---

## 8. Checklist del día

- [ ] IP autorizada en Brevo + clave rotada + bloqueo activo.
- [ ] Reenviador arriba y persistente (`10.10.1.1:2587`).
- [ ] Firewall probado (CT → host:2587).
- [ ] 6 LXC creados/restaurados en `vmbr1`.
- [ ] Base migrada y verificada; Redis con password.
- [ ] Backend `.env` (600) y systemd.
- [ ] Front servido y `/api` proxiado.
- [ ] App publicada al equipo y probada.
- [ ] Relay con `relayhost=[10.10.1.1]:2587`.
- [ ] Alta → `delivered`; recuperación → OTP.
- [ ] Rollback anotado.

---

## 9. Anexos

### Anexo A — Variables del backend (`.env` real, permisos `600`)

| Variable                  | Para qué                           |
|---------------------------|------------------------------------|
| `AUTH_MODE`               | Modo de autenticación (`postgres`) |
| `DATABASE_URL` / `DB_DSN` | Conexión a PostgreSQL              |
| `HTTP_ADDR`               | Puerto de la API (`:8080`)         |
| `JWT_SECRET`              | Firma de tokens                    |
| `TOTP_ENCRYPTION_KEY`     | Cifrado de los secretos 2FA        |
| `ENABLE_SWAGGER`          | Expone `/swagger`                  |
| `ALLOWED_ORIGINS`         | CORS                               |
| `PROXMOX_*`               | Integración con Proxmox            |
| `REDIS_*`                 | Conexión a Redis                   |
| `EMAIL_PROVIDER=smtp`     | Proveedor de correo real           |
| `SMTP_HOST`               | **El relay** (10.10.1.15)          |
| `SMTP_PORT`               | `25` (interno)                     |
| `SMTP_USER/PASS`          | Credenciales SMTP de Brevo         |
| `SMTP_FROM`               | Remitente                          |

### Anexo B — Config del relay

- `docker/relay/main.cf.tpl` → plantilla que se usa (adaptada al CT; **sin Docker**).
- `docker/relay/README.md` → el "por qué" de cada ajuste.

### Anexo C — Comandos clave

```bash
# Reenviador de correo en el host (usuario, sin root)
socat TCP-LISTEN:2587,bind=10.10.1.1,fork,reuseaddr TCP:smtp-relay.brevo.com:587

# Migrar la base (desde casa)
ssh root@proxmox 'pct exec 101 -- su - postgres -c "pg_dump -d centinela_test -Fc -f /tmp/centinela_test.dump"'
pct push 201 /ruta/centinela_test.dump /root/centinela_test.dump
pg_restore -d centinela_test --clean --if-exists /root/centinela_test.dump

# Publicar la app al equipo
socat TCP-LISTEN:8443,bind=100.116.49.111,fork,reuseaddr TCP:10.10.1.10:443
```

### Anexo D — Certificado del relay (STARTTLS interno)

```bash
openssl req -x509 -newkey rsa:2048 -nodes -keyout relay.key -out relay.crt -days 3650 \
  -subj "/CN=centinela-smtp-brevo" -addext "subjectAltName=IP:<IP-del-relay>"
cp relay.crt /usr/local/share/ca-certificates/centinela-smtp-brevo.crt
update-ca-certificates
```

### Anexo E — Evidencia de la verificación (06/10/2026, solo lectura)

| Chequeo          | Resultado                                   |
|------------------|---------------------------------------------|
| Usuario          | nicolas.lencina, sin root (sudo pide clave) |
| vmbr1 del host   | 10.10.1.1/24                                |
| Herramientas     | socat, python3, tmux, cron OK               |
| Firewall         | pve-firewall activo                         |
| IP de egreso     | 170.210.70.248 (estable en 3 muestras)      |
| Brevo 587        | alcanzable; cert regional sudamérica OK     |
| Brevo 25         | bloqueado (esperado)                        |
| Postfix del host | localhost:25 (notif. PVE; no tocar)         |

---

## 10. Diferencias respecto de casa (resumen)

| Tema              | Casa (hoy)               | C1 (objetivo)                             |
|-------------------|--------------------------|-------------------------------------------|
| Contenedores      | 6 LXC separados          | 6 LXC separados (igual)                   |
| Salida a internet | NAT propio               | reenviador de usuario (solo 587)          |
| IP de salida      | domiciliaria             | la del campus (170.210.70.248)            |
| Relay             | CT107 (`10.10.10.60`)    | CT205 (`10.10.1.15`)                      |
| Brevo IPs         | 1 autorizada (casa)      | + la IP del campus                        |
| Acceso            | Tailscale propia         | Tailscale del nodo + publicación por host |
| Deploy            | scripts                  | los mismos scripts (cambiando host/CT)    |
| CI                | runners en el hipervisor | **no** en el nodo (CPU)                   |

---

## 11. Repositorio y autorizaciones

- El repo principal es **`micaelaalvaradomendez/centinela`**: unifica front y back.
- **Los despliegues los autorizan únicamente Micaela y Nicolás.** Nada se publica,
  comenta, pushea ni despliega sin esa autorización.
