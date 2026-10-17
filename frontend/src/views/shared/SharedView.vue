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
.shared-view-page {
  min-height: 100vh;
  background-color: var(--color-neutral-50);
  display: flex;
  flex-direction: column;
}

.shared-header {
  background-color: var(--color-neutral-100);
  border-bottom: 1px solid var(--color-neutral-200);
  padding: 1rem 1.5rem;
}

.shared-header-content {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.brand {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.brand h1 {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-primary);
  margin: 0;
}

.shared-badge {
  font-size: 0.75rem;
  font-weight: 500;
  background-color: var(--color-neutral-200);
  color: var(--color-neutral-700);
  padding: 0.25rem 0.5rem;
  border-radius: var(--radius-sm);
}

.attribution {
  font-size: 0.8125rem;
  color: var(--color-neutral-500);
}

.shared-container {
  flex: 1;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
  padding: 1.5rem;
}

.loading-state, .error-state {
  text-align: center;
  padding: 4rem 2rem;
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  color: var(--color-neutral-500);
}

.error-state svg { color: var(--color-danger); margin-bottom: 1rem; }
.error-state h2 { font-size: 1.25rem; font-weight: 600; color: var(--color-neutral-900); margin-bottom: 0.5rem; }

.account-summary-card {
  background: linear-gradient(135deg, var(--color-primary), var(--color-primary-dark));
  color: white;
  border-radius: var(--radius);
  padding: 2rem;
  box-shadow: var(--shadow);
  margin-bottom: 1.5rem;
}

.account-meta { display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.5rem; }
.account-name { font-size: 1.25rem; font-weight: 600; }
.account-currency { font-size: 0.875rem; opacity: 0.8; background: rgba(255,255,255,0.2); padding: 0.125rem 0.5rem; border-radius: var(--radius-sm); }
.account-balance { font-size: 2.25rem; font-weight: 700; }
.account-balance.text-danger { color: #ff8080; }

.transactions-section h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-neutral-900);
  margin-bottom: 1rem;
}

.transactions-list {
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  overflow: hidden;
}

.transaction-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--color-neutral-200);
}

.transaction-item:last-child { border-bottom: none; }

.tx-left { display: flex; align-items: center; gap: 1rem; }

.tx-category-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.125rem;
  flex-shrink: 0;
}

.tx-info { display: flex; flex-direction: column; gap: 0.125rem; }

.tx-payee { font-weight: 500; color: var(--color-neutral-900); font-size: 0.9375rem; }
.tx-note { font-size: 0.8125rem; color: var(--color-neutral-500); }
.tx-date { font-size: 0.75rem; color: var(--color-neutral-400); }

.tx-right { text-align: right; }
.tx-amount { font-weight: 600; font-size: 0.9375rem; }
.tx-amount.income { color: var(--color-success); }
.tx-amount.expense { color: var(--color-neutral-900); }

.empty-transactions {
  text-align: center;
  padding: 3rem;
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  color: var(--color-neutral-500);
}
</style>