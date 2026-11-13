import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { localDb } from '@/storage/repository'
import { useUIStore } from './ui'
import { useI18n } from '@/locales'
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
  const { t } = useI18n()

  const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.limit))

  const matchesFilters = (t: Transaction) => {
    const f = filters.value
    if (f.account_id && t.from_account_id !== f.account_id && t.to_account_id !== f.account_id) return false
    if (f.category_id && t.category_id !== f.category_id) return false
    if (f.type && t.type !== f.type) return false
    if (f.status && t.status !== f.status) return false
    if (f.date_from && t.date_time < f.date_from) return false
    if (f.date_to && t.date_time > f.date_to) return false
    if (f.amount_min && t.amount < Number(f.amount_min)) return false
    if (f.amount_max && t.amount > Number(f.amount_max)) return false
    if (f.payee && !(t.payee_payer ?? '').toLowerCase().includes(f.payee.toLowerCase())) return false
    if (f.q) {
      const q = f.q.toLowerCase()
      const text = [t.note, t.payee_payer, t.place_name].join(' ').toLowerCase()
      if (!text.includes(q)) return false
    }
    if (f.label) {
      const labels = Array.isArray(t.labels) ? t.labels : String(t.labels ?? '').split(',')
      if (!labels.map((l) => l.trim()).includes(f.label)) return false
    }
    return true
  }

  const fetchTransactions = async (page = 1) => {
    loading.value = true
    try {
      const all = await localDb.list<Transaction>('transaction')
      const filtered = all.filter(matchesFilters)
      filtered.sort((a, b) => {
        const d = filters.value.order === 'DESC' ? -1 : 1
        const key = filters.value.sort as keyof Transaction
        const av = String(a[key] ?? '')
        const bv = String(b[key] ?? '')
        return av < bv ? -d : av > bv ? d : 0
      })
      pagination.value.total = filtered.length
      pagination.value.page = page
      const start = (page - 1) * pagination.value.limit
      transactions.value = filtered.slice(start, start + pagination.value.limit)
    } catch (e) {
      uiStore.error(t('transactions.fetchError'))
    } finally {
      loading.value = false
    }
  }

  const createTransaction = async (data: Partial<Transaction>) => {
    try {
      const id = crypto.randomUUID()
      const now = new Date().toISOString()
      const transaction: Transaction = {
        id,
        type: 'expense',
        date_time: now,
        amount: 0,
        currency: '',
        from_account_id: null,
        to_account_id: null,
        category_id: null,
        labels: '',
        note: '',
        payee_payer: '',
        payment_type: '',
        warranty_until: null,
        status: 'cleared',
        place_name: '',
        place_lat: null,
        place_lng: null,
        attachments: '',
        template_id: null,
        created_at: now,
        updated_at: now,
        ...data
      } as Transaction
      await localDb.put('transaction', id, transaction)
      await fetchTransactions(pagination.value.page)
      return id
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to create transaction')
      throw e
    }
  }

  const updateTransaction = async (id: string, data: Partial<Transaction>) => {
    try {
      const current = await localDb.get<Transaction>('transaction', id)
      if (!current) throw new Error('Transaction not found')
      await localDb.put('transaction', id, { ...current, ...data, updated_at: new Date().toISOString() } as Transaction)
      await fetchTransactions(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to update transaction')
      throw e
    }
  }

  const deleteTransaction = async (id: string) => {
    try {
      await localDb.remove('transaction', id)
      await fetchTransactions(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to delete transaction')
      throw e
    }
  }

  const batchOperation = async (ids: string[], action: 'delete' | 'recategorize' | 'relabel', category_id?: string, labels?: string) => {
    try {
      for (const id of ids) {
        if (action === 'delete') {
          await localDb.remove('transaction', id)
          continue
        }
        const current = await localDb.get<Transaction>('transaction', id)
        if (!current) continue
        const next = { ...current, updated_at: new Date().toISOString() }
        if (action === 'recategorize') next.category_id = category_id ?? null
        if (action === 'relabel') next.labels = labels ?? ''
        await localDb.put('transaction', id, next)
      }
      await fetchTransactions(pagination.value.page)
      return { updated: ids.length }
    } catch (e: any) {
      uiStore.error(e?.message || 'Batch operation failed')
      throw e
    }
  }

  const fetchSummary = async () => {
    try {
      const all = await localDb.list<Transaction>('transaction')
      const filtered = all.filter(matchesFilters)
      const result: TransactionSummary = {}
      for (const t of filtered) {
        const key = t.currency || 'unknown'
        if (!result[key]) result[key] = { income: 0, expense: 0, net: 0 }
        if (t.type === 'income') result[key].income += t.amount
        else if (t.type === 'expense') result[key].expense += t.amount
        result[key].net = result[key].income - result[key].expense
      }
      summary.value = result
    } catch (e) {
      console.error('Failed to compute summary:', e)
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
