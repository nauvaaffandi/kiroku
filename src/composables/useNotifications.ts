import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { notificationRepository } from '@/repository'

export function useNotifications() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'notification'

  const create = async (data: unknown) => {
    try {
      console.log('[useNotifications] create start', data)
      const id = await notificationRepository.create(data as any)
      console.log('[useNotifications] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useNotifications] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useNotifications] update start', { id, data })
      await notificationRepository.update(id, data as any)
      console.log('[useNotifications] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useNotifications] update error', msg)
      throw err
    }
  }

  const deleteNotification = async (id: string) => {
    try {
      console.log('[useNotifications] delete start', { id })
      await notificationRepository.delete(id)
      console.log('[useNotifications] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useNotifications] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useNotifications] findById start', { id })
      const result = await notificationRepository.findById(id)
      console.log('[useNotifications] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useNotifications] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useNotifications] findAll start', { filters })
      const result = await notificationRepository.findAll(filters)
      console.log('[useNotifications] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useNotifications] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useNotifications] count start', { filters })
      const result = await notificationRepository.count(filters)
      console.log('[useNotifications] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useNotifications] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteNotification, findById, findAll, count }
}
