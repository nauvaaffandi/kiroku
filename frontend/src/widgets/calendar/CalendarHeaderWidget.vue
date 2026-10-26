<template>
  <div class="calendar-header">
    <div class="nav-controls">
      <button class="nav-btn" @click="$emit('prev')">&lt; 前の月</button>
      <h2 class="month-title">{{ year }}年 {{ month }}月</h2>
      <button class="nav-btn" @click="$emit('next')">次の月 &gt;</button>
      <button class="today-btn" @click="$emit('today')">今日</button>
    </div>
    <div class="filter-controls">
      <select :value="selectedAccount" @change="$emit('update:selectedAccount', ($event.target as HTMLSelectElement).value)" class="filter-select">
        <option value="">すべての口座</option>
        <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
          {{ acc.name }}
        </option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Account } from '@/types'

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