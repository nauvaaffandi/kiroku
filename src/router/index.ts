import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: () => import('@/views/misc/DashboardView.vue')
    },
    {
      path: '/dashboard',
      redirect: '/'
    },
    {
      path: '/notifications',
      name: 'Notifications',
      component: () => import('@/views/notifications/NotificationsView.vue')
    },
    {
      path: '/transactions',
      name: 'Transactions',
      component: () => import('@/views/transactions/TransactionsView.vue')
    },
    {
      path: '/transactions/scan',
      name: 'ScanReceipt',
      component: () => import('@/views/transactions/ScanReceiptView.vue')
    },
    {
      path: '/budgets',
      name: 'Budgets',
      component: () => import('@/views/budgets/BudgetsView.vue')
    },
    {
      path: '/accounts',
      name: 'Accounts',
      component: () => import('@/views/accounts/AccountsView.vue')
    },
    {
      path: '/categories',
      name: 'Categories',
      component: () => import('@/views/categories/CategoriesView.vue')
    },
    {
      path: '/categories/sub/:parentName',
      name: 'SubCategories',
      component: () => import('@/views/categories/SubCategoriesView.vue')
    },
    {
      path: '/scheduled',
      name: 'Scheduled',
      component: () => import('@/views/scheduled/ScheduledView.vue')
    },
    {
      path: '/reports',
      name: 'Reports',
      component: () => import('@/views/reports/ReportsView.vue')
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('@/views/settings/SettingsView.vue')
    },
    {
      path: '/about',
      name: 'About',
      component: () => import('@/views/misc/AboutView.vue')
    },
    {
      path: '/shared-access',
      name: 'SharedAccess',
      component: () => import('@/views/shared/SharedAccessView.vue')
    },
    {
      path: '/labels',
      name: 'Labels',
      component: () => import('@/views/labels/LabelsView.vue')
    },
    {
      path: '/shared/:token',
      name: 'SharedView',
      component: () => import('@/views/shared/SharedView.vue')
    },
    {
      path: '/calendar',
      name: 'Calendar',
      component: () => import('@/views/misc/CalendarView.vue')
    },
    {
      path: '/analytics/balance',
      name: 'AnalyticsBalance',
      component: () => import('@/views/analytics/AnalyticsBalance.vue')
    },
    {
      path: '/analytics/outlook',
      name: 'AnalyticsOutlook',
      component: () => import('@/views/analytics/AnalyticsOutlook.vue')
    },
    {
      path: '/analytics/cashflow',
      name: 'AnalyticsCashflow',
      component: () => import('@/views/analytics/AnalyticsCashflow.vue')
    },
    {
      path: '/analytics/spending',
      name: 'AnalyticsSpending',
      component: () => import('@/views/analytics/AnalyticsSpending.vue')
    },
    {
      path: '/goals',
      name: 'Goals',
      component: () => import('@/views/goals/GoalsView.vue')
    }
  ]
})

export default router