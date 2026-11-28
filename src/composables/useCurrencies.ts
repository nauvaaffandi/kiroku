import { ref, onMounted } from 'vue'
import { useI18n } from '@/locales'
import { currencyRepository } from '@/repository'
import { useRepositoryProxy } from './useRepositoryProxy'

export function useCurrencies() {
  const currencies = ref<string[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const { t } = useI18n()

  const repo = useRepositoryProxy(currencyRepository, 'currency')

  const fetchCurrencies = async () => {
    loading.value = true
    error.value = null
    try {
      const picked = await repo.findPicked()
      currencies.value = picked.map((c) => c.code)
      console.log('[useCurrencies] fetchCurrencies', {
        count: currencies.value.length,
        codes: currencies.value,
      })
    } catch (e: any) {
      error.value = e.message || t('repository.error.fetch', { entity: 'currency' })
      currencies.value = []
      console.error('[useCurrencies] fetchCurrencies error', error.value)
    } finally {
      loading.value = false
    }
  }

  const addCurrency = async (code: string): Promise<boolean> => {
    try {
      await repo.pick(code)
      await fetchCurrencies()
      console.log('[useCurrencies] addCurrency', { code, success: true })
      return true
    } catch (e: any) {
      error.value = e.message || t('repository.error.save', { entity: 'currency' })
      console.error('[useCurrencies] addCurrency error', {
        code,
        error: error.value,
      })
      return false
    }
  }

  const removeCurrency = async (code: string): Promise<boolean> => {
    try {
      await repo.unpick(code)
      await fetchCurrencies()
      console.log('[useCurrencies] removeCurrency', {
        code,
        success: true,
      })
      return true
    } catch (e: any) {
      error.value = e.message || t('repository.error.delete', { entity: 'currency' })
      console.error('[useCurrencies] removeCurrency error', {
        code,
        error: error.value,
      })
      return false
    }
  }

  onMounted(() => {
    fetchCurrencies()
  })

  return {
    currencies,
    loading,
    error,
    fetchCurrencies,
    addCurrency,
    removeCurrency,
  }
}
