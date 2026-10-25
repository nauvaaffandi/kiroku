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
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-header {
  padding: 1.5rem;
  background-color: var(--color-neutral-100);
  border-bottom: 1px solid var(--color-neutral-200);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1rem;
}

.header-content h1 { font-size: 1.5rem; font-weight: 600; color: var(--color-neutral-900); }

.shared-access-content {
  flex: 1;
  padding: 1.5rem;
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
}

.shares-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.share-card {
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 1.5rem;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.share-info { flex: 1; min-width: 280px; }

.share-account { margin-bottom: 1rem; }

.account-name { font-weight: 600; color: var(--color-neutral-900); font-size: 1.125rem; }

.account-badge { font-size: 0.75rem; color: var(--color-neutral-500); margin-left: 0.5rem; font-family: monospace; }

.share-details { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 0.75rem; }

.detail-item { display: flex; flex-direction: column; gap: 0.125rem; }

.detail-label { font-size: 0.75rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; color: var(--color-neutral-500); }

.detail-value { font-size: 0.875rem; color: var(--color-neutral-900); }

.detail-value.text-muted { color: var(--color-neutral-400); }

.token-display { display: flex; align-items: center; gap: 0.5rem; }

.token { font-family: monospace; font-size: 0.8125rem; color: var(--color-neutral-600); background-color: var(--color-neutral-50); padding: 0.25rem 0.5rem; border-radius: var(--radius-sm); flex: 1; word-break: break-all; }

.share-actions { display: flex; gap: 0.5rem; flex-shrink: 0; }

.btn-sm { padding: 0.375rem 0.75rem; font-size: 0.8125rem; }

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  color: var(--color-neutral-500);
}

.empty-state svg { margin-bottom: 1rem; opacity: 0.5; }

.empty-state h3 { font-size: 1.125rem; font-weight: 600; color: var(--color-neutral-700); margin-bottom: 0.5rem; }

.empty-state p { font-size: 0.875rem; margin-bottom: 1.5rem; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background-color: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; padding: 1rem; z-index: 1000; animation: fadeIn 0.2s ease; }
.modal { background-color: var(--color-neutral-100); border-radius: var(--radius); box-shadow: var(--shadow-hover); width: 100%; max-width: 480px; max-height: 90vh; overflow-y: auto; animation: slideUp 0.2s ease; }
.modal-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.5rem; border-bottom: 1px solid var(--color-neutral-200); }
.modal-header h2 { font-size: 1.125rem; font-weight: 600; }
.modal-close { padding: 0.25rem; border-radius: var(--radius-sm); color: var(--color-neutral-500); transition: background-color 0.2s ease, color 0.2s ease; background: none; border: none; cursor: pointer; }
.modal-close:hover { background-color: var(--color-neutral-200); color: var(--color-neutral-900); }
.modal-form { padding: 1.5rem; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; font-size: 0.875rem; font-weight: 500; color: var(--color-neutral-700); margin-bottom: 0.375rem; }
.form-input { width: 100%; padding: 0.5rem 0.75rem; border: 1px solid var(--color-neutral-300); border-radius: var(--radius-sm); background-color: var(--color-neutral-100); color: var(--color-neutral-900); font-size: 0.875rem; transition: border-color 0.2s ease, box-shadow 0.2s ease; }
.form-input:focus { outline: none; border-color: transparent; box-shadow: 0 0 0 2px var(--color-primary); }
.form-hint { font-size: 0.75rem; color: var(--color-neutral-500); margin-top: 0.25rem; }
.modal-actions { display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.5rem; }
.btn-primary, .btn-secondary, .btn-danger { padding: 0.5rem 1rem; border-radius: var(--radius-sm); font-weight: 500; font-size: 0.875rem; cursor: pointer; transition: background-color 0.2s ease, opacity 0.2s ease; display: inline-flex; align-items: center; justify-content: center; gap: 0.5rem; border: none; }
.btn-primary { background-color: var(--color-primary); color: white; }
.btn-primary:hover:not(:disabled) { background-color: var(--color-primary-hover); }
.btn-secondary { background-color: var(--color-neutral-200); color: var(--color-neutral-700); }
.btn-secondary:hover:not(:disabled) { background-color: var(--color-neutral-300); }
.btn-danger { background-color: var(--color-danger); color: white; }
.btn-danger:hover:not(:disabled) { background-color: #c82333; }
.btn-primary:disabled, .btn-secondary:disabled, .btn-danger:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-icon { padding: 0.375rem; border-radius: var(--radius-sm); background: none; border: none; color: var(--color-neutral-500); cursor: pointer; transition: background-color 0.2s ease, color 0.2s ease; }
.btn-icon:hover { background-color: var(--color-neutral-200); color: var(--color-neutral-900); }

.share-result p { color: var(--color-neutral-600); margin-bottom: 0.75rem; }
.share-url-display { display: flex; gap: 0.5rem; margin-bottom: 0.75rem; }
.share-url-display .form-input { flex: 1; background-color: var(--color-neutral-50); }
.share-warning { font-size: 0.8125rem; color: var(--color-warning); margin: 0; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
</style>