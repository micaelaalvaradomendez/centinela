# Informe de Auditoría de Infraestructura y Red — Servidor Final (Nodo G1: pve-is3-g1)

**Fecha:** 29 de Octubre de 2026  
**Auditor:** Jump Host `centi-mica` (LXC 106) / Tailscale  
**IP Nodo Proxmox Host:** `100.116.49.111` (`pve-is3-g1`)  
**IP Bastión Tailscale:** `100.114.44.18` (`centi-mica`)  
**URL de acceso directo a la App:** [http://centi-mica.tail878ccc.ts.net/](http://centi-mica.tail878ccc.ts.net/)  

---

## 1. Resumen de la Auditoría

Se desplegó y configuró con éxito el contenedor bastión **`centi-mica` (LXC 106)** dentro del nodo Proxmox, con interfaces en las 3 redes virtuales internas y conectado a Tailscale.

Gracias a este salto, se auditó la conectividad real y el estado de cada servicio desde adentro de la red del nodo, confirmando que **los contenedores y el código de Centinela están corriendo correctamente**, pero estaban completamente aislados del exterior debido a la arquitectura de red original de Proxmox.

---

## 2. Mapa y Estado Real de los Servicios

| VMID | Nombre | IPs Internas | Puertos Abiertos | Estado del Servicio |
| :--- | :--- | :--- | :--- | :--- |
| **100** | `centinela-proxy` | `10.10.1.10`, `10.10.2.10` | 80/TCP | **OK (HTTP 200)**. Nginx ruteando `/` al front y `/api/` a la API. |
| **101** | `centinela-bd` | `10.10.3.60` | 5432/TCP | **OK**. PostgreSQL 16 escuchando y aceptando conexiones. |
| **102** | `centinela-api` | `10.10.1.30`, `10.10.3.30` | 8080/TCP | **OK**. Backend Golang activo (`/api/version` responde `{"builtAt":"dev","commit":"dev"}`). |
| **103** | `centinela-front` | `10.10.1.20` | 80/TCP | **OK (HTTP 200)**. Nginx sirviendo la SPA de React/Vite. |
| **104** | `centinela-redis` | `10.10.1.15` | 6379/TCP | **OK**. Servidor Redis 7 activo y escuchando. |
| **105** | `centinela-smtp-brevo` | `10.10.1.60` | — | Activo en red `vmbr1`. |
| **106** | `centi-mica` | `10.10.1.250`, `10.10.2.250`, `10.10.3.250` | 22/TCP, 80/TCP | **Jump Host & Gateway Tailscale** activo. |
| **110** | `prueba-red` | `10.10.1.110`, `10.10.2.110` | — | Residual (Running innecesariamente). |
| **210 / 310** | `prueba-red-xxx` | — | — | Residuales (Stopped). |

---

## 3. Aclaración Clave: ¿Qué es `centinela-api`?

En el código fuente local (`compose.yaml`), el backend del sistema está definido como el servicio `backend`.
Quien implementó los contenedores en Proxmox decidió nombrarlo **`centinela-api` (LXC 102)** porque ejecuta el **binario compilado de Golang** que atiende la API REST en el puerto `8080`. 

Ambos términos (`backend` y `api`) refieren exactamente al mismo componente en este proyecto.

---

## 4. Diagnóstico de la Red: ¿Por qué estaba mal?

1. **Inaccesibilidad desde Tailscale (Resuelto vía `centi-mica`):**  
   El Proxmox original solo exponía el puerto `8006` en su IP de Tailscale `100.116.49.111`. Los contenedores de Centinela estaban encerrados en redes privadas virtuales sin ninguna regla de reenvío (DNAT).  
   *Solución aplicada:* A través de `centi-mica` se habilitó el reenvío de tráfico Tailscale (`tailscale serve`) hacia el proxy interno `10.10.1.10:80`.
2. **Ubicación de Redis:**  
   `centinela-redis` (`10.10.1.15`) fue ubicado en la red `vmbr1` (donde convive con el front y el proxy) en vez de estar en `vmbr3` (red interna de datos con la BD). Como Redis almacena tokens de sesión sensibles, debería aislarse de la red del frontend.
3. **Salida a Internet para Brevo:**  
   `vmbr1` tiene Source NAT (`snat: 1`) hacia el exterior mediante la zona SDN `labnat`, lo que permite salida a internet (comprobado porque pudimos instalar paquetes y Tailscale en `centi-mica`).
4. **Desperdicio de Recursos:**  
   El contenedor `110 (prueba-red)` continúa consumiendo 256 MB de RAM de los escasos 4.1 GB del nodo y debería apagarse.

---

## 5. Acceso Inmediato

Para abrir la aplicación desplegada en el servidor final directamente en tu navegador:
👉 **[http://centi-mica.tail878ccc.ts.net/](http://centi-mica.tail878ccc.ts.net/)** (o `http://100.114.44.18/`)


La aplicación es inaccesible desde afuera: En el servidor de pruebas (100.93.128.111), Nginx responde directo en el puerto 80/443. En este servidor final (100.116.49.111), esos puertos están cerrados. Los contenedores están encerrados en redes virtuales aisladas (10.10.x.x) y no hay ningún port forwarding (DNAT) en el host para que Tailscale llegue al proxy.
Mezcla de capas de seguridad: Supuestamente armaron 3 redes (vmbr1 para lógica, vmbr2 para clientes, vmbr3 para datos). Pero pusieron a Redis en vmbr1 compartiendo red con el front y el proxy, en lugar de aislarlo en la red de datos junto a la base de datos.
Brevo sin salida a Internet: centinela-smtp-brevo necesita conectarse a los servidores externos de Brevo, pero las redes internas de Proxmox no tienen gateway ni salida a Internet hacia vmbr0.
Desorden de configuración: Interfaces llamadas eth3 en vez de eth0, gateways ficticios y contenedores de prueba ocupando recursos en un nodo de apenas 4 GB de RAM.