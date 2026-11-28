import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { scheduledRepository } from '@/repository'

export function useScheduled() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'scheduled'

  const create = async (data: unknown) => {
    try {
      console.log('[useScheduled] create start', data)
      const id = await scheduledRepository.create(data as any)
      console.log('[useScheduled] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useScheduled] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useScheduled] update start', { id, data })
      await scheduledRepository.update(id, data as any)
      console.log('[useScheduled] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useScheduled] update error', msg)
      throw err
    }
  }

  const deleteScheduled = async (id: string) => {
    try {
      console.log('[useScheduled] delete start', { id })
      await scheduledRepository.delete(id)
      console.log('[useScheduled] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useScheduled] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useScheduled] findById start', { id })
      const result = await scheduledRepository.findById(id)
      console.log('[useScheduled] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useScheduled] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useScheduled] findAll start', { filters })
      const result = await scheduledRepository.findAll(filters)
      console.log('[useScheduled] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useScheduled] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useScheduled] count start', { filters })
      const result = await scheduledRepository.count(filters)
      console.log('[useScheduled] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useScheduled] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteScheduled, findById, findAll, count }
}
