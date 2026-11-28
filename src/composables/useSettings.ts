import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { settingRepository } from '@/repository'

export function useSettings() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'setting'

  const create = async (data: unknown) => {
    try {
      console.log('[useSettings] create start', data)
      const id = await settingRepository.create(data as any)
      console.log('[useSettings] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useSettings] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useSettings] update start', { id, data })
      await settingRepository.update(id, data as any)
      console.log('[useSettings] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useSettings] update error', msg)
      throw err
    }
  }

  const deleteSetting = async (id: string) => {
    try {
      console.log('[useSettings] delete start', { id })
      await settingRepository.delete(id)
      console.log('[useSettings] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useSettings] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useSettings] findById start', { id })
      const result = await settingRepository.findById(id)
      console.log('[useSettings] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useSettings] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useSettings] findAll start', { filters })
      const result = await settingRepository.findAll(filters)
      console.log('[useSettings] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useSettings] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useSettings] count start', { filters })
      const result = await settingRepository.count(filters)
      console.log('[useSettings] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useSettings] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteSetting, findById, findAll, count }
}
