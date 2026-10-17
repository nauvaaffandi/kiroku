# Hide Amount Toggle Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a toggle setting to hide/show transaction amounts in the UI for privacy purposes.

**Architecture:** 
- Add a new boolean setting `hide_amounts` (default false) to the settings table via backend defaults.
- Extend the frontend settings store to include this setting.
- Add a toggle switch in the Settings view under the Data section.
- Create a utility function `formatAmountForDisplay` that returns masked amounts when the setting is enabled.
- Replace all amount formatting calls in the UI with this utility.

**Tech Stack:** Go (backend), Vue 3/Pinia (frontend), SQLite

**Spec:** /data/data/com.termux/files/home/developments/kiroku/PRD.md (Sections 7.12 Settings, 9.2 UX Principles)

## Global Constraints
- No new dependencies may be added (YAGNI)
- Must maintain existing functionality when setting is disabled
- Setting must persist across sessions
- Amount masking must preserve sign indicator (income/expense/transfer) where applicable
- Backend must store setting as string "true"/"false" in settings table

## Review Focus
- Ensure setting defaults to false on fresh install
- Verify masked amounts appear consistently across all views (transactions, reports, analytics, shared)
- Confirm setting persists after page reload
- Check that amount inputs in forms remain unmasked for editing
- Validate that masking does not break layout (fixed width placeholder)

---
### Task 1: Backend - Add default for hide_amounts setting

**Files:**
- Modify: `/data/data/com.termux/files/home/developments/kiroku/backend/internal/settings/handler.go:68-80`

**Interfaces:**
- Consumes: None
- Produces: Updated settings GET endpoint returning hide_amounts default

- [ ] **Step 1: Write the failing test**
  ```go
  // Test that GetAll returns hide_amounts default
  func TestSettingsGetAllIncludesHideAmountsDefault(t *testing.T) {
      // Setup test DB and handler
      // Call GetAll
      // Check that settings map contains "hide_amounts": "false"
  }
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `go test ./internal/settings -run TestSettingsGetAllIncludesHideAmountsDefault -v`
  Expected: FAIL with "undefined: hide_amounts" or missing key

- [ ] **Step 3: Implement default addition in GetAll function**
  Add `"hide_amounts": "false"` to the defaults map in handler.go lines 68-80

- [ ] **Step 4: Run test to verify it passes**
  Run: `go test ./internal/settings -run TestSettingsGetAllIncludesHideAmountsDefault -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add backend/internal/settings/handler.go
  git commit -m "settings: add default for hide_amounts setting"
  ```

### Task 2: Frontend - Add hide_amounts to settings store and types

**Files:**
- Create: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/utils/amountDisplay.ts`
- Modify: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/types/settings.ts:15-16`
- Modify: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/stores/settings.ts:15-16,30-33`

**Interfaces:**
- Consumes: None
- Produces: Extended settings interface and store with hide_amounts boolean

- [ ] **Step 1: Write the failing test**
  ```ts
  // Test that settings store includes hide_amounts boolean
  import { useSettingsStore } from '@/stores/settings'
  const store = useSettingsStore()
  expect(typeof store.settings.hide_amounts).toBe('boolean')
  expect(store.settings.hide_amounts).toBe(false)
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `npm test -- --testNamePattern="settings store includes hide_amounts" -v`
  Expected: FAIL (property missing)

- [ ] **Step 3: Implement settings extension**
  - Add `hide_amounts: boolean` to Settings interface in types/settings.ts
  - Add `hide_amounts: false` to settings object in stores/settings.ts line 15
  - Add `hide_amounts: computed(() => settings.value.hide_amounts ?? false)` line 30
  - Add `setHideAmounts: async (value: boolean) => { await setSetting('hide_amounts', value) }` line 33

- [ ] **Step 4: Run test to verify it passes**
  Run: `npm test -- --testNamePattern="settings store includes hide_amounts" -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add frontend/src/types/settings.ts frontend/src/stores/settings.ts
  git commit -m "settings: add hide_amounts boolean setting"
  ```

### Task 3: Frontend - Create amount display utility

**Files:**
- Create: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/utils/amountDisplay.ts`

