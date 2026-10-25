---
name: Kiroku
version: 1.0
description: >
  Quiet Ledger. A calm, paper-and-ink design system for a personal finance
  tracker. Modern but restrained. Light mode is warm paper, dark mode is zinc.

colors:
  # Light mode: paper and ink
  canvas: "#f6f5f2"
  surface: "#fcfbf9"
  surface-sunken: "#efede8"
  border: "#e6e3dd"
  border-strong: "#d2cec5"
  ink: "#1f1e1b"
  ink-secondary: "#55524c"
  ink-muted: "#6f6b63"
  on-ink: "#fcfbf9"
  accent: "#2b67ad"
  accent-hover: "#245a99"
  accent-soft: "#e6eef7"
  income: "#0d7a50"
  income-soft: "#e3f2ea"
  expense: "#c0213f"
  expense-soft: "#f9e6ea"
  warning: "#a15c07"
  warning-soft: "#faf2e4"

  # Dark mode: zinc. Same roles, same names, "-dark" suffix.
  canvas-dark: "#121212"
  surface-dark: "#18181b"
  surface-sunken-dark: "#27272a"
  border-dark: "#27272a"
  border-strong-dark: "#3f3f46"
  ink-dark: "#f4f4f5"
  ink-secondary-dark: "#a1a1aa"
  ink-muted-dark: "#93939c"
  on-ink-dark: "#18181b"
  accent-dark: "#7db7ee"
  accent-hover-dark: "#9ccbf4"
  accent-soft-dark: "#16222e"
  income-dark: "#34d399"
  income-soft-dark: "#0f2a20"
  expense-dark: "#fb7185"
  expense-soft-dark: "#2e1519"
  warning-dark: "#fbbf24"
  warning-soft-dark: "#2b2110"

typography:
  display:
    fontFamily: &sans '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Hiragino Sans", "Noto Sans JP", "Yu Gothic UI", sans-serif'
    fontSize: 32px
    fontWeight: 600
    lineHeight: 1.1
    letterSpacing: -0.02em
    fontFeature: '"tnum" 1'
  amount-lg:
    fontFamily: *sans
    fontSize: 24px
    fontWeight: 600
    lineHeight: 1.2
    letterSpacing: -0.01em
    fontFeature: '"tnum" 1'
  amount:
    fontFamily: *sans
    fontSize: 15px
    fontWeight: 600
    lineHeight: 1.3
    fontFeature: '"tnum" 1'
  title:
    fontFamily: *sans
    fontSize: 20px
    fontWeight: 600
    lineHeight: 1.3
    letterSpacing: -0.01em
  heading:
    fontFamily: *sans
    fontSize: 16px
    fontWeight: 600
    lineHeight: 1.4
  body:
    fontFamily: *sans
    fontSize: 14px
    fontWeight: 400
    lineHeight: 1.55
  body-strong:
    fontFamily: *sans
    fontSize: 14px
    fontWeight: 500
    lineHeight: 1.55
  label:
    fontFamily: *sans
    fontSize: 13px
    fontWeight: 500
    lineHeight: 1.3
  caption:
    fontFamily: *sans
    fontSize: 12px
    fontWeight: 400
    lineHeight: 1.4

rounded:
  sm: 6px
  md: 10px
  lg: 14px
  full: 9999px

spacing:
  unit: 4px
  xs: 4px
  sm: 8px
  md: 12px
  lg: 16px
  xl: 24px
  2xl: 32px
  3xl: 48px
  gutter: 16px
  touch: 44px

