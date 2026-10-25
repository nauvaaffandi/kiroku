# Test Plan: Mobile-Only UI Refactoring

## 1. Scope and Strategy
- **Objective:** Verify that the mobile-only UI refactoring preserves all existing application functionality while successfully enforcing a single-column mobile view, eliminating layout overflows, and resolving UI glitches.
- **Approach:**
  - Automated test execution of the existing frontend test suite (`npm run test:run` / vitest) to confirm logic preservation.
  - Lint and type check verification (`npm run lint` / `npm run typecheck`) to ensure structural code integrity.
  - Smoke tests and CSS constraint checks to ensure desktop residues (sidebars, multi-column max-width containers) are removed.

## 2. Test Cases

### TC-01: Frontend Unit & Integration Tests
- **Target:** All frontend component and store logic tests.
- **Command:** `npm run test:run` (run within `frontend/` directory)
- **Acceptance Criteria:** 100% of existing tests pass successfully.

### TC-02: Type Safety Check
- **Target:** TypeScript compilation and type definitions.
- **Command:** `npm run typecheck` or `vue-tsc --noEmit` (run within `frontend/` directory)
- **Acceptance Criteria:** Zero type errors.

### TC-03: Linter Compliance
- **Target:** Code style and syntax checks.
- **Command:** `npm run lint` (run within `frontend/` directory)
- **Acceptance Criteria:** Zero lint errors or warnings introduced.

### TC-04: Layout & Overflow Verification (Static Check)
- **Target:** `AppShell.vue` and global CSS assets (`main.css`).
- **Criteria:** No desktop sidebar spacing offsets, no desktop media queries forcing multi-column views, standard mobile padding and container constraints applied.
