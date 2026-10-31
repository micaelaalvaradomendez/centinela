#!/usr/bin/env bash
# ============================================================================
# conectar-unix.sh - Asistente paso a paso para conectar al nodo G1
#                    (pve-is3-g1) desde Linux o macOS.
#
# Es un asistente interactivo, sin menu: en cada paso la persona lee la
# explicacion, hace la accion y vuelve a la ventana a presionar ENTER.
#   1) Tailscale instalado (lo instala si falta).
#   2) Iniciar sesion en Tailscale.
#   3) Aceptar la invitacion del nodo.
#   4) Prueba de conexion (ping y servicio web). Automatica, sin pausa.
#   5) Resumen con [OK]/[X].
#
# Uso (desde la carpeta de conexion, la que contiene la Guia de conexion (.docx)):
#   bash scripts/conectar-unix.sh
# ============================================================================
set -uo pipefail

# --- Datos del nodo (NO cambiar sin coordinar con infraestructura) ----------
NODE_NAME="pve-is3-g1"
NODE_IP="100.116.49.111"
# URL del servicio web: SOLO se usa para comprobar que responde. No se imprime,
# no se muestra como destino y no lleva usuario/realm/contrasena: la
# administracion vive en los materiales de infraestructura del grupo.
WEB_URL="https://${NODE_IP}:8006/"
SHARE_LINK="https://login.tailscale.com/admin/invite/7hgHA4tQXPYjS22vFJx421"

# --- Estado de cada chequeo (1 = paso OK) -----------------------------------
OK_TS=0
OK_SESSION=0
OK_NODE=0
OK_PING=0
OK_WEB=0

# Binario de Tailscale (se completa en find_tailscale).
TS=""
OS="$(uname -s)"

# Prefijo de privilegios (sudo solo si hace falta) y salida del ultimo run_barra.
SUDO=""
if [ "$(id -u)" != "0" ] && command -v sudo >/dev/null 2>&1; then SUDO="sudo"; fi
RUN_LOG=""
LOGIN_PID=""

# --- Estado estructurado de Tailscale (lo completa estado_estructurado) -----
JSON_PARSER=""
ST_STATE=""
ST_ONLINE="false"
ST_SELFIPS=""
ST_AUTHURL=""
ST_PEERIPS=""

# ============================================================================
# Estilo de consola (banner al iniciar, tarjeta al cerrar).
# Se anima solo si hay una terminal interactiva; si no, queda el encabezado
# simple de siempre. Decode y colores de las demos (24 bits); bash 3.2+ (Ubuntu/macOS).
# ============================================================================
UI_ESTILO=0
if [ -t 1 ] && [ "${TERM:-}" != "dumb" ]; then UI_ESTILO=1; fi
# Si la terminal es angosta (< 68 columnas), modo simple.
if [ "$UI_ESTILO" = 1 ]; then
  COLS_TERM="$(tput cols 2>/dev/null || echo 80)"
  case "$COLS_TERM" in ''|*[!0-9]*) COLS_TERM=80 ;; esac
  if [ "$COLS_TERM" -lt 68 ]; then UI_ESTILO=0; fi
fi
# Consola limpia al arrancar (solo con estilo de UI completo: interactiva y
# ancha). En modo simple o no interactivo no se emite ninguna secuencia.
# Limpieza visible propia SIN `clear`: el comando `clear` puede emitir ESC[3J
# (borra el scrollback) en terminales con capacidad E3. ESC[H + ESC[2J solo
# limpia la pantalla visible.
ESC_SEQ=$'\033'
if [ "$UI_ESTILO" = 1 ]; then printf '%s[H%s[2J' "$ESC_SEQ" "$ESC_SEQ"; fi
RESET=""; BOLD=""; DIM=""
if [ "$UI_ESTILO" = 1 ]; then
  RESET="${ESC_SEQ}[0m"; BOLD="${ESC_SEQ}[1m"; DIM="${ESC_SEQ}[2m"
fi

fg() { if [ "$UI_ESTILO" = 1 ]; then printf '%s[38;2;%s;%s;%sm' "$ESC_SEQ" "$1" "$2" "$3"; fi; }

# Margen de lectura del contenido (nada queda pegado al borde).
marg() { if [ "$UI_ESTILO" = 1 ]; then printf '%s%s\n' "$MARGEN" "$1"; else printf '%s\n' "$1"; fi; }

# Enlaces: en azul para distinguirlos del texto comun. $1 = prefijo literal,
# $2 = URL (azul), $3 = sufijo literal opcional. Sin estilo queda texto plano.
marg_url() {
  if [ "$UI_ESTILO" = 1 ]; then
    printf '%s%s%s%s%s%s\n' "$MARGEN" "$1" "$(fg 90 170 255)" "$2" "$RESET" "${3:-}"
  else
    printf '%s%s%s\n' "$1" "$2" "${3:-}"
  fi
}

# Estados de linea: [OK] verde, [X] rojo, [!] ambar, [..] gris, [i] celeste.
estado() {  # $1 = etiqueta ; $2 = texto
  if [ "$UI_ESTILO" != 1 ]; then
    printf '[%s] %s\n' "$1" "$2"
    return 0
  fi
  case "$1" in
    OK) printf '%s%s[OK]%s %s\n' "$MARGEN" "$(fg 0 230 130)" "$RESET" "$2" ;;
    X) printf '%s%s[X]%s %s\n' "$MARGEN" "$(fg 255 90 90)" "$RESET" "$2" ;;
    '!') printf '%s%s[!]%s %s\n' "$MARGEN" "$(fg 255 196 0)" "$RESET" "$2" ;;
    '..') printf '%s%s[..]%s %s\n' "$MARGEN" "$(fg 150 150 160)" "$RESET" "$2" ;;
    i) printf '%s%s[i]%s %s\n' "$MARGEN" "$(fg 0 200 255)" "$RESET" "$2" ;;
    *) printf '%s[%s] %s\n' "$MARGEN" "$1" "$2" ;;
  esac
}

# Cabecera de paso: separadores en gris y titulo destacado.
seccion() {  # $1 = titulo
  if [ "$UI_ESTILO" = 1 ]; then
    printf '%s%s------------------------------------------%s\n' "$MARGEN" "$(fg 95 95 105)" "$RESET"
    printf '%s%s%s %s%s\n' "$MARGEN" "$BOLD" "$(fg 220 220 240)" "$1" "$RESET"
    printf '%s%s------------------------------------------%s\n' "$MARGEN" "$(fg 95 95 105)" "$RESET"
  else
    printf -- '------------------------------------------\n'
    printf ' %s\n' "$1"
    printf -- '------------------------------------------\n'
  fi
}

# Cabecera del resumen final (borde con =).
seccion_fin() {  # $1 = titulo
  if [ "$UI_ESTILO" = 1 ]; then
    printf '%s%s==========================================%s\n' "$MARGEN" "$(fg 95 95 105)" "$RESET"
    printf '%s%s%s %s%s\n' "$MARGEN" "$BOLD" "$(fg 220 220 240)" "$1" "$RESET"
    printf '%s%s==========================================%s\n' "$MARGEN" "$(fg 95 95 105)" "$RESET"
  else
    printf -- '==========================================\n'
    printf ' %s\n' "$1"
    printf -- '==========================================\n'
  fi
}

