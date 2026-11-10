# Resultados de las pruebas backend

**Fecha:** 10/11/2026. **Backend:** submódulo en `d0f3f57` (últimos commits de `origin/main`: `d0f3f57`/`df4cada` protección de VMID reservados y pertenencia al borrar; `d92a568` **`FIX-139`** reconciliación de tareas al arranque antes de abrir la API; `5b3ac80` **`FIX-137`** eliminación de endpoints fantasma e inclusión de auditorías globales del sistema; `4011b08` **`FIX-143`** vocabulario cerrado de `accion` y `resultado`, enums en Swagger y normalización histórica).

## 1. Suite interna del submódulo (`backend/`)

**Comando:** `go test -C backend ./...`

| Métrica | Valor |
|---|---:|
| Paquetes probados | **14** |
| Aprueban | **14 (100 %)** |
| Fallan | **0** |

Paquetes en verde: `cmd/api`, `cmd/proxmox-simulador`, `docs`, `internal/adapters/primary/http`, `internal/adapters/primary/http/middleware`, `internal/adapters/secondary/email`, `internal/adapters/secondary/memoria`, `internal/adapters/secondary/postgres`, `internal/adapters/secondary/proxmox`, `internal/adapters/secondary/redis`, `internal/core/domain/topology`, `internal/core/services`, `internal/infrastructure/crypto`, `migrations`.

---

## 2. Suite de aceptación e inspección (`test/back/`)

**Comando:** `go test -C test/back -v ./...`

### 2.1 Tareas de Backend verificadas en verde (100 % completadas en `terminado3.md`)

| Tarea | Prueba en `test/back` | Resultado | Detalle |
|---|---|:---:|---|
| `FIX-137` | `TestFIX137_CodigoInerteYAuditoriaDeSistema` (3/3 subtests) | ✅ **Pasa** | Commit `5b3ac80`: se eliminó el stub `PausarInstancia` y la ruta `/pause` de Swagger, `audit_repository.go` incluye eventos globales (`a.organizacion_id IS NULL`) y se limpiaron los tipos DTO huérfanos en `contrato_etapa1.go` |
| `FIX-139` | `TestFIX139_ReconciliacionDeTareasAlArranque` (3/3 subtests) | ✅ **Pasa** | Commit `d92a568`: `RecuperarYReconciliarAlArranque` en `seguimiento_tareas.go`, invocada en `cmd/api/main.go` antes de abrir el servidor HTTP, con pruebas unitarias en `seguimiento_tareas_test.go` |
| `FIX-143` | `TestFIX143_VocabularioAuditoriaContratoCerrado` (4/4 subtests) + `TestContratoVocabularioAuditoria` | ✅ **Pasa** | Commit `4011b08`: conjunto cerrado único `EXITO`/`FALLA`/`PENDIENTE` (`ResultadosAuditoriaValidos()`) y `AccionesAuditoriaValidas()` en `ports/audit_port.go`, eliminación de constantes muertas, validación `ErrAuditoriaFueraDeContrato`, enums en `docs/swagger.json` y `docs/contrato-etapa1.md`, test `TestContratoVocabularioAuditoria` y migración `migrations/20261104_vocabulario_auditoria.sql` |
| `FIX-134` | `TestFIX134_InstanciasPausadas` (4/4 subtests) | ✅ **Pasa** | `EstadoInstanciaPausada = "PAUSED"`, soporte de `stop` y `resume` sobre instancias pausadas y desbloqueo de `EliminarInstancia` |
| `FIX-135` | `TestFIX135_AlertasParadasExternas` (3/3 subtests) | ✅ **Pasa** | Reconciliación con tareas de Proxmox en `explicadaPorCentinela`, auditoría `SYSTEM_INSTANCE_STOPPED_EXTERNALLY` y variable `ALERT_IGNORE_EXTERNAL_SHUTDOWN` |
| `FIX-136` | `TestFIX136_UmbralesDeAlertaConfigurables` (4/4 subtests) | ✅ **Pasa** | `ConfigAlertas` inyectable desde variables de entorno `ALERT_THRESHOLD_*` y `ALERT_CYCLES_COUNT` |
| `FIX-128` | `TestFIX128VLAN_VlanCeroEsTraficoSinEtiquetar` (3/3 subtests) | ✅ **Pasa** | `*input.VLAN == 0` se normaliza a `nil` y no emite `,tag=0` en `net0` |
| `FIX-133` | `TestFIX133_CuotaVCPUSinOvercommit` (2/2 subtests) | ✅ **Pasa** | Escalado vertical permite overcommit de vCPU cuando $\Delta\text{RAM} \le 0$ |
| `FIX-138` | `TestFIX138_CaracteresEspecialesNIST` (1/1 subtest) | ✅ **Pasa** | `CaracteresEspeciales` ampliado según NIST SP 800-63B |
| `FIX-140` | `TestFIX140_ProveedorDeEmailVisible` (2/2 subtests) | ✅ **Pasa** | `GET /api/version` informa `emailProvider` y añade `X-Centinela-Email-Provider: mock` |
| `BAC-54` | `TestEtapa3BAC54_ParametrizacionEntornoServidoresHeterogeneos` | ✅ **Pasa** | Lectura dinámica de `PROXMOX_NODE`, `PROXMOX_URL` y `PROXMOX_PROTECTED_VMIDS` |
| `BAC-55` / `ADP-05` / `ADP-06` | `TestEtapa3BAC55_SwitchPerfilVersionSimulador`, `TestADP05_ModoDualSimulador`, `TestADP06_CatalogoDebian12Simulador` | ✅ **Pasa** | Simulador Proxmox dual PVE 8.4.21 / 9.2.2 y catálogo Debian 12 (`debian-12-standard_12.7-1_amd64.tar.zst`) |
| `ADP-01` a `ADP-04` | `TestADP01_...` a `TestADP04_...` | ✅ **Pasa** | Parámetros de aprovisionamiento QEMU/LXC, purga en borrado, parser defensivo de `memory` y `PROXMOX_PROTECTED_VMIDS` |
| `FIX-129` (`INF-11`) | `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS` (3/3 subtests) | ✅ **Pasa** | `docker/nginx-edge.conf` redirige HTTP->HTTPS con 301, TLS 1.2/1.3 + WSS/SSE sin buffering, y cabeceras HSTS (`includeSubDomains`), `X-Frame-Options "DENY"`, `X-Content-Type-Options "nosniff"` y `Content-Security-Policy` |

