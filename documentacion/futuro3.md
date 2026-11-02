# Etapa 3: correcciones pendientes (FIX)

Acá se listan los fixes y mejoras técnicas correspondientes a las tareas de la Etapa 3 identificadas durante la verificación y ejecución de suites de aceptación.

---

## 🛠️ Especificación Detallada de Fixes (Etapa 3)


### `FIX-130` - Documentación exhaustiva de storages compatibles con Snapshots y sanitización de secretos (`INF-13`) (Infraestructura / Documentación)

- **Área:** Infraestructura / Documentación
- **Asignado:** Nico / Tayra
- **Estimación:** 0.5 h
- **Depende de:** `INF-13`.
- **Problema técnico detallado:**
  En `documentacion/api-proxmox.md`:
  1. No se documenta la compatibilidad ni las restricciones de formato de disco de Proxmox VE para snapshots (`lvmthin`, `zfspool` y `qcow2` soportan snapshots en caliente; `raw` sobre LVM tradicional carece de soporte de snapshots y genera fallos `qm snapshot`).
  2. En la sección 1 (`.env`), la variable `PROXMOX_TOKEN_SECRET` expone un UUID real (`39e714a5-967a-46c2-8de9-af048ae0020f`) en texto plano en lugar de un placeholder seguro (ej. `<uuid-secreto-proxmox>`).
- **Qué debe hacer el equipo de infraestructura:**
  1. En `documentacion/api-proxmox.md`, reemplazar el token real por un placeholder genérico.
  2. Agregar una tabla o sección descriptiva detallando los tipos de storage aptos para snapshots (`lvmthin`, `zfspool`, `qcow2`) y advirtiendo sobre la incompatibilidad de `raw` sobre LVM estándar.
- **Criterio de éxito:**
  La prueba `TestEtapa2INF13_StorageCompatibleConSnapshots` en `test/back/etapa3_acceptance_test.go` valida la documentación sin errores.
---


---
