# Resultados de las pruebas backend

**Fecha:** 01/10/2026. **Backend:** submódulo en `860b3c9` (último commit de `main`). **Comando:** `go test -v -count=1 ./...`. 2 corridas completas con el mismo resultado (~2,5 min) y sin contenedores residuales. El caso `FIX-37` se agregó después y se corrió aparte, junto con su bloque.

| Métrica | Valor |
|---|---:|
| Casos | 45 |
| Aprueban | **37** |
| Fallan | **8** |
| Omitidos | 0 |

## Por caso

| Caso | Tarea | Resultado |
|---|---|---|
| BAC-01 a BAC-12, BAC-09, BAC-05/06/06B, prefijo `/api/admin/users`, LOGIN-03 | regresión | ✅ 12/12 |
| LOGIN-04 circuito desde Go | regresión | ✅ |
| LOGIN-04 front ↔ back (`TestLOGIN04IntegracionFrontBack`) | LOGIN-04 | ✅ 10/10 pasos |
| INF-05 CORS | regresión | ✅ |
| **FIX-31 TLS de Nginx versionado** | FIX-31 | ❌ |
| INF-06A Redis local | regresión | ✅ |
| **INF-08B credenciales SMTP** | INF-08B / FIX-34 | ❌ (valores de ejemplo) |
| **FIX-34 credenciales de referencia contra Brevo** (nuevo) | FIX-34 | ❌ (`525 Unauthorized IP address`) |
| BAC-16B envío por SMTP (STARTTLS obligatorio) | BAC-16B | ✅ **nuevo** |
| BAC-17A adaptador de Redis, BAC-17B una sesión = un registro | regresión | ✅ 2/2 |
| **BAC-18B índice parcial y particiones** | BAC-18B | ❌ |
| BAC-19, BAC-20, BAC-16, BAC-21, BAC-18, FIX-17, BAC-13, BAC-15 | regresión | ✅ 8/8 |
| BAC-07, FIX-16/BAC-08 (incluido `INSTANCE_PROTECTED`), BAC-14, SEC-04 (2) | regresión | ✅ 5/5 |
| **FIX-37 nivel de acceso en `GET /account/profile`** (nuevo) | FIX-37 / FIX-30 | ❌ |
| **BAC-21B** (3 casos) | BAC-21B | ❌ 0/3 |
| BAC-21C (2 casos) | regresión | ✅ 2/2 |
| BAC-17, SEC-01, integración SEC-01/SEC-02, FIX-08 | regresión | ✅ 4/4 |
| FIX-28 integración del logout del frontend | FIX-28 | ✅ **nuevo** |

## Fallos y causa

| Caso | Causa |
|---|---|
| FIX-31 | No hay ninguna configuración de Nginx versionada con `listen 443 ssl` |
| INF-08B | `SMTP_USER` y `SMTP_PASS` tienen valores de ejemplo en `backend/.env.example` |
| FIX-34 | Brevo rechaza la IP de la máquina de pruebas (IP autorizadas activadas en la cuenta); además, `.env.example` no coincide con `smtp-brevo.env` |
| BAC-18B | Sin índice parcial (sobre `jti_access`), sin particiones y sin índices compuestos |
| FIX-37 | `UsuarioDetalleDTO` solo tiene `instanciasPermitidas`; falta `permisos [{ vmid, nivelAcceso }]` |
| BAC-21B | Faltan los campos nuevos en `GET /instances`; `/status/:action` → 404; `DELETE /instances/:vmid` → 405 |

## Cambios en la suite

- **Mailpit con STARTTLS obligatorio:** usa el certificado de `mailpit-tls/` y `backend-smtp` confía en esa CA mediante `SSL_CERT_FILE`.
- **Extracción de la clave temporal:** se toma del cuerpo del correo sin limpiar, porque puede contener `@`, `*`, `<`, `=`, etc.
- **Marca del código de recuperación:** ahora es `"código de seguridad es:"` (texto nuevo del mock).
- **Caso nuevo `FIX-34`:** autentica contra Brevo con `smtp-brevo.env` (credenciales reales, versionadas a propósito) y compara `backend/.env.example` contra ese archivo.
- **Caso nuevo `FIX-37`:** verifica que el perfil informe el nivel de acceso por instancia.
- **Integración FIX-28:** el logout se envía con Bearer y cookie, como lo hace el frontend desde `cfb88f7`.

## Lo que la suite no cubre

- `INF-06B` e `INF-08A`: son configuración del servidor.
- La conexión a un SMTP y a un Proxmox reales. La comparación entre el simulador y el Proxmox real se hizo a mano; ver `../informe.md` §4.
- La purga horaria y el tiempo de consulta de BAC-18B.
