#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT

MODE="${1:---dev}"

if [[ "$MODE" == "stop" ]]; then
    echo "Stopping all Kiroku processes..."
    pm2 delete all 2>/dev/null || true
    echo "All processes stopped."
    exit 0
fi

if [[ "$MODE" != "--dev" && "$MODE" != "--watch" && "$MODE" != "--prod" ]]; then
    echo "Usage: $0 [--dev|--watch|--prod|stop]"
    exit 1
fi

cd "$PROJECT_ROOT"

pm2 delete all 2>/dev/null || true

mkdir -p config backend/data backend/bin backend/logs frontend/logs

BACKEND_SCRIPT="/data/data/com.termux/files/home/go/bin/air"
if [[ "$MODE" == "--prod" ]]; then
    echo "Building production binary..."
    cd "$PROJECT_ROOT/backend"
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server cmd/server/main.go
    cd "$PROJECT_ROOT"
    BACKEND_SCRIPT="./bin/server"
fi

cat > config/ecosystem.config.js <<EOF
module.exports = {
  apps: [
    {
      name: "kiroku-backend",
      cwd: process.env.PROJECT_ROOT + "/backend",
      script: "${BACKEND_SCRIPT}",
      interpreter: "none",
      watch: false,
      env: {
        PORT: "8080",
        HOST: "0.0.0.0",
        INSITU_ADDR: ":8080",
        INSITU_DATA_DIR: process.env.PROJECT_ROOT + "/backend/data",
        CGO_ENABLED: "0"
      },
      log_date_format: "YYYY-MM-DD HH:mm:ss Z",
      merge_logs: true,
      error_file: process.env.PROJECT_ROOT + "/backend/logs/backend-error.log",
      out_file: process.env.PROJECT_ROOT + "/backend/logs/backend-out.log"
    },
    {
      name: "kiroku-frontend",
      cwd: process.env.PROJECT_ROOT + "/frontend",
      script: "npm",
      args: "run dev -- --host 0.0.0.0 --port 5173",
      interpreter: "none",
      watch: false,
      env: {
        HOST: "0.0.0.0",
        PORT: "5173"
      },
      log_date_format: "YYYY-MM-DD HH:mm:ss Z",
      merge_logs: true,
      error_file: process.env.PROJECT_ROOT + "/frontend/logs/frontend-error.log",
      out_file: process.env.PROJECT_ROOT + "/frontend/logs/frontend-out.log"
    }
  ]
};
EOF

pm2 start config/ecosystem.config.js

echo "Unified services started in background via PM2 (${MODE})."
echo "Cloudflare Tunnel and Upstash discovery active."
echo "Check status: pm2 status"
echo "View logs: pm2 logs"
echo "Stop all: $0 stop"
