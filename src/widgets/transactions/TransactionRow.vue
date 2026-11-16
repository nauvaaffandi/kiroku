<template>
  <tr :class="['transaction-row', `type-${transaction.type}`, transaction.status]">
    <td class="cell-date">
      <time>{{ formatDate(transaction.date_time) }}</time>
    </td>
    <td class="cell-type">
      <span :class="['type-badge', typeBadgeClass(transaction.type)]">{{ typeLabel(transaction.type) }}</span>
    </td>
    <td class="cell-amount" :class="amountClass">
      {{ formatRowAmount(transaction.amount, transaction.type) }}
    </td>
    <td class="cell-category">
      <span v-if="category" class="category-cell">
        <span class="category-dot" :style="{ backgroundColor: category.color || 'var(--text-muted)' }" aria-hidden="true" />
        <span class="category-name">{{ translate(category) }}</span>
      </span>
      <span v-else class="placeholder-text">—</span>
    </td>
    <td class="cell-account">
      <span v-if="fromAccount">{{ fromAccount.name }}</span>
      <span v-else class="placeholder-text">—</span>
      <span v-if="transaction.type === 'transfer' && toAccount" class="transfer-arrow">→ {{ toAccount.name }}</span>
    </td>
    <td class="cell-note">{{ transaction.note || '—' }}</td>
    <td class="cell-status">
      <span :class="['status-badge', statusBadgeClass(transaction.status)]">{{ statusLabel(transaction.status) }}</span>
    </td>
    <td class="cell-actions">
      <div class="action-buttons">
        <button
          @click="$emit('edit', transaction)"
          class="btn-icon"
          :aria-label="t('transactions.editTransaction')"
        >
          <Edit :size="18" stroke-width="1.75" />
        </button>
        <button
          @click="$emit('delete', transaction.id)"
          class="btn-delete"
          :aria-label="t('transactions.deleteTransaction')"
        >
          <Trash2 :size="18" stroke-width="1.75" />
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
import { Edit, Trash2 } from '@lucide/vue'
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import type { Transaction } from '@/types'
import dayjs from 'dayjs'

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
const settingsStore = useSettingsStore()
const { translate } = useCategoryName()

const { t } = useI18n()

if (import.meta.env.DEV) console.log('[TransactionRow] rendering transaction', {
  id: props.transaction.id,
  type: props.transaction.type,
  amount: props.transaction.amount
})

const category = computed(() => categoriesStore.categories.find(c => c.id === props.transaction.category_id))
const fromAccount = computed(() => accountsStore.accounts.find(a => a.id === props.transaction.from_account_id))
const toAccount = computed(() => accountsStore.accounts.find(a => a.id === props.transaction.to_account_id))

const typeLabel = (type: string) => {
  return t(`transactions.type${type.charAt(0).toUpperCase() + type.slice(1)}`) || type
}

const typeBadgeClass = (type: string) => {
  return type === 'income' ? 'income' : type === 'expense' ? 'expense' : 'transfer'
}

const statusLabel = (status: string) => {
  return status === 'pending' ? t('transactions.statusPending') : t('transactions.statusCleared')
}

const statusBadgeClass = (status: string) => {
  return status === 'pending' ? 'pending' : 'cleared'
}

const amountClass = computed(() => {
  if (props.transaction.type === 'income') return 'positive'
  if (props.transaction.type === 'expense') return 'negative'
  return 'neutral'
})

const formatRowAmount = (amount: number, type: 'income' | 'expense' | 'transfer') => {
  if (import.meta.env.DEV) console.log('[TransactionRow] formatting amount', { amount, type })
  return formatAmountForDisplay(amount, type)
}

const formatDate = (dateStr: string) => {
  const fmt = settingsStore.dateFormat || 'YYYY-MM-DD'
  return dayjs(dateStr).format(fmt)
}
</script>

<style scoped>
@reference "tailwindcss";

.transaction-row {
  @apply hover:bg-[var(--bg-hover)] transition-colors duration-120 ease-out;
}

.cell-date {
  @apply whitespace-nowrap text-xs text-[var(--text-primary)];
}

.cell-type {
  @apply whitespace-nowrap;
}

.type-badge {
  @apply inline-block px-2 py-0.5 text-[0.625rem] font-medium rounded-full whitespace-nowrap;
}

.type-badge.income {
  @apply bg-[var(--success-light)] text-[var(--success)];
}

.type-badge.expense {
  @apply bg-[var(--error-light)] text-[var(--error)];
}

.type-badge.transfer {
  @apply bg-[var(--accent-light)] text-[var(--accent)];
}

.cell-amount {
  @apply font-semibold tabular-nums text-right text-xs;
}

.cell-amount.positive {
  @apply text-[var(--success)];
}

.cell-amount.negative,
.cell-amount.neutral {
  @apply text-[var(--text-primary)];
}

.cell-category {
  @apply max-w-[140px] text-xs;
}

.category-cell {
  @apply flex items-center gap-1.5;
}

.category-dot {
  @apply w-1.5 h-1.5 rounded-full shrink-0;
}

.category-name {
  @apply truncate text-[var(--text-primary)];
}

.placeholder-text {
  @apply text-[var(--text-muted)];
}

.cell-account {
  @apply text-xs max-w-[130px];
}

.transfer-arrow {
  @apply ml-2 font-medium text-[var(--accent)];
}

.cell-note {
  @apply max-w-[160px] text-xs overflow-hidden text-ellipsis whitespace-nowrap text-[var(--text-secondary)];
}

.cell-status {
  @apply whitespace-nowrap;
}

.status-badge {
  @apply inline-block px-2 py-0.5 text-[0.625rem] font-medium rounded-full whitespace-nowrap;
}

.status-badge.cleared {
  @apply bg-[var(--success-light)] text-[var(--success)];
}

.status-badge.pending {
  @apply bg-[var(--warning-light)] text-[var(--warning)];
}

.cell-actions {
  @apply whitespace-nowrap w-[60px];
}

.action-buttons {
  @apply flex justify-end gap-1;
}

.btn-icon {
  @apply flex items-center justify-center p-2.5 bg-transparent border-none cursor-pointer transition-colors duration-150 ease-in-out rounded-md;
  color: var(--text-muted);
  min-width: 44px;
  min-height: 44px;
}

.btn-icon:hover {
  @apply bg-[var(--bg-hover)];
  color: var(--text-primary);
}

.btn-delete {
  @apply flex items-center justify-center p-2.5 bg-transparent border-none cursor-pointer transition-colors duration-150 ease-in-out rounded-md;
  color: var(--text-muted);
  min-width: 44px;
  min-height: 44px;
}

.btn-delete:hover {
  @apply bg-[var(--error-light)];
  color: var(--error);
}
</style>
