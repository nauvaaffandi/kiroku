<template>
  <Teleport to="body">
    <div v-if="goal" class="modal-overlay" @click.self="$emit('close')">
      <div class="modal-container">
        <div class="modal-header">
          <div class="modal-title">
            <span class="goal-icon" :style="{ backgroundColor: goal.color_theme }">{{ goal.icon }}</span>
            <h2>{{ goal.name }}</h2>
          </div>
          <button @click="$emit('close')" class="modal-close" aria-label="Close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="detail-summary">
            <div class="summary-item">
              <span class="summary-label">{{ t('goals.targetAmount') }}</span>
              <span class="summary-value">{{ formatCurrency(goal.target_amount) }}</span>
            </div>
            <div class="summary-item">
              <span class="summary-label">{{ t('goals.currentAmount') }}</span>
              <span class="summary-value">{{ formatCurrency(goal.saved_already) }}</span>
            </div>
            <div class="summary-item">
              <span class="summary-label">{{ t('goals.remaining') }}</span>
              <span class="summary-value" :class="remaining > 0 ? 'negative' : 'positive'">
                {{ formatCurrency(remaining) }}
              </span>
            </div>
            <div class="summary-item">
              <span class="summary-label">{{ t('goals.progress') }}</span>
              <span class="summary-value">{{ progress.toFixed(1) }}%</span>
            </div>
          </div>
          <div class="projection-summary" v-if="projection">
            <div class="summary-item">
              <span class="summary-label">{{ t('goals.weeklyRate') }}</span>
              <span class="summary-value">{{ formatCurrency(projection.weekly_rate) }}</span>
            </div>
            <div class="summary-item">
              <span class="summary-label">{{ t('goals.estimatedWeeks') }}</span>
              <span class="summary-value">{{ projection.estimated_weeks }}</span>
            </div>
            <div class="summary-item" v-if="projection.estimated_reach_date">
              <span class="summary-label">{{ t('goals.projectedCompletion') }}</span>
              <span class="summary-value">{{ projection.estimated_reach_date }}</span>
            </div>
          </div>
          <div class="projection-chart" ref="projectionChartRef" style="height: 300px;" />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed, onBeforeUnmount } from 'vue'
import { api } from '@/api/client'
import { useSettingsStore } from '@/stores/settings'
import { useExchangeRate } from '@/composables/useExchangeRate'
import { getCurrencySymbol } from '@/utils/amountDisplay'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { useI18n } from '@/locales'
import type { Goal, GoalProjection } from '@/types/goal'

echarts.use([LineChart, TitleComponent, TooltipComponent, GridComponent, CanvasRenderer])

const { t } = useI18n()
const settingsStore = useSettingsStore()
const { toUsd } = useExchangeRate()
const baseCurrency = computed(() => settingsStore.baseCurrency || 'USD')
const projection = ref<GoalProjection | null>(null)
const projectionChartRef = ref<HTMLElement | null>(null)
let projectionChartInstance: echarts.ECharts | null = null

function formatCurrency(val: number): string {
  const usd = toUsd(val)
  const symbol = getCurrencySymbol(baseCurrency.value)
  const locale = settingsStore.numberFormat || 'en-US'
  return `${symbol}${new Intl.NumberFormat(locale).format(Math.round(usd))}`
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

const props = defineProps<{
  goal: Goal | null
}>()

defineEmits<{
  close: []
}>()

const remaining = computed(() => Math.max(props.goal ? props.goal.target_amount - props.goal.saved_already : 0, 0))
const progress = computed(() => {
  if (!props.goal || props.goal.target_amount <= 0) return 0
  return (props.goal.saved_already / props.goal.target_amount) * 100
})

async function renderProjectionChart(goal: Goal) {
  projection.value = null
  projectionChartInstance?.dispose()
  projectionChartInstance = null

  try {
    const res = await api.get<GoalProjection>(`/api/v1/goals/${goal.id}/projection`)
    projection.value = res.data
    if (!projectionChartRef.value) return

    projectionChartInstance = echarts.init(projectionChartRef.value)
    projectionChartInstance.setOption({
      title: { text: t('goals.predictionTitle'), left: 'center' },
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true },
      xAxis: { type: 'category', data: [t('goals.currentAmount'), t('goals.targetAmount')] },
      yAxis: { type: 'value', axisLabel: { formatter: (v: number) => `${formatCurrencyShort(v)} ${goal.currency}` } },
      series: [{
        name: t('goals.progress'),
        type: 'line',
        data: [res.data.current_saved, res.data.target_amount],
        smooth: true,
        itemStyle: { color: goal.color_theme || 'var(--accent)' },
        lineStyle: { width: 3 }
      }]
    })
  } catch (e) {
    console.error('Failed to fetch projection:', e)
  }
}

watch(() => props.goal, (newGoal) => {
  if (newGoal) {
    renderProjectionChart(newGoal)
  } else {
    projection.value = null
    projectionChartInstance?.dispose()
    projectionChartInstance = null
  }
})

onBeforeUnmount(() => {
  projectionChartInstance?.dispose()
})
</script>

<style scoped>
@reference "tailwindcss";

.modal-overlay {
  @apply fixed inset-0 flex items-center justify-center p-4 z-[1000];
  background-color: rgba(0, 0, 0, 0.5);
}
.modal-container {
  @apply border rounded-xl w-full max-w-[500px] max-h-[90vh] overflow-hidden;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}
.modal-header {
  @apply flex items-center justify-between py-4 px-6 border-b;
  border-color: var(--border-primary);
}
.modal-title {
  @apply flex items-center gap-3;
}
.modal-title h2 {
  @apply text-lg font-semibold;
  color: var(--text-primary);
}
.goal-icon {
  @apply w-[2.5rem] h-[2.5rem] rounded-xl flex items-center justify-center text-xl;
}
.modal-close {
  @apply p-2 rounded-sm bg-transparent border-none;
  color: var(--text-secondary);
}
.modal-body {
  @apply p-6;
}
.detail-summary,
.projection-summary {
  @apply grid grid-cols-2 gap-4 mb-6;
}
.summary-item {
  @apply flex flex-col gap-1;
}
.summary-label {
  @apply text-xs uppercase tracking-wider;
  color: var(--text-secondary);
}
.summary-value {
  @apply text-base font-semibold;
  color: var(--text-primary);
}
.summary-value.positive {
  color: var(--success);
}
.summary-value.negative {
  color: var(--error);
}
.projection-chart {
  @apply w-full;
}
</style>