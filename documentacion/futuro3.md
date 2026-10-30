# Etapa 3: correcciones pendientes (FIX)

Acá se listan los fixes y mejoras técnicas correspondientes a las tareas de la Etapa 3 identificadas durante la verificación y ejecución de suites de aceptación.

---

## 🛠️ Especificación Detallada de Fixes (Etapa 3)

### `FIX-128` - Opciones completas de acciones de auditoría Etapa 3 y selectores en Auditoria.tsx (`FRN-31`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-31`, `BAC-44`.
- **Problema técnico detallado:**
  En `frontend/centinela/src/utils/auditActionLabels.ts`, el listado de acciones en `AUDIT_ACTION_LABELS` y `AUDIT_ACTION_ORDER` omite las acciones `SNAPSHOT_DELETE` ("Snapshot eliminado") y `ALERT_TRIGGERED` ("Alerta disparada"). Asimismo, `SNAPSHOT_ROLLBACK` está rotulado como `"Snapshot restaurado"`, lo que no coincide con las expresiones esperadas de reversión/rollback (`/revertir.*snapshot|rollback|snapshot_rollback/i`). Por último, `Auditoria.tsx` renderiza un único botón estático de `"Exportar CSV"` en lugar de ofrecer la opción de descarga tanto en CSV como en JSON integrada con `BAC-45`.
- **Qué debe hacer el equipo de frontend:**
  1. En `frontend/centinela/src/utils/auditActionLabels.ts`:
     - Incorporar en `AUDIT_ACTION_LABELS` y `AUDIT_ACTION_ORDER`:
       - `SNAPSHOT_DELETE: 'Snapshot eliminado'`
       - `ALERT_TRIGGERED: 'Alerta disparada'`
     - Asegurar que `SNAPSHOT_ROLLBACK` sea etiquetado como `"Revertir snapshot"` o `"Snapshot revertido"` para satisfacer la búsqueda por vocabulario canónico.
  2. En `frontend/centinela/src/pages/Auditoria.tsx`:
     - Incorporar en el filtro de tipo de recurso la opción `SNAPSHOT`.
     - Añadir botón o menú de exportación que permita seleccionar entre `"Exportar como CSV"` y `"Exportar como JSON"`, disparando la descarga de `auditoria_centinela_YYYYMMDD.json` con `format=json`.
- **Criterio de éxito:**
  La suite `test/front/frn-31-audit-etapa3.test.tsx` pasa 4 de 4 pruebas en verde (100% ✅).

---

### `FIX-129` - Endurecimiento de cabeceras de seguridad en Nginx de producción (`INF-11`) (Infraestructura)

- **Área:** Infraestructura
- **Asignado:** Nico
- **Estimación:** 1.0 h
- **Depende de:** `INF-11`.
- **Problema técnico detallado:**
  En el archivo de configuración `docker/nginx-edge.conf:36`:
  1. La cabecera `Strict-Transport-Security` está configurada como `"max-age=31536000"` pero omite la directiva obligatoria `includeSubDomains`.
  2. Faltan las directivas de seguridad para mitigar clickjacking y sniffing de MIME:
     - `add_header X-Frame-Options "DENY" always;`
     - `add_header X-Content-Type-Options "nosniff" always;`
  3. Falta la directiva `Content-Security-Policy` ajustada para scripts, fuentes locales y canales bidireccionales WebSocket (`connect-src 'self' wss: ws:;`).
- **Qué debe hacer el equipo de infraestructura:**
  En `docker/nginx-edge.conf`, agregar y ajustar las cabeceras requeridas en el bloque `server` HTTPS (puerto 443):
  ```nginx
  add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
  add_header X-Frame-Options "DENY" always;
  add_header X-Content-Type-Options "nosniff" always;
  add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; connect-src 'self' ws: wss:; img-src 'self' data:;" always;
  ```
- **Criterio de éxito:**
  La prueba `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS` en `test/back/etapa3_acceptance_test.go` pasa todos sus subtests en verde.

---

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

### `FIX-131` - Incorporación de métricas y gráfico de Tráfico de Red (`netin`/`netout`) en detalle de instancia (`CU-05`) (Frontend)

