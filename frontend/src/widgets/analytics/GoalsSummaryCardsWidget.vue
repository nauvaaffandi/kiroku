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
.summary-cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(140px, 1fr)); gap: 1rem; }
.summary-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); padding: 1rem; display: flex; flex-direction: column; gap: 0.25rem; }
.summary-label { font-size: 0.75rem; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.05em; }
.summary-value { font-size: 1.25rem; font-weight: 600; color: var(--text-primary); }
.summary-value.positive { color: var(--success); }
.summary-value.negative { color: var(--error); }
</style>