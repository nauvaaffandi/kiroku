import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<{ user_id: string } | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const showPasswordOverlay = ref(false)
  const isChangingPassword = ref(false)

  const isAuthenticated = computed(() => !!user.value)

const init = async () => {
  try {
    await fetchUser()
    showPasswordOverlay.value = false
  } catch {
    showPasswordOverlay.value = true
  }
}

  const fetchUser = async () => {
    const response = await api.get('/api/v1/auth/me')
    user.value = response.data
  }

  const login = async (password: string) => {
    loading.value = true
    error.value = null
    try {
      await api.post('/api/v1/auth/login', { password })
      await fetchUser()
      showPasswordOverlay.value = false
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Invalid password'
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
      showPasswordOverlay.value = true
    }
  }

  const changePassword = async (oldPassword: string, newPassword: string) => {
    isChangingPassword.value = true
    error.value = null
    try {
      await api.post('/api/v1/auth/password', {
        old_password: oldPassword,
        new_password: newPassword
      })
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to change password'
      return false
    } finally {
      isChangingPassword.value = false
    }
  }

  return {
    user,
    loading,
    error,
    showPasswordOverlay,
    isChangingPassword,
    isAuthenticated,
    init,
    fetchUser,
    login,
    logout,
    changePassword
  }
})