**Interfaces:**
- Consumes: settings store hide_amounts boolean
- Produces: formatted amount string (masked or normal)

- [ ] **Step 1: Write the failing test**
  ```ts
  // Test that utility returns masked string when hide_amounts=true
  import { useSettingsStore } from '@/stores/settings'
  const settingsStore = useSettingsStore()
  settingsStore.settings.hide_amounts = true
  expect(formatAmountForDisplay(1000, 'income')).toBe('+¥*****')
  expect(formatAmountForDisplay(500, 'expense')).toBe('-¥*****')
  expect(formatAmountForDisplay(0, 'transfer')).toBe('¥*****')
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `npm test -- --testNamePattern="amount display utility masks when enabled" -v`
  Expected: FAIL (utility not implemented)

- [ ] **Step 3: Implement utility function**
  Create amountDisplay.ts with:
  ```ts
  import { useSettingsStore } from '@/stores/settings'

  export function formatAmountForDisplay(amount: number, type: 'income' | 'expense' | 'transfer' = 'transfer'): string {
    const settingsStore = useSettingsStore()
    if (settingsStore.settings.hide_amounts) {
      const sign = type === 'income' ? '+' : type === 'expense' ? '-' : ''
      return `${sign}¥*****`
    }
    // Format as original: sign + ¥ + locale string
    const sign = type === 'income' ? '+' : type === 'expense' ? '-' : ''
    return `${sign}¥${amount.toLocaleString()}`
  }
  ```

- [ ] **Step 4: Run test to verify it passes**
  Run: `npm test -- --testNamePattern="amount display utility masks when enabled" -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add frontend/src/utils/amountDisplay.ts
  git commit -m "utils: add amount display utility for hide amounts setting"
  ```

### Task 4: Frontend - Add toggle to Settings view

**Files:**
- Modify: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/views/settings/SettingsView.vue:118-136` (Data section)
- Modify: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/widgets/settings/SettingItem.vue` (if needed for toggle, but we'll use BaseSelect or create toggle)
  Actually, we'll add a BaseSelect with boolean options or use a custom toggle. Let's check if we have a toggle component.

  Looking at the SettingsView, we see BaseSelect used for other settings. We'll add a BaseSelect for hide_amounts with options:
    - { value: false, label: 'Show Amounts' }
    - { value: true, label: 'Hide Amounts' }

- [ ] **Step 1: Write the failing test**
  ```ts
  // Test that SettingsView renders hide_amounts toggle
  // Mount SettingsView
  // Expect to see a select with options for hide amounts
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `npm test -- --testNamePattern="SettingsView renders hide_amounts toggle" -v`
  Expected: FAIL (toggle not present)

- [ ] **Step 3: Implement toggle in SettingsView**
  Add a new SettingItem in the Data section (after backup, before clearData) with:
    - Label: t('settings.hideAmounts') (need to add translation)
    - Description: t('settings.hideAmountsDesc')
    - BaseSelect with v-model="settings.hide_amounts" and options as above

  Also add translation keys in locales/en.ts (and ja.ts, id.ts) under settings:
    hideAmounts: 'Hide Amounts',
    hideAmountsDesc: 'Replace transaction amounts with asterisks for privacy when sharing screen'

- [ ] **Step 4: Run test to verify it passes**
  Run: `npm test -- --testNamePattern="SettingsView renders hide_amounts toggle" -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add frontend/src/views/settings/SettingsView.vue frontend/src/locales/en.ts frontend/src/locales/ja.ts frontend/src/locales/id.ts
  git commit -m "settings: add hide_amounts toggle to Settings view"
  ```

### Task 5: Frontend - Update TransactionRow to use amount display utility

