<template>
  <header class="page-header">
    <div>
      <h1>{{ t('transactions.title') }}</h1>
      <p class="page-subtitle">{{ t('transactions.addTransaction') }}</p>
    </div>
    <div class="header-actions">
      <BaseButton variant="secondary" @click="$emit('export-csv')" :disabled="loading">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
          <polyline points="7 10 12 15 17 10"></polyline>
          <line x1="12" y1="15" x2="12" y2="3"></line>
        </svg>
        {{ t('transactions.exportCsv') }}
      </BaseButton>
      <BaseButton variant="secondary" @click="$emit('export-excel')" :disabled="loading">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
          <polyline points="14 2 14 8 20 8"></polyline>
        </svg>
        {{ t('transactions.exportExcel') }}
      </BaseButton>
      <BaseButton variant="secondary" @click="$emit('scan-receipt')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 18a5 5 0 0 0-10 0"></path>
          <path d="M12 2a10 10 0 1 0 10 10"></path>
        </svg>
        {{ t('transactions.scanReceipt') }}
      </BaseButton>
      <BaseButton @click="$emit('create')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        {{ t('transactions.addTransaction') }}
      </BaseButton>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { BaseButton } from '@/components/ui'

const props = defineProps<{
  loading: boolean
}>()

const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.page-header {
  @apply flex items-start justify-between p-4 flex-wrap gap-3;
  background-color: var(--bg-card);
  border-bottom: 1px solid var(--border-primary);
}

.page-header h1 {
  @apply text-xl font-semibold;
  color: var(--text-primary);
}

.page-subtitle {
  @apply text-sm mt-1;
  color: var(--text-muted);
}

.header-actions {
  @apply flex gap-2 flex-wrap;
}

.header-actions svg {
  @apply mr-2;
}

@media (max-width: 768px) {
  .page-header {
    @apply flex-col items-stretch p-3;
  }
  .header-actions {
    @apply justify-stretch flex-col;
  }
  .header-actions button {
    @apply w-full justify-center;
  }
}
</style>
