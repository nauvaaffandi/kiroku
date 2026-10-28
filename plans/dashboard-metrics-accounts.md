# Dashboard: Metrics Highlight + Compact Account List

## Status

- Type: frontend UI/UX implementation plan
- Status: ready for implementation
- Scope: dashboard’s key financial metrics section and compact per-account overview
- Explicitly excluded: FAB changes, backend/API changes, changes to calculations, changes to account management flows
- Primary screen: mobile-first dashboard, app shell content column capped at 480px

## 1. Request and product intent

The dashboard’s four financial figures—monthly income, monthly expense, monthly net, and total balance—are important and should be visually prominent as one coherent section, not four weak standalone values. The dashboard should also show a compact list of accounts so the user can quickly inspect account-specific information without navigating away.

The account section is a dashboard overview, not a replacement for the Accounts screen. It should provide high information density while remaining easy to scan, visually quiet, and subordinate to the four key metrics.

The FAB is explicitly out of scope and must remain untouched.

## 2. Goals

1. Present all four financial metrics together in a clear, emphasized section.
2. Preserve existing metric values and calculation behavior exactly.
3. Show a compact row for each account using only account data already loaded by the dashboard.
4. Make account name, account type, currency, and balance scannable at a glance.
5. Keep balances readable, correctly formatted, right-aligned, tabular, and compatible with the hide-amounts setting.
6. Follow the project’s Kiroku design tokens in both light and dark themes.
7. Keep the dashboard mobile-first and usable at 360px width, including Japanese and Indonesian labels.
8. Add no dependency, API endpoint, database field, speculative feature, or new business logic.
9. Keep the implementation reviewable and limited to the dashboard view, dashboard metric widget, and locale files needed by new strings.

## 3. Non-goals

- Do not change the FAB or its placement, interactions, labels, animation, or styles.
- Do not change account create/edit/archive/delete flows.
- Do not add account filtering, sorting controls, pagination controls, search, or account actions to the dashboard.
- Do not add new API requests. Reuse the `accountsStore.accounts` data already fetched by the dashboard.
- Do not derive account balances from transactions or reconcile account balances.
- Do not change monthly income, monthly expense, net, total-balance, or recent-transaction calculations.
- Do not add a new global design token or modify `main.css` / `DESIGN.md` as part of implementation.
- Do not add account-type translations unless the existing account-type translation set is inadequate and implementation confirms a need. Search and reuse existing translations first.
- Do not add icons or emoji as decorative UI chrome. Account icon may be shown only if the existing UI conventions and the account’s user-selected icon field justify it; keep it small and optional.
- Do not change chart data, chart calculations, chart routing, or chart lifecycle.
- Do not rewrite unrelated dashboard widgets or the FAB.

## 4. Repository constraints to follow

Read and apply `AGENTS.md` and `DESIGN.md` before implementation.

### Styling

- Component styles must use `<style scoped>` and start with `@reference "tailwindcss";`.
- Template classes must be semantic names, not Tailwind utilities.
- Use Tailwind through `@apply` for supported layout, spacing, sizing, typography, borders, and states.
- Use CSS variables for color and existing design tokens only. Do not add hex colors or invent variable names.
- `@reference "tailwindcss"` loads stock Tailwind only; project-specific typography utilities such as `text-heading`, `text-caption`, `text-body`, or `text-label` may not resolve. Use supported stock utilities via `@apply` and ordinary CSS declarations for exact token typography values when needed.
- Avoid OS-driven `dark:` styling: the app theme is controlled by `data-theme`. Prefer semantic CSS variables that switch automatically.
- Do not use utility classes in template `class` attributes.
- Card surfaces follow the design system: `--bg-card`, 1px `--border-primary`, `--radius-lg`, no resting shadow. Dashboard emphasis must come from hierarchy, layout, and restrained surface contrast—not unauthorized color or decorative effects.

### Internationalization

- Every new user-visible string must use `t(...)`.
- Add new locale keys to `frontend/src/locales/en.ts`, `ja.ts`, and `id.ts`.
- Keep translations concise enough for 360px and verify the longest label does not collide with the value.

### Financial values

- Displayed amounts must use `utils/amountDisplay` so hide-amounts works consistently.
- Respect settings for number format, base currency, and currency display. Avoid a hard-coded `ja-JP` formatter for the new account values.
- Right-align changing amounts and use tabular numerals.
- Distinguish income from expense with sign and text/color treatment; never rely on color alone.
- Do not reinterpret account balances as income/expense. They are neutral balances.
- Multi-currency balances cannot be meaningfully summed without conversion. This plan does not introduce any new aggregate or exchange-rate calculation.

