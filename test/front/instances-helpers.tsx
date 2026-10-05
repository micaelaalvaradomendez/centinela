import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { expect, vi } from 'vitest';
import InstancesPage from '@/pages/Instances';
import { Toaster } from '@/components/ui/toast';
import { withProtectedProviders } from './app-providers';

// Helpers compartidos por las pruebas de la tabla de instancias (FRN-15, FRN-16).
// Las tareas no fijan si las acciones son botones en la fila o un menú "Acciones": se aceptan las dos
// formas. Sí se exige que cada acción tenga un nombre accesible (texto o aria-label).

export type Instance = {
  id: number; name: string; type: 'vm' | 'lxc'; node: string; status: string; ip: string | null;
  cpuUsage: number | null; ramUsage: number | null; maxRam: number | null; nivelAcceso: string; activeTask: null;
};

export const inventory: Instance[] = [
  { id: 101, name: 'servidor-web', type: 'vm', node: 'pve', status: 'running', ip: '192.168.1.50', cpuUsage: 25, ramUsage: 4294967296, maxRam: 17179869184, nivelAcceso: 'FULL_ACCESS', activeTask: null },
  { id: 102, name: 'base-datos', type: 'lxc', node: 'pve', status: 'stopped', ip: null, cpuUsage: null, ramUsage: null, maxRam: null, nivelAcceso: 'READ_ONLY', activeTask: null },
];

export const ACTIONS = {
  start: /^(iniciar|encender|start)\b/i,
  shutdown: /^(apagar|apagado ordenado|shutdown)\b(?!.*forz)/i,
  stop: /^(detener|forzar|stop)\b|apagado forzado/i,
  reboot: /^(reiniciar|reboot)\b/i,
  delete: /^(eliminar|borrar|delete)\b/i,
};
export type Action = keyof typeof ACTIONS;

export type Session = { rol: 'ADMIN' | 'OPERATOR'; permisos?: { vmid: number; nivelAcceso: 'FULL_ACCESS' | 'READ_ONLY' }[] };

export function renderInstances({ rol, permisos = [] }: Session) {
  window.sessionStorage.setItem('centinela_access', 'access-token-instancias');
  window.localStorage.setItem('centinela_user', JSON.stringify({
    id: 'user-test', organizacionId: 'org-1', nombreCompleto: 'Usuario de prueba', email: 'prueba@centinela.local',
    rol, instanciasPermitidas: permisos.map((p) => p.vmid), permisos, tiene2FA: true,
  }));
  return render(withProtectedProviders(
    <Toaster>
      <MemoryRouter initialEntries={['/instances']}>
        <InstancesPage />
      </MemoryRouter>
    </Toaster>,
  ));
}

// fetch simulado: GET /instances devuelve el inventario; el resto lo resuelve `onCommand`.
export function stubFetch(onCommand: (url: string, init: RequestInit) => Promise<Response>) {
  const fetchMock = vi.fn().mockImplementation((url: string, init: RequestInit = {}) => {
    const method = (init.method ?? 'GET').toUpperCase();
    if (method === 'GET' && /\/instances(\?|$)/.test(String(url))) {
      return Promise.resolve(new Response(JSON.stringify(inventory), { status: 200, headers: { 'Content-Type': 'application/json' } }));
    }
    if (method === 'GET') return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }));
    return onCommand(String(url), init);
  });
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

// Órdenes sobre instancias. Se excluye el resto (por ejemplo POST /events/ticket del cliente de eventos).
export const commands = (fetchMock: ReturnType<typeof vi.fn>) =>
  (fetchMock.mock.calls as [string, RequestInit?][]).filter(([url, init]) =>
    /\/instances\//.test(String(url)) && ['POST', 'DELETE', 'PUT', 'PATCH'].includes((init?.method ?? 'GET').toUpperCase()));

export async function rowOf(name: string) {
  const cell = await screen.findByText(name, {}, { timeout: 3000 });
  const row = cell.closest('tr, [role="row"]');
  expect(row, `"${name}" no está dentro de una fila de tabla`).not.toBeNull();
  return row as HTMLElement;
}

const named = (element: HTMLElement, pattern: RegExp) => {
  const name = (element.getAttribute('aria-label') ?? element.textContent ?? '').trim();
  return pattern.test(name);
};

// Acciones visibles de la fila: botones directos y, si hay un menú de acciones, sus ítems.
export async function actionsOf(row: HTMLElement): Promise<HTMLElement[]> {
  const user = userEvent.setup();
  const direct = within(row).queryAllByRole('button');
  const menuTrigger = direct.find((button) => button.getAttribute('aria-haspopup') || named(button, /^(acciones|m[aá]s|opciones)\b/i));
  if (!menuTrigger) return direct;
  await user.click(menuTrigger);
  const items = await waitFor(() => {
    const found = screen.queryAllByRole('menuitem');
    expect(found.length).toBeGreaterThan(0);
    return found;
  });
  return [...direct, ...items];
}

export async function findAction(row: HTMLElement, action: Action): Promise<HTMLElement | undefined> {
  return (await actionsOf(row)).find((element) => named(element, ACTIONS[action]));
}

export async function requireAction(row: HTMLElement, action: Action): Promise<HTMLElement> {
  const element = await findAction(row, action);
  expect(element, `la fila no ofrece la acción "${action}" (se busca un botón o ítem de menú con nombre ${ACTIONS[action]})`).toBeDefined();
  return element!;
}

// Busca texto fuera de la tabla (toasts, avisos), para no confundirlo con encabezados o celdas.
export async function findOutsideTable(pattern: RegExp, timeout = 3000): Promise<string> {
  return waitFor(() => {
    const matches = screen.queryAllByText(pattern).filter((node) => !node.closest('table, [role="table"], [role="grid"]'));
    expect(matches.length, `no apareció un mensaje que coincida con ${pattern}`).toBeGreaterThan(0);
    return matches[0].textContent ?? '';
  }, { timeout });
}

// El modal puede ser role="dialog" o role="alertdialog" (Radix usa el segundo para confirmaciones).
export const queryDialog = () => screen.queryByRole('dialog') ?? screen.queryByRole('alertdialog');
export async function findDialog(): Promise<HTMLElement> {
  return waitFor(() => {
    const dialog = queryDialog();
    expect(dialog, 'no se abrió ningún modal (role="dialog" o "alertdialog")').not.toBeNull();
    return dialog as HTMLElement;
  }, { timeout: 3000 });
}
