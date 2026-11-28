import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { goalRepository } from '@/repository'

export function useGoals() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'goal'

  const create = async (data: unknown) => {
    try {
      console.log('[useGoals] create start', data)
      const id = await goalRepository.create(data as any)
      console.log('[useGoals] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useGoals] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useGoals] update start', { id, data })
      await goalRepository.update(id, data as any)
      console.log('[useGoals] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useGoals] update error', msg)
      throw err
    }
  }

  const deleteGoal = async (id: string) => {
    try {
      console.log('[useGoals] delete start', { id })
      await goalRepository.delete(id)
      console.log('[useGoals] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useGoals] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useGoals] findById start', { id })
      const result = await goalRepository.findById(id)
      console.log('[useGoals] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useGoals] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useGoals] findAll start', { filters })
      const result = await goalRepository.findAll(filters)
      console.log('[useGoals] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useGoals] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useGoals] count start', { filters })
      const result = await goalRepository.count(filters)
      console.log('[useGoals] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useGoals] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteGoal, findById, findAll, count }
}
