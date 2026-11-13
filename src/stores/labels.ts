import { defineStore } from 'pinia'
import { ref } from 'vue'
import { localDb } from '@/storage/repository'
import { useUIStore } from './ui'
import type { Label, CreateLabelRequest, UpdateLabelRequest } from '@/types'

export const useLabelsStore = defineStore('labels', () => {
  const labels = ref<Label[]>([])
  const loading = ref(false)

  const fetchLabels = async () => {
    loading.value = true
    try {
      const rows = await localDb.list<Label>('label')
      labels.value = rows
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e?.message || 'Failed to load labels')
      throw e
    } finally {
      loading.value = false
    }
  }

  const createLabel = async (data: CreateLabelRequest) => {
    try {
      const id = crypto.randomUUID()
      const label: Label = { id, ...data } as Label
      await localDb.put('label', id, label)
      labels.value.push(label)
      return id
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e?.message || 'Failed to create label')
      throw e
    }
  }

  const updateLabel = async (id: string, data: UpdateLabelRequest) => {
    try {
      const current = await localDb.get<Label>('label', id)
      if (!current) throw new Error('Label not found')
      const updated: Label = { ...current, ...data } as Label
      await localDb.put('label', id, updated)
      const idx = labels.value.findIndex(l => l.id === id)
      if (idx >= 0) labels.value[idx] = updated
      return updated
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e?.message || 'Failed to update label')
      throw e
    }
  }

  const deleteLabel = async (id: string) => {
    try {
      await localDb.remove('label', id)
      labels.value = labels.value.filter(l => l.id !== id)
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e?.message || 'Failed to delete label')
      throw e
    }
  }

  const reorderLabels = async (items: { id: string; sort_order: number }[]) => {
    try {
      for (const item of items) {
        const current = await localDb.get<Label>('label', item.id)
        if (current) await localDb.put('label', item.id, { ...current, sort_order: item.sort_order } as Label)
      }
      await fetchLabels()
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e?.message || 'Failed to reorder labels')
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
