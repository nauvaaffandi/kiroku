import { defineStore } from 'pinia'
import { ref } from 'vue'
import { localDb } from '@/storage/repository'
import { useTransactionsStore } from './transactions'
import { useUIStore } from './ui'
import type { ApplyTemplateResponse, Template } from '@/types'

export const useTemplatesStore = defineStore('templates', () => {
  const templates = ref<Template[]>([])
  const loading = ref(false)
  const uiStore = useUIStore()

  const fetchTemplates = async () => {
    loading.value = true
    try {
      templates.value = await localDb.list<Template>('template')
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to load templates')
    } finally {
      loading.value = false
    }
  }

  const createTemplate = async (data: Partial<Template>) => {
    try {
      const id = crypto.randomUUID()
      const now = new Date().toISOString()
      const template: Template = {
        id,
        created_at: now,
        updated_at: now,
        ...data
      } as Template
      await localDb.put('template', id, template)
      await fetchTemplates()
      return id
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to create template')
      throw e
    }
  }

  const updateTemplate = async (id: string, data: Partial<Template>) => {
    try {
      const current = await localDb.get<Template>('template', id)
      if (!current) throw new Error('Template not found')
      const updated: Template = { ...current, ...data, updated_at: new Date().toISOString() } as Template
      await localDb.put('template', id, updated)
      await fetchTemplates()
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to update template')
      throw e
    }
  }

  const deleteTemplate = async (id: string) => {
    try {
      await localDb.remove('template', id)
      await fetchTemplates()
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to delete template')
      throw e
    }
  }

  const applyTemplate = async (id: string): Promise<ApplyTemplateResponse | null> => {
    try {
      const template = await localDb.get<Template>('template', id)
      if (!template) throw new Error('Template not found')
      const transactionsStore = useTransactionsStore()
      const transactionId = crypto.randomUUID()
      const now = new Date().toISOString()
      const transaction = {
        id: transactionId,
        type: template.type,
        date_time: now,
        amount: template.amount ?? 0,
        currency: template.currency ?? null,
        from_account_id: template.from_account_id,
        to_account_id: template.to_account_id,
        category_id: template.category_id,
        labels: template.labels,
        note: template.note,
        payee_payer: template.payee_payer,
        payment_type: template.payment_type,
        warranty_until: template.warranty_until,
        status: template.status,
        place_name: template.place_name,
        place_lat: template.place_lat,
        place_lng: template.place_lng,
        attachments: template.attachments,
        template_id: id,
        created_at: now,
        updated_at: now
      }
      await localDb.put('transaction', transactionId, transaction)
      await transactionsStore.fetchTransactions()
      return { ...transaction, template_id: id } as ApplyTemplateResponse
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to apply template')
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