import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '@/views/misc/DashboardView.vue'
import NotificationsView from '@/views/notifications/NotificationsView.vue'
import TransactionsView from '@/views/transactions/TransactionsView.vue'
import ScanReceiptView from '@/views/transactions/ScanReceiptView.vue'
import BudgetsView from '@/views/budgets/BudgetsView.vue'
import AccountsView from '@/views/accounts/AccountsView.vue'
import CategoriesView from '@/views/categories/CategoriesView.vue'
import SubCategoriesView from '@/views/categories/SubCategoriesView.vue'
import ScheduledView from '@/views/scheduled/ScheduledView.vue'
import ReportsView from '@/views/reports/ReportsView.vue'
import SettingsView from '@/views/settings/SettingsView.vue'
import AboutView from '@/views/misc/AboutView.vue'
import SharedAccessView from '@/views/shared/SharedAccessView.vue'
import LabelsView from '@/views/labels/LabelsView.vue'
import SharedView from '@/views/shared/SharedView.vue'
import CalendarView from '@/views/misc/CalendarView.vue'
import AnalyticsBalance from '@/views/analytics/AnalyticsBalance.vue'
import AnalyticsOutlook from '@/views/analytics/AnalyticsOutlook.vue'
import AnalyticsCashflow from '@/views/analytics/AnalyticsCashflow.vue'
import AnalyticsSpending from '@/views/analytics/AnalyticsSpending.vue'
import GoalsView from '@/views/goals/GoalsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: DashboardView
    },
    {
      path: '/dashboard',
      redirect: '/'
    },
    {
      path: '/notifications',
      name: 'Notifications',
      component: NotificationsView
    },
    {
      path: '/transactions',
      name: 'Transactions',
      component: TransactionsView
    },
    {
      path: '/transactions/scan',
      name: 'ScanReceipt',
      component: ScanReceiptView
    },
    {
      path: '/budgets',
      name: 'Budgets',
      component: BudgetsView
    },
    {
      path: '/accounts',
      name: 'Accounts',
      component: AccountsView
    },
    {
      path: '/categories',
      name: 'Categories',
      component: CategoriesView
    },
    {
      path: '/categories/sub/:parentName',
      name: 'SubCategories',
      component: SubCategoriesView
    },
    {
      path: '/scheduled',
      name: 'Scheduled',
      component: ScheduledView
    },
    {
      path: '/reports',
      name: 'Reports',
      component: ReportsView
    },
    {
      path: '/settings',
      name: 'Settings',
      component: SettingsView
    },
    {
      path: '/about',
      name: 'About',
      component: AboutView
    },
    {
      path: '/shared-access',
      name: 'SharedAccess',
      component: SharedAccessView
    },
    {
      path: '/labels',
      name: 'Labels',
      component: LabelsView
    },
    {
      path: '/shared/:token',
      name: 'SharedView',
      component: SharedView
    },
    {
      path: '/calendar',
      name: 'Calendar',
      component: CalendarView
    },
    {
      path: '/analytics/balance',
      name: 'AnalyticsBalance',
      component: AnalyticsBalance
    },
    {
      path: '/analytics/outlook',
      name: 'AnalyticsOutlook',
      component: AnalyticsOutlook
    },
    {
      path: '/analytics/cashflow',
      name: 'AnalyticsCashflow',
      component: AnalyticsCashflow
    },
    {
      path: '/analytics/spending',
      name: 'AnalyticsSpending',
      component: AnalyticsSpending
    },
    {
      path: '/goals',
      name: 'Goals',
      component: GoalsView
    }
  ]
})

export default router