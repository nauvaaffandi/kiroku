<template>
  <div class="filters-bar">
    <div class="filter-group">
      <BaseInput
        v-model="localFilters.q"
        type="search"
        placeholder="Search scheduled transactions..."
        @input="emitSearch"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.type"
        :options="typeOptions"
        placeholder="All Types"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.frequency"
        :options="frequencyOptions"
        placeholder="All Frequencies"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.active"
        :options="activeOptions"
        placeholder="All Status"
        @change="emitChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { BaseInput, BaseSelect } from '@/components/ui'
import type { ScheduledFilters } from '@/types'

interface Props {
  modelValue: ScheduledFilters
  searchQuery: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: ScheduledFilters]
  change: []
  search: [value: string]
}>()

const localFilters = ref<ScheduledFilters>({ ...props.modelValue })

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal }
}, { deep: true })

const typeOptions = [
  { value: 'income', label: 'Income' },
  { value: 'expense', label: 'Expense' },
  { value: 'transfer', label: 'Transfer' }
]

const frequencyOptions = [
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: 'Monthly' },
  { value: 'yearly', label: 'Yearly' },
  { value: 'custom', label: 'Custom' }
]

const activeOptions = [
  { value: '', label: 'All' },
  { value: 'true', label: 'Active' },
  { value: 'false', label: 'Inactive' }
]

const emitChange = () => {
  emit('update:modelValue', { ...localFilters.value })
  emit('change')
}

const emitSearch = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('search', target.value)
}
</script>

<style scoped>
@reference "tailwindcss";
.filters-bar {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  @apply grid gap-4 mb-6 p-4;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.filter-group {
  @apply flex flex-col;
}
</style>