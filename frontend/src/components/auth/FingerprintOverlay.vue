<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-[1000] flex items-end" @click="handleBackdropClick">
    <div class="w-full bg-white rounded-t-[24px] p-[1.5rem_2rem_2.5rem_2rem] shadow-[0_-10px_30px_rgba(0,0,0,0.15)] flex flex-col items-center" @click.stop>
      <div class="w-[40px] h-[4px] bg-[#ccc] rounded-[2px] mb-[1.5rem]"></div>
      
      <div class="text-center mb-[2rem]">
        <h2 class="text-lg font-semibold text-[#212529] mb-[0.375rem]">{{ t('auth.biometricAuth') }}</h2>
        <p class="text-sm text-[#6c757d]">{{ t('auth.touchSensor') }}</p>
      </div>

      <div class="relative w-[96px] h-[96px] rounded-full bg-[#e7f1fc] flex items-center justify-center cursor-pointer mb-[1.5rem] hover:scale-105 transition-transform duration-200" @click="triggerScan">
        <div class="absolute inset-[-4px] rounded-full border-2 border-dashed border-[#4a90d9] opacity-40" :class="{ 'animate-spin': isScanning }"></div>
        <div class="text-[#4a90d9] z-2" :class="{ 'animate-pulse': isScanning }">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 10a2 2 0 0 0-2 2c0 1.02-.1 2.51-.81 3.5"></path>
            <path d="M15.5 13.5c0-1.8-1.5-3.5-3.5-3.5s-3.5 1.7-3.5 3.5"></path>
            <path d="M19 11c0-3.86-3.14-7-7-7s-7 3.14-7 7c0 1.55.3 3 .83 4.33"></path>
            <path d="M21.5 14C22 12.8 22.5 11.5 22.5 10c0-5.52-4.48-10-10-10S2.5 4.48 2.5 10c0 2.21.73 4.25 1.97 5.88"></path>
            <path d="M12 18v2"></path>
            <path d="M9 19.5c0 .83.67 1.5 1.5 1.5h3c.83 0 1.5-.67 1.5-1.5V18"></path>
          </svg>
        </div>
      </div>

      <div v-if="errorMessage" class="px-[0.5rem] py-[1rem] bg-[#f8d7da] text-[#dc3545] rounded-[8px] text-sm text-center mb-[1rem]">
        {{ errorMessage }}
      </div>

      <div class="w-full flex justify-center">
        <button class="px-[0.75rem] py-[1.5rem] bg-[#e7f1fc] text-[#4a90d9] font-medium rounded-[10px] border border-[#4a90d9] cursor-pointer w-full max-w-[200px] mb-[0.75rem] hover:bg-[#4a90d9] hover:text-white transition-colors duration-200" @click="triggerScan">{{ t('auth.tryAgain') }}</button>
        <button class="px-[0.75rem] py-[1.5rem] bg-[#f1f3f5] text-[#495057] font-medium rounded-[10px] border-0 cursor-pointer w-full max-w-[200px] hover:bg-[#e9ecef] transition-colors duration-200" @click="close">{{ t('auth.usePin') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from '@/locales'

const { t } = useI18n()

interface Props {
  isOpen: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  authenticate: []
}>()

const isScanning = ref(false)
const errorMessage = ref<string | null>(null)

const close = () => {
  emit('close')
}

const handleBackdropClick = () => {
  close()
}

const triggerScan = async () => {
  if (isScanning.value) return
  isScanning.value = true
  errorMessage.value = null

  try {
    // Emit authenticate event to parent which will call the biometric composable
    await emit('authenticate')
  } catch (err: any) {
    errorMessage.value = err.message || t('auth.biometricFailed')
  } finally {
    isScanning.value = false
  }
}

watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    // Auto trigger scan when opened
    setTimeout(() => {
      triggerScan()
    }, 400)
  } else {
    // Reset error when closed
    errorMessage.value = null
    isScanning.value = false
  }
})
</script>
