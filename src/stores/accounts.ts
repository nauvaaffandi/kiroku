import { defineStore } from 'pinia'
import { ref } from 'vue'
import { localDb } from '@/storage/repository'
import { useUIStore } from './ui'
import { useI18n } from '@/locales'
import type { Account, AccountFilters } from '@/types'

export const useAccountsStore = defineStore('accounts', () => {
  const accounts = ref<Account[]>([])
  const loading = ref(false)
  const uiStore = useUIStore()
  const { t } = useI18n()
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
      const rows = await localDb.list<Account>('account')
      let list = rows.filter((a) => !a.archived)
      if (filters.value.q) {
        const q = filters.value.q.toLowerCase()
        list = list.filter((a) => a.name.toLowerCase().includes(q))
      }
      if (filters.value.type) list = list.filter((a) => a.type === filters.value.type)
      if (filters.value.currency) list = list.filter((a) => a.currency === filters.value.currency)
      if (filters.value.archived === 'true') list = rows
      list.sort((a, b) => {
        const d = filters.value.order === 'DESC' ? -1 : 1
        return a[filters.value.sort as keyof Account] < b[filters.value.sort as keyof Account] ? -d : d
      })
      accounts.value = list
      pagination.value = { ...pagination.value, page, total: list.length }
    } catch (e) {
      uiStore.error(t('accounts.fetchError'))
    } finally {
      loading.value = false
    }
  }

  const createAccount = async (data: Partial<Account>) => {
    try {
      const id = crypto.randomUUID()
      await localDb.put('account', id, { ...data, id, archived: false } as Account)
      await fetchAccounts(pagination.value.page)
      return id
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to create account')
      throw e
    }
  }

  const updateAccount = async (id: string, data: Partial<Account>) => {
    try {
      const current = await localDb.get<Account>('account', id)
      if (!current) throw new Error('Account not found')
      await localDb.put('account', id, { ...current, ...data } as Account)
      await fetchAccounts(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to update account')
      throw e
    }
  }

  const deleteAccount = async (id: string) => {
    try {
      await localDb.remove('account', id)
      await fetchAccounts(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to delete account')
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