# Components reference light-mode tokens. In dark mode every color token X
# resolves to X-dark through the CSS variables, so components never branch.
components:
  button-primary:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.on-ink}"
    typography: "{typography.label}"
    rounded: "{rounded.md}"
    height: 44px
    padding: "0 16px"
  button-primary-hover:
    opacity: 0.9
  button-secondary:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    typography: "{typography.label}"
    rounded: "{rounded.md}"
    border: "1px solid {colors.border-strong}"
    height: 44px
    padding: "0 16px"
  button-secondary-hover:
    backgroundColor: "{colors.surface-sunken}"
  button-ghost:
    backgroundColor: transparent
    textColor: "{colors.ink-secondary}"
    typography: "{typography.label}"
    rounded: "{rounded.md}"
    height: 44px
    padding: "0 12px"
  button-ghost-hover:
    backgroundColor: "{colors.surface-sunken}"
    textColor: "{colors.ink}"
  button-danger:
    backgroundColor: "{colors.expense}"
    textColor: "{colors.on-ink}"
    typography: "{typography.label}"
    rounded: "{rounded.md}"
    height: 44px
    padding: "0 16px"
  button-sm:
    height: 36px
    padding: "0 12px"
  input:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    typography: "{typography.body}"
    rounded: "{rounded.md}"
    border: "1px solid {colors.border-strong}"
    height: 44px
    padding: "0 12px"
  input-focus:
    border: "1px solid {colors.accent}"
    ring: "2px {colors.accent}"
  input-error:
    border: "1px solid {colors.expense}"
  card:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.lg}"
    border: "1px solid {colors.border}"
    padding: "{spacing.lg}"
    shadow: none
  list-row:
    backgroundColor: transparent
    textColor: "{colors.ink}"
    typography: "{typography.body}"
    border-bottom: "1px solid {colors.border}"
    padding: "{spacing.md} 0"
    minHeight: 56px
  amount-income:
    textColor: "{colors.income}"
    typography: "{typography.amount}"
  amount-expense:
    textColor: "{colors.ink}"
    typography: "{typography.amount}"
  amount-transfer:
    textColor: "{colors.ink-secondary}"
    typography: "{typography.amount}"
  stat:
    backgroundColor: transparent
    textColor: "{colors.ink}"
    typography: "{typography.amount-lg}"
    label: "{typography.caption} in {colors.ink-muted}"
  badge:
    backgroundColor: "{colors.surface-sunken}"
    textColor: "{colors.ink-secondary}"
    typography: "{typography.caption}"
    rounded: "{rounded.sm}"
    padding: "2px 8px"
  badge-income:
    backgroundColor: "{colors.income-soft}"
    textColor: "{colors.income}"
  badge-expense:
    backgroundColor: "{colors.expense-soft}"
    textColor: "{colors.expense}"
  badge-warning:
    backgroundColor: "{colors.warning-soft}"
    textColor: "{colors.warning}"
  nav-header:
    backgroundColor: "{colors.surface}"
    border-bottom: "1px solid {colors.border}"
    height: 56px
  nav-bottom:
    backgroundColor: "{colors.surface}"
    border-top: "1px solid {colors.border}"
    height: 60px
  nav-item:
    textColor: "{colors.ink-muted}"
    typography: "{typography.caption}"
  nav-item-active:
    textColor: "{colors.ink}"
    indicator: "2px {colors.accent} at top edge"
  add-button:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.on-ink}"
    rounded: "{rounded.lg}"
    size: 52px
    shadow: float
  sheet:
    backgroundColor: "{colors.surface}"
    rounded: "{rounded.lg} {rounded.lg} 0 0"
    shadow: overlay
    padding: "{spacing.lg}"
  toast:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    typography: "{typography.body-strong}"
    rounded: "{rounded.md}"
    border: "1px solid {colors.border}"
    shadow: float
    marker: ribbon
---

# Kiroku Design System

> Read this before writing or changing any UI. It is the contract for how Kiroku looks and behaves. The YAML block above is the **schema layer** (tokens). Everything below is the **sections layer** (how to use them).

## Overview

**Concept: Quiet Ledger.** 記録 means "record". Kiroku should feel like a well-kept notebook: calm paper, ink-colored type, hairline rules, numbers lined up in tidy columns. People open it for a few seconds, one-handed, to log a coffee or check a balance. Nothing should compete with the numbers.

