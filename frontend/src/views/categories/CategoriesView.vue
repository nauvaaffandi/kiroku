<template>
  <div class="categories-page">
    <header class="page-header">
      <div>
        <h1>{{ t('categories.title') }}</h1>
        <p class="page-subtitle">{{ t('categories.subtitle') }}</p>
      </div>
      <BaseButton @click="openCreateModal">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        {{ t('categories.addCategory') }}
      </BaseButton>
    </header>

    <CategoryFilters
      v-model="localFilters"
      :searchQuery="searchQuery"
      @change="handleFilterChange"
      @search="handleSearch"
    />

    <CategoryTree
      :categories="categories"
      @edit="handleEdit"
      @create-child="handleCreateChild"
      @delete="handleDelete"
      @create="openCreateModal"
    />

    <CategoryFormModal
      v-model="showFormModal"
      :category="editingCategory"
      :parent="creatingChildParent"
      @submit="handleFormSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useCategoriesStore } from '@/stores/categories'
import { useI18n } from '@/locales'
import type { Category, CategoryFilters } from '@/types'
import { CategoryTree, CategoryFormModal } from '@/widgets/categories'
import { BaseButton } from '@/components/ui'

const categoriesStore = useCategoriesStore()
const { t } = useI18n()

const categories = computed(() => categoriesStore.categories)

const showFormModal = ref(false)
const editingCategory = ref<Category | null>(null)
const creatingChildParent = ref<Category | null>(null)

const localFilters = ref<CategoryFilters>({
  q: '',
  type: '',
  archived: '',
  sort: 'sort_order',
  order: 'ASC'
})

const searchQuery = ref('')

let searchTimeout: ReturnType<typeof setTimeout>
watch(searchQuery, (value) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    categoriesStore.filters.q = value
    fetchCategories(1)
  }, 300)
})

const fetchCategories = async (page = 1) => {
  await categoriesStore.fetchCategories(page)
}

onMounted(() => {
  fetchCategories(1)
})

const handleFilterChange = () => {
  categoriesStore.filters = { ...localFilters.value }
  fetchCategories(1)
}

const handleSearch = (value: string) => {
  searchQuery.value = value
}

const handleEdit = (category: Category) => {
  editingCategory.value = category
  creatingChildParent.value = null
  showFormModal.value = true
}

const handleCreateChild = (parent: Category) => {
  editingCategory.value = null
  creatingChildParent.value = parent
  showFormModal.value = true
}

const handleDelete = (id: string) => {
  if (confirm('Are you sure you want to delete this category?')) {
    categoriesStore.deleteCategory(id)
  }
}

const openCreateModal = () => {
  editingCategory.value = null
  creatingChildParent.value = null
  showFormModal.value = true
}

const handleFormSubmit = async (data: Partial<Category>) => {
  if (editingCategory.value) {
    await categoriesStore.updateCategory(editingCategory.value.id, data)
  } else if (creatingChildParent.value) {
    await categoriesStore.createCategory({ ...data, parent_id: creatingChildParent.value.id })
  } else {
    await categoriesStore.createCategory(data)
  }
  showFormModal.value = false
  editingCategory.value = null
  creatingChildParent.value = null
  await fetchCategories(1)
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