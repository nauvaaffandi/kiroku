<template>
  <div class="sub-categories-page">
    <header class="page-header">
      <button class="back-button" @click="goBack" :aria-label="t('common.back')">
        <ChevronLeft :size="20" />
      </button>
    </header>

    <section class="parent-action-bar">
      <h1 class="parent-title">{{ t('categories.subcategoriesTitle', { name: parentName }) }}</h1>
      <div class="parent-actions">
        <BaseButton size="sm" @click="openCreateModal">
          <Plus :size="16" />
          {{ t('categories.addChild') }}
        </BaseButton>
        <button type="button" class="action-btn" @click="handleEdit(parentCategory!)">
          <Pencil :size="16" />
        </button>
        <button type="button" class="action-btn danger" @click="handleDelete(parentCategory!.id)">
          <Trash2 :size="16" />
        </button>
      </div>
    </section>

    <div class="sub-categories-list">
      <CategoryChildRow
        v-for="sub in subs"
        :key="sub.id"
        :category="sub"
        @edit="handleEdit"
        @delete="handleDelete"
      />

      <div v-if="subs.length === 0" class="empty-state">
        <p>{{ t('categories.noCategories') }}</p>
        <p class="empty-desc">{{ t('categories.noCategoriesDescription') }}</p>
        <BaseButton size="sm" @click="openCreateModal">
          {{ t('categories.addChild') }}
        </BaseButton>
      </div>
    </div>

    <CategoryFormModal
      v-model="showFormModal"
      :category="editingCategory"
      :parent="parentCategory"
      @submit="handleFormSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useCategoriesStore } from '@/stores/categories'
import { useI18n } from '@/locales'
import type { Category } from '@/types'
import { CategoryChildRow, CategoryFormModal } from '@/widgets/categories'
import { BaseButton } from '@/components/ui'
import { ChevronLeft, Plus, Pencil, Trash2 } from '@lucide/vue'

const categoriesStore = useCategoriesStore()
const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const parentName = ref(decodeURIComponent(route.params.parentName as string))

const categories = computed(() => categoriesStore.categories)

const parentCategory = computed(() => {
  return categories.value.find(c => c.name === parentName.value && c.parent_id === null) || null
})

const subs = computed(() => {
  if (!parentCategory.value) return []
  return categories.value.filter(c => c.parent_id === parentCategory.value!.id)
})

const showFormModal = ref(false)
const editingCategory = ref<Category | null>(null)

watch(() => route.params.parentName, (newName) => {
  parentName.value = decodeURIComponent(newName as string)
})

onMounted(() => {
  parentName.value = decodeURIComponent(route.params.parentName as string)
})

const goBack = () => {
  router.push('/categories')
}

const openCreateModal = () => {
  editingCategory.value = null
  showFormModal.value = true
}

const handleEdit = (category: Category) => {
  editingCategory.value = category
  showFormModal.value = true
}

const handleDelete = async (id: string) => {
  if (!confirm(t('common.confirm') + '?')) return
  await categoriesStore.deleteCategory(id)
}

const handleFormSubmit = async (data: Partial<Category>) => {
  if (editingCategory.value) {
    await categoriesStore.updateCategory(editingCategory.value.id, data)
  } else if (parentCategory.value) {
    await categoriesStore.createCategory({ ...data, parent_id: parentCategory.value.id })
  }

  showFormModal.value = false
  editingCategory.value = null
}
</script>

<style scoped>
@reference "tailwindcss";

.sub-categories-page {
  @apply px-4 py-6 w-full;
}

.page-header {
  @apply flex items-center justify-between mb-6 pb-3 border-b;
  border-color: var(--border-primary);
}

.back-button {
  @apply w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0;
  color: var(--text-muted);
  background-color: transparent;
  border: none;
  cursor: pointer;
}

.back-button:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.back-button:focus-visible {
  @apply ring-2 ring-offset-2 rounded-lg;
  ring-offset-color: var(--bg-primary);
  ring-color: var(--accent);
}

.parent-action-bar {
  @apply mb-4;
  padding: 1rem;
  background-color: var(--bg-card);
  border-bottom: 1px solid var(--border-primary);
}

.parent-title {
  @apply text-base font-bold m-0 leading-snug text-left max-w-md;
  color: var(--text-primary);
  margin-bottom: 0.75rem;
}

.parent-actions {
  @apply flex items-center gap-2 flex-wrap justify-end;
}

.action-btn {
  @apply p-2 rounded-lg transition-colors flex items-center justify-center bg-transparent border;
  color: var(--text-muted);
  background-color: var(--bg-hover);
  border-color: var(--border-secondary);
  cursor: pointer;
  min-width: 44px;
  min-height: 44px;
}

.action-btn:hover {
  background-color: var(--bg-card);
  color: var(--text-primary);
}

.action-btn.danger:hover {
  background-color: var(--error-light);
  color: var(--error);
  border-color: var(--error-light);
}

.action-btn:focus-visible {
  @apply ring-2 ring-offset-2 rounded-lg;
  ring-offset-color: var(--bg-primary);
  ring-color: var(--accent);
}

.sub-categories-list {
  @apply space-y-2;
}

.empty-state {
  background-color: var(--bg-card);
  border: 1px dashed var(--border-primary);
  color: var(--text-muted);
  @apply text-center py-12 px-6 rounded-xl mt-6;
}

.empty-state p {
  @apply mb-4;
}

.empty-desc {
  @apply text-sm;
  color: var(--text-secondary);
}
</style>