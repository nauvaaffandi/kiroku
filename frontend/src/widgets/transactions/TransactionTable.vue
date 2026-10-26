<template>
  <div class="table-container">
    <table class="data-table">
      <thead>
        <tr>
          <th>Date</th>
          <th>Type</th>
          <th>Amount</th>
          <th>Category</th>
          <th>Account</th>
          <th>Note</th>
          <th>Status</th>
          <th class="text-right">Actions</th>
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
          <td colspan="8" class="text-center py-8 text-muted">Loading transactions...</td>
        </tr>
        <tr v-else>
          <td colspan="8" class="text-center py-8 text-muted">No transactions found</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import TransactionRow from './TransactionRow.vue'
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
</script>

<style scoped>
@reference "tailwindcss";
.table-container {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  @apply w-full overflow-x-auto;
}

.data-table {
  @apply w-full border-collapse text-left text-sm;
}

.data-table th {
  background-color: var(--bg-secondary);
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-primary);
  @apply px-4 py-3 font-semibold whitespace-nowrap;
}

.data-table td {
  border-bottom: 1px solid var(--border-primary);
  @apply px-4 py-3;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.text-right {
  @apply text-right;
}

.text-center {
  @apply text-center;
}

.py-8 {
  @apply py-8;
}

.text-muted {
  color: var(--text-muted);
}
</style>