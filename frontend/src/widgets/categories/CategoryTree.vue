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
.categories-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.empty-state {
  text-align: center;
  padding: 3rem 1.5rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  color: var(--color-text-muted);
}

.empty-icon {
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state p {
  margin-bottom: 1rem;
}
</style>