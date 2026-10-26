<template>
  <div class="summary-cards" v-if="!loading">
    <div class="summary-card">
      <span class="summary-label">総予算</span>
      <span class="summary-value">{{ formatCurrency(summary.total_budget) }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">総実績</span>
      <span class="summary-value">{{ formatCurrency(summary.total_actual) }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">差額</span>
      <span class="summary-value" :class="summary.total_variance >= 0 ? 'positive' : 'negative'">
        {{ summary.total_variance >= 0 ? '+' : '' }}{{ formatCurrency(summary.total_variance) }}
      </span>
    </div>
    <div class="summary-card">
      <span class="summary-label">超過カテゴリ</span>
      <span class="summary-value negative">{{ summary.categories_over_budget }}</span>
    </div>
    <div class="summary-card">
      <span class="summary-label">未達カテゴリ</span>
      <span class="summary-value positive">{{ summary.categories_under_budget }}</span>
    </div>
  </div>
  <div class="summary-cards" v-else>
    <div class="summary-card" v-for="i in 5" :key="i">
      <div class="skeleton skeleton-text" style="width: 50%;" />
      <div class="skeleton skeleton-text" style="width: 80%; height: 1.5rem;" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Summary {
  total_budget: number
  total_actual: number
  total_variance: number
  categories_over_budget: number
  categories_under_budget: number
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
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
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
