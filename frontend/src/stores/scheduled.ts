import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/client'
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
      const params = new URLSearchParams({
        page: String(page),
        pageSize: String(pagination.value.pageSize),
        sort: filters.value.sort,
        order: filters.value.order,
        ...Object.fromEntries(
          Object.entries(filters.value).filter(([, v]) => v !== '')
        )
      })
      const response = await api.get(`/api/v1/scheduled?${params}`)
      scheduled.value = response.data
      pagination.value = {
        ...pagination.value,
        page,
        total: scheduled.value.length
      }
    } catch (e) {
      uiStore.error('Failed to fetch scheduled transactions')
    } finally {
      loading.value = false
    }
  }

  const createScheduled = async (data: Partial<ScheduledTransaction>) => {
    try {
      const response = await api.post('/api/v1/scheduled', data)
      await fetchScheduled(pagination.value.page)
      return response.data.id
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to create scheduled transaction')
      throw e
    }
  }

  const updateScheduled = async (id: string, data: Partial<ScheduledTransaction>) => {
    try {
      await api.put(`/api/v1/scheduled/${id}`, data)
      await fetchScheduled(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to update scheduled transaction')
      throw e
    }
  }

  const deleteScheduled = async (id: string) => {
    try {
      await api.delete(`/api/v1/scheduled/${id}`)
      await fetchScheduled(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to delete scheduled transaction')
      throw e
    }
  }

  const toggleScheduled = async (id: string) => {
    try {
      await api.post(`/api/v1/scheduled/${id}/toggle`)
      await fetchScheduled(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to toggle scheduled transaction')
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