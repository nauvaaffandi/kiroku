<template>
  <Transition name="modal">
    <div v-if="showCalculator" class="modal-overlay" @click.self="close">
      <div class="modal modal-sm calculator-modal" @click.stop>
        <div class="modal-header">
          <h2>Calculator</h2>
          <button @click="close" class="modal-close" aria-label="Close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
        <div class="calc-body">
          <div class="calc-display">{{ calcDisplay }}</div>
          <div class="calc-keypad">
            <button @click="calcClear" class="calc-key op">C</button>
            <button @click="calcAppend('%')" class="calc-key op">%</button>
            <button @click="calcAppend('/')" class="calc-key op">÷</button>
            <button @click="calcBackspace" class="calc-key op">⌫</button>
            
            <button @click="calcAppend('7')" class="calc-key">7</button>
            <button @click="calcAppend('8')" class="calc-key">8</button>
            <button @click="calcAppend('9')" class="calc-key">9</button>
            <button @click="calcAppend('*')" class="calc-key op">×</button>
            
            <button @click="calcAppend('4')" class="calc-key">4</button>
            <button @click="calcAppend('5')" class="calc-key">5</button>
            <button @click="calcAppend('6')" class="calc-key">6</button>
            <button @click="calcAppend('-')" class="calc-key op">-</button>
            
            <button @click="calcAppend('1')" class="calc-key">1</button>
            <button @click="calcAppend('2')" class="calc-key">2</button>
            <button @click="calcAppend('3')" class="calc-key">3</button>
            <button @click="calcAppend('+')" class="calc-key op">+</button>
            
            <button @click="calcAppend('0')" class="calc-key span-2">0</button>
            <button @click="calcAppend('.')" class="calc-key">.</button>
            <button @click="calcEvaluate" class="calc-key eq">=</button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useUIStore } from '@/stores/ui'

const props = defineProps<{
  modelValue: boolean
  currentAmount: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'result': [amount: number]
}>()

const uiStore = useUIStore()

const showCalculator = ref(false)
const calcDisplay = ref('0')

watch(() => props.modelValue, (val) => {
  showCalculator.value = val
  if (val) {
    calcDisplay.value = String(props.currentAmount || 0)
  }
}, { immediate: true })

const close = () => {
  showCalculator.value = false
  emit('update:modelValue', false)
}

const calcClear = () => { calcDisplay.value = '0' }
const calcAppend = (char: string) => {
  if (calcDisplay.value === '0' && !['+', '-', '*', '/', '%', '.'].includes(char)) {
    calcDisplay.value = char
  } else {
    calcDisplay.value += char
  }
}
const calcBackspace = () => {
  calcDisplay.value = calcDisplay.value.length > 1 ? calcDisplay.value.slice(0, -1) : '0'
}
const calcEvaluate = () => {
  try {
    const sanitized = calcDisplay.value.replace(/[^0-9+\-*/().%]/g, '')
    const res = Function(`'use strict'; return (${sanitized})`)()
    const amount = Math.round(Number(res) || 0)
    emit('result', amount)
    close()
  } catch (e) {
    uiStore.error('Invalid expression')
  }
}
</script>

<style scoped>
@reference "tailwindcss";
.calculator-modal { max-width: 300px; }
.calc-body { padding: 1rem; }
.calc-display {
  background-color: var(--color-neutral-900);
  color: white;
  font-size: 1.5rem;
  font-weight: 600;
  text-align: right;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-sm);
  margin-bottom: 1rem;
  font-variant-numeric: tabular-nums;
  overflow: hidden;
}
.calc-keypad { display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.5rem; }
.calc-key { padding: 0.75rem; font-size: 1rem; font-weight: 600; background-color: var(--color-neutral-200); color: var(--color-neutral-900); border-radius: var(--radius-sm); border: none; cursor: pointer; transition: background-color 0.2s ease; }
.calc-key:hover { background-color: var(--color-neutral-300); }
.calc-key.op { background-color: var(--color-neutral-300); color: var(--color-primary); }
.calc-key.eq { background-color: var(--color-primary); color: white; }
.calc-key.span-2 { grid-column: span 2; }
</style>