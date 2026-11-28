<template>
  <div class="filters-bar">
    <div class="filter-group">
      <BaseInput
        v-model="localFilters.q"
        type="search"
        :placeholder="t('accounts.searchPlaceholder')"
        @input="emitSearch"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.type"
        :options="typeOptions"
        :placeholder="t('accounts.allTypes')"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.currency"
        :options="currencyOptions"
        :placeholder="t('accounts.allCurrencies')"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.archived"
        :options="archivedOptions"
        :placeholder="t('accounts.allStatus')"
        @change="emitChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseInput, BaseSelect } from '@/components/ui'
import { useI18n } from '@/locales'
import type { AccountFilters } from '@/types'
import { useCurrencies } from '@/composables/useCurrencies'

const { currencies } = useCurrencies()
const { t } = useI18n()

const currencyOptions = computed(() => currencies.value.map(c => ({ value: c, label: c })))

interface Props {
  modelValue: AccountFilters
  searchQuery: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: AccountFilters]
  change: []
  search: [value: string]
}>()

const localFilters = ref<AccountFilters>({ ...props.modelValue })

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal }
}, { deep: true })

const typeOptions = computed(() => [
  { value: 'cash', label: t('accounts.typeCash') },
  { value: 'bank', label: t('accounts.typeBank') },
  { value: 'credit_card', label: t('accounts.typeCreditCard') },
  { value: 'investment', label: t('accounts.typeInvestment') },
  { value: 'ewallet', label: t('accounts.typeEWallet') },
  { value: 'other', label: t('accounts.typeOther') }
])

const archivedOptions = computed(() => [
  { value: '', label: t('accounts.filterAll') },
  { value: 'false', label: t('categories.filterActive') },
  { value: 'true', label: t('accounts.archived') }
])

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

.filter-group {
  @apply flex flex-col;
}
</style>