# Kiroku Enterprise Testing Architecture & Master Plan

## 1. Executive Summary & Architecture Vision

Kiroku is a self-hosted personal finance tracker built as a Go 1.27 backend (`chi` v5, SQLite `modernc.org/sqlite`) and a Vue 3 frontend (Vite, Pinia, Tailwind CSS v4, PWA). This document defines the enterprise-grade End-to-End (E2E), integration, and unit testing suite designed to ensure 100% reliability, zero regressions, and robust cross-platform mobile verification.

### Core Testing Pillars
1. **Backend Integration Testing**: Pure Go `net/http/httptest` coupled with ephemeral SQLite in-memory databases, automated schema migration application, and transaction rollback isolation.
2. **Frontend Component & Store Testing**: Vitest running under Bun with `@vue/test-utils`, Pinia store hydration checks, and DOM assertions against Vue components.
3. **End-to-End (E2E) Browser Flows**: Playwright running against compiled Go binaries and Vite dev/preview servers, validating full user journeys inside mobile viewports (380px–480px).
4. **Data Isolation Strategy**: Dedicated per-test isolated SQLite instances with deterministic fixture seeding via JSON snapshots.
5. **CI/CD Integration**: GitHub Actions pipeline executing Go test suites, Bun Vitest runs, and Playwright E2E suites with artifact collection.

---

## 2. Directory Structure (`./tests/`)

```
tests/
├── README.md                          # Testing guidelines & execution commands
├── backend/                           # Go integration and API test suites
│   ├── go.mod                         # Test-specific or module integration
│   ├── internal/                      # Test helpers, database setup, assertions
│   │   ├── db/                        # Test SQLite memory hooks & migrators
│   │   ├── fixtures/                  # JSON/SQL test seed fixtures
│   │   └── server/                    # Test HTTP server harness initializer
│   ├── auth_test.go                   # Authentication & PIN verification tests
│   ├── accounts_test.go               # Account CRUD, currency, balances
│   ├── categories_test.go             # Category hierarchy & icons
│   ├── transactions_test.go           # Ledger entry, splits, amount math
│   ├── templates_test.go              # Quick transaction templates
│   ├── scheduled_test.go              # Recurring transaction jobs
│   ├── reports_test.go                # Cashflow, summary, analytics endpoints
│   ├── fx_test.go                     # Currency exchange rates & conversion
│   ├── settings_test.go               # User preferences & config sync
│   ├── shared_test.go                 # Shared public token views
│   ├── receipts_test.go               # OCR receipt parsing jobs
│   └── backup_test.go                 # SQLite export/import & Git backup
├── frontend/                          # Bun + Vitest unit & component specs
│   ├── package.json                   # Bun test scripts & dependencies
│   ├── tsconfig.json                  # TypeScript config for test environment
│   ├── setup.ts                       # Vitest global setup & mocks (localStorage, matchMedia)
│   ├── stores/                        # Pinia store unit tests
│   │   ├── auth.spec.ts               # Auth state, PIN token persistence
│   │   ├── accounts.spec.ts           # Account fetching & filtering
│   │   ├── transactions.spec.ts       # Ledger state & pagination
│   │   └── settings.spec.ts           # Theme, currency, amount hiding
│   ├── components/                    # Vue component specs
│   │   ├── Navigation.spec.ts         # Mobile tab bar & routing
│   │   ├── AmountDisplay.spec.ts      # Privacy toggle & currency formatting
│   │   └── ThemeToggle.spec.ts        # Dark/light attribute toggling
│   └── views/                         # View-level integration specs
│       ├── Dashboard.spec.ts          # Widget rendering & summary cards
│       └── TransactionsView.spec.ts   # Ledger filtering & modal flows
├── e2e/                               # Playwright end-to-end browser specs
│   ├── playwright.config.ts           # Playwright runner configuration
│   ├── fixtures/                      # Browser auth state & seed states
│   ├── helpers/                       # Locator extensions & API seeding helpers
│   ├── auth.spec.ts                   # Login, PIN setup, session expiry
│   ├── ledger.spec.ts                 # Adding income/expense on mobile viewports
│   ├── settings.spec.ts               # Dark mode, amount hiding toggle
│   └── receipts.spec.ts               # Receipt upload and OCR processing flow
└── fixtures/                          # Shared cross-suite JSON fixtures
    ├── users.json
    ├── accounts.json
    ├── categories.json
    └── transactions.json
```

---

## 3. Bun Runner Configuration (`package.json`)

The frontend testing suite utilizes Bun as the test runner combined with Vitest for blazing-fast execution.

### `tests/frontend/package.json`
```json
{
  "name": "kiroku-tests-frontend",
  "private": true,
  "type": "module",
  "scripts": {
    "test": "vitest run",
    "test:watch": "vitest",
    "test:coverage": "vitest run --coverage"
  },
  "devDependencies": {
    "@pinia/testing": "^0.1.4",
    "@vue/test-utils": "^2.4.6",
    "jsdom": "^24.0.0",
    "typescript": "^5.3.3",
    "vitest": "^1.2.0",
    "vue": "^3.4.15",
    "vue-router": "^4.2.5"
  }
}
```

### `tests/frontend/tsconfig.json`
```json
{
  "compilerOptions": {
    "target": "ESNext",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "jsx": "preserve",
    "sourceMap": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "esModuleInterop": true,
    "lib": ["ESNext", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "paths": {
      "@/*": ["../../frontend/src/*"]
    }
  },
  "include": ["./**/*.ts", "./**/*.tsx", "./**/*.vue"]
}
```

