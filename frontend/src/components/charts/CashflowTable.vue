  <template>
  <div class="cashflow-table-widget">
    <div class="widget-header">
      <h2>{{ t('reports.cashflowTable.title') }}</h2>
      <div class="widget-controls">
        <select v-model="interval" @change="fetchData" class="interval-select">
          <option value="day">{{ t('reports.cashflowTable.intervalDay') }}</option>
          <option value="week">{{ t('reports.cashflowTable.intervalWeek') }}</option>
          <option value="month">{{ t('reports.cashflowTable.intervalMonth') }}</option>
          <option value="quarter">{{ t('reports.cashflowTable.intervalQuarter') }}</option>
          <option value="year">{{ t('reports.cashflowTable.intervalYear') }}</option>
        </select>
        <select v-model="groupBy" @change="fetchData" class="interval-select">
          <option value="none">{{ t('reports.cashflowTable.groupNone') }}</option>
          <option value="category">{{ t('reports.cashflowTable.groupCategory') }}</option>
          <option value="account">{{ t('reports.cashflowTable.groupAccount') }}</option>
        </select>
      </div>
    </div>

    <div class="table-container" v-if="!loading">
      <table class="cashflow-table">
        <thead>
          <tr>
            <th>{{ t('reports.cashflowTable.period') }}</th>
            <th colspan="4" class="col-group">{{ t('reports.cashflowTable.income') }}</th>
            <th colspan="4" class="col-group">{{ t('reports.cashflowTable.expense') }}</th>
            <th>{{ t('reports.cashflowTable.net') }}</th>
            <th>{{ t('reports.cashflowTable.transactionCount') }}</th>
          </tr>
          <tr>
            <th></th>
            <th>{{ t('reports.cashflowTable.count') }}</th>
            <th>{{ t('reports.cashflowTable.avgPerDay') }}</th>
            <th>{{ t('reports.cashflowTable.avgPerTxn') }}</th>
            <th>{{ t('reports.cashflowTable.total') }}</th>
            <th>{{ t('reports.cashflowTable.count') }}</th>
            <th>{{ t('reports.cashflowTable.avgPerDay') }}</th>
            <th>{{ t('reports.cashflowTable.avgPerTxn') }}</th>
            <th>{{ t('reports.cashflowTable.total') }}</th>
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
            <td><strong>{{ t('reports.cashflowTable.total') }}</strong></td>
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
        <button @click="exportCSV" class="btn-secondary">{{ t('reports.cashflowTable.exportCsv') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import { useTransactionsStore } from '@/stores/transactions'
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
const { t } = useI18n()

async function fetchData() {
  loading.value = true
  try {
    const transactionsStore = useTransactionsStore()
    await transactionsStore.fetchTransactions(1)
    const list = transactionsStore.transactions
    const from = props.dateFrom
    const to = props.dateTo
    const filtered = list.filter((tx: any) => {
      const d = String(tx.date_time ?? '').slice(0, 10)
      if (!d) return false
      if (from && d < from) return false
      if (to && d > to) return false
      return true
    })
    const byPeriod: Record<string, { income: { total: number; count: number }; expense: { total: number; count: number } }> = {}
    let income = 0, expense = 0, count = 0
    for (const tx of filtered) {
      const d = String(tx.date_time ?? '').slice(0, 7)
      if (!byPeriod[d]) byPeriod[d] = { income: { total: 0, count: 0 }, expense: { total: 0, count: 0 } }
      if (tx.type === 'income') { byPeriod[d].income.total += tx.amount; byPeriod[d].income.count += 1; income += tx.amount }
      else if (tx.type === 'expense') { byPeriod[d].expense.total += tx.amount; byPeriod[d].expense.count += 1; expense += tx.amount }
      count += 1
    }
    const periods = Object.keys(byPeriod).sort().map((p) => {
      const inc = byPeriod[p].income
      const exp = byPeriod[p].expense
      return {
        period: p,
        income: { count: inc.count, avg_per_day: 0, avg_per_txn: inc.count ? inc.total / inc.count : 0, total: inc.total },
        expense: { count: exp.count, avg_per_day: 0, avg_per_txn: exp.count ? exp.total / exp.count : 0, total: exp.total },
        net: inc.total - exp.total,
        transaction_count: inc.count + exp.count,
        children: []
      }
    })
    tableData.value = periods
    totals.value = { income, expense, net: income - expense, total_txns: count }
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
  return val.toLocaleString(settingsStore.numberFormat)
}

function exportCSV() {
  const headers = [
    t('reports.cashflowTable.period'),
    `${t('reports.cashflowTable.income')} ${t('reports.cashflowTable.count')}`,
    `${t('reports.cashflowTable.income')} ${t('reports.cashflowTable.avgPerDay')}`,
    `${t('reports.cashflowTable.income')} ${t('reports.cashflowTable.avgPerTxn')}`,
    `${t('reports.cashflowTable.income')} ${t('reports.cashflowTable.total')}`,
    `${t('reports.cashflowTable.expense')} ${t('reports.cashflowTable.count')}`,
    `${t('reports.cashflowTable.expense')} ${t('reports.cashflowTable.avgPerDay')}`,
    `${t('reports.cashflowTable.expense')} ${t('reports.cashflowTable.avgPerTxn')}`,
    `${t('reports.cashflowTable.expense')} ${t('reports.cashflowTable.total')}`,
    t('reports.cashflowTable.net'),
    t('reports.cashflowTable.transactionCount')
  ]
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
</script>

<style scoped>
@reference "tailwindcss";
.cashflow-table-widget { @apply flex flex-col gap-4; }
.widget-header { @apply flex flex-col gap-3 p-4 sm:flex-row sm:items-center sm:justify-between border border-b-0; background-color: var(--bg-secondary); border-color: var(--border-primary); border-radius: var(--radius) var(--radius) 0 0; }
.widget-header h2 { @apply text-base sm:text-lg font-semibold; color: var(--text-primary); }
.widget-controls { @apply flex w-full flex-wrap gap-2 sm:w-auto; }
.interval-select { @apply flex-1 sm:flex-initial min-w-[120px] px-2.5 py-1.5 text-sm border rounded-sm; border-color: var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); }
.table-container { @apply w-full overflow-x-auto border; border-color: var(--border-primary); border-radius: 0 0 var(--radius) var(--radius); background-color: var(--bg-card); }
.cashflow-table { @apply w-full border-collapse text-xs; }
.cashflow-table th { @apply p-2 sm:p-3 text-left font-semibold whitespace-nowrap border-b-2; background-color: var(--bg-secondary); color: var(--text-secondary); border-color: var(--border-primary); }
.cashflow-table th.col-group { @apply text-center; }
.cashflow-table td { @apply p-2 sm:p-3 border-b whitespace-nowrap; border-color: var(--border-primary); color: var(--text-primary); }
.cashflow-table td.number { @apply text-right font-mono text-xs; }
.cashflow-table td.income { @apply font-semibold; color: var(--success); }
.cashflow-table td.expense { @apply font-semibold; color: var(--error); }
.cashflow-table tbody tr:hover { background-color: var(--bg-hover); }
.period-cell { @apply font-medium; }
.total-row { @apply border-t-2; background-color: var(--bg-secondary); border-color: var(--border-primary); }
.table-actions { @apply p-4 border-t flex justify-end; border-color: var(--border-primary); }
.btn-secondary { @apply inline-flex items-center gap-2 px-4 py-2 text-sm font-medium border rounded-sm cursor-pointer transition-colors; border-color: var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); }
.btn-secondary:hover { background-color: var(--bg-hover); }
</style>