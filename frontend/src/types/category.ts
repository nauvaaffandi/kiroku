// Category Types
export interface CategoryFilters {
  q: string
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
