<template>
  <div class="category-parent-row">
    <button class="parent-main" @click="emit('select', category)" :aria-label="t('categories.parentCategory')">
      <div class="parent-icon">
        <span v-if="category.icon">{{ category.icon }}</span>
      </div>
      <div class="parent-lead">
        <div class="parent-info">
          <span class="parent-name">{{ translate(category) }}</span>
          <span v-if="category.kind" class="kind-badge" :class="'kind-' + category.kind">{{ kindLabel }}</span>
        </div>
        <div class="parent-meta">
          <span v-if="category.archived" class="meta-chip archived">{{ t('categories.labelArchived') }}</span>
        </div>
      </div>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, withDefaults } from 'vue'
import { useI18n } from '@/locales'
import { useCategoryName } from '@/utils/categoryName'
import type { Category } from '@/types'

interface Props {
  category: Category & { _count?: number }
}

interface Emits {
  select: [category: Category]
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

.category-parent-row {
  @apply mb-2 border border-transparent rounded-xl;
  transition: border-color 0.15s ease, background-color 0.15s ease;
}

.category-parent-row:hover {
  border-color: var(--border-secondary);
  background-color: var(--bg-hover);
}

.parent-main {
  @apply flex items-center gap-3 min-h-[48px] w-full text-left bg-transparent border-none cursor-pointer outline-none px-3 py-2;
}

.parent-main:focus-visible {
  @apply ring-2 ring-offset-2 rounded-lg;
  ring-offset-color: var(--bg-primary);
  ring-color: var(--accent);
}

.parent-icon {
  @apply flex-shrink-0 w-10 h-10 rounded-lg flex items-center justify-center text-white;
  background-color: var(--accent);
}

.parent-lead {
  @apply flex flex-1 items-center gap-3 min-w-0;
}

.parent-info {
  @apply flex-1 min-w-0;
}

.parent-name {
  @apply truncate text-base font-semibold leading-snug;
  color: var(--text-primary);
}

.parent-meta {
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