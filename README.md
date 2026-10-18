# Kiroku — Personal Finance Tracker

Self-hosted, multi-account, multi-currency personal finance tracker with Japanese household taxonomy and precision.

## Core Features
- **Multi-Account & Currency**: Track bank accounts, cash, investments, and liabilities with base currency conversion.
- **Japanese Household Taxonomy**: Hierarchical category structure (食費, 交通費, etc.) with budgeting rules.
- **Fast Entry & Automation**: Templates, recurring transactions (RRULE), and AI/OCR receipt scanning.
- **Analytics & Forecasting**: Cashflow forecasting, spending breakdowns, trends, and multi-merchant analytics.
- **Security**: JWT auth, Argon2, TOTP 2FA, Telegram approval bot, and temporary read-only sharing links.
- **Backups & Portability**: Automated backups, Excel/ZIP export/import.

## Quick Start

### Development
```bash
# Install dependencies
npm install

# Start dev servers (backend + frontend)
npm run dev
```

Backend runs on `http://localhost:8080`, frontend on `http://localhost:5173`.

### Production
```bash
# Build and start production via PM2
./start.sh --prod
```

## Project Structure
```
kiroku/
├── backend/          # Go API server (Chi, SQLite, Zerolog)
├── frontend/         # Vue 3, Vite, Pinia, Tailwind, ECharts
├── .agent/           # Agent plans and workflows
├── docs/             # Documentation & API specs
└── start.sh          # Unified start script
```

## Tech Stack
- **Backend**: Go 1.22+, Chi router, SQLite (pure Go `modernc.org/sqlite`), Zerolog
- **Frontend**: Vue 3, Vite, Pinia, Tailwind CSS, ECharts, Leaflet
- **Deployment**: PM2, Cloudflare Tunnel, Upstash Redis, Telegram Bot

## Environment
Copy `backend/.env.example` to `backend/.env` and configure:
- `JWT_SECRET`: Min 32 characters for token signing
- `UPSTASH_REDIS_REST_URL` / `UPSTASH_REDIS_REST_TOKEN`: Cloudflare tunnel URL discovery
- `BACKUP_PASSPHRASE`: Encryption passphrase for backups
- `TELEGRAM_BOT_TOKEN` / `TELEGRAM_CHAT_ID`: Telegram integration