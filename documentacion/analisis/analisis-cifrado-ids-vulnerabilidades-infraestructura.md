# Auditoría Técnica: Cifrado de Identificadores, Fuga de Metadatos y Blindaje de Infraestructura

**Fecha:** 19 de Octubre de 2026  
**Área:** Ciberseguridad, Arquitectura de Backend y Seguridad de la Infraestructura  
**Archivos y Submódulos Auditados:**  
*   `backend` (`centinela-back`, commits `8daaf87`, `bd22c3d`, `1f63959`)
*   `frontend` (`centinela-front`, commit `91eabf4`)
*   `documentacion/futuro2.md` (`FIX-105` al `FIX-110`)
*   Suites de integración en `test/back/` y `test/front/`

---

## 1. Introducción y Contexto del Análisis

El presente informe documenta la investigación técnica realizada tras la revisión del código fuente del backend y la consulta mantenida con el equipo de desarrollo acerca de la exposición de identificadores técnicos (`UPID` de Proxmox VE, `VMID` de instancias, identificadores de tareas y nodos).

### Pregunta Inicial y Premisa de Seguridad:
Durante el desarrollo del requerimiento **RF-08** (Auditoría inmutable de aprovisionamiento y cambios de recursos), surgió la duda sobre si el `UPID` de Proxmox debía enviarse al frontend. La premisa acordada establecía que **el frontend no debía recibir identificadores que revelen la infraestructura real o que permitan vulnerar el servidor donde reside la plataforma en caso de que la aplicación web sea comprometida**.

A continuación, se detalla:
1. Cómo se detectó el estado real de la encriptación en el backend.
2. Dónde reside y dónde falta dicha lógica.
3. Qué vulnerabilidades concretas se presentan ante identificadores sin cifrar/ofuscar.
4. Por qué cada una de ellas es crítica para la supervivencia del servidor anfitrión.

---

## 2. Cómo se Detectaron los Detalles de Encriptación en el Código

Para verificar fehacientemente qué mecanismos de cifrado estaban operando, se ejecutó una inspección exhaustiva de código estático, historial de git y dependencias criptográficas en el backend.

