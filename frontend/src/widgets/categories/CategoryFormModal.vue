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
.category-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.form-checkbox {
  flex-direction: row;
  align-items: center;
}

.form-checkbox .form-label {
  margin-bottom: 0;
}

.form-checkbox-input {
  width: 1rem;
  height: 1rem;
  accent-color: var(--color-primary);
}

.form-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text);
}

.form-input {
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  font-family: inherit;
  color: var(--color-text);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-focus);
}

input[type="color"].form-input {
  width: 3rem;
  height: 2.5rem;
  padding: 0;
  border-radius: var(--radius);
  cursor: pointer;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}
</style>