**Modern, but simple.** Modern means crisp type, generous but disciplined spacing, flat surfaces and fast, quiet motion. It does not mean gradients, glass, glow or decoration. When in doubt, remove something.

**Principles**

1. **Numbers are the interface.** Amounts get the strongest typography on every screen. Everything else supports them.
2. **Ink does the work, color is rare.** The primary control is ink on paper (near-black in light, near-white in dark). Color appears only when it carries meaning: money in, money out, a warning, focus.
3. **Flat by default.** Structure comes from hairline borders, spacing and alignment. Shadow is reserved for layers that float.
4. **Calm over clever.** One idea per screen. Short copy. No celebration, no mascots, no filler.
5. **Fast to scan, fast to tap.** One-handed, 360px wide, thumb zone first.

**Signature details (use sparingly, they only work if rare)**

- The **ribbon marker**, a small bookmark notch in the top-left of a toast, is the only decorative element in the system.
- **Tabular, right-aligned numerals** everywhere money appears.
- A **2px accent line** at the top edge of the active bottom-nav item.

**Languages.** The UI ships in English, 日本語 and Bahasa Indonesia. Layouts must survive all three without truncating amounts.

### Token to CSS variable map

Tokens live in `frontend/src/assets/main.css` as CSS variables and flip with `[data-theme]`. **Use only these names in components** (see AGENTS.md for how to apply them with `@reference "tailwindcss"`).

| Token (light / `-dark`)          | CSS variable                              |
| -------------------------------- | ----------------------------------------- |
| `canvas`                         | `--bg-primary`                            |
| `surface`                        | `--bg-card`, `--bg-secondary`             |
| `surface-sunken`                 | `--bg-tertiary`, `--bg-hover`             |
| `border`                         | `--border-primary`                        |
| `border-strong`                  | `--border-secondary`                      |
| `ink`, `ink-secondary`, `ink-muted` | `--text-primary`, `--text-secondary`, `--text-muted` |
| `on-ink`                         | `--text-inverse`                          |
| `accent`, `accent-hover`, `accent-soft` | `--accent`, `--accent-hover`, `--accent-light` (focus ring uses `--border-focus` = `accent`) |
| `income`, `income-soft`          | `--success`, `--success-light`            |
| `expense`, `expense-soft`        | `--error`, `--error-light`                |
| `warning`, `warning-soft`        | `--warning`, `--warning-light`            |
| informational                    | `--info` = `accent`, `--info-light` = `accent-soft` |
| `rounded.sm / md / lg / full`    | `--radius-sm`, `--radius`, `--radius-lg`, `--radius-full` |

> **Sync note.** The values in this file are the target. If `main.css` still holds older values (cool grays, `#4a90d9` accent, 12px radius), updating `main.css` to match is its own task. Change both in the same commit and never let them drift.

## Colors

### Light: warm paper

| Token            | Hex       | Role                                                      |
| ---------------- | --------- | --------------------------------------------------------- |
| `canvas`         | `#f6f5f2` | Page background                                           |
| `surface`        | `#fcfbf9` | Cards, sheets, header, inputs                             |
| `surface-sunken` | `#efede8` | Hover, pressed, badges, skeletons, quiet fills            |
| `border`         | `#e6e3dd` | Hairline dividers, card outlines                          |
| `border-strong`  | `#d2cec5` | Input outlines, secondary buttons                         |
| `ink`            | `#1f1e1b` | Primary text, primary buttons                             |
| `ink-secondary`  | `#55524c` | Supporting text                                           |
| `ink-muted`      | `#6f6b63` | Captions, placeholders, inactive nav                      |
| `accent`         | `#2b67ad` | Focus ring, links, active-nav line, forecast/info series  |
| `income`         | `#0d7a50` | Money in                                                  |
| `expense`        | `#c0213f` | Errors, destructive actions, money-out emphasis           |
| `warning`        | `#a15c07` | Warnings, budget nearly spent                             |

