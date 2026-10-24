import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import InstanceDetailPage from '@/pages/InstanceDetailPage';
import { Toaster } from '@/components/ui/toast';
import { withProtectedProviders } from './app-providers';

// FRN-27 y FRN-28: Modal y panel de edición de recursos (CPU y RAM) con validación de cuotas
// e integración con PUT /api/instances/:instanceRef/config y aviso de reinicio pendiente.
// Entregables (documentacion/actual.md):
//   FRN-27:
//     1. Modal o sección "Editar Recursos" en InstanceDetailPage, visible/habilitado solo para
//        ADMIN u OPERATOR con FULL_ACCESS (oculto o deshabilitado para READ_ONLY).
//     2. Campos numéricos precargados con los valores actuales de la instancia (cores y ramMb).
//     3. Validación en cliente (cores >= 1, ramMb entero positivo, límites del host).
//     4. Cálculo reactivo de incremento (DeltaRAM > 0) contra allocatable.freeRamMb del host,
//        bloqueando el guardado cuando excede la memoria libre segura.
//     5. Nota informativa amarilla cuando la instancia está encendida ("Si la instancia está
//        encendida y no cuenta con módulos hot-plug activos...").
//   FRN-28:
//     1. Servicio updateInstanceConfiguration invocando PUT /api/instances/:instanceRef/config.
//     2. requiresReboot: false -> toast verde de éxito y actualización inmediata de CPU/RAM.
//     3. requiresReboot: true -> toast amarillo de advertencia y aviso de reinicio pendiente.
//     4. Manejo de errores 409 RESOURCE_QUOTA_EXCEEDED y 409 INSTANCE_BUSY mostrando el mensaje.

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

const GIB = 1024 ** 3;
const MIB = 1024 ** 2;

class SilentStream {
  readyState = 1;
  onopen: unknown = null;
  onmessage: unknown = null;
  onerror: unknown = null;
  onclose: unknown = null;
  constructor(public url: string) {}
  addEventListener() {}
  removeEventListener() {}
  send() {}
  close() { this.readyState = 3; }
}

class FakeResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

let fetchMock: ReturnType<typeof vi.fn>;

interface BackendStubOptions {
  instanceRef?: string;
  status?: 'running' | 'stopped';
  nivelAcceso?: 'FULL_ACCESS' | 'READ_ONLY';
  cores?: number;
  ramMb?: number;
  maxCpuCores?: number;
  freeRamMb?: number;
  putConfigResponse?: () => Response | Promise<Response>;
}

function stubBackend(options: BackendStubOptions = {}) {
  const ref = options.instanceRef ?? 'inst_k8Pq2v';
  const status = options.status ?? 'running';
  const nivelAcceso = options.nivelAcceso ?? 'FULL_ACCESS';
  let currentCores = options.cores ?? 2;
  let currentRamMb = options.ramMb ?? 2048;
  const maxCpuCores = options.maxCpuCores ?? 8;
  const freeRamMb = options.freeRamMb ?? 4096;

  fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
    const target = String(url);
    const method = (init.method ?? 'GET').toUpperCase();

    if ((method === 'PUT' || method === 'PATCH') && new RegExp(`/instances/${ref}/config$`).test(target)) {
      if (options.putConfigResponse) {
        return Promise.resolve(options.putConfigResponse());
      }
      const payload = JSON.parse(String(init.body ?? '{}')) as { cores?: number; ramMb?: number };
      if (typeof payload.cores === 'number') currentCores = payload.cores;
      if (typeof payload.ramMb === 'number') currentRamMb = payload.ramMb;
      return Promise.resolve(json({
        requiresReboot: false,
        message: 'La configuración de recursos se actualizó correctamente.',
      }, 200));
    }

    if (new RegExp(`/instances/${ref}/metrics`).test(target)) {
      return Promise.resolve(json({ data: [] }));
    }

    if (new RegExp(`/instances/${ref}(\\?|$)`).test(target)) {
      return Promise.resolve(json({
        vmid: ref,
        nombre: 'db-principal',
        nodo: 'nodo-a1',
        estado: status,
        tipo: 'qemu',
        cpu: status === 'running' ? 18 : 0,
        mem: status === 'running' ? 1024 * MIB : 0,
        maxCpu: currentCores,
        maxMem: currentRamMb * MIB,
      }));
    }

    if (/\/instances(\?|$)/.test(target)) {
      return Promise.resolve(json([
        {
          id: ref,
          name: 'db-principal',
          type: 'vm',
          node: 'nodo-a1',
          status,
          ip: '10.0.0.25',
          cpuUsage: status === 'running' ? 18 : null,
          ramUsage: status === 'running' ? 1024 * MIB : null,
          maxRam: currentRamMb * MIB,
          nivelAcceso,
          activeTask: null,
        },
      ]));
    }

    if (/\/node\/status(\?|$)/.test(target)) {
      return Promise.resolve(json({
        cpu: { usagePercent: 15, cores: maxCpuCores },
        ram: { usedGb: 4, totalGb: 16, usagePercent: 25 },
        storage: { usedGb: 40, totalGb: 200, usagePercent: 20 },
        uptimeSeconds: 7200,
        allocatable: {
          maxCpuCores,
          freeRamMb,
          freeStorageGb: 120,
          headroomPercent: 15,
        },
        instancesSummary: {
          vms: { running: 1, stopped: 0, paused: 0, total: 1 },
          lxc: { running: 0, stopped: 0, paused: 0, total: 0 },
        },
        stale: false,
        fetchedAt: '2026-10-24T12:00:00Z',
      }));
    }

    if (target.includes('/events/ticket') || target.includes('/auth/ws-ticket')) {
      return Promise.resolve(json({ ticket: 'ticket-frn27' }));
    }

    return Promise.resolve(json({}));
  });

  vi.stubGlobal('fetch', fetchMock);
}