---

### 2.2 Pruebas de Integración y Seguridad del Sistema (100 % en verde)

| Suite de Integración | Resultado | Subtests Clave Verificados |
|---|:---:|---|
| `TestTareasBackendEIntegracion` | ✅ **PASS (100 %)** | `BAC-01` (esquema y usuario seed), `BAC-02` (bcrypt), `BAC-03` (JWT emitido y firmado), `BAC-04` (unificación de 400/401), `LOGIN-01` (TOTP), `BAC-10`/`BAC-11` (anti-reemplazo y persistencia), `BAC-12` (cambio obligatorio de clave temporal), `BAC-09` (roles), `BAC-05`/`BAC-06`/`BAC-06B` (gestión de usuarios con referencias opacas e `instanciasPermitidas`), `LOGIN-03` (protección y recorrido). |
| `TestHitoSeguridadDeSesionesYCookies` | ✅ **PASS (100 %)** | `BAC-17` (logout revoca atómicamente access y refresh en BD y registra auditoría `LOGOUT`), `SEC-01` (refresh token en cookie HttpOnly/Secure/SameSite=Strict), `SEC-01`/`SEC-02` (contratos HTTP de frontend), `FIX-28`/`FRN-13` (logout con frontend), `FIX-08` (formato estándar de error con `errorCode` y `message`). |
| `TestHitoControlDeAccesoBasadoEnRecursos` | ✅ **PASS (100 %)** | `BAC-07` (asignación y reemplazo atómico de permisos con referencias opacas `inst_...`), `FIX-16`/`BAC-08` (403 `INSTANCE_ACCESS_DENIED` sin delegar a Proxmox e `INSTANCE_PROTECTED` en VMID protegidos). |
| `TestAdversarialSystemIntegrity` | ✅ **PASS (100 %)** | Resistencia a inyecciones en parámetros Proxmox, mitigación de colisiones de VMID, exclusión de tokens huérfanos, atomicidad transaccional y rechazo con `409 INSTANCE_BUSY` ante ráfagas concurrentes. |
| `TestFixesPendientes` | ✅ **PASS (100 %)** | 9 subtests de regresión en verde (`FIX-107`, `FIX-110`, `FIX-111A`, `FIX-112`, `FIX-113`, `FIX-114`, `FIX-115A`, `FIX-116`, `FIX-117A`). |
| `TestPuenteEtapa1` | ✅ **PASS (100 %)** | 10 subtests de integración continua y compatibilidad de base de datos. |

