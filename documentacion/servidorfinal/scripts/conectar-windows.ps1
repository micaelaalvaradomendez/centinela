# ============================================================================
# conectar-windows.ps1 - Asistente paso a paso para conectar al nodo G1
#                        (pve-is3-g1) desde Windows.
#
# Es un asistente interactivo, sin menu: en cada paso la persona lee la
# explicacion, hace la accion y vuelve a la ventana a presionar ENTER.
#   1) Tailscale instalado (lo instala con winget si falta).
#   2) Iniciar sesion en Tailscale.
#   3) Aceptar la invitacion del nodo.
#   4) Prueba de conexion (ping y servicio web). Automatica, sin pausa.
#   5) Resumen con [OK]/[X].
#
# Uso (PowerShell), desde la carpeta de conexion (la que contiene la Guia de conexion (.docx)):
#   powershell -ExecutionPolicy Bypass -File .\scripts\conectar-windows.ps1
# ============================================================================

$ErrorActionPreference = "Continue"

# --- Datos del nodo (NO cambiar sin coordinar con infraestructura) ----------
$NodeName  = "pve-is3-g1"
$NodeIp    = "100.116.49.111"
# URL del servicio web: SOLO se usa para comprobar que responde. No se imprime,
# no se muestra como destino y no lleva usuario/realm/contrasena: la
# administracion vive en los materiales de infraestructura del grupo.
$WebUrl    = "https://${NodeIp}:8006/"
$ShareLink = "https://login.tailscale.com/admin/invite/7hgHA4tQXPYjS22vFJx421"

# --- Estado de cada chequeo ------------------------------------------------
$script:OkTs      = $false
$script:OkSession = $false
$script:OkNode    = $false
$script:OkWeb     = $false
$script:OkPing    = $false

# Binario de Tailscale (se completa en Find-Tailscale).
$script:Ts = $null

# --- Estado estructurado de Tailscale (ConvertFrom-Json nativo) -------------
$script:StObj     = $null
$script:StParseOk = $false
$script:LoginProc = $null

# Acelera las esperas en pruebas con mocks (nunca se define en uso real).
function Start-Pausa {
    param([int]$Segundos)
    if ($env:ING3_TEST_FAST -eq "1") { Start-Sleep -Milliseconds 5; return }
    Start-Sleep -Seconds $Segundos
}

# ============================================================================
# Estilo de consola (banner al iniciar, tarjeta al cerrar).
# Se anima solo con terminal interactiva y soporte VT; si no, arte en texto
# plano o el encabezado simple. Decode y colores de las demos (24 bits).
# ============================================================================
$script:Esc = [char]27

# --- Habilitar VT (ENABLE_VIRTUAL_TERMINAL_PROCESSING) en esta consola -------
function Enable-VtConsole {
    try {
        if (-not ("Win32Vt.ConsoleMode" -as [type])) {
            Add-Type -Namespace Win32Vt -Name ConsoleMode -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("kernel32.dll", SetLastError = true)]
public static extern System.IntPtr GetStdHandle(int nStdHandle);
[System.Runtime.InteropServices.DllImport("kernel32.dll", SetLastError = true)]
public static extern bool GetConsoleMode(System.IntPtr hConsoleHandle, out uint lpMode);
[System.Runtime.InteropServices.DllImport("kernel32.dll", SetLastError = true)]
public static extern bool SetConsoleMode(System.IntPtr hConsoleHandle, uint dwMode);
'@ -ErrorAction Stop
        }
        $h = [Win32Vt.ConsoleMode]::GetStdHandle(-11)
        [uint32]$mode = 0
        if ([Win32Vt.ConsoleMode]::GetConsoleMode($h, [ref]$mode)) {
            # 0x0004 = ENABLE_VIRTUAL_TERMINAL_PROCESSING
            [void][Win32Vt.ConsoleMode]::SetConsoleMode($h, [uint32]($mode -bor 0x0004))
            return $true
        }
    } catch { }
    return $false
}

$script:UiEstilo = $false
if (-not [Console]::IsOutputRedirected) { $script:UiEstilo = $true }
$script:VtOk = $false
if ($script:UiEstilo) { $script:VtOk = Enable-VtConsole }

# Si la terminal es angosta (< 68 columnas), modo simple.
if ($script:UiEstilo) {
    try { if ([Console]::WindowWidth -lt 68) { $script:UiEstilo = $false; $script:VtOk = $false } }
    catch { $script:UiEstilo = $false; $script:VtOk = $false }
}

# Limpia la pantalla visible sin borrar el scrollback (ESC[H + ESC[2J). No se
# usa Clear-Host: en Windows Terminal puede emitir ESC[3J.
function Clear-PantallaVisible {
    if ($script:VtOk) { [Console]::Write("$script:Esc[H$script:Esc[2J"); return }
    try { Clear-Host } catch { }
}

# Consola limpia al arrancar (solo con terminal interactiva).
if ($script:UiEstilo) { Clear-PantallaVisible }

$script:Reset = ""
$script:Bold = ""
if ($script:VtOk) {
    $script:Reset = "$script:Esc[0m"
    $script:Bold = "$script:Esc[1m"
}

function Get-Fg([int]$r, [int]$g, [int]$b) { if (-not $script:VtOk) { return "" } return "$script:Esc[38;2;$r;$g;${b}m" }

# Margen de lectura del contenido (nada queda pegado al borde).
function Marg {
    param([string]$Msg)
    if ($script:UiEstilo) { Write-Host "$Margen$Msg" } else { Write-Host $Msg }
}

# Enlaces: en azul para distinguirlos del texto comun. $Pre = prefijo literal,
# $Url = direccion (azul), $Post = sufijo literal opcional.
function MargUrl {
    param([string]$Pre, [string]$Url, [string]$Post = "")
    if ($script:VtOk) {
        Marg "$Pre$(Get-Fg 90 170 255)$Url$script:Reset$Post"
    } else {
        Marg "$Pre$Url$Post"
    }
}

# Estados de linea: [OK] verde, [X] rojo, [!] ambar, [..] gris, [i] celeste.
function Estado {
    param([string]$Tag, [string]$Msg)
    if ($script:VtOk) {
        $rgb = switch ($Tag) {
            "OK" { @(0,230,130) } "X" { @(255,90,90) } "!" { @(255,196,0) } ".." { @(150,150,160) } "i" { @(0,200,255) }
            default { @(255,255,255) }
        }
        [Console]::Write("$Margen$(Get-Fg $rgb[0] $rgb[1] $rgb[2])[$Tag]$script:Reset $Msg`n")
    } else {
        $col = switch ($Tag) {
            "OK" { "Green" } "X" { "Red" } "!" { "Yellow" } ".." { "Gray" } "i" { "Cyan" }
            default { "White" }
        }
        Write-Host "[$Tag] $Msg" -ForegroundColor $col
    }
}

