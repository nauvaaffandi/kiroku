import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { labelRepository } from '@/repository'

export function useLabels() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'label'

  const create = async (data: unknown) => {
    try {
      console.log('[useLabels] create start', data)
      const id = await labelRepository.create(data as any)
      console.log('[useLabels] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useLabels] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useLabels] update start', { id, data })
      await labelRepository.update(id, data as any)
      console.log('[useLabels] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useLabels] update error', msg)
      throw err
    }
  }

  const deleteLabel = async (id: string) => {
    try {
      console.log('[useLabels] delete start', { id })
      await labelRepository.delete(id)
      console.log('[useLabels] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useLabels] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useLabels] findById start', { id })
      const result = await labelRepository.findById(id)
      console.log('[useLabels] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useLabels] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useLabels] findAll start', { filters })
      const result = await labelRepository.findAll(filters)
      console.log('[useLabels] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useLabels] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useLabels] count start', { filters })
      const result = await labelRepository.count(filters)
      console.log('[useLabels] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useLabels] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteLabel, findById, findAll, count }
}