---

### 2.3 Tareas de Infraestructura / Documentación pendientes en `actual.md`

| Tarea | Prueba en `test/back` | Resultado | Causa verificada |
|---|---|:---:|---|
| `FIX-145` | `setupIntegrationFixtures` (mitigación en harness) | 🟡 **Definida en `actual.md`** | Corrección de borrado estricto `403 INSTANCE_NOT_OWNED` para permitir borrar instancias preexistentes o huérfanas protegiendo únicamente `PROXMOX_PROTECTED_VMIDS`. |
| `FIX-123` | `TestFIX123_PlantillasEISOsEnProxmoxReal` | ⚪ **Omitida sin host físico** | Requiere variables de entorno del servidor Proxmox VE real (`PROXMOX_REAL_URL`, `PROXMOX_REAL_TOKEN_ID`, `PROXMOX_REAL_TOKEN_SECRET`). |
| `INF-14` / `ADP-09` / `ADP-10` | `TestEtapa3INF14_AuditoriaYSondaPVE8Real`, `TestADP09_ADP10_InfraestructuraPVE8` | ❌ **Falla (requiere host físico)** | Falta ejecutar `scripts/pve8-audit-probe.sh` en el hipervisor PVE 8.4.21 y configurar el token/volúmenes en el servidor físico. |

---

## 3. Regresión de Etapas 1, 2 y 3

Pasan sin regresiones todas las pruebas estáticas, de simulación y de integración de Etapas 1, 2 y 3:
- **Etapa 1 y Simulador:** `TestINF09WebSocketsNginxEdge`, `TestINF07CPermisosTokenProxmox`, `TestEtapa1Ola3BAC32`, `TestFIX81_OrdenMigracionesBaseLimpia`, `TestEtapa1SimuladorProxmox`, `TestEtapa2BAC50_SimuladorSnapshotsYRollback`.
- **Etapa 2:** `TestEtapa2BAC34_AdaptadorMetricasRRD`, `TestEtapa2BAC36_MetricasHistoricasEndpoint`, `TestEtapa2BAC39_SeguimientoCreacionYPermisos`, `TestEtapa2BAC40_EndpointConfigRecursos`, `TestEtapa2BAC41_AuditoriaAprovisionamientoYConfig`, `TestEtapa2BAC51_GuardaConcurrenciaInstanceBusy`, `TestEtapa2BAC52_BrokerSSEEventosAlertaGlobal`, `TestEtapa2BAC53_ParticionadoPreventivoEIndexacionAuditoria`, `TestFIX118_SinCredencialesPorDefecto`, y todos los subtests de `TestFIX107_OfuscacionDeVMIDs`, `TestFIX110_SaneamientoDeErrores`, `TestFIX112_RateLimitingAutenticacion`, `TestFIX113_PurgaDePermisosAlEliminarInstancia`, `TestFIX115A_ContrasenasDeHasta128Caracteres` y `TestFIX116_InyeccionDeParametrosProxmox`.
- **Etapa 3:** `TestEtapa3BAC42` (subtest 1), `TestEtapa3BAC43` (subtest 1), `TestEtapa3BAC44` (subtest 1), `TestEtapa3BAC45` (subtest 1), `TestEtapa3BAC46_WorkerEvaluadorUmbralesHisteresis` (2/2), `TestEtapa3BAC47_DespachoMulticanalAlertasSSEyCorreo` (2/2), `TestEtapa3BAC48_HardeningAutenticacionCookiesRateLimiting` (2/2), `TestEtapa3BAC54_ParametrizacionEntornoServidoresHeterogeneos`, `TestEtapa3BAC55_SwitchPerfilVersionSimulador`, y la totalidad de `etapa3_fixes_pendientes_test.go` (`TestFIX128VLAN`, `TestFIX133`..`TestFIX140`, `TestFIX143`).
