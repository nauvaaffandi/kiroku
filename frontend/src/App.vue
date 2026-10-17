<template>
  <AppShell />
  <InstallPrompt />
  <PasswordOverlay :modelValue="authStore.showPasswordOverlay" @update:modelValue="authStore.showPasswordOverlay = $event" />
  <ErudaWidget />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { onBeforeMount } from 'vue'
import AppShell from '@/components/layout/AppShell.vue'
import InstallPrompt from '@/components/forms/InstallPrompt.vue'
import PasswordOverlay from '@/components/auth/PasswordOverlay.vue'
import ErudaWidget from '@/components/common/ErudaWidget.vue'
import { initApiClient } from '@/api/client'

const authStore = useAuthStore()

onBeforeMount(async () => {
  await initApiClient()
  authStore.init()
})
</script>

<style>
/* Page transitions */
.page-enter-active,
.page-leave-active {
  transition: opacity 0.2s ease;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>