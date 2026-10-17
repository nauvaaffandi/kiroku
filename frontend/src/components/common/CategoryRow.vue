<template>
  <div class="category-row" :style="{ paddingLeft: level * 24 + 'px' }">
    <div class="tree-row">
      <div class="tree-col name">
        <span class="tree-toggle" v-if="category.children.length > 0" @click="expanded = !expanded">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline :points="expanded ? '6 9 12 15 18 9' : '9 6 15 12 9 18'"></polyline>
          </svg>
        </span>
        <span class="tree-toggle placeholder" v-else></span>
        <div class="category-info">
          <span class="category-icon" :style="{ backgroundColor: category.color }" v-if="category.icon">{{ category.icon }}</span>
          <span class="category-icon default" v-else>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
            </svg>
          </span>
          <span class="category-name">{{ category.name }}</span>
        </div>
      </div>
      <div class="tree-col kind" v-if="category.kind">
        <span class="kind-badge" :class="category.kind">{{ category.kind }}</span>
      </div>
      <div class="tree-col kind" v-else>—</div>
      <div class="tree-col color">
        <span class="color-swatch" :style="{ backgroundColor: category.color }"></span>
        <span class="color-hex">{{ category.color }}</span>
      </div>
      <div class="tree-col actions">
        <button @click="$emit('edit', category)" class="btn-icon" aria-label="Edit">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
          </svg>
        </button>
        <button @click="$emit('delete', category.id)" class="btn-icon text-danger" aria-label="Delete">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
        </button>
      </div>
    </div>
    <div v-if="expanded && category.children.length > 0" class="category-children">
      <CategoryRow v-for="child in category.children" :key="child.id" :category="child" :level="level + 1" @edit="(cat) => $emit('edit', cat)" @delete="(id) => $emit('delete', id)" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Category {
  id: string
  parent_id: string | null
  name: string
  icon: string
  color: string
  kind: string | null
  sort_order: number
  depth: number
  path: string
  children: Category[]
}

defineProps<{
  category: Category
  level: number
}>()

defineEmits<{
  edit: [category: Category]
  delete: [id: string]
}>()

const expanded = ref(true)
</script>

<style scoped>
.category-row {
  border-radius: var(--radius-sm);
}

.tree-row {
  display: grid;
  grid-template-columns: 1fr 120px 80px 120px;
  gap: 1rem;
  padding: 0.5rem 1rem;
  background-color: var(--color-neutral-100);
  border-radius: var(--radius-sm);
  align-items: center;
  transition: background-color 0.2s ease;
}

.tree-row:hover {
  background-color: var(--color-neutral-200);
}

.category-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tree-toggle {
  width: 1.5rem;
  height: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-neutral-500);
  cursor: pointer;
  flex-shrink: 0;
}

.tree-toggle.placeholder {
  visibility: hidden;
}

.category-icon {
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 0.875rem;
  flex-shrink: 0;
}

.category-icon.default {
  color: var(--color-neutral-400);
  background-color: transparent;
}

.category-name {
  font-weight: 500;
  color: var(--color-neutral-900);
}

.kind-badge {
  display: inline-block;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.625rem;
  font-weight: 600;
  text-transform: capitalize;
}

.kind-badge.must {
  background-color: rgba(220, 53, 69, 0.15);
  color: var(--color-danger);
}
.kind-badge.have {
  background-color: rgba(255, 193, 7, 0.15);
  color: var(--color-warning);
}
.kind-badge.want {
  background-color: rgba(23, 162, 184, 0.15);
  color: var(--color-info);
}

.color-swatch {
  display: inline-block;
  width: 1rem;
  height: 1rem;
  border-radius: 0.25rem;
  border: 1px solid var(--color-neutral-300);
  vertical-align: middle;
}

.color-hex {
  font-size: 0.75rem;
  color: var(--color-neutral-500);
  font-family: monospace;
  margin-left: 0.5rem;
}

.btn-icon {
  padding: 0.375rem;
  border-radius: var(--radius-sm);
  color: var(--color-neutral-500);
  transition: background-color 0.2s ease, color 0.2s ease;
}
.btn-icon:hover { background-color: var(--color-neutral-200); color: var(--color-neutral-900); }
.btn-icon.text-danger:hover { background-color: rgba(220, 53, 69, 0.15); color: var(--color-danger); }

.category-children {
  margin-top: 0.25rem;
}

@media (max-width: 768px) {
  .tree-row {
    grid-template-columns: 1fr 80px 60px 100px;
    font-size: 0.875rem;
  }
  .color-hex { display: none; }
}
</style>