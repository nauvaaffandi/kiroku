<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>{{ t('outlook.plannedByTimeline') }}</h2>
      <select :value="interval" @change="$emit('update:interval', ($event.target as HTMLSelectElement).value)" class="interval-select">
        <option value="day">{{ t('outlook.intervalDay') }}</option>
        <option value="week">{{ t('outlook.intervalWeek') }}</option>
        <option value="month">{{ t('outlook.intervalMonth') }}</option>
      </select>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="chart-container" ref="chartRef" />
    </div>
    <div v-else class="widget-content">
      <div class="skeleton skeleton-chart" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import { useI18n } from '@/locales'

interface TimelinePoint {
  period: string
  planned_income: number
  planned_expense: number
  net: number
  transaction_count: number
}

const { t } = useI18n()

const props = defineProps<{
  loading: boolean
  data: TimelinePoint[]
  interval: string
}>()

defineEmits<{
  'update:interval': [val: string]
}>()

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

function renderChart() {
  if (!chartRef.value || props.data.length === 0) return
  if (chart) chart.dispose()

  chart = echarts.init(chartRef.value)
  const periods = props.data.map(d => d.period)
  const income = props.data.map(d => d.planned_income)
  const expense = props.data.map(d => d.planned_expense)
  const net = props.data.map(d => d.net)

  chart.setOption({
    tooltip: {
      trigger: 'axis',
formatter: (params: any) => {
          const p = params[0]
          return `${p.name}<br/>${t('outlook.income')}: ${formatCurrencyShort(p.value)}<br/>${t('outlook.expense')}: ${formatCurrencyShort(expense[p.dataIndex])}<br/>${t('outlook.net')}: ${formatCurrencyShort(net[p.dataIndex])}`
        },
    },
    legend: { data: [t('outlook.income'), t('outlook.expense'), t('outlook.net')], bottom: 0 },
    grid: { left: 40, right: 20, top: 10, bottom: 40 },
    xAxis: { type: 'category', data: periods },
    yAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    series: [
      {
        name: '収入',
        type: 'line',
        data: income,
        smooth: true,
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(40,167,69,0.3)' }, { offset: 1, color: 'rgba(40,167,69,0)' }] } },
        lineStyle: { color: '#28a745', width: 2 },
      },
      {
        name: '支出',
        type: 'line',
        data: expense,
        smooth: true,
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(220,53,69,0.3)' }, { offset: 1, color: 'rgba(220,53,69,0)' }] } },
        lineStyle: { color: '#dc3545', width: 2 },
      },
      {
        name: '純額',
        type: 'line',
        data: net,
        smooth: true,
        lineStyle: { color: '#4a90d9', width: 3 },
      },
    ],
  })
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

watch(() => [props.data, props.interval], () => {
  nextTick(() => renderChart())
}, { deep: true })
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
