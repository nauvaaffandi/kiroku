import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<{ user_id: string }>({ user_id: 'authenticated' })
  const loading = ref(false)
  const error = ref<string | null>(null)
  const isChangingPin = ref(false)

  const isAuthenticated = computed(() => true)

  const init = async () => {
    // No session in local-only offline app; always authenticated.
  }

  const fetchUser = async () => {
    user.value = { user_id: 'authenticated' }
  }

  const login = async (_pin: string) => {
    user.value = { user_id: 'authenticated' }
    return true
  }

  const logout = async () => {
    user.value = { user_id: 'authenticated' }
  }

  const changePin = async (_currentPin: string, _newPin: string) => {
    return true
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