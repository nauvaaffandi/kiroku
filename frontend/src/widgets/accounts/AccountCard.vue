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
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.account-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.account-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.25rem;
}

.account-badges {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.account-type-badge {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  background: var(--color-background);
  color: var(--color-text-muted);
  text-transform: capitalize;
}

.account-type-badge.archived {
  background: var(--color-danger-bg);
  color: var(--color-danger);
}

.account-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.account-balance {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.account-currency {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  margin: 0;
}

.account-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}
</style>