# Cabecera de paso: separadores en gris y titulo destacado.
function Seccion {
    param([string]$Titulo)
    if ($script:VtOk) {
        [Console]::Write("$Margen$(Get-Fg 95 95 105)------------------------------------------$script:Reset`n")
        [Console]::Write("$Margen$script:Bold$(Get-Fg 220 220 240) $Titulo$script:Reset`n")
        [Console]::Write("$Margen$(Get-Fg 95 95 105)------------------------------------------$script:Reset`n")
    } else {
        Write-Host "------------------------------------------"
        Write-Host " $Titulo"
        Write-Host "------------------------------------------"
    }
}

# Cabecera del resumen final (borde con =).
function SeccionFin {
    param([string]$Titulo)
    if ($script:VtOk) {
        [Console]::Write("$Margen$(Get-Fg 95 95 105)==========================================$script:Reset`n")
        [Console]::Write("$Margen$script:Bold$(Get-Fg 220 220 240) $Titulo$script:Reset`n")
        [Console]::Write("$Margen$(Get-Fg 95 95 105)==========================================$script:Reset`n")
    } else {
        Write-Host "=========================================="
        Write-Host " $Titulo"
        Write-Host "=========================================="
    }
}

# Linea de exito suelta (sin etiqueta), en verde.
function AvisoOk {
    param([string]$Msg)
    if ($script:VtOk) { [Console]::Write("$Margen$(Get-Fg 0 230 130)$Msg$script:Reset`n") }
    else { Write-Host $Msg -ForegroundColor Green }
}

# Gradiente de la barra (se arma una sola vez por corrida).
$script:GradLista = $false
$script:GradBarra = @()

# Barra de progreso con brillo (estilo demos). $Msg = mensaje ; $Segs = segundos.
function Barra-Progreso {
    param([string]$Msg, [int]$Segs)
    if (-not $script:VtOk) { return }
    if (-not $script:GradLista) {
        for ($j = 0; $j -lt 24; $j++) {
            $t = [int][math]::Floor(100 * $j / 24)
            $r = [int][math]::Floor((255 - 0) * $t / 100)
            $g = [int][math]::Floor(229 + (0 - 229) * $t / 100)
            $b = [int][math]::Floor(255 + (128 - 255) * $t / 100)
            if ($r -lt 0) { $r = 0 } elseif ($r -gt 255) { $r = 255 }
            if ($g -lt 0) { $g = 0 } elseif ($g -gt 255) { $g = 255 }
            if ($b -lt 0) { $b = 0 } elseif ($b -gt 255) { $b = 255 }
            $script:GradBarra += (Get-Fg ([int]$r) ([int]$g) ([int]$b))
        }
        $script:GradLista = $true
    }
    $width = 24
    $pasos = $Segs * 12
    for ($i = 0; $i -le $pasos; $i++) {
        $pct = [int][math]::Floor(100 * $i / $pasos)
        $filled = [int][math]::Floor($width * $i / $pasos)
        $shine = [int][math]::Floor(($i * ($width + 8)) / $pasos) - 8
        $barra = ""
        for ($k = 0; $k -lt $width; $k++) {
            if ($k -lt $filled) {
                if ($k -ge $shine -and $k -lt ($shine + 5)) {
                    $barra += "$(Get-Fg 255 236 236)▓$script:Reset"
                } elseif ($k -eq ($filled - 1) -and ($i % 2) -eq 0) {
                    $barra += "$(Get-Fg 255 150 152)▓$script:Reset"
                } else {
                    $barra += "$($script:GradBarra[$k])█$script:Reset"
                }
            } else {
                $barra += "$(Get-Fg 70 70 90)░$script:Reset"
            }
        }
        [Console]::Write("`r$Margen$(Get-Fg 160 160 190)$Msg [$barra]$script:Reset $('{0,3}' -f $pct)%$script:Esc[K")
        Start-Sleep -Milliseconds 80
    }
    [Console]::Write("`n")
}

# Corre un comando mostrando spinner; devuelve su salida como texto.
function Invoke-Spin {
    param([string]$Msg, [string]$File, [string[]]$ArgList, [switch]$OcultarErr)
    if (-not $script:VtOk) {
        if ($OcultarErr) { return ((& $File @ArgList 2>$null) | Out-String) }
        return ((& $File @ArgList 2>&1) | Out-String)
    }
    $tmpO = [System.IO.Path]::GetTempFileName()
    $tmpE = [System.IO.Path]::GetTempFileName()
    try {
        try {
            $p = Start-Process -FilePath $File -ArgumentList $ArgList -NoNewWindow -RedirectStandardOutput $tmpO -RedirectStandardError $tmpE -PassThru -ErrorAction Stop
        } catch {
            if ($OcultarErr) { return ((& $File @ArgList 2>$null) | Out-String) }
            return ((& $File @ArgList 2>&1) | Out-String)
        }
        $frames = @("⠋","⠙","⠹","⠸","⠼","⠴","⠦","⠧","⠇","⠏")
        $i = 0
        while (-not $p.HasExited -or $i -lt 11) {
            [Console]::Write("`r$Margen$(Get-Fg 0 200 255)$($frames[$i % 10])$script:Reset $Msg$script:Esc[K")
            Start-Sleep -Milliseconds 90
            $i++
        }
        [Console]::Write("`r$script:Esc[K")
        $texto = ""
        if (Test-Path $tmpO) { $texto += (Get-Content $tmpO -Raw -ErrorAction SilentlyContinue) }
        if ((-not $OcultarErr) -and (Test-Path $tmpE)) { $texto += (Get-Content $tmpE -Raw -ErrorAction SilentlyContinue) }
        return $texto
    } finally {
        Remove-Item $tmpO, $tmpE -ErrorAction SilentlyContinue
    }
}

# Construye el gradiente de la barra (una sola vez por corrida).
function Ensure-GradBarra {
    if ($script:GradLista) { return }
    for ($j = 0; $j -lt 24; $j++) {
        $t = [int][math]::Floor(100 * $j / 24)
        $r = [int][math]::Floor((255 - 0) * $t / 100)
        $g = [int][math]::Floor(229 + (0 - 229) * $t / 100)
        $b = [int][math]::Floor(255 + (128 - 255) * $t / 100)
        if ($r -lt 0) { $r = 0 } elseif ($r -gt 255) { $r = 255 }
        if ($g -lt 0) { $g = 0 } elseif ($g -gt 255) { $g = 255 }
        if ($b -lt 0) { $b = 0 } elseif ($b -gt 255) { $b = 255 }
        $script:GradBarra += (Get-Fg ([int]$r) ([int]$g) ([int]$b))
    }
    $script:GradLista = $true
}

