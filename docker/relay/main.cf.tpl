# ============================================================================
#  Configuración de postfix para el relay  smtp-relay-brevo
# ----------------------------------------------------------------------------
#  Reenvía el correo del backend hacia Brevo.
#  Las líneas marcadas con WHY explican el motivo de cada ajuste, para que no
#  parezcan "líneas mágicas" cuando las lea otra persona del equipo.
# ============================================================================

# Identidad del relay (cómo se presenta al saludar).
myhostname = ${RELAY_MYHOSTNAME}

# No somos el destino final de ningún dominio: solo reenviamos.
mydestination =

# Escuchamos en todas las interfaces del contenedor; el acceso se restringe
# con mynetworks (abajo) y con cómo se publica el puerto.
inet_interfaces = all

# WHY: solo aceptamos correo desde redes internas, nunca desde internet.
#      En el server real se limita a la red de Tailscale.
mynetworks = ${RELAY_MYNETWORKS}

# WHY: si el origen no está en mynetworks, rechazamos (no somos "open relay").
smtpd_relay_restrictions = permit_mynetworks, reject

# WHY: acá es donde reenviamos: al servidor SMTP de Brevo.
relayhost = [${BREVO_HOST}]:${BREVO_PORT}

# WHY: nos autenticamos con la clave SMTP de Brevo (usuario + clave).
smtp_sasl_auth_enable = yes
smtp_sasl_password_maps = static:${BREVO_USER}:${BREVO_PASS}
smtp_sasl_security_options = noanonymous

# WHY: ciframos la conexión hacia Brevo (STARTTLS).
smtp_tls_security_level = encrypt
smtp_tls_CAfile = /etc/ssl/certs/ca-certificates.crt

# WHY: mandamos el log de postfix a la salida estándar para poder verlo con
#      "docker logs" (en el contenedor no hay syslog).
maillog_file = /dev/stdout
