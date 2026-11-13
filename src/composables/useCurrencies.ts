import { ref, onMounted } from 'vue'

export function useCurrencies() {
  const currencies = ref<string[]>(['USD'])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchCurrencies = async () => {
    loading.value = true
    error.value = null
    try {
      currencies.value = ['USD']
    } catch (e: any) {
      error.value = e.message || 'Failed to load currencies'
      currencies.value = ['USD']
    } finally {
      loading.value = false
    }
  }

  const addCurrency = async (_code: string) => {
    return false
  }

  onMounted(() => {
    fetchCurrencies()
  })

  return { currencies, loading, error, fetchCurrencies, addCurrency }
}
