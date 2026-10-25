<template>
  <div class="min-h-screen flex flex-col pb-20">
    <header class="p-4 bg-zinc-100 dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800">
      <h1 class="text-xl font-bold text-zinc-900 dark:text-zinc-100">{{ t('reports.title') }}</h1>
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 mt-3">
        <select v-model="selectedCurrency" @change="fetchSpending" class="w-full sm:w-auto px-3 py-2 bg-white dark:bg-zinc-800 border border-zinc-300 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:outline-none focus:ring-2 focus:ring-sky-500">
          <option v-for="c of currencies" :key="c" :value="c">{{ c }}</option>
        </select>
        <button @click="exportCSV" class="flex items-center justify-center gap-2 px-4 py-2 bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 text-zinc-800 dark:text-zinc-200 text-sm font-medium rounded-lg transition-colors cursor-pointer border-none">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
          Export CSV
        </button>
        <button @click="exportExcel" class="flex items-center justify-center gap-2 px-4 py-2 bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 text-zinc-800 dark:text-zinc-200 text-sm font-medium rounded-lg transition-colors cursor-pointer border-none">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
          </svg>
          {{ t('reports.exportExcel') }}
        </button>
      </div>
    </header>

    <div class="flex flex-col lg:grid lg:grid-cols-[280px_1fr] gap-4 p-4 flex-1">
      <aside class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl p-4 shadow-sm h-fit">
        <section class="mb-6">
          <h3 class="text-xs font-bold text-zinc-500 dark:text-zinc-400 uppercase tracking-wider mb-3">{{ t('reports.spendingReport') }}</h3>
          <div class="space-y-3">
            <div>
              <label class="block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1">{{ t('reports.period') }}</label>
              <select v-model="spendingParams.period" @change="fetchSpending" class="w-full px-3 py-2 bg-white dark:bg-zinc-800 border border-zinc-300 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:outline-none focus:ring-2 focus:ring-sky-500">
                <option value="this_week">This Week</option>
                <option value="this_month">This Month</option>
                <option value="this_year">This Year</option>
                <option value="last_month">Last Month</option>
                <option value="last_year">Last Year</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1">{{ t('reports.groupBy') }}</label>
              <select v-model="spendingParams.group_by" @change="fetchSpending" class="w-full px-3 py-2 bg-white dark:bg-zinc-800 border border-zinc-300 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:outline-none focus:ring-2 focus:ring-sky-500">
                <option value="category">{{ t('reports.category') }}</option>
                <option value="parent_category">{{ t('reports.parentCategory') }}</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1">{{ t('reports.accountOptional') }}</label>
              <select v-model="spendingParams.account_id" @change="fetchSpending" class="w-full px-3 py-2 bg-white dark:bg-zinc-800 border border-zinc-300 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:outline-none focus:ring-2 focus:ring-sky-500">
                <option value="">{{ t('reports.allAccounts') }}</option>
                <option v-for="acc in accountsStore.accounts" :key="acc.id" :value="acc.id">
                  {{ acc.name }}
                </option>
              </select>
            </div>
          </div>
        </section>

        <section>
          <h3 class="text-xs font-bold text-zinc-500 dark:text-zinc-400 uppercase tracking-wider mb-3">{{ t('reports.trendReport') }}</h3>
          <div>
            <label class="block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1">{{ t('reports.months') }}</label>
            <select v-model="trendParams.months" @change="fetchTrend" class="w-full px-3 py-2 bg-white dark:bg-zinc-800 border border-zinc-300 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:outline-none focus:ring-2 focus:ring-sky-500">
              <option value="6">{{ 6 }} {{ t('reports.months') }}</option>
              <option value="12" selected>{{ 12 }} {{ t('reports.months') }}</option>
              <option value="24">{{ 24 }} {{ t('reports.months') }}</option>
            </select>
          </div>
        </section>
      </aside>

      <main class="flex flex-col gap-4">
        <div class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden shadow-sm">
          <div class="flex bg-zinc-100 dark:bg-zinc-800 border-b border-zinc-200 dark:border-zinc-700 overflow-x-auto">
            <button :class="['flex-1 px-4 py-3 text-xs sm:text-sm font-medium whitespace-nowrap transition-colors cursor-pointer border-b-2', activeTab === 'spending' ? 'text-sky-500 border-sky-500 bg-white dark:bg-zinc-900' : 'text-zinc-600 dark:text-zinc-400 border-transparent hover:text-zinc-900 dark:hover:text-zinc-100']" @click="activeTab = 'spending'">{{ t('reports.spendingReport') }}</button>
            <button :class="['flex-1 px-4 py-3 text-xs sm:text-sm font-medium whitespace-nowrap transition-colors cursor-pointer border-b-2', activeTab === 'trend' ? 'text-sky-500 border-sky-500 bg-white dark:bg-zinc-900' : 'text-zinc-600 dark:text-zinc-400 border-transparent hover:text-zinc-900 dark:hover:text-zinc-100']" @click="activeTab = 'trend'">{{ t('reports.trendReport') }}</button>
            <button :class="['flex-1 px-4 py-3 text-xs sm:text-sm font-medium whitespace-nowrap transition-colors cursor-pointer border-b-2', activeTab === 'cashflow' ? 'text-sky-500 border-sky-500 bg-white dark:bg-zinc-900' : 'text-zinc-600 dark:text-zinc-400 border-transparent hover:text-zinc-900 dark:hover:text-zinc-100']" @click="activeTab = 'cashflow'">{{ t('reports.cashflowTable') }}</button>
          </div>

          <div v-if="activeTab === 'spending'" class="p-4">
            <section class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl p-4 mb-4 border border-zinc-200 dark:border-zinc-700">
              <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
                <h2 class="text-base font-semibold text-zinc-900 dark:text-zinc-100">{{ t('reports.spendingByCategory') }}</h2>
                <span class="text-xs text-zinc-600 dark:text-zinc-400 bg-white dark:bg-zinc-800 px-2.5 py-1 rounded-full border border-zinc-200 dark:border-zinc-700">{{ spendingParams.period }}</span>
              </div>
              <div class="w-full h-[300px] sm:h-[350px] mb-4" ref="spendingChart"></div>
              <div class="overflow-x-auto" v-if="spendingReport.groups">
                <table class="w-full text-left text-sm border-collapse">
                  <thead>
                    <tr class="border-b border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 bg-white dark:bg-zinc-800">
                      <th class="p-3 font-semibold">Category</th>
                      <th class="p-3 font-semibold text-right">Amount</th>
                      <th class="p-3 font-semibold text-right">Count</th>
                      <th class="p-3 font-semibold text-right">%</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-zinc-200 dark:divide-zinc-700">
                    <tr v-for="group in spendingReport.groups" :key="group.id" class="hover:bg-zinc-100 dark:hover:bg-zinc-800/50">
                      <td class="p-3">
                        <div class="flex items-center gap-2">
                          <span class="w-2.5 h-2.5 rounded-full flex-shrink-0" :style="{ backgroundColor: group.color }"></span>
                          <span class="font-medium text-zinc-900 dark:text-zinc-100">{{ group.name }}</span>
                        </div>
                        <div class="ml-4 pl-2 border-l-2 border-zinc-200 dark:border-zinc-700 mt-1 space-y-1" v-if="group.children && group.children.length > 0">
                          <div v-for="child in group.children" :key="child.id" class="flex items-center gap-2 text-xs text-zinc-600 dark:text-zinc-400 py-0.5">
                            <span class="w-2 h-2 rounded-full flex-shrink-0" :style="{ backgroundColor: child.color }"></span>
                            <span>{{ child.name }}: {{ formatDisplayNumber(child.amount, 'expense') }}</span>
                          </div>
                        </div>
                      </td>
                      <td class="p-3 text-right font-medium text-zinc-900 dark:text-zinc-100">{{ formatDisplayNumber(group.amount, 'expense') }}</td>
                      <td class="p-3 text-right text-zinc-600 dark:text-zinc-400">{{ group.count }}</td>
                      <td class="p-3 text-right text-zinc-600 dark:text-zinc-400">{{ totalAmount > 0 ? ((group.amount / totalAmount) * 100).toFixed(1) : 0 }}%</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </section>
          </div>

          <div v-if="activeTab === 'trend'" class="p-4">
            <section class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl p-4 border border-zinc-200 dark:border-zinc-700">
              <h2 class="text-base font-semibold text-zinc-900 dark:text-zinc-100 mb-4">{{ t('reports.monthlyTrend') }}</h2>
              <div class="w-full h-[300px] sm:h-[350px]" ref="trendChart"></div>
            </section>
          </div>

          <div v-if="activeTab === 'cashflow'" class="p-4">
            <CashflowTable :date-from="spendingParams.date_from" :date-to="spendingParams.date_to" />
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from '@/locales'
import { useReportsStore } from '@/stores/reports'
import { useAccountsStore } from '@/stores/accounts'
import { useUIStore } from '@/stores/ui'
import { useSettingsStore } from '@/stores/settings'
import { useCurrencies } from '@/composables/useCurrencies'
import { formatAmountForDisplay, formatNumber as formatBaseNumber } from '@/utils/amountDisplay'
import { api } from '@/api/client'
import { exportToExcel } from '@/utils/exportExcel'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import CashflowTable from '@/components/charts/CashflowTable.vue'

