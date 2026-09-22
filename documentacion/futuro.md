# Plan de próximas tareas

## Criterio de planificación

Este plan parte del estado registrado en `actual.md`: ninguna tarea de `BAC-01` a `BAC-04` ni de `FRN-01` a `FRN-04` está confirmada como terminada. Por eso, primero se debe cerrar el prototipo básico de login con usuario de prueba, JWT y flujo 2FA/TOTP en memoria.

Las fechas y horarios siguientes son una propuesta de ejecución desde el jueves 10/09/2026. Las horas indicadas son horas de trabajo estimadas y las franjas permiten visualizar dependencias y tareas paralelas. Una tarea solo pasa a `terminado.md` cuando se valida su criterio de éxito.

### Decisiones de alcance

- El 2FA se configura por usuario y cada cuenta conserva su propio secreto TOTP cifrado. No forma parte del rol.
- `ADMIN` y `OPERATOR` determinan qué acciones puede realizar una cuenta. `user_instances` solo determina sobre qué instancias puede actuar un operador.
- Los permisos granulares por acción, por ejemplo permitir `start` y prohibir `delete`, no están modelados por `user_instances` ni son parte del RF-09 actual. Si se requieren, deben planificarse como una ampliación separada del modelo de autorización.
- Antes de completar el cambio de contraseña y la verificación o vinculación 2FA, el backend solo debe emitir una sesión o token restringido para continuar el proceso de autenticación, nunca un JWT con acceso completo a la aplicación.
- **RF-12 (organizaciones) queda descartado del MVP.** El sistema se planifica y se implementa como monoorganización; ninguna tarea de la fase base ni de las etapas posteriores debe reservar campos de `organization_id`.
- **RF-13 (recuperación de contraseña por el propio usuario) se suma a la fase base**, para no dejar el flujo de recuperación cubierto solo con el reset administrativo (`BAC-15`).
- **La auditoría (RF-08) es una base transversal, no una tarea aislada de la Etapa 3.** Desde `BAC-18` la tabla `audit_logs` debe diseñarse como append-only (sin `UPDATE`/`DELETE` habilitados para la aplicación) y con columnas genéricas (`resource_type`, `resource_id`, `upid` nullable) para que la Etapa 1 (energía de instancias), la Etapa 2 (aprovisionamiento) y la Etapa 3 (snapshots) reutilicen la misma tabla sin migrar el esquema ni reconstruir historial. Toda tarea que agregue una acción auditable en cualquier etapa debe registrar contra esta tabla, no crear una paralela.
- **El contrato del canal de eventos/notificaciones (base de RF-11) se define en la fase base**, aunque el motor de alertas completo se construya recién en la Etapa 3, para que el poller de UPID de la Etapa 1 y las alertas de saturación de la Etapa 3 compartan el mismo esquema de evento desde el principio (ver `BAC-21`).
- **Servidor SMTP real descartado:** La integración con un proveedor SMTP real queda formalmente descartada del alcance. La entrega de correos y credenciales temporales queda resuelta de forma definitiva y completa mediante `MockEmailService` por consola a través del puerto `EmailService` (`BAC-16`), cumpliendo con mínimos privilegios (Zero-Trust), logs sanitizados y control estructurado de fallos sin dependencias de infraestructura externa.

## Fase 0. Cierre del prototipo de login


### `FIX-08` - Códigos de error desalineados entre frontend y backend

- **Área:** Frontend / Backend
- **Asignados:** Cristian y Lisandro
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (junto con `FIX-07`).
- **Depende de:** `FIX-07`.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 6, hallazgo H2):** además del nombre del campo (`FIX-07`), los valores concretos de los códigos de error que espera el frontend no coinciden todos con los que emite el backend.
- **Entregable:** acordar y documentar la lista cerrada de `errorCode` que puede emitir el backend, y actualizar el frontend para manejar exactamente esos valores.
- **Criterio de éxito:** cada `errorCode` que devuelve el backend tiene un manejo explícito en el frontend, sin casos que caigan en el mensaje genérico por desconocer el código.

### `FIX-09` - `/signup` sin endpoint de backend (404)

- **Área:** Frontend / Backend
- **Asignados:** Cristian y Tayra
- **Estimación:** 1 h (decisión) + implementación si corresponde
- **Ventana propuesta:** A definir.
- **Depende de:** ninguna.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 6, hallazgo H4):** el frontend tiene una vista o llamada a `/signup` que no tiene endpoint correspondiente en el backend y responde `404`.
- **Entregable:** decidir si el autorregistro forma parte del alcance (no lo es: `RF-09` dice que los usuarios los crea un administrador) y, si no lo es, quitar la ruta/llamada del frontend; si se decide que sí, crear el endpoint correspondiente.
- **Criterio de éxito:** no queda ninguna ruta del frontend apuntando a un endpoint inexistente; la decisión de alcance queda reflejada en este documento.

