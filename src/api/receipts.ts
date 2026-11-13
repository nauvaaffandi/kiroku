import { api } from '@/api/client'
import type { ReceiptJob, ReceiptJobEvent, CreateReceiptJobReq, AINormalizationResult } from '@/types'

export const receiptsApi = {
  async scanReceipt(data: CreateReceiptJobReq) {
    const response = await api.post('/api/v1/receipts/scan', data)
    return response.data as { job_id: string }
  },

  async getJob(jobId: string) {
    const response = await api.get(`/api/v1/receipts/jobs/${jobId}`)
    return response.data as { job: ReceiptJob; events: ReceiptJobEvent[] }
  },

  async confirmTransaction(jobId: string, data: AINormalizationResult) {
    const response = await api.post(`/api/v1/receipts/jobs/${jobId}/confirm`, data)
    return response.data as { id: string }
  }
}