<template>
  <div class="cashflow-table-widget">
    <div class="widget-header">
      <h2>キャッシュフロー テーブル</h2>
      <div class="widget-controls">
        <select v-model="interval" @change="fetchData" class="interval-select">
          <option value="day">日次</option>
          <option value="week">週次</option>
          <option value="month">月次</option>
          <option value="quarter">四半期</option>
          <option value="year">年次</option>
        </select>
        <select v-model="groupBy" @change="fetchData" class="interval-select">
          <option value="none">グループなし</option>
          <option value="category">カテゴリ別</option>
          <option value="account">口座別</option>
        </select>
      </div>
    </div>

    <div class="table-container" v-if="!loading">
      <table class="cashflow-table">
        <thead>
          <tr>
            <th>期間</th>
            <th colspan="4" class="col-group">収入</th>
            <th colspan="4" class="col-group">支出</th>
            <th>純額</th>
            <th>取引数</th>
          </tr>
          <tr>
            <th></th>
            <th>件数</th>
            <th>日平均</th>
            <th>取引平均</th>
            <th>合計</th>
            <th>件数</th>
            <th>日平均</th>
            <th>取引平均</th>
            <th>合計</th>
            <th></th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <template v-for="item in tableData" :key="item.period">
            <tr @click="toggleExpanded(item.period)" class="period-row" :class="{ expanded: expandedPeriods.has(item.period) }">
              <td class="period-cell">
                <span v-if="groupBy !== 'none'" class="expand-icon">{{ expandedPeriods.has(item.period) ? '▼' : '▶' }}</span>
                {{ item.period }}
              </td>
              <td class="number">{{ item.income.count }}</td>
              <td class="number">{{ formatCurrencyShort(item.income.avg_per_day, 'income') }}</td>
              <td class="number">{{ formatCurrencyShort(item.income.avg_per_txn, 'income') }}</td>
              <td class="number income">{{ formatCurrencyShort(item.income.total, 'income') }}</td>
              <td class="number">{{ item.expense.count }}</td>
              <td class="number">{{ formatCurrencyShort(item.expense.avg_per_day, 'expense') }}</td>
              <td class="number">{{ formatCurrencyShort(item.expense.avg_per_txn, 'expense') }}</td>
              <td class="number expense">{{ formatCurrencyShort(item.expense.total, 'expense') }}</td>
              <td class="number" :class="item.net >= 0 ? 'income' : 'expense'">{{ formatCurrencyShort(item.net, item.net >= 0 ? 'income' : 'expense') }}</td>
              <td class="number">{{ item.transaction_count }}</td>
            </tr>
          </template>
        </tbody>
        <tfoot>
          <tr class="total-row">
            <td><strong>合計</strong></td>
            <td class="number"><strong>{{ totals.income_count || 0 }}</strong></td>
            <td class="number"><strong>{{ formatCurrencyShort(totals.income_avg_day || 0, 'income') }}</strong></td>
            <td class="number"><strong>{{ formatCurrencyShort(totals.income_avg_txn || 0, 'income') }}</strong></td>
            <td class="number income"><strong>{{ formatCurrencyShort(totals.income || 0, 'income') }}</strong></td>
            <td class="number"><strong>{{ totals.expense_count || 0 }}</strong></td>
            <td class="number"><strong>{{ formatCurrencyShort(totals.expense_avg_day || 0, 'expense') }}</strong></td>
            <td class="number"><strong>{{ formatCurrencyShort(totals.expense_avg_txn || 0, 'expense') }}</strong></td>
            <td class="number expense"><strong>{{ formatCurrencyShort(totals.expense || 0, 'expense') }}</strong></td>
            <td class="number" :class="totals.net >= 0 ? 'income' : 'expense'"><strong>{{ formatCurrencyShort(totals.net || 0, totals.net >= 0 ? 'income' : 'expense') }}</strong></td>
            <td class="number"><strong>{{ totals.total_txns || 0 }}</strong></td>
          </tr>
        </tfoot>
      </table>

      <div class="table-actions">
        <button @click="exportCSV" class="btn-secondary">CSVエクスポート</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watchEffect } from 'vue'
import { api } from '@/api/client'
import { useSettingsStore } from '@/stores/settings'
import { formatAmountForDisplay, formatNumber as formatBaseNumber } from '@/utils/amountDisplay'

interface PeriodData {
  period: string
  income: { count: number; avg_per_day: number; avg_per_txn: number; total: number }
  expense: { count: number; avg_per_day: number; avg_per_txn: number; total: number }
  net: number
  transaction_count: number
  children?: any[]
}

