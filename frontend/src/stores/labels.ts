import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/client'
import { useUIStore } from './ui'
import type { Label, CreateLabelRequest, UpdateLabelRequest } from '@/types'

export const useLabelsStore = defineStore('labels', () => {
  const labels = ref<Label[]>([])
  const loading = ref(false)

  const fetchLabels = async () => {
    loading.value = true
    try {
      const res = await api.get('/api/v1/labels')
      labels.value = Array.isArray(res.data) ? res.data : []
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e.response?.data?.error || 'Failed to fetch labels')
      throw e
    } finally {
      loading.value = false
    }
  }

  const createLabel = async (data: CreateLabelRequest) => {
    try {
      const res = await api.post('/api/v1/labels', data)
      labels.value.push(res.data)
      return res.data.id
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e.response?.data?.error || 'Failed to create label')
      throw e
    }
  }

  const updateLabel = async (id: string, data: UpdateLabelRequest) => {
    try {
      const res = await api.put(`/api/v1/labels/${id}`, data)
      const idx = labels.value.findIndex(l => l.id === id)
      if (idx >= 0) labels.value[idx] = res.data
      return res.data
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e.response?.data?.error || 'Failed to update label')
      throw e
    }
  }

  const deleteLabel = async (id: string) => {
    try {
      await api.delete(`/api/v1/labels/${id}`)
      labels.value = labels.value.filter(l => l.id !== id)
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e.response?.data?.error || 'Failed to delete label')
      throw e
    }
  }

  const reorderLabels = async (items: { id: string; sort_order: number }[]) => {
    try {
      await api.put('/api/v1/labels/reorder', { items })
      await fetchLabels()
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e.response?.data?.error || 'Failed to reorder labels')
      throw e
    }
  }

  return {
    labels,
    loading,
    fetchLabels,
    createLabel,
    updateLabel,
    deleteLabel,
    reorderLabels
  }
})
