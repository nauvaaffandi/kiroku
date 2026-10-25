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
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow: var(--shadow);
}
.weekdays-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
  text-align: center;
  font-weight: 600;
  color: var(--text-secondary);
  font-size: 0.875rem;
  padding: 0.5rem 0;
}
.days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  grid-auto-rows: 1fr;
}
.day-cell {
  min-height: 110px;
  padding: 0.375rem;
  border-bottom: 1px solid var(--border-primary);
  border-right: 1px solid var(--border-primary);
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}
.day-cell:hover { background-color: var(--bg-hover); }
.day-cell.other-month { opacity: 0.4; background-color: var(--bg-secondary); }
.day-cell.today { background-color: rgba(74, 144, 217, 0.1); border: 2px solid var(--accent); }
.day-header { display: flex; justify-content: space-between; align-items: center; font-size: 0.75rem; }
.day-number { font-weight: 700; color: var(--text-primary); }
.day-balance { font-size: 0.625rem; font-weight: 500; }
.day-balance.text-success { color: var(--success); }
.day-balance.text-danger { color: var(--error); }
.day-content { display: flex; flex-direction: column; gap: 0.25rem; overflow: hidden; flex: 1; }
.txn-pill { padding: 0.125rem 0.375rem; border-radius: var(--radius-sm); font-size: 0.625rem; display: flex; justify-content: space-between; align-items: center; gap: 0.25rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.txn-pill.pill-income { background-color: var(--success-light); color: var(--success); }
.txn-pill.pill-expense { background-color: var(--error-light); color: var(--error); }
.txn-pill.pill-scheduled { background-color: var(--warning-light); color: var(--warning); border: 1px dashed var(--warning); }
.txn-name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; }
.txn-amt { font-weight: 600; white-space: nowrap; }
.more-indicator { font-size: 0.625rem; color: var(--text-muted); text-align: center; font-weight: 500; }
</style>