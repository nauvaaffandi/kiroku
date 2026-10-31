<template>
  <section class="accounts-card mb-6" aria-labelledby="accounts-heading">
    <div class="accounts-header">
      <h2 id="accounts-heading" class="accounts-heading">{{ t('dashboard.accounts') }}</h2>
      <router-link to="/accounts" class="view-all-link">{{ t('dashboard.viewAll') }}</router-link>
    </div>

    <div v-if="accounts.length === 0" class="accounts-empty">
      <p class="accounts-empty-text">{{ t('dashboard.noAccounts') }}</p>
      <router-link to="/accounts" class="accounts-empty-link">{{ t('accounts.addAccount') }}</router-link>
    </div>

    <div v-else class="accounts-list">
      <div
        v-for="account in accounts"
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
          <span class="account-amount" :class="{ 'is-hidden': isHideAmounts }">
            {{ formatAccountBalance(account.balance, account.currency) }}
          </span>
          <span class="account-currency">{{ account.currency }}</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import { computed } from 'vue'
import { getCurrencySymbol } from '@/utils/amountDisplay'
import type { Account } from '@/types'

const { t } = useI18n()
const settingsStore = useSettingsStore()

const isHideAmounts = computed(() => settingsStore.isHideAmounts)
const numberFormat = computed(() => settingsStore.numberFormat)

defineProps<{
  accounts: Account[]
}>()

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
</script>

<style scoped>
@reference "tailwindcss";

.accounts-card {
  @apply rounded-[14px] border p-4;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  color: var(--text-primary);
}

.accounts-header {
  @apply flex items-center justify-between mb-4;
}

.accounts-heading {
  @apply text-sm font-semibold m-0;
  color: var(--text-muted);
}

.view-all-link {
  @apply text-xs font-medium no-underline;
  color: var(--accent);
}

.view-all-link:hover {
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

.account-amount.is-hidden {
  letter-spacing: 0.1em;
}

.account-currency {
  @apply text-xs;
  color: var(--text-muted);
}
</style>