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
@reference "tailwindcss";

.goal-card {
  @apply border rounded-xl p-[1.25rem] cursor-pointer;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}
.goal-card:hover {
  @apply -translate-y-[2px];
  box-shadow: var(--shadow-lg);
  border-color: var(--accent);
}
.goal-header {
  @apply flex items-center justify-between mb-3;
}
.goal-icon {
  @apply w-[2.5rem] h-[2.5rem] rounded-xl flex items-center justify-center text-xl;
}
.goal-status {
  @apply px-[0.625rem] py-[0.25rem] rounded-full text-[0.6875rem] font-semibold uppercase;
}
.goal-status.on_track {
  background-color: var(--success-light);
  color: var(--success);
}
.goal-status.behind {
  background-color: var(--warning-light);
  color: var(--warning);
}
.goal-status.achieved {
  background-color: var(--accent-light);
  color: var(--accent);
}
.goal-name {
  @apply text-base font-semibold mb-4;
  color: var(--text-primary);
}
.goal-progress-ring {
  @apply relative w-[100px] h-[100px] mx-auto mb-4;
}
.progress-svg {
  @apply w-full h-full;
}
.progress-fill {
  transition: stroke-dashoffset 0.5s ease;
}
.progress-text {
  @apply absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-center;
  pointer-events: none;
}
.progress-pct {
  @apply block text-xl font-bold;
  color: var(--text-primary);
}
.progress-amount {
  @apply block text-[0.6875rem];
  color: var(--text-secondary);
}
.goal-details {
  @apply flex flex-col gap-[0.375rem] pt-3 border-t;
  border-color: var(--border-primary);
}
.detail-row {
  @apply flex justify-between text-xs;
}
.detail-label {
  color: var(--text-secondary);
}
.detail-value {
  @apply font-medium;
  color: var(--text-primary);
}
</style>
