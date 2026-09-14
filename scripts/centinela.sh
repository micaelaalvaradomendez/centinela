#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

command -v docker >/dev/null 2>&1 || {
  echo "Docker no está instalado o no está disponible en PATH." >&2
  exit 1
}

case "${1:-up}" in
  up)
    git submodule update --init --recursive
    docker compose up --build --detach --wait
    docker compose ps
    echo "Centinela disponible en http://localhost:${FRONTEND_PORT:-8088}"
    ;;
  down)
    docker compose down
    ;;
  logs)
    docker compose logs --follow
    ;;
  status)
    docker compose ps
    ;;
  rebuild)
    git submodule update --init --recursive
    docker compose up --build --detach --wait --force-recreate
    docker compose ps
    ;;
  reset)
    docker compose down --volumes
    ;;
  *)
    echo "Uso: $0 {up|down|logs|status|rebuild|reset}" >&2
    exit 2
    ;;
esac