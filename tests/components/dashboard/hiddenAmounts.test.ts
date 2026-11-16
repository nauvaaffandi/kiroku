import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import { useSettingsStore } from '@/stores/settings'
import { useAccountsStore } from '@/stores/accounts'
import DashboardAccountsWidget from '@/widgets/dashboard/DashboardAccountsWidget.vue'

const BULLET_MASK = '••••'

const routerLinkStub = {
  template: '<a><slot /></a>',
  props: ['to']
}

const account = {
  id: 'acc-1',
  name: 'Cash',
  type: 'cash' as const,
  currency: 'USD',
  balance: 1500,
  icon: '',
  color: '#ffffff',
  sort_order: 0,
  archived: false
}

describe('dashboard hidden amounts mask consistency (L1)', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    const settingsStore = useSettingsStore()
    settingsStore.settings.hide_amounts = true
    settingsStore.settings.base_currency = 'USD'
    settingsStore.settings.number_format = 'en-US'
  })

  test('DashboardAccountsWidget masks account balance with ••••', () => {
    const accountsStore = useAccountsStore()
    accountsStore.accounts = [account]

    const wrapper = mount(DashboardAccountsWidget, {
      global: { stubs: { 'router-link': routerLinkStub } }
    })

    const amountEl = wrapper.find('.account-amount')
    expect(amountEl.text()).toBe('$••••')
    expect(amountEl.text()).not.toContain('*****')
  })
})
