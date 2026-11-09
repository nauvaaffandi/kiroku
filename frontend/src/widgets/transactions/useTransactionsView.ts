import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import { useCategoriesStore } from '@/stores/categories'
import { useUIStore } from '@/stores/ui'
import type { Transaction, TransactionFilters as ITransactionFilters } from '@/types'

export function useTransactionsView() {
  const router = useRouter()
  const route = useRoute()
  const transactionsStore = useTransactionsStore()
  const accountsStore = useAccountsStore()
  const categoriesStore = useCategoriesStore()
  const uiStore = useUIStore()

  const transactions = computed(() => transactionsStore.transactions || [])
  const accounts = computed(() => accountsStore.accounts || [])
  const categories = computed(() => categoriesStore.categories || [])
  const currentPage = computed(() => transactionsStore.pagination.page)
  const totalPages = computed(() => transactionsStore.totalPages)
  const showFormModal = ref(false)
  const editingTransaction = ref<Transaction | null>(null)
  const showDeleteConfirm = ref(false)
  const deleteId = ref<string | null>(null)
  const initialDateTime = computed(() => {
    const d = route.query.date as string | undefined
    return d ? `${d}T00:00` : undefined
  })
  const localFilters = ref<ITransactionFilters>({
    account_id: '', category_id: '', type: '', status: '', date_from: '', date_to: '',
    amount_min: '', amount_max: '', q: '', label: '', payee: '', sort: 'date_time', order: 'DESC'
  })
  const searchQuery = ref('')

  const fetchTransactions = async (page = 1) => {
    await transactionsStore.fetchTransactions(page)
    await transactionsStore.fetchSummary()
  }

  let searchTimeout: ReturnType<typeof setTimeout>
  watch(searchQuery, (value) => {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
      transactionsStore.filters.q = value
      fetchTransactions(1)
    }, 300)
  })

  const handleFilterChange = () => {
    transactionsStore.filters = { ...localFilters.value }
    fetchTransactions(1)
  }
  const handleSearch = (value: string) => { searchQuery.value = value }
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
  const openScanReceipt = () => { router.push('/transactions/scan') }
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
  const prevPage = () => {
    if (currentPage.value > 1) fetchTransactions(currentPage.value - 1)
  }
  const nextPage = () => {
    if (currentPage.value < totalPages.value) fetchTransactions(currentPage.value + 1)
  }

  onMounted(async () => {
    await Promise.all([
      fetchTransactions(),
      accountsStore.fetchAccounts(),
      categoriesStore.fetchCategories()
    ])
    if (route.query.create === 'true' || route.query.type) {
      openCreateModal()
      router.replace({ query: { ...route.query, create: undefined, date: undefined } })
    }
  })

  return {
    transactionsStore, transactions, accounts, categories, currentPage, totalPages,
    showFormModal, editingTransaction, showDeleteConfirm, initialDateTime, localFilters,
    searchQuery, handleFilterChange, handleSearch, handleEdit, handleDelete,
    openCreateModal, openScanReceipt, closeDeleteConfirm, executeDelete,
    handleFormSubmit, prevPage, nextPage
  }
}