# Linea de exito suelta (sin etiqueta), en verde.
ok_line() {  # $1 = texto
  if [ "$UI_ESTILO" = 1 ]; then
    printf '%s%s%s%s\n' "$MARGEN" "$(fg 0 230 130)" "$1" "$RESET"
  else
    printf '%s\n' "$1"
  fi
}

# Gradiente de la barra (se arma una sola vez por corrida).
BARRA_LISTA=0
BARRA_GRAD=()
preparar_barra() {
  local j tt r g b
  if [ "$BARRA_LISTA" = 1 ]; then return 0; fi
  for ((j=0; j<24; j++)); do
    tt=$(( 100 * j / 24 ))
    r=$(( 0 + (255 - 0) * tt / 100 ))
    g=$(( 229 + (0 - 229) * tt / 100 ))
    b=$(( 255 + (128 - 255) * tt / 100 ))
    BARRA_GRAD[$j]="$(fg "$r" "$g" "$b")"
  done
  BARRA_LISTA=1
}

# Barra de progreso con brillo (estilo demos). $1 = mensaje ; $2 = segundos.
barra_progreso() {
  local msg="$1" segs="$2" width=24 pasos i k pct filled shine barra
  local cShine cHead cEmpty cBase
  if [ "$UI_ESTILO" != 1 ]; then return 0; fi
  preparar_barra
  cShine="$(fg 255 236 236)"
  cHead="$(fg 255 150 152)"
  cEmpty="${DIM}$(fg 70 70 90)"
  cBase="$(fg 160 160 190)"
  pasos=$(( segs * 12 ))
  for ((i=0; i<=pasos; i++)); do
    pct=$(( 100 * i / pasos ))
    filled=$(( width * i / pasos ))
    shine=$(( (i * (width + 8) / pasos) - 8 ))
    barra=""
    for ((k=0; k<width; k++)); do
      if [ "$k" -lt "$filled" ]; then
        if [ "$k" -ge "$shine" ] && [ "$k" -lt "$((shine + 5))" ]; then
          barra+="${cShine}▓${RESET}"
        elif [ "$k" -eq "$((filled - 1))" ] && [ $((i % 2)) -eq 0 ]; then
          barra+="${cHead}▓${RESET}"
        else
          barra+="${BARRA_GRAD[$k]}█${RESET}"
        fi
      else
        barra+="${cEmpty}░${RESET}"
      fi
    done
    printf '\r%s%s%s [%s]%s %3d%%\033[K' "$MARGEN" "$cBase" "$msg" "$barra" "$RESET" "$pct"
    sleep 0.08
  done
  printf '\n'
}

# Corre un comando mostrando spinner; imprime su salida (stderr oculto).
run_spin() {  # $1 = mensaje ; resto = comando
  local msg="$1"; shift
  local tmp pid rc i=0
  local FRAMES
  FRAMES=("⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏")
  tmp="$(mktemp 2>/dev/null)" || tmp="/tmp/spin.$$"
  if [ "$UI_ESTILO" = 1 ]; then
    "$@" >"$tmp" 2>/dev/null &
    pid=$!
    while kill -0 "$pid" 2>/dev/null || [ "$i" -lt 11 ]; do
      printf '\r%s%s%s%s %s\033[K' "$MARGEN" "$(fg 0 200 255)" "${FRAMES[$((i % 10))]}" "$RESET" "$msg" >&2
      sleep 0.09
      i=$((i + 1))
    done
    wait "$pid"; rc=$?
    printf '\r\033[K' >&2
  else
    "$@" >"$tmp" 2>/dev/null
    rc=$?
  fi
  cat "$tmp"
  rm -f "$tmp"
  return $rc
}

# Corre un comando con spinner sin emitir su salida: el stdout que run_spin
# imprime al final se descarta; solo el spinner (stderr) queda visible y se
# conserva el codigo de salida del comando.
run_spin_silencioso() {  # $1 = mensaje ; resto = comando
  local msg="$1"; shift
  run_spin "$msg" "$@" >/dev/null
}

# Restaura colores y cursor (defensivo en salida final y por senales).
restaurar_terminal() {
  if [ "$UI_ESTILO" = 1 ]; then printf '%s[0m%s[?25h' "$ESC_SEQ" "$ESC_SEQ"; fi
}

# Dibuja una linea de barra (sin newline) con el estilo de barra_progreso.
# $1 = mensaje ; $2 = porcentaje (0-100) ; $3 = frame (para el brillo).
dibuja_barra() {
  local msg="$1" pct="$2" f="$3" width=24 filled shine barra k
  local cShine cHead cEmpty cBase
  cShine="$(fg 255 236 236)"
  cHead="$(fg 255 150 152)"
  cEmpty="${DIM}$(fg 70 70 90)"
  cBase="$(fg 160 160 190)"
  filled=$(( width * pct / 100 ))
  if [ "$filled" -gt "$width" ]; then filled=$width; fi
  shine=$(( (f * 3) % (width + 8) - 8 ))
  barra=""
  for ((k=0; k<width; k++)); do
    if [ "$k" -lt "$filled" ]; then
      if [ "$k" -ge "$shine" ] && [ "$k" -lt "$((shine + 5))" ]; then
        barra+="${cShine}▓${RESET}"
      elif [ "$k" -eq "$((filled - 1))" ] && [ $((f % 2)) -eq 0 ]; then
        barra+="${cHead}▓${RESET}"
      else
        barra+="${BARRA_GRAD[$k]}█${RESET}"
      fi
    else
      barra+="${cEmpty}░${RESET}"
    fi
  done
  printf '\r%s%s%s [%s]%s %3d%%\033[K' "$MARGEN" "$cBase" "$msg" "$barra" "$RESET" "$pct"
}

# Corre un comando mostrando una barra de progreso (salida+errores a RUN_LOG).
# Devuelve el rc del comando. Minimo visible ~1s.
run_barra() {  # $1 = mensaje ; resto = comando
  local msg="$1"; shift
  local tmp pid rc i pct
  RUN_LOG=""
  if [ "$UI_ESTILO" != 1 ]; then
    printf '[..] %s\n' "$msg"
    tmp="$(mktemp 2>/dev/null)" || tmp="/tmp/barra.$$"
    RUN_LOG="$tmp"
    "$@" >"$tmp" 2>&1
    rc=$?
    return $rc
  fi
  preparar_barra
  tmp="$(mktemp 2>/dev/null)" || tmp="/tmp/barra.$$"
  RUN_LOG="$tmp"
  "$@" >"$tmp" 2>&1 &
  pid=$!
  i=0
  while kill -0 "$pid" 2>/dev/null || [ "$i" -lt 11 ]; do
    pct=$(( 95 * i / (i + 12) ))
    dibuja_barra "$msg" "$pct" "$i"
    sleep 0.09
    i=$((i + 1))
  done
  wait "$pid"; rc=$?
  # Completar a 100% (3 ticks) y cerrar la linea.
  dibuja_barra "$msg" 97 "$i"; sleep 0.09; i=$((i + 1))
  dibuja_barra "$msg" 99 "$i"; sleep 0.09; i=$((i + 1))
  dibuja_barra "$msg" 100 "$i"
  printf '\n'
  return $rc
}