There is no pure white and no pure black anywhere in light mode.

### Dark: zinc

Dark mode uses the zinc scale and nothing else for neutrals.

| Token            | Hex       | Zinc step | Role                              |
| ---------------- | --------- | --------- | --------------------------------- |
| `canvas-dark`    | `#121212` | (custom)  | Page background                   |
| `surface-dark`   | `#18181b` | 900       | Cards, sheets, header, inputs     |
| `surface-sunken-dark` | `#27272a` | 800  | Hover, pressed, badges            |
| `border-dark`    | `#27272a` | 800       | Hairlines                         |
| `border-strong-dark` | `#3f3f46` | 700   | Input outlines                    |
| `ink-dark`       | `#f4f4f5` | 100       | Primary text, primary buttons     |
| `ink-secondary-dark` | `#a1a1aa` | 400   | Supporting text                   |
| `ink-muted-dark` | `#93939c` | between 400 and 500 | Captions, placeholders  |
| `accent-dark`    | `#7db7ee` | (custom)  | Focus, links, active line         |
| `income-dark`    | `#34d399` | emerald 400 | Money in                        |
| `expense-dark`   | `#fb7185` | rose 400  | Errors, destructive               |
| `warning-dark`   | `#fbbf24` | amber 400 | Warnings                          |

Dark mode is **flat**: depth comes from the lightness step between `canvas`, `surface` and `surface-sunken`, not from shadows.

### Contrast (WCAG 2.x, computed)

| Pair                                  | Light | Dark |
| ------------------------------------- | ----- | ---- |
| `ink` on `surface`                    | 16.1  | 16.1 |
| `ink-secondary` on `surface`          | 7.5   | 6.9  |
| `ink-muted` on `canvas`               | 4.9   | 6.2  |
| `ink-muted` on `surface-sunken`       | 4.5   | 4.9  |
| `accent` on `surface`                 | 5.6   | 8.3  |
| `income` on `surface`                 | 5.2   | 9.2  |
| `expense` on `surface`                | 5.8   | 6.6  |
| `income` on `income-soft`             | 4.6   | 8.0  |
| `warning` on `warning-soft`           | 4.7   | 9.5  |
| `on-ink` on `ink` (primary button)    | 16.1  | 16.1 |

Body text and anything interactive must stay at 4.5:1 or better. Borders are decorative hairlines (about 1.2:1 to 1.5:1 on purpose), so **never use a border as the only cue** for state or error.

### Usage rules

- **One accent.** `accent` appears for focus, links, the active-nav line and forecast or informational chart series. Not on buttons, not on cards, not as a decoration.
- **Meaning-only color.** `income` and `expense` color money and money status only. Do not use green or red for branding, section headers or icons.
- **Expense amounts stay in `ink`.** A list full of red numbers is loud and tiring. Expenses are shown with a minus sign in `ink`, income with a plus sign in `income`. Use `expense` for errors, destructive actions, overspent budgets and negative net totals.
- **Never color alone.** Pair color with a sign (+ or −), an icon or a label.
- **Tints are for badges and banners only**, using the `-soft` token with its matching text color.
- **No gradients, anywhere.** Not on backgrounds, buttons, text, charts or skeletons.
- **No new colors.** Category and account colors chosen by the user are data. Show them as small dots (8px) or thin bars, never as large filled areas.

### Charts

Series order: `ink-secondary`, `accent`, `income`, `warning`, `ink-muted`. Income vs expense charts: `income` and `ink` (expense as ink bars) or `expense` when highlighting overspend. Area fills are a flat 8% tint of the line color. Grid lines are `border`, axis text is `ink-muted`.

## Typography

**One family, system stack.** Zero font downloads (the app is an offline-capable PWA), native feel, and correct CJK glyphs.

