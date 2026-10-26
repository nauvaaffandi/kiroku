<template>
  <BaseModal
    v-model="isOpen"
    :title="isEditing ? 'Edit Category' : (isChild ? 'New Subcategory' : 'New Category')"
    size="lg"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <form @submit.prevent="handleSubmit" class="category-form">
      <div class="form-group">
        <label class="form-label">Name *</label>
        <BaseInput
          v-model="formData.name"
          type="text"
          required
          placeholder="Category name"
        />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">Type *</label>
          <BaseSelect
            v-model="formData.type"
            :options="typeOptions"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label">Parent</label>
          <BaseSelect
            v-model="formData.parent_id"
            :options="parentOptions"
            placeholder="No parent (root category)"
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">Icon</label>
          <BaseInput
            v-model="formData.icon"
            type="text"
            placeholder="Emoji or leave blank"
            maxlength="4"
          />
        </div>
        <div class="form-group">
          <label class="form-label">Color</label>
          <input
            v-model="formData.color"
            type="color"
            class="form-input"
          />
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">Sort Order</label>
        <BaseInput
          v-model="formData.sort_order"
          type="number"
          :min="0"
          placeholder="0"
        />
      </div>
      <div class="form-group form-checkbox">
        <label class="form-label">Archived</label>
        <input
          v-model="formData.archived"
          type="checkbox"
          class="form-checkbox-input"
        />
      </div>
      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="isOpen = false">
          Cancel
        </BaseButton>
        <BaseButton type="submit">
          {{ isEditing ? 'Save Changes' : (isChild ? 'Create Subcategory' : 'Create Category') }}
        </BaseButton>
      </div>
    </form>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseModal, BaseInput, BaseSelect, BaseButton } from '@/components/ui'
import type { Category, CategoryType } from '@/types'

interface Props {
  modelValue: boolean
  category?: Category | null
  parent?: Category | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [data: Partial<Category>]
}>()

const isOpen = ref(props.modelValue)
const isEditing = computed(() => !!props.category)
const isChild = computed(() => !!props.parent && !props.category)

interface FormData {
  name: string
  type: CategoryType
  parent_id: string
  icon: string
  color: string
  sort_order: string
  archived: boolean
}

const defaultForm = (): FormData => ({
  name: '',
  type: 'expense',
  parent_id: props.parent?.id || '',
  icon: '',
  color: '#2563eb',
  sort_order: '0',
  archived: false
})

const formData = ref<FormData>(defaultForm())

const typeOptions = [
  { value: 'income', label: 'Income' },
  { value: 'expense', label: 'Expense' },
  { value: 'transfer', label: 'Transfer' }
]

const parentOptions = computed(() => (props.category
  ? [] // Can't be parent of itself
  : [])) // We'll pass parent options from parent if needed

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val && props.category) {
    formData.value = {
      name: props.category.name,
      type: props.category.type,
      parent_id: props.category.parent_id || '',
      icon: props.category.icon || '',
      color: props.category.color || '#2563eb',
      sort_order: String(props.category.sort_order || 0),
      archived: props.category.archived || false
    }
  } else if (val) {
    formData.value = defaultForm()
  }
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = () => {
  emit('submit', {
    name: formData.value.name,
    type: formData.value.type,
    parent_id: formData.value.parent_id || null,
    icon: formData.value.icon,
    color: formData.value.color,
    sort_order: Number(formData.value.sort_order),
    archived: formData.value.archived
  })
  isOpen.value = false
}
</script>

<style scoped>
@reference "tailwindcss";
.category-form { @apply flex flex-col gap-5; }
.form-row { @apply grid grid-cols-2 gap-4; }
.form-group { @apply flex flex-col gap-1.5; }
.form-checkbox { @apply flex-row items-center; }
.form-checkbox .form-label { @apply mb-0; }
.form-checkbox-input { accent-color: var(--accent); @apply w-4 h-4; }
.form-label {
  color: var(--text-primary);
  @apply text-[0.8125rem] font-medium;
}
.form-input {
  color: var(--text-primary);
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  @apply px-3 py-2 text-sm font-sans rounded transition-[border-color,box-shadow] duration-150 ease-in-out;
}
.form-input:focus {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px var(--accent-light);
  @apply outline-none;
}
input[type="color"].form-input {
  @apply w-12 h-10 p-0 cursor-pointer rounded;
}
.modal-actions { @apply flex justify-end gap-3 mt-4; }
</style>