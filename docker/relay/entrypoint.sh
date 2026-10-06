#!/bin/sh
# ============================================================================
#  Punto de entrada del relay  smtp-relay-brevo
# ----------------------------------------------------------------------------
#  Arma la configuración de postfix desde variables de entorno (así la clave
#  SMTP nunca queda escrita en la imagen ni en el repositorio) y arranca
#  postfix en primer plano.
# ============================================================================
set -e

: "${BREVO_HOST:=smtp-relay.brevo.com}"
: "${BREVO_PORT:=587}"
: "${RELAY_MYHOSTNAME:=relay.centinela.local}"
: "${RELAY_MYNETWORKS:=127.0.0.0/8 10.0.0.0/8 172.16.0.0/12 192.168.0.0/16}"

if [ -z "${BREVO_USER:-}" ] || [ -z "${BREVO_PASS:-}" ]; then
  echo "[smtp-relay-brevo] ERROR: faltan BREVO_USER y/o BREVO_PASS" >&2
  exit 1
fi

export BREVO_HOST BREVO_PORT BREVO_USER BREVO_PASS RELAY_MYHOSTNAME RELAY_MYNETWORKS
envsubst < /etc/postfix/main.cf.tpl > /etc/postfix/main.cf

# WHY: postfix corre enjaulado (chroot) por defecto y, adentro del jaque, no
#      ve /etc/resolv.conf, así que NO resuelve smtp-relay.brevo.com y el
#      correo queda en cola ("Host or domain name not found"). Lo desactivamos
#      y, por las dudas, copiamos el resolv.conf al directorio del jaque.
postconf -F '*/*/chroot = n' >/dev/null 2>&1 || true
mkdir -p /var/spool/postfix/etc
cp -f /etc/resolv.conf /var/spool/postfix/etc/resolv.conf 2>/dev/null || true

echo "[smtp-relay-brevo] relayhost=[$BREVO_HOST]:$BREVO_PORT | myhostname=$RELAY_MYHOSTNAME | mynetworks=$RELAY_MYNETWORKS"
exec postfix start-fg
