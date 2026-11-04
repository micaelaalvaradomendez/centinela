# Resultados de las pruebas backend

**Fecha:** 04/11/2026. **Backend:** submódulo en `4011b08` (últimos commits de `main`: `4011b08` **`FIX-143`** vocabulario cerrado de `accion` y `resultado`, enums en Swagger y normalización histórica; `08126ae` **`FIX-136`** parametrización de umbrales y persistencia de alertas por entorno; `5f8bb5e` **`FIX-135`** alertas por paradas externas de instancias; `42093c3` **`FIX-134`** soporte para instancias pausadas en inventario/energía/borrado).

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

**Comando:** `go test -v ./...`

### 2.1 Tareas recién verificadas en verde (movidas a `terminado3.md`)

| Tarea | Prueba en `test/back` | Resultado | Detalle |
|---|---|:---:|---|
| `FIX-143` | `TestFIX143_VocabularioAuditoriaContratoCerrado` (4/4 subtests) + `TestContratoVocabularioAuditoria` | ✅ **Pasa** | Commit `4011b08`: conjunto cerrado único `EXITO`/`FALLA`/`PENDIENTE` (`ResultadosAuditoriaValidos()`) y `AccionesAuditoriaValidas()` en `ports/audit_port.go`, eliminación de constantes muertas, validación `ErrAuditoriaFueraDeContrato`, enums en `docs/swagger.json` y `docs/contrato-etapa1.md`, test `TestContratoVocabularioAuditoria` y migración `migrations/20261104_vocabulario_auditoria.sql` |
| `FIX-134` | `TestFIX134_InstanciasPausadas` (4/4 subtests) | ✅ **Pasa** | `EstadoInstanciaPausada = "PAUSED"`, mapeo `"paused" -> PAUSED` en `instance_service.go`, guarda previa en `EjecutarAccionEnergia` y `EliminarInstancia` que rechaza `PAUSED` con `409 INSTANCE_INVALID_STATE`, y soporte de `resume` en `proxmox/client.go` |
| `FIX-135` | `TestFIX135_AlertasParadasExternas` (3/3 subtests) | ✅ **Pasa** | Detección de transición `RUNNING -> STOPPED` sin tarea activa ni orden en Redis (`instance:stop_intent:{vmid}`), auditoría `EXTERNAL_STOP_DETECTED` y alerta crítica SSE + email `EXTERNAL_INSTANCE_STOPPED` |
| `FIX-136` | `TestFIX136_UmbralesDeAlertaConfigurables` (4/4 subtests) | ✅ **Pasa** | Lectura de `ALERT_CPU_THRESHOLD`, `ALERT_RAM_THRESHOLD`, `ALERT_STORAGE_THRESHOLD`, `ALERT_HYSTERESIS_MARGIN` y `ALERT_EVALUATION_SAMPLES` con validación $0 < \text{threshold} \le 100$ y $\text{hysteresis} < \text{threshold}$, e inyección en `NuevoEvaluadorAlertasService` |
| `FIX-128` | `TestFIX128VLAN_VlanCeroEsTraficoSinEtiquetar` (3/3 subtests) | ✅ **Pasa** | `VlanTag == 0` se normaliza a `nil` y no emite `,tag=0` en `net0` |
| `FIX-133` | `TestFIX133_CuotaVCPUSinOvercommit` (2/2 subtests) | ✅ **Pasa** | `FreeCores` descuenta los vCPU asignados a instancias del nodo sin overcommit |
| `FIX-138` | `TestFIX138_CaracteresEspecialesNIST` (1/1 subtest) | ✅ **Pasa** | `ValidarComplejidad` acepta cualquier carácter imprimible no alfanumérico (`unicode.IsPunct` / `unicode.IsSymbol`) |
| `FIX-140` | `TestFIX140_ProveedorDeEmailVisible` (2/2 subtests) | ✅ **Pasa** | `cmd/api/main.go` registra el proveedor de correo activo y advierte si `EMAIL_PROVIDER=mock` en producción |
| `BAC-54` | `TestEtapa3BAC54_ParametrizacionEntornoServidoresHeterogeneos` | ✅ **Pasa** | Parametrización de entorno para servidores heterogéneos (`PROXMOX_STORAGE_DEFAULT`, `PROXMOX_BRIDGE_DEFAULT`, etc.) |
| `BAC-55` / `ADP-05` / `ADP-06` | `TestEtapa3BAC55_SwitchPerfilVersionSimulador`, `TestADP05_ModoDualSimulador`, `TestADP06_CatalogoDebian12Simulador` | ✅ **Pasa** | Simulador Proxmox dual PVE 7.4 / 8.4 y catálogo Debian 12 (`debian-12-standard_12.7-1_amd64.tar.zst`) |
| `ADP-01` a `ADP-04` | `TestADP01_...` a `TestADP04_...` | ✅ **Pasa** | Parámetros de aprovisionamiento QEMU/LXC (`ostype=l26`, `scsihw=virtio-scsi-single`), purga en borrado, parser defensivo de `memory` y `PROXMOX_PROTECTED_VMIDS` |

