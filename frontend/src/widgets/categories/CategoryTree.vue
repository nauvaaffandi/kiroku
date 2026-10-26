<template>
  <div class="categories-tree">
    <div v-if="categories.length === 0" class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="empty-icon">
        <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path>
        <line x1="7" y1="7" x2="7.01" y2="7"></line>
      </svg>
      <p>No categories yet</p>
      <BaseButton @click="$emit('create')">Add Category</BaseButton>
    </div>

    <div v-else class="categories-list">
      <CategoryCard
        v-for="category in categories"
        :key="category.id"
        :category="category"
        @edit="$emit('edit', $event)"
        @create-child="$emit('create-child', $event)"
        @delete="$emit('delete', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import CategoryCard from './CategoryCard.vue'
import type { Category } from '@/types'

defineProps<{
  categories: Category[]
}>()

defineEmits<{
  edit: [category: Category]
  'create-child': [parent: Category]
  delete: [id: string]
  create: []
}>()
</script>

<style scoped>
@reference "tailwindcss";
.categories-list { @apply flex flex-col gap-2; }
.empty-state {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  color: var(--text-muted);
  @apply text-center py-12 px-6 rounded-2xl;
}
.empty-icon { @apply mb-4 opacity-50 mx-auto; }
.empty-state p { @apply mb-4; }
</style>