import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import { useSettingsStore } from '@/stores/settings'
import RecentTransactionsWidget from '@/widgets/dashboard/RecentTransactionsWidget.vue'

const routerLinkStub = {
  template: '<a><slot /></a>',
  props: ['to']
}

const transaction = {
  id: 'tx-1',
  type: 'income' as const,
  date_time: '2026-01-15T10:00:00Z',
  amount: 1500,
  currency: 'USD',
  from_account_id: null,
  to_account_id: 'acc-1',
  category_id: 'cat_income_salary',
  labels: '',
  note: '',
  payee_payer: 'Salary',
  payment_type: '',
  warranty_until: null,
  status: 'cleared' as const,
  place_name: '',
  place_lat: null,
  place_lng: null,
  attachments: '',
  template_id: null,
  created_at: '2026-01-15T10:00:00Z',
  updated_at: '2026-01-15T10:00:00Z'
}

const category = {
  id: 'cat_income_salary',
  name: 'Salary',
  type: 'income',
  color: '#ffffff',
  icon: '',
  sort_order: 0,
  archived: false,
  path: ''
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

describe('RecentTransactionsWidget formatRowAmount (L10)', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    const settingsStore = useSettingsStore()
    settingsStore.settings.hide_amounts = false
    settingsStore.settings.base_currency = 'USD'
    settingsStore.settings.number_format = 'en-US'
  })

  test('visible income amount includes currency symbol via formatAmountForDisplay', () => {
    const wrapper = mount(RecentTransactionsWidget, {
      props: {
        recentTransactions: [transaction],
        categories: [category],
        accounts: [account]
      },
      global: { stubs: { 'router-link': routerLinkStub } }
    })

    const amountEl = wrapper.find('.row-amount')
    expect(amountEl.text()).toBe('+$1,500')
    expect(amountEl.text()).not.toBe('+1,500')
  })

  test('visible expense amount includes currency symbol via formatAmountForDisplay', () => {
    const wrapper = mount(RecentTransactionsWidget, {
      props: {
        recentTransactions: [{ ...transaction, type: 'expense' as const, amount: 750 }],
        categories: [category],
        accounts: [account]
      },
      global: { stubs: { 'router-link': routerLinkStub } }
    })

    const amountEl = wrapper.find('.row-amount')
    expect(amountEl.text()).toBe('−$750')
    expect(amountEl.text()).not.toBe('−750')
  })
})
