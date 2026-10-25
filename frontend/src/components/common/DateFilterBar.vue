<template>
  <div class="date-filter-bar">
    <div class="preset-buttons">
      <button
        v-for="preset in presets"
        :key="preset.value"
        :class="['preset-btn', { active: currentPreset === preset.value }]"
        @click="selectPreset(preset.value)"
      >
        {{ preset.label }}
      </button>
    </div>
    <div class="custom-range">
      <div class="date-input-group">
        <label>From</label>
        <input type="date" v-model="dateFrom" @change="onDateChange" class="date-input" />
      </div>
      <div class="date-input-group">
        <label>To</label>
        <input type="date" v-model="dateTo" @change="onDateChange" class="date-input" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'

interface Props {
  defaultFrom?: string
  defaultTo?: string
}

interface Emits {
  (e: 'change', value: { dateFrom: string; dateTo: string }): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const dateFrom = ref(props.defaultFrom || getDefaultFrom())
const dateTo = ref(props.defaultTo || getDefaultTo())
const currentPreset = ref<string>('')

const presets = [
  { value: 'thisMonth', label: '今月' },
  { value: 'lastMonth', label: '先月' },
  { value: 'last3Months', label: '直近3ヶ月' },
  { value: 'last6Months', label: '直近6ヶ月' },
  { value: 'thisYear', label: '今年' },
  { value: 'lastYear', label: '昨年' },
  { value: 'allTime', label: '全期間' },
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
.date-filter-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  padding: 1rem;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
}

:where([data-theme="dark"], [data-theme="dark"] *) .date-filter-bar {
  background-color: #18181b;
  border-color: #27272a;
}

.preset-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
}

.preset-btn {
  padding: 0.375rem 0.75rem;
  font-size: 0.875rem;
  font-weight: 600;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
  transition: background-color 0.2s ease, color 0.2s ease;
  white-space: nowrap;
}

:where([data-theme="dark"], [data-theme="dark"] *) .preset-btn {
  background-color: #18181b;
  border-color: #27272a;
  color: #f4f4f5;
}

.preset-btn:hover {
  background-color: var(--bg-hover);
}

:where([data-theme="dark"], [data-theme="dark"] *) .preset-btn:hover {
  background-color: #27272a;
}

.preset-btn.active {
  background-color: var(--accent);
  color: white;
  border-color: var(--accent);
}

.custom-range {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.date-input-group {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.date-input-group label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-secondary);
  white-space: nowrap;
}

:where([data-theme="dark"], [data-theme="dark"] *) .date-input-group label {
  color: #d4d4d8;
}

.date-input {
  padding: 0.375rem 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}

:where([data-theme="dark"], [data-theme="dark"] *) .date-input {
  background-color: #121212;
  border-color: #27272a;
  color: #f4f4f5;
}
</style>