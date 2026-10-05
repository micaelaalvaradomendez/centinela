import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { commands, findDialog, findOutsideTable, renderInstances, requireAction, rowOf, stubFetch } from './instances-helpers';

// FRN-16: máquina de estados "Operación en progreso" por instancia (RF-04).
// Criterio de éxito: es imposible disparar una segunda acción sobre la misma instancia mientras hay una
// orden en curso. Cada uno de los seis códigos (D2) muestra su propio mensaje, y cualquier otro código,
// un mensaje genérico. Entregable: al confirmar el modal (FRN-15) se envía la orden por la ruta del
// contrato (BAC-29: /instances/:vmid/:accion o /instances/:vmid/status/:accion); el botón muestra un
// spinner y se deshabilitan todos los de la fila; ante un error la fila se desbloquea.

afterEach(() => {
  vi.unstubAllGlobals();
});

const json = (body: unknown, status: number) => new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });

async function confirmStop() {
  const user = userEvent.setup();
  const row = await rowOf('servidor-web');
  await user.click(await requireAction(row, 'stop'));
  const dialog = await findDialog();
  await user.click(within(dialog).getByRole('button', { name: /confirmar|detener|forzar|aceptar/i }));
  return { user, row };
}

const busy = (row: HTMLElement) => Boolean(
  row.getAttribute('aria-busy') === 'true'
  || within(row).queryByRole('status')
  || within(row).queryByRole('progressbar')
  || row.querySelector('.animate-spin, [data-loading="true"], [aria-busy="true"]'),
);

describe('FRN-16 - Máquina de estados "Operación en progreso" por instancia', () => {
  it('al confirmar envía la orden por la ruta del contrato y, mientras está en curso, bloquea toda la fila con un spinner', async () => {
    let release!: (response: Response) => void;
    const fetchMock = stubFetch(() => new Promise<Response>((resolve) => { release = resolve; }));
    renderInstances({ rol: 'ADMIN' });

    const { user, row } = await confirmStop();

    await waitFor(() => expect(commands(fetchMock)).toHaveLength(1));
    const [url, init] = commands(fetchMock)[0];
    expect(String(init?.method).toUpperCase()).toBe('POST');
    expect(url, 'la orden no va a la ruta de energía de la 101').toMatch(/\/instances\/101\/(status\/)?stop(\?|$)/);
    expect(new Headers(init?.headers).get('Authorization')).toBe('Bearer access-token-instancias');

    await waitFor(() => {
      const actions = within(row).queryAllByRole('button');
      expect(actions.length).toBeGreaterThan(0);
      for (const action of actions) expect(action, `"${action.getAttribute('aria-label') ?? action.textContent}" sigue habilitado con una orden en curso`).toBeDisabled();
    });
    expect(busy(row), 'la fila no muestra spinner ni estado de carga mientras la orden está en curso').toBe(true);

    // Una segunda acción sobre la misma instancia no envía nada.
    for (const action of within(row).queryAllByRole('button')) await user.click(action).catch(() => {});
    expect(commands(fetchMock), 'se pudo disparar una segunda orden sobre la misma instancia').toHaveLength(1);

    release(json({ upid: 'UPID:pve:1:stop:101:', tareaId: 't-1' }, 202));
  });

  const cases = [
    { status: 409, code: 'INSTANCE_INVALID_STATE', message: /estado actual|ya est[aá] (encendida|apagada|detenida|en ejecuci)|no corresponde/i },
    { status: 409, code: 'INSTANCE_BUSY', message: /otra tarea|ocupad/i },
    { status: 403, code: 'INSTANCE_PROTECTED', message: /protegid|infraestructura/i },
    { status: 403, code: 'INSTANCE_ACCESS_DENIED', message: /permiso/i },
    { status: 502, code: 'PROXMOX_UNAVAILABLE', message: /no est[aá] disponible|no disponible/i },
    { status: 504, code: 'PROXMOX_TIMEOUT', message: /no respondi[oó] a tiempo|tiempo de espera|timeout/i },
  ];
  const shown = new Map<string, string>();

  it.each(cases)('D2: $status $code muestra su mensaje y desbloquea la fila', async ({ status, code, message }) => {
    // El cuerpo trae un mensaje neutro: el texto que ve el usuario lo arma el frontend según el código.
    stubFetch(() => Promise.resolve(json({ errorCode: code, message: 'error' }, status)));
    renderInstances({ rol: 'ADMIN' });
    const { row } = await confirmStop();

    shown.set(code, await findOutsideTable(message));
    await waitFor(async () => {
      const stop = await requireAction(row, 'stop');
      expect(stop, 'después del error la fila tiene que desbloquearse').toBeEnabled();
    });
  });

  it('cualquier otro código muestra un mensaje genérico, distinto de los seis de D2', async () => {
    stubFetch(() => Promise.resolve(json({ errorCode: 'INTERNAL_ERROR', message: 'error' }, 500)));
    renderInstances({ rol: 'ADMIN' });
    await confirmStop();
    const generic = await findOutsideTable(/error|no se pudo|fall[oó]|inesperad/i);
    expect(cases.some(({ message }) => message.test(generic)), `el mensaje genérico coincide con el de un código de D2: "${generic}"`).toBe(false);
  });

  it('los seis mensajes de D2 son distintos entre sí', () => {
    expect(shown.size, 'algún caso de D2 no mostró su mensaje (ver los casos anteriores)').toBe(cases.length);
    expect(new Set(shown.values()).size, `hay mensajes repetidos: ${JSON.stringify(Object.fromEntries(shown))}`).toBe(cases.length);
  });
});

