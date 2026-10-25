# AGENTS.md

Instructions for AI coding agents working on Kiroku (記録), a self-hosted personal finance tracker. Humans: see [README.md](./README.md).

Read this file fully before changing anything. For any UI work, also read [DESIGN.md](./DESIGN.md).

## 1. Non-negotiable rules

1. **Styling: Tailwind first, and only through the `<style>` tag.** Start every component style block with `@reference "tailwindcss";`, then write normal CSS rules whose declarations use `@apply`. See [section 3](#3-styling-rule-tailwind-via-style-tag).
2. **UI follows DESIGN.md.** Do not invent colors, radii, spacing, shadows or CSS variable names. If DESIGN.md has no answer, pick the closest existing token and say so in your summary.
3. **Never open `.env`.** Read `.env.example` only. Never print, log or commit secrets, tokens, PINs, `backend/data/`, `*.db*` or tunnel URLs.
4. **Keep diffs small.** Do the task asked. No drive-by refactors, renames or reformatting of files you were not asked to touch.
5. **Every user-facing string goes through i18n** (`en`, `ja`, `id`). No hard-coded UI text.
6. **Verify before claiming done** (see [section 7](#7-definition-of-done)).

## 2. Project snapshot

| Part       | Stack                                                                                     |
| ---------- | ----------------------------------------------------------------------------------------- |
| `backend/` | Go 1.27, chi v5, SQLite via `modernc.org/sqlite` (pure Go, `CGO_ENABLED=0`), zerolog, koanf |
| `frontend/`| Vue 3 + `<script setup lang="ts">`, Vite, Pinia, Vue Router, Tailwind CSS v4, ECharts, PWA |
| Ops        | `start.sh` generates a PM2 config; `cloudflared` tunnel; Upstash Redis publishes backend URL |

Single user, many accounts and currencies. Mobile-first PWA: the app shell is a centered column (max 480px).

### Commands

```bash
# backend (run from backend/)
KIROKU_TUNNEL_ENABLED=false go run ./cmd/server
go build ./... && go test ./...
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server

# frontend (run from frontend/)
npm run dev
npm run build        # vue-tsc + vite build, this is the type-check
npm run test:run     # vitest, tests belong in frontend/tests/
```

## 3. Styling rule: Tailwind via style tag

Tailwind is used **only** from the component's `<style scoped>` block, through `@apply`. The first line of the block is always `@reference "tailwindcss";`.

```vue
<template>
  <section class="hero">
    <h1 class="hero-title">{{ t('dashboard.title') }}</h1>
  </section>
</template>

<style scoped>
@reference "tailwindcss";

.hero {
  @apply m-auto flex max-w-md flex-col gap-4 px-4 py-6;
  background-color: var(--bg-card);
}

.hero-title {
  @apply text-xl font-semibold tracking-tight;
  color: var(--text-primary);
}
</style>
```

### Rules

- **Class names in the template are semantic** (`.hero`, `.stat`, `.row-amount`). Do not put utility classes (`flex`, `px-4`, `dark:bg-zinc-900`) in `class="..."`.
- **Use `@apply` for anything Tailwind can express**: layout, flex/grid, spacing, sizing, typography, radius, borders, transitions, states (`hover:`, `focus-visible:`, `disabled:`), responsive prefixes (`md:`).
- **Plain CSS is allowed only for what `@apply` cannot do**: `@keyframes`, `var(--token)` colors, `env(safe-area-inset-*)`, `clip-path`, complex selectors. Keep it in the same style block, next to the rule that needs it.
- **Always `<style scoped>`.** Global styles belong in `src/assets/main.css` and need a reason.
- **Never** write `@import "tailwindcss"` inside a component. `main.css` already does it once.
- **State via semantic classes:** `:class="{ active: isActive }"` in the template, `.item.active { @apply ... }` in the style block.
- **Arbitrary values are fine** when no utility exists: `@apply w-[7.5rem]`. Prefer scale values first.

### Scope caveat: what `@reference "tailwindcss"` can see

`@reference "tailwindcss"` loads **stock Tailwind only**. It does not load `src/assets/main.css`, so inside `@apply`:

- Project `@theme` tokens (`bg-primary`, `text-danger`, ...) are **not** available. Use `var(--token)` in a normal declaration instead, as in the example above.
- The custom `dark` variant from `main.css` is **not** available. A `dark:` utility inside `@apply` would follow the OS setting, not the in-app theme switch. Do not use it.
- Stock utilities are available, including the stock palette (`zinc-*`, `emerald-*`, ...). Prefer semantic CSS variables over palette classes for color.

### Dark mode

The theme is controlled by `data-theme="light" | "dark"` on `<html>` (no attribute means follow the OS). Semantic CSS variables in `main.css` flip automatically, so **most components need zero dark-specific code**. If a rule truly must differ, use the selector already used in the codebase:

```css
:where([data-theme="dark"], [data-theme="dark"] *) .hero {
  @apply border-zinc-800;
}
```

### CSS variables you may use

These are the only tokens defined in `main.css`. Use them by name; do not invent others.

| Role           | Variables                                                                  |
| -------------- | -------------------------------------------------------------------------- |
| Backgrounds    | `--bg-primary` `--bg-secondary` `--bg-tertiary` `--bg-card` `--bg-hover`   |
| Text           | `--text-primary` `--text-secondary` `--text-muted` `--text-inverse`        |
| Accent         | `--accent` `--accent-hover` `--accent-light`                               |
| Borders        | `--border-primary` `--border-secondary` `--border-focus`                   |
| Status         | `--success` `--warning` `--error` `--info` and their `-light` variants     |
| Shadow, radius | `--shadow-sm` `--shadow` `--shadow-md` `--shadow-lg` `--shadow-xl`, `--radius-sm` `--radius` `--radius-lg` `--radius-full` |

### Legacy styling in the repo

The codebase is mid-migration. 107 of 113 `.vue` files use hand-written CSS in `<style scoped>`, and a few (`ThemeToggle`, `ToastContainer`, `DashboardStatsWidget`) use utility classes in the template. Only `CashFlowProjectionChart.vue` uses `@reference` today.

- **New components** follow section 3 from line one.
- **Existing components:** convert to `@reference` + `@apply` only when you are already editing that component's styling. Do not mass-migrate.
- **Broken tokens:** many files use variables that are never defined (`--color-surface`, `--color-border`, `--color-text`, `--color-text-muted`, `--color-neutral-*`, `--color-surface-hover`, `--color-*-bg`, ...). They resolve to nothing. When you touch such a rule, replace them:

| Undefined (legacy)               | Use instead                         |
| -------------------------------- | ----------------------------------- |
| `--color-surface`                | `--bg-card`                         |
| `--color-background`             | `--bg-primary`                      |
| `--color-surface-hover`          | `--bg-hover`                        |
| `--color-border`                 | `--border-primary`                  |
| `--color-text`                   | `--text-primary`                    |
| `--color-text-secondary`         | `--text-secondary`                  |
| `--color-text-muted`             | `--text-muted`                      |
| `--color-danger-bg`, `--error-bg`| `--error-light`                     |
| `--color-success-bg`             | `--success-light`                   |
| `--color-warning-bg`             | `--warning-light`                   |
| `--color-primary-bg`             | `--accent-light`                    |
| `--color-neutral-*`              | pick by role from the table above   |

## 4. Frontend conventions

- **Component shape:** `<template>`, then `<script setup lang="ts">`, then `<style scoped>`. Typed `defineProps` / `defineEmits`.
- **Folders:** `views/` route pages, `widgets/<domain>/` feature blocks (each folder has an `index.ts`), `components/ui` base controls, `components/layout`, `components/common`, `components/charts`, `components/icons`. Match the neighbors before creating a new folder.
- **Imports:** use the `@/` alias for `src/`.
- **State:** Pinia setup stores in `src/stores/`. Components call stores, stores call the API.
- **API:** always use the shared axios instance from `@/api/client`. The base URL is resolved at runtime (Upstash, env, cache, then `http://localhost:8080`), so never hard-code a host. Paths start with `/api/v1`.
- **i18n:** `const { t } = useI18n()` from `@/locales`. Add every new key to **all three** files (`en.ts`, `ja.ts`, `id.ts`). `en` is the fallback.
- **Icons:** use `@lucide/vue` (already a dependency) for new icons. Do not paste new inline `<svg>` blocks. Emoji are for user-chosen category and account icons only, never UI chrome.
- **Money:** format with `Intl.NumberFormat` using the user's number format, and route displayed amounts through `utils/amountDisplay` so the **hide amounts** setting keeps working. Right-align amounts and use tabular numerals.
- **Dates:** `dayjs`. Respect the `date_format` and `week_start` settings.
- **Theme:** read and write through `useSettingsStore` / `useUIStore`. Do not toggle `data-theme` by hand elsewhere.
- **Accessibility:** real `<button>` and `<a>`, labels on icon-only buttons, visible `:focus-visible`, touch targets at least 44px, never rely on color alone for income vs expense.
- **PWA:** the service worker caches `/api/v1/*` with NetworkFirst. Do not add caching that can serve stale financial data as truth.

## 5. Backend conventions

### Layering (the live path)

```
domain/                 entities, input types, repository interfaces
repository/sqlite/      implements domain.*Repository with database/sql
delivery/http/          chi handlers, each exposes Routes() chi.Router
cmd/server/main.go      wiring: repos -> handlers -> mounts under /api/v1
```

Handlers depend on `domain` interfaces, never on `database/sql` directly. Keep new code on this path.

### Adding an endpoint

1. Add or extend types and the repository interface in `domain/`.
2. Implement it in `repository/sqlite/`.
3. Add the handler in `delivery/http/` with a `Routes()` method.
4. Mount it in `cmd/server/main.go` under `/api/v1` **and behind `auth.AuthMiddleware(cfg)`** (see [Known gaps](#known-gaps)).
5. If the schema changes, add a migration.
6. Add a frontend type in `src/types/` and a store action. Keep JSON field names `snake_case`.

### Database

- Migrations are embedded SQL in `internal/db/migrations/`, named `NNN_description.sql`, applied in filename order when `NNN` is greater than the stored `schema_version`.
- **Never edit an applied migration.** Add a new one with the next free number.
- **Do not reuse a number.** `008` already exists twice (`008_backups_users_fts.sql` and `008_receipt_jobs.sql`). The runner tracks one version per number, so on a fresh database the second file (`008_receipt_jobs.sql`) is skipped. Check this before relying on the receipt jobs table.
- SQLite runs in WAL mode with foreign keys on. Keep it pure Go: no cgo dependencies.
- Transactions are soft-deleted (`deleted_at`). Queries must filter them out unless the feature is about trash or restore.
- The `transactions.amount` column is `INTEGER` while the domain type uses `float64`. Read how `repository/sqlite/transactions.go` converts before touching any money math, and add a test for it.

### Config and logging

- `backend/config.yaml` is loaded from the working directory. Run the server from `backend/`.
- koanf maps `KIROKU_*` env vars by turning every `_` into `.`, so only single-word keys (like `KIROKU_ADDR`) override the YAML. Multi-word keys such as `jwt_secret` must be set in `config.yaml`. Telegram, Upstash, backup git token, log retention and tunnel flags have explicit `os.Getenv` fallbacks in `internal/config/config.go`.
- Use zerolog (`github.com/rs/zerolog/log`). Do not use `fmt.Println` or the std `log`.
- Return JSON for success and `http.Error` with a short message for failures. Never leak internals or secrets in error bodies.

### Tests

- Table-driven tests with testify. Run `go test ./...` from `backend/`.
- `.gitignore` currently excludes `backend/**/*_test.go`, so new tests are silently untracked. Tell the user instead of changing `.gitignore` on your own.

## 6. Known gaps

Do not assume these are fine. Flag them to the user when your work touches them.

1. **Data routes have no authentication.** In `cmd/server/main.go` only part of `/api/v1/auth/*` has `AuthMiddleware`. Under `/api/v1`, `/accounts`, `/categories`, `/transactions`, `/templates`, `/scheduled`, `/reports`, `/settings` and `/fx` are mounted without it. The top-level `/shared` mount is also open, including create, list and revoke (not only the public view). The backend is exposed through a public tunnel.
2. **Default PIN.** The first login auto-creates the user with PIN `123456`.
3. **CORS reflects any `Origin` with credentials**, and auth cookies use `SameSite=None`. Treat as unsafe until restricted.
4. **Unmounted backend packages.** `internal/{accounts,analytics,backup,categories,fx,impexp,receipts,reports,scheduled,settings,shared,templates,transactions}` are not imported by `main.go`. The frontend already calls routes such as `/api/v1/backup/*`, analytics and receipts that are not served. Do not copy patterns from these packages into new code.
5. **`.env.example` is partly decorative.** `PORT`, `DB_PATH`, `JWT_SECRET`, the `ARGON2_*` values and the 9Router/OCR keys are not read from the environment. They must go in `config.yaml`.
6. **CI is stale.** `.github/workflows/go.yml` uses Go 1.20 and runs `go build ./...` from the repo root, but `go.mod` is in `backend/` and requires Go 1.27.
7. **No frontend tests.** `vitest` is configured for `tests/**` but the folder does not exist.
8. **`start.sh` assumes Termux** for `air` (`/data/data/com.termux/files/home/go/bin/air`).
9. **Swagger spec is empty.** No handler has swag annotations.

## 7. Definition of done

- Backend change: `go build ./...` and `go test ./...` pass from `backend/`.
- Frontend change: `npm run build` passes from `frontend/` (this runs `vue-tsc`).
- UI change: checked in **light and dark**, at 360px width, with long text (Japanese and Indonesian strings are longer or wider than English), and with the hide-amounts setting on.
- New strings exist in `en`, `ja` and `id`.
- Styling follows section 3 and uses only tokens from DESIGN.md.
- Your summary lists what changed, what you could not verify, and any known gap you hit.
