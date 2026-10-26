<template>
  <div class="date-range-picker">
    <div class="picker-group">
      <label class="picker-label">From</label>
      <BaseInput
        v-model="localRange.from"
        type="date"
        @change="emitChange"
      />
    </div>
    <div class="picker-group">
      <label class="picker-label">To</label>
      <BaseInput
        v-model="localRange.to"
        type="date"
        @change="emitChange"
      />
    </div>
    <div class="picker-presets">
      <BaseButton
        v-for="preset in presets"
        :key="preset.value"
        :variant="localRange.preset === preset.value ? 'primary' : 'secondary'"
        size="sm"
        @click="selectPreset(preset.value)"
      >
        {{ preset.label }}
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { BaseInput, BaseButton } from '@/components/ui'
import type { DateRange } from '@/types'

interface Props {
  modelValue: DateRange
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: DateRange]
  change: []
}>()

const localRange = ref<DateRange>({ ...props.modelValue })

watch(() => props.modelValue, (newVal) => {
  localRange.value = { ...newVal }
}, { deep: true })

const presets = [
  { value: 'today', label: 'Today' },
  { value: 'week', label: 'This Week' },
  { value: 'month', label: 'This Month' },
  { value: 'quarter', label: 'This Quarter' },
  { value: 'year', label: 'This Year' },
  { value: 'custom', label: 'Custom' }
] as const

const emitChange = () => {
  emit('update:modelValue', { ...localRange.value })
  emit('change')
}

const selectPreset = (preset: typeof presets[number]['value']) => {
  localRange.value.preset = preset
  const now = new Date()
  
  switch (preset) {
    case 'today':
      localRange.value.from = formatDate(now)
      localRange.value.to = formatDate(now)
      break
    case 'week':
      const startOfWeek = new Date(now)
      startOfWeek.setDate(now.getDate() - now.getDay())
      localRange.value.from = formatDate(startOfWeek)
      localRange.value.to = formatDate(now)
      break
    case 'month':
      localRange.value.from = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`
      localRange.value.to = formatDate(now)
      break
    case 'quarter':
      const quarterStart = new Date(now.getFullYear(), Math.floor(now.getMonth() / 3) * 3, 1)
      localRange.value.from = formatDate(quarterStart)
      localRange.value.to = formatDate(now)
      break
    case 'year':
      localRange.value.from = `${now.getFullYear()}-01-01`
      localRange.value.to = formatDate(now)
      break
    case 'custom':
      // Keep current custom dates
      break
  }
  emitChange()
}

const formatDate = (date: Date) => {
  return date.toISOString().split('T')[0]
}
</script>

<style scoped>
@reference "tailwindcss";

.date-range-picker {
  @apply flex flex-col gap-4;
}

.picker-group {
  @apply flex flex-col gap-[0.375rem];
}

.picker-label {
  @apply text-[0.8125rem] font-medium;
  color: var(--text-primary);
}

.picker-presets {
  @apply flex flex-wrap gap-2;
}
</style>
