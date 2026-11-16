<template>
  <div class="table-container">
    <table class="data-table">
      <thead>
        <tr>
          <th class="cell-date-head">{{ t('transactions.columnDate') }}</th>
          <th class="cell-type-head">{{ t('transactions.columnType') }}</th>
          <th class="cell-amount-head">{{ t('transactions.columnAmount') }}</th>
          <th class="cell-category-head">{{ t('transactions.columnCategory') }}</th>
          <th class="cell-account-head">{{ t('transactions.columnAccount') }}</th>
          <th class="cell-note-head">{{ t('transactions.columnNote') }}</th>
          <th class="cell-status-head">{{ t('transactions.columnStatus') }}</th>
          <th class="cell-actions-head">{{ t('transactions.columnActions') }}</th>
        </tr>
      </thead>
      <tbody>
        <template v-if="transactions.length > 0">
          <TransactionRow
            v-for="tx in transactions"
            :key="tx.id"
            :transaction="tx"
            @edit="$emit('edit', $event)"
            @delete="$emit('delete', $event)"
          />
        </template>
        <tr v-else-if="loading">
          <td colspan="8" class="td-empty">{{ t('common.loading') }}</td>
        </tr>
        <tr v-else>
          <td colspan="8" class="td-empty">{{ t('transactions.noTransactions') }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import TransactionRow from './TransactionRow.vue'
import { useI18n } from '@/locales'
import type { Transaction } from '@/types'

interface Props {
  transactions: Transaction[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false
})

defineEmits<{
  edit: [transaction: Transaction]
  delete: [id: string]
}>()

const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.table-container {
  @apply w-full overflow-x-auto;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
}

.data-table {
  @apply w-full border-collapse text-left text-sm;
}

.data-table th {
  @apply bg-[var(--bg-secondary)] text-[var(--text-secondary)] border-b-[var(--border-primary)] px-4 py-3 font-semibold whitespace-nowrap;
  border-bottom-width: 1px;
  border-bottom-style: solid;
}

.data-table td {
  @apply border-b-[var(--border-primary)] px-4 py-3;
  border-bottom-width: 1px;
  border-bottom-style: solid;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.cell-amount-head,
.cell-actions-head {
  @apply text-right;
}

.td-empty {
  @apply text-center py-8 text-[var(--text-muted)];
}
</style>
