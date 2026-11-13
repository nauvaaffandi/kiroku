// Template Types
export interface Template {
  id: string
  name: string
  type: string
  amount: number | null
  currency: string | null
  from_account_id: string | null
  to_account_id: string | null
  category_id: string | null
  labels: string
  note: string
  payee_payer: string
  payment_type: string
  warranty_until: string | null
  status: string
  place_name: string
  place_lat: number | null
  place_lng: number | null
  attachments: string
  created_at: string
  updated_at: string
}

export interface CreateTemplateRequest {
  name: string
  type: string
  amount?: number | null
  currency?: string | null
  from_account_id?: string | null
  to_account_id?: string | null
  category_id?: string | null
  labels?: string
  note?: string
  payee_payer?: string
  payment_type?: string
  warranty_until?: string | null
  status?: string
  place_name?: string
  place_lat?: number | null
  place_lng?: number | null
  attachments?: string
}

export interface UpdateTemplateRequest {
  name?: string
  type?: string
  amount?: number | null
  currency?: string | null
  from_account_id?: string | null
  to_account_id?: string | null
  category_id?: string | null
  labels?: string
  note?: string
  payee_payer?: string
  payment_type?: string
  warranty_until?: string | null
  status?: string
  place_name?: string
  place_lat?: number | null
  place_lng?: number | null
  attachments?: string
}

export interface ApplyTemplateResponse {
  type: string
  date_time: string
  amount: number
  currency: string | null
  from_account_id: string | null
  to_account_id: string | null
  category_id: string | null
  labels: string
  note: string
  payee_payer: string
  payment_type: string
  warranty_until: string | null
  status: string
  place_name: string
  place_lat: number | null
  place_lng: number | null
  attachments: string
  template_id: string
}