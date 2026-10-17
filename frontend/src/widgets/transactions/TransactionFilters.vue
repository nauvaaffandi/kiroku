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
.filters-bar {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
  background: var(--color-surface);
  padding: 1rem;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

.filter-group {
  display: flex;
  flex-direction: column;
}
</style>