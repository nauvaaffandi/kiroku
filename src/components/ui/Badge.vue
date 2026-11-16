<template>
  <span class="badge" :class="`badge-${variant}`" :aria-label="ariaLabel">{{ label }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  label: string
  variant?: 'income' | 'expense' | 'warning' | 'neutral'
}

const props = defineProps<Props>()

const variant = computed(() => props.variant || 'neutral')

const ariaLabel = computed(() => `${props.label} (${variant.value})`)
</script>

<style scoped>
@reference "tailwindcss";

.badge {
  @apply inline-block px-2 py-0.5 text-[0.625rem] font-medium rounded-[6px] whitespace-nowrap;
  background-color: var(--bg-tertiary);
  color: var(--text-secondary);
}

.badge.income {
  @apply bg-[var(--success-light)] text-[var(--success)];
}

.badge.expense {
  @apply bg-[var(--error-light)] text-[var(--error)];
}

.badge.warning {
  @apply bg-[var(--warning-light)] text-[var(--warning)];
}
</style>
