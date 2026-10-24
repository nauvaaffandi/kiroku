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
        <select :value="settingsStore.settings.base_currency || 'JPY'" @change="changeCurrency" aria-label="Base currency">
          <option value="JPY">JPY ¥</option>
          <option value="USD">USD $</option>
          <option value="EUR">EUR €</option>
          <option value="GBP">GBP £</option>
        </select>
      </div>
      
      <button @click="toggleTheme" class="icon-btn" aria-label="Toggle theme">
        <svg v-if="isDark" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="5"></circle>
          <line x1="12" y1="1" x2="12" y2="3"></line>
          <line x1="12" y1="21" x2="12" y2="23"></line>
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
          <line x1="1" y1="12" x2="3" y2="12"></line>
          <line x1="21" y1="12" x2="23" y2="12"></line>
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
        </svg>
        <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
        </svg>
      </button>
      
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

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()

const showMoreMenu = ref(false)
const isDark = ref(false)

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

const toggleTheme = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  showMoreMenu.value = false
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
  isDark.value = document.documentElement.classList.contains('dark') || localStorage.getItem('theme') === 'dark'
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  }
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
.page-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-neutral-900);
  margin: 0;
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
  cursor: pointer;
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
  color: var(--color-neutral-700);
  cursor: pointer;
}
.dropdown-item:hover {
  background-color: var(--color-neutral-200);
}
.dropdown-item.danger {
  color: var(--color-danger);
}
.dropdown-divider {
  height: 1px;
  background-color: var(--color-neutral-200);
  margin: 0.25rem 0;
}
@media (min-width: 768px) {
  /* Removed desktop hiding of menu button */
}
</style>