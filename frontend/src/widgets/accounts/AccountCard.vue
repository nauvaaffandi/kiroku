<template>
  <div class="account-card">
    <div class="account-header">
      <div class="account-icon" :style="{ backgroundColor: account.color }">
        <span v-if="account.icon">{{ account.icon }}</span>
        <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="2" y="4" width="20" height="16" rx="2"></rect>
          <line x1="2" y1="10" x2="22" y2="10"></line>
        </svg>
      </div>
      <div class="account-badges">
        <span v-if="account.archived" class="account-type-badge archived">Archived</span>
        <div class="account-type-badge" :class="account.type">{{ account.type }}</div>
      </div>
    </div>
    <h3 class="account-name">{{ account.name }}</h3>
    <p class="account-balance">¥{{ formatNumber(account.balance) }}</p>
    <p class="account-currency">{{ account.currency }}</p>
    <div class="account-actions">
      <BaseButton variant="secondary" size="sm" @click="$emit('edit', account)">
        Edit
      </BaseButton>
      <BaseButton variant="danger" size="sm" @click="$emit('delete', account.id)">
        Delete
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import type { Account } from '@/types'

defineProps<{
  account: Account
}>()

defineEmits<{
  edit: [account: Account]
  delete: [id: string]
}>()

const formatNumber = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(num)
}
</script>

<style scoped>
@reference "tailwindcss";
.account-card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  @apply rounded-2xl p-5 flex flex-col gap-3;
}

.account-header {
  @apply flex justify-between items-center;
}

.account-icon {
  @apply w-10 h-10 rounded flex items-center justify-center text-white text-xl;
}

.account-badges {
  @apply flex gap-2 items-center;
}

.account-type-badge {
  background-color: var(--bg-primary);
  color: var(--text-muted);
  @apply text-xs px-2 py-0.5 rounded-full capitalize;
}

.account-type-badge.archived {
  color: var(--error);
}

.account-name {
  color: var(--text-primary);
  @apply text-base font-semibold m-0;
}

.account-balance {
  color: var(--text-primary);
  @apply text-[1.5rem] font-bold m-0;
}

.account-currency {
  color: var(--text-muted);
  @apply text-[0.8125rem] m-0;
}

.account-actions {
  @apply flex gap-2 mt-2;
}
</style>