**Files:**
- Modify: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/widgets/transactions/TransactionRow.vue:84-86` (formatAmount method)

**Interfaces:**
- Consumes: amountDisplay utility
- Produces: formatted amount string for table cell

- [ ] **Step 1: Write the failing test**
  ```ts
  // Test that TransactionRow shows masked amount when hide_amounts=true
  // Mount TransactionRow with transaction amount=1000, type='income'
  // Set settings.store.hide_amounts = true
  // Expect cell-amount text to be '+¥*****'
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `npm test -- --testNamePattern="TransactionRow masks amount when setting enabled" -v`
  Expected: FAIL (still shows real amount)

- [ ] **Step 3: Implement updated formatAmount**
  Replace the formatAmount method (lines 84-86) with:
  ```ts
  import { formatAmountForDisplay } from '@/utils/amountDisplay'
  // ...
  const formatAmount = (amount: number) => {
    return formatAmountForDisplay(amount, props.transaction.type as 'income' | 'expense' | 'transfer')
  }
  ```

- [ ] **Step 4: Run test to verify it passes**
  Run: `npm test -- --testNamePattern="TransactionRow masks amount when setting enabled" -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add frontend/src/widgets/transactions/TransactionRow.vue
  git commit -m "transactions: use amount display utility in TransactionRow"
  ```

### Task 6: Frontend - Update SharedView to use amount display utility

**Files:**
- Modify: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/views/shared/SharedView.vue:25-35` (formatCurrency function and its usage)
- Also update the template where formatCurrency is called (line ~70)

**Interfaces:**
- Consumes: amountDisplay utility (but note: SharedView currently adds sign separately)
- Produces: formatted amount string (with sign handled by utility)

- [ ] **Step 1: Write the failing test**
  ```ts
  // Test that SharedView shows masked amount when hide_amounts=true
  // Mount SharedView with a transaction
  // Set settings.store.hide_amounts = true
  // Expect tx-amount text to be '+¥*****' or '-¥*****' (depending on type)
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `npm test -- --testNamePattern="SharedView masks amount when setting enabled" -v`
  Expected: FAIL (still shows real amount)

- [ ] **Step 3: Implement changes**
  - Remove the existing formatCurrency function (lines 25-35)
  - Import formatAmountForDisplay from '@/utils/amountDisplay'
  - In the template, replace:
        {{ tx.type === 'income' ? '+' : '-' }}{{ formatCurrency(tx.amount, data.account.currency) }}
    with:
        {{ formatAmountForDisplay(tx.amount, tx.type as 'income' | 'expense' | 'transfer') }}
    (Note: we need to cast tx.type to the literal type)

- [ ] **Step 4: Run test to verify it passes**
  Run: `npm test -- --testNamePattern="SharedView masks amount when setting enabled" -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add frontend/src/views/shared/SharedView.vue
  git commit -m "shared: use amount display utility for transaction amounts"
  ```

### Task 7: Frontend - Update ReportsView to use amount display utility

**Files:**
- Modify: `/data/data/com.termux/files/home/developments/kiroku/frontend/src/views/reports/ReportsView.vue` (multiple places: lines with formatNumber and chart data)

**Interfaces:**
- Consumes: amountDisplay utility
- Produces: formatted amount string for table and chart

