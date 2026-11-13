<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>{{ t('analytics.cashflow.trend.title') }}</h2>
      <select :value="interval" @change="$emit('update:interval', ($event.target as HTMLSelectElement).value)" class="interval-select">
        <option value="day">{{ t('analytics.cashflow.trend.intervalDay') }}</option>
        <option value="week">{{ t('analytics.cashflow.trend.intervalWeek') }}</option>
        <option value="month">{{ t('analytics.cashflow.trend.intervalMonth') }}</option>
      </select>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="chart-container" ref="trendChartRef" />
    </div>
    <div v-else class="widget-content">
      <div class="skeleton skeleton-chart" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, withDefaults } from 'vue'
import * as echarts from 'echarts'
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import { formatAmountForDisplay } from '@/utils/amountDisplay'

interface TrendPoint {
  period: string
  income: number
  expense: number
  net: number
  cumulative_net: number
}

const props = withDefaults(defineProps<{
  interval: string
  loading: boolean
  data: TrendPoint[]
}>(), {
  data: () => [] as TrendPoint[]
})

const { t } = useI18n()
const settingsStore = useSettingsStore()

defineEmits<{
  'update:interval': [val: string]
}>()

const trendChartRef = ref<HTMLElement>()
let trendChart: echarts.ECharts | null = null

function renderTrendChart() {
  if (!trendChartRef.value || !props.data || props.data.length === 0) return
  if (trendChart) trendChart.dispose()

  trendChart = echarts.init(trendChartRef.value)
  const periods = props.data.map(d => d.period)
  const income = props.data.map(d => d.income)
  const expense = props.data.map(d => d.expense)
  const net = props.data.map(d => d.net)
  const cumulative = props.data.map(d => d.cumulative_net)

  trendChart.setOption({
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        let res = `${params[0].name}<br/>`
        params.forEach((p: any) => {
          res += `${p.marker} ${p.seriesName}: ${formatCurrencyShort(p.value)}<br/>`
        })
        return res
      },
    },
    legend: { data: [t('analytics.cashflow.trend.income'), t('analytics.cashflow.trend.expense'), t('analytics.cashflow.trend.net'), t('analytics.cashflow.trend.cumulative')], bottom: 0 },
    grid: { left: 40, right: 20, top: 10, bottom: 40 },
    xAxis: { type: 'category', data: periods },
    yAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    series: [
      {
        name: t('analytics.cashflow.trend.income'),
        type: 'line',
        data: income,
        smooth: true,
        lineStyle: { color: '#28a745', width: 2 },
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(40,167,69,0.3)' }, { offset: 1, color: 'rgba(40,167,69,0)' }] } },
      },
      {
        name: t('analytics.cashflow.trend.expense'),
        type: 'line',
        data: expense,
        smooth: true,
        lineStyle: { color: '#dc3545', width: 2 },
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(220,53,69,0.3)' }, { offset: 1, color: 'rgba(220,53,69,0)' }] } },
      },
      {
        name: t('analytics.cashflow.trend.net'),
        type: 'line',
        data: net,
        smooth: true,
        lineStyle: { color: '#4a90d9', width: 3 },
      },
      {
        name: t('analytics.cashflow.trend.cumulative'),
        type: 'line',
        data: cumulative,
        smooth: true,
        lineStyle: { color: '#6c757d', width: 2, type: 'dashed' },
      },
    ],
  })
}

function formatCurrencyShort(val: number): string {
  if (settingsStore.isHideAmounts) return formatAmountForDisplay(val, 'transfer')
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

watch(() => props.data, () => {
  nextTick(() => renderTrendChart())
}, { deep: true })

onMounted(() => {
  nextTick(() => renderTrendChart())
  window.addEventListener('resize', () => trendChart?.resize())
})
</script>

<style scoped>
@reference "tailwindcss";

.widget-card {
  @apply border rounded-xl overflow-hidden;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  box-shadow: var(--shadow);
}
.widget-header {
  @apply flex items-center justify-between p-4 border-b;
  background-color: var(--bg-secondary);
  border-color: var(--border-primary);
}
.widget-header h2 {
  @apply text-lg font-semibold;
  color: var(--text-primary);
}
.interval-select {
  @apply px-2 py-1 text-sm border rounded-sm;
  border-color: var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}
.widget-content {
  @apply p-4;
}
.chart-container {
  @apply w-full;
  height: 300px;
}
</style>
