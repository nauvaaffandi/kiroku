// Category Types
export type CategoryType = 'income' | 'expense' | 'transfer'

export interface CategoryFilters {
  q: string
  type: string
  archived: string
  sort: string
  order: string
}

export interface Category {
  id: string
  parent_id: string | null
  name: string
  icon: string
  color: string
  type: CategoryType
  sort_order: number
  depth: number
  path: string
  children: Category[]
  archived: boolean
}

export interface CategoryTree extends Category {}

export interface CreateCategoryRequest {
  parent_id?: string | null
  name: string
  icon?: string
  color?: string
  kind?: string | null
  sort_order?: number
}

export interface UpdateCategoryRequest {
  parent_id?: string | null
  name?: string
  icon?: string
  color?: string
  kind?: string | null
  sort_order?: number
}