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
           <span class="category-icon default" v-else>
             <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
               <circle cx="12" cy="12" r="10"></circle>
             </svg>
           </span>
           <span class="category-name">{{ translate(category) }}</span>
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
import { useCategoryName } from '@/utils/categoryName'
import type { Category } from '@/types'

defineProps<{
  category: Category
  level: number
}>()

const { translate } = useCategoryName()

defineEmits<{
  edit: [category: Category]
  delete: [id: string]
}>()

const expanded = ref(true)
</script>

<style scoped>
@reference "tailwindcss";

.category-row {
  @apply rounded;
}

.tree-row {
  @apply grid grid-cols-[1fr_120px_80px_120px] gap-4 px-4 py-2 rounded items-center transition-colors border;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}

.tree-row:hover {
  background-color: var(--bg-hover);
}

.category-info {
  @apply flex items-center gap-2;
}

.tree-toggle {
  @apply w-6 h-6 flex items-center justify-center cursor-pointer shrink-0;
  color: var(--text-muted);
}

.tree-toggle.placeholder {
  @apply invisible;
}

.category-icon {
  @apply w-6 h-6 rounded flex items-center justify-center text-sm shrink-0 text-white;
}

.category-icon.default {
  color: var(--text-muted);
  @apply bg-transparent;
}

.category-name {
  @apply font-semibold;
  color: var(--text-primary);
}

.kind-badge {
  @apply inline-block px-2 py-0.5 rounded-full text-[0.625rem] font-bold capitalize;
}

.kind-badge.must {
  background-color: var(--error-light);
  color: var(--error);
}
.kind-badge.have {
  background-color: var(--warning-light);
  color: var(--warning);
}
.kind-badge.want {
  background-color: var(--info-light);
  color: var(--info);
}

.color-swatch {
  @apply inline-block w-4 h-4 rounded border align-middle;
  border-color: var(--border-primary);
}

.color-hex {
  @apply text-xs font-medium font-mono ml-2;
  color: var(--text-muted);
}

.btn-icon {
  @apply p-1.5 rounded transition-colors;
  color: var(--text-muted);
}

.btn-icon:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.btn-icon.text-danger:hover {
  background-color: var(--error-light);
  color: var(--error);
}

.category-children {
  @apply mt-1;
}

@media (max-width: 768px) {
  .tree-row {
    grid-template-columns: 1fr 80px 60px 100px;
    @apply text-sm;
  }
  .color-hex {
    @apply hidden;
  }
}
</style>