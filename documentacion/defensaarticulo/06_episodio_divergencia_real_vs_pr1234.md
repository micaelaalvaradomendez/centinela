# 6. Episodios Reales de Divergencia: La Sustitución Empírica del "PR #1234"

**Contexto en el Artículo:** Sección V.C (*"Verificación formal contra el código"*). Sustituye el ejemplo hipotético del `"PR #1234: divergencia en campo currency; falta manejo de error 403"` por los casos empíricos reales registrados en los repositorios de *El Centinela*.

---

## 6.1. El Caso Original en el Manuscrito vs. La Realidad Empírica

En la versión preliminar del artículo, en la Sección V.C, se expuso un escenario didáctico pero sintético:
> *"El PM no necesita revisar el código línea por línea: recibe un reporte sintético que dice 'PR #1234: divergencia en campo currency; falta manejo de error 403'. Con esa información, negocia con el equipo de backend la corrección antes de integrar. El tablero decía 'completado'; la auditoría dijo 'no cumple'. La IA hizo visible la brecha."*

El desarrollo del proyecto **El Centinela** ofrece **cuatro episodios reales documentados en Git** (dos en las Etapas 1-2 y dos en el cierre de la Etapa 3 el 02-03/11/2026), con idéntica estructura conceptual pero con el valor inestimable de la evidencia primaria verificable.

---

## 6.2. Episodio Real 1: Divergencia de Contrato HTTP en el Borrado de Instancias (`BAC-24B` y Decisión `D2`)