- **Área:** Frontend
- **Asignado:** Belinda / Cristian
- **Estimación:** 1.5 h
- **Depende de:** `FRN-22`, `FRN-23`.
- **Problema técnico detallado:**
  En el flujo de **CU-05** (`casosdeuso.md:275`), se especifica que:
  > *"El frontend renderiza gráficos de líneas (CPU, RAM y Tráfico de Red) mostrando la evolución temporal."*
  
  El backend en `domain.PuntoMetricaDTO` ya expone `netInBps` y `netOutBps` provenientes de Proxmox RRDdata. Sin embargo, en el frontend (`frontend/centinela/src/components/features/instances/types/instanceDetail.ts`), la interfaz `InstanceMetricPoint` omite las propiedades de red, y `InstanceResourceUsage.tsx` sólo instancia columnas y curvas para CPU, Memoria y Disco, dejando sin cobertura la gráfica de tráfico de red exigida por el caso de uso.
- **Qué debe hacer el equipo de frontend:**
  1. En `instanceDetail.ts`:
     - Agregar a `InstanceMetricPoint`: `netInBps: number | null; netOutBps: number | null;`.
  2. En `instanceDetailService.ts`:
     - Normalizar los campos `netInBps` y `netOutBps` en el mapeador de puntos métricos.
  3. En `InstanceResourceUsage.tsx`:
     - Incorporar la métrica y columna de **Tráfico de Red** utilizando el icono `Network` o `ArrowUpDown`, formateando las tasas con una función auxiliar `formatNetworkRate(bytesPerSec)` (`KB/s`, `MB/s`) e instanciando un `InstanceResourceChart` para la serie de tráfico de red combinado (`(netInBps ?? 0) + (netOutBps ?? 0)` o líneas segregadas de entrada/salida).
- **Criterio de éxito:**
  Al ingresar a la vista de detalle de una instancia con métricas históricas (2h o 24h), la grilla de uso de recursos renderiza las 4 métricas (CPU, Memoria, Disco y Tráfico de Red) con sus respectivos gráficos de evolución temporal.

---

### `FIX-132` - Discriminación de errores 400 en recuperación de contraseña para no quemar intentos de OTP ante contraseñas inválidas o duplicadas (`CU-12`) (Frontend)

- **Área:** Frontend
- **Asignado:** Luz / Belinda
- **Estimación:** 1.0 h
- **Depende de:** `FRN-34`.
- **Problema técnico detallado:**
  En el flujo de excepciones de **CU-12** (`casosdeuso.md:612`), se establece taxativamente:
  > *"Contraseña débil o igual a la anterior (`400 WEAK_PASSWORD` / `400 SAME_PASSWORD`): El backend rechaza el cambio sin quemar el código OTP para que el usuario pueda corregir la contraseña elegida."*
  
  Sin embargo, en `frontend/centinela/src/pages/RecoverPassword.tsx:130-135`:
  ```tsx
  if (error.errorCode === 'RESET_FAILED' || error.status === 400) {
      const nextAttempts = failedAttempts + 1;
      ...
      setCodeError('El código ingresado es incorrecto o ha expirado.');
  ```
  La condición `error.status === 400` captura de forma indiscriminada cualquier error 400 devuelto por `POST /api/auth/password/reset` (incluyendo `400 SAME_PASSWORD`, `400 WEAK_PASSWORD` o errores de validación de contraseña). Esto provoca que la interfaz acuse erróneamente al código OTP de ser incorrecto, incremente `failedAttempts` y bloquee el proceso (`setResetBlocked(true)`) al 3.er intento, quemando en el cliente un código OTP que es completamente válido en el servidor.
- **Qué debe hacer el equipo de frontend:**
  En `RecoverPassword.tsx`:
  1. Capturar específicamente los códigos de error de contraseña previa al conteo de intentos:
     ```tsx
     if (error.errorCode === 'SAME_PASSWORD' || error.errorCode === 'WEAK_PASSWORD' || error.errorCode === 'PASSWORD_RECENTLY_USED') {
         setPasswordError(error.message || 'La nueva contraseña no puede ser igual a la anterior ni violar las políticas de seguridad.');
         return;
     }
     ```
  2. Incrementar `failedAttempts` y aplicar el mensaje `'El código ingresado es incorrecto o ha expirado.'` **únicamente** cuando `error.errorCode === 'RESET_FAILED'`.
- **Criterio de éxito:**
  Si el usuario ingresa un código OTP válido pero una contraseña igual a la actual o no conforme a la política del backend, la interfaz muestra el error sobre el campo de contraseña, no incrementa los intentos fallidos del OTP y permite reintentar con otra contraseña sin bloquear el proceso.