### `FIX-10` - Test de `FRN-03` desactualizado (usa `localStorage` en vez de `sessionStorage`)

- **Área:** Frontend / Testing
- **Asignados:** Cristian y Luz
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir.
- **Depende de:** ninguna.
- **Problema (evidencia en `ANALISIS-READINESS.md`, sección 5):** `navigation.test.tsx` escribe el token en `localStorage`, pero la app lo lee de `sessionStorage` (`src/storage/tokenStorage.ts`) desde el commit `432d747`. El producto funciona; el test quedó desactualizado y falla.
- **Entregable:** actualizar `navigation.test.tsx` para que escriba el token en `sessionStorage`.
- **Criterio de éxito:** `FRN-03` pasa en `test/front` sin tocar el código de producto.
- **Propuesta de mejora de seguridad:** Se implementó el almacenamiento de los JWT mediante `sessionStorage`. Como mejora de seguridad, se sugiere evaluar el uso de cookies `HttpOnly`, principalmente para el `refreshToken`, ya que evita que JavaScript pueda acceder directamente a este token y reduce su exposición ante ataques XSS. Esto requeriría adaptar el manejo del refresh en el backend (ver `siguientesetapas.md`).

### `FIX-11` - Re-verificación completa tras el fix crítico

- **Área:** Backend / Frontend
- **Asignados:** Cristian, Tayra y Lisandro
- **Estimación:** 1 h
- **Ventana propuesta:** A definir (posterior a `FIX-01` y `FIX-03`).
- **Depende de:** `FIX-01` y `FIX-03`.
- **Problema:** con el login roto por `FIX-01`, `BAC-06`, `BAC-06B`, `BAC-09`, `BAC-10`, `BAC-11` y `LOGIN-03` no se pudieron verificar de punta a punta, aunque figuren como cerradas en `terminado.md`.
- **Entregable:** correr `test/back` completo contra el commit con los fixes aplicados y actualizar el estado real de cada tarea afectada en `terminado.md` (o devolverla a este documento si el criterio de éxito no se cumple).
- **Criterio de éxito:** `test/back` no reporta ningún fallo de cascada originado en el login; cada tarea que figura en `terminado.md` tiene su criterio de éxito confirmado por la suite, no solo por inspección de código.

### `FIX-12` - Falta guard de rol administrativo en rutas de gestión de usuarios (`FRN-05`)

- **Área:** Frontend
- **Asignada:** Belinda / Cristian
- **Estimación:** 0,5 h
- **Ventana propuesta:** A definir (alta prioridad).
- **Depende de:** `FIX-08`.
- **Problema (evidencia en `test/front/admin-users.test.tsx`):** en `applicationRoutes.tsx`, las rutas administrativas `/users`, `/users/new` y `/users/:userId` solo cuentan con `loader: loadProtectedSession` (que comprueba que exista token JWT), pero no verifican el rol del usuario autenticado. Esto permite que un usuario con rol `OPERATOR` ingrese al panel de administración de usuarios en lugar de ser redirigido a `/dashboard`, incumpliendo el requerimiento `RF-09` y el criterio de éxito de `FRN-05` ("el administrador ve el listado real y un operador no puede acceder a la vista"). El guard requerido (`loadAdminSession`) ya existe en `sessionGuard.ts`, pero no está enlazado a las rutas.
- **Entregable:**
  1. Importar `loadAdminSession` en `src/routes/applicationRoutes.tsx`.
  2. Asignar `loader: loadAdminSession` a las rutas `/users`, `/users/new` y `/users/:userId` (o agruparlas bajo un layout/sub-ruta administrativa con dicho guard).
- **Criterio de éxito:** un usuario con rol `OPERATOR` que navega hacia `/users` es redirigido inmediatamente a `/dashboard`; la prueba `un usuario con rol OPERATOR no debe poder ver el panel de administración` de `test/front/admin-users.test.tsx` pasa exitosamente.

### `FIX-13` - Desalineación de endpoints y PR pendiente en recuperación de contraseñas (`BAC-19` / `BAC-20`)

