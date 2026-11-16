<template>
  <div class="stat">
    <div class="stat-label">{{ label }}</div>
    <div class="stat-value" :aria-label="ariaLabel">{{ displayValue }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatAmountForDisplay } from '@/utils/amountDisplay'
import { useSettingsStore } from '@/stores/settings'
import { useI18n } from '@/locales'

interface Props {
  label: string
  value: number
  type?: 'income' | 'expense' | 'transfer'
  hideAmounts?: boolean
}

const props = defineProps<Props>()
const settingsStore = useSettingsStore()
const { t } = useI18n()

const displayValue = computed(() => {
  if (props.hideAmounts || settingsStore.isHideAmounts) {
    return '••••'
  }
  return formatAmountForDisplay(props.value, props.type || 'transfer')
})

const ariaLabel = computed(() => {
  const v = displayValue.value
  return v === '••••' ? `${props.label} ${t('stat.amountHidden')}` : `${props.label} ${v}`
})
</script>

<style scoped>
@reference "tailwindcss";

.stat {
  @apply flex flex-col items-start gap-1;
}

.stat-label {
  @apply text-[0.75rem] font-medium text-[var(--text-muted)];
}

.stat-value {
  @apply text-[24px] font-semibold leading-[1.2] tabular-nums text-[var(--text-primary)];
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
