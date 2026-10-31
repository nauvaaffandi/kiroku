<template>
  <div class="dashboard-page">
    <main class="dashboard-main">
      <!-- Monthly Summary — integrated into dashboard, not a widget -->
      <section class="dashboard-summary" aria-labelledby="summary-heading">
        <h2 id="summary-heading" class="dashboard-heading">{{ t('dashboard.summary') }}</h2>
        <div class="summary-grid">
          <div class="summary-item">
            <p class="summary-label">{{ t('dashboard.monthlyIncome') }}</p>
            <p class="summary-value income">{{ formatIncome(monthlyIncome) }}</p>
          </div>
          <div class="summary-item">
            <p class="summary-label">{{ t('dashboard.monthlyExpense') }}</p>
            <p class="summary-value expense">{{ formatExpense(monthlyExpense) }}</p>
          </div>
          <div class="summary-item">
            <p class="summary-label">{{ t('dashboard.net') }}</p>
            <p class="summary-value" :class="netClass">{{ formatNet(monthlyNet) }}</p>
          </div>
          <div class="summary-item">
            <p class="summary-label">{{ t('dashboard.totalBalance') }}</p>
            <p class="summary-value balance">{{ formatBalance(totalBalance) }}</p>
          </div>
        </div>
      </section>

      <!-- Accounts — continuous flow below summary -->
      <section class="dashboard-accounts" aria-labelledby="accounts-heading">
        <div class="accounts-header">
          <h2 id="accounts-heading" class="dashboard-heading">{{ t('dashboard.accounts') }}</h2>
          <router-link to="/accounts" class="view-all">{{ t('dashboard.viewAll') }}</router-link>
        </div>

        <div v-if="accountsStore.accounts.length === 0" class="accounts-empty">
          <p class="accounts-empty-text">{{ t('dashboard.noAccounts') }}</p>
          <router-link to="/accounts" class="accounts-empty-link">{{ t('accounts.addAccount') }}</router-link>
        </div>

        <div v-else class="accounts-list">
          <div
            v-for="account in accountsStore.accounts"
            :key="account.id"
            class="account-row"
            :class="{ 'is-archived': account.archived }"
          >
            <div class="account-info">
              <div class="account-dot" v-if="account.color" :style="{ backgroundColor: account.color }"></div>
              <div class="account-details">
                <span class="account-name">{{ account.name }}</span>
                <span class="account-type">{{ tAccountType(account.type) }}</span>
                <span v-if="account.archived" class="account-archived">{{ t('categories.filterArchived') }}</span>
              </div>
            </div>
            <div class="account-balance">
              <span class="account-amount">{{ formatAccountBalance(account.balance, account.currency) }}</span>
              <span class="account-currency">{{ account.currency }}</span>
            </div>
          </div>
        </div>
      </section>

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
import { useSettingsStore } from '@/stores/settings'
import { useI18n } from '@/locales'
import { useCategoryName } from '@/utils/categoryName'
import { getCurrencySymbol } from '@/utils/amountDisplay'
import RecentTransactionsWidget from '@/widgets/dashboard/RecentTransactionsWidget.vue'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([BarChart, LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

const { t } = useI18n()
const { translate } = useCategoryName()

const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const categoriesStore = useCategoriesStore()
const settingsStore = useSettingsStore()

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

const isHideAmounts = computed(() => settingsStore.isHideAmounts)
const baseCurrency = computed(() => settingsStore.baseCurrency)
const numberFormat = computed(() => settingsStore.numberFormat)
const currencySymbol = computed(() => getCurrencySymbol(baseCurrency.value))

const netClass = computed(() => {
  if (monthlyNet.value > 0) return 'income'
  if (monthlyNet.value < 0) return 'expense'
  return 'neutral'
})

const formatNumber = (value: number): string => {
  return new Intl.NumberFormat(numberFormat.value).format(Math.abs(value))
}

const maskAmount = (sign: string): string => {
  return `${sign}${currencySymbol.value}*****`
}

const formatIncome = (value: number): string => {
  if (isHideAmounts.value) return maskAmount('+')
  return `+${currencySymbol.value}${formatNumber(value)}`
}

const formatExpense = (value: number): string => {
  if (isHideAmounts.value) return maskAmount('−')
  return `−${currencySymbol.value}${formatNumber(value)}`
}

const formatNet = (value: number): string => {
  if (isHideAmounts.value) {
    const sign = value >= 0 ? '+' : '−'
    return maskAmount(sign)
  }
  const sign = value >= 0 ? '+' : '−'
  return `${sign}${currencySymbol.value}${formatNumber(value)}`
}

const formatBalance = (value: number): string => {
  if (isHideAmounts.value) {
    const sign = value < 0 ? '−' : ''
    return `${sign}${currencySymbol.value}*****`
  }
  const sign = value < 0 ? '−' : '+'
  return `${sign}${currencySymbol.value}${formatNumber(value)}`
}

const tAccountType = (type: string): string => {
  const keyMap: Record<string, string> = {
    checking: 'accounts.typeChecking',
    savings: 'accounts.typeSavings',
    credit: 'accounts.typeCredit',
    cash: 'accounts.typeCash',
    investment: 'accounts.typeInvestment',
    other: 'accounts.typeOther'
  }
  const key = keyMap[type]
  return key ? t(key) : type
}

const formatAccountBalance = (balance: number, currency: string): string => {
  const symbol = getCurrencySymbol(currency)
  if (isHideAmounts.value) {
    const sign = balance < 0 ? '−' : ''
    return `${sign}${symbol}*****`
  }
  const formatted = new Intl.NumberFormat(numberFormat.value).format(Math.abs(balance))
  const sign = balance < 0 ? '−' : ''
  return `${sign}${symbol}${formatted}`
}

const formatCurrency = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(Math.abs(num))
}

const initCharts = async () => {
  if (spendingChart.value && categoriesStore.categories.length > 0) {
    spendingChartInstance = echarts.init(spendingChart.value)
    const expenseCats = categoriesStore.categories
      .slice(0, 8)

    const data = expenseCats.map(c => {
      const total = (transactionsStore.transactions || [])
        .filter(t => t.category_id === c.id)
        .reduce((sum, t) => sum + t.amount, 0)
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
.dashboard-main { @apply flex-1 px-4 py-6 mx-auto w-full max-w-[1200px]; }

/* Monthly Summary — integrated into dashboard, no widget card mindset */
.dashboard-summary {
  @apply rounded-[14px] border p-6 mb-4;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  color: var(--text-primary);
}

.dashboard-heading {
  @apply text-xs font-medium m-0 block mb-4;
  color: var(--text-muted);
}

.summary-grid {
  @apply grid gap-4;
  grid-template-columns: repeat(2, 1fr);
}

.summary-item {
  @apply flex flex-col gap-1;
}

.summary-label {
  @apply text-xs font-medium m-0 block mb-0.5;
  color: var(--text-muted);
}

.summary-value {
  @apply text-[24px] font-semibold leading-tight m-0 tabular-nums text-right;
}

.summary-value.income {
  color: var(--success);
}

.summary-value.expense {
  color: var(--text-primary);
}

.summary-value.balance {
  color: var(--text-primary);
}

/* Accounts — continuous flow below summary */
.dashboard-accounts {
  @apply rounded-[14px] border p-6 mb-6;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  color: var(--text-primary);
}

.accounts-header {
  @apply flex items-center justify-between mb-4;
}

.view-all {
  @apply text-xs font-medium no-underline;
  color: var(--accent);
}

.view-all:hover {
  color: var(--accent-hover);
  text-decoration: underline;
}

.accounts-empty {
  @apply text-center py-8;
}

.accounts-empty-text {
  @apply text-sm mb-4 m-0;
  color: var(--text-muted);
}

.accounts-empty-link {
  @apply text-sm font-medium no-underline;
  color: var(--accent);
}

.accounts-empty-link:hover {
  color: var(--accent-hover);
  text-decoration: underline;
}

.accounts-list {
  @apply flex flex-col;
}

.account-row {
  @apply flex items-center py-3 gap-3;
  border-bottom: 1px solid var(--border-primary);
}

.account-row:last-child {
  border-bottom: none;
}

.account-row.is-archived {
  opacity: 0.6;
}

.account-info {
  @apply flex items-center gap-2 flex-1 min-w-0;
}

.account-dot {
  @apply w-2 h-2 rounded-full shrink-0;
}

.account-details {
  @apply flex flex-col min-w-0;
}

.account-name {
  @apply text-sm font-medium truncate;
  color: var(--text-primary);
}

.account-type {
  @apply text-xs truncate;
  color: var(--text-muted);
}

.account-archived {
  @apply text-[10px] font-medium ml-1;
  color: var(--error);
}

.account-balance {
  @apply flex flex-col items-end shrink-0;
}

.account-amount {
  @apply text-sm font-semibold tabular-nums;
  color: var(--text-primary);
}

.account-currency {
  @apply text-xs;
  color: var(--text-muted);
}

/* Charts */
.charts-section { @apply grid grid-cols-[repeat(auto-fit,minmax(400px,1fr))] gap-6 mb-6; }
@media (max-width: 768px) { .charts-section { @apply grid-cols-1; } }
.chart-card { @apply rounded-xl shadow-sm p-[1.5rem]; background-color: var(--bg-card); }
.chart-card h2 { @apply text-base font-semibold mb-4; color: var(--text-primary); }
.chart-container { @apply w-full h-[300px]; }
</style>