```css
font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
  "Hiragino Sans", "Noto Sans JP", "Yu Gothic UI", sans-serif;
```

### Scale

| Token         | Size / line | Weight | Use                                      |
| ------------- | ----------- | ------ | ---------------------------------------- |
| `display`     | 32 / 1.1    | 600    | Total balance, one per screen at most    |
| `amount-lg`   | 24 / 1.2    | 600    | Stat values, account balance             |
| `title`       | 20 / 1.3    | 600    | Page title in content, sheet titles      |
| `heading`     | 16 / 1.4    | 600    | Section headings, card titles            |
| `body`        | 14 / 1.55   | 400    | Default text                             |
| `body-strong` | 14 / 1.55   | 500    | Row titles, emphasis                     |
| `amount`      | 15 / 1.3    | 600    | Amounts in lists and tables              |
| `label`       | 13 / 1.3    | 500    | Form labels, buttons                     |
| `caption`     | 12 / 1.4    | 400    | Metadata, helper text, nav labels        |

### Rules

- **Weights: 400, 500, 600 only.** No 700, 800 or `font-black`. Hierarchy comes from size, weight and color together.
- **Tabular numerals for every number** that can change or align: amounts, dates, percentages, counters. Use `font-variant-numeric: tabular-nums` (Tailwind: `tabular-nums`).
- **Right-align amounts** and keep the currency symbol attached to the number. Format with `Intl.NumberFormat` and the user's number format. JPY has no decimals.
- **Minus is a real minus sign** (−, U+2212), not a hyphen.
- **Sentence case.** No ALL CAPS, no letter-spaced micro labels. Capitalize like a sentence.
- **Japanese:** `line-height` at least 1.6 for body text. Do not apply negative letter-spacing to CJK.
- **Line length:** keep prose under about 60 characters wide. In a 480px column this happens naturally.
- **Hidden amounts:** when "hide amounts" is on, render `••••` in the same typography so layout does not jump.

## Layouts

### Frame

- **Mobile-first single column**, centered, max width **480px**. Do not build desktop multi-column dashboards unless explicitly asked. If a wider view is requested, cap content at 720px and keep it a single column.
- **Gutter:** 16px (`spacing.gutter`) on both sides.
- **Header:** sticky, 56px, `surface` with a bottom hairline. Page title left, utilities right.
- **Bottom nav:** fixed, 60px plus `env(safe-area-inset-bottom)`. Four destinations and the add action.
- **Page padding:** content starts 16px below the header and ends with enough bottom padding to clear the nav (about 84px plus the safe area).

### Spacing

Base unit is 4px. Use the scale `4 / 8 / 12 / 16 / 24 / 32 / 48` and nothing in between.

| Use                                     | Value |
| --------------------------------------- | ----- |
| Icon to label                           | 8     |
| Inside rows and compact controls        | 12    |
| Inside cards, between form fields       | 16    |
| Between sections on a page              | 24    |
| Between a page title and its content    | 24    |
| Large breaks (empty states, onboarding) | 32–48 |

### Structure

- **Lists beat cards.** Repeated items (transactions, accounts, categories, scheduled items) are **rows separated by hairlines**, not a stack of boxed cards. Use a card only for a single self-contained group: a form, a chart, a summary.
- **Group by date**, not by decoration. Date headers are `caption` in `ink-muted` with the day total aligned right.
- **Top of page:** the one number that matters (`display` or `amount-lg`), a caption under it, then content. Left-aligned.
- **Grids:** two columns at most, only for paired figures (for example income and expense). Never a three-up feature grid.
- **Forms:** one column, labels above fields, 16px between fields, primary action pinned at the bottom of the sheet.
- **Tables** (reports, cash flow): scroll horizontally inside their own container, first column sticky, numeric columns right-aligned.
- **Touch targets:** at least 44×44px. Small visual controls get extra hit area through padding.
- **Z-order:** header and nav 100, drawer 200, toast 400, sheet and modal 500.

