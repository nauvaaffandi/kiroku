<template>
  <div class="shared-access-content">
    <div class="shares-list" v-if="shares.length > 0">
      <div v-for="share in shares" :key="share.id" class="share-card">
        <div class="share-info">
          <div class="share-account">
            <span class="account-name">{{ share.account_name }}</span>
            <span class="account-badge">{{ share.account_id }}</span>
          </div>
          <div class="share-details">
            <div class="detail-item">
              <span class="detail-label">{{ t('shared.access.link') }}</span>
              <div class="token-display">
                <span class="token">{{ share.token }}</span>
                <button @click="copyShareUrl(share.token)" class="btn-icon" :title="t('shared.access.copyFullLink')">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                  </svg>
                </button>
              </div>
            </div>
            <div class="detail-item">
              <span class="detail-label">{{ t('shared.access.created') }}</span>
              <span class="detail-value">{{ formatDate(share.created_at) }}</span>
            </div>
            <div class="detail-item" v-if="share.expires_at">
              <span class="detail-label">{{ t('shared.access.expires') }}</span>
              <span class="detail-value">{{ formatDate(share.expires_at) }}</span>
            </div>
            <div class="detail-item" v-if="share.last_accessed_at">
              <span class="detail-label">{{ t('shared.access.lastAccessed') }}</span>
            </div>
            <div class="detail-item" v-else>
              <span class="detail-label">{{ t('shared.access.lastAccessed') }}</span>
              <span class="detail-value text-muted">{{ t('shared.access.never') }}</span>
            </div>
          </div>
        </div>
        <div class="share-actions">
          <button @click="copyShareUrl(share.token)" class="btn-secondary btn-sm">{{ t('shared.access.copyLink') }}</button>
          <button @click="revokeShare(share.id)" class="btn-danger btn-sm">{{ t('shared.access.revoke') }}</button>
        </div>
      </div>
    </div>
    <div v-else class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
        <polyline points="22 6 12 16 10 14"></polyline>
      </svg>
      <h3>{{ t('shared.access.noShares') }}</h3>
      <p>{{ t('shared.access.noSharesDesc') }}</p>
      <button @click="$emit('create')" class="btn-primary">{{ t('shared.access.createShare') }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import type { Share } from '@/widgets/shared/useShare'

const props = defineProps<{
  shares: Share[]
  copyShareUrl: (token: string) => void
  revokeShare: (id: string) => void
  formatDate: (dateStr?: string | null) => string
}>()
const emit = defineEmits<{ create: [] }>()
const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.shared-access-content { @apply flex-1 p-6 max-w-[900px] mx-auto w-full; }
.shares-list { @apply flex flex-col gap-4; }
.share-card {
  @apply rounded-[var(--radius)] p-6 flex items-start justify-between gap-6 flex-wrap shadow-[var(--shadow)];
  background-color: var(--bg-card);
}
.share-info { @apply flex-1 min-w-[280px]; }
.share-account { @apply mb-4; }
.account-name { @apply font-semibold text-lg; color: var(--text-primary); }
.account-badge { @apply text-xs ml-2 font-mono; color: var(--text-muted); }
.share-details { @apply grid gap-3; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); }
.detail-item { @apply flex flex-col gap-0.5; }
.detail-label { @apply text-xs font-semibold uppercase tracking-wider; color: var(--text-muted); }
.detail-value { @apply text-sm; color: var(--text-primary); }
.detail-value.text-muted { color: var(--text-muted); }
.token-display { @apply flex items-center gap-2; }
.token { @apply font-mono text-[0.8125rem] px-2 py-1 rounded-[var(--radius-sm)] flex-1 break-all; color: var(--text-secondary); background-color: var(--bg-tertiary); }
.share-actions { @apply flex gap-2 shrink-0; }
.btn-sm { @apply py-1.5 px-3 text-[0.8125rem]; }
.btn-secondary { @apply text-[var(--text-secondary)]; background-color: var(--bg-tertiary); border: 1px solid var(--border-primary); }
.btn-secondary:hover:not(:disabled) { background-color: var(--bg-hover); }
.btn-danger { @apply text-white; background-color: var(--error); border: 1px solid var(--border-primary); }
.btn-danger:hover:not(:disabled) { background-color: var(--error-light); }
.btn-icon { @apply p-1.5 rounded-[var(--radius-sm)] bg-transparent border-none cursor-pointer transition-colors duration-200; color: var(--text-muted); }
.btn-icon:hover { @apply text-[var(--text-primary)]; background-color: var(--bg-hover); }
.empty-state { @apply text-center py-16 px-8 rounded-[var(--radius)] shadow-[var(--shadow)]; background-color: var(--bg-card); color: var(--text-muted); }
.empty-state svg { @apply mb-4 opacity-50; }
.empty-state h3 { @apply text-lg font-semibold mb-2; color: var(--text-primary); }
.empty-state p { @apply text-sm mb-6; }
.btn-primary { @apply py-2 px-4 rounded-[var(--radius-sm)] font-medium text-sm cursor-pointer transition-opacity duration-200 inline-flex items-center justify-center gap-2 border-none text-white; background-color: var(--accent); }
.btn-primary:hover:not(:disabled) { background-color: var(--accent-hover); }
</style>
