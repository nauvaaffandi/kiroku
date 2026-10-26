<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>キャッシュフロー・ウォーターフォール</h2>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="chart-container" ref="flowChartRef" />
    </div>
    <div v-else class="widget-content">
      <div class="skeleton skeleton-chart" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'

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

const props = defineProps<{
  loading: boolean
  data: FlowResponse | null
}>()

const flowChartRef = ref<HTMLElement>()
let flowChart: echarts.ECharts | null = null

function renderFlowChart() {
  if (!flowChartRef.value || !props.data) return
  if (flowChart) flowChart.dispose()

  flowChart = echarts.init(flowChartRef.value)
  const data = props.data
  const categories = data.by_category.filter(c => c.net !== 0).slice(0, 8)

  const seriesData = [
    { name: '期首残高', value: data.opening_balance },
    ...categories.flatMap(c => [
      { name: `収入: ${c.name}`, value: c.income },
      { name: `支出: ${c.name}`, value: -c.expense },
    ]),
    { name: '期末残高', value: data.closing_balance },
  ]

  flowChart.setOption({
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const p = params[0]
        return `${p.name}<br/>${formatCurrencyShort(p.value)}`
      },
    },
    grid: { left: 10, right: 10, top: 10, bottom: 40 },
    xAxis: { type: 'category', data: seriesData.map(d => d.name), axisLabel: { interval: 0, rotate: 30 } },
    yAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    series: [
      {
        name: '金額',
        type: 'bar',
        data: seriesData.map((d, i) => ({
          value: d.value,
          itemStyle: {
            color: i === 0 ? '#4a90d9' : i === seriesData.length - 1 ? '#28a745' : d.value > 0 ? 'rgba(40,167,69,0.7)' : 'rgba(220,53,69,0.7)',
          },
        })),
        label: { show: true, position: 'top', formatter: (p: any) => formatCurrencyShort(p.value) },
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
  nextTick(() => renderFlowChart())
}, { deep: true })

onMounted(() => {
  nextTick(() => renderFlowChart())
  window.addEventListener('resize', () => flowChart?.resize())
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
.widget-content {
  @apply p-4;
}
.chart-container {
  @apply w-full;
  height: 300px;
}
</style>
