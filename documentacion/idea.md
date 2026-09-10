# El Centinela

## Orquestación y gestión de entornos virtualizados y contenedores

## Resumen

### Problema

Administrar un servidor donde funcionan máquinas virtuales y contenedores es una tarea compleja y delicada. No todas las personas deben tener permisos para gestionar un servidor y, cuando alguien recién empieza, es fácil perderse entre las interfaces o cometer errores, como apagar el servidor o detener todo lo que está funcionando en él.

### Contexto

Para gestionar servidores se utilizan hipervisores como Proxmox VE. Estas plataformas funcionan como un administrador de tareas para servidores: permiten comprobar si un servicio falló, reiniciarlo o desplegar nuevos recursos en producción.

### Solución

El Centinela será un panel de control que consumirá la API de Proxmox VE para mostrar información en tiempo real y ejecutar operaciones autorizadas sobre la infraestructura. Tendrá una interfaz clara y simple, con permisos controlados, para que una persona pueda consultar el estado de los servicios sin poner en riesgo el resto del servidor.

El panel incluirá:

- **Dashboard de estado:** monitoreo en tiempo real del servidor, con el consumo de RAM, CPU y disco de cada máquina virtual o contenedor.
- **Botones de control rápido:** opciones directas para encender, apagar y reiniciar instancias sin ingresar a la interfaz nativa de Proxmox VE.
- **Despliegue simplificado (wizard):** formulario de tres pasos para crear una máquina virtual o un contenedor sin configurar archivos manualmente.
- **Aislamiento por usuarios y roles:** sistema de permisos mediante el cual un administrador asigna qué instancias puede gestionar cada usuario y bloquea acciones peligrosas.

## División de trabajo

- **Infraestructura:** preparar el servidor host con Proxmox VE y Linux de prueba, configurar las máquinas virtuales y contenedores base, establecer las redes virtuales y dejar listas las credenciales y los tokens de API para que el backend se conecte.
- **Backend:** desarrollar la API, gestionar la autenticación, coordinar la comunicación asíncrona con Proxmox VE y levantar el servidor de WebSockets para las métricas en vivo.
- **Frontend:** diseñar el dashboard, los gráficos interactivos alimentados mediante WebSockets y los formularios de creación.

## Desarrollo

Centinela consiste en desarrollar un sistema web de orquestación y gestión de entornos virtualizados y contenedores, con el propósito de simplificar la administración, el monitoreo y el aprovisionamiento de infraestructura sobre servidores Linux en tiempo real.

El objetivo principal es construir una plataforma web de orquestación y abstracción para hipervisores y clústeres, como Proxmox VE o VMware. La plataforma permitirá a los administradores de infraestructura desplegar, controlar y monitorear máquinas virtuales (VMs) o contenedores de manera remota, segura y automatizada, mediante un entorno visual más limpio e intuitivo.

En lugar de construir un hipervisor desde cero a nivel del kernel, el proyecto consiste en desarrollar un panel de control unificado, personalizado y simplificado que interactúe directamente con las API de las plataformas de virtualización subyacentes.

## Problemática que soluciona

- **Curva de aprendizaje y complejidad visual:** las interfaces nativas de plataformas como Proxmox VE o vSphere están pensadas para administradores de infraestructura sénior. Presentan una alta densidad de opciones avanzadas de red, almacenamiento en clúster y parámetros del kernel que pueden resultar abrumadoras para desarrolladores, personal de soporte o clientes finales.
- **Riesgo operativo por exceso de permisos:** dar acceso directo a la consola principal de Proxmox VE a usuarios que solo necesitan apagar, encender o monitorear sus máquinas virtuales o contenedores (LXC) aumenta el riesgo de configuraciones erróneas a nivel del nodo o del servidor físico.
- **Falta de personalización y experiencia de usuario modernizada:** los paneles tradicionales suelen priorizar la funcionalidad sobre la usabilidad y ofrecen interfaces densas, poco enfocadas en flujos de trabajo ágiles o métricas claras de uso diario.
- **Aislamiento simplificado:** ofrecer a cada equipo o cliente un entorno donde solo vea sus recursos asignados requiere una capa propia de aislamiento, con un panel despejado y acciones directas sobre sus instancias.

