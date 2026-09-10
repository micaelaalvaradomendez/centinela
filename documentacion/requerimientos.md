# Requerimientos de El Centinela

## Requerimientos funcionales

### RF-01. Autenticación y control de acceso

El sistema debe permitir el registro y el inicio de sesión mediante usuario o correo electrónico y contraseña. Después de validar las credenciales, debe exigir autenticación de doble factor (2FA) mediante contraseñas temporales basadas en tiempo (TOTP).

Si el usuario no tiene configurado el 2FA, el sistema debe mostrar la pantalla de vinculación inicial mediante código QR y permitir también la carga manual del secreto. Si ya está vinculado, debe solicitar un código TOTP de seis dígitos. Debe existir un mecanismo para solicitar la regeneración o revinculación del 2FA en caso de pérdida de acceso a la aplicación autenticadora.

El control de acceso debe estar basado en roles, con al menos los roles de administrador y operador. El sistema debe generar y validar tokens JWT para gestionar sesiones seguras.

**Frontend**

- Diseñar y maquetar los formularios de registro e inicio de sesión, con validación de campos.
- Implementar el almacenamiento seguro del JWT y el gestor de estado global de la sesión.
- Restringir las vistas según el rol incluido en el JWT.
- Agregar el encabezado de autorización a las peticiones HTTP.
- Diseñar y renderizar la pantalla o modal de vinculación inicial del 2FA mediante código QR.
- Diseñar la vista para ingresar el código TOTP de seis dígitos y solicitar su regeneración o revinculación.

**Backend**

- Crear los endpoints de autenticación y verificación del 2FA.
- Proteger el almacenamiento de contraseñas y secretos TOTP.
- Diseñar la firma y verificación de los JWT, incluyendo el identificador de usuario, el rol, el estado de verificación del 2FA y el tiempo de expiración.
- Validar los permisos según el rol solicitado.
- Generar secretos TOTP y códigos QR.
- Validar los códigos ingresados y gestionar la revinculación o generación de nuevas claves 2FA.

**Infraestructura**

- Desplegar y mantener la persistencia de usuarios y credenciales.
- Configurar de forma segura las claves utilizadas para firmar los JWT.
- Garantizar conexiones HTTPS/TLS entre el cliente y el servidor.

### RF-02. Dashboard de estado del nodo

Al iniciar el sistema, debe mostrar una vista con el estado de salud y el consumo de recursos del nodo físico o clúster de Proxmox VE. Debe incluir:

- CPU total, porcentaje de uso y cantidad de núcleos.
- Memoria RAM utilizada, porcentaje y capacidad total.
- Almacenamiento utilizado, porcentaje y capacidad total.
- Tiempo de actividad del nodo.
- Cantidad de instancias VM y LXC, agrupadas por estado: en ejecución, detenidas u otros estados disponibles.

Este requisito permite conocer la salud del hipervisor y determinar si el servidor está sobrecargado o dispone de recursos.

**Frontend**

- Diseñar e integrar tarjetas de resumen y medidores para el consumo del host.
- Actualizar los valores sin recargar la página mediante el mecanismo de tiempo real definido por el equipo (polling, WebSockets o una combinación de ambos).
- Implementar un indicador visual del estado del hipervisor, por ejemplo, normal, degradado o inaccesible.

**Backend**

- Consultar periódicamente el estado de Proxmox VE.
- Procesar, normalizar y distribuir la información hacia el frontend.

**Infraestructura**

- Garantizar que el token de API de Proxmox VE tenga permisos de lectura sobre el estado del nodo.
- Configurar una caché con el último estado conocido para responder rápidamente durante períodos de alta concurrencia.

### RF-03. Inventario de instancias

El sistema debe presentar una lista de todas las instancias gestionadas por Proxmox VE. Cada instancia debe mostrar su ID, nombre, tipo (VM o LXC), estado, dirección IP y un resumen del uso de CPU y RAM. La lista debe permitir buscar y filtrar por nombre, estado y tipo.

**Frontend**

- Diseñar una vista clara que combine máquinas virtuales y contenedores LXC.
- Implementar un buscador dinámico y filtros por tipo y estado.
- Diseñar etiquetas visuales para los estados.
- Sincronizar el inventario con los eventos de cambio de estado recibidos del backend.

**Backend**

- Desarrollar el endpoint de consulta de instancias.
- Filtrar las instancias según la matriz de acceso del usuario.
- Normalizar las propiedades de las VMs y los LXC en una respuesta JSON común.

**Infraestructura**