---

## 4. Backend Endpoint Test Matrix & Use Cases

Every backend endpoint implemented in Go must be tested with positive cases, boundary conditions, malformed input handling, and authorization checks.

| Subsystem | Endpoint Path | Method | Test Scenarios & Use Cases |
| :--- | :--- | :--- | :--- |
| **Auth** | `/api/v1/auth/login` | POST | • Valid PIN returns JWT & sets cookie<br>• Invalid PIN returns 401 Unauthorized<br>• Rate-limiting after 5 consecutive failures |
| **Auth** | `/api/v1/auth/verify` | GET | • Active session cookie returns 200 OK<br>• Missing/expired token returns 401 |
| **Accounts** | `/api/v1/accounts` | GET / POST | • List accounts ordered by sort index<br>• Create account with valid balance/currency<br>• Reject duplicate account names or invalid currency codes |
| **Accounts** | `/api/v1/accounts/{id}` | PUT / DELETE | • Update balance/name successfully<br>• Soft delete account; ensure dependent queries omit deleted items |
| **Categories** | `/api/v1/categories` | GET / POST | • Fetch tree structure (parent/child categories)<br>• Create subcategory with valid parent ID<br>• Prevent circular parent-child references |
| **Transactions**| `/api/v1/transactions` | GET / POST | • Create income, expense, and transfer records<br>• Verify integer amount storage vs float representation<br>• Filter by date range, category, account, and search query |
| **Templates** | `/api/v1/templates` | GET / POST / DELETE | • Store recurring/frequent transaction template<br>• Instantiate template into actual ledger transaction |
| **Scheduled** | `/api/v1/scheduled` | GET / POST | • Create recurring rule (daily, weekly, monthly)<br>• Verify automatic generation trigger logic |
| **Reports** | `/api/v1/reports/cashflow` | GET | • Calculate net cash flow over custom date intervals<br>• Aggregate income vs expenses accurately |
| **FX** | `/api/v1/fx/rates` | GET / POST | • Retrieve cached exchange rates<br>• Convert amounts across multi-currency accounts |
| **Settings** | `/api/v1/settings` | GET / PUT | • Persist user preferences (theme, date format, week start)<br>• Reject invalid timezone or format strings |
| **Shared** | `/api/v1/shared` | POST / GET | • Generate secure public token link for filtered view<br>• Revoke sharing token |
| **Receipts** | `/api/v1/receipts` | POST / GET | • Upload receipt image blob<br>• Trigger OCR parser mock and verify extracted metadata |
| **Backups** | `/api/v1/backup/export` | GET | • Export encrypted or raw SQLite backup stream<br>• Import and restore database snapshot |

---

## 5. Frontend E2E & Component Testing Specifications

### 5.1 Pinia Store Tests
- **Auth Store**: Verify PIN entry state transition, cookie/localStorage persistence, and automatic logout on 401 responses.
- **Accounts Store**: Validate balance summation across multi-currency accounts using exchange rate multipliers.
- **Settings Store**: Test theme switching (`data-theme="dark"` / `"light"`), amount hiding flag toggling, and locale selection (`en`, `ja`, `id`).

### 5.2 Component & View Tests
- **AmountDisplay.vue**: Verify that when amount hiding is active, numbers render as `••••••` while maintaining screen reader accessibility labels.
- **Navigation.vue**: Test mobile tab bar routing, active state highlights, and touch target size compliance (>= 44px).
- **DashboardView.vue**: Validate widget grid rendering, chart mount lifecycle with ECharts mocks, and empty state rendering when zero transactions exist.

---

## 6. Database Test Isolation Strategy

1. **In-Memory SQLite Instance**: Each Go test package initializes an isolated connection string (`file::memory:?cache=shared`).
2. **Migration Runner**: Prior to running test suites, the test helper executes embedded migrations from `internal/db/migrations/` sequentially.
3. **Fixture Seeding**: Deterministic JSON seed files (`tests/fixtures/*.json`) are loaded via helper functions into the ephemeral database before test execution begins.
4. **Transaction Rollback / Clean Table Truncation**: Tests execute within isolated transactions or truncate tables between test cases to prevent state pollution.

---

## 7. CI/CD Pipeline Integration (GitHub Actions)

### `.github/workflows/test.yml`
```yaml
name: Kiroku Comprehensive Test Suite

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  backend-test:
    name: Backend Integration Tests (Go)
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.27'

      - name: Run Backend Tests
        run: |
          cd backend
          CGO_ENABLED=0 go test -v -race -coverprofile=coverage.out ./...

  frontend-test:
    name: Frontend Unit & Component Tests (Bun)
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Setup Bun
        uses: oven-sh/setup-bun@v1
        with:
          bun-version: latest

      - name: Install Frontend Dependencies
        run: |
          cd frontend
          bun install

      - name: Run Vitest Suite
        run: |
          cd frontend
          bun run test

  e2e-test:
    name: End-to-End Browser Tests (Playwright)
    needs: [backend-test, frontend-test]
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Setup Go & Bun
        uses: actions/setup-go@v5
        with:
          go-version: '1.27'
      - uses: oven-sh/setup-bun@v1

      - name: Install Playwright Browsers
        run: |
          cd tests/e2e
          bun install
          bunx playwright install --with-deps

      - name: Run E2E Suite
        run: |
          cd tests/e2e
          bun test
```
