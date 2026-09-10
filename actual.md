
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

#### ⚙️ Lista Backend (Meta Fin de Semana: Login Funcional)

* `BAC-01`: **Script Docker Compose de BD Local:** Levantar contenedor PostgreSQL con script `.sql` inicial (`users`, roles `ADMIN`/`OPERATOR`) y usuario seed de prueba. *(Asignado: Tayra / Lucas - 2h)*


* `BAC-02`: **Modelo y Hashing de Contraseñas:** Configurar librería de hash (bcrypt o argon2) para verificar contraseñas. *(Asignado: Lisandro - 2h)*


* `BAC-03`: **Endpoint POST /api/auth/login:** Recibir usuario/contraseña, validar contra PostgreSQL y devolver JWT firmado con claims (`id`, `role`). *(Asignado: Tayra - 3h)*


* `BAC-04`: **Manejo de Errores HTTP:** Retornar 401 para credenciales inválidas y 400 para payloads incompletos con respuestas JSON unificadas. *(Asignado: Lisandro - 1.5h)*



#### 🎨 Lista Frontend (Meta Domingo: Pantallas Visibles)

* `FRN-01`: **Maquetado de Formulario Login:** Campos de email/usuario, contraseña, botón de envío y estado de carga visual. *(Asignado: Belinda - 3h)*


* `FRN-02`: **Validación en Cliente (Login):** Validar formato de campos y mensajes de error antes de enviar la petición. *(Asignado: Belinda - 2h)*


* `FRN-03`: **Navbar y Rutas Base:** Crear barra de navegación con botones principales y vistas esqueleto (Dashboard, Instancias) para verificar flujo de navegación. *(Asignado: Cristian / Luz - 3h)*


* `FRN-04`: **Conexión Login con Localhost:** Configurar llamada HTTP (`fetch`/`axios`) apuntando al endpoint local del backend y almacenar JWT en memoria. *(Asignado: Cristian - 3h)*

