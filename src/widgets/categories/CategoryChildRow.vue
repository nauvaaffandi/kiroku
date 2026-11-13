<template>
  <div class="category-child-row">
    <div class="child-main">
      <div class="child-lead">
        <div class="child-icon" :style="{ backgroundColor: category.color || 'var(--accent)' }">
          <span v-if="category.icon">{{ category.icon }}</span>
          <Tag v-else :size="18" />
        </div>
        <div class="child-info">
          <span class="child-name">{{ translate(category) }}</span>
          <span v-if="category.kind" class="kind-badge" :class="'kind-' + category.kind">{{ kindLabel }}</span>
        </div>
        <div class="child-meta" v-if="category.archived">
          <span class="meta-chip archived">{{ t('categories.labelArchived') }}</span>
        </div>
      </div>
      <div class="child-actions">
        <button type="button" class="action-btn" @click="emit('edit', category)" :aria-label="t('categories.edit')">
          <Pencil :size="16" />
        </button>
        <button type="button" class="action-btn danger" @click="emit('delete', category.id)" :aria-label="t('categories.delete')">
          <Trash2 :size="16" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, withDefaults } from 'vue'
import { Pencil, Trash2, Tag } from '@lucide/vue'
import { useI18n } from '@/locales'
import { useCategoryName } from '@/utils/categoryName'
import type { Category } from '@/types'

interface Props {
  category: Category & { _count?: number }
}

interface Emits {
  edit: [category: Category]
  delete: [id: string]
}

const props = withDefaults(defineProps<Props>(), {})

const emit = defineEmits<Emits>()

const { t } = useI18n()
const { translate } = useCategoryName()

const kindLabel = computed(() => {
  if (!props.category.kind) return ''
  const key = props.category.kind.charAt(0).toUpperCase() + props.category.kind.slice(1)
  return t(`categories.kind${key}`)
})
</script>

<style scoped>
@reference "tailwindcss";

.category-child-row {
  @apply border border-transparent rounded-xl mb-1.5;
  transition: background-color 0.15s ease, border-color 0.15s ease;
}

.category-child-row:hover {
  background-color: var(--bg-hover);
  border-color: var(--border-secondary);
}

.child-main {
  @apply flex items-center gap-3 min-h-[44px] px-3 py-2;
}

.child-lead {
  @apply flex flex-1 items-center gap-2 min-w-0;
}

.child-icon {
  @apply flex-shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-white text-sm;
}

.child-info {
  @apply flex-1 min-w-0;
}

.child-name {
  @apply truncate text-sm font-semibold leading-snug;
  color: var(--text-primary);
}

.child-meta {
  @apply flex flex-wrap items-center gap-1 mt-0.5;
}

.meta-chip {
  @apply inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[0.75rem] font-medium leading-none;
  color: var(--text-muted);
}

.meta-chip.archived {
  background-color: var(--warning-light);
  color: var(--warning);
}

.child-actions {
  @apply flex-shrink-0 ml-2 flex items-center gap-1;
}

.action-btn {
  @apply p-1.5 rounded-lg transition-colors flex items-center justify-center;
  color: var(--text-muted);
  background-color: transparent;
  border: none;
  cursor: pointer;
  min-width: 36px;
  min-height: 36px;
}

.action-btn:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.action-btn.danger:hover {
  background-color: var(--error-light);
  color: var(--error);
}

.action-btn:focus-visible {
  @apply ring-2 ring-offset-2 rounded-lg;
  ring-offset-color: var(--bg-primary);
  ring-color: var(--accent);
}

.kind-badge {
  @apply inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[0.7rem] font-semibold leading-none uppercase tracking-tighter;
  color: var(--text-primary);
}

.kind-badge.kind-must {
  background-color: var(--error-light);
}

.kind-badge.kind-have {
  background-color: var(--warning-light);
}

.kind-badge.kind-want {
  background-color: var(--info-light);
}
</style>