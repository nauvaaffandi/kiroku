<template>
  <div class="shared-access-page">
    <header class="page-header">
      <div class="header-content">
        <h1>Shared Access</h1>
        <button @click="showCreateModal = true" class="btn-primary">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          Create Share
        </button>
      </div>
    </header>

    <div class="shared-access-content">
      <div class="shares-list" v-if="shares.length > 0">
        <div class="share-card" v-for="share in shares" :key="share.id">
          <div class="share-info">
            <div class="share-account">
              <span class="account-name">{{ share.account_name }}</span>
              <span class="account-badge">{{ share.account_id }}</span>
            </div>
            <div class="share-details">
              <div class="detail-item">
                <span class="detail-label">Link</span>
                <div class="token-display">
                  <span class="token">{{ share.token }}</span>
                  <button @click="copyShareUrl(share.token)" class="btn-icon" title="Copy full link">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                    </svg>
                  </button>
                </div>
              </div>
              <div class="detail-item">
                <span class="detail-label">Created</span>
                <span class="detail-value">{{ formatDate(share.created_at) }}</span>
              </div>
              <div class="detail-item" v-if="share.expires_at">
                <span class="detail-label">Expires</span>
                <span class="detail-value">{{ formatDate(share.expires_at) }}</span>
              </div>
              <div class="detail-item" v-if="share.last_accessed_at">
                <span class="detail-label">Last Accessed</span>
                <span class="detail-value">{{ formatDate(share.last_accessed_at) }}</span>
              </div>
              <div class="detail-item" v-else>
                <span class="detail-label">Last Accessed</span>
                <span class="detail-value text-muted">Never</span>
              </div>
            </div>
          </div>
          <div class="share-actions">
            <button @click="copyShareUrl(share.token)" class="btn-secondary btn-sm">Copy Link</button>
            <button @click="revokeShare(share.id)" class="btn-danger btn-sm">Revoke</button>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
          <polyline points="22 6 12 16 10 14"></polyline>
        </svg>
        <h3>No shared links yet</h3>
        <p>Create a share link to give read-only access to an account</p>
        <button @click="showCreateModal = true" class="btn-primary">Create Share</button>
      </div>
    </div>

    <Transition name="modal">
      <div v-if="showCreateModal" class="modal-overlay" @click.self="closeCreateModal">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h2>Create Share Link</h2>
            <button @click="closeCreateModal" class="modal-close" aria-label="Close">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
          <form @submit.prevent="createShare" class="modal-form">
            <div class="form-group">
              <label for="account_id">Account</label>
              <select id="account_id" v-model="createForm.account_id" required class="form-input">
                <option value="">Select an account</option>
                <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                  {{ acc.name }} ({{ acc.currency }})
                </option>
              </select>
            </div>
            <div class="form-group">
              <label for="expires_at">Expires At (optional)</label>
              <input
                id="expires_at"
                v-model="createForm.expires_at"
                type="datetime-local"
                class="form-input"
              />
              <p class="form-hint">Leave empty for no expiration</p>
            </div>
            <div class="modal-actions">
              <button type="button" @click="closeCreateModal" class="btn-secondary">Cancel</button>
              <button type="submit" :disabled="loading" class="btn-primary">
                <span v-if="loading">Creating...</span>
                <span v-else>Create Share</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <Transition name="modal">
      <div v-if="createdShare" class="modal-overlay" @click.self="closeCreatedModal">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h2>Share Created</h2>
          </div>
          <div class="modal-form">
            <div class="share-result">
              <p>Share link created successfully. Copy and share this link:</p>
              <div class="share-url-display">
                <input :value="fullShareUrl" readonly class="form-input" />
                <button @click="copyFullShareUrl" class="btn-secondary">Copy</button>
              </div>
              <p class="share-warning">Anyone with this link can view the account (read-only).</p>
            </div>
            <div class="modal-actions">
              <button @click="closeCreatedModal" class="btn-primary">Done</button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useUIStore } from '@/stores/ui'
import { api } from '@/api/client'

const uiStore = useUIStore()

