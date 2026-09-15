
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

Las tareas que ya tenés cargadas en tus capturas (con vencimiento viernes/mañana) corresponden al **bloque base del Login** (`BAC-01` a `BAC-04` y `FRN-01` a `FRN-04`).

Las tareas que te quedan por cargar para completar el circuito son **el flujo 2FA/TOTP y la integración final**:

* `LOGIN-01`: Flujo 2FA/TOTP en memoria (Backend).


* `LOGIN-02`: Pantalla y validación del código TOTP (Frontend).


* `LOGIN-03`: Integración y prueba del prototipo (Front + Back).



---

### 1. Mapa de Correlaciones y Dependencias (Para configurar en ClickUp)

| Lista en ClickUp | Tarea a Cargar | ⚠️ Bloqueada por (Waiting on) | ⛔ Bloquea a (Blocking) |
| --- | --- | --- | --- |
| **backend**<br> | **`LOGIN-01`: Flujo 2FA/TOTP en memoria**<br> | `BAC-03` (POST /login) y `BAC-04` (Errores HTTP)

 | `LOGIN-02`<br> |
| **frontend**<br> | **`LOGIN-02`: Pantalla y validación del código TOTP**<br> | `LOGIN-01` (Back) y `FRN-01` (Maquetado)

 | `LOGIN-03`<br> |
| **frontend** o **backend** | **`LOGIN-03`: Integración y prueba del prototipo**<br> | `BAC-04`, `FRN-04` y `LOGIN-02`<br> | *(Cierra el prototipo funcional)*<br> |

(Nota: en tu tarea ya creada **`FRN-04`**, asegurate de que esté bloqueada por `BAC-03`, `BAC-04` y `FRN-01`).

---

#### ⚙️ En la Lista `backend`:

**Nombre:** `LOGIN-01: Flujo 2FA/TOTP en memoria`

* **Descripción para pegar en la tarjeta:**
```text
OBJETIVO:
Generar y validar el código TOTP de 6 dígitos tras validar usuario y contraseña (en memoria).

PASO A PASO:
1. Generar clave secreta TOTP temporal al autenticar.
2. Endpoint o validación que reciba el código de 6 dígitos.
3. Rechazar códigos vencidos o inválidos (retornar error unificado).

CRITERIO DE ÉXITO (DONE):
El login exige el código de 6 dígitos y solo entrega acceso si el código es válido.

```



---

#### 🎨 En la Lista `frontend`:

**Nombre:** `LOGIN-02: Pantalla y validación de código TOTP`

* **Descripción para pegar en la tarjeta:**
```text
OBJETIVO:
Pantalla para que el usuario ingrese los 6 dígitos del código 2FA.

PASO A PASO:
1. Maquetar vista/input para el código de 6 dígitos.
2. Mostrar estado de carga y mensajes de error si el código es incorrecto.
3. Al confirmar éxito, permitir la navegación a la pantalla base (Dashboard).

CRITERIO DE ÉXITO (DONE):
La pantalla permite tipear los 6 números, muestra error visual si falla y solo avanza si el back lo aprueba.

```



---

#### 🚀 Integración:

**Nombre:** `LOGIN-03: Integración y prueba punta a punta del prototipo`

* **Descripción para pegar en la tarjeta:**
```text
OBJETIVO:
Probar el flujo completo de inicio a fin en una llamada de 10-15 minutos.

CIRCUITO A PROBAR:
Usuario seed -> Ingreso de contraseña -> Solicitud de TOTP -> Validación -> Redirección a Dashboard.

CRITERIO DE ÉXITO (DONE):
- Credenciales correctas + TOTP correcto = Entra a Dashboard y guarda sesión.
- Credenciales o código incorrectos = Muestra error claro en pantalla sin romperse.

```




### `BAC-01` - Script Docker Compose de base de datos local

- **Área:** Backend / Infraestructura
- **Asignados:** Tayra y Lucas
- **Estimación:** 2 h
- **Ventana propuesta:** 10/09/2026, 14:00-16:00
- **Depende de:** ninguna.
- **Entregable:** PostgreSQL local mediante Docker Compose, tablas mínimas de usuarios y roles `ADMIN`/`OPERATOR`, y usuario de prueba.
- **Criterio de éxito:** el equipo puede levantar la base con `docker compose up` y conectarse usando las credenciales documentadas.

### `BAC-02` - Modelo y hashing de contraseñas

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 10/09/2026, 14:00-16:00
- **Depende de:** ninguna; puede ejecutarse en paralelo con `BAC-01`.
- **Entregable:** modelo de usuario y verificación de contraseñas con bcrypt o Argon2.
- **Criterio de éxito:** la contraseña nunca se persiste en texto plano y una credencial válida se distingue de una inválida mediante pruebas automatizadas o manuales documentadas.

