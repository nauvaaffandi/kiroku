import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  PushSubscription,
  PushPermissionState,
  PushSubscriptionStatus,
  PushNotificationPayload
} from '@/types'
import { pushSubscribe, pushUnsubscribe, pushStatus, pushSend, pushGetVAPIDKey } from '@/api/push'

export const usePushStore = defineStore('push', () => {
  const subscription = ref<PushSubscription | null>(null)
  const permission = ref<PushPermissionState>({
    granted: false,
    denied: false,
    prompt: true,
    supported: false,
  })
  const loading = ref(false)
  const error = ref<string | null>(null)
  const initializing = ref(false)

  const vapidPublicKey = computed<string>(() => {
    const key = import.meta.env.VITE_VAPID_PUBLIC_KEY
    return key || ''
  })

  const isSubscribed = computed(() => subscription.value !== null)
  const isSupported = computed(() => permission.value.supported)
  const canSubscribe = computed(() => isSupported.value && !permission.value.denied && !isSubscribed.value)

  // Check if Push API, Notifications API, and Service Worker are supported
  const checkSupport = () => {
    const supported =
      typeof window !== 'undefined' &&
      'PushManager' in window &&
      'Notification' in window &&
      'serviceWorker' in navigator

    permission.value.supported = supported

    if (typeof window !== 'undefined' && 'Notification' in window) {
      const state = Notification.permission
      permission.value.granted = state === 'granted'
      permission.value.denied = state === 'denied'
      permission.value.prompt = state === 'default'
    }
  }

  const getVapidKey = async (): Promise<Uint8Array | null> => {
    // Fetch from backend first (VAPID keys auto-generated at runtime)
    try {
      const config = await pushGetVAPIDKey()
      if (config.publicKey) {
        return urlBase64ToUint8Array(config.publicKey)
      }
    } catch {
      // Fall back to env var if backend unavailable
    }

    const key = vapidPublicKey.value
    if (!key) return null
    try {
      return urlBase64ToUint8Array(key)
    } catch {
      return null
    }
  }

  const initializePush = async (): Promise<void> => {
    initializing.value = true
    try {
      if (!isSupported.value) {
        checkSupport()
        if (!isSupported.value) return
      }
      const registration = await navigator.serviceWorker.ready
      const existing = await registration.pushManager.getSubscription()
      if (existing) {
        subscription.value = {
          endpoint: existing.endpoint,
          keys: {
            p256dh: arrayBufferToBase64(existing.getKey('p256dh') as ArrayBuffer),
            auth: arrayBufferToBase64(existing.getKey('auth') as ArrayBuffer),
          },
          expirationTime: existing.expirationTime,
        }
      }
    } catch {
      // Ignore init errors
    } finally {
      initializing.value = false
    }
  }

  const subscribe = async (): Promise<boolean> => {
    loading.value = true
    error.value = null

    try {
      if (!isSupported.value) {
        throw new Error('Web Push not supported in this browser')
      }
      if (permission.value.denied) {
        throw new Error('Notification permission denied')
      }

      // Request notification permission
      const perm = await Notification.requestPermission()
      permission.value.granted = perm === 'granted'
      permission.value.denied = perm === 'denied'
      permission.value.prompt = perm === 'default'

      if (perm !== 'granted') {
        throw new Error('Notification permission not granted')
      }

      // Get service worker registration
      const registration = await navigator.serviceWorker.ready

      // Check if already subscribed
      const existing = await registration.pushManager.getSubscription()

      // Subscribe with VAPID key
      const vapidKey = await getVapidKey()
      if (!vapidKey) {
        throw new Error('VAPID public key not configured')
      }

      const newSubscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: vapidKey as unknown as ArrayBuffer,
      })

      const subData = {
        endpoint: newSubscription.endpoint,
        keys: {
          p256dh: arrayBufferToBase64(newSubscription.getKey('p256dh') as ArrayBuffer),
          auth: arrayBufferToBase64(newSubscription.getKey('auth') as ArrayBuffer),
        },
        expirationTime: newSubscription.expirationTime,
      }

      // Sync with backend first, then set state on success
      await pushSubscribe({
        endpoint: subData.endpoint,
        keys: {
          p256dh: subData.keys.p256dh,
          auth: subData.keys.auth,
        },
        user_agent: navigator.userAgent,
      })

      // If already subscribed in browser, the pushSubscribe confirms backend sync;
      // otherwise subscription was just created in browser, so both paths set state here
      if (existing) {
        subscription.value = subData
      } else {
        subscription.value = subData
      }

      return true
    } catch (e: any) {
      error.value = e?.message || 'Failed to subscribe to push notifications'
      return false
    } finally {
      loading.value = false
    }
  }

  const unsubscribe = async (): Promise<boolean> => {
    loading.value = true
    error.value = null

    try {
      const registration = await navigator.serviceWorker.ready
      const existing = await registration.pushManager.getSubscription()

      if (existing) {
        // Browser unsubscribe first, then backend cleanup
        await existing.unsubscribe()
        await pushUnsubscribe(existing.endpoint)
      }

      subscription.value = null
      return true
    } catch (e: any) {
      error.value = e?.message || 'Failed to unsubscribe from push notifications'
      return false
    } finally {
      loading.value = false
    }
  }

  const checkSubscriptionStatus = async (): Promise<PushSubscriptionStatus> => {
    try {
      const status = await pushStatus()
      return status
    } catch {
      return { subscribed: false, endpoint: null }
    }
  }

  const sendTestNotification = async (payload: PushNotificationPayload): Promise<boolean> => {
    try {
      await pushSend(payload)
      return true
    } catch (e: any) {
      error.value = e?.message || 'Failed to send test notification'
      return false
    }
  }

  const clearError = () => {
    error.value = null
  }

  return {
    subscription,
    permission,
    loading,
    error,
    initializing,
    vapidPublicKey,
    isSubscribed,
    isSupported,
    canSubscribe,
    checkSupport,
    initializePush,
    subscribe,
    unsubscribe,
    checkSubscriptionStatus,
    sendTestNotification,
    clearError,
  }
})

// --- helpers ---

function urlBase64ToUint8Array(base64String: string): Uint8Array {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/')
  const rawData = window.atob(base64)
  const outputArray = new Uint8Array(rawData.length)
  for (let i = 0; i < rawData.length; i++) {
    outputArray[i] = rawData.charCodeAt(i)
  }
  return outputArray
}

function arrayBufferToBase64(buffer: ArrayBuffer): string {
  let binary = ''
  const bytes = new Uint8Array(buffer)
  const len = bytes.byteLength
  for (let i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return window.btoa(binary)
}