- [ ] **Step 1: Write the failing test**
  ```ts
  // Test that ReportsView shows masked amount when hide_amounts=true
  // Mount ReportsView with spending report data
  // Set settings.store.hide_amounts = true
  // Expect table cells and chart tooltips to show masked amounts
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `npm test -- --testNamePattern="ReportsView masks amount when setting enabled" -v`
  Expected: FAIL (still shows real amounts)

- [ ] **Step 3: Implement changes**
  - Import formatAmountForDisplay
  - Replace formatNumber calls with formatAmountForDisplay (need to determine type - for spending report, likely expense)
  - Replace chart data mapping: instead of `d.amount`, use formatted string? But chart expects numbers. 
    Actually, the chart in ReportsView uses:
        data: data.map(d => ({ value: d.amount, itemStyle: { color: d.color } }))
    We cannot put strings in chart value. We need to hide the amount in the chart as well? 
    The requirement is to hide amounts in the UI - the chart is part of the UI. 
    We have two options:
      a) Hide the chart entirely when hide_amounts=true (maybe not desired)
      b) Show the chart with masked values? But chart needs numbers to draw.
    Looking at the PRD, the hide amount toggle is for privacy when showing screen. 
    The chart reveals amounts via visual length of bars. Hiding the numbers but showing bars still reveals amounts.
    Therefore, we should likely hide the chart when hide_amounts=true, or show a placeholder chart.
    However, to keep it simple and given the scope, let's assume we only hide the textual amounts and leave the chart as is (since the chart is not precise without numbers? Actually, the chart uses the actual amount for bar length).
    This is a design decision. We'll note in the plan that the chart may still reveal approximate amounts via bar length, and if full hiding is required, we would need to hide the chart or show dummy data.
    For now, we'll update only the textual amounts (table and tooltip formatter) and leave the chart data as numbers.
    The tooltip formatter uses: formatter: '¥{c}' -> we can change that to use our masked formatter.

  So:
    - In the table: replace ¥{{ formatNumber(child.amount) }} with formatted string
    - In the tooltip formatter: change to return formatted amount instead of '¥{c}'
    - In the chart data: keep d.amount (so bars reflect real amounts, but tooltip will be masked if we change formatter)

  However, the tooltip formatter is called with {c} which is the value (the amount). We can change the formatter to a function that returns the masked string.

  Let's do:
    - Tooltip formatter: change to a function that params => {
        const amount = params.value; // this is the number
        return formatAmountForDisplay(amount, 'expense'); // assuming expense for spending report
      }

  - For the chart data, we leave it as d.amount so the chart draws correctly.

- [ ] **Step 4: Run test to verify it passes**
  Run: `npm test -- --testNamePattern="ReportsView masks amount when setting enabled" -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add frontend/src/views/reports/ReportsView.vue
  git commit -m "reports: use amount display utility for amounts"
  ```

### Task 8: Frontend - Update Analytics widgets to use amount display utility

**Files:**
- Modify: Multiple widget files in `/data/data/com.termux/files/home/developments/kiroku/frontend/src/widgets/analytics/` that display amounts (e.g., BalanceTrendWidget.vue, CashflowTrendWidget.vue, etc.)

**Interfaces:**
- Consumes: amountDisplay utility
- Produces: formatted amount string for widget display

- [ ] **Step 1: Write the failing test**
  ```ts
  // Test that an analytics widget shows masked amount when hide_amounts=true
  // Mount BalanceTrendWidget with sample data
  // Set settings.store.hide_amounts = true
  // Expect amount text to be masked
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `npm test -- --testNamePattern="Analytics widget masks amount when setting enabled" -v`
  Expected: FAIL (still shows real amounts)

