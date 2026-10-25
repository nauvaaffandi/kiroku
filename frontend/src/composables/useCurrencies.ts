import { ref, onMounted } from 'vue'
import { api } from '@/api/client'

export function useCurrencies() {
  const currencies = ref<string[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchCurrencies = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await api.get('/api/v1/settings/currencies')
      currencies.value = Array.isArray(res.data) ? res.data : ['USD']
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to load currencies'
      currencies.value = ['USD']
    } finally {
      loading.value = false
    }
  }

  const addCurrency = async (code: string) => {
    loading.value = true
    try {
      const res = await api.post('/api/v1/settings/currencies', { code })
      currencies.value = Array.isArray(res.data) ? res.data : currencies.value
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to add currency'
      return false
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchCurrencies()
  })

  return { currencies, loading, error, fetchCurrencies, addCurrency }
}
