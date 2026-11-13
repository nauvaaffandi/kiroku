<template>
  <div class="shared-access-page">
    <header class="page-header">
      <div class="header-content">
        <h1>{{ t('shared.access.title') }}</h1>
        <button @click="showCreateModal = true" class="btn-primary" :aria-label="t('shared.access.createShare')">{{ t('shared.access.createShare') }}</button>
      </div>
    </header>
    <SharedAccessList
      :shares="shares"
      :copy-share-url="copyShareUrl"
      :revoke-share="revokeShare"
      :format-date="formatDate"
      @create="showCreateModal = true"
    />
    <CreateShareModal
      :accounts="accounts"
      :show-create-modal="showCreateModal"
      :loading="loading"
      :form="createForm"
      @close="closeCreateModal"
      @submit="createShare"
    />
    <ShareResultModal
      :share="createdShare"
      :full-share-url="fullShareUrl"
      :copy-full-share-url="copyFullShareUrl"
      @close="closeCreatedModal"
    />
  </div>
</template>

<script setup lang="ts">
import { useShare } from '@/widgets/shared/useShare'
import { CreateShareModal, ShareResultModal } from '@/widgets/shared'
import SharedAccessList from '@/components/shared/SharedAccessList.vue'

const {
  accounts, closeCreateModal, closeCreatedModal, copyFullShareUrl, copyShareUrl,
  createForm, createShare, createdShare, formatDate, fullShareUrl, loading,
  revokeShare, shares, showCreateModal, t
} = useShare()
</script>

<style scoped>
@reference "tailwindcss";

.shared-access-page { @apply min-h-screen flex flex-col; }
.page-header {
  @apply p-6;
  background-color: var(--bg-card);
  border-bottom: 1px solid var(--border-primary);
}
.header-content { @apply flex items-center justify-between flex-wrap gap-4; }
.header-content h1 {
  @apply text-2xl font-semibold;
  color: var(--text-primary);
}
.btn-primary {
  @apply py-2 px-4 rounded-[var(--radius-sm)] font-medium text-sm cursor-pointer transition-opacity duration-200 inline-flex items-center justify-center gap-2 border-none text-white;
  background-color: var(--accent);
}
.btn-primary:hover:not(:disabled) { background-color: var(--accent-hover); }
</style>
