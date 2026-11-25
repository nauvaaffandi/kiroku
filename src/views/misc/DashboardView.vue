<template>
  <div class="dashboard-page">
    <main class="dashboard-main">
      <DashboardSummaryWidget />
      <DashboardAccountsWidget
        @view-accounts="router.push('/accounts')"
        @add-account="router.push('/accounts')"
      />
      <DashboardChartsWidget />
      <RecentTransactionsWidget
        :recentTransactions="recentTransactions"
        :categories="categoriesStore.categories"
        :accounts="accountsStore.accounts"
      />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import { useCategoriesStore } from '@/stores/categories'
import RecentTransactionsWidget from '@/widgets/dashboard/RecentTransactionsWidget.vue'
import DashboardSummaryWidget from '@/widgets/dashboard/DashboardSummaryWidget.vue'
import DashboardAccountsWidget from '@/widgets/dashboard/DashboardAccountsWidget.vue'
import DashboardChartsWidget from '@/widgets/dashboard/DashboardChartsWidget.vue'

const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const categoriesStore = useCategoriesStore()
const router = useRouter()

console.log('[DashboardView] reading store data', {
  transactionCount: transactionsStore.transactions?.length,
  accountCount: accountsStore.accounts?.length,
  categoryCount: categoriesStore.categories?.length
})

const recentTransactions = computed(() => (transactionsStore.transactions || []).slice(0, 5))
</script>

<style scoped>
@reference "tailwindcss";
.dashboard-page { @apply min-h-screen flex flex-col; }
.dashboard-main { @apply flex-1 px-4 py-4 w-full; }
</style>