## Elevation & Depth

Flat by default. Depth is communicated in this order of preference: spacing, hairline border, surface lightness step, and only last, shadow.

| Level    | Used for                             | Light                                   | Dark                                  |
| -------- | ------------------------------------ | --------------------------------------- | ------------------------------------- |
| 0 flat   | Page, cards, rows, header, nav       | `canvas`, hairline `border`             | `canvas-dark`, hairline `border-dark` |
| 1 sunken | Hover, pressed, inputs at rest       | `surface-sunken`                        | `surface-sunken-dark`                 |
| 2 float  | Dropdown, toast, add button          | `0 4px 16px rgb(31 30 27 / 0.08)` plus hairline | hairline only, `surface-dark` |
| 3 overlay| Sheet, modal, drawer                 | `0 12px 40px rgb(31 30 27 / 0.14)`      | scrim, `surface-dark`, hairline       |

- **Scrim:** `rgb(0 0 0 / 0.4)` in light, `rgb(0 0 0 / 0.6)` in dark. No backdrop blur.
- **Cards never have shadows.** A card on `canvas` is already separated by its fill and hairline.
- **Never stack** border, shadow and a tinted fill on the same element. Pick one or two.
- **No colored or glowing shadows**, no inner shadows, no "elevated on hover" effects.

## Shapes

| Token  | Radius | Use                                                     |
| ------ | ------ | ------------------------------------------------------- |
| `sm`   | 6px    | Badges, small chips, table cells with fills             |
| `md`   | 10px   | Buttons, inputs, selects, toasts, dropdowns             |
| `lg`   | 14px   | Cards, sheets (top corners), add button                 |
| `full` | 9999px | Toggle switches, color dots, avatars. Nothing else.     |

- **Hairlines:** 1px. No 2px borders except the focus ring and the active-nav line.
- **Focus ring:** `:focus-visible` only, 2px `accent`, 2px offset. Never remove outlines without replacing them.
- **Icons:** `@lucide/vue`, 20px, stroke 1.75, `currentColor`. Same stroke width everywhere. Icons stand alone; they are not placed inside tinted squares or circles.
- **Category and account marks:** an 8px color dot next to the name. If the user picked an emoji, show it at body size inline with the name, not inside a colored tile.
- **Corners are consistent.** Nested radii shrink toward the inside (a 10px button inside a 14px card keeps 16px of padding).

## Components

### Button

- **Primary:** `ink` fill, `on-ink` text, 44px, radius `md`. One per view.
- **Secondary:** `surface` fill, `border-strong` outline, `ink` text.
- **Ghost:** transparent, `ink-secondary` text, `surface-sunken` on hover. For toolbars and low-priority actions.
- **Danger:** `expense` fill, `on-ink` text. Only in the confirm step of a destructive flow; the first tap is a secondary or ghost button.
- **Link:** `accent` text, underline on hover. Inline only.
- **Sizes:** default 44px; `sm` 36px for dense toolbars (extend the hit area).
- **States:** hover primary `opacity: 0.9`; hover secondary and ghost `surface-sunken`; pressed adds a 1px downward shift, no scale; disabled `opacity: 0.5` with `cursor: not-allowed`; loading replaces the label with a 16px spinner and keeps the width.
- **No** icon-only primary buttons without an accessible label. **No** gradients, glow or pill shape (radius is `md`).

### Input and select

- 44px high, `surface` fill, 1px `border-strong`, radius `md`, `body` text, `ink-muted` placeholder.
- Label above in `label` style. Helper or error text below in `caption`.
- Focus: border turns `accent` and a 2px `accent` ring appears.
- Error: `expense` border **plus** an icon and a message in `expense`. Never a red border alone.
- Amount fields use `inputmode="decimal"`, tabular numerals and right alignment.
- Prefix and suffix (currency symbol, unit) sit inside the field in `ink-muted`.

