import { defineStore } from 'pinia'
import { ref } from 'vue'
import { localDb } from '@/storage/repository'
import { useI18n } from '@/locales'
import type { NotificationItem } from '@/types'

const NOTIFICATION_KIND = 'notification'

export const useNotificationStore = defineStore('notifications', () => {
  const notifications = ref<NotificationItem[]>([])
  const unreadCount = ref<number>(0)
  const loading = ref<boolean>(false)
  const error = ref<string | null>(null)
  const { t } = useI18n()

  const loadNotifications = async (unreadOnly = false) => {
    loading.value = true
    error.value = null
    try {
      const rows = await localDb.list<NotificationItem>(NOTIFICATION_KIND)
      let list = rows.filter((n) => !n.deleted_at)
      if (unreadOnly) list = list.filter((n) => !n.read_at)
      notifications.value = list
      unreadCount.value = list.filter((n) => !n.read_at).length
    } catch (e: any) {
      error.value = e?.message || t('notifications.fetchError')
    } finally {
      loading.value = false
    }
  }

  const loadUnreadCount = async () => {
    try {
      const rows = await localDb.list<NotificationItem>(NOTIFICATION_KIND)
      unreadCount.value = rows.filter((n) => !n.deleted_at && !n.read_at).length
    } catch {
      // Silently fail — badge is best-effort
    }
  }

  const markAsRead = async (id: string) => {
    try {
      const item = await localDb.get<NotificationItem>(NOTIFICATION_KIND, id)
      if (item) {
        const now = new Date().toISOString()
        await localDb.put(NOTIFICATION_KIND, id, { ...item, read_at: now, updated_at: now })
      }
      await loadNotifications()
    } catch (e: any) {
      error.value = e?.message || 'Failed to mark as read'
    }
  }

  const markAllAsRead = async () => {
    try {
      const now = new Date().toISOString()
      const rows = await localDb.list<NotificationItem>(NOTIFICATION_KIND)
      for (const n of rows) {
        if (!n.read_at) {
          await localDb.put(NOTIFICATION_KIND, n.id, { ...n, read_at: now, updated_at: now })
        }
      }
      await loadNotifications()
    } catch (e: any) {
      error.value = e?.message || 'Failed to mark all as read'
    }
  }

  const remove = async (id: string) => {
    try {
      const now = new Date().toISOString()
      const item = await localDb.get<NotificationItem>(NOTIFICATION_KIND, id)
      if (item) await localDb.put(NOTIFICATION_KIND, id, { ...item, deleted_at: now, updated_at: now })
      await loadNotifications()
    } catch (e: any) {
      error.value = e?.message || 'Failed to delete notification'
    }
  }

  const clearRead = async () => {
    try {
      const now = new Date().toISOString()
      const rows = await localDb.list<NotificationItem>(NOTIFICATION_KIND)
      for (const n of rows) {
        if (n.read_at) await localDb.put(NOTIFICATION_KIND, n.id, { ...n, deleted_at: now, updated_at: now })
      }
      await loadNotifications()
    } catch (e: any) {
      error.value = e?.message || 'Failed to clear read notifications'
    }
  }

  const initPolling = () => {
    // Offline mode: no server polling
  }

  const stopPolling = () => {
    // Offline mode: nothing to stop
  }

  return {
    notifications,
    unreadCount,
    loading,
    error,
    loadNotifications,
    loadUnreadCount,
    markAsRead,
    markAllAsRead,
    remove,
    clearRead,
    initPolling,
    stopPolling
  }
})