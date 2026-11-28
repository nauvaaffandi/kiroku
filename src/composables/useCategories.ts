import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { categoryRepository } from '@/repository'

export function useCategories() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'category'

  const create = async (data: unknown) => {
    try {
      console.log('[useCategories] create start', data)
      const id = await categoryRepository.create(data as any)
      console.log('[useCategories] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useCategories] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useCategories] update start', { id, data })
      await categoryRepository.update(id, data as any)
      console.log('[useCategories] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useCategories] update error', msg)
      throw err
    }
  }

  const deleteCategory = async (id: string) => {
    try {
      console.log('[useCategories] delete start', { id })
      await categoryRepository.delete(id)
      console.log('[useCategories] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useCategories] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useCategories] findById start', { id })
      const result = await categoryRepository.findById(id)
      console.log('[useCategories] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useCategories] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useCategories] findAll start', { filters })
      const result = await categoryRepository.findAll(filters)
      console.log('[useCategories] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useCategories] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useCategories] count start', { filters })
      const result = await categoryRepository.count(filters)
      console.log('[useCategories] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useCategories] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteCategory, findById, findAll, count }
}
