<template>
  <div class="category-card" :class="{ 'has-children': category.children && category.children.length > 0 }">
    <div class="category-header">
      <div class="category-icon" :style="{ backgroundColor: category.color }">
        <span v-if="category.icon">{{ category.icon }}</span>
        <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path>
          <line x1="7" y1="7" x2="7.01" y2="7"></line>
        </svg>
      </div>
      <div class="category-badges">
        <div class="category-type-badge" :class="category.type">{{ category.type }}</div>
        <span v-if="category.archived" class="category-type-badge archived">Archived</span>
      </div>
    </div>
    <h3 class="category-name">{{ category.path || category.name }}</h3>
    <p class="category-count" v-if="category._count">{{ category._count }} transactions</p>
    <div class="category-actions">
      <BaseButton variant="secondary" size="sm" @click="$emit('edit', category)">
        Edit
      </BaseButton>
      <BaseButton variant="secondary" size="sm" @click="$emit('create-child', category)">
        Add Child
      </BaseButton>
      <BaseButton variant="danger" size="sm" @click="$emit('delete', category.id)">
        Delete
      </BaseButton>
    </div>

    <div v-if="category.children && category.children.length > 0" class="category-children">
      <CategoryCard
        v-for="child in category.children"
        :key="child.id"
        :category="child"
        @edit="$emit('edit', $event)"
        @create-child="$emit('create-child', $event)"
        @delete="$emit('delete', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import type { Category } from '@/types'

interface Props {
  category: Category & { _count?: number }
}

defineProps<Props>()

defineEmits<{
  edit: [category: Category]
  'create-child': [parent: Category]
  delete: [id: string]
}>()
</script>

<style scoped>
@reference "tailwindcss";
.category-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); @apply p-4 mb-2 rounded-2xl; }
.category-card.has-children { @apply pb-0 rounded-b-none; }
.category-header { @apply flex justify-between items-center; }
.category-icon { @apply w-10 h-10 rounded flex items-center justify-center text-white text-xl; }
.category-badges { @apply flex gap-2 items-center; }
.category-type-badge { background-color: var(--bg-primary); color: var(--text-muted); @apply text-xs px-2 py-0.5 rounded-full capitalize; }
.category-type-badge.income { background-color: var(--success-light); color: var(--success); }
.category-type-badge.expense { background-color: var(--error-light); color: var(--error); }
.category-type-badge.transfer { background-color: var(--accent-light); color: var(--accent); }
.category-type-badge.archived { background-color: var(--warning-light); color: var(--warning); }
.category-name { color: var(--text-primary); @apply text-base font-semibold mt-2 mb-1; }
.category-count { color: var(--text-muted); @apply text-[0.8125rem] m-0; }
.category-actions { @apply flex gap-2 mt-3 pt-3 border-t; border-color: var(--border-primary); }
.category-children { @apply ml-6 border-l-2 pl-4 mt-2; border-color: var(--border-primary); }
</style>