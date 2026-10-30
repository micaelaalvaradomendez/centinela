# SMTP y Brevo — egreso del relay

> **Estado**: resuelto y verificado el 30/10/2026 en el entorno de casa.
> Sin credenciales. El relay es `centinela-smtp-brevo` (CT107 en casa, CT105 en el nodo de la
> universidad).

## Síntoma

Los correos quedaban en cola:

```
delivery temporarily suspended: connect to smtp-relay.brevo.com[1.179.117.1]:587:
Connection timed out
```

El alta de usuarios y la recuperación de contraseña dependen de ese correo.

## Causa

Brevo solo acepta envíos desde IPs autorizadas. El relay salía por la IP pública del host
(`179.62.61.127`), que estaba autorizada, pero el firewall del host solo dejaba salir al proxy:
las reglas de 587/465 que habían funcionado quedaron en un snapshot viejo
(`/etc/iptables/rules.v4`) y no en el ruleset vivo ni en `/etc/network/interfaces`.

## Fix

Acotar el egreso al relay, sin abrirlo a todos los contenedores:

```bash
iptables -t nat -A POSTROUTING -s <IP-RELAY>/32 -o <UPLINK> -p tcp -m multiport --dports 587,465 -j MASQUERADE
iptables -A FORWARD -i <BRIDGE> -o <UPLINK> -s <IP-RELAY>/32 -p tcp -m multiport --dports 587,465 -j ACCEPT
```

El `ACCEPT` va **antes** del `DROP` general. En casa quedó persistido en
`/etc/network/interfaces`, con sus `post-down`.

## Verificación

- `getent hosts smtp-relay.brevo.com` resuelve.
- TCP 587 conecta; STARTTLS negocia **TLS 1.3** con
  `smtp-relay-southamerica-east1-v2.brevo.com`.
- `postqueue -f` vació la cola: 8 mensajes con `250 2.0.0 OK`.
- Prueba nueva a un destinatario real: `status=sent` en 2,3 s.
- La matriz de red quedó igual: solo el proxy sale a internet; el relay solo puede 587/465.

## En el nodo de la universidad

El relay vive en `vmbr1`, la red con SNAT, así que sale con la IP fija del campus
(`170.210.70.248`): se autoriza esa IP en Brevo y no hace falta Tailscale.

Si el 587 está bloqueado, probar **465**. Si todo el SMTP saliente está bloqueado, relay en un
VPS por **443** o API HTTP de Brevo. La matriz de prueba completa vive en el documento de
réplica de la materia, fuera de este repo.
