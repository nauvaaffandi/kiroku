# Kiroku — Personal Finance Tracker

Self-hosted, single-user, multi-account/multi-currency personal finance tracker with Japanese precision.

## Quick Start

### Development
```bash
# Install dependencies
npm install

# Start dev servers (backend + frontend)
npm run dev
```

Backend runs on http://localhost:8080
Frontend runs on http://localhost:5173

### Production
```bash
# Build and start production
./start.sh --prod
```

## Project Structure
```
kiroku/
├── backend/          # Go API server
│   ├── cmd/server/   # Entry point
│   ├── internal/     # Private packages
│   ├── pkg/          # Public packages
│   └── data/         # SQLite database (gitignored)
├── frontend/         # Vue 3 + Vite + Pinia
│   ├── src/
│   └── public/
├── config/           # PM2 ecosystem config (generated)
├── .tasks/           # Task definitions
├── docs/             # Documentation
├── start.sh          # Unified start script
├── PRD.md            # Product requirements
└── AGENTS.md         # AI agent directives
```

## Tech Stack
- **Backend**: Go 1.22+, Chi router, SQLite (modernc.org/sqlite), sqlc
- **Frontend**: Vue 3, Vite, Pinia, ECharts, Leaflet
- **Deployment**: PM2, Cloudflare Tunnel, Upstash Redis
- **PWA**: Workbox via Vite PWA plugin

## Scripts
| Command | Description |
|---------|-------------|
| `npm run dev` | Start both dev servers |
| `npm run build` | Build both apps |
| `npm run lint` | Lint both apps |
| `npm run test` | Run tests |
| `./start.sh --dev` | Dev mode with PM2 |
| `./start.sh --prod` | Production mode |

## Environment
Copy `.env.example` to `.env` and configure:
- `JWT_SECRET`: Min 32 chars for token signing
- `UPSTASH_REDIS_REST_URL/TOKEN`: For Cloudflare tunnel URL discovery
- `BACKUP_PASSPHRASE`: Optional encryption for backups