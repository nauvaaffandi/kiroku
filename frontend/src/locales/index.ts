import en from './en'
import ja from './ja'
import id from './id'
import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

type Locale = 'en' | 'ja' | 'id'

const locales: Record<Locale, typeof en> = { en, ja, id }

const SUPPORTED_LOCALES: Locale[] = ['en', 'ja', 'id']

function getInitialLocale(): Locale {
  const stored = localStorage.getItem('locale') as Locale | null
  if (stored && SUPPORTED_LOCALES.includes(stored)) {
    return stored
  }
  const browserLang = navigator.language.slice(0, 2) as Locale
  if (SUPPORTED_LOCALES.includes(browserLang)) {
    return browserLang
  }
  return 'en'
}

function applyLocale(locale: Locale) {
  document.documentElement.lang = locale
  localStorage.setItem('locale', locale)
}

export const useI18nStore = defineStore('i18n', () => {
  const currentLocale = ref<Locale>(getInitialLocale())

  const t = (key: string, params?: Record<string, string | number>): string => {
    const keys = key.split('.')
    let result: unknown = locales[currentLocale.value]

    for (const k of keys) {
      if (result && typeof result === 'object' && k in result) {
        result = (result as Record<string, unknown>)[k]
      } else {
        let fallback: unknown = locales.en
        for (const fk of keys) {
          if (fallback && typeof fallback === 'object' && fk in fallback) {
            fallback = (fallback as Record<string, unknown>)[fk]
          } else {
            console.warn(`[i18n] Missing translation key: ${key}`)
            return key
          }
        }
        result = fallback
        break
      }
    }

    if (typeof result !== 'string') {
      console.warn(`[i18n] Translation key not a string: ${key}`)
      return key
    }

    if (params) {
      return result.replace(/\{(\w+)\}/g, (_, k) => String(params[k] ?? ''))
    }

    return result
  }

  function setLocale(locale: Locale) {
    if (!SUPPORTED_LOCALES.includes(locale)) {
      console.warn(`[i18n] Unsupported locale: ${locale}`)
      return
    }
    currentLocale.value = locale
    applyLocale(locale)
  }

  const availableLocales = computed(() =>
    SUPPORTED_LOCALES.map(l => ({
      code: l,
      name: l === 'en' ? 'English' : l === 'ja' ? '日本語' : 'Bahasa Indonesia'
    }))
  )

  applyLocale(currentLocale.value)

  return { currentLocale, t, setLocale, availableLocales }
})

export function useI18n() {
  const store = useI18nStore()
  return { t: store.t, locale: store.currentLocale, setLocale: store.setLocale }
}

export type { Locale }
export { SUPPORTED_LOCALES }