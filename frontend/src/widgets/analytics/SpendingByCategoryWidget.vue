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
.widget-controls {
  @apply flex gap-2;
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
.category-legend {
  @apply flex flex-col gap-2 mt-4 pt-4 border-t;
  border-color: var(--border-primary);
}
.legend-item {
  @apply flex items-center gap-2 text-sm;
}
.legend-color {
  @apply w-[12px] h-[12px] rounded-[2px] shrink-0;
}
.legend-name {
  @apply flex-1 truncate;
  color: var(--text-primary);
}
.legend-value {
  @apply text-xs whitespace-nowrap;
  color: var(--text-secondary);
}
</style>
