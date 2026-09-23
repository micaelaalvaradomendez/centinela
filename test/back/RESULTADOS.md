# Resultados de las pruebas backend

**Fecha:** 23/09/2026. **Backend:** `15032da` (`origin/main`). **Comando:** `go test -v -count=1 ./...` (2 corridas, mismo resultado, ~75 s).

| Métrica | Valor |
|---|---:|
| Casos | 28 |
| Aprueban | 24 |
| Fallan | **4** |
| Omitidos | 0 |
| Contenedores residuales | ninguno |

## Por caso

| Caso | Tarea | Resultado | Qué se verifica |
|---|---|---|---|
| BAC-01 crea esquema… | BAC-01 | ✅ | Tablas, seed del admin con hash y columna `rol` |
| BAC-05 protege permisos_instancia… | BAC-05 | ✅ | Restricción UNIQUE sobre `(usuario_id, vmid_proxmox)` |
| BAC-02 persiste bcrypt… | BAC-02 | ✅ | Hash `$2…` sin texto plano; credencial inválida → 401 |
| BAC-03 login emite JWT… | BAC-03 | ✅ | Claims `sub/jti/rol/org_id/exp`, HS256, token pre-auth |
| BAC-04 unifica errores… | BAC-04 | ✅ | 400 `INVALID_REQUEST` y 401 `AUTH_FAILED` con `message` |
| LOGIN-01 genera persiste y valida TOTP | LOGIN-01 | ✅ | QR PNG, secreto cifrado, códigos corto/incorrecto/vencido rechazados y refresh entregado (en JSON o cookie) |
| BAC-10 BAC-11 impiden reemplazar TOTP… | BAC-10/11 | ✅ | Sin QR nuevo con 2FA activo; login con el secreto persistido; anti-replay |
| BAC-12 cambio obligatorio… | BAC-12 | ✅ | 403 antes del cambio; `PUT /account/password`; re-login |
| BAC-09 expone roles | BAC-09 | ✅ | `ADMIN` y `OPERATOR` |
| BAC-05 BAC-06 BAC-06B gestiona usuarios… | BAC-06/06B/16 | ✅ | CRUD, filtros, permisos, 403 al operador en todos los endpoints y **sin `contrasenaTemp` en el JSON** |
| BAC-06 BAC-06B exigen el prefijo /api/admin/users | BAC-06 | ✅ | `/api/admin/users` → 200 y la ruta vieja `/api/users` → 404 |
| LOGIN-03 completa el recorrido… | LOGIN-03 | ✅ | Claims de acceso; 401 sin token y 403 con token pre-auth |
| LOGIN-04 circuito completo | LOGIN-04 | ✅ | Los 7 pasos con clave temporal real por correo, incluidos el 403 sobre instancia no asignada y el inventario filtrado a `[101]` |
| BAC-19 solicitud de recuperación… | BAC-19 | ✅ | 400 por formato; respuesta idéntica exista o no la cuenta; código de 6 dígitos; **el código anterior queda invalidado** |
| BAC-20 confirmación… | BAC-20 | ✅ | Código incorrecto, usado o vencido → 400; la clave vieja deja de funcionar; **sesiones previas revocadas** |
| BAC-16 clave temporal solo por correo | BAC-16 | ✅ | Sin campos de contraseña en el JSON de alta ni de reset; la clave del correo permite el login; hash bcrypt |
| BAC-21 contrato RealtimeEvent | BAC-21 | ✅ | Mismos campos y enums en `event_port.go` y `notifications.ts` |
| **BAC-18 auditoría append-only** | BAC-18 | ❌ | Ver abajo |
| BAC-13 reset de 2FA | BAC-13 | ✅ | 403 al operador; `totp_vinculado=false`; sesiones revocadas; re-enrolamiento; **el secreto viejo se rechaza** |
| BAC-15 reset de contraseña | BAC-15 | ✅ | 403 al operador; clave temporal por correo; la vieja → 401; `cambioContrasenaRequerido=true`; sesiones revocadas |
| BAC-07 GET/PUT permissions | BAC-07 | ✅ | Reemplazo atómico, deduplicación, 403 al operador en GET y PUT, auditoría `ASIGNAR_PERMISOS` |
| FIX-16 BAC-08 rechaza no asignadas | FIX-16 | ✅ | 403 `INSTANCE_ACCESS_DENIED` en GET, start y stop; el stub no recibe la orden; lo asignado y el ADMIN pasan |
| BAC-14 inventario normalizado y filtrado | BAC-14 | ✅ | 401 sin token; el ADMIN ve exactamente 101/102/103 con `name/type(vm/lxc)/node/status`; el operador ve solo lo asignado y existente; 404 `INSTANCE_NOT_FOUND` |
| **SEC-04 niveles de acceso** | SEC-04 | ❌ | Ver abajo |
| BAC-17 logout atómico | BAC-17 | ✅ | Logout sin Bearer → 401; con Bearer → 204; access → 401 `TOKEN_REVOKED`; refresh → 401; ambas sesiones inactivas; 1 registro `LOGOUT` nuevo |
| **SEC-01 refresh solo en cookie HttpOnly** | SEC-01 | ❌ | Ver abajo |
| **SEC-01 SEC-02 integración** | SEC-01/SEC-02 | ❌ | Ver abajo |
| FIX-08 estructura estándar de error | FIX-08 | ✅ | 4 casos con `errorCode` y `message` |

