# Integración y Alcance de la API de Proxmox VE

Este documento define la configuración, el alcance operativo y las directivas de seguridad para el consumo de la API REST de Proxmox VE desde el backend de **El Centinela**, garantizando la estabilidad del servidor hipervisor y aislando los servicios productivos del entorno.

---

## 1. Configuración de Entorno (.env)

Para que el backend interactúe con la API de Proxmox, debe utilizar las siguientes variables de entorno:

```env
# Conexión Proxmox VE
PROXMOX_URL=https://100.81.49.19:8006/api2/json
PROXMOX_NODE=proxmox
PROXMOX_TOKEN_ID=centi-api@pve!backend-token
PROXMOX_TOKEN_SECRET=39e714a5-967a-46c2-8de9-af048ae0020f
PROXMOX_INSECURE_SKIP_VERIFY=true
```

### Cabecera de Autenticación
Cada petición HTTP enviada a la API debe incluir la cabecera estándar de API Token de Proxmox:

```http
Authorization: PVEAPIToken=centi-api@pve!backend-token=39e714a5-967a-46c2-8de9-af048ae0020f
```

> **Requisito de Red:** La IP `100.81.49.19` pertenece a la red privada Tailscale. La máquina o contenedor que ejecute el backend debe estar autenticado en la red para alcanzar el puerto `8006`.

---

## 2. Alcance Operativo (Scope)

El token `centi-api@pve!backend-token` tiene alcance de administración sobre el cluster, lo que permite implementar los requerimientos funcionales de la plataforma:

| Requerimiento | Operación API | Descripción |
|---|---|---|
| **RF-02: Dashboard de estado** | `GET /nodes/{node}/status` | Lectura de telemetría del host (CPU, memoria, disco, uptime). |
| **RF-03: Inventario de instancias** | `GET /nodes/{node}/qemu`<br>`GET /nodes/{node}/lxc` | Listado y métricas de VMs y contenedores LXC. |
| **RF-04: Control de instancias** | `POST /nodes/{node}/{type}/{vmid}/status/{action}` | Iniciar, apagar, forzar apagado o reiniciar instancias asignadas. |
| **RF-07: Asistente de creación** | `POST /nodes/{node}/qemu`<br>`POST /nodes/{node}/lxc` | Aprovisionamiento de nuevas instancias con recursos acotados. |
| **Eliminación y limpieza** | `DELETE /nodes/{node}/{type}/{vmid}` | Destrucción de instancias y purga de sus volúmenes asociados. |
| **Asignación de VMID** | `GET /cluster/nextid` | Reserva automática del siguiente identificador libre. |

---

## 3. Directivas de Protección: Operar sin modificar el Servidor

Para no alterar la infraestructura base ni interrumpir los servicios del servidor, el backend debe respetar estrictamente las siguientes reglas:

### 3.1. Rango de VMIDs Protegidos (NO TOCAR)
Los siguientes IDs corresponden a la infraestructura interna de Centinela desplegada en el servidor:

| VMID | Nombre | Rol en la infraestructura |
|---|---|---|
| **100** | `centinela-proxy` | Nginx reverse proxy y certificados TLS |
| **101** | `centinela-db` | Base de datos PostgreSQL (`centinela_test`) |
| **102** | `centinela-api-test` | API Go del entorno de pruebas |
| **103** | `centinela-front-test` | Frontend del entorno de pruebas |
| **104** | `centinela-api-stable` | API Go del entorno estable |
| **105** | `centinela-front-stable` | Frontend del entorno estable |

> ⛔ **Regla de Aislamiento:** El backend **bajo ninguna circunstancia** debe enviar peticiones de modificación, apagado o eliminación (`DELETE` o `POST status/stop`) dirigidas a los IDs `100`, `101`, `102`, `103`, `104` o `105`. Toda creación dinámica debe usar IDs asignados por `/cluster/nextid` (rango >= 106).