- **Área:** Backend
- **Asignados:** Lisandro y Tayra
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (alta prioridad).
- **Depende de:** rama remota `origin/feat/gestion-credenciales`.
- **Problema (evidencia en `test/back/password_recovery_acceptance_test.go`):** los endpoints de recuperación de contraseña devuelven `404 page not found`. La auditoría técnica del repositorio reveló que el código fue implementado en la rama remota `origin/feat/gestion-credenciales`, pero nunca fue integrado a `main` vía Pull Request. Además, en dicha rama los endpoints se nombraron de forma inconsistente con la especificación de `RF-13` y `actual.md` (`/api/auth/password/forgot` y `/api/auth/password/reset` en lugar de `/api/auth/password-recovery/request` y `/api/auth/password-recovery/confirm`).
- **Entregable:**
  1. Alinear en la rama `feat/gestion-credenciales` los paths de los endpoints al contrato oficial de `actual.md` (`POST /api/auth/password-recovery/request` y `POST /api/auth/password-recovery/confirm`), o bien exponer alias de compatibilidad en `main.go`.
  2. Abrir PR de `origin/feat/gestion-credenciales` a `main`, validar que no rompa regresiones e integrarlo.
- **Criterio de éxito:** las peticiones a `POST /api/auth/password-recovery/request` y `POST /api/auth/password-recovery/confirm` responden según contrato; las pruebas de `password_recovery_acceptance_test.go` pasan en verde sobre `main`.

### `FIX-15` - Implementación de cookies HttpOnly para refreshToken (Mitigación XSS)

- **Área:** Frontend / Backend
- **Asignados:** Cristian y Lisandro
- **Estimación:** 2 h
- **Ventana propuesta:** A definir.
- **Depende de:** `FIX-10`.
- **Problema / Propuesta:** Se implementó el almacenamiento de los JWT mediante `sessionStorage`. Como mejora de seguridad, se debe evaluar e implementar el uso de cookies `HttpOnly`, principalmente para el `refreshToken`, ya que evita que JavaScript pueda acceder directamente a este token y reduce su exposición ante ataques XSS. Esto requeriría adaptar el manejo del refresh en el backend.
- **Entregable:**
  1. Backend: configurar `/api/auth/refresh` y login para emitir el `refreshToken` mediante cookie `Set-Cookie` con flags `HttpOnly`, `Secure` y `SameSite=Strict`, y consumirlo desde la cookie en lugar de body JSON.
  2. Frontend: adaptar el interceptor HTTP y cliente de autenticación para enviar credenciales automáticamente (`credentials: 'include'`) sin almacenar `refreshToken` en memoria/storage accesible por JS.
- **Criterio de éxito:** el `refreshToken` no es accesible desde JavaScript en el cliente; las peticiones de renovación de sesión son transparentes y seguras mediante la cookie HttpOnly.

## Fase 2. Cierre de autenticación y asignación de máquinas

*(Nota de alcance: `BAC-16` se encuentra finalizada y validada en `terminado.md` con `MockEmailService`; la integración con un servidor SMTP real quedó formalmente descartada).*

---

*(Nota: `BAC-19`, `BAC-20` y `BAC-21` fueron completadas, integradas y verificadas 100% en `terminado.md`).*



### `INF-05` - CORS y TLS en el borde

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1,5 h
- **Ventana propuesta:** A definir (junto con `INF-04`).
- **Depende de:** `INF-04`.
- **Entregable:** whitelist de orígenes permitidos en CORS y certificados TLS configurados en Nginx para todo el tráfico hacia el frontend y la API.
- **Criterio de éxito:** una petición desde un origen no autorizado es rechazada por CORS y el tráfico hacia el sistema se sirve únicamente sobre HTTPS.

### `LOGIN-04` - Prueba integral de autenticación y autorización

- **Área:** Frontend / Backend
- **Asignados:** Cristian, Tayra y Lisandro
- **Estimación:** 3 h
- **Ventana propuesta:** 18/09/2026, 14:00-17:00
- **Depende de:** `FRN-07`, `FRN-09`, `FRN-10`, `FRN-11`, `FRN-12`, `BAC-14`, `BAC-16`, `BAC-17`, `BAC-18`, `BAC-20` y `BAC-21`.
- **Entregable:** pruebas documentadas o automatizadas de creación de usuario, entrega y cambio de clave temporal, enrolamiento y login con TOTP, acceso según rol, filtro por instancias, recuperación administrativa de contraseña y 2FA, recuperación de contraseña por el propio usuario, logout/revocación de sesión y verificación de que las acciones administrativas quedan auditadas.
- **Criterio de éxito:** todos los recorridos válidos terminan con el acceso esperado y los intentos de omitir pasos, usar credenciales anteriores, acceder con otro rol, consultar una instancia no asignada o reutilizar un token revocado son rechazados con códigos HTTP controlados.

