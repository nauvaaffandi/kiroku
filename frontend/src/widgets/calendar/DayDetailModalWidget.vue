<template>
  <div v-if="activeDay" class="modal-backdrop" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <h3>{{ t('calendar.dayTransactions', { date: activeDay.date }) }}</h3>
      <div class="detail-list">
        <div v-for="txn in activeDay.transactions" :key="txn.id" class="detail-item">
          <span :class="['badge', txn.type === 'income' ? 'badge-success' : 'badge-danger']">
            {{ txn.type === 'income' ? t('transactions.typeIncome') : t('transactions.typeExpense') }}
          </span>
          <div class="detail-info">
            <span class="detail-name">{{ txn.name }}</span>
            <span class="detail-cat">{{ txn.category }}</span>
          </div>
          <span :class="['detail-amt', txn.type === 'income' ? 'text-success' : 'text-danger']">
            {{ txn.type === 'income' ? '+' : '-' }}{{ formatCurrency(txn.amount) }}
          </span>
        </div>
        <div v-for="txn in activeDay.scheduled" :key="txn.id" class="detail-item scheduled-item">
          <span class="badge badge-warning">{{ t('calendar.scheduledBadge') }}</span>
          <div class="detail-info">
            <span class="detail-name">{{ txn.name }}</span>
            <span class="detail-cat">{{ txn.category }}</span>
          </div>
          <span class="detail-amt">{{ formatCurrency(txn.amount) }}</span>
        </div>
        <div v-if="activeDay.transactions.length === 0 && activeDay.scheduled.length === 0" class="empty-msg">
          {{ t('calendar.noTransactions') }}
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn-primary" @click="$emit('add-transaction', activeDay.date)">{{ t('calendar.addTransaction') }}</button>
        <button class="btn-secondary" @click="$emit('close')">{{ t('common.close') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DisplayDay } from '@/types'
import { useI18n } from '@/locales'

const { t } = useI18n()

defineProps<{
  activeDay: DisplayDay | null
}>()

defineEmits(['close', 'add-transaction'])

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
.modal-backdrop {
  background-color: rgba(0, 0, 0, 0.5);
  @apply fixed inset-0 flex items-center justify-center z-50 p-4;
}
.modal-content {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  box-shadow: var(--shadow-xl);
  @apply rounded max-w-md w-full p-6 flex flex-col gap-4;
}
.modal-content h3 {
  color: var(--text-primary);
  @apply text-lg font-semibold m-0;
}
.detail-list {
  @apply flex flex-col gap-2 max-h-[300px] overflow-y-auto;
}
.detail-item {
  background-color: var(--bg-secondary);
  @apply flex items-center justify-between p-2 rounded-sm gap-3;
}
.detail-info {
  @apply flex flex-col flex-1 min-w-0;
}
.detail-name {
  color: var(--text-primary);
  @apply text-sm font-medium;
}
.detail-cat {
  color: var(--text-muted);
  @apply text-xs;
}
.detail-amt {
  @apply text-sm font-bold;
}
.detail-amt.text-success {
  color: var(--success);
}
.detail-amt.text-danger {
  color: var(--error);
}
.badge {
  @apply px-2 py-0.5 text-xs rounded-full font-medium;
}
.badge-success {
  background-color: var(--success-light);
  color: var(--success);
}
.badge-danger {
  background-color: var(--error-light);
  color: var(--error);
}
.badge-warning {
  background-color: var(--warning-light);
  color: var(--warning);
}
.empty-msg {
  color: var(--text-muted);
  @apply text-center p-4 text-sm;
}
.modal-actions {
  @apply flex justify-end gap-2 mt-2;
}
.btn-primary {
  background-color: var(--accent);
  color: white;
  @apply px-4 py-2 rounded-sm text-sm font-medium cursor-pointer;
}
.btn-primary:hover {
  background-color: var(--accent-hover);
}
.btn-secondary {
  background-color: var(--bg-card);
  color: var(--text-primary);
  border: 1px solid var(--border-primary);
  @apply px-4 py-2 rounded-sm text-sm font-medium cursor-pointer;
}
.btn-secondary:hover {
  background-color: var(--bg-hover);
}
</style>