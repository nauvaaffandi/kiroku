import axios from 'axios'

export const api = axios.create({
  baseURL: '',
  withCredentials: true,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

let resolvePromise: ((url: string) => void) | null = null
let initPromise: Promise<string> | null = null

const setBaseUrl = (url: string) => {
  api.defaults.baseURL = url
  if (resolvePromise) {
    resolvePromise(url)
    resolvePromise = null
  }
}

const waitForBaseUrl = (): Promise<string> => {
  if (api.defaults.baseURL) {
    return Promise.resolve(api.defaults.baseURL)
  }
  if (!initPromise) {
    initPromise = new Promise((resolve) => {
      resolvePromise = resolve
    })
  }
  return initPromise
}

api.interceptors.request.use(async (config) => {
  await waitForBaseUrl()
  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      try {
        const { useAuthStore } = await import('@/stores/auth')
        const authStore = useAuthStore()
        authStore.logout()
      } catch {
        // logout failed
      }
    }
    return Promise.reject(error)
  }
)

export const initApiClient = async () => {
  const { useBackendUrlStore } = await import('@/stores/backendUrl')
  const backendUrlStore = useBackendUrlStore()
  const url = await backendUrlStore.fetchUrl()
  setBaseUrl(url)
}

export const refreshBackendUrl = async () => {
  const { useBackendUrlStore } = await import('@/stores/backendUrl')
  const backendUrlStore = useBackendUrlStore()
  backendUrlStore.clearCache()
  const url = await backendUrlStore.fetchUrl()
  setBaseUrl(url)
}