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
import { api } from '@/api/client'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import BalanceTrendWidget from '@/widgets/analytics/BalanceTrendWidget.vue'
import BalanceByAccountsWidget from '@/widgets/analytics/BalanceByAccountsWidget.vue'
import BalanceByCurrenciesWidget from '@/widgets/analytics/BalanceByCurrenciesWidget.vue'
import type { BalanceTrendPoint, AccountBalanceSummary, CurrencyBalanceSummary } from '@/types/report'

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
    const res = await api.get<{ series: BalanceTrendPoint[] }>(
      `/api/v1/analytics/balance/trend?date_from=${dateFrom.value}&date_to=${dateTo.value}&interval=${trendInterval.value}`
    )
    trendData.value = res.data.series
  } catch (e) {
    console.error('Failed to fetch balance trend:', e)
  } finally {
    trendLoading.value = false
  }
}

async function fetchBalanceByAccounts() {
  accountsLoading.value = true
  try {
    const res = await api.get<{ accounts: AccountBalanceSummary[]; total_by_currency: Record<string, number> }>(
      `/api/v1/analytics/balance/by-accounts?date_from=${dateFrom.value}&date_to=${dateTo.value}`
    )
    accountsData.value = res.data.accounts
  } catch (e) {
    console.error('Failed to fetch balance by accounts:', e)
  } finally {
    accountsLoading.value = false
  }
}

async function fetchBalanceByCurrencies() {
  currenciesLoading.value = true
  try {
    const res = await api.get<{ currencies: CurrencyBalanceSummary[] }>(
      `/api/v1/analytics/balance/by-currencies?date_from=${dateFrom.value}&date_to=${dateTo.value}`
    )
    currenciesData.value = res.data.currencies
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