### El Contexto y el Requerimiento:
En la Etapa 1, la tarea **`BAC-24B`** establecía la implementación del endpoint para destruir máquinas virtuales y contenedores LXC (`DELETE /api/instances/:vmid`).  
En la especificación de contratos ([`backend/docs/contrato-etapa1.md`](file:///home/micaela/centinela/backend/docs/contrato-etapa1.md)) y en las **Decisiones de Diseño `D2`**, se había acordado taxativamente la matriz canónica de errores:
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

**Acción del Gestor:** La PM no necesitó depurar el código de Go manualmente. Con el reporte objetivo en mano:
1.  Rechazó la promoción de `BAC-24B` a `terminado-1.md`.
2.  Generó el ticket atómico de corrección en ClickUp.
3.  El desarrollador unificó el código en el enumerador canónico, el test pasó a **PASS (verde)** y recién entonces se fusionó el puntero en el repositorio central.

---

## 6.3. Episodio Real 2: Falsa Remediación de IDs Ofuscados y Brecha de Inyección (`FIX-108` / `FIX-124`)

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

**Resolución del Gestor:** La PM asignó el ticket **`FIX-124`** exigiendo la expresión regular estricta `INSTANCE_REF_REGEX = /^(?:[1-9]\d{2,5}|inst_[a-zA-Z0-9_-]{12,64})$/`. Tras el commit `c35445f`, las 11 pruebas de `fix-108` y las 14 pruebas de `ids-ofuscados.test.tsx` pasaron al 100 % en verde.

---

## 6.4. Episodio Real 3 (Etapa 3): Rechazo de Tráfico Nativo `vlan: 0` en Creación de Instancias (`FIX-128` Backend)

### El Contexto y la Divergencia:
En la creación de máquinas virtuales y contenedores (`POST /api/instances`, `RF-07`), cuando el operador no especifica una VLAN (tráfico nativo *untagged* según IEEE 802.1Q), los formularios o serializadores JSON envían `"vlan": 0`.
Sin embargo, en `backend/internal/adapters/primary/http/instance_handler.go`, el DTO `VLAN` era un puntero `*int` y el validador ejecutaba:
```go
if input.VLAN != nil && (*input.VLAN < 1 || *input.VLAN > 4094) {
    return ports.CrearInstanciaProxmoxInput{}, fmt.Errorf("vlan debe estar entre 1 y 4094")
}
```
Al recibir `0`, el puntero no era `nil` y la validación rechazaba la petición con `400 Bad Request`, bloqueando el aprovisionamiento de cualquier instancia sin VLAN etiquetada.

### Detección y Remediación (`TestFIX128VLAN_VlanCeroEsTraficoSinEtiquetar`):
La PM documentó el diagnóstico en [`actual.md`](file:///home/micaela/centinela/documentacion/actual.md#L89-L122), creó la prueba en [`test/back/etapa3_fixes_pendientes_test.go`](file:///home/micaela/centinela/test/back/etapa3_fixes_pendientes_test.go#L59-L84) y el equipo de backend la resolvió en el commit `b2d5fec` (*"fix(instancias): aceptar vlan 0 o null como red sin etiquetar"*), normalizando `*input.VLAN == 0` a `nil` y pasando los 3 subtests en verde.

---

## 6.5. Episodio Real 4 (Etapa 3): Quema Indebida de Intentos OTP ante Errores 400 de Contraseña (`FIX-132` Frontend)

### El Contexto y la Divergencia:
En el Caso de Uso **CU-12** (Recuperación de contraseña con OTP), la especificación establecía que si el usuario ingresa un código OTP válido pero elige una contraseña débil o igual a la anterior (`400 WEAK_PASSWORD` o `400 SAME_PASSWORD`), el sistema debe rechazar la clave **sin quemar los 3 intentos del código OTP**.
El backend cumplía el contrato, pero en [`frontend/centinela/src/pages/RecoverPassword.tsx`](file:///home/micaela/centinela/frontend/centinela/src/pages/RecoverPassword.tsx) se había programado:
```tsx
if (error.errorCode === 'RESET_FAILED' || error.status === 400) {
    const nextAttempts = failedAttempts + 1;
    setCodeError('El código ingresado es incorrecto o ha expirado.');
```
Al capturar cualquier `error.status === 400`, cuando un usuario ingresaba una clave repetida 3 veces, **el frontend acusaba falsamente al código OTP de ser inválido y bloqueaba la pantalla**, obligando a pedir un código nuevo por correo.

### Detección y Remediación (`etapa3-pendientes.test.tsx`):
La suite [`test/front/etapa3-pendientes.test.tsx`](file:///home/micaela/centinela/test/front/etapa3-pendientes.test.tsx#L97-L148) simuló 4 intentos consecutivos con `SAME_PASSWORD`, `WEAK_PASSWORD` y `PASSWORD_RECENTLY_USED`, exponiendo el bloqueo indebido en el 3.er intento. Tras el ticket **`FIX-132`**, Cristian corrigió la discriminación de errores en el commit `3d06236` (PR #179) y las 4 pruebas pasaron al 100 % en verde.

---

## 6.6. Texto Propuesto para Reemplazar en la Sección V.C del Artículo

Para enriquecer directamente el manuscrito, se sugiere sustituir el párrafo hipotético del artículo por el siguiente texto redactado en estilo IEEE:

```markdown
La verificación automatizada contra contratos de interfaz demostró ser determinante en la práctica. 
Durante el desarrollo del proyecto *El Centinela* (plataforma de orquestación sobre Proxmox VE), 
se constató un caso paradigmático en la tarea BAC-24B (destrucción de instancias). 
El tablero de seguimiento registraba la funcionalidad como "completada" por el equipo de desarrollo. 
Sin embargo, la suite de pruebas de aceptación en Go detectó que, ante una máquina encendida, 
el endpoint respondía HTTP 409 con el código INSTANCE_NOT_STOPPED en lugar de INSTANCE_INVALID_STATE, 
violando el contrato de la API e impidiendo que el frontend desplegara el diálogo de advertencia. 
Del mismo modo, en la tarea FIX-132 (recuperación de clave con OTP), las pruebas de aceptación 
detectaron que el frontend capturaba cualquier HTTP 400 (como SAME_PASSWORD) e incrementaba 
erróneamente el contador de fallos del código OTP, bloqueando al usuario al tercer intento. 
En ambos casos, el arnés bloqueó la integración y emitió un diagnóstico estructurado. Con esa 
evidencia objetiva, la PM negoció la corrección antes de actualizar el repositorio integrador: 
la auditoría formal hizo visible la brecha que el tablero administrativo ocultaba.
```
