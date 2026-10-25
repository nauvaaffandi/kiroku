<template>
  <BaseModal
    v-model="isOpen"
    :title="title"
    size="md"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <div class="totp-setup">
      <div class="totp-qr">
        <img :src="qrCode" alt="TOTP QR Code" />
        <p class="totp-secret">Secret: {{ secret }}</p>
      </div>
      <div class="form-group">
        <label>{{ codeLabel }}</label>
        <BaseInput
          v-model="code"
          type="text"
          placeholder="Enter 6-digit code"
          maxlength="6"
          autocomplete="one-time-code"
        />
      </div>
      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="isOpen = false">
          Cancel
        </BaseButton>
        <BaseButton type="submit" @click="handleSubmit">
          {{ submitLabel }}
        </BaseButton>
      </div>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { BaseModal, BaseInput, BaseButton } from '@/components/ui'

interface Props {
  modelValue: boolean
  title: string
  qrCode: string
  secret: string
  codeLabel?: string
  submitLabel?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [code: string]
}>()

const isOpen = ref(props.modelValue)
const code = ref('')

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val) code.value = ''
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = () => {
  emit('submit', code.value)
  isOpen.value = false
}
</script>

<style scoped>
@reference "tailwindcss";
.totp-setup {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.totp-qr {
  text-align: center;
}

.totp-qr img {
  max-width: 200px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
}

.totp-secret {
  font-family: monospace;
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin-top: 0.5rem;
  word-break: break-all;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.form-group label {
  font-weight: 500;
  color: var(--color-text);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}
</style>