# Imprime hasta 2 de las ultimas lineas no vacias del ultimo run_barra.
detalle_log() {
  [ -n "$RUN_LOG" ] || return 0
  [ -s "$RUN_LOG" ] || return 0
  grep -v '^[[:space:]]*$' "$RUN_LOG" 2>/dev/null | tail -2 | while IFS= read -r l; do
    printf '    Detalle: %.120s\n' "$l"
  done
}

# Banner (letras de bloque) del nodo (el gradiente se arma en mostrar_banner).
BANNER_FILAS=(
" ████ █████ ██  ██ ██████ ████ ██  ██ █████ ██     ████ "
"██    ██    ███ ██   ██    ██  ███ ██ ██    ██    ██  ██"
"██    ████  ██ ███   ██    ██  ██ ███ ████  ██    ██████"
"██    ██    ██  ██   ██    ██  ██  ██ ██    ██    ██  ██"
" ████ █████ ██  ██   ██   ████ ██  ██ █████ █████ ██  ██"
)

# Margen izquierdo del bloque del logo (4 columnas) y lineas de identidad.
MARGEN='    '
IDENT_FILAS=(
"Ingreso a la red de Tailscale"
"Carrera Ingenieria de Software III"
"Universidad Nacional de Tierra del Fuego, Antartida e Islas del Atlantico Sur"
"Hecho por Lencina Nicolas"
)

