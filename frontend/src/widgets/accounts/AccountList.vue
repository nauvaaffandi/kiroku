<template>
  <div class="accounts-list">
    <div v-if="accounts.length === 0" class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="empty-icon">
        <rect x="2" y="4" width="20" height="16" rx="2"></rect>
        <line x1="2" y1="10" x2="22" y2="10"></line>
      </svg>
      <p>No accounts yet</p>
      <BaseButton @click="$emit('create')">Add Account</BaseButton>
    </div>

    <div v-else class="accounts-grid">
      <AccountCard
        v-for="account in accounts"
        :key="account.id"
        :account="account"
        @edit="$emit('edit', $event)"
        @delete="$emit('delete', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import AccountCard from './AccountCard.vue'
import type { Account } from '@/types'

defineProps<{
  accounts: Account[]
}>()

defineEmits<{
  edit: [account: Account]
  delete: [id: string]
  create: []
}>()
</script>

<style scoped>
@reference "tailwindcss";

.accounts-grid {
  @apply grid gap-4;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

.empty-state {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  color: var(--text-muted);
  border-radius: var(--radius-lg);
  @apply text-center py-12 px-6;
}

.empty-icon {
  @apply mb-4 opacity-50;
}

.empty-state p {
  @apply mb-4;
}
</style>