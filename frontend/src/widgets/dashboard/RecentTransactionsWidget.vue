<template>
  <section class="card">
    <div class="section-header">
      <h2 class="section-title">{{ t('dashboard.recentTransactions') }}</h2>
      <router-link to="/transactions" class="action-link">{{ t('dashboard.viewAll') }}</router-link>
    </div>
    <div v-if="recentTransactions.length > 0">
      <div v-for="tx in recentTransactions" :key="tx.id" class="transaction-row">
        <div class="category-avatar" :style="{ backgroundColor: getCategoryColor(tx.category_id) }">
          <span v-if="getCategoryIcon(tx.category_id)">{{ getCategoryIcon(tx.category_id) }}</span>
          <span v-else class="text-muted">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
            </svg>
          </span>
        </div>
        <div class="flex-expand">
          <p class="payee-name">{{ getCategoryName(tx.category_id) || translateName(tx.payee_payer || '') || 'Transaction' }}</p>
          <p class="meta-text">{{ formatDate(tx.date_time) }} • {{ getAccountName(tx.from_account_id || tx.to_account_id) }}</p>
        </div>
        <div :class="['amount', tx.type === 'income' ? 'income-text' : 'expense-text']">
          {{ tx.type === 'income' ? '+' : '-' }}{{ formatCurrency(tx.amount) }}
        </div>
      </div>
    </div>
    <div v-else class="empty-state">
      <p class="empty-title">{{ t('dashboard.noTransactions') }}</p>
      <router-link to="/transactions" class="btn-add-transaction">{{ t('common.add') }} {{ t('nav.transactions') }}</router-link>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { useCategoryName } from '@/utils/categoryName'
import type { Transaction, Category, Account } from '@/types'

const { t } = useI18n()
const { translate, translateName } = useCategoryName()

const props = defineProps<{
  recentTransactions: Transaction[]
  categories: Category[]
  accounts: Account[]
}>()

const getCategoryName = (id: string | null) => {
  if (!id) return ''
  const cat = props.categories.find(c => c.id === id)
  return cat ? translate(cat) : ''
}

const getCategoryIcon = (id: string | null) => {
  if (!id) return ''
  const cat = props.categories.find(c => c.id === id)
  return cat?.icon || ''
}

const getCategoryColor = (id: string | null) => {
  if (!id) return '#cbd5e1'
  const cat = props.categories.find(c => c.id === id)
  return cat?.color || '#cbd5e1'
}

const getAccountName = (id: string | null) => {
  if (!id) return ''
  const acc = props.accounts.find(a => a.id === id)
  return acc?.name || ''
}

const formatCurrency = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(Math.abs(num))
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('ja-JP', { month: 'short', day: 'numeric' })
}
</script>

<style scoped>
@reference "tailwindcss";

/* Card */
.card {
  @apply bg-white dark:bg-zinc-900 rounded-xl shadow-sm p-6 mb-6 border border-transparent dark:border-zinc-800;
}

/* Section header */
.section-header {
  @apply flex items-center justify-between mb-4;
}

/* Section title */
.section-title {
  @apply text-base font-bold tracking-tight text-zinc-900 dark:text-zinc-100;
}

/* Action link */
.action-link {
  @apply text-blue-600 dark:text-blue-400 text-sm font-semibold hover:underline;
}

/* Transaction row */
.transaction-row {
  @apply flex items-center gap-4 py-3 border-b border-zinc-100 dark:border-zinc-800 last:border-b-0;
}

/* Category avatar */
.category-avatar {
  @apply w-10 h-10 rounded-lg flex items-center justify-center text-white text-lg shrink-0 shadow-xs;
}

/* Icon fallback text */
.text-muted {
  @apply text-zinc-300 dark:text-zinc-200;
}

/* Payee name */
.payee-name {
  @apply font-semibold text-zinc-900 dark:text-zinc-100 truncate leading-tight;
}

/* Meta text */
.meta-text {
  @apply text-xs text-zinc-500 dark:text-zinc-400 truncate font-medium;
}

/* Flex grow container */
.flex-expand {
  @apply flex-1 min-w-0;
}

/* Income/expense text */
.income-text {
  @apply text-emerald-600 dark:text-emerald-400;
}
.expense-text {
  @apply text-rose-600 dark:text-rose-400;
}

/* Empty state */
.empty-state {
  @apply text-center py-8 text-zinc-500 dark:text-zinc-300 font-medium;
}

/* Empty state title */
.empty-title {
  @apply text-base font-medium;
}

/* Add transaction button */
.btn-add-transaction {
  @apply inline-block mt-4 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium text-sm transition-colors;
}

/* Amount text */
.amount {
  @apply font-bold text-sm shrink-0 leading-none;
}
</style>