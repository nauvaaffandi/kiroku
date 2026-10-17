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
                localStorage.setItem('kiroku_backend_url', data.result)
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

      const envUrl = import.meta.env.VITE_API_BASE_URL
      if (envUrl && !envUrl.includes('localhost')) {
        url.value = envUrl
        localStorage.setItem('kiroku_backend_url', envUrl)
        return envUrl
      }

      const cached = localStorage.getItem('kiroku_backend_url')
      if (cached) {
        url.value = cached
        return cached
      }

      const fallback = 'http://localhost:8080/api/v1'
      url.value = fallback
      return fallback
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Unknown error'
      error.value = msg
      console.warn('Failed to fetch backend URL:', e)
      
      const cached = localStorage.getItem('kiroku_backend_url')
      if (cached) {
        url.value = cached
        return cached
      }
      
      const fallback = 'http://localhost:8080/api/v1'
      url.value = fallback
      return fallback
    } finally {
      loading.value = false
      initialized.value = true
    }
  }

  const setUrl = (newUrl: string) => {
    url.value = newUrl
    localStorage.setItem('kiroku_backend_url', newUrl)
  }

  const clearCache = () => {
    url.value = null
    localStorage.removeItem('kiroku_backend_url')
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