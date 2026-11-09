<template>
  <section class="charts-section">
    <div class="chart-card">
      <h2>{{ t('dashboard.spendingByCategory') }}</h2>
      <div class="chart-container" ref="spendingChart"></div>
    </div>
    <div class="chart-card">
      <h2>{{ t('dashboard.monthlyTrend') }}</h2>
      <div class="chart-container" ref="trendChart"></div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useI18n } from '@/locales'
import { useCategoriesStore } from '@/stores/categories'
import { useTransactionsStore } from '@/stores/transactions'
import { useSettingsStore } from '@/stores/settings'
import { useCategoryName } from '@/utils/categoryName'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([BarChart, LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const { t } = useI18n()
const { translate } = useCategoryName()
const categoriesStore = useCategoriesStore()
const transactionsStore = useTransactionsStore()
const settingsStore = useSettingsStore()

const spendingChart = ref<HTMLElement | null>(null)
const trendChart = ref<HTMLElement | null>(null)
let spendingChartInstance: echarts.ECharts | null = null
let trendChartInstance: echarts.ECharts | null = null

const numberFormat = computed(() => settingsStore.numberFormat)
const formatCurrency = (num: number) => new Intl.NumberFormat(numberFormat.value).format(Math.abs(num))

const initCharts = () => {
  if (spendingChart.value && categoriesStore.categories.length > 0) {
    spendingChartInstance = echarts.init(spendingChart.value)
    const expenseCats = categoriesStore.categories.slice(0, 8)
    const data = expenseCats.map(c => {
      const total = (transactionsStore.transactions || []).filter(t => t.category_id === c.id).reduce((sum, t) => sum + t.amount, 0)
      return { name: translate(c), value: total, color: c.color }
    }).sort((a, b) => b.value - a.value)

    spendingChartInstance.setOption({
      title: { text: t('dashboard.spendingByCategory'), left: 'center' },
      tooltip: { trigger: 'item', formatter: '{b}: ' + formatCurrency(0).replace('0', '{c}') },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'value', axisLabel: { formatter: (val: number) => formatCurrency(val) } },
      yAxis: { type: 'category', data: data.map(d => d.name), inverse: true },
      series: [{
        type: 'bar',
        data: data.map(d => ({ value: d.value, itemStyle: { color: d.color } })),
        label: { show: true, position: 'right', formatter: (params: any) => formatCurrency(params.value) }
      }]
    })
  }

  if (trendChart.value) {
    trendChartInstance = echarts.init(trendChart.value)
    trendChartInstance.setOption({
      title: { text: t('dashboard.monthlyTrend'), left: 'center' },
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: [] },
      yAxis: { type: 'value' },
      series: []
    })
  }
}

const handleResize = () => {
  spendingChartInstance?.resize()
  trendChartInstance?.resize()
}

onMounted(() => {
  initCharts()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  spendingChartInstance?.dispose()
  trendChartInstance?.dispose()
})
</script>

<style scoped>
@reference "tailwindcss";
.charts-section { @apply grid grid-cols-[repeat(auto-fit,minmax(400px,1fr))] gap-6 mb-6; }
@media (max-width: 768px) { .charts-section { @apply grid-cols-1; } }
.chart-card { @apply rounded-xl shadow-sm p-[1.5rem]; background-color: var(--bg-card); }
.chart-card h2 { @apply text-base font-semibold mb-4; color: var(--text-primary); }
.chart-container { @apply w-full h-[300px]; }
</style>