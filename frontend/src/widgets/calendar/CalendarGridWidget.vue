<template>
  <div class="calendar-grid-wrapper">
    <div class="weekdays-header">
      <div v-for="day in weekdays" :key="day" class="weekday-cell">{{ day }}</div>
    </div>
    <div class="days-grid">
      <div
        v-for="(day, index) in calendarDays"
        :key="index"
        :class="['day-cell', { 'other-month': !day.isCurrentMonth, 'today': day.isToday }]"
        @click="$emit('open-detail', day)"
      >
        <div class="day-header">
          <span class="day-number">{{ day.dayNumber }}</span>
          <span v-if="day.projectedBalance !== undefined" :class="['day-balance', day.projectedBalance >= 0 ? 'text-success' : 'text-danger']">
            {{ formatCurrency(day.projectedBalance) }}
          </span>
        </div>
        <div class="day-content">
          <div
            v-for="txn in day.transactions.slice(0, 2)"
            :key="txn.id"
            :class="['txn-pill', txn.type === 'income' ? 'pill-income' : 'pill-expense']"
          >
            <span class="txn-name">{{ txn.name }}</span>
            <span class="txn-amt">{{ formatCurrency(txn.amount) }}</span>
          </div>
          <div
            v-for="txn in day.scheduled.slice(0, 1)"
            :key="txn.id"
            class="txn-pill pill-scheduled"
          >
            <span class="txn-name">予定: {{ txn.name }}</span>
            <span class="txn-amt">{{ formatCurrency(txn.amount) }}</span>
          </div>
          <div v-if="day.transactions.length + day.scheduled.length > 3" class="more-indicator">
            +{{ day.transactions.length + day.scheduled.length - 3 }} 件
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DisplayDay } from '@/types'

defineProps<{
  calendarDays: DisplayDay[]
}>()

defineEmits(['open-detail'])

const weekdays = ['日', '月', '火', '水', '木', '金', '土']

function formatCurrency(val: number): string {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(val)
}
</script>

<style scoped>
@reference "tailwindcss";

.calendar-grid-wrapper {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  @apply overflow-hidden rounded shadow;
}

.weekdays-header {
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
  color: var(--text-secondary);
  @apply grid grid-cols-7 text-center font-semibold text-sm px-0 py-2;
}

.days-grid {
  @apply grid grid-cols-7 grid-rows-[1fr];
}

.day-cell {
  border-color: var(--border-primary);
  @apply min-h-[110px] p-1.5 flex flex-col gap-1 cursor-pointer transition-colors duration-200 ease-in-out border-b border-r;
}

.day-cell:hover {
  background-color: var(--bg-hover);
}

.day-cell.other-month {
  background-color: var(--bg-secondary);
  @apply opacity-40;
}

.day-cell.today {
  border-color: var(--accent);
  @apply border-2;
  background-color: rgba(74, 144, 217, 0.1);
}

.day-header {
  @apply flex justify-between items-center text-xs;
}

.day-number {
  color: var(--text-primary);
  @apply font-bold;
}

.day-balance {
  @apply text-[0.625rem] font-medium;
}

.day-balance.text-success {
  color: var(--success);
}

.day-balance.text-danger {
  color: var(--error);
}

.day-content {
  @apply flex flex-col gap-2 overflow-hidden flex-1;
}

.txn-pill {
  @apply px-1.5 py-0.5 text-[0.625rem] flex items-center justify-between gap-1 rounded-sm overflow-hidden whitespace-nowrap text-ellipsis;
}

.txn-pill.pill-income {
  background-color: var(--success-light);
  color: var(--success);
}

.txn-pill.pill-expense {
  background-color: var(--error-light);
  color: var(--error);
}

.txn-pill.pill-scheduled {
  background-color: var(--warning-light);
  color: var(--warning);
  border: 1px dashed var(--warning);
}

.txn-name {
  @apply overflow-hidden text-ellipsis whitespace-nowrap flex-1;
}

.txn-amt {
  @apply font-semibold whitespace-nowrap;
}

.more-indicator {
  color: var(--text-muted);
  @apply text-[0.625rem] font-medium text-center;
}
</style>