echarts.use([BarChart, LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const reportsStore = useReportsStore()
const accountsStore = useAccountsStore()
const uiStore = useUIStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()

const spendingChart = ref<HTMLElement | null>(null)
const trendChart = ref<HTMLElement | null>(null)
let spendingChartInstance: echarts.ECharts | null = null
let trendChartInstance: echarts.ECharts | null = null

const activeTab = ref('spending')

const spendingParams = ref({ 
  period: 'this_month', 
  group_by: 'category', 
  account_id: '',
  date_from: '',
  date_to: '',
})
const trendParams = ref({ months: '12', account_id: '', category_id: '' })

const { currencies } = useCurrencies()
const selectedCurrency = ref(settingsStore.baseCurrency)

const spendingReport = ref({ period: '', groups: [] as any[] })
const trendReport = ref({ series: [] as any[] })

const totalAmount = computed(() => spendingReport.value.groups.reduce((sum, g) => sum + g.amount, 0))

const formatDisplayNumber = (num: number, type?: 'income' | 'expense' | 'transfer') => {
  if (settingsStore.isHideAmounts) {
    if (type) {
      return formatAmountForDisplay(num, type)
    }
    const symbol = formatBaseNumber(0).replace(/\d/g, '')
    return `${symbol}*****`
  }
  return formatBaseNumber(num)
}

const fetchSpending = async () => {
  const res = await reportsStore.getSpending(spendingParams.value)
  spendingReport.value = res
  updateSpendingChart()
}

const fetchTrend = async () => {
  const res = await reportsStore.getTrend({ months: parseInt(trendParams.value.months) })
  trendReport.value = res
  updateTrendChart()
}

const updateSpendingChart = () => {
  if (!spendingChart.value) return
  if (spendingChartInstance) spendingChartInstance.dispose()
  spendingChartInstance = echarts.init(spendingChart.value)
  const data = spendingReport.value.groups.slice(0, 10).reverse()
  spendingChartInstance.setOption({
    title: { text: 'Top Categories', left: 'center', textStyle: { fontSize: 13, color: 'inherit' } },
    tooltip: { trigger: 'item', formatter: '{b}: ¥{c} ({d}%)' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'value', axisLabel: { formatter: '¥{value}' } },
    yAxis: { type: 'category', data: data.map(d => d.name), inverse: true },
    series: [{ type: 'bar', data: data.map(d => ({ value: d.amount, itemStyle: { color: d.color } })), label: { show: true, position: 'right', formatter: '¥{c}' } }]
  })
}

const updateTrendChart = () => {
  if (!trendChart.value) return
  if (trendChartInstance) trendChartInstance.dispose()
  trendChartInstance = echarts.init(trendChart.value)
  const data = trendReport.value.series
  trendChartInstance.setOption({
    title: { text: 'Monthly Income vs Expense', left: 'center', textStyle: { fontSize: 13, color: 'inherit' } },
    tooltip: { trigger: 'axis' },
    legend: { data: ['Income', 'Expense', 'Net'], bottom: 0, textStyle: { fontSize: 11 } },
    grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true },
    xAxis: { type: 'category', data: data.map(d => d.month) },
    yAxis: { type: 'value', axisLabel: { formatter: '¥{value}' } },
    series: [
      { name: 'Income', type: 'line', data: data.map(d => d.income), smooth: true, itemStyle: { color: '#28a745' } },
      { name: 'Expense', type: 'line', data: data.map(d => d.expense), smooth: true, itemStyle: { color: '#dc3545' } },
      { name: 'Net', type: 'line', data: data.map(d => d.net), smooth: true, itemStyle: { color: '#4a90d9' } }
    ]
  })
}

const exportCSV = async () => {
  try {
    const params = { ...spendingParams.value, currency: selectedCurrency.value }
    const response = await api.get('/api/v1/reports/spending/export', { 
      params,
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `spending-report-${spendingParams.value.period}-${new Date().toISOString().split('T')[0]}.csv`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    uiStore.success('CSV exported successfully')
  } catch (e) {
    uiStore.error('Failed to export CSV')
  }
}

const exportExcel = async () => {
  try {
    const params = { ...spendingParams.value, currency: selectedCurrency.value }
    const response = await reportsStore.getSpending(params)
    const groups = Array.isArray(response.groups) ? response.groups : []
    const headers = ['Category', 'Amount', 'Count', '%']
    const totalAmount = groups.reduce((sum, g) => sum + (g.amount || 0), 0)
    const rows = groups.map((g: any) => [
      g.name || '',
      g.amount || 0,
      g.count || 0,
      totalAmount > 0 ? (((g.amount || 0) / totalAmount) * 100).toFixed(1) + '%' : '0%'
    ])
    exportToExcel(headers, rows, `spending-report-${spendingParams.value.period}-${new Date().toISOString().split('T')[0]}`)
    uiStore.success('Excel exported successfully')
  } catch (e) {
    uiStore.error('Failed to export Excel')
  }
}

onMounted(async () => {
  await accountsStore.fetchAccounts()
  
  const now = new Date()
  spendingParams.value.date_to = now.toISOString().split('T')[0]
  const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1)
  spendingParams.value.date_from = lastMonth.toISOString().split('T')[0]
  
  await Promise.all([fetchSpending(), fetchTrend()])
})

onMounted(() => {
  window.addEventListener('resize', () => {
    spendingChartInstance?.resize()
    trendChartInstance?.resize()
  })
})
</script>