# Dibuja una linea de barra (sin newline) con el estilo de Barra-Progreso.
function Write-BarraLinea {
    param([string]$Msg, [int]$Pct, [int]$Frame)
    $width = 24
    $filled = [int][math]::Floor($width * $Pct / 100)
    if ($filled -gt $width) { $filled = $width }
    $shine = (($Frame * 3) % ($width + 8)) - 8
    $barra = ""
    for ($k = 0; $k -lt $width; $k++) {
        if ($k -lt $filled) {
            if ($k -ge $shine -and $k -lt ($shine + 5)) {
                $barra += "$(Get-Fg 255 236 236)▓$script:Reset"
            } elseif ($k -eq ($filled - 1) -and ($Frame % 2) -eq 0) {
                $barra += "$(Get-Fg 255 150 152)▓$script:Reset"
            } else {
                $barra += "$($script:GradBarra[$k])█$script:Reset"
            }
        } else {
            $barra += "$(Get-Fg 70 70 90)░$script:Reset"
        }
    }
    [Console]::Write("`r$Margen$(Get-Fg 160 160 190)$Msg [$barra]$script:Reset $('{0,3}' -f $Pct)%$script:Esc[K")
}

# Corre un comando con barra estilo Barra-Progreso; devuelve su ExitCode.
function Invoke-ConBarra {
    param([string]$Msg, [string]$File, [string[]]$ArgList)
    if (-not $script:VtOk) {
        Write-Host "[..] $Msg"
        $tmpO = [System.IO.Path]::GetTempFileName()
        $tmpE = [System.IO.Path]::GetTempFileName()
        try {
            $p = Start-Process -FilePath $File -ArgumentList $ArgList -NoNewWindow -RedirectStandardOutput $tmpO -RedirectStandardError $tmpE -PassThru -ErrorAction Stop
            $p.WaitForExit()
            return $p.ExitCode
        } catch {
            try { & $File @ArgList 2>$null | Out-Null; return $LASTEXITCODE } catch { return 1 }
        } finally {
            Remove-Item $tmpO, $tmpE -ErrorAction SilentlyContinue
        }
    }
    Ensure-GradBarra
    $tmpO = [System.IO.Path]::GetTempFileName()
    $tmpE = [System.IO.Path]::GetTempFileName()
    try {
        try {
            $p = Start-Process -FilePath $File -ArgumentList $ArgList -NoNewWindow -RedirectStandardOutput $tmpO -RedirectStandardError $tmpE -PassThru -ErrorAction Stop
        } catch {
            try { & $File @ArgList 2>$null | Out-Null; return $LASTEXITCODE } catch { return 1 }
        }
        $i = 0
        while (-not $p.HasExited -or $i -lt 11) {
            $pct = [int][math]::Floor(95 * $i / ($i + 12))
            Write-BarraLinea $Msg $pct $i
            Start-Sleep -Milliseconds 90
            $i++
        }
        Write-BarraLinea $Msg 97 $i; Start-Sleep -Milliseconds 90; $i++
        Write-BarraLinea $Msg 99 $i; Start-Sleep -Milliseconds 90; $i++
        Write-BarraLinea $Msg 100 $i
        [Console]::Write("`n")
        return $p.ExitCode
    } finally {
        Remove-Item $tmpO, $tmpE -ErrorAction SilentlyContinue
    }
}

# Anima la barra mientras el proceso corre; minimo visible ~1s.
function Watch-ConBarra {
    param([string]$Msg, $Proc)
    if (-not $script:VtOk) {
        if ($Proc) { try { $Proc.WaitForExit() } catch { } }
        return
    }
    Ensure-GradBarra
    $i = 0
    $running = $true
    while ($running -or $i -lt 11) {
        $pct = [int][math]::Floor(95 * $i / ($i + 12))
        Write-BarraLinea $Msg $pct $i
        Start-Sleep -Milliseconds 90
        $i++
        if ($Proc) { try { $running = -not $Proc.HasExited } catch { $running = $false } } else { $running = $false }
    }
    Write-BarraLinea $Msg 97 $i; Start-Sleep -Milliseconds 90; $i++
    Write-BarraLinea $Msg 99 $i; Start-Sleep -Milliseconds 90; $i++
    Write-BarraLinea $Msg 100 $i
    [Console]::Write("`n")
}

# URL del MSI segun la arquitectura (ARM64: x86/arm64 si existe; default amd64).
function Get-TailscaleMsiUrl {
    $arch = "$env:PROCESSOR_ARCHITECTURE"
    $cands = @()
    if ($arch -eq "ARM64") {
        $cands = @(
            "https://pkgs.tailscale.com/stable/tailscale-setup-latest-x86.msi",
            "https://pkgs.tailscale.com/stable/tailscale-setup-latest-arm64.msi",
            "https://pkgs.tailscale.com/stable/tailscale-setup-latest-amd64.msi"
        )
    } else {
        $cands = @("https://pkgs.tailscale.com/stable/tailscale-setup-latest-amd64.msi")
    }
    foreach ($u in $cands) {
        try {
            $code = & curl.exe -L -s -o NUL -w "%{http_code}" -I $u 2>$null
            if ("$code".Trim() -eq "200") { return $u }
        } catch { }
    }
    return "https://pkgs.tailscale.com/stable/tailscale-setup-latest-amd64.msi"
}

# Margen izquierdo del bloque del logo (4 columnas) y lineas de identidad.
$Margen = "    "
$IdentFilas = @(
    "Ingreso a la red de Tailscale",
    "Carrera Ingenieria de Software III",
    "Universidad Nacional de Tierra del Fuego, Antartida e Islas del Atlantico Sur",
    "Hecho por Lencina Nicolas"
)

# Banner CENTINELA (letras de bloque) y rampa de color de este nodo.
$BannerFilas = @(
    ' ████ █████ ██  ██ ██████ ████ ██  ██ █████ ██     ████ ',
    '██    ██    ███ ██   ██    ██  ███ ██ ██    ██    ██  ██',
    '██    ████  ██ ███   ██    ██  ██ ███ ████  ██    ██████',
    '██    ██    ██  ██   ██    ██  ██  ██ ██    ██    ██  ██',
    ' ████ █████ ██  ██   ██   ████ ██  ██ █████ █████ ██  ██'
)

