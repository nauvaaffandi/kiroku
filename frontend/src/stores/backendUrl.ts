import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useBackendUrlStore = defineStore('backendUrl', () => {
  const url = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const initialized = ref(false)

  const hasUrl = computed(() => !!url.value)

  const fetchUrl = async (): Promise<string> => {
    if (url.value && !loading.value) {
      return url.value
    }

    loading.value = true
    error.value = null

    try {
      const upstashUrl = import.meta.env.VITE_UPSTASH_URL
      const upstashToken = import.meta.env.VITE_FRONTEND_UPSTASH_REDIS_REST_TOKEN

      if (upstashUrl && upstashToken) {
        for (let attempt = 0; attempt < 5; attempt++) {
          try {
            const response = await fetch(`${upstashUrl}/get/kiroku:backend:url`, {
              headers: { Authorization: `Bearer ${upstashToken}` }
            })
            if (response.ok) {
              const data = await response.json()
              if (data.result) {
                url.value = data.result
                return data.result
              }
            }
          } catch {
            // retry on network error
          }
          if (attempt < 4) {
            await new Promise((r) => setTimeout(r, 1000))
          }
        }
      }

      throw new Error('Failed to fetch backend URL from Upstash after 5 attempts')
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Unknown error'
      error.value = msg
      console.warn('Failed to fetch backend URL:', e)
      throw e
    } finally {
      loading.value = false
      initialized.value = true
    }
  }

  const setUrl = (newUrl: string) => {
    url.value = newUrl
  }

  const clearCache = () => {
    url.value = null
  }

  return {
    url,
    loading,
    error,
    initialized,
    hasUrl,
    fetchUrl,
    setUrl,
    clearCache
  }
})