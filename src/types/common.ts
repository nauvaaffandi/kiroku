// API Response Wrappers
export interface ApiResponse<T> {
  data: T
  error?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}

// Common Types
export type SortOrder = 'ASC' | 'DESC'

export interface PaginationParams {
  page: number
  limit: number
  sort?: string
  order?: SortOrder
}

export interface DateRange {
  from: string
  to: string
  preset?: string
}

export interface SelectOption {
  value: string
  label: string
  disabled?: boolean
}

// Re-export for convenience
export type { SortOrder as SortOrderType }
export type { PaginationParams as PaginationParamsType }