function Mostrar-Banner {
    param([string]$Sub)
    if (-not $script:UiEstilo) {
        Write-Host "=========================================="
        Write-Host " $Sub"
        Write-Host "=========================================="
        foreach ($l in $IdentFilas) { Write-Host " $l" }
        Write-Host ""
        return
    }
    if (-not $script:VtOk) {
        # Terminal sin VT: arte en texto plano (sin secuencias).
        Write-Host ""
        Write-Host ""
        foreach ($l in $BannerFilas) { Write-Host "$Margen$l" }
        Write-Host ""
        Write-Host ""
        foreach ($l in $IdentFilas) { Write-Host "$Margen$l" }
        Write-Host ""
        Write-Host "$Margen$Sub"
        Write-Host ""
        return
    }
    $cGl = Get-Fg 90 90 110
    $shades = @("░","▒","▓")
    $rnd = New-Object System.Random
    $cols = $BannerFilas[0].Length
    $GCol = @()
    for ($j = 0; $j -lt $cols; $j++) {
        $t = [int][math]::Floor(100 * $j / $cols)
        $r = (255 - 0) * $t / 100
        $g = 229 + (0 - 229) * $t / 100
        $b = 255 + (128 - 255) * $t / 100
        $GCol += (Get-Fg ([int]$r) ([int]$g) ([int]$b))
    }
    $filas = $BannerFilas.Count
    [Console]::Write("$script:Esc[?25l")
    [Console]::Write("`n`n")
    # Decode: el arte se revela de izquierda a derecha con gradiente cian-rosa.
    for ($frame = 0; $frame -le 16; $frame++) {
        $reveal = [int][math]::Floor($cols * $frame / 16)
        for ($fila = 0; $fila -lt $filas; $fila++) {
            $linea = $BannerFilas[$fila]
            $out = ""
            for ($j = 0; $j -lt $linea.Length; $j++) {
                $ch = $linea[$j]
                if ($j -lt $reveal) { $out += "$($GCol[$j])$ch$script:Reset" }
                elseif ($ch -eq " ") { $out += " " }
                else { $out += "$cGl$($shades[$rnd.Next(3)])$script:Reset" }
            }
            [Console]::Write("$script:Esc[K$Margen$out`n")
        }
        if ($frame -lt 16) { [Console]::Write("$script:Esc[$($filas)A"); Start-Sleep -Milliseconds 60 }
    }
    [Console]::Write("$script:Esc[?25h")
    # Dos lineas de aire entre el logo y el texto que lo sigue.
    [Console]::Write("`n`n")
    foreach ($l in $IdentFilas) {
        [Console]::Write("$Margen$(Get-Fg 130 130 160)$l$script:Reset`n")
    }
    [Console]::Write("`n$Margen$(Get-Fg 255 255 255)$Sub$script:Reset`n`n")
}

# --- Tarjeta final ----------------------------------------------------------
function Mostrar-Card {
    param([string]$Titulo, [string]$Linea, [string]$Firma = "")
    if (-not $script:VtOk) {
        Write-Host ""
        Write-Host "== $Titulo =="
        Write-Host $Linea
        if ($Firma) { Write-Host $Firma }
        Write-Host ""
        return
    }
    $borde = "─" * 44
    $t1 = " {0,-43}" -f $Titulo
    $t2 = " {0,-43}" -f $Linea
    [Console]::Write("$Margen$(Get-Fg 130 80 255)╭$borde╮$script:Reset`n")
    [Console]::Write("$Margen$(Get-Fg 130 80 255)│$script:Reset$script:Bold$(Get-Fg 240 240 250)$t1$script:Reset$(Get-Fg 130 80 255)│$script:Reset`n")
    [Console]::Write("$Margen$(Get-Fg 130 80 255)│$script:Reset$(Get-Fg 160 160 190)$t2$script:Reset$(Get-Fg 130 80 255)│$script:Reset`n")
    if ($Firma) {
        $t3 = " {0,-43}" -f $Firma
        [Console]::Write("$Margen$(Get-Fg 130 80 255)│$script:Reset$(Get-Fg 160 160 190)$t3$script:Reset$(Get-Fg 130 80 255)│$script:Reset`n")
    }
    [Console]::Write("$Margen$(Get-Fg 130 80 255)╰$borde╯$script:Reset`n")
    # Pulso del check final, como en las demos.
    foreach ($paso in @((Get-Fg 40 90 60),(Get-Fg 0 180 110),(Get-Fg 0 255 140))) {
        [Console]::Write("`r$Margen$paso$script:Bold✓$script:Reset")
        Start-Sleep -Milliseconds 90
    }
    [Console]::Write(" Listo.`n`n")
}

# --- Pausa interactiva ------------------------------------------------------
# Read-Host espera; cualquier entrada (incluso vacia) sigue.
function Pause-Step {
    if ($script:VtOk) {
        [Console]::Write("$Margen$script:Bold$(Get-Fg 220 220 240)Cuando termines, volve a esta ventana y presiona ENTER para seguir:$script:Reset ")
        [void](Read-Host)
    } else {
        [void](Read-Host "Cuando termines, volve a esta ventana y presiona ENTER para seguir")
    }
}

# Ante una falla: reintentar solo con confirmacion explicita. ENTER, EOF o
# cualquier otra entrada = salir. Evita bucles con stdin cerrado.
function Ask-Reintentar {
    $r = Read-Host 'Reintentar o salir? (escriba "reintentar" para volver a probar; ENTER = salir)'
    if ($null -eq $r) { return $false }
    if ($r.Trim().ToLower() -eq "reintentar") { return $true }
    return $false
}

# Corre un paso; ante fallo ofrece UN reintento explicito (vuelve a chequear).
function Invoke-PasoConReintento {
    param([string]$Nombre, [scriptblock]$Accion)
    if (& $Accion) { return $true }
    if (Ask-Reintentar) {
        Write-Host ""
        Seccion "Reintento - $Nombre"
        return (& $Accion)
    }
    return $false
}

# --- Intro del asistente ----------------------------------------------------
function Show-Intro {
    Mostrar-Banner "Se esta conectando al nodo G1 ($NodeName) - conexion desde Windows"
    if ($script:VtOk) { Barra-Progreso "preparando la consola" 1; [Console]::Write("`n") }
    Marg "Este asistente te lleva paso a paso."
    Marg "En cada paso: lees la explicacion, haces la accion, y volves a esta"
    Marg "ventana a presionar ENTER para seguir. Nada se hace sin avisarte."
    Write-Host ""
    Pause-Step
}

# --- Localizar Tailscale ----------------------------------------------------
function Find-Tailscale {
    $cmd = Get-Command tailscale.exe -ErrorAction SilentlyContinue
    if ($cmd) { $script:Ts = $cmd.Source; return $true }
    if ($env:ProgramFiles) {
        $p = Join-Path $env:ProgramFiles "Tailscale\tailscale.exe"
        if (Test-Path -LiteralPath $p) { $script:Ts = $p; return $true }
    }
    return $false
}

