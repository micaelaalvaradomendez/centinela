# Comparativa Arquitectónica: Servidor de Prueba vs. Servidor Final

**Fecha:** 29 de Octubre de 2026  
**Proyecto:** El Centinela — Monitoreo de Infraestructura Proxmox VE  
**Servidor de Prueba:** `centinela.tail6bb3f3.ts.net` (`100.93.128.111`)  
**Servidor Final (Grupo 1):** `pve-is3-g1` (`100.116.49.111`) / Gateway `centi-mica` (`100.114.44.18`)  

---

## 1. Tabla Comparativa General

| Dimensión | Servidor de Prueba (`100.93.128.111`) | Servidor Final (`100.116.49.111`) | Diagnóstico / Impacto |
| :--- | :--- | :--- | :--- |
| **Exposición Web a Tailscale** | **Abierta**: Puertos 80 (HTTP) y 443 (HTTPS) responden directamente. | **Cerrada**: Solo responde el puerto 8006 (Proxmox VE). Puertos 80/443 cerrados. | **Crítico**: El servidor final no es accesible para usuarios finales sin el Jump Host `centi-mica`. |
| **Versión del Código (Backend)** | **Actualizada**: `commit 908522e` (29/10/2026 14:36). Rama `main`. | **Desactualizada / Dev**: `commit dev`, compilado localmente sin flags. | **Medio**: El backend del servidor final no tiene los últimos fixes de producción. |
| **Arquitectura de Ejecución** | Servidor Nginx integrado como frontend y reverse proxy hacia la API local. | **Multi-Contenedor (LXC)**: 6 contenedores independientes (`proxy`, `bd`, `api`, `front`, `redis`, `smtp`). | **Positivo conceptualmente**, pero con graves fallas de interconexión. |
| **Esquema de Redes Internas** | Red local única / comunicación directa sin segmentación estricta. | **3 Bridges virtuales**: `vmbr1` (App), `vmbr2` (Cliente/DMZ), `vmbr3` (Datos). | **Problema**: Segmentación a medias; Redis expuesto en `vmbr1`, y `vmbr2` sin salida. |
| **Bases de Datos** | Base de datos activa con usuarios de prueba cargados. | Dos bases creadas (`centinela_test` y `centinela_stable`). | **Pendiente**: Falta insertar los usuarios administradores en la base activa. |

---

## 2. Análisis Detallado: ¿Qué está mal en el Servidor Final y por qué?

### A. Fallas de Red y Enrutamiento (Capa de Infraestructura)

#### 1. Inaccesibilidad desde Tailscale hacia la Aplicación
* **El problema:** En el servidor final, el Host Proxmox tiene la IP de Tailscale `100.116.49.111`. Pero la aplicación vive adentro de los contenedores LXC (`10.10.1.10`, etc.). El Host no tiene configurado ningún reenvío de puertos (Port Forwarding / DNAT) en `iptables` desde `tailscale0` hacia `10.10.1.10:80`.
* **Consecuencia:** Nadie puede entrar a la URL del nodo en el navegador; solo se ve el panel de Proxmox en el puerto 8006.
* **Dónde y cómo corregir:**  
  En el host Proxmox (o manteniendo `centi-mica` como gateway activo con `tailscale serve`), ejecutar en el host:
  ```bash
  iptables -t nat -A PREROUTING -i tailscale0 -p tcp --dport 80 -j DNAT --to-destination 10.10.1.10:80
  iptables -t nat -A PREROUTING -i tailscale0 -p tcp --dport 443 -j DNAT --to-destination 10.10.1.10:443
  ```

#### 2. Mala ubicación de Redis (Seguridad de Capas)
* **El problema:** El contenedor `104 (centinela-redis)` fue conectado a la interfaz `vmbr1` con la IP `10.10.1.15`. En esa misma red convive el frontend (`103: 10.10.1.20`) y la pata interna del proxy (`100: 10.10.1.10`).
* **Consecuencia:** Redis almacena tokens de sesión JWT, bloqueos y datos de caché. En una arquitectura de 3 capas (DMZ, Aplicación, Datos), Redis debe estar en la capa de datos (`vmbr3`) accesible únicamente por el backend (`centinela-api`), nunca expuesto a la red del frontend.
* **Dónde y cómo corregir:**  
  En la configuración de red del contenedor `104` en Proxmox: cambiar el bridge de `vmbr1` a `vmbr3` (ej. IP `10.10.3.15/24`) y actualizar la variable `REDIS_ADDR=10.10.3.15:6379` en el `.env` de `centinela-api`.

#### 3. Salida a Internet del Microservicio de Correo (`centinela-smtp-brevo`)
* **El problema:** `centinela-smtp-brevo` (LXC 105) tiene credenciales para `smtp-relay.brevo.com`. Las redes `vmbr1`, `vmbr2` y `vmbr3` por diseño del laboratorio están aisladas.
* **Diagnóstico:** Si bien `vmbr1` tiene una zona SDN `labnat` con `snat: 1`, si el host no permite el reenvío de paquetes desde `vmbr1` hacia `vmbr0` (uplink físico a internet), el contenedor no podrá despachar correos reales de recuperación de clave.

---

### B. Fallas de Aplicación y Base de Datos (Capa de Software)

#### 1. Versión desactualizada de la API en el Servidor Final
* El servidor de prueba tiene el commit `908522e` (último deploy oficial de la rama `main`).
* El servidor final tiene una compilación vieja `dev`.
* **Dónde y cómo corregir:**  
  Compilar y desplegar la versión de `origin/main` en el contenedor `102` usando el script de despliegue oficial:
  ```bash
  ./backend/scripts/deploy-back.sh test
  ```

#### 2. Por qué falló el Login con "Credenciales inválidas"
* El contenedor `102 (centinela-api)` está conectado a la base de datos `centinela_test` (como se confirmó en su archivo `.env`):
  `DB_DSN=postgres://centinela_test_app:...@10.10.3.60:5432/centinela_test`
* Al revisar el contenedor `101 (centinela-bd)`, se consultó la tabla de organizaciones pero **aún no se ejecutó la inserción del usuario**.
* La API consulta la tabla `usuarios` buscando `micaelaalvaradomendez@gmail.com`. Al no existir la fila, responde correctamente con HTTP 401 `AUTH_FAILED: credenciales inválidas`.
* **Dónde y cómo corregir:**  
  Ejecutar la consulta de inserción (INSERT) en `centinela-bd` (LXC 101) para que el registro exista.

---

## 3. Plan de Remediación Paso a Paso

1. **Paso 1 (Inmediato - Login):** Ejecutar el `INSERT` en `101 (centinela-bd)` para crear el usuario admin y probar el flujo de login y 2FA.
2. **Paso 2 (Red):** Definir si el acceso permanente a la app será mediante el Jump Host `centi-mica` (con `tailscale serve`) o solicitando al administrador del laboratorio la regla de port forwarding en el host Proxmox.
3. **Paso 3 (Despliegue):** Actualizar el binario de `centinela-api` al último commit de `main` (`908522e`).
4. **Paso 4 (Aislamiento):** Mover `centinela-redis` a `vmbr3`.
5. **Paso 5 (Limpieza):** Detener y eliminar los contenedores residuales `110`, `210` y `310` para recuperar memoria RAM y disco.

