import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { budgetRepository } from '@/repository'

export function useBudgets() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'budget'

  const create = async (data: unknown) => {
    try {
      console.log('[useBudgets] create start', data)
      const id = await budgetRepository.create(data as any)
      console.log('[useBudgets] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useBudgets] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useBudgets] update start', { id, data })
      await budgetRepository.update(id, data as any)
      console.log('[useBudgets] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useBudgets] update error', msg)
      throw err
    }
  }

  const deleteBudget = async (id: string) => {
    try {
      console.log('[useBudgets] delete start', { id })
      await budgetRepository.delete(id)
      console.log('[useBudgets] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useBudgets] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useBudgets] findById start', { id })
      const result = await budgetRepository.findById(id)
      console.log('[useBudgets] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useBudgets] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useBudgets] findAll start', { filters })
      const result = await budgetRepository.findAll(filters)
      console.log('[useBudgets] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useBudgets] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useBudgets] count start', { filters })
      const result = await budgetRepository.count(filters)
      console.log('[useBudgets] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useBudgets] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteBudget, findById, findAll, count }
}
