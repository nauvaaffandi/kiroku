<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>口座別残高</h2>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="chart-container" ref="accountsChartRef" style="height: 250px;" />
      <div class="accounts-table">
        <table>
          <thead>
            <tr>
              <th>口座</th>
              <th class="text-right">期首残高</th>
              <th class="text-right">期末残高</th>
              <th class="text-right">増減</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="acc in data" :key="acc.id">
              <td>
                <div class="account-cell">
                  <span class="account-color" :style="{ backgroundColor: acc.color }" />
                  {{ acc.name }}
                </div>
              </td>
              <td class="text-right">{{ formatCurrency(acc.opening_balance, acc.currency) }}</td>
              <td class="text-right">{{ formatCurrency(acc.closing_balance, acc.currency) }}</td>
              <td class="text-right" :class="acc.net_change >= 0 ? 'text-success' : 'text-error'">
                {{ acc.net_change >= 0 ? '+' : '' }}{{ formatCurrency(acc.net_change, acc.currency) }}
              </td>
            </tr>
            <tr class="total-row">
              <td><strong>合計</strong></td>
              <td class="text-right"><strong>{{ formatCurrency(totalOpening) }}</strong></td>
              <td class="text-right"><strong>{{ formatCurrency(totalClosing) }}</strong></td>
              <td class="text-right" :class="totalNet >= 0 ? 'text-success' : 'text-error'">
                <strong>{{ totalNet >= 0 ? '+' : '' }}{{ formatCurrency(totalNet) }}</strong>
              </td>
            </tr>
          </tbody>
        </table>
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
import type { AccountBalanceSummary } from '@/types/report'

const props = defineProps<{
  loading: boolean
  data: AccountBalanceSummary[]
}>()

const accountsChartRef = ref<HTMLElement>()
let accountsChart: echarts.ECharts | null = null

const totalOpening = computed(() => props.data.reduce((sum, a) => sum + a.opening_balance, 0))
const totalClosing = computed(() => props.data.reduce((sum, a) => sum + a.closing_balance, 0))
const totalNet = computed(() => props.data.reduce((sum, a) => sum + a.net_change, 0))

function renderAccountsChart() {
  if (!accountsChartRef.value || props.data.length === 0) return
  if (accountsChart) accountsChart.dispose()

  accountsChart = echarts.init(accountsChartRef.value)
  const names = props.data.map(d => d.name)
  const closing = props.data.map(d => d.closing_balance)
  const colors = props.data.map(d => d.color || '#4a90d9')

  accountsChart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: 10, right: 10, top: 10, bottom: 30 },
    xAxis: { type: 'value', axisLabel: { formatter: (v: number) => formatCurrencyShort(v) } },
    yAxis: { type: 'category', data: names, axisLabel: { interval: 0 } },
    series: [
      {
        name: '期末残高',
        type: 'bar',
        data: closing.map((v, i) => ({ value: v, itemStyle: { color: colors[i] } })),
        label: { show: true, position: 'right', formatter: (p: any) => formatCurrencyShort(p.value) },
      },
    ],
  })
}

function formatCurrency(val: number, currency: string = 'JPY'): string {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(val)
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

watch(() => props.data, () => {
  nextTick(() => renderAccountsChart())
}, { deep: true })

onMounted(() => {
  nextTick(() => renderAccountsChart())
  window.addEventListener('resize', () => accountsChart?.resize())
})
</script>

<style scoped>
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; border-bottom: 1px solid var(--border-primary); background-color: var(--bg-secondary); }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.widget-content { padding: 1rem; }
.chart-container { width: 100%; }
.accounts-table { margin-top: 1rem; }
.accounts-table table { width: 100%; font-size: 0.875rem; }
.accounts-table th { text-align: left; padding: 0.5rem 0.75rem; font-weight: 500; color: var(--text-secondary); border-bottom: 1px solid var(--border-primary); }
.accounts-table td { padding: 0.5rem 0.75rem; border-bottom: 1px solid var(--border-primary); }
.accounts-table .total-row { background-color: var(--bg-secondary); }
.accounts-table .text-right { text-align: right; }
.account-cell { display: flex; align-items: center; gap: 0.5rem; }
.account-color { width: 0.75rem; height: 0.75rem; border-radius: 50%; }
.text-success { color: var(--success); }
.text-error { color: var(--error); }
</style>