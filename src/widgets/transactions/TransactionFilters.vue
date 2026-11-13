<template>
  <div class="filters-bar">
    <div class="filter-group">
      <BaseInput
        v-model="localSearchQuery"
        type="search"
        placeholder="Search transactions..."
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
        v-model="localFilters.category_id"
        :options="categoryOptions"
        placeholder="All Categories"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseInput
        v-model="localFilters.date_from"
        type="date"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseInput
        v-model="localFilters.date_to"
        type="date"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.status"
        :options="statusOptions"
        placeholder="All Status"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.account_id"
        :options="accountOptions"
        placeholder="All Accounts"
        @change="emitChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { BaseInput, BaseSelect } from '@/components/ui'
import type { TransactionFilters, Category, Account } from '@/types'

interface Props {
  modelValue: TransactionFilters
  categories: Category[]
  accounts: Account[]
  searchQuery: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: TransactionFilters]
  change: []
  search: [value: string]
}>()

const localFilters = ref<TransactionFilters>({ ...props.modelValue })
const localSearchQuery = ref(props.searchQuery)

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal }
}, { deep: true })

watch(() => props.searchQuery, (newVal) => {
  localSearchQuery.value = newVal
})

watch(localSearchQuery, (value) => {
  emit('search', value)
})

const typeOptions = [
  { value: 'income', label: 'Income' },
  { value: 'expense', label: 'Expense' },
  { value: 'transfer', label: 'Transfer' }
]

const statusOptions = [
  { value: 'cleared', label: 'Cleared' },
  { value: 'pending', label: 'Pending' }
]

const categoryOptions = props.categories.map(c => ({
  value: c.id,
  label: c.path || c.name,
  disabled: c.children && c.children.length > 0
}))

const accountOptions = props.accounts.map(a => ({
  value: a.id,
  label: a.name
}))

const emitChange = () => {
  emit('update:modelValue', { ...localFilters.value })
  emit('change')
}
</script>

<style scoped>
@reference "tailwindcss";
.filters-bar {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  @apply flex flex-col gap-3 p-3 mb-4;
}

.filter-group {
  @apply flex flex-col w-full;
}

@media (min-width: 640px) {
  .filters-bar {
    @apply grid p-4 gap-4 mb-6;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  }
}
</style>