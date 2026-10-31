# centi-mica ( 106-CT106 )

infra@pve no tiene permisos para crear usuarios administradores ni elevar privilegios en Proxmox, entonces hice un Bastión en Tailscale:

cree un contenedor liviano tipo "jump host" con SSH unido a mi red de Tailscale para entrar directo por consola sin pasar por la interfaz de Proxmox.

# pasos 
cree el contenedor 106 por la API de Proxmox, conectándolo a las tres redes internas (vmbr1, vmbr2 y vmbr3) para que tenga alcance directo a todos los contenedores (bd, api, proxy, front, etc.) y salida a Internet vía vmbr1 (con el SNAT que tiene configurado).

en la consola web de ese contenedor en Proxmox puse el comando de instalación de Tailscale con SSH activado (tailscale up --ssh).

Una vez unido puedo conectarme directo por terminal a ese jump host y auditar cada servicio, archivo de configuración y base de datos sin hacer todo manual.