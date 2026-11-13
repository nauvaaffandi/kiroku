<template>
  <div class="analytics-cashflow-view">
    <div class="page-header">
      <h1>{{ t('analytics.cashflow.title') }}</h1>
      <p class="page-subtitle">{{ t('analytics.cashflow.subtitle') }}</p>
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
import { ref, onMounted } from 'vue'
import { useI18n } from '@/locales'
import { localDb } from '@/storage/repository'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import CashflowWaterfallWidget from '@/widgets/analytics/CashflowWaterfallWidget.vue'
import CashflowTrendWidget from '@/widgets/analytics/CashflowTrendWidget.vue'
import CashflowComparisonWidget from '@/widgets/analytics/CashflowComparisonWidget.vue'
import ComparisonPeriodPickerWidget from '@/widgets/analytics/ComparisonPeriodPickerWidget.vue'
import type { Transaction } from '@/types'

const { t } = useI18n()

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

const flowData = ref<any>(null)
const trendData = ref<any>([])
const comparisonData = ref<any>(null)

function onDateChange(payload: { dateFrom: string; dateTo: string }) {
  dateFrom.value = payload.dateFrom
  dateTo.value = payload.dateTo
  fetchFlow()
  fetchCashFlowTrend()
}

async function fetchFlow() {
  flowLoading.value = true
  try {
    const transactions = await localDb.list<Transaction>('transaction')
    const inRange = transactions.filter((t) => {
      const d = t.date_time.slice(0, 10)
      return (!dateFrom.value || d >= dateFrom.value) && (!dateTo.value || d <= dateTo.value)
    })
    const opening_balance = 0
    const total_income = inRange.filter((t) => t.type === 'income').reduce((s, t) => s + t.amount, 0)
    const total_expense = inRange.filter((t) => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
    const net_flow = total_income - total_expense
    const closing_balance = opening_balance + net_flow
    const byCategory = inRange.reduce((acc, t) => {
      if (t.type === 'expense') {
        const key = t.category_id ?? 'uncategorized'
        if (!acc[key]) acc[key] = { category_id: key, name: t.category_id, income: 0, expense: 0, net: 0 }
        acc[key].expense += t.amount
        acc[key].net -= t.amount
      } else {
        const key = t.category_id ?? 'uncategorized'
        if (!acc[key]) acc[key] = { category_id: key, name: t.category_id, income: 0, expense: 0, net: 0 }
        acc[key].income += t.amount
        acc[key].net += t.amount
      }
      return acc
    }, {} as Record<string, any>)

    flowData.value = {
      opening_balance,
      total_income,
      total_expense,
      net_flow,
      closing_balance,
      by_category: Object.values(byCategory)
    }
  } catch (e) {
    console.error('Failed to compute cash flow:', e)
  } finally {
    flowLoading.value = false
  }
}

async function fetchCashFlowTrend() {
  trendLoading.value = true
  try {
    const transactions = await localDb.list<Transaction>('transaction')
    const now = new Date()
    const months: string[] = []
    const count = trendInterval.value === 'year' ? 12 : 6
    for (let i = count - 1; i >= 0; i--) {
      const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
      months.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
    }
    const series: any[] = months.map((period) => {
      const inMonth = transactions.filter((t) => t.date_time.slice(0, 7) === period)
      const income = inMonth.filter((t) => t.type === 'income').reduce((s, t) => s + t.amount, 0)
      const expense = inMonth.filter((t) => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
      return {
        period,
        income,
        expense,
        net: income - expense,
        cumulative_net: 0
      }
    })
    trendData.value = series
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
    const transA = (await localDb.list<Transaction>('transaction')).filter((t) => {
      const d = t.date_time.slice(0, 10)
      return (!periodAFrom.value || d >= periodAFrom.value) && (!periodATo.value || d <= periodATo.value)
    })
    const transB = (await localDb.list<Transaction>('transaction')).filter((t) => {
      const d = t.date_time.slice(0, 10)
      return (!periodBFrom.value || d >= periodBFrom.value) && (!periodBTo.value || d <= periodBTo.value)
    })
    const compute = (transactions: Transaction[]) => ({
      income: transactions.filter((t) => t.type === 'income').reduce((s, t) => s + t.amount, 0),
      expense: transactions.filter((t) => t.type === 'expense').reduce((s, t) => s + t.amount, 0),
      net: transactions.filter((t) => t.type === 'income').reduce((s, t) => s + t.amount, 0) - transactions.filter((t) => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
    })
    const pa = compute(transA)
    const pb = compute(transB)
    comparisonData.value = {
      period_a: {
        income: pa.income,
        expense: pa.expense,
        net: pa.net,
        avg_daily_income: pa.income / Math.max(1, periodAFrom.value ? new Date(periodATo.value).getDate() - new Date(periodAFrom.value).getDate() + 1 : 1),
        avg_daily_expense: pa.expense / Math.max(1, periodAFrom.value ? new Date(periodATo.value).getDate() - new Date(periodAFrom.value).getDate() + 1 : 1),
        transaction_count: transA.length
      },
      period_b: {
        income: pb.income,
        expense: pb.expense,
        net: pb.net,
        avg_daily_income: pb.income / Math.max(1, periodBFrom.value ? new Date(periodBTo.value).getDate() - new Date(periodBFrom.value).getDate() + 1 : 1),
        avg_daily_expense: pb.expense / Math.max(1, periodBFrom.value ? new Date(periodBTo.value).getDate() - new Date(periodBFrom.value).getDate() + 1 : 1),
        transaction_count: transB.length
      },
      delta: {
        income_pct: pa.income === 0 ? 0 : (pb.income - pa.income) / pa.income * 100,
        expense_pct: pa.expense === 0 ? 0 : (pb.expense - pa.expense) / pa.expense * 100,
        net_pct: pa.net === 0 ? 0 : (pb.net - pa.net) / Math.abs(pa.net) * 100
      }
    }
  } catch (e) {
    console.error('Failed to fetch cash flow comparison:', e)
  } finally {
    comparisonLoading.value = false
  }
}

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