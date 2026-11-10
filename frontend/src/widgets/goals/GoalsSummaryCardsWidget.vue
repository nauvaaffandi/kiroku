<template>
  <div class="summary-cards">
    <div class="summary-card">
      <span class="summary-label">{{ t('goals.totalTarget') }}</span>
      <span class="summary-value">{{ formatCurrencyShort(totalTarget) }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('goals.currentAmount') }}</span>
      <span class="summary-value">{{ formatCurrencyShort(currentAmount) }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('goals.overallProgress') }}</span>
      <span class="summary-value">{{ overallProgress.toFixed(1) }}%</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('goals.achieved') }}</span>
      <span class="summary-value positive">{{ achieved }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('goals.onTrack') }}</span>
      <span class="summary-value positive">{{ onTrack }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('goals.behind') }}</span>
      <span class="summary-value negative">{{ behind }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { computed } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useExchangeRate } from '@/composables/useExchangeRate'
import { getCurrencySymbol } from '@/utils/amountDisplay'
import type { Goal } from '@/types/goal'

interface GoalItem extends Goal {
  status: 'on_track' | 'behind' | 'achieved'
  progress_pct: number
  remaining: number
}

const { t } = useI18n()
const settingsStore = useSettingsStore()
const { toUsd } = useExchangeRate()
const baseCurrency = computed(() => settingsStore.baseCurrency || 'USD')

const props = defineProps<{
  goals: GoalItem[]
}>()

const totalTarget = computed(() =>
  props.goals.reduce((sum: number, g: GoalItem) => sum + g.target_amount, 0)
)
const currentAmount = computed(() =>
  props.goals.reduce((sum: number, g: GoalItem) => sum + g.saved_already, 0)
)
const overallProgress = computed(() => {
  const t = totalTarget.value
  return t > 0 ? (currentAmount.value / t) * 100 : 0
})
const achieved = computed(() =>
  props.goals.filter((g) => g.status === 'achieved').length
)
const onTrack = computed(() =>
  props.goals.filter((g) => g.status === 'on_track').length
)
const behind = computed(() =>
  props.goals.filter((g) => g.status === 'behind').length
)

function formatCurrencyShort(val: number): string {
  const usd = toUsd(val)
  const symbol = getCurrencySymbol(baseCurrency.value)
  if (Math.abs(usd) >= 1e8) return `${symbol}${(usd / 1e8).toFixed(1)}億`
  if (Math.abs(usd) >= 1e4) return `${symbol}${(usd / 1e4).toFixed(1)}万`
  const locale = settingsStore.numberFormat || 'en-US'
  return `${symbol}${usd.toLocaleString(locale, { maximumFractionDigits: 0 })}`
}
</script>

<style scoped>
@reference "tailwindcss";

.summary-cards {
  @apply grid gap-4;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
}
.summary-card {
  @apply border rounded-xl p-4 flex flex-col gap-1;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}
.summary-label {
  @apply text-xs uppercase tracking-wider;
  color: var(--text-secondary);
}
.summary-value {
  @apply text-xl font-semibold;
  color: var(--text-primary);
}
.summary-value.positive {
  color: var(--success);
}
.summary-value.negative {
  color: var(--error);
}
</style>