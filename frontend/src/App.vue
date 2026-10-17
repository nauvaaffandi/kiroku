<template>
  <AppShell />
  <InstallPrompt />
  <PasswordOverlay :modelValue="authStore.showPasswordOverlay" @update:modelValue="authStore.showPasswordOverlay = $event" />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useBackendUrlStore } from '@/stores/backendUrl'
import { onBeforeMount } from 'vue'
import AppShell from '@/components/layout/AppShell.vue'
import InstallPrompt from '@/components/forms/InstallPrompt.vue'
import PasswordOverlay from '@/components/auth/PasswordOverlay.vue'

const authStore = useAuthStore()
const backendUrlStore = useBackendUrlStore()

onBeforeMount(async () => {
  await backendUrlStore.fetchUrl()
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