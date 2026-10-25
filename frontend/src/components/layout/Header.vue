<template>
  <div class="header-container">
    <div class="header-left">
      <button @click="$emit('toggle-sidebar')" class="menu-btn" aria-label="Toggle sidebar">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="3" y1="12" x2="21" y2="12"></line>
          <line x1="3" y1="6" x2="21" y2="6"></line>
          <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
      </button>
      <h1 class="page-title">{{ pageTitle }}</h1>
    </div>
    
    <div class="header-right">
      <div class="currency-selector">
        <select :value="settingsStore.settings.base_currency || 'USD'" @change="changeCurrency" aria-label="Base currency">
          <option v-for="c of currencies" :key="c" :value="c">{{ c }}</option>
        </select>
      </div>
      

      
      <div class="more-menu-container">
        <button @click="toggleMoreMenu" class="icon-btn" aria-label="More options">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="1"></circle>
            <circle cx="19" cy="12" r="1"></circle>
            <circle cx="5" cy="12" r="1"></circle>
          </svg>
        </button>
        
        <div v-if="showMoreMenu" class="more-dropdown">
          <div class="dropdown-section">
            <button @click="navigateTo('Settings')" class="dropdown-item">
              <span>{{ t('nav.settings') }}</span>
            </button>
            <button @click="navigateTo('About')" class="dropdown-item">
              <span>{{ t('nav.about') }}</span>
            </button>
          </div>
          <div class="dropdown-divider" />
          <div class="dropdown-section">
            <button @click="logout" class="dropdown-item danger">
              <span>{{ t('nav.logout') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { useI18n } from '@/locales'
import { useCurrencies } from '@/composables/useCurrencies'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()
const { currencies } = useCurrencies()

const showMoreMenu = ref(false)

const pageTitle = computed(() => {
  const name = route.name as string
  if (!name) return t('nav.dashboard')
  const key = name.charAt(0).toLowerCase() + name.slice(1)
  return t(`nav.${key}`) || name
})

const changeCurrency = async (e: Event) => {
  const val = (e.target as HTMLSelectElement).value
  await settingsStore.updateSettings({ base_currency: val })
}

const toggleMoreMenu = () => {
  showMoreMenu.value = !showMoreMenu.value
}

const navigateTo = (name: string) => {
  router.push({ name })
  showMoreMenu.value = false
}

const logout = async () => {
  await authStore.logout()
  router.push('/login')
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.more-menu-container')) {
    showMoreMenu.value = false
  }
}

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

defineEmits(['toggle-sidebar'])
</script>

<style scoped>
@reference "tailwindcss";
.header-container {
  @apply flex items-center justify-between px-4 h-[64px];
}
.header-left {
  @apply flex items-center gap-3;
}
.menu-btn {
  @apply flex items-center justify-center p-2 rounded border-0 bg-transparent cursor-pointer transition-colors;
  color: var(--text-secondary);
}
.menu-btn:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}
.page-title {
  @apply text-xl font-bold m-0;
  color: var(--text-primary);
}
.header-right {
  @apply flex items-center gap-3;
}
.currency-selector select {
  @apply px-3 py-1.5 rounded text-sm font-medium cursor-pointer border;
  border-color: var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}
.icon-btn {
  @apply flex items-center justify-center w-10 h-10 rounded-full border-0 bg-transparent cursor-pointer transition-colors;
  color: var(--text-secondary);
}
.icon-btn:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}
.more-menu-container {
  position: relative;
}
.more-dropdown {
  position: absolute;
  top: calc(100% + 0.5rem);
  right: 0;
  width: 14rem;
  @apply rounded-lg py-2;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  box-shadow: var(--shadow-lg);
  z-index: 300;
}
.dropdown-section {
  @apply py-1;
}
.dropdown-item {
  @apply flex items-center gap-3 w-full px-4 py-2.5 text-left bg-transparent border-0 text-sm font-medium cursor-pointer transition-colors;
  color: var(--text-secondary);
}
.dropdown-item:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}
.dropdown-item.danger {
  color: var(--error);
}
.dropdown-divider {
  height: 1px;
  background-color: var(--border-primary);
  @apply my-1;
}
</style>