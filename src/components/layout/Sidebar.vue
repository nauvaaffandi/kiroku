<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <router-link to="/dashboard" class="logo">
        <AppIcon class="app-icon-logo" />
        <span>{{ t('common.appName') }}</span>
      </router-link>
      <button v-show="isMobile" @click="$emit('close')" class="sidebar-close" :aria-label="t('sidebar.close')">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
    <SidebarNav />
    <SidebarFooter @navigate="closeMobileSidebar" />
  </div>
</template>

<script setup lang="ts">
import AppIcon from '@/components/icons/AppIcon.vue'
import SidebarNav from './SidebarNav.vue'
import SidebarFooter from './SidebarFooter.vue'
import { useSidebarNav } from '@/composables/useSidebarNav'
import { useI18n } from '@/locales'

defineEmits(['close'])
const { isMobile, closeMobileSidebar } = useSidebarNav()
const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";
.sidebar {
  @apply flex flex-col h-full;
}

.sidebar-header {
  @apply flex items-center justify-between px-6 py-4 border-b;
  border-color: var(--border-primary);
}

.logo {
  @apply flex items-center gap-3 font-bold text-lg no-underline;
  color: var(--accent);
}

.app-icon-logo {
  width: 2rem;
  height: 2rem;
  flex-shrink: 0;
  color: var(--accent);
}

.sidebar-close {
  @apply p-2 rounded-sm cursor-pointer transition-colors;
  color: var(--text-muted);
}

.sidebar-close:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.sidebar-close:focus-visible {
  @apply outline-none ring-2 ring-[var(--accent)];
}
</style>