### 2.1. Criptografía Realmente Existente en el Backend
Al auditar el paquete [`backend/internal/infrastructure/crypto`](file:///home/micaela/centinela/backend/internal/infrastructure/crypto):
*   **Contraseñas:** En [`password.go`](file:///home/micaela/centinela/backend/internal/infrastructure/crypto/password.go), se utiliza **bcrypt** (`golang.org/x/crypto/bcrypt`) con factor de costo 12 (`Cost: 12`).
*   **Secretos 2FA / TOTP:** En [`totp.go`](file:///home/micaela/centinela/backend/internal/infrastructure/crypto/totp.go#L13-L45), existe cifrado simétrico real **AES-256-GCM** mediante la función `CifrarSecreto(raw, key)` y `DescifrarSecreto(enc, key)` con la variable de entorno `TOTP_ENCRYPTION_KEY` (32 bytes).
*   **Tokens de Sesión:** En [`jwt.go`](file:///home/micaela/centinela/backend/internal/infrastructure/crypto/jwt.go), firma criptográfica HMAC-SHA256 (`HS256`).

> ⚠️ **Hallazgo 1:** **No existe ninguna función de cifrado, encriptación reversible o hashing para IDs de entidades.**  
> Los identificadores de bases de datos, máquinas virtuales, tareas y usuarios se almacenan y transmiten sin ningún procesamiento criptográfico.

---

### 2.2. Qué Hizo Realmente el Desarrollador (Análisis de los Commits `bd22c3d` y `1f63959`)

Al revisar el historial de Git del backend:
```bash
commit bd22c3d00b6c0d6838b847bc558e18c0f4466774
Author: tayraag <tayraag07@gmail.com>
Date:   Sun Oct 18 21:07:44 2026 -0300
Subject: fix(seguridad): ocultar UPID tras el identificador de tarea
```

Se constató que:
1.  **No hubo encriptación:** El desarrollador no aplicó AES ni ningún algoritmo de cifrado sobre el `UPID`.
2.  **Sustitución por ID Opaco:** En vez de encriptar el UPID, implementó un **enmascaramiento por base de datos**. Al iniciar una acción en Proxmox:
    *   Genera un UUIDv4 aleatorio (`tareaId`) en la tabla `tareas_asincronas`.
    *   Almacena el UPID privado de Proxmox en la columna interna `tareas_asincronas.upid_proxmox`.
    *   Modificó el DTO de respuesta [`AccionAceptadaResponse`](file:///home/micaela/centinela/backend/internal/adapters/primary/http/contrato_etapa1.go#L50-L55) para responder únicamente:
        ```json
        { "tareaId": "3f2504e0-4f89-11d3-9a0c-0305e82c3301" }
        ```
    *   Eliminó el campo `upid` tanto del 202 como de la tabla visible de auditoría ([`audit_port.go`](file:///home/micaela/centinela/backend/internal/core/ports/audit_port.go#L43)).

---

### 2.3. Detección de la Ruptura Inmediata en el Frontend y Tests

Al contrastar este cambio del backend con el resto del repositorio:

1.  **Frontend Roto en [`instanceService.ts`](file:///home/micaela/centinela/frontend/centinela/src/components/features/instances/services/instanceService.ts#L88):**
    En las funciones de ciclo de vida (`requestInstancePowerAction`, `deleteInstance` y `createInstance`), el frontend tenía validaciones defensivas rígidas que exigían la presencia obligatoria de `upid`:
    ```typescript
    // instanceService.ts:88-90
    if (!response || typeof response.upid !== 'string' || !response.upid.trim()) {
        throw new ApiRequestError('La orden fue aceptada, pero no se pudo interpretar su seguimiento...');
    }
    ```
    Al dejar de recibir el campo `upid`, **todas las órdenes en la interfaz gráfica (start, stop, reboot, delete, create) lanzan un error bloqueante en pantalla**.

2.  **Tests de Integración Rotos en [`test/back/int01_acceptance_test.go:58`](file:///home/micaela/centinela/test/back/int01_acceptance_test.go#L58-L60):**
    La suite de aceptación exige expresamente:
    ```go
    upid := requiredString(t, body, "upid")
    if !strings.HasPrefix(upid, "UPID:") {
        t.Errorf("campo upid debe iniciar con 'UPID:', recibido %q", upid)
    }
    ```
    Falla al no estar sincronizada con el nuevo contrato de seguridad.

---

## 3. Matriz de Identificadores: Vulnerabilidades Detectadas y por qué son Críticas

Si un atacante vulnera o compromete la interfaz web de Centinela, **ningún identificador debe proveer información que permita atacar directamente el sistema operativo o el hipervisor Proxmox VE**.

A continuación se analizan los identificadores presentes en el sistema, sus riesgos y su nivel de severidad:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        SUPERFICIE DE ATAQUE                            │
├───────────────────┬────────────────────────────┬───────────────────────┤
│ Identificador     │ Vulnerabilidad Principal   │ Clasificación OWASP   │
├───────────────────┼────────────────────────────┼───────────────────────┤
│ UPID de Proxmox   │ Fuga de Topología / PIDs   │ CWE-200 (Info Leak)   │
│ VMID de Instancia │ Enumeración / IDOR         │ CWE-639 (IDOR)        │
│ Nombre de Nodo    │ Mapeo de Red Interna       │ CWE-200 (Info Leak)   │
│ ID de Storage     │ Fuga de FileSystem / Pools │ CWE-200 (Info Leak)   │
│ Errores de Perl   │ Stack Trace / Path Leak    │ CWE-209 (Error Leak)  │
└───────────────────┴────────────────────────────┴───────────────────────┘
```

---

### 3.1. UPID de Proxmox VE
*   **Ubicación en Código:** [`backend/internal/adapters/secondary/proxmox/client.go`](file:///home/micaela/centinela/backend/internal/adapters/secondary/proxmox/client.go), tabla `tareas_asincronas.upid_proxmox`.
*   **Estructura del Dato:**  
    `UPID:<nodo>:<pid_hex>:<pstart_hex>:<starttime_hex>:<tipo_tarea>:<vmid>:<usuario>:<comentario>`  
    *Ejemplo real:* `UPID:proxmox:00083D71:0131F1BF:6A84EA54:qmstart:110:root@pam:`
*   **Vulnerabilidad:** **Fuga Crítica de Información de Infraestructura (CWE-200)**.
*   **Información Expuesta:**
    1.  `proxmox`: Hostname interno del servidor físico en la red.
    2.  `00083D71`: PID (Process ID) del proceso en el kernel de Linux. Permite deducir la carga del sistema y planificar ataques de agotamiento de PIDs o *PID recycling*.
    3.  `root@pam`: Revela la cuenta de administración privilegiada bajo la cual corre la tarea.
    4.  `6A84EA54`: Timestamp de arranque en milisegundos, permitiendo sincronización temporal en ataques de condición de carrera.
*   **Por qué es Importante:**
    Si el UPID llega al cliente, un usuario malicioso o una extensión de navegador infectada obtienen la radiografía de la arquitectura interna del datacenter sin necesidad de acceder a la consola. La decisión de mantenerlo **privado en el backend bajo un UUID opaco (`tareaId`)** es la solución más sólida porque no filtra nada de esta información.

---

### 3.2. VMID de Instancia (`vmid`: 100, 101, 102...)
*   **Ubicación en Código:** Rutas `/instances/:vmid`, respuestas JSON en `InstanciaInventarioResponse.id`.
*   **Estructura del Dato:** Números enteros secuenciales (`int`).
*   **Vulnerabilidad:** **Referencia Directa Insegura a Objetos (IDOR / CWE-639) y Enumeración Secuencial**.
*   **Riesgo e Impacto:**
    1.  **Enumeración de máquinas:** Los atacantes pueden sondear `/instances/100`, `/instances/101`, `/instances/102` para saber con precisión cuántas máquinas existen, cuáles fueron creadas primero y con qué velocidad se crean nuevas instancias.
    2.  **Explotación de fallas de autorización:** Si algún endpoint nuevo o modificado llegara a omitir el guard `RequireInstanceAccess`, el atacante sólo necesita iterar enteros del 100 al 200 para tomar control de instancias de otros clientes u organizaciones.
*   **Por qué es Importante:**
    Los VMIDs son los identificadores físicos directos que maneja QEMU y KVM en el kernel. **Deben ser cifrados u ofuscados** (mediante Hashids con clave secreta o identificadores públicos opacos `inst_9kX2mB`) para que el exterior nunca conozca el ID secuencial de hardware.

---

### 3.3. Nombres de Nodo del Hipervisor (`proxmox`, `pve-01`)
*   **Ubicación en Código:** `InstanciaInventarioResponse.node`, query parameters en rutas de métricas `/nodes/:node/rrddata`.
*   **Vulnerabilidad:** **Host Topography Disclosure (CWE-200)**.
*   **Información Expuesta:**
    Los nombres reales asignados a los hipervisores en el archivo `/etc/hosts` o en el DNS corporativo interno.
*   **Por qué es Importante:**
    En un ataque dirigido por etapas (kill-chain), la primera fase es el reconocimiento. Exponer que el host físico se llama `pve-cluster-node01.datacenter.interno` permite al atacante dirigir ataques de escaneo de puertos (SSH, IPMI, Proxmox Web GUI puerto 8006) directamente a la IP/DNS del host físico si comparte red o segmento de gestión. Se debe sustituir por un alias abstracto (`"nodo-principal"` o `"cluster-node-1"`).

---

### 3.4. Identificadores de Almacenamiento y Pools de Disco (`storageId`)
*   **Ubicación en Código:** Asistente de creación de instancias (`/instances`), DTOs de aprovisionamiento (`local-lvm`, `local`, `zfs-pool`).
*   **Vulnerabilidad:** **Storage Topology & Driver Disclosure**.
*   **Información Expuesta:**
    Revela la tecnología de almacenamiento utilizada: si es LVM-Thin (`local-lvm`), almacenamiento basado en archivos (`local`), ZFS o Ceph.
*   **Por qué es Importante:**
    Diferentes tecnologías de storage tienen vectores de ataque específicos (ej. fragmentación ZFS, vulnerabilidades en drivers de LVM, saturación de IOPS por falta de cuotas a nivel bloque). Exponer los nombres de pools permite diseñar cargas maliciosas dirigidas a colapsar el almacenamiento compartido del hipervisor.

---

### 3.5. Trazas de Error Crudas de la API de Proxmox
*   **Ubicación en Código:** [`backend/internal/adapters/primary/http/instance_handler.go`](file:///home/micaela/centinela/backend/internal/adapters/primary/http/instance_handler.go) (`mapearErrorProxmox`).
*   **Vulnerabilidad:** **Fuga de Rutas de Sistema de Archivos y Código Fuente (CWE-209)**.
*   **Información Expuesta:**
    Cuando Proxmox rechaza una operación, devuelve cadenas como:  
    `unable to create VM 115: storage 'local-lvm' does not exist at /usr/share/perl5/PVE/API2/Qemu.pm line 450`
*   **Por qué es Importante:**
    Revela:
    1.  Rutas absolutas del sistema operativo (`/usr/share/perl5/...`).
    2.  Versiones internas de módulos Perl de Proxmox PVE.
    3.  Líneas exactas de código en el servidor anfitrión donde falló la llamada.
    Toda respuesta de error debe ser sanitizada y canonizada antes de abandonar el backend.

---

## 4. Fundamentación del Plan de Remediación (`FIX-105` al `FIX-110`)

Para garantizar que el sistema sea hermético sin romper la experiencia de usuario ni la suite de pruebas, se definió la siguiente secuencia atómica en [`documentacion/futuro2.md`](file:///home/micaela/centinela/documentacion/futuro2.md):

| Fix | Componente | Acción de Seguridad | Objetivo contra Ataques |
| :--- | :---: | :--- | :--- |
| **`FIX-105`** | **Frontend** | Consumir `tareaId` y remover la validación de `upid`. | **Desbloqueo inmediato:** repara la interfaz gráfica y adopta el identificador opaco. |
| **`FIX-106`** | **Tests / Backend** | Validar `tareaId` y asegurar que `upid` nunca esté en el JSON. | **Garantía de regresión:** valida que ninguna tarea futura vuelva a exponer UPIDs por descuido. |
| **`FIX-107`** | **Backend** | Cifrado/Ofuscación reversible de VMIDs (Hashids con salt secreto). | **Protección contra IDOR:** ningún atacante puede deducir ni enumerar números de VM. |
| **`FIX-108`** | **Frontend** | Manejar `id` como string opaco en TypeScript y rutas. | **Compatibilidad cliente:** permite que la UI opere con identificadores ofuscados transparentemente. |
| **`FIX-109`** | **Backend** | Enmascaramiento de Nodos y Pools de Almacenamiento. | **Anti-Reconocimiento:** oculta la topología física de red y storage del datacenter. |
| **`FIX-110`** | **Backend** | Sanitización estricta de trazas de error de Proxmox y WebSockets. | **Cero fuga de sistema de archivos:** bloquea rutas de Perl, Linux y metadatos del kernel. |

---

## 5. Conclusión y Recomendación Final

1.  **Validación de la Opción A:** Ocultar el UPID detrás de un UUID opaco (`tareaId`) administrado en base de datos es la arquitectura correcta. No introduce la complejidad de descifrado simétrico en cada endpoint y garantiza un **aislamiento absoluto** de los metadatos de Proxmox hacia el exterior.
2.  **Extensión a Todos los IDs:** La ofuscación de los VMIDs enteros (`FIX-107`) y el enmascaramiento de nombres de host y storage (`FIX-109`) cerrarán los canales de reconocimiento que hoy quedan abiertos.
3.  **Prioridad 1:** Implementar de inmediato **`FIX-105`**, ya que el frontend se encuentra actualmente bloqueado por exigir el campo `upid` que el backend dejó de enviar.
