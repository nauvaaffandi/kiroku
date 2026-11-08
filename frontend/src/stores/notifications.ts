import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { NotificationFilter, NotificationItem } from '@/types'
import {
  fetchNotifications,
  fetchUnreadCount,
  markNotificationRead,
  markAllNotificationsRead,
  deleteNotification,
  deleteReadNotifications
} from '@/api/notifications'

export const useNotificationStore = defineStore('notifications', () => {
  const notifications = ref<NotificationItem[]>([])
  const unreadCount = ref<number>(0)
  const loading = ref<boolean>(false)
  const error = ref<string | null>(null)

  let pollingTimer: ReturnType<typeof setInterval> | null = null

  const loadNotifications = async (unreadOnly = false) => {
    loading.value = true
    error.value = null
    try {
      const filter: NotificationFilter | undefined = unreadOnly
        ? { unread: true }
        : undefined
      const data = await fetchNotifications(filter)
      notifications.value = data
    } catch (e: any) {
      error.value = e?.message || 'Failed to load notifications'
    } finally {
      loading.value = false
    }
  }

  const loadUnreadCount = async () => {
    try {
      const data = await fetchUnreadCount()
      unreadCount.value = data.unread_count
    } catch {
      // Silently fail — badge is best-effort
    }
  }

  const markAsRead = async (id: string) => {
    try {
      await markNotificationRead(id)
      const item = notifications.value.find((n) => n.id === id)
      if (item && !item.read_at) {
        item.read_at = new Date().toISOString()
        if (unreadCount.value > 0) unreadCount.value--
      }
    } catch (e: any) {
      error.value = e?.message || 'Failed to mark as read'
    }
  }

  const markAllAsRead = async () => {
    try {
      await markAllNotificationsRead()
      const now = new Date().toISOString()
      notifications.value.forEach((n) => {
        if (!n.read_at) n.read_at = now
      })
      unreadCount.value = 0
    } catch (e: any) {
      error.value = e?.message || 'Failed to mark all as read'
    }
  }

  const remove = async (id: string) => {
    try {
      const item = notifications.value.find((n) => n.id === id)
      await deleteNotification(id)
      notifications.value = notifications.value.filter((n) => n.id !== id)
      if (item && !item.read_at && unreadCount.value > 0) unreadCount.value--
    } catch (e: any) {
      error.value = e?.message || 'Failed to delete notification'
    }
  }

  const clearRead = async () => {
    try {
      await deleteReadNotifications()
      notifications.value = notifications.value.filter(
        (n) => !n.read_at
      )
    } catch (e: any) {
      error.value = e?.message || 'Failed to clear read notifications'
    }
  }

  const initPolling = (intervalMs = 30000) => {
    stopPolling()
    pollingTimer = setInterval(() => {
      loadUnreadCount()
    }, intervalMs)
  }

  const stopPolling = () => {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
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