# --- Instalar si falta (multi-metodo, con barra y fallback) -----------------
function Install-Tailscale {
    # Metodo 1: winget.
    $rc = Invoke-ConBarra "instalando con winget (1/2)" "winget.exe" @("install","--id","Tailscale.Tailscale","-e","--accept-source-agreements","--accept-package-agreements")
    if ($rc -ne 0) {
        # El alias de winget puede no dejarse lanzar por Start-Process: inline.
        try {
            & winget install --id Tailscale.Tailscale -e --accept-source-agreements --accept-package-agreements 2>$null | Out-Null
            $rc = $LASTEXITCODE
        } catch { $rc = 1 }
    }
    Start-Pausa 5
    [void](Find-Tailscale)
    if ($script:Ts) {
        Estado "OK" "Tailscale instalado con winget."
        Start-Pausa 1
        return $true
    }
    Estado "X" "No se pudo instalar con winget."
    Start-Pausa 1

    # Metodo 2: instalador oficial (MSI) con elevacion (UAC).
    $msi = Join-Path $env:TEMP "tailscale-setup-ing3.msi"
    $url = Get-TailscaleMsiUrl
    [void](Invoke-ConBarra "descargando el instalador oficial (2/2)" "curl.exe" @("-L","-fsS","-o",$msi,$url))
    if ((Test-Path -LiteralPath $msi) -and ((Get-Item -LiteralPath $msi).Length -gt 0)) {
        Estado ".." "Instalando en silencio; Windows puede pedir confirmacion (UAC)."
        $p = $null
        try {
            $p = Start-Process -FilePath "msiexec.exe" -ArgumentList @("/i", $msi, "/qn", "/norestart") -Verb RunAs -PassThru -ErrorAction Stop
        } catch { $p = $null }
        if ($p) {
            Watch-ConBarra "instalando con el instalador oficial" $p
        } else {
            try { Start-Process -FilePath "msiexec.exe" -ArgumentList @("/i", $msi, "/qn", "/norestart") -Verb RunAs -Wait -ErrorAction Stop } catch { }
        }
        Start-Pausa 5
        [void](Find-Tailscale)
        if ($script:Ts) {
            Estado "OK" "Tailscale instalado con el instalador oficial."
            Start-Pausa 1
            return $true
        }
        Estado "X" "No se pudo instalar con el instalador oficial."
        Start-Pausa 1
    } else {
        Estado "X" "No se pudo instalar con el instalador oficial."
        Start-Pausa 1
    }

    Estado "X" "No se pudo instalar Tailscale automaticamente."
    Marg "    Instalelo a mano desde la pagina oficial:"
    MargUrl "      " "https://tailscale.com/download/windows"
    Start-Pausa 1
    return $false
}

# --- Estado estructurado de Tailscale (ConvertFrom-Json nativo) -------------
function Get-StructuredStatus {
    # $true si 'status --json' devolvio JSON parseable.
    $script:StParseOk = $false
    $script:StObj = $null
    if (-not $script:Ts) { return $false }
    $out = $null
    $rc = 1
    try {
        $out = & $script:Ts status --json 2>$null
        $rc = $LASTEXITCODE
    } catch { return $false }
    if ($rc -ne 0) { return $false }
    $text = ($out | Out-String).Trim()
    if (-not $text) { return $false }
    try { $obj = $text | ConvertFrom-Json -ErrorAction Stop } catch { return $false }
    if ($null -eq $obj) { return $false }
    $script:StObj = $obj
    $script:StParseOk = $true
    return $true
}

function Test-ValidIp([string]$ip) {
    if (-not $ip) { return $false }
    return ($ip -match '^[0-9]{1,3}(\.[0-9]{1,3}){3}$')
}

# Prueba de sesion: Running + Self.Online=true + IP valida + sin AuthURL.
function Test-AuthState {
    if (-not (Get-StructuredStatus)) { return $false }
    $o = $script:StObj
    if ("$($o.BackendState)" -ne "Running") { return $false }
    if ($o.Self.Online -ne $true) { return $false }
    if ("$($o.AuthURL)") { return $false }
    foreach ($ip in @($o.Self.TailscaleIPs)) {
        if (Test-ValidIp "$ip") { return $true }
    }
    return $false
}

# --- IP de tailnet actual (nula si no hay sesion verificada) ----------------
function Get-TailscaleIp {
    if (-not (Get-StructuredStatus)) { return $null }
    foreach ($ip in @($script:StObj.Self.TailscaleIPs)) {
        if (Test-ValidIp "$ip") { return "$ip" }
    }
    return $null
}

# Diagnostico del estado observado.
function Show-StateHelp {
    $o = $script:StObj
    if (-not $o) {
        Marg "      - Tailscale no reporta estado (servicio apagado o sin iniciar)."
        return
    }
    switch ("$($o.BackendState)") {
        "NeedsLogin"       { Marg "      - Falta iniciar sesion (Tailscale pide autenticacion)." }
        "NeedsMachineAuth" { Marg "      - Falta aprobacion administrativa de este equipo." }
        "Stopped"          { Marg "      - El servicio de Tailscale esta detenido." }
        "Starting"         { Marg "      - Tailscale esta arrancando; espere unos segundos." }
        "InUseOtherUser"   { Marg "      - Tailscale esta en uso por otro usuario de Windows." }
        "Running"          { Marg "      - El servicio corre, pero la sesion no esta verificada (Online=false o sin IP)." }
        ""                 { Marg "      - Tailscale no reporta estado (servicio apagado o sin iniciar)." }
        default            { Marg "      - Estado reportado: $($o.BackendState)." }
    }
    if ("$($o.AuthURL)") { Marg "      - Hay un enlace de autenticacion pendiente (esta arriba)." }
}

# --- Iniciar sesion: captura el enlace, lo imprime y abre el navegador ------
function Clear-LoginProc {
    if ($script:LoginProc) {
        try { if (-not $script:LoginProc.HasExited) { $script:LoginProc.Kill() } } catch { }
    }
    $script:LoginProc = $null
}

