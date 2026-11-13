import { api } from '@/api/client'
import type {
  SubscribeRequest,
  PushSubscriptionStatus,
  PushNotificationPayload,
  VAPIDConfig
} from '@/types'

export const pushSubscribe = async (req: SubscribeRequest): Promise<{ status: string; id?: string; endpoint?: string }> => {
  const response = await api.post('/api/v1/push/subscribe', req)
  return response.data
}

export const pushUnsubscribe = async (endpoint: string): Promise<{ status: string }> => {
  const response = await api.post('/api/v1/push/unsubscribe', { endpoint })
  return response.data
}

export const pushStatus = async (): Promise<PushSubscriptionStatus> => {
  const response = await api.get('/api/v1/push/status')
  return response.data
}

export const pushSend = async (payload: PushNotificationPayload): Promise<{ status: string; sent: number }> => {
  const response = await api.post('/api/v1/push/send', payload)
  return response.data
}

export const pushGetVAPIDKey = async (): Promise<VAPIDConfig> => {
  const response = await api.get('/api/v1/push/vapid-public-key')
  return response.data
}
