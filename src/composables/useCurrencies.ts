import { ref, onMounted } from 'vue'
import { useI18n } from '@/locales'
import { currencyRepository } from '@/repository'

export function useCurrencies() {
	const currencies = ref<string[]>([])
	const loading = ref(false)
	const error = ref<string | null>(null)
	const { t } = useI18n()

	const fetchCurrencies = async () => {
		loading.value = true
		error.value = null
		try {
			const picked = await currencyRepository.findPicked()
			currencies.value = picked.map((c) => c.code)
			console.log('[useCurrencies] fetchCurrencies', {
				count: currencies.value.length,
				codes: currencies.value,
			})
		} catch (e: any) {
			error.value = e.message || t('settings.fetchCurrenciesError')
			currencies.value = []
			console.error('[useCurrencies] fetchCurrencies error', error.value)
		} finally {
			loading.value = false
		}
	}

	const addCurrency = async (code: string): Promise<boolean> => {
		try {
			await currencyRepository.pick(code)
			await fetchCurrencies()
			console.log('[useCurrencies] addCurrency', { code, success: true })
			return true
		} catch (e: any) {
			error.value = e.message || t('settings.fetchCurrenciesError')
			console.error('[useCurrencies] addCurrency error', {
				code,
				error: error.value,
			})
			return false
		}
	}

	const removeCurrency = async (code: string): Promise<boolean> => {
		try {
			await currencyRepository.unpick(code)
			await fetchCurrencies()
			console.log('[useCurrencies] removeCurrency', {
				code,
				success: true,
			})
			return true
		} catch (e: any) {
			error.value = e.message || t('settings.fetchCurrenciesError')
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
