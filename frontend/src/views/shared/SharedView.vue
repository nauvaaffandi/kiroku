<template>
  <div class="shared-view-page">
    <header class="shared-header">
      <div class="shared-header-content">
        <div class="brand">
          <h1>Kiroku</h1>
          <span class="shared-badge">Shared Account</span>
        </div>
        <div class="attribution" v-if="data?.attribution">
          {{ data.attribution }}
        </div>
      </div>
    </header>

    <div class="shared-container" v-if="loading">
      <div class="loading-state">Loading shared account...</div>
    </div>

    <div class="shared-container" v-else-if="error">
      <div class="error-state">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
        <h2>Unable to load shared account</h2>
        <p>{{ error }}</p>
      </div>
    </div>

    <div class="shared-container" v-else-if="data">
      <div class="account-summary-card">
        <div class="account-meta">
          <span class="account-name">{{ data.account.name }}</span>
          <span class="account-currency">{{ data.account.currency }}</span>
        </div>
        <div class="account-balance" :class="{ 'text-danger': data.account.balance < 0 }">
          {{ formatCurrency(data.account.balance, data.account.currency) }}
        </div>
      </div>

      <div class="transactions-section">
        <h3>Transactions</h3>
        <div class="transactions-list" v-if="data.transactions && data.transactions.length > 0">
          <div class="transaction-item" v-for="tx in data.transactions" :key="tx.id">
            <div class="tx-left">
              <span class="tx-category-icon" :style="{ backgroundColor: tx.category_color || '#e2e8f0' }">
                {{ tx.category_icon || '📦' }}
              </span>
              <div class="tx-info">
                <span class="tx-payee">{{ tx.payee || tx.category || 'Transaction' }}</span>
                <span class="tx-note" v-if="tx.note">{{ tx.note }}</span>
                <span class="tx-date">{{ formatDate(tx.date) }}</span>
              </div>
            </div>
            <div class="tx-right">
              <span class="tx-amount" :class="tx.type">
                {{ tx.type === 'income' ? '+' : '-' }}{{ formatCurrency(tx.amount, data.account.currency) }}
              </span>
            </div>
          </div>
        </div>
        <div class="empty-transactions" v-else>
          <p>No transactions found in this account.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSettingsStore } from '@/stores/settings'
import { formatAmountForDisplay, getCurrencySymbol } from '@/utils/amountDisplay'
import axios from 'axios'

const route = useRoute()
const token = route.params.token as string
const settingsStore = useSettingsStore()

const data = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const fetchSharedData = async () => {
  loading.value = true
  error.value = null
  try {
    const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'
    const res = await axios.get(`${backendUrl}/shared/public/${token}`)
    data.value = res.data
  } catch (e: any) {
    error.value = e.response?.data || 'Link may be invalid or expired.'
  } finally {
    loading.value = false
  }
}

const formatCurrency = (amount: number, currency = 'JPY', type?: 'income' | 'expense' | 'transfer') => {
  if (settingsStore.isHideAmounts && type) {
    return formatAmountForDisplay(amount, type)
  }
  if (settingsStore.isHideAmounts && !type) {
    const symbol = getCurrencySymbol(currency)
    return `${symbol}*****`
  }
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency }).format(amount)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('ja-JP', { year: 'numeric', month: 'short', day: 'numeric' })
}

onMounted(fetchSharedData)
</script>

<style scoped>
@reference "tailwindcss";

.shared-view-page {
  @apply min-h-screen flex flex-col;
  background-color: var(--bg-secondary);
}

.shared-header {
  @apply px-6 py-4;
  background-color: var(--bg-card);
  border-bottom: 1px solid var(--border-primary);
}

.shared-header-content {
  @apply mx-auto max-w-[800px] flex items-center justify-between w-full;
}

.brand {
  @apply flex items-center gap-4;
}

.brand h1 {
  @apply text-xl font-bold m-0;
  color: var(--accent);
}

.shared-badge {
  @apply text-xs font-medium px-2 py-1 rounded-sm;
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.attribution {
  @apply text-[0.8125rem];
  color: var(--text-secondary);
}

.shared-container {
  @apply flex-1 max-w-[800px] mx-auto w-full p-6;
}

.loading-state, .error-state {
  @apply text-center p-16 rounded-xl shadow-sm;
  background-color: var(--bg-card);
  color: var(--text-secondary);
}

.error-state svg { 
  @apply mb-4; 
  color: var(--error); 
}
.error-state h2 { 
  @apply text-xl font-semibold mb-2; 
  color: var(--text-primary); 
}

.account-summary-card {
  @apply rounded-xl p-8 mb-6 shadow-lg text-white;
  background: linear-gradient(135deg, var(--accent), var(--accent-hover));
}

.account-meta { 
  @apply flex items-center justify-between mb-2; 
}
.account-name { 
  @apply text-xl font-semibold; 
}
.account-currency { 
  @apply text-sm opacity-80 px-2 py-[0.125rem] rounded-sm; 
  background-color: rgba(255,255,255,0.2); 
}
.account-balance { 
  @apply text-[2.25rem] font-bold; 
}
.account-balance.text-danger { 
  color: #ff8080; 
}

.transactions-section h3 { 
  @apply text-lg font-semibold mb-4; 
  color: var(--text-primary); 
}

.transactions-list {
  @apply rounded-xl overflow-hidden shadow-sm;
  background-color: var(--bg-card);
}

.transaction-item {
  @apply flex items-center justify-between px-5 py-4;
  border-bottom: 1px solid var(--border-primary);
}

.transaction-item:last-child { 
  @apply border-b-0; 
}

.tx-left { 
  @apply flex items-center gap-4; 
}

.tx-category-icon { 
  @apply w-[2.5rem] h-[2.5rem] rounded-full flex items-center justify-center shrink-0 text-lg; 
}

.tx-info { 
  @apply flex flex-col gap-0.5; 
}

.tx-payee { 
  @apply font-medium text-[0.9375rem]; 
  color: var(--text-primary); 
}
.tx-note { 
  @apply text-[0.8125rem]; 
  color: var(--text-secondary); 
}
.tx-date { 
  @apply text-[0.75rem]; 
  color: var(--text-muted); 
}

.tx-right { 
  @apply text-right; 
}
.tx-amount { 
  @apply font-semibold text-[0.9375rem]; 
}
.tx-amount.income { 
  color: var(--success); 
}
.tx-amount.expense { 
  color: var(--text-primary); 
}

.empty-transactions {
  @apply text-center p-12 rounded-xl shadow-sm;
  background-color: var(--bg-card);
  color: var(--text-secondary);
}
</style>