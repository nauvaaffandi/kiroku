# 記録 (Kiroku)

### Enterprise-Grade, Self-Hosted Personal Finance Tracker & Wealth Management Platform
*A privacy-first, offline-capable mobile Progressive Web App (PWA) powered by a modular Go backend and a high-performance Vue 3 dashboard.*

[![Go Version](https://img.shields.io/badge/Go-1.27-%2300ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.5-%234FC08D?style=for-the-badge&logo=vue.js)](https://vuejs.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-v4-%2338BDF8?style=for-the-badge&logo=tailwindcss)](https://tailwindcss.com)
[![SQLite](https://img.shields.io/badge/SQLite-WAL-%23003B57?style=for-the-badge&logo=sqlite)](https://sqlite.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

---

## 📖 Table of Contents

1. [Introduction: What is Kiroku?](#-introduction-what-is-kiroku)
2. [Core Philosophy & Design Principles](#-core-philosophy--design-principles)
3. [Key Features & Capabilities](#-key-features--capabilities)
4. [Monorepo Architecture Overview](#-monorepo-architecture-overview)
5. [Deep Dive: Backend Architecture (`backend/`)](#-deep-dive-backend-architecture-backend)
   - [Composition Root (`cmd/server`)](#composition-root-cmdserver)
   - [Domain Layer & Business Entities](#domain-layer--business-entities)
   - [SQLite Persistence & Pure Go Driver (`modernc.org/sqlite`)](#sqlite-persistence--pure-go-driver-moderncorgsqlite)
   - [Embedded Migrations & Schema Management](#embedded-migrations--schema-management)
   - [HTTP Delivery Layer (Chi v5)](#http-delivery-layer-chi-v5)
   - [Configuration & Environment Management (`Koanf`)](#configuration--environment-management-koanf)
   - [Logging (`zerolog`)](#logging-zerolog)
6. [Deep Dive: Frontend Architecture (`frontend/`)](#-deep-dive-frontend-architecture-frontend)
   - [Vue 3 Composition API & TypeScript](#vue-3-composition-api--typescript)
   - [Pinia State Management Stores](#pinia-state-management-stores)
   - [Routing & Navigation Guards](#routing--navigation-guards)
   - [Tailwind CSS v4 Styling & Semantic Design Tokens](#tailwind-css-v4-styling--semantic-design-tokens)
   - [Internationalization (i18n): EN, JA, ID](#internationalization-i18n-en-ja-id)
   - [Progressive Web App (PWA) & Offline Capabilities](#progressive-web-app-pwa--offline-capabilities)
7. [Complete Database Schema Reference](#-complete-database-schema-reference)
8. [Dependencies](#-dependencies)
   - [Backend (Go modules)](#backend-go-modules)
   - [Frontend (npm packages)](#frontend-npm-packages)
9. [Getting Started & Local Development](#-getting-started--local-development)
   - [Prerequisites](#prerequisites)
   - [Running the Backend](#running-the-backend)
   - [Running the Frontend](#running-the-frontend)
   - [Running Tests & Builds](#running-tests--builds)
10. [Deployment & Operations Guide](#-deployment--operations-guide)
11. [Known Security Considerations & Architectural Gaps](#-known-security-considerations--architectural-gaps)
12. [Contributing & License](#-contributing--license)

---

## 🌟 Introduction: What is Kiroku?

**Kiroku (記録 — meaning "Record" or "Log" in Japanese)** is an advanced, privacy-first, self-hosted personal finance tracking and wealth management web application. In an era where financial data is routinely surrendered to third-party subscription services, Kiroku offers an uncompromising alternative: **complete data sovereignty combined with an enterprise-grade user experience**.

Designed from the ground up to operate seamlessly across mobile and desktop devices as a Progressive Web App (PWA), Kiroku provides meticulous tracking for multi-currency accounts, transactions, transfers, budgets, cash flow projections, and receipt scanning. Because it runs on your own infrastructure backed by a lightweight SQLite database, your financial records never leave your control.

---

## 💡 Core Philosophy & Design Principles

Kiroku is built upon several non-negotiable engineering and product tenets:

1. **Privacy Above All:** No third-party analytics, no cloud telemetry, and no mandatory external SaaS dependencies. Your financial data stays entirely on your server or device.
2. **Offline-First Resilience:** Financial tracking should not depend on uninterrupted internet connectivity. Thanks to robust PWA service worker caching (`NetworkFirst` strategy) and local storage adapters, Kiroku works seamlessly offline.
3. **Impeccable Ergonomics:** Inspired by modern consumer fintech apps, the user interface features a centered column layout optimized for mobile thumbs, immediate keyboard navigation, clear visual hierarchies, and dark/light mode reactivity.
4. **Zero CGO Complexity:** The backend leverages pure Go SQLite drivers (`modernc.org/sqlite`), enabling seamless cross-platform compilation without requiring external C compilers or complex build toolchains.
5. **Strict Internationalization:** Every user-facing string is fully internationalized across English (`en`), Japanese (`ja`), and Indonesian (`id`), ensuring native accessibility for diverse user bases.

---

## 🚀 Key Features & Capabilities

- **Multi-Account & Multi-Currency Management:** Group accounts by type (Checking, Savings, Credit Cards, Investments, Cash) and manage transactions across multiple global currencies with automatic or manual exchange rates (`/fx`).
- **Granular Ledger & Categories:** Categorize income and expenses hierarchically with custom icons and color tags. Support for split transactions, transfers between accounts, and recurring templates.
- **Scheduled & Recurring Transactions:** Automate recurring bills, salaries, and subscriptions with configurable frequency rules and execution logs.
- **Cash Flow Projections:** Visualize future financial health using ECharts-driven projection curves that extrapolate recurring income and committed expenses.
- **Budgeting & Spending Limits:** Establish monthly or custom-period budgets per category with real-time visual progress indicators and threshold alerts.
- **Receipt Scanning & OCR Integration:** Capture paper receipts via camera or file upload, processed asynchronously through backend ingestion pipelines (`receipt_jobs`).
- **Spending Reports & Cashflow Table:** Category-based spending breakdowns with parent category grouping, per-account filtering, and interval-based cashflow summaries (day/week/month/quarter/year).
- **Secure Authentication:** PIN and password-protected sessions powered by Argon2 hashing and secure JWT authentication cookies.
- **Shared Access:** Generate time-limited public share tokens to expose specific account activity without exposing full credentials.

---

## 📂 Monorepo Architecture Overview

Kiroku is structured as a tightly integrated monorepo containing both the Go backend server and the Vue 3 frontend client:

- **`backend/`**: A modular Go 1.27 backend following clean architecture principles. It exposes a RESTful JSON API under `/api/v1` and manages persistence via SQLite.
- **`frontend/`**: A Vue 3 single-page application built with Vite, Pinia, Vue Router, and Tailwind CSS v4, compiled into an installable PWA.
- **`start.sh`**: Operational wrapper script for process management and cloud tunnel orchestration.
- **`AGENTS.md` / `DESIGN.md`**: Binding project instructions and design system rules for AI coding agents and contributors.

---

## 🛠️ Deep Dive: Backend Architecture (`backend/`)

The backend is structured as a modular monolith. It deliberately avoids heavy framework magic in favor of explicit, testable, and maintainable Go code.

```
backend/
├── cmd/server/main.go          # Application composition root (wiring + router setup)
├── internal/
│   ├── config/                 # Koanf + dotenv configuration loader
│   ├── db/                     # SQLite pool, migration runner, WAL mode
│   │   └── migrations/         # Versioned SQL files (001_initial.sql ... 020_backup_multi_provider_state.sql)
│   ├── models/                 # Shared DTOs / request-response models
│   ├── auth/                   # JWT, PIN, Argon2, TOTP, genesis user, rate limiting
│   ├── logger/                 # zerolog structured logging (console + lumberjack rotation)
│   ├── tunnel/                 # Cloudflared tunnel manager
│   ├── delivery/http/          # Chi HTTP routing & REST JSON handlers
│   ├── domain/                 # Core business entities & repository interfaces
│   ├── repository/sqlite/      # Concrete SQLite repository adapters
│   ├── analytics/              # Analytics cashflow, balance, outlook, budget, goals, stats
│   ├── transactions/           # Transaction CRUD, list, batch, search, summary
│   ├── accounts/               # Account CRUD & summary
│   ├── categories/             # Category tree CRUD & helpers
│   ├── scheduled/              # Scheduled/recurring transaction CRUD
│   ├── settings/               # Settings CRUD
│   ├── fx/                     # Foreign exchange rates & conversion
│   ├── templates/              # Transaction templates
│   ├── labels/                 # Transaction labels
│   ├── push/                   # Web Push subscriptions
│   ├── shared/                 # Public share token management (authenticated + public view)
│   ├── receipts/               # Receipt ingestion & OCR pipeline
│   ├── impexp/                 # Import/export pipeline
│   ├── backup/                 # Encrypted backup scheduler & providers (local/git/gdrive)
│   └── telegram/               # Telegram bot service
├── docs/                       # Swagger handler stubs
└── go.mod                      # Go module definition (module kiroku/backend)
```

### Composition Root (`cmd/server`)
`cmd/server/main.go` wires the application together:
1. Loads configuration (`config.Load()`).
2. Opens SQLite (`db.Open()`).
3. Ensures genesis user (`auth.EnsureGenesisUser`).
4. Instantiates repositories (`NewAccountsRepo`, `NewCategoriesRepo`, `NewTransactionsRepo`, `NewReportsRepo`, etc.).
5. Mounts handlers under `/api/v1` behind `auth.AuthMiddleware(cfg)`.
6. Mounts shared public routes at `/shared/public/{token}`.
7. Starts the HTTP server and (optionally) the Cloudflare tunnel manager.

### Domain Layer & Business Entities
The `domain/` package houses pure Go structs and interface contracts. By decoupling domain logic from HTTP delivery and SQL persistence, Kiroku ensures high testability. Key entities include:
- `User` — Authentication and user profile settings.
- `Account` — Financial accounts with balances and currency definitions.
- `Category` — Income/expense classification trees with parent/child hierarchy.
- `Transaction` — Ledger entries linking accounts, categories, amounts (stored as integers to prevent floating-point inaccuracies), dates, and notes.
- `ScheduledTransaction`, `Template`, `Label`, `FXRate`, `Share`, `PushSubscription`.

### SQLite Persistence & Pure Go Driver (`modernc.org/sqlite`)
To maintain portability across diverse server environments (Linux, macOS, Windows, ARM/x86), Kiroku uses `modernc.org/sqlite`:
- **WAL Mode:** SQLite connection pools are explicitly configured in Write-Ahead Logging (`WAL`) mode to support concurrent reads while writes are executed.
- **Foreign Keys:** Enforced at the database connection level (`PRAGMA foreign_keys = ON;`).
- **Integer Money Representation:** Monetary amounts are stored as raw `INTEGER` values (representing minor currency units / cents) in SQLite to guarantee zero precision loss during arithmetic operations.

### Embedded Migrations & Schema Management
Database schema evolution is handled deterministically via embedded SQL migration scripts stored in `internal/db/migrations/`:
- Migrations are embedded into the Go binary using `embed.FS`.
- At startup, the migration runner queries the `schema_version` table and applies any outstanding `NNN_description.sql` files in ascending numerical order.
- Current migration count: `001` through `020` (21 files).

### HTTP Delivery Layer (Chi v5)
The HTTP server uses `go-chi/chi/v5` for lightweight, idiomatic routing:
- Handlers in `delivery/http/` decode JSON payloads, validate input structures, invoke domain repository methods, and marshal responses back to JSON.
- Middleware handles structured logging (`rs/zerolog`), CORS policies, real-IP detection, and authentication verification.

### Configuration & Environment Management (`Koanf`)
Configuration is managed via `github.com/knadh/koanf/v2` combined with `.env` file loading:
- Default values loaded from `config.yaml` / dotenv source.
- Environment variables prefixed with `KIROKU_` override YAML settings by translating underscores into nested configuration keys.
- Key environment variables include:
  - `PORT` / `KIROKU_ADDR` — HTTP listen address.
  - `DATA_DIR` — Absolute path to the data directory (database, persistent files).
  - `JWT_SECRET` — JWT signing secret.
  - `FX_BASE_CURRENCY` — Base currency for exchange rates (default `JPY`).
  - `TELEGRAM_BOT_TOKEN` — Telegram bot token.
  - `UPSTASH_REDIS_REST_URL` / `BACKEND_UPSTASH_REDIS_REST_TOKEN` — Upstash Redis REST for tunnel URL publishing.
  - `CLOUDFLARE_TUNNEL_ENABLED` — Toggle cloudflared tunnel.
  - `LOG_LEVEL`, `LOG_FORMAT`, `LOG_RETENTION` — Logging controls.
  - `ARGON2_TIME`, `ARGON2_MEMORY`, `ARGON2_THREADS` — Argon2 password hashing costs.
  - `VAPID_PUBLIC_KEY` / `VAPID_PRIVATE_KEY` — Web Push VAPID keys.
  - `BACKUP_PASSPHRASE`, `BACKUP_REPO_PATH`, `BACKUP_GIT_TOKEN` — Backup encryption & storage.

### Logging (`zerolog`)
Backend code uses `github.com/rs/zerolog/log` at appropriate severity levels (`trace`, `debug`, `info`, `warn`, `error`, `panic`, `fatal`). Structured request logging is emitted via `logger.RequestLoggingMiddleware`, recording component, method, path, status, duration, client IP, and user agent. No sensitive request data (bodies, credentials, tokens) is logged.

---

## 🖥️ Deep Dive: Frontend Architecture (`frontend/`)

The frontend is a modern Single Page Application (SPA) designed as an installable Progressive Web App (PWA).

```
frontend/
├── src/
│   ├── api/                    # Axios client instance & interceptors
│   ├── assets/                 # Global stylesheet (main.css) & design tokens
│   ├── components/             # Reusable UI controls, widgets, charts, icons, layout
│   ├── composables/            # Vue composables (useBiometric, useCurrencies, useDebounce, useHaptics, usePagination)
│   ├── locales/                # i18n dictionaries (en.ts, ja.ts, id.ts, index.ts)
│   ├── router/                 # Vue Router config & navigation guards
│   ├── stores/                 # Pinia state stores (auth, accounts, categories, transactions, scheduled, settings, reports, ui, labels, templates, push, backendUrl)
│   ├── types/                  # TypeScript interface definitions (account, auth, category, report, settings, scheduled, transaction, ui, etc.)
│   ├── utils/                  # Formatters (amountDisplay), exportExcel, categoryName
│   ├── views/                  # Route views (Dashboard, Accounts, Transactions, Categories, Scheduled, Reports, Settings, Analytics*, Labels, Shared, Calendar, About, Login)
│   ├── App.vue                 # Root shell component
│   └── main.ts                 # Application bootstrapping
├── index.html                  # HTML entry point with PWA meta tags
├── vite.config.ts              # Vite & Vite PWA plugin configuration
├── tailwind.config.js          # Tailwind CSS v4 configuration
└── package.json                # npm scripts & dependencies
```

### Vue 3 Composition API & TypeScript
All components use Vue 3 `<script setup lang="ts">` syntax, providing strict type safety, reactive refs, computed properties, and reusable composables.

### Pinia State Management Stores
Application state is organized into modular Pinia stores located in `src/stores/`:
- Handles API communication via shared Axios clients (`src/api/client.ts`).
- Manages reactive caching, optimistic UI updates, and synchronization with backend endpoints.
- Stores: `auth`, `accounts`, `categories`, `transactions`, `scheduled`, `settings`, `reports`, `ui`, `labels`, `templates`, `push`, `backendUrl`, `i18n`.

### Routing & Navigation Guards
Vue Router manages client-side page transitions between core views:
- **Dashboard:** Financial overview, net worth summary, recent transactions.
- **Accounts:** Detailed breakdown of asset and liability accounts.
- **Transactions:** Filterable, searchable ledger management table (incl. receipt scan).
- **Categories:** Category tree management.
- **Scheduled:** Recurring transaction management.
- **Reports:** Spending report (category / parent_category / account filtering) and cashflow table.
- **Analytics:** Balance trend, accounts, currencies, outlook (category/timeline/forecast), budget vs actual, cashflow projection, spending analytics, goals progress.
- **Settings:** User preferences, currency settings, theme toggle, data export.
- **SharedAccess / SharedView:** Create & view time-limited public share tokens.
- **Labels:** Transaction label management.
- **About / Calendar:** Miscellaneous utility views.

### Tailwind CSS v4 Styling & Semantic Design Tokens
Kiroku adheres strictly to a design token system defined in `src/assets/main.css`:
- **Tailwind via Style Tag:** Tailwind is consumed exclusively inside component `<style scoped>` blocks via `@reference "tailwindcss";` and `@apply`.
- **Semantic Color Variables:** Components reference semantic tokens (`--bg-primary`, `--bg-card`, `--text-primary`, `--accent`, `--border-primary`, `--success`, `--error`, etc.) rather than hardcoded hex codes. This guarantees instant, flicker-free light and dark mode transitions across all UI elements.
- **Theme Control:** Theme is controlled by `data-theme="light" | "dark"` on `<html>` (no attribute means follow OS). Dark mode selector is `:where([data-theme="dark"], [data-theme="dark"] *)`.

### Internationalization (i18n): EN, JA, ID
To support global users, all UI text is externalized into three translation dictionaries managed by `src/locales/index.ts`:
- `en.ts` (English — Default fallback)
- `ja.ts` (Japanese)
- `id.ts` (Indonesian)

Every new user-facing key must be added to all three files. The `useI18n()` composable provides the `t(key)` function with fallback to English.

### Progressive Web App (PWA) & Offline Capabilities
Configured via Vite PWA plugin (`vite-plugin-pwa`), Kiroku caches static assets and API responses (`NetworkFirst` strategy for `/api/v1/*`), allowing users to record transactions and review accounts even when completely offline.

---

## 📊 Complete Database Schema Reference

Below is an overview of the core database tables managed by Kiroku's SQLite migration engine:

| Table Name | Description | Key Columns |
|---|---|---|
| `schema_version` | Tracks applied migration files | `version` (INTEGER, PK), `applied_at` |
| `users` | User accounts & credentials | `id`, `username`, `pin_hash`, `created_at` |
| `accounts` | Financial accounts & wallets | `id`, `user_id`, `name`, `type`, `currency`, `balance`, `deleted_at` |
| `categories` | Income/Expense category tree | `id`, `name`, `type`, `parent_id`, `icon`, `color`, `archived` |
| `transactions` | Financial ledger movements | `id`, `from_account_id`, `to_account_id`, `category_id`, `amount` (INTEGER in cents), `type`, `date_time`, `notes`, `deleted_at` |
| `budgets` | Spending limits & thresholds | `id`, `category_id`, `amount`, `period`, `start_date`, `end_date` |
| `scheduled_transactions` | Recurrence & automation rules | `id`, `account_id`, `category_id`, `amount`, `frequency`, `next_run`, `active` |
| `templates` | Reusable transaction templates | `id`, `account_id`, `category_id`, `amount`, `type`, `notes` |
| `labels` | Transaction labels | `id`, `name`, `color`, `sort_order` |
| `transaction_labels` | Transaction-label junction | `transaction_id`, `label_id` |
| `fx_rates` | Foreign exchange rates | `base`, `quote`, `rate`, `updated_at` |
| `settings` | User settings key-value | `key`, `value` |
| `shares` | Public share tokens | `id`, `account_id`, `token`, `expires_at`, `revoked_at` |
| `telegram_bindings` | Telegram chat binding | `user_id`, `chat_id`, `token` |
| `push_subscriptions` | Web Push subscriptions | `id`, `endpoint`, `p256dh`, `auth` |
| `login_attempts` | Login attempt log | `id`, `identifier`, `success`, `attempted_at` |
| `receipts` | Receipt images & OCR metadata | `id`, `image_path`, `merchant`, `total`, `parsed_data` |
| `receipt_jobs` | Background OCR processing queue | `id`, `status`, `payload`, `error`, `created_at` |
| `backup_state` | Backup provider state | provider-specific state columns |
| `backup_multi_provider_state` | Multi-provider backup state | provider-specific state columns |

---

## 📦 Dependencies

### Backend (Go modules)
Module path: `kiroku/backend`

| Package | Purpose |
|---|---|
| `github.com/go-chi/chi/v5` | HTTP router |
| `github.com/golang-jwt/jwt/v5` | JWT signing & validation |
| `github.com/rs/zerolog` | Structured logging |
| `modernc.org/sqlite` | Pure Go SQLite driver |
| `github.com/natefinch/lumberjack` | Log rotation |
| `github.com/pquerna/otp` | TOTP (two-factor) |
| `golang.org/x/crypto` | Argon2 password hashing |
| `github.com/teambition/rrule-go` | Recurrence rule generation |
| `github.com/xuri/excelize/v2` | Excel export/import |
| `github.com/swaggo/http-swagger/v2` | Swagger docs |
| `github.com/stretchr/testify` | Test assertions |

### Frontend (npm packages)
Package name: `kiroku-frontend`

| Package | Purpose |
|---|---|
| `vue` | Vue 3 framework |
| `vue-router` | Client-side routing |
| `pinia` | State management |
| `axios` | HTTP client |
| `echarts` / `vue-echarts` | Charts & visualizations |
| `@lucide/vue` | Icon set |
| `dayjs` | Date manipulation |
| `tailwindcss` ^4 | CSS framework |
| `@tailwindcss/vite` | Tailwind Vite integration |
| `vite-plugin-pwa` | PWA service worker generation |
| `vite` | Build tool & dev server |
| `vue-tsc` / `typescript` | TypeScript type-checking |
| `vitest` | Unit testing |

---

## 💻 Getting Started & Local Development

### Prerequisites
Ensure your development environment has the following installed:
- **Go 1.27+** (`go version`)
- **Node.js 18+ & npm** (`node -v && npm -v`)

### Running the Backend

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Configure your runtime settings (edit `config.yaml` or set environment variables).
3. Start the Go server (with tunnel disabled for local testing):
   ```bash
   KIROKU_TUNNEL_ENABLED=false go run ./cmd/server
   ```
4. Build the static binary:
   ```bash
   CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server
   ```

### Running the Frontend

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the Vite development server:
   ```bash
   npm run dev
   ```

### Running Tests & Builds

- **Backend Tests & Compilation:**
  ```bash
  cd backend
  go build ./...
  go test ./...
  ```
- **Frontend Type-Check & Production Build:**
  ```bash
  cd frontend
  npm run build
  ```
- **Frontend Unit Tests:**
  ```bash
  cd frontend
  npm run test:run
  ```

---

## 🚢 Deployment & Operations Guide

For production self-hosting:

1. Compile the backend statically with CGO disabled:
   ```bash
   cd backend
   CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server
   ```
2. Build the frontend static assets:
   ```bash
   cd frontend
   npm run build
   ```
3. Use the provided `start.sh` script to orchestrate process management:
   ```bash
   ./start.sh start be
   ./start.sh start fe
   ```
4. For remote access, enable the Cloudflare tunnel (`CLOUDFLARE_TUNNEL_ENABLED=true`) and configure Upstash Redis for tunnel URL publishing.

---

## ⚠️ Known Security Considerations & Architectural Gaps

Operators deploying Kiroku in production should review the following known limitations:

1. **Data Route Authentication:** In `cmd/server/main.go`, certain data routes under `/api/v1` (such as `/accounts`, `/categories`, `/transactions`, `/reports`) are mounted behind `auth.AuthMiddleware(cfg)`, but full coverage should be verified against all mounted sub-routers.
2. **Default PIN:** Initial user provisioning automatically provisions the default user with a random 6-digit PIN written to `{DATA_DIR}/persistent/genesis-pin.txt`. The admin must retrieve this PIN from the file on first startup and change it immediately via the UI.
3. **CORS Policy:** The default CORS configuration reflects incoming `Origin` headers with credentials allowed (`SameSite=None`), which should be restricted in strict multi-tenant or public deployments.
4. **Unmounted Backend Packages:** Internal backend packages such as `backup` are not imported by `main.go`. The frontend already calls some of these routes that may not be served.
5. **`.env.example` is partly decorative:** `PORT`, `DB_PATH`, `JWT_SECRET`, and `ARGON2_*` values are not all read from the environment; they must go in `config.yaml` or environment directly.

---

## 📜 Contributing & License

Contributions, bug reports, and pull requests are welcome. Please ensure that all changes adhere to the codebase conventions specified in `AGENTS.md`.

Distributed under the **MIT License**. See `LICENSE` for more information.

---