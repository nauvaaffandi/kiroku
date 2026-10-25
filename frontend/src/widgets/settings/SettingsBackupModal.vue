<template>
  <BaseModal
    v-model="isOpen"
    :title="title"
    size="lg"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <div class="backup-modal">
      <div class="backup-tabs">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          :class="['tab-btn', { active: activeTab === tab.id }]"
          @click="activeTab = tab.id"
        >
          {{ tab.label }}
        </button>
      </div>

      <div v-if="activeTab === 'backup'" class="tab-content">
        <div class="backup-options">
          <h3>Create Backup</h3>
          <p>Download a complete backup of your data including transactions, accounts, categories, and settings.</p>
          <BaseButton @click="createBackup" :loading="backupLoading">
            Download Backup
          </BaseButton>
        </div>
      </div>

      <div v-else-if="activeTab === 'restore'" class="tab-content">
        <div class="restore-options">
          <h3>Restore from Backup</h3>
          <p class="warning">Warning: This will overwrite all existing data. This action cannot be undone.</p>
          <div class="form-group">
            <label>Backup File</label>
            <input
              type="file"
              accept=".json"
              @change="onFileSelect"
              class="form-input"
            />
          </div>
          <BaseButton @click="restoreBackup" :loading="restoreLoading" variant="danger" :disabled="!selectedFile">
            Restore Backup
          </BaseButton>
        </div>
      </div>

      <div v-else-if="activeTab === 'export'" class="tab-content">
        <div class="export-options">
          <h3>Export Data</h3>
          <p>Export your data in various formats for external use.</p>
          <div class="export-formats">
            <BaseButton @click="exportCSV('transactions')" :loading="exportLoading">Export Transactions (CSV)</BaseButton>
            <BaseButton @click="exportCSV('accounts')" :loading="exportLoading">Export Accounts (CSV)</BaseButton>
            <BaseButton @click="exportCSV('categories')" :loading="exportLoading">Export Categories (CSV)</BaseButton>
          </div>
        </div>
      </div>

      <div v-else-if="activeTab === 'import'" class="tab-content">
        <div class="import-options">
          <h3>Import Data</h3>
          <p>Import transactions from CSV files.</p>
          <div class="form-group">
            <label>CSV File</label>
            <input
              type="file"
              accept=".csv"
              @change="onImportFileSelect"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>Import Type</label>
            <BaseSelect
              v-model="importType"
              :options="importTypeOptions"
            />
          </div>
          <BaseButton @click="importCSV" :loading="importLoading" :disabled="!importFile">
            Import Transactions
          </BaseButton>
        </div>
      </div>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { BaseModal, BaseSelect, BaseButton } from '@/components/ui'

interface Props {
  modelValue: boolean
  title: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'create-backup': []
  'restore-backup': [file: File]
  'export-csv': [type: string]
  'import-csv': [file: File, type: string]
}>()

const isOpen = ref(props.modelValue)
const activeTab = ref('backup')
const backupLoading = ref(false)
const restoreLoading = ref(false)
const exportLoading = ref(false)
const importLoading = ref(false)
const selectedFile = ref<File | null>(null)
const importFile = ref<File | null>(null)
const importType = ref('transactions')

const tabs = [
  { id: 'backup', label: 'Backup' },
  { id: 'restore', label: 'Restore' },
  { id: 'export', label: 'Export' },
  { id: 'import', label: 'Import' }
]

const importTypeOptions = [
  { value: 'transactions', label: 'Transactions' },
  { value: 'accounts', label: 'Accounts' },
  { value: 'categories', label: 'Categories' }
]

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val) {
    activeTab.value = 'backup'
    selectedFile.value = null
    importFile.value = null
  }
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const onFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    selectedFile.value = target.files[0]
  }
}

const onImportFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    importFile.value = target.files[0]
  }
}

const createBackup = () => {
  emit('create-backup')
}

const restoreBackup = () => {
  if (selectedFile.value) {
    emit('restore-backup', selectedFile.value)
  }
}

const exportCSV = (type: string) => {
  emit('export-csv', type)
}

const importCSV = () => {
  if (importFile.value) {
    emit('import-csv', importFile.value, importType.value)
  }
}
</script>

<style scoped>
@reference "tailwindcss";
.backup-modal {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.backup-tabs {
  display: flex;
  gap: 0.5rem;
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 0.5rem;
}

.tab-btn {
  padding: 0.5rem 1rem;
  background: none;
  border: none;
  color: var(--color-text-muted);
  font-weight: 500;
  cursor: pointer;
  border-radius: var(--radius);
  transition: all 0.15s ease;
}

.tab-btn:hover {
  color: var(--color-text);
  background: var(--color-background);
}

.tab-btn.active {
  color: var(--color-primary);
  background: var(--color-primary-bg);
}

.tab-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.backup-options,
.restore-options,
.export-options,
.import-options {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.backup-options h3,
.restore-options h3,
.export-options h3,
.import-options h3 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.backup-options p,
.restore-options p,
.export-options p,
.import-options p {
  color: var(--color-text-muted);
  margin: 0;
  font-size: 0.875rem;
}

.warning {
  color: var(--color-danger) !important;
  font-weight: 500;
}

.export-formats {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.form-group label {
  font-weight: 500;
  color: var(--color-text);
}

.form-input {
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  font-family: inherit;
  color: var(--color-text);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-focus);
}
</style>