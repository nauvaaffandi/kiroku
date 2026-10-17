<template>
  <tr :class="['transaction-row', `type-${transaction.type}`, transaction.status]">
    <td class="cell-date">
      <time>{{ formatDate(transaction.date_time) }}</time>
    </td>
    <td class="cell-type">
      <span :class="['type-badge', transaction.type]">{{ transaction.type }}</span>
    </td>
    <td class="cell-amount" :class="amountClass">
      {{ formatAmount(transaction.amount) }}
    </td>
    <td class="cell-category">
      <span v-if="category" class="category-badge" :style="{ backgroundColor: category.color + '20', color: category.color }">
        <span class="category-icon" :style="{ backgroundColor: category.color }">{{ category.icon }}</span>
        {{ category.name }}
      </span>
      <span v-else class="text-muted">—</span>
    </td>
    <td class="cell-account">
      <span v-if="fromAccount">{{ fromAccount.name }}</span>
      <span v-else class="text-muted">—</span>
      <span v-if="transaction.type === 'transfer' && toAccount" class="transfer-arrow">→ {{ toAccount.name }}</span>
    </td>
    <td class="cell-note">{{ transaction.note || '—' }}</td>
    <td class="cell-status">
      <span :class="['status-badge', transaction.status]">{{ transaction.status }}</span>
    </td>
    <td class="cell-actions">
      <div class="action-buttons">
        <button @click="$emit('edit', transaction)" class="btn-icon" aria-label="Edit">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
          </svg>
        </button>
        <button @click="$emit('delete', transaction.id)" class="btn-icon text-danger" aria-label="Delete">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
        </button>
      </div>
    </td>
  </tr>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCategoriesStore } from '@/stores/categories'
import { useAccountsStore } from '@/stores/accounts'
import { formatAmountForDisplay } from '@/utils/amountDisplay'
import type { Transaction } from '@/types'

interface Props {
  transaction: Transaction
}

const props = defineProps<Props>()
defineEmits<{
  edit: [transaction: Transaction]
  delete: [id: string]
}>()

const categoriesStore = useCategoriesStore()
const accountsStore = useAccountsStore()

const category = computed(() => categoriesStore.categories.find(c => c.id === props.transaction.category_id))
const fromAccount = computed(() => accountsStore.accounts.find(a => a.id === props.transaction.from_account_id))
const toAccount = computed(() => accountsStore.accounts.find(a => a.id === props.transaction.to_account_id))

const amountClass = computed(() => {
  if (props.transaction.type === 'income') return 'positive'
  if (props.transaction.type === 'expense') return 'negative'
  return 'neutral'
})

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}

const formatAmount = (amount: number) => {
  return formatAmountForDisplay(amount, props.transaction.type as 'income' | 'expense' | 'transfer')
}
</script>

<style scoped>
.transaction-row:hover {
  background: var(--color-surface-hover);
}

.cell-date {
  white-space: nowrap;
  font-size: 0.8125rem;
  color: var(--color-text);
}

.cell-type {
  white-space: nowrap;
}

.type-badge {
  display: inline-block;
  padding: 0.125rem 0.5rem;
  font-size: 0.6875rem;
  font-weight: 600;
  text-transform: uppercase;
  border-radius: var(--radius-full);
}

.type-badge.income {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.type-badge.expense {
  background: var(--color-danger-bg);
  color: var(--color-danger);
}

.type-badge.transfer {
  background: var(--color-primary-bg);
  color: var(--color-primary);
}

.cell-amount {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.cell-amount.positive {
  color: var(--color-success);
}

.cell-amount.negative {
  color: var(--color-danger);
}

.cell-amount.neutral {
  color: var(--color-text);
}

.cell-category {
  max-width: 180px;
}

.category-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.125rem 0.5rem;
  font-size: 0.8125rem;
  border-radius: var(--radius);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.category-icon {
  width: 1rem;
  height: 1rem;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
}

.cell-account {
  font-size: 0.8125rem;
  max-width: 160px;
}

.transfer-arrow {
  margin-left: 0.5rem;
  color: var(--color-primary);
  font-weight: 500;
}

.cell-note {
  max-width: 200px;
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cell-status {
  white-space: nowrap;
}

.status-badge {
  display: inline-block;
  padding: 0.125rem 0.5rem;
  font-size: 0.6875rem;
  font-weight: 600;
  text-transform: uppercase;
  border-radius: var(--radius-full);
}

.status-badge.cleared {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.status-badge.pending {
  background: var(--color-warning-bg);
  color: var(--color-warning);
}

.cell-actions {
  white-space: nowrap;
  width: 80px;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 0.25rem;
}

.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: var(--radius);
  color: var(--color-text-muted);
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.btn-icon:hover {
  background: var(--color-surface-hover);
  color: var(--color-text);
}

.btn-icon.text-danger:hover {
  color: var(--color-danger);
  background: var(--color-danger-bg);
}

.text-muted {
  color: var(--color-text-muted);
}
</style>