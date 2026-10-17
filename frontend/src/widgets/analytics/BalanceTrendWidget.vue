<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>残高推移</h2>
      <select :value="interval" @change="$emit('update:interval', ($event.target as HTMLSelectElement).value)" class="interval-select">
        <option value="day">日次</option>
        <option value="week">週次</option>
        <option value="month">月次</option>
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
import { ref, onMounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import type { BalanceTrendPoint } from '@/types/report'
import { useSettingsStore } from '@/stores/settings'
import { formatAmountForDisplay } from '@/utils/amountDisplay'

const props = defineProps<{
  interval: string
  loading: boolean
  data: BalanceTrendPoint[]
}>()

defineEmits<{
  'update:interval': [val: string]
}>()

const trendChartRef = ref<HTMLElement>()
let trendChart: echarts.ECharts | null = null
const settingsStore = useSettingsStore()

function renderTrendChart() {
  if (!trendChartRef.value || props.data.length === 0) return
  if (trendChart) trendChart.dispose()

  trendChart = echarts.init(trendChartRef.value)
  const periods = props.data.map(d => d.period)
  const closing = props.data.map(d => d.closing_balance)
  const net = props.data.map(d => d.net_change)

  trendChart.setOption({
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const p = params[0]
        return `${p.name}<br/>期末: ${formatCurrency(p.value)}<br/>増減: ${p.value >= 0 ? '+' : ''}${formatCurrency(net[p.dataIndex])}`
      },
    },
    grid: { left: 40, right: 20, top: 10, bottom: 40 },
    xAxis: { type: 'category', data: periods, boundaryGap: false },
    yAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    series: [
      {
        name: '残高',
        type: 'line',
        data: closing,
        smooth: true,
        lineStyle: { color: '#4a90d9', width: 2 },
        areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(74,144,217,0.3)' }, { offset: 1, color: 'rgba(74,144,217,0)' }] } },
      },
    ],
  })
}

function formatCurrency(val: number, currency: string = 'JPY'): string {
  if (settingsStore.isHideAmounts) {
    return formatAmountForDisplay(val, 'transfer')
  }
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(val)
}

function formatCurrencyShort(val: number): string {
  if (settingsStore.isHideAmounts) {
    return formatAmountForDisplay(val, 'transfer')
  }
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
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; border-bottom: 1px solid var(--border-primary); background-color: var(--bg-secondary); }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.interval-select { padding: 0.25rem 0.5rem; font-size: 0.875rem; border: 1px solid var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); border-radius: var(--radius-sm); }
.widget-content { padding: 1rem; }
.chart-container { width: 100%; height: 300px; }
</style>