---

### 2.2 Tareas pendientes en `actual.md` (fallos esperados en `test/back`)

| Tarea | Prueba en `test/back` | Resultado | Causa verificada |
|---|---|:---:|---|
| `FIX-137` | `TestFIX137_CodigoInerteYAuditoriaDeSistema` (0/3) | ❌ **Falla** | 1) `PausarInstancia` (`501 Not Implemented`) sigue expuesto en `instance_handler.go` y Swagger; 2) `audit_repository.go` oculta eventos con `usuario_id IS NULL AND organizacion_id IS NULL` cuando el admin tiene `organizacion_id`; 3) `contrato_etapa1.go` mantiene tipos DTO sin referencias |
| `FIX-139` (`FIX-123`) | `TestFIX139_ReconciliacionDeTareasAlArranque` (0/3) | ❌ **Falla** | No existe `RecuperarYReconciliarAlArranque` en `seguimiento_tareas.go` ni su invocación en `cmd/api/main.go`; tareas `CREATE`/`DELETE` en `RUNNING` al reiniciar no disparan sus callbacks de permisos |
| `FIX-129` (`INF-11`) | `TestEtapa3INF11_NginxProduccionTLS13HSTSCSPyWSS` (2/3 pasan, 1/3 falla) | ❌ **Falla parcial** | `docker/nginx-edge.conf` redirige HTTP->HTTPS y soporta TLS 1.2/1.3 + WSS/SSE, pero faltan `Strict-Transport-Security` (`includeSubDomains`), `Content-Security-Policy`, `X-Frame-Options` y `X-Content-Type-Options` |
| `INF-14` / `ADP-09` / `ADP-10` | `TestEtapa3INF14_AuditoriaYSondaPVE8Real`, `TestADP09_ADP10_InfraestructuraPVE8` | ❌ **Falla** | Falta crear `scripts/pve8-audit-probe.sh` y completar la documentación de auditoría PVE 8.4 y tipos de storage compatibles con snapshots en `documentacion/api-proxmox.md` |

---

## 3. Regresión de Etapas 1, 2 y 3

Pasan sin regresiones todas las pruebas estáticas y de simulador de Etapas 1, 2 y 3:
- **Etapa 1 y Simulador:** `TestINF09WebSocketsNginxEdge`, `TestINF07CPermisosTokenProxmox`, `TestEtapa1Ola3BAC32`, `TestFIX81_OrdenMigracionesBaseLimpia`, `TestEtapa1SimuladorProxmox`, `TestEtapa2BAC50_SimuladorSnapshotsYRollback`.
- **Etapa 2:** `TestEtapa2BAC34_AdaptadorMetricasRRD`, `TestEtapa2BAC36_MetricasHistoricasEndpoint`, `TestEtapa2BAC39_SeguimientoCreacionYPermisos`, `TestEtapa2BAC40_EndpointConfigRecursos`, `TestEtapa2BAC41_AuditoriaAprovisionamientoYConfig`, `TestEtapa2BAC51_GuardaConcurrenciaInstanceBusy`, `TestEtapa2BAC52_BrokerSSEEventosAlertaGlobal`, `TestEtapa2BAC53_ParticionadoPreventivoEIndexacionAuditoria`, `TestFIX118_SinCredencialesPorDefecto`, y todos los subtests estáticos de `TestFIX107_OfuscacionDeVMIDs`, `TestFIX110_SaneamientoDeErrores`, `TestFIX112_RateLimitingAutenticacion`, `TestFIX113_PurgaDePermisosAlEliminarInstancia`, `TestFIX115A_ContrasenasDeHasta128Caracteres` y `TestFIX116_InyeccionDeParametrosProxmox`.
- **Etapa 3:** `TestEtapa3BAC42` (subtest 1), `TestEtapa3BAC43` (subtest 1), `TestEtapa3BAC44` (subtest 1), `TestEtapa3BAC45` (subtest 1), `TestEtapa3BAC46_WorkerEvaluadorUmbralesHisteresis` (2/2), `TestEtapa3BAC47_DespachoMulticanalAlertasSSEyCorreo` (2/2), `TestEtapa3BAC48_HardeningAutenticacionCookiesRateLimiting` (2/2).
