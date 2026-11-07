import { ref, computed } from 'vue'
import { useSettingsStore } from '@/stores/settings'

const FRANKFURTER_BASE = 'https://api.frankfurter.app/latest'
const CACHE_TTL_MS = 60 * 60 * 1000 // 1 jam

interface RateState {
  rate: number | null // rate dari baseCurrency ke USD
  loading: boolean
  error: string | null
  lastFetched: number | null
}

const state = ref<RateState>({
  rate: null,
  loading: false,
  error: null,
  lastFetched: null
})

function isStale(): boolean {
  if (!state.value.lastFetched) return true
  return Date.now() - state.value.lastFetched > CACHE_TTL_MS
}

export function useExchangeRate() {
  const settingsStore = useSettingsStore()
  const baseCurrency = computed(() => settingsStore.baseCurrency || 'USD')

  const rate = computed(() => state.value.rate)
  const loading = computed(() => state.value.loading)
  const error = computed(() => state.value.error)

  async function fetchRate(): Promise<number> {
    const base = baseCurrency.value
    if (base === 'USD') {
      state.value.rate = 1
      state.value.lastFetched = Date.now()
      return 1
    }
    if (!isStale() && state.value.rate !== null) {
      return state.value.rate
    }

    state.value.loading = true
    state.value.error = null
    try {
      const res = await fetch(`${FRANKFURTER_BASE}?from=${base}&to=USD`)
      if (!res.ok) throw new Error(`Frankfurter API ${res.status}`)
      const data = await res.json()
      const usdRate = data.rates?.USD
      if (typeof usdRate !== 'number') throw new Error('USD rate not found')
      state.value.rate = usdRate
      state.value.lastFetched = Date.now()
      return usdRate
    } catch (e: any) {
      const msg = e instanceof Error ? e.message : 'Unknown error'
      state.value.error = msg
      console.warn('Failed to fetch exchange rate:', e)
      return 1 // fallback 1:1 agar tidak crash UI
    } finally {
      state.value.loading = false
    }
  }

  function toUsd(amount: number): number {
    return amount * (state.value.rate || 1)
  }

  return { rate, loading, error, fetchRate, toUsd }
}
