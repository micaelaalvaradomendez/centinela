# 6. Episodios Reales de Divergencia: La Sustitución Empírica del "PR #1234"

**Contexto en el Artículo:** Sección V.C (*"Verificación formal contra el código"*). Sustituye el ejemplo hipotético del `"PR #1234: divergencia en campo currency; falta manejo de error 403"` por los casos empíricos reales registrados en los repositorios de *El Centinela*.

---

## 6.1. El Caso Original en el Manuscrito vs. La Realidad Empírica

En la versión preliminar del artículo, en la Sección V.C, se expuso un escenario didáctico pero sintético:
> *"El PM no necesita revisar el código línea por línea: recibe un reporte sintético que dice 'PR #1234: divergencia en campo currency; falta manejo de error 403'. Con esa información, negocia con el equipo de backend la corrección antes de integrar. El tablero decía 'completado'; la auditoría dijo 'no cumple'. La IA hizo visible la brecha."*

El desarrollo del proyecto **El Centinela** ofrece **dos episodios reales documentados en Git**, con idéntica estructura conceptual pero con el valor inestimable de la evidencia primaria verificable.

---

## 6.2. Episodio Real 1: Divergencia de Contrato HTTP en el Borrado de Instancias (`BAC-24B` y Decisión `D2`)

### El Contexto y el Requerimiento:
En la Etapa 1, la tarea **`BAC-24B`** establecía la implementación del endpoint para destruir máquinas virtuales y contenedores LXC (`DELETE /api/instances/:vmid`).  
En la especificación de contratos ([`docs/contrato-etapa1.md`](file:///home/micaela/centinela/documentacion/etapa1.md)) y en las **Decisiones de Diseño `D2`**, se había acordado taxativamente la matriz canónica de errores:
*   Si una máquina virtual se encuentra encendida (`running`), la API debe rechazar la orden de borrado devolviendo un código **HTTP `409 Conflict`** con el cuerpo JSON:
    ```json
    {
      "codigo": "INSTANCE_INVALID_STATE",
      "mensaje": "La instancia debe estar detenida para ser eliminada"
    }
    ```
*   El frontend dependía de este código exacto para capturar el error y mostrar al operador un modal específico indicándole que debía apagar el recurso antes de destruirlo.

### La Divergencia en el Código:
El desarrollador de backend implementó la validación y marcó la tarea `BAC-24B` como **"Terminada"** en el tablero de ClickUp tras confirmar el commit `0850ec7`.  
Sin embargo, al implementar la condición en el controlador Go, escribió:
```go
// Fragmento defectuoso en backend/internal/adapters/primary/http/handlers/instance_handler.go
if status == "running" {
    c.JSON(http.StatusConflict, gin.H{
        "codigo": "INSTANCE_NOT_STOPPED", // <--- DIVERGENCIA RESPECTO AL CONTRATO D2
        "mensaje": "No se puede eliminar una instancia en ejecución"
    })
    return
}
```

### La Detección Automatizada:
El arnés de pruebas independiente en [`test/back/etapa1_acceptance_test.go`](file:///home/micaela/centinela/test/back/etapa1_acceptance_test.go) ejecutó la suite de validación contra el endpoint:

```
=== RUN   TestBAC24B_EliminacionInstanciaValidacionEstado
    etapa1_acceptance_test.go:184: [FAIL] DELETE /api/instances/101 con VM encendida:
        esperaba código de error: "INSTANCE_INVALID_STATE"
        obtuvo código de error:   "INSTANCE_NOT_STOPPED"
        HTTP Status: 409 Conflict
--- FAIL: TestBAC24B_EliminacionInstanciaValidacionEstado (0.42s)
```

### El Reporte y la Intervención del AI Manager:
La automatización bloqueó el merge hacia el repositorio orquestador y emitió el reporte de fallo documentado en [`test/back/RESULTADOS.md`](file:///home/micaela/centinela/test/back/RESULTADOS.md):
> **Dictamen del Test:** *BAC-24B DELETE: Avanzó en `0850ec7`: devuelve 202 { upid, tareaId } con seguimiento y TASK_FINISHED. Pero responde `409 INSTANCE_NOT_STOPPED` en lugar de `INSTANCE_INVALID_STATE` (violación de D2).*

**Acción del Gestor:** El PM no necesitó depurar el código de Go manualmente. Con el reporte objetivo en mano, intervino ante el desarrollador:
1.  Rechazó la promoción de `BAC-24B` a `terminado-1.md`.
2.  Generó el ticket atómico de corrección en ClickUp.
3.  El desarrollador unificó el código en el enum canónico, el test pasó a **PASS (verde)** y recién entonces se fusionó el puntero en el repositorio central.

---

## 6.3. Episodio Real 2: Falsa Remediación de IDs Ofuscados y Brecha de Inyección (`FIX-108`)

### El Contexto y el Requerimiento:
Tras la auditoría de seguridad, se emitió la directiva **`FIX-107` / `FIX-108`**: el sistema debía dejar de exponer VMIDs numéricos de Proxmox (como `101`) para erradicar vulnerabilidades IDOR, reemplazándolos por referencias opacas `InstanceRef` (`inst_...`).
*   Backend implementó la ofuscación en el commit `c51c77d`.
*   El equipo de frontend modificó las interfaces TypeScript para aceptar `ref: string | number` en el commit `b5466ac`, y declaró en la reunión de seguimiento que la tarea estaba *"100 % integrada y adaptada"*.

### La Divergencia en el Código:
En el frontend, la función de validación previa en [`frontend/centinela/src/components/features/instances/services/instanceService.ts`](file:///home/micaela/centinela/frontend/centinela/src/components/features/instances/services/instanceService.ts) fue implementada de forma laxa por el asistente de IA:
```typescript
// Código defectuoso implementado en el frontend
export function validateInstanceRef(ref: string | number): boolean {
    // Verificación superficial: no valida formato ni caracteres peligrosos
    return typeof ref === 'string' && ref.trim() !== '';
}
```

### La Detección Automatizada:
El arnés independiente ejecutó la suite [`test/front/fix-108-ids-ofuscados.test.tsx`](file:///home/micaela/centinela/test/front/fix-108-ids-ofuscados.test.tsx), arrojando **6 pruebas fallidas sobre 11 casos evaluados**:

```
FAIL test/front/fix-108-ids-ofuscados.test.tsx > FIX-108 Frontend Obfuscated IDs
  × rechaza referencias con directory traversal ../../instances/101 (esperaba false, obtuvo true)
  × rechaza referencias con espacios intermedios "inst 101" (esperaba false, obtuvo true)
  × rechaza caracteres no alfanuméricos "inst_101$#" (esperaba false, obtuvo true)
  × valida regex canónica ^inst_[a-zA-Z0-9_-]{12,64}$ (falló en cadenas de longitud incorrecta)

Tests: 6 failed, 5 passed, 11 total
```

### Diagnóstico de Causa Raíz en [`test/informe.md`](file:///home/micaela/centinela/test/informe.md):
> *"`FIX-108` — Adaptación del Frontend para consumo de identificadores ofuscados: Parcialmente implementada. validateInstanceRef en instanceService.ts:12-16 únicamente comprueba typeof ref === 'string' && ref.trim() !== '', dejando pasar inyecciones de ruta ('../') y omitiendo la regex canónica. Bloquea el merge hacia etapa de cierre."*

**Resolución del Gestor:** La automatización impidió que una validación superficial generara agujeros de seguridad en producción. El PM asignó inmediatamente el ticket **`FIX-124`** exigiendo la expresión regular estricta:
```typescript
const INSTANCE_REF_REGEX = /^(?:[1-9]\d{2,5}|inst_[a-zA-Z0-9_-]{12,64})$/;
```
Tras el commit `c35445f`, las 11 pruebas de `fix-108` y las 14 pruebas de `ids-ofuscados.test.tsx` pasaron al 100 % en verde.

---

## 6.4. Texto Propuesto para Reemplazar en la Sección V.C del Artículo

Para enriquecer directamente el manuscrito, se sugiere sustituir el párrafo hipotético del artículo por el siguiente texto redactado en estilo IEEE:

```markdown
La verificación automatizada contra contratos de interfaz demostró ser determinante en la práctica. 
Durante el desarrollo del proyecto *El Centinela* (plataforma de orquestación sobre Proxmox VE), 
se constató un caso paradigmático en la tarea BAC-24B (destrucción de instancias). 
El tablero de seguimiento registraba la funcionalidad como "completada" por el equipo de desarrollo. 
Sin embargo, la suite de pruebas de aceptación en Go detectó que, ante una máquina encendida, 
el endpoint respondía HTTP 409 con el código INSTANCE_NOT_STOPPED en lugar de INSTANCE_INVALID_STATE, 
violando el contrato de la API e impidiendo que el frontend desplegara el diálogo de advertencia. 
El arnés bloqueó automáticamente la integración y emitió un reporte estructurado de divergencia. 
Con este diagnóstico objetivo, el PM negoció la alineación del enumerador con el equipo de backend 
antes de fusionar el código a la rama principal. La auditoría formal hizo visible la brecha que el 
tablero administrativo ocultaba.
```
