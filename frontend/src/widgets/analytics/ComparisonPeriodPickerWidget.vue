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

.comparison-periods {
  @apply flex flex-wrap gap-6 p-4 border rounded-xl;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}
.period-picker {
  @apply flex flex-col gap-2 min-w-[250px];
}
.period-picker label {
  @apply text-sm font-medium;
  color: var(--text-secondary);
}
.period-inputs {
  @apply flex items-center gap-2;
}
.period-inputs .date-input {
  @apply flex-1;
}
</style>
