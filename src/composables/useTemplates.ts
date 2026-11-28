import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { templateRepository } from '@/repository'

export function useTemplates() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'template'

  const create = async (data: unknown) => {
    try {
      console.log('[useTemplates] create start', data)
      const id = await templateRepository.create(data as any)
      console.log('[useTemplates] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useTemplates] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useTemplates] update start', { id, data })
      await templateRepository.update(id, data as any)
      console.log('[useTemplates] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useTemplates] update error', msg)
      throw err
    }
  }

  const deleteTemplate = async (id: string) => {
    try {
      console.log('[useTemplates] delete start', { id })
      await templateRepository.delete(id)
      console.log('[useTemplates] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useTemplates] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useTemplates] findById start', { id })
      const result = await templateRepository.findById(id)
      console.log('[useTemplates] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useTemplates] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useTemplates] findAll start', { filters })
      const result = await templateRepository.findAll(filters)
      console.log('[useTemplates] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useTemplates] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useTemplates] count start', { filters })
      const result = await templateRepository.count(filters)
      console.log('[useTemplates] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useTemplates] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteTemplate, findById, findAll, count }
}
