<template>
  <div class="scheduled-page">
    <header class="page-header">
      <div>
        <h1>{{ t('scheduled.title') }}</h1>
        <p class="page-subtitle">{{ t('scheduled.subtitle') }}</p>
      </div>
      <BaseButton @click="openCreateModal">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        {{ t('scheduled.addScheduled') }}
      </BaseButton>
    </header>

    <ScheduledFilters
      v-model="localFilters"
      :searchQuery="searchQuery"
      @change="handleFilterChange"
      @search="handleSearch"
    />

    <ScheduledList
      :scheduled="scheduled"
      @edit="handleEdit"
      @toggle="handleToggle"
      @delete="handleDelete"
      @create="openCreateModal"
    />

    <ScheduledFormModal
      v-model="showFormModal"
      :scheduled="editingScheduled"
      @submit="handleFormSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useScheduledStore } from '@/stores/scheduled'
import { useI18n } from '@/locales'
import type { ScheduledTransaction, ScheduledFilters } from '@/types'
import { ScheduledList, ScheduledFormModal } from '@/widgets/scheduled'
import { BaseButton } from '@/components/ui'

const scheduledStore = useScheduledStore()
const { t } = useI18n()

const scheduled = computed(() => scheduledStore.scheduled)

const showFormModal = ref(false)
const editingScheduled = ref<ScheduledTransaction | null>(null)

const localFilters = ref<ScheduledFilters>({
  q: '',
  type: '',
  frequency: '',
  active: '',
  sort: 'next_run',
  order: 'ASC'
})

const searchQuery = ref('')

let searchTimeout: ReturnType<typeof setTimeout>
watch(searchQuery, (value) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    scheduledStore.filters.q = value
    fetchScheduled(1)
  }, 300)
})

const fetchScheduled = async (page = 1) => {
  await scheduledStore.fetchScheduled(page)
}

onMounted(() => {
  fetchScheduled(1)
})

const handleFilterChange = () => {
  scheduledStore.filters = { ...localFilters.value }
  fetchScheduled(1)
}

const handleSearch = (value: string) => {
  searchQuery.value = value
}

const handleEdit = (s: ScheduledTransaction) => {
  editingScheduled.value = s
  showFormModal.value = true
}

const handleToggle = (s: ScheduledTransaction) => {
  scheduledStore.toggleScheduled(s.id)
}

const handleDelete = (id: string) => {
  if (confirm('Are you sure you want to delete this scheduled transaction?')) {
    scheduledStore.deleteScheduled(id)
  }
}

const openCreateModal = () => {
  editingScheduled.value = null
  showFormModal.value = true
}

const handleFormSubmit = async (data: Partial<ScheduledTransaction>) => {
  if (editingScheduled.value) {
    await scheduledStore.updateScheduled(editingScheduled.value.id, data)
  } else {
    await scheduledStore.createScheduled(data)
  }
  showFormModal.value = false
  editingScheduled.value = null
  await fetchScheduled(1)
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  gap: 1rem;
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.page-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin: 0.25rem 0 0;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
  flex-shrink: 0;
}
</style>