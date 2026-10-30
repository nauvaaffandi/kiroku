<template>
  <div class="filters-bar">
    <div class="filter-group">
      <BaseInput
        v-model="localFilters.q"
        type="search"
        :placeholder="t('common.search')"
        @input="emitSearch"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.archived"
        :options="archivedOptions"
        :placeholder="t('common.status')"
        @change="emitChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { BaseInput, BaseSelect } from '@/components/ui'
import { useI18n } from '@/locales'
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

const { t } = useI18n()

const localFilters = ref<CategoryFilters>({ ...props.modelValue })

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal }
}, { deep: true })

const archivedOptions = [
  { value: '', label: t('common.all') },
  { value: 'false', label: t('categories.filterActive') || 'Active' },
  { value: 'true', label: t('categories.filterArchived') || 'Archived' }
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
