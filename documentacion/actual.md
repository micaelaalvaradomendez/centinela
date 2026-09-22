
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:


### `BAC-14` - Lectura mínima del inventario de Proxmox

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 3 h
- **Ventana propuesta:** 17/09/2026, 11:00-14:00
- **Depende de:** `BAC-08` y de las credenciales de lectura de Proxmox VE.
- **Entregable:** `GET /api/instances` consumiendo `/cluster/resources` o `/nodes/{node}/resources`, con una respuesta normalizada mínima que incluya ID, nombre, tipo, nodo y estado. Un administrador recibe todo el inventario y un operador solo las instancias asignadas.
- **Criterio de éxito:** `FRN-07` puede cargar IDs reales de VMs y LXC; una cuenta no puede descubrir instancias fuera de su alcance y los errores de Proxmox se traducen a una respuesta HTTP controlada.



### `FRN-10` - Cambio obligatorio de contraseña temporal

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 2 h
- **Ventana propuesta:** 16/09/2026, 11:00-13:00
- **Depende de:** `BAC-12` y `FRN-04`.
- **Entregable:** vista de nueva contraseña y confirmación que detecte `must_change_password`, bloquee la navegación general y llame a `POST /api/auth/change-password`.
- **Criterio de éxito:** una cuenta con contraseña temporal solo puede cerrar sesión o cambiarla; después del cambio continúa al enrolamiento o validación 2FA que corresponda.

### `FRN-11` - Acciones administrativas de recuperación

- **Área:** Frontend
- **Asignada:** Luz
- **Estimación:** 2 h
- **Ventana propuesta:** 18/09/2026, 12:00-14:00
- **Depende de:** `FRN-05`, `BAC-13` y `BAC-15`.
- **Entregable:** acciones separadas para restablecer contraseña y 2FA desde el panel de usuarios, ambas con confirmación explícita, estado de carga y notificación del resultado.
- **Criterio de éxito:** un administrador puede iniciar cada recuperación sin confundir sus efectos; la tabla refleja que el 2FA quedó desvinculado y nunca muestra secretos ni hashes.

### `FRN-12` - Vistas de recuperación de contraseña (RF-13)

- **Área:** Frontend
- **Asignada:** Belinda
- **Estimación:** 3 h
- **Ventana propuesta:** A definir (posterior a `BAC-20`).
- **Depende de:** `BAC-19`, `BAC-20` y el maquetado existente de `RecoverPassword.tsx`.
- **Entregable:** conectar `RecoverPassword.tsx` al flujo real: paso de ingreso de correo, paso de ingreso del código de seis dígitos y paso de nueva contraseña, con manejo de errores del backend en cada paso.
- **Criterio de éxito:** una cuenta puede recuperar el acceso sin intervención de un administrador, y los errores de código inválido o vencido se muestran junto al campo correspondiente.

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

**Hito:** *Control de Acceso Basado en Recursos*
* El administrador abre un usuario en `FRN-07`, el frontend lista las instancias de Proxmox (`BAC-14`), selecciona un subconjunto y guarda (`BAC-07`).
* Al iniciar sesión como ese Operador, `GET /api/instances` solo devuelve las instancias que tiene permitidas.
* Si intenta forzar una petición sobre una instancia no asignada, el backend rechaza con `403` (`BAC-08`) y el frontend muestra el error en un toast sin cerrar la sesión (`FRN-08`).
---


