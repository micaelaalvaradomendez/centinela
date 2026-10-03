# Etapa 1: correcciones pendientes (FIX)

Correcciones de tareas de la Etapa 1 que ya están implementadas pero no cumplen del todo su criterio de éxito. Las tareas originales están en [`terminado-1.md`](terminado-1.md).

---

## 🛠️ Fixes detectados en la verificación del 03/10/2026

### `FIX-43` - Corregir maquetado, visualización de IP, badges de estado y columnas en tabla de instancias (`FRN-20A` / RF-03) (Frontend)

- **Área:** Frontend
- **Asignada:** Luz / Cristian (PR #74)
- **Estimación:** 1.0 h
- **Depende de:** `FRN-20A` (en `terminado-1.md`).
- **Problema y evidencia:**
  El PR #74 implementó la integración viva de `Instances.tsx` con `useInstances.ts` y `instanceService.ts`, consumiendo `GET /instances` y aplicando permisos. Sin embargo, la suite de pruebas de aceptación (`test/front/instances-table.test.tsx`) falla 5/5 por discrepancias de maquetado e interfaz:
  1. **Nombre y VMID:** `Instances.tsx:113` renderiza `<p>{instance.name} ({instance.id})</p>`. La prueba y el diseño requieren que el nombre y el VMID se presenten como elementos claramente identificables en la fila (o celdas diferenciadas), permitiendo consultar `within(row).getByText('101')` y `screen.findByText('servidor-web')` de forma unívoca.
  2. **Columna de IP:** `Instances.tsx:125` tiene hardcodeado un guión fijo (`<td className="px-4 py-4">—</td>`). `instanceService.ts` debe leer `instance.ip` de la respuesta, y `Instances.tsx` debe renderizar la IP (ej. `192.168.1.50`) o el texto `"No detectada"` cuando sea `null`, junto a un botón interactivo para copiar la dirección al portapapeles (`navigator.clipboard.writeText`).
  3. **Badges de estado:** `Instances.tsx:119` renderiza `{instance.status}` como texto simple sin estilos distintivos. El criterio de aceptación exige badges diferenciados con estilos visuales estándar: verde para `running` / `en ejecución` y gris para `stopped` / `detenida`.
  4. **Botonera de acciones:** cada fila debe contar con su botonera de acciones presente en la tabla, independientemente de si la fila tiene IP o si las acciones operativas están condicionadas.
- **Entregable:**
  1. En `instanceService.ts`, agregar el campo opcional `ip?: string | null` en `InventoryInstance` y leerlo en el mapeo de `fetchInstanceInventory`.
  2. En `Instances.tsx`, ajustar el renderizado de la columna de nombre para que el texto del nombre (`instance.name`) y el VMID (`instance.id`) estén en elementos o nodos de texto separados.
  3. Renderizar la columna IP mostrando `instance.ip` si existe (con botón de copia con icono y `aria-label="Copiar IP"`) o `"No detectada"` si es `null`.
  4. Envolver el estado en un badge con clases de Tailwind que apliquen fondo y texto verde para `running` (ej. `bg-green-100 text-green-700` o variante shadcn correspondiente) y gris para `stopped`.
  5. Asegurar que la columna tipo exponga claramente `VM` o `LXC`.
- **Criterio de éxito:**
  - Los 5 casos de `test/front/instances-table.test.tsx` pasan 100% en verde.
