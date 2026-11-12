<template>
  <AppShell />
   <InstallPrompt />
   <ErudaWidget />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { onBeforeMount } from 'vue'
import AppShell from '@/components/layout/AppShell.vue'
import InstallPrompt from '@/components/forms/InstallPrompt.vue'
import ErudaWidget from '@/components/common/ErudaWidget.vue'
import { initApiClient } from '@/api/client'

const authStore = useAuthStore()

onBeforeMount(async () => {
  await initApiClient()
  authStore.init()
})
</script>