## Fase 3. Despliegue de persistencia y red




## Relación y orden de ejecución

```mermaid
flowchart TD
	BAC01[BAC-01 BD local] --> BAC03[BAC-03 Login JWT]
	BAC02[BAC-02 Hashing] --> BAC03
	BAC03 --> BAC04[BAC-04 Errores HTTP]
	BAC04 --> LOGIN01[LOGIN-01 TOTP en memoria]
	FRN01[FRN-01 Login visual] --> FRN02[FRN-02 Validación cliente]
	BAC03 --> FRN04[FRN-04 Conexión localhost]
	FRN01 --> FRN04
	LOGIN01 --> LOGIN02[LOGIN-02 Pantalla TOTP]
	FRN01 --> LOGIN02
	BAC04 --> LOGIN03[LOGIN-03 Integración login]
	FRN04 --> LOGIN03
	LOGIN02 --> LOGIN03
	LOGIN03 --> BAC05[BAC-05 Esquema RBAC]
	BAC05 --> BAC06[BAC-06 CRUD usuarios]
	BAC05 --> BAC07[BAC-07 Permisos por recurso]
	BAC07 --> BAC08[BAC-08 Guard por recurso]
	BAC05 --> BAC10[BAC-10 Setup 2FA]
	BAC10 --> BAC11[BAC-11 Persistencia y validación TOTP]
	BAC11 --> FRN09[FRN-09 Vinculación QR]
	BAC06 --> BAC12[BAC-12 Cambio de contraseña]
	BAC12 --> FRN10[FRN-10 Cambio obligatorio]
	BAC11 --> BAC13[BAC-13 Reset 2FA]
	BAC08 --> BAC13
	BAC08 --> BAC14[BAC-14 Inventario Proxmox]
	BAC06 --> BAC15[BAC-15 Reset contraseña]
	BAC12 --> BAC15
	BAC15 --> BAC16[BAC-16 Entrega credenciales]
	BAC13 --> FRN11[FRN-11 Recuperación admin]
	BAC15 --> FRN11
	LOGIN03 --> FRN05[FRN-05 Panel usuarios]
	BAC06 --> FRN05
	FRN05 --> FRN06[FRN-06 Alta y baja]
	BAC07 --> FRN07[FRN-07 Selector instancias]
	BAC14 --> FRN07
	FRN05 --> FRN07
	BAC08 --> FRN08[FRN-08 Interceptor 403]
	FRN07 --> LOGIN04[LOGIN-04 Prueba integral]
	FRN09 --> LOGIN04
	FRN10 --> LOGIN04
	FRN11 --> LOGIN04
	BAC16 --> LOGIN04
	BAC05 --> INF03[INF-03 PostgreSQL servidor]
	INF03 --> INF04[INF-04 Nginx y vmbr1]
```

### Tareas que pueden ejecutarse en paralelo

- `BAC-01`, `BAC-02`, `FRN-01` y `FRN-03` no dependen entre sí.
- `FRN-05` puede comenzar con el contrato definido de `BAC-06`, aunque necesita el endpoint para validación final.
- `FRN-06` puede avanzar con datos simulados mientras termina `BAC-06`.
- `BAC-10`/`BAC-11`, `BAC-12` y `BAC-14` pueden desarrollarse en paralelo una vez disponibles sus dependencias.
- `FRN-09` y `FRN-10` pueden maquetarse con contratos simulados, pero requieren `BAC-11` y `BAC-12` para validar los bloqueos de navegación.
- `BAC-13` y `BAC-15` pueden desarrollarse en paralelo; `FRN-11` integra ambas recuperaciones una vez definidos sus contratos.
- `INF-03` debe esperar el esquema de `BAC-05`, pero puede preparar el LXC, la red y los secretos antes de desplegarlo.

### Secuencia crítica

`BAC-01` + `BAC-02` -> `BAC-03` -> `BAC-04` -> `LOGIN-01` + `LOGIN-02` -> `FRN-04` -> `LOGIN-03` -> `BAC-05` -> `BAC-06`/`BAC-07` -> `BAC-08` -> `BAC-10`/`BAC-11` + `BAC-12` + `BAC-14` -> `BAC-13`/`BAC-15` + `FRN-07`/`FRN-09`/`FRN-10` -> `BAC-16`/`FRN-11` -> `LOGIN-04`.

