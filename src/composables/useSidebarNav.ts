import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'
import DashboardIcon from '@/components/icons/DashboardIcon.vue'
import TransactionsIcon from '@/components/icons/TransactionsIcon.vue'
import CalendarIcon from '@/components/icons/CalendarIcon.vue'
import AccountsIcon from '@/components/icons/AccountsIcon.vue'
import CategoriesIcon from '@/components/icons/CategoriesIcon.vue'
import LabelsIcon from '@/components/icons/LabelsIcon.vue'
import ScheduledIcon from '@/components/icons/ScheduledIcon.vue'
import ReportsIcon from '@/components/icons/ReportsIcon.vue'
import SharedAccessIcon from '@/components/icons/SharedAccessIcon.vue'
import BudgetIcon from '@/components/icons/BudgetIcon.vue'
import GoalsIcon from '@/components/icons/GoalsIcon.vue'

export function useSidebarNav() {
  const route = useRoute()
  const uiStore = useUIStore()
  const { t } = useI18n()
  const isMobile = computed(() => window.innerWidth < 768)

  const mainNavItems = [
    { path: '/', label: t('nav.dashboard'), icon: DashboardIcon, color: '#3b82f6' },
    { path: '/accounts', label: t('nav.accounts'), icon: AccountsIcon, color: '#06b6d4' },
    { path: '/transactions', label: t('nav.transactions'), icon: TransactionsIcon, color: '#10b981' }
  ]

  const analyticsItems = [
    { path: '/analytics/balance', label: t('nav.analyticsBalance'), color: '#6366f1' },
    { path: '/analytics/outlook', label: t('nav.analyticsOutlook'), color: '#8b5cf6' },
    { path: '/analytics/cashflow', label: t('nav.analyticsCashflow'), color: '#ec4899' },
    { path: '/analytics/spending', label: t('nav.analyticsSpending'), color: '#f43f5e' }
  ]

  const goalsItems = [
    { path: '/calendar', label: t('nav.calendar'), icon: CalendarIcon, color: '#f59e0b' },
    { path: '/budgets', label: t('nav.budgets'), icon: BudgetIcon, color: '#14b8a6' },
    { path: '/goals', label: t('nav.goals'), icon: GoalsIcon, color: '#eab308' }
  ]

  const section3NavItems = [
    { path: '/scheduled', label: t('nav.scheduled'), icon: ScheduledIcon, color: '#0ea5e9' },
    { path: '/reports', label: t('nav.reports'), icon: ReportsIcon, color: '#a855f7' },
    { path: '/shared-access', label: t('nav.sharedAccess'), icon: SharedAccessIcon, color: '#f97316' }
  ]

  const section4NavItems = [
    { path: '/categories', label: t('nav.categories'), icon: CategoriesIcon, color: '#84cc16' },
    { path: '/labels', label: t('nav.labels'), icon: LabelsIcon, color: '#6366f1' }
  ]

  const isActive = (path: string) => route.path.startsWith(path)

  const closeMobileSidebar = () => {
    if (isMobile.value) uiStore.sidebarOpen = false
  }

  return {
    mainNavItems, analyticsItems, goalsItems, section3NavItems, section4NavItems,
    isActive, closeMobileSidebar, isMobile, t
  }
}
