<template>
  <section class="recent-card">
    <div class="section-header">
      <h2 class="section-title">{{ t('dashboard.recentTransactions') }}</h2>
      <router-link to="/transactions" class="action-link">{{ t('dashboard.viewAll') }}</router-link>
    </div>
    <div v-if="recentTransactions.length > 0" class="transactions-list">
      <div v-for="tx in recentTransactions" :key="tx.id" class="transaction-row">
        <div class="category-dot-wrapper">
          <span
            v-if="getCategoryIcon(tx.category_id)"
            class="category-emoji"
          >{{ getCategoryIcon(tx.category_id) }}</span>
          <span
            v-else
            class="category-dot"
            :style="{ backgroundColor: getCategoryColor(tx.category_id) }"
          />
        </div>
        <div class="row-details">
          <p class="payee-name">{{ getCategoryName(tx.category_id) || translateName(tx.payee_payer || '') || t('common.account') }}</p>
          <p class="meta-text">{{ formatDate(tx.date_time) }} · {{ getAccountName(tx.from_account_id || tx.to_account_id) }}</p>
        </div>
        <div :class="['row-amount', tx.type === 'income' ? 'income' : 'expense']">
          {{ formatRowAmount(tx.amount, tx.type) }}
        </div>
      </div>
    </div>
    <div v-else class="empty-state">
      <p class="empty-title">{{ t('dashboard.noTransactions') }}</p>
      <router-link to="/transactions?create=true" class="btn-add-transaction">
        {{ t('common.addTransaction') }}
      </router-link>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { useCategoryName } from '@/utils/categoryName'
import { formatAmountForDisplay } from '@/utils/amountDisplay'
import { useSettingsStore } from '@/stores/settings'
import type { Transaction, Category, Account } from '@/types'
import dayjs from 'dayjs'

const { t } = useI18n()
const { translate, translateName } = useCategoryName()
const settingsStore = useSettingsStore()

const props = defineProps<{
  recentTransactions: Transaction[]
  categories: Category[]
  accounts: Account[]
}>()

console.log('[RecentTransactionsWidget] rendered transactions', {
  count: props.recentTransactions.length
})

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
  if (!id) return 'var(--border-secondary)'
  const cat = props.categories.find(c => c.id === id)
  return cat?.color || 'var(--border-secondary)'
}

const getAccountName = (id: string | null) => {
  if (!id) return ''
  const acc = props.accounts.find(a => a.id === id)
  return acc?.name || ''
}

const formatRowAmount = (amount: number, type: 'income' | 'expense' | 'transfer') => {
  return formatAmountForDisplay(amount, type)
}

const formatDate = (dateStr: string) => {
  const fmt = settingsStore.dateFormat || 'YYYY-MM-DD'
  return dayjs(dateStr).format(fmt)
}
</script>

<style scoped>
@reference "tailwindcss";

.recent-card {
  @apply rounded-[14px] border p-4 mb-6 bg-[var(--bg-card)] border-[var(--border-primary)] text-[var(--text-primary)];
}

.section-header {
  @apply flex items-center justify-between mb-4;
}

.section-title {
  @apply text-base font-semibold text-[var(--text-primary)];
}

.action-link {
  @apply text-xs font-medium no-underline text-[var(--accent)];
}

.action-link:hover {
  @apply text-[var(--accent-hover)] underline;
}

.transactions-list {
  @apply flex flex-col;
}

.transaction-row {
  @apply flex items-center gap-3 py-3;
  min-height: 56px;
  border-bottom: 1px solid var(--border-primary);
}

.transaction-row:last-child {
  border-bottom: none;
}

.category-dot-wrapper {
  @apply w-5 h-5 flex items-center justify-center shrink-0;
}

.category-dot {
  @apply w-2 h-2 rounded-full;
}

.category-emoji {
  @apply text-base leading-none;
}

.row-details {
  @apply flex-1 min-w-0;
}

.payee-name {
  @apply text-sm font-medium truncate text-[var(--text-primary)];
}

.meta-text {
  @apply text-xs text-[var(--text-muted)] truncate;
}

.row-amount {
  @apply text-[15px] font-semibold shrink-0 tabular-nums text-right text-[var(--text-primary)];
}

.row-amount.income {
  @apply text-[var(--success)];
}

.empty-state {
  @apply text-center py-8;
}

.empty-title {
  @apply text-sm font-medium mb-3 text-[var(--text-muted)];
}

.btn-add-transaction {
  @apply inline-flex items-center justify-center px-4 py-2 text-sm font-medium rounded-md text-[var(--text-inverse)] bg-[var(--text-primary)] transition-opacity;
  min-height: 44px;
}

.btn-add-transaction:hover {
  @apply opacity-90;
}
</style>