mostrar_banner() {  # $1 = linea operativa
  local sub="$1" fila i frame j ch out reveal linea cGl cols filas j2 tt r g b
  local SHADES
  if [ "$UI_ESTILO" != 1 ]; then
    echo "=========================================="
    echo " ${sub}"
    echo "=========================================="
    for i in 0 1 2 3; do echo " ${IDENT_FILAS[$i]}"; done
    echo
    return 0
  fi
  cGl="$(fg 90 90 110)"
  SHADES=("░" "▒" "▓")
  filas=5
  cols=${#BANNER_FILAS[0]}
  local GCOL
  GCOL=()
  for ((j2=0; j2<cols; j2++)); do
    tt=$(( 100 * j2 / cols ))
    r=$(( (255 - 0) * tt / 100 ))
    g=$(( 229 + (0 - 229) * tt / 100 ))
    b=$(( 255 + (128 - 255) * tt / 100 ))
    GCOL[$j2]="$(fg "$r" "$g" "$b")"
  done
  printf '%s[?25l' "$ESC_SEQ"
  echo
  echo
  # Decode: el arte se revela de izquierda a derecha con gradiente cian-rosa.
  for ((frame=0; frame<=16; frame++)); do
    reveal=$(( cols * frame / 16 ))
    for ((fila=0; fila<filas; fila++)); do
      linea="${BANNER_FILAS[$fila]}"
      out=""
      for ((j=0; j<${#linea}; j++)); do
        ch="${linea:j:1}"
        if [ "$j" -lt "$reveal" ]; then
          out+="${GCOL[$j]}${ch}${RESET}"
        elif [ "$ch" = " " ]; then
          out+=" "
        else
          out+="${DIM}${cGl}${SHADES[$((RANDOM % 3))]}${RESET}"
        fi
      done
      printf '\033[K%s%s\n' "$MARGEN" "$out"
    done
    if [ "$frame" -lt 16 ]; then
      printf '%s[%dA' "$ESC_SEQ" "$filas"
      sleep 0.06
    fi
  done
  printf '%s[?25h' "$ESC_SEQ"
  # Dos lineas de aire entre el logo y el texto que lo sigue.
  echo
  echo
  for i in 0 1 2 3; do
    printf '%s%s%s%s\n' "$MARGEN" "$(fg 130 130 160)" "${IDENT_FILAS[$i]}" "$RESET"
  done
  # Linea operativa actual, en dim.
  printf '\n%s%s%s%s\n' "$MARGEN" "$(fg 255 255 255)" "$sub" "$RESET"
  echo
}

tarjeta_final() {  # $1 = titulo ; $2 = linea ; $3 = firma (opcional)
  local t1="$1" t2="$2" t3="${3:-}" borde="" i=0
  if [ "$UI_ESTILO" != 1 ]; then
    echo
    echo "== $1 =="
    echo "$2"
    if [ -n "$t3" ]; then echo "$t3"; fi
    echo
    return 0
  fi
  while [ "$i" -lt 44 ]; do borde="${borde}─"; i=$((i + 1)); done
  printf -v t1 ' %-43s' "$1"
  printf -v t2 ' %-43s' "$2"
  printf '%s%s╭%s╮%s\n' "$MARGEN" "$(fg 130 80 255)" "$borde" "$RESET"
  printf '%s%s│%s%s%s%s%s│%s\n' "$MARGEN" "$(fg 130 80 255)" "$RESET" "${BOLD}$(fg 240 240 250)" "$t1" "$RESET" "$(fg 130 80 255)" "$RESET"
  printf '%s%s│%s%s%s%s%s│%s\n' "$MARGEN" "$(fg 130 80 255)" "$RESET" "$(fg 160 160 190)" "$t2" "$RESET" "$(fg 130 80 255)" "$RESET"
  if [ -n "$t3" ]; then
    printf -v t3 ' %-43s' "$t3"
    printf '%s%s│%s%s%s%s%s│%s\n' "$MARGEN" "$(fg 130 80 255)" "$RESET" "$(fg 160 160 190)" "$t3" "$RESET" "$(fg 130 80 255)" "$RESET"
  fi
  printf '%s%s╰%s╯%s\n' "$MARGEN" "$(fg 130 80 255)" "$borde" "$RESET"
  # Pulso del check final, como en las demos.
  local paso
  for paso in "$(fg 40 90 60)" "$(fg 0 180 110)" "$(fg 0 255 140)"; do
    printf '\r%s%s%s✓%s' "$MARGEN" "$paso" "$BOLD" "$RESET"
    sleep 0.09
  done
  printf ' Listo.\n'
  echo
}

# ============================================================================
# Utilidades del asistente
# ============================================================================

# --- Pausa interactiva ------------------------------------------------------
# Imprime el aviso y lee una linea. Cualquier entrada (incluso vacia) sigue.
pause_step() {
  if [ "$UI_ESTILO" = 1 ]; then
    printf '%s%s%s%s%s ' "$MARGEN" "$BOLD" "$(fg 220 220 240)" "Cuando termines, volve a esta ventana y presiona ENTER para seguir:" "$RESET"
  else
    printf '%s' "Cuando termines, volve a esta ventana y presiona ENTER para seguir: "
  fi
  IFS= read -r _ || true
  echo
}

# Ante una falla: reintentar solo con confirmacion explicita. ENTER, EOF o
# cualquier otra entrada = salir. Evita bucles cuando no hay TTY.
# Devuelve 0 = reintentar, 1 = salir.
reintentar_o_salir() {
  local r
  printf '%sReintentar o salir? (escriba "reintentar" para volver a probar; ENTER = salir): ' "$MARGEN"
  IFS= read -r r || return 1
  r="$(printf '%s' "$r" | tr '[:upper:]' '[:lower:]' | tr -d '[:space:]')"
  [ "$r" = "reintentar" ] && return 0
  return 1
}

# Corre un paso; ante fallo ofrece UN reintento explicito (vuelve a chequear).
paso_con_reintento() {  # $1 = nombre ; $2 = funcion
  local nombre="$1" fn="$2"
  if "$fn"; then return 0; fi
  if reintentar_o_salir; then
    echo
    seccion "Reintento - ${nombre}"
    "$fn"
    return $?
  fi
  return 1
}

# --- Intro del asistente ----------------------------------------------------
show_intro() {
  mostrar_banner "Se esta conectando al nodo G1 (${NODE_NAME}) - conexion desde ${OS}"
  if [ "$UI_ESTILO" = 1 ]; then barra_progreso "preparando la consola" 1; echo; fi
  marg "Este asistente te lleva paso a paso."
  marg "En cada paso: lees la explicacion, haces la accion, y volves a esta"
  marg "ventana a presionar ENTER para seguir. Nada se hace sin avisarte."
  echo
  pause_step
}

# --- Localizar el binario de tailscale --------------------------------------
find_tailscale() {
  if command -v tailscale >/dev/null 2>&1; then TS="tailscale"; return 0; fi
  # macOS: app oficial (cask / App Store).
  if [ -x "/Applications/Tailscale.app/Contents/MacOS/Tailscale" ]; then
    TS="/Applications/Tailscale.app/Contents/MacOS/Tailscale"; return 0
  fi
  return 1
}

# --- Metodos de instalacion (cada uno es una funcion) -----------------------
metodo_mac_brew() {
  brew install --cask tailscale-app
}

metodo_mac_pkg() {
  curl -fsSL -o /tmp/Tailscale-ing3.pkg \
    https://pkgs.tailscale.com/stable/Tailscale-latest-macos.pkg &&
  $SUDO installer -pkg /tmp/Tailscale-ing3.pkg -target /
}

metodo_linux_script() {
  curl -fsSL https://tailscale.com/install.sh | sh
}

metodo_linux_apt() {
  . /etc/os-release
  curl -fsSL "https://pkgs.tailscale.com/stable/$ID/$VERSION_CODENAME.noarmor.gpg" |
    $SUDO tee /usr/share/keyrings/tailscale-archive-keyring.gpg >/dev/null &&
  curl -fsSL "https://pkgs.tailscale.com/stable/$ID/$VERSION_CODENAME.tailscale-keyring.list" |
    $SUDO tee /etc/apt/sources.list.d/tailscale.list >/dev/null &&
  $SUDO apt-get update &&
  $SUDO apt-get install -y tailscale tailscale-archive-keyring
}

# Aplica el metodo 2 de Linux? (apt-get + ID ubuntu/debian + codename).
hay_metodo_linux_apt() {
  command -v apt-get >/dev/null 2>&1 || return 1
  [ -r /etc/os-release ] || return 1
  local ID="" VERSION_CODENAME=""
  . /etc/os-release 2>/dev/null || return 1
  case "${ID:-}" in ubuntu|debian) ;; *) return 1 ;; esac
  [ -n "${VERSION_CODENAME:-}" ] || return 1
  return 0
}

# --- Instalar Tailscale si falta (metodos visibles con barra) ---------------
install_tailscale() {
  local total=0 n=0 rc
  local nombres=()
  local etiquetas=()
  local funciones=()

  # 1) Pre-autenticacion de sudo: los sudo internos del comando en segundo
  #    plano no deben pedir la contrasena detras de la barra.
  if [ "$SUDO" = "sudo" ]; then
    estado '..' "La instalacion necesita permisos de administrador: sudo va a pedir la contrasena (una sola vez)."
    if sudo -v; then
      estado OK "Permisos de administrador listos."
      sleep 1
    else
      estado X "No se pudieron obtener los permisos de administrador."
      sleep 1
      return 1
    fi
  fi

  # 2) Armar la lista de metodos disponibles (1 o 2) con sus etiquetas exactas.
  if [ "$OS" = "Darwin" ]; then
    if command -v brew >/dev/null 2>&1; then
      nombres[0]="Homebrew"; etiquetas[0]="instalando con Homebrew (1/2)"; funciones[0]="metodo_mac_brew"
      nombres[1]="el paquete oficial"; etiquetas[1]="instalando con el paquete oficial (2/2)"; funciones[1]="metodo_mac_pkg"
    else
      nombres[0]="el paquete oficial"; etiquetas[0]="instalando con el paquete oficial (1/1)"; funciones[0]="metodo_mac_pkg"
    fi
  else
    nombres[0]="el script oficial"; etiquetas[0]="instalando con el script oficial (1/2)"; funciones[0]="metodo_linux_script"
    if hay_metodo_linux_apt; then
      nombres[1]="el repositorio oficial"; etiquetas[1]="instalando con el repositorio oficial (2/2)"; funciones[1]="metodo_linux_apt"
    fi
  fi
  total=${#nombres[@]}

  # 3) Probar los metodos en orden: un intento visible por metodo.
  for ((n=0; n<total; n++)); do
    run_barra "${etiquetas[$n]}" "${funciones[$n]}"
    rc=$?
    if [ "$rc" = 0 ]; then
      estado OK "Tailscale instalado con ${nombres[$n]}."
      sleep 1
      return 0
    fi
    estado X "No se pudo instalar con ${nombres[$n]}."
    detalle_log
    sleep 1
  done

  # 4) Ningun metodo funciono: fallback a la pagina oficial.
  estado X "No se pudo instalar Tailscale automaticamente."
  marg "    Instalelo a mano desde la pagina oficial:"
  if [ "$OS" = "Darwin" ]; then
    marg_url "      " "https://tailscale.com/download/mac"
  else
    marg_url "      " "https://tailscale.com/download/linux"
    marg_url "      (alternativa: curl -fsSL " "https://tailscale.com/install.sh" " | sh)"
  fi
  sleep 1
  return 1
}

# --- Lector JSON (jq o python3/python) --------------------------------------
# El estado de Tailscale se prueba con 'status --json' y un parser real: nunca
# por coincidencia de texto. Prerrequisito: jq o python3/python instalado.
detectar_parser_json() {
  if command -v jq >/dev/null 2>&1; then JSON_PARSER="jq"; return 0; fi
  if command -v python3 >/dev/null 2>&1; then JSON_PARSER="python3"; return 0; fi
  if command -v python >/dev/null 2>&1; then JSON_PARSER="python"; return 0; fi
  JSON_PARSER=""
  return 1
}

# Normaliza el JSON de 'status --json' a lineas clave=valor (parser real).
# STATE, ONLINE, SELFIPS, AUTHURL y una linea PEER por cada peer.
_json_norm() {
  case "$JSON_PARSER" in
    jq)
      jq -r '
        "STATE=" + ((.BackendState // "") | tostring),
        "ONLINE=" + (((.Self.Online // false)) | tostring),
        "SELFIPS=" + (((.Self.TailscaleIPs // [])) | join(",")),
        "AUTHURL=" + ((.AuthURL // "") | tostring),
        (((.Peer // {}) | to_entries[]) |
          "PEER=" + (((.value.Online // false)) | tostring) + " " +
          (((.value.TailscaleIPs // [])) | join(",")))
      ' 2>/dev/null
      ;;
    python3|python)
      "$JSON_PARSER" -c '
import sys, json
try:
    d = json.load(sys.stdin)
except Exception:
    sys.exit(3)
if not isinstance(d, dict):
    sys.exit(3)
selfd = d.get("Self") or {}

def flag(v):
    return "true" if v is True else "false"

ips = [x for x in (selfd.get("TailscaleIPs") or []) if isinstance(x, str)]
print("STATE=" + str(d.get("BackendState") or ""))
print("ONLINE=" + flag(selfd.get("Online")))
print("SELFIPS=" + ",".join(ips))
print("AUTHURL=" + str(d.get("AuthURL") or ""))
peers = d.get("Peer") or {}
if isinstance(peers, dict):
    for p in peers.values():
        if not isinstance(p, dict):
            continue
        pips = [x for x in (p.get("TailscaleIPs") or []) if isinstance(x, str)]
        print("PEER=" + flag(p.get("Online")) + " " + ",".join(pips))
' 2>/dev/null
      ;;
  esac
}

# IPv4 con forma valida (Tailscale asigna 100.x.y.z).
ip_valida() {
  printf '%s' "$1" | grep -Eq '^[0-9]{1,3}(\.[0-9]{1,3}){3}$'
}

# La lista separada por comas contiene el valor exacto $1?
lista_contiene() {  # $1 = valor ; $2 = lista
  case ",$2," in *",$1,"*) return 0 ;; esac
  return 1
}

# Hay al menos una IP valida en la lista $1 (separada por comas)?
hay_ip_valida() {
  printf '%s' "$1" | grep -Eq '(^|,)[0-9]{1,3}(\.[0-9]{1,3}){3}(,|$)'
}

# Corre 'status --json' y deja el estado en las variables ST_*.
# rc: 0 = parseado OK ; 2 = CLI o JSON invalido ; 3 = sin parser.
estado_estructurado() {
  local json rc norm line
  detectar_parser_json || return 3
  json="$("$TS" status --json 2>/dev/null)"; rc=$?
  [ "$rc" -eq 0 ] || return 2
  [ -n "$json" ] || return 2
  norm="$(printf '%s' "$json" | _json_norm)" || return 3
  [ -n "$norm" ] || return 3
  ST_STATE=""
  ST_ONLINE="false"
  ST_SELFIPS=""
  ST_AUTHURL=""
  ST_PEERIPS=""
  while IFS= read -r line; do
    case "$line" in
      STATE=*) ST_STATE="${line#STATE=}" ;;
      ONLINE=*) ST_ONLINE="${line#ONLINE=}" ;;
      SELFIPS=*) ST_SELFIPS="${line#SELFIPS=}" ;;
      AUTHURL=*) ST_AUTHURL="${line#AUTHURL=}" ;;
      PEER=*) ST_PEERIPS="${ST_PEERIPS}${line#PEER=}"$'\n' ;;
    esac
  done <<EOF
$norm
EOF
  return 0
}

# Prueba de sesion: Running + Online=true + IP valida + sin AuthURL pendiente.
auth_ok() {
  estado_estructurado || return 1
  [ "$ST_STATE" = "Running" ] || return 1
  [ "$ST_ONLINE" = "true" ] || return 1
  [ -z "$ST_AUTHURL" ] || return 1
  hay_ip_valida "$ST_SELFIPS" || return 1
  return 0
}

# Primera IP valida del estado (vacia si no hay).
primera_ip() {
  local ip
  while IFS= read -r ip; do
    if ip_valida "$ip"; then printf '%s' "$ip"; return 0; fi
  done <<EOF
$(printf '%s' "$ST_SELFIPS" | tr ',' '\n')
EOF
  return 1
}

# Explica el estado observado (diagnostico accionable).
explicar_estado() {
  case "$ST_STATE" in
    NeedsLogin) marg "      - Falta iniciar sesion (Tailscale pide autenticacion)." ;;
    NeedsMachineAuth) marg "      - Falta aprobacion administrativa de este equipo." ;;
    Stopped) marg "      - El servicio de Tailscale esta detenido." ;;
    Starting) marg "      - Tailscale esta arrancando; espere unos segundos." ;;
    InUseOtherUser) marg "      - Tailscale esta en uso por otro usuario de Windows." ;;
    Running) marg "      - El servicio corre, pero la sesion no esta verificada (Online=false o sin IP)." ;;
    NoState|"") marg "      - Tailscale no reporta estado (servicio apagado o sin iniciar)." ;;
    *) marg "      - Estado reportado: ${ST_STATE}." ;;
  esac
  if [ -n "$ST_AUTHURL" ]; then
    marg "      - Hay un enlace de autenticacion pendiente (esta arriba)."
  fi
}

# Abre una URL en el navegador; rc 0 si se pudo lanzar.
abrir_url() {  # $1 = url
  if command -v xdg-open >/dev/null 2>&1; then xdg-open "$1" >/dev/null 2>&1 & return 0; fi
  if command -v open >/dev/null 2>&1; then open "$1" >/dev/null 2>&1 & return 0; fi
  return 1
}

# Limpia un proceso de login en segundo plano si quedo vivo.
limpiar_login_bg() {
  if [ -n "$LOGIN_PID" ] && kill -0 "$LOGIN_PID" 2>/dev/null; then
    kill "$LOGIN_PID" 2>/dev/null || true
    wait "$LOGIN_PID" 2>/dev/null || true
  fi
  LOGIN_PID=""
}

# Lanza 'tailscale up' en segundo plano (proceso propio), captura el enlace de
# inicio de sesion de stdout+stderr y del AuthURL de status, lo imprime copiable
# y abre el navegador best-effort. No espera a que el proceso termine.
iniciar_login_con_enlace() {
  local log pid url i
  limpiar_login_bg
  if [ "$SUDO" = "sudo" ]; then
    marg "     La sesion se inicia como administrador: sudo puede pedir la contrasena."
    sudo -v || true
  fi
  log="$(mktemp 2>/dev/null)" || log="/tmp/login.$$"
  : >"$log"
  $SUDO "$TS" up >"$log" 2>&1 </dev/null &
  pid=$!
  LOGIN_PID=$pid
  url=""
  i=0
  while [ "$i" -lt 60 ]; do
    url="$(grep -Eo 'https://login\.tailscale\.com/[A-Za-z0-9/._-]+' "$log" 2>/dev/null | head -1)"
    if [ -z "$url" ] && estado_estructurado && [ -n "$ST_AUTHURL" ]; then
      url="$ST_AUTHURL"
    fi
    [ -n "$url" ] && break
    kill -0 "$pid" 2>/dev/null || break
    sleep 0.5
    i=$((i + 1))
  done
  if [ -n "$url" ]; then
    marg "     Enlace de inicio de sesion (copiable):"
    marg_url "     " "$url"
    if abrir_url "$url"; then
      marg "     Se intento abrir el navegador; si no aparece, copielo a mano."
    else
      marg "     No hay comando de navegador: copielo y peguelo a mano."
    fi
  else
    marg "     Tailscale no mostro un enlace de inicio de sesion todavia."
    grep -v '^[[:space:]]*$' "$log" 2>/dev/null | tail -4 | while IFS= read -r l; do marg "     $l"; done
  fi
  rm -f "$log" 2>/dev/null
}

# --- Iniciar sesion (captura e imprime el enlace) ---------------------------
start_login() {
  if [ "$OS" = "Darwin" ]; then
    marg "     En macOS tambien podes iniciar sesion desde la app Tailscale (boton 'Log in')."
  fi
  # Se inicia en segundo plano y se muestra/abre el enlace de autenticacion.
  iniciar_login_con_enlace
}

# --- Espera acotada de la sesion (max intentos x 5 s), sonda estructurada ---
wait_session() {
  local max="$1" i
  for i in $(seq 1 "$max"); do
    if [ "$UI_ESTILO" = 1 ]; then
      barra_progreso "esperando la sesion (${i}/${max})" 5
    else
      sleep 5
    fi
    if auth_ok; then return 0; fi
    if [ "$UI_ESTILO" != 1 ]; then
      echo "     Esperando la sesion de Tailscale... (${i}/${max})"
    fi
  done
  return 1
}

# --- El nodo es visible? (IP exacta en Peer, Online=true; nunca texto) ------
node_visible() {
  estado_estructurado || return 1
  local line on ips
  while IFS= read -r line; do
    [ -n "$line" ] || continue
    on="${line%% *}"
    ips="${line#* }"
    if [ "$on" = "true" ] && lista_contiene "$NODE_IP" "$ips"; then return 0; fi
  done <<EOF
$ST_PEERIPS
EOF
  return 1
}

# --- Mostrar el enlace del grupo: siempre visible y copiable a mano --------
mostrar_enlace() {
  echo
  marg "Enlace de invitacion del grupo (queda en pantalla para copiarlo):"
  echo
  marg_url "    " "$SHARE_LINK"
  echo
  marg "Vigencia del enlace: 30 dias desde el 27/10/2026 hasta el 26/11/2026; si vence, pida uno nuevo."
  marg "Se puede usar mas de una vez. Si la ventana del navegador se cerro,"
  marg "abrio con otra cuenta o no se abrio, copie este enlace y peguelo en"
  marg "el navegador donde inicio sesion en el paso 2."
  echo
}

# --- Copiar el enlace al portapapeles (mejor esfuerzo; nunca fatal) ---------
copiar_enlace() {
  if [ "$OS" = "Darwin" ] && command -v pbcopy >/dev/null 2>&1; then
    printf '%s' "$SHARE_LINK" | pbcopy >/dev/null 2>&1 && return 0
  fi
  if [ -n "${WAYLAND_DISPLAY:-}" ] && command -v wl-copy >/dev/null 2>&1; then
    printf '%s' "$SHARE_LINK" | wl-copy >/dev/null 2>&1 && return 0
  fi
  if [ -n "${DISPLAY:-}" ] && command -v xclip >/dev/null 2>&1; then
    printf '%s' "$SHARE_LINK" | xclip -selection clipboard >/dev/null 2>&1 && return 0
  fi
  if [ -n "${DISPLAY:-}" ] && command -v xsel >/dev/null 2>&1; then
    printf '%s' "$SHARE_LINK" | xsel --clipboard --input >/dev/null 2>&1 && return 0
  fi
  if command -v clip.exe >/dev/null 2>&1; then
    printf '%s' "$SHARE_LINK" | clip.exe >/dev/null 2>&1 && return 0
  fi
  return 1
}

# --- Abrir la invitacion en el navegador (1 si no hay comando) --------------
open_browser() {
  if command -v xdg-open >/dev/null 2>&1; then
    xdg-open "$SHARE_LINK" >/dev/null 2>&1 &
    return 0
  fi
  if command -v open >/dev/null 2>&1; then
    open "$SHARE_LINK" >/dev/null 2>&1 &
    return 0
  fi
  return 1
}

# --- Ensenar el enlace: siempre visible, copiado y abierto (mejor esfuerzo) -
presentar_enlace() {
  mostrar_enlace
  if copiar_enlace; then
    marg "Tambien quedo copiado al portapapeles: puede pegarlo con CTRL+V."
    echo
  fi
  if ! open_browser; then
    estado '!' "No se encontro un comando para abrir el navegador."
    marg "    Copie el enlace de arriba y abralo a mano en el navegador."
    return 1
  fi
  return 0
}

# --- Pausa del paso 3: ENTER sigue; "e" vuelve a mostrar y abrir el enlace --
pausa_enlace() {
  local r
  while true; do
    if [ "$UI_ESTILO" = 1 ]; then
      printf '%s%s%s%s%s ' "$MARGEN" "$BOLD" "$(fg 220 220 240)" \
        "Cuando termines, volve a esta ventana y presiona ENTER (o escribi e + ENTER para ver y reabrir el enlace):" "$RESET"
    else
      printf '%s' "Cuando termines, volve a esta ventana y presiona ENTER (o escribi e + ENTER para ver y reabrir el enlace): "
    fi
    if ! IFS= read -r r; then
      echo
      return 0
    fi
    case "$(printf '%s' "$r" | tr '[:upper:]' '[:lower:]' | tr -d '[:space:]')" in
      e|enlace)
        presentar_enlace
        ;;
      *)
        echo
        return 0
        ;;
    esac
  done
}

# --- Espera acotada de la invitacion (max intentos x 15 s) ------------------
wait_node() {
  local max="$1" i
  for i in $(seq 1 "$max"); do
    if [ "$UI_ESTILO" = 1 ]; then
      barra_progreso "esperando el nodo (${i}/${max})" 15
    else
      sleep 15
    fi
    if node_visible; then return 0; fi
    if [ "$UI_ESTILO" != 1 ]; then
      echo "     Esperando que el nodo aparezca... (${i}/${max})"
    fi
  done
  return 1
}

# ============================================================================
# Pasos del asistente
# ============================================================================

# --- Paso 1/5: Tailscale instalado ------------------------------------------
step1_tailscale() {
  echo
  seccion "Paso 1/5 - Tailscale instalado"
  marg "Tailscale es la red privada por la que se llega al nodo."
  marg "Tiene que estar instalado y SIEMPRE encendido mientras trabajas."
  echo

  if find_tailscale; then
    OK_TS=1
    estado OK "Tailscale encontrado: ${TS}"
    return 0
  fi

  estado '..' "Tailscale no esta instalado. Se intentara instalarlo..."
  if install_tailscale; then
    if find_tailscale; then
      OK_TS=1
      estado OK "Tailscale encontrado: ${TS}"
      sleep 1
      return 0
    fi
  fi

  # Ruta manual: solo si la instalacion automatica no dejo Tailscale.
  echo
  pause_step

  if find_tailscale; then
    OK_TS=1
    estado OK "Tailscale encontrado: ${TS}"
    return 0
  fi

  estado X "Tailscale sigue sin estar instalado."
  marg "    Instalelo a mano y vuelva a correr el script:"
  marg_url "      Linux:   curl -fsSL " "https://tailscale.com/install.sh" " | sh"
  marg_url "      macOS:   " "https://tailscale.com/download/mac"
  return 1
}

# --- Paso 2/5: Iniciar sesion -----------------------------------------------
step2_session() {
  local ip
  echo
  seccion "Paso 2/5 - Iniciar sesion"
  marg "Hay que iniciar sesion con TU cuenta de Tailscale."
  marg "Si no tenia cuenta, se crea gratis en el navegador."
  marg "IMPORTANTE: deje ABIERTA la ventana del navegador que se abre."
  marg "La va a necesitar en el paso 3 (aceptar la invitacion)."
  echo

  if ! detectar_parser_json; then
    estado X "No se encontro un lector JSON (jq o python3) para verificar el estado."
    marg "    Instale jq o python3 y vuelva a correr el script."
    return 1
  fi

  if auth_ok; then
    OK_SESSION=1
    ip="$(primera_ip || true)"
    estado OK "Sesion de Tailscale activa - su IP: $ip"
    limpiar_login_bg
    return 0
  fi

  estado '..' "Tailscale no esta conectado. Vamos a iniciar sesion."
  if estado_estructurado; then explicar_estado; fi
  start_login
  echo
  pause_step
  wait_session 12
  if auth_ok; then
    OK_SESSION=1
    ip="$(primera_ip || true)"
    estado OK "Sesion de Tailscale activa - su IP: $ip"
    limpiar_login_bg
    return 0
  fi

  estado '..' "Todavia no hay sesion verificada. Causas probables:"
  if estado_estructurado; then explicar_estado; fi
  marg "    Vamos con un reintento acotado."
  start_login
  echo
  pause_step
  wait_session 12
  if auth_ok; then
    OK_SESSION=1
    ip="$(primera_ip || true)"
    estado OK "Sesion de Tailscale activa - su IP: $ip"
    limpiar_login_bg
    return 0
  fi

  estado X "Tailscale sigue sin sesion verificada."
  marg "    Termine de iniciar sesion (o apruebe el equipo) y vuelva a correr el script."
  limpiar_login_bg
  return 1
}

# --- Paso 3/5: Aceptar la invitacion ----------------------------------------
step3_invite() {
  echo
  seccion "Paso 3/5 - Aceptar la invitacion"

  if ! auth_ok; then
    estado X "Se perdio o no se completo la sesion de Tailscale."
    marg "    Vuelva al paso 2 e inicie sesion antes de aceptar la invitacion."
    return 1
  fi

  if node_visible; then
    OK_NODE=1
    estado OK "El nodo ${NODE_NAME} (${NODE_IP}) ya es visible en su tailnet."
    return 0
  fi

  marg "Se va a mostrar el enlace de invitacion del grupo y se va a intentar"
  marg "abrir el navegador con el. El enlace queda en pantalla: si la ventana"
  marg "se cierra o abre con otra cuenta, se puede volver a copiar."
  marg "RECOMENDACION: use el MISMO navegador donde inicio sesion en"
  marg "Tailscale (la ventana que dejo abierta en el paso 2) y no la"
  marg "cierre, asi la invitacion abre con su sesion ya iniciada."
  marg "El enlace es del grupo: tratalo como una contrasena."
  presentar_enlace
  pausa_enlace
  wait_node 12
  if ! auth_ok; then
    estado X "Se perdio la sesion de Tailscale durante la espera."
    marg "    Vuelva a iniciar sesion y a correr el script."
    return 1
  fi
  if node_visible; then
    OK_NODE=1
    estado OK "El nodo ${NODE_NAME} (${NODE_IP}) ya es visible."
    return 0
  fi

  estado '..' "El nodo no aparece todavia. Causas probables:"
  marg "      - No acepto la invitacion."
  marg "      - La acepto con OTRA cuenta de Tailscale."
  marg "      - Tailscale esta apagado o sin sesion."
  marg "    Vamos con un reintento. El enlace se muestra de nuevo."
  presentar_enlace
  pausa_enlace
  wait_node 12
  if ! auth_ok; then
    estado X "Se perdio la sesion de Tailscale durante la espera."
    marg "    Vuelva a iniciar sesion y a correr el script."
    return 1
  fi
  if node_visible; then
    OK_NODE=1
    estado OK "El nodo ${NODE_NAME} (${NODE_IP}) ya es visible."
    return 0
  fi

  estado X "El nodo ${NODE_NAME} sigue sin aparecer."
  marg "    Si la invitacion no se acepto, o se acepto con OTRA cuenta, abra"
  marg "    de nuevo este enlace con la cuenta del paso 2:"
  mostrar_enlace
  marg "    Despues vuelva a correr el script."
  return 1
}

# --- Ping (Tailscale e ICMP, best-effort) -----------------------------------
# Ping corto: maximo 5 segundos en total. Imprime la salida; rc 0 si hubo pong.
run_ts_ping() {  # $1 = args extra ("--c 3" o vacio)
  local args="$1" out=""
  if command -v timeout >/dev/null 2>&1; then
    out="$(timeout 5 "$TS" ping $args "$NODE_IP" 2>&1)"
  elif command -v gtimeout >/dev/null 2>&1; then
    out="$(gtimeout 5 "$TS" ping $args "$NODE_IP" 2>&1)"
  elif [ "$OS" = "Darwin" ] && command -v perl >/dev/null 2>&1; then
    out="$(perl -e 'alarm shift; exec @ARGV' 5 "$TS" ping $args "$NODE_IP" 2>&1)"
  else
    # Ultimo recurso: un solo intento, acotado por el propio CLI (<=5 s).
    out="$("$TS" ping --c 1 --timeout 5s "$NODE_IP" 2>&1)"
  fi
  printf '%s' "$out"
  case "$out" in *pong*) return 0 ;; *) return 1 ;; esac
}

ping_node() {
  local out ruta icmp_ok=0
  echo
  marg "Prueba de alcance (ping, maximo 5 segundos):"

  out="$(run_spin "probando el nodo (tailscale ping)..." run_ts_ping "--c 3")"
  case "$out" in
    *"not defined"*|*"unknown flag"*)
      # Version vieja sin --c: se reintenta sin flags.
      out="$(run_spin "probando el nodo (tailscale ping)..." run_ts_ping "")"
      ;;
  esac
  if [ -n "$out" ]; then
    printf '%s\n' "$out" | grep -v '^[[:space:]]*$' | head -3 | while IFS= read -r _l; do
      case "$_l" in
        *pong*)
          # Pong: tilde verde centrada respecto del [OK] (texto alineado).
          if [ "$UI_ESTILO" = 1 ]; then
            printf '%s  %s✓%s  %s\n' "$MARGEN" "$(fg 0 230 130)" "$RESET" "$_l"
          else
            printf '%s\n' "$_l"
          fi
          ;;
        *)
          marg "     $_l"
          ;;
      esac
    done
  fi
  case "$out" in
    *pong*)
      case "$out" in
        *DERP*|*derp*) ruta="relay (DERP)" ;;
        *via*)         ruta="directo" ;;
        *)             ruta="respuesta recibida" ;;
      esac
      estado OK "Tailscale ping: responde (${ruta})."
      OK_PING=1
      ;;
    *)
      estado X "✗ Tailscale ping: sin respuesta en 5 segundos (parece fuera de la red)."
      ;;
  esac

  if command -v ping >/dev/null 2>&1; then
    icmp_ok=0
    if [ "$OS" = "Darwin" ]; then
      # En macOS el timeout por paquete va en milisegundos. El spinner es
      # visible; el stdout de ping no se emite.
      run_spin_silencioso "comprobando el alcance (ping ICMP)..." ping -c 3 -W 2000 "$NODE_IP" && icmp_ok=1
    else
      run_spin_silencioso "comprobando el alcance (ping ICMP)..." ping -c 3 -W 2 "$NODE_IP" && icmp_ok=1
    fi
    if [ "$icmp_ok" = 1 ]; then
      estado OK "Ping ICMP: responde."
    else
      estado '!' "Ping ICMP sin respuesta (Tailscale puede bloquearlo; no es fatal)."
    fi
  else
    estado i "No hay comando ping en este sistema; se omite el ping ICMP."
  fi
}

