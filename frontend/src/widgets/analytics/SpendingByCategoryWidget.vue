<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>カテゴリ別支出</h2>
      <div class="widget-controls">
        <select :value="groupBy" @change="$emit('update:groupBy', ($event.target as HTMLSelectElement).value)" class="interval-select">
          <option value="category">カテゴリ別</option>
          <option value="parent_category">親カテゴリ別</option>
        </select>
      </div>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="chart-container" ref="chartRef" />
      <div class="category-legend" v-if="data.length > 0">
        <div v-for="cat in data.slice(0, 10)" :key="cat.id" class="legend-item">
          <span class="legend-color" :style="{ backgroundColor: cat.color || '#4a90d9' }" />
          <span class="legend-name">{{ cat.name }}</span>
          <span class="legend-value">{{ formatCurrencyShort(cat.amount) }} ({{ cat.percentage.toFixed(1) }}%)</span>
        </div>
      </div>
    </div>
    <div v-else class="widget-content">
      <div class="skeleton skeleton-chart" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import * as echarts from 'echarts'

interface CategorySpending {
  id: string
  name: string
  parent_id: string
  color: string
  kind: string
  amount: number
  count: number
  percentage: number
  children: any[]
}

const props = defineProps<{
  loading: boolean
  data: CategorySpending[]
  groupBy: string
}>()

defineEmits<{
  'update:groupBy': [val: string]
}>()

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

function renderChart() {
  if (!chartRef.value || props.data.length === 0) return
  if (chart) chart.dispose()

  chart = echarts.init(chartRef.value)
  const data = props.data.slice(0, 10)
  const colors = data.map(d => d.color || '#4a90d9')

  chart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: (params: any) => `${params.name}<br/>${formatCurrencyShort(params.value)} (${params.percent}%)`,
    },
    series: [
      {
        name: '支出',
        type: 'treemap',
        data: data.map((d, i) => ({
          name: d.name,
          value: d.amount,
          itemStyle: { color: colors[i] },
          label: { show: true, formatter: '{b}\n{c}', position: 'inside' },
        })),
        breadcrumb: { show: false },
        nodeClick: false,
      },
    ],
  })
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

watch(() => [props.data, props.groupBy], () => {
  nextTick(() => renderChart())
}, { deep: true })
</script>

<style scoped>
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; border-bottom: 1px solid var(--border-primary); background-color: var(--bg-secondary); }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.widget-controls { display: flex; gap: 0.5rem; }
.interval-select { padding: 0.25rem 0.5rem; font-size: 0.875rem; border: 1px solid var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); border-radius: var(--radius-sm); }
.widget-content { padding: 1rem; }
.chart-container { width: 100%; height: 300px; }
.category-legend { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 1rem; padding-top: 1rem; border-top: 1px solid var(--border-primary); }
.legend-item { display: flex; align-items: center; gap: 0.5rem; font-size: 0.875rem; }
.legend-color { width: 12px; height: 12px; border-radius: 2px; flex-shrink: 0; }
.legend-name { flex: 1; color: var(--text-primary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.legend-value { font-size: 0.75rem; color: var(--text-secondary); white-space: nowrap; }
</style>