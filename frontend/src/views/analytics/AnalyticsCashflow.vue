<template>
  <div class="analytics-cashflow-view">
    <div class="page-header">
      <h1>キャッシュフロー分析</h1>
      <p class="page-subtitle">キャッシュフローのウォーターフォール・トレンド・期間比較</p>
    </div>

    <div class="filter-section">
      <DateFilterBar @change="onDateChange" />
      <ComparisonPeriodPickerWidget
        :showComparisonPeriods="showComparisonPeriods"
        v-model:periodAFrom="periodAFrom"
        v-model:periodATo="periodATo"
        v-model:periodBFrom="periodBFrom"
        v-model:periodBTo="periodBTo"
      />
    </div>

    <div class="widgets-grid">
      <CashflowWaterfallWidget :loading="flowLoading" :data="flowData" />

      <CashflowTrendWidget
        v-model:interval="trendInterval"
        :loading="trendLoading"
        :data="trendData"
        @update:interval="fetchCashFlowTrend"
      />

      <CashflowComparisonWidget :loading="comparisonLoading" :data="comparisonData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { api } from '@/api/client'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import CashflowWaterfallWidget from '@/widgets/analytics/CashflowWaterfallWidget.vue'
import CashflowTrendWidget from '@/widgets/analytics/CashflowTrendWidget.vue'
import CashflowComparisonWidget from '@/widgets/analytics/CashflowComparisonWidget.vue'
import ComparisonPeriodPickerWidget from '@/widgets/analytics/ComparisonPeriodPickerWidget.vue'

interface FlowResponse {
  opening_balance: number
  total_income: number
  total_expense: number
  net_flow: number
  closing_balance: number
  by_category: Array<{
    category_id: string
    name: string
    income: number
    expense: number
    net: number
  }>
}

interface TrendPoint {
  period: string
  income: number
  expense: number
  net: number
  cumulative_net: number
}

interface ComparisonResponse {
  period_a: {
    income: number
    expense: number
    net: number
    avg_daily_income: number
    avg_daily_expense: number
    transaction_count: number
  }
  period_b: {
    income: number
    expense: number
    net: number
    avg_daily_income: number
    avg_daily_expense: number
    transaction_count: number
  }
  delta: {
    income_pct: number
    expense_pct: number
    net_pct: number
  }
}

const dateFrom = ref('')
const dateTo = ref('')
const showComparisonPeriods = ref(false)
const periodAFrom = ref('')
const periodATo = ref('')
const periodBFrom = ref('')
const periodBTo = ref('')
const trendInterval = ref('month')

const flowLoading = ref(true)
const trendLoading = ref(true)
const comparisonLoading = ref(true)

const flowData = ref<FlowResponse | null>(null)
const trendData = ref<TrendPoint[]>([])
const comparisonData = ref<ComparisonResponse | null>(null)

function onDateChange(payload: { dateFrom: string; dateTo: string }) {
  dateFrom.value = payload.dateFrom
  dateTo.value = payload.dateTo
  fetchFlow()
  fetchCashFlowTrend()
}

async function fetchFlow() {
  flowLoading.value = true
  try {
    const res = await api.get<FlowResponse>(
      `/analytics/cashflow/flow?date_from=${dateFrom.value}&date_to=${dateTo.value}`
    )
    flowData.value = res.data
  } catch (e) {
    console.error('Failed to fetch cash flow:', e)
  } finally {
    flowLoading.value = false
  }
}

async function fetchCashFlowTrend() {
  trendLoading.value = true
  try {
    const res = await api.get<{ series: TrendPoint[] }>(
      `/analytics/cashflow/trend?date_from=${dateFrom.value}&date_to=${dateTo.value}&interval=${trendInterval.value}`
    )
    trendData.value = res.data.series
  } catch (e) {
    console.error('Failed to fetch cash flow trend:', e)
  } finally {
    trendLoading.value = false
  }
}

async function fetchComparison() {
  if (!periodAFrom.value || !periodATo.value || !periodBFrom.value || !periodBTo.value) return
  comparisonLoading.value = true
  try {
    const res = await api.get<ComparisonResponse>(
      `/analytics/cashflow/comparison?period_a_from=${periodAFrom.value}&period_a_to=${periodATo.value}&period_b_from=${periodBFrom.value}&period_b_to=${periodBTo.value}`
    )
    comparisonData.value = res.data
  } catch (e) {
    console.error('Failed to fetch cash flow comparison:', e)
  } finally {
    comparisonLoading.value = false
  }
}

watch(() => [periodAFrom.value, periodATo.value, periodBFrom.value, periodBTo.value], () => {
  fetchComparison()
})

onMounted(async () => {
  const now = new Date()
  dateTo.value = now.toISOString().split('T')[0]
  const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1)
  dateFrom.value = lastMonth.toISOString().split('T')[0]

  periodAFrom.value = new Date(now.getFullYear(), now.getMonth() - 2, 1).toISOString().split('T')[0]
  periodATo.value = new Date(now.getFullYear(), now.getMonth() - 1, 0).toISOString().split('T')[0]
  periodBFrom.value = new Date(now.getFullYear(), now.getMonth() - 1, 1).toISOString().split('T')[0]
  periodBTo.value = new Date(now.getFullYear(), now.getMonth(), 0).toISOString().split('T')[0]

  await Promise.all([fetchFlow(), fetchCashFlowTrend()])
  fetchComparison()
})
</script>

<style scoped>
@reference "tailwindcss";
.analytics-cashflow-view {
  @apply max-w-[80rem] mx-auto p-4 flex flex-col gap-4;
}
.page-header { 
  @apply flex flex-col gap-1; 
}
.page-header h1 { 
  @apply text-xl font-bold; 
  color: var(--text-primary); 
}
.page-subtitle { 
  color: var(--text-muted); 
}
.filter-section { 
  @apply flex flex-col gap-4; 
}
.widgets-grid { 
  @apply grid gap-4; 
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); 
}
</style>