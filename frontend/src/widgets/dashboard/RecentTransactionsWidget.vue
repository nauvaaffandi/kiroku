<template>
  <section class="recent-transactions">
    <div class="section-header">
      <h2>{{ t('dashboard.recentTransactions') }}</h2>
      <router-link to="/transactions" class="view-all">{{ t('dashboard.viewAll') }}</router-link>
    </div>
    <div class="transaction-list" v-if="recentTransactions.length > 0">
      <div class="transaction-item" v-for="tx in recentTransactions" :key="tx.id">
        <div class="tx-icon" :style="{ backgroundColor: getCategoryColor(tx.category_id) }">
          <span v-if="getCategoryIcon(tx.category_id)">{{ getCategoryIcon(tx.category_id) }}</span>
          <span v-else class="default-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
            </svg>
          </span>
        </div>
        <div class="tx-details">
          <p class="tx-name">{{ getCategoryName(tx.category_id) || tx.payee_payer || 'Transaction' }}</p>
          <p class="tx-meta">{{ formatDate(tx.date_time) }} • {{ getAccountName(tx.from_account_id || tx.to_account_id) }}</p>
        </div>
        <div class="tx-amount" :class="tx.type === 'income' ? 'text-success' : 'text-danger'">
          {{ tx.type === 'income' ? '+' : '-' }}{{ formatCurrency(tx.amount) }}
        </div>
      </div>
    </div>
    <div v-else class="empty-state">
      <p>{{ t('dashboard.noTransactions') }}</p>
      <router-link to="/transactions" class="btn-primary">{{ t('common.add') }} {{ t('nav.transactions') }}</router-link>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import type { Transaction, Category, Account } from '@/types'

const { t } = useI18n()

const props = defineProps<{
  recentTransactions: Transaction[]
  categories: Category[]
  accounts: Account[]
}>()

const getCategoryName = (id: string | null) => {
  if (!id) return ''
  const cat = props.categories.find(c => c.id === id)
  return cat?.name || ''
}

const getCategoryIcon = (id: string | null) => {
  if (!id) return ''
  const cat = props.categories.find(c => c.id === id)
  return cat?.icon || ''
}

const getCategoryColor = (id: string | null) => {
  if (!id) return 'var(--color-neutral-300)'
  const cat = props.categories.find(c => c.id === id)
  return cat?.color || 'var(--color-neutral-300)'
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
.recent-transactions {
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 1.5rem;
}
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1rem; }
.section-header h2 { font-size: 1rem; font-weight: 600; color: var(--color-neutral-900); }
.view-all { color: var(--color-primary); font-size: 0.875rem; font-weight: 500; }
.transaction-list { display: flex; flex-direction: column; }
.transaction-item { display: flex; align-items: center; gap: 1rem; padding: 1rem 0; border-bottom: 1px solid var(--color-neutral-200); }
.transaction-item:last-child { border-bottom: none; }
.tx-icon { width: 2.5rem; height: 2.5rem; border-radius: var(--radius-sm); display: flex; align-items: center; justify-content: center; color: white; font-size: 1.25rem; }
.tx-icon .default-icon { color: var(--color-neutral-400); }
.tx-details { flex: 1; min-width: 0; }
.tx-name { font-weight: 500; color: var(--color-neutral-900); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.tx-meta { font-size: 0.75rem; color: var(--color-neutral-500); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.tx-amount { font-weight: 600; font-size: 0.875rem; }
.text-success { color: var(--color-success); }
.text-danger { color: var(--color-danger); }
.empty-state { text-align: center; padding: 2rem; color: var(--color-neutral-500); }
.btn-primary { display: inline-block; margin-top: 1rem; padding: 0.625rem 1.25rem; background-color: var(--color-primary); color: white; border-radius: var(--radius-sm); font-weight: 500; text-decoration: none; }
</style>