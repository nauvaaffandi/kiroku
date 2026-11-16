<template>
  <SettingSection :title="t('settings.data')">
    <SettingItem :label="t('settings.downloadBackup')" :description="t('settings.backupDesc')">
      <BaseButton variant="primary" @click="$emit('downloadBackup')">{{ t('settings.downloadBackup') }}</BaseButton>
    </SettingItem>
    <SettingItem :label="t('settings.clearData')" :description="t('settings.clearDataDesc')">
      <BaseButton variant="danger" @click="$emit('clearData')">{{ t('settings.clearAllData') }}</BaseButton>
    </SettingItem>
  </SettingSection>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { BaseButton } from '@/components/ui'
import { SettingItem, SettingSection } from '@/widgets/settings'

defineProps<{}>()
const emit = defineEmits<{ downloadBackup: []; clearData: [] }>()
const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";
</style>