interface Totals {
  income: number
  expense: number
  net: number
  total_txns: number
  income_count?: number
  expense_count?: number
  income_avg_day?: number
  expense_avg_day?: number
  income_avg_txn?: number
  expense_avg_txn?: number
}

const props = defineProps<{ dateFrom: string; dateTo: string }>()

const interval = ref('month')
const groupBy = ref('none')
const loading = ref(true)

const tableData = ref<PeriodData[]>([])
const totals = ref<Totals>({ income: 0, expense: 0, net: 0, total_txns: 0 })
const expandedPeriods = ref<Set<string>>(new Set())
const settingsStore = useSettingsStore()

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get<{ periods: PeriodData[]; totals: Totals }>(
      `/analytics/reports/cashflow-table?date_from=${props.dateFrom}&date_to=${props.dateTo}&interval=${interval.value}&group_by=${groupBy.value}`
    )
    tableData.value = res.data.periods
    totals.value = res.data.totals
  } catch (e) {
    console.error('Failed to fetch cashflow table:', e)
  } finally {
    loading.value = false
  }
}

function toggleExpanded(period: string) {
  const set = new Set(expandedPeriods.value)
  if (set.has(period)) set.delete(period)
  else set.add(period)
  expandedPeriods.value = set
}

function formatCurrencyShort(val: number, type?: 'income' | 'expense' | 'transfer'): string {
  if (settingsStore.isHideAmounts) {
    if (type) {
      return formatAmountForDisplay(val, type)
    }
    const symbol = formatBaseNumber(0).replace(/\d/g, '')
    return `${symbol}*****`
  }
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

function exportCSV() {
  const headers = ['期間', '収入件数', '収入日平均', '収入取引平均', '収入合計', '支出件数', '支出日平均', '支出取引平均', '支出合計', '純額', '取引数']
  const rows = tableData.value.map(p => [
    p.period, p.income.count, formatCurrencyShort(p.income.avg_per_day), formatCurrencyShort(p.income.avg_per_txn), formatCurrencyShort(p.income.total),
    p.expense.count, formatCurrencyShort(p.expense.avg_per_day), formatCurrencyShort(p.expense.avg_per_txn), formatCurrencyShort(p.expense.total),
    formatCurrencyShort(p.net), p.transaction_count
  ])
  const csv = [headers, ...rows].map(r => r.map(v => `"${v}"`).join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `cashflow-table-${props.dateFrom}-${props.dateTo}.csv`
  link.click()
}

onMounted(fetchData)
watchEffect(fetchData)
</script>

<style scoped>
.cashflow-table-widget { display: flex; flex-direction: column; gap: 1rem; }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; background-color: var(--bg-secondary); border: 1px solid var(--border-primary); border-radius: var(--radius) var(--radius) 0 0; }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.widget-controls { display: flex; gap: 0.5rem; }
.interval-select { padding: 0.25rem 0.5rem; font-size: 0.875rem; border: 1px solid var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); border-radius: var(--radius-sm); }
.table-container { overflow-x: auto; border: 1px solid var(--border-primary); border-radius: 0 0 var(--radius) var(--radius); background-color: var(--bg-card); }
.cashflow-table { width: 100%; border-collapse: collapse; font-size: 0.8125rem; }
.cashflow-table th { background-color: var(--bg-secondary); padding: 0.5rem 0.75rem; text-align: left; font-weight: 600; color: var(--text-secondary); border-bottom: 2px solid var(--border-primary); white-space: nowrap; }
.cashflow-table th.col-group { text-align: center; }
.cashflow-table td { padding: 0.5rem 0.75rem; border-bottom: 1px solid var(--border-primary); color: var(--text-primary); white-space: nowrap; }
.cashflow-table td.number { text-align: right; font-family: monospace; font-size: 0.75rem; }
.cashflow-table td.income { color: var(--success); font-weight: 600; }
.cashflow-table td.expense { color: var(--error); font-weight: 600; }
.cashflow-table tbody tr:hover { background-color: var(--bg-hover); }
.period-cell { font-weight: 500; }
.total-row { background-color: var(--bg-secondary); border-top: 2px solid var(--border-primary); }
.table-actions { padding: 1rem; border-top: 1px solid var(--border-primary); display: flex; justify-content: flex-end; }
.btn-secondary { display: inline-flex; align-items: center; gap: 0.5rem; padding: 0.5rem 1rem; border: 1px solid var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); border-radius: var(--radius-sm); font-size: 0.875rem; font-weight: 500; cursor: pointer; }
.btn-secondary:hover { background-color: var(--bg-hover); }
</style>