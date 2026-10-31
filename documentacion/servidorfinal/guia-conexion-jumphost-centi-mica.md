# Guía de Conexión al Jump Host `centi-mica`

Esta guía detalla cómo acceder y utilizar el contenedor bastión (**Jump Host**) `centi-mica` desplegado en el nodo Proxmox del Grupo 1 (`pve-is3-g1`).

Permite a cualquier integrante del equipo o revisor:
1. **Ver y probar la aplicación web** directamente en el navegador.
2. **Acceder a la terminal del Jump Host** para diagnóstico y auditoría.
3. **Saltar hacia los contenedores internos** (Proxy, Front, API, Redis, BD) que se encuentran aislados en las redes virtuales privadas.

---

## 1. Ficha Técnica del Bastión

| Parámetro | Valor |
| :--- | :--- |
| **Identificador Proxmox** | LXC `106` (`centi-mica`) en nodo `pve-is3-g1` |
| **IP Tailscale** | `100.114.44.18` |
| **FQDN MagicDNS** | `centi-mica.tail878ccc.ts.net` |
| **Usuario SSH** | `root` |
| **Servicios activos** | OpenSSH (`22`), Tailscale SSH, `tailscale serve` (HTTP `80` $\to$ `10.10.1.10:80`) |
| **IPs en redes internas** | `10.10.1.250` (`vmbr1`), `10.10.2.250` (`vmbr2`), `10.10.3.250` (`vmbr3`) |

---

## 2. Acceso Web a la Aplicación (Navegador)

El Jump Host tiene configurado un reenvío mediante Tailscale Serve que publica el Nginx del sistema (`centinela-proxy` en `10.10.1.10:80`).

### Si estás conectado a la misma Tailnet o el dispositivo te fue compartido:
Ingresar directamente a cualquiera de estas dos URLs:
- **[http://centi-mica.tail878ccc.ts.net/](http://centi-mica.tail878ccc.ts.net/)**
- **[http://100.114.44.18/](http://100.114.44.18/)**

> **Nota:** La aplicación responde el frontend React y enruta `/api/` hacia el backend Golang.

---

## 3. Acceso por Terminal (SSH) al Jump Host

Existen dos vías de acceso según los permisos y herramientas disponibles:

### Opción A: Vía Tailscale SSH (Recomendada)

Si tienes acceso al nodo dentro de Tailscale:

```bash
# Método directo con el cliente de Tailscale:
tailscale ssh root@centi-mica

# O mediante SSH estándar apuntando a la IP de Tailscale:
ssh root@100.114.44.18
```

> **Compartir el acceso con otra cuenta de Tailscale:**
> Desde la consola de administración de Tailscale ([login.tailscale.com/admin/machines](https://login.tailscale.com/admin/machines)):
> 1. Ubicar la máquina `centi-mica`.
> 2. Seleccionar el menú de opciones (`...`) $\to$ **Share...**.
> 3. Generar el enlace de invitación (Device Sharing) y enviárselo a la otra persona para que acepte el nodo en su tailnet.

---

### Opción B: Vía Consola Web de Proxmox VE

Si la persona cuenta con las credenciales de gestión del nodo Proxmox:

1. Abrir la interfaz web de Proxmox (con Tailscale activo):
   - **URL:** `https://100.116.49.111:8006/`
   - **Usuario:** `infra`
   - **Realm:** `Proxmox VE authentication server`
   - **Contraseña:** `Infra2026.`
2. En el árbol izquierdo, desplegar el nodo `pve-is3-g1`.
3. Seleccionar el contenedor **`106 (centi-mica)`**.
4. Hacer clic en **`>_ Consola`** en la barra superior. Se abrirá una sesión interactiva root de inmediato sin requerir claves SSH.

---

## 4. Uso como Jump Host: Conexión a Contenedores Internos

Debido a que `centi-mica` tiene interfaces en `vmbr1`, `vmbr2` y `vmbr3`, actúa como puente hacia los demás contenedores del clúster:

### Mapa de IPs Internas

| Contenedor | VMID | IP Interna | Rol |
| :--- | :--- | :--- | :--- |
| `centinela-proxy` | 100 | `10.10.1.10` | Nginx Reverse Proxy (Frontend + API) |
| `centinela-front` | 103 | `10.10.1.20` | Frontend React / SPA |
| `centinela-api` | 102 | `10.10.1.30` | Backend API REST (Golang, puerto 8080) |
| `centinela-redis` | 104 | `10.10.1.15` | Cache / Sesiones Redis (puerto 6379) |
| `centinela-bd` | 101 | `10.10.3.60` | Base de datos PostgreSQL 16 (puerto 5432) |

---

### Salto Directo desde la Consola de `centi-mica`

Una vez conectado dentro de `centi-mica`, puedes interactuar o saltar directamente a cualquier componente:

```bash
# Probar respuesta del Proxy HTTP
curl -I http://10.10.1.10/

# Probar la API de Golang directamente
curl http://10.10.1.30:8080/api/version

# Conectarse por SSH al contenedor del Backend o Proxy
ssh root@10.10.1.30
ssh root@10.10.1.10

# Consultar PostgreSQL directamente
psql -h 10.10.3.60 -U centinela -d centinela
```

---

### Salto Transparente desde tu Máquina Local (`ProxyJump`)

Puedes configurar tu archivo local `~/.ssh/config` para conectarte a los contenedores internos en un solo comando:

```ssh-config
# Bastión Jump Host
Host centi-mica
    HostName 100.114.44.18
    User root

# Contenedores internos a través del Jump Host
Host centinela-*
    ProxyJump centi-mica
    User root

Host centinela-proxy
    HostName 10.10.1.10

Host centinela-api
    HostName 10.10.1.30

Host centinela-bd
    HostName 10.10.3.60
```

Con esa configuración, puedes acceder directamente con:
```bash
ssh centinela-api
```

---

## 5. Verificación de Estado y Diagnóstico en `centi-mica`

Si necesitas comprobar que el enrutamiento y Tailscale siguen operando:

```bash
# Ver estado del proxy reverso Tailscale Serve
tailscale serve status

# Ver conectividad de interfaces de red
ip -br a

# Comprobar estado del nodo en Tailscale
tailscale status
```

