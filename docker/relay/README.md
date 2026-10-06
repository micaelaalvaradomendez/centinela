# smtp-relay-brevo — Repartidor de correo

Contenedor con **postfix** que recibe el correo del backend y lo **reenvía a Brevo**.

## Qué problema resuelve

Brevo solo acepta envíos desde **IPs autorizadas** (Settings → Security → Authorized IPs).
El backend salía a internet por la **IP pública de la red local, que cambia**. Cuando
cambia, Brevo rechaza el envío con:

```
525 "5.7.1 Unauthorized IP address"
```

y se rompe el **alta de usuarios** y la **recuperación de contraseña** (el correo con la
clave temporal y los códigos OTP no salía). Antes se "arreglaba" desactivando el bloqueo
de IPs en Brevo, lo que deja la clave SMTP usable desde cualquier lugar.

## Cómo lo resuelve

Todos los correos salen por **un único punto con IP fija** (este relay). Brevo ve siempre
esa IP, se autoriza **una sola** en el panel y se puede **reactivar** el bloqueo.

```
Backend (CT102) --Tailscale--> smtp-relay-brevo --Internet--> Brevo --> destinatario
```

## Qué NO cambia

- El backend sigue mandando SMTP igual; **solo cambia `SMTP_HOST`** para que apunte a
  este relay en vez de `smtp-relay.brevo.com`.
- La clave SMTP de Brevo vive en el `.env.relay` de este contenedor, no en el backend.

## Uso

```bash
cp .env.relay.example .env.relay   # completar BREVO_USER y BREVO_PASS
docker compose -f compose.relay.yaml up -d --build
```

En el servidor real, el puerto se publica **solo para la red de Tailscale** (no a internet).

## Detalles que parecen raros (y por qué están)

- **`dns: 1.1.1.1 / 8.8.8.8`**: el contenedor no puede usar el resolver del host
  (`127.0.0.53`, systemd-resolved) porque no es alcanzable desde adentro; sin esto,
  postfix no resuelve `smtp-relay.brevo.com` y el correo queda en cola.
- **`chroot = n`** (en `entrypoint.sh`): postfix corre enjaulado y adentro del jaque no ve
  `/etc/resolv.conf`, así que tampoco resuelve DNS. Se desactiva el chroot y se copia el
  `resolv.conf` al directorio del jaque.
- **`maillog_file = /dev/stdout`**: para ver los logs de postfix con `docker logs`.