const shares = ref<any[]>([])
const accounts = ref<any[]>([])
const loading = ref(false)
const showCreateModal = ref(false)
const createdShare = ref<any>(null)
const createForm = ref({ account_id: '', expires_at: '' })

const loadData = async () => {
  try {
    const [sharesRes, accountsRes] = await Promise.all([
      api.get('/api/v1/shares'),
      api.get('/api/v1/accounts')
    ])
    shares.value = sharesRes.data || []
    accounts.value = accountsRes.data.filter((a: any) => !a.archived) || []
  } catch (e) {
    uiStore.error('Failed to load data')
  }
}

const createShare = async () => {
  loading.value = true
  try {
    const payload: any = { account_id: createForm.value.account_id }
    if (createForm.value.expires_at) {
      payload.expires_at = new Date(createForm.value.expires_at).toISOString()
    }
    const res = await api.post('/api/v1/shares', payload)
    createdShare.value = res.data
    showCreateModal.value = false
    await loadData()
  } catch (e: any) {
    uiStore.error(e.response?.data?.error || 'Failed to create share')
  } finally {
    loading.value = false
  }
}

const revokeShare = async (id: string) => {
  if (!confirm('Are you sure you want to revoke this share link?')) return
  try {
    await api.delete(`/api/v1/shares/${id}`)
    uiStore.success('Share revoked')
    await loadData()
  } catch (e) {
    uiStore.error('Failed to revoke share')
  }
}

const copyShareUrl = async (token: string) => {
  const url = `${window.location.origin}/shared/${token}`
  await navigator.clipboard.writeText(url)
  uiStore.success('Link copied to clipboard')
}

const copyFullShareUrl = async () => {
  await navigator.clipboard.writeText(fullShareUrl.value)
  uiStore.success('Link copied to clipboard')
}

const closeCreateModal = () => {
  showCreateModal.value = false
  createForm.value = { account_id: '', expires_at: '' }
}

const closeCreatedModal = () => {
  createdShare.value = null
}

