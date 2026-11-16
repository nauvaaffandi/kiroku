<template>
  <div class="scheduled-page">
    <header class="page-header">
      <div>
        <h1>{{ t('scheduled.title') }}</h1>
        <p class="page-subtitle">{{ t('scheduled.subtitle') }}</p>
      </div>
      <BaseButton @click="openCreateModal">
        <Plus :size="20" stroke-width="1.75" aria-hidden="true" />
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
import type { ScheduledTransaction, ScheduledFilters as ScheduledFiltersType } from '@/types'
import { ScheduledList, ScheduledFormModal, ScheduledFilters } from '@/widgets/scheduled'
import { BaseButton } from '@/components/ui'
import { Plus } from '@lucide/vue'

const scheduledStore = useScheduledStore()
const { t } = useI18n()

const scheduled = computed(() => scheduledStore.scheduled)

const showFormModal = ref(false)
const editingScheduled = ref<ScheduledTransaction | null>(null)

const localFilters = ref<ScheduledFiltersType>({
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
  console.log('[ScheduledView] fetching scheduled', { page })
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
  if (confirm(t('scheduled.deleteConfirm'))) {
    if (import.meta.env.DEV) console.log('[ScheduledView] deleting scheduled', { id })
    scheduledStore.deleteScheduled(id)
  }
}

const openCreateModal = () => {
  editingScheduled.value = null
  showFormModal.value = true
}

const handleFormSubmit = async (data: Partial<ScheduledTransaction>) => {
  if (import.meta.env.DEV) console.log('[ScheduledView] saving scheduled', { id: editingScheduled.value?.id, data })
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
@reference "tailwindcss";
.page-header {
  @apply flex justify-between items-start mb-6 gap-4;
}

.page-header svg,
.page-header [class*="lucide-icon"],
.page-header i {
  @apply mr-2;
}

.page-header h1 {
  @apply text-[20px] font-semibold m-0;
  color: var(--text-primary);
}

.page-subtitle {
  @apply text-sm mt-1 mb-0;
  color: var(--text-muted);
}

.header-actions {
  @apply flex gap-3 shrink-0;
}
</style>
