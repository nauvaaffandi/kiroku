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
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; border-bottom: 1px solid var(--border-primary); background-color: var(--bg-secondary); }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.widget-content { padding: 1rem; }
.chart-container { width: 100%; height: 300px; }
</style>