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
.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  height: 64px;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.menu-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  border-radius: var(--radius-sm);
  color: var(--color-neutral-600);
  background: none;
  border: none;
  cursor: pointer;
}
:where([data-theme="dark"], [data-theme="dark"] *) .menu-btn {
  color: #d4d4d8;
}
.page-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-neutral-900);
  margin: 0;
}
:where([data-theme="dark"], [data-theme="dark"] *) .page-title {
  color: #f4f4f5;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.currency-selector select {
  padding: 0.375rem 0.75rem;
  border: 1px solid var(--color-neutral-300);
  border-radius: var(--radius-sm);
  background-color: var(--color-neutral-100);
  color: var(--color-neutral-900);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
}
:where([data-theme="dark"], [data-theme="dark"] *) .currency-selector select {
  background-color: #121212;
  border-color: #27272a;
  color: #f4f4f5;
}
.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background: none;
  border: none;
  color: var(--color-neutral-600);
  cursor: pointer;
}
:where([data-theme="dark"], [data-theme="dark"] *) .icon-btn {
  color: #d4d4d8;
}
.more-menu-container {
  position: relative;
}
.more-dropdown {
  position: absolute;
  top: calc(100% + 0.5rem);
  right: 0;
  width: 14rem;
  background-color: var(--color-neutral-100);
  border: 1px solid var(--color-neutral-200);
  border-radius: var(--radius);
  box-shadow: var(--shadow-lg);
  padding: 0.5rem 0;
  z-index: 300;
}
:where([data-theme="dark"], [data-theme="dark"] *) .more-dropdown {
  background-color: #18181b;
  border-color: #27272a;
}
.dropdown-section {
  padding: 0.25rem 0;
}
.dropdown-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  padding: 0.625rem 1rem;
  text-align: left;
  background: none;
  border: none;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-neutral-700);
  cursor: pointer;
}
:where([data-theme="dark"], [data-theme="dark"] *) .dropdown-item {
  color: #d4d4d8;
}
.dropdown-item:hover {
  background-color: var(--color-neutral-200);
}
:where([data-theme="dark"], [data-theme="dark"] *) .dropdown-item:hover {
  background-color: #27272a;
  color: #f4f4f5;
}
.dropdown-item.danger {
  color: var(--color-danger);
}
.dropdown-divider {
  height: 1px;
  background-color: var(--color-neutral-200);
  margin: 0.25rem 0;
}
:where([data-theme="dark"], [data-theme="dark"] *) .dropdown-divider {
  background-color: #27272a;
}
@media (min-width: 768px) {
  /* Removed desktop hiding of menu button */
}
</style>