## Resumen para ClickUp

| ID | Área | Asignado | Horas | Fecha propuesta | Dependencias |
|---|---|---|---:|---|---|
| BAC-01 | Backend/Infra | Tayra / Lucas | 2 | 10/09/2026 | - |
| BAC-02 | Backend | Lisandro | 2 | 10/09/2026 | - |
| BAC-03 | Backend | Tayra | 3 | 11/09/2026 | BAC-01, BAC-02 |
| BAC-04 | Backend | Lisandro | 1,5 | 11/09/2026 | BAC-03 |
| LOGIN-01 | Backend | Tayra / Lisandro | 2 | 11/09/2026 | BAC-03, BAC-04 |
| FRN-01 | Frontend | Belinda | 3 | 10/09/2026 | - |
| FRN-02 | Frontend | Belinda | 2 | 10/09/2026 | FRN-01 |
| FRN-03 | Frontend | Cristian / Luz | 3 | 10/09/2026 | - |
| FRN-04 | Frontend | Cristian | 3 | 11/09/2026 | BAC-03, BAC-04, FRN-01 |
| LOGIN-02 | Frontend | Belinda / Luz | 2 | 11/09/2026 | LOGIN-01, FRN-01 |
| LOGIN-03 | Integración | Cristian / Tayra / Lisandro | 2 | 11/09/2026 | FRN-04, LOGIN-02 |
| BAC-05 | Backend | Tayra / Lisandro | 3 | 14/09/2026 | BAC-01, LOGIN-03 |
| BAC-06 | Backend | Lisandro | 4 | 14/09/2026 | BAC-05 |
| BAC-07 | Backend | Tayra | 3 | 15/09/2026 | BAC-05, BAC-06 |
| BAC-08 | Backend | Lisandro | 3 | 15/09/2026 | BAC-07 |
| FRN-05 | Frontend | Belinda | 4 | 14/09/2026 | LOGIN-03, BAC-06 |
| FRN-06 | Frontend | Luz | 3 | 14/09/2026 | FRN-05, BAC-06 |
| FRN-07 | Frontend | Cristian | 4 | 15/09/2026 | FRN-05, BAC-07, BAC-14 |
| FRN-08 | Frontend | Cristian | 2 | 15/09/2026 | BAC-08, FRN-04 |
| BAC-10 | Backend | Tayra | 2 | 16/09/2026 | BAC-05, LOGIN-03 |
| BAC-11 | Backend | Lisandro | 3 | 16/09/2026 | BAC-10, BAC-05 |
| FRN-09 | Frontend | Belinda / Luz | 3 | 16/09/2026 | BAC-10, BAC-11, LOGIN-02 |
| BAC-12 | Backend | Lisandro | 2 | 16/09/2026 | BAC-02, BAC-05, BAC-06 |
| FRN-10 | Frontend | Belinda | 2 | 16/09/2026 | BAC-12, FRN-04 |
| BAC-13 | Backend | Tayra | 2 | 17/09/2026 | BAC-06, BAC-08, BAC-11 |
| BAC-14 | Backend | Tayra | 3 | 17/09/2026 | BAC-08, credenciales Proxmox |
| BAC-15 | Backend | Lisandro | 2 | 17/09/2026 | BAC-06, BAC-12 |
| BAC-16 | Backend/Infra | Lisandro / Nico | 3 | 18/09/2026 | BAC-06, BAC-15, SMTP |
| FRN-11 | Frontend | Luz | 2 | 18/09/2026 | FRN-05, BAC-13, BAC-15 |
| LOGIN-04 | Integración | Cristian / Tayra / Lisandro | 3 | 18/09/2026 | FRN-07, FRN-09, FRN-10, FRN-11, BAC-14, BAC-16 |
| INF-03 | Infraestructura | Lucas / Nico | 3 | 16/09/2026 | BAC-05 |
| INF-04 | Infraestructura | Nico | 3 | 16/09/2026 | INF-03 |
| BAC-17 | Backend | Lisandro | 2 | A definir | BAC-03, BAC-05 |
| BAC-18 | Backend | Tayra | 4 | A definir | BAC-05 |
| BAC-19 | Backend | Lisandro | 2 | A definir | BAC-05, BAC-16 |
| BAC-20 | Backend | Tayra | 2 | A definir | BAC-19, BAC-17 |
| BAC-21 | Backend | Tayra / Lisandro | 2 | A definir | BAC-05 |
| FRN-12 | Frontend | Belinda | 3 | A definir | BAC-19, BAC-20 |
| INF-05 | Infraestructura | Nico | 1,5 | A definir | INF-04 |

