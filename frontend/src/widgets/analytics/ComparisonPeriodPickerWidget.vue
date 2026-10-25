<template>
  <div class="comparison-periods" v-if="showComparisonPeriods">
    <div class="period-picker">
      <label>期間 A (比較元)</label>
      <div class="period-inputs">
        <input type="date" :value="periodAFrom" @input="$emit('update:periodAFrom', ($event.target as HTMLInputElement).value)" class="date-input" />
        <span>〜</span>
        <input type="date" :value="periodATo" @input="$emit('update:periodATo', ($event.target as HTMLInputElement).value)" class="date-input" />
      </div>
    </div>
    <div class="period-picker">
      <label>期間 B (比較先)</label>
      <div class="period-inputs">
        <input type="date" :value="periodBFrom" @input="$emit('update:periodBFrom', ($event.target as HTMLInputElement).value)" class="date-input" />
        <span>〜</span>
        <input type="date" :value="periodBTo" @input="$emit('update:periodBTo', ($event.target as HTMLInputElement).value)" class="date-input" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  showComparisonPeriods: boolean
  periodAFrom: string
  periodATo: string
  periodBFrom: string
  periodBTo: string
}>()

defineEmits<{
  'update:periodAFrom': [val: string]
  'update:periodATo': [val: string]
  'update:periodBFrom': [val: string]
  'update:periodBTo': [val: string]
}>()
</script>

<style scoped>
@reference "tailwindcss";
.comparison-periods { display: flex; flex-wrap: wrap; gap: 1.5rem; padding: 1rem; background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); }
.period-picker { display: flex; flex-direction: column; gap: 0.5rem; min-width: 250px; }
.period-picker label { font-size: 0.875rem; font-weight: 500; color: var(--text-secondary); }
.period-inputs { display: flex; align-items: center; gap: 0.5rem; }
.period-inputs .date-input { flex: 1; }
</style>