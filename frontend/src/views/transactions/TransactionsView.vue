<template>
  <div class="transactions-page">
    <header class="page-header">
      <div>
        <h1>{{ t('transactions.title') }}</h1>
        <p class="page-subtitle">{{ t('transactions.addTransaction') }}</p>
      </div>
      <div class="header-actions">
        <BaseButton variant="secondary" @click="exportCSV" :disabled="transactionsStore.loading">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
          {{ t('transactions.exportCsv') }}
        </BaseButton>
        <BaseButton variant="secondary" @click="exportExcel" :disabled="transactionsStore.loading">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
          </svg>
          {{ t('transactions.exportExcel') }}
        </BaseButton>
        <BaseButton variant="secondary" @click="openScanReceipt">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
            <path d="M17 18a5 5 0 0 0-10 0"></path>
            <path d="M12 2a10 10 0 1 0 10 10"></path>
          </svg>
          {{ t('transactions.scanReceipt') }}
        </BaseButton>
        <BaseButton @click="openCreateModal">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          {{ t('transactions.addTransaction') }}
        </BaseButton>
      </div>
    </header>

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

    <div class="pagination" v-if="totalPages > 1">
      <BaseButton variant="secondary" @click="prevPage" :disabled="currentPage === 1">
        {{ t('common.prev') }}
      </BaseButton>
      <span class="page-info">{{ t('common.page', { current: currentPage, total: totalPages }) }}</span>
      <BaseButton variant="secondary" @click="nextPage" :disabled="currentPage >= totalPages">
        {{ t('common.next') }}
      </BaseButton>
    </div>

    <TransactionFormModal
      v-model="showFormModal"
      :transaction="editingTransaction"
      :accounts="accounts"
      :categories="categories"
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import { useCategoriesStore } from '@/stores/categories'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'
import { api } from '@/api/client'
import { exportToExcel } from '@/utils/exportExcel'
import type { Transaction, TransactionFilters as ITransactionFilters } from '@/types'
import { TransactionTable, TransactionFormModal, TransactionFilters } from '@/widgets/transactions'
import { BaseButton, BaseModal } from '@/components/ui'

const router = useRouter()
const route = useRoute()
const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const categoriesStore = useCategoriesStore()
const uiStore = useUIStore()
const { t } = useI18n()

const transactions = computed(() => transactionsStore.transactions || [])
const accounts = computed(() => accountsStore.accounts || [])
const categories = computed(() => categoriesStore.categories || [])
const currentPage = computed(() => transactionsStore.pagination.page)
const totalPages = computed(() => transactionsStore.totalPages)

const showFormModal = ref(false)
const editingTransaction = ref<Transaction | null>(null)
const showDeleteConfirm = ref(false)
const deleteId = ref<string | null>(null)

const localFilters = ref<ITransactionFilters>({
  account_id: '',
  category_id: '',
  type: '',
  status: '',
  date_from: '',
  date_to: '',
  amount_min: '',
  amount_max: '',
  q: '',
  label: '',
  payee: '',
  sort: 'date_time',
  order: 'DESC'
})

const searchQuery = ref('')

let searchTimeout: ReturnType<typeof setTimeout>
watch(searchQuery, (value) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    transactionsStore.filters.q = value
    fetchTransactions(1)
  }, 300)
})

const fetchTransactions = async (page = 1) => {
  await transactionsStore.fetchTransactions(page)
  await transactionsStore.fetchSummary()
}

const handleFilterChange = () => {
  transactionsStore.filters = { ...localFilters.value }
  fetchTransactions(1)
}

const handleSearch = (value: string) => {
  searchQuery.value = value
}

const handleEdit = (tx: Transaction) => {
  editingTransaction.value = tx
  showFormModal.value = true
}

const handleDelete = (id: string) => {
  deleteId.value = id
  showDeleteConfirm.value = true
}

const openCreateModal = () => {
  editingTransaction.value = null
  showFormModal.value = true
}

const openScanReceipt = () => {
  router.push('/transactions/scan')
}

const closeDeleteConfirm = () => {
  showDeleteConfirm.value = false
  deleteId.value = null
}

const executeDelete = async () => {
  if (deleteId.value) {
    try {
      await transactionsStore.deleteTransaction(deleteId.value)
      uiStore.success('Transaction deleted')
    } catch (e) {
      // Error handled in store
    }
    closeDeleteConfirm()
  }
}

const handleFormSubmit = async (data: Partial<Transaction>) => {
  try {
    if (editingTransaction.value) {
      await transactionsStore.updateTransaction(editingTransaction.value.id, data)
      uiStore.success('Transaction updated')
    } else {
      await transactionsStore.createTransaction(data)
      uiStore.success('Transaction created')
    }
    showFormModal.value = false
    editingTransaction.value = null
  } catch (e) {
    // Error handled in store
  }
}

const exportCSV = async () => {
  try {
    const params = { ...transactionsStore.filters }
    const response = await api.get('/api/v1/transactions/export', {
      params,
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `transactions-${new Date().toISOString().split('T')[0]}.csv`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    uiStore.success('CSV exported successfully')
  } catch (e) {
    uiStore.error('Failed to export CSV')
  }
}

const exportExcel = async () => {
  try {
    const params = { page: 1, limit: 10000, ...transactionsStore.filters }
    const response = await api.get('/api/v1/transactions', { params })
    const data = Array.isArray(response.data) ? response.data : []
    const headers = ['Date', 'Type', 'Amount', 'Category', 'Account', 'Note', 'Status']
    const rows = data.map((tx: any) => [
      tx.date_time,
      tx.type,
      tx.amount,
      tx.category_name || '',
      tx.from_account_name || '',
      tx.note || '',
      tx.status
    ])
    exportToExcel(headers, rows, `transactions-${new Date().toISOString().split('T')[0]}`)
    uiStore.success('Excel exported successfully')
  } catch (e) {
    uiStore.error('Failed to export Excel')
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    fetchTransactions(currentPage.value - 1)
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    fetchTransactions(currentPage.value + 1)
  }
}

onMounted(async () => {
  await Promise.all([
    fetchTransactions(),
    accountsStore.fetchAccounts(),
    categoriesStore.fetchCategories()
  ])

  if (route.query.create === 'true' || route.query.type) {
    openCreateModal()
    router.replace({ query: { ...route.query, create: undefined, type: undefined } })
  }
})
</script>

<style scoped>
.transactions-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1rem;
  background-color: var(--color-neutral-100);
  border-bottom: 1px solid var(--color-neutral-200);
  flex-wrap: wrap;
  gap: 0.75rem;
}

.page-header h1 {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-neutral-900);
}

.page-subtitle {
  color: var(--color-neutral-500);
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.table-wrap {
  width: 100%;
  overflow-x: auto;
  margin: 1rem;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 1rem;
  border-top: 1px solid var(--color-neutral-200);
}

.page-info {
  color: var(--color-neutral-600);
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    padding: 0.75rem;
  }
  .header-actions {
    justify-content: stretch;
    flex-direction: column;
  }
  .header-actions button {
    width: 100%;
    justify-content: center;
  }
  .table-wrap {
    margin: 0.75rem 0;
  }
  .pagination {
    padding: 0.75rem;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  .pagination button {
    flex: 1;
  }
  .page-info {
    flex: 1;
    text-align: center;
  }
}
</style>