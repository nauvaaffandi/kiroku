import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  PushSubscription,
  PushPermissionState,
  PushSubscriptionStatus,
  PushNotificationPayload
} from '@/types'

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

  const vapidPublicKey = computed<string>(() => '')

  const isSubscribed = computed(() => false)
  const isSupported = computed(() => false)
  const canSubscribe = computed(() => false)

  const checkSupport = () => {
    permission.value.supported = false
    permission.value.granted = false
    permission.value.denied = false
    permission.value.prompt = false
  }

  const initializePush = async (): Promise<void> => {
    initializing.value = false
  }

  const subscribe = async (): Promise<boolean> => {
    error.value = 'Push notifications not available offline'
    return false
  }

  const unsubscribe = async (): Promise<boolean> => {
    error.value = 'Push notifications not available offline'
    return false
  }

  const checkSubscriptionStatus = async (): Promise<PushSubscriptionStatus> => ({ subscribed: false, endpoint: null })

  const sendTestNotification = async (_payload: PushNotificationPayload): Promise<boolean> => {
    error.value = 'Push notifications not available offline'
    return false
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