**Resultado esperado:** al completar `LOGIN-04` e `INF-04`, el equipo tendrá cerrado el circuito técnico de RF-01 y RF-09: contraseña temporal entregada y reemplazada, 2FA individual persistido, JWT posterior a la verificación, usuarios, roles, recuperación administrativa, permisos binarios por instancia, inventario real y bloqueo por recurso. Con `BAC-17` a `BAC-21`, `FRN-12` e `INF-05` también queda cerrado RF-13 (recuperación de contraseña por el propio usuario), la revocación de sesiones, la base append-only de auditoría (parte de RF-08) y el contrato de eventos que la Etapa 1 y la Etapa 3 van a reutilizar (base de RF-11), sin dejar backfill ni rediseño pendiente. La protección de las operaciones de Proxmox debe validarse antes de habilitar acciones destructivas.

**Fuera de alcance:** no se incluyen permisos granulares por acción porque el requerimiento actual solo define roles y asignación de instancias. Incorporarlos requiere una ampliación explícita del modelo de autorización y nuevos criterios para cada operación. **RF-12 (organizaciones/multi-tenant) queda descartado del MVP** y no debe planificarse en ninguna etapa posterior salvo decisión explícita en contrario.

### 🧱 Bloque 3: Persistencia Definitiva de 2FA y Ciclo de Vida de Credenciales

Este bloque formaliza la seguridad completa (Fase 2) y cierra el hito integral `LOGIN-04`.

#### 1. Tareas de desarrollo paralelo

* **Backend:**
* `BAC-10` y `BAC-11` — Enrolamiento 2FA, QR (`GET /api/auth/2fa/qr`) y persistencia del secreto cifrado con AES-256 (`POST /api/auth/2fa/verify`).


* `BAC-12` — Cambio obligatorio de contraseña temporal (`POST /api/auth/change-password`).


* `BAC-13` y `BAC-15` — Restablecimiento administrativo de 2FA y contraseña (`POST /api/admin/users/{id}/reset-...`).


* `BAC-16` — Entrega segura de credenciales temporales vía SMTP.




* **Frontend:**
* `FRN-09` — Modal obligatorio de vinculación 2FA con escaneo de QR y confirmación.


* `FRN-10` — Pantalla obligatoria de cambio de contraseña cuando `must_change_password: true`.


* `FRN-11` — Acciones de restablecimiento de contraseña y 2FA desde el panel de administración.





#### 🎯 Prueba de Integración final (`LOGIN-04`):

> **Hito:** *Circuito Completo de Seguridad y Autenticación*
> 1. El Admin crea un usuario -> El usuario recibe su clave temporal por correo (`BAC-16`).
> 
> 
> 2. El usuario inicia sesión -> El Front detecta `must_change_password` y lo bloquea hasta que cambie su clave temporal (`FRN-10` / `BAC-12`).
> 
> 
> 3. Pasa al enrolamiento de 2FA: escanea el QR en Google Authenticator (`FRN-09` / `BAC-10`) y confirma su código (`BAC-11`).
> 
> 
> 4. Ingresa al sistema con su JWT definitivo según su rol y permisos de instancias.
> 
> 
> 5. El Admin puede revocarle el 2FA o la clave (`FRN-11` / `BAC-13` / `BAC-15`) forzándolo a reiniciar el ciclo en su próximo acceso.
> 
> 
> 
> 

---

### 📋 Resumen del orden de ejecución para el equipo:

1. **Cerrar Bloque 0:** Resolver el fix de dígitos OTP en `TwoFactor.tsx` (`LOGIN-02` / `LOGIN-03` al 100%).


2. **Asignar Bloque 1:** CRUD de Usuarios y Roles (`BAC-05`, `BAC-06`, `BAC-06B`, `BAC-09` + `FRN-05`, `FRN-06`, `FRN-06B`) -> **Prueba de Integración 1**.


3. **Asignar Bloque 2:** Permisos por Recurso e Inventario (`BAC-07`, `BAC-08`, `BAC-14` + `FRN-07`, `FRN-08`) -> **Prueba de Integración 2**.


4. **Asignar Bloque 3:** 2FA Real, Recuperación y Cambio de Clave (`BAC-10` a `BAC-16` + `FRN-09` a `FRN-11`) -> **Prueba de Integración `LOGIN-04**`.