<template>
  <div class="dashboard-page">
    <main class="dashboard-main">
      <DashboardSummaryWidget />
      <DashboardAccountsWidget />
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

const recentTransactions = computed(() => (transactionsStore.transactions || []).slice(0, 5))
</script>

<style scoped>
@reference "tailwindcss";
.dashboard-page { @apply min-h-screen flex flex-col; }
.dashboard-main { @apply flex-1 px-4 py-6 mx-auto w-full max-w-[1200px]; }
</style>