function renderDetailAs(
  rol: 'ADMIN' | 'OPERATOR',
  nivelAcceso: 'FULL_ACCESS' | 'READ_ONLY' = 'FULL_ACCESS',
  instanceRef = 'inst_k8Pq2v',
) {
  window.sessionStorage.setItem('centinela_access', 'access-token-frn27');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: `${rol.toLowerCase()}-1`,
    organizacionId: 'org-1',
    nombreCompleto: `${rol} Test`,
    email: `${rol.toLowerCase()}@centinela.local`,
    rol,
    instanciasPermitidas: [instanceRef],
    permisos: [{ vmid: instanceRef, nivelAcceso }],
    tiene2FA: true,
  }));

  return render(withProtectedProviders(
    <MemoryRouter initialEntries={[`/instances/${instanceRef}`]}>
      <Toaster />
      <Routes>
        <Route path="/instances/:vmid" element={<InstanceDetailPage />} />
      </Routes>
    </MemoryRouter>,
  ));
}

async function openResourcesTab(user: ReturnType<typeof userEvent.setup>) {
  const tab = await screen.findByRole('tab', { name: /recursos/i });
  await user.click(tab);
}

beforeEach(() => {
  HTMLDialogElement.prototype.showModal ??= function showModalPolyfill(this: HTMLDialogElement) {
    this.setAttribute('open', '');
  };
  HTMLDialogElement.prototype.close ??= function closePolyfill(this: HTMLDialogElement) {
    this.removeAttribute('open');
  };
  vi.stubGlobal('EventSource', SilentStream);
  vi.stubGlobal('WebSocket', SilentStream);
  vi.stubGlobal('ResizeObserver', FakeResizeObserver);
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('FRN-27 - Formulario y modal de edición de recursos (CPU y RAM) con validación de cuotas', () => {
  it('precarga cores y ramMb actuales, muestra la nota amarilla si está encendida y calcula DeltaRAM reactivamente', async () => {
    stubBackend({ cores: 2, ramMb: 2048, freeRamMb: 4096, maxCpuCores: 8, status: 'running' });
    renderDetailAs('ADMIN');
    const user = userEvent.setup();

    await openResourcesTab(user);

    const cpuInput = await screen.findByLabelText(/n[uú]cleos \(vcpu\)/i) as HTMLInputElement;
    const ramInput = await screen.findByLabelText(/memoria ram \(mb\)/i) as HTMLInputElement;

    expect(cpuInput.value).toBe('2');
    expect(ramInput.value).toBe('2048');

    // Nota amarilla de advertencia cuando la instancia está encendida
    expect(screen.getByText(/si la instancia est[aá] encendida y no cuenta con m[oó]dulos hot-plug activos/i)).toBeVisible();

    // Aumentar RAM a 3072 MB (+1024 MB de DeltaRAM)
    fireEvent.change(ramInput, { target: { value: '3072' } });
    expect(await screen.findByText(/incremento solicitado:\s*\+1024\s*mb/i)).toBeVisible();
    expect(screen.getByText(/quedar[aá]n\s*3072\s*mb libres en el host/i)).toBeVisible();
  });

  it('valida límites en cliente (cores >= 1, máximo del host y exceso de cuota de RAM libre) y bloquea el guardado', async () => {
    stubBackend({ cores: 2, ramMb: 2048, freeRamMb: 1024, maxCpuCores: 4, status: 'running' });
    renderDetailAs('ADMIN');
    const user = userEvent.setup();

    await openResourcesTab(user);

    // Abrir el modal de edición de recursos
    await user.click(await screen.findByRole('button', { name: /editar recursos/i }));
    const dialog = await screen.findByRole('dialog');

    const cpuInput = within(dialog).getByLabelText(/n[uú]cleos \(vcpu\)/i);
    const ramInput = within(dialog).getByLabelText(/memoria ram \(mb\)/i);
    const saveButton = within(dialog).getByRole('button', { name: /guardar cambios/i });

    // 1. CPU por encima del máximo del host (maxCpuCores = 4)
    fireEvent.change(cpuInput, { target: { value: '8' } });
    expect(await within(dialog).findByText(/el m[aá]ximo del host es 4 n[uú]cleos/i)).toBeVisible();
    expect(saveButton).toBeDisabled();

    // Restaurar CPU a un valor válido (4) y exceder la memoria libre del host (+2048 MB > 1024 MB libres)
    fireEvent.change(cpuInput, { target: { value: '4' } });
    fireEvent.change(ramInput, { target: { value: '4096' } });

    expect(await within(dialog).findByText(/el incremento solicitado \(\+2048 mb\) excede la memoria disponible en el host \(1024 mb libres\)/i)).toBeVisible();
    expect(saveButton).toBeDisabled();
  });

  it('oculta o deshabilita la edición de recursos para un OPERATOR con permiso READ_ONLY', async () => {
    stubBackend({ nivelAcceso: 'READ_ONLY', cores: 2, ramMb: 2048 });
    renderDetailAs('OPERATOR', 'READ_ONLY');
    const user = userEvent.setup();

    await openResourcesTab(user);

    const cpuInput = await screen.findByLabelText(/n[uú]cleos \(vcpu\)/i);
    const ramInput = await screen.findByLabelText(/memoria ram \(mb\)/i);

    expect(cpuInput).toBeDisabled();
    expect(ramInput).toBeDisabled();
    expect(screen.queryByRole('button', { name: /editar recursos/i })).toBeNull();
    expect(screen.getByRole('button', { name: /guardar cambios/i })).toBeDisabled();
  });
});

describe('FRN-28 - Integración con PUT /api/instances/:instanceRef/config y aviso de reinicio pendiente', () => {
  it('envía PUT /api/instances/:instanceRef/config y con requiresReboot: false muestra toast de éxito y actualiza la vista', async () => {
    stubBackend({ cores: 2, ramMb: 2048, freeRamMb: 4096, maxCpuCores: 8 });
    renderDetailAs('OPERATOR', 'FULL_ACCESS');
    const user = userEvent.setup();

    await openResourcesTab(user);
    await user.click(await screen.findByRole('button', { name: /editar recursos/i }));

    const dialog = await screen.findByRole('dialog');
    const cpuInput = within(dialog).getByLabelText(/n[uú]cleos \(vcpu\)/i);
    const ramInput = within(dialog).getByLabelText(/memoria ram \(mb\)/i);

    fireEvent.change(cpuInput, { target: { value: '4' } });
    fireEvent.change(ramInput, { target: { value: '3072' } });

    const saveButton = within(dialog).getByRole('button', { name: /guardar cambios/i });
    await waitFor(() => expect(saveButton).toBeEnabled());
    await user.click(saveButton);

    // Paso de confirmación dentro del modal
    const confirmButton = await within(dialog).findByRole('button', { name: /confirmar cambios/i });
    await user.click(confirmButton);

    // Verifica llamada PUT/PATCH /instances/inst_k8Pq2v/config con { cores: 4, ramMb: 3072 }
    await waitFor(() => {
      const configCall = (fetchMock.mock.calls as [string, RequestInit][]).find(
        ([url, init]) => ['PUT', 'PATCH'].includes((init?.method ?? 'GET').toUpperCase())
          && /\/instances\/inst_k8Pq2v\/config$/.test(String(url)),
      );
      expect(configCall).toBeDefined();
      expect(JSON.parse(String(configCall![1].body))).toEqual({ cores: 4, ramMb: 3072 });
    });

    // Toast de éxito, sin aviso de reinicio pendiente, y valores actualizados tanto en inputs como en la tarjeta Hardware del Resumen
    expect(await screen.findByText(/cambios guardados/i)).toBeVisible();
    expect(await screen.findByText(/la configuraci[oó]n de recursos se actualiz[oó] correctamente/i)).toBeVisible();
    expect(screen.queryByText(/los recursos guardados se aplicar[aá]n al reiniciar la instancia/i)).toBeNull();
    expect((await screen.findByLabelText(/n[uú]cleos \(vcpu\)/i) as HTMLInputElement).value).toBe('4');
    expect((await screen.findByLabelText(/memoria ram \(mb\)/i) as HTMLInputElement).value).toBe('3072');

    // Al volver a la pestaña Resumen, la tarjeta Hardware refleja reactivamente 4 vCPU y 3 GB de RAM
    await user.click(screen.getByRole('tab', { name: /resumen/i }));
    const hardwareHeading = await screen.findByRole('heading', { name: /^hardware$/i });
    const hardwareCard = hardwareHeading.closest('[data-slot="card"]') ?? hardwareHeading.parentElement!;
    expect(within(hardwareCard as HTMLElement).getByText('4')).toBeVisible();
    expect(within(hardwareCard as HTMLElement).getByText(/3(\.0+)?\s*gb/i)).toBeVisible();
  });

  it('cuando el backend responde requiresReboot: true muestra advertencia amarilla y aviso de reinicio pendiente', async () => {
    const rebootMsg = 'Configuración registrada en Proxmox. Los cambios de CPU/RAM requieren reiniciar la instancia para aplicarse por completo.';
    stubBackend({
      cores: 2,
      ramMb: 2048,
      putConfigResponse: () => json({
        requiresReboot: true,
        message: rebootMsg,
      }, 200),
    });

    renderDetailAs('ADMIN');
    const user = userEvent.setup();

    await openResourcesTab(user);
    await user.click(await screen.findByRole('button', { name: /editar recursos/i }));

    const dialog = await screen.findByRole('dialog');
    fireEvent.change(within(dialog).getByLabelText(/n[uú]cleos \(vcpu\)/i), { target: { value: '4' } });

    await user.click(within(dialog).getByRole('button', { name: /guardar cambios/i }));
    await user.click(await within(dialog).findByRole('button', { name: /confirmar cambios/i }));

    // Muestra el toast de advertencia con el mensaje del backend y el banner de reinicio pendiente en la página
    expect(await screen.findByText(rebootMsg)).toBeVisible();
    const rebootNotice = await screen.findByRole('status');
    expect(rebootNotice).toHaveTextContent(/los recursos guardados se aplicar[aá]n al reiniciar la instancia/i);
  });

  it('maneja errores 409 RESOURCE_QUOTA_EXCEEDED y 409 INSTANCE_BUSY mostrando el mensaje descriptivo en la interfaz', async () => {
    let errorMode: 'QUOTA' | 'BUSY' = 'QUOTA';
    const quotaErrorMsg = 'El incremento de memoria solicitado excede la capacidad disponible del nodo.';
    const busyErrorMsg = 'La instancia tiene una operación asíncrona en curso. Intente nuevamente cuando finalice.';

    stubBackend({
      cores: 2,
      ramMb: 2048,
      putConfigResponse: () => (errorMode === 'QUOTA'
        ? json({ errorCode: 'RESOURCE_QUOTA_EXCEEDED', message: quotaErrorMsg }, 409)
        : json({ errorCode: 'INSTANCE_BUSY', message: busyErrorMsg }, 409)),
    });

    renderDetailAs('ADMIN');
    const user = userEvent.setup();

    await openResourcesTab(user);
    await user.click(await screen.findByRole('button', { name: /editar recursos/i }));

    const dialog = await screen.findByRole('dialog');
    fireEvent.change(within(dialog).getByLabelText(/memoria ram \(mb\)/i), { target: { value: '3072' } });

    await user.click(within(dialog).getByRole('button', { name: /guardar cambios/i }));
    await user.click(await within(dialog).findByRole('button', { name: /confirmar cambios/i }));

    expect((await screen.findAllByText(quotaErrorMsg)).length).toBeGreaterThan(0);

    // Ahora probamos 409 INSTANCE_BUSY (errorMapper traduce INSTANCE_BUSY a su mensaje canónico)
    errorMode = 'BUSY';
    await user.click(within(dialog).getByRole('button', { name: /confirmar cambios/i }));

    expect((await screen.findAllByText(/la instancia se encuentra ejecutando otra tarea|operaci[oó]n as[ií]ncrona en curso/i)).length).toBeGreaterThan(0);
  });
});
