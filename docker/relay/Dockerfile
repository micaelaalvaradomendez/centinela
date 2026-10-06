# ============================================================================
#  Imagen del relay  smtp-relay-brevo
#  Recibe SMTP y lo reenvía a Brevo usando postfix.
#  Las credenciales NO van en la imagen: se inyectan por variables de entorno.
# ============================================================================
FROM debian:bookworm-slim

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
      postfix libsasl2-modules ca-certificates gettext-base \
 && rm -rf /var/lib/apt/lists/*

COPY main.cf.tpl /etc/postfix/main.cf.tpl
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod 644 /etc/postfix/main.cf.tpl && chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 25
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
