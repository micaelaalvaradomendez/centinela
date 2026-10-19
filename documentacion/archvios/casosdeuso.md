Actores
usuario: persona que interactúa a través de la interfaz web
api proxmox: sistema externo que interactúa a través de una api y rest/json
Casos de uso
CU-01 autenticarse e iniciar sesión
requerimiento relacionado: RF-01
actor principal: usuario (administrador / operador).
descripción: el usuario ingresa credenciales básicas y, posteriormente, valida su identidad mediante un código de autenticación de doble factor (“FA) para acceder a la plataforma.
transacciones:
el usuario solicita la pantalla de ingreso y el sistema la renderiza.
el usuario envía sus credenciales (usuario/contraseña) al backend.
el backend valida las credenciales, y verifica en la base de datos si el usuario tiene el 2FA vinculado.
CONDICIONAL:
sin vinculación: el backend genera la clave secreta y la URL del código QR, notificando al frontend
sin vinculación: el frontend renderiza la pantalla con el código QR para que el usuario vincule su aplicacion autenticadora
con vinculación: el frontend renderiza la pantalla para solicitar el código TOTP de 6 dígitos
el usuario ingresa el código del 2FA y lo envía al backend
el backend valida el código 2FA ingresado, firma el token JWT con la sesion completamente autenticada y responde exitosamente
el frontend almacena el JWT y redirige al dashboard

CU-02: visualizar dashboard general del servidor
requerimiento asociado: RF-02.
actor principal: usuario (administrador / operador).
descripción: el usuario consulta el estado de salud e hipervisor (cpu, ram, disco, instancias).
transacciones:
el usuario navega a la vista principal/dashboard.
el frontend solicita el estado de salud al backend.
el backend consulta la api de proxmox para leer recursos del host.
proxmox responde con la telemetría del servidor.
el backend formatea la información y la envía al frontend.
el frontend actualiza los componentes gráficos del panel.

CU-03: filtrar e inventariar instancias (vms y lxc)
requerimiento asociado: RF-03.
actor principal: usuario (administrador / operador).
descripción: muestra la lista filtrable de instancias asociadas al rol del usuario.
transacciones:
el usuario solicita la lista de instancias.
el backend intercepta la solicitud y valida los permisos del rol.
el backend solicita el inventario a la api de proxmox.
el backend unifica las propiedades json (vms y lxc).
el frontend renderiza la tabla de inventario.
el usuario aplica filtros (tipo, estado o búsqueda por nombre).
el frontend procesa y filtra la lista.

CU-04: ejecutar acción de control de energía (manejo de UPID)
requerimiento asociado: RF-04.
actor principal: usuario (administrador / operador).
descripción: iniciar, apagar, forzar apagado o reiniciar una máquina o contenedor de forma asíncrona.
transacciones:
el usuario presiona una acción de energía (ej. apagar).
la interfaz abre un modal solicitando confirmación explícita.
el usuario confirma la acción.
el frontend envía la orden al backend.
el backend valida los permisos e invoca la orden a proxmox.
proxmox retorna un identificador de tarea asíncrona (UPID).
el frontend bloquea los controles de esa instancia y muestra el estado "en progreso".
el backend monitorea en segundo plano la evolución del upid en proxmox.
el backend notifica la finalización de la tarea al frontend.
el frontend desbloquea los botones e informa el resultado.
el backend registra la acción en la tabla de auditoría.

CU-05: consultar telemetría y métricas por instancia
requerimiento asociado: RF-05.
actor principal: usuario (administrador / operador).
descripción: visualización de gráficos de consumo (cpu, ram, disco, red) en tiempo real e histórico corto.
transacciones:
el usuario ingresa a la vista de detalle de una instancia.
el frontend establece suscripción de métricas.
el backend consulta métricas en vivo e histórico a proxmox/caché.
el backend emite el flujo de métricas hacia el frontend.
el frontend renderiza los gráficos de líneas en vivo.
el usuario selecciona un rango histórico diferente (ej. últimas 2 horas).

