<template>
  <div class="reports-page">
    <header class="page-header">
      <h1>Reports</h1>
      <div class="header-actions">
        <select v-model="selectedCurrency" @change="fetchSpending" class="form-input" style="width: auto;">
          <option v-for="c of currencies" :key="c" :value="c">{{ c }}</option>
        </select>
        <button @click="exportCSV" class="btn-secondary">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
          Export CSV
        </button>
      </div>
    </header>

    <div class="reports-content">
      <aside class="reports-sidebar">
        <section class="sidebar-section">
          <h3>Spending Report</h3>
          <div class="form-group">
            <label>Period</label>
            <select v-model="spendingParams.period" @change="fetchSpending" class="form-input">
              <option value="this_week">This Week</option>
              <option value="this_month">This Month</option>
              <option value="this_year">This Year</option>
              <option value="last_month">Last Month</option>
              <option value="last_year">Last Year</option>
            </select>
          </div>
          <div class="form-group">
            <label>Group By</label>
            <select v-model="spendingParams.group_by" @change="fetchSpending" class="form-input">
              <option value="category">Category</option>
              <option value="parent_category">Parent Category</option>
            </select>
          </div>
          <div class="form-group">
            <label>Account (optional)</label>
            <select v-model="spendingParams.account_id" @change="fetchSpending" class="form-input">
              <option value="">All Accounts</option>
              <option v-for="acc in accountsStore.accounts" :key="acc.id" :value="acc.id">
                {{ acc.name }}
              </option>
            </select>
          </div>
        </section>

        <section class="sidebar-section">
          <h3>Trend Report</h3>
          <div class="form-group">
            <label>Months</label>
            <select v-model="trendParams.months" @change="fetchTrend" class="form-input">
              <option value="6">6 Months</option>
              <option value="12" selected>12 Months</option>
              <option value="24">24 Months</option>
            </select>
          </div>
        </section>
      </aside>

      <main class="reports-main">
        <div class="report-tabs">
          <div class="tab-buttons">
            <button :class="['tab-btn', { active: activeTab === 'spending' }]" @click="activeTab = 'spending'">支出レポート</button>
            <button :class="['tab-btn', { active: activeTab === 'trend' }]" @click="activeTab = 'trend'">月次トレンド</button>
            <button :class="['tab-btn', { active: activeTab === 'cashflow' }]" @click="activeTab = 'cashflow'">キャッシュフローテーブル</button>
          </div>

          <div v-if="activeTab === 'spending'" class="tab-panel">
            <section class="report-section">
              <div class="section-header">
                <h2>Spending by Category</h2>
                <span class="period-label">{{ spendingParams.period }}</span>
              </div>
              <div class="chart-container" ref="spendingChart"></div>
              <div class="spending-table" v-if="spendingReport.groups">
                <table>
                  <thead>
                    <tr>
                      <th>Category</th>
                      <th class="text-right">Amount</th>
                      <th class="text-right">Count</th>
                      <th class="text-right">%</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="group in spendingReport.groups" :key="group.id">
                      <td>
                        <div class="category-cell">
                          <span class="category-dot" :style="{ backgroundColor: group.color }"></span>
                          {{ group.name }}
                        </div>
                        <div class="subcategories" v-if="group.children && group.children.length > 0">
                          <div v-for="child in group.children" :key="child.id" class="subcategory-row">
                            <span class="category-dot" :style="{ backgroundColor: child.color }"></span>
                            {{ child.name }}: {{ formatDisplayNumber(child.amount, 'expense') }}
                          </div>
                        </div>
                      </td>
                      <td class="text-right">{{ formatDisplayNumber(group.amount, 'expense') }}</td>
                      <td class="text-right">{{ group.count }}</td>
                      <td class="text-right">{{ totalAmount > 0 ? ((group.amount / totalAmount) * 100).toFixed(1) : 0 }}%</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </section>
          </div>

          <div v-if="activeTab === 'trend'" class="tab-panel">
            <section class="report-section">
              <h2>Monthly Trend</h2>
              <div class="chart-container" ref="trendChart"></div>
            </section>
          </div>

          <div v-if="activeTab === 'cashflow'" class="tab-panel">
            <CashflowTable :date-from="spendingParams.date_from" :date-to="spendingParams.date_to" />
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useReportsStore } from '@/stores/reports'
import { useAccountsStore } from '@/stores/accounts'
import { useUIStore } from '@/stores/ui'
import { useSettingsStore } from '@/stores/settings'
import { useCurrencies } from '@/composables/useCurrencies'
import { formatAmountForDisplay, formatNumber as formatBaseNumber } from '@/utils/amountDisplay'
import { api } from '@/api/client'
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
  spendingChartInstance = echarts.init(spendingChart.value)
  const data = spendingReport.value.groups.slice(0, 10).reverse()
  spendingChartInstance.setOption({
    title: { text: 'Top Categories', left: 'center' },
    tooltip: { trigger: 'item', formatter: '{b}: ¥{c} ({d}%)' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'value', axisLabel: { formatter: '¥{value}' } },
    yAxis: { type: 'category', data: data.map(d => d.name), inverse: true },
    series: [{ type: 'bar', data: data.map(d => ({ value: d.amount, itemStyle: { color: d.color } })), label: { show: true, position: 'right', formatter: '¥{c}' } }]
  })
}

