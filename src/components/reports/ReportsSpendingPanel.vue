<template>
  <section class="chart-section">
    <div class="chart-header">
      <h2 class="chart-title">{{ t('reports.spendingByCategory') }}</h2>
      <span class="period-badge">{{ period }}</span>
    </div>
    <div class="spending-chart" ref="el"></div>
    <div class="table-wrap" v-if="groups">
      <table class="data-table">
        <thead>
          <tr class="table-header-row">
            <th class="table-head">{{ t('common.category') }}</th>
            <th class="table-head table-cell--numeric">{{ t('common.amount') }}</th>
            <th class="table-head table-cell--numeric">{{ t('reports.cashflowTable.count') }}</th>
            <th class="table-head table-cell--numeric">{{ t('reports.percent') }}</th>
          </tr>
        </thead>
        <tbody class="table-body">
          <tr v-for="group in groups" :key="group.id" class="table-row">
            <td class="table-cell">
              <div class="category-row">
                <span class="color-dot" :style="{ backgroundColor: group.color }"></span>
                <span class="category-name">{{ group.name }}</span>
              </div>
              <div class="children-list" v-if="group.children && group.children.length > 0">
                <div v-for="child in group.children" :key="child.id" class="child-row">
                  <span class="color-dot color-dot--sm" :style="{ backgroundColor: child.color }"></span>
                  <span>{{ child.name }}: {{ formatDisplayNumber(child.amount, 'expense') }}</span>
                </div>
              </div>
            </td>
            <td class="table-cell table-cell--amount">{{ formatDisplayNumber(group.amount, 'expense') }}</td>
            <td class="table-cell table-cell--count">{{ group.count }}</td>
            <td class="table-cell table-cell--pct">{{ totalAmount > 0 ? ((group.amount / totalAmount) * 100).toFixed(1) : 0 }}%</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { useI18n } from '@/locales'

echarts.use([BarChart, TitleComponent, TooltipComponent, GridComponent, CanvasRenderer])

const props = defineProps<{
  period: string
  groups: any[]
  totalAmount: number
  formatDisplayNumber: (num: number, type?: 'income' | 'expense' | 'transfer') => string
}>()

const el = ref<HTMLElement | null>(null)
const { t } = useI18n()
let instance: echarts.ECharts | null = null

const draw = () => {
  if (!el.value) return
  if (instance) instance.dispose()
  instance = echarts.init(el.value)
  const data = props.groups.slice(0, 10).reverse()
  instance.setOption({
    title: { text: t('reports.topCategories'), left: 'center', textStyle: { fontSize: 13, color: 'inherit' } },
    tooltip: { trigger: 'item', formatter: '{b}: ¥{c} ({d}%)' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'value', axisLabel: { formatter: '¥{value}' } },
    yAxis: { type: 'category', data: data.map((d: any) => d.name), inverse: true },
    series: [{ type: 'bar', data: data.map((d: any) => ({ value: d.amount, itemStyle: { color: d.color } })), label: { show: true, position: 'right', formatter: '¥{c}' } }]
  })
}

watch(() => props.groups, draw, { deep: true })
onMounted(draw)
window.addEventListener('resize', () => instance?.resize())
</script>

<style scoped>
@reference "tailwindcss";

.chart-section { @apply bg-zinc-50 dark:bg-zinc-800/50 rounded-xl p-4 border border-zinc-200 dark:border-zinc-700 mb-4; }
.chart-header { @apply flex items-center justify-between mb-4 flex-wrap gap-2; }
.chart-title { @apply text-base font-semibold text-zinc-900 dark:text-zinc-100; }
.period-badge { @apply text-xs text-zinc-600 dark:text-zinc-400 bg-white dark:bg-zinc-800 px-2.5 py-1 rounded-full border border-zinc-200 dark:border-zinc-700; }
.table-wrap { @apply overflow-x-auto; }
.data-table { @apply w-full text-left text-sm border-collapse; }
.table-header-row { @apply border-b border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 bg-white dark:bg-zinc-800; }
.table-head { @apply p-3 font-semibold; }
.table-cell { @apply p-3; }
.table-cell--numeric { @apply text-right; }
.table-cell--amount { @apply text-right font-medium text-zinc-900 dark:text-zinc-100; }
.table-cell--count, .table-cell--pct { @apply text-right text-zinc-600 dark:text-zinc-400; }
.table-body { @apply divide-y divide-zinc-200 dark:divide-zinc-700; }
.table-row { @apply hover:bg-zinc-100 dark:hover:bg-zinc-800/50; }
.category-row { @apply flex items-center gap-2; }
.color-dot { @apply w-2.5 h-2.5 rounded-full flex-shrink-0; }
.color-dot--sm { @apply w-2 h-2 rounded-full flex-shrink-0; }
.category-name { @apply font-medium text-zinc-900 dark:text-zinc-100; }
.children-list { @apply ml-4 pl-2 border-l-2 border-zinc-200 dark:border-zinc-700 mt-1 space-y-1; }
.child-row { @apply flex items-center gap-2 text-xs text-zinc-600 dark:text-zinc-400 py-0.5; }
.spending-chart { @apply w-full mb-4; height: 300px; }
@media (min-width: 640px) { .spending-chart { height: 350px; } }
</style>
