<template>
  <div class="scan-receipt-page">
    <header class="page-header">
      <div>
        <h1>Scan Receipt</h1>
        <p class="page-subtitle">Convert receipt images into transaction drafts automatically</p>
      </div>
      <BaseButton variant="secondary" @click="router.push('/transactions')">
        Back to Transactions
      </BaseButton>
    </header>

    <div class="scan-container">
      <ReceiptUploadCard v-if="step === 'upload'" @file-selected="processFile" />

      <div v-if="step === 'processing'" class="card processing-card">
        <h3>Processing Receipt...</h3>
        <p class="text-muted">Please wait while we extract text with OCR and normalize with AI.</p>
        <ReceiptProcessingCard :originalSize="originalSize" :compressedSize="compressedSize" :events="events" />
      </div>

      <ReceiptReviewCard 
        v-if="step === 'review' && draft" 
        :draft="draft" 
        :accounts="accounts" 
        :categories="categories" 
        :submitting="submitting"
        @confirm="confirmTransaction" 
        @reset="resetScan" 
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUIStore } from '@/stores/ui'
import { receiptsApi } from '@/api/receipts'
import { api } from '@/api/client'
import { BaseButton } from '@/components/ui'
import type { ReceiptJobEvent, AINormalizationResult, Category, Account } from '@/types'
import ReceiptUploadCard from './components/ReceiptUploadCard.vue'
import ReceiptProcessingCard from './components/ReceiptProcessingCard.vue'
import ReceiptReviewCard from './components/ReceiptReviewCard.vue'

const router = useRouter()
const uiStore = useUIStore()

const step = ref<'upload' | 'processing' | 'review'>('upload')
const originalSize = ref(0)
const compressedSize = ref(0)
const jobId = ref<string | null>(null)
const events = ref<ReceiptJobEvent[]>([])
const draft = ref<AINormalizationResult | null>(null)
const submitting = ref(false)

const categories = ref<Category[]>([])
const accounts = ref<Account[]>([])

let pollTimer: ReturnType<typeof setInterval> | null = null

const compressImageIfNeeded = async (file: File): Promise<{ base64: string; originalSize: number; compressedSize: number }> => {
  originalSize.value = file.size
  const threshold = 2 * 1024 * 1024 // 2MB

  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (event) => {
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement('canvas')
        canvas.width = img.width
        canvas.height = img.height
        const ctx = canvas.getContext('2d')
        if (!ctx) {
          reject(new Error('Canvas context failed'))
          return
        }
        ctx.drawImage(img, 0, 0)

        let quality = 0.9
        let base64 = canvas.toDataURL('image/jpeg', quality)

        if (file.size > threshold) {
          while (base64.length * 0.75 > threshold && quality > 0.3) {
            quality -= 0.1
            base64 = canvas.toDataURL('image/jpeg', quality)
          }
        }

        const base64Data = base64.split(',')[1]
        const compressedBytes = Math.round(base64Data.length * 0.75)
        compressedSize.value = compressedBytes

        resolve({
          base64: base64Data,
          originalSize: file.size,
          compressedSize: compressedBytes
        })
      }
      img.onerror = reject
      img.src = event.target?.result as string
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

const processFile = async (file: File) => {
  try {
    uiStore.info('Compressing image...')
    const { base64, originalSize: orig, compressedSize: comp } = await compressImageIfNeeded(file)
    originalSize.value = orig
    compressedSize.value = comp

    step.value = 'processing'
    uiStore.success('Image prepared successfully')

    const res = await receiptsApi.scanReceipt({
      image_base64: base64,
      original_size: orig,
      compressed_size: comp
    })

    jobId.value = res.job_id
    startPolling(res.job_id)
  } catch (e: any) {
    uiStore.error(e.message || 'Failed to process image')
    step.value = 'upload'
  }
}

const startPolling = (id: string) => {
  pollTimer = setInterval(async () => {
    try {
      const res = await receiptsApi.getJob(id)
      events.value = res.events || []

      if (res.job.status === 'draft_ready') {
        if (pollTimer) clearInterval(pollTimer)
        if (res.job.transaction_draft) {
          draft.value = JSON.parse(res.job.transaction_draft)
          if (draft.value && draft.value.date_time) {
            draft.value.date_time = new Date(draft.value.date_time).toISOString().slice(0, 16)
          }
        }
        step.value = 'review'
        uiStore.success('Transaction draft ready for review!')
      } else if (res.job.status === 'failed') {
        if (pollTimer) clearInterval(pollTimer)
        uiStore.error(res.job.error_message || 'Receipt processing failed')
        step.value = 'upload'
      }
    } catch (e) {
      // ignore
    }
  }, 2000)
}

const confirmTransaction = async (draftData: AINormalizationResult) => {
  if (!jobId.value) return
  submitting.value = true
  try {
    await receiptsApi.confirmTransaction(jobId.value, draftData)
    uiStore.success('Transaction created successfully from receipt!')
    router.push('/transactions')
  } catch (e: any) {
    uiStore.error(e.response?.data?.error || 'Failed to save transaction')
  } finally {
    submitting.value = false
  }
}

const resetScan = () => {
  if (pollTimer) clearInterval(pollTimer)
  jobId.value = null
  events.value = []
  draft.value = null
  step.value = 'upload'
}

onMounted(async () => {
  try {
    const [catRes, accRes] = await Promise.all([
      api.get('/api/v1/categories'),
      api.get('/api/v1/accounts')
    ])
    categories.value = catRes.data || []
    accounts.value = accRes.data || []
  } catch (e) {
    // ignore
  }
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<style scoped>
.scan-receipt-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1.5rem;
  background-color: var(--color-neutral-100);
  border-bottom: 1px solid var(--color-neutral-200);
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-neutral-900);
}

.page-subtitle {
  color: var(--color-neutral-500);
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.scan-container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 0 1rem;
  width: 100%;
}

.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  padding: 2rem;
  box-shadow: var(--shadow);
}
</style>