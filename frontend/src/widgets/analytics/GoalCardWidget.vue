<template>
  <div class="goal-card" @click="$emit('click', goal)">
    <div class="goal-header">
      <span class="goal-icon" :style="{ backgroundColor: goal.color }">{{ goal.icon }}</span>
      <span class="goal-status" :class="goal.status">{{ statusLabel(goal.status) }}</span>
    </div>
    <div class="goal-name">{{ goal.category_name }}</div>
    <div class="goal-progress-ring" :style="progressRingStyle(goal.progress_pct)">
      <svg class="progress-svg" viewBox="0 0 100 100">
        <circle 
          class="progress-bg" 
          cx="50" cy="50" r="45" 
          fill="none" stroke="var(--bg-tertiary)" stroke-width="8"
        />
        <circle 
          class="progress-fill" 
          cx="50" cy="50" r="45" 
          fill="none" :stroke="goal.color" stroke-width="8" 
          :stroke-dasharray="circumference" 
          :stroke-dashoffset="strokeOffset(goal.progress_pct)"
          stroke-linecap="round"
          transform="rotate(-90 50 50)"
        />
      </svg>
      <div class="progress-text">
        <span class="progress-pct">{{ goal.progress_pct.toFixed(0) }}%</span>
        <span class="progress-amount">{{ formatCurrencyShort(goal.current_amount) }} / {{ formatCurrencyShort(goal.target_amount) }}</span>
      </div>
    </div>
    <div class="goal-details">
      <div class="detail-row">
        <span class="detail-label">残り</span>
        <span class="detail-value">{{ formatCurrencyShort(goal.remaining) }}</span>
      </div>
      <div class="detail-row" v-if="goal.projected_completion">
        <span class="detail-label">予測完了</span>
        <span class="detail-value">{{ goal.projected_completion }}</span>
      </div>
      <div class="detail-row" v-else>
        <span class="detail-label">期間</span>
        <span class="detail-value">{{ goal.target_period }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
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

defineProps<{
  goal: GoalItem
}>()

defineEmits<{
  click: [goal: GoalItem]
}>()

function statusLabel(status: string): string {
  const labels: Record<string, string> = {
    on_track: '順調',
    behind: '遅れ',
    achieved: '達成'
  }
  return labels[status] || status
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

function progressRingStyle(progressPct: number) {
  return {
    '--progress-pct': progressPct / 100
  }
}

const circumference = 2 * Math.PI * 45

function strokeOffset(progressPct: number): number {
  return circumference * (1 - Math.min(progressPct, 100) / 100)
}
</script>

<style scoped>
.goal-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); padding: 1.25rem; cursor: pointer; transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease; }
.goal-card:hover { transform: translateY(-2px); box-shadow: var(--shadow-lg); border-color: var(--accent); }
.goal-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.75rem; }
.goal-icon { width: 2.5rem; height: 2.5rem; border-radius: var(--radius); display: flex; align-items: center; justify-content: center; font-size: 1.25rem; }
.goal-status { padding: 0.25rem 0.625rem; border-radius: 9999px; font-size: 0.6875rem; font-weight: 600; text-transform: uppercase; }
.goal-status.on_track { background-color: var(--success-bg); color: var(--success); }
.goal-status.behind { background-color: var(--warning-bg); color: var(--warning); }
.goal-status.achieved { background-color: var(--accent-bg); color: var(--accent); }
.goal-name { font-size: 1rem; font-weight: 600; color: var(--text-primary); margin-bottom: 1rem; }
.goal-progress-ring { position: relative; width: 100px; height: 100px; margin: 0 auto 1rem; }
.progress-svg { width: 100%; height: 100%; }
.progress-fill { transition: stroke-dashoffset 0.5s ease; }
.progress-text { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; pointer-events: none; }
.progress-pct { display: block; font-size: 1.25rem; font-weight: 700; color: var(--text-primary); }
.progress-amount { display: block; font-size: 0.6875rem; color: var(--text-secondary); }
.goal-details { display: flex; flex-direction: column; gap: 0.375rem; padding-top: 0.75rem; border-top: 1px solid var(--border-primary); }
.detail-row { display: flex; justify-content: space-between; font-size: 0.75rem; }
.detail-label { color: var(--text-secondary); }
.detail-value { font-weight: 500; color: var(--text-primary); }
</style>