- [ ] **Step 3: Implement changes for each widget**
  For each widget file:
    - Import formatAmountForDisplay
    - Replace any direct amount formatting (toLocaleString, Intl.NumberFormat, etc.) with formatAmountForDisplay
    - Determine the appropriate type (income/expense/transfer) based on widget context

  We'll need to survey the widgets. Based on the widget names:
    - BalanceTrendWidget.vue: likely shows balance (net amount) -> type? we can use 'transfer' or determine from context (balance can be positive/negative). We'll pass type as 'transfer' and let the utility return sign based on amount? Actually, our utility requires the type to determine sign. 
      We can change the utility to accept an optional `showSign` boolean and compute sign from amount if not provided? 
      But to minimize changes, we'll compute the sign outside and pass type accordingly.
      For balance, we can do: type = amount >= 0 ? 'income' : 'expense' (but note: balance is not income/expense, it's net). 
      Alternatively, we can change the utility to take a `sign` parameter or compute sign from amount.
      Let's adjust the utility: instead of requiring type, we can have an optional `signOverride` (string: '+', '-', or '') and if provided, use that; else compute from type.
      However, to keep the utility simple and avoid breaking existing uses, we'll create two utilities or extend the current one.

      Given the time, let's change the utility to accept an optional `sign` parameter (string) and if provided, use that; else compute from type.
      We'll update the utility in Task 3.

  We'll do this in Task 3 step 3: modify the utility to:
    export function formatAmountForDisplay(
      amount: number, 
      type: 'income' | 'expense' | 'transfer' = 'transfer',
      signOverride?: '+' | '-' | ''
    ): string {
      const settingsStore = useSettingsStore()
      if (settingsStore.settings.hide_amounts) {
        const sign = signOverride ?? (type === 'income' ? '+' : type === 'expense' ? '-' : '')
        return `${sign}¥*****`
      }
      const sign = signOverride ?? (type === 'income' ? '+' : type === 'expense' ? '-' : '')
      return `${sign}¥${amount.toLocaleString()}`
    }

  Then for balance, we can compute signOverride = amount >= 0 ? '+' : '-'

- [ ] **Step 4: Run test to verify it passes**
  Run: `npm test -- --testNamePattern="Analytics widget masks amount when setting enabled" -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add frontend/src/widgets/analytics/BalanceTrendWidget.vue frontend/src/widgets/analytics/CashflowTrendWidget.vue /* ... */ 
  git commit -m "analytics: use amount display utility for widgets"
  ```

### Task 9: Frontend - Update any remaining amount displays (e.g., dropdowns, modals)

**Files:**
- Modify: Any remaining Vue files that display amounts (e.g., TransactionFormModal? but we decided not to mask inputs, only display)

**Interfaces:**
- Consumes: amountDisplay utility
- Produces: formatted amount string

- [ ] **Step 1: Write the failing test**
  ```ts
  // Test that no remaining places show real amounts when hide_amounts=true
  // Traverse all components and check for plain amount display
  ```

- [ ] **Step 2: Run test to verify it fails**
  Run: `npm test -- --testNamePattern="no raw amounts displayed when hide amounts enabled" -v`
  Expected: FAIL if any raw amounts found

- [ ] **Step 3: Implement changes**
  - Use the same pattern: import utility and replace formatting calls
  - Pay attention to places like:
        - Dropdown options that show amounts (unlikely)
        - Empty states
        - Tooltips
        - Skeleton loaders (should not show amounts)

- [ ] **Step 4: Run test to verify it passes**
  Run: `npm test -- --testNamePattern="no raw amounts displayed when hide amounts enabled" -v`
  Expected: PASS

- [ ] **Step 5: Commit**
  ```bash
  git add frontend/src/components/forms/TransactionFormModal.vue /* ... other files ... */
  git commit -m "ui: apply amount display utility to remaining amount displays"
  ```

### Task 10: End-to-end verification

**Files:**
- None (verification tasks)

**Interfaces:**
- Consumes: All implemented changes
- Produces: Confirmed working feature

- [ ] **Step 1: Verify setting persists**
  - Toggle hide amounts to true
  - Reload page
  - Confirm setting remains true and amounts are masked

- [ ] **Step 2: Verify masking works across views**
  - Navigate to Transactions, Reports, Analytics, Shared view
  - Confirm all amounts are masked (showing ***** with appropriate sign)

- [ ] **Step 3: Verify unmasked when disabled**
  - Toggle hide amounts to false
  - Confirm real amounts reappear

- [ ] **Step 4: Verify form inputs unaffected**
  - Open transaction form
  - Confirm amount input field allows editing and shows real numbers

- [ ] **Step 5: Commit**
  ```bash
  git commit -m "feat: hide amounts toggle implementation complete"
  ```
