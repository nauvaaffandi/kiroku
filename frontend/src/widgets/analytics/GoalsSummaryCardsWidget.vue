<template>
  <div class="summary-cards" v-if="!loading">
    <div class="summary-card">
      <span class="summary-label">総目標額</span>
      <span class="summary-value">{{ formatCurrency(summary.total_targets) }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">現在額</span>
      <span class="summary-value">{{ formatCurrency(summary.total_current) }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">全体進捗</span>
      <span class="summary-value">{{ summary.overall_progress.toFixed(1) }}%</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">達成</span>
      <span class="summary-value positive">{{ summary.goals_achieved }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">順調</span>
      <span class="summary-value positive">{{ summary.goals_on_track }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">遅れ</span>
      <span class="summary-value negative">{{ summary.goals_behind }}</span>
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
interface Summary {
  total_targets: number
  total_current: number
  overall_progress: number
  goals_achieved: number
  goals_on_track: number
  goals_behind: number
}

defineProps<{
  loading: boolean
  summary: Summary
}>()

function formatCurrency(val: number): string {
  return new Intl.NumberFormat('ja-JP', { 
    style: 'currency', 
    currency: 'JPY',
    maximumFractionDigits: 0 
  }).format(val)
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