### `BAC-03` - Endpoint `POST /api/auth/login`

- **Área:** Backend
- **Asignado:** Tayra
- **Estimación:** 3 h
- **Ventana propuesta:** 11/09/2026, 09:00-12:00
- **Depende de:** `BAC-01` y `BAC-02`.
- **Entregable:** endpoint que recibe usuario o correo y contraseña, valida contra PostgreSQL y devuelve un JWT firmado con `id` y `role`.
- **Criterio de éxito:** el usuario de prueba puede iniciar sesión y el JWT se valida correctamente con la clave configurada.

### `BAC-04` - Manejo de errores HTTP

- **Área:** Backend
- **Asignado:** Lisandro
- **Estimación:** 1,5 h
- **Ventana propuesta:** 11/09/2026, 12:00-13:30
- **Depende de:** `BAC-03`.
- **Entregable:** respuestas JSON unificadas con `400` para payload incompleto y `401` para credenciales inválidas.
- **Criterio de éxito:** los casos válidos e inválidos devuelven código HTTP y estructura JSON documentados.

### `LOGIN-01` - Flujo 2FA/TOTP en memoria

- **Área:** Backend
- **Asignado:** Tayra y Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 11/09/2026, 13:30-15:30
- **Depende de:** `BAC-03` y `BAC-04`.
- **Entregable:** generación de secreto TOTP temporal, validación de código de seis dígitos y estado de verificación en memoria. No reemplaza todavía la persistencia segura definitiva del RF-01.
- **Criterio de éxito:** el login solicita el TOTP después de validar la contraseña y rechaza códigos inválidos o vencidos.

### `FRN-01` - Maquetado del formulario de login

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 3 h
- **Ventana propuesta:** 10/09/2026, 14:00-17:00
- **Depende de:** ninguna.
- **Entregable:** campos de usuario o correo, contraseña, botón de envío y estado visual de carga.
- **Criterio de éxito:** la vista se renderiza y permite recorrer los estados inicial, cargando, éxito y error.

### `FRN-02` - Validación en cliente del login

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 2 h
- **Ventana propuesta:** 10/09/2026, 17:00-19:00
- **Depende de:** `FRN-01`.
- **Entregable:** validación de formato y mensajes de error antes de enviar la petición.
- **Criterio de éxito:** no se envían formularios incompletos y cada error se muestra junto al campo correspondiente.

### `FRN-03` - Navbar y rutas base

- **Área:** Frontend
- **Asignados:** Cristian y Luz
- **Estimación:** 3 h
- **Ventana propuesta:** 10/09/2026, 14:00-17:00
- **Depende de:** ninguna.
- **Entregable:** navegación y vistas iniciales de Dashboard e Instancias, aunque sean estáticas.
- **Criterio de éxito:** el usuario puede navegar entre las pantallas base sin romper la aplicación.

### `FRN-04` - Conexión del login con localhost

- **Área:** Frontend
- **Asignado:** Cristian
- **Estimación:** 3 h
- **Ventana propuesta:** 11/09/2026, 15:30-18:30
- **Depende de:** `BAC-03`, `BAC-04` y `FRN-01`.
- **Entregable:** petición HTTP al backend local y almacenamiento del JWT en memoria.
- **Criterio de éxito:** el usuario de prueba puede iniciar sesión desde el frontend y el cliente conserva el JWT durante la sesión.

### `LOGIN-02` - Pantalla y validación del código TOTP

- **Área:** Frontend
- **Asignados:** Belinda y Luz
- **Estimación:** 2 h
- **Ventana propuesta:** 11/09/2026, 15:30-17:30
- **Depende de:** `LOGIN-01` y `FRN-01`.
- **Entregable:** pantalla de ingreso del código TOTP de seis dígitos y mensajes de error.
- **Criterio de éxito:** el frontend permite ingresar el código, muestra el estado de verificación y habilita la navegación solo cuando la validación es correcta.

### `LOGIN-03` - Integración y prueba del prototipo

- **Área:** Frontend / Backend
- **Asignados:** Cristian, Tayra y Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** 11/09/2026, 18:30-20:30
- **Depende de:** `BAC-04`, `FRN-04` y `LOGIN-02`.
- **Entregable:** recorrido completo usuario de prueba -> contraseña -> TOTP en memoria -> pantalla base.
- **Criterio de éxito:** el flujo válido termina en Dashboard, el flujo inválido muestra un error y no permite acceder a las vistas protegidas.