<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>残高予測</h2>
      <div class="confidence-toggles">
        <label v-for="level in confidenceLevels" :key="level" class="toggle-label">
          <input 
            type="checkbox" 
            :value="level" 
            :checked="selectedConfidence.includes(level)"
            @change="toggleConfidence(level)"
          />
          <span :class="['confidence-badge', level]">{{ level }}</span>
        </label>
      </div>
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

interface ForecastPoint {
  date: string
  projected_balance: number
  confidence_level: string
}

const props = defineProps<{
  loading: boolean
  data: ForecastPoint[]
  selectedConfidence: string[]
  confidenceLevels: string[]
}>()

const emit = defineEmits<{
  'update:selectedConfidence': [val: string[]]
}>()

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

function renderChart() {
  if (!chartRef.value || props.data.length === 0) return
  if (chart) chart.dispose()

  chart = echarts.init(chartRef.value)
  const dates = props.data.map(d => d.date)
  const balances = props.data.map(d => d.projected_balance)

  chart.setOption({
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const p = params[0]
        return `${p.name}<br/>予測残高: ${formatCurrencyShort(p.value)}`
      },
    },
    grid: { left: 40, right: 20, top: 10, bottom: 40 },
    xAxis: { type: 'category', data: dates, boundaryGap: false },
    yAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    series: [
      {
        name: '予測残高',
        type: 'line',
        data: balances,
        smooth: true,
        lineStyle: { color: '#4a90d9', width: 2 },
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(74,144,217,0.2)' }, { offset: 1, color: 'rgba(74,144,217,0)' }] } },
      },
    ],
  })
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

function toggleConfidence(level: string) {
  const current = [...props.selectedConfidence]
  const index = current.indexOf(level)
  if (index >= 0) {
    current.splice(index, 1)
  } else {
    current.push(level)
  }
  emit('update:selectedConfidence', current)
}

watch(() => props.data, () => {
  nextTick(() => renderChart())
}, { deep: true })
</script>

<style scoped>
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; border-bottom: 1px solid var(--border-primary); background-color: var(--bg-secondary); }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.confidence-toggles { display: flex; gap: 0.5rem; }
.toggle-label { display: flex; align-items: center; gap: 0.25rem; cursor: pointer; }
.toggle-label input[type="checkbox"] { accent-color: var(--accent); }
.confidence-badge { padding: 0.125rem 0.5rem; font-size: 0.75rem; border-radius: 9999px; font-weight: 500; }
.confidence-badge.high { background-color: var(--success-light); color: var(--success); }
.confidence-badge.medium { background-color: var(--warning-light); color: var(--warning); }
.confidence-badge.low { background-color: var(--error-light); color: var(--error); }
.widget-content { padding: 1rem; }
.chart-container { width: 100%; height: 300px; }
</style>