<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>カテゴリ別予定</h2>
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

const props = defineProps<{
  loading: boolean
  data: PlannedCategory[]
}>()

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

function renderChart() {
  if (!chartRef.value || props.data.length === 0) return
  if (chart) chart.dispose()

  chart = echarts.init(chartRef.value)
  const names = props.data.map(d => d.name)
  const income = props.data.map(d => d.planned_income)
  const expense = props.data.map(d => d.planned_expense)
  const colors = props.data.map(d => d.color || '#4a90d9')

  chart.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      formatter: (params: any) => {
        const p = params[0]
        return `${p.name}<br/>収入: ${formatCurrencyShort(p.value)}<br/>支出: ${formatCurrencyShort(expense[p.dataIndex])}`
      },
    },
    legend: { data: ['収入', '支出'], bottom: 0 },
    grid: { left: 10, right: 10, top: 10, bottom: 40 },
    xAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    yAxis: { type: 'category', data: names, axisLabel: { interval: 0 } },
    series: [
      {
        name: '収入',
        type: 'bar',
        stack: 'total',
        data: income.map((v, i) => ({ value: v, itemStyle: { color: colors[i] } })),
        label: { show: false },
      },
      {
        name: '支出',
        type: 'bar',
        stack: 'total',
        data: expense.map(v => ({ value: v, itemStyle: { color: '#dc3545' } })),
        label: { show: false },
      },
    ],
  })
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

watch(() => props.data, () => {
  nextTick(() => renderChart())
}, { deep: true })
</script>

<style scoped>
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; border-bottom: 1px solid var(--border-primary); background-color: var(--bg-secondary); }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.widget-content { padding: 1rem; }
.chart-container { width: 100%; height: 300px; }
</style>