function Start-Login {
    param([int]$WaitSeconds = 30)
    $tmpO = [System.IO.Path]::GetTempFileName()
    $tmpE = [System.IO.Path]::GetTempFileName()
    $proc = $null
    $url = ""
    try {
        try {
            $proc = Start-Process -FilePath $script:Ts -ArgumentList @("up") -NoNewWindow -RedirectStandardOutput $tmpO -RedirectStandardError $tmpE -PassThru -ErrorAction Stop
            $script:LoginProc = $proc
        } catch { $proc = $null }
        $deadline = (Get-Date).AddSeconds($WaitSeconds)
        while ((Get-Date) -lt $deadline) {
            $txt = ""
            try { if (Test-Path $tmpO) { $txt += (Get-Content $tmpO -Raw -ErrorAction SilentlyContinue) } } catch { }
            try { if (Test-Path $tmpE) { $txt += (Get-Content $tmpE -Raw -ErrorAction SilentlyContinue) } } catch { }
            if ($txt) {
                $m = [regex]::Match($txt, 'https://login\.tailscale\.com/[A-Za-z0-9/._-]+')
                if ($m.Success) { $url = $m.Value; break }
            }
            if ((Get-StructuredStatus) -and ("$($script:StObj.AuthURL)")) { $url = "$($script:StObj.AuthURL)"; break }
            if ($proc) { try { if ($proc.HasExited) { break } } catch { } }
            Start-Sleep -Milliseconds 500
        }
    } finally {
        Remove-Item $tmpO, $tmpE -ErrorAction SilentlyContinue
    }
    if ($url) {
        Marg "     Enlace de inicio de sesion (copiable):"
        MargUrl "     " "$url"
        try {
            Start-Process $url -ErrorAction Stop
            Marg "     Se intento abrir el navegador; si no aparece, copielo a mano."
        } catch {
            Marg "     Copielo y peguelo a mano en el navegador."
        }
        return $true
    }
    Marg "     Tailscale no mostro un enlace de inicio de sesion todavia."
    return $false
}

# --- Espera acotada de la sesion (max intentos x 5 s), sonda estructurada ---
function Wait-Session {
    param([int]$Max)
    for ($i = 1; $i -le $Max; $i++) {
        if ($script:VtOk) { Barra-Progreso "esperando la sesion ($i/$Max)" 5 }
        else { Start-Pausa 5 }
        if (Test-AuthState) { return $true }
        if (-not $script:VtOk) { Write-Host "     Esperando la sesion de Tailscale... ($i/$Max)" }
    }
    return $false
}

# --- El nodo es visible? (IP exacta en Peer, Online=true; nunca texto) ------
function Test-NodeVisible {
    if (-not (Get-StructuredStatus)) { return $false }
    $o = $script:StObj
    if (-not $o.Peer) { return $false }
    foreach ($p in $o.Peer.PSObject.Properties) {
        $peer = $p.Value
        if ($peer.Online -ne $true) { continue }
        foreach ($ip in @($peer.TailscaleIPs)) {
            if ("$ip" -eq $NodeIp) { return $true }
        }
    }
    return $false
}

# --- Mostrar el enlace del grupo: siempre visible y copiable a mano --------
function Mostrar-Enlace {
    Write-Host ""
    Marg "Enlace de invitacion del grupo (queda en pantalla para copiarlo):"
    Write-Host ""
    MargUrl "    " "$ShareLink"
    Write-Host ""
    Marg "Vigencia del enlace: 30 dias desde el 27/10/2026 hasta el 26/11/2026; si vence, pida uno nuevo."
    Marg "Se puede usar mas de una vez. Si la ventana del navegador se cerro,"
    Marg "abrio con otra cuenta o no se abrio, copie este enlace y peguelo en"
    Marg "el navegador donde inicio sesion en el paso 2."
    Write-Host ""
}

# --- Copiar el enlace al portapapeles (mejor esfuerzo; nunca fatal) ---------
function Copiar-Enlace {
    try {
        Set-Clipboard -Value $ShareLink -ErrorAction Stop
        return $true
    } catch {
        return $false
    }
}

# --- Ensenar el enlace: siempre visible, copiado y abierto (mejor esfuerzo) -
function Presentar-Enlace {
    Mostrar-Enlace
    if (Copiar-Enlace) {
        Marg "Tambien quedo copiado al portapapeles: puede pegarlo con CTRL+V."
        Write-Host ""
    }
    try {
        Start-Process $ShareLink -ErrorAction Stop
        return $true
    } catch {
        Estado "!" "No se pudo abrir el navegador automaticamente."
        Marg "    Copie el enlace de arriba y abralo a mano en el navegador."
        return $false
    }
}

# --- Pausa del paso 3: ENTER sigue; "e" vuelve a mostrar y abrir el enlace --
function Pause-StepEnlace {
    while ($true) {
        if ($script:VtOk) {
            [Console]::Write("$Margen$script:Bold$(Get-Fg 220 220 240)Cuando termines, volve a esta ventana y presiona ENTER (o escribi e + ENTER para ver y reabrir el enlace):$script:Reset ")
            $r = Read-Host
        } else {
            $r = Read-Host 'Cuando termines, volve a esta ventana y presiona ENTER (o escribi e + ENTER para ver y reabrir el enlace)'
        }
        if ($null -eq $r) { Write-Host ""; return }
        $t = "$r".Trim().ToLower()
        if ($t -eq "e" -or $t -eq "enlace") {
            [void](Presentar-Enlace)
            continue
        }
        Write-Host ""
        return
    }
}

# --- Espera acotada de la invitacion (max intentos x 15 s) ------------------
function Wait-Node {
    param([int]$Max)
    for ($i = 1; $i -le $Max; $i++) {
        if ($script:VtOk) { Barra-Progreso "esperando el nodo ($i/$Max)" 15 }
        else { Start-Pausa 15 }
        if (Test-NodeVisible) { return $true }
        if (-not $script:VtOk) { Write-Host "     Esperando que el nodo aparezca... ($i/$Max)" }
    }
    return $false
}

# --- Paso 1/5: Tailscale instalado ------------------------------------------
function Invoke-Step1Tailscale {
    Write-Host ""
    Seccion "Paso 1/5 - Tailscale instalado"
    Marg "Tailscale es la red privada por la que se llega al nodo."
    Marg "Tiene que estar instalado y SIEMPRE encendido mientras trabajas."
    Write-Host ""

    if (Find-Tailscale) {
        $script:OkTs = $true
        Estado "OK" "Tailscale encontrado: $script:Ts"
        return $true
    }

    Estado ".." "Tailscale no esta instalado. Se intentara instalarlo..."
    [void](Install-Tailscale)
    if (Find-Tailscale) {
        $script:OkTs = $true
        Estado "OK" "Tailscale encontrado: $script:Ts"
        Start-Pausa 1
        return $true
    }

    # Ruta manual: solo si la instalacion automatica no dejo Tailscale.
    Write-Host ""
    Pause-Step
    [void](Find-Tailscale)

    if ($script:Ts) {
        $script:OkTs = $true
        Estado "OK" "Tailscale encontrado: $script:Ts"
        return $true
    }

    Estado "X" "Tailscale sigue sin estar instalado."
    Marg "    Instalelo a mano y vuelva a correr el script:"
    MargUrl "      " "https://tailscale.com/download/windows"
    return $false
}

