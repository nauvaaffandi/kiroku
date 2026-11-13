<template>
  <div class="reports-page">
    <header class="reports-header">
      <h1 class="reports-title">{{ t('reports.title') }}</h1>
      <div class="toolbar">
        <select v-model="selectedCurrency" @change="fetchSpending" class="currency-select">
          <option v-for="c of currencies" :key="c" :value="c">{{ c }}</option>
        </select>
        <button @click="exportCSV" class="export-btn">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
          Export CSV
        </button>
        <button @click="exportExcel" class="export-btn">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
          </svg>
          {{ t('reports.exportExcel') }}
        </button>
      </div>
    </header>

    <div class="reports-body">
      <ReportsSidebar
        :spending-params="spendingParams"
        :trend-params="trendParams"
        :accounts="accountsStore.accounts"
        @spending="fetchSpending"
        @trend="fetchTrend"
      />
      <main class="reports-main">
        <div class="tab-panel">
          <div class="tab-bar">
            <button :class="['tab-btn', { 'tab-btn--active': activeTab === 'spending' }]" @click="activeTab = 'spending'">{{ t('reports.spendingReport') }}</button>
            <button :class="['tab-btn', { 'tab-btn--active': activeTab === 'trend' }]" @click="activeTab = 'trend'">{{ t('reports.trendReport') }}</button>
            <button :class="['tab-btn', { 'tab-btn--active': activeTab === 'cashflow' }]" @click="activeTab = 'cashflow'">{{ t('reports.cashflowTable.title') }}</button>
          </div>
          <div v-if="activeTab === 'spending'" class="tab-content">
            <ReportsSpendingPanel
              :period="spendingParams.period"
              :groups="spendingReport.groups"
              :total-amount="totalAmount"
              :format-display-number="formatDisplayNumber"
            />
          </div>
          <div v-if="activeTab === 'trend'" class="tab-content">
            <ReportsTrendPanel :series="trendReport.series" />
          </div>
          <div v-if="activeTab === 'cashflow'" class="tab-content">
            <CashflowTable :date-from="spendingParams.date_from" :date-to="spendingParams.date_to" />
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from '@/locales'
import { useReportsPage } from '@/composables/useReportsPage'
import ReportsSidebar from '@/components/reports/ReportsSidebar.vue'
import ReportsSpendingPanel from '@/components/reports/ReportsSpendingPanel.vue'
import ReportsTrendPanel from '@/components/reports/ReportsTrendPanel.vue'
import CashflowTable from '@/components/charts/CashflowTable.vue'

const { t } = useI18n()
const activeTab = ref('spending')
const {
  spendingParams, trendParams, accountsStore, spendingReport, trendReport, totalAmount,
  formatDisplayNumber, fetchSpending, fetchTrend, exportCSV, exportExcel, currencies,
  selectedCurrency
} = useReportsPage()
</script>

<style scoped>
@reference "tailwindcss";

.reports-page { @apply min-h-screen flex flex-col pb-20; }
.reports-header { @apply p-4 bg-zinc-100 dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800; }
.reports-title { @apply text-xl font-bold text-zinc-900 dark:text-zinc-100; }
.toolbar { @apply flex flex-col sm:flex-row items-stretch sm:items-center gap-3 mt-3; }
.currency-select {
  @apply w-full sm:w-auto px-3 py-2 bg-white dark:bg-zinc-800 border border-zinc-300 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:outline-none focus:ring-2 focus:ring-sky-500;
}
.export-btn {
  @apply flex items-center justify-center gap-2 px-4 py-2 bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 text-zinc-800 dark:text-zinc-200 text-sm font-medium rounded-lg transition-colors cursor-pointer border-none;
}
.reports-body { @apply flex flex-col lg:grid lg:grid-cols-[280px_1fr] gap-4 p-4 flex-1; }
.reports-main { @apply flex flex-col gap-4; }
.tab-panel { @apply bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden shadow-sm; }
.tab-bar { @apply flex bg-zinc-100 dark:bg-zinc-800 border-b border-zinc-200 dark:border-zinc-700 overflow-x-auto; }
.tab-btn { @apply flex-1 px-4 py-3 text-xs sm:text-sm font-medium whitespace-nowrap transition-colors cursor-pointer border-b-2; }
.tab-btn--active { @apply text-sky-500 border-sky-500 bg-white dark:bg-zinc-900; }
.tab-content { @apply p-4; }
</style>
