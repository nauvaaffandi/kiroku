<template>
  <div v-if="activeDay" class="modal-backdrop" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <h3>{{ activeDay.date }} の取引一覧</h3>
      <div class="detail-list">
        <div v-for="txn in activeDay.transactions" :key="txn.id" class="detail-item">
          <span :class="['badge', txn.type === 'income' ? 'badge-success' : 'badge-danger']">
            {{ txn.type === 'income' ? '収入' : '支出' }}
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
          <span class="badge badge-warning">予定</span>
          <div class="detail-info">
            <span class="detail-name">{{ txn.name }}</span>
            <span class="detail-cat">{{ txn.category }}</span>
          </div>
          <span class="detail-amt">{{ formatCurrency(txn.amount) }}</span>
        </div>
        <div v-if="activeDay.transactions.length === 0 && activeDay.scheduled.length === 0" class="empty-msg">
          この日の取引はありません
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn-primary" @click="$emit('add-transaction', activeDay.date)">取引を追加</button>
        <button class="btn-secondary" @click="$emit('close')">閉じる</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DisplayDay } from '@/types'

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
.modal-backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
  padding: 1rem;
}
.modal-content {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 1.5rem;
  max-width: 28rem;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  box-shadow: var(--shadow-xl);
}
.modal-content h3 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); margin: 0; }
.detail-list { display: flex; flex-direction: column; gap: 0.5rem; max-height: 300px; overflow-y: auto; }
.detail-item { display: flex; align-items: center; justify-content: space-between; padding: 0.5rem; border-radius: var(--radius-sm); background-color: var(--bg-secondary); gap: 0.75rem; }
.detail-info { display: flex; flex-direction: column; flex: 1; min-width: 0; }
.detail-name { font-size: 0.875rem; font-weight: 500; color: var(--text-primary); }
.detail-cat { font-size: 0.75rem; color: var(--text-muted); }
.detail-amt { font-size: 0.875rem; font-weight: 700; }
.detail-amt.text-success { color: var(--success); }
.detail-amt.text-danger { color: var(--error); }
.badge { padding: 0.125rem 0.5rem; font-size: 0.75rem; border-radius: 9999px; font-weight: 500; }
.badge-success { background-color: var(--success-light); color: var(--success); }
.badge-danger { background-color: var(--error-light); color: var(--error); }
.badge-warning { background-color: var(--warning-light); color: var(--warning); }
.empty-msg { text-align: center; color: var(--text-muted); padding: 1rem; font-size: 0.875rem; }
.modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 0.5rem; }
.btn-primary { padding: 0.5rem 1rem; background-color: var(--accent); color: white; border-radius: var(--radius-sm); font-size: 0.875rem; font-weight: 500; cursor: pointer; }
.btn-primary:hover { background-color: var(--accent-hover); }
.btn-secondary { padding: 0.5rem 1rem; border: 1px solid var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); border-radius: var(--radius-sm); font-size: 0.875rem; font-weight: 500; cursor: pointer; }
.btn-secondary:hover { background-color: var(--bg-hover); }
</style>