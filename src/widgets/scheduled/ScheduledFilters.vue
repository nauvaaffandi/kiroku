<template>
  <div class="filters-bar">
    <div class="filter-group">
      <BaseInput
        v-model="localFilters.q"
        type="search"
        :placeholder="t('scheduled.searchPlaceholder')"
        @input="emitSearch"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.type"
        :options="typeOptions"
        :placeholder="t('scheduled.allTypes')"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.frequency"
        :options="frequencyOptions"
        :placeholder="t('scheduled.allFrequencies')"
        @change="emitChange"
      />
    </div>
    <div class="filter-group">
      <BaseSelect
        v-model="localFilters.active"
        :options="activeOptions"
        :placeholder="t('scheduled.allStatus')"
        @change="emitChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseInput, BaseSelect } from '@/components/ui'
import { useI18n } from '@/locales'
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

const { t } = useI18n()

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal }
}, { deep: true })

const typeOptions = computed(() => [
  { value: 'income', label: t('transactions.typeIncome') },
  { value: 'expense', label: t('transactions.typeExpense') },
  { value: 'transfer', label: t('transactions.typeTransfer') }
])

const frequencyOptions = computed(() => [
  { value: 'daily', label: t('scheduled.freqDaily') },
  { value: 'weekly', label: t('scheduled.freqWeekly') },
  { value: 'monthly', label: t('scheduled.freqMonthly') },
  { value: 'yearly', label: t('scheduled.freqYearly') },
  { value: 'custom', label: t('scheduled.freqCustom') }
])

const activeOptions = computed(() => [
  { value: '', label: t('scheduled.filterAll') },
  { value: 'true', label: t('scheduled.filterActive') },
  { value: 'false', label: t('scheduled.filterInactive') }
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
  border-radius: var(--radius-lg);
  @apply grid gap-4 mb-6 p-4;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.filter-group {
  @apply flex flex-col;
}
</style>