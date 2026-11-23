import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { mount } from '@vue/test-utils'
import SidebarNav from '@/components/layout/SidebarNav.vue'
import ThemeToggle from '@/components/common/ThemeToggle.vue'

const routes = [
  { path: '/', component: { template: '<div />' } },
  { path: '/transactions', component: { template: '<div />' } },
  { path: '/transactions/scan', component: { template: '<div />' } },
  { path: '/budgets', component: { template: '<div />' } },
  { path: '/accounts', component: { template: '<div />' } },
  { path: '/categories', component: { template: '<div />' } },
  { path: '/categories/sub/:parentName', component: { template: '<div />' } },
  { path: '/scheduled', component: { template: '<div />' } },
  { path: '/reports', component: { template: '<div />' } },
  { path: '/settings', component: { template: '<div />' } },
  { path: '/about', component: { template: '<div />' } },
  { path: '/labels', component: { template: '<div />' } },
  { path: '/calendar', component: { template: '<div />' } },
  { path: '/analytics/balance', component: { template: '<div />' } },
  { path: '/analytics/outlook', component: { template: '<div />' } },
  { path: '/analytics/cashflow', component: { template: '<div />' } },
  { path: '/analytics/spending', component: { template: '<div />' } },
  { path: '/goals', component: { template: '<div />' } }
]

let router: ReturnType<typeof createRouter>

beforeEach(() => {
  setActivePinia(createPinia())
  router = createRouter({ history: createWebHistory(import.meta.env.BASE_URL), routes })
})

const mountNav = async (routePath = '/', collapsed = false) => {
  await router.push(routePath)
  return mount(SidebarNav, {
    props: { collapsed },
    global: {
      plugins: [router],
      stubs: { ThemeToggle: true }
    }
  })
}

describe('SidebarNav — collapsible nav (planned API, TDD red phase)', () => {
  test('renders all nav items with icons', async () => {
    const wrapper = await mountNav('/')
    expect(wrapper.findAll('.nav-icon').length).toBeGreaterThanOrEqual(3)
  })

  test('collapsed: labels hidden, icons present, title tooltip on nav links', async () => {
    const wrapper = await mountNav('/', true)
    expect(wrapper.findAll('.nav-link span').length).toBe(0)
    expect(wrapper.findAll('.nav-icon').length).toBeGreaterThan(0)
    const links = wrapper.findAll('.nav-link')
    links.forEach((link) => {
      expect(link.attributes('title')).toBeTruthy()
    })
  })

  test('expanded: labels shown', async () => {
    const wrapper = await mountNav('/', false)
    expect(wrapper.findAll('.nav-link span').length).toBeGreaterThan(0)
    expect(wrapper.text()).toContain('Dashboard')
    expect(wrapper.text()).toContain('Transactions')
  })

  test('active state applied to current route link', async () => {
    const wrapper = await mountNav('/transactions')
    const links = wrapper.findAll('.nav-link')
    const tx = links.find((el) => el.attributes('href') === '/transactions')
    expect(tx).toBeDefined()
    expect(tx.classes()).toContain('active')
    const acc = links.find((el) => el.attributes('href') === '/accounts')
    expect(acc).toBeDefined()
    expect(acc.classes()).not.toContain('active')
  })

  test('analytics sub-nav expand/collapse toggles rotate-180', async () => {
    const wrapper = await mountNav('/')
    const toggle = wrapper.find('.nav-group-toggle')
    expect(toggle.exists()).toBe(true)
    const icon = wrapper.find('.nav-expand-icon')
    expect(icon.classes()).not.toContain('rotate-180')
    await toggle.trigger('click')
    expect(icon.classes()).toContain('rotate-180')
    await toggle.trigger('click')
    expect(icon.classes()).not.toContain('rotate-180')
  })
})