### Scope safety

- Never read `.env`; do not touch backend data, database files, credentials, tunnel URLs, or secret values.
- Do not commit changes.
- Preserve unrelated worktree changes. The current worktree already contains untracked backend label files; do not modify or stage them.

## 5. Current implementation facts

### Dashboard composition

`frontend/src/views/misc/DashboardView.vue` currently composes:

1. `DashboardStatsWidget`, receiving `monthlyIncome`, `monthlyExpense`, `monthlyNet`, `totalBalance`.
2. Spending-by-category and monthly-trend chart panels.
3. `RecentTransactionsWidget`, receiving recent transactions, categories, and accounts.

The view already creates `useAccountsStore()` and calls `accountsStore.fetchAccounts()` in its existing `onMounted` data-fetch sequence. The account section should consume the existing `accountsStore.accounts` array; do not introduce a duplicate fetch.

### Existing metric calculations

- `monthlyIncome`: sum of current-month transaction amounts whose `type` is `income`.
- `monthlyExpense`: sum of current-month transaction amounts whose `type` is `expense`.
- `monthlyNet`: monthly income minus monthly expense.
- `totalBalance`: current reduction of all loaded account balances.

These expressions and their source data must remain unchanged in this task. If the current aggregate has mixed-currency semantics, do not silently alter it here; note it as a separate product/data issue rather than extending scope.

### Account model

`frontend/src/types/account.ts` defines `Account` with these fields:

- `id: string`
- `name: string`
- `type: AccountType`
- `currency: string`
- `balance: number`
- `icon: string`
- `color: string`
- `sort_order: number`
- `archived: boolean`

`AccountType` is currently `'checking' | 'savings' | 'credit' | 'cash' | 'investment' | 'other'`.

Do not invent account fields such as institution, masked account number, available credit, transaction count, or last activity.

### Settings/utilities

- `frontend/src/stores/settings.ts` exposes settings, including `numberFormat`, `baseCurrency`, and `isHideAmounts`.
- `frontend/src/utils/amountDisplay.ts` exposes currency-symbol lookup and `formatAmountForDisplay(...)`. Its current formatter uses the settings store’s base currency, so inspect this utility carefully before using it for accounts denominated in different currencies.
- Do not show a foreign-currency account balance with the base currency’s symbol. For each account, its own `currency` field is authoritative for its balance display.
- Prefer the existing amount display utility if it can represent the account’s currency correctly without broadening utility behavior. If it cannot, choose the smallest safe approach that honors hide-amounts and account currency, or record the mismatch and seek approval before changing shared financial formatting behavior.

## 6. Proposed user experience

### 6.1 Page hierarchy

Order the dashboard content as follows unless the existing product flow or app shell imposes a conflicting order:

1. **Key metrics section** — primary dashboard summary; visually prominent and cohesive.
2. **Accounts section** — compact account-by-account overview; secondary but immediately scannable.
3. Existing charts and recent transactions, retaining their current order relative to each other unless implementation reveals a specific layout constraint.

Do not duplicate account totals inside the account list header unless the existing metric’s total balance is clearly not an aggregate. The top metrics section already shows total balance; the account section’s header should identify the list and optionally provide a route to the full Accounts screen only if a suitable existing route/string is confirmed.

### 6.2 Key metrics section

- Keep all four existing values and props. Do not change calculations or data dependencies.
- Group the four metrics within a single cohesive section/card.
- Create hierarchy: a concise section heading, then four clearly labeled amounts.
- Use a mobile two-column arrangement so all four values are visible without excessive vertical space; at 360px ensure columns and large amounts do not collide.
- Use paired rows/columns with subtle hairline separators, not four separate cards.
- Use `--bg-card`, `--border-primary`, and `--radius-lg` for the containing surface.
- Use clear labels in sentence case, `--text-muted`/`--text-secondary`, and amounts at an appropriately strong scale (design system `amount-lg` 24px or a carefully constrained equivalent).
- Income may use `--success`, with `+` sign. Expenses should remain ink-colored per DESIGN.md; use `−` as a sign and do not turn routine spending into a red error state. Net may use income color only when positive and expense color only when negative, always with a sign; zero should remain neutral or use the documented positive convention only if existing behavior requires it. Total balance remains neutral.
- Hide-amounts must mask each value while preserving its currency and sign behavior according to the existing utility contract. Do not hand-roll an inconsistent masking format.
- Ensure CJK wrapping and long localized labels. Labels may wrap to two lines; values must remain right readable, without clipping.
- Do not use icon tiles, gradients, glow, colored shadows, all-caps labels, font weight over 600, or a card shadow.