# --- Paso 2/5: Iniciar sesion -----------------------------------------------
function Invoke-Step2Session {
    Write-Host ""
    Seccion "Paso 2/5 - Iniciar sesion"
    Marg "Hay que iniciar sesion con TU cuenta de Tailscale."
    Marg "Si no tenia cuenta, se crea gratis en el navegador."
    Marg "IMPORTANTE: deje ABIERTA la ventana del navegador que se abre."
    Marg "La va a necesitar en el paso 3 (aceptar la invitacion)."
    Write-Host ""

    if (Test-AuthState) {
        $script:OkSession = $true
        Estado "OK" "Sesion de Tailscale activa - su IP: $(Get-TailscaleIp)"
        Clear-LoginProc
        return $true
    }

    Estado ".." "Tailscale no esta conectado. Vamos a iniciar sesion."
    if (Get-StructuredStatus) { Show-StateHelp }
    [void](Start-Login -WaitSeconds 30)
    Write-Host ""
    Pause-Step
    [void](Wait-Session 12)
    if (Test-AuthState) {
        $script:OkSession = $true
        Estado "OK" "Sesion de Tailscale activa - su IP: $(Get-TailscaleIp)"
        Clear-LoginProc
        return $true
    }

    Estado ".." "Todavia no hay sesion verificada. Causas probables:"
    if (Get-StructuredStatus) { Show-StateHelp }
    Marg "    Vamos con un reintento acotado."
    [void](Start-Login -WaitSeconds 30)
    Write-Host ""
    Pause-Step
    [void](Wait-Session 12)
    if (Test-AuthState) {
        $script:OkSession = $true
        Estado "OK" "Sesion de Tailscale activa - su IP: $(Get-TailscaleIp)"
        Clear-LoginProc
        return $true
    }

    Estado "X" "Tailscale sigue sin sesion verificada."
    Marg "    Termine de iniciar sesion (o apruebe el equipo) y vuelva a correr el script."
    Clear-LoginProc
    return $false
}

# --- Paso 3/5: Aceptar la invitacion ----------------------------------------
function Invoke-Step3Invite {
    Write-Host ""
    Seccion "Paso 3/5 - Aceptar la invitacion"

    if (-not (Test-AuthState)) {
        Estado "X" "Se perdio o no se completo la sesion de Tailscale."
        Marg "    Vuelva al paso 2 e inicie sesion antes de aceptar la invitacion."
        return $false
    }

    if (Test-NodeVisible) {
        $script:OkNode = $true
        Estado "OK" "El nodo $NodeName ($NodeIp) ya es visible en su tailnet."
        return $true
    }

    Marg "Se va a mostrar el enlace de invitacion del grupo y se va a intentar"
    Marg "abrir el navegador con el. El enlace queda en pantalla: si la ventana"
    Marg "se cierra o abre con otra cuenta, se puede volver a copiar."
    Marg "RECOMENDACION: use el MISMO navegador donde inicio sesion en"
    Marg "Tailscale (la ventana que dejo abierta en el paso 2) y no la"
    Marg "cierre, asi la invitacion abre con su sesion ya iniciada."
    Marg "El enlace es del grupo: tratalo como una contrasena."
    [void](Presentar-Enlace)
    Pause-StepEnlace
    [void](Wait-Node 12)
    if (-not (Test-AuthState)) {
        Estado "X" "Se perdio la sesion de Tailscale durante la espera."
        Marg "    Vuelva a iniciar sesion y a correr el script."
        return $false
    }
    if (Test-NodeVisible) {
        $script:OkNode = $true
        Estado "OK" "El nodo $NodeName ($NodeIp) ya es visible."
        return $true
    }

    Estado ".." "El nodo no aparece todavia. Causas probables:"
    Marg "      - No acepto la invitacion."
    Marg "      - La acepto con OTRA cuenta de Tailscale."
    Marg "      - Tailscale esta apagado o sin sesion."
    Marg "    Vamos con un reintento. El enlace se muestra de nuevo."
    [void](Presentar-Enlace)
    Pause-StepEnlace
    [void](Wait-Node 12)
    if (-not (Test-AuthState)) {
        Estado "X" "Se perdio la sesion de Tailscale durante la espera."
        Marg "    Vuelva a iniciar sesion y a correr el script."
        return $false
    }
    if (Test-NodeVisible) {
        $script:OkNode = $true
        Estado "OK" "El nodo $NodeName ($NodeIp) ya es visible."
        return $true
    }

    Estado "X" "El nodo $NodeName sigue sin aparecer."
    Marg "    Si la invitacion no se acepto, o se acepto con OTRA cuenta, abra"
    Marg "    de nuevo este enlace con la cuenta del paso 2:"
    Mostrar-Enlace
    Marg "    Despues vuelva a correr el script."
    return $false
}

# --- Paso 4/5: Prueba de conexion (automatica) ------------------------------
function Invoke-Step4Tests {
    Write-Host ""
    Seccion "Paso 4/5 - Prueba de conexion"
    Marg "Estas pruebas son automaticas: no hay que hacer nada."
    Invoke-PingChecks
    Invoke-WebCheck
}

# Ping ICMP con indicador visible, en el mismo hilo (sin jobs ni runspaces, para
# no invocar cmdlets reales fuera de los overrides del driver de pruebas).
# Mantiene la semantica y el conteo de Test-Connection; limpia la linea en finally.
function Invoke-IcmpAnimado {
    param([string]$Ip)
    $animar = $script:VtOk
    $frames = @("⠋","⠙","⠹","⠸")
    $visto = $false
    try {
        for ($k = 0; $k -lt 3; $k++) {
            if ($animar) {
                [Console]::Write("`r$Margen$(Get-Fg 0 200 255)$($frames[$k % $frames.Count])$script:Reset comprobando el alcance (ping ICMP)...$script:Esc[K")
            }
            if (Test-Connection -Count 1 -ComputerName $Ip -Quiet -ErrorAction Stop) { $visto = $true }
        }
    } finally {
        if ($animar) { [Console]::Write("`r$script:Esc[K") }
    }
    return $visto
}