- Verificar que el token de API tenga permisos para consultar las direcciones IP de cada VM o LXC.

### RF-04. Control del ciclo de vida de instancias

El sistema debe permitir iniciar, apagar ordenadamente, forzar el apagado y reiniciar instancias. Cada acción debe solicitar confirmación explícita. El estado de la instancia debe actualizarse en tiempo real mientras la operación progresa.

El backend debe manejar el UPID devuelto por Proxmox VE, monitorear la tarea en segundo plano y notificar al frontend cuando finalice. El frontend debe bloquear acciones duplicadas mientras la tarea esté en progreso.

**Frontend**

- Diseñar botones de control rápido con modales de confirmación.
- Mostrar el estado «operación en progreso» y bloquear los controles de la instancia durante la tarea.
- Escuchar el resultado asociado al UPID, volver a habilitar los controles e informar si la operación tuvo éxito o falló.

**Backend**

- Crear el endpoint para enviar órdenes de energía a las instancias.
- Capturar el UPID devuelto por Proxmox VE.
- Implementar un monitor en segundo plano para consultar el estado de la tarea.
- Notificar la finalización y registrar la operación en la auditoría.

**Infraestructura**

- Asignar al token de API los permisos de gestión de energía necesarios.
- Configurar colas o un bus de mensajes para procesar los UPID en segundo plano sin bloquear la API.

### RF-05. Monitor de métricas por instancia

El sistema debe mostrar gráficos dinámicos con las métricas de consumo de cada instancia en tiempo real. Debe incluir CPU, RAM, tráfico de red y uso de disco, además de un histórico reciente.

Este requisito es diferente del RF-02: permite identificar qué instancia consume más recursos y detectar picos de uso.

**Frontend**

- Diseñar vistas de detalle por instancia con gráficos de líneas.
- Procesar flujos de datos en tiempo real sin degradar el consumo de memoria del navegador.
- Ofrecer selectores para consultar el histórico corto, por ejemplo, las últimas dos o 24 horas.

**Backend**

- Leer los datos de rendimiento desde la API de Proxmox VE.
- Transmitir las métricas en vivo por instancia.
- Mantener una caché de histórico corto para responder consultas de tendencias.

**Infraestructura**

- Verificar que Proxmox VE mantenga activos los servicios de recolección de métricas.
- Monitorear el almacenamiento temporal para evitar un consumo excesivo de memoria o disco.

### RF-06. Gestión de snapshots

El sistema debe permitir administrar snapshots de las instancias. Debe permitir crear un snapshot con nombre y descripción, listar los snapshots de una instancia y revertir una instancia a un snapshot seleccionado. La reversión debe tratarse como una acción destructiva.

**Frontend**

- Construir una vista o modal de gestión de snapshots dentro del detalle de la instancia.
- Diseñar el formulario de creación con nombre y descripción.
- Mostrar una tabla con el historial de puntos de restauración.
- Implementar un flujo de confirmación estricto y una advertencia visual para la reversión.

**Backend**

- Crear endpoints para listar, crear y solicitar la reversión de snapshots.
- Gestionar la reversión como una tarea asíncrona, capturando el UPID y notificando su evolución.
- Validar los permisos del usuario sobre la instancia antes de autorizar la reversión.

**Infraestructura**

- Otorgar al token de API los permisos necesarios para gestionar almacenamiento y snapshots.
- Garantizar que los almacenamientos configurados en Proxmox VE tengan habilitado el soporte nativo de snapshots.

### RF-07. Asistente de creación de instancias

El sistema debe proporcionar un asistente paso a paso para crear nuevas VMs o contenedores LXC. El wizard debe guiar al usuario en la selección de:

- Tipo de instancia: VM o LXC.
- Recursos: vCPUs, RAM y disco.
- Imagen o plantilla base: ISO para VMs o template para LXC.
- Configuración de red: dirección IP y puerta de enlace.

El asistente debe ocultar la complejidad de la configuración nativa de Proxmox VE. Antes de crear la instancia, el backend debe validar que existan recursos suficientes de CPU, RAM y disco. Si varias solicitudes compiten por los mismos recursos, debe rechazarse la solicitud que llegue cuando estos ya no estén disponibles, evitando condiciones de carrera.

**Frontend**

- Diseñar el wizard multipaso con navegación progresiva.
- Mostrar indicadores dinámicos de las cuotas disponibles mientras se ajustan los recursos.
- Bloquear el formulario durante el envío.
- Procesar errores por cuotas insuficientes y mostrar un mensaje amigable ante un conflicto de recursos, por ejemplo, HTTP 409.

