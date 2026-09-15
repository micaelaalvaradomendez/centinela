
### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:



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

---