# --- Comprobar que el servicio web del nodo responde ------------------------
check_web() {
  local code rc
  echo
  marg "Servicio web del nodo:"
  if ! command -v curl >/dev/null 2>&1; then
    estado '!' "No hay curl; no se pudo comprobar el servicio web del nodo."
    marg "    Los datos de administracion estan en los materiales de infraestructura."
    return 0
  fi
  # -k: el certificado del nodo es propio. curl imprime el codigo HTTP en stdout.
  # Solo una respuesta de servicio valida cuenta como exito: 2xx, o 3xx de
  # redireccion. Vacio, 000, texto inesperado, 1xx, 4xx y 5xx son fallos y el
  # resumen no debe quedar "conexion lista".
  code="$(run_spin "consultando el servicio web..." curl -k -s -o /dev/null -w '%{http_code}' --max-time 15 "$WEB_URL")"
  rc=$?
  code="${code//[!0-9]/}"
  if [ "$rc" -eq 0 ] && [ "${#code}" -eq 3 ] && [ "$code" -ge 200 ] && [ "$code" -le 399 ] 2>/dev/null; then
    OK_WEB=1
    estado OK "Servicio web del nodo: responde (HTTP ${code})."
  else
    estado X "Servicio web del nodo: no respondio (codigo ${code:-000})."
    marg "    Causa probable: Tailscale apagado o el nodo no esta arriba."
  fi
}

