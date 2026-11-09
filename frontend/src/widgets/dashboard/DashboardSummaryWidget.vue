<template>
  <div class="dashboard-summary" aria-labelledby="summary-heading">
    <h2 id="summary-heading" class="dashboard-heading">{{ t('dashboard.summary') }}</h2>
    <div class="summary-grid">
      <div class="summary-item"><p class="summary-label">{{ t('dashboard.monthlyIncome') }}</p><p class="summary-value income">{{ formatIncome(monthlyIncome) }}</p></div>
      <div class="summary-item"><p class="summary-label">{{ t('dashboard.monthlyExpense') }}</p><p class="summary-value expense">{{ formatExpense(monthlyExpense) }}</p></div>
      <div class="summary-item"><p class="summary-label">{{ t('dashboard.net') }}</p><p class="summary-value" :class="netClass">{{ formatNet(monthlyNet) }}</p></div>
      <div class="summary-item"><p class="summary-label">{{ t('dashboard.totalBalance') }}</p><p class="summary-value balance">{{ formatBalance(totalBalance) }}</p></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from '@/locales'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import { useDashboardFormatting } from '@/widgets/dashboard/useDashboardFormatting'

const { t } = useI18n()
const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const { formatIncome, formatExpense, formatNet, formatBalance } = useDashboardFormatting()
const monthlyIncome = computed(() => {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), 1)
  return (transactionsStore.transactions || []).filter(t => t.type === 'income' && new Date(t.date_time) >= start).reduce((sum, t) => sum + t.amount, 0)
})
const monthlyExpense = computed(() => {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), 1)
  return (transactionsStore.transactions || []).filter(t => t.type === 'expense' && new Date(t.date_time) >= start).reduce((sum, t) => sum + t.amount, 0)
})
const monthlyNet = computed(() => monthlyIncome.value - monthlyExpense.value)
const totalBalance = computed(() => (accountsStore.accounts || []).reduce((sum, account) => sum + account.balance, 0))
const netClass = computed(() => monthlyNet.value > 0 ? 'income' : monthlyNet.value < 0 ? 'expense' : 'neutral')
</script>

<style scoped>
@reference "tailwindcss";
.dashboard-summary { @apply rounded-[14px] border p-6 mb-4; background-color: var(--bg-card); border-color: var(--border-primary); color: var(--text-primary); }
.dashboard-heading { @apply text-xs font-medium m-0 block mb-4; color: var(--text-muted); }
.summary-grid { @apply grid gap-4; grid-template-columns: repeat(2, 1fr); }
.summary-item { @apply flex flex-col gap-1; }
.summary-label { @apply text-xs font-medium m-0 block mb-0.5; color: var(--text-muted); }
.summary-value { @apply text-[24px] font-semibold leading-tight m-0 tabular-nums text-right; }
.summary-value.income { color: var(--success); }
.summary-value.expense, .summary-value.balance { color: var(--text-primary); }
</style>
