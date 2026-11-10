<template>
  <div class="goals-view">
    <div class="page-header">
      <div class="header-left">
        <h1>{{ t('goals.title') }}</h1>
        <p class="page-subtitle">{{ t('goals.subtitle') }}</p>
      </div>
      <div class="header-actions">
        <button class="btn-primary" @click="openCreateDialog">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          {{ t('goals.addGoal') }}
        </button>
      </div>
    </div>

    <div class="goals-grid" v-if="!loading && goals.length > 0">
      <GoalCardWidget
        v-for="goal in goals"
        :key="goal.id"
        :goal="goal"
        @click="openDetail"
      />
    </div>

    <div class="goals-empty" v-if="!loading && goals.length === 0">
      <div class="empty-state">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M12 2L2 7l10 5 10-5-10-5z"></path>
          <path d="M2 17l10 5 10-5"></path>
          <path d="M2 12l10 5 10-5"></path>
        </svg>
        <h3>{{ t('goals.emptyTitle') }}</h3>
        <p>{{ t('goals.emptyDesc') }}</p>
        <button class="btn-primary" @click="openCreateDialog">{{ t('goals.setGoal') }}</button>
      </div>
    </div>

    <div v-else-if="loading" class="goals-grid">
      <div class="skeleton skeleton-card" v-for="i in 4" :key="i" />
    </div>

    <GoalDetailModalWidget
      :goal="selectedGoal"
      @close="closeDetail"
    />

    <GoalCreateModal
      :visible="showCreateDialog"
      @close="showCreateDialog = false"
      @created="onGoalCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect } from 'vue'
import { api } from '@/api/client'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'
import GoalCardWidget from '@/widgets/goals/GoalCardWidget.vue'
import GoalDetailModalWidget from '@/widgets/goals/GoalDetailModalWidget.vue'
import GoalCreateModal from '@/widgets/goals/GoalCreateModal.vue'
import type { Goal } from '@/types/goal'

interface GoalItem extends Goal {
  status: 'on_track' | 'behind' | 'achieved'
  progress_pct: number
  remaining: number
  color: string
}

const uiStore = useUIStore()
const { t } = useI18n()
const loading = ref(true)
const showCreateDialog = ref(false)
const selectedGoal = ref<Goal | null>(null)
const goals = ref<GoalItem[]>([])

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get<GoalItem[]>('/api/v1/goals/')
    goals.value = res.data
  } catch (e) {
    console.error('Failed to fetch goals:', e)
    uiStore.error(t('goals.fetchError'))
  } finally {
    loading.value = false
  }
}

function openDetail(goal: Goal) {
  selectedGoal.value = goal
}

function closeDetail() {
  selectedGoal.value = null
}

function openCreateDialog() {
  showCreateDialog.value = true
}

function onGoalCreated() {
  showCreateDialog.value = false
  fetchData()
}

watchEffect(() => {
  fetchData()
})
</script>

<style scoped>
@reference "tailwindcss";

.goals-view {
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
  border: none;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.btn-primary:hover {
  @apply opacity-90;
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