### 3.2. Prohibición de Modificar el Host Físico
- **No alterar configuraciones del nodo:** Queda prohibido llamar a endpoints destructivos sobre el host (`/nodes/proxmox/network`, `/nodes/proxmox/dns`, `/nodes/proxmox/disks`, `/nodes/proxmox/certificates`).
- **No reiniciar el host:** No invocar `/nodes/proxmox/status` con acciones de `reboot` o `shutdown`.
- **Almacenamiento permitido:** Solo utilizar `local-lvm` para discos virtuales de instancias y `local` para lectura de plantillas ISO (`local:iso/...`) o LXC (`local:vztmpl/...`). No crear, borrar ni redimensionar pools de almacenamiento.

---

## 4. Gestión Asíncrona de Tareas (UPID)

Proxmox ejecuta operaciones de creación, eliminación y ciclo de vida de forma **asíncrona**. En lugar de esperar a que la acción termine, la API responde de inmediato con un identificador de tarea (`UPID`):

```json
{
  "data": "UPID:proxmox:0014323F:064AB169:6AB943C6:qmcreate:110:centi-api@pve!backend-token:"
}
```

### Ciclo de vida de una orden en el Backend:
1. **Enviar la orden:** El backend valida permisos de negocio y emite la llamada (ej. `qmcreate` o `vmstart`).
2. **Capturar UPID:** Almacena el `UPID` retornado y devuelve al cliente `HTTP 202 Accepted`.
3. **Monitoreo en segundo plano:** El worker consulta periódicamente:
   ```http
   GET /api2/json/nodes/proxmox/tasks/{UPID}/status
   ```
4. **Verificación de finalización:** La tarea finaliza cuando `status` es `"stopped"`:
   - Si `exitstatus == "OK"`: la operación concluyó exitosamente.
   - Si `exitstatus != "OK"`: la tarea falló (capturar el error en la tabla de auditoría).

---

## 5. Ejemplos Prácticos de Integración

### 5.1. Consulta de Salud del Host (RF-02)
```bash
curl -k -s -H "Authorization: PVEAPIToken=centi-api@pve!backend-token=39e714a5-967a-46c2-8de9-af048ae0020f" \
  https://100.81.49.19:8006/api2/json/nodes/proxmox/status
```

### 5.2. Obtención del siguiente VMID libre
```bash
curl -k -s -H "Authorization: PVEAPIToken=centi-api@pve!backend-token=39e714a5-967a-46c2-8de9-af048ae0020f" \
  https://100.81.49.19:8006/api2/json/cluster/nextid
```

### 5.3. Creación de una Máquina Virtual (QEMU)
```bash
curl -k -s -X POST \
  -H "Authorization: PVEAPIToken=centi-api@pve!backend-token=39e714a5-967a-46c2-8de9-af048ae0020f" \
  -H "Content-Type: application/json" \
  -d '{
    "vmid": 110,
    "name": "vm-demo",
    "memory": 2048,
    "cores": 2,
    "sockets": 1,
    "net0": "virtio,bridge=vmbr0",
    "scsihw": "virtio-scsi-pci",
    "scsi0": "local-lvm:20",
    "ide2": "local:iso/debian-13.6.0-amd64-netinst.iso,media=cdrom",
    "ostype": "l26"
  }' \
  https://100.81.49.19:8006/api2/json/nodes/proxmox/qemu
```

### 5.4. Destrucción Limpia de una Instancia de Prueba
```bash
curl -k -s -X DELETE \
  -H "Authorization: PVEAPIToken=centi-api@pve!backend-token=39e714a5-967a-46c2-8de9-af048ae0020f" \
  "https://100.81.49.19:8006/api2/json/nodes/proxmox/qemu/110?purge=1&destroy-unreferenced-disks=1"
```

---

## 6. Resumen de Seguridad para Desarrolladores

1. **Nunca exponer el secreto en repositorios:** `PROXMOX_TOKEN_SECRET` debe existir únicamente en el archivo `.env` local (ignorado por `.gitignore`).
2. **Validar antes de invocar:** Comprobar memoria y espacio en disco libres antes de mandar órdenes de creación para evitar saturar el almacenamiento de Proxmox.
3. **Control por UPID:** Toda operación de cambio de estado debe monitorear su respectivo UPID antes de actualizar el estado visual en la interfaz de usuario.
