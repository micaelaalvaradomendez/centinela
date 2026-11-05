# Etapa 3: correcciones pendientes (FIX)

Acá se listan los fixes y mejoras técnicas correspondientes a las tareas de la Etapa 3 identificadas durante la verificación y ejecución de suites de aceptación.

---

## 🛠️ Especificación Detallada de Fixes (Etapa 3)

### `FIX-145` - Eliminación de la restricción de pertenencia rígida en borrado (`403 INSTANCE_NOT_OWNED`) y desacoplamiento mediante configuración de instancias protegidas (`RF-07` / `BAC-24B`) (Backend)

- **Área:** Backend
- **Asignado:** Lucas / Nico
- **Estimación:** 1.5 h
- **Depende de:** `RF-07`, `BAC-24B`, `PR #30` (`df4cada`).
- **Problema técnico detallado:**
  En el commit `df4cada` (PR #30), se introdujo en `instance_handler.go:EliminarInstancia` una consulta bloqueante `FueCreadaPorCentinela` contra la tabla `tareas_asincronas`, rechazando con `403 INSTANCE_NOT_OWNED` cualquier orden `DELETE` sobre instancias cuya tarea de creación no figure en la base de datos local.
  Esto genera un problema grave de desacople y operación en el servidor real:
  1. Impide a los administradores gestionar o destruir instancias preexistentes en Proxmox, instancias migradas o recreadas tras limpiezas de base de datos.
  2. En el servidor físico los nombres, identificadores y rangos de VMIDs cambian constantemente respecto a los entornos locales de desarrollo.
  3. La protección contra la destrucción accidental de las instancias propias del proyecto Centinela (el propio nodo Docker, DNS, base de datos del host) debe definirse explícitamente en la variable de entorno `PROXMOX_PROTECTED_VMIDS`, que ya responde de forma canónica con `403 INSTANCE_PROTECTED`.
- **Qué debe hacer el equipo de backend:**
  1. En `backend/internal/adapters/primary/http/instance_handler.go`:
     - Retirar de `EliminarInstancia` la guarda `FueCreadaPorCentinela` y la emisión del código `403 INSTANCE_NOT_OWNED`.
     - Mantener la verificación de estado operativo (`stopped` o `paused`) para responder `409 INSTANCE_INVALID_STATE` si la instancia está encendida.
     - Confiar la protección de infraestructura exclusivamente a la lista declarativa `PROXMOX_PROTECTED_VMIDS` en `RejectProtectedInstance(vmidsProtegidos)`.
  2. En `backend/docs/contrato-etapa1.md`:
     - Retirar la cláusula `INSTANCE_NOT_OWNED` del contrato y de la documentación de errores de `DELETE /instances/{instanceRef}`.
  3. En `backend/internal/core/ports/tarea_port.go` y `tarea_repository.go`:
     - Deprecar el método `FueCreadaPorCentinela` si no tiene otros consumidores funcionales.
- **Criterio de éxito:**
  Un usuario con rol `ADMIN` puede solicitar `DELETE` sobre cualquier instancia detenida del hipervisor que no forme parte de `PROXMOX_PROTECTED_VMIDS`, recibiendo `202 Accepted` y despachando el borrado en Proxmox VE.



