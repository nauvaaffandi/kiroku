<template>
  <div class="analytics-outlook-view">
    <div class="page-header">
      <h1>Outlook 分析</h1>
      <p class="page-subtitle">予定支払いのカテゴリ別・タイムライン別・残高予測</p>
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
import { api } from '@/api/client'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import PlannedByCategoryWidget from '@/widgets/analytics/PlannedByCategoryWidget.vue'
import PlannedByTimelineWidget from '@/widgets/analytics/PlannedByTimelineWidget.vue'
import BalanceForecastWidget from '@/widgets/analytics/BalanceForecastWidget.vue'

interface PlannedCategory {
  id: string
  name: string
  color: string
  kind: string
  planned_income: number
  planned_expense: number
  net: number
  transaction_count: number
}

interface TimelinePoint {
  period: string
  planned_income: number
  planned_expense: number
  net: number
  transaction_count: number
}

interface ForecastPoint {
  date: string
  projected_balance: number
  confidence_level: string
}

const dateFrom = ref('')
const dateTo = ref('')
const timelineInterval = ref('month')
const selectedConfidence = ref(['high', 'medium', 'low'])
const confidenceLevels = ['high', 'medium', 'low']

const categoryLoading = ref(true)
const timelineLoading = ref(true)
const forecastLoading = ref(true)

const categoryData = ref<PlannedCategory[]>([])
const timelineData = ref<TimelinePoint[]>([])
const forecastData = ref<ForecastPoint[]>([])

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
    const res = await api.get<{ categories: PlannedCategory[] }>(
      `/analytics/outlook/planned-by-category?date_from=${dateFrom.value}&date_to=${dateTo.value}`
    )
    categoryData.value = res.data.categories
  } catch (e) {
    console.error('Failed to fetch planned by category:', e)
  } finally {
    categoryLoading.value = false
  }
}

async function fetchPlannedByTimeline() {
  timelineLoading.value = true
  try {
    const res = await api.get<{ timeline: TimelinePoint[] }>(
      `/analytics/outlook/planned-by-timeline?date_from=${dateFrom.value}&date_to=${dateTo.value}&interval=${timelineInterval.value}`
    )
    timelineData.value = res.data.timeline
  } catch (e) {
    console.error('Failed to fetch planned by timeline:', e)
  } finally {
    timelineLoading.value = false
  }
}

async function fetchBalanceForecast() {
  forecastLoading.value = true
  try {
    const res = await api.get<{ forecast: ForecastPoint[] }>(
      `/analytics/outlook/balance-forecast?date_from=${dateFrom.value}&date_to=${dateTo.value}`
    )
    forecastData.value = res.data.forecast
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
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); 
}
</style>