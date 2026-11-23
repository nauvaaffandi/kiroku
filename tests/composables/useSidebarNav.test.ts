import { describe, test, expect, beforeEach, afterEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({ path: '/' }))
}))

import { useRoute } from 'vue-router'
import { useUIStore } from '@/stores/ui'
import { useSidebarNav } from '@/composables/useSidebarNav'
import DashboardIcon from '@/components/icons/DashboardIcon.vue'
import AccountsIcon from '@/components/icons/AccountsIcon.vue'
import TransactionsIcon from '@/components/icons/TransactionsIcon.vue'
import CalendarIcon from '@/components/icons/CalendarIcon.vue'
import BudgetIcon from '@/components/icons/BudgetIcon.vue'
import GoalsIcon from '@/components/icons/GoalsIcon.vue'
import ScheduledIcon from '@/components/icons/ScheduledIcon.vue'
import ReportsIcon from '@/components/icons/ReportsIcon.vue'
import CategoriesIcon from '@/components/icons/CategoriesIcon.vue'
import LabelsIcon from '@/components/icons/LabelsIcon.vue'

const mockUseRoute = vi.mocked(useRoute)

beforeEach(() => {
  setActivePinia(createPinia())
  mockUseRoute.mockReturnValue({ path: '/' })
})

afterEach(() => {
  ;(window as any).innerWidth = 1024
})

describe('useSidebarNav — nav items, isActive, closeMobileSidebar (TDD red phase)', () => {
  test('mainNavItems have correct paths/labels/icons', () => {
    const { mainNavItems } = useSidebarNav()
    expect(mainNavItems[0]).toMatchObject({ path: '/', label: 'Dashboard', color: '#3b82f6' })
    expect(mainNavItems[0].icon).toBe(DashboardIcon)
    expect(mainNavItems[1]).toMatchObject({ path: '/accounts', label: 'Accounts', color: '#06b6d4' })
    expect(mainNavItems[1].icon).toBe(AccountsIcon)
    expect(mainNavItems[2]).toMatchObject({ path: '/transactions', label: 'Transactions', color: '#10b981' })
    expect(mainNavItems[2].icon).toBe(TransactionsIcon)
  })

  test('analyticsItems have correct paths/labels', () => {
    const { analyticsItems } = useSidebarNav()
    expect(analyticsItems[0]).toMatchObject({ path: '/analytics/balance', label: 'Balance', color: '#6366f1' })
    expect(analyticsItems[1]).toMatchObject({ path: '/analytics/outlook', label: 'Outlook', color: '#8b5cf6' })
    expect(analyticsItems[2]).toMatchObject({ path: '/analytics/cashflow', label: 'Cashflow', color: '#ec4899' })
    expect(analyticsItems[3]).toMatchObject({ path: '/analytics/spending', label: 'Spending', color: '#f43f5e' })
    expect(analyticsItems[0].icon).toBeDefined()
    expect(analyticsItems[1].icon).toBeDefined()
    expect(analyticsItems[2].icon).toBeDefined()
    expect(analyticsItems[3].icon).toBeDefined()
  })

  test('goalsItems have correct paths/labels/icons', () => {
    const { goalsItems } = useSidebarNav()
    expect(goalsItems[0]).toMatchObject({ path: '/calendar', label: 'Calendar', color: '#f59e0b' })
    expect(goalsItems[0].icon).toBe(CalendarIcon)
    expect(goalsItems[1]).toMatchObject({ path: '/budgets', label: 'Budgets', color: '#14b8a6' })
    expect(goalsItems[1].icon).toBe(BudgetIcon)
    expect(goalsItems[2]).toMatchObject({ path: '/goals', label: 'Goals', color: '#eab308' })
    expect(goalsItems[2].icon).toBe(GoalsIcon)
  })

  test('section3NavItems and section4NavItems have correct paths/labels/icons', () => {
    const { section3NavItems, section4NavItems } = useSidebarNav()
    expect(section3NavItems[0]).toMatchObject({ path: '/scheduled', label: 'Scheduled', color: '#0ea5e9' })
    expect(section3NavItems[0].icon).toBe(ScheduledIcon)
    expect(section3NavItems[1]).toMatchObject({ path: '/reports', label: 'Reports', color: '#a855f7' })
    expect(section3NavItems[1].icon).toBe(ReportsIcon)
    expect(section4NavItems[0]).toMatchObject({ path: '/categories', label: 'Categories', color: '#84cc16' })
    expect(section4NavItems[0].icon).toBe(CategoriesIcon)
    expect(section4NavItems[1]).toMatchObject({ path: '/labels', label: 'Labels', color: '#6366f1' })
    expect(section4NavItems[1].icon).toBe(LabelsIcon)
  })

  test('isActive matches current route path', () => {
    mockUseRoute.mockReturnValue({ path: '/transactions' })
    const { isActive } = useSidebarNav()
    expect(isActive('/transactions')).toBe(true)
    expect(isActive('/accounts')).toBe(false)
  })

  test('isActive treats "/" as prefix match for all paths (current behavior)', () => {
    mockUseRoute.mockReturnValue({ path: '/settings' })
    const { isActive } = useSidebarNav()
    expect(isActive('/')).toBe(true)
    expect(isActive('/settings')).toBe(true)
  })

  test('closeMobileSidebar closes sidebar on mobile, not on desktop', () => {
    const uiStore = useUIStore()

    // mobile: fresh composable so isMobile caches true (window.innerWidth not reactive)
    ;(window as any).innerWidth = 360
    const mobileNav = useSidebarNav()
    uiStore.sidebarOpen = true
    mobileNav.closeMobileSidebar()
    expect(uiStore.sidebarOpen).toBe(false)

    // desktop: fresh composable so isMobile caches false
    ;(window as any).innerWidth = 1024
    const desktopNav = useSidebarNav()
    uiStore.sidebarOpen = true
    desktopNav.closeMobileSidebar()
    expect(uiStore.sidebarOpen).toBe(true)
  })
})
