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
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 1rem;
}
.nav-controls {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.nav-btn, .today-btn {
  padding: 0.375rem 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
  cursor: pointer;
}
.nav-btn:hover {
  background-color: var(--bg-hover);
}
.today-btn {
  background-color: var(--accent);
  color: white;
  border-color: var(--accent);
}
.month-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-primary);
}
.filter-select {
  padding: 0.375rem 0.75rem;
  font-size: 0.875rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}
</style>