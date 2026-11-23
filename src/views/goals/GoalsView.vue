<template>
  <div class="goals-view">
    <div class="page-header">
      <div class="header-left">
        <h1>{{ t('goals.title') }}</h1>
        <p class="page-subtitle">{{ t('goals.subtitle') }}</p>
      </div>
      <div class="header-actions">
        <BaseButton @click="openCreateDialog">
          <Plus :size="16" stroke-width="1.75" aria-hidden="true" />
          {{ t('goals.addGoal') }}
        </BaseButton>
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

    <div v-if="!loading && goals.length === 0" class="goals-empty">
      <EmptyState
        :message="t('goals.emptyTitle')"
        :action="t('goals.setGoal')"
        @action="openCreateDialog"
      />
    </div>

    <div v-else-if="loading" class="goals-grid">
      <div class="skeleton skeleton-card" v-for="i in 4" :key="i"></div>
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
import { goalRepository } from '@/repository'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'
import GoalCardWidget from '@/widgets/goals/GoalCardWidget.vue'
import GoalDetailModalWidget from '@/widgets/goals/GoalDetailModalWidget.vue'
import GoalCreateModal from '@/widgets/goals/GoalCreateModal.vue'
import { EmptyState } from '@/components/ui'
import { Plus } from '@lucide/vue'
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
  console.log('[GoalsView] fetching goals')
  try {
    const rows = await goalRepository.findAll()
    goals.value = rows.map((g: any) => {
      const target = g.target_amount ?? 0
      const saved = g.saved_already ?? 0
      const progress_pct = target ? (saved / target) * 100 : 0
      const remaining = Math.max(0, target - saved)
      return {
        id: g.id,
        name: g.name,
        target_amount: target,
        saved_already: saved,
        currency: g.currency,
        desired_date: g.desired_date,
        icon: g.icon ?? '',
        color: g.color_theme ?? '',
        created_at: g.created_at ?? '',
        updated_at: g.updated_at ?? '',
        status: progress_pct >= 100 ? 'achieved' : progress_pct > 50 ? 'on_track' : 'behind',
        progress_pct,
        remaining,
      }
    })
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
  @apply flex flex-col gap-4 p-4 w-full max-w-[480px] mx-auto;
}

.page-header {
  @apply flex items-start justify-between pb-4 flex-wrap gap-4 w-full;
  border-bottom: 1px solid var(--border-primary);
}

.header-left h1 {
  @apply text-[20px] font-semibold;
  color: var(--text-primary);
}

.page-subtitle {
  @apply text-[0.875rem] mt-1;
  color: var(--text-secondary);
}

.header-actions {
  @apply flex gap-2 items-center;
}

.goals-grid {
  @apply grid gap-4;
  grid-template-columns: minmax(0, 1fr);
}

.goals-empty {
  @apply col-span-full;
}

.skeleton-card {
  @apply h-[300px];
}
</style>
