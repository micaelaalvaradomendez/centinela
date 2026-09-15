
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




