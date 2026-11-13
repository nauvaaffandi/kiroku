<template>
  <div class="form-field">
    <label v-if="label" :for="id" class="form-label">{{ label }}</label>
    <div class="select-wrapper">
      <select
        :id="id"
        :value="modelValue"
        :disabled="disabled"
        :required="required"
        :aria-invalid="error ? 'true' : 'false'"
        :aria-describedby="error ? `${id}-error` : hint ? `${id}-hint` : undefined"
        class="form-select"
        @change="handleChange"
        @blur="handleBlur"
        @focus="handleFocus"
      >
        <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
        <option
          v-for="option in options"
          :key="option.value"
          :value="option.value"
          :disabled="option.disabled"
        >
          {{ option.label }}
        </option>
      </select>
      <svg class="select-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
        <polyline points="6 9 12 15 18 9"></polyline>
      </svg>
    </div>
    <p v-if="error" :id="`${id}-error`" class="form-error" role="alert">{{ error }}</p>
    <p v-else-if="hint" :id="`${id}-hint`" class="form-hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
interface Option {
  value: string
  label: string
  disabled?: boolean
}

interface Props {
  modelValue: string
  options: Option[]
  label?: string
  placeholder?: string
  disabled?: boolean
  required?: boolean
  error?: string
  hint?: string
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  required: false
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  blur: [event: FocusEvent]
  focus: [event: FocusEvent]
}>()

const id = props.id || `select-${Math.random().toString(36).slice(2)}`

const handleChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
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

.select-wrapper {
  @apply relative;
}

.form-select {
  @apply w-full px-3 py-2 text-sm pl-3 pr-10 font-sans rounded border appearance-none cursor-pointer transition-colors;
  color: var(--text-primary);
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  font-family: inherit;
}

.form-select:focus {
  @apply outline-none;
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent);
}

.form-select:disabled {
  @apply cursor-not-allowed;
  background-color: var(--bg-hover);
  color: var(--text-muted);
}

.select-icon {
  @apply absolute right-[0.75rem] top-[50%] w-[1rem] h-[1rem] pointer-events-none;
  transform: translateY(-50%);
  color: var(--text-muted);
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