const updateTrendChart = () => {
  if (!trendChart.value) return
  trendChartInstance = echarts.init(trendChart.value)
  const data = trendReport.value.series
  trendChartInstance.setOption({
    title: { text: 'Monthly Income vs Expense', left: 'center' },
    tooltip: { trigger: 'axis' },
    legend: { data: ['Income', 'Expense', 'Net'], bottom: 0 },
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

onMounted(async () => {
  await accountsStore.fetchAccounts()
  
  // Set default date range for cashflow table
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

<style scoped>
.reports-page { min-height: 100vh; display: flex; flex-direction: column; }

.page-header { padding: 1.5rem; background-color: var(--bg-secondary); border-bottom: 1px solid var(--border-primary); }
.page-header h1 { font-size: 1.5rem; font-weight: 600; color: var(--text-primary); }

.header-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
}

.reports-content { display: grid; grid-template-columns: 280px 1fr; gap: 1.5rem; padding: 1.5rem; flex: 1; }

.reports-sidebar { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); box-shadow: var(--shadow); padding: 1.5rem; height: fit-content; position: sticky; top: 1.5rem; }

.sidebar-section { margin-bottom: 1.5rem; }
.sidebar-section h3 { font-size: 0.875rem; font-weight: 600; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 1rem; }

.form-group { margin-bottom: 1rem; }
.form-group label { display: block; font-size: 0.875rem; font-weight: 500; color: var(--text-secondary); margin-bottom: 0.375rem; }
.form-input { width: 100%; padding: 0.5rem 0.75rem; border: 1px solid var(--border-primary); border-radius: var(--radius-sm); background-color: var(--bg-card); color: var(--text-primary); font-size: 0.875rem; transition: border-color 0.2s ease, box-shadow 0.2s ease; }
.form-input:focus { outline: none; border-color: var(--accent); box-shadow: 0 0 0 2px var(--accent-light); }

.reports-main { display: flex; flex-direction: column; gap: 1.5rem; }

.report-tabs {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  overflow: hidden;
}

.tab-buttons {
  display: flex;
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
}

.tab-btn {
  flex: 1;
  padding: 0.75rem 1rem;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: color 0.2s ease, background-color 0.2s ease;
  border-bottom: 3px solid transparent;
}

.tab-btn:hover {
  color: var(--text-primary);
  background-color: var(--bg-hover);
}

.tab-btn.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
  background-color: var(--bg-card);
}

.tab-panel {
  padding: 1.5rem;
}

.report-section { background-color: var(--bg-secondary); border-radius: var(--radius); padding: 1.5rem; margin-bottom: 1rem; }

.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1rem; }
.section-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.period-label { font-size: 0.875rem; color: var(--text-secondary); background-color: var(--bg-card); padding: 0.25rem 0.75rem; border-radius: 9999px; }

.chart-container { width: 100%; height: 350px; margin-bottom: 1.5rem; }

.spending-table { overflow-x: auto; }
.spending-table table { width: 100%; border-collapse: collapse; font-size: 0.875rem; }
.spending-table th, .spending-table td { padding: 0.75rem 1rem; text-align: left; border-bottom: 1px solid var(--border-primary); }
.spending-table th { font-weight: 600; color: var(--text-secondary); background-color: var(--bg-card); }
.spending-table .text-right { text-align: right; }
.category-cell { display: flex; flex-direction: column; gap: 0.25rem; }
.subcategories { margin-left: 1.5rem; padding-left: 0.75rem; border-left: 2px solid var(--border-primary); }
.subcategory-row { display: flex; align-items: center; gap: 0.5rem; font-size: 0.8125rem; color: var(--text-secondary); padding: 0.25rem 0; }
.category-dot { width: 0.5rem; height: 0.5rem; border-radius: 50%; flex-shrink: 0; }

@media (max-width: 1024px) {
  .reports-content { grid-template-columns: 1fr; }
  .reports-sidebar { position: static; }
}
</style>