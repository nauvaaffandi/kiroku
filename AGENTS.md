# AGENTS.md

Instructions for AI coding agents working on Kiroku (記録), a self-hosted personal finance tracker. Humans: see [README.md](./README.md).

Read this file fully before changing anything. For any UI work, also read [DESIGN.md](./DESIGN.md).

## 1. Non-negotiable rules

1. **Styling: Tailwind first, and only through the `<style>` tag.** Start every component style block with `@reference "tailwindcss";`, then write normal CSS rules whose declarations use `@apply`. See [section 3](#3-styling-rule-tailwind-via-style-tag).
2. **UI follows DESIGN.md.** Do not invent colors, radii, spacing, shadows or CSS variable names. If DESIGN.md has no answer, pick the closest existing token and say so in your summary.
3. **Never open `.env`.** Read `.env.example` only. Never print, log or commit secrets, tokens, PINs, `data/`, `*.db*` or tunnel URLs.
4. **Keep diffs small.** Do the task asked. No drive-by refactors, renames or reformatting of files you were not asked to touch.
5. **Every user-facing string goes through i18n** (`en`, `ja`, `id`). No hard-coded UI text.
6. **Verify before claiming done** (see [section 6](#6-definition-of-done)).
7. **Full modular & reusable.** Always keep code modular and reusable, at most 300 lines per file. If a file exceeds 350 lines, split it into small reusable functions or files used elsewhere. No abstractions for single use, no "for later" configuration.

## 2. Project snapshot

| Part       | Stack                                                                                     |
| ---------- | ----------------------------------------------------------------------------------------- |
| frontend   | Vue 3 + `<script setup lang="ts">`, Vite, Pinia, Vue Router, Tailwind CSS v4, ECharts   |
| storage    | Capacitor + `@capacitor-community/sqlite`, local SQLite (`src/storage/repository.ts`), offline-first |

Single user, many accounts and currencies. Mobile-first PWA / Capacitor Android app: the app shell is a centered column (max 480px). There is **no backend** — data lives in local SQLite only.

### Commands

```bash
# frontend (run from repo root)
npm run dev
npm run build        # vue-tsc + vite build, this is the type-check
npm run test:run     # vitest, tests belong in tests/ (folder does not exist yet)
```

Android CI: `.github/workflows/build-app.yml` (Node 22, Java 21, Android SDK, `npm ci`, `npm run build`, copy dist to android assets, `./gradlew assembleDebug`, upload APK).

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
  - `components/ui/` — base UI controls and primitives. Has `index.ts` barrel.
  - `components/layout/` — layout-related components.
  - `components/common/` — shared components that do not belong to a more specific category.
  - `components/charts/` — chart components.
  - `components/icons/` — icon components.
  - `components/auth/` — auth modals (PIN, fingerprint, Telegram).
  - `components/forms/` — form controls (FAB, TransactionForm, etc.).
  - `components/reports/` — report panels.
  - `components/settings/` — settings UI.
  - `components/shared/` — shared access modals.
  - `src/widgets/` — self-contained, reusable UI compositions (NOT `components/widgets/`). Each widget folder has an `index.ts` entry point. Barrel at `src/widgets/index.ts`. Widgets are **not** feature-based or domain-based containers.
  - `src/stores/` — Pinia setup stores (14 stores). Components call stores; stores call the API or local storage.
  - `src/storage/` — local SQLite via Capacitor (`LocalDatabase`, `localDb`). Offline-first data layer.
  - `src/composables/` — reusable composables + `index.ts` barrel.
  - `src/types/` — TypeScript types + `index.ts` + `pwa.d.ts`.
  - `src/utils/` — `amountDisplay.ts`, `categoryName.ts`, `exportExcel.ts`.
  - `src/api/` — `client.ts` (shared Axios), `notifications.ts`, `push.ts`, `receipts.ts`.
  - `src/locales/` — `en.ts`, `ja.ts`, `id.ts`, `index.ts` (custom Pinia-based i18n store, `useI18n()`).
- **Imports:** use the `@/` alias for `src/`.
- **State:** Pinia stores in `src/stores/`. Components call stores; stores call the API or `src/storage/repository.ts`.
- **API:** always use the shared Axios instance from `@/api/client`. **Offline app:** `baseURL` is `''`, `initApiClient()` and `refreshBackendUrl()` both return `''`. There is **no remote backend** — do not hard-code a host, do not expect `/api/v1` endpoints. Data writes/reads go through Pinia stores → `src/storage/repository.ts` (local SQLite).
- **i18n:** use `const { t } = useI18n()` from `@/locales`. Add every new translation key to all three locale files: `en.ts`, `ja.ts`, and `id.ts`. English (`en`) is the fallback.
- **Icons:** use `@lucide/vue` (already a dependency) for new icons. Do not introduce inline `<svg>` blocks. Emoji are reserved for user-selected category and account icons, never UI chrome.
- **Money:** format values with `Intl.NumberFormat` using the user's number-format settings. Route displayed amounts through `utils/amountDisplay` so the **hide amounts** setting remains effective. Right-align monetary values and use tabular numerals.
- **Dates:** use `dayjs`. Respect the `date_format` and `week_start` settings.
- **Theme:** read and update theme state through `useSettingsStore` / `useUIStore`. Never manipulate `data-theme` directly elsewhere.
- **Storage:** entities are soft-deleted (`deleted = 1`). Queries must filter `deleted = 0` unless the feature is about trash or restore.
- **Accessibility:** use semantic `<button>` and `<a>` elements, provide accessible labels for icon-only buttons, maintain visible `:focus-visible` states, ensure touch targets are at least 44px, and never rely on color alone to distinguish income from expenses.
- **Logging:** Full frontend, no complex logging level. Always use `console.log()` for all data flow operations. Every **get**, **save**, **update**, and **delete** process MUST have `console.log()` to trace input payload, response data, and errors. Never log secrets, PINs, or raw auth tokens.

## 5. Known gaps

Do not assume these are fine. Flag them to the user when your work touches them.

1. **Backend removed.** `backend/` (Go server, SQLite, migrations) has been deleted. No backend routes, auth, JWT, or tunnel remain. The app is offline-first with local SQLite only.
2. **No `start.sh`.** Previously referenced `start.sh` (PM2 + cloudflared); it does not exist. Do not reference it.
3. **No `go.yml` CI.** `.github/workflows/go.yml` does not exist. Only `build-app.yml` (Android APK) is active.
4. **No frontend tests.** `vitest` is configured for `tests/**` but the folder does not exist.
5. **`.env.example` is partly decorative.** It contains frontend vars (`VITE_BACKEND_URL`, `VITE_VAPID_PUBLIC_KEY`, etc.) plus commented-out backend reference vars (PORT, JWT_SECRET, ARGON2_*, 9Router, OCR) that are not read from the environment. No `config.yaml` exists.
6. **Offline-only API client.** `@/api/client` returns `''` for base URL — no remote backend. Do not write code assuming `/api/v1` endpoints exist.

## 6. Definition of done

- Frontend change: `npm run build` passes from repo root (this runs `vue-tsc`).
- UI change: checked in **light and dark**, at 360px width, with long text (Japanese and Indonesian strings are longer or wider than English), and with the hide-amounts setting on.
- New strings exist in `en`, `ja` and `id`.
- Styling follows section 3 and uses only tokens from DESIGN.md.
- Your summary lists what changed, what you could not verify, and any known gap you hit.
