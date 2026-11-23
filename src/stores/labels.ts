import { defineStore } from 'pinia'
import { ref } from 'vue'
import { labelRepository } from '@/repository'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'
import type { Label, CreateLabelRequest, UpdateLabelRequest } from '@/types'

export const useLabelsStore = defineStore('labels', () => {
  const labels = ref<Label[]>([])
  const loading = ref(false)
  const { t } = useI18n()

  const fetchLabels = async () => {
    loading.value = true
    try {
      const rows = await labelRepository.findAll()
      labels.value = rows
    } catch (e: any) {
      const uiStore = useUIStore()
      uiStore.error(e?.message || t('labels.fetchError'))
      throw e
    } finally {
      loading.value = false
    }
  }

  const createLabel = async (data: CreateLabelRequest) => {
    try {
      const id = crypto.randomUUID()
      const label: Label = { id, ...data } as Label
      await labelRepository.create(label)
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
      const current = await labelRepository.findById(id)
      if (!current) throw new Error('Label not found')
      const updated: Label = { ...current, ...data } as Label
      await labelRepository.update(id, updated)
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
      await labelRepository.delete(id)
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
        const current = await labelRepository.findById(item.id)
        if (current) await labelRepository.update(item.id, { ...current, sort_order: item.sort_order } as Label)
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
