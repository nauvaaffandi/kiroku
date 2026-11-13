<template>
  <div class="transactions-page">
    <TransactionsHeader
      :loading="transactionsStore.loading"
      @export-csv="exportCSV"
      @export-excel="exportExcel"
      @scan-receipt="openScanReceipt"
      @create="openCreateModal"
    />
    <TransactionFilters
      v-model="localFilters"
      :categories="categories"
      :accounts="accounts"
      :searchQuery="searchQuery"
      @change="handleFilterChange"
      @search="handleSearch"
    />
    <div class="table-wrap">
      <TransactionTable
        :transactions="transactions"
        :loading="transactionsStore.loading"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </div>
    <TransactionsPagination
      :currentPage="currentPage"
      :totalPages="totalPages"
      @prev="prevPage"
      @next="nextPage"
    />
    <TransactionFormModal
      v-model="showFormModal"
      :transaction="editingTransaction"
      :accounts="accounts"
      :categories="categories"
      :initial-date-time="initialDateTime"
      @submit="handleFormSubmit"
    />
    <BaseModal v-model="showDeleteConfirm" title="{{ t('transactions.deleteTransaction') }}" size="sm">
      <p>{{ t('transactions.deleteConfirm') }}</p>
      <template #actions>
        <BaseButton variant="secondary" @click="closeDeleteConfirm">{{ t('common.cancel') }}</BaseButton>
        <BaseButton variant="danger" @click="executeDelete" :disabled="transactionsStore.loading">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { TransactionTable, TransactionFormModal, TransactionFilters, TransactionsHeader, TransactionsPagination, useTransactionsView, useExport } from '@/widgets/transactions'
import { BaseButton, BaseModal } from '@/components/ui'

const { t } = useI18n()
const { exportCSV, exportExcel } = useExport()
const {
  transactionsStore, transactions, accounts, categories, currentPage, totalPages,
  showFormModal, editingTransaction, showDeleteConfirm, initialDateTime, localFilters,
  searchQuery, handleFilterChange, handleSearch, handleEdit, handleDelete,
  openCreateModal, openScanReceipt, closeDeleteConfirm, executeDelete,
  handleFormSubmit, prevPage, nextPage
} = useTransactionsView()
</script>

<style scoped>
@reference "tailwindcss";

.transactions-page {
  @apply min-h-screen flex flex-col;
}

.table-wrap {
  @apply w-full overflow-x-auto m-4;
}

@media (max-width: 768px) {
  .table-wrap {
    @apply my-3 mx-0;
  }
}
</style>