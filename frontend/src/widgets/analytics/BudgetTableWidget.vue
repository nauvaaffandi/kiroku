<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>カテゴリ別 予算 vs 実績</h2>
      <button @click="$emit('export')" class="btn-secondary" style="padding: 0.375rem 0.75rem; font-size: 0.75rem;">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
          <polyline points="7 10 12 15 17 10"></polyline>
          <line x1="12" y1="15" x2="12" y2="3"></line>
        </svg>
        CSV
      </button>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="budget-table-container">
        <table class="budget-table">
          <thead>
            <tr>
              <th>カテゴリ</th>
              <th class="text-right">予算</th>
              <th class="text-right">実績</th>
              <th class="text-right">差額</th>
              <th class="text-right">%</th>
              <th>ステータス</th>
              <th>進捗バー</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="budget in budgets" :key="budget.category_id">
              <td>
                <div class="category-cell">
                  <span class="category-dot" :style="{ backgroundColor: budget.color }"></span>
                  <span>{{ budget.icon }} {{ budget.category_name }}</span>
                </div>
              </td>
              <td class="text-right">{{ formatCurrency(budget.budget_amount) }}</td>
              <td class="text-right">{{ formatCurrency(budget.actual_amount) }}</td>
              <td class="text-right" :class="budget.variance >= 0 ? 'positive' : 'negative'">
                {{ budget.variance >= 0 ? '+' : '' }}{{ formatCurrency(budget.variance) }}
              </td>
              <td class="text-right" :class="budget.variance_pct >= 0 ? 'positive' : 'negative'">
                {{ budget.variance_pct >= 0 ? '+' : '' }}{{ budget.variance_pct.toFixed(1) }}%
              </td>
              <td>
                <span class="status-badge" :class="budget.status">
                  {{ statusLabel(budget.status) }}
                </span>
              </td>
              <td>
                <div class="budget-bar-container">
                  <div class="budget-bar-bg">
                    <div 
                      class="budget-bar-fill" 
                      :style="{ width: barWidth(budget), backgroundColor: barColor(budget) }"
                    />
                  </div>
                  <div class="budget-bar-target" :style="{ left: targetPosition(budget) + '%' }" />
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div v-else class="widget-content">
      <div class="skeleton skeleton-table-row" v-for="i in 5" :key="i" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface BudgetItem {
  category_id: string
  category_name: string
  icon: string
  color: string
  period: string
  budget_amount: number
  actual_amount: number
  variance: number
  variance_pct: number
  status: 'under' | 'over' | 'on_track'
}

defineProps<{
  loading: boolean
  budgets: BudgetItem[]
}>()

defineEmits<{
  export: []
}>()

function statusLabel(status: string): string {
  const labels: Record<string, string> = {
    under: '未達',
    over: '超過',
    on_track: '順調'
  }
  return labels[status] || status
}

function barWidth(budget: BudgetItem): string {
  const pct = budget.budget_amount > 0 ? (budget.actual_amount / budget.budget_amount) * 100 : 0
  return Math.min(pct, 100) + '%'
}

function targetPosition(_budget: BudgetItem): string {
  return '100%'
}

function barColor(budget: BudgetItem): string {
  if (budget.status === 'over') return 'var(--error)'
  if (budget.status === 'under') return 'var(--warning)'
  return 'var(--success)'
}

function formatCurrency(val: number): string {
  return new Intl.NumberFormat('ja-JP', { 
    style: 'currency', 
    currency: 'JPY',
    maximumFractionDigits: 0 
  }).format(val)
}
</script>

<style scoped>
@reference "tailwindcss";

.widget-card {
  @apply border rounded-xl overflow-hidden;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}
.widget-header {
  @apply flex items-center justify-between p-4 border-b;
  background-color: var(--bg-secondary);
  border-color: var(--border-primary);
}
.widget-header h2 {
  @apply text-base font-semibold;
  color: var(--text-primary);
}
.widget-content {
  @apply p-4;
}
.budget-table-container {
  @apply overflow-x-auto;
}
.budget-table {
  @apply w-full border-collapse;
  font-size: 0.8125rem;
}
.budget-table th {
  @apply text-left py-3 px-2 font-semibold border-b-2 whitespace-nowrap;
  color: var(--text-secondary);
  border-color: var(--border-primary);
}
.budget-table td {
  @apply py-3 px-2 border-b whitespace-nowrap;
  color: var(--text-primary);
  border-color: var(--border-primary);
}
.budget-table td.text-right {
  @apply text-right font-mono text-xs;
}
.budget-table td.positive {
  color: var(--success);
}
.budget-table td.negative {
  color: var(--error);
}
.budget-table tbody tr:hover {
  background-color: var(--bg-hover);
}
.category-cell {
  @apply flex items-center gap-2;
}
.category-dot {
  @apply w-[0.625rem] h-[0.625rem] rounded-full shrink-0;
}
.status-badge {
  @apply inline-block px-2 py-[0.125rem] rounded-full text-[0.6875rem] font-medium uppercase;
}
.status-badge.under {
  background-color: var(--warning-light);
  color: var(--warning);
}
.status-badge.over {
  background-color: var(--error-light);
  color: var(--error);
}
.status-badge.on_track {
  background-color: var(--success-light);
  color: var(--success);
}
.budget-bar-container {
  @apply relative;
  width: 200px;
  height: 1.5rem;
}
.budget-bar-bg {
  @apply w-full h-full rounded-full overflow-hidden;
  background-color: var(--bg-tertiary);
}
.budget-bar-fill {
  @apply h-full rounded-full;
  transition: width 0.3s ease;
}
.budget-bar-target {
  @apply absolute w-[2px] z-[1];
  top: -2px;
  bottom: -2px;
  background-color: var(--text-primary);
}
.btn-secondary {
  @apply inline-flex items-center gap-[0.375rem] border rounded-sm font-medium;
  border-color: var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
  transition: background-color 0.2s ease;
}
.btn-secondary:hover {
  background-color: var(--bg-hover);
}
</style>
