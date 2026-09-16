import { describe, it } from 'vitest';

describe('FRN-05 - panel de gestión de usuarios', () => {
  it.todo('expone la ruta /admin/users únicamente para una sesión ADMIN');
  it.todo('rechaza la navegación de un usuario OPERATOR con estado 403 o vista denegada');
  it.todo('consulta GET /api/users con Authorization Bearer y renderiza summary y users');
  it.todo('permite filtrar usuarios por rol, activo y búsqueda');
});

describe('FRN-06 - alta y desactivación de usuarios', () => {
  it.todo('abre el modal de alta con nombre, username, email y rol');
  it.todo('envía POST /api/users y muestra una sola vez la contraseña temporal recibida');
  it.todo('actualiza la tabla después de crear un usuario');
  it.todo('confirma la desactivación y envía DELETE /api/users/:id');
  it.todo('muestra errores 400, 403 y 409 de forma comprensible');
});

describe('FRN-06B - edición de usuario y cambio de rol', () => {
  it.todo('precarga los datos del usuario seleccionado en el modal de edición');
  it.todo('envía PUT /api/users/:id con nombre, email, rol y activo');
  it.todo('actualiza la fila y muestra una confirmación después de editar');
});

describe('FRN-07 - selector de asignación de instancias', () => {
  it.todo('consulta GET /api/instances con Authorization Bearer');
  it.todo('muestra instancias disponibles y permisos actuales del usuario');
  it.todo('envía el array de VMIDs seleccionado al endpoint de permisos');
  it.todo('recarga la selección guardada al abrir nuevamente el usuario');
});

describe('FRN-08 - manejo de 403 en recursos protegidos', () => {
  it.todo('envía el JWT Bearer en las peticiones de instancias y permisos');
  it.todo('muestra un error amigable ante 403 sin borrar la sesión local');
  it.todo('mantiene al usuario en la vista después de un 403');
});
