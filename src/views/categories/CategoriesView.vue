<template>
  <div class="categories-page">
    <header class="page-header">
      <div>
        <h1>{{ t('categories.title') }}</h1>
        <p class="page-subtitle">{{ t('categories.subtitle') }}</p>
      </div>
      <BaseButton @click="openCreateModal">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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

    <section class="parent-section">
      <header class="section-header">
        <h2>{{ t('categories.parentCategory') }}</h2>
      </header>

      <CategoryParentRow
        v-for="parent in parents"
        :key="parent.id"
        :category="parent"
        @select="() => goToSub(parent.name)"
      />

      <div v-if="parents.length === 0" class="empty-state">
        <p>{{ t('categories.noCategories') }}</p>
        <p class="empty-desc">{{ t('categories.noCategoriesDescription') }}</p>
      </div>
    </section>

    <CategoryFormModal
      v-model="showFormModal"
      :category="editingCategory"
      :parent="creatingChildParent"
      @submit="handleFormSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useCategoriesStore } from '@/stores/categories'
import { useI18n } from '@/locales'
import type { Category, CategoryFilters as CategoryFiltersType } from '@/types'
import { CategoryParentRow, CategoryFormModal, CategoryFilters } from '@/widgets/categories'
import { BaseButton } from '@/components/ui'

const categoriesStore = useCategoriesStore()
const { t } = useI18n()
const router = useRouter()

const categories = computed(() => categoriesStore.categories)

const parents = computed(() => categories.value.filter(c => c.parent_id === null))

const showFormModal = ref(false)
const editingCategory = ref<Category | null>(null)
const creatingChildParent = ref<Category | null>(null)

const localFilters = ref<CategoryFiltersType>({
  q: '',
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

function goToSub(parentName: string) {
  router.push(`/categories/sub/${encodeURIComponent(parentName)}`)
}

const openCreateModal = () => {
  editingCategory.value = null
  creatingChildParent.value = null
  showFormModal.value = true
}

const handleFormSubmit = async (data: Partial<Category>) => {
  await categoriesStore.createCategory(data)

  showFormModal.value = false
  editingCategory.value = null
  creatingChildParent.value = null
  await fetchCategories(1)
}
</script>

<style scoped>
@reference "tailwindcss";

.categories-page {
  @apply px-4 py-6;
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

.parent-section {
  @apply mb-6;
}

.section-header {
  @apply mb-4 pb-3 border-b;
  border-color: var(--border-primary);
}

.section-header h2 {
  @apply text-lg font-semibold m-0;
  color: var(--text-primary);
}

.empty-state {
  background-color: var(--bg-card);
  border: 1px dashed var(--border-primary);
  color: var(--text-muted);
  @apply text-center py-8 px-6 rounded-xl;
}

.empty-state p {
  @apply mb-2;
}

.empty-desc {
  @apply text-sm;
  color: var(--text-secondary);
}
</style>