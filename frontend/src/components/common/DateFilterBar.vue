<template>
  <div class="date-filter-bar">
    <div class="preset-buttons">
      <button
        v-for="preset in presets"
        :key="preset.value"
        :class="['preset-btn', { active: currentPreset === preset.value }]"
        @click="selectPreset(preset.value)"
      >
        {{ preset.label() }}
      </button>
    </div>
    <div class="custom-range">
      <div class="date-input-group">
        <label>{{ t('common.from') }}</label>
        <input type="date" v-model="dateFrom" @change="onDateChange" class="date-input" />
      </div>
      <div class="date-input-group">
        <label>{{ t('common.to') }}</label>
        <input type="date" v-model="dateTo" @change="onDateChange" class="date-input" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useI18n } from '@/locales'

interface Props {
  defaultFrom?: string
  defaultTo?: string
}

interface Emits {
  (e: 'change', value: { dateFrom: string; dateTo: string }): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const { t } = useI18n()

const dateFrom = ref(props.defaultFrom || getDefaultFrom())
const dateTo = ref(props.defaultTo || getDefaultTo())
const currentPreset = ref<string>('')

const presets = [
  { value: 'thisMonth', label: () => t('common.thisMonth') },
  { value: 'lastMonth', label: () => t('common.lastMonth') },
  { value: 'last3Months', label: () => t('common.last3Months') },
  { value: 'last6Months', label: () => t('common.last6Months') },
  { value: 'thisYear', label: () => t('common.thisYear') },
  { value: 'lastYear', label: () => t('common.lastYear') },
  { value: 'allTime', label: () => t('common.allTime') },
]

function getDefaultFrom(): string {
  const d = new Date()
  d.setMonth(d.getMonth() - 1)
  return d.toISOString().split('T')[0]
}

function getDefaultTo(): string {
  return new Date().toISOString().split('T')[0]
}

function selectPreset(value: string) {
  currentPreset.value = value
  const { from, to } = getPresetRange(value)
  dateFrom.value = from
  dateTo.value = to
  emit('change', { dateFrom: from, dateTo: to })
}

function getPresetRange(preset: string): { from: string; to: string } {
  const now = new Date()
  const today = now.toISOString().split('T')[0]
  let from: string

  switch (preset) {
    case 'thisMonth':
      from = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0]
      break
    case 'lastMonth':
      from = new Date(now.getFullYear(), now.getMonth() - 1, 1).toISOString().split('T')[0]
      return { from, to: new Date(now.getFullYear(), now.getMonth(), 0).toISOString().split('T')[0] }
    case 'last3Months':
      from = new Date(now.getFullYear(), now.getMonth() - 2, 1).toISOString().split('T')[0]
      break
    case 'last6Months':
      from = new Date(now.getFullYear(), now.getMonth() - 5, 1).toISOString().split('T')[0]
      break
    case 'thisYear':
      from = new Date(now.getFullYear(), 0, 1).toISOString().split('T')[0]
      break
    case 'lastYear':
      from = new Date(now.getFullYear() - 1, 0, 1).toISOString().split('T')[0]
      return { from, to: new Date(now.getFullYear() - 1, 11, 31).toISOString().split('T')[0] }
    case 'allTime':
      from = '2020-01-01'
      break
    default:
      from = getDefaultFrom()
  }
  return { from, to: today }
}

function onDateChange() {
  currentPreset.value = ''
  emit('change', { dateFrom: dateFrom.value, dateTo: dateTo.value })
}

watch([dateFrom, dateTo], () => {
  if (!currentPreset.value) {
    emit('change', { dateFrom: dateFrom.value, dateTo: dateTo.value })
  }
})

onMounted(() => {
  for (const p of presets) {
    const { from, to } = getPresetRange(p.value)
    if (from === dateFrom.value && to === dateTo.value) {
      currentPreset.value = p.value
      break
    }
  }
})
</script>

<style scoped>
@reference "tailwindcss";

.date-filter-bar {
  @apply flex flex-wrap items-center justify-between gap-3 p-4 rounded;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
}

.preset-buttons {
  @apply flex flex-wrap gap-1.5;
}

.preset-btn {
  @apply px-3 py-1.5 text-sm font-semibold rounded border whitespace-nowrap transition-colors;
  background-color: var(--bg-card);
  color: var(--text-primary);
  border-color: var(--border-primary);
}

.preset-btn:hover {
  background-color: var(--bg-hover);
}

.preset-btn.active {
  background-color: var(--accent);
  color: var(--text-inverse);
  border-color: var(--accent);
}

.custom-range {
  @apply flex items-center gap-3 flex-wrap;
}

.date-input-group {
  @apply flex items-center gap-1.5;
}

.date-input-group label {
  @apply text-sm font-semibold whitespace-nowrap;
  color: var(--text-secondary);
}

.date-input {
  @apply px-3 py-1.5 text-sm font-medium rounded border;
  background-color: var(--bg-card);
  color: var(--text-primary);
  border-color: var(--border-primary);
}
</style>