### Card

`surface` fill, 1px `border`, radius `lg`, 16px padding, no shadow. Title in `heading`, content below. A card never contains another card.

### Transaction row (the core component)

```
[dot] Title (body-strong)                         −¥1,280 (amount)
      Account · note (caption, ink-muted)
```

- Min height 56px, 12px vertical padding, hairline below, full row is the tap target.
- Title is the payee, or the category name if no payee. Subtitle is account and note, truncated to one line.
- Left: an 8px category dot. Right: the amount, tabular, right-aligned.
- Expense `−` in `ink`; income `+` in `income`; transfer shows `→ Account` in `ink-secondary` with no sign.
- Status appears only when it is not the default (for example a `pending` badge).
- Edit and delete live in a swipe action or the detail sheet, not as two icons on every row.

### Amount

Always `amount`, `amount-lg` or `display` with tabular numerals. Currency symbol directly attached. Honor the hide-amounts setting. Show `—` for missing values, never `0` unless it is zero.

### Stat (summary figure)

Replaces boxed "icon + number" tiles. A stat is a `caption` label in `ink-muted` above an `amount-lg` value, with no icon, no fill and no border. Place two or three stats in a row separated by a vertical hairline, or inside one card. Net totals use `income` when positive and `expense` when negative, always with a sign.

### Badge and chip

`surface-sunken` fill, `ink-secondary` text, `caption`, radius `sm`, padding `2px 8px`. Semantic variants use the matching `-soft` fill and text color. A badge states a fact (pending, archived, recurring); it is not decoration. Filter chips are the same shape with `ink` text and a `border-strong` outline when selected.

### Navigation

- **Header:** menu button, page title (`heading`), currency selector and overflow menu on the right. No logo lockups, no gradient bars.
- **Bottom nav:** four items (Home, Transactions, Accounts, Reports) with the add action between the second and third. Inactive: `ink-muted`. Active: `ink`, with a 2px `accent` line on the top edge of the item. Labels are full words that fit in all three languages.
- **Add button:** 52px, radius `lg`, `ink` fill, `on-ink` plus icon, level 2 shadow in light mode. One, not a speed-dial fan of colored circles.
- **Drawer:** `surface`, hairline on the edge, rows like a list. Active item gets `surface-sunken` and `ink` text.

### Sheet and modal

On mobile, dialogs are **bottom sheets**: full shell width, `surface`, top corners `lg`, max height 90vh, scrim behind. Header has the title and a close button. Primary action pinned at the bottom. Close on scrim tap unless the form is dirty. Trap focus, restore it on close, support Escape.

### Toast

`surface` fill, 1px `border`, radius `md`, level 2 shadow. The **ribbon marker** at the top-left carries the semantic color (`income` success, `expense` error, `warning`, `accent` info). Text in `body-strong`, one line, with a dismiss button. Auto-dismiss in 4s. Slide in from the bottom, 180ms. This is the only place the ribbon appears.

### Charts (ECharts)

Flat fills at 8% tint, 1.5px lines, `border` gridlines, `ink-muted` axis labels at `caption` size, no gradients, no 3D, no shadows, no animated count-ups. Bars have at most 2px corner radius. Tooltips use `surface`, a hairline and `body` text. Always show the value with its sign and currency. Provide a text or table alternative for key charts.

### Skeleton

`surface-sunken` shapes that match the real layout. Animate with a slow opacity pulse (100% to 60%, 1.6s ease-in-out), not a moving gradient. Disable under `prefers-reduced-motion`.

### Empty state

One short sentence that says what will appear here, one primary action. No illustration, no emoji, no exclamation mark. Left-aligned or centered within a card, 32–48px of vertical space.

### Motion

| Purpose                    | Duration | Easing                      |
| -------------------------- | -------- | --------------------------- |
| Hover, press, color change | 120ms    | `ease-out`                  |
| Toast, dropdown            | 180ms    | `cubic-bezier(0.2, 0, 0, 1)`|
| Sheet, drawer              | 240ms    | `cubic-bezier(0.2, 0, 0, 1)`|

