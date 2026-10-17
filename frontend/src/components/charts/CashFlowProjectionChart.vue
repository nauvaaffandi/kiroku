<template>
  <div class="cashflow-chart-container">
    <div class="chart-controls">
      <div class="control-group">
        <label class="control-label">期間</label>
        <div class="preset-buttons">
          <button
            v-for="preset in presets"
            :key="preset.value"
            :class="['preset-btn', { active: dateRangePreset === preset.value }]"
            @click="setPreset(preset.value)"
          >
            {{ preset.label }}
          </button>
        </div>
      </div>
      <div class="control-group">
        <label class="control-label">表示</label>
        <div class="toggle-group">
          <label class="toggle-item">
            <input type="checkbox" v-model="showIncome" @change="updateChart" />
            <span>収入</span>
          </label>
          <label class="toggle-item">
            <input type="checkbox" v-model="showExpense" @change="updateChart" />
            <span>支出</span>
          </label>
          <label class="toggle-item">
            <input type="checkbox" v-model="showNet" @change="updateChart" />
            <span>残高</span>
          </label>
        </div>
      </div>
      <div class="control-group">
        <label class="control-label">口座</label>
        <select v-model="selectedAccount" @change="fetchData" class="filter-select">
          <option value="">すべて</option>
          <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
            {{ acc.name }} ({{ acc.currency }})
          </option>
        </select>
      </div>
      <div class="control-group">
        <label class="control-label">通貨</label>
        <select v-model="selectedCurrency" @change="fetchData" class="filter-select">
          <option v-for="cur in currencies" :key="cur" :value="cur">{{ cur }}</option>
        </select>
      </div>
    </div>

    <div class="chart-wrapper" ref="chartRef" />
    
    <div v-if="loading" class="chart-loading">
      <div class="skeleton skeleton-chart" style="height: 300px;" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick, computed } from 'vue'
import * as echarts from 'echarts'
import { useAccountsStore } from '@/stores/accounts'
import { api } from '@/api/client'
import type { ProjectionPoint, CashFlowResponse } from '@/types/report'

const props = defineProps<{
  initialDateFrom?: string
  initialDateTo?: string
}>()

const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null

const accountsStore = useAccountsStore()

const loading = ref(true)
const projectionData = ref<ProjectionPoint[]>([])
const accounts = computed(() => accountsStore.accounts)
// const baseCurrency = computed(() => settingsStore.settings.base_currency || 'JPY')

const dateRangePreset = ref('90d')
const showIncome = ref(true)
const showExpense = ref(true)
const showNet = ref(true)
const selectedAccount = ref('')
const selectedCurrency = ref('JPY')

const presets = [
  { label: '30日', value: '30d' },
  { label: '90日', value: '90d' },
  { label: '180日', value: '180d' },
  { label: '1年', value: '1y' },
  { label: 'カスタム', value: 'custom' }
]

const currencies = ['JPY', 'USD', 'EUR', 'GBP', 'CNY', 'KRW']

const dateFrom = computed(() => {
  const now = new Date()
  switch (dateRangePreset.value) {
    case '30d': now.setDate(now.getDate() - 30); break
    case '90d': now.setDate(now.getDate() - 90); break
    case '180d': now.setDate(now.getDate() - 180); break
    case '1y': now.setFullYear(now.getFullYear() - 1); break
    default: return props.initialDateFrom || ''
  }
  return now.toISOString().split('T')[0]
})

const dateTo = computed(() => {
  const now = new Date()
  if (dateRangePreset.value === 'custom') return props.initialDateTo || ''
  return now.toISOString().split('T')[0]
})

function setPreset(preset: string) {
  dateRangePreset.value = preset
  if (preset !== 'custom') {
    fetchData()
  }
}

async function fetchData() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (dateFrom.value) params.append('date_from', dateFrom.value)
    if (dateTo.value) params.append('date_to', dateTo.value)
    if (selectedAccount.value) params.append('account_id', selectedAccount.value)
    if (selectedCurrency.value) params.append('currency', selectedCurrency.value)

    const res = await api.get<CashFlowResponse>(`/analytics/cashflow/projection?${params}`)
    projectionData.value = res.data.projected_balance || []
    selectedCurrency.value = res.data.projected_balance[0]?.currency || 'JPY'
  } catch (e) {
    console.error('Failed to fetch cashflow projection:', e)
  } finally {
    loading.value = false
    await nextTick()
    updateChart()
  }
}

