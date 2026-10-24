import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/client'
import { useUIStore } from './ui'
import type { Category, CategoryFilters } from '@/types'

export const useCategoriesStore = defineStore('categories', () => {
  const categories = ref<Category[]>([])
  const tree = ref<Category[]>([])
  const loading = ref(false)
  const uiStore = useUIStore()
  const pagination = ref({ page: 1, pageSize: 50, total: 0 })
  const filters = ref<CategoryFilters>({
    q: '',
    type: '',
    archived: '',
    sort: 'sort_order',
    order: 'ASC'
  })

  const fetchCategories = async (page = 1) => {
    loading.value = true
    try {
      const params = new URLSearchParams({
        page: String(page),
        pageSize: String(pagination.value.pageSize),
        sort: filters.value.sort,
        order: filters.value.order,
        ...Object.fromEntries(
          Object.entries(filters.value).filter(([, v]) => v !== '')
        )
      })
      const [listRes, treeRes] = await Promise.all([
        api.get(`/api/v1/categories?${params}`),
        api.get('/api/v1/categories/tree')
      ])
      categories.value = Array.isArray(listRes.data) ? listRes.data : []
      tree.value = Array.isArray(treeRes.data) ? treeRes.data : []
      pagination.value = {
        ...pagination.value,
        page,
        total: categories.value.length
      }
    } catch (e) {
      uiStore.error('Failed to fetch categories')
    } finally {
      loading.value = false
    }
  }

  const createCategory = async (data: Partial<Category>) => {
    try {
      const response = await api.post('/api/v1/categories', data)
      await fetchCategories(pagination.value.page)
      return response.data.id
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to create category')
      throw e
    }
  }

  const updateCategory = async (id: string, data: Partial<Category>) => {
    try {
      await api.put(`/api/v1/categories/${id}`, data)
      await fetchCategories(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to update category')
      throw e
    }
  }

  const deleteCategory = async (id: string) => {
    try {
      await api.delete(`/api/v1/categories/${id}`)
      await fetchCategories(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to delete category')
      throw e
    }
  }

  const reorderCategories = async (items: { id: string; sort_order: number }[]) => {
    try {
      await api.put('/api/v1/categories/reorder', { items })
      await fetchCategories(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || 'Failed to reorder categories')
      throw e
    }
  }

  return {
    categories,
    tree,
    loading,
    pagination,
    filters,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    reorderCategories
  }
})