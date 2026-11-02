#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT

MODE="${1:-}"

case "$MODE" in
  stop)
    echo "Stopping all Kiroku processes..."
    pkill -f "kiroku-(backend|frontend)" 2>/dev/null || true
    echo "All processes stopped."
    exit 0
    ;;
  start)
    shift
    SERVICE="${1:-}"
    if [[ -z "$SERVICE" ]]; then
      echo "Usage: $0 start [fe|be]"
      exit 1
    fi
    ;;
  *)
    echo "Usage: $0 start [fe|be]"
    exit 1
    ;;
esac

case "$SERVICE" in
  fe)
    set -a
    source "$PROJECT_ROOT/frontend/.env"
    set +a

    cd "$PROJECT_ROOT/frontend"
    npm run dev -- --host 0.0.0.0 --port 5173
    ;;
  be)
    set -a
    source "$PROJECT_ROOT/backend/.env"
    set +a

    cd "$PROJECT_ROOT/backend"
    export DATA_DIR="$PROJECT_ROOT/backend/$DATA_DIR"
    CGO_ENABLED=0 air
    ;;
  *)
    echo "Usage: $0 start [fe|be]"
    exit 1
    ;;
esac