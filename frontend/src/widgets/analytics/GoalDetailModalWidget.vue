<template>
  <Teleport to="body">
    <div v-if="goal" class="modal-overlay" @click.self="$emit('close')">
      <div class="modal-container">
        <div class="modal-header">
          <div class="modal-title">
            <span class="goal-icon" :style="{ backgroundColor: goal.color }">{{ goal.icon }}</span>
            <h2>{{ goal.category_name }}</h2>
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
              <span class="summary-label">目標額</span>
              <span class="summary-value">{{ formatCurrency(goal.target_amount) }}</span>
            </div>
            <div class="summary-item">
              <span class="summary-label">現在額</span>
              <span class="summary-value">{{ formatCurrency(goal.current_amount) }}</span>
            </div>
            <div class="summary-item">
              <span class="summary-label">残り</span>
              <span class="summary-value" :class="goal.remaining > 0 ? 'negative' : 'positive'">
                {{ formatCurrency(goal.remaining) }}
              </span>
            </div>
            <div class="summary-item">
              <span class="summary-label">進捗</span>
              <span class="summary-value">{{ goal.progress_pct.toFixed(1) }}%</span>
            </div>
          </div>
          <div class="projection-chart" ref="projectionChartRef" style="height: 300px;" />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { api } from '@/api/client'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([LineChart, TitleComponent, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

interface GoalItem {
  category_id: string
  category_name: string
  icon: string
  color: string
  target_amount: number
  target_period: string
  current_amount: number
  progress_pct: number
  remaining: number
  projected_completion: string
  status: 'on_track' | 'behind' | 'achieved'
}

const props = defineProps<{
  goal: GoalItem | null
  dateFrom: string
  dateTo: string
}>()

defineEmits<{
  close: []
}>()

const projectionChartRef = ref<HTMLElement | null>(null)
let projectionChartInstance: echarts.ECharts | null = null

async function renderProjectionChart(goal: GoalItem) {
  if (!projectionChartRef.value) return
  if (projectionChartInstance) projectionChartInstance.dispose()
  
  try {
    const res = await api.get<{ projection: any[] }>(
      `/analytics/goals/projection?category_id=${goal.category_id}&date_from=${props.dateFrom}&date_to=${props.dateTo}`
    )
    const projection = res.data.projection

    projectionChartInstance = echarts.init(projectionChartRef.value)
    projectionChartInstance.setOption({
      title: { text: '進捗予測', left: 'center' },
      tooltip: { trigger: 'axis' },
      legend: { data: ['実績累計', '目標累計', '目標総額'], bottom: 0 },
      grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true },
      xAxis: { type: 'category', data: projection.map(p => p.period) },
      yAxis: { type: 'value', axisLabel: { formatter: '{value}円' } },
      series: [
        {
          name: '実績累計',
          type: 'line',
          data: projection.map(p => p.actual_cumulative),
          smooth: true,
          itemStyle: { color: goal.color },
          lineStyle: { width: 3 }
        },
        {
          name: '目標累計',
          type: 'line',
          data: projection.map(p => p.target_cumulative),
          smooth: true,
          itemStyle: { color: '#666' },
          lineStyle: { type: 'dashed', width: 2 }
        },
        {
          name: '目標総額',
          type: 'line',
          data: projection.map(p => p.target_total),
          smooth: true,
          itemStyle: { color: '#999' },
          lineStyle: { type: 'dotted', width: 1 }
        }
      ]
    })
  } catch (e) {
    console.error('Failed to fetch projection:', e)
  }
}

function formatCurrency(val: number): string {
  return new Intl.NumberFormat('ja-JP', { 
    style: 'currency', 
    currency: 'JPY',
    maximumFractionDigits: 0 
  }).format(val)
}

watch(() => props.goal, async (newGoal) => {
  if (newGoal) {
    await nextTick()
    renderProjectionChart(newGoal)
  } else {
    projectionChartInstance?.dispose()
    projectionChartInstance = null
  }
})

onMounted(() => {
  if (props.goal) {
    nextTick(() => renderProjectionChart(props.goal!))
  }
})
</script>

<style scoped>
@reference "tailwindcss";

.modal-overlay {
  @apply fixed inset-0 flex items-center justify-center p-4 z-[1000];
  background-color: rgba(0, 0, 0, 0.5);
  animation: fadeIn 0.2s ease;
}
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
.modal-container {
  @apply border rounded-xl w-full max-w-[500px] max-h-[90vh] overflow-hidden;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  animation: slideUp 0.2s ease;
}
@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
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
  transition: background-color 0.2s ease, color 0.2s ease;
}
.modal-close:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}
.modal-body {
  @apply p-6;
}
.detail-summary {
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