CU-06: administrar snapshots (puntos de restauración)
requerimiento asociado: RF-06.
actor principal: usuario (administrador / operador).
descripción: crear, listar y revertir un punto de restauración de una instancia.
transacciones:
el usuario abre la sección de snapshots de una instancia.
el backend obtiene la lista de snapshots desde proxmox y la envía al frontend.
el usuario completa el formulario (nombre y descripción) para crear uno nuevo.
el backend envía la orden de creación de snapshot a proxmox.
el usuario solicita la reversión a un snapshot existente.
el sistema muestra una alerta de confirmación destructiva.
el backend captura el upid del rollback y notifica la finalización.

CU-07: aprovisionar nueva instancia (wizard)
requerimiento asociado: RF-07.
actor principal: usuario (administrador / operador).
descripción: asistente paso a paso para crear vms/lxc ocultando la complejidad con validación previa de cuotas.
transacciones:
el usuario inicia el asistente de creación (wizard).
el usuario selecciona el tipo de instancia (vm o lxc).
el usuario define la cuota de recursos (vcpus, ram, disco).
el usuario selecciona la plantilla/iso e ingresa configuración de red.
el usuario confirma el resumen de creación.
el backend consulta a proxmox la capacidad libre real del nodo físico.
el backend valida las cuotas y aprueba/rechaza la solicitud.
el backend manda la orden a proxmox, rastrea el upid y asocia la vm al usuario.

CU-08: consultar registros de auditoría
requerimiento asociado: RF-08.
actor principal: usuario administrador.
descripción: visualización e inspección inalterable de los logs de acciones de la plataforma.
transacciones:
el administrador accede a la vista de auditoría.
el frontend consulta los registros paginados.
el backend ejecuta la consulta indexada sobre la base de datos.
el administrador aplica un filtro (por usuario, fecha o acción).
el backend devuelve la lista ordenada y filtrada.

CU-09: gestionar usuarios, roles y permisos de instancia
requerimiento asociado: RF-09.
actor principal: usuario administrador.
descripción: crear o eliminar usuarios, cambiar roles y vincular qué operador puede controlar qué instancia.
transacciones:
el administrador accede al panel de usuarios.
el administrador crea un nuevo usuario y le asigna un rol.
el backend guarda las credenciales en la base de datos aplicando un hash seguro.
el administrador selecciona un operador y le asigna un listado de instancias permitidas.
el backend actualiza la tabla intermedia de permisos en la base de datos.

CU-10: modificar recursos de instancia
requerimiento asociado: RF-10.
actor principal: usuario (administrador / operador con permisos).
descripción: reconfiguración de vcpus y ram asignadas a una instancia con validación de disponibilidad.
transacciones:
el usuario abre el formulario de modificación de recursos de la máquina.
el usuario ajusta los valores de vcpus y memoria ram.
el backend calcula el diferencial de recursos solicitados.
el backend verifica la ram libre en el servidor físico.
el backend aprueba la actualización y envía la orden de reconfiguración a proxmox.
CU-11: notificaciones
requerimiento asociado: RF-11.
actor principal: usuario (administrador / operador con permisos)
descripción:el sistema monitorea el estado de la infraestructura y emite alertas emergentes en tiempo real en la interfaz. Al hacer clic en la notificación relacionada con una instancia, la interfaz redirige automáticamente al usuario a la pantalla de detalle de dicha instancia, donde se presenta las métricas en vivo y el panel de control para ejecutar acciones correctivas como apagar, reiniciar o ajustar recursos.
transacciones:
El usuario está navegando en cualquier sección del panel web.
Una máquina virtual se satura o se apaga inesperadamente en Proxmox o un proceso termina
backend detecta evento en proxmox
backend emite evento con tipo de alerta, severidad e instancia afectada
frontend renderiza componente mostrando mensaje y enlace de redireccionamiento
el usuario al hacer clic en ver detalle
el frontend captura el evento de navegación y redirige a la vista específica de la instancia
frontend renderiza métricas en vivo y habilita botones de control
el backend registra  emisión de alerta
