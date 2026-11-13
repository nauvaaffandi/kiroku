export interface ReceiptJob {
  id: string
  user_id: string
  status: 'created' | 'processing' | 'draft_ready' | 'confirmed' | 'failed' | 'cancelled'
  stage: string
  original_size: number
  compressed_size: number
  ocr_text?: string
  ai_raw_output?: string
  transaction_draft?: string
  transaction_id?: string
  error_message?: string
  created_at: string
  updated_at: string
}

export interface ReceiptJobEvent {
  id: string
  job_id: string
  event_type: string
  stage: string
  attempt: number
  status: 'success' | 'warning' | 'error' | 'processing'
  message?: string
  metadata?: string
  created_at: string
}

export interface CreateReceiptJobReq {
  image_base64: string
  original_size: number
  compressed_size: number
}

export interface AINormalizationResult {
  type: string
  date_time: string
  amount: number
  currency: string
  from_account_id?: string | null
  to_account_id?: string | null
  category_id?: string | null
  note?: string
  payee_payer?: string
  payment_type?: string
  warranty_until?: string | null
  status?: string
  place_name?: string
  place_lat?: number | null
  place_lng?: number | null
  labels?: string
}