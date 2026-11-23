import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { mount } from '@vue/test-utils'
import { useUIStore } from '@/stores/ui'
import Sidebar from '@/components/layout/Sidebar.vue'

const routes = [
  { path: '/', component: { template: '<div />' } },
  { path: '/transactions', component: { template: '<div />' } }
]

let router: ReturnType<typeof createRouter>

beforeEach(() => {
  setActivePinia(createPinia())
  router = createRouter({ history: createWebHistory(import.meta.env.BASE_URL), routes })
})

const mountSidebar = (sidebarOpen = false) => {
  const uiStore = useUIStore()
  uiStore.sidebarOpen = sidebarOpen
  return mount(Sidebar, {
    global: {
      plugins: [router],
      stubs: { SidebarNav: true, SidebarFooter: true }
    }
  })
}

describe('Sidebar — Collapsible Icon Sidebar (planned API, TDD red phase)', () => {
  test('collapsed by default: ChevronRight toggle, logo/app name hidden', () => {
    const wrapper = mountSidebar(false)
    const toggle = wrapper.find('button.sidebar-toggle')
    expect(toggle.exists()).toBe(true)
    expect(wrapper.find('.lucide-chevron-right').exists()).toBe(true)
    expect(wrapper.find('.lucide-chevron-left').exists()).toBe(false)
    expect(wrapper.text()).not.toContain('Kiroku')
  })

  test('expanded: ChevronLeft toggle, logo/app name shown', () => {
    const wrapper = mountSidebar(true)
    const toggle = wrapper.find('button.sidebar-toggle')
    expect(toggle.exists()).toBe(true)
    expect(wrapper.find('.lucide-chevron-left').exists()).toBe(true)
    expect(wrapper.find('.lucide-chevron-right').exists()).toBe(false)
    expect(wrapper.text()).toContain('Kiroku')
  })

  test('clicking toggle switches collapsed/expanded state', async () => {
    const wrapper = mountSidebar(false)
    const toggle = wrapper.find('button.sidebar-toggle')
    expect(toggle.exists()).toBe(true)
    toggle.trigger('click')
    await wrapper.vm.$nextTick()
    expect(wrapper.find('.lucide-chevron-left').exists()).toBe(true)
    expect(wrapper.text()).toContain('Kiroku')
  })

  test('toggle button title attribute shows correct tooltip per state', () => {
    const collapsed = mountSidebar(false)
    const cToggle = collapsed.find('button.sidebar-toggle')
    expect(cToggle.exists()).toBe(true)
    expect(cToggle.attributes('title')).toBe('Expand sidebar')

    const expanded = mountSidebar(true)
    const eToggle = expanded.find('button.sidebar-toggle')
    expect(eToggle.exists()).toBe(true)
    expect(eToggle.attributes('title')).toBe('Collapse sidebar')
  })
})
