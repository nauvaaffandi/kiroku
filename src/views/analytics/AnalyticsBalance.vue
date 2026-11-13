<template>
  <div class="analytics-balance-view">
    <div class="page-header">
      <h1>{{ t('balanceAnalytics.title') }}</h1>
      <p class="page-subtitle">{{ t('balanceAnalytics.subtitle') }}</p>
    </div>

    <DateFilterBar @change="onDateChange" />

    <div class="widgets-grid">
      <BalanceTrendWidget
        v-model:interval="trendInterval"
        :loading="trendLoading"
        :data="trendData"
        @update:interval="fetchBalanceTrend"
      />

      <BalanceByAccountsWidget
        :loading="accountsLoading"
        :data="accountsData"
      />

      <BalanceByCurrenciesWidget
        :loading="currenciesLoading"
        :data="currenciesData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from '@/locales'
import { localDb } from '@/storage/repository'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import BalanceTrendWidget from '@/widgets/analytics/BalanceTrendWidget.vue'
import BalanceByAccountsWidget from '@/widgets/analytics/BalanceByAccountsWidget.vue'
import BalanceByCurrenciesWidget from '@/widgets/analytics/BalanceByCurrenciesWidget.vue'
import type { BalanceTrendPoint, AccountBalanceSummary, CurrencyBalanceSummary } from '@/types/report'
import type { Account, Transaction } from '@/types'

const { t } = useI18n()

const dateFrom = ref('')
const dateTo = ref('')
const trendInterval = ref('month')

const trendLoading = ref(true)
const accountsLoading = ref(true)
const currenciesLoading = ref(true)

const trendData = ref<BalanceTrendPoint[]>([])
const accountsData = ref<AccountBalanceSummary[]>([])
const currenciesData = ref<CurrencyBalanceSummary[]>([])

function onDateChange(payload: { dateFrom: string; dateTo: string }) {
  dateFrom.value = payload.dateFrom
  dateTo.value = payload.dateTo
  fetchAll()
}

async function fetchAll() {
  await Promise.all([
    fetchBalanceTrend(),
    fetchBalanceByAccounts(),
    fetchBalanceByCurrencies(),
  ])
}

async function fetchBalanceTrend() {
  trendLoading.value = true
  try {
    const accounts = await localDb.list<Account>('account')
    const transactions = await localDb.list<Transaction>('transaction')
    const now = new Date()
    const months: string[] = []
    const count = trendInterval.value === 'year' ? 12 : 6
    for (let i = count - 1; i >= 0; i--) {
      const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
      months.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
    }
    const series: BalanceTrendPoint[] = months.map((period) => {
      const inMonth = transactions.filter((t) => t.date_time.slice(0, 7) === period)
      const net = inMonth.reduce((s, t) => t.type === 'income' ? s + t.amount : t.type === 'expense' ? s - t.amount : s, 0)
      return {
        period,
        opening_balance: 0,
        closing_balance: accounts.reduce((s, a) => s + a.balance, 0) + net,
        net_change: net,
        currency: accounts[0]?.currency ?? 'USD'
      }
    })
    trendData.value = series
  } catch (e) {
    console.error('Failed to compute balance trend:', e)
  } finally {
    trendLoading.value = false
  }
}

async function fetchBalanceByAccounts() {
  accountsLoading.value = true
  try {
    const accounts = await localDb.list<Account>('account')
    accountsData.value = accounts.map((a) => ({
      id: Number(a.id) || 0,
      name: a.name,
      currency: a.currency,
      opening_balance: a.balance,
      closing_balance: a.balance,
      net_change: 0,
      color: a.color
    }))
  } catch (e) {
    console.error('Failed to fetch balance by accounts:', e)
  } finally {
    accountsLoading.value = false
  }
}

async function fetchBalanceByCurrencies() {
  currenciesLoading.value = true
  try {
    const accounts = await localDb.list<Account>('account')
    const byCurrency = new Map<string, { balance: number }>()
    for (const a of accounts) {
      const entry = byCurrency.get(a.currency) ?? { balance: 0 }
      entry.balance += a.balance
      byCurrency.set(a.currency, entry)
    }
    currenciesData.value = [...byCurrency.entries()].map(([currency, entry]) => ({
      currency,
      opening_balance: entry.balance,
      closing_balance: entry.balance,
      net_change: 0,
      fx_rate_to_base: 1,
      converted_total: entry.balance
    }))
  } catch (e) {
    console.error('Failed to fetch balance by currencies:', e)
  } finally {
    currenciesLoading.value = false
  }
}

onMounted(async () => {
  const now = new Date()
  dateTo.value = now.toISOString().split('T')[0]
  const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1)
  dateFrom.value = lastMonth.toISOString().split('T')[0]
  await fetchAll()
})
</script>

<style scoped>
@reference "tailwindcss";
.analytics-balance-view {
  @apply max-w-[80rem] mx-auto p-4 flex flex-col gap-4 w-full max-w-full overflow-x-hidden box-border;
}
.page-header { 
  @apply flex flex-col gap-1; 
}
.page-header h1 { 
  @apply text-xl font-bold; 
  color: var(--text-primary); 
}
.page-subtitle { 
  color: var(--text-muted); 
}
.widgets-grid { 
  @apply grid gap-4 w-full; 
  grid-template-columns: minmax(0, 1fr);
}
@media (min-width: 768px) {
  .widgets-grid {
    grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
  }
}
</style>