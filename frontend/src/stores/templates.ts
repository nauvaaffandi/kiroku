import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/client'
import { useUIStore } from './ui'
import type { Template } from '@/types'

export const useTemplatesStore = defineStore('templates', () => {
  const templates = ref<Template[]>([])
  const loading = ref(false)
  const uiStore = useUIStore()

  const fetchTemplates = async () => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/templates')
      templates.value = response.data
    } catch (e) {
      uiStore.error('Failed to fetch templates')
    } finally {
      loading.value = false
    }
  }

  const createTemplate = async (data: Partial<Template>) => {
    try {
      const response = await api.post('/api/v1/templates', data)
      await fetchTemplates()
      return response.data.id
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to create template')
      throw e
    }
  }

  const updateTemplate = async (id: string, data: Partial<Template>) => {
    try {
      await api.put(`/api/v1/templates/${id}`, data)
      await fetchTemplates()
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to update template')
      throw e
    }
  }

  const deleteTemplate = async (id: string) => {
    try {
      await api.delete(`/api/v1/templates/${id}`)
      await fetchTemplates()
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to delete template')
      throw e
    }
  }

  const applyTemplate = async (id: string) => {
    try {
      const response = await api.post(`/api/v1/templates/${id}/apply`)
      return response.data
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to apply template')
      throw e
    }
  }

  return {
    templates,
    loading,
    fetchTemplates,
    createTemplate,
    updateTemplate,
    deleteTemplate,
    applyTemplate
  }
})