### 6.3 Compact account list section

- Present accounts as a single compact section/card, not one card per account.
- Use a concise localized heading (for example `dashboard.accounts`) and optionally a localized “view all” link to the existing Accounts route after verifying that route. Avoid adding a new route.
- Render one row per account from `accountsStore.accounts`.
- Each row should show:
  - account name as primary text;
  - localized account type as secondary metadata, provided a translation exists or is added in all locales;
  - account-specific currency plus formatted balance, right-aligned and tabular;
  - optional small user-selected account icon/color marker only if it is already used consistently and does not create a large colored tile. The 8px dot convention is preferred where a marker is useful.
- Account rows should be dense but accessible: align name/type and amount; use at least 44px row height/tap target only if the row is interactive. If rows are not links, keep a reasonable compact visual height and do not make them look clickable.
- Use 12px/16px spacing units, a hairline divider, `--text-primary`, `--text-muted`, `--border-primary`, `--bg-card`, and `--radius-lg` only as appropriate.
- No account actions, selection controls, or account-specific navigation unless there is an existing, verified account detail destination and it fits the request. The simplest version is read-only.
- Archived account behavior must be deliberate. The store currently has `filters.archived` initially empty and the backend response behavior must not be assumed. Inspect the actual fetched account data and existing Accounts screen conventions. Do not add a new archive filter or silently hide records. If archived accounts are included, display a localized archived marker only if that distinction is already relevant and translated; otherwise preserve the existing list data behavior.
- Empty state: if there are no accounts, show a short localized sentence and optionally a link/action only if an existing, safe account-creation route is verified. No decorative illustration or emoji.
- Long account names and labels must truncate gracefully without pushing balances off-screen. Use `min-width: 0`, flexible name column, and a non-shrinking amount column with responsive type sizing if needed.
- For many accounts, do not invent a maximum count or pagination. Render the available store list; if dashboard height becomes problematic, propose a “show more” behavior as a separate decision rather than hiding entries silently.

## 7. Files expected to change

Expected minimum:

1. `frontend/src/widgets/dashboard/DashboardStatsWidget.vue`
   - Refine current metrics into one clear, highlighted section.
   - Keep its existing props and numeric formatting/calculation responsibilities unchanged unless the user explicitly approves a narrowly scoped format fix.
2. `frontend/src/views/misc/DashboardView.vue`
   - Add/import the compact account section component or inline section, passing `accountsStore.accounts`.
   - Do not modify fetch sequence, computed calculations, chart logic, or FAB.
3. `frontend/src/locales/en.ts`, `frontend/src/locales/ja.ts`, `frontend/src/locales/id.ts`
   - Add strings only for new section title, account type names, empty state, and any “view all” link that implementation actually uses.
4. Optional `frontend/src/widgets/dashboard/AccountsOverviewWidget.vue`
   - Create only if a separate component materially improves clarity/reuse. It is reasonable because the section has its own list rendering and styles. Keep it narrowly scoped; do not create an `index.ts` or abstraction without a project convention requiring it.

Do not edit `DashboardView.vue` script logic beyond imports/props needed to render the section. Do not edit FAB code/files.

## 8. Implementation sequence

### Step 1 — Confirm live conventions and destination

1. Read `AGENTS.md` and `DESIGN.md` in full before UI changes.
2. Inspect `DashboardView.vue`, current `DashboardStatsWidget.vue`, account types/store, Accounts view/widget, account type translations, and account detail routes.
3. Inspect current `amountDisplay` behavior and callers. Determine if it supports arbitrary account currency and current hide-amounts preference.
4. Inspect `frontend/src/assets/main.css` variable definitions to confirm actual token values, but do not edit them for this task.
5. Inspect current worktree status. Preserve all unrelated modifications and untracked files.
6. Record account sort/order and archived data behavior based on existing code; do not infer undocumented API behavior.

Verification: each planned field and route must be evidenced by source code; no fabricated information.

