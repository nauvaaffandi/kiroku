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