**Backend**

- Desarrollar el endpoint de aprovisionamiento.
- Validar las cuotas y la disponibilidad de RAM y disco antes de llamar a Proxmox VE.
- Aplicar un mecanismo de reserva o bloqueo transaccional para evitar condiciones de carrera.
- Responder con un error específico, como HTTP 409, si los recursos se agotaron.
- Enviar el comando a Proxmox VE, rastrear el UPID y asociar la instancia con el usuario cuando existan recursos suficientes.

**Infraestructura**

- Otorgar al token de API permisos para crear máquinas y asignar almacenamiento.
- Mantener disponibles las plantillas LXC y las imágenes ISO base en los almacenamientos de Proxmox VE.

### RF-08. Registro de auditoría

El sistema debe registrar de forma inmutable todas las acciones realizadas a través de la plataforma. Cada registro debe incluir usuario, fecha y hora, acción, instancia afectada (ID y nombre) y resultado (éxito o fallo). Los registros deben poder consultarse, filtrarse por usuario, fecha, tipo de acción, resultado o instancia, y exportarse.

**Frontend**

- Diseñar una tabla de auditoría exclusiva para administradores.
- Implementar filtros por fecha, usuario, tipo de acción y resultado.
- Diseñar paginación eficiente.
- Integrar un botón de exportación que conserve los filtros aplicados.

**Backend**

- Diseñar la tabla de auditoría con los índices correspondientes.
- Crear un servicio que registre cada operación de forma inalterable.
- Exponer un endpoint con filtros, paginación y ordenamiento.
- Crear un endpoint de exportación que genere el archivo según los filtros solicitados.

**Infraestructura**

- Garantizar la persistencia segura y los respaldos periódicos de la base de datos.
- Aislar la base de datos para evitar modificaciones externas directas.

### RF-09. Gestión de usuarios y permisos

Los administradores deben poder crear y eliminar usuarios, asignarles roles y definir las instancias a las que tendrá acceso cada operador. Al crear una cuenta, el sistema debe generar una contraseña temporal, enviarla al nuevo usuario y exigir su cambio durante el primer inicio de sesión.

**Frontend**

- Diseñar el panel de administración de usuarios y permisos.
- Crear una interfaz para vincular operadores con sus instancias permitidas.
- Maquetar las vistas de creación, edición de roles y eliminación con modales de confirmación.

**Backend**

- Desarrollar el CRUD de usuarios, accesible solo para administradores.
- Generar una contraseña temporal segura y enviarla por correo junto con las instrucciones de acceso.
- Marcar la contraseña como temporal y obligar a cambiarla durante el primer inicio de sesión.
- Diseñar la tabla intermedia para la relación entre usuarios e instancias.
- Implementar autorización por recurso para verificar el acceso de un operador a cada instancia.
- Configurar el servicio SMTP o el proveedor de correo para entregar las credenciales iniciales.

**Infraestructura**

- Asegurar que la lógica de roles permanezca aislada dentro de la plataforma, sin modificar usuarios del sistema operativo ni de Proxmox VE.

### RF-10. Edición de recursos de instancias

El sistema debe permitir a los usuarios autorizados modificar la cantidad de vCPUs y RAM asignadas a una instancia, en caliente o después de un reinicio según las capacidades de la instancia. El backend debe validar la disponibilidad de recursos antes de aplicar el cambio.

**Frontend**

- Diseñar un modal o formulario de edición dentro del detalle de la instancia.
- Proporcionar controles numéricos precargados con los valores actuales.
- Mostrar advertencias si la modificación requiere reiniciar la instancia.

**Backend**

- Crear el endpoint de reconfiguración de RAM y vCPUs.
- Calcular la diferencia de recursos y verificar la disponibilidad de RAM y CPU en el nodo.
- Rechazar el cambio si no hay recursos suficientes o enviar la orden a Proxmox VE si la validación es correcta.

**Infraestructura**

- Otorgar al token de API permisos de configuración de CPU y memoria.
- Habilitar el soporte de hot-plug en las máquinas virtuales cuando se requiera ajustar recursos sin apagarlas.

### RF-11. Notificaciones

El sistema debe emitir notificaciones en la interfaz sin recargar la pantalla y enviar alertas por correo ante situaciones de alta severidad. Debe notificar cambios de estado, creación de instancias, saturación de recursos y finalización de tareas, indicando si tuvieron éxito o fallaron.

Cada notificación relacionada con una VM o LXC debe incluir un enlace hacia el detalle de la instancia y sus acciones disponibles.

