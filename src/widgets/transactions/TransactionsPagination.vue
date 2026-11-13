<template>
  <div class="pagination" v-if="totalPages > 1">
    <BaseButton variant="secondary" @click="$emit('prev')" :disabled="currentPage === 1">
      {{ t('common.prev') }}
    </BaseButton>
    <span class="page-info">{{ t('common.page', { current: currentPage, total: totalPages }) }}</span>
    <BaseButton variant="secondary" @click="$emit('next')" :disabled="currentPage >= totalPages">
      {{ t('common.next') }}
    </BaseButton>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { BaseButton } from '@/components/ui'

const props = defineProps<{
  currentPage: number
  totalPages: number
}>()

const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.pagination {
  @apply flex items-center justify-center gap-4 p-4;
  border-top: 1px solid var(--border-primary);
}

.page-info {
  @apply text-sm;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .pagination {
    @apply p-3 flex-wrap gap-2;
  }
  .pagination button {
    @apply flex-1;
  }
  .page-info {
    @apply flex-1 text-center;
  }
}
</style>
