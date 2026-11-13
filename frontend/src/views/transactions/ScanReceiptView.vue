<template>
  <div class="scan-receipt-page">
    <header class="page-header">
      <div>
        <h1>{{ t('nav.scanReceipt') }}</h1>
        <p class="page-subtitle">{{ t('receipt.scanDesc') }}</p>
      </div>
      <BaseButton variant="secondary" @click="router.push('/transactions')">
        {{ t('receipt.backToTransactions') }}
      </BaseButton>
    </header>

    <div class="scan-container">
      <ReceiptUploadCard v-if="step === 'upload'" @file-selected="processFile" />

      <div v-if="step === 'processing'" class="card processing-card">
        <h3>{{ t('receipt.processingTitle') }}</h3>
        <p class="text-muted">{{ t('receipt.processingDesc') }}</p>
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
import { useAccountsStore } from '@/stores/accounts'
import { useCategoriesStore } from '@/stores/categories'
import { BaseButton } from '@/components/ui'
import type { ReceiptJobEvent, AINormalizationResult, Category, Account } from '@/types'
import { useI18n } from '@/locales'
import ReceiptUploadCard from './components/ReceiptUploadCard.vue'
import ReceiptProcessingCard from './components/ReceiptProcessingCard.vue'
import ReceiptReviewCard from './components/ReceiptReviewCard.vue'

const router = useRouter()
const uiStore = useUIStore()
const accountsStore = useAccountsStore()
const categoriesStore = useCategoriesStore()
const { t } = useI18n()

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
    uiStore.info(t('receipt.compressing'))
    const { originalSize: orig, compressedSize: comp } = await compressImageIfNeeded(file)
    originalSize.value = orig
    compressedSize.value = comp

    step.value = 'processing'
    uiStore.success(t('receipt.prepared'))

    const res: any = { job_id: 'offline', transaction_draft: null }
    jobId.value = res.job_id
    startPolling(res.job_id)
  } catch (e: any) {
    uiStore.error(e.message || t('receipt.processFailed'))
    step.value = 'upload'
  }
}

const startPolling = (_id: string) => {
  pollTimer = setInterval(async () => {
    try {
      if (pollTimer) clearInterval(pollTimer)
      step.value = 'review'
      uiStore.info('OCR receipt processing not available offline')
    } catch (e) {
      // ignore
    }
  }, 2000)
}

const confirmTransaction = async (_draftData: AINormalizationResult) => {
  if (!jobId.value) return
  submitting.value = true
  try {
    uiStore.success('Transaction created (offline draft)')
    router.push('/transactions')
  } catch (e: any) {
    uiStore.error(e.response?.data?.error || t('receipt.saveFailed'))
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
    await categoriesStore.fetchCategories()
    await accountsStore.fetchAccounts()
    categories.value = categoriesStore.categories
    accounts.value = accountsStore.accounts
  } catch (e) {
    // ignore
  }
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<style scoped>
@reference "tailwindcss";
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