**Frontend**

- Implementar un contenedor persistente de notificaciones flotantes, por ejemplo, un componente toast.
- Asociar el ID de la instancia afectada para habilitar la navegación.
- Escuchar las alertas recibidas del backend y filtrarlas por instancia cuando corresponda.
- Usar códigos de color para diferenciar los tipos de evento.

**Backend**

- Emitir eventos de cambio de estado asociados con la instancia afectada.
- Evaluar la salud y el consumo de recursos del nodo o clúster para detectar saturación.
- Emitir alertas con el resultado de las acciones solicitadas por el frontend.
- Integrar un servicio SMTP o un proveedor externo con disponibilidad suficiente para enviar correos.

**Infraestructura**

- Asegurar que el bus interno de datos, como Kafka, transmita los eventos sin demoras ni pérdidas.

### RF-12. Gestión de organizaciones

El sistema debe soportar múltiples organizaciones operando de forma aislada.

**Frontend**

- Diseñar los formularios de creación de organizaciones y edición de sus datos, como el nombre.
- Implementar la gestión de miembros, permitiendo que un administrador cree usuarios o envíe invitaciones.
- Mostrar la organización asignada como un dato no modificable para los usuarios estándar.

**Backend**

- Crear organizaciones y asignar automáticamente el rol de administrador inicial al usuario creador.
- Asociar nuevos usuarios a la organización del administrador que los creó o invitó.
- Rechazar cualquier petición de usuarios estándar que intente modificar su organización.
- Vincular en la base de datos las instancias, roles y registros de auditoría con una organización específica.
- Validar en todos los endpoints que cada administrador solo consulte y gestione recursos de su propia organización.
- Permitir la edición de los datos de la organización únicamente a sus administradores.

**Infraestructura**

- Garantizar el aislamiento lógico y de datos entre organizaciones.
- Configurar respaldos y controles de acceso que preserven la separación de la información de cada organización.

### RF-13. Recuperación de contraseña

El sistema debe permitir restablecer una contraseña olvidada mediante un código temporal de verificación enviado al correo electrónico del usuario.

**Frontend**

- Diseñar el formulario para ingresar el correo y solicitar la recuperación.
- Implementar la vista para ingresar el código temporal de seis dígitos.
- Procesar visualmente el resultado de la solicitud y los errores del backend.

**Backend**

- Validar el formato del correo y buscar la cuenta asociada.
- Generar una solicitud de recuperación y un código temporal de seis dígitos.
- Asociar el código con la cuenta y la solicitud, estableciendo una fecha y hora de expiración.
- Enviar el código al correo de la cuenta.
- Invalidar los códigos anteriores al generar uno nuevo y rechazar códigos vencidos.
- Devolver respuestas y errores estructurados para que el frontend gestione el flujo.

**Infraestructura**

- Configurar y proteger el servicio SMTP o un proveedor de correo externo.
- Garantizar la persistencia temporal de los códigos y solicitudes en la base de datos o caché.

## Requerimientos no funcionales

### RNF-01. Seguridad

**Mínimo privilegio:** el backend no debe utilizar credenciales de administrador o root de Proxmox VE. Debe comunicarse mediante tokens de API con permisos limitados a las operaciones necesarias y con la separación de privilegios activada.

**Credenciales:** las contraseñas deben almacenarse mediante funciones de hashing seguras. Los secretos TOTP y los tokens de la API de Proxmox VE deben protegerse y nunca almacenarse en texto plano.

**Frontend**

- Manipular únicamente el JWT de sesión emitido por la plataforma.
- No almacenar contraseñas, secretos TOTP ni tokens maestros en texto plano.
- Limpiar el estado de sesión y redirigir al inicio de sesión cuando expire la sesión.

**Backend**

- Cifrar los secretos TOTP y el token de API de Proxmox VE mediante claves gestionadas de forma segura.
- Utilizar hashing seguro para las contraseñas.
- No exponer secretos TOTP, tokens ni claves maestras en las respuestas JSON.

**Infraestructura**

- Crear un usuario de servicio exclusivo en Proxmox VE sin permisos de administración del host.
- Generar un token de API con la separación de privilegios activada.
- Asignar ACL específicas únicamente sobre los recursos que gestionará el sistema.
- Garantizar conexiones HTTPS/TLS para proteger las credenciales y los códigos TOTP en tránsito.

### RNF-02. Rendimiento

El sistema debe manejar múltiples solicitudes y operaciones simultáneas sin degradación significativa. En condiciones normales, las respuestas deben producirse en un máximo de cinco segundos.

