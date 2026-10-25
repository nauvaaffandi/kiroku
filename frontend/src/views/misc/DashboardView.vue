<template>
  <div class="dashboard-page">
    <main class="dashboard-main">
      <DashboardStatsWidget
        :monthlyIncome="monthlyIncome"
        :monthlyExpense="monthlyExpense"
        :monthlyNet="monthlyNet"
        :totalBalance="totalBalance"
      />

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

      <RecentTransactionsWidget
        :recentTransactions="recentTransactions"
        :categories="categoriesStore.categories"
        :accounts="accountsStore.accounts"
      />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import { useCategoriesStore } from '@/stores/categories'
import { useI18n } from '@/locales'
import DashboardStatsWidget from '@/widgets/dashboard/DashboardStatsWidget.vue'
import RecentTransactionsWidget from '@/widgets/dashboard/RecentTransactionsWidget.vue'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([BarChart, LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const { t } = useI18n()

const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const categoriesStore = useCategoriesStore()

const spendingChart = ref<HTMLElement | null>(null)
const trendChart = ref<HTMLElement | null>(null)
let spendingChartInstance: echarts.ECharts | null = null
let trendChartInstance: echarts.ECharts | null = null

const recentTransactions = computed(() => (transactionsStore.transactions || []).slice(0, 5))

const monthlyIncome = computed(() => {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), 1)
  return (transactionsStore.transactions || [])
    .filter(t => t.type === 'income' && new Date(t.date_time) >= start)
    .reduce((sum, t) => sum + t.amount, 0)
})

const monthlyExpense = computed(() => {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), 1)
  return (transactionsStore.transactions || [])
    .filter(t => t.type === 'expense' && new Date(t.date_time) >= start)
    .reduce((sum, t) => sum + t.amount, 0)
})

const monthlyNet = computed(() => monthlyIncome.value - monthlyExpense.value)

const totalBalance = computed(() => {
  return (accountsStore.accounts || []).reduce((sum, a) => sum + a.balance, 0)
})

const formatCurrency = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(Math.abs(num))
}

const initCharts = async () => {
  if (spendingChart.value && categoriesStore.categories.length > 0) {
    spendingChartInstance = echarts.init(spendingChart.value)
    const expenseCats = categoriesStore.categories
      .filter(c => c.type === 'expense')
      .slice(0, 8)
    
    const data = expenseCats.map(c => {
      const total = (transactionsStore.transactions || [])
        .filter(t => t.category_id === c.id && t.type === 'expense')
        .reduce((sum, t) => sum + t.amount, 0)
      return { name: c.name, value: total, color: c.color }
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

onMounted(async () => {
  await Promise.all([
    transactionsStore.fetchTransactions(),
    accountsStore.fetchAccounts(),
    categoriesStore.fetchCategories()
  ])
  await initCharts()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
@reference "tailwindcss";
.dashboard-page { @apply min-h-screen flex flex-col; }
.dashboard-main { @apply flex-1 p-[1.5rem] mx-auto w-full max-w-[1200px]; }
.charts-section { @apply grid grid-cols-[repeat(auto-fit,minmax(400px,1fr))] gap-6 mb-6; }
@media (max-width: 768px) { .charts-section { @apply grid-cols-1; } }
.chart-card { @apply rounded-xl shadow-sm p-[1.5rem]; background-color: var(--bg-card); }
.chart-card h2 { @apply text-base font-semibold mb-4; color: var(--text-primary); }
.chart-container { @apply w-full h-[300px]; }
</style>