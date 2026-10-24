# Implementation Plan: Mobile-Only UI Refactoring

## 1. Objectives
- Remove all desktop layouts, sidebars, and responsive residues, enforcing a strict **mobile-only** single-column layout.
- Resolve UI glitches (horizontal overflow, incorrect padding, max-width constraints).
- Preserve 100% of the existing application logic, state management (Pinia), routing, and API integrations.

## 2. Steps

### Step 1: Enforce Mobile Layout in `AppShell.vue` & Styles
- **File:** `frontend/src/components/layout/AppShell.vue`, `frontend/src/assets/main.css`
- **Actions:** Remove desktop sidebar integration, remove media queries enforcing desktop side margins (e.g. `margin-left: 280px`), and simplify shell to header + scrollable content + bottom nav.
- **Verification:** Confirm mobile layout fills the viewport width without desktop shifts.

### Step 2: Fix Container & Spacing Constraints
- **File:** `frontend/src/assets/main.css` or layout containers
- **Actions:** Remove large desktop max-width constraints (e.g., `max-width: 1200px`) and standardize mobile padding (`1rem`).
- **Verification:** Check visual alignment on mobile viewports.

### Step 3: Fix Table & Widget Overflow Glitches
- **File:** Relevant widgets and views under `frontend/src/widgets/` and `frontend/src/views/`
- **Actions:** Convert overflowing tables/grids into stacked mobile-friendly card lists where necessary.
- **Verification:** Run `npm run test:run` in `frontend/` to ensure tests pass.

### Step 4: Final Verification
- **Actions:** Run frontend test suite (`npm run test:run`) and typecheck/lint.
