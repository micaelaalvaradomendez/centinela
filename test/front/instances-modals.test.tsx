import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { actionsOf, ACTIONS, commands, findAction, findDialog, queryDialog, renderInstances, requireAction, rowOf, stubFetch } from './instances-helpers';

// FRN-15: modales de confirmación antierror para acciones operativas (RF-04).
// Criterio de éxito:
//   - Ninguna acción se dispara sin pasar por el modal, y cancelar no envía peticiones.
//   - Un OPERATOR con READ_ONLY no ve botones de energía; uno con FULL_ACCESS sí.
//   - El OPERATOR nunca ve Delete.
// Entregable: Start (confirmación estándar), Shutdown/Reboot (aviso de apagado o reinicio del SO
// huésped), Stop (advertencia en rojo por posible pérdida de datos), Delete (tipear ID o nombre, solo ADMIN).
// La forma del modal es libre: se busca role="dialog" o role="alertdialog".

let fetchMock: ReturnType<typeof vi.fn>;

beforeEach(() => {
  fetchMock = stubFetch(() => Promise.resolve(new Response(JSON.stringify({ upid: 'UPID:pve:1:start:101:', tareaId: 't-1' }), { status: 202, headers: { 'Content-Type': 'application/json' } })));
});

afterEach(() => {
  vi.unstubAllGlobals();
});


async function openAndCancel(name: string, action: keyof typeof ACTIONS) {
  const user = userEvent.setup();
  const row = await rowOf(name);
  await user.click(await requireAction(row, action));
  const dialog = await findDialog();
  expect(commands(fetchMock), `"${action}" envió la orden antes de confirmar el modal`).toHaveLength(0);
  return { user, dialog };
}

async function cancel(user: ReturnType<typeof userEvent.setup>, dialog: HTMLElement) {
  await user.click(within(dialog).getByRole('button', { name: /cancelar|volver|cerrar/i }));
  await waitFor(() => expect(queryDialog()).toBeNull());
  expect(commands(fetchMock), 'cancelar el modal envió una petición').toHaveLength(0);
}

describe('FRN-15 - Modales de confirmación antierror para acciones operativas', () => {
  it('Start pide confirmación y cancelar no envía ninguna petición', async () => {
    renderInstances({ rol: 'ADMIN' });
    const { user, dialog } = await openAndCancel('base-datos', 'start');
    expect(within(dialog).getByRole('button', { name: /confirmar|iniciar|encender|aceptar/i })).toBeInTheDocument();
    await cancel(user, dialog);
  });

  it.each(['shutdown', 'reboot'] as const)('%s avisa que se apaga o reinicia el sistema operativo huésped, y cancelar no envía nada', async (action) => {
    renderInstances({ rol: 'ADMIN' });
    const { user, dialog } = await openAndCancel('servidor-web', action);
    expect(dialog.textContent).toMatch(/sistema operativo|hu[eé]sped|guest/i);
    await cancel(user, dialog);
  });

  it('Stop advierte en rojo sobre la posible pérdida de datos, y cancelar no envía nada', async () => {
    renderInstances({ rol: 'ADMIN' });
    const { user, dialog } = await openAndCancel('servidor-web', 'stop');
    expect(dialog.textContent).toMatch(/p[eé]rdida de datos|perder (los )?datos|datos no guardados/i);
    const red = [dialog, ...Array.from(dialog.querySelectorAll<HTMLElement>('*'))].some((element) =>
      /(^|\s|:)(text|bg|border)-(red|rose)-\d|destructive|danger/.test(element.getAttribute('class') ?? '')
      || /color:\s*(red|#(dc2626|ef4444|b91c1c|f87171))|rgb\((220, 38, 38|239, 68, 68)\)/i.test(element.getAttribute('style') ?? ''));
    expect(red, 'la advertencia de Stop no tiene un estilo rojo (clases red/rose/destructive o color rojo)').toBe(true);
    await cancel(user, dialog);
  });

  it('Delete pide tipear el ID o el nombre: con otro texto no se habilita, y cancelar no envía nada', async () => {
    renderInstances({ rol: 'ADMIN' });
    const { user, dialog } = await openAndCancel('base-datos', 'delete');
    const input = within(dialog).getByRole('textbox');
    const confirm = within(dialog).getByRole('button', { name: /eliminar|borrar|confirmar/i });
    expect(confirm, 'el botón de confirmar Delete tiene que empezar deshabilitado').toBeDisabled();

    await user.type(input, 'otra-cosa');
    expect(confirm, 'un texto que no es el ID ni el nombre habilitó el borrado').toBeDisabled();

    // Se acepta el nombre o el ID (la tarea dice "el ID o el nombre").
    await user.clear(input);
    await user.type(input, 'base-datos');
    if ((confirm as HTMLButtonElement).disabled) {
      await user.clear(input);
      await user.type(input, '102');
    }
    expect(confirm, 'ni el nombre ni el ID de la instancia habilitan el borrado').toBeEnabled();
    expect(commands(fetchMock), 'tipear la confirmación envió la orden sin pulsar el botón').toHaveLength(0);
    await cancel(user, dialog);
  });

  it('un OPERATOR ve acciones de energía solo donde tiene FULL_ACCESS y nunca ve Delete', async () => {
    renderInstances({ rol: 'OPERATOR', permisos: [{ vmid: 101, nivelAcceso: 'FULL_ACCESS' }, { vmid: 102, nivelAcceso: 'READ_ONLY' }] });
    const energy = [ACTIONS.start, ACTIONS.shutdown, ACTIONS.stop, ACTIONS.reboot];
    const isEnergy = (element: HTMLElement) => energy.some((pattern) => pattern.test((element.getAttribute('aria-label') ?? element.textContent ?? '').trim()));

    const full = await actionsOf(await rowOf('servidor-web'));
    expect(full.filter(isEnergy).length, 'con FULL_ACCESS sobre la 101 el OPERATOR tiene que ver acciones de energía').toBeGreaterThan(0);
    expect(full.filter(isEnergy).every((element) => !(element as HTMLButtonElement).disabled), 'con FULL_ACCESS las acciones de energía no pueden estar deshabilitadas').toBe(true);

    const readOnly = await actionsOf(await rowOf('base-datos'));
    expect(readOnly.filter(isEnergy).map((element) => element.getAttribute('aria-label') ?? element.textContent), 'con READ_ONLY sobre la 102 el OPERATOR no tiene que ver acciones de energía').toHaveLength(0);

    for (const name of ['servidor-web', 'base-datos']) {
      expect(await findAction(await rowOf(name), 'delete'), `el OPERATOR ve Delete en ${name}`).toBeUndefined();
    }
  });
});
