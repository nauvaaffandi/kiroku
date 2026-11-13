<template>
  <div class="analytics-outlook-view">
    <div class="page-header">
      <h1>{{ t('outlook.title') }}</h1>
      <p class="page-subtitle">{{ t('outlook.subtitle') }}</p>
    </div>

    <DateFilterBar @change="onDateChange" />

    <div class="widgets-grid">
      <PlannedByCategoryWidget
        :loading="categoryLoading"
        :data="categoryData"
      />

      <PlannedByTimelineWidget
        :loading="timelineLoading"
        :data="timelineData"
        v-model:interval="timelineInterval"
        @update:interval="fetchPlannedByTimeline"
      />

      <BalanceForecastWidget
        :loading="forecastLoading"
        :data="forecastData"
        v-model:selectedConfidence="selectedConfidence"
        :confidenceLevels="confidenceLevels"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from '@/locales'
import { localDb } from '@/storage/repository'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import PlannedByCategoryWidget from '@/widgets/analytics/PlannedByCategoryWidget.vue'
import PlannedByTimelineWidget from '@/widgets/analytics/PlannedByTimelineWidget.vue'
import BalanceForecastWidget from '@/widgets/analytics/BalanceForecastWidget.vue'

const { t } = useI18n()

const dateFrom = ref('')
const dateTo = ref('')
const timelineInterval = ref('month')

const selectedConfidence = ref(['high', 'medium', 'low'])
const confidenceLevels = ['high', 'medium', 'low']

const categoryLoading = ref(true)
const timelineLoading = ref(true)
const forecastLoading = ref(true)

const categoryData = ref<any[]>([])
const timelineData = ref<any[]>([])
const forecastData = ref<any[]>([])

function onDateChange(payload: { dateFrom: string; dateTo: string }) {
  dateFrom.value = payload.dateFrom
  dateTo.value = payload.dateTo
  fetchAll()
}

async function fetchAll() {
  await Promise.all([
    fetchPlannedByCategory(),
    fetchPlannedByTimeline(),
    fetchBalanceForecast(),
  ])
}

async function fetchPlannedByCategory() {
  categoryLoading.value = true
  try {
    const categories = await localDb.list<any>('category')
    const allTransactions = await localDb.list<any>('transaction')
    const byCategory = new Map<string, any>()
    for (const t of allTransactions) {
      const cat = categories.find((c: any) => c.id === t.category_id)
      const name = cat?.name ?? 'Uncategorized'
      const entry = byCategory.get(name) ?? { id: cat?.id ?? '', name, color: cat?.color ?? '', kind: cat?.kind ?? 'uncategorized', planned_income: 0, planned_expense: 0, net: 0, transaction_count: 0 }
      entry.planned_income += t.amount
      entry.planned_expense += t.amount
      entry.net = entry.planned_income - entry.planned_expense
      entry.transaction_count += 1
      byCategory.set(name, entry)
    }
    categoryData.value = [...byCategory.values()]
  } catch (e) {
    console.error('Failed to fetch planned by category:', e)
  } finally {
    categoryLoading.value = false
  }
}

async function fetchPlannedByTimeline() {
  timelineLoading.value = true
  try {
    const allTransactions = await localDb.list<any>('transaction')
    const months: string[] = []
    const count = timelineInterval.value === 'year' ? 12 : 6
    for (let i = count - 1; i >= 0; i--) {
      const d = new Date(new Date().getFullYear(), new Date().getMonth() - i, 1)
      months.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
    }
    const series = months.map((period) => {
      const inMonth = allTransactions.filter((t) => t.date_time.slice(0, 7) === period)
      const income = inMonth.filter((t) => t.type === 'income').reduce((s, t) => s + t.amount, 0)
      const expense = inMonth.filter((t) => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
      return {
        period,
        planned_income: income,
        planned_expense: expense,
        net: income - expense,
        transaction_count: inMonth.length
      }
    })
    timelineData.value = series
  } catch (e) {
    console.error('Failed to fetch planned by timeline:', e)
  } finally {
    timelineLoading.value = false
  }
}

async function fetchBalanceForecast() {
  forecastLoading.value = true
  try {
    const allTransactions = await localDb.list<any>('transaction')
    const series = []
    for (let i = 0; i < 12; i++) {
      const d = new Date(new Date().getFullYear(), new Date().getMonth() + i, 1)
      const date = d.toISOString().split('T')[0]
      const inMonth = allTransactions.filter((t) => t.date_time.slice(0, 10) === date)
      const cumulative = inMonth.reduce((s, t) => t.type === 'income' ? s + t.amount : t.type === 'expense' ? s - t.amount : s, 0)
      series.push({
        date,
        projected_balance: cumulative,
        confidence_level: i < 6 ? 'high' : 'medium'
      })
    }
    forecastData.value = series
  } catch (e) {
    console.error('Failed to fetch balance forecast:', e)
  } finally {
    forecastLoading.value = false
  }
}

onMounted(async () => {
  const now = new Date()
  dateTo.value = new Date(now.getFullYear(), now.getMonth() + 3, 0).toISOString().split('T')[0]
  dateFrom.value = now.toISOString().split('T')[0]
  await fetchAll()
})
</script>

<style scoped>
@reference "tailwindcss";
.analytics-outlook-view {
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
.widgets-grid { 
  @apply grid gap-4; 
  grid-template-columns: 1fr;
}
@media (min-width: 768px) {
  .widgets-grid {
    grid-template-columns: minmax(400px, 1fr);
  }
}
</style>