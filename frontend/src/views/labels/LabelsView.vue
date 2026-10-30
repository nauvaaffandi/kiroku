<template>
  <div class="labels-page">
    <header class="page-header">
      <div>
        <h1>{{ t('labels.title') }}</h1>
        <p class="page-subtitle">{{ t('labels.empty') }}</p>
      </div>
      <BaseButton @click="openCreateModal">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        {{ t('labels.addLabel') }}
      </BaseButton>
    </header>

    <div v-if="labelsStore.loading" class="loading-state">
      {{ t('common.loading') }}
    </div>

    <div v-else-if="labels.length === 0" class="empty-state">
      <p>{{ t('labels.empty') }}</p>
    </div>

    <div v-else class="labels-list">
      <LabelCard
        v-for="label in labels"
        :key="label.id"
        :label="label"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </div>

    <LabelFormModal
      v-model="showFormModal"
      :label="editingLabel"
      @submit="handleFormSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useLabelsStore } from '@/stores/labels'
import { useI18n } from '@/locales'
import type { Label, CreateLabelRequest } from '@/types'
import { LabelCard, LabelFormModal } from '@/widgets/labels'
import { BaseButton } from '@/components/ui'

const labelsStore = useLabelsStore()
const { t } = useI18n()

const labels = computed(() => labelsStore.labels)
const showFormModal = ref(false)
const editingLabel = ref<Label | null>(null)

const fetchLabels = async () => {
  await labelsStore.fetchLabels()
}

onMounted(() => {
  fetchLabels()
})

const openCreateModal = () => {
  editingLabel.value = null
  showFormModal.value = true
}

const handleEdit = (label: Label) => {
  editingLabel.value = label
  showFormModal.value = true
}

const handleDelete = async (id: string) => {
  if (confirm('Delete this label?')) {
    await labelsStore.deleteLabel(id)
  }
}

const handleFormSubmit = async (data: CreateLabelRequest) => {
  if (editingLabel.value) {
    await labelsStore.updateLabel(editingLabel.value.id, data)
  } else {
    await labelsStore.createLabel(data)
  }
  showFormModal.value = false
  editingLabel.value = null
  await fetchLabels()
}
</script>

<style scoped>
@reference "tailwindcss";

.labels-page {
  @apply max-w-4xl mx-auto px-4 py-6;
}

.page-header {
  @apply flex justify-between items-start mb-6 gap-4;
}

.page-header h1 {
  @apply text-2xl font-bold m-0;
  color: var(--text-primary);
}

.page-subtitle {
  @apply text-sm mt-1 mb-0;
  color: var(--text-muted);
}

.loading-state,
.empty-state {
  @apply text-center py-12;
  color: var(--text-muted);
}
</style>