### Step 2 — Lock scope and visual structure

1. Keep metric props and existing computed values intact.
2. Define a single metrics card containing four values, using a compact two-column layout.
3. Define one accounts card containing heading and rows, with the account name/type on the left and currency/balance on the right.
4. Confirm section order relative to charts and recent transactions.
5. Confirm which user-facing strings require new locale keys.
6. Keep FAB explicitly outside touched files/components.

Verification: implementation plan is limited to the files listed in section 7; no backend/API/data changes.

### Step 3 — Implement metrics presentation

1. Update only markup/classes/styles in `DashboardStatsWidget.vue`.
2. Keep the component’s props identical.
3. Add one accessible section heading using an existing dashboard heading key if appropriate; otherwise add the new key to all three locales.
4. Apply the status treatment and typography rules from section 6.2.
5. Use supported Tailwind utilities through `@apply`; use CSS variables for tokens.
6. Remove unused classes or attributes introduced by this UI change only.
7. Do not change calculations or add icon dependencies.

Verification: typecheck/build; inspect the compiled styling for unsupported `@apply` utilities.

### Step 4 — Implement compact account overview

1. Prefer a focused `AccountsOverviewWidget.vue` if it keeps `DashboardView.vue` clean. If adding the file, inspect neighboring widget conventions first.
2. Type props as `Account[]` from `@/types`; no loose `any` types.
3. Use the existing accounts array; do not fetch again.
4. Render account name, localized type, account currency, and balance only.
5. Format each balance with account currency; honor `isHideAmounts` through the established utility if it supports this currency contract. If unsupported, do not silently use the base currency’s symbol—pause or use the smallest safe local formatter that preserves masking and explicit account currency without changing global behavior.
6. Use `t(...)` for heading, account types, empty state, archived marker, and route link if applicable.
7. Keep name and amount columns responsive; truncate only the name/type text, never the amount or currency.
8. Add separators and compact spacing; do not create cards for individual accounts.
9. Add section to DashboardView after key metrics and before charts unless verified layout evidence favors another placement.
10. Reuse verified account route for “view all” only if it exists and is consistent; otherwise omit the link rather than invent a route.

Verification: no missing locale keys; no undefined account fields; no duplicate API request; hide-amounts behavior checked.

### Step 5 — Add/confirm translations

1. Search for existing keys for account types, account list, empty states, and “view all”. Reuse existing keys if they fit exact meaning.
2. Add only missing keys to `en.ts`, `ja.ts`, and `id.ts` with equivalent concise wording.
3. Avoid changing unrelated locale keys or copy.
4. Ensure `en` fallback does not conceal missing Japanese/Indonesian entries; verify all three explicitly.

Verification: locale object keys remain type-compatible and all new template keys exist in all locales.

### Step 6 — Build and inspect diff

1. Run `npm run build` in `frontend/` (runs `vue-tsc` and Vite build).
2. Review `git diff --check` and `git diff` for only requested UI and locale changes.
3. Confirm `DashboardView.vue` metric computed values, fetch list, chart initialization, resize listener, and lifecycle cleanup are unchanged.
4. Confirm FAB source is untouched and FAB remains in existing app shell.
5. Check for unsupported Tailwind utilities: `@reference` only sees stock Tailwind, not project typography names.
6. Check for hardcoded hex colors, undefined CSS variables, utility classes in template attributes, font weights above 600, and shadows on metric/account cards.
7. Check `formatAmountForDisplay` usages and ensure currency label does not falsely imply an amount is shown in base currency when it is an account currency.

Verification: build succeeds; diff scope is deliberate; no unintended behavior change.

### Step 7 — Visual QA

Use the available browser workflow when a runnable app/browser is available. Do not claim visual verification if it was not performed.

Test the dashboard at:

- 360px viewport width;
- light theme;
- dark theme via app’s `data-theme` setting;
- English, Japanese, and Indonesian locales;
- hide-amounts off and on;
- no accounts, one account, and multiple accounts;
- long account name;
- large account balance;
- accounts with different currencies;
- no account icon/color and configured icon/color;
- archived account presence if the API/store returns archived entries.

Check:

- metric section visually reads as the main financial summary;
- all four metrics remain clearly distinguished and legible;
- no clipping, horizontal overflow, or collision at 360px;
- labels and currency strings accommodate Japanese/Indonesian length;
- account balances remain right-aligned and never wrap into the name column;
- list rows remain compact and scannable;
- income/expense cues include signs and are not color-only;
- hide-amounts masks all displayed balances in this new section;
- FAB remains unchanged and unobstructed.

