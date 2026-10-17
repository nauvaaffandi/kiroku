import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/client'
import { useUIStore } from './ui'
import type { Account, AccountFilters } from '@/types'

export const useAccountsStore = defineStore('accounts', () => {
  const accounts = ref<Account[]>([])
  const loading = ref(false)
  const uiStore = useUIStore()
  const pagination = ref({ page: 1, pageSize: 20, total: 0 })
  const filters = ref<AccountFilters>({
    q: '',
    type: '',
    currency: '',
    archived: '',
    sort: 'name',
    order: 'ASC'
  })

  const fetchAccounts = async (page = 1) => {
    loading.value = true
    try {
      const params = new URLSearchParams({
        page: String(page),
        pageSize: String(pagination.value.pageSize),
        sort: filters.value.sort,
        order: filters.value.order,
        ...Object.fromEntries(
          Object.entries(filters.value).filter(([, v]) => v !== '')
        )
      })
      const response = await api.get(`/api/v1/accounts?${params}`)
      accounts.value = response.data
      pagination.value = {
        ...pagination.value,
        page,
        total: accounts.value.length
      }
    } catch (e) {
      uiStore.error('Failed to fetch accounts')
    } finally {
      loading.value = false
    }
  }

  const createAccount = async (data: Partial<Account>) => {
    try {
      const response = await api.post('/api/v1/accounts', data)
      await fetchAccounts(pagination.value.page)
      return response.data.id
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to create account')
      throw e
    }
  }

  const updateAccount = async (id: string, data: Partial<Account>) => {
    try {
      await api.put(`/api/v1/accounts/${id}`, data)
      await fetchAccounts(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to update account')
      throw e
    }
  }

  const deleteAccount = async (id: string) => {
    try {
      await api.delete(`/api/v1/accounts/${id}`)
      await fetchAccounts(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to delete account')
      throw e
    }
  }

  return {
    accounts,
    loading,
    pagination,
    filters,
    fetchAccounts,
    createAccount,
    updateAccount,
    deleteAccount
  }
})