Animate only `opacity`, `transform` and color. No bounce, spring, parallax or hover lift. Page transitions are a 200ms fade. Everything respects `prefers-reduced-motion: reduce` by dropping to an instant change.

## Do's and Don'ts

### Do

- Start from the content: what number or action does the person need in the first two seconds?
- Use hairlines, spacing and alignment to build structure.
- Keep amounts right-aligned, tabular and the largest thing in their row.
- Use `ink` for primary actions and let color mean something.
- Check every screen in light **and** dark, at 360px, in Japanese and Indonesian, with long account names and large amounts (¥12,345,678).
- Write plain copy: say what happened and what to do next. "Couldn't save. Check your connection and try again."
- Support keyboard, screen readers and reduced motion by default.
- Reuse existing components and tokens before creating new ones. If something new is truly needed, add the token to this file first.
- Keep one primary button per view and one accent usage per screen region.
- Use skeletons that match the final layout instead of spinners for page loads.

### Don't (the AI-slop list)

- **No gradients.** Not on backgrounds, buttons, headings, charts or skeletons. No purple-to-blue anything.
- **No glassmorphism**, backdrop blur, frosted panels, glow, neon or colored shadows.
- **No decorative blobs**, orbs, mesh backgrounds, sparkles, grain overlays or abstract hero art.
- **No icon-in-a-tinted-square** on every card, row or stat. No emoji as UI icons or in copy.
- **No three-column feature-card grids**, no "bento" layouts, no cards inside cards, no card per list item.
- **No oversized radii** (`rounded-2xl`, `rounded-3xl`) and no pill-shaped everything. Radius is 6, 10 or 14px.
- **No shadow on resting cards**, no hover lift, no scale-on-hover, no bouncy or spring animation, no count-up number animation.
- **No heavy type:** no `font-bold` or `font-black` blocks, no ALL CAPS micro labels with wide tracking, no gradient text.
- **No second brand color.** No default indigo, violet or teal "because it looks modern".
- **No red-for-every-expense.** Red is for errors and overspend, not for routine spending.
- **No color-only meaning.** Always pair with a sign, icon or text.
- **No generic copy:** "Welcome back 👋", "Oops!", "Let's get started", "Supercharge your finances", fake testimonials, lorem ipsum, placeholder metrics.
- **No hard-coded colors, radii or spacing.** No hex values and no invented `--color-*` variables in components. Use the tokens above through the CSS variables listed in the map.
- **No centered-everything layouts.** Left-align text and numbers; center only isolated empty states.
- **No new dependencies for UI** (component libraries, animation libraries, icon sets other than `@lucide/vue`) without asking.

### Pre-flight checklist

Before you hand back any UI work, answer yes to all of these:

1. Does it still look right with the dark theme and with the light theme, with no dark-specific hacks?
2. Is every color, radius and space value a token from this file?
3. Is there at most one primary button and at most one accent usage in view?
4. Are all amounts tabular, right-aligned and honoring hide-amounts?
5. Could you remove any border, shadow, fill or icon without losing meaning? If yes, remove it.
6. Does it work at 360px with Japanese and Indonesian strings?
7. Are touch targets at least 44px and focus visible?
8. Would this screenshot look at home next to the Transactions screen, or does it look like a template?

---

**How to apply these tokens in code:** follow `AGENTS.md` section 3. Component styles live in `<style scoped>`, start with `@reference "tailwindcss";`, use `@apply` for layout, spacing, type and radius, and read colors from the CSS variables.

```vue
<style scoped>
@reference "tailwindcss";

.card {
  @apply flex flex-col gap-3 rounded-[14px] border p-4;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  color: var(--text-primary);
}

.card-title {
  @apply text-base font-semibold;
}
</style>
```
