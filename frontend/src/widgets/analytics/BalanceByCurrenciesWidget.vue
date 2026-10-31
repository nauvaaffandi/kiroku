<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>{{ t('balanceAnalytics.currencies') }}</h2>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="chart-container" ref="currenciesChartRef" />
      <div class="currencies-summary">
        <div
          v-for="curr in data"
          :key="curr.currency"
          class="currency-row"
        >
          <span class="currency-code">{{ curr.currency }}</span>
          <div class="currency-values">
            <span>{{ formatCurrency(curr.opening_balance, curr.currency) }}</span>
            <span class="arrow">→</span>
            <span>{{ formatCurrency(curr.closing_balance, curr.currency) }}</span>
            <span :class="curr.net_change >= 0 ? 'text-success' : 'text-error'">
              ({{ curr.net_change >= 0 ? '+' : '' }}{{ formatCurrency(curr.net_change, curr.currency) }})
            </span>
          </div>
        </div>
        <div class="currency-row total">
          <span><strong>{{ t('balanceAnalytics.convertedTotal') }}</strong></span>
          <div class="currency-values">
            <strong>{{ formatCurrency(totalConvertedOpening) }}</strong>
            <span class="arrow">→</span>
            <strong>{{ formatCurrency(totalConvertedClosing) }}</strong>
            <span :class="totalConvertedNet >= 0 ? 'text-success' : 'text-error'">
              ({{ totalConvertedNet >= 0 ? '+' : '' }}{{ formatCurrency(totalConvertedNet) }})
            </span>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="widget-content">
      <div class="skeleton skeleton-chart" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, computed } from 'vue'
import * as echarts from 'echarts'
import type { CurrencyBalanceSummary } from '@/types/report'
import { useSettingsStore } from '@/stores/settings'
import { formatAmountForDisplay } from '@/utils/amountDisplay'
import { useI18n } from '@/locales'

const { t } = useI18n()
const settingsStore = useSettingsStore()

const props = defineProps<{
  loading: boolean
  data: CurrencyBalanceSummary[]
}>()

const currenciesChartRef = ref<HTMLElement>()
let currenciesChart: echarts.ECharts | null = null

const totalConvertedOpening = computed(() => props.data.reduce((sum, c) => sum + c.opening_balance * c.fx_rate_to_base, 0))
const totalConvertedClosing = computed(() => props.data.reduce((sum, c) => sum + c.closing_balance * c.fx_rate_to_base, 0))
const totalConvertedNet = computed(() => props.data.reduce((sum, c) => sum + c.net_change * c.fx_rate_to_base, 0))

function renderCurrenciesChart() {
  if (!currenciesChartRef.value || props.data.length === 0) return
  if (currenciesChart) currenciesChart.dispose()

  currenciesChart = echarts.init(currenciesChartRef.value)
  const names = props.data.map(d => d.currency)
  const opening = props.data.map(d => d.opening_balance)
  const closing = props.data.map(d => d.closing_balance)

  currenciesChart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { data: [t('balanceAnalytics.legendOpening'), t('balanceAnalytics.legendClosing')], bottom: 0 },
    grid: { left: 10, right: 20, top: 10, bottom: 50, containLabel: true },
    xAxis: { type: 'category', data: names },
    yAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    series: [
      { name: t('balanceAnalytics.legendOpening'), type: 'bar', stack: 'total', data: opening, itemStyle: { color: '#999' } },
      { name: t('balanceAnalytics.legendClosing'), type: 'bar', stack: 'total', data: closing, itemStyle: { color: '#4a90d9' } },
    ],
  })
}

function formatCurrency(val: number, currency: string = 'JPY'): string {
  if (settingsStore.isHideAmounts) {
    return formatAmountForDisplay(val, 'transfer')
  }
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(val)
}

function formatCurrencyShort(val: number): string {
  if (settingsStore.isHideAmounts) {
    return formatAmountForDisplay(val, 'transfer')
  }
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

watch(() => props.data, () => {
  nextTick(() => currenciesChart?.resize())
}, { deep: true })

onMounted(() => {
  nextTick(() => {
    renderCurrenciesChart()
    currenciesChart?.resize()
  })
  window.addEventListener('resize', () => currenciesChart?.resize())
})
</script>

<style scoped>
@reference "tailwindcss";

.widget-card {
  @apply border rounded-xl overflow-hidden w-full max-w-full box-border;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  box-shadow: var(--shadow);
}
.widget-header {
  @apply flex items-center justify-between p-4 border-b;
  background-color: var(--bg-secondary);
  border-color: var(--border-primary);
}
.widget-header h2 {
  @apply text-lg font-semibold truncate;
  color: var(--text-primary);
}
.widget-content {
  @apply p-4 overflow-hidden w-full max-w-full box-border;
}
.chart-container {
  @apply w-full max-w-full;
  height: 240px;
}
.currencies-summary {
  @apply mt-4 w-full;
}
.currency-row {
  @apply flex flex-col sm:flex-row gap-1 justify-between items-start sm:items-center p-2 border-b;
  border-color: var(--border-primary);
  white-space: normal;
}
.currency-row.total {
  @apply border-t-2;
  background-color: var(--bg-secondary);
  border-color: var(--border-primary);
}
.currency-code {
  @apply font-mono font-medium w-full sm:w-auto truncate;
  color: var(--text-primary);
}
.currency-values {
  @apply flex flex-col sm:flex-row gap-2 items-start sm:items-center text-sm w-full sm:w-auto;
  color: var(--text-secondary);
}
.currency-values .arrow {
  color: var(--text-muted);
}
.text-success {
  color: var(--success);
}
.text-error {
  color: var(--error);
}
</style>
