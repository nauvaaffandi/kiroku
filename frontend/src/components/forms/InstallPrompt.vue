<template>
  <div v-if="showInstallPrompt" class="install-banner" role="dialog" aria-label="Install app">
    <div class="install-content">
      <div class="install-icon">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
          <polyline points="7 10 12 15 17 10"></polyline>
          <line x1="12" y1="15" x2="12" y2="3"></line>
        </svg>
      </div>
      <div class="install-text">
        <strong>Install Kiroku</strong>
        <span>Add to home screen for quick access</span>
      </div>
    </div>
    <div class="install-actions">
      <button @click="dismissInstall" class="btn-ghost" aria-label="Dismiss">Not now</button>
      <button @click="installApp" class="btn-primary">Install</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useUIStore } from '@/stores/ui'

const uiStore = useUIStore()

const showInstallPrompt = ref(false)
let deferredPrompt: BeforeInstallPromptEvent | null = null

const handleBeforeInstallPrompt = (e: Event) => {
  e.preventDefault()
  deferredPrompt = e as BeforeInstallPromptEvent
  // Show prompt after a short delay if user hasn't dismissed recently
  const dismissed = localStorage.getItem('installDismissed')
  const dismissTime = dismissed ? parseInt(dismissed, 10) : 0
  const now = Date.now()
  // Show again after 7 days
  if (now - dismissTime > 7 * 24 * 60 * 60 * 1000) {
    showInstallPrompt.value = true
  }
}

const handleAppInstalled = () => {
  showInstallPrompt.value = false
  deferredPrompt = null
  uiStore.success('Kiroku installed!')
}

const installApp = async () => {
  if (!deferredPrompt) return
  deferredPrompt.prompt()
  const { outcome } = await deferredPrompt.userChoice
  if (outcome === 'accepted') {
    uiStore.success('Installing Kiroku...')
  }
  deferredPrompt = null
  showInstallPrompt.value = false
}

const dismissInstall = () => {
  showInstallPrompt.value = false
  localStorage.setItem('installDismissed', Date.now().toString())
}

onMounted(() => {
  window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  window.addEventListener('appinstalled', handleAppInstalled)
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  window.removeEventListener('appinstalled', handleAppInstalled)
})
</script>

<style scoped>
@reference "tailwindcss";
.install-banner {
  position: fixed;
  bottom: 1rem;
  left: 1rem;
  right: 1rem;
  max-width: 400px;
  margin: 0 auto;
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--color-neutral-200);
  padding: 1rem;
  z-index: 1000;
  animation: slideUp 0.3s ease;
}

.install-content {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.install-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: var(--radius);
  background: linear-gradient(135deg, var(--color-primary), var(--color-primary-light));
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.install-text {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.install-text strong {
  font-size: 0.9375rem;
  color: var(--color-neutral-900);
}

.install-text span {
  font-size: 0.8125rem;
  color: var(--color-neutral-500);
}

.install-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.btn-ghost {
  padding: 0.5rem 1rem;
  border-radius: var(--radius-sm);
  font-weight: 500;
  font-size: 0.875rem;
  cursor: pointer;
  background: none;
  border: 1px solid var(--color-neutral-300);
  color: var(--color-neutral-600);
  transition: background-color 0.2s ease, color 0.2s ease;
}

.btn-ghost:hover {
  background-color: var(--color-neutral-100);
  color: var(--color-neutral-900);
}

.btn-primary {
  padding: 0.5rem 1rem;
  border-radius: var(--radius-sm);
  font-weight: 500;
  font-size: 0.875rem;
  cursor: pointer;
  background-color: var(--color-primary);
  color: white;
  border: none;
  transition: background-color 0.2s ease;
}

.btn-primary:hover {
  background-color: var(--color-primary-hover);
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (min-width: 768px) {
  .install-banner {
    bottom: 2rem;
    left: auto;
    right: 2rem;
    max-width: 360px;
  }
}
</style>