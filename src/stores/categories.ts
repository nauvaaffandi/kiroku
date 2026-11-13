import { defineStore } from 'pinia'
import { ref } from 'vue'
import { localDb } from '@/storage/repository'
import { useUIStore } from './ui'
import { useI18n } from '@/locales'
import type { Category, CategoryFilters } from '@/types'

export const useCategoriesStore = defineStore('categories', () => {
  const categories = ref<Category[]>([])
  const tree = ref<Category[]>([])
  const loading = ref(false)
  const uiStore = useUIStore()
  const { t } = useI18n()
  const pagination = ref({ page: 1, pageSize: 50, total: 0 })
  const filters = ref<CategoryFilters>({
    q: '',
    archived: '',
    sort: 'sort_order',
    order: 'ASC'
  })

  const buildTree = (list: Category[]) => {
    const byId = new Map<string, Category>()
    list.forEach((c) => byId.set(c.id, { ...c, children: [] }))
    const roots: Category[] = []
    byId.forEach((node) => {
      const parent = node.parent_id ? byId.get(node.parent_id) : null
      if (parent) parent.children.push(node)
      else roots.push(node)
    })
    return roots
  }

  const fetchCategories = async (page = 1) => {
    loading.value = true
    try {
      const rows = await localDb.list<Category>('category')
      let list = rows.filter((c) => !c.archived)
      if (filters.value.q) {
        const q = filters.value.q.toLowerCase()
        list = list.filter((c) => c.name.toLowerCase().includes(q))
      }
      if (filters.value.archived === 'true') list = rows
      list.sort((a, b) => {
        const d = filters.value.order === 'DESC' ? -1 : 1
        return a.sort_order - b.sort_order === 0 ? 0 : (a.sort_order < b.sort_order ? -1 : 1) * d
      })
      categories.value = list
      tree.value = buildTree(rows)
      pagination.value = { ...pagination.value, page, total: list.length }
    } catch (e) {
      uiStore.error(t('categories.fetchError'))
    } finally {
      loading.value = false
    }
  }

  const createCategory = async (data: Partial<Category>) => {
    try {
      const id = crypto.randomUUID()
      const node: Category = {
        id,
        parent_id: data.parent_id ?? null,
        parent_name: null,
        name: data.name ?? '',
        icon: data.icon ?? '',
        color: data.color ?? '',
        sort_order: data.sort_order ?? 0,
        depth: 0,
        path: '',
        children: [],
        archived: false,
        ...data
      } as Category
      await localDb.put('category', id, node)
      await fetchCategories(pagination.value.page)
      return id
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to create category')
      throw e
    }
  }

  const updateCategory = async (id: string, data: Partial<Category>) => {
    try {
      const current = await localDb.get<Category>('category', id)
      if (!current) throw new Error('Category not found')
      await localDb.put('category', id, { ...current, ...data } as Category)
      await fetchCategories(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to update category')
      throw e
    }
  }

  const deleteCategory = async (id: string) => {
    try {
      await localDb.remove('category', id)
      await fetchCategories(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to delete category')
      throw e
    }
  }

  const reorderCategories = async (items: { id: string; sort_order: number }[]) => {
    try {
      for (const item of items) {
        const current = await localDb.get<Category>('category', item.id)
        if (current) await localDb.put('category', item.id, { ...current, sort_order: item.sort_order } as Category)
      }
      await fetchCategories(pagination.value.page)
    } catch (e: any) {
      uiStore.error(e?.message || 'Failed to reorder categories')
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
