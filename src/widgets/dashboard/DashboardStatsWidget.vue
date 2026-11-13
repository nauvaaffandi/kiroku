<template>
  <section class="summary-card" aria-labelledby="summary-heading">
    <h2 id="summary-heading" class="summary-heading">{{ t('dashboard.summary') }}</h2>
    <div class="summary-grid">
      <div class="summary-item">
        <p class="summary-label">{{ t('dashboard.monthlyIncome') }}</p>
        <p class="summary-value income" :class="{ 'is-hidden': isHideAmounts }">
          {{ formatIncome(monthlyIncome) }}
        </p>
      </div>
      <div class="summary-item">
        <p class="summary-label">{{ t('dashboard.monthlyExpense') }}</p>
        <p class="summary-value expense" :class="{ 'is-hidden': isHideAmounts }">
          {{ formatExpense(monthlyExpense) }}
        </p>
      </div>
      <div class="summary-item">
        <p class="summary-label">{{ t('dashboard.net') }}</p>
        <p class="summary-value" :class="[netClass, { 'is-hidden': isHideAmounts }]">
          {{ formatNet(monthlyNet) }}
        </p>
      </div>
      <div class="summary-item">
        <p class="summary-label">{{ t('dashboard.totalBalance') }}</p>
        <p class="summary-value balance" :class="{ 'is-hidden': isHideAmounts }">
          {{ formatBalance(totalBalance) }}
        </p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import { getCurrencySymbol } from '@/utils/amountDisplay'
import { computed } from 'vue'

const { t } = useI18n()
const settingsStore = useSettingsStore()

const isHideAmounts = computed(() => settingsStore.isHideAmounts)
const baseCurrency = computed(() => settingsStore.baseCurrency)
const numberFormat = computed(() => settingsStore.numberFormat)
const currencySymbol = computed(() => getCurrencySymbol(baseCurrency.value))

const props = defineProps<{
  monthlyIncome: number
  monthlyExpense: number
  monthlyNet: number
  totalBalance: number
}>()

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

const netClass = computed(() => {
  if (props.monthlyNet > 0) return 'net-positive'
  if (props.monthlyNet < 0) return 'net-negative'
  return 'net-neutral'
})
</script>

<style scoped>
@reference "tailwindcss";

.summary-card {
  @apply rounded-[14px] border p-8 mb-6;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  color: var(--text-primary);
  border-top-width: 2px;
  border-top-color: var(--accent);
}

.summary-heading {
  @apply text-xs font-medium m-0 block mb-6;
  color: var(--text-muted);
}

.summary-grid {
  @apply grid gap-5;
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

.summary-value.net-positive {
  color: var(--success);
}

.summary-value.net-negative {
  color: var(--error);
}

.summary-value.net-neutral {
  color: var(--text-primary);
}

.summary-value.balance {
  color: var(--text-primary);
}

.summary-value.is-hidden {
  letter-spacing: 0.1em;
}
</style>