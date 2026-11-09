<template>
  <div class="sidebar-footer">
    <div class="user-info" v-if="authStore.user">
      <div class="user-avatar">
        {{ authStore.user.user_id.charAt(0).toUpperCase() }}
      </div>
      <div class="user-details">
        <span class="user-name text-zinc-900 dark:text-white font-bold">{{ authStore.user.user_id.slice(0, 8) }}</span>
        <span class="user-role text-zinc-500 dark:text-zinc-300 font-medium">User</span>
      </div>
    </div>
    <router-link to="/settings" class="footer-link text-zinc-600 dark:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800 font-semibold" @click="$emit('navigate')">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="color: #64748b">
        <circle cx="12" cy="12" r="3"></circle>
        <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06-.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
      </svg>
      <span>Settings</span>
    </router-link>
    <ThemeToggle />
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import ThemeToggle from '@/components/common/ThemeToggle.vue'

defineEmits(['navigate'])
const authStore = useAuthStore()
</script>

<style scoped>
@reference "tailwindcss";

.sidebar-footer {
  @apply px-6 py-4 border-t;
  border-color: var(--border-primary);
}

.user-info {
  @apply flex items-center gap-3 p-2 rounded-sm mb-4;
  background-color: var(--bg-card);
}

.user-avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background-color: var(--accent);
  color: white;
  @apply flex items-center justify-center font-semibold text-sm;
}

.user-details {
  @apply flex-1 min-w-0 flex flex-col;
}

.user-name {
  @apply text-sm font-medium truncate;
  color: var(--text-primary);
}

.user-role {
  @apply text-xs;
  color: var(--text-muted);
}

.footer-link {
  @apply flex items-center gap-3 px-4 py-3 rounded-sm no-underline transition-colors;
  color: var(--text-secondary);
}

.footer-link:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.footer-link svg {
  flex-shrink: 0;
}
</style>
