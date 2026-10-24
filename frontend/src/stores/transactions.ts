import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/client'
import { useUIStore } from './ui'
import type { Transaction, TransactionFilters, TransactionSummary } from '@/types'

export const useTransactionsStore = defineStore('transactions', () => {
  const transactions = ref<Transaction[]>([])
  const loading = ref(false)
  const pagination = ref({ page: 1, limit: 20, total: 0 })
  const filters = ref<TransactionFilters>({
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
  const summary = ref<TransactionSummary>({})
  const uiStore = useUIStore()

  const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.limit))

  const fetchTransactions = async (page = 1) => {
    loading.value = true
    try {
      const params = { page, limit: pagination.value.limit, ...filters.value }
      const response = await api.get('/api/v1/transactions', { params })
      transactions.value = Array.isArray(response.data) ? response.data : []
      pagination.value.page = page
    } catch (e) {
      uiStore.error('Failed to fetch transactions')
    } finally {
      loading.value = false
    }
  }

  const createTransaction = async (data: Partial<Transaction>) => {
    try {
      const response = await api.post('/api/v1/transactions', data)
      await fetchTransactions(pagination.value.page)
      return response.data.id
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to create transaction')
      throw e
    }
  }

  const updateTransaction = async (id: string, data: Partial<Transaction>) => {
    try {
      await api.put(`/api/v1/transactions/${id}`, data)
      await fetchTransactions(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to update transaction')
      throw e
    }
  }

  const deleteTransaction = async (id: string) => {
    try {
      await api.delete(`/api/v1/transactions/${id}`)
      await fetchTransactions(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to delete transaction')
      throw e
    }
  }

  const batchOperation = async (ids: string[], action: 'delete' | 'recategorize' | 'relabel', category_id?: string, labels?: string) => {
    try {
      const response = await api.post('/api/v1/transactions/batch', { ids, action, category_id, labels })
      await fetchTransactions(pagination.value.page)
      return response.data
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Batch operation failed')
      throw e
    }
  }

  const fetchSummary = async () => {
    try {
      const response = await api.get('/api/v1/transactions/summary', { params: filters.value })
      summary.value = response.data.reduce((acc: TransactionSummary, item: any) => {
        acc[item.currency] = { income: item.income, expense: item.expense, net: item.net }
        return acc
      }, {})
    } catch (e) {
      console.error('Failed to fetch summary:', e)
    }
  }

  return {
    transactions,
    loading,
    pagination,
    filters,
    summary,
    totalPages,
    fetchTransactions,
    createTransaction,
    updateTransaction,
    deleteTransaction,
    batchOperation,
    fetchSummary
  }
})