## ¿Qué implica?

El desarrollo de este sistema exige integrar comunicaciones en tiempo real, conectividad con interfaces del sistema y protocolos de gestión remota:

- **Integración con la API del hipervisor:** consumir los endpoints de Proxmox VE de forma segura mediante tokens de API para consultar el inventario, disparar operaciones y leer estados de salud.
- **Manejo de estados y operaciones asíncronas:** las tareas de infraestructura, como crear una VM, tomar un snapshot o migrar un contenedor, llevan tiempo y no deben bloquear el uso del sistema. Se gestionarán mediante colas de tareas y confirmación de eventos.
- **Streaming de métricas:** implementar la infraestructura necesaria para canalizar en tiempo real la telemetría de CPU, RAM, red y disco.
- **Motor de reglas y permisos internos:** desarrollar un esquema propio de control de acceso que permita asociar usuarios de la plataforma con recursos específicos del hipervisor sin exponer credenciales administrativas.

## Alcance

El alcance de El Centinela se limita estrictamente al aprovisionamiento, la configuración de hardware virtual (vCPU, RAM, disco y red) y la gestión del ciclo de vida (encendido, apagado, reinicio y snapshots) de entornos virtualizados (VMs y LXC) sobre Proxmox VE.

El despliegue de código fuente, la subida de archivos de aplicaciones finales y la configuración de entornos de ejecución dentro de los contenedores o máquinas virtuales creadas quedan expresamente fuera del alcance de esta plataforma. Esas tareas de integración y despliegue continuo (CI/CD) corresponden a herramientas externas, como GitHub Actions o Jenkins, o a la gestión directa del desarrollador sobre la instancia.

## Funcionalidades mínimas

### 1. Tablero general del servidor o nodo

- **Estado de salud global:** visión resumida del consumo total del servidor físico, incluyendo CPU, memoria RAM, espacio en disco y tiempo de actividad.
- **Inventario de instancias:** lista clara y filtrable que distinga entre máquinas virtuales (VMs) y contenedores (LXC), mostrando su estado: `Running`, `Stopped` o `Paused`.

### 2. Gestión del ciclo de vida de instancias

- **Controles de energía directos:** acciones para iniciar (`Start`), apagar ordenadamente (`Shutdown`), apagar de forma forzada (`Power Off`) y reiniciar (`Reboot`) una instancia.
- **Edición básica de recursos:** capacidad de reconfigurar la memoria RAM o las vCPUs asignadas a una instancia, en caliente o después de un reinicio según las capacidades del recurso.
- **Gestión de snapshots:** posibilidad de crear, listar y revertir un punto de restauración antes de realizar cambios críticos en una máquina.

### 3. Telemetría y monitoreo en tiempo real

- **Gráficos dinámicos por instancia:** monitor de recursos en vivo con consumo porcentual de vCPU, memoria asignada y consumida, tráfico de red entrante y saliente, y lectura y escritura de disco.
- **Histórico corto:** gráficos de tendencias para analizar picos de uso durante las últimas horas.

### 4. Creación y configuración simplificada

- **Asistente de aprovisionamiento (wizard):** flujo guiado en pocos pasos para desplegar un nuevo contenedor o VM seleccionando recursos clave, como vCPUs, RAM, tamaño de disco y plantilla o imagen ISO base.

### 5. Seguridad y auditoría

- **Autenticación y roles:** acceso mediante usuario y contraseña, con roles diferenciados, por ejemplo, administrador del sistema y operador o usuario de instancia.
- **Registro de actividad (audit log):** historial cronológico de quién realizó cada acción sobre la infraestructura, por ejemplo: «Usuario X reinició la VM 102 el día Y».

## Funcionalidades excluidas

### Consola remota web

La consola remota web (terminal SSH/VNC, xterm.js o noVNC) se excluye intencionalmente de la plataforma para mitigar riesgos operativos y evitar la ejecución de comandos arbitrarios a nivel del sistema operativo. El acceso por consola directa queda reservado exclusivamente para administradores de infraestructura mediante las herramientas nativas de Proxmox VE.