**Frontend**

- Optimizar el renderizado para evitar actualizaciones innecesarias ante eventos masivos.
- Implementar paginación y carga diferida en tablas con gran volumen de datos.

**Backend**

- Utilizar un marco de trabajo asíncrono y concurrente con bucle de eventos.
- Separar la recolección de métricas de las peticiones REST.

**Infraestructura**

- Asegurar recursos suficientes de memoria y ancho de banda para cachés y colas.
- Configurar el servidor web o proxy de entrada para gestionar conexiones concurrentes y persistentes.

### RNF-03. Usabilidad

La interfaz debe ser limpia, moderna e intuitiva, especialmente para desarrolladores. Los flujos del wizard deben ser guiados y las acciones críticas deben incluir confirmaciones explícitas y advertencias claras.

**Frontend**

- Construir una interfaz homogénea mediante una librería moderna de componentes.
- Ocultar parámetros avanzados del kernel y configuraciones complejas de red.
- Diseñar confirmaciones con advertencias visuales para acciones críticas o destructivas.

**Backend**

- Estandarizar los mensajes de error y traducir las excepciones técnicas de Proxmox VE a textos comprensibles.

**Infraestructura**

- Desplegar el sistema en un dominio o dirección accesible internamente por el equipo.

### RNF-04. Confiabilidad

El sistema debe manejar errores de la API de Proxmox VE, incluidos timeouts, tokens inválidos, recursos inexistentes y operaciones fallidas. Los mensajes deben explicar claramente qué ocurrió.

Todas las operaciones que generen una tarea de Proxmox VE (UPID) deben ejecutarse de forma asíncrona. El backend debe iniciar la tarea, devolver inmediatamente un identificador, monitorear el proceso en segundo plano y notificar al frontend cuando finalice o falle. El frontend debe reflejar el estado «en progreso».

**Frontend**

- Bloquear acciones duplicadas sobre una misma instancia.
- Mostrar indicadores de carga durante las tareas.
- Mostrar notificaciones flotantes adaptativas con el resultado final.

**Backend**

- Capturar excepciones de las llamadas a la API de Proxmox VE y devolver respuestas estructuradas.
- Implementar el motor de tareas en segundo plano, con reintentos y tiempos de espera al consultar UPID.
- Mantener un estado degradado controlado durante caídas temporales de Proxmox VE.

**Infraestructura**

- Configurar políticas de reinicio automático para la base de datos, la caché y el backend.

### RNF-05. Mantenibilidad

El código debe estar documentado, especialmente la API del backend, para facilitar el mantenimiento y la incorporación de nuevas funcionalidades.

**Frontend**

- Estructurar el código con componentes reutilizables.
- Centralizar los servicios de comunicación con la API y los WebSockets.

**Backend**

- Implementar una arquitectura clara y fácil de mantener.
- Configurar la generación automática de documentación interactiva de la API.

**Infraestructura**

- Documentar los procedimientos de despliegue, los scripts de base de datos y los archivos de orquestación.

### RNF-06. Compatibilidad

La interfaz debe ser compatible con las versiones recientes de Google Chrome, Mozilla Firefox, Brave, Microsoft Edge y Safari.

**Frontend**

- Probar gráficos, estilos y WebSockets de forma uniforme en navegadores modernos.

**Backend**

- Configurar encabezados CORS para permitir peticiones únicamente desde orígenes autorizados.

**Infraestructura**

- Asegurar que el proxy inverso gestione HTTP, HTTPS y los encabezados necesarios para actualizar la conexión a WebSockets.

### RNF-07. Escalabilidad

El sistema debe soportar la gestión de múltiples nodos Proxmox VE y permitir el escalamiento horizontal sin modificar la lógica central.

**Frontend**

- Diseñar vistas preparadas para seleccionar o agrupar recursos por nodo sin rehacer los componentes.

**Backend**

- Desacoplar la identificación del nodo de la lógica de servicios para permitir el enrutamiento dinámico hacia múltiples endpoints de Proxmox VE.

**Infraestructura**

- Preparar la arquitectura de contenedores y servicios de backend para permitir réplicas cuando aumente la carga.

## Funcionalidades excluidas

### Consola remota web

La consola remota web, incluyendo terminal SSH/VNC, xterm.js y noVNC, queda fuera del alcance para mitigar el riesgo de configuraciones erróneas y evitar la ejecución de comandos arbitrarios en el sistema operativo. El acceso por consola directa queda reservado a los administradores de infraestructura mediante las herramientas nativas de Proxmox VE.
