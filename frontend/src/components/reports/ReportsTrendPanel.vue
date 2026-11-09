<template>
  <section class="trend-section">
    <h2 class="trend-title">{{ t('reports.monthlyTrend') }}</h2>
    <div class="trend-chart" ref="el"></div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { useI18n } from '@/locales'

echarts.use([LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{ series: any[] }>()
const el = ref<HTMLElement | null>(null)
const { t } = useI18n()
let instance: echarts.ECharts | null = null

const draw = () => {
  if (!el.value) return
  if (instance) instance.dispose()
  instance = echarts.init(el.value)
  const data = props.series
  instance.setOption({
    title: { text: 'Monthly Income vs Expense', left: 'center', textStyle: { fontSize: 13, color: 'inherit' } },
    tooltip: { trigger: 'axis' },
    legend: { data: ['Income', 'Expense', 'Net'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true },
    xAxis: { type: 'category', data: data.map((d: any) => d.month) },
    yAxis: { type: 'value', axisLabel: { formatter: '¥{value}' } },
    series: [
      { name: 'Income', type: 'line', data: data.map((d: any) => d.income), smooth: true, itemStyle: { color: '#28a745' } },
      { name: 'Expense', type: 'line', data: data.map((d: any) => d.expense), smooth: true, itemStyle: { color: '#dc3545' } },
      { name: 'Net', type: 'line', data: data.map((d: any) => d.net), smooth: true, itemStyle: { color: '#4a90d9' } }
    ]
  })
}

watch(() => props.series, draw, { deep: true })
onMounted(draw)
window.addEventListener('resize', () => instance?.resize())
</script>

<style scoped>
@reference "tailwindcss";

.trend-section { @apply bg-zinc-50 dark:bg-zinc-800/50 rounded-xl p-4 border border-zinc-200 dark:border-zinc-700; }
.trend-title { @apply text-base font-semibold text-zinc-900 dark:text-zinc-100 mb-4; }
.trend-chart { @apply w-full; height: 300px; }
@media (min-width: 640px) { .trend-chart { height: 350px; } }
</style>
