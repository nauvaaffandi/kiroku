<template>
  <section class="trend-section">
    <h2 class="trend-title">{{ t('reports.monthlyTrend') }}</h2>
    <div class="trend-chart" ref="el"></div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { getCurrencySymbol } from '@/utils/amountDisplay'
import { useSettingsStore } from '@/stores/settings'
import { useI18n } from '@/locales'

echarts.use([LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{ series: any[] }>()
const el = ref<HTMLElement | null>(null)
const { t } = useI18n()
const settingsStore = useSettingsStore()
const currencySymbol = computed(() => getCurrencySymbol(settingsStore.baseCurrency))
let instance: echarts.ECharts | null = null

const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

const cssVar = (name: string): string => {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || '#000000'
}

const draw = () => {
  if (!el.value) return
  if (instance) instance.dispose()
  instance = echarts.init(el.value)
  const data = props.series
  const symbol = currencySymbol.value
  const incomeColor = cssVar('--success')
  const expenseColor = cssVar('--error')
  const netColor = cssVar('--accent')
  instance.setOption({
    title: { text: t('reports.monthlyTrend'), left: 'center', textStyle: { fontSize: 13, color: 'inherit' } },
    tooltip: { trigger: 'axis', formatter: (params: any) => {
      const line = params.map((p: any) => `${p.seriesName}: ${symbol}${p.value}`).join('<br/>')
      return line
    } },
    legend: { data: [t('analytics.income'), t('analytics.expense'), t('analytics.net')], bottom: 0, textStyle: { fontSize: 11 } },
    animation: !reducedMotion,
    grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true },
    xAxis: { type: 'category', data: data.map((d: any) => d.month) },
    yAxis: { type: 'value', axisLabel: { formatter: `${symbol}{value}` } },
    series: [
      { name: t('analytics.income'), type: 'line', data: data.map((d: any) => d.income), smooth: true, itemStyle: { color: incomeColor } },
      { name: t('analytics.expense'), type: 'line', data: data.map((d: any) => d.expense), smooth: true, itemStyle: { color: expenseColor } },
      { name: t('analytics.net'), type: 'line', data: data.map((d: any) => d.net), smooth: true, itemStyle: { color: netColor } }
    ]
  })
}

watch(() => props.series, draw, { deep: true })
onMounted(draw)
window.addEventListener('resize', () => instance?.resize())
</script>

<style scoped>
@reference "tailwindcss";

.trend-section {
  @apply bg-[var(--bg-card)] border border-[var(--border-primary)] rounded-[14px] p-4;
}

.trend-title { @apply text-base font-semibold text-[var(--text-primary)] mb-4; }

.trend-chart { @apply w-full; min-height: 300px; }
</style>
