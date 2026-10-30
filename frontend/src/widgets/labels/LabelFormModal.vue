<template>
  <BaseModal
    v-model="isOpen"
    :title="isEditing ? t('common.edit') : t('labels.addLabel')"
    size="md"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <form @submit.prevent="handleSubmit" class="label-form">
      <div class="form-group">
        <label class="form-label">{{ t('labels.name') }} *</label>
        <BaseInput
          v-model="formData.name"
          type="text"
          required
          :placeholder="t('labels.name')"
        />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('labels.color') }}</label>
          <div class="color-picker-wrap">
            <input
              v-model="formData.color"
              type="color"
              class="color-input"
            />
            <BaseInput
              v-model="formData.color"
              type="text"
              placeholder="#6366f1"
            />
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('labels.icon') }}</label>
          <BaseInput
            v-model="formData.icon"
            type="text"
            placeholder="🏷️"
            maxlength="4"
          />
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">{{ t('labels.description') }}</label>
        <BaseInput
          v-model="formData.description"
          type="text"
          :placeholder="t('labels.description')"
        />
      </div>
      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="isOpen = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton type="submit">
          {{ t('labels.save') }}
        </BaseButton>
      </div>
    </form>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseModal, BaseInput, BaseButton } from '@/components/ui'
import { useI18n } from '@/locales'
import type { Label, CreateLabelRequest } from '@/types'

interface Props {
  modelValue: boolean
  label?: Label | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [data: CreateLabelRequest]
}>()

const { t } = useI18n()
const isOpen = ref(props.modelValue)
const isEditing = computed(() => !!props.label)

const formData = ref({
  name: '',
  color: '#6366f1',
  icon: '',
  description: ''
})

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val) {
    if (props.label) {
      formData.value = {
        name: props.label.name,
        color: props.label.color || '#6366f1',
        icon: props.label.icon || '',
        description: props.label.description || ''
      }
    } else {
      formData.value = {
        name: '',
        color: '#6366f1',
        icon: '',
        description: ''
      }
    }
  }
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = () => {
  if (!formData.value.name.trim()) return
  emit('submit', { ...formData.value })
  isOpen.value = false
}
</script>

<style scoped>
@reference "tailwindcss";

.label-form {
  @apply flex flex-col gap-4;
}

.form-group {
  @apply flex flex-col gap-1.5;
}

.form-row {
  @apply grid grid-cols-2 gap-4;
}

.form-label {
  @apply text-sm font-medium;
  color: var(--text-primary);
}

.color-picker-wrap {
  @apply flex items-center gap-2;
}

.color-input {
  @apply w-10 h-10 p-0 border rounded cursor-pointer;
  border-color: var(--border-primary);
  background-color: transparent;
}

.modal-actions {
  @apply flex justify-end gap-3 mt-4 pt-4 border-t;
  border-color: var(--border-primary);
}
</style>
