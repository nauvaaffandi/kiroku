export interface PushSubscription {
  endpoint: string
  keys: {
    p256dh: string
    auth: string
  }
  expirationTime: number | null
}

export interface PushPermissionState {
  granted: boolean
  denied: boolean
  prompt: boolean
  supported: boolean
}

export interface PushSubscriptionStatus {
  subscribed: boolean
  endpoint: string | null
}

export interface PushNotificationPayload {
  title: string
  body?: string
  url?: string
}

export interface SubscribeRequest {
  endpoint: string
  keys: {
    p256dh: string
    auth: string
  }
  user_agent?: string
}

export interface VAPIDConfig {
  publicKey: string
}

export type PushSubscriptionAction = 'subscribe' | 'unsubscribe'

// Notification types for in-app notifications
export interface NotificationItem {
  id: string
  user_id: string
  type: string
  title: string
  body: string
  data?: Record<string, any>
  action_url?: string
  read_at?: string | null
  deleted_at?: string | null
  idempotency_key?: string
  created_at: string
  updated_at: string
}

export interface NotificationFilter {
  unread?: boolean
  limit?: number
  offset?: number
}
