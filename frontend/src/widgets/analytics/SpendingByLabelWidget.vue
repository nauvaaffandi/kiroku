<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>ラベル別支出</h2>
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

interface LabelSpending {
  label: string
  amount: number
  count: number
  percentage: number
}

const props = defineProps<{
  loading: boolean
  data: LabelSpending[]
}>()

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

function renderChart() {
  if (!chartRef.value || props.data.length === 0) return
  if (chart) chart.dispose()

  chart = echarts.init(chartRef.value)
  const data = props.data

  chart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: (params: any) => `${params.name}<br/>${formatCurrencyShort(params.value)} (${params.percent}%)`,
    },
    legend: { orient: 'vertical', right: 0, top: 'center' },
    series: [
      {
        name: 'ラベル別支出',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        data: data.map((d, i) => ({
          name: d.label,
          value: d.amount,
          itemStyle: { color: getColor(i) },
        })),
        label: { show: false },
        emphasis: {
          label: { show: true, fontSize: 14, fontWeight: 'bold' },
        },
      },
    ],
  })
}

function getColor(index: number): string {
  const colors = ['#4a90d9', '#28a745', '#ffc107', '#dc3545', '#17a2b8', '#6f42c1', '#e83e8c', '#fd7e14', '#20c997', '#6c757d']
  return colors[index % colors.length]
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
