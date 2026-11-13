import { ref, onMounted } from 'vue'
import { useI18n } from '@/locales'

export function useCurrencies() {
  const currencies = ref<string[]>(['USD'])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const { t } = useI18n()

  const fetchCurrencies = async () => {
    loading.value = true
    error.value = null
    try {
      currencies.value = ['USD']
    } catch (e: any) {
      error.value = e.message || t('settings.fetchCurrenciesError')
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