const fullShareUrl = computed(() => {
  if (!createdShare.value) return ''
  return `${window.location.origin}/shared/${createdShare.value.token}`
})

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('ja-JP', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(loadData)
</script>

<style scoped>
@reference "tailwindcss";

.shared-access-page {
  @apply min-h-screen flex flex-col;
}

.page-header {
  @apply p-6;
  background-color: var(--bg-card);
  border-bottom: 1px solid var(--border-primary);
}

.header-content {
  @apply flex items-center justify-between flex-wrap gap-4;
}

.header-content h1 { 
  @apply text-2xl font-semibold; 
  color: var(--text-primary); 
}

.shared-access-content {
  @apply flex-1 p-6 max-w-[900px] mx-auto w-full;
}

.shares-list {
  @apply flex flex-col gap-4;
}

.share-card {
  @apply rounded-[var(--radius)] p-6 flex items-start justify-between gap-6 flex-wrap shadow-[var(--shadow)];
  background-color: var(--bg-card);
}

.share-info { 
  @apply flex-1 min-w-[280px]; 
}

.share-account { 
  @apply mb-4; 
}

.account-name { 
  @apply font-semibold text-lg; 
  color: var(--text-primary); 
}

.account-badge { 
  @apply text-xs ml-2 font-mono; 
  color: var(--text-muted); 
}

.share-details { 
  @apply grid gap-3;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.detail-item { 
  @apply flex flex-col gap-0.5; 
}

.detail-label { 
  @apply text-xs font-semibold uppercase tracking-wider; 
  color: var(--text-muted); 
}

.detail-value { 
  @apply text-sm; 
  color: var(--text-primary); 
}

.detail-value.text-muted { 
  color: var(--text-muted); 
}

.token-display { 
  @apply flex items-center gap-2; 
}

.token { 
  @apply font-mono text-[0.8125rem] px-2 py-1 rounded-[var(--radius-sm)] flex-1 break-all; 
  color: var(--text-secondary);
  background-color: var(--bg-tertiary);
}

.share-actions { 
  @apply flex gap-2 shrink-0; 
}

.btn-sm { 
  @apply py-1.5 px-3 text-[0.8125rem]; 
}

.empty-state {
  @apply text-center py-16 px-8 rounded-[var(--radius)] shadow-[var(--shadow)];
  background-color: var(--bg-card);
  color: var(--text-muted);
}

.empty-state svg { 
  @apply mb-4 opacity-50; 
}

.empty-state h3 { 
  @apply text-lg font-semibold mb-2; 
  color: var(--text-primary); 
}

.empty-state p { 
  @apply text-sm mb-6; 
}

/* Modal */
.modal-overlay { 
  @apply fixed inset-0 flex items-center justify-center p-4 z-[1000]; 
  background-color: rgba(0, 0, 0, 0.5); 
  animation: fadeIn 0.2s ease; 
}

.modal { 
  @apply rounded-[var(--radius)] w-full max-w-[480px] max-h-[90vh] overflow-y-auto shadow-[var(--shadow-md)]; 
  background-color: var(--bg-card);
  animation: slideUp 0.2s ease; 
}

.modal-header { 
  @apply flex items-center justify-between py-4 px-6; 
  border-bottom: 1px solid var(--border-primary); 
}

.modal-header h2 { 
  @apply text-lg font-semibold; 
  color: var(--text-primary);
}

.modal-close { 
  @apply p-1 rounded-[var(--radius-sm)] bg-transparent border-none cursor-pointer transition-colors duration-200; 
  color: var(--text-muted); 
}

.modal-close:hover { 
  @apply text-[var(--text-primary)]; 
  background-color: var(--bg-hover); 
}

.modal-form { 
  @apply p-6; 
}

.form-group { 
  @apply mb-4; 
}

.form-group label { 
  @apply block text-sm font-medium mb-1.5; 
  color: var(--text-secondary); 
}

.form-input { 
  @apply w-full py-2 px-3 text-sm rounded-[var(--radius-sm)] transition-all duration-200; 
  border: 1px solid var(--border-primary); 
  background-color: var(--bg-primary); 
  color: var(--text-primary); 
}

.form-input:focus { 
  @apply outline-none border-transparent ring-2 ring-[var(--accent)]; 
}

.form-hint { 
  @apply text-xs mt-1; 
  color: var(--text-muted); 
}

.modal-actions { 
  @apply flex justify-end gap-3 mt-6; 
}

.btn-primary, .btn-secondary, .btn-danger { 
  @apply py-2 px-4 rounded-[var(--radius-sm)] font-medium text-sm cursor-pointer transition-opacity duration-200 inline-flex items-center justify-center gap-2 border-none; 
}

.btn-primary { 
  @apply text-white; 
  background-color: var(--accent); 
}

.btn-primary:hover:not(:disabled) { 
  background-color: var(--accent-hover); 
}

.btn-secondary { 
  @apply text-[var(--text-secondary)]; 
  background-color: var(--bg-tertiary); 
}

.btn-secondary:hover:not(:disabled) { 
  background-color: var(--bg-hover); 
}

.btn-danger { 
  @apply text-white; 
  background-color: var(--error); 
}

.btn-danger:hover:not(:disabled) { 
  background-color: var(--error-light); 
}

.btn-primary:disabled, .btn-secondary:disabled, .btn-danger:disabled { 
  @apply opacity-50 cursor-not-allowed; 
}

.btn-icon { 
  @apply p-1.5 rounded-[var(--radius-sm)] bg-transparent border-none cursor-pointer transition-colors duration-200; 
  color: var(--text-muted); 
}

.btn-icon:hover { 
  @apply text-[var(--text-primary)]; 
  background-color: var(--bg-hover); 
}

.share-result p { 
  @apply mb-3; 
  color: var(--text-secondary); 
}

.share-url-display { 
  @apply flex gap-2 mb-3; 
}

.share-url-display .form-input { 
  @apply flex-1; 
  background-color: var(--bg-tertiary); 
}

.share-warning { 
  @apply text-[0.8125rem] m-0; 
  color: var(--warning); 
}

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
</style>