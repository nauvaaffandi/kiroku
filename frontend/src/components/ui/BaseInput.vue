<template>
  <div class="form-field">
    <label v-if="label" :for="id" class="form-label">{{ label }}</label>
    <div class="input-wrapper">
      <span v-if="prefix" class="input-prefix">{{ prefix }}</span>
      <input
        :id="id"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :required="required"
        :aria-invalid="error ? 'true' : 'false'"
        :aria-describedby="error ? `${id}-error` : hint ? `${id}-hint` : undefined"
        class="form-input"
        @input="handleInput"
        @blur="handleBlur"
        @focus="handleFocus"
      />
      <span v-if="suffix" class="input-suffix">{{ suffix }}</span>
    </div>
    <p v-if="error" :id="`${id}-error`" class="form-error" role="alert">{{ error }}</p>
    <p v-else-if="hint" :id="`${id}-hint`" class="form-hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string
  label?: string
  type?: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'search' | 'date' | 'datetime-local' | 'time'
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  required?: boolean
  error?: string
  hint?: string
  prefix?: string
  suffix?: string
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false,
  readonly: false,
  required: false
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  blur: [event: FocusEvent]
  focus: [event: FocusEvent]
}>()

const id = props.id || `input-${Math.random().toString(36).slice(2)}`

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}

const handleBlur = (event: FocusEvent) => {
  emit('blur', event)
}

const handleFocus = (event: FocusEvent) => {
  emit('focus', event)
}
</script>

<style scoped>
@reference "tailwindcss";

.form-field {
  @apply flex flex-col gap-1;
}

.form-label {
  @apply text-xs font-medium;
  color: var(--text-primary);
}

.input-wrapper {
  @apply relative flex items-center;
}

.form-input {
  @apply w-full px-3 py-2 text-sm font-sans rounded border transition-colors;
  color: var(--text-primary);
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  font-family: inherit;
}

.form-input:focus {
  @apply outline-none;
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent);
}

.form-input:disabled {
  @apply cursor-not-allowed;
  background-color: var(--bg-hover);
  color: var(--text-muted);
}

.form-input::placeholder {
  color: var(--text-muted);
}

.input-prefix,
.input-suffix {
  @apply absolute flex items-center pointer-events-none text-[0.875rem];
  color: var(--text-muted);
}

.input-prefix {
  @apply left-[0.75rem];
}

.input-suffix {
  @apply right-[0.75rem];
}

.form-input:has(+ .input-suffix) {
  @apply pr-[2.5rem];
}

.form-error {
  @apply m-0 text-[0.75rem];
  color: var(--error);
}

.form-hint {
  @apply m-0 text-[0.75rem];
  color: var(--text-muted);
}
</style>