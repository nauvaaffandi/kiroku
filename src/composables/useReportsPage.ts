import { ref, onMounted, computed, type Ref } from 'vue'
import { useReportsStore } from '@/stores/reports'
import { useAccountsStore } from '@/stores/accounts'
import { useUIStore } from '@/stores/ui'
import { useSettingsStore } from '@/stores/settings'
import { useCurrencies } from '@/composables/useCurrencies'
import { formatAmountForDisplay, formatNumber as formatBaseNumber } from '@/utils/amountDisplay'
import { exportToExcel } from '@/utils/exportExcel'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([BarChart, LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

export function useReportsPage() {
  const reportsStore = useReportsStore()
  const accountsStore = useAccountsStore()
  const uiStore = useUIStore()
  const settingsStore = useSettingsStore()

  const spendingChart = ref<HTMLElement | null>(null)
  const trendChart = ref<HTMLElement | null>(null)
  let spendingChartInstance: echarts.ECharts | null = null
  let trendChartInstance: echarts.ECharts | null = null

  const spendingParams = ref({
    period: 'this_month', group_by: 'category', account_id: '', date_from: '', date_to: ''
  })
  const trendParams = ref({ months: '12', account_id: '', category_id: '' })
  const { currencies } = useCurrencies()
  const selectedCurrency = ref(settingsStore.baseCurrency)
  const spendingReport = ref({ period: '', groups: [] as any[] })
  const trendReport = ref({ series: [] as any[] })
  const totalAmount = computed(() => (spendingReport.value.groups ?? []).reduce((sum, g) => sum + g.amount, 0))

  const formatDisplayNumber = (num: number, type?: 'income' | 'expense' | 'transfer') => {
    if (settingsStore.isHideAmounts) {
      if (type) return formatAmountForDisplay(num, type)
      const symbol = formatBaseNumber(0).replace(/\d/g, '')
      return `${symbol}*****`
    }
    return formatBaseNumber(num)
  }

  const updateSpendingChart = () => {
    if (!spendingChart.value) return
    if (spendingChartInstance) spendingChartInstance.dispose()
    spendingChartInstance = echarts.init(spendingChart.value)
    const data = (spendingReport.value.groups ?? []).slice(0, 10).reverse()
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

  const fetchSpending = async () => {
    spendingReport.value = await reportsStore.getSpending(spendingParams.value)
    updateSpendingChart()
  }

  const fetchTrend = async () => {
    trendReport.value = await reportsStore.getTrend({ months: parseInt(trendParams.value.months) })
    updateTrendChart()
  }

  const exportCSV = async () => {
    try {
      const groups = Array.isArray(spendingReport.value.groups) ? spendingReport.value.groups : []
      const headers = ['Category', 'Amount', 'Count', '%']
      const sum = groups.reduce((s: number, g: any) => s + (g.amount || 0), 0)
      const rows = groups.map((g: any) => [
        g.name || '', g.amount || 0, g.count || 0,
        sum > 0 ? (((g.amount || 0) / sum) * 100).toFixed(1) + '%' : '0%'
      ])
      const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
      const url = window.URL.createObjectURL(new Blob([csv], { type: 'text/csv' }))
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
      const sum = groups.reduce((s: number, g: any) => s + (g.amount || 0), 0)
      const rows = groups.map((g: any) => [
        g.name || '', g.amount || 0, g.count || 0,
        sum > 0 ? (((g.amount || 0) / sum) * 100).toFixed(1) + '%' : '0%'
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
    window.addEventListener('resize', () => {
      spendingChartInstance?.resize()
      trendChartInstance?.resize()
    })
  })

  return {
    spendingParams, trendParams, accountsStore, spendingReport, trendReport,
    totalAmount, formatDisplayNumber, fetchSpending, fetchTrend, exportCSV, exportExcel,
    currencies, selectedCurrency, spendingChart, trendChart
  }
}

export type ReportsPage = ReturnType<typeof useReportsPage>
export type ChartRef = Ref<HTMLElement | null>