If browser verification is unavailable, state that visual QA remains unverified; do not present build success as proof of appearance.

## 9. Acceptance criteria

### Functional

- [ ] Metrics section shows monthly income, monthly expense, monthly net, and total balance.
- [ ] Existing metric values and calculations are unchanged.
- [ ] Account list is sourced from the already fetched `accountsStore.accounts` array.
- [ ] Each account row shows account name, localized account type (if implemented), its currency, and its balance.
- [ ] No account data is fabricated or fetched through a new endpoint.
- [ ] Hide-amounts setting masks account balances consistently.
- [ ] Account currencies are not mislabeled with the base-currency symbol.
- [ ] Empty account state is localized and visually calm.
- [ ] Account list does not add unrelated actions or interaction.
- [ ] FAB is untouched.

### Visual/design

- [ ] All four metrics appear as one cohesive, emphasized section.
- [ ] The metrics have a clear hierarchy and remain compact on mobile.
- [ ] Account list is a compact single card/section with row dividers, not individual cards.
- [ ] Light and dark themes use semantic CSS variables and remain legible.
- [ ] Uses only DESIGN.md tokens; no new color/radius/spacing/shadow tokens.
- [ ] No utility classes in templates; component styles begin with `@reference "tailwindcss";`.
- [ ] Accessible heading hierarchy, semantic list markup, readable text, and visible focus for any links.
- [ ] 360px layout has no horizontal scrolling, overlap, or clipped amounts.
- [ ] Japanese and Indonesian labels fit without breaking amount alignment.

### Verification

- [ ] All new strings exist in `en.ts`, `ja.ts`, and `id.ts`.
- [ ] `npm run build` passes from `frontend/`.
- [ ] `git diff --check` passes.
- [ ] Diff contains no unintended changes to business logic, FAB, backend, API, or generated build output.
- [ ] Visual QA is completed or explicitly reported as unavailable.

## 10. Risks and mitigations

### Mixed currencies in total balance

The existing total-balance computation adds account balances without currency conversion. This task must preserve that existing behavior but must not imply that mixed-currency balances are converted. If the current metric is misleading for multi-currency accounts, raise a separate product decision; do not silently change the value in this UI task.

### Hide-amounts and account currency

`formatAmountForDisplay` currently uses the settings store’s base currency. A naïve call for an account balance could show the wrong currency symbol. Verify and address this carefully before coding; masking and currency labeling must not contradict each other.

### Account type translations

The type union is confirmed, but locale keys may not exist. Search all locales and account screens before adding keys; avoid duplicating translations if equivalent keys already exist.

### Long values and narrow screens

A long account name or large formatted balance can compete for width. Keep the amount column non-shrinking, use `min-width: 0` and truncate the name/type column, and verify 360px with CJK/Indonesian strings. Do not shrink amounts to unreadable sizes.

### Archived accounts

Store filters start with `archived: ''`; inclusion semantics depend on API behavior. Do not change filtering or assume archived accounts are absent. Preserve the fetched list; add a marker only if needed and localized.

### Design-token mismatch

Some global CSS values may lag behind `DESIGN.md`. Use the existing semantic variable names; do not update global tokens as a side effect. Call out any discrepancy separately if it affects this section.

### App shell owns FAB

The FAB may be rendered by the app shell rather than this view. Do not change shell files. Verify no new section overlays or obscures it at small widths.

## 11. Suggested manual review checklist

- [ ] At first glance, are the four key metrics the dashboard’s most important summary?
- [ ] Can a user compare income, expense, and net without decoding color alone?
- [ ] Is total balance visually distinct from monthly flow values without looking secondary or disabled?
- [ ] Does the account section feel like a quick overview rather than a second Accounts page?
- [ ] Can the user scan account name/type and balance in one pass?
- [ ] Are different currencies explicit and accurate?
- [ ] Does the compact design still have adequate text size and touch/accessibility affordances?
- [ ] Is the visual rhythm consistent with Kiroku’s paper-and-ink design?
- [ ] Did the implementation avoid touching the FAB and unrelated local changes?

## 12. Handoff notes

Implementation is not complete until the build is verified and visual limitations are honestly reported. Keep changes small, preserve existing worktree state, and do not commit unless explicitly asked.
