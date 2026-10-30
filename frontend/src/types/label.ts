export interface Label {
  id: string
  name: string
  color: string
  icon?: string
  description?: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface CreateLabelRequest {
  name: string
  color?: string
  icon?: string
  description?: string
  sort_order?: number
}

export interface UpdateLabelRequest {
  name?: string
  color?: string
  icon?: string
  description?: string
  sort_order?: number
}
