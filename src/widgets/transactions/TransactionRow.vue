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
        {{ translate(category) }}
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
import { useCategoryName } from '@/utils/categoryName'
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
const { translate } = useCategoryName()

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
@reference "tailwindcss";
.transaction-row:hover {
  background-color: var(--bg-hover);
}

.cell-date {
  color: var(--text-primary);
  @apply whitespace-nowrap text-xs;
}

.cell-type {
  @apply whitespace-nowrap;
}

.type-badge {
  @apply inline-block px-2 py-0.5 text-[0.625rem] font-semibold uppercase rounded-full;
}

.type-badge.income {
  background-color: var(--success-light);
  color: var(--success);
}

.type-badge.expense {
  background-color: var(--error-light);
  color: var(--error);
}

.type-badge.transfer {
  background-color: var(--accent-light);
  color: var(--accent);
}

.cell-amount {
  @apply font-semibold whitespace-nowrap text-[0.8125rem] tabular-nums;
}

.cell-amount.positive {
  color: var(--success);
}

.cell-amount.negative {
  color: var(--error);
}

.cell-amount.neutral {
  color: var(--text-primary);
}

.cell-category {
  @apply max-w-[140px] text-xs;
}

.category-badge {
  border-radius: var(--radius);
  @apply inline-flex items-center gap-1.5 px-2 py-0.5 text-xs whitespace-nowrap overflow-hidden text-ellipsis;
}

.category-icon {
  border-radius: var(--radius-sm);
  @apply w-3.5 h-3.5 flex items-center justify-center text-[0.6875rem];
}

.cell-account {
  @apply text-xs max-w-[130px];
}

.transfer-arrow {
  color: var(--accent);
  @apply ml-2 font-medium;
}

.cell-note {
  color: var(--text-secondary);
  @apply max-w-[160px] text-xs overflow-hidden text-ellipsis whitespace-nowrap;
}

.cell-status {
  @apply whitespace-nowrap;
}

.status-badge {
  @apply inline-block px-2 py-0.5 text-[0.625rem] font-semibold uppercase rounded-full;
}

.status-badge.cleared {
  background-color: var(--success-light);
  color: var(--success);
}

.status-badge.pending {
  background-color: var(--warning-light);
  color: var(--warning);
}

.cell-actions {
  @apply whitespace-nowrap w-[60px];
}

.action-buttons {
  @apply flex justify-end gap-1;
}

.btn-icon {
  color: var(--text-muted);
  border-radius: var(--radius);
  @apply flex items-center justify-center w-6 h-6 border-none bg-transparent cursor-pointer transition-colors duration-150 ease-in-out;
}

.btn-icon:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.btn-icon.text-danger:hover {
  color: var(--error);
  background-color: var(--error-light);
}

.text-muted {
  color: var(--text-muted);
}

@media (max-width: 480px) {
  .cell-date { @apply text-[0.6875rem]; }
  .cell-amount { @apply text-xs; }
  .cell-category { @apply max-w-[100px]; }
  .cell-account { @apply max-w-[90px]; }
  .cell-note { @apply max-w-[120px]; }
  .type-badge { @apply text-[0.5625rem] px-1.5 py-0.5; }
  .status-badge { @apply text-[0.5625rem] px-1.5 py-0.5; }
  .category-badge { @apply text-[0.6875rem]; }
  .category-icon { @apply w-3 h-3 text-[0.625rem]; }
  .btn-icon { @apply w-[22px] h-[22px]; }
}
</style>