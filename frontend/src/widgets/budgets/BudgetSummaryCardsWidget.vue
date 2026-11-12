<template>
  <div class="summary-cards" v-if="!loading">
    <div class="summary-card">
      <span class="summary-label">{{ t('budgets.summary.totalBudgets') }}</span>
      <span class="summary-value">{{ summary.totalBudgets }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('budgets.summary.overBudgets') }}</span>
      <span class="summary-value negative">{{ summary.overBudgets }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">{{ t('budgets.summary.underBudgets') }}</span>
      <span class="summary-value positive">{{ summary.underBudgets }}</span>
    </div>
  </div>
  <div class="summary-cards" v-else>
    <div class="summary-card" v-for="i in 3" :key="i">
      <div class="skeleton skeleton-text skeleton-total" />
      <div class="skeleton skeleton-text skeleton-label" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'

interface Summary {
  totalBudgets: number
  overBudgets: number
  underBudgets: number
}

const { t } = useI18n()

defineProps<{
  loading: boolean
  summary: Summary
}>()
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

.skeleton-total {
  @apply w-[50%];
}
.skeleton-label {
  @apply w-[80%] h-[1.5rem];
}
</style>