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
  @apply flex flex-col gap-6;
}

.totp-qr {
  @apply text-center;
}

.totp-qr img {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  @apply max-w-[200px] mx-auto;
}

.totp-secret {
  color: var(--text-muted);
  @apply font-mono text-sm mt-2 break-all;
}

.form-group {
  @apply flex flex-col gap-1.5;
}

.form-group label {
  color: var(--text-primary);
  @apply font-medium;
}

.modal-actions {
  @apply flex justify-end gap-3 mt-4;
}
</style>