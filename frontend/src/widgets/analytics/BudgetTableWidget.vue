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
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; background-color: var(--bg-secondary); border-bottom: 1px solid var(--border-primary); }
.widget-header h2 { font-size: 1rem; font-weight: 600; color: var(--text-primary); }
.widget-content { padding: 1rem; }
.budget-table-container { overflow-x: auto; }
.budget-table { width: 100%; border-collapse: collapse; font-size: 0.8125rem; }
.budget-table th { padding: 0.75rem 0.5rem; text-align: left; font-weight: 600; color: var(--text-secondary); border-bottom: 2px solid var(--border-primary); white-space: nowrap; }
.budget-table td { padding: 0.75rem 0.5rem; border-bottom: 1px solid var(--border-primary); color: var(--text-primary); white-space: nowrap; }
.budget-table td.text-right { text-align: right; font-family: monospace; font-size: 0.75rem; }
.budget-table td.positive { color: var(--success); }
.budget-table td.negative { color: var(--error); }
.budget-table tbody tr:hover { background-color: var(--bg-hover); }
.category-cell { display: flex; align-items: center; gap: 0.5rem; }
.category-dot { width: 0.625rem; height: 0.625rem; border-radius: 50%; flex-shrink: 0; }
.status-badge { display: inline-block; padding: 0.125rem 0.5rem; border-radius: 9999px; font-size: 0.6875rem; font-weight: 500; text-transform: uppercase; }
.status-badge.under { background-color: var(--warning-bg); color: var(--warning); }
.status-badge.over { background-color: var(--error-bg); color: var(--error); }
.status-badge.on_track { background-color: var(--success-bg); color: var(--success); }
.budget-bar-container { position: relative; width: 200px; height: 1.5rem; }
.budget-bar-bg { width: 100%; height: 100%; background-color: var(--bg-tertiary); border-radius: 9999px; overflow: hidden; }
.budget-bar-fill { height: 100%; border-radius: 9999px; transition: width 0.3s ease; }
.budget-bar-target { position: absolute; top: -2px; bottom: -2px; width: 2px; background-color: var(--text-primary); z-index: 1; }
.btn-secondary { display: inline-flex; align-items: center; gap: 0.375rem; border: 1px solid var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); border-radius: var(--radius-sm); font-weight: 500; transition: background-color 0.2s ease; }
.btn-secondary:hover { background-color: var(--bg-hover); }
</style>