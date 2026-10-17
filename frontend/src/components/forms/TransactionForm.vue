<template>
  <div class="transaction-form-container">
    <TransactionFormFieldsWidget
      :initialData="initialData"
      :isEditing="isEditing"
      :categories="categories"
      :accounts="accounts"
      @save="handleSave"
      @cancel="$emit('cancel')"
    />

    <TransactionCalculatorWidget
      v-model:modelValue="showCalculator"
      :currentAmount="formAmount"
      @result="handleCalculatorResult"
    />

    <TransactionMapWidget
      v-model:modelValue="showMapModal"
      :placeLat="formPlaceLat"
      :placeLng="formPlaceLng"
      :placeName="formPlaceName"
      @update:placeLat="updatePlaceLat"
      @update:placeLng="updatePlaceLng"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import TransactionFormFieldsWidget from '@/widgets/transaction-form/TransactionFormFieldsWidget.vue'
import TransactionCalculatorWidget from '@/widgets/transaction-form/TransactionCalculatorWidget.vue'
import TransactionMapWidget from '@/widgets/transaction-form/TransactionMapWidget.vue'
import type { Category, Account } from '@/types'

const props = defineProps<{
  initialData?: any
  isEditing?: boolean
  categories: Category[]
  accounts: Account[]
}>()

const emit = defineEmits(['save', 'cancel'])

const showCalculator = ref(false)
const showMapModal = ref(false)

// These refs need to be connected to the child component's form state
// We'll use a custom event or provide/inject pattern
const formAmount = ref(0)
const formPlaceLat = ref<number | null>(null)
const formPlaceLng = ref<number | null>(null)
const formPlaceName = ref('')

const handleSave = (payload: any) => {
  emit('save', payload)
}

const handleCalculatorResult = (_amount: number) => {
  showCalculator.value = false
}

const updatePlaceLat = (lat: number | null) => {
  formPlaceLat.value = lat
}

const updatePlaceLng = (lng: number | null) => {
  formPlaceLng.value = lng
}

watch(() => props.initialData, (newData) => {
  if (newData) {
    formAmount.value = newData.amount || 0
    formPlaceLat.value = newData.place_lat || null
    formPlaceLng.value = newData.place_lng || null
    formPlaceName.value = newData.place_name || ''
  }
}, { deep: true, immediate: true })
</script>

<style scoped>
.transaction-form-container { display: flex; flex-direction: column; gap: 1rem; }
</style>