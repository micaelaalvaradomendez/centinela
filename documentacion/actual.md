### 4. Desglose en ClickUp: Tareas chiquitas y puntuales (Paso a paso)

Para cumplir con la directiva de desglosar más el tablero y que nadie pueda escudarse en que una tarea es "demasiado grande" o "depende de otro", dividí las épicas en subtareas de 2 a 4 horas:

> [!NOTE]
> **Estado al 06/10/2026** (backend `eec77ff`, frontend `5a86dce`; detalle en [test/informe.md](../test/informe.md)).
> - **Pasaron a [`terminado.md`](terminado.md):**
>   - `FIX-29` y `FIX-40`, completas;
>   - `SEC-03`, `FIX-41`, `FIX-42` y `FIX-38`, completas;
>   - `FIX-36`, completa (PR #77, frontend `070e96b`);
>   - `FIX-37` y `FIX-39`, completas (backend `4e204f1`);
>   - `BAC-18B`, implementada con problemas en backend `4e204f1` (particionamiento trimestral OK, purga periódica OK; falta indexar `jti_access`, derivado a `FIX-44` en `futuro.md`);
>   - `FRN-17C`, implementada con problemas en frontend `5a86dce` (solicitud de ticket efímero y reconexión OK; falta redirección inmediata a `/login` ante 401, derivado a `FIX-45` en `futuro.md`).
> - **Pasaron a [`terminado-1.md`](terminado-1.md):**
>   - `INF-07B`, referencia local verificada (pendiente despliegue CT 103 en `INT-03`);
>   - `FRN-20A`, implementada con problemas en frontend `070e96b` (maquetado, IP y badges resueltos con `FIX-43` en frontend `5a86dce`, tests 5/5 en verde);
>   - `FRN-19A`, completa (PR #78, frontend `907efe5`, medidores de host en `ResourceMeter.tsx`);
>   - `BAC-25A`, completa (backend `eec77ff`, pool acotado `UPID_WORKERS` implementado y verificado, test en verde);
>   - `BAC-29`, implementada con problemas en backend `eec77ff` (contrato y Swagger publicados; alineación de especificación de permisos derivada a `FIX-46` en `futuro-1.md`);
>   - `BAC-22`, implementada con problemas en backend `eec77ff` (telemetría con Redis y stale reading implementados; recuperación ante caída Proxmox derivada a `FIX-47` en `futuro-1.md`);
>   - `BAC-23A`, implementada con problemas en backend `eec77ff` (adaptador de inventario y resolución de IP implementados; rutas literales en pruebas unitarias derivada a `FIX-48` en `futuro-1.md`);
>   - `FRN-17A`, implementada con problemas en frontend `5a86dce` (contexto y consumo de eventos implementados; deduplicación, conexión compartida y suscripción granular derivadas a `FIX-49` en `futuro-1.md`).
> - **Tareas que permanecen en este archivo:**
>   Todas las tareas asignadas de la Ola 1 y Fase Base han sido implementadas en los últimos commits y trasladadas a `terminado.md` o `terminado-1.md`. El trabajo pendiente activo se gestiona a través de los FIXes en [`futuro.md`](futuro.md) y [`futuro-1.md`](futuro-1.md).

---

Hito: Gestión Administrativa de Usuarios

Un administrador entra a /admin/users, ve la tabla real provista por GET /api/admin/users.  
Crea un usuario desde el modal (POST), el Back genera su clave temporal y must_change_password: true, y la tabla se actualiza.  
Modifica su rol o lo desactiva (PUT/DELETE).  
Si un usuario con rol OPERATOR intenta consultar estos endpoints o la vista, recibe un 403 Forbidden.  

> **Estado verificado (01/10/2026):** el recorrido completo funciona en ambos lados: tabla real, alta con confirmación, edición, baja y 403 al `OPERATOR`. El mensaje ante `502 EMAIL_DELIVERY_FAILED` en el alta quedó resuelto con `FIX-27` (en `terminado.md`). Verificado adicionalmente el 06/10/2026: baja lógica conserva unicidad (409 USER_CONFLICT ante reutilización de email/nombre) y alta de nuevos usuarios opera correctamente (201).

---

# ETAPA 1
> Las tareas de la Etapa 1 ya verificadas están en [`terminado-1.md`](terminado-1.md), y sus correcciones en [`futuro-1.md`](futuro-1.md).
---

> Todas las tareas de la Ola 1 de [`etapa1.md`](etapa1.md) (`BAC-29`, `BAC-22`, `BAC-23A`, `BAC-25A`, `FRN-17A`, `FRN-19A`, `FRN-20A`) están implementadas en los submódulos (verificado el 06/10/2026). Las completadas y las implementadas con observaciones fueron movidas a [`terminado-1.md`](terminado-1.md), y sus FIXes correctivos se encuentran registrados en [`futuro-1.md`](futuro-1.md).

