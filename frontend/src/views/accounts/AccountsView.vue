<template>
  <div class="accounts-page">
    <header class="page-header">
      <div>
        <h1>{{ t('accounts.title') }}</h1>
        <p class="page-subtitle">{{ t('accounts.subtitle') }}</p>
      </div>
      <BaseButton @click="openCreateModal">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        {{ t('accounts.addAccount') }}
      </BaseButton>
    </header>

    <AccountFilters
      v-model="localFilters"
      :searchQuery="searchQuery"
      @change="handleFilterChange"
      @search="handleSearch"
    />

    <AccountList
      :accounts="accounts"
      @edit="handleEdit"
      @delete="handleDelete"
      @create="openCreateModal"
    />

    <AccountFormModal
      v-model="showFormModal"
      :account="editingAccount"
      @submit="handleFormSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
import { useI18n } from '@/locales'
import type { Account, AccountFilters } from '@/types'
import { AccountList, AccountFormModal } from '@/widgets/accounts'
import { BaseButton } from '@/components/ui'

const accountsStore = useAccountsStore()
const { t } = useI18n()

const accounts = computed(() => accountsStore.accounts)

const showFormModal = ref(false)
const editingAccount = ref<Account | null>(null)

const localFilters = ref<AccountFilters>({
  q: '',
  type: '',
  currency: '',
  archived: '',
  sort: 'name',
  order: 'ASC'
})

const searchQuery = ref('')

let searchTimeout: ReturnType<typeof setTimeout>
watch(searchQuery, (value) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    accountsStore.filters.q = value
    fetchAccounts(1)
  }, 300)
})

const fetchAccounts = async (page = 1) => {
  await accountsStore.fetchAccounts(page)
}

onMounted(() => {
  fetchAccounts(1)
})

const handleFilterChange = () => {
  accountsStore.filters = { ...localFilters.value }
  fetchAccounts(1)
}

const handleSearch = (value: string) => {
  searchQuery.value = value
}

const handleEdit = (account: Account) => {
  editingAccount.value = account
  showFormModal.value = true
}

const handleDelete = (id: string) => {
  if (confirm('Are you sure you want to delete this account?')) {
    accountsStore.deleteAccount(id)
  }
}

const openCreateModal = () => {
  editingAccount.value = null
  showFormModal.value = true
}

const handleFormSubmit = async (data: Partial<Account>) => {
  if (editingAccount.value) {
    await accountsStore.updateAccount(editingAccount.value.id, data)
  } else {
    await accountsStore.createAccount(data)
  }
  showFormModal.value = false
  editingAccount.value = null
  await fetchAccounts(1)
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  gap: 1rem;
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.page-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin: 0.25rem 0 0;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
  flex-shrink: 0;
}
</style>