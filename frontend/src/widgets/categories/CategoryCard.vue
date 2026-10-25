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
.category-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 1rem;
  margin-bottom: 0.5rem;
}

.category-card.has-children {
  padding-bottom: 0;
  border-bottom-left-radius: 0;
  border-bottom-right-radius: 0;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.category-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.25rem;
}

.category-badges {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.category-type-badge {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  background: var(--color-background);
  color: var(--color-text-muted);
  text-transform: capitalize;
}

.category-type-badge.income {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.category-type-badge.expense {
  background: var(--color-danger-bg);
  color: var(--color-danger);
}

.category-type-badge.transfer {
  background: var(--color-primary-bg);
  color: var(--color-primary);
}

.category-type-badge.archived {
  background: var(--color-warning-bg);
  color: var(--color-warning);
}

.category-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0.5rem 0 0.25rem;
}

.category-count {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  margin: 0;
}

.category-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--color-border);
}

.category-children {
  margin-left: 1.5rem;
  border-left: 2px solid var(--color-border);
  padding-left: 1rem;
  margin-top: 0.5rem;
}
</style>