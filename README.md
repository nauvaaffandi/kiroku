# Kiroku (記録)

A self-hosted personal finance tracker. Record income, expenses and transfers across multiple accounts and currencies from a mobile-first PWA, backed by a single Go binary and one SQLite file.

Single user, your data stays on your machine.

## Features

- **Transactions**: income, expense and transfer, with search, batch operations, templates and soft delete
- **Accounts and categories**: multi-currency accounts, hierarchical categories
- **Scheduled transactions**: recurring entries via RRULE
- **Reports**: spending by category and trend charts (ECharts)
- **Auth**: 6-digit PIN with JWT sessions, optional TOTP, PIN reset through a Telegram bot
- **Shared links**: public read-only views of selected data
- **FX rates**: fetched and stored against a configurable base currency
- **i18n and theme**: English, 日本語, Bahasa Indonesia; light and dark mode
- **PWA**: installable, works offline for cached views

> **Status:** the backend packages for analytics, backup/restore, receipt scanning (OCR + LLM) and CSV/Excel import-export exist under `backend/internal/` but are not mounted in `cmd/server/main.go` yet. See [AGENTS.md](./AGENTS.md#known-gaps).

## Tech stack

| Layer    | Choice                                                                                      |
| -------- | ------------------------------------------------------------------------------------------- |
| Backend  | Go 1.27+, chi, SQLite (`modernc.org/sqlite`, pure Go, WAL), zerolog, koanf, Argon2id, JWT   |
| Frontend | Vue 3 (`<script setup lang="ts">`), Vite, Pinia, Vue Router, Tailwind CSS v4, ECharts, PWA  |
| Ops      | PM2, Cloudflare Tunnel (`cloudflared`), Upstash Redis for backend URL discovery             |

## Getting started

### Prerequisites

- Go 1.27+
- Node.js 18+ and npm
- Optional: [`air`](https://github.com/air-verse/air) for backend live reload, `pm2` and `cloudflared` for `start.sh`

> `start.sh` currently points `--dev`/`--watch` at `air` in the Termux path (`/data/data/com.termux/files/home/go/bin/air`). Edit `BACKEND_SCRIPT` in `start.sh` if you run elsewhere.

### Run locally

```bash
# 1. Backend (http://localhost:8080)
cd backend
KIROKU_TUNNEL_ENABLED=false go run ./cmd/server

# 2. Frontend (http://localhost:5173)
cd frontend
npm install
npm run dev
```

Without Upstash credentials the frontend falls back to `http://localhost:8080`, so no extra setup is needed in development.

> **First login:** the user is created on first sign-in with the default PIN `123456`. Change it right away from **Settings**, especially before exposing the backend through a tunnel.

### One-command start (PM2)

```bash
./start.sh --dev     # backend with air + frontend dev server
./start.sh --prod    # builds backend/bin/server, then starts everything
./start.sh stop
```

`start.sh` loads `.env` from the repo root and regenerates `config/ecosystem.config.js`.

## Configuration

Backend settings live in `backend/config.yaml` (address, data dir, JWT TTLs, Argon2 params, FX base currency, backup schedule and retention).

```yaml
addr: ":8080"
data_dir: "./data"
db_path: "kiroku.db"
jwt_secret: "change-me-to-at-least-32-random-chars"
fx_base_currency: "JPY"
backup_schedule: "0 3 * * *"
```

Environment variables (via `.env`, see `.env.example`) are read for these:

| Variable                                                           | Purpose                                      |
| ------------------------------------------------------------------ | -------------------------------------------- |
| `TELEGRAM_BOT_TOKEN`, `TELEGRAM_BOT_USERNAME`, `TELEGRAM_WEBHOOK_SECRET` | Telegram bot for PIN reset and change   |
| `UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN`               | Publish the tunnel URL for the frontend      |
| `KIROKU_TUNNEL_ENABLED` / `CLOUDFLARE_TUNNEL_ENABLED`              | Set to `false` to disable the tunnel         |
| `KIROKU_LOG_RETENTION`                                             | Log retention, default `90d`                 |
| `KIROKU_BACKUP_GIT_TOKEN`                                          | Token for the git backup repo                |
| `KIROKU_DEV`                                                       | `true` for console logging                   |
| `VITE_UPSTASH_URL`, `VITE_FRONTEND_UPSTASH_REDIS_REST_TOKEN`       | Frontend: look up the backend URL            |
| `VITE_API_BASE_URL`                                                | Frontend: fixed backend URL (non-localhost)  |

**Set `jwt_secret` in `config.yaml` before exposing the app to a network.**

## Usage

1. Open the app, sign in with the default PIN (`123456`), then change it in **Settings**.
2. Create accounts, then categories (a default set is seeded on first run).
3. Add transactions from the `+` button. Use templates for repeated entries.
4. Add scheduled transactions for rent, subscriptions and salary.
5. Check **Reports** for spending and trends. Switch theme or language in **Settings**.

Swagger UI is mounted at `http://localhost:8080/api/docs/` (the spec has no annotated endpoints yet). Health check: `GET /healthz`. The REST API lives under `/api/v1`.

## Project structure

```
kiroku/
├── backend/
│   ├── cmd/server/          # entrypoint, router, middleware
│   ├── docs/                # generated Swagger
│   └── internal/
│       ├── domain/          # entities + repository interfaces
│       ├── repository/sqlite/   # repository implementations
│       ├── delivery/http/   # HTTP handlers
│       ├── auth/            # PIN, JWT, TOTP, Telegram
│       ├── db/migrations/   # embedded SQL migrations
│       └── config, logger, tunnel, telegram, models
├── frontend/
│   └── src/
│       ├── views/           # route-level pages
│       ├── widgets/         # feature blocks per domain
│       ├── components/      # ui, layout, common, charts, icons
│       ├── stores/          # Pinia
│       ├── locales/         # en, ja, id
│       └── assets/main.css  # Tailwind entry + theme tokens
├── config/                  # PM2 config (generated)
├── start.sh                 # dev / prod process manager
└── DESIGN.md                # UI design system
```

## Scripts

| Where       | Command                              | What it does                    |
| ----------- | ------------------------------------ | ------------------------------- |
| `backend/`  | `go run ./cmd/server`                | Run the API                     |
| `backend/`  | `go test ./...`                      | Backend tests                   |
| `backend/`  | `CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server` | Production binary |
| `frontend/` | `npm run dev`                        | Vite dev server                 |
| `frontend/` | `npm run build`                      | Type-check and production build |
| `frontend/` | `npm run test:run`                   | Vitest (tests go in `tests/`)   |

## Deployment

Production splits the two halves:

- **Backend** runs from `start.sh --prod` under PM2. It starts `cloudflared`, reads the random `trycloudflare.com` URL from its output, and writes it to Upstash Redis under `kiroku:backend:url`.
- **Frontend** is a static build (for example on Vercel). On load it reads the current backend URL from Upstash, so a tunnel restart needs no redeploy.

## Contributing

Working with an AI coding agent? Read [AGENTS.md](./AGENTS.md) for conventions and [DESIGN.md](./DESIGN.md) for the UI rules.

## License

MIT
