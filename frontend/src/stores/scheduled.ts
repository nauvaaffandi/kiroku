import { defineStore } from 'pinia'
import { ref } from 'vue'
import { localDb } from '@/storage/repository'
import { useUIStore } from './ui'
import type { ScheduledTransaction, ScheduledFilters } from '@/types'

export const useScheduledStore = defineStore('scheduled', () => {
  const scheduled = ref<ScheduledTransaction[]>([])
  const loading = ref(false)
  const uiStore = useUIStore()
  const pagination = ref({ page: 1, pageSize: 20, total: 0 })
  const filters = ref<ScheduledFilters>({
    q: '',
    type: '',
    frequency: '',
    active: '',
    sort: 'next_run',
    order: 'ASC'
  })

  const fetchScheduled = async (page = 1) => {
    loading.value = true
    try {
      const rows = await localDb.list<ScheduledTransaction>('scheduled')
      let list = [...rows]
      if (filters.value.q) {
        const q = filters.value.q.toLowerCase()
        list = list.filter((s) => s.name.toLowerCase().includes(q))
      }
      if (filters.value.type) list = list.filter((s) => s.type === filters.value.type)
      if (filters.value.frequency) list = list.filter((s) => s.frequency === filters.value.frequency)
      if (filters.value.active === 'true') list = list.filter((s) => s.active)
      list.sort((a, b) => {
        const d = filters.value.order === 'DESC' ? -1 : 1
        const key = filters.value.sort as keyof ScheduledTransaction
        const av = String(a[key] ?? '')
        const bv = String(b[key] ?? '')
        return av < bv ? -d : av > bv ? d : 0
      })
      scheduled.value = list
      pagination.value = { ...pagination.value, page, total: list.length }
    } catch (e) {
      uiStore.error('Failed to load scheduled transactions')
    } finally {
      loading.value = false
    }
  }

  const createScheduled = async (data: Partial<ScheduledTransaction>) => {
    try {
      const id = crypto.randomUUID()
      const now = new Date().toISOString()
      const row: ScheduledTransaction = {
        id,
        name: '',
        type: 'expense',
        frequency: 'monthly',
        start_date: now,
        end_date: null,
        cron_expression: null,
        amount: 0,
        currency: '',
        from_account_id: null,
        to_account_id: null,
        category_id: null,
        description: '',
        active: true,
        next_run: null,
        occurrences_done: 0,
        created_at: now,
        updated_at: now,
        ...data
      } as ScheduledTransaction
      await localDb.put('scheduled', id, row)
      await fetchScheduled(pagination.value.page)
      return id
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to create scheduled transaction')
      throw e
    }
  }

  const updateScheduled = async (id: string, data: Partial<ScheduledTransaction>) => {
    try {
      const current = await localDb.get<ScheduledTransaction>('scheduled', id)
      if (!current) throw new Error('Scheduled transaction not found')
      await localDb.put('scheduled', id, { ...current, ...data, updated_at: new Date().toISOString() } as ScheduledTransaction)
      await fetchScheduled(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to update scheduled transaction')
      throw e
    }
  }

  const deleteScheduled = async (id: string) => {
    try {
      await localDb.remove('scheduled', id)
      await fetchScheduled(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to delete scheduled transaction')
      throw e
    }
  }

  const toggleScheduled = async (id: string) => {
    try {
      const current = await localDb.get<ScheduledTransaction>('scheduled', id)
      if (current) {
        await localDb.put('scheduled', id, { ...current, active: !current.active, updated_at: new Date().toISOString() } as ScheduledTransaction)
      }
      await fetchScheduled(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to toggle scheduled transaction')
      throw e
    }
  }

  return {
    scheduled,
    loading,
    pagination,
    filters,
    fetchScheduled,
    createScheduled,
    updateScheduled,
    deleteScheduled,
    toggleScheduled
  }
})