# --- Paso 4/5: Prueba de conexion (automatica) ------------------------------
step4_tests() {
  echo
  seccion "Paso 4/5 - Prueba de conexion"
  marg "Estas pruebas son automaticas: no hay que hacer nada."
  ping_node
  check_web
}

# --- Paso 5/5: Resumen final ------------------------------------------------
# Imprime la seccion y las cinco lineas de estado (sin la tarjeta).
resumen_estados() {
  seccion_fin "Paso 5/5 - Resumen"
  if [ "$OK_TS" = 1 ]; then estado OK "Tailscale instalado"; else estado X " Tailscale no instalado"; fi
  if [ "$OK_SESSION" = 1 ]; then estado OK "Sesion de Tailscale activa y verificada"; else estado X " Tailscale sin sesion verificada"; fi
  if [ "$OK_NODE" = 1 ]; then estado OK "Nodo ${NODE_NAME} visible"; else estado X " Nodo ${NODE_NAME} no visible"; fi
  if [ "$OK_PING" = 1 ]; then estado OK "Tailscale ping: responde"; else estado X " Tailscale ping: sin respuesta"; fi
  if [ "$OK_WEB" = 1 ]; then estado OK "Servicio web del nodo alcanzable"; else estado X " Servicio web del nodo no alcanzable"; fi
}

