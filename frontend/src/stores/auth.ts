import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<{ user_id: string } | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const isChangingPin = ref(false)

  const isAuthenticated = computed(() => !!user.value)

  const init = async () => {
    try {
      await fetchUser()
    } catch {
      // unauthenticated -> router redirects to /login
    }
  }

  const fetchUser = async () => {
    const response = await api.get('/api/v1/auth/me')
    user.value = response.data
  }

  const login = async (pin: string) => {
    loading.value = true
    error.value = null
    try {
      await api.post('/api/v1/auth/login', { pin })
      await fetchUser()
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Invalid PIN'
      return false
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    try {
      await api.post('/api/v1/auth/logout')
    } finally {
      user.value = null
    }
  }

  const changePin = async (currentPin: string, newPin: string) => {
    isChangingPin.value = true
    error.value = null
    try {
      await api.post('/api/v1/auth/pin', {
        current_pin: currentPin,
        new_pin: newPin
      })
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to change PIN'
      return false
    } finally {
      isChangingPin.value = false
    }
  }

  return {
    user,
    loading,
    error,
    isChangingPin,
    isAuthenticated,
    init,
    fetchUser,
    login,
    logout,
    changePin
  }
})