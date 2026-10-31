<template>
  <div class="calendar-header">
    <div class="nav-controls">
      <button class="nav-btn" @click="$emit('prev')" :aria-label="t('calendar.prevMonth')">
        {{ t('calendar.prevMonth') }}
      </button>
      <h2 class="month-title">{{ t('calendar.monthTitle', { year, month: t(`months.${month}`) }) }}</h2>
      <button class="nav-btn" @click="$emit('next')" :aria-label="t('calendar.nextMonth')">
        {{ t('calendar.nextMonth') }}
      </button>
      <button class="today-btn" @click="$emit('today')" :aria-label="t('calendar.today')">
        {{ t('calendar.today') }}
      </button>
    </div>
    <div class="filter-controls">
      <select :value="selectedAccount" @change="$emit('update:selectedAccount', ($event.target as HTMLSelectElement).value)" class="filter-select">
        <option value="">{{ t('calendar.allAccounts') }}</option>
        <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
          {{ acc.name }}
        </option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Account } from '@/types'
import { useI18n } from '@/locales'

const { t } = useI18n()

defineProps<{
  year: number
  month: number
  selectedAccount: string
  accounts: Account[]
}>()

defineEmits(['prev', 'next', 'today', 'update:selectedAccount'])
</script>

<style scoped>
@reference "tailwindcss";
.calendar-header {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  @apply flex flex-wrap justify-between items-center gap-4 p-4 rounded;
}
.nav-controls { @apply flex items-center gap-3; }
.nav-btn, .today-btn {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  @apply px-3 py-1.5 text-sm font-medium cursor-pointer transition-colors duration-200 ease-in-out rounded-sm;
}
.nav-btn:hover { background-color: var(--bg-hover); }
.today-btn { background-color: var(--accent); color: white; border-color: var(--accent); }
.month-title { color: var(--text-primary); @apply text-xl font-black; }
.filter-select {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  @apply px-3 py-1.5 text-sm cursor-pointer rounded-sm;
}
</style>