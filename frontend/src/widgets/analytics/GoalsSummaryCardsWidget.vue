<template>
  <div class="summary-cards" v-if="!loading">
    <div class="summary-card">
      <span class="summary-label">{{ t('analytics.goalsSummary.totalTarget') }}</span>
      <span class="summary-value">{{ formatCurrency(displayGoalsTotal) }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('analytics.goalsSummary.currentAmount') }}</span>
      <span class="summary-value">{{ formatCurrency(displayGoalsCurrent) }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('analytics.goalsSummary.overallProgress') }}</span>
      <span class="summary-value">{{ displayOverallProgress.toFixed(1) }}%</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('analytics.goalsSummary.achieved') }}</span>
      <span class="summary-value positive">{{ displayAchieved }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('analytics.goalsSummary.onTrack') }}</span>
      <span class="summary-value positive">{{ displayOnTrack }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('analytics.goalsSummary.behind') }}</span>
      <span class="summary-value negative">{{ displayBehind }}</span>
    </div>
  </div>
  <div class="summary-cards" v-else>
    <div class="summary-card" v-for="i in 6" :key="i">
      <div class="skeleton skeleton-text" style="width: 50%;" />
      <div class="skeleton skeleton-text" style="width: 80%; height: 1.5rem;" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { computed } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useExchangeRate } from '@/composables/useExchangeRate'
import { getCurrencySymbol } from '@/utils/amountDisplay'

interface Summary {
  total_targets: number
  total_current: number
  overall_progress: number
  goals_achieved: number
  goals_on_track: number
  goals_behind: number
}

interface GoalItem {
  target_amount: number
  current_amount: number
  status: string
}

const { t } = useI18n()
const settingsStore = useSettingsStore()
const { toUsd } = useExchangeRate()
const baseCurrency = computed(() => settingsStore.baseCurrency || 'USD')

const props = defineProps<{
  loading: boolean
  summary: Summary
  goals: GoalItem[]
}>()

const displayGoalsTotal = computed(() =>
  props.goals.reduce((sum: number, g: GoalItem) => sum + g.target_amount, 0)
)
const displayGoalsCurrent = computed(() =>
  props.goals.reduce((sum: number, g: GoalItem) => sum + g.current_amount, 0)
)
const displayOverallProgress = computed(() => {
  const t = displayGoalsTotal.value
  return t > 0 ? (displayGoalsCurrent.value / t) * 100 : 0
})
const displayAchieved = computed(() =>
  props.goals.filter((g) => g.status === 'achieved').length
)
const displayOnTrack = computed(() =>
  props.goals.filter((g) => g.status === 'on_track').length
)
const displayBehind = computed(() =>
  props.goals.filter((g) => g.status === 'behind').length
)

function formatCurrency(val: number): string {
  const usd = toUsd(val)
  const symbol = getCurrencySymbol(baseCurrency.value)
  const locale = settingsStore.numberFormat || 'en-US'
  return `${symbol}${new Intl.NumberFormat(locale).format(Math.round(usd))}`
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
