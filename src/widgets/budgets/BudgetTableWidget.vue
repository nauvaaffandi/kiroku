<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>{{ t('budgets.table.title') }}</h2>
      <button @click="$emit('export')" class="btn-secondary">CSV</button>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="budget-table-container">
        <table class="budget-table">
          <thead>
            <tr>
              <th>{{ t('budgets.table.name') }}</th>
              <th>{{ t('budgets.table.period') }}</th>
              <th class="text-right">{{ t('budgets.table.budget') }}</th>
              <th class="text-right">{{ t('budgets.table.spent') }}</th>
              <th class="text-right">{{ t('budgets.table.remaining') }}</th>
              <th class="text-right">{{ t('budgets.table.usage') }}</th>
              <th>{{ t('budgets.table.status') }}</th>
              <th>{{ t('budgets.table.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="budget in budgets" :key="budget.id">
              <td>{{ budget.name }}</td>
              <td><span class="period-badge" :class="budget.period">{{ t(`budgets.period.${budget.period}`) }}</span></td>
              <td class="text-right">{{ formatCurrency(budget.amount) }}</td>
               <td class="text-right">{{ formatCurrency(budget.spent) }}</td>
               <td class="text-right" :class="budget.remaining < 0 ? 'negative' : ''">{{ formatCurrency(budget.remaining) }}</td>
              <td class="text-right">{{ budget.percentageUsed }}%</td>
              <td><span class="status-badge" :class="budget.status">{{ statusLabel(budget.status) }}</span></td>
              <td class="text-right">
                <button @click="$emit('edit', budget.id)" class="icon-btn edit" aria-label="Edit">{{ t('common.edit') }}</button>
                <button @click="$emit('delete', budget)" class="icon-btn delete" aria-label="Delete">{{ t('common.delete') }}</button>
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
import { useI18n } from '@/locales'
import { formatAmountForDisplay } from '@/utils/amountDisplay'

interface BudgetRow {
  id: string
  name: string
  period: string
  amount: number
  currency: string
  spent: number
  remaining: number
  percentageUsed: number
  isOverBudget: boolean
  status: 'over' | 'under'
}

const { t } = useI18n()

defineProps<{
  loading: boolean
  budgets: BudgetRow[]
}>()

defineEmits<{
  export: []
  edit: [id: string]
  delete: [row: BudgetRow]
}>()

function statusLabel(status: string): string {
  return t(`budgets.status.${status}`)
}

function formatCurrency(val: number): string {
  return formatAmountForDisplay(val, 'transfer')
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
.budget-table td.negative {
  color: var(--error);
}
.budget-table tbody tr:hover {
  background-color: var(--bg-hover);
}
.period-badge {
  @apply inline-block px-2 py-[0.125rem] rounded-full text-[0.6875rem] font-medium uppercase;
}
.period-badge.weekly { background-color: var(--info-light); color: var(--info); }
.period-badge.monthly { background-color: var(--accent-light); color: var(--accent); }
.period-badge.yearly { background-color: var(--warning-light); color: var(--warning); }
.period-badge.one_time { background-color: var(--text-muted); color: var(--text-secondary); }
.status-badge {
  @apply inline-block px-2 py-[0.125rem] rounded-full text-[0.6875rem] font-medium uppercase;
}
.status-badge.under { background-color: var(--success-light); color: var(--success); }
.status-badge.over { background-color: var(--error-light); color: var(--error); }
.icon-btn {
  @apply inline-flex items-center justify-center px-2 py-1 rounded-sm border font-medium text-xs no-underline;
  border-color: var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
  margin-left: 0.25rem;
}
.icon-btn:hover { background-color: var(--bg-hover); }
.icon-btn.edit { color: var(--accent); border-color: var(--accent-light); }
.icon-btn.delete { color: var(--error); border-color: var(--error-light); }
.btn-secondary {
  @apply inline-flex items-center gap-[0.375rem] border rounded-sm font-medium text-sm;
  border-color: var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}
.skeleton-table-row { @apply h-10 animate-pulse bg-zinc-200 dark:bg-zinc-800 rounded; }
</style>