# --- Ping (Tailscale e ICMP, best-effort) -----------------------------------
function Invoke-PingChecks {
    Write-Host ""
    Marg "Prueba de alcance (ping, maximo 5 segundos):"

    $text = Invoke-Spin "probando el nodo (tailscale ping)..." $script:Ts @("ping","--c","1","--timeout","5s",$NodeIp)
    if ($text -and ($text -match 'pong')) {
        foreach ($l in (($text -split "\r?\n") | Where-Object { $_ -ne "" } | Select-Object -First 3)) {
            if ($l -match 'pong') {
                if ($script:VtOk) { [Console]::Write("$Margen  $(Get-Fg 0 230 130)✓$script:Reset  $l`n") }
                else { Write-Host $l }
            } else {
                Marg "     $l"
            }
        }
        if ($text -match 'DERP') {
            $ruta = "relay (DERP)"
        } elseif ($text -match 'via') {
            $ruta = "directo"
        } else {
            $ruta = ""
        }
        if ($ruta) { Estado "OK" "Tailscale ping: responde ($ruta)." }
        else { Estado "OK" "Tailscale ping: responde." }
        $script:OkPing = $true
    } else {
        if ($text -and $text.Trim()) {
            (($text -split "\r?\n") | Where-Object { $_ -ne "" } | Select-Object -First 2) | ForEach-Object { Marg "     $_" }
        }
        Estado "X" "✗ Tailscale ping: sin respuesta en 5 segundos (parece fuera de la red)."
    }

    try {
        if (Invoke-IcmpAnimado -Ip $NodeIp) {
            Estado "OK" "Ping ICMP: responde."
        } else {
            Estado "!" "Ping ICMP sin respuesta (Tailscale puede bloquearlo; no es fatal)."
        }
    } catch {
        Estado "i" "No se pudo ejecutar el ping ICMP (se omite)."
    }
}

# --- Comprobar que el servicio web del nodo responde ------------------------
function Invoke-WebCheck {
    Write-Host ""
    Marg "Servicio web del nodo:"
    $code = "000"
    try {
        # curl.exe viene incluido en Windows 10 y 11. -k: el certificado del nodo
        # es propio. Solo una respuesta de servicio valida (2xx, o 3xx de
        # redireccion) cuenta como exito: 000, vacio, texto, 1xx, 4xx y 5xx son
        # fallos y el resumen no debe quedar "conexion lista".
        $raw = Invoke-Spin "consultando el servicio web..." "curl.exe" @("-k","-s","-o","NUL","-w","%{http_code}","--max-time","15",$WebUrl) -OcultarErr
        if ($null -eq $raw) { $raw = "000" }
        $code = "$raw".Trim()
    } catch { $code = "000" }
    $okCode = $false
    if ($code -match '^[0-9]{3}$') {
        $n = [int]$code
        if ($n -ge 200 -and $n -le 399) { $okCode = $true }
    }
    if ($okCode) {
        $script:OkWeb = $true
        Estado "OK" "Servicio web del nodo: responde (HTTP $code)."
    } else {
        Estado "X" "Servicio web del nodo: no respondio (codigo $code)."
        Marg "    Causa probable: Tailscale apagado o el nodo no esta arriba."
    }
}

# --- Paso 5/5: Resumen final ------------------------------------------------
# Imprime la seccion y las cinco lineas de estado (sin la tarjeta).
function Write-ResumenEstados {
    SeccionFin "Paso 5/5 - Resumen"
    if ($script:OkTs) { Estado "OK" "Tailscale instalado" } else { Estado "X" " Tailscale no instalado" }
    if ($script:OkSession) { Estado "OK" "Sesion de Tailscale activa y verificada" } else { Estado "X" " Tailscale sin sesion verificada" }
    if ($script:OkNode) { Estado "OK" "Nodo $NodeName visible" } else { Estado "X" " Nodo $NodeName no visible" }
    if ($script:OkPing) { Estado "OK" "Tailscale ping: responde" } else { Estado "X" " Tailscale ping: sin respuesta" }
    if ($script:OkWeb) { Estado "OK" "Servicio web del nodo alcanzable" } else { Estado "X" " Servicio web del nodo no alcanzable" }
}

# Pantalla final interactiva: restaura cursor/colores, avisa del historial y
# muestra el resumen compacto a continuacion, con espaciado normal (sin
# limpieza, padding ni reposicion del cursor). Resumen y tarjeta una sola vez.
function Show-PantallaFinal {
    if ($script:VtOk) { [Console]::Write("$script:Reset$script:Esc[?25h") }
    Marg "El detalle de los pasos queda en el historial; ahora se muestra el resumen final."
    Write-Host ""
    Write-ResumenEstados
    AvisoOk "Conexion lista."
    Marg "Tailscale tiene que estar siempre encendido mientras trabaja."
    Marg "La administracion del nodo esta en los materiales de infraestructura:"
    Marg "'Nodo G1 - Mapa y guia de infraestructura.docx' y 'Tu nodo G1.md'."
    Mostrar-Card "CONEXION LISTA" "Nodo G1 - todo en orden" "Lencina Nicolas"
}

function Show-Summary {
    param([switch]$Final)
    $todo = ($script:OkTs -and $script:OkSession -and $script:OkNode -and $script:OkPing -and $script:OkWeb)
    if ($Final -and $script:UiEstilo -and $todo) {
        Show-PantallaFinal
        return $true
    }
    Write-Host ""
    Write-ResumenEstados
    Write-Host ""
    if ($todo) {
        AvisoOk "Conexion lista."
        Marg "Tailscale tiene que estar siempre encendido mientras trabaja."
        Marg "La administracion del nodo esta en los materiales de infraestructura:"
        Marg "'Nodo G1 - Mapa y guia de infraestructura.docx' y 'Tu nodo G1.md'."
        Mostrar-Card "CONEXION LISTA" "Nodo G1 - todo en orden" "Lencina Nicolas"
        return $true
    }
    Marg "La conexion no quedo completa. Revise los [X] de arriba y la seccion"
    Marg "'Problemas comunes' de su guia, y vuelva a correr el script."
    return $false
}

# ============================================================================
# Ejecucion
# ============================================================================
try {
    Show-Intro
    if (-not (Invoke-PasoConReintento "Paso 1 - Tailscale instalado" { Invoke-Step1Tailscale })) { Clear-LoginProc; [void](Show-Summary); exit 1 }
    if (-not (Invoke-PasoConReintento "Paso 2 - Iniciar sesion" { Invoke-Step2Session }))   { Clear-LoginProc; [void](Show-Summary); exit 1 }
    if (-not (Invoke-PasoConReintento "Paso 3 - Aceptar la invitacion" { Invoke-Step3Invite })) { Clear-LoginProc; [void](Show-Summary); exit 1 }
    Invoke-Step4Tests
    Clear-LoginProc

    if (Show-Summary -Final) { exit 0 }
    exit 1
} finally {
    # Cursor visible al salir, sin ocultar el codigo de salida real.
    if ($script:VtOk) { [Console]::Write("$script:Reset$script:Esc[?25h") }
}
