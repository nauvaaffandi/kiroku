<template>
  <button
    type="button"
    class="color-swatch"
    :style="{ backgroundColor: color }"
    :class="{ active: colorActive }"
    @click="$emit('update:color', color)"
    :aria-pressed="colorActive"
    :aria-label="t('colorPicker.select', { color })"
    :disabled="disabled"
  />
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'

interface Props {
  color: string
  colorActive: boolean
  disabled?: boolean
}

defineProps<Props>()
defineEmits<{
  'update:color': [value: string]
}>()

const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.color-swatch {
  @apply flex shrink-0 cursor-pointer;
  width: 100%;
  height: 3rem;
  border-radius: var(--radius-sm);
  border: 2px solid transparent;
}

.color-swatch.active {
  border-color: var(--accent);
}

.color-swatch:focus-visible {
  outline: 2px solid var(--border-focus);
  outline-offset: -2px;
}

.color-swatch:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
</style>