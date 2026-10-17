# Task objective
Refactor `frontend/src/views/transactions/ScanReceiptView.vue` (482 lines) into smaller, modular, and readable single-file components to reduce cognitive load and separate concerns.

# Scope and exclusions
- **Scope:** Extract sub-components from `ScanReceiptView.vue`:
  1. `ReceiptUploadCard.vue`
  2. `ReceiptProcessingCard.vue`
  3. `ReceiptReviewCard.vue`
- **Exclusions:** Backend changes, API client changes, style changes.

# Existing architecture and relevant files
- `frontend/src/views/transactions/ScanReceiptView.vue`

# Current behavior and desired behavior
- **Current:** Everything (upload, compression logic, polling, review form) is in one 482-line file.
- **Desired:** `ScanReceiptView.vue` acts as the orchestrator managing state, image compression, and API calls, while rendering modular child components for each step.

# Root-cause analysis
- Monolithic component design in frontend views increases cognitive load and makes maintenance difficult.

# Assumptions and unresolved questions
- None.

# Numbered implementation steps
### Step 1: Create components directory and extract `ReceiptUploadCard.vue`
### Step 2: Extract `ReceiptProcessingCard.vue`
### Step 3: Extract `ReceiptReviewCard.vue`
### Step 4: Refactor `ScanReceiptView.vue` to use extracted components

# File-level change descriptions
- Create `frontend/src/views/transactions/components/ReceiptUploadCard.vue`
- Create `frontend/src/views/transactions/components/ReceiptProcessingCard.vue`
- Create `frontend/src/views/transactions/components/ReceiptReviewCard.vue`
- Modify `frontend/src/views/transactions/ScanReceiptView.vue`

# Dependencies and execution order
- Steps 1, 2, 3 can be done in parallel or sequence, followed by Step 4.

# Acceptance criteria
- `ScanReceiptView.vue` is significantly smaller (< 200 lines).
- Extracted components are clean and single-purpose.
- Frontend builds successfully without TypeScript or Vue compilation errors.

# Verification commands
- Run `npm run build` inside the frontend directory.

# Regression risks
- Broken event passing between parent and child components.

# Rollback considerations
- Use git checkout to revert changes if the build or verification fails.
