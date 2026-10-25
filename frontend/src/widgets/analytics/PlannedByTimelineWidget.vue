<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>タイムライン別予定</h2>
      <select :value="interval" @change="$emit('update:interval', ($event.target as HTMLSelectElement).value)" class="interval-select">
        <option value="day">日次</option>
        <option value="week">週次</option>
        <option value="month">月次</option>
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

interface TimelinePoint {
  period: string
  planned_income: number
  planned_expense: number
  net: number
  transaction_count: number
}

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
        return `${p.name}<br/>収入: ${formatCurrencyShort(p.value)}<br/>支出: ${formatCurrencyShort(expense[p.dataIndex])}<br/>純額: ${formatCurrencyShort(net[p.dataIndex])}`
      },
    },
    legend: { data: ['収入', '支出', '純額'], bottom: 0 },
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
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; border-bottom: 1px solid var(--border-primary); background-color: var(--bg-secondary); }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.interval-select { padding: 0.25rem 0.5rem; font-size: 0.875rem; border: 1px solid var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); border-radius: var(--radius-sm); }
.widget-content { padding: 1rem; }
.chart-container { width: 100%; height: 300px; }
</style>