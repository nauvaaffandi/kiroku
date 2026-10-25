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
  width: 100%;
  overflow-x: auto;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
}

.data-table th {
  padding: 0.75rem 1rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  background: var(--color-surface-variant);
  border-bottom: 1px solid var(--color-border);
  white-space: nowrap;
}

.data-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--color-border);
}

.data-table tr:last-child td {
  border-bottom: none;
}

.text-right {
  text-align: right;
}

.text-center {
  text-align: center;
}

.py-8 {
  padding-top: 2rem;
  padding-bottom: 2rem;
}

.text-muted {
  color: var(--color-text-muted);
}
</style>