function updateChart() {
  if (!chartRef.value || !chart) return

  const dates = projectionData.value.map(p => p.date)
  const incomeData = projectionData.value.map(p => p.income)
  const expenseData = projectionData.value.map(p => p.expense)
  const netData = projectionData.value.map(p => p.net)
  const balanceData = projectionData.value.map(p => p.running_balance)

  const series: any[] = []

  if (showIncome.value) {
    series.push({
      name: '収入',
      type: 'bar',
      data: incomeData,
      itemStyle: { color: '#28a745' },
      yAxisIndex: 0
    })
  }
  if (showExpense.value) {
    series.push({
      name: '支出',
      type: 'bar',
      data: expenseData.map(v => -v),
      itemStyle: { color: '#dc3545' },
      yAxisIndex: 0
    })
  }
  if (showNet.value) {
    series.push({
      name: '純増減',
      type: 'line',
      data: netData,
      itemStyle: { color: '#4a90d9' },
      lineStyle: { width: 2 },
      yAxisIndex: 0
    })
  }
  if (showNet.value) {
    series.push({
      name: '予測残高',
      type: 'line',
      data: balanceData,
      itemStyle: { color: '#6ab0f3' },
      lineStyle: { width: 3, type: 'dashed' },
      yAxisIndex: 1
    })
  }

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' },
      formatter: (params: any[]) => {
        let html = params[0].axisValue + '<br/>'
        params.forEach(p => {
          const color = p.color
          let val = p.value
          if (p.seriesName === '支出') val = -val
          html += `<span style="display:inline-block;margin-right:5px;border-radius:10px;width:10px;height:10px;background-color:${color};"></span>`
          html += `${p.seriesName}: ${formatCurrency(val)}<br/>`
        })
        return html
      }
    },
    legend: {
      data: series.map(s => s.name),
      bottom: 0
    },
    grid: { left: 60, right: 40, bottom: 40, top: 20 },
    xAxis: {
      type: 'category',
      data: dates,
      boundaryGap: false,
      axisLabel: { interval: Math.max(1, Math.floor(dates.length / 10)), rotate: 30 }
    },
    yAxis: [
      {
        type: 'value',
        name: '金額',
        axisLabel: { formatter: (v: number) => formatCurrency(v) }
      },
      {
        type: 'value',
        name: '残高',
        axisLabel: { formatter: (v: number) => formatCurrency(v) },
        offset: 0
      }
    ],
    series,
    dataZoom: [
      { type: 'inside', start: 0, end: 100 },
      { type: 'slider', start: 0, end: 100, bottom: -10 }
    ]
  }

  chart.setOption(option)
}

function formatCurrency(val: number): string {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: selectedCurrency.value,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(val)
}

onMounted(async () => {
  await accountsStore.fetchAccounts()
  if (chartRef.value) {
    chart = echarts.init(chartRef.value)
    window.addEventListener('resize', () => chart?.resize())
  }
  await fetchData()
})

onUnmounted(() => {
  chart?.dispose()
  window.removeEventListener('resize', () => chart?.resize())
})

watch([dateFrom, dateTo, selectedAccount, selectedCurrency], () => {
  if (dateRangePreset.value !== 'custom') {
    fetchData()
  }
})
</script>

<style scoped>
@reference "tailwindcss";

.cashflow-chart-container {
  @apply bg-card border border-border-primary rounded-xl p-4;
}

.chart-controls {
  @apply flex flex-wrap gap-4 mb-4 items-end;
}

.control-group {
  @apply flex flex-col gap-1 min-w-[120px];
}

.control-label {
  @apply text-xs text-text-muted font-medium;
}

.preset-buttons {
  @apply flex gap-1;
}

.preset-btn {
  @apply px-3 py-1 text-sm rounded-full border border-border-primary bg-bg-card text-text-primary transition-all;
}

.preset-btn.active {
  @apply bg-accent text-white border-accent;
}

.toggle-group {
  @apply flex gap-3;
}

.toggle-item {
  @apply flex items-center gap-1 cursor-pointer;
}

.toggle-item input {
  @apply w-4 h-4 accent-accent;
}

.filter-select {
  @apply px-3 py-1 text-sm rounded border border-border-primary bg-bg-card text-text-primary;
}

.chart-wrapper {
  @apply w-full;
  height: 350px;
}

.chart-loading {
  @apply w-full;
  height: 350px;
}
</style>