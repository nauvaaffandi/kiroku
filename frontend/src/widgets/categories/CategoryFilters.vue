<template>
  <div class="filters-bar">
    <div class="filter-group">
      <BaseInput
        v-model="localFilters.q"
        type="search"
        placeholder="Search categories..."
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
        v-model="localFilters.archived"
        :options="archivedOptions"
        placeholder="All Status"
        @change="emitChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { BaseInput, BaseSelect } from '@/components/ui'
import type { CategoryFilters } from '@/types'

interface Props {
  modelValue: CategoryFilters
  searchQuery: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: CategoryFilters]
  change: []
  search: [value: string]
}>()

const localFilters = ref<CategoryFilters>({ ...props.modelValue })

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal }
}, { deep: true })

const typeOptions = [
  { value: 'income', label: 'Income' },
  { value: 'expense', label: 'Expense' },
  { value: 'transfer', label: 'Transfer' }
]

const archivedOptions = [
  { value: '', label: 'All' },
  { value: 'false', label: 'Active' },
  { value: 'true', label: 'Archived' }
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
  @apply grid gap-4 mb-6 p-4 rounded-2xl;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}
.filter-group { @apply flex flex-col; }
</style>