## Fallos y causa

Los 4 fallos son **del producto**. Las pruebas reflejan el criterio de éxito documentado.

### BAC-18: la tabla `auditoria` no es append-only a nivel de motor
```
password_recovery_acceptance_test.go:227: la base permitió UPDATE sobre auditoria con el usuario de la aplicación; la tabla no es append-only
password_recovery_acceptance_test.go:230: la base permitió DELETE sobre auditoria con el usuario de la aplicación; la tabla no es append-only
```
El criterio de `terminado.md` exige que `UPDATE` y `DELETE` fallen "a nivel de base de datos, no solo por convención de código". En el backend no hay trigger, `REVOKE` ni regla: la palabra "append-only" solo aparece en comentarios (`audit_port.go:106`, `models.go:86`). **Es una regresión de una tarea marcada como terminada**, así que requiere un fix de backend.

### SEC-04: no implementada
```
resource_access_acceptance_test.go:184: SEC-04 no implementada: la tabla permisos_instancia no tiene la columna nivel_acceso (entregable 1)
```
`domain.PermisoInstancia` no tiene `NivelAcceso` (`models.go:79-83`). `VerificarAcceso(ctx, userID, vmid)` no recibe el nivel y `RequireInstanceAccess(repo, param)` no tiene `requiredLevel`. Cuando exista la columna, la prueba ya verifica el resto: default `FULL_ACCESS`, CHECK de valores, READ_ONLY → GET 200 y start/stop 403, FULL_ACCESS → start 202.

### SEC-01: no implementada
```
session_security_acceptance_test.go:83: POST /auth/2fa/verify todavía devuelve refreshToken en el JSON (entregable 3)
session_security_acceptance_test.go:87: SEC-01 no implementada: POST /auth/2fa/verify no emite Set-Cookie con el refresh token (cabeceras: [])
```
Cuando exista la cookie, la prueba ya verifica el resto:
- los atributos `HttpOnly`, `SameSite` y `Path`;
- refresh solo con la cookie, sin `refreshToken` en la respuesta;
- 401 con una cookie inválida;
- logout que borra la cookie, y que esa cookie ya no renueve la sesión;
- que el token no aparezca en los logs.

### Integración SEC-01 / SEC-02: el backend rechaza las peticiones que ya envía el frontend
```
session_security_acceptance_test.go:146: POST /auth/refresh con body {} (petición real del frontend): esperado 200, recibido 400: {"errorCode":"INVALID_REQUEST","message":"Se requiere el campo refreshToken."}
session_security_acceptance_test.go:151: POST /auth/logout con Bearer y body {} (petición real del frontend): esperado 204, recibido 400: {"errorCode":"INVALID_REQUEST","message":"El formato de la petición es incorrecto. Se requiere refreshToken."}
session_security_acceptance_test.go:154: tras el logout iniciado desde el frontend el access token debe quedar revocado: esperado 401, recibido 200
```
El frontend (`origin/main`, SEC-02) ya no envía el refresh token. Con ambos `main` desplegados, **el logout de la interfaz no revoca la sesión en el servidor y la renovación silenciosa falla**. Se resuelve con SEC-01.

## Lo que la suite no cubre

- `502 PROXMOX_UNAVAILABLE` en `/instances`: el stub siempre responde. El mapeo está en `instance_handler.go:39-48`.
- El envío real de correo: está fuera de alcance según BAC-16, que usa `MockEmailService`.