# Pantalla final interactiva: avisa donde quedo el detalle y muestra el resumen
# compacto a continuacion, con espaciado normal (sin limpieza, padding ni
# reposicion del cursor). Imprime resumen y tarjeta una sola vez.
pantalla_final() {
  restaurar_terminal
  printf '%s%sEl detalle de los pasos queda en el historial; ahora se muestra el resumen final.%s\n' "$MARGEN" "$(fg 220 220 240)" "$RESET"
  echo
  resumen_estados
  echo
  ok_line "Conexion lista."
  marg "Tailscale tiene que estar siempre encendido mientras trabaja."
  marg "La administracion del nodo esta en los materiales de infraestructura:"
  marg "'Nodo G1 - Mapa y guia de infraestructura.docx' y 'Tu nodo G1.md'."
  tarjeta_final "CONEXION LISTA" "Nodo G1 - todo en orden" "Lencina Nicolas"
}

summary() {  # $1 = "final" para la pantalla interactiva de exito (opcional)
  local final="${1:-}" todo_ok=0
  if [ "$OK_TS" = 1 ] && [ "$OK_SESSION" = 1 ] && [ "$OK_NODE" = 1 ] && [ "$OK_PING" = 1 ] && [ "$OK_WEB" = 1 ]; then
    todo_ok=1
  fi
  if [ "$final" = "final" ] && [ "$UI_ESTILO" = 1 ] && [ "$todo_ok" = 1 ]; then
    pantalla_final
    return 0
  fi
  echo
  resumen_estados
  echo
  if [ "$todo_ok" = 1 ]; then
    ok_line "Conexion lista."
    marg "Tailscale tiene que estar siempre encendido mientras trabaja."
    marg "La administracion del nodo esta en los materiales de infraestructura:"
    marg "'Nodo G1 - Mapa y guia de infraestructura.docx' y 'Tu nodo G1.md'."
    tarjeta_final "CONEXION LISTA" "Nodo G1 - todo en orden" "Lencina Nicolas"
    return 0
  fi
  marg "La conexion no quedo completa. Revise los [X] de arriba y la seccion"
  marg "'Problemas comunes' de su guia, y vuelva a correr el script."
  return 1
}

# ============================================================================
# Ejecucion
# ============================================================================
# Cursor/colores defensivos al salir o ante senales; no cambia el codigo de
# salida de la ruta normal.
trap 'restaurar_terminal' EXIT
trap 'restaurar_terminal; limpiar_login_bg; exit 130' INT
trap 'restaurar_terminal; limpiar_login_bg; exit 143' TERM
show_intro
if ! paso_con_reintento "Paso 1 - Tailscale instalado" step1_tailscale; then
  limpiar_login_bg; summary; exit 1
fi
if ! paso_con_reintento "Paso 2 - Iniciar sesion" step2_session; then
  limpiar_login_bg; summary; exit 1
fi
if ! paso_con_reintento "Paso 3 - Aceptar la invitacion" step3_invite; then
  limpiar_login_bg; summary; exit 1
fi
step4_tests
limpiar_login_bg

if summary final; then
  exit 0
fi
exit 1
