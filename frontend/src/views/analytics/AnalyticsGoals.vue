<template>
  <div class="analytics-goals-view">
    <div class="page-header">
      <div class="header-left">
        <h1>目標管理</h1>
        <p class="page-subtitle">「Want」カテゴリの予算目標と進捗</p>
      </div>
      <div class="header-actions">
        <router-link to="/categories" class="btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          目標を追加
        </router-link>
      </div>
    </div>

    <DateFilterBar @change="onDateChange" />

    <div class="controls-bar">
      <select v-model="statusFilter" class="filter-select">
        <option value="all">すべて</option>
        <option value="on_track">順調</option>
        <option value="behind">遅れ</option>
        <option value="achieved">達成</option>
      </select>
    </div>

    <GoalsSummaryCardsWidget :loading="loading" :summary="summary" />

    <div class="goals-grid" v-if="!loading && filteredGoals.length > 0">
      <GoalCardWidget
        v-for="goal in filteredGoals"
        :key="goal.category_id"
        :goal="goal"
        @click="openDetail"
      />
    </div>

    <div class="goals-empty" v-if="!loading && filteredGoals.length === 0">
      <div class="empty-state">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M12 2L2 7l10 5 10-5-10-5z"></path>
          <path d="M2 17l10 5 10-5"></path>
          <path d="M2 12l10 5 10-5"></path>
        </svg>
        <h3>目標が設定されていません</h3>
        <p>カテゴリ設定で「Want」種類かつ予算額を設定すると、ここに目標として表示されます</p>
        <router-link to="/categories" class="btn-primary">カテゴリを設定する</router-link>
      </div>
    </div>

    <div v-else-if="loading" class="goals-grid">
      <div class="skeleton skeleton-card" v-for="i in 4" :key="i" />
    </div>

    <GoalDetailModalWidget
      :goal="selectedGoal"
      :dateFrom="dateFrom"
      :dateTo="dateTo"
      @close="closeDetail"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watchEffect, computed } from 'vue'
import { api } from '@/api/client'
import { useUIStore } from '@/stores/ui'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import GoalsSummaryCardsWidget from '@/widgets/analytics/GoalsSummaryCardsWidget.vue'
import GoalCardWidget from '@/widgets/analytics/GoalCardWidget.vue'
import GoalDetailModalWidget from '@/widgets/analytics/GoalDetailModalWidget.vue'
import type { GoalItem, GoalSummary } from '@/types/report'

const props = defineProps<{
  dateFrom: string
  dateTo: string
}>()

const uiStore = useUIStore()
const loading = ref(true)
const statusFilter = ref('all')
const selectedGoal = ref<GoalItem | null>(null)

const goals = ref<GoalItem[]>([])
const summary = ref<GoalSummary>({
  total_targets: 0,
  total_current: 0,
  overall_progress: 0,
  goals_achieved: 0,
  goals_on_track: 0,
  goals_behind: 0,
})

const filteredGoals = computed(() => {
  if (statusFilter.value === 'all') return goals.value
  return goals.value.filter(g => g.status === statusFilter.value)
})

function onDateChange(_payload: { dateFrom: string; dateTo: string }) {
  // Parent handles date changes
}

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get<{ goals: GoalItem[]; summary: GoalSummary }>(
      `/analytics/goals/progress?date_from=${props.dateFrom}&date_to=${props.dateTo}`
    )
    goals.value = res.data.goals
    summary.value = res.data.summary
  } catch (e) {
    console.error('Failed to fetch goals:', e)
    uiStore.error('目標の取得に失敗しました')
  } finally {
    loading.value = false
  }
}

function openDetail(goal: GoalItem) {
  selectedGoal.value = goal
}

function closeDetail() {
  selectedGoal.value = null
}

onMounted(() => {
  fetchData()
})

watchEffect(() => {
  fetchData()
})
</script>

<style scoped>
@reference "tailwindcss";
.analytics-goals-view {
  @apply flex flex-col gap-4 p-6;
}

.page-header {
  @apply flex items-start justify-between pb-4 flex-wrap gap-4 w-full;
  border-bottom: 1px solid var(--border-primary);
}

.header-left h1 {
  @apply text-xl font-semibold;
  color: var(--text-primary);
}

.page-subtitle {
  @apply text-[0.875rem] mt-1;
  color: var(--text-secondary);
}

.header-actions {
  @apply flex gap-2 items-center;
}

.btn-primary {
  @apply inline-flex items-center gap-2 px-4 py-2 rounded-sm font-medium text-sm no-underline;
  background-color: var(--accent);
  color: var(--text-inverse);
  transition: opacity 0.2s ease;
}

.btn-primary:hover {
  @apply opacity-90;
}

.controls-bar {
  @apply flex gap-4 items-center p-4 rounded-xl;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
}

.filter-select {
  @apply px-4 py-2 rounded-sm text-sm;
  border: 1px solid var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}

.goals-grid {
  @apply grid gap-4;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
}

.goals-empty {
  @apply col-span-full;
}

.empty-state {
  @apply text-center p-12 rounded-xl;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  color: var(--text-secondary);
}

.empty-state svg {
  @apply mb-4;
  color: var(--text-muted);
}

.empty-state h3 {
  @apply text-lg font-semibold mb-2;
  color: var(--text-primary);
}

.empty-state p {
  @apply text-sm mb-6 max-w-[300px] mx-auto;
}

.skeleton-card {
  @apply h-[300px];
}
</style>