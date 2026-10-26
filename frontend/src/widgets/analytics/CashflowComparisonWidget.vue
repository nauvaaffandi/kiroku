<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>期間比較</h2>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="chart-container" ref="comparisonChartRef" />
      <div class="delta-summary">
        <div class="delta-card" :class="delta.income_pct >= 0 ? 'positive' : 'negative'">
          <span class="delta-label">収入変化</span>
          <span class="delta-value">{{ delta.income_pct >= 0 ? '+' : '' }}{{ delta.income_pct.toFixed(1) }}%</span>
        </div>
        <div class="delta-card" :class="delta.expense_pct >= 0 ? 'negative' : 'positive'">
          <span class="delta-label">支出変化</span>
          <span class="delta-value">{{ delta.expense_pct >= 0 ? '+' : '' }}{{ delta.expense_pct.toFixed(1) }}%</span>
        </div>
        <div class="delta-card" :class="delta.net_pct >= 0 ? 'positive' : 'negative'">
          <span class="delta-label">純額変化</span>
          <span class="delta-value">{{ delta.net_pct >= 0 ? '+' : '' }}{{ delta.net_pct.toFixed(1) }}%</span>
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

interface ComparisonResponse {
  period_a: {
    income: number
    expense: number
    net: number
    avg_daily_income: number
    avg_daily_expense: number
    transaction_count: number
  }
  period_b: {
    income: number
    expense: number
    net: number
    avg_daily_income: number
    avg_daily_expense: number
    transaction_count: number
  }
  delta: {
    income_pct: number
    expense_pct: number
    net_pct: number
  }
}

const props = defineProps<{
  loading: boolean
  data: ComparisonResponse | null
}>()

const comparisonChartRef = ref<HTMLElement>()
let comparisonChart: echarts.ECharts | null = null

const delta = computed(() => props.data?.delta || { income_pct: 0, expense_pct: 0, net_pct: 0 })

function renderComparisonChart() {
  if (!comparisonChartRef.value || !props.data) return
  if (comparisonChart) comparisonChart.dispose()

  comparisonChart = echarts.init(comparisonChartRef.value)
  const data = props.data

  comparisonChart.setOption({
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        let res = `${params[0].name}<br/>`
        params.forEach((p: any) => {
          res += `${p.marker} ${p.seriesName}: ${formatCurrencyShort(p.value)}<br/>`
        })
        return res
      },
    },
    legend: { data: ['期間A 収入', '期間A 支出', '期間B 収入', '期間B 支出'], bottom: 0 },
    grid: { left: 40, right: 20, top: 10, bottom: 40 },
    xAxis: { type: 'category', data: ['収入', '支出', '純額'] },
    yAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    series: [
      { name: '期間A 収入', type: 'bar', data: [data.period_a.income, 0, 0], itemStyle: { color: '#28a745' } },
      { name: '期間A 支出', type: 'bar', data: [0, data.period_a.expense, 0], itemStyle: { color: '#dc3545' } },
      { name: '期間B 収入', type: 'bar', data: [data.period_b.income, 0, 0], itemStyle: { color: '#28a745', opacity: 0.6 } },
      { name: '期間B 支出', type: 'bar', data: [0, data.period_b.expense, 0], itemStyle: { color: '#dc3545', opacity: 0.6 } },
    ],
  })
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

watch(() => props.data, () => {
  nextTick(() => renderComparisonChart())
}, { deep: true })

onMounted(() => {
  nextTick(() => renderComparisonChart())
  window.addEventListener('resize', () => comparisonChart?.resize())
})
</script>

<style scoped>
@reference "tailwindcss";

.widget-card {
  @apply border rounded-xl overflow-hidden;
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
  @apply text-lg font-semibold;
  color: var(--text-primary);
}
.widget-content {
  @apply p-4;
}
.chart-container {
  @apply w-full;
  height: 300px;
}
.delta-summary {
  @apply flex gap-4 mt-4 flex-wrap;
}
.delta-card {
  @apply flex-1 min-w-[120px] py-3 px-4 rounded-sm text-center border-l-4;
  background-color: var(--bg-secondary);
}
.delta-card.positive {
  border-color: var(--success);
}
.delta-card.negative {
  border-color: var(--error);
}
.delta-label {
  @apply block text-xs mb-1;
  color: var(--text-muted);
}
.delta-value {
  @apply block text-lg font-bold;
  color: var(--text-primary);
}
</style>
