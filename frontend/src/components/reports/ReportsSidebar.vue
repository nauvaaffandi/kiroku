<template>
  <aside class="reports-sidebar">
    <section class="sidebar-section">
      <h3 class="sidebar-label">{{ t('reports.spendingReport') }}</h3>
      <div class="sidebar-group">
        <div>
          <label class="field-label">{{ t('reports.period') }}</label>
          <select v-model="spendingParams.period" @change="$emit('spending')" class="field-select">
            <option value="this_week">{{ t('reports.thisWeek') }}</option>
            <option value="this_month">{{ t('reports.thisMonth') }}</option>
            <option value="this_year">{{ t('reports.thisYear') }}</option>
            <option value="last_month">{{ t('reports.lastMonth') }}</option>
            <option value="last_year">{{ t('reports.lastYear') }}</option>
          </select>
        </div>
        <div>
          <label class="field-label">{{ t('reports.groupBy') }}</label>
          <select v-model="spendingParams.group_by" @change="$emit('spending')" class="field-select">
            <option value="category">{{ t('reports.category') }}</option>
            <option value="parent_category">{{ t('reports.parentCategory') }}</option>
          </select>
        </div>
        <div>
          <label class="field-label">{{ t('reports.accountOptional') }}</label>
          <select v-model="spendingParams.account_id" @change="$emit('spending')" class="field-select">
            <option value="">{{ t('reports.allAccounts') }}</option>
            <option v-for="acc in accounts" :key="acc.id" :value="acc.id">{{ acc.name }}</option>
          </select>
        </div>
      </div>
    </section>
    <section class="sidebar-section">
      <h3 class="sidebar-label">{{ t('reports.trendReport') }}</h3>
      <div>
        <label class="field-label">{{ t('reports.months') }}</label>
        <select v-model="trendParams.months" @change="$emit('trend')" class="field-select">
          <option value="6">{{ 6 }} {{ t('reports.months') }}</option>
          <option value="12" selected>{{ 12 }} {{ t('reports.months') }}</option>
          <option value="24">{{ 24 }} {{ t('reports.months') }}</option>
        </select>
      </div>
    </section>
  </aside>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'

defineProps<{
  spendingParams: { period: string; group_by: string; account_id: string; date_from: string; date_to: string }
  trendParams: { months: string; account_id: string; category_id: string }
  accounts: { id: string; name: string }[]
}>()
defineEmits<{ spending: []; trend: [] }>()
const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.reports-sidebar {
  @apply bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl p-4 shadow-sm h-fit;
}

.sidebar-section { @apply mb-6; }

.sidebar-label {
  @apply text-xs font-bold text-zinc-500 dark:text-zinc-400 uppercase tracking-wider mb-3;
}

.sidebar-group { @apply space-y-3; }

.field-label { @apply block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1; }

.field-select {
  @apply w-full sm:w-auto px-3 py-2 bg-white dark:bg-zinc-800 border border-zinc-300 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:outline-none focus:ring-2 focus:ring-sky-500;
}
</style>
