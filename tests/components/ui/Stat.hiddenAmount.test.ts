import { describe, test, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import { useSettingsStore } from '@/stores/settings'
import { useI18nStore } from '@/locales'
import Stat from '@/components/ui/Stat.vue'

describe('Stat hidden amount aria label (L2)', () => {
  let i18n: ReturnType<typeof useI18nStore>
  let tSpy: ReturnType<typeof vi.spyOn>

  beforeEach(() => {
    setActivePinia(createPinia())
    i18n = useI18nStore()
    i18n.setLocale('ja')
    tSpy = vi.spyOn(i18n, 't')
  })

  test('aria-label uses t("stat.amountHidden") not hard-coded English', () => {
    const settingsStore = useSettingsStore()
    settingsStore.settings.hide_amounts = true

    const wrapper = mount(Stat, {
      props: { label: 'Total Balance', value: 1000, type: 'income', hideAmounts: true }
    })

    const valueEl = wrapper.find('.stat-value')
    const ariaLabel = valueEl.attributes('aria-label')

    expect(tSpy).toHaveBeenCalledWith('stat.amountHidden')
    expect(ariaLabel).not.toContain('(amount hidden)')
  })

  test('visible amount aria-label does not call amountHidden key', () => {
    const settingsStore = useSettingsStore()
    settingsStore.settings.hide_amounts = false

    const wrapper = mount(Stat, {
      props: { label: 'Monthly Income', value: 5000, type: 'income', hideAmounts: false }
    })

    const valueEl = wrapper.find('.stat-value')
    const ariaLabel = valueEl.attributes('aria-label')

    expect(ariaLabel).toBe('Monthly Income +$5,000')
    expect(tSpy).not.toHaveBeenCalledWith('stat.amountHidden')
  })
})
