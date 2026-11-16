<template>
  <div class="reports-page">
    <header class="reports-header">
      <h1 class="reports-title">{{ t('reports.title') }}</h1>
      <div class="toolbar">
        <select v-model="selectedCurrency" @change="fetchSpending" class="currency-select">
          <option v-for="c of currencies" :key="c" :value="c">{{ c }}</option>
        </select>
        <button @click="exportCSV" class="export-btn" :aria-label="t('reports.exportCsv')">
          <Download :size="18" stroke-width="1.75" aria-hidden="true" />
          {{ t('reports.exportCsv') }}
        </button>
        <button @click="exportExcel" class="export-btn" :aria-label="t('reports.exportExcel')">
          <FileText :size="18" stroke-width="1.75" aria-hidden="true" />
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
import { Download, FileText } from '@lucide/vue'

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

.reports-page { @apply min-h-screen flex flex-col pb-20 w-full; }

.reports-header { @apply p-4 border-b; background-color: var(--bg-card); border-color: var(--border-primary); }

.reports-title { @apply text-[20px] font-semibold text-[var(--text-primary)]; }

.toolbar { @apply flex flex-col sm:flex-row items-stretch sm:items-center gap-3 mt-3; }

.currency-select {
  @apply w-full sm:w-auto px-3 py-2 bg-[var(--bg-card)] border border-[var(--border-secondary)] rounded-[10px] text-sm text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)];
}

.export-btn {
  @apply flex items-center justify-center gap-2 px-4 py-2 bg-[var(--bg-tertiary)] border border-[var(--border-primary)] rounded-[10px] text-sm font-medium transition-colors duration-120 ease-out cursor-pointer border-none;
  color: var(--text-primary);
}

.export-btn:hover {
  @apply bg-[var(--bg-hover)];
}

.reports-body { @apply flex flex-col gap-4 p-4 flex-1; }

.reports-main { @apply flex flex-col gap-4; }

.tab-panel {
  @apply bg-[var(--bg-card)] border border-[var(--border-primary)] rounded-[14px] overflow-hidden;
}

.tab-bar { @apply flex border-b overflow-x-auto; border-color: var(--border-primary); }

.tab-btn { @apply flex-1 px-4 py-3 text-[0.75rem] font-medium whitespace-nowrap transition-colors duration-120 ease-out cursor-pointer border-b-2 border-transparent; color: var(--text-secondary); }

.tab-btn--active { @apply text-[var(--accent)] border-b-2 border-[var(--accent)] bg-[var(--accent-light)]; }

.tab-content { @apply p-4; }
</style>
