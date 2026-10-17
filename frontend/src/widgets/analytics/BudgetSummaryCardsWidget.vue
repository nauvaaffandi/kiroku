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
.summary-cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 1rem; }
.summary-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); padding: 1rem; display: flex; flex-direction: column; gap: 0.25rem; }
.summary-label { font-size: 0.75rem; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.05em; }
.summary-value { font-size: 1.25rem; font-weight: 600; color: var(--text-primary); }
.summary-value.positive { color: var(--success); }
.summary-value.negative { color: var(--error); }
</style>