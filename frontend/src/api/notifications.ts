import { api } from '@/api/client'
import type { NotificationFilter, NotificationItem } from '@/types'

export const fetchNotifications = async (
  filter?: NotificationFilter
): Promise<NotificationItem[]> => {
  const params: Record<string, string | number | boolean> = {}
  if (filter?.unread !== undefined) params.unread = filter.unread
  if (filter?.limit !== undefined) params.limit = filter.limit
  if (filter?.offset !== undefined) params.offset = filter.offset

  const { data } = await api.get<NotificationItem[]>('/api/v1/notifications', {
    params
  })
  return data
}

export const fetchUnreadCount = async (): Promise<{
  unread_count: number
}> => {
  const { data } = await api.get<{ unread_count: number }>(
    '/api/v1/notifications/unread-count'
  )
  return data
}

export const markNotificationRead = async (
  id: string
): Promise<{ status: string }> => {
  const { data } = await api.patch<{ status: string }>(
    `/api/v1/notifications/${id}/read`
  )
  return data
}

export const markAllNotificationsRead = async (): Promise<{
  status: string
}> => {
  const { data } = await api.patch<{ status: string }>(
    '/api/v1/notifications/read-all'
  )
  return data
}

export const deleteNotification = async (id: string): Promise<{
  status: string
}> => {
  const { data } = await api.delete<{ status: string }>(
    `/api/v1/notifications/${id}`
  )
  return data
}

export const deleteReadNotifications = async (): Promise<{
  status: string
}> => {
  const { data } = await api.delete<{ status: string }>(
    '/api/v1/notifications/read'
  )
  return data
}
