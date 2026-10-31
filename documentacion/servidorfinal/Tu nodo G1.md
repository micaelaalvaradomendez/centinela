# Tu nodo G1 — Guía de referencia

Esta guía describe el nodo del **Grupo 1** y qué puede hacer el rol de infraestructura dentro de
él. Léala **después** de completar la sección de su sistema operativo dentro de `Guía de conexión - Nodo G1.docx`: **Sección 2 — Windows**, **Sección 3 — Linux / Ubuntu** o **Sección 4 — macOS**, según su equipo.

---

## 1. Qué es este nodo

Es un nodo **Proxmox VE independiente del Grupo 1** (`pve-is3-g1`): el sandbox del grupo para el
proyecto de Ingeniería de Software 3. Cada grupo tiene el suyo, y los nodos **no forman un
clúster** entre sí: lo que ocurre en este nodo no afecta a los demás.

Usted entra a su administración por la **Web UI** y crea ahí sus máquinas virtuales y
contenedores.

---

## 2. Recursos disponibles

| Recurso           | Valor            |
|-------------------|------------------|
| Nodo              | `pve-is3-g1`     |
| IP LAN (gestión)  | `10.7.93.124`    |
| IP de tailnet     | `100.116.49.111` |
| CPU               | 2 vCPU           |
| Memoria           | ~4,1 GB de RAM   |
| Disco `local`     | 37,6 GB          |
| Disco `local-lvm` | 48,6 GB          |
| Disco total       | ~86 GB           |
| Proxmox VE        | 8.4.21           |

**Detalle:** `local` y `local-lvm` son los dos almacenamientos disponibles (discos, plantillas e
imágenes ISO). Los recursos son acotados: son un entorno educativo, no un servidor de producción.

---

## 3. Cómo entrar

Primero confirme que **Tailscale esté encendido y conectado**. Después abra la Web UI:

```
URL:         https://100.116.49.111:8006/
Usuario:     infra
Realm:       Proxmox VE authentication server
Contraseña:  Infra2026.
```

**Enlace de invitación (Device Sharing):** `https://login.tailscale.com/admin/invite/7hgHA4tQXPYjS22vFJx421`

**Vigencia del enlace:** 30 días desde el **27/10/2026**, hasta el **26/11/2026**. Si vence, pida uno nuevo.

> El navegador va a advertir por el **certificado**. Es esperado: el nodo usa un certificado
> propio, así que continúe.

**Detalle:** la cuenta `infra@pve` es la del **rol de infraestructura** del Grupo 1. Úsela solo
dentro del grupo.

---

## 4. Qué puede hacer el rol de infraestructura

| Capacidad                            | ¿Permitido? |
|--------------------------------------|-------------|
| Ver todo el nodo (solo lectura)      | Sí          |
| Crear y administrar VM y LXC         | Sí          |
| Encender, apagar y reiniciar         | Sí          |
| Clonar, respaldar y borrar           | Sí          |
| Usar `local` y `local-lvm`           | Sí          |
| Conectar a `vmbr1`, `vmbr2`, `vmbr3` | Sí          |
| Crear o cambiar usuarios             | No          |
| Cambiar permisos                     | No          |
| Conectar algo a `vmbr0` (campus)     | No          |
| Consola o shell del nodo             | No          |
| Cambiar la configuración del sistema | No          |
| Crear otras redes                    | No          |

**Detalle:** "no puede" significa que el nodo **deniega** la operación (un error de permiso, sin
importar lo que muestre la interfaz). En particular, la interfaz puede dejarle abrir un diálogo
que después falla: que un diálogo se abra no significa que la operación esté permitida.

---

## 5. Redes internas del nodo

| Red     | Rango          |
|---------|----------------|
| `vmbr1` | `10.10.1.0/24` |
| `vmbr2` | `10.10.2.0/24` |
| `vmbr3` | `10.10.3.0/24` |

**Detalle:** son **redes internas aisladas**: no tienen salida a internet, no tienen acceso al
campus y el nodo **no routea** entre ellas. Sirven para practicar comunicación entre sus propias
máquinas y contenedores dentro de cada red. La red del campus es `vmbr0`, y con este rol **no se
puede** conectar nada a ella.

**Host contra máquinas:** el **host** (el nodo Proxmox) **sí** sale a la red del campus y a
internet por su uplink físico en `vmbr0`, usando su IP LAN `10.7.93.124` y la puerta de enlace del
campus `10.7.93.1`; por esa salida corre Tailscale. Sus **máquinas**, en cambio, están en las redes
internas **sin uplink, sin gateway ni NAT**, así que **no tienen salida** al campus ni a internet.
Darles salida más adelante requiere un diseño de routing/NAT administrado; no está habilitado y no
deben conectar sus máquinas a `vmbr0`.

---

## 6. Reglas de convivencia

- No toque **usuarios ni permisos**: no son parte de su rol.
- No conecte nada a `vmbr0` (la red del campus): está fuera del alcance del grupo.
- Ante un error o una duda, **deténgase y consulte** antes de forzar nada.
- No comparta la contraseña ni el enlace de invitación fuera del grupo.
- Tailscale tiene que estar **siempre encendido** para trabajar en el nodo.

---

## 7. Consultas

Ante cualquier duda, consulte con el **referente de infraestructura de su grupo** o con el
**administrador del laboratorio**.

---

**Última revisión:** 2026-09-29.
