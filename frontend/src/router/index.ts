import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { guest: true }
    },
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('@/views/misc/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/transactions',
      name: 'Transactions',
      component: () => import('@/views/transactions/TransactionsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/transactions/scan',
      name: 'ScanReceipt',
      component: () => import('@/views/transactions/ScanReceiptView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/accounts',
      name: 'Accounts',
      component: () => import('@/views/accounts/AccountsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/categories',
      name: 'Categories',
      component: () => import('@/views/categories/CategoriesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/scheduled',
      name: 'Scheduled',
      component: () => import('@/views/scheduled/ScheduledView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/reports',
      name: 'Reports',
      component: () => import('@/views/reports/ReportsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('@/views/settings/SettingsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/about',
      name: 'About',
      component: () => import('@/views/misc/AboutView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/shared-access',
      name: 'SharedAccess',
      component: () => import('@/views/shared/SharedAccessView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/shared/:token',
      name: 'SharedView',
      component: () => import('@/views/shared/SharedView.vue'),
      meta: { guest: true }
    },
    {
      path: '/calendar',
      name: 'Calendar',
      component: () => import('@/views/misc/CalendarView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/analytics/balance',
      name: 'AnalyticsBalance',
      component: () => import('@/views/analytics/AnalyticsBalance.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/analytics/outlook',
      name: 'AnalyticsOutlook',
      component: () => import('@/views/analytics/AnalyticsOutlook.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/analytics/cashflow',
      name: 'AnalyticsCashflow',
      component: () => import('@/views/analytics/AnalyticsCashflow.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/analytics/spending',
      name: 'AnalyticsSpending',
      component: () => import('@/views/analytics/AnalyticsSpending.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/analytics/budget',
      name: 'AnalyticsBudget',
      component: () => import('@/views/analytics/AnalyticsBudget.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/analytics/goals',
      name: 'AnalyticsGoals',
      component: () => import('@/views/analytics/AnalyticsGoals.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.guest && authStore.isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
