<template>
  <SettingSection :title="t('settings.security')">
    <SettingItem :label="t('settings.changePin')" :description="t('settings.changePinDesc')">
      <template #default>
        <BaseButton variant="primary" @click="emit('changePin')" :loading="props.totpLoading">{{ t('settings.requestToken') }}</BaseButton>
      </template>
    </SettingItem>
    <SettingItem :label="t('settings.twoFactorAuth')" :description="settingsStore.isTotpEnabled ? t('settings.twoFactorEnabled') : t('settings.twoFactorDesc')">
      <template #default>
        <BaseButton :variant="totpButtonVariant" @click="emit('toggleTotp')" :loading="props.totpLoading">{{ settingsStore.isTotpEnabled ? t('settings.disable2FA') : t('settings.enable2FA') }}</BaseButton>
      </template>
    </SettingItem>
    <SettingItem :label="t('settings.biometricAuth')" :description="status.enabled ? t('settings.biometricEnabledDesc') : t('settings.biometricDisabledDesc')">
      <template #default>
        <button @click="emit('toggleBiometric')" class="biometric-toggle" :class="{ active: status.enabled }" :disabled="!status.supported || props.biometricLoading">{{ status.enabled ? t('settings.disableBiometric') : t('settings.enableBiometric') }}</button>
        <span v-if="!status.supported" class="biometric-unsupported">{{ t('settings.biometricNotSupported') }}</span>
      </template>
    </SettingItem>
  </SettingSection>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import { BaseButton } from '@/components/ui'
import { SettingItem, SettingSection } from '@/widgets/settings'

const props = defineProps<{
  totpButtonVariant: 'primary' | 'danger'
  totpLoading: boolean
  biometricLoading: boolean
  status: any
}>()
const emit = defineEmits<{ changePin: []; toggleTotp: []; toggleBiometric: [] }>()
const settingsStore = useSettingsStore()
const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.biometric-unsupported {
  @apply text-xs ml-2 text-red-600 dark:text-red-400;
}
</style>
