# AGENTS.md

Instructions for AI coding agents working on Kiroku (記録), a self-hosted personal finance tracker. Humans: see [README.md](./README.md).

Read this file fully before changing anything. For any UI work, also read [DESIGN.md](./DESIGN.md).

## 1. Non-negotiable rules

1. **Styling: Tailwind first, and only through the `<style>` tag.** Start every component style block with `@reference "tailwindcss";`, then write normal CSS rules whose declarations use `@apply`. See [section 3](#3-styling-rule-tailwind-via-style-tag).
2. **UI follows DESIGN.md.** Do not invent colors, radii, spacing, shadows or CSS variable names. If DESIGN.md has no answer, pick the closest existing token and say so in your summary.
3. **Never open `.env`.** Read `.env.example` only. Never print, log or commit secrets, tokens, PINs, `data/`, `*.db*` or tunnel URLs.
4. **Keep diffs small.** Do the task asked. No drive-by refactors, renames or reformatting of files you were not asked to touch.
5. **Every user-facing string goes through i18n** (`en`, `ja`, `id`). No hard-coded UI text.
6. **Verify before claiming done** (see [section 7](#7-definition-of-done)).
7. **Full modular & reusable.** Always keep code modular and reusable, at most 300 lines per file. If a file exceeds 350 lines, split it into small reusable functions or files used elsewhere. No abstractions for single use, no "for later" configuration.

## 2. Project snapshot

| Part       | Stack                                                                                     |
| ---------- | ----------------------------------------------------------------------------------------- |
| Part       | Stack                                                                                     |
| ---------- | ----------------------------------------------------------------------------------------- |
| frontend   | Vue 3 + `<script setup lang="ts">`, Vite, Pinia, Vue Router, Tailwind CSS v4, ECharts   |
| Ops        | `start.sh` generates a PM2 config; `cloudflared` tunnel; Upstash Redis publishes backend URL |

Single user, many accounts and currencies. Mobile-first PWA: the app shell is a centered column (max 480px).

### Commands

```bash
# frontend (run from repo root)
npm run dev
npm run build        # vue-tsc + vite build, this is the type-check
npm run test:run     # vitest, tests belong in tests/
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
  @apply m-auto flex max-w-md flex-col gap-4 px-4 py-6 bg-[var(--bg-card)];
}

.hero-title {
  @apply text-xl font-semibold tracking-tight text-[var(--text-primary)];
}
</style>
```

### Rules

- **Class names in the template are semantic** (`.hero`, `.stat`, `.row-amount`). Do not put utility classes (`flex`, `px-4`, `dark:bg-zinc-900`) in `class="..."`.
- **Use `@apply` for anything Tailwind can express**: layout, flex/grid, spacing, sizing, typography, radius, borders, transitions, states (`hover:`, `focus-visible:`, `disabled:`), responsive prefixes (`md:`).
- **Plain CSS is allowed only for what `@apply` cannot do**: `@keyframes`, `env(safe-area-inset-*)`, `clip-path`, complex selectors. Keep it in the same style block, next to the rule that needs it.
- **Always `<style scoped>`.** Global styles belong in `src/assets/main.css` and need a reason.
- **Never** write `@import "tailwindcss"` inside a component. `main.css` already does it once.
- **State via semantic classes:** `:class="{ active: isActive }"` in the template, `.item.active { @apply ... }` in the style block.
- **Arbitrary values are fine** when no utility exists: `@apply w-[7.5rem]`. Prefer scale values first.

### Scope caveat: what `@reference "tailwindcss"` can see

`@reference "tailwindcss"` loads **stock Tailwind only**. It does not load `src/assets/main.css`, so inside `@apply`:

- Project `@theme` tokens (`bg-primary`, `text-danger`, etc.) are **not** available. Use explicit CSS variable arbitrary-value syntax: `bg-[var(--bg-card)]`, `text-[var(--text-primary)]`, `border-[var(--border-primary)]`, `rounded-[var(--radius)]`.
- The custom `dark` variant from `main.css` is **not** available. Use Tailwind's built-in `dark:` variant directly in `@apply` for dark-mode styling.
- Stock utilities are available, including the stock palette (`zinc-*`, `emerald-*`, etc.). Prefer explicit semantic CSS variable utilities such as `bg-[var(--bg-card)]` and `text-[var(--text-primary)]` over palette classes for color.

### Dark mode

Use Tailwind's built-in `dark:` variant for dark-mode styling. Apply it directly through `@apply` inside the component's `<style scoped>` block.

```css
.hero {
  @apply bg-white text-zinc-900 dark:bg-zinc-950 dark:text-zinc-100;
}

.hero-title {
  @apply text-zinc-900 dark:text-zinc-100;
}
```

## 4. Frontend conventions

- **Component shape:** `<template>`, then `<script setup lang="ts">`, then `<style scoped>`. Use typed `defineProps` / `defineEmits`.
- **Folders:**
  - `views/` — route-level pages.
  - `components/widgets/` — self-contained, reusable UI compositions that can be dropped directly into other components. A widget must be ready to render its complete UI without requiring the parent component to manage its internal state, business logic, or implementation details. Encapsulate its own logic and dependencies internally; expose only the necessary props and events as its public interface. Widgets are **not** feature-based or domain-based containers.
  - `components/ui/` — base UI controls and primitives.
  - `components/layout/` — layout-related components.
  - `components/common/` — shared components that do not belong to a more specific category.
  - `components/charts/` — chart components.
  - `components/icons/` — icon components.
  - Each widget folder must have an `index.ts` entry point. Match the existing directory structure before creating new folders.
- **Imports:** use the `@/` alias for `src/`.
- **State:** Pinia setup stores in `src/stores/`. Components call stores; stores call the API.
- **API:** always use the shared Axios instance from `@/api/client`. The base URL is resolved at runtime (Upstash, environment variables, cache, then `http://localhost:8080`), so never hard-code a host. API paths must start with `/api/v1`.
- **i18n:** use `const { t } = useI18n()` from `@/locales`. Add every new translation key to all three locale files: `en.ts`, `ja.ts`, and `id.ts`. English (`en`) is the fallback.
- **Icons:** use `@lucide/vue` (already a dependency) for new icons. Do not introduce inline `<svg>` blocks. Emoji are reserved for user-selected category and account icons, never UI chrome.
- **Money:** format values with `Intl.NumberFormat` using the user's number-format settings. Route displayed amounts through `utils/amountDisplay` so the **hide amounts** setting remains effective. Right-align monetary values and use tabular numerals.
- **Dates:** use `dayjs`. Respect the `date_format` and `week_start` settings.
- **Theme:** read and update theme state through `useSettingsStore` / `useUIStore`. Never manipulate `data-theme` directly elsewhere.
- **Accessibility:** use semantic `<button>` and `<a>` elements, provide accessible labels for icon-only buttons, maintain visible `:focus-visible` states, ensure touch targets are at least 44px, and never rely on color alone to distinguish income from expenses.
- **PWA:** the service worker uses NetworkFirst for `/api/v1/*`. Do not introduce caching behavior that could present stale financial data as current or authoritative.

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

- Configuration uses environment variables; use `.env.example` to understand available variable names and expected values. **Never open or read `.env`**; do not assume its contents.
- Return JSON for success and `http.Error` with a short message for failures. Never leak internals or secrets in error bodies.

### Logging levels

Use the project's available logging mechanism and its severity levels. Apply levels consistently and log operations where failures are likely or consequential; use trace/debug for useful troubleshooting detail and info for important lifecycle events, without noisy or repetitive logs.

- `trace` — very detailed execution tracing; use selectively
- `debug` — diagnostic detail for troubleshooting
- `info` — important normal operation events, such as startup and completion; avoid request-by-request noise unless useful
- `warn` — unexpected but non-fatal conditions
- `error` — failures that prevent an operation from succeeding; include relevant error context
- `panic` / `fatal` — unrecoverable failures that require process termination, only when appropriate

Never log secrets, PINs, tokens, raw request bodies, or `data/` paths.

### Tests

- Frontend tests: vitest configured for `tests/**`.

## 6. Known gaps

Do not assume these are fine. Flag them to the user when your work touches them.

1. **Backend removed.** `backend/` (Go server, SQLite, migrations) has been deleted from the repo. No backend routes, auth, JWT, or tunnel remain.
2. **CORS reflects any `Origin` with credentials** — N/A after backend removal (no auth cookies remain).
3. **`.env.example` is partly decorative.** `PORT`, `DB_PATH`, `JWT_SECRET`, the `ARGON2_*` values and the 9Router/OCR keys are not read from the environment. They must go in `config.yaml` if a backend is ever re-added.
4. **CI is stale.** `.github/workflows/go.yml` references Go and `backend/`; only `build-app.yml` (Android) is active now.
5. **No frontend tests.** `vitest` is configured for `tests/**` but the folder does not exist.
6. **`start.sh` assumes Termux** for `air` (`/data/data/com.termux/files/home/go/bin/air`).
7. **Swagger spec is empty.** N/A after backend removal.

## 7. Definition of done

- Frontend change: `npm run build` passes from repo root (this runs `vue-tsc`).
- UI change: checked in **light and dark**, at 360px width, with long text (Japanese and Indonesian strings are longer or wider than English), and with the hide-amounts setting on.
- New strings exist in `en`, `ja` and `id`.
- Styling follows section 3 and uses only tokens from DESIGN.md.
- Your summary lists what changed, what you could not verify, and any known gap you hit.
