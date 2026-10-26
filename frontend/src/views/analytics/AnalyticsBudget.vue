<template>
  <div class="analytics-budget-view">
    <div class="page-header">
      <h1>{{ t('budgetAnalysis.title') }}</h1>
      <p class="page-subtitle">{{ t('budgetAnalysis.subtitle') }}</p>
    </div>

    <DateFilterBar @change="onDateChange" />

    <div class="controls-bar">
      <select v-model="budgetPeriod" @change="fetchData" class="period-select">
        <option value="weekly">{{ t('budgetAnalysis.weekly') }}</option>
        <option value="monthly">{{ t('budgetAnalysis.monthly') }}</option>
        <option value="yearly">{{ t('budgetAnalysis.yearly') }}</option>
      </select>
    </div>

    <BudgetSummaryCardsWidget :loading="loading" :summary="summary" />

    <BudgetTableWidget
      :loading="loading"
      :budgets="budgets"
      @export="exportCSV"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watchEffect } from 'vue'
import { api } from '@/api/client'
import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import BudgetSummaryCardsWidget from '@/widgets/analytics/BudgetSummaryCardsWidget.vue'
import BudgetTableWidget from '@/widgets/analytics/BudgetTableWidget.vue'

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

interface Summary {
  total_budget: number
  total_actual: number
  total_variance: number
  categories_over_budget: number
  categories_under_budget: number
}

const props = defineProps<{
  dateFrom: string
  dateTo: string
}>()

const uiStore = useUIStore()
const { t } = useI18n()
const loading = ref(true)
const budgetPeriod = ref('monthly')

const budgets = ref<BudgetItem[]>([])
const summary = ref<Summary>({
  total_budget: 0,
  total_actual: 0,
  total_variance: 0,
  categories_over_budget: 0,
  categories_under_budget: 0,
})

function onDateChange(_payload: { dateFrom: string; dateTo: string }) {
  // Parent handles date changes
}

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get<{ budgets: BudgetItem[]; summary: Summary }>(
      `/analytics/budget/analysis?date_from=${props.dateFrom}&date_to=${props.dateTo}&budget_period=${budgetPeriod.value}`
    )
    budgets.value = res.data.budgets
    summary.value = res.data.summary
  } catch (e) {
    console.error('Failed to fetch budget analysis:', e)
    uiStore.error('Failed to fetch budget analysis')
  } finally {
    loading.value = false
  }
}

function exportCSV() {
  function formatCurrency(val: number): string {
    return new Intl.NumberFormat('ja-JP', { 
      style: 'currency', 
      currency: 'JPY',
      maximumFractionDigits: 0 
    }).format(val)
  }
  function statusLabel(status: string): string {
    const labels: Record<string, string> = { under: 'Below Budget', over: 'Over Budget', on_track: 'On Track' }
    return labels[status] || status
  }
  const headers = ['Category', 'Budget', 'Actual', 'Variance', 'Variance%', 'Status']
  const rows = budgets.value.map(b => [
    b.category_name,
    formatCurrency(b.budget_amount),
    formatCurrency(b.actual_amount),
    formatCurrency(b.variance),
    b.variance_pct.toFixed(1) + '%',
    statusLabel(b.status)
  ])
  rows.push(['合計', formatCurrency(summary.value.total_budget), formatCurrency(summary.value.total_actual), formatCurrency(summary.value.total_variance), '', ''])

  const csv = [headers, ...rows].map(r => r.map(v => `"${v}"`).join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `budget-analysis-${props.dateFrom}-${props.dateTo}.csv`
  link.click()
}

onMounted(() => {
  fetchData()
})

watchEffect(() => {
  fetchData()
})
</script>

<style scoped>
@reference "tailwindcss";
.analytics-budget-view {
  @apply flex flex-col gap-4 p-6;
}

.page-header {
  @apply pb-4;
  border-bottom: 1px solid var(--border-primary);
}

.page-header h1 {
  @apply text-xl font-semibold;
  color: var(--text-primary);
}

.page-subtitle {
  @apply text-[0.875rem] mt-1;
  color: var(--text-secondary);
}

.controls-bar {
  @apply flex gap-4 items-center p-4 rounded-xl;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
}

.period-select {
  @apply px-4 py-